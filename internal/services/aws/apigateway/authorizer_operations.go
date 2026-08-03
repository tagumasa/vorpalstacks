// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/apigateway"
)

// CreateAuthorizer creates a new authorizer for the specified REST API.
func (s *APIGatewayService) CreateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)

	in := &AuthorizerInput{
		Name:     request.GetStringParam(req.Parameters, "name"),
		Type:     request.GetStringParam(req.Parameters, "type"),
		AuthType: request.GetStringParam(req.Parameters, "authType"),
	}
	if v := request.GetStringParam(req.Parameters, "authorizerUri"); v != "" {
		in.AuthorizerUri = v
	}
	if v := request.GetStringParam(req.Parameters, "authorizerCredentials"); v != "" {
		in.AuthorizerCredentials = v
	}
	if v := request.GetStringParam(req.Parameters, "identitySource"); v != "" {
		in.IdentitySource = v
	}
	if v := request.GetStringParam(req.Parameters, "identityValidationExpression"); v != "" {
		in.IdentityValidationExpression = v
	}
	if _, ok := req.Parameters["authorizerResultTtlInSeconds"]; ok {
		v := int32(request.GetIntParam(req.Parameters, "authorizerResultTtlInSeconds"))
		in.AuthorizerResultTtlInSeconds = &v
	}
	if providerArns, ok := req.Parameters["providerARNs"].([]interface{}); ok {
		for _, item := range providerArns {
			if str, ok := item.(string); ok {
				in.ProviderArns = append(in.ProviderArns, str)
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createAuthorizerCore(stores, apiId, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toAuthorizerResponse(result), nil
}

// GetAuthorizer retrieves an authorizer by its ID for the specified REST API.
func (s *APIGatewayService) GetAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	if authorizerId == "" {
		authorizerId = getPathParam(req, "authorizerId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.getAuthorizerCore(stores, apiId, authorizerId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toAuthorizerResponse(result), nil
}

// UpdateAuthorizer updates an existing authorizer for the specified REST API.
func (s *APIGatewayService) UpdateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	if authorizerId == "" {
		authorizerId = getPathParam(req, "authorizerId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	existing, err := s.updateAuthorizerCore(stores, apiId, authorizerId, parsePatchOperations(req.Parameters))
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toAuthorizerResponse(existing), nil
}

// DeleteAuthorizer deletes an authorizer from the specified REST API.
func (s *APIGatewayService) DeleteAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	if authorizerId == "" {
		authorizerId = getPathParam(req, "authorizerId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteAuthorizerCore(stores, apiId, authorizerId); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// GetAuthorizers lists all authorizers for the specified REST API.
func (s *APIGatewayService) GetAuthorizers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	authorizers, err := s.listAuthorizersCore(stores, apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	items := make([]interface{}, 0, len(authorizers))
	for _, a := range authorizers {
		items = append(items, s.toAuthorizerResponse(a))
	}

	page, nextPos, found := paginateItems(items, position, limit)
	if !found {
		return nil, NewBadRequestException("Invalid position: " + position)
	}
	result := map[string]interface{}{
		"item": page,
	}
	if nextPos != "" {
		result["position"] = nextPos
	}
	return result, nil
}

func (s *APIGatewayService) toAuthorizerResponse(a *apigateway.Authorizer) map[string]interface{} {
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
