package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// awsEpochSeconds renders a timestamp for the AWS JSON wire. Step
// Functions timestamps carry sub-second precision ("zero, three, six or
// nine digits as necessary" — context-object documentation), so a
// whole-second Unix() truncation would collapse the ordering of nearby
// events onto one wire value.
func awsEpochSeconds(t time.Time) float64 {
	return float64(t.UnixNano()) / 1e9
}

// parsePageLimit reads the maxResults parameter. Zero-as-unset
// normalisation lives in the list Cores (normaliseListLimit) so both
// protocol planes share one default; the range check also runs there
// (validateMaxResults).
func parsePageLimit(req *request.ParsedRequest) int32 {
	return int32(request.GetIntParam(req.Parameters, "maxResults"))
}

// normaliseListLimit applies the documented default page size when a list
// request omits maxResults. Zero means unset on both protocol planes (the
// AWS SDK leaves the member absent; the admin console's optional proto
// field reads as zero), and the Core — not a calling handler — owns the
// default: an offset-paged Core that slices with a raw zero emits empty
// pages and a "0" nextToken that reproduces them forever.
func normaliseListLimit(maxResults int32) int32 {
	if maxResults == 0 {
		return sfnstore.DefaultPageSize
	}
	return maxResults
}

// computeRedriveStatus derives the AWS redriveStatus / redriveStatusReason
// for DescribeExecution from the same eligibility verdict RedriveExecution
// enforces, so the two operations cannot disagree.
//
// Per AWS spec (ExecutionRedriveStatus enum), redriveStatus is
// REDRIVABLE / NOT_REDRIVABLE / REDRIVABLE_BY_MAP_RUN: an unsuccessful
// STANDARD execution is REDRIVABLE while it is within the redrive window
// and below the history-event ceiling. For a Distributed Map, redriveStatus
// "indicates whether or not the Map Run can redrive child workflow
// executions": a child in PENDING_REDRIVE (parked while waiting for a Map
// Run concurrency slot) or in a failed terminal status is
// REDRIVABLE_BY_MAP_RUN — never REDRIVABLE, because a child "can only be
// redriven by its Map Run". Failed or timed out EXPRESS children are also
// redriven by their Map Run (restarted via StartExecution); only children
// that completed successfully are not.
//
// This is a derived field, never stored.
func computeRedriveStatus(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) (status string, reason string) {
	if exec.MapRunArn != "" {
		if exec.Status == "PENDING_REDRIVE" {
			return "REDRIVABLE_BY_MAP_RUN", ""
		}
		if !isRedrivableStatus(exec.Status) {
			return "NOT_REDRIVABLE", "Execution is in " + exec.Status + " status and cannot be redriven"
		}
		if ok, why, err := redrivePeriodChecks(ctx, store, exec); err != nil {
			return "NOT_REDRIVABLE", "Execution history could not be read to determine redrive eligibility"
		} else if !ok {
			return "NOT_REDRIVABLE", "Execution " + why
		}
		return "REDRIVABLE_BY_MAP_RUN", ""
	}
	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil || sm == nil {
		if err != nil && !errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return "NOT_REDRIVABLE", "The state machine of the execution could not be read, so redrive eligibility is unknown"
		}
		return "NOT_REDRIVABLE", "The state machine of the execution does not exist, so the execution cannot be redriven"
	}
	eligible, why, err := evaluateRedriveEligibility(ctx, store, sm, exec)
	if err != nil {
		return "NOT_REDRIVABLE", "Execution history could not be read to determine redrive eligibility"
	}
	if eligible {
		return "REDRIVABLE", ""
	}
	return "NOT_REDRIVABLE", "Execution " + why
}

