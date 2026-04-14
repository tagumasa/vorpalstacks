package cloudwatch

import (
	"context"

	"vorpalstacks/internal/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// GetMetricData retrieves metric data values for one or more metrics.
func (s *CloudWatchService) GetMetricData(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	startTime := parseTimestampFromMap(req.Parameters, "StartTime")
	endTime := parseTimestampFromMap(req.Parameters, "EndTime")

	if startTime.IsZero() || endTime.IsZero() {
		return nil, ErrInvalidParameter
	}

	var metricDataQueries []cwstore.MetricDataQuery
	if queriesRaw, ok := req.Parameters["MetricDataQueries"]; ok {
		metricDataQueries = parseMetricDataQueries(queriesRaw)
	} else if queriesRaw, ok := req.Parameters["metricDataQueries"]; ok {
		metricDataQueries = parseMetricDataQueries(queriesRaw)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	for _, query := range metricDataQueries {
		if query.MetricStat != nil {
			mq := cwstore.MetricQuery{
				Namespace:  query.MetricStat.Metric.Namespace,
				MetricName: query.MetricStat.Metric.MetricName,
				Dimensions: query.MetricStat.Metric.Dimensions,
				StartTime:  startTime,
				EndTime:    endTime,
				Period:     query.MetricStat.Period,
				Statistics: []string{query.MetricStat.Stat},
			}

			stats, err := store.metrics.GetMetricStatistics(mq)
			if err != nil {
				continue
			}

			result := buildMetricDataResult(query, stats)
			results = append(results, result)
		}
	}

	return map[string]interface{}{
		"MetricDataResults": results,
		"Messages":          []interface{}{},
	}, nil
}

// GetMetricWidgetImage retrieves an image for a metric widget.
func (s *CloudWatchService) GetMetricWidgetImage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	widgetDefinition := request.GetStringParam(req.Parameters, "metricWidget")
	if widgetDefinition == "" {
		widgetDefinition = request.GetStringParam(req.Parameters, "MetricWidget")
	}
	_ = widgetDefinition

	format := request.GetStringParam(req.Parameters, "outputFormat")
	if format == "" {
		format = "png"
	}
	_ = format

	placeholderPNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
	}

	return map[string]interface{}{
		"MetricWidgetImage": placeholderPNG,
	}, nil
}
