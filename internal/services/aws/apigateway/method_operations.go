// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	store "vorpalstacks/internal/store/aws/apigateway"
)

// PutMethod creates or updates a method for a resource in API Gateway.
func (s *APIGatewayService) PutMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	authorizationType := request.GetStringParam(req.Parameters, "authorizationType")
	if authorizationType == "" {
		authorizationType = "NONE"
	}

	method := &store.Method{
		HttpMethod:         httpMethod,
		AuthorizationType:  authorizationType,
		AuthorizerId:       request.GetStringParam(req.Parameters, "authorizerId"),
		ApiKeyRequired:     request.GetBoolParam(req.Parameters, "apiKeyRequired"),
		RequestValidatorId: request.GetStringParam(req.Parameters, "requestValidatorId"),
		OperationName:      request.GetStringParam(req.Parameters, "operationName"),
	}

	if reqParams, ok := req.Parameters["requestParameters"].(map[string]interface{}); ok {
		method.RequestParameters = make(map[string]bool)
		for k, v := range reqParams {
			if vb, ok := v.(bool); ok {
				method.RequestParameters[k] = vb
			}
		}
	}

	if reqModels, ok := req.Parameters["requestModels"].(map[string]interface{}); ok {
		method.RequestModels = make(map[string]string)
		for k, v := range reqModels {
			if vs, ok := v.(string); ok {
				method.RequestModels[k] = vs
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.restApis.PutMethod(apiId, resourceId, method)
	if err != nil {
		return nil, err
	}

	return s.toMethodResponse(created), nil
}

// GetMethod retrieves a method for a resource in API Gateway.
func (s *APIGatewayService) GetMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	method, err := stores.restApis.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toMethodResponse(method), nil
}

// DeleteMethod deletes a method from a resource in API Gateway.
func (s *APIGatewayService) DeleteMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	if err := stores.restApis.DeleteMethod(apiId, resourceId, httpMethod); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

func (s *APIGatewayService) toMethodResponse(m *store.Method) map[string]interface{} {
	response := map[string]interface{}{
		"httpMethod":        m.HttpMethod,
		"authorizationType": m.AuthorizationType,
		"apiKeyRequired":    m.ApiKeyRequired,
	}

	if m.AuthorizerId != "" {
		response["authorizerId"] = m.AuthorizerId
	}
	if m.RequestValidatorId != "" {
		response["requestValidatorId"] = m.RequestValidatorId
	}
	if m.OperationName != "" {
		response["operationName"] = m.OperationName
	}
	if len(m.RequestParameters) > 0 {
		response["requestParameters"] = m.RequestParameters
	}
	if len(m.RequestModels) > 0 {
		response["requestModels"] = m.RequestModels
	}
	if m.MethodIntegration != nil {
		response["methodIntegration"] = s.toIntegrationResponse(m.MethodIntegration)
	}
	if len(m.MethodResponses) > 0 {
		response["methodResponses"] = s.toMethodResponsesMap(m.MethodResponses)
	}

	return response
}

func (s *APIGatewayService) toMethodResponsesMap(responses map[string]*store.MethodResponse) map[string]interface{} {
	result := make(map[string]interface{})
	for statusCode, resp := range responses {
		result[statusCode] = s.toMethodResponseSingle(resp)
	}
	return result
}

func (s *APIGatewayService) toMethodResponseSingle(resp *store.MethodResponse) map[string]interface{} {
	result := map[string]interface{}{
		"statusCode": resp.StatusCode,
	}
	if len(resp.ResponseParameters) > 0 {
		result["responseParameters"] = resp.ResponseParameters
	}
	if len(resp.ResponseModels) > 0 {
		result["responseModels"] = resp.ResponseModels
	}
	return result
}
