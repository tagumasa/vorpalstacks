package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// AddUserPoolClientSecret adds a new secret to a user pool client (multi-secret support).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AddUserPoolClientSecret.html
func (s *CognitoService) AddUserPoolClientSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	clientID := req.GetParam("ClientId")
	secretValue := req.GetParam("ClientSecret")
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	now := time.Now().UTC()
	descriptor := cognitostore.ClientSecretDescriptor{
		ClientSecretID:         "secret-" + generateSecretID(),
		ClientSecretValue:      secretValue,
		ClientSecretCreateDate: now,
	}
	if descriptor.ClientSecretValue == "" {
		descriptor.ClientSecretValue = generateSecretValue()
	}

	client.ClientSecrets = append(client.ClientSecrets, descriptor)
	client.LastModifiedDate = now
	if err := store.UpdateUserPoolClient(client); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"ClientSecretDescriptor": formatClientSecretDescriptor(descriptor),
	}, nil
}

// DeleteUserPoolClientSecret removes a secret from a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPoolClientSecret.html
func (s *CognitoService) DeleteUserPoolClientSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	clientID := req.GetParam("ClientId")
	secretID := req.GetParam("ClientSecretId")
	if userPoolID == "" || clientID == "" || secretID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	found := false
	filtered := client.ClientSecrets[:0]
	for _, s := range client.ClientSecrets {
		if s.ClientSecretID == secretID {
			found = true
			continue
		}
		filtered = append(filtered, s)
	}
	if !found {
		return nil, ErrResourceNotFound
	}

	client.ClientSecrets = filtered
	client.LastModifiedDate = time.Now().UTC()
	if err := store.UpdateUserPoolClient(client); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ListUserPoolClientSecrets lists secrets for a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPoolClientSecrets.html
func (s *CognitoService) ListUserPoolClientSecrets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	clientID := req.GetParam("ClientId")
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	secrets := make([]map[string]interface{}, 0, len(client.ClientSecrets))
	for _, d := range client.ClientSecrets {
		secrets = append(secrets, formatClientSecretDescriptor(d))
	}

	return map[string]interface{}{
		"ClientSecrets": secrets,
	}, nil
}

func formatClientSecretDescriptor(d cognitostore.ClientSecretDescriptor) map[string]interface{} {
	return map[string]interface{}{
		"ClientSecretId":         d.ClientSecretID,
		"ClientSecretValue":      d.ClientSecretValue,
		"ClientSecretCreateDate": d.ClientSecretCreateDate.Unix(),
	}
}

func generateSecretID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateSecretValue() string {
	b := make([]byte, 40)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
