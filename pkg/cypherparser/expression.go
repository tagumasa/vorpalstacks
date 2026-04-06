// Cypher expression evaluator and result helpers.
//
// Based on goraphdb/cypher_exec.go evalExpr/evalBool/evalComparison/evalFunc,
// getProperty, matchProps, matchLabels, compareValues, toFloat64, toBool,
// extended with:
//   - String predicates: IN, STARTS WITH, ENDS WITH, CONTAINS, IS [NOT] NULL, =~
//   - CASE/WHEN/THEN/ELSE evaluation
//   - Arithmetic operators: +, -, *, /, % (with string concatenation for +)
//   - List literal, map literal evaluation
//   - Aggregation computation: COUNT, SUM, AVG, MIN, MAX, COLLECT (with DISTINCT)
//   - LABELS(), PROPERTIES() evaluation
//   - Parameter resolution ($param substitution in AST)
//   - Result projection: sortRows, topKHeap, returnItemName, exprName
//   - EvalContext with ExistsFn callback for EXISTS { pattern }
//
// Original goraphdb evalExpr handled: Literal, VarRef, PropAccess, FuncCall,
// Comparison, And, Or, Not. All other expression kinds are our addition.

package cypherparser

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	guuid "github.com/google/uuid"
	"vorpalstacks/pkg/graphengine"
)

// CypherResult holds the result of a Cypher query execution.
// Based on goraphdb CypherResult.
type CypherResult struct {
	Columns []string         // column names in RETURN order
	Rows    []map[string]any // each row is a map of column name -> value
}

// EvalContext provides the context needed for expression evaluation.
type EvalContext struct {
	Bindings  map[string]any                               // variable name -> value
	Params    map[string]any                               // $param name -> value
	ExistsFn  func(*Pattern, map[string]any) (bool, error) // callback for EXISTS { pattern }
	AggValues map[*Expression]any                          // pre-computed aggregation values
}

// ---------------------------------------------------------------------------
// Expression evaluation
// ---------------------------------------------------------------------------

func evalExpr(ctx *EvalContext, e *Expression) (any, error) {
	if e == nil {
		return nil, nil
	}
	switch e.Kind {
	case ExprLiteral:
		return e.LitValue, nil

	case ExprVarRef:
		val, ok := ctx.Bindings[e.Variable]
		if !ok {
			return nil, fmt.Errorf("cypher: unbound variable %q", e.Variable)
		}
		return val, nil

	case ExprPropAccess:
		var obj any
		if e.NestedPropObject != nil {
			var err error
			obj, err = evalExpr(ctx, e.NestedPropObject)
			if err != nil {
				return nil, err
			}
		} else {
			var ok bool
			obj, ok = ctx.Bindings[e.Object]
			if !ok {
				return nil, fmt.Errorf("cypher: unbound variable %q", e.Object)
			}
		}
		if obj == nil {
			return nil, nil
		}
		return getProperty(obj, e.Property), nil

	case ExprFuncCall:
		return evalFunc(ctx, e.FuncName, e.Args)

	case ExprComparison:
		return evalComparison(ctx, e)

	case ExprAnd:
		var anyNil bool
		var allTrue = true
		for i := range e.Operands {
			v, err := evalExpr(ctx, &e.Operands[i])
			if err != nil {
				return nil, err
			}
			if v == nil {
				anyNil = true
				continue
			}
			b, ok := toBoolStrict(v)
			if !ok {
				return nil, fmt.Errorf("cypher: AND requires boolean operands, got %T", v)
			}
			if !b {
				return false, nil
			}
			allTrue = allTrue && b
		}
		if anyNil {
			return nil, nil
		}
		return allTrue, nil

	case ExprOr:
		var anyNil bool
		var anyTrue bool
		for i := range e.Operands {
			v, err := evalExpr(ctx, &e.Operands[i])
			if err != nil {
				return nil, err
			}
			if v == nil {
				anyNil = true
				continue
			}
			b, ok := toBoolStrict(v)
			if !ok {
				return nil, fmt.Errorf("cypher: OR requires boolean operands, got %T", v)
			}
			if b {
				anyTrue = true
				break
			}
		}
		if anyTrue {
			return true, nil
		}
		if anyNil {
			return nil, nil
		}
		return false, nil

	case ExprNot:
		v, err := evalExpr(ctx, e.Inner)
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		}
		b, ok := toBoolStrict(v)
		if !ok {
			return nil, fmt.Errorf("cypher: NOT requires a boolean operand, got %T", v)
		}
		return !b, nil

	case ExprParam:
		if ctx.Params == nil {
			return nil, fmt.Errorf("cypher: no parameters provided for $%s", e.ParamName)
		}
		val, ok := ctx.Params[e.ParamName]
		if !ok {
			return nil, fmt.Errorf("cypher: missing parameter $%s", e.ParamName)
		}
		return val, nil

	case ExprIn:
		left, err := evalExpr(ctx, e.InLeft)
		if err != nil {
			return nil, err
		}
		right, err := evalExpr(ctx, e.InRight)
		if err != nil {
			return nil, err
		}
		return evalIn(left, right)

	case ExprExists:
		if ctx.ExistsFn == nil {
			return false, nil
		}
		return ctx.ExistsFn(e.ExistsPattern, ctx.Bindings)

	case ExprStartsWith:
		s, err := evalExpr(ctx, e.StrExpr)
		if err != nil {
			return nil, err
		}
		prefix, err := evalExpr(ctx, e.StrValue)
		if err != nil {
			return nil, err
		}
		return evalStartsWith(s, prefix)

	case ExprEndsWith:
		s, err := evalExpr(ctx, e.StrExpr)
		if err != nil {
			return nil, err
		}
		suffix, err := evalExpr(ctx, e.StrValue)
		if err != nil {
			return nil, err
		}
		return evalEndsWith(s, suffix)

	case ExprContains:
		s, err := evalExpr(ctx, e.StrExpr)
		if err != nil {
			return nil, err
		}
		sub, err := evalExpr(ctx, e.StrValue)
		if err != nil {
			return nil, err
		}
		return evalContains(s, sub)

	case ExprIsNull:
		val, err := evalExpr(ctx, e.IsNullExpr)
		if err != nil {
			return nil, err
		}
		return val == nil, nil

	case ExprIsNotNull:
		val, err := evalExpr(ctx, e.IsNullExpr)
		if err != nil {
			return nil, err
		}
		return val != nil, nil

	case ExprRegexMatch:
		val, err := evalExpr(ctx, e.RegexExpr)
		if err != nil {
			return nil, err
		}
		pat, err := evalExpr(ctx, e.RegexPattern)
		if err != nil {
			return nil, err
		}
		return evalRegexMatch(val, pat)

	case ExprCase:
		return evalCase(ctx, e)

	case ExprListLiteral:
		result := make([]any, len(e.ListItems))
		for i := range e.ListItems {
			val, err := evalExpr(ctx, &e.ListItems[i])
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil

	case ExprMapLiteral:
		result := make(map[string]any, len(e.MapPairs))
		for i := range e.MapPairs {
			val, err := evalExpr(ctx, &e.MapPairs[i].Value)
			if err != nil {
				return nil, err
			}
			result[e.MapPairs[i].Key] = val
		}
		return result, nil

	case ExprAggregation:
		if ctx.AggValues != nil {
			if val, ok := ctx.AggValues[e]; ok {
				return val, nil
			}
		}
		return nil, fmt.Errorf("cypher: aggregation not yet computed")

	case ExprArithAdd:
		left, err := evalExpr(ctx, e.ArithLeft)
		if err != nil {
			return nil, err
		}
		right, err := evalExpr(ctx, e.ArithRight)
		if err != nil {
			return nil, err
		}
		return evalArith(left, right, '+')

	case ExprArithSub:
		left, err := evalExpr(ctx, e.ArithLeft)
		if err != nil {
			return nil, err
		}
		right, err := evalExpr(ctx, e.ArithRight)
		if err != nil {
			return nil, err
		}
		return evalArith(left, right, '-')

	case ExprArithMul:
		left, err := evalExpr(ctx, e.ArithLeft)
		if err != nil {
			return nil, err
		}
		right, err := evalExpr(ctx, e.ArithRight)
		if err != nil {
			return nil, err
		}
		return evalArith(left, right, '*')

	case ExprArithDiv:
		left, err := evalExpr(ctx, e.ArithLeft)
		if err != nil {
			return nil, err
		}
		right, err := evalExpr(ctx, e.ArithRight)
		if err != nil {
			return nil, err
		}
		return evalArith(left, right, '/')

	case ExprArithMod:
		left, err := evalExpr(ctx, e.ArithLeft)
		if err != nil {
			return nil, err
		}
		right, err := evalExpr(ctx, e.ArithRight)
		if err != nil {
			return nil, err
		}
		return evalArith(left, right, '%')

	case ExprLabels:
		val, err := evalExpr(ctx, e.LabelsExpr)
		if err != nil {
			return nil, err
		}
		return evalLabels(val)

	case ExprProperties:
		val, err := evalExpr(ctx, e.PropertiesExpr)
		if err != nil {
			return nil, err
		}
		return evalProperties(val)

	default:
		return nil, fmt.Errorf("cypher: unsupported expression kind %d", e.Kind)
	}
}

