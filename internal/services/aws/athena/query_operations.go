package athena

import (
	"context"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	athenastore "vorpalstacks/internal/store/aws/athena"

	"github.com/google/uuid"
	"vorpalstacks/internal/core/resilience"
)

const (
	maxResultRows      = 1000
	maxQueryStringSize = 262144
)

// StartQueryExecution starts a new query execution in Athena.
func (s *AthenaService) StartQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString == "" {
		return nil, ErrInvalidRequestException
	}

	if len(queryString) > maxQueryStringSize {
		return nil, ErrInvalidRequestException
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	wg, err := st.workGroupStore.GetWorkGroup(workGroup)
	if err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	var queryExecutionContext *athenastore.QueryExecutionContext
	contextMap := request.GetMapParamCaseInsensitive(req.Parameters, "QueryExecutionContext")
	if contextMap != nil {
		queryExecutionContext = &athenastore.QueryExecutionContext{}
		if db, ok := contextMap["Database"].(string); ok {
			queryExecutionContext.Database = db
		}
		if catalog, ok := contextMap["Catalog"].(string); ok {
			queryExecutionContext.Catalog = catalog
		}
	}

	var resultConfiguration *athenastore.ResultConfiguration
	resultConfigMap := request.GetMapParamCaseInsensitive(req.Parameters, "ResultConfiguration")

	if wg.Configuration != nil && wg.Configuration.EnforceWorkGroupConfiguration {
		if wg.Configuration.ResultConfiguration != nil {
			resultConfiguration = wg.Configuration.ResultConfiguration
		}
	} else if resultConfigMap != nil {
		resultConfiguration = s.parseResultConfiguration(resultConfigMap)
	} else if wg.Configuration != nil && wg.Configuration.ResultConfiguration != nil {
		resultConfiguration = wg.Configuration.ResultConfiguration
	}

	statementType := s.detectStatementType(queryString)

	now := time.Now().UTC()
	queryExecution := &athenastore.QueryExecution{
		QueryExecutionId:      uuid.New().String(),
		Query:                 queryString,
		StatementType:         statementType,
		WorkGroup:             workGroup,
		QueryExecutionContext: queryExecutionContext,
		ResultConfiguration:   resultConfiguration,
		Status: &athenastore.QueryExecutionStatus{
			State:              athenastore.QueryExecutionStateQueued,
			SubmissionDateTime: now,
		},
		Statistics: &athenastore.QueryExecutionStatistics{
			EngineExecutionTimeInMillis:   0,
			DataScannedInBytes:            0,
			TotalExecutionTimeInMillis:    0,
			QueryQueueTimeInMillis:        0,
			QueryPlanningTimeInMillis:     0,
			ServiceProcessingTimeInMillis: 0,
		},
	}

	if err := st.queryExecutionStore.CreateQueryExecution(queryExecution); err != nil {
		return nil, err
	}

	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer func() { resilience.RecoverPanic("athena async query") }()
		s.executeQueryAsync(reqCtx, queryExecution)
	}()

	return map[string]interface{}{
		"QueryExecutionId": queryExecution.QueryExecutionId,
	}, nil
}

