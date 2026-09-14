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
// AWS Step Functions does not re-execute states that already succeeded:
// the resume point is the state that follows the last top-level state that
// exited successfully. Branch and iteration inner states write their events
// into the parent history, so only exits outside any Parallel/Map container
// count — a failure inside a branch resumes at the container itself, which
// re-runs only the failed branches/iterations from its checkpoints. If no
// top-level state exited successfully, the execution resumes from the start
// state with the original execution input.
func determineResumePoint(ctx context.Context, store *sfnstore.StepFunctionStore, executionArn string, definition *sfnstore.StateMachineDefinition) (*resumePoint, error) {
	// The paginated history API caps any single call at the page ceiling,
	// so the marker chain is followed until the history is exhausted: a
	// clamped single call would compute the last event id and the last
	// exited state from a truncated history, and the events the resume then
	// appends would overwrite persisted ones (ids are blind Put keys).
	var events []*sfnstore.ExecutionHistoryEvent
	token := ""
	for {
		page, next, err := store.GetExecutionHistory(ctx, executionArn, sfnstore.MaxPageSize, token, false)
		if err != nil {
			return nil, fmt.Errorf("failed to load execution history: %w", err)
		}
		events = append(events, page...)
		if next == "" {
			break
		}
		token = next
	}

	var lastExitedState string
	var lastExitedOutput string
	var lastExitedNextState string
	var lastEventId int64
	containerDepth := 0

	for _, event := range events {
		if event.EventId > lastEventId {
			lastEventId = event.EventId
		}
		switch event.Type {
		case "ParallelStateEntered", "MapStateEntered":
			containerDepth++
		case "ParallelStateExited", "MapStateExited":
			if containerDepth > 0 {
				containerDepth--
			}
		}
		if containerDepth > 0 {
			continue
		}

		name, output := extractStateExitInfo(event)
		if name != "" {
			lastExitedState = name
			lastExitedOutput = output
			// A Choice exit records the transition it took; the resume
			// follows that record instead of consulting the definition,
			// whose Choices rules would re-evaluate against data the
			// rerun no longer holds.
			lastExitedNextState = extractStateExitNextState(event)
		}
	}

	if lastExitedState == "" {
		return &resumePoint{
			StateName:   definition.StartAt,
			Input:       "",
			LastEventId: lastEventId,
		}, nil
	}

	if lastExitedNextState != "" {
		return &resumePoint{
			StateName:   lastExitedNextState,
			Input:       lastExitedOutput,
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
// event, or empty strings if the event is not a StateExited event. Choice
// exits additionally carry the taken transition on the details' internal
// NextState member; determineResumePoint reads it through
// extractStateExitNextState.
func extractStateExitInfo(event *sfnstore.ExecutionHistoryEvent) (string, string) {
	switch event.Type {
	case "PassStateExited", "TaskStateExited", "ChoiceStateExited", "WaitStateExited",
		"MapStateExited", "ParallelStateExited", "SucceedStateExited":
		if d := event.StateExitedEventDetails; d != nil {
			return d.Name, d.Output
		}
	}
	return "", ""
}

// extractStateExitNextState returns the transition a Choice state recorded
// on its exit event, or the empty string for other exits.
func extractStateExitNextState(event *sfnstore.ExecutionHistoryEvent) string {
	if event.Type == "ChoiceStateExited" && event.StateExitedEventDetails != nil {
		return event.StateExitedEventDetails.NextState
	}
	return ""
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
