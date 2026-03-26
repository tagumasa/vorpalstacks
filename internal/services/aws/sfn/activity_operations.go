package sfn

import (
	"context"
	"errors"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// GetActivityTask retrieves a task from an activity for a worker to process.
func (s *StepFunctionService) GetActivityTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	activityArn := request.GetParamLowerFirst(req.Parameters, "activityArn")
	workerName := request.GetParamLowerFirst(req.Parameters, "workerName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetActivity(ctx, activityArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return nil, NewActivityDoesNotExist("Activity Does not exist: " + activityArn)
		}
		return nil, err
	}

	task, err := store.GetActivityTask(ctx, activityArn, workerName)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return map[string]interface{}{
			"taskToken": "",
			"input":     "",
		}, nil
	}

	return map[string]interface{}{
		"taskToken": task.TaskToken,
		"input":     task.Input,
	}, nil
}

// SendTaskSuccess reports that a task completed successfully.
func (s *StepFunctionService) SendTaskSuccess(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskToken := request.GetParamLowerFirst(req.Parameters, "taskToken")
	output := request.GetParamLowerFirst(req.Parameters, "output")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.CompleteActivityTask(taskToken, output); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// SendTaskFailure reports that a task failed.
func (s *StepFunctionService) SendTaskFailure(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskToken := request.GetParamLowerFirst(req.Parameters, "taskToken")
	errorMsg := request.GetParamLowerFirst(req.Parameters, "error")
	cause := request.GetParamLowerFirst(req.Parameters, "cause")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.FailActivityTask(taskToken, errorMsg, cause); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// SendTaskHeartbeat reports that a task is still running.
func (s *StepFunctionService) SendTaskHeartbeat(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskToken := request.GetParamLowerFirst(req.Parameters, "taskToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.GetActivityTaskByToken(taskToken)
	if err != nil {
		return nil, err
	}

	if task == nil || task.Status != "RUNNING" {
		return nil, sfnstore.ErrTaskNotRunning
	}

	return response.EmptyResponse(), nil
}
