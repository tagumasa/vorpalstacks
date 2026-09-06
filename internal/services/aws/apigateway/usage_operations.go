// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateApiKey creates a new API key in API Gateway.
func (s *APIGatewayService) CreateApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &ApiKeyInput{
		Name:        request.GetStringParam(req.Parameters, "name"),
		Description: request.GetStringParam(req.Parameters, "description"),
		Enabled:     true,
		CustomerId:  request.GetStringParam(req.Parameters, "customerId"),
		Value:       request.GetStringParam(req.Parameters, "value"),
	}
	if v, ok := req.Parameters["enabled"]; ok {
		if b, ok := v.(bool); ok {
			in.Enabled = b
		}
	}
	if stageKeys, ok := req.Parameters["stageKeys"].([]interface{}); ok {
		for _, sk := range stageKeys {
			if sks, ok := sk.(string); ok {
				in.StageKeys = append(in.StageKeys, sks)
			}
		}
	}
	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		in.Tags = tagutil.MapInterfaceToTags(tags)
	}
	// generateDistinctId=false forces the caller-supplied value to become
	// the key id; the Core enforces that contract for both planes.
	if v, ok := req.Parameters["generateDistinctId"].(bool); ok {
		in.GenerateDistinctId = &v
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createApiKeyCore(stores, in)
	if err != nil {
		return nil, err
	}
	return s.toApiKeyResponseWithIncludeValue(created, true), nil
}

// GetApiKey retrieves an API key from API Gateway.
func (s *APIGatewayService) GetApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiKeyId := request.GetStringParam(req.Parameters, "apiKey")
	if apiKeyId == "" {
		apiKeyId = getPathParam(req, "apiKey")
	}
	includeValue := request.GetBoolParam(req.Parameters, "includeValue")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	apiKey, err := s.getApiKeyCore(stores, apiKeyId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toApiKeyResponseWithIncludeValue(apiKey, includeValue), nil
}

