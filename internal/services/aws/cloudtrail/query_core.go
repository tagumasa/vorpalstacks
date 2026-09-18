package cloudtrail

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/core/resilience"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	"vorpalstacks/pkg/sqlparser"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// StartQueryInput carries the StartQuery members. A query is stated either
// through QueryStatement or through the dashboard QueryAlias form; the
// delivery target, when set, receives the finished results.
type StartQueryInput struct {
	QueryStatement  string
	QueryAlias      string
	QueryParameters []string
	DeliveryS3URI   string
	OwnerAccountID  string
}

// GetQueryResultsInput carries the pagination members for GetQueryResults.
type GetQueryResultsInput struct {
	QueryID         string
	MaxQueryResults int
	NextToken       string
}

// DescribeQueryInput carries the query identifier for DescribeQuery.
type DescribeQueryInput struct {
	QueryID string
}

// CancelQueryInput carries the query identifier for CancelQuery.
type CancelQueryInput struct {
	QueryID string
}

// ListQueriesInput carries the filter and pagination members for
// ListQueries. Time members keep both wire forms (RFC3339 string and Unix
// epoch number) because the JSON 1.1 protocol serialises timestamps as
// epochs while query strings arrive as RFC3339 text.
type ListQueriesInput struct {
	EventDataStore string
	QueryStatus    string
	StartTimeStr   string
	StartTimeRaw   interface{}
	EndTimeStr     string
	EndTimeRaw     interface{}
	MaxResults     int
	NextToken      string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// startQueryCore is the single entry point for StartQuery: it validates
// the statement/alias request shape, the delivery target, and the owner
// account, admits the query under the concurrent-query bound, records it
// as RUNNING, and executes it asynchronously under the query deadline.
func (s *CloudTrailService) startQueryCore(ctx context.Context, store cloudtrailstore.CloudTrailStoreInterface, in StartQueryInput) (map[string]interface{}, error) {
	// StartQuery takes either a QueryStatement or the dashboard form
	// (QueryAlias with QueryParameters); the two forms do not mix, and
	// parameters name no referent without an alias.
	if in.QueryAlias != "" && in.QueryStatement != "" {
		return nil, newInvalidParameterException(
			"Specify either QueryStatement or QueryAlias, not both")
	}
	if in.QueryAlias == "" && len(in.QueryParameters) > 0 {
		return nil, newInvalidParameterException(
			"QueryParameters are only valid with QueryAlias")
	}
	if in.QueryStatement == "" && in.QueryAlias == "" {
		return nil, newInvalidQueryStatementException(
			"QueryStatement is required")
	}
	if in.QueryAlias != "" {
		// The alias form resolves dashboard query templates, which land
		// with the dashboard operations; a template cannot be resolved
		// until they exist.
		return nil, newUnsupportedOperationException(
			"QueryAlias requires a dashboard query template, which is not supported")
	}

	pq, err := parseQueryStatement(in.QueryStatement)
	if err != nil {
		return nil, err
	}

	// The owner account, when stated, must own the event data store this
	// store serves; event data stores of other accounts are not queryable
	// here.
	if in.OwnerAccountID != "" && in.OwnerAccountID != store.GetAccountID() {
		return nil, newEventDataStoreNotFoundException(
			fmt.Sprintf("No event data store owned by account %s was found", in.OwnerAccountID))
	}

	// Verify the EDS exists.
	eds, err := store.GetEventDataStore(pq.edsID)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	if eds.Status == "PENDING_DELETION" {
		return nil, newOperationNotPermittedException(
			"Cannot query a PENDING_DELETION event data store")
	}

	// A delivery target must be a well-formed S3 URI naming an existing
	// bucket; malformed URIs fail the model's bucket-name rule and missing
	// buckets fail the operation's declared not-found shape.
	if in.DeliveryS3URI != "" {
		invoker := s.s3Invoker()
		if invoker == nil {
			return nil, newOperationNotPermittedException(
				"S3 query result delivery is not available")
		}
		bucket, _, err := splitDeliveryS3URI(in.DeliveryS3URI)
		if err != nil {
			return nil, err
		}
		exists, err := invoker.BucketExists(ctx, store.GetRegion(), bucket)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, newS3BucketDoesNotExistException(
				fmt.Sprintf("The specified S3 bucket does not exist: %s", bucket))
		}
	}

	queryID := uuid.NewString()
	now := time.Now().UTC()

	// Admit the query record with RUNNING status under the store's
	// concurrent-query bound; a refusal never records the query.
	qr := &cloudtrailstore.QueryRecord{
		QueryID:        queryID,
		EventDataStore: eds.EventDataStoreID,
		QueryStatement: in.QueryStatement,
		QueryStatus:    "RUNNING",
		StartTime:      now,
	}
	if in.DeliveryS3URI != "" {
		qr.DeliveryS3URI = in.DeliveryS3URI
		qr.DeliveryStatus = "PENDING"
	}

	if err := store.AdmitQuery(qr, cloudtrailstore.MaxConcurrentQueries); err != nil {
		return nil, s.mapStoreError(err)
	}

	// Execute the query asynchronously. The terminal save coordinates with
	// CancelQuery through the store mutex: a cancel that lands first has
	// already written CANCELLED, and the executor's verdict never
	// overwrites it.
	deadline := time.Now().Add(cloudtrailstore.LakeQueryDeadline)
	if s.queryDeadline != 0 {
		deadline = time.Now().Add(s.queryDeadline)
	}
	go func() {
		defer func() {
			if r := resilience.RecoverPanic("cloudtrail.StartQuery"); r != nil {
				s.finaliseQueryExecution(store, queryID, nil, nil, false,
					fmt.Errorf("internal error: panic recovered: %v", r))
			}
		}()

		results, stats, timedOut, execErr := s.executeQuery(store, pq, deadline)
		final := s.finaliseQueryExecution(store, queryID, results, stats, timedOut, execErr)
		// Delivery runs only for queries that finished whole: partial
		// results of a TIMED_OUT query are retrievable through
		// GetQueryResults but are never delivered to S3.
		if final != nil && final.QueryStatus == "FINISHED" && final.DeliveryS3URI != "" {
			s.deliverQueryResults(store, final)
		}
	}()

	return map[string]interface{}{
		"QueryId":                      queryID,
		"EventDataStoreOwnerAccountId": store.GetAccountID(),
	}, nil
}

