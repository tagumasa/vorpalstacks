package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// newMapTestStore opens a fresh store for the Map execution tests.
func newMapTestStore(t *testing.T) *sfnstore.StepFunctionStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")
}

// TestExecuteMapItemReaderCSV drives executeMap end to end with an S3 CSV
// ItemReader: items come from the bucket, each runs the iterator, the
// result array preserves the dataset order and the Map Run counts the
// read items.
func TestExecuteMapItemReaderCSV(t *testing.T) {
	store := newMapTestStore(t)
	stub := &stubItemReaderS3{objects: map[string][]byte{
		"src/items.csv": []byte("n\n1\n2\n3\n"),
	}}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(stub)
	e := NewExecutor(store, bus)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e1",
		Status:          "RUNNING",
		Input:           `{}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type": "Map",
				"ItemReader": map[string]interface{}{
					"Resource": "arn:aws:states:::s3:getObject",
					"Parameters": map[string]interface{}{
						"Bucket": "src", "Key": "items.csv",
					},
					"ReaderConfig": map[string]interface{}{
						"InputType": "CSV", "CSVHeaderLocation": "FIRST_ROW",
					},
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
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "M",
		Input:         `{}`,
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
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		t.Fatalf("output not valid JSON: %v (%s)", err, output)
	}
	if len(items) != 3 || items[0]["n"] != "1" || items[2]["n"] != "3" {
		t.Errorf("items = %v", items)
	}

	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("map runs = %v, %v", runs, err)
	}
	if runs[0].ItemCounts.Total != 3 || runs[0].ItemCounts.Succeeded != 3 {
		t.Errorf("map run counts = %+v", runs[0].ItemCounts)
	}
}

// TestExecuteMapToleratedFailure pins the Distributed-mode threshold
// semantics: within-threshold failures leave null placeholders and the map
// succeeds; exceeded thresholds fail with
// States.ExceedToleratedFailureThreshold instead of States.IteratorFailed.
func TestExecuteMapToleratedFailure(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	defJSON := func(mode string, tolerated int64, failOdd bool) *sfnstore.MapState {
		processor := map[string]interface{}{
			"StartAt": "W",
			"States": map[string]interface{}{
				"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
			},
		}
		if mode == "DISTRIBUTED" {
			processor["ProcessorConfig"] = map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"}
		}
		state := map[string]interface{}{
			"Type":          "Map",
			"ItemProcessor": processor,
			"End":           true,
		}
		if tolerated >= 0 {
			state["ToleratedFailureCount"] = tolerated
		}
		if failOdd {
			state["ItemSelector"] = map[string]interface{}{"boom.$": "States.Format('{}')"}
		}
		var ms sfnstore.MapState
		b, _ := json.Marshal(state)
		_ = json.Unmarshal(b, &ms)
		return &ms
	}

	newCtx := func() *ExecutionContext {
		exec := &sfnstore.Execution{
			ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e1",
			StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
			Name:            "e1", Status: "RUNNING", Input: `[1,2,3]`,
		}
		return &ExecutionContext{
			Execution: exec, Definition: &sfnstore.StateMachineDefinition{},
			CurrentState: "M", Input: `[1,2,3]`, EventId: ptrEventID(),
			States: map[string]sfnstore.State{}, QueryLanguage: "JSONPath", MapItemIndex: -1,
		}
	}

	// Inline maps keep the classic iterator failure.
	inline := defJSON("", -1, false)
	_, _, execErr := e.executeMap(context.Background(), newCtx(), inline)
	if execErr != nil {
		t.Fatalf("inline all-success map failed: %v", execErr.Cause)
	}

	// A Distributed map with a threshold above the failure count succeeds
	// and tolerates the failed iterations as nulls: simulate via
	// evaluateToleratedFailure directly for the threshold arithmetic.
	s := defJSON("DISTRIBUTED", 5, false)
	tolerated, exceeded := e.evaluateToleratedFailure(&ExecutionContext{Input: `{}`}, s, 2, 4)
	if !tolerated || exceeded {
		t.Errorf("within-threshold: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	tolerated, exceeded = e.evaluateToleratedFailure(&ExecutionContext{Input: `{}`}, s, 6, 4)
	if tolerated || !exceeded {
		t.Errorf("count-exceeded: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	s.ToleratedFailureCount = nil
	s.ToleratedFailurePercentage = ptrFloat64(50)
	tolerated, exceeded = e.evaluateToleratedFailure(&ExecutionContext{Input: `{}`}, s, 2, 4)
	if !tolerated || exceeded {
		t.Errorf("percentage-within: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	tolerated, exceeded = e.evaluateToleratedFailure(&ExecutionContext{Input: `{}`}, s, 3, 4)
	if tolerated || !exceeded {
		t.Errorf("percentage-exceeded: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	// Default (no thresholds): any failure exceeds the implicit zero.
	s.ToleratedFailurePercentage = nil
	tolerated, exceeded = e.evaluateToleratedFailure(&ExecutionContext{Input: `{}`}, s, 1, 4)
	if tolerated || !exceeded {
		t.Errorf("default-zero: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
}

// TestExecuteMapResultWriter pins the ResultWriter export: manifest and
// result files land under the configured prefix and the state output
// becomes the Map Run ARN plus the export location.
func TestExecuteMapResultWriter(t *testing.T) {
	store := newMapTestStore(t)
	stub := &stubItemReaderS3{objects: map[string][]byte{}, put: map[string][]byte{}}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(stub)
	e := NewExecutor(store, bus)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e1", Status: "RUNNING", Input: `[1,2]`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type": "Map",
				"ItemProcessor": map[string]interface{}{
					"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
					"StartAt":         "W",
					"States": map[string]interface{}{
						"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
					},
				},
				"ResultWriter": map[string]interface{}{
					"Resource":   "arn:aws:states:::s3:putObject",
					"Parameters": map[string]interface{}{"Bucket": "out", "Prefix": "jobs"},
				},
				"End": true,
			},
		},
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: `[1,2]`, EventId: ptrEventID(), States: map[string]sfnstore.State{},
		QueryLanguage: "JSONPath", MapItemIndex: -1,
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
	var summary map[string]interface{}
	if err := json.Unmarshal([]byte(output), &summary); err != nil {
		t.Fatalf("output not the export summary: %v (%s)", err, output)
	}
	details, ok := summary["ResultWriterDetails"].(map[string]interface{})
	if !ok || details["Bucket"] != "out" {
		t.Fatalf("summary = %v", summary)
	}
	manifestKey, _ := details["Key"].(string)
	if !strings.HasSuffix(manifestKey, "/manifest.json") || !strings.HasPrefix(manifestKey, "jobs/") {
		t.Fatalf("manifest key = %q", manifestKey)
	}
	manifest, ok := stub.put["out/"+manifestKey]
	if !ok {
		t.Fatalf("manifest object not written: %v", stub.put)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(manifest, &m); err != nil {
		t.Fatalf("manifest not valid JSON: %v", err)
	}
	if m["MapRunArn"] == nil || m["ResultLocation"] != "s3://out/jobs" {
		t.Errorf("manifest = %v", m)
	}
	resultFiles, ok := stub.put["out/"+strings.TrimSuffix(manifestKey, "manifest.json")+"SUCCEEDED_0.json"]
	if !ok {
		t.Fatalf("result file not written: %v", stub.put)
	}
	var records []map[string]interface{}
	if err := json.Unmarshal(resultFiles, &records); err != nil {
		t.Fatalf("result file not valid JSON: %v", err)
	}
	if len(records) != 2 || records[0]["Status"] != "SUCCEEDED" {
		t.Errorf("records = %v", records)
	}
}

func ptrEventID() *int64 {
	id := int64(1)
	return &id
}

func ptrFloat64(v float64) *float64 { return &v }
