package cloudtrail

import (
	"context"
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// applyTags parses tags from the raw wire value — the JSON-protocol list
// form or a JSON string holding the same list — into dst, initialising a nil
// map. Entries with an empty key are ignored. It is the single tag-apply
// path shared by every tag-bearing CloudTrail resource (trails, event data
// stores, channels), whose stores all hold tags as map[string]string.
func applyTags(dst *map[string]string, raw interface{}) {
	if *dst == nil {
		*dst = make(map[string]string)
	}
	var tagsList []interface{}
	switch v := raw.(type) {
	case []interface{}:
		tagsList = v
	case string:
		if err := json.Unmarshal([]byte(v), &tagsList); err != nil {
			return
		}
	default:
		return
	}
	for _, item := range tagsList {
		if m, ok := item.(map[string]interface{}); ok {
			key, _ := m["Key"].(string)
			val, _ := m["Value"].(string)
			if key != "" {
				(*dst)[key] = val
			}
		}
	}
}

// formatTagsList renders a tag map as the TagsList wire shape (a list of
// {Key, Value} maps).
func formatTagsList(tags map[string]string) []interface{} {
	tagsList := make([]interface{}, 0, len(tags))
	for k, v := range tags {
		tagsList = append(tagsList, map[string]interface{}{
			"Key":   k,
			"Value": v,
		})
	}
	return tagsList
}

// cloudTrailTagLimits applies the standard AWS tag bounds with the aws:
// reservation compared case-sensitively, preserving the established
// CloudTrail behaviour.
var cloudTrailTagLimits = tagutil.TagLimits{
	MaxCount:              tagutil.MaxTagsPerResource,
	MinKeyLength:          1,
	MaxKeyLength:          tagutil.MaxTagKeyLength,
	MaxValueLength:        tagutil.MaxTagValueLength,
	ReservedPrefix:        "aws:",
	ReservedCaseSensitive: true,
}

// validateCloudTrailTags validates tag count, key length, value length, and
// reserved prefix against AWS CloudTrail limits.
func validateCloudTrailTags(tagList []tagutil.Tag) error {
	switch v, _ := tagutil.CheckTags(tagList, cloudTrailTagLimits); v {
	case tagutil.TooManyTags:
		return newTagsLimitExceededException(
			fmt.Sprintf("Number of tags exceeds the limit of %d", tagutil.MaxTagsPerResource))
	case tagutil.ReservedTagKey:
		return newInvalidTagParameterException(
			"Tag keys starting with 'aws:' are reserved")
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return newInvalidTagParameterException(
			fmt.Sprintf("Tag key length must be between 1 and %d", tagutil.MaxTagKeyLength))
	case tagutil.TagValueTooLong:
		return newInvalidTagParameterException(
			fmt.Sprintf("Tag value length must not exceed %d", tagutil.MaxTagValueLength))
	}
	return nil
}

func cloudTrailMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError:
		return ErrInvalidParameter
	case *tagutil.MissingTagsError:
		return ErrInvalidParameter
	case *tagutil.MissingTagKeysError:
		return ErrInvalidParameter
	}
	return err
}

// cloudTrailTagConfig builds the shared tag-handler configuration for the
// CloudTrail tag operations, resolving trail resources by ARN.
func cloudTrailTagConfig(store cloudtrailstore.CloudTrailStoreInterface, mapErr func(error) error, requireTagKeys bool) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: func() tagutil.TagOperationConfig {
			c := tagutil.CloudTrailConfig
			c.RequireTagKeys = requireTagKeys
			return c
		}(),
		ResourceKey: func(arn string) string {
			trail, err := store.GetTrailByARN(arn)
			if err != nil {
				return ""
			}
			return trail.Name
		},
		ValidateResource: func(_ context.Context, resourceKey string) error {
			_, err := store.GetTrail(resourceKey)
			if err != nil {
				return mapErr(err)
			}
			return nil
		},
		TagFunc: func(_ context.Context, resourceKey string, tag []tagutil.Tag) error {
			if err := validateCloudTrailTags(tag); err != nil {
				return err
			}
			return store.Tag(resourceKey, tagutil.ToMap(tag))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Untag(resourceKey, tagKeys)
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			keys := tagutil.ParseTagKeysAsSlice(params, "TagKeyList")
			if len(keys) > 0 {
				return keys
			}
			return tagutil.ParseTagKeysWithKeyName(params, "TagsList", "Key")
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return store.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tag []tagutil.Tag, rawKey string) (interface{}, error) {
			return map[string]interface{}{
				"ResourceTagList": []map[string]interface{}{
					{
						"ResourceId": rawKey,
						"TagsList":   tagutil.ToResponse(tag),
					},
				},
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: cloudTrailMapError,
	}
}
