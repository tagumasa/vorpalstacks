package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestMapThresholdsJSONataExpressionsAccepted pins the creation-time side
// of the Distributed-mode field contracts: "In JSONata states, you can
// specify a JSONata expression that evaluates to an integer" (MaxConcurrency,
// ToleratedFailureCount, ToleratedFailurePercentage) and the ItemBatcher
// sentence "For JSONata-based states, you can also provide a JSONata
// expression that evaluates to a positive integer" / "you can provide
// JSONata expressions directly to BatchInput".
func TestMapThresholdsJSONataExpressionsAccepted(t *testing.T) {
	base := func(replacements map[string]string) string {
		members := []string{
			`"Items": [1, 2, 3, 4]`,
			`"MaxConcurrency": "{% $states.input.limit %}"`,
			`"ToleratedFailureCount": "{% $states.input.allowed %}"`,
			`"ToleratedFailurePercentage": "{% $states.input.pct %}"`,
			`"ItemBatcher": {"MaxItemsPerBatch": "{% $states.input.batch %}", "BatchInput": "{% {'factCheck': $states.input.date} %}"}`,
		}
		var kept []string
		for _, m := range members {
			key := m[:strings.Index(m, ":")]
			if _, drop := replacements[key]; !drop {
				kept = append(kept, m)
			}
		}
		return `{"QueryLanguage": "JSONata", "StartAt": "M", "States": {"M": {"Type": "Map", ` + strings.Join(kept, ", ") + `, "ItemProcessor": {"StartAt": "W", "States": {"W": {"Type": "Pass", "End": true}}}, "End": true}}}`
	}
	if err := validateDefinitionStructure(base(nil), "STANDARD"); err != nil {
		t.Fatalf("the documented JSONata expression forms were rejected: %v", err)
	}

	// The same strings are not legal in JSONPath states.
	jsonPath := strings.Replace(base(nil), `"QueryLanguage": "JSONata", `, "", 1)
	if err := validateDefinitionStructure(jsonPath, "STANDARD"); err == nil {
		t.Fatal("expression strings in a JSONPath Map were accepted")
	}
}

