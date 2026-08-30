// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// UpdateResource updates an existing resource in API Gateway.
func (s *APIGatewayService) UpdateResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	resourceId := getResourceId(req)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	resource, err := s.updateResourceCore(stores, apiId, resourceId, ops)
	if err != nil {
		return nil, err
	}
	return s.toResourceResponse(resource), nil
}

// UpdateMethod updates an existing method in API Gateway.
func (s *APIGatewayService) UpdateMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	resourceId := getResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	method, err := s.updateMethodCore(stores, apiId, resourceId, httpMethod, ops)
	if err != nil {
		return nil, err
	}
	return s.toMethodResponse(method), nil
}

// UpdateIntegration updates an existing integration in API Gateway.
func (s *APIGatewayService) UpdateIntegration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	resourceId := getResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	integration, err := s.updateIntegrationCore(stores, apiId, resourceId, httpMethod, ops)
	if err != nil {
		return nil, err
	}
	return s.toIntegrationResponse(integration), nil
}

// UpdateIntegrationResponse updates an existing integration response in API Gateway.
func (s *APIGatewayService) UpdateIntegrationResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	resourceId := getResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	statusCode := request.GetStringParam(req.Parameters, "statusCode")
	if statusCode == "" {
		statusCode = getPathParam(req, "statusCode")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	intResp, err := s.updateIntegrationResponseCore(stores, apiId, resourceId, httpMethod, statusCode, ops)
	if err != nil {
		return nil, err
	}
	return s.toIntegrationResponseResponse(intResp), nil
}

// UpdateDeployment updates an existing deployment in API Gateway.
func (s *APIGatewayService) UpdateDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	deploymentId := request.GetStringParam(req.Parameters, "deploymentId")
	if deploymentId == "" {
		deploymentId = getPathParam(req, "deploymentId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	deployment, err := s.updateDeploymentCore(stores, apiId, deploymentId, ops)
	if err != nil {
		return nil, err
	}
	return s.toDeploymentResponse(deployment), nil
}

// UpdateModel updates an existing model in API Gateway.
func (s *APIGatewayService) UpdateModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	modelName := request.GetStringParam(req.Parameters, "modelName")
	if modelName == "" {
		modelName = getPathParam(req, "modelName")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	model, err := s.updateModelCore(stores, apiId, modelName, ops)
	if err != nil {
		return nil, err
	}
	return s.toModelResponse(model), nil
}