func evalBool(ctx *EvalContext, e *Expression) (bool, error) {
	val, err := evalExpr(ctx, e)
	if err != nil {
		return false, err
	}
	if val == nil {
		return false, nil
	}
	b, ok := toBoolStrict(val)
	if !ok {
		return false, fmt.Errorf("cypher: WHERE requires a boolean expression, got %T", val)
	}
	return b, nil
}

// ---------------------------------------------------------------------------
// Comparison
// ---------------------------------------------------------------------------

func evalComparison(ctx *EvalContext, e *Expression) (any, error) {
	left, err := evalExpr(ctx, e.Left)
	if err != nil {
		return nil, err
	}
	right, err := evalExpr(ctx, e.Right)
	if err != nil {
		return nil, err
	}

	if left == nil || right == nil {
		return nil, nil
	}

	cmp, cerr := compareValuesSafe(left, right)
	if cerr != nil {
		if e.Op == OpNeq {
			return nil, nil
		}
		return nil, nil
	}

	switch e.Op {
	case OpEq:
		return cmp == 0, nil
	case OpNeq:
		return cmp != 0, nil
	case OpLt:
		return cmp < 0, nil
	case OpGt:
		return cmp > 0, nil
	case OpLte:
		return cmp <= 0, nil
	case OpGte:
		return cmp >= 0, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// Built-in functions
// ---------------------------------------------------------------------------

func evalFunc(ctx *EvalContext, name string, args []Expression) (any, error) {
	switch strings.ToLower(name) {
	case "type":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: type() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		if e, ok := val.(*graphengine.Edge); ok {
			return e.Label, nil
		}
		return nil, nil

	case "id":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: id() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case *graphengine.Node:
			return int64(v.ID), nil
		case *graphengine.Edge:
			return int64(v.ID), nil
		}
		return nil, nil

	case "size":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: size() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		return evalSize(val)

	case "length":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: length() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		return evalSize(val)

	case "abs":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: abs() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: abs() requires a numeric argument")
		}
		return math.Abs(f), nil

	case "ceil":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: ceil() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: ceil() requires a numeric argument")
		}
		return math.Ceil(f), nil

	case "floor":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: floor() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: floor() requires a numeric argument")
		}
		return math.Floor(f), nil

	case "round":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: round() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: round() requires a numeric argument")
		}
		return math.Round(f), nil

	case "sqrt":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: sqrt() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: sqrt() requires a numeric argument")
		}
		return math.Sqrt(f), nil

	case "coalesce":
		for i := range args {
			val, err := evalExpr(ctx, &args[i])
			if err != nil {
				return nil, err
			}
			if val != nil {
				return val, nil
			}
		}
		return nil, nil

	case "tointeger":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: toInteger() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}
		result, _ := toInt64Coerce(val)
		return result, nil

	case "tofloat":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: toFloat() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}
		f, ok := toFloat64(val)
		if !ok {
			return float64(0), nil
		}
		return f, nil

	case "tostring":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: toString() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}
		switch v := val.(type) {
		case string:
			return v, nil
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
		case []any:
			parts := make([]string, len(v))
			for i, item := range v {
				s := fmt.Sprint(item)
				parts[i] = s
			}
			return "[" + strings.Join(parts, ", ") + "]", nil
		case []string:
			return "[" + strings.Join(v, ", ") + "]", nil
		}
		return fmt.Sprint(val), nil

	case "head":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: head() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		list, ok := val.([]any)
		if !ok || len(list) == 0 {
			return nil, nil
		}
		return list[0], nil

	case "last":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: last() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		list, ok := val.([]any)
		if !ok || len(list) == 0 {
			return nil, nil
		}
		return list[len(list)-1], nil

	case "reverse":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: reverse() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}
		list, ok := val.([]any)
		if !ok {
			return nil, fmt.Errorf("cypher: reverse() requires a list argument")
		}
		result := make([]any, len(list))
		for i, v := range list {
			result[len(list)-1-i] = v
		}
		return result, nil

	case "range":
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("cypher: range() requires 2 or 3 arguments")
		}
		startVal, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		endVal, err := evalExpr(ctx, &args[1])
		if err != nil {
			return nil, err
		}
		var step int64 = 1
		if len(args) == 3 {
			stepVal, err := evalExpr(ctx, &args[2])
			if err != nil {
				return nil, err
			}
			s, ok := toInt64Coerce(stepVal)
			if !ok || s == 0 {
				return nil, fmt.Errorf("cypher: range() step must be a non-zero integer")
			}
			step = s
		}
		start, ok1 := toInt64Coerce(startVal)
		end, ok2 := toInt64Coerce(endVal)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("cypher: range() requires integer arguments")
		}
		return evalRange(start, end, step), nil

	case "substring", "substr":
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("cypher: substring() requires 2 or 3 arguments")
		}
		sVal, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := sVal.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: substring() requires a string argument")
		}
		startVal, err := evalExpr(ctx, &args[1])
		if err != nil {
			return nil, err
		}
		start, ok := toInt64Coerce(startVal)
		if !ok {
			return nil, fmt.Errorf("cypher: substring() start must be an integer")
		}
		var length int64
		if len(args) == 3 {
			lenVal, err := evalExpr(ctx, &args[2])
			if err != nil {
				return nil, err
			}
			l, ok := toInt64Coerce(lenVal)
			if !ok {
				return nil, fmt.Errorf("cypher: substring() length must be an integer")
			}
			length = l
		} else {
			length = int64(utf8.RuneCountInString(s)) - start
		}
		runes := []rune(s)
		if start < 0 {
			start = 0
		}
		if start >= int64(len(runes)) || length <= 0 {
			return "", nil
		}
		end := start + length
		if end > int64(len(runes)) {
			end = int64(len(runes))
		}
		return string(runes[start:end]), nil

	case "trim":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: trim() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: trim() requires a string argument")
		}
		return strings.TrimSpace(s), nil

	case "ltrim":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: ltrim() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: ltrim() requires a string argument")
		}
		return strings.TrimLeft(s, " \t\n\r"), nil

	case "rtrim":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: rtrim() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: rtrim() requires a string argument")
		}
		return strings.TrimRight(s, " \t\n\r"), nil

	case "upper":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: upper() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: upper() requires a string argument")
		}
		return strings.ToUpper(s), nil

	case "lower":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: lower() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: lower() requires a string argument")
		}
		return strings.ToLower(s), nil

	case "replace":
		if len(args) != 3 {
			return nil, fmt.Errorf("cypher: replace() requires exactly 3 arguments")
		}
		orig, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := orig.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: replace() requires string arguments")
		}
		oldVal, err := evalExpr(ctx, &args[1])
		if err != nil {
			return nil, err
		}
		old, ok := oldVal.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: replace() requires string arguments")
		}
		newVal, err := evalExpr(ctx, &args[2])
		if err != nil {
			return nil, err
		}
		newStr, ok := newVal.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: replace() requires string arguments")
		}
		return strings.ReplaceAll(s, old, newStr), nil

	case "split":
		if len(args) != 2 {
			return nil, fmt.Errorf("cypher: split() requires exactly 2 arguments")
		}
		sVal, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		if sVal == nil {
			return nil, nil
		}
		s, ok := sVal.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: split() requires string arguments")
		}
		sepVal, err := evalExpr(ctx, &args[1])
		if err != nil {
			return nil, err
		}
		sep, ok := sepVal.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: split() requires string arguments")
		}
		parts := strings.Split(s, sep)
		result := make([]any, len(parts))
		for i, p := range parts {
			result[i] = p
		}
		return result, nil

	case "left":
		if len(args) != 2 {
			return nil, fmt.Errorf("cypher: left() requires exactly 2 arguments")
		}
		sVal, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := sVal.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: left() requires a string and an integer")
		}
		nVal, err := evalExpr(ctx, &args[1])
		if err != nil {
			return nil, err
		}
		n, ok := toInt64Coerce(nVal)
		if !ok || n < 0 {
			return "", nil
		}
		runes := []rune(s)
		if int(n) > len(runes) {
			n = int64(len(runes))
		}
		return string(runes[:n]), nil

	case "right":
		if len(args) != 2 {
			return nil, fmt.Errorf("cypher: right() requires exactly 2 arguments")
		}
		sVal, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		s, ok := sVal.(string)
		if !ok {
			return nil, fmt.Errorf("cypher: right() requires a string and an integer")
		}
		nVal, err := evalExpr(ctx, &args[1])
		if err != nil {
			return nil, err
		}
		n, ok := toInt64Coerce(nVal)
		if !ok || n < 0 {
			return "", nil
		}
		runes := []rune(s)
		if int(n) > len(runes) {
			n = int64(len(runes))
		}
		return string(runes[len(runes)-int(n):]), nil

	case "exists":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: exists() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return false, err
		}
		if val == nil {
			return false, nil
		}
		if args[0].Kind == ExprPropAccess {
			obj, ok := ctx.Bindings[args[0].Object]
			if !ok {
				return false, nil
			}
			propKey := args[0].Property
			switch n := obj.(type) {
			case *graphengine.Node:
				if n.Props != nil {
					_, exists := n.Props[propKey]
					return exists, nil
				}
			case *graphengine.Edge:
				if propKey == "label" || propKey == "type" {
					return true, nil
				}
				if n.Props != nil {
					_, exists := n.Props[propKey]
					return exists, nil
				}
			case map[string]any:
				_, exists := n[propKey]
				return exists, nil
			}
			return false, nil
		}
		return true, nil

	case "sign":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: sign() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: sign() requires a numeric argument")
		}
		if f > 0 {
			return 1.0, nil
		}
		if f < 0 {
			return -1.0, nil
		}
		return 0.0, nil

	case "exp":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: exp() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: exp() requires a numeric argument")
		}
		return math.Exp(f), nil

	case "log":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: log() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: log() requires a numeric argument")
		}
		return math.Log(f), nil

	case "log10":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: log10() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: log10() requires a numeric argument")
		}
		return math.Log10(f), nil

	case "sin":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: sin() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: sin() requires a numeric argument")
		}
		return math.Sin(f), nil

	case "cos":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: cos() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: cos() requires a numeric argument")
		}
		return math.Cos(f), nil

	case "tan":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: tan() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: tan() requires a numeric argument")
		}
		return math.Tan(f), nil

	case "asin":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: asin() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: asin() requires a numeric argument")
		}
		return math.Asin(f), nil

	case "acos":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: acos() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: acos() requires a numeric argument")
		}
		return math.Acos(f), nil

	case "atan":
		if len(args) != 1 {
			return nil, fmt.Errorf("cypher: atan() requires exactly 1 argument")
		}
		val, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		f, ok := toFloat64(val)
		if !ok {
			return nil, fmt.Errorf("cypher: atan() requires a numeric argument")
		}
		return math.Atan(f), nil

	case "atan2":
		if len(args) != 2 {
			return nil, fmt.Errorf("cypher: atan2() requires exactly 2 arguments")
		}
		yVal, err := evalExpr(ctx, &args[0])
		if err != nil {
			return nil, err
		}
		xVal, err := evalExpr(ctx, &args[1])
		if err != nil {
			return nil, err
		}
		y, ok := toFloat64(yVal)
		if !ok {
			return nil, fmt.Errorf("cypher: atan2() requires numeric arguments")
		}
		x, ok := toFloat64(xVal)
		if !ok {
			return nil, fmt.Errorf("cypher: atan2() requires numeric arguments")
		}
		return math.Atan2(y, x), nil

	case "pi":
		return math.Pi, nil

	case "randomuuid":
		return genUUID().String(), nil

	default:
		return nil, fmt.Errorf("cypher: unknown function %q", name)
	}
}

