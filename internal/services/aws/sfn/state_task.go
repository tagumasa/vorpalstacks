package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	baseInput, ipErr := e.applyInputPath(execCtx, execCtx.Input, state.GetInputPath())
	if ipErr != nil {
		return "", "", ipErr
	}

	// evaluateTaskInput applies Arguments/Parameters for a single attempt.
	// It runs per attempt rather than once per state because the attempt's
	// token and $$.State.RetryCount must reflect the attempt being
	// dispatched, and each activity attempt is stored under its own token.
	// The token reaches the payload through both dialects: JSONPath
	// Parameters resolve $$.Task.Token from the taskToken parameter, and
	// JSONata Arguments resolve $states.context.Task.Token from the context
	// object built with that token. taskToken is empty for resources
	// without a token (activities and .waitForTaskToken callbacks mint
	// one), where any Task.Token reference fails the evaluation.
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
			applied, evalErr := e.applyParameters(execCtx, taskToken, baseInput, state.Parameters)
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
	if e.taskUsesToken(state.Resource) {
		taskToken = generateTaskToken()
	}
	execCtx.TaskToken = taskToken
	processedInput, evalErr := evaluateTaskInput(taskToken)
	if evalErr != nil {
		return "", "", evalErr
	}
	taskCredentials, credErr := e.evaluateTaskCredentials(ctx, execCtx, state, baseInput)
	if credErr != nil {
		return "", "", credErr
	}

	timeoutSeconds := state.GetTimeoutSeconds()
	heartbeatSeconds := state.GetHeartbeatSeconds()

	if !isJSONata {
		// TimeoutSecondsPath and HeartbeatSecondsPath are reference
		// paths selecting positive integers from the state input; each
		// excludes its literal form.
		if state.TimeoutSecondsPath != "" {
			resolved, err := resolveTaskSecondsPath(processedInput, state.TimeoutSecondsPath, "TimeoutSecondsPath")
			if err != nil {
				return "", "", err
			}
			timeoutSeconds = resolved
		}
		if state.HeartbeatSecondsPath != "" {
			resolved, err := resolveTaskSecondsPath(processedInput, state.HeartbeatSecondsPath, "HeartbeatSecondsPath")
			if err != nil {
				return "", "", err
			}
			heartbeatSeconds = resolved
		}
	}

	if isJSONata {
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		vars := buildVarsMap(statesVar, execCtx.VariableScope)

		// TimeoutSeconds and HeartbeatSeconds carry a {% %} expression in
		// the JSONata dialect; both resolve to a positive whole number of
		// seconds through the one evaluation path.
		evalSecondsField := func(field string, raw interface{}, slot *int32) *ExecutionError {
			s, ok := raw.(string)
			if !ok || !IsExpression(s) {
				return nil
			}
			result, err := EvaluateJSONata(ctx, UnwrapExpression(s), nil, vars)
			if err != nil {
				return e.newQueryEvalError(ctx, execCtx, field, err.Error())
			}
			seconds, evalErr := queryEvalSecondsValue(field, result)
			if evalErr != nil {
				return e.newQueryEvalError(ctx, execCtx, field, evalErr.Error())
			}
			*slot = seconds
			return nil
		}
		if execErr := evalSecondsField("TimeoutSeconds", state.TimeoutSeconds, &timeoutSeconds); execErr != nil {
			return "", "", execErr
		}
		if execErr := evalSecondsField("HeartbeatSeconds", state.HeartbeatSeconds, &heartbeatSeconds); execErr != nil {
			return "", "", execErr
		}
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:             execCtx.Execution.ExecutionArn,
		EventId:                  eventId,
		PreviousEventId:          eventId - 1,
		Type:                     "TaskStateEntered",
		Timestamp:                time.Now().UTC(),
		StateEnteredEventDetails: stateEnteredDetails(execCtx, processedInput),
	})

	var output string
	var taskErr error
	var errorCode, cause string

	resourceType := taskResourceType(state.Resource)

	// The integration pattern arrives as the resource suffix. It stays on
	// state.Resource everywhere the definition's identifier is reported
	// (events, errors), while taskResource — the suffix-stripped identifier
	// — is what the integration executors parse.
	integrationPattern := ""
	taskResource := state.Resource
	if strings.HasSuffix(taskResource, ".waitForTaskToken") {
		integrationPattern = "waitForTaskToken"
		taskResource = strings.TrimSuffix(taskResource, ".waitForTaskToken")
	} else if strings.HasSuffix(taskResource, ".sync:2") {
		integrationPattern = "sync2"
		taskResource = strings.TrimSuffix(taskResource, ".sync:2")
	} else if strings.HasSuffix(taskResource, ".sync") {
		integrationPattern = "sync"
		taskResource = strings.TrimSuffix(taskResource, ".sync")
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
			if e.taskUsesToken(state.Resource) {
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

		// Each attempt records the modelled scheduled/started pair for its
		// resource class. Activity tasks record ActivityScheduled when the
		// task record is created inside executeActivityTask and
		// ActivityStarted when a worker claims the task, so they schedule
		// nothing here.
		if arnutil.IsLambdaARN(state.Resource) {
			eventId = execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Type:            "LambdaFunctionScheduled",
				Timestamp:       time.Now().UTC(),
				LambdaFunctionScheduledEventDetails: &sfnstore.LambdaFunctionScheduledEventDetails{
					Resource:         state.Resource,
					Input:            attemptInput,
					TimeoutInSeconds: timeoutSeconds,
					TaskCredentials:  taskCredentials,
				},
			})
			eventId = execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Type:            "LambdaFunctionStarted",
				Timestamp:       time.Now().UTC(),
			})
		} else if !isActivityResource(state.Resource) {
			eventId = execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Type:            "TaskScheduled",
				Timestamp:       time.Now().UTC(),
				TaskScheduledEventDetails: &sfnstore.TaskScheduledEventDetails{
					Resource:           state.Resource,
					ResourceType:       resourceType,
					Region:             e.region,
					Parameters:         attemptInput,
					TimeoutInSeconds:   timeoutSeconds,
					HeartbeatInSeconds: heartbeatSeconds,
					TaskCredentials:    taskCredentials,
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
					Resource:     state.Resource,
					ResourceType: resourceType,
				},
			})
		}

		// Credentials "specifies a target role the state machine's
		// execution role must assume before invoking the specified
		// Resource" — the assumption and the action authorisation run
		// before every dispatch, and a denial fails the attempt as
		// States.Permissions before any integration call is made. The
		// authorisation's resource resolution rides the attempt's carrier
		// so the dispatch reuses it instead of repeating its lookups.
		resolved := &resolvedInvocation{}
		if taskCredentials != nil {
			taskErr = e.authoriseTaskCredentialsInvocation(taskCtx, taskResource, attemptInput, taskCredentials.RoleArn, resolved)
		} else {
			taskErr = nil
		}

		if taskErr == nil {
			switch {
			case integrationPattern == "waitForTaskToken":
				output, taskErr = e.executeCallbackTask(taskCtx, execCtx, state, taskResource, attemptInput, attemptToken, timeoutSeconds, heartbeatSeconds, resolved)
			case arnutil.IsLambdaARN(taskResource):
				output, taskErr = e.executeLambdaTask(taskCtx, execCtx, state, attemptInput)
			case isActivityResource(taskResource):
				output, taskErr = e.executeActivityTask(taskCtx, execCtx, state, attemptInput, attemptToken, timeoutSeconds, heartbeatSeconds)
			case taskResource == "arn:aws:states:::lambda:invoke" || taskResource == "arn:aws:states:::aws-sdk:lambda:invoke":
				output, taskErr = e.invokeLambdaIntegration(taskCtx, attemptInput)
			case strings.HasPrefix(taskResource, "arn:aws:states:::sqs:") || strings.HasPrefix(taskResource, "arn:aws:states:::aws-sdk:sqs:"):
				output, taskErr = e.executeSQSTask(taskCtx, taskResource, attemptInput, resolved)
			case strings.HasPrefix(taskResource, "arn:aws:states:::sns:") || strings.HasPrefix(taskResource, "arn:aws:states:::aws-sdk:sns:"):
				output, taskErr = e.executeSNSTask(taskCtx, taskResource, attemptInput)
			case strings.HasPrefix(taskResource, "arn:aws:states:::events:") || strings.HasPrefix(taskResource, "arn:aws:states:::aws-sdk:eventbridge:"):
				output, taskErr = e.executeEventsTask(taskCtx, execCtx, taskResource, attemptInput)
			case strings.HasPrefix(taskResource, "arn:aws:states:::dynamodb:") || strings.HasPrefix(taskResource, "arn:aws:states:::aws-sdk:dynamodb:"):
				output, taskErr = e.executeDynamoDBTask(taskCtx, taskResource, attemptInput)
			case (integrationPattern == "sync" || integrationPattern == "sync2") && taskResource == "arn:aws:states:::states:startExecution":
				// Run a Job is an optimised-integration pattern: the AWS SDK
				// namespace carries only Request Response and the callback
				// pattern (integrations table), so the aws-sdk form of
				// startExecution refuses the .sync suffixes like every other
				// aws-sdk resource.
				output, taskErr = e.executeStartExecutionSyncTask(taskCtx, execCtx, attemptInput, integrationPattern == "sync2")
			case e.isStartExecutionResource(taskResource):
				output, taskErr = e.executeStartExecutionTask(taskCtx, execCtx, attemptInput)
			default:
				taskErr = fmt.Errorf("unsupported resource type: %s", state.Resource)
			}
		}
		if taskErr != nil {
			// An execution-level interruption (StopExecution, the shutdown
			// sweep, or the machine-level TimeoutSeconds deadline) is not a
			// task failure: no failure event, retry or Catch runs for it —
			// the terminal path owns the history and records the
			// execution's own terminal pair (ExecutionAborted or
			// ExecutionTimedOut). The vehicle wraps context.Canceled for
			// cancellations so Map and Parallel aggregates recognise the
			// interruption at every layer; a deadline keeps the
			// States.Timeout identity.
			if ctx.Err() != nil {
				if cancel != nil {
					cancel()
				}
				return "", "", executionInterruptedError(ctx, "Execution interrupted")
			}
			// Integrations other than Lambda classify their invocation
			// failures here: the attempt never produced a result, so it
			// failed to start under the documented wildcard name for task
			// failures.
			var tf *taskFailure
			if !errors.As(taskErr, &tf) {
				taskErr = &taskFailure{code: "States.TaskFailed", cause: taskErr.Error(), failedToStart: true}
			}
		}

		if taskErr == nil {
			if cancel != nil {
				cancel()
			}
			break
		}

		// A failed attempt keeps its error identity: a Lambda function
		// error carries the function's own error name, a transport failure
		// the service class, and unclassified failures the wildcard name.
		// A timeout is recognised only from the typed signals (the task's
		// OWN deadline wrapper or the waiter layer's timedOut flag) — never
		// from the error text, which user-controlled causes may quote. The
		// cancel guard is the discriminator: a task without TimeoutSeconds
		// has no wrapper, so its taskCtx IS the machine context, and a
		// machine-level deadline expiring there must not masquerade as the
		// task's own timeout event.
		timedOut := cancel != nil && taskCtx.Err() == context.DeadlineExceeded
		var startFailed bool
		var tf *taskFailure
		if errors.As(taskErr, &tf) {
			timedOut = timedOut || tf.timedOut
			errorCode, cause, startFailed = tf.code, tf.cause, tf.failedToStart
		} else {
			errorCode, cause = "States.TaskFailed", taskErr.Error()
		}
		if timedOut {
			errorCode = "States.Timeout"
			cause = "Task timed out"
		}

		// Every failed attempt records its resource-classed failure event;
		// activity tasks record their own Activity* events when the worker
		// report or the timeout arrives, so they write nothing here.
		if !isActivityResource(state.Resource) {
			e.logTaskAttemptFailureEvent(ctx, execCtx, state, resourceType, errorCode, cause, timedOut, startFailed)
		}

		if len(state.Retry) > 0 {
			matchedRetry := e.findMatchingRetryPolicy(state.Retry, errorCode)
			if matchedRetry != nil && execCtx.RetryCount < matchedRetry.MaxAttempts {
				if cancel != nil {
					cancel()
				}
				if e.sleepForRetry(ctx, matchedRetry, attempt) {
					return "", "", executionInterruptedError(ctx, "Execution interrupted during retry")
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
				if isJSONata {
					return e.executeJSONataCatch(ctx, execCtx, processedInput, errorCode, cause, catchPolicy)
				}

				catchOutput, coErr := e.buildCatchOutput(processedInput, errorCode, cause, catchPolicy.ResultPath)
				if coErr != nil {
					return "", "", coErr
				}
				if len(catchPolicy.Assign) > 0 {
					if err := e.applyJSONPathCatchAssign(execCtx, catchPolicy.Assign, catchOutput); err != nil {
						return "", "", err
					}
				}
				return catchOutput, catchPolicy.Next, nil
			}
		}
		return "", "", &ExecutionError{ErrorCode: errorCode, Cause: cause}
	}

	// JSONPath Task Assign evaluates against the raw result ("In Task, Map,
	// Parallel states, $ refers to the API/sub-workflow result"); new values
	// become visible in the next state.
	if len(state.Assign) > 0 && !isJSONata {
		var assignRoot interface{}
		if err := json.Unmarshal([]byte(output), &assignRoot); err != nil {
			assignRoot = nil
		}
		evaluated, err := e.evaluateJSONPathAssign(execCtx, state.Assign, assignRoot)
		if err != nil {
			return "", "", newJSONPathEvalError("Assign", err)
		}
		execCtx.PendingAssign = evaluated
	}

	if isJSONata {
		jsonataOutput, jerr := jsonataOutputOf(state)
		if jerr != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Output", jerr.Error())
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

		if jsonataOutput != nil {
			resolved, err := e.applyJSONataOutput(ctx, jsonataOutput, statesVar, execCtx.VariableScope)
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
		selected, selErr := e.applyResultSelector(execCtx, output, state.GetResultSelector(), attemptToken)
		if selErr != nil {
			return "", "", selErr
		}
		output = selected
		folded, rpErr := e.applyResultPath(processedInput, output, state.ResultPath)
		if rpErr != nil {
			return "", "", rpErr
		}
		output = folded
		selectedOut, opErr := e.applyOutputPath(execCtx, output, state.GetOutputPath())
		if opErr != nil {
			return "", "", opErr
		}
		output = selectedOut
	}

	switch {
	case arnutil.IsLambdaARN(state.Resource):
		eventId = execCtx.nextEventId()
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    execCtx.Execution.ExecutionArn,
			EventId:         eventId,
			PreviousEventId: eventId - 1,
			Type:            "LambdaFunctionSucceeded",
			Timestamp:       time.Now().UTC(),
			LambdaFunctionSucceededEventDetails: &sfnstore.LambdaFunctionSucceededEventDetails{
				Output: output,
			},
		})
	case !isActivityResource(state.Resource):
		eventId = execCtx.nextEventId()
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    execCtx.Execution.ExecutionArn,
			EventId:         eventId,
			PreviousEventId: eventId - 1,
			Type:            "TaskSucceeded",
			Timestamp:       time.Now().UTC(),
			TaskSucceededEventDetails: &sfnstore.TaskSucceededEventDetails{
				Resource:     state.Resource,
				ResourceType: resourceType,
				Output:       output,
			},
		})
	}
	// Activity tasks record ActivitySucceeded inside executeActivityTask
	// when the worker's report arrives.

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:            execCtx.Execution.ExecutionArn,
		EventId:                 eventId,
		PreviousEventId:         eventId - 1,
		Type:                    "TaskStateExited",
		Timestamp:               time.Now().UTC(),
		StateExitedEventDetails: stateExitedDetails(execCtx, output),
	})

	return output, state.Next, nil
}

