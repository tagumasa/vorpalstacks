package sfn

import (
	"context"
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
