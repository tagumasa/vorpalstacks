package cloudwatchlogs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// queryState is the lifecycle record of one query. Every field is guarded
// by mu: the executing worker writes status/results/stats while the read
// plane (GetQueryResults, DescribeQueries, StopQuery, lookup-table
// resolution) inspects the same state concurrently. cancelled is the
// cancellation flag StopQuery raises; the worker honours it at its stage
// checkpoints and never overwrites a terminal status. region names the
// store the state persists to, so every transition writes through to the
// record that survives a restart.
type queryState struct {
	mu                  sync.RWMutex
	queryId             string
	logGroupNames       []string
	logGroupIdentifiers []string
	// scannedGroups is the group set the query actually analyses — the
	// explicit list at start, replaced by the SOURCE-resolved set once
	// execution fixes it. ListLogGroupsForQuery serves it.
	scannedGroups []string
	startTime     int64
	endTime       int64
	queryString   string
	queryLanguage string
	status        string
	cancelled     bool
	results       []queryResultRow
	stats         queryStats
	createdAt     time.Time
	region        string
	// userIdentity is the ARN of the principal that started the query —
	// the member DescribeQueries reports back; empty when the caller is
	// anonymous. durationMs is the wall-clock duration stamped at the
	// terminal transition; zero while the query runs.
	userIdentity string
	durationMs   int64
}

// markTerminal stamps the terminal status together with the duration the
// query ran for, the value DescribeQueries keeps reporting after the
// execution ends. Callers hold the write lock.
func (qs *queryState) markTerminal(status string) {
	qs.status = status
	qs.durationMs = time.Since(qs.createdAt).Milliseconds()
}

// queryStatusRunning and its siblings are the QueryStatus vocabulary the
// Smithy model enumerates for the status members. The platform produces
// Running, Complete, Failed, Cancelled and Timeout; Scheduled is the
// status of an in-progress scheduled-query execution (surfaced through
// DescribeQueries), while Unknown stays unproduced — the platform has no
// indeterminate states. Timeout is the sixty-minute runtime bound
// ("Queries time out after 60 minutes of runtime").
const (
	queryStatusRunning   = "Running"
	queryStatusComplete  = "Complete"
	queryStatusFailed    = "Failed"
	queryStatusCancelled = "Cancelled"
	queryStatusScheduled = "Scheduled"
	queryStatusTimeout   = "Timeout"
)

// queryRuntimeLimit is the documented outer bound: "Queries time out
// after 60 minutes of runtime."
const queryRuntimeLimit = 60 * time.Minute

// isTerminalQueryStatus reports whether the status admits no further
// transition. A terminal status always wins: neither the worker's
// finalisation nor a failure path overwrites it.
func isTerminalQueryStatus(status string) bool {
	return status == queryStatusComplete || status == queryStatusFailed ||
		status == queryStatusCancelled || status == queryStatusTimeout
}

// StartQuery initiates a CloudWatch Logs Insights query.
func (s *LogsService) StartQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")
	startTime := int64(request.GetIntParam(req.Parameters, "StartTime"))
	endTime := int64(request.GetIntParam(req.Parameters, "EndTime"))

	// The three required members' zero values are inside their
	// documented domains, so the presence flags — derived over both wire
	// key casings — carry requiredness and the values stay as read.
	present := func(keys ...string) bool {
		for _, key := range keys {
			if _, ok := req.Parameters[key]; ok {
				return true
			}
		}
		return false
	}
	queryStringSet := present("QueryString", "queryString")
	startTimeSet := present("StartTime", "startTime")
	endTimeSet := present("EndTime", "endTime")

	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	logGroupNames := request.GetStringList(req.Parameters, "LogGroupNames")
	logGroupIdentifiers := request.GetStringList(req.Parameters, "LogGroupIdentifiers")

	// The query records the ARN of the principal that started it — the
	// userIdentity DescribeQueries reports back. Anonymous requests carry
	// no identity and the member stays absent.
	caller := ""
	builder := svcarn.NewARNBuilder(reqCtx.AccountID, reqCtx.GetRegion())
	switch reqCtx.PrincipalType {
	case request.PrincipalTypeUser:
		if reqCtx.Principal != "" {
			caller = builder.IAM().User(reqCtx.Principal)
		}
	case request.PrincipalTypeRole:
		if reqCtx.Principal != "" {
			caller = builder.IAM().Role(reqCtx.Principal)
		}
	case request.PrincipalTypeRoot:
		caller = builder.IAM().Root()
	}

	queryId, err := s.startQueryCore(&StartQueryInput{
		StartTime:           startTime,
		EndTime:             endTime,
		QueryString:         queryString,
		QueryLanguage:       queryLanguage,
		StartTimeSet:        startTimeSet,
		EndTimeSet:          endTimeSet,
		QueryStringSet:      queryStringSet,
		LogGroupName:        logGroupName,
		LogGroupNames:       logGroupNames,
		LogGroupIdentifiers: logGroupIdentifiers,
		Limit:               int64(request.GetIntParam(req.Parameters, "Limit")),
		Region:              reqCtx.GetRegion(),
		UserIdentity:        caller,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"queryId": queryId,
	}, nil
}

