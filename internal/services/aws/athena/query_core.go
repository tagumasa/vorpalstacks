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
	"vorpalstacks/internal/core/resilience"
	athenastore "vorpalstacks/internal/store/aws/athena"

	"github.com/google/uuid"
)

// --- DTOs (transport-neutral inputs carrying the raw wire members) ---

// StartQueryExecutionInput carries the parsed wire members of a
// StartQueryExecution request. QueryExecutionContext and ResultConfiguration
// travel as the raw wire maps so the Core applies the same parse order the
// handler previously applied inline.
type StartQueryExecutionInput struct {
	QueryString           string
	WorkGroup             string
	ClientRequestToken    string
	QueryExecutionContext map[string]interface{}
	ResultConfiguration   map[string]interface{}
}

// ListQueryExecutionsInput carries the workgroup filter plus the raw
// MaxResults window (presence-flagged) and pagination marker.
type ListQueryExecutionsInput struct {
	WorkGroup     string
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// BatchGetQueryExecutionInput carries the raw QueryExecutionIds wire array;
// the Core applies the entry-count gates and string filtering.
type BatchGetQueryExecutionInput struct {
	QueryExecutionIds []interface{}
}

// GetQueryResultsInput carries the execution id plus the raw MaxResults
// window (presence-flagged) and pagination marker.
type GetQueryResultsInput struct {
	QueryExecutionId string
	MaxResults       int
	HasMaxResults    bool
	NextToken        string
}

// GetQueryResultsResult carries the paged rows and the raw result set so the
// handler serialises column metadata without touching the store.
type GetQueryResultsResult struct {
	EmptyResultSet bool
	PageRows       []map[string]interface{}
	NextMarker     string
	Result         *athenastore.QueryResult
}

// GetQueryRuntimeStatisticsResult carries everything the handler needs to
// serialise the QueryRuntimeStatistics response shape.
type GetQueryRuntimeStatisticsResult struct {
	Statistics        *athenastore.QueryExecutionStatistics
	DataScannedInByte int64
	OutputRows        int64
	OutputBytes       int64
}

// --- Core functions ---

// startQueryExecutionCore validates the request, reserves the
// ClientRequestToken idempotency slot, creates the QUEUED execution record
// and launches the async worker. The HTTP handler and any future admin path
// share this single sequence. The store is acquired only after the request
// has been validated so an invalid request never depends on storage, the
// order the original handler applied.
func (s *AthenaService) startQueryExecutionCore(reqCtx *request.RequestContext, input StartQueryExecutionInput) (string, error) {
	if input.QueryString == "" {
		return "", ErrInvalidRequestException
	}

	// QueryString @length(1-262144) counts Unicode characters.
	if err := validateQueryStringSize(input.QueryString); err != nil {
		return "", err
	}

	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	if input.ClientRequestToken != "" {
		if err := validateClientRequestToken(input.ClientRequestToken); err != nil {
			return "", err
		}
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	// Idempotency: reserve the ClientRequestToken atomically before any
	// validation or creation. StoreClientRequestTokenIfAbsent uses crtMu
	// internally, so concurrent requests with the same token are
	// serialised — no TOCTOU window. If the token is already mapped to
	// a successful query, return that query's ID. If the token is new,
	// it is reserved and a deferred rollback releases it when any
	// subsequent step fails, preventing phantom ID mappings.
	queryExecutionId := uuid.New().String()
	tokenReserved := false
	queryCreated := false

	if input.ClientRequestToken != "" {
		existingId, stored, err := st.queryExecutionStore.StoreClientRequestTokenIfAbsent(input.ClientRequestToken, queryExecutionId)
		if err != nil {
			return "", fmt.Errorf("failed to store ClientRequestToken: %w", err)
		}
		if !stored {
			return existingId, nil
		}
		tokenReserved = true
	}

	defer func() {
		if tokenReserved && !queryCreated {
			st.queryExecutionStore.ReleaseClientRequestToken(input.ClientRequestToken)
		}
	}()

	wg, err := st.workGroupStore.GetWorkGroup(workGroup)
	if err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return "", workGroupNotFound(workGroup)
		}
		return "", err
	}

	if wg.State == athenastore.WorkGroupStateDisabled {
		return "", awserrors.NewAWSError("InvalidRequestException",
			fmt.Sprintf("WorkGroup %s is DISABLED.", workGroup), http.StatusBadRequest)
	}

	var queryExecutionContext *athenastore.QueryExecutionContext
	if input.QueryExecutionContext != nil {
		queryExecutionContext = &athenastore.QueryExecutionContext{}
		if db, ok := input.QueryExecutionContext["Database"].(string); ok {
			queryExecutionContext.Database = db
		}
		if catalog, ok := input.QueryExecutionContext["Catalog"].(string); ok {
			queryExecutionContext.Catalog = catalog
		}
	}

	var resultConfiguration *athenastore.ResultConfiguration

	if wg.Configuration != nil && wg.Configuration.EnforceWorkGroupConfiguration {
		if wg.Configuration.ResultConfiguration != nil {
			resultConfiguration = wg.Configuration.ResultConfiguration
		}
	} else if input.ResultConfiguration != nil {
		resultConfiguration, err = s.parseResultConfiguration(input.ResultConfiguration)
		if err != nil {
			return "", err
		}
	} else if wg.Configuration != nil && wg.Configuration.ResultConfiguration != nil {
		resultConfiguration = wg.Configuration.ResultConfiguration
	}

	var bytesScannedCutoff int64
	if wg.Configuration != nil {
		bytesScannedCutoff = wg.Configuration.BytesScannedCutoffPerQuery
	}

	statementType := s.detectStatementType(input.QueryString)

	now := time.Now().UTC()
	queryExecution := &athenastore.QueryExecution{
		QueryExecutionId:      queryExecutionId,
		Query:                 input.QueryString,
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
		return "", err
	}
	queryCreated = true

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

	return queryExecution.QueryExecutionId, nil
}

