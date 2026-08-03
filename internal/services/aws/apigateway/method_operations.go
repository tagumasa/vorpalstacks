// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/apigateway"
)

// PutMethod creates or updates a method for a resource in API Gateway.
func (s *APIGatewayService) PutMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")

	in := &MethodInput{
		AuthorizationType:  request.GetStringParam(req.Parameters, "authorizationType"),
		AuthorizerId:       request.GetStringParam(req.Parameters, "authorizerId"),
		ApiKeyRequired:     request.GetBoolParam(req.Parameters, "apiKeyRequired"),
		RequestValidatorId: request.GetStringParam(req.Parameters, "requestValidatorId"),
		OperationName:      request.GetStringParam(req.Parameters, "operationName"),
	}
	if rp, ok := req.Parameters["requestParameters"].(map[string]interface{}); ok {
		in.RequestParameters = make(map[string]bool)
		for k, v := range rp {
			switch tv := v.(type) {
			case bool:
				in.RequestParameters[k] = tv
			case string:
				in.RequestParameters[k] = tv == "true"
			}
		}
	}
	if rm, ok := req.Parameters["requestModels"].(map[string]interface{}); ok {
		in.RequestModels = make(map[string]string)
		for k, v := range rm {
			if vs, ok := v.(string); ok {
				in.RequestModels[k] = vs
			}
		}
	}
	if scopes, ok := req.Parameters["authorizationScopes"].([]interface{}); ok {
		for _, scope := range scopes {
			if ss, ok := scope.(string); ok {
				in.AuthorizationScopes = append(in.AuthorizationScopes, ss)
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.putMethodCore(stores, apiId, resourceId, httpMethod, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toMethodResponse(created), nil
}

// GetMethod retrieves a method for a resource in API Gateway.
func (s *APIGatewayService) GetMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	method, err := s.getMethodCore(stores, apiId, resourceId, httpMethod)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toMethodResponse(method), nil
}

// DeleteMethod deletes a method from a resource in API Gateway.
func (s *APIGatewayService) DeleteMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteMethodCore(stores, apiId, resourceId, httpMethod); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

func (s *APIGatewayService) toMethodResponse(m *apigateway.Method) map[string]interface{} {
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
	if len(m.AuthorizationScopes) > 0 {
		response["authorizationScopes"] = m.AuthorizationScopes
	}
	if m.MethodIntegration != nil {
		response["methodIntegration"] = s.toIntegrationResponse(m.MethodIntegration)
	}
	if len(m.MethodResponses) > 0 {
		response["methodResponses"] = s.toMethodResponsesMap(m.MethodResponses)
	}

	return response
}

func (s *APIGatewayService) toMethodResponsesMap(responses map[string]*apigateway.MethodResponse) map[string]interface{} {
	result := make(map[string]interface{})
	for statusCode, resp := range responses {
		result[statusCode] = s.toMethodResponseSingle(resp)
	}
	return result
}

func (s *APIGatewayService) toMethodResponseSingle(resp *apigateway.MethodResponse) map[string]interface{} {
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
