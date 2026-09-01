package cloudwatchlogs

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// tagHandlerConfig binds the shared tag-trio machinery to the CloudWatch Logs
// tag store. The closures carry the store calls, so the builder lives at the
// Core layer; handlers only pass the resolved per-region store through.
// Every tag operation first resolves the log group behind the ARN, so a tag
// against a nonexistent log group fails with the modelled
// ResourceNotFoundException instead of silently persisting tags under an
// unowned key.
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
		ValidateResource: func(_ context.Context, resourceKey string) error {
			// CloudWatch Logs tags log groups; the resource field is
			// "log-group:<name>". DescribeLogGroups reports the object ARN
			// with a trailing ":*" (the log-stream namespace) — the tag
			// operations address the group itself, so the suffix is
			// tolerated here.
			_, _, _, _, resource := svcarn.SplitARN(resourceKey)
			name, ok := strings.CutPrefix(resource, "log-group:")
			if !ok || name == "" {
				return ErrLogGroupNotFound
			}
			name = strings.TrimSuffix(name, ":*")
			if _, err := store.GetLogGroup(name); err != nil {
				return mapStoreError(err)
			}
			return nil
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
