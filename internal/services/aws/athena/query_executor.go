package athena

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

func (s *AthenaService) executeQueryAsync(reqCtx *request.RequestContext, ctx context.Context, qe *athenastore.QueryExecution, bytesScannedCutoff int64) {
	defer func() {
		if r := resilience.RecoverPanic("executeQueryAsync"); r != nil {
			qe.Status.State = athenastore.QueryExecutionStateFailed
			qe.Status.StateChangeReason = fmt.Sprintf("internal panic: %v", r)
			qe.Status.CompletionDateTime = time.Now().UTC()
			st, _ := s.store(reqCtx)
			if st != nil {
				_ = st.queryExecutionStore.UpdateQueryExecution(qe)
			}
		}
	}()
	startTime := time.Now().UTC()

	st, err := s.store(reqCtx)
	if err != nil {
		qe.Status.State = athenastore.QueryExecutionStateFailed
		qe.Status.StateChangeReason = err.Error()
		qe.Status.CompletionDateTime = time.Now().UTC()
		qe.Status.AthenaError = &athenastore.AthenaError{
			ErrorCategory: 2,
			ErrorType:     "INTERNAL_ERROR",
			Retryable:     true,
			ErrorMessage:  err.Error(),
		}
		qe.Statistics = &athenastore.QueryExecutionStatistics{
			TotalExecutionTimeInMillis: time.Since(startTime).Milliseconds(),
			ResultReuseInformation:     &athenastore.ResultReuseInformation{ReusedPreviousResult: false},
		}
		logs.Error("Failed to get store in executeQueryAsync", logs.Err(err))
		return
	}
	// Check if StopQueryExecution already cancelled the query while it
	// was in QUEUED state (before this goroutine reached RUNNING).
	current, err := st.queryExecutionStore.GetQueryExecution(qe.QueryExecutionId)
	if err != nil {
		logs.Error("Failed to re-read query execution before RUNNING transition", logs.String("id", qe.QueryExecutionId), logs.Err(err))
	} else if current.Status.State == athenastore.QueryExecutionStateCancelled {
		return
	}

	qe.Status.State = athenastore.QueryExecutionStateRunning
	qe.Status.StateChangeReason = ""
	if err := st.queryExecutionStore.UpdateQueryExecution(qe); err != nil {
		logs.Error("Failed to update query execution to RUNNING", logs.String("id", qe.QueryExecutionId), logs.Err(err))
	}

	if strings.Contains(qe.Query, "/* SLOW */") {
		select {
		case <-ctx.Done():
			qe.Status.State = athenastore.QueryExecutionStateCancelled
			qe.Status.CompletionDateTime = time.Now().UTC()
			if err := st.queryExecutionStore.UpdateQueryExecution(qe); err != nil {
				logs.Error("Failed to update query execution to CANCELLED", logs.String("id", qe.QueryExecutionId), logs.Err(err))
			}
			return
		case <-time.After(200 * time.Millisecond):
		}
	}

	result, stats, err := s.executeSQLQuery(reqCtx, ctx, qe.Query, qe.QueryExecutionContext)

	if ctx.Err() != nil {
		qe.Status.State = athenastore.QueryExecutionStateCancelled
		qe.Status.CompletionDateTime = time.Now().UTC()
		qe.Statistics = &athenastore.QueryExecutionStatistics{
			TotalExecutionTimeInMillis: time.Since(startTime).Milliseconds(),
			ResultReuseInformation:     &athenastore.ResultReuseInformation{ReusedPreviousResult: false},
		}
		if err := st.queryExecutionStore.UpdateQueryExecution(qe); err != nil {
			logs.Error("Failed to update query execution to CANCELLED after context cancellation", logs.String("id", qe.QueryExecutionId), logs.Err(err))
		}
		return
	}

	if err != nil {
		qe.Status.State = athenastore.QueryExecutionStateFailed
		qe.Status.StateChangeReason = err.Error()
		qe.Status.CompletionDateTime = time.Now().UTC()
		if qe.Status.AthenaError == nil {
			qe.Status.AthenaError = &athenastore.AthenaError{
				ErrorCategory: 1,
				ErrorType:     "SYNTAX_ERROR",
				Retryable:     false,
				ErrorMessage:  err.Error(),
			}
		}
	} else {
		qe.Status.State = athenastore.QueryExecutionStateSucceeded
		qe.Status.CompletionDateTime = time.Now().UTC()

		queryResult := &athenastore.QueryResult{
			QueryExecutionId: qe.QueryExecutionId,
			ResultSet:        result,
		}
		if err := st.resultStore.StoreResult(qe.QueryExecutionId, queryResult); err != nil {
			logs.Error("Failed to store query result", logs.String("id", qe.QueryExecutionId), logs.Err(err))
		}

		if qe.ResultConfiguration != nil && qe.ResultConfiguration.OutputLocation != "" &&
			qe.StatementType == athenastore.StatementTypeDML {
			if writeErr := s.writeQueryResultsToS3(ctx, reqCtx.GetRegion(), qe.QueryExecutionId, queryResult, qe.ResultConfiguration.OutputLocation); writeErr != nil {
				logs.Warn("Failed to write query results to S3", logs.Err(writeErr))
			}
		}
	}

	endTime := time.Now().UTC()
	dataScanned := int64(0)
	queryPlanningTime := int64(0)
	if stats != nil {
		dataScanned = stats.DataScannedInBytes
		queryPlanningTime = stats.QueryPlanningTimeInMillis
	}

	if bytesScannedCutoff > 0 && dataScanned > bytesScannedCutoff {
		qe.Status.State = athenastore.QueryExecutionStateCancelled
		qe.Status.StateChangeReason = fmt.Sprintf("Query exceeded the data scan limit of %d bytes configured for the workgroup.", bytesScannedCutoff)
		qe.Status.CompletionDateTime = time.Now().UTC()
		st.resultStore.DeleteResultsByIDs([]string{qe.QueryExecutionId})
		logs.Warn("Query cancelled due to BytesScannedCutoffPerQuery exceeded",
			logs.String("id", qe.QueryExecutionId),
			logs.Int64("scanned", dataScanned),
			logs.Int64("cutoff", bytesScannedCutoff))
	}

	qe.Statistics = &athenastore.QueryExecutionStatistics{
		EngineExecutionTimeInMillis:   endTime.Sub(startTime).Milliseconds(),
		DataScannedInBytes:            dataScanned,
		TotalExecutionTimeInMillis:    endTime.Sub(qe.Status.SubmissionDateTime).Milliseconds(),
		QueryQueueTimeInMillis:        0,
		QueryPlanningTimeInMillis:     queryPlanningTime,
		ServiceProcessingTimeInMillis: 0,
		ResultReuseInformation:        &athenastore.ResultReuseInformation{ReusedPreviousResult: false},
	}

	if err := st.queryExecutionStore.UpdateQueryExecution(qe); err != nil {
		logs.Error("Failed to update query execution final state", logs.String("id", qe.QueryExecutionId), logs.Err(err))
	}
}

