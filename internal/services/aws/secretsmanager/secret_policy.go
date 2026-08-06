package secretsmanager

import (
	"context"
	"net/http"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
)

// GetResourcePolicy returns the resource policy for a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_GetResourcePolicy.html
func (s *SecretsManagerService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	policyStr := request.GetStringParam(req.Parameters, "ResourcePolicy")
	if policyStr == "" {
		return nil, errors.ErrMissingParameter
	}
	if err := validateResourcePolicyLength(policyStr); err != nil {
		return nil, err
	}

	// Resolve the secret first so that non-existent secrets return
	// ResourceNotFoundException before any policy validation.
	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}
	if err := checkNotDeleted(secret); err != nil {
		return nil, err
	}

	// Parse and structurally validate the policy JSON.  Malformed
	// JSON returns MalformedPolicyDocumentException.
	doc, err := ensurePolicyJSONValid(policyStr)
	if err != nil {
		return nil, errors.NewAWSError("MalformedPolicyDocumentException",
			err.Error(), http.StatusBadRequest)
	}

	// BlockPublicPolicy: when true, Secrets Manager rejects resource policies
	// that grant broad access (Principal "*" without a restricting Condition).
	// By default (false / unset) public policies are accepted.
	blockPublicPolicy := request.GetBoolParam(req.Parameters, "BlockPublicPolicy")
	if blockPublicPolicy && isPolicyPublic(doc) {
		return nil, errors.NewAWSError("PublicPolicyException",
			"The BlockPublicPolicy parameter is set to true, and the resource policy did not prevent broad access to the secret.",
			http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.PutResourcePolicy(secret.Name, policyStr); err != nil {
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

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}
	if err := checkNotDeleted(secret); err != nil {
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
// Runs multiple validation checks (syntax, missing version, public access)
// and returns all failures in ValidationErrors, matching the AWS behaviour
// where a single call may report multiple issues.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_ValidateResourcePolicy.html
func (s *SecretsManagerService) ValidateResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	policyStr := request.GetStringParam(req.Parameters, "ResourcePolicy")

	if policyStr == "" {
		return nil, errors.ErrMissingParameter
	}

	if secretId != "" {
		_, err := s.resolveSecretForMetadata(reqCtx, secretId)
		if err != nil {
			return nil, err
		}
	}

	checks := validatePolicyDocument(policyStr)

	result := map[string]interface{}{}
	if len(checks) == 0 {
		result["PolicyValidationPassed"] = true
	} else {
		result["PolicyValidationPassed"] = false
		result["ValidationErrors"] = policyChecksToResponse(checks)
	}

	return result, nil
}
