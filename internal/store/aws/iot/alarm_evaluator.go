// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package iot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// inputPropertyPattern captures the inputName and attribute path from the
// AWS IoT Events $input.<inputName>.<attributePath> reference format used
// in SimpleRule.inputProperty and SimpleRule.threshold (H-SM2).
var inputPropertyPattern = regexp.MustCompile(`^\$input\.([A-Za-z0-9_-]+)\.(.+)$`)

// parsedInputProperty holds the parsed components of a $input reference.
type parsedInputProperty struct {
	inputName     string
	attributePath string
}

// parseInputProperty parses a $input.<inputName>.<attributePath> reference.
// Returns nil if the string is not a valid $input reference.
func parseInputProperty(s string) *parsedInputProperty {
	m := inputPropertyPattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return nil
	}
	return &parsedInputProperty{
		inputName:     m[1],
		attributePath: m[2],
	}
}

// simpleRuleConfig holds the pre-parsed SimpleRule extracted from an
// AlarmModelDefinition. Pre-parsing at load time avoids re-parsing on
// every message evaluation.
type simpleRuleConfig struct {
	comparisonOperator string
	inputProperty      *parsedInputProperty
	thresholdValue     float64              // valid when thresholdIsRef is false
	thresholdRef       *parsedInputProperty // non-nil when threshold is a $input reference
	thresholdIsRef     bool
}

// extractSimpleRule walks the AlarmModelDefinition JSON to find the
// SimpleRule. Returns nil if the definition has no alarmRule.simpleRule.
//
// AWS structure (verified via API docs):
//
//	{
//	  "alarmRule": {
//	    "simpleRule": {
//	      "comparisonOperator": "GREATER",
//	      "inputProperty": "$input.tempInput.temperature",
//	      "threshold": "30"
//	    }
//	  }
//	}
func extractSimpleRule(def map[string]interface{}) *simpleRuleConfig {
	if def == nil {
		return nil
	}
	alarmRule, ok := def["alarmRule"].(map[string]interface{})
	if !ok {
		return nil
	}
	sr, ok := alarmRule["simpleRule"].(map[string]interface{})
	if !ok {
		return nil
	}
	op, _ := sr["comparisonOperator"].(string)
	if op == "" {
		return nil
	}
	inputPropStr, _ := sr["inputProperty"].(string)
	inputProp := parseInputProperty(inputPropStr)
	if inputProp == nil {
		return nil
	}
	thresholdStr, _ := sr["threshold"].(string)
	if thresholdStr == "" {
		return nil
	}
	cfg := &simpleRuleConfig{
		comparisonOperator: op,
		inputProperty:      inputProp,
	}
	if ref := parseInputProperty(thresholdStr); ref != nil {
		cfg.thresholdRef = ref
		cfg.thresholdIsRef = true
	} else if v, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
		cfg.thresholdValue = v
	} else {
		return nil
	}
	return cfg
}

// loadedAlarmModel holds an alarm model name, its key attribute, and the
// pre-parsed SimpleRule used for evaluation.
type loadedAlarmModel struct {
	name string
	key  string
	rule *simpleRuleConfig
}

// LoadAlarmModel registers an alarm model for message evaluation. If the
// model has no SimpleRule, it is silently skipped (non-simple alarm rules
// are not supported). Re-loading with the same name replaces the previous
// definition (H-SM2).
func (sm *AlarmStateMachine) LoadAlarmModel(am *AlarmModel) {
	if am == nil {
		return
	}
	rule := extractSimpleRule(am.AlarmModelDefinition)
	if rule == nil {
		sm.mu.Lock()
		// Even without a rule, track the model name so DescribeAlarm/ListAlarms
		// can distinguish "model exists but no instances" from "model unknown".
		if _, exists := sm.modelDefs[am.AlarmModelName]; !exists {
			sm.modelDefs[am.AlarmModelName] = true
		}
		sm.mu.Unlock()
		return
	}
	sm.mu.Lock()
	sm.modelDefs[am.AlarmModelName] = true
	sm.alarmRules[am.AlarmModelName] = &loadedAlarmModel{
		name: am.AlarmModelName,
		key:  am.Key,
		rule: rule,
	}
	sm.mu.Unlock()
}

// UnloadAlarmModel removes an alarm model from evaluation. Existing alarm
// instances are preserved so their state survives model deletion (H-SM2).
func (sm *AlarmStateMachine) UnloadAlarmModel(name string) {
	sm.mu.Lock()
	delete(sm.modelDefs, name)
	delete(sm.alarmRules, name)
	sm.mu.Unlock()
}

