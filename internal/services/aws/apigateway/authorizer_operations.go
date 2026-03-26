// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	store "vorpalstacks/internal/store/aws/apigateway"
)

// CreateAuthorizer creates a new authorizer for the specified REST API.
func (s *APIGatewayService) CreateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	authType := request.GetStringParam(req.Parameters, "type")
	if authType == "" {
		authType = "TOKEN"
	}

	authorizer := &store.Authorizer{
		Name: name,
		Type: authType,
	}

	if providerArns, ok := req.Parameters["providerARNs"].([]interface{}); ok {
		for _, item := range providerArns {
			if str, ok := item.(string); ok {
				authorizer.ProviderArns = append(authorizer.ProviderArns, str)
			}
		}
	}

	if v := request.GetStringParam(req.Parameters, "authType"); v != "" {
		authorizer.AuthType = v
	}
	if v := request.GetStringParam(req.Parameters, "authorizerUri"); v != "" {
		authorizer.AuthorizerUri = v
	}
	if v := request.GetStringParam(req.Parameters, "authorizerCredentials"); v != "" {
		authorizer.AuthorizerCredentials = v
	}
	if v := request.GetStringParam(req.Parameters, "identitySource"); v != "" {
		authorizer.IdentitySource = v
	}
	if v := request.GetStringParam(req.Parameters, "identityValidationExpression"); v != "" {
		authorizer.IdentityValidationExpression = v
	}

	ttl := request.GetIntParam(req.Parameters, "authorizerResultTtlInSeconds")
	if ttl > 0 {
		authorizer.AuthorizerResultTtlInSeconds = int32(ttl)
	} else {
		authorizer.AuthorizerResultTtlInSeconds = 300
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.restApis.CreateAuthorizer(apiId, authorizer)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toAuthorizerResponse(result), nil
}

// GetAuthorizer retrieves an authorizer by its ID for the specified REST API.
func (s *APIGatewayService) GetAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	if authorizerId == "" {
		authorizerId = getPathParam(req, "authorizerId")
	}
	if authorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.restApis.GetAuthorizer(apiId, authorizerId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toAuthorizerResponse(result), nil
}

// UpdateAuthorizer updates an existing authorizer for the specified REST API.
func (s *APIGatewayService) UpdateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	if authorizerId == "" {
		authorizerId = getPathParam(req, "authorizerId")
	}
	if authorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	existing, err := stores.restApis.GetAuthorizer(apiId, authorizerId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	patchOps, ok := req.Parameters["patchOperations"].([]interface{})
	if ok {
		for _, op := range patchOps {
			if opMap, ok := op.(map[string]interface{}); ok {
				path := ""
				value := ""
				if p, ok := opMap["path"].(string); ok {
					path = p
				}
				if v, ok := opMap["value"].(string); ok {
					value = v
				}

				switch path {
				case "/name":
					existing.Name = value
				case "/type":
					existing.Type = value
				case "/authType":
					existing.AuthType = value
				case "/authorizerUri":
					existing.AuthorizerUri = value
				case "/authorizerCredentials":
					existing.AuthorizerCredentials = value
				case "/identitySource":
					existing.IdentitySource = value
				case "/identityValidationExpression":
					existing.IdentityValidationExpression = value
				case "/authorizerResultTtlInSeconds":
					existing.AuthorizerResultTtlInSeconds = int32(parseInt64(value))
				}
			}
		}
	}

	if err := stores.restApis.UpdateAuthorizer(apiId, existing); err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toAuthorizerResponse(existing), nil
}

// DeleteAuthorizer deletes an authorizer from the specified REST API.
func (s *APIGatewayService) DeleteAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	if authorizerId == "" {
		authorizerId = getPathParam(req, "authorizerId")
	}
	if authorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteAuthorizer(apiId, authorizerId); err != nil {
		return nil, toApiGatewayError(err)
	}

	return response.EmptyResponse(), nil
}

// GetAuthorizers lists all authorizers for the specified REST API.
func (s *APIGatewayService) GetAuthorizers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	authorizers, err := stores.restApis.ListAuthorizers(apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	items := make([]interface{}, 0, len(authorizers))
	for _, a := range authorizers {
		items = append(items, s.toAuthorizerResponse(a))
	}

	return map[string]interface{}{
		"item": items,
	}, nil
}

func (s *APIGatewayService) toAuthorizerResponse(a *store.Authorizer) map[string]interface{} {
	response := map[string]interface{}{
		"id":                           a.Id,
		"name":                         a.Name,
		"type":                         a.Type,
		"authorizerResultTtlInSeconds": a.AuthorizerResultTtlInSeconds,
	}

	if len(a.ProviderArns) > 0 {
		response["providerARNs"] = a.ProviderArns
	}
	if a.AuthType != "" {
		response["authType"] = a.AuthType
	}
	if a.AuthorizerUri != "" {
		response["authorizerUri"] = a.AuthorizerUri
	}
	if a.AuthorizerCredentials != "" {
		response["authorizerCredentials"] = a.AuthorizerCredentials
	}
	if a.IdentitySource != "" {
		response["identitySource"] = a.IdentitySource
	}
	if a.IdentityValidationExpression != "" {
		response["identityValidationExpression"] = a.IdentityValidationExpression
	}

	return response
}

func toApiGatewayError(err error) *ApiGatewayError {
	if err == nil {
		return nil
	}
	return GetApiGatewayError(err)
}
