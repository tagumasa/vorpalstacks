package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestExecuteMapChildExecutions pins the Distributed Map child-execution
// dispatch: every item runs as its own child workflow execution under the
// parent state machine's namespace with the parent-scoped naming, its own
// history and terminal status, and ListExecutions scoped by the Map Run
// ARN returns exactly those children.
func TestExecuteMapChildExecutions(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent",
		Status:          "RUNNING",
		Input:           `{"v":[1,2,3]}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"Label":     "dist",
				"ItemsPath": "$.v",
				"ItemProcessor": map[string]interface{}{
					"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
					"StartAt":         "W",
					"States": map[string]interface{}{
						"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
					},
				},
				"End": true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "M",
		Input:         exec.Input,
		EventId:       ptrEventID(),
		States:        map[string]sfnstore.State{},
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}

	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("map runs = %v, %v", runs, err)
	}
	mapRunArn := runs[0].MapRunArn

	children, err := store.ListAllExecutions(context.Background(), "", "", mapRunArn, "")
	if err != nil {
		t.Fatalf("ListAllExecutions failed: %v", err)
	}
	if len(children) != 3 {
		t.Fatalf("expected 3 child executions, got %d: %+v", len(children), children)
	}
	for i, child := range children {
		wantName := fmt.Sprintf("parent:dist-%d", i)
		if child.Name != wantName {
			t.Errorf("child %d Name = %q, want %q", i, child.Name, wantName)
		}
		if child.Status != "SUCCEEDED" {
			t.Errorf("child %d Status = %q", i, child.Status)
		}
		if child.MapRunArn != mapRunArn {
			t.Errorf("child %d MapRunArn = %q", i, child.MapRunArn)
		}
		if child.ItemCount != 1 {
			t.Errorf("child %d ItemCount = %d, want 1", i, child.ItemCount)
		}
		var input interface{}
		if err := json.Unmarshal([]byte(child.Input), &input); err != nil {
			t.Errorf("child %d input not JSON: %v (%s)", i, err, child.Input)
		}
		history, _, herr := store.GetExecutionHistory(context.Background(), child.ExecutionArn, 100, "", false)
		if herr != nil {
			t.Fatalf("child history failed: %v", herr)
		}
		first, last := history[0], history[len(history)-1]
		if first.Type != "ExecutionStarted" || last.Type != "ExecutionSucceeded" {
			t.Errorf("child %d history endpoints = %s..%s, want ExecutionStarted..ExecutionSucceeded", i, first.Type, last.Type)
		}
	}

	// Listing by the parent ARN includes the parent and its children.
	all, err := store.ListAllExecutions(context.Background(), exec.StateMachineArn, "", "", "")
	if err != nil {
		t.Fatalf("ListAllExecutions(sm) failed: %v", err)
	}
	if len(all) != 4 {
		t.Errorf("expected parent + 3 children in the state machine listing, got %d", len(all))
	}
}

// TestExecuteMapChildExecutionsBatched pins the child dispatch under an
// ItemBatcher: one child per batch, itemCount carries the batch's item
// count and the child input is the {"Items": [...]} payload.
func TestExecuteMapChildExecutionsBatched(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent2",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent2",
		Status:          "RUNNING",
		Input:           `{"v":[1,2,3]}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":          "Map",
				"ItemsPath":     "$.v",
				"ItemProcessor": distributedProcessor(),
				"ItemBatcher": map[string]interface{}{
					"MaxItemsPerBatch": 2,
				},
				"End": true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "M",
		Input:         exec.Input,
		EventId:       ptrEventID(),
		States:        map[string]sfnstore.State{},
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}

	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("map runs = %v, %v", runs, err)
	}
	children, err := store.ListAllExecutions(context.Background(), "", "", runs[0].MapRunArn, "")
	if err != nil {
		t.Fatalf("ListAllExecutions failed: %v", err)
	}
	if len(children) != 2 {
		t.Fatalf("expected 2 child executions (batches of 2 and 1), got %d", len(children))
	}
	if children[0].ItemCount != 2 || children[1].ItemCount != 1 {
		t.Errorf("child item counts = %d, %d; want 2, 1", children[0].ItemCount, children[1].ItemCount)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(children[0].Input), &payload); err != nil {
		t.Fatalf("first child input not JSON: %v (%s)", err, children[0].Input)
	}
	if items, ok := payload["Items"].([]interface{}); !ok || len(items) != 2 {
		t.Errorf("first child input Items = %v", payload["Items"])
	}
	if runs[0].ItemCounts.Total != 3 || runs[0].ExecutionCounts.Total != 2 {
		t.Errorf("counts: items %+v executions %+v", runs[0].ItemCounts, runs[0].ExecutionCounts)
	}
}

// TestExecuteMapChildExecutionsInlineNoChildren pins the boundary: inline
// Map iterations never create child execution records.
func TestExecuteMapChildExecutionsInlineNoChildren(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/inline",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "inline",
		Status:          "RUNNING",
		Input:           `{"v":[1,2]}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"ItemsPath": "$.v",
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
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "M",
		Input:         exec.Input,
		EventId:       ptrEventID(),
		States:        map[string]sfnstore.State{},
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}

	all, err := store.ListAllExecutions(context.Background(), exec.StateMachineArn, "", "", "")
	if err != nil {
		t.Fatalf("ListAllExecutions failed: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("inline map must not create child executions, got %d records", len(all))
	}
}
