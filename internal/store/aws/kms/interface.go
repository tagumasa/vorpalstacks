package kms

// KeyStoreInterface defines operations for managing KMS keys.
type KeyStoreInterface interface {
	Create(keyID string, keyUsage KeyUsage, keySpec KeySpec, description string, origin OriginType, multiRegion bool) (*Key, error)
	Get(keyID string) (*Key, error)
	Update(key *Key) error
	Delete(keyID string) error
	Exists(keyID string) bool
	List(marker string, maxItems int) (*KeyListResult, error)
	Enable(keyID string) error
	Disable(keyID string) error
	ScheduleDeletion(keyID string, pendingWindowInDays int) error
	CancelDeletion(keyID string) error
	// Raw returns the underlying key store.
	Raw() *KeyStore
}

// AliasStoreInterface defines operations for managing KMS key aliases.
type AliasStoreInterface interface {
	Create(aliasName, targetKeyID string) (*Alias, error)
	Get(aliasName string) (*Alias, error)
	Delete(aliasName string) error
	Exists(aliasName string) bool
	UpdateTarget(aliasName, targetKeyID string) error
	List(marker string, maxItems int, keyID string) (*AliasListResult, error)
	ListForKeyID(keyID string) ([]*Alias, error)
	Raw() *AliasStore
}

// GrantStoreInterface defines operations for managing KMS grants.
type GrantStoreInterface interface {
	Create(keyID, granteePrincipal, retiringPrincipal string, operations []string, name string, constraints *GrantConstraints) (*Grant, error)
	Get(grantID string) (*Grant, error)
	Delete(grantID string) error
	List(keyID string, granteePrincipal string, marker string, maxItems int) (*GrantListResult, error)
	ListByRetiringPrincipal(principal string, marker string, maxItems int) (*GrantListResult, error)
	DeleteByKeyID(keyID string) error
	Count(keyID string) (int, error)
	Raw() *GrantStore
}

// KeyPolicyStoreInterface defines operations for managing KMS key policies.
type KeyPolicyStoreInterface interface {
	Create(keyID, policyName, policyDocument string) (*KeyPolicy, error)
	Get(keyID, policyName string) (*KeyPolicy, error)
	Put(keyID, policyName, policyDocument string) error
	Delete(keyID, policyName string) error
	List(keyID string) ([]string, error)
	DeleteAllForKey(keyID string) error
	GetDefault(keyID string) (*KeyPolicy, error)
	PutDefault(keyID, policyDocument string) error
	Raw() *KeyPolicyStore
}

// KMSStoresInterface defines access to all KMS stores.
type KMSStoresInterface interface {
	KeyStore() KeyStoreInterface
	AliasStore() AliasStoreInterface
	GrantStore() GrantStoreInterface
	KeyPolicyStore() KeyPolicyStoreInterface
}

// KMSStores provides access to all KMS stores.
type KMSStores struct {
	keys        *KeyStore
	aliases     *AliasStore
	grants      *GrantStore
	keyPolicies *KeyPolicyStore
}

// NewKMSStores creates a new KMSStores with the given stores.
func NewKMSStores(keys *KeyStore, aliases *AliasStore, grants *GrantStore, keyPolicies *KeyPolicyStore) *KMSStores {
	return &KMSStores{
		keys:        keys,
		aliases:     aliases,
		grants:      grants,
		keyPolicies: keyPolicies,
	}
}

// KeyStore returns the key store.
func (s *KMSStores) KeyStore() KeyStoreInterface {
	if s.keys == nil {
		return nil
	}
	return s.keys
}

// AliasStore returns the alias store.
func (s *KMSStores) AliasStore() AliasStoreInterface {
	if s.aliases == nil {
		return nil
	}
	return s.aliases
}

// GrantStore returns the grant store.
func (s *KMSStores) GrantStore() GrantStoreInterface {
	if s.grants == nil {
		return nil
	}
	return s.grants
}

// KeyPolicyStore returns the key policy store.
func (s *KMSStores) KeyPolicyStore() KeyPolicyStoreInterface {
	if s.keyPolicies == nil {
		return nil
	}
	return s.keyPolicies
}

// Raw returns the underlying key store.
func (s *KeyStore) Raw() *KeyStore { return s }

// Raw returns the underlying alias store.
func (s *AliasStore) Raw() *AliasStore { return s }

// Raw returns the underlying grant store.
func (s *GrantStore) Raw() *GrantStore { return s }

// Raw returns the underlying key policy store.
func (s *KeyPolicyStore) Raw() *KeyPolicyStore { return s }
