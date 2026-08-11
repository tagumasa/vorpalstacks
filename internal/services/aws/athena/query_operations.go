package athena

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	athenastore "vorpalstacks/internal/store/aws/athena"

	"github.com/google/uuid"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
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

	clientRequestToken := request.GetParamCaseInsensitive(req.Parameters, "ClientRequestToken")
	if clientRequestToken != "" {
		if err := validateClientRequestToken(clientRequestToken); err != nil {
			return nil, err
		}
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Idempotency: if a ClientRequestToken was provided and a previous
	// StartQueryExecution with the same token already succeeded, return
	// the existing QueryExecutionId instead of creating a duplicate.
	if clientRequestToken != "" {
		if existingId, found := st.queryExecutionStore.GetQueryExecutionIdByClientRequestToken(clientRequestToken); found {
			return map[string]interface{}{
				"QueryExecutionId": existingId,
			}, nil
		}
	}

	wg, err := st.workGroupStore.GetWorkGroup(workGroup)
	if err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return nil, workGroupNotFound(workGroup)
		}
		return nil, err
	}

	if wg.State == athenastore.WorkGroupStateDisabled {
		return nil, awserrors.NewAWSError("InvalidRequestException",
			fmt.Sprintf("WorkGroup %s is DISABLED.", workGroup), http.StatusBadRequest)
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
		resultConfiguration, err = s.parseResultConfiguration(resultConfigMap)
		if err != nil {
			return nil, err
		}
	} else if wg.Configuration != nil && wg.Configuration.ResultConfiguration != nil {
		resultConfiguration = wg.Configuration.ResultConfiguration
	}

	var bytesScannedCutoff int64
	if wg.Configuration != nil {
		bytesScannedCutoff = wg.Configuration.BytesScannedCutoffPerQuery
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

	if clientRequestToken != "" {
		if err := st.queryExecutionStore.StoreClientRequestToken(clientRequestToken, queryExecution.QueryExecutionId); err != nil {
			logs.Warn("Failed to store ClientRequestToken mapping", logs.Err(err))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.setCancelFunc(queryExecution.QueryExecutionId, cancel)

	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer func() { resilience.RecoverPanic("athena async query") }()
		defer cancel()
		defer s.getAndRemoveCancelFunc(queryExecution.QueryExecutionId)
		s.executeQueryAsync(reqCtx, ctx, queryExecution, bytesScannedCutoff)
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
			return nil, queryExecutionNotFound(queryExecutionId)
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
			return nil, queryExecutionNotFound(queryExecutionId)
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
		return nil, ErrInvalidRequestException
	}

	return response.EmptyResponse(), nil
}

// ListQueryExecutions returns a list of query executions.
func (s *AthenaService) ListQueryExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")

	maxResults, err := validateMaxResults(req.Parameters, athenaListMaxResults, 0, 50)
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allIds, err := st.queryExecutionStore.ListQueryExecutionIDs(workGroup, 0)
	if err != nil {
		return nil, err
	}

	sort.Strings(allIds)

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	pageResult := pagination.PaginateSlice(allIds, marker, maxResults, func(id string) string {
		return id
	})

	return pagination.BuildListResponse("QueryExecutionIds", pageResult.Items, pageResult.NextMarker), nil
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
			return nil, queryExecutionNotFound(queryExecutionId)
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

	maxRows := maxResultRows
	if val, ok := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults"); ok && val > 0 && val < maxResultRows {
		maxRows = val
	}

	if result.ResultSet == nil || len(result.ResultSet.Rows) == 0 {
		return map[string]interface{}{
			"ResultSet": map[string]interface{}{
				"Rows":              []interface{}{},
				"ResultSetMetadata": map[string]interface{}{"ColumnInfo": []interface{}{}},
			},
			"QueryExecutionId": queryExecutionId,
		}, nil
	}

	type indexedRow struct {
		idx int
		row map[string]interface{}
	}

	allRows := make([]indexedRow, len(result.ResultSet.Rows))
	for i, row := range result.ResultSet.Rows {
		var data []map[string]interface{}
		for _, datum := range row.Data {
			data = append(data, map[string]interface{}{
				"VarCharValue": datum.VarCharValue,
			})
		}
		allRows[i] = indexedRow{
			idx: i,
			row: map[string]interface{}{"Data": data},
		}
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	pageResult := pagination.PaginateSlice(allRows, marker, maxRows, func(item indexedRow) string {
		return strconv.Itoa(item.idx)
	})

	rows := make([]map[string]interface{}, len(pageResult.Items))
	for i, item := range pageResult.Items {
		rows[i] = item.row
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

	resp := map[string]interface{}{
		"UpdateCount": 0,
		"ResultSet": map[string]interface{}{
			"Rows": rows,
			"ResultSetMetadata": map[string]interface{}{
				"ColumnInfo": columnInfo,
			},
		},
		"QueryExecutionId": queryExecutionId,
	}
	if pageResult.NextMarker != "" {
		resp["NextToken"] = pageResult.NextMarker
	}
	return resp, nil
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
			return nil, queryExecutionNotFound(queryExecutionId)
		}
		return nil, err
	}

	if queryExecution.Status.State != athenastore.QueryExecutionStateSucceeded {
		return nil, ErrInvalidRequestException
	}

	outputRows := int64(0)
	outputBytes := int64(0)
	if result, err := st.resultStore.GetResult(queryExecutionId); err == nil && result != nil && result.ResultSet != nil {
		outputRows = int64(len(result.ResultSet.Rows))
		for _, row := range result.ResultSet.Rows {
			for _, d := range row.Data {
				outputBytes += int64(len(d.VarCharValue))
			}
		}
	}

	dataScanned := int64(0)
	if queryExecution.Statistics != nil {
		dataScanned = queryExecution.Statistics.DataScannedInBytes
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
				// InputRows is approximate — our engine does not track source row
				// counts separately from output rows. InputBytes equals
				// DataScannedInBytes (the source data volume). OutputBytes is
				// computed from the actual result set content.
				"InputRows":   outputRows,
				"InputBytes":  dataScanned,
				"OutputRows":  outputRows,
				"OutputBytes": outputBytes,
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