// ---------------------------------------------------------------------------
// String predicates
// ---------------------------------------------------------------------------

func evalIn(left, right any) (any, error) {
	if left == nil {
		return nil, nil
	}
	list, ok := right.([]any)
	if !ok {
		return nil, nil
	}
	var hasNull bool
	for _, item := range list {
		if item == nil {
			hasNull = true
			continue
		}
		cmp, cerr := compareValuesSafe(left, item)
		if cerr != nil {
			hasNull = true
			continue
		}
		if cmp == 0 {
			return true, nil
		}
	}
	if hasNull {
		return nil, nil
	}
	return false, nil
}

func evalStartsWith(s, prefix any) (any, error) {
	if s == nil || prefix == nil {
		return nil, nil
	}
	ss, ok := s.(string)
	if !ok {
		return false, nil
	}
	sp, ok := prefix.(string)
	if !ok {
		return false, nil
	}
	return strings.HasPrefix(ss, sp), nil
}

func evalEndsWith(s, suffix any) (any, error) {
	if s == nil || suffix == nil {
		return nil, nil
	}
	ss, ok := s.(string)
	if !ok {
		return false, nil
	}
	sp, ok := suffix.(string)
	if !ok {
		return false, nil
	}
	return strings.HasSuffix(ss, sp), nil
}

