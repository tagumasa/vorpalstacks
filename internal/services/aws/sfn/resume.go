package sfn

import (
	"context"
	"encoding/json"
	"fmt"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// resumePoint identifies the state from which a redriven execution should
// resume, along with the input that state should receive.
type resumePoint struct {
	StateName   string
	Input       string
	LastEventId int64
}

// determineResumePoint inspects the execution history to decide where a
// redriven execution should resume.
//
// AWS Step Functions does not re-execute states that already succeeded.
// The resume point is the first state that was entered but never exited
// (i.e. the state that was running when the execution failed). Its input
// is the output of the last successfully-exited state. If no state exited
// successfully, the execution resumes from the start state with the
// original execution input.
func determineResumePoint(ctx context.Context, store *sfnstore.StepFunctionStore, executionArn string, definition *sfnstore.StateMachineDefinition) (*resumePoint, error) {
	events, _, err := store.GetExecutionHistory(ctx, executionArn, 100000, "", false)
	if err != nil {
		return nil, fmt.Errorf("failed to load execution history: %w", err)
	}

	var lastExitedState string
	var lastExitedOutput string
	var lastEventId int64

	for _, event := range events {
		if event.EventId > lastEventId {
			lastEventId = event.EventId
		}

		name, output := extractStateExitInfo(event)
		if name != "" {
			lastExitedState = name
			lastExitedOutput = output
		}
	}

	if lastExitedState == "" {
		return &resumePoint{
			StateName:   definition.StartAt,
			Input:       "",
			LastEventId: lastEventId,
		}, nil
	}

	nextState := lookupNextState(definition, lastExitedState)
	if nextState == "" {
		return &resumePoint{
			StateName:   definition.StartAt,
			Input:       "",
			LastEventId: lastEventId,
		}, nil
	}

	return &resumePoint{
		StateName:   nextState,
		Input:       lastExitedOutput,
		LastEventId: lastEventId,
	}, nil
}

// extractStateExitInfo returns the state name and output from a StateExited
// event, or empty strings if the event is not a StateExited event.
func extractStateExitInfo(event *sfnstore.ExecutionHistoryEvent) (string, string) {
	switch event.Type {
	case "PassStateExited":
		if d := event.PassStateExitedEventDetails; d != nil {
			return d.Name, d.Output
		}
	case "TaskStateExited":
		if d := event.TaskStateExitedEventDetails; d != nil {
			return d.Name, d.Output
		}
	case "ChoiceStateExited":
		if d := event.ChoiceStateExitedEventDetails; d != nil {
			return d.Name, d.Output
		}
	case "WaitStateExited":
		if d := event.WaitStateExitedEventDetails; d != nil {
			return d.Name, d.Output
		}
	case "MapStateExited":
		if d := event.MapStateExitedEventDetails; d != nil {
			return d.Name, d.Output
		}
	case "ParallelStateExited":
		if d := event.ParallelStateExitedEventDetails; d != nil {
			return d.Name, d.Output
		}
	}
	return "", ""
}

// lookupNextState finds the Next field of a state in the definition. Returns
// empty string if the state has End=true or does not exist.
func lookupNextState(definition *sfnstore.StateMachineDefinition, stateName string) string {
	if definition == nil {
		return ""
	}
	raw, ok := definition.States[stateName]
	if !ok {
		return ""
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return ""
	}
	if isEnd, ok := m["End"].(bool); ok && isEnd {
		return ""
	}
	next, _ := m["Next"].(string)
	return next
}

// parseStateMachineDefinition deserialises the JSON definition string stored
// on a StateMachine record into a typed StateMachineDefinition.
func parseStateMachineDefinition(definitionJSON string) (*sfnstore.StateMachineDefinition, error) {
	var def sfnstore.StateMachineDefinition
	if err := json.Unmarshal([]byte(definitionJSON), &def); err != nil {
		return nil, fmt.Errorf("invalid state machine definition JSON: %w", err)
	}
	return &def, nil
}
