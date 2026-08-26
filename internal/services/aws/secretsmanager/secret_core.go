package secretsmanager

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode"

	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// ---------------------------------------------------------------------------
// Input structs — transport-agnostic DTOs shared by HTTP API and admin handler
// ---------------------------------------------------------------------------

// SecretFilter is the transport-agnostic form of one ListSecrets/BatchGetSecretValue
// Filter entry.
type SecretFilter struct {
	Key    string
	Values []string
}

// CreateSecretInput carries all fields needed for CreateSecret.
type CreateSecretInput struct {
	Name         string
	SecretString string
	SecretBinary []byte
	// SecretBinaryB64 carries the HTTP wire form (base64 text); it is decoded
	// and validated in the Core. The admin plane fills SecretBinary directly
	// because proto bytes need no decoding.
	SecretBinaryB64             string
	Description                 string
	KmsKeyId                    string
	Type                        string
	ClientRequestToken          string
	Tags                        map[string]string
	AddReplicaRegions           []replicaRegion
	ForceOverwriteReplicaSecret bool
	Region                      string
}

// GetSecretValueInput carries all fields needed for GetSecretValue.
type GetSecretValueInput struct {
	SecretId     string
	VersionId    string
	VersionStage string
}

// UpdateSecretInput carries all fields needed for UpdateSecret.
type UpdateSecretInput struct {
	SecretId           string
	SecretString       string
	SecretBinaryB64    string
	Description        string
	KmsKeyId           string
	Type               string
	ClientRequestToken string
	// Has* flags preserve the explicit-presence semantics of the wire
	// parameters: an absent Description must not overwrite the stored value.
	HasDescription bool
	HasKmsKeyId    bool
	HasType        bool
}

// DeleteSecretInput carries all fields needed for DeleteSecret.
type DeleteSecretInput struct {
	SecretId                   string
	ForceDeleteWithoutRecovery bool
	RecoveryWindowInDays       int
	HasRecoveryWindow          bool
}

// ListSecretsInput carries pagination, sorting and filter parameters for
// ListSecrets.
type ListSecretsInput struct {
	MaxResults             *int
	NextToken              string
	SortBy                 string
	SortOrder              string
	Filters                []SecretFilter
	IncludePlannedDeletion bool
}

// DescribeSecretInput carries all fields needed for DescribeSecret.
type DescribeSecretInput struct {
	SecretId string
}

// ---------------------------------------------------------------------------
// Result structs — transport-agnostic results
// ---------------------------------------------------------------------------

// CreateSecretResult holds the transport-agnostic result of CreateSecret.
type CreateSecretResult struct {
	ARN               string
	Name              string
	VersionID         string
	ReplicationStatus []secretsmanagerstore.ReplicationStatus
}

// GetSecretValueResult holds the transport-agnostic result of GetSecretValue.
type GetSecretValueResult struct {
	Secret  *secretsmanagerstore.Secret
	Version *secretsmanagerstore.SecretVersion
}

// UpdateSecretResult holds the transport-agnostic result of UpdateSecret.
type UpdateSecretResult struct {
	ARN       string
	Name      string
	VersionID string
}

// DeleteSecretResult holds the transport-agnostic result of DeleteSecret.
type DeleteSecretResult struct {
	ARN          string
	Name         string
	DeletionDate time.Time
}

// ListSecretsResult holds the transport-agnostic result of ListSecrets.
// Stages carries the per-secret version-to-stages mapping for the returned
// page so the wire serialisers never need store access.
type ListSecretsResult struct {
	Secrets   []*secretsmanagerstore.Secret
	Stages    map[string]map[string][]string
	NextToken string
}

// DescribeSecretResult holds the transport-agnostic result of DescribeSecret.
type DescribeSecretResult struct {
	Secret             *secretsmanagerstore.Secret
	VersionIdsToStages map[string][]string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// bytesEqual reports whether two byte slices are equal, treating nil and
// empty slices as equivalent.
func bytesEqual(a, b []byte) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return bytes.Equal(a, b)
}

