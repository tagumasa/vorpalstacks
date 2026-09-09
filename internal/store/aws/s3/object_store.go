package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

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
	// ErrVersioningNotEnabled is returned when an operation that requires a
	// specific object version runs against a bucket whose versioning was
	// never enabled.
	ErrVersioningNotEnabled = common.NewStoreError("s3", "versioning_not_enabled", common.ErrInvalidState)
	// ErrRetentionNotFound is returned when an object version carries no
	// retention configuration — the NoSuchObjectRetention condition.
	ErrRetentionNotFound = common.NewStoreError("s3", "retention_not_found", common.ErrNotFound)
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

// recordNotFoundOrErr maps a record read to its caller-facing error: an
// absent key is the not-found sentinel, while any other failure (I/O,
// corruption, unmarshal) propagates — a store fault must never surface as
// a missing object.
func recordNotFoundOrErr(err error) error {
	if errors.Is(err, common.ErrNotFound) {
		return ErrObjectNotFound
	}
	return err
}

// resolveObjectRecordPB loads the protobuf record addressed by the layout
// rules and reports the storage key it was read from, so mutation paths
// write the record back where it lives. The layout decides, not the
// bucket's current versioning status: an explicit versionId addresses
// that record directly (the null version included), while the current
// version is the "_latest" pointer when the key carries versioned layout
// and the null record otherwise — suspension rewrites nothing, and a
// suspended write replaces the pointer with its null record. Unlike
// resolveObjectMetaPB it does not reject delete markers; the caller
// decides what a marker means for its operation.
func (s *ObjectStore) resolveObjectRecordPB(bucket, key, versionId string) (*pb.Object, string, error) {
	var pbObj pb.Object
	if versionId != "" {
		storageKey := s.versionedStorageKey(bucket, key, versionId)
		if err := s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			return nil, "", recordNotFoundOrErr(err)
		}
		return &pbObj, storageKey, nil
	}
	latestKey := s.latestKeyStorageKey(bucket, key)
	if err := s.BaseStore.GetProto(latestKey, &pbObj); err != nil {
		if !errors.Is(err, common.ErrNotFound) {
			return nil, "", err
		}
		// Fallback: the key may predate versioning enablement or carry
		// a suspended write — only the null-version record exists.
		nullKey := s.versionedStorageKey(bucket, key, "null")
		if err2 := s.BaseStore.GetProto(nullKey, &pbObj); err2 != nil {
			return nil, "", recordNotFoundOrErr(err2)
		}
		return &pbObj, nullKey, nil
	}
	return &pbObj, latestKey, nil
}

// resolveObjectMetaPB reads one object record through the shared layout
// resolution (see resolveObjectRecordPB). Delete markers surface as
// ErrObjectNotFound: a marker is not a readable object.
func (s *ObjectStore) resolveObjectMetaPB(bucket, key, versionId string) (*pb.Object, error) {
	pbObj, _, err := s.resolveObjectRecordPB(bucket, key, versionId)
	if err != nil {
		return nil, err
	}
	if pbObj.IsDeleteMarker {
		return nil, ErrObjectNotFound
	}
	return pbObj, nil
}

// mutateObjectRecord is the single mutation skeleton for object metadata:
// it resolves the addressed record (explicit versionId, or the current
// version by layout), rejects delete markers, applies mutate, and commits
// the mutated record together with its shadow copies in one atomic batch —
// a crash between writes can never leave the record and its shadows out of
// step. Callers that need extra writes in the same commit (the restore
// index) use enqueueRecordWithShadows directly instead.
func (s *ObjectStore) mutateObjectRecord(bucket, key, versionId string, mutate func(*Object) error) error {
	lockKey := bucket + keySep + key
	return s.keyLocker.WithLock(lockKey, func() error {
		pbObj, storageKey, err := s.resolveObjectRecordPB(bucket, key, versionId)
		if err != nil {
			return err
		}
		if pbObj.IsDeleteMarker {
			return ErrObjectNotFound
		}

		obj := ProtoToObject(pbObj)
		if err := mutate(obj); err != nil {
			return err
		}

		batchBucket, ok := s.BaseStore.Bucket().(storage.BatchBucket)
		if !ok {
			return fmt.Errorf("s3: storage bucket does not support atomic batches")
		}
		batch := batchBucket.NewBatch()
		defer batch.Close()

		if err := s.enqueueRecordWithShadows(batch, bucket, key, storageKey, versionId, obj); err != nil {
			return err
		}
		return batch.Commit()
	})
}

