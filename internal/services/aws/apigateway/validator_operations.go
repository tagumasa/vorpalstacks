// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/apigateway"
)

// CreateRequestValidator creates a new request validator in API Gateway.
func (s *APIGatewayService) CreateRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	in := &RequestValidatorInput{
		Name:                      request.GetStringParam(req.Parameters, "name"),
		ValidateRequestBody:       request.GetBoolParam(req.Parameters, "validateRequestBody"),
		ValidateRequestParameters: request.GetBoolParam(req.Parameters, "validateRequestParameters"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createRequestValidatorCore(stores, apiId, in)
	if err != nil {
		return nil, err
	}
	return s.toRequestValidatorResponse(created), nil
}

// GetRequestValidator retrieves a request validator from API Gateway.
func (s *APIGatewayService) GetRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	validatorId := request.GetStringParam(req.Parameters, "requestValidatorId")
	if validatorId == "" {
		validatorId = getPathParam(req, "requestValidatorId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	validator, err := s.getRequestValidatorCore(stores, apiId, validatorId)
	if err != nil {
		return nil, err
	}
	return s.toRequestValidatorResponse(validator), nil
}

// DeleteRequestValidator deletes a request validator from API Gateway.
func (s *APIGatewayService) DeleteRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	validatorId := request.GetStringParam(req.Parameters, "requestValidatorId")
	if validatorId == "" {
		validatorId = getPathParam(req, "requestValidatorId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteRequestValidatorCore(stores, apiId, validatorId); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UpdateRequestValidator updates an existing request validator in API Gateway.
func (s *APIGatewayService) UpdateRequestValidator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	validatorId := request.GetStringParam(req.Parameters, "requestValidatorId")
	if validatorId == "" {
		validatorId = getPathParam(req, "requestValidatorId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	validator, err := s.updateRequestValidatorCore(stores, apiId, validatorId, ops)
	if err != nil {
		return nil, err
	}
	return s.toRequestValidatorResponse(validator), nil
}

// GetRequestValidators retrieves all request validators for an API in API Gateway.
func (s *APIGatewayService) GetRequestValidators(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	validators, err := s.listRequestValidatorsCore(stores, apiId)
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

func (s *APIGatewayService) toRequestValidatorResponse(v *apigateway.RequestValidator) map[string]interface{} {
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
	in := &ModelInput{
		Name:        request.GetStringParam(req.Parameters, "name"),
		Description: request.GetStringParam(req.Parameters, "description"),
		Schema:      request.GetStringParam(req.Parameters, "schema"),
		ContentType: request.GetStringParam(req.Parameters, "contentType"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createModelCore(stores, apiId, in)
	if err != nil {
		return nil, err
	}
	return s.toModelResponse(created), nil
}

// GetModel retrieves a model from API Gateway.
func (s *APIGatewayService) GetModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	modelName := request.GetStringParam(req.Parameters, "modelName")
	if modelName == "" {
		modelName = getPathParam(req, "modelName")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	model, err := s.getModelCore(stores, apiId, modelName)
	if err != nil {
		return nil, err
	}
	return s.toModelResponse(model), nil
}

// DeleteModel deletes a model from API Gateway.
func (s *APIGatewayService) DeleteModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	modelName := request.GetStringParam(req.Parameters, "modelName")
	if modelName == "" {
		modelName = getPathParam(req, "modelName")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteModelCore(stores, apiId, modelName); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetModels retrieves all models for an API in API Gateway.
func (s *APIGatewayService) GetModels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	models, err := s.listModelsCore(stores, apiId)
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

func (s *APIGatewayService) toModelResponse(m *apigateway.Model) map[string]interface{} {
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
