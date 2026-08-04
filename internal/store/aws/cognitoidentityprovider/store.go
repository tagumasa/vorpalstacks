// Package cognito provides storage layer for AWS Cognito service entities
// including user pools, users, groups, and tokens.
package cognitoidentityprovider

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"

	"github.com/google/uuid"
	"sync"
)

// CognitoStore provides Cognito storage operations.
type CognitoStore struct {
	*common.BaseStore
	usersStore             *common.BaseStore
	groupsStore            *common.BaseStore
	clientsStore           *common.BaseStore
	refreshTokensStore     *common.BaseStore
	idTokensStore          *common.BaseStore
	accessTokensStore      *common.BaseStore
	challengeSessionsStore *common.BaseStore
	devicesStore           *common.BaseStore
	authEventsStore        *common.BaseStore
	userImportJobsStore    *common.BaseStore
	webauthnStore          *common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	groupMu    sync.Mutex
	createMu   sync.Mutex
}

// NewCognitoStore creates a new Cognito identity provider store.
func NewCognitoStore(store storage.BasicStorage, accountID, region string) *CognitoStore {
	return &CognitoStore{
		BaseStore:              common.NewBaseStore(store.Bucket(userPoolBucketName(region)), "cognito-userpools"),
		usersStore:             common.NewBaseStore(store.Bucket(userBucketName(region)), "cognito-users"),
		groupsStore:            common.NewBaseStore(store.Bucket(groupBucketName(region)), "cognito-groups"),
		clientsStore:           common.NewBaseStore(store.Bucket(clientBucketName(region)), "cognito-clients"),
		refreshTokensStore:     common.NewBaseStore(store.Bucket(refreshTokenBucketName(region)), "cognito-refreshtokens"),
		idTokensStore:          common.NewBaseStore(store.Bucket(idTokenBucketName(region)), "cognito-idtokens"),
		accessTokensStore:      common.NewBaseStore(store.Bucket(accessTokenBucketName(region)), "cognito-accesstokens"),
		challengeSessionsStore: common.NewBaseStore(store.Bucket(challengeSessionBucketName(region)), "cognito-challengesessions"),
		devicesStore:           common.NewBaseStore(store.Bucket(deviceBucketName(region)), "cognito-devices"),
		authEventsStore:        common.NewBaseStore(store.Bucket(authEventBucketName(region)), "cognito-authevents"),
		userImportJobsStore:    common.NewBaseStore(store.Bucket(userImportJobBucketName(region)), "cognito-userimportjobs"),
		webauthnStore:          common.NewBaseStore(store.Bucket(webauthnCredentialBucketName(region)), "cognito-webauthn"),
		TagStore:               common.NewTagStoreWithRegion(store, "cognito", region),
		arnBuilder:             svcarn.NewARNBuilder(accountID, region),
		accountID:              accountID,
		region:                 region,
	}
}

func (s *CognitoStore) buildUserPoolArn(userPoolID string) string {
	return s.arnBuilder.Cognito().UserPool(userPoolID)
}

// CreateUserPool creates a new Cognito user pool.
func (s *CognitoStore) CreateUserPool(userPool *UserPool) (*UserPool, error) {
	// Generate RSA key outside createMu to avoid blocking CreateUserPool
	// and CreateGroup operations during the 100-300ms key generation.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()
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
	userPool.JwtKeyID = uuid.New().String()[:8]
	userPool.JwtPrivateKey = encodePrivateKeyToPEM(privateKey)
	userPool.JwtPublicKey = encodePublicKeyToPEM(&privateKey.PublicKey)
	if userPool.Status == "" {
		userPool.Status = "ACTIVE"
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
		return nil, ErrUserPoolNotFound
	}
	return &userPool, nil
}

// UpdateUserPool updates an existing Cognito user pool.
func (s *CognitoStore) UpdateUserPool(userPool *UserPool) error {
	if !s.Exists(userPool.ID) {
		return ErrUserPoolNotFound
	}
	userPool.LastModifiedDate = time.Now().UTC()
	return s.Put(userPool.ID, userPool)
}

