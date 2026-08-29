package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagResource adds tags to an SESv2 resource.
func (s *SESv2Service) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleTag(ctx, req, sesv2TagConfig(s, reqCtx))
}

// UntagResource removes tags from an SESv2 resource.
func (s *SESv2Service) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleUntag(ctx, req, sesv2TagConfig(s, reqCtx))
}

// ListTagsForResource lists tags for an SESv2 resource.
func (s *SESv2Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleList(ctx, req, sesv2TagConfig(s, reqCtx))
}
