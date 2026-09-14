package sfn

// This file implements the .waitForTaskToken callback pattern and the Step
// Functions self-integration (states:startExecution, plain and AWS SDK
// forms). The callback substrate is the task-token machinery activities
// use: a record keyed by the attempt's token, resolved by
// SendTaskSuccess/SendTaskFailure/SendTaskHeartbeat, and WaitForTaskResult
// for the suspension. Registration deliberately skips the activity poll
// queues — the token reaches its consumer through the integration payload.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// taskTimeoutDefaultSeconds bounds a token-waiting task (an activity
// worker report or a .waitForTaskToken callback) when the definition sets
// no TimeoutSeconds. The Task-state reference fixes the default at
// 99,999,999 seconds; executions themselves are capped by the one-year
// service quota, which is the effective ceiling for a waiting task.
const taskTimeoutDefaultSeconds = 99999999

// taskWaitTimeout derives a token-waiting attempt's effective wait: the
// definition's TimeoutSeconds, or the documented default when unset.
func taskWaitTimeout(timeoutSeconds int32) time.Duration {
	if timeoutSeconds <= 0 {
		return taskTimeoutDefaultSeconds * time.Second
	}
	return time.Duration(timeoutSeconds) * time.Second
}

// childSyncPollInterval is how often a .sync self-integration checks the
// child execution's status. Step Functions itself learns of completion
// through events and DescribeExecution polling; a short poll keeps the
// observable latency close to that without an event-rule dependency.
const childSyncPollInterval = 200 * time.Millisecond

// executeCallbackTask runs one attempt of a .waitForTaskToken integration:
// submit the payload (which already carries the attempt token) to the
// integrated service, record the modelled TaskSubmitted event, then
// suspend until the token returns with SendTaskSuccess or SendTaskFailure
// or the timeout fires. A timeout or a retried attempt spends the token —
// the next attempt mints a fresh one and re-evaluates the payload so the
// consumer receives the new token wherever the definition references it.
func (e *Executor) executeCallbackTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, resource, input, taskToken string, timeoutSeconds, heartbeatSeconds int32, resolved *resolvedInvocation) (string, error) {
	// Registration precedes the submit: the record must exist before any
	// consumer can present the token back, and the ActivityArn it carries
	// is the definition's resource — never a pollable activity ARN, so no
	// GetActivityTask can ever claim it.
	task := &sfnstore.ActivityTask{
		ActivityArn:  state.Resource,
		ExecutionArn: execCtx.Execution.ExecutionArn,
		Input:        input,
		TaskToken:    taskToken,
	}
	if err := e.store.CreateCallbackTask(task); err != nil {
		return "", fmt.Errorf("failed to register callback task: %w", err)
	}
	// Panic safety for the registration: the record must not outlive the
	// attempt as RUNNING when the executor unwinds abnormally — the
	// execution runner recovers panics above this frame, and nothing else
	// between registration and the wait cleans up. Normal exits leave the
	// record in place: the wait itself settles it (a report leaves
	// SUCCEEDED/FAILED, a timeout flips it to TIMED_OUT), and a duplicate
	// report must keep meeting the closed-task answer (TaskTimedOut —
	// "the task token has either expired or the task associated with the
	// token has already been closed"), not TaskDoesNotExist.
	attemptSettled := false
	defer func() {
		if !attemptSettled {
			e.store.DeleteCallbackTask(taskToken)
		}
	}()

	submitOutput, submitErr := e.submitCallbackPayload(ctx, execCtx, resource, input, resolved)
	if submitErr != nil {
		return "", submitErr
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "TaskSubmitted",
		Timestamp:       time.Now().UTC(),
		TaskSubmittedEventDetails: &sfnstore.TaskSubmittedEventDetails{
			Resource:     state.Resource,
			ResourceType: taskResourceType(state.Resource),
			Output:       submitOutput,
		},
	})

	timeout := taskWaitTimeout(timeoutSeconds)
	hbTimeout := time.Duration(heartbeatSeconds) * time.Second

	// Every WaitForTaskResult exit settles the record (terminal status via
	// the report paths, TIMED_OUT via the store's own marking, or the
	// record is already gone) — from here the deferred cleanup must not
	// fire.
	result, err := e.store.WaitForTaskResult(ctx, taskToken, timeout, hbTimeout)
	attemptSettled = true
	if err != nil {
		if err == sfnstore.ErrTaskTimeout || err == sfnstore.ErrHeartbeatTimeout {
			return "", &taskFailure{code: "States.Timeout", cause: "Task timed out", timedOut: true}
		}
		return "", err
	}

	if result.Error != nil {
		code := result.Error.Error()
		if code == "" {
			code = "States.TaskFailed"
		}
		return "", &taskFailure{code: code, cause: result.Cause}
	}
	return result.Output, nil
}