func (s *AthenaService) executeSQLQuery(reqCtx *request.RequestContext, ctx context.Context, queryString string, context *athenastore.QueryExecutionContext) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	startTime := time.Now().UTC()

	upperQuery := strings.ToUpper(strings.TrimSpace(queryString))

	if strings.HasPrefix(upperQuery, "CREATE DATABASE") || strings.HasPrefix(upperQuery, "CREATE SCHEMA") {
		return s.executeCreateDatabase(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "DROP DATABASE") || strings.HasPrefix(upperQuery, "DROP SCHEMA") {
		return s.executeDropDatabase(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "CREATE EXTERNAL TABLE") {
		return s.executeCreateTable(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "CREATE TABLE") {
		return s.executeCreateTable(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "DROP TABLE") {
		return s.executeDropTable(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "INSERT INTO") {
		return s.executeInsert(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "SHOW DATABASES") || strings.HasPrefix(upperQuery, "SHOW SCHEMAS") {
		return s.executeShowDatabases(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "SHOW TABLES") {
		return s.executeShowTables(reqCtx, queryString, context, startTime)
	}
	if strings.HasPrefix(upperQuery, "DESCRIBE") || strings.HasPrefix(upperQuery, "DESC ") {
		return s.executeDescribe(reqCtx, queryString, context, startTime)
	}

	opts := sqlparser.ParserOptions{
		Dialect: sqlparser.DialectAthena,
	}

	processedSQL := sqlparser.NewAthenaPreprocessor().Process(queryString)

	stmt, err := sqlparser.ParseWithOptions(processedSQL, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("SQL parse error: %w", err)
	}

	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		if unionStmt, ok := stmt.(*sqlparser.Union); ok {
			return s.executeUnion(reqCtx, ctx, unionStmt, context, startTime)
		}
		return nil, nil, fmt.Errorf("only SELECT and UNION statements are supported")
	}

	catalog := "AwsDataCatalog"
	database := "default"
	if context != nil {
		if context.Catalog != "" {
			catalog = context.Catalog
		}
		if context.Database != "" {
			database = context.Database
		}
	}

	hasJoin := s.hasJoin(selectStmt)
	hasGroupBy := len(selectStmt.GroupBy) > 0
	hasAggregates := s.hasAggregateFunctions(selectStmt.SelectExprs)

	if hasJoin || hasGroupBy || hasAggregates {
		tableData, err := s.executeAdvancedSelect(reqCtx, selectStmt, context)
		if err != nil {
			return nil, nil, err
		}

		return s.buildResultSetFromStoredTable(tableData, startTime)
	}

	tableName, err := s.extractTableName(selectStmt)
	if err != nil {
		return s.executeSelectWithoutFrom(selectStmt, startTime)
	}

	if tableName == "dual" {
		return s.executeSelectWithoutFrom(selectStmt, startTime)
	}

	tableData, err := s.getTableData(reqCtx, catalog, database, tableName)
	if err != nil {
		return nil, nil, fmt.Errorf("TABLE_NOT_FOUND: table %s.%s not found in catalog %s", database, tableName, catalog)
	}

	rows, err := s.applyQuery(selectStmt, tableData)
	if err != nil {
		return nil, nil, err
	}

	columnInfo := s.buildColumnInfoFromSelect(selectStmt)
	var columns []athenastore.Column
	if len(columnInfo) > 0 {
		for _, ci := range columnInfo {
			columns = append(columns, athenastore.Column{Name: ci.Name, Type: ci.Type})
		}
	} else {
		columns = tableData.Columns
	}

	projectedTable := &athenastore.StoredTable{
		Columns: columns,
		Rows:    rows,
	}

	return s.buildResultSetFromStoredTable(projectedTable, startTime)
}

func (s *AthenaService) getTableData(reqCtx *request.RequestContext, catalog, database, tableName string) (*athenastore.StoredTable, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	table, err := stores.tableStore.GetTable(catalog, database, tableName)
	if err != nil {
		if err != athenastore.ErrTableNotFound {
			return nil, fmt.Errorf("failed to get table metadata for %s.%s: %w", database, tableName, err)
		}
		return stores.tableDataStore.GetTableData(catalog, database, tableName)
	}

	if table.Parameters != nil {
		if location, ok := table.Parameters["LOCATION"]; ok && strings.HasPrefix(location, "s3://") {
			if s.hasS3Support() {
				data, err := s.loadExternalTableData(reqCtx, catalog, database, tableName)
				if err != nil {
					return nil, fmt.Errorf("failed to load external table data from S3: %w", err)
				}
				return s.convertS3DataToStoredTable(data, table.Columns), nil
			}
		}
	}

	storedTable, err := stores.tableDataStore.GetTableData(catalog, database, tableName)
	if err != nil {
		return &athenastore.StoredTable{
			DatabaseName: database,
			TableName:    tableName,
			Columns:      table.Columns,
			Rows:         []*athenastore.StoredRow{},
		}, nil
	}
	return storedTable, nil
}

// executeSelectWithoutFrom handles SELECT statements without a FROM clause (e.g. SELECT 1, SELECT 1+2).
func (s *AthenaService) executeSelectWithoutFrom(selectStmt *sqlparser.Select, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	var columns []athenastore.ColumnInfo
	var headerData []athenastore.Datum
	var rowData []athenastore.Datum

	for _, expr := range selectStmt.SelectExprs {
		aliased, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		colName := s.extractColumnName(aliased.Expr)
		outputName := colName
		if !aliased.As.IsEmpty() {
			outputName = aliased.As.String()
		}
		if outputName == "" {
			outputName = "_col0"
		}

		colType := "varchar"
		if sqlVal, ok := aliased.Expr.(*sqlparser.SQLVal); ok {
			switch sqlVal.Type {
			case sqlparser.IntVal:
				colType = "integer"
			case sqlparser.FloatVal:
				colType = "double"
			}
		}

		columns = append(columns, athenastore.ColumnInfo{
			Name:  outputName,
			Label: outputName,
			Type:  colType,
		})
		headerData = append(headerData, athenastore.Datum{VarCharValue: outputName})

		val := s.evaluateExpr(aliased.Expr)
		rowData = append(rowData, athenastore.Datum{VarCharValue: val})
	}

	rows := []athenastore.Row{
		{Data: headerData},
		{Data: rowData},
	}

	return &athenastore.ResultSet{
			Rows:              rows,
			ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: columns},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        0,
		}, nil
}

// evaluateExpr evaluates a simple SQL expression to a string value.
func (s *AthenaService) evaluateExpr(expr sqlparser.Expr) string {
	switch e := expr.(type) {
	case *sqlparser.SQLVal:
		return string(e.Val)
	case *sqlparser.ColName:
		return e.Name.String()
	default:
		return sqlparser.String(expr)
	}
}

func (s *AthenaService) extractTableName(selectStmt *sqlparser.Select) (string, error) {
	if len(selectStmt.From) == 0 {
		return "", fmt.Errorf("no table specified in FROM clause")
	}

	aliasedTableExpr, ok := selectStmt.From[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return "", fmt.Errorf("unsupported table expression")
	}

	tableNameExpr, ok := aliasedTableExpr.Expr.(sqlparser.TableName)
	if !ok {
		return "", fmt.Errorf("unsupported table name format")
	}

	return tableNameExpr.Name.String(), nil
}

func (s *AthenaService) applyQuery(selectStmt *sqlparser.Select, tableData *athenastore.StoredTable) ([]*athenastore.StoredRow, error) {
	var rows []*athenastore.StoredRow

	for _, row := range tableData.Rows {
		if selectStmt.Where != nil {
			match, err := s.evaluateWhere(selectStmt.Where.Expr, row.Values)
			if err != nil {
				return nil, err
			}
			if !match {
				continue
			}
		}
		rows = append(rows, row)
	}

	if len(selectStmt.OrderBy) > 0 {
		s.applyOrderBy(rows, selectStmt.OrderBy)
	}

	if selectStmt.Limit != nil {
		limited, err := s.applyLimit(rows, selectStmt.Limit)
		if err != nil {
			return nil, err
		}
		rows = limited
	}

	return s.projectColumns(rows, selectStmt.SelectExprs), nil
}

// executeUnion handles UNION / UNION ALL / UNION DISTINCT set operations.
// Each side is executed as an independent query and results are merged.
// pkg/sqlparser does not support INTERSECT or EXCEPT.
func (s *AthenaService) executeUnion(reqCtx *request.RequestContext, ctx context.Context, unionStmt *sqlparser.Union, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	leftSQL := sqlparser.String(unionStmt.Left)
	rightSQL := sqlparser.String(unionStmt.Right)

	leftResult, leftStats, err := s.executeSQLQuery(reqCtx, ctx, leftSQL, context)
	if err != nil {
		return nil, nil, fmt.Errorf("UNION left side error: %w", err)
	}

	rightResult, rightStats, err := s.executeSQLQuery(reqCtx, ctx, rightSQL, context)
	if err != nil {
		return nil, nil, fmt.Errorf("UNION right side error: %w", err)
	}

	merged := s.mergeUnionResults(leftResult, rightResult, unionStmt.Type)

	dataScanned := int64(0)
	if leftStats != nil {
		dataScanned += leftStats.DataScannedInBytes
	}
	if rightStats != nil {
		dataScanned += rightStats.DataScannedInBytes
	}

	return merged, &athenastore.QueryExecutionStatistics{
		QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
		DataScannedInBytes:        dataScanned,
	}, nil
}

// mergeUnionResults combines two result sets based on the set operation type.
// Only UNION and UNION ALL are supported — INTERSECT/EXCEPT are not parsed
// by pkg/sqlparser.
func (s *AthenaService) mergeUnionResults(left, right *athenastore.ResultSet, unionType string) *athenastore.ResultSet {
	upperType := strings.ToUpper(strings.TrimSpace(unionType))

	var columnInfo []athenastore.ColumnInfo
	if left.ResultSetMetadata != nil {
		columnInfo = left.ResultSetMetadata.ColumnInfo
	} else if right.ResultSetMetadata != nil {
		columnInfo = right.ResultSetMetadata.ColumnInfo
	}

	makeRowKey := func(row athenastore.Row) string {
		var parts []string
		for _, d := range row.Data {
			parts = append(parts, d.VarCharValue)
		}
		return strings.Join(parts, "\x00")
	}

	// UNION ALL: concatenate without deduplication
	if strings.Contains(upperType, "ALL") {
		rows := make([]athenastore.Row, 0, len(left.Rows)+len(right.Rows))
		rows = append(rows, left.Rows...)
		rows = append(rows, right.Rows...)
		return &athenastore.ResultSet{Rows: rows, ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: columnInfo}}
	}

	// UNION (distinct): concatenate with deduplication
	rows := make([]athenastore.Row, 0, len(left.Rows)+len(right.Rows))
	seen := make(map[string]bool)
	for _, l := range left.Rows {
		key := makeRowKey(l)
		if !seen[key] {
			rows = append(rows, l)
			seen[key] = true
		}
	}
	for _, r := range right.Rows {
		key := makeRowKey(r)
		if !seen[key] {
			rows = append(rows, r)
			seen[key] = true
		}
	}
	return &athenastore.ResultSet{Rows: rows, ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: columnInfo}}
}
