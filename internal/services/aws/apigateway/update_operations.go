// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"strconv"
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

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	resource, err := stores.restApis.GetResource(apiId, resourceId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	ops := parsePatchOperations(req.Parameters)
	for _, po := range ops {
		switch po.Path {
		case "/pathPart":
			if po.Value == "" {
				return nil, NewBadRequestException("pathPart cannot be empty")
			}
			if !validatePathPart(po.Value) {
				return nil, NewBadRequestException("Invalid pathPart: malformed path parameter")
			}
			// Check for path collision with siblings under the same parent.
			siblings, err := stores.restApis.ListResources(apiId)
			if err != nil {
				return nil, err
			}
			for _, sib := range siblings {
				if sib.Id != resourceId && sib.ParentId == resource.ParentId && sib.PathPart == po.Value {
					return nil, NewConflictException("pathPart already exists under this parent")
				}
			}
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

	if err := stores.restApis.UpdateResourceCascade(apiId, resource); err != nil {
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

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	method, err := stores.restApis.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch {
		case po.Path == "/authorizationType":
			if !validateAuthorizationType(po.Value) {
				return nil, NewBadRequestException("Invalid authorization type: " + po.Value)
			}
			method.AuthorizationType = po.Value
		case po.Path == "/authorizerId":
			method.AuthorizerId = po.Value
		case po.Path == "/apiKeyRequired":
			method.ApiKeyRequired = po.Value == "true"
		case po.Path == "/requestValidatorId":
			method.RequestValidatorId = po.Value
		case po.Path == "/operationName":
			method.OperationName = po.Value
		case strings.HasPrefix(po.Path, "/requestParameters/"):
			if method.RequestParameters == nil {
				method.RequestParameters = make(map[string]bool)
			}
			paramName := strings.TrimPrefix(po.Path, "/requestParameters/")
			paramName = strings.ReplaceAll(paramName, "~1", "/")
			if po.Op == "remove" {
				delete(method.RequestParameters, paramName)
			} else {
				method.RequestParameters[paramName] = po.Value == "true"
			}
		case strings.HasPrefix(po.Path, "/requestModels/"):
			if method.RequestModels == nil {
				method.RequestModels = make(map[string]string)
			}
			modelName := strings.TrimPrefix(po.Path, "/requestModels/")
			modelName = strings.ReplaceAll(modelName, "~1", "/")
			if po.Op == "remove" {
				delete(method.RequestModels, modelName)
			} else {
				method.RequestModels[modelName] = po.Value
			}
		case strings.HasPrefix(po.Path, "/authorizationScopes"):
			if po.Op == "remove" {
				if idx, err := strconv.Atoi(strings.TrimPrefix(po.Path, "/authorizationScopes/")); err == nil && idx < len(method.AuthorizationScopes) {
					method.AuthorizationScopes = append(method.AuthorizationScopes[:idx], method.AuthorizationScopes[idx+1:]...)
				}
			} else {
				idxStr := strings.TrimPrefix(po.Path, "/authorizationScopes/")
				if idxStr == "-" || idxStr == "" {
					method.AuthorizationScopes = append(method.AuthorizationScopes, po.Value)
				} else if idx, err := strconv.Atoi(idxStr); err == nil {
					if idx >= len(method.AuthorizationScopes) {
						method.AuthorizationScopes = append(method.AuthorizationScopes, po.Value)
					} else {
						method.AuthorizationScopes = append(method.AuthorizationScopes[:idx], append([]string{po.Value}, method.AuthorizationScopes[idx:]...)...)
					}
				} else {
					method.AuthorizationScopes = append(method.AuthorizationScopes, po.Value)
				}
			}
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

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	integration, err := stores.restApis.GetIntegration(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch {
		case po.Path == "/uri":
			integration.Uri = po.Value
		case po.Path == "/type":
			if !validateIntegrationType(po.Value) {
				return nil, NewBadRequestException("Invalid integration type: " + po.Value)
			}
			integration.Type = po.Value
		case po.Path == "/httpMethod":
			if po.Value != "" && !validateHTTPMethod(po.Value) {
				return nil, NewBadRequestException("Invalid integration HTTP method: " + po.Value)
			}
			integration.IntegrationHttpMethod = po.Value
		case po.Path == "/credentials":
			integration.Credentials = po.Value
		case po.Path == "/passthroughBehavior":
			if !validatePassthroughBehavior(po.Value) {
				return nil, NewBadRequestException("Invalid passthroughBehavior: " + po.Value)
			}
			integration.PassthroughBehavior = po.Value
		case po.Path == "/contentHandling":
			if !validateContentHandling(po.Value) {
				return nil, NewBadRequestException("Invalid contentHandling: " + po.Value)
			}
			integration.ContentHandling = po.Value
		case po.Path == "/cacheNamespace":
			integration.CacheNamespace = po.Value
		case po.Path == "/connectionType":
			if !validateConnectionType(po.Value) {
				return nil, NewBadRequestException("Invalid connectionType: " + po.Value)
			}
			integration.ConnectionType = po.Value
		case po.Path == "/connectionId":
			integration.ConnectionId = po.Value
		case po.Path == "/timeoutInMillis":
			v, err := parseInt32(po.Value)
			if err != nil {
				return nil, NewBadRequestException("invalid timeoutInMillis: not a number")
			}
			if v <= 0 {
				v = 29000
			}
			if v < 50 || v > 30000 {
				return nil, NewBadRequestException("timeoutInMillis must be between 50 and 30000")
			}
			integration.TimeoutInMillis = v
		case strings.HasPrefix(po.Path, "/requestParameters/"):
			if integration.RequestParameters == nil {
				integration.RequestParameters = make(map[string]string)
			}
			paramName := strings.TrimPrefix(po.Path, "/requestParameters/")
			paramName = strings.ReplaceAll(paramName, "~1", "/")
			if po.Op == "remove" {
				delete(integration.RequestParameters, paramName)
			} else {
				integration.RequestParameters[paramName] = po.Value
			}
		case strings.HasPrefix(po.Path, "/requestTemplates/"):
			if integration.RequestTemplates == nil {
				integration.RequestTemplates = make(map[string]string)
			}
			tplName := strings.TrimPrefix(po.Path, "/requestTemplates/")
			tplName = strings.ReplaceAll(tplName, "~1", "/")
			if po.Op == "remove" {
				delete(integration.RequestTemplates, tplName)
			} else {
				integration.RequestTemplates[tplName] = po.Value
			}
		case strings.HasPrefix(po.Path, "/cacheKeyParameters"):
			if po.Op == "remove" {
				if idx, err := strconv.Atoi(strings.TrimPrefix(po.Path, "/cacheKeyParameters/")); err == nil && idx < len(integration.CacheKeyParameters) {
					integration.CacheKeyParameters = append(integration.CacheKeyParameters[:idx], integration.CacheKeyParameters[idx+1:]...)
				}
			} else {
				integration.CacheKeyParameters = append(integration.CacheKeyParameters, po.Value)
			}
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

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	response, err := stores.restApis.GetIntegrationResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch {
		case po.Path == "/selectionPattern":
			response.SelectionPattern = po.Value
		case po.Path == "/contentHandling":
			if !validateContentHandling(po.Value) {
				return nil, NewBadRequestException("Invalid contentHandling: " + po.Value)
			}
			response.ContentHandling = po.Value
		case strings.HasPrefix(po.Path, "/responseParameters/"):
			if response.ResponseParameters == nil {
				response.ResponseParameters = make(map[string]string)
			}
			paramName := strings.TrimPrefix(po.Path, "/responseParameters/")
			paramName = strings.ReplaceAll(paramName, "~1", "/")
			if po.Op == "remove" {
				delete(response.ResponseParameters, paramName)
			} else {
				response.ResponseParameters[paramName] = po.Value
			}
		case strings.HasPrefix(po.Path, "/responseTemplates/"):
			if response.ResponseTemplates == nil {
				response.ResponseTemplates = make(map[string]string)
			}
			tplName := strings.TrimPrefix(po.Path, "/responseTemplates/")
			tplName = strings.ReplaceAll(tplName, "~1", "/")
			if po.Op == "remove" {
				delete(response.ResponseTemplates, tplName)
			} else {
				response.ResponseTemplates[tplName] = po.Value
			}
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

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

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

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	model, err := stores.restApis.GetModel(apiId, modelName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/description":
			model.Description = po.Value
		case "/schema":
			if !validateModelSchemaSize(po.Value) {
				return nil, NewBadRequestException("schema must not exceed 400 KB")
			}
			model.Schema = po.Value
		}
	}

	if err := stores.restApis.UpdateModel(apiId, model); err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toModelResponse(model), nil
}
