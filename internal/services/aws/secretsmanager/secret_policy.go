package secretsmanager

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// GetResourcePolicy returns the resource policy for a secret.
func (s *SecretsManagerService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getResourcePolicyCore(ctx, store, GetResourcePolicyInput{
		SecretId: request.GetStringParam(req.Parameters, "SecretId"),
	})
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}
	if result.ResourcePolicy != "" {
		response["ResourcePolicy"] = result.ResourcePolicy
	}

	return response, nil
}

// PutResourcePolicy attaches a resource policy to a secret.
func (s *SecretsManagerService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putResourcePolicyCore(ctx, store, PutResourcePolicyInput{
		SecretId:          request.GetStringParam(req.Parameters, "SecretId"),
		ResourcePolicy:    request.GetStringParam(req.Parameters, "ResourcePolicy"),
		BlockPublicPolicy: request.GetBoolParam(req.Parameters, "BlockPublicPolicy"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}, nil
}

// DeleteResourcePolicy deletes the resource policy from a secret.
func (s *SecretsManagerService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.deleteResourcePolicyCore(ctx, store, DeleteResourcePolicyInput{
		SecretId: request.GetStringParam(req.Parameters, "SecretId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}, nil
}

// ValidateResourcePolicy validates a resource policy for a secret.
// Runs multiple validation checks (syntax, missing version, public access)
// and returns all failures in ValidationErrors, matching the AWS behaviour
// where a single call may report multiple issues.
func (s *SecretsManagerService) ValidateResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.validateResourcePolicyCore(ctx, store, ValidateResourcePolicyInput{
		SecretId:       request.GetStringParam(req.Parameters, "SecretId"),
		ResourcePolicy: request.GetStringParam(req.Parameters, "ResourcePolicy"),
	})
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"PolicyValidationPassed": result.PolicyValidationPassed,
	}
	if !result.PolicyValidationPassed {
		response["ValidationErrors"] = policyChecksToResponse(result.Checks)
	}

	return response, nil
}
