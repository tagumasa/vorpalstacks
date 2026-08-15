package sfn

import (
	"context"
	"strings"
	"testing"
)

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
