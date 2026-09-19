// Package cognitoidentity provides Cognito Identity Pool storage functionality for vorpalstacks.
package cognitoidentity

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
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
	keyLocker  common.KeyLocker
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
}

// MaxIdentityPoolsPerAccount is the default AWS resource quota for identity
// pools per account ("Quotas in Amazon Cognito": 1,000, adjustable).
const MaxIdentityPoolsPerAccount = 1000

// identityPoolQuotaLockKey serialises quota-checked identity pool creation so
// concurrent creators cannot overshoot the per-account quota.
const identityPoolQuotaLockKey = "identity-pool-quota"

// keySeparator joins the segments of the composite keys in the identity,
// developer identity and principal tag buckets.
const keySeparator = "#"

// identityIndexKeyPrefix prefixes the identity-ID index entries. Pool IDs
// always contain a colon (the Smithy IdentityPoolId pattern is
// ^[\w-]+:[0-9a-f-]+$), so the colon-free prefix "idx" can never collide with
// a pool-prefixed identity key.
const identityIndexKeyPrefix = "idx" + keySeparator

// identityIndexKey maps an identity ID to the ID of the pool that stores it,
// allowing direct lookup without scanning the identity bucket.
func identityIndexKey(identityID string) string {
	return identityIndexKeyPrefix + identityID
}

// loginIndexKeyPrefix prefixes the login-claim index entries. Like the
// identity-ID index prefix it carries no colon, so it can never collide with
// a pool-prefixed identity key.
const loginIndexKeyPrefix = "loginidx" + keySeparator

// loginIndexKey maps a (provider, token) pair to the ID of the identity that
// claims it within a pool, letting the login-conflict check and the
// by-logins identity lookup resolve in one read instead of a sweep over
// every identity the pool holds. The pair is hashed: provider names and
// tokens may both contain the key separator.
func loginIndexKey(poolID, provider, token string) string {
	sum := sha256.Sum256([]byte(provider + "\x00" + token))
	return loginIndexKeyPrefix + poolID + keySeparator + hex.EncodeToString(sum[:])
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
		TagStore:          common.NewTagStoreWithRegion(store, "cognito-identity", region, common.StandardTagBudget("InvalidParameterException")),
		keyLocker:         common.KeyLocker{},
		arnBuilder:        svcarn.NewARNBuilder(accountID, region),
		accountID:         accountID,
		region:            region,
	}
}

func (s *CognitoIdentityStore) buildIdentityPoolArn(poolID string) string {
	return s.arnBuilder.Cognito().IdentityPool(poolID)
}

// CreateIdentityPool creates a new Identity Pool in the store. The per-account
// quota documented by AWS (default 1,000 pools) is enforced atomically with
// the creation itself.
// Returns the created Identity Pool or an error if creation fails.
func (s *CognitoIdentityStore) CreateIdentityPool(pool *IdentityPool) (*IdentityPool, error) {
	if pool.Name == "" {
		return nil, ErrInvalidIdentityPoolName
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

	// The quota count, the duplicate check and the write run under one lock so
	// concurrent creators cannot overshoot the per-account quota.
	if err := s.keyLocker.WithLock(identityPoolQuotaLockKey, func() error {
		count, err := s.countIdentityPools()
		if err != nil {
			return err
		}
		if count >= MaxIdentityPoolsPerAccount {
			return ErrTooManyIdentityPools
		}
		if s.Exists(pool.ID) {
			return ErrIdentityPoolAlreadyExists
		}
		return s.Put(pool.ID, pool)
	}); err != nil {
		return nil, err
	}

	return pool, nil
}

// countIdentityPools counts the pool records in the primary bucket without
// deserialising the values.
func (s *CognitoIdentityStore) countIdentityPools() (int, error) {
	count := 0
	err := s.BaseStore.ScanPrefix("", func(_ string, _ []byte) error {
		count++
		return nil
	})
	return count, err
}

// GetIdentityPool retrieves an Identity Pool by its ID.
// Returns the Identity Pool or an error if not found.
func (s *CognitoIdentityStore) GetIdentityPool(id string) (*IdentityPool, error) {
	var pool IdentityPool
	if err := s.BaseStore.Get(id, &pool); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrIdentityPoolNotFound
		}
		return nil, err
	}
	return &pool, nil
}

