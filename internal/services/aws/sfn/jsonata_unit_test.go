package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	gnata "github.com/recolabs/gnata"
)

func TestEvaluateJSONata_BasicExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		data     interface{}
		vars     map[string]interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "string literal",
			expr:     `"hello"`,
			expected: "hello",
		},
		{
			name:     "numeric literal",
			expr:     `42`,
			expected: 42.0,
		},
		{
			name:     "null literal",
			expr:     `null`,
			expected: nil,
		},
		{
			name:     "boolean true",
			expr:     `true`,
			expected: true,
		},
		{
			name:     "boolean false",
			expr:     `false`,
			expected: false,
		},
		{
			name:     "simple field access",
			expr:     `$.name`,
			data:     map[string]interface{}{"name": "Alice"},
			expected: "Alice",
		},
		{
			name:     "nested field access",
			expr:     `$.a.b.c`,
			data:     map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": "deep"}}},
			expected: "deep",
		},
		{
			name:     "arithmetic",
			expr:     `2 + 3 * 4`,
			expected: 14.0,
		},
		{
			name:     "string comparison",
			expr:     `$.status = "ACTIVE"`,
			data:     map[string]interface{}{"status": "ACTIVE"},
			expected: true,
		},
		{
			name:     "numeric comparison",
			expr:     `$.count > 5`,
			data:     map[string]interface{}{"count": 10.0},
			expected: true,
		},
		{
			name:     "array construction",
			expr:     `[1, 2, 3]`,
			expected: []interface{}{1.0, 2.0, 3.0},
		},
		{
			name:     "object construction",
			expr:     `{"x": 1, "y": 2}`,
			expected: map[string]interface{}{"x": 1.0, "y": 2.0},
		},
		{
			name:     "array filter",
			expr:     `items[price > 15]`,
			data:     map[string]interface{}{"items": []interface{}{map[string]interface{}{"price": 10.0}, map[string]interface{}{"price": 20.0}}},
			expected: map[string]interface{}{"price": 20.0},
		},
		{
			name:     "$sum function",
			expr:     `$sum(items.price)`,
			data:     map[string]interface{}{"items": []interface{}{map[string]interface{}{"price": 10.0}, map[string]interface{}{"price": 20.0}}},
			expected: 30.0,
		},
		{
			name:     "$count returns float64",
			expr:     `$count(items)`,
			data:     map[string]interface{}{"items": []interface{}{1.0, 2.0, 3.0}},
			expected: 3.0,
		},
		{
			name:     "$exists true",
			expr:     `$exists($.name)`,
			data:     map[string]interface{}{"name": "test"},
			expected: true,
		},
		{
			name:     "$exists false",
			expr:     `$exists($.missing)`,
			data:     map[string]interface{}{"name": "test"},
			expected: false,
		},
		{
			name:     "variable via $states",
			expr:     `$states.result`,
			vars:     map[string]interface{}{"states": map[string]interface{}{"input": nil, "result": "done"}},
			expected: "done",
		},
		{
			name:     "variable via custom var",
			expr:     `$myVar + 1`,
			vars:     map[string]interface{}{"myVar": 9.0},
			expected: 10.0,
		},
		{
			name:    "syntax error",
			expr:    `{{invalid`,
			wantErr: true,
		},
		{
			name:    "type error",
			expr:    `$length(42)`,
			wantErr: true,
		},
		{
			name:     "undefined variable returns nil",
			expr:     `$nonexistent`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EvaluateJSONata(context.Background(), tt.expr, tt.data, tt.vars)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Fatalf("expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

func TestEvaluateJSONata_AWSFunctions(t *testing.T) {
	t.Run("$uuid format", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$uuid()`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		s, ok := result.(string)
		if !ok {
			t.Fatalf("expected string, got %T", result)
		}
		if len(s) != 36 {
			t.Fatalf("expected UUID format (36 chars), got: %s", s)
		}
	})

	t.Run("$partition basic", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$partition([1,2,3,4,5], 2)`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := []interface{}{[]interface{}{1.0, 2.0}, []interface{}{3.0, 4.0}, []interface{}{5.0}}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$partition exact", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$partition([1,2,3,4], 2)`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := []interface{}{[]interface{}{1.0, 2.0}, []interface{}{3.0, 4.0}}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$partition single", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$partition([1,2,3], 10)`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := []interface{}{[]interface{}{1.0, 2.0, 3.0}}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$partition empty array", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$partition([], 2)`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		arr, ok := result.([]interface{})
		if !ok || len(arr) != 0 {
			t.Fatalf("expected empty array, got %v (%T)", result, result)
		}
	})

	t.Run("$partition too few args", func(t *testing.T) {
		_, err := EvaluateJSONata(context.Background(), `$partition([1,2])`, nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("$partition non-array first arg", func(t *testing.T) {
		_, err := EvaluateJSONata(context.Background(), `$partition("hello", 2)`, nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("$range ascending", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$range(1, 5)`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := []interface{}{1.0, 2.0, 3.0, 4.0}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$range with delta", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$range(0, 10, 3)`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := []interface{}{0.0, 3.0, 6.0, 9.0}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$range descending", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$range(5, 1, -1)`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := []interface{}{5.0, 4.0, 3.0, 2.0}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$range zero delta", func(t *testing.T) {
		_, err := EvaluateJSONata(context.Background(), `$range(0, 5, 0)`, nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("$hash MD5", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$hash("hello", "MD5")`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != "5d41402abc4b2a76b9719d911017c592" {
			t.Fatalf("expected MD5 hash, got %v", result)
		}
	})

	t.Run("$hash SHA-256", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$hash("hello", "SHA-256")`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
			t.Fatalf("expected SHA-256 hash, got %v", result)
		}
	})

	t.Run("$hash SHA-512", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$hash("hello", "SHA-512")`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043" {
			t.Fatalf("expected SHA-512 hash, got %v", result)
		}
	})

	t.Run("$hash unsupported algorithm", func(t *testing.T) {
		_, err := EvaluateJSONata(context.Background(), `$hash("hello", "BLAKE2")`, nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("$hash non-string first arg", func(t *testing.T) {
		_, err := EvaluateJSONata(context.Background(), `$hash(42, "MD5")`, nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("$parse JSON object", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$parse('{"foo":"bar"}')`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := map[string]interface{}{"foo": "bar"}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$parse JSON array", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$parse('[1,2,3]')`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		expected := []interface{}{1.0, 2.0, 3.0}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("$parse invalid JSON", func(t *testing.T) {
		_, err := EvaluateJSONata(context.Background(), `$parse('not json')`, nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("$parse non-string arg", func(t *testing.T) {
		_, err := EvaluateJSONata(context.Background(), `$parse(42)`, nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("$random returns 0-1", func(t *testing.T) {
		result, err := EvaluateJSONata(context.Background(), `$random()`, nil, nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		f, ok := result.(float64)
		if !ok {
			t.Fatalf("expected float64, got %T", result)
		}
		if f < 0 || f >= 1 {
			t.Fatalf("expected 0 <= x < 1, got %v", f)
		}
	})
}

func TestBuildStatesVar(t *testing.T) {
	t.Run("input only", func(t *testing.T) {
		input := map[string]interface{}{"x": 1.0}
		result := BuildStatesVar(input, nil, nil, nil)
		states := result["states"].(map[string]interface{})
		if !deepEqual(states["input"], input) {
			t.Fatalf("input mismatch: got %v", states["input"])
		}
		if states["result"] != nil {
			t.Fatal("result should be nil")
		}
	})

	t.Run("input and result", func(t *testing.T) {
		input := map[string]interface{}{"x": 1.0}
		resultVal := map[string]interface{}{"status": "OK"}
		result := BuildStatesVar(input, resultVal, nil, nil)
		states := result["states"].(map[string]interface{})
		if !deepEqual(states["input"], input) {
			t.Fatalf("input mismatch: got %v", states["input"])
		}
		if !deepEqual(states["result"], resultVal) {
			t.Fatalf("result mismatch: got %v", states["result"])
		}
		if _, ok := states["errorOutput"]; ok {
			t.Fatal("errorOutput should not exist")
		}
	})

	t.Run("with errorOutput", func(t *testing.T) {
		errOut := map[string]interface{}{"Error": "test"}
		result := BuildStatesVar(nil, nil, errOut, nil)
		states := result["states"].(map[string]interface{})
		if !deepEqual(states["errorOutput"], errOut) {
			t.Fatalf("errorOutput mismatch: got %v", states["errorOutput"])
		}
	})

	t.Run("with context", func(t *testing.T) {
		ctx := map[string]interface{}{"State": map[string]interface{}{"Name": "TestState"}}
		result := BuildStatesVar(nil, nil, nil, ctx)
		states := result["states"].(map[string]interface{})
		if !deepEqual(states["context"], ctx) {
			t.Fatalf("context mismatch: got %v", states["context"])
		}
	})

	t.Run("all fields", func(t *testing.T) {
		input := "data"
		resultVal := 42.0
		errOut := map[string]interface{}{"Error": "fail"}
		ctx := map[string]interface{}{"Execution": map[string]interface{}{"Id": "exec-1"}}
		result := BuildStatesVar(input, resultVal, errOut, ctx)
		states := result["states"].(map[string]interface{})
		if states["input"] != input {
			t.Fatalf("input mismatch")
		}
		if states["result"] != resultVal {
			t.Fatalf("result mismatch")
		}
		if !deepEqual(states["errorOutput"], errOut) {
			t.Fatalf("errorOutput mismatch: got %v", states["errorOutput"])
		}
		if !deepEqual(states["context"], ctx) {
			t.Fatalf("context mismatch: got %v", states["context"])
		}
	})
}

func TestResolveTemplate_Walker(t *testing.T) {
	vars := map[string]interface{}{
		"states": map[string]interface{}{
			"input":  map[string]interface{}{"x": 10.0, "y": 20.0},
			"result": nil,
		},
	}

	t.Run("expression string evaluates", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "{% $states.input.x + $states.input.y %}", nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != 30.0 {
			t.Fatalf("expected 30, got %v", result)
		}
	})

	t.Run("non-expression string passes through", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "static string", nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != "static string" {
			t.Fatalf("expected 'static string', got %v", result)
		}
	})

	t.Run("partial expression string passes through (AWS spec)", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "prefix {% $x %} suffix", nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != "prefix {% $x %} suffix" {
			t.Fatalf("expected literal, got %v", result)
		}
	})

	t.Run("expression returning null", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "{% null %}", nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})

	t.Run("expression returning number", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "{% 42 %}", nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != 42.0 {
			t.Fatalf("expected 42, got %v", result)
		}
	})

	t.Run("non-string non-expression passes through", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), 99.0, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != 99.0 {
			t.Fatalf("expected 99, got %v", result)
		}
	})

	t.Run("object with mixed values", func(t *testing.T) {
		input := map[string]interface{}{
			"expr": "{% $states.input.x %}",
			"lit":  "static",
			"num":  5.0,
			"null": nil,
			"bool": true,
		}
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("expected map, got %T", result)
		}
		if m["expr"] != 10.0 {
			t.Fatalf("expr: expected 10, got %v", m["expr"])
		}
		if m["lit"] != "static" {
			t.Fatalf("lit: expected 'static', got %v", m["lit"])
		}
		if m["num"] != 5.0 {
			t.Fatalf("num: expected 5, got %v", m["num"])
		}
		if m["null"] != nil {
			t.Fatalf("null: expected nil, got %v", m["null"])
		}
		if m["bool"] != true {
			t.Fatalf("bool: expected true, got %v", m["bool"])
		}
	})

	t.Run("array with mixed values", func(t *testing.T) {
		input := []interface{}{"{% $states.input.x %}", "static", "{% $states.input.y * 2 %}"}
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		arr, ok := result.([]interface{})
		if !ok {
			t.Fatalf("expected array, got %T", result)
		}
		if arr[0] != 10.0 {
			t.Fatalf("[0]: expected 10, got %v", arr[0])
		}
		if arr[1] != "static" {
			t.Fatalf("[1]: expected 'static', got %v", arr[1])
		}
		if arr[2] != 40.0 {
			t.Fatalf("[2]: expected 40, got %v", arr[2])
		}
	})

	t.Run("nested object with expressions", func(t *testing.T) {
		input := map[string]interface{}{
			"nested": map[string]interface{}{
				"value": "{% $states.input.x %}",
				"deep": map[string]interface{}{
					"x": "{% $states.input.y %}",
				},
			},
		}
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		m := result.(map[string]interface{})
		nested := m["nested"].(map[string]interface{})
		if nested["value"] != 10.0 {
			t.Fatalf("nested.value: expected 10, got %v", nested["value"])
		}
		deep := nested["deep"].(map[string]interface{})
		if deep["x"] != 20.0 {
			t.Fatalf("nested.deep.x: expected 20, got %v", deep["x"])
		}
	})

	t.Run("undefined path returns nil not error", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "{% $states.input.nonexistent.bad.path %}", nil, vars)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != nil {
			t.Fatalf("expected nil for undefined path, got %v", result)
		}
	})

	t.Run("empty map", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), map[string]interface{}{}, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("expected map, got %T", result)
		}
		if len(m) != 0 {
			t.Fatalf("expected empty map, got %v", m)
		}
	})

	t.Run("empty array", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), []interface{}{}, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		arr, ok := result.([]interface{})
		if !ok {
			t.Fatalf("expected array, got %T", result)
		}
		if len(arr) != 0 {
			t.Fatalf("expected empty array, got %v", arr)
		}
	})
}

