// Package athena provides AWS Athena service operations for vorpalstacks.
package athena

import (
	"context"
	"io"
	"sync"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// S3ObjectStore defines the interface for s3 object operations.
type S3ObjectStore interface {
	Get(ctx context.Context, bucket, key string) (io.ReadCloser, *s3store.Object, error)
	Put(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string) (*s3store.Object, error)
	List(bucket, prefix, delimiter, marker string, maxKeys int) (*s3store.ObjectListResult, error)
}

// athenaStores holds the various Athena stores.
type athenaStores struct {
	workGroupStore         *athenastore.WorkGroupStore
	namedQueryStore        *athenastore.NamedQueryStore
	preparedStatementStore *athenastore.PreparedStatementStore
	queryExecutionStore    *athenastore.QueryExecutionStore
	resultStore            *athenastore.ResultStore
	dataCatalogStore       *athenastore.DataCatalogStore
	databaseStore          *athenastore.DatabaseStore
	tableStore             *athenastore.TableStore
	tableDataStore         *athenastore.TableDataStore
}

// Service provides AWS Athena operations.
type Service struct {
	accountID     string
	serverHost    string
	s3ObjectStore S3ObjectStore
	asyncWg       sync.WaitGroup
	cancelMu      sync.Mutex
	cancelFuncs   map[string]context.CancelFunc
	stores        sync.Map // region → *athenaStores
}

// NewService creates a new Athena service instance.
func NewService(accountID, serverHost string) *Service {
	return &Service{
		accountID:   accountID,
		serverHost:  serverHost,
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

// NewServiceWithS3 creates a new Athena service with s3 object store.
func NewServiceWithS3(accountID, serverHost string, s3ObjectStore S3ObjectStore) *Service {
	return &Service{
		accountID:     accountID,
		serverHost:    serverHost,
		s3ObjectStore: s3ObjectStore,
		asyncWg:       sync.WaitGroup{},
		cancelFuncs:   make(map[string]context.CancelFunc),
	}
}

func (s *Service) setCancelFunc(id string, fn context.CancelFunc) {
	s.cancelMu.Lock()
	s.cancelFuncs[id] = fn
	s.cancelMu.Unlock()
}

func (s *Service) getAndRemoveCancelFunc(id string) (context.CancelFunc, bool) {
	s.cancelMu.Lock()
	fn, ok := s.cancelFuncs[id]
	delete(s.cancelFuncs, id)
	s.cancelMu.Unlock()
	return fn, ok
}

func (s *Service) store(reqCtx *request.RequestContext) (*athenaStores, error) {
	if athenaStore := reqCtx.GetAthenaStore(); athenaStore != nil {
		return &athenaStores{
			workGroupStore:         athenaStore.WorkGroupStoreRaw(),
			namedQueryStore:        athenaStore.NamedQueryStoreRaw(),
			preparedStatementStore: athenaStore.PreparedStatementStoreRaw(),
			queryExecutionStore:    athenaStore.QueryExecutionStoreRaw(),
			resultStore:            athenaStore.ResultStoreRaw(),
			dataCatalogStore:       athenaStore.DataCatalogStoreRaw(),
			databaseStore:          athenaStore.DatabaseStoreRaw(),
			tableStore:             athenaStore.TableStoreRaw(),
			tableDataStore:         athenaStore.TableDataStoreRaw(),
		}, nil
	}

	region := reqCtx.GetRegion()
	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(*athenaStores); ok {
			return typed, nil
		}
	}

	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}

	stores := &athenaStores{
		workGroupStore:         athenastore.NewWorkGroupStore(storage, s.accountID, region),
		namedQueryStore:        athenastore.NewNamedQueryStore(storage, region),
		preparedStatementStore: athenastore.NewPreparedStatementStore(storage, region),
		queryExecutionStore:    athenastore.NewQueryExecutionStore(storage, region),
		resultStore:            athenastore.NewResultStore(storage, region),
		dataCatalogStore:       athenastore.NewDataCatalogStore(storage, region),
		databaseStore:          athenastore.NewDatabaseStore(storage, region),
		tableStore:             athenastore.NewTableStore(storage, region),
		tableDataStore:         athenastore.NewTableDataStore(storage, region),
	}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		if typed, ok := actual.(*athenaStores); ok {
			return typed, nil
		}
	}
	return stores, nil
}

// Shutdown gracefully shuts down the Athena service by waiting for all asynchronous operations to complete.
func (s *Service) Shutdown() {
	s.asyncWg.Wait()
}

// RegisterHandlers registers the Athena service handlers with the dispatcher.
func (s *Service) RegisterHandlers(d *dispatcher.Dispatcher) {
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
}
