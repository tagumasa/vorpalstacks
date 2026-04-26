// Cypher string predicates, CASE expressions, arithmetic, LABELS/PROPERTIES, and aggregation.
//
// Based on goraphdb/cypher_exec.go, extended with string predicates (IN, STARTS WITH,
// ENDS WITH, CONTAINS, =~), CASE/WHEN/THEN/ELSE, arithmetic operators, and aggregation
// computation (COUNT, SUM, AVG, MIN, MAX, COLLECT with DISTINCT).

package cypherparser

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"vorpalstacks/internal/core/storage/graphengine"
)

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

func evalCase(ctx *EvalContext, e *Expression) (any, error) {
	if e.CaseOperand != nil {
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
