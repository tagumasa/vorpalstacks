package cloudtrail

import (
	"context"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// GetResourcePolicy retrieves the resource policy for a CloudTrail resource.
func (s *CloudTrailService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := req.GetParam("ResourceArn")

	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	_, err = store.GetTrailByARN(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	policy, err := store.GetResourcePolicy(resourceARN)
	if err != nil {
		return nil, awserrors.NewAWSError("ResourcePolicyNotFoundException",
			"Resource policy not found", 404)
	}

	return map[string]interface{}{
		"ResourceArn":    policy.ResourceARN,
		"ResourcePolicy": policy.Policy,
	}, nil
}

// PutResourcePolicy attaches a resource policy to a CloudTrail resource.
func (s *CloudTrailService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := req.GetParam("ResourceArn")
	policy := req.GetParam("ResourcePolicy")

	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := s.verifyPolicyResource(store, resourceARN); err != nil {
		return nil, err
	}

	if err := store.PutResourcePolicy(resourceARN, policy); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ResourceArn": resourceARN,
	}, nil
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

// DeleteResourcePolicy removes the resource policy from a CloudTrail resource.
func (s *CloudTrailService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := req.GetParam("ResourceArn")

	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// Verify the resource exists before deleting its policy.
	if err := s.verifyPolicyResource(store, resourceARN); err != nil {
		return nil, err
	}

	// Verify the policy exists before deleting.
	if _, err := store.GetResourcePolicy(resourceARN); err != nil {
		return nil, awserrors.NewAWSError("ResourcePolicyNotFoundException",
			"Resource policy not found", 404)
	}

	if err := store.DeleteResourcePolicy(resourceARN); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}
