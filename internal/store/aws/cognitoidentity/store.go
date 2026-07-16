// Package cognitoidentity provides Cognito Identity Pool storage functionality for vorpalstacks.
package cognitoidentity

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CognitoIdentityStore provides storage operations for Cognito Identity Pools and Identities.
type CognitoIdentityStore struct {
	*common.BaseStore
	identitiesStore   *common.BaseStore
	developerIdStore  *common.BaseStore
	principalTagStore *common.BaseStore
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

func developerIdBucketName(region string) string {
	return "cognito-developerids-" + region
}

func principalTagBucketName(region string) string {
	return "cognito-principaltags-" + region
}

// NewCognitoIdentityStore creates a new CognitoIdentityStore instance.
func NewCognitoIdentityStore(store storage.BasicStorage, accountID, region string) *CognitoIdentityStore {
	return &CognitoIdentityStore{
		BaseStore:         common.NewBaseStore(store.Bucket(identityPoolBucketName(region)), "cognito-identitypools"),
		identitiesStore:   common.NewBaseStore(store.Bucket(identityBucketName(region)), "cognito-identities"),
		developerIdStore:  common.NewBaseStore(store.Bucket(developerIdBucketName(region)), "cognito-developerids"),
		principalTagStore: common.NewBaseStore(store.Bucket(principalTagBucketName(region)), "cognito-principaltags"),
		TagStore:          common.NewTagStoreWithRegion(store, "cognito-identity", region),
		arnBuilder:        svcarn.NewARNBuilder(accountID, region),
		accountID:         accountID,
		region:            region,
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

	if err := s.identitiesStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if err := s.identitiesStore.Delete(key); err != nil {
			logs.Warn("failed to delete identity during pool deletion",
				logs.String("identityKey", key), logs.String("poolId", id), logs.Err(err))
		}
		return nil
	}); err != nil {
		logs.Warn("failed to scan identities during pool deletion",
			logs.String("poolId", id), logs.Err(err))
	}

	if err := s.developerIdStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if err := s.developerIdStore.Delete(key); err != nil {
			logs.Warn("failed to delete developer identity during pool deletion",
				logs.String("devKey", key), logs.String("poolId", id), logs.Err(err))
		}
		return nil
	}); err != nil {
		logs.Warn("failed to scan developer identities during pool deletion",
			logs.String("poolId", id), logs.Err(err))
	}

	if err := s.principalTagStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if err := s.principalTagStore.Delete(key); err != nil {
			logs.Warn("failed to delete principal tag attribute map during pool deletion",
				logs.String("ptKey", key), logs.String("poolId", id), logs.Err(err))
		}
		return nil
	}); err != nil {
		logs.Warn("failed to scan principal tags during pool deletion",
			logs.String("poolId", id), logs.Err(err))
	}

	if err := s.TagStore.Delete(s.arnBuilder.Cognito().IdentityPool(id)); err != nil {
		logs.Warn("failed to delete resource tags during pool deletion",
			logs.String("poolId", id), logs.Err(err))
	}

	return s.BaseStore.Delete(id)
}

// ListIdentityPools returns a list of Identity Pools with server-side pagination.
// Returns the list of Identity Pools or an error if the operation fails.
func (s *CognitoIdentityStore) ListIdentityPools(opts common.ListOptions) (*common.ListResult[IdentityPool], error) {
	return common.List[IdentityPool](s.BaseStore, opts, nil)
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

	key := IdentityPoolIdentityKey(identity.IdentityPoolID, identity.ID)
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
	key := IdentityPoolIdentityKey(poolID, identityID)
	var identity Identity
	if err := s.identitiesStore.Get(key, &identity); err != nil {
		return nil, ErrIdentityNotFound
	}
	return &identity, nil
}

