package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagResource adds tags to a Cognito User Pool or User Pool Client resource.
func (s *CognitoService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	tags := tagutil.GetTags(req.Parameters, tagutil.StandardConfig)
	if len(tags) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.TagResource(resourceArn, tagutil.ToMap(tags)); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a Cognito User Pool or User Pool Client resource.
func (s *CognitoService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	tagKeys := tagutil.GetTagKeys(req.Parameters, tagutil.StandardConfig)
	if len(tagKeys) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.UntagResource(resourceArn, tagKeys); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for a Cognito User Pool or User Pool Client resource.
func (s *CognitoService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tags, err := store.ListTagsAsSlice(resourceArn)
	if err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"Tags": tagutil.ToMap(tags),
	}, nil
}