// DeleteApiKey deletes an API key from API Gateway.
func (s *APIGatewayService) DeleteApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiKeyId := request.GetStringParam(req.Parameters, "apiKey")
	if apiKeyId == "" {
		apiKeyId = getPathParam(req, "apiKey")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteApiKeyCore(stores, apiKeyId); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// UpdateApiKey updates an existing API key in API Gateway.
func (s *APIGatewayService) UpdateApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiKeyId := request.GetStringParam(req.Parameters, "apiKey")
	if apiKeyId == "" {
		apiKeyId = getPathParam(req, "apiKey")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	apiKey, err := s.updateApiKeyCore(stores, apiKeyId, ops)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toApiKeyResponse(apiKey), nil
}

// GetApiKeys retrieves all API keys from API Gateway.
func (s *APIGatewayService) GetApiKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listApiKeysCore(stores, limit, request.GetStringParam(req.Parameters, "position"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, k := range result.Items {
		items = append(items, s.toApiKeyResponse(k))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toApiKeyResponse(k *apigateway.ApiKey) map[string]interface{} {
	return s.toApiKeyResponseWithIncludeValue(k, false)
}

func (s *APIGatewayService) toApiKeyResponseWithIncludeValue(k *apigateway.ApiKey, includeValue bool) map[string]interface{} {
	response := map[string]interface{}{
		"id":              k.Id,
		"name":            k.Name,
		"enabled":         k.Enabled,
		"createdDate":     timeutils.FormatEpochSeconds(k.CreatedDate),
		"lastUpdatedDate": timeutils.FormatEpochSeconds(k.LastUpdatedDate),
	}

	if includeValue && k.Value != "" {
		response["value"] = k.Value
	}
	if k.Description != "" {
		response["description"] = k.Description
	}
	if len(k.StageKeys) > 0 {
		response["stageKeys"] = k.StageKeys
	}
	if k.CustomerId != "" {
		response["customerId"] = k.CustomerId
	}
	if len(k.Tags) > 0 {
		response["tags"] = tagutil.ToMap(k.Tags)
	}

	return response
}

// CreateUsagePlan creates a new usage plan in API Gateway.
func (s *APIGatewayService) CreateUsagePlan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &UsagePlanInput{
		Name:        request.GetStringParam(req.Parameters, "name"),
		Description: request.GetStringParam(req.Parameters, "description"),
	}
	if apiStages, ok := req.Parameters["apiStages"].([]interface{}); ok {
		for _, as := range apiStages {
			if asMap, ok := as.(map[string]interface{}); ok {
				stage := ApiStageInput{}
				if apiId, ok := asMap["apiId"].(string); ok {
					stage.ApiId = apiId
				}
				if stageName, ok := asMap["stage"].(string); ok {
					stage.Stage = stageName
				}
				if throttleMap, ok := asMap["throttle"].(map[string]interface{}); ok {
					stage.Throttle = make(map[string]*apigateway.Throttle)
					for k, v := range throttleMap {
						if ts, ok := v.(map[string]interface{}); ok {
							t := &apigateway.Throttle{}
							if bl, ok := ts["burstLimit"]; ok {
								switch bv := bl.(type) {
								case int:
									t.BurstLimit = int64(bv)
								case float64:
									t.BurstLimit = int64(bv)
								}
							}
							if rl, ok := ts["rateLimit"]; ok {
								if rv, ok := rl.(float64); ok {
									t.RateLimit = rv
								}
							}
							stage.Throttle[k] = t
						}
					}
				}
				in.ApiStages = append(in.ApiStages, stage)
			}
		}
	}
	if quotaMap, ok := req.Parameters["quota"].(map[string]interface{}); ok {
		in.Quota = &QuotaInput{}
		if limit, ok := quotaMap["limit"]; ok {
			switch v := limit.(type) {
			case int:
				in.Quota.Limit = int64(v)
			case float64:
				in.Quota.Limit = int64(v)
			}
		}
		if period, ok := quotaMap["period"].(string); ok {
			in.Quota.Period = period
		}
		if offset, ok := quotaMap["offset"]; ok {
			switch v := offset.(type) {
			case int:
				in.Quota.Offset = int64(v)
			case float64:
				in.Quota.Offset = int64(v)
			}
		}
	}
	if throttleMap, ok := req.Parameters["throttle"].(map[string]interface{}); ok {
		in.Throttle = &ThrottleInput{}
		if burst, ok := throttleMap["burstLimit"]; ok {
			switch v := burst.(type) {
			case int:
				in.Throttle.BurstLimit = int64(v)
			case float64:
				in.Throttle.BurstLimit = int64(v)
			}
		}
		if rate, ok := throttleMap["rateLimit"]; ok {
			if v, ok := rate.(float64); ok {
				in.Throttle.RateLimit = v
			}
		}
	}
	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		in.Tags = tagutil.MapInterfaceToTags(tags)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createUsagePlanCore(stores, in)
	if err != nil {
		return nil, err
	}
	return s.toUsagePlanResponse(created), nil
}

// GetUsagePlan retrieves a usage plan from API Gateway.
func (s *APIGatewayService) GetUsagePlan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	plan, err := s.getUsagePlanCore(stores, usagePlanId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toUsagePlanResponse(plan), nil
}

// DeleteUsagePlan deletes a usage plan from API Gateway.
func (s *APIGatewayService) DeleteUsagePlan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteUsagePlanCore(stores, usagePlanId); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// UpdateUsagePlan updates an existing usage plan in API Gateway.
func (s *APIGatewayService) UpdateUsagePlan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	plan, err := s.updateUsagePlanCore(stores, usagePlanId, ops)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toUsagePlanResponse(plan), nil
}

// GetUsagePlans retrieves all usage plans from API Gateway.
func (s *APIGatewayService) GetUsagePlans(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listUsagePlansCore(stores, limit, request.GetStringParam(req.Parameters, "position"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, s.toUsagePlanResponse(p))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toUsagePlanResponse(p *apigateway.UsagePlan) map[string]interface{} {
	response := map[string]interface{}{
		"id":   p.Id,
		"name": p.Name,
	}

	if p.Description != "" {
		response["description"] = p.Description
	}
	if len(p.ApiStages) > 0 {
		apiStages := make([]interface{}, 0, len(p.ApiStages))
		for _, as := range p.ApiStages {
			stageMap := map[string]interface{}{
				"apiId": as.ApiId,
				"stage": as.Stage,
			}
			if len(as.Throttle) > 0 {
				throttleMap := make(map[string]interface{}, len(as.Throttle))
				for k, v := range as.Throttle {
					throttleMap[k] = map[string]interface{}{
						"burstLimit": v.BurstLimit,
						"rateLimit":  v.RateLimit,
					}
				}
				stageMap["throttle"] = throttleMap
			}
			apiStages = append(apiStages, stageMap)
		}
		response["apiStages"] = apiStages
	}
	if p.Quota != nil {
		response["quota"] = map[string]interface{}{
			"limit":  p.Quota.Limit,
			"offset": p.Quota.Offset,
			"period": p.Quota.Period,
		}
	}
	if p.Throttle != nil {
		response["throttle"] = map[string]interface{}{
			"burstLimit": p.Throttle.BurstLimit,
			"rateLimit":  p.Throttle.RateLimit,
		}
	}
	if p.ProductCode != "" {
		response["productCode"] = p.ProductCode
	}
	if len(p.Tags) > 0 {
		response["tags"] = tagutil.ToMap(p.Tags)
	}

	return response
}

// CreateUsagePlanKey associates an API key with a usage plan in API Gateway.
func (s *APIGatewayService) CreateUsagePlanKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	in := &UsagePlanKeyInput{
		KeyId:   request.GetStringParam(req.Parameters, "keyId"),
		KeyType: request.GetStringParam(req.Parameters, "keyType"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createUsagePlanKeyCore(stores, usagePlanId, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toUsagePlanKeyResponse(created), nil
}

// GetUsagePlanKey retrieves a usage plan key from API Gateway.
func (s *APIGatewayService) GetUsagePlanKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	keyId := request.GetStringParam(req.Parameters, "keyId")
	if keyId == "" {
		keyId = getPathParam(req, "keyId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.getUsagePlanKeyCore(stores, usagePlanId, keyId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toUsagePlanKeyResponse(key), nil
}

// DeleteUsagePlanKey removes an API key association from a usage plan in API Gateway.
func (s *APIGatewayService) DeleteUsagePlanKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	keyId := request.GetStringParam(req.Parameters, "keyId")
	if keyId == "" {
		keyId = getPathParam(req, "keyId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteUsagePlanKeyCore(stores, usagePlanId, keyId); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// GetUsagePlanKeys retrieves all API keys associated with a usage plan in API Gateway.
func (s *APIGatewayService) GetUsagePlanKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listUsagePlanKeysCore(stores, usagePlanId, limit, request.GetStringParam(req.Parameters, "position"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, k := range result.Items {
		items = append(items, s.toUsagePlanKeyResponse(k))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toUsagePlanKeyResponse(k *apigateway.UsagePlanKey) map[string]interface{} {
	return map[string]interface{}{
		"id":    k.Id,
		"type":  k.Type,
		"value": k.Value,
		"name":  k.Name,
	}
}

// GetUsage retrieves usage data for a usage plan in API Gateway.
func (s *APIGatewayService) GetUsage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	keyId := request.GetStringParam(req.Parameters, "keyId")
	if keyId == "" {
		keyId = getPathParam(req, "keyId")
	}
	startDate := request.GetStringParam(req.Parameters, "startDate")
	endDate := request.GetStringParam(req.Parameters, "endDate")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getUsageCore(stores, usagePlanId, keyId, startDate, endDate)
}
