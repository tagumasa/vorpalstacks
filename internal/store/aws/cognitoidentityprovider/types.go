// Package cognito provides Cognito data store functionality for vorpalstacks.
package cognitoidentityprovider

import (
	"time"

	types "vorpalstacks/internal/common/tags"

	"github.com/google/uuid"
)

// UserPool represents a Cognito user pool.
type UserPool struct {
	ID                            string                       `json:"id"`
	Name                          string                       `json:"name"`
	Arn                           string                       `json:"arn"`
	Status                        string                       `json:"status"`
	CreationDate                  time.Time                    `json:"creationDate"`
	LastModifiedDate              time.Time                    `json:"lastModifiedDate"`
	AliasAttributes               []string                     `json:"aliasAttributes,omitempty"`
	UsernameAttributes            []string                     `json:"usernameAttributes,omitempty"`
	AutoVerifiedAttributes        []string                     `json:"autoVerifiedAttributes,omitempty"`
	SchemaAttributes              []SchemaAttributeType        `json:"schemaAttributes,omitempty"`
	MfaConfiguration              string                       `json:"mfaConfiguration,omitempty"`
	PasswordPolicy                *PasswordPolicy              `json:"passwordPolicy,omitempty"`
	LambdaConfig                  *LambdaConfig                `json:"lambdaConfig,omitempty"`
	Tags                          []types.Tag                  `json:"tags,omitempty"`
	EstimatedNumberOfUsers        int64                        `json:"estimatedNumberOfUsers,omitempty"`
	JwtPrivateKey                 string                       `json:"jwtPrivateKey,omitempty"`
	JwtPublicKey                  string                       `json:"jwtPublicKey,omitempty"`
	JwtKeyID                      string                       `json:"jwtKeyId,omitempty"`
	EmailConfiguration            *EmailConfiguration          `json:"emailConfiguration,omitempty"`
	SmsConfiguration              *SmsConfiguration            `json:"smsConfiguration,omitempty"`
	AdminCreateUserConfig         *AdminCreateUserConfig       `json:"adminCreateUserConfig,omitempty"`
	VerificationMessageTemplate   *VerificationMessageTemplate `json:"verificationMessageTemplate,omitempty"`
	DeletionProtection            string                       `json:"deletionProtection,omitempty"`
	UserPoolAddOns                *UserPoolAddOns              `json:"userPoolAddOns,omitempty"`
	AccountRecoverySetting        *AccountRecoverySetting      `json:"accountRecoverySetting,omitempty"`
	UsernameConfiguration         *UsernameConfiguration       `json:"usernameConfiguration,omitempty"`
	DeviceConfiguration           *DeviceConfiguration         `json:"deviceConfiguration,omitempty"`
	EmailVerificationMessage      string                       `json:"emailVerificationMessage,omitempty"`
	EmailVerificationSubject      string                       `json:"emailVerificationSubject,omitempty"`
	SmsVerificationMessage        string                       `json:"smsVerificationMessage,omitempty"`
	SmsAuthenticationMessage      string                       `json:"smsAuthenticationMessage,omitempty"`
	UserAttributeUpdateSettings   *UserAttributeUpdateSettings `json:"userAttributeUpdateSettings,omitempty"`
	MfaConfigurationSms           *SmsMfaConfig                `json:"mfaConfigurationSms,omitempty"`
	MfaConfigurationSoftwareToken *MfaConfigurationType        `json:"mfaConfigurationSoftwareToken,omitempty"`
	EmailMfaConfig                *EmailMfaConfig              `json:"emailMfaConfig,omitempty"`
	WebAuthnConfiguration         *WebAuthnConfiguration       `json:"webAuthnConfiguration,omitempty"`
	IssuerConfiguration           *IssuerConfiguration         `json:"issuerConfiguration,omitempty"`
	KeyConfiguration              *KeyConfiguration            `json:"keyConfiguration,omitempty"`
	UserPoolTier                  string                       `json:"userPoolTier,omitempty"`
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
}

// EmailConfiguration represents the email configuration for a user pool.
type EmailConfiguration struct {
	SourceArn           string `json:"sourceArn,omitempty"`
	ReplyToEmailAddress string `json:"replyToEmailAddress,omitempty"`
	EmailSendingAccount string `json:"emailSendingAccount,omitempty"`
	From                string `json:"from,omitempty"`
	ConfigurationSet    string `json:"configurationSet,omitempty"`
}