// DeleteUserPool deletes a Cognito user pool by ID and cascades to users,
// groups, clients, tokens, and challenge sessions.
func (s *CognitoStore) DeleteUserPool(userPoolID string) error {
	if !s.Exists(userPoolID) {
		return ErrUserPoolNotFound
	}

	users, _ := s.ListUsers(userPoolID)
	for _, u := range users {
		prefix := userPoolID + "#" + u.ID + "#"
		_ = s.refreshTokensStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			return s.refreshTokensStore.Delete(key)
		})
		_ = s.idTokensStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			return s.idTokensStore.Delete(key)
		})
		_ = s.accessTokensStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			return s.accessTokensStore.Delete(key)
		})
		_ = s.DeleteUser(userPoolID, u.Username)
	}

	groups, _ := s.ListGroups(userPoolID)
	for _, g := range groups {
		_ = s.DeleteGroup(userPoolID, g.Name)
	}

	clients, _ := s.ListUserPoolClients(userPoolID)
	for _, c := range clients {
		_ = s.DeleteUserPoolClient(userPoolID, c.ClientID)
	}

	_ = s.challengeSessionsStore.ScanPrefix(userPoolID+"#", func(key string, _ []byte) error {
		return s.challengeSessionsStore.Delete(key)
	})

	// Cascade: delete resource servers.
	rsPrefix := resourceServerPrefix(userPoolID)
	_ = s.BaseStore.ScanPrefix(rsPrefix, func(key string, _ []byte) error {
		return s.BaseStore.Delete(key)
	})

	// Cascade: delete identity providers.
	idpPrefix := identityProviderPrefix(userPoolID)
	_ = s.BaseStore.ScanPrefix(idpPrefix, func(key string, _ []byte) error {
		return s.BaseStore.Delete(key)
	})

	// Cascade: delete domains associated with this pool.
	_ = s.BaseStore.ScanPrefix("domain:", func(key string, value []byte) error {
		var entry UserPoolDomain
		if err := json.Unmarshal(value, &entry); err == nil && entry.UserPoolID == userPoolID {
			return s.BaseStore.Delete(key)
		}
		return nil
	})

	if pool, err := s.GetUserPool(userPoolID); err == nil && pool.Arn != "" {
		_ = s.TagStore.Delete(pool.Arn)
	}

	return s.BaseStore.Delete(userPoolID)
}

// ListUsersPaginated lists users in a Cognito user pool with server-side pagination.
// The filter callback is applied during iteration; filtered-out items do not count
// toward MaxItems, allowing efficient pagination even with selective filters.
func (s *CognitoStore) ListUsersPaginated(userPoolID string, opts common.ListOptions, filter func(*User) bool) (*common.ListResult[User], error) {
	opts.Prefix = userPoolID + "#"
	return common.List[User](s.usersStore, opts, filter)
}

// ListGroupsPaginated lists groups in a Cognito user pool with server-side pagination.
func (s *CognitoStore) ListGroupsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[Group], error) {
	opts.Prefix = userPoolID + "#"
	return common.List[Group](s.groupsStore, opts, nil)
}

// ListUserPoolClientsPaginated lists clients for a Cognito user pool with server-side pagination.
func (s *CognitoStore) ListUserPoolClientsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserPoolClient], error) {
	opts.Prefix = userPoolID + "#"
	return common.List[UserPoolClient](s.clientsStore, opts, nil)
}

// ListResourceServersPaginated lists resource servers for a user pool with server-side pagination.
func (s *CognitoStore) ListResourceServersPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[ResourceServer], error) {
	opts.Prefix = resourceServerPrefix(userPoolID)
	return common.List[ResourceServer](s.BaseStore, opts, nil)
}

// ListIdentityProvidersPaginated lists identity providers for a user pool with server-side pagination.
func (s *CognitoStore) ListIdentityProvidersPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[IdentityProvider], error) {
	opts.Prefix = identityProviderPrefix(userPoolID)
	return common.List[IdentityProvider](s.BaseStore, opts, nil)
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

// CreateGroup creates a new Cognito group.
func (s *CognitoStore) CreateGroup(group *Group) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if group.Name == "" {
		return ErrInvalidGroupName
	}

	key := userPoolGroupKey(group.UserPoolID, group.Name)
	if s.groupsStore.Exists(key) {
		return ErrGroupAlreadyExists
	}

	now := time.Now().UTC()
	group.CreationDate = now
	group.LastModifiedDate = now

	return s.groupsStore.Put(key, group)
}

// GetGroup retrieves a Cognito group by user pool ID and group name.
func (s *CognitoStore) GetGroup(userPoolID, groupName string) (*Group, error) {
	key := userPoolGroupKey(userPoolID, groupName)
	var group Group
	if err := s.groupsStore.Get(key, &group); err != nil {
		return nil, ErrGroupNotFound
	}
	return &group, nil
}

// UpdateGroup updates an existing Cognito group.
func (s *CognitoStore) UpdateGroup(group *Group) error {
	key := userPoolGroupKey(group.UserPoolID, group.Name)
	if !s.groupsStore.Exists(key) {
		return ErrGroupNotFound
	}
	group.LastModifiedDate = time.Now().UTC()
	return s.groupsStore.Put(key, group)
}

// DeleteGroup deletes a Cognito group.
func (s *CognitoStore) DeleteGroup(userPoolID, groupName string) error {
	s.groupMu.Lock()
	defer s.groupMu.Unlock()
	key := userPoolGroupKey(userPoolID, groupName)
	if !s.groupsStore.Exists(key) {
		return ErrGroupNotFound
	}
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return err
	}
	for _, member := range group.Members {
		user, err := s.GetUser(userPoolID, member)
		if err != nil {
			continue
		}
		var newGroups []string
		for _, g := range user.Groups {
			if g != groupName {
				newGroups = append(newGroups, g)
			}
		}
		user.Groups = newGroups
		if err := s.usersStore.Put(userPoolUserKey(userPoolID, member), user); err != nil {
			logs.Warn("failed to update user after group deletion", logs.String("user", member), logs.Err(err))
		}
	}
	return s.groupsStore.Delete(key)
}