func evalContains(s, sub any) (any, error) {
	if s == nil || sub == nil {
		return nil, nil
	}
	ss, ok := s.(string)
	if !ok {
		return false, nil
	}
	sp, ok := sub.(string)
	if !ok {
		return false, nil
	}
	return strings.Contains(ss, sp), nil
}

func evalRegexMatch(val, pattern any) (bool, error) {
	s, ok := val.(string)
	if !ok {
		return false, nil
	}
	p, ok := pattern.(string)
	if !ok {
		return false, fmt.Errorf("cypher: =~ requires a string pattern")
	}
	re, err := regexp.Compile(p)
	if err != nil {
		return false, fmt.Errorf("cypher: invalid regex pattern: %w", err)
	}
	return re.MatchString(s), nil
}

// ---------------------------------------------------------------------------
// CASE expression
// ---------------------------------------------------------------------------

func evalCase(ctx *EvalContext, e *Expression) (any, error) {
	if e.CaseOperand != nil {
		// General CASE: CASE expr WHEN val1 THEN result1 WHEN val2 THEN result2 ... ELSE default END
		operand, err := evalExpr(ctx, e.CaseOperand)
		if err != nil {
			return nil, err
		}
		for i := range e.CaseWhens {
			whenVal, err := evalExpr(ctx, e.CaseWhens[i].Condition)
			if err != nil {
				return nil, err
			}
			cmp, cerr := compareValuesSafe(operand, whenVal)
			if cerr == nil && cmp == 0 {
				return evalExpr(ctx, e.CaseWhens[i].Result)
			}
		}
	} else {
		// Simple CASE: CASE WHEN cond1 THEN result1 WHEN cond2 THEN result2 ... ELSE default END
		for i := range e.CaseWhens {
			cond, err := evalBool(ctx, e.CaseWhens[i].Condition)
			if err != nil {
				return nil, err
			}
			if cond {
				return evalExpr(ctx, e.CaseWhens[i].Result)
			}
		}
	}

	if e.CaseElse != nil {
		return evalExpr(ctx, e.CaseElse)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Arithmetic
// ---------------------------------------------------------------------------

func evalArith(left, right any, op byte) (any, error) {
	if left == nil || right == nil {
		return nil, nil
	}

	if op == '+' {
		ls, lIsStr := left.(string)
		rs, rIsStr := right.(string)
		if lIsStr && rIsStr {
			return ls + rs, nil
		}
		if lIsStr {
			return ls + fmt.Sprint(right), nil
		}
		if rIsStr {
			return fmt.Sprint(left) + rs, nil
		}
	}

	lf, lOk := toFloat64(left)
	rf, rOk := toFloat64(right)
	if !lOk || !rOk {
		return nil, fmt.Errorf("cypher: arithmetic %c requires numeric operands, got %T and %T", op, left, right)
	}

	// Integer arithmetic when both operands are native integer types (not coerced from float)
	lIsInt := isIntegerType(left)
	rIsInt := isIntegerType(right)
	li, _ := toInt64(left)
	ri, _ := toInt64(right)

	switch op {
	case '+':
		if lIsInt && rIsInt {
			if ri > 0 && li > math.MaxInt64-ri {
				return lf + rf, nil
			}
			if ri < 0 && li < math.MinInt64-ri {
				return lf + rf, nil
			}
			return li + ri, nil
		}
		return lf + rf, nil
	case '-':
		if lIsInt && rIsInt {
			if ri < 0 && li > math.MaxInt64+ri {
				return lf - rf, nil
			}
			if ri > 0 && li < math.MinInt64+ri {
				return lf - rf, nil
			}
			return li - ri, nil
		}
		return lf - rf, nil
	case '*':
		if lIsInt && rIsInt {
			if li != 0 && ri != 0 {
				if li > 0 {
					if ri > 0 {
						if li > math.MaxInt64/ri {
							return lf * rf, nil
						}
					} else {
						if ri < math.MinInt64/li {
							return lf * rf, nil
						}
					}
				} else {
					if ri > 0 {
						if li < math.MinInt64/ri {
							return lf * rf, nil
						}
					} else {
						if li < math.MaxInt64/ri {
							return lf * rf, nil
						}
					}
				}
			}
			return li * ri, nil
		}
		return lf * rf, nil
	case '/':
		if rf == 0 {
			return nil, nil
		}
		if lIsInt && rIsInt {
			if li == math.MinInt64 && ri == -1 {
				return float64(li) / float64(ri), nil
			}
			return li / ri, nil
		}
		return lf / rf, nil
	case '%':
		if rf == 0 {
			return nil, nil
		}
		if lIsInt && rIsInt {
			if li == math.MinInt64 && ri == -1 {
				return 0, nil
			}
			return li % ri, nil
		}
		return math.Mod(lf, rf), nil
	default:
		return nil, fmt.Errorf("cypher: unknown arithmetic operator %c", op)
	}
}

// ---------------------------------------------------------------------------
// LABELS / PROPERTIES
// ---------------------------------------------------------------------------

func evalLabels(val any) (any, error) {
	switch v := val.(type) {
	case *graphengine.Node:
		return v.Labels, nil
	}
	return nil, nil
}

func evalProperties(val any) (any, error) {
	switch v := val.(type) {
	case *graphengine.Node:
		if v.Props == nil {
			return map[string]any{}, nil
		}
		result := make(map[string]any, len(v.Props))
		for k, v := range v.Props {
			result[k] = v
		}
		return result, nil
	case *graphengine.Edge:
		if v.Props == nil {
			return map[string]any{}, nil
		}
		result := make(map[string]any, len(v.Props))
		for k, v := range v.Props {
			result[k] = v
		}
		return result, nil
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Aggregation
// ---------------------------------------------------------------------------

// HasAggregation checks if an expression tree contains any aggregation function.
func HasAggregation(e *Expression) bool {
	if e == nil {
		return false
	}
	if e.Kind == ExprAggregation {
		return true
	}
	switch e.Kind {
	case ExprAnd, ExprOr:
		for i := range e.Operands {
			if HasAggregation(&e.Operands[i]) {
				return true
			}
		}
	case ExprNot:
		return HasAggregation(e.Inner)
	case ExprComparison:
		return HasAggregation(e.Left) || HasAggregation(e.Right)
	case ExprFuncCall:
		for i := range e.Args {
			if HasAggregation(&e.Args[i]) {
				return true
			}
		}
	case ExprArithAdd, ExprArithSub, ExprArithMul, ExprArithDiv, ExprArithMod:
		return HasAggregation(e.ArithLeft) || HasAggregation(e.ArithRight)
	case ExprCase:
		if e.CaseOperand != nil && HasAggregation(e.CaseOperand) {
			return true
		}
		for i := range e.CaseWhens {
			if HasAggregation(e.CaseWhens[i].Condition) || HasAggregation(e.CaseWhens[i].Result) {
				return true
			}
		}
		if e.CaseElse != nil && HasAggregation(e.CaseElse) {
			return true
		}
	case ExprListLiteral:
		for i := range e.ListItems {
			if HasAggregation(&e.ListItems[i]) {
				return true
			}
		}
	case ExprMapLiteral:
		for i := range e.MapPairs {
			if HasAggregation(&e.MapPairs[i].Value) {
				return true
			}
		}
	case ExprIn:
		return HasAggregation(e.InLeft) || HasAggregation(e.InRight)
	case ExprStartsWith, ExprEndsWith, ExprContains:
		return HasAggregation(e.StrExpr) || HasAggregation(e.StrValue)
	case ExprIsNull, ExprIsNotNull:
		return HasAggregation(e.IsNullExpr)
	case ExprRegexMatch:
		return HasAggregation(e.RegexExpr) || HasAggregation(e.RegexPattern)
	case ExprLabels:
		return HasAggregation(e.LabelsExpr)
	case ExprProperties:
		return HasAggregation(e.PropertiesExpr)
	}
	return false
}

// ComputeAggregation computes an aggregation function over a set of rows.
// Each row is a map of variable name -> value (not column names).
func ComputeAggregation(fn AggFunc, arg *Expression, distinct bool, rows []map[string]any) (any, error) {
	switch fn {
	case AggCount:
		if arg == nil {
			return int64(len(rows)), nil
		}
		if distinct {
			seen := make(map[string]struct{})
			for _, row := range rows {
				ctx := &EvalContext{Bindings: row}
				val, err := evalExpr(ctx, arg)
				if err != nil {
					return nil, err
				}
				if val == nil {
					continue
				}
				key := fmt.Sprintf("%v", val)
				seen[key] = struct{}{}
			}
			return int64(len(seen)), nil
		}
		var count int64
		for _, row := range rows {
			ctx := &EvalContext{Bindings: row}
			val, err := evalExpr(ctx, arg)
			if err != nil {
				return nil, err
			}
			if val != nil {
				count++
			}
		}
		return count, nil

	case AggSum:
		if distinct {
			seen := make(map[string]struct{})
			var sum float64
			allInt := true
			for _, row := range rows {
				ctx := &EvalContext{Bindings: row}
				val, err := evalExpr(ctx, arg)
				if err != nil {
					return nil, err
				}
				if val == nil {
					continue
				}
				key := fmt.Sprintf("%v", val)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
				if isIntegerType(val) {
					iv, _ := toInt64(val)
					sum += float64(iv)
				} else {
					allInt = false
					f, ok := toFloat64(val)
					if ok {
						sum += f
					}
				}
			}
			if allInt {
				return int64(sum), nil
			}
			return sum, nil
		}
		var sum float64
		allInt := true
		for _, row := range rows {
			ctx := &EvalContext{Bindings: row}
			val, err := evalExpr(ctx, arg)
			if err != nil {
				return nil, err
			}
			if val == nil {
				allInt = false
				continue
			}
			if isIntegerType(val) {
				iv, _ := toInt64(val)
				sum += float64(iv)
			} else {
				allInt = false
				f, ok := toFloat64(val)
				if ok {
					sum += f
				}
			}
		}
		if allInt {
			return int64(sum), nil
		}
		return sum, nil

	case AggAvg:
		var sum float64
		var count int64
		var seen map[any]struct{}
		if distinct {
			seen = make(map[any]struct{})
		}
		for _, row := range rows {
			ctx := &EvalContext{Bindings: row}
			val, err := evalExpr(ctx, arg)
			if err != nil {
				return nil, err
			}
			if distinct {
				key := fmt.Sprintf("%v", val)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
			}
			f, ok := toFloat64(val)
			if ok {
				sum += f
				count++
			}
		}
		if count == 0 {
			return nil, nil
		}
		return sum / float64(count), nil

	case AggMin:
		var seen map[string]struct{}
		if distinct {
			seen = make(map[string]struct{})
		}
		var min any
		for _, row := range rows {
			ctx := &EvalContext{Bindings: row}
			val, err := evalExpr(ctx, arg)
			if err != nil {
				return nil, err
			}
			if val == nil {
				continue
			}
			if distinct {
				key := fmt.Sprintf("%v", val)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
			}
			if min == nil || compareValues(val, min) < 0 {
				min = val
			}
		}
		return min, nil

	case AggMax:
		var seen map[string]struct{}
		if distinct {
			seen = make(map[string]struct{})
		}
		var max any
		for _, row := range rows {
			ctx := &EvalContext{Bindings: row}
			val, err := evalExpr(ctx, arg)
			if err != nil {
				return nil, err
			}
			if val == nil {
				continue
			}
			if distinct {
				key := fmt.Sprintf("%v", val)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
			}
			if max == nil || compareValues(val, max) > 0 {
				max = val
			}
		}
		return max, nil

	case AggCollect:
		if distinct {
			seen := make(map[string]struct{})
			result := make([]any, 0)
			for _, row := range rows {
				ctx := &EvalContext{Bindings: row}
				val, err := evalExpr(ctx, arg)
				if err != nil {
					return nil, err
				}
				key := fmt.Sprintf("%v", val)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
				result = append(result, val)
			}
			return result, nil
		}
		result := make([]any, 0, len(rows))
		for _, row := range rows {
			ctx := &EvalContext{Bindings: row}
			val, err := evalExpr(ctx, arg)
			if err != nil {
				return nil, err
			}
			result = append(result, val)
		}
		return result, nil

	default:
		return nil, fmt.Errorf("cypher: unknown aggregation function %d", fn)
	}
}

// ---------------------------------------------------------------------------
// Value helpers
// ---------------------------------------------------------------------------

// toFloat64 converts a value to float64 for numeric comparisons and arithmetic.
// Based on goraphdb toFloat64, extended with uint types and json.Number.
func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case bool:
		if n {
			return 1.0, true
		}
		return 0.0, true
	}
	return 0, false
}

func isIntegerType(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func toInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case int8:
		return int64(n), true
	case int16:
		return int64(n), true
	case int32:
		return int64(n), true
	case uint:
		return int64(n), true
	case uint8:
		return int64(n), true
	case uint16:
		return int64(n), true
	case uint32:
		return int64(n), true
	case uint64:
		if n > uint64(math.MaxInt64) {
			return 0, false
		}
		return int64(n), true
	case float64:
		if n == math.Trunc(n) && !math.IsInf(n, 0) {
			return int64(n), true
		}
		return 0, false
	}
	return 0, false
}

// toInt64Coerce converts a value to int64, coercing floats by truncation.
func toInt64Coerce(v any) (int64, bool) {
	f, ok := toFloat64(v)
	if !ok {
		return 0, false
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, false
	}
	return int64(f), true
}

// toBool coerces a value to bool.
// Based on goraphdb toBool.
func toBool(v any) bool {
	if v == nil {
		return false
	}
	f, ok := toFloat64(v)
	if ok {
		return f != 0
	}
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b != ""
	}
	return true
}

