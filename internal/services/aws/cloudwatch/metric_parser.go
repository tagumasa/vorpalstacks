// Package cloudwatch provides CloudWatch service operations for vorpalstacks.
package cloudwatch

import (
	"time"
	"unicode"
	"unicode/utf8"

	"vorpalstacks/internal/services/aws/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

func getFloatValue(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case float32:
			return float64(val)
		case int:
			return float64(val)
		case int64:
			return float64(val)
		}
	}
	return 0
}

func parseMetricData(metricDataRaw interface{}) []cwstore.MetricDatum {
	var metricData []cwstore.MetricDatum
	metricDataList, ok := metricDataRaw.([]interface{})
	if !ok {
		return metricData
	}
	for _, md := range metricDataList {
		mdMap, ok := md.(map[string]interface{})
		if !ok {
			continue
		}
		metricName := request.GetStringParam(mdMap, "MetricName")
		if metricName == "" {
			metricName = request.GetStringParam(mdMap, "metricName")
		}

		datum := cwstore.MetricDatum{
			MetricName: metricName,
			Value:      getFloatValue(mdMap, "Value"),
			Unit:       cwstore.StandardUnit(request.GetStringParam(mdMap, "Unit")),
		}
		if datum.Value == 0 {
			datum.Value = getFloatValue(mdMap, "value")
		}
		if datum.Unit == "" {
			datum.Unit = cwstore.StandardUnit(request.GetStringParam(mdMap, "unit"))
		}
		datum.StorageResolution = int32(request.GetIntParam(mdMap, "StorageResolution"))
		if datum.StorageResolution == 0 {
			datum.StorageResolution = int32(request.GetIntParam(mdMap, "storageResolution"))
		}

		if ts := parseTimestampFromMap(mdMap, "Timestamp"); !ts.IsZero() {
			datum.Timestamp = ts
		} else if ts := parseTimestampFromMap(mdMap, "timestamp"); !ts.IsZero() {
			datum.Timestamp = ts
		}

		datum.Dimensions = parseDimensions(mdMap["Dimensions"], mdMap["dimensions"])
		datum.Values = parseFloatValues(mdMap["Values"], mdMap["values"])
		datum.Counts = parseFloatValues(mdMap["Counts"], mdMap["counts"])

		if svRaw, ok := mdMap["StatisticValues"]; ok {
			if svMap, ok := svRaw.(map[string]interface{}); ok {
				datum.StatisticValues = &cwstore.StatisticSet{
					SampleCount: getFloatValue(svMap, "SampleCount"),
					Sum:         getFloatValue(svMap, "Sum"),
					Minimum:     getFloatValue(svMap, "Minimum"),
					Maximum:     getFloatValue(svMap, "Maximum"),
				}
			}
		} else if svRaw, ok := mdMap["statisticValues"]; ok {
			if svMap, ok := svRaw.(map[string]interface{}); ok {
				datum.StatisticValues = &cwstore.StatisticSet{
					SampleCount: getFloatValue(svMap, "sampleCount"),
					Sum:         getFloatValue(svMap, "sum"),
					Minimum:     getFloatValue(svMap, "minimum"),
					Maximum:     getFloatValue(svMap, "maximum"),
				}
			}
		}

		metricData = append(metricData, datum)
	}
	return metricData
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func parseDimensions(dimsRaw1, dimsRaw2 interface{}) []cwstore.Dimension {
	var dimensions []cwstore.Dimension
	if dimsList, ok := dimsRaw1.([]interface{}); ok {
		for _, d := range dimsList {
			if dimMap, ok := d.(map[string]interface{}); ok {
				dimensions = append(dimensions, cwstore.Dimension{
					Name:  request.GetStringParam(dimMap, "Name"),
					Value: request.GetStringParam(dimMap, "Value"),
				})
			}
		}
	} else if dimsList, ok := dimsRaw2.([]interface{}); ok {
		for _, d := range dimsList {
			if dimMap, ok := d.(map[string]interface{}); ok {
				dimensions = append(dimensions, cwstore.Dimension{
					Name:  request.GetStringParam(dimMap, "name"),
					Value: request.GetStringParam(dimMap, "value"),
				})
			}
		}
	}
	return dimensions
}

func parseFloatValues(valuesRaw1, valuesRaw2 interface{}) []float64 {
	var values []float64
	if valuesList, ok := valuesRaw1.([]interface{}); ok {
		for _, v := range valuesList {
			if vf, ok := v.(float64); ok {
				values = append(values, vf)
			}
		}
	} else if valuesList, ok := valuesRaw2.([]interface{}); ok {
		for _, v := range valuesList {
			if vf, ok := v.(float64); ok {
				values = append(values, vf)
			}
		}
	}
	return values
}

func parseStatistics(statsRaw interface{}) []string {
	var statistics []string
	if statsList, ok := statsRaw.([]interface{}); ok {
		for _, s := range statsList {
			if ss, ok := s.(string); ok {
				statistics = append(statistics, ss)
			}
		}
	}
	return statistics
}

func parseMetricStat(metricStatRaw interface{}) *cwstore.MetricStat {
	if metricStat, ok := metricStatRaw.(map[string]interface{}); ok {
		stat := &cwstore.MetricStat{
			Period: int32(request.GetIntParam(metricStat, "Period")),
			Stat:   request.GetStringParam(metricStat, "Stat"),
			Unit:   cwstore.StandardUnit(request.GetStringParam(metricStat, "Unit")),
		}
		if stat.Period == 0 {
			stat.Period = int32(request.GetIntParam(metricStat, "period"))
		}
		if stat.Stat == "" {
			stat.Stat = request.GetStringParam(metricStat, "stat")
		}
		if stat.Unit == "" {
			stat.Unit = cwstore.StandardUnit(request.GetStringParam(metricStat, "unit"))
		}
		if metric, ok := metricStat["Metric"].(map[string]interface{}); ok {
			stat.Metric = cwstore.MetricRef{
				Namespace:  request.GetStringParam(metric, "Namespace"),
				MetricName: request.GetStringParam(metric, "MetricName"),
			}
			if stat.Metric.Namespace == "" {
				stat.Metric.Namespace = request.GetStringParam(metric, "namespace")
			}
			if stat.Metric.MetricName == "" {
				stat.Metric.MetricName = request.GetStringParam(metric, "metricName")
			}
			stat.Metric.Dimensions = parseDimensions(metric["Dimensions"], metric["dimensions"])
		} else if metric, ok := metricStat["metric"].(map[string]interface{}); ok {
			stat.Metric = cwstore.MetricRef{
				Namespace:  request.GetStringParam(metric, "namespace"),
				MetricName: request.GetStringParam(metric, "metricName"),
			}
			if stat.Metric.Namespace == "" {
				stat.Metric.Namespace = request.GetStringParam(metric, "Namespace")
			}
			if stat.Metric.MetricName == "" {
				stat.Metric.MetricName = request.GetStringParam(metric, "MetricName")
			}
			stat.Metric.Dimensions = parseDimensions(metric["Dimensions"], metric["dimensions"])
		}
		return stat
	}
	return nil
}

func parseMetricDataQueries(queriesRaw interface{}) []cwstore.MetricDataQuery {
	var metricDataQueries []cwstore.MetricDataQuery
	if queriesList, ok := queriesRaw.([]interface{}); ok {
		for _, q := range queriesList {
			if qMap, ok := q.(map[string]interface{}); ok {
				query := cwstore.MetricDataQuery{
					Id: request.GetStringParam(qMap, "Id"),
				}
				if query.Id == "" {
					query.Id = request.GetStringParam(qMap, "id")
				}

				if metricStatRaw, ok := qMap["MetricStat"]; ok {
					query.MetricStat = parseMetricStat(metricStatRaw)
				} else if metricStatRaw, ok := qMap["metricStat"]; ok {
					query.MetricStat = parseMetricStat(metricStatRaw)
				}

				metricDataQueries = append(metricDataQueries, query)
			}
		}
	}
	return metricDataQueries
}

func parseTimestampFromMap(m map[string]interface{}, key string) time.Time {
	if v, ok := m[key]; ok {
		switch ts := v.(type) {
		case time.Time:
			return ts
		case int:
			return time.Unix(int64(ts), 0)
		case int64:
			return time.Unix(ts, 0)
		case float64:
			return time.Unix(int64(ts), 0)
		}
	}
	if ts := request.GetIntParam(m, key); ts > 0 {
		return time.Unix(int64(ts), 0)
	}
	var timeStr string
	timeStr = request.GetStringParam(m, key)
	if timeStr == "" {
		timeStr = request.GetStringParam(m, lowerCaseFirst(key))
	}
	if timeStr != "" {
		var err error
		ts, err := time.Parse(time.RFC3339, timeStr)
		if err == nil {
			return ts
		}
	}
	return time.Time{}
}

func lowerCaseFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToLower(r)) + s[size:]
}