// GetQueryExecution retrieves the details of a query execution.
func (s *AthenaService) GetQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionId := request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId")
	if queryExecutionId == "" {
		return nil, ErrInvalidRequestException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryExecution, err := st.queryExecutionStore.GetQueryExecution(queryExecutionId)
	if err != nil {
		if err == athenastore.ErrQueryExecutionNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return map[string]interface{}{
		"QueryExecution": s.queryExecutionToResponse(queryExecution),
	}, nil
}

// StopQueryExecution stops a running or queued query execution.
func (s *AthenaService) StopQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionId := request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId")
	if queryExecutionId == "" {
		return nil, ErrInvalidRequestException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryExecution, err := st.queryExecutionStore.GetQueryExecution(queryExecutionId)
	if err != nil {
		if err == athenastore.ErrQueryExecutionNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	if queryExecution.Status.State == athenastore.QueryExecutionStateRunning ||
		queryExecution.Status.State == athenastore.QueryExecutionStateQueued {
		if cancelFn, ok := s.getAndRemoveCancelFunc(queryExecutionId); ok {
			cancelFn()
		}
		queryExecution.Status.State = athenastore.QueryExecutionStateCancelled
		queryExecution.Status.CompletionDateTime = time.Now().UTC()
		if err := st.queryExecutionStore.UpdateQueryExecution(queryExecution); err != nil {
			return nil, err
		}
	} else if queryExecution.Status.State != athenastore.QueryExecutionStateCancelled {
		return nil, awserrors.NewBadRequestException("Query execution is not in a cancellable state")
	}

	return response.EmptyResponse(), nil
}

// ListQueryExecutions returns a list of query executions.
func (s *AthenaService) ListQueryExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")

	maxResults := 50
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 {
			maxResults = val
		}
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		if val, err := strconv.Atoi(nextToken); err == nil && val >= 0 {
			offset = val
		}
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allIds, err := st.queryExecutionStore.ListQueryExecutionIDs(workGroup, 0)
	if err != nil {
		return nil, err
	}

	var ids []string
	for i, id := range allIds {
		if i < offset {
			continue
		}
		if len(ids) >= maxResults {
			break
		}
		ids = append(ids, id)
	}

	result := map[string]interface{}{
		"QueryExecutionIds": ids,
	}

	if offset+maxResults < len(allIds) {
		result["NextToken"] = strconv.Itoa(offset + maxResults)
	}

	return result, nil
}

// ListQueryExecutionsWithPagination returns a list of query executions with pagination support.
func (s *AthenaService) ListQueryExecutionsWithPagination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.ListQueryExecutions(ctx, reqCtx, req)
}

// BatchGetQueryExecution retrieves details for multiple query executions in a single call.
func (s *AthenaService) BatchGetQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionIdsRaw := request.GetArrayParam(req.Parameters, "QueryExecutionIds")
	if len(queryExecutionIdsRaw) == 0 {
		return nil, ErrInvalidRequestException
	}

	if len(queryExecutionIdsRaw) > 50 {
		return nil, ErrInvalidRequestException
	}

	var queryExecutionIds []string
	for _, id := range queryExecutionIdsRaw {
		if idStr, ok := id.(string); ok {
			queryExecutionIds = append(queryExecutionIds, idStr)
		}
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var queryExecutions []map[string]interface{}
	var unprocessedIds []map[string]interface{}

	for _, id := range queryExecutionIds {
		queryExecution, err := st.queryExecutionStore.GetQueryExecution(id)
		if err != nil {
			unprocessedIds = append(unprocessedIds, map[string]interface{}{
				"QueryExecutionId": id,
			})
			continue
		}
		queryExecutions = append(queryExecutions, s.queryExecutionToResponse(queryExecution))
	}

	return map[string]interface{}{
		"QueryExecutions":              queryExecutions,
		"UnprocessedQueryExecutionIds": unprocessedIds,
	}, nil
}

