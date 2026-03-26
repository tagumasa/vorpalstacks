package cloudtrail

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	"vorpalstacks/internal/services/aws/common/tags"
)

// AddTags adds tags to a CloudTrail trail.
func (s *CloudTrailService) AddTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := getParam(req, "ResourceId")
	if resourceARN == "" {
		if resourceIdList, ok := req.Parameters["ResourceIdList"].([]interface{}); ok && len(resourceIdList) > 0 {
			if id, ok := resourceIdList[0].(string); ok {
				resourceARN = id
			}
		}
	}

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

	tagList := tags.ParseTags(req.Parameters, "TagsList")
	tagMap := make(map[string]string)
	for _, t := range tagList {
		tagMap[t.Key] = t.Value
	}

	if err := store.TagResource(trail.Name, tagMap); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// RemoveTags removes tags from a CloudTrail trail.
func (s *CloudTrailService) RemoveTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := getParam(req, "ResourceId")
	if resourceARN == "" {
		if resourceIdList, ok := req.Parameters["ResourceIdList"].([]interface{}); ok && len(resourceIdList) > 0 {
			if id, ok := resourceIdList[0].(string); ok {
				resourceARN = id
			}
		}
	}

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
		return response.EmptyResponse(), nil
	}

	if err := store.UntagResource(trail.Name, tagKeys); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTags retrieves the tags associated with a CloudTrail trail.
func (s *CloudTrailService) ListTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := getParam(req, "ResourceId")
	if resourceARN == "" {
		if resourceIdList, ok := req.Parameters["ResourceIdList"].([]interface{}); ok && len(resourceIdList) > 0 {
			if id, ok := resourceIdList[0].(string); ok {
				resourceARN = id
			}
		}
	}

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

	formattedTags := make([]map[string]interface{}, 0, len(tagSlice))
	for _, t := range tagSlice {
		formattedTags = append(formattedTags, map[string]interface{}{
			"Key":   t.Key,
			"Value": t.Value,
		})
	}

	return map[string]interface{}{
		"ResourceTagList": []map[string]interface{}{
			{
				"ResourceId": resourceARN,
				"TagsList":   formattedTags,
			},
		},
	}, nil
}
