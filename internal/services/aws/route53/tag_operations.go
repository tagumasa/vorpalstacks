package route53

import (
	"context"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

// rawResourceId reads the resource identifier, accepting both the
// documented ResourceId member and the ResourceID spelling some clients
// send.
func rawResourceId(params map[string]interface{}) string {
	resourceId := request.GetStringParam(params, "ResourceId")
	if resourceId == "" {
		resourceId = request.GetStringParam(params, "ResourceID")
	}
	return resourceId
}

// parseAddTags extracts the AddTags member in every wire shape the API
// accepts (flat list, query fallback, or wrapped single Tag map).
func parseAddTags(params map[string]interface{}) []tags.Tag {
	addTags := tags.ParseTagsWithKeyNames(params, "AddTags", "Key", "Value")
	if len(addTags) == 0 {
		addTags = tags.ParseTags(params, "Tags")
	}
	if len(addTags) == 0 {
		addTags = tags.ParseTagsWithQueryFallback(params, "AddTags")
	}
	if len(addTags) == 0 {
		if addTagsMap, ok := params["AddTags"].(map[string]interface{}); ok {
			if tagList, ok := addTagsMap["Tag"]; ok {
				switch tl := tagList.(type) {
				case []interface{}:
					for _, t := range tl {
						if tagMap, ok := t.(map[string]interface{}); ok {
							addTags = append(addTags, tags.Tag{
								Key:   request.GetStringParam(tagMap, "Key"),
								Value: request.GetStringParam(tagMap, "Value"),
							})
						}
					}
				case map[string]interface{}:
					addTags = append(addTags, tags.Tag{
						Key:   request.GetStringParam(tl, "Key"),
						Value: request.GetStringParam(tl, "Value"),
					})
				}
			}
		}
	}
	return addTags
}

// parseRemoveTagKeys extracts the RemoveTagKeys member in every wire shape
// the API accepts (query fallback, key slice, or wrapped single Key map).
func parseRemoveTagKeys(params map[string]interface{}) []string {
	removeTagKeys := tags.ParseTagKeysWithQueryFallback(params, "RemoveTagKeys")
	if len(removeTagKeys) == 0 {
		removeTagKeys = tags.ParseTagKeysAsSlice(params, "RemoveTagKeys")
	}
	if len(removeTagKeys) == 0 {
		if rkMap, ok := params["RemoveTagKeys"].(map[string]interface{}); ok {
			if k, ok := rkMap["Key"].(string); ok {
				removeTagKeys = []string{k}
			} else if k, ok := rkMap["key"].(string); ok {
				removeTagKeys = []string{k}
			}
		}
	}
	return removeTagKeys
}

// ChangeTagsForResource adds or removes tags for a Route 53 resource.
func (s *Route53Service) ChangeTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := changeTagsForResourceCore(st, ChangeTagsForResourceInput{
		ResourceType:  request.GetStringParam(req.Parameters, "ResourceType"),
		ResourceId:    rawResourceId(req.Parameters),
		AddTags:       parseAddTags(req.Parameters),
		RemoveTagKeys: parseRemoveTagKeys(req.Parameters),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for a Route 53 resource.
func (s *Route53Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := listTagsForResourceCore(st, ListTagsForResourceInput{
		ResourceType: request.GetStringParam(req.Parameters, "ResourceType"),
		ResourceId:   rawResourceId(req.Parameters),
	})
	if err != nil {
		return nil, err
	}

	tagItems := make([]interface{}, 0, len(result.Tags))
	for _, t := range result.Tags {
		tagItems = append(tagItems, map[string]interface{}{
			"Key":   t.Key,
			"Value": t.Value,
		})
	}

	return map[string]interface{}{
		"ResourceTagSet": map[string]interface{}{
			"ResourceType": result.ResourceType,
			"ResourceId":   result.ResourceId,
			"Tags":         protocol.XMLElements{ElementName: "Tag", Items: tagItems},
		},
	}, nil
}