// EvaluateMessages processes a batch of input messages against all loaded
// alarm models. For each message whose inputName matches an alarm model's
// inputProperty reference, the SimpleRule is evaluated and the alarm state
// is transitioned accordingly:
//   - Condition true  → NORMAL transitions to ACTIVE
//   - Condition false → ACTIVE transitions to NORMAL
//
// Other state transitions (ACKNOWLEDGED, SNOOZE, DISABLED) are driven by
// explicit API calls and are not affected by evaluation.
func (sm *AlarmStateMachine) EvaluateMessages(messages []InputMessage) {
	sm.mu.RLock()
	rules := make([]*loadedAlarmModel, 0, len(sm.alarmRules))
	for _, r := range sm.alarmRules {
		rules = append(rules, r)
	}
	sm.mu.RUnlock()

	for _, msg := range messages {
		if msg.InputName == "" || msg.Payload == nil {
			continue
		}
		for _, model := range rules {
			if model.rule.inputProperty.inputName != msg.InputName {
				continue
			}
			inputVal, ok := extractNumericValue(msg.Payload, model.rule.inputProperty.attributePath)
			if !ok {
				continue
			}
			threshold, ok := sm.resolveThreshold(model.rule, messages, msg)
			if !ok {
				continue
			}
			triggered := evaluateComparison(model.rule.comparisonOperator, inputVal, threshold)
			keyValue := deriveAlarmKey(model, msg)
			if triggered {
				sm.transitionToActive(model.name, keyValue)
			} else {
				sm.transitionToNormal(model.name, keyValue)
			}
		}
	}
}

// resolveThreshold resolves the threshold value. If the threshold is a
// fixed numeric, it returns the stored value. If it is a $input reference,
// it looks up the corresponding message payload.
func (sm *AlarmStateMachine) resolveThreshold(rule *simpleRuleConfig, messages []InputMessage, currentMsg InputMessage) (float64, bool) {
	if !rule.thresholdIsRef {
		return rule.thresholdValue, true
	}
	// First try the current message (most common: same input).
	if val, ok := extractNumericValue(currentMsg.Payload, rule.thresholdRef.attributePath); ok {
		return val, true
	}
	// Fall back to searching other messages for the referenced input.
	for _, msg := range messages {
		if msg.InputName == rule.thresholdRef.inputName {
			if val, ok := extractNumericValue(msg.Payload, rule.thresholdRef.attributePath); ok {
				return val, true
			}
		}
	}
	return 0, false
}

// transitionToActive moves an alarm to ACTIVE state. If the alarm instance
// does not yet exist (first message for this key), it is created in ACTIVE
// state — AWS IoT Events auto-creates alarm instances on first message
// arrival via BatchPutMessage. Disabled alarms are skipped.
func (sm *AlarmStateMachine) transitionToActive(modelName, keyValue string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.alarms[modelName] == nil {
		sm.alarms[modelName] = make(map[string]*alarmInstanceState)
	}
	now := time.Now().UTC()
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		// Auto-create: first message for this key triggers alarm creation.
		alarm = &alarmInstanceState{
			stateName:    AlarmStateActive,
			creationTime: now,
			lastUpdate:   now,
		}
		sm.alarms[modelName][keyValue] = alarm
		sm.persist(modelName, keyValue, alarm)
		return
	}
	if alarm.stateName == AlarmStateDisabled {
		return
	}
	if alarm.stateName == AlarmStateActive {
		return
	}
	alarm.stateName = AlarmStateActive
	alarm.lastUpdate = now
	sm.persist(modelName, keyValue, alarm)
}

// transitionToNormal moves an alarm to NORMAL state. If the alarm instance
// does not yet exist, it is created in NORMAL state — AWS IoT Events
// auto-creates alarm instances on first message arrival via BatchPutMessage.
func (sm *AlarmStateMachine) transitionToNormal(modelName, keyValue string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.alarms[modelName] == nil {
		sm.alarms[modelName] = make(map[string]*alarmInstanceState)
	}
	now := time.Now().UTC()
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		// Auto-create: first message for this key triggers alarm creation.
		alarm = &alarmInstanceState{
			stateName:    AlarmStateNormal,
			creationTime: now,
			lastUpdate:   now,
		}
		sm.alarms[modelName][keyValue] = alarm
		sm.persist(modelName, keyValue, alarm)
		return
	}
	if alarm.stateName != AlarmStateActive {
		return
	}
	alarm.stateName = AlarmStateNormal
	alarm.lastUpdate = now
	sm.persist(modelName, keyValue, alarm)
}

// deriveAlarmKey resolves the alarm instance key for the given model and
// message. If the model has no key attribute, an empty string is used
// (single-instance alarm). Otherwise, the key attribute value is extracted
// from the payload.
func deriveAlarmKey(model *loadedAlarmModel, msg InputMessage) string {
	if model.key == "" {
		return ""
	}
	val, ok := msg.Payload[model.key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

// extractNumericValue traverses a nested JSON path (e.g. "data.temperature")
// in the payload and converts the result to float64 for numeric comparison.
func extractNumericValue(payload map[string]interface{}, path string) (float64, bool) {
	parts := strings.Split(path, ".")
	var current interface{} = payload
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return 0, false
		}
		current, ok = m[part]
		if !ok {
			return 0, false
		}
	}
	return toFloat64(current)
}

// toFloat64 converts common JSON numeric types to float64.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(n, 64)
		return f, err == nil
	}
	return 0, false
}

// evaluateComparison applies the AWS SimpleRule comparison operator to
// the input value and threshold. Both values must be numeric.
func evaluateComparison(op string, inputVal, threshold float64) bool {
	switch op {
	case "GREATER":
		return inputVal > threshold
	case "GREATER_OR_EQUAL":
		return inputVal >= threshold
	case "LESS":
		return inputVal < threshold
	case "LESS_OR_EQUAL":
		return inputVal <= threshold
	case "EQUAL":
		return inputVal == threshold
	case "NOT_EQUAL":
		return inputVal != threshold
	}
	return false
}
