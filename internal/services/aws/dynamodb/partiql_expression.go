package dynamodb

import (
	"regexp"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

func filterItemsByExpr(items []*dbstore.Item, expr sqlparser.Expr, params *partiQLParams) []*dbstore.Item {
	var result []*dbstore.Item

	for _, item := range items {
		if evaluateExpr(item.Attributes, expr, params) {
			result = append(result, item)
		}
	}

	return result
}

func evaluateExpr(attrs map[string]*dbstore.AttributeValue, expr sqlparser.Expr, params *partiQLParams) bool {
	switch e := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return evaluateComparison(attrs, e, params)
	case *sqlparser.AndExpr:
		return evaluateExpr(attrs, e.Left, params) && evaluateExpr(attrs, e.Right, params)
	case *sqlparser.OrExpr:
		return evaluateExpr(attrs, e.Left, params) || evaluateExpr(attrs, e.Right, params)
	case *sqlparser.NotExpr:
		return !evaluateExpr(attrs, e.Expr, params)
	case *sqlparser.IsExpr:
		return evaluateIsExpr(attrs, e, params)
	case *sqlparser.RangeCond:
		return evaluateRangeCond(attrs, e, params)
	default:
		return false
	}
}

// evaluateComparison applies DynamoDB type semantics: equality requires the
// same type with the same payload (a string "1" never equals a number 1,
// a binary never equals the empty string), and ordering is defined only for
// same-type S/N/B operands. A cross-type <> holds because the operands are
// not equal.
func evaluateComparison(attrs map[string]*dbstore.AttributeValue, cmp *sqlparser.ComparisonExpr, params *partiQLParams) bool {
	attrName := extractColName(cmp.Left)
	if attrName == "" {
		return false
	}

	attr, exists := attrs[attrName]
	if !exists {
		return false
	}

	right := whereOperand(cmp.Right, attrs, params)
	if right == nil {
		return false
	}

	switch cmp.Operator {
	case sqlparser.EqualStr:
		return attributeValuesEqual(attr, right)
	case sqlparser.NotEqualStr:
		return !attributeValuesEqual(attr, right)
	case sqlparser.LessThanStr:
		c, ok := compareOrderedValues(attr, right)
		return ok && c < 0
	case sqlparser.GreaterThanStr:
		c, ok := compareOrderedValues(attr, right)
		return ok && c > 0
	case sqlparser.LessEqualStr:
		c, ok := compareOrderedValues(attr, right)
		return ok && c <= 0
	case sqlparser.GreaterEqualStr:
		c, ok := compareOrderedValues(attr, right)
		return ok && c >= 0
	case sqlparser.InStr:
		return evaluateIn(attrs, cmp, params)
	case sqlparser.LikeStr:
		return evaluateLike(attr, right)
	}
	return false
}

// whereOperand resolves the right-hand side of a comparison to an
// AttributeValue: a literal or bound parameter through the shared value
// materialiser, a column reference to the referenced attribute's value.
// An unresolvable operand never matches.
func whereOperand(expr sqlparser.Expr, attrs map[string]*dbstore.AttributeValue, params *partiQLParams) *dbstore.AttributeValue {
	if col, ok := expr.(*sqlparser.ColName); ok {
		return attrs[col.Name.String()]
	}
	v, err := exprToAttributeValueWithParams(expr, params)
	if err != nil {
		return nil
	}
	return v
}

func extractColName(expr sqlparser.Expr) string {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		return e.Name.String()
	case *sqlparser.SQLVal:
		return string(e.Val)
	}
	return ""
}

func evaluateIsExpr(attrs map[string]*dbstore.AttributeValue, is *sqlparser.IsExpr, params *partiQLParams) bool {
	attrName := extractColName(is.Expr)
	if attrName == "" {
		return false
	}

	attr, exists := attrs[attrName]
	switch is.Operator {
	case sqlparser.IsNullStr:
		return !exists || attr == nil || (attr.NULL != nil && *attr.NULL)
	case sqlparser.IsNotNullStr:
		return exists && attr != nil && (attr.NULL == nil || !*attr.NULL)
	}
	return false
}

func evaluateRangeCond(attrs map[string]*dbstore.AttributeValue, rc *sqlparser.RangeCond, params *partiQLParams) bool {
	attrName := extractColName(rc.Left)
	if attrName == "" {
		return false
	}

	attr, exists := attrs[attrName]
	if !exists {
		return false
	}

	from := whereOperand(rc.From, attrs, params)
	to := whereOperand(rc.To, attrs, params)
	if from == nil || to == nil {
		return false
	}

	if rc.Operator == sqlparser.BetweenStr {
		cFrom, okFrom := compareOrderedValues(attr, from)
		cTo, okTo := compareOrderedValues(attr, to)
		return okFrom && okTo && cFrom >= 0 && cTo <= 0
	}
	// NOT BETWEEN: an operand pair that cannot be ordered does not match.
	cFrom, okFrom := compareOrderedValues(attr, from)
	if okFrom && cFrom < 0 {
		return true
	}
	cTo, okTo := compareOrderedValues(attr, to)
	return okTo && cTo > 0
}

func evaluateIn(attrs map[string]*dbstore.AttributeValue, cmp *sqlparser.ComparisonExpr, params *partiQLParams) bool {
	attrName := extractColName(cmp.Left)
	if attrName == "" {
		return false
	}

	attr, exists := attrs[attrName]
	if !exists {
		return false
	}

	tuple, ok := cmp.Right.(sqlparser.ValTuple)
	if !ok {
		return false
	}

	for _, item := range tuple {
		if v := whereOperand(item, attrs, params); v != nil && attributeValuesEqual(attr, v) {
			return true
		}
	}
	return false
}

// evaluateLike matches SQL LIKE patterns; only string operands participate.
func evaluateLike(value, pattern *dbstore.AttributeValue) bool {
	if value.S == nil || pattern.S == nil {
		return false
	}
	regex := likeToRegex(*pattern.S)
	matched, _ := regexp.MatchString("^"+regex+"$", *value.S)
	return matched
}

func likeToRegex(pattern string) string {
	var result strings.Builder
	for _, ch := range pattern {
		switch ch {
		case '%':
			result.WriteString(".*")
		case '_':
			result.WriteString(".")
		default:
			result.WriteString(regexp.QuoteMeta(string(ch)))
		}
	}
	return result.String()
}
