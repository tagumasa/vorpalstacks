package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_s3"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"

	"github.com/google/uuid"
)

func objectBucketName(region string) string {
	return "s3_objects-" + region
}

func multipartBucketName(region string) string {
	return "s3_multipart-" + region
}

const (
	// keySep is the internal delimiter for Pebble storage keys. Using \x00
	// avoids collisions with S3 object keys, which cannot contain null bytes.
	keySep = "\x00"
	// restoreIndexPrefix keys the restore index, which tracks every object
	// version with an active temporary restored copy so the expiry sweep
	// scans the index instead of listing all bucket versions. The \x01
	// marker keeps index keys disjoint from object records, whose keys
	// start with a bucket name (bucket names cannot contain control or
	// null bytes).
	restoreIndexPrefix = "\x01restore" + keySep
)

var (
	// ErrObjectNotFound is returned when the specified object does not exist.
	ErrObjectNotFound = common.NewStoreError("s3", "object_not_found", common.ErrNotFound)
	// ErrBucketHasObjects is returned when attempting to delete a bucket that contains objects.
	ErrBucketHasObjects = common.NewStoreError("s3", "bucket_has_objects", common.ErrConflict)
	// ErrUploadNotFound is returned when the specified multipart upload does not exist.
	ErrUploadNotFound = common.NewStoreError("s3", "upload_not_found", common.ErrNotFound)
)

// ObjectStore manages S3 object storage and retrieval.
type ObjectStore struct {
	*common.BaseStore
	storage         storage.BasicStorage
	blobStore       storage.BlobStore
	arnBuilder      *svcarn.S3Builder
	bucketStore     *BucketStore
	keyLocker       common.KeyLocker
	versioningCache *VersioningCache
	region          string
}

// NewObjectStore creates a new ObjectStore instance.
func NewObjectStore(store storage.BasicStorage, blobStore storage.BlobStore, bucketStore *BucketStore, accountId, region string) (*ObjectStore, error) {
	cache, err := NewVersioningCache()
	if err != nil {
		slog.Warn("s3: versioning cache creation failed, retrying in 500ms", "error", err)
		time.Sleep(500 * time.Millisecond)
		cache, err = NewVersioningCache()
		if err != nil {
			return nil, fmt.Errorf("failed to create versioning cache after retry: %w", err)
		}
	}

	os := &ObjectStore{
		BaseStore:       common.NewBaseStore(store.Bucket(objectBucketName(region)), "s3"),
		storage:         store,
		blobStore:       blobStore,
		bucketStore:     bucketStore,
		arnBuilder:      svcarn.NewARNBuilder(accountId, region).S3(),
		versioningCache: cache,
		region:          region,
	}

	bucketStore.SetVersioningCallback(func(bucket string, enabled bool) {
		cache.Set(bucket, enabled)
	})

	bucketStore.SetOnDeleteCallback(func(bucket string) {
		cache.Delete(bucket)
		os.keyLocker.DeleteByPrefix(bucket + keySep)
	})

	return os, nil
}

// Close closes the object store and releases resources.
func (s *ObjectStore) Close() {
	if s.versioningCache != nil {
		s.versioningCache.Close()
	}
}

func (s *ObjectStore) versionedStorageKey(bucket, key, versionId string) string {
	if versionId == "" {
		versionId = "null"
	}
	return bucket + keySep + key + keySep + versionId
}

func (s *ObjectStore) latestKeyStorageKey(bucket, key string) string {
	return bucket + keySep + key + keySep + "_latest"
}

// restoreIndexKey builds the restore-index key for one object version. An
// empty versionId is normalised to "null", matching the object record key.
func restoreIndexKey(bucket, key, versionId string) string {
	if versionId == "" {
		versionId = "null"
	}
	return restoreIndexPrefix + bucket + keySep + key + keySep + versionId
}

// parseRestoreIndexKey decodes a restore-index key into its entry. It
// reports false for keys that do not carry the index layout.
func parseRestoreIndexKey(key string) (RestoreIndexEntry, bool) {
	rest := strings.TrimPrefix(key, restoreIndexPrefix)
	if len(rest) == len(key) {
		return RestoreIndexEntry{}, false
	}
	parts := strings.SplitN(rest, keySep, 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return RestoreIndexEntry{}, false
	}
	return RestoreIndexEntry{Bucket: parts[0], Key: parts[1], VersionID: parts[2]}, true
}

