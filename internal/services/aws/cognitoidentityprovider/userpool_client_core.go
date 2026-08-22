package cognitoidentityprovider

import (
	"errors"

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

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > listLimitMax {
		maxResults = listLimitMax
	}

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
		return nil, ErrClientNotFound
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
// pre-constructed store-level UserPoolClient so that the HTTP handler can
// apply its wire-specific parameter parsing before persistence. The admin
// handler builds the client via the Input DTO path in admin_handler_convert.go.
func (s *CognitoService) createUserPoolClientCore(region string, client *cognitostore.UserPoolClient) (*cognitostore.UserPoolClient, error) {
	if client.UserPoolID == "" || client.ClientName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateClientNamePattern(client.ClientName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetUserPool(client.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.CreateUserPoolClient(client); err != nil {
		if errors.Is(err, cognitostore.ErrClientAlreadyExists) {
			return nil, ErrClientAlreadyExists
		}
		return nil, err
	}
	return client, nil
}

// updateUserPoolClientCore persists updates to a user pool client.
func (s *CognitoService) updateUserPoolClientCore(region string, client *cognitostore.UserPoolClient) error {
	if !validateClientNamePattern(client.ClientName) {
		return ErrInvalidParameter
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
