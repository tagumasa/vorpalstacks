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
