// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/crypto"
)

type s3Stores struct {
	buckets s3store.BucketStoreInterface
	objects s3store.ObjectStoreInterface
}

// S3Service provides S3 service operations.
type S3Service struct {
	s3Store             s3store.S3StoreInterface
	storage             storage.BasicStorage
	blobStore           storage.BlobStore
	accountID           string
	accessController    *AccessController
	credentialsProvider crypto.CredentialsProvider
	encryptionManager   *EncryptionManager
	fallbackCache       sync.Map
}

// NewS3Service creates a new S3Service with the given store and blob store.
func NewS3Service(s3Store s3store.S3StoreInterface, blobStore storage.BlobStore, accountID string) *S3Service {
	return &S3Service{
		s3Store:           s3Store,
		blobStore:         blobStore,
		accountID:         accountID,
		accessController:  NewAccessController(accountID),
		encryptionManager: NewEncryptionManager(),
	}
}

// NewS3ServiceWithStorage creates a new S3Service with the given storage.
func NewS3ServiceWithStorage(store storage.BasicStorage, blobStore storage.BlobStore, accountID, region string) *S3Service {
	return &S3Service{
		storage:           store,
		blobStore:         blobStore,
		accountID:         accountID,
		accessController:  NewAccessController(accountID),
		encryptionManager: NewEncryptionManager(),
	}
}

// SetCredentialsProvider sets the credentials provider for the S3 service.
func (s *S3Service) SetCredentialsProvider(provider crypto.CredentialsProvider) {
	s.credentialsProvider = provider
}

// SetKMSClient sets the KMS client for server-side encryption.
func (s *S3Service) SetKMSClient(kmsClient KMSClient) {
	s.encryptionManager = NewEncryptionManagerWithKMS(kmsClient)
}

func (s *S3Service) store(ctx *request.RequestContext) (*s3Stores, error) {
	region := ctx.GetRegion()

	if s.s3Store != nil {
		buckets := s.s3Store.Buckets(region)
		objects := s.s3Store.Objects(region)
		if buckets == nil || objects == nil {
			return nil, fmt.Errorf("failed to get stores for region %s", region)
		}
		return &s3Stores{buckets: buckets, objects: objects}, nil
	}

	if v, ok := s.fallbackCache.Load(region); ok {
		return v.(*s3Stores), nil
	}

	store, err := ctx.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %w", err)
	}

	bucketStore := s3store.NewBucketStore(store, s.accountID, region)
	objectStore, err := s3store.NewObjectStore(store, s.blobStore, bucketStore, s.accountID, region)
	if err != nil {
		return nil, fmt.Errorf("failed to create object store: %w", err)
	}

	stores := &s3Stores{
		buckets: bucketStore,
		objects: objectStore,
	}

	s.fallbackCache.Store(region, stores)
	return stores, nil
}

// AccountId returns the account ID for the S3 service.
func (s *S3Service) AccountId() string {
	return s.accountID
}
