package timestreamquery

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/pkg/sqlparser"
)

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
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.queryCore(ctx, stores, QueryInput{
		QueryString:  request.GetParamCaseInsensitive(req.Parameters, "QueryString"),
		MaxRowsRaw:   request.GetParamCaseInsensitive(req.Parameters, "MaxRows"),
		NextTokenRaw: request.GetParamCaseInsensitive(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	formattedRows := s.formatRowsForResponse(result.Rows, result.ColumnInfo)

	response := map[string]interface{}{
		"QueryId":    result.QueryID,
		"Rows":       formattedRows,
		"ColumnInfo": result.ColumnInfo,
		"QueryStatus": map[string]interface{}{
			"ProgressPercentage":     100.0,
			"CumulativeBytesScanned": int64(0),
			"CumulativeBytesMetered": int64(0),
		},
	}

	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}

// CancelQuery cancels a running query.
func (s *TimestreamQueryService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.cancelQueryCore(stores, request.GetParamCaseInsensitive(req.Parameters, "QueryId")); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CancellationMessage": "Query has been cancelled",
	}, nil
}

// PrepareQuery prepares a query for execution and validates its syntax.
func (s *TimestreamQueryService) PrepareQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.prepareQueryCore(PrepareQueryInput{
		QueryString:     request.GetParamCaseInsensitive(req.Parameters, "QueryString"),
		ValidateOnlyRaw: req.Parameters["ValidateOnly"],
	})
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"QueryString": result.QueryString,
	}

	if result.ValidateOnly {
		// ValidateOnly=true: syntax validated, skip column/parameter extraction.
		return response, nil
	}

	response["Parameters"] = result.Parameters
	response["Columns"] = result.Columns

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
			if colName == "time" {
				scalarType = "TIMESTAMP"
			} else if strings.HasPrefix(colName, "measure_value::") {
				if strings.Contains(colName, "double") {
					scalarType = "DOUBLE"
				} else if strings.Contains(colName, "bigint") {
					scalarType = "BIGINT"
				} else if strings.Contains(colName, "boolean") {
					scalarType = "BOOLEAN"
				} else if strings.Contains(colName, "timestamp") {
					scalarType = "TIMESTAMP"
				}
			} else if colName == "measure_name" {
				scalarType = "VARCHAR"
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
