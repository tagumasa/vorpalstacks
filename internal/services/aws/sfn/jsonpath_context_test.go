package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestContextObjectFullNodeSet pins that the JSONPath context object
// resolves the full documented node set (Execution incl. Input,
// RedriveCount, RedriveTime; State; Map.Item incl. Source) — the same
// object the JSONata plane exposes as $states.context.
func TestContextObjectFullNodeSet(t *testing.T) {
	e := &Executor{}
	now := time.Now().UTC()
	execCtx := &ExecutionContext{
		Execution: &sfnstore.Execution{
			ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:e1",
			Name:         "e1",
			Input:        `{"k":42}`,
			StartDate:    now,
			RedriveCount: 2,
			RedriveDate:  now,
		},
		CurrentState:     "S",
		StateEnteredTime: now,
		RetryCount:       1,
		MapItemIndex:     3,
		MapItemValue:     map[string]interface{}{"who": "joe"},
	}

	cases := map[string]interface{}{
		"$$.Execution.Id":           execCtx.Execution.ExecutionArn,
		"$$.Execution.Name":         "e1",
		"$$.Execution.Input":        map[string]interface{}{"k": float64(42)},
		"$$.Execution.RedriveCount": int64(2),
		"$$.Execution.RedriveTime":  now.Format(time.RFC3339),
		"$$.State.Name":             "S",
		"$$.State.RetryCount":       int32(1),
		"$$.Map.Item.Index":         3,
		"$$.Map.Item.Value":         map[string]interface{}{"who": "joe"},
		"$$.Map.Item.Source":        "STATE_DATA",
	}
	for path, want := range cases {
		got, err := e.getContextValue(execCtx, "", path)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("%s = %#v, want %#v", path, got, want)
		}
	}

	// Map context exists only inside a Map iteration (MapItemIndex's
	// out-of-Map convention is -1).
	if _, err := e.getContextValue(&ExecutionContext{Execution: execCtx.Execution, MapItemIndex: -1}, "", "$$.Map.Item.Index"); err == nil {
		t.Error("$$.Map.Item.Index outside a Map iteration must not resolve")
	}
}

// TestContextRootSelectsWholeContextObject pins the bare "$$": the Amazon
// States Language defines $$-rooted paths by stripping the first dollar
// sign and applying the remainder as the JSONPath against the context
// object, so the bare form — the root path — selects the entire context
// object on every documented $$. field (InputPath, payload templates,
// intrinsic arguments), never an input-plane key miss.
func TestContextRootSelectsWholeContextObject(t *testing.T) {
	e := &Executor{}
	now := time.Now().UTC()
	execCtx := &ExecutionContext{
		Execution: &sfnstore.Execution{
			ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:e1",
			Name:         "e1",
			Input:        `{"k":42}`,
			StartDate:    now,
		},
		CurrentState:     "S",
		StateEnteredTime: now,
	}

	whole, err := e.getContextValue(execCtx, "", "$$")
	if err != nil {
		t.Fatalf("the bare context root must resolve: %v", err)
	}
	wholeMap, ok := whole.(map[string]interface{})
	if !ok {
		t.Fatalf("the bare context root selects the whole object, got %T", whole)
	}
	if got := wholeMap["Execution"].(map[string]interface{})["Id"]; got != execCtx.Execution.ExecutionArn {
		t.Errorf("context root Execution.Id = %v, want %s", got, execCtx.Execution.ExecutionArn)
	}

	selected, execErr := e.applyInputPath(execCtx, `{"k":1}`, "$$")
	if execErr != nil {
		t.Fatalf("InputPath $$ must select the context object: %v", execErr.Cause)
	}
	var inputSel map[string]interface{}
	if err := json.Unmarshal([]byte(selected), &inputSel); err != nil {
		t.Fatalf("InputPath $$ output not valid JSON: %v (%s)", err, selected)
	}
	if got := inputSel["Execution"].(map[string]interface{})["Name"]; got != "e1" {
		t.Errorf("InputPath $$ Execution.Name = %v, want e1", got)
	}
	if _, hasInputKey := inputSel["k"]; hasInputKey {
		t.Error("InputPath $$ replaces the input with the context object, the input key must be gone")
	}

	value, found, resErr := e.resolvePayloadTemplateValue(execCtx, "", "$$", nil)
	if resErr != nil || !found {
		t.Fatalf("a payload-template value of $$ must resolve (found=%v): %v", found, resErr)
	}
	if got := value.(map[string]interface{})["Execution"].(map[string]interface{})["Name"]; got != "e1" {
		t.Errorf("payload template $$ Execution.Name = %v, want e1", got)
	}

	arg, okArg, argErr := e.resolveIntrinsicPathArg(execCtx, "", "$$", nil)
	if argErr != nil || !okArg {
		t.Fatalf("an intrinsic argument of $$ must resolve (ok=%v): %v", okArg, argErr)
	}
	if got := arg.(map[string]interface{})["State"].(map[string]interface{})["Name"]; got != "S" {
		t.Errorf("intrinsic argument $$ State.Name = %v, want S", got)
	}
}

// TestMapItemSelectorContextNodes pins the documented Map ItemSelector
// example: $$.Map.Item.Index / Value / Source deliver the iteration index,
// the item and its provenance to every iteration's payload.
func TestMapItemSelectorContextNodes(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:ctxmap1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "ctxmap1",
		Status:          "RUNNING",
		Input:           `[{"who":"bob"},{"who":"meg"}]`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type": "Map",
				"ItemSelector": map[string]interface{}{
					"ContextIndex.$":  "$$.Map.Item.Index",
					"ContextValue.$":  "$$.Map.Item.Value",
					"ContextSource.$": "$$.Map.Item.Source",
				},
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
		t.Fatalf("persist execution failed: %v", err)
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	eventId := int64(0)
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: exec.Input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	output, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}

	var iterations []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &iterations); err != nil {
		t.Fatalf("output not valid JSON: %v (%s)", err, output)
	}
	if len(iterations) != 2 {
		t.Fatalf("iterations = %d, want 2 (%s)", len(iterations), output)
	}
	for i, who := range []string{"bob", "meg"} {
		if got := iterations[i]["ContextIndex"]; fmt.Sprint(got) != fmt.Sprint(i) {
			t.Errorf("iteration %d ContextIndex = %v, want %d", i, got, i)
		}
		value, _ := iterations[i]["ContextValue"].(map[string]interface{})
		if value == nil || value["who"] != who {
			t.Errorf("iteration %d ContextValue = %v, want the item {who:%s}", i, iterations[i]["ContextValue"], who)
		}
		if iterations[i]["ContextSource"] != "STATE_DATA" {
			t.Errorf("iteration %d ContextSource = %v, want STATE_DATA", i, iterations[i]["ContextSource"])
		}
	}
}
