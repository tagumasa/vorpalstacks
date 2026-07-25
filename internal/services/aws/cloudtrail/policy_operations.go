package cloudtrail

import (
	"context"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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

	switch {
	case strings.Contains(resourceARN, ":trail/"):
		if _, err := store.GetTrailByARN(resourceARN); err != nil {
			return nil, s.mapStoreError(err)
		}
	case strings.Contains(resourceARN, ":eventdata-store/"):
		if _, err := store.GetEventDataStore(resourceARN); err != nil {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException",
				"Event data store not found", 404)
		}
	case strings.Contains(resourceARN, ":channel/"):
		if _, err := store.GetChannel(resourceARN); err != nil {
			return nil, awserrors.NewAWSError("ChannelNotFoundException",
				"Channel not found", 404)
		}
	default:
		return nil, awserrors.NewAWSError("ResourceARNNotValidException",
			"The resource ARN is not valid", 400)
	}

	if err := store.PutResourcePolicy(resourceARN, policy); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ResourceArn": resourceARN,
	}, nil
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

	if err := store.DeleteResourcePolicy(resourceARN); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}
