package secretsmanager

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// PutSecretValue stores a secret value in a secret.
func (s *SecretsManagerService) PutSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	secretString := request.GetStringParam(req.Parameters, "SecretString")
	secretBinaryStr := request.GetStringParam(req.Parameters, "SecretBinary")
	clientRequestToken := request.GetStringParam(req.Parameters, "ClientRequestToken")
	versionStages := request.GetStringList(req.Parameters, "VersionStages")

	// RotationToken is an optional parameter used only for cross-account
	// rotation (IAM assumed role). It is validated against the Smithy
	// RotationTokenType constraints (@length 36-256, @pattern) when present
	// but is not required for staging label assignment. AWS does not
	// enforce AWSPENDING at the API level — the standard same-account
	// rotation Lambda template calls PutSecretValue with
	// VersionStages=["AWSPENDING"] and no RotationToken.
	rotationToken := request.GetStringParam(req.Parameters, "RotationToken")
	if err := validateRotationToken(rotationToken); err != nil {
		return nil, err
	}

	if err := validateSecretValueMutex(secretString, secretBinaryStr); err != nil {
		return nil, err
	}
	if secretString == "" && secretBinaryStr == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You must include either a SecretString or SecretBinary parameter.", http.StatusBadRequest)
	}
	if err := validateSecretStringLength(secretString); err != nil {
		return nil, err
	}
	if err := validateClientRequestToken(clientRequestToken); err != nil {
		return nil, err
	}

	secretBinary, err := decodeAndValidateSecretBinary(secretBinaryStr)
	if err != nil {
		return nil, err
	}

	// Determine final staging labels.
	finalStages := versionStages
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// ClientRequestToken idempotency: if the token is provided and a
	// version with that ID already exists, check whether the values match
	// (idempotent success) or differ (error). This mirrors the AWS
	// PutSecretValue idempotency rules.
	if clientRequestToken != "" {
		existingVer, verErr := store.GetSecretVersion(secret.Name, clientRequestToken)
		if verErr == nil {
			if existingVer.SecretString == secretString &&
				bytesEqual(existingVer.SecretBinary, secretBinary) {
				return map[string]interface{}{
					"ARN":           secret.ARN,
					"Name":          secret.Name,
					"VersionId":     clientRequestToken,
					"VersionStages": existingVer.VersionStages,
				}, nil
			}
			return nil, awserrors.NewAWSError("InvalidRequestException",
				fmt.Sprintf("You can't modify an existing secret version. The ClientRequestToken %s is already associated with a different version value.", clientRequestToken),
				http.StatusBadRequest)
		}
	}

	var resultVersionId string
	var resultStages []string

	if hasAWSCURRENT {
		// AWSCURRENT is in the stages — store.UpdateSecret handles version
		// creation, AWSCURRENT assignment, and old-version demotion.
		secret.SecretString = secretString
		secret.SecretBinary = secretBinary
		secret.InitialVersionId = clientRequestToken
		updated, err := store.UpdateSecret(secret)
		if err != nil {
			return nil, mapStoreError(err)
		}

		resultVersionId = updated.CurrentVersion

		// If the user specified additional stages alongside AWSCURRENT, apply them.
		if len(versionStages) > 0 {
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
		newVersion := secretsmanagerstore.NewSecretVersion(clientRequestToken)
		newVersion.SecretName = secret.Name
		newVersion.SecretString = secretString
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

	return map[string]interface{}{
		"ARN":           secret.ARN,
		"Name":          secret.Name,
		"VersionId":     resultVersionId,
		"VersionStages": resultStages,
	}, nil
}

// ListSecrets lists the secrets in AWS Secrets Manager.
func (s *SecretsManagerService) ListSecrets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")
	if maxResults > 100 {
		maxResults = 100
	}
	includePlannedDeletion := request.GetBoolParam(req.Parameters, "IncludePlannedDeletion")
	sortBy := request.GetStringParam(req.Parameters, "SortBy")
	sortOrder := request.GetStringParam(req.Parameters, "SortOrder")
	if err := validateSortBy(sortBy); err != nil {
		return nil, err
	}
	if err := validateSortOrder(sortOrder); err != nil {
		return nil, err
	}
	filters := request.GetListParam(req.Parameters, "Filters")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// When SortBy is not specified, AWS defaults to sorting by name
	// in ascending order.
	if sortBy == "" {
		sortBy = "name"
	}

	var secrets []*secretsmanagerstore.Secret
	var nextMarker string
	var isTruncated bool

	if sortBy != "" {
		// Fetch all items for client-side sorting.
		coreResult, err := s.listSecretsCore(ctx, store, ListSecretsInput{
			MaxResults:             -1,
			Filters:                filters,
			IncludePlannedDeletion: includePlannedDeletion,
		})
		if err != nil {
			return nil, err
		}
		sortSecrets(coreResult.Secrets, sortBy, sortOrder)

		skipCount := 0
		if nextToken != "" {
			if n, e := fmt.Sscanf(nextToken, "%d", &skipCount); n != 1 || e != nil {
				skipCount = 0
			}
		}
		paged := coreResult.Secrets[skipCount:]
		isTruncated = maxResults < len(paged)
		if isTruncated {
			paged = paged[:maxResults]
		}
		secrets = paged
		if isTruncated {
			nextMarker = fmt.Sprintf("%d", skipCount+maxResults)
		}
	} else {
		coreResult, err := s.listSecretsCore(ctx, store, ListSecretsInput{
			MaxResults:             maxResults,
			NextToken:              nextToken,
			Filters:                filters,
			IncludePlannedDeletion: includePlannedDeletion,
		})
		if err != nil {
			return nil, err
		}
		secrets = coreResult.Secrets
		nextMarker = coreResult.NextToken
		isTruncated = nextMarker != ""
	}

	secretList := make([]interface{}, 0, len(secrets))
	for _, secret := range secrets {
		entry := map[string]interface{}{
			"ARN":                    secret.ARN,
			"Name":                   secret.Name,
			"CreatedDate":            secret.CreatedDate.Unix(),
			"LastChangedDate":        secret.LastChangedDate.Unix(),
			"SecretVersionsToStages": s.buildSecretVersionsToStages(reqCtx, secret),
		}
		if secret.Description != "" {
			entry["Description"] = secret.Description
		}
		if secret.KmsKeyId != "" {
			entry["KmsKeyId"] = secret.KmsKeyId
		}
		if !secret.LastAccessedDate.IsZero() {
			entry["LastAccessedDate"] = secret.LastAccessedDate.Unix()
		}
		if secret.Type != "" {
			entry["Type"] = secret.Type
		}
		s.addRotationFields(entry, secret)
		if len(secret.Tags) > 0 {
			entry["Tags"] = s.buildTagsList(secret)
		}
		if len(secret.ReplicationStatus) > 0 {
			entry["ReplicationStatus"] = buildReplicationStatusResponse(secret.ReplicationStatus)
		}
		secretList = append(secretList, entry)
	}

	response := map[string]interface{}{
		"SecretList": secretList,
	}
	if isTruncated {
		pagination.SetNextToken(response, "NextToken", nextMarker)
	}
	return response, nil
}

