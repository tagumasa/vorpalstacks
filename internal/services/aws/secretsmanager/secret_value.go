package secretsmanager

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/store/aws/common"
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

	opts := common.ListOptions{MaxItems: maxResults}
	if nextToken != "" {
		opts.Marker = nextToken
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListSecrets(opts)
	if err != nil {
		return nil, err
	}

	secretList := []interface{}{}
	for _, secret := range result.Items {
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
		s.addRotationFields(entry, secret)
		secretList = append(secretList, entry)
	}

	response := map[string]interface{}{
		"SecretList": secretList,
	}
	if result.NextMarker != "" {
		response["NextToken"] = result.NextMarker
	}
	return response, nil
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
		"ARN":                    secret.ARN,
		"Name":                   secret.Name,
		"SecretId":               secret.ARN,
		"Description":            secret.Description,
		"KmsKeyId":               secret.KmsKeyId,
		"CreatedDate":            secret.CreatedDate.Unix(),
		"LastChangedDate":        secret.LastChangedDate.Unix(),
		"LastAccessedDate":       secret.LastAccessedDate.Unix(),
		"SecretVersionsToStages": s.buildSecretVersionsToStages(reqCtx, secret),
		"Tags":                   s.buildTagsList(secret),
	}
	s.addRotationFields(result, secret)
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

	if removeFromVersionId != "" {
		existingVer, err := store.GetSecretVersion(secret.Name, removeFromVersionId)
		if err != nil {
			return nil, awserrors.NewAWSError("ResourceNotFoundException",
				fmt.Sprintf("Secrets Manager can't find the version %s", removeFromVersionId), http.StatusNotFound)
		}
		newStages := []string{}
		for _, s := range existingVer.VersionStages {
			if s != versionStage {
				newStages = append(newStages, s)
			}
		}
		existingVer.VersionStages = newStages
		if err := store.UpdateSecretVersionStage(secret.Name, removeFromVersionId, newStages); err != nil {
			return nil, mapStoreError(err)
		}
	}

	targetVerId := moveToVersionId
	if targetVerId == "" {
		currentVer, err := store.GetSecretVersionByStage(secret.Name, versionStage)
		if err != nil {
			return nil, awserrors.NewAWSError("ResourceNotFoundException",
				fmt.Sprintf("Secrets Manager can't find the version with stage %s", versionStage), http.StatusNotFound)
		}
		targetVerId = currentVer.VersionId
	}

	targetVer, err := store.GetSecretVersion(secret.Name, targetVerId)
	if err != nil {
		return nil, awserrors.NewAWSError("ResourceNotFoundException",
			fmt.Sprintf("Secrets Manager can't find the version %s", targetVerId), http.StatusNotFound)
	}

	found := false
	for _, s := range targetVer.VersionStages {
		if s == versionStage {
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
