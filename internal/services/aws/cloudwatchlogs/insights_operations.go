package cloudwatchlogs

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
)

type queryState struct {
	queryId             string
	logGroupNames       []string
	logGroupIdentifiers []string
	startTime           int64
	endTime             int64
	queryString         string
	queryLanguage       string
	status              string
	results             []queryResultRow
	stats               queryStats
	createdAt           time.Time
}

// StartQuery initiates a CloudWatch Logs Insights query.
func (s *LogsService) StartQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	startTime := int64(request.GetIntParam(req.Parameters, "StartTime"))
	endTime := int64(request.GetIntParam(req.Parameters, "EndTime"))
	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")

	if queryString == "" || startTime <= 0 || endTime <= 0 {
		return nil, ErrMissingParameter
	}

	var logGroupNames []string
	if name := request.GetParamLowerFirst(req.Parameters, "LogGroupName"); name != "" {
		logGroupNames = append(logGroupNames, name)
	}
	if names, ok := req.Parameters["logGroupNames"]; ok {
		if arr, ok := names.([]interface{}); ok {
			for _, item := range arr {
				if n, ok := item.(string); ok {
					logGroupNames = append(logGroupNames, n)
				}
			}
		}
	}

	var logGroupIdentifiers []string
	if idents, ok := req.Parameters["logGroupIdentifiers"]; ok {
		if arr, ok := idents.([]interface{}); ok {
			for _, item := range arr {
				if id, ok := item.(string); ok {
					logGroupIdentifiers = append(logGroupIdentifiers, id)
				}
			}
		}
	}

	limit := int64(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 10000
	}

	if len(logGroupNames) == 0 && len(logGroupIdentifiers) == 0 {
		return nil, ErrMissingParameter
	}

	allGroups := logGroupNames
	allGroups = append(allGroups, logGroupIdentifiers...)

	queryId := fmt.Sprintf("query-%d", time.Now().UnixNano())
	region := reqCtx.GetRegion()

	qs := &queryState{
		queryId:             queryId,
		logGroupNames:       allGroups,
		logGroupIdentifiers: logGroupIdentifiers,
		startTime:           startTime,
		endTime:             endTime,
		queryString:         queryString,
		queryLanguage:       queryLanguage,
		status:              "Running",
		createdAt:           time.Now(),
	}
	s.queries.Store(queryId, qs)

	go s.executeQuery(region, queryId, queryString, allGroups, startTime, endTime, limit)

	return map[string]interface{}{
		"queryId": queryId,
	}, nil
}

func (s *LogsService) executeQuery(region, queryId, queryString string, logGroupNames []string, startTime, endTime, limit int64) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		s.failQuery(queryId, fmt.Sprintf("store error: %v", err))
		return
	}

	var allEvents []logEventWithContext
	for _, lgName := range logGroupNames {
		streams, _, _ := store.ListLogStreams(lgName, "", "", 1000)
		for _, ls := range streams {
			events, _, _, _ := store.GetLogEvents(lgName, ls.Name, startTime, endTime, int(limit), true, "")
			for _, evt := range events {
				allEvents = append(allEvents, logEventWithContext{
					timestamp:     evt.Timestamp,
					message:       evt.Message,
					ingestionTime: evt.IngestionTime,
					logGroup:      lgName,
					logStream:     ls.Name,
				})
			}
		}
	}

	rows, stats := executeQuery(queryString, allEvents)

	val, ok := s.queries.Load(queryId)
	if !ok {
		return
	}
	qs := val.(*queryState)
	qs.results = rows
	qs.stats = stats
	qs.status = "Complete"
}

func (s *LogsService) failQuery(queryId, message string) {
	val, ok := s.queries.Load(queryId)
	if !ok {
		return
	}
	qs := val.(*queryState)
	qs.status = "Failed"
	_ = message
}

