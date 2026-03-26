package secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
)

// GetResourcePolicy returns the resource policy for a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_GetResourcePolicy.html
func (s *SecretsManagerService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	policy, err := store.GetResourcePolicy(secret.Name)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}
	if policy != "" {
		result["ResourcePolicy"] = policy
	}

	return result, nil
}

// PutResourcePolicy attaches a resource policy to a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_PutResourcePolicy.html
func (s *SecretsManagerService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	policy := request.GetStringParam(req.Parameters, "ResourcePolicy")
	if policy == "" {
		return nil, errors.ErrMissingParameter
	}

	var js interface{}
	if err := json.Unmarshal([]byte(policy), &js); err != nil {
		return nil, fmt.Errorf("invalid resource policy: %w", err)
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.PutResourcePolicy(secret.Name, policy); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}

// DeleteResourcePolicy deletes the resource policy from a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_DeleteResourcePolicy.html
func (s *SecretsManagerService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	if err := store.DeleteResourcePolicy(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}

// ValidateResourcePolicy validates a resource policy for a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_ValidateResourcePolicy.html
func (s *SecretsManagerService) ValidateResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	policy := request.GetStringParam(req.Parameters, "ResourcePolicy")

	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	_, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"PolicyValidationPassed": true,
		"ValidationErrors":       []interface{}{},
	}

	if policy != "" {
		var js interface{}
		if err := json.Unmarshal([]byte(policy), &js); err != nil {
			log.Printf("Resource policy JSON validation failed: %v", err)
			result["PolicyValidationPassed"] = false
			result["ValidationErrors"] = []interface{}{
				map[string]interface{}{
					"CheckName":    "ResourcePolicySyntax",
					"ErrorMessage": "Invalid JSON syntax in resource policy document",
				},
			}
		}
	}

	return result, nil
}
