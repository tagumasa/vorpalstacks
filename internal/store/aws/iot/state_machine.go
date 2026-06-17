package iot

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// inputReferencePattern captures the input name from AWS IoT Events
// $input.<inputName>.<field> condition references. The grammar allows
// letters, digits, underscores and dashes in input names, matching the
// AWS IoT Events CreateInput InputName rules.
var inputReferencePattern = regexp.MustCompile(`\$input\.([A-Za-z0-9_-]+)\.`)

// detectorState represents the current state of a single detector instance.
type detectorState struct {
	stateName   string
	lastUpdated time.Time
}

// DetectorStateMachine manages state transitions for detector instances.
// Each detector model can have multiple instances, keyed by the detector key
// (derived from the input message's key field).
type DetectorStateMachine struct {
	mu        sync.RWMutex
	detectors map[string]map[string]*detectorState
	models    map[string]*DetectorModel
	onAction  func(modelName, key, actionType string, payload map[string]interface{})
}

// InputMessage represents a single message for batch evaluation.
type InputMessage struct {
	InputName string
	Payload   map[string]interface{}
}

// NewDetectorStateMachine creates a new state machine with the given action callback.
// The onAction callback is invoked when a state transition triggers an action;
// it receives the model name, detector key, action type, and action parameters.
func NewDetectorStateMachine(onAction func(string, string, string, map[string]interface{})) *DetectorStateMachine {
	return &DetectorStateMachine{
		detectors: make(map[string]map[string]*detectorState),
		models:    make(map[string]*DetectorModel),
		onAction:  onAction,
	}
}

// LoadModel registers a detector model for event evaluation.
func (sm *DetectorStateMachine) LoadModel(dm *DetectorModel) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.models[dm.DetectorModelName] = dm
}

// UnloadModel removes a detector model and all its detector instances.
func (sm *DetectorStateMachine) UnloadModel(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.models, name)
	delete(sm.detectors, name)
}

// EvaluateEvent processes an input message against the named detector model.
// This method performs its own locking; callers must NOT hold sm.mu.
//
// The state mutation runs under sm.mu so detector instances are not
// corrupted by concurrent events, but action dispatch is deferred until
// after the lock is released. Holding the write lock while invoking
// external services (SQS/SNS/Lambda/iotTopicPublish through the eventbus)
// would serialise every detector evaluation behind the slowest downstream
// call.
//
// Evaluation order per AWS IoT Events semantics:
//  1. Evaluate onInput.transitionEvents of the current state (in order).
//     If a transition matches, execute its actions, then onExit of old state,
//     then onEnter of new state. Stop processing after the first match.
//  2. If no transition matched, evaluate onInput.events (non-transition).
//     If an event condition matches, execute its actions without state change.
//     Stop processing after the first match.
func (sm *DetectorStateMachine) EvaluateEvent(ctx context.Context, modelName, key string, payload map[string]interface{}) error {
	type pendingAction struct {
		actionType string
		params     map[string]interface{}
	}
	var pending []pendingAction
	collect := func(_, _, actionType string, params map[string]interface{}) {
		pending = append(pending, pendingAction{actionType: actionType, params: params})
	}

	// Phase 1: under lock, compute the transition and collect actions.
	if err := sm.evaluateUnderLock(modelName, key, payload, collect); err != nil {
		return err
	}

	// Phase 2: outside the lock, dispatch the collected actions so that
	// slow downstream calls do not block concurrent evaluations.
	for _, p := range pending {
		if sm.onAction != nil {
			sm.onAction(modelName, key, p.actionType, p.params)
		}
	}
	return nil
}

