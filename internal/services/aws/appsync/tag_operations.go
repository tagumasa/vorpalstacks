package appsync

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
)

// TagResource adds or updates tags on an AppSync resource.
// POST /v1/tags/{resourceArn}
// Body: {"tags": {"key": "value", ...}}
func (s *AppSyncService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, NewBadRequestException("resourceArn is required")
	}

	newTags := parseTags(req.Parameters)
	if len(newTags) == 0 {
		return nil, NewBadRequestException("tags are required")
	}

	if err := s.tagResource(store, resourceArn, newTags); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// UntagResource removes tags from an AppSync resource.
// DELETE /v1/tags/{resourceArn}?tagKeys=key1&tagKeys=key2
func (s *AppSyncService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, NewBadRequestException("resourceArn is required")
	}

	tagKeys := parseTagKeysFromQuery(req)
	if len(tagKeys) == 0 {
		return nil, NewBadRequestException("tagKeys are required")
	}

	if err := s.untagResource(store, resourceArn, tagKeys); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListTagsForResource returns the tags for an AppSync resource.
// GET /v1/tags/{resourceArn}
func (s *AppSyncService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, NewBadRequestException("resourceArn is required")
	}

	tags, err := s.listResourceTags(store, resourceArn)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tags,
	}, nil
}

// parseTagKeysFromQuery extracts tag keys from the query string.
// AppSync sends tag keys as repeated query parameters: ?tagKeys=key1&tagKeys=key2
func parseTagKeysFromQuery(req *request.ParsedRequest) []string {
	keys := req.QueryParams["tagKeys"]
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		if k != "" {
			result = append(result, k)
		}
	}
	return result
}
