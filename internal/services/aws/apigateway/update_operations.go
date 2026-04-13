// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"strings"
	"vorpalstacks/internal/common/request"
)

// UpdateResource updates an existing resource in API Gateway.
func (s *APIGatewayService) UpdateResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	resourceId := getResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resource, err := stores.restApis.GetResource(apiId, resourceId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	ops := parsePatchOperations(req.Parameters)
	for _, po := range ops {
		switch po.Path {
		case "/pathPart":
			resource.PathPart = po.Value
			if resource.ParentId == "" {
				resource.Path = "/" + po.Value
			} else {
				parent, err := stores.restApis.GetResource(apiId, resource.ParentId)
				if err == nil {
					parentPath := strings.TrimRight(parent.Path, "/")
					resource.Path = parentPath + "/" + po.Value
				}
			}
		}
	}

	if err := stores.restApis.UpdateResource(apiId, resource); err != nil {
		return nil, toApiGatewayError(err)
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

	if apiId == "" || resourceId == "" || httpMethod == "" {
		return nil, NewBadRequestException("restApiId, resourceId, and httpMethod are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	method, err := stores.restApis.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/authorizationType":
			method.AuthorizationType = po.Value
		case "/authorizerId":
			method.AuthorizerId = po.Value
		case "/apiKeyRequired":
			method.ApiKeyRequired = po.Value == "true"
		case "/requestValidatorId":
			method.RequestValidatorId = po.Value
		case "/operationName":
			method.OperationName = po.Value
		}
	}

	_, err = stores.restApis.PutMethod(apiId, resourceId, method)
	if err != nil {
		return nil, toApiGatewayError(err)
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

	if apiId == "" || resourceId == "" || httpMethod == "" {
		return nil, NewBadRequestException("restApiId, resourceId, and httpMethod are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	integration, err := stores.restApis.GetIntegration(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/uri":
			integration.Uri = po.Value
		case "/type":
			integration.Type = po.Value
		case "/httpMethod":
			integration.IntegrationHttpMethod = po.Value
		case "/credentials":
			integration.Credentials = po.Value
		case "/passthroughBehavior":
			integration.PassthroughBehavior = po.Value
		case "/contentHandling":
			integration.ContentHandling = po.Value
		case "/cacheNamespace":
			integration.CacheNamespace = po.Value
		case "/connectionType":
			integration.ConnectionType = po.Value
		case "/connectionId":
			integration.ConnectionId = po.Value
		case "/timeoutInMillis":
			integration.TimeoutInMillis = parseInt32(po.Value)
		}
	}

	if err := stores.restApis.UpdateIntegration(apiId, resourceId, httpMethod, integration); err != nil {
		return nil, toApiGatewayError(err)
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

	if apiId == "" || resourceId == "" || httpMethod == "" || statusCode == "" {
		return nil, NewBadRequestException("missing required parameters")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	response, err := stores.restApis.GetIntegrationResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/selectionPattern":
			response.SelectionPattern = po.Value
		case "/contentHandling":
			response.ContentHandling = po.Value
		}
	}

	if err := stores.restApis.UpdateIntegrationResponse(apiId, resourceId, httpMethod, statusCode, response); err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toIntegrationResponseResponse(response), nil
}

// UpdateDeployment updates an existing deployment in API Gateway.
func (s *APIGatewayService) UpdateDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	deploymentId := request.GetStringParam(req.Parameters, "deploymentId")
	if deploymentId == "" {
		deploymentId = getPathParam(req, "deploymentId")
	}

	if apiId == "" || deploymentId == "" {
		return nil, NewBadRequestException("restApiId and deploymentId are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	deployment, err := stores.restApis.GetDeployment(apiId, deploymentId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/description":
			deployment.Description = po.Value
		}
	}

	if err := stores.restApis.UpdateDeployment(apiId, deployment); err != nil {
		return nil, toApiGatewayError(err)
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

	if apiId == "" || modelName == "" {
		return nil, NewBadRequestException("restApiId and modelName are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	model, err := stores.restApis.GetModel(apiId, modelName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/description":
			model.Description = po.Value
		case "/schema":
			model.Schema = po.Value
		}
	}

	if err := stores.restApis.UpdateModel(apiId, model); err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toModelResponse(model), nil
}
