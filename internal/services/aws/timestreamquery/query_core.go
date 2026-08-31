package timestreamquery

import (
	"context"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/pkg/sqlparser"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Query / CancelQuery / PrepareQuery — DTO inputs and Core functions. The
// HTTP handlers lift the wire parameters into the DTOs; validation, query
// execution and query-info persistence live here.
// ---------------------------------------------------------------------------

// QueryInput is the DTO input for the Query operation. The pagination
// members stay in their raw wire form; the Core owns their parsing and
// range validation.
type QueryInput struct {
	QueryString  string
	MaxRowsRaw   string
	NextTokenRaw string
}

// QueryOutput is the DTO result for the Query operation: one page of rows
// with the next-page token when more rows remain.
type QueryOutput struct {
	QueryID    string
	Rows       []map[string]interface{}
	ColumnInfo []ColumnInfo
	NextToken  string
}

// parseMaxRows converts the raw MaxRows wire member into the page size.
// MaxRows targets the integer MaxQueryResults shape: a present value that
// is not an integer is a wire-type violation rejected with
// SerializationException, an integer outside the modelled range is a
// ValidationException, and an absent member keeps the default page size.
func parseMaxRows(raw string) (int, error) {
	if raw == "" {
		return maxQueryRows, nil
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, ErrSerializationException
	}
	if err := validateMaxResultsInRange(val, "MaxRows", rangeMaxQueryResults); err != nil {
		return 0, err
	}
	return val, nil
}

// queryCore validates the query string, executes the query, records the
// query lifecycle in the query-info store and returns one page of rows per
// the MaxRows and NextToken parameters.
func (s *TimestreamQueryService) queryCore(ctx context.Context, stores *tsQueryStores, input QueryInput) (*QueryOutput, error) {
	queryString := input.QueryString
	if queryString == "" {
		return nil, ErrValidationException
	}
	if err := validateQueryString(queryString); err != nil {
		return nil, err
	}

	// MaxRows is parsed before execution: a type-violating value is
	// rejected at deserialisation time in the aws-json protocol, so the
	// query never runs for such a request.
	maxRows, err := parseMaxRows(input.MaxRowsRaw)
	if err != nil {
		return nil, err
	}

	queryID := strings.ReplaceAll(uuid.New().String(), "-", "")
	now := time.Now().UTC()

	queryCtx, queryCancel := context.WithCancel(ctx)
	s.cancelFuncs.Store(queryID, queryCancel)
	defer func() {
		s.cancelFuncs.Delete(queryID)
		queryCancel()
	}()

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

	result, execErr := s.executeSQLQuery(queryCtx, stores, queryString)

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
		logs.Warn("Failed to update query info", logs.String("queryId", queryID), logs.Err(err))
	}

	offset := 0
	if input.NextTokenRaw != "" {
		val, err := strconv.Atoi(input.NextTokenRaw)
		if err != nil || val < 0 {
			return nil, ErrValidationException
		}
		offset = val
	}

	totalRows := len(result.Rows)
	var pagedRows []map[string]interface{}
	for i := offset; i < totalRows && len(pagedRows) < maxRows; i++ {
		pagedRows = append(pagedRows, result.Rows[i])
	}

	output := &QueryOutput{
		QueryID:    queryID,
		Rows:       pagedRows,
		ColumnInfo: result.ColumnInfo,
	}
	if offset+maxRows < totalRows {
		output.NextToken = strconv.Itoa(offset + maxRows)
	}

	return output, nil
}

// cancelQueryCore validates the QueryId, cancels a running query through
// the in-process cancellation registry and persists the cancelled state.
func (s *TimestreamQueryService) cancelQueryCore(stores *tsQueryStores, queryID string) error {
	if queryID == "" {
		return ErrValidationException
	}
	if err := validateQueryID(queryID); err != nil {
		return err
	}

	var queryInfo QueryInfo
	if err := stores.queryInfoStore.Get(queryID, &queryInfo); err != nil {
		return ErrResourceNotFound
	}

	if cancelFn, ok := s.cancelFuncs.LoadAndDelete(queryID); ok {
		if fn, ok := cancelFn.(context.CancelFunc); ok {
			fn()
		}
	}

	queryInfo.Cancelled = true
	queryInfo.Status = QueryStatusCancelled
	queryInfo.CompletionTime = time.Now().UTC()
	if err := stores.queryInfoStore.Put(queryID, &queryInfo); err != nil {
		logs.Warn("Failed to update cancelled query info", logs.Err(err))
	}

	return nil
}

// PrepareQueryInput is the DTO input for the PrepareQuery operation. The
// ValidateOnly member keeps its raw wire form (bool or string); the Core
// owns its parsing.
type PrepareQueryInput struct {
	QueryString     string
	ValidateOnlyRaw interface{}
}

// PrepareQueryOutput is the DTO result for the PrepareQuery operation.
// Parameters and Columns stay nil when ValidateOnly is set: AWS validates
// syntax only and does not compute column metadata in that mode.
type PrepareQueryOutput struct {
	QueryString  string
	ValidateOnly bool
	Parameters   []map[string]interface{}
	Columns      []map[string]interface{}
}

// prepareQueryCore validates the query string, checks its syntax against
// the Timestream SQL dialect and derives the parameter and column metadata.
func (s *TimestreamQueryService) prepareQueryCore(input PrepareQueryInput) (*PrepareQueryOutput, error) {
	queryString := input.QueryString
	if queryString == "" {
		return nil, ErrValidationException
	}
	if err := validateQueryString(queryString); err != nil {
		return nil, err
	}

	// Parse ValidateOnly. When true, AWS validates syntax only and
	// does not compute column metadata.
	validateOnly := false
	if b, ok := input.ValidateOnlyRaw.(bool); ok {
		validateOnly = b
	} else if str, ok := input.ValidateOnlyRaw.(string); ok {
		validateOnly = strings.EqualFold(str, "true")
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

	output := &PrepareQueryOutput{
		QueryString:  queryString,
		ValidateOnly: validateOnly,
	}

	if validateOnly {
		// ValidateOnly=true: syntax validated, skip column/parameter extraction.
		return output, nil
	}

	output.Parameters = s.extractParameters(queryString)

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
	output.Columns = columns

	return output, nil
}
