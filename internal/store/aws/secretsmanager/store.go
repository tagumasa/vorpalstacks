// Package secretsmanager provides Secrets Manager storage functionality for vorpalstacks.
package secretsmanager

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"

	"github.com/google/uuid"
)

// SecretStore provides Secrets Manager secret storage functionality.
type SecretStore struct {
	*common.BaseStore
	versionsStore *common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.RWMutex
}

// NewSecretStore creates a new SecretStore instance with the specified storage, account ID, and region.
func NewSecretStore(store storage.BasicStorage, accountID, region string) *SecretStore {
	return &SecretStore{
		BaseStore:     common.NewBaseStore(store.Bucket("secretsmanager-secrets-"+region), "secretsmanager-secrets"),
		versionsStore: common.NewBaseStore(store.Bucket("secretsmanager-versions-"+region), "secretsmanager-versions"),
		TagStore:      common.NewTagStoreWithRegion(store, "secretsmanager", region),
		arnBuilder:    svcarn.NewARNBuilder(accountID, region),
		accountID:     accountID,
		region:        region,
	}
}

// GetAccountID returns the AWS account ID associated with this secret store.
func (s *SecretStore) GetAccountID() string {
	return s.accountID
}

// GetRegion returns the AWS region associated with this secret store.
func (s *SecretStore) GetRegion() string {
	return s.region
}

// GetBaseStore returns the underlying base store for secrets.
func (s *SecretStore) GetBaseStore() *common.BaseStore {
	return s.BaseStore
}

func (s *SecretStore) buildSecretARN(name string) string {
	return s.arnBuilder.SecretsManager().Secret(name)
}

func (s *SecretStore) buildSecretKey(name string) string {
	return fmt.Sprintf("%s/%s", s.accountID, name)
}

func (s *SecretStore) buildVersionKey(secretName, versionId string) string {
	return fmt.Sprintf("%s/%s/%s", s.accountID, secretName, versionId)
}

// CreateSecret creates a new secret in the store with the specified details.
// Returns the created secret or an error if the secret already exists or validation fails.
func (s *SecretStore) CreateSecret(secret *Secret) (*Secret, error) {
	if secret.Name == "" {
		return nil, ErrInvalidSecretName
	}

	key := s.buildSecretKey(secret.Name)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(key) {
		return nil, ErrSecretAlreadyExists
	}

	now := time.Now().UTC()
	secret.ARN = s.buildSecretARN(secret.Name)
	secret.CreatedDate = now
	secret.LastChangedDate = now

	if secret.Tags == nil {
		secret.Tags = make(map[string]string)
	}

	if secret.VersionIDs == nil {
		secret.VersionIDs = []string{}
	}

	if secret.SecretString != "" || len(secret.SecretBinary) > 0 {
		versionId := generateVersionId()
		version := NewSecretVersion(versionId)
		version.SecretName = secret.Name
		version.SecretString = secret.SecretString
		version.SecretBinary = secret.SecretBinary

		if err := s.versionsStore.Put(s.buildVersionKey(secret.Name, versionId), version); err != nil {
			return nil, err
		}

		secret.VersionIDs = append(secret.VersionIDs, versionId)
		secret.CurrentVersion = versionId
	}

	if err := s.Put(key, secret); err != nil {
		return nil, err
	}

	if len(secret.Tags) > 0 {
		if err := s.TagStore.Tag(key, secret.Tags); err != nil {
			return nil, err
		}
	}

	return secret, nil
}

// GetSecret retrieves a secret by its name from the store.
// Returns the secret or an error if not found.
func (s *SecretStore) GetSecret(name string) (*Secret, error) {
	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return nil, ErrSecretNotFound
	}
	if secret.DeletedDate != nil || secret.ScheduledDeletionDate != nil {
		return nil, ErrSecretNotFound
	}
	tags, err := s.TagStore.List(key)
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	secret.Tags = tags
	return &secret, nil
}