func (s *ObjectStore) generateVersionId() string {
	return uuid.New().String()
}

func (s *ObjectStore) isVersioningEnabled(bucket string) bool {
	if enabled, ok := s.versioningCache.Get(bucket); ok {
		return enabled
	}
	b, err := s.bucketStore.Get(bucket)
	if err != nil {
		return false
	}
	enabled := b.VersioningStatus == BucketVersioningEnabled
	s.versioningCache.Set(bucket, enabled)
	return enabled
}

// resolveObjectMetaPB loads the protobuf record for a bucket/key/versionId
// request using the shared version-resolution rules: an explicit version ID
// reads that version's record; an implicit request on a versioned bucket
// reads the latest pointer and falls back to the pre-versioning null-version
// record when the pointer is absent; a non-versioned bucket reads the
// null-version record directly. Delete markers surface as ErrObjectNotFound.
func (s *ObjectStore) resolveObjectMetaPB(bucket, key, versionId string) (*pb.Object, error) {
	isVersioned := s.isVersioningEnabled(bucket)
	effectiveVersionId := versionId
	if !isVersioned && versionId == "null" {
		effectiveVersionId = ""
	}

	var pbObj pb.Object
	if effectiveVersionId != "" {
		if err := s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, effectiveVersionId), &pbObj); err != nil {
			return nil, ErrObjectNotFound
		}
	} else if isVersioned {
		if err := s.BaseStore.GetProto(s.latestKeyStorageKey(bucket, key), &pbObj); err != nil {
			// Fallback: object may predate versioning enablement, in which
			// case only the null-version record exists.
			if err2 := s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, "null"), &pbObj); err2 != nil {
				return nil, ErrObjectNotFound
			}
		}
	} else {
		if err := s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, "null"), &pbObj); err != nil {
			return nil, ErrObjectNotFound
		}
	}
	if pbObj.IsDeleteMarker {
		return nil, ErrObjectNotFound
	}
	return &pbObj, nil
}

func (s *ObjectStore) multipartKey(uploadId string) string {
	return uploadId
}

func multipartIndexBucketName(region string) string {
	return "s3_multipart_index-" + region
}

func (s *ObjectStore) multipartIndexKey(bucket, key, uploadId string) string {
	return bucket + keySep + key + keySep + uploadId
}

func validateS3Key(key string) error {
	if strings.Contains(key, "..") {
		return fmt.Errorf("invalid object key: path traversal detected")
	}
	if strings.Contains(key, "\x00") {
		return fmt.Errorf("invalid object key: null byte detected")
	}
	return nil
}

// Get retrieves an object from the store.
func (s *ObjectStore) Get(ctx context.Context, bucket, key string) (io.ReadCloser, *Object, error) {
	return s.GetWithVersion(ctx, bucket, key, "")
}

// GetMetadata retrieves metadata for an object.
func (s *ObjectStore) GetMetadata(bucket, key string) (*Object, error) {
	var obj pb.Object

	if s.isVersioningEnabled(bucket) {
		latestKey := s.latestKeyStorageKey(bucket, key)
		if err := s.BaseStore.GetProto(latestKey, &obj); err != nil {
			// Fallback: object may have been created before versioning was enabled.
			nullKey := s.versionedStorageKey(bucket, key, "null")
			if err2 := s.BaseStore.GetProto(nullKey, &obj); err2 != nil {
				return nil, err2
			}
		}
	} else {
		if err := s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, "null"), &obj); err != nil {
			return nil, err
		}
	}

	if obj.IsDeleteMarker {
		return nil, ErrObjectNotFound
	}
	return ProtoToObject(&obj), nil
}

// Put stores an object in the store.
func (s *ObjectStore) Put(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string) (*Object, error) {
	return s.PutWithVersioning(ctx, bucket, key, reader, contentType, metadata, false, StorageClassStandard, nil)
}

// Delete removes an object from the store.
func (s *ObjectStore) Delete(ctx context.Context, bucket, key string) error {
	_, err := s.DeleteWithVersion(ctx, bucket, key, "")
	return err
}

// Exists checks whether an object exists in the store.
func (s *ObjectStore) Exists(ctx context.Context, bucket, key string) (bool, error) {
	return s.blobStore.Exists(ctx, bucket, key)
}

// Head retrieves metadata for an object without the body content.
func (s *ObjectStore) Head(ctx context.Context, bucket, key string) (*Object, error) {
	return s.HeadWithVersion(ctx, bucket, key, "")
}

