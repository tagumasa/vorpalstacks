// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/apigateway"
)

// methodResponseRequestParams extracts the method-response identity members
// from the wire request: api id, resource id, HTTP method (query, then path
// label) and status code (query, then path label).
func methodResponseRequestParams(req *request.ParsedRequest) (apiId, resourceId, httpMethod, statusCode string) {
	apiId = getRestApiId(req)
	resourceId = getResourceId(req)
	httpMethod = request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	statusCode = request.GetStringParam(req.Parameters, "statusCode")
	if statusCode == "" {
		statusCode = getPathParam(req, "statusCode")
	}
	return apiId, resourceId, httpMethod, statusCode
}

// PutMethodResponse creates or updates a method response for a method in API Gateway.
func (s *APIGatewayService) PutMethodResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId, httpMethod, statusCode := methodResponseRequestParams(req)

	in := &MethodResponseInput{
		StatusCode: statusCode,
	}
	if respParams, ok := req.Parameters["responseParameters"].(map[string]interface{}); ok {
		in.ResponseParameters = make(map[string]bool)
		for k, v := range respParams {
			if b, ok := v.(bool); ok {
				in.ResponseParameters[k] = b
			}
		}
	}
	if respModels, ok := req.Parameters["responseModels"].(map[string]interface{}); ok {
		in.ResponseModels = make(map[string]string)
		for k, v := range respModels {
			if str, ok := v.(string); ok {
				in.ResponseModels[k] = str
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.putMethodResponseCore(stores, apiId, resourceId, httpMethod, statusCode, in)
	if err != nil {
		return nil, err
	}
	return s.toMethodResponseResponse(result), nil
}

// GetMethodResponse retrieves a method response for a method in API Gateway.
func (s *APIGatewayService) GetMethodResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId, httpMethod, statusCode := methodResponseRequestParams(req)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.getMethodResponseCore(stores, apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, err
	}
	return s.toMethodResponseResponse(result), nil
}

// DeleteMethodResponse deletes a method response from a method in API Gateway.
func (s *APIGatewayService) DeleteMethodResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId, httpMethod, statusCode := methodResponseRequestParams(req)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteMethodResponseCore(stores, apiId, resourceId, httpMethod, statusCode); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UpdateMethodResponse updates a method response via patch operations.
func (s *APIGatewayService) UpdateMethodResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId, httpMethod, statusCode := methodResponseRequestParams(req)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	updatedResp, err := s.updateMethodResponseCore(stores, apiId, resourceId, httpMethod, statusCode, ops)
	if err != nil {
		return nil, err
	}
	return s.toMethodResponseResponse(updatedResp), nil
}

func (s *APIGatewayService) toMethodResponseResponse(r *apigateway.MethodResponse) map[string]interface{} {
	response := map[string]interface{}{
		"statusCode": r.StatusCode,
	}

	if len(r.ResponseParameters) > 0 {
		response["responseParameters"] = r.ResponseParameters
	}
	if len(r.ResponseModels) > 0 {
		response["responseModels"] = r.ResponseModels
	}

	return response
}
