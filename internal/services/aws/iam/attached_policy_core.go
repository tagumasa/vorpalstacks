// Transport-agnostic Core functions for IAM managed-policy attachment:
// validation and store operations shared by the AWS-compatible HTTP API
// handlers and any admin surface (the xxxCore pattern).
package iam

import (
	"errors"
	"strings"

	"vorpalstacks/internal/common/pagination"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// AttachPolicyInput holds the parameters for attaching or detaching a
// managed policy on a principal (user, group, or role).
type AttachPolicyInput struct {
	PrincipalType string
	PrincipalName string
	PolicyArn     string
}

// AttachedPolicyEntry is one managed-policy entry of a principal's
// attachment list; Path participates in prefix filtering and is not part of
// the wire response.
type AttachedPolicyEntry struct {
	PolicyName string
	PolicyArn  string
	Path       string
}

// AttachedPolicyListResult is the Core result for the list-attached-policies
// family.
type AttachedPolicyListResult struct {
	Entries     []AttachedPolicyEntry
	IsTruncated bool
	NextMarker  string
}

// attachPolicyCore validates input and attaches a managed policy to a
// principal. Attaching an already-attached policy is an idempotent success.
func (s *IAMService) attachPolicyCore(store *iamstore.IAMStore, input *AttachPolicyInput) error {
	if input.PrincipalName == "" {
		return principalNameRequiredError(input.PrincipalType)
	}
	if input.PolicyArn == "" {
		return NewValidationError("PolicyArn")
	}
	if err := validateAttachPolicyArn(input.PolicyArn); err != nil {
		return err
	}

	if !principalExists(store, input.PrincipalType, input.PrincipalName) {
		return newPrincipalNotFoundError(input.PrincipalType, input.PrincipalName)
	}
	if !store.Policies().Exists(input.PolicyArn) {
		return NewNoSuchPolicyError(input.PolicyArn)
	}

	if store.AttachedPolicies().IsAttached(input.PrincipalType, input.PrincipalName, input.PolicyArn) {
		return nil
	}

	if err := store.AttachedPolicies().Attach(input.PrincipalType, input.PrincipalName, input.PolicyArn); err != nil {
		return err
	}
	if err := store.Policies().IncrementAttachmentCount(input.PolicyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Detach(input.PrincipalType, input.PrincipalName, input.PolicyArn); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	return nil
}

// detachPolicyCore validates input and detaches a managed policy from a
// principal.
func (s *IAMService) detachPolicyCore(store *iamstore.IAMStore, input *AttachPolicyInput) error {
	if input.PrincipalName == "" {
		return principalNameRequiredError(input.PrincipalType)
	}
	if input.PolicyArn == "" {
		return NewValidationError("PolicyArn")
	}
	if err := validateAttachPolicyArn(input.PolicyArn); err != nil {
		return err
	}

	if !store.AttachedPolicies().IsAttached(input.PrincipalType, input.PrincipalName, input.PolicyArn) {
		return NewPolicyNotAttachedError(input.PolicyArn)
	}

	if err := store.AttachedPolicies().Detach(input.PrincipalType, input.PrincipalName, input.PolicyArn); err != nil {
		return err
	}
	if err := store.Policies().DecrementAttachmentCount(input.PolicyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Attach(input.PrincipalType, input.PrincipalName, input.PolicyArn); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	return nil
}

// listAttachedPoliciesCore lists the managed policies attached to a
// principal, filtered by path prefix and paginated by Marker and MaxItems.
func (s *IAMService) listAttachedPoliciesCore(store *iamstore.IAMStore, principalType, principalName, pathPrefix, marker string, maxItems int) (*AttachedPolicyListResult, error) {
	if principalName == "" {
		return nil, principalNameRequiredError(principalType)
	}

	if !principalExists(store, principalType, principalName) {
		return nil, newPrincipalNotFoundError(principalType, principalName)
	}

	policyArns, err := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	if err != nil {
		return nil, err
	}

	entries := make([]AttachedPolicyEntry, 0, len(policyArns))
	for _, arn := range policyArns {
		policy, err := store.Policies().Get(arn)
		if err != nil {
			continue
		}
		if pathPrefix != "" && !strings.HasPrefix(policy.Path, pathPrefix) {
			continue
		}
		entries = append(entries, AttachedPolicyEntry{
			PolicyName: policy.PolicyName,
			PolicyArn:  policy.Arn,
			Path:       policy.Path,
		})
	}

	result := pagination.PaginateSlice(entries, marker, maxItems, func(e AttachedPolicyEntry) string {
		return e.PolicyArn
	})

	return &AttachedPolicyListResult{
		Entries:     result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
