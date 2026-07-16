package cloudtrail

import (
	"context"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	"vorpalstacks/internal/utils/aws/types"
)

// CloudTrail tag limits per AWS spec (Smithy model).
const (
	maxCloudTrailTags = 50
	maxCTTagKeyLen    = 128
	maxCTTagValueLen  = 256
)

// validateCloudTrailTags validates tag count, key length, value length, and
// reserved prefix against AWS CloudTrail limits.
func validateCloudTrailTags(tagList []types.Tag) error {
	if len(tagList) > maxCloudTrailTags {
		return awserrors.NewAWSError("TagsLimitExceededException",
			fmt.Sprintf("Number of tags exceeds the limit of %d", maxCloudTrailTags), 400)
	}
	for _, t := range tagList {
		if strings.HasPrefix(t.Key, "aws:") {
			return awserrors.NewAWSError("InvalidTagKeyException",
				"Tag keys starting with 'aws:' are reserved", 400)
		}
		if len(t.Key) < 1 || len(t.Key) > maxCTTagKeyLen {
			return awserrors.NewAWSError("InvalidTagKeyException",
				fmt.Sprintf("Tag key length must be between 1 and %d", maxCTTagKeyLen), 400)
		}
		if len(t.Value) > maxCTTagValueLen {
			return awserrors.NewAWSError("InvalidTagValueException",
				fmt.Sprintf("Tag value length must not exceed %d", maxCTTagValueLen), 400)
		}
	}
	return nil
}

// validateCloudTrailTagMap validates a tag map for count, size, and reserved
// prefix limits.
func validateCloudTrailTagMap(tagMap map[string]string) error {
	if len(tagMap) > maxCloudTrailTags {
		return awserrors.NewAWSError("TagsLimitExceededException",
			fmt.Sprintf("Number of tags exceeds the limit of %d", maxCloudTrailTags), 400)
	}
	for k, v := range tagMap {
		if strings.HasPrefix(k, "aws:") {
			return awserrors.NewAWSError("InvalidTagKeyException",
				"Tag keys starting with 'aws:' are reserved", 400)
		}
		if len(k) < 1 || len(k) > maxCTTagKeyLen {
			return awserrors.NewAWSError("InvalidTagKeyException",
				fmt.Sprintf("Tag key length must be between 1 and %d", maxCTTagKeyLen), 400)
		}
		if len(v) > maxCTTagValueLen {
			return awserrors.NewAWSError("InvalidTagValueException",
				fmt.Sprintf("Tag value length must not exceed %d", maxCTTagValueLen), 400)
		}
	}
	return nil
}

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
		ValidateResource: func(_ context.Context, resourceKey string) error {
			_, err := store.GetTrail(resourceKey)
			if err != nil {
				return mapErr(err)
			}
			return nil
		},
		TagFunc: func(_ context.Context, resourceKey string, tag []types.Tag) error {
			if err := validateCloudTrailTags(tag); err != nil {
				return err
			}
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
