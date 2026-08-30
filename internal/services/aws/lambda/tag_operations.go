package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// esmTagResourceKey namespaces an event source mapping's tags inside the
// shared tag store so mapping UUIDs cannot collide with function names.
func esmTagResourceKey(uuid string) string {
	return "event-source-mapping/" + uuid
}

// TagResource adds tags to a Lambda function.
func (s *LambdaService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleTag(ctx, req, lambdaTagConfig(s, reqCtx))
}

// UntagResource removes tags from a Lambda function.
func (s *LambdaService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleUntag(ctx, req, lambdaTagConfig(s, reqCtx))
}

// ListTags lists the tags for a Lambda function.
func (s *LambdaService) ListTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleList(ctx, req, lambdaTagConfig(s, reqCtx))
}