// TestResolveMapNumericMember unit-pins the member resolver: literal
// numbers pass through, expression strings evaluate, and non-integer or
// out-of-range results fail.
func TestResolveMapNumericMember(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"
	execCtx := &ExecutionContext{
		Execution:     &sfnstore.Execution{ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:mn"},
		EventId:       ptrEventID(),
		QueryLanguage: "JSONata",
		Input:         `{"limit": 7, "pct": 40}`,
	}
	ctx := context.Background()

	v, present, err := e.resolveMapNumericMember(ctx, execCtx, "MaxConcurrency", float64(5), `{}`, 0, 0)
	if err != nil || !present || v != 5 {
		t.Fatalf("literal: got (%v, %v, %v)", v, present, err)
	}
	v, present, err = e.resolveMapNumericMember(ctx, execCtx, "MaxConcurrency", "{% $states.input.limit %}", execCtx.Input, 0, 0)
	if err != nil || !present || v != 7 {
		t.Fatalf("expression: got (%v, %v, %v)", v, present, err)
	}
	if _, _, err = e.resolveMapNumericMember(ctx, execCtx, "ToleratedFailurePercentage", "{% $states.input.limit * 1000 %}", execCtx.Input, 0, 100); err == nil {
		t.Fatal("an out-of-range expression result was accepted for a bounded member")
	}
	if _, _, err = e.resolveMapNumericMember(ctx, execCtx, "MaxConcurrency", "{% $states.input %}", execCtx.Input, 0, 0); err == nil {
		t.Fatal("a non-integer expression result was accepted")
	}
	if _, _, err = e.resolveMapNumericMember(ctx, execCtx, "MaxConcurrency", "not an expression", execCtx.Input, 0, 0); err == nil {
		t.Fatal("a non-expression string was accepted")
	}
}

// TestToleratedFailurePathErrorsSurface pins the fail-closed threshold
// resolution: a ToleratedFailure*Path that selects nothing or selects a
// non-number, and a literal outside its range, are execution errors —
// never a silent fallback to the unconfigured fail-fast default, which
// would invert the operator's configured tolerance.
func TestToleratedFailurePathErrorsSurface(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"
	execCtx := &ExecutionContext{
		Execution: &sfnstore.Execution{ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:tf"},
		EventId:   ptrEventID(),
		Input:     `{"tolerated": "5"}`,
	}
	ctx := context.Background()

	state := &sfnstore.MapState{ToleratedFailureCountPath: "$.tolerated"}
	if _, err := e.resolveToleratedFailureThresholds(ctx, execCtx, state, execCtx.Input); err == nil || err.ErrorCode != "States.InvalidInput" {
		t.Fatalf("non-number threshold path: got %v, want States.InvalidInput", err)
	}
	state = &sfnstore.MapState{ToleratedFailurePercentagePath: "$.missing"}
	if _, err := e.resolveToleratedFailureThresholds(ctx, execCtx, state, execCtx.Input); err == nil || err.ErrorCode != "States.InvalidInput" {
		t.Fatalf("threshold path selecting nothing: got %v, want States.InvalidInput", err)
	}
	state = &sfnstore.MapState{ToleratedFailureCountPath: "$.tolerated"}
	if _, err := e.resolveToleratedFailureThresholds(ctx, execCtx, state, `[1,2]`); err == nil || err.ErrorCode != "States.InvalidInput" {
		t.Fatalf("threshold path against non-object input: got %v, want States.InvalidInput", err)
	}
	state = &sfnstore.MapState{ToleratedFailurePercentage: float64(150)}
	if _, err := e.resolveToleratedFailureThresholds(ctx, execCtx, state, execCtx.Input); err == nil || err.ErrorCode != "States.InvalidInput" {
		t.Fatalf("out-of-range literal: got %v, want States.InvalidInput", err)
	}
	state = &sfnstore.MapState{ToleratedFailureCountPath: "$.tolerated"}
	thresholds, err := e.resolveToleratedFailureThresholds(ctx, execCtx, state, `{"tolerated": 3}`)
	if err != nil || !thresholds.hasCount || thresholds.count != 3 {
		t.Fatalf("valid threshold path: got (%+v, %v)", thresholds, err)
	}
}

// TestMapBatchInputExpressionResolved pins the runtime side: a JSONata
// ItemBatcher BatchInput expression is resolved and merged into each
// batch unit.
func TestMapBatchInputExpressionResolved(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"
	execCtx := &ExecutionContext{
		Execution:     &sfnstore.Execution{ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:bi"},
		EventId:       ptrEventID(),
		QueryLanguage: "JSONata",
		Input:         `{"date": "December 2022", "batch": 2}`,
		VariableScope: NewVariableScope(nil),
	}
	state := &sfnstore.MapState{
		ItemBatcher: &sfnstore.ItemBatcherConfig{
			MaxItemsPerBatch: "{% $states.input.batch %}",
			BatchInput:       "{% {'factCheck': $states.input.date} %}",
		},
		QueryLanguage: "JSONata",
	}
	items := []interface{}{float64(1), float64(2), float64(3)}
	units, err := e.buildMapWorkUnits(context.Background(), execCtx, state, execCtx.Input, items, items)
	if err != nil {
		t.Fatalf("batching with expression BatchInput failed: %v", err)
	}
	if len(units) != 2 {
		t.Fatalf("produced %d units, want 2 (three items in batches of two)", len(units))
	}
	var first map[string]interface{}
	if jerr := json.Unmarshal([]byte(units[0].InputJSON), &first); jerr != nil {
		t.Fatalf("first unit input %s is not JSON: %v", units[0].InputJSON, jerr)
	}
	nested, _ := first["BatchInput"].(map[string]interface{})
	if nested["factCheck"] != "December 2022" {
		t.Fatalf("BatchInput = %v, want the resolved expression object", first["BatchInput"])
	}
	if arr, ok := first["Items"].([]interface{}); !ok || len(arr) != 2 {
		t.Fatalf("Items = %v, want the first two items", first["Items"])
	}
}

// TestMapMaxConcurrencyExpressionRuns pins that an inline JSONata Map
// carrying the expression form runs end to end.
func TestMapMaxConcurrencyExpressionRuns(t *testing.T) {
	state := map[string]interface{}{
		"Type":           "Map",
		"MaxConcurrency": "{% $states.input.limit %}",
		"Items":          "{% $states.input.items %}",
		"ItemProcessor":  mapEchoProcessor(),
		"End":            true,
	}
	items := runMapItems(t, mapDef(state), `{"limit": 2, "items": [1, 2, 3]}`, true)
	if len(items) != 3 {
		t.Fatalf("produced %d items, want 3", len(items))
	}
}

// TestValidatorDiagnosticsDeterministic pins that the unknown-field
// diagnostics are emitted in sorted member order: two creations of the
// same definition produce the identical first error.
func TestValidatorDiagnosticsDeterministic(t *testing.T) {
	def := `{"StartAt": "A", "States": {"A": {"Type": "Pass", "Zed": 1, "Abc": 2, "Mid": 3, "End": true}}}`
	first := validateDefinitionStructure(def, "STANDARD")
	if first == nil {
		t.Fatal("expected a validation error")
	}
	for i := 0; i < 20; i++ {
		next := validateDefinitionStructure(def, "STANDARD")
		if next == nil || next.Error() != first.Error() {
			t.Fatalf("run %d message diverged: %v vs %v", i, next, first)
		}
	}
}

// TestVariableReferencesSortedAndNested pins the describe-time map: each
// state's variable list is sorted, and references inside Parallel branches
// and Map processors are collected.
func TestVariableReferencesSortedAndNested(t *testing.T) {
	def := `{
		"QueryLanguage": "JSONata",
		"StartAt": "P",
		"States": {
			"P": {
				"Type": "Parallel",
				"Branches": [{
					"StartAt": "Inner",
					"States": {
						"Inner": {"Type": "Pass", "Output": {"a": "{% $zeta %}", "b": "{% $alpha %}", "c": "{% $mid %}"}, "End": true}
					}
				}],
				"End": true
			}
		}
	}`
	refs := extractVariableReferences(def)
	inner, ok := refs["Inner"]
	if !ok {
		t.Fatalf("branch state references missing: %v", refs)
	}
	want := []string{"alpha", "mid", "zeta"}
	for i, name := range want {
		if inner[i] != name {
			t.Fatalf("refs = %v, want sorted %v", inner, want)
		}
	}

	mapDef := `{
		"QueryLanguage": "JSONata",
		"StartAt": "M",
		"States": {
			"M": {
				"Type": "Map",
				"Items": [1],
				"ItemProcessor": {
					"StartAt": "W",
					"States": {"W": {"Type": "Pass", "Output": {"v": "{% $counter %}"}, "End": true}}
				},
				"End": true
			}
		}
	}`
	refs = extractVariableReferences(mapDef)
	if got := refs["W"]; len(got) != 1 || got[0] != "counter" {
		t.Fatalf("ItemProcessor state references = %v, want [counter]", got)
	}
}

// TestQueryLanguageSeparationTables pins the three separation contracts:
// the Map *Path fields and ItemSelector are JSONPath-only, Output is
// JSONata-only in JSONPath states, the state-level Input member is not a
// documented field on either dialect, and a JSONata machine cannot revert
// a state to JSONPath.
func TestQueryLanguageSeparationTables(t *testing.T) {
	jsonataMapPaths := `{"QueryLanguage": "JSONata", "StartAt": "M", "States": {"M": {"Type": "Map", "MaxConcurrencyPath": "$.m", "Items": [1], "ItemProcessor": {"StartAt": "W", "States": {"W": {"Type": "Pass", "End": true}}}, "End": true}}}`
	if err := validateDefinitionStructure(jsonataMapPaths, "STANDARD"); err == nil {
		t.Fatal("MaxConcurrencyPath in a JSONata Map was accepted")
	}
	jsonataThresholdPath := `{"QueryLanguage": "JSONata", "StartAt": "M", "States": {"M": {"Type": "Map", "ToleratedFailureCountPath": "$.t", "Items": [1], "ItemProcessor": {"StartAt": "W", "States": {"W": {"Type": "Pass", "End": true}}}, "End": true}}}`
	if err := validateDefinitionStructure(jsonataThresholdPath, "STANDARD"); err == nil {
		t.Fatal("ToleratedFailureCountPath in a JSONata Map was accepted")
	}
	jsonataItemSelector := `{"QueryLanguage": "JSONata", "StartAt": "M", "States": {"M": {"Type": "Map", "ItemSelector": {"x": 1}, "Items": [1], "ItemProcessor": {"StartAt": "W", "States": {"W": {"Type": "Pass", "End": true}}}, "End": true}}}`
	if err := validateDefinitionStructure(jsonataItemSelector, "STANDARD"); err == nil {
		t.Fatal("ItemSelector in a JSONata Map was accepted")
	}

	jsonPathOutput := `{"StartAt": "A", "States": {"A": {"Type": "Pass", "Output": {"y": 2}, "End": true}}}`
	if err := validateDefinitionStructure(jsonPathOutput, "STANDARD"); err == nil {
		t.Fatal("Output in a JSONPath Pass was accepted")
	}
	jsonPathInput := `{"StartAt": "A", "States": {"A": {"Type": "Pass", "Input": {"y": 2}, "End": true}}}`
	if err := validateDefinitionStructure(jsonPathInput, "STANDARD"); err == nil {
		t.Fatal("Input in a JSONPath Pass was accepted")
	}
	jsonataInput := `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Pass", "Input": {"y": 2}, "End": true}}}`
	if err := validateDefinitionStructure(jsonataInput, "STANDARD"); err == nil {
		t.Fatal("a state-level Input member has no documented basis on either dialect and must be rejected")
	}

	jsonataMachineReverting := `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Pass", "QueryLanguage": "JSONPath", "InputPath": "$", "End": true}}}`
	if err := validateDefinitionStructure(jsonataMachineReverting, "STANDARD"); err == nil {
		t.Fatal("a JSONPath state override inside a JSONata machine was accepted")
	}
	jsonPathMachineUpgrading := `{"StartAt": "A", "States": {"A": {"Type": "Pass", "QueryLanguage": "JSONata", "Output": {"y": "{% $states.input.x %}"}, "End": true}}}`
	if err := validateDefinitionStructure(jsonPathMachineUpgrading, "STANDARD"); err != nil {
		t.Fatalf("a JSONata state override inside a JSONPath machine was rejected: %v", err)
	}
}
