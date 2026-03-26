package ssm

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
)

// AddTagsToResource adds tags to an SSM resource.
func (s *SSMService) AddTagsToResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = getParam(req, "ResourceType")
	resourceID := getParam(req, "ResourceId")
	if resourceID == "" {
		return nil, ErrInvalidParameterName
	}

	tags := tagutil.GetTags(req.Parameters, tagutil.StandardConfig)
	if len(tags) > 0 {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		if err := store.AddTagsToResource(resourceID, tagutil.ToMap(tags)); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// RemoveTagsFromResource removes tags from an SSM resource.
func (s *SSMService) RemoveTagsFromResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = getParam(req, "ResourceType")
	resourceID := getParam(req, "ResourceId")
	if resourceID == "" {
		return nil, ErrInvalidParameterName
	}

	tagKeys := tagutil.GetTagKeys(req.Parameters, tagutil.StandardConfig)
	if len(tagKeys) > 0 {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		if err := store.RemoveTagsFromResource(resourceID, tagKeys); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for an SSM resource.
func (s *SSMService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = getParam(req, "ResourceType")
	resourceID := getParam(req, "ResourceId")
	if resourceID == "" {
		return nil, ErrInvalidParameterName
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tags, err := store.ListTagsForResource(resourceID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TagList": tagutil.MapToResponse(tags),
	}, nil
}
