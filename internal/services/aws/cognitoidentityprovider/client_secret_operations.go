package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// AddUserPoolClientSecret adds a new secret to a user pool client (multi-secret support).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AddUserPoolClientSecret.html
func (s *CognitoService) AddUserPoolClientSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	descriptor, err := s.addUserPoolClientSecretCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), req.GetParam("ClientId"), req.GetParam("ClientSecret"))
	if err != nil {
		return nil, err
	}

	formatted := formatClientSecretDescriptor(descriptor)
	// The actual value is returned only when creating a secret the service
	// generated; a caller-supplied secret is never echoed.
	if descriptor.Generated {
		formatted["ClientSecretValue"] = descriptor.ClientSecretValue
	}
	return map[string]interface{}{
		"ClientSecretDescriptor": formatted,
	}, nil
}

// DeleteUserPoolClientSecret removes a secret from a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPoolClientSecret.html
func (s *CognitoService) DeleteUserPoolClientSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteUserPoolClientSecretCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), req.GetParam("ClientId"), req.GetParam("ClientSecretId")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListUserPoolClientSecrets lists secrets for a user pool client.
// The AWS API does not expose a Limit parameter for this operation; page
// size is controlled server-side (maxClientSecretsPerPage, durations.go).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPoolClientSecrets.html
func (s *CognitoService) ListUserPoolClientSecrets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.listUserPoolClientSecretsCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), req.GetParam("ClientId"), req.GetParam("NextToken"))
	if err != nil {
		return nil, err
	}

	secrets := make([]map[string]interface{}, 0, len(result.Secrets))
	for _, d := range result.Secrets {
		secrets = append(secrets, formatClientSecretDescriptor(d))
	}

	resp := map[string]interface{}{"ClientSecrets": secrets}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}

// formatClientSecretDescriptor renders a descriptor's non-secret members.
// The stored value is never disclosed here: the list response never reveals
// it, and the Add response includes it only for service-generated secrets,
// which the Add handler adds separately.
func formatClientSecretDescriptor(d cognitostore.ClientSecretDescriptor) map[string]interface{} {
	return map[string]interface{}{
		"ClientSecretId":         d.ClientSecretID,
		"ClientSecretCreateDate": d.ClientSecretCreateDate.Unix(),
	}
}
