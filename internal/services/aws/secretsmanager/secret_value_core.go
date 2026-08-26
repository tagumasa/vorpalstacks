package secretsmanager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// ---------------------------------------------------------------------------
// Input structs — transport-agnostic DTOs shared by HTTP API and admin handler
// ---------------------------------------------------------------------------

// PutSecretValueInput carries all fields needed for PutSecretValue.
type PutSecretValueInput struct {
	SecretId           string
	SecretString       string
	SecretBinaryB64    string
	ClientRequestToken string
	VersionStages      []string
	RotationToken      string
}

// ListSecretVersionIdsInput carries all fields needed for ListSecretVersionIds.
type ListSecretVersionIdsInput struct {
	SecretId          string
	MaxResults        *int
	NextToken         string
	IncludeDeprecated bool
}

// UpdateSecretVersionStageInput carries all fields needed for
// UpdateSecretVersionStage.
type UpdateSecretVersionStageInput struct {
	SecretId            string
	VersionStage        string
	MoveToVersionId     string
	RemoveFromVersionId string
}

// BatchGetSecretValueInput carries all fields needed for BatchGetSecretValue.
type BatchGetSecretValueInput struct {
	SecretIdList []string
	Filters      []SecretFilter
	MaxResults   *int
	NextToken    string
}

// ---------------------------------------------------------------------------
// Result structs — transport-agnostic results
// ---------------------------------------------------------------------------

// PutSecretValueResult holds the transport-agnostic result of PutSecretValue.
type PutSecretValueResult struct {
	ARN           string
	Name          string
	VersionID     string
	VersionStages []string
}

// ListSecretVersionIdsResult holds the transport-agnostic result of
// ListSecretVersionIds.
type ListSecretVersionIdsResult struct {
	ARN       string
	Name      string
	Versions  []secretsmanagerstore.SecretVersion
	NextToken string
}

// UpdateSecretVersionStageResult holds the transport-agnostic result of
// UpdateSecretVersionStage.
type UpdateSecretVersionStageResult struct {
	ARN  string
	Name string
}

// BatchSecretValueEntry pairs a secret with the retrieved version for one
// successful BatchGetSecretValue entry.
type BatchSecretValueEntry struct {
	Secret  *secretsmanagerstore.Secret
	Version *secretsmanagerstore.SecretVersion
}

// BatchSecretValueError describes one failed BatchGetSecretValue entry.
type BatchSecretValueError struct {
	SecretId  string
	ErrorCode string
	Message   string
}

