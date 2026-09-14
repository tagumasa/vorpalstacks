package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	sfnstore "vorpalstacks/internal/store/aws/sfn"

	"vorpalstacks/internal/config"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// getJSONPathValueRaw resolves a dotted JSONPath against an object root,
// reporting misses as a boolean. Object members resolve by name exactly as
// on the any-root resolver; the map signature serves the callers whose
// input is already typed as an object.
func getJSONPathValueRaw(data map[string]interface{}, path string) (interface{}, bool) {
	return getJSONPathValueAny(data, path)
}

// getJSONPathValueAny resolves a dotted JSONPath against any root value:
// object members by name and array elements by bare or bracketed numeric
// segment ($.0, $.[0], $.results.[1].total). Payload templates whose data is
// an array — the Parallel and Map combined result, a Task result that is an
// array — need the numeric form; object roots behave exactly like
// getJSONPathValueRaw.
func getJSONPathValueAny(data interface{}, path string) (interface{}, bool) {
	if path == "" || path == "$" {
		return data, true
	}

	path = strings.TrimPrefix(path, "$.")
	parts := strings.Split(path, ".")
	current := data
	for _, part := range parts {
		if part == "" {
			continue
		}
		key := part
		if strings.HasPrefix(key, "[") && strings.HasSuffix(key, "]") {
			key = key[1 : len(key)-1]
		}
		switch node := current.(type) {
		case map[string]interface{}:
			value, exists := node[key]
			if !exists {
				return nil, false
			}
			current = value
		case []interface{}:
			index, err := strconv.Atoi(key)
			if err != nil || index < 0 || index >= len(node) {
				return nil, false
			}
			current = node[index]
		default:
			return nil, false
		}
	}

	return current, true
}

