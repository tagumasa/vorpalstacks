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
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// PutSecretValue stores a secret value in a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_PutSecretValue.html
func (s *SecretsManagerService) PutSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, awserrors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	secretString := request.GetStringParam(req.Parameters, "SecretString")
	secretBinaryStr := request.GetStringParam(req.Parameters, "SecretBinary")
	versionStages := request.GetStringList(req.Parameters, "VersionStages")

	if secretString != "" {
		secret.SecretString = secretString
		secret.SecretBinary = nil
	}
	if secretBinaryStr != "" {
		decoded, err := decodeSecretBinary(secretBinaryStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SecretBinary encoding: %w", err)
		}
		secret.SecretBinary = decoded
		secret.SecretString = ""
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	updated, err := store.UpdateSecret(secret)
	if err != nil {
		return nil, mapStoreError(err)
	}

	version, err := store.GetSecretVersion(updated.Name, updated.CurrentVersion)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if len(versionStages) > 0 {
		customStages := []string{}
		for _, st := range versionStages {
			if st != "AWSCURRENT" {
				customStages = append(customStages, st)
			}
		}
		customStages = append(customStages, "AWSCURRENT")
		if err := store.UpdateSecretVersionStage(updated.Name, updated.CurrentVersion, customStages); err != nil {
			return nil, fmt.Errorf("failed to apply version stages: %w", err)
		}
		version.VersionStages = customStages
	}

	return map[string]interface{}{
		"ARN":           updated.ARN,
		"Name":          updated.Name,
		"VersionId":     updated.CurrentVersion,
		"VersionStages": version.VersionStages,
	}, nil
}

// ListSecrets lists the secrets in AWS Secrets Manager.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_ListSecrets.html
func (s *SecretsManagerService) ListSecrets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 {
		maxResults = 100
	}
	includePlannedDeletion := request.GetBoolParam(req.Parameters, "IncludePlannedDeletion")
	sortBy := request.GetStringParam(req.Parameters, "SortBy")
	sortOrder := request.GetStringParam(req.Parameters, "SortOrder")
	filters := request.GetListParam(req.Parameters, "Filters")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Load all secrets from the store without pagination, then apply
	// filtering (IncludePlannedDeletion, Filters) and sorting in the handler
	// before paginating with an offset-based NextToken. This avoids the
	// previous bug where filtering was applied after store-level pagination,
	// causing NextToken to never be returned.
	result, err := store.ListSecrets(common.ListOptions{})
	if err != nil {
		return nil, err
	}

	filtered := result.Items
	if !includePlannedDeletion {
		kept := filtered[:0]
		for _, sec := range filtered {
			if sec.DeletedDate == nil && sec.ScheduledDeletionDate == nil {
				kept = append(kept, sec)
			}
		}
		filtered = kept
	}

	if len(filters) > 0 {
		filtered = applySecretFilters(filtered, filters)
	}

	if sortBy != "" {
		sortSecrets(filtered, sortBy, sortOrder)
	}

	// Parse the offset from the NextToken (encoded as a decimal integer).
	// Subsequent requests resume from where the previous page ended.
	skipCount := 0
	if nextToken != "" {
		if n, err := fmt.Sscanf(nextToken, "%d", &skipCount); n != 1 || err != nil {
			skipCount = 0
		}
	}

	paged := filtered[skipCount:]
	isTruncated := maxResults < len(paged)
	if isTruncated {
		paged = paged[:maxResults]
	}

	secretList := []interface{}{}
	for _, secret := range paged {
		entry := map[string]interface{}{
			"ARN":                    secret.ARN,
			"Name":                   secret.Name,
			"Description":            secret.Description,
			"KmsKeyId":               secret.KmsKeyId,
			"CreatedDate":            secret.CreatedDate.Unix(),
			"LastChangedDate":        secret.LastChangedDate.Unix(),
			"LastAccessedDate":       secret.LastAccessedDate.Unix(),
			"SecretVersionsToStages": s.buildSecretVersionsToStages(reqCtx, secret),
		}
		if secret.Type != "" {
			entry["Type"] = secret.Type
		}
		s.addRotationFields(entry, secret)
		if len(secret.ReplicationStatus) > 0 {
			entry["ReplicationStatus"] = buildReplicationStatusResponse(secret.ReplicationStatus)
		}
		secretList = append(secretList, entry)
	}

	response := map[string]interface{}{
		"SecretList": secretList,
	}
	if isTruncated {
		response["NextToken"] = fmt.Sprintf("%d", skipCount+maxResults)
	}
	return response, nil
}

