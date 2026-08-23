// Package kinesis provides AWS Kinesis stream service operations for vorpalstacks.
package kinesis

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// KinesisService provides AWS Kinesis stream operations.
type KinesisService struct {
	accountID      string
	storageManager *storage.RegionStorageManager
	stores         sync.Map // region → *kinesisstore.KinesisStore
}

// NewKinesisService creates a new Kinesis service instance.
func NewKinesisService(accountID, region string) *KinesisService {
	return &KinesisService{
		accountID: accountID,
	}
}

// SetStorageManager injects the region storage manager for lazy store
// creation in GetStoreForRegion.
func (s *KinesisService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetKinesisStore pre-populates the store cache for the given region so that
// the same store instance is used by both the service handlers and
// cross-service integrations (Lambda ESM, CloudWatch Logs).
func (s *KinesisService) SetKinesisStore(region string, store *kinesisstore.KinesisStore) {
	s.stores.Store(region, store)
}

func (s *KinesisService) store(reqCtx *request.RequestContext) (*kinesisstore.KinesisStore, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*kinesisstore.KinesisStore, error) {
		basicStore, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		tstore, ok := basicStore.(storage.TransactionalStorageWith2PC)
		if !ok {
			return nil, fmt.Errorf("storage does not support TransactionalStorageWith2PC")
		}

		return kinesisstore.NewKinesisStore(tstore, reqCtx.GetAccountID(), reqCtx.GetRegion()), nil
	})
}

// GetStoreForRegion returns the KinesisStore for the given region, creating
// it from the region storage on first use. The service handlers, the admin
// console, and cross-service integrations (the eventbus kinesis invoker)
// resolve stores through this method so that all callers share one store
// instance per region.
func (s *KinesisService) GetStoreForRegion(region string) (*kinesisstore.KinesisStore, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*kinesisstore.KinesisStore), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("kinesis store not initialised for region %s", region)
	}
	regionStorage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	tstore, ok := regionStorage.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("kinesis: storage for region %s does not support 2PC", region)
	}
	store := kinesisstore.NewKinesisStore(tstore, s.accountID, region)
	if actual, loaded := s.stores.LoadOrStore(region, store); loaded {
		return actual.(*kinesisstore.KinesisStore), nil
	}
	return store, nil
}

// RegisterHandlers registers the Kinesis service handlers with the dispatcher.
func (s *KinesisService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("kinesis", "CreateStream", s.CreateStream)
	d.RegisterHandlerForService("kinesis", "DeleteStream", s.DeleteStream)
	d.RegisterHandlerForService("kinesis", "DescribeStream", s.DescribeStream)
	d.RegisterHandlerForService("kinesis", "DescribeStreamSummary", s.DescribeStreamSummary)
	d.RegisterHandlerForService("kinesis", "ListStreams", s.ListStreams)
	d.RegisterHandlerForService("kinesis", "UpdateStreamMode", s.UpdateStreamMode)

	d.RegisterHandlerForService("kinesis", "PutRecord", s.PutRecord)
	d.RegisterHandlerForService("kinesis", "PutRecords", s.PutRecords)
	d.RegisterHandlerForService("kinesis", "GetRecords", s.GetRecords)
	d.RegisterHandlerForService("kinesis", "GetShardIterator", s.GetShardIterator)

	d.RegisterHandlerForService("kinesis", "ListShards", s.ListShards)
	d.RegisterHandlerForService("kinesis", "SplitShard", s.SplitShard)
	d.RegisterHandlerForService("kinesis", "MergeShards", s.MergeShards)
	d.RegisterHandlerForService("kinesis", "UpdateShardCount", s.UpdateShardCount)

	d.RegisterHandlerForService("kinesis", "RegisterStreamConsumer", s.RegisterStreamConsumer)
	d.RegisterHandlerForService("kinesis", "DeregisterStreamConsumer", s.DeregisterStreamConsumer)
	d.RegisterHandlerForService("kinesis", "DescribeStreamConsumer", s.DescribeStreamConsumer)
	d.RegisterHandlerForService("kinesis", "ListStreamConsumers", s.ListStreamConsumers)
	d.RegisterHandlerForService("kinesis", "SubscribeToShard", s.SubscribeToShard)

	d.RegisterHandlerForService("kinesis", "AddTagsToStream", s.AddTagsToStream)
	d.RegisterHandlerForService("kinesis", "RemoveTagsFromStream", s.RemoveTagsFromStream)
	d.RegisterHandlerForService("kinesis", "ListTagsForStream", s.ListTagsForStream)
	d.RegisterHandlerForService("kinesis", "TagResource", s.TagResource)
	d.RegisterHandlerForService("kinesis", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("kinesis", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("kinesis", "IncreaseStreamRetentionPeriod", s.IncreaseStreamRetentionPeriod)
	d.RegisterHandlerForService("kinesis", "DecreaseStreamRetentionPeriod", s.DecreaseStreamRetentionPeriod)
	d.RegisterHandlerForService("kinesis", "DescribeLimits", s.DescribeLimits)
	d.RegisterHandlerForService("kinesis", "DescribeAccountSettings", s.DescribeAccountSettings)
	d.RegisterHandlerForService("kinesis", "UpdateAccountSettings", s.UpdateAccountSettings)

	d.RegisterHandlerForService("kinesis", "EnableEnhancedMonitoring", s.EnableEnhancedMonitoring)
	d.RegisterHandlerForService("kinesis", "DisableEnhancedMonitoring", s.DisableEnhancedMonitoring)

	d.RegisterHandlerForService("kinesis", "StartStreamEncryption", s.StartStreamEncryption)
	d.RegisterHandlerForService("kinesis", "StopStreamEncryption", s.StopStreamEncryption)

	d.RegisterHandlerForService("kinesis", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("kinesis", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("kinesis", "DeleteResourcePolicy", s.DeleteResourcePolicy)

	d.RegisterHandlerForService("kinesis", "UpdateMaxRecordSize", s.UpdateMaxRecordSize)
	d.RegisterHandlerForService("kinesis", "UpdateStreamWarmThroughput", s.UpdateStreamWarmThroughput)
}