// SmsConfiguration represents the SMS configuration for a user pool.
type SmsConfiguration struct {
	SnsCallerArn string `json:"snsCallerArn,omitempty"`
	ExternalId   string `json:"externalId,omitempty"`
	SnsRegion    string `json:"snsRegion,omitempty"`
}

// AdminCreateUserConfig represents the admin create user configuration for a user pool.
type AdminCreateUserConfig struct {
	AllowAdminCreateUserOnly  bool             `json:"allowAdminCreateUserOnly"`
	UnusedAccountValidityDays int              `json:"unusedAccountValidityDays,omitempty"`
	InviteMessageTemplate     *MessageTemplate `json:"inviteMessageTemplate,omitempty"`
}

// MessageTemplate represents the message template for invite messages.
type MessageTemplate struct {
	SMSMessage   string `json:"smsMessage,omitempty"`
	EmailMessage string `json:"emailMessage,omitempty"`
	EmailSubject string `json:"emailSubject,omitempty"`
}

// VerificationMessageTemplate represents the verification message template.
type VerificationMessageTemplate struct {
	SmsMessage         string `json:"smsMessage,omitempty"`
	EmailMessage       string `json:"emailMessage,omitempty"`
	EmailSubject       string `json:"emailSubject,omitempty"`
	EmailMessageByLink string `json:"emailMessageByLink,omitempty"`
	EmailSubjectByLink string `json:"emailSubjectByLink,omitempty"`
	DefaultEmailOption string `json:"defaultEmailOption,omitempty"`
}

// UserPoolAddOns represents advanced security configuration for a user pool.
type UserPoolAddOns struct {
	AdvancedSecurityMode string `json:"advancedSecurityMode,omitempty"`
}

// AccountRecoverySetting represents the account recovery setting for a user pool.
type AccountRecoverySetting struct {
	RecoveryMechanisms []RecoveryMechanism `json:"recoveryMechanisms,omitempty"`
}

// RecoveryMechanism represents a recovery mechanism.
type RecoveryMechanism struct {
	Priority int    `json:"priority"`
	Name     string `json:"name,omitempty"`
}

// UsernameConfiguration represents the username configuration.
type UsernameConfiguration struct {
	CaseSensitive bool `json:"caseSensitive"`
}

// DeviceConfiguration represents the device configuration.
type DeviceConfiguration struct {
	ChallengeRequiredOnNewDevice     bool `json:"challengeRequiredOnNewDevice"`
	DeviceOnlyRememberedOnUserPrompt bool `json:"deviceOnlyRememberedOnUserPrompt"`
}

// UserAttributeUpdateSettings represents the user attribute update settings.
type UserAttributeUpdateSettings struct {
	AttributesRequireVerificationBeforeUpdate []string `json:"attributesRequireVerificationBeforeUpdate,omitempty"`
}

// MfaConfigurationType represents MFA configuration details (SMS or SoftwareToken).
type MfaConfigurationType struct {
	Enabled bool `json:"enabled"`
}

// SmsMfaConfig holds SMS MFA settings per Smithy SmsMfaConfigType.
type SmsMfaConfig struct {
	SmsAuthenticationMessage string            `json:"smsAuthenticationMessage,omitempty"`
	SmsConfiguration         *SmsConfiguration `json:"smsConfiguration,omitempty"`
}

// EmailMfaConfig holds email-based MFA message templates per Smithy EmailMfaConfigType.
type EmailMfaConfig struct {
	Message string `json:"message,omitempty"`
	Subject string `json:"subject,omitempty"`
}

// WebAuthnConfiguration holds passkey/WebAuthn MFA settings.
type WebAuthnConfiguration struct {
	RelyingPartyId   string `json:"relyingPartyId,omitempty"`
	UserVerification string `json:"userVerification,omitempty"`
}

