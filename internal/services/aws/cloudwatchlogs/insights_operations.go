package cloudwatchlogs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	corelogs "vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
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
	errorMessage        string
	results             []queryResultRow
	stats               queryStats
	createdAt           time.Time
}

// StartQuery initiates a CloudWatch Logs Insights query.
func (s *LogsService) StartQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")
	startTime := int64(request.GetIntParam(req.Parameters, "StartTime"))
	endTime := int64(request.GetIntParam(req.Parameters, "EndTime"))

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

	limit32, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 10000, 100000)
	if err != nil {
		return nil, err
	}

	queryId, err := s.startQueryCore(&StartQueryInput{
		StartTime:           startTime,
		EndTime:             endTime,
		QueryString:         queryString,
		QueryLanguage:       queryLanguage,
		LogGroupNames:       logGroupNames,
		LogGroupIdentifiers: logGroupIdentifiers,
		Limit:               int64(limit32),
		Region:              reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"queryId": queryId,
	}, nil
}

func (s *LogsService) executeQuery(region, queryId, queryString string, logGroupNames []string, startTime, endTime, limit int64) {
	defer func() {
		if r := recover(); r != nil {
			s.failQuery(queryId, fmt.Sprintf("panic: %v", r))
		}
	}()

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		s.failQuery(queryId, fmt.Sprintf("store error: %v", err))
		return
	}

	ctx := &execContext{
		startTime:     startTime,
		endTime:       endTime,
		accountID:     s.accountID,
		defaultGroups: logGroupNames,
		events:        fetchLogEvents(store, logGroupNames, startTime, endTime),
		fetchEvents: func(groups []string, start, end int64) ([]logEventWithContext, error) {
			return fetchLogEvents(store, groups, start, end), nil
		},
		listLogGroups: func() ([]sourceGroupInfo, error) {
			return listSourceGroups(store), nil
		},
		getLookupTable: func(name string) (*parsedLookupTable, error) {
			lt, err := store.GetLookupTable(name)
			if err != nil {
				return nil, fmt.Errorf("lookup table %s not found", name)
			}
			body, err := s.lookupTablePlainBody(lt, region)
			if err != nil {
				return nil, fmt.Errorf("lookup table %s is unavailable: %v", name, err)
			}
			columns, records, err := parseLookupCSV(body)
			if err != nil {
				return nil, fmt.Errorf("lookup table %s is invalid: %v", name, err)
			}
			return newParsedLookupTable(columns, records), nil
		},
		subqueryCache: map[string][]interface{}{},
	}

	rows, err := executeQueryContext(ctx, queryString)
	if err != nil {
		s.failQuery(queryId, err.Error())
		return
	}
	if ctx.sourceError != nil {
		s.failQuery(queryId, fmt.Sprintf("source error: %v", ctx.sourceError))
		return
	}
	if int64(len(rows)) > limit {
		rows = rows[:limit]
	}

	stats := queryStats{
		recordsScanned: int64(len(ctx.events)),
	}
	for _, e := range ctx.events {
		stats.bytesScanned += int64(len(e.message))
	}
	stats.recordsMatched = int64(len(rows))

	val, ok := s.queries.Load(queryId)
	if !ok {
		return
	}
	qs := val.(*queryState)
	qs.results = rows
	qs.stats = stats
	qs.status = "Complete"
}

