package sfn

import (
	"context"
	"encoding/json"
	"testing"
)

func TestEvaluateAssign(t *testing.T) {
	statesVar := BuildStatesVar(
		map[string]interface{}{"x": 10.0},
		nil, nil, nil,
	)

	t.Run("basic assign", func(t *testing.T) {
		assign := map[string]interface{}{
			"sum":  "{% $states.input.x + 5 %}",
			"lit":  "static",
			"num":  42.0,
			"null": nil,
		}
		result, err := evaluateAssign(context.Background(), assign, statesVar, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result["sum"] != 15.0 {
			t.Fatalf("sum: expected 15, got %v", result["sum"])
		}
		if result["lit"] != "static" {
			t.Fatalf("lit: expected 'static', got %v", result["lit"])
		}
		if result["num"] != 42.0 {
			t.Fatalf("num: expected 42, got %v", result["num"])
		}
		if result["null"] != nil {
			t.Fatalf("null: expected nil, got %v", result["null"])
		}
	})

	t.Run("empty assign", func(t *testing.T) {
		result, err := evaluateAssign(context.Background(), map[string]interface{}{}, statesVar, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected empty, got %v", result)
		}
	})

	t.Run("assign with scope variables", func(t *testing.T) {
		scope := NewVariableScope(nil)
		scope.SetAll(map[string]interface{}{"prev": 5.0})

		assign := map[string]interface{}{
			"next": "{% $prev + 1 %}",
		}
		result, err := evaluateAssign(context.Background(), assign, statesVar, scope)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result["next"] != 6.0 {
			t.Fatalf("next: expected 6, got %v", result["next"])
		}
	})

	t.Run("assign with bad expression type error", func(t *testing.T) {
		assign := map[string]interface{}{
			"bad": "{% $length(42) %}",
		}
		_, err := evaluateAssign(context.Background(), assign, statesVar, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestValidateDefinitionJSONataFields(t *testing.T) {
	t.Run("JSONPath-only fields in JSONata state rejected", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "InputPath": "$"}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err == nil {
			t.Fatal("expected error for InputPath in JSONata state")
		}
	})

	t.Run("JSONata-only fields in JSONPath state rejected", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Task", "Arguments": {"x": 1}}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err == nil {
			t.Fatal("expected error for Arguments in JSONPath state")
		}
	})

	t.Run("JSONata Items in JSONPath Map rejected", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Map", "Items": "{% $states.input %}"}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err == nil {
			t.Fatal("expected error for Items in JSONPath Map state")
		}
	})

	t.Run("JSONata Condition in JSONPath Choice rejected", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Choice", "Choices": [{"Condition": "{% true %}"}]}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err == nil {
			t.Fatal("expected error for Condition in JSONPath Choice state")
		}
	})

	t.Run("valid JSONata definition passes", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Pass", "Output": "{% $states.input %}", "End": true}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid JSONPath definition passes", func(t *testing.T) {
		def := `{"StartAt": "A", "States": {"A": {"Type": "Pass", "InputPath": "$", "Result": "hello", "End": true}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("mixed mode: JSONata state in JSONPath machine", func(t *testing.T) {
		def := `{"StartAt": "A", "States": {"A": {"Type": "Pass", "InputPath": "$", "Next": "B"}, "B": {"Type": "Pass", "QueryLanguage": "JSONata", "Output": "{% $states.input %}", "End": true}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("mixed mode: JSONPath state in JSONata machine rejected", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "InputPath": "$"}}}`
		err := validateDefinitionStructure(def, "STANDARD")
		if err == nil {
			t.Fatal("expected error for InputPath in JSONata-default machine")
		}
	})

	t.Run("JSONPath-only fields per state type", func(t *testing.T) {
		tests := []struct {
			name string
			def  string
		}{
			{"Task ResultPath", `{"QueryLanguage": "JSONata", "States": {"T": {"Type": "Task", "Resource": "arn", "ResultPath": "$.r"}}}`},
			{"Task ResultSelector", `{"QueryLanguage": "JSONata", "States": {"T": {"Type": "Task", "Resource": "arn", "ResultSelector": {"x.$": "$.y"}}}}`},
			{"Task TimeoutSecondsPath", `{"QueryLanguage": "JSONata", "States": {"T": {"Type": "Task", "Resource": "arn", "TimeoutSecondsPath": "$.t"}}}`},
			{"Task HeartbeatSecondsPath", `{"QueryLanguage": "JSONata", "States": {"T": {"Type": "Task", "Resource": "arn", "HeartbeatSecondsPath": "$.h"}}}`},
			{"Pass ResultPath", `{"QueryLanguage": "JSONata", "States": {"P": {"Type": "Pass", "ResultPath": "$.r"}}}`},
			{"Pass ResultSelector", `{"QueryLanguage": "JSONata", "States": {"P": {"Type": "Pass", "ResultSelector": {"x.$": "$.y"}}}}`},
			{"Map ItemsPath", `{"QueryLanguage": "JSONata", "States": {"M": {"Type": "Map", "ItemsPath": "$.items"}}}`},
			{"Map ResultPath", `{"QueryLanguage": "JSONata", "States": {"M": {"Type": "Map", "Items": "{% [] %}", "ResultPath": "$.r"}}}`},
			{"Parallel ResultPath", `{"QueryLanguage": "JSONata", "States": {"Pa": {"Type": "Parallel", "Branches": [[]], "ResultPath": "$.r"}}}`},
			{"Fail CausePath", `{"QueryLanguage": "JSONata", "States": {"F": {"Type": "Fail", "CausePath": "$.c"}}}`},
			{"Fail ErrorPath", `{"QueryLanguage": "JSONata", "States": {"F": {"Type": "Fail", "ErrorPath": "$.e"}}}`},
			{"Wait SecondsPath", `{"QueryLanguage": "JSONata", "States": {"W": {"Type": "Wait", "SecondsPath": "$.s"}}}`},
			{"Wait TimestampPath", `{"QueryLanguage": "JSONata", "States": {"W": {"Type": "Wait", "TimestampPath": "$.t"}}}`},
			{"Choice InputPath", `{"QueryLanguage": "JSONata", "States": {"C": {"Type": "Choice", "Choices": [{"Variable": "$.x", "StringEquals": "a"}], "InputPath": "$"}}}`},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				err := validateDefinitionStructure(tt.def, "STANDARD")
				if err == nil {
					t.Fatal("expected error")
				}
			})
		}
	})
}

func TestExtractVariableReferences(t *testing.T) {
	t.Run("single variable in Assign", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Assign": {"x": "{% $myVar + 1 %}"}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 || refs["A"][0] != "myVar" {
			t.Fatalf("expected [myVar], got %v", refs)
		}
	})

	t.Run("multiple variables in Assign", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Assign": {"x": "{% $var1 %}", "y": "{% $var2 %}"}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 2 {
			t.Fatalf("expected 2 refs, got %v", refs)
		}
	})

	t.Run("variable in Output", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Output": "{% $myVar %}"}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 || refs["A"][0] != "myVar" {
			t.Fatalf("expected [myVar], got %v", refs)
		}
	})

	t.Run("variable in Arguments", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Task", "Resource": "arn", "Arguments": {"x": "{% $myVar %}"}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 || refs["A"][0] != "myVar" {
			t.Fatalf("expected [myVar], got %v", refs)
		}
	})

	t.Run("variable in Items", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Map", "Items": "{% $myItems %}", "ItemProcessor": {"StartAt": "X", "States": {"X": {"Type": "Pass"}}}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 || refs["A"][0] != "myItems" {
			t.Fatalf("expected [myItems], got %v", refs)
		}
	})

	t.Run("variable in Condition", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Choice", "Choices": [{"Condition": "{% $myVar > 0 %}"}]}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 || refs["A"][0] != "myVar" {
			t.Fatalf("expected [myVar], got %v", refs)
		}
	})

	t.Run("$states is not a variable reference", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Assign": {"x": "{% $states.input %}"}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 0 {
			t.Fatalf("$states should not be tracked, got %v", refs)
		}
	})

	t.Run("$context is not a variable reference", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Assign": {"x": "{% $context.Execution.Id %}"}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 0 {
			t.Fatalf("$context should not be tracked, got %v", refs)
		}
	})

	t.Run("duplicate variable deduplicated", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Assign": {"x": "{% $myVar %}", "y": "{% $myVar %}"}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 {
			t.Fatalf("expected 1 unique ref, got %v", refs)
		}
	})

	t.Run("no variable references", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Pass", "Output": "{% $states.input %}", "End": true}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 0 {
			t.Fatalf("expected no refs, got %v", refs)
		}
	})

	t.Run("multiple states", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Assign": {"x": "{% $var1 %}"}}, "B": {"Type": "Pass", "Assign": {"y": "{% $var2 %}"}}}}`
		refs := extractVariableReferences(def)
		if refs["A"][0] != "var1" || refs["B"][0] != "var2" {
			t.Fatalf("expected var1 in A, var2 in B, got %v", refs)
		}
	})

	t.Run("variable in nested object", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Output": {"nested": {"val": "{% $myVar %}"}}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 || refs["A"][0] != "myVar" {
			t.Fatalf("expected [myVar], got %v", refs)
		}
	})

	t.Run("variable in array", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Output": ["{% $myVar %}"]}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 1 || refs["A"][0] != "myVar" {
			t.Fatalf("expected [myVar], got %v", refs)
		}
	})

	t.Run("non-expression strings ignored", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Assign": {"x": "plain string"}}}}`
		refs := extractVariableReferences(def)
		if len(refs["A"]) != 0 {
			t.Fatalf("expected no refs for plain strings, got %v", refs)
		}
	})
}

