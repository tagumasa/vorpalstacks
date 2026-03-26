package scheduler

import (
	"context"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
)

// TagResource adds tags to an EventBridge Scheduler schedule.
func (s *SchedulerService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		if name := request.GetStringParam(req.Parameters, "Name"); name != "" {
			store, err := s.store(reqCtx)
			if err != nil {
				return nil, err
			}
			resourceArn = store.BuildScheduleARNFromName(name)
		}
	}

	if resourceArn == "" {
		return nil, ErrValidation
	}

	tags := tagutil.GetTags(req.Parameters, tagutil.StandardConfig)
	if len(tags) == 0 {
		return nil, ErrValidation
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.TagResourceFromSlice(resourceArn, tags); err != nil {
		logs.Debug("Failed to tag resource", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from an EventBridge Scheduler schedule.
func (s *SchedulerService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, ErrValidation
	}

	tagKeys := tagutil.GetTagKeys(req.Parameters, tagutil.StandardConfig)
	if len(tagKeys) == 0 {
		return nil, ErrValidation
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.UntagResource(resourceArn, tagKeys); err != nil {
		logs.Debug("Failed to untag resource", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for an EventBridge Scheduler schedule.
func (s *SchedulerService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := tagutil.GetResourceKey(req.Parameters, tagutil.StandardConfig)
	if resourceArn == "" {
		return nil, ErrValidation
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tagList, err := store.ListTagsAsSlice(resourceArn)
	if err != nil {
		logs.Debug("Failed to list tags", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	return map[string]interface{}{
		"Tags": tagutil.ToResponse(tagList),
	}, nil
}
