// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"

	"vorpalstacks/internal/services/aws/iam/policy"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// IAMValidator validates IAM policies and role trust relationships.
type IAMValidator struct {
	roleProvider RolePolicyProvider
	policyCache  *PolicyCache
	accountID    string
}

// NewIAMValidator creates a new IAM validator.
func NewIAMValidator(roleProvider RolePolicyProvider, accountID string) *IAMValidator {
	return &IAMValidator{
		roleProvider: roleProvider,
		accountID:    accountID,
	}
}

// NewIAMValidatorWithCache creates a new IAM validator with a policy cache.
func NewIAMValidatorWithCache(roleProvider RolePolicyProvider, accountID string, cache *PolicyCache) *IAMValidator {
	return &IAMValidator{
		roleProvider: roleProvider,
		accountID:    accountID,
		policyCache:  cache,
	}
}

// RoleErrorFactories defines error factory functions for role validation errors.
type RoleErrorFactories struct {
	RoleNotFoundError        func(roleArn string) error
	RoleCannotBeAssumedError func(roleArn string) error
	InvalidArnError          func(roleArn string) error
}

// ValidateRoleForService validates if a role can be assumed by a service principal.
func (v *IAMValidator) ValidateRoleForService(ctx context.Context, roleArn string, servicePrincipal ServicePrincipal) error {
	return v.ValidateRoleForServiceWithErrors(ctx, roleArn, servicePrincipal, &RoleErrorFactories{
		RoleNotFoundError:        NewLambdaRoleNotFoundError,
		RoleCannotBeAssumedError: NewLambdaRoleCannotBeAssumedError,
		InvalidArnError:          NewInvalidRoleArnError,
	})
}

// ValidateRoleForServiceWithErrors validates if a role can be assumed by a service principal with custom error factories.
func (v *IAMValidator) ValidateRoleForServiceWithErrors(ctx context.Context, roleArn string, servicePrincipal ServicePrincipal, errFactory *RoleErrorFactories) error {
	parsed, err := arnutil.ParseARN(roleArn)
	if err != nil {
		return errFactory.InvalidArnError(roleArn)
	}

	if !arnutil.IsValidRoleARN(roleArn) {
		return errFactory.InvalidArnError(roleArn)
	}

	if parsed.IsCrossAccount(v.accountID) {
		return nil
	}

	policyDoc, err := v.roleProvider.GetAssumeRolePolicyDocument(parsed.RoleName)
	if err != nil {
		return errFactory.RoleNotFoundError(roleArn)
	}

	evalCtx := BuildEvaluationContext(v.accountID, string(servicePrincipal))
	if err := v.ValidateTrustPolicy(parsed.RoleName, policyDoc, servicePrincipal, evalCtx); err != nil {
		return errFactory.RoleCannotBeAssumedError(roleArn)
	}
	return nil
}

// ValidateRoleExists checks if a role exists.
func (v *IAMValidator) ValidateRoleExists(roleArn string) error {
	parsed, err := arnutil.ParseARN(roleArn)
	if err != nil {
		return NewInvalidRoleArnError(roleArn)
	}

	if !arnutil.IsValidRoleARN(roleArn) {
		return NewInvalidRoleArnError(roleArn)
	}

	if parsed.IsCrossAccount(v.accountID) {
		return nil
	}

	if !v.roleProvider.RoleExists(parsed.RoleName) {
		return NewLambdaRoleNotFoundError(roleArn)
	}

	return nil
}

// ValidateTrustPolicy validates the trust policy of a role for a service principal.
func (v *IAMValidator) ValidateTrustPolicy(roleName, policyDoc string, servicePrincipal ServicePrincipal, evalCtx *policy.EvaluationContext) error {
	doc, err := v.getOrParsePolicy(roleName, policyDoc)
	if err != nil {
		return NewLambdaRoleCannotBeAssumedError("")
	}

	return EvaluateTrustPolicy(doc, string(servicePrincipal), evalCtx)
}

// InvalidateCache invalidates the policy cache for a role.
func (v *IAMValidator) InvalidateCache(roleName string) {
	if v.policyCache != nil {
		v.policyCache.Delete(roleName)
	}
}

func (v *IAMValidator) getOrParsePolicy(roleName, policyDoc string) (*policy.Document, error) {
	if v.policyCache != nil {
		if doc, ok := v.policyCache.Get(roleName); ok {
			return doc, nil
		}
	}

	doc, err := policy.ParseDocument(policyDoc)
	if err != nil {
		return nil, err
	}

	if v.policyCache != nil {
		v.policyCache.SetWithTTL(roleName, doc, DefaultCacheTTL)
	}

	return doc, nil
}
