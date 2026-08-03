// Package timestreamwrite provides AWS Timestream Write service operations for vorpalstacks.
// This package is part of the Timestream service family and shares the underlying
// store with "vorpalstacks/internal/store/aws/timestream" and the sibling package
// "vorpalstacks/internal/services/aws/timestreamquery".
package timestreamwrite

import (
	"context"
	"fmt"
	"sync"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// Smithy model pagination limits (timestream-write-2018-11-01.json):
// PaginationLimit range {min:1, max:20} for ListDatabases/ListTables.
// PageLimit range {min:1, max:100} for ListBatchLoadTasks.
const (
	maxListDatabasesResults      = 20
	maxListTablesResults         = 20
	maxListBatchLoadTasksResults = 100
)

// tsWriteStores holds the various Timestream Write stores.
type tsWriteStores struct {
	store          *tsstore.Store
	tableStore     *tsstore.TableStore
	recordStore    *tsstore.RecordStore
	batchLoadStore *tsstore.BatchLoadTaskStore
}

// Close stops background goroutines in the RecordStore.
func (s *tsWriteStores) Close() {
	if s.recordStore != nil {
		s.recordStore.Close()
	}
}

// Service provides AWS Timestream Write operations.
type TimestreamWriteService struct {
	accountID      string
	serverHost     string
	dataPath       string
	stores         sync.Map // region → *tsWriteStores
	batchWg        sync.WaitGroup
	storageManager *storage.RegionStorageManager
	s3Invoker      eventbus.S3Invoker
	batchCtx       context.Context
	batchCancel    context.CancelFunc
}

// NewService creates a new Timestream Write service instance.
func NewTimestreamWriteService(accountID, serverHost, dataPath string) *TimestreamWriteService {
	ctx, cancel := context.WithCancel(context.Background())
	return &TimestreamWriteService{
		accountID:   accountID,
		serverHost:  serverHost,
		dataPath:    dataPath,
		batchCtx:    ctx,
		batchCancel: cancel,
	}
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *TimestreamWriteService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetS3Invoker injects the S3 reader for batch load tasks that load data from S3.
func (s *TimestreamWriteService) SetS3Invoker(invoker eventbus.S3Invoker) {
	s.s3Invoker = invoker
}

func (s *TimestreamWriteService) createStoreGroup(region string) (*tsWriteStores, error) {
	if s.storageManager == nil {
		return nil, fmt.Errorf("timestream write storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}
	tsStore := tsstore.NewStore(st, s.accountID, region)
	tableStore := tsstore.NewTableStore(st, tsStore, s.accountID, region)
	recordStore, err := tsstore.NewRecordStoreWithIndex(st, tableStore, region, s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create record store: %w", err)
	}
	return &tsWriteStores{
		store:          tsStore,
		tableStore:     tableStore,
		recordStore:    recordStore,
		batchLoadStore: tsstore.NewBatchLoadTaskStore(st, tableStore, region),
	}, nil
}

// GetDatabaseStoreForRegion returns the cached Store (database-level) for the given region,
// creating a new store group if not already cached.
func (s *TimestreamWriteService) GetDatabaseStoreForRegion(region string) (*tsstore.Store, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*tsWriteStores).store, nil
	}
	stores, err := s.createStoreGroup(region)
	if err != nil {
		return nil, err
	}
	actual, loaded := s.stores.LoadOrStore(region, stores)
	if loaded {
		stores.Close()
	}
	return actual.(*tsWriteStores).store, nil
}

// GetTableStoreForRegion returns the cached TableStore for the given region,
// creating a new store group if not already cached.
func (s *TimestreamWriteService) GetTableStoreForRegion(region string) (*tsstore.TableStore, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*tsWriteStores).tableStore, nil
	}
	stores, err := s.createStoreGroup(region)
	if err != nil {
		return nil, err
	}
	actual, loaded := s.stores.LoadOrStore(region, stores)
	if loaded {
		stores.Close()
	}
	return actual.(*tsWriteStores).tableStore, nil
}

func (s *TimestreamWriteService) store(ctx *request.RequestContext) (*tsWriteStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, ctx.GetRegion(), func() (*tsWriteStores, error) {
		storage, err := ctx.GetStorage()
		if err != nil {
			return nil, err
		}
		tsStore := tsstore.NewStore(storage, s.accountID, ctx.GetRegion())
		tableStore := tsstore.NewTableStore(storage, tsStore, s.accountID, ctx.GetRegion())
		recordStore, err := tsstore.NewRecordStoreWithIndex(storage, tableStore, ctx.GetRegion(), s.dataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create record store: %w", err)
		}
		return &tsWriteStores{
			store:          tsStore,
			tableStore:     tableStore,
			recordStore:    recordStore,
			batchLoadStore: tsstore.NewBatchLoadTaskStore(storage, tableStore, ctx.GetRegion()),
		}, nil
	})
}

// RegisterHandlers registers the Timestream Write service handlers with the dispatcher.
func (s *TimestreamWriteService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("timestream-write", "DescribeEndpoints", s.DescribeEndpoints)
	d.RegisterHandlerForService("timestream-write", "CreateDatabase", s.CreateDatabase)
	d.RegisterHandlerForService("timestream-write", "DescribeDatabase", s.DescribeDatabase)
	d.RegisterHandlerForService("timestream-write", "ListDatabases", s.ListDatabases)
	d.RegisterHandlerForService("timestream-write", "UpdateDatabase", s.UpdateDatabase)
	d.RegisterHandlerForService("timestream-write", "DeleteDatabase", s.DeleteDatabase)
	d.RegisterHandlerForService("timestream-write", "CreateTable", s.CreateTable)
	d.RegisterHandlerForService("timestream-write", "DescribeTable", s.DescribeTable)
	d.RegisterHandlerForService("timestream-write", "ListTables", s.ListTables)
	d.RegisterHandlerForService("timestream-write", "UpdateTable", s.UpdateTable)
	d.RegisterHandlerForService("timestream-write", "DeleteTable", s.DeleteTable)
	d.RegisterHandlerForService("timestream-write", "WriteRecords", s.WriteRecords)
	d.RegisterHandlerForService("timestream-write", "TagResource", s.TagResource)
	d.RegisterHandlerForService("timestream-write", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("timestream-write", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("timestream-write", "CreateBatchLoadTask", s.CreateBatchLoadTask)
	d.RegisterHandlerForService("timestream-write", "DescribeBatchLoadTask", s.DescribeBatchLoadTask)
	d.RegisterHandlerForService("timestream-write", "ListBatchLoadTasks", s.ListBatchLoadTasks)
	d.RegisterHandlerForService("timestream-write", "ResumeBatchLoadTask", s.ResumeBatchLoadTask)
}

func (s *TimestreamWriteService) mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if mapped != err {
		return mapped
	}
	return ErrInternalServer
}

// Close waits for any in-flight batch load simulation goroutines to
// finish and stops background goroutines in all cached stores.
func (s *TimestreamWriteService) Close() {
	if s.batchCancel != nil {
		s.batchCancel()
	}
	s.batchWg.Wait()
	s.stores.Range(func(_, v any) bool {
		if c, ok := v.(interface{ Close() }); ok {
			c.Close()
		}
		return true
	})
}
