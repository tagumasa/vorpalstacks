package cognitoidentityprovider

import (
	"encoding/base64"
	"strconv"
	"time"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Core functions for the user-pool client multi-secret family. The handlers
// extract the wire members; validation, secret minting and store access
// live here.

// addUserPoolClientSecretCore attaches a new secret descriptor to a user
// pool client. An empty secretValue triggers generation of a fresh secret.
func (s *CognitoService) addUserPoolClientSecretCore(region, userPoolID, clientID, secretValue string) (cognitostore.ClientSecretDescriptor, error) {
	if userPoolID == "" || clientID == "" {
		return cognitostore.ClientSecretDescriptor{}, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return cognitostore.ClientSecretDescriptor{}, err
	}

	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return cognitostore.ClientSecretDescriptor{}, ErrResourceNotFound
	}

	now := time.Now().UTC()
	secretID, err := generateSecretID()
	if err != nil {
		return cognitostore.ClientSecretDescriptor{}, ErrInternalError
	}
	descriptor := cognitostore.ClientSecretDescriptor{
		ClientSecretID:         "secret-" + secretID,
		ClientSecretValue:      secretValue,
		ClientSecretCreateDate: now,
	}
	if descriptor.ClientSecretValue == "" {
		generated, gerr := generateSecretValue()
		if gerr != nil {
			return cognitostore.ClientSecretDescriptor{}, ErrInternalError
		}
		descriptor.ClientSecretValue = generated
	}
	client.ClientSecrets = append(client.ClientSecrets, descriptor)
	client.LastModifiedDate = now
	if err := store.UpdateUserPoolClient(client); err != nil {
		return cognitostore.ClientSecretDescriptor{}, ErrInternalError
	}

	return descriptor, nil
}

// deleteUserPoolClientSecretCore removes one secret descriptor from a user
// pool client; an unknown secret id is a not-found error.
func (s *CognitoService) deleteUserPoolClientSecretCore(region, userPoolID, clientID, secretID string) error {
	if userPoolID == "" || clientID == "" || secretID == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return ErrResourceNotFound
	}

	found := false
	filtered := client.ClientSecrets[:0]
	for _, sec := range client.ClientSecrets {
		if sec.ClientSecretID == secretID {
			found = true
			continue
		}
		filtered = append(filtered, sec)
	}
	if !found {
		return ErrResourceNotFound
	}

	client.ClientSecrets = filtered
	client.LastModifiedDate = time.Now().UTC()
	if err := store.UpdateUserPoolClient(client); err != nil {
		return ErrInternalError
	}

	return nil
}

// ListUserPoolClientSecretsResult is one page of client secret descriptors.
type ListUserPoolClientSecretsResult struct {
	Secrets   []cognitostore.ClientSecretDescriptor
	NextToken string
}

// listUserPoolClientSecretsCore pages through a client's secret descriptors,
// decoding and re-encoding the base64 offset continuation token.
func (s *CognitoService) listUserPoolClientSecretsCore(region, userPoolID, clientID, nextToken string) (*ListUserPoolClientSecretsResult, error) {
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	start := 0
	if nextToken != "" {
		if decoded, err := decodeOffsetToken(nextToken); err == nil && decoded >= 0 {
			start = decoded
		}
	}

	total := len(client.ClientSecrets)
	if start > total {
		start = total
	}
	end := start + maxClientSecretsPerPage
	if end > total {
		end = total
	}

	return &ListUserPoolClientSecretsResult{
		Secrets:   client.ClientSecrets[start:end],
		NextToken: encodeOffsetToken(end, total),
	}, nil
}

// decodeOffsetToken decodes a base64 raw-URL offset continuation token.
func decodeOffsetToken(token string) (int, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(decoded))
}

// encodeOffsetToken mints the continuation token for the next page; an
// empty string means the listing is exhausted.
func encodeOffsetToken(end, total int) string {
	if end < total {
		return base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(end)))
	}
	return ""
}
