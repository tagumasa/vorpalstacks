// Package athena provides AWS Athena service operations for vorpalstacks.
package athena

import (
	"context"
	"io"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/aws/athena"
	storecommon "vorpalstacks/internal/store/aws/common"
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
	workGroupStore         *athena.WorkGroupStore
	namedQueryStore        *athena.NamedQueryStore
	preparedStatementStore *athena.PreparedStatementStore
	queryExecutionStore    *athena.QueryExecutionStore
	resultStore            *athena.ResultStore
	dataCatalogStore       *athena.DataCatalogStore
	databaseStore          *athena.DatabaseStore
	tableStore             *athena.TableStore
	tableDataStore         *athena.TableDataStore
}

// Service provides AWS Athena operations.
type Service struct {
	accountID     string
	s3ObjectStore S3ObjectStore
	asyncWg       sync.WaitGroup
	cancelMu      sync.Mutex
	cancelFuncs   map[string]context.CancelFunc
	stores        sync.Map // region → *athenaStores
}

// NewService creates a new Athena service instance.
func NewService(accountID string) *Service {
	return &Service{
		accountID:   accountID,
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

// NewServiceWithS3 creates a new Athena service with s3 object store.
func NewServiceWithS3(accountID string, s3ObjectStore S3ObjectStore) *Service {
	return &Service{
		accountID:     accountID,
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
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*athenaStores, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return &athenaStores{
			workGroupStore:         athena.NewWorkGroupStore(storage, s.accountID, reqCtx.GetRegion()),
			namedQueryStore:        athena.NewNamedQueryStore(storage, reqCtx.GetRegion()),
			preparedStatementStore: athena.NewPreparedStatementStore(storage, reqCtx.GetRegion()),
			queryExecutionStore:    athena.NewQueryExecutionStore(storage, reqCtx.GetRegion()),
			resultStore:            athena.NewResultStore(storage, reqCtx.GetRegion()),
			dataCatalogStore:       athena.NewDataCatalogStore(storage, reqCtx.GetRegion()),
			databaseStore:          athena.NewDatabaseStore(storage, reqCtx.GetRegion()),
			tableStore:             athena.NewTableStore(storage, reqCtx.GetRegion()),
			tableDataStore:         athena.NewTableDataStore(storage, reqCtx.GetRegion()),
		}, nil
	})
}

// Shutdown gracefully shuts down the Athena service by waiting for all asynchronous operations to complete.
func (s *Service) Shutdown() {
	s.asyncWg.Wait()
}

// RegisterHandlers registers the Athena service handlers with the dispatcher.
func (s *Service) RegisterHandlers(d handler.Registrar) {
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
