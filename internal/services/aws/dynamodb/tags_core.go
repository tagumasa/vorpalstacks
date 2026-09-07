package dynamodb

import (
	"context"
	"sort"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func dynamodbMapError(err error) error {
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

// tagResourceCore owns the TagResource validation and persistence: the
// resource-ARN shape, tag-key uniqueness, and the tag key/value character
// and length rules, before the shared tag machinery persists them.
func (s *DynamoDBService) tagResourceCore(ctx context.Context, req *request.ParsedRequest, store dynamodbstore.DynamoDBStoreInterface) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	tags := tagutil.ParseTags(req.Parameters, "Tags")
	if tagutil.HasDuplicateKeys(tags) {
		return nil, ErrInvalidParameter
	}
	for _, tag := range tags {
		if !validateTagKey(tag.Key) {
			return nil, ErrInvalidParameter
		}
		if !validateTagValue(tag.Value) {
			return nil, ErrInvalidParameter
		}
	}
	return tagutil.HandleTag(ctx, req, dynamodbTagConfig(store))
}

// untagResourceCore owns the UntagResource resource-ARN validation; the
// shared tag machinery resolves the resource and removes the keys.
func (s *DynamoDBService) untagResourceCore(ctx context.Context, req *request.ParsedRequest, store dynamodbstore.DynamoDBStoreInterface) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}
	return tagutil.HandleUntag(ctx, req, dynamodbTagConfig(store))
}

// listTagsForResourceCore owns the ListTagsForResource read and pagination:
// resource-ARN validation, resource resolution, and the key-ordered marker
// pagination AWS applies to tag listings.
func (s *DynamoDBService) listTagsForResourceCore(ctx context.Context, req *request.ParsedRequest, store dynamodbstore.DynamoDBStoreInterface) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	cfg := dynamodbTagConfig(store)
	rawKey := tagutil.GetResourceKey(req.Parameters, cfg.Param)
	if rawKey == "" {
		return nil, dynamodbMapError(&tagutil.MissingResourceError{Param: "ResourceArn"})
	}
	resourceKey := rawKey
	if cfg.ResourceKey != nil {
		resourceKey = cfg.ResourceKey(rawKey)
	}
	if cfg.ValidateResource != nil {
		if err := cfg.ValidateResource(ctx, resourceKey); err != nil {
			return nil, dynamodbMapError(err)
		}
	}

	allTags, err := cfg.ListFunc(ctx, resourceKey)
	if err != nil {
		return nil, err
	}

	// Sort tags for deterministic pagination.
	sort.Slice(allTags, func(i, j int) bool {
		return allTags[i].Key < allTags[j].Key
	})

	// Apply NextToken pagination (marker = tag key).
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	startIdx := 0
	if nextToken != "" {
		for i, t := range allTags {
			if t.Key == nextToken {
				startIdx = i + 1
				break
			}
		}
	}

	pageSize := listTagsForResourceDefaultPageSize
	remaining := allTags[startIdx:]
	hasMore := len(remaining) > pageSize
	if hasMore {
		remaining = remaining[:pageSize]
	}

	resp := map[string]interface{}{
		"Tags": tagutil.ConvertToMapSlice(remaining),
	}
	if hasMore && len(remaining) > 0 {
		resp["NextToken"] = remaining[len(remaining)-1].Key
	}

	return resp, nil
}

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
