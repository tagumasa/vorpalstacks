package sfn

import (
	"context"
	"testing"

	gnata "github.com/recolabs/gnata"
)

func TestGnataTemplateResolution(t *testing.T) {
	t.Run("simple expression string", func(t *testing.T) {
		expr, err := gnata.Compile(`"hello"`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "hello" {
			t.Fatalf("expected hello, got %v", result)
		}
	})

	t.Run("non-expression string pass-through", func(t *testing.T) {
		expr, err := gnata.Compile(`"Hello {% $name %}"`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "Hello {% $name %}" {
			t.Fatalf("AWS spec: no interpolation, expected literal string, got %v", result)
		}
	})

	t.Run("numeric literal", func(t *testing.T) {
		expr, err := gnata.Compile(`42`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != 42.0 {
			t.Fatalf("expected 42, got %v", result)
		}
	})

	t.Run("null", func(t *testing.T) {
		expr, err := gnata.Compile(`null`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		normalized := normalizeResult(result)
		if normalized != nil {
			t.Fatalf("expected nil, got %v", normalized)
		}
	})

	t.Run("boolean true", func(t *testing.T) {
		expr, err := gnata.Compile(`true`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != true {
			t.Fatalf("expected true, got %v", result)
		}
	})

	t.Run("array literal", func(t *testing.T) {
		expr, err := gnata.Compile(`[1, 2, 3]`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		expected := []any{1.0, 2.0, 3.0}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})
}

func TestGnataTemplateWalker(t *testing.T) {
	statesVar := map[string]any{
		"input":  map[string]any{"x": 10.0, "y": 20.0},
		"result": nil,
	}
	vars := map[string]any{"states": statesVar}

	t.Run("object with expressions", func(t *testing.T) {
		input := map[string]any{
			"sum":  "{% $states.input.x + $states.input.y %}",
			"diff": "{% $states.input.y - $states.input.x %}",
			"lit":  "static",
		}
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("ResolveTemplate: %v", err)
		}
		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T", result)
		}
		if m["sum"] != 30.0 {
			t.Fatalf("expected sum=30, got %v", m["sum"])
		}
		if m["diff"] != 10.0 {
			t.Fatalf("expected diff=10, got %v", m["diff"])
		}
		if m["lit"] != "static" {
			t.Fatalf("expected lit=static, got %v", m["lit"])
		}
	})

	t.Run("array with expressions", func(t *testing.T) {
		input := []any{"{% $states.input.x %}", "static", "{% $states.input.y * 2 %}"}
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("ResolveTemplate: %v", err)
		}
		arr, ok := result.([]any)
		if !ok {
			t.Fatalf("expected array, got %T", result)
		}
		if arr[0] != 10.0 {
			t.Fatalf("expected [0]=10, got %v", arr[0])
		}
		if arr[1] != "static" {
			t.Fatalf("expected [1]=static, got %v", arr[1])
		}
		if arr[2] != 40.0 {
			t.Fatalf("expected [2]=40, got %v", arr[2])
		}
	})

	t.Run("nested object with expressions", func(t *testing.T) {
		input := map[string]any{
			"nested": map[string]any{
				"value": "{% $states.input.x %}",
			},
		}
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("ResolveTemplate: %v", err)
		}
		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T", result)
		}
		nested, ok := m["nested"].(map[string]any)
		if !ok {
			t.Fatalf("expected nested map, got %T", m["nested"])
		}
		if nested["value"] != 10.0 {
			t.Fatalf("expected nested.value=10, got %v", nested["value"])
		}
	})

	t.Run("non-expression string passthrough", func(t *testing.T) {
		input := "Hello World"
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("ResolveTemplate: %v", err)
		}
		if result != "Hello World" {
			t.Fatalf("expected 'Hello World', got %v", result)
		}
	})

	t.Run("expression returning null", func(t *testing.T) {
		input := "{% null %}"
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("ResolveTemplate: %v", err)
		}
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})

	t.Run("non-string non-expression value passthrough", func(t *testing.T) {
		input := 42.0
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("ResolveTemplate: %v", err)
		}
		if result != 42.0 {
			t.Fatalf("expected 42, got %v", result)
		}
	})

	t.Run("no interpolation in strings", func(t *testing.T) {
		input := "prefix {% $states.input.x %} suffix"
		result, err := ResolveTemplate(context.Background(), input, nil, vars)
		if err != nil {
			t.Fatalf("ResolveTemplate: %v", err)
		}
		if result != "prefix {% $states.input.x %} suffix" {
			t.Fatalf("AWS spec: no interpolation, expected literal, got %v", result)
		}
	})
}

func TestGnataExpressionSyntaxDetection(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"{% $states.input %}", true},
		{"{% $x + 1 %}", true},
		{"{% $sum(items) %}", true},
		{"{% null %}", true},
		{"{% 42 %}", true},
		{"{% true %}", true},
		{"  {% $x %}  ", false},
		{"{% $x %} extra", false},
		{"extra {% $x %}", false},
		{"not an expression", false},
		{"", false},
		{"plain string", false},
		{"42", false},
		{"true", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsExpression(tt.input)
			if result != tt.expected {
				t.Fatalf("IsExpression(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGnataExpressionUnwrap(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{% $states.input %}", "$states.input"},
		{"{% $x + 1 %}", "$x + 1"},
		{"{% null %}", "null"},
		{"{%  %}", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnwrapExpression(tt.input)
			if result != tt.expected {
				t.Fatalf("UnwrapExpression(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
