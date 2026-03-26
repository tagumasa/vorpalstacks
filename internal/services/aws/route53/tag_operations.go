package route53

import (
	"context"
	"strings"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
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

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceKey := resourceType + "/" + resourceId

	addTags := tagutil.ParseTagsWithKeyNames(req.Parameters, "AddTags", "Key", "Value")
	if len(addTags) == 0 {
		addTags = tagutil.ParseTags(req.Parameters, "Tags")
	}
	if len(addTags) == 0 {
		addTags = tagutil.ParseTagsWithQueryFallback(req.Parameters, "AddTags")
	}

	removeTagKeys := tagutil.ParseTagKeysWithQueryFallback(req.Parameters, "RemoveTagKeys")
	if len(removeTagKeys) == 0 {
		removeTagKeys = tagutil.ParseTagKeysAsSlice(req.Parameters, "RemoveTagKeys")
	}

	if len(addTags) > 0 {
		if err := st.Tags().TagResource(resourceKey, addTags); err != nil {
			return nil, NewAPIError("TagResource", err.Error(), 500)
		}
	}

	if len(removeTagKeys) > 0 {
		existingTags, err := st.Tags().ListTagsForResource(resourceKey)
		if err != nil {
			return nil, NewAPIError("ListTags", err.Error(), 500)
		}
		keysToRemoveMap := make(map[string]bool)
		for _, k := range removeTagKeys {
			keysToRemoveMap[k] = true
		}
		remainingTags := tagutil.Remove(existingTags, keysToRemoveMap)
		if err := st.Tags().TagResource(resourceKey, remainingTags); err != nil {
			return nil, NewAPIError("UntagResource", err.Error(), 500)
		}
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

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceKey := resourceType + "/" + resourceId

	tags, err := st.Tags().ListTagsForResource(resourceKey)
	if err != nil {
		return nil, NewAPIError("ListTags", err.Error(), 500)
	}

	normalizedResourceType := strings.ToLower(resourceType)
	switch normalizedResourceType {
	case "healthcheck", "hostedzone":
	default:
		normalizedResourceType = resourceType
	}

	var tagItems []interface{}
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
