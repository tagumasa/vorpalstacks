package secretsmanager

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// bytesEqual reports whether two byte slices are equal, treating nil and
// empty slices as equivalent.
func bytesEqual(a, b []byte) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return bytes.Equal(a, b)
}

// CreateSecret creates a new secret in AWS Secrets Manager.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_CreateSecret.html
func (s *SecretsManagerService) CreateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if err := validateSecretName(name); err != nil {
		return nil, err
	}

	secretBinary, err := decodeAndValidateSecretBinary(request.GetStringParam(req.Parameters, "SecretBinary"))
	if err != nil {
		return nil, err
	}

	result, err := s.createSecretCore(ctx, store, CreateSecretInput{
		Name:               name,
		SecretString:       request.GetStringParam(req.Parameters, "SecretString"),
		SecretBinary:       secretBinary,
		Description:        request.GetStringParam(req.Parameters, "Description"),
		KmsKeyId:           request.GetStringParam(req.Parameters, "KmsKeyId"),
		Type:               request.GetStringParam(req.Parameters, "Type"),
		ClientRequestToken: request.GetStringParam(req.Parameters, "ClientRequestToken"),
		Tags:               tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")),
	})
	if err != nil {
		return nil, err
	}

	// AddReplicaRegions — when provided, replicate the secret to
	// the specified regions immediately after creation (same as calling
	// ReplicateSecretToRegions separately).
	addReplicaRegions, err := parseReplicaRegions(req.Parameters, "AddReplicaRegions")
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"ARN":       result.ARN,
		"Name":      result.Name,
		"VersionId": result.VersionID,
	}

	if len(addReplicaRegions) > 0 {
		if s.storageManager == nil {
			return nil, errors.NewAWSError("InvalidRequestException",
				"Replication is not configured for this service.", http.StatusBadRequest)
		}
		created, metaErr := store.GetSecretForMetadata(result.Name)
		if metaErr != nil {
			return nil, mapStoreError(metaErr)
		}
		forceOverwrite := request.GetBoolParam(req.Parameters, "ForceOverwriteReplicaSecret")
		s.replicateSecretToRegions(store, created, addReplicaRegions, forceOverwrite, reqCtx.GetRegion())
		if err := store.UpdateSecretMetadata(created); err != nil {
			return nil, mapStoreError(err)
		}
		if len(created.ReplicationStatus) > 0 {
			response["ReplicationStatus"] = buildReplicationStatusResponse(created.ReplicationStatus)
		}
	}

	return response, nil
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

	var version *secretsmanagerstore.SecretVersion
	if versionId == "" && versionStage != "" {
		var err error
		version, err = store.GetSecretVersionByStage(secret.Name, versionStage)
		if err != nil {
			return nil, mapStoreError(err)
		}
	} else if versionId == "" {
		if secret.CurrentVersion == "" {
			return nil, errors.NewAWSError("ResourceNotFoundException",
				"Secrets Manager can't find the version for the secret because no version has been created.", http.StatusNotFound)
		}
		var err error
		version, err = store.GetSecretVersion(secret.Name, secret.CurrentVersion)
		if err != nil {
			return nil, mapStoreError(err)
		}
	} else {
		if isStageLabel(versionId) {
			if secret.CurrentVersion == "" {
				return nil, errors.NewAWSError("ResourceNotFoundException",
					"Secrets Manager can't find the version for the secret because no version has been created.", http.StatusNotFound)
			}
			var err error
			version, err = store.GetSecretVersionByStage(secret.Name, versionId)
			if err != nil {
				return nil, mapStoreError(err)
			}
		} else {
			var err error
			version, err = store.GetSecretVersion(secret.Name, versionId)
			if err != nil {
				return nil, mapStoreError(err)
			}
		}
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
	if err := validateSecretValueMutex(secretString, secretBinaryStr); err != nil {
		return nil, err
	}
	if err := validateSecretStringLength(secretString); err != nil {
		return nil, err
	}
	description := request.GetStringParam(req.Parameters, "Description")
	if err := validateDescription(description); err != nil {
		return nil, err
	}
	kmsKeyId := request.GetStringParam(req.Parameters, "KmsKeyId")
	secretType := request.GetStringParam(req.Parameters, "Type")
	clientRequestToken := request.GetStringParam(req.Parameters, "ClientRequestToken")
	if err := validateClientRequestToken(clientRequestToken); err != nil {
		return nil, err
	}

	hasSecretValue := secretString != "" || secretBinaryStr != ""

	if secretString != "" {
		secret.SecretString = secretString
		secret.SecretBinary = nil
	} else if secretBinaryStr != "" {
		decoded, decErr := decodeAndValidateSecretBinary(secretBinaryStr)
		if decErr != nil {
			return nil, decErr
		}
		secret.SecretBinary = decoded
		secret.SecretString = ""
	}
	if request.HasParam(req.Parameters, "Description") {
		secret.Description = description
	}
	if request.HasParam(req.Parameters, "KmsKeyId") {
		secret.KmsKeyId = kmsKeyId
	}
	if request.HasParam(req.Parameters, "Type") {
		secret.Type = secretType
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var versionId string
	if hasSecretValue {
		// ClientRequestToken becomes the version ID when a new
		// version is created.
		secret.InitialVersionId = clientRequestToken
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

	result := map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}
	if versionId != "" {
		result["VersionId"] = versionId
	}
	return result, nil
}

// DeleteSecret deletes a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_DeleteSecret.html
func (s *SecretsManagerService) DeleteSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.deleteSecretCore(ctx, store, DeleteSecretInput{
		SecretId:                   request.GetStringParam(req.Parameters, "SecretId"),
		ForceDeleteWithoutRecovery: request.GetBoolParam(req.Parameters, "ForceDeleteWithoutRecovery"),
		RecoveryWindowInDays:       request.GetIntParam(req.Parameters, "RecoveryWindowInDays"),
		HasRecoveryWindow:          request.HasParam(req.Parameters, "RecoveryWindowInDays"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":          result.ARN,
		"Name":         result.Name,
		"DeletionDate": result.DeletionDate.Unix(),
	}, nil
}

func isStageLabel(s string) bool {
	switch s {
	case "AWSCURRENT", "AWSPREVIOUS", "AWSPENDING":
		return true
	default:
		return false
	}
}
