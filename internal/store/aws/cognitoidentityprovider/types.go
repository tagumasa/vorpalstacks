// Package cognito provides Cognito data store functionality for vorpalstacks.
package cognitoidentityprovider

import (
	"time"

	"vorpalstacks/internal/store/aws/common"

	"github.com/google/uuid"
)

// UserPool represents a Cognito user pool.
type UserPool struct {
	ID                     string          `json:"id"`
	Name                   string          `json:"name"`
	Arn                    string          `json:"arn"`
	Status                 string          `json:"status"`
	CreationDate           time.Time       `json:"creationDate"`
	LastModifiedDate       time.Time       `json:"lastModifiedDate"`
	AliasAttributes        []string        `json:"aliasAttributes,omitempty"`
	UsernameAttributes     []string        `json:"usernameAttributes,omitempty"`
	AutoVerifiedAttributes []string        `json:"autoVerifiedAttributes,omitempty"`
	Schema                 string          `json:"schema,omitempty"`
	MfaConfiguration       string          `json:"mfaConfiguration,omitempty"`
	PasswordPolicy         *PasswordPolicy `json:"passwordPolicy,omitempty"`
	LambdaConfig           *LambdaConfig   `json:"lambdaConfig,omitempty"`
	Tags                   []common.Tag    `json:"tags,omitempty"`
	EstimatedNumberOfUsers int64           `json:"estimatedNumberOfUsers,omitempty"`
	JwtPrivateKey          string          `json:"jwtPrivateKey,omitempty"`
	JwtPublicKey           string          `json:"jwtPublicKey,omitempty"`
	JwtKeyID               string          `json:"jwtKeyId,omitempty"`
}

// PasswordPolicy represents the password policy for a Cognito user pool.
type PasswordPolicy struct {
	MinimumLength                 int  `json:"minimumLength,omitempty"`
	RequireUppercase              bool `json:"requireUppercase,omitempty"`
	RequireLowercase              bool `json:"requireLowercase,omitempty"`
	RequireNumbers                bool `json:"requireNumbers,omitempty"`
	RequireSymbols                bool `json:"requireSymbols,omitempty"`
	TemporaryPasswordValidityDays int  `json:"temporaryPasswordValidityDays,omitempty"`
	PasswordHistorySize           int  `json:"passwordHistorySize,omitempty"`
	MaxPasswordAge                int  `json:"maxPasswordAge,omitempty"`
	PreventUserExistenceErrors    bool `json:"preventUserExistenceErrors,omitempty"`
}

// LambdaConfig represents the Lambda trigger configuration for a Cognito user pool.
type LambdaConfig struct {
	PreSignUp                   string `json:"preSignUp,omitempty"`
	CustomMessage               string `json:"customMessage,omitempty"`
	PostConfirmation            string `json:"postConfirmation,omitempty"`
	PreAuthentication           string `json:"preAuthentication,omitempty"`
	PostAuthentication          string `json:"postAuthentication,omitempty"`
	DefineAuthChallenge         string `json:"defineAuthChallenge,omitempty"`
	CreateAuthChallenge         string `json:"createAuthChallenge,omitempty"`
	VerifyAuthChallengeResponse string `json:"verifyAuthChallengeResponse,omitempty"`
	PreTokenGeneration          string `json:"preTokenGeneration,omitempty"`
	UserMigration               string `json:"userMigration,omitempty"`
}

// User represents a Cognito user.
type User struct {
	ID                     string                    `json:"id"`
	UserPoolID             string                    `json:"userPoolId"`
	Username               string                    `json:"username"`
	Enabled                bool                      `json:"enabled"`
	UserStatus             string                    `json:"userStatus"`
	CreatedDate            time.Time                 `json:"createdDate"`
	LastModifiedDate       time.Time                 `json:"lastModifiedDate"`
	Attributes             map[string]string         `json:"attributes,omitempty"`
	PasswordHash           string                    `json:"passwordHash,omitempty"`
	Groups                 []string                  `json:"groups,omitempty"`
	MFAOptions             []*MFAOptionType          `json:"mfaOptions,omitempty"`
	ConfirmationCode       string                    `json:"confirmationCode,omitempty"`
	ConfirmationCodeExpiry time.Time                 `json:"confirmationCodeExpiry,omitempty"`
	SoftwareTokenMfa       *SoftwareTokenMfaSettings `json:"softwareTokenMfa,omitempty"`
}

// GetID returns the unique identifier of the user.
func (u *User) GetID() string {
	return u.ID
}

// GetUsername returns the username of the user.
func (u *User) GetUsername() string {
	return u.Username
}

// GetGroups returns the groups the user belongs to.
func (u *User) GetGroups() []string {
	return u.Groups
}

