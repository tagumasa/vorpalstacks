package secretsmanager

import (
	"context"
	"encoding/base64"

	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// maxResultsParam extracts a page-size parameter with explicit presence
// semantics: nil means the parameter was absent (the Core applies its
// default), a non-nil pointer carries the caller's value verbatim so an
// explicit 0 can be rejected as out of range.
func maxResultsParam(params map[string]interface{}) *int {
	if !request.HasParam(params, "MaxResults") {
		return nil
	}
	v := request.GetIntParam(params, "MaxResults")
	return &v
}

// parseSecretFilters converts the wire Filter list into the
// transport-agnostic SecretFilter DTO.
func parseSecretFilters(params map[string]interface{}, key string) []SecretFilter {
	raw := request.GetListParam(params, key)
	filters := make([]SecretFilter, 0, len(raw))
	for _, f := range raw {
		filters = append(filters, SecretFilter{
			Key:    request.GetStringParam(f, "Key"),
			Values: request.GetStringList(f, "Values"),
		})
	}
	return filters
}

// PutSecretValue stores a secret value in a secret.
func (s *SecretsManagerService) PutSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putSecretValueCore(ctx, store, PutSecretValueInput{
		SecretId:           request.GetStringParam(req.Parameters, "SecretId"),
		SecretString:       request.GetStringParam(req.Parameters, "SecretString"),
		SecretBinaryB64:    request.GetStringParam(req.Parameters, "SecretBinary"),
		ClientRequestToken: request.GetStringParam(req.Parameters, "ClientRequestToken"),
		VersionStages:      request.GetStringList(req.Parameters, "VersionStages"),
		RotationToken:      request.GetStringParam(req.Parameters, "RotationToken"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":           result.ARN,
		"Name":          result.Name,
		"VersionId":     result.VersionID,
		"VersionStages": result.VersionStages,
	}, nil
}

// ListSecrets lists the secrets in AWS Secrets Manager.
func (s *SecretsManagerService) ListSecrets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listSecretsCore(ctx, store, ListSecretsInput{
		MaxResults:             maxResultsParam(req.Parameters),
		NextToken:              pagination.GetMarker(req.Parameters, "NextToken"),
		SortBy:                 request.GetStringParam(req.Parameters, "SortBy"),
		SortOrder:              request.GetStringParam(req.Parameters, "SortOrder"),
		Filters:                parseSecretFilters(req.Parameters, "Filters"),
		IncludePlannedDeletion: request.GetBoolParam(req.Parameters, "IncludePlannedDeletion"),
	})
	if err != nil {
		return nil, err
	}

	secretList := make([]interface{}, 0, len(result.Secrets))
	for _, secret := range result.Secrets {
		entry := map[string]interface{}{
			"ARN":                    secret.ARN,
			"Name":                   secret.Name,
			"CreatedDate":            secret.CreatedDate.Unix(),
			"LastChangedDate":        secret.LastChangedDate.Unix(),
			"SecretVersionsToStages": result.Stages[secret.Name],
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
		if secret.OwningService != "" {
			entry["OwningService"] = secret.OwningService
		}
		if secret.PrimaryRegion != "" {
			entry["PrimaryRegion"] = secret.PrimaryRegion
		}
		s.addRotationFields(entry, secret)
		if len(secret.Tags) > 0 {
			entry["Tags"] = s.buildTagsList(secret)
		}
		if len(secret.ReplicationStatus) > 0 {
			entry["ReplicationStatus"] = buildReplicationStatusResponse(secret.ReplicationStatus)
		}
		addExternalRotationFields(entry, secret)
		secretList = append(secretList, entry)
	}

	response := map[string]interface{}{
		"SecretList": secretList,
	}
	if result.NextToken != "" {
		pagination.SetNextToken(response, "NextToken", result.NextToken)
	}
	return response, nil
}

// DescribeSecret returns the metadata for a secret.
func (s *SecretsManagerService) DescribeSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeSecretCore(ctx, store, DescribeSecretInput{
		SecretId: request.GetStringParam(req.Parameters, "SecretId"),
	})
	if err != nil {
		return nil, err
	}

	secret := result.Secret
	response := map[string]interface{}{
		"ARN":                secret.ARN,
		"Name":               secret.Name,
		"CreatedDate":        secret.CreatedDate.Unix(),
		"LastChangedDate":    secret.LastChangedDate.Unix(),
		"VersionIdsToStages": result.VersionIdsToStages,
	}
	if secret.Description != "" {
		response["Description"] = secret.Description
	}
	if secret.KmsKeyId != "" {
		response["KmsKeyId"] = secret.KmsKeyId
	}
	if !secret.LastAccessedDate.IsZero() {
		response["LastAccessedDate"] = secret.LastAccessedDate.Unix()
	}
	tags := s.buildTagsList(secret)
	if len(tags) > 0 {
		response["Tags"] = tags
	}
	if secret.DeletedDate != nil {
		response["DeletedDate"] = secret.DeletedDate.Unix()
	}
	s.addRotationFields(response, secret)
	if secret.OwningService != "" {
		response["OwningService"] = secret.OwningService
	}
	if secret.PrimaryRegion != "" {
		response["PrimaryRegion"] = secret.PrimaryRegion
	}
	if secret.Type != "" {
		response["Type"] = secret.Type
	}
	if !secret.NextRotationDate.IsZero() {
		response["NextRotationDate"] = secret.NextRotationDate.Unix()
	}
	if len(secret.ReplicationStatus) > 0 {
		response["ReplicationStatus"] = buildReplicationStatusResponse(secret.ReplicationStatus)
	}
	addExternalRotationFields(response, secret)
	return response, nil
}

// ListSecretVersionIds lists the versions of a secret.
// Supports MaxResults (1-100), NextToken pagination, and IncludeDeprecated
// filtering.
func (s *SecretsManagerService) ListSecretVersionIds(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listSecretVersionIdsCore(ctx, store, ListSecretVersionIdsInput{
		SecretId:          request.GetStringParam(req.Parameters, "SecretId"),
		MaxResults:        maxResultsParam(req.Parameters),
		NextToken:         pagination.GetMarker(req.Parameters, "NextToken"),
		IncludeDeprecated: request.GetBoolParam(req.Parameters, "IncludeDeprecated"),
	})
	if err != nil {
		return nil, err
	}

	versionList := make([]interface{}, 0, len(result.Versions))
	for _, version := range result.Versions {
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

	response := map[string]interface{}{
		"ARN":      result.ARN,
		"Name":     result.Name,
		"Versions": versionList,
	}
	if result.NextToken != "" {
		pagination.SetNextToken(response, "NextToken", result.NextToken)
	}
	return response, nil
}

// UpdateSecretVersionStage modifies the staging labels attached to a version.
// Supports three modes per AWS documentation:
//  1. Add a label to a version (MoveToVersionId only).
//  2. Remove a label from a version (RemoveFromVersionId only).
//  3. Move a label between versions (both MoveToVersionId and RemoveFromVersionId).
func (s *SecretsManagerService) UpdateSecretVersionStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.updateSecretVersionStageCore(ctx, store, UpdateSecretVersionStageInput{
		SecretId:            request.GetStringParam(req.Parameters, "SecretId"),
		VersionStage:        request.GetStringParam(req.Parameters, "VersionStage"),
		MoveToVersionId:     request.GetStringParam(req.Parameters, "MoveToVersionId"),
		RemoveFromVersionId: request.GetStringParam(req.Parameters, "RemoveFromVersionId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}, nil
}

// BatchGetSecretValue retrieves multiple secret values.
// You must include either SecretIdList or Filters, but not both.
// MaxResults (1-20) and NextToken pagination are supported when using
// Filters.
func (s *SecretsManagerService) BatchGetSecretValue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.batchGetSecretValueCore(ctx, store, BatchGetSecretValueInput{
		SecretIdList: request.GetStringList(req.Parameters, "SecretIdList"),
		Filters:      parseSecretFilters(req.Parameters, "Filters"),
		MaxResults:   maxResultsParam(req.Parameters),
		NextToken:    pagination.GetMarker(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	secretValues := make([]interface{}, 0, len(result.SecretValues))
	for _, entry := range result.SecretValues {
		e := map[string]interface{}{
			"ARN":           entry.Secret.ARN,
			"Name":          entry.Secret.Name,
			"VersionId":     entry.Version.VersionId,
			"VersionStages": entry.Version.VersionStages,
			"CreatedDate":   entry.Version.CreatedDate.Unix(),
		}

		if entry.Version.SecretString != "" {
			e["SecretString"] = entry.Version.SecretString
		}
		if len(entry.Version.SecretBinary) > 0 {
			e["SecretBinary"] = base64.StdEncoding.EncodeToString(entry.Version.SecretBinary)
		}

		secretValues = append(secretValues, e)
	}

	response := map[string]interface{}{
		"SecretValues": secretValues,
	}
	if len(result.Errors) > 0 {
		errors := make([]interface{}, 0, len(result.Errors))
		for _, apiErr := range result.Errors {
			errors = append(errors, map[string]interface{}{
				"SecretId":  apiErr.SecretId,
				"ErrorCode": apiErr.ErrorCode,
				"Message":   apiErr.Message,
			})
		}
		response["Errors"] = errors
	}
	if result.NextToken != "" {
		pagination.SetNextToken(response, "NextToken", result.NextToken)
	}

	return response, nil
}

// addExternalRotationFields emits the managed external secret rotation
// members on the DescribeSecret and ListSecrets entry shapes.
func addExternalRotationFields(m map[string]interface{}, secret *secretsmanagerstore.Secret) {
	if len(secret.ExternalSecretRotationMetadata) > 0 {
		items := make([]interface{}, 0, len(secret.ExternalSecretRotationMetadata))
		for _, item := range secret.ExternalSecretRotationMetadata {
			items = append(items, map[string]interface{}{
				"Key":   item.Key,
				"Value": item.Value,
			})
		}
		m["ExternalSecretRotationMetadata"] = items
	}
	if secret.ExternalSecretRotationRoleArn != "" {
		m["ExternalSecretRotationRoleArn"] = secret.ExternalSecretRotationRoleArn
	}
}