// StopQuery stops a running query.
func (s *LogsService) StopQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryId := request.GetParamLowerFirst(req.Parameters, "QueryId")
	if queryId == "" {
		return nil, ErrMissingParameter
	}

	val, ok := s.queries.Load(queryId)
	if !ok {
		return nil, NewLogsError("ResourceNotFoundException",
			fmt.Sprintf("Query %s not found", queryId), 404)
	}

	qs := val.(*queryState)
	qs.status = "Cancelled"

	return map[string]interface{}{
		"success": true,
	}, nil
}

// DescribeQueries lists queries.
func (s *LogsService) DescribeQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statusFilter := request.GetParamLowerFirst(req.Parameters, "Status")
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))
	if maxResults <= 0 {
		maxResults = 50
	}

	var allQueries []*queryState
	s.queries.Range(func(key, value interface{}) bool {
		qs := value.(*queryState)
		if statusFilter != "" && qs.status != statusFilter {
			return true
		}
		if logGroupName != "" {
			found := false
			for _, n := range qs.logGroupNames {
				if n == logGroupName {
					found = true
					break
				}
			}
			if !found {
				return true
			}
		}
		allQueries = append(allQueries, qs)
		return true
	})

	result := pagination.PaginateSlice(allQueries, nextToken, int(maxResults), func(qs *queryState) string {
		return qs.queryId
	})

	queries := make([]map[string]interface{}, len(result.Items))
	for i, qs := range result.Items {
		queries[i] = formatQueryInfo(qs)
	}

	resp := map[string]interface{}{
		"queries": queries,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}

	return resp, nil
}

// GetQueryResults retrieves the results of a completed query.
func (s *LogsService) GetQueryResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryId := request.GetParamLowerFirst(req.Parameters, "QueryId")
	if queryId == "" {
		return nil, ErrMissingParameter
	}

	val, ok := s.queries.Load(queryId)
	if !ok {
		return nil, NewLogsError("ResourceNotFoundException",
			fmt.Sprintf("Query %s not found", queryId), 404)
	}

	qs := val.(*queryState)

	limit := int(request.GetIntParam(req.Parameters, "MaxItems"))
	if limit <= 0 {
		limit = 10000
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	offset := 0
	if nextToken != "" {
		if n, err := parseInt(nextToken); err == nil {
			offset = n
		}
	}

	endIdx := offset + limit
	if endIdx > len(qs.results) {
		endIdx = len(qs.results)
	}
	if offset > len(qs.results) {
		offset = len(qs.results)
	}

	pageResults := qs.results[offset:endIdx]

	resultRows := make([][]map[string]interface{}, len(pageResults))
	for i, row := range pageResults {
		fields := make([]map[string]interface{}, 0, len(row.fields))
		for k, v := range row.fields {
			fields = append(fields, map[string]interface{}{
				"field": k,
				"value": v,
			})
		}
		resultRows[i] = fields
	}

	statsMap := map[string]interface{}{
		"recordsMatched":          qs.stats.recordsMatched,
		"recordsScanned":          qs.stats.recordsScanned,
		"bytesScanned":            qs.stats.bytesScanned,
		"estimatedRecordsSkipped": 0,
		"estimatedBytesSkipped":   0,
		"logGroupsScanned":        0,
	}

	resp := map[string]interface{}{
		"queryId":    queryId,
		"status":     qs.status,
		"results":    resultRows,
		"statistics": statsMap,
	}
	if qs.queryLanguage != "" {
		resp["queryLanguage"] = qs.queryLanguage
	}

	if endIdx < len(qs.results) {
		resp["nextToken"] = fmt.Sprintf("%d", endIdx)
	}

	return resp, nil
}

func formatQueryInfo(qs *queryState) map[string]interface{} {
	logGroupName := ""
	if len(qs.logGroupNames) > 0 {
		logGroupName = qs.logGroupNames[0]
	}
	result := map[string]interface{}{
		"queryId":      qs.queryId,
		"queryString":  qs.queryString,
		"status":       qs.status,
		"createTime":   qs.createdAt.UnixMilli(),
		"logGroupName": logGroupName,
	}
	if qs.queryLanguage != "" {
		result["queryLanguage"] = qs.queryLanguage
	}
	return result
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
