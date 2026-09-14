package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

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
	states, err := extractStatesFromDefinition(def)
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
		pre := history[len(history)-2]
		if last.EventId != pre.EventId+1 {
			t.Errorf("child %d terminal EventId = %d, want %d — the terminal event must take the next ID, not reuse the last state event's", i, last.EventId, pre.EventId+1)
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

// TestExecuteMapChildRedriveKeepsHistorySequence pins the child reclaim on
// redrive: the re-run continues the child's event-id sequence past the prior
// attempt's events instead of restarting at one — history keys embed the
// event id, so a restarted sequence would overwrite the earlier attempt's
// events and strand its terminal event mid-sequence.
func TestExecuteMapChildRedriveKeepsHistorySequence(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent3",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent3",
		Status:          "RUNNING",
		Input:           `{"v":[1]}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":          "Map",
				"ItemsPath":     "$.v",
				"ItemProcessor": distributedProcessor(),
				"End":           true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	parentEventId := int64(0)
	runMap := func(isRedrive bool) {
		execCtx := &ExecutionContext{
			Execution:     exec,
			Definition:    def,
			CurrentState:  "M",
			Input:         exec.Input,
			EventId:       &parentEventId,
			States:        map[string]sfnstore.State{},
			QueryLanguage: "JSONPath",
			MapItemIndex:  -1,
			IsRedrive:     isRedrive,
		}
		states, err := extractStatesFromDefinition(def)
		if err != nil {
			t.Fatalf("extract states failed: %v", err)
		}
		execCtx.States = states
		if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr == nil {
			t.Fatal("the failing iterator must fail the map")
		}
	}
	// Replace the shared processor's Pass with a Fail so the unit fails.
	if m := def.States["M"].(map[string]interface{}); true {
		proc := m["ItemProcessor"].(map[string]interface{})
		proc["States"] = map[string]interface{}{
			"F": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
		}
		proc["StartAt"] = "F"
	}

	runMap(false)
	runMap(true)

	children, err := store.ListAllExecutions(context.Background(), "", "", "", "")
	if err != nil {
		t.Fatalf("ListAllExecutions failed: %v", err)
	}
	if len(children) != 2 {
		t.Fatalf("expected parent + 1 child, got %d records", len(children))
	}
	var child *sfnstore.Execution
	for _, c := range children {
		if c.Name == "parent3:M-0" {
			child = c
		}
	}
	if child == nil {
		t.Fatal("child execution parent3:M-0 not found")
	}
	if child.RedriveCount != 1 {
		t.Errorf("child RedriveCount = %d, want 1", child.RedriveCount)
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), child.ExecutionArn, 100, "", false)
	if herr != nil {
		t.Fatalf("child history failed: %v", herr)
	}
	assertSequentialIds(t, history)
	started, redriven, failed := 0, 0, 0
	for _, ev := range history {
		switch ev.Type {
		case "ExecutionStarted":
			started++
		case "ExecutionRedriven":
			redriven++
			if ev.EventId != 4 {
				t.Errorf("ExecutionRedrived EventId = %d, want 4 — the re-run must append past the failed attempt's terminal event", ev.EventId)
			}
		case "ExecutionFailed":
			failed++
		}
	}
	if started != 1 {
		t.Errorf("child history has %d ExecutionStarted events, want exactly 1 from the original attempt", started)
	}
	if redriven != 1 || failed != 2 {
		var types []string
		for _, ev := range history {
			types = append(types, fmt.Sprintf("%d:%s", ev.EventId, ev.Type))
		}
		t.Errorf("child history has %d ExecutionRedriven / %d ExecutionFailed events, want 1 / 2 (both attempts' terminals survive); history: %v (arn %s)", redriven, failed, types, child.ExecutionArn)
	}

	// A zero-completed-results redrive reclaims the SAME Map Run — the
	// run's ARN is its identity across redrives, so no second run record
	// may appear and the parent's lifecycle event says MapRunRedriven.
	runs, runsErr := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if runsErr != nil || len(runs) != 1 {
		t.Fatalf("map runs after the redrive = %d (err %v), want 1 — the redrive must reclaim the prior run", len(runs), runsErr)
	}
	parentHistory, _, pherr := store.GetExecutionHistory(context.Background(), exec.ExecutionArn, 1000, "", false)
	if pherr != nil {
		t.Fatalf("parent history failed: %v", pherr)
	}
	runStarted, runRedriven := 0, 0
	for _, ev := range parentHistory {
		switch ev.Type {
		case "MapRunStarted":
			runStarted++
		case "MapRunRedriven":
			runRedriven++
		}
	}
	if runStarted != 1 || runRedriven != 1 {
		t.Errorf("parent map-run lifecycle = %d MapRunStarted / %d MapRunRedriven, want 1 / 1", runStarted, runRedriven)
	}
}

// TestExecuteMapChildRedriveResumesAtFailedState pins that a reclaimed
// child continues from its own failed state: the child's history is a
// complete execution record, so the iteration state that already succeeded
// must not re-run when the unit is redriven.
func TestExecuteMapChildRedriveResumesAtFailedState(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent4",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent4",
		Status:          "RUNNING",
		Input:           `{"v":[1]}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"ItemsPath": "$.v",
				"ItemProcessor": map[string]interface{}{
					"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
					"StartAt":         "Step1",
					"States": map[string]interface{}{
						"Step1": map[string]interface{}{"Type": "Pass", "Result": map[string]interface{}{"done": true}, "Next": "Step2"},
						"Step2": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
					},
				},
				"End": true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	parentEventId := int64(0)
	runMap := func(isRedrive bool) {
		execCtx := &ExecutionContext{
			Execution:     exec,
			Definition:    def,
			CurrentState:  "M",
			Input:         exec.Input,
			EventId:       &parentEventId,
			States:        map[string]sfnstore.State{},
			QueryLanguage: "JSONPath",
			MapItemIndex:  -1,
			IsRedrive:     isRedrive,
		}
		states, err := extractStatesFromDefinition(def)
		if err != nil {
			t.Fatalf("extract states failed: %v", err)
		}
		execCtx.States = states
		if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr == nil {
			t.Fatal("the failing iterator must fail the map")
		}
	}

	runMap(false)
	runMap(true)

	children, err := store.ListAllExecutions(context.Background(), "", "", "", "")
	if err != nil {
		t.Fatalf("ListAllExecutions failed: %v", err)
	}
	var child *sfnstore.Execution
	for _, c := range children {
		if c.Name == "parent4:M-0" {
			child = c
		}
	}
	if child == nil {
		t.Fatal("child execution parent4:M-0 not found")
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), child.ExecutionArn, 100, "", false)
	if herr != nil {
		t.Fatalf("child history failed: %v", herr)
	}
	assertSequentialIds(t, history)
	step1Entered, step2Entered := 0, 0
	for _, ev := range history {
		if ev.StateEnteredEventDetails == nil {
			continue
		}
		switch ev.StateEnteredEventDetails.Name {
		case "Step1":
			if ev.Type == "PassStateEntered" {
				step1Entered++
			}
		case "Step2":
			if ev.Type == "FailStateEntered" {
				step2Entered++
			}
		}
	}
	if step1Entered != 1 {
		t.Errorf("Step1 entered %d times across the original run and the redrive, want 1 — the succeeded state must not re-run", step1Entered)
	}
	if step2Entered != 2 {
		t.Errorf("Step2 entered %d times, want 2 — both attempts run the failed state", step2Entered)
	}
}

