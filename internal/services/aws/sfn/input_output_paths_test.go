package sfn

import (
	"context"
	"encoding/json"
	"fmt"
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
	states, err := extractStatesFromDefinition(def)
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
	states, err := extractStatesFromDefinition(def)
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
	states, err := extractStatesFromDefinition(def)
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
	states, err := extractStatesFromDefinition(def)
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

// TestChoiceRuleOperatorWireForms pins the choice-operator wire form
// through the real parse-and-evaluate path: each operator carries its
// comparison value directly ("the value may be a string, number, boolean,
// or timestamp"). A variable-keyed map is not a documented operator form
// and decodes no comparison operand, so such a rule never matches and the
// routing falls through to Default.
func TestChoiceRuleOperatorWireForms(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		choices []map[string]interface{}
	}{
		{
			name:  "shorthand string equals",
			input: `{"v":"match"}`,
			want:  "hit",
			choices: []map[string]interface{}{
				{"Variable": "$.v", "StringEquals": "match", "Next": "Hit"},
			},
		},
		{
			name:  "shorthand numeric equals",
			input: `{"n":7}`,
			want:  "hit",
			choices: []map[string]interface{}{
				{"Variable": "$.n", "NumericEquals": float64(7), "Next": "Hit"},
			},
		},
		{
			name:  "shorthand boolean equals",
			input: `{"flag":true}`,
			want:  "hit",
			choices: []map[string]interface{}{
				{"Variable": "$.flag", "BooleanEquals": true, "Next": "Hit"},
			},
		},
		{
			name:  "variable-keyed map carries no operand and does not match",
			input: `{"n":7}`,
			want:  "miss",
			choices: []map[string]interface{}{
				{"Variable": "$.n", "NumericEquals": map[string]interface{}{"$.n": 7}, "Next": "Hit"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := newMapTestStore(t)
			e := NewExecutor(store, nil)
			e.region = "us-east-1"

			def := &sfnstore.StateMachineDefinition{
				StartAt: "C",
				States: map[string]interface{}{
					"C":    map[string]interface{}{"Type": "Choice", "Choices": tc.choices, "Default": "Miss"},
					"Hit":  map[string]interface{}{"Type": "Pass", "Result": "hit", "End": true},
					"Miss": map[string]interface{}{"Type": "Pass", "Result": "miss", "End": true},
				},
			}
			states, err := extractStatesFromDefinition(def)
			if err != nil {
				t.Fatalf("extract states failed: %v", err)
			}
			exec := &sfnstore.Execution{
				ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/choice-" + strings.ReplaceAll(tc.name, " ", "-"),
				StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
				Name:            "choice",
				Status:          "RUNNING",
				Input:           tc.input,
			}
			execCtx := &ExecutionContext{
				Execution:     exec,
				Definition:    def,
				CurrentState:  "C",
				Input:         tc.input,
				EventId:       ptrEventID(),
				States:        states,
				QueryLanguage: "JSONPath",
				MapItemIndex:  -1,
			}
			if err := e.executeStates(context.Background(), execCtx); err != nil {
				t.Fatalf("executeStates failed: %v", err)
			}
			want := `"` + tc.want + `"`
			if execCtx.Output != want {
				t.Errorf("output = %s, want %s (routing verdict)", execCtx.Output, want)
			}
		})
	}
}

// TestParameterPathFailureOnMissingPath pins the payload-template failure
// contract: "If the path is legal but cannot be applied successfully, the
// interpreter fails the machine execution" with States.ParameterPathFailure
// — a path that selects no value is never a silent key omission.
func TestParameterPathFailureOnMissingPath(t *testing.T) {
	e := &Executor{}
	params := &sfnstore.Parameters{Values: map[string]interface{}{
		"QueueUrl.$":    "$.queueUrl",
		"MessageBody.$": "$.missing",
	}}
	_, evalErr := e.applyParameters(nil, "", `{"queueUrl":"https://sqs"}`, params)
	if evalErr == nil {
		t.Fatal("a payload-template path that selects no value must fail")
	}
	if evalErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("error code = %s, want States.ParameterPathFailure", evalErr.ErrorCode)
	}

	selector := &sfnstore.ResultSelector{Fields: map[string]interface{}{"total.$": "$.absent"}}
	if _, selErr := e.applyResultSelector(nil, `{"total":1}`, selector, ""); selErr == nil || selErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("ResultSelector miss = %v, want States.ParameterPathFailure", selErr)
	}
}

