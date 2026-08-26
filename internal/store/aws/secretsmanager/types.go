// Package secretsmanager provides Secrets Manager storage functionality for vorpalstacks.
package secretsmanager

import (
	"time"
)

// Secret represents a Secrets Manager secret.
//
// SecretString and SecretBinary are transient fields (json:"-") used to pass
// the initial or updated secret value from the service handler to the store.
// The persisted value lives exclusively in SecretVersion, never in the Secret
// metadata — this prevents accidental exposure of secret material through
// DescribeSecret or ListSecrets.
//
// InitialVersionId is likewise transient (json:"-"); it carries the
// ClientRequestToken so the store can use it as the version ID during
// CreateSecret or UpdateSecret.
type Secret struct {
	Name              string              `json:"name"`
	ARN               string              `json:"arn"`
	SecretString      string              `json:"-"`
	SecretBinary      []byte              `json:"-"`
	VersionIDs        []string            `json:"versionIds"`
	CurrentVersion    string              `json:"currentVersion"`
	CreatedDate       time.Time           `json:"createdDate"`
	LastChangedDate   time.Time           `json:"lastChangedDate"`
	LastAccessedDate  time.Time           `json:"lastAccessedDate"`
	Description       string              `json:"description,omitempty"`
	KmsKeyId          string              `json:"kmsKeyId,omitempty"`
	Tags              map[string]string   `json:"tags"`
	ReplicationStatus []ReplicationStatus `json:"replicationStatus,omitempty"`
	RotationEnabled   bool                `json:"rotationEnabled"`
	RotationLambdaARN string              `json:"rotationLambdaARN,omitempty"`
	RotationRules     *RotationRules      `json:"rotationRules,omitempty"`
	LastRotatedDate   time.Time           `json:"lastRotatedDate,omitempty"`
	ResourcePolicy    string              `json:"resourcePolicy,omitempty"`
	DeletedDate       *time.Time          `json:"deletedDate,omitempty"`
	Type              string              `json:"type,omitempty"`
	OwningService     string              `json:"owningService,omitempty"`
	PrimaryRegion     string              `json:"primaryRegion,omitempty"`
	// Managed external secret rotation metadata (partner-supplied) and the
	// role ARN that permits rotating a partner-held secret; both are
	// configuration storage with no behavioural consumer on-platform.
	ExternalSecretRotationMetadata []ExternalSecretRotationMetadataItem `json:"externalSecretRotationMetadata,omitempty"`
	ExternalSecretRotationRoleArn  string                               `json:"externalSecretRotationRoleArn,omitempty"`
	NextRotationDate               time.Time                            `json:"nextRotationDate,omitempty"`
	InitialVersionId               string                               `json:"-"`
}

// RotationRules defines the rotation rules for a secret.
type RotationRules struct {
	AutomaticallyAfterDays int    `json:"automaticallyAfterDays,omitempty"`
	Duration               string `json:"duration,omitempty"`
	ScheduleExpression     string `json:"scheduleExpression,omitempty"`
}

// ExternalSecretRotationMetadataItem is one partner-supplied key/value pair
// of the managed external secret rotation metadata.
type ExternalSecretRotationMetadataItem struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// SecretVersion represents a version of a Secrets Manager secret.
type SecretVersion struct {
	VersionId        string    `json:"versionId"`
	SecretName       string    `json:"secretName"`
	SecretString     string    `json:"secretString,omitempty"`
	SecretBinary     []byte    `json:"secretBinary,omitempty"`
	VersionStages    []string  `json:"versionStages"`
	CreatedDate      time.Time `json:"createdDate"`
	LastAccessedDate time.Time `json:"lastAccessedDate"`
	// KmsKeyIds lists the KMS keys used to encrypt this version.  On
	// vorpalstacks secrets are encrypted with a single platform key so
	// this slice is typically empty in stored data, but the field exists
	// for API output compatibility.
	KmsKeyIds []string `json:"kmsKeyIds,omitempty"`
}

// ReplicationStatus represents the replication status of a secret.
type ReplicationStatus struct {
	Region           string    `json:"region"`
	KmsKeyId         string    `json:"kmsKeyId"`
	Status           string    `json:"status"`
	StatusMessage    string    `json:"statusMessage"`
	LastAccessedDate time.Time `json:"lastAccessedDate"`
}

// NewSecret creates a new Secret with the specified name and default values.
func NewSecret(name string) *Secret {
	now := time.Now().UTC()
	return &Secret{
		Name:            name,
		VersionIDs:      []string{},
		Tags:            make(map[string]string),
		CreatedDate:     now,
		LastChangedDate: now,
	}
}

// NewSecretVersion creates a new SecretVersion with the specified version ID and default values.
func NewSecretVersion(versionId string) *SecretVersion {
	return &SecretVersion{
		VersionId:     versionId,
		VersionStages: []string{"AWSCURRENT"},
		CreatedDate:   time.Now().UTC(),
	}
}
