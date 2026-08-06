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
		if secret.InitialVersionId != "" {
			versionId = secret.InitialVersionId
		}
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

	if err := s.writeARNIndex(secret.ARN, secret.Name); err != nil {
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
	if secret.DeletedDate != nil {
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

// arnIndexEntry is the JSON structure stored for ARN index entries.
// The field uses a distinct JSON key "_arn_name" that does not overlap with
// any Secret field, so when a ForEach/List scan deserialises the value as a
// Secret the result has Name="" and is filtered out.
type arnIndexEntry struct {
	Name string `json:"_arn_name"`
}

// arnIndexKey returns the index key for a given ARN. The prefix "#" marks it
// as an index entry so that ForEach scans skip it automatically.
func (s *SecretStore) arnIndexKey(arn string) string {
	return "#arn:" + arn
}

// writeARNIndex creates a secondary index entry mapping ARN → secret name.
func (s *SecretStore) writeARNIndex(arn, name string) error {
	return s.Put(s.arnIndexKey(arn), &arnIndexEntry{Name: name})
}

// lookupNameByARN resolves a secret name from its ARN via the secondary index.
// Falls back to a full scan for secrets created before the index was added.
func (s *SecretStore) lookupNameByARN(arn string) (string, error) {
	var entry arnIndexEntry
	if err := s.Get(s.arnIndexKey(arn), &entry); err == nil && entry.Name != "" {
		return entry.Name, nil
	}
	var name string
	err := s.ForEach(func(key string, value []byte) error {
		if len(key) > 0 && key[0] == '#' {
			return nil
		}
		var secret Secret
		if err := json.Unmarshal(value, &secret); err != nil {
			return err
		}
		if secret.ARN == arn {
			name = secret.Name
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", ErrSecretNotFound
	}
	return name, nil
}

// LookupNameByARN resolves a secret name from its ARN via the secondary index.
// Unlike GetSecretByARN this does not check soft-deletion status.
func (s *SecretStore) LookupNameByARN(arn string) (string, error) {
	return s.lookupNameByARN(arn)
}

// GetSecretByARN retrieves a secret by its ARN from the store.
// Uses a secondary index for O(1) lookup with a full-scan fallback.
func (s *SecretStore) GetSecretByARN(arn string) (*Secret, error) {
	name, err := s.lookupNameByARN(arn)
	if err != nil {
		return nil, err
	}
	secret, err := s.GetSecret(name)
	if err != nil {
		// The secret may be soft-deleted; try metadata lookup.
		secret, err = s.GetSecretForMetadata(name)
		if err != nil {
			return nil, err
		}
		if secret.DeletedDate != nil {
			return nil, ErrSecretNotFound
		}
	}
	return secret, nil
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
		if secret.InitialVersionId != "" {
			versionId = secret.InitialVersionId
		}
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

	var secret Secret
	if err := s.BaseStore.Get(key, &secret); err != nil {
		return ErrSecretNotFound
	}

	for _, versionId := range secret.VersionIDs {
		if err := s.versionsStore.Delete(s.buildVersionKey(name, versionId)); err != nil {
			return err
		}
	}

	if err := s.TagStore.Delete(key); err != nil {
		return err
	}

	if err := s.BaseStore.Delete(key); err != nil {
		return err
	}

	// Clean up the ARN index entry to prevent orphan accumulation.
	_ = s.BaseStore.Delete(s.arnIndexKey(secret.ARN))
	return nil
}

// ListSecrets returns a list of secrets from the store using the specified list options.
// The filter callback is applied inside the Pebble iterator so that filtered-out items
// do not count towards MaxItems, enabling true store-level pagination even with filters.
// Pass nil to accept all named secrets.
func (s *SecretStore) ListSecrets(opts common.ListOptions, filter func(*Secret) bool) (*common.ListResult[Secret], error) {
	return common.List[Secret](s.BaseStore, opts, func(secret *Secret) bool {
		if secret.Name == "" {
			return false
		}
		if filter != nil && !filter(secret) {
			return false
		}
		return true
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
	secret, err := s.GetSecretForMetadata(name)
	if err != nil {
		return nil, err
	}
	for _, versionId := range secret.VersionIDs {
		var version SecretVersion
		if err := s.versionsStore.Get(s.buildVersionKey(name, versionId), &version); err != nil {
			continue
		}
		for _, st := range version.VersionStages {
			if st == stage {
				return &version, nil
			}
		}
	}
	return nil, ErrInvalidVersionId
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
	// AWS updates LastChangedDate when a secret is restored.
	secret.LastChangedDate = time.Now().UTC()
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
	// AWS spec: DeletedDate in DescribeSecret = the scheduled deletion date
	// (request time + RecoveryWindowInDays), not the time the request was made.
	scheduled := deletionDate.UTC()
	secret.DeletedDate = &scheduled
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

// MoveStage atomically moves a staging label from one version to another.
// All version stage reads and writes execute under a single mutex lock so that
// no concurrent request can observe an intermediate state.
func (s *SecretStore) MoveStage(secret *Secret, versionStage, moveToVersionId, removeFromVersionId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if removeFromVersionId != "" && removeFromVersionId != moveToVersionId {
		removeKey := s.buildVersionKey(secret.Name, removeFromVersionId)
		var removeVer SecretVersion
		if err := s.versionsStore.Get(removeKey, &removeVer); err != nil {
			return fmt.Errorf("failed to read remove-from version during MoveStage: %w", err)
		}
		newStages := make([]string, 0, len(removeVer.VersionStages))
		for _, st := range removeVer.VersionStages {
			if st != versionStage {
				newStages = append(newStages, st)
			}
		}
		removeVer.VersionStages = newStages
		if err := s.versionsStore.Put(removeKey, &removeVer); err != nil {
			return fmt.Errorf("failed to remove stage from version during MoveStage: %w", err)
		}
	}

	targetKey := s.buildVersionKey(secret.Name, moveToVersionId)
	var targetVer SecretVersion
	if err := s.versionsStore.Get(targetKey, &targetVer); err != nil {
		return fmt.Errorf("failed to read target version during MoveStage: %w", err)
	}

	found := false
	for _, st := range targetVer.VersionStages {
		if st == versionStage {
			found = true
			break
		}
	}
	if !found {
		targetVer.VersionStages = append(targetVer.VersionStages, versionStage)
	}

	if versionStage == "AWSCURRENT" {
		oldPrevious, prevErr := s.getSecretVersionByStageLocked(secret.Name, "AWSPREVIOUS")
		if prevErr == nil && oldPrevious.VersionId != moveToVersionId {
			prevStages := make([]string, 0, len(oldPrevious.VersionStages))
			for _, st := range oldPrevious.VersionStages {
				if st != "AWSPREVIOUS" {
					prevStages = append(prevStages, st)
				}
			}
			prevKey := s.buildVersionKey(secret.Name, oldPrevious.VersionId)
			oldPrevious.VersionStages = prevStages
			if err := s.versionsStore.Put(prevKey, &oldPrevious); err != nil {
				return fmt.Errorf("failed to clean AWSPREVIOUS from old previous version during MoveStage: %w", err)
			}
		}

		if removeFromVersionId != "" && removeFromVersionId != moveToVersionId {
			oldCurrentKey := s.buildVersionKey(secret.Name, removeFromVersionId)
			var oldCurrentVer SecretVersion
			if err := s.versionsStore.Get(oldCurrentKey, &oldCurrentVer); err == nil {
				hasPrev := false
				for _, st := range oldCurrentVer.VersionStages {
					if st == "AWSPREVIOUS" {
						hasPrev = true
						break
					}
				}
				if !hasPrev {
					oldCurrentVer.VersionStages = append(oldCurrentVer.VersionStages, "AWSPREVIOUS")
					if err := s.versionsStore.Put(oldCurrentKey, &oldCurrentVer); err != nil {
						return fmt.Errorf("failed to add AWSPREVIOUS to old current version during MoveStage: %w", err)
					}
				}
			}
		}

		secret.CurrentVersion = moveToVersionId
		secret.LastChangedDate = time.Now().UTC()
		secretKey := s.buildSecretKey(secret.Name)
		if err := s.Put(secretKey, secret); err != nil {
			return fmt.Errorf("failed to update secret metadata during MoveStage: %w", err)
		}
	}

	if err := s.versionsStore.Put(targetKey, &targetVer); err != nil {
		return fmt.Errorf("failed to add stage to target version during MoveStage: %w", err)
	}

	return nil
}

// CreateVersionDirect creates a secret version with a specific version ID (used for replication).
func (s *SecretStore) CreateVersionDirect(secretName string, version *SecretVersion) error {
	if version.VersionId == "" {
		version.VersionId = generateVersionId()
	}
	key := s.buildVersionKey(secretName, version.VersionId)
	return s.versionsStore.Put(key, version)
}
