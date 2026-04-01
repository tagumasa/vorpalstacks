// Package timestreamwrite provides AWS Timestream Write service operations for vorpalstacks.
package timestreamwrite

import (
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// tsWriteStores holds the various Timestream Write stores.
type tsWriteStores struct {
	store          *tsstore.Store
	tableStore     *tsstore.TableStore
	recordStore    *tsstore.RecordStore
	batchLoadStore *tsstore.BatchLoadTaskStore
}

// Service provides AWS Timestream Write operations.
type Service struct {
	accountID  string
	serverHost string
	dataPath   string
}

// NewService creates a new Timestream Write service instance.
func NewService(accountID, serverHost, dataPath string) *Service {
	return &Service{
		accountID:  accountID,
		serverHost: serverHost,
		dataPath:   dataPath,
	}
}

func (s *Service) store(ctx *request.RequestContext) (*tsWriteStores, error) {
	if stores := ctx.GetTimestreamStores(); stores != nil {
		return &tsWriteStores{
			store:          stores.DatabaseStore().Raw(),
			tableStore:     stores.TableStore().Raw(),
			recordStore:    stores.RecordStore().Raw(),
			batchLoadStore: stores.BatchLoadTaskStore().Raw(),
		}, nil
	}
	storage, err := ctx.GetStorage()
	if err != nil {
		return nil, err
	}
	region := ctx.GetRegion()
	tsStore := tsstore.NewStore(storage, s.accountID, region)
	tableStore := tsstore.NewTableStore(storage, tsStore, s.accountID, region)
	recordStore, err := tsstore.NewRecordStoreWithIndex(storage, tableStore, region, s.dataPath)
	if err != nil {
		return nil, err
	}
	return &tsWriteStores{
		store:          tsStore,
		tableStore:     tableStore,
		recordStore:    recordStore,
		batchLoadStore: tsstore.NewBatchLoadTaskStore(storage, tableStore, region),
	}, nil
}

// RegisterHandlers registers the Timestream Write service handlers with the dispatcher.
func (s *Service) RegisterHandlers(d *dispatcher.Dispatcher) {
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
	d.RegisterHandlerForService("timestream-write", "DeleteBatchLoadTask", s.DeleteBatchLoadTask)
}

func (s *Service) mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case "database not found":
		return ErrResourceNotFound
	case "database already exists":
		return ErrResourceAlreadyExists
	case "database is not empty":
		return ErrValidationException
	case "table not found":
		return ErrResourceNotFound
	case "table already exists":
		return ErrResourceAlreadyExists
	case "batch load task not found":
		return ErrResourceNotFound
	case "batch load task already exists":
		return ErrResourceAlreadyExists
	default:
		return ErrInternalServer
	}
}
