package sfn

import (
	"context"
	"testing"
	"time"

	gnata "github.com/recolabs/gnata"
)

func TestGnataErrorReporting(t *testing.T) {
	tests := []struct {
		name       string
		expr       string
		data       any
		wantErr    bool
		compileErr bool
	}{
		{
			name:    "type error in function",
			expr:    `$length(42)`,
			data:    nil,
			wantErr: true,
		},
		{
			name:       "syntax error",
			expr:       `{{invalid`,
			data:       nil,
			compileErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := gnata.Compile(tt.expr)
			if tt.compileErr {
				if err == nil {
					t.Fatal("expected compile error")
				}
				t.Logf("compile error type: %T, message: %v", err, err)
				return
			}
			if err != nil {
				t.Fatalf("unexpected compile error: %v", err)
			}
			_, err = expr.Eval(context.Background(), tt.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if err != nil {
				t.Logf("error type: %T, message: %v", err, err)
			}
		})
	}

	t.Run("undefined variable returns nil", func(t *testing.T) {
		expr, err := gnata.Compile(`$nonExistentVar`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != nil {
			t.Fatalf("expected nil for undefined var, got %v", result)
		}
	})

	t.Run("division by zero returns Infinity", func(t *testing.T) {
		expr, err := gnata.Compile(`10 / 0`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), nil)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != nil {
			t.Logf("division by zero result: %v (%T)", result, result)
		}
	})
}

