package s3

// Package s3 provides S3 (Simple Storage Service) data store implementations
// for vorpalstacks.

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_s3"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

func bucketBucketName(region string) string {
	return "s3_buckets-" + region
}

var (
	// ErrBucketNotFound is returned when the specified bucket does not exist.
	ErrBucketNotFound = common.NewStoreError("s3", "bucket_not_found", common.ErrNotFound)
	// ErrBucketAlreadyExists is returned when attempting to create a bucket that already exists.
	ErrBucketAlreadyExists = common.NewStoreError("s3", "bucket_already_exists", common.ErrAlreadyExists)
	// ErrBucketNotEmpty is returned when attempting to delete a bucket that contains objects.
	ErrBucketNotEmpty = common.NewStoreError("s3", "bucket_not_empty", common.ErrConflict)
)

// BucketStore manages S3 bucket metadata storage and retrieval.
type BucketStore struct {
	*common.BaseStore
	arnBuilder         *svcarn.S3Builder
	createMu           sync.Mutex
	callbackMu         sync.RWMutex
	keyLocker          common.KeyLocker
	region             string
	onVersioningChange func(bucket string, enabled bool)
	onDelete           func(bucket string)
}

// NewBucketStore creates a new BucketStore instance.
// It initialises the store with the provided storage, account ID, and region.
func NewBucketStore(store storage.BasicStorage, accountId, region string) *BucketStore {
	return &BucketStore{
		BaseStore:  common.NewBaseStore(store.Bucket(bucketBucketName(region)), "s3"),
		arnBuilder: svcarn.NewARNBuilder(accountId, region).S3(),
		region:     region,
	}
}

// Get retrieves a bucket by name.
func (s *BucketStore) Get(name string) (*Bucket, error) {
	var bucket pb.Bucket
	if err := s.BaseStore.GetProto(name, &bucket); err != nil {
		return nil, err
	}
	return ProtoToBucket(&bucket), nil
}

// Create creates a new bucket with the specified name and region.
// Returns the created bucket's metadata or an error if it already exists.
func (s *BucketStore) Create(name, region string) (*Bucket, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	if s.Exists(name) {
		return nil, ErrBucketAlreadyExists
	}

	bucket := &Bucket{
		Name:         name,
		Region:       region,
		CreationDate: time.Now().UTC(),
		ACL:          nil,
	}

	if err := s.BaseStore.PutProto(name, BucketToProto(bucket)); err != nil {
		return nil, err
	}
	return bucket, nil
}

// Put updates an existing bucket's metadata in the store.
func (s *BucketStore) Put(bucket *Bucket) error {
	return s.BaseStore.PutProto(bucket.Name, BucketToProto(bucket))
}

// Delete removes a bucket from the store.
// Returns an error if the bucket does not exist.
func (s *BucketStore) Delete(name string) error {
	s.callbackMu.RLock()
	callback := s.onDelete
	s.callbackMu.RUnlock()
	if callback != nil {
		callback(name)
	}
	s.keyLocker.Delete(name)
	return s.BaseStore.Delete(name)
}

// Exists checks whether a bucket exists in the store.
func (s *BucketStore) Exists(name string) bool {
	return s.BaseStore.Exists(name)
}

// atomicUpdate performs a locked Get-Modify-Put cycle on a bucket.
// The modify function receives the bucket and may return an error to abort the update.
func (s *BucketStore) atomicUpdate(name string, modify func(*Bucket) error) error {
	return s.keyLocker.WithLock(name, func() error {
		bucket, err := s.Get(name)
		if err != nil {
			return err
		}
		if err := modify(bucket); err != nil {
			return err
		}
		return s.Put(bucket)
	})
}

// List returns a list of all buckets in the store.
func (s *BucketStore) List() ([]*Bucket, error) {
	result, err := common.ListProto[*pb.Bucket](s.BaseStore, common.ListOptions{MaxItems: 10000}, func() *pb.Bucket { return &pb.Bucket{} }, nil)
	if err != nil {
		return nil, err
	}
	buckets := make([]*Bucket, len(result.Items))
	for i, b := range result.Items {
		buckets[i] = ProtoToBucket(b)
	}
	return buckets, nil
}

// SetVersioning sets the versioning status for a bucket.
func (s *BucketStore) SetVersioning(name string, status BucketVersioningStatus) error {
	return s.keyLocker.WithLock(name, func() error {
		bucket, err := s.Get(name)
		if err != nil {
			return err
		}
		bucket.VersioningStatus = status
		if err := s.Put(bucket); err != nil {
			return err
		}
		s.callbackMu.RLock()
		callback := s.onVersioningChange
		s.callbackMu.RUnlock()
		if callback != nil {
			callback(name, status == BucketVersioningEnabled)
		}
		return nil
	})
}

// SetVersioningCallback sets a callback function that is invoked when
// the versioning status of a bucket changes.
func (s *BucketStore) SetVersioningCallback(fn func(bucket string, enabled bool)) {
	s.callbackMu.Lock()
	defer s.callbackMu.Unlock()
	s.onVersioningChange = fn
}

