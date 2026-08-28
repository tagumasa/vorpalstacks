package neptunegraph

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// ListTagsForResource returns all tags associated with the specified resource ARN.
func (s *NeptuneGraphService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ListTagsForResourceInput{
		ResourceArn: request.GetStringParam(req.Parameters, "resourceArn"),
	}

	resourceTags, err := s.listTagsForResourceCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": resourceTags,
	}, nil
}

// TagResource adds or updates tags on the specified resource, validating key and value constraints.
func (s *NeptuneGraphService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &TagResourceInput{
		ResourceArn: request.GetStringParam(req.Parameters, "resourceArn"),
		Tags:        parseTagsFromParams(req.Parameters),
	}

	resourceTags, err := s.tagResourceCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": resourceTags,
	}, nil
}

// UntagResource removes the specified tag keys from a resource's tag set.
func (s *NeptuneGraphService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &UntagResourceInput{
		ResourceArn: request.GetStringParam(req.Parameters, "resourceArn"),
		TagKeys:     tags.ParseTagKeysAsSlice(req.Parameters, "tagKeys"),
	}

	resourceTags, err := s.untagResourceCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": resourceTags,
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