// taskFailure carries one task attempt's failure identity: the error name
// AWS reports for it (a Lambda function's own error type, or a service-class
// name) and its cause, plus whether the attempt failed before the task ever
// produced a result, and whether it ended in a timeout (the waiter layer's
// task/heartbeat timeout — the classification consumes the flag instead of
// pattern-matching the error text, which user-controlled causes can quote).
type taskFailure struct {
	code          string
	cause         string
	failedToStart bool
	timedOut      bool
}

func (f *taskFailure) Error() string { return f.code + ": " + f.cause }

// taskResourceType derives the resourceType member of the task event
// details from the resource identifier: the integration family name.
func taskResourceType(resource string) string {
	switch {
	case arnutil.IsLambdaARN(resource):
		return "lambda"
	case strings.HasPrefix(resource, "arn:aws:states:::lambda:invoke"):
		return "lambda"
	case strings.HasPrefix(resource, "arn:aws:states:::sqs:"):
		return "sqs"
	case strings.HasPrefix(resource, "arn:aws:states:::sns:"):
		return "sns"
	case strings.HasPrefix(resource, "arn:aws:states:::events:"):
		return "events"
	case strings.HasPrefix(resource, "arn:aws:states:::dynamodb:"):
		return "dynamodb"
	case strings.HasPrefix(resource, "arn:aws:states:::states:"):
		return "states"
	case strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:"):
		// The AWS SDK namespace reports the SDK service identifier as the
		// family name.
		rest := strings.TrimPrefix(resource, "arn:aws:states:::aws-sdk:")
		if i := strings.Index(rest, ":"); i > 0 {
			return rest[:i]
		}
		return "task"
	case isActivityResource(resource):
		return "activity"
	default:
		return "task"
	}
}

