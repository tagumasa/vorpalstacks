package sqs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	"vorpalstacks/internal/utils/aws/types"
)

func sqsMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError:
		return ErrMissingParameter
	case *tagutil.MissingTagsError:
		return ErrMissingParameter
	case *tagutil.MissingTagKeysError:
		return ErrMissingParameter
	}
	return convertStoreError(err)
}

func sqsTagConfig(store sqsstore.SQSStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.SQSConfig,
		ValidateResource: func(ctx context.Context, rawKey string) error {
			_, err := store.GetQueue(rawKey)
			if err != nil {
				return convertStoreError(err)
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []types.Tag) error {
			return store.TagQueue(resourceKey, tagutil.ToMap(tags))
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.UntagQueue(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			m, err := store.ListQueueTags(resourceKey)
			if err != nil {
				return nil, err
			}
			return tagutil.MapToTags(m), nil
		},
		FormatResponse: func(tags []types.Tag, rawResourceKey string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ToMap(tags),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: sqsMapError,
	}
}

// TagQueue adds or overwrites tags on an SQS queue.
func (s *SQSService) TagQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, sqsTagConfig(store))
}

// UntagQueue removes the specified tags from an SQS queue.
func (s *SQSService) UntagQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, sqsTagConfig(store))
}

// ListQueueTags lists all tags assigned to an SQS queue.
func (s *SQSService) ListQueueTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, sqsTagConfig(store))
}