// getJSONPathValue resolves a $-rooted dotted JSONPath against an object
// root, reporting misses as errors for the callers that surface them in
// runtime diagnostics. A path without the $ root is rejected outright: the
// definition validator admits only $/$$-rooted paths, so a bare name here
// is an internal hand-off defect, not a user path to resolve leniently.
func getJSONPathValue(data map[string]interface{}, path string) (interface{}, error) {
	if !strings.HasPrefix(path, "$") {
		return nil, fmt.Errorf("invalid JSONPath: %s", path)
	}
	value, found := getJSONPathValueAny(data, path)
	if !found {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	return value, nil
}

// isContextPath reports whether a path addresses the context object. The
// Amazon States Language defines $$-rooted paths by stripping the first
// dollar sign and applying the remainder as the JSONPath against the
// context object, so both "$$.Node" and the bare "$$" (the root path)
// address the context object.
func isContextPath(path string) bool {
	return path == "$$" || strings.HasPrefix(path, "$$.")
}

// getContextValue resolves a $$. context-object path against the same
// Context object the JSONata plane exposes as $states.context, so the two
// dialects cannot diverge on the node set: Execution (Id, Input, Name,
// RoleArn, StartTime, RedriveCount, RedriveTime), State (EnteredTime,
// Name, RetryCount), StateMachine (Id, Name), Task (Token) and, inside a
// Map iteration, Map.Item (Index, Key, Value, Source). taskToken carries
// the token minted for the current attempt (empty everywhere else);
// threading it explicitly keeps concurrent Map/Parallel branches and
// retry attempts independent. Referencing $$.Task.Token without a backing
// token fails the evaluation: a fabricated random token would match no
// stored task record and send the worker an unanswerable token, whereas
// the runtime error surfaces the definition mistake the way Step
// Functions reports an unavailable context node.
func (e *Executor) getContextValue(execCtx *ExecutionContext, taskToken string, path string) (interface{}, error) {
	if taskToken == "" && strings.HasPrefix(path, "$$.Task") {
		return nil, errTaskTokenUnavailable
	}

	ctxObj := e.buildContextObject(execCtx)
	if taskToken != "" {
		taskSection, _ := ctxObj["Task"].(map[string]interface{})
		if taskSection == nil {
			taskSection = map[string]interface{}{}
			ctxObj["Task"] = taskSection
		}
		taskSection["Token"] = taskToken
	}

	if path == "$$" {
		// The bare context root selects the whole context object: the
		// language strips the first dollar sign of a $$-rooted path and
		// applies the remainder — the root path — to the context object.
		return ctxObj, nil
	}
	rest := strings.TrimPrefix(path, "$$.")
	if rest == "" {
		return ctxObj, nil
	}
	value, found := getJSONPathValueAny(ctxObj, "$."+rest)
	if !found {
		return nil, fmt.Errorf("context object path %s selected no value", path)
	}
	return value, nil
}

// errTaskTokenUnavailable marks a $$.Task.Token reference evaluated outside
// a task attempt that carries a token (non-activity tasks, Pass or Map
// states).
var errTaskTokenUnavailable = errors.New("the Task.Token context object is not available in this state")

func (e *Executor) extractExecutionRoleArn() string {
	if e.currentRoleArn != "" {
		return e.currentRoleArn
	}
	if e.accountID != "" {
		return arnutil.NewARNBuilder(e.accountID, "").IAM().Role("StepFunctionsExecutionRole")
	}
	return arnutil.NewARNBuilder(config.AWSAccountID(), "").IAM().Role("StepFunctionsExecutionRole")
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

// definitionForExecution returns the definition text the execution must
// run: the pinned version snapshot when the execution was started with a
// version or alias ARN, the live state machine definition otherwise. The
// empty string reports that no definition could be resolved; a storage
// fault other than absence surfaces as itself instead of masquerading as
// a missing definition.
func (e *Executor) definitionForExecution(ctx context.Context, execution *sfnstore.Execution) (string, error) {
	if execution.StateMachineVersionArn != "" {
		version, err := e.store.GetStateMachineVersion(ctx, execution.StateMachineVersionArn)
		if err != nil {
			if !errors.Is(err, sfnstore.ErrStateMachineVersionNotFound) {
				return "", err
			}
		} else {
			return version.Definition, nil
		}
	}
	sm, err := e.store.GetStateMachine(ctx, execution.StateMachineArn)
	if err != nil {
		if !errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return "", err
		}
		return "", nil
	}
	return sm.Definition, nil
}

// parseDefinitionForExecution parses the definition the execution must
// run — the pinned version snapshot when the execution is version-
// qualified, the live state machine definition otherwise.
func (e *Executor) parseDefinitionForExecution(ctx context.Context, execution *sfnstore.Execution) (*sfnstore.StateMachineDefinition, map[string]sfnstore.State, error) {
	definition, err := e.definitionForExecution(ctx, execution)
	if err != nil {
		return nil, nil, err
	}
	if definition == "" {
		return nil, nil, fmt.Errorf("no definition found for state machine %s", execution.StateMachineArn)
	}
	return parseDefinitionJSON(definition)
}

// parseDefinitionJSON parses a state machine definition string into the
// definition struct and its typed states.
func parseDefinitionJSON(definition string) (*sfnstore.StateMachineDefinition, map[string]sfnstore.State, error) {
	var def sfnstore.StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return nil, nil, fmt.Errorf("failed to parse state machine definition: %w", err)
	}

	states, err := extractStatesFromDefinition(&def)
	if err != nil {
		return nil, nil, err
	}

	return &def, states, nil
}

func resolveOutputRaw(state sfnstore.State, stateMap map[string]interface{}) error {
	outputVal, hasOutput := stateMap["Output"]
	if !hasOutput {
		return nil
	}

	outputJSON, err := json.Marshal(outputVal)
	if err != nil {
		return err
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

// extractStatesFromDefinition decodes a definition's raw state values into
// their typed structs; every state-decoding caller shares this one loop.
func extractStatesFromDefinition(definition *sfnstore.StateMachineDefinition) (map[string]sfnstore.State, error) {
	states := make(map[string]sfnstore.State)
	for name, stateData := range definition.States {
		state, err := parseState(name, stateData)
		if err != nil {
			return nil, err
		}
		states[name] = state
	}
	return states, nil
}

// parseState decodes one state from its raw definition value: the typed
// struct is chosen by Type, the raw map is re-marshalled into it, and the
// Output literal is preserved as raw JSON.
func parseState(name string, stateData interface{}) (sfnstore.State, error) {
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
