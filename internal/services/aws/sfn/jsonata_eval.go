package sfn

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	gnata "github.com/recolabs/gnata"
)

// awsCustomFuncs registers AWS-specific custom functions for JSONata evaluation.
// These provide additional functionality beyond the standard JSONata library:
// $partition, $range, $hash, $random, $uuid, and $parse.
var awsCustomFuncs = map[string]gnata.CustomFunc{
	"partition": awsPartitionFunc,
	"range":     awsRangeFunc,
	"hash":      awsHashFunc,
	"random":    awsRandomFunc,
	"uuid":      awsUUIDFunc,
	"parse":     awsParseFunc,
}

// awsCustomEnv is the pre-built JSONata custom environment containing AWS custom functions.
var awsCustomEnv any

func init() {
	awsCustomEnv = gnata.NewCustomEnv(awsCustomFuncs)
}

// IsExpression reports whether a string is a JSONata inline expression (wrapped in {% ... %}).
func IsExpression(s string) bool {
	return strings.HasPrefix(s, "{%") && strings.HasSuffix(s, "%}")
}

// UnwrapExpression strips the {% and %} delimiters from a JSONata inline expression.
func UnwrapExpression(s string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(s, "%}"), "{%"))
}

// EvaluateJSONata compiles and evaluates a JSONata expression against the given data.
func EvaluateJSONata(ctx context.Context, expression string, data interface{}, vars map[string]interface{}) (interface{}, error) {
	expr, err := gnata.Compile(expression)
	if err != nil {
		return nil, fmt.Errorf("compile: %w", err)
	}
	return evaluateCompiled(ctx, expr, data, vars)
}

func evaluateCompiled(ctx context.Context, expr *gnata.Expression, data interface{}, vars map[string]interface{}) (interface{}, error) {
	result, err := expr.EvalWithEnvAndVars(ctx, data, awsCustomEnv, vars)
	if err != nil {
		return nil, err
	}
	normalized := gnata.NormalizeValue(result)
	if normalized != nil && reflect.TypeOf(normalized).Kind() == reflect.Func {
		return nil, nil
	}
	return normalized, nil
}

// ResolveTemplate recursively resolves JSONata inline expressions within strings,
// maps, and arrays against the given data.
func ResolveTemplate(ctx context.Context, value interface{}, data interface{}, vars map[string]interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		if IsExpression(v) {
			exprStr := UnwrapExpression(v)
			return EvaluateJSONata(ctx, exprStr, data, vars)
		}
		return v, nil
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			resolved, err := ResolveTemplate(ctx, val, data, vars)
			if err != nil {
				return nil, err
			}
			result[key] = resolved
		}
		return result, nil
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			resolved, err := ResolveTemplate(ctx, val, data, vars)
			if err != nil {
				return nil, err
			}
			result[i] = resolved
		}
		return result, nil
	default:
		return v, nil
	}
}

// BuildStatesVar constructs the $states context object used in JSONata state path expressions.
// Contains input, result, optional errorOutput, and optional context object.
func BuildStatesVar(input, result, errorOutput, contextObj interface{}) map[string]interface{} {
	states := map[string]interface{}{
		"input":  input,
		"result": result,
	}
	if errorOutput != nil {
		states["errorOutput"] = errorOutput
	}
	if contextObj != nil {
		states["context"] = contextObj
	}
	return map[string]interface{}{"states": states}
}

// IsJSONataExpressionValue reports whether a value is a JSONata inline expression string.
func IsJSONataExpressionValue(v interface{}) bool {
	s, ok := v.(string)
	return ok && IsExpression(s)
}

// EvaluateExpressionValue evaluates a value; if it is a JSONata expression string,
// it is compiled and evaluated against the given data.
func EvaluateExpressionValue(ctx context.Context, v interface{}, data interface{}, vars map[string]interface{}) (interface{}, error) {
	s, ok := v.(string)
	if !ok {
		return v, nil
	}
	if IsExpression(s) {
		return EvaluateJSONata(ctx, UnwrapExpression(s), data, vars)
	}
	return s, nil
}

// NormalizeResult normalises a JSONata result value using the gnata normaliser.
func NormalizeResult(v interface{}) interface{} {
	return gnata.NormalizeValue(v)
}

func awsUUIDFunc(args []interface{}, focus interface{}) (interface{}, error) {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return nil, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

func awsPartitionFunc(args []interface{}, focus interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("$partition requires 2 arguments")
	}
	arr, ok := args[0].([]interface{})
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
	var result []interface{}
	for i := 0; i < len(arr); i += cs {
		end := i + cs
		if end > len(arr) {
			end = len(arr)
		}
		chunk := make([]interface{}, end-i)
		copy(chunk, arr[i:end])
		result = append(result, chunk)
	}
	return result, nil
}

func awsRangeFunc(args []interface{}, focus interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("$range requires at least 2 arguments")
	}
	toFloat := func(v interface{}) (float64, error) {
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
	var result []interface{}
	for v := start; (delta > 0 && v < end) || (delta < 0 && v > end); v += delta {
		result = append(result, v)
	}
	return result, nil
}

func awsHashFunc(args []interface{}, focus interface{}) (interface{}, error) {
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
}

func awsRandomFunc(args []interface{}, focus interface{}) (interface{}, error) {
	if len(args) > 0 {
		if seed, ok := toFloat64(args[0]); ok {
			n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			if err != nil {
				return nil, err
			}
			seeded := float64(n.Int64()) / float64(math.MaxInt64)
			return (seeded + seed) - math.Floor(seeded+seed), nil
		}
	}
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	return float64(n.Int64()) / float64(math.MaxInt64), nil
}

func awsParseFunc(args []interface{}, focus interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("$parse requires 1 argument")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("$parse argument must be a string")
	}
	var result interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, fmt.Errorf("$parse invalid JSON: %w", err)
	}
	return result, nil
}
