package cloudwatch

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/monitoring"
	monitoringconnect "vorpalstacks/internal/pb/aws/monitoring/monitoringconnect"
	cloudwatchstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// AdminHandler provides CloudWatch service administration functionality.
// It implements the CloudWatchServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	alarmStore  *cloudwatchstore.AlarmStore
	metricStore *cloudwatchstore.MetricChunkStore
}

var _ monitoringconnect.CloudWatchServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudWatch AdminHandler.
// It initialises the handler with the provided alarm store and metric store.
func NewAdminHandler(alarmStore *cloudwatchstore.AlarmStore, metricStore *cloudwatchstore.MetricChunkStore) *AdminHandler {
	return &AdminHandler{
		alarmStore:  alarmStore,
		metricStore: metricStore,
	}
}

// ListMetrics lists metrics in CloudWatch.
// It returns metrics for the specified namespace and metric name.
func (h *AdminHandler) ListMetrics(ctx context.Context, req *connect.Request[pb.ListMetricsInput]) (*connect.Response[pb.ListMetricsOutput], error) {
	var dimensions []cloudwatchstore.Dimension
	for _, d := range req.Msg.Dimensions {
		dimensions = append(dimensions, cloudwatchstore.Dimension{
			Name:  d.Name,
			Value: d.Value,
		})
	}

	metrics, err := h.metricStore.ListMetrics(req.Msg.Namespace, req.Msg.Metricname, dimensions)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbMetrics := make([]*pb.Metric, len(metrics))
	for i, m := range metrics {
		pbDims := make([]*pb.Dimension, len(m.Dimensions))
		for j, d := range m.Dimensions {
			pbDims[j] = &pb.Dimension{
				Name:  d.Name,
				Value: d.Value,
			}
		}
		pbMetrics[i] = &pb.Metric{
			Namespace:  m.Namespace,
			Metricname: m.MetricName,
			Dimensions: pbDims,
		}
	}

	return connect.NewResponse(&pb.ListMetricsOutput{
		Metrics: pbMetrics,
	}), nil
}

// DescribeAlarms lists alarms in CloudWatch.
// It supports filtering by alarm name prefix.
func (h *AdminHandler) DescribeAlarms(ctx context.Context, req *connect.Request[pb.DescribeAlarmsInput]) (*connect.Response[pb.DescribeAlarmsOutput], error) {
	alarms, err := h.alarmStore.ListAlarms(req.Msg.Alarmnameprefix)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbAlarms := make([]*pb.MetricAlarm, 0, len(alarms))
	for _, a := range alarms {
		pbAlarms = append(pbAlarms, toPbMetricAlarm(a))
	}

	return connect.NewResponse(&pb.DescribeAlarmsOutput{
		Metricalarms: pbAlarms,
	}), nil
}

