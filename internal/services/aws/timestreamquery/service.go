// Package timestreamquery provides AWS Timestream Query service operations for vorpalstacks.
// This package is part of the Timestream service family and shares the underlying
// store with "vorpalstacks/internal/store/aws/timestream" and the sibling package
// "vorpalstacks/internal/services/aws/timestreamwrite".
package timestreamquery

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/sqlparser"
)

// Smithy model pagination limits (timestream-query-2018-11-01.json):
// MaxQueryResults range {min:1, max:1000} for Query.MaxRows.
// MaxScheduledQueriesResults range {min:1, max:1000} for ListScheduledQueries.MaxResults.
// MaxTagsForResourceResult range {min:1, max:200} for ListTagsForResource.MaxResults.
const (
	maxQueryRows            = 1000
	maxListScheduledQueries = 1000
	maxListTagsForResource  = 200
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

// Close stops background goroutines in the RecordStore.
func (s *tsQueryStores) Close() {
	if s.recordStore != nil {
		s.recordStore.Close()
	}
}

// Service represents the Timestream Query service.
type TimestreamQueryService struct {
	accountID      string
	serverHost     string
	dataPath       string
	preprocessor   *sqlparser.Preprocessor
	stores         sync.Map // region → *tsQueryStores
	storageManager *storage.RegionStorageManager
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

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *TimestreamQueryService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// Close stops background goroutines in all cached stores.
func (s *TimestreamQueryService) Close() {
	s.stores.Range(func(_, v any) bool {
		if c, ok := v.(interface{ Close() }); ok {
			c.Close()
		}
		return true
	})
}

func (s *TimestreamQueryService) createStoreGroup(region string) (*tsQueryStores, error) {
	if s.storageManager == nil {
		return nil, fmt.Errorf("timestream query storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}
	dbStore := tsstore.NewStore(st, s.accountID, region)
	tableStore := tsstore.NewTableStore(st, dbStore, s.accountID, region)
	recordStore, err := tsstore.NewRecordStoreWithIndex(st, tableStore, region, s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create record store: %w", err)
	}
	return &tsQueryStores{
		recordStore:            recordStore,
		tableStore:             tableStore,
		dbStore:                dbStore,
		scheduledQueryStore:    tsstore.NewScheduledQueryStore(st, s.accountID, region),
		scheduledQueryRunStore: tsstore.NewScheduledQueryRunStore(st, region),
		accountSettingsStore:   tsstore.NewAccountSettingsStore(st, s.accountID, region),
		queryInfoStore:         storecommon.NewBaseStore(st.Bucket("timestream-query-info-"+region), "timestream-query"),
		arnBuilder:             svcarn.NewARNBuilder(s.accountID, region),
	}, nil
}

// GetScheduledQueryStoreForRegion returns the cached ScheduledQueryStore for the
// given region, creating a new store group if not already cached.
func (s *TimestreamQueryService) GetScheduledQueryStoreForRegion(region string) (*tsstore.ScheduledQueryStore, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*tsQueryStores).scheduledQueryStore, nil
	}
	stores, err := s.createStoreGroup(region)
	if err != nil {
		return nil, err
	}
	actual, loaded := s.stores.LoadOrStore(region, stores)
	if loaded {
		stores.Close()
	}
	return actual.(*tsQueryStores).scheduledQueryStore, nil
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
