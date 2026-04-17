package athena

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"vorpalstacks/pkg/sqlparser"
)

func (s *AthenaService) evaluateWhere(expr sqlparser.Expr, row map[string]interface{}) bool {
	switch e := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return s.evaluateComparison(e, row)
	case *sqlparser.AndExpr:
		return s.evaluateWhere(e.Left, row) && s.evaluateWhere(e.Right, row)
	case *sqlparser.OrExpr:
		return s.evaluateWhere(e.Left, row) || s.evaluateWhere(e.Right, row)
	case *sqlparser.ParenExpr:
		return s.evaluateWhere(e.Expr, row)
	case *sqlparser.IsExpr:
		return s.evaluateIs(e, row)
	}
	return true
}

func (s *AthenaService) evaluateComparison(expr *sqlparser.ComparisonExpr, row map[string]interface{}) bool {
	leftVal := s.getExprValue(expr.Left, row)
	rightVal := s.getExprValue(expr.Right, row)

	switch expr.Operator {
	case sqlparser.EqualStr:
		return s.compareValues(leftVal, rightVal) == 0
	case sqlparser.NotEqualStr:
		return s.compareValues(leftVal, rightVal) != 0
	case sqlparser.LessThanStr:
		return s.compareValues(leftVal, rightVal) < 0
	case sqlparser.LessEqualStr:
		return s.compareValues(leftVal, rightVal) <= 0
	case sqlparser.GreaterThanStr:
		return s.compareValues(leftVal, rightVal) > 0
	case sqlparser.GreaterEqualStr:
		return s.compareValues(leftVal, rightVal) >= 0
	case sqlparser.LikeStr:
		return s.matchLike(fmt.Sprintf("%v", leftVal), fmt.Sprintf("%v", rightVal))
	case sqlparser.NotLikeStr:
		return !s.matchLike(fmt.Sprintf("%v", leftVal), fmt.Sprintf("%v", rightVal))
	}
	return false
}

func (s *AthenaService) compareValues(left, right interface{}) int {
	leftFloat, leftErr := s.toFloat(left)
	rightFloat, rightErr := s.toFloat(right)

	if leftErr == nil && rightErr == nil {
		if leftFloat < rightFloat {
			return -1
		} else if leftFloat > rightFloat {
			return 1
		}
		return 0
	}

	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)
	if leftStr < rightStr {
		return -1
	} else if leftStr > rightStr {
		return 1
	}
	return 0
}

func (s *AthenaService) evaluateIs(expr *sqlparser.IsExpr, row map[string]interface{}) bool {
	val := s.getExprValue(expr.Expr, row)
	isNull := val == nil
	switch expr.Operator {
	case sqlparser.IsNullStr:
		return isNull
	case sqlparser.IsNotNullStr:
		return !isNull
	}
	return false
}

func (s *AthenaService) getExprValue(expr sqlparser.Expr, row map[string]interface{}) interface{} {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		colName := e.Name.String()
		return row[colName]
	case *sqlparser.SQLVal:
		if e.Type == sqlparser.StrVal {
			return string(e.Val)
		} else if e.Type == sqlparser.IntVal {
			if val, err := strconv.ParseInt(string(e.Val), 10, 64); err == nil {
				return val
			}
		} else if e.Type == sqlparser.FloatVal {
			if val, err := strconv.ParseFloat(string(e.Val), 64); err == nil {
				return val
			}
		}
		return string(e.Val)
	case *sqlparser.NullVal:
		return nil
	}
	return nil
}

func (s *AthenaService) matchLike(value, pattern string) bool {
	pattern = strings.Trim(pattern, "'")
	pattern = strings.ToLower(pattern)
	value = strings.ToLower(value)

	regexPattern := s.likePatternToRegex(pattern)
	matched, _ := regexp.MatchString("^"+regexPattern+"$", value)
	return matched
}

func (s *AthenaService) likePatternToRegex(pattern string) string {
	var result strings.Builder
	escaped := false

	for _, ch := range pattern {
		if escaped {
			switch ch {
			case '%', '_', '\\':
				result.WriteRune(ch)
			default:
				result.WriteRune('\\')
				result.WriteRune(ch)
			}
			escaped = false
			continue
		}

		switch ch {
		case '\\':
			escaped = true
		case '%':
			result.WriteString(".*")
		case '_':
			result.WriteString(".")
		case '.', '^', '$', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|':
			result.WriteRune('\\')
			result.WriteRune(ch)
		default:
			result.WriteRune(ch)
		}
	}

	if escaped {
		result.WriteRune('\\')
	}

	return result.String()
}
