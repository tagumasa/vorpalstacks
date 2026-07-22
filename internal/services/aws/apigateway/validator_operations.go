// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	store "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateRequestValidator creates a new request validator in API Gateway.
func (s *APIGatewayService) CreateRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	validator := &store.RequestValidator{
		Name:                      name,
		ValidateRequestBody:       request.GetBoolParam(req.Parameters, "validateRequestBody"),
		ValidateRequestParameters: request.GetBoolParam(req.Parameters, "validateRequestParameters"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := stores.restApis.CreateRequestValidator(apiId, validator)
	if err != nil {
		return nil, err
	}

	return s.toRequestValidatorResponse(created), nil
}

// GetRequestValidator retrieves a request validator from API Gateway.
func (s *APIGatewayService) GetRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	validatorId := request.GetStringParam(req.Parameters, "requestValidatorId")
	if validatorId == "" {
		validatorId = getPathParam(req, "requestValidatorId")
	}
	if validatorId == "" {
		return nil, NewBadRequestException("requestValidatorId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	validator, err := stores.restApis.GetRequestValidator(apiId, validatorId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toRequestValidatorResponse(validator), nil
}

// DeleteRequestValidator deletes a request validator from API Gateway.
func (s *APIGatewayService) DeleteRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	validatorId := request.GetStringParam(req.Parameters, "requestValidatorId")
	if validatorId == "" {
		validatorId = getPathParam(req, "requestValidatorId")
	}
	if validatorId == "" {
		return nil, NewBadRequestException("requestValidatorId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteRequestValidator(apiId, validatorId); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// UpdateRequestValidator updates an existing request validator in API Gateway.
func (s *APIGatewayService) UpdateRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	validatorId := request.GetStringParam(req.Parameters, "requestValidatorId")
	if validatorId == "" {
		validatorId = getPathParam(req, "requestValidatorId")
	}
	if validatorId == "" {
		return nil, NewBadRequestException("requestValidatorId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stores.keyLocker.Lock(apiId + ":" + validatorId)
	defer stores.keyLocker.Unlock(apiId + ":" + validatorId)

	validator, err := stores.restApis.GetRequestValidator(apiId, validatorId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/name":
			validator.Name = po.Value
		case "/validateRequestBody":
			validator.ValidateRequestBody = po.Value == "true"
		case "/validateRequestParameters":
			validator.ValidateRequestParameters = po.Value == "true"
		}
	}

	if err := stores.restApis.UpdateRequestValidator(apiId, validator); err != nil {
		return nil, err
	}

	return s.toRequestValidatorResponse(validator), nil
}

// GetRequestValidators retrieves all request validators for an API in API Gateway.
func (s *APIGatewayService) GetRequestValidators(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	limit := request.GetIntParam(req.Parameters, "limit")
	if limit <= 0 {
		limit = 25
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	validators, err := stores.restApis.ListRequestValidators(apiId)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(validators))
	for _, v := range validators {
		items = append(items, s.toRequestValidatorResponse(v))
	}

	page, nextPos := paginateItems(items, position, limit)
	result := map[string]interface{}{
		"item": page,
	}
	if nextPos != "" {
		result["position"] = nextPos
	}
	return result, nil
}

func (s *APIGatewayService) toRequestValidatorResponse(v *store.RequestValidator) map[string]interface{} {
	return map[string]interface{}{
		"id":                        v.Id,
		"name":                      v.Name,
		"validateRequestBody":       v.ValidateRequestBody,
		"validateRequestParameters": v.ValidateRequestParameters,
	}
}

// CreateModel creates a new model in API Gateway.
func (s *APIGatewayService) CreateModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	model := &store.Model{
		Name:        name,
		Description: request.GetStringParam(req.Parameters, "description"),
		Schema:      request.GetStringParam(req.Parameters, "schema"),
		ContentType: request.GetStringParam(req.Parameters, "contentType"),
	}
	if model.ContentType == "" {
		model.ContentType = "application/json"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.restApis.CreateModel(apiId, model)
	if err != nil {
		return nil, err
	}

	return s.toModelResponse(created), nil
}

// GetModel retrieves a model from API Gateway.
func (s *APIGatewayService) GetModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	modelName := request.GetStringParam(req.Parameters, "modelName")
	if modelName == "" {
		modelName = getPathParam(req, "modelName")
	}
	if modelName == "" {
		return nil, NewBadRequestException("modelName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	model, err := stores.restApis.GetModel(apiId, modelName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toModelResponse(model), nil
}

// DeleteModel deletes a model from API Gateway.
func (s *APIGatewayService) DeleteModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	modelName := request.GetStringParam(req.Parameters, "modelName")
	if modelName == "" {
		modelName = getPathParam(req, "modelName")
	}
	if modelName == "" {
		return nil, NewBadRequestException("modelName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteModel(apiId, modelName); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// GetModels retrieves all models for an API in API Gateway.
func (s *APIGatewayService) GetModels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	limit := request.GetIntParam(req.Parameters, "limit")
	if limit <= 0 {
		limit = 25
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	models, err := stores.restApis.ListModels(apiId)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(models))
	for _, m := range models {
		items = append(items, s.toModelResponse(m))
	}

	page, nextPos := paginateItems(items, position, limit)
	result := map[string]interface{}{
		"item": page,
	}
	if nextPos != "" {
		result["position"] = nextPos
	}
	return result, nil
}

func (s *APIGatewayService) toModelResponse(m *store.Model) map[string]interface{} {
	response := map[string]interface{}{
		"id":     m.Id,
		"name":   m.Name,
		"schema": m.Schema,
	}

	if m.Description != "" {
		response["description"] = m.Description
	}
	if m.ContentType != "" {
		response["contentType"] = m.ContentType
	}

	return response
}

// TagResource adds tags to a resource in API Gateway.
func (s *APIGatewayService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arnStr := request.GetStringParam(req.Parameters, "resourceArn")

	tagsMap := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags"))
	if len(tagsMap) == 0 {
		return nil, NewBadRequestException("tags is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.tagResource(stores, arnStr, tagsMap); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a resource in API Gateway.
func (s *APIGatewayService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arnStr := request.GetStringParam(req.Parameters, "resourceArn")

	tagKeys := tagutil.ParseTagKeysWithQueryFallback(req.Parameters, "tagKeys")
	if len(tagKeys) == 0 {
		return nil, NewBadRequestException("tagKeys is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.untagResource(stores, arnStr, tagKeys); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for a resource in API Gateway.
func (s *APIGatewayService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arnStr := request.GetStringParam(req.Parameters, "resourceArn")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tagsList, err := s.getResourceTags(stores, arnStr)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tagutil.ToMap(tagsList),
	}, nil
}

// tagResource dispatches tag operations based on the resource ARN pattern.
func (s *APIGatewayService) tagResource(stores *apiGatewayStores, arnStr string, tagsMap map[string]string) error {
	switch {
	case strings.Contains(arnStr, "/stages/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		stageName := extractResourceFromArn(arnStr, "/stages/")
		if apiId == "" || stageName == "" {
			return NewBadRequestException("invalid stage ARN")
		}
		return stores.restApis.TagStage(apiId, stageName, tagsMap)

	case strings.Contains(arnStr, "/usageplans/"):
		usagePlanId := extractResourceFromArn(arnStr, "/usageplans/")
		if usagePlanId == "" {
			return NewBadRequestException("invalid usage plan ARN")
		}
		return stores.usage.TagUsagePlan(usagePlanId, tagsMap)

	case strings.Contains(arnStr, "/apikeys/"):
		apiKeyId := extractResourceFromArn(arnStr, "/apikeys/")
		if apiKeyId == "" {
			return NewBadRequestException("invalid API key ARN")
		}
		return stores.usage.TagApiKey(apiKeyId, tagsMap)

	case strings.Contains(arnStr, "/domainnames/"):
		domainName := extractResourceFromArn(arnStr, "/domainnames/")
		if domainName == "" {
			return NewBadRequestException("invalid domain name ARN")
		}
		return stores.domains.TagDomainName(domainName, tagsMap)

	case strings.Contains(arnStr, "/restapis/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		if apiId == "" {
			return NewBadRequestException("invalid REST API ARN")
		}
		return stores.restApis.Tag(apiId, tagsMap)

	default:
		return NewBadRequestException("resourceArn is required")
	}
}

// untagResource dispatches untag operations based on the resource ARN pattern.
func (s *APIGatewayService) untagResource(stores *apiGatewayStores, arnStr string, tagKeys []string) error {
	switch {
	case strings.Contains(arnStr, "/stages/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		stageName := extractResourceFromArn(arnStr, "/stages/")
		if apiId == "" || stageName == "" {
			return NewBadRequestException("invalid stage ARN")
		}
		return stores.restApis.UntagStage(apiId, stageName, tagKeys)

	case strings.Contains(arnStr, "/usageplans/"):
		usagePlanId := extractResourceFromArn(arnStr, "/usageplans/")
		if usagePlanId == "" {
			return NewBadRequestException("invalid usage plan ARN")
		}
		return stores.usage.UntagUsagePlan(usagePlanId, tagKeys)

	case strings.Contains(arnStr, "/apikeys/"):
		apiKeyId := extractResourceFromArn(arnStr, "/apikeys/")
		if apiKeyId == "" {
			return NewBadRequestException("invalid API key ARN")
		}
		return stores.usage.UntagApiKey(apiKeyId, tagKeys)

	case strings.Contains(arnStr, "/domainnames/"):
		domainName := extractResourceFromArn(arnStr, "/domainnames/")
		if domainName == "" {
			return NewBadRequestException("invalid domain name ARN")
		}
		return stores.domains.UntagDomainName(domainName, tagKeys)

	case strings.Contains(arnStr, "/restapis/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		if apiId == "" {
			return NewBadRequestException("invalid REST API ARN")
		}
		return stores.restApis.Untag(apiId, tagKeys)

	default:
		return NewBadRequestException("resourceArn is required")
	}
}

// getResourceTags dispatches get-tags operations based on the resource ARN pattern.
func (s *APIGatewayService) getResourceTags(stores *apiGatewayStores, arnStr string) ([]types.Tag, error) {
	switch {
	case strings.Contains(arnStr, "/stages/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		stageName := extractResourceFromArn(arnStr, "/stages/")
		if apiId == "" || stageName == "" {
			return nil, NewBadRequestException("invalid stage ARN")
		}
		return stores.restApis.GetStageTags(apiId, stageName)

	case strings.Contains(arnStr, "/usageplans/"):
		usagePlanId := extractResourceFromArn(arnStr, "/usageplans/")
		if usagePlanId == "" {
			return nil, NewBadRequestException("invalid usage plan ARN")
		}
		return stores.usage.GetUsagePlanTags(usagePlanId)

	case strings.Contains(arnStr, "/apikeys/"):
		apiKeyId := extractResourceFromArn(arnStr, "/apikeys/")
		if apiKeyId == "" {
			return nil, NewBadRequestException("invalid API key ARN")
		}
		return stores.usage.GetApiKeyTags(apiKeyId)

	case strings.Contains(arnStr, "/domainnames/"):
		domainName := extractResourceFromArn(arnStr, "/domainnames/")
		if domainName == "" {
			return nil, NewBadRequestException("invalid domain name ARN")
		}
		return stores.domains.GetDomainNameTags(domainName)

	case strings.Contains(arnStr, "/restapis/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		if apiId == "" {
			return nil, NewBadRequestException("invalid REST API ARN")
		}
		return stores.restApis.GetResourceTags(apiId)

	default:
		return nil, NewBadRequestException("resourceArn is required")
	}
}

func extractResourceFromArn(arnStr, suffix string) string {
	idx := strings.Index(arnStr, suffix)
	if idx < 0 {
		return ""
	}
	rest := arnStr[idx+len(suffix):]
	if slashIdx := strings.Index(rest, "/"); slashIdx >= 0 {
		return rest[:slashIdx]
	}
	return rest
}