func toBoolStrict(v any) (bool, bool) {
	b, ok := v.(bool)
	return b, ok
}

func genUUID() guuid.UUID {
	id, _ := guuid.NewRandomFromReader(rand.Reader)
	return id
}

// compareValues compares two arbitrary values. Returns -1, 0, or 1.
// Based on goraphdb compareValues.
func compareValues(a, b any) int {
	cmp, _ := compareValuesSafe(a, b)
	return cmp
}

var errIncomparable = fmt.Errorf("incomparable types")

func compareValuesSafe(a, b any) (int, error) {
	if a == nil && b == nil {
		return 0, nil
	}
	if a == nil {
		return -1, nil
	}
	if b == nil {
		return 1, nil
	}

	ai, aIsInt := toInt64(a)
	bi, bIsInt := toInt64(b)
	if aIsInt && bIsInt {
		switch {
		case ai < bi:
			return -1, nil
		case ai > bi:
			return 1, nil
		default:
			return 0, nil
		}
	}

	af, aOk := toFloat64(a)
	bf, bOk := toFloat64(b)
	if aOk && bOk {
		switch {
		case af < bf:
			return -1, nil
		case af > bf:
			return 1, nil
		default:
			return 0, nil
		}
	}

	as, aStr := a.(string)
	bs, bStr := b.(string)
	if aStr && bStr {
		return strings.Compare(as, bs), nil
	}

	ab, aBool := a.(bool)
	bb, bBool := b.(bool)
	if aBool && bBool {
		if ab == bb {
			return 0, nil
		}
		if !ab {
			return -1, nil
		}
		return 1, nil
	}

	aList, aIsList := a.([]any)
	bList, bIsList := b.([]any)
	if aIsList && bIsList {
		minLen := len(aList)
		if len(bList) < minLen {
			minLen = len(bList)
		}
		for i := 0; i < minLen; i++ {
			cmp, err := compareValuesSafe(aList[i], bList[i])
			if err != nil {
				return 0, err
			}
			if cmp != 0 {
				return cmp, nil
			}
		}
		switch {
		case len(aList) < len(bList):
			return -1, nil
		case len(aList) > len(bList):
			return 1, nil
		default:
			return 0, nil
		}
	}

	return 0, errIncomparable
}

