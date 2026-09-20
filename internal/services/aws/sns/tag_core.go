package sns

import (
	"context"
	"errors"
	"fmt"

	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

func snsMapError(err error) error {
	switch e := err.(type) {
	case *tagutil.MissingResourceError:
		return NewInvalidParameter(e.Param + " is required")
	case *tagutil.MissingTagsError:
		return NewInvalidParameter(e.Param + " is required")
	case *tagutil.MissingTagKeysError:
		return NewInvalidParameter(e.Param + " is required")
	}
	return err
}

// validateTopicTagSlice enforces the AWS-wide tag constraints on the tag
// operations' typed list form (the shared framework's ValidateTagsFunc
// hook): at most fifty tags, keys of 1-128 characters not starting with
// the reserved aws: prefix, values of at most 256 characters. The
// fifty-tag total of the MERGED set is the tag store's budget (wired at
// construction); this check bounds the request itself and the per-tag
// shapes the store does not see.
func validateTopicTagSlice(tags []tagutil.Tag) error {
	switch v, offender := tagutil.CheckTags(tags, tagutil.StandardLimits()); v {
	case tagutil.TooManyTags:
		return ErrTagLimitExceeded
	case tagutil.ReservedTagKey:
		return NewInvalidParameter(fmt.Sprintf("Invalid parameter: Tag key %q: tag keys starting with aws: are reserved", offender))
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return NewInvalidParameter(fmt.Sprintf("Invalid parameter: Tag key %q: tag keys must be between 1 and 128 characters", offender))
	case tagutil.TagValueTooLong:
		return NewInvalidParameter(fmt.Sprintf("Invalid parameter: Tag value for key %q: tag values must be at most 256 characters", offender))
	}
	return nil
}

// validateTopicTagMap enforces the same tag constraints on the map form
// CreateTopic receives in its Tags member.
func validateTopicTagMap(tags map[string]string) error {
	if len(tags) == 0 {
		return nil
	}
	return validateTopicTagSlice(tagutil.MapToTags(tags))
}

// validateTopicTagKeys enforces the tag-key shapes on the UntagResource
// removal list through the shared checker's key-only form — the same
// constraints validateTopicTagKeys' slice sibling derives from CheckTags,
// with no count bound because the fifty-tag budget governs the stored set,
// not a removal request.
func validateTopicTagKeys(tagKeys []string) error {
	switch v, offender := tagutil.CheckTagKeys(tagKeys, tagutil.StandardLimits()); v {
	case tagutil.ReservedTagKey:
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid parameter: Tag key %q: tag keys starting with aws: are reserved", offender))
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid parameter: Tag key %q: tag keys must be between 1 and 128 characters", offender))
	}
	return nil
}

func snsTagConfig(store snsstore.SNSStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			// The tag operations' ResourceArn member targets the model's
			// AmazonResourceName shape, "@length": {"min": 1, "max": 1011}.
			if len(resourceKey) > snsstore.MaxResourceArnLength {
				return NewInvalidParameter(fmt.Sprintf(
					"Invalid parameter: ResourceArn must be at most %d characters", snsstore.MaxResourceArnLength))
			}
			// SNS tags topics only; the resource key is the topic ARN, and
			// the model attaches ResourceNotFoundException (awsQueryError
			// code "ResourceNotFound") to every tag operation, so an ARN
			// addressing no existing topic is rejected with that code —
			// not silently persisted under an unowned key.
			if _, err := store.GetTopic(resourceKey); err != nil {
				if errors.Is(err, snsstore.ErrTopicNotFound) {
					return ErrResourceNotFound
				}
				return mapStoreError(err)
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
		// TagResource's response is empty — the model's TagResourceResponse
		// shape carries no member, so a successful tag write returns no
		// body (the tags themselves are read back through
		// ListTagsForResource).
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		ValidateTagsFunc:    validateTopicTagSlice,
		ValidateTagKeysFunc: validateTopicTagKeys,
		MapError:            snsMapError,
	}
}