func TestEvaluateExpressionValue(t *testing.T) {
	vars := map[string]interface{}{"states": map[string]interface{}{"input": map[string]interface{}{"x": 5.0}}}

	t.Run("expression string evaluates", func(t *testing.T) {
		result, err := EvaluateExpressionValue(context.Background(), "{% $states.input.x + 1 %}", nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != 6.0 {
			t.Fatalf("expected 6, got %v", result)
		}
	})

	t.Run("non-expression string passes through", func(t *testing.T) {
		result, err := EvaluateExpressionValue(context.Background(), "hello", nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != "hello" {
			t.Fatalf("expected 'hello', got %v", result)
		}
	})

	t.Run("non-string passes through", func(t *testing.T) {
		result, err := EvaluateExpressionValue(context.Background(), 42.0, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != 42.0 {
			t.Fatalf("expected 42, got %v", result)
		}
	})

	t.Run("nil passes through", func(t *testing.T) {
		result, err := EvaluateExpressionValue(context.Background(), nil, nil, vars)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})
}

func TestIsJSONataExpressionValue(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"{% $x %}", true},
		{"{% null %}", true},
		{"not expression", false},
		{"", false},
		{42, false},
		{nil, false},
		{true, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if IsJSONataExpressionValue(tt.input) != tt.expected {
				t.Fatalf("IsJSONataExpressionValue(%v) = %v, want %v", tt.input, !tt.expected, tt.expected)
			}
		})
	}
}

