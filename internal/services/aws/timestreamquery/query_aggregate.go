package timestreamquery

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"vorpalstacks/pkg/sqlparser"
)

func (s *TimestreamQueryService) applyGroupBy(rows []map[string]interface{}, selectStmt *sqlparser.Select) []map[string]interface{} {
	if !s.hasAggregate(selectStmt) {
		return rows
	}

	groupKeys := s.extractGroupKeys(selectStmt)
	groups := s.groupRows(rows, groupKeys)

	var result []map[string]interface{}
	for _, group := range groups {
		aggRow := s.computeAggregates(group.Rows, group.KeyValues, selectStmt)
		result = append(result, aggRow)
	}
	return result
}

func (s *TimestreamQueryService) hasAggregate(selectStmt *sqlparser.Select) bool {
	if len(selectStmt.GroupBy) > 0 {
		return true
	}
	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			if fn, ok := aliased.Expr.(*sqlparser.FuncExpr); ok && fn.IsAggregate() {
				return true
			}
		}
	}
	return false
}

func (s *TimestreamQueryService) extractGroupKeys(selectStmt *sqlparser.Select) []string {
	var keys []string
	for _, expr := range selectStmt.GroupBy {
		keys = append(keys, s.extractColumnName(expr))
	}
	return keys
}

type rowGroup struct {
	KeyValues map[string]interface{}
	Rows      []map[string]interface{}
}

func (s *TimestreamQueryService) groupRows(rows []map[string]interface{}, groupKeys []string) []rowGroup {
	if len(groupKeys) == 0 {
		return []rowGroup{{Rows: rows}}
	}

	groupMap := make(map[string]*rowGroup)
	var groupOrder []string

	for _, row := range rows {
		key := s.buildGroupKey(row, groupKeys)
		if g, ok := groupMap[key]; ok {
			g.Rows = append(g.Rows, row)
		} else {
			keyValues := make(map[string]interface{})
			for _, k := range groupKeys {
				keyValues[k] = row[k]
			}
			g := &rowGroup{KeyValues: keyValues, Rows: []map[string]interface{}{row}}
			groupMap[key] = g
			groupOrder = append(groupOrder, key)
		}
	}

	var groups []rowGroup
	for _, key := range groupOrder {
		groups = append(groups, *groupMap[key])
	}
	return groups
}

func (s *TimestreamQueryService) buildGroupKey(row map[string]interface{}, groupKeys []string) string {
	var parts []string
	for _, k := range groupKeys {
		parts = append(parts, fmt.Sprintf("%v", row[k]))
	}
	return strings.Join(parts, "\x00")
}

func (s *TimestreamQueryService) computeAggregates(rows []map[string]interface{}, keyValues map[string]interface{}, selectStmt *sqlparser.Select) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range keyValues {
		result[k] = v
	}

	for _, expr := range selectStmt.SelectExprs {
		aliased, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}
		fn, ok := aliased.Expr.(*sqlparser.FuncExpr)
		if !ok || !fn.IsAggregate() {
			continue
		}

		outputName := s.getOutputName(aliased.Expr)
		if !aliased.As.IsEmpty() {
			outputName = aliased.As.String()
		}

		colName := ""
		if len(fn.Exprs) > 0 {
			if arg, ok := fn.Exprs[0].(*sqlparser.AliasedExpr); ok {
				colName = s.extractColumnName(arg.Expr)
			}
		}

		funcName := strings.ToLower(fn.Name.String())
		result[outputName] = s.evaluateAggregateFunc(funcName, colName, rows)
	}
	return result
}

func (s *TimestreamQueryService) evaluateAggregateFunc(funcName, colName string, rows []map[string]interface{}) interface{} {
	switch funcName {
	case "count":
		if colName == "" {
			return len(rows)
		}
		count := 0
		for _, row := range rows {
			if v, ok := row[colName]; ok && v != nil {
				count++
			}
		}
		return count
	case "sum":
		return s.aggregateNumeric(colName, rows, func(sum, val float64) float64 { return sum + val })
	case "avg":
		sum := s.aggregateNumeric(colName, rows, func(sum, val float64) float64 { return sum + val })
		if sum == nil {
			return nil
		}
		count := 0
		for _, row := range rows {
			if v, ok := row[colName]; ok && v != nil {
				if _, err := toFloat64(v); err == nil {
					count++
				}
			}
		}
		if count == 0 {
			return nil
		}
		return sum.(float64) / float64(count)
	case "min":
		return s.aggregateNumeric(colName, rows, func(current, val float64) float64 {
			return math.Min(current, val)
		})
	case "max":
		return s.aggregateNumeric(colName, rows, func(current, val float64) float64 {
			return math.Max(current, val)
		})
	}
	return nil
}

type aggregateFn func(current, val float64) float64

func (s *TimestreamQueryService) aggregateNumeric(colName string, rows []map[string]interface{}, fn aggregateFn) interface{} {
	var result float64
	started := false
	for _, row := range rows {
		v, ok := row[colName]
		if !ok || v == nil {
			continue
		}
		fval, err := toFloat64(v)
		if err != nil {
			continue
		}
		if !started {
			result = fval
			started = true
		} else {
			result = fn(result, fval)
		}
	}
	if !started {
		return nil
	}
	return result
}

func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int64:
		return float64(val), nil
	case int:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	}
	return 0, fmt.Errorf("cannot convert %T to float64", v)
}

func (s *TimestreamQueryService) projectAggregateColumns(rows []map[string]interface{}, selectStmt *sqlparser.Select) []map[string]interface{} {
	colNames := s.extractSelectColumnNames(selectStmt)
	var result []map[string]interface{}
	for _, row := range rows {
		projected := make(map[string]interface{})
		for _, name := range colNames {
			if v, ok := row[name]; ok {
				projected[name] = v
			}
		}
		result = append(result, projected)
	}
	return result
}

func (s *TimestreamQueryService) extractSelectColumnNames(selectStmt *sqlparser.Select) []string {
	var names []string
	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			if !aliased.As.IsEmpty() {
				names = append(names, aliased.As.String())
			} else {
				names = append(names, s.getOutputName(aliased.Expr))
			}
		}
	}
	return names
}
