package cognitoidentityprovider

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/pkg/vsjwt"
)

// CreateUserPool creates a new Cognito user pool.
func (s *CognitoStore) CreateUserPool(userPool *UserPool) (*UserPool, error) {
	// Generate RSA key outside recordMu to avoid blocking CreateUserPool
	// and CreateGroup operations during the 100-300ms key generation.
	privateKey, err := vsjwt.GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	s.recordMu.Lock()
	defer s.recordMu.Unlock()
	if userPool.Name == "" {
		return nil, ErrInvalidUserPoolName
	}

	if s.Exists(userPool.ID) {
		return nil, ErrUserPoolAlreadyExists
	}

	now := time.Now().UTC()
	userPool.Arn = s.buildUserPoolArn(userPool.ID)
	userPool.CreationDate = now
	userPool.LastModifiedDate = now
	userPool.JwtKeyID = uuid.New().String()[:keyIDPrefixLength]
	userPool.JwtPrivateKey = vsjwt.EncodePrivateKeyToPEM(privateKey)
	// An encode failure aborts pool creation: a pool persisted with an
	// empty public key would pass creation and surface the corruption only
	// at the first token validation.
	publicKeyPEM, err := vsjwt.EncodePublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	userPool.JwtPublicKey = publicKeyPEM
	// The model's StatusType enum is Enabled | Disabled; a pool is created
	// enabled and no operation transitions it (the wire member is retained
	// but marked as no longer used).
	if userPool.Status == "" {
		userPool.Status = UserPoolStatusEnabled
	}

	if err := s.Put(userPool.ID, userPool); err != nil {
		return nil, err
	}

	return userPool, nil
}

// GetUserPool retrieves a Cognito user pool by ID.
func (s *CognitoStore) GetUserPool(userPoolID string) (*UserPool, error) {
	var userPool UserPool
	if err := s.BaseStore.Get(userPoolID, &userPool); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrUserPoolNotFound
		}
		return nil, err
	}
	return &userPool, nil
}

// updateUserPoolUnlocked stamps and persists a user-pool record; the caller
// must hold the pool lock. The write also refreshes the cached username-case
// sensitivity so pool mutations take effect on the next username-keyed call
// instead of serving the pre-update value.
func (s *CognitoStore) updateUserPoolUnlocked(userPool *UserPool) error {
	if !s.Exists(userPool.ID) {
		return ErrUserPoolNotFound
	}
	userPool.LastModifiedDate = time.Now().UTC()
	if err := s.Put(userPool.ID, userPool); err != nil {
		return err
	}
	s.usernameCaseCache.Store(userPool.ID,
		userPool.UsernameConfiguration != nil && userPool.UsernameConfiguration.CaseSensitive)
	return nil
}

// UpdateUserPoolFunc reads the pool, applies mutate and persists the result
// as one serialised mutation under the pool's key lock: concurrent updates
// cannot lose each other's field changes, and a pool deleted before or
// during the mutation surfaces ErrUserPoolNotFound instead of being
// overwritten. A mutate error aborts without writing.
func (s *CognitoStore) UpdateUserPoolFunc(userPoolID string, mutate func(*UserPool) error) error {
	return s.poolKeyLocker.WithLock(userPoolID, func() error {
		pool, err := s.GetUserPool(userPoolID)
		if err != nil {
			return err
		}
		if err := mutate(pool); err != nil {
			return err
		}
		return s.updateUserPoolUnlocked(pool)
	})
}

// DeleteUserPool deletes a Cognito user pool by ID and cascades to every
// family the pool owns — users with their tokens, groups, clients, challenge
// sessions, resource servers, identity providers, the domain binding, tags,
// devices, auth events, WebAuthn credentials, import jobs, and the prefixed
// settings families in the pool bucket. The deletion runs under the pool
// lock so pool-record mutations (updates, MFA configuration, schema
// additions) cannot interleave and resurrect deleted state.
func (s *CognitoStore) DeleteUserPool(userPoolID string) error {
	return s.poolKeyLocker.WithLock(userPoolID, func() error {
		return s.deleteUserPoolUnlocked(userPoolID)
	})
}

// deleteUserPoolUnlocked cascades the pool deletion across every bucket; the
// caller must hold the pool lock. Any failure aborts before the pool record
// itself is removed, so a retry completes the cascade; an entity deleted by
// a racing request between its list and its delete is skipped, not an error.
// The family sweeps run a second time after the pool record is removed —
// child writers that resolved the pool before the deletion began and land
// between the first sweep and the tombstone are caught by the second pass,
// and the pool-checked creators (users, clients) cannot land after the
// tombstone at all because their existence check runs under this same pool
// lock.
func (s *CognitoStore) deleteUserPoolUnlocked(userPoolID string) error {
	pool, err := s.GetUserPool(userPoolID)
	if err != nil {
		return err
	}

	if err := s.sweepPoolFamilies(userPoolID); err != nil {
		return err
	}

	if pool.Arn != "" {
		if err := s.TagStore.Delete(pool.Arn); err != nil {
			return err
		}
	}

	// The pool record is gone; its cached username-case sensitivity must not
	// outlive it (a recreated pool of the same ID may configure differently).
	s.usernameCaseCache.Delete(userPoolID)
	if err := s.BaseStore.Delete(userPoolID); err != nil {
		return err
	}
	return s.sweepPoolFamilies(userPoolID)
}

