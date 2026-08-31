package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagResource adds or overwrites tags on a Cognito Identity pool.
func (s *CognitoIdentityService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, cognitoIdentityTagConfig(store))
}

// UntagResource removes the specified tags from a Cognito Identity pool.
func (s *CognitoIdentityService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, cognitoIdentityTagConfig(store))
}

// ListTagsForResource lists all tags assigned to a Cognito Identity pool.
func (s *CognitoIdentityService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, cognitoIdentityTagConfig(store))
}
