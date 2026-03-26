// Package cognitoidentity provides Cognito Identity Pool storage functionality for vorpalstacks.
package cognitoidentity

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CognitoIdentityStore provides storage operations for Cognito Identity Pools and Identities.
type CognitoIdentityStore struct {
	*common.BaseStore
	identitiesStore *common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
}

func identityPoolBucketName(region string) string {
	return "cognito-identitypools-" + region
}

func identityBucketName(region string) string {
	return "cognito-identities-" + region
}

// NewCognitoIdentityStore creates a new CognitoIdentityStore instance.
func NewCognitoIdentityStore(store storage.BasicStorage, accountID, region string) *CognitoIdentityStore {
	return &CognitoIdentityStore{
		BaseStore:       common.NewBaseStore(store.Bucket(identityPoolBucketName(region)), "cognito-identitypools"),
		identitiesStore: common.NewBaseStore(store.Bucket(identityBucketName(region)), "cognito-identities"),
		TagStore:        common.NewTagStoreWithRegion(store, "cognito-identity", region),
		arnBuilder:      svcarn.NewARNBuilder(accountID, region),
		accountID:       accountID,
		region:          region,
	}
}

func (s *CognitoIdentityStore) buildIdentityPoolArn(poolID string) string {
	return s.arnBuilder.Cognito().IdentityPool(poolID)
}

// CreateIdentityPool creates a new Identity Pool in the store.
// Returns the created Identity Pool or an error if creation fails.
func (s *CognitoIdentityStore) CreateIdentityPool(pool *IdentityPool) (*IdentityPool, error) {
	if pool.Name == "" {
		return nil, ErrInvalidIdentityPoolName
	}

	if s.Exists(pool.ID) {
		return nil, ErrIdentityPoolAlreadyExists
	}

	now := time.Now().UTC()
	pool.Arn = s.buildIdentityPoolArn(pool.ID)
	pool.CreationDate = now
	pool.LastModifiedDate = now

	if pool.Tags == nil {
		pool.Tags = make(map[string]string)
	}
	if pool.CognitoIdentityProviders == nil {
		pool.CognitoIdentityProviders = []CognitoIdentityProvider{}
	}
	if pool.SupportedLoginProviders == nil {
		pool.SupportedLoginProviders = make(map[string]string)
	}
	if pool.OpenIdConnectProviderARNs == nil {
		pool.OpenIdConnectProviderARNs = []string{}
	}
	if pool.SamlProviderARNs == nil {
		pool.SamlProviderARNs = []string{}
	}
	if pool.RoleMappings == nil {
		pool.RoleMappings = make(map[string]RoleMapping)
	}

	if err := s.Put(pool.ID, pool); err != nil {
		return nil, err
	}

	return pool, nil
}

// GetIdentityPool retrieves an Identity Pool by its ID.
// Returns the Identity Pool or an error if not found.
func (s *CognitoIdentityStore) GetIdentityPool(id string) (*IdentityPool, error) {
	var pool IdentityPool
	if err := s.BaseStore.Get(id, &pool); err != nil {
		return nil, ErrIdentityPoolNotFound
	}
	return &pool, nil
}

// UpdateIdentityPool updates an existing Identity Pool.
// Returns an error if the Identity Pool does not exist.
func (s *CognitoIdentityStore) UpdateIdentityPool(pool *IdentityPool) error {
	if !s.Exists(pool.ID) {
		return ErrIdentityPoolNotFound
	}
	pool.LastModifiedDate = time.Now().UTC()
	return s.Put(pool.ID, pool)
}

// DeleteIdentityPool deletes an Identity Pool and all its associated identities.
// Returns an error if the Identity Pool does not exist.
func (s *CognitoIdentityStore) DeleteIdentityPool(id string) error {
	if !s.Exists(id) {
		return ErrIdentityPoolNotFound
	}

	prefix := id + "#"
	_ = s.identitiesStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.identitiesStore.Delete(key)
	})

	return s.BaseStore.Delete(id)
}

