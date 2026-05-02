package timestreamquery

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	tsstore "vorpalstacks/internal/store/aws/timestream"
	"vorpalstacks/pkg/sqlparser"
)

func (s *TimestreamQueryService) executeSQLQuery(ctx context.Context, reqCtx *request.RequestContext, queryString string) (*QueryResult, error) {
	processedSQL := queryString

	if strings.Contains(queryString, "::") {
		processedSQL = s.convertCastOperator(queryString)
	}

	opts := sqlparser.ParserOptions{
		Dialect: sqlparser.DialectTimestream,
	}
	stmt, err := sqlparser.ParseWithOptions(processedSQL, opts)
	if err != nil {
		return nil, fmt.Errorf("SQL parse error: %w", err)
	}

	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, fmt.Errorf("only SELECT statements are supported")
	}

	databaseName, tableName, err := s.extractTableInfo(selectStmt)
	if err != nil {
		return nil, err
	}

	if databaseName == "" || tableName == "" {
		return s.executeExpressionQuery(selectStmt)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}
	records, err := st.recordStore.QueryRecords(databaseName, tableName, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, err
	}

	columnInfo := s.buildColumnInfo(selectStmt, records)

	rows := s.applyQuery(selectStmt, records)

	return &QueryResult{
		QueryID:    "",
		Rows:       rows,
		ColumnInfo: columnInfo,
	}, nil
}

func (s *TimestreamQueryService) executeExpressionQuery(selectStmt *sqlparser.Select) (*QueryResult, error) {
	var columns []ColumnInfo
	var row map[string]interface{}

	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			colName := s.extractColumnName(aliased.Expr)
			if !aliased.As.IsEmpty() {
				colName = aliased.As.String()
			}

			var value interface{}
			scalarType := "INTEGER"

			switch e := aliased.Expr.(type) {
			case *sqlparser.SQLVal:
				switch e.Type {
				case sqlparser.IntVal:
					value = string(e.Val)
					scalarType = "INTEGER"
				case sqlparser.FloatVal:
					value = string(e.Val)
					scalarType = "DOUBLE"
				case sqlparser.StrVal:
					value = string(e.Val)
					scalarType = "VARCHAR"
				default:
					value = string(e.Val)
				}
			case *sqlparser.FuncExpr:
				funcName := strings.ToLower(e.Name.String())
				switch funcName {
				case "now", "current_timestamp", "current_date", "current_time":
					value = time.Now().UTC().Format("2006-01-02 15:04:05.000000")
					scalarType = "TIMESTAMP"
				default:
					value = sqlparser.String(aliased.Expr)
					scalarType = "VARCHAR"
				}
			default:
				value = sqlparser.String(aliased.Expr)
			}

			if row == nil {
				row = make(map[string]interface{})
			}
			row[colName] = map[string]interface{}{
				"ScalarValue": value,
			}
			columns = append(columns, ColumnInfo{
				Name: colName,
				Type: ColumnTypeInfo{ScalarType: scalarType},
			})
		}
	}

	var rows []map[string]interface{}
	if row != nil {
		rows = append(rows, row)
	}

	return &QueryResult{
		QueryID:    "",
		Rows:       rows,
		ColumnInfo: columns,
	}, nil
}

func (s *TimestreamQueryService) convertCastOperator(sql string) string {
	var result strings.Builder
	inString := false
	i := 0

	for i < len(sql) {
		ch := sql[i]

		if ch == '\'' {
			inString = !inString
			result.WriteByte(ch)
			i++
			continue
		}

		if inString {
			result.WriteByte(ch)
			i++
			continue
		}

		if i+1 < len(sql) && sql[i:i+2] == "::" {
			j := i + 2
			for j < len(sql) && (isAlphaNum(sql[j]) || sql[j] == '_') {
				j++
			}
			i = j
			continue
		}

		result.WriteByte(ch)
		i++
	}

	return result.String()
}

func (s *TimestreamQueryService) extractTableInfo(selectStmt *sqlparser.Select) (databaseName, tableName string, err error) {
	if len(selectStmt.From) == 0 {
		return "", "", fmt.Errorf("no table specified in FROM clause")
	}

	aliasedTableExpr, ok := selectStmt.From[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return "", "", fmt.Errorf("unsupported table expression")
	}

	tableNameExpr, ok := aliasedTableExpr.Expr.(sqlparser.TableName)
	if !ok {
		return "", "", fmt.Errorf("unsupported table name format")
	}

	if !tableNameExpr.Qualifier.IsEmpty() {
		databaseName = tableNameExpr.Qualifier.String()
	}
	tableName = tableNameExpr.Name.String()

	return databaseName, tableName, nil
}

