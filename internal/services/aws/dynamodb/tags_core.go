package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// dynamodbTagConfig builds the table-tagging configuration of the shared tag
// machinery: the store-backed resource validation and tag persistence the
// tag operations run through.
func dynamodbTagConfig(store dynamodbstore.DynamoDBStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam: "ResourceArn",
			TagsParam:     "Tags",
			TagKeysParam:  "TagKeys",
			TagKeyName:    "Key",
			TagValueName:  "Value",
		},
		ResourceKey: func(rawKey string) string {
			return svcarn.ParseTableARN(rawKey)
		},
		ValidateResource: func(_ context.Context, resourceKey string) error {
			if resourceKey == "" {
				return ErrResourceNotFound
			}
			if _, err := store.Tables().Get(resourceKey); err != nil {
				return ErrResourceNotFound
			}
			return nil
		},
		ParseTags: func(params map[string]interface{}) []tagutil.Tag {
			return tagutil.ParseTags(params, "Tags")
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			m := tagutil.ParseTagKeys(params, "TagKeys")
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			return keys
		},
		TagFunc: func(_ context.Context, resourceKey string, tag []tagutil.Tag) error {
			return store.Tables().Tags().Tag(resourceKey, tagutil.ToMap(tag))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Tables().Tags().Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return store.Tables().Tags().ListAsSlice(resourceKey)
		},
		FormatResponse: func(tags []tagutil.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ConvertToMapSlice(tags),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: dynamodbMapError,
	}
}