// DeleteIdentity deletes an Identity from the store.
// Returns an error if the Identity does not exist.
func (s *CognitoIdentityStore) DeleteIdentity(poolID, identityID string) error {
	key := IdentityPoolIdentityKey(poolID, identityID)
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

// IdentityPoolIdentityKey returns the composite storage key for a Cognito identity within a pool.
func IdentityPoolIdentityKey(poolID, identityID string) string {
	return poolID + "#" + identityID
}

// FindIdentityByLogins scans identities in a pool and returns the first whose
// logins map matches every entry in the requested logins. Returns
// ErrIdentityNotFound when no match exists.
func (s *CognitoIdentityStore) FindIdentityByLogins(poolID string, logins map[string]string) (*Identity, error) {
	prefix := poolID + "#"
	var found *Identity

	err := s.identitiesStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var identity Identity
		if err := json.Unmarshal(value, &identity); err != nil {
			return err
		}
		for provider, token := range logins {
			if identity.Logins[provider] != token {
				return nil
			}
		}
		found = &identity
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrIdentityNotFound
	}
	return found, nil
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

// ListIdentitiesByPool retrieves identities for a given pool with pagination support.
func (s *CognitoIdentityStore) ListIdentitiesByPool(poolID string, maxResults int, nextToken string) ([]*Identity, string, error) {
	result, err := common.List[Identity](s.identitiesStore, common.ListOptions{
		Prefix:   poolID + "#",
		Marker:   nextToken,
		MaxItems: maxResults,
	}, nil)
	if err != nil {
		return nil, "", err
	}

	var token string
	if result.IsTruncated {
		token = result.NextMarker
	}

	return result.Items, token, nil
}

// DeleteIdentitiesBatch deletes multiple identities and returns any that could not be deleted.
func (s *CognitoIdentityStore) DeleteIdentitiesBatch(poolID string, identityIDs []string) ([]string, error) {
	var unprocessed []string
	for _, id := range identityIDs {
		key := IdentityPoolIdentityKey(poolID, id)
		if !s.identitiesStore.Exists(key) {
			unprocessed = append(unprocessed, id)
			continue
		}
		if err := s.identitiesStore.Delete(key); err != nil {
			unprocessed = append(unprocessed, id)
		}
	}
	return unprocessed, nil
}

// UnlinkLogins removes specified login providers from an identity.
func (s *CognitoIdentityStore) UnlinkLogins(poolID, identityID string, loginsToRemove []string) error {
	key := IdentityPoolIdentityKey(poolID, identityID)
	var identity Identity
	if err := s.identitiesStore.Get(key, &identity); err != nil {
		return ErrIdentityNotFound
	}
	for _, login := range loginsToRemove {
		delete(identity.Logins, login)
	}
	identity.LastModifiedDate = time.Now().UTC()
	return s.identitiesStore.Put(key, identity)
}

func developerIdentityKey(poolID, providerName, devUserID string) string {
	return poolID + "#" + providerName + "#" + devUserID
}

// LinkDeveloperIdentity creates or updates a mapping between a developer user identifier and an identity.
func (s *CognitoIdentityStore) LinkDeveloperIdentity(di *DeveloperIdentity) error {
	if !s.Exists(di.IdentityPoolID) {
		return ErrIdentityPoolNotFound
	}
	key := developerIdentityKey(di.IdentityPoolID, di.DeveloperProviderName, di.DeveloperUserIdentifier)
	return s.developerIdStore.Put(key, di)
}

// LookupDeveloperIdentity looks up developer identity mappings with pagination support.
func (s *CognitoIdentityStore) LookupDeveloperIdentity(poolID string, identityID, devUserID string, maxResults int, nextToken string) (matchedIdentityID string, devUserIDs []string, nextTokenOut string, err error) {
	prefix := poolID + "#"

	type entry struct {
		key        string
		devUserID  string
		identityID string
	}
	var entries []entry

	scanErr := s.developerIdStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var di DeveloperIdentity
		if err := json.Unmarshal(value, &di); err != nil {
			return err
		}
		if devUserID != "" && di.DeveloperUserIdentifier != devUserID {
			return nil
		}
		if identityID != "" && di.IdentityID != identityID {
			return nil
		}
		entries = append(entries, entry{key: key, devUserID: di.DeveloperUserIdentifier, identityID: di.IdentityID})
		return nil
	})
	if scanErr != nil {
		return "", nil, "", scanErr
	}

	startIdx := 0
	if nextToken != "" {
		for i, e := range entries {
			if e.key == nextToken {
				startIdx = i + 1
				break
			}
		}
	}

	end := startIdx + maxResults
	if end > len(entries) {
		end = len(entries)
	}

	for _, e := range entries[startIdx:end] {
		devUserIDs = append(devUserIDs, e.devUserID)
		if e.identityID != "" {
			matchedIdentityID = e.identityID
		}
	}

	if end < len(entries) {
		nextTokenOut = entries[end-1].key
	}

	return matchedIdentityID, devUserIDs, nextTokenOut, nil
}

// UnlinkDeveloperIdentity removes a developer identity mapping.
func (s *CognitoIdentityStore) UnlinkDeveloperIdentity(poolID, providerName, devUserID string) error {
	key := developerIdentityKey(poolID, providerName, devUserID)
	if !s.developerIdStore.Exists(key) {
		return ErrIdentityNotFound
	}
	return s.developerIdStore.Delete(key)
}

// GetDeveloperIdentity looks up a specific developer identity mapping.
func (s *CognitoIdentityStore) GetDeveloperIdentity(poolID, providerName, devUserID string) (*DeveloperIdentity, error) {
	key := developerIdentityKey(poolID, providerName, devUserID)
	var di DeveloperIdentity
	if err := s.developerIdStore.Get(key, &di); err != nil {
		return nil, ErrIdentityNotFound
	}
	return &di, nil
}

func principalTagKey(poolID, providerName string) string {
	return poolID + "#" + providerName
}

// SetPrincipalTagAttributeMap stores the principal tag attribute mapping for an identity provider.
func (s *CognitoIdentityStore) SetPrincipalTagAttributeMap(poolID, providerName string, principalTags map[string]string, useDefaults bool) error {
	if !s.Exists(poolID) {
		return ErrIdentityPoolNotFound
	}
	ptam := &PrincipalTagAttributeMap{
		IdentityPoolID:       poolID,
		IdentityProviderName: providerName,
		PrincipalTags:        principalTags,
		UseDefaults:          useDefaults,
	}
	if ptam.PrincipalTags == nil {
		ptam.PrincipalTags = make(map[string]string)
	}
	key := principalTagKey(poolID, providerName)
	return s.principalTagStore.Put(key, ptam)
}

// GetPrincipalTagAttributeMap retrieves the principal tag attribute mapping for an identity provider.
func (s *CognitoIdentityStore) GetPrincipalTagAttributeMap(poolID, providerName string) (*PrincipalTagAttributeMap, error) {
	key := principalTagKey(poolID, providerName)
	var ptam PrincipalTagAttributeMap
	if err := s.principalTagStore.Get(key, &ptam); err != nil {
		return nil, ErrIdentityNotFound
	}
	if ptam.PrincipalTags == nil {
		ptam.PrincipalTags = make(map[string]string)
	}
	return &ptam, nil
}