func TestValueSize(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected int64
	}{
		{"string", "hello", 7},
		{"number", 42.0, 2},
		{"null", nil, 4},
		{"bool true", true, 4},
		{"bool false", false, 5},
		{"empty string", "", 2},
		{"empty object", map[string]interface{}{}, 2},
		{"empty array", []interface{}{}, 2},
		{"object with values", map[string]interface{}{"a": 1.0}, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := valueSize(tt.value)
			if size != tt.expected {
				b, _ := json.Marshal(tt.value)
				t.Fatalf("valueSize(%v) = %d, want %d (json: %s)", tt.value, size, tt.expected, string(b))
			}
		})
	}
}

func TestDeepCopyValue(t *testing.T) {
	t.Run("deep copy map", func(t *testing.T) {
		orig := map[string]interface{}{"nested": map[string]interface{}{"x": 1.0}}
		copy := deepCopyValue(orig).(map[string]interface{})
		orig["nested"].(map[string]interface{})["x"] = 99.0
		if copy["nested"].(map[string]interface{})["x"] != 1.0 {
			t.Fatal("deep copy should be independent")
		}
	})

	t.Run("deep copy array", func(t *testing.T) {
		orig := []interface{}{map[string]interface{}{"x": 1.0}}
		copy := deepCopyValue(orig).([]interface{})
		orig[0].(map[string]interface{})["x"] = 99.0
		if copy[0].(map[string]interface{})["x"] != 1.0 {
			t.Fatal("deep copy should be independent")
		}
	})

	t.Run("scalar passes through", func(t *testing.T) {
		orig := "hello"
		copy := deepCopyValue(orig)
		if copy != orig {
			t.Fatalf("scalar should pass through, got %v", copy)
		}
	})
}

