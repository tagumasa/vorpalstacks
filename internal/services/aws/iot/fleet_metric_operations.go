package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, err := s.bulkCreate(reqCtx, "fleetMetric", req, "metricName", map[string]interface{}{
		"queryString":      request.GetParamCaseInsensitive(req.Parameters, "queryString"),
		"aggregationType":  request.GetMapParamCaseInsensitive(req.Parameters, "aggregationType"),
		"period":           int64(request.GetIntParam(req.Parameters, "period")),
		"aggregationField": request.GetParamCaseInsensitive(req.Parameters, "aggregationField"),
		"unit":             request.GetParamCaseInsensitive(req.Parameters, "unit"),
		"indexName":        request.GetParamCaseInsensitive(req.Parameters, "indexName"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName": rec["name"],
		"metricArn":  iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
func (s *IoTService) DeleteFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "fleetMetric", req, "metricName"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "fleetMetric", req, "metricName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"queryString":      rec["queryString"],
		"aggregationType":  rec["aggregationType"],
		"period":           rec["period"],
		"aggregationField": rec["aggregationField"],
		"unit":             rec["unit"],
		"indexName":        rec["indexName"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
		"version":          int64(1),
	}, nil
}
func (s *IoTService) ListFleetMetrics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	items, err := s.bulkList(reqCtx, "fleetMetric")
	if err != nil {
		return nil, err
	}
	// Transform internal records to the AWS FleetMetricNameAndArn summary
	// shape. Without this the response items are empty objects because the
	// internal field names ("name") do not match the expected output members
	// ("metricName", "metricArn").
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		name := bulkName(item)
		result = append(result, map[string]interface{}{
			"metricName": name,
			"metricArn":  iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
		})
	}
	return paginatedMaps("fleetMetrics", result, req.Parameters), nil
}
func (s *IoTService) UpdateFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "fleetMetric", req, "metricName", map[string]interface{}{
		"queryString":      request.GetParamCaseInsensitive(req.Parameters, "queryString"),
		"aggregationType":  request.GetMapParamCaseInsensitive(req.Parameters, "aggregationType"),
		"period":           int64(request.GetIntParam(req.Parameters, "period")),
		"aggregationField": request.GetParamCaseInsensitive(req.Parameters, "aggregationField"),
		"unit":             request.GetParamCaseInsensitive(req.Parameters, "unit"),
		"indexName":        request.GetParamCaseInsensitive(req.Parameters, "indexName"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
