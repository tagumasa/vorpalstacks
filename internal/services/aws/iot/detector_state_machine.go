package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// detectorState represents the current state of a single detector instance.
type detectorState struct {
	stateName   string
	lastUpdated time.Time
}

// detectorStateMachine manages state transitions for detector instances.
// Each detector model can have multiple instances, keyed by the detector key
// (derived from the input message's key field).
type detectorStateMachine struct {
	mu        sync.RWMutex
	detectors map[string]map[string]*detectorState
	models    map[string]*iotstore.DetectorModel
	onAction  func(modelName, key, actionType string, payload map[string]interface{})
}

func newDetectorStateMachine(onAction func(string, string, string, map[string]interface{})) *detectorStateMachine {
	return &detectorStateMachine{
		detectors: make(map[string]map[string]*detectorState),
		models:    make(map[string]*iotstore.DetectorModel),
		onAction:  onAction,
	}
}

func (sm *detectorStateMachine) LoadModel(dm *iotstore.DetectorModel) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.models[dm.DetectorModelName] = dm
}

func (sm *detectorStateMachine) UnloadModel(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.models, name)
	delete(sm.detectors, name)
}

// EvaluateEvent processes an input message against the named detector model.
// This method performs its own locking; callers must NOT hold sm.mu.
func (sm *detectorStateMachine) EvaluateEvent(ctx context.Context, modelName, key string, payload map[string]interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	model, ok := sm.models[modelName]
	if !ok {
		return fmt.Errorf("detector model %q not found", modelName)
	}

	instanceMap := sm.detectors[modelName]
	if instanceMap == nil {
		instanceMap = make(map[string]*detectorState)
		sm.detectors[modelName] = instanceMap
	}

	current, exists := instanceMap[key]
	if !exists {
		// Determine initial state from model definition
		initialState := initialStateFromDefinition(model.DetectorModelDefinition)
		current = &detectorState{stateName: initialState, lastUpdated: time.Now().UTC()}
		instanceMap[key] = current
	}

	currentStateName := current.stateName

	// Find transition matching the event in current state.
	// Read model definition under lock to avoid racing with LoadModel/UnloadModel.
	def := model.DetectorModelDefinition
	states := getStatesFromDefinition(def)

	for _, state := range states {
		if stateName(state) != currentStateName {
			continue
		}
		transitions := getTransitions(state)
		for _, transition := range transitions {
			if evaluateTransitionCondition(transition, payload) {
				newState := transitionDestination(transition)
				current.stateName = newState
				current.lastUpdated = time.Now().UTC()

				// Execute onExit actions for old state
				executeActions(state, "onExit", sm.onAction, modelName, key)
				// Execute onEnter actions for new state
				newStateDef := findState(states, newState)
				if newStateDef != nil {
					executeActions(newStateDef, "onEnter", sm.onAction, modelName, key)
				}
				return nil
			}
		}
	}

	// No transition matched — stay in current state
	return nil
}

// GetState returns the current state of a detector instance.
func (sm *detectorStateMachine) GetState(modelName, key string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	instanceMap, ok := sm.detectors[modelName]
	if !ok {
		return "", false
	}
	ds, ok := instanceMap[key]
	if !ok {
		return "", false
	}
	return ds.stateName, true
}

// --- Definition parsing helpers ---

func initialStateFromDefinition(def map[string]interface{}) string {
	if def == nil {
		return ""
	}
	states := getStatesFromDefinition(def)
	if len(states) > 0 {
		return stateName(states[0])
	}
	return ""
}