// resolveSecret loads the full secret (including values) by name or ARN.
func resolveSecret(store secretsmanagerstore.SecretStoreInterface, secretId string) (*secretsmanagerstore.Secret, error) {
	if strings.HasPrefix(secretId, arnPrefix) {
		secret, err := store.GetSecretByARN(secretId)
		return secret, mapStoreError(err)
	}
	secret, err := store.GetSecret(secretId)
	return secret, mapStoreError(err)
}

// resolveSecretForMetadata loads the secret metadata by name or ARN.
func resolveSecretForMetadata(store secretsmanagerstore.SecretStoreInterface, secretId string) (*secretsmanagerstore.Secret, error) {
	if strings.HasPrefix(secretId, arnPrefix) {
		name, err := store.LookupNameByARN(secretId)
		if err != nil {
			return nil, mapStoreError(err)
		}
		return store.GetSecretForMetadata(name)
	}
	secret, err := store.GetSecretForMetadata(secretId)
	return secret, mapStoreError(err)
}

// versionsToStages builds the version-id-to-stages mapping for a secret.
func versionsToStages(store secretsmanagerstore.SecretStoreInterface, secret *secretsmanagerstore.Secret) map[string][]string {
	result := make(map[string][]string)
	for _, versionId := range secret.VersionIDs {
		version, err := store.GetSecretVersion(secret.Name, versionId)
		if err == nil && len(version.VersionStages) > 0 {
			result[versionId] = version.VersionStages
		} else if versionId == secret.CurrentVersion {
			result[versionId] = []string{"AWSCURRENT"}
		}
	}
	return result
}

// createSecretCore is the single entry point for secret creation logic
// shared by the HTTP API and the admin gRPC handler. It performs all field
// validation, constructs the Secret struct, persists it, applies inline
// replication when requested, and returns the creation result.
func (s *SecretsManagerService) createSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in CreateSecretInput) (*CreateSecretResult, error) {
	if err := validateSecretName(in.Name); err != nil {
		return nil, err
	}
	secretBinary := in.SecretBinary
	if len(secretBinary) == 0 && in.SecretBinaryB64 != "" {
		decoded, err := decodeAndValidateSecretBinary(in.SecretBinaryB64)
		if err != nil {
			return nil, err
		}
		secretBinary = decoded
	}
	if in.SecretString != "" && len(secretBinary) > 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You can't specify both SecretString and SecretBinary in the same request.", http.StatusBadRequest)
	}
	if err := validateSecretStringLength(in.SecretString); err != nil {
		return nil, err
	}
	if err := validateSecretBinaryLength(secretBinary); err != nil {
		return nil, err
	}
	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}
	if err := validateKmsKeyId(in.KmsKeyId); err != nil {
		return nil, err
	}
	if err := validateClientRequestToken(in.ClientRequestToken); err != nil {
		return nil, err
	}

	// Validate tags when provided as a map (admin handler path).
	if len(in.Tags) > 0 {
		tagList := make([]types.Tag, 0, len(in.Tags))
		for k, v := range in.Tags {
			tagList = append(tagList, types.Tag{Key: k, Value: v})
		}
		if err := validateSecretTags(tagList); err != nil {
			return nil, err
		}
	}

	// ClientRequestToken idempotency: if the token is provided and a
	// version with that ID already exists for a secret with this name,
	// check whether the values match (idempotent success) or differ
	// (error). This mirrors the AWS CreateSecret idempotency rules.
	if in.ClientRequestToken != "" {
		existing, metaErr := store.GetSecretForMetadata(in.Name)
		if metaErr == nil && len(existing.VersionIDs) > 0 {
			for _, vid := range existing.VersionIDs {
				if vid == in.ClientRequestToken {
					existingVer, verErr := store.GetSecretVersion(in.Name, in.ClientRequestToken)
					if verErr != nil {
						break
					}
					if existingVer.SecretString == in.SecretString &&
						bytesEqual(existingVer.SecretBinary, secretBinary) {
						return &CreateSecretResult{
							ARN:       existing.ARN,
							Name:      existing.Name,
							VersionID: in.ClientRequestToken,
						}, nil
					}
					return nil, awserrors.NewAWSError("InvalidRequestException",
						fmt.Sprintf("You can't modify an existing secret version. The ClientRequestToken %s is already associated with a different version value.", in.ClientRequestToken),
						http.StatusBadRequest)
				}
			}
		}
	}

	secret := secretsmanagerstore.NewSecret(in.Name)
	secret.SecretString = in.SecretString
	secret.SecretBinary = secretBinary
	secret.Description = in.Description
	secret.KmsKeyId = in.KmsKeyId
	secret.Type = in.Type
	secret.InitialVersionId = in.ClientRequestToken
	if len(in.Tags) > 0 {
		secret.Tags = in.Tags
	}

	created, err := store.CreateSecret(secret)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := &CreateSecretResult{
		ARN:       created.ARN,
		Name:      created.Name,
		VersionID: created.CurrentVersion,
	}

	// AddReplicaRegions — when provided, replicate the secret to
	// the specified regions immediately after creation (same as calling
	// ReplicateSecretToRegions separately).
	if len(in.AddReplicaRegions) > 0 {
		if s.storageManager == nil {
			return nil, awserrors.NewAWSError("InvalidRequestException",
				"Replication is not configured for this service.", http.StatusBadRequest)
		}
		created, metaErr := store.GetSecretForMetadata(result.Name)
		if metaErr != nil {
			return nil, mapStoreError(metaErr)
		}
		s.replicateSecretToRegions(store, created, in.AddReplicaRegions, in.ForceOverwriteReplicaSecret, in.Region)
		if err := store.UpdateSecretMetadata(created); err != nil {
			return nil, mapStoreError(err)
		}
		result.ReplicationStatus = created.ReplicationStatus
	}

	return result, nil
}

