package cloudtrail

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

func resolveTrailResourceID(req *request.ParsedRequest) string {
	id := req.GetParam("ResourceId")
	if id != "" {
		return id
	}
	if resourceIdList, ok := req.Parameters["ResourceIdList"].([]interface{}); ok && len(resourceIdList) > 0 {
		if v, ok := resourceIdList[0].(string); ok {
			return v
		}
	}
	return ""
}

// AddTags adds tags to a CloudTrail trail.
func (s *CloudTrailService) AddTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := resolveTrailResourceID(req)

	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := store.GetTrailByARN(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	tagMap := tags.ToMap(tags.ParseTags(req.Parameters, "TagsList"))

	if err := store.TagResource(trail.Name, tagMap); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// RemoveTags removes tags from a CloudTrail trail.
func (s *CloudTrailService) RemoveTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := resolveTrailResourceID(req)

	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := store.GetTrailByARN(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	tagKeys := tags.ParseTagKeysAsSlice(req.Parameters, "TagKeyList")
	if len(tagKeys) == 0 {
		tagKeys = tags.ParseTagKeysWithKeyName(req.Parameters, "TagsList", "Key")
	}

	if len(tagKeys) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := store.UntagResource(trail.Name, tagKeys); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTags retrieves the tags associated with a CloudTrail trail.
func (s *CloudTrailService) ListTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := resolveTrailResourceID(req)

	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := store.GetTrailByARN(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	tagSlice, err := store.ListTagsAsSlice(trail.Name)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	formattedTags := tags.ToResponse(tagSlice)

	return map[string]interface{}{
		"ResourceTagList": []map[string]interface{}{
			{
				"ResourceId": resourceARN,
				"TagsList":   formattedTags,
			},
		},
	}, nil
}
