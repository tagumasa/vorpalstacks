// Cypher value conversion helpers, property/label matching, and result naming utilities.
//
// Based on goraphdb/cypher_exec.go getProperty, matchProps, matchLabels, compareValues,
// toFloat64, toBool. Extended with type coercion, parameter resolution support, and
// result projection helpers (returnItemName, exprName, aggFuncName, evalSize, evalRange).

package cypherparser

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	guuid "github.com/google/uuid"
	"vorpalstacks/pkg/graphengine"
)

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

func returnItemName(item ReturnItem) string {
	if item.Alias != "" {
		return item.Alias
	}
	return exprName(item.Expr)
}

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
