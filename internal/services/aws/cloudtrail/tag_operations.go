package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// AddTags adds or overwrites tags on a CloudTrail trail.
func (s *CloudTrailService) AddTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return tagutil.HandleTag(ctx, req, cloudTrailTagConfig(store, s.mapStoreError, true))
}

// RemoveTags removes the specified tags from a CloudTrail trail.
func (s *CloudTrailService) RemoveTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return tagutil.HandleUntag(ctx, req, cloudTrailTagConfig(store, s.mapStoreError, false))
}

// ListTags lists all tags assigned to a CloudTrail trail.
func (s *CloudTrailService) ListTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return tagutil.HandleList(ctx, req, cloudTrailTagConfig(store, s.mapStoreError, false))
}