// GetQueryResults retrieves the results of a completed query execution.
func (s *AthenaService) GetQueryResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionId := request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId")
	if queryExecutionId == "" {
		return nil, ErrInvalidRequestException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryExecution, err := st.queryExecutionStore.GetQueryExecution(queryExecutionId)
	if err != nil {
		if err == athenastore.ErrQueryExecutionNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	if queryExecution.Status.State != athenastore.QueryExecutionStateSucceeded {
		return nil, ErrInvalidRequestException
	}

	result, err := st.resultStore.GetResult(queryExecutionId)
	if err != nil {
		return nil, err
	}

	nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken")
	offset := 0
	if nextToken != "" {
		if val, err := strconv.Atoi(nextToken); err == nil {
			offset = val
		}
	}

	maxRows := maxResultRows
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 && val < maxResultRows {
			maxRows = val
		}
	}

	if result.ResultSet == nil || len(result.ResultSet.Rows) == 0 {
		return map[string]interface{}{
			"ResultSet": map[string]interface{}{
				"Rows":              []interface{}{},
				"ResultSetMetadata": map[string]interface{}{"ColumnInfo": []interface{}{}},
			},
			"QueryExecutionId": queryExecutionId,
			"NextToken":        "",
		}, nil
	}

	end := offset + maxRows
	if end > len(result.ResultSet.Rows) {
		end = len(result.ResultSet.Rows)
	}

	pagedRows := result.ResultSet.Rows[offset:end]
	responseNextToken := ""
	if end < len(result.ResultSet.Rows) {
		responseNextToken = strconv.Itoa(end)
	}

	var rows []map[string]interface{}
	for _, row := range pagedRows {
		var data []map[string]interface{}
		for _, datum := range row.Data {
			data = append(data, map[string]interface{}{
				"VarCharValue": datum.VarCharValue,
			})
		}
		rows = append(rows, map[string]interface{}{
			"Data": data,
		})
	}

	var columnInfo []map[string]interface{}
	if result.ResultSet.ResultSetMetadata != nil {
		for _, col := range result.ResultSet.ResultSetMetadata.ColumnInfo {
			columnInfo = append(columnInfo, map[string]interface{}{
				"Label":         col.Label,
				"Name":          col.Name,
				"Type":          col.Type,
				"Precision":     col.Precision,
				"Scale":         col.Scale,
				"Nullable":      col.Nullable,
				"CaseSensitive": col.CaseSensitive,
				"SchemaName":    col.SchemaName,
				"TableName":     col.TableName,
				"CatalogName":   col.CatalogName,
			})
		}
	}

	return map[string]interface{}{
		"UpdateCount": 0,
		"ResultSet": map[string]interface{}{
			"Rows": rows,
			"ResultSetMetadata": map[string]interface{}{
				"ColumnInfo": columnInfo,
			},
		},
		"QueryExecutionId": queryExecutionId,
		"NextToken":        responseNextToken,
	}, nil
}

// GetQueryRuntimeStatistics retrieves runtime statistics for a query execution.
func (s *AthenaService) GetQueryRuntimeStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionId := request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId")
	if queryExecutionId == "" {
		return nil, ErrInvalidRequestException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryExecution, err := st.queryExecutionStore.GetQueryExecution(queryExecutionId)
	if err != nil {
		if err == athenastore.ErrQueryExecutionNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	if queryExecution.Status.State != athenastore.QueryExecutionStateSucceeded {
		return nil, ErrInvalidRequestException
	}

	return map[string]interface{}{
		"QueryRuntimeStatistics": map[string]interface{}{
			"Timeline": map[string]interface{}{
				"QueryQueueTimeInMillis":        queryExecution.Statistics.QueryQueueTimeInMillis,
				"QueryPlanningTimeInMillis":     queryExecution.Statistics.QueryPlanningTimeInMillis,
				"EngineExecutionTimeInMillis":   queryExecution.Statistics.EngineExecutionTimeInMillis,
				"ServiceProcessingTimeInMillis": queryExecution.Statistics.ServiceProcessingTimeInMillis,
				"TotalExecutionTimeInMillis":    queryExecution.Statistics.TotalExecutionTimeInMillis,
			},
			"Rows": map[string]interface{}{
				"InputRows":   queryExecution.Statistics.DataScannedInBytes / 100,
				"InputBytes":  queryExecution.Statistics.DataScannedInBytes,
				"OutputRows":  queryExecution.Statistics.DataScannedInBytes / 100,
				"OutputBytes": queryExecution.Statistics.DataScannedInBytes / 10,
			},
			"OutputStage": map[string]interface{}{
				"StageId":         0,
				"State":           "FINISHED",
				"Done":            true,
				"Nodes":           1,
				"TotalSplits":     1,
				"QueuedSplits":    0,
				"CompletedSplits": 1,
				"RuntimeInMillis": queryExecution.Statistics.EngineExecutionTimeInMillis,
			},
		},
	}, nil
}
