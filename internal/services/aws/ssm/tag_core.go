package ssm

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

var ssmParamConfig tagutil.TagOperationConfig

func init() {
	ssmParamConfig = tagutil.StandardConfig
	ssmParamConfig.ResourceParam = "ResourceId"
}

func ssmMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError:
		return ErrInvalidParameterName
	case *tagutil.MissingTagsError:
		return ErrInvalidParameterName
	case *tagutil.MissingTagKeysError:
		return ErrInvalidParameterName
	}
	if errors.Is(err, ssmstore.ErrParameterNotFound) {
		return ErrParameterNotFound
	}
	return err
}

// ssmTagConfig builds the tag-engine configuration whose store closures are
// the single tag persistence path for both protocol planes.
func ssmTagConfig(store ssmstore.SSMStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: ssmParamConfig,
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			return store.AddTagsToResource(resourceKey, tagutil.ToMap(tags))
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.RemoveTagsFromResource(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			m, err := store.ListTagsForResource(resourceKey)
			if err != nil {
				return nil, err
			}
			return tagutil.MapToTags(m), nil
		},
		FormatResponse: func(tags []tagutil.Tag, rawResourceKey string) (interface{}, error) {
			return map[string]interface{}{
				"TagList": tagutil.MapToResponse(tagutil.ToMap(tags)),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: ssmMapError,
	}
}

// validateResourceType enforces the Smithy ResourceTypeForTagging contract of
// the three tag operations: ResourceType is a required member, and only
// "Parameter" is implemented in this edge/on-prem platform. A missing member
// is rejected with ValidationException (AWS rejects null required members)
// and any other value with InvalidResourceType.
func validateResourceType(req *request.ParsedRequest) error {
	rt := req.GetParam("ResourceType")
	if rt == "" {
		return ErrValidationException
	}
	if rt != "Parameter" {
		return ErrInvalidResourceType
	}
	return nil
}

// addTagsToResourceCore is the single entry point for SSM tag application:
// it enforces the ResourceTypeForTagging enum and delegates to the shared
// tag engine with the SSM parameter store closures.
func (s *SSMService) addTagsToResourceCore(ctx context.Context, store ssmstore.SSMStoreInterface, req *request.ParsedRequest) (interface{}, error) {
	if err := validateResourceType(req); err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, ssmTagConfig(store))
}

// removeTagsFromResourceCore is the single entry point for SSM tag removal.
func (s *SSMService) removeTagsFromResourceCore(ctx context.Context, store ssmstore.SSMStoreInterface, req *request.ParsedRequest) (interface{}, error) {
	if err := validateResourceType(req); err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, ssmTagConfig(store))
}

// listTagsForResourceCore is the single entry point for SSM tag listing.
func (s *SSMService) listTagsForResourceCore(ctx context.Context, store ssmstore.SSMStoreInterface, req *request.ParsedRequest) (interface{}, error) {
	if err := validateResourceType(req); err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, ssmTagConfig(store))
}
