package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// type definitions for credentials and signing.

import "time"

// Credentials holds AWS credentials for API requests.
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string // #nosec G117
	SessionToken    string // #nosec G117
	Region          string
	Service         string
	Expiration      *time.Time
}

// CredentialsProvider defines an interface for providing credentials.
type CredentialsProvider interface {
	GetCredentials() (*Credentials, error)
}

// StaticCredentialsProvider provides static credentials.
type StaticCredentialsProvider struct {
	credentials *Credentials
}

// NewStaticCredentialsProvider creates a new static credentials provider.
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

// AuthorizationHeader represents the components of an AWS authorization header.
type AuthorizationHeader struct {
	Algorithm     string
	Credential    string
	SignedHeaders string
	Signature     string
}
