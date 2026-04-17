package timestreamquery

import (
	"fmt"
	"sort"
	"strconv"

	"vorpalstacks/pkg/sqlparser"
)

func (s *TimestreamQueryService) applyOrderBy(rows []map[string]interface{}, orderBy sqlparser.OrderBy) {
	sort.Slice(rows, func(i, j int) bool {
		for _, order := range orderBy {
			colName := s.extractColumnName(order.Expr)
			valI := rows[i][colName]
			valJ := rows[j][colName]

			if s.compareEqual(valI, valJ) {
				continue
			}

			strI := fmt.Sprintf("%v", valI)
			strJ := fmt.Sprintf("%v", valJ)

			less := s.compareLess(valI, valJ, strI, strJ)
			if order.Direction == sqlparser.DescScr {
				return !less
			}
			return less
		}
		return false
	})
}

func (s *TimestreamQueryService) applyLimit(rows []map[string]interface{}, limit *sqlparser.Limit) []map[string]interface{} {
	count := 0
	if rowcount := limit.Rowcount; rowcount != nil {
		if sqlVal, ok := rowcount.(*sqlparser.SQLVal); ok {
			count, _ = strconv.Atoi(string(sqlVal.Val))
		}
	}

	offset := 0
	if limit.Offset != nil {
		if sqlVal, ok := limit.Offset.(*sqlparser.SQLVal); ok {
			offset, _ = strconv.Atoi(string(sqlVal.Val))
		}
	}

	if count <= 0 {
		return rows
	}

	end := offset + count
	if end > len(rows) {
		end = len(rows)
	}
	if offset >= len(rows) {
		return []map[string]interface{}{}
	}
	return rows[offset:end]
}

func (s *TimestreamQueryService) projectColumns(rows []map[string]interface{}, selectExprs sqlparser.SelectExprs) []map[string]interface{} {
	if len(selectExprs) == 1 {
		if _, isStar := selectExprs[0].(*sqlparser.StarExpr); isStar {
			return rows
		}
	}

	var result []map[string]interface{}
	for _, row := range rows {
		projectedRow := make(map[string]interface{})
		for _, expr := range selectExprs {
			if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
				sourceColName := s.extractColumnName(aliased.Expr)
				if val, exists := row[sourceColName]; exists {
					outputName := sourceColName
					if !aliased.As.IsEmpty() {
						outputName = aliased.As.String()
					} else if _, isFunc := aliased.Expr.(*sqlparser.FuncExpr); isFunc {
						outputName = s.getOutputName(aliased.Expr)
					}
					projectedRow[outputName] = val
				}
			}
		}
		result = append(result, projectedRow)
	}
	return result
}

func (s *TimestreamQueryService) getOutputName(expr sqlparser.Expr) string {
	return sqlparser.String(expr)
}
