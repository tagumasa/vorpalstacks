package sfn

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
	"reflect"
	"strings"
	"time"

	gnata "github.com/recolabs/gnata"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// awsCustomFuncs registers AWS-specific custom functions for JSONata evaluation.
// These provide additional functionality beyond the standard JSONata library:
// $partition, $range, $hash, $random, $uuid, and $parse. The eval override
// shadows the stdlib binding (customs bind after RegisterAll) because AWS
// removes the function: "$eval is not available—use $parse instead".
var awsCustomFuncs = map[string]gnata.CustomFunc{
	"partition": awsPartitionFunc,
	"range":     awsRangeFunc,
	"hash":      awsHashFunc,
	"random":    awsRandomFunc,
	"uuid":      awsUUIDFunc,
	"parse":     awsParseFunc,
	"eval":      awsEvalDisabledFunc,
}

// awsEvalDisabledFunc replaces the JSONata stdlib $eval for AWS JSONata
// states.
func awsEvalDisabledFunc(args []interface{}, focus interface{}) (interface{}, error) {
	return nil, errors.New("$eval is not available—use $parse instead")
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

// jsonataEvalTimeout bounds a single JSONata expression evaluation: "A
// JSONata expression that takes longer than 1 second to evaluate will fail
// with an Expression evaluation timeout error." It is a variable so tests
// can shrink the budget; production always runs at the documented second.
var jsonataEvalTimeout = time.Second

// EvaluateJSONata compiles and evaluates a JSONata expression against the given data.
func EvaluateJSONata(ctx context.Context, expression string, data interface{}, vars map[string]interface{}) (interface{}, error) {
	expr, err := gnata.Compile(expression)
	if err != nil {
		return nil, fmt.Errorf("compile: %w", err)
	}
	evalCtx, cancel := context.WithTimeout(ctx, jsonataEvalTimeout)
	defer cancel()
	result, err := evaluateCompiled(evalCtx, expr, data, vars)
	if evalCtx.Err() != nil && ctx.Err() == nil {
		// The one-second expression budget expired; an outer
		// cancellation (a task or machine timeout) keeps its own
		// error class.
		return nil, errors.New("Expression evaluation timeout")
	}
	return result, err
}

// errUndefinedResult marks a JSONata expression that evaluated to undefined.
// JSON cannot represent an undefined value, so the state fails; every call
// site wraps this message as States.QueryEvaluationError via
// newQueryEvalError.
var errUndefinedResult = errors.New("the expression returned an undefined result")

// jsonNullSentinel is gnata's JSON null value. The evaluator distinguishes
// it from Go nil (JSONata undefined), but the platform decodes state JSON
// with encoding/json, which maps null to Go nil — so nested nils are
// converted to the sentinel before evaluation to keep input nulls distinct
// from undefined paths.
var jsonNullSentinel = func() interface{} {
	v, err := gnata.DecodeJSON([]byte("null"))
	if err != nil {
		panic("gnata: decoding the null literal failed: " + err.Error())
	}
	return v
}()

// nullsToSentinels rebuilds v with every nested Go nil replaced by the JSON
// null sentinel. Containers are copied only when a conversion happened;
// otherwise the original value is returned unchanged.
func nullsToSentinels(v interface{}) (interface{}, bool) {
	switch val := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		changed := false
		for k, mv := range val {
			if mv == nil {
				out[k] = jsonNullSentinel
				changed = true
				continue
			}
			cv, c := nullsToSentinels(mv)
			out[k] = cv
			changed = changed || c
		}
		if !changed {
			return val, false
		}
		return out, true
	case []interface{}:
		out := make([]interface{}, len(val))
		changed := false
		for i, ev := range val {
			if ev == nil {
				out[i] = jsonNullSentinel
				changed = true
				continue
			}
			cv, c := nullsToSentinels(ev)
			out[i] = cv
			changed = changed || c
		}
		if !changed {
			return val, false
		}
		return out, true
	}
	return v, false
}

