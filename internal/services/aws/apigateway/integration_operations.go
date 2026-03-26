// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	store "vorpalstacks/internal/store/aws/apigateway"
)

// PutIntegration creates or updates an integration for a resource in API Gateway.
func (s *APIGatewayService) PutIntegration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	integrationType := request.GetStringParam(req.Parameters, "type")
	if integrationType == "" {
		return nil, NewBadRequestException("type is required")
	}

	integration := &store.Integration{
		Type:                  integrationType,
		IntegrationHttpMethod: request.GetStringParam(req.Parameters, "integrationHttpMethod"),
		Uri:                   request.GetStringParam(req.Parameters, "uri"),
		Credentials:           request.GetStringParam(req.Parameters, "credentials"),
		PassthroughBehavior:   request.GetStringParam(req.Parameters, "passthroughBehavior"),
		ContentHandling:       request.GetStringParam(req.Parameters, "contentHandling"),
		CacheNamespace:        request.GetStringParam(req.Parameters, "cacheNamespace"),
		TimeoutInMillis:       int32(request.GetIntParam(req.Parameters, "timeoutInMillis")),
		ConnectionType:        request.GetStringParam(req.Parameters, "connectionType"),
		ConnectionId:          request.GetStringParam(req.Parameters, "connectionId"),
	}

	if reqParams, ok := req.Parameters["requestParameters"].(map[string]interface{}); ok {
		integration.RequestParameters = make(map[string]string)
		for k, v := range reqParams {
			if vs, ok := v.(string); ok {
				integration.RequestParameters[k] = vs
			}
		}
	}

	if reqTemplates, ok := req.Parameters["requestTemplates"].(map[string]interface{}); ok {
		integration.RequestTemplates = make(map[string]string)
		for k, v := range reqTemplates {
			if vs, ok := v.(string); ok {
				integration.RequestTemplates[k] = vs
			}
		}
	}

	if cacheKeyParams, ok := req.Parameters["cacheKeyParameters"].([]interface{}); ok {
		for _, p := range cacheKeyParams {
			if ps, ok := p.(string); ok {
				integration.CacheKeyParameters = append(integration.CacheKeyParameters, ps)
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.restApis.PutIntegration(apiId, resourceId, httpMethod, integration)
	if err != nil {
		return nil, err
	}

	return s.toIntegrationResponse(created), nil
}

// GetIntegration retrieves an integration for a resource in API Gateway.
func (s *APIGatewayService) GetIntegration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	integration, err := stores.restApis.GetIntegration(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toIntegrationResponse(integration), nil
}

// DeleteIntegration deletes an integration from a resource in API Gateway.
func (s *APIGatewayService) DeleteIntegration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteIntegration(apiId, resourceId, httpMethod); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// PutIntegrationResponse creates or updates an integration response for a method in API Gateway.
func (s *APIGatewayService) PutIntegrationResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	statusCode := request.GetStringParam(req.Parameters, "statusCode")
	if statusCode == "" {
		return nil, NewBadRequestException("statusCode is required")
	}

	response := &store.IntegrationResponse{
		StatusCode:       statusCode,
		SelectionPattern: request.GetStringParam(req.Parameters, "selectionPattern"),
		ContentHandling:  request.GetStringParam(req.Parameters, "contentHandling"),
	}

	if respParams, ok := req.Parameters["responseParameters"].(map[string]interface{}); ok {
		response.ResponseParameters = make(map[string]string)
		for k, v := range respParams {
			if vs, ok := v.(string); ok {
				response.ResponseParameters[k] = vs
			}
		}
	}

	if respTemplates, ok := req.Parameters["responseTemplates"].(map[string]interface{}); ok {
		response.ResponseTemplates = make(map[string]string)
		for k, v := range respTemplates {
			if vs, ok := v.(string); ok {
				response.ResponseTemplates[k] = vs
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.restApis.PutIntegrationResponse(apiId, resourceId, httpMethod, statusCode, response)
	if err != nil {
		return nil, err
	}

	return s.toIntegrationResponseResponse(created), nil
}

// GetIntegrationResponse retrieves an integration response for a method in API Gateway.
func (s *APIGatewayService) GetIntegrationResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	statusCode := request.GetStringParam(req.Parameters, "statusCode")
	if statusCode == "" {
		statusCode = getPathParam(req, "statusCode")
	}
	if statusCode == "" {
		return nil, NewBadRequestException("statusCode is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	response, err := stores.restApis.GetIntegrationResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toIntegrationResponseResponse(response), nil
}

// DeleteIntegrationResponse deletes an integration response from a method in API Gateway.
func (s *APIGatewayService) DeleteIntegrationResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	statusCode := request.GetStringParam(req.Parameters, "statusCode")
	if statusCode == "" {
		statusCode = getPathParam(req, "statusCode")
	}
	if statusCode == "" {
		return nil, NewBadRequestException("statusCode is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteIntegrationResponse(apiId, resourceId, httpMethod, statusCode); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

func (s *APIGatewayService) toIntegrationResponse(i *store.Integration) map[string]interface{} {
	response := map[string]interface{}{
		"type": i.Type,
	}

	if i.IntegrationHttpMethod != "" {
		response["httpMethod"] = i.IntegrationHttpMethod
	}
	if i.Uri != "" {
		response["uri"] = i.Uri
	}
	if i.Credentials != "" {
		response["credentials"] = i.Credentials
	}
	if i.PassthroughBehavior != "" {
		response["passthroughBehavior"] = i.PassthroughBehavior
	}
	if i.ContentHandling != "" {
		response["contentHandling"] = i.ContentHandling
	}
	if i.CacheNamespace != "" {
		response["cacheNamespace"] = i.CacheNamespace
	}
	if i.ConnectionType != "" {
		response["connectionType"] = i.ConnectionType
	}
	if i.ConnectionId != "" {
		response["connectionId"] = i.ConnectionId
	}
	if len(i.CacheKeyParameters) > 0 {
		response["cacheKeyParameters"] = i.CacheKeyParameters
	}
	if i.TimeoutInMillis > 0 {
		response["timeoutInMillis"] = i.TimeoutInMillis
	}
	if len(i.RequestParameters) > 0 {
		response["requestParameters"] = i.RequestParameters
	}
	if len(i.RequestTemplates) > 0 {
		response["requestTemplates"] = i.RequestTemplates
	}

	return response
}

func (s *APIGatewayService) toIntegrationResponseResponse(r *store.IntegrationResponse) map[string]interface{} {
	response := map[string]interface{}{
		"statusCode": r.StatusCode,
	}

	if r.SelectionPattern != "" {
		response["selectionPattern"] = r.SelectionPattern
	}
	if r.ContentHandling != "" {
		response["contentHandling"] = r.ContentHandling
	}
	if len(r.ResponseParameters) > 0 {
		response["responseParameters"] = r.ResponseParameters
	}
	if len(r.ResponseTemplates) > 0 {
		response["responseTemplates"] = r.ResponseTemplates
	}

	return response
}
