package athena

import (
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

func (s *AthenaService) executeCreateDatabase(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	catalog := "AwsDataCatalog"
	if context != nil && context.Catalog != "" {
		catalog = context.Catalog
	}

	dbName := s.extractDatabaseNameFromDDL(queryString, "CREATE DATABASE", "CREATE SCHEMA")
	if dbName == "" {
		return nil, nil, fmt.Errorf("database name not specified")
	}

	db := &athenastore.Database{
		Name:        dbName,
		Description: "",
		Parameters:  map[string]string{"EXTERNAL": "TRUE"},
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	if err := stores.databaseStore.CreateDatabase(catalog, db); err != nil {
		return nil, nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &athenastore.ResultSet{
			Rows:              []athenastore.Row{},
			ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: []athenastore.ColumnInfo{}},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        00,
		}, nil
}

func (s *AthenaService) executeDropDatabase(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	catalog := "AwsDataCatalog"
	if context != nil && context.Catalog != "" {
		catalog = context.Catalog
	}

	dbName := s.extractDatabaseNameFromDDL(queryString, "DROP DATABASE", "DROP SCHEMA")
	if dbName == "" {
		return nil, nil, fmt.Errorf("database name not specified")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	if err := stores.databaseStore.DeleteDatabase(catalog, dbName); err != nil {
		return nil, nil, fmt.Errorf("failed to drop database: %w", err)
	}

	return &athenastore.ResultSet{
			Rows:              []athenastore.Row{},
			ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: []athenastore.ColumnInfo{}},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        00,
		}, nil
}

func (s *AthenaService) executeCreateTable(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
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

	tableName, parsedDB, columns, location, format := s.parseCreateTableStatementWithLocation(queryString)
	if tableName == "" {
		return nil, nil, fmt.Errorf("table name not specified")
	}

	if parsedDB != "" {
		database = parsedDB
	}

	table := &athenastore.TableMetadata{
		Name:         tableName,
		DatabaseName: database,
		TableType:    "EXTERNAL_TABLE",
		Columns:      columns,
		Parameters:   map[string]string{"EXTERNAL": "TRUE"},
	}

	if location != "" {
		table.Parameters["LOCATION"] = location
	}
	if format != "" {
		table.Parameters["format"] = format
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	if err := stores.tableStore.CreateTable(catalog, database, table); err != nil {
		return nil, nil, fmt.Errorf("failed to create table: %w", err)
	}

	if location == "" || !strings.HasPrefix(location, "s3://") {
		storedTable := &athenastore.StoredTable{
			DatabaseName: database,
			TableName:    tableName,
			Columns:      columns,
			Rows:         []*athenastore.StoredRow{},
		}
		if err := stores.tableDataStore.StoreTableData(catalog, database, tableName, storedTable); err != nil {
			logs.Warn("failed to store table data after CREATE TABLE", logs.String("table", tableName), logs.Err(err))
		}
	}

	return &athenastore.ResultSet{
			Rows:              []athenastore.Row{},
			ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: []athenastore.ColumnInfo{}},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        0,
		}, nil
}

func (s *AthenaService) executeDropTable(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
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

	tableName := s.extractTableNameFromDrop(queryString)
	if tableName == "" {
		return nil, nil, fmt.Errorf("table name not specified")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	if err := stores.tableStore.DeleteTable(catalog, database, tableName); err != nil {
		return nil, nil, fmt.Errorf("failed to drop table: %w", err)
	}

	if err := stores.tableDataStore.DeleteTableData(catalog, database, tableName); err != nil {
		logs.Warn("failed to delete table data after DROP TABLE", logs.String("table", tableName), logs.Err(err))
	}

	return &athenastore.ResultSet{
			Rows:              []athenastore.Row{},
			ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: []athenastore.ColumnInfo{}},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        0,
		}, nil
}

func (s *AthenaService) executeInsert(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
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

	tableName, columnNames, rawValues := s.parseInsertStatement(queryString)
	if tableName == "" {
		return nil, nil, fmt.Errorf("table name not specified")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	tableData, err := stores.tableDataStore.GetTableData(catalog, database, tableName)
	if err != nil {
		return nil, nil, fmt.Errorf("table not found: %w", err)
	}

	for _, rowValues := range rawValues {
		newRow := &athenastore.StoredRow{Values: make(map[string]interface{})}

		rowVals, ok := rowValues.([]interface{})
		if !ok {
			rowVals = []interface{}{rowValues}
		}

		if len(columnNames) > 0 {
			for i, colName := range columnNames {
				if i < len(rowVals) {
					newRow.Values[colName] = rowVals[i]
				}
			}
		} else {
			for i, col := range tableData.Columns {
				if i < len(rowVals) {
					newRow.Values[col.Name] = rowVals[i]
				}
			}
		}

		tableData.Rows = append(tableData.Rows, newRow)
	}

	if err := stores.tableDataStore.StoreTableData(catalog, database, tableName, tableData); err != nil {
		return nil, nil, fmt.Errorf("failed to store data: %w", err)
	}

	return &athenastore.ResultSet{
			Rows:              []athenastore.Row{},
			ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: []athenastore.ColumnInfo{}},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        00,
		}, nil
}

func (s *AthenaService) executeShowDatabases(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	catalog := "AwsDataCatalog"
	if context != nil && context.Catalog != "" {
		catalog = context.Catalog
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	databases, err := stores.databaseStore.ListDatabases(catalog)
	if err != nil {
		return nil, nil, err
	}

	var rows []athenastore.Row
	for _, db := range databases {
		rows = append(rows, athenastore.Row{
			Data: []athenastore.Datum{{VarCharValue: db.Name}},
		})
	}

	if catalog == "AwsDataCatalog" {
		rows = append([]athenastore.Row{{Data: []athenastore.Datum{{VarCharValue: "default"}}}}, rows...)
	}

	return &athenastore.ResultSet{
			Rows: rows,
			ResultSetMetadata: &athenastore.ResultSetMetadata{
				ColumnInfo: []athenastore.ColumnInfo{
					{Name: "Database", Type: "string"},
				},
			},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        0,
		}, nil
}

func (s *AthenaService) executeShowTables(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
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

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	tables, err := stores.tableStore.ListTables(catalog, database)
	if err != nil {
		return nil, nil, err
	}

	var rows []athenastore.Row
	for _, tbl := range tables {
		rows = append(rows, athenastore.Row{
			Data: []athenastore.Datum{{VarCharValue: tbl.Name}},
		})
	}

	return &athenastore.ResultSet{
			Rows: rows,
			ResultSetMetadata: &athenastore.ResultSetMetadata{
				ColumnInfo: []athenastore.ColumnInfo{
					{Name: "Table", Type: "string"},
				},
			},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        0,
		}, nil
}

func (s *AthenaService) executeDescribe(reqCtx *request.RequestContext, queryString string, context *athenastore.QueryExecutionContext, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
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

	upperQuery := strings.ToUpper(queryString)
	prefix := "DESCRIBE "
	if strings.HasPrefix(upperQuery, "DESC ") {
		prefix = "DESC "
	}

	tableName := strings.TrimSpace(queryString[len(prefix):])
	if strings.Contains(tableName, ";") {
		tableName = strings.TrimSpace(strings.Split(tableName, ";")[0])
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	table, err := stores.tableStore.GetTable(catalog, database, tableName)
	if err != nil {
		return nil, nil, fmt.Errorf("table not found: %w", err)
	}

	var rows []athenastore.Row
	for _, col := range table.Columns {
		rows = append(rows, athenastore.Row{
			Data: []athenastore.Datum{
				{VarCharValue: col.Name},
				{VarCharValue: col.Type},
				{VarCharValue: col.Comment},
			},
		})
	}

	return &athenastore.ResultSet{
			Rows: rows,
			ResultSetMetadata: &athenastore.ResultSetMetadata{
				ColumnInfo: []athenastore.ColumnInfo{
					{Name: "Column", Type: "string"},
					{Name: "Type", Type: "string"},
					{Name: "Comment", Type: "string"},
				},
			},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        0,
		}, nil
}
