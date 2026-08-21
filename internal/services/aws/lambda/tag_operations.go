package lambda

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

// esmTagResourceKey namespaces an event source mapping's tags inside the
// shared tag store so mapping UUIDs cannot collide with function names.
func esmTagResourceKey(uuid string) string {
	return "event-source-mapping/" + uuid
}

// splitTaggableResource splits a taggable-resource ARN into its resource
// type, identifier and qualifier, following the TaggableResource
// pattern: function:<name> with an optional :qualifier, or
// event-source-mapping:<uuid>.
func splitTaggableResource(rawARN string) (resourceType, id, qualifier string) {
	_, _, _, _, resource := svcarn.SplitARN(rawARN)
	parts := strings.SplitN(resource, ":", 3)
	if len(parts) < 2 {
		return "", "", ""
	}
	if len(parts) == 3 {
		qualifier = parts[2]
	}
	return parts[0], parts[1], qualifier
}

func lambdaTagConfig(s *LambdaService, reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.LambdaConfig,
		MapError: func(err error) error {
			switch e := err.(type) {
			case *tagutil.MissingResourceError:
				return NewInvalidParameter("Resource", e.Param+" is required")
			case *tagutil.MissingTagsError:
				return NewInvalidParameter("Tags", e.Param+" is required")
			case *tagutil.MissingTagKeysError:
				return NewInvalidParameter("TagKeys", e.Param+" is required")
			}
			return err
		},
		ResourceKey: func(rawKey string) string {
			resourceType, id, qualifier := splitTaggableResource(rawKey)
			switch resourceType {
			case "event-source-mapping":
				return esmTagResourceKey(id)
			case "function":
				if qualifier != "" {
					// Keep the qualifier in the key so validation can
					// reject it: "Lambda does not support adding tags to
					// function aliases or versions."
					return id + ":" + qualifier
				}
				return id
			}
			return svcarn.ExtractFunctionNameFromARN(rawKey)
		},
		ParseTags: func(params map[string]interface{}) []types.Tag {
			return tagutil.ParseTags(params, "Tags")
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			seen := make(map[string]bool)
			for _, paramKey := range []string{"TagKeys", "tagKeys"} {
				if v, ok := params[paramKey]; ok {
					switch val := v.(type) {
					case []interface{}:
						for _, k := range val {
							if ks, ok := k.(string); ok {
								seen[ks] = true
							}
						}
					case []string:
						for _, ks := range val {
							seen[ks] = true
						}
					case string:
						if val != "" {
							seen[val] = true
						}
					}
				}
			}
			result := make([]string, 0, len(seen))
			for k := range seen {
				result = append(result, k)
			}
			return result
		},
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			if strings.Contains(resourceKey, ":") {
				// A qualified function ARN: "Lambda does not support
				// adding tags to function aliases or versions."
				return NewInvalidParameter("Resource",
					"Lambda does not support adding tags to function aliases or versions")
			}
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			if uuid, ok := strings.CutPrefix(resourceKey, "event-source-mapping/"); ok {
				if _, err := store.EventSources.Get(uuid); err != nil {
					return ErrResourceNotFound
				}
				return nil
			}
			_, err = store.Functions.Get(resourceKey)
			if err != nil {
				return ErrResourceNotFound
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []types.Tag) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.Functions.TagStore.TagFromSlice(resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.Functions.TagStore.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			store, err := s.store(reqCtx)
			if err != nil {
				return nil, err
			}
			return store.Functions.TagStore.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tags []types.Tag, rawResourceKey string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ToMap(tags),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) { return response.EmptyResponse(), nil },
	}
}

// TagResource adds tags to a Lambda function.
func (s *LambdaService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleTag(ctx, req, lambdaTagConfig(s, reqCtx))
}

// UntagResource removes tags from a Lambda function.
func (s *LambdaService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleUntag(ctx, req, lambdaTagConfig(s, reqCtx))
}

// ListTags lists the tags for a Lambda function.
func (s *LambdaService) ListTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleList(ctx, req, lambdaTagConfig(s, reqCtx))
}
