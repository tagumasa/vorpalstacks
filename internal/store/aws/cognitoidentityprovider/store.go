// Package cognito provides storage layer for AWS Cognito service entities
// including user pools, users, groups, and tokens.
package cognitoidentityprovider

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"strings"
	"time"

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
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	groupMu    sync.Mutex
	createMu   sync.Mutex
}

func userPoolBucketName(region string) string {
	return "cognito-userpools-" + region
}

func userBucketName(region string) string {
	return "cognito-users-" + region
}

func groupBucketName(region string) string {
	return "cognito-groups-" + region
}

func clientBucketName(region string) string {
	return "cognito-clients-" + region
}

func refreshTokenBucketName(region string) string {
	return "cognito-refreshtokens-" + region
}

func idTokenBucketName(region string) string {
	return "cognito-idtokens-" + region
}

func accessTokenBucketName(region string) string {
	return "cognito-accesstokens-" + region
}

func challengeSessionBucketName(region string) string {
	return "cognito-challengesessions-" + region
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
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if userPool.Name == "" {
		return nil, ErrInvalidUserPoolName
	}

	if s.Exists(userPool.ID) {
		return nil, ErrUserPoolAlreadyExists
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
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

// DeleteUserPool deletes a Cognito user pool by ID.
func (s *CognitoStore) DeleteUserPool(userPoolID string) error {
	if !s.Exists(userPoolID) {
		return ErrUserPoolNotFound
	}
	return s.BaseStore.Delete(userPoolID)
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

	return s.usersStore.Put(key, user)
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
	var found *User
	err := s.usersStore.ForEach(func(key string, value []byte) error {
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
	return found, nil
}

// UpdateUser updates an existing Cognito user.
func (s *CognitoStore) UpdateUser(user *User) error {
	key := userPoolUserKey(user.UserPoolID, user.Username)
	if !s.usersStore.Exists(key) {
		return ErrUserNotFound
	}
	user.LastModifiedDate = time.Now().UTC()
	return s.usersStore.Put(key, user)
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
		s.groupsStore.Put(userPoolGroupKey(userPoolID, groupName), group)
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
		s.usersStore.Put(userPoolUserKey(userPoolID, member), user)
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

	return s.clientsStore.Put(key, client)
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

// CreateRefreshToken creates a new Cognito refresh token.
// CreateRefreshToken creates a new Cognito refresh token.
func (s *CognitoStore) CreateRefreshToken(token *RefreshToken) error {
	key := refreshTokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.refreshTokensStore.Put(key, token)
}

// GetRefreshToken retrieves a Cognito refresh token.
func (s *CognitoStore) GetRefreshToken(userPoolID, userID, token string) (*RefreshToken, error) {
	key := refreshTokenKey(userPoolID, userID, token)
	var rt RefreshToken
	if err := s.refreshTokensStore.Get(key, &rt); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(rt.Expires) {
		return nil, ErrTokenExpired
	}
	return &rt, nil
}

// GetRefreshTokenByValue retrieves a Cognito refresh token by its token value.
func (s *CognitoStore) GetRefreshTokenByValue(token string) (*RefreshToken, error) {
	var found *RefreshToken
	err := s.refreshTokensStore.ForEach(func(key string, value []byte) error {
		var rt RefreshToken
		if err := json.Unmarshal(value, &rt); err != nil {
			return err
		}
		if rt.Token == token {
			found = &rt
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(found.Expires) {
		return nil, ErrTokenExpired
	}
	return found, nil
}

// DeleteRefreshToken deletes a Cognito refresh token.
func (s *CognitoStore) DeleteRefreshToken(userPoolID, userID, token string) error {
	key := refreshTokenKey(userPoolID, userID, token)
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
// CreateIDToken creates a new Cognito ID token.
func (s *CognitoStore) CreateIDToken(token *IDToken) error {
	key := idTokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.idTokensStore.Put(key, token)
}

// GetIDToken retrieves a Cognito ID token.
// GetIDToken retrieves a Cognito ID token.
func (s *CognitoStore) GetIDToken(userPoolID, userID, token string) (*IDToken, error) {
	key := idTokenKey(userPoolID, userID, token)
	var it IDToken
	if err := s.idTokensStore.Get(key, &it); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(it.Expires) {
		return nil, ErrTokenExpired
	}
	return &it, nil
}

// GetIDTokenByValue retrieves a Cognito ID token by its token value.
// GetIDTokenByValue retrieves a Cognito ID token by its token value.
func (s *CognitoStore) GetIDTokenByValue(token string) (*IDToken, error) {
	var found *IDToken
	err := s.idTokensStore.ForEach(func(key string, value []byte) error {
		var it IDToken
		if err := json.Unmarshal(value, &it); err != nil {
			return err
		}
		if it.Token == token {
			found = &it
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(found.Expires) {
		return nil, ErrTokenExpired
	}
	return found, nil
}

// DeleteIDToken deletes a Cognito ID token.
func (s *CognitoStore) DeleteIDToken(userPoolID, userID, token string) error {
	key := idTokenKey(userPoolID, userID, token)
	return s.idTokensStore.Delete(key)
}

// CreateAccessToken creates a new Cognito access token.
// CreateAccessToken creates a new Cognito access token.
func (s *CognitoStore) CreateAccessToken(token *AccessToken) error {
	key := accessTokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.accessTokensStore.Put(key, token)
}

// GetAccessToken retrieves a Cognito access token.
// GetAccessToken retrieves a Cognito access token.
func (s *CognitoStore) GetAccessToken(userPoolID, userID, token string) (*AccessToken, error) {
	key := accessTokenKey(userPoolID, userID, token)
	var at AccessToken
	if err := s.accessTokensStore.Get(key, &at); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(at.Expires) {
		return nil, ErrTokenExpired
	}
	return &at, nil
}

// GetAccessTokenByValue retrieves a Cognito access token by its token value.
// GetAccessTokenByValue retrieves a Cognito access token by its token value.
func (s *CognitoStore) GetAccessTokenByValue(token string) (*AccessToken, error) {
	var found *AccessToken
	err := s.accessTokensStore.ForEach(func(key string, value []byte) error {
		var at AccessToken
		if err := json.Unmarshal(value, &at); err != nil {
			return err
		}
		if at.Token == token {
			found = &at
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(found.Expires) {
		return nil, ErrTokenExpired
	}
	return found, nil
}

// DeleteAccessToken deletes a Cognito access token.
func (s *CognitoStore) DeleteAccessToken(userPoolID, userID, token string) error {
	key := accessTokenKey(userPoolID, userID, token)
	return s.accessTokensStore.Delete(key)
}

func userPoolUserKey(userPoolID, username string) string {
	return userPoolID + "#" + username
}

func userPoolGroupKey(userPoolID, groupName string) string {
	return userPoolID + "#" + groupName
}

func userPoolClientKey(userPoolID, clientID string) string {
	return userPoolID + "#" + clientID
}

func refreshTokenKey(userPoolID, userID, token string) string {
	return userPoolID + "#" + userID + "#" + token
}

func idTokenKey(userPoolID, userID, token string) string {
	return userPoolID + "#" + userID + "#" + token
}

func accessTokenKey(userPoolID, userID, token string) string {
	return userPoolID + "#" + userID + "#" + token
}

func encodePrivateKeyToPEM(key *rsa.PrivateKey) string {
	der := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	}
	return string(pem.EncodeToMemory(block))
}

func encodePublicKeyToPEM(key *rsa.PublicKey) string {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return ""
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return string(pem.EncodeToMemory(block))
}

// DeleteUserTokens deletes all tokens for a user.
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
// SaveChallengeSession saves a challenge session.
func (s *CognitoStore) SaveChallengeSession(session *ChallengeSession) error {
	return s.challengeSessionsStore.Put(session.SessionID, session)
}

// GetChallengeSession retrieves a challenge session by ID.
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

func (s *CognitoStore) SetUserPoolDomain(domain string, entry *UserPoolDomain) error {
	return s.BaseStore.Put(domainKey(domain), entry)
}

func (s *CognitoStore) GetUserPoolDomain(domain string) (*UserPoolDomain, error) {
	var entry UserPoolDomain
	if err := s.BaseStore.Get(domainKey(domain), &entry); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &entry, nil
}

func (s *CognitoStore) DeleteUserPoolDomain(domain string) error {
	return s.BaseStore.Delete(domainKey(domain))
}

func (s *CognitoStore) CreateResourceServer(rs *ResourceServer) error {
	key := resourceServerKey(rs.UserPoolID, rs.Identifier)
	if s.BaseStore.Exists(key) {
		return ErrResourceAlreadyExists
	}
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

func (s *CognitoStore) CreateIdentityProvider(ip *IdentityProvider) error {
	key := identityProviderKey(ip.UserPoolID, ip.ProviderName)
	if s.BaseStore.Exists(key) {
		return ErrResourceAlreadyExists
	}
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

func resourceServerKey(userPoolID, identifier string) string {
	return "resourceserver:" + userPoolID + "#" + identifier
}

func resourceServerPrefix(userPoolID string) string {
	return "resourceserver:" + userPoolID + "#"
}

func identityProviderKey(userPoolID, providerName string) string {
	return "identityprovider:" + userPoolID + "#" + providerName
}

func identityProviderPrefix(userPoolID string) string {
	return "identityprovider:" + userPoolID + "#"
}

func domainKey(domain string) string {
	return "domain:" + domain
}
