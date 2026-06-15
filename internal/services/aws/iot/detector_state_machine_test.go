package iot

import (
	"context"
	"sync"
	"testing"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

func TestStateMachine_TwoStateTransition(t *testing.T) {
	var mu sync.Mutex
	var actions []string
	onAction := func(modelName, key, actionType string, _ map[string]interface{}) {
		mu.Lock()
		actions = append(actions, actionType)
		mu.Unlock()
	}

	sm := newDetectorStateMachine(onAction)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Idle",
				"onInput": {
					"events": [
						{
							"transitions": [
								{
									"condition": "true",
									"nextState": "Active"
								}
							]
						}
					]
				},
				"onExit": [
					{"action": "log", "parameters": {"msg": "exiting Idle"}}
				]
			},
			{
				"stateName": "Active",
				"onEnter": [
					{"action": "notify", "parameters": {"msg": "entered Active"}}
				]
			}
		]
	}`)

	dm := &iotstore.DetectorModel{
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

	_ = sm.EvaluateEvent(context.Background(), "test-model", "key1", map[string]interface{}{"temp": 60.0})
	_, _ = sm.GetState("test-model", "key1")

	mu.Lock()
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions (exit + enter), got %d: %v", len(actions), actions)
	}
	if actions[0] != "log" {
		t.Fatalf("expected first action 'log', got %q", actions[0])
	}
	if actions[1] != "notify" {
		t.Fatalf("expected second action 'notify', got %q", actions[1])
	}
	mu.Unlock()
}

func TestStateMachine_ConditionMatch(t *testing.T) {
	sm := newDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Normal",
				"onInput": {
					"events": [
						{
							"transitions": [
								{
									"condition": "temp > 50",
									"nextState": "HighTemp"
								}
							]
						}
					]
				}
			},
			{
				"stateName": "HighTemp"
			}
		]
	}`)

	sm.LoadModel(&iotstore.DetectorModel{
		DetectorModelName:       "temp-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	_ = sm.EvaluateEvent(context.Background(), "temp-model", "sensor1", map[string]interface{}{"temp": 30.0})
	state, _ := sm.GetState("temp-model", "sensor1")
	if state != "Normal" {
		t.Fatalf("expected 'Normal' for temp=30, got %q", state)
	}

	_ = sm.EvaluateEvent(context.Background(), "temp-model", "sensor1", map[string]interface{}{"temp": 75.0})
	state, _ = sm.GetState("temp-model", "sensor1")
	if state != "HighTemp" {
		t.Fatalf("expected 'HighTemp' for temp=75, got %q", state)
	}
}

func TestStateMachine_MultipleInstances(t *testing.T) {
	sm := newDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Off",
				"onInput": {
					"events": [{"transitions": [{"condition": "trigger == true", "nextState": "On"}]}]
				}
			},
			{"stateName": "On"}
		]
	}`)

	sm.LoadModel(&iotstore.DetectorModel{
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
	sm := newDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "A",
				"onInput": {
					"events": [{"transitions": [{"condition": "true", "nextState": "B"}]}]
				}
			},
			{"stateName": "B"}
		]
	}`)

	sm.LoadModel(&iotstore.DetectorModel{
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
	sm := newDetectorStateMachine(nil)
	err := sm.EvaluateEvent(context.Background(), "no-such-model", "k1", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for nonexistent model")
	}
}

func TestStateMachine_StringConditionEquality(t *testing.T) {
	sm := newDetectorStateMachine(nil)

	def := detectorModelDefinitionJSON(`{
		"states": [
			{
				"stateName": "Locked",
				"onInput": {
					"events": [{"transitions": [{"condition": "status == 'unlocked'", "nextState": "Unlocked"}]}]
				}
			},
			{"stateName": "Unlocked"}
		]
	}`)

	sm.LoadModel(&iotstore.DetectorModel{
		DetectorModelName:       "door-model",
		DetectorModelDefinition: def,
		Status:                  "ACTIVE",
	})

	_ = sm.EvaluateEvent(context.Background(), "door-model", "door-1", map[string]interface{}{"status": "locked"})
	state, _ := sm.GetState("door-model", "door-1")
	if state != "Locked" {
		t.Fatalf("expected 'Locked', got %q", state)
	}

	_ = sm.EvaluateEvent(context.Background(), "door-model", "door-1", map[string]interface{}{"status": "unlocked"})
	state, _ = sm.GetState("door-model", "door-1")
	if state != "Unlocked" {
		t.Fatalf("expected 'Unlocked', got %q", state)
	}
}