// GetSecretByARN retrieves a secret by its ARN from the store.
// Returns the secret or an error if not found.
func (s *SecretStore) GetSecretByARN(arn string) (*Secret, error) {
	var found *Secret
	err := s.ForEach(func(key string, value []byte) error {
		if len(key) > 0 && key[0] == '#' {
			return nil
		}
		var secret Secret
		if err := json.Unmarshal(value, &secret); err != nil {
			return err
		}
		if secret.ARN == arn {
			if secret.DeletedDate != nil || secret.ScheduledDeletionDate != nil {
				return nil
			}
			found = &secret
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrSecretNotFound
	}
	key := s.buildSecretKey(found.Name)
	tags, err := s.TagStore.List(key)
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	found.Tags = tags
	return found, nil
}

// GetSecretForMetadata retrieves a secret by name including soft-deleted ones.
// Used for operations like DescribeSecret, RestoreSecret that work on deleted secrets.
func (s *SecretStore) GetSecretForMetadata(name string) (*Secret, error) {
	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return nil, ErrSecretNotFound
	}
	tags, err := s.TagStore.List(key)
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	secret.Tags = tags
	return &secret, nil
}

// UpdateSecret updates an existing secret in the store.
// Returns the updated secret or an error if the secret does not exist.
func (s *SecretStore) UpdateSecret(secret *Secret) (*Secret, error) {
	key := s.buildSecretKey(secret.Name)

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(key) {
		return nil, ErrSecretNotFound
	}

	secret.LastChangedDate = time.Now().UTC()

	if secret.SecretString != "" || len(secret.SecretBinary) > 0 {
		oldCurrentVersion := secret.CurrentVersion

		versionId := generateVersionId()
		version := NewSecretVersion(versionId)
		version.SecretName = secret.Name
		version.SecretString = secret.SecretString
		version.SecretBinary = secret.SecretBinary

		if err := s.versionsStore.Put(s.buildVersionKey(secret.Name, versionId), version); err != nil {
			return nil, err
		}

		if oldCurrentVersion != "" {
			oldVersion, err := s.GetSecretVersion(secret.Name, oldCurrentVersion)
			if err == nil {
				newStages := []string{}
				for _, stage := range oldVersion.VersionStages {
					if stage != "AWSCURRENT" {
						newStages = append(newStages, stage)
					}
				}
				if len(newStages) == 0 {
					newStages = []string{"AWSPREVIOUS"}
				} else {
					found := false
					for _, s := range newStages {
						if s == "AWSPREVIOUS" {
							found = true
							break
						}
					}
					if !found {
						newStages = append(newStages, "AWSPREVIOUS")
					}
				}
				oldVersion.VersionStages = newStages
				if err := s.versionsStore.Put(s.buildVersionKey(secret.Name, oldCurrentVersion), oldVersion); err != nil {
					return nil, err
				}
			}
		}

		secret.VersionIDs = append(secret.VersionIDs, versionId)
		secret.CurrentVersion = versionId
	}

	if err := s.Put(key, secret); err != nil {
		return nil, err
	}

	return secret, nil
}

// DeleteSecret deletes a secret by its name from the store.
// Returns an error if the secret does not exist.
func (s *SecretStore) DeleteSecret(name string) error {
	key := s.buildSecretKey(name)

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(key) {
		return ErrSecretNotFound
	}

	opts := common.ListOptions{MaxItems: 10000}
	result, err := common.List[SecretVersion](s.versionsStore, opts, func(v *SecretVersion) bool {
		return v.SecretName == name
	})
	if err != nil {
		return err
	}
	for _, version := range result.Items {
		if err := s.versionsStore.Delete(s.buildVersionKey(name, version.VersionId)); err != nil {
			return err
		}
	}

	if err := s.TagStore.Delete(key); err != nil {
		return err
	}

	return s.BaseStore.Delete(key)
}

// ListSecrets returns a list of secrets from the store using the specified list options.
func (s *SecretStore) ListSecrets(opts common.ListOptions) (*common.ListResult[Secret], error) {
	return common.List[Secret](s.BaseStore, opts, func(secret *Secret) bool {
		return secret.DeletedDate == nil && secret.ScheduledDeletionDate == nil
	})
}

// GetSecretVersion retrieves a specific version of a secret by name and version ID.
// If version ID is empty, retrieves the current version.
func (s *SecretStore) GetSecretVersion(name, versionId string) (*SecretVersion, error) {
	if versionId == "" {
		secret, err := s.GetSecret(name)
		if err != nil {
			return nil, err
		}
		versionId = secret.CurrentVersion
	}

	key := s.buildVersionKey(name, versionId)
	var version SecretVersion
	if err := s.versionsStore.Get(key, &version); err != nil {
		return nil, ErrInvalidVersionId
	}
	return &version, nil
}

// GetSecretVersionByStage retrieves a version of a secret by its stage (AWSCURRENT, AWSPREVIOUS, etc.)
func (s *SecretStore) GetSecretVersionByStage(name, stage string) (*SecretVersion, error) {
	opts := common.ListOptions{MaxItems: 10000}
	result, err := common.List[SecretVersion](s.versionsStore, opts, func(v *SecretVersion) bool {
		if v.SecretName != name {
			return false
		}
		for _, s := range v.VersionStages {
			if s == stage {
				return true
			}
		}
		return false
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list secret versions: %w", err)
	}
	if len(result.Items) == 0 {
		return nil, ErrInvalidVersionId
	}
	return result.Items[0], nil
}

// ListSecretTags retrieves the tags associated with a secret.
func (s *SecretStore) ListSecretTags(name string) (map[string]string, error) {
	key := s.buildSecretKey(name)
	return s.TagStore.List(key)
}

// TagSecret adds tags to a secret.
func (s *SecretStore) TagSecret(name string, tags map[string]string) error {
	key := s.buildSecretKey(name)
	return s.TagStore.Tag(key, tags)
}

// UntagSecret removes tags from a secret.
func (s *SecretStore) UntagSecret(name string, tagKeys []string) error {
	key := s.buildSecretKey(name)
	return s.TagStore.Untag(key, tagKeys)
}

func generateVersionId() string {
	return uuid.New().String()[:32]
}

// CancelRotation cancels the rotation of a secret.
func (s *SecretStore) CancelRotation(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return ErrSecretNotFound
	}
	secret.RotationEnabled = false
	secret.RotationLambdaARN = ""
	secret.RotationRules = nil
	return s.Put(key, &secret)
}

// GetResourcePolicy retrieves the resource policy for a secret.
func (s *SecretStore) GetResourcePolicy(name string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return "", ErrSecretNotFound
	}
	return secret.ResourcePolicy, nil
}