// GetEmail returns the email address of the user if available.
func (u *User) GetEmail() string {
	if u.Attributes != nil {
		if email, ok := u.Attributes["email"]; ok {
			return email
		}
	}
	return ""
}

// GetCustomClaims returns custom claims for JWT token generation.
func (u *User) GetCustomClaims() map[string]interface{} {
	claims := make(map[string]interface{})
	if u.Attributes != nil {
		for k, v := range u.Attributes {
			if k != "email" && k != "sub" {
				claims[k] = v
			}
		}
	}
	return claims
}

// MFAOptionType represents MFA options for a Cognito user.
type MFAOptionType struct {
	DeliveryMedium string `json:"deliveryMedium,omitempty"`
	AttributeName  string `json:"attributeName,omitempty"`
}

// SoftwareTokenMfaSettings represents software token MFA settings for a user.
type SoftwareTokenMfaSettings struct {
	Enabled      bool   `json:"enabled"`
	PreferredMfa bool   `json:"preferredMfa"`
	SecretKey    string `json:"secretKey,omitempty"`
	Verified     bool   `json:"verified,omitempty"`
}

// Group represents a Cognito user group.
type Group struct {
	ID               string    `json:"id"`
	UserPoolID       string    `json:"userPoolId"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	RoleArn          string    `json:"roleArn,omitempty"`
	Precedence       int       `json:"precedence,omitempty"`
	CreationDate     time.Time `json:"creationDate"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
	Members          []string  `json:"members,omitempty"`
}

// UserPoolClient represents a Cognito user pool client.
type UserPoolClient struct {
	ClientID                        string    `json:"clientId"`
	UserPoolID                      string    `json:"userPoolId"`
	ClientName                      string    `json:"clientName"`
	ClientSecret                    string    `json:"clientSecret,omitempty"`
	RefreshTokenValidity            int       `json:"refreshTokenValidity"`
	AccessTokenValidity             int       `json:"accessTokenValidity"`
	IDTokenValidity                 int       `json:"idTokenValidity"`
	ExplicitAuthFlows               []string  `json:"explicitAuthFlows,omitempty"`
	AllowedOAuthFlows               []string  `json:"allowedOAuthFlows,omitempty"`
	CallbackURLs                    []string  `json:"callbackURLs,omitempty"`
	LogoutURLs                      []string  `json:"logoutURLs,omitempty"`
	DefaultRedirectURI              string    `json:"defaultRedirectUri,omitempty"`
	SupportedIdentityProviders      []string  `json:"supportedIdentityProviders,omitempty"`
	AllowedOAuthScopes              []string  `json:"allowedOAuthScopes,omitempty"`
	AllowedOAuthFlowsUserPoolClient bool      `json:"allowedOAuthFlowsUserPoolClient"`
	PreventUserExistenceErrors      string    `json:"preventUserExistenceErrors,omitempty"`
	CreationDate                    time.Time `json:"creationDate"`
	LastModifiedDate                time.Time `json:"lastModifiedDate"`
}

// RefreshToken represents a Cognito refresh token.
type RefreshToken struct {
	Token       string    `json:"token"`
	Expires     time.Time `json:"expires"`
	Scope       string    `json:"scope"`
	ClientID    string    `json:"clientId"`
	UserPoolID  string    `json:"userPoolId"`
	UserID      string    `json:"userId"`
	CreatedDate time.Time `json:"createdDate"`
}

// IDToken represents a Cognito ID token.
type IDToken struct {
	Token       string    `json:"token"`
	Expires     time.Time `json:"expires"`
	Scope       string    `json:"scope"`
	ClientID    string    `json:"clientId"`
	UserPoolID  string    `json:"userPoolId"`
	UserID      string    `json:"userId"`
	CreatedDate time.Time `json:"createdDate"`
	Groups      []string  `json:"groups,omitempty"`
}

// AccessToken represents a Cognito access token.
type AccessToken struct {
	Token       string    `json:"token"`
	Expires     time.Time `json:"expires"`
	Scope       string    `json:"scope"`
	ClientID    string    `json:"clientId"`
	UserPoolID  string    `json:"userPoolId"`
	UserID      string    `json:"userId"`
	CreatedDate time.Time `json:"createdDate"`
}

// NewUserPool creates a new Cognito user pool with the specified name.
func NewUserPool(name string, region string) *UserPool {
	now := time.Now().UTC()
	return &UserPool{
		ID:                     generateUserPoolID(region),
		Name:                   name,
		Status:                 "ACTIVE",
		CreationDate:           now,
		LastModifiedDate:       now,
		AliasAttributes:        []string{},
		UsernameAttributes:     []string{},
		AutoVerifiedAttributes: []string{},
		Schema:                 "{}",
		MfaConfiguration:       "OFF",
		Tags:                   []common.Tag{},
	}
}