func TestVariableScope_Basic(t *testing.T) {
	t.Run("new scope has no variables", func(t *testing.T) {
		scope := NewVariableScope(nil)
		if _, ok := scope.Get("x"); ok {
			t.Fatal("new scope should have no variables")
		}
	})

	t.Run("set and get", func(t *testing.T) {
		scope := NewVariableScope(nil)
		err := scope.SetAll(map[string]interface{}{"x": 1.0, "y": "hello"})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		v, ok := scope.Get("x")
		if !ok || v != 1.0 {
			t.Fatalf("x: expected 1, got %v", v)
		}
		v, ok = scope.Get("y")
		if !ok || v != "hello" {
			t.Fatalf("y: expected 'hello', got %v", v)
		}
	})

	t.Run("get missing returns false", func(t *testing.T) {
		scope := NewVariableScope(nil)
		scope.SetAll(map[string]interface{}{"x": 1.0})
		_, ok := scope.Get("missing")
		if ok {
			t.Fatal("missing variable should not exist")
		}
	})

	t.Run("set empty map is no-op", func(t *testing.T) {
		scope := NewVariableScope(nil)
		err := scope.SetAll(map[string]interface{}{})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if len(scope.Snapshot()) != 0 {
			t.Fatal("scope should remain empty")
		}
	})
}

