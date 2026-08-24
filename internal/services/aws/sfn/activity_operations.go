package sfn

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// GetActivityTask retrieves a task from an activity for a worker to process.
func (s *StepFunctionService) GetActivityTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.getActivityTaskCore(ctx, store, GetActivityTaskInput{
		ActivityArn: request.GetParamLowerFirst(req.Parameters, "activityArn"),
		WorkerName:  request.GetParamLowerFirst(req.Parameters, "workerName"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"taskToken": result.TaskToken,
		"input":     result.Input,
	}, nil
}

// SendTaskSuccess reports that a task completed successfully.
func (s *StepFunctionService) SendTaskSuccess(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.sendTaskSuccessCore(ctx, store, SendTaskSuccessInput{
		TaskToken: request.GetParamLowerFirst(req.Parameters, "taskToken"),
		Output:    request.GetParamLowerFirst(req.Parameters, "output"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// SendTaskFailure reports that a task failed.
func (s *StepFunctionService) SendTaskFailure(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.sendTaskFailureCore(ctx, store, SendTaskFailureInput{
		TaskToken: request.GetParamLowerFirst(req.Parameters, "taskToken"),
		Error:     request.GetParamLowerFirst(req.Parameters, "error"),
		Cause:     request.GetParamLowerFirst(req.Parameters, "cause"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// SendTaskHeartbeat reports that a task is still running.
func (s *StepFunctionService) SendTaskHeartbeat(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.sendTaskHeartbeatCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "taskToken")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// CreateActivity creates a new activity.
func (s *StepFunctionService) CreateActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	encryptionConfig, err := parseEncryptionConfigurationFromJSON(req.Parameters["encryptionConfiguration"])
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.createActivityCore(ctx, store, CreateActivityInput{
		Name:             request.GetParamLowerFirst(req.Parameters, "name"),
		Tags:             tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags")),
		EncryptionConfig: encryptionConfig,
	})
}

// DeleteActivity deletes an activity.
func (s *StepFunctionService) DeleteActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteActivityCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "activityArn")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DescribeActivity returns the details of an activity.
func (s *StepFunctionService) DescribeActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeActivityCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "activityArn"))
}

// ListActivities returns a list of activities.
func (s *StepFunctionService) ListActivities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listActivitiesCore(ctx, store, limit, request.GetParamLowerFirst(req.Parameters, "nextToken"))
}
