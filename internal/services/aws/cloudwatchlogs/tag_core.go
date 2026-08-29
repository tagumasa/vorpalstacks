package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// tagHandlerConfig binds the shared tag-trio machinery to the CloudWatch Logs
// tag store. The closures carry the store calls, so the builder lives at the
// Core layer; handlers only pass the resolved per-region store through.
func (s *LogsService) tagHandlerConfig(store *logsstore.Store) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:      "ResourceArn",
			TagsParam:          "Tags",
			TagKeysParam:       "TagKeys",
			TagKeyName:         "Key",
			TagValueName:       "Value",
			RequireTags:        false,
			RequireTagKeys:     false,
			RequireResource:    true,
			CaseInsensitiveRes: true,
		},
		ParseTags: func(params map[string]interface{}) []tagutil.Tag {
			return tagutil.MapToTags(tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(params, "Tags")))
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return request.GetStringList(params, "TagKeys")
		},
		TagFunc: func(_ context.Context, resourceKey string, tagSlice []tagutil.Tag) error {
			return store.Tags().TagFromSlice(resourceKey, tagSlice)
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Tags().Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return store.Tags().ListAsSlice(resourceKey)
		},
		FormatResponse: func(tagSlice []tagutil.Tag, _ string) (interface{}, error) {
			m := tagutil.ToMap(tagSlice)
			if m == nil {
				m = make(map[string]string)
			}
			return map[string]interface{}{
				"tags": m,
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: mapStoreError,
	}
}
