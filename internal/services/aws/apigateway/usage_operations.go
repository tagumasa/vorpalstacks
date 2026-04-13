// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"strings"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	store "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateApiKey creates a new API key in API Gateway.
func (s *APIGatewayService) CreateApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	apiKey := &store.ApiKey{
		Name:        name,
		Description: request.GetStringParam(req.Parameters, "description"),
		Enabled:     request.GetBoolParam(req.Parameters, "enabled"),
		CustomerId:  request.GetStringParam(req.Parameters, "customerId"),
	}

	if stageKeys, ok := req.Parameters["stageKeys"].([]interface{}); ok {
		for _, sk := range stageKeys {
			if sks, ok := sk.(string); ok {
				apiKey.StageKeys = append(apiKey.StageKeys, sks)
			}
		}
	}

	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		apiKey.Tags = tagutil.MapInterfaceToTags(tags)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.usage.CreateApiKey(apiKey)
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
	if apiKeyId == "" {
		return nil, NewBadRequestException("apiKey is required")
	}

	includeValue := request.GetBoolParam(req.Parameters, "includeValue")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	apiKey, err := stores.usage.GetApiKey(apiKeyId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toApiKeyResponseWithIncludeValue(apiKey, includeValue), nil
}

// DeleteApiKey deletes an API key from API Gateway.
func (s *APIGatewayService) DeleteApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiKeyId := request.GetStringParam(req.Parameters, "apiKey")
	if apiKeyId == "" {
		apiKeyId = getPathParam(req, "apiKey")
	}
	if apiKeyId == "" {
		return nil, NewBadRequestException("apiKey is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.usage.DeleteApiKey(apiKeyId); err != nil {
		return nil, ErrNotFoundException
	}

	if s.runtimeServer != nil {
		s.runtimeServer.RemoveApiKey(apiKeyId)
	}

	return response.EmptyResponse(), nil
}