// ListGroups lists all groups in a Cognito user pool.
func (s *CognitoStore) ListGroups(userPoolID string) ([]*Group, error) {
	var groups []*Group
	prefix := userPoolID + "#"
	err := s.groupsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var group Group
		if err := json.Unmarshal(value, &group); err != nil {
			return err
		}
		groups = append(groups, &group)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// AddUserToGroup adds a user to a Cognito group.
func (s *CognitoStore) AddUserToGroup(userPoolID, groupName, username string) error {
	s.groupMu.Lock()
	defer s.groupMu.Unlock()
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return err
	}

	for _, member := range group.Members {
		if member == username {
			return ErrUserAlreadyInGroup
		}
	}

	group.Members = append(group.Members, username)
	if err := s.UpdateGroup(group); err != nil {
		return err
	}

	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return err
	}

	for _, g := range user.Groups {
		if g == groupName {
			return nil
		}
	}

	user.Groups = append(user.Groups, groupName)
	return s.UpdateUser(user)
}

// RemoveUserFromGroup removes a user from a Cognito group.
func (s *CognitoStore) RemoveUserFromGroup(userPoolID, groupName, username string) error {
	s.groupMu.Lock()
	defer s.groupMu.Unlock()
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return err
	}

	found := false
	var newMembers []string
	for _, member := range group.Members {
		if member == username {
			found = true
		} else {
			newMembers = append(newMembers, member)
		}
	}

	if !found {
		return ErrUserNotInGroup
	}

	group.Members = newMembers
	if err := s.UpdateGroup(group); err != nil {
		return err
	}

	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return err
	}

	var newGroups []string
	for _, g := range user.Groups {
		if g != groupName {
			newGroups = append(newGroups, g)
		}
	}

	user.Groups = newGroups
	return s.UpdateUser(user)
}

// ListGroupsForUser lists all groups for a Cognito user.
func (s *CognitoStore) ListGroupsForUser(userPoolID, username string) ([]*Group, error) {
	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return nil, err
	}

	var groups []*Group
	for _, groupName := range user.Groups {
		group, err := s.GetGroup(userPoolID, groupName)
		if err == nil {
			groups = append(groups, group)
		}
	}
	return groups, nil
}

// CreateUserPoolClient creates a new Cognito user pool client.
func (s *CognitoStore) CreateUserPoolClient(client *UserPoolClient) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if client.ClientName == "" {
		return ErrInvalidParameter
	}

	key := userPoolClientKey(client.UserPoolID, client.ClientID)
	if s.clientsStore.Exists(key) {
		return ErrClientAlreadyExists
	}

	now := time.Now().UTC()
	client.CreationDate = now
	client.LastModifiedDate = now

	if err := s.clientsStore.Put(key, client); err != nil {
		return err
	}
	// Write secondary index for O(1) GetUserPoolByClientID lookup.
	_ = s.clientsStore.Put(clientIndexKey(client.ClientID), client.UserPoolID)
	return nil
}

// GetUserPoolClient retrieves a Cognito user pool client by client ID.
func (s *CognitoStore) GetUserPoolClient(userPoolID, clientID string) (*UserPoolClient, error) {
	key := userPoolClientKey(userPoolID, clientID)
	var client UserPoolClient
	if err := s.clientsStore.Get(key, &client); err != nil {
		return nil, ErrClientNotFound
	}
	return &client, nil
}

// GetUserPoolClientByName retrieves a Cognito user pool client by client name.
func (s *CognitoStore) GetUserPoolClientByName(userPoolID, clientName string) (*UserPoolClient, error) {
	var found *UserPoolClient
	prefix := userPoolID + "#"
	err := s.clientsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var client UserPoolClient
		if err := json.Unmarshal(value, &client); err != nil {
			return err
		}
		if client.ClientName == clientName {
			found = &client
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrClientNotFound
	}
	return found, nil
}

// UpdateUserPoolClient updates an existing Cognito user pool client.
func (s *CognitoStore) UpdateUserPoolClient(client *UserPoolClient) error {
	key := userPoolClientKey(client.UserPoolID, client.ClientID)
	if !s.clientsStore.Exists(key) {
		return ErrClientNotFound
	}
	client.LastModifiedDate = time.Now().UTC()
	return s.clientsStore.Put(key, client)
}

// DeleteUserPoolClient deletes a Cognito user pool client.
func (s *CognitoStore) DeleteUserPoolClient(userPoolID, clientID string) error {
	key := userPoolClientKey(userPoolID, clientID)
	if !s.clientsStore.Exists(key) {
		return ErrClientNotFound
	}
	// Clean up secondary index.
	_ = s.clientsStore.Delete(clientIndexKey(clientID))
	return s.clientsStore.Delete(key)
}

