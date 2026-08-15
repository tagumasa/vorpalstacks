package sfn

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// GetActivityTask retrieves a task from an activity for a worker to process.
func (s *StepFunctionService) GetActivityTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	activityArn := request.GetParamLowerFirst(req.Parameters, "activityArn")
	workerName := request.GetParamLowerFirst(req.Parameters, "workerName")
	if err := validateArnRequired(activityArn, "activityArn"); err != nil {
		return nil, err
	}

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

	// GetActivityTask is a long poll: the Step Functions API holds the
	// request open for up to ActivityTaskPollTimeout (60 seconds per the
	// API reference) and returns an empty taskToken when no task arrives
	// within that window, instead of blocking until the client gives up.
	pollCtx, cancel := context.WithTimeout(ctx, sfnstore.ActivityTaskPollTimeout)
	defer cancel()

	task, err := store.GetActivityTask(pollCtx, activityArn, workerName)
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
	if taskToken == "" {
		return nil, NewInvalidToken("taskToken is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.CompleteActivityTask(taskToken, output); err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotRunning) {
			// The task already reported, or its attempt timed out and a
			// newer attempt now holds a different token.
			return nil, ErrTaskTimedOut
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// SendTaskFailure reports that a task failed.
func (s *StepFunctionService) SendTaskFailure(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskToken := request.GetParamLowerFirst(req.Parameters, "taskToken")
	if taskToken == "" {
		return nil, NewInvalidToken("taskToken is required")
	}
	errorMsg := request.GetParamLowerFirst(req.Parameters, "error")
	cause := request.GetParamLowerFirst(req.Parameters, "cause")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.FailActivityTask(taskToken, errorMsg, cause); err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotFound) {
			return nil, ErrTaskDoesNotExist
		}
		if errors.Is(err, sfnstore.ErrTaskNotRunning) {
			return nil, ErrTaskTimedOut
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// SendTaskHeartbeat reports that a task is still running.
func (s *StepFunctionService) SendTaskHeartbeat(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskToken := request.GetParamLowerFirst(req.Parameters, "taskToken")
	if taskToken == "" {
		return nil, NewInvalidToken("taskToken is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.GetActivityTaskByToken(taskToken)
	if err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotFound) {
			return nil, ErrTaskDoesNotExist
		}
		return nil, err
	}

	if task == nil || task.Status != "RUNNING" {
		return nil, ErrTaskDoesNotExist
	}

	if err := store.HeartbeatActivityTask(taskToken); err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotFound) {
			return nil, ErrTaskDoesNotExist
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// CreateActivity creates a new activity.
func (s *StepFunctionService) CreateActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "name")

	if err := validateResourceName(name); err != nil {
		return nil, err
	}

	activity := &sfnstore.Activity{
		Name: name,
	}

	if ec, err := parseEncryptionConfigurationFromJSON(req.Parameters["encryptionConfiguration"]); err != nil {
		return nil, err
	} else {
		activity.EncryptionConfiguration = ec
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags"))
	if len(tags) > 0 {
		if err := store.Tag(activity.ActivityArn, tags); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"activityArn":  activity.ActivityArn,
		"creationDate": activity.CreationDate.Unix(),
	}, nil
}

// DeleteActivity deletes an activity.
func (s *StepFunctionService) DeleteActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "activityArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteActivity(ctx, arn); err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return nil, NewActivityDoesNotExist("Activity Does not exist: " + arn)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeActivity returns the details of an activity.
func (s *StepFunctionService) DescribeActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "activityArn")
	if err := validateArnRequired(arn, "activityArn"); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	activity, err := store.GetActivity(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return nil, NewActivityDoesNotExist("Activity Does not exist: " + arn)
		}
		return nil, err
	}

	return activityToResponse(activity), nil
}

// GetActivity returns the details of an activity (alias for DescribeActivity).
func (s *StepFunctionService) GetActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeActivity(ctx, reqCtx, req)
}

// ListActivities returns a list of activities.
func (s *StepFunctionService) ListActivities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListActivities(ctx, limit, nextToken)
	if err != nil {
		return nil, err
	}

	activities := make([]map[string]interface{}, len(result.Activities))
	for i, activity := range result.Activities {
		activities[i] = activityToResponse(activity)
	}

	response := map[string]interface{}{
		"activities": activities,
	}

	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}

	return response, nil
}
