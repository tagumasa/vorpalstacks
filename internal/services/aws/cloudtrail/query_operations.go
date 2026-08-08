package cloudtrail

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/resilience"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	"vorpalstacks/pkg/sqlparser"
)

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
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"QueryStatement is required", 400)
	}

	matches := selectPattern.FindStringSubmatch(stmt)
	if matches == nil {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"QueryStatement must contain SELECT ... FROM ...", 400)
	}

	colPart := strings.TrimSpace(matches[1])
	edsID := strings.TrimSpace(matches[2])

	// Extract EDS ID from FROM (strip ARN prefix if present).
	if idx := strings.LastIndex(edsID, "/"); idx >= 0 {
		edsID = edsID[idx+1:]
	}
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
	}

	// Parse WHERE clause using pkg/sqlparser PartiQL dialect for full
	// operator support (=, !=, >, <, >=, <=, LIKE, IN, BETWEEN, IS NULL,
	// NOT, AND, OR). Previously only `=` was supported via a single regex.
	var whereExpr sqlparser.Expr
	if wm := wherePattern.FindStringSubmatch(stmt); wm != nil {
		whereRaw := strings.TrimSpace(wm[1])
		// Wrap in a minimal SELECT so the parser has a complete statement.
		parsed, err := sqlparser.ParseWithOptions(
			"SELECT * FROM t WHERE "+whereRaw,
			sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL},
		)
		if err == nil {
			if sel, ok := parsed.(*sqlparser.Select); ok && sel.Where != nil {
				whereExpr = sel.Where.Expr
			}
		}
	}

	return &parsedQuery{
		edsID:     edsID,
		columns:   columns,
		whereExpr: whereExpr,
	}, nil
}

// executeQuery runs the parsed query against the event store and returns the
// result rows in CloudTrail Lake format ([][]map[string]string).
func (s *CloudTrailService) executeQuery(store cloudtrailstore.CloudTrailStoreInterface, pq *parsedQuery) ([][]map[string]string, error) {
	query := cloudtrailstore.NewEventQuery()
	query.MaxResults = 1000

	// Extract index-friendly conditions from the AST to pre-filter via the
	// store index where possible.  This narrows the candidate set before
	// the full per-event WHERE evaluation runs.
	if pq.whereExpr != nil {
		extractQueryConditions(pq.whereExpr, &query)
	}

	events, _, err := store.LookupEvents(query)
	if err != nil {
		return nil, err
	}
	rows := make([][]map[string]string, 0, len(events))
	for _, e := range events {
		formatted := s.formatEvent(e)

		// Evaluate the full WHERE expression against the formatted event.
		// This supports all SQL operators (=, !=, >, <, LIKE, IN, BETWEEN,
		// IS NULL, NOT, AND, OR) and all event fields.
		if pq.whereExpr != nil {
			matched, err := evaluateWhere(pq.whereExpr, formatted)
			if err != nil {
				continue
			}
			if !matched {
				continue
			}
		}

		var row []map[string]string
		for _, col := range pq.columns {
			val := ""
			if v, ok := formatted[col]; ok {
				val = fmt.Sprintf("%v", v)
			}
			row = append(row, map[string]string{col: val})
		}
		if len(pq.columns) == 1 && pq.columns[0] == "*" {
			row = row[:0]
			for k, v := range formatted {
				row = append(row, map[string]string{k: fmt.Sprintf("%v", v)})
			}
		}
		rows = append(rows, row)
	}

	return rows, nil
}

