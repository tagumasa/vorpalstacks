package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// ---- Custom Metrics -----------------------------------------------

func (s *IoTService) CreateCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	result, err := s.createCustomMetricCore(store, CreateCustomMetricInput{
		MetricName:         request.GetParamCaseInsensitive(req.Parameters, "metricName"),
		MetricType:         request.GetParamCaseInsensitive(req.Parameters, "metricType"),
		DisplayName:        request.GetParamCaseInsensitive(req.Parameters, "displayName"),
		ClientRequestToken: request.GetParamCaseInsensitive(req.Parameters, "clientRequestToken"),
		Tags:               recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName": result.MetricName,
		"metricArn":  result.MetricArn,
	}, nil
}
func (s *IoTService) DeleteCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteCustomMetricCore(store, request.GetParamCaseInsensitive(req.Parameters, "metricName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeCustomMetricCore(store, request.GetParamCaseInsensitive(req.Parameters, "metricName"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName":       rec.Rec["name"],
		"metricArn":        rec.Arn,
		"metricType":       rec.Rec["metricType"],
		"displayName":      rec.Rec["displayName"],
		"creationDate":     rec.Rec["creationDate"],
		"lastModifiedDate": rec.Rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListCustomMetrics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy: ListCustomMetricsResponse.metricNames is a list of MetricName (string).
	items, err := s.bulkList(reqCtx, "customMetric")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		if name, ok := item["name"].(string); ok {
			names = append(names, name)
		}
	}
	return paginatedStrings("metricNames", names, req.Parameters)
}
func (s *IoTService) UpdateCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.updateCustomMetricCore(store, UpdateCustomMetricInput{
		MetricName:  request.GetParamCaseInsensitive(req.Parameters, "metricName"),
		DisplayName: request.GetParamCaseInsensitive(req.Parameters, "displayName"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName":       result.Rec["name"],
		"metricArn":        result.Arn,
		"metricType":       result.Rec["metricType"],
		"displayName":      result.Rec["displayName"],
		"creationDate":     result.Rec["creationDate"],
		"lastModifiedDate": result.Rec["lastModifiedDate"],
	}, nil
}
