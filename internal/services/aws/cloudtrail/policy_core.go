package cloudtrail

import (
	"encoding/json"
	"strings"

	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// ResourcePolicyInput carries the resource ARN for the resource-policy
// operations.
type ResourcePolicyInput struct {
	ResourceARN string
}

// PutResourcePolicyInput carries the resource ARN and the policy document.
type PutResourcePolicyInput struct {
	ResourceARN string
	Policy      string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getResourcePolicyCore is the single entry point for GetResourcePolicy.
func (s *CloudTrailService) getResourcePolicyCore(store cloudtrailstore.CloudTrailStoreInterface, in ResourcePolicyInput) (map[string]interface{}, error) {
	if in.ResourceARN == "" {
		return nil, ErrInvalidParameter
	}

	// The operation is documented for trails, event data stores, and
	// channels; dispatch on the ARN resource type exactly like the Put and
	// Delete paths so the Get/Put round trip works for every type.
	if err := s.verifyPolicyResource(store, in.ResourceARN); err != nil {
		return nil, err
	}

	policy, err := store.GetResourcePolicy(in.ResourceARN)
	if err != nil {
		return nil, ErrResourcePolicyNotFound
	}

	return map[string]interface{}{
		"ResourceArn":    policy.ResourceARN,
		"ResourcePolicy": policy.Policy,
	}, nil
}

// putResourcePolicyCore is the single entry point for PutResourcePolicy.
func (s *CloudTrailService) putResourcePolicyCore(store cloudtrailstore.CloudTrailStoreInterface, in PutResourcePolicyInput) (map[string]interface{}, error) {
	if in.ResourceARN == "" {
		return nil, ErrInvalidParameter
	}

	if err := s.verifyPolicyResource(store, in.ResourceARN); err != nil {
		return nil, err
	}

	if err := validateResourcePolicyDocument(in.Policy); err != nil {
		return nil, err
	}

	if err := store.PutResourcePolicy(in.ResourceARN, in.Policy); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ResourceArn": in.ResourceARN,
	}, nil
}

// deleteResourcePolicyCore is the single entry point for
// DeleteResourcePolicy.
func (s *CloudTrailService) deleteResourcePolicyCore(store cloudtrailstore.CloudTrailStoreInterface, in ResourcePolicyInput) error {
	if in.ResourceARN == "" {
		return ErrInvalidParameter
	}

	// Verify the resource exists before deleting its policy.
	if err := s.verifyPolicyResource(store, in.ResourceARN); err != nil {
		return err
	}

	// Verify the policy exists before deleting.
	if _, err := store.GetResourcePolicy(in.ResourceARN); err != nil {
		return ErrResourcePolicyNotFound
	}

	if err := store.DeleteResourcePolicy(in.ResourceARN); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// verifyPolicyResource confirms that the resource identified by resourceARN
// exists before a resource-policy operation, dispatching on the ARN resource
// field. The documented ARN set is "the CloudTrail event data store,
// dashboard, or channel" (GetResourcePolicy/PutResourcePolicy/
// DeleteResourcePolicy ResourceArn): event data store and channel ARNs
// resolve through the store, a dashboard ARN addresses a supported type the
// platform does not serve yet and maps to the declared
// ResourceNotFoundException, and every other ARN shape — trails included —
// is rejected with ResourceTypeNotSupportedException.
func (s *CloudTrailService) verifyPolicyResource(store cloudtrailstore.CloudTrailStoreInterface, resourceARN string) error {
	parsed, err := svcarn.ParseARN(resourceARN)
	if err != nil {
		return newResourceARNNotValidException(
			"The resource ARN is not valid")
	}
	switch {
	case strings.HasPrefix(parsed.Resource, "eventdatastore/"):
		if _, err := store.GetEventDataStore(resourceARN); err != nil {
			return s.mapStoreError(err)
		}
	case strings.HasPrefix(parsed.Resource, "channel/"):
		if _, err := store.GetChannel(resourceARN); err != nil {
			return s.mapStoreError(err)
		}
	case strings.HasPrefix(parsed.Resource, "dashboard/"):
		return newResourceNotFoundException(
			"The specified resource was not found")
	default:
		return newResourceTypeNotSupportedException(
			"The specified resource type is not supported by CloudTrail")
	}
	return nil
}

// validateResourcePolicyDocument checks a resource-based policy document
// before it is stored: the document must be a JSON object whose Statement
// entries each carry a Principal member. The operation declares
// ResourcePolicyNotValidException for "syntax errors, or contains a
// principal that is not valid".
func validateResourcePolicyDocument(policy string) error {
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(policy), &doc); err != nil {
		return newResourcePolicyNotValidException(
			"The resource-based policy has syntax errors")
	}
	stmtRaw, ok := doc["Statement"]
	if !ok {
		return newResourcePolicyNotValidException(
			"The resource-based policy must contain a Statement member")
	}
	statements, ok := stmtRaw.([]interface{})
	if !ok {
		// A single statement object is the documented shorthand for a
		// one-element list.
		if _, isMap := stmtRaw.(map[string]interface{}); !isMap {
			return newResourcePolicyNotValidException(
				"The resource-based policy Statement must be a list or an object")
		}
		statements = []interface{}{stmtRaw}
	}
	for _, stmt := range statements {
		m, ok := stmt.(map[string]interface{})
		if !ok {
			return newResourcePolicyNotValidException(
				"Each policy statement must be an object")
		}
		if _, hasPrincipal := m["Principal"]; !hasPrincipal {
			if _, hasNotPrincipal := m["NotPrincipal"]; !hasNotPrincipal {
				return newResourcePolicyNotValidException(
					"Each policy statement must contain a Principal")
			}
		}
	}
	return nil
}
