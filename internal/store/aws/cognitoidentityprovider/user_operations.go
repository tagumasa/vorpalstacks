package cognitoidentityprovider

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// ListUsersPaginated lists users in a Cognito user pool with server-side pagination.
// The filter callback is applied during iteration; filtered-out items do not count
// toward MaxItems, allowing efficient pagination even with selective filters.
func (s *CognitoStore) ListUsersPaginated(userPoolID string, opts common.ListOptions, filter func(*User) bool) (*common.ListResult[User], error) {
	opts.Prefix = userPoolID + "#"
	return common.List[User](s.usersStore, opts, filter)
}

// userCIIndexKey is the case-insensitive username index: it maps the
// lowercased username to the canonical stored one so that lookups and
// uniqueness checks honour UsernameConfiguration. The index is written
// for every user but only consulted when the pool does not opt into
// case-sensitive usernames (the AWS default is case-insensitive).
func userCIIndexKey(userPoolID, username string) string {
	return "userci:" + userPoolID + "#" + strings.ToLower(username)
}

// usernameCaseSensitive reports whether the pool enforces case-sensitive
// usernames. A missing UsernameConfiguration means case-insensitive. The
// flag is served from the store-wide cache and only decoded from the pool
// record on first request — every username-keyed call consults it, and the
// pool-record write paths keep the cached value exact.
func (s *CognitoStore) usernameCaseSensitive(userPoolID string) bool {
	if cached, ok := s.usernameCaseCache.Load(userPoolID); ok {
		return cached.(bool)
	}
	pool, err := s.GetUserPool(userPoolID)
	if err != nil {
		return false
	}
	caseSensitive := pool.UsernameConfiguration != nil && pool.UsernameConfiguration.CaseSensitive
	s.usernameCaseCache.Store(userPoolID, caseSensitive)
	return caseSensitive
}

// resolveUsername maps the supplied username to the canonical stored
// username. Case-insensitive pools resolve through the lowercased index;
// when no index entry exists (an exact match, or a record predating the
// index) the input is returned unchanged.
func (s *CognitoStore) resolveUsername(userPoolID, username string) string {
	if s.usernameCaseSensitive(userPoolID) {
		return username
	}
	var canonical string
	if err := s.usersStore.Get(userCIIndexKey(userPoolID, username), &canonical); err == nil && canonical != "" {
		return canonical
	}
	return username
}

// CreateUser creates a new Cognito user.
func (s *CognitoStore) CreateUser(user *User) error {
	return s.createUser(user, false)
}

// CreateUserMigrateAliasClaims creates a new Cognito user, taking over alias
// values another user already claims instead of rejecting the conflict (the
// AdminCreateUser ForceAliasCreation semantics).
func (s *CognitoStore) CreateUserMigrateAliasClaims(user *User) error {
	return s.createUser(user, true)
}