// ListUserPoolClients lists all clients for a Cognito user pool.
func (s *CognitoStore) ListUserPoolClients(userPoolID string) ([]*UserPoolClient, error) {
	var clients []*UserPoolClient
	prefix := userPoolID + "#"
	err := s.clientsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var client UserPoolClient
		if err := json.Unmarshal(value, &client); err != nil {
			return err
		}
		clients = append(clients, &client)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// GetUserPoolByClientID retrieves the user pool associated with a client ID.
func (s *CognitoStore) GetUserPoolByClientID(clientID string) (*UserPool, error) {
	// Use secondary index for O(1) lookup.
	var poolID string
	if err := s.clientsStore.Get(clientIndexKey(clientID), &poolID); err == nil && poolID != "" {
		return s.GetUserPool(poolID)
	}
	// Fallback: full scan for legacy clients without index
	pools, err := s.ListUserPools()
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		clients, err := s.ListUserPoolClients(pool.ID)
		if err != nil {
			continue
		}
		for _, client := range clients {
			if client.ClientID == clientID {
				// Backfill index
				_ = s.clientsStore.Put(clientIndexKey(clientID), pool.ID)
				return pool, nil
			}
		}
	}

	return nil, ErrClientNotFound
}

// ListUsersInGroup lists all users in a Cognito group.
func (s *CognitoStore) ListUsersInGroup(userPoolID, groupName string) ([]*User, error) {
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return nil, err
	}

	var users []*User
	for _, username := range group.Members {
		user, err := s.GetUser(userPoolID, username)
		if err == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

// ListUsersInGroupPaginated returns a page of users belonging to the
// specified group. The Marker is a zero-based index into group.Members.
func (s *CognitoStore) ListUsersInGroupPaginated(userPoolID, groupName string, opts common.ListOptions) (*common.ListResult[User], error) {
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return nil, err
	}

	members := group.Members
	start := 0
	if opts.Marker != "" {
		if idx, perr := strconv.Atoi(opts.Marker); perr == nil && idx >= 0 && idx < len(members) {
			start = idx
		}
	}

	limit := opts.MaxItems
	if limit <= 0 {
		limit = 60
	}

	end := start + limit
	if end > len(members) {
		end = len(members)
	}

	var users []*User
	for i := start; i < end; i++ {
		user, err := s.GetUser(userPoolID, members[i])
		if err == nil {
			users = append(users, user)
		}
	}

	result := &common.ListResult[User]{
		Items:       users,
		IsTruncated: end < len(members),
	}
	if end < len(members) {
		result.NextMarker = strconv.Itoa(end)
	}
	return result, nil
}

// ListGroupsForUserPaginated returns a page of groups the specified user
// belongs to. The Marker is a zero-based index into user.Groups.
func (s *CognitoStore) ListGroupsForUserPaginated(userPoolID, username string, opts common.ListOptions) (*common.ListResult[Group], error) {
	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return nil, err
	}

	groupNames := user.Groups
	start := 0
	if opts.Marker != "" {
		if idx, perr := strconv.Atoi(opts.Marker); perr == nil && idx >= 0 && idx < len(groupNames) {
			start = idx
		}
	}

	limit := opts.MaxItems
	if limit <= 0 {
		limit = 60
	}

	end := start + limit
	if end > len(groupNames) {
		end = len(groupNames)
	}

	var groups []*Group
	for i := start; i < end; i++ {
		group, err := s.GetGroup(userPoolID, groupNames[i])
		if err == nil {
			groups = append(groups, group)
		}
	}

	result := &common.ListResult[Group]{
		Items:       groups,
		IsTruncated: end < len(groupNames),
	}
	if end < len(groupNames) {
		result.NextMarker = strconv.Itoa(end)
	}
	return result, nil
}

// CreateRefreshToken creates a new Cognito refresh token.
func (s *CognitoStore) CreateRefreshToken(token *RefreshToken) error {
	key := tokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.refreshTokensStore.Put(key, token)
}

// GetRefreshToken retrieves a Cognito refresh token.
func (s *CognitoStore) GetRefreshToken(userPoolID, userID, token string) (*RefreshToken, error) {
	key := tokenKey(userPoolID, userID, token)
	var rt RefreshToken
	if err := s.refreshTokensStore.Get(key, &rt); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(rt.Expires) {
		_ = s.refreshTokensStore.Delete(key)
		return nil, ErrTokenExpired
	}
	return &rt, nil
}

// GetRefreshTokenByValue retrieves a Cognito refresh token by its token value.
func (s *CognitoStore) GetRefreshTokenByValue(token string) (*RefreshToken, error) {
	return findTokenByValue(s.refreshTokensStore, token, func(t *RefreshToken) string { return t.Token }, func(t *RefreshToken) time.Time { return t.Expires })
}

