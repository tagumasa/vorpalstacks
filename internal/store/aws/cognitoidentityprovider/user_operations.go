package cognitoidentityprovider

import (
	"encoding/json"
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

// CreateUser creates a new Cognito user.
func (s *CognitoStore) CreateUser(user *User) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if user.Username == "" {
		return ErrInvalidUsername
	}

	if !s.Exists(user.UserPoolID) {
		return ErrUserPoolNotFound
	}

	key := userPoolUserKey(user.UserPoolID, user.Username)
	if s.usersStore.Exists(key) {
		return ErrUserAlreadyExists
	}

	now := time.Now().UTC()
	user.CreatedDate = now
	user.LastModifiedDate = now

	if err := s.usersStore.Put(key, user); err != nil {
		return err
	}
	// Write secondary index for O(1) GetUserByID lookup.
	_ = s.usersStore.Put(userIndexKey(user.ID), user.UserPoolID+"#"+user.Username)
	// Write secondary index for O(1) GetUserByProvider lookup.
	if user.ProviderName != "" && user.ProviderAttributeValue != "" {
		_ = s.usersStore.Put(providerIndexKey(user.UserPoolID, user.ProviderName, user.ProviderAttributeValue), user.Username)
	}
	return nil
}

// GetUser retrieves a Cognito user by user pool ID and username.
func (s *CognitoStore) GetUser(userPoolID, username string) (*User, error) {
	key := userPoolUserKey(userPoolID, username)
	var user User
	if err := s.usersStore.Get(key, &user); err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// GetUserByID retrieves a Cognito user by user ID.
func (s *CognitoStore) GetUserByID(userID string) (*User, error) {
	// Use secondary index for O(1) lookup.
	var idx string
	if err := s.usersStore.Get(userIndexKey(userID), &idx); err == nil && idx != "" {
		parts := strings.SplitN(idx, "#", 2)
		if len(parts) == 2 {
			return s.GetUser(parts[0], parts[1])
		}
	}
	// Fallback: full scan for legacy users without index
	var found *User
	err := s.usersStore.ForEach(func(key string, value []byte) error {
		if strings.HasPrefix(key, "useridx:") || strings.HasPrefix(key, "providx:") {
			return nil
		}
		var user User
		if err := json.Unmarshal(value, &user); err != nil {
			return err
		}
		if user.ID == userID {
			found = &user
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrUserNotFound
	}
	// Backfill index for future lookups
	_ = s.usersStore.Put(userIndexKey(found.ID), found.UserPoolID+"#"+found.Username)
	return found, nil
}

// GetUserByProvider scans users in a pool and returns the one linked to the
// given federated provider with the matching attribute value. Returns
// ErrUserNotFound if no user matches.
func (s *CognitoStore) GetUserByProvider(userPoolID, providerName, providerAttrValue string) (*User, error) {
	// Use secondary index for O(1) lookup.
	var username string
	if err := s.usersStore.Get(providerIndexKey(userPoolID, providerName, providerAttrValue), &username); err == nil && username != "" {
		return s.GetUser(userPoolID, username)
	}
	// Fallback: prefix scan for legacy users without index
	var found *User
	prefix := userPoolID + "#"
	err := s.usersStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if strings.HasPrefix(key, "useridx:") || strings.HasPrefix(key, "providx:") {
			return nil
		}
		var user User
		if err := json.Unmarshal(value, &user); err != nil {
			return err
		}
		if user.ProviderName == providerName && user.ProviderAttributeValue == providerAttrValue {
			found = &user
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrUserNotFound
	}
	// Backfill index
	_ = s.usersStore.Put(providerIndexKey(userPoolID, providerName, providerAttrValue), found.Username)
	return found, nil
}

// UpdateUser updates an existing Cognito user.
func (s *CognitoStore) UpdateUser(user *User) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := userPoolUserKey(user.UserPoolID, user.Username)
	if !s.usersStore.Exists(key) {
		return ErrUserNotFound
	}
	user.LastModifiedDate = time.Now().UTC()
	if err := s.usersStore.Put(key, user); err != nil {
		return err
	}
	// Update provider index if provider info is set.
	if user.ProviderName != "" && user.ProviderAttributeValue != "" {
		_ = s.usersStore.Put(providerIndexKey(user.UserPoolID, user.ProviderName, user.ProviderAttributeValue), user.Username)
	}
	return nil
}

// DeleteUser deletes a Cognito user.
func (s *CognitoStore) DeleteUser(userPoolID, username string) error {
	s.groupMu.Lock()
	defer s.groupMu.Unlock()
	key := userPoolUserKey(userPoolID, username)
	if !s.usersStore.Exists(key) {
		return ErrUserNotFound
	}
	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return s.usersStore.Delete(key)
	}
	for _, groupName := range user.Groups {
		group, err := s.GetGroup(userPoolID, groupName)
		if err != nil {
			continue
		}
		var newMembers []string
		for _, m := range group.Members {
			if m != username {
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
	if user.ProviderName != "" && user.ProviderAttributeValue != "" {
		_ = s.usersStore.Delete(providerIndexKey(userPoolID, user.ProviderName, user.ProviderAttributeValue))
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
