package sns

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	snsstore "vorpalstacks/internal/store/aws/sns"
	"vorpalstacks/internal/utils/aws/types"
)

func snsMapError(err error) error {
	switch e := err.(type) {
	case *tagutil.MissingResourceError:
		return awserrors.NewInvalidParameterException(e.Param + " is required")
	case *tagutil.MissingTagsError:
		return awserrors.NewInvalidParameterException(e.Param + " is required")
	case *tagutil.MissingTagKeysError:
		return awserrors.NewInvalidParameterException(e.Param + " is required")
	}
	return err
}

func snsTagConfig(store snsstore.SNSStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		TagFunc: func(ctx context.Context, resourceKey string, tags []types.Tag) error {
			return store.Tag(resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			return store.ListTagsForResource(resourceKey)
		},
		TagResponse: func(ctx context.Context, resourceKey string) (interface{}, error) {
			return tagutil.HandleListSimple(ctx, tagutil.StandardConfig, resourceKey,
				func(key string) ([]types.Tag, error) { return store.ListTagsForResource(key) },
			)
		},
		MapError: snsMapError,
	}
}

// TagResource adds or overwrites tags on an SNS topic.
func (s *SNSService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, snsTagConfig(store))
}

// UntagResource removes the specified tags from an SNS topic.
func (s *SNSService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, snsTagConfig(store))
}

// ListTagsForResource lists all tags assigned to an SNS topic.
func (s *SNSService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, snsTagConfig(store))
}