// updateIdentityPoolUnlocked stamps and persists an Identity Pool; the caller
// must hold the pool lock.
func (s *CognitoIdentityStore) updateIdentityPoolUnlocked(pool *IdentityPool) error {
	if !s.Exists(pool.ID) {
		return ErrIdentityPoolNotFound
	}
	pool.LastModifiedDate = time.Now().UTC()
	return s.Put(pool.ID, pool)
}

// UpdateIdentityPoolFunc reads the pool, applies mutate and persists the
// result as one serialised mutation under the pool's key lock: concurrent
// updates cannot lose each other's field changes, and a pool deleted before
// or during the mutation surfaces ErrIdentityPoolNotFound instead of being
// overwritten. A mutate error aborts without writing.
func (s *CognitoIdentityStore) UpdateIdentityPoolFunc(id string, mutate func(*IdentityPool) error) error {
	return s.keyLocker.WithLock(id, func() error {
		pool, err := s.GetIdentityPool(id)
		if err != nil {
			return err
		}
		if err := mutate(pool); err != nil {
			return err
		}
		return s.updateIdentityPoolUnlocked(pool)
	})
}

// DeleteIdentityPool deletes an Identity Pool and all its associated
// identities, developer identities, principal tag attribute maps and resource
// tags. The deletion runs under the pool lock, so pool content mutations
// (identity creation, developer identity links, principal tag maps, pool
// updates) cannot interleave and leave orphaned or resurrected records behind.
// Each bucket is cleared with one atomic batch and any failure aborts
// before the pool record is removed, so the deletion can be retried instead
// of silently leaving orphaned records behind.
// Returns an error if the Identity Pool does not exist.
func (s *CognitoIdentityStore) DeleteIdentityPool(id string) error {
	return s.keyLocker.WithLock(id, func() error {
		return s.deleteIdentityPoolUnlocked(id)
	})
}

// deleteIdentityPoolUnlocked cascades the pool deletion across every bucket;
// the caller must hold the pool lock.
func (s *CognitoIdentityStore) deleteIdentityPoolUnlocked(id string) error {
	if !s.Exists(id) {
		return ErrIdentityPoolNotFound
	}

	prefix := identityPoolPrefix(id)

	if err := deletePoolIdentitiesBatch(s.identitiesStore, prefix); err != nil {
		logs.Error("failed to delete identities during pool deletion",
			logs.String("poolId", id), logs.Err(err))
		return err
	}
	if err := deletePrefixBatch(s.developerIdStore, prefix); err != nil {
		logs.Error("failed to delete developer identities during pool deletion",
			logs.String("poolId", id), logs.Err(err))
		return err
	}
	if err := deletePrefixBatch(s.principalTagStore, prefix); err != nil {
		logs.Error("failed to delete principal tag attribute maps during pool deletion",
			logs.String("poolId", id), logs.Err(err))
		return err
	}

	if err := s.TagStore.Delete(s.arnBuilder.Cognito().IdentityPool(id)); err != nil {
		logs.Error("failed to delete resource tags during pool deletion",
			logs.String("poolId", id), logs.Err(err))
		return err
	}

	return s.BaseStore.Delete(id)
}

// deletePrefixBatch deletes every key under prefix in one atomic batch, so a
// partial failure commits nothing and the caller can retry the deletion.
func deletePrefixBatch(store *common.BaseStore, prefix string) error {
	batchBucket, ok := store.Bucket().(storage.BatchBucket)
	if !ok {
		return errors.New("cognitoidentity: storage bucket does not support atomic batches")
	}
	batch := batchBucket.NewBatch()
	defer batch.Close()
	if err := store.ScanPrefix(prefix, func(key string, _ []byte) error {
		return batch.Delete([]byte(key))
	}); err != nil {
		return err
	}
	return batch.Commit()
}

