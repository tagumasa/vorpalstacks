package rules

import (
	"fmt"
	"math"
	"strings"
)

// FunctionHandler evaluates a built-in SQL function call.
// The evaluator is provided so handlers can evaluate their own arguments
// (lazily, with error propagation).
type FunctionHandler func(e *Evaluator, args []Expr) (interface{}, error)

var functionHandlers map[string]FunctionHandler

func init() {
	functionHandlers = map[string]FunctionHandler{
		"topic":     fnTopic,
		"timestamp": fnTimestamp,
		"clientid":  fnClientID,
		"concat":    fnConcat,
		"length":    fnLength,
		"upper":     fnUpper,
		"lower":     fnLower,
		"trim":      fnTrim,
		"replace":   fnReplace,
		"substring": fnSubstring,
		"cast":      fnCast,
		"abs":       fnAbs,
		"ceil":      fnCeil,
		"floor":     fnFloor,
	}
}

func fnTopic(e *Evaluator, _ []Expr) (interface{}, error)     { return e.Topic, nil }
func fnTimestamp(e *Evaluator, _ []Expr) (interface{}, error) { return e.Timestamp.Unix(), nil }
func fnClientID(e *Evaluator, _ []Expr) (interface{}, error)  { return e.ClientID, nil }

func fnConcat(e *Evaluator, args []Expr) (interface{}, error) {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		val, err := e.Eval(arg)
		if err != nil {
			return nil, err
		}
		parts = append(parts, toString(val))
	}
	return strings.Join(parts, ""), nil
}

func fnLength(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("evaluator: length() requires 1 argument")
	}
	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case string:
		return float64(len(v)), nil
	case []interface{}:
		return float64(len(v)), nil
	case map[string]interface{}:
		return float64(len(v)), nil
	default:
		return float64(len(toString(v))), nil
	}
}

func fnUpper(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 1 {
		return nil, nil
	}
	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	return strings.ToUpper(toString(val)), nil
}

func fnLower(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 1 {
		return nil, nil
	}
	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	return strings.ToLower(toString(val)), nil
}

func fnTrim(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 1 {
		return nil, nil
	}
	val, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}
	return strings.TrimSpace(toString(val)), nil
}

func fnReplace(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 3 {
		return nil, nil
	}
	str, _ := e.Eval(args[0])
	old, _ := e.Eval(args[1])
	repl, _ := e.Eval(args[2])
	return strings.ReplaceAll(toString(str), toString(old), toString(repl)), nil
}

func fnSubstring(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 2 {
		return nil, nil
	}
	str, _ := e.Eval(args[0])
	startVal, _ := e.Eval(args[1])
	start := int(toFloat(startVal)) - 1
	s := toString(str)
	length := len(s) - start
	if len(args) >= 3 {
		lenVal, _ := e.Eval(args[2])
		length = int(toFloat(lenVal))
	}
	if start < 0 || start >= len(s) {
		return "", nil
	}
	if start+length > len(s) {
		return s[start:], nil
	}
	if length < 0 {
		return s[start:], nil
	}
	return s[start : start+length], nil
}

func fnCast(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 2 {
		return nil, nil
	}
	val, _ := e.Eval(args[0])
	typeName := strings.ToUpper(fmt.Sprintf("%v", evalExprName(args[1])))
	return e.castValue(val, typeName)
}

func fnAbs(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 1 {
		return nil, nil
	}
	val, _ := e.Eval(args[0])
	return math.Abs(toFloat(val)), nil
}

func fnCeil(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 1 {
		return nil, nil
	}
	val, _ := e.Eval(args[0])
	return math.Ceil(toFloat(val)), nil
}

func fnFloor(e *Evaluator, args []Expr) (interface{}, error) {
	if len(args) < 1 {
		return nil, nil
	}
	val, _ := e.Eval(args[0])
	return math.Floor(toFloat(val)), nil
}
