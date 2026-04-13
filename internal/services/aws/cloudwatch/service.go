// Package cloudwatch provides CloudWatch service operations for vorpalstacks.
package cloudwatch

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/server/eventbus"
	"vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/services/aws/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// cloudwatchStores holds the stores for CloudWatch resources.
type cloudwatchStores struct {
	metrics    *cwstore.MetricChunkStore
	alarms     *cwstore.AlarmStore
	dashboards *cwstore.DashboardStore
}

// CloudWatchService provides CloudWatch operations.
type CloudWatchService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	region         string
	dataPath       string
	bus            eventbus.Bus
	lambdaInvoker  common.LambdaInvoker
	evaluator      *alarmEvaluator
	logger         logs.Logger
	stores         sync.Map // region → *cloudwatchStores
}

// NewCloudWatchService creates a new CloudWatch service.
//
// Parameters:
//   - accountID: The AWS account ID
//   - region: The AWS region
//   - dataPath: The data path for chunk storage
//
// Returns:
//   - *CloudWatchService: A new CloudWatch service
func NewCloudWatchService(accountID, region, dataPath string) *CloudWatchService {
	return &CloudWatchService{
		accountID: accountID,
		region:    region,
		dataPath:  dataPath,
	}
}

// SetEventBus sets the event bus used for publishing alarm state change
// events. The alarm evaluator is started when both the bus and a logger
// are available.
func (s *CloudWatchService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// SetLambdaInvoker sets the Lambda invoker used for dispatching alarm
// actions to Lambda function targets.
func (s *CloudWatchService) SetLambdaInvoker(invoker common.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// SetStorageManager injects the region storage manager for multi-region
// alarm evaluation.
func (s *CloudWatchService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetLogger sets the structured logger used by the alarm evaluator for
// diagnostic output.
func (s *CloudWatchService) SetLogger(logger logs.Logger) {
	s.logger = logger
}

// StartEvaluator creates and starts the background alarm evaluation loop.
// This should be called once during server initialisation after SetEventBus
// and SetLambdaInvoker have been wired.
func (s *CloudWatchService) StartEvaluator(ctx context.Context) {
	s.evaluator = newAlarmEvaluator(0, 0, s.logger)
	s.evaluator.Start(ctx, s)
}

// StopEvaluator gracefully shuts down the alarm evaluation loop, waiting
// for any in-flight evaluations to complete. This should be called during
// server shutdown.
func (s *CloudWatchService) StopEvaluator() {
	if s.evaluator != nil {
		s.evaluator.Stop()
	}
}

// store returns the CloudWatch stores for a given request context.
func (s *CloudWatchService) store(reqCtx *request.RequestContext) (*cloudwatchStores, error) {
	if reqCtx.GetCloudWatchMetricStore() != nil {
		if ms, ok := reqCtx.GetCloudWatchMetricStore().(*cwstore.MetricChunkStore); ok {
			return &cloudwatchStores{metrics: ms}, nil
		}
	}

	region := reqCtx.GetRegion()
	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(*cloudwatchStores); ok {
			return typed, nil
		}
	}

	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}

	metricStore, err := cwstore.NewMetricChunkStoreWithIndex(storage, region, s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric store: %w", err)
	}
	alarmStore := cwstore.NewAlarmStore(storage, s.accountID, region)
	dashboardStore := cwstore.NewDashboardStore(storage, s.accountID, region)

	stores := &cloudwatchStores{
		metrics:    metricStore,
		alarms:     alarmStore,
		dashboards: dashboardStore,
	}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		if typed, ok := actual.(*cloudwatchStores); ok {
			return typed, nil
		}
	}
	return stores, nil
}

// RegisterHandlers registers CloudWatch handlers with the dispatcher.
func (s *CloudWatchService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("monitoring", "PutMetricData", s.PutMetricData)
	d.RegisterHandlerForService("monitoring", "GetMetricStatistics", s.GetMetricStatistics)
	d.RegisterHandlerForService("monitoring", "ListMetrics", s.ListMetrics)
	d.RegisterHandlerForService("monitoring", "GetMetricData", s.GetMetricData)
	d.RegisterHandlerForService("monitoring", "PutMetricAlarm", s.PutMetricAlarm)
	d.RegisterHandlerForService("monitoring", "DescribeAlarms", s.DescribeAlarms)
	d.RegisterHandlerForService("monitoring", "DescribeAlarmsForMetric", s.DescribeAlarmsForMetric)
	d.RegisterHandlerForService("monitoring", "DeleteAlarms", s.DeleteAlarms)
	d.RegisterHandlerForService("monitoring", "SetAlarmState", s.SetAlarmState)
	d.RegisterHandlerForService("monitoring", "GetMetricWidgetImage", s.GetMetricWidgetImage)
	d.RegisterHandlerForService("monitoring", "PutDashboard", s.PutDashboard)
	d.RegisterHandlerForService("monitoring", "GetDashboard", s.GetDashboard)
	d.RegisterHandlerForService("monitoring", "ListDashboards", s.ListDashboards)
	d.RegisterHandlerForService("monitoring", "DeleteDashboards", s.DeleteDashboards)
	d.RegisterHandlerForService("monitoring", "TagResource", s.TagResource)
	d.RegisterHandlerForService("monitoring", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("monitoring", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("monitoring", "EnableAlarmActions", s.EnableAlarmActions)
	d.RegisterHandlerForService("monitoring", "DisableAlarmActions", s.DisableAlarmActions)
	d.RegisterHandlerForService("monitoring", "DescribeAlarmHistory", s.DescribeAlarmHistory)
	d.RegisterHandlerForService("monitoring", "PutCompositeAlarm", s.PutCompositeAlarm)
}