func applySecretFilters(secrets []*secretsmanagerstore.Secret, filters []map[string]interface{}) []*secretsmanagerstore.Secret {
	for _, f := range filters {
		key := request.GetStringParam(f, "Key")
		values := request.GetStringList(f, "Values")
		if key == "" || len(values) == 0 {
			continue
		}
		switch key {
		case "name":
			kept := secrets[:0]
			for _, sec := range secrets {
				for _, v := range values {
					if strings.Contains(sec.Name, v) {
						kept = append(kept, sec)
						break
					}
				}
			}
			secrets = kept
		case "description":
			kept := secrets[:0]
			for _, sec := range secrets {
				for _, v := range values {
					if strings.Contains(sec.Description, v) {
						kept = append(kept, sec)
						break
					}
				}
			}
			secrets = kept
		case "tag-key":
			kept := secrets[:0]
			for _, sec := range secrets {
				for _, v := range values {
					if _, ok := sec.Tags[v]; ok {
						kept = append(kept, sec)
						break
					}
				}
			}
			secrets = kept
		case "tag-value":
			kept := secrets[:0]
			for _, sec := range secrets {
				for _, v := range values {
					found := false
					for _, tv := range sec.Tags {
						if strings.Contains(tv, v) {
							found = true
							break
						}
					}
					if found {
						kept = append(kept, sec)
						break
					}
				}
			}
			secrets = kept
		case "primary-region":
			kept := secrets[:0]
			for _, sec := range secrets {
				for _, v := range values {
					if sec.PrimaryRegion == v {
						kept = append(kept, sec)
						break
					}
				}
			}
			secrets = kept
		case "owning-service":
			kept := secrets[:0]
			for _, sec := range secrets {
				for _, v := range values {
					if sec.OwningService == v {
						kept = append(kept, sec)
						break
					}
				}
			}
			secrets = kept
		}
	}
	return secrets
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
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_DescribeSecret.html
func (s *SecretsManagerService) DescribeSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, awserrors.ErrMissingParameter
	}

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"ARN":                secret.ARN,
		"Name":               secret.Name,
		"SecretId":           secret.ARN,
		"Description":        secret.Description,
		"KmsKeyId":           secret.KmsKeyId,
		"CreatedDate":        secret.CreatedDate.Unix(),
		"LastChangedDate":    secret.LastChangedDate.Unix(),
		"LastAccessedDate":   secret.LastAccessedDate.Unix(),
		"VersionIdsToStages": s.buildSecretVersionsToStages(reqCtx, secret),
		"Tags":               s.buildTagsList(secret),
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
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_ListSecretVersionIds.html
func (s *SecretsManagerService) ListSecretVersionIds(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, awserrors.ErrMissingParameter
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

	versionList := []interface{}{}
	for _, version := range versions {
		entry := map[string]interface{}{
			"VersionId":     version.VersionId,
			"VersionStages": version.VersionStages,
			"CreatedDate":   version.CreatedDate.Unix(),
		}
		versionList = append(versionList, entry)
	}

	return map[string]interface{}{
		"ARN":       secret.ARN,
		"Name":      secret.Name,
		"Versions":  versionList,
		"NextToken": "",
	}, nil
}

