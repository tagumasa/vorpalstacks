package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// distributedProcessor renders an ItemProcessor whose Pass state echoes the
// unit input it receives.
func distributedProcessor() map[string]interface{} {
	return map[string]interface{}{
		"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
		"StartAt":         "W",
		"States": map[string]interface{}{
			"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
		},
	}
}

// TestExecuteMapItemBatcher pins the batching data plane: MaxItemsPerBatch
// groups consecutive items into {"Items": [...]} unit inputs, the result
// array carries one element per unit, the item counters count items and
// the execution counters count units.
func TestExecuteMapItemBatcher(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e2",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e2",
		Status:          "RUNNING",
		Input:           `{"v":[1,2,3,4,5]}`,
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

	output, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}
	var units []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &units); err != nil {
		t.Fatalf("output not valid JSON: %v (%s)", err, output)
	}
	if len(units) != 3 {
		t.Fatalf("expected 3 batch outputs, got %d (%s)", len(units), output)
	}
	for i, want := range []int{2, 2, 1} {
		items, ok := units[i]["Items"].([]interface{})
		if !ok || len(items) != want {
			t.Errorf("unit %d Items = %v, want %d items", i, units[i]["Items"], want)
		}
	}
	if _, has := units[0]["BatchInput"]; has {
		t.Errorf("BatchInput key must be absent when not configured: %v", units[0])
	}

	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("map runs = %v, %v", runs, err)
	}
	if runs[0].ItemCounts.Total != 5 || runs[0].ItemCounts.Succeeded != 5 {
		t.Errorf("item counts = %+v", runs[0].ItemCounts)
	}
	if runs[0].ExecutionCounts.Total != 3 || runs[0].ExecutionCounts.Succeeded != 3 {
		t.Errorf("execution counts = %+v", runs[0].ExecutionCounts)
	}

	if execCtx.AfterItemBatcher == nil || !strings.Contains(*execCtx.AfterItemBatcher, `"Items":[1,2]`) {
		t.Errorf("AfterItemBatcher = %v", execCtx.AfterItemBatcher)
	}
}

// TestExecuteMapItemBatcherBatchInput pins the BatchInput merge: every
// unit input carries the fixed object under the BatchInput key next to
// the Items array.
func TestExecuteMapItemBatcherBatchInput(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e3",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e3",
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
					"BatchInput":       map[string]interface{}{"factCheck": "December 2022"},
				},
				"End": true,
			},
		},
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

	output, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}
	var units []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &units); err != nil {
		t.Fatalf("output not valid JSON: %v (%s)", err, output)
	}
	if len(units) != 2 {
		t.Fatalf("expected 2 batch outputs, got %d (%s)", len(units), output)
	}
	for i := range units {
		batchInput, ok := units[i]["BatchInput"].(map[string]interface{})
		if !ok || batchInput["factCheck"] != "December 2022" {
			t.Errorf("unit %d BatchInput = %v", i, units[i]["BatchInput"])
		}
	}
}

// TestBuildMapWorkUnitsByteCap pins the byte-cap batching: units break
// before the rendered {"Items": ...} payload would exceed
// MaxInputBytesPerBatch, and a single item larger than the cap still
// forms its own unit.
func TestBuildMapWorkUnitsByteCap(t *testing.T) {
	e := &Executor{}
	state := &sfnstore.MapState{
		ItemBatcher: &sfnstore.ItemBatcherConfig{
			MaxInputBytesPerBatch: ptrInt64(20),
		},
	}
	// Each item renders as "x" (3 bytes); the wrapper {"Items":[]} is 11
	// bytes and every separator adds one.
	items := []interface{}{"x", "x", "x", "x"}
	units, execErr := e.buildMapWorkUnits(context.Background(), &ExecutionContext{}, state, `{}`, items, items)
	if execErr != nil {
		t.Fatalf("buildMapWorkUnits failed: %v", execErr.Cause)
	}
	if len(units) != 2 {
		t.Fatalf("expected the byte cap to split the items into 2 units, got %d: %+v", len(units), units)
	}
	for _, u := range units {
		if int64(len(u.InputJSON)) > 20 {
			t.Errorf("unit input %q exceeds the 20 byte cap (%d bytes)", u.InputJSON, len(u.InputJSON))
		}
	}

	oversized := &sfnstore.MapState{
		ItemBatcher: &sfnstore.ItemBatcherConfig{
			MaxInputBytesPerBatch: ptrInt64(10),
		},
	}
	big := strings.Repeat("y", 40)
	units, execErr = e.buildMapWorkUnits(context.Background(), &ExecutionContext{}, oversized, `{}`, []interface{}{big, "z"}, []interface{}{big, "z"})
	if execErr != nil {
		t.Fatalf("buildMapWorkUnits failed: %v", execErr.Cause)
	}
	if len(units) != 2 {
		t.Fatalf("expected the oversized item to form its own unit, got %+v", units)
	}
}

// TestBuildMapWorkUnitsPathVariants pins the reference-path forms: the
// per-unit item ceiling and the fixed batch input resolve from the state
// input.
func TestBuildMapWorkUnitsPathVariants(t *testing.T) {
	e := &Executor{}
	state := &sfnstore.MapState{
		ItemBatcher: &sfnstore.ItemBatcherConfig{
			MaxItemsPerBatchPath: "$.max",
			BatchInputPath:       "$.fixed",
		},
	}
	items := []interface{}{1, 2, 3, 4, 5}
	input := `{"max": 2, "fixed": {"k": "v"}}`
	units, execErr := e.buildMapWorkUnits(context.Background(), &ExecutionContext{}, state, input, items, items)
	if execErr != nil {
		t.Fatalf("buildMapWorkUnits failed: %v", execErr.Cause)
	}
	if len(units) != 3 {
		t.Fatalf("expected 3 units from MaxItemsPerBatchPath 2, got %d", len(units))
	}
	if units[0].ItemCount != 2 {
		t.Errorf("first unit ItemCount = %d, want 2", units[0].ItemCount)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(units[0].InputJSON), &payload); err != nil {
		t.Fatalf("unit input not valid JSON: %v (%s)", err, units[0].InputJSON)
	}
	if fixed, ok := payload["BatchInput"].(map[string]interface{}); !ok || fixed["k"] != "v" {
		t.Errorf("BatchInput = %v", payload["BatchInput"])
	}
}