// BatchGetSecretValueResult holds the transport-agnostic result of
// BatchGetSecretValue.
type BatchGetSecretValueResult struct {
	SecretValues []BatchSecretValueEntry
	Errors       []BatchSecretValueError
	NextToken    string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// putSecretValueCore is the single entry point for storing a secret value.
func (s *SecretsManagerService) putSecretValueCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in PutSecretValueInput) (*PutSecretValueResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecret(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	// RotationToken is an optional parameter used only for cross-account
	// rotation (IAM assumed role). It is validated against the Smithy
	// RotationTokenType constraints (@length 36-256, @pattern) when present
	// but is not required for staging label assignment. AWS does not
	// enforce AWSPENDING at the API level — the standard same-account
	// rotation Lambda template calls PutSecretValue with
	// VersionStages=["AWSPENDING"] and no RotationToken.
	if err := validateRotationToken(in.RotationToken); err != nil {
		return nil, err
	}

	if err := validateSecretValueMutex(in.SecretString, in.SecretBinaryB64); err != nil {
		return nil, err
	}
	if in.SecretString == "" && in.SecretBinaryB64 == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You must include either a SecretString or SecretBinary parameter.", http.StatusBadRequest)
	}
	if err := validateSecretStringLength(in.SecretString); err != nil {
		return nil, err
	}
	if err := validateClientRequestToken(in.ClientRequestToken); err != nil {
		return nil, err
	}

	secretBinary, err := decodeAndValidateSecretBinary(in.SecretBinaryB64)
	if err != nil {
		return nil, err
	}

	// Determine final staging labels.
	finalStages := in.VersionStages
	if len(finalStages) == 0 {
		finalStages = []string{"AWSCURRENT"}
	}

	// Check whether AWSCURRENT is explicitly requested.
	hasAWSCURRENT := false
	for _, st := range finalStages {
		if st == "AWSCURRENT" {
			hasAWSCURRENT = true
			break
		}
	}

	// ClientRequestToken idempotency: if the token is provided and a
	// version with that ID already exists, check whether the values match
	// (idempotent success) or differ (error). This mirrors the AWS
	// PutSecretValue idempotency rules.
	if in.ClientRequestToken != "" {
		existingVer, verErr := store.GetSecretVersion(secret.Name, in.ClientRequestToken)
		if verErr == nil {
			if existingVer.SecretString == in.SecretString &&
				bytesEqual(existingVer.SecretBinary, secretBinary) {
				return &PutSecretValueResult{
					ARN:           secret.ARN,
					Name:          secret.Name,
					VersionID:     in.ClientRequestToken,
					VersionStages: existingVer.VersionStages,
				}, nil
			}
			return nil, awserrors.NewAWSError("InvalidRequestException",
				fmt.Sprintf("You can't modify an existing secret version. The ClientRequestToken %s is already associated with a different version value.", in.ClientRequestToken),
				http.StatusBadRequest)
		}
	}

	var resultVersionId string
	var resultStages []string

	if hasAWSCURRENT {
		// AWSCURRENT is in the stages — store.UpdateSecret handles version
		// creation, AWSCURRENT assignment, and old-version demotion.
		secret.SecretString = in.SecretString
		secret.SecretBinary = secretBinary
		secret.InitialVersionId = in.ClientRequestToken
		updated, err := store.UpdateSecret(secret)
		if err != nil {
			return nil, mapStoreError(err)
		}

		resultVersionId = updated.CurrentVersion

		// If the user specified additional stages alongside AWSCURRENT, apply them.
		if len(in.VersionStages) > 0 {
			if err := store.UpdateSecretVersionStage(updated.Name, updated.CurrentVersion, finalStages); err != nil {
				return nil, mapStoreError(err)
			}
			resultStages = finalStages
		} else {
			version, err := store.GetSecretVersion(updated.Name, updated.CurrentVersion)
			if err != nil {
				return nil, mapStoreError(err)
			}
			resultStages = version.VersionStages
		}
	} else {
		// AWSCURRENT is NOT in the stages — create a version without moving
		// AWSCURRENT from the existing current version.
		newVersion := secretsmanagerstore.NewSecretVersion(in.ClientRequestToken)
		newVersion.SecretName = secret.Name
		newVersion.SecretString = in.SecretString
		newVersion.SecretBinary = secretBinary
		newVersion.VersionStages = finalStages

		if err := store.CreateVersionDirect(secret.Name, newVersion); err != nil {
			return nil, mapStoreError(err)
		}

		secret.VersionIDs = append(secret.VersionIDs, newVersion.VersionId)
		secret.LastChangedDate = time.Now().UTC()
		if err := store.UpdateSecretMetadata(secret); err != nil {
			return nil, mapStoreError(err)
		}

		resultVersionId = newVersion.VersionId
		resultStages = finalStages
	}

	return &PutSecretValueResult{
		ARN:           secret.ARN,
		Name:          secret.Name,
		VersionID:     resultVersionId,
		VersionStages: resultStages,
	}, nil
}

