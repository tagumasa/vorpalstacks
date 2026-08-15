package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---- Custom Metrics -----------------------------------------------

func (s *IoTService) CreateCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	rec, err := s.bulkCreate(reqCtx, "customMetric", req, "metricName", map[string]interface{}{
		"metricType":         request.GetParamCaseInsensitive(req.Parameters, "metricType"),
		"displayName":        request.GetParamCaseInsensitive(req.Parameters, "displayName"),
		"clientRequestToken": request.GetParamCaseInsensitive(req.Parameters, "clientRequestToken"),
		"tags":               recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName": rec["name"],
		"metricArn":  iotstore.BuildCustomMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
func (s *IoTService) DeleteCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "customMetric", req, "metricName"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "customMetric", req, "metricName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCustomMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildCustomMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"metricType":       rec["metricType"],
		"displayName":      rec["displayName"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
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
	return paginatedStrings("metricNames", names, req.Parameters), nil
}
func (s *IoTService) UpdateCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "customMetric", req, "metricName", map[string]interface{}{
		"displayName": request.GetParamCaseInsensitive(req.Parameters, "displayName"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCustomMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildCustomMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"metricType":       rec["metricType"],
		"displayName":      rec["displayName"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
