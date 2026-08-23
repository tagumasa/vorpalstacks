package cognitoidentityprovider

import (
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
)

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
	DeviceOperations
	AuthEventOperations
	LogDeliveryOperations
	RiskConfigurationOperations
	UICustomizationOperations
	UserImportJobOperations
	WebAuthnOperations
	ManagedLoginBrandingOperations
	TermsOperations
	UserPoolReplicaOperations
	Raw() *CognitoStore
}

// DeviceOperations defines operations for managing tracked devices.
type DeviceOperations interface {
	CreateDevice(d *Device) error
	GetDevice(userPoolID, userID, devKey string) (*Device, error)
	UpdateDevice(d *Device) error
	DeleteDevice(userPoolID, userID, devKey string) error
	ListDevicesPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[Device], error)
}

// AuthEventOperations defines operations for managing authentication events.
type AuthEventOperations interface {
	CreateAuthEvent(e *AuthEvent) error
	GetAuthEvent(userPoolID, userID, eventID string) (*AuthEvent, error)
	UpdateAuthEvent(e *AuthEvent) error
	ListAuthEventsPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[AuthEvent], error)
}

// LogDeliveryOperations defines operations for managing log delivery configuration.
type LogDeliveryOperations interface {
	GetLogDeliveryConfiguration(userPoolID string) (*LogDeliveryConfiguration, error)
	SaveLogDeliveryConfiguration(cfg *LogDeliveryConfiguration) error
}

// RiskConfigurationOperations defines operations for managing risk configuration.
type RiskConfigurationOperations interface {
	GetRiskConfiguration(userPoolID, clientID string) (*RiskConfiguration, error)
	SaveRiskConfiguration(cfg *RiskConfiguration) error
}

// UICustomizationOperations defines operations for managing UI customisation.
type UICustomizationOperations interface {
	GetUICustomization(userPoolID, clientID string) (*UICustomization, error)
	SaveUICustomization(ui *UICustomization) error
}

// UserImportJobOperations defines operations for managing user import jobs.
type UserImportJobOperations interface {
	CreateUserImportJob(job *UserImportJob) error
	GetUserImportJob(userPoolID, jobID string) (*UserImportJob, error)
	UpdateUserImportJob(job *UserImportJob) error
	ListUserImportJobsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserImportJob], error)
	// ListUserImportJobsAll returns every import job in the regional store
	// across all user pools, walking all pages. The start guard that limits
	// the account to one active import job needs the cross-pool view.
	ListUserImportJobsAll() ([]*UserImportJob, error)
	// StartUserImportJobIfEligible atomically moves a Created job to
	// Pending when no other job in the account is active.
	StartUserImportJobIfEligible(userPoolID, jobID string) (*UserImportJob, error)
	// TransitionUserImportJobStatus atomically moves a job from an
	// expected status to a new one; concurrent Start, Stop, and worker
	// finalisation serialise on the import-job lock.
	TransitionUserImportJobStatus(userPoolID, jobID, from, to string, mutate func(*UserImportJob)) (*UserImportJob, error)
	// UpdateUserImportJobProgress applies a mutation to a running job's
	// counters under the import-job lock.
	UpdateUserImportJobProgress(userPoolID, jobID string, mutate func(*UserImportJob)) error
}

// WebAuthnOperations defines operations for managing WebAuthn credentials.
type WebAuthnOperations interface {
	CreateWebAuthnCredential(c *WebAuthnCredential) error
	GetWebAuthnCredential(userPoolID, userID, credID string) (*WebAuthnCredential, error)
	DeleteWebAuthnCredential(userPoolID, userID, credID string) error
	ListWebAuthnCredentialsPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[WebAuthnCredential], error)
}

// ManagedLoginBrandingOperations defines operations for managing managed login branding.
type ManagedLoginBrandingOperations interface {
	SaveManagedLoginBranding(b *ManagedLoginBranding) error
	GetManagedLoginBranding(userPoolID, brandingID string) (*ManagedLoginBranding, error)
	GetManagedLoginBrandingByClient(userPoolID, clientID string) (*ManagedLoginBranding, error)
	DeleteManagedLoginBranding(userPoolID, brandingID string) error
	ListManagedLoginBrandings(userPoolID string) ([]*ManagedLoginBranding, error)
}

// TermsOperations defines operations for managing terms documents.
type TermsOperations interface {
	SaveTerms(t *Terms) error
	GetTerms(userPoolID, termsID string) (*Terms, error)
	DeleteTerms(userPoolID, termsID string) error
	ListTermsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[Terms], error)
}

