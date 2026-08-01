// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"strconv"
	"strings"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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

	if !validateAuthorizerType(authType) {
		return nil, NewBadRequestException("Invalid authorizer type: " + authType)
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

	// authorizerUri is required for TOKEN/REQUEST; COGNITO_USER_POOLS uses providerARNs.
	if authType != "COGNITO_USER_POOLS" && authorizer.AuthorizerUri == "" {
		return nil, NewBadRequestException("authorizerUri is required for " + authType + " authorizers")
	}

	// Default identitySource for TOKEN/REQUEST (AWS defaults to this header).
	if authorizer.IdentitySource == "" && authType != "COGNITO_USER_POOLS" {
		authorizer.IdentitySource = "method.request.header.Authorization"
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
		if ttl > 3600 {
			return nil, NewBadRequestException("authorizerResultTtlInSeconds must be between 0 and 3600")
		}
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

	stores.keyLocker.Lock(apiId + ":" + authorizerId)
	defer stores.keyLocker.Unlock(apiId + ":" + authorizerId)

	existing, err := stores.restApis.GetAuthorizer(apiId, authorizerId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch {
		case po.Path == "/name":
			existing.Name = po.Value
		case po.Path == "/type":
			if !validateAuthorizerType(po.Value) {
				return nil, NewBadRequestException("Invalid authorizer type: " + po.Value)
			}
			existing.Type = po.Value
		case po.Path == "/authType":
			existing.AuthType = po.Value
		case po.Path == "/authorizerUri":
			if po.Value != "" && !strings.HasPrefix(po.Value, "arn:") {
				return nil, NewBadRequestException("authorizerUri must be a valid ARN")
			}
			existing.AuthorizerUri = po.Value
		case po.Path == "/authorizerCredentials":
			existing.AuthorizerCredentials = po.Value
		case po.Path == "/identitySource":
			existing.IdentitySource = po.Value
		case po.Path == "/identityValidationExpression":
			existing.IdentityValidationExpression = po.Value
		case po.Path == "/authorizerResultTtlInSeconds":
			v, err := parseInt64(po.Value)
			if err != nil {
				return nil, NewBadRequestException("invalid authorizerResultTtlInSeconds: not a number")
			}
			if v < 0 || v > 3600 {
				return nil, NewBadRequestException("authorizerResultTtlInSeconds must be between 0 and 3600")
			}
			existing.AuthorizerResultTtlInSeconds = int32(v)
		case strings.HasPrefix(po.Path, "/providerARNs"):
			if po.Op == "remove" {
				if idx, err := strconv.Atoi(strings.TrimPrefix(po.Path, "/providerARNs/")); err == nil && idx < len(existing.ProviderArns) {
					existing.ProviderArns = append(existing.ProviderArns[:idx], existing.ProviderArns[idx+1:]...)
				}
			} else {
				if !sliceContains(existing.ProviderArns, po.Value) {
					existing.ProviderArns = append(existing.ProviderArns, po.Value)
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

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	position := request.GetStringParam(req.Parameters, "position")

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
