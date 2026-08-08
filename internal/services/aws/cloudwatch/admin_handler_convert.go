package cloudwatch

import (
	"google.golang.org/protobuf/proto"

	pb "vorpalstacks/internal/pb/aws/cloudwatch"
	cloudwatchstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/utils/timeutils"
)

// toPbMetricAlarm converts a store Alarm to the proto MetricAlarm type
// for the admin gRPC-Web console.
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
		Evaluationperiods:                  proto.Int32(alarm.EvaluationPeriods),
		Period:                             proto.Int32(alarm.Period),
		Statistic:                          toPbStatistic(alarm.Statistic),
		Treatmissingdata:                   alarm.TreatMissingData,
		Statevalue:                         toPbStateValue(alarm.State),
		Stateupdatedtimestamp:              alarm.StateUpdatedTimestamp.Format(timeutils.ISO8601UTCFormat),
		Alarmconfigurationupdatedtimestamp: alarm.CreatedAt.Format(timeutils.ISO8601UTCFormat),
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

func fromPbComparisonOperator(op pb.ComparisonOperator) string {
	switch op {
	case pb.ComparisonOperator_COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD:
		return "GreaterThanOrEqualToThreshold"
	case pb.ComparisonOperator_COMPARISON_OPERATOR_GREATERTHANTHRESHOLD:
		return "GreaterThanThreshold"
	case pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANTHRESHOLD:
		return "LessThanThreshold"
	case pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD:
		return "LessThanOrEqualToThreshold"
	case pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD:
		return "LessThanLowerOrGreaterThanUpperThreshold"
	case pb.ComparisonOperator_COMPARISON_OPERATOR_LESSTHANLOWERTHRESHOLD:
		return "LessThanLowerThreshold"
	case pb.ComparisonOperator_COMPARISON_OPERATOR_GREATERTHANUPPERTHRESHOLD:
		return "GreaterThanUpperThreshold"
	default:
		return "GreaterThanOrEqualToThreshold"
	}
}

func fromPbStatistic(stat pb.Statistic) string {
	switch stat {
	case pb.Statistic_STATISTIC_AVERAGE:
		return "Average"
	case pb.Statistic_STATISTIC_SUM:
		return "Sum"
	case pb.Statistic_STATISTIC_SAMPLECOUNT:
		return "SampleCount"
	case pb.Statistic_STATISTIC_MAXIMUM:
		return "Maximum"
	case pb.Statistic_STATISTIC_MINIMUM:
		return "Minimum"
	default:
		return "Average"
	}
}

// dimensionsToStore converts proto Dimensions to store Dimension type.
func dimensionsToStore(pbDims []*pb.Dimension) []cloudwatchstore.Dimension {
	dims := make([]cloudwatchstore.Dimension, len(pbDims))
	for i, d := range pbDims {
		dims[i] = cloudwatchstore.Dimension{Name: d.Name, Value: d.Value}
	}
	return dims
}

// dimensionFiltersToStore converts proto DimensionFilters to store
// Dimension type (used by ListMetrics which uses DimensionFilter).
func dimensionFiltersToStore(pbDims []*pb.DimensionFilter) []cloudwatchstore.Dimension {
	dims := make([]cloudwatchstore.Dimension, len(pbDims))
	for i, d := range pbDims {
		dims[i] = cloudwatchstore.Dimension{Name: d.Name, Value: d.Value}
	}
	return dims
}

// metricsToPb converts store MetricDatum to proto Metric type.
func metricsToPb(metrics []cloudwatchstore.MetricDatum) []*pb.Metric {
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
	return pbMetrics
}
