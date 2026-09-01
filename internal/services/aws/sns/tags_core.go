package sns

import (
	"context"
	"errors"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	snsstore "vorpalstacks/internal/store/aws/sns"
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
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			// SNS tags topics only; the resource key is the topic ARN, and
			// the model attaches ResourceNotFoundException to every tag
			// operation, so an ARN addressing no existing topic is rejected
			// instead of silently persisting tags under an unowned key.
			if _, err := store.GetTopic(resourceKey); err != nil {
				if errors.Is(err, snsstore.ErrTopicNotFound) {
					return awserrors.NewAWSError("ResourceNotFoundException", "Topic does not exist", http.StatusNotFound)
				}
				return err
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			return store.Tag(resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return store.ListTagsForResource(resourceKey)
		},
		TagResponse: func(ctx context.Context, resourceKey string) (interface{}, error) {
			return tagutil.HandleListSimple(ctx, tagutil.StandardConfig, resourceKey,
				func(key string) ([]tagutil.Tag, error) { return store.ListTagsForResource(key) },
			)
		},
		MapError: snsMapError,
	}
}
