package cloudwatchlogs

import (
	"context"
	"errors"
	"fmt"
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
// Every tag operation first resolves the taggable resource behind the ARN —
// a log group or a destination — so a tag against a nonexistent resource
// fails with the modelled ResourceNotFoundException instead of silently
// persisting tags under an unowned key.
func (s *LogsService) tagHandlerConfig(store *logsstore.Store) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam: "ResourceArn",
			TagsParam:     "Tags",
			TagKeysParam:  "TagKeys",
			TagKeyName:    "Key",
			TagValueName:  "Value",
			RequireTags:   true,
			// tagKeys carries "Minimum number of 0 items" (live reference;
			// the TagKeyList target's @length min 0) — an empty list is a
			// valid no-op, so no required-member rejection here; the 50
			// ceiling is checked with the element traits below.
			RequireTagKeys:     false,
			RequireResource:    true,
			CaseInsensitiveRes: true,
		},
		ValidateResource: func(_ context.Context, resourceKey string) error {
			// "Currently, the only CloudWatch Logs resources that can be
			// tagged are log groups and destinations" (TagResource): the
			// resource field is "log-group:<name>" or "destination:<name>".
			// DescribeLogGroups reports the object ARN with a trailing ":*"
			// (the log-stream namespace) — the tag operations address the
			// group itself, so the suffix is tolerated on groups alone.
			_, _, _, _, resource := svcarn.SplitARN(resourceKey)
			if name, ok := strings.CutPrefix(resource, "log-group:"); ok && name != "" {
				name = strings.TrimSuffix(name, ":*")
				if _, err := store.GetLogGroup(name); err != nil {
					return mapStoreError(err)
				}
				return nil
			}
			if name, ok := strings.CutPrefix(resource, "destination:"); ok && name != "" {
				if _, err := store.GetDestination(name); err != nil {
					return mapStoreError(err)
				}
				return nil
			}
			return ErrLogGroupNotFound
		},
		ParseTags: func(params map[string]interface{}) []tagutil.Tag {
			return tagutil.MapToTags(tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(params, "Tags")))
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return request.GetStringList(params, "TagKeys")
		},
		ValidateTagsFunc: func(tags []tagutil.Tag) error {
			// The Tags map's own length trait is 1-50. The empty set never
			// reaches this validator — the shared machinery's RequireTags
			// rejects it first as a missing member — so the trait's
			// enforceable half here is the ceiling: the 51st entry is the
			// operation's declared TooManyTagsException.
			if len(tags) > tagutil.MaxTagsPerResource {
				return NewLogsError("TooManyTagsException",
					fmt.Sprintf("A resource can have a maximum of %d tags", tagutil.MaxTagsPerResource), 400)
			}
			return validateTagEntries(tagutil.ToMap(tags))
		},
		ValidateTagKeysFunc: func(tagKeys []string) error {
			if len(tagKeys) > tagutil.MaxTagsPerResource {
				return NewLogsError("InvalidParameterException",
					fmt.Sprintf("tagKeys must contain between 0 and %d entries", tagutil.MaxTagsPerResource), 400)
			}
			return validateTagKeyElements(tagKeys)
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
		MapError: mapTagHandlerError,
	}
}

// mapTagHandlerError routes the shared tag machinery's sentinel-shaped
// errors to this service's modelled identities before the store mapping:
// the missing-tags member is InvalidParameterException, and unknown
// errors fall through to mapStoreError unchanged. The tagKeys member has
// no missing-key sentinel to map — its minimum is zero (the config's
// empty-list-is-no-op ruling), so the machinery never raises
// MissingTagKeysError for this service.
func mapTagHandlerError(err error) error {
	var missingTags *tagutil.MissingTagsError
	if errors.As(err, &missingTags) {
		return NewLogsError("InvalidParameterException",
			"The tags member is required and must not be empty", 400)
	}
	return mapStoreError(err)
}

// parseTagsFromParams reads the wire form of a Tags member — a map from
// the admin plane, a list-shaped map from the AWS plane — into the
// string map the tag-trio cores consume.
func parseTagsFromParams(params map[string]interface{}) map[string]string {
	tags := make(map[string]string)
	for k, v := range request.GetMapParamLowerFirst(params, "Tags") {
		if s, ok := v.(string); ok {
			tags[k] = s
		}
	}
	return tags
}