// TestExecuteMapChildRedriveParksQueuedChildren pins the Map Run side of
// the child redrive contract on a concurrency-limited redrive: a unit that
// re-acquires a slot flips its child to RUNNING on disk immediately (not
// at unit completion), a unit still queued parks its prior child as
// PENDING_REDRIVE — surfaced to ListExecutions and reported
// REDRIVABLE_BY_MAP_RUN — and the parked child is reclaimed once a slot
// frees.
func TestExecuteMapChildRedriveParksQueuedChildren(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent6",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent6",
		Status:          "RUNNING",
		Input:           `{"v":[1,2]}`,
	}
	iterator := func(fail bool) map[string]interface{} {
		if fail {
			return map[string]interface{}{
				"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
				"StartAt":         "F",
				"States": map[string]interface{}{
					"F": map[string]interface{}{"Type": "Fail", "Error": "ItemFailed", "Cause": "deliberate"},
				},
			}
		}
		return map[string]interface{}{
			"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
			"StartAt":         "W",
			"States": map[string]interface{}{
				"W": map[string]interface{}{"Type": "Wait", "Seconds": 2, "Next": "P"},
				"P": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
			},
		}
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":           "Map",
				"ItemsPath":      "$.v",
				"MaxConcurrency": 1,
				"ItemProcessor":  iterator(true),
				"End":            true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	parentEventId := int64(0)
	newExecCtx := func(isRedrive bool) *ExecutionContext {
		execCtx := &ExecutionContext{
			Execution:     exec,
			Definition:    def,
			CurrentState:  "M",
			Input:         exec.Input,
			EventId:       &parentEventId,
			States:        map[string]sfnstore.State{},
			QueryLanguage: "JSONPath",
			MapItemIndex:  -1,
			IsRedrive:     isRedrive,
		}
		states, err := extractStatesFromDefinition(def)
		if err != nil {
			t.Fatalf("extract states failed: %v", err)
		}
		execCtx.States = states
		return execCtx
	}
	runMap := func(isRedrive bool) error {
		ctx := newExecCtx(isRedrive)
		_, _, execErr := e.executeMap(context.Background(), ctx, ctx.States["M"].(*sfnstore.MapState))
		if execErr == nil {
			return nil
		}
		return execErr
	}
	probeCtx := newExecCtx(false)
	arn0 := e.mapChildARN(probeCtx, probeCtx.States["M"].(*sfnstore.MapState), 0, 0)
	arn1 := e.mapChildARN(probeCtx, probeCtx.States["M"].(*sfnstore.MapState), 1, 0)

	// The original attempt fails both children under the concurrency limit.
	if err := runMap(false); err == nil {
		t.Fatal("the failing iterator must fail the map")
	}
	for _, arn := range []string{arn0, arn1} {
		child, err := store.GetExecution(context.Background(), arn)
		if err != nil {
			t.Fatalf("get child %s: %v", arn, err)
		}
		if child.Status != "FAILED" {
			t.Fatalf("child %s status = %s, want FAILED after the original attempt", arn, child.Status)
		}
	}

	// The redrive re-runs both units through a two-second iterator, so the
	// mid-run window is observable: the slot holder is RUNNING on disk and
	// the queued unit's child is parked PENDING_REDRIVE.
	def.States["M"].(map[string]interface{})["ItemProcessor"] = iterator(false)
	type runResult struct {
		err error
	}
	done := make(chan runResult, 1)
	go func() {
		done <- runResult{runMap(true)}
	}()

	sawRunning, sawParked := false, false
	var parkedRecord *sfnstore.Execution
	deadline := time.Now().Add(5 * time.Second)
	for (!sawRunning || !sawParked) && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
		first, err1 := store.GetExecution(context.Background(), arn0)
		second, err2 := store.GetExecution(context.Background(), arn1)
		if err1 != nil || err2 != nil {
			t.Fatalf("probe children: %v %v", err1, err2)
		}
		// Either unit may win the race for the single slot, so the probe is
		// order-agnostic: one child is RUNNING while the other is parked.
		for _, probe := range []*sfnstore.Execution{first, second} {
			if probe.Status == "RUNNING" {
				sawRunning = true
			}
			if probe.Status == "PENDING_REDRIVE" {
				sawParked = true
				parkedRecord = probe
			}
		}
	}
	if !sawRunning {
		t.Error("the reclaimed child never appeared as RUNNING on disk mid-run — the reclaim must persist RUNNING immediately")
	}
	if !sawParked {
		t.Error("the queued unit's child never appeared as PENDING_REDRIVE — a redriven unit waiting for a Map Run slot must park its child")
	}
	if sawParked {
		if status, reason := computeRedriveStatus(context.Background(), store, parkedRecord); status != "REDRIVABLE_BY_MAP_RUN" {
			t.Errorf("parked child redriveStatus = %s (%s), want REDRIVABLE_BY_MAP_RUN", status, reason)
		}
	}

	select {
	case res := <-done:
		if res.err != nil {
			t.Fatalf("the redriven map must succeed, got %v", res.err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("the redriven map did not complete")
	}
	for _, arn := range []string{arn0, arn1} {
		child, err := store.GetExecution(context.Background(), arn)
		if err != nil {
			t.Fatalf("get child %s: %v", arn, err)
		}
		if child.Status != "SUCCEEDED" {
			t.Errorf("child %s final status = %s, want SUCCEEDED", arn, child.Status)
		}
		if child.RedriveCount != 1 {
			t.Errorf("child %s RedriveCount = %d, want 1", arn, child.RedriveCount)
		}
	}
}

// TestMapChildReclaimRejectsNonReclaimableStatus pins the reclaim guard: a
// Map Run may only reclaim a child from an unsuccessful terminal status,
// PENDING_REDRIVE, or a crash-leftover RUNNING record — reclaiming a
// SUCCEEDED child would re-run completed work, so the transition is
// refused and the record left untouched.
func TestMapChildReclaimRejectsNonReclaimableStatus(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent7",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent7",
		Status:          "RUNNING",
		Input:           `{"v":[1]}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":          "Map",
				"ItemsPath":     "$.v",
				"ItemProcessor": distributedProcessor(),
				"End":           true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	parentEventId := int64(0)
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "M",
		Input:         exec.Input,
		EventId:       &parentEventId,
		States:        map[string]sfnstore.State{},
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
		IsRedrive:     true,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}
	mapState := execCtx.States["M"].(*sfnstore.MapState)
	childArn := e.mapChildARN(execCtx, mapState, 0, 0)
	child, err := store.GetExecution(context.Background(), childArn)
	if err != nil {
		t.Fatalf("get child: %v", err)
	}
	if child.Status != "SUCCEEDED" {
		t.Fatalf("child status = %s, want SUCCEEDED before the reclaim attempt", child.Status)
	}

	_, _, _, rerr := e.beginMapChildExecution(context.Background(), execCtx, mapState, child.MapRunArn, mapWorkUnit{StartIndex: 0, ItemCount: 1, InputJSON: "1"}, 0, mapState.GetIterator(), 0)
	if !errors.Is(rerr, sfnstore.ErrExecutionNotRedrivable) {
		t.Fatalf("reclaim of a SUCCEEDED child = %v, want ErrExecutionNotRedrivable", rerr)
	}
	after, err := store.GetExecution(context.Background(), childArn)
	if err != nil {
		t.Fatalf("get child after the rejected reclaim: %v", err)
	}
	if after.Status != "SUCCEEDED" || after.RedriveCount != 0 {
		t.Errorf("rejected reclaim must leave the child untouched, got status %s redriveCount %d", after.Status, after.RedriveCount)
	}
}

// TestExecuteMapRetryRedispatchesAllUnits pins the Distributed Map retry
// contract: "The retry policy applies to all of the child workflow
// executions and not just the failed execution. ... When you retry a Map
// state, it creates a new Map Run." A state-level retry must therefore
// dispatch every unit as a fresh child of a new Map Run — succeeded units
// included — never collide with the prior attempt's children.
func TestExecuteMapRetryRedispatchesAllUnits(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent8",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent8",
		Status:          "RUNNING",
		Input:           `{"v":[{"n":2},{"n":1}]}`,
	}
	iterator := map[string]interface{}{
		"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
		"StartAt":         "IsLast",
		"States": map[string]interface{}{
			"IsLast": map[string]interface{}{
				"Type": "Choice",
				"Choices": []interface{}{
					map[string]interface{}{"Variable": "$.n", "NumericEquals": 1, "Next": "Explode"},
				},
				"Default": "Done",
			},
			"Explode": map[string]interface{}{"Type": "Fail", "Error": "ItemFailed", "Cause": "deliberate unit failure"},
			"Done":    map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
		},
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":           "Map",
				"ItemsPath":      "$.v",
				"MaxConcurrency": 2,
				"ItemProcessor":  iterator,
				"Retry": []interface{}{
					map[string]interface{}{"ErrorEquals": []interface{}{"States.ALL"}, "MaxAttempts": 1, "IntervalSeconds": 0},
				},
				"End": true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	parentEventId := int64(0)
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "M",
		Input:         exec.Input,
		EventId:       &parentEventId,
		States:        map[string]sfnstore.State{},
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	_, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr == nil {
		t.Fatal("the map with one permanently failing unit must fail once retries are exhausted")
	}
	// The zero-default thresholds make the failed unit exceed the
	// threshold; the reported failure must still carry the iteration's own
	// error, never an internal dispatch error from a succeeded unit's
	// reclaim rejection.
	if execErr.ErrorCode != "States.ExceedToleratedFailureThreshold" || !strings.Contains(execErr.Cause, "ItemFailed") {
		t.Errorf("map error = %s (%s), want States.ExceedToleratedFailureThreshold carrying the iteration's ItemFailed failure", execErr.ErrorCode, execErr.Cause)
	}

	mapState := execCtx.States["M"].(*sfnstore.MapState)
	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil {
		t.Fatalf("list map runs: %v", err)
	}
	if len(runs) != 2 {
		t.Fatalf("ListMapRunsByExecution returned %d runs, want 2 (one per attempt)", len(runs))
	}
	if runs[0].Name != "M" || runs[1].Name != "M" {
		t.Errorf("run names = %q/%q, want M/M", runs[0].Name, runs[1].Name)
	}
	if runs[0].Attempt != 0 || runs[1].Attempt != 1 {
		t.Errorf("run attempts = %d/%d, want 0/1", runs[0].Attempt, runs[1].Attempt)
	}
	if runs[1].Status != "FAILED" || runs[1].ExecutionCounts.Succeeded != 1 || runs[1].ExecutionCounts.Failed != 1 {
		t.Errorf("attempt-1 run status = %s counts %+v, want FAILED with 1 succeeded and 1 failed", runs[1].Status, runs[1].ExecutionCounts)
	}

	mustChild := func(ordinal int, attempt int64) *sfnstore.Execution {
		t.Helper()
		arn := e.mapChildARN(execCtx, mapState, ordinal, attempt)
		child, cerr := store.GetExecution(context.Background(), arn)
		if cerr != nil || child == nil {
			t.Fatalf("get child (ordinal %d, attempt %d) %s: %v", ordinal, attempt, arn, cerr)
		}
		return child
	}

	// The prior attempt's children keep their own outcomes: the succeeded
	// unit is untouched (never re-run, never reclassified), the failed one
	// keeps its failure.
	first0 := mustChild(0, 0)
	if first0.Status != "SUCCEEDED" || first0.RedriveCount != 0 {
		t.Errorf("attempt-0 unit 0 status = %s redriveCount = %d, want SUCCEEDED/0", first0.Status, first0.RedriveCount)
	}
	first1 := mustChild(1, 0)
	if first1.Status != "FAILED" || first1.Error != "ItemFailed" || first1.RedriveCount != 0 {
		t.Errorf("attempt-0 unit 1 status = %s error = %s redriveCount = %d, want FAILED/ItemFailed/0 (a retry never redrives a prior attempt's child)", first1.Status, first1.Error, first1.RedriveCount)
	}
	if first1.MapRunArn == runs[1].MapRunArn {
		t.Error("attempt-0 unit 1 must belong to the attempt-0 Map Run, not the retry's")
	}

	// The retry's children are fresh executions of the new run: no
	// redrive markers, the succeeded unit re-ran and succeeded again, and
	// both belong to the attempt-1 run.
	second0 := mustChild(0, 1)
	if second0.Status != "SUCCEEDED" || second0.RedriveCount != 0 {
		t.Errorf("attempt-1 unit 0 status = %s redriveCount = %d, want SUCCEEDED/0 (a fresh execution)", second0.Status, second0.RedriveCount)
	}
	second1 := mustChild(1, 1)
	if second1.Status != "FAILED" || second1.Error != "ItemFailed" || second1.RedriveCount != 0 {
		t.Errorf("attempt-1 unit 1 status = %s error = %s redriveCount = %d, want FAILED/ItemFailed/0 (a fresh execution)", second1.Status, second1.Error, second1.RedriveCount)
	}
	if second0.MapRunArn != runs[1].MapRunArn || second1.MapRunArn != runs[1].MapRunArn {
		t.Error("the retry's children must carry the new Map Run's ARN")
	}

	// A fresh child's history starts with ExecutionStarted and holds no
	// ExecutionRedrived event — the collision-with-prior-attempt path is
	// exactly what the attempt-stamped names exist to prevent.
	history, _, herr := store.GetExecutionHistory(context.Background(), second0.ExecutionArn, 100, "", false)
	if herr != nil {
		t.Fatalf("get attempt-1 unit 0 history: %v", herr)
	}
	if len(history) == 0 || history[0].Type != "ExecutionStarted" {
		t.Errorf("attempt-1 unit 0 history head = %+v, want ExecutionStarted", history)
	}
	for _, evt := range history {
		if evt.Type == "ExecutionRedriven" {
			t.Error("attempt-1 unit 0 history must not contain ExecutionRedriven — a retry dispatches a fresh execution, not a redrive")
		}
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
	states, err := extractStatesFromDefinition(def)
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
	states, err := extractStatesFromDefinition(def)
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

// TestMapMidTaskCancellationClassifiesAborted pins the mid-task
// cancellation contract for a Distributed Map: a StopExecution landing
// while a unit's task is in flight must abort every layer — the child
// execution ABORTED, the Map Run ABORTED with the MapRunAborted event —
// and no Catcher may consume the interruption.
func TestMapMidTaskCancellationClassifiesAborted(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/parent9",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "parent9",
		Status:          "RUNNING",
		Input:           `{"v":[{"n":1}]}`,
	}
	iterator := map[string]interface{}{
		"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
		"StartAt":         "Work",
		"States": map[string]interface{}{
			"Work": map[string]interface{}{"Type": "Task", "Resource": tokenTestActivityARN, "TimeoutSeconds": 10, "End": true},
		},
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":          "Map",
				"ItemsPath":     "$.v",
				"ItemProcessor": iterator,
				"Catch": []interface{}{
					map[string]interface{}{"ErrorEquals": []interface{}{"States.ALL"}, "Next": "Fallback"},
				},
			},
			"Fallback": map[string]interface{}{"Type": "Pass", "Result": "recovered", "End": true},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	parentEventId := int64(0)
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "M",
		Input:         exec.Input,
		EventId:       &parentEventId,
		States:        map[string]sfnstore.State{},
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- e.executeStates(ctx, execCtx) }()

	// A claimed activity task means the unit's child execution is in
	// flight and waiting on the worker.
	workerCtx, wcancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer wcancel()
	if _, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1"); err != nil {
		t.Fatalf("no activity task became available: %v", err)
	}

	cancel()

	select {
	case runErr := <-done:
		if !errors.Is(runErr, context.Canceled) {
			t.Fatalf("executeStates returned %v, want a context.Canceled-classified error", runErr)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeStates did not return after the cancellation")
	}

	childArn := e.mapChildARN(execCtx, execCtx.States["M"].(*sfnstore.MapState), 0, 0)
	child, cerr := store.GetExecution(context.Background(), childArn)
	if cerr != nil {
		t.Fatalf("get child: %v", cerr)
	}
	if child.Status != "ABORTED" {
		t.Errorf("child status = %s, want ABORTED — a cancelled unit is aborted, not failed", child.Status)
	}
	childHistory, _, herr := store.GetExecutionHistory(context.Background(), childArn, 200, "", false)
	if herr != nil {
		t.Fatalf("child history failed: %v", herr)
	}
	childAborted, childFailed := false, false
	for _, ev := range childHistory {
		switch ev.Type {
		case "ExecutionAborted":
			childAborted = true
		case "ExecutionFailed":
			childFailed = true
		}
	}
	if !childAborted || childFailed {
		t.Errorf("child history aborted=%v failed=%v, want the ExecutionAborted terminal alone", childAborted, childFailed)
	}

	runs, rerr := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if rerr != nil || len(runs) != 1 {
		t.Fatalf("list map runs: %v (len %d)", rerr, len(runs))
	}
	if runs[0].Status != "ABORTED" {
		t.Errorf("map run status = %s, want ABORTED", runs[0].Status)
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), exec.ExecutionArn, 200, "", false)
	if herr != nil {
		t.Fatalf("history failed: %v", herr)
	}
	sawAborted, sawFailed, sawFallback := false, false, false
	for _, ev := range history {
		switch ev.Type {
		case "MapRunAborted":
			sawAborted = true
		case "MapRunFailed":
			sawFailed = true
		}
		if ev.StateEnteredEventDetails != nil && ev.StateEnteredEventDetails.Name == "Fallback" {
			sawFallback = true
		}
	}
	if !sawAborted {
		t.Error("history lacks MapRunAborted — a stopped run is ABORTED, matching the MapRunAborted history event")
	}
	if sawFailed {
		t.Error("history records MapRunFailed for a cancellation — the interruption is not a run failure")
	}
	if sawFallback {
		t.Error("the States.ALL Catcher consumed the cancellation — no recovery transition may run on a stopping execution")
	}
}
