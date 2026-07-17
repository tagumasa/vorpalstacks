// Package athena provides AWS Athena service operations for vorpalstacks.
package athena

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/store/aws/athena"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// athenaStores holds the various Athena stores.
type athenaStores struct {
	workGroupStore           *athena.WorkGroupStore
	namedQueryStore          *athena.NamedQueryStore
	preparedStatementStore   *athena.PreparedStatementStore
	queryExecutionStore      *athena.QueryExecutionStore
	resultStore              *athena.ResultStore
	dataCatalogStore         *athena.DataCatalogStore
	databaseStore            *athena.DatabaseStore
	tableStore               *athena.TableStore
	tableDataStore           *athena.TableDataStore
	capacityReservationStore *athena.CapacityReservationStore
}

const cleanupInterval = 24 * time.Hour

// AthenaService provides AWS Athena operations.
type AthenaService struct {
	accountID      string
	s3Invoker      eventbus.S3Invoker
	region         string
	testMode       bool
	asyncWg        sync.WaitGroup
	cancelMu       sync.Mutex
	cancelFuncs    map[string]context.CancelFunc
	stores         sync.Map
	regionCleanups sync.Map
	storageManager *storage.RegionStorageManager
}

// NewAthenaService creates a new Athena service instance.
func NewAthenaService(accountID string) *AthenaService {
	return &AthenaService{
		accountID:   accountID,
		cancelFuncs: make(map[string]context.CancelFunc),
		testMode:    os.Getenv("TEST_MODE") == "true",
	}
}

// SetS3Invoker injects the S3 invoker for cross-service S3 operations.
func (s *AthenaService) SetS3Invoker(invoker eventbus.S3Invoker) {
	s.s3Invoker = invoker
}