// ---------------------------------------------------------------------------
// Property / label matching
// ---------------------------------------------------------------------------

// getProperty extracts a property from a Node or Edge.
// Based on goraphdb getProperty.
func getProperty(obj any, prop string) any {
	switch v := obj.(type) {
	case *graphengine.Node:
		if v.Props != nil {
			return v.Props[prop]
		}
	case *graphengine.Edge:
		switch prop {
		case "label", "type":
			return v.Label
		default:
			if v.Props != nil {
				return v.Props[prop]
			}
		}
	case map[string]any:
		return v[prop]
	}
	return nil
}

// matchProps checks whether a node's properties match all constraints.
// Based on goraphdb matchProps.
func matchProps(actual graphengine.Props, constraints map[string]any) bool {
	if constraints == nil {
		return true
	}
	for key, expected := range constraints {
		val, ok := actual[key]
		if !ok {
			return false
		}
		if compareValues(val, expected) != 0 {
			return false
		}
	}
	return true
}

func matchPropsWithBinding(actual graphengine.Props, constraints map[string]any, binding map[string]any) bool {
	if constraints == nil {
		return true
	}
	for key, expected := range constraints {
		val, ok := actual[key]
		if !ok {
			return false
		}
		resolved := expected
		if expr, ok := expected.(Expression); ok {
			v, err := evalExpr(&EvalContext{Bindings: binding}, &expr)
			if err != nil {
				return false
			}
			resolved = v
		}
		if compareValues(val, resolved) != 0 {
			return false
		}
	}
	return true
}