func evaluateCompiled(ctx context.Context, expr *gnata.Expression, data interface{}, vars map[string]interface{}) (interface{}, error) {
	evalData, _ := nullsToSentinels(data)
	evalVars := vars
	if vars != nil {
		evalVars = make(map[string]interface{}, len(vars))
		for k, v := range vars {
			cv, _ := nullsToSentinels(v)
			evalVars[k] = cv
		}
	}
	result, err := expr.EvalWithEnvAndVars(ctx, evalData, awsCustomEnv, evalVars)
	if err != nil {
		return nil, err
	}
	// JSONata undefined arrives as Go nil, distinct from the JSON null
	// sentinel (which NormalizeValue maps to nil only below): JSON cannot
	// represent an undefined value, so it fails the state.
	if result == nil {
		return nil, errUndefinedResult
	}
	normalized := gnata.NormalizeValue(result)
	if normalized != nil && reflect.TypeOf(normalized).Kind() == reflect.Func {
		// A function reference has no JSON representation either.
		return nil, errUndefinedResult
	}
	return normalized, nil
}

// queryEvalSecondsValue validates a JSONata-evaluated seconds field
// (TimeoutSeconds, HeartbeatSeconds): the field accepts an integer within
// 1..MaxWaitSeconds, and any other type or range fails the state as a
// query-evaluation error.
func queryEvalSecondsValue(field string, result interface{}) (int32, error) {
	f, ok := toFloat64(result)
	if !ok {
		return 0, fmt.Errorf("%s requires a numeric value, got %v", field, result)
	}
	if f != math.Trunc(f) {
		return 0, fmt.Errorf("%s requires an integer value, got %v", field, result)
	}
	if f < 1 || f > sfnstore.MaxWaitSeconds {
		return 0, fmt.Errorf("%s value %v is outside the acceptable range 1-%d", field, result, sfnstore.MaxWaitSeconds)
	}
	return int32(f), nil
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

// The AWS context functions with a States intrinsic twin delegate to the
// single intrinsic implementation (applyIntrinsic); the wrappers only
// adapt the gnata function-registry signature.
func awsUUIDFunc(args []interface{}, focus interface{}) (interface{}, error) {
	return applyIntrinsic("UUID", nil)
}

func awsPartitionFunc(args []interface{}, focus interface{}) (interface{}, error) {
	return applyIntrinsic("ArrayPartition", args)
}

func awsRangeFunc(args []interface{}, focus interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("$range requires at least 2 arguments")
	}
	toRangeFloat := func(v interface{}) (float64, error) {
		f, ok := toFloat64(v)
		if !ok {
			return 0, fmt.Errorf("expected number, got %T", v)
		}
		return f, nil
	}
	start, err := toRangeFloat(args[0])
	if err != nil {
		return nil, err
	}
	end, err := toRangeFloat(args[1])
	if err != nil {
		return nil, err
	}
	delta := 1.0
	if len(args) >= 3 {
		delta, err = toRangeFloat(args[2])
		if err != nil {
			return nil, err
		}
	}
	if delta == 0 {
		return nil, fmt.Errorf("$range delta must not be zero")
	}
	// $range is the JSONata equivalent of States.ArrayRange: elements run
	// from start until the end value is reached or exceeded (inclusive),
	// bounded by the documented element limit.
	var result []interface{}
	for v := start; (delta > 0 && v <= end) || (delta < 0 && v >= end); v += delta {
		result = append(result, v)
		if len(result) > sfnstore.MaxArrayRangeElements {
			return nil, fmt.Errorf("$range result exceeds %d elements", sfnstore.MaxArrayRangeElements)
		}
	}
	return result, nil
}

func awsHashFunc(args []interface{}, focus interface{}) (interface{}, error) {
	return applyIntrinsic("Hash", args)
}

func awsRandomFunc(args []interface{}, focus interface{}) (interface{}, error) {
	// "The function takes an optional integer argument representing the
	// seed value of the random function. If you use this function with
	// the same seed value, it returns an identical number."
	if len(args) >= 1 {
		seed, ok := toFloat64(args[0])
		if !ok {
			return nil, fmt.Errorf("$random seed must be a number, got %T", args[0])
		}
		return mrand.New(mrand.NewSource(int64(seed))).Float64(), nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	return float64(n.Int64()) / float64(math.MaxInt64), nil
}

func awsParseFunc(args []interface{}, focus interface{}) (interface{}, error) {
	return applyIntrinsic("StringToJson", args)
}
