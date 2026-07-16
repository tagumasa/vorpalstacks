package sfn

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	gnata "github.com/recolabs/gnata"
)

func TestGnataBasicEvaluation(t *testing.T) {
	expr, err := gnata.Compile(`$.name`)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	data := map[string]any{"name": "Alice", "age": 30}
	result, err := expr.Eval(context.Background(), data)
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if result != "Alice" {
		t.Fatalf("expected Alice, got %v", result)
	}
}

func TestGnataNestedPath(t *testing.T) {
	expr, err := gnata.Compile(`Account.Order.Product.Price`)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	data := map[string]any{
		"Account": map[string]any{
			"Order": map[string]any{
				"Product": map[string]any{"Price": 34.45},
			},
		},
	}

	result, err := expr.Eval(context.Background(), data)
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if result != 34.45 {
		t.Fatalf("expected 34.45, got %v", result)
	}
}

func TestGnataVariableBinding(t *testing.T) {
	statesVar := map[string]any{
		"input":  map[string]any{"foo": "bar"},
		"result": map[string]any{"status": "SUCCESS"},
	}

	expr, err := gnata.Compile(`$states.result.status`)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := expr.EvalWithVars(context.Background(), map[string]any{}, map[string]any{
		"states": statesVar,
	})
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if result != "SUCCESS" {
		t.Fatalf("expected SUCCESS, got %v", result)
	}
}

func TestGnataVariableBindingWithStateInput(t *testing.T) {
	data := map[string]any{"comment": "hello"}
	statesVar := map[string]any{
		"input":  data,
		"result": nil,
	}

	expr, err := gnata.Compile(`$states.input.comment`)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := expr.EvalWithVars(context.Background(), data, map[string]any{
		"states": statesVar,
	})
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if result != "hello" {
		t.Fatalf("expected hello, got %v", result)
	}
}