// UpdateApiKey updates an existing API key in API Gateway.
func (s *APIGatewayService) UpdateApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiKeyId := request.GetStringParam(req.Parameters, "apiKey")
	if apiKeyId == "" {
		apiKeyId = getPathParam(req, "apiKey")
	}
	if apiKeyId == "" {
		return nil, NewBadRequestException("apiKey is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	apiKey, err := stores.usage.GetApiKey(apiKeyId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	patchOps, ok := req.Parameters["patchOperations"].([]interface{})
	if ok {
		for _, op := range patchOps {
			if opMap, ok := op.(map[string]interface{}); ok {
				path := ""
				value := ""
				if p, ok := opMap["path"].(string); ok {
					path = p
				}
				if v, ok := opMap["value"].(string); ok {
					value = v
				}

				switch path {
				case "/name":
					apiKey.Name = value
				case "/description":
					apiKey.Description = value
				case "/enabled":
					apiKey.Enabled = value == "true"
				}
				if strings.HasPrefix(path, "/stageKeys/") {
					if apiKey.StageKeys == nil {
						apiKey.StageKeys = []string{}
					}
					stageKey := strings.TrimPrefix(path, "/stageKeys/")
					if opMap["op"] == "add" || opMap["op"] == "replace" {
						if !containsString(apiKey.StageKeys, stageKey) {
							apiKey.StageKeys = append(apiKey.StageKeys, stageKey)
						}
					} else if opMap["op"] == "remove" {
						newStageKeys := []string{}
						for _, sk := range apiKey.StageKeys {
							if sk != stageKey {
								newStageKeys = append(newStageKeys, sk)
							}
						}
						apiKey.StageKeys = newStageKeys
					}
				}
			}
		}
	}

	if err := stores.usage.UpdateApiKey(apiKey); err != nil {
		return nil, err
	}

	return s.toApiKeyResponse(apiKey), nil
}

// GetApiKeys retrieves all API keys from API Gateway.
func (s *APIGatewayService) GetApiKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := request.GetIntParam(req.Parameters, "limit")
	if maxItems <= 0 {
		maxItems = 25
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.usage.ListApiKeys(common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
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

func (s *APIGatewayService) toApiKeyResponse(k *store.ApiKey) map[string]interface{} {
	return s.toApiKeyResponseWithIncludeValue(k, false)
}

func (s *APIGatewayService) toApiKeyResponseWithIncludeValue(k *store.ApiKey, includeValue bool) map[string]interface{} {
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
	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	usagePlan := &store.UsagePlan{
		Name:        name,
		Description: request.GetStringParam(req.Parameters, "description"),
		ProductCode: request.GetStringParam(req.Parameters, "productCode"),
	}

	if apiStages, ok := req.Parameters["apiStages"].([]interface{}); ok {
		for _, as := range apiStages {
			if asMap, ok := as.(map[string]interface{}); ok {
				stage := store.ApiStage{}
				if apiId, ok := asMap["apiId"].(string); ok {
					stage.ApiId = apiId
				}
				if stageName, ok := asMap["stage"].(string); ok {
					stage.Stage = stageName
				}
				usagePlan.ApiStages = append(usagePlan.ApiStages, stage)
			}
		}
	}

	if quotaMap, ok := req.Parameters["quota"].(map[string]interface{}); ok {
		usagePlan.Quota = &store.Quota{}
		if limit, ok := quotaMap["limit"]; ok {
			switch v := limit.(type) {
			case int:
				usagePlan.Quota.Limit = int64(v)
			case float64:
				usagePlan.Quota.Limit = int64(v)
			}
		}
		if period, ok := quotaMap["period"].(string); ok {
			usagePlan.Quota.Period = period
		}
	}

	if throttleMap, ok := req.Parameters["throttle"].(map[string]interface{}); ok {
		usagePlan.Throttle = &store.Throttle{}
		if burst, ok := throttleMap["burstLimit"]; ok {
			switch v := burst.(type) {
			case int:
				usagePlan.Throttle.BurstLimit = int64(v)
			case float64:
				usagePlan.Throttle.BurstLimit = int64(v)
			}
		}
		if rate, ok := throttleMap["rateLimit"]; ok {
			switch v := rate.(type) {
			case float64:
				usagePlan.Throttle.RateLimit = v
			}
		}
	}

	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		usagePlan.Tags = tagutil.MapInterfaceToTags(tags)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.usage.CreateUsagePlan(usagePlan)
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
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	usagePlan, err := stores.usage.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toUsagePlanResponse(usagePlan), nil
}

// DeleteUsagePlan deletes a usage plan from API Gateway.
func (s *APIGatewayService) DeleteUsagePlan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.usage.DeleteUsagePlan(usagePlanId); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// UpdateUsagePlan updates an existing usage plan in API Gateway.
func (s *APIGatewayService) UpdateUsagePlan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	usagePlan, err := stores.usage.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	patchOps, ok := req.Parameters["patchOperations"].([]interface{})
	if ok {
		for _, op := range patchOps {
			if opMap, ok := op.(map[string]interface{}); ok {
				path := ""
				value := ""
				if p, ok := opMap["path"].(string); ok {
					path = p
				}
				if v, ok := opMap["value"].(string); ok {
					value = v
				}

				switch path {
				case "/name":
					usagePlan.Name = value
				case "/description":
					usagePlan.Description = value
				case "/quota/limit":
					if usagePlan.Quota == nil {
						usagePlan.Quota = &store.Quota{}
					}
					usagePlan.Quota.Limit = parseInt64(value)
				case "/quota/period":
					if usagePlan.Quota == nil {
						usagePlan.Quota = &store.Quota{}
					}
					usagePlan.Quota.Period = value
				case "/quota/offset":
					if usagePlan.Quota == nil {
						usagePlan.Quota = &store.Quota{}
					}
					usagePlan.Quota.Offset = parseInt64(value)
				case "/throttle/burstLimit":
					if usagePlan.Throttle == nil {
						usagePlan.Throttle = &store.Throttle{}
					}
					usagePlan.Throttle.BurstLimit = parseInt64(value)
				case "/throttle/rateLimit":
					if usagePlan.Throttle == nil {
						usagePlan.Throttle = &store.Throttle{}
					}
					usagePlan.Throttle.RateLimit = parseFloat64(value)
				}
			}
		}
	}

	if err := stores.usage.UpdateUsagePlan(usagePlan); err != nil {
		return nil, err
	}

	return s.toUsagePlanResponse(usagePlan), nil
}

// GetUsagePlans retrieves all usage plans from API Gateway.
func (s *APIGatewayService) GetUsagePlans(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := request.GetIntParam(req.Parameters, "limit")
	if maxItems <= 0 {
		maxItems = 25
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.usage.ListUsagePlans(common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
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

func (s *APIGatewayService) toUsagePlanResponse(p *store.UsagePlan) map[string]interface{} {
	response := map[string]interface{}{
		"id":   p.Id,
		"name": p.Name,
	}

	if p.Description != "" {
		response["description"] = p.Description
	}
	if len(p.ApiStages) > 0 {
		response["apiStages"] = p.ApiStages
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
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	keyId := request.GetStringParam(req.Parameters, "keyId")
	keyType := request.GetStringParam(req.Parameters, "keyType")
	if keyId == "" {
		return nil, NewBadRequestException("keyId is required")
	}
	if keyType == "" {
		return nil, NewBadRequestException("keyType is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiKey, err := stores.usage.GetApiKey(keyId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	key := &store.UsagePlanKey{
		Id:    keyId,
		Type:  keyType,
		Value: apiKey.Value,
		Name:  apiKey.Name,
	}

	created, err := stores.usage.CreateUsagePlanKey(usagePlanId, key)
	if err != nil {
		return nil, err
	}

	return s.toUsagePlanKeyResponse(created), nil
}

// GetUsagePlanKey retrieves a usage plan key from API Gateway.
func (s *APIGatewayService) GetUsagePlanKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	keyId := request.GetStringParam(req.Parameters, "keyId")
	if keyId == "" {
		keyId = getPathParam(req, "keyId")
	}
	if keyId == "" {
		return nil, NewBadRequestException("keyId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := stores.usage.GetUsagePlanKey(usagePlanId, keyId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toUsagePlanKeyResponse(key), nil
}

// DeleteUsagePlanKey removes an API key association from a usage plan in API Gateway.
func (s *APIGatewayService) DeleteUsagePlanKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	keyId := request.GetStringParam(req.Parameters, "keyId")
	if keyId == "" {
		keyId = getPathParam(req, "keyId")
	}
	if keyId == "" {
		return nil, NewBadRequestException("keyId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.usage.DeleteUsagePlanKey(usagePlanId, keyId); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// GetUsagePlanKeys retrieves all API keys associated with a usage plan in API Gateway.
func (s *APIGatewayService) GetUsagePlanKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	usagePlanId := request.GetStringParam(req.Parameters, "usagePlanId")
	if usagePlanId == "" {
		usagePlanId = getPathParam(req, "usagePlanId")
	}
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	maxItems := request.GetIntParam(req.Parameters, "limit")
	if maxItems <= 0 {
		maxItems = 25
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.usage.ListUsagePlanKeys(usagePlanId, common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
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

func (s *APIGatewayService) toUsagePlanKeyResponse(k *store.UsagePlanKey) map[string]interface{} {
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
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	apiKeyId := request.GetStringParam(req.Parameters, "apiKeyId")
	if apiKeyId == "" {
		apiKeyId = getPathParam(req, "apiKeyId")
	}

	startDate := request.GetStringParam(req.Parameters, "startDate")
	endDate := request.GetStringParam(req.Parameters, "endDate")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = stores.usage.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	var apiKeys []string
	if apiKeyId != "" {
		apiKeys = []string{apiKeyId}
	} else {
		apiKeysResult, err := stores.usage.ListUsagePlanKeys(usagePlanId, common.ListOptions{MaxItems: 1000})
		if err != nil {
			return nil, err
		}
		for _, key := range apiKeysResult.Items {
			apiKeys = append(apiKeys, key.Id)
		}
	}

	usageCounts := make(map[string]int64)
	for _, keyId := range apiKeys {
		records, err := stores.usage.ListUsageRecordsForAPIKey(usagePlanId, keyId, startDate, endDate)
		if err != nil {
			logs.Warn("GetUsage: failed to list usage records for key", logs.String("keyId", keyId), logs.Err(err))
			continue
		}
		for _, record := range records {
			usageCounts[record.Date] += record.RequestCount
		}
	}

	items := make([]interface{}, 0, len(usageCounts))
	for date, count := range usageCounts {
		items = append(items, map[string]interface{}{
			"date":         date,
			"requestCount": count,
		})
	}

	return map[string]interface{}{
		"usagePlanId": usagePlanId,
		"startDate":   startDate,
		"endDate":     endDate,
		"items":       items,
	}, nil
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
