// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DynamoDBService provides DynamoDB operations for managing tables, items, and other resources.
type DynamoDBService struct {
	accountID               string
	stores                  sync.Map // region → dynamodbstore.DynamoDBStoreInterface
	storageManager          *storage.RegionStorageManager
	busStoreFactory         *dynamodbstore.DynamoDBStoreFactory
	bus                     eventbus.Bus
	bgCtx                   context.Context
	bgCancel                context.CancelFunc
	bgWg                    sync.WaitGroup
	idempotencySweepOnce    sync.Once
	clientRequestTokenLocks [clientRequestTokenLockShards]sync.Mutex
}

// NewDynamoDBService creates a new DynamoDB service instance.
func NewDynamoDBService(accountID string) *DynamoDBService {
	ctx, cancel := context.WithCancel(context.Background())
	return &DynamoDBService{
		accountID: accountID,
		bgCtx:     ctx,
		bgCancel:  cancel,
	}
}

// Close stops background goroutines (TTL workers, state-transition
// goroutines) in all cached stores and waits for them to finish.
func (s *DynamoDBService) Close() {
	if s.bgCancel != nil {
		s.bgCancel()
	}
	s.bgWg.Wait()
	s.stores.Range(func(_, v any) bool {
		if c, ok := v.(interface{ Close() }); ok {
			c.Close()
		}
		return true
	})
}

// SetStorageManager sets the storage manager for cross-service store access
// (e.g. from the EventBus DynamoDBInvoker).
func (s *DynamoDBService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
	s.busStoreFactory = dynamodbstore.NewDynamoDBStoreFactory(sm, s.accountID)
}

// SetEventBus sets the EventBus for cross-service invoker access (e.g. S3
// import/export operations).
func (s *DynamoDBService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// GetStoreForRegion returns the DynamoDB store for the given region.
// This is used by the DynamoDBInvoker adapter for cross-service access.
func (s *DynamoDBService) GetStoreForRegion(region string) (dynamodbstore.DynamoDBStoreInterface, error) {
	if s.busStoreFactory == nil {
		return nil, fmt.Errorf("storage manager not configured for DynamoDB service")
	}
	return s.busStoreFactory.GetStore(region)
}

// GetCachedStoreForRegion returns the cached DynamoDB store for the given
// region, creating a new store instance if not already cached. This ensures
// the admin console handlers work without requiring a prior HTTP API request
// to initialise the store.
func (s *DynamoDBService) GetCachedStoreForRegion(region string) (dynamodbstore.DynamoDBStoreInterface, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(dynamodbstore.DynamoDBStoreInterface), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("dynamodb storage manager not initialised")
	}
	basicStorage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}
	txnStorage, ok := basicStorage.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("storage does not implement TransactionalStorageWith2PC")
	}
	store := dynamodbstore.NewDynamoDBStore(txnStorage, s.accountID, region)
	actual, loaded := s.stores.LoadOrStore(region, store)
	if loaded {
		store.Close()
	}
	return actual.(dynamodbstore.DynamoDBStoreInterface), nil
}

func (s *DynamoDBService) store(reqCtx *request.RequestContext) (dynamodbstore.DynamoDBStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (dynamodbstore.DynamoDBStoreInterface, error) {
		basicStorage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		txnStorage, ok := basicStorage.(storage.TransactionalStorageWith2PC)
		if !ok {
			return nil, fmt.Errorf("storage does not implement TransactionalStorageWith2PC")
		}
		return dynamodbstore.NewDynamoDBStore(txnStorage, s.accountID, reqCtx.GetRegion()), nil
	})
}