// UserPoolReplicaOperations defines operations for managing user pool replicas.
type UserPoolReplicaOperations interface {
	SaveUserPoolReplica(r *UserPoolReplica) error
	GetUserPoolReplica(userPoolID, regionName string) (*UserPoolReplica, error)
	DeleteUserPoolReplica(userPoolID, regionName string) error
	ListUserPoolReplicasPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserPoolReplica], error)
}

// UserPoolOperations defines operations for managing user pools.
type UserPoolOperations interface {
	CreateUserPool(userPool *UserPool) (*UserPool, error)
	GetUserPool(userPoolID string) (*UserPool, error)
	UpdateUserPool(userPool *UserPool) error
	DeleteUserPool(userPoolID string) error
	ListUserPools() ([]*UserPool, error)
	ListUserPoolsPaginated(opts common.ListOptions) (*common.ListResult[UserPool], error)
}

// UserOperations defines operations for managing users.
type UserOperations interface {
	CreateUser(user *User) error
	GetUser(userPoolID, username string) (*User, error)
	GetUserByID(userID string) (*User, error)
	GetUserByProvider(userPoolID, providerName, providerAttrValue string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(userPoolID, username string) error
	ListUsers(userPoolID string) ([]*User, error)
	ListUsersPaginated(userPoolID string, opts common.ListOptions, filter func(*User) bool) (*common.ListResult[User], error)
}

// GroupOperations defines operations for managing groups.
type GroupOperations interface {
	CreateGroup(group *Group) error
	GetGroup(userPoolID, groupName string) (*Group, error)
	UpdateGroup(group *Group) error
	DeleteGroup(userPoolID, groupName string) error
	ListGroups(userPoolID string) ([]*Group, error)
	ListGroupsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[Group], error)
	AddUserToGroup(userPoolID, groupName, username string) error
	RemoveUserFromGroup(userPoolID, groupName, username string) error
	ListGroupsForUser(userPoolID, username string) ([]*Group, error)
	ListUsersInGroup(userPoolID, groupName string) ([]*User, error)
	ListUsersInGroupPaginated(userPoolID, groupName string, opts common.ListOptions) (*common.ListResult[User], error)
	ListGroupsForUserPaginated(userPoolID, username string, opts common.ListOptions) (*common.ListResult[Group], error)
}

// ClientOperations defines operations for managing user pool clients.
type ClientOperations interface {
	CreateUserPoolClient(client *UserPoolClient) error
	GetUserPoolClient(userPoolID, clientID string) (*UserPoolClient, error)
	GetUserPoolClientByName(userPoolID, clientName string) (*UserPoolClient, error)
	UpdateUserPoolClient(client *UserPoolClient) error
	DeleteUserPoolClient(userPoolID, clientID string) error
	ListUserPoolClients(userPoolID string) ([]*UserPoolClient, error)
	ListUserPoolClientsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserPoolClient], error)
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
	List(resourceArn string) (map[string]string, error)
	ListAsSlice(resourceArn string) ([]types.Tag, error)
	Tag(resourceArn string, tags map[string]string) error
	Untag(resourceArn string, tagKeys []string) error
}

// Raw returns the underlying Cognito store.
func (s *CognitoStore) Raw() *CognitoStore {
	return s
}

// DomainOperations defines operations for managing user pool domains.
type DomainOperations interface {
	SetUserPoolDomain(domain string, entry *UserPoolDomain) error
	GetUserPoolDomain(domain string) (*UserPoolDomain, error)
	GetUserPoolDomainByPool(userPoolID string) (*UserPoolDomain, error)
	DeleteUserPoolDomain(domain string) error
}

// ResourceServerOperations defines operations for managing resource servers.
type ResourceServerOperations interface {
	CreateResourceServer(rs *ResourceServer) error
	GetResourceServer(userPoolID, identifier string) (*ResourceServer, error)
	UpdateResourceServer(rs *ResourceServer) error
	DeleteResourceServer(userPoolID, identifier string) error
	ListResourceServers(userPoolID string) ([]*ResourceServer, error)
	ListResourceServersPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[ResourceServer], error)
}

// IdentityProviderOperations defines operations for managing identity providers.
type IdentityProviderOperations interface {
	CreateIdentityProvider(ip *IdentityProvider) error
	GetIdentityProvider(userPoolID, providerName string) (*IdentityProvider, error)
	UpdateIdentityProvider(ip *IdentityProvider) error
	DeleteIdentityProvider(userPoolID, providerName string) error
	ListIdentityProviders(userPoolID string) ([]*IdentityProvider, error)
	ListIdentityProvidersPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[IdentityProvider], error)
}
