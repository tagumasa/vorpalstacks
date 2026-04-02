// Package cloudwatch provides CloudWatch service operations for vorpalstacks.
package cloudwatch

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
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
	storage   storage.BasicStorage
	accountID string
	dataPath  string
}

// NewCloudWatchService creates a new CloudWatch service.
//
// Parameters:
//   - store: The storage instance
//   - accountID: The AWS account ID
//   - region: The AWS region
//   - dataPath: The data path for chunk storage
//
// Returns:
//   - *CloudWatchService: A new CloudWatch service
func NewCloudWatchService(store storage.BasicStorage, accountID, region, dataPath string) *CloudWatchService {
	return &CloudWatchService{
		storage:   store,
		accountID: accountID,
		dataPath:  dataPath,
	}
}

// store returns the CloudWatch stores for a given request context.
func (s *CloudWatchService) store(reqCtx *request.RequestContext) (*cloudwatchStores, error) {
	var metricStore *cwstore.MetricChunkStore
	var alarmStore *cwstore.AlarmStore
	var dashboardStore *cwstore.DashboardStore

	if reqCtx.GetCloudWatchMetricStore() != nil {
		if ms, ok := reqCtx.GetCloudWatchMetricStore().(*cwstore.MetricChunkStore); ok {
			metricStore = ms
		}
	}
	if reqCtx.GetCloudWatchAlarmStore() != nil {
		if as, ok := reqCtx.GetCloudWatchAlarmStore().(*cwstore.AlarmStore); ok {
			alarmStore = as
		}
	}

	if metricStore == nil || alarmStore == nil || dashboardStore == nil {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		region := reqCtx.GetRegion()

		if metricStore == nil {
			var err error
			metricStore, err = cwstore.NewMetricChunkStoreWithIndex(storage, region, s.dataPath)
			if err != nil {
				return nil, err
			}
		}

		if alarmStore == nil {
			alarmStore = cwstore.NewAlarmStore(storage, s.accountID, region)
		}

		if dashboardStore == nil {
			dashboardStore = cwstore.NewDashboardStore(storage, s.accountID, region)
		}
	}

	return &cloudwatchStores{
		metrics:    metricStore,
		alarms:     alarmStore,
		dashboards: dashboardStore,
	}, nil
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