// RegisterHandlers registers all DynamoDB operation handlers with the dispatcher.
func (s *DynamoDBService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("dynamodb", "CreateTable", s.CreateTable)
	d.RegisterHandlerForService("dynamodb", "DeleteTable", s.DeleteTable)
	d.RegisterHandlerForService("dynamodb", "DescribeTable", s.DescribeTable)
	d.RegisterHandlerForService("dynamodb", "ListTables", s.ListTables)
	d.RegisterHandlerForService("dynamodb", "UpdateTable", s.UpdateTable)
	d.RegisterHandlerForService("dynamodb", "PutItem", s.PutItem)
	d.RegisterHandlerForService("dynamodb", "GetItem", s.GetItem)
	d.RegisterHandlerForService("dynamodb", "DeleteItem", s.DeleteItem)
	d.RegisterHandlerForService("dynamodb", "UpdateItem", s.UpdateItem)
	d.RegisterHandlerForService("dynamodb", "Query", s.Query)
	d.RegisterHandlerForService("dynamodb", "Scan", s.Scan)
	d.RegisterHandlerForService("dynamodb", "BatchGetItem", s.BatchGetItem)
	d.RegisterHandlerForService("dynamodb", "BatchWriteItem", s.BatchWriteItem)
	d.RegisterHandlerForService("dynamodb", "TagResource", s.TagResource)
	d.RegisterHandlerForService("dynamodb", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("dynamodb", "ListTagsOfResource", s.ListTagsForResource)

	d.RegisterHandlerForService("dynamodb", "TransactGetItems", s.TransactGetItems)
	d.RegisterHandlerForService("dynamodb", "TransactWriteItems", s.TransactWriteItems)

	d.RegisterHandlerForService("dynamodb", "DescribeTimeToLive", s.DescribeTimeToLive)
	d.RegisterHandlerForService("dynamodb", "UpdateTimeToLive", s.UpdateTimeToLive)

	d.RegisterHandlerForService("dynamodb", "DescribeEndpoints", s.DescribeEndpoints)

	// DynamoDB Streams operations (registered under "dynamodb" because the
	// signing service for DynamoDB Streams is "dynamodb").
	d.RegisterHandlerForService("dynamodb", "DescribeStream", s.DescribeStream)
	d.RegisterHandlerForService("dynamodb", "GetShardIterator", s.GetShardIterator)
	d.RegisterHandlerForService("dynamodb", "GetRecords", s.GetRecords)
	d.RegisterHandlerForService("dynamodb", "ListStreams", s.ListStreams)
	d.RegisterHandlerForService("dynamodb", "DescribeLimits", s.DescribeLimits)

	d.RegisterHandlerForService("dynamodb", "CreateBackup", s.CreateBackup)
	d.RegisterHandlerForService("dynamodb", "DeleteBackup", s.DeleteBackup)
	d.RegisterHandlerForService("dynamodb", "DescribeBackup", s.DescribeBackup)
	d.RegisterHandlerForService("dynamodb", "ListBackups", s.ListBackups)
	d.RegisterHandlerForService("dynamodb", "RestoreTableFromBackup", s.RestoreTableFromBackup)
	d.RegisterHandlerForService("dynamodb", "RestoreTableToPointInTime", s.RestoreTableToPointInTime)

	d.RegisterHandlerForService("dynamodb", "CreateGlobalTable", s.CreateGlobalTable)
	d.RegisterHandlerForService("dynamodb", "DescribeGlobalTable", s.DescribeGlobalTable)
	d.RegisterHandlerForService("dynamodb", "DescribeGlobalTableSettings", s.DescribeGlobalTableSettings)
	d.RegisterHandlerForService("dynamodb", "ListGlobalTables", s.ListGlobalTables)
	d.RegisterHandlerForService("dynamodb", "UpdateGlobalTable", s.UpdateGlobalTable)
	d.RegisterHandlerForService("dynamodb", "UpdateGlobalTableSettings", s.UpdateGlobalTableSettings)

	d.RegisterHandlerForService("dynamodb", "DeleteResourcePolicy", s.DeleteResourcePolicy)
	d.RegisterHandlerForService("dynamodb", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("dynamodb", "PutResourcePolicy", s.PutResourcePolicy)

	d.RegisterHandlerForService("dynamodb", "DescribeContinuousBackups", s.DescribeContinuousBackups)
	d.RegisterHandlerForService("dynamodb", "UpdateContinuousBackups", s.UpdateContinuousBackups)

	d.RegisterHandlerForService("dynamodb", "DescribeKinesisStreamingDestination", s.DescribeKinesisStreamingDestination)
	d.RegisterHandlerForService("dynamodb", "EnableKinesisStreamingDestination", s.EnableKinesisStreamingDestination)
	d.RegisterHandlerForService("dynamodb", "DisableKinesisStreamingDestination", s.DisableKinesisStreamingDestination)
	d.RegisterHandlerForService("dynamodb", "UpdateKinesisStreamingDestination", s.UpdateKinesisStreamingDestination)

	d.RegisterHandlerForService("dynamodb", "ExportTableToPointInTime", s.ExportTableToPointInTime)
	d.RegisterHandlerForService("dynamodb", "DescribeExport", s.DescribeExport)
	d.RegisterHandlerForService("dynamodb", "ListExports", s.ListExports)
	d.RegisterHandlerForService("dynamodb", "ImportTable", s.ImportTable)
	d.RegisterHandlerForService("dynamodb", "DescribeImport", s.DescribeImport)
	d.RegisterHandlerForService("dynamodb", "ListImports", s.ListImports)

	d.RegisterHandlerForService("dynamodb", "DescribeContributorInsights", s.DescribeContributorInsights)
	d.RegisterHandlerForService("dynamodb", "ListContributorInsights", s.ListContributorInsights)
	d.RegisterHandlerForService("dynamodb", "UpdateContributorInsights", s.UpdateContributorInsights)

	d.RegisterHandlerForService("dynamodb", "DescribeTableReplicaAutoScaling", s.DescribeTableReplicaAutoScaling)
	d.RegisterHandlerForService("dynamodb", "UpdateTableReplicaAutoScaling", s.UpdateTableReplicaAutoScaling)

	d.RegisterHandlerForService("dynamodb", "ExecuteStatement", s.ExecuteStatement)
	d.RegisterHandlerForService("dynamodb", "ExecuteTransaction", s.ExecuteTransaction)
	d.RegisterHandlerForService("dynamodb", "BatchExecuteStatement", s.BatchExecuteStatement)
}