// SetRegion sets the default region for S3 operations.
func (s *AthenaService) SetRegion(region string) {
	s.region = region
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *AthenaService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

func (s *AthenaService) setCancelFunc(id string, fn context.CancelFunc) {
	s.cancelMu.Lock()
	s.cancelFuncs[id] = fn
	s.cancelMu.Unlock()
}

func (s *AthenaService) getAndRemoveCancelFunc(id string) (context.CancelFunc, bool) {
	s.cancelMu.Lock()
	fn, ok := s.cancelFuncs[id]
	delete(s.cancelFuncs, id)
	s.cancelMu.Unlock()
	return fn, ok
}

// GetWorkGroupStoreForRegion returns the cached WorkGroupStore for the given
// region, creating a new store group if not already cached.
func (s *AthenaService) GetWorkGroupStoreForRegion(region string) (*athena.WorkGroupStore, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*athenaStores).workGroupStore, nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("athena storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}
	stores := &athenaStores{
		workGroupStore:           athena.NewWorkGroupStore(st, s.accountID, region),
		namedQueryStore:          athena.NewNamedQueryStore(st, region),
		preparedStatementStore:   athena.NewPreparedStatementStore(st, region),
		queryExecutionStore:      athena.NewQueryExecutionStore(st, region),
		resultStore:              athena.NewResultStore(st, region),
		dataCatalogStore:         athena.NewDataCatalogStore(st, s.accountID, region),
		databaseStore:            athena.NewDatabaseStore(st, region),
		tableStore:               athena.NewTableStore(st, region),
		tableDataStore:           athena.NewTableDataStore(st, region),
		capacityReservationStore: athena.NewCapacityReservationStore(st, s.accountID, region),
	}
	actual, _ := s.stores.LoadOrStore(region, stores)
	return actual.(*athenaStores).workGroupStore, nil
}

func (s *AthenaService) store(reqCtx *request.RequestContext) (*athenaStores, error) {
	st, err := storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*athenaStores, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return &athenaStores{
			workGroupStore:           athena.NewWorkGroupStore(storage, s.accountID, reqCtx.GetRegion()),
			namedQueryStore:          athena.NewNamedQueryStore(storage, reqCtx.GetRegion()),
			preparedStatementStore:   athena.NewPreparedStatementStore(storage, reqCtx.GetRegion()),
			queryExecutionStore:      athena.NewQueryExecutionStore(storage, reqCtx.GetRegion()),
			resultStore:              athena.NewResultStore(storage, reqCtx.GetRegion()),
			dataCatalogStore:         athena.NewDataCatalogStore(storage, s.accountID, reqCtx.GetRegion()),
			databaseStore:            athena.NewDatabaseStore(storage, reqCtx.GetRegion()),
			tableStore:               athena.NewTableStore(storage, reqCtx.GetRegion()),
			tableDataStore:           athena.NewTableDataStore(storage, reqCtx.GetRegion()),
			capacityReservationStore: athena.NewCapacityReservationStore(storage, s.accountID, reqCtx.GetRegion()),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	now := time.Now()
	lastVal, loaded := s.regionCleanups.LoadOrStore(reqCtx.GetRegion(), now)
	shouldCleanup := !loaded
	if loaded {
		interval := cleanupInterval
		if s.testMode {
			interval = 5 * time.Minute
		}
		if now.Sub(lastVal.(time.Time)) >= interval {
			shouldCleanup = true
		}
	}
	if shouldCleanup {
		s.regionCleanups.Store(reqCtx.GetRegion(), now)
		s.cleanupExpiredQueryExecutions(st)
	}
	return st, nil
}

// Shutdown gracefully shuts down the Athena service by waiting for all asynchronous operations to complete.
func (s *AthenaService) Shutdown() {
	s.asyncWg.Wait()
}

// RegisterHandlers registers the Athena service handlers with the dispatcher.
func (s *AthenaService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("athena", "CreateWorkGroup", s.CreateWorkGroup)
	d.RegisterHandlerForService("athena", "GetWorkGroup", s.GetWorkGroup)
	d.RegisterHandlerForService("athena", "UpdateWorkGroup", s.UpdateWorkGroup)
	d.RegisterHandlerForService("athena", "DeleteWorkGroup", s.DeleteWorkGroup)
	d.RegisterHandlerForService("athena", "ListWorkGroups", s.ListWorkGroups)
	d.RegisterHandlerForService("athena", "TagResource", s.TagResource)
	d.RegisterHandlerForService("athena", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("athena", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("athena", "StartQueryExecution", s.StartQueryExecution)
	d.RegisterHandlerForService("athena", "GetQueryExecution", s.GetQueryExecution)
	d.RegisterHandlerForService("athena", "StopQueryExecution", s.StopQueryExecution)
	d.RegisterHandlerForService("athena", "ListQueryExecutions", s.ListQueryExecutions)
	d.RegisterHandlerForService("athena", "BatchGetQueryExecution", s.BatchGetQueryExecution)
	d.RegisterHandlerForService("athena", "GetQueryResults", s.GetQueryResults)
	d.RegisterHandlerForService("athena", "GetQueryRuntimeStatistics", s.GetQueryRuntimeStatistics)
	d.RegisterHandlerForService("athena", "CreateNamedQuery", s.CreateNamedQuery)
	d.RegisterHandlerForService("athena", "GetNamedQuery", s.GetNamedQuery)
	d.RegisterHandlerForService("athena", "DeleteNamedQuery", s.DeleteNamedQuery)
	d.RegisterHandlerForService("athena", "ListNamedQueries", s.ListNamedQueries)
	d.RegisterHandlerForService("athena", "UpdateNamedQuery", s.UpdateNamedQuery)
	d.RegisterHandlerForService("athena", "BatchGetNamedQuery", s.BatchGetNamedQuery)
	d.RegisterHandlerForService("athena", "CreatePreparedStatement", s.CreatePreparedStatement)
	d.RegisterHandlerForService("athena", "GetPreparedStatement", s.GetPreparedStatement)
	d.RegisterHandlerForService("athena", "DeletePreparedStatement", s.DeletePreparedStatement)
	d.RegisterHandlerForService("athena", "ListPreparedStatements", s.ListPreparedStatements)
	d.RegisterHandlerForService("athena", "UpdatePreparedStatement", s.UpdatePreparedStatement)
	d.RegisterHandlerForService("athena", "BatchGetPreparedStatement", s.BatchGetPreparedStatement)
	d.RegisterHandlerForService("athena", "ListEngineVersions", s.ListEngineVersions)
	d.RegisterHandlerForService("athena", "ListDataCatalogs", s.ListDataCatalogs)
	d.RegisterHandlerForService("athena", "GetDataCatalog", s.GetDataCatalog)
	d.RegisterHandlerForService("athena", "CreateDataCatalog", s.CreateDataCatalog)
	d.RegisterHandlerForService("athena", "DeleteDataCatalog", s.DeleteDataCatalog)
	d.RegisterHandlerForService("athena", "UpdateDataCatalog", s.UpdateDataCatalog)
	d.RegisterHandlerForService("athena", "ListDatabases", s.ListDatabases)
	d.RegisterHandlerForService("athena", "GetDatabase", s.GetDatabase)
	d.RegisterHandlerForService("athena", "ListTableMetadata", s.ListTableMetadata)
	d.RegisterHandlerForService("athena", "GetTableMetadata", s.GetTableMetadata)
	d.RegisterHandlerForService("athena", "CreateCapacityReservation", s.CreateCapacityReservation)
	d.RegisterHandlerForService("athena", "GetCapacityReservation", s.GetCapacityReservation)
	d.RegisterHandlerForService("athena", "ListCapacityReservations", s.ListCapacityReservations)
	d.RegisterHandlerForService("athena", "UpdateCapacityReservation", s.UpdateCapacityReservation)
	d.RegisterHandlerForService("athena", "CancelCapacityReservation", s.CancelCapacityReservation)
	d.RegisterHandlerForService("athena", "DeleteCapacityReservation", s.DeleteCapacityReservation)
}
