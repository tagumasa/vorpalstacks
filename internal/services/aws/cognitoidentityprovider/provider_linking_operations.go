package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// Handlers for the federated provider linking family; validation and store
// access live in provider_linking_core.go.

// AdminDisableProviderForUser disables a federated provider for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminDisableProviderForUser.html
func (s *CognitoService) AdminDisableProviderForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	user, _ := req.Parameters["User"].(map[string]interface{})

	if err := s.adminDisableProviderForUserCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), user); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// AdminLinkProviderForUser links a federated provider to an existing user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminLinkProviderForUser.html
func (s *CognitoService) AdminLinkProviderForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	destinationUser, _ := req.Parameters["DestinationUser"].(map[string]interface{})
	sourceUser, _ := req.Parameters["SourceUser"].(map[string]interface{})

	if err := s.adminLinkProviderForUserCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), destinationUser, sourceUser); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