// LambdaConfig represents the Lambda trigger configuration for a Cognito user pool.
type LambdaConfig struct {
	PreSignUp                   string               `json:"preSignUp,omitempty"`
	CustomMessage               string               `json:"customMessage,omitempty"`
	PostConfirmation            string               `json:"postConfirmation,omitempty"`
	PreAuthentication           string               `json:"preAuthentication,omitempty"`
	PostAuthentication          string               `json:"postAuthentication,omitempty"`
	DefineAuthChallenge         string               `json:"defineAuthChallenge,omitempty"`
	CreateAuthChallenge         string               `json:"createAuthChallenge,omitempty"`
	VerifyAuthChallengeResponse string               `json:"verifyAuthChallengeResponse,omitempty"`
	PreTokenGeneration          string               `json:"preTokenGeneration,omitempty"`
	UserMigration               string               `json:"userMigration,omitempty"`
	KMSKeyID                    string               `json:"kmsKeyId,omitempty"`
	CustomEmailSender           *LambdaVersionConfig `json:"customEmailSender,omitempty"`
	CustomSMSSender             *LambdaVersionConfig `json:"customSmsSender,omitempty"`
	PreTokenGenerationConfig    *LambdaVersionConfig `json:"preTokenGenerationConfig,omitempty"`
	InboundFederation           *LambdaVersionConfig `json:"inboundFederation,omitempty"`
}

// LambdaVersionConfig pairs a Lambda ARN with its version for advanced
// trigger types (CustomEmailSender, CustomSMSSender, etc.).
type LambdaVersionConfig struct {
	LambdaArn     string `json:"lambdaArn,omitempty"`
	LambdaVersion string `json:"lambdaVersion,omitempty"`
}

// IssuerConfiguration controls token issuance settings.
type IssuerConfiguration struct {
	Type string `json:"type,omitempty"`
}

// KeyConfiguration holds encryption key settings for a user pool.
type KeyConfiguration struct {
	KeyType   string `json:"keyType,omitempty"`
	KmsKeyArn string `json:"kmsKeyArn,omitempty"`
}

// User represents a Cognito user.
type User struct {
	ID               string            `json:"id"`
	UserPoolID       string            `json:"userPoolId"`
	Username         string            `json:"username"`
	Enabled          bool              `json:"enabled"`
	UserStatus       string            `json:"userStatus"`
	CreatedDate      time.Time         `json:"createdDate"`
	LastModifiedDate time.Time         `json:"lastModifiedDate"`
	Attributes       map[string]string `json:"attributes,omitempty"`
	PasswordHash     string            `json:"passwordHash,omitempty"`
	// PasswordHashAlgo names the algorithm an imported PasswordHash was
	// encoded with (BCRYPT, SCRYPT, ARGON2ID, PBKDF2_SHA256). It is set by
	// CSV user import and cleared once the credentials migrate to the
	// native bcrypt+SRP pair at first successful sign-in.
	PasswordHashAlgo string `json:"passwordHashAlgo,omitempty"`
	// SrpSalt is the hex-encoded 16-byte random salt used to derive the SRP
	// verifier. It is sent to clients in the SALT ChallengeParameter.
	SrpSalt string `json:"srpSalt,omitempty"`
	// SrpVerifier is the hex-encoded SRP verifier v = g^x mod N. It is a
	// long-term secret stored at password-set time and never sent to clients.
	SrpVerifier                string                            `json:"srpVerifier,omitempty"`
	Groups                     []string                          `json:"groups,omitempty"`
	MFAOptions                 []*MFAOptionType                  `json:"mfaOptions,omitempty"`
	ConfirmationCode           string                            `json:"confirmationCode,omitempty"`
	ConfirmationCodeExpiry     time.Time                         `json:"confirmationCodeExpiry,omitempty"`
	SoftwareTokenMfa           *SoftwareTokenMfaSettings         `json:"softwareTokenMfa,omitempty"`
	SmsMfa                     *SmsMfaSettings                   `json:"smsMfa,omitempty"`
	EmailMfa                   *EmailMfaSettings                 `json:"emailMfa,omitempty"`
	WebAuthnMfaEnabled         bool                              `json:"webAuthnMfaEnabled,omitempty"`
	ProviderName               string                            `json:"providerName,omitempty"`
	ProviderAttributeName      string                            `json:"providerAttributeName,omitempty"`
	ProviderAttributeValue     string                            `json:"providerAttributeValue,omitempty"`
	AttributeVerificationCodes map[string]*AttributeVerification `json:"attributeVerificationCodes,omitempty"`
}

