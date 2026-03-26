package tags

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
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

// ListTagsMapHandler is a function that lists tags for a resource as a map.
type ListTagsMapHandler func(resourceKey string) (map[string]string, error)

// TagOperationFuncs holds the handler functions for tag operations.
type TagOperationFuncs struct {
	Tag     TagResourceHandler
	Untag   UntagResourceHandler
	List    ListTagsHandler
	ListMap ListTagsMapHandler
}

// HandleTag processes a TagResource request using the configured handlers.
func (f *TagOperationFuncs) HandleTag(
	ctx context.Context,
	params map[string]interface{},
	config TagOperationConfig,
	onError func(msg string) error,
) (interface{}, error) {
	resourceKey := GetResourceKey(params, config)
	if resourceKey == "" && config.RequireResource {
		return nil, onError(config.ResourceParam + " is required")
	}

	tags := GetTags(params, config)
	if len(tags) == 0 && config.RequireTags {
		return nil, onError(config.TagsParam + " is required")
	}

	if len(tags) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := f.Tag(resourceKey, tags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// HandleUntag processes an UntagResource request using the configured handlers.
func (f *TagOperationFuncs) HandleUntag(
	ctx context.Context,
	params map[string]interface{},
	config TagOperationConfig,
	onError func(msg string) error,
) (interface{}, error) {
	resourceKey := GetResourceKey(params, config)
	if resourceKey == "" && config.RequireResource {
		return nil, onError(config.ResourceParam + " is required")
	}

	tagKeys := GetTagKeys(params, config)
	if len(tagKeys) == 0 && config.RequireTagKeys {
		return nil, onError(config.TagKeysParam + " is required")
	}

	if len(tagKeys) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := f.Untag(resourceKey, tagKeys); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// HandleList processes a ListTags request using the configured handlers.
func (f *TagOperationFuncs) HandleList(
	ctx context.Context,
	params map[string]interface{},
	config TagOperationConfig,
	onError func(msg string) error,
) (interface{}, error) {
	resourceKey := GetResourceKey(params, config)
	if resourceKey == "" && config.RequireResource {
		return nil, onError(config.ResourceParam + " is required")
	}

	if f.ListMap != nil {
		tags, err := f.ListMap(resourceKey)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"Tags": tags,
		}, nil
	}

	tags, err := f.List(resourceKey)
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

// HandleTagSimple processes a simple TagResource request without full error handling.
func HandleTagSimple(
	ctx context.Context,
	params map[string]interface{},
	config TagOperationConfig,
	resourceKey string,
	tagFunc TagResourceHandler,
) (interface{}, error) {
	tags := GetTags(params, config)
	if len(tags) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := tagFunc(resourceKey, tags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// HandleUntagSimple processes a simple UntagResource request without full error handling.
func HandleUntagSimple(
	ctx context.Context,
	params map[string]interface{},
	config TagOperationConfig,
	resourceKey string,
	untagFunc UntagResourceHandler,
) (interface{}, error) {
	tagKeys := GetTagKeys(params, config)
	if len(tagKeys) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := untagFunc(resourceKey, tagKeys); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

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

// HandleListMapSimple processes a simple ListTags request and returns tags as a map.
func HandleListMapSimple(
	resourceKey string,
	listFunc ListTagsMapHandler,
) (interface{}, error) {
	tags, err := listFunc(resourceKey)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Tags": tags,
	}, nil
}
