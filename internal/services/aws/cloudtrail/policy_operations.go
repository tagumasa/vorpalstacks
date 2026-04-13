package cloudtrail

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetResourcePolicy retrieves the resource policy for a CloudTrail resource.
func (s *CloudTrailService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := getParam(req, "ResourceArn")

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
		return map[string]interface{}{
			"ResourceArn":    resourceARN,
			"ResourcePolicy": "",
		}, nil
	}

	return map[string]interface{}{
		"ResourceArn":    policy.ResourceARN,
		"ResourcePolicy": policy.Policy,
	}, nil
}

// PutResourcePolicy attaches a resource policy to a CloudTrail resource.
func (s *CloudTrailService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := getParam(req, "ResourceArn")
	policy := getParam(req, "ResourcePolicy")

	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if strings.Contains(resourceARN, ":trail/") {
		if _, err := store.GetTrailByARN(resourceARN); err != nil {
			return nil, s.mapStoreError(err)
		}
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
	resourceARN := getParam(req, "ResourceArn")

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
