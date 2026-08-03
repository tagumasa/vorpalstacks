package athena

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

func emptyDDLResult(startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics) {
	return &athenastore.ResultSet{
			Rows:              []athenastore.Row{},
			ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: []athenastore.ColumnInfo{}},
		}, &athenastore.QueryExecutionStatistics{
			QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
			DataScannedInBytes:        0,
		}
}

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

	rs, stats := emptyDDLResult(startTime)
	return rs, stats, nil
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

	stores.tableStore.DeleteTablesByDatabase(catalog, dbName)
	stores.tableDataStore.DeleteTableDataByDatabase(catalog, dbName)

	if err := stores.databaseStore.DeleteDatabase(catalog, dbName); err != nil {
		return nil, nil, fmt.Errorf("failed to drop database: %w", err)
	}

	rs, stats := emptyDDLResult(startTime)
	return rs, stats, nil
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

	rs, stats := emptyDDLResult(startTime)
	return rs, stats, nil
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
	if strings.Contains(tableName, ".") {
		parts := strings.Split(tableName, ".")
		switch len(parts) {
		case 2:
			database = strings.Trim(parts[0], "`\"';")
			tableName = strings.Trim(parts[1], "`\"';")
		case 3:
			catalog = strings.Trim(parts[0], "`\"';")
			database = strings.Trim(parts[1], "`\"';")
			tableName = strings.Trim(parts[2], "`\"';")
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	if err := stores.tableStore.DeleteTable(catalog, database, tableName); err != nil {
		return nil, nil, fmt.Errorf("failed to drop table: %w", err)
	}

	if err := stores.tableDataStore.DeleteTableData(catalog, database, tableName); err != nil {
		return nil, nil, fmt.Errorf("failed to delete table data for %s.%s: %w", database, tableName, err)
	}

	rs, stats := emptyDDLResult(startTime)
	return rs, stats, nil
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

	table, tblErr := stores.tableStore.GetTable(catalog, database, tableName)
	if tblErr == nil && table.Parameters != nil {
		if location, ok := table.Parameters["LOCATION"]; ok && strings.HasPrefix(location, "s3://") {
			return nil, nil, fmt.Errorf("INSERT INTO external S3 tables is not supported")
		}
	}

	// Table data may not exist yet if CREATE TABLE didn't initialise it.
	// Only fall back to table metadata columns for ErrTableNotFound; any
	// other storage error must be propagated, not swallowed.
	var columns []athenastore.Column
	tableData, err := stores.tableDataStore.GetTableData(catalog, database, tableName)
	if err != nil {
		if !errors.Is(err, athenastore.ErrTableNotFound) {
			return nil, nil, fmt.Errorf("get table data for %s.%s: %w", database, tableName, err)
		}
		if tblErr != nil {
			return nil, nil, fmt.Errorf("table not found: %w", tblErr)
		}
		for _, col := range table.Columns {
			columns = append(columns, athenastore.Column{Name: col.Name, Type: col.Type})
		}
	} else {
		columns = tableData.Columns
	}

	var newRows []*athenastore.StoredRow
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
			for i, col := range columns {
				if i < len(rowVals) {
					newRow.Values[col.Name] = rowVals[i]
				}
			}
		}

		newRows = append(newRows, newRow)
	}

	if err := stores.tableDataStore.AppendRows(catalog, database, tableName, newRows); err != nil {
		return nil, nil, fmt.Errorf("failed to store data: %w", err)
	}

	rs, stats := emptyDDLResult(startTime)
	return rs, stats, nil
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
		hasDefault := false
		for _, db := range databases {
			if db.Name == "default" {
				hasDefault = true
				break
			}
		}
		if !hasDefault {
			rows = append([]athenastore.Row{{Data: []athenastore.Datum{{VarCharValue: "default"}}}}, rows...)
		}
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
	if strings.HasPrefix(tableName, ".") {
		tableName = strings.TrimPrefix(tableName, ".")
	}
	if strings.Contains(tableName, ".") {
		parts := strings.Split(tableName, ".")
		switch len(parts) {
		case 2:
			database = strings.Trim(parts[0], "`\"';")
			tableName = strings.Trim(parts[1], "`\"';")
		case 3:
			catalog = strings.Trim(parts[0], "`\"';")
			database = strings.Trim(parts[1], "`\"';")
			tableName = strings.Trim(parts[2], "`\"';")
		}
	} else {
		tableName = strings.Trim(tableName, "`\"';")
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
