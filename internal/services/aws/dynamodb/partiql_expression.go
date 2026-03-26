package dynamodb

import (
	"math/big"
	"regexp"
	"strconv"
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

func evaluateComparison(attrs map[string]*dbstore.AttributeValue, cmp *sqlparser.ComparisonExpr, params *partiQLParams) bool {
	attrName := extractColName(cmp.Left)
	if attrName == "" {
		return false
	}

	attr, exists := attrs[attrName]
	if !exists {
		return false
	}

	leftVal := getAttrValue(attr)
	rightVal := extractValue(cmp.Right, params)

	switch cmp.Operator {
	case sqlparser.EqualStr:
		return valuesEqual(leftVal, rightVal)
	case sqlparser.NotEqualStr:
		return !valuesEqual(leftVal, rightVal)
	case sqlparser.LessThanStr:
		return compareValues(leftVal, rightVal) < 0
	case sqlparser.GreaterThanStr:
		return compareValues(leftVal, rightVal) > 0
	case sqlparser.LessEqualStr:
		return compareValues(leftVal, rightVal) <= 0
	case sqlparser.GreaterEqualStr:
		return compareValues(leftVal, rightVal) >= 0
	case sqlparser.InStr:
		return evaluateIn(attrs, cmp, params)
	case sqlparser.LikeStr:
		return evaluateLike(leftVal, rightVal)
	}
	return false
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

func extractValue(expr sqlparser.Expr, params *partiQLParams) string {
	switch e := expr.(type) {
	case *sqlparser.SQLVal:
		switch e.Type {
		case sqlparser.StrVal:
			return string(e.Val)
		case sqlparser.IntVal, sqlparser.FloatVal:
			return string(e.Val)
		case sqlparser.ValArg:
			if strings.HasPrefix(string(e.Val), ":") {
				idxStr := strings.TrimPrefix(string(e.Val), ":v")
				if idx, err := strconv.Atoi(idxStr); err == nil && params != nil && idx > 0 && idx <= len(params.Parameters) {
					return paramToString(params.Parameters[idx-1])
				}
			}
			return string(e.Val)
		}
	case *sqlparser.ColName:
		return e.Name.String()
	}
	return ""
}

func getAttrValue(attr *dbstore.AttributeValue) string {
	switch {
	case attr.S != nil:
		return *attr.S
	case attr.N != nil:
		return *attr.N
	case attr.BOOL != nil:
		if *attr.BOOL {
			return "true"
		}
		return "false"
	case attr.NULL != nil && *attr.NULL:
		return "null"
	}
	return ""
}

func compareValues(a, b string) int {
	numA, okA := new(big.Rat).SetString(a)
	numB, okB := new(big.Rat).SetString(b)
	if okA && okB {
		return numA.Cmp(numB)
	}
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func valuesEqual(a, b string) bool {
	numA, okA := new(big.Rat).SetString(a)
	numB, okB := new(big.Rat).SetString(b)
	if okA && okB {
		return numA.Cmp(numB) == 0
	}
	return a == b
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

	val := getAttrValue(attr)
	from := extractValue(rc.From, params)
	to := extractValue(rc.To, params)

	if rc.Operator == sqlparser.BetweenStr {
		return compareValues(val, from) >= 0 && compareValues(val, to) <= 0
	}
	return compareValues(val, from) < 0 || compareValues(val, to) > 0
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

	val := getAttrValue(attr)

	tuple, ok := cmp.Right.(sqlparser.ValTuple)
	if !ok {
		return false
	}

	for _, item := range tuple {
		if valuesEqual(extractValue(item, params), val) {
			return true
		}
	}
	return false
}

func evaluateLike(value, pattern string) bool {
	regex := likeToRegex(pattern)
	matched, _ := regexp.MatchString("^"+regex+"$", value)
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

type whereCondition struct {
	attribute string
	operator  string
	value     string
}

func parseWhereConditions(whereClause string) []whereCondition {
	var conditions []whereCondition

	re := regexp.MustCompile(`["']?(\w+)["']?\s*(=|<>|!=|<|>|<=|>=)\s*(\?|["'][^"']*["']|\S+)`)
	matches := re.FindAllStringSubmatch(whereClause, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			conditions = append(conditions, whereCondition{
				attribute: match[1],
				operator:  match[2],
				value:     strings.Trim(match[3], `"'`),
			})
		}
	}

	return conditions
}

func matchesConditions(attrs map[string]*dbstore.AttributeValue, conditions []whereCondition, params *partiQLParams) bool {
	paramIdx := 0
	for _, cond := range conditions {
		attr, exists := attrs[cond.attribute]
		if !exists {
			return false
		}

		var compareValue string
		if cond.value == "?" && params != nil && paramIdx < len(params.Parameters) {
			compareValue = paramToString(params.Parameters[paramIdx])
			paramIdx++
		} else {
			compareValue = strings.Trim(cond.value, `"'`)
		}

		if !compareAttribute(attr, cond.operator, compareValue) {
			return false
		}
	}
	return true
}

func compareAttribute(attr *dbstore.AttributeValue, operator string, value string) bool {
	switch operator {
	case "=", "==":
		return attributeMatchesValue(attr, value)
	case "<>", "!=":
		return !attributeMatchesValue(attr, value)
	case "<":
		return compareAttributeValue(attr, value) < 0
	case ">":
		return compareAttributeValue(attr, value) > 0
	case "<=":
		return compareAttributeValue(attr, value) <= 0
	case ">=":
		return compareAttributeValue(attr, value) >= 0
	}
	return false
}

func attributeMatchesValue(attr *dbstore.AttributeValue, value string) bool {
	switch {
	case attr.S != nil:
		return *attr.S == value
	case attr.N != nil:
		return numbersEqual(*attr.N, value)
	case attr.BOOL != nil:
		if *attr.BOOL {
			return value == "true"
		}
		return value == "false"
	default:
		return false
	}
}

func compareAttributeValue(attr *dbstore.AttributeValue, value string) int {
	if attr.N != nil {
		return compareNumberStrings(*attr.N, value)
	}
	if attr.S != nil {
		if *attr.S < value {
			return -1
		} else if *attr.S > value {
			return 1
		}
		return 0
	}
	return 0
}

func paramToString(param interface{}) string {
	switch v := param.(type) {
	case map[string]interface{}:
		if s, ok := v["S"].(string); ok {
			return s
		}
		if n, ok := v["N"].(string); ok {
			return n
		}
		if b, ok := v["BOOL"].(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	}
	return "null"
}
