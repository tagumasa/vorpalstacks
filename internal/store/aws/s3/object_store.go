package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

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

	os.migrateKeyDelimiter(store, region)

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
		if err := s.BaseStore.GetProto(s.latestKeyStorageKey(bucket, key), &obj); err != nil {
			return nil, err
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
		var storageKey string
		if versionId != "" {
			storageKey = s.versionedStorageKey(bucket, key, versionId)
		} else if s.isVersioningEnabled(bucket) {
			storageKey = s.latestKeyStorageKey(bucket, key)
		} else {
			storageKey = s.versionedStorageKey(bucket, key, "null")
		}

		var obj pb.Object
		if err := s.BaseStore.GetProto(storageKey, &obj); err != nil {
			return err
		}
		obj.StorageClass = objectStorageClassToProto(storageClass)
		return s.BaseStore.PutProto(storageKey, &obj)
	})
}

const migrationFlagKey = "__key_delimiter_migrated_v1"

func (s *ObjectStore) migrateKeyDelimiter(store storage.BasicStorage, region string) {
	objBucket := store.Bucket(objectBucketName(region))
	mpIdxBucket := store.Bucket(multipartIndexBucketName(region))

	flag, _ := objBucket.Get([]byte(migrationFlagKey))
	if flag != nil {
		return
	}

	oldSep := "#"
	migrated := 0

	for _, bucket := range []storage.Bucket{objBucket, mpIdxBucket} {
		iter := bucket.ScanPrefix(nil)
		for iter.Next() {
			k := iter.Key()
			if string(k) == migrationFlagKey {
				continue
			}
			if !bytes.Contains(k, []byte(oldSep)) {
				continue
			}
			v := iter.Value()
			newKey := bytes.ReplaceAll(k, []byte(oldSep), []byte(keySep))
			if err := bucket.Put(newKey, v); err != nil {
				slog.Error("s3 migration: failed to write new key", "old", string(k), "error", err)
				continue
			}
			if err := bucket.Delete(k); err != nil {
				slog.Error("s3 migration: failed to delete old key", "old", string(k), "error", err)
			}
			migrated++
		}
		iter.Close()
	}

	objBucket.Put([]byte(migrationFlagKey), []byte("1"))

	if migrated > 0 {
		slog.Info("s3: migrated key delimiter", "region", region, "keys", migrated)
	}
}
