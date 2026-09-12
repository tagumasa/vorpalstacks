package cognitoidentityprovider

import (
	"errors"
	"slices"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// ListUserPoolClientsInput carries the pagination parameters for
// ListUserPoolClients in a wire-protocol-independent format.
type ListUserPoolClientsInput struct {
	UserPoolID string
	MaxResults int
	NextToken  string
}

// UserPoolClientSummary is the summary entry returned by ListUserPoolClients.
type UserPoolClientSummary struct {
	ClientID   string
	UserPoolID string
	ClientName string
}

// ListUserPoolClientsResult is the paginated result of ListUserPoolClients.
type ListUserPoolClientsResult struct {
	Clients   []UserPoolClientSummary
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// listUserPoolClientsCore lists user pool clients for a given pool with
// pagination. The store is resolved internally via GetStoreForRegion.
func (s *CognitoService) listUserPoolClientsCore(region string, in ListUserPoolClientsInput) (*ListUserPoolClientsResult, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	maxResults := applyListLimitDefaults(in.MaxResults)

	result, err := store.ListUserPoolClientsPaginated(in.UserPoolID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, err
	}

	clients := make([]UserPoolClientSummary, 0, len(result.Items))
	for _, c := range result.Items {
		clients = append(clients, UserPoolClientSummary{
			ClientID:   c.ClientID,
			UserPoolID: c.UserPoolID,
			ClientName: c.ClientName,
		})
	}

	return &ListUserPoolClientsResult{
		Clients:   clients,
		NextToken: result.NextMarker,
	}, nil
}

// describeUserPoolClientCore retrieves a single user pool client by ID.
// Returns the store-level UserPoolClient for format conversion by callers.
func (s *CognitoService) describeUserPoolClientCore(region, userPoolID, clientID string) (*cognitostore.UserPoolClient, error) {
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		if errors.Is(err, cognitostore.ErrClientNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, ErrInternalError
	}
	return client, nil
}

// deleteUserPoolClientCore deletes a user pool client by ID.
func (s *CognitoService) deleteUserPoolClientCore(region, userPoolID, clientID string) error {
	if userPoolID == "" || clientID == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteUserPoolClient(userPoolID, clientID); err != nil {
		return ErrClientNotFound
	}
	return nil
}

// createUserPoolClientCore creates a new user pool client. It accepts a
// pre-constructed store-level UserPoolClient carrying the wire-parsed
// configuration and the create-only ClientSecret member's value; the
// assembled configuration is validated here, the single validation path
// shared by every transport that creates a client.
func (s *CognitoService) createUserPoolClientCore(region string, client *cognitostore.UserPoolClient, customSecret string) (*cognitostore.UserPoolClient, error) {
	if client.UserPoolID == "" || client.ClientName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateClientNamePattern(client.ClientName) {
		return nil, ErrInvalidParameter
	}
	if err := validateUserPoolClientConfig(client); err != nil {
		return nil, err
	}

	// A custom ClientSecret replaces the constructor's generated value and
	// must satisfy the ClientSecretType bounds; it cannot be combined with
	// GenerateSecret ("You cannot specify both GenerateSecret as true and
	// provide a ClientSecret value"). A client with a custom secret is a
	// secret-bearing client, so it takes the GenerateSecret side of the
	// suppression rule below.
	if customSecret != "" {
		if client.GenerateSecret {
			return nil, ErrInvalidParameter
		}
		if secretLen := len(customSecret); secretLen < clientSecretMinLength || secretLen > clientSecretMaxLength || !clientSecretPattern.MatchString(customSecret) {
			return nil, ErrInvalidParameter
		}
		client.ClientSecret = customSecret
		client.GenerateSecret = true
	}

	// A client created without GenerateSecret is a public client: the
	// stored record must carry no secret. This create-time rule belongs
	// to the Core so every transport creating a client gets it.
	s.suppressClientSecretCore(client)

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetUserPool(client.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.CreateUserPoolClient(client); err != nil {
		return nil, err
	}
	return client, nil
}

// updateUserPoolClientCore validates the assembled client configuration and
// persists the update.
func (s *CognitoService) updateUserPoolClientCore(region string, client *cognitostore.UserPoolClient) error {
	if !validateClientNamePattern(client.ClientName) {
		return ErrInvalidParameter
	}
	if err := validateUserPoolClientConfig(client); err != nil {
		return err
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.UpdateUserPoolClient(client); err != nil {
		return ErrInternalError
	}
	return nil
}

// validateUserPoolClientConfig is the single model-derived validation entry
// point for the whole client configuration, the counterpart of
// validateUserPoolConfig for user pools. It runs on the assembled
// store-level client, so a value is checked exactly once regardless of which
// transport assembled it.
func validateUserPoolClientConfig(client *cognitostore.UserPoolClient) error {
	if !validateRefreshTokenValidity(client.RefreshTokenValidity) {
		return ErrInvalidParameter
	}
	if !validateAccessTokenValidity(client.AccessTokenValidity) {
		return ErrInvalidParameter
	}
	if !validateIdTokenValidity(client.IDTokenValidity) {
		return ErrInvalidParameter
	}
	for _, f := range client.ExplicitAuthFlows {
		if !validateExplicitAuthFlow(f) {
			return ErrInvalidParameter
		}
	}
	// The machine-to-machine flow authorises an application rather than a
	// user, so it must be the only authentication flow the client carries.
	if slices.Contains(client.ExplicitAuthFlows, "ALLOW_CLIENT_TOKEN_AUTH") &&
		len(client.ExplicitAuthFlows) > 1 {
		return ErrInvalidParameter
	}
	for _, f := range client.AllowedOAuthFlows {
		if !validateOAuthFlow(f) {
			return ErrInvalidParameter
		}
	}
	if v := client.PreventUserExistenceErrors; v != "" {
		if !validatePreventUserExistenceErrors(v) {
			return ErrInvalidParameter
		}
	}
	if tvu := client.TokenValidityUnits; tvu != nil {
		for _, u := range []string{tvu.AccessToken, tvu.IdToken, tvu.RefreshToken} {
			if u != "" && !validateTimeUnit(u) {
				return ErrInvalidParameter
			}
		}
	}
	if rtr := client.RefreshTokenRotation; rtr != nil {
		if rtr.Feature != "" && !validateFeatureType(rtr.Feature) {
			return ErrInvalidParameter
		}
	}
	return nil
}

// newUserPoolClientCore validates the create-path required members and
// constructs the store object. Store constructors are store-package calls,
// so construction lives behind Core.
func (s *CognitoService) newUserPoolClientCore(userPoolID, clientName string) (*cognitostore.UserPoolClient, error) {
	if userPoolID == "" || clientName == "" {
		return nil, ErrInvalidParameter
	}
	return cognitostore.NewUserPoolClient(userPoolID, clientName), nil
}

// suppressClientSecretCore clears the client secret when the client was not
// configured to generate one, so the stored record carries no secret.
func (s *CognitoService) suppressClientSecretCore(client *cognitostore.UserPoolClient) {
	if !client.GenerateSecret {
		client.ClientSecret = ""
	}
}