func TestVariableScope_ParentChild(t *testing.T) {
	t.Run("child reads parent variable", func(t *testing.T) {
		parent := NewVariableScope(nil)
		parent.SetAll(map[string]interface{}{"x": 10.0})
		child := parent.NewChild()

		v, ok := child.Get("x")
		if !ok || v != 10.0 {
			t.Fatalf("child should see parent x=10, got %v", v)
		}
	})

	t.Run("child has own variables", func(t *testing.T) {
		parent := NewVariableScope(nil)
		parent.SetAll(map[string]interface{}{"x": 10.0})
		child := parent.NewChild()
		child.SetAll(map[string]interface{}{"y": 20.0})

		v, ok := child.Get("y")
		if !ok || v != 20.0 {
			t.Fatalf("child should have y=20, got %v", v)
		}
		v, ok = child.Get("x")
		if !ok || v != 10.0 {
			t.Fatalf("child should see parent x=10, got %v", v)
		}
	})

	t.Run("parent cannot see child variable", func(t *testing.T) {
		parent := NewVariableScope(nil)
		parent.SetAll(map[string]interface{}{"x": 10.0})
		child := parent.NewChild()
		child.SetAll(map[string]interface{}{"y": 20.0})

		_, ok := parent.Get("y")
		if ok {
			t.Fatal("parent should not see child variable y")
		}
	})

	t.Run("child overwrites own variable", func(t *testing.T) {
		parent := NewVariableScope(nil)
		parent.SetAll(map[string]interface{}{"x": 10.0})
		child := parent.NewChild()

		err := child.SetAll(map[string]interface{}{"x": 99.0})
		if err == nil {
			t.Fatal("expected shadowing error")
		}
	})

	t.Run("GetAll includes parent variables", func(t *testing.T) {
		parent := NewVariableScope(nil)
		parent.SetAll(map[string]interface{}{"x": 1.0})
		child := parent.NewChild()
		child.SetAll(map[string]interface{}{"y": 2.0})

		all := child.GetAll()
		if all["x"] != 1.0 {
			t.Fatalf("GetAll should include parent x, got %v", all["x"])
		}
		if all["y"] != 2.0 {
			t.Fatalf("GetAll should include child y, got %v", all["y"])
		}
	})

	t.Run("GetAll child cannot override parent", func(t *testing.T) {
		parent := NewVariableScope(nil)
		parent.SetAll(map[string]interface{}{"x": 1.0})
		child := parent.NewChild()

		err := child.SetAll(map[string]interface{}{"x": 2.0})
		if err == nil {
			t.Fatal("expected shadowing error")
		}

		all := child.GetAll()
		if all["x"] != 1.0 {
			t.Fatalf("GetAll should still show parent x=1, got %v", all["x"])
		}
	})

	t.Run("shadowing parent variable is rejected", func(t *testing.T) {
		parent := NewVariableScope(nil)
		parent.SetAll(map[string]interface{}{"x": 1.0})
		child := parent.NewChild()

		err := child.SetAll(map[string]interface{}{"x": 99.0})
		if err == nil {
			t.Fatal("expected shadowing error")
		}
	})
}