func executionToResponse(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) map[string]interface{} {
	response := map[string]interface{}{
		"executionArn":    exec.ExecutionArn,
		"stateMachineArn": exec.StateMachineArn,
		"name":            exec.Name,
		"status":          exec.Status,
		"startDate":       awsEpochSeconds(exec.StartDate),
	}

	// CloudWatchEventsExecutionDataDetails carries one member: included,
	// "Always true for API calls".
	response["inputDetails"] = map[string]interface{}{"included": true}
	if exec.Output != "" {
		response["outputDetails"] = map[string]interface{}{"included": true}
	}

	if exec.Input != "" {
		response["input"] = exec.Input
	}
	if exec.Output != "" {
		response["output"] = exec.Output
	}
	if !exec.StopDate.IsZero() {
		response["stopDate"] = awsEpochSeconds(exec.StopDate)
	}
	if exec.Error != "" {
		response["error"] = exec.Error
	}
	if exec.Cause != "" {
		response["cause"] = exec.Cause
	}
	if exec.TraceHeader != "" {
		response["traceHeader"] = exec.TraceHeader
	}
	if exec.StateMachineVersionArn != "" {
		response["stateMachineVersionArn"] = exec.StateMachineVersionArn
	}
	if exec.StateMachineAliasArn != "" {
		response["stateMachineAliasArn"] = exec.StateMachineAliasArn
	}
	if exec.MapRunArn != "" {
		response["mapRunArn"] = exec.MapRunArn
	}
	if exec.RedriveCount != 0 {
		response["redriveCount"] = exec.RedriveCount
	}
	if !exec.RedriveDate.IsZero() {
		response["redriveDate"] = awsEpochSeconds(exec.RedriveDate)
	}
	if rs, reason := computeRedriveStatus(ctx, store, exec); rs != "" {
		response["redriveStatus"] = rs
		if reason != "" {
			response["redriveStatusReason"] = reason
		}
	}

	return response
}

func activityToResponse(activity *sfnstore.Activity) map[string]interface{} {
	response := map[string]interface{}{
		"activityArn":  activity.ActivityArn,
		"name":         activity.Name,
		"creationDate": awsEpochSeconds(activity.CreationDate),
	}
	if activity.EncryptionConfiguration != nil {
		response["encryptionConfiguration"] = activity.EncryptionConfiguration
	}
	return response
}

// executionDataSensitiveFields lists the fields that contain input/output
// data and should be omitted when includeExecutionData is false.
var executionDataSensitiveFields = map[string]bool{
	"input":  true,
	"output": true,
}

// jsonataOutputSlot returns pointers to a JSONata state's Output pair: the
// resolved template and the raw Output literal it lazily resolves from.
// States without an Output member return a nil pair.
func jsonataOutputSlot(state sfnstore.State) (*interface{}, *json.RawMessage) {
	switch s := state.(type) {
	case *sfnstore.TaskState:
		return &s.JSONataOutput, &s.OutputRaw
	case *sfnstore.PassState:
		return &s.JSONataOutput, &s.OutputRaw
	case *sfnstore.ParallelState:
		return &s.JSONataOutput, &s.OutputRaw
	case *sfnstore.MapState:
		return &s.JSONataOutput, &s.OutputRaw
	case *sfnstore.SucceedState:
		return &s.JSONataOutput, &s.OutputRaw
	}
	return nil, nil
}

// jsonataOutputOf returns the state's resolved JSONata Output template,
// lazily resolving it from the raw Output literal on first use and writing
// the resolution back, so later visits reuse it.
func jsonataOutputOf(state sfnstore.State) (interface{}, error) {
	output, raw := jsonataOutputSlot(state)
	if output == nil {
		return nil, nil
	}
	if *output == nil && len(*raw) > 0 {
		resolved, err := resolveJSONataOutput(state)
		if err != nil {
			return nil, err
		}
		*output = resolved
	}
	return *output, nil
}