// submitCallbackPayload performs the pattern's service call. The
// optimised integrations reject the callback pattern at validation time
// where the integrations table marks it unsupported (DynamoDB, the
// Run-a-Job families); the AWS SDK namespace carries the callback pattern
// for every service, of which the executor implements the families with a
// callback consumer here — the remaining aws-sdk families fail at run
// time with a cause naming the resource.
func (e *Executor) submitCallbackPayload(ctx context.Context, execCtx *ExecutionContext, resource, input string, resolved *resolvedInvocation) (string, error) {
	switch {
	case resource == "arn:aws:states:::lambda:invoke" || resource == "arn:aws:states:::aws-sdk:lambda:invoke":
		out, err := e.invokeLambdaIntegration(ctx, input)
		if err != nil {
			// A function error aborts the submit itself: the token-carrying
			// payload never reached a consumer that could answer it, so the
			// attempt is a submit failure whatever the function's own error
			// name says.
			var tf *taskFailure
			if errors.As(err, &tf) && !tf.failedToStart {
				tf.failedToStart = true
			}
			return "", err
		}
		return out, nil
	case strings.HasPrefix(resource, "arn:aws:states:::sqs:") || strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:sqs:"):
		return e.executeSQSTask(ctx, resource, input, resolved)
	case strings.HasPrefix(resource, "arn:aws:states:::sns:") || strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:sns:"):
		return e.executeSNSTask(ctx, resource, input)
	case strings.HasPrefix(resource, "arn:aws:states:::events:") || strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:eventbridge:"):
		return e.executeEventsTask(ctx, execCtx, resource, input)
	case e.isStartExecutionResource(resource):
		child, err := e.launchIntegrationChildExecution(ctx, input)
		if err != nil {
			return "", err
		}
		return startExecutionResponseJSON(child), nil
	default:
		return "", fmt.Errorf("the .waitForTaskToken pattern is not supported for resource: %s", resource)
	}
}

// invokeLambdaIntegration runs the optimised and AWS SDK lambda:invoke
// forms: FunctionName selects the function, Payload is the invocation
// payload, and the task result carries the documented invoke response
// wrapper (Payload, StatusCode, ExecutedVersion).
func (e *Executor) invokeLambdaIntegration(ctx context.Context, input string) (string, error) {
	if e.bus == nil || e.bus.LambdaInvoker() == nil {
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "lambda invoker not configured", failedToStart: true}
	}

	var params map[string]interface{}
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "lambda invoke parameters must be a JSON object", failedToStart: true}
	}

	functionName := getStr(params, "FunctionName")
	if functionName == "" {
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "FunctionName is required for the lambda:invoke integration", failedToStart: true}
	}

	var payloadBytes []byte
	if payload, ok := params["Payload"]; ok {
		b, err := json.Marshal(payload)
		if err != nil {
			return "", &taskFailure{code: "Lambda.ServiceException", cause: "Payload is not serialisable: " + err.Error(), failedToStart: true}
		}
		payloadBytes = b
	} else {
		payloadBytes = []byte("{}")
	}

	statusCode, respPayload, err := e.bus.LambdaInvoker().InvokeForGateway(ctx, functionName, payloadBytes)
	if err != nil {
		// The invocation transport itself failed; the function never ran.
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "failed to invoke Lambda function: " + err.Error(), failedToStart: true}
	}
	if statusCode != 200 {
		code, cause := lambdaFunctionError(respPayload)
		return "", &taskFailure{code: code, cause: cause}
	}

	var payloadValue interface{}
	if err := json.Unmarshal(respPayload, &payloadValue); err != nil {
		payloadValue = string(respPayload)
	}
	result := map[string]interface{}{
		"Payload":         payloadValue,
		"StatusCode":      statusCode,
		"ExecutedVersion": "$LATEST",
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "failed to marshal invoke response: " + err.Error(), failedToStart: true}
	}
	return string(resultJSON), nil
}

