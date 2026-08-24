package sfn

import (
	"context"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestDirectInputOutputPathApplied pins the core Amazon States Language
// contract: InputPath and OutputPath are direct state members and the
// state input/output filtering applies to them. A Pass state that keeps
// $.keep and then projects $.keep.n must yield the inner value.
func TestDirectInputOutputPathApplied(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "P",
		States: map[string]interface{}{
			"P": map[string]interface{}{
				"Type":       "Pass",
				"InputPath":  "$.keep",
				"OutputPath": "$.n",
				"End":        true,
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/io1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "io1",
		Status:          "RUNNING",
		Input:           `{"keep":{"n":42},"drop":1}`,
	}
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "P",
		Input:         exec.Input,
		EventId:       ptrEventID(),
		States:        states,
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	if err := e.executeStates(context.Background(), execCtx); err != nil {
		t.Fatalf("executeStates failed: %v", err)
	}
	if execCtx.Output != "42" {
		t.Errorf("output = %s, want 42", execCtx.Output)
	}
}

// TestDirectOutputPathOnSucceed pins that a Succeed state also honours a
// direct OutputPath member on its way out.
func TestDirectOutputPathOnSucceed(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "S",
		States: map[string]interface{}{
			"S": map[string]interface{}{
				"Type":       "Succeed",
				"OutputPath": "$.value",
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/io2",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "io2",
		Status:          "RUNNING",
		Input:           `{"value":{"deep":true},"other":2}`,
	}
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "S",
		Input:         exec.Input,
		EventId:       ptrEventID(),
		States:        states,
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	if err := e.executeStates(context.Background(), execCtx); err != nil {
		t.Fatalf("executeStates failed: %v", err)
	}
	if execCtx.Output != `{"deep":true}` {
		t.Errorf("output = %s, want {\"deep\":true}", execCtx.Output)
	}
}

// TestFailPathResolution pins the Fail ErrorPath and CausePath contract:
// reference paths select the error name and cause strings from the state
// input.
func TestFailPathResolution(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "F",
		States: map[string]interface{}{
			"F": map[string]interface{}{
				"Type":      "Fail",
				"ErrorPath": "$.code",
				"CausePath": "$.reason",
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/fp1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "fp1",
		Status:          "RUNNING",
		Input:           `{"code":"CustomError","reason":"something broke"}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "F",
		Input: exec.Input, EventId: ptrEventID(), States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	err = e.executeStates(context.Background(), execCtx)
	if err == nil {
		t.Fatal("Fail state must fail the execution")
	}
	msg := err.Error()
	if !strings.Contains(msg, "CustomError") || !strings.Contains(msg, "something broke") {
		t.Errorf("Fail error = %q, want resolved CustomError / something broke", msg)
	}
}

// TestMapMaxConcurrencyPathResolution pins that MaxConcurrencyPath reads
// the concurrency ceiling from the state input.
func TestMapMaxConcurrencyPathResolution(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":               "Map",
				"ItemsPath":          "$.v",
				"MaxConcurrencyPath": "$.mc",
				"ItemProcessor": map[string]interface{}{
					"StartAt": "W",
					"States": map[string]interface{}{
						"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
					},
				},
				"End": true,
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/mcp1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "mcp1",
		Status:          "RUNNING",
		Input:           `{"v":[1,2,3],"mc":2}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: exec.Input, EventId: ptrEventID(), States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	if _, _, execErr := e.executeMap(context.Background(), execCtx, states["M"].(*sfnstore.MapState)); execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}
	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("map runs = %v, %v", runs, err)
	}
	if runs[0].MaxConcurrency != 2 {
		t.Errorf("MapRun MaxConcurrency = %d, want 2 from MaxConcurrencyPath", runs[0].MaxConcurrency)
	}
}

// TestResolveTaskSecondsPath pins the TimeoutSecondsPath resolver: a
// reference path selecting a positive integer resolves, non-integer and
// missing selections are runtime input errors.
func TestResolveTaskSecondsPath(t *testing.T) {
	v, err := resolveTaskSecondsPath(`{"cfg":{"timeout":30}}`, "$.cfg.timeout", "TimeoutSecondsPath")
	if err != nil || v != 30 {
		t.Errorf("resolution = (%d, %v), want (30, nil)", v, err)
	}
	if _, err := resolveTaskSecondsPath(`{"cfg":{"timeout":1.5}}`, "$.cfg.timeout", "TimeoutSecondsPath"); err == nil {
		t.Error("non-integer timeout accepted")
	}
	if _, err := resolveTaskSecondsPath(`{"cfg":{}}`, "$.cfg.timeout", "TimeoutSecondsPath"); err == nil {
		t.Error("missing timeout accepted")
	}
	if _, err := resolveTaskSecondsPath(`{"cfg":{"timeout":0}}`, "$.cfg.timeout", "TimeoutSecondsPath"); err == nil {
		t.Error("zero timeout accepted")
	}
}
