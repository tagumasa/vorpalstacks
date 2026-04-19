// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/common"
	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
	storesqs "vorpalstacks/internal/store/aws/sqs"
)

type s3Stores struct {
	buckets s3store.BucketStoreInterface
	objects s3store.ObjectStoreInterface
}

// S3Service provides S3 service operations.
type S3Service struct {
	s3Store             s3store.S3StoreInterface
	blobStore           storage.BlobStore
	storageManager      *storage.RegionStorageManager
	accountID           string
	accessController    *AccessController
	credentialsProvider auth.CredentialsProvider
	encryptionManager   *EncryptionManager
	fallbackCache       sync.Map
	bus                 eventbus.Bus
	sqsStore            *storesqs.SQSStore
	lambdaInvoker       common.LambdaInvoker
	busUnsubscribe      func()
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
func NewS3ServiceWithStorage(blobStore storage.BlobStore, accountID string) *S3Service {
	return &S3Service{
		blobStore:         blobStore,
		accountID:         accountID,
		accessController:  NewAccessController(accountID),
		encryptionManager: NewEncryptionManager(),
	}
}

// SetCredentialsProvider sets the credentials provider for the S3 service.
func (s *S3Service) SetCredentialsProvider(provider auth.CredentialsProvider) {
	s.credentialsProvider = provider
}

// SetKMSClient sets the KMS client for server-side encryption.
func (s *S3Service) SetKMSClient(kmsClient KMSClient) {
	s.encryptionManager = NewEncryptionManagerWithKMS(kmsClient)
}

// SetStorageManager sets the storage manager for persisting SSE-S3 bucket keys.
func (s *S3Service) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// RestoreSSE3Keys restores persisted SSE-S3 bucket keys from storage.
func (s *S3Service) RestoreSSE3Keys() {
	if s.storageManager == nil {
		return
	}
	rs, err := s.storageManager.GetStorage(s.accountID)
	if err != nil {
		logs.Warn("failed to get storage for SSE-S3 key restore", logs.Err(err))
		return
	}
	bucket := rs.Bucket("s3_sse_keys")
	prefix := []byte("bucket:")
	iter := bucket.ScanPrefix(prefix)
	for iter.Next() {
		bucketName := string(iter.Key()[len(prefix):])
		data := iter.Value()
		if err := s.encryptionManager.LoadSSE3BucketKey(bucketName, data); err != nil {
			logs.Warn("failed to restore SSE-S3 bucket key", logs.String("bucket", bucketName), logs.Err(err))
		}
	}
	iter.Close()
}

// SetEventBus sets the event bus and registers the S3 notification handler.
// The handler is subscribed asynchronously so that object operations return
// immediately without waiting for notification delivery.
func (s *S3Service) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	if bus != nil {
		subID, err := eventbus.SubscribeTyped[*eventbus.S3ObjectEvent](bus, s.handleS3Notification, eventbus.WithAsync())
		if err == nil {
			s.busUnsubscribe = func() {
				_ = bus.Unsubscribe(subID)
			}
		}
	}
}

// SetSQSStore sets the SQS store used for dispatching S3 event notifications
// to SQS queue destinations.
func (s *S3Service) SetSQSStore(store *storesqs.SQSStore) {
	s.sqsStore = store
}

// SetLambdaInvoker sets the Lambda invoker used for dispatching S3 event
// notifications to Lambda function destinations.
func (s *S3Service) SetLambdaInvoker(invoker common.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// publishObjectNotification publishes an S3ObjectEvent to the event bus after
// a successful object operation. The event region is taken from the request
// context to support multi-region buckets.
func (s *S3Service) publishObjectNotification(ctx context.Context, reqCtx *request.RequestContext, bucket, key string, size int64, versionID, etag string, op eventbus.S3ObjectOp) {
	if s.bus == nil {
		return
	}

	region := reqCtx.GetRegion()
	_ = ctx

	evt := &eventbus.S3ObjectEvent{
		EventBase: eventbus.EventBase{
			Timestamp: time.Now().UTC(),
			Source:    "aws:s3",
			Region:    region,
			AccountID: s.accountID,
			Caller: eventbus.CallerContext{
				ServicePrincipal: "s3.amazonaws.com",
				AccountID:        s.accountID,
			},
		},
		Bucket:    bucket,
		Key:       key,
		Size:      size,
		VersionID: versionID,
		ETag:      etag,
		Op:        op,
	}

	if err := s.bus.Publish(context.Background(), evt); err != nil {
		logs.Warn("s3: event bus publish failed", logs.String("bucket", bucket), logs.String("key", key), logs.Err(err))
	}
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

// Close unsubscribes from the event bus, persists SSE-S3 bucket keys, and releases resources.
func (s *S3Service) Close() {
	if s.busUnsubscribe != nil {
		s.busUnsubscribe()
		s.busUnsubscribe = nil
	}
	s.persistSSE3Keys()
}

func (s *S3Service) persistSSE3Keys() {
	if s.storageManager == nil {
		return
	}
	rs, err := s.storageManager.GetStorage(s.accountID)
	if err != nil {
		logs.Warn("failed to get storage for SSE-S3 key persistence", logs.Err(err))
		return
	}
	bucket := rs.Bucket("s3_sse_keys")
	s.encryptionManager.ForEachSSE3BucketKey(func(bucketName string, keyData []byte) {
		if err := bucket.Put([]byte("bucket:"+bucketName), keyData); err != nil {
			logs.Warn("failed to persist SSE-S3 bucket key", logs.String("bucket", bucketName), logs.Err(err))
		}
	})
}
