package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
)

// TagResource adds tags to a Cognito identity pool resource.
func (s *CognitoIdentityService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

// UntagResource removes tags from a Cognito identity pool resource.
func (s *CognitoIdentityService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

// ListTagsForResource lists tags for a Cognito identity pool resource.
func (s *CognitoIdentityService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tags, err := store.ListTags(resourceArn)
	if err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"Tags": tagutil.MapToResponse(tags),
	}, nil
}