// DeleteRefreshToken deletes a Cognito refresh token.
func (s *CognitoStore) DeleteRefreshToken(userPoolID, userID, token string) error {
	key := tokenKey(userPoolID, userID, token)
	return s.refreshTokensStore.Delete(key)
}

// DeleteAllRefreshTokensForUser deletes all refresh tokens for a user.
func (s *CognitoStore) DeleteAllRefreshTokensForUser(userPoolID, userID string) error {
	prefix := userPoolID + "#" + userID + "#"
	return s.refreshTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.refreshTokensStore.Delete(key)
	})
}

// CreateIDToken creates a new Cognito ID token.
func (s *CognitoStore) CreateIDToken(token *IDToken) error {
	key := tokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.idTokensStore.Put(key, token)
}

// GetIDToken retrieves a Cognito ID token.
func (s *CognitoStore) GetIDToken(userPoolID, userID, token string) (*IDToken, error) {
	key := tokenKey(userPoolID, userID, token)
	var it IDToken
	if err := s.idTokensStore.Get(key, &it); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(it.Expires) {
		_ = s.idTokensStore.Delete(key)
		return nil, ErrTokenExpired
	}
	return &it, nil
}

// GetIDTokenByValue retrieves a Cognito ID token by its token value.
func (s *CognitoStore) GetIDTokenByValue(token string) (*IDToken, error) {
	return findTokenByValue(s.idTokensStore, token, func(t *IDToken) string { return t.Token }, func(t *IDToken) time.Time { return t.Expires })
}

// DeleteIDToken deletes a Cognito ID token.
func (s *CognitoStore) DeleteIDToken(userPoolID, userID, token string) error {
	key := tokenKey(userPoolID, userID, token)
	return s.idTokensStore.Delete(key)
}

// CreateAccessToken creates a new Cognito access token.
func (s *CognitoStore) CreateAccessToken(token *AccessToken) error {
	key := tokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.accessTokensStore.Put(key, token)
}

// GetAccessToken retrieves a Cognito access token.
func (s *CognitoStore) GetAccessToken(userPoolID, userID, token string) (*AccessToken, error) {
	key := tokenKey(userPoolID, userID, token)
	var at AccessToken
	if err := s.accessTokensStore.Get(key, &at); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(at.Expires) {
		_ = s.accessTokensStore.Delete(key)
		return nil, ErrTokenExpired
	}
	return &at, nil
}

// GetAccessTokenByValue retrieves a Cognito access token by its token value.
func (s *CognitoStore) GetAccessTokenByValue(token string) (*AccessToken, error) {
	return findTokenByValue(s.accessTokensStore, token, func(t *AccessToken) string { return t.Token }, func(t *AccessToken) time.Time { return t.Expires })
}

// DeleteAccessToken deletes a Cognito access token.
func (s *CognitoStore) DeleteAccessToken(userPoolID, userID, token string) error {
	key := tokenKey(userPoolID, userID, token)
	return s.accessTokensStore.Delete(key)
}

// DeleteUserTokens deletes all tokens for a user.
func (s *CognitoStore) DeleteUserTokens(userPoolID, userID string) error {
	prefix := userPoolID + "#" + userID + "#"
	var firstErr error

	if err := s.refreshTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.refreshTokensStore.Delete(key)
	}); err != nil && firstErr == nil {
		firstErr = err
	}

	if err := s.idTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.idTokensStore.Delete(key)
	}); err != nil && firstErr == nil {
		firstErr = err
	}

	if err := s.accessTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.accessTokensStore.Delete(key)
	}); err != nil && firstErr == nil {
		firstErr = err
	}

	return firstErr
}

// SaveChallengeSession saves a challenge session.
func (s *CognitoStore) SaveChallengeSession(session *ChallengeSession) error {
	return s.challengeSessionsStore.Put(session.SessionID, session)
}