func TestVariableScope_SizeLimits(t *testing.T) {
	t.Run("single variable within limit", func(t *testing.T) {
		scope := NewVariableScope(nil)
		small := strings.Repeat("x", 1000)
		err := scope.SetAll(map[string]interface{}{"data": small})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	})

	t.Run("single variable exceeds per-variable limit", func(t *testing.T) {
		scope := NewVariableScope(nil)
		huge := strings.Repeat("x", 300*1024)
		err := scope.SetAll(map[string]interface{}{"data": huge})
		if err == nil {
			t.Fatal("expected size limit error")
		}
	})

	t.Run("combined Assign exceeds per-assign limit", func(t *testing.T) {
		scope := NewVariableScope(nil)
		big1 := strings.Repeat("a", 200*1024)
		big2 := strings.Repeat("b", 200*1024)
		err := scope.SetAll(map[string]interface{}{"a": big1, "b": big2})
		if err == nil {
			t.Fatal("expected combined size limit error")
		}
	})
}

func TestVariableScope_Snapshot(t *testing.T) {
	t.Run("snapshot is independent copy", func(t *testing.T) {
		scope := NewVariableScope(nil)
		scope.SetAll(map[string]interface{}{"x": map[string]interface{}{"nested": true}})
		snap := scope.Snapshot()

		original := scope.Snapshot()
		origMap := original["x"].(map[string]interface{})
		origMap["nested"] = false

		snapMap := snap["x"].(map[string]interface{})
		if snapMap["nested"] != true {
			t.Fatal("snapshot should be independent copy")
		}
	})
}