// sweepPoolFamilies removes every pool-scoped child family once: users with
// their tokens, groups, clients, the alias and case-insensitive username
// indexes, challenge sessions, devices, auth events, WebAuthn credentials,
// import jobs, the prefixed settings families, resource servers, identity
// providers and the domain binding. Already-absent entities are skipped.
func (s *CognitoStore) sweepPoolFamilies(userPoolID string) error {
	users, err := s.ListUsers(userPoolID)
	if err != nil {
		return err
	}
	for _, u := range users {
		prefix := userPoolID + "#" + u.ID + "#"
		if err := s.refreshTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
			return deleteTokenEntry(s.refreshTokensStore, key, value)
		}); err != nil {
			return err
		}
		if err := s.idTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
			return deleteTokenEntry(s.idTokensStore, key, value)
		}); err != nil {
			return err
		}
		if err := s.accessTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
			return deleteTokenEntry(s.accessTokensStore, key, value)
		}); err != nil {
			return err
		}
		if err := s.DeleteUser(userPoolID, u.Username); err != nil && !errors.Is(err, ErrUserNotFound) {
			return err
		}
	}

	groups, err := s.ListGroups(userPoolID)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if err := s.DeleteGroup(userPoolID, g.Name); err != nil && !errors.Is(err, ErrGroupNotFound) {
			return err
		}
	}

	// Cascade: delete this pool's alias index entries. The per-user deletes
	// above already release their own claims; the sweep catches entries left
	// behind by pool configuration changes.
	if err := s.usersStore.ScanPrefix("aliasidx:"+userPoolID+"#", func(key string, _ []byte) error {
		return s.usersStore.Delete(key)
	}); err != nil {
		return err
	}

	// Cascade: the case-insensitive username index shares the users bucket;
	// the same defensive sweep applies — a per-user delete interrupted or a
	// record predating the index must not leave pool-scoped entries behind.
	if err := s.usersStore.ScanPrefix("userci:"+userPoolID+"#", func(key string, _ []byte) error {
		return s.usersStore.Delete(key)
	}); err != nil {
		return err
	}

	clients, err := s.ListUserPoolClients(userPoolID)
	if err != nil {
		return err
	}
	for _, c := range clients {
		if err := s.DeleteUserPoolClient(userPoolID, c.ClientID); err != nil && !errors.Is(err, ErrClientNotFound) {
			return err
		}
	}

	// Cascade: delete challenge sessions. The pool lives in the record, not
	// in the key (keys are bare session IDs or "webauthn-reg#…" composites),
	// so the sweep filters by the stored pool. An undecodable record is
	// logged and skipped — it can no longer be attributed to a pool.
	if err := s.challengeSessionsStore.ForEach(func(key string, value []byte) error {
		var session ChallengeSession
		if err := json.Unmarshal(value, &session); err != nil {
			logs.Warn("cognito user-pool delete: skipping an undecodable challenge session", logs.String("key", key), logs.Err(err))
			return nil
		}
		if session.UserPoolID == userPoolID {
			return s.challengeSessionsStore.Delete(key)
		}
		return nil
	}); err != nil {
		return err
	}

	// Cascade: devices, auth events, WebAuthn credentials, and import jobs
	// are keyed "<poolID>#…" in their own buckets — one prefix batch per
	// bucket removes the pool's whole family.
	for _, bucket := range []*common.BaseStore{s.devicesStore, s.authEventsStore, s.webauthnStore, s.userImportJobsStore} {
		if err := bucket.ScanPrefix(userPoolID+"#", func(key string, _ []byte) error {
			return bucket.Delete(key)
		}); err != nil {
			return err
		}
	}

	// Cascade: the pool bucket holds six more families under their own key
	// prefixes — log delivery, risk configuration, and UI customisation
	// (each keyed "<prefix><poolID>" with optional "#<client>" per-client
	// keys), plus managed-login branding, terms, and replication regions.
	for _, prefix := range []string{
		logDeliveryKey(userPoolID),
		riskConfigKey(userPoolID, ""),
		uiCustomizationKey(userPoolID, ""),
		managedLoginBrandingPrefix(userPoolID),
		termsPrefix(userPoolID),
		userPoolReplicaPrefix(userPoolID),
	} {
		if err := s.BaseStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			return s.BaseStore.Delete(key)
		}); err != nil {
			return err
		}
	}

	// Cascade: delete resource servers.
	rsPrefix := resourceServerPrefix(userPoolID)
	if err := s.BaseStore.ScanPrefix(rsPrefix, func(key string, _ []byte) error {
		return s.BaseStore.Delete(key)
	}); err != nil {
		return err
	}

	// Cascade: delete identity providers.
	idpPrefix := identityProviderPrefix(userPoolID)
	if err := s.BaseStore.ScanPrefix(idpPrefix, func(key string, _ []byte) error {
		return s.BaseStore.Delete(key)
	}); err != nil {
		return err
	}

	// Cascade: delete domains associated with this pool. An undecodable
	// record is logged and skipped, like the challenge-session sweep.
	if err := s.BaseStore.ScanPrefix("domain:", func(key string, value []byte) error {
		var entry UserPoolDomain
		if err := json.Unmarshal(value, &entry); err != nil {
			logs.Warn("cognito user-pool delete: skipping an undecodable domain binding", logs.String("key", key), logs.Err(err))
			return nil
		}
		if entry.UserPoolID == userPoolID {
			return s.BaseStore.Delete(key)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// ListUserPools lists all Cognito user pools.
func (s *CognitoStore) ListUserPools() ([]*UserPool, error) {
	var userPools []*UserPool
	// The user pool bucket also contains other entities (domains, identity
	// providers, etc.) keyed by prefixed names like "domain:..." or
	// "identityprovider:...". Only keys matching the user pool ID format
	// ({region}_{uuid}) are actual user pool entries.
	err := s.ForEach(func(key string, value []byte) error {
		if !strings.HasPrefix(key, s.region+"_") {
			return nil
		}
		var userPool UserPool
		if err := json.Unmarshal(value, &userPool); err != nil {
			return err
		}
		userPools = append(userPools, &userPool)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userPools, nil
}

// ListUserPoolsPaginated lists Cognito user pools with server-side pagination.
// The user pool bucket contains other entities keyed by prefixed names; only
// keys matching the "{region}_{uuid}" pattern are actual user pool entries.
func (s *CognitoStore) ListUserPoolsPaginated(opts common.ListOptions) (*common.ListResult[UserPool], error) {
	return common.List[UserPool](s.BaseStore, opts, func(pool *UserPool) bool {
		return strings.HasPrefix(pool.ID, s.region+"_")
	})
}

// SetUserPoolDomain stores a domain configuration for a user pool,
// enforcing the domain binding rules under domainMu: a domain string
// already bound to another pool cannot be claimed, and a pool that already
// owns a domain cannot bind a second one. Re-binding a pool's own domain
// (the update path) is allowed. The checks and the write run as one
// serialised sequence so concurrent requests cannot interleave into a
// duplicate binding.
func (s *CognitoStore) SetUserPoolDomain(domain string, entry *UserPoolDomain) error {
	s.domainMu.Lock()
	defer s.domainMu.Unlock()
	var existing UserPoolDomain
	if err := s.BaseStore.Get(domainKey(domain), &existing); err == nil {
		if existing.UserPoolID != entry.UserPoolID {
			return ErrUserPoolDomainInUse
		}
	} else if !common.IsNotFound(err) {
		return err
	}
	if current, err := s.GetUserPoolDomainByPool(entry.UserPoolID); err == nil {
		if current.Domain != domain {
			return ErrUserPoolAlreadyHasDomain
		}
	} else if !errors.Is(err, ErrUserPoolDomainNotFound) {
		return err
	}
	return s.BaseStore.Put(domainKey(domain), entry)
}

// GetUserPoolDomain retrieves the domain configuration bound to a domain
// string. A missing domain reports ErrUserPoolDomainNotFound — a distinct
// sentinel from a missing pool, so callers can tell the two entities apart.
func (s *CognitoStore) GetUserPoolDomain(domain string) (*UserPoolDomain, error) {
	var entry UserPoolDomain
	if err := s.BaseStore.Get(domainKey(domain), &entry); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrUserPoolDomainNotFound
		}
		return nil, err
	}
	return &entry, nil
}

// GetUserPoolDomainByPool retrieves the domain bound to a user pool, if any.
func (s *CognitoStore) GetUserPoolDomainByPool(userPoolID string) (*UserPoolDomain, error) {
	var found *UserPoolDomain
	if err := s.BaseStore.ScanPrefix("domain:", func(key string, value []byte) error {
		var entry UserPoolDomain
		if err := json.Unmarshal(value, &entry); err == nil && entry.UserPoolID == userPoolID {
			found = &entry
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrUserPoolDomainNotFound
	}
	return found, nil
}

// DeleteUserPoolDomain removes a user pool's domain binding. Both the
// domain and the owning pool ID are required: the entry is deleted only
// when it exists and belongs to the named pool, so one pool cannot delete
// another pool's domain. A missing or foreign domain reports
// ErrUserPoolDomainNotFound.
func (s *CognitoStore) DeleteUserPoolDomain(userPoolID, domain string) error {
	s.domainMu.Lock()
	defer s.domainMu.Unlock()
	var entry UserPoolDomain
	if err := s.BaseStore.Get(domainKey(domain), &entry); err != nil {
		if common.IsNotFound(err) {
			return ErrUserPoolDomainNotFound
		}
		return err
	}
	if entry.UserPoolID != userPoolID {
		return ErrUserPoolDomainNotFound
	}
	return s.BaseStore.Delete(domainKey(domain))
}
