// Package timestreamquery provides AWS Timestream Query service operations for vorpalstacks.
package timestreamquery

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
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
	queryInfoStore         *storecommon.BaseStore
	arnBuilder             *svcarn.ARNBuilder
}

// Service represents the Timestream Query service.
type Service struct {
	accountID    string
	serverHost   string
	dataPath     string
	preprocessor *sqlparser.Preprocessor
	stores       sync.Map // region → *tsQueryStores
}

// NewService creates a new Timestream Query service.
func NewService(accountID, serverHost, dataPath string) *Service {
	return &Service{
		accountID:    accountID,
		serverHost:   serverHost,
		dataPath:     dataPath,
		preprocessor: sqlparser.NewTimestreamPreprocessor(),
	}
}

func (s *Service) store(ctx *request.RequestContext) (*tsQueryStores, error) {
	if stores := ctx.GetTimestreamStores(); stores != nil {
		storage, err := ctx.GetStorage()
		if err != nil {
			return nil, err
		}
		region := ctx.GetRegion()
		return &tsQueryStores{
			recordStore:            stores.RecordStore().Raw(),
			tableStore:             stores.TableStore().Raw(),
			dbStore:                stores.DatabaseStore().Raw(),
			scheduledQueryStore:    stores.ScheduledQueryStore().Raw(),
			scheduledQueryRunStore: stores.ScheduledQueryRunStore().Raw(),
			accountSettingsStore:   stores.AccountSettingsStore().Raw(),
			queryInfoStore:         storecommon.NewBaseStore(storage.Bucket("timestream-query-info-"+region), "timestream-query"),
			arnBuilder:             svcarn.NewARNBuilder(s.accountID, ctx.GetRegion()),
		}, nil
	}

	region := ctx.GetRegion()
	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(*tsQueryStores); ok {
			return typed, nil
		}
	}

	storage, err := ctx.GetStorage()
	if err != nil {
		return nil, err
	}

	dbStore := tsstore.NewStore(storage, s.accountID, region)
	tableStore := tsstore.NewTableStore(storage, dbStore, s.accountID, region)
	recordStore, err := tsstore.NewRecordStoreWithIndex(storage, tableStore, region, s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create record store: %w", err)
	}

	stores := &tsQueryStores{
		recordStore:            recordStore,
		tableStore:             tableStore,
		dbStore:                dbStore,
		scheduledQueryStore:    tsstore.NewScheduledQueryStore(storage, s.accountID, region),
		scheduledQueryRunStore: tsstore.NewScheduledQueryRunStore(storage, region),
		accountSettingsStore:   tsstore.NewAccountSettingsStore(storage, s.accountID, region),
		queryInfoStore:         storecommon.NewBaseStore(storage.Bucket("timestream-query-info-"+region), "timestream-query"),
		arnBuilder:             svcarn.NewARNBuilder(s.accountID, region),
	}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		if typed, ok := actual.(*tsQueryStores); ok {
			return typed, nil
		}
	}
	return stores, nil
}

// RegisterHandlers registers the Timestream Query service handlers with the dispatcher.
func (s *Service) RegisterHandlers(d handler.Registrar) {
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
	d.RegisterHandlerForService("timestream-query", "TagResource", s.TagResource)
	d.RegisterHandlerForService("timestream-query", "UntagResource", s.UntagResource)
}