// SetOnDeleteCallback registers a callback function to be invoked when a bucket is deleted.
// The callback receives the name of the deleted bucket.
func (s *BucketStore) SetOnDeleteCallback(fn func(bucket string)) {
	s.callbackMu.Lock()
	defer s.callbackMu.Unlock()
	s.onDelete = fn
}

// SetEncryption sets the encryption configuration for a bucket.
func (s *BucketStore) SetEncryption(name string, config *EncryptionConfig) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.EncryptionConfig = config; return nil })
}

// GetEncryptionConfiguration retrieves the encryption configuration for a bucket.
func (s *BucketStore) GetEncryptionConfiguration(name string) (*EncryptionConfig, error) {
	bucket, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return bucket.EncryptionConfig, nil
}

// SetPolicy sets the bucket policy for a bucket.
func (s *BucketStore) SetPolicy(name, policy string) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.Policy = policy; return nil })
}

// SetCORS sets the CORS configuration for a bucket.
func (s *BucketStore) SetCORS(name string, config *CORSConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.CORSConfiguration = config; return nil })
}

// SetPublicAccessBlock sets the public access block configuration for a bucket.
func (s *BucketStore) SetPublicAccessBlock(name string, config *PublicAccessBlockConfig) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.PublicAccessBlock = config; return nil })
}

// SetTags sets the tags for a bucket.
func (s *BucketStore) SetTags(name string, tags []types.Tag) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.Tags = tags; return nil })
}

// SetLifecycleConfiguration sets the lifecycle configuration for a bucket.
func (s *BucketStore) SetLifecycleConfiguration(name string, config *LifecycleConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.LifecycleConfiguration = config; return nil })
}

// ARNBuilder returns the S3 ARN builder for this bucket store.
func (s *BucketStore) ARNBuilder() *svcarn.S3Builder {
	return s.arnBuilder
}

// SetACL sets the access control list for a bucket.
func (s *BucketStore) SetACL(name string, acp *AccessControlPolicy) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.ACL = acp; return nil })
}

// GetPublicAccessBlock retrieves the public access block configuration for a bucket.
func (s *BucketStore) GetPublicAccessBlock(name string) (*PublicAccessBlockConfig, error) {
	bucket, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return bucket.PublicAccessBlock, nil
}

// SetWebsiteConfiguration sets the website configuration for a bucket.
func (s *BucketStore) SetWebsiteConfiguration(name string, config *WebsiteConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.WebsiteConfiguration = config; return nil })
}

// SetObjectLockConfiguration sets the object lock configuration for a bucket.
func (s *BucketStore) SetObjectLockConfiguration(name string, config *ObjectLockConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.ObjectLockConfig = config; return nil })
}

// SetNotificationConfiguration sets the notification configuration for a bucket.
func (s *BucketStore) SetNotificationConfiguration(name string, config *NotificationConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.NotificationConfiguration = config; return nil })
}

// GetNotificationConfiguration retrieves the notification configuration for a bucket.
func (s *BucketStore) GetNotificationConfiguration(name string) (*NotificationConfiguration, error) {
	bucket, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return bucket.NotificationConfiguration, nil
}

// SetLoggingConfiguration sets the logging configuration for a bucket.
func (s *BucketStore) SetLoggingConfiguration(name string, config *LoggingConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.LoggingConfiguration = config; return nil })
}

// GetLoggingConfiguration retrieves the logging configuration for a bucket.
func (s *BucketStore) GetLoggingConfiguration(name string) (*LoggingConfiguration, error) {
	bucket, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return bucket.LoggingConfiguration, nil
}

// SetOwnershipControls sets the ownership controls for a bucket.
func (s *BucketStore) SetOwnershipControls(name string, config *OwnershipControls) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.OwnershipControls = config; return nil })
}

// GetOwnershipControls retrieves the ownership controls for a bucket.
func (s *BucketStore) GetOwnershipControls(name string) (*OwnershipControls, error) {
	bucket, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return bucket.OwnershipControls, nil
}

// SetRequestPayment sets the request payment configuration for a bucket.
func (s *BucketStore) SetRequestPayment(name string, config *RequestPaymentConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.RequestPayment = config; return nil })
}

// GetRequestPayment retrieves the request payment configuration for a bucket.
func (s *BucketStore) GetRequestPayment(name string) (*RequestPaymentConfiguration, error) {
	bucket, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return bucket.RequestPayment, nil
}

// SetAccelerateConfiguration sets the transfer acceleration configuration for a bucket.
func (s *BucketStore) SetAccelerateConfiguration(name string, config *AccelerateConfiguration) error {
	return s.atomicUpdate(name, func(b *Bucket) error { b.AccelerateConfiguration = config; return nil })
}

// GetAccelerateConfiguration retrieves the transfer acceleration configuration for a bucket.
func (s *BucketStore) GetAccelerateConfiguration(name string) (*AccelerateConfiguration, error) {
	bucket, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return bucket.AccelerateConfiguration, nil
}
