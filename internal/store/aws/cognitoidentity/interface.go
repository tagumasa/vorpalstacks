package cognitoidentity

import (
	"vorpalstacks/internal/store/aws/common"
)

// CognitoIdentityStoreInterface defines operations for managing Cognito Identity pools.
type CognitoIdentityStoreInterface interface {
	CreateIdentityPool(pool *IdentityPool) (*IdentityPool, error)
	GetIdentityPool(id string) (*IdentityPool, error)
	UpdateIdentityPool(pool *IdentityPool) error
	DeleteIdentityPool(id string) error
	ListIdentityPools() ([]*IdentityPool, error)
	CreateIdentity(identity *Identity) error
	GetIdentity(poolID, identityID string) (*Identity, error)
	DeleteIdentity(poolID, identityID string) error
	SetIdentityPoolRoles(poolID string, authRole, unauthRole string, mappings map[string]RoleMapping) error
	GetIdentityPoolRoles(poolID string) (authRole, unauthRole string, mappings map[string]RoleMapping, err error)
	GetIdentityByID(identityID string) (*Identity, error)
	Exists(id string) bool
	ListTags(resourceKey string) (map[string]string, error)
	TagResource(resourceKey string, tags map[string]string) error
	UntagResource(resourceKey string, tagKeys []string) error
	Raw() *CognitoIdentityStore
	Identities() *common.BaseStore
	ListIdentitiesByPool(poolID string, maxResults int, nextToken string) ([]*Identity, string, error)
	DeleteIdentitiesBatch(poolID string, identityIDs []string) ([]string, error)
	UnlinkLogins(poolID, identityID string, loginsToRemove []string) error
	LinkDeveloperIdentity(di *DeveloperIdentity) error
	LookupDeveloperIdentity(poolID string, identityID, devUserID string, maxResults int) (string, []string, []string, error)
	UnlinkDeveloperIdentity(poolID, providerName, devUserID string) error
	GetDeveloperIdentity(poolID, providerName, devUserID string) (*DeveloperIdentity, error)
	SetPrincipalTagAttributeMap(poolID, providerName string, principalTags map[string]string, useDefaults bool) error
	GetPrincipalTagAttributeMap(poolID, providerName string) (*PrincipalTagAttributeMap, error)
}

// Raw returns the underlying Cognito Identity store.
func (s *CognitoIdentityStore) Raw() *CognitoIdentityStore {
	return s
}

// Identities returns the underlying identities store for direct access.
func (s *CognitoIdentityStore) Identities() *common.BaseStore {
	return s.identitiesStore
}

var _ CognitoIdentityStoreInterface = (*CognitoIdentityStore)(nil)
