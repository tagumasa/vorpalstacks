// Package logs provides AWS CloudWatch Logs service operations for vorpalstacks.
package cloudwatchlogs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// LogsService provides CloudWatch Logs service operations.
type LogsService struct {
	storageManager  *storage.RegionStorageManager
	accountID       string
	dataPath        string
	logsStores      sync.Map // region → *logsstore.Store
	cwMetricInvoker eventbus.CloudWatchMetricInvoker
	bus             eventbus.Bus
	kms             eventbus.KMSInvoker
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	queries         sync.Map // queryId → *queryState (per-service, not global)
}

// NewLogsService creates a new CloudWatch Logs service.
// Optional cross-service dependencies (Lambda invoker, Kinesis store) should be
// injected via setter methods before registering handlers.
func NewLogsService(storageMgr *storage.RegionStorageManager, accountID, dataPath string) *LogsService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &LogsService{
		storageManager: storageMgr,
		accountID:      accountID,
		dataPath:       dataPath,
		ctx:            ctx,
		cancel:         cancel,
	}
	s.startRetentionPurger()
	s.startScheduledQueryWorker()
	return s
}

// SetEventBus injects the event bus and registers handlers for CloudWatch Logs
// delivery, Lambda log writes, API Gateway access logs, and direct log event
// ingestion from EventBridge/Scheduler targets.
func (s *LogsService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	if bus != nil {
		s.kms = bus.KMSInvoker()
	}
	_, _ = eventbus.SubscribeTyped[*eventbus.CloudWatchLogDeliveryEvent](bus, s.handleBusDelivery, eventbus.WithAsync())
	_, _ = eventbus.SubscribeTyped[*eventbus.LambdaLogWriteEvent](bus, s.handleLambdaLogWrite, eventbus.WithAsync())
	_, _ = eventbus.SubscribeTyped[*eventbus.APIGatewayAccessLogEvent](bus, s.handleAPIGatewayAccessLog, eventbus.WithAsync())
	_, _ = eventbus.SubscribeTyped[*eventbus.CloudWatchLogsPutEvent](bus, s.handleDirectPutLogEvents, eventbus.WithAsync())
}

func (s *LogsService) store(reqCtx *request.RequestContext) (*logsstore.Store, error) {
	return storecommon.GetOrCreateStoreE(&s.logsStores, reqCtx.GetRegion(), func() (*logsstore.Store, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		store, err := logsstore.NewStore(storage, storage.Bucket("logs-"+reqCtx.GetRegion()), s.accountID, reqCtx.GetRegion(), s.dataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create CloudWatch Logs store: %w", err)
		}
		return store, nil
	})
}

// getLogsStoreByRegion resolves the CloudWatch Logs store for the given region.
// Used by bus handlers that operate outside of an HTTP request context.
func (s *LogsService) getLogsStoreByRegion(region string) (*logsstore.Store, error) {
	if cached, ok := s.logsStores.Load(region); ok {
		if typed, ok := cached.(*logsstore.Store); ok {
			return typed, nil
		}
	}
	regionStorage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %q: %w", region, err)
	}
	store, err := logsstore.NewStore(regionStorage, regionStorage.Bucket("logs-"+region), s.accountID, region, s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudWatch Logs store: %w", err)
	}
	if actual, loaded := s.logsStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(*logsstore.Store); ok {
			return typed, nil
		}
	}
	return store, nil
}

// GetStoreForRegion resolves the CloudWatch Logs store for the given region,
// creating it on first use. Cross-service consumers (the eventbus logs
// invoker) resolve stores through this method so that every writer and the
// API read plane share one store instance per region.
func (s *LogsService) GetStoreForRegion(region string) (*logsstore.Store, error) {
	return s.getLogsStoreByRegion(region)
}

// SetCloudWatchMetricInvoker injects the CloudWatch metric invoker for emitting
// metric data when metric filters match log events.
func (s *LogsService) SetCloudWatchMetricInvoker(invoker eventbus.CloudWatchMetricInvoker) {
	s.cwMetricInvoker = invoker
}