func TestValidateVariableName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple", "myVar", false},
		{"underscore start", "_private", false},
		{"dollar start", "$dollar", false},
		{"with digits", "var123", false},
		{"single letter", "x", false},
		{"empty", "", true},
		{"starts with digit", "1var", true},
		{"starts with hyphen", "-var", true},
		{"contains space", "my var", true},
		{"contains dot", "my.var", true},
		{"contains hyphen", "my-var", true},
		{"starts with uppercase", "MyVar", false},
		{"unicode letter start", "変数", false},
		{"too long", strings.Repeat("a", 81), true},
		{"max length", strings.Repeat("a", 80), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVariableName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateVariableName(%q) err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestBuildVarsMap(t *testing.T) {
	t.Run("with scope", func(t *testing.T) {
		scope := NewVariableScope(nil)
		scope.SetAll(map[string]interface{}{"myVar": "value"})
		statesVar := BuildStatesVar(map[string]interface{}{"x": 1.0}, nil, nil, nil)

		result := buildVarsMap(statesVar, scope)
		statesInner, ok := result["states"].(map[string]interface{})
		if !ok {
			t.Fatal("states key missing or wrong type")
		}
		input := statesInner["input"].(map[string]interface{})
		if input["x"] != 1.0 {
			t.Fatal("states.input.x mismatch")
		}
		if result["myVar"] != "value" {
			t.Fatal("myVar mismatch")
		}
	})

	t.Run("nil scope", func(t *testing.T) {
		statesVar := BuildStatesVar(nil, nil, nil, nil)
		result := buildVarsMap(statesVar, nil)
		if _, ok := result["states"]; !ok {
			t.Fatal("states key missing")
		}
		if len(result) != 1 {
			t.Fatalf("expected only states, got %v", result)
		}
	})

	t.Run("scope variable does not override states", func(t *testing.T) {
		scope := NewVariableScope(nil)
		scope.SetAll(map[string]interface{}{"states": "should not override"})
		statesVar := BuildStatesVar(nil, nil, nil, nil)

		result := buildVarsMap(statesVar, scope)
		statesInner, ok := result["states"].(map[string]interface{})
		if !ok {
			t.Fatal("states key should be a map, not overridden by scope")
		}
		_ = statesInner
	})
}

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
		err := validateDefinitionJSONataFields(def)
		if err == nil {
			t.Fatal("expected error for InputPath in JSONata state")
		}
	})

	t.Run("JSONata-only fields in JSONPath state rejected", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Task", "Arguments": {"x": 1}}}}`
		err := validateDefinitionJSONataFields(def)
		if err == nil {
			t.Fatal("expected error for Arguments in JSONPath state")
		}
	})

	t.Run("JSONata Items in JSONPath Map rejected", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Map", "Items": "{% $states.input %}"}}}`
		err := validateDefinitionJSONataFields(def)
		if err == nil {
			t.Fatal("expected error for Items in JSONPath Map state")
		}
	})

	t.Run("JSONata Condition in JSONPath Choice rejected", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Choice", "Choices": [{"Condition": "{% true %}"}]}}}`
		err := validateDefinitionJSONataFields(def)
		if err == nil {
			t.Fatal("expected error for Condition in JSONPath Choice state")
		}
	})

	t.Run("valid JSONata definition passes", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Output": "{% $states.input %}"}}}`
		err := validateDefinitionJSONataFields(def)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid JSONPath definition passes", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Pass", "InputPath": "$", "Result": "hello"}}}`
		err := validateDefinitionJSONataFields(def)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("mixed mode: JSONata state in JSONPath machine", func(t *testing.T) {
		def := `{"States": {"A": {"Type": "Pass", "InputPath": "$"}, "B": {"Type": "Pass", "QueryLanguage": "JSONata", "Output": "{% $states.input %}"}}}`
		err := validateDefinitionJSONataFields(def)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("mixed mode: JSONPath state in JSONata machine rejected", func(t *testing.T) {
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "InputPath": "$"}}}`
		err := validateDefinitionJSONataFields(def)
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
				err := validateDefinitionJSONataFields(tt.def)
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
		def := `{"QueryLanguage": "JSONata", "States": {"A": {"Type": "Pass", "Output": "{% $states.input %}"}}}`
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

