package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// Handler for the machine-to-machine token family; validation, scope
// resolution and token issuance live in client_token_core.go.

// GetClientToken issues an M2M access token for a confidential app client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetClientToken.html
func (s *CognitoService) GetClientToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getClientTokenCore(ctx, reqCtx, GetClientTokenInput{
		ClientID:       req.GetParam("ClientId"),
		Secret:         req.GetParam("Secret"),
		Scopes:         getStringSliceParam(req, "Scopes"),
		ClientMetadata: parseClientMetadata(req),
	})
}