// deletePoolIdentitiesBatch deletes every pool identity in one atomic batch,
// dropping the identity-ID and login index entries of the deleted identities
// in the same commit so the indexes never point at removed records.
func deletePoolIdentitiesBatch(store *common.BaseStore, prefix string) error {
	batchBucket, ok := store.Bucket().(storage.BatchBucket)
	if !ok {
		return errors.New("cognitoidentity: storage bucket does not support atomic batches")
	}
	batch := batchBucket.NewBatch()
	defer batch.Close()
	if err := store.ScanPrefix(prefix, func(key string, value []byte) error {
		if err := batch.Delete([]byte(key)); err != nil {
			return err
		}
		if err := batch.Delete([]byte(identityIndexKey(strings.TrimPrefix(key, prefix)))); err != nil {
			return err
		}
		var identity Identity
		if err := json.Unmarshal(value, &identity); err != nil {
			return err
		}
		for provider, token := range identity.Logins {
			if token == "" {
				continue
			}
			if err := batch.Delete([]byte(loginIndexKey(identity.IdentityPoolID, provider, token))); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return batch.Commit()
}

// ListIdentityPools returns a list of Identity Pools with server-side pagination.
// Returns the list of Identity Pools or an error if the operation fails.
func (s *CognitoIdentityStore) ListIdentityPools(opts common.ListOptions) (*common.ListResult[IdentityPool], error) {
	return common.List[IdentityPool](s.BaseStore, opts, nil)
}

// CreateIdentity creates a new Identity in the specified Identity Pool.
// The write is serialised with the pool's other mutations through the pool
// lock, so it cannot interleave with a concurrent pool deletion.
// Returns an error if the Identity Pool does not exist or the Identity already exists.
func (s *CognitoIdentityStore) CreateIdentity(identity *Identity) error {
	if identity.IdentityPoolID == "" {
		return ErrInvalidIdentityPoolID
	}
	return s.keyLocker.WithLock(identity.IdentityPoolID, func() error {
		return s.createIdentityUnlocked(identity)
	})
}

// createIdentityUnlocked stores a new Identity and its ID index entry; the
// caller must hold the pool lock.
func (s *CognitoIdentityStore) createIdentityUnlocked(identity *Identity) error {
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

	if err := s.identitiesStore.Put(key, identity); err != nil {
		return err
	}
	// Maintain the identity-ID index so lookups by ID alone resolve directly.
	// A stale entry (identity deleted through the pool cascade) is harmless:
	// the indexed pool lookup then reports the identity as missing.
	if err := s.identitiesStore.Put(identityIndexKey(identity.ID), identity.IdentityPoolID); err != nil {
		return err
	}
	// Claim each login's index entry so the conflict check and the
	// by-logins lookup resolve without a pool sweep.
	for provider, token := range identity.Logins {
		if token == "" {
			continue
		}
		if err := s.identitiesStore.Put(loginIndexKey(identity.IdentityPoolID, provider, token), identity.ID); err != nil {
			return err
		}
	}
	return nil
}

// GetIdentity retrieves an Identity by its pool ID and identity ID.
// Returns the Identity or an error if not found.
func (s *CognitoIdentityStore) GetIdentity(poolID, identityID string) (*Identity, error) {
	key := IdentityPoolIdentityKey(poolID, identityID)
	var identity Identity
	if err := s.identitiesStore.Get(key, &identity); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrIdentityNotFound
		}
		return nil, err
	}
	return &identity, nil
}

// DeleteIdentity deletes an Identity from the store, serialised with the
// pool's other mutations.
// Returns an error if the Identity does not exist.
func (s *CognitoIdentityStore) DeleteIdentity(poolID, identityID string) error {
	return s.keyLocker.WithLock(poolID, func() error {
		return s.deleteIdentityUnlocked(poolID, identityID)
	})
}

// deleteIdentityUnlocked removes an Identity with its ID index and login
// index entries; the caller must hold the pool lock.
func (s *CognitoIdentityStore) deleteIdentityUnlocked(poolID, identityID string) error {
	key := IdentityPoolIdentityKey(poolID, identityID)
	identity, err := s.GetIdentity(poolID, identityID)
	if err != nil {
		return err
	}
	if err := s.identitiesStore.Delete(key); err != nil {
		return err
	}
	if err := s.identitiesStore.Delete(identityIndexKey(identityID)); err != nil {
		return err
	}
	for provider, token := range identity.Logins {
		if token == "" {
			continue
		}
		if err := s.identitiesStore.Delete(loginIndexKey(poolID, provider, token)); err != nil {
			return err
		}
	}
	return nil
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
	return s.keyLocker.WithLock(poolID, func() error {
		pool, err := s.GetIdentityPool(poolID)
		if err != nil {
			return err
		}

		pool.AuthenticatedRoleArn = authRole
		pool.UnauthenticatedRoleArn = unauthRole
		pool.RoleMappings = mappings

		return s.updateIdentityPoolUnlocked(pool)
	})
}

// GetIdentityPoolRoles retrieves the roles configured for an Identity Pool.
func (s *CognitoIdentityStore) GetIdentityPoolRoles(poolID string) (authRole, unauthRole string, mappings map[string]RoleMapping, err error) {
	pool, err := s.GetIdentityPool(poolID)
	if err != nil {
		return "", "", nil, err
	}
	return pool.AuthenticatedRoleArn, pool.UnauthenticatedRoleArn, pool.RoleMappings, nil
}

// identityPoolPrefix returns the key prefix shared by every record that
// belongs to the identity pool, in any of the per-pool buckets.
func identityPoolPrefix(poolID string) string {
	return poolID + keySeparator
}

// IdentityPoolIdentityKey returns the composite storage key for a Cognito identity within a pool.
func IdentityPoolIdentityKey(poolID, identityID string) string {
	return poolID + keySeparator + identityID
}

// findIdentityByLogins resolves the pool identity whose logins match every
// entry in the requested set through the login index: an identity holding
// all the requested (provider, token) pairs has claimed each pair's index
// entry, so every pair must resolve to the same identity. A pair with no
// index entry matches nothing — a matching identity would have claimed it.
// The resolved record is re-checked against the full set so a stale index
// entry can never manufacture a match. Returns ErrIdentityNotFound when no
// match exists.
func (s *CognitoIdentityStore) findIdentityByLogins(poolID string, logins map[string]string) (*Identity, error) {
	var matched *Identity
	for provider, token := range logins {
		if token == "" {
			continue
		}
		var identityID string
		if err := s.identitiesStore.Get(loginIndexKey(poolID, provider, token), &identityID); err != nil {
			if common.IsNotFound(err) {
				return nil, ErrIdentityNotFound
			}
			return nil, err
		}
		identity, err := s.GetIdentity(poolID, identityID)
		if err != nil {
			// A stale index entry naming a deleted identity reads as no
			// match — the create path re-points the index. Any other store
			// failure propagates, so a transient read cannot make
			// GetOrCreateIdentityByLogins mint a duplicate identity.
			if errors.Is(err, ErrIdentityNotFound) {
				return nil, ErrIdentityNotFound
			}
			return nil, err
		}
		if matched != nil && matched.ID != identity.ID {
			return nil, ErrIdentityNotFound
		}
		matched = identity
	}
	if matched == nil {
		return nil, ErrIdentityNotFound
	}
	for provider, token := range logins {
		if token != "" && matched.Logins[provider] != token {
			return nil, ErrIdentityNotFound
		}
	}
	return matched, nil
}

// GetOrCreateIdentityByLogins returns the pool identity whose logins match the
// supplied set, creating a new identity when no match exists. The resolution
// runs under the pool's key lock so concurrent callers reuse the same
// identity instead of creating duplicates.
func (s *CognitoIdentityStore) GetOrCreateIdentityByLogins(poolID string, logins map[string]string) (*Identity, error) {
	var matched *Identity
	err := s.keyLocker.WithLock(poolID, func() error {
		if len(logins) > 0 {
			existing, err := s.findIdentityByLogins(poolID, logins)
			if err == nil {
				matched = existing
				return nil
			}
			if !errors.Is(err, ErrIdentityNotFound) {
				return err
			}
		}
		identity := NewIdentity(poolID)
		if len(logins) > 0 {
			identity.Logins = logins
		}
		if err := s.createIdentityUnlocked(identity); err != nil {
			return err
		}
		matched = identity
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matched, nil
}

// GetIdentityByID retrieves an identity by its ID alone through the
// identity-ID index, resolving the owning pool without scanning the identity
// bucket.
func (s *CognitoIdentityStore) GetIdentityByID(identityID string) (*Identity, error) {
	var poolID string
	if err := s.identitiesStore.Get(identityIndexKey(identityID), &poolID); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrIdentityNotFound
		}
		return nil, err
	}
	return s.GetIdentity(poolID, identityID)
}

// ListIdentitiesByPool retrieves identities for a given pool with pagination support.
func (s *CognitoIdentityStore) ListIdentitiesByPool(poolID string, maxResults int, nextToken string) ([]*Identity, string, error) {
	result, err := common.List[Identity](s.identitiesStore, common.ListOptions{
		Prefix:   identityPoolPrefix(poolID),
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

// UnlinkLogins removes specified login providers from an identity, serialised
// with the pool's other mutations, releasing each removed provider's login
// index entry.
func (s *CognitoIdentityStore) UnlinkLogins(poolID, identityID string, loginsToRemove []string) error {
	return s.keyLocker.WithLock(poolID, func() error {
		key := IdentityPoolIdentityKey(poolID, identityID)
		var identity Identity
		if err := s.identitiesStore.Get(key, &identity); err != nil {
			if common.IsNotFound(err) {
				return ErrIdentityNotFound
			}
			return err
		}
		for _, login := range loginsToRemove {
			if old, held := identity.Logins[login]; held && old != "" {
				if err := s.identitiesStore.Delete(loginIndexKey(poolID, login, old)); err != nil {
					return err
				}
			}
			delete(identity.Logins, login)
		}
		identity.LastModifiedDate = time.Now().UTC()
		return s.identitiesStore.Put(key, identity)
	})
}

// LinkLogins attaches public-provider login entries to an identity,
// serialised with the pool's other mutations. A login that is already
// linked to a different identity of the pool is rejected with
// ErrLoginConflict — answered by the login index in one read per login,
// without sweeping the pool's identities; entries already carried by the
// identity are refreshed in place, with the superseded token's index entry
// released.
func (s *CognitoIdentityStore) LinkLogins(poolID, identityID string, logins map[string]string) error {
	return s.keyLocker.WithLock(poolID, func() error {
		key := IdentityPoolIdentityKey(poolID, identityID)
		var identity Identity
		if err := s.identitiesStore.Get(key, &identity); err != nil {
			if common.IsNotFound(err) {
				return ErrIdentityNotFound
			}
			return err
		}
		for provider, token := range logins {
			if token == "" {
				continue
			}
			var ownerID string
			if err := s.identitiesStore.Get(loginIndexKey(poolID, provider, token), &ownerID); err == nil {
				if ownerID != identityID {
					return ErrLoginConflict
				}
			} else if !common.IsNotFound(err) {
				return err
			}
		}
		if identity.Logins == nil {
			identity.Logins = make(map[string]string)
		}
		for provider, token := range logins {
			if old, held := identity.Logins[provider]; held && old != "" && old != token {
				if err := s.identitiesStore.Delete(loginIndexKey(poolID, provider, old)); err != nil {
					return err
				}
			}
			identity.Logins[provider] = token
			if token != "" {
				if err := s.identitiesStore.Put(loginIndexKey(poolID, provider, token), identityID); err != nil {
					return err
				}
			}
		}
		identity.LastModifiedDate = time.Now().UTC()
		return s.identitiesStore.Put(key, identity)
	})
}

// MergeLogins merges login entries into an identity's stored Logins map,
// reading and writing under the pool lock so concurrent credential fetches
// for one identity cannot lose each other's provider links. Entries
// already carried by the identity are refreshed in place. The merged map
// read inside the critical section is returned, so the caller resolves
// roles against the post-merge state rather than a pre-merge local copy.
func (s *CognitoIdentityStore) MergeLogins(poolID, identityID string, logins map[string]string) (map[string]string, error) {
	var merged map[string]string
	err := s.keyLocker.WithLock(poolID, func() error {
		key := IdentityPoolIdentityKey(poolID, identityID)
		var identity Identity
		if err := s.identitiesStore.Get(key, &identity); err != nil {
			if common.IsNotFound(err) {
				return ErrIdentityNotFound
			}
			return err
		}
		if identity.Logins == nil {
			identity.Logins = make(map[string]string)
		}
		for provider, token := range logins {
			if old, held := identity.Logins[provider]; held && old != "" && old != token {
				if err := s.identitiesStore.Delete(loginIndexKey(poolID, provider, old)); err != nil {
					return err
				}
			}
			identity.Logins[provider] = token
			if token != "" {
				if err := s.identitiesStore.Put(loginIndexKey(poolID, provider, token), identityID); err != nil {
					return err
				}
			}
		}
		identity.LastModifiedDate = time.Now().UTC()
		if err := s.identitiesStore.Put(key, identity); err != nil {
			return err
		}
		merged = identity.Logins
		return nil
	})
	return merged, err
}

func developerIdentityKey(poolID, providerName, devUserID string) string {
	return poolID + keySeparator + providerName + keySeparator + devUserID
}

// linkDeveloperIdentityUnlocked writes a developer identity mapping; the
// caller must hold the pool lock.
func (s *CognitoIdentityStore) linkDeveloperIdentityUnlocked(di *DeveloperIdentity) error {
	if !s.Exists(di.IdentityPoolID) {
		return ErrIdentityPoolNotFound
	}
	key := developerIdentityKey(di.IdentityPoolID, di.DeveloperProviderName, di.DeveloperUserIdentifier)
	return s.developerIdStore.Put(key, di)
}

// PutIdentity persists an identity under its canonical composite key,
// serialised with the pool's other mutations. The identity ID and its pool
// association never change, so the identity-ID index needs no maintenance
// here.
func (s *CognitoIdentityStore) PutIdentity(identity *Identity) error {
	return s.keyLocker.WithLock(identity.IdentityPoolID, func() error {
		key := IdentityPoolIdentityKey(identity.IdentityPoolID, identity.ID)
		return s.identitiesStore.Put(key, identity)
	})
}

// EnsureDeveloperIdentity resolves the identity linked to the developer user,
// creating and linking a fresh identity when no link exists yet. The whole
// resolution runs under the pool lock, so it cannot interleave with a
// concurrent pool deletion and concurrent callers for one developer user
// serialise onto a single identity. A supplied identity ID must either match
// the existing link or reference an existing identity of the pool.
func (s *CognitoIdentityStore) EnsureDeveloperIdentity(poolID, providerName, devUserID, suppliedIdentityID string) (string, error) {
	identityID := suppliedIdentityID
	err := s.keyLocker.WithLock(poolID, func() error {
		existing, err := s.GetDeveloperIdentity(poolID, providerName, devUserID)
		if err == nil {
			if suppliedIdentityID != "" && existing.IdentityID != suppliedIdentityID {
				return ErrDeveloperIdentityConflict
			}
			identityID = existing.IdentityID
			return nil
		}
		if !errors.Is(err, ErrIdentityNotFound) {
			return err
		}
		if suppliedIdentityID != "" {
			if _, err := s.GetIdentity(poolID, suppliedIdentityID); err != nil {
				return err
			}
		} else {
			identity := NewIdentity(poolID)
			if err := s.createIdentityUnlocked(identity); err != nil {
				return err
			}
			identityID = identity.ID
		}
		return s.linkDeveloperIdentityUnlocked(&DeveloperIdentity{
			DeveloperUserIdentifier: devUserID,
			DeveloperProviderName:   providerName,
			IdentityPoolID:          poolID,
			IdentityID:              identityID,
		})
	})
	if err != nil {
		return "", err
	}
	return identityID, nil
}

// MergeDeveloperIdentities moves a developer user's link to the destination
// user's identity, merges the source identity's logins into it and removes
// the source identity. The whole merge runs under the pool lock, and the
// link moves before any identity record is destroyed so a failure at any
// step leaves every developer identity referencing a live identity.
func (s *CognitoIdentityStore) MergeDeveloperIdentities(poolID, providerName, sourceUserID, destUserID string) (string, error) {
	return s.mergeDeveloperIdentities(poolID, providerName, sourceUserID, destUserID, s.linkDeveloperIdentityUnlocked)
}

// mergeDeveloperIdentities is the seam behind MergeDeveloperIdentities: the
// link step is a parameter so tests can inject its failure and pin the
// ordering invariant.
func (s *CognitoIdentityStore) mergeDeveloperIdentities(poolID, providerName, sourceUserID, destUserID string, link func(*DeveloperIdentity) error) (string, error) {
	destIdentityID := ""
	err := s.keyLocker.WithLock(poolID, func() error {
		sourceDI, err := s.GetDeveloperIdentity(poolID, providerName, sourceUserID)
		if err != nil {
			return err
		}
		destDI, err := s.GetDeveloperIdentity(poolID, providerName, destUserID)
		if err != nil {
			return err
		}
		destIdentityID = destDI.IdentityID

		// Move the developer identity link first: it is the authoritative
		// association, and until it succeeds no identity record is destroyed,
		// so a failure at any later step cannot leave a developer identity
		// referencing a deleted identity.
		if err := link(&DeveloperIdentity{
			DeveloperUserIdentifier: sourceUserID,
			DeveloperProviderName:   providerName,
			IdentityPoolID:          poolID,
			IdentityID:              destDI.IdentityID,
		}); err != nil {
			return err
		}

		if sourceDI.IdentityID == "" || destDI.IdentityID == "" || sourceDI.IdentityID == destDI.IdentityID {
			return nil
		}

		// Merge the source identity's logins into the destination identity so
		// that public provider links (Facebook, Google, etc.) are not lost.
		sourceIdentity, err := s.GetIdentity(poolID, sourceDI.IdentityID)
		if err != nil {
			return err
		}
		destIdentity, err := s.GetIdentity(poolID, destDI.IdentityID)
		if err != nil {
			return err
		}
		if destIdentity.Logins == nil {
			destIdentity.Logins = make(map[string]string)
		}
		for provider, token := range sourceIdentity.Logins {
			if _, exists := destIdentity.Logins[provider]; !exists {
				destIdentity.Logins[provider] = token
			}
		}
		destIdentity.LastModifiedDate = time.Now().UTC()

		// The source identity has no remaining references once its logins are
		// merged and the developer identity link has moved. The destination
		// write and the source's removal (record plus identity-ID index)
		// commit as one atomic batch: a mid-sequence failure commits
		// neither, so a merge can never destroy a source identity whose
		// logins were not handed over. The login index entries of the
		// handed-over logins re-point at the destination in the same commit.
		destKey := IdentityPoolIdentityKey(poolID, destDI.IdentityID)
		sourceKey := IdentityPoolIdentityKey(poolID, sourceDI.IdentityID)
		destValue, err := json.Marshal(destIdentity)
		if err != nil {
			return err
		}
		batchBucket, ok := s.identitiesStore.Bucket().(storage.BatchBucket)
		if !ok {
			return errors.New("cognitoidentity: storage bucket does not support atomic batches")
		}
		batch := batchBucket.NewBatch()
		defer batch.Close()
		if err := batch.Put([]byte(destKey), destValue); err != nil {
			return err
		}
		if err := batch.Delete([]byte(sourceKey)); err != nil {
			return err
		}
		if err := batch.Delete([]byte(identityIndexKey(sourceDI.IdentityID))); err != nil {
			return err
		}
		for provider, token := range sourceIdentity.Logins {
			if token == "" {
				continue
			}
			if held, exists := destIdentity.Logins[provider]; exists && held != token {
				// The destination keeps its own token for the provider, so
				// the source's token is discarded with the source identity —
				// its index entry goes in the same commit, or it would name
				// a deleted identity forever.
				if err := batch.Delete([]byte(loginIndexKey(poolID, provider, token))); err != nil {
					return err
				}
				continue
			}
			idxValue, merr := json.Marshal(destDI.IdentityID)
			if merr != nil {
				return merr
			}
			if err := batch.Put([]byte(loginIndexKey(poolID, provider, token)), idxValue); err != nil {
				return err
			}
		}
		return batch.Commit()
	})
	if err != nil {
		return "", err
	}
	return destIdentityID, nil
}

// LookupDeveloperIdentity looks up developer identity mappings with pagination support.
func (s *CognitoIdentityStore) LookupDeveloperIdentity(poolID string, identityID, devUserID string, maxResults int, nextToken string) (matchedIdentityID string, devUserIDs []string, nextTokenOut string, err error) {
	filter := func(di *DeveloperIdentity) bool {
		if devUserID != "" && di.DeveloperUserIdentifier != devUserID {
			return false
		}
		if identityID != "" && di.IdentityID != identityID {
			return false
		}
		return true
	}

	result, err := common.List[DeveloperIdentity](s.developerIdStore, common.ListOptions{
		Prefix:   identityPoolPrefix(poolID),
		Marker:   nextToken,
		MaxItems: maxResults,
	}, filter)
	if err != nil {
		return "", nil, "", err
	}

	devUserIDs = make([]string, 0, len(result.Items))
	for _, di := range result.Items {
		devUserIDs = append(devUserIDs, di.DeveloperUserIdentifier)
		if di.IdentityID != "" {
			matchedIdentityID = di.IdentityID
		}
	}

	if result.IsTruncated {
		nextTokenOut = result.NextMarker
	}

	return matchedIdentityID, devUserIDs, nextTokenOut, nil
}

// UnlinkDeveloperIdentity removes a developer identity mapping, serialised
// with the pool's other mutations.
func (s *CognitoIdentityStore) UnlinkDeveloperIdentity(poolID, providerName, devUserID string) error {
	return s.keyLocker.WithLock(poolID, func() error {
		key := developerIdentityKey(poolID, providerName, devUserID)
		if !s.developerIdStore.Exists(key) {
			return ErrIdentityNotFound
		}
		return s.developerIdStore.Delete(key)
	})
}

// GetDeveloperIdentity looks up a specific developer identity mapping.
func (s *CognitoIdentityStore) GetDeveloperIdentity(poolID, providerName, devUserID string) (*DeveloperIdentity, error) {
	key := developerIdentityKey(poolID, providerName, devUserID)
	var di DeveloperIdentity
	if err := s.developerIdStore.Get(key, &di); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrIdentityNotFound
		}
		return nil, err
	}
	return &di, nil
}

func principalTagKey(poolID, providerName string) string {
	return poolID + keySeparator + providerName
}

// SetPrincipalTagAttributeMap stores the principal tag attribute mapping for an identity provider,
// serialised with the pool's other mutations so it cannot orphan a record
// under a pool that is being deleted.
func (s *CognitoIdentityStore) SetPrincipalTagAttributeMap(poolID, providerName string, principalTags map[string]string, useDefaults bool) error {
	return s.keyLocker.WithLock(poolID, func() error {
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
	})
}

// GetPrincipalTagAttributeMap retrieves the principal tag attribute mapping for an identity provider.
func (s *CognitoIdentityStore) GetPrincipalTagAttributeMap(poolID, providerName string) (*PrincipalTagAttributeMap, error) {
	key := principalTagKey(poolID, providerName)
	var ptam PrincipalTagAttributeMap
	if err := s.principalTagStore.Get(key, &ptam); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrIdentityNotFound
		}
		return nil, err
	}
	if ptam.PrincipalTags == nil {
		ptam.PrincipalTags = make(map[string]string)
	}
	return &ptam, nil
}
