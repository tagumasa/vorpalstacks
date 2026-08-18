package s3

import (
	"context"
	"io"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

// BucketStoreInterface defines operations for managing S3 buckets.
type BucketStoreInterface interface {
	Get(name string) (*Bucket, error)
	Create(name, region string) (*Bucket, error)
	Put(bucket *Bucket) error
	Delete(name string) error
	Exists(name string) bool
	List() ([]*Bucket, error)
	SetVersioning(name string, status BucketVersioningStatus, mfaDelete string) error
	SetVersioningCallback(fn func(bucket string, enabled bool))
	SetOnDeleteCallback(fn func(bucket string))
	SetEncryption(name string, config *EncryptionConfig) error
	GetEncryptionConfiguration(name string) (*EncryptionConfig, error)
	SetPolicy(name, policy string) error
	SetCORS(name string, config *CORSConfiguration) error
	SetPublicAccessBlock(name string, config *PublicAccessBlockConfig) error
	SetTags(name string, tags []types.Tag) error
	SetLifecycleConfiguration(name string, config *LifecycleConfiguration) error
	SetWebsiteConfiguration(name string, config *WebsiteConfiguration) error
	SetObjectLockConfiguration(name string, config *ObjectLockConfiguration) error
	SetNotificationConfiguration(name string, config *NotificationConfiguration) error
	GetNotificationConfiguration(name string) (*NotificationConfiguration, error)
	SetLoggingConfiguration(name string, config *LoggingConfiguration) error
	GetLoggingConfiguration(name string) (*LoggingConfiguration, error)
	SetOwnershipControls(name string, config *OwnershipControls) error
	GetOwnershipControls(name string) (*OwnershipControls, error)
	SetRequestPayment(name string, config *RequestPaymentConfiguration) error
	GetRequestPayment(name string) (*RequestPaymentConfiguration, error)
	SetAccelerateConfiguration(name string, config *AccelerateConfiguration) error
	GetAccelerateConfiguration(name string) (*AccelerateConfiguration, error)
	SetACL(name string, acp *AccessControlPolicy) error
	GetPublicAccessBlock(name string) (*PublicAccessBlockConfig, error)
	SetReplication(name string, config *ReplicationConfiguration) error
	GetReplication(name string) (*ReplicationConfiguration, error)
	ARNBuilder() *svcarn.S3Builder
}

// ObjectStoreInterface defines operations for managing S3 objects.
type ObjectStoreInterface interface {
	Close()
	Get(ctx context.Context, bucket, key string) (io.ReadCloser, *Object, error)
	GetMetadata(bucket, key string) (*Object, error)
	Put(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string) (*Object, error)
	Delete(ctx context.Context, bucket, key string) error
	Exists(ctx context.Context, bucket, key string) (bool, error)
	Head(ctx context.Context, bucket, key string) (*Object, error)
	GetWithVersion(ctx context.Context, bucket, key, versionId string) (io.ReadCloser, *Object, error)
	PutWithVersioning(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error)
	DeleteWithVersion(ctx context.Context, bucket, key, versionId string) (*Object, error)
	HeadWithVersion(ctx context.Context, bucket, key, versionId string) (*Object, error)
	GetRangeWithVersion(ctx context.Context, bucket, key, versionId string, offset, length int64) (io.ReadCloser, *Object, error)
	SetACLWithVersion(bucket, key, versionId string, acp *AccessControlPolicy) error
	GetACLWithVersion(bucket, key, versionId string) (*AccessControlPolicy, error)
	CreateMultipartUpload(ctx context.Context, bucket, key, contentType string, metadata map[string]string, sseType SSEType, kmsKeyID, customerKeyMD5 string, sseMetadata *SSEObjectMetadata, plaintextDataKey []byte, storageClass ObjectStorageClass, acl *AccessControlPolicy) (*MultipartUpload, error)
	GetMultipartUpload(uploadId string) (*MultipartUpload, error)
	UploadPart(ctx context.Context, bucket, key, uploadId string, partNumber int, reader io.Reader, encryptedSize int64, plainSize int64, contentNonce, dataKey []byte) (*ObjectPart, error)
	ListParts(ctx context.Context, bucket, key, uploadId string, partNumberMarker, maxParts int) ([]ObjectPart, int, bool, error)
	CompleteMultipartUpload(ctx context.Context, bucket, key, uploadId string, parts []ObjectPart) (*Object, error)
	AbortMultipartUpload(ctx context.Context, bucket, key, uploadId string) error
	ListMultipartUploads(bucket, prefix, keyMarker, uploadIdMarker string, maxUploads int) (*MultipartUploadListResult, error)
	Copy(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string, storageClass ObjectStorageClass) (*Object, error)
	CopyWithMetadata(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string, contentType string, metadata map[string]string, storageClass ObjectStorageClass) (*Object, error)
	CopyWithVersion(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string, storageClass ObjectStorageClass) (*Object, error)
	CopyWithVersionAndMetadata(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string, contentType string, metadata map[string]string, storageClass ObjectStorageClass) (*Object, error)
	PutEncrypted(ctx context.Context, bucket, key string, encryptedData []byte, contentType string, metadata map[string]string, sseMetadata *SSEObjectMetadata, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error)
	PutEncryptedWithVersioning(ctx context.Context, bucket, key string, encryptedData []byte, contentType string, metadata map[string]string, sseMetadata *SSEObjectMetadata, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error)
	GetEncrypted(ctx context.Context, bucket, key, versionId string) ([]byte, *Object, error)
	UpdateObjectEncryption(ctx context.Context, bucket, key, versionId string, encryptedData []byte, sseMetadata *SSEObjectMetadata) (*Object, error)
	List(bucket, prefix, delimiter, marker string, maxKeys int) (*ObjectListResult, error)
	ListObjectVersions(bucket, prefix, delimiter, keyMarker, versionIdMarker string, maxKeys int) (*ObjectListResult, error)
	CountByBucket(bucket string) (int, error)
	CountMultipartUploadsByBucket(bucket string) (int, error)
	SetObjectLegalHold(ctx context.Context, bucket, key, versionId string, legalHold *ObjectLockLegalHold) error
	GetObjectLegalHold(ctx context.Context, bucket, key, versionId string) (*ObjectLockLegalHold, error)
	SetObjectRetention(ctx context.Context, bucket, key, versionId string, retention *ObjectLockRetention) error
	GetObjectRetention(ctx context.Context, bucket, key, versionId string) (*ObjectLockRetention, error)
	GetRange(ctx context.Context, bucket, key string, offset, length int64) (io.ReadCloser, *Object, error)
	SetTags(bucket, key, versionId string, tags []types.Tag) error
	SetACL(bucket, key string, acp *AccessControlPolicy) error
	GetACL(bucket, key string) (*AccessControlPolicy, error)
	SetStorageClass(bucket, key, versionId string, storageClass ObjectStorageClass) error
	SetRestoreState(bucket, key, versionId string, expiry *time.Time) error
	ActiveRestores() ([]RestoreIndexEntry, error)
	SetReplicationStatus(bucket, key, versionId, status string) error
}

// S3StoreInterface defines access to S3 stores.
type S3StoreInterface interface {
	Buckets(region string) BucketStoreInterface
	Objects(region string) ObjectStoreInterface
	// FindBucket searches all regions for a bucket with the given name.
	// Returns the bucket and the region in which it was found.
	FindBucket(name string) (*Bucket, string)
}

type regionStores struct {
	buckets *BucketStore
	objects *ObjectStore
}

type s3Store struct {
	storageManager *storage.RegionStorageManager
	blobStore      storage.BlobStore
	accountID      string
	stores         sync.Map
}

// NewS3Store creates a new S3 store.
func NewS3Store(storageManager *storage.RegionStorageManager, blobStore storage.BlobStore, accountID string) S3StoreInterface {
	return &s3Store{
		storageManager: storageManager,
		blobStore:      blobStore,
		accountID:      accountID,
	}
}

func (s *s3Store) getStoresForRegion(region string) (*regionStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*regionStores, error) {
		store, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}

		bucketStore := NewBucketStore(store, s.accountID, region)
		objectStore, err := NewObjectStore(store, s.blobStore, bucketStore, s.accountID, region)
		if err != nil {
			return nil, err
		}

		return &regionStores{
			buckets: bucketStore,
			objects: objectStore,
		}, nil
	})
}