// evaluateTaskCredentials resolves the Credentials payload template to
// the role ARN the invocation runs under ("you can also specify a JSONPath
// value or an intrinsic function that resolves to an IAM role ARN at
// runtime based on the execution input"). The resolved ARN rides the
// scheduled event's taskCredentials member and drives the per-task
// assume-role authorisation the dispatch path runs before invoking the
// Resource.
func (e *Executor) evaluateTaskCredentials(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, baseInput string) (*sfnstore.TaskScheduledCredentials, *ExecutionError) {
	if state.Credentials == nil {
		return nil, nil
	}
	var resolved interface{}
	if IsJSONataState(state, execCtx.QueryLanguage) {
		var inputData interface{}
		if err := json.Unmarshal([]byte(baseInput), &inputData); err != nil {
			return nil, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		value, err := ResolveTemplate(ctx, state.Credentials, nil, buildVarsMap(statesVar, execCtx.VariableScope))
		if err != nil {
			return nil, e.newQueryEvalError(ctx, execCtx, "Credentials", err.Error())
		}
		resolved = value
	} else {
		credsMap, ok := state.Credentials.(map[string]interface{})
		if !ok {
			return nil, &ExecutionError{ErrorCode: "States.Runtime", Cause: "Credentials must be a payload template object"}
		}
		applied, perr := e.applyParameters(execCtx, "", baseInput, &sfnstore.Parameters{Values: credsMap})
		if perr != nil {
			return nil, perr
		}
		if jerr := json.Unmarshal([]byte(applied), &resolved); jerr != nil {
			return nil, &ExecutionError{ErrorCode: "States.Runtime", Cause: "Credentials did not resolve to a JSON object"}
		}
	}
	creds, ok := resolved.(map[string]interface{})
	if !ok {
		return nil, &ExecutionError{ErrorCode: "States.Runtime", Cause: "Credentials must resolve to an object with RoleArn"}
	}
	roleArn, _ := creds["RoleArn"].(string)
	if roleArn == "" {
		return nil, &ExecutionError{ErrorCode: "States.Runtime", Cause: "Credentials must resolve to a RoleArn"}
	}
	return &sfnstore.TaskScheduledCredentials{RoleArn: roleArn}, nil
}

// taskUsesToken reports whether a task attempt mints a task token: activity
// tasks and .waitForTaskToken callbacks both deliver the token through the
// payload ($$.Task.Token / $states.context.Task.Token) and resume only on
// SendTaskSuccess or SendTaskFailure.
func (e *Executor) taskUsesToken(resource string) bool {
	return isActivityResource(resource) || strings.HasSuffix(resource, ".waitForTaskToken")
}

// isStartExecutionResource reports whether the resource is the Step
// Functions self-integration (optimised or AWS SDK form) that starts a
// child execution.
func (e *Executor) isStartExecutionResource(resource string) bool {
	return resource == "arn:aws:states:::states:startExecution" || resource == "arn:aws:states:::aws-sdk:stepfunctions:startExecution"
}

// iterationErrorName derives the error name a Map state reports for a
// failed iteration: an execution error keeps its own identity, and
// anything else is an engine-level failure, which the documented runtime
// error name covers.
func iterationErrorName(err error) string {
	var stateErr *ExecutionError
	if errors.As(err, &stateErr) && stateErr.ErrorCode != "" {
		return stateErr.ErrorCode
	}
	return "States.Runtime"
}

// logTaskAttemptFailureEvent records one failed task attempt with the
// resource-classed event type: the Lambda family, the timeout family and the
// start-failure family each have their own HistoryEventType members.
func (e *Executor) logTaskAttemptFailureEvent(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, resourceType, errorCode, cause string, timedOut, startFailed bool) {
	eventId := execCtx.nextEventId()
	evt := &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Timestamp:       time.Now().UTC(),
	}
	isLambda := arnutil.IsLambdaARN(state.Resource)
	switch {
	case isLambda && timedOut:
		evt.Type = "LambdaFunctionTimedOut"
		evt.LambdaFunctionTimedOutEventDetails = &sfnstore.LambdaFunctionTimedOutEventDetails{Error: errorCode, Cause: cause}
	case isLambda && startFailed:
		evt.Type = "LambdaFunctionStartFailed"
		evt.LambdaFunctionStartFailedEventDetails = &sfnstore.LambdaFunctionStartFailedEventDetails{Error: errorCode, Cause: cause}
	case isLambda:
		evt.Type = "LambdaFunctionFailed"
		evt.LambdaFunctionFailedEventDetails = &sfnstore.LambdaFunctionFailedEventDetails{Error: errorCode, Cause: cause}
	case timedOut:
		evt.Type = "TaskTimedOut"
		evt.TaskTimedOutEventDetails = &sfnstore.TaskTimedOutEventDetails{Resource: state.Resource, ResourceType: resourceType, Error: errorCode, Cause: cause}
	case strings.HasSuffix(state.Resource, ".waitForTaskToken") && startFailed:
		// A callback task whose submission to the integrated service failed:
		// the model dedicates TaskSubmitFailed to this stage, before any
		// token callback could exist.
		evt.Type = "TaskSubmitFailed"
		evt.TaskSubmitFailedEventDetails = &sfnstore.TaskSubmitFailedEventDetails{Resource: state.Resource, ResourceType: resourceType, Error: errorCode, Cause: cause}
	case startFailed:
		evt.Type = "TaskStartFailed"
		evt.TaskStartFailedEventDetails = &sfnstore.TaskStartFailedEventDetails{Resource: state.Resource, ResourceType: resourceType, Error: errorCode, Cause: cause}
	default:
		evt.Type = "TaskFailed"
		evt.TaskFailedEventDetails = &sfnstore.TaskFailedEventDetails{Resource: state.Resource, ResourceType: resourceType, Error: errorCode, Cause: cause}
	}
	e.logHistoryEvent(ctx, execCtx.Execution, evt)
}

