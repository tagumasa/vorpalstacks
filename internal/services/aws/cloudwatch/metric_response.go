package cloudwatch

import (
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

func buildDatapointResponse(stats []*cwstore.MetricStatistics) []map[string]interface{} {
	datapoints := make([]map[string]interface{}, len(stats))
	for i, dp := range stats {
		point := map[string]interface{}{
			"Timestamp": dp.Timestamp.UnixMilli(),
		}
		point["SampleCount"] = dp.SampleCount
		point["Average"] = dp.Average
		point["Sum"] = dp.Sum
		point["Minimum"] = dp.Minimum
		point["Maximum"] = dp.Maximum
		if dp.Unit != "" {
			point["Unit"] = string(dp.Unit)
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
	var timestamps []float64
	var values []float64
	for _, dp := range stats {
		timestamps = append(timestamps, float64(dp.Timestamp.UnixMilli()))
		var val float64
		switch query.MetricStat.Stat {
		case "Average":
			val = dp.Average
		case "Sum":
			val = dp.Sum
		case "Minimum":
			val = dp.Minimum
		case "Maximum":
			val = dp.Maximum
		case "SampleCount":
			val = dp.SampleCount
		default:
			val = dp.Average
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
