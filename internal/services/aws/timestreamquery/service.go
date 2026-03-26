// Package timestreamquery provides AWS Timestream Query service operations for vorpalstacks.
package timestreamquery

import (
	"sync"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	tsstore "vorpalstacks/internal/store/aws/timestream"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/sqlparser"
)

type tsQueryStores struct {
	recordStore            *tsstore.RecordStore
	tableStore             *tsstore.TableStore
	dbStore                *tsstore.Store
	scheduledQueryStore    *tsstore.ScheduledQueryStore
	scheduledQueryRunStore *tsstore.ScheduledQueryRunStore
	accountSettingsStore   *tsstore.AccountSettingsStore
	arnBuilder             *svcarn.ARNBuilder
}

// Service represents the Timestream Query service.
type Service struct {
	accountID      string
	serverHost     string
	dataPath       string
	runningQueries map[string]*QueryInfo
	queryMutex     sync.RWMutex
	preprocessor   *sqlparser.Preprocessor
}

// NewService creates a new Timestream Query service.
func NewService(accountID, serverHost, dataPath string) *Service {
	return &Service{
		accountID:      accountID,
		serverHost:     serverHost,
		dataPath:       dataPath,
		runningQueries: make(map[string]*QueryInfo),
		preprocessor:   sqlparser.NewTimestreamPreprocessor(),
	}
}

func (s *Service) store(ctx *request.RequestContext) (*tsQueryStores, error) {
	if stores := ctx.GetTimestreamStores(); stores != nil {
		return &tsQueryStores{
			recordStore:            stores.RecordStore().Raw(),
			tableStore:             stores.TableStore().Raw(),
			dbStore:                stores.DatabaseStore().Raw(),
			scheduledQueryStore:    stores.ScheduledQueryStore().Raw(),
			scheduledQueryRunStore: stores.ScheduledQueryRunStore().Raw(),
			accountSettingsStore:   stores.AccountSettingsStore().Raw(),
			arnBuilder:             svcarn.NewARNBuilder(s.accountID, ctx.GetRegion()),
		}, nil
	}
	storage, err := ctx.GetStorage()
	if err != nil {
		return nil, err
	}
	region := ctx.GetRegion()
	dbStore := tsstore.NewStore(storage, s.accountID, region)
	tableStore := tsstore.NewTableStore(storage, dbStore, s.accountID, region)
	recordStore, err := tsstore.NewRecordStoreWithIndex(storage, tableStore, region, s.dataPath)
	if err != nil {
		return nil, err
	}
	return &tsQueryStores{
		recordStore:            recordStore,
		tableStore:             tableStore,
		dbStore:                dbStore,
		scheduledQueryStore:    tsstore.NewScheduledQueryStore(storage, s.accountID, region),
		scheduledQueryRunStore: tsstore.NewScheduledQueryRunStore(storage, region),
		accountSettingsStore:   tsstore.NewAccountSettingsStore(storage, s.accountID, region),
		arnBuilder:             svcarn.NewARNBuilder(s.accountID, region),
	}, nil
}

// RegisterHandlers registers the Timestream Query service handlers with the dispatcher.
func (s *Service) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("timestream-query", "DescribeEndpoints", s.DescribeEndpoints)
	d.RegisterHandlerForService("timestream-query", "Query", s.Query)
	d.RegisterHandlerForService("timestream-query", "CancelQuery", s.CancelQuery)
	d.RegisterHandlerForService("timestream-query", "PrepareQuery", s.PrepareQuery)
	d.RegisterHandlerForService("timestream-query", "CreateScheduledQuery", s.CreateScheduledQuery)
	d.RegisterHandlerForService("timestream-query", "DeleteScheduledQuery", s.DeleteScheduledQuery)
	d.RegisterHandlerForService("timestream-query", "DescribeScheduledQuery", s.DescribeScheduledQuery)
	d.RegisterHandlerForService("timestream-query", "ListScheduledQueries", s.ListScheduledQueries)
	d.RegisterHandlerForService("timestream-query", "UpdateScheduledQuery", s.UpdateScheduledQuery)
	d.RegisterHandlerForService("timestream-query", "ExecuteScheduledQuery", s.ExecuteScheduledQuery)
	d.RegisterHandlerForService("timestream-query", "UpdateAccountSettings", s.UpdateAccountSettings)
	d.RegisterHandlerForService("timestream-query", "DescribeAccountSettings", s.DescribeAccountSettings)
	d.RegisterHandlerForService("timestream-query", "ListTagsForResource", s.ListTagsForResource)
}
