// Cypher parameter resolution ($param substitution in AST).
//
// Walks the parsed CypherQuery and replaces ParameterRef expressions with
// concrete values from the params map. Also resolves expression parameters
// in node/relationship patterns.

package cypherparser

import "fmt"

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