// buildSecretFilter creates a store-level filter callback that combines
// IncludePlannedDeletion and the ListSecrets Filter parameter.
// Returns nil when no filtering is needed (includePlannedDeletion=true, no filters).
func buildSecretFilter(includePlannedDeletion bool, filters []map[string]interface{}) func(*secretsmanagerstore.Secret) bool {
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

// secretMatchesFilters checks whether a single secret matches all filter criteria.
func secretMatchesFilters(sec *secretsmanagerstore.Secret, filters []map[string]interface{}) bool {
	for _, f := range filters {
		key := request.GetStringParam(f, "Key")
		values := request.GetStringList(f, "Values")
		if key == "" || len(values) == 0 {
			continue
		}
		switch key {
		case "name":
			if !matchesAny(sec.Name, values, strings.Contains) {
				return false
			}
		case "description":
			if !matchesAny(sec.Description, values, strings.Contains) {
				return false
			}
		case "tag-key":
			found := false
			for _, v := range values {
				if _, ok := sec.Tags[v]; ok {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case "tag-value":
			found := false
			for _, v := range values {
				for _, tv := range sec.Tags {
					if strings.Contains(tv, v) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				return false
			}
		case "primary-region":
			if !matchesAny(sec.PrimaryRegion, values, func(a, b string) bool { return a == b }) {
				return false
			}
		case "owning-service":
			if !matchesAny(sec.OwningService, values, func(a, b string) bool { return a == b }) {
				return false
			}
		}
	}
	return true
}

func matchesAny(s string, values []string, cmp func(string, string) bool) bool {
	for _, v := range values {
		if cmp(s, v) {
			return true
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

// DescribeSecret returns the metadata for a secret.
func (s *SecretsManagerService) DescribeSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"ARN":                secret.ARN,
		"Name":               secret.Name,
		"CreatedDate":        secret.CreatedDate.Unix(),
		"LastChangedDate":    secret.LastChangedDate.Unix(),
		"VersionIdsToStages": s.buildSecretVersionsToStages(reqCtx, secret),
	}
	if secret.Description != "" {
		result["Description"] = secret.Description
	}
	if secret.KmsKeyId != "" {
		result["KmsKeyId"] = secret.KmsKeyId
	}
	if !secret.LastAccessedDate.IsZero() {
		result["LastAccessedDate"] = secret.LastAccessedDate.Unix()
	}
	tags := s.buildTagsList(secret)
	if len(tags) > 0 {
		result["Tags"] = tags
	}
	if secret.DeletedDate != nil {
		result["DeletedDate"] = secret.DeletedDate.Unix()
	}
	s.addRotationFields(result, secret)
	if secret.OwningService != "" {
		result["OwningService"] = secret.OwningService
	}
	if secret.PrimaryRegion != "" {
		result["PrimaryRegion"] = secret.PrimaryRegion
	}
	if secret.Type != "" {
		result["Type"] = secret.Type
	}
	if !secret.NextRotationDate.IsZero() {
		result["NextRotationDate"] = secret.NextRotationDate.Unix()
	}
	if len(secret.ReplicationStatus) > 0 {
		result["ReplicationStatus"] = buildReplicationStatusResponse(secret.ReplicationStatus)
	}
	return result, nil
}

// ListSecretVersionIds lists the versions of a secret.
// Supports MaxResults (1-100), NextToken pagination, and IncludeDeprecated
// filtering.
func (s *SecretsManagerService) ListSecretVersionIds(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
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
	includeDeprecated := request.GetBoolParam(req.Parameters, "IncludeDeprecated")
	if !includeDeprecated {
		filtered := make([]secretsmanagerstore.SecretVersion, 0, len(versions))
		for _, v := range versions {
			if len(v.VersionStages) > 0 {
				filtered = append(filtered, v)
			}
		}
		versions = filtered
	}

	// Pagination: offset-based NextToken matching the ListSecrets sorted-
	// result pattern.  Smithy MaxResultsType range is 1-100; clamp
	// explicitly because GetMaxItems allows up to AbsoluteMaxItems (1000).
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")
	if maxResults > 100 {
		maxResults = 100
	}
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	skipCount := 0
	if nextToken != "" {
		if n, e := fmt.Sscanf(nextToken, "%d", &skipCount); n != 1 || e != nil || skipCount < 0 || skipCount >= len(versions) {
			return nil, awserrors.NewAWSError("InvalidNextTokenException",
				"Your request has an invalid next token.", http.StatusBadRequest)
		}
	}

	end := skipCount + maxResults
	if end > len(versions) {
		end = len(versions)
	}
	paged := versions[skipCount:end]

	versionList := make([]interface{}, 0, len(paged))
	for _, version := range paged {
		entry := map[string]interface{}{
			"VersionId":     version.VersionId,
			"VersionStages": version.VersionStages,
			"CreatedDate":   version.CreatedDate.Unix(),
		}
		if !version.LastAccessedDate.IsZero() {
			entry["LastAccessedDate"] = version.LastAccessedDate.Unix()
		}
		if len(version.KmsKeyIds) > 0 {
			entry["KmsKeyIds"] = version.KmsKeyIds
		}
		versionList = append(versionList, entry)
	}

	result := map[string]interface{}{
		"ARN":      secret.ARN,
		"Name":     secret.Name,
		"Versions": versionList,
	}
	if end < len(versions) {
		pagination.SetNextToken(result, "NextToken", fmt.Sprintf("%d", end))
	}
	return result, nil
}

// UpdateSecretVersionStage modifies the staging labels attached to a version.
// Supports three modes per AWS documentation:
//  1. Add a label to a version (MoveToVersionId only).
//  2. Remove a label from a version (RemoveFromVersionId only).
//  3. Move a label between versions (both MoveToVersionId and RemoveFromVersionId).
func (s *SecretsManagerService) UpdateSecretVersionStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}
	versionStage := request.GetStringParam(req.Parameters, "VersionStage")
	if versionStage == "" {
		return nil, awserrors.ErrMissingParameter
	}
	moveToVersionId := request.GetStringParam(req.Parameters, "MoveToVersionId")
	removeFromVersionId := request.GetStringParam(req.Parameters, "RemoveFromVersionId")

	if moveToVersionId == "" && removeFromVersionId == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You must specify either MoveToVersionId or RemoveFromVersionId.", http.StatusBadRequest)
	}

	// MoveToVersionId == RemoveFromVersionId is a no-op "move" that
	// AWS rejects with InvalidParameterException.
	if moveToVersionId != "" && removeFromVersionId != "" && moveToVersionId == removeFromVersionId {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"MoveToVersionId and RemoveFromVersionId must not be the same version.", http.StatusBadRequest)
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Mode 2: Remove-only — remove the staging label from the specified version.
	if moveToVersionId == "" {
		// AWSCURRENT can only be moved (RemoveFromVersionId + MoveToVersionId),
		// never removed in isolation. Every secret must always have exactly one
		// AWSCURRENT version.
		if versionStage == "AWSCURRENT" {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("To remove %s from a version, you must also specify the version to move it to.", versionStage),
				http.StatusBadRequest)
		}
		version, err := store.GetSecretVersion(secret.Name, removeFromVersionId)
		if err != nil {
			return nil, mapStoreError(err)
		}
		newStages := make([]string, 0, len(version.VersionStages))
		for _, st := range version.VersionStages {
			if st != versionStage {
				newStages = append(newStages, st)
			}
		}
		if len(newStages) == len(version.VersionStages) {
			return nil, awserrors.NewAWSError("InvalidRequestException",
				fmt.Sprintf("Version %s does not have staging label %s", removeFromVersionId, versionStage), http.StatusBadRequest)
		}
		if err := store.UpdateSecretVersionStage(secret.Name, removeFromVersionId, newStages); err != nil {
			return nil, mapStoreError(err)
		}
		return map[string]interface{}{
			"ARN":  secret.ARN,
			"Name": secret.Name,
		}, nil
	}

	// Mode 1 or 3: Add or move — auto-resolve RemoveFromVersionId if not provided.
	removeVersionId := removeFromVersionId
	if removeVersionId == "" {
		existingVer, err := store.GetSecretVersionByStage(secret.Name, versionStage)
		if err != nil {
			return nil, awserrors.NewAWSError("ResourceNotFoundException",
				fmt.Sprintf("Secrets Manager can't find the version with stage %s", versionStage), http.StatusNotFound)
		}
		removeVersionId = existingVer.VersionId
	}

	targetVerId := moveToVersionId
	if err := store.MoveStage(secret, versionStage, targetVerId, removeVersionId); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}

// BatchGetSecretValue retrieves multiple secret values.
// You must include either SecretIdList or Filters, but not both.
// MaxResults (1-20) and NextToken pagination are supported when using
// Filters.
func (s *SecretsManagerService) BatchGetSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretIdList := request.GetStringList(req.Parameters, "SecretIdList")
	filters := request.GetListParam(req.Parameters, "Filters")

	// Mutual exclusion: exactly one of SecretIdList or Filters must be
	// provided.  AWS: "You must include Filters or SecretIdList,
	// but not both."
	if len(secretIdList) == 0 && len(filters) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You must include either Filters or SecretIdList, but not both.", http.StatusBadRequest)
	}
	if len(secretIdList) > 0 && len(filters) > 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You can't specify both Filters and SecretIdList in the same request.", http.StatusBadRequest)
	}

	if err := validateSecretIdList(secretIdList); err != nil {
		return nil, err
	}
	if err := validateMaxFilters(len(filters)); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Build the list of secret IDs to retrieve.
	//
	// When using Filters, we list all matching secrets first, apply
	// pagination, then retrieve values for the current page.
	var targetIds []string
	var nextToken string
	var isTruncated bool

	if len(filters) > 0 {
		secretFilter := buildSecretFilter(true, filters)
		result, err := store.ListSecrets(common.ListOptions{}, secretFilter)
		if err != nil {
			return nil, err
		}

		// Smithy MaxResultsBatchType range is 1-20; clamp explicitly.
		maxResults := pagination.GetMaxItems(req.Parameters, 20, "MaxResults")
		if maxResults > 20 {
			maxResults = 20
		}
		token := pagination.GetMarker(req.Parameters, "NextToken")

		skipCount := 0
		if token != "" {
			if n, e := fmt.Sscanf(token, "%d", &skipCount); n != 1 || e != nil || skipCount < 0 || skipCount >= len(result.Items) {
				return nil, awserrors.NewAWSError("InvalidNextTokenException",
					"Your request has an invalid next token.", http.StatusBadRequest)
			}
		}

		end := skipCount + maxResults
		if end > len(result.Items) {
			end = len(result.Items)
		}
		for _, sec := range result.Items[skipCount:end] {
			targetIds = append(targetIds, sec.Name)
		}
		if end < len(result.Items) {
			isTruncated = true
			nextToken = fmt.Sprintf("%d", end)
		}
	} else {
		targetIds = secretIdList
	}

	secretValues := []interface{}{}
	apiErrors := []interface{}{}

	for _, secretId := range targetIds {
		secret, err := s.resolveSecret(reqCtx, secretId)
		if err != nil {
			apiErrors = append(apiErrors, map[string]interface{}{
				"SecretId":  secretId,
				"ErrorCode": "ResourceNotFoundException",
				"Message":   "Secrets Manager can't find the specified secret.",
			})
			continue
		}

		version, err := store.GetSecretVersion(secret.Name, secret.CurrentVersion)
		if err != nil {
			apiErrors = append(apiErrors, map[string]interface{}{
				"SecretId":  secretId,
				"ErrorCode": "ResourceNotFoundException",
				"Message":   fmt.Sprintf("Secrets Manager can't find the version for secret %s", secret.Name),
			})
			continue
		}

		entry := map[string]interface{}{
			"ARN":           secret.ARN,
			"Name":          secret.Name,
			"VersionId":     version.VersionId,
			"VersionStages": version.VersionStages,
			"CreatedDate":   version.CreatedDate.Unix(),
		}

		if version.SecretString != "" {
			entry["SecretString"] = version.SecretString
		}
		if len(version.SecretBinary) > 0 {
			entry["SecretBinary"] = base64.StdEncoding.EncodeToString(version.SecretBinary)
		}

		secretValues = append(secretValues, entry)
	}

	result := map[string]interface{}{
		"SecretValues": secretValues,
	}
	if len(apiErrors) > 0 {
		result["Errors"] = apiErrors
	}
	if isTruncated {
		pagination.SetNextToken(result, "NextToken", nextToken)
	}

	return result, nil
}
