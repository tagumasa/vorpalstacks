package timestreamquery

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/pkg/sqlparser"

	"github.com/google/uuid"
)

// queryIDPattern matches the Smithy QueryId trait: ^[a-zA-Z0-9]+$, length 1-64.
var queryIDPattern = regexp.MustCompile(`^[a-zA-Z0-9]{1,64}$`)

// QueryStatus represents the status of a query execution.
type QueryStatus string

// QueryStatusRunning indicates the query is currently running.
const QueryStatusRunning QueryStatus = "RUNNING"

// QueryStatusSucceeded indicates the query completed successfully.
const QueryStatusSucceeded QueryStatus = "SUCCEEDED"

// QueryStatusFailed indicates the query failed.
const QueryStatusFailed QueryStatus = "FAILED"

// QueryStatusCancelled indicates the query was cancelled.
const QueryStatusCancelled QueryStatus = "CANCELLED"

// QueryStatusTimedOut indicates the query timed out.
const QueryStatusTimedOut QueryStatus = "TIMED_OUT"

// QueryInfo contains information about a query execution.
type QueryInfo struct {
	QueryID        string       `json:"queryId"`
	QueryString    string       `json:"queryString"`
	Status         QueryStatus  `json:"status"`
	SubmitTime     time.Time    `json:"submitTime"`
	CompletionTime time.Time    `json:"completionTime,omitempty"`
	Error          string       `json:"error,omitempty"`
	Cancelled      bool         `json:"cancelled"`
	CachedResult   *QueryResult `json:"cachedResult,omitempty"`
}

// QueryResult contains the results of a query execution.
type QueryResult struct {
	QueryID    string                   `json:"queryId"`
	Rows       []map[string]interface{} `json:"rows"`
	ColumnInfo []ColumnInfo             `json:"columnInfo"`
}

// ColumnInfo contains information about a column in query results.
type ColumnInfo struct {
	Name string         `json:"Name"`
	Type ColumnTypeInfo `json:"Type"`
}

// ColumnTypeInfo contains type information for a column.
type ColumnTypeInfo struct {
	ScalarType string `json:"ScalarType,omitempty"`
}

// DescribeEndpoints returns the endpoints for the Timestream Query service.
func (s *TimestreamQueryService) DescribeEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"Endpoints": []map[string]interface{}{
			{
				"Address":              s.serverHost,
				"CachePeriodInMinutes": 1440,
			},
		},
	}, nil
}

// Query executes a query and returns the results.
func (s *TimestreamQueryService) Query(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString == "" {
		return nil, ErrValidationException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryID := strings.ReplaceAll(uuid.New().String(), "-", "")
	now := time.Now().UTC()

	queryInfo := &QueryInfo{
		QueryID:     queryID,
		QueryString: queryString,
		Status:      QueryStatusRunning,
		SubmitTime:  now,
		Cancelled:   false,
	}

	if err := stores.queryInfoStore.Put(queryID, queryInfo); err != nil {
		logs.Warn("Failed to store query info", logs.Err(err))
	}

	result, execErr := s.executeSQLQuery(ctx, reqCtx, queryString)

	var latestInfo QueryInfo
	if getErr := stores.queryInfoStore.Get(queryID, &latestInfo); getErr == nil {
		if latestInfo.Cancelled {
			logs.Info("Timestream query was cancelled", logs.String("queryId", queryID))
			return nil, ErrQueryExecutionError
		}
	}

	if execErr != nil {
		queryInfo.Status = QueryStatusFailed
		queryInfo.Error = execErr.Error()
		queryInfo.CompletionTime = time.Now().UTC()
		if err := stores.queryInfoStore.Put(queryID, queryInfo); err != nil {
			logs.Warn("Failed to persist failed query info", logs.String("queryId", queryID), logs.Err(err))
		}
		return nil, execErr
	}

	queryInfo.Status = QueryStatusSucceeded
	queryInfo.CompletionTime = time.Now().UTC()
	queryInfo.CachedResult = result
	if err := stores.queryInfoStore.Put(queryID, queryInfo); err != nil {
		logs.Warn("Failed to update query info", logs.Err(err))
	}

	maxRows := maxQueryRows
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxRows"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 {
			maxRows = val
		}
	}
	if maxRows > maxQueryRows {
		maxRows = maxQueryRows
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		if val, err := strconv.Atoi(nextToken); err == nil && val >= 0 {
			offset = val
		}
	}

	totalRows := len(result.Rows)
	var pagedRows []map[string]interface{}
	for i := offset; i < totalRows && len(pagedRows) < maxRows; i++ {
		pagedRows = append(pagedRows, result.Rows[i])
	}

	formattedRows := s.formatRowsForResponse(pagedRows, result.ColumnInfo)

	response := map[string]interface{}{
		"QueryId":    queryID,
		"Rows":       formattedRows,
		"ColumnInfo": result.ColumnInfo,
		"QueryStatus": map[string]interface{}{
			"ProgressPercentage":     100.0,
			"CumulativeBytesScanned": int64(0),
			"CumulativeBytesMetered": int64(0),
		},
	}

	if offset+maxRows < totalRows {
		response["NextToken"] = strconv.Itoa(offset + maxRows)
	}

	return response, nil
}