func TestEvaluateJSONataWithInputAndVars(t *testing.T) {
	statesVar := BuildStatesVar(
		map[string]interface{}{"msg": "hello jsonata"},
		nil, nil,
		map[string]interface{}{
			"Execution": map[string]interface{}{"Id": "arn:test"},
			"State":     map[string]interface{}{"Name": "Test"},
		},
	)

	t.Run("ResolveTemplate with $states.input", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "{% $states.input %}", nil, statesVar)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		m, ok := result.(map[string]interface{})
		if !ok || m["msg"] != "hello jsonata" {
			t.Fatalf("expected input with msg, got %v (%T)", result, result)
		}
	})

	t.Run("ResolveTemplate with $states.input.msg", func(t *testing.T) {
		result, err := ResolveTemplate(context.Background(), "{% $states.input.msg %}", nil, statesVar)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != "hello jsonata" {
			t.Fatalf("expected 'hello jsonata', got %v (%T)", result, result)
		}
	})

	t.Run("EvalWithVars nil data (mimics integration)", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.input`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithVars(context.Background(), nil, statesVar)
		if err != nil {
			t.Fatalf("eval err: %v", err)
		}
		if result == nil {
			t.Fatalf("expected non-nil result")
		}
		t.Logf("result type: %T, value: %v", result, result)
	})

	t.Run("EvalWithVars with data", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.input`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithVars(context.Background(), map[string]interface{}{"x": 1}, statesVar)
		if err != nil {
			t.Fatalf("eval err: %v", err)
		}
		if result == nil {
			t.Fatalf("expected non-nil result")
		}
		t.Logf("result type: %T, value: %v", result, result)
	})

	t.Run("EvalWithEnvAndVars nil data (mimics integration)", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.input`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithEnvAndVars(context.Background(), nil, awsCustomEnv, statesVar)
		if err != nil {
			t.Fatalf("eval err: %v", err)
		}
		if result == nil {
			t.Fatalf("expected non-nil result")
		}
		t.Logf("result type: %T, value: %v", result, result)
	})

	t.Run("EvalWithEnvAndVars with data", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.input`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithEnvAndVars(context.Background(), map[string]interface{}{"x": 1}, awsCustomEnv, statesVar)
		if err != nil {
			t.Fatalf("eval err: %v", err)
		}
		if result == nil {
			t.Fatalf("expected non-nil result")
		}
		t.Logf("result type: %T, value: %v", result, result)
	})
}
