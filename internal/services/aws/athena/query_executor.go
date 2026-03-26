package athena

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

func (s *Service) executeQueryAsync(reqCtx *request.RequestContext, qe *athenastore.QueryExecution) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	qe.Status.State = athenastore.QueryExecutionStateRunning
	qe.Status.StateChangeReason = ""
	_ = st.queryExecutionStore.UpdateQueryExecution(qe)

	result, stats, err := s.executeSQLQuery(reqCtx, ctx, qe.Query, qe.QueryExecutionContext)
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
		_ = st.resultStore.StoreResult(qe.QueryExecutionId, queryResult)

		if qe.ResultConfiguration != nil && qe.ResultConfiguration.OutputLocation != "" {
			if writeErr := s.writeQueryResultsToS3(ctx, qe.QueryExecutionId, queryResult, qe.ResultConfiguration.OutputLocation); writeErr != nil {
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
	qe.Statistics = &athenastore.QueryExecutionStatistics{
		EngineExecutionTimeInMillis:   endTime.Sub(startTime).Milliseconds(),
		DataScannedInBytes:            dataScanned,
		TotalExecutionTimeInMillis:    endTime.Sub(qe.Status.SubmissionDateTime).Milliseconds(),
		QueryQueueTimeInMillis:        0,
		QueryPlanningTimeInMillis:     queryPlanningTime,
		ServiceProcessingTimeInMillis: 0,
		ResultReuseInformation:        &athenastore.ResultReuseInformation{ReusedPreviousResult: false},
	}

	_ = st.queryExecutionStore.UpdateQueryExecution(qe)
}

func (s *Service) executeSQLQuery(reqCtx *request.RequestContext, ctx context.Context, queryString string, context *athenastore.QueryExecutionContext) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
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

	processedSQL := queryString
	if strings.Contains(queryString, "::") {
		processedSQL = s.convertCastOperator(queryString)
	}

	stmt, err := sqlparser.ParseWithOptions(processedSQL, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("SQL parse error: %w", err)
	}

	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, nil, fmt.Errorf("only SELECT statements are supported")
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
		return nil, nil, err
	}

	tableData, err := s.getTableData(reqCtx, catalog, database, tableName)
	if err != nil {
		columnInfo := s.buildColumnInfoFromSelect(selectStmt)
		return &athenastore.ResultSet{
				Rows:              []athenastore.Row{},
				ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: columnInfo},
			}, &athenastore.QueryExecutionStatistics{
				QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
				DataScannedInBytes:        0,
			}, nil
	}

	rows := s.applyQuery(selectStmt, tableData)

	columnInfo := s.buildColumnInfoFromSelect(selectStmt)
	if len(columnInfo) == 0 && len(tableData.Columns) > 0 {
		for _, col := range tableData.Columns {
			columnInfo = append(columnInfo, athenastore.ColumnInfo{
				Label: col.Name,
				Name:  col.Name,
				Type:  col.Type,
			})
		}
	}

	var resultRows []athenastore.Row
	for _, row := range rows {
		var data []athenastore.Datum
		for _, col := range columnInfo {
			val := ""
			if v, ok := row.Values[col.Name]; ok {
				val = fmt.Sprintf("%v", v)
			}
			data = append(data, athenastore.Datum{VarCharValue: val})
		}
		resultRows = append(resultRows, athenastore.Row{Data: data})
	}

	headerRow := athenastore.Row{}
	for _, col := range columnInfo {
		headerRow.Data = append(headerRow.Data, athenastore.Datum{VarCharValue: col.Name})
	}
	resultRows = append([]athenastore.Row{headerRow}, resultRows...)

	stats := &athenastore.QueryExecutionStatistics{
		QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
		DataScannedInBytes:        int64(len(tableData.Rows) * 100),
	}

	return &athenastore.ResultSet{
		Rows:              resultRows,
		ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: columnInfo},
	}, stats, nil
}

func (s *Service) getTableData(reqCtx *request.RequestContext, catalog, database, tableName string) (*athenastore.StoredTable, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	table, err := stores.tableStore.GetTable(catalog, database, tableName)
	if err != nil {
		return stores.tableDataStore.GetTableData(catalog, database, tableName)
	}

	if table.Parameters != nil {
		if location, ok := table.Parameters["LOCATION"]; ok && strings.HasPrefix(location, "s3://") {
			if s.hasS3Support() {
				data, err := s.loadExternalTableData(reqCtx, catalog, database, tableName)
				if err != nil {
					logs.Warn("Failed to load external table data from S3", logs.Err(err))
				} else {
					return s.convertS3DataToStoredTable(data, table.Columns), nil
				}
			}
		}
	}

	return stores.tableDataStore.GetTableData(catalog, database, tableName)
}

func (s *Service) extractTableName(selectStmt *sqlparser.Select) (string, error) {
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

func (s *Service) applyQuery(selectStmt *sqlparser.Select, tableData *athenastore.StoredTable) []*athenastore.StoredRow {
	var rows []*athenastore.StoredRow

	for _, row := range tableData.Rows {
		if selectStmt.Where != nil {
			if !s.evaluateWhere(selectStmt.Where.Expr, row.Values) {
				continue
			}
		}
		rows = append(rows, row)
	}

	if len(selectStmt.OrderBy) > 0 {
		s.applyOrderBy(rows, selectStmt.OrderBy)
	}

	if selectStmt.Limit != nil {
		rows = s.applyLimit(rows, selectStmt.Limit)
	}

	return s.projectColumns(rows, selectStmt.SelectExprs)
}
