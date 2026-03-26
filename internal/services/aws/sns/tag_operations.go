package sns

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	"vorpalstacks/internal/utils/aws/types"
)

// TagResource adds tags to an SNS topic.
func (s *SNSService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, NewInvalidParameterException("ResourceArn is required")
	}

	tags := tagutil.GetTags(req.Parameters, tagutil.StandardConfig)
	if len(tags) == 0 {
		return nil, NewInvalidParameterException("Tags is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.TagResource(resourceArn, tags); err != nil {
		return nil, err
	}

	return tagutil.HandleListSimple(ctx, tagutil.StandardConfig, resourceArn,
		func(key string) ([]types.Tag, error) { return store.ListTagsForResource(key) },
	)
}

// UntagResource removes tags from an SNS topic.
func (s *SNSService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, NewInvalidParameterException("ResourceArn is required")
	}

	tagKeys := tagutil.GetTagKeys(req.Parameters, tagutil.StandardConfig)
	if len(tagKeys) == 0 {
		return nil, NewInvalidParameterException("TagKeys is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.UntagResource(resourceArn, tagKeys); err != nil {
		return nil, err
	}

	return nil, nil
}

// ListTagsForResource lists tags for an SNS topic.
func (s *SNSService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, NewInvalidParameterException("ResourceArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return tagutil.HandleListSimple(ctx, tagutil.StandardConfig, resourceArn,
		func(key string) ([]types.Tag, error) { return store.ListTagsForResource(key) },
	)
}