func (s *TimestreamQueryService) buildColumnInfo(selectStmt *sqlparser.Select, records []*tsstore.StoredRecord) []ColumnInfo {
	var columns []ColumnInfo

	if len(selectStmt.SelectExprs) == 1 {
		if _, isStar := selectStmt.SelectExprs[0].(*sqlparser.StarExpr); isStar {
			if len(records) > 0 {
				record := records[0]
				columns = append(columns, ColumnInfo{
					Name: "time",
					Type: ColumnTypeInfo{ScalarType: "TIMESTAMP"},
				})
				for _, dim := range record.Dimensions {
					columns = append(columns, ColumnInfo{
						Name: dim.Name,
						Type: ColumnTypeInfo{ScalarType: "VARCHAR"},
					})
				}
				columns = append(columns, ColumnInfo{
					Name: "measure_value::double",
					Type: ColumnTypeInfo{ScalarType: "DOUBLE"},
				})
				columns = append(columns, ColumnInfo{
					Name: "measure_name",
					Type: ColumnTypeInfo{ScalarType: "VARCHAR"},
				})
			}
			return columns
		}
	}

	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			colName := s.extractColumnName(aliased.Expr)
			if !aliased.As.IsEmpty() {
				colName = aliased.As.String()
			}
			scalarType := "VARCHAR"
			if strings.Contains(strings.ToLower(colName), "time") {
				scalarType = "TIMESTAMP"
			} else if strings.Contains(colName, "measure_value") {
				scalarType = "DOUBLE"
			}
			if _, isFunc := aliased.Expr.(*sqlparser.FuncExpr); isFunc {
				if fn, ok := aliased.Expr.(*sqlparser.FuncExpr); ok && fn.IsAggregate() {
					scalarType = "DOUBLE"
				}
			}
			columns = append(columns, ColumnInfo{
				Name: colName,
				Type: ColumnTypeInfo{ScalarType: scalarType},
			})
		}
	}

	return columns
}

func (s *TimestreamQueryService) extractColumnName(expr sqlparser.Expr) string {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		return e.Name.String()
	case *sqlparser.SQLVal:
		return string(e.Val)
	case *sqlparser.FuncExpr:
		funcName := strings.ToLower(e.Name.String())
		if funcName == "cast" && len(e.Exprs) > 0 {
			if aliased, ok := e.Exprs[0].(*sqlparser.AliasedExpr); ok {
				return s.extractColumnName(aliased.Expr)
			}
		}
		return sqlparser.String(expr)
	default:
		return sqlparser.String(expr)
	}
}

func (s *TimestreamQueryService) applyQuery(selectStmt *sqlparser.Select, records []*tsstore.StoredRecord) []map[string]interface{} {
	var rows []map[string]interface{}

	for _, record := range records {
		row := s.recordToRow(record)
		if selectStmt.Where != nil {
			if !s.evaluateWhere(selectStmt.Where.Expr, row) {
				continue
			}
		}
		rows = append(rows, row)
	}

	if s.hasAggregate(selectStmt) {
		rows = s.applyGroupBy(rows, selectStmt)

		if len(selectStmt.OrderBy) > 0 {
			s.applyOrderBy(rows, selectStmt.OrderBy)
		}
		if selectStmt.Limit != nil {
			rows = s.applyLimit(rows, selectStmt.Limit)
		}
		return s.projectAggregateColumns(rows, selectStmt)
	}

	if len(selectStmt.OrderBy) > 0 {
		s.applyOrderBy(rows, selectStmt.OrderBy)
	}

	if selectStmt.Limit != nil {
		rows = s.applyLimit(rows, selectStmt.Limit)
	}

	return s.projectColumns(rows, selectStmt.SelectExprs)
}

func (s *TimestreamQueryService) recordToRow(record *tsstore.StoredRecord) map[string]interface{} {
	row := map[string]interface{}{
		"time":         record.Timestamp.Format(time.RFC3339),
		"measure_name": record.MeasureName,
	}

	for _, dim := range record.Dimensions {
		row[dim.Name] = dim.Value
	}

	var measureValue interface{}
	switch record.MeasureValueType {
	case tsstore.MeasureValueTypeDouble:
		if val, err := strconv.ParseFloat(record.MeasureValue, 64); err == nil {
			measureValue = val
			row["measure_value::double"] = val
		}
	case tsstore.MeasureValueTypeBigint:
		if val, err := strconv.ParseInt(record.MeasureValue, 10, 64); err == nil {
			measureValue = val
			row["measure_value::bigint"] = val
		}
	case tsstore.MeasureValueTypeBoolean:
		measureValue = strings.ToLower(record.MeasureValue) == "true"
		row["measure_value::boolean"] = measureValue
	case tsstore.MeasureValueTypeTimestamp:
		measureValue = record.MeasureValue
		row["measure_value::timestamp"] = record.MeasureValue
	default:
		measureValue = record.MeasureValue
		row["measure_value::varchar"] = record.MeasureValue
	}

	if measureValue != nil {
		row["measure_value"] = measureValue
	}

	return row
}
