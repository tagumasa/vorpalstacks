// Package cloudwatch implements AWS CloudWatch API operations including
// alarms, metrics, dashboards, log groups, and metric widget image rendering.
package cloudwatch

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cloudwatch"
	cloudwatchconnect "vorpalstacks/internal/pb/aws/cloudwatch/cloudwatchconnect"
	cloudwatchstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// AdminHandler implements the CloudWatch gRPC-Web admin console handler. It
// exposes list and describe operations for alarms and metrics for the Flutter
// management UI.
type AdminHandler struct {
	cloudwatchconnect.UnimplementedCloudWatchServiceHandler
	alarmStore  *cloudwatchstore.AlarmStore
	metricStore *cloudwatchstore.MetricChunkStore
}

var _ cloudwatchconnect.CloudWatchServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudWatch admin console handler backed by the
// given alarm and metric stores.
func NewAdminHandler(alarmStore *cloudwatchstore.AlarmStore, metricStore *cloudwatchstore.MetricChunkStore) *AdminHandler {
	return &AdminHandler{
		alarmStore:  alarmStore,
		metricStore: metricStore,
	}
}

// ListMetrics lists the specified metrics within a namespace, optionally
// filtered by metric name and dimensions.
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

// DescribeAlarms retrieves the alarms for the specified alarm name prefix.
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

// NewConnectHandler creates a gRPC-Web connect handler for the CloudWatch admin console.
func NewConnectHandler(st storage.BasicStorage, accountID, region, dataPath string) (string, http.Handler) {
	return cloudwatchconnect.NewCloudWatchServiceHandler(NewAdminHandler(cloudwatchstore.NewAlarmStore(st, accountID, region), cloudwatchstore.NewMetricChunkStore(st, region, dataPath)))
}
