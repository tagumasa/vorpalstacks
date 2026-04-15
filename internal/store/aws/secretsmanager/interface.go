package secretsmanager

import (
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// SecretStoreInterface defines operations for managing Secrets Manager secrets.
type SecretStoreInterface interface {
	GetAccountID() string
	GetRegion() string
	GetBaseStore() *common.BaseStore
	CreateSecret(secret *Secret) (*Secret, error)
	GetSecret(name string) (*Secret, error)
	GetSecretByARN(arn string) (*Secret, error)
	GetSecretForMetadata(name string) (*Secret, error)
	UpdateSecret(secret *Secret) (*Secret, error)
	DeleteSecret(name string) error
	ListSecrets(opts common.ListOptions) (*common.ListResult[Secret], error)
	GetSecretVersion(name, versionId string) (*SecretVersion, error)
	GetSecretVersionByStage(name, stage string) (*SecretVersion, error)
	ListSecretTags(name string) (map[string]string, error)
	TagSecret(name string, tags map[string]string) error
	UntagSecret(name string, tagKeys []string) error
	CancelRotation(name string) error
	GetResourcePolicy(name string) (string, error)
	PutResourcePolicy(name, policy string) error
	DeleteResourcePolicy(name string) error
	RestoreSecret(name string) error
	RotateSecret(name string) error
	ScheduleDeletion(name string, deletionDate time.Time) error
	ListSecretVersions(name string) ([]SecretVersion, error)
	UpdateSecretVersionStage(name, versionId string, stages []string) error
	UpdateSecretMetadata(secret *Secret) error
	FinishRotation(secret *Secret, pendingVersionID string) error
	CreateVersionDirect(secretName string, version *SecretVersion) error
}

var _ SecretStoreInterface = (*SecretStore)(nil)
