package timestreamquery

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/pkg/sqlparser"

	"github.com/google/uuid"
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
	QueryID        string                   `json:"queryId"`
	QueryString    string                   `json:"queryString"`
	Status         QueryStatus              `json:"status"`
	SubmitTime     time.Time                `json:"submitTime"`
	CompletionTime time.Time                `json:"completionTime,omitempty"`
	Error          string                   `json:"error,omitempty"`
	Cancelled      bool                     `json:"cancelled"`
	CachedResult   *QueryResult             `json:"cachedResult,omitempty"`
	CachedRows     []map[string]interface{} `json:"cachedRows,omitempty"`
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
func (s *Service) DescribeEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
func (s *Service) Query(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString == "" {
		return nil, ErrValidationException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryID := uuid.New().String()
	now := time.Now().UTC()

	queryInfo := &QueryInfo{
		QueryID:     queryID,
		QueryString: queryString,
		Status:      QueryStatusRunning,
		SubmitTime:  now,
		Cancelled:   false,
	}

	stores.queryInfoStore.Put(queryID, queryInfo)

	result, execErr := s.executeSQLQuery(ctx, reqCtx, queryString)

	var latestInfo QueryInfo
	if getErr := stores.queryInfoStore.Get(queryID, &latestInfo); getErr == nil {
		if latestInfo.Cancelled {
			return nil, fmt.Errorf("query was cancelled")
		}
	}

	if execErr != nil {
		queryInfo.Status = QueryStatusFailed
		queryInfo.Error = execErr.Error()
		queryInfo.CompletionTime = time.Now().UTC()
		stores.queryInfoStore.Put(queryID, queryInfo)
		return nil, execErr
	}

	queryInfo.Status = QueryStatusSucceeded
	queryInfo.CompletionTime = time.Now().UTC()
	queryInfo.CachedResult = result
	queryInfo.CachedRows = result.Rows
	stores.queryInfoStore.Put(queryID, queryInfo)

	maxRows := 100
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxRows"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 {
			maxRows = val
		}
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
func (s *Service) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryID := request.GetParamCaseInsensitive(req.Parameters, "QueryId")
	if queryID == "" {
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
	stores.queryInfoStore.Put(queryID, &queryInfo)

	return map[string]interface{}{
		"CancellationMessage": "Query has been cancelled",
	}, nil
}

// PrepareQuery prepares a query for execution and validates its syntax.
func (s *Service) PrepareQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString == "" {
		return nil, ErrValidationException
	}

	processedSQL := s.preprocessor.Process(queryString)

	opts := sqlparser.ParserOptions{
		Dialect: sqlparser.DialectTimestream,
	}
	_, err := sqlparser.ParseWithOptions(processedSQL, opts)
	if err != nil {
		return nil, fmt.Errorf("SQL parse error: %w", err)
	}

	params := s.extractParameters(queryString)

	return map[string]interface{}{
		"QueryString":        queryString,
		"Parameters":         params,
		"QueryId":            uuid.New().String(),
		"QueryParsingStatus": "SUCCESSFUL",
	}, nil
}

func (s *Service) extractParameters(queryString string) []map[string]interface{} {
	var params []map[string]interface{}

	for i := 0; i < len(queryString); i++ {
		if queryString[i] == '@' {
			j := i + 1
			for j < len(queryString) && (isAlphaNum(queryString[j]) || queryString[j] == '_') {
				j++
			}
			paramName := queryString[i+1 : j]
			params = append(params, map[string]interface{}{
				"Name": paramName,
				"Type": "VARCHAR",
			})
			i = j - 1
		}
	}

	return params
}

// GetQueryStatus returns the status of a query.
func (s *Service) GetQueryStatus(ctx context.Context, reqCtx *request.RequestContext, queryID string) (*QueryInfo, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var queryInfo QueryInfo
	if err := stores.queryInfoStore.Get(queryID, &queryInfo); err != nil {
		return nil, ErrResourceNotFound
	}

	return &queryInfo, nil
}

// ListQueryResults returns the results of a completed query.
func (s *Service) ListQueryResults(ctx context.Context, reqCtx *request.RequestContext, queryID string) (*QueryResult, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var queryInfo QueryInfo
	if err := stores.queryInfoStore.Get(queryID, &queryInfo); err != nil {
		return nil, ErrResourceNotFound
	}

	if queryInfo.Status != QueryStatusSucceeded {
		return nil, fmt.Errorf("query not completed")
	}

	if queryInfo.CachedResult != nil {
		return queryInfo.CachedResult, nil
	}

	if queryInfo.CachedRows == nil {
		return nil, fmt.Errorf("query result not found")
	}

	result := &QueryResult{
		QueryID:    queryID,
		Rows:       queryInfo.CachedRows,
		ColumnInfo: []ColumnInfo{},
	}
	return result, nil
}