// PutResourcePolicy sets the resource policy for a secret.
func (s *SecretStore) PutResourcePolicy(name, policy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return ErrSecretNotFound
	}
	secret.ResourcePolicy = policy
	return s.Put(key, &secret)
}

// DeleteResourcePolicy deletes the resource policy for a secret.
func (s *SecretStore) DeleteResourcePolicy(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return ErrSecretNotFound
	}
	secret.ResourcePolicy = ""
	return s.Put(key, &secret)
}

// RestoreSecret restores a previously deleted secret.
func (s *SecretStore) RestoreSecret(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return ErrSecretNotFound
	}
	secret.DeletedDate = nil
	secret.ScheduledDeletionDate = nil
	return s.Put(key, &secret)
}

// RotateSecret rotates the secret by updating its last rotated date.
func (s *SecretStore) RotateSecret(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return ErrSecretNotFound
	}
	secret.LastRotatedDate = time.Now().UTC()
	return s.Put(key, &secret)
}

// ScheduleDeletion schedules a secret for deletion.
func (s *SecretStore) ScheduleDeletion(name string, deletionDate time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return ErrSecretNotFound
	}
	now := time.Now().UTC()
	secret.DeletedDate = &now
	secret.ScheduledDeletionDate = &deletionDate
	return s.Put(key, &secret)
}

// ListSecretVersions returns all versions of a secret.
func (s *SecretStore) ListSecretVersions(name string) ([]SecretVersion, error) {
	secret, err := s.GetSecretForMetadata(name)
	if err != nil {
		return nil, err
	}

	var versions []SecretVersion
	for _, versionId := range secret.VersionIDs {
		key := s.buildVersionKey(name, versionId)
		var version SecretVersion
		if err := s.versionsStore.Get(key, &version); err == nil {
			versions = append(versions, version)
		}
	}

	return versions, nil
}

// UpdateSecretVersionStage updates the version stages for a specific secret version.
func (s *SecretStore) UpdateSecretVersionStage(name, versionId string, stages []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildVersionKey(name, versionId)
	var version SecretVersion
	if err := s.versionsStore.Get(key, &version); err != nil {
		return ErrInvalidVersionId
	}
	version.VersionStages = stages
	return s.versionsStore.Put(key, &version)
}

