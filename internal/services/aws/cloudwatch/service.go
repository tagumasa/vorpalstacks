// Package cloudwatch provides CloudWatch service operations for vorpalstacks.
package cloudwatch

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// cloudwatchStores holds the stores for CloudWatch resources.
type cloudwatchStores struct {
	metrics          *cwstore.MetricChunkStore
	alarms           *cwstore.AlarmStore
	dashboards       *cwstore.DashboardStore
	anomalyDetectors *cwstore.AnomalyDetectorStore
	insightRules     *cwstore.InsightRuleStore
	alarmMuteRules   *cwstore.AlarmMuteRuleStore
}

// CloudWatchService provides CloudWatch operations.
type CloudWatchService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	region         string
	dataPath       string
	bus            eventbus.ServiceBus
	evaluator      *alarmEvaluator
	logger         logs.Logger
	stores         sync.Map // region → *cloudwatchStores

	// Global service state (not per-region). Protected by globalMu.
	globalMu   sync.Mutex
	otelStatus string            // OTel enrichment status: "STARTED" or "STOPPED"
	datasetKMS map[string]string // datasetIdentifier → kmsKeyArn
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
		accountID:  accountID,
		region:     region,
		dataPath:   dataPath,
		otelStatus: "STOPPED",
		datasetKMS: make(map[string]string),
	}
}

// SetEventBus sets the event bus used for publishing alarm state change
// events and subscribes the S3 request-metrics ingest. The alarm evaluator
// is started when both the bus and a logger are available.
func (s *CloudWatchService) SetEventBus(bus eventbus.ServiceBus) {
	s.bus = bus
	if bus != nil {
		if _, err := eventbus.SubscribeTyped[*eventbus.S3RequestMetricsEvent](bus, s.handleS3RequestMetrics, eventbus.WithAsync()); err != nil {
			logs.Warn("cloudwatch: subscribe s3 request metrics failed", logs.Err(err))
		}
	}
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
// has been wired.
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
	stores := &cloudwatchStores{
		metrics:          metricStore,
		alarms:           cwstore.NewAlarmStore(storage, s.accountID, region),
		dashboards:       cwstore.NewDashboardStore(storage, s.accountID, region),
		anomalyDetectors: cwstore.NewAnomalyDetectorStore(storage, s.accountID, region),
		insightRules:     cwstore.NewInsightRuleStore(storage, region),
		alarmMuteRules:   cwstore.NewAlarmMuteRuleStore(storage, s.accountID, region),
	}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		metricStore.Close()
		if typed, ok := actual.(*cloudwatchStores); ok {
			return typed, nil
		}
	}
	return stores, nil
}

// GetStoreForRegion returns the cached CloudWatch stores for the given region,
// creating a new store instance if not already cached.
func (s *CloudWatchService) GetStoreForRegion(region string) (*cloudwatchStores, error) {
	if v, ok := s.stores.Load(region); ok {
		if typed, ok := v.(*cloudwatchStores); ok {
			return typed, nil
		}
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("cloudwatch storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	metricStore, err := cwstore.NewMetricChunkStoreWithIndex(st, region, s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric store: %w", err)
	}
	stores := &cloudwatchStores{
		metrics:          metricStore,
		alarms:           cwstore.NewAlarmStore(st, s.accountID, region),
		dashboards:       cwstore.NewDashboardStore(st, s.accountID, region),
		anomalyDetectors: cwstore.NewAnomalyDetectorStore(st, s.accountID, region),
		insightRules:     cwstore.NewInsightRuleStore(st, region),
		alarmMuteRules:   cwstore.NewAlarmMuteRuleStore(st, s.accountID, region),
	}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		metricStore.Close()
		if typed, ok := actual.(*cloudwatchStores); ok {
			return typed, nil
		}
	}
	return stores, nil
}

// MetricStoreForRegion returns the metric store for the given region,
// creating the store group if not already cached. Cross-service consumers
// (the eventbus CloudWatch metric invoker) resolve stores through this
// method so that every writer and the API read plane share one
// MetricChunkStore per region — the store tracks chunk files in per-
// instance memory, so a second instance would orphan chunks.
func (s *CloudWatchService) MetricStoreForRegion(region string) (*cwstore.MetricChunkStore, error) {
	stores, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	return stores.metrics, nil
}

// AlarmStoreForRegion returns the alarm store for the given region,
// creating the store group if not already cached. Cross-service consumers
// (the eventbus CloudWatch alarm invoker) resolve stores through this
// method for the same single-instance-per-region guarantee.
func (s *CloudWatchService) AlarmStoreForRegion(region string) (*cwstore.AlarmStore, error) {
	stores, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	return stores.alarms, nil
}

// RegisterHandlers registers CloudWatch handlers with the dispatcher.
func (s *CloudWatchService) RegisterHandlers(d handler.Registrar) {
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
	d.RegisterHandlerForService("monitoring", "PutAnomalyDetector", s.PutAnomalyDetector)
	d.RegisterHandlerForService("monitoring", "DeleteAnomalyDetector", s.DeleteAnomalyDetector)
	d.RegisterHandlerForService("monitoring", "DescribeAnomalyDetectors", s.DescribeAnomalyDetectors)
	d.RegisterHandlerForService("monitoring", "PutInsightRule", s.PutInsightRule)
	d.RegisterHandlerForService("monitoring", "DeleteInsightRules", s.DeleteInsightRules)
	d.RegisterHandlerForService("monitoring", "DescribeInsightRules", s.DescribeInsightRules)
	d.RegisterHandlerForService("monitoring", "EnableInsightRules", s.EnableInsightRules)
	d.RegisterHandlerForService("monitoring", "DisableInsightRules", s.DisableInsightRules)
	d.RegisterHandlerForService("monitoring", "GetInsightRuleReport", s.GetInsightRuleReport)
	d.RegisterHandlerForService("monitoring", "PutManagedInsightRules", s.PutManagedInsightRules)
	d.RegisterHandlerForService("monitoring", "ListManagedInsightRules", s.ListManagedInsightRules)
	d.RegisterHandlerForService("monitoring", "PutAlarmMuteRule", s.PutAlarmMuteRule)
	d.RegisterHandlerForService("monitoring", "DeleteAlarmMuteRule", s.DeleteAlarmMuteRule)
	d.RegisterHandlerForService("monitoring", "GetAlarmMuteRule", s.GetAlarmMuteRule)
	d.RegisterHandlerForService("monitoring", "ListAlarmMuteRules", s.ListAlarmMuteRules)
	d.RegisterHandlerForService("monitoring", "GetOTelEnrichment", s.GetOTelEnrichment)
	d.RegisterHandlerForService("monitoring", "StartOTelEnrichment", s.StartOTelEnrichment)
	d.RegisterHandlerForService("monitoring", "StopOTelEnrichment", s.StopOTelEnrichment)
	d.RegisterHandlerForService("monitoring", "GetDataset", s.GetDataset)
	d.RegisterHandlerForService("monitoring", "AssociateDatasetKmsKey", s.AssociateDatasetKmsKey)
	d.RegisterHandlerForService("monitoring", "DisassociateDatasetKmsKey", s.DisassociateDatasetKmsKey)
	d.RegisterHandlerForService("monitoring", "DescribeAlarmContributors", s.DescribeAlarmContributors)
	d.RegisterHandlerForService("monitoring", "PutLogAlarm", s.PutLogAlarm)
}