// extractEqualityConditions walks the top level of a WHERE expression and
// populates EventQuery fields for simple column = 'value' conditions. This
// enables the store's index-based pre-filter. Complex conditions (>, <, LIKE,
// IN, BETWEEN, OR) are left for per-event evaluation in executeQuery.
func extractQueryConditions(expr sqlparser.Expr, query *cloudtrailstore.EventQuery) {
	switch e := expr.(type) {
	case *sqlparser.AndExpr:
		extractQueryConditions(e.Left, query)
		extractQueryConditions(e.Right, query)
	case *sqlparser.ComparisonExpr:
		colName := strings.ToLower(getColName(e.Left))
		if colName == "" {
			return
		}
		switch e.Operator {
		case sqlparser.EqualStr:
			valStr := getExprString(e.Right)
			if valStr == "" {
				return
			}
			switch colName {
			case "eventname":
				query.EventNames = append(query.EventNames, valStr)
			case "username":
				query.Username = valStr
			case "eventsource":
				query.EventSource = valStr
			case "resourcename":
				query.ResourceNames = append(query.ResourceNames, valStr)
			case "resourcetype":
				query.ResourceType = valStr
			case "accesskeyid":
				query.AccessKeyID = valStr
			case "eventid":
				query.EventID = valStr
			case "readonly":
				query.ReadOnly = valStr
			}
		case sqlparser.GreaterEqualStr, sqlparser.GreaterThanStr:
			if colName == "eventtime" {
				if t := parseEpochFromExpr(e.Right); t != nil {
					query.StartTime = t
				}
			}
		case sqlparser.LessEqualStr, sqlparser.LessThanStr:
			if colName == "eventtime" {
				if t := parseEpochFromExpr(e.Right); t != nil {
					query.EndTime = t
				}
			}
		case sqlparser.InStr:
			if colName == "eventname" {
				tuple, ok := e.Right.(sqlparser.ValTuple)
				if !ok {
					return
				}
				for _, item := range tuple {
					if v := getExprString(item); v != "" {
						query.EventNames = append(query.EventNames, v)
					}
				}
			}
		}
	case *sqlparser.ParenExpr:
		extractQueryConditions(e.Expr, query)
	}
}

func parseEpochFromExpr(expr sqlparser.Expr) *time.Time {
	v, ok := expr.(*sqlparser.SQLVal)
	if !ok {
		return nil
	}
	epoch, err := strconv.ParseInt(string(v.Val), 10, 64)
	if err != nil || epoch <= 0 {
		return nil
	}
	t := time.Unix(epoch, 0).UTC()
	return &t
}

// evaluateWhere evaluates a WHERE expression against a formatted event row.
// Returns (matched, error). An error indicates an unsupported expression
// type, which causes the event to be skipped (fail-closed).
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

func evaluateComparison(expr *sqlparser.ComparisonExpr, row map[string]interface{}) bool {
	leftVal := getExprValue(expr.Left, row)

	switch expr.Operator {
	case sqlparser.EqualStr:
		return compareValues(leftVal, getExprValue(expr.Right, row)) == 0
	case sqlparser.NotEqualStr:
		return compareValues(leftVal, getExprValue(expr.Right, row)) != 0
	case sqlparser.LessThanStr:
		return compareValues(leftVal, getExprValue(expr.Right, row)) < 0
	case sqlparser.LessEqualStr:
		return compareValues(leftVal, getExprValue(expr.Right, row)) <= 0
	case sqlparser.GreaterThanStr:
		return compareValues(leftVal, getExprValue(expr.Right, row)) > 0
	case sqlparser.GreaterEqualStr:
		return compareValues(leftVal, getExprValue(expr.Right, row)) >= 0
	case sqlparser.LikeStr:
		return matchLike(fmt.Sprintf("%v", leftVal), fmt.Sprintf("%v", getExprValue(expr.Right, row)))
	case sqlparser.NotLikeStr:
		return !matchLike(fmt.Sprintf("%v", leftVal), fmt.Sprintf("%v", getExprValue(expr.Right, row)))
	case sqlparser.InStr:
		tuple, ok := expr.Right.(sqlparser.ValTuple)
		if !ok {
			return false
		}
		for _, item := range tuple {
			if compareValues(leftVal, getExprValue(item, row)) == 0 {
				return true
			}
		}
		return false
	case sqlparser.NotInStr:
		tuple, ok := expr.Right.(sqlparser.ValTuple)
		if !ok {
			return true
		}
		for _, item := range tuple {
			if compareValues(leftVal, getExprValue(item, row)) == 0 {
				return false
			}
		}
		return true
	}
	return false
}

func evaluateRangeCond(expr *sqlparser.RangeCond, row map[string]interface{}) bool {
	val := getExprValue(expr.Left, row)
	fromVal := getExprValue(expr.From, row)
	toVal := getExprValue(expr.To, row)

	inRange := compareValues(val, fromVal) >= 0 && compareValues(val, toVal) <= 0

	switch expr.Operator {
	case sqlparser.BetweenStr:
		return inRange
	case sqlparser.NotBetweenStr:
		return !inRange
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
		colName := e.Name.String()
		if !e.Qualifier.IsEmpty() {
			qualifiedKey := e.Qualifier.Name.String() + "." + colName
			if val, exists := row[qualifiedKey]; exists {
				return val
			}
		}
		return row[colName]
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
	case *sqlparser.NullVal:
		return nil
	}
	return nil
}