func toPbMetricAlarm(alarm *cloudwatchstore.Alarm) *pb.MetricAlarm {
	pbDims := make([]*pb.Dimension, len(alarm.Dimensions))
	for i, d := range alarm.Dimensions {
		pbDims[i] = &pb.Dimension{
			Name:  d.Name,
			Value: d.Value,
		}
	}

	return &pb.MetricAlarm{
		Alarmname:                          alarm.Name,
		Alarmarn:                           alarm.ARN,
		Namespace:                          alarm.Namespace,
		Metricname:                         alarm.MetricName,
		Dimensions:                         pbDims,
		Comparisonoperator:                 toPbComparisonOperator(alarm.ComparisonOperator),
		Threshold:                          alarm.Threshold,
		Evaluationperiods:                  alarm.EvaluationPeriods,
		Period:                             alarm.Period,
		Statistic:                          toPbStatistic(alarm.Statistic),
		Treatmissingdata:                   alarm.TreatMissingData,
		Statevalue:                         toPbStateValue(alarm.State),
		Stateupdatedtimestamp:              alarm.StateUpdatedTimestamp.Format("2006-01-02T15:04:05Z"),
		Alarmconfigurationupdatedtimestamp: alarm.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toPbComparisonOperator(op string) pb.ComparisonOperator {
	switch op {
	case "GreaterThanOrEqualToThreshold":
		return pb.ComparisonOperator_COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD
	case "GreaterThanThreshold":
		return pb.ComparisonOperator_COMPARISON_OPERATOR_GREATERTHANTHRESHOLD
	case "LessThanThreshold":
		return pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANTHRESHOLD
	case "LessThanOrEqualToThreshold":
		return pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD
	case "LessThanLowerOrGreaterThanUpperThreshold":
		return pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD
	case "LessThanLowerThreshold":
		return pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANLOWERTHRESHOLD
	case "GreaterThanUpperThreshold":
		return pb.ComparisonOperator_COMPARISON_OPERATOR_GREATERTHANUPPERTHRESHOLD
	default:
		return pb.ComparisonOperator_COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD
	}
}

func toPbStatistic(stat string) pb.Statistic {
	switch stat {
	case "Average":
		return pb.Statistic_STATISTIC_AVERAGE
	case "Sum":
		return pb.Statistic_STATISTIC_SUM
	case "SampleCount":
		return pb.Statistic_STATISTIC_SAMPLECOUNT
	case "Maximum":
		return pb.Statistic_STATISTIC_MAXIMUM
	case "Minimum":
		return pb.Statistic_STATISTIC_MINIMUM
	default:
		return pb.Statistic_STATISTIC_AVERAGE
	}
}

func toPbStateValue(state string) pb.StateValue {
	switch state {
	case "OK":
		return pb.StateValue_STATE_VALUE_OK
	case "ALARM":
		return pb.StateValue_STATE_VALUE_ALARM
	case "INSUFFICIENT_DATA":
		return pb.StateValue_STATE_VALUE_INSUFFICIENT_DATA
	default:
		return pb.StateValue_STATE_VALUE_INSUFFICIENT_DATA
	}
}

// DeleteAlarms deletes one or more alarms.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAlarms(ctx context.Context, req *connect.Request[pb.DeleteAlarmsInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteAlarmMuteRule deletes an alarm mute rule.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAlarmMuteRule(ctx context.Context, req *connect.Request[pb.DeleteAlarmMuteRuleInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteAnomalyDetector deletes an anomaly detector.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAnomalyDetector(ctx context.Context, req *connect.Request[pb.DeleteAnomalyDetectorInput]) (*connect.Response[pb.DeleteAnomalyDetectorOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteDashboards deletes one or more dashboards.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDashboards(ctx context.Context, req *connect.Request[pb.DeleteDashboardsInput]) (*connect.Response[pb.DeleteDashboardsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteInsightRules deletes one or more insight rules.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteInsightRules(ctx context.Context, req *connect.Request[pb.DeleteInsightRulesInput]) (*connect.Response[pb.DeleteInsightRulesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteMetricStream deletes a metric stream.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteMetricStream(ctx context.Context, req *connect.Request[pb.DeleteMetricStreamInput]) (*connect.Response[pb.DeleteMetricStreamOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DescribeAlarmContributors describes contributors for an alarm.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeAlarmContributors(ctx context.Context, req *connect.Request[pb.DescribeAlarmContributorsInput]) (*connect.Response[pb.DescribeAlarmContributorsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeAlarmHistory describes the history of an alarm.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeAlarmHistory(ctx context.Context, req *connect.Request[pb.DescribeAlarmHistoryInput]) (*connect.Response[pb.DescribeAlarmHistoryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeAlarmsForMetric describes alarms for a metric.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeAlarmsForMetric(ctx context.Context, req *connect.Request[pb.DescribeAlarmsForMetricInput]) (*connect.Response[pb.DescribeAlarmsForMetricOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeAnomalyDetectors describes anomaly detectors.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeAnomalyDetectors(ctx context.Context, req *connect.Request[pb.DescribeAnomalyDetectorsInput]) (*connect.Response[pb.DescribeAnomalyDetectorsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeInsightRules describes insight rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeInsightRules(ctx context.Context, req *connect.Request[pb.DescribeInsightRulesInput]) (*connect.Response[pb.DescribeInsightRulesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableAlarmActions disables actions for an alarm.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DisableAlarmActions(ctx context.Context, req *connect.Request[pb.DisableAlarmActionsInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DisableInsightRules disables insight rules.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DisableInsightRules(ctx context.Context, req *connect.Request[pb.DisableInsightRulesInput]) (*connect.Response[pb.DisableInsightRulesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// EnableAlarmActions enables actions for an alarm.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) EnableAlarmActions(ctx context.Context, req *connect.Request[pb.EnableAlarmActionsInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// EnableInsightRules enables insight rules.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) EnableInsightRules(ctx context.Context, req *connect.Request[pb.EnableInsightRulesInput]) (*connect.Response[pb.EnableInsightRulesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// GetAlarmMuteRule returns the mute rule for an alarm.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAlarmMuteRule(ctx context.Context, req *connect.Request[pb.GetAlarmMuteRuleInput]) (*connect.Response[pb.GetAlarmMuteRuleOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetDashboard returns details about a dashboard.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDashboard(ctx context.Context, req *connect.Request[pb.GetDashboardInput]) (*connect.Response[pb.GetDashboardOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetInsightRuleReport returns a report for an insight rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetInsightRuleReport(ctx context.Context, req *connect.Request[pb.GetInsightRuleReportInput]) (*connect.Response[pb.GetInsightRuleReportOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetMetricData returns metric data for a query.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMetricData(ctx context.Context, req *connect.Request[pb.GetMetricDataInput]) (*connect.Response[pb.GetMetricDataOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetMetricStatistics returns statistics for a metric.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMetricStatistics(ctx context.Context, req *connect.Request[pb.GetMetricStatisticsInput]) (*connect.Response[pb.GetMetricStatisticsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetMetricStream returns details about a metric stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMetricStream(ctx context.Context, req *connect.Request[pb.GetMetricStreamInput]) (*connect.Response[pb.GetMetricStreamOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetMetricWidgetImage returns an image of a metric widget.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMetricWidgetImage(ctx context.Context, req *connect.Request[pb.GetMetricWidgetImageInput]) (*connect.Response[pb.GetMetricWidgetImageOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListAlarmMuteRules lists alarm mute rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAlarmMuteRules(ctx context.Context, req *connect.Request[pb.ListAlarmMuteRulesInput]) (*connect.Response[pb.ListAlarmMuteRulesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListDashboards lists dashboards.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDashboards(ctx context.Context, req *connect.Request[pb.ListDashboardsInput]) (*connect.Response[pb.ListDashboardsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListManagedInsightRules lists managed insight rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListManagedInsightRules(ctx context.Context, req *connect.Request[pb.ListManagedInsightRulesInput]) (*connect.Response[pb.ListManagedInsightRulesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListMetricStreams lists metric streams.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListMetricStreams(ctx context.Context, req *connect.Request[pb.ListMetricStreamsInput]) (*connect.Response[pb.ListMetricStreamsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTagsForResource lists tags for a CloudWatch resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceInput]) (*connect.Response[pb.ListTagsForResourceOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutAlarmMuteRule creates or updates an alarm mute rule.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutAlarmMuteRule(ctx context.Context, req *connect.Request[pb.PutAlarmMuteRuleInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutAnomalyDetector creates or updates an anomaly detector.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutAnomalyDetector(ctx context.Context, req *connect.Request[pb.PutAnomalyDetectorInput]) (*connect.Response[pb.PutAnomalyDetectorOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutCompositeAlarm creates or updates a composite alarm.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutCompositeAlarm(ctx context.Context, req *connect.Request[pb.PutCompositeAlarmInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutDashboard creates or updates a dashboard.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutDashboard(ctx context.Context, req *connect.Request[pb.PutDashboardInput]) (*connect.Response[pb.PutDashboardOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutInsightRule creates or updates a contributor insight rule.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutInsightRule(ctx context.Context, req *connect.Request[pb.PutInsightRuleInput]) (*connect.Response[pb.PutInsightRuleOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutManagedInsightRules creates or updates managed insight rules.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutManagedInsightRules(ctx context.Context, req *connect.Request[pb.PutManagedInsightRulesInput]) (*connect.Response[pb.PutManagedInsightRulesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutMetricAlarm creates or updates a metric alarm.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutMetricAlarm(ctx context.Context, req *connect.Request[pb.PutMetricAlarmInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutMetricData publishes metric data points to CloudWatch.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutMetricData(ctx context.Context, req *connect.Request[pb.PutMetricDataInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutMetricStream creates or updates a metric stream.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PutMetricStream(ctx context.Context, req *connect.Request[pb.PutMetricStreamInput]) (*connect.Response[pb.PutMetricStreamOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// SetAlarmState temporarily sets the state of an alarm.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) SetAlarmState(ctx context.Context, req *connect.Request[pb.SetAlarmStateInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// StartMetricStreams starts the specified metric streams.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) StartMetricStreams(ctx context.Context, req *connect.Request[pb.StartMetricStreamsInput]) (*connect.Response[pb.StartMetricStreamsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// StopMetricStreams stops the specified metric streams.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) StopMetricStreams(ctx context.Context, req *connect.Request[pb.StopMetricStreamsInput]) (*connect.Response[pb.StopMetricStreamsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// TagResource adds tags to a CloudWatch resource.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceInput]) (*connect.Response[pb.TagResourceOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// UntagResource removes tags from a CloudWatch resource.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceInput]) (*connect.Response[pb.UntagResourceOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}
