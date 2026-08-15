package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

func (e *Executor) executeTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState) (string, string, *ExecutionError) {
	// The attempt token lives on the execution context so every evaluation
	// dialect sees the same value; it is cleared on exit so later states'
	// context objects carry no stale Task section.
	defer func() { execCtx.TaskToken = "" }()

	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)
	baseInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

	// evaluateTaskInput applies Arguments/Parameters for a single attempt.
	// It runs per attempt rather than once per state because the attempt's
	// token and $$.State.RetryCount must reflect the attempt being
	// dispatched, and each activity attempt is stored under its own token.
	// The token reaches the payload through both dialects: JSONPath
	// Parameters resolve $$.Task.Token from the taskToken parameter, and
	// JSONata Arguments resolve $states.context.Task.Token from the context
	// object built with that token. taskToken is empty for non-activity
	// resources, where any Task.Token reference fails the evaluation.
	evaluateTaskInput := func(taskToken string) (string, *ExecutionError) {
		if isJSONata && state.Arguments != nil {
			var inputData interface{}
			if err := json.Unmarshal([]byte(baseInput), &inputData); err != nil {
				return "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
			}
			statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
			argsInput, err := e.applyJSONataArguments(ctx, state.Arguments, statesVar, execCtx.VariableScope)
			if err != nil {
				return "", e.newQueryEvalError(ctx, execCtx, "Arguments", err.Error())
			}
			execCtx.AfterArguments = &argsInput
			return argsInput, nil
		}
		if state.Parameters != nil {
			applied, evalErr := e.applyParameters(taskToken, baseInput, state.Parameters)
			if evalErr != nil {
				return "", evalErr
			}
			return applied, nil
		}
		return baseInput, nil
	}

	// The attempt-1 evaluation feeds the history events and the output
	// processing below; retries re-evaluate inside the loop.
	taskToken := ""
	if e.isActivityResource(state.Resource) {
		taskToken = generateTaskToken()
	}
	execCtx.TaskToken = taskToken
	processedInput, evalErr := evaluateTaskInput(taskToken)
	if evalErr != nil {
		return "", "", evalErr
	}

	timeoutSeconds := state.GetTimeoutSeconds()
	heartbeatSeconds := state.GetHeartbeatSeconds()

	if isJSONata {
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
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
		Type:            "TaskStateEntered",
		Timestamp:       time.Now().UTC(),
		TaskStateEnteredEventDetails: &sfnstore.TaskStateEnteredEventDetails{
			Input: processedInput,
			Name:  execCtx.CurrentState,
		},
	})

	eventId = execCtx.nextEventId()
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

	var retryPolicies []*sfnstore.RetryPolicy
	if len(state.Retry) > 0 {
		retryPolicies = state.Retry
	}

	attempt := int32(0)
	attemptInput := processedInput
	attemptToken := taskToken

	for {
		attempt++
		execCtx.RetryCount = attempt - 1

		if attempt > 1 {
			// Every retry schedules a fresh task: the activity record is
			// keyed by token, so a new token keeps attempts as separate
			// records and invalidates the previous attempt's token; the
			// input is re-evaluated so the worker receives the new token
			// wherever the token appears in the payload.
			attemptToken = ""
			if e.isActivityResource(state.Resource) {
				attemptToken = generateTaskToken()
			}
			execCtx.TaskToken = attemptToken
			var evalErr *ExecutionError
			attemptInput, evalErr = evaluateTaskInput(attemptToken)
			if evalErr != nil {
				return "", "", evalErr
			}
		}

		taskCtx := ctx
		var cancel context.CancelFunc
		if timeoutSeconds > 0 {
			taskCtx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		}

		if arnutil.IsLambdaARN(state.Resource) {
			output, taskErr = e.executeLambdaTask(taskCtx, execCtx, state, attemptInput)
		} else if e.isActivityResource(state.Resource) {
			output, _, taskErr = e.executeActivityTask(taskCtx, execCtx, state, attemptInput, attemptToken, timeoutSeconds, heartbeatSeconds)
		} else if strings.HasPrefix(state.Resource, "arn:aws:states:::sqs:") {
			output, taskErr = e.executeSQSTask(taskCtx, execCtx, state, attemptInput)
		} else if strings.HasPrefix(state.Resource, "arn:aws:states:::sns:") {
			output, taskErr = e.executeSNSTask(taskCtx, execCtx, state, attemptInput)
		} else if strings.HasPrefix(state.Resource, "arn:aws:states:::events:") {
			output, taskErr = e.executeEventsTask(taskCtx, execCtx, state, attemptInput)
		} else if strings.HasPrefix(state.Resource, "arn:aws:states:::dynamodb:") {
			output, taskErr = e.executeDynamoDBTask(taskCtx, execCtx, state, attemptInput)
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

		if len(retryPolicies) > 0 {
			matchedRetry := e.findMatchingRetryPolicy(retryPolicies, errorCode)
			if matchedRetry != nil && execCtx.RetryCount < matchedRetry.MaxAttempts {
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
		if err := json.Unmarshal([]byte(output), &resultData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidOutput", Cause: "failed to parse output JSON"}
		}
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
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
		selected, selErr := e.applyResultSelector(output, state.GetResultSelector(), attemptToken)
		if selErr != nil {
			return "", "", selErr
		}
		output = selected
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

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "TaskStateExited",
		Timestamp:       time.Now().UTC(),
		TaskStateExitedEventDetails: &sfnstore.TaskStateExitedEventDetails{
			Output: output,
			Name:   execCtx.CurrentState,
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
	if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
		return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
	}
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
