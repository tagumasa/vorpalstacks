package timestreamquery

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/pkg/sqlparser"
)

func (s *Service) evaluateWhere(expr sqlparser.Expr, row map[string]interface{}) bool {
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

func (s *Service) evaluateComparison(expr *sqlparser.ComparisonExpr, row map[string]interface{}) bool {
	leftVal := s.getExprValue(expr.Left, row)
	rightVal := s.getExprValue(expr.Right, row)

	if leftVal == nil || rightVal == nil {
		if expr.Operator == sqlparser.EqualStr {
			return leftVal == nil && rightVal == nil
		}
		return false
	}

	leftStr := fmt.Sprintf("%v", leftVal)
	rightStr := fmt.Sprintf("%v", rightVal)

	switch expr.Operator {
	case sqlparser.EqualStr:
		return s.compareEqual(leftVal, rightVal)
	case sqlparser.NotEqualStr:
		return !s.compareEqual(leftVal, rightVal)
	case sqlparser.LessThanStr:
		return s.compareLess(leftVal, rightVal, leftStr, rightStr)
	case sqlparser.LessEqualStr:
		return s.compareLess(leftVal, rightVal, leftStr, rightStr) || s.compareEqual(leftVal, rightVal)
	case sqlparser.GreaterThanStr:
		return s.compareGreater(leftVal, rightVal, leftStr, rightStr)
	case sqlparser.GreaterEqualStr:
		return s.compareGreater(leftVal, rightVal, leftStr, rightStr) || s.compareEqual(leftVal, rightVal)
	case sqlparser.LikeStr:
		return s.matchLike(leftStr, rightStr)
	case sqlparser.NotLikeStr:
		return !s.matchLike(leftStr, rightStr)
	}
	return false
}

func (s *Service) parseTimestampValue(v string) (time.Time, error) {
	if ts, err := time.Parse(time.RFC3339, v); err == nil {
		return ts, nil
	}
	if ts, err := time.Parse(time.RFC3339Nano, v); err == nil {
		return ts, nil
	}
	if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
		if ts > 1e12 {
			return time.UnixMilli(ts), nil
		}
		return time.Unix(ts, 0), nil
	}
	return time.Time{}, fmt.Errorf("not a timestamp")
}

func (s *Service) evaluateIs(expr *sqlparser.IsExpr, row map[string]interface{}) bool {
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

func (s *Service) getExprValue(expr sqlparser.Expr, row map[string]interface{}) interface{} {
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

func (s *Service) matchLike(value, pattern string) bool {
	regexPattern := regexp.QuoteMeta(pattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\%`, `.*`)
	regexPattern = strings.ReplaceAll(regexPattern, `\_`, `.`)
	regexPattern = `(?i)^` + regexPattern + `$`
	matched, _ := regexp.MatchString(regexPattern, value)
	return matched
}
