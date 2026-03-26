// Package cognitoidentity provides Cognito Identity Pool storage functionality for vorpalstacks.
package cognitoidentity

import (
	"time"

	"vorpalstacks/internal/store/aws/common"

	"github.com/google/uuid"
)

// IdentityPool represents an Amazon Cognito Identity Pool.
type IdentityPool struct {
	ID                             string                    `json:"id"`
	Name                           string                    `json:"name"`
	Arn                            string                    `json:"arn"`
	AllowUnauthenticatedIdentities bool                      `json:"allowUnauthenticatedIdentities"`
	AllowClassicFlow               bool                      `json:"allowClassicFlow,omitempty"`
	CognitoIdentityProviders       []CognitoIdentityProvider `json:"cognitoIdentityProviders,omitempty"`
	DeveloperProviderName          string                    `json:"developerProviderName,omitempty"`
	SupportedLoginProviders        map[string]string         `json:"supportedLoginProviders,omitempty"`
	OpenIdConnectProviderARNs      []string                  `json:"openIdConnectProviderARNs,omitempty"`
	SamlProviderARNs               []string                  `json:"samlProviderARNs,omitempty"`
	Tags                           map[string]string         `json:"tags,omitempty"`
	CreationDate                   time.Time                 `json:"creationDate"`
	LastModifiedDate               time.Time                 `json:"lastModifiedDate"`
	AuthenticatedRoleArn           string                    `json:"authenticatedRoleArn,omitempty"`
	UnauthenticatedRoleArn         string                    `json:"unauthenticatedRoleArn,omitempty"`
	RoleMappings                   map[string]RoleMapping    `json:"roleMappings,omitempty"`
}

// Identity represents an Amazon Cognito Identity.
type Identity struct {
	ID               string            `json:"id"`
	IdentityPoolID   string            `json:"identityPoolId"`
	Logins           map[string]string `json:"logins,omitempty"`
	CreationDate     time.Time         `json:"creationDate"`
	LastModifiedDate time.Time         `json:"lastModifiedDate"`
}

// CognitoIdentityProvider represents a Cognito Identity Provider configuration.
type CognitoIdentityProvider struct {
	ProviderName         string `json:"providerName"`
	ClientID             string `json:"clientId"`
	ServerSideTokenCheck bool   `json:"serverSideTokenCheck"`
}

// RoleMapping represents a role mapping configuration for an Identity Pool.
type RoleMapping struct {
	Type                    string              `json:"type"`
	AmbiguousRoleResolution string              `json:"ambiguousRoleResolution,omitempty"`
	RulesConfiguration      *RulesConfiguration `json:"rulesConfiguration,omitempty"`
}

// RulesConfiguration represents the rules configuration for a role mapping.
type RulesConfiguration struct {
	Rules []MappingRule `json:"rules,omitempty"`
}

// MappingRule represents a single mapping rule for role resolution.
type MappingRule struct {
	Claim     string `json:"claim"`
	MatchType string `json:"matchType"`
	Value     string `json:"value"`
	RoleARN   string `json:"roleArn"`
}

// NewIdentityPool creates a new Identity Pool with the specified name and configuration.
func NewIdentityPool(name string, allowUnauthenticated bool, region string) *IdentityPool {
	now := time.Now().UTC()
	id := generateIdentityPoolID(region)
	return &IdentityPool{
		ID:                             id,
		Name:                           name,
		AllowUnauthenticatedIdentities: allowUnauthenticated,
		AllowClassicFlow:               false,
		CognitoIdentityProviders:       []CognitoIdentityProvider{},
		SupportedLoginProviders:        make(map[string]string),
		OpenIdConnectProviderARNs:      []string{},
		SamlProviderARNs:               []string{},
		Tags:                           make(map[string]string),
		CreationDate:                   now,
		LastModifiedDate:               now,
		RoleMappings:                   make(map[string]RoleMapping),
	}
}

// NewIdentity creates a new Identity for the specified Identity Pool.
func NewIdentity(identityPoolID string) *Identity {
	now := time.Now().UTC()
	return &Identity{
		ID:               generateIdentityID(),
		IdentityPoolID:   identityPoolID,
		Logins:           make(map[string]string),
		CreationDate:     now,
		LastModifiedDate: now,
	}
}

func generateIdentityPoolID(region string) string {
	return region + ":" + uuid.New().String()
}

func generateIdentityID() string {
	return uuid.New().String()
}

// Tag represents a tag for an Amazon Cognito Identity Pool.
type Tag = common.Tag
