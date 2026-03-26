package secretsmanager

import (
	"context"

	"vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
)

// RestoreSecret restores a deleted secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_RestoreSecret.html
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

// RotateSecret rotates a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_RotateSecret.html
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
	if err := store.RotateSecret(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	updated, err := store.GetSecret(secret.Name)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":       updated.ARN,
		"Name":      updated.Name,
		"VersionId": updated.CurrentVersion,
	}, nil
}

// CancelRotateSecret cancels a secret rotation.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_CancelRotateSecret.html
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
	if err := store.CancelRotation(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}
