package cloudtrail

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// queryFieldsPattern extracts the SELECT and FROM portions of a CloudTrail
// Lake SQL-like QueryStatement. Example: "SELECT eventID, eventTime FROM
// <eds-id> WHERE eventName = 'PutObject'".
var (
	selectPattern = regexp.MustCompile(`(?i)^SELECT\s+(.+?)\s+FROM\s+([^\s]+)`)
	wherePattern  = regexp.MustCompile(`(?i)WHERE\s+(.+)$`)
)

// parseQueryStatement extracts the EDS ID, selected columns, and WHERE
// conditions from a CloudTrail Lake QueryStatement.
type parsedQuery struct {
	edsID    string
	columns  []string
	whereRaw string
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

	whereRaw := ""
	if wm := wherePattern.FindStringSubmatch(stmt); wm != nil {
		whereRaw = strings.TrimSpace(wm[1])
	}

	return &parsedQuery{
		edsID:    edsID,
		columns:  columns,
		whereRaw: whereRaw,
	}, nil
}

// executeQuery runs the parsed query against the event store and returns the
// result rows in CloudTrail Lake format ([][]map[string]string).
func (s *CloudTrailService) executeQuery(store cloudtrailstore.CloudTrailStoreInterface, pq *parsedQuery) ([][]map[string]string, error) {
	query := cloudtrailstore.NewEventQuery()
	query.MaxResults = 1000

	if pq.whereRaw != "" {
		parseWhereClause(pq.whereRaw, &query)
	}

	events, _, err := store.LookupEvents(query)
	if err != nil {
		return nil, err
	}

	rows := make([][]map[string]string, 0, len(events))
	for _, e := range events {
		formatted := s.formatEvent(e)
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

// parseWhereClause extracts simple WHERE conditions from the query string.
// Supports: eventName = 'x', eventSource = 'x', username = 'x',
// accessKeyId = 'x', readOnly = 'true'.
var (
	whereEqualPattern = regexp.MustCompile(`(?i)(\w+)\s*=\s*'([^']*)'`)
)

func parseWhereClause(whereRaw string, query *cloudtrailstore.EventQuery) {
	matches := whereEqualPattern.FindAllStringSubmatch(whereRaw, -1)
	for _, m := range matches {
		field := strings.ToLower(m[1])
		value := m[2]
		switch field {
		case "eventname":
			query.EventNames = append(query.EventNames, value)
		case "eventsource":
			query.EventSource = value
		case "username":
			query.Username = value
		case "accesskeyid":
			query.AccessKeyID = value
		case "readonly":
			query.ReadOnly = value
		case "eventid":
			query.EventID = value
		}
	}

	// Extract time range conditions (eventTime > epoch_seconds).
	if m := regexp.MustCompile(`(?i)eventTime\s*>=?\s*(\d+)`).FindStringSubmatch(whereRaw); m != nil {
		if epoch, err := strconv.ParseInt(m[1], 10, 64); err == nil && epoch > 0 {
			t := time.Unix(epoch, 0).UTC()
			query.StartTime = &t
		}
	}
	if m := regexp.MustCompile(`(?i)eventTime\s*<=?\s*(\d+)`).FindStringSubmatch(whereRaw); m != nil {
		if epoch, err := strconv.ParseInt(m[1], 10, 64); err == nil && epoch > 0 {
			t := time.Unix(epoch, 0).UTC()
			query.EndTime = &t
		}
	}
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
		_ = store.SaveQuery(qr)
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
			// qr.QueryResultRows[i] is []map[string]string, already in
			// the correct SDK format ([][]map[string]string).
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

// CancelQuery cancels a running or finished CloudTrail Lake query.
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

	if qr.QueryStatus == "FINISHED" || qr.QueryStatus == "FAILED" || qr.QueryStatus == "CANCELLED" {
		return nil, awserrors.NewAWSError("OperationNotPermittedException",
			"Cannot cancel a query that has already finished or been cancelled", 400)
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