func (s *LogsService) executeQuery(ctx context.Context, region, queryId, queryString string, logGroupNames []string, startTime, endTime, limit int64) {
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
	if s.queryCancelled(ctx, queryId) {
		return
	}

	execCtx := &execContext{
		startTime:     startTime,
		endTime:       endTime,
		accountID:     s.accountID,
		defaultGroups: logGroupNames,
		events:        fetchGroupEventsForQuery(store, logGroupNames, startTime, endTime),
		fetchEvents: func(groups []string, start, end int64) ([]logEventWithContext, error) {
			return fetchGroupEventsForQuery(store, groups, start, end), nil
		},
		listLogGroups: func() ([]sourceGroupInfo, error) {
			return listSourceGroups(store), nil
		},
		getLookupTable: func(name string) (*parsedLookupTable, error) {
			return s.loadParsedLookupTable(store, region, name)
		},
		subqueryCache: map[string][]interface{}{},
	}
	execCtx.startedMs = time.Now().UnixMilli()
	execCtx.deadline = queryDeadlineNow().Add(queryRuntimeLimit)

	// The error path is the failure path — compilation, validation and
	// source errors all return through err — while a query that breached
	// the sixty-minute bound returns without one and stamps the Timeout
	// status, a sibling of Failed rather than a failure identity.
	rows, err := executeQueryContext(execCtx, queryString)
	if err != nil {
		s.failQuery(queryId, err.Error())
		return
	}
	if execCtx.timedOut {
		s.timeoutQuery(queryId)
		return
	}
	// Cancellation checkpoint between the execution and result commit: a
	// StopQuery that landed mid-run has already set Cancelled, and the
	// terminal status wins over the completion write.
	if s.queryCancelled(ctx, queryId) {
		return
	}
	if int64(len(rows)) > limit {
		rows = rows[:limit]
	}

	stats := queryStats{
		recordsScanned: int64(len(execCtx.events)),
	}
	for _, e := range execCtx.events {
		stats.bytesScanned += int64(len(e.message))
	}
	stats.recordsMatched = execCtx.recordsMatched

	val, ok := s.queries.Load(queryId)
	if !ok {
		return
	}
	qs := val.(*queryState)
	qs.mu.Lock()
	if isTerminalQueryStatus(qs.status) {
		qs.mu.Unlock()
		return
	}
	// A SOURCE command replaced the default group set during execution;
	// the scanned set the listing serves is the resolved one.
	if execCtx.effectiveGroups != nil {
		qs.scannedGroups = execCtx.effectiveGroups
	}
	qs.results = rows
	qs.stats = stats
	qs.markTerminal(queryStatusComplete)
	qs.mu.Unlock()
	// The completed state writes through so GetQueryResults and
	// DescribeQueries survive a restart; the map stays the live registry
	// and a persistence failure degrades to the in-memory result, logged
	// with the query id.
	if err := s.persistQueryState(qs); err != nil {
		logs.Error("Failed to persist completed query state",
			logs.String("queryId", queryId), logs.Err(err))
	}
}

// queryCancelled reports whether the query should stop before its next
// stage: the service context is done (shutdown drains the worker without a
// spurious Failed status) or StopQuery raised the cancellation flag.
func (s *LogsService) queryCancelled(ctx context.Context, queryId string) bool {
	if ctx != nil && ctx.Err() != nil {
		return true
	}
	val, ok := s.queries.Load(queryId)
	if !ok {
		return true
	}
	qs := val.(*queryState)
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	return qs.cancelled || isTerminalQueryStatus(qs.status)
}

