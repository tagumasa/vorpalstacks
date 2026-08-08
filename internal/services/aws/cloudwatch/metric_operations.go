package cloudwatch

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
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

	// EntityMetricData (OTel) and StrictEntityValidation are accepted but
	// not processed — OTel entity-based metrics require the OTel pipeline.
	// The StrictEntityValidation flag is accepted as a no-op.

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

	// AWS requires Period to be a positive integer.
	if period <= 0 {
		return nil, awserrors.NewInvalidParameterValueException("Period must be a positive integer")
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

	// AWS requires at least one of Statistics or ExtendedStatistics.
	if len(statistics) == 0 && len(extendedStats) == 0 {
		return nil, awserrors.NewMissingParameter("Must specify either Statistics or ExtendedStatistics")
	}

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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input := &ListMetricsInput{
		Namespace:      getMetricStringParam(req.Parameters, "Namespace", "namespace"),
		MetricName:     getMetricStringParam(req.Parameters, "MetricName", "metricName"),
		Dimensions:     parseDimensions(req.Parameters["Dimensions"], req.Parameters["dimensions"]),
		NextToken:      getMetricStringParam(req.Parameters, "NextToken", "nextToken"),
		RecentlyActive: getMetricStringParam(req.Parameters, "RecentlyActive", "recentlyActive"),
		MaxResults:     request.GetIntParam(req.Parameters, "MaxResults"),
	}

	result, err := s.listMetricsCore(store, input)
	if err != nil {
		return nil, err
	}

	resp := buildListMetricsResponse(input.Namespace, result.Metrics)
	if result.IsTruncated {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}
