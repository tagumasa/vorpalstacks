package tags

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/aws/types"
)

// TagOperationConfig defines the configuration for tag operations.
type TagOperationConfig struct {
	ResourceParam      string
	TagsParam          string
	TagKeysParam       string
	TagKeyName         string
	TagValueName       string
	ResourceParamAlt   string
	RequireTags        bool
	RequireTagKeys     bool
	RequireResource    bool
	UseQueryFallback   bool
	CaseInsensitiveRes bool
}

// StandardConfig is the default tag operation configuration.
var StandardConfig = TagOperationConfig{
	ResourceParam:    "ResourceArn",
	TagsParam:        "Tags",
	TagKeysParam:     "TagKeys",
	TagKeyName:       "Key",
	TagValueName:     "Value",
	RequireTags:      true,
	RequireTagKeys:   true,
	RequireResource:  true,
	UseQueryFallback: true,
}

// StandardARNConfig is the tag operation configuration for ARN-based resources.
var StandardARNConfig = TagOperationConfig{
	ResourceParam:    "ResourceARN",
	TagsParam:        "Tags",
	TagKeysParam:     "TagKeys",
	TagKeyName:       "Key",
	TagValueName:     "Value",
	RequireTags:      true,
	RequireTagKeys:   true,
	RequireResource:  true,
	UseQueryFallback: false,
}

// KMSConfig is the tag operation configuration for KMS keys.
var KMSConfig = TagOperationConfig{
	ResourceParam:    "KeyId",
	TagsParam:        "Tags",
	TagKeysParam:     "TagKeys",
	TagKeyName:       "TagKey",
	TagValueName:     "TagValue",
	RequireTags:      false,
	RequireTagKeys:   false,
	RequireResource:  true,
	UseQueryFallback: false,
}

// LambdaConfig is the tag operation configuration for Lambda functions.
var LambdaConfig = TagOperationConfig{
	ResourceParam:    "Resource",
	TagsParam:        "Tags",
	TagKeysParam:     "TagKeys",
	TagKeyName:       "Key",
	TagValueName:     "Value",
	RequireTags:      true,
	RequireTagKeys:   true,
	RequireResource:  true,
	UseQueryFallback: false,
}

// SQSConfig is the tag operation configuration for SQS queues.
var SQSConfig = TagOperationConfig{
	ResourceParam:      "QueueUrl",
	TagsParam:          "Tags",
	TagKeysParam:       "TagKeys",
	TagKeyName:         "Key",
	TagValueName:       "Value",
	RequireTags:        false,
	RequireTagKeys:     false,
	RequireResource:    true,
	UseQueryFallback:   true,
	CaseInsensitiveRes: true,
}

// KinesisStreamConfig is the tag operation configuration for Kinesis streams.
var KinesisStreamConfig = TagOperationConfig{
	ResourceParam:    "StreamName",
	ResourceParamAlt: "StreamARN",
	TagsParam:        "Tags",
	TagKeysParam:     "TagKeys",
	TagKeyName:       "Key",
	TagValueName:     "Value",
	RequireTags:      true,
	RequireTagKeys:   true,
	RequireResource:  true,
	UseQueryFallback: false,
}

// CloudTrailConfig is the tag operation configuration for CloudTrail trails.
var CloudTrailConfig = TagOperationConfig{
	ResourceParam:    "ResourceId",
	ResourceParamAlt: "ResourceIdList",
	TagsParam:        "TagsList",
	TagKeysParam:     "TagKeyList",
	TagKeyName:       "Key",
	TagValueName:     "Value",
	RequireTags:      true,
	RequireTagKeys:   true,
	RequireResource:  true,
	UseQueryFallback: false,
}

// GetResourceKey extracts the resource key from parameters using the given configuration.
func GetResourceKey(params map[string]interface{}, config TagOperationConfig) string {
	if config.CaseInsensitiveRes {
		return request.GetParamCaseInsensitive(params, config.ResourceParam)
	}
	key := request.GetParamLowerFirst(params, config.ResourceParam)
	if key == "" && config.ResourceParamAlt != "" {
		if altList, ok := params[config.ResourceParamAlt].([]interface{}); ok && len(altList) > 0 {
			if s, ok := altList[0].(string); ok {
				return s
			}
		}
	}
	return key
}

// GetTags extracts tags from parameters using the given configuration.
func GetTags(params map[string]interface{}, config TagOperationConfig) []types.Tag {
	if config.UseQueryFallback {
		tags := ParseTagsWithQueryFallback(params, config.TagsParam)
		if len(tags) == 0 {
			tags = ParseTagsWithPrefix(params, "Tag")
		}
		return tags
	}
	return ParseTagsWithKeyNames(params, config.TagsParam, config.TagKeyName, config.TagValueName)
}

// GetTagKeys extracts tag keys from parameters using the given configuration.
func GetTagKeys(params map[string]interface{}, config TagOperationConfig) []string {
	if config.UseQueryFallback {
		return ParseTagKeysWithQueryFallback(params, config.TagKeysParam)
	}
	return ParseTagKeysAsSlice(params, config.TagKeysParam)
}

// TagResourceHandler is a function that tags a resource.
type TagResourceHandler func(resourceKey string, tags []types.Tag) error

// UntagResourceHandler is a function that untags a resource.
type UntagResourceHandler func(resourceKey string, tagKeys []string) error

// ListTagsHandler is a function that lists tags for a resource.
type ListTagsHandler func(resourceKey string) ([]types.Tag, error)

// HandleListSimple processes a simple ListTags request without full error handling.
func HandleListSimple(
	ctx context.Context,
	config TagOperationConfig,
	resourceKey string,
	listFunc ListTagsHandler,
) (interface{}, error) {
	tags, err := listFunc(resourceKey)
	if err != nil {
		return nil, err
	}

	var tagResponse []map[string]interface{}
	if config.TagKeyName == "Key" && config.TagValueName == "Value" {
		tagResponse = ToResponse(tags)
	} else {
		tagResponse = ToResponseWithKeyNames(tags, config.TagKeyName, config.TagValueName)
	}

	return map[string]interface{}{
		config.TagsParam: tagResponse,
	}, nil
}