// Buckets returns the bucket store for the given region.
// Returns nil when the underlying storage has not been initialised for
// the requested region (e.g. region is not configured).
func (s *s3Store) Buckets(region string) BucketStoreInterface {
	stores, err := s.getStoresForRegion(region)
	if err != nil {
		return nil
	}
	return stores.buckets
}

// Objects returns the object store for the given region.
// Returns nil when the underlying storage has not been initialised for
// the requested region (e.g. region is not configured).
func (s *s3Store) Objects(region string) ObjectStoreInterface {
	stores, err := s.getStoresForRegion(region)
	if err != nil {
		return nil
	}
	return stores.objects
}

// FindBucket searches all initialised region stores for a bucket with the
// given name. Returns the bucket metadata and the region in which it was
// found, or an error if the bucket does not exist in any region.
func (s *s3Store) FindBucket(name string) (*Bucket, string) {
	var found *Bucket
	var foundRegion string
	s.stores.Range(func(key, value interface{}) bool {
		region := key.(string)
		rs := value.(*regionStores)
		bucket, err := rs.buckets.Get(name)
		if err == nil && bucket != nil {
			found = bucket
			foundRegion = region
			return false
		}
		return true
	})
	if found != nil {
		return found, foundRegion
	}
	return nil, ""
}

var (
	_ BucketStoreInterface = (*BucketStore)(nil)
	_ ObjectStoreInterface = (*ObjectStore)(nil)
	_ S3StoreInterface     = (*s3Store)(nil)
)