// fetchLogEvents reads all events of the given log groups within the time
// window. Stream listing and event reads surface their errors through the
// server log; unresolvable groups yield no events, matching the documented
// behaviour of querying a group without matching events.
func fetchLogEvents(store *logsstore.Store, groups []string, startTime, endTime int64) []logEventWithContext {
	var allEvents []logEventWithContext
	for _, lgName := range groups {
		streams, _, err := store.ListLogStreams(lgName, "", "", 1000)
		if err != nil {
			corelogs.Error("Failed to list log streams for query",
				corelogs.String("logGroup", lgName), corelogs.Err(err))
			continue
		}
		for _, ls := range streams {
			events, _, _, err := store.GetLogEvents(lgName, ls.Name, startTime, endTime, 10000, true, "")
			if err != nil {
				corelogs.Error("Failed to read log events for query",
					corelogs.String("logGroup", lgName), corelogs.String("logStream", ls.Name), corelogs.Err(err))
				continue
			}
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
	return allEvents
}

// listSourceGroups returns the log group inventory for SOURCE selection.
func listSourceGroups(store *logsstore.Store) []sourceGroupInfo {
	var out []sourceGroupInfo
	marker := ""
	for {
		groups, next, err := store.ListLogGroups("", marker, 50)
		if err != nil {
			return out
		}
		for _, g := range groups {
			out = append(out, sourceGroupInfo{Name: g.Name, Class: g.LogGroupClass, Tags: g.Tags})
		}
		if next == "" || len(groups) == 0 {
			break
		}
		marker = next
	}
	return out
}

func (s *LogsService) failQuery(queryId, message string) {
	val, ok := s.queries.Load(queryId)
	if !ok {
		return
	}
	qs := val.(*queryState)
	qs.status = "Failed"
	qs.errorMessage = message
}

// StopQuery stops a running query.
func (s *LogsService) StopQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryId := request.GetParamLowerFirst(req.Parameters, "QueryId")
	if err := s.stopQueryCore(queryId); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

// DescribeQueries lists queries.
func (s *LogsService) DescribeQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "MaxResults")), 50, 1000)
	if err != nil {
		return nil, err
	}

	items, nextToken, err := s.describeQueriesCore(&DescribeQueriesInput{
		StatusFilter: request.GetParamLowerFirst(req.Parameters, "Status"),
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		NextToken:    request.GetParamLowerFirst(req.Parameters, "NextToken"),
		MaxResults:   maxResults,
	})
	if err != nil {
		return nil, err
	}

	queries := make([]map[string]interface{}, len(items))
	for i, qs := range items {
		queries[i] = formatQueryInfo(qs)
	}

	resp := map[string]interface{}{
		"queries": queries,
	}
	if nextToken != "" {
		resp["nextToken"] = nextToken
	}

	return resp, nil
}

// resultTimestampLayout is the rendering query results use for
// timestamp-typed values.
const resultTimestampLayout = "2006-01-02 15:04:05.000"

// formatResultTimestamp renders the internal epoch-millisecond value of
// @timestamp and @ingestionTime in the "2006-01-02 15:04:05.000" form that
// query results present. Non-numeric values pass through unchanged.
func formatResultTimestamp(v string) string {
	ms, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return v
	}
	return time.UnixMilli(ms).UTC().Format(resultTimestampLayout)
}

// parseResultTimestamp parses the "2006-01-02 15:04:05.000" rendering back
// into epoch milliseconds, so stored timestamp values round-trip through
// later commands such as dateceil over a binned column.
func parseResultTimestamp(s string) (int64, bool) {
	t, err := time.Parse(resultTimestampLayout, strings.TrimSpace(s))
	if err != nil {
		return 0, false
	}
	return t.UnixMilli(), true
}

// GetQueryResults retrieves the results of a completed query.
func (s *LogsService) GetQueryResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit32, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "MaxItems")), 10000, 10000)
	if err != nil {
		return nil, err
	}

	result, err := s.getQueryResultsCore(&GetQueryResultsInput{
		QueryId:   request.GetParamLowerFirst(req.Parameters, "QueryId"),
		NextToken: request.GetParamLowerFirst(req.Parameters, "NextToken"),
		MaxItems:  limit32,
	})
	if err != nil {
		return nil, err
	}

	resultRows := make([][]map[string]interface{}, len(result.Results))
	for i := range result.Results {
		row := &result.Results[i]
		fields := make([]map[string]interface{}, 0, len(row.fields))
		for _, k := range row.ordered() {
			v := row.fields[k]
			if k == "@timestamp" || k == "@ingestionTime" {
				v = formatResultTimestamp(v)
			}
			fields = append(fields, map[string]interface{}{
				"field": k,
				"value": v,
			})
		}
		resultRows[i] = fields
	}

	statsMap := map[string]interface{}{
		"recordsMatched":          result.Stats.recordsMatched,
		"recordsScanned":          result.Stats.recordsScanned,
		"bytesScanned":            result.Stats.bytesScanned,
		"estimatedRecordsSkipped": 0,
		"estimatedBytesSkipped":   0,
		"logGroupsScanned":        0,
	}

	resp := map[string]interface{}{
		"queryId":    result.QueryId,
		"status":     result.Status,
		"results":    resultRows,
		"statistics": statsMap,
	}
	if result.QueryLanguage != "" {
		resp["queryLanguage"] = result.QueryLanguage
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
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
	if qs.errorMessage != "" {
		result["errorMessage"] = qs.errorMessage
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
