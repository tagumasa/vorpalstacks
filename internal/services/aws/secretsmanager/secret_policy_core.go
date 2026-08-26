package secretsmanager

import (
	"context"
	"net/http"

	"vorpalstacks/internal/common/errors"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// ---------------------------------------------------------------------------
// Input structs — transport-agnostic DTOs shared by HTTP API and admin handler
// ---------------------------------------------------------------------------

// GetResourcePolicyInput carries all fields needed for GetResourcePolicy.
type GetResourcePolicyInput struct {
	SecretId string
}

// PutResourcePolicyInput carries all fields needed for PutResourcePolicy.
type PutResourcePolicyInput struct {
	SecretId          string
	ResourcePolicy    string
	BlockPublicPolicy bool
}

// DeleteResourcePolicyInput carries all fields needed for DeleteResourcePolicy.
type DeleteResourcePolicyInput struct {
	SecretId string
}

// ValidateResourcePolicyInput carries all fields needed for
// ValidateResourcePolicy.
type ValidateResourcePolicyInput struct {
	SecretId       string
	ResourcePolicy string
}

// ---------------------------------------------------------------------------
// Result structs — transport-agnostic results
// ---------------------------------------------------------------------------

// GetResourcePolicyResult holds the transport-agnostic result of
// GetResourcePolicy.
type GetResourcePolicyResult struct {
	ARN            string
	Name           string
	ResourcePolicy string
}

// PolicyARNNameResult is the shared ARN/Name-only result shape of
// PutResourcePolicy and DeleteResourcePolicy.
type PolicyARNNameResult struct {
	ARN  string
	Name string
}

// ValidateResourcePolicyResult holds the transport-agnostic result of
// ValidateResourcePolicy.
type ValidateResourcePolicyResult struct {
	PolicyValidationPassed bool
	Checks                 []policyCheck
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getResourcePolicyCore is the single entry point for reading a secret's
// resource policy.
func (s *SecretsManagerService) getResourcePolicyCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in GetResourcePolicyInput) (*GetResourcePolicyResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecretForMetadata(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	policy, err := store.GetResourcePolicy(secret.Name)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &GetResourcePolicyResult{
		ARN:            secret.ARN,
		Name:           secret.Name,
		ResourcePolicy: policy,
	}, nil
}

// putResourcePolicyCore is the single entry point for attaching a resource
// policy to a secret.
func (s *SecretsManagerService) putResourcePolicyCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in PutResourcePolicyInput) (*PolicyARNNameResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	if in.ResourcePolicy == "" {
		return nil, errors.ErrMissingParameter
	}
	if err := validateResourcePolicyLength(in.ResourcePolicy); err != nil {
		return nil, err
	}

	// Resolve the secret first so that non-existent secrets return
	// ResourceNotFoundException before any policy validation.
	secret, err := resolveSecretForMetadata(store, in.SecretId)
	if err != nil {
		return nil, err
	}
	if err := checkNotDeleted(secret); err != nil {
		return nil, err
	}

	// Parse and structurally validate the policy JSON.  Malformed
	// JSON returns MalformedPolicyDocumentException.
	doc, err := ensurePolicyJSONValid(in.ResourcePolicy)
	if err != nil {
		return nil, errors.NewAWSError("MalformedPolicyDocumentException",
			err.Error(), http.StatusBadRequest)
	}

	// BlockPublicPolicy: when true, Secrets Manager rejects resource policies
	// that grant broad access (Principal "*" without a restricting Condition).
	// By default (false / unset) public policies are accepted.
	if in.BlockPublicPolicy && isPolicyPublic(doc) {
		return nil, errors.NewAWSError("PublicPolicyException",
			"The BlockPublicPolicy parameter is set to true, and the resource policy did not prevent broad access to the secret.",
			http.StatusBadRequest)
	}

	if err := store.PutResourcePolicy(secret.Name, in.ResourcePolicy); err != nil {
		return nil, mapStoreError(err)
	}

	return &PolicyARNNameResult{
		ARN:  secret.ARN,
		Name: secret.Name,
	}, nil
}

// deleteResourcePolicyCore is the single entry point for deleting the
// resource policy from a secret.
func (s *SecretsManagerService) deleteResourcePolicyCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in DeleteResourcePolicyInput) (*PolicyARNNameResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecretForMetadata(store, in.SecretId)
	if err != nil {
		return nil, err
	}
	if err := checkNotDeleted(secret); err != nil {
		return nil, err
	}

	if err := store.DeleteResourcePolicy(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	return &PolicyARNNameResult{
		ARN:  secret.ARN,
		Name: secret.Name,
	}, nil
}

// validateResourcePolicyCore is the single entry point for validating a
// resource policy. Runs multiple validation checks (syntax, missing version,
// public access) and returns all failures so a single call may report
// multiple issues, matching the AWS behaviour.
func (s *SecretsManagerService) validateResourcePolicyCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in ValidateResourcePolicyInput) (*ValidateResourcePolicyResult, error) {
	if in.ResourcePolicy == "" {
		return nil, errors.ErrMissingParameter
	}
	// The identical NonEmptyResourcePolicyType @length(1,20480) constraint
	// governs both PutResourcePolicy and ValidateResourcePolicy; a policy
	// PutResourcePolicy rejects must not pass validation-only either.
	if err := validateResourcePolicyLength(in.ResourcePolicy); err != nil {
		return nil, err
	}

	if in.SecretId != "" {
		if err := validateSecretId(in.SecretId); err != nil {
			return nil, err
		}
		_, err := resolveSecretForMetadata(store, in.SecretId)
		if err != nil {
			return nil, err
		}
	}

	checks := validatePolicyDocument(in.ResourcePolicy)

	return &ValidateResourcePolicyResult{
		PolicyValidationPassed: len(checks) == 0,
		Checks:                 checks,
	}, nil
}
