package sfn

import (
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func executionToResponse(exec *sfnstore.Execution) map[string]interface{} {
	response := map[string]interface{}{
		"executionArn":    exec.ExecutionArn,
		"stateMachineArn": exec.StateMachineArn,
		"name":            exec.Name,
		"status":          exec.Status,
		"startDate":       exec.StartDate.Unix(),
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
	if sm.Description != "" {
		response["description"] = sm.Description
	}
	if sm.Status != "" {
		response["status"] = sm.Status
	}
	if !sm.UpdateDate.IsZero() {
		response["updateDate"] = sm.UpdateDate.Unix()
	}

	return response
}

func activityToResponse(activity *sfnstore.Activity) map[string]interface{} {
	return map[string]interface{}{
		"activityArn":  activity.ActivityArn,
		"name":         activity.Name,
		"creationDate": activity.CreationDate.Unix(),
	}
}

func historyEventToResponse(event *sfnstore.ExecutionHistoryEvent) map[string]interface{} {
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
	}

	return response
}
