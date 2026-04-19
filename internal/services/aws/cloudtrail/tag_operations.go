package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	"vorpalstacks/internal/utils/aws/types"
)

func cloudTrailMapError(err error) error {
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

func cloudTrailTagConfig(store cloudtrailstore.CloudTrailStoreInterface, mapErr func(error) error, requireTagKeys bool) tags.TagHandlerConfig {
	return tags.TagHandlerConfig{
		Param: func() tags.TagOperationConfig {
			c := tags.CloudTrailConfig
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
		ValidateResource: func(_ context.Context, rawKey string) error {
			_, err := store.GetTrailByARN(rawKey)
			if err != nil {
				return mapErr(err)
			}
			return nil
		},
		TagFunc: func(_ context.Context, resourceKey string, tag []types.Tag) error {
			return store.Tag(resourceKey, tags.ToMap(tag))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Untag(resourceKey, tagKeys)
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			keys := tags.ParseTagKeysAsSlice(params, "TagKeyList")
			if len(keys) > 0 {
				return keys
			}
			return tags.ParseTagKeysWithKeyName(params, "TagsList", "Key")
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			return store.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tag []types.Tag, rawKey string) (interface{}, error) {
			return map[string]interface{}{
				"ResourceTagList": []map[string]interface{}{
					{
						"ResourceId": rawKey,
						"TagsList":   tags.ToResponse(tag),
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

// AddTags adds or overwrites tags on a CloudTrail trail.
func (s *CloudTrailService) AddTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return tags.HandleTag(ctx, req, cloudTrailTagConfig(store, s.mapStoreError, true))
}

// RemoveTags removes the specified tags from a CloudTrail trail.
func (s *CloudTrailService) RemoveTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return tags.HandleUntag(ctx, req, cloudTrailTagConfig(store, s.mapStoreError, false))
}

// ListTags lists all tags assigned to a CloudTrail trail.
func (s *CloudTrailService) ListTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return tags.HandleList(ctx, req, cloudTrailTagConfig(store, s.mapStoreError, false))
}
