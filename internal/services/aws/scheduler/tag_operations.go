package scheduler

import (
	"context"
	"strings"

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
		// ValidateResource ensures the ARN refers to an existing schedule
		// or schedule group before tag operations proceed (M12). Without
		// this, tags could be applied to non-existent resources.
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			// Scheduler ARNs are formatted as:
			//   arn:aws:scheduler:<region>:<account>:schedule/<group>/<name>
			//   arn:aws:scheduler:<region>:<account>:schedule-group/<name>
			if strings.Contains(resourceKey, ":schedule-group/") {
				groupName := extractResourceName(resourceKey, "schedule-group/")
				if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
					return ErrScheduleGroupNotFound
				}
			} else if strings.Contains(resourceKey, ":schedule/") {
				groupName, schedName := extractScheduleNames(resourceKey)
				if _, err := store.GetSchedule(ctx, groupName, schedName); err != nil {
					return ErrScheduleNotFound
				}
			}
			return nil
		},
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

// extractResourceName extracts the resource name from an ARN after the
// given prefix (e.g. "schedule-group/").
func extractResourceName(arn, prefix string) string {
	idx := strings.Index(arn, prefix)
	if idx < 0 {
		return ""
	}
	return arn[idx+len(prefix):]
}

// extractScheduleNames extracts the group name and schedule name from a
// schedule ARN: ...:schedule/<group>/<name>
func extractScheduleNames(arn string) (groupName, scheduleName string) {
	idx := strings.Index(arn, ":schedule/")
	if idx < 0 {
		return "", ""
	}
	rest := arn[idx+len(":schedule/"):]
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) >= 1 {
		groupName = parts[0]
	}
	if groupName == "" {
		groupName = "default"
	}
	if len(parts) >= 2 {
		scheduleName = parts[1]
	}
	return groupName, scheduleName
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
