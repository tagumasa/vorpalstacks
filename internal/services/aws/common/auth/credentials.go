// Package auth provides AWS authentication utilities for vorpalstacks.
package auth

// Credentials represents AWS credentials for authentication.
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string // #nosec G117
	SessionToken    string // #nosec G117
	Region          string
	Service         string
}

// CredentialsProvider defines the interface for providing AWS credentials.
type CredentialsProvider interface {
	GetCredentials() (*Credentials, error)
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
