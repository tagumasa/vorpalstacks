package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
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
	states, err := extractStatesFromDefinition(def)
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
// States.ExceedToleratedFailureThreshold rather than the iteration's own
// error identity.
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
	// and tolerates the failed iterations as nulls: simulate via the
	// threshold resolver and evaluateToleratedFailure directly for the
	// threshold arithmetic.
	s := defJSON("DISTRIBUTED", 5, false)
	thresholds, terr := e.resolveToleratedFailureThresholds(context.Background(), &ExecutionContext{Input: `{}`}, s, `{}`)
	if terr != nil {
		t.Fatalf("threshold resolution failed: %v", terr.Cause)
	}
	tolerated, exceeded := e.evaluateToleratedFailure(s, thresholds, 2, 4)
	if !tolerated || exceeded {
		t.Errorf("within-threshold: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	tolerated, exceeded = e.evaluateToleratedFailure(s, thresholds, 6, 4)
	if tolerated || !exceeded {
		t.Errorf("count-exceeded: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	s.ToleratedFailureCount = nil
	s.ToleratedFailurePercentage = float64(50)
	thresholds, terr = e.resolveToleratedFailureThresholds(context.Background(), &ExecutionContext{Input: `{}`}, s, `{}`)
	if terr != nil {
		t.Fatalf("percentage threshold resolution failed: %v", terr.Cause)
	}
	tolerated, exceeded = e.evaluateToleratedFailure(s, thresholds, 2, 4)
	if !tolerated || exceeded {
		t.Errorf("percentage-within: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	tolerated, exceeded = e.evaluateToleratedFailure(s, thresholds, 3, 4)
	if tolerated || !exceeded {
		t.Errorf("percentage-exceeded: tolerated=%v exceeded=%v", tolerated, exceeded)
	}
	// Default (no thresholds): any failure exceeds the implicit zero.
	s.ToleratedFailurePercentage = nil
	thresholds, _ = e.resolveToleratedFailureThresholds(context.Background(), &ExecutionContext{Input: `{}`}, s, `{}`)
	tolerated, exceeded = e.evaluateToleratedFailure(s, thresholds, 1, 4)
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
	states, err := extractStatesFromDefinition(def)
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

// TestMapItemsPathOnArrayInputRejected pins the ASL contract: ItemsPath
// resolves against object input, so a Map state carrying a non-trivial
// ItemsPath over an array input fails with States.InvalidItemsPath
// instead of silently resolving the path against the first element.
func TestMapItemsPathOnArrayInputRejected(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"ItemsPath": "$.items",
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
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/itemspath-array",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "itemspath-array",
		Status:          "RUNNING",
		Input:           `[{"a":1}]`,
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

	_, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr == nil {
		t.Fatal("ItemsPath over array input was accepted")
	}
	if execErr.ErrorCode != "States.InvalidItemsPath" {
		t.Errorf("error code = %s, want States.InvalidItemsPath", execErr.ErrorCode)
	}
}

// TestExecuteMapResultWriterPreview pins the WriterConfig-only preview
// form: with no Resource or Parameters "the results are passed on to the
// next state" — the formatted child results become the state output and
// nothing is written to S3 (the executor here carries no S3 integration
// at all).
func TestExecuteMapResultWriterPreview(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
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
						"W": map[string]interface{}{"Type": "Pass", "Result": []interface{}{map[string]interface{}{"x": "a"}}, "ResultPath": "$", "End": true},
					},
				},
				"ResultWriter": map[string]interface{}{
					"WriterConfig": map[string]interface{}{"Transformation": "FLATTEN", "OutputType": "JSON"},
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
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	output, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("preview ResultWriter failed: %v", execErr.Cause)
	}
	// FLATTEN splices each child's array result into one array.
	if output != `[{"x":"a"},{"x":"a"}]` {
		t.Fatalf("preview output = %s, want the flattened child results", output)
	}
}

// TestExecuteMapResultWriterReferencePaths pins the reference-path
// Parameters form ("Bucket.$" selecting from the Distributed Map state
// input) together with the WriterConfig transformation of the exported
// result file: the SUCCEEDED file carries the formatted child results,
// not the untransformed execution records.
func TestExecuteMapResultWriterReferencePaths(t *testing.T) {
	store := newMapTestStore(t)
	stub := &stubItemReaderS3{objects: map[string][]byte{}, put: map[string][]byte{}}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(stub)
	e := NewExecutor(store, bus)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e2",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e2", Status: "RUNNING", Input: `{"items":[1,2],"resultBucket":"exports","pfx":"jobs"}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"ItemsPath": "$.items",
				"ItemProcessor": map[string]interface{}{
					"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
					"StartAt":         "W",
					"States": map[string]interface{}{
						"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
					},
				},
				"ResultWriter": map[string]interface{}{
					"Resource":     "arn:aws:states:::s3:putObject",
					"Parameters":   map[string]interface{}{"Bucket.$": "$.resultBucket", "Prefix.$": "$.pfx"},
					"WriterConfig": map[string]interface{}{"Transformation": "FLATTEN", "OutputType": "JSON"},
				},
				"End": true,
			},
		},
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: exec.Input, EventId: ptrEventID(), States: map[string]sfnstore.State{},
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	output, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("reference-path ResultWriter failed: %v", execErr.Cause)
	}
	var summary map[string]interface{}
	if err := json.Unmarshal([]byte(output), &summary); err != nil {
		t.Fatalf("output not the export summary: %v (%s)", err, output)
	}
	details, _ := summary["ResultWriterDetails"].(map[string]interface{})
	if details["Bucket"] != "exports" {
		t.Fatalf("export bucket = %v, want the reference-path-resolved exports", details["Bucket"])
	}
	manifestKey, _ := details["Key"].(string)
	if !strings.HasPrefix(manifestKey, "jobs/") {
		t.Fatalf("manifest key = %q, want the resolved jobs/ prefix", manifestKey)
	}
	succeeded, ok := stub.put["exports/"+strings.TrimSuffix(manifestKey, "manifest.json")+"SUCCEEDED_0.json"]
	if !ok {
		t.Fatalf("result file not written: %v", stub.put)
	}
	// The FLATTEN transformation formats the exported results: the child
	// outputs, not the per-execution records.
	var payload []interface{}
	if err := json.Unmarshal(succeeded, &payload); err != nil {
		t.Fatalf("result file not valid JSON: %v", err)
	}
	if len(payload) != 2 {
		t.Fatalf("formatted payload = %s, want the two child outputs", succeeded)
	}
	for _, el := range payload {
		if record, isObject := el.(map[string]interface{}); isObject {
			if _, isRecord := record["ExecutionArn"]; isRecord {
				t.Fatalf("formatted payload carries execution records: %s", succeeded)
			}
		}
	}
}

// TestDistributedMapAbortRecordsAbortedStatuses pins the abort bookkeeping:
// a stopped Distributed Map run is ABORTED (MapRunStatus), and a cancelled
// child workflow execution — "the user cancelled the execution" — is
// ABORTED with the platform's abort terminal pair, never FAILED/States
// .Runtime.
func TestDistributedMapAbortRecordsAbortedStatuses(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e3",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e3", Status: "RUNNING", Input: `[1]`,
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
						"W": map[string]interface{}{"Type": "Wait", "Seconds": 30, "End": true},
					},
				},
				"End": true,
			},
		},
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: `[1]`, EventId: ptrEventID(), States: map[string]sfnstore.State{},
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	var mapErr *ExecutionError
	go func() {
		defer close(done)
		_, _, mapErr = e.executeMap(ctx, execCtx, execCtx.States["M"].(*sfnstore.MapState))
	}()

	childArn := ""
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		execs, lerr := store.ListAllExecutions(ctx, exec.StateMachineArn, "", "", "")
		if lerr == nil {
			for _, xc := range execs {
				if xc.ExecutionArn != exec.ExecutionArn && xc.MapRunArn != "" {
					childArn = xc.ExecutionArn
				}
			}
		}
		if childArn != "" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if childArn == "" {
		cancel()
		t.Fatal("child execution never appeared")
	}
	cancel()
	<-done
	if mapErr == nil {
		t.Fatal("the cancelled map must return an error")
	}

	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) == 0 {
		t.Fatalf("map runs: %v (%d)", err, len(runs))
	}
	if last := runs[len(runs)-1]; last.Status != "ABORTED" {
		t.Errorf("map run status = %s, want ABORTED", last.Status)
	}

	child, err := store.GetExecution(context.Background(), childArn)
	if err != nil {
		t.Fatalf("get child: %v", err)
	}
	if child.Status != "ABORTED" {
		t.Errorf("child status = %s, want ABORTED", child.Status)
	}
	if child.Error != "ExecutionAborted" {
		t.Errorf("child error = %q, want the abort vocabulary", child.Error)
	}
	history, _, herr := store.GetExecutionHistory(context.Background(), childArn, 100, "", false)
	if herr != nil {
		t.Fatalf("child history: %v", herr)
	}
	abortedEvent := false
	for _, ev := range history {
		if ev.Type == "ExecutionAborted" {
			abortedEvent = true
		}
		if ev.Type == "ExecutionFailed" {
			t.Errorf("child history carries ExecutionFailed for a cancelled unit")
		}
	}
	if !abortedEvent {
		t.Error("child history lacks the ExecutionAborted terminal event")
	}
}

// TestInlineMapRunStaysInvisible pins the Map Run resource boundary:
// "When you run a Map state in Distributed mode, Step Functions creates a
// Map Run resource" — an inline map keeps only its internal redrive
// checkpoint, which neither lists nor describes as a Map Run.
func TestInlineMapRunStaysInvisible(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e4",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e4", Status: "RUNNING", Input: `[1]`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":          "Map",
				"ItemProcessor": map[string]interface{}{"StartAt": "P", "States": map[string]interface{}{"P": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true}}},
				"End":           true,
			},
		},
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: `[1]`, EventId: ptrEventID(), States: map[string]sfnstore.State{},
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}
	if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr != nil {
		t.Fatalf("inline map failed: %v", execErr.Cause)
	}

	// The engine keeps its checkpoint record…
	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("engine checkpoint records: %v (%d)", err, len(runs))
	}
	if !runs[0].Inline {
		t.Error("the inline map's record must carry the inline marker")
	}

	// …but the Map Run API surface reports nothing.
	svc := &StepFunctionService{}
	resp, err := svc.listMapRunsCore(context.Background(), store, exec.ExecutionArn, 100, "")
	if err != nil {
		t.Fatalf("list map runs: %v", err)
	}
	if listed, _ := resp["mapRuns"].([]map[string]interface{}); len(listed) != 0 {
		t.Errorf("inline map listed %d map runs, want none", len(listed))
	}
	if _, err := svc.describeMapRunCore(context.Background(), store, runs[0].MapRunArn); err == nil {
		t.Error("describing an inline map's checkpoint must not find a Map Run")
	}
	if err := svc.updateMapRunCore(context.Background(), store, UpdateMapRunInput{MapRunArn: runs[0].MapRunArn, MaxConcurrency: ptrInt64(2)}); err == nil {
		t.Error("updating an inline map's checkpoint must not find a Map Run")
	}
}

// TestMapRunRedriveBookkeeping pins the redrive counters: the reclaimed
// run's redrive count is "always greater than 0" with a redriveDate, and
// the MapRunRedrived event carries the incremented count.
func TestMapRunRedriveBookkeeping(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e5",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e5", Status: "RUNNING", Input: `[1]`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type": "Map",
				"ItemProcessor": map[string]interface{}{
					"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
					"StartAt":         "F",
					"States": map[string]interface{}{
						"F": map[string]interface{}{"Type": "Fail", "Error": "IterFail", "Cause": "iteration failed"},
					},
				},
				"End": true,
			},
		},
	}
	newContext := func() *ExecutionContext {
		execCtx := &ExecutionContext{
			Execution: exec, Definition: def, CurrentState: "M",
			Input: `[1]`, EventId: ptrEventID(), States: map[string]sfnstore.State{},
			QueryLanguage: "JSONPath", MapItemIndex: -1,
		}
		states, err := extractStatesFromDefinition(def)
		if err != nil {
			t.Fatalf("extract states failed: %v", err)
		}
		execCtx.States = states
		return execCtx
	}

	first := newContext()
	if _, _, execErr := e.executeMap(context.Background(), first, first.States["M"].(*sfnstore.MapState)); execErr == nil {
		t.Fatal("the failing iteration must fail the map")
	}
	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("map runs after first attempt: %v (%d)", err, len(runs))
	}
	if runs[0].Status != "FAILED" || runs[0].RedriveCount != 0 {
		t.Fatalf("first attempt run = %s redriveCount=%d, want FAILED/0", runs[0].Status, runs[0].RedriveCount)
	}

	second := newContext()
	second.IsRedrive = true
	if _, _, execErr := e.executeMap(context.Background(), second, second.States["M"].(*sfnstore.MapState)); execErr == nil {
		t.Fatal("the redriven attempt still fails its iteration")
	}
	runs, err = store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 1 {
		t.Fatalf("map runs after redrive: %v (%d)", err, len(runs))
	}
	if runs[0].RedriveCount != 1 {
		t.Errorf("redriveCount = %d, want 1 (a redriven run always counts greater than 0)", runs[0].RedriveCount)
	}
	if runs[0].RedriveDate == 0 {
		t.Error("redriveDate was never recorded")
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), exec.ExecutionArn, 1000, "", false)
	if herr != nil {
		t.Fatalf("history: %v", herr)
	}
	redrivedEvent := false
	for _, ev := range history {
		if ev.Type == "MapRunRedriven" && ev.MapRunRedrivenEventDetails != nil {
			redrivedEvent = true
			if ev.MapRunRedrivenEventDetails.RedriveCount != 1 {
				t.Errorf("MapRunRedriven redriveCount = %d, want 1", ev.MapRunRedrivenEventDetails.RedriveCount)
			}
		}
	}
	if !redrivedEvent {
		t.Error("history lacks the MapRunRedrived event")
	}
}

// TestMapRunReclaimSelectsHighestSequence pins the reclaim's ordering
// contract: a state redriven after several run cycles reclaims the run
// with the highest ARN sequence, not the last list match — a
// lexicographic key order sorts mapRun-12 before mapRun-9, so trusting
// the list order would resurrect an older cycle.
func TestMapRunReclaimSelectsHighestSequence(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/e6",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e6", Status: "RUNNING", Input: `[1]`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type": "Map",
				"ItemProcessor": map[string]interface{}{
					"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
					"StartAt":         "F",
					"States": map[string]interface{}{
						"F": map[string]interface{}{"Type": "Fail", "Error": "IterFail", "Cause": "iteration failed"},
					},
				},
				"End": true,
			},
		},
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	base := "arn:aws:states:us-east-1:000000000000:mapRun:sm/M/mapRun-"
	run9 := base + "9-20260901000000"
	run12 := base + "12-20260902000000"
	for _, mr := range []*sfnstore.MapRun{
		{MapRunArn: run9, ExecutionArn: exec.ExecutionArn, Name: "M", Status: "FAILED", StartDate: 100, ItemCounts: sfnstore.MapRunItemCounts{Total: 1}, CompletedResults: map[int]string{}},
		{MapRunArn: run12, ExecutionArn: exec.ExecutionArn, Name: "M", Status: "FAILED", StartDate: 200, ItemCounts: sfnstore.MapRunItemCounts{Total: 1}, CompletedResults: map[int]string{}},
	} {
		if err := store.CreateMapRun(context.Background(), mr); err != nil {
			t.Fatalf("seed map run %s: %v", mr.MapRunArn, err)
		}
	}

	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: `[1]`, EventId: ptrEventID(), States: map[string]sfnstore.State{},
		QueryLanguage: "JSONPath", MapItemIndex: -1, IsRedrive: true,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states
	if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr == nil {
		t.Fatal("the failing iteration must fail the map")
	}

	runs, err := store.ListMapRunsByExecution(context.Background(), exec.ExecutionArn)
	if err != nil || len(runs) != 2 {
		t.Fatalf("map runs after reclaim: %v (%d)", err, len(runs))
	}
	for _, mr := range runs {
		switch mr.MapRunArn {
		case run12:
			if mr.RedriveCount != 1 {
				t.Errorf("seq-12 run redriveCount = %d, want 1 — the highest sequence is the cycle to reclaim", mr.RedriveCount)
			}
		case run9:
			if mr.RedriveCount != 0 {
				t.Errorf("seq-9 run redriveCount = %d, want 0 — an older cycle must stay untouched", mr.RedriveCount)
			}
		}
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), exec.ExecutionArn, 1000, "", false)
	if herr != nil {
		t.Fatalf("history: %v", herr)
	}
	for _, ev := range history {
		if ev.Type != "MapRunRedriven" || ev.MapRunRedrivenEventDetails == nil {
			continue
		}
		if ev.MapRunRedrivenEventDetails.MapRunArn != run12 {
			t.Errorf("MapRunRedriven mapRunArn = %s, want the seq-12 run %s", ev.MapRunRedrivenEventDetails.MapRunArn, run12)
		}
	}
}

// TestMapItemSourceFromStandardReaderForm pins $$Map.Item.Source for an S3
// ItemReader written in the standard reference-path parameter form
// ("Bucket.$": "$.b", "Key.$": "$.k"): "For all the other input types, the
// value will be the Amazon S3 URI. For example: S3://bucket-name/object-key"
// — and the bucket URI for a list read: "S3://bucket-name".
func TestMapItemSourceFromStandardReaderForm(t *testing.T) {
	store := newMapTestStore(t)
	stub := &stubItemReaderS3{objects: map[string][]byte{
		"src/data.json": []byte(`[{"id":1},{"id":2}]`),
	}}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(stub)
	e := NewExecutor(store, bus)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type": "Map",
				"ItemReader": map[string]interface{}{
					"Resource": "arn:aws:states:::s3:getObject",
					"Parameters": map[string]interface{}{
						"Bucket.$": "$.b", "Key.$": "$.k",
					},
					"ReaderConfig": map[string]interface{}{"InputType": "JSON"},
				},
				"ItemSelector": map[string]interface{}{
					"src.$": "$$.Map.Item.Source",
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
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/src-keyed",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "src-keyed", Status: "RUNNING", Input: `{"b":"src","k":"data.json"}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: exec.Input, EventId: ptrEventID(), States: map[string]sfnstore.State{},
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	states, err := extractStatesFromDefinition(def)
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
		t.Fatalf("output not an array: %v (%s)", err, output)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 iterations, got %d (%s)", len(items), output)
	}
	for i, item := range items {
		if item["src"] != "S3://src/data.json" {
			t.Errorf("item %d Source = %v, want S3://src/data.json — the standard Bucket.$/Key.$ form must report the object URI", i, item["src"])
		}
	}
}

// TestMapItemSourceFromListReaderReportsBucket pins the listing variant:
// an s3:listObjectsV2 reader reports the bucket URI alone
// ("S3://bucket-name").
func TestMapItemSourceFromListReaderReportsBucket(t *testing.T) {
	store := newMapTestStore(t)
	stub := &stubItemReaderS3{
		objects: map[string][]byte{
			"jobs/a.json": []byte(`{}`),
		},
		entries: []invokers.S3ObjectEntry{{Key: "jobs/a.json", Size: 2}},
	}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(stub)
	e := NewExecutor(store, bus)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type": "Map",
				"ItemReader": map[string]interface{}{
					"Resource": "arn:aws:states:::s3:listObjectsV2",
					"Parameters": map[string]interface{}{
						"Bucket.$": "$.b", "Prefix.$": "$.p",
					},
				},
				"ItemSelector": map[string]interface{}{
					"src.$": "$$.Map.Item.Source",
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
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/src-list",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "src-list", Status: "RUNNING", Input: `{"b":"src","p":"jobs/"}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: exec.Input, EventId: ptrEventID(), States: map[string]sfnstore.State{},
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	output, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &items); err != nil || len(items) == 0 {
		t.Fatalf("output not a non-empty array: %v (%s)", err, output)
	}
	if items[0]["src"] != "S3://src" {
		t.Errorf("listing Source = %v, want S3://src — a list read reports the bucket URI", items[0]["src"])
	}
}

// TestToleratedFailureThresholdsResolveOnProcessedInput pins the resolution
// stage: the *Path threshold members select from the input after InputPath
// has been applied (the same stage MaxConcurrency resolves at), not from
// the raw execution input.
func TestToleratedFailureThresholdsResolveOnProcessedInput(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)

	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":                      "Map",
				"InputPath":                 "$.cfg",
				"ToleratedFailureCountPath": "$.maxFail",
				"ItemProcessor": map[string]interface{}{
					"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
					"StartAt":         "W",
					"States": map[string]interface{}{
						"W": map[string]interface{}{"Type": "Pass", "End": true},
					},
				},
				"End": true,
			},
		},
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	ms := states["M"].(*sfnstore.MapState)
	if ms.ToleratedFailureCountPath != "$.maxFail" {
		t.Fatalf("threshold path not parsed: %q", ms.ToleratedFailureCountPath)
	}

	// The raw input nests the value under cfg; resolution against the raw
	// input would find nothing, against the processed input ($.cfg) it
	// finds maxFail.
	thresholds, terr := e.resolveToleratedFailureThresholds(context.Background(),
		&ExecutionContext{Input: `{"cfg":{"maxFail":5}}`}, ms, `{"maxFail":5}`)
	if terr != nil {
		t.Fatalf("threshold resolution failed: %v", terr.Cause)
	}
	if !thresholds.hasCount || thresholds.count != 5 {
		t.Errorf("thresholds = %+v, want count 5 resolved from the processed input", thresholds)
	}
}
