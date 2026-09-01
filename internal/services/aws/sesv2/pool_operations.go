package sesv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

// CreateDedicatedIpPool creates a new dedicated IP pool.
func (s *SESv2Service) CreateDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.createDedicatedIpPoolCore(store, CreateDedicatedIpPoolInput{
		PoolName:    request.GetStringParam(req.Parameters, "PoolName"),
		ScalingMode: request.GetStringParam(req.Parameters, "ScalingMode"),
		Tags:        tags.ParseTags(req.Parameters, "Tags"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteDedicatedIpPool deletes a dedicated IP pool.
func (s *SESv2Service) DeleteDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteDedicatedIpPoolCore(store, request.GetStringParam(req.Parameters, "PoolName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetDedicatedIpPool retrieves the details of a dedicated IP pool.
func (s *SESv2Service) GetDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getDedicatedIpPoolCore(store, request.GetStringParam(req.Parameters, "PoolName"))
}

// ListDedicatedIpPools returns a list of dedicated IP pools.
func (s *SESv2Service) ListDedicatedIpPools(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listDedicatedIpPoolsCore(store,
		pagination.GetMaxItems(req.Parameters, 100, "PageSize"),
		pagination.GetMarker(req.Parameters, "NextToken"))
}