func getStatesFromDefinition(def map[string]interface{}) []map[string]interface{} {
	if def == nil {
		return nil
	}
	statesRaw, ok := def["states"]
	if !ok {
		return nil
	}
	switch s := statesRaw.(type) {
	case []interface{}:
		result := make([]map[string]interface{}, 0, len(s))
		for _, item := range s {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
		return result
	case []map[string]interface{}:
		return s
	}
	return nil
}

func stateName(state map[string]interface{}) string {
	if v, ok := state["stateName"].(string); ok {
		return v
	}
	return ""
}

func getTransitions(state map[string]interface{}) []map[string]interface{} {
	onInputRaw, ok := state["onInput"]
	if !ok {
		return nil
	}
	onInput, ok := onInputRaw.(map[string]interface{})
	if !ok {
		return nil
	}
	eventsRaw, ok := onInput["events"]
	if !ok {
		return nil
	}
	switch e := eventsRaw.(type) {
	case []interface{}:
		// Each event can have transitions
		var transitions []map[string]interface{}
		for _, item := range e {
			if event, ok := item.(map[string]interface{}); ok {
				transRaw, ok := event["transitions"]
				if ok {
					switch t := transRaw.(type) {
					case []interface{}:
						for _, tr := range t {
							if m, ok := tr.(map[string]interface{}); ok {
								transitions = append(transitions, m)
							}
						}
					case []map[string]interface{}:
						transitions = append(transitions, t...)
					}
				}
			}
		}
		return transitions
	}
	return nil
}

func evaluateTransitionCondition(transition map[string]interface{}, payload map[string]interface{}) bool {
	// A transition with no condition always matches (true condition).
	condRaw, ok := transition["condition"]
	if !ok {
		return true
	}
	cond, ok := condRaw.(string)
	if !ok || cond == "" {
		return true
	}
	// Basic condition evaluation: check if the condition string is "true"
	if cond == "true" {
		return true
	}
	// For simple equality checks like "payload.temp > 50",
	// we do a basic string match against payload values.
	// Full expression evaluation is out of scope for this basic implementation.
	return evaluateSimpleCondition(cond, payload)
}

func evaluateSimpleCondition(condition string, payload map[string]interface{}) bool {
	// Parse simple conditions of the form: field op value
	// e.g., "temp > 50" or "status == 'active'"
	// This is a minimal evaluator — full SQL-like evaluation is deferred.
	type condParts struct {
		field string
		op    string
		value string
	}

	parts := &condParts{}
	n, err := fmt.Sscanf(condition, "%s %s %s", &parts.field, &parts.op, &parts.value)
	if err != nil || n != 3 {
		return false
	}

	fieldVal, ok := payload[parts.field]
	if !ok {
		return false
	}

	switch v := fieldVal.(type) {
	case float64:
		var threshold float64
		fmt.Sscanf(parts.value, "%f", &threshold)
		switch parts.op {
		case ">", ">=":
			if parts.op == ">" {
				return v > threshold
			}
			return v >= threshold
		case "<", "<=":
			if parts.op == "<" {
				return v < threshold
			}
			return v <= threshold
		case "==", "=":
			return v == threshold
		}
	case bool:
		switch parts.op {
		case "==", "=":
			switch parts.value {
			case "true":
				return v
			case "false":
				return !v
			}
		case "!=":
			switch parts.value {
			case "true":
				return !v
			case "false":
				return v
			}
		}
	case string:
		target := trimQuotes(parts.value)
		switch parts.op {
		case "==", "=":
			return v == target
		case "!=":
			return v != target
		}
	}
	return false
}

func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func transitionDestination(transition map[string]interface{}) string {
	if v, ok := transition["nextState"].(string); ok {
		return v
	}
	return ""
}

func findState(states []map[string]interface{}, name string) map[string]interface{} {
	for _, s := range states {
		if stateName(s) == name {
			return s
		}
	}
	return nil
}

func executeActions(state map[string]interface{}, actionType string, onAction func(string, string, string, map[string]interface{}), modelName, key string) {
	actionsRaw, ok := state[actionType]
	if !ok {
		return
	}
	switch actions := actionsRaw.(type) {
	case []interface{}:
		for _, a := range actions {
			if action, ok := a.(map[string]interface{}); ok {
				actionName, _ := action["action"].(string)
				params, _ := action["parameters"].(map[string]interface{})
				if onAction != nil && actionName != "" {
					onAction(modelName, key, actionName, params)
				}
			}
		}
	case []map[string]interface{}:
		for _, action := range actions {
			actionName, _ := action["action"].(string)
			params, _ := action["parameters"].(map[string]interface{})
			if onAction != nil && actionName != "" {
				onAction(modelName, key, actionName, params)
			}
		}
	}
}

// detectorModelDefinitionJSON is a helper for tests that builds a minimal
// detector model definition from a JSON string.
func detectorModelDefinitionJSON(jsonStr string) map[string]interface{} {
	var def map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &def); err != nil {
		return nil
	}
	return def
}
