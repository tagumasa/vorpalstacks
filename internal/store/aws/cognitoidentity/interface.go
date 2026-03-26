package cognitoidentity

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
}

// Raw returns the underlying Cognito Identity store.
func (s *CognitoIdentityStore) Raw() *CognitoIdentityStore {
	return s
}

var _ CognitoIdentityStoreInterface = (*CognitoIdentityStore)(nil)
