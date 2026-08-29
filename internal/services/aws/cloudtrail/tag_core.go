package cloudtrail

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

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
		return awserrors.NewAWSError("TagsLimitExceededException",
			fmt.Sprintf("Number of tags exceeds the limit of %d", tagutil.MaxTagsPerResource), 400)
	case tagutil.ReservedTagKey:
		return awserrors.NewAWSError("InvalidTagKeyException",
			"Tag keys starting with 'aws:' are reserved", 400)
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return awserrors.NewAWSError("InvalidTagKeyException",
			fmt.Sprintf("Tag key length must be between 1 and %d", tagutil.MaxTagKeyLength), 400)
	case tagutil.TagValueTooLong:
		return awserrors.NewAWSError("InvalidTagValueException",
			fmt.Sprintf("Tag value length must not exceed %d", tagutil.MaxTagValueLength), 400)
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
