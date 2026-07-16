// Package cloudwatch provides CloudWatch service operations for vorpalstacks.
package cloudwatch

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// CloudWatch PutMetricData limits per AWS spec.
const (
	maxMetricDataPerRequest = 1000
	maxValuesPerDatum       = 150
	maxMetricDimensions     = 30
)

func getMetricStringParam(params map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// PutMetricData publishes metric data points to CloudWatch.
func (s *CloudWatchService) PutMetricData(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namespace := request.GetStringParam(req.Parameters, "Namespace")
	if namespace == "" {
		namespace = request.GetStringParam(req.Parameters, "namespace")
	}
	if namespace == "" {
		return nil, ErrInvalidParameter
	}

	var metricDataRaw interface{}
	if md, ok := req.Parameters["MetricData"]; ok {
		metricDataRaw = md
	} else if md, ok := req.Parameters["metricData"]; ok {
		metricDataRaw = md
	}

	metricData := parseMetricData(metricDataRaw)
	if len(metricData) == 0 {
		return nil, ErrInvalidParameter
	}
	if len(metricData) > maxMetricDataPerRequest {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of MetricData entries must not exceed %d", maxMetricDataPerRequest))
	}
	for _, datum := range metricData {
		if err := validateMetricDatum(datum); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.metrics.PutMetricData(namespace, metricData); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// validateMetricDatum checks that a single MetricDatum adheres to AWS
// constraints: at most one of Value, Values+Counts, or StatisticValues
// is provided; Values and Counts must have matching lengths; Values
// must not exceed 150 entries; and Dimensions must not exceed 30.
func validateMetricDatum(datum cwstore.MetricDatum) error {
	if len(datum.Dimensions) > maxMetricDimensions {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("A MetricDatum can have at most %d Dimensions", maxMetricDimensions))
	}

	hasValue := datum.HasValue
	hasValues := len(datum.Values) > 0
	hasStatSet := datum.StatisticValues != nil

	// AWS rejects requests that specify more than one of Value, Values,
	// or StatisticValues simultaneously.
	modeCount := 0
	if hasValue {
		modeCount++
	}
	if hasValues {
		modeCount++
	}
	if hasStatSet {
		modeCount++
	}
	if modeCount > 1 {
		return awserrors.NewInvalidParameterValueException(
			"A MetricDatum must not specify more than one of Value, Values, or StatisticValues")
	}

	if len(datum.Values) > maxValuesPerDatum {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("A MetricDatum Values array must not exceed %d entries", maxValuesPerDatum))
	}

	if hasValues && len(datum.Counts) > 0 && len(datum.Values) != len(datum.Counts) {
		return awserrors.NewInvalidParameterValueException(
			"Values and Counts arrays must have the same length")
	}

	return nil
}

// GetMetricStatistics retrieves statistics for a metric.
func (s *CloudWatchService) GetMetricStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namespace := request.GetStringParam(req.Parameters, "Namespace")
	if namespace == "" {
		namespace = request.GetStringParam(req.Parameters, "namespace")
	}
	metricName := request.GetStringParam(req.Parameters, "MetricName")
	if metricName == "" {
		metricName = request.GetStringParam(req.Parameters, "metricName")
	}

	startTime := parseTimestampFromMap(req.Parameters, "StartTime")
	endTime := parseTimestampFromMap(req.Parameters, "EndTime")

	period := int32(request.GetIntParam(req.Parameters, "Period"))
	if period == 0 {
		period = int32(request.GetIntParam(req.Parameters, "period"))
	}

	if namespace == "" || metricName == "" || startTime.IsZero() || endTime.IsZero() {
		return nil, ErrInvalidParameter
	}

	var statistics []string
	if statsRaw, ok := req.Parameters["Statistics"]; ok {
		statistics = parseStatistics(statsRaw)
	} else if statsRaw, ok := req.Parameters["statistics"]; ok {
		statistics = parseStatistics(statsRaw)
	}

	var extendedStats []string
	if extRaw, ok := req.Parameters["ExtendedStatistics"]; ok {
		extendedStats = parseStatistics(extRaw)
	} else if extRaw, ok := req.Parameters["extendedStatistics"]; ok {
		extendedStats = parseStatistics(extRaw)
	}

	dimensions := parseDimensions(req.Parameters["Dimensions"], req.Parameters["dimensions"])

	query := cwstore.MetricQuery{
		Namespace:          namespace,
		MetricName:         metricName,
		Dimensions:         dimensions,
		StartTime:          startTime,
		EndTime:            endTime,
		Period:             period,
		Statistics:         statistics,
		ExtendedStatistics: extendedStats,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stats, err := store.metrics.GetMetricStatistics(query)
	if err != nil {
		return nil, err
	}

	datapoints := buildDatapointResponse(stats, statistics)
	if datapoints == nil {
		datapoints = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"Label":      metricName,
		"Datapoints": datapoints,
	}, nil
}

// ListMetrics returns a list of available metrics.
func (s *CloudWatchService) ListMetrics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namespace := getMetricStringParam(req.Parameters, "Namespace", "namespace")
	metricName := getMetricStringParam(req.Parameters, "MetricName", "metricName")

	dimensions := parseDimensions(req.Parameters["Dimensions"], req.Parameters["dimensions"])
	nextToken := getMetricStringParam(req.Parameters, "NextToken", "nextToken")
	maxResults := 500
	if mr := request.GetIntParam(req.Parameters, "MaxResults"); mr > 0 {
		maxResults = mr
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	metrics, nextMarker, isTruncated, err := store.metrics.ListMetricsPaginated(namespace, metricName, dimensions, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	result := buildListMetricsResponse(namespace, metrics)
	if isTruncated {
		result["NextToken"] = nextMarker
	}
	return result, nil
}