// getQueryResultsCore is the single entry point for GetQueryResults.
func (s *CloudTrailService) getQueryResultsCore(store cloudtrailstore.CloudTrailStoreInterface, in GetQueryResultsInput) (map[string]interface{}, error) {
	if in.QueryID == "" {
		return nil, newInvalidParameterException(
			"QueryId is required")
	}

	qr, err := store.GetQuery(in.QueryID)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	maxResults := in.MaxQueryResults
	if maxResults < 0 {
		return nil, newInvalidMaxResultsException(
			"MaxQueryResults must be between 1 and 1000")
	}
	if maxResults == 0 {
		maxResults = cloudtrailstore.DefaultQueryResultsPageSize
	} else if maxResults > cloudtrailstore.MaxGetQueryResultsResults {
		return nil, newInvalidMaxResultsException(
			fmt.Sprintf("MaxQueryResults exceeds the maximum of %d", cloudtrailstore.MaxGetQueryResultsResults))
	}

	paged := pagination.PaginateSliceByPosition(qr.QueryResultRows, in.NextToken, maxResults)

	// Initialise with make to ensure JSON serialises as [] not null.
	rows := make([]interface{}, 0, len(paged.Items))
	for _, row := range paged.Items {
		rows = append(rows, row)
	}

	result := map[string]interface{}{
		// GetQueryResultsResponse declares QueryStatus, QueryResultRows,
		// QueryStatistics, ErrorMessage and NextToken — no query
		// identifier.
		"QueryStatus":     qr.QueryStatus,
		"QueryResultRows": rows,
		"QueryStatistics": map[string]interface{}{
			"ResultsCount":      qr.ResultsCount,
			"TotalResultsCount": qr.ResultsCount,
			"BytesScanned":      qr.BytesScanned,
		},
	}

	if paged.IsTruncated {
		result["NextToken"] = paged.NextMarker
	}

	if qr.ErrorMessage != "" {
		result["ErrorMessage"] = qr.ErrorMessage
	}

	return result, nil
}