// applyContainerStateOutput applies the output tail the Map and Parallel
// container states share: the JSONata dialect evaluates Assign and the
// Output template against the combined states variable, while the JSONPath
// dialect applies ResultSelector, ResultPath and OutputPath in order. The
// selector receives the container's combined result array as its data,
// before ResultPath folds it into the state input.
func (e *Executor) applyContainerStateOutput(ctx context.Context, execCtx *ExecutionContext, state sfnstore.State, processedInput, output string) (string, *ExecutionError) {
	if IsJSONataState(state, execCtx.QueryLanguage) {
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
		var resultData interface{}
		if err := json.Unmarshal([]byte(output), &resultData); err != nil {
			return "", &ExecutionError{ErrorCode: "States.InvalidOutput", Cause: "failed to parse output JSON"}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, resultData, nil)

		assign := containerAssign(state)
		if len(assign) > 0 {
			evaluated, err := evaluateAssign(ctx, assign, statesVar, execCtx.VariableScope)
			if err != nil {
				return "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
			}
			execCtx.PendingAssign = evaluated
		}

		jsonataOutput, err := jsonataOutputOf(state)
		if err != nil {
			return "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
		}
		if jsonataOutput == nil {
			return output, nil
		}
		resolved, err := e.applyJSONataOutput(ctx, jsonataOutput, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
		}
		outputJSON, err := json.Marshal(resolved)
		if err != nil {
			return "", e.newQueryEvalError(ctx, execCtx, "Output", fmt.Sprintf("failed to marshal: %s", err.Error()))
		}
		return string(outputJSON), nil
	}

	// JSONPath container Assign evaluates against the raw combined result
	// ("In Task, Map, Parallel states, $ refers to the API/sub-workflow
	// result"); new values become visible in the next state.
	if assign := containerAssign(state); len(assign) > 0 {
		var root interface{}
		if err := json.Unmarshal([]byte(output), &root); err != nil {
			root = nil
		}
		evaluated, err := e.evaluateJSONPathAssign(execCtx, assign, root)
		if err != nil {
			return "", newJSONPathEvalError("Assign", err)
		}
		execCtx.PendingAssign = evaluated
	}

	selector, resultPath, outputPath := containerFilters(state)
	if selector != nil {
		selected, selErr := e.applyResultSelector(execCtx, output, selector, "")
		if selErr != nil {
			return "", selErr
		}
		output = selected
	}
	if resultPath != "" {
		folded, rpErr := e.applyResultPath(processedInput, output, resultPath)
		if rpErr != nil {
			return "", rpErr
		}
		output = folded
	}
	return e.applyOutputPath(execCtx, output, outputPath)
}

// containerAssign returns the Assign block of a container state.
func containerAssign(state sfnstore.State) map[string]interface{} {
	switch s := state.(type) {
	case *sfnstore.ParallelState:
		return s.Assign
	case *sfnstore.MapState:
		return s.Assign
	}
	return nil
}

// containerFilters returns the JSONPath output filters of a container
// state: ResultSelector, ResultPath and OutputPath.
func containerFilters(state sfnstore.State) (*sfnstore.ResultSelector, string, string) {
	switch s := state.(type) {
	case *sfnstore.ParallelState:
		return s.ResultSelector, s.ResultPath, s.GetOutputPath()
	case *sfnstore.MapState:
		return s.ResultSelector, s.ResultPath, s.GetOutputPath()
	}
	return nil, "", ""
}

// stateEnteredDetails builds the generic state-entry detail shape every
// StateEntered event carries.
func stateEnteredDetails(execCtx *ExecutionContext, input string) *sfnstore.StateEnteredEventDetails {
	return &sfnstore.StateEnteredEventDetails{
		Input: input,
		Name:  execCtx.CurrentState,
	}
}

// stateExitedDetails builds the generic state-exit detail shape, carrying
// the variables the state assigned when it processed an Assign block.
func stateExitedDetails(execCtx *ExecutionContext, output string) *sfnstore.StateExitedEventDetails {
	details := &sfnstore.StateExitedEventDetails{
		Output: output,
		Name:   execCtx.CurrentState,
	}
	if len(execCtx.PendingAssign) > 0 {
		details.AssignedVariables = execCtx.PendingAssign
	}
	return details
}

