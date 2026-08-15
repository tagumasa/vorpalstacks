// Package auth provides AWS authentication utilities for vorpalstacks.
package auth

import (
	"errors"
	"time"
)

// ErrSessionExpired is a sentinel error indicating that the caller's
// temporary credentials (STS session) have expired. The auth middleware
// uses errors.Is to detect this and return an ExpiredTokenException
// response so that AWS SDK clients can trigger credential refresh.
var ErrSessionExpired = errors.New("session expired")

// Credentials represents AWS credentials for authentication.
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string // #nosec G117
	SessionToken    string // #nosec G117
	Region          string
	Service         string
	Expiration      *time.Time
}

// CredentialsProvider defines the interface for providing AWS credentials.
type CredentialsProvider interface {
	GetCredentials() (*Credentials, error)
}

// SessionCredentials represents temporary session credentials from STS.
// In addition to the bare credentials, callers (notably the request
// authoriser) receive the session-scoped principal identity, tags and
// source identity so that ABAC condition keys and sts:SourceIdentity
// references can be resolved at request time.
type SessionCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string

	// PrincipalArn is the ARN of the principal that owns this session
	// (AssumedRole, SAML, WebIdentity, FederatedUser, Root, etc.).
	PrincipalArn string
	// PrincipalType mirrors the session store's PrincipalType field so
	// that consumers can distinguish role sessions from federated
	// sessions without re-reading the underlying Session record.
	PrincipalType string
	// PrincipalName mirrors the session store's PrincipalName field. For
	// GetSessionToken sessions created with the root user's permanent key
	// it is the root user name, which lets the authoriser recognise
	// legitimate root sessions instead of trusting the ARN suffix alone.
	PrincipalName string
	// Tags carries the caller-supplied session tags and any transitive
	// tags forwarded from a previous role session. They populate the
	// EvaluationContext.SessionContext map for policy evaluation.
	Tags map[string]string
	// SourceIdentity is the caller-supplied SourceIdentity value, which
	// surfaces as the aws:SourceIdentity / sts:SourceIdentity condition
	// variable in policy evaluation.
	SourceIdentity string
	// Policy and PolicyArns carry the session-scoping policy through to
	// the authorisation layer, where it is intersected with the
	// assumed role's identity-based policies.
	Policy     string
	PolicyArns []string
}

// SessionResolver resolves temporary session credentials by access key ID.
// Implementations typically look up STS sessions.
type SessionResolver interface {
	ResolveSession(accessKeyId string) (*SessionCredentials, error)
}

// StaticCredentialsProvider provides static credentials for authentication.
type StaticCredentialsProvider struct {
	credentials *Credentials
}

// NewStaticCredentialsProvider creates a new static credentials provider with the specified credentials.
func NewStaticCredentialsProvider(accessKeyID, secretAccessKey, region, service string) *StaticCredentialsProvider {
	return &StaticCredentialsProvider{
		credentials: &Credentials{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
			Region:          region,
			Service:         service,
		},
	}
}

// GetCredentials returns the static credentials.
func (p *StaticCredentialsProvider) GetCredentials() (*Credentials, error) {
	return p.credentials, nil
}
