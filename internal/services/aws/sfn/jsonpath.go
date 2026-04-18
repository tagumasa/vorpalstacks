package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"

	"vorpalstacks/internal/config"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

func getJSONPathValueRaw(data map[string]interface{}, path string) (interface{}, bool) {
	if path == "" || path == "$" {
		return data, true
	}

	path = strings.TrimPrefix(path, "$.")
	parts := strings.Split(path, ".")

	current := interface{}(data)
	for _, part := range parts {
		if part == "" {
			continue
		}
		if m, ok := current.(map[string]interface{}); ok {
			var exists bool
			current, exists = m[part]
			if !exists {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	return current, true
}

func getJSONPathValue(data map[string]interface{}, path string) (interface{}, error) {
	if path == "$" || path == "" {
		return data, nil
	}

	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid JSONPath: %s", path)
	}

	current := interface{}(data)
	for _, part := range parts[1:] {
		if m, ok := current.(map[string]interface{}); ok {
			var exists bool
			current, exists = m[part]
			if !exists {
				return nil, fmt.Errorf("path not found: %s", part)
			}
		} else {
			return nil, fmt.Errorf("cannot access property on non-object")
		}
	}

	return current, nil
}

func (e *Executor) getContextValue(path string) interface{} {
	path = strings.TrimPrefix(path, "$$.")
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return nil
	}

	switch parts[0] {
	case "Execution":
		if len(parts) < 2 {
			return nil
		}
		switch parts[1] {
		case "Id":
			return e.extractExecutionId()
		case "Name":
			return e.extractExecutionName()
		case "RoleArn":
			return e.extractExecutionRoleArn()
		case "StartTime":
			return time.Now().UTC().Format(time.RFC3339)
		}
	case "StateMachine":
		if len(parts) < 2 {
			return nil
		}
		switch parts[1] {
		case "Id":
			return e.extractStateMachineId()
		case "Name":
			return e.extractStateMachineName()
		}
	case "Task":
		if len(parts) < 2 {
			return nil
		}
		switch parts[1] {
		case "Token":
			return fmt.Sprintf("task-token-%d", time.Now().UnixNano())
		}
	}
	return nil
}

func (e *Executor) extractExecutionId() string {
	if e.currentExecution != nil {
		return e.currentExecution.ExecutionArn
	}
	return "execution-id"
}

func (e *Executor) extractExecutionName() string {
	if e.currentExecution != nil {
		return e.currentExecution.Name
	}
	return "execution-name"
}

func (e *Executor) extractExecutionRoleArn() string {
	if e.currentRoleArn != "" {
		return e.currentRoleArn
	}
	if e.accountID != "" {
		return arnutil.NewARNBuilder(e.accountID, "").IAM().Role("StepFunctionsExecutionRole")
	}
	return arnutil.NewARNBuilder(config.AWSAccountID(), "").IAM().Role("StepFunctionsExecutionRole")
}

func (e *Executor) extractStateMachineId() string {
	if e.currentStateMachine != nil {
		return e.currentStateMachine.StateMachineArn
	}
	if e.currentExecution != nil {
		return e.currentExecution.StateMachineArn
	}
	return "state-machine-id"
}

func (e *Executor) extractStateMachineName() string {
	if e.currentStateMachine != nil {
		return e.currentStateMachine.Name
	}
	return "state-machine-name"
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case json.Number:
		f, err := val.Float64()
		return f, err == nil
	}
	return 0, false
}

func (e *Executor) parseDefinition(ctx context.Context, stateMachineArn string) (*sfnstore.StateMachineDefinition, map[string]sfnstore.State, error) {
	sm, err := e.store.GetStateMachine(ctx, stateMachineArn)
	if err != nil {
		return nil, nil, err
	}

	var definition sfnstore.StateMachineDefinition
	if err := json.Unmarshal([]byte(sm.Definition), &definition); err != nil {
		return nil, nil, fmt.Errorf("failed to parse state machine definition: %w", err)
	}

	states := make(map[string]sfnstore.State)
	for name, stateData := range definition.States {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("invalid state format for %s", name)
		}

		stateType, ok := stateMap["Type"].(string)
		if !ok {
			return nil, nil, fmt.Errorf("state type not specified for %s", name)
		}

		var state sfnstore.State
		switch stateType {
		case "Pass":
			state = &sfnstore.PassState{Name: name}
		case "Task":
			state = &sfnstore.TaskState{Name: name}
		case "Choice":
			state = &sfnstore.ChoiceState{Name: name}
		case "Wait":
			state = &sfnstore.WaitState{Name: name}
		case "Parallel":
			state = &sfnstore.ParallelState{Name: name}
		case "Map":
			state = &sfnstore.MapState{Name: name}
		case "Fail":
			state = &sfnstore.FailState{Name: name}
		case "Succeed":
			state = &sfnstore.SucceedState{Name: name}
		default:
			return nil, nil, fmt.Errorf("unknown state type: %s", stateType)
		}

		stateJSON, err := json.Marshal(stateMap)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal state data: %w", err)
		}

		if err := json.Unmarshal(stateJSON, state); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal state: %w", err)
		}

		if err := resolveOutputRaw(state, stateMap); err != nil {
			return nil, nil, fmt.Errorf("failed to resolve output for %s: %w", name, err)
		}

		states[name] = state
	}

	return &definition, states, nil
}

