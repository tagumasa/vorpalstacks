package secretsmanager

import (
	"context"
	"encoding/base64"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateSecret creates a new secret in AWS Secrets Manager.
func (s *SecretsManagerService) CreateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
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

	result, err := s.createSecretCore(ctx, store, CreateSecretInput{
		Name:                        request.GetStringParam(req.Parameters, "Name"),
		SecretString:                request.GetStringParam(req.Parameters, "SecretString"),
		SecretBinaryB64:             request.GetStringParam(req.Parameters, "SecretBinary"),
		Description:                 request.GetStringParam(req.Parameters, "Description"),
		KmsKeyId:                    request.GetStringParam(req.Parameters, "KmsKeyId"),
		Type:                        request.GetStringParam(req.Parameters, "Type"),
		ClientRequestToken:          request.GetStringParam(req.Parameters, "ClientRequestToken"),
		Tags:                        tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")),
		AddReplicaRegions:           addReplicaRegions,
		ForceOverwriteReplicaSecret: request.GetBoolParam(req.Parameters, "ForceOverwriteReplicaSecret"),
		Region:                      reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"ARN":       result.ARN,
		"Name":      result.Name,
		"VersionId": result.VersionID,
	}

	if len(result.ReplicationStatus) > 0 {
		response["ReplicationStatus"] = buildReplicationStatusResponse(result.ReplicationStatus)
	}

	return response, nil
}

// GetSecretValue returns the secret value for a secret.
func (s *SecretsManagerService) GetSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getSecretValueCore(ctx, store, GetSecretValueInput{
		SecretId:     request.GetStringParam(req.Parameters, "SecretId"),
		VersionId:    request.GetStringParam(req.Parameters, "VersionId"),
		VersionStage: request.GetStringParam(req.Parameters, "VersionStage"),
	})
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"ARN":           result.Secret.ARN,
		"Name":          result.Secret.Name,
		"VersionId":     result.Version.VersionId,
		"VersionStages": result.Version.VersionStages,
		"CreatedDate":   result.Version.CreatedDate.Unix(),
	}

	if result.Version.SecretString != "" {
		response["SecretString"] = result.Version.SecretString
	}
	if len(result.Version.SecretBinary) > 0 {
		response["SecretBinary"] = base64.StdEncoding.EncodeToString(result.Version.SecretBinary)
	}

	return response, nil
}

// UpdateSecret updates an existing secret.
func (s *SecretsManagerService) UpdateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.updateSecretCore(ctx, store, UpdateSecretInput{
		SecretId:           request.GetStringParam(req.Parameters, "SecretId"),
		SecretString:       request.GetStringParam(req.Parameters, "SecretString"),
		SecretBinaryB64:    request.GetStringParam(req.Parameters, "SecretBinary"),
		Description:        request.GetStringParam(req.Parameters, "Description"),
		KmsKeyId:           request.GetStringParam(req.Parameters, "KmsKeyId"),
		Type:               request.GetStringParam(req.Parameters, "Type"),
		ClientRequestToken: request.GetStringParam(req.Parameters, "ClientRequestToken"),
		HasDescription:     request.HasParam(req.Parameters, "Description"),
		HasKmsKeyId:        request.HasParam(req.Parameters, "KmsKeyId"),
		HasType:            request.HasParam(req.Parameters, "Type"),
	})
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}
	if result.VersionID != "" {
		response["VersionId"] = result.VersionID
	}
	return response, nil
}

// DeleteSecret deletes a secret.
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