// SystemMetadata holds S3 object system-level metadata (content-type, size, etc.).
type SystemMetadata struct {
	ContentEncoding    string
	ContentLanguage    string
	ContentDisposition string
	CacheControl       string
}

func newObject(key, bucket, contentType string, metadata map[string]string, versionId string, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata) *Object {
	sc := storageClass
	if sc == "" {
		sc = StorageClassStandard
	}
	obj := &Object{
		Key:            key,
		BucketName:     bucket,
		Size:           0,
		ETag:           "",
		LastModified:   time.Now().UTC(),
		ContentType:    contentType,
		Metadata:       metadata,
		StorageClass:   sc,
		IsLatest:       true,
		IsDeleteMarker: isDeleteMarker,
		VersionID:      versionId,
	}
	if sysMeta != nil {
		obj.ContentEncoding = sysMeta.ContentEncoding
		obj.ContentLanguage = sysMeta.ContentLanguage
		obj.ContentDisposition = sysMeta.ContentDisposition
		obj.CacheControl = sysMeta.CacheControl
	}
	return obj
}

// SetStorageClass updates the storage class of an object.
func (s *ObjectStore) SetStorageClass(bucket, key, versionId string, storageClass ObjectStorageClass) error {
	return s.keyLocker.WithLock(bucket+keySep+key, func() error {
		isVersioned := s.isVersioningEnabled(bucket)

		var storageKey string
		if versionId != "" {
			storageKey = s.versionedStorageKey(bucket, key, versionId)
		} else if isVersioned {
			storageKey = s.latestKeyStorageKey(bucket, key)
		} else {
			storageKey = s.versionedStorageKey(bucket, key, "null")
		}

		var pbObj pb.Object
		if err := s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			// Fallback: object may predate versioning enablement and have no
			// _latest pointer. Mirror SetACLWithVersion behaviour.
			if isVersioned && versionId == "" {
				nullKey := s.versionedStorageKey(bucket, key, "null")
				if err2 := s.BaseStore.GetProto(nullKey, &pbObj); err2 != nil {
					return ErrObjectNotFound
				}
				storageKey = nullKey
			} else {
				return err
			}
		}
		pbObj.StorageClass = objectStorageClassToProto(storageClass)
		obj := ProtoToObject(&pbObj)

		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return err
		}

		// When versioning is enabled, keep the versioned record and the
		// _latest pointer in sync, mirroring updateObjectLockMetadata.
		if isVersioned {
			vid := versionId
			if vid == "" {
				vid = obj.VersionID
			}
			versionedKey := s.versionedStorageKey(bucket, key, vid)
			if versionedKey != storageKey {
				if err := s.BaseStore.PutProto(versionedKey, ObjectToProto(obj)); err != nil {
					return err
				}
			}
			if obj.IsLatest {
				latestKey := s.latestKeyStorageKey(bucket, key)
				if latestKey != storageKey {
					if err := s.BaseStore.PutProto(latestKey, ObjectToProto(obj)); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

// RestoreIndexEntry identifies one object version with an active restore,
// as recorded in the restore index.
type RestoreIndexEntry struct {
	Bucket    string
	Key       string
	VersionID string
	Expiry    time.Time
}

// SetRestoreState records the expiry of the temporary restored copy of an
// archived object; a nil expiry clears the restored state. The object's
// storage class is never modified: an archived object stays in its archive
// class while restored, and the expiry alone gates reads of the copy. The
// restore index is written and cleared alongside the object record — in a
// single atomic batch — so the expiry sweep only ever visits objects with
// an active restore and the record and index can never diverge.
func (s *ObjectStore) SetRestoreState(bucket, key, versionId string, expiry *time.Time) error {
	return s.keyLocker.WithLock(bucket+keySep+key, func() error {
		isVersioned := s.isVersioningEnabled(bucket)

		var storageKey string
		if versionId != "" {
			storageKey = s.versionedStorageKey(bucket, key, versionId)
		} else if isVersioned {
			storageKey = s.latestKeyStorageKey(bucket, key)
		} else {
			storageKey = s.versionedStorageKey(bucket, key, "null")
		}

		var pbObj pb.Object
		if err := s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			if isVersioned && versionId == "" {
				nullKey := s.versionedStorageKey(bucket, key, "null")
				if err2 := s.BaseStore.GetProto(nullKey, &pbObj); err2 != nil {
					if expiry == nil {
						// The object record is gone; forget the restore
						// state entirely instead of failing the sweep.
						return s.BaseStore.Delete(restoreIndexKey(bucket, key, versionId))
					}
					return ErrObjectNotFound
				}
				storageKey = nullKey
			} else {
				if expiry == nil {
					return s.BaseStore.Delete(restoreIndexKey(bucket, key, versionId))
				}
				return err
			}
		}

		if expiry != nil {
			pbObj.RestoreExpiry = timestamppb.New(*expiry)
		} else {
			pbObj.RestoreExpiry = nil
		}
		obj := ProtoToObject(&pbObj)

		vid := versionId
		if vid == "" {
			vid = obj.VersionID
		}

		// Buffer every mutation — the addressed record, the restore index
		// entry, and the shadow copies of the record — and commit them in
		// one batch so a crash between writes can never leave the object
		// record and the restore index out of step.
		batchBucket, ok := s.BaseStore.Bucket().(storage.BatchBucket)
		if !ok {
			return fmt.Errorf("s3: storage bucket does not support atomic batches")
		}
		batch := batchBucket.NewBatch()
		defer batch.Close()

		objBytes, err := proto.Marshal(ObjectToProto(obj))
		if err != nil {
			return err
		}
		if err := batch.Put([]byte(storageKey), objBytes); err != nil {
			return err
		}

		indexKey := restoreIndexKey(bucket, key, vid)
		if expiry != nil {
			if err := batch.Put([]byte(indexKey), []byte(strconv.FormatInt(expiry.UnixNano(), 10))); err != nil {
				return err
			}
		} else if err := batch.Delete([]byte(indexKey)); err != nil {
			return err
		}

		if isVersioned {
			versionedKey := s.versionedStorageKey(bucket, key, vid)
			if versionedKey != storageKey {
				if err := batch.Put([]byte(versionedKey), objBytes); err != nil {
					return err
				}
			}
			if obj.IsLatest {
				latestKey := s.latestKeyStorageKey(bucket, key)
				if latestKey != storageKey {
					if err := batch.Put([]byte(latestKey), objBytes); err != nil {
						return err
					}
				}
			}
		}

		return batch.Commit()
	})
}

// ActiveRestores lists every object version with an active restore, in
// index order. The restore expiry sweep consumes this list instead of
// listing all versions of every bucket.
func (s *ObjectStore) ActiveRestores() ([]RestoreIndexEntry, error) {
	var entries []RestoreIndexEntry
	err := s.BaseStore.ScanPrefix(restoreIndexPrefix, func(key string, value []byte) error {
		entry, ok := parseRestoreIndexKey(key)
		if !ok {
			return nil
		}
		nanos, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return nil
		}
		entry.Expiry = time.Unix(0, nanos)
		entries = append(entries, entry)
		return nil
	})
	return entries, err
}

// SetReplicationStatus updates the replication status of an object.
// Valid statuses are "PENDING", "COMPLETED", and "FAILED".
func (s *ObjectStore) SetReplicationStatus(bucket, key, versionId, status string) error {
	return s.keyLocker.WithLock(bucket+keySep+key, func() error {
		isVersioned := s.isVersioningEnabled(bucket)

		var storageKey string
		if versionId != "" {
			storageKey = s.versionedStorageKey(bucket, key, versionId)
		} else if isVersioned {
			storageKey = s.latestKeyStorageKey(bucket, key)
		} else {
			storageKey = s.versionedStorageKey(bucket, key, "null")
		}

		var pbObj pb.Object
		if err := s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			if isVersioned && versionId == "" {
				nullKey := s.versionedStorageKey(bucket, key, "null")
				if err2 := s.BaseStore.GetProto(nullKey, &pbObj); err2 != nil {
					return ErrObjectNotFound
				}
				storageKey = nullKey
			} else {
				return err
			}
		}
		pbObj.ReplicationStatus = status
		obj := ProtoToObject(&pbObj)

		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return err
		}

		if isVersioned {
			vid := versionId
			if vid == "" {
				vid = obj.VersionID
			}
			versionedKey := s.versionedStorageKey(bucket, key, vid)
			if versionedKey != storageKey {
				if err := s.BaseStore.PutProto(versionedKey, ObjectToProto(obj)); err != nil {
					return err
				}
			}
			if obj.IsLatest {
				latestKey := s.latestKeyStorageKey(bucket, key)
				if latestKey != storageKey {
					if err := s.BaseStore.PutProto(latestKey, ObjectToProto(obj)); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}
