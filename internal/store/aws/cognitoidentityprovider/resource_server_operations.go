package cognitoidentityprovider

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// ListResourceServersPaginated lists resource servers for a user pool with server-side pagination.
func (s *CognitoStore) ListResourceServersPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[ResourceServer], error) {
	opts.Prefix = resourceServerPrefix(userPoolID)
	return common.List[ResourceServer](s.BaseStore, opts, nil)
}

// CreateResourceServer persists a new resource server for a user pool.
func (s *CognitoStore) CreateResourceServer(rs *ResourceServer) error {
	key := resourceServerKey(rs.UserPoolID, rs.Identifier)
	if s.BaseStore.Exists(key) {
		return ErrResourceAlreadyExists
	}
	now := time.Now().UTC()
	rs.CreationDate = now
	rs.LastModifiedDate = now
	return s.BaseStore.Put(key, rs)
}

// GetResourceServer retrieves a resource server by user pool ID and identifier.
func (s *CognitoStore) GetResourceServer(userPoolID, identifier string) (*ResourceServer, error) {
	key := resourceServerKey(userPoolID, identifier)
	var rs ResourceServer
	if err := s.BaseStore.Get(key, &rs); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &rs, nil
}

// UpdateResourceServer updates an existing resource server in the store.
func (s *CognitoStore) UpdateResourceServer(rs *ResourceServer) error {
	key := resourceServerKey(rs.UserPoolID, rs.Identifier)
	rs.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(key, rs)
}

// DeleteResourceServer removes a resource server from the store by user pool ID and identifier.
func (s *CognitoStore) DeleteResourceServer(userPoolID, identifier string) error {
	key := resourceServerKey(userPoolID, identifier)
	if !s.BaseStore.Exists(key) {
		return ErrUserPoolNotFound
	}
	return s.BaseStore.Delete(key)
}

// ListResourceServers lists all resource servers for a user pool.
func (s *CognitoStore) ListResourceServers(userPoolID string) ([]*ResourceServer, error) {
	var servers []*ResourceServer
	prefix := resourceServerPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var rs ResourceServer
		if err := json.Unmarshal(value, &rs); err != nil {
			return err
		}
		servers = append(servers, &rs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// ListIdentityProvidersPaginated lists identity providers for a user pool with server-side pagination.
func (s *CognitoStore) ListIdentityProvidersPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[IdentityProvider], error) {
	opts.Prefix = identityProviderPrefix(userPoolID)
	return common.List[IdentityProvider](s.BaseStore, opts, nil)
}

// CreateIdentityProvider persists a new identity provider for a user pool.
func (s *CognitoStore) CreateIdentityProvider(ip *IdentityProvider) error {
	key := identityProviderKey(ip.UserPoolID, ip.ProviderName)
	if s.BaseStore.Exists(key) {
		return ErrResourceAlreadyExists
	}
	now := time.Now().UTC()
	ip.CreationDate = now
	ip.LastModifiedDate = now
	return s.BaseStore.Put(key, ip)
}

// GetIdentityProvider retrieves an identity provider by user pool ID and provider name.
func (s *CognitoStore) GetIdentityProvider(userPoolID, providerName string) (*IdentityProvider, error) {
	key := identityProviderKey(userPoolID, providerName)
	var ip IdentityProvider
	if err := s.BaseStore.Get(key, &ip); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &ip, nil
}

// UpdateIdentityProvider updates an existing identity provider in the store.
func (s *CognitoStore) UpdateIdentityProvider(ip *IdentityProvider) error {
	key := identityProviderKey(ip.UserPoolID, ip.ProviderName)
	ip.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(key, ip)
}

// DeleteIdentityProvider removes an identity provider from the store by user pool ID and provider name.
func (s *CognitoStore) DeleteIdentityProvider(userPoolID, providerName string) error {
	key := identityProviderKey(userPoolID, providerName)
	if !s.BaseStore.Exists(key) {
		return ErrUserPoolNotFound
	}
	return s.BaseStore.Delete(key)
}

// ListIdentityProviders lists all identity providers for a user pool.
func (s *CognitoStore) ListIdentityProviders(userPoolID string) ([]*IdentityProvider, error) {
	var providers []*IdentityProvider
	prefix := identityProviderPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var ip IdentityProvider
		if err := json.Unmarshal(value, &ip); err != nil {
			return err
		}
		providers = append(providers, &ip)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return providers, nil
}
