package cognitoidentityprovider

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// ListUserPoolClientsPaginated lists clients for a Cognito user pool with server-side pagination.
func (s *CognitoStore) ListUserPoolClientsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserPoolClient], error) {
	opts.Prefix = userPoolID + "#"
	return common.List[UserPoolClient](s.clientsStore, opts, nil)
}

// CreateUserPoolClient creates a new Cognito user pool client.
func (s *CognitoStore) CreateUserPoolClient(client *UserPoolClient) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if client.ClientName == "" {
		return ErrInvalidParameter
	}

	key := userPoolClientKey(client.UserPoolID, client.ClientID)
	if s.clientsStore.Exists(key) {
		return ErrClientAlreadyExists
	}

	now := time.Now().UTC()
	client.CreationDate = now
	client.LastModifiedDate = now

	if err := s.clientsStore.Put(key, client); err != nil {
		return err
	}
	// Write secondary index for O(1) GetUserPoolByClientID lookup.
	_ = s.clientsStore.Put(clientIndexKey(client.ClientID), client.UserPoolID)
	return nil
}

// GetUserPoolClient retrieves a Cognito user pool client by client ID.
func (s *CognitoStore) GetUserPoolClient(userPoolID, clientID string) (*UserPoolClient, error) {
	key := userPoolClientKey(userPoolID, clientID)
	var client UserPoolClient
	if err := s.clientsStore.Get(key, &client); err != nil {
		return nil, ErrClientNotFound
	}
	return &client, nil
}

// GetUserPoolClientByName retrieves a Cognito user pool client by client name.
func (s *CognitoStore) GetUserPoolClientByName(userPoolID, clientName string) (*UserPoolClient, error) {
	var found *UserPoolClient
	prefix := userPoolID + "#"
	err := s.clientsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var client UserPoolClient
		if err := json.Unmarshal(value, &client); err != nil {
			return err
		}
		if client.ClientName == clientName {
			found = &client
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrClientNotFound
	}
	return found, nil
}

// UpdateUserPoolClient updates an existing Cognito user pool client.
func (s *CognitoStore) UpdateUserPoolClient(client *UserPoolClient) error {
	key := userPoolClientKey(client.UserPoolID, client.ClientID)
	if !s.clientsStore.Exists(key) {
		return ErrClientNotFound
	}
	client.LastModifiedDate = time.Now().UTC()
	return s.clientsStore.Put(key, client)
}

// DeleteUserPoolClient deletes a Cognito user pool client.
func (s *CognitoStore) DeleteUserPoolClient(userPoolID, clientID string) error {
	key := userPoolClientKey(userPoolID, clientID)
	if !s.clientsStore.Exists(key) {
		return ErrClientNotFound
	}
	// Clean up secondary index.
	_ = s.clientsStore.Delete(clientIndexKey(clientID))
	return s.clientsStore.Delete(key)
}

// ListUserPoolClients lists all clients for a Cognito user pool.
func (s *CognitoStore) ListUserPoolClients(userPoolID string) ([]*UserPoolClient, error) {
	var clients []*UserPoolClient
	prefix := userPoolID + "#"
	err := s.clientsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var client UserPoolClient
		if err := json.Unmarshal(value, &client); err != nil {
			return err
		}
		clients = append(clients, &client)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// GetUserPoolByClientID retrieves the user pool associated with a client ID.
func (s *CognitoStore) GetUserPoolByClientID(clientID string) (*UserPool, error) {
	// Use secondary index for O(1) lookup.
	var poolID string
	if err := s.clientsStore.Get(clientIndexKey(clientID), &poolID); err == nil && poolID != "" {
		return s.GetUserPool(poolID)
	}
	// Fallback: full scan for legacy clients without index
	pools, err := s.ListUserPools()
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		clients, err := s.ListUserPoolClients(pool.ID)
		if err != nil {
			continue
		}
		for _, client := range clients {
			if client.ClientID == clientID {
				// Backfill index
				_ = s.clientsStore.Put(clientIndexKey(clientID), pool.ID)
				return pool, nil
			}
		}
	}

	return nil, ErrClientNotFound
}
