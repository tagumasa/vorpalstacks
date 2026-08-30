package athena

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"

	"vorpalstacks/internal/core/logs"
)

const (
	maxResultRows             = 1000
	maxQueryStringSize        = 262144
	queryHistoryRetentionDays = 45
	testModeRetentionHours    = 2
	athenaListMaxResults      = 50
)

// StartQueryExecution starts a new query execution in Athena.
func (s *AthenaService) StartQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := StartQueryExecutionInput{
		QueryString:           request.GetParamCaseInsensitive(req.Parameters, "QueryString"),
		WorkGroup:             request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		ClientRequestToken:    request.GetParamCaseInsensitive(req.Parameters, "ClientRequestToken"),
		QueryExecutionContext: request.GetMapParamCaseInsensitive(req.Parameters, "QueryExecutionContext"),
		ResultConfiguration:   request.GetMapParamCaseInsensitive(req.Parameters, "ResultConfiguration"),
	}

	queryExecutionId, err := s.startQueryExecutionCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"QueryExecutionId": queryExecutionId,
	}, nil
}

// GetQueryExecution retrieves the details of a query execution.
func (s *AthenaService) GetQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionId := request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId")

	queryExecution, err := s.getQueryExecutionCore(reqCtx, queryExecutionId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"QueryExecution": s.queryExecutionToResponse(queryExecution),
	}, nil
}

// StopQueryExecution stops a running or queued query execution.
func (s *AthenaService) StopQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionId := request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId")

	if err := s.stopQueryExecutionCore(reqCtx, queryExecutionId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListQueryExecutions returns a list of query executions.
func (s *AthenaService) ListQueryExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListQueryExecutionsInput{
		WorkGroup:     request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	ids, nextMarker, err := s.listQueryExecutionsCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return pagination.BuildListResponse("QueryExecutionIds", ids, nextMarker), nil
}

// BatchGetQueryExecution retrieves details for multiple query executions in a single call.
func (s *AthenaService) BatchGetQueryExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := BatchGetQueryExecutionInput{
		QueryExecutionIds: request.GetArrayParam(req.Parameters, "QueryExecutionIds"),
	}

	queryExecutions, unprocessedIds, err := s.batchGetQueryExecutionCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	var executionResponses []map[string]interface{}
	for _, qe := range queryExecutions {
		executionResponses = append(executionResponses, s.queryExecutionToResponse(qe))
	}

	var unprocessedResponses []map[string]interface{}
	for _, id := range unprocessedIds {
		unprocessedResponses = append(unprocessedResponses, map[string]interface{}{
			"QueryExecutionId": id,
		})
	}

	return map[string]interface{}{
		"QueryExecutions":              executionResponses,
		"UnprocessedQueryExecutionIds": unprocessedResponses,
	}, nil
}

// GetQueryResults retrieves the results of a completed query execution.
func (s *AthenaService) GetQueryResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := GetQueryResultsInput{
		QueryExecutionId: request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId"),
		MaxResults:       maxResults,
		HasMaxResults:    hasMaxResults,
		NextToken:        pagination.GetMarker(req.Parameters, "NextToken"),
	}

	result, err := s.getQueryResultsCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	if result.EmptyResultSet {
		return map[string]interface{}{
			"ResultSet": map[string]interface{}{
				"Rows":              []interface{}{},
				"ResultSetMetadata": map[string]interface{}{"ColumnInfo": []interface{}{}},
			},
			"QueryExecutionId": input.QueryExecutionId,
		}, nil
	}

	var columnInfo []map[string]interface{}
	if result.Result.ResultSet.ResultSetMetadata != nil {
		for _, col := range result.Result.ResultSet.ResultSetMetadata.ColumnInfo {
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

	resp := map[string]interface{}{
		"UpdateCount": 0,
		"ResultSet": map[string]interface{}{
			"Rows": result.PageRows,
			"ResultSetMetadata": map[string]interface{}{
				"ColumnInfo": columnInfo,
			},
		},
		"QueryExecutionId": input.QueryExecutionId,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// GetQueryRuntimeStatistics retrieves runtime statistics for a query execution.
func (s *AthenaService) GetQueryRuntimeStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryExecutionId := request.GetParamCaseInsensitive(req.Parameters, "QueryExecutionId")

	result, err := s.getQueryRuntimeStatisticsCore(reqCtx, queryExecutionId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"QueryRuntimeStatistics": map[string]interface{}{
			"Timeline": map[string]interface{}{
				"QueryQueueTimeInMillis":        result.Statistics.QueryQueueTimeInMillis,
				"QueryPlanningTimeInMillis":     result.Statistics.QueryPlanningTimeInMillis,
				"EngineExecutionTimeInMillis":   result.Statistics.EngineExecutionTimeInMillis,
				"ServiceProcessingTimeInMillis": result.Statistics.ServiceProcessingTimeInMillis,
				"TotalExecutionTimeInMillis":    result.Statistics.TotalExecutionTimeInMillis,
			},
			"Rows": map[string]interface{}{
				// InputRows: source row tracking is not implemented; 0 is the
				// honest value until it is. InputBytes equals DataScannedInBytes
				// (the source data volume). OutputBytes is computed from the
				// actual result set content.
				"InputRows":   int64(0),
				"InputBytes":  result.DataScannedInByte,
				"OutputRows":  result.OutputRows,
				"OutputBytes": result.OutputBytes,
			},
			"OutputStage": map[string]interface{}{
				"StageId":         0,
				"State":           "FINISHED",
				"Done":            true,
				"Nodes":           1,
				"TotalSplits":     1,
				"QueuedSplits":    0,
				"CompletedSplits": 1,
				"RuntimeInMillis": result.Statistics.EngineExecutionTimeInMillis,
			},
		},
	}, nil
}

// cleanupExpiredQueryExecutions removes query executions older than the AWS-specified
// 45-day retention period. Athena keeps query history for 45 days per AWS documentation.
// In testMode, a shorter retention (testModeRetentionHours) is applied to prevent
// unbounded growth across repeated test runs without destroying in-flight executions.
func (s *AthenaService) cleanupExpiredQueryExecutions(st *athenaStores) {
	cutoff := time.Now().UTC().AddDate(0, 0, -queryHistoryRetentionDays)
	retentionLabel := fmt.Sprintf("%d days", queryHistoryRetentionDays)
	if s.testMode {
		cutoff = time.Now().UTC().Add(-time.Duration(testModeRetentionHours) * time.Hour)
		retentionLabel = fmt.Sprintf("%d hours", testModeRetentionHours)
	}
	deleted, deletedIds, err := st.queryExecutionStore.DeleteExpiredQueryExecutions(cutoff)
	if err != nil {
		return
	}
	st.resultStore.DeleteResultsByIDs(deletedIds)
	if deleted > 0 {
		logs.Info(fmt.Sprintf("athena: cleaned up %d expired query executions (older than %s)", deleted, retentionLabel))
	}
}