// listSecretVersionIdsCore is the single entry point for listing the
// versions of a secret. Supports MaxResults (1-100), NextToken pagination,
// and IncludeDeprecated filtering.
func (s *SecretsManagerService) listSecretVersionIdsCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in ListSecretVersionIdsInput) (*ListSecretVersionIdsResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecretForMetadata(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	versions, err := store.ListSecretVersions(secret.Name)
	if err != nil {
		return nil, mapStoreError(err)
	}

	// IncludeDeprecated defaults to false: versions without staging
	// labels are considered deprecated and excluded unless explicitly
	// requested.
	if !in.IncludeDeprecated {
		filtered := make([]secretsmanagerstore.SecretVersion, 0, len(versions))
		for _, v := range versions {
			if len(v.VersionStages) > 0 {
				filtered = append(filtered, v)
			}
		}
		versions = filtered
	}

	// Pagination: offset-based NextToken matching the ListSecrets sorted-
	// result pattern.  Smithy MaxResultsType range is 1-100; out-of-range
	// and explicit-zero values are rejected.
	maxResults, err := resolveListMaxResults(in.MaxResults, maxListSecretsResults, maxListSecretsResults)
	if err != nil {
		return nil, err
	}

	skipCount, err := parseOffsetNextToken(in.NextToken, len(versions))
	if err != nil {
		return nil, err
	}

	end := skipCount + maxResults
	if end > len(versions) {
		end = len(versions)
	}
	paged := versions[skipCount:end]

	result := &ListSecretVersionIdsResult{
		ARN:      secret.ARN,
		Name:     secret.Name,
		Versions: paged,
	}
	if end < len(versions) {
		result.NextToken = fmt.Sprintf("%d", end)
	}
	return result, nil
}

// updateSecretVersionStageCore is the single entry point for modifying the
// staging labels attached to a version. Supports three modes per AWS
// documentation:
//  1. Add a label to a version (MoveToVersionId only).
//  2. Remove a label from a version (RemoveFromVersionId only).
//  3. Move a label between versions (both MoveToVersionId and RemoveFromVersionId).
func (s *SecretsManagerService) updateSecretVersionStageCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in UpdateSecretVersionStageInput) (*UpdateSecretVersionStageResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}
	if in.VersionStage == "" {
		return nil, awserrors.ErrMissingParameter
	}

	if in.MoveToVersionId == "" && in.RemoveFromVersionId == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You must specify either MoveToVersionId or RemoveFromVersionId.", http.StatusBadRequest)
	}

	// MoveToVersionId == RemoveFromVersionId is a no-op "move" that
	// AWS rejects with InvalidParameterException.
	if in.MoveToVersionId != "" && in.RemoveFromVersionId != "" && in.MoveToVersionId == in.RemoveFromVersionId {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"MoveToVersionId and RemoveFromVersionId must not be the same version.", http.StatusBadRequest)
	}

	secret, err := resolveSecret(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	// Mode 2: Remove-only — remove the staging label from the specified version.
	if in.MoveToVersionId == "" {
		// AWSCURRENT can only be moved (RemoveFromVersionId + MoveToVersionId),
		// never removed in isolation. Every secret must always have exactly one
		// AWSCURRENT version.
		if in.VersionStage == "AWSCURRENT" {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("To remove %s from a version, you must also specify the version to move it to.", in.VersionStage),
				http.StatusBadRequest)
		}
		version, err := store.GetSecretVersion(secret.Name, in.RemoveFromVersionId)
		if err != nil {
			return nil, mapStoreError(err)
		}
		newStages := make([]string, 0, len(version.VersionStages))
		for _, st := range version.VersionStages {
			if st != in.VersionStage {
				newStages = append(newStages, st)
			}
		}
		if len(newStages) == len(version.VersionStages) {
			return nil, awserrors.NewAWSError("InvalidRequestException",
				fmt.Sprintf("Version %s does not have staging label %s", in.RemoveFromVersionId, in.VersionStage), http.StatusBadRequest)
		}
		if err := store.UpdateSecretVersionStage(secret.Name, in.RemoveFromVersionId, newStages); err != nil {
			return nil, mapStoreError(err)
		}
		return &UpdateSecretVersionStageResult{
			ARN:  secret.ARN,
			Name: secret.Name,
		}, nil
	}

	// Mode 1 or 3: Add or move. AWS API_UpdateSecretVersionStage: a label
	// not attached anywhere is added by MoveToVersionId alone (Example 1);
	// a label already attached to a different version requires
	// RemoveFromVersionId to match the holding version or the operation
	// fails. The documentation does not name the failure error, so the
	// omission case uses the operation's parameter-contract error and the
	// mismatch case reuses the remove-only path's state-mismatch error.
	holderVer, holderErr := store.GetSecretVersionByStage(secret.Name, in.VersionStage)
	if holderErr != nil && !errors.Is(holderErr, secretsmanagerstore.ErrInvalidVersionId) {
		return nil, mapStoreError(holderErr)
	}
	if holderErr == nil {
		// The label is currently attached to holderVer.
		if in.RemoveFromVersionId != "" && in.RemoveFromVersionId != holderVer.VersionId {
			return nil, awserrors.NewAWSError("InvalidRequestException",
				fmt.Sprintf("Version %s does not have staging label %s", in.RemoveFromVersionId, in.VersionStage), http.StatusBadRequest)
		}
		if holderVer.VersionId != in.MoveToVersionId && in.RemoveFromVersionId == "" {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("Staging label %s is already attached to a different version; you must also specify RemoveFromVersionId.", in.VersionStage), http.StatusBadRequest)
		}
		// A no-op re-add (holder == MoveToVersionId) and a real move are
		// both safe here; equal IDs skip the removal inside MoveStage.
		if err := store.MoveStage(secret, in.VersionStage, in.MoveToVersionId, holderVer.VersionId); err != nil {
			return nil, mapStoreError(err)
		}
	} else {
		// The label is not attached to any version — a pure add.
		if in.RemoveFromVersionId != "" {
			return nil, awserrors.NewAWSError("InvalidRequestException",
				fmt.Sprintf("Version %s does not have staging label %s", in.RemoveFromVersionId, in.VersionStage), http.StatusBadRequest)
		}
		if err := store.MoveStage(secret, in.VersionStage, in.MoveToVersionId, ""); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return &UpdateSecretVersionStageResult{
		ARN:  secret.ARN,
		Name: secret.Name,
	}, nil
}