// UpdateSecretMetadata updates the secret metadata without creating a new version.
func (s *SecretStore) UpdateSecretMetadata(secret *Secret) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildSecretKey(secret.Name)
	return s.Put(key, secret)
}

// FinishRotation atomically promotes AWSPENDING to AWSCURRENT and demotes the
// old AWSCURRENT to AWSPREVIOUS.  All stage-transition writes and the metadata
// update execute under a single mutex lock so that a crash mid-rotation cannot
// leave no version carrying the AWSCURRENT stage.
func (s *SecretStore) FinishRotation(secret *Secret, pendingVersionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	versionKey := s.buildVersionKey(secret.Name, pendingVersionID)
	var pendingVersion SecretVersion
	if err := s.versionsStore.Get(versionKey, &pendingVersion); err != nil {
		return fmt.Errorf("AWSPENDING version %s not found: %w", pendingVersionID, err)
	}

	oldPrevious, prevErr := s.getSecretVersionByStageLocked(secret.Name, "AWSPREVIOUS")
	if prevErr == nil && oldPrevious.VersionId != pendingVersionID {
		cleanedStages := make([]string, 0, len(oldPrevious.VersionStages))
		for _, st := range oldPrevious.VersionStages {
			if st != "AWSPREVIOUS" {
				cleanedStages = append(cleanedStages, st)
			}
		}
		prevKey := s.buildVersionKey(secret.Name, oldPrevious.VersionId)
		oldPrevious.VersionStages = cleanedStages
		if err := s.versionsStore.Put(prevKey, &oldPrevious); err != nil {
			return fmt.Errorf("failed to clean AWSPREVIOUS from old previous version during rotation: %w", err)
		}
	}

	oldCurrentID := secret.CurrentVersion
	if oldCurrentID != "" && oldCurrentID != pendingVersionID {
		currentKey := s.buildVersionKey(secret.Name, oldCurrentID)
		var currentVersion SecretVersion
		if err := s.versionsStore.Get(currentKey, &currentVersion); err != nil {
			return fmt.Errorf("failed to read current version during rotation: %w", err)
		}
		currentVersion.VersionStages = []string{"AWSPREVIOUS"}
		if err := s.versionsStore.Put(currentKey, &currentVersion); err != nil {
			return fmt.Errorf("failed to demote current version to AWSPREVIOUS during rotation: %w", err)
		}
	}

	pendingVersion.VersionStages = []string{"AWSCURRENT"}
	if err := s.versionsStore.Put(versionKey, &pendingVersion); err != nil {
		return fmt.Errorf("failed to promote pending version to AWSCURRENT during rotation: %w", err)
	}

	secret.CurrentVersion = pendingVersionID
	secretKey := s.buildSecretKey(secret.Name)
	if err := s.Put(secretKey, secret); err != nil {
		return fmt.Errorf("failed to update secret metadata after rotation: %w", err)
	}

	return nil
}

// getSecretVersionByStageLocked reads a version by stage without acquiring the
// mutex (caller must already hold s.mu).
func (s *SecretStore) getSecretVersionByStageLocked(name, stage string) (*SecretVersion, error) {
	secret, err := s.getSecretForMetadataLocked(name)
	if err != nil {
		return nil, err
	}
	for _, versionId := range secret.VersionIDs {
		key := s.buildVersionKey(name, versionId)
		var version SecretVersion
		if err := s.versionsStore.Get(key, &version); err != nil {
			continue
		}
		for _, st := range version.VersionStages {
			if st == stage {
				return &version, nil
			}
		}
	}
	return nil, fmt.Errorf("no version with stage %s", stage)
}

// getSecretForMetadataLocked reads secret metadata without acquiring the mutex.
func (s *SecretStore) getSecretForMetadataLocked(name string) (*Secret, error) {
	key := s.buildSecretKey(name)
	var secret Secret
	if err := s.Get(key, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

// CreateVersionDirect creates a secret version with a specific version ID (used for replication).
func (s *SecretStore) CreateVersionDirect(secretName string, version *SecretVersion) error {
	if version.VersionId == "" {
		version.VersionId = generateVersionId()
	}
	key := s.buildVersionKey(secretName, version.VersionId)
	return s.versionsStore.Put(key, version)
}
