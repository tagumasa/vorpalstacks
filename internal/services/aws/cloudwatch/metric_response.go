package cloudwatch

import (
	"strings"
	"time"

	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/pkg/metricmath"
)

func buildDatapointResponse(stats []*cwstore.MetricStatistics, requestedStats []string) []map[string]interface{} {
	requested := make(map[string]bool, len(requestedStats))
	for _, s := range requestedStats {
		requested[strings.ToLower(s)] = true
	}
	allRequested := len(requested) == 0

	datapoints := make([]map[string]interface{}, len(stats))
	for i, dp := range stats {
		point := map[string]interface{}{
			"Timestamp": dp.Timestamp.UnixMilli(),
		}
		for stat, include := range requested {
			if !include && !allRequested {
				continue
			}
			switch stat {
			case "samplecount":
				point["SampleCount"] = dp.SampleCount
			case "average":
				point["Average"] = dp.Average
			case "sum":
				point["Sum"] = dp.Sum
			case "minimum":
				point["Minimum"] = dp.Minimum
			case "maximum":
				point["Maximum"] = dp.Maximum
			}
		}
		if allRequested {
			point["SampleCount"] = dp.SampleCount
			point["Average"] = dp.Average
			point["Sum"] = dp.Sum
			point["Minimum"] = dp.Minimum
			point["Maximum"] = dp.Maximum
		}
		if dp.Unit != "" {
			point["Unit"] = string(dp.Unit)
		}
		if len(dp.ExtendedStats) > 0 {
			for statName, statVal := range dp.ExtendedStats {
				point[statName] = statVal
			}
		}
		datapoints[i] = point
	}
	return datapoints
}

func buildListMetricsResponse(namespace string, metrics []cwstore.MetricDatum) map[string]interface{} {
	result := make([]map[string]interface{}, len(metrics))
	for i, m := range metrics {
		ns := namespace
		if ns == "" && m.Namespace != "" {
			ns = m.Namespace
		}
		metric := map[string]interface{}{
			"Namespace":  ns,
			"MetricName": m.MetricName,
		}
		if len(m.Dimensions) > 0 {
			dims := make([]map[string]interface{}, len(m.Dimensions))
			for j, d := range m.Dimensions {
				dims[j] = map[string]interface{}{
					"Name":  d.Name,
					"Value": d.Value,
				}
			}
			metric["Dimensions"] = dims
		}
		result[i] = metric
	}
	return map[string]interface{}{
		"Metrics": result,
	}
}

func buildMetricDataResult(query cwstore.MetricDataQuery, stats []*cwstore.MetricStatistics) map[string]interface{} {
	var timestamps []time.Time
	var values []float64
	statName := query.MetricStat.Stat
	statLower := strings.ToLower(statName)
	for _, dp := range stats {
		timestamps = append(timestamps, dp.Timestamp)
		var val float64
		switch statLower {
		case "average":
			val = dp.Average
		case "sum":
			val = dp.Sum
		case "minimum":
			val = dp.Minimum
		case "maximum":
			val = dp.Maximum
		case "samplecount":
			val = dp.SampleCount
		default:
			// Extended statistic (percentile, IQM, etc.)
			if dp.ExtendedStats != nil {
				if ev, ok := dp.ExtendedStats[statName]; ok {
					val = ev
				}
			}
		}
		values = append(values, val)
	}

	return map[string]interface{}{
		"Id":         query.Id,
		"Label":      query.MetricStat.Metric.MetricName,
		"Timestamps": timestamps,
		"Values":     values,
		"StatusCode": "Complete",
	}
}

// buildMetricDataResultFromPoints constructs the response entry for a
// MetricStat query from pre-evaluated DataPoints. This avoids the lossy
// conversion through MetricStatistics that only preserves Average.
func buildMetricDataResultFromPoints(query cwstore.MetricDataQuery, dataPoints []metricmath.DataPoint) map[string]interface{} {
	timestamps := make([]time.Time, 0, len(dataPoints))
	values := make([]float64, 0, len(dataPoints))
	for _, dp := range dataPoints {
		timestamps = append(timestamps, dp.Timestamp)
		values = append(values, dp.Value)
	}
	label := query.MetricStat.Metric.MetricName
	if query.Label != "" {
		label = query.Label
	}
	return map[string]interface{}{
		"Id":         query.Id,
		"Label":      label,
		"Timestamps": timestamps,
		"Values":     values,
		"StatusCode": "Complete",
	}
}

// buildExpressionResult constructs the response entry for a metric math
// expression query. The label is taken from the query's Label field, or
// the expression string if no label is set.
func buildExpressionResult(query cwstore.MetricDataQuery, dataPoints []metricmath.DataPoint) map[string]interface{} {
	timestamps := make([]time.Time, 0, len(dataPoints))
	values := make([]float64, 0, len(dataPoints))
	for _, dp := range dataPoints {
		timestamps = append(timestamps, dp.Timestamp)
		values = append(values, dp.Value)
	}
	label := query.Label
	if label == "" {
		label = query.Expression
	}
	return map[string]interface{}{
		"Id":         query.Id,
		"Label":      label,
		"Timestamps": timestamps,
		"Values":     values,
		"StatusCode": "Complete",
	}
}
