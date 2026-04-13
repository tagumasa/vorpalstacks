// Package cloudwatch provides CloudWatch service operations for vorpalstacks.
package cloudwatch

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// PutMetricData publishes metric data points to CloudWatch.
func (s *CloudWatchService) PutMetricData(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namespace := request.GetStringParam(req.Parameters, "Namespace")
	if namespace == "" {
		namespace = request.GetStringParam(req.Parameters, "namespace")
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

	var statistics []string
	if statsRaw, ok := req.Parameters["Statistics"]; ok {
		statistics = parseStatistics(statsRaw)
	} else if statsRaw, ok := req.Parameters["statistics"]; ok {
		statistics = parseStatistics(statsRaw)
	}

	dimensions := parseDimensions(req.Parameters["Dimensions"], req.Parameters["dimensions"])

	query := cwstore.MetricQuery{
		Namespace:  namespace,
		MetricName: metricName,
		Dimensions: dimensions,
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     period,
		Statistics: statistics,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stats, err := store.metrics.GetMetricStatistics(query)
	if err != nil {
		return nil, err
	}

	datapoints := buildDatapointResponse(stats)
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
	namespace := request.GetStringParam(req.Parameters, "Namespace")
	if namespace == "" {
		namespace = request.GetStringParam(req.Parameters, "namespace")
	}
	metricName := request.GetStringParam(req.Parameters, "MetricName")
	if metricName == "" {
		metricName = request.GetStringParam(req.Parameters, "metricName")
	}

	dimensions := parseDimensions(req.Parameters["Dimensions"], req.Parameters["dimensions"])

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	metrics, err := store.metrics.ListMetrics(namespace, metricName, dimensions)
	if err != nil {
		return nil, err
	}

	return buildListMetricsResponse(namespace, metrics), nil
}