// AttributeVerification holds a per-attribute verification code, independent
// from the SignUp confirmation code. Keyed by attribute name (e.g. "email",
// "phone_number") so that codes for multiple attributes can coexist.
type AttributeVerification struct {
	Code   string    `json:"code"`
	Expiry time.Time `json:"expiry"`
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

// SmsMfaSettings represents SMS-based MFA settings for a user.
type SmsMfaSettings struct {
	Enabled      bool `json:"enabled"`
	PreferredMfa bool `json:"preferredMfa"`
}

// EmailMfaSettings represents email-based MFA settings for a user.
type EmailMfaSettings struct {
	Enabled      bool `json:"enabled"`
	PreferredMfa bool `json:"preferredMfa"`
}

// Group represents a Cognito user group.
type Group struct {
	ID               string    `json:"id"`
	UserPoolID       string    `json:"userPoolId"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	RoleArn          string    `json:"roleArn,omitempty"`
	Precedence       *int      `json:"precedence,omitempty"`
	CreationDate     time.Time `json:"creationDate"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
	Members          []string  `json:"members,omitempty"`
}

// UserPoolClient represents a Cognito user pool client.
type UserPoolClient struct {
	ClientID                                 string                   `json:"clientId"`
	UserPoolID                               string                   `json:"userPoolId"`
	ClientName                               string                   `json:"clientName"`
	ClientSecret                             string                   `json:"clientSecret,omitempty"`
	RefreshTokenValidity                     int                      `json:"refreshTokenValidity"`
	AccessTokenValidity                      int                      `json:"accessTokenValidity"`
	IDTokenValidity                          int                      `json:"idTokenValidity"`
	ExplicitAuthFlows                        []string                 `json:"explicitAuthFlows,omitempty"`
	AllowedOAuthFlows                        []string                 `json:"allowedOAuthFlows,omitempty"`
	CallbackURLs                             []string                 `json:"callbackURLs,omitempty"`
	LogoutURLs                               []string                 `json:"logoutURLs,omitempty"`
	DefaultRedirectURI                       string                   `json:"defaultRedirectUri,omitempty"`
	SupportedIdentityProviders               []string                 `json:"supportedIdentityProviders,omitempty"`
	AllowedOAuthScopes                       []string                 `json:"allowedOAuthScopes,omitempty"`
	AllowedOAuthFlowsUserPoolClient          bool                     `json:"allowedOAuthFlowsUserPoolClient"`
	PreventUserExistenceErrors               string                   `json:"preventUserExistenceErrors,omitempty"`
	ClientSecrets                            []ClientSecretDescriptor `json:"clientSecrets,omitempty"`
	CreationDate                             time.Time                `json:"creationDate"`
	LastModifiedDate                         time.Time                `json:"lastModifiedDate"`
	AnalyticsConfiguration                   *AnalyticsConfiguration  `json:"analyticsConfiguration,omitempty"`
	AuthSessionValidity                      int                      `json:"authSessionValidity,omitempty"`
	EnablePropagateAdditionalUserContextData bool                     `json:"enablePropagateAdditionalUserContextData,omitempty"`
	EnableTokenRevocation                    bool                     `json:"enableTokenRevocation,omitempty"`
	GenerateSecret                           bool                     `json:"generateSecret,omitempty"`
	ReadAttributes                           []string                 `json:"readAttributes,omitempty"`
	RefreshTokenRotation                     *RefreshTokenRotation    `json:"refreshTokenRotation,omitempty"`
	TokenValidityUnits                       *TokenValidityUnits      `json:"tokenValidityUnits,omitempty"`
	WriteAttributes                          []string                 `json:"writeAttributes,omitempty"`
}

// AnalyticsConfiguration holds Amazon Pinpoint analytics configuration.
type AnalyticsConfiguration struct {
	ApplicationArn string `json:"applicationArn,omitempty"`
	ApplicationId  string `json:"applicationId,omitempty"`
	ExternalId     string `json:"externalId,omitempty"`
	RoleArn        string `json:"roleArn,omitempty"`
	UserDataShared bool   `json:"userDataShared,omitempty"`
}

// RefreshTokenRotation controls refresh token rotation behaviour.
type RefreshTokenRotation struct {
	Feature                 string `json:"feature,omitempty"`
	RetryGracePeriodSeconds int    `json:"retryGracePeriodSeconds,omitempty"`
}

// TokenValidityUnits specifies the time units for token validity periods.
type TokenValidityUnits struct {
	AccessToken  string `json:"accessToken,omitempty"`
	IdToken      string `json:"idToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
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
		MfaConfiguration:       "OFF",
		Tags:                   []types.Tag{},
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

// Default token validity periods for a user pool client, applied when the
// client (or its validity members) are absent: access and ID tokens last
// sixty minutes, refresh tokens thirty days.
const (
	DefaultAccessTokenValidityMinutes = 60
	DefaultIDTokenValidityMinutes     = 60
	DefaultRefreshTokenValidityDays   = 30
)

// NewUserPoolClient creates a new Cognito user pool client for the specified user pool.
func NewUserPoolClient(userPoolID, clientName string) *UserPoolClient {
	now := time.Now().UTC()
	return &UserPoolClient{
		ClientID:             generateClientID(),
		UserPoolID:           userPoolID,
		ClientName:           clientName,
		ClientSecret:         generateClientSecret(),
		RefreshTokenValidity: DefaultRefreshTokenValidityDays,
		AccessTokenValidity:  60,
		IDTokenValidity:      DefaultIDTokenValidityMinutes,
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
	// SRP state for PASSWORD_VERIFIER challenges. SrpA, SrpB and SrpPrivateB
	// are hex strings (matching the on-the-wire format of SRP_A/SRP_B); only
	// SecretBlock is base64 because it is opaque binary returned to the client
	// verbatim. SrpPrivateB is the server's secret scalar b and must never be
	// disclosed.
	SrpA        string `json:"srpA,omitempty"`
	SrpB        string `json:"srpB,omitempty"`
	SrpPrivateB string `json:"srpPrivateB,omitempty"`
	SecretBlock string `json:"secretBlock,omitempty"`
	// ChallengeData stores the opaque challenge bytes (base64 RawURL encoded)
	// for challenge types that require server-side verification (e.g.
	// WEB_AUTHN_REGISTRATION: the WebAuthn challenge that the client must sign
	// with their authenticator).
	ChallengeData string `json:"challengeData,omitempty"`
	// RelyingPartyID stores the user pool relying party id under which a
	// WEB_AUTHN_REGISTRATION challenge was issued, so the completion verifies
	// origin and rpIdHash against the exact value offered at Start.
	RelyingPartyID string `json:"relyingPartyId,omitempty"`
	// FailedAttempts counts wrong verification answers accepted for this
	// session. Sessions exceeding the attempt budget are invalidated so
	// short numeric codes cannot be brute-forced within one session.
	FailedAttempts int `json:"failedAttempts,omitempty"`
}

func generateToken() string {
	return uuid.New().String() + uuid.New().String()
}

// UserPoolDomain represents a custom domain assigned to a Cognito user pool.
type UserPoolDomain struct {
	Domain              string              `json:"domain"`
	UserPoolID          string              `json:"userPoolId"`
	CloudFrontDomain    string              `json:"cloudFrontDomain"`
	CreatedDate         time.Time           `json:"createdDate"`
	ManagedLoginVersion *int                `json:"managedLoginVersion,omitempty"`
	CustomDomainConfig  *CustomDomainConfig `json:"customDomainConfig,omitempty"`
	Routing             *Routing            `json:"routing,omitempty"`
	Status              string              `json:"status,omitempty"`
}

// CustomDomainConfig holds the certificate and security policy for a custom
// domain.
type CustomDomainConfig struct {
	CertificateArn string `json:"certificateArn,omitempty"`
	SecurityPolicy string `json:"securityPolicy,omitempty"`
}

// Routing holds the failover configuration for a user pool domain.
type Routing struct {
	Failover *FailoverType `json:"failover,omitempty"`
}

// FailoverType holds the failover detail for a domain routing configuration.
type FailoverType struct {
	SecondaryRegion             string `json:"secondaryRegion,omitempty"`
	PrimaryRoute53HealthCheckId string `json:"primaryRoute53HealthCheckId,omitempty"`
}

// ResourceServer represents an OAuth 2.0 resource server within a Cognito user pool.
type ResourceServer struct {
	UserPoolID       string                `json:"userPoolId"`
	Identifier       string                `json:"identifier"`
	Name             string                `json:"name"`
	Scopes           []ResourceServerScope `json:"scopes,omitempty"`
	CreationDate     time.Time             `json:"creationDate"`
	LastModifiedDate time.Time             `json:"lastModifiedDate"`
}

// ResourceServerScope defines a scope within an OAuth 2.0 resource server.
type ResourceServerScope struct {
	ScopeName        string `json:"scopeName"`
	ScopeDescription string `json:"scopeDescription"`
}

// IdentityProvider represents an external identity provider linked to a Cognito user pool.
type IdentityProvider struct {
	UserPoolID       string            `json:"userPoolId"`
	ProviderName     string            `json:"providerName"`
	ProviderType     string            `json:"providerType"`
	ProviderDetails  map[string]string `json:"providerDetails,omitempty"`
	AttributeMapping map[string]string `json:"attributeMapping,omitempty"`
	IdpIdentifiers   []string          `json:"idpIdentifiers,omitempty"`
	CreationDate     time.Time         `json:"creationDate"`
	LastModifiedDate time.Time         `json:"lastModifiedDate"`
}

// Device represents a tracked device in a Cognito user pool.
type Device struct {
	DeviceKey                   string            `json:"deviceKey"`
	UserPoolID                  string            `json:"userPoolId"`
	UserID                      string            `json:"userId"`
	DeviceName                  string            `json:"deviceName,omitempty"`
	DeviceAttributes            map[string]string `json:"deviceAttributes,omitempty"`
	DeviceSecretVerifierB       string            `json:"deviceSecretVerifierB,omitempty"`
	DeviceSaltVerifier          string            `json:"deviceSaltVerifier,omitempty"`
	DeviceCreateDate            time.Time         `json:"deviceCreateDate"`
	DeviceLastModifiedDate      time.Time         `json:"deviceLastModifiedDate"`
	DeviceLastAuthenticatedDate time.Time         `json:"deviceLastAuthenticatedDate,omitempty"`
	DeviceRememberedStatus      string            `json:"deviceRememberedStatus,omitempty"`
}

// AuthEvent represents a user authentication event recorded by Cognito.
type AuthEvent struct {
	EventID            string                  `json:"eventId"`
	UserName           string                  `json:"userName,omitempty"`
	ClientID           string                  `json:"clientId,omitempty"`
	UserPoolID         string                  `json:"userPoolId"`
	UserID             string                  `json:"userId"`
	EventType          string                  `json:"eventType"`
	CreationDate       time.Time               `json:"creationDate"`
	EventResponse      string                  `json:"eventResponse"`
	RiskDecision       string                  `json:"riskDecision,omitempty"`
	RiskLevel          string                  `json:"riskLevel,omitempty"`
	CompromisedFlag    bool                    `json:"compromisedCredentialsDetected,omitempty"`
	ChallengeResponses []ChallengeResponsePair `json:"challengeResponses,omitempty"`
	ContextIPAddress   string                  `json:"contextIpAddress,omitempty"`
	ContextDeviceName  string                  `json:"contextDeviceName,omitempty"`
	ContextCity        string                  `json:"contextCity,omitempty"`
	ContextCountry     string                  `json:"contextCountry,omitempty"`
	ContextTimezone    string                  `json:"contextTimezone,omitempty"`
	FeedbackValue      string                  `json:"feedbackValue,omitempty"`
	FeedbackProvider   string                  `json:"feedbackProvider,omitempty"`
	FeedbackDate       time.Time               `json:"feedbackDate,omitempty"`
}

// ChallengeResponsePair represents a single challenge-response entry in an auth event.
type ChallengeResponsePair struct {
	ChallengeName     string `json:"challengeName"`
	ChallengeResponse string `json:"challengeResponse"`
}

// ClientSecretDescriptor represents a client secret in multi-secret support.
type ClientSecretDescriptor struct {
	ClientSecretID         string    `json:"clientSecretId"`
	ClientSecretValue      string    `json:"clientSecretValue"`
	ClientSecretCreateDate time.Time `json:"clientSecretCreateDate"`
}

// LogDeliveryConfiguration stores Cognito log delivery settings for a user pool.
type LogDeliveryConfiguration struct {
	UserPoolID        string             `json:"userPoolId"`
	LogConfigurations []LogConfiguration `json:"logConfigurations,omitempty"`
}

// LogConfiguration represents a single log delivery rule.
type LogConfiguration struct {
	LogLevel                    string                `json:"logLevel"`
	EventSource                 string                `json:"eventSource"`
	CloudWatchLogsConfiguration *CloudWatchLogsConfig `json:"cloudWatchLogsConfiguration,omitempty"`
	S3Configuration             *S3Config             `json:"s3Configuration,omitempty"`
	FirehoseConfiguration       *FirehoseConfig       `json:"firehoseConfiguration,omitempty"`
}

// CloudWatchLogsConfig holds the CloudWatch Logs destination ARN.
type CloudWatchLogsConfig struct {
	LogGroupArn string `json:"logGroupArn,omitempty"`
}

// S3Config holds the S3 destination ARN.
type S3Config struct {
	BucketArn string `json:"bucketArn,omitempty"`
}

// FirehoseConfig holds the Firehose destination ARN.
type FirehoseConfig struct {
	StreamArn string `json:"streamArn,omitempty"`
}

// RiskConfiguration stores advanced security risk configuration for a user pool.
type RiskConfiguration struct {
	UserPoolID                        string    `json:"userPoolId"`
	ClientID                          string    `json:"clientId,omitempty"`
	CompromisedCredentialsEventFilter []string  `json:"ccEventFilter,omitempty"`
	CompromisedCredentialsEventAction string    `json:"ccEventAction,omitempty"`
	AccountTakeoverLowAction          string    `json:"atLowAction,omitempty"`
	AccountTakeoverLowNotify          bool      `json:"atLowNotify,omitempty"`
	AccountTakeoverMediumAction       string    `json:"atMediumAction,omitempty"`
	AccountTakeoverMediumNotify       bool      `json:"atMediumNotify,omitempty"`
	AccountTakeoverHighAction         string    `json:"atHighAction,omitempty"`
	AccountTakeoverHighNotify         bool      `json:"atHighNotify,omitempty"`
	NotifyFrom                        string    `json:"notifyFrom,omitempty"`
	NotifyReplyTo                     string    `json:"notifyReplyTo,omitempty"`
	NotifySourceArn                   string    `json:"notifySourceArn,omitempty"`
	NotifyBlockEmailSubject           string    `json:"notifyBlockEmailSubject,omitempty"`
	NotifyBlockEmailHtml              string    `json:"notifyBlockEmailHtml,omitempty"`
	NotifyNoActionEmailSubject        string    `json:"notifyNoActionEmailSubject,omitempty"`
	NotifyNoActionEmailHtml           string    `json:"notifyNoActionEmailHtml,omitempty"`
	NotifyMfaEmailSubject             string    `json:"notifyMfaEmailSubject,omitempty"`
	NotifyMfaEmailHtml                string    `json:"notifyMfaEmailHtml,omitempty"`
	BlockedIPRangeList                []string  `json:"blockedIpRangeList,omitempty"`
	SkippedIPRangeList                []string  `json:"skippedIpRangeList,omitempty"`
	LastModifiedDate                  time.Time `json:"lastModifiedDate,omitempty"`
}

// UICustomization stores CSS and image customisation for hosted UI.
type UICustomization struct {
	UserPoolID       string    `json:"userPoolId"`
	ClientID         string    `json:"clientId,omitempty"`
	CSS              string    `json:"css,omitempty"`
	CSSVersion       string    `json:"cssVersion,omitempty"`
	ImageFile        []byte    `json:"imageFile,omitempty"`
	CreationDate     time.Time `json:"creationDate,omitempty"`
	LastModifiedDate time.Time `json:"lastModifiedDate,omitempty"`
}

// SchemaAttributeType represents a custom schema attribute for a user pool.
type SchemaAttributeType struct {
	Name                       string                      `json:"name"`
	AttributeDataType          string                      `json:"attributeDataType"`
	DeveloperOnlyAttribute     bool                        `json:"developerOnlyAttribute,omitempty"`
	Mutable                    bool                        `json:"mutable,omitempty"`
	Required                   bool                        `json:"required,omitempty"`
	NumberAttributeConstraints *NumberAttributeConstraints `json:"numberAttributeConstraints,omitempty"`
	StringAttributeConstraints *StringAttributeConstraints `json:"stringAttributeConstraints,omitempty"`
}

// NumberAttributeConstraints defines min/max for number attributes.
type NumberAttributeConstraints struct {
	MinValue string `json:"minValue,omitempty"`
	MaxValue string `json:"maxValue,omitempty"`
}

// StringAttributeConstraints defines min/max length for string attributes.
type StringAttributeConstraints struct {
	MinLength string `json:"minLength,omitempty"`
	MaxLength string `json:"maxLength,omitempty"`
}

// UserImportJob represents a CSV user import job.
type UserImportJob struct {
	JobID                 string    `json:"jobId"`
	JobName               string    `json:"jobName"`
	UserPoolID            string    `json:"userPoolId"`
	PreSignedUrl          string    `json:"preSignedUrl,omitempty"`
	CreationDate          time.Time `json:"creationDate"`
	StartDate             time.Time `json:"startDate,omitempty"`
	CompletionDate        time.Time `json:"completionDate,omitempty"`
	Status                string    `json:"status"`
	CloudWatchLogsRoleArn string    `json:"cloudWatchLogsRoleArn,omitempty"`
	// PasswordHashingAlgorithm records the algorithm the CSV password_hash
	// column is encoded with (BCRYPT, SCRYPT, ARGON2ID, PBKDF2_SHA256).
	// Empty when the import carries no password hashes.
	PasswordHashingAlgorithm string `json:"passwordHashingAlgorithm,omitempty"`
	ImportedUsers            int64  `json:"importedUsers"`
	SkippedUsers             int64  `json:"skippedUsers"`
	FailedUsers              int64  `json:"failedUsers"`
	CompletionMessage        string `json:"completionMessage,omitempty"`
}

// WebAuthnCredential represents a registered FIDO2/WebAuthn credential.
type WebAuthnCredential struct {
	CredentialID            string    `json:"credentialId"`
	FriendlyName            string    `json:"friendlyName,omitempty"`
	UserPoolID              string    `json:"userPoolId"`
	UserID                  string    `json:"userId"`
	PublicKey               string    `json:"publicKey"`
	SignCount               uint32    `json:"signCount"`
	RelyingPartyID          string    `json:"relyingPartyId,omitempty"`
	AuthenticatorAttachment string    `json:"authenticatorAttachment,omitempty"`
	CreatedAt               time.Time `json:"createdAt"`
}

// ManagedLoginBranding stores branding configuration for managed login.
type ManagedLoginBranding struct {
	ManagedLoginBrandingId   string                 `json:"managedLoginBrandingId"`
	UserPoolID               string                 `json:"userPoolId"`
	ClientID                 string                 `json:"clientId,omitempty"`
	UseCognitoProvidedValues bool                   `json:"useCognitoProvidedValues"`
	Settings                 map[string]interface{} `json:"settings,omitempty"`
	Assets                   []BrandingAsset        `json:"assets,omitempty"`
	CreationDate             time.Time              `json:"creationDate,omitempty"`
	LastModifiedDate         time.Time              `json:"lastModifiedDate,omitempty"`
}

// BrandingAsset represents a branding asset (logo, colour, etc).
type BrandingAsset struct {
	Category  string `json:"category,omitempty"`
	Color     string `json:"color,omitempty"`
	Extension string `json:"extension,omitempty"`
	Bytes     string `json:"bytes,omitempty"`
}

// Terms represents a terms document for a user pool.
type Terms struct {
	TermsID          string                 `json:"termsId"`
	UserPoolID       string                 `json:"userPoolId"`
	ClientID         string                 `json:"clientId,omitempty"`
	TermsName        string                 `json:"termsName"`
	TermsSource      string                 `json:"termsSource,omitempty"`
	Enforcement      string                 `json:"enforcement,omitempty"`
	Links            map[string]interface{} `json:"links,omitempty"`
	CreationDate     time.Time              `json:"creationDate,omitempty"`
	LastModifiedDate time.Time              `json:"lastModifiedDate,omitempty"`
}

// UserPoolReplica represents a cross-region replica of a user pool.
type UserPoolReplica struct {
	UserPoolID   string      `json:"userPoolId"`
	RegionName   string      `json:"regionName"`
	Status       string      `json:"status"`
	Role         string      `json:"role"`
	UserPoolArn  string      `json:"userPoolArn,omitempty"`
	CreationDate time.Time   `json:"creationDate,omitempty"`
	Tags         []types.Tag `json:"tags,omitempty"`
}
