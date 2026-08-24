package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// This file holds the activity-operation Core path shared by the HTTP API
// and the admin console handler: activity lifecycle, the long-poll task
// worker API and the task report operations.

// CreateActivityInput carries every field that CreateActivity needs.
type CreateActivityInput struct {
	Name             string
	Tags             map[string]string
	EncryptionConfig *sfnstore.EncryptionConfiguration
}

// GetActivityTaskInput carries the parameters for GetActivityTask.
type GetActivityTaskInput struct {
	ActivityArn string
	WorkerName  string
}

// GetActivityTaskResult carries the GetActivityTaskOutput members; an
// empty task token reports that no task arrived within the poll window.
type GetActivityTaskResult struct {
	TaskToken string
	Input     string
}

// SendTaskSuccessInput carries the parameters for SendTaskSuccess.
type SendTaskSuccessInput struct {
	TaskToken string
	Output    string
}

// SendTaskFailureInput carries the parameters for SendTaskFailure.
type SendTaskFailureInput struct {
	TaskToken string
	Error     string
	Cause     string
}

// createActivityCore is the single entry point for CreateActivity: name
// validation, encryption-configuration validation, the tag quota with
// rollback on tagging failure, and the duplicate-name mapping to
// ActivityAlreadyExists all run here, mirroring the state machine
// creation Core.
func (s *StepFunctionService) createActivityCore(ctx context.Context, store *sfnstore.StepFunctionStore, in CreateActivityInput) (map[string]interface{}, error) {
	if err := validateResourceName(in.Name); err != nil {
		return nil, err
	}
	if in.EncryptionConfig != nil {
		if err := validateEncryptionConfiguration(in.EncryptionConfig); err != nil {
			return nil, err
		}
	}
	if len(in.Tags) > sfnstore.MaxTagsPerResource {
		return nil, NewTooManyTags(fmt.Sprintf("Too many tags: %d, maximum allowed %d", len(in.Tags), sfnstore.MaxTagsPerResource))
	}

	activity := &sfnstore.Activity{
		Name:                    in.Name,
		EncryptionConfiguration: in.EncryptionConfig,
	}
	if err := store.CreateActivity(ctx, activity); err != nil {
		if errors.Is(err, sfnstore.ErrActivityAlreadyExists) {
			return nil, NewActivityAlreadyExists("Activity already exists: " + in.Name)
		}
		return nil, err
	}

	// Apply tags through the tag store; a tagging failure rolls the
	// creation back to avoid an untagged orphan.
	if len(in.Tags) > 0 {
		if err := store.Tag(activity.ActivityArn, in.Tags); err != nil {
			_ = store.DeleteActivity(ctx, activity.ActivityArn)
			return nil, err
		}
	}

	return map[string]interface{}{
		"activityArn":  activity.ActivityArn,
		"creationDate": activity.CreationDate.Unix(),
	}, nil
}

// deleteActivityCore is the single entry point for DeleteActivity.
func (s *StepFunctionService) deleteActivityCore(ctx context.Context, store *sfnstore.StepFunctionStore, activityArn string) error {
	if err := validateArnRequired(activityArn, "activityArn"); err != nil {
		return err
	}
	if err := store.DeleteActivity(ctx, activityArn); err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return NewActivityDoesNotExist("Activity Does not exist: " + activityArn)
		}
		return err
	}
	return nil
}

// describeActivityCore is the single entry point for DescribeActivity.
func (s *StepFunctionService) describeActivityCore(ctx context.Context, store *sfnstore.StepFunctionStore, activityArn string) (map[string]interface{}, error) {
	if err := validateArnRequired(activityArn, "activityArn"); err != nil {
		return nil, err
	}
	activity, err := store.GetActivity(ctx, activityArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return nil, NewActivityDoesNotExist("Activity Does not exist: " + activityArn)
		}
		return nil, err
	}
	return activityToResponse(activity), nil
}

// listActivitiesCore is the single entry point for ListActivities. The
// ActivityListItem shape carries activityArn, name and creationDate only.
func (s *StepFunctionService) listActivitiesCore(ctx context.Context, store *sfnstore.StepFunctionStore, maxResults int32, nextToken string) (map[string]interface{}, error) {
	if err := validateMaxResults(maxResults, 0, sfnstore.MaxPageSize, "maxResults"); err != nil {
		return nil, err
	}
	result, err := store.ListActivities(ctx, maxResults, nextToken)
	if err != nil {
		return nil, err
	}

	activities := make([]map[string]interface{}, len(result.Activities))
	for i, activity := range result.Activities {
		activities[i] = map[string]interface{}{
			"activityArn":  activity.ActivityArn,
			"name":         activity.Name,
			"creationDate": activity.CreationDate.Unix(),
		}
	}

	response := map[string]interface{}{"activities": activities}
	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}
	return response, nil
}