func (s *CognitoStore) createUser(user *User, migrateAliases bool) error {
	// The pool-existence check and the write run under the pool lock so a
	// concurrent DeleteUserPool cascade cannot slip a user between its
	// family sweep and the pool-record tombstone. Lock order:
	// poolKeyLocker → recordMu, the same order the cascade takes.
	return s.poolKeyLocker.WithLock(user.UserPoolID, func() error {
		s.recordMu.Lock()
		defer s.recordMu.Unlock()
		if user.Username == "" {
			return ErrInvalidUsername
		}

		pool, err := s.GetUserPool(user.UserPoolID)
		if err != nil {
			return ErrUserPoolNotFound
		}

		key := userPoolUserKey(user.UserPoolID, user.Username)
		if s.usersStore.Exists(key) {
			return ErrUserAlreadyExists
		}
		// Uniqueness is case-insensitive unless the pool opts into
		// case-sensitive usernames.
		ciKey := userCIIndexKey(user.UserPoolID, user.Username)
		if !s.usernameCaseSensitive(user.UserPoolID) && s.usersStore.Exists(ciKey) {
			return ErrUserAlreadyExists
		}

		now := time.Now().UTC()
		user.CreatedDate = now
		user.LastModifiedDate = now

		// Claim alias values before the record lands so a rejected claim leaves
		// no user behind.
		if err := s.claimUserAliases(pool, user, migrateAliases); err != nil {
			return err
		}

		if err := s.usersStore.Put(key, user); err != nil {
			return err
		}
		// The record and its indexes commit together: a silently failed
		// index write would leave the user invisible to GetUserByID,
		// GetUserByProvider or the case-insensitive lookup, so it fails
		// the create.
		if err := s.usersStore.Put(userIndexKey(user.ID), user.UserPoolID+"#"+user.Username); err != nil {
			return err
		}
		if err := s.usersStore.Put(ciKey, user.Username); err != nil {
			return err
		}
		if user.ProviderName != "" && user.ProviderAttributeValue != "" {
			if err := s.usersStore.Put(providerIndexKey(user.UserPoolID, user.ProviderName, user.ProviderAttributeValue), user.Username); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetUser retrieves a Cognito user by user pool ID and username. The
// username matches case-insensitively unless the pool opts into
// case-sensitive usernames; when no user is registered under the literal
// username, the identifier is resolved through the alias index — the model's
// Username-parameter documentation accepts any alias attribute in place of
// the username, which is also how alias sign-in reaches the user record.
// Only a storage miss routes to the alias leg or reports ErrUserNotFound;
// an I/O failure is returned as-is and never masked as a missing user.
func (s *CognitoStore) GetUser(userPoolID, username string) (*User, error) {
	key := userPoolUserKey(userPoolID, s.resolveUsername(userPoolID, username))
	var user User
	if err := s.usersStore.Get(key, &user); err != nil {
		if !common.IsNotFound(err) {
			return nil, err
		}
		if alias := s.resolveAliasUsername(userPoolID, username); alias != "" {
			if err := s.usersStore.Get(userPoolUserKey(userPoolID, alias), &user); err != nil {
				if common.IsNotFound(err) {
					return nil, ErrUserNotFound
				}
				return nil, err
			}
			return &user, nil
		}
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// GetUserByID retrieves a Cognito user by user ID through the secondary
// index, resolving the owning pool and username without scanning the users
// bucket. The index is written together with the record, so an identifier
// without an index entry identifies no user.
func (s *CognitoStore) GetUserByID(userID string) (*User, error) {
	var idx string
	if err := s.usersStore.Get(userIndexKey(userID), &idx); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	parts := strings.SplitN(idx, "#", 2)
	if len(parts) != 2 {
		return nil, ErrUserNotFound
	}
	return s.GetUser(parts[0], parts[1])
}

// GetUserByProvider returns the user in a pool linked to the given federated
// provider with the matching attribute value, resolved through the provider
// index written at creation. Returns ErrUserNotFound if no user matches.
func (s *CognitoStore) GetUserByProvider(userPoolID, providerName, providerAttrValue string) (*User, error) {
	var username string
	if err := s.usersStore.Get(providerIndexKey(userPoolID, providerName, providerAttrValue), &username); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return s.GetUser(userPoolID, username)
}

// UpdateUser updates an existing Cognito user's attribute and identity
// state. Group membership is not part of the update: the membership
// operations own it and write it under the same lock from records they read
// inside their critical section, so an attribute update arriving with a
// stale membership snapshot refreshes it from the stored record instead of
// resurrecting or dropping groups.
func (s *CognitoStore) UpdateUser(user *User) error {
	s.recordMu.Lock()
	defer s.recordMu.Unlock()
	return s.updateUserUnlocked(user)
}

// UpdateUserFunc reads the user, applies mutate and persists the result as
// one serialised read-modify-write under recordMu: the mutation sees the
// freshest stored record, so concurrent writers cannot lose each other's
// changes — a whole-record UpdateUser of a stale snapshot is exactly that
// race (a concurrent password change's history silently overwritten by the
// older snapshot). Group membership is refreshed from the stored record as
// UpdateUser does, and a mutate error aborts without writing.
func (s *CognitoStore) UpdateUserFunc(userPoolID, username string, mutate func(*User) error) error {
	s.recordMu.Lock()
	defer s.recordMu.Unlock()
	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return err
	}
	if err := mutate(user); err != nil {
		return err
	}
	return s.updateUserUnlocked(user)
}

// updateUserUnlocked persists an updated user record together with its alias
// and identity indexes; the caller must hold recordMu.
func (s *CognitoStore) updateUserUnlocked(user *User) error {
	key := userPoolUserKey(user.UserPoolID, s.resolveUsername(user.UserPoolID, user.Username))
	var previous User
	if err := s.usersStore.Get(key, &previous); err != nil {
		return ErrUserNotFound
	}
	user.Groups = previous.Groups
	pool, err := s.GetUserPool(user.UserPoolID)
	if err != nil {
		return ErrUserPoolNotFound
	}
	// Newly claimed alias values are checked before anything is released or
	// written, so a rejected claim leaves both the record and the index
	// untouched; claims the new state no longer holds release their index
	// entries afterwards.
	if err := s.claimUserAliases(pool, user, false); err != nil {
		return err
	}
	s.releaseAliasClaims(user.UserPoolID, user.Username,
		userAliasClaims(pool, previous.Attributes), userAliasClaims(pool, user.Attributes))
	user.LastModifiedDate = time.Now().UTC()
	if err := s.usersStore.Put(key, user); err != nil {
		return err
	}
	// Keep the case-insensitive index pointing at the canonical username.
	if err := s.usersStore.Put(userCIIndexKey(user.UserPoolID, user.Username), user.Username); err != nil {
		return err
	}
	// The provider index must track the record's current federated identity:
	// when the triplet changes or is cleared, the entry the previous state
	// wrote is removed so the index can never resolve a user that no longer
	// carries that identity. Index failures propagate: a stale or missing
	// entry would silently break GetUserByProvider lookups.
	previousProvider := previous.ProviderName != "" && previous.ProviderAttributeValue != ""
	providerLinked := user.ProviderName != "" && user.ProviderAttributeValue != ""
	if previousProvider &&
		(!providerLinked || user.ProviderName != previous.ProviderName || user.ProviderAttributeValue != previous.ProviderAttributeValue) {
		if err := s.usersStore.Delete(providerIndexKey(user.UserPoolID, previous.ProviderName, previous.ProviderAttributeValue)); err != nil {
			return err
		}
	}
	if providerLinked {
		if err := s.usersStore.Put(providerIndexKey(user.UserPoolID, user.ProviderName, user.ProviderAttributeValue), user.Username); err != nil {
			return err
		}
	}
	return nil
}

// DeleteUser deletes a Cognito user.
func (s *CognitoStore) DeleteUser(userPoolID, username string) error {
	s.recordMu.Lock()
	defer s.recordMu.Unlock()
	canonical := s.resolveUsername(userPoolID, username)
	key := userPoolUserKey(userPoolID, canonical)
	if !s.usersStore.Exists(key) {
		return ErrUserNotFound
	}
	user, err := s.GetUser(userPoolID, canonical)
	if err != nil {
		// Exists proved the record a moment ago, so only a genuine
		// mid-window delete surfaces as not-found here — removing the
		// record is then the whole job. Any other failure is an I/O or
		// decode error: deleting the record anyway would corrupt every
		// index and group reference it owns, so the error propagates and
		// a retry completes the cleanup.
		if errors.Is(err, ErrUserNotFound) {
			return s.usersStore.Delete(key)
		}
		return err
	}
	for _, groupName := range user.Groups {
		group, err := s.GetGroup(userPoolID, groupName)
		if err != nil {
			continue
		}
		var newMembers []string
		for _, m := range group.Members {
			if m != canonical {
				newMembers = append(newMembers, m)
			}
		}
		group.Members = newMembers
		if err := s.groupsStore.Put(userPoolGroupKey(userPoolID, groupName), group); err != nil {
			logs.Warn("failed to update group after user deletion", logs.String("group", groupName), logs.Err(err))
		}
	}
	// Clean up secondary indexes.
	_ = s.usersStore.Delete(userIndexKey(user.ID))
	_ = s.usersStore.Delete(userCIIndexKey(userPoolID, canonical))
	if user.ProviderName != "" && user.ProviderAttributeValue != "" {
		_ = s.usersStore.Delete(providerIndexKey(userPoolID, user.ProviderName, user.ProviderAttributeValue))
	}
	// Release the alias values the user claims so the pool can reassign
	// them.
	if pool, err := s.GetUserPool(userPoolID); err == nil {
		s.releaseAliasClaims(userPoolID, canonical, userAliasClaims(pool, user.Attributes), nil)
	}
	return s.usersStore.Delete(key)
}

// ListUsers lists all users in a Cognito user pool.
func (s *CognitoStore) ListUsers(userPoolID string) ([]*User, error) {
	var users []*User
	prefix := userPoolID + "#"
	err := s.usersStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var user User
		if err := json.Unmarshal(value, &user); err != nil {
			return err
		}
		users = append(users, &user)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}