// describeQueryCore is the single entry point for DescribeQuery.
func (s *CloudTrailService) describeQueryCore(store cloudtrailstore.CloudTrailStoreInterface, in DescribeQueryInput) (map[string]interface{}, error) {
	if in.QueryID == "" {
		return nil, newInvalidParameterException(
			"QueryId is required")
	}

	qr, err := store.GetQuery(in.QueryID)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	result := map[string]interface{}{
		"QueryId":     qr.QueryID,
		"QueryStatus": qr.QueryStatus,
		"QueryString": qr.QueryStatement,
		// DescribeQuery reports the QueryStatisticsForDescribeQuery member
		// set: matched/scanned event counts, honest scan bytes, run time,
		// and the query's creation time.
		"QueryStatistics": map[string]interface{}{
			"EventsMatched":         qr.EventsMatched,
			"EventsScanned":         qr.EventsScanned,
			"BytesScanned":          qr.BytesScanned,
			"ExecutionTimeInMillis": qr.ExecutionTimeInMillis,
			"CreationTime":          qr.StartTime,
		},
		"EventDataStoreOwnerAccountId": store.GetAccountID(),
	}

	if qr.DeliveryS3URI != "" {
		result["DeliveryS3Uri"] = qr.DeliveryS3URI
	}
	if qr.DeliveryStatus != "" {
		result["DeliveryStatus"] = qr.DeliveryStatus
	}

	if qr.ErrorMessage != "" {
		result["ErrorMessage"] = qr.ErrorMessage
	}

	return result, nil
}

// finaliseQueryExecution records the terminal outcome of a query execution
// and returns the resulting record. The transition runs under the store
// mutex and only from the RUNNING status: a CancelQuery that landed while
// the executor ran has already written CANCELLED, and that verdict is
// never overwritten — the results are discarded instead. A timed-out
// query keeps the rows the scan produced before the deadline; they remain
// retrievable through GetQueryResults.
func (s *CloudTrailService) finaliseQueryExecution(store cloudtrailstore.CloudTrailStoreInterface, queryID string, results [][]map[string]string, stats *lakeExecutionStats, timedOut bool, execErr error) *cloudtrailstore.QueryRecord {
	final, err := store.MutateQuery(queryID, func(qr *cloudtrailstore.QueryRecord) error {
		if qr.QueryStatus != "RUNNING" {
			return cloudtrailstore.ErrUnchanged
		}
		endTime := time.Now().UTC()
		qr.EndTime = &endTime
		if execErr != nil {
			qr.QueryStatus = "FAILED"
			qr.ErrorMessage = execErr.Error()
			return nil
		}
		if timedOut {
			qr.QueryStatus = "TIMED_OUT"
		} else {
			qr.QueryStatus = "FINISHED"
		}
		qr.QueryResultRows = results
		qr.ResultsCount = int32(len(results))
		if stats != nil {
			qr.EventsMatched = stats.eventsMatched
			qr.EventsScanned = stats.eventsScanned
			qr.BytesScanned = stats.bytesScanned
			qr.ExecutionTimeInMillis = stats.executionTimeMs
		}
		return nil
	})
	if err != nil {
		slog.Error("Failed to save query result", "queryId", queryID, "error", err)
		return nil
	}
	if final != nil && final.QueryStatus == "CANCELLED" {
		slog.Info("Query was cancelled during execution; results discarded", "queryId", queryID)
	}
	return final
}

