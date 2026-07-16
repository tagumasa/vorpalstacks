package iot

import (
	"context"
	"sync"
	"testing"
)

func TestStateMachine_TwoStateTransition(t *testing.T) {
	var mu sync.Mutex
	var actions []string
	onAction := func(modelName, key, actionType string, _ map[string]interface{}) {
		mu.Lock()
		actions = append(actions, actionType)
		mu.Unlock()
	}

	sm := NewDetectorStateMachine(onAction)

	// AWS wire format: transitionEvents is a sibling of events inside onInput.
	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"transitionEvents": [
						{
							"condition": "true",
							"nextState": "Active"
						}
					]
				},
				"onExit": {
					"events": [
						{
							"actions": [{"iotTopicPublish": {"mqttTopic": "exit/idle"}}]
						}
					]
				}
			},
			{
				"stateName": "Active",
				"onEnter": {
					"events": [
						{
							"actions": [{"sns": {"topicArn": "arn:aws:sns:us-east-1:000000000000:enter-active"}}]
						}
					]
				}
			}
		]
	}`)

	dm := &DetectorModel{
		DetectorModelName:       "test-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	}
	sm.LoadModel(dm)

	_ = sm.EvaluateEvent(context.Background(), "test-model", "key1", map[string]interface{}{})

	state, ok := sm.GetState("test-model", "key1")
	if !ok {
		t.Fatal("expected state after first EvaluateEvent")
	}
	if state != "Active" {
		t.Fatalf("expected state 'Active' after transition, got %q", state)
	}

	// Second evaluation — no transition matches (Active has no transitionEvents).
	_ = sm.EvaluateEvent(context.Background(), "test-model", "key1", map[string]interface{}{"temp": 60.0})
	_, _ = sm.GetState("test-model", "key1")

	mu.Lock()
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions (onExit iotTopicPublish + onEnter sns), got %d: %v", len(actions), actions)
	}
	if actions[0] != "iotTopicPublish" {
		t.Fatalf("expected first action 'iotTopicPublish' (onExit), got %q", actions[0])
	}
	if actions[1] != "sns" {
		t.Fatalf("expected second action 'sns' (onEnter), got %q", actions[1])
	}
	mu.Unlock()
}

func TestStateMachine_ConditionMatch(t *testing.T) {
	sm := NewDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Normal",
				"onInput": {
					"transitionEvents": [
						{
							"condition": "$input.TempInput.temperature > 50",
							"nextState": "HighTemp"
						}
					]
				}
			},
			{
				"stateName": "HighTemp"
			}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "temp-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	// temperature=30 does NOT match > 50, stays Normal.
	_ = sm.EvaluateEvent(context.Background(), "temp-model", "sensor1", map[string]interface{}{"temperature": 30.0})
	state, _ := sm.GetState("temp-model", "sensor1")
	if state != "Normal" {
		t.Fatalf("expected 'Normal' for temp=30, got %q", state)
	}

	// temperature=75 matches > 50, transitions to HighTemp.
	_ = sm.EvaluateEvent(context.Background(), "temp-model", "sensor1", map[string]interface{}{"temperature": 75.0})
	state, _ = sm.GetState("temp-model", "sensor1")
	if state != "HighTemp" {
		t.Fatalf("expected 'HighTemp' for temp=75, got %q", state)
	}
}

func TestStateMachine_InputCondition_Numeric(t *testing.T) {
	sm := NewDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"transitionEvents": [
						{
							"condition": "$input.Sensor.pressure >= 100",
							"nextState": "Alert"
						}
					]
				}
			},
			{
				"stateName": "Alert"
			}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "pressure-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	// pressure=99 does NOT match >= 100.
	_ = sm.EvaluateEvent(context.Background(), "pressure-model", "p1", map[string]interface{}{"pressure": 99.0})
	state, _ := sm.GetState("pressure-model", "p1")
	if state != "Idle" {
		t.Fatalf("expected 'Idle' for pressure=99, got %q", state)
	}

	// pressure=100 matches >= 100.
	_ = sm.EvaluateEvent(context.Background(), "pressure-model", "p1", map[string]interface{}{"pressure": 100.0})
	state, _ = sm.GetState("pressure-model", "p1")
	if state != "Alert" {
		t.Fatalf("expected 'Alert' for pressure=100, got %q", state)
	}
}

func TestStateMachine_InputCondition_String(t *testing.T) {
	sm := NewDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Normal",
				"onInput": {
					"transitionEvents": [
						{
							"condition": "$input.S.status == 'critical'",
							"nextState": "Alarm"
						}
					]
				}
			},
			{
				"stateName": "Alarm"
			}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "status-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	// status=normal does NOT match 'critical'.
	_ = sm.EvaluateEvent(context.Background(), "status-model", "d1", map[string]interface{}{"status": "normal"})
	state, _ := sm.GetState("status-model", "d1")
	if state != "Normal" {
		t.Fatalf("expected 'Normal' for status=normal, got %q", state)
	}

	// status=critical matches 'critical'.
	_ = sm.EvaluateEvent(context.Background(), "status-model", "d1", map[string]interface{}{"status": "critical"})
	state, _ = sm.GetState("status-model", "d1")
	if state != "Alarm" {
		t.Fatalf("expected 'Alarm' for status=critical, got %q", state)
	}
}

func TestStateMachine_OnInputNonTransitionActions(t *testing.T) {
	var mu sync.Mutex
	var actions []string
	onAction := func(modelName, key, actionType string, _ map[string]interface{}) {
		mu.Lock()
		actions = append(actions, actionType)
		mu.Unlock()
	}

	sm := NewDetectorStateMachine(onAction)

	// Non-transition event in onInput.events fires actions without state change.
	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"events": [
						{
							"condition": "$input.Sensor.value > 30",
							"actions": [
								{"iotTopicPublish": {"mqttTopic": "data/alert"}}
							]
						}
					]
				}
			}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "notify-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	_ = sm.EvaluateEvent(context.Background(), "notify-model", "s1", map[string]interface{}{"value": 50.0})

	// State should NOT have changed.
	state, _ := sm.GetState("notify-model", "s1")
	if state != "Idle" {
		t.Fatalf("expected state still 'Idle', got %q", state)
	}

	mu.Lock()
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d: %v", len(actions), actions)
	}
	if actions[0] != "iotTopicPublish" {
		t.Fatalf("expected 'iotTopicPublish', got %q", actions[0])
	}
	mu.Unlock()
}

func TestStateMachine_TransitionFalse_NoStateChange(t *testing.T) {
	sm := NewDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"transitionEvents": [
						{
							"condition": "$input.Sensor.value > 100",
							"nextState": "Alarm"
						}
					]
				}
			},
			{
				"stateName": "Alarm"
			}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "threshold-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	// value=50 does NOT match > 100 — no transition.
	_ = sm.EvaluateEvent(context.Background(), "threshold-model", "x1", map[string]interface{}{"value": 50.0})
	state, _ := sm.GetState("threshold-model", "x1")
	if state != "Idle" {
		t.Fatalf("expected 'Idle' (no transition), got %q", state)
	}
}

func TestStateMachine_MultipleInstances(t *testing.T) {
	sm := NewDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Off",
				"onInput": {
					"transitionEvents": [{"condition": "trigger == true", "nextState": "On"}]
				}
			},
			{"stateName": "On"}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "switch-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	_ = sm.EvaluateEvent(context.Background(), "switch-model", "switch-A", map[string]interface{}{"trigger": true})
	_ = sm.EvaluateEvent(context.Background(), "switch-model", "switch-B", map[string]interface{}{"trigger": false})

	stateA, _ := sm.GetState("switch-model", "switch-A")
	if stateA != "On" {
		t.Fatalf("expected switch-A 'On', got %q", stateA)
	}

	stateB, _ := sm.GetState("switch-model", "switch-B")
	if stateB != "Off" {
		t.Fatalf("expected switch-B still 'Off', got %q", stateB)
	}
}

func TestStateMachine_UnloadModel(t *testing.T) {
	sm := NewDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "A",
				"onInput": {
					"transitionEvents": [{"condition": "true", "nextState": "B"}]
				}
			},
			{"stateName": "B"}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "transient",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	_ = sm.EvaluateEvent(context.Background(), "transient", "k1", map[string]interface{}{})
	state, ok := sm.GetState("transient", "k1")
	if !ok || state != "B" {
		t.Fatalf("expected 'B', got %q (ok=%v)", state, ok)
	}

	sm.UnloadModel("transient")

	_, ok = sm.GetState("transient", "k1")
	if ok {
		t.Fatal("expected state to be gone after UnloadModel")
	}
}

func TestStateMachine_EvaluateNonexistentModel(t *testing.T) {
	sm := NewDetectorStateMachine(nil)
	err := sm.EvaluateEvent(context.Background(), "no-such-model", "k1", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for nonexistent model")
	}
}

func TestStateMachine_DeriveDetectorKey(t *testing.T) {
	sm := NewDetectorStateMachine(nil)
	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "keyed-model",
		DetectorModelDefinition: map[string]interface{}{},
		Key:                     "deviceId",
		Status:                  "ACTIVE",
	})

	key := sm.DeriveDetectorKey("keyed-model", "input-1", map[string]interface{}{"deviceId": "dev-42"})
	if key != "dev-42" {
		t.Fatalf("expected key 'dev-42', got %q", key)
	}

	key = sm.DeriveDetectorKey("keyed-model", "input-1", map[string]interface{}{"other": "val"})
	if key != "input-1" {
		t.Fatalf("expected fallback to inputName 'input-1', got %q", key)
	}
}

func TestStateMachine_DeriveDetectorKey_NoModel(t *testing.T) {
	sm := NewDetectorStateMachine(nil)
	key := sm.DeriveDetectorKey("nonexistent", "input-1", map[string]interface{}{})
	if key != "input-1" {
		t.Fatalf("expected fallback to inputName, got %q", key)
	}
}

func TestStateMachine_BatchEvaluate(t *testing.T) {
	sm := NewDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"transitionEvents": [{"condition": "true", "nextState": "Active"}]
				}
			},
			{"stateName": "Active"}
		]
	}`)

	sm.LoadModel(&DetectorModel{
		DetectorModelName:       "batch-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	messages := []InputMessage{
		{InputName: "input-A", Payload: map[string]interface{}{}},
		{InputName: "input-B", Payload: map[string]interface{}{}},
	}

	errs := sm.BatchEvaluate(context.Background(), messages)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	state, ok := sm.GetState("batch-model", "input-A")
	if !ok || state != "Active" {
		t.Fatalf("expected input-A 'Active', got %q (ok=%v)", state, ok)
	}

	state, ok = sm.GetState("batch-model", "input-B")
	if !ok || state != "Active" {
		t.Fatalf("expected input-B 'Active', got %q (ok=%v)", state, ok)
	}
}