// GetChallengeSession retrieves a challenge session by ID.
func (s *CognitoStore) GetChallengeSession(sessionID string) (*ChallengeSession, error) {
	var session ChallengeSession
	if err := s.challengeSessionsStore.Get(sessionID, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteChallengeSession deletes a challenge session.
func (s *CognitoStore) DeleteChallengeSession(sessionID string) error {
	return s.challengeSessionsStore.Delete(sessionID)
}

// SetUserPoolDomain stores a domain configuration for a user pool.
func (s *CognitoStore) SetUserPoolDomain(domain string, entry *UserPoolDomain) error {
	return s.BaseStore.Put(domainKey(domain), entry)
}

// GetUserPoolDomain retrieves the domain configuration for a user pool.
func (s *CognitoStore) GetUserPoolDomain(domain string) (*UserPoolDomain, error) {
	var entry UserPoolDomain
	if err := s.BaseStore.Get(domainKey(domain), &entry); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &entry, nil
}

// GetUserPoolDomainByPool retrieves the custom domain for a user pool, if any.
func (s *CognitoStore) GetUserPoolDomainByPool(userPoolID string) (*UserPoolDomain, error) {
	var found *UserPoolDomain
	_ = s.BaseStore.ScanPrefix("domain:", func(key string, value []byte) error {
		var entry UserPoolDomain
		if err := json.Unmarshal(value, &entry); err == nil && entry.UserPoolID == userPoolID {
			found = &entry
		}
		return nil
	})
	if found == nil {
		return nil, ErrUserPoolNotFound
	}
	return found, nil
}

// DeleteUserPoolDomain removes the domain configuration for a user pool.
func (s *CognitoStore) DeleteUserPoolDomain(domain string) error {
	return s.BaseStore.Delete(domainKey(domain))
}

// CreateResourceServer persists a new resource server for a user pool.
func (s *CognitoStore) CreateResourceServer(rs *ResourceServer) error {
	key := resourceServerKey(rs.UserPoolID, rs.Identifier)
	if s.BaseStore.Exists(key) {
		return ErrResourceAlreadyExists
	}
	now := time.Now().UTC()
	rs.CreationDate = now
	rs.LastModifiedDate = now
	return s.BaseStore.Put(key, rs)
}

// GetResourceServer retrieves a resource server by user pool ID and identifier.
func (s *CognitoStore) GetResourceServer(userPoolID, identifier string) (*ResourceServer, error) {
	key := resourceServerKey(userPoolID, identifier)
	var rs ResourceServer
	if err := s.BaseStore.Get(key, &rs); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &rs, nil
}

// UpdateResourceServer updates an existing resource server in the store.
func (s *CognitoStore) UpdateResourceServer(rs *ResourceServer) error {
	key := resourceServerKey(rs.UserPoolID, rs.Identifier)
	rs.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(key, rs)
}

// DeleteResourceServer removes a resource server from the store by user pool ID and identifier.
func (s *CognitoStore) DeleteResourceServer(userPoolID, identifier string) error {
	key := resourceServerKey(userPoolID, identifier)
	if !s.BaseStore.Exists(key) {
		return ErrUserPoolNotFound
	}
	return s.BaseStore.Delete(key)
}

// ListResourceServers lists all resource servers for a user pool.
func (s *CognitoStore) ListResourceServers(userPoolID string) ([]*ResourceServer, error) {
	var servers []*ResourceServer
	prefix := resourceServerPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var rs ResourceServer
		if err := json.Unmarshal(value, &rs); err != nil {
			return err
		}
		servers = append(servers, &rs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// CreateIdentityProvider persists a new identity provider for a user pool.
func (s *CognitoStore) CreateIdentityProvider(ip *IdentityProvider) error {
	key := identityProviderKey(ip.UserPoolID, ip.ProviderName)
	if s.BaseStore.Exists(key) {
		return ErrResourceAlreadyExists
	}
	now := time.Now().UTC()
	ip.CreationDate = now
	ip.LastModifiedDate = now
	return s.BaseStore.Put(key, ip)
}

// GetIdentityProvider retrieves an identity provider by user pool ID and provider name.
func (s *CognitoStore) GetIdentityProvider(userPoolID, providerName string) (*IdentityProvider, error) {
	key := identityProviderKey(userPoolID, providerName)
	var ip IdentityProvider
	if err := s.BaseStore.Get(key, &ip); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &ip, nil
}

// UpdateIdentityProvider updates an existing identity provider in the store.
func (s *CognitoStore) UpdateIdentityProvider(ip *IdentityProvider) error {
	key := identityProviderKey(ip.UserPoolID, ip.ProviderName)
	ip.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(key, ip)
}

// DeleteIdentityProvider removes an identity provider from the store by user pool ID and provider name.
func (s *CognitoStore) DeleteIdentityProvider(userPoolID, providerName string) error {
	key := identityProviderKey(userPoolID, providerName)
	if !s.BaseStore.Exists(key) {
		return ErrUserPoolNotFound
	}
	return s.BaseStore.Delete(key)
}

// ListIdentityProviders lists all identity providers for a user pool.
func (s *CognitoStore) ListIdentityProviders(userPoolID string) ([]*IdentityProvider, error) {
	var providers []*IdentityProvider
	prefix := identityProviderPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var ip IdentityProvider
		if err := json.Unmarshal(value, &ip); err != nil {
			return err
		}
		providers = append(providers, &ip)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return providers, nil
}

// CreateDevice stores a new device record.
func (s *CognitoStore) CreateDevice(d *Device) error {
	return s.devicesStore.Put(deviceKey(d.UserPoolID, d.UserID, d.DeviceKey), d)
}

// GetDevice retrieves a device by its key.
func (s *CognitoStore) GetDevice(userPoolID, userID, devKey string) (*Device, error) {
	var d Device
	if err := s.devicesStore.Get(deviceKey(userPoolID, userID, devKey), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// UpdateDevice updates an existing device record.
func (s *CognitoStore) UpdateDevice(d *Device) error {
	return s.devicesStore.Put(deviceKey(d.UserPoolID, d.UserID, d.DeviceKey), d)
}

// DeleteDevice removes a device record.
func (s *CognitoStore) DeleteDevice(userPoolID, userID, devKey string) error {
	return s.devicesStore.Delete(deviceKey(userPoolID, userID, devKey))
}

// ListDevicesPaginated lists devices for a user with server-side pagination.
func (s *CognitoStore) ListDevicesPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[Device], error) {
	opts.Prefix = devicePrefix(userPoolID, userID)
	return common.List[Device](s.devicesStore, opts, nil)
}

// CreateAuthEvent stores a new authentication event.
func (s *CognitoStore) CreateAuthEvent(e *AuthEvent) error {
	return s.authEventsStore.Put(authEventKey(e.UserPoolID, e.UserID, e.EventID), e)
}

// GetAuthEvent retrieves an auth event by its ID.
func (s *CognitoStore) GetAuthEvent(userPoolID, userID, eventID string) (*AuthEvent, error) {
	var e AuthEvent
	if err := s.authEventsStore.Get(authEventKey(userPoolID, userID, eventID), &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateAuthEvent updates an existing auth event record.
func (s *CognitoStore) UpdateAuthEvent(e *AuthEvent) error {
	return s.authEventsStore.Put(authEventKey(e.UserPoolID, e.UserID, e.EventID), e)
}

// ListAuthEventsPaginated lists auth events for a user with server-side pagination.
func (s *CognitoStore) ListAuthEventsPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[AuthEvent], error) {
	opts.Prefix = authEventPrefix(userPoolID, userID)
	return common.List[AuthEvent](s.authEventsStore, opts, nil)
}

// GetLogDeliveryConfiguration retrieves the log delivery configuration for a user pool.
func (s *CognitoStore) GetLogDeliveryConfiguration(userPoolID string) (*LogDeliveryConfiguration, error) {
	var cfg LogDeliveryConfiguration
	if err := s.BaseStore.Get(logDeliveryKey(userPoolID), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveLogDeliveryConfiguration stores the log delivery configuration for a user pool.
func (s *CognitoStore) SaveLogDeliveryConfiguration(cfg *LogDeliveryConfiguration) error {
	return s.BaseStore.Put(logDeliveryKey(cfg.UserPoolID), cfg)
}

// GetRiskConfiguration retrieves the risk configuration for a user pool/client.
func (s *CognitoStore) GetRiskConfiguration(userPoolID, clientID string) (*RiskConfiguration, error) {
	var cfg RiskConfiguration
	if err := s.BaseStore.Get(riskConfigKey(userPoolID, clientID), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveRiskConfiguration stores the risk configuration for a user pool/client.
func (s *CognitoStore) SaveRiskConfiguration(cfg *RiskConfiguration) error {
	cfg.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(riskConfigKey(cfg.UserPoolID, cfg.ClientID), cfg)
}

// GetUICustomization retrieves UI customisation for a user pool/client.
func (s *CognitoStore) GetUICustomization(userPoolID, clientID string) (*UICustomization, error) {
	var ui UICustomization
	if err := s.BaseStore.Get(uiCustomizationKey(userPoolID, clientID), &ui); err != nil {
		return nil, err
	}
	return &ui, nil
}

// SaveUICustomization stores UI customisation for a user pool/client.
// CSSVersion is incremented only when the CSS content actually changes,
// matching AWS behaviour for cache-busting on the hosted UI.
func (s *CognitoStore) SaveUICustomization(ui *UICustomization) error {
	key := uiCustomizationKey(ui.UserPoolID, ui.ClientID)
	now := time.Now().UTC()

	var prev UICustomization
	hasPrev := s.BaseStore.Get(key, &prev) == nil

	if !hasPrev {
		ui.CreationDate = now
		ui.CSSVersion = "20200101"
	} else {
		if ui.CreationDate.IsZero() {
			ui.CreationDate = prev.CreationDate
		}
		if ui.CSS != prev.CSS {
			if prev.CSSVersion == "" {
				ui.CSSVersion = "20200101"
			} else if v, err := strconv.Atoi(prev.CSSVersion); err == nil {
				ui.CSSVersion = strconv.Itoa(v + 1)
			} else {
				ui.CSSVersion = "20200101"
			}
		} else {
			ui.CSSVersion = prev.CSSVersion
		}
	}

	ui.LastModifiedDate = now
	return s.BaseStore.Put(key, ui)
}

// CreateUserImportJob stores a new user import job.
func (s *CognitoStore) CreateUserImportJob(job *UserImportJob) error {
	return s.userImportJobsStore.Put(userImportJobKey(job.UserPoolID, job.JobID), job)
}

// GetUserImportJob retrieves an import job by ID.
func (s *CognitoStore) GetUserImportJob(userPoolID, jobID string) (*UserImportJob, error) {
	var job UserImportJob
	if err := s.userImportJobsStore.Get(userImportJobKey(userPoolID, jobID), &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// UpdateUserImportJob updates an import job.
func (s *CognitoStore) UpdateUserImportJob(job *UserImportJob) error {
	return s.userImportJobsStore.Put(userImportJobKey(job.UserPoolID, job.JobID), job)
}

// ListUserImportJobs lists import jobs for a user pool.
func (s *CognitoStore) ListUserImportJobsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserImportJob], error) {
	opts.Prefix = userImportJobPrefix(userPoolID)
	return common.List[UserImportJob](s.userImportJobsStore, opts, nil)
}

// ===================== WebAuthn =====================

func (s *CognitoStore) CreateWebAuthnCredential(c *WebAuthnCredential) error {
	return s.webauthnStore.Put(webauthnKey(c.UserPoolID, c.UserID, c.CredentialID), c)
}

func (s *CognitoStore) GetWebAuthnCredential(userPoolID, userID, credID string) (*WebAuthnCredential, error) {
	var c WebAuthnCredential
	if err := s.webauthnStore.Get(webauthnKey(userPoolID, userID, credID), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CognitoStore) DeleteWebAuthnCredential(userPoolID, userID, credID string) error {
	return s.webauthnStore.Delete(webauthnKey(userPoolID, userID, credID))
}

func (s *CognitoStore) ListWebAuthnCredentialsPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[WebAuthnCredential], error) {
	opts.Prefix = webauthnPrefix(userPoolID, userID)
	return common.List[WebAuthnCredential](s.webauthnStore, opts, nil)
}

// ===================== Managed Login Branding =====================

func (s *CognitoStore) SaveManagedLoginBranding(b *ManagedLoginBranding) error {
	if b.CreationDate.IsZero() {
		b.CreationDate = time.Now().UTC()
	}
	b.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(managedLoginBrandingKey(b.UserPoolID, b.ManagedLoginBrandingId), b)
}

func (s *CognitoStore) GetManagedLoginBranding(userPoolID, brandingID string) (*ManagedLoginBranding, error) {
	var b ManagedLoginBranding
	if err := s.BaseStore.Get(managedLoginBrandingKey(userPoolID, brandingID), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *CognitoStore) GetManagedLoginBrandingByClient(userPoolID, clientID string) (*ManagedLoginBranding, error) {
	var found *ManagedLoginBranding
	prefix := managedLoginBrandingPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var b ManagedLoginBranding
		if err := json.Unmarshal(value, &b); err != nil {
			return err
		}
		if b.ClientID == clientID {
			found = &b
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrNotFound
	}
	return found, nil
}

func (s *CognitoStore) DeleteManagedLoginBranding(userPoolID, brandingID string) error {
	return s.BaseStore.Delete(managedLoginBrandingKey(userPoolID, brandingID))
}

func (s *CognitoStore) ListManagedLoginBrandings(userPoolID string) ([]*ManagedLoginBranding, error) {
	var brandings []*ManagedLoginBranding
	prefix := managedLoginBrandingPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var b ManagedLoginBranding
		if err := json.Unmarshal(value, &b); err != nil {
			return err
		}
		brandings = append(brandings, &b)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return brandings, nil
}

// ===================== Terms =====================

func (s *CognitoStore) SaveTerms(t *Terms) error {
	if t.CreationDate.IsZero() {
		t.CreationDate = time.Now().UTC()
	}
	t.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(termsKey(t.UserPoolID, t.TermsID), t)
}

func (s *CognitoStore) GetTerms(userPoolID, termsID string) (*Terms, error) {
	var t Terms
	if err := s.BaseStore.Get(termsKey(userPoolID, termsID), &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *CognitoStore) DeleteTerms(userPoolID, termsID string) error {
	return s.BaseStore.Delete(termsKey(userPoolID, termsID))
}

func (s *CognitoStore) ListTermsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[Terms], error) {
	opts.Prefix = termsPrefix(userPoolID)
	return common.List[Terms](s.BaseStore, opts, nil)
}

// ===================== User Pool Replicas =====================

func (s *CognitoStore) SaveUserPoolReplica(r *UserPoolReplica) error {
	return s.BaseStore.Put(userPoolReplicaKey(r.UserPoolID, r.RegionName), r)
}

func (s *CognitoStore) GetUserPoolReplica(userPoolID, regionName string) (*UserPoolReplica, error) {
	var r UserPoolReplica
	if err := s.BaseStore.Get(userPoolReplicaKey(userPoolID, regionName), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *CognitoStore) DeleteUserPoolReplica(userPoolID, regionName string) error {
	return s.BaseStore.Delete(userPoolReplicaKey(userPoolID, regionName))
}

func (s *CognitoStore) ListUserPoolReplicasPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserPoolReplica], error) {
	opts.Prefix = userPoolReplicaPrefix(userPoolID)
	return common.List[UserPoolReplica](s.BaseStore, opts, nil)
}
