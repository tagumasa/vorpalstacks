// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	store "vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
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
	validator, err := stores.restApis.GetRequestValidator(apiId, validatorId)
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
					validator.Name = value
				case "/validateRequestBody":
					validator.ValidateRequestBody = value == "true"
				case "/validateRequestParameters":
					validator.ValidateRequestParameters = value == "true"
				}
			}
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

	return map[string]interface{}{
		"item": items,
	}, nil
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

	return map[string]interface{}{
		"item": items,
	}, nil
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
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	arnStr := request.GetStringParam(req.Parameters, "resourceArn")

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags"))
	if len(tags) == 0 {
		return nil, NewBadRequestException("tags is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if strings.Contains(arnStr, "/stages/") {
		stageName := extractResourceFromArn(arnStr, "/stages/")
		if stageName != "" {
			if err := stores.restApis.TagStage(apiId, stageName, tags); err != nil {
				return nil, err
			}
			return response.EmptyResponse(), nil
		}
	}

	if err := stores.restApis.TagResource(apiId, tags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a resource in API Gateway.
func (s *APIGatewayService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	arnStr := request.GetStringParam(req.Parameters, "resourceArn")

	tagKeys := tagutil.ParseTagKeysWithQueryFallback(req.Parameters, "tagKeys")
	if len(tagKeys) == 0 {
		return nil, NewBadRequestException("tagKeys is required")
	}

	keys := tagKeys

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if strings.Contains(arnStr, "/stages/") {
		stageName := extractResourceFromArn(arnStr, "/stages/")
		if stageName != "" {
			if err := stores.restApis.UntagStage(apiId, stageName, keys); err != nil {
				return nil, err
			}
			return response.EmptyResponse(), nil
		}
	}

	if err := stores.restApis.UntagResource(apiId, keys); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for a resource in API Gateway.
func (s *APIGatewayService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	arnStr := request.GetStringParam(req.Parameters, "resourceArn")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var tags []storecommon.Tag
	if strings.Contains(arnStr, "/stages/") {
		stageName := extractResourceFromArn(arnStr, "/stages/")
		if stageName != "" {
			tags, err = stores.restApis.GetStageTags(apiId, stageName)
		}
	}
	if tags == nil {
		tags, err = stores.restApis.GetResourceTags(apiId)
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tagutil.ToMap(tags),
	}, nil
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