// listSourceGroups returns the log group inventory for SOURCE selection:
// the shared fetch-all walk plus the per-group tag enrichment.
func listSourceGroups(store *logsstore.Store) []sourceGroupInfo {
	groups, err := fetchAllLogGroups(store)
	if err != nil {
		return nil
	}
	out := make([]sourceGroupInfo, 0, len(groups))
	for _, g := range groups {
		// SOURCE tag filters read the LIVE tags from the tag store —
		// the stored record's copy is write-once and goes stale the
		// moment a tag mutation lands.
		tags, _ := store.Tags().List(g.ARN)
		out = append(out, sourceGroupInfo{Name: g.Name, Class: g.LogGroupClass, Tags: tags})
	}
	return out
}

// timeoutQuery stamps the sixty-minute runtime breach: Timeout is the
// status the vocabulary carries for it, a sibling of Failed rather than
// a failure identity.
func (s *LogsService) timeoutQuery(queryId string) {
	val, ok := s.queries.Load(queryId)
	if !ok {
		return
	}
	qs := val.(*queryState)
	qs.mu.Lock()
	if isTerminalQueryStatus(qs.status) {
		qs.mu.Unlock()
		return
	}
	qs.markTerminal(queryStatusTimeout)
	qs.mu.Unlock()
	if err := s.persistQueryState(qs); err != nil {
		logs.Error("Failed to persist timed-out query state",
			logs.String("queryId", queryId), logs.Err(err))
	}
}

func (s *LogsService) failQuery(queryId, message string) {
	val, ok := s.queries.Load(queryId)
	if !ok {
		return
	}
	qs := val.(*queryState)
	qs.mu.Lock()
	// Terminal status wins: a query StopQuery already cancelled (or that
	// completed before the failure surfaced) keeps its status.
	if isTerminalQueryStatus(qs.status) {
		qs.mu.Unlock()
		return
	}
	qs.markTerminal(queryStatusFailed)
	qs.mu.Unlock()
	if err := s.persistQueryState(qs); err != nil {
		logs.Error("Failed to persist failed query state",
			logs.String("queryId", queryId), logs.Err(err))
	}
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
	items, nextToken, err := s.describeQueriesCore(&DescribeQueriesInput{
		StatusFilter: request.GetParamLowerFirst(req.Parameters, "Status"),
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		QueryLanguage: request.GetParamLowerFirst(
			req.Parameters, "QueryLanguage"),
		NextToken:  request.GetParamLowerFirst(req.Parameters, "NextToken"),
		MaxResults: int32(request.GetIntParam(req.Parameters, "MaxResults")),
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
	result, err := s.getQueryResultsCore(&GetQueryResultsInput{
		QueryId:   request.GetParamLowerFirst(req.Parameters, "QueryId"),
		NextToken: request.GetParamLowerFirst(req.Parameters, "NextToken"),
		MaxItems:  int32(request.GetIntParam(req.Parameters, "MaxItems")),
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
		"logGroupsScanned":        result.LogGroupsScanned,
		"resultCount":             result.ResultCount,
	}

	resp := map[string]interface{}{
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
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	logGroupName := ""
	if len(qs.logGroupNames) > 0 {
		logGroupName = qs.logGroupNames[0]
	}
	// createTime is epoch seconds — the Timestamp shape's wire format,
	// matching the documented example ("createTime": 1540923785). The
	// duration carries the terminal stamp once execution has ended and
	// the elapsed wall clock while the query runs (the worked example
	// shows a Running query carrying its duration so far).
	duration := qs.durationMs
	if duration == 0 && qs.status == queryStatusRunning {
		duration = time.Since(qs.createdAt).Milliseconds()
	}
	entry := map[string]interface{}{
		"queryId":       qs.queryId,
		"queryString":   qs.queryString,
		"status":        qs.status,
		"createTime":    qs.createdAt.Unix(),
		"logGroupName":  logGroupName,
		"queryDuration": duration,
		"bytesScanned":  qs.stats.bytesScanned,
		"queryLanguage": qs.queryLanguage,
	}
	if qs.userIdentity != "" {
		entry["userIdentity"] = qs.userIdentity
	}
	return entry
}