// NewUser creates a new Cognito user in the specified user pool.
func NewUser(userPoolID, username string) *User {
	now := time.Now().UTC()
	return &User{
		ID:               generateID(),
		UserPoolID:       userPoolID,
		Username:         username,
		Enabled:          true,
		UserStatus:       "UNCONFIRMED",
		CreatedDate:      now,
		LastModifiedDate: now,
		Attributes:       make(map[string]string),
		Groups:           []string{},
	}
}

// NewGroup creates a new Cognito user group in the specified user pool.
func NewGroup(userPoolID, name string) *Group {
	now := time.Now().UTC()
	return &Group{
		ID:               generateID(),
		UserPoolID:       userPoolID,
		Name:             name,
		CreationDate:     now,
		LastModifiedDate: now,
		Members:          []string{},
	}
}

// NewUserPoolClient creates a new Cognito user pool client for the specified user pool.
func NewUserPoolClient(userPoolID, clientName string) *UserPoolClient {
	now := time.Now().UTC()
	return &UserPoolClient{
		ClientID:             generateClientID(),
		UserPoolID:           userPoolID,
		ClientName:           clientName,
		ClientSecret:         generateClientSecret(),
		RefreshTokenValidity: 30,
		AccessTokenValidity:  60,
		IDTokenValidity:      60,
		ExplicitAuthFlows:    []string{"ALLOW_USER_SRP_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"},
		CreationDate:         now,
		LastModifiedDate:     now,
	}
}

func generateClientID() string {
	return uuid.New().String()
}

func generateClientSecret() string {
	return uuid.New().String() + uuid.New().String()
}

// NewRefreshToken creates a new Cognito refresh token.
func NewRefreshToken(userPoolID, userID, clientID, scope string, expires time.Time) *RefreshToken {
	now := time.Now().UTC()
	return &RefreshToken{
		Token:       generateToken(),
		Expires:     expires,
		Scope:       scope,
		ClientID:    clientID,
		UserPoolID:  userPoolID,
		UserID:      userID,
		CreatedDate: now,
	}
}

// NewIDToken creates a new Cognito ID token.
func NewIDToken(userPoolID, userID, clientID, scope string, expires time.Time, groups []string) *IDToken {
	now := time.Now().UTC()
	return &IDToken{
		Token:       generateToken(),
		Expires:     expires,
		Scope:       scope,
		ClientID:    clientID,
		UserPoolID:  userPoolID,
		UserID:      userID,
		CreatedDate: now,
		Groups:      groups,
	}
}

// NewAccessToken creates a new Cognito access token.
func NewAccessToken(userPoolID, userID, clientID, scope string, expires time.Time) *AccessToken {
	now := time.Now().UTC()
	return &AccessToken{
		Token:       generateToken(),
		Expires:     expires,
		Scope:       scope,
		ClientID:    clientID,
		UserPoolID:  userPoolID,
		UserID:      userID,
		CreatedDate: now,
	}
}

func generateUserPoolID(region string) string {
	id := uuid.New().String()
	return region + "_" + id[:8]
}

func generateID() string {
	return uuid.New().String()
}

// ChallengeSession represents a Cognito auth challenge session.
type ChallengeSession struct {
	SessionID     string    `json:"sessionId"`
	UserPoolID    string    `json:"userPoolId"`
	ClientID      string    `json:"clientId"`
	Username      string    `json:"username"`
	ChallengeName string    `json:"challengeName"`
	CreatedAt     time.Time `json:"createdAt"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

func generateToken() string {
	return uuid.New().String() + uuid.New().String()
}

// UserPoolDomain represents a custom domain assigned to a Cognito user pool.
type UserPoolDomain struct {
	Domain           string    `json:"domain"`
	UserPoolID       string    `json:"userPoolId"`
	CloudFrontDomain string    `json:"cloudFrontDomain"`
	CreatedDate      time.Time `json:"createdDate"`
}

// ResourceServer represents an OAuth 2.0 resource server within a Cognito user pool.
type ResourceServer struct {
	UserPoolID string                `json:"userPoolId"`
	Identifier string                `json:"identifier"`
	Name       string                `json:"name"`
	Scopes     []ResourceServerScope `json:"scopes,omitempty"`
}

// ResourceServerScope defines a scope within an OAuth 2.0 resource server.
type ResourceServerScope struct {
	ScopeName        string `json:"scopeName"`
	ScopeDescription string `json:"scopeDescription"`
}

// IdentityProvider represents an external identity provider linked to a Cognito user pool.
type IdentityProvider struct {
	UserPoolID      string            `json:"userPoolId"`
	ProviderName    string            `json:"providerName"`
	ProviderType    string            `json:"providerType"`
	ProviderDetails map[string]string `json:"providerDetails,omitempty"`
}
