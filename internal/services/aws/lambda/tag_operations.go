package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

func lambdaTagConfig(s *LambdaService, reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.LambdaConfig,
		ResourceKey: func(rawKey string) string {
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
		ValidateResource: func(ctx context.Context, rawKey string) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			functionName := svcarn.ExtractFunctionNameFromARN(rawKey)
			_, err = store.Functions.Get(functionName)
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
