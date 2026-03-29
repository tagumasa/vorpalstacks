package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState) (string, string, *ExecutionError) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

	if isJSONata && state.Arguments != nil {
		var inputData interface{}
		json.Unmarshal([]byte(processedInput), &inputData)
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		argsInput, err := e.applyJSONataArguments(ctx, state.Arguments, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Arguments", err.Error())
		}
		processedInput = argsInput
		execCtx.AfterArguments = &processedInput
	} else if state.Parameters != nil {
		processedInput = e.applyParameters(processedInput, state.Parameters)
	}

	timeoutSeconds := state.GetTimeoutSeconds()
	heartbeatSeconds := state.GetHeartbeatSeconds()

	if isJSONata {
		var inputData interface{}
		json.Unmarshal([]byte(processedInput), &inputData)
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		vars := buildVarsMap(statesVar, execCtx.VariableScope)

		if s, ok := state.TimeoutSeconds.(string); ok && IsExpression(s) {
			result, err := EvaluateJSONata(ctx, UnwrapExpression(s), nil, vars)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "TimeoutSeconds", err.Error())
			}
			if f, ok := toFloat64(result); ok {
				timeoutSeconds = int32(f)
			}
		}

		if s, ok := state.HeartbeatSeconds.(string); ok && IsExpression(s) {
			result, err := EvaluateJSONata(ctx, UnwrapExpression(s), nil, vars)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "HeartbeatSeconds", err.Error())
			}
			if f, ok := toFloat64(result); ok {
				heartbeatSeconds = int32(f)
			}
		}
	}

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
		execCtx.RetryCount = attempt - 1

		taskCtx := ctx
		var cancel context.CancelFunc
		if timeoutSeconds > 0 {
			taskCtx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		}

		if strings.HasPrefix(state.Resource, "arn:aws:lambda:") {
			output, taskErr = e.executeLambdaTask(taskCtx, execCtx, state, processedInput)
		} else if e.isActivityResource(state.Resource) {
			var next string
			output, next, taskErr = e.executeActivityTask(taskCtx, execCtx, state, processedInput, timeoutSeconds, heartbeatSeconds)
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

				if isJSONata {
					return e.executeTaskJSONataCatch(ctx, execCtx, state, processedInput, errorCode, cause, catchPolicy)
				}

				catchOutput := e.buildCatchOutput(processedInput, errorCode, cause, catchPolicy.ResultPath)
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

	if isJSONata {
		if state.JSONataOutput == nil && len(state.OutputRaw) > 0 {
			var err error
			state.JSONataOutput, err = resolveJSONataOutput(state)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
			}
		}

		var resultData interface{}
		json.Unmarshal([]byte(output), &resultData)
		var inputData interface{}
		json.Unmarshal([]byte(processedInput), &inputData)
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, resultData, nil)

		if len(state.Assign) > 0 {
			evaluated, err := evaluateAssign(ctx, state.Assign, statesVar, execCtx.VariableScope)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
			}
			execCtx.PendingAssign = evaluated
		}

		if state.JSONataOutput != nil {
			resolved, err := e.applyJSONataOutput(ctx, state.JSONataOutput, statesVar, execCtx.VariableScope)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
			}
			outputJSON, err := json.Marshal(resolved)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", fmt.Sprintf("failed to marshal: %s", err.Error()))
			}
			output = string(outputJSON)
		}
	} else {
		output = e.applyResultSelector(output, state.GetResultSelector())
		output = e.applyResultPath(processedInput, output, state.ResultPath)
		output = e.applyOutputPath(output, state.GetOutputPath())
	}

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

func (e *Executor) executeTaskJSONataCatch(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, processedInput, errorCode, cause string, catchPolicy *sfnstore.CatchPolicy) (string, string, *ExecutionError) {
	errorOutput := map[string]interface{}{
		"Error": errorCode,
		"Cause": cause,
	}

	var inputData interface{}
	json.Unmarshal([]byte(processedInput), &inputData)
	statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, errorOutput)

	if len(catchPolicy.Assign) > 0 {
		evaluated, err := evaluateAssign(ctx, catchPolicy.Assign, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Catch.Assign", err.Error())
		}
		execCtx.PendingAssign = evaluated
	}

	if catchPolicy.Output != nil {
		resolved, err := e.applyJSONataOutput(ctx, catchPolicy.Output, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Catch.Output", err.Error())
		}
		outputJSON, err := json.Marshal(resolved)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Catch.Output", fmt.Sprintf("failed to marshal: %s", err.Error()))
		}
		return string(outputJSON), catchPolicy.Next, nil
	}

	errorJSON, _ := json.Marshal(errorOutput)
	return string(errorJSON), catchPolicy.Next, nil
}