// executeStartExecutionTask runs the plain states:startExecution
// integration: start the child execution and return the StartExecution
// response pair (executionArn, startDate) without waiting.
func (e *Executor) executeStartExecutionTask(ctx context.Context, execCtx *ExecutionContext, input string) (string, error) {
	child, err := e.launchIntegrationChildExecution(ctx, input)
	if err != nil {
		return "", err
	}
	return startExecutionResponseJSON(child), nil
}

// executeStartExecutionSyncTask runs the .sync self-integration: start the
// child execution and wait for its terminal status. The result carries the
// DescribeExecution members with the child's Output — as a string for the
// legacy .sync form and as parsed JSON for .sync:2. A child that did not
// succeed fails the task with the child's own error identity.
func (e *Executor) executeStartExecutionSyncTask(ctx context.Context, execCtx *ExecutionContext, input string, jsonOutput bool) (string, error) {
	child, err := e.launchIntegrationChildExecution(ctx, input)
	if err != nil {
		return "", err
	}

	ticker := time.NewTicker(childSyncPollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			exec, gerr := e.store.GetExecution(ctx, child.ExecutionArn)
			if gerr != nil || exec == nil {
				continue
			}
			switch exec.Status {
			case "SUCCEEDED", "FAILED", "TIMED_OUT", "ABORTED":
				return e.childSyncResult(exec, jsonOutput)
			}
		case <-ctx.Done():
			// The parent no longer waits: stopping the nested execution is
			// the documented best-effort cancellation of a .sync job.
			e.store.CancelExecution(child.ExecutionArn)
			return "", ctx.Err()
		}
	}
}

// childSyncResult renders the terminal child record as the .sync task
// result.
func (e *Executor) childSyncResult(exec *sfnstore.Execution, jsonOutput bool) (string, error) {
	if exec.Status != "SUCCEEDED" {
		code := exec.Error
		if code == "" {
			code = "States.TaskFailed"
		}
		// "If a nested state machine throws a States.Timeout, the parent
		// will receive a States.TaskFailed error" — the timeout identity
		// alone is remapped; every other child error keeps its own code,
		// and the cause always passes through untouched.
		if code == "States.Timeout" {
			code = "States.TaskFailed"
		}
		return "", &taskFailure{code: code, cause: exec.Cause}
	}

	outputValue := interface{}(exec.Output)
	if jsonOutput && exec.Output != "" {
		var parsed interface{}
		if err := json.Unmarshal([]byte(exec.Output), &parsed); err == nil {
			outputValue = parsed
		}
	}

	result := map[string]interface{}{
		"ExecutionArn":    exec.ExecutionArn,
		"StateMachineArn": exec.StateMachineArn,
		"Name":            exec.Name,
		"Status":          exec.Status,
		"StartDate":       awsEpochSeconds(exec.StartDate),
		"StopDate":        awsEpochSeconds(exec.StopDate),
		"Input":           exec.Input,
		"Output":          outputValue,
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", &taskFailure{code: "States.TaskFailed", cause: "failed to marshal child execution result: " + err.Error(), failedToStart: true}
	}
	return string(resultJSON), nil
}