func getColName(expr sqlparser.Expr) string {
	if cn, ok := expr.(*sqlparser.ColName); ok {
		return cn.Name.String()
	}
	return ""
}

func getExprString(expr sqlparser.Expr) string {
	if v, ok := expr.(*sqlparser.SQLVal); ok {
		return string(v.Val)
	}
	return ""
}

func compareValues(left, right interface{}) int {
	leftFloat, leftErr := toFloat(left)
	rightFloat, rightErr := toFloat(right)

	if leftErr == nil && rightErr == nil {
		if leftFloat < rightFloat {
			return -1
		} else if leftFloat > rightFloat {
			return 1
		}
		return 0
	}

	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)
	if leftStr < rightStr {
		return -1
	} else if leftStr > rightStr {
		return 1
	}
	return 0
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

// StartQuery starts a CloudTrail Lake query.
func (s *CloudTrailService) StartQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stmt := request.GetStringParam(req.Parameters, "QueryStatement")
	if stmt == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"QueryStatement is required", 400)
	}

	pq, err := parseQueryStatement(stmt)
	if err != nil {
		return nil, err
	}

	// Verify the EDS exists.
	eds, err := store.GetEventDataStore(pq.edsID)
	if err != nil {
		return nil, awserrors.NewAWSError("EventDataStoreNotFoundException",
			"Event data store not found", 404)
	}
	if eds.Status == "PENDING_DELETION" {
		return nil, awserrors.NewAWSError("OperationNotPermittedException",
			"Cannot query a PENDING_DELETION event data store", 400)
	}

	queryID := uuid.NewString()
	now := time.Now().UTC()

	// Save the query record with RUNNING status before returning.
	qr := &cloudtrailstore.QueryRecord{
		QueryID:        queryID,
		EventDataStore: eds.EventDataStoreID,
		QueryStatement: stmt,
		QueryStatus:    "RUNNING",
		StartTime:      now,
	}

	if err := store.SaveQuery(qr); err != nil {
		return nil, s.mapStoreError(err)
	}

	// Execute the query asynchronously.
	go func() {
		defer func() {
			if r := resilience.RecoverPanic("cloudtrail.StartQuery"); r != nil {
				endTime := time.Now().UTC()
				qr.EndTime = &endTime
				qr.QueryStatus = "FAILED"
				qr.ErrorMessage = fmt.Sprintf("internal error: panic recovered: %v", r)
				_ = store.SaveQuery(qr)
			}
		}()

		results, execErr := s.executeQuery(store, pq)
		endTime := time.Now().UTC()
		qr.EndTime = &endTime
		if execErr != nil {
			qr.QueryStatus = "FAILED"
			qr.ErrorMessage = execErr.Error()
		} else {
			qr.QueryStatus = "FINISHED"
			qr.QueryResultRows = results
			qr.ResultsCount = int32(len(results))
			qr.BytesScanned = int64(len(results) * 500)
		}
		if err := store.SaveQuery(qr); err != nil {
			slog.Error("Failed to save query result", "queryId", queryID, "error", err)
		}
	}()

	return map[string]interface{}{
		"QueryId": queryID,
	}, nil
}

// GetQueryResults retrieves the results of a CloudTrail Lake query.
func (s *CloudTrailService) GetQueryResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	queryID := request.GetStringParam(req.Parameters, "QueryId")
	if queryID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"QueryId is required", 400)
	}

	qr, err := store.GetQuery(queryID)
	if err != nil {
		return nil, awserrors.NewAWSError("QueryIdNotFoundException",
			"Query not found", 404)
	}

	maxResults := request.GetIntParam(req.Parameters, "MaxQueryResults")
	if maxResults <= 0 {
		maxResults = 50
	}

	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	offset := 0
	if nextToken != "" {
		if n, err := strconv.Atoi(nextToken); err == nil {
			offset = n
		}
	}

	end := offset + maxResults
	if end > len(qr.QueryResultRows) {
		end = len(qr.QueryResultRows)
	}

	// Initialise with make to ensure JSON serialises as [] not null.
	rows := make([]interface{}, 0)
	if offset < len(qr.QueryResultRows) {
		for i := offset; i < end; i++ {
			rows = append(rows, qr.QueryResultRows[i])
		}
	}

	result := map[string]interface{}{
		"QueryId":         qr.QueryID,
		"QueryStatus":     qr.QueryStatus,
		"QueryResultRows": rows,
		"QueryStatistics": map[string]interface{}{
			"ResultsCount":      qr.ResultsCount,
			"TotalResultsCount": qr.ResultsCount,
			"BytesScanned":      qr.BytesScanned,
		},
	}

	if end < len(qr.QueryResultRows) {
		result["NextToken"] = strconv.Itoa(end)
	}

	if qr.ErrorMessage != "" {
		result["ErrorMessage"] = qr.ErrorMessage
	}

	return result, nil
}

