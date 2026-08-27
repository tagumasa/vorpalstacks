package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

func (s *IoTService) CreateFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := fleetMetricInputFromRequest(req, false)
	result, err := s.createFleetMetricCore(store, in)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName": result.Record["name"],
		"metricArn":  result.ARN,
	}, nil
}
func (s *IoTService) DeleteFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	metricName := request.GetParamCaseInsensitive(req.Parameters, "metricName")
	if err := s.deleteFleetMetricCore(store, metricName); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	metricName := request.GetParamCaseInsensitive(req.Parameters, "metricName")
	result, err := s.describeFleetMetricCore(store, metricName)
	if err != nil {
		return nil, err
	}
	rec := result.Record
	resp := map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        result.ARN,
		"queryString":      rec["queryString"],
		"aggregationType":  rec["aggregationType"],
		"period":           rec["period"],
		"aggregationField": rec["aggregationField"],
		"unit":             rec["unit"],
		"indexName":        rec["indexName"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
		"version":          fleetMetricRecordVersion(rec),
	}
	if description, ok := rec["description"]; ok {
		resp["description"] = description
	}
	if queryVersion, ok := rec["queryVersion"]; ok {
		resp["queryVersion"] = queryVersion
	}
	return resp, nil
}
func (s *IoTService) ListFleetMetrics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	summaries, err := s.listFleetMetricsCore(store)
	if err != nil {
		return nil, err
	}
	// Transform internal records to the AWS FleetMetricNameAndArn summary
	// shape. Without this the response items are empty objects because the
	// internal field names ("name") do not match the expected output members
	// ("metricName", "metricArn").
	result := make([]map[string]interface{}, 0, len(summaries))
	for _, summary := range summaries {
		result = append(result, map[string]interface{}{
			"metricName": summary.Name,
			"metricArn":  summary.ARN,
		})
	}
	return paginatedMaps("fleetMetrics", result, req.Parameters)
}
func (s *IoTService) UpdateFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := fleetMetricInputFromRequest(req, true)
	if _, err := s.updateFleetMetricCore(store, in); err != nil {
		return nil, err
	}
	// API_UpdateFleetMetric returns an HTTP 200 with an empty body.
	return map[string]interface{}{}, nil
}

// fleetMetricInputFromRequest projects the wire members onto the
// fleet-metric DTO. Presence is detected per member so the update core can
// preserve omitted optional members; the create path fills every provided
// member the same way, with requiredness enforced by the core validation.
func fleetMetricInputFromRequest(req *request.ParsedRequest, forUpdate bool) FleetMetricInput {
	in := FleetMetricInput{
		MetricName: request.GetParamCaseInsensitive(req.Parameters, "metricName"),
		IndexName:  request.GetParamCaseInsensitive(req.Parameters, "indexName"),
	}
	if str, ok := optionalStringParam(req.Parameters, "queryString"); ok {
		in.QueryString = &str
	}
	if str, ok := optionalStringParam(req.Parameters, "aggregationField"); ok {
		in.AggregationField = &str
	}
	if str, ok := optionalStringParam(req.Parameters, "unit"); ok {
		in.Unit = &str
	}
	if str, ok := optionalStringParam(req.Parameters, "description"); ok {
		in.Description = &str
	}
	if str, ok := optionalStringParam(req.Parameters, "queryVersion"); ok {
		in.QueryVersion = &str
	}
	if m := request.GetMapParamCaseInsensitive(req.Parameters, "aggregationType"); m != nil {
		in.AggregationType = m
	}
	if v, ok := request.GetIntParamCaseInsensitive(req.Parameters, "period"); ok {
		period := int64(v)
		in.Period = &period
	}
	if forUpdate {
		if v, ok := request.GetIntParamCaseInsensitive(req.Parameters, "expectedVersion"); ok {
			expected := int64(v)
			in.ExpectedVersion = &expected
		}
	} else {
		in.Tags = tagListParam(req.Parameters)
	}
	return in
}

// optionalStringParam returns the member's value when the key is present;
// the boolean distinguishes an omitted member from an empty string.
func optionalStringParam(params map[string]interface{}, key string) (string, bool) {
	if !hasParam(params, key) {
		return "", false
	}
	return request.GetParamCaseInsensitive(params, key), true
}

// hasParam checks the three key spellings the parameter getters use
// (exact, lower-first-letter, all-lowercase).
func hasParam(params map[string]interface{}, key string) bool {
	return request.HasParam(params, key) ||
		request.HasParam(params, request.LowerFirst(key)) ||
		request.HasParam(params, lowerAll(key))
}

func lowerAll(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

// tagListParam extracts the JSON-protocol TagList member
// ([{"Key": ..., "Value": ...}]) into a map.
func tagListParam(params map[string]interface{}) map[string]string {
	list := request.GetListParamLowerFirst(params, "tags")
	if len(list) == 0 {
		return nil
	}
	tags := make(map[string]string, len(list))
	for _, entry := range list {
		key, _ := entry["Key"].(string)
		if key == "" {
			key, _ = entry["key"].(string)
		}
		value, _ := entry["Value"].(string)
		if value == "" {
			value, _ = entry["value"].(string)
		}
		if key != "" {
			tags[key] = value
		}
	}
	return tags
}