// UpdateSecretVersionStage moves a staging label from one version to another.
func (s *SecretsManagerService) UpdateSecretVersionStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, awserrors.ErrMissingParameter
	}
	versionStage := request.GetStringParam(req.Parameters, "VersionStage")
	if versionStage == "" {
		return nil, awserrors.ErrMissingParameter
	}
	moveToVersionId := request.GetStringParam(req.Parameters, "MoveToVersionId")
	removeFromVersionId := request.GetStringParam(req.Parameters, "RemoveFromVersionId")

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	removeVersionId := removeFromVersionId
	if removeVersionId == "" {
		existingVer, err := store.GetSecretVersionByStage(secret.Name, versionStage)
		if err != nil {
			return nil, awserrors.NewAWSError("ResourceNotFoundException",
				fmt.Sprintf("Secrets Manager can't find the version with stage %s", versionStage), http.StatusNotFound)
		}
		removeVersionId = existingVer.VersionId
	}

	if removeVersionId != moveToVersionId {
		existingVer, err := store.GetSecretVersion(secret.Name, removeVersionId)
		if err != nil {
			return nil, awserrors.NewAWSError("ResourceNotFoundException",
				fmt.Sprintf("Secrets Manager can't find the version %s", removeVersionId), http.StatusNotFound)
		}
		newStages := []string{}
		for _, st := range existingVer.VersionStages {
			if st != versionStage {
				newStages = append(newStages, st)
			}
		}
		if err := store.UpdateSecretVersionStage(secret.Name, removeVersionId, newStages); err != nil {
			return nil, mapStoreError(err)
		}
	}

	targetVerId := moveToVersionId
	if targetVerId == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You must specify MoveToVersionId.", http.StatusBadRequest)
	}

	targetVer, err := store.GetSecretVersion(secret.Name, targetVerId)
	if err != nil {
		return nil, awserrors.NewAWSError("ResourceNotFoundException",
			fmt.Sprintf("Secrets Manager can't find the version %s", targetVerId), http.StatusNotFound)
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
		secret.CurrentVersion = targetVerId
		secret.LastChangedDate = time.Now().UTC()

		oldPrevious, err := store.GetSecretVersionByStage(secret.Name, "AWSPREVIOUS")
		if err == nil && oldPrevious.VersionId != targetVerId {
			prevStages := []string{}
			for _, st := range oldPrevious.VersionStages {
				if st != "AWSPREVIOUS" {
					prevStages = append(prevStages, st)
				}
			}
			if err := store.UpdateSecretVersionStage(secret.Name, oldPrevious.VersionId, prevStages); err != nil {
				return nil, mapStoreError(err)
			}
		}

		if removeVersionId != "" && removeVersionId != targetVerId {
			oldCurrentVer, err := store.GetSecretVersion(secret.Name, removeVersionId)
			if err == nil {
				hasPrev := false
				for _, st := range oldCurrentVer.VersionStages {
					if st == "AWSPREVIOUS" {
						hasPrev = true
						break
					}
				}
				if !hasPrev {
					oldCurrentVer.VersionStages = append(oldCurrentVer.VersionStages, "AWSPREVIOUS")
					if err := store.UpdateSecretVersionStage(secret.Name, removeVersionId, oldCurrentVer.VersionStages); err != nil {
						return nil, mapStoreError(err)
					}
				}
			}
		}

		if err := store.UpdateSecretMetadata(secret); err != nil {
			return nil, mapStoreError(err)
		}
	}

	if err := store.UpdateSecretVersionStage(secret.Name, targetVerId, targetVer.VersionStages); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}

// BatchGetSecretValue retrieves multiple secret values.
func (s *SecretsManagerService) BatchGetSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretIdList := request.GetStringList(req.Parameters, "SecretIdList")
	if len(secretIdList) == 0 {
		return nil, awserrors.ErrMissingParameter
	}
	if len(secretIdList) > 20 {
		return nil, awserrors.NewAWSError("ValidationException",
			"You can include up to 20 secrets in a batch.", http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	secretValues := []interface{}{}
	apiErrors := []interface{}{}

	for _, secretId := range secretIdList {
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

	return result, nil
}