// TestPayloadTemplateFailureParity pins the payload-template discipline
// across every builder: a ".$" key with a non-string value fails with
// States.ParameterPathFailure in ResultSelector and ItemSelector (no
// silent drop), and a nested path miss classifies identically to the flat
// form — a Catcher on States.ParameterPathFailure must not depend on the
// template's nesting level.
func TestPayloadTemplateFailureParity(t *testing.T) {
	e := &Executor{}

	selector := &sfnstore.ResultSelector{Fields: map[string]interface{}{"total.$": float64(5)}}
	if _, selErr := e.applyResultSelector(nil, `{"total":1}`, selector, ""); selErr == nil || selErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("ResultSelector non-string path value = %v, want States.ParameterPathFailure", selErr)
	}

	itemSel := map[string]interface{}{"flag.$": true}
	if _, selErr := e.applyItemSelectorJSONPath(nil, itemSel, map[string]interface{}{"v": 1}); selErr == nil || selErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("ItemSelector non-string path value = %v, want States.ParameterPathFailure", selErr)
	}

	nested := &sfnstore.Parameters{Values: map[string]interface{}{
		"outer": map[string]interface{}{"inner.$": "$.missing"},
	}}
	_, nestedErr := e.applyParameters(nil, "", `{"v":1}`, nested)
	if nestedErr == nil || nestedErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("nested path miss = %v, want States.ParameterPathFailure (flat/nested parity)", nestedErr)
	}

	nestedSelector := &sfnstore.ResultSelector{Fields: map[string]interface{}{
		"outer": map[string]interface{}{"inner.$": "$.missing"},
	}}
	if _, selErr := e.applyResultSelector(nil, `{"total":1}`, nestedSelector, ""); selErr == nil || selErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("ResultSelector nested path miss = %v, want States.ParameterPathFailure", selErr)
	}

	nestedItem := map[string]interface{}{"outer": map[string]interface{}{"inner.$": "$.missing"}}
	if _, selErr := e.applyItemSelectorJSONPath(nil, nestedItem, map[string]interface{}{"v": 1}); selErr == nil || selErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("ItemSelector nested path miss = %v, want States.ParameterPathFailure", selErr)
	}
}

// TestParameterValueSideDollarSuffixIsLiteral pins that the ".$" convention
// is key-side only: a literal string value ending in ".$" is data.
func TestParameterValueSideDollarSuffixIsLiteral(t *testing.T) {
	e := &Executor{}
	params := &sfnstore.Parameters{Values: map[string]interface{}{
		"version": "v1.0.$",
	}}
	out, evalErr := e.applyParameters(nil, "", `{"v1.0":42}`, params)
	if evalErr != nil {
		t.Fatalf("applyParameters: %v", evalErr)
	}
	if !strings.Contains(out, `"v1.0.$"`) {
		t.Fatalf("literal value ending in .$ was treated as a path: %s", out)
	}
}

// TestParametersNonStringPathKeyFails pins that a payload-template ".$" key
// whose value is not a string fails with States.ParameterPathFailure —
// the suffix marks a path reference, so a non-string value is an invalid
// template at every nesting level, never an entry to drop silently.
func TestParametersNonStringPathKeyFails(t *testing.T) {
	e := &Executor{}
	params := &sfnstore.Parameters{Values: map[string]interface{}{
		"count.$": float64(5),
	}}
	_, evalErr := e.applyParameters(nil, "", `{}`, params)
	if evalErr == nil {
		t.Fatal("a non-string .$/key value must fail the payload template")
	}
	if evalErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("error code = %s, want States.ParameterPathFailure", evalErr.ErrorCode)
	}

	nested := &sfnstore.Parameters{Values: map[string]interface{}{
		"outer": map[string]interface{}{"flag.$": true},
	}}
	_, nestedErr := e.applyParameters(nil, "", `{}`, nested)
	if nestedErr == nil {
		t.Fatal("a nested non-string .$/key value must fail the payload template")
	}
	if nestedErr.ErrorCode != "States.ParameterPathFailure" {
		t.Fatalf("nested error code = %s, want States.ParameterPathFailure", nestedErr.ErrorCode)
	}
}

