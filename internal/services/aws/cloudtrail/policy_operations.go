package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetResourcePolicy retrieves the resource policy for a CloudTrail resource.
func (s *CloudTrailService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.getResourcePolicyCore(store, ResourcePolicyInput{
		ResourceARN: req.GetParam("ResourceArn"),
	})
}

// PutResourcePolicy attaches a resource policy to a CloudTrail resource.
func (s *CloudTrailService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.putResourcePolicyCore(store, PutResourcePolicyInput{
		ResourceARN: req.GetParam("ResourceArn"),
		Policy:      req.GetParam("ResourcePolicy"),
	})
}

// DeleteResourcePolicy removes the resource policy from a CloudTrail resource.
func (s *CloudTrailService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := s.deleteResourcePolicyCore(store, ResourcePolicyInput{
		ResourceARN: req.GetParam("ResourceArn"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