// batchGetSecretValueCore is the single entry point for retrieving multiple
// secret values. You must include either SecretIdList or Filters, but not
// both. MaxResults (1-20) and NextToken pagination are supported when using
// Filters.
func (s *SecretsManagerService) batchGetSecretValueCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in BatchGetSecretValueInput) (*BatchGetSecretValueResult, error) {
	// Mutual exclusion: exactly one of SecretIdList or Filters must be
	// provided.  AWS: "You must include Filters or SecretIdList,
	// but not both."
	if len(in.SecretIdList) == 0 && len(in.Filters) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You must include either Filters or SecretIdList, but not both.", http.StatusBadRequest)
	}
	if len(in.SecretIdList) > 0 && len(in.Filters) > 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You can't specify both Filters and SecretIdList in the same request.", http.StatusBadRequest)
	}
	// API_BatchGetSecretValue MaxResults: "To use this parameter, you
	// must also use the Filters parameter." SecretIdList mode is capped
	// at 20 entries and never paginates, so the pairing is rejected
	// outright rather than silently ignored.
	if len(in.SecretIdList) > 0 && in.MaxResults != nil {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"MaxResults can only be used together with the Filters parameter.", http.StatusBadRequest)
	}

	if err := validateSecretIdList(in.SecretIdList); err != nil {
		return nil, err
	}
	if err := validateMaxFilters(len(in.Filters)); err != nil {
		return nil, err
	}
	if err := validateSecretFilters(in.Filters); err != nil {
		return nil, err
	}

	// Smithy MaxResultsBatchType range is 1-20; out-of-range and
	// explicit-zero values are rejected whichever selection mode is used.
	maxResults, err := resolveListMaxResults(in.MaxResults, maxBatchSecretsResults, maxBatchSecretsResults)
	if err != nil {
		return nil, err
	}

	// Build the list of secret IDs to retrieve.
	//
	// When using Filters, we list all matching secrets first, apply
	// pagination, then retrieve values for the current page.
	var targetIds []string
	nextToken := ""

	if len(in.Filters) > 0 {
		secretFilter := buildSecretFilter(true, in.Filters)
		result, err := store.ListSecrets(common.ListOptions{}, secretFilter)
		if err != nil {
			return nil, err
		}

		skipCount, err := parseOffsetNextToken(in.NextToken, len(result.Items))
		if err != nil {
			return nil, err
		}

		end := skipCount + maxResults
		if end > len(result.Items) {
			end = len(result.Items)
		}
		for _, sec := range result.Items[skipCount:end] {
			targetIds = append(targetIds, sec.Name)
		}
		if end < len(result.Items) {
			nextToken = fmt.Sprintf("%d", end)
		}
	} else {
		targetIds = in.SecretIdList
	}

	result := &BatchGetSecretValueResult{}

	for _, secretId := range targetIds {
		secret, err := resolveSecret(store, secretId)
		if err != nil {
			result.Errors = append(result.Errors, BatchSecretValueError{
				SecretId:  secretId,
				ErrorCode: "ResourceNotFoundException",
				Message:   "Secrets Manager can't find the specified secret.",
			})
			continue
		}

		version, err := store.GetSecretVersion(secret.Name, secret.CurrentVersion)
		if err != nil {
			result.Errors = append(result.Errors, BatchSecretValueError{
				SecretId:  secretId,
				ErrorCode: "ResourceNotFoundException",
				Message:   fmt.Sprintf("Secrets Manager can't find the version for secret %s", secret.Name),
			})
			continue
		}

		// Advance LastAccessedDate for the same documented reason as
		// GetSecretValue; best-effort, never fails the retrieval.
		secret.LastAccessedDate = time.Now().UTC()
		_ = store.UpdateSecretMetadata(secret)

		result.SecretValues = append(result.SecretValues, BatchSecretValueEntry{
			Secret:  secret,
			Version: version,
		})
	}

	result.NextToken = nextToken
	return result, nil
}
