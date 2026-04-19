package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	tagutil "vorpalstacks/internal/common/tags"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

func dynamodbMapError(err error) error {
	switch err.(type) {
	case *tags.MissingResourceError:
		return ErrInvalidParameter
	case *tags.MissingTagsError:
		return ErrInvalidParameter
	case *tags.MissingTagKeysError:
		return ErrInvalidParameter
	}
	return err
}

func dynamodbTagConfig(store dynamodbstore.DynamoDBStoreInterface) tags.TagHandlerConfig {
	return tags.TagHandlerConfig{
		Param: tags.TagOperationConfig{
			ResourceParam: "ResourceArn",
			TagsParam:     "Tags",
			TagKeysParam:  "TagKeys",
			TagKeyName:    "Key",
			TagValueName:  "Value",
		},
		ResourceKey: func(rawKey string) string {
			return svcarn.ParseTableARN(rawKey)
		},
		ValidateResource: func(_ context.Context, rawKey string) error {
			tableName := svcarn.ParseTableARN(rawKey)
			if tableName == "" {
				return ErrResourceNotFound
			}
			if _, err := store.Tables().Get(tableName); err != nil {
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, dynamodbTagConfig(store))
}

// UntagResource removes the specified tags from a DynamoDB table.
func (s *DynamoDBService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleUntag(ctx, req, dynamodbTagConfig(store))
}

// ListTagsForResource lists all tags assigned to a DynamoDB table.
func (s *DynamoDBService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleList(ctx, req, dynamodbTagConfig(store))
}
