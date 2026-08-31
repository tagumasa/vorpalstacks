package ssm

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// AddTagsToResource adds or overwrites tags on an SSM parameter.
func (s *SSMService) AddTagsToResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.addTagsToResourceCore(ctx, store, req)
}

// RemoveTagsFromResource removes the specified tags from an SSM parameter.
func (s *SSMService) RemoveTagsFromResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.removeTagsFromResourceCore(ctx, store, req)
}

// ListTagsForResource lists all tags assigned to an SSM parameter.
func (s *SSMService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listTagsForResourceCore(ctx, store, req)
}