// RegisterHandlers registers the CloudWatch Logs service handlers with the dispatcher.
func (s *LogsService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("logs", "CreateLogGroup", s.CreateLogGroup)
	d.RegisterHandlerForService("logs", "DeleteLogGroup", s.DeleteLogGroup)
	d.RegisterHandlerForService("logs", "DescribeLogGroups", s.DescribeLogGroups)
	d.RegisterHandlerForService("logs", "ListLogGroups", s.ListLogGroups)
	d.RegisterHandlerForService("logs", "PutRetentionPolicy", s.PutRetentionPolicy)
	d.RegisterHandlerForService("logs", "DeleteRetentionPolicy", s.DeleteRetentionPolicy)
	d.RegisterHandlerForService("logs", "TagResource", s.TagResource)
	d.RegisterHandlerForService("logs", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("logs", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("logs", "TagLogGroup", s.TagLogGroup)
	d.RegisterHandlerForService("logs", "ListTagsLogGroup", s.ListTagsLogGroup)
	d.RegisterHandlerForService("logs", "UntagLogGroup", s.UntagLogGroup)

	d.RegisterHandlerForService("logs", "CreateLogStream", s.CreateLogStream)
	d.RegisterHandlerForService("logs", "DeleteLogStream", s.DeleteLogStream)
	d.RegisterHandlerForService("logs", "DescribeLogStreams", s.DescribeLogStreams)
	d.RegisterHandlerForService("logs", "ListLogStreams", s.ListLogStreams)

	d.RegisterHandlerForService("logs", "PutLogEvents", s.PutLogEvents)
	d.RegisterHandlerForService("logs", "GetLogEvents", s.GetLogEvents)
	d.RegisterHandlerForService("logs", "FilterLogEvents", s.FilterLogEvents)

	d.RegisterHandlerForService("logs", "PutMetricFilter", s.PutMetricFilter)
	d.RegisterHandlerForService("logs", "DeleteMetricFilter", s.DeleteMetricFilter)
	d.RegisterHandlerForService("logs", "DescribeMetricFilters", s.DescribeMetricFilters)
	d.RegisterHandlerForService("logs", "TestMetricFilter", s.TestMetricFilter)

	d.RegisterHandlerForService("logs", "PutSubscriptionFilter", s.PutSubscriptionFilter)
	d.RegisterHandlerForService("logs", "DeleteSubscriptionFilter", s.DeleteSubscriptionFilter)
	d.RegisterHandlerForService("logs", "DescribeSubscriptionFilters", s.DescribeSubscriptionFilters)

	d.RegisterHandlerForService("logs", "PutDestination", s.PutDestination)
	d.RegisterHandlerForService("logs", "PutDestinationPolicy", s.PutDestinationPolicy)
	d.RegisterHandlerForService("logs", "DeleteDestination", s.DeleteDestination)
	d.RegisterHandlerForService("logs", "DescribeDestinations", s.DescribeDestinations)

	d.RegisterHandlerForService("logs", "AssociateKmsKey", s.AssociateKmsKey)
	d.RegisterHandlerForService("logs", "DisassociateKmsKey", s.DisassociateKmsKey)
	d.RegisterHandlerForService("logs", "PutLogGroupDeletionProtection", s.PutLogGroupDeletionProtection)

	d.RegisterHandlerForService("logs", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("logs", "DeleteResourcePolicy", s.DeleteResourcePolicy)
	d.RegisterHandlerForService("logs", "DescribeResourcePolicies", s.DescribeResourcePolicies)

	d.RegisterHandlerForService("logs", "PutAccountPolicy", s.PutAccountPolicy)
	d.RegisterHandlerForService("logs", "DeleteAccountPolicy", s.DeleteAccountPolicy)
	d.RegisterHandlerForService("logs", "DescribeAccountPolicies", s.DescribeAccountPolicies)

	d.RegisterHandlerForService("logs", "PutDataProtectionPolicy", s.PutDataProtectionPolicy)
	d.RegisterHandlerForService("logs", "GetDataProtectionPolicy", s.GetDataProtectionPolicy)
	d.RegisterHandlerForService("logs", "DeleteDataProtectionPolicy", s.DeleteDataProtectionPolicy)

	d.RegisterHandlerForService("logs", "PutQueryDefinition", s.PutQueryDefinition)
	d.RegisterHandlerForService("logs", "DeleteQueryDefinition", s.DeleteQueryDefinition)
	d.RegisterHandlerForService("logs", "DescribeQueryDefinitions", s.DescribeQueryDefinitions)

	d.RegisterHandlerForService("logs", "GetLogGroupFields", s.GetLogGroupFields)
	d.RegisterHandlerForService("logs", "GetLogRecord", s.GetLogRecord)
	d.RegisterHandlerForService("logs", "GetLogObject", s.GetLogObject)
	d.RegisterHandlerForService("logs", "GetLogFields", s.GetLogFields)

	d.RegisterHandlerForService("logs", "CreateExportTask", s.CreateExportTask)
	d.RegisterHandlerForService("logs", "DescribeExportTasks", s.DescribeExportTasks)
	d.RegisterHandlerForService("logs", "CancelExportTask", s.CancelExportTask)

	d.RegisterHandlerForService("logs", "CreateImportTask", s.CreateImportTask)
	d.RegisterHandlerForService("logs", "DescribeImportTasks", s.DescribeImportTasks)
	d.RegisterHandlerForService("logs", "CancelImportTask", s.CancelImportTask)
	d.RegisterHandlerForService("logs", "DescribeImportTaskBatches", s.DescribeImportTaskBatches)

	d.RegisterHandlerForService("logs", "StartQuery", s.StartQuery)
	d.RegisterHandlerForService("logs", "StopQuery", s.StopQuery)
	d.RegisterHandlerForService("logs", "DescribeQueries", s.DescribeQueries)
	d.RegisterHandlerForService("logs", "GetQueryResults", s.GetQueryResults)

	d.RegisterHandlerForService("logs", "CreateScheduledQuery", s.CreateScheduledQuery)
	d.RegisterHandlerForService("logs", "DeleteScheduledQuery", s.DeleteScheduledQuery)
	d.RegisterHandlerForService("logs", "UpdateScheduledQuery", s.UpdateScheduledQuery)
	d.RegisterHandlerForService("logs", "GetScheduledQuery", s.GetScheduledQuery)
	d.RegisterHandlerForService("logs", "GetScheduledQueryHistory", s.GetScheduledQueryHistory)
	d.RegisterHandlerForService("logs", "ListScheduledQueries", s.ListScheduledQueries)

	d.RegisterHandlerForService("logs", "CreateLookupTable", s.CreateLookupTable)
	d.RegisterHandlerForService("logs", "DeleteLookupTable", s.DeleteLookupTable)
	d.RegisterHandlerForService("logs", "GetLookupTable", s.GetLookupTable)
	d.RegisterHandlerForService("logs", "UpdateLookupTable", s.UpdateLookupTable)
	d.RegisterHandlerForService("logs", "DescribeLookupTables", s.DescribeLookupTables)
}

func (s *LogsService) startRetentionPurger() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("PANIC in retention purger, restarting",
					logs.Any("panic", r))
				go s.startRetentionPurger()
			}
		}()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.purgeAllRegions()
			}
		}
	}()
}

func (s *LogsService) purgeAllRegions() {
	s.logsStores.Range(func(key, value interface{}) bool {
		store, ok := value.(*logsstore.Store)
		if !ok {
			return true
		}
		if err := store.PurgeAllExpiredChunks(); err != nil {
			logs.Error("Failed to purge expired chunks", logs.Err(err))
		}
		return true
	})
}

// AccountID returns the AWS account ID for this service.
func (s *LogsService) AccountID() string {
	return s.accountID
}

// Stop stops the CloudWatch Logs service by canceling the context and waiting for goroutines to complete.
func (s *LogsService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
