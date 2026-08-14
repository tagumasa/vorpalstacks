package sfn

import (
	"strconv"

	"vorpalstacks/internal/common/request"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// parsePageLimit parses the maxResults parameter and validates it against
// the AWS spec range [0, 1000]. Returns the default of 100 when not
// provided. Returns an error if the value is outside the range.
func parsePageLimit(req *request.ParsedRequest) (int32, error) {
	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit == 0 {
		return sfnstore.DefaultPageSize, nil
	}
	if limit < 0 || limit > sfnstore.MaxPageSize {
		return 0, NewValidationException("maxResults must be in [0, 1000], got " + strconv.FormatInt(int64(limit), 10))
	}
	return limit, nil
}

// computeRedriveStatus derives the AWS redriveStatus / redriveStatusReason
// for DescribeExecution from the current execution state.
//
// Per AWS spec (ExecutionRedriveStatus enum), redriveStatus is
// REDRIVABLE / NOT_REDRIVABLE / REDRIVABLE_BY_MAP_RUN. It mirrors what
// RedriveExecution would return: NOT_REDRIVABLE with a reason when the
// execution cannot be redriven.
//
// This is a derived field, never stored.
func computeRedriveStatus(exec *sfnstore.Execution) (status string, reason string) {
	switch exec.Status {
	case "FAILED", "TIMED_OUT", "PENDING_REDRIVE":
		return "REDRIVABLE", ""
	case "RUNNING":
		return "NOT_REDRIVABLE", "Execution is RUNNING and cannot be redriven"
	case "SUCCEEDED":
		return "NOT_REDRIVABLE", "Execution is SUCCEEDED and cannot be redriven"
	case "ABORTED":
		return "NOT_REDRIVABLE", "Execution is ABORTED and cannot be redriven"
	}
	return "", ""
}

func executionToResponse(exec *sfnstore.Execution) map[string]interface{} {
	response := map[string]interface{}{
		"executionArn":    exec.ExecutionArn,
		"stateMachineArn": exec.StateMachineArn,
		"name":            exec.Name,
		"status":          exec.Status,
		"startDate":       exec.StartDate.Unix(),
	}

	if exec.InputDetails != nil {
		response["inputDetails"] = map[string]interface{}{
			"included": exec.InputDetails.Included,
			"type":     exec.InputDetails.Type,
		}
	}
	if exec.OutputDetails != nil {
		response["outputDetails"] = map[string]interface{}{
			"included": exec.OutputDetails.Included,
			"type":     exec.OutputDetails.Type,
		}
	}

	if exec.Input != "" {
		response["input"] = exec.Input
	}
	if exec.Output != "" {
		response["output"] = exec.Output
	}
	if !exec.StopDate.IsZero() {
		response["stopDate"] = exec.StopDate.Unix()
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
	if exec.ItemCount != 0 {
		response["itemCount"] = exec.ItemCount
	}
	if exec.RedriveCount != 0 {
		response["redriveCount"] = exec.RedriveCount
	}
	if !exec.RedriveDate.IsZero() {
		response["redriveDate"] = exec.RedriveDate.Unix()
	}
	if rs, reason := computeRedriveStatus(exec); rs != "" {
		response["redriveStatus"] = rs
		if reason != "" {
			response["redriveStatusReason"] = reason
		}
	}

	return response
}

func stateMachineToResponse(sm *sfnstore.StateMachine) map[string]interface{} {
	response := map[string]interface{}{
		"stateMachineArn": sm.StateMachineArn,
		"name":            sm.Name,
		"type":            sm.Type,
		"creationDate":    sm.CreationDate.Unix(),
	}

	if sm.Definition != "" {
		response["definition"] = sm.Definition
	}
	if sm.RoleArn != "" {
		response["roleArn"] = sm.RoleArn
	}
	if sm.Status != "" {
		response["status"] = sm.Status
	}
	if !sm.UpdateDate.IsZero() {
		response["updateDate"] = sm.UpdateDate.Unix()
	}
	if sm.RevisionId != "" {
		response["revisionId"] = sm.RevisionId
	}
	if sm.Label != "" {
		response["label"] = sm.Label
	}
	if sm.LoggingConfiguration != nil {
		response["loggingConfiguration"] = sm.LoggingConfiguration
	}
	if sm.EncryptionConfiguration != nil {
		response["encryptionConfiguration"] = sm.EncryptionConfiguration
	}
	if sm.TracingConfiguration != nil {
		response["tracingConfiguration"] = sm.TracingConfiguration
	}

	return response
}

func activityToResponse(activity *sfnstore.Activity) map[string]interface{} {
	response := map[string]interface{}{
		"activityArn":  activity.ActivityArn,
		"name":         activity.Name,
		"creationDate": activity.CreationDate.Unix(),
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

func historyEventToResponse(event *sfnstore.ExecutionHistoryEvent, includeExecutionData bool) map[string]interface{} {
	response := map[string]interface{}{
		"id":              event.EventId,
		"previousEventId": event.PreviousEventId,
		"timestamp":       event.Timestamp.Unix(),
		"type":            event.Type,
	}

	switch event.Type {
	case "ExecutionStarted":
		if event.ExecutionStartedEventDetails != nil {
			response["executionStartedEventDetails"] = map[string]interface{}{
				"input":           event.ExecutionStartedEventDetails.Input,
				"roleArn":         event.ExecutionStartedEventDetails.RoleArn,
				"stateMachineArn": event.ExecutionStartedEventDetails.StateMachineArn,
				"name":            event.ExecutionStartedEventDetails.Name,
			}
		}
	case "ExecutionSucceeded":
		if event.ExecutionSucceededEventDetails != nil {
			response["executionSucceededEventDetails"] = map[string]interface{}{
				"output": event.ExecutionSucceededEventDetails.Output,
			}
		}
	case "ExecutionFailed":
		if event.ExecutionFailedEventDetails != nil {
			response["executionFailedEventDetails"] = map[string]interface{}{
				"error": event.ExecutionFailedEventDetails.Error,
				"cause": event.ExecutionFailedEventDetails.Cause,
			}
		}
	case "TaskStarted":
		if event.TaskStartedEventDetails != nil {
			response["taskStartedEventDetails"] = map[string]interface{}{
				"resource": event.TaskStartedEventDetails.Resource,
				"input":    event.TaskStartedEventDetails.Input,
			}
		}
	case "TaskSucceeded":
		if event.TaskSucceededEventDetails != nil {
			response["taskSucceededEventDetails"] = map[string]interface{}{
				"resource": event.TaskSucceededEventDetails.Resource,
				"output":   event.TaskSucceededEventDetails.Output,
			}
		}
	case "TaskFailed":
		if event.TaskFailedEventDetails != nil {
			response["taskFailedEventDetails"] = map[string]interface{}{
				"resource": event.TaskFailedEventDetails.Resource,
				"error":    event.TaskFailedEventDetails.Error,
				"cause":    event.TaskFailedEventDetails.Cause,
			}
		}
	case "PassStateEntered":
		if event.PassStateEnteredEventDetails != nil {
			response["passStateEnteredEventDetails"] = map[string]interface{}{
				"input": event.PassStateEnteredEventDetails.Input,
				"name":  event.PassStateEnteredEventDetails.Name,
			}
		}
	case "PassStateExited":
		if event.PassStateExitedEventDetails != nil {
			response["passStateExitedEventDetails"] = map[string]interface{}{
				"output": event.PassStateExitedEventDetails.Output,
				"name":   event.PassStateExitedEventDetails.Name,
			}
		}
	case "ChoiceStateEntered":
		if event.ChoiceStateEnteredEventDetails != nil {
			response["choiceStateEnteredEventDetails"] = map[string]interface{}{
				"input": event.ChoiceStateEnteredEventDetails.Input,
				"name":  event.ChoiceStateEnteredEventDetails.Name,
			}
		}
	case "ChoiceStateExited":
		if event.ChoiceStateExitedEventDetails != nil {
			response["choiceStateExitedEventDetails"] = map[string]interface{}{
				"output":    event.ChoiceStateExitedEventDetails.Output,
				"name":      event.ChoiceStateExitedEventDetails.Name,
				"nextState": event.ChoiceStateExitedEventDetails.NextState,
			}
		}
	case "WaitStateEntered":
		if event.WaitStateEnteredEventDetails != nil {
			response["waitStateEnteredEventDetails"] = map[string]interface{}{
				"input":     event.WaitStateEnteredEventDetails.Input,
				"name":      event.WaitStateEnteredEventDetails.Name,
				"seconds":   event.WaitStateEnteredEventDetails.Seconds,
				"timestamp": event.WaitStateEnteredEventDetails.Timestamp,
			}
		}
	case "WaitStateExited":
		if event.WaitStateExitedEventDetails != nil {
			response["waitStateExitedEventDetails"] = map[string]interface{}{
				"output": event.WaitStateExitedEventDetails.Output,
				"name":   event.WaitStateExitedEventDetails.Name,
			}
		}
	case "ParallelStateEntered":
		if event.ParallelStateEnteredEventDetails != nil {
			response["parallelStateEnteredEventDetails"] = map[string]interface{}{
				"input":    event.ParallelStateEnteredEventDetails.Input,
				"name":     event.ParallelStateEnteredEventDetails.Name,
				"branches": event.ParallelStateEnteredEventDetails.Branches,
			}
		}
	case "ParallelStateExited":
		if event.ParallelStateExitedEventDetails != nil {
			response["parallelStateExitedEventDetails"] = map[string]interface{}{
				"output": event.ParallelStateExitedEventDetails.Output,
				"name":   event.ParallelStateExitedEventDetails.Name,
			}
		}
	case "MapStateEntered":
		if event.MapStateEnteredEventDetails != nil {
			response["mapStateEnteredEventDetails"] = map[string]interface{}{
				"input":          event.MapStateEnteredEventDetails.Input,
				"name":           event.MapStateEnteredEventDetails.Name,
				"itemsProcessed": event.MapStateEnteredEventDetails.ItemsProcessed,
				"itemsFailed":    event.MapStateEnteredEventDetails.ItemsFailed,
			}
		}
	case "MapStateExited":
		if event.MapStateExitedEventDetails != nil {
			response["mapStateExitedEventDetails"] = map[string]interface{}{
				"output":         event.MapStateExitedEventDetails.Output,
				"name":           event.MapStateExitedEventDetails.Name,
				"itemsProcessed": event.MapStateExitedEventDetails.ItemsProcessed,
				"itemsFailed":    event.MapStateExitedEventDetails.ItemsFailed,
			}
		}
	case "FailStateEntered":
		if event.FailStateEnteredEventDetails != nil {
			response["failStateEnteredEventDetails"] = map[string]interface{}{
				"input": event.FailStateEnteredEventDetails.Input,
				"name":  event.FailStateEnteredEventDetails.Name,
				"error": event.FailStateEnteredEventDetails.Error,
				"cause": event.FailStateEnteredEventDetails.Cause,
			}
		}
	case "SucceedStateEntered":
		if event.SucceedStateEnteredEventDetails != nil {
			response["succeedStateEnteredEventDetails"] = map[string]interface{}{
				"input": event.SucceedStateEnteredEventDetails.Input,
				"name":  event.SucceedStateEnteredEventDetails.Name,
			}
		}
	case "ActivityTaskScheduled":
		if event.ActivityTaskScheduledEventDetails != nil {
			response["activityTaskScheduledEventDetails"] = map[string]interface{}{
				"resource":  event.ActivityTaskScheduledEventDetails.Resource,
				"input":     event.ActivityTaskScheduledEventDetails.Input,
				"taskToken": event.ActivityTaskScheduledEventDetails.TaskToken,
			}
		}
	case "ActivityTaskStarted":
		if event.ActivityTaskStartedEventDetails != nil {
			response["activityTaskStartedEventDetails"] = map[string]interface{}{
				"workerName": event.ActivityTaskStartedEventDetails.WorkerName,
			}
		}
	case "ActivityTaskSucceeded":
		if event.ActivityTaskSucceededEventDetails != nil {
			response["activityTaskSucceededEventDetails"] = map[string]interface{}{
				"output": event.ActivityTaskSucceededEventDetails.Output,
			}
		}
	case "ActivityTaskFailed":
		if event.ActivityTaskFailedEventDetails != nil {
			response["activityTaskFailedEventDetails"] = map[string]interface{}{
				"error": event.ActivityTaskFailedEventDetails.Error,
				"cause": event.ActivityTaskFailedEventDetails.Cause,
			}
		}
	case "ActivityTaskTimedOut":
		if event.ActivityTaskTimedOutEventDetails != nil {
			response["activityTaskTimedOutEventDetails"] = map[string]interface{}{
				"error": event.ActivityTaskTimedOutEventDetails.Error,
				"cause": event.ActivityTaskTimedOutEventDetails.Cause,
			}
		}
	case "EvaluationFailed":
		if event.EvaluationFailedEventDetails != nil {
			response["evaluationFailedEventDetails"] = map[string]interface{}{
				"state":    event.EvaluationFailedEventDetails.State,
				"cause":    event.EvaluationFailedEventDetails.Cause,
				"error":    event.EvaluationFailedEventDetails.Error,
				"location": event.EvaluationFailedEventDetails.Location,
			}
		}
	case "TaskStateEntered":
		if event.TaskStateEnteredEventDetails != nil {
			response["taskStateEnteredEventDetails"] = map[string]interface{}{
				"input": event.TaskStateEnteredEventDetails.Input,
				"name":  event.TaskStateEnteredEventDetails.Name,
			}
		}
	case "TaskStateExited":
		if event.TaskStateExitedEventDetails != nil {
			response["taskStateExitedEventDetails"] = map[string]interface{}{
				"output": event.TaskStateExitedEventDetails.Output,
				"name":   event.TaskStateExitedEventDetails.Name,
			}
		}
	case "ExecutionRedriven":
		if event.ExecutionRedrivedEventDetails != nil {
			response["executionRedrivedEventDetails"] = map[string]interface{}{
				"redriveDate":     event.ExecutionRedrivedEventDetails.RedriveDate.Unix(),
				"stateMachineArn": event.ExecutionRedrivedEventDetails.StateMachineArn,
				"executionArn":    event.ExecutionRedrivedEventDetails.ExecutionArn,
			}
		}
	}

	if !includeExecutionData {
		for key := range response {
			if executionDataSensitiveFields[key] {
				delete(response, key)
			}
		}
		for _, detailsKey := range []string{
			"executionStartedEventDetails",
			"executionSucceededEventDetails",
			"taskStartedEventDetails",
			"taskSucceededEventDetails",
			"taskFailedEventDetails",
			"passStateEnteredEventDetails",
			"passStateExitedEventDetails",
			"choiceStateEnteredEventDetails",
			"choiceStateExitedEventDetails",
			"waitStateEnteredEventDetails",
			"waitStateExitedEventDetails",
			"parallelStateEnteredEventDetails",
			"parallelStateExitedEventDetails",
			"mapStateEnteredEventDetails",
			"mapStateExitedEventDetails",
			"failStateEnteredEventDetails",
			"succeedStateEnteredEventDetails",
			"activityTaskScheduledEventDetails",
			"activityTaskSucceededEventDetails",
			"taskStateEnteredEventDetails",
			"taskStateExitedEventDetails",
			"executionRedrivedEventDetails",
		} {
			if details, ok := response[detailsKey].(map[string]interface{}); ok {
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
