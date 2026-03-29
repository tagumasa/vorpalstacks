// Package logs provides AWS CloudWatch Logs service operations for vorpalstacks.
package cloudwatchlogs

import (
	"context"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/services/aws/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
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

func (s *LogsService) store(reqCtx *request.RequestContext) (*logsstore.Store, error) {
	if store := reqCtx.GetCloudWatchLogsStore(); store != nil {
		return store.Raw(), nil
	}
	region := reqCtx.GetRegion()
	if cached, ok := s.logsStores.Load(region); ok {
		if typed, ok := cached.(*logsstore.Store); ok {
			return typed, nil
		}
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	store := logsstore.NewStore(storage.Bucket("logs-"+region), s.accountID, region, s.dataPath)
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

// RegisterHandlers registers the CloudWatch Logs service handlers with the dispatcher.
func (s *LogsService) RegisterHandlers(d *dispatcher.Dispatcher) {
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