func resolveOutputRaw(state sfnstore.State, stateMap map[string]interface{}) error {
	outputVal, hasOutput := stateMap["Output"]
	if !hasOutput {
		return nil
	}

	outputJSON, err := json.Marshal(outputVal)
	if err != nil {
		return nil
	}

	switch s := state.(type) {
	case *sfnstore.PassState:
		s.OutputRaw = outputJSON
	case *sfnstore.TaskState:
		s.OutputRaw = outputJSON
	case *sfnstore.ParallelState:
		s.OutputRaw = outputJSON
	case *sfnstore.MapState:
		s.OutputRaw = outputJSON
	case *sfnstore.SucceedState:
		s.OutputRaw = outputJSON
	}
	return nil
}

func resolveJSONataOutput(state sfnstore.State) (interface{}, error) {
	var outputRaw json.RawMessage
	switch s := state.(type) {
	case *sfnstore.PassState:
		outputRaw = s.OutputRaw
	case *sfnstore.TaskState:
		outputRaw = s.OutputRaw
	case *sfnstore.ParallelState:
		outputRaw = s.OutputRaw
	case *sfnstore.MapState:
		outputRaw = s.OutputRaw
	case *sfnstore.SucceedState:
		outputRaw = s.OutputRaw
	default:
		return nil, nil
	}

	if len(outputRaw) == 0 {
		return nil, nil
	}

	var result interface{}
	if err := json.Unmarshal(outputRaw, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Executor) extractStatesFromDefinition(definition *sfnstore.StateMachineDefinition) (map[string]sfnstore.State, error) {
	states := make(map[string]sfnstore.State)
	for name, stateData := range definition.States {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid state format for %s", name)
		}

		stateType, ok := stateMap["Type"].(string)
		if !ok {
			return nil, fmt.Errorf("state type not specified for %s", name)
		}

		var state sfnstore.State
		switch stateType {
		case "Pass":
			state = &sfnstore.PassState{Name: name}
		case "Task":
			state = &sfnstore.TaskState{Name: name}
		case "Choice":
			state = &sfnstore.ChoiceState{Name: name}
		case "Wait":
			state = &sfnstore.WaitState{Name: name}
		case "Parallel":
			state = &sfnstore.ParallelState{Name: name}
		case "Map":
			state = &sfnstore.MapState{Name: name}
		case "Fail":
			state = &sfnstore.FailState{Name: name}
		case "Succeed":
			state = &sfnstore.SucceedState{Name: name}
		default:
			return nil, fmt.Errorf("unknown state type: %s", stateType)
		}

		stateJSON, err := json.Marshal(stateMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal state data: %w", err)
		}

		if err := json.Unmarshal(stateJSON, state); err != nil {
			return nil, fmt.Errorf("failed to unmarshal state: %w", err)
		}

		if err := resolveOutputRaw(state, stateMap); err != nil {
			return nil, fmt.Errorf("failed to resolve output for %s: %w", name, err)
		}

		states[name] = state
	}

	return states, nil
}

func (e *Executor) parseState(name string, stateData interface{}) (sfnstore.State, error) {
	stateMap, ok := stateData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid state format for %s", name)
	}

	stateType, ok := stateMap["Type"].(string)
	if !ok {
		return nil, fmt.Errorf("state type not specified for %s", name)
	}

	var state sfnstore.State
	switch stateType {
	case "Pass":
		state = &sfnstore.PassState{Name: name}
	case "Task":
		state = &sfnstore.TaskState{Name: name}
	case "Choice":
		state = &sfnstore.ChoiceState{Name: name}
	case "Wait":
		state = &sfnstore.WaitState{Name: name}
	case "Parallel":
		state = &sfnstore.ParallelState{Name: name}
	case "Map":
		state = &sfnstore.MapState{Name: name}
	case "Fail":
		state = &sfnstore.FailState{Name: name}
	case "Succeed":
		state = &sfnstore.SucceedState{Name: name}
	default:
		return nil, fmt.Errorf("unknown state type: %s", stateType)
	}

	stateJSON, err := json.Marshal(stateMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state data: %w", err)
	}

	if err := json.Unmarshal(stateJSON, state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	if err := resolveOutputRaw(state, stateMap); err != nil {
		return nil, fmt.Errorf("failed to resolve output for %s: %w", name, err)
	}

	return state, nil
}