// evaluateUnderLock performs the state-machine transition under sm.mu and
// routes every action that should fire into the supplied collector. The
// collector lets EvaluateEvent buffer actions and dispatch them after the
// lock is released, avoiding head-of-line blocking on slow eventbus calls.
func (sm *DetectorStateMachine) evaluateUnderLock(modelName, key string, payload map[string]interface{}, collect func(string, string, string, map[string]interface{})) error {
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
		initialState := initialStateFromDefinition(model.DetectorModelDefinition)
		current = &detectorState{stateName: initialState, lastUpdated: time.Now().UTC()}
		instanceMap[key] = current
	}

	currentStateName := current.stateName

	def := model.DetectorModelDefinition
	states := getStatesFromDefinition(def)
	stateDef := findState(states, currentStateName)
	if stateDef == nil {
		return nil
	}

	// Step 1: Evaluate transitionEvents in onInput for state transitions.
	transitionEvents := getTransitionEvents(stateDef, "onInput")
	for _, te := range transitionEvents {
		if evaluateEventCondition(te, payload) {
			newState := transitionDestination(te)
			if newState == "" {
				continue
			}

			// Execute the transition event's own actions.
			dispatchActionMaps(te, collect, modelName, key)

			// Execute onExit actions of the current state.
			executeLifecycleActions(stateDef, "onExit", collect, modelName, key)

			current.stateName = newState
			current.lastUpdated = time.Now().UTC()

			// Execute onEnter actions of the new state.
			newStateDef := findState(states, newState)
			if newStateDef != nil {
				executeLifecycleActions(newStateDef, "onEnter", collect, modelName, key)
			}
			return nil
		}
	}

	// Step 2: No transition matched — evaluate onInput non-transition events.
	events := getContainerEvents(stateDef, "onInput")
	for _, event := range events {
		if evaluateEventCondition(event, payload) {
			dispatchActionMaps(event, collect, modelName, key)
			return nil // first matching event wins
		}
	}

	return nil
}

// GetState returns the current state of a detector instance.
func (sm *DetectorStateMachine) GetState(modelName, key string) (string, bool) {
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

// BatchEvaluate processes a batch of input messages against all loaded detector models.
// For each message, it derives the detector key and evaluates the event against
// every loaded model that declares the message's input. Returns error entries
// for any failed evaluations.
//
// A model "declares" an input when any of its state conditions reference it
// via the $input.<inputName>.<field> syntax. Models whose conditions never
// reference an input (e.g. "true") are evaluated against every message so
// that trivial transition guards keep firing.
func (sm *DetectorStateMachine) BatchEvaluate(ctx context.Context, messages []InputMessage) []map[string]interface{} {
	sm.mu.RLock()
	models := make([]string, 0, len(sm.models))
	consumedInputs := make(map[string]map[string]bool, len(sm.models))
	for modelName, dm := range sm.models {
		models = append(models, modelName)
		consumedInputs[modelName] = extractDeclaredInputs(dm.DetectorModelDefinition)
	}
	sm.mu.RUnlock()

	var errs []map[string]interface{}
	for _, msg := range messages {
		if msg.InputName == "" || msg.Payload == nil {
			continue
		}
		for _, modelName := range models {
			inputs := consumedInputs[modelName]
			if len(inputs) > 0 && !inputs[msg.InputName] {
				continue
			}
			key := sm.DeriveDetectorKey(modelName, msg.InputName, msg.Payload)
			if err := sm.EvaluateEvent(ctx, modelName, key, msg.Payload); err != nil {
				errs = append(errs, map[string]interface{}{
					"errorCode":    "InvalidRequestException",
					"errorMessage": err.Error(),
				})
			}
		}
	}
	return errs
}

// extractDeclaredInputs walks a detector model definition and returns the
// set of input names referenced via the $input.<inputName>.<field> syntax.
// Returns nil when the definition contains no $input references at all,
// which lets BatchEvaluate treat the model as a catch-all (matching the
// behaviour of trivial "true" conditions that do not name an input).
func extractDeclaredInputs(def map[string]interface{}) map[string]bool {
	var inputs map[string]bool
	walkStrings(def, func(s string) {
		for _, match := range inputReferencePattern.FindAllStringSubmatch(s, -1) {
			if inputs == nil {
				inputs = make(map[string]bool)
			}
			inputs[match[1]] = true
		}
	})
	return inputs
}

// walkStrings traverses the dynamic definition tree and invokes fn for
// every string value encountered. Conditions and action parameters in
// the AWS IoT Events definition are arbitrary JSON, so the walker keeps
// the traversal generic instead of assuming a fixed schema.
func walkStrings(v interface{}, fn func(string)) {
	switch t := v.(type) {
	case string:
		fn(t)
	case map[string]interface{}:
		for _, child := range t {
			walkStrings(child, fn)
		}
	case []interface{}:
		for _, child := range t {
			walkStrings(child, fn)
		}
	}
}

// DeriveDetectorKey resolves the detector instance key for the given model.
// If the model has a Key configured, that attribute is extracted from the
// payload. Otherwise the inputName is used as the key so that each input
// gets its own detector instance.
func (sm *DetectorStateMachine) DeriveDetectorKey(modelName, inputName string, payload map[string]interface{}) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	dm := sm.models[modelName]
	if dm != nil && dm.Key != "" {
		if v, ok := payload[dm.Key].(string); ok && v != "" {
			return v
		}
	}
	return inputName
}

