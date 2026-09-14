package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// runJSONPathStates drives executeStates over a raw definition and returns
// the context (whose Output and VariableScope carry the result).
func runJSONPathStates(t *testing.T, def *sfnstore.StateMachineDefinition, input string) *ExecutionContext {
	t.Helper()
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:vars-" + def.StartAt,
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "vars-" + def.StartAt,
		Status:          "RUNNING",
		Input:           input,
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
		Execution: exec, Definition: def, CurrentState: def.StartAt,
		Input: input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
		VariableScope: NewVariableScope(nil),
	}
	if execErr := e.executeStates(context.Background(), execCtx); execErr != nil {
		t.Fatalf("executeStates failed: %v", execErr)
	}
	return execCtx
}

// TestJSONPathAssignAndReference pins the documented JSONPath variables
// round trip: "Assign": {"x.$": "$.v"} stores the variable, and a LATER
// state references it as $x ("newly assigned values will be available in
// the next state").
func TestJSONPathAssignAndReference(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Store",
		States: map[string]interface{}{
			"Store": map[string]interface{}{
				"Type":   "Pass",
				"Assign": map[string]interface{}{"x.$": "$.v"},
				"Next":   "Use",
			},
			"Use": map[string]interface{}{
				"Type":       "Pass",
				"Parameters": map[string]interface{}{"out.$": "$x"},
				"End":        true,
			},
		},
	}
	execCtx := runJSONPathStates(t, def, `{"v":42}`)
	if execCtx.Output != `{"out":42}` {
		t.Fatalf("output = %s, want {\"out\":42}", execCtx.Output)
	}
}

// TestJSONPathVariableSubpathAndIntrinsic pins "$x.field" references and
// variables inside intrinsic arguments: "States.Format('The order number
// is {}', $order.number)".
func TestJSONPathVariableSubpathAndIntrinsic(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Store",
		States: map[string]interface{}{
			"Store": map[string]interface{}{
				"Type":   "Pass",
				"Assign": map[string]interface{}{"order.$": "$.order"},
				"Next":   "Use",
			},
			"Use": map[string]interface{}{
				"Type": "Pass",
				"Parameters": map[string]interface{}{
					"sub.$": "$order.number",
					"msg.$": "States.Format('The order number is {}', $order.number)",
				},
				"End": true,
			},
		},
	}
	execCtx := runJSONPathStates(t, def, `{"order":{"number":"o-77"}}`)
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(execCtx.Output), &out); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, execCtx.Output)
	}
	if out["sub"] != "o-77" {
		t.Errorf("subpath reference = %v, want o-77", out["sub"])
	}
	if out["msg"] != "The order number is o-77" {
		t.Errorf("intrinsic with variable = %v, want the formatted string", out["msg"])
	}
}

// TestJSONPathChoiceVariableAndRuleAssign pins Choice states: the rule's
// Variable may reference a workflow variable, a *Path operand may compare
// variable against variable, and the matched rule's Assign stages values
// for the next state.
func TestJSONPathChoiceVariableAndRuleAssign(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Store",
		States: map[string]interface{}{
			"Store": map[string]interface{}{
				"Type":   "Pass",
				"Assign": map[string]interface{}{"want.$": "$.want", "limit.$": "$.limit"},
				"Next":   "Pick",
			},
			"Pick": map[string]interface{}{
				"Type": "Choice",
				"Choices": []interface{}{
					map[string]interface{}{
						"Variable":         "$want",
						"StringEqualsPath": "$limit",
						"Assign":           map[string]interface{}{"picked.$": "$.tag"},
						"Next":             "Report",
					},
				},
				"Default": "Report",
			},
			"Report": map[string]interface{}{
				"Type":       "Pass",
				"Parameters": map[string]interface{}{"picked.$": "$picked"},
				"End":        true,
			},
		},
	}
	execCtx := runJSONPathStates(t, def, `{"want":"go","limit":"go","tag":"chosen"}`)
	if execCtx.Output != `{"picked":"chosen"}` {
		t.Fatalf("output = %s, want {\"picked\":\"chosen\"} (variable comparison + rule Assign)", execCtx.Output)
	}
}

// TestJSONPathCatchAssign pins the Catch handler's Assign: the error
// output feeds the payload template and the variables become visible in
// the fallback state.
func TestJSONPathCatchAssign(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Boom",
		States: map[string]interface{}{
			"Boom": map[string]interface{}{
				"Type":     "Task",
				"Resource": "arn:aws:states:::batch:submitJob",
				"Catch": []interface{}{
					map[string]interface{}{
						"ErrorEquals": []interface{}{"States.ALL"},
						"Assign":      map[string]interface{}{"errName.$": "$.Error"},
						"Next":        "Report",
					},
				},
			},
			"Report": map[string]interface{}{
				"Type":       "Pass",
				"Parameters": map[string]interface{}{"caught.$": "$errName"},
				"End":        true,
			},
		},
	}
	execCtx := runJSONPathStates(t, def, `{}`)
	if execCtx.Output != `{"caught":"States.TaskFailed"}` {
		t.Fatalf("output = %s, want {\"caught\":\"States.TaskFailed\"} (the catch Assign captured the error output)", execCtx.Output)
	}
}