// cancelQueryCore is the single entry point for CancelQuery. The terminal
// check and the CANCELLED write run under the store mutex as one step, so
// the acknowledgement the client receives is the status the record holds.
func (s *CloudTrailService) cancelQueryCore(store cloudtrailstore.CloudTrailStoreInterface, in CancelQueryInput) (map[string]interface{}, error) {
	if in.QueryID == "" {
		return nil, newInvalidParameterException(
			"QueryId is required")
	}

	qr, err := store.MutateQuery(in.QueryID, func(qr *cloudtrailstore.QueryRecord) error {
		if cloudtrailstore.QueryTerminalStatus(qr.QueryStatus) {
			return newInactiveQueryException(
				"Cannot cancel a query that has already finished, been cancelled, or timed out")
		}
		// A query carrying a delivery target no longer needs it once the
		// query itself is cancelled.
		if qr.DeliveryStatus == "PENDING" {
			qr.DeliveryStatus = "CANCELLED"
		}
		qr.QueryStatus = "CANCELLED"
		now := time.Now().UTC()
		qr.EndTime = &now
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"QueryId":     qr.QueryID,
		"QueryStatus": qr.QueryStatus,
	}, nil
}

// listQueriesCore is the single entry point for ListQueries.
func (s *CloudTrailService) listQueriesCore(store cloudtrailstore.CloudTrailStoreInterface, in ListQueriesInput) (map[string]interface{}, error) {
	edsID := in.EventDataStore
	if edsID == "" {
		return nil, newInvalidParameterException(
			"EventDataStore is required")
	}

	edsID = cloudtrailstore.ExtractEventDataStoreID(edsID)

	queries, err := store.ListQueriesByEDS(edsID)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// Optional filters. StartTime/EndTime bound the listing to queries run
	// within the period; both wire time forms are accepted like the
	// LookupEvents bounds.
	statusFilter := in.QueryStatus
	if statusFilter != "" {
		if err := validateQueryStatus(statusFilter); err != nil {
			return nil, err
		}
	}
	startFilter, err := parseWireTime(in.StartTimeStr, in.StartTimeRaw)
	if err != nil {
		return nil, err
	}
	endFilter, err := parseWireTime(in.EndTimeStr, in.EndTimeRaw)
	if err != nil {
		return nil, err
	}
	// ListQueries declares InvalidDateRangeException — "Be sure that the
	// start time is chronologically before the end time" — for an
	// out-of-order period; the LookupEvents time-range error is a
	// different shape on a different operation.
	if startFilter != nil && endFilter != nil && endFilter.Before(*startFilter) {
		return nil, newInvalidDateRangeException(
			"The start time must be chronologically before the end time")
	}
	maxResults := in.MaxResults
	if maxResults < 0 {
		return nil, newInvalidMaxResultsException(
			"MaxResults must be between 1 and 1000")
	}
	if maxResults == 0 {
		maxResults = cloudtrailstore.DefaultQueryResultsPageSize
	} else if maxResults > cloudtrailstore.MaxListQueriesResults {
		return nil, newInvalidMaxResultsException(
			fmt.Sprintf("MaxResults exceeds the maximum of %d", cloudtrailstore.MaxListQueriesResults))
	}

	// Filter queries by status and run-time period. The listing carries
	// the operation's own window regardless of the caller's period:
	// "Returns a list of queries and query statuses for the past seven
	// days" (ListQueries).
	windowStart := time.Now().UTC().Add(-cloudtrailstore.ListQueriesWindow)
	var filtered []*cloudtrailstore.QueryRecord
	for _, qr := range queries {
		if qr.StartTime.Before(windowStart) {
			continue
		}
		if statusFilter != "" && qr.QueryStatus != statusFilter {
			continue
		}
		if startFilter != nil && qr.StartTime.Before(*startFilter) {
			continue
		}
		if endFilter != nil && qr.StartTime.After(*endFilter) {
			continue
		}
		filtered = append(filtered, qr)
	}

	// Paginate the filtered results.
	paged := pagination.PaginateSliceByPosition(filtered, in.NextToken, maxResults)
	queryList := make([]map[string]interface{}, 0, len(paged.Items))
	for _, qr := range paged.Items {
		// The list item is the Query shape: creation time, identifier,
		// status — nothing else.
		queryList = append(queryList, map[string]interface{}{
			"QueryId":      qr.QueryID,
			"QueryStatus":  qr.QueryStatus,
			"CreationTime": qr.StartTime.Unix(),
		})
	}

	result := map[string]interface{}{
		"Queries": queryList,
	}
	if paged.IsTruncated {
		result["NextToken"] = paged.NextMarker
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// CloudTrail Lake SQL engine
// ---------------------------------------------------------------------------

// queryFieldsPattern extracts the SELECT and FROM portions of a CloudTrail
// Lake SQL-like QueryStatement. Example: "SELECT eventID, eventTime FROM
// <eds-id> WHERE eventName = 'PutObject'".
var (
	selectPattern = regexp.MustCompile(`(?i)^SELECT\s+(.+?)\s+FROM\s+([^\s]+)`)
	wherePattern  = regexp.MustCompile(`(?i)WHERE\s+(.+)$`)
)

// parsedQuery holds the structured representation of a CloudTrail Lake query.
type parsedQuery struct {
	edsID     string
	columns   []string
	whereExpr sqlparser.Expr
}

func parseQueryStatement(stmt string) (*parsedQuery, error) {
	stmt = strings.TrimSpace(stmt)
	if stmt == "" {
		return nil, newInvalidQueryStatementException(
			"QueryStatement is required")
	}
	// The length bound is character-counted: the Smithy length trait
	// measures a string in Unicode scalar values, not bytes.
	if utf8.RuneCountInString(stmt) > cloudtrailstore.MaxQueryStatementChars {
		return nil, newInvalidQueryStatementException(fmt.Sprintf(
			"QueryStatement must contain at most %d characters", cloudtrailstore.MaxQueryStatementChars))
	}

	matches := selectPattern.FindStringSubmatch(stmt)
	if matches == nil {
		return nil, newInvalidQueryStatementException(
			"QueryStatement must contain SELECT ... FROM ...")
	}

	colPart := strings.TrimSpace(matches[1])
	edsID := strings.TrimSpace(matches[2])

	// Extract EDS ID from FROM (strip ARN prefix if present).
	edsID = cloudtrailstore.ExtractEventDataStoreID(edsID)
	// Strip any surrounding quotes or backticks.
	edsID = strings.Trim(edsID, "\"`'")

	var columns []string
	if colPart == "*" {
		columns = []string{"*"}
	} else {
		for _, c := range strings.Split(colPart, ",") {
			c = strings.TrimSpace(c)
			c = strings.Trim(c, "\"`'")
			columns = append(columns, c)
		}
		// Every projection must be a Lake schema column; the row keys echo
		// the statement's own spelling, the values resolve through the
		// vocabulary (case-insensitive).
		for _, c := range columns {
			if !lakeColumns[lakeColumnCanonical(c)] {
				return nil, newInvalidQueryStatementException(
					fmt.Sprintf("Unknown column: %s", c))
			}
		}
	}

	// Parse WHERE clause using pkg/sqlparser PartiQL dialect for full
	// operator support (=, !=, >, <, >=, <=, LIKE, IN, BETWEEN, IS NULL,
	// NOT, AND, OR). Previously only `=` was supported via a single regex.
	// The parse is fail-closed: a syntactically invalid WHERE clause is
	// rejected with StartQuery's declared InvalidQueryStatementException
	// instead of silently degrading to an unfiltered scan.
	var whereExpr sqlparser.Expr
	if wm := wherePattern.FindStringSubmatch(stmt); wm != nil {
		whereRaw := strings.TrimSpace(wm[1])
		// Wrap in a minimal SELECT so the parser has a complete statement.
		parsed, err := sqlparser.ParseWithOptions(
			"SELECT * FROM t WHERE "+whereRaw,
			sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL},
		)
		if err != nil {
			return nil, newInvalidQueryStatementException(
				fmt.Sprintf("Invalid WHERE clause: %v", err))
		}
		sel, ok := parsed.(*sqlparser.Select)
		if !ok || sel.Where == nil {
			return nil, newInvalidQueryStatementException(
				"QueryStatement WHERE clause is not a valid condition")
		}
		whereExpr = sel.Where.Expr
		if err := walkColumnRefs(whereExpr, func(cn *sqlparser.ColName) error {
			qualifier, name := colNameParts(cn)
			return validateLakeColumnRef(qualifier, name)
		}); err != nil {
			return nil, err
		}
	}

	return &parsedQuery{
		edsID:     edsID,
		columns:   columns,
		whereExpr: whereExpr,
	}, nil
}

// lakeExecutionStats carries the honest execution accounting recorded on
// the query: every event the scan examined, the bytes of record JSON it
// read, the events the WHERE clause matched, and the wall-clock time of
// the scan.
type lakeExecutionStats struct {
	eventsScanned   int64
	eventsMatched   int64
	bytesScanned    int64
	executionTimeMs int64
}

// executeQuery runs the parsed query against the event data store named in
// the statement and returns the result rows in CloudTrail Lake format
// ([][]map[string]string), the execution statistics, and whether the
// deadline expired. The walk reads only that store's ingested event copies,
// in ascending event-time order; the WHERE expression is evaluated wholly in
// the engine because the per-EDS buckets carry no secondary indexes to
// pre-filter through. The scan follows the store's NextToken to exhaustion:
// LakeQueryScanBound is the per-page scan size, not a cap on the scanned
// set — a query examines the whole store. When the deadline passes the scan
// stops between pages and the rows it produced so far are the partial
// result.
func (s *CloudTrailService) executeQuery(store cloudtrailstore.CloudTrailStoreInterface, pq *parsedQuery, deadline time.Time) ([][]map[string]string, *lakeExecutionStats, bool, error) {
	query := cloudtrailstore.EDSQuery{
		MaxResults: cloudtrailstore.LakeQueryScanBound,
	}

	started := time.Now()
	stats := &lakeExecutionStats{}
	rows := make([][]map[string]string, 0)
	timedOut := false
	for {
		if !deadline.IsZero() && time.Now().After(deadline) {
			timedOut = true
			break
		}
		events, nextToken, err := store.LookupEDSEvents(pq.edsID, query)
		if err != nil {
			return nil, nil, false, err
		}
		for _, e := range events {
			stats.eventsScanned++
			stats.bytesScanned += int64(len(e.CloudTrailEvent))

			lake := lakeRow(e)

			// Evaluate the full WHERE expression against the Lake row. This
			// supports all SQL operators (=, !=, >, <, LIKE, IN, BETWEEN,
			// IS NULL, NOT, AND, OR) over the Lake column vocabulary. An
			// expression type the evaluator does not cover aborts the whole
			// query: silently dropping the row would shrink the result set
			// without notice, which is fail-open, not fail-closed — the
			// query fails with the expression type named instead.
			if pq.whereExpr != nil {
				matched, err := evaluateWhere(pq.whereExpr, lake)
				if err != nil {
					return nil, nil, false, err
				}
				if !matched {
					continue
				}
			}
			stats.eventsMatched++

			var row []map[string]string
			if len(pq.columns) == 1 && pq.columns[0] == "*" {
				for _, col := range lakeColumnList {
					if v, ok := lake[col]; ok {
						row = append(row, map[string]string{col: lakeValueString(v)})
					}
				}
			} else {
				for _, col := range pq.columns {
					val := lakeValueString(resolveLakeColumn(lake, "", lakeColumnCanonical(col)))
					row = append(row, map[string]string{col: val})
				}
			}
			rows = append(rows, row)
		}
		if nextToken == "" {
			break
		}
		query.NextToken = nextToken
	}
	stats.executionTimeMs = time.Since(started).Milliseconds()

	return rows, stats, timedOut, nil
}

// evaluateWhere evaluates a WHERE expression against a formatted event row.
// Returns (matched, error). An error indicates an unsupported expression
// type; the caller aborts the query on it — a silently skipped row would
// shrink the result set without notice.
func evaluateWhere(expr sqlparser.Expr, row map[string]interface{}) (bool, error) {
	switch e := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return evaluateComparison(e, row), nil
	case *sqlparser.AndExpr:
		left, err := evaluateWhere(e.Left, row)
		if err != nil {
			return false, err
		}
		if !left {
			return false, nil
		}
		return evaluateWhere(e.Right, row)
	case *sqlparser.OrExpr:
		left, err := evaluateWhere(e.Left, row)
		if err != nil {
			return false, err
		}
		if left {
			return true, nil
		}
		return evaluateWhere(e.Right, row)
	case *sqlparser.ParenExpr:
		return evaluateWhere(e.Expr, row)
	case *sqlparser.IsExpr:
		return evaluateIs(e, row), nil
	case *sqlparser.NotExpr:
		result, err := evaluateWhere(e.Expr, row)
		if err != nil {
			return false, err
		}
		return !result, nil
	case *sqlparser.RangeCond:
		return evaluateRangeCond(e, row), nil
	default:
		return false, fmt.Errorf("unsupported WHERE expression type: %T", expr)
	}
}

// evaluateComparison applies SQL three-valued-logic semantics: a NULL
// operand makes every operator — including the negated ones, which would
// otherwise turn an unknown into a match — unsatisfied; IS NULL is the
// null test. Ordering operators additionally require two orderable
// operands: the array columns (the resources fields) only carry equality.
// The IN forms delegate to the per-item membership test.
func evaluateComparison(expr *sqlparser.ComparisonExpr, row map[string]interface{}) bool {
	leftVal := getExprValue(expr.Left, row)
	if leftVal == nil {
		return false
	}
	if tuple, ok := expr.Right.(sqlparser.ValTuple); ok {
		return evaluateInComparison(expr.Operator, leftVal, tuple, row)
	}
	rightVal := getExprValue(expr.Right, row)
	if rightVal == nil {
		return false
	}

	switch expr.Operator {
	case sqlparser.EqualStr:
		return valuesEqual(leftVal, rightVal)
	case sqlparser.NotEqualStr:
		return !valuesEqual(leftVal, rightVal)
	case sqlparser.LessThanStr:
		cmp, ok := compareForOrder(leftVal, rightVal)
		return ok && cmp < 0
	case sqlparser.LessEqualStr:
		cmp, ok := compareForOrder(leftVal, rightVal)
		return ok && cmp <= 0
	case sqlparser.GreaterThanStr:
		cmp, ok := compareForOrder(leftVal, rightVal)
		return ok && cmp > 0
	case sqlparser.GreaterEqualStr:
		cmp, ok := compareForOrder(leftVal, rightVal)
		return ok && cmp >= 0
	case sqlparser.LikeStr:
		return matchLikeAny(leftVal, rightVal)
	case sqlparser.NotLikeStr:
		return !matchLikeAny(leftVal, rightVal)
	}
	return false
}

// evaluateInComparison decides the IN / NOT IN membership forms item by
// item under the same three-valued logic: a NULL list item leaves the
// negated form unknown, so it cannot claim the row either.
func evaluateInComparison(operator string, leftVal interface{}, tuple sqlparser.ValTuple, row map[string]interface{}) bool {
	switch operator {
	case sqlparser.InStr:
		for _, item := range tuple {
			if valuesEqual(leftVal, getExprValue(item, row)) {
				return true
			}
		}
		return false
	case sqlparser.NotInStr:
		for _, item := range tuple {
			itemVal := getExprValue(item, row)
			if itemVal == nil || valuesEqual(leftVal, itemVal) {
				return false
			}
		}
		return true
	}
	return false
}

// evaluateRangeCond applies BETWEEN and NOT BETWEEN with the same
// three-valued-logic semantics: a NULL bound or value, or an array column
// (which carries no order), leaves the range test unknown, and neither the
// positive nor the negated form claims the row.
func evaluateRangeCond(expr *sqlparser.RangeCond, row map[string]interface{}) bool {
	val := getExprValue(expr.Left, row)
	fromCmp, fromOK := compareForOrder(val, getExprValue(expr.From, row))
	toCmp, toOK := compareForOrder(val, getExprValue(expr.To, row))
	if !fromOK || !toOK {
		return false
	}

	switch expr.Operator {
	case sqlparser.BetweenStr:
		return fromCmp >= 0 && toCmp <= 0
	case sqlparser.NotBetweenStr:
		return fromCmp < 0 || toCmp > 0
	}
	return false
}

func evaluateIs(expr *sqlparser.IsExpr, row map[string]interface{}) bool {
	val := getExprValue(expr.Expr, row)
	isNull := val == nil
	switch expr.Operator {
	case sqlparser.IsNullStr:
		return isNull
	case sqlparser.IsNotNullStr:
		return !isNull
	}
	return false
}

func getExprValue(expr sqlparser.Expr, row map[string]interface{}) interface{} {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		qualifier, name := colNameParts(e)
		return resolveLakeColumn(row, qualifier, name)
	case *sqlparser.SQLVal:
		if e.Type == sqlparser.StrVal {
			return string(e.Val)
		} else if e.Type == sqlparser.IntVal {
			if val, err := strconv.ParseInt(string(e.Val), 10, 64); err == nil {
				return val
			}
		} else if e.Type == sqlparser.FloatVal {
			if val, err := strconv.ParseFloat(string(e.Val), 64); err == nil {
				return val
			}
		}
		return string(e.Val)
	case sqlparser.BoolVal:
		return bool(e)
	case *sqlparser.NullVal:
		return nil
	}
	return nil
}