// getSecretValueCore is the single entry point for secret value retrieval.
func (s *SecretsManagerService) getSecretValueCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in GetSecretValueInput) (*GetSecretValueResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	if in.VersionId != "" && in.VersionStage != "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You can't specify both VersionStage and VersionId.", http.StatusBadRequest)
	}

	secret, err := resolveSecret(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	var version *secretsmanagerstore.SecretVersion
	if in.VersionId == "" && in.VersionStage != "" {
		version, err = store.GetSecretVersionByStage(secret.Name, in.VersionStage)
		if err != nil {
			return nil, mapStoreError(err)
		}
	} else if in.VersionId == "" {
		if secret.CurrentVersion == "" {
			return nil, awserrors.NewAWSError("ResourceNotFoundException",
				"Secrets Manager can't find the version for the secret because no version has been created.", http.StatusNotFound)
		}
		version, err = store.GetSecretVersion(secret.Name, secret.CurrentVersion)
		if err != nil {
			return nil, mapStoreError(err)
		}
	} else {
		if isStageLabel(in.VersionId) {
			if secret.CurrentVersion == "" {
				return nil, awserrors.NewAWSError("ResourceNotFoundException",
					"Secrets Manager can't find the version for the secret because no version has been created.", http.StatusNotFound)
			}
			version, err = store.GetSecretVersionByStage(secret.Name, in.VersionId)
			if err != nil {
				return nil, mapStoreError(err)
			}
		} else {
			version, err = store.GetSecretVersion(secret.Name, in.VersionId)
			if err != nil {
				return nil, mapStoreError(err)
			}
		}
	}

	// Advance LastAccessedDate on value retrieval: DescribeSecret documents
	// the member as "the date that the secret was last accessed in the
	// Region" and omits it only when the secret was never retrieved. AWS
	// derives the date from access logs with lag; a synchronous best-effort
	// metadata update satisfies the documented semantics on-platform, and a
	// write failure must not fail the retrieval itself.
	secret.LastAccessedDate = time.Now().UTC()
	_ = store.UpdateSecretMetadata(secret)

	return &GetSecretValueResult{Secret: secret, Version: version}, nil
}