func TestNormalizeResult(t *testing.T) {
	t.Run("normal nil stays nil", func(t *testing.T) {
		if NormalizeResult(nil) != nil {
			t.Fatal("nil should stay nil")
		}
	})

	t.Run("string passes through", func(t *testing.T) {
		if NormalizeResult("hello") != "hello" {
			t.Fatal("string should pass through")
		}
	})

	t.Run("number passes through", func(t *testing.T) {
		if NormalizeResult(42.0) != 42.0 {
			t.Fatal("number should pass through")
		}
	})
}

func TestBuildStatesVarWithContext(t *testing.T) {
	statesVar := BuildStatesVar(map[string]interface{}{"msg": "hello"}, nil, nil, map[string]interface{}{
		"Execution": map[string]interface{}{"Id": "arn:test"},
		"State":     map[string]interface{}{"Name": "Test"},
	})

	t.Run("states.input is accessible", func(t *testing.T) {
		statesMap := statesVar["states"].(map[string]interface{})
		input := statesMap["input"]
		m, ok := input.(map[string]interface{})
		if !ok || m["msg"] != "hello" {
			t.Fatalf("expected input with msg=hello, got %v", input)
		}
	})

	result, err := EvaluateJSONata(context.Background(), "$states.input.msg", nil, statesVar)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if result != "hello" {
		t.Fatalf("expected 'hello', got %v (%T)", result, result)
	}
}
