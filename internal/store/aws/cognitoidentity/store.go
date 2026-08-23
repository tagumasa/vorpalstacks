// Package cognitoidentity provides Cognito Identity Pool storage functionality for vorpalstacks.
package cognitoidentity

import (
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

// UpdateIdentityPool updates an existing Identity Pool, serialised with the
// pool's other mutations so a racing deletion cannot be overwritten and the
// pool record cannot be resurrected.
// Returns an error if the Identity Pool does not exist.
func (s *CognitoIdentityStore) UpdateIdentityPool(pool *IdentityPool) error {
	return s.keyLocker.WithLock(pool.ID, func() error {
		return s.updateIdentityPoolUnlocked(pool)
	})
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
// dropping the identity-ID index entries of the deleted identities in the
// same commit so the index never points at removed records.
func deletePoolIdentitiesBatch(store *common.BaseStore, prefix string) error {
	batchBucket, ok := store.Bucket().(storage.BatchBucket)
	if !ok {
		return errors.New("cognitoidentity: storage bucket does not support atomic batches")
	}
	batch := batchBucket.NewBatch()
	defer batch.Close()
	if err := store.ScanPrefix(prefix, func(key string, _ []byte) error {
		if err := batch.Delete([]byte(key)); err != nil {
			return err
		}
		return batch.Delete([]byte(identityIndexKey(strings.TrimPrefix(key, prefix))))
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
	return s.identitiesStore.Put(identityIndexKey(identity.ID), identity.IdentityPoolID)
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

// deleteIdentityUnlocked removes an Identity and its ID index entry; the
// caller must hold the pool lock.
func (s *CognitoIdentityStore) deleteIdentityUnlocked(poolID, identityID string) error {
	key := IdentityPoolIdentityKey(poolID, identityID)
	if !s.identitiesStore.Exists(key) {
		return ErrIdentityNotFound
	}
	if err := s.identitiesStore.Delete(key); err != nil {
		return err
	}
	return s.identitiesStore.Delete(identityIndexKey(identityID))
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

// FindIdentityByLogins scans identities in a pool and returns the first whose
// logins map matches every entry in the requested logins. Returns
// ErrIdentityNotFound when no match exists.
func (s *CognitoIdentityStore) FindIdentityByLogins(poolID string, logins map[string]string) (*Identity, error) {
	prefix := identityPoolPrefix(poolID)
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

// GetOrCreateIdentityByLogins returns the pool identity whose logins match the
// supplied set, creating a new identity when no match exists. The resolution
// runs under the pool's key lock so concurrent callers reuse the same
// identity instead of creating duplicates.
func (s *CognitoIdentityStore) GetOrCreateIdentityByLogins(poolID string, logins map[string]string) (*Identity, error) {
	var matched *Identity
	err := s.keyLocker.WithLock(poolID, func() error {
		if len(logins) > 0 {
			existing, err := s.FindIdentityByLogins(poolID, logins)
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
		return nil, ErrIdentityNotFound
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
// with the pool's other mutations.
func (s *CognitoIdentityStore) UnlinkLogins(poolID, identityID string, loginsToRemove []string) error {
	return s.keyLocker.WithLock(poolID, func() error {
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
	})
}

func developerIdentityKey(poolID, providerName, devUserID string) string {
	return poolID + keySeparator + providerName + keySeparator + devUserID
}

// LinkDeveloperIdentity creates or updates a mapping between a developer user identifier and an identity,
// serialised with the pool's other mutations.
func (s *CognitoIdentityStore) LinkDeveloperIdentity(di *DeveloperIdentity) error {
	return s.keyLocker.WithLock(di.IdentityPoolID, func() error {
		return s.linkDeveloperIdentityUnlocked(di)
	})
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
		if err := s.identitiesStore.Put(IdentityPoolIdentityKey(poolID, destDI.IdentityID), destIdentity); err != nil {
			return err
		}

		// The source identity has no remaining references once its logins are
		// merged and the developer identity link has moved; delete it last so
		// failures in the steps above leave it recoverable.
		return s.deleteIdentityUnlocked(poolID, sourceDI.IdentityID)
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
