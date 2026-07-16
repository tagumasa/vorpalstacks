package neptunegraph

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// ListTagsForResource returns all tags associated with the specified resource ARN.
func (s *NeptuneGraphService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	tags, err := store.GetTags(resourceArn)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tags,
	}, nil
}

// TagResource adds or updates tags on the specified resource, validating key and value constraints.
func (s *NeptuneGraphService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	tags := parseTagsFromParams(req.Parameters)
	if len(tags) > maxTags {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "too many tags")
	}
	for k, v := range tags {
		if len(k) < 1 || len(k) > 128 || strings.HasPrefix(k, "aws:") || !tagKeyRegex.MatchString(k) {
			return nil, newValidationException("ILLEGAL_ARGUMENT", "invalid tag key")
		}
		if len(v) > 256 {
			return nil, newValidationException("ILLEGAL_ARGUMENT", "tag value too long")
		}
	}

	if err := store.AddTags(resourceArn, tags); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tags,
	}, nil
}

// UntagResource removes the specified tag keys from a resource's tag set.
func (s *NeptuneGraphService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	tagKeys := tags.ParseTagKeysAsSlice(req.Parameters, "tagKeys")

	if err := store.RemoveTags(resourceArn, tagKeys); err != nil {
		return nil, err
	}

	tags, _ := store.GetTags(resourceArn)
	return map[string]interface{}{
		"tags": tags,
	}, nil
}

func parseTagsFromParams(params map[string]interface{}) map[string]string {
	v, ok := params["tags"]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	tags := make(map[string]string)
	key, hasKey := m["key"]
	val, hasVal := m["value"]
	if hasKey && hasVal {
		if ks, ok := key.(string); ok {
			if vs, ok := val.(string); ok {
				tags[ks] = vs
				return tags
			}
		}
	}
	keyUpper, hasKeyUpper := m["Key"]
	valUpper, hasValUpper := m["Value"]
	if hasKeyUpper && hasValUpper {
		if ks, ok := keyUpper.(string); ok {
			if vs, ok := valUpper.(string); ok {
				tags[ks] = vs
				return tags
			}
		}
	}
	for k, v := range m {
		if s, ok := v.(string); ok {
			tags[k] = s
		}
	}
	return tags
}