// DescribeQuery retrieves metadata about a CloudTrail Lake query.
func (s *CloudTrailService) DescribeQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	queryID := request.GetStringParam(req.Parameters, "QueryId")
	if queryID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"QueryId is required", 400)
	}

	qr, err := store.GetQuery(queryID)
	if err != nil {
		return nil, awserrors.NewAWSError("QueryIdNotFoundException",
			"Query not found", 404)
	}

	result := map[string]interface{}{
		"QueryId":     qr.QueryID,
		"QueryStatus": qr.QueryStatus,
		"QueryString": qr.QueryStatement,
		"QueryStatistics": map[string]interface{}{
			"ResultsCount":      qr.ResultsCount,
			"TotalResultsCount": qr.ResultsCount,
			"BytesScanned":      qr.BytesScanned,
		},
	}

	if qr.ErrorMessage != "" {
		result["ErrorMessage"] = qr.ErrorMessage
	}

	return result, nil
}

// CancelQuery cancels a running CloudTrail Lake query.
func (s *CloudTrailService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	queryID := request.GetStringParam(req.Parameters, "QueryId")
	if queryID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"QueryId is required", 400)
	}

	qr, err := store.GetQuery(queryID)
	if err != nil {
		return nil, awserrors.NewAWSError("QueryIdNotFoundException",
			"Query not found", 404)
	}

	if qr.QueryStatus == "FINISHED" || qr.QueryStatus == "FAILED" ||
		qr.QueryStatus == "CANCELLED" || qr.QueryStatus == "TIMED_OUT" {
		return nil, awserrors.NewAWSError("OperationNotPermittedException",
			"Cannot cancel a query that has already finished, been cancelled, or timed out", 400)
	}

	qr.QueryStatus = "CANCELLED"
	now := time.Now().UTC()
	qr.EndTime = &now

	if err := store.SaveQuery(qr); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"QueryId":     qr.QueryID,
		"QueryStatus": qr.QueryStatus,
	}, nil
}

// ListQueries lists queries for an event data store.
func (s *CloudTrailService) ListQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	edsID := request.GetStringParam(req.Parameters, "EventDataStore")
	if edsID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"EventDataStore is required", 400)
	}

	if idx := strings.LastIndex(edsID, "/"); idx >= 0 {
		edsID = edsID[idx+1:]
	}

	queries, err := store.ListQueriesByEDS(edsID)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// Optional filters.
	statusFilter := request.GetStringParam(req.Parameters, "QueryStatus")
	if statusFilter != "" {
		if err := validateQueryStatus(statusFilter); err != nil {
			return nil, err
		}
	}
	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 {
		maxResults = 50
	}
	offset := 0
	if nt := request.GetStringParam(req.Parameters, "NextToken"); nt != "" {
		if n, err := strconv.Atoi(nt); err == nil && n > 0 {
			offset = n
		}
	}

	// Filter queries by status.
	var filtered []*cloudtrailstore.QueryRecord
	for _, qr := range queries {
		if statusFilter != "" && qr.QueryStatus != statusFilter {
			continue
		}
		filtered = append(filtered, qr)
	}

	// Paginate the filtered results.
	queryList := make([]map[string]interface{}, 0)
	end := offset + maxResults
	if end > len(filtered) {
		end = len(filtered)
	}
	for i := offset; i < end; i++ {
		qr := filtered[i]
		queryList = append(queryList, map[string]interface{}{
			"QueryId":        qr.QueryID,
			"QueryStatus":    qr.QueryStatus,
			"StartTime":      qr.StartTime.Unix(),
			"EventDataStore": qr.EventDataStore,
		})
	}

	result := map[string]interface{}{
		"Queries": queryList,
	}
	if end < len(filtered) {
		result["NextToken"] = strconv.Itoa(end)
	}

	return result, nil
}
