// Transport-agnostic Core functions for IAM inline policies: validation and
// store operations shared by the AWS-compatible HTTP API handlers and any
// admin surface (the xxxCore pattern).
package iam

import (
	"vorpalstacks/internal/common/pagination"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// PutInlinePolicyInput holds the parameters for putting an inline policy on
// a principal (user, group, or role).
type PutInlinePolicyInput struct {
	PrincipalType  string
	PrincipalName  string
	PolicyName     string
	PolicyDocument string
}

// principalNameRequiredError maps the principal type to the error returned
// when the principal name parameter is missing.
func principalNameRequiredError(principalType string) error {
	switch principalType {
	case PrincipalTypeUser:
		return ErrNoSuchUser
	case PrincipalTypeGroup:
		return ErrNoSuchGroup
	case PrincipalTypeRole:
		return ErrNoSuchRole
	}
	return NewValidationError("PrincipalName")
}

// principalExists reports whether a principal of the given type exists.
func principalExists(store *iamstore.IAMStore, principalType, name string) bool {
	switch principalType {
	case PrincipalTypeUser:
		return store.Users().Exists(name)
	case PrincipalTypeGroup:
		return store.Groups().Exists(name)
	case PrincipalTypeRole:
		return store.Roles().Exists(name)
	}
	return false
}

// newPrincipalNotFoundError maps the principal type to its NoSuchEntity
// error.
func newPrincipalNotFoundError(principalType, name string) error {
	switch principalType {
	case PrincipalTypeUser:
		return NewNoSuchUserError(name)
	case PrincipalTypeGroup:
		return NewNoSuchGroupError(name)
	case PrincipalTypeRole:
		return NewNoSuchRoleError(name)
	}
	return NewNoSuchEntityError("Principal", name)
}

// putInlinePolicyCore validates input and creates or updates an inline
// policy for a principal.
func (s *IAMService) putInlinePolicyCore(store *iamstore.IAMStore, input *PutInlinePolicyInput) error {
	if input.PrincipalName == "" {
		return principalNameRequiredError(input.PrincipalType)
	}
	if input.PolicyName == "" {
		return NewValidationError("PolicyName")
	}
	if err := validateEntityName128(input.PolicyName, "PolicyName"); err != nil {
		return err
	}
	if !validatePolicyDocument(input.PolicyDocument) {
		return ErrMalformedPolicyDocument
	}

	if !principalExists(store, input.PrincipalType, input.PrincipalName) {
		return newPrincipalNotFoundError(input.PrincipalType, input.PrincipalName)
	}

	return store.InlinePolicies().Put(input.PrincipalType, input.PrincipalName, input.PolicyName, input.PolicyDocument)
}

// getInlinePolicyCore retrieves an inline policy for a principal.
func (s *IAMService) getInlinePolicyCore(store *iamstore.IAMStore, principalType, principalName, policyName string) (*iamstore.InlinePolicy, error) {
	if principalName == "" {
		return nil, principalNameRequiredError(principalType)
	}
	if policyName == "" {
		return nil, NewValidationError("PolicyName")
	}

	if !principalExists(store, principalType, principalName) {
		return nil, newPrincipalNotFoundError(principalType, principalName)
	}

	policy, err := store.InlinePolicies().Get(principalType, principalName, policyName)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyName)
	}
	return policy, nil
}

// deleteInlinePolicyCore deletes an inline policy from a principal.
func (s *IAMService) deleteInlinePolicyCore(store *iamstore.IAMStore, principalType, principalName, policyName string) error {
	if principalName == "" {
		return principalNameRequiredError(principalType)
	}
	if policyName == "" {
		return NewValidationError("PolicyName")
	}

	if !principalExists(store, principalType, principalName) {
		return newPrincipalNotFoundError(principalType, principalName)
	}
	if !store.InlinePolicies().Exists(principalType, principalName, policyName) {
		return NewNoSuchPolicyError(policyName)
	}

	return store.InlinePolicies().Delete(principalType, principalName, policyName)
}

// listInlinePoliciesCore lists the inline policy names for a principal,
// paginated by Marker and MaxItems.
func (s *IAMService) listInlinePoliciesCore(store *iamstore.IAMStore, principalType, principalName, marker string, maxItems int) (*pagination.SliceResult[string], error) {
	if principalName == "" {
		return nil, principalNameRequiredError(principalType)
	}

	if !principalExists(store, principalType, principalName) {
		return nil, newPrincipalNotFoundError(principalType, principalName)
	}

	allNames, err := store.InlinePolicies().List(principalType, principalName)
	if err != nil {
		return nil, err
	}

	result := pagination.PaginateSlice(allNames, marker, maxItems, func(name string) string {
		return name
	})
	return &result, nil
}
