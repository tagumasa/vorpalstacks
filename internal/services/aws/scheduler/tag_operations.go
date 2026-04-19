package scheduler

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	"vorpalstacks/internal/utils/aws/types"
)

func schedulerMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError, *tagutil.MissingTagsError, *tagutil.MissingTagKeysError:
		return ErrValidation
	}
	return err
}

func schedulerTagConfig(store *schedulerstore.SchedulerStore) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		TagFunc: func(ctx context.Context, resourceKey string, tags []types.Tag) error {
			if err := store.TagFromSlice(resourceKey, tags); err != nil {
				logs.Debug("Failed to tag resource", logs.String("arn", resourceKey), logs.String("error", err.Error()))
				return ErrInternalServer
			}
			return nil
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			if err := store.Untag(resourceKey, tagKeys); err != nil {
				logs.Debug("Failed to untag resource", logs.String("arn", resourceKey), logs.String("error", err.Error()))
				return ErrInternalServer
			}
			return nil
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			tags, err := store.ListAsSlice(resourceKey)
			if err != nil {
				logs.Debug("Failed to list tags", logs.String("arn", resourceKey), logs.String("error", err.Error()))
				return nil, ErrInternalServer
			}
			return tags, nil
		},
		EmptyResponse: func() (interface{}, error) { return response.EmptyResponse(), nil },
		MapError:      schedulerMapError,
	}
}

// TagResource adds or overwrites tags on an EventBridge Scheduler schedule.
func (s *SchedulerService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, schedulerTagConfig(store))
}

// UntagResource removes the specified tags from an EventBridge Scheduler schedule.
func (s *SchedulerService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, schedulerTagConfig(store))
}

// ListTagsForResource lists all tags assigned to an EventBridge Scheduler schedule.
func (s *SchedulerService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, schedulerTagConfig(store))
}