// TestStateMachine_BatchEvaluate_InputRouting verifies that BatchEvaluate
// only delivers a message to detector models whose conditions reference
// the message's input via the $input.<inputName>.<field> syntax. A model
// that consumes TempInput must not react to a PressureInput message even
// if both payloads share a field name.
func TestStateMachine_BatchEvaluate_InputRouting(t *testing.T) {
	var mu sync.Mutex
	var fired []string
	onAction := func(modelName, key, actionType string, _ map[string]interface{}) {
		mu.Lock()
		fired = append(fired, modelName+"/"+actionType)
		mu.Unlock()
	}
	sm := NewDetectorStateMachine(onAction)

	tempModel := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"events": [
						{
							"condition": "$input.TempInput.temperature > 50",
							"actions": [{"action": "iotTopicPublish", "parameters": {"topic": "t/temp"}}]
						}
					]
				}
			}
		]
	}`)
	pressureModel := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"events": [
						{
							"condition": "$input.PressureInput.temperature > 50",
							"actions": [{"action": "iotTopicPublish", "parameters": {"topic": "t/pressure"}}]
						}
					]
				}
			}
		]
	}`)
	sm.LoadModel(&DetectorModel{DetectorModelName: "temp-model", DetectorModelDefinition: tempModel, Status: "ACTIVE"})
	sm.LoadModel(&DetectorModel{DetectorModelName: "pressure-model", DetectorModelDefinition: pressureModel, Status: "ACTIVE"})

	// Message declares InputName=PressureInput but carries a "temperature"
	// field that would match the temp-model condition if routing were
	// purely field-name based. Only pressure-model should fire.
	errs := sm.BatchEvaluate(context.Background(), []InputMessage{
		{InputName: "PressureInput", Payload: map[string]interface{}{"temperature": 99.0}},
	})
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 || fired[0] != "pressure-model/iotTopicPublish" {
		t.Fatalf("expected only pressure-model to fire, got %v", fired)
	}
}
