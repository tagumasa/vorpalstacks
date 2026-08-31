package kms

// policy_core.go carries the Core functions of the KMS key-policy family.
// AWS supports only the "default" policy name on customer managed keys;
// the Cores enforce that together with policy size, JSON well-formedness
// and the root-lockout rule before persisting.

import (
	"encoding/json"

	"vorpalstacks/internal/common/pagination"

	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// getKeyPolicyCore is the single entry point for reading a key policy.
func (s *KMSService) getKeyPolicyCore(stores *kmsStores, key *kmsstore.Key, policyName string) (*kmsstore.KeyPolicy, error) {
	if policyName == "" {
		policyName = "default"
	}
	// AWS: only "default" is a valid policy name for KMS keys.
	if policyName != "default" {
		return nil, ErrValidation
	}

	return stores.keyPolicies.Get(key.KeyID, policyName)
}

// putKeyPolicyCore is the single entry point for writing a key policy.
// The BypassPolicyLockoutSafetyCheck parameter, when false (default),
// enforces that the policy must not lock out the root principal.
func (s *KMSService) putKeyPolicyCore(stores *kmsStores, key *kmsstore.Key, policyName, policyDocument string, bypassPolicyLockoutSafetyCheck bool, accountID string) error {
	if policyName == "" {
		policyName = "default"
	}
	if policyName != "default" {
		return ErrValidation
	}

	if policyDocument == "" {
		return ErrMalformedPolicy
	}
	if err := validatePolicySize(policyDocument); err != nil {
		return err
	}

	var js interface{}
	if err := json.Unmarshal([]byte(policyDocument), &js); err != nil {
		return ErrMalformedPolicy
	}

	if !bypassPolicyLockoutSafetyCheck {
		if err := validatePolicyDoesNotLockOutRoot(policyDocument, accountID); err != nil {
			return err
		}
	}

	return stores.keyPolicies.Put(key.KeyID, policyName, policyDocument)
}

// PolicyListResult is the paginated Core result of ListKeyPolicies.
type PolicyListResult struct {
	Items       []string
	IsTruncated bool
	NextMarker  string
}

// listKeyPoliciesCore is the single entry point for listing the policy
// names attached to a key. The marker validation runs after the store
// list, matching the original failure precedence.
func (s *KMSService) listKeyPoliciesCore(stores *kmsStores, key *kmsstore.Key, marker string, maxItems int) (*PolicyListResult, error) {
	policies, err := stores.keyPolicies.List(key.KeyID)
	if err != nil {
		return nil, err
	}

	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}

	result := pagination.PaginateSlice(policies, marker, maxItems, func(p string) string {
		return p
	})
	return &PolicyListResult{
		Items:       result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