// valuesEqual reports SQL equality of two non-null operand values.
// Array columns (the resources fields) satisfy equality when any entry
// equals the operand; everything else is decided by the ordering
// comparator's zero result — timestamps as time, numbers numerically, and
// the remainder lexically over the Lake string rendering.
func valuesEqual(left, right interface{}) bool {
	if left == nil || right == nil {
		return false
	}
	if vals, ok := left.([]string); ok {
		return anyOfValues(vals, right)
	}
	if vals, ok := right.([]string); ok {
		return anyOfValues(vals, left)
	}
	cmp, _ := compareForOrder(left, right)
	return cmp == 0
}

// compareForOrder orders two operand values for the ordering operators,
// reporting whether both sides are orderable at all: a null operand or an
// array column (the resources fields) carries no order, and every ordering
// comparison on it is unsatisfied rather than satisfied-by-accident.
// Timestamps compare as time when the other operand parses as an epoch or
// a timestamp literal; everything else compares numerically when both
// sides parse as numbers, else lexically over the Lake string rendering.
func compareForOrder(left, right interface{}) (int, bool) {
	if left == nil || right == nil {
		return 0, false
	}
	if _, ok := left.([]string); ok {
		return 0, false
	}
	if _, ok := right.([]string); ok {
		return 0, false
	}
	if lt, ok := left.(time.Time); ok {
		if rt, ok := toTime(right); ok {
			return compareEpoch(lt.Unix(), rt.Unix()), true
		}
	}
	if rt, ok := right.(time.Time); ok {
		if lt, ok := toTime(left); ok {
			return compareEpoch(lt.Unix(), rt.Unix()), true
		}
	}

	leftFloat, leftErr := toFloat(left)
	rightFloat, rightErr := toFloat(right)

	if leftErr == nil && rightErr == nil {
		if leftFloat < rightFloat {
			return -1, true
		} else if leftFloat > rightFloat {
			return 1, true
		}
		return 0, true
	}

	leftStr := lakeValueString(left)
	rightStr := lakeValueString(right)
	if leftStr < rightStr {
		return -1, true
	} else if leftStr > rightStr {
		return 1, true
	}
	return 0, true
}

