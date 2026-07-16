// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	store "vorpalstacks/internal/store/aws/apigateway"
)

// PutMethodResponse creates or updates a method response for a method in API Gateway.
func (s *APIGatewayService) PutMethodResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	response := &store.MethodResponse{
		StatusCode: statusCode,
	}

	if respParams, ok := req.Parameters["responseParameters"].(map[string]interface{}); ok {
		response.ResponseParameters = make(map[string]bool)
		for k, v := range respParams {
			if b, ok := v.(bool); ok {
				response.ResponseParameters[k] = b
			}
		}
	}

	if respModels, ok := req.Parameters["responseModels"].(map[string]interface{}); ok {
		response.ResponseModels = make(map[string]string)
		for k, v := range respModels {
			if s, ok := v.(string); ok {
				response.ResponseModels[k] = s
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.restApis.PutMethodResponse(apiId, resourceId, httpMethod, statusCode, response)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toMethodResponseResponse(result), nil
}

// GetMethodResponse retrieves a method response for a method in API Gateway.
func (s *APIGatewayService) GetMethodResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	result, err := stores.restApis.GetMethodResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toMethodResponseResponse(result), nil
}

// DeleteMethodResponse deletes a method response from a method in API Gateway.
func (s *APIGatewayService) DeleteMethodResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	if err := stores.restApis.DeleteMethodResponse(apiId, resourceId, httpMethod, statusCode); err != nil {
		return nil, toApiGatewayError(err)
	}

	return response.EmptyResponse(), nil
}

func (s *APIGatewayService) toMethodResponseResponse(r *store.MethodResponse) map[string]interface{} {
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
