package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// TagResource adds or overwrites tags on an AppSync API.
func (s *AppSyncService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}
	// Pre-validate tags so that specific validation errors (e.g. aws:
	// prefix reservation, 50-tag limit, key/value length) surface to the
	// caller. Without this, the ParseTags callback inside
	// appsyncTagConfig swallows the error and the handler converts the
	// resulting empty tag slice into a generic "tags are required".
	if _, err := parseTags(req.Parameters); err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, appsyncTagConfig(store, req))
}

// UntagResource removes the specified tags from an AppSync API.
func (s *AppSyncService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}
	return tags.HandleUntag(ctx, req, appsyncTagConfig(store, req))
}

// ListTagsForResource lists all tags assigned to an AppSync API.
func (s *AppSyncService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}
	return tags.HandleList(ctx, req, appsyncTagConfig(store, req))
}