// executeJSONataCatch builds a JSONata state's catch output: the error and
// cause pair exposed through $states, the catch's Assign, and its Output
// template, falling back to the marshalled pair when no template is set.
// Task, Parallel and Map states share it — their catch contract is
// identical.
func (e *Executor) executeJSONataCatch(ctx context.Context, execCtx *ExecutionContext, processedInput, errorCode, cause string, catchPolicy *sfnstore.CatchPolicy) (string, string, *ExecutionError) {
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

	errorJSON, err := json.Marshal(errorOutput)
	if err != nil {
		errorJSON = []byte(`{"error":"failed to marshal error output"}`)
	}
	return string(errorJSON), catchPolicy.Next, nil
}

// resolveTaskSecondsPath resolves a TimeoutSecondsPath or
// HeartbeatSecondsPath reference path against the state input; the
// selected field must hold a positive integer.
func resolveTaskSecondsPath(processedInput, path, field string) (int32, *ExecutionError) {
	var data interface{}
	if err := json.Unmarshal([]byte(processedInput), &data); err != nil {
		return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON for " + field}
	}
	inputMap, ok := data.(map[string]interface{})
	if !ok {
		return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: field + " requires object input"}
	}
	value, err := getJSONPathValue(inputMap, path)
	if err != nil {
		return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: field + " failed to resolve: " + err.Error()}
	}
	number, ok := value.(float64)
	if !ok || number != math.Trunc(number) || number < 1 || number > math.MaxInt32 {
		return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: field + " must resolve to a positive integer"}
	}
	return int32(number), nil
}