func historyEventToResponse(event *sfnstore.ExecutionHistoryEvent, includeExecutionData bool) map[string]interface{} {
	response := map[string]interface{}{
		"id":              event.EventId,
		"previousEventId": event.PreviousEventId,
		"timestamp":       awsEpochSeconds(event.Timestamp),
		"type":            event.Type,
	}

	// dataDetails renders an inputDetails / outputDetails member.
	// HistoryEventExecutionDataDetails (and the assigned-variables
	// variant) carries one member, truncated, "Always false for API
	// calls"; whether the paired payload member appears at all is
	// governed by the includeExecutionData scrub below, not by this
	// detail.
	dataDetails := func() map[string]interface{} {
		return map[string]interface{}{"truncated": false}
	}

	switch event.Type {
	case "ExecutionStarted":
		if d := event.ExecutionStartedEventDetails; d != nil {
			details := map[string]interface{}{
				"input":        d.Input,
				"inputDetails": dataDetails(),
				"roleArn":      d.RoleArn,
			}
			if d.StateMachineAliasArn != "" {
				details["stateMachineAliasArn"] = d.StateMachineAliasArn
			}
			if d.StateMachineVersionArn != "" {
				details["stateMachineVersionArn"] = d.StateMachineVersionArn
			}
			response["executionStartedEventDetails"] = details
		}
	case "ExecutionSucceeded":
		if d := event.ExecutionSucceededEventDetails; d != nil {
			response["executionSucceededEventDetails"] = map[string]interface{}{
				"output":        d.Output,
				"outputDetails": dataDetails(),
			}
		}
	case "ExecutionFailed":
		if d := event.ExecutionFailedEventDetails; d != nil {
			response["executionFailedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "ExecutionAborted":
		if d := event.ExecutionAbortedEventDetails; d != nil {
			response["executionAbortedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "ExecutionTimedOut":
		if d := event.ExecutionTimedOutEventDetails; d != nil {
			response["executionTimedOutEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "ExecutionRedriven":
		if d := event.ExecutionRedrivenEventDetails; d != nil {
			response["executionRedrivenEventDetails"] = map[string]interface{}{
				"redriveCount": d.RedriveCount,
			}
		}
	case "TaskScheduled":
		if d := event.TaskScheduledEventDetails; d != nil {
			details := map[string]interface{}{
				"resource":     d.Resource,
				"resourceType": d.ResourceType,
			}
			if d.Region != "" {
				details["region"] = d.Region
			}
			if d.Parameters != nil {
				details["parameters"] = d.Parameters
			}
			if d.TimeoutInSeconds > 0 {
				details["timeoutInSeconds"] = d.TimeoutInSeconds
			}
			if d.HeartbeatInSeconds > 0 {
				details["heartbeatInSeconds"] = d.HeartbeatInSeconds
			}
			if d.TaskCredentials != nil {
				details["taskCredentials"] = map[string]interface{}{"roleArn": d.TaskCredentials.RoleArn}
			}
			response["taskScheduledEventDetails"] = details
		}
	case "TaskStarted":
		if d := event.TaskStartedEventDetails; d != nil {
			response["taskStartedEventDetails"] = map[string]interface{}{
				"resource":     d.Resource,
				"resourceType": d.ResourceType,
			}
		}
	case "TaskStartFailed":
		if d := event.TaskStartFailedEventDetails; d != nil {
			response["taskStartFailedEventDetails"] = map[string]interface{}{
				"resource":     d.Resource,
				"resourceType": d.ResourceType,
				"error":        d.Error,
				"cause":        d.Cause,
			}
		}
	case "TaskSubmitted":
		if d := event.TaskSubmittedEventDetails; d != nil {
			response["taskSubmittedEventDetails"] = map[string]interface{}{
				"resource":      d.Resource,
				"resourceType":  d.ResourceType,
				"output":        d.Output,
				"outputDetails": dataDetails(),
			}
		}
	case "TaskSubmitFailed":
		if d := event.TaskSubmitFailedEventDetails; d != nil {
			response["taskSubmitFailedEventDetails"] = map[string]interface{}{
				"resource":     d.Resource,
				"resourceType": d.ResourceType,
				"error":        d.Error,
				"cause":        d.Cause,
			}
		}
	case "TaskSucceeded":
		if d := event.TaskSucceededEventDetails; d != nil {
			response["taskSucceededEventDetails"] = map[string]interface{}{
				"resource":      d.Resource,
				"resourceType":  d.ResourceType,
				"output":        d.Output,
				"outputDetails": dataDetails(),
			}
		}
	case "TaskFailed":
		if d := event.TaskFailedEventDetails; d != nil {
			response["taskFailedEventDetails"] = map[string]interface{}{
				"resource":     d.Resource,
				"resourceType": d.ResourceType,
				"error":        d.Error,
				"cause":        d.Cause,
			}
		}
	case "TaskTimedOut":
		if d := event.TaskTimedOutEventDetails; d != nil {
			response["taskTimedOutEventDetails"] = map[string]interface{}{
				"resource":     d.Resource,
				"resourceType": d.ResourceType,
				"error":        d.Error,
				"cause":        d.Cause,
			}
		}
	case "LambdaFunctionScheduled":
		if d := event.LambdaFunctionScheduledEventDetails; d != nil {
			details := map[string]interface{}{
				"resource":     d.Resource,
				"input":        d.Input,
				"inputDetails": dataDetails(),
			}
			if d.TimeoutInSeconds > 0 {
				details["timeoutInSeconds"] = d.TimeoutInSeconds
			}
			if d.TaskCredentials != nil {
				details["taskCredentials"] = map[string]interface{}{"roleArn": d.TaskCredentials.RoleArn}
			}
			response["lambdaFunctionScheduledEventDetails"] = details
		}
	case "LambdaFunctionScheduleFailed":
		if d := event.LambdaFunctionScheduleFailedEventDetails; d != nil {
			response["lambdaFunctionScheduleFailedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "LambdaFunctionStartFailed":
		if d := event.LambdaFunctionStartFailedEventDetails; d != nil {
			response["lambdaFunctionStartFailedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "LambdaFunctionFailed":
		if d := event.LambdaFunctionFailedEventDetails; d != nil {
			response["lambdaFunctionFailedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "LambdaFunctionTimedOut":
		if d := event.LambdaFunctionTimedOutEventDetails; d != nil {
			response["lambdaFunctionTimedOutEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "LambdaFunctionSucceeded":
		if d := event.LambdaFunctionSucceededEventDetails; d != nil {
			response["lambdaFunctionSucceededEventDetails"] = map[string]interface{}{
				"output":        d.Output,
				"outputDetails": dataDetails(),
			}
		}
	case "ActivityScheduled":
		if d := event.ActivityTaskScheduledEventDetails; d != nil {
			// Member set per the model's ActivityScheduledEventDetails; the
			// task token travels in GetActivityTask's response, not in the
			// history event.
			details := map[string]interface{}{
				"resource":     d.Resource,
				"input":        d.Input,
				"inputDetails": dataDetails(),
			}
			if d.TimeoutInSeconds > 0 {
				details["timeoutInSeconds"] = d.TimeoutInSeconds
			}
			if d.HeartbeatSeconds > 0 {
				details["heartbeatInSeconds"] = d.HeartbeatSeconds
			}
			response["activityScheduledEventDetails"] = details
		}
	case "ActivityScheduleFailed":
		if d := event.ActivityScheduleFailedEventDetails; d != nil {
			response["activityScheduleFailedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "ActivityStarted":
		if d := event.ActivityTaskStartedEventDetails; d != nil {
			response["activityStartedEventDetails"] = map[string]interface{}{
				"workerName": d.WorkerName,
			}
		}
	case "ActivitySucceeded":
		if d := event.ActivityTaskSucceededEventDetails; d != nil {
			response["activitySucceededEventDetails"] = map[string]interface{}{
				"output":        d.Output,
				"outputDetails": dataDetails(),
			}
		}
	case "ActivityFailed":
		if d := event.ActivityTaskFailedEventDetails; d != nil {
			response["activityFailedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "ActivityTimedOut":
		if d := event.ActivityTaskTimedOutEventDetails; d != nil {
			response["activityTimedOutEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "PassStateEntered", "ChoiceStateEntered", "WaitStateEntered", "ParallelStateEntered",
		"MapStateEntered", "FailStateEntered", "SucceedStateEntered", "TaskStateEntered":
		if d := event.StateEnteredEventDetails; d != nil {
			response["stateEnteredEventDetails"] = map[string]interface{}{
				"input":        d.Input,
				"inputDetails": dataDetails(),
				"name":         d.Name,
			}
		}
	case "PassStateExited", "ChoiceStateExited", "WaitStateExited", "ParallelStateExited",
		"MapStateExited", "SucceedStateExited", "TaskStateExited":
		if d := event.StateExitedEventDetails; d != nil {
			details := map[string]interface{}{
				"name":          d.Name,
				"output":        d.Output,
				"outputDetails": dataDetails(),
			}
			if len(d.AssignedVariables) > 0 {
				details["assignedVariables"] = d.AssignedVariables
				varDetails := make(map[string]interface{}, len(d.AssignedVariables))
				for name := range d.AssignedVariables {
					varDetails[name] = dataDetails()
				}
				details["assignedVariablesDetails"] = varDetails
			}
			response["stateExitedEventDetails"] = details
		}
	case "MapStateStarted":
		if d := event.MapStateStartedEventDetails; d != nil {
			response["mapStateStartedEventDetails"] = map[string]interface{}{
				"length": d.Length,
			}
		}
	case "MapRunStarted":
		if d := event.MapRunStartedEventDetails; d != nil {
			response["mapRunStartedEventDetails"] = map[string]interface{}{
				"mapRunArn": d.MapRunArn,
			}
		}
	case "MapRunFailed":
		if d := event.MapRunFailedEventDetails; d != nil {
			response["mapRunFailedEventDetails"] = map[string]interface{}{
				"error": d.Error,
				"cause": d.Cause,
			}
		}
	case "MapRunRedriven":
		if d := event.MapRunRedrivenEventDetails; d != nil {
			response["mapRunRedrivenEventDetails"] = map[string]interface{}{
				"mapRunArn":    d.MapRunArn,
				"redriveCount": d.RedriveCount,
			}
		}
	case "MapIterationStarted", "MapIterationSucceeded", "MapIterationFailed", "MapIterationAborted":
		if d := event.MapIterationEventDetails; d != nil {
			// The four iteration events share the MapIterationEventDetails
			// shape; the member key carries the event's own suffix.
			key := "mapIteration" + event.Type[len("MapIteration"):] + "EventDetails"
			response[key] = map[string]interface{}{
				"index": d.Index,
				"name":  d.Name,
			}
		}
	case "EvaluationFailed":
		if d := event.EvaluationFailedEventDetails; d != nil {
			response["evaluationFailedEventDetails"] = map[string]interface{}{
				"state":    d.State,
				"cause":    d.Cause,
				"error":    d.Error,
				"location": d.Location,
			}
		}
	}

	if !includeExecutionData {
		for key, value := range response {
			if executionDataSensitiveFields[key] {
				delete(response, key)
				continue
			}
			// Every detail map the switch above serialises lives under a
			// key ending in "EventDetails"; scrubbing by suffix keeps this
			// in step with the event vocabulary without a second,
			// hand-maintained enumeration of it.
			if details, ok := value.(map[string]interface{}); ok && strings.HasSuffix(key, "EventDetails") {
				for f := range details {
					if executionDataSensitiveFields[f] {
						delete(details, f)
					}
				}
			}
		}
	}

	return response
}
