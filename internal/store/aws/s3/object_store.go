package s3

import (
	"context"
	"fmt"
	"io"
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
		os.keyLocker.DeleteByPrefix(bucket + "#")
	})

	return os, nil
}

// Close closes the object store and releases resources.
func (s *ObjectStore) Close() {
	if s.versioningCache != nil {
		s.versioningCache.Close()
	}
}

func (s *ObjectStore) storageKey(bucket, key string) string {
	return bucket + "#" + key
}

func (s *ObjectStore) versionedStorageKey(bucket, key, versionId string) string {
	if versionId == "" {
		versionId = "null"
	}
	return bucket + "#" + key + "#" + versionId
}

func (s *ObjectStore) latestKeyStorageKey(bucket, key string) string {
	return bucket + "#" + key + "#_latest"
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
	return bucket + "#" + key + "#" + uploadId
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

func (s *ObjectStore) wasVersioned(bucket string) bool {
	b, err := s.bucketStore.Get(bucket)
	if err != nil {
		return false
	}
	return b.VersioningStatus == BucketVersioningEnabled || b.VersioningStatus == "Suspended"
}

// GetMetadata retrieves metadata for an object.
func (s *ObjectStore) GetMetadata(bucket, key string) (*Object, error) {
	var obj pb.Object
	var err error

	b, bucketErr := s.bucketStore.Get(bucket)
	versioningEnabled := bucketErr == nil && b.VersioningStatus == BucketVersioningEnabled
	versioningSuspended := bucketErr == nil && b.VersioningStatus == "Suspended"

	if versioningEnabled {
		storageKey := s.latestKeyStorageKey(bucket, key)
		if err = s.BaseStore.GetProto(storageKey, &obj); err != nil {
			return nil, err
		}
		if obj.IsDeleteMarker {
			return nil, ErrObjectNotFound
		}
	} else if versioningSuspended {
		latestKey := s.latestKeyStorageKey(bucket, key)
		if err = s.BaseStore.GetProto(latestKey, &obj); err != nil {
			storageKey := s.storageKey(bucket, key)
			if err = s.BaseStore.GetProto(storageKey, &obj); err != nil {
				return nil, err
			}
		}
		if obj.IsDeleteMarker {
			return nil, ErrObjectNotFound
		}
	} else {
		storageKey := s.storageKey(bucket, key)
		if err = s.BaseStore.GetProto(storageKey, &obj); err != nil {
			return nil, err
		}
	}
	return ProtoToObject(&obj), nil
}

// Put stores an object in the store.
func (s *ObjectStore) Put(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string) (*Object, error) {
	return s.PutWithVersioning(ctx, bucket, key, reader, contentType, metadata, false)
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

func newObject(key, bucket, contentType string, metadata map[string]string, versionId string, isDeleteMarker bool) *Object {
	return &Object{
		Key:            key,
		BucketName:     bucket,
		Size:           0,
		ETag:           "",
		LastModified:   time.Now().UTC(),
		ContentType:    contentType,
		Metadata:       metadata,
		StorageClass:   StorageClassStandard,
		IsLatest:       true,
		IsDeleteMarker: isDeleteMarker,
		VersionID:      versionId,
	}
}
