package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagResource adds or overwrites tags on a Cognito IdP resource.
func (s *CognitoService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	tags := tagutil.ParseTagsAsMap(req.Parameters, "Tags")
	if err := s.tagResourceCore(reqCtx.GetRegion(), resourceArn, tags); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UntagResource removes the specified tags from a Cognito IdP resource.
func (s *CognitoService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	tagKeys := tagutil.ParseTagKeysAsSlice(req.Parameters, "TagKeys")
	if err := s.untagResourceCore(reqCtx.GetRegion(), resourceArn, tagKeys); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListTagsForResource lists all tags assigned to a Cognito IdP resource.
func (s *CognitoService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	tags, err := s.listTagsForResourceCore(reqCtx.GetRegion(), resourceArn)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"Tags": tags}, nil
}
