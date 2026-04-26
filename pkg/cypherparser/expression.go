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
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"vorpalstacks/internal/core/storage/graphengine"
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
