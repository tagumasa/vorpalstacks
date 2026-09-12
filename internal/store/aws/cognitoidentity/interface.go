package cognitoidentity

import (
	"vorpalstacks/internal/store/aws/common"
)

// CognitoIdentityStoreInterface defines operations for managing Cognito Identity pools.
type CognitoIdentityStoreInterface interface {
	CreateIdentityPool(pool *IdentityPool) (*IdentityPool, error)
	GetIdentityPool(id string) (*IdentityPool, error)
	UpdateIdentityPoolFunc(id string, mutate func(*IdentityPool) error) error
	DeleteIdentityPool(id string) error
	ListIdentityPools(opts common.ListOptions) (*common.ListResult[IdentityPool], error)
	CreateIdentity(identity *Identity) error
	GetIdentity(poolID, identityID string) (*Identity, error)
	DeleteIdentity(poolID, identityID string) error
	SetIdentityPoolRoles(poolID string, authRole, unauthRole string, mappings map[string]RoleMapping) error
	GetIdentityPoolRoles(poolID string) (authRole, unauthRole string, mappings map[string]RoleMapping, err error)
	GetIdentityByID(identityID string) (*Identity, error)
	GetOrCreateIdentityByLogins(poolID string, logins map[string]string) (*Identity, error)
	PutIdentity(identity *Identity) error
	Exists(id string) bool
	List(resourceKey string) (map[string]string, error)
	Tag(resourceKey string, tags map[string]string) error
	Untag(resourceKey string, tagKeys []string) error
	Replace(resourceKey string, tags map[string]string) error
	ListIdentitiesByPool(poolID string, maxResults int, nextToken string) ([]*Identity, string, error)
	UnlinkLogins(poolID, identityID string, loginsToRemove []string) error
	LinkLogins(poolID, identityID string, logins map[string]string) error
	MergeLogins(poolID, identityID string, logins map[string]string) (map[string]string, error)
	EnsureDeveloperIdentity(poolID, providerName, devUserID, suppliedIdentityID string) (string, error)
	MergeDeveloperIdentities(poolID, providerName, sourceUserID, destUserID string) (destIdentityID string, err error)
	LookupDeveloperIdentity(poolID string, identityID, devUserID string, maxResults int, nextToken string) (matchedIdentityID string, devUserIDs []string, nextTokenOut string, err error)
	UnlinkDeveloperIdentity(poolID, providerName, devUserID string) error
	GetDeveloperIdentity(poolID, providerName, devUserID string) (*DeveloperIdentity, error)
	SetPrincipalTagAttributeMap(poolID, providerName string, principalTags map[string]string, useDefaults bool) error
	GetPrincipalTagAttributeMap(poolID, providerName string) (*PrincipalTagAttributeMap, error)
}

var _ CognitoIdentityStoreInterface = (*CognitoIdentityStore)(nil)
