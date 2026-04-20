package secretsmanager

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// RestoreSecret restores a previously deleted secret. Secrets Manager keeps
// deleted secrets recoverable for a minimum of 30 days.
func (s *SecretsManagerService) RestoreSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.RestoreSecret(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}

// RotateSecret configures rotation for a secret and, if a RotationLambdaARN
// is available, performs the full AWS rotation protocol:
//
//  1. Update rotation metadata (RotationEnabled, RotationLambdaARN, RotationRules).
//  2. If a rotation Lambda is configured, execute the multi-step rotation:
//     a. createSecret — Lambda generates new secret material with AWSPENDING stage.
//     b. setSecret    — Lambda configures the target resource with the new value.
//     c. testSecret   — Lambda verifies the new secret works.
//     d. finishSecret — Service promotes AWSPENDING to AWSCURRENT.
//  3. If no rotation Lambda is set, fall back to copying the current version
//     (metadata-only rotation, preserving backward compatibility).
func (s *SecretsManagerService) RotateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rotationLambdaARN := request.GetStringParam(req.Parameters, "RotationLambdaARN")
	automaticallyAfterDays := request.GetIntParam(req.Parameters, "RotationRules.AutomaticallyAfterDays")

	secret.RotationEnabled = true
	if rotationLambdaARN != "" {
		secret.RotationLambdaARN = rotationLambdaARN
	}
	if automaticallyAfterDays > 0 {
		if secret.RotationRules == nil {
			secret.RotationRules = &secretsmanagerstore.RotationRules{}
		}
		secret.RotationRules.AutomaticallyAfterDays = automaticallyAfterDays
	}

	var versionId string

	if secret.RotationLambdaARN != "" && s.bus != nil {
		if rotErr := s.executeRotation(ctx, store, secret); rotErr != nil {
			return nil, errors.NewAWSError("InvalidRequestException",
				rotErr.Error(), 400)
		}
		versionId = secret.CurrentVersion
	} else {
		versionId, err = s.executeMetadataOnlyRotation(store, secret)
		if err != nil {
			return nil, err
		}
	}

	secret.LastRotatedDate = storeClock()
	if automaticallyAfterDays > 0 {
		secret.NextRotationDate = secret.LastRotatedDate.AddDate(0, 0, automaticallyAfterDays)
	}

	if err := store.UpdateSecretMetadata(secret); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":       secret.ARN,
		"Name":      secret.Name,
		"VersionId": versionId,
	}, nil
}

// executeMetadataOnlyRotation performs a rotation without invoking a Lambda
// function. It copies the current version to a new version, mirroring the
// pre-existing behaviour when no rotation Lambda is configured.
func (s *SecretsManagerService) executeMetadataOnlyRotation(store secretsmanagerstore.SecretStoreInterface, secret *secretsmanagerstore.Secret) (string, error) {
	if secret.CurrentVersion == "" {
		return "", nil
	}

	currentVer, verErr := store.GetSecretVersion(secret.Name, secret.CurrentVersion)
	if verErr != nil {
		return "", verErr
	}

	newVersion := secretsmanagerstore.NewSecretVersion("")
	newVersion.SecretName = secret.Name
	newVersion.SecretString = currentVer.SecretString
	newVersion.SecretBinary = make([]byte, len(currentVer.SecretBinary))
	copy(newVersion.SecretBinary, currentVer.SecretBinary)
	newVersion.VersionStages = []string{"AWSCURRENT"}

	if createErr := store.CreateVersionDirect(secret.Name, newVersion); createErr != nil {
		return "", createErr
	}

	oldPrevious, prevErr := store.GetSecretVersionByStage(secret.Name, "AWSPREVIOUS")
	if prevErr == nil && oldPrevious.VersionId != newVersion.VersionId {
		prevStages := []string{}
		for _, st := range oldPrevious.VersionStages {
			if st != "AWSPREVIOUS" {
				prevStages = append(prevStages, st)
			}
		}
		if err := store.UpdateSecretVersionStage(secret.Name, oldPrevious.VersionId, prevStages); err != nil {
			return "", fmt.Errorf("failed to clean AWSPREVIOUS from old previous version: %w", err)
		}
	}

	if err := store.UpdateSecretVersionStage(secret.Name, secret.CurrentVersion, []string{"AWSPREVIOUS"}); err != nil {
		return "", fmt.Errorf("failed to demote current version to AWSPREVIOUS: %w", err)
	}

	secret.VersionIDs = append(secret.VersionIDs, newVersion.VersionId)
	secret.CurrentVersion = newVersion.VersionId

	return newVersion.VersionId, nil
}

// CancelRotateSecret disables automatic rotation for a secret. If a version
// with the AWSPENDING stage exists, that stage label is removed before the
// rotation configuration is cleared.
func (s *SecretsManagerService) CancelRotateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pendingVer, err := store.GetSecretVersionByStage(secret.Name, "AWSPENDING")
	versionId := ""
	if err == nil {
		versionId = pendingVer.VersionId
		newStages := []string{}
		for _, st := range pendingVer.VersionStages {
			if st != "AWSPENDING" {
				newStages = append(newStages, st)
			}
		}
		if err := store.UpdateSecretVersionStage(secret.Name, pendingVer.VersionId, newStages); err != nil {
			return nil, fmt.Errorf("failed to remove AWSPENDING stage during cancel: %w", err)
		}
	}

	if err := store.CancelRotation(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	result := map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}
	if versionId != "" {
		result["VersionId"] = versionId
	}
	return result, nil
}