// TestUndefinedVariableFails pins that referencing an unassigned variable
// fails the evaluation: "If $x had not been previously assigned, the
// example would fail because $x would be undefined."
func TestUndefinedVariableFails(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "Use",
		States: map[string]interface{}{
			"Use": map[string]interface{}{
				"Type":       "Pass",
				"Parameters": map[string]interface{}{"out.$": "$never"},
				"End":        true,
			},
		},
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:vars-undef",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "vars-undef",
		Status:          "RUNNING",
		Input:           `{}`,
	}
	eventId := int64(0)
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "Use",
		Input: exec.Input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
		VariableScope: NewVariableScope(nil),
	}
	execErr := e.executeStates(context.Background(), execCtx)
	if execErr == nil {
		t.Fatal("an undefined variable reference must fail the state")
	}
	if !strings.Contains(execErr.Error(), "never") {
		t.Errorf("error %q must name the undefined variable", execErr.Error())
	}
}

// TestVariableRefScanSeesUnicodeNames pins that the reference scanner uses
// the same name grammar as ValidateVariableName: a $変数 reference in a
// JSONata template is visible to variableReferences, not just to the
// validator.
func TestVariableRefScanSeesUnicodeNames(t *testing.T) {
	if err := ValidateVariableName("変数"); err != nil {
		t.Fatalf("a CJK variable name is legal: %v", err)
	}
	def := `{"QueryLanguage":"JSONata","StartAt":"P","States":{"P":{"Type":"Pass","Output":"{% $変数 %}","End":true}}}`
	refs := extractVariableReferences(def)
	if got := refs["P"]; len(got) != 1 || got[0] != "変数" {
		t.Fatalf("state P references = %v, want [変数]", got)
	}
}

// TestVariableBudgetSharedAcrossScopes pins the per-execution aggregate:
// "The total size of all stored variables cannot exceed 10MiB per
// execution" — sibling scopes (Parallel branches, Map iterations) cannot
// each draw the full budget, and a parent's late assignment counts
// against children created earlier.
func TestVariableBudgetSharedAcrossScopes(t *testing.T) {
	parent := NewVariableScope(nil)
	childA := parent.NewChild()
	childB := parent.NewChild()

	// Each Assign stays under the per-Assign ceiling; the accumulated
	// total is what the shared budget must bound.
	chunk := strings.Repeat("a", 200*1024)
	assign := func(i int) map[string]interface{} {
		return map[string]interface{}{fmt.Sprintf("v%d", i): chunk}
	}
	for i := 0; i < 30; i++ {
		if err := childA.SetAll(assign(i)); err != nil {
			t.Fatalf("child A assignment %d rejected: %v", i, err)
		}
	}
	// child A alone holds ~6MB; the sibling must now hit the shared
	// execution ceiling, not a fresh one.
	rejected := false
	for i := 0; i < 30; i++ {
		if err := childB.SetAll(assign(i)); err != nil {
			rejected = true
			break
		}
	}
	if !rejected {
		t.Error("the sibling scope exhausted 30 further assignments without ever hitting the per-execution total")
	}
	// The parent is on the same budget as its children.
	if err := parent.SetAll(map[string]interface{}{"late": strings.Repeat("b", 200*1024)}); err == nil {
		t.Error("the parent's late assignment bypassed the exhausted execution budget")
	}
}

// TestVariableScopeReleaseReturnsBudgetBytes pins the live-storage reading
// of the 10MiB execution total at scope level: a completed child scope's
// variables are out of scope, so Release returns their bytes — and a
// released scope's lookups fall through to the parent exactly as the
// out-of-scope wording describes.
func TestVariableScopeReleaseReturnsBudgetBytes(t *testing.T) {
	root := NewVariableScope(nil)
	if err := root.SetAll(map[string]interface{}{"shared": 1}); err != nil {
		t.Fatalf("root assign: %v", err)
	}
	child := root.NewChild()

	chunk := strings.Repeat("a", 250*1024)
	assign := func(i int) error {
		return child.SetAll(map[string]interface{}{fmt.Sprintf("v%d", i): chunk})
	}
	// Forty dead-scope assignments of ~250KiB approach but do not reach the
	// 10MiB total; the next one exhausts it.
	for i := 0; i < 40; i++ {
		if err := assign(i); err != nil {
			t.Fatalf("fill assignment %d: %v", i, err)
		}
	}
	if err := assign(40); err == nil {
		t.Fatal("the execution total must bound the sum while the scope holds its variables")
	}

	child.Release()
	if err := assign(40); err != nil {
		t.Fatalf("a released scope's bytes must return to the budget: %v", err)
	}
	if _, found := child.Get("v0"); found {
		t.Error("a released scope's variables are no longer accessible")
	}
	if v, found := child.Get("shared"); !found || v != 1 {
		t.Error("a released scope's lookups fall through to the parent")
	}
	child.Release() // idempotent: a second release subtracts nothing
	if err := assign(41); err != nil {
		t.Fatalf("double release must not corrupt the budget: %v", err)
	}
}

