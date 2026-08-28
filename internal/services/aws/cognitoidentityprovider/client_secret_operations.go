package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"encoding/base64"

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

	return map[string]interface{}{
		"ClientSecretDescriptor": formatClientSecretDescriptor(descriptor),
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

const maxClientSecretsPerPage = 60

// ListUserPoolClientSecrets lists secrets for a user pool client.
// The AWS API does not expose a Limit parameter for this operation; page
// size is controlled server-side.
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

func formatClientSecretDescriptor(d cognitostore.ClientSecretDescriptor) map[string]interface{} {
	return map[string]interface{}{
		"ClientSecretId":         d.ClientSecretID,
		"ClientSecretValue":      d.ClientSecretValue,
		"ClientSecretCreateDate": d.ClientSecretCreateDate.Unix(),
	}
}

func generateSecretID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateSecretValue() (string, error) {
	b := make([]byte, 40)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
