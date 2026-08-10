package timestreamquery

import (
	"context"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	tsstore "vorpalstacks/internal/store/aws/timestream"
	"vorpalstacks/pkg/sqlparser"
)

func (s *TimestreamQueryService) executeSQLQuery(ctx context.Context, stores *tsQueryStores, queryString string) (*QueryResult, error) {
	processedSQL := queryString

	if strings.Contains(queryString, "::") {
		processedSQL = s.convertCastOperator(queryString)
	}

	opts := sqlparser.ParserOptions{
		Dialect: sqlparser.DialectTimestream,
	}
	stmt, err := sqlparser.ParseWithOptions(processedSQL, opts)
	if err != nil {
		logs.Debug("Timestream SQL parse error", logs.String("query", processedSQL), logs.Err(err))
		return nil, ErrValidationException
	}

	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		logs.Debug("Timestream query rejected: not a SELECT statement", logs.String("query", processedSQL))
		return nil, ErrValidationException
	}

	databaseName, tableName, err := s.extractTableInfo(selectStmt)
	if err != nil {
		return nil, err
	}

	if databaseName == "" || tableName == "" {
		return s.executeExpressionQuery(selectStmt)
	}

	if stores == nil {
		return nil, ErrInternalServer
	}

	startTime, endTime := extractTimeRange(selectStmt.Where)
	records, err := stores.recordStore.QueryRecords(databaseName, tableName, startTime, endTime)
	if err != nil {
		return nil, err
	}

	columnInfo := s.buildColumnInfo(selectStmt, records)

	rows := s.applyQuery(ctx, selectStmt, records)

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
			row[colName] = value
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
		logs.Debug("Timestream query rejected: no table specified in FROM clause")
		return "", "", ErrQueryExecutionError
	}

	aliasedTableExpr, ok := selectStmt.From[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		logs.Debug("Timestream query rejected: unsupported table expression")
		return "", "", ErrQueryExecutionError
	}

	tableNameExpr, ok := aliasedTableExpr.Expr.(sqlparser.TableName)
	if !ok {
		logs.Debug("Timestream query rejected: unsupported table name format")
		return "", "", ErrQueryExecutionError
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
			scalarType := inferScalarType(colName, aliased.Expr, records)
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

// extractTimeRange inspects the WHERE clause for time predicates
// (e.g. "time > '2024-01-01'") and returns the narrowest time range
// that satisfies all conditions. If no time predicate is found,
// returns a wide default range.
func extractTimeRange(where *sqlparser.Where) (time.Time, time.Time) {
	defaultStart := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultEnd := time.Now().Add(24 * time.Hour)

	if where == nil || where.Expr == nil {
		return defaultStart, defaultEnd
	}

	start := defaultStart
	end := defaultEnd

	walkTimePredicates(where.Expr, &start, &end)

	return start, end
}

func walkTimePredicates(expr sqlparser.Expr, start, end *time.Time) {
	switch e := expr.(type) {
	case *sqlparser.AndExpr:
		walkTimePredicates(e.Left, start, end)
		walkTimePredicates(e.Right, start, end)
	case *sqlparser.ParenExpr:
		walkTimePredicates(e.Expr, start, end)
	case *sqlparser.ComparisonExpr:
		_, isTimeCol := getTimeColumnName(e.Left)
		if !isTimeCol {
			_, isTimeCol = getTimeColumnName(e.Right)
			if !isTimeCol {
				return
			}
		}
		ts, tsOk := getComparisonTimestamp(e)
		if !tsOk {
			return
		}
		switch e.Operator {
		case sqlparser.GreaterThanStr:
			if ts.After(*start) {
				*start = ts.Add(time.Nanosecond)
			}
		case sqlparser.GreaterEqualStr:
			if ts.After(*start) {
				*start = ts
			}
		case sqlparser.LessThanStr:
			if ts.Before(*end) {
				*end = ts.Add(-time.Nanosecond)
			}
		case sqlparser.LessEqualStr:
			if ts.Before(*end) {
				*end = ts
			}
		case sqlparser.EqualStr:
			if ts.After(*start) {
				*start = ts
			}
			if ts.Before(*end) {
				*end = ts
			}
		}
	}
}

func getTimeColumnName(expr sqlparser.Expr) (string, bool) {
	colName, ok := expr.(*sqlparser.ColName)
	if !ok {
		return "", false
	}
	name := colName.Name.String()
	if name == "time" {
		return name, true
	}
	return "", false
}

func getComparisonTimestamp(expr *sqlparser.ComparisonExpr) (time.Time, bool) {
	for _, val := range []sqlparser.Expr{expr.Left, expr.Right} {
		if sqlVal, ok := val.(*sqlparser.SQLVal); ok {
			s := string(sqlVal.Val)
			if sqlVal.Type == sqlparser.StrVal {
				if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
					return t, true
				}
				if t, err := time.Parse(time.RFC3339, s); err == nil {
					return t, true
				}
				if t, err := time.Parse("2006-01-02 15:04:05.000000", s); err == nil {
					return t, true
				}
				if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
					return t, true
				}
				if t, err := time.Parse("2006-01-02", s); err == nil {
					return t, true
				}
			} else if sqlVal.Type == sqlparser.IntVal {
				if n, err := strconv.ParseInt(s, 10, 64); err == nil {
					if n > 1e12 {
						return time.UnixMilli(n), true
					}
					return time.Unix(n, 0), true
				}
			}
		}
	}
	return time.Time{}, false
}

// inferScalarType determines the ScalarType for a column by examining the
// actual record data first, then falling back to expression type analysis.
func inferScalarType(colName string, expr sqlparser.Expr, records []*tsstore.StoredRecord) string {
	switch colName {
	case "time":
		return "TIMESTAMP"
	case "measure_name":
		return "VARCHAR"
	case "measure_value::double":
		return "DOUBLE"
	case "measure_value::bigint":
		return "BIGINT"
	case "measure_value::boolean":
		return "BOOLEAN"
	case "measure_value::timestamp":
		return "TIMESTAMP"
	case "measure_value::varchar":
		return "VARCHAR"
	case "measure_value":
		if len(records) > 0 {
			switch records[0].MeasureValueType {
			case tsstore.MeasureValueTypeDouble:
				return "DOUBLE"
			case tsstore.MeasureValueTypeBigint:
				return "BIGINT"
			case tsstore.MeasureValueTypeBoolean:
				return "BOOLEAN"
			case tsstore.MeasureValueTypeTimestamp:
				return "TIMESTAMP"
			}
		}
		return "VARCHAR"
	}

	if len(records) > 0 {
		for _, dim := range records[0].Dimensions {
			if dim.Name == colName {
				return "VARCHAR"
			}
		}
	}

	if fn, ok := expr.(*sqlparser.FuncExpr); ok && fn.IsAggregate() {
		return "DOUBLE"
	}

	if sqlVal, ok := expr.(*sqlparser.SQLVal); ok {
		switch sqlVal.Type {
		case sqlparser.IntVal:
			return "INTEGER"
		case sqlparser.FloatVal:
			return "DOUBLE"
		case sqlparser.StrVal:
			return "VARCHAR"
		}
	}

	return "VARCHAR"
}

func (s *TimestreamQueryService) applyQuery(ctx context.Context, selectStmt *sqlparser.Select, records []*tsstore.StoredRecord) []map[string]interface{} {
	var rows []map[string]interface{}

	for i, record := range records {
		if i&0xFF == 0 {
			select {
			case <-ctx.Done():
				return rows
			default:
			}
		}
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
