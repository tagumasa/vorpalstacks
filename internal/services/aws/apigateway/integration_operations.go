// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/apigateway"
)

// PutIntegration creates or updates an integration for a resource in API Gateway.
func (s *APIGatewayService) PutIntegration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")

	in := &IntegrationInput{
		Type:                  request.GetStringParam(req.Parameters, "type"),
		IntegrationHttpMethod: request.GetStringParam(req.Parameters, "integrationHttpMethod"),
		Uri:                   request.GetStringParam(req.Parameters, "uri"),
		Credentials:           request.GetStringParam(req.Parameters, "credentials"),
		PassthroughBehavior:   request.GetStringParam(req.Parameters, "passthroughBehavior"),
		ContentHandling:       request.GetStringParam(req.Parameters, "contentHandling"),
		CacheNamespace:        request.GetStringParam(req.Parameters, "cacheNamespace"),
		TimeoutInMillis:       int32(request.GetIntParam(req.Parameters, "timeoutInMillis")),
		ConnectionType:        request.GetStringParam(req.Parameters, "connectionType"),
		ConnectionId:          request.GetStringParam(req.Parameters, "connectionId"),
		ResponseTransferMode:  request.GetStringParam(req.Parameters, "responseTransferMode"),
		IntegrationTarget:     request.GetStringParam(req.Parameters, "integrationTarget"),
	}
	if rp, ok := req.Parameters["requestParameters"].(map[string]interface{}); ok {
		in.RequestParameters = make(map[string]string)
		for k, v := range rp {
			if vs, ok := v.(string); ok {
				in.RequestParameters[k] = vs
			}
		}
	}
	if rt, ok := req.Parameters["requestTemplates"].(map[string]interface{}); ok {
		in.RequestTemplates = make(map[string]string)
		for k, v := range rt {
			if vs, ok := v.(string); ok {
				in.RequestTemplates[k] = vs
			}
		}
	}
	if ckp, ok := req.Parameters["cacheKeyParameters"].([]interface{}); ok {
		for _, p := range ckp {
			if ps, ok := p.(string); ok {
				in.CacheKeyParameters = append(in.CacheKeyParameters, ps)
			}
		}
	}
	if tlsConfigMap, ok := req.Parameters["tlsConfig"].(map[string]interface{}); ok {
		in.TlsConfig = &apigateway.TlsConfig{}
		if v, ok := tlsConfigMap["insecureSkipVerification"].(bool); ok {
			in.TlsConfig.InsecureSkipVerification = v
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.putIntegrationCore(stores, apiId, resourceId, httpMethod, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toIntegrationResponse(created), nil
}

// GetIntegration retrieves an integration for a resource in API Gateway.
func (s *APIGatewayService) GetIntegration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	integration, err := s.getIntegrationCore(stores, apiId, resourceId, httpMethod)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toIntegrationResponse(integration), nil
}

// DeleteIntegration deletes an integration from a resource in API Gateway.
func (s *APIGatewayService) DeleteIntegration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteIntegrationCore(stores, apiId, resourceId, httpMethod); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// PutIntegrationResponse creates or updates an integration response for a method in API Gateway.
func (s *APIGatewayService) PutIntegrationResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		httpMethod = getPathParam(req, "httpMethod")
	}
	statusCode := request.GetStringParam(req.Parameters, "statusCode")

	in := &IntegrationResponseInput{
		SelectionPattern: request.GetStringParam(req.Parameters, "selectionPattern"),
		ContentHandling:  request.GetStringParam(req.Parameters, "contentHandling"),
	}
	if rp, ok := req.Parameters["responseParameters"].(map[string]interface{}); ok {
		in.ResponseParameters = make(map[string]string)
		for k, v := range rp {
			if vs, ok := v.(string); ok {
				in.ResponseParameters[k] = vs
			}
		}
	}
	if rt, ok := req.Parameters["responseTemplates"].(map[string]interface{}); ok {
		in.ResponseTemplates = make(map[string]string)
		for k, v := range rt {
			if vs, ok := v.(string); ok {
				in.ResponseTemplates[k] = vs
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.putIntegrationResponseCore(stores, apiId, resourceId, httpMethod, statusCode, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toIntegrationResponseResponse(created), nil
}

// GetIntegrationResponse retrieves an integration response for a method in API Gateway.
func (s *APIGatewayService) GetIntegrationResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
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
	response, err := s.getIntegrationResponseCore(stores, apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toIntegrationResponseResponse(response), nil
}

// DeleteIntegrationResponse deletes an integration response from a method in API Gateway.
func (s *APIGatewayService) DeleteIntegrationResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
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
	if err := s.deleteIntegrationResponseCore(stores, apiId, resourceId, httpMethod, statusCode); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

func (s *APIGatewayService) toIntegrationResponse(i *apigateway.Integration) map[string]interface{} {
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
	response["timeoutInMillis"] = i.TimeoutInMillis
	if len(i.RequestParameters) > 0 {
		response["requestParameters"] = i.RequestParameters
	}
	if len(i.RequestTemplates) > 0 {
		response["requestTemplates"] = i.RequestTemplates
	}
	if i.TlsConfig != nil {
		response["tlsConfig"] = map[string]interface{}{
			"insecureSkipVerification": i.TlsConfig.InsecureSkipVerification,
		}
	}
	if i.ResponseTransferMode != "" {
		response["responseTransferMode"] = i.ResponseTransferMode
	}
	if i.IntegrationTarget != "" {
		response["integrationTarget"] = i.IntegrationTarget
	}
	if len(i.IntegrationResponses) > 0 {
		irMap := make(map[string]interface{})
		for code, ir := range i.IntegrationResponses {
			irMap[code] = s.toIntegrationResponseResponse(ir)
		}
		response["integrationResponses"] = irMap
	}

	return response
}

func (s *APIGatewayService) toIntegrationResponseResponse(r *apigateway.IntegrationResponse) map[string]interface{} {
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