// launchIntegrationChildExecution starts one child execution for the Step
// Functions self-integration. The parameters follow the StartExecution
// API: StateMachineArn (required), Name and Input optional. The child runs
// on its own executor goroutine, registered for StopExecution, with the
// same panic isolation the top-level launch path applies.
func (e *Executor) launchIntegrationChildExecution(ctx context.Context, input string) (*sfnstore.Execution, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return nil, &taskFailure{code: "States.TaskFailed", cause: "startExecution parameters must be a JSON object", failedToStart: true}
	}

	smArn := getStr(params, "StateMachineArn")
	if smArn == "" {
		return nil, &taskFailure{code: "States.TaskFailed", cause: "StateMachineArn is required for the startExecution integration", failedToStart: true}
	}
	parsed, err := arnutil.ParseARN(smArn)
	if err != nil {
		return nil, &taskFailure{code: "States.TaskFailed", cause: "invalid StateMachineArn: " + smArn, failedToStart: true}
	}
	if _, err := e.store.GetStateMachine(ctx, smArn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, &taskFailure{code: "StateMachineDoesNotExist", cause: "state machine not found: " + smArn, failedToStart: true}
		}
		return nil, &taskFailure{code: "States.TaskFailed", cause: "failed to read the state machine: " + err.Error(), failedToStart: true}
	}

	childName := getStr(params, "Name")
	if childName == "" {
		childName = generateExecutionName()
	}

	childInput := ""
	if raw, ok := params["Input"]; ok {
		if s, isStr := raw.(string); isStr {
			childInput = s
		} else {
			b, merr := json.Marshal(raw)
			if merr != nil {
				return nil, &taskFailure{code: "States.TaskFailed", cause: "Input is not serialisable: " + merr.Error(), failedToStart: true}
			}
			childInput = string(b)
		}
	}

	region, accountID := parsed.Region, parsed.AccountID
	if region == "" {
		region = e.region
	}
	if accountID == "" {
		accountID = e.accountID
	}
	smName := arnutil.ExtractStateMachineNameFromARN(smArn)
	childArn := arnutil.NewARNBuilder(accountID, region).StepFunctions().Execution(smName, childName)

	now := time.Now().UTC()
	child := &sfnstore.Execution{
		ExecutionArn:    childArn,
		StateMachineArn: smArn,
		Name:            childName,
		Status:          "RUNNING",
		Input:           childInput,
		StartDate:       now,
	}
	if err := e.store.CreateExecution(ctx, child); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil, &taskFailure{code: "ExecutionAlreadyExists", cause: "an execution named " + childName + " already exists", failedToStart: true}
		}
		return nil, &taskFailure{code: "States.TaskFailed", cause: "failed to create child execution: " + err.Error(), failedToStart: true}
	}

	childExecCtx, cancel := context.WithCancel(context.Background())
	handle := e.store.RegisterExecution(childArn, cancel)
	childExecutor := NewExecutorWithStores(e.store, e.bus, e.accountID, e.region, e.taskAuthz)
	go func() {
		defer e.store.UnregisterExecution(childArn, handle)
		defer func() {
			if r := recover(); r != nil {
				// A panic must never leave a RUNNING child behind.
				logs.Error("sfn: panic in integration child execution", logs.String("arn", childArn), logs.Any("panic", r))
				child.Status = "FAILED"
				child.Error = "States.Runtime"
				child.Cause = fmt.Sprintf("internal panic: %v", r)
				child.StopDate = time.Now().UTC()
				appendTerminalFailureEvent(context.Background(), e.store, child, child.Error, child.Cause)
				_ = e.store.UpdateExecution(context.Background(), child)
			}
		}()
		if err := childExecutor.ExecuteStateMachine(childExecCtx, child); err != nil {
			logs.Error("sfn: integration child execution failed", logs.String("arn", childArn), logs.Err(err))
		}
	}()

	return child, nil
}

// startExecutionResponseJSON renders the StartExecution API response pair
// for the plain and callback self-integration forms.
func startExecutionResponseJSON(child *sfnstore.Execution) string {
	result := map[string]interface{}{
		"ExecutionArn": child.ExecutionArn,
		"StartDate":    awsEpochSeconds(child.StartDate),
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "{}"
	}
	return string(resultJSON)
}
