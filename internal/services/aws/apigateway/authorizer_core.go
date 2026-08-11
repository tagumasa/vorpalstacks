package apigateway

import (
	"strconv"
	"strings"

	"vorpalstacks/internal/store/aws/apigateway"
)

// AuthorizerInput is the transport-agnostic input for creating or replacing
// an authorizer.
type AuthorizerInput struct {
	Name                         string
	Type                         string
	AuthType                     string
	AuthorizerUri                string
	AuthorizerCredentials        string
	IdentitySource               string
	IdentityValidationExpression string
	AuthorizerResultTtlInSeconds *int32
	ProviderArns                 []string
}

// createAuthorizerCore persists an authorizer. Centralises the type
// default ("TOKEN"), the type validation, the authorizerUri requirement
// for TOKEN/REQUEST and the identitySource default ("method.request.header.Authorization").
func (s *APIGatewayService) createAuthorizerCore(
	stores *apiGatewayStores,
	apiId string,
	in *AuthorizerInput,
) (*apigateway.Authorizer, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}

	authType := in.Type
	if authType == "" {
		authType = "TOKEN"
	}
	if !validateAuthorizerType(authType) {
		return nil, NewBadRequestException("Invalid authorizer type: " + authType)
	}
	if authType != "COGNITO_USER_POOLS" && in.AuthorizerUri == "" {
		return nil, NewBadRequestException("authorizerUri is required for " + authType + " authorizers")
	}

	identitySource := in.IdentitySource
	if identitySource == "" && authType != "COGNITO_USER_POOLS" {
		identitySource = "method.request.header.Authorization"
	}

	// Distinguish "unset" (nil → default 300) from "explicitly 0"
	// (cache disabled): AuthorizerResultTtlInSeconds is *int32 in the
	// Smithy model so the optional TTL field can express the AWS
	// "cache disabled" state, which an int32 zero value could not.
	ttl := int32(300)
	if in.AuthorizerResultTtlInSeconds != nil {
		ttl = *in.AuthorizerResultTtlInSeconds
	}
	if !validateAuthorizerTtl(ttl) {
		return nil, NewBadRequestException("authorizerResultTtlInSeconds must be between 0 and 3600")
	}

	authorizer := &apigateway.Authorizer{
		Name:                         in.Name,
		Type:                         authType,
		AuthType:                     in.AuthType,
		AuthorizerUri:                in.AuthorizerUri,
		AuthorizerCredentials:        in.AuthorizerCredentials,
		IdentitySource:               identitySource,
		IdentityValidationExpression: in.IdentityValidationExpression,
		AuthorizerResultTtlInSeconds: ttl,
		ProviderArns:                 in.ProviderArns,
	}

	created, err := stores.restApis.CreateAuthorizer(apiId, authorizer)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getAuthorizerCore retrieves an authorizer by id.
func (s *APIGatewayService) getAuthorizerCore(stores *apiGatewayStores, apiId, authorizerId string) (*apigateway.Authorizer, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if authorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required")
	}
	result, err := stores.restApis.GetAuthorizer(apiId, authorizerId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return result, nil
}

// deleteAuthorizerCore removes an authorizer.
func (s *APIGatewayService) deleteAuthorizerCore(stores *apiGatewayStores, apiId, authorizerId string) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if authorizerId == "" {
		return NewBadRequestException("authorizerId is required")
	}
	if err := stores.restApis.DeleteAuthorizer(apiId, authorizerId); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}

// listAuthorizersCore returns all authorizers for an api id.
func (s *APIGatewayService) listAuthorizersCore(stores *apiGatewayStores, apiId string) ([]*apigateway.Authorizer, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	authorizers, err := stores.restApis.ListAuthorizers(apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return authorizers, nil
}

// updateAuthorizerCore applies patch operations to an authorizer under
// the per-authorizer key locker.
func (s *APIGatewayService) updateAuthorizerCore(
	stores *apiGatewayStores,
	apiId, authorizerId string,
	patches []PatchOperation,
) (*apigateway.Authorizer, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if authorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required")
	}

	stores.keyLocker.Lock(apiId + ":" + authorizerId)
	defer stores.keyLocker.Unlock(apiId + ":" + authorizerId)

	existing, err := stores.restApis.GetAuthorizer(apiId, authorizerId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range patches {
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
			if !validateAuthorizerTtl(int32(v)) {
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
	return existing, nil
}
