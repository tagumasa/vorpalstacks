package athena

import (
	"fmt"
	"sort"
	"strconv"

	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

func (s *Service) applyOrderBy(rows []*athenastore.StoredRow, orderBy sqlparser.OrderBy) {
	sort.Slice(rows, func(i, j int) bool {
		for _, order := range orderBy {
			colName := s.extractColumnName(order.Expr)
			valI := rows[i].Values[colName]
			valJ := rows[j].Values[colName]

			strI := fmt.Sprintf("%v", valI)
			strJ := fmt.Sprintf("%v", valJ)

			if strI == strJ {
				continue
			}

			less := strI < strJ
			if order.Direction == sqlparser.DescScr {
				return !less
			}
			return less
		}
		return false
	})
}

func (s *Service) applyLimit(rows []*athenastore.StoredRow, limit *sqlparser.Limit) []*athenastore.StoredRow {
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
		return []*athenastore.StoredRow{}
	}
	return rows[offset:end]
}

func (s *Service) projectColumns(rows []*athenastore.StoredRow, selectExprs sqlparser.SelectExprs) []*athenastore.StoredRow {
	if len(selectExprs) == 1 {
		if _, isStar := selectExprs[0].(*sqlparser.StarExpr); isStar {
			return rows
		}
	}

	var result []*athenastore.StoredRow
	for _, row := range rows {
		projectedRow := &athenastore.StoredRow{Values: make(map[string]interface{})}
		for _, expr := range selectExprs {
			if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
				sourceColName := s.extractColumnName(aliased.Expr)
				if val, exists := row.Values[sourceColName]; exists {
					outputName := sourceColName
					if !aliased.As.IsEmpty() {
						outputName = aliased.As.String()
					}
					projectedRow.Values[outputName] = val
				}
			}
		}
		result = append(result, projectedRow)
	}
	return result
}

func (s *Service) extractColumnName(expr sqlparser.Expr) string {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		return e.Name.String()
	case *sqlparser.SQLVal:
		return string(e.Val)
	case *sqlparser.FuncExpr:
		return sqlparser.String(expr)
	default:
		return sqlparser.String(expr)
	}
}

func (s *Service) buildColumnInfoFromSelect(selectStmt *sqlparser.Select) []athenastore.ColumnInfo {
	var columns []athenastore.ColumnInfo

	if len(selectStmt.SelectExprs) == 1 {
		if _, isStar := selectStmt.SelectExprs[0].(*sqlparser.StarExpr); isStar {
			return columns
		}
	}

	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			colName := s.extractColumnName(aliased.Expr)
			outputName := colName
			if !aliased.As.IsEmpty() {
				outputName = aliased.As.String()
			}
			columns = append(columns, athenastore.ColumnInfo{
				Label: outputName,
				Name:  outputName,
				Type:  "string",
			})
		}
	}

	return columns
}