// CancelQuery cancels a running query.
func (s *TimestreamQueryService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryID := request.GetParamCaseInsensitive(req.Parameters, "QueryId")
	if queryID == "" || !queryIDPattern.MatchString(queryID) {
		return nil, ErrValidationException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var queryInfo QueryInfo
	if err := stores.queryInfoStore.Get(queryID, &queryInfo); err != nil {
		return nil, ErrResourceNotFound
	}

	queryInfo.Cancelled = true
	queryInfo.Status = QueryStatusCancelled
	queryInfo.CompletionTime = time.Now().UTC()
	if err := stores.queryInfoStore.Put(queryID, &queryInfo); err != nil {
		logs.Warn("Failed to update cancelled query info", logs.Err(err))
	}

	return map[string]interface{}{
		"CancellationMessage": "Query has been cancelled",
	}, nil
}

// PrepareQuery prepares a query for execution and validates its syntax.
func (s *TimestreamQueryService) PrepareQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString == "" {
		return nil, ErrValidationException
	}

	// M6: Parse ValidateOnly. When true, AWS validates syntax only and
	// does not compute column metadata.
	validateOnly := false
	if val, ok := req.Parameters["ValidateOnly"]; ok {
		if b, ok := val.(bool); ok {
			validateOnly = b
		} else if s, ok := val.(string); ok {
			validateOnly = strings.EqualFold(s, "true")
		}
	}

	processedSQL := s.preprocessor.Process(queryString)

	opts := sqlparser.ParserOptions{
		Dialect: sqlparser.DialectTimestream,
	}
	stmt, err := sqlparser.ParseWithOptions(processedSQL, opts)
	if err != nil {
		logs.Debug("Timestream prepared statement SQL parse error", logs.String("query", processedSQL), logs.Err(err))
		return nil, ErrValidationException
	}

	response := map[string]interface{}{
		"QueryString": queryString,
	}

	if validateOnly {
		// ValidateOnly=true: syntax validated, skip column/parameter extraction.
		return response, nil
	}

	params := s.extractParameters(queryString)
	response["Parameters"] = params

	var columns []map[string]interface{}
	if selectStmt, ok := stmt.(*sqlparser.Select); ok {
		columnInfo := s.buildColumnInfoForPrepare(selectStmt)
		for _, ci := range columnInfo {
			columns = append(columns, map[string]interface{}{
				"Name": ci.Name,
				"Type": map[string]interface{}{
					"ScalarType": ci.Type.ScalarType,
				},
			})
		}
	}
	response["Columns"] = columns

	return response, nil
}

func (s *TimestreamQueryService) extractParameters(queryString string) []map[string]interface{} {
	seen := make(map[string]string)
	var order []string

	for i := 0; i < len(queryString); i++ {
		if queryString[i] == '@' {
			j := i + 1
			for j < len(queryString) && (isAlphaNum(queryString[j]) || queryString[j] == '_') {
				j++
			}
			paramName := queryString[i+1 : j]

			// Infer type from surrounding SQL context.
			paramType := inferParamType(queryString, j)

			if _, exists := seen[paramName]; !exists {
				order = append(order, paramName)
			}
			// Prefer a more specific type over VARCHAR.
			if existing, exists := seen[paramName]; !exists || (existing == "VARCHAR" && paramType != "VARCHAR") {
				seen[paramName] = paramType
			}
			i = j - 1
		}
	}

	params := make([]map[string]interface{}, 0, len(order))
	for _, name := range order {
		params = append(params, map[string]interface{}{
			"Name": name,
			"Type": seen[name],
		})
	}
	return params
}

// inferParamType attempts to infer a parameter's type from the SQL text
// immediately following the parameter. Timestream uses the :: cast operator
// (e.g. @param::double). Falls back to VARCHAR when no type hint is found.
func inferParamType(queryString string, pos int) string {
	if pos+1 < len(queryString) && queryString[pos] == ':' && queryString[pos+1] == ':' {
		start := pos + 2
		end := start
		for end < len(queryString) && (isAlphaNum(queryString[end]) || queryString[end] == '_') {
			end++
		}
		if end > start {
			castType := strings.ToUpper(queryString[start:end])
			switch castType {
			case "BIGINT", "INTEGER", "INT":
				return "BIGINT"
			case "DOUBLE", "FLOAT":
				return "DOUBLE"
			case "BOOLEAN", "BOOL":
				return "BOOLEAN"
			case "TIMESTAMP", "DATE", "TIME":
				return "TIMESTAMP"
			case "VARCHAR", "STRING":
				return "VARCHAR"
			default:
				return castType
			}
		}
	}
	return "VARCHAR"
}

func (s *TimestreamQueryService) buildColumnInfoForPrepare(selectStmt *sqlparser.Select) []ColumnInfo {
	var columns []ColumnInfo
	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			colName := s.extractColumnName(aliased.Expr)
			if !aliased.As.IsEmpty() {
				colName = aliased.As.String()
			}
			scalarType := "VARCHAR"
			if strings.Contains(strings.ToLower(colName), "time") {
				scalarType = "TIMESTAMP"
			}
			if _, isFunc := aliased.Expr.(*sqlparser.FuncExpr); isFunc {
				if fn, ok := aliased.Expr.(*sqlparser.FuncExpr); ok && fn.IsAggregate() {
					scalarType = "DOUBLE"
				}
			}
			if _, isVal := aliased.Expr.(*sqlparser.SQLVal); isVal {
				scalarType = "INTEGER"
			}
			columns = append(columns, ColumnInfo{
				Name: colName,
				Type: ColumnTypeInfo{ScalarType: scalarType},
			})
		}
	}
	return columns
}