// updateSecretCore is the single entry point for secret updates.
func (s *SecretsManagerService) updateSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in UpdateSecretInput) (*UpdateSecretResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecret(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	if err := validateSecretValueMutex(in.SecretString, in.SecretBinaryB64); err != nil {
		return nil, err
	}
	if err := validateSecretStringLength(in.SecretString); err != nil {
		return nil, err
	}
	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}
	if err := validateKmsKeyId(in.KmsKeyId); err != nil {
		return nil, err
	}
	if err := validateClientRequestToken(in.ClientRequestToken); err != nil {
		return nil, err
	}

	hasSecretValue := in.SecretString != "" || in.SecretBinaryB64 != ""

	if in.SecretString != "" {
		secret.SecretString = in.SecretString
		secret.SecretBinary = nil
	} else if in.SecretBinaryB64 != "" {
		decoded, decErr := decodeAndValidateSecretBinary(in.SecretBinaryB64)
		if decErr != nil {
			return nil, decErr
		}
		secret.SecretBinary = decoded
		secret.SecretString = ""
	}
	if in.HasDescription {
		secret.Description = in.Description
	}
	if in.HasKmsKeyId {
		secret.KmsKeyId = in.KmsKeyId
	}
	if in.HasType {
		secret.Type = in.Type
	}

	var versionId string
	if hasSecretValue {
		// ClientRequestToken becomes the version ID when a new
		// version is created.
		secret.InitialVersionId = in.ClientRequestToken
		updated, err := store.UpdateSecret(secret)
		if err != nil {
			return nil, mapStoreError(err)
		}
		versionId = updated.CurrentVersion
	} else {
		// Metadata-only update: persist changes without creating a new version.
		secret.LastChangedDate = time.Now().UTC()
		if err := store.UpdateSecretMetadata(secret); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return &UpdateSecretResult{
		ARN:       secret.ARN,
		Name:      secret.Name,
		VersionID: versionId,
	}, nil
}

// deleteSecretCore is the single entry point for secret deletion logic
// shared by the HTTP API and the admin gRPC handler.
func (s *SecretsManagerService) deleteSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in DeleteSecretInput) (*DeleteSecretResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}
	if err := validateRecoveryWindow(in.RecoveryWindowInDays, in.HasRecoveryWindow, in.ForceDeleteWithoutRecovery); err != nil {
		return nil, err
	}

	// Resolve the secret by name or ARN.
	secret, err := resolveSecretForMetadata(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	// You can't delete a primary secret that is replicated to other Regions.
	if len(secret.ReplicationStatus) > 0 {
		return nil, awserrors.NewAWSError("InvalidRequestException",
			"You can't delete a primary secret that is replicated to other Regions. Remove the replicas first.", http.StatusBadRequest)
	}

	var deletionDate time.Time
	if in.ForceDeleteWithoutRecovery {
		if err := store.DeleteSecret(secret.Name); err != nil {
			return nil, mapStoreError(err)
		}
		deletionDate = time.Now().UTC()
	} else {
		rw := in.RecoveryWindowInDays
		if rw == 0 {
			rw = 30
		}
		deletionDate = time.Now().UTC().AddDate(0, 0, rw)
		if err := store.ScheduleDeletion(secret.Name, deletionDate); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return &DeleteSecretResult{
		ARN:          secret.ARN,
		Name:         secret.Name,
		DeletionDate: deletionDate,
	}, nil
}

// listSecretsCore is the single entry point for listing secrets,
// shared by the HTTP API and the admin gRPC handler. It validates the
// request, fetches all matching secrets, sorts them, and applies offset
// pagination so the caller receives one page plus the continuation token.
func (s *SecretsManagerService) listSecretsCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in ListSecretsInput) (*ListSecretsResult, error) {
	if err := validateSortBy(in.SortBy); err != nil {
		return nil, err
	}
	if err := validateSortOrder(in.SortOrder); err != nil {
		return nil, err
	}
	if err := validateMaxFilters(len(in.Filters)); err != nil {
		return nil, err
	}
	if err := validateSecretFilters(in.Filters); err != nil {
		return nil, err
	}

	// When SortBy is not specified, AWS lists secrets by CreatedDate
	// (API_ListSecrets SortBy: "If not specified, secrets are listed by
	// CreatedDate.").
	sortBy := in.SortBy
	if sortBy == "" {
		sortBy = "created-date"
	}

	maxResults, err := resolveListMaxResults(in.MaxResults, maxListSecretsResults, maxListSecretsResults)
	if err != nil {
		return nil, err
	}

	secretFilter := buildSecretFilter(in.IncludePlannedDeletion, in.Filters)

	// Fetch all items for client-side sorting.
	opts := common.ListOptions{}
	result, err := store.ListSecrets(opts, secretFilter)
	if err != nil {
		return nil, mapStoreError(err)
	}

	sortSecrets(result.Items, sortBy, in.SortOrder)

	skipCount, err := parseOffsetNextToken(in.NextToken, len(result.Items))
	if err != nil {
		return nil, err
	}
	paged := result.Items[skipCount:]
	isTruncated := maxResults < len(paged)
	if isTruncated {
		paged = paged[:maxResults]
	}

	nextToken := ""
	if isTruncated {
		nextToken = fmt.Sprintf("%d", skipCount+maxResults)
	}

	stages := make(map[string]map[string][]string, len(paged))
	for _, secret := range paged {
		stages[secret.Name] = versionsToStages(store, secret)
	}

	return &ListSecretsResult{
		Secrets:   paged,
		Stages:    stages,
		NextToken: nextToken,
	}, nil
}