// ListIdentityPools returns a list of all Identity Pools.
// Returns the list of Identity Pools or an error if the operation fails.
func (s *CognitoIdentityStore) ListIdentityPools() ([]*IdentityPool, error) {
	var pools []*IdentityPool
	err := s.ForEach(func(key string, value []byte) error {
		var pool IdentityPool
		if err := json.Unmarshal(value, &pool); err != nil {
			return err
		}
		pools = append(pools, &pool)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pools, nil
}

// CreateIdentity creates a new Identity in the specified Identity Pool.
// Returns an error if the Identity Pool does not exist or the Identity already exists.
func (s *CognitoIdentityStore) CreateIdentity(identity *Identity) error {
	if identity.IdentityPoolID == "" {
		return ErrInvalidIdentityPoolID
	}

	if !s.Exists(identity.IdentityPoolID) {
		return ErrIdentityPoolNotFound
	}

	key := identityPoolIdentityKey(identity.IdentityPoolID, identity.ID)
	if s.identitiesStore.Exists(key) {
		return ErrIdentityAlreadyExists
	}

	now := time.Now().UTC()
	identity.CreationDate = now
	identity.LastModifiedDate = now

	if identity.Logins == nil {
		identity.Logins = make(map[string]string)
	}

	return s.identitiesStore.Put(key, identity)
}

// GetIdentity retrieves an Identity by its pool ID and identity ID.
// Returns the Identity or an error if not found.
func (s *CognitoIdentityStore) GetIdentity(poolID, identityID string) (*Identity, error) {
	key := identityPoolIdentityKey(poolID, identityID)
	var identity Identity
	if err := s.identitiesStore.Get(key, &identity); err != nil {
		return nil, ErrIdentityNotFound
	}
	return &identity, nil
}

// DeleteIdentity deletes an Identity from the store.
// Returns an error if the Identity does not exist.
func (s *CognitoIdentityStore) DeleteIdentity(poolID, identityID string) error {
	key := identityPoolIdentityKey(poolID, identityID)
	if !s.identitiesStore.Exists(key) {
		return ErrIdentityNotFound
	}
	return s.identitiesStore.Delete(key)
}

// SetIdentityPoolRoles sets the authentication and unauthentication roles for an Identity Pool.
//
// Parameters:
//   - poolID: The identity pool ID
//   - authRole: The authenticated role ARN
//   - unauthRole: The unauthenticated role ARN
//   - mappings: The role mappings
//
// Returns:
//   - error: An error if the operation fails
func (s *CognitoIdentityStore) SetIdentityPoolRoles(poolID string, authRole, unauthRole string, mappings map[string]RoleMapping) error {
	pool, err := s.GetIdentityPool(poolID)
	if err != nil {
		return err
	}

	pool.AuthenticatedRoleArn = authRole
	pool.UnauthenticatedRoleArn = unauthRole
	pool.RoleMappings = mappings

	return s.UpdateIdentityPool(pool)
}

// GetIdentityPoolRoles retrieves the roles configured for an Identity Pool.
func (s *CognitoIdentityStore) GetIdentityPoolRoles(poolID string) (authRole, unauthRole string, mappings map[string]RoleMapping, err error) {
	pool, err := s.GetIdentityPool(poolID)
	if err != nil {
		return "", "", nil, err
	}
	return pool.AuthenticatedRoleArn, pool.UnauthenticatedRoleArn, pool.RoleMappings, nil
}

func identityPoolIdentityKey(poolID, identityID string) string {
	return poolID + "#" + identityID
}

// GetIdentityByID retrieves a Cognito identity by its ID.
func (s *CognitoIdentityStore) GetIdentityByID(identityID string) (*Identity, error) {
	var foundIdentity *Identity
	err := s.identitiesStore.ScanPrefix("", func(key string, value []byte) error {
		var identity Identity
		if err := json.Unmarshal(value, &identity); err != nil {
			return err
		}
		if identity.ID == identityID {
			foundIdentity = &identity
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if foundIdentity == nil {
		return nil, ErrIdentityNotFound
	}
	return foundIdentity, nil
}
