package secretsmanager

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
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
	name := request.GetStringParam(req.Parameters, "Name")
	if err := validateSecretName(name); err != nil {
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
	parsedTags := tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	if err := validateSecretTags(parsedTags); err != nil {
		return nil, err
	}
	tags := tagutil.ToMap(parsedTags)

	secretBinary, err := decodeAndValidateSecretBinary(secretBinaryStr)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// ClientRequestToken idempotency: if the token is provided and a
	// version with that ID already exists for a secret with this name,
	// check whether the values match (idempotent success) or differ
	// (error). This mirrors the AWS CreateSecret idempotency rules.
	if clientRequestToken != "" {
		existing, metaErr := store.GetSecretForMetadata(name)
		if metaErr == nil && len(existing.VersionIDs) > 0 {
			for _, vid := range existing.VersionIDs {
				if vid == clientRequestToken {
					existingVer, verErr := store.GetSecretVersion(name, clientRequestToken)
					if verErr != nil {
						break
					}
					if existingVer.SecretString == secretString &&
						bytesEqual(existingVer.SecretBinary, secretBinary) {
						return map[string]interface{}{
							"ARN":       existing.ARN,
							"Name":      existing.Name,
							"VersionId": clientRequestToken,
						}, nil
					}
					return nil, errors.NewAWSError("InvalidRequestException",
						fmt.Sprintf("You can't modify an existing secret version. The ClientRequestToken %s is already associated with a different version value.", clientRequestToken),
						http.StatusBadRequest)
				}
			}
		}
	}

	secret := secretsmanagerstore.NewSecret(name)
	secret.SecretString = secretString
	secret.SecretBinary = secretBinary
	secret.Description = description
	secret.KmsKeyId = kmsKeyId
	secret.Type = secretType
	secret.InitialVersionId = clientRequestToken
	if len(tags) > 0 {
		secret.Tags = tags
	}

	created, err := store.CreateSecret(secret)
	if err != nil {
		return nil, mapStoreError(err)
	}

	// AddReplicaRegions — when provided, replicate the secret to
	// the specified regions immediately after creation (same as calling
	// ReplicateSecretToRegions separately).
	addReplicaRegions, err := parseReplicaRegions(req.Parameters, "AddReplicaRegions")
	if err != nil {
		return nil, err
	}
	if len(addReplicaRegions) > 0 {
		if s.storageManager == nil {
			return nil, errors.NewAWSError("InvalidRequestException",
				"Replication is not configured for this service.", http.StatusBadRequest)
		}
		forceOverwrite := request.GetBoolParam(req.Parameters, "ForceOverwriteReplicaSecret")
		s.replicateSecretToRegions(store, created, addReplicaRegions, forceOverwrite, reqCtx.GetRegion())
		if err := store.UpdateSecretMetadata(created); err != nil {
			return nil, mapStoreError(err)
		}
	}

	// Include ReplicationStatus in the response when replication was
	// requested.
	result := map[string]interface{}{
		"ARN":       created.ARN,
		"Name":      created.Name,
		"VersionId": created.CurrentVersion,
	}
	if len(created.ReplicationStatus) > 0 {
		result["ReplicationStatus"] = buildReplicationStatusResponse(created.ReplicationStatus)
	}
	return result, nil
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
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	recoveryWindowInDays := request.GetIntParam(req.Parameters, "RecoveryWindowInDays")
	forceDeleteWithoutRecovery := request.GetBoolParam(req.Parameters, "ForceDeleteWithoutRecovery")
	hasRecoveryWindow := request.HasParam(req.Parameters, "RecoveryWindowInDays")
	if err := validateRecoveryWindow(recoveryWindowInDays, hasRecoveryWindow, forceDeleteWithoutRecovery); err != nil {
		return nil, err
	}

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	// You can't delete a primary secret that is replicated to other Regions.
	if len(secret.ReplicationStatus) > 0 {
		return nil, errors.NewAWSError("InvalidRequestException",
			"You can't delete a primary secret that is replicated to other Regions. Remove the replicas first.", http.StatusBadRequest)
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
		if !hasRecoveryWindow {
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

func isStageLabel(s string) bool {
	switch s {
	case "AWSCURRENT", "AWSPREVIOUS", "AWSPENDING":
		return true
	default:
		return false
	}
}