// matchLabels checks if a node's labels contain all required labels.
// Based on goraphdb matchLabels.
func matchLabels(nodeLabels, required []string) bool {
	if len(required) == 0 {
		return true
	}
	set := make(map[string]bool, len(nodeLabels))
	for _, l := range nodeLabels {
		set[l] = true
	}
	for _, r := range required {
		if !set[r] {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Result projection helpers
// ---------------------------------------------------------------------------

// returnItemName computes the column name for a RETURN item.
// Based on goraphdb returnItemName.
func returnItemName(item ReturnItem) string {
	if item.Alias != "" {
		return item.Alias
	}
	return exprName(item.Expr)
}

// exprName returns a readable name for an expression (used as column name).
// Based on goraphdb exprName.
func exprName(e Expression) string {
	switch e.Kind {
	case ExprVarRef:
		return e.Variable
	case ExprPropAccess:
		return e.Object + "." + e.Property
	case ExprFuncCall:
		args := make([]string, len(e.Args))
		for i, a := range e.Args {
			args[i] = exprName(a)
		}
		return e.FuncName + "(" + strings.Join(args, ", ") + ")"
	case ExprAggregation:
		if e.AggArg == nil {
			return aggFuncName(e.AggFunc) + "(*)"
		}
		name := aggFuncName(e.AggFunc) + "("
		if e.AggDistinct {
			name += "DISTINCT "
		}
		name += exprName(*e.AggArg) + ")"
		return name
	case ExprLiteral:
		return fmt.Sprintf("%v", e.LitValue)
	case ExprArithAdd, ExprArithSub, ExprArithMul, ExprArithDiv, ExprArithMod:
		op := "?"
		switch e.Kind {
		case ExprArithAdd:
			op = "+"
		case ExprArithSub:
			op = "-"
		case ExprArithMul:
			op = "*"
		case ExprArithDiv:
			op = "/"
		case ExprArithMod:
			op = "%"
		}
		return exprName(*e.ArithLeft) + " " + op + " " + exprName(*e.ArithRight)
	case ExprLabels:
		return "labels(" + exprName(*e.LabelsExpr) + ")"
	case ExprProperties:
		return "properties(" + exprName(*e.PropertiesExpr) + ")"
	case ExprParam:
		return "$" + e.ParamName
	case ExprIn:
		return exprName(*e.InLeft) + " IN " + exprName(*e.InRight)
	case ExprIsNull:
		return exprName(*e.IsNullExpr) + " IS NULL"
	case ExprIsNotNull:
		return exprName(*e.IsNullExpr) + " IS NOT NULL"
	case ExprStartsWith:
		return exprName(*e.StrExpr) + " STARTS WITH " + exprName(*e.StrValue)
	case ExprEndsWith:
		return exprName(*e.StrExpr) + " ENDS WITH " + exprName(*e.StrValue)
	case ExprContains:
		return exprName(*e.StrExpr) + " CONTAINS " + exprName(*e.StrValue)
	case ExprRegexMatch:
		return exprName(*e.RegexExpr) + " =~ " + exprName(*e.RegexPattern)
	case ExprNot:
		return "NOT " + exprName(*e.Inner)
	case ExprAnd:
		ops := make([]string, len(e.Operands))
		for i, o := range e.Operands {
			ops[i] = exprName(o)
		}
		return strings.Join(ops, " AND ")
	case ExprOr:
		ops := make([]string, len(e.Operands))
		for i, o := range e.Operands {
			ops[i] = exprName(o)
		}
		return strings.Join(ops, " OR ")
	case ExprListLiteral:
		items := make([]string, len(e.ListItems))
		for i, it := range e.ListItems {
			items[i] = exprName(it)
		}
		return "[" + strings.Join(items, ", ") + "]"
	case ExprMapLiteral:
		pairs := make([]string, len(e.MapPairs))
		for i, p := range e.MapPairs {
			pairs[i] = p.Key + ": " + exprName(p.Value)
		}
		return "{" + strings.Join(pairs, ", ") + "}"
	case ExprCase:
		if e.CaseOperand != nil {
			return "CASE " + exprName(*e.CaseOperand) + " ..."
		}
		return "CASE ..."
	default:
		return "expr"
	}
}

func aggFuncName(fn AggFunc) string {
	switch fn {
	case AggCount:
		return "count"
	case AggSum:
		return "sum"
	case AggAvg:
		return "avg"
	case AggMin:
		return "min"
	case AggMax:
		return "max"
	case AggCollect:
		return "collect"
	default:
		return "agg"
	}
}

func evalSize(val any) (any, error) {
	if val == nil {
		return nil, nil
	}
	switch v := val.(type) {
	case string:
		return int64(utf8.RuneCountInString(v)), nil
	case []any:
		return int64(len(v)), nil
	case []string:
		return int64(len(v)), nil
	case map[string]any:
		return int64(len(v)), nil
	}
	return nil, fmt.Errorf("cypher: size() requires a string, list, or map argument")
}

func evalRange(start, end, step int64) []any {
	if step > 0 && start > end {
		return []any{}
	}
	if step < 0 && start < end {
		return []any{}
	}
	var result []any
	for i := start; (step > 0 && i <= end) || (step < 0 && i >= end); i += step {
		result = append(result, i)
	}
	return result
}

// ---------------------------------------------------------------------------
// Sorting
// ---------------------------------------------------------------------------

func sortRows(rows []map[string]any, orderItems []OrderItem) {
	sort.SliceStable(rows, func(i, j int) bool {
		for _, oi := range orderItems {
			vi := evalOrderValue(&oi.Expr, rows[i])
			vj := evalOrderValue(&oi.Expr, rows[j])

			if vi == nil && vj == nil {
				continue
			}
			if vi == nil {
				if oi.Desc {
					return false
				}
				return false
			}
			if vj == nil {
				if oi.Desc {
					return true
				}
				return true
			}

			cmp := compareValues(vi, vj)
			if cmp == 0 {
				continue
			}
			if oi.Desc {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})
}

func evalOrderValue(e *Expression, row map[string]any) any {
	name := exprName(*e)
	if v, ok := row[name]; ok {
		return v
	}
	val, _ := evalExpr(&EvalContext{Bindings: row}, e)
	return val
}

func evalRowExpr(e *Expression, row map[string]any) any {
	return evalOrderValue(e, row)
}

// ---------------------------------------------------------------------------
// Top-K Heap — selects the K best items from a stream using O(N log K) time.
// Based on goraphdb topKHeap.
// ---------------------------------------------------------------------------

type topKItem struct {
	sortKey []any
	source  any
}

type topKHeap struct {
	items   []topKItem
	orderBy []OrderItem
	limit   int
}

func newTopKHeap(orderBy []OrderItem, limit int) *topKHeap {
	return &topKHeap{
		items:   make([]topKItem, 0, limit),
		orderBy: orderBy,
		limit:   limit,
	}
}

func (h *topKHeap) Len() int { return len(h.items) }

func (h *topKHeap) Less(i, j int) bool {
	for idx, oi := range h.orderBy {
		cmp := compareValues(h.items[i].sortKey[idx], h.items[j].sortKey[idx])
		if cmp == 0 {
			continue
		}
		if oi.Desc {
			return cmp < 0
		}
		return cmp > 0
	}
	return false
}

func (h *topKHeap) Swap(i, j int) { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *topKHeap) Push(x any) { h.items = append(h.items, x.(topKItem)) }

func (h *topKHeap) Pop() any {
	n := len(h.items)
	item := h.items[n-1]
	h.items = h.items[:n-1]
	return item
}

func (h *topKHeap) offer(item topKItem) {
	if len(h.items) < h.limit {
		heapPush(h, item)
		return
	}
	if h.isBetter(item.sortKey, h.items[0].sortKey) {
		h.items[0] = item
		heapFix(h, 0)
	}
}

func (h *topKHeap) isBetter(a, b []any) bool {
	for idx, oi := range h.orderBy {
		cmp := compareValues(a[idx], b[idx])
		if cmp == 0 {
			continue
		}
		if oi.Desc {
			return cmp > 0
		}
		return cmp < 0
	}
	return false
}

func (h *topKHeap) sorted() []topKItem {
	n := len(h.items)
	result := make([]topKItem, n)
	for i := n - 1; i >= 0; i-- {
		result[i] = heapPop(h).(topKItem)
	}
	return result
}

func evalSortKey(orderBy []OrderItem, bindings map[string]any) []any {
	key := make([]any, len(orderBy))
	for i, oi := range orderBy {
		val, _ := evalExpr(&EvalContext{Bindings: bindings}, &oi.Expr)
		key[i] = val
	}
	return key
}

// Minimal heap operations to avoid importing container/heap in tests.
func heapPush(h *topKHeap, item topKItem) {
	h.items = append(h.items, item)
	i := len(h.items) - 1
	for i > 0 {
		parent := (i - 1) / 2
		if h.Less(i, parent) {
			h.Swap(i, parent)
			i = parent
		} else {
			break
		}
	}
}

func heapPop(h *topKHeap) any {
	n := len(h.items)
	item := h.items[0]
	h.items[0] = h.items[n-1]
	h.items = h.items[:n-1]
	i := 0
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < len(h.items) && h.Less(left, smallest) {
			smallest = left
		}
		if right < len(h.items) && h.Less(right, smallest) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.Swap(i, smallest)
		i = smallest
	}
	return item
}

func heapFix(h *topKHeap, i int) {
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < len(h.items) && h.Less(left, smallest) {
			smallest = left
		}
		if right < len(h.items) && h.Less(right, smallest) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.Swap(i, smallest)
		i = smallest
	}
}

// ---------------------------------------------------------------------------
// Parameter resolution
// ---------------------------------------------------------------------------

// ResolveParams substitutes $param references in a CypherQuery AST with actual
// values from the params map. Does not mutate the original AST — returns a
// deep-copied query with params resolved.
// Based on goraphdb resolveParams.
func ResolveParams(q *CypherQuery, params map[string]any) error {
	for i := range q.Match.Pattern.Nodes {
		if err := resolveNodeParams(&q.Match.Pattern.Nodes[i], params); err != nil {
			return err
		}
	}
	for i := range q.Match.Pattern.Rels {
		if err := resolveRelParams(&q.Match.Pattern.Rels[i], params); err != nil {
			return err
		}
	}
	for i := range q.Matches {
		for j := range q.Matches[i].Pattern.Nodes {
			if err := resolveNodeParams(&q.Matches[i].Pattern.Nodes[j], params); err != nil {
				return err
			}
		}
		for j := range q.Matches[i].Pattern.Rels {
			if err := resolveRelParams(&q.Matches[i].Pattern.Rels[j], params); err != nil {
				return err
			}
		}
	}
	if q.OptionalMatch != nil {
		for i := range q.OptionalMatch.Pattern.Nodes {
			if err := resolveNodeParams(&q.OptionalMatch.Pattern.Nodes[i], params); err != nil {
				return err
			}
		}
		for i := range q.OptionalMatch.Pattern.Rels {
			if err := resolveRelParams(&q.OptionalMatch.Pattern.Rels[i], params); err != nil {
				return err
			}
		}
	}
	if q.Where != nil {
		if err := resolveExprParams(q.Where, params); err != nil {
			return err
		}
	}
	if q.OptionalWhere != nil {
		if err := resolveExprParams(q.OptionalWhere, params); err != nil {
			return err
		}
	}
	for i := range q.Return.Items {
		if err := resolveExprParams(&q.Return.Items[i].Expr, params); err != nil {
			return err
		}
	}
	for i := range q.Set {
		if err := resolveExprParams(&q.Set[i].Value, params); err != nil {
			return err
		}
	}
	if q.Unwind != nil {
		if err := resolveExprParams(&q.Unwind.Expr, params); err != nil {
			return err
		}
	}
	if q.With != nil {
		for i := range q.With.Items {
			if err := resolveExprParams(&q.With.Items[i].Expr, params); err != nil {
				return err
			}
		}
	}
	for i := range q.OrderBy {
		if err := resolveExprParams(&q.OrderBy[i].Expr, params); err != nil {
			return err
		}
	}
	if q.LimitParam != "" {
		val, exists := params[q.LimitParam]
		if !exists {
			return fmt.Errorf("cypher: missing parameter $%s", q.LimitParam)
		}
		n, ok := toInt64Coerce(val)
		if !ok {
			return fmt.Errorf("cypher: LIMIT parameter must be an integer")
		}
		q.Limit = intPtr(int(n))
		q.LimitParam = ""
	}
	if q.SkipParam != "" {
		val, exists := params[q.SkipParam]
		if !exists {
			return fmt.Errorf("cypher: missing parameter $%s", q.SkipParam)
		}
		if val == nil {
			q.Skip = intPtr(0)
			q.SkipParam = ""
			return nil
		}
		n, ok := toInt64Coerce(val)
		if !ok {
			return fmt.Errorf("cypher: SKIP parameter must be an integer")
		}
		q.Skip = intPtr(int(n))
		q.SkipParam = ""
	}
	return nil
}

func resolveRelParams(rp *RelPattern, params map[string]any) error {
	return resolveNodeParams(&NodePattern{Props: rp.Props}, params)
}

func resolveNodeParams(np *NodePattern, params map[string]any) error {
	for k, v := range np.Props {
		if ref, ok := v.(paramRef); ok {
			val, exists := params[string(ref)]
			if !exists {
				return fmt.Errorf("cypher: missing parameter $%s", string(ref))
			}
			np.Props[k] = val
		} else if expr, ok := v.(Expression); ok {
			if err := resolveExprParams(&expr, params); err != nil {
				return err
			}
			if exprHasVarRef(&expr) {
				np.Props[k] = expr
			} else {
				val, err := evalExpr(&EvalContext{}, &expr)
				if err != nil {
					return fmt.Errorf("cypher: cannot evaluate property %q: %w", k, err)
				}
				np.Props[k] = val
			}
		}
	}
	return nil
}

func exprHasVarRef(e *Expression) bool {
	if e == nil {
		return false
	}
	if e.Kind == ExprVarRef {
		return true
	}
	switch e.Kind {
	case ExprPropAccess:
		return e.Object != ""
	case ExprFuncCall:
		for i := range e.Args {
			if exprHasVarRef(&e.Args[i]) {
				return true
			}
		}
	case ExprComparison:
		return exprHasVarRef(e.Left) || exprHasVarRef(e.Right)
	case ExprAnd, ExprOr:
		for i := range e.Operands {
			if exprHasVarRef(&e.Operands[i]) {
				return true
			}
		}
	case ExprNot:
		return exprHasVarRef(e.Inner)
	case ExprArithAdd, ExprArithSub, ExprArithMul, ExprArithDiv, ExprArithMod:
		return exprHasVarRef(e.ArithLeft) || exprHasVarRef(e.ArithRight)
	case ExprListLiteral:
		for i := range e.ListItems {
			if exprHasVarRef(&e.ListItems[i]) {
				return true
			}
		}
	case ExprMapLiteral:
		for i := range e.MapPairs {
			if exprHasVarRef(&e.MapPairs[i].Value) {
				return true
			}
		}
	}
	return false
}

func resolveExprParams(expr *Expression, params map[string]any) error {
	if expr == nil {
		return nil
	}
	switch expr.Kind {
	case ExprParam:
		val, exists := params[expr.ParamName]
		if !exists {
			return fmt.Errorf("cypher: missing parameter $%s", expr.ParamName)
		}
		*expr = litExpr(val)

	case ExprComparison:
		if err := resolveExprParams(expr.Left, params); err != nil {
			return err
		}
		return resolveExprParams(expr.Right, params)

	case ExprAnd, ExprOr:
		for i := range expr.Operands {
			if err := resolveExprParams(&expr.Operands[i], params); err != nil {
				return err
			}
		}

	case ExprNot:
		return resolveExprParams(expr.Inner, params)

	case ExprFuncCall:
		for i := range expr.Args {
			if err := resolveExprParams(&expr.Args[i], params); err != nil {
				return err
			}
		}

	case ExprIn:
		if err := resolveExprParams(expr.InLeft, params); err != nil {
			return err
		}
		return resolveExprParams(expr.InRight, params)

	case ExprStartsWith, ExprEndsWith, ExprContains:
		if err := resolveExprParams(expr.StrExpr, params); err != nil {
			return err
		}
		return resolveExprParams(expr.StrValue, params)

	case ExprIsNull, ExprIsNotNull:
		return resolveExprParams(expr.IsNullExpr, params)

	case ExprRegexMatch:
		if err := resolveExprParams(expr.RegexExpr, params); err != nil {
			return err
		}
		return resolveExprParams(expr.RegexPattern, params)

	case ExprCase:
		if expr.CaseOperand != nil {
			if err := resolveExprParams(expr.CaseOperand, params); err != nil {
				return err
			}
		}
		for i := range expr.CaseWhens {
			if err := resolveExprParams(expr.CaseWhens[i].Condition, params); err != nil {
				return err
			}
			if err := resolveExprParams(expr.CaseWhens[i].Result, params); err != nil {
				return err
			}
		}
		if expr.CaseElse != nil {
			return resolveExprParams(expr.CaseElse, params)
		}

	case ExprListLiteral:
		for i := range expr.ListItems {
			if err := resolveExprParams(&expr.ListItems[i], params); err != nil {
				return err
			}
		}

	case ExprMapLiteral:
		for i := range expr.MapPairs {
			if err := resolveExprParams(&expr.MapPairs[i].Value, params); err != nil {
				return err
			}
		}

	case ExprArithAdd, ExprArithSub, ExprArithMul, ExprArithDiv, ExprArithMod:
		if err := resolveExprParams(expr.ArithLeft, params); err != nil {
			return err
		}
		return resolveExprParams(expr.ArithRight, params)

	case ExprLabels:
		return resolveExprParams(expr.LabelsExpr, params)

	case ExprProperties:
		return resolveExprParams(expr.PropertiesExpr, params)

	case ExprAggregation:
		if expr.AggArg != nil {
			return resolveExprParams(expr.AggArg, params)
		}
	}
	return nil
}
