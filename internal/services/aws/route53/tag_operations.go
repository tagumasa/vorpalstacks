package route53

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"
)

// ChangeTagsForResource adds or removes tags for a Route 53 resource.
func (s *Route53Service) ChangeTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceType := request.GetStringParam(req.Parameters, "ResourceType")
	resourceId := request.GetStringParam(req.Parameters, "ResourceId")
	if resourceId == "" {
		resourceId = request.GetStringParam(req.Parameters, "ResourceID")
	}

	if resourceType == "" || resourceId == "" {
		return nil, NewAPIError("InvalidParameter", "ResourceType and ResourceId are required", 400)
	}

	normalizedType := strings.ToLower(resourceType)
	if normalizedType != "hostedzone" && normalizedType != "healthcheck" {
		return nil, NewAPIError("InvalidParameter", "ResourceType must be 'hostedzone' or 'healthcheck'", 400)
	}

	resourceId = strings.TrimPrefix(resourceId, "/hostedzone/")
	resourceId = strings.TrimPrefix(resourceId, "/healthcheck/")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceKey := normalizedType + "/" + resourceId

	addTags := tags.ParseTagsWithKeyNames(req.Parameters, "AddTags", "Key", "Value")
	if len(addTags) == 0 {
		addTags = tags.ParseTags(req.Parameters, "Tags")
	}
	if len(addTags) == 0 {
		addTags = tags.ParseTagsWithQueryFallback(req.Parameters, "AddTags")
	}
	if len(addTags) == 0 {
		if addTagsMap, ok := req.Parameters["AddTags"].(map[string]interface{}); ok {
			if tagList, ok := addTagsMap["Tag"]; ok {
				switch tl := tagList.(type) {
				case []interface{}:
					for _, t := range tl {
						if tagMap, ok := t.(map[string]interface{}); ok {
							addTags = append(addTags, types.Tag{
								Key:   request.GetStringParam(tagMap, "Key"),
								Value: request.GetStringParam(tagMap, "Value"),
							})
						}
					}
				case map[string]interface{}:
					addTags = append(addTags, types.Tag{
						Key:   request.GetStringParam(tl, "Key"),
						Value: request.GetStringParam(tl, "Value"),
					})
				}
			}
		}
	}

	removeTagKeys := tags.ParseTagKeysWithQueryFallback(req.Parameters, "RemoveTagKeys")
	if len(removeTagKeys) == 0 {
		removeTagKeys = tags.ParseTagKeysAsSlice(req.Parameters, "RemoveTagKeys")
	}
	if len(removeTagKeys) == 0 {
		if rkMap, ok := req.Parameters["RemoveTagKeys"].(map[string]interface{}); ok {
			if k, ok := rkMap["Key"].(string); ok {
				removeTagKeys = []string{k}
			} else if k, ok := rkMap["key"].(string); ok {
				removeTagKeys = []string{k}
			}
		}
	}

	if len(addTags) > 0 {
		if err := st.Tags().TagResource(resourceKey, addTags); err != nil {
			return nil, NewAPIError("TagResource", err.Error(), 500)
		}
	}

	if len(removeTagKeys) > 0 {
		if err := st.Tags().Raw().UntagResource(resourceKey, removeTagKeys); err != nil {
			return nil, NewAPIError("UntagResource", err.Error(), 500)
		}
		if len(addTags) > 0 {
			if err := st.Tags().TagResource(resourceKey, addTags); err != nil {
				return nil, NewAPIError("TagResource", err.Error(), 500)
			}
		}
		return response.EmptyResponse(), nil
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for a Route 53 resource.
func (s *Route53Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceType := request.GetStringParam(req.Parameters, "ResourceType")
	resourceId := request.GetStringParam(req.Parameters, "ResourceId")
	if resourceId == "" {
		resourceId = request.GetStringParam(req.Parameters, "ResourceID")
	}

	if resourceType == "" || resourceId == "" {
		return nil, NewAPIError("InvalidParameter", "ResourceType and ResourceId are required", 400)
	}

	normalizedResourceType := strings.ToLower(resourceType)
	if normalizedResourceType != "hostedzone" && normalizedResourceType != "healthcheck" {
		return nil, NewAPIError("InvalidParameter", "ResourceType must be 'hostedzone' or 'healthcheck'", 400)
	}

	resourceId = strings.TrimPrefix(resourceId, "/hostedzone/")
	resourceId = strings.TrimPrefix(resourceId, "/healthcheck/")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceKey := normalizedResourceType + "/" + resourceId

	tags, err := st.Tags().ListTagsForResource(resourceKey)
	if err != nil {
		return nil, NewAPIError("ListTags", err.Error(), 500)
	}

	tagItems := make([]interface{}, 0, len(tags))
	for _, t := range tags {
		tagItems = append(tagItems, map[string]interface{}{
			"Key":   t.Key,
			"Value": t.Value,
		})
	}

	return map[string]interface{}{
		"ResourceTagSet": map[string]interface{}{
			"ResourceType": normalizedResourceType,
			"ResourceId":   extractResourceId(resourceId),
			"Tags":         protocol.XMLElements{ElementName: "Tag", Items: tagItems},
		},
	}, nil
}

func extractResourceId(resourceId string) string {
	if strings.HasPrefix(resourceId, "/") {
		parts := strings.Split(strings.Trim(resourceId, "/"), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return resourceId
}
