package dynamodb

import (
	"context"
	"sort"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
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
		ParseTags: func(params map[string]interface{}) []types.Tag {
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
		TagFunc: func(_ context.Context, resourceKey string, tag []types.Tag) error {
			return store.Tables().Tags().Tag(resourceKey, tagutil.ToMap(tag))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Tables().Tags().Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			return store.Tables().Tags().ListAsSlice(resourceKey)
		},
		FormatResponse: func(tags []types.Tag, _ string) (interface{}, error) {
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

// TagResource adds or overwrites tags on a DynamoDB table.
func (s *DynamoDBService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

// UntagResource removes the specified tags from a DynamoDB table.
func (s *DynamoDBService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, dynamodbTagConfig(store))
}

// ListTagsForResource lists tags assigned to a DynamoDB table with pagination.
// AWS paginates tags; the marker is the tag key of the last returned entry.
func (s *DynamoDBService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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