// getActivityTaskCore is the single entry point for the worker long poll:
// it holds the request open for up to sixty seconds and returns an empty
// task token when no task arrives within that window.
func (s *StepFunctionService) getActivityTaskCore(ctx context.Context, store *sfnstore.StepFunctionStore, in GetActivityTaskInput) (*GetActivityTaskResult, error) {
	if err := validateArnRequired(in.ActivityArn, "activityArn"); err != nil {
		return nil, err
	}

	if _, err := store.GetActivity(ctx, in.ActivityArn); err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return nil, NewActivityDoesNotExist("Activity Does not exist: " + in.ActivityArn)
		}
		return nil, err
	}

	pollCtx, cancel := context.WithTimeout(ctx, sfnstore.ActivityTaskPollTimeout)
	defer cancel()

	task, err := store.GetActivityTask(pollCtx, in.ActivityArn, in.WorkerName)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return &GetActivityTaskResult{TaskToken: "", Input: ""}, nil
	}
	return &GetActivityTaskResult{TaskToken: task.TaskToken, Input: task.Input}, nil
}

// validateTaskToken enforces the TaskToken shape: required, at most 2048
// characters. The operation-attached error for token problems is
// InvalidToken.
func validateTaskToken(taskToken string) error {
	if taskToken == "" {
		return NewInvalidToken("taskToken is required")
	}
	if len(taskToken) > sfnstore.MaxTaskTokenLength {
		return NewInvalidToken(fmt.Sprintf("taskToken must be at most %d characters, got %d", sfnstore.MaxTaskTokenLength, len(taskToken)))
	}
	return nil
}

// sendTaskSuccessCore is the single entry point for SendTaskSuccess. The
// output must be valid JSON within the SensitiveData bound; unknown
// tokens map to TaskDoesNotExist and superseded ones to TaskTimedOut.
func (s *StepFunctionService) sendTaskSuccessCore(ctx context.Context, store *sfnstore.StepFunctionStore, in SendTaskSuccessInput) error {
	if err := validateTaskToken(in.TaskToken); err != nil {
		return err
	}
	if len(in.Output) > sfnstore.MaxExecutionDataBytes {
		return NewInvalidOutput(fmt.Sprintf("Invalid Output: output must be at most %d bytes, got %d", sfnstore.MaxExecutionDataBytes, len(in.Output)))
	}
	if !json.Valid([]byte(in.Output)) {
		return NewInvalidOutput("Invalid Output: output must be valid JSON")
	}

	if err := store.CompleteActivityTask(in.TaskToken, in.Output); err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotFound) {
			return ErrTaskDoesNotExist
		}
		if errors.Is(err, sfnstore.ErrTaskNotRunning) {
			// The task already reported, or its attempt timed out and a
			// newer attempt now holds a different token.
			return ErrTaskTimedOut
		}
		return err
	}
	return nil
}

// sendTaskFailureCore is the single entry point for SendTaskFailure; the
// error and cause strings obey the SensitiveError and SensitiveCause
// bounds.
func (s *StepFunctionService) sendTaskFailureCore(ctx context.Context, store *sfnstore.StepFunctionStore, in SendTaskFailureInput) error {
	if err := validateTaskToken(in.TaskToken); err != nil {
		return err
	}
	if len(in.Error) > sfnstore.MaxErrorLength {
		return NewValidationException(fmt.Sprintf("error must be at most %d characters, got %d", sfnstore.MaxErrorLength, len(in.Error)))
	}
	if len(in.Cause) > sfnstore.MaxCauseLength {
		return NewValidationException(fmt.Sprintf("cause must be at most %d characters, got %d", sfnstore.MaxCauseLength, len(in.Cause)))
	}

	if err := store.FailActivityTask(in.TaskToken, in.Error, in.Cause); err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotFound) {
			return ErrTaskDoesNotExist
		}
		if errors.Is(err, sfnstore.ErrTaskNotRunning) {
			return ErrTaskTimedOut
		}
		return err
	}
	return nil
}

// sendTaskHeartbeatCore is the single entry point for SendTaskHeartbeat.
func (s *StepFunctionService) sendTaskHeartbeatCore(ctx context.Context, store *sfnstore.StepFunctionStore, taskToken string) error {
	if err := validateTaskToken(taskToken); err != nil {
		return err
	}

	task, err := store.GetActivityTaskByToken(taskToken)
	if err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotFound) {
			return ErrTaskDoesNotExist
		}
		return err
	}
	if task == nil || task.Status != "RUNNING" {
		return ErrTaskDoesNotExist
	}

	if err := store.HeartbeatActivityTask(taskToken); err != nil {
		if errors.Is(err, sfnstore.ErrTaskNotFound) {
			return ErrTaskDoesNotExist
		}
		return err
	}
	return nil
}
