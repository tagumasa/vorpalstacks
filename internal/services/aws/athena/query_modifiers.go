package athena

import (
	"sort"
	"strconv"
	"strings"

	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

func (s *AthenaService) applyOrderBy(rows []*athenastore.StoredRow, orderBy sqlparser.OrderBy) {
	sort.Slice(rows, func(i, j int) bool {
		for _, order := range orderBy {
			colName := s.extractColumnName(order.Expr)
			valI := rows[i].Values[colName]
			valJ := rows[j].Values[colName]

			cmp := s.compareValues(valI, valJ)
			if cmp == 0 {
				continue
			}

			if order.Direction == sqlparser.DescScr {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})
}

func (s *AthenaService) applyLimit(rows []*athenastore.StoredRow, limit *sqlparser.Limit) []*athenastore.StoredRow {
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

func (s *AthenaService) projectColumns(rows []*athenastore.StoredRow, selectExprs sqlparser.SelectExprs) []*athenastore.StoredRow {
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

				if colName, ok := aliased.Expr.(*sqlparser.ColName); ok && !colName.Qualifier.IsEmpty() {
					qualifiedKey := colName.Qualifier.Name.String() + "." + sourceColName
					if val, exists := row.Values[qualifiedKey]; exists {
						outputName := sourceColName
						if !aliased.As.IsEmpty() {
							outputName = aliased.As.String()
						}
						projectedRow.Values[outputName] = val
						continue
					}
				}

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

func (s *AthenaService) extractColumnName(expr sqlparser.Expr) string {
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

func (s *AthenaService) buildColumnInfoFromSelect(selectStmt *sqlparser.Select) []athenastore.ColumnInfo {
	return s.buildColumnInfoFromSelectWithTypes(selectStmt, nil)
}

func (s *AthenaService) buildColumnInfoFromSelectWithTypes(selectStmt *sqlparser.Select, tableCols []athenastore.Column) []athenastore.ColumnInfo {
	typeMap := make(map[string]string)
	for _, c := range tableCols {
		typeMap[strings.ToLower(c.Name)] = c.Type
	}

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
				Type:  s.inferExprType(aliased.Expr, typeMap),
			})
		}
	}

	return columns
}

// inferExprType determines the Athena data type for a SQL expression.
func (s *AthenaService) inferExprType(expr sqlparser.Expr, typeMap map[string]string) string {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		if t, ok := typeMap[strings.ToLower(e.Name.String())]; ok {
			return t
		}
		return "string"

	case *sqlparser.SQLVal:
		switch e.Type {
		case sqlparser.IntVal:
			return "integer"
		case sqlparser.FloatVal:
			return "double"
		default:
			return "string"
		}

	case *sqlparser.FuncExpr:
		upperName := strings.ToUpper(e.Name.String())
		switch upperName {
		case "COUNT":
			return "bigint"
		case "SUM", "AVG":
			return "double"
		case "MIN", "MAX":
			if len(e.Exprs) > 0 {
				if aliased, ok := e.Exprs[0].(*sqlparser.AliasedExpr); ok {
					return s.inferExprType(aliased.Expr, typeMap)
				}
			}
			return "string"
		default:
			return "string"
		}

	default:
		return "string"
	}
}