// enqueueRecordWithShadows buffers obj's record and its shadow copies into
// batch: the addressed storage key, the versioned record when the address
// was the "_latest" pointer, and the pointer itself when the record is the
// current version and the pointer already exists — a suspended key never
// resurrects its pointer, and a never-versioned key never grows one.
func (s *ObjectStore) enqueueRecordWithShadows(batch storage.Batch, bucket, key, storageKey, versionId string, obj *Object) error {
	objBytes, err := proto.Marshal(ObjectToProto(obj))
	if err != nil {
		return err
	}
	if err := batch.Put([]byte(storageKey), objBytes); err != nil {
		return err
	}
	vid := versionId
	if vid == "" {
		vid = obj.VersionID
	}
	versionedKey := s.versionedStorageKey(bucket, key, vid)
	if versionedKey != storageKey {
		if err := batch.Put([]byte(versionedKey), objBytes); err != nil {
			return err
		}
	}
	latestKey := s.latestKeyStorageKey(bucket, key)
	if obj.IsLatest && latestKey != storageKey && s.BaseStore.Exists(latestKey) {
		if err := batch.Put([]byte(latestKey), objBytes); err != nil {
			return err
		}
	}
	return nil
}

func multipartIndexBucketName(region string) string {
	return "s3_multipart_index-" + region
}

func (s *ObjectStore) multipartIndexKey(bucket, key, uploadId string) string {
	return bucket + keySep + key + keySep + uploadId
}

// Get retrieves an object from the store.
func (s *ObjectStore) Get(ctx context.Context, bucket, key string) (io.ReadCloser, *Object, error) {
	return s.GetWithVersion(ctx, bucket, key, "")
}

// GetMetadata retrieves metadata for an object.
func (s *ObjectStore) GetMetadata(bucket, key string) (*Object, error) {
	return s.getVersionedObjectMeta(bucket, key, "")
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

// GetRange retrieves a range of bytes from an object.
func (s *ObjectStore) GetRange(ctx context.Context, bucket, key string, offset, length int64) (io.ReadCloser, *Object, error) {
	return s.GetRangeWithVersion(ctx, bucket, key, "", offset, length)
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
	return s.mutateObjectRecord(bucket, key, versionId, func(obj *Object) error {
		obj.StorageClass = storageClass
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
	lockKey := bucket + keySep + key
	return s.keyLocker.WithLock(lockKey, func() error {
		pbObj, storageKey, err := s.resolveObjectRecordPB(bucket, key, versionId)
		if err != nil {
			if expiry == nil {
				// The object record is gone; forget the restore
				// state entirely instead of failing the sweep.
				return s.BaseStore.Delete(restoreIndexKey(bucket, key, versionId))
			}
			return ErrObjectNotFound
		}

		obj := ProtoToObject(pbObj)
		obj.RestoreExpiry = expiry

		vid := versionId
		if vid == "" {
			vid = obj.VersionID
		}

		batchBucket, ok := s.BaseStore.Bucket().(storage.BatchBucket)
		if !ok {
			return fmt.Errorf("s3: storage bucket does not support atomic batches")
		}
		batch := batchBucket.NewBatch()
		defer batch.Close()

		if err := s.enqueueRecordWithShadows(batch, bucket, key, storageKey, versionId, obj); err != nil {
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

// Replication status values carried by the x-amz-replication-status
// header: source-side PENDING/COMPLETED/FAILED, destination-side REPLICA.
const (
	ReplicationStatusPending   = "PENDING"
	ReplicationStatusCompleted = "COMPLETED"
	ReplicationStatusFailed    = "FAILED"
	ReplicationStatusReplica   = "REPLICA"
)

// SetReplicationStatus updates the replication status of an object.
// On a replication source the valid statuses are "PENDING", "COMPLETED",
// and "FAILED"; on a replication destination the status is "REPLICA",
// marking the object as a copy rather than an original upload.
func (s *ObjectStore) SetReplicationStatus(bucket, key, versionId, status string) error {
	return s.mutateObjectRecord(bucket, key, versionId, func(obj *Object) error {
		obj.ReplicationStatus = status
		return nil
	})
}