// describeSecretCore is the single entry point for secret metadata
// retrieval.
func (s *SecretsManagerService) describeSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in DescribeSecretInput) (*DescribeSecretResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecretForMetadata(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	return &DescribeSecretResult{
		Secret:             secret,
		VersionIdsToStages: versionsToStages(store, secret),
	}, nil
}

// buildSecretFilter creates a store-level filter callback that combines
// IncludePlannedDeletion and the ListSecrets Filter parameter.
// Returns nil when no filtering is needed (includePlannedDeletion=true, no filters).
func buildSecretFilter(includePlannedDeletion bool, filters []SecretFilter) func(*secretsmanagerstore.Secret) bool {
	needsDeletionCheck := !includePlannedDeletion
	needsFilterCheck := len(filters) > 0
	if !needsDeletionCheck && !needsFilterCheck {
		return nil
	}

	return func(sec *secretsmanagerstore.Secret) bool {
		if needsDeletionCheck && sec.DeletedDate != nil {
			return false
		}
		if needsFilterCheck && !secretMatchesFilters(sec, filters) {
			return false
		}
		return true
	}
}

// secretMatchesFilters checks whether a single secret matches all filter
// criteria. API_Filter documents the per-key semantics: description is a
// prefix match that is not case-sensitive; name, tag-key, tag-value,
// primary-region and owning-service are case-sensitive prefix matches;
// "all" breaks each value into words and searches every attribute, not
// case-sensitive. A value prefixed with "!" negates: the secret must not
// match it. Values within one filter are OR-ed; filters are AND-ed.
func secretMatchesFilters(sec *secretsmanagerstore.Secret, filters []SecretFilter) bool {
	for _, f := range filters {
		if f.Key == "" || len(f.Values) == 0 {
			continue
		}
		if !filterMatches(sec, f.Key, f.Values) {
			return false
		}
	}
	return true
}

// filterMatches applies the positive/negation split across one filter's
// values: the secret matches when it matches at least one positive value
// (unless there are none) and matches none of the negated values.
func filterMatches(sec *secretsmanagerstore.Secret, key string, values []string) bool {
	var positives, negatives []string
	for _, v := range values {
		if strings.HasPrefix(v, "!") {
			negatives = append(negatives, strings.TrimPrefix(v, "!"))
			continue
		}
		positives = append(positives, v)
	}
	for _, v := range negatives {
		if secretMatchesFilterValue(sec, key, v) {
			return false
		}
	}
	if len(positives) == 0 {
		return true
	}
	for _, v := range positives {
		if secretMatchesFilterValue(sec, key, v) {
			return true
		}
	}
	return false
}

// secretMatchesFilterValue applies one value under one key's documented
// match kind.
func secretMatchesFilterValue(sec *secretsmanagerstore.Secret, key, value string) bool {
	switch key {
	case "name":
		return strings.HasPrefix(sec.Name, value)
	case "description":
		return strings.HasPrefix(strings.ToLower(sec.Description), strings.ToLower(value))
	case "tag-key":
		for k := range sec.Tags {
			if strings.HasPrefix(k, value) {
				return true
			}
		}
		return false
	case "tag-value":
		for _, tv := range sec.Tags {
			if strings.HasPrefix(tv, value) {
				return true
			}
		}
		return false
	case "primary-region":
		return strings.HasPrefix(sec.PrimaryRegion, value)
	case "owning-service":
		return strings.HasPrefix(sec.OwningService, value)
	case "all":
		return secretMatchesAll(sec, value)
	}
	return false
}

// secretMatchesAll implements the "all" key: "Breaks the filter value
// string into words and then searches all attributes for matches. Not
// case-sensitive." Words are split on non-alphanumeric boundaries and each
// word is matched as a case-insensitive substring against every searchable
// attribute; a secret matches when any word hits any attribute.
func secretMatchesAll(sec *secretsmanagerstore.Secret, value string) bool {
	words := strings.FieldsFunc(value, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	for _, w := range words {
		lw := strings.ToLower(w)
		if strings.Contains(strings.ToLower(sec.Name), lw) ||
			strings.Contains(strings.ToLower(sec.Description), lw) ||
			strings.Contains(strings.ToLower(sec.PrimaryRegion), lw) ||
			strings.Contains(strings.ToLower(sec.OwningService), lw) {
			return true
		}
		for k, v := range sec.Tags {
			if strings.Contains(strings.ToLower(k), lw) || strings.Contains(strings.ToLower(v), lw) {
				return true
			}
		}
	}
	return false
}

func sortSecrets(secrets []*secretsmanagerstore.Secret, sortBy, sortOrder string) {
	desc := sortOrder == "desc"
	switch sortBy {
	case "name":
		sort.Slice(secrets, func(i, j int) bool {
			if desc {
				return secrets[i].Name > secrets[j].Name
			}
			return secrets[i].Name < secrets[j].Name
		})
	case "created-date":
		sort.Slice(secrets, func(i, j int) bool {
			if desc {
				return secrets[i].CreatedDate.After(secrets[j].CreatedDate)
			}
			return secrets[i].CreatedDate.Before(secrets[j].CreatedDate)
		})
	case "last-accessed-date":
		sort.Slice(secrets, func(i, j int) bool {
			if desc {
				return secrets[i].LastAccessedDate.After(secrets[j].LastAccessedDate)
			}
			return secrets[i].LastAccessedDate.Before(secrets[j].LastAccessedDate)
		})
	case "last-changed-date":
		sort.Slice(secrets, func(i, j int) bool {
			if desc {
				return secrets[i].LastChangedDate.After(secrets[j].LastChangedDate)
			}
			return secrets[i].LastChangedDate.Before(secrets[j].LastChangedDate)
		})
	}
}
