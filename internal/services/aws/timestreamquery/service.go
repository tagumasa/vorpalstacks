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
type TimestreamQueryService struct {
	accountID    string
	serverHost   string
	dataPath     string
	preprocessor *sqlparser.Preprocessor
	stores       sync.Map // region → *tsQueryStores
}

// NewService creates a new Timestream Query service.
func NewTimestreamQueryService(accountID, serverHost, dataPath string) *TimestreamQueryService {
	return &TimestreamQueryService{
		accountID:    accountID,
		serverHost:   serverHost,
		dataPath:     dataPath,
		preprocessor: sqlparser.NewTimestreamPreprocessor(),
	}
}

func (s *TimestreamQueryService) store(ctx *request.RequestContext) (*tsQueryStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, ctx.GetRegion(), func() (*tsQueryStores, error) {
		storage, err := ctx.GetStorage()
		if err != nil {
			return nil, err
		}
		dbStore := tsstore.NewStore(storage, s.accountID, ctx.GetRegion())
		tableStore := tsstore.NewTableStore(storage, dbStore, s.accountID, ctx.GetRegion())
		recordStore, err := tsstore.NewRecordStoreWithIndex(storage, tableStore, ctx.GetRegion(), s.dataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create record store: %w", err)
		}
		return &tsQueryStores{
			recordStore:            recordStore,
			tableStore:             tableStore,
			dbStore:                dbStore,
			scheduledQueryStore:    tsstore.NewScheduledQueryStore(storage, s.accountID, ctx.GetRegion()),
			scheduledQueryRunStore: tsstore.NewScheduledQueryRunStore(storage, ctx.GetRegion()),
			accountSettingsStore:   tsstore.NewAccountSettingsStore(storage, s.accountID, ctx.GetRegion()),
			queryInfoStore:         storecommon.NewBaseStore(storage.Bucket("timestream-query-info-"+ctx.GetRegion()), "timestream-query"),
			arnBuilder:             svcarn.NewARNBuilder(s.accountID, ctx.GetRegion()),
		}, nil
	})
}

// RegisterHandlers registers the Timestream Query service handlers with the dispatcher.
func (s *TimestreamQueryService) RegisterHandlers(d handler.Registrar) {
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
