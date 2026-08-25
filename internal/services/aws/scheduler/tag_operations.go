package scheduler

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

func schedulerMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError, *tagutil.MissingTagsError, *tagutil.MissingTagKeysError:
		return ErrValidation
	}
	return err
}

// TagResource adds or overwrites tags on an EventBridge Scheduler schedule
// group.
func (s *SchedulerService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleTag(ctx, req, s.scheduleGroupTagConfig(reqCtx))
}

// UntagResource removes the specified tags from an EventBridge Scheduler
// schedule group.
func (s *SchedulerService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleUntag(ctx, req, s.scheduleGroupTagConfig(reqCtx))
}

// ListTagsForResource lists all tags assigned to an EventBridge Scheduler
// schedule group.
func (s *SchedulerService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleList(ctx, req, s.scheduleGroupTagConfig(reqCtx))
}