// anyOfValues reports an array-column equality: it holds when any entry
// equals the operand.
func anyOfValues(vals []string, operand interface{}) bool {
	for _, v := range vals {
		if valuesEqual(v, operand) {
			return true
		}
	}
	return false
}

func compareEpoch(left, right int64) int {
	if left < right {
		return -1
	} else if left > right {
		return 1
	}
	return 0
}

// toTime coerces an operand to a timestamp: time values pass through,
// numbers read as epoch seconds, strings parse as epoch seconds or as a
// timestamp literal in the layouts the Lake record format and AWS query
// examples use.
func toTime(v interface{}) (time.Time, bool) {
	switch val := v.(type) {
	case time.Time:
		return val, true
	case int64:
		return time.Unix(val, 0).UTC(), true
	case int:
		return time.Unix(int64(val), 0).UTC(), true
	case float64:
		return time.Unix(int64(val), 0).UTC(), true
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return time.Unix(int64(f), 0).UTC(), true
		}
		for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
			if t, err := time.Parse(layout, val); err == nil {
				return t.UTC(), true
			}
		}
	}
	return time.Time{}, false
}

func toFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int64:
		return float64(val), nil
	case int:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	}
	return 0, fmt.Errorf("cannot convert %T to float", v)
}

// matchLikeAny applies a LIKE pattern to a column value, rendering both
// sides through the Lake string form; an array column (a resources field)
// matches when any entry matches.
func matchLikeAny(value, pattern interface{}) bool {
	patternStr := lakeValueString(pattern)
	if vals, ok := value.([]string); ok {
		for _, v := range vals {
			if matchLike(v, patternStr) {
				return true
			}
		}
		return false
	}
	return matchLike(lakeValueString(value), patternStr)
}

func matchLike(value, pattern string) bool {
	pattern = strings.Trim(pattern, "'")
	pattern = strings.ToLower(pattern)
	value = strings.ToLower(value)
	regexPattern := likePatternToRegex(pattern)
	matched, _ := regexp.MatchString("^"+regexPattern+"$", value)
	return matched
}

func likePatternToRegex(pattern string) string {
	var result strings.Builder
	escaped := false

	for _, ch := range pattern {
		if escaped {
			switch ch {
			case '%', '_', '\\':
				result.WriteRune(ch)
			default:
				result.WriteRune('\\')
				result.WriteRune(ch)
			}
			escaped = false
			continue
		}

		switch ch {
		case '\\':
			escaped = true
		case '%':
			result.WriteString(".*")
		case '_':
			result.WriteString(".")
		case '.', '^', '$', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|':
			result.WriteRune('\\')
			result.WriteRune(ch)
		default:
			result.WriteRune(ch)
		}
	}

	if escaped {
		result.WriteRune('\\')
	}

	return result.String()
}
