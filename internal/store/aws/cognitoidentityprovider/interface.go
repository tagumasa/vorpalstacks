package cognitoidentityprovider

import "vorpalstacks/internal/store/aws/common"

// CognitoStoreInterface defines operations for managing Cognito User Pools.
type CognitoStoreInterface interface {
	UserPoolOperations
	UserOperations
	GroupOperations
	ClientOperations
	TokenOperations
	ChallengeOperations
	TagOperations
	DomainOperations
	ResourceServerOperations
	IdentityProviderOperations
	Raw() *CognitoStore
}

// UserPoolOperations defines operations for managing user pools.
type UserPoolOperations interface {
	CreateUserPool(userPool *UserPool) (*UserPool, error)
	GetUserPool(userPoolID string) (*UserPool, error)
	UpdateUserPool(userPool *UserPool) error
	DeleteUserPool(userPoolID string) error
	ListUserPools() ([]*UserPool, error)
}

// UserOperations defines operations for managing users.
type UserOperations interface {
	CreateUser(user *User) error
	GetUser(userPoolID, username string) (*User, error)
	GetUserByID(userID string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(userPoolID, username string) error
	ListUsers(userPoolID string) ([]*User, error)
}

// GroupOperations defines operations for managing groups.
type GroupOperations interface {
	CreateGroup(group *Group) error
	GetGroup(userPoolID, groupName string) (*Group, error)
	UpdateGroup(group *Group) error
	DeleteGroup(userPoolID, groupName string) error
	ListGroups(userPoolID string) ([]*Group, error)
	AddUserToGroup(userPoolID, groupName, username string) error
	RemoveUserFromGroup(userPoolID, groupName, username string) error
	ListGroupsForUser(userPoolID, username string) ([]*Group, error)
	ListUsersInGroup(userPoolID, groupName string) ([]*User, error)
}

// ClientOperations defines operations for managing user pool clients.
type ClientOperations interface {
	CreateUserPoolClient(client *UserPoolClient) error
	GetUserPoolClient(userPoolID, clientID string) (*UserPoolClient, error)
	GetUserPoolClientByName(userPoolID, clientName string) (*UserPoolClient, error)
	UpdateUserPoolClient(client *UserPoolClient) error
	DeleteUserPoolClient(userPoolID, clientID string) error
	ListUserPoolClients(userPoolID string) ([]*UserPoolClient, error)
	GetUserPoolByClientID(clientID string) (*UserPool, error)
}

// TokenOperations defines operations for managing tokens.
type TokenOperations interface {
	CreateRefreshToken(token *RefreshToken) error
	GetRefreshToken(userPoolID, userID, token string) (*RefreshToken, error)
	GetRefreshTokenByValue(token string) (*RefreshToken, error)
	DeleteRefreshToken(userPoolID, userID, token string) error
	DeleteAllRefreshTokensForUser(userPoolID, userID string) error

	CreateIDToken(token *IDToken) error
	GetIDToken(userPoolID, userID, token string) (*IDToken, error)
	GetIDTokenByValue(token string) (*IDToken, error)
	DeleteIDToken(userPoolID, userID, token string) error

	CreateAccessToken(token *AccessToken) error
	GetAccessToken(userPoolID, userID, token string) (*AccessToken, error)
	GetAccessTokenByValue(token string) (*AccessToken, error)
	DeleteAccessToken(userPoolID, userID, token string) error

	DeleteUserTokens(userPoolID, userID string) error
}

// ChallengeOperations defines operations for managing authentication challenges.
type ChallengeOperations interface {
	SaveChallengeSession(session *ChallengeSession) error
	GetChallengeSession(sessionID string) (*ChallengeSession, error)
	DeleteChallengeSession(sessionID string) error
}

// TagOperations defines operations for managing tags.
type TagOperations interface {
	ListTags(resourceArn string) (map[string]string, error)
	ListTagsAsSlice(resourceArn string) ([]common.Tag, error)
	TagResource(resourceArn string, tags map[string]string) error
	UntagResource(resourceArn string, tagKeys []string) error
}

// Raw returns the underlying Cognito store.
func (s *CognitoStore) Raw() *CognitoStore {
	return s
}

// DomainOperations defines operations for managing user pool domains.
type DomainOperations interface {
	SetUserPoolDomain(domain string, entry *UserPoolDomain) error
	GetUserPoolDomain(domain string) (*UserPoolDomain, error)
	DeleteUserPoolDomain(domain string) error
}

// ResourceServerOperations defines operations for managing resource servers.
type ResourceServerOperations interface {
	CreateResourceServer(rs *ResourceServer) error
	ListResourceServers(userPoolID string) ([]*ResourceServer, error)
}

// IdentityProviderOperations defines operations for managing identity providers.
type IdentityProviderOperations interface {
	CreateIdentityProvider(ip *IdentityProvider) error
	ListIdentityProviders(userPoolID string) ([]*IdentityProvider, error)
}
