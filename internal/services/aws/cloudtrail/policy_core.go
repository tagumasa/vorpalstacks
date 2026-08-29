package cloudtrail

import (
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
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
		return nil, awserrors.NewAWSError("ResourcePolicyNotFoundException",
			"Resource policy not found", 404)
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
		return awserrors.NewAWSError("ResourcePolicyNotFoundException",
			"Resource policy not found", 404)
	}

	if err := store.DeleteResourcePolicy(in.ResourceARN); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// verifyPolicyResource confirms that the resource identified by resourceARN
// exists before a resource-policy operation, dispatching on the ARN resource
// field (trail/, eventdata-store/, channel/) to the matching store getter.
// Any other ARN shape is rejected with ResourceARNNotValidException.
func (s *CloudTrailService) verifyPolicyResource(store cloudtrailstore.CloudTrailStoreInterface, resourceARN string) error {
	_, _, _, _, resource := svcarn.SplitARN(resourceARN)
	switch {
	case strings.HasPrefix(resource, "trail/"):
		if _, err := store.GetTrailByARN(resourceARN); err != nil {
			return s.mapStoreError(err)
		}
	case strings.HasPrefix(resource, "eventdata-store/"):
		if _, err := store.GetEventDataStore(resourceARN); err != nil {
			return awserrors.NewAWSError("EventDataStoreNotFoundException",
				"Event data store not found", 404)
		}
	case strings.HasPrefix(resource, "channel/"):
		if _, err := store.GetChannel(resourceARN); err != nil {
			return awserrors.NewAWSError("ChannelNotFoundException",
				"Channel not found", 404)
		}
	default:
		return awserrors.NewAWSError("ResourceARNNotValidException",
			"The resource ARN is not valid", 400)
	}
	return nil
}
