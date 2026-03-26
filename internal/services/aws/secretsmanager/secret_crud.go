package secretsmanager

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// CreateSecret creates a new secret in AWS Secrets Manager.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_CreateSecret.html
func (s *SecretsManagerService) CreateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, errors.ErrMissingParameter
	}

	secretString := request.GetStringParam(req.Parameters, "SecretString")
	secretBinaryStr := request.GetStringParam(req.Parameters, "SecretBinary")
	description := request.GetStringParam(req.Parameters, "Description")
	kmsKeyId := request.GetStringParam(req.Parameters, "KmsKeyId")
	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	var secretBinary []byte
	if secretBinaryStr != "" {
		var err error
		secretBinary, err = decodeSecretBinary(secretBinaryStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SecretBinary encoding: %w", err)
		}
	}

	secret := secretsmanagerstore.NewSecret(name)
	secret.SecretString = secretString
	secret.SecretBinary = secretBinary
	secret.Description = description
	secret.KmsKeyId = kmsKeyId
	if len(tags) > 0 {
		secret.Tags = tags
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.CreateSecret(secret)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":       created.ARN,
		"Name":      created.Name,
		"VersionId": created.CurrentVersion,
	}, nil
}

// GetSecretValue returns the secret value for a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_GetSecretValue.html
func (s *SecretsManagerService) GetSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	versionId := request.GetStringParam(req.Parameters, "VersionId")
	versionStage := request.GetStringParam(req.Parameters, "VersionStage")

	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	if versionId != "" && versionStage != "" {
		return nil, errors.NewAWSError("InvalidParameterException",
			"You can't specify both VersionStage and VersionId.", http.StatusBadRequest)
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if versionId == "" && versionStage != "" {
		version, err := store.GetSecretVersionByStage(secret.Name, versionStage)
		if err != nil {
			return nil, mapStoreError(err)
		}
		versionId = version.VersionId
	} else if versionId == "" {
		versionId = secret.CurrentVersion
	}

	version, err := store.GetSecretVersion(secret.Name, versionId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := map[string]interface{}{
		"ARN":           secret.ARN,
		"Name":          secret.Name,
		"VersionId":     version.VersionId,
		"VersionStages": version.VersionStages,
		"CreatedDate":   version.CreatedDate.Unix(),
	}

	if version.SecretString != "" {
		result["SecretString"] = version.SecretString
	}
	if len(version.SecretBinary) > 0 {
		result["SecretBinary"] = base64.StdEncoding.EncodeToString(version.SecretBinary)
	}

	return result, nil
}

// UpdateSecret updates an existing secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_UpdateSecret.html
func (s *SecretsManagerService) UpdateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	secretString := request.GetStringParam(req.Parameters, "SecretString")
	secretBinaryStr := request.GetStringParam(req.Parameters, "SecretBinary")
	description := request.GetStringParam(req.Parameters, "Description")
	kmsKeyId := request.GetStringParam(req.Parameters, "KmsKeyId")

	hasSecretValue := secretString != "" || secretBinaryStr != ""

	if secretString != "" {
		secret.SecretString = secretString
		secret.SecretBinary = nil
	}
	if secretBinaryStr != "" {
		decoded, err := decodeSecretBinary(secretBinaryStr)
		if err != nil {
			return nil, err
		}
		secret.SecretBinary = decoded
		secret.SecretString = ""
	}
	if description != "" {
		secret.Description = description
	}
	if kmsKeyId != "" {
		secret.KmsKeyId = kmsKeyId
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	updated, err := store.UpdateSecret(secret)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := map[string]interface{}{
		"ARN":  updated.ARN,
		"Name": updated.Name,
	}
	if hasSecretValue && updated.CurrentVersion != "" {
		result["VersionId"] = updated.CurrentVersion
	}
	return result, nil
}

// DeleteSecret deletes a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_DeleteSecret.html
func (s *SecretsManagerService) DeleteSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	recoveryWindowInDays := request.GetIntParam(req.Parameters, "RecoveryWindowInDays")
	forceDeleteWithoutRecovery := request.GetBoolParam(req.Parameters, "ForceDeleteWithoutRecovery")

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var deletionDate time.Time
	if forceDeleteWithoutRecovery {
		if err := store.DeleteSecret(secret.Name); err != nil {
			return nil, mapStoreError(err)
		}
		deletionDate = time.Now().UTC()
	} else {
		if recoveryWindowInDays <= 0 {
			recoveryWindowInDays = 30
		}
		deletionDate = time.Now().UTC().AddDate(0, 0, recoveryWindowInDays)
		if err := store.ScheduleDeletion(secret.Name, deletionDate); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return map[string]interface{}{
		"ARN":          secret.ARN,
		"Name":         secret.Name,
		"DeletionDate": deletionDate.Unix(),
	}, nil
}
