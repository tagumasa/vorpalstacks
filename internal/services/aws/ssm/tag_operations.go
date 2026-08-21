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

// AddTagsToResource adds or overwrites tags on an SSM parameter.
func (s *SSMService) AddTagsToResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := validateResourceType(req); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, ssmTagConfig(store))
}

// RemoveTagsFromResource removes the specified tags from an SSM parameter.
func (s *SSMService) RemoveTagsFromResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := validateResourceType(req); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, ssmTagConfig(store))
}

// ListTagsForResource lists all tags assigned to an SSM parameter.
func (s *SSMService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := validateResourceType(req); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, ssmTagConfig(store))
}

// validateResourceType enforces the Smithy ResourceTypeForTagging enum. Only
// "Parameter" is implemented in this edge/on-prem platform; any other value
// is rejected with InvalidResourceType, matching the AWS contract.
func validateResourceType(req *request.ParsedRequest) error {
	if rt := req.GetParam("ResourceType"); rt != "" && rt != "Parameter" {
		return ErrInvalidResourceType
	}
	return nil
}
