// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	store "vorpalstacks/internal/store/aws/apigateway"
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
	if len(name) > 128 {
		return nil, NewBadRequestException("name must be 1 to 128 characters")
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

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
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

	page, nextPos, found := paginateItems(items, position, limit)
	if !found {
		return nil, NewBadRequestException("Invalid position: " + position)
	}
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
	if !validateModelName(name) {
		return nil, NewBadRequestException("name must be alphanumeric")
	}

	model := &store.Model{
		Name:        name,
		Description: request.GetStringParam(req.Parameters, "description"),
		Schema:      request.GetStringParam(req.Parameters, "schema"),
		ContentType: request.GetStringParam(req.Parameters, "contentType"),
	}
	if !validateModelSchemaSize(model.Schema) {
		return nil, NewBadRequestException("schema must not exceed 400 KB")
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

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
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

	page, nextPos, found := paginateItems(items, position, limit)
	if !found {
		return nil, NewBadRequestException("Invalid position: " + position)
	}
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
