// Package logs provides AWS CloudWatch Logs service operations for vorpalstacks.
package cloudwatchlogs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/common"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// LogsService provides CloudWatch Logs service operations.
type LogsService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	dataPath       string
	logsStores     sync.Map // region → *logsstore.Store
	metricStores   sync.Map // region → *cwstore.MetricChunkStore
	lambdaStores   sync.Map // region → common.LambdaInvoker
	kinesisStores  sync.Map // region → kinesisstore.KinesisStoreInterface
	bus            eventbus.Bus
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
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
	return s
}

// SetLambdaInvoker registers a Lambda invoker for a given region for subscription filter delivery.
func (s *LogsService) SetLambdaInvoker(region string, invoker common.LambdaInvoker) {
	if invoker != nil {
		s.lambdaStores.Store(region, invoker)
	}
}

// SetKinesisStore registers a Kinesis store for a given region for subscription filter delivery.
func (s *LogsService) SetKinesisStore(region string, store kinesisstore.KinesisStoreInterface) {
	if store != nil {
		s.kinesisStores.Store(region, store)
	}
}

// SetEventBus injects the event bus and registers handlers for CloudWatch Logs
// delivery, Lambda log writes, API Gateway access logs, and direct log event
// ingestion from EventBridge/Scheduler targets.
func (s *LogsService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
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

func (s *LogsService) getLambdaInvoker(region string) common.LambdaInvoker {
	if cached, ok := s.lambdaStores.Load(region); ok {
		if typed, ok := cached.(common.LambdaInvoker); ok {
			return typed
		}
	}
	return nil
}

func (s *LogsService) getKinesisStore(region string) (kinesisstore.KinesisStoreInterface, bool) {
	if cached, ok := s.kinesisStores.Load(region); ok {
		if typed, ok := cached.(kinesisstore.KinesisStoreInterface); ok {
			return typed, true
		}
	}
	regionStorage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, false
	}
	tstore, ok := regionStorage.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, false
	}
	store := kinesisstore.NewKinesisStore(tstore, s.accountID, region)
	if actual, loaded := s.kinesisStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(*kinesisstore.KinesisStore); ok {
			return typed, true
		}
	}
	return store, true
}

func (s *LogsService) getMetricStore(reqCtx *request.RequestContext) (*cwstore.MetricChunkStore, error) {
	region := reqCtx.GetRegion()
	if cached, ok := s.metricStores.Load(region); ok {
		if typed, ok := cached.(*cwstore.MetricChunkStore); ok {
			return typed, nil
		}
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	store, err := cwstore.NewMetricChunkStoreWithIndex(storage, region, s.dataPath)
	if err != nil {
		return nil, err
	}
	if actual, loaded := s.metricStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(*cwstore.MetricChunkStore); ok {
			return typed, nil
		}
	}
	return store, nil
}

// getMetricStoreByRegion resolves the CloudWatch metric store for the given region.
// Used by bus handlers that operate outside of an HTTP request context.
func (s *LogsService) getMetricStoreByRegion(region string) (*cwstore.MetricChunkStore, error) {
	if cached, ok := s.metricStores.Load(region); ok {
		if typed, ok := cached.(*cwstore.MetricChunkStore); ok {
			return typed, nil
		}
	}
	regionStorage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %q: %w", region, err)
	}
	store, err := cwstore.NewMetricChunkStoreWithIndex(regionStorage, region, s.dataPath)
	if err != nil {
		return nil, err
	}
	if actual, loaded := s.metricStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(*cwstore.MetricChunkStore); ok {
			return typed, nil
		}
	}
	return store, nil
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
}

func (s *LogsService) startRetentionPurger() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
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
		store := value.(*logsstore.Store)
		if err := store.PurgeAllExpiredChunks(); err != nil {
			logs.Error("Failed to purge expired chunks", logs.Err(err))
		}
		return true
	})
}

// Stop stops the CloudWatch Logs service by canceling the context and waiting for goroutines to complete.
func (s *LogsService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