// TestMapIterationsReleaseVariableBudget pins the execution-wide variable
// budget against live storage end to end: inline Map iterations whose
// cumulative assignments exceed the 10MiB total — but never hold more
// than one iteration's worth live, because MaxConcurrency 1 makes that
// bound structural rather than scheduler luck ("Step Functions doesn't
// start a new iteration until it completes the previous iteration") —
// succeed, since a completed iteration's variables are out of scope and
// no longer count against "the total size of all stored variables".
func TestMapIterationsReleaseVariableBudget(t *testing.T) {
	svc, store := newRecoveryService(t)

	payload := strings.Repeat("a", 250*1024)
	items := make([]interface{}, 42)
	for i := range items {
		items[i] = i
	}
	inputBytes, err := json.Marshal(map[string]interface{}{"items": items, "payload": payload})
	if err != nil {
		t.Fatal(err)
	}

	definition := `{
		"StartAt": "Seed",
		"States": {
			"Seed": {"Type": "Pass", "Assign": {"big.$": "$.payload"}, "Next": "M"},
			"M": {
				"Type": "Map",
				"ItemsPath": "$.items",
				"MaxConcurrency": 1,
				"ItemProcessor": {
					"StartAt": "Copy",
					"States": {"Copy": {"Type": "Pass", "Assign": {"local.$": "$big"}, "Result": "ok", "End": true}}
				},
				"End": true
			}
		}
	}`
	sm := &sfnstore.StateMachine{Name: "budget-map", Definition: definition}
	if err := store.CreateStateMachine(t.Context(), sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "exec-budget-map", string(inputBytes), "")
	exec.ExecutionArn = sm.StateMachineArn + ":exec-budget-map"
	if err := store.CreateExecution(t.Context(), exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}
	svc.launchExecution(store, exec)

	final := waitForExecution(t, store, exec.ExecutionArn)
	if final.Status != "SUCCEEDED" {
		t.Fatalf("status = %s (error %q, cause %q) — completed iterations must not keep consuming the 10MiB total", final.Status, final.Error, final.Cause)
	}
}

// TestVariableRefScanFollowsUAX31Grammar pins the scanner's parity with
// the UAX #31 name grammar: a mark-carrying or letter-number name scans
// like the validator accepts it, an underscore continues a name, and a
// leading underscore never begins one — so "$_hidden" carries no variable
// reference at all.
func TestVariableRefScanFollowsUAX31Grammar(t *testing.T) {
	marks := `{"QueryLanguage":"JSONata","StartAt":"P","States":{"P":{"Type":"Pass","Output":"{% $cafe\u0301 %}","End":true}}}`
	refs := extractVariableReferences(marks)
	if got := refs["P"]; len(got) != 1 || got[0] != "cafe\u0301" {
		t.Fatalf("state P references = %q, want the mark-carrying name", got)
	}

	numeral := `{"QueryLanguage":"JSONata","StartAt":"P","States":{"P":{"Type":"Pass","Output":"{% $Ⅰcount %}","End":true}}}`
	refs = extractVariableReferences(numeral)
	if got := refs["P"]; len(got) != 1 || got[0] != "Ⅰcount" {
		t.Fatalf("state P references = %q, want the letter-number-starting name", got)
	}

	continues := `{"QueryLanguage":"JSONata","StartAt":"P","States":{"P":{"Type":"Pass","Output":"{% $my_var + 1 %}","End":true}}}`
	refs = extractVariableReferences(continues)
	if got := refs["P"]; len(got) != 1 || got[0] != "my_var" {
		t.Fatalf("state P references = %q, want [my_var]", got)
	}

	leading := `{"QueryLanguage":"JSONata","StartAt":"P","States":{"P":{"Type":"Pass","Output":"{% $_hidden %}","End":true}}}`
	refs = extractVariableReferences(leading)
	if got := refs["P"]; len(got) != 0 {
		t.Fatalf("state P references = %q, want none — a leading underscore never begins a name", got)
	}
}