// TestItemSelectorNestedTemplate pins that non-path keys in an ItemSelector
// still resolve nested payload templates instead of passing the ".$" keys
// through verbatim.
func TestItemSelectorNestedTemplate(t *testing.T) {
	e := &Executor{}
	selector := map[string]interface{}{
		"static": map[string]interface{}{"inner.$": "$.v"},
	}
	out, evalErr := e.applyItemSelectorJSONPath(nil, selector, map[string]interface{}{"v": 7})
	if evalErr != nil {
		t.Fatalf("applyItemSelectorJSONPath: %v", evalErr)
	}
	obj, _ := out.(map[string]interface{})
	nested, _ := obj["static"].(map[string]interface{})
	if nested == nil || fmt.Sprint(nested["inner"]) != "7" {
		t.Fatalf("nested template not resolved: %+v", out)
	}
}

// TestResultPathMatchFailure pins the spec example: "Suppose a state's
// input is the string "foo", and its "ResultPath" field has the value
// "$.x" … Then ResultPath cannot apply and the interpreter fails the
// machine with an Error Name of "States.ResultPathMatchFailure"."
func TestResultPathMatchFailure(t *testing.T) {
	e := &Executor{}
	if _, err := e.applyResultPath(`"foo"`, `{"r":1}`, "$.x"); err == nil || err.ErrorCode != "States.ResultPathMatchFailure" {
		t.Fatalf("ResultPath on non-object input = %v, want States.ResultPathMatchFailure", err)
	}
	if _, err := e.buildCatchOutput(`5`, "CustomError", "boom", "$.err"); err == nil || err.ErrorCode != "States.ResultPathMatchFailure" {
		t.Fatalf("Catch ResultPath on non-object input = %v, want States.ResultPathMatchFailure", err)
	}
}

// TestBuildCatchOutputEscapesQuotes pins that a Catch handler's error
// output is always valid JSON: quotes inside the error code or cause are
// escaped by marshalling, never pasted raw into the document.
func TestBuildCatchOutputEscapesQuotes(t *testing.T) {
	e := &Executor{}
	out, execErr := e.buildCatchOutput(`{}`, `Code"Quote`, `he said "hi"`, "")
	if execErr != nil {
		t.Fatalf("buildCatchOutput errored: %v", execErr)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("catch output is not valid JSON: %v (%s)", err, out)
	}
	if decoded["Error"] != `Code"Quote` || decoded["Cause"] != `he said "hi"` {
		t.Errorf("catch output round-trip = %v, want the quoted strings preserved", decoded)
	}
}

// TestInputPathUnresolvableFails pins that an InputPath that cannot apply
// fails the state with States.Runtime ("errors at runtime, such as
// attempting to apply InputPath or OutputPath on a null JSON payload")
// instead of silently substituting an empty object.
func TestInputPathUnresolvableFails(t *testing.T) {
	e := &Executor{}
	if _, err := e.applyInputPath(nil, `{"a":1}`, "$.missing"); err == nil || err.ErrorCode != "States.Runtime" {
		t.Fatalf("missing InputPath selection = %v, want States.Runtime", err)
	}
	if _, err := e.applyInputPath(nil, `null`, "$.a"); err == nil || err.ErrorCode != "States.Runtime" {
		t.Fatalf("InputPath on a null payload = %v, want States.Runtime", err)
	}
	if got, err := e.applyInputPath(nil, `[1,2,3]`, "$.1"); err != nil || got != "2" {
		t.Fatalf("array-root InputPath = (%s, %v), want the second element", got, err)
	}
}
