package secretsmanager

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagResource adds or overwrites tags on a Secrets Manager secret.
func (s *SecretsManagerService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, secretsManagerTagConfig(store, s))
}

// UntagResource removes the specified tags from a Secrets Manager secret.
func (s *SecretsManagerService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, secretsManagerTagConfig(store, s))
}

// ListTagsForResource lists all tags assigned to a Secrets Manager secret.
func (s *SecretsManagerService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, secretsManagerTagConfig(store, s))
}