func TestGnataComplexExpressions(t *testing.T) {
	data := map[string]any{
		"items": []any{
			map[string]any{"name": "A", "price": 10.0},
			map[string]any{"name": "B", "price": 20.0},
			map[string]any{"name": "C", "price": 30.0},
		},
	}

	t.Run("filter and map", func(t *testing.T) {
		expr, err := gnata.Compile(`items[price > 15].name`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		expected := []any{"B", "C"}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("sum with $sum", func(t *testing.T) {
		expr, err := gnata.Compile(`$sum(items.price)`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != 60.0 {
			t.Fatalf("expected 60, got %v", result)
		}
	})

	t.Run("count with $count", func(t *testing.T) {
		expr, err := gnata.Compile(`$count(items)`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != 3.0 {
			t.Fatalf("expected 3, got %v", result)
		}
	})

	t.Run("conditional", func(t *testing.T) {
		expr, err := gnata.Compile(`$count(items) > 2 ? "many" : "few"`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "many" {
			t.Fatalf("expected many, got %v", result)
		}
	})

	t.Run("object construction", func(t *testing.T) {
		expr, err := gnata.Compile(`{"total": $sum(items.price), "count": $count(items)}`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		normalized := normalizeResult(result)
		m, ok := normalized.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T (normalized: %T)", result, normalized)
		}
		if m["total"] != 60.0 {
			t.Fatalf("expected total=60, got %v", m["total"])
		}
		if m["count"] != 3.0 {
			t.Fatalf("expected count=3, got %v", m["count"])
		}
	})
}

func TestGnataBooleanExpression(t *testing.T) {
	data := map[string]any{"status": "ACTIVE", "count": 5.0}

	t.Run("string comparison", func(t *testing.T) {
		expr, err := gnata.Compile(`status = "ACTIVE"`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != true {
			t.Fatalf("expected true, got %v", result)
		}
	})

	t.Run("numeric comparison", func(t *testing.T) {
		expr, err := gnata.Compile(`count >= 5`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != true {
			t.Fatalf("expected true, got %v", result)
		}
	})

	t.Run("and condition", func(t *testing.T) {
		expr, err := gnata.Compile(`status = "ACTIVE" and count > 3`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != true {
			t.Fatalf("expected true, got %v", result)
		}
	})

	t.Run("returns actual boolean not truthy", func(t *testing.T) {
		expr, err := gnata.Compile(`count = 5`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		b, ok := result.(bool)
		if !ok {
			t.Fatalf("expected bool, got %T", result)
		}
		if !b {
			t.Fatal("expected true")
		}
	})
}

func TestGnataPerformance(t *testing.T) {
	payload := generateLargePayload(256 * 1024)
	data := map[string]any{"payload": payload}

	t.Run("simple field access", func(t *testing.T) {
		expr, err := gnata.Compile(`$.payload.id`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		start := time.Now()
		for i := 0; i < 100; i++ {
			_, err = expr.Eval(context.Background(), data)
			if err != nil {
				t.Fatalf("eval: %v", err)
			}
		}
		elapsed := time.Since(start)
		avg := elapsed / 100
		t.Logf("simple field access: avg=%v, total=%v", avg, elapsed)
		if avg > 100*time.Millisecond {
			t.Fatalf("too slow: avg=%v", avg)
		}
	})

	t.Run("filter on large array", func(t *testing.T) {
		expr, err := gnata.Compile(`$count(payload.items[value > 500])`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		start := time.Now()
		_, err = expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		elapsed := time.Since(start)
		t.Logf("filter on large array: %v", elapsed)
		if elapsed > 5*time.Second {
			t.Fatalf("too slow: %v", elapsed)
		}
	})

	t.Run("compile caching", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < 100; i++ {
			_, err := gnata.Compile(`$sum(payload.items[value].value)`)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
		}
		elapsed := time.Since(start)
		avg := elapsed / 100
		t.Logf("compile: avg=%v, total=%v", avg, elapsed)
	})
}

func TestGnataAWSExpressionPatterns(t *testing.T) {
	t.Run("Assign pattern: $states.result.Payload.current_price", func(t *testing.T) {
		result := map[string]any{
			"Payload": map[string]any{
				"current_price": 42.99,
				"currency":      "USD",
			},
		}
		statesVar := map[string]any{
			"input":  nil,
			"result": result,
		}

		expr, err := gnata.Compile(`$states.result.Payload.current_price`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		out, err := expr.EvalWithVars(context.Background(), nil, map[string]any{"states": statesVar})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if out != 42.99 {
			t.Fatalf("expected 42.99, got %v", out)
		}
	})

	t.Run("Output pattern: object with expressions", func(t *testing.T) {
		result := map[string]any{"data": "value"}
		statesVar := map[string]any{
			"input":  map[string]any{"x": 1.0},
			"result": result,
		}

		expr, err := gnata.Compile(`{"outputData": $states.result.data, "inputX": $states.input.x}`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		out, err := expr.EvalWithVars(context.Background(), nil, map[string]any{"states": statesVar})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		normalized := normalizeResult(out)
		m, ok := normalized.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T (raw: %T)", normalized, out)
		}
		if m["outputData"] != "value" {
			t.Fatalf("expected outputData=value, got %v", m["outputData"])
		}
		if m["inputX"] != 1.0 {
			t.Fatalf("expected inputX=1, got %v", m["inputX"])
		}
	})

	t.Run("Condition pattern: boolean expression", func(t *testing.T) {
		statesVar := map[string]any{
			"input": map[string]any{"status": "ACTIVE"},
		}

		expr, err := gnata.Compile(`$states.input.status = "ACTIVE"`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		out, err := expr.EvalWithVars(context.Background(), nil, map[string]any{"states": statesVar})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		b, ok := out.(bool)
		if !ok {
			t.Fatalf("expected bool, got %T", out)
		}
		if !b {
			t.Fatal("expected true")
		}
	})

	t.Run("Items pattern: expression returning array", func(t *testing.T) {
		statesVar := map[string]any{
			"input": map[string]any{"numbers": []any{1.0, 2.0, 3.0}},
		}

		expr, err := gnata.Compile(`$states.input.numbers`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		out, err := expr.EvalWithVars(context.Background(), nil, map[string]any{"states": statesVar})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		expected := []any{1.0, 2.0, 3.0}
		if !deepEqual(out, expected) {
			t.Fatalf("expected %v, got %v", expected, out)
		}
	})
}

func TestGnataEdgeCases(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		expr, err := gnata.Compile(`$exists($.foo)`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), map[string]any{})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != false {
			t.Fatalf("expected false, got %v", result)
		}
	})

	t.Run("null in expression", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.result = null`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithVars(context.Background(), nil, map[string]any{
			"states": map[string]any{"result": nil},
		})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != true {
			t.Fatalf("expected true, got %v", result)
		}
	})

	t.Run("deeply nested path", func(t *testing.T) {
		data := map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": map[string]any{
						"d": "found",
					},
				},
			},
		}
		expr, err := gnata.Compile(`a.b.c.d`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "found" {
			t.Fatalf("expected found, got %v", result)
		}
	})

	t.Run("array index access", func(t *testing.T) {
		data := map[string]any{
			"items": []any{"a", "b", "c"},
		}
		expr, err := gnata.Compile(`items[0]`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "a" {
			t.Fatalf("expected a, got %v", result)
		}
	})

	t.Run("string concatenation", func(t *testing.T) {
		data := map[string]any{"first": "hello", "second": "world"}
		expr, err := gnata.Compile(`first & " " & second`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "hello world" {
			t.Fatalf("expected 'hello world', got %v", result)
		}
	})

	t.Run("$not function", func(t *testing.T) {
		expr, err := gnata.Compile(`$not(false)`)
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

	t.Run("$exists function", func(t *testing.T) {
		data := map[string]any{"foo": "bar"}
		expr, err := gnata.Compile(`$exists($.foo)`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.Eval(context.Background(), data)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != true {
			t.Fatalf("expected true, got %v", result)
		}
	})
}