// ExtractPayload converts the raw payload field from a BatchPutMessage entry
// into a map. AWS sends payload as a base64-encoded blob; the decoded bytes
// are expected to be a JSON object.
func ExtractPayload(raw interface{}) map[string]interface{} {
	switch p := raw.(type) {
	case map[string]interface{}:
		return p
	case string:
		if p == "" {
			return nil
		}
		decoded, err := base64.StdEncoding.DecodeString(p)
		if err != nil {
			decoded = []byte(p)
		}
		var m map[string]interface{}
		if err := json.Unmarshal(decoded, &m); err != nil {
			slog.Warn("iot events payload parse failed", "error", err)
			return nil
		}
		return m
	}
	return nil
}

// --- Definition parsing helpers ---

func initialStateFromDefinition(def map[string]interface{}) string {
	if def == nil {
		return ""
	}
	if name, ok := def["initialStateName"].(string); ok && name != "" {
		return name
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

// getTransitionEvents returns transition events from a lifecycle container
// (onInput, onEnter, or onExit) of a state definition.
// AWS wire format places transitionEvents as a sibling of events inside
// each lifecycle container:
//
//	"onInput": {
//	  "events": [...],              // non-transition events
//	  "transitionEvents": [...]     // transition events (with nextState)
//	}
func getTransitionEvents(state map[string]interface{}, containerKey string) []map[string]interface{} {
	container, ok := state[containerKey].(map[string]interface{})
	if !ok {
		return nil
	}
	raw, ok := container["transitionEvents"]
	if !ok {
		return nil
	}
	return toMapSlice(raw)
}

// getContainerEvents returns non-transition events from a lifecycle container.
// These events fire actions when their condition matches, without causing
// state transitions.
func getContainerEvents(state map[string]interface{}, containerKey string) []map[string]interface{} {
	container, ok := state[containerKey].(map[string]interface{})
	if !ok {
		return nil
	}
	raw, ok := container["events"]
	if !ok {
		return nil
	}
	return toMapSlice(raw)
}

// toMapSlice converts a raw interface value (expected to be []interface{} of
// maps) into a typed slice of map[string]interface{}.
func toMapSlice(raw interface{}) []map[string]interface{} {
	switch v := raw.(type) {
	case []interface{}:
		result := make([]map[string]interface{}, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
		return result
	case []map[string]interface{}:
		return v
	}
	return nil
}

// evaluateEventCondition checks whether an event's condition matches the
// given payload. A missing or empty condition is treated as "true".
func evaluateEventCondition(event map[string]interface{}, payload map[string]interface{}) bool {
	condRaw, ok := event["condition"]
	if !ok {
		return true
	}
	cond, ok := condRaw.(string)
	if !ok || cond == "" {
		return true
	}
	if cond == "true" {
		return true
	}
	return evaluateSimpleCondition(cond, payload)
}

// evaluateSimpleCondition evaluates a condition expression against a payload.
//
// Supports:
//   - Simple field references: "temperature > 50"
//   - $input references: "$input.TempInput.temperature > 50"
//   - String comparisons with quotes: "$input.S.status == 'critical'"
//   - Operators: >, >=, <, <=, ==, =, !=
//
// For $input.<InputName>.<field> syntax, the InputName segment is skipped
// because the payload already contains the resolved fields from the input.
func evaluateSimpleCondition(condition string, payload map[string]interface{}) bool {
	field, op, value, ok := parseCondition(condition)
	if !ok {
		return false
	}

	fieldVal, exists := resolveField(field, payload)
	if !exists {
		return false
	}

	return compareValues(fieldVal, op, value)
}

// parseCondition splits a condition string into field, operator, and value.
// Tries longest operators first to avoid ambiguity (>= before >).
func parseCondition(condition string) (field, op, value string, ok bool) {
	condition = strings.TrimSpace(condition)
	for _, candidate := range []string{">=", "<=", "!=", "==", ">", "<", "="} {
		idx := strings.Index(condition, candidate)
		if idx < 0 {
			continue
		}
		f := strings.TrimSpace(condition[:idx])
		v := strings.TrimSpace(condition[idx+len(candidate):])
		if f != "" && v != "" {
			return f, candidate, v, true
		}
	}
	return "", "", "", false
}

// resolveField extracts a field value from the payload.
// For $input.<InputName>.<field> references, the $input. prefix and
// InputName segment are stripped, leaving just the field name to look
// up in the payload (which already contains the resolved input fields).
func resolveField(field string, payload map[string]interface{}) (interface{}, bool) {
	field = strings.TrimPrefix(field, "$input.")
	// Skip the InputName segment (e.g. "TempInput" in "TempInput.temperature").
	parts := strings.SplitN(field, ".", 2)
	if len(parts) == 2 {
		field = parts[1]
	}
	val, ok := payload[field]
	return val, ok
}

// compareValues compares a resolved field value against a string threshold
// using the given operator. Supports numeric (float64) and string comparisons.
func compareValues(fieldVal interface{}, op, value string) bool {
	target := trimQuotes(value)

	switch v := fieldVal.(type) {
	case float64:
		threshold, perr := strconv.ParseFloat(target, 64)
		if perr != nil {
			return false
		}
		switch op {
		case ">":
			return v > threshold
		case ">=":
			return v >= threshold
		case "<":
			return v < threshold
		case "<=":
			return v <= threshold
		case "==", "=":
			return v == threshold
		case "!=":
			return v != threshold
		}
	case bool:
		switch op {
		case "==", "=":
			switch target {
			case "true":
				return v
			case "false":
				return !v
			}
		case "!=":
			switch target {
			case "true":
				return !v
			case "false":
				return v
			}
		}
	case string:
		switch op {
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

// executeLifecycleActions extracts and dispatches actions from a lifecycle
// container (onEnter/onExit). The container format is:
//
//	"onEnter": {
//	  "events": [{"eventName":"...", "condition":"...", "actions":[...]}]
//	}
//
// All actions from all events are dispatched (lifecycle events typically
// have condition "true" or no condition).
func executeLifecycleActions(state map[string]interface{}, containerKey string, onAction func(string, string, string, map[string]interface{}), modelName, key string) {
	container := state[containerKey]
	if container == nil {
		return
	}

	var events []map[string]interface{}
	switch v := container.(type) {
	case map[string]interface{}:
		raw := v["events"]
		if raw == nil {
			return
		}
		events = toMapSlice(raw)
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				events = append(events, m)
			}
		}
	case []map[string]interface{}:
		events = v
	default:
		return
	}

	for _, event := range events {
		dispatchActionMaps(event, onAction, modelName, key)
	}
}

// dispatchActionMaps extracts action maps from an event or transition and
// invokes the onAction callback for each. Detector model actions use the
// format where the action type is the key:
//
//	{"sqs": {"queueUrl": "...", "payload": {...}}}
//	{"iotTopicPublish": {"mqttTopic": "topic/name"}}
//	{"lambda": {"functionArn": "..."}}
//
// This also handles the IoT rule action format with an explicit "action" key:
//
//	{"action": "sqs", "parameters": {...}}
func dispatchActionMaps(event map[string]interface{}, onAction func(string, string, string, map[string]interface{}), modelName, key string) {
	actionsRaw := event["actions"]
	if actionsRaw == nil {
		return
	}
	actionMaps := toMapSlice(actionsRaw)

	for _, action := range actionMaps {
		// IoT rule action format: explicit "action" field.
		if actionName, ok := action["action"].(string); ok && actionName != "" {
			params, _ := action["parameters"].(map[string]interface{})
			if onAction != nil {
				onAction(modelName, key, actionName, params)
			}
			continue
		}
		// Detector model action format: action type is the key itself.
		for actionName, cfg := range action {
			params, _ := cfg.(map[string]interface{})
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