// getQueryExecutionCore fetches a query execution, mapping the store
// not-found sentinel onto the API error. The store is acquired after the
// identifier validation, the order the original handler applied.
func (s *AthenaService) getQueryExecutionCore(reqCtx *request.RequestContext, queryExecutionId string) (*athenastore.QueryExecution, error) {
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
	return queryExecution, nil
}

// stopQueryExecutionCore stops a running or queued query execution,
// cancelling the async worker and transitioning the state atomically.
func (s *AthenaService) stopQueryExecutionCore(reqCtx *request.RequestContext, queryExecutionId string) error {
	if queryExecutionId == "" {
		return ErrInvalidRequestException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	queryExecution, err := st.queryExecutionStore.GetQueryExecution(queryExecutionId)
	if err != nil {
		if err == athenastore.ErrQueryExecutionNotFound {
			return queryExecutionNotFound(queryExecutionId)
		}
		return err
	}

	if queryExecution.Status.State == athenastore.QueryExecutionStateRunning ||
		queryExecution.Status.State == athenastore.QueryExecutionStateQueued {
		if cancelFn, ok := s.getAndRemoveCancelFunc(queryExecutionId); ok {
			cancelFn()
		}
		// Transition atomically so the async worker's QUEUED -> RUNNING
		// write can never overwrite the cancelled state after this
		// response returns.
		cancelled, transitioned, err := st.queryExecutionStore.TransitionQueryExecutionState(
			queryExecutionId, athenastore.QueryExecutionStateCancelled,
			athenastore.QueryExecutionStateRunning, athenastore.QueryExecutionStateQueued)
		if err != nil {
			return err
		}
		if !transitioned && cancelled.Status.State != athenastore.QueryExecutionStateCancelled {
			return ErrInvalidRequestException
		}
	} else if queryExecution.Status.State != athenastore.QueryExecutionStateCancelled {
		return ErrInvalidRequestException
	}

	return nil
}

// listQueryExecutionsCore lists query-execution ids for a workgroup with the
// documented window semantics (default 50, range 0-50) applied before the
// store walk, matching the original validation position.
func (s *AthenaService) listQueryExecutionsCore(reqCtx *request.RequestContext, input ListQueryExecutionsInput) ([]string, string, error) {
	maxResults, err := resolveMaxResults(input.MaxResults, input.HasMaxResults, athenaListMaxResults, 0, 50)
	if err != nil {
		return nil, "", err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	allIds, err := st.queryExecutionStore.ListQueryExecutionIDs(input.WorkGroup, 0)
	if err != nil {
		return nil, "", err
	}

	sort.Strings(allIds)

	pageResult := pagination.PaginateSlice(allIds, input.NextToken, maxResults, func(id string) string {
		return id
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// batchGetQueryExecutionCore applies the entry-count gates (1-50 raw ids),
// fetches each execution and partitions found/unprocessed.
func (s *AthenaService) batchGetQueryExecutionCore(reqCtx *request.RequestContext, input BatchGetQueryExecutionInput) ([]*athenastore.QueryExecution, []string, error) {
	if len(input.QueryExecutionIds) == 0 {
		return nil, nil, ErrInvalidRequestException
	}

	if len(input.QueryExecutionIds) > 50 {
		return nil, nil, ErrInvalidRequestException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	var queryExecutionIds []string
	for _, id := range input.QueryExecutionIds {
		if idStr, ok := id.(string); ok {
			queryExecutionIds = append(queryExecutionIds, idStr)
		}
	}

	var found []*athenastore.QueryExecution
	var unprocessed []string

	for _, id := range queryExecutionIds {
		queryExecution, err := st.queryExecutionStore.GetQueryExecution(id)
		if err != nil {
			unprocessed = append(unprocessed, id)
			continue
		}
		found = append(found, queryExecution)
	}

	return found, unprocessed, nil
}

// getQueryResultsCore fetches a succeeded execution's result set, enforces
// the documented MaxResults window (1-1000) at the position the original
// handler read the member and pages the rows by index.
func (s *AthenaService) getQueryResultsCore(reqCtx *request.RequestContext, input GetQueryResultsInput) (*GetQueryResultsResult, error) {
	if input.QueryExecutionId == "" {
		return nil, ErrInvalidRequestException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryExecution, err := st.queryExecutionStore.GetQueryExecution(input.QueryExecutionId)
	if err != nil {
		if err == athenastore.ErrQueryExecutionNotFound {
			return nil, queryExecutionNotFound(input.QueryExecutionId)
		}
		return nil, err
	}

	if queryExecution.Status.State != athenastore.QueryExecutionStateSucceeded {
		return nil, ErrInvalidRequestException
	}

	result, err := st.resultStore.GetResult(input.QueryExecutionId)
	if err != nil {
		return nil, err
	}

	// MaxQueryResults carries the documented 1-1000 range: an explicitly
	// provided out-of-window value is rejected rather than silently
	// clamped to the default page size.
	if input.HasMaxResults && (input.MaxResults < 1 || input.MaxResults > maxResultRows) {
		return nil, awserrors.NewAWSError("InvalidRequestException",
			fmt.Sprintf("MaxResults must be between 1 and %d (got %d)", maxResultRows, input.MaxResults), http.StatusBadRequest)
	}
	maxRows := maxResultRows
	if input.HasMaxResults {
		maxRows = input.MaxResults
	}

	if result.ResultSet == nil || len(result.ResultSet.Rows) == 0 {
		return &GetQueryResultsResult{EmptyResultSet: true}, nil
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

	pageResult := pagination.PaginateSlice(allRows, input.NextToken, maxRows, func(item indexedRow) string {
		return strconv.Itoa(item.idx)
	})

	rows := make([]map[string]interface{}, len(pageResult.Items))
	for i, item := range pageResult.Items {
		rows[i] = item.row
	}

	return &GetQueryResultsResult{
		PageRows:   rows,
		NextMarker: pageResult.NextMarker,
		Result:     result,
	}, nil
}

// getQueryRuntimeStatisticsCore fetches a succeeded execution and computes
// the output-row/output-byte figures from the stored result set.
func (s *AthenaService) getQueryRuntimeStatisticsCore(reqCtx *request.RequestContext, queryExecutionId string) (*GetQueryRuntimeStatisticsResult, error) {
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

	return &GetQueryRuntimeStatisticsResult{
		Statistics:        queryExecution.Statistics,
		DataScannedInByte: dataScanned,
		OutputRows:        outputRows,
		OutputBytes:       outputBytes,
	}, nil
}
