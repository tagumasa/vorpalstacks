package route53

import (
	"context"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// maxTagsPerRoute53Resource is the ChangeTagsForResource API tag quota: the
// API and CLI accept at most 10 tags per hosted zone or health check (the
// newer console allows 50, but the API this server implements does not).
const maxTagsPerRoute53Resource = 10

// route53TagLimits applies the API tag count quota on top of the standard
// key and value length bounds.
var route53TagLimits = tags.TagLimits{
	MaxCount:       maxTagsPerRoute53Resource,
	MaxKeyLength:   tags.MaxTagKeyLength,
	MaxValueLength: tags.MaxTagValueLength,
}

// parseResourceParams validates and normalises ResourceType/ResourceId,
// returning (normalisedType, bareId, resourceKey, error).
func parseResourceParams(params map[string]interface{}) (string, string, string, error) {
	resourceType := request.GetStringParam(params, "ResourceType")
	resourceId := request.GetStringParam(params, "ResourceId")
	if resourceId == "" {
		resourceId = request.GetStringParam(params, "ResourceID")
	}

	if resourceType == "" || resourceId == "" {
		return "", "", "", awserrors.NewAWSError("InvalidParameter", "ResourceType and ResourceId are required", 400)
	}

	normalizedType := strings.ToLower(resourceType)
	if normalizedType != "hostedzone" && normalizedType != "healthcheck" {
		return "", "", "", awserrors.NewAWSError("InvalidParameter", "ResourceType must be 'hostedzone' or 'healthcheck'", 400)
	}

	resourceId = strings.TrimPrefix(resourceId, "/hostedzone/")
	resourceId = strings.TrimPrefix(resourceId, "/healthcheck/")

	return normalizedType, resourceId, normalizedType + "/" + resourceId, nil
}

// ChangeTagsForResource adds or removes tags for a Route 53 resource.
func (s *Route53Service) ChangeTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	normalizedType, resourceId, resourceKey, err := parseResourceParams(req.Parameters)
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Verify resource exists before applying tag operations.
	if err := s.verifyResourceExists(st, normalizedType, resourceId); err != nil {
		return nil, err
	}

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

	if v, _ := tags.CheckTags(addTags, route53TagLimits); v != tags.OK {
		switch v {
		case tags.TooManyTags:
			return nil, awserrors.NewAWSError("InvalidInput",
				fmt.Sprintf("Number of tags must not exceed %d", maxTagsPerRoute53Resource), 400)
		case tags.TagKeyTooLong:
			return nil, awserrors.NewAWSError("InvalidInput", "Tag key must not exceed 128 characters", 400)
		case tags.TagValueTooLong:
			return nil, awserrors.NewAWSError("InvalidInput", "Tag value must not exceed 256 characters", 400)
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

	// Enforce 50-tag limit — compute the resulting key set
	// after applying both AddTags and RemoveTagKeys.
	existingTags, _ := st.Tags().ListTagsForResource(resourceKey)
	keySet := make(map[string]bool)
	for _, t := range existingTags {
		keySet[t.Key] = true
	}
	for _, k := range removeTagKeys {
		delete(keySet, k)
	}
	for _, t := range addTags {
		keySet[t.Key] = true
	}
	if len(keySet) > 50 {
		return nil, awserrors.NewAWSError("InvalidInput", "Maximum of 50 tags allowed per resource", 400)
	}

	if len(addTags) > 0 {
		if err := st.Tags().Tag(resourceKey, addTags); err != nil {
			return nil, awserrors.NewAWSError("TagResource", err.Error(), 500)
		}
	}

	if len(removeTagKeys) > 0 {
		if err := st.Tags().Raw().Untag(resourceKey, removeTagKeys); err != nil {
			return nil, awserrors.NewAWSError("UntagResource", err.Error(), 500)
		}
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for a Route 53 resource.
func (s *Route53Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	normalizedType, resourceId, resourceKey, err := parseResourceParams(req.Parameters)
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Verify resource exists before listing tags.
	if err := s.verifyResourceExists(st, normalizedType, resourceId); err != nil {
		return nil, err
	}

	tags, err := st.Tags().ListTagsForResource(resourceKey)
	if err != nil {
		return nil, awserrors.NewAWSError("ListTags", err.Error(), 500)
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
			"ResourceType": normalizedType,
			"ResourceId":   resourceId,
			"Tags":         protocol.XMLElements{ElementName: "Tag", Items: tagItems},
		},
	}, nil
}

// verifyResourceExists checks whether the specified resource exists
// in the store, returning NoSuchHostedZone or NoSuchHealthCheck if not.
func (s *Route53Service) verifyResourceExists(st *route53store.Route53Stores, resourceType, resourceId string) error {
	switch resourceType {
	case "hostedzone":
		if !st.HostedZones().Exists(resourceId) {
			return NewNoSuchHostedZoneError(resourceId)
		}
	case "healthcheck":
		if !st.HealthChecks().Exists(resourceId) {
			return NewNoSuchHealthCheckError(resourceId)
		}
	}
	return nil
}