func TestGnataCustomFunctionRegistration(t *testing.T) {
	customFuncs := map[string]gnata.CustomFunc{
		"uuid": func(args []any, focus any) (any, error) {
			uuid := make([]byte, 16)
			if _, err := rand.Read(uuid); err != nil {
				return nil, err
			}
			uuid[6] = (uuid[6] & 0x0f) | 0x40
			uuid[8] = (uuid[8] & 0x3f) | 0x80
			return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
				uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
		},
		"partition": func(args []any, focus any) (any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("$partition requires 2 arguments")
			}
			arr, ok := args[0].([]any)
			if !ok {
				return nil, fmt.Errorf("$partition first argument must be an array")
			}
			var chunkSize float64
			switch v := args[1].(type) {
			case float64:
				chunkSize = v
			case int:
				chunkSize = float64(v)
			default:
				return nil, fmt.Errorf("$partition second argument must be a number")
			}
			if chunkSize < 1 {
				return nil, fmt.Errorf("$partition chunk size must be >= 1")
			}
			cs := int(chunkSize)
			var result []any
			for i := 0; i < len(arr); i += cs {
				end := i + cs
				if end > len(arr) {
					end = len(arr)
				}
				chunk := make([]any, end-i)
				copy(chunk, arr[i:end])
				result = append(result, chunk)
			}
			return result, nil
		},
		"range": func(args []any, focus any) (any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("$range requires at least 2 arguments")
			}
			toFloat := func(v any) (float64, error) {
				switch n := v.(type) {
				case float64:
					return n, nil
				case int:
					return float64(n), nil
				default:
					return 0, fmt.Errorf("expected number, got %T", v)
				}
			}
			start, err := toFloat(args[0])
			if err != nil {
				return nil, err
			}
			end, err := toFloat(args[1])
			if err != nil {
				return nil, err
			}
			delta := 1.0
			if len(args) >= 3 {
				delta, err = toFloat(args[2])
				if err != nil {
					return nil, err
				}
			}
			if delta == 0 {
				return nil, fmt.Errorf("$range delta must not be zero")
			}
			var result []any
			for v := start; (delta > 0 && v < end) || (delta < 0 && v > end); v += delta {
				result = append(result, v)
			}
			return result, nil
		},
		"hash": func(args []any, focus any) (any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("$hash requires 2 arguments")
			}
			s, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("$hash first argument must be a string")
			}
			alg, ok := args[1].(string)
			if !ok {
				return nil, fmt.Errorf("$hash second argument must be a string")
			}
			var hash []byte
			switch strings.ToUpper(alg) {
			case "MD5":
				h := md5.Sum([]byte(s))
				hash = h[:]
			case "SHA-1", "SHA1":
				h := sha1.Sum([]byte(s))
				hash = h[:]
			case "SHA-256", "SHA256":
				h := sha256.Sum256([]byte(s))
				hash = h[:]
			case "SHA-384", "SHA384":
				h := sha512.Sum384([]byte(s))
				hash = h[:]
			case "SHA-512", "SHA512":
				h := sha512.Sum512([]byte(s))
				hash = h[:]
			default:
				return nil, fmt.Errorf("$hash unsupported algorithm: %s", alg)
			}
			return hex.EncodeToString(hash), nil
		},
		"parse": func(args []any, focus any) (any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("$parse requires 1 argument")
			}
			s, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("$parse argument must be a string")
			}
			var result any
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				return nil, fmt.Errorf("$parse invalid JSON: %w", err)
			}
			return result, nil
		},
	}

	env := gnata.NewCustomEnv(customFuncs)

	tests := []struct {
		name     string
		expr     string
		data     any
		expected any
	}{
		{
			name:     "uuid returns string",
			expr:     `$uuid()`,
			data:     nil,
			expected: nil,
		},
		{
			name:     "partition basic",
			expr:     `$partition([1,2,3,4,5], 2)`,
			data:     nil,
			expected: []any{[]any{1.0, 2.0}, []any{3.0, 4.0}, []any{5.0}},
		},
		{
			name:     "partition exact",
			expr:     `$partition([1,2,3,4], 2)`,
			data:     nil,
			expected: []any{[]any{1.0, 2.0}, []any{3.0, 4.0}},
		},
		{
			name:     "range ascending",
			expr:     `$range(1, 5)`,
			data:     nil,
			expected: []any{1.0, 2.0, 3.0, 4.0},
		},
		{
			name:     "range with delta",
			expr:     `$range(0, 10, 3)`,
			data:     nil,
			expected: []any{0.0, 3.0, 6.0, 9.0},
		},
		{
			name:     "hash MD5",
			expr:     `$hash("hello", "MD5")`,
			data:     nil,
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "hash SHA-256",
			expr:     `$hash("hello", "SHA-256")`,
			data:     nil,
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:     "parse JSON string",
			expr:     `$parse('{"foo":"bar"}')`,
			data:     nil,
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:     "parse JSON array",
			expr:     `$parse('[1,2,3]')`,
			data:     nil,
			expected: []any{1.0, 2.0, 3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := gnata.Compile(tt.expr)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			result, err := expr.EvalWithCustomFuncs(context.Background(), tt.data, env)
			if err != nil {
				t.Fatalf("eval: %v", err)
			}
			if tt.name == "uuid returns string" {
				if _, ok := result.(string); !ok {
					t.Fatalf("expected string, got %T: %v", result, result)
				}
				if len(result.(string)) != 36 {
					t.Fatalf("expected UUID format (36 chars), got: %s", result)
				}
				return
			}
			if !deepEqual(result, tt.expected) {
				t.Fatalf("expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

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

func TestGnataWithVarsAndData(t *testing.T) {
	data := map[string]any{
		"inputVal": "test-input",
		"nested":   map[string]any{"key": "value"},
	}

	statesVar := map[string]any{
		"input":  data,
		"result": map[string]any{"output": "result-data"},
		"context": map[string]any{
			"Execution": map[string]any{
				"Id":   "arn:aws:states:us-east-1:123:execution:sm:exec1",
				"Name": "exec1",
			},
			"StateMachine": map[string]any{
				"Id":   "arn:aws:states:us-east-1:123:stateMachine:sm",
				"Name": "sm",
			},
			"State": map[string]any{
				"Name":        "TestState",
				"EnteredTime": "2025-01-01T00:00:00Z",
				"RetryCount":  0,
			},
		},
	}

	t.Run("access input via $states", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.input.inputVal`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithVars(context.Background(), data, map[string]any{
			"states": statesVar,
		})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "test-input" {
			t.Fatalf("expected test-input, got %v", result)
		}
	})

	t.Run("access result via $states", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.result.output`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithVars(context.Background(), data, map[string]any{
			"states": statesVar,
		})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "result-data" {
			t.Fatalf("expected result-data, got %v", result)
		}
	})

	t.Run("access context via $states", func(t *testing.T) {
		expr, err := gnata.Compile(`$states.context.State.Name`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithVars(context.Background(), data, map[string]any{
			"states": statesVar,
		})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if result != "TestState" {
			t.Fatalf("expected TestState, got %v", result)
		}
	})

	t.Run("combine data access and $states", func(t *testing.T) {
		expr, err := gnata.Compile(`{"dataKey": nested.key, "statesResult": $states.result.output}`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithVars(context.Background(), data, map[string]any{
			"states": statesVar,
		})
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		normalized := normalizeResult(result)
		m, ok := normalized.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T (raw: %T)", normalized, result)
		}
		if m["dataKey"] != "value" {
			t.Fatalf("expected dataKey=value, got %v", m["dataKey"])
		}
		if m["statesResult"] != "result-data" {
			t.Fatalf("expected statesResult=result-data, got %v", m["statesResult"])
		}
	})
}

func TestGnataCustomFuncsWithVarsWorkaround(t *testing.T) {
	data := map[string]any{"items": []any{1.0, 2.0, 3.0, 4.0, 5.0}}

	customFuncs := map[string]gnata.CustomFunc{
		"partition": func(args []any, focus any) (any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("$partition requires 2 arguments")
			}
			arr, ok := args[0].([]any)
			if !ok {
				return nil, fmt.Errorf("$partition first argument must be an array, got %T: %v", args[0], args[0])
			}
			var chunkSize float64
			switch v := args[1].(type) {
			case float64:
				chunkSize = v
			case int:
				chunkSize = float64(v)
			default:
				return nil, fmt.Errorf("$partition second argument must be a number")
			}
			cs := int(chunkSize)
			var result []any
			for i := 0; i < len(arr); i += cs {
				end := i + cs
				if end > len(arr) {
					end = len(arr)
				}
				chunk := make([]any, end-i)
				copy(chunk, arr[i:end])
				result = append(result, chunk)
			}
			return result, nil
		},
	}

	env := gnata.NewCustomEnv(customFuncs)

	t.Run("custom function with data as input", func(t *testing.T) {
		expr, err := gnata.Compile(`$partition(items, 2)`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		result, err := expr.EvalWithCustomFuncs(context.Background(), data, env)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		expected := []any{[]any{1.0, 2.0}, []any{3.0, 4.0}, []any{5.0}}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("custom functions and variables combined", func(t *testing.T) {
		expr, err := gnata.Compile(`$partition(items, 2)`)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		statesVar := map[string]any{"input": data}
		vars := map[string]any{"states": statesVar}
		result, err := expr.EvalWithEnvAndVars(context.Background(), data, env, vars)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		expected := []any{[]any{1.0, 2.0}, []any{3.0, 4.0}, []any{5.0}}
		if !deepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
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

func normalizeResult(v any) any {
	return gnata.NormalizeValue(v)
}

func deepEqual(a, b any) bool {
	ja, _ := json.Marshal(normalizeResult(a))
	jb, _ := json.Marshal(normalizeResult(b))
	return string(ja) == string(jb)
}

func generateLargePayload(targetSize int) map[string]any {
	items := make([]any, 0, 1000)
	for i := 0; i < 1000; i++ {
		items = append(items, map[string]any{
			"id":    fmt.Sprintf("item-%d", i),
			"value": float64(i),
			"data":  strings.Repeat("x", 200),
		})
	}
	payload := map[string]any{
		"id":    "large-payload",
		"items": items,
	}
	currentSize, _ := json.Marshal(payload)
	paddingNeeded := targetSize - len(currentSize)
	if paddingNeeded > 0 {
		payload["padding"] = strings.Repeat("p", paddingNeeded)
	}
	return payload
}

func TestGnataRandomFunction(t *testing.T) {
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		t.Fatalf("rand: %v", err)
	}
	_ = n

	seed := make([]byte, 8)
	if _, err := rand.Read(seed); err != nil {
		t.Fatalf("rand: %v", err)
	}
	_ = base64.StdEncoding.EncodeToString(seed)
}

func TestGnataIsFastPath(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		wantFast bool
	}{
		{"simple path", "foo.bar", true},
		{"dotted path", "a.b.c.d", true},
		{"function call", "$sum(items)", false},
		{"filter", "items[value > 0]", false},
		{"comparison", "status = 'ACTIVE'", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := gnata.Compile(tt.expr)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			if expr.IsFastPath() != tt.wantFast {
				t.Fatalf("IsFastPath(%q) = %v, want %v", tt.expr, expr.IsFastPath(), tt.wantFast)
			}
		})
	}
}
