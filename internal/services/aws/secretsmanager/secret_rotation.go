package secretsmanager

import (
	"context"

	"vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

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
	secret.RotationLambdaARN = rotationLambdaARN
	if automaticallyAfterDays > 0 {
		if secret.RotationRules == nil {
			secret.RotationRules = &secretsmanagerstore.RotationRules{}
		}
		secret.RotationRules.AutomaticallyAfterDays = automaticallyAfterDays
	}

	if secret.CurrentVersion != "" {
		currentVer, verErr := store.GetSecretVersion(secret.Name, secret.CurrentVersion)
		if verErr == nil {
			newVersion := secretsmanagerstore.NewSecretVersion("")
			newVersion.SecretName = secret.Name
			newVersion.SecretString = currentVer.SecretString
			newVersion.SecretBinary = make([]byte, len(currentVer.SecretBinary))
			copy(newVersion.SecretBinary, currentVer.SecretBinary)
			newVersion.VersionStages = []string{"AWSCURRENT"}

			if createErr := store.CreateVersionDirect(secret.Name, newVersion); createErr == nil {
				oldPrevious, prevErr := store.GetSecretVersionByStage(secret.Name, "AWSPREVIOUS")
				if prevErr == nil && oldPrevious.VersionId != newVersion.VersionId {
					prevStages := []string{}
					for _, st := range oldPrevious.VersionStages {
						if st != "AWSPREVIOUS" {
							prevStages = append(prevStages, st)
						}
					}
					_ = store.UpdateSecretVersionStage(secret.Name, oldPrevious.VersionId, prevStages)
				}

				_ = store.UpdateSecretVersionStage(secret.Name, secret.CurrentVersion, []string{"AWSPREVIOUS"})

				secret.VersionIDs = append(secret.VersionIDs, newVersion.VersionId)
				secret.CurrentVersion = newVersion.VersionId
			}
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
		"VersionId": secret.CurrentVersion,
	}, nil
}

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
		_ = store.UpdateSecretVersionStage(secret.Name, pendingVer.VersionId, newStages)
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
