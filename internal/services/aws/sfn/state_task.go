package sfn

import (
	"context"
	"fmt"
	"strings"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState) (string, string, *ExecutionError) {
	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())
	processedInput = e.applyParameters(processedInput, state.Parameters)

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "TaskStarted",
		Timestamp:       time.Now().UTC(),
		TaskStartedEventDetails: &sfnstore.TaskStartedEventDetails{
			Resource: state.Resource,
			Input:    processedInput,
		},
	})

	var output string
	var taskErr error
	var errorCode, cause string

	maxAttempts := int32(1)
	var retryPolicies []*sfnstore.RetryPolicy
	if len(state.Retry) > 0 {
		retryPolicies = state.Retry
		maxAttempts = 3
		for _, rp := range retryPolicies {
			if rp.MaxAttempts > maxAttempts {
				maxAttempts = rp.MaxAttempts
			}
		}
	}

	attempt := int32(0)

	for attempt < maxAttempts {
		attempt++

		taskCtx := ctx
		var cancel context.CancelFunc
		if state.TimeoutSeconds > 0 {
			taskCtx, cancel = context.WithTimeout(ctx, time.Duration(state.TimeoutSeconds)*time.Second)
		}

		if strings.HasPrefix(state.Resource, "arn:aws:lambda:") {
			output, taskErr = e.executeLambdaTask(taskCtx, execCtx, state, processedInput)
		} else if e.isActivityResource(state.Resource) {
			var next string
			output, next, taskErr = e.executeActivityTask(taskCtx, execCtx, state, processedInput)
			if taskErr == nil && next != "" {
				state.Next = next
			}
		} else if strings.HasPrefix(state.Resource, "arn:aws:states:::sqs:") {
			output, taskErr = e.executeSQSTask(taskCtx, execCtx, state, processedInput)
		} else if strings.HasPrefix(state.Resource, "arn:aws:states:::sns:") {
			output, taskErr = e.executeSNSTask(taskCtx, execCtx, state, processedInput)
		} else if strings.HasPrefix(state.Resource, "arn:aws:states:::events:") {
			output, taskErr = e.executeEventsTask(taskCtx, execCtx, state, processedInput)
		} else {
			taskErr = fmt.Errorf("unsupported resource type: %s", state.Resource)
		}

		if taskErr == nil {
			if cancel != nil {
				cancel()
			}
			break
		}

		if taskCtx.Err() == context.DeadlineExceeded || strings.Contains(taskErr.Error(), "States.Timeout") {
			errorCode = "States.Timeout"
			cause = "Task timed out"
		} else {
			errorCode = "States.TaskFailed"
			cause = taskErr.Error()
		}

		if attempt < maxAttempts && len(retryPolicies) > 0 {
			matchedRetry := e.findMatchingRetryPolicy(retryPolicies, errorCode)
			if matchedRetry != nil {
				interval := e.calculateBackoffInterval(matchedRetry, attempt)
				eventId = execCtx.nextEventId()
				e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
					ExecutionArn:    execCtx.Execution.ExecutionArn,
					EventId:         eventId,
					PreviousEventId: eventId - 1,
					Type:            "TaskRetried",
					Timestamp:       time.Now().UTC(),
				})

				if cancel != nil {
					cancel()
				}
				timer := time.NewTimer(interval)
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return "", "", &ExecutionError{ErrorCode: "States.Timeout", Cause: "Execution interrupted during retry"}
				}
				continue
			}
		}

		if cancel != nil {
			cancel()
		}
		break
	}

	if taskErr != nil {
		if len(state.Catch) > 0 {
			catchPolicy := e.findMatchingCatchPolicy(state.Catch, errorCode)
			if catchPolicy != nil {
				catchOutput := e.buildCatchOutput(processedInput, errorCode, cause, catchPolicy.ResultPath)
				eventId = execCtx.nextEventId()
				e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
					ExecutionArn:    execCtx.Execution.ExecutionArn,
					EventId:         eventId,
					PreviousEventId: eventId - 1,
					Type:            "TaskFailed",
					Timestamp:       time.Now().UTC(),
					TaskFailedEventDetails: &sfnstore.TaskFailedEventDetails{
						Resource: state.Resource,
						Error:    errorCode,
						Cause:    cause,
					},
				})
				return catchOutput, catchPolicy.Next, nil
			}
		}

		eventId = execCtx.nextEventId()
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    execCtx.Execution.ExecutionArn,
			EventId:         eventId,
			PreviousEventId: eventId - 1,
			Type:            "TaskFailed",
			Timestamp:       time.Now().UTC(),
			TaskFailedEventDetails: &sfnstore.TaskFailedEventDetails{
				Resource: state.Resource,
				Error:    errorCode,
				Cause:    cause,
			},
		})
		return "", "", &ExecutionError{ErrorCode: errorCode, Cause: cause}
	}

	output = e.applyResultSelector(output, state.GetResultSelector())
	output = e.applyResultPath(processedInput, output, state.ResultPath)
	output = e.applyOutputPath(output, state.GetOutputPath())

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "TaskSucceeded",
		Timestamp:       time.Now().UTC(),
		TaskSucceededEventDetails: &sfnstore.TaskSucceededEventDetails{
			Resource: state.Resource,
			Output:   output,
		},
	})

	return output, state.Next, nil
}
