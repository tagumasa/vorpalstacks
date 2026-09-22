package cloudwatchlogs

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// StartQueryInput holds validated parameters for StartQuery.
type StartQueryInput struct {
	StartTime     int64
	EndTime       int64
	QueryString   string
	QueryLanguage string
	// The three required members' zero values sit inside their
	// documented domains (startTime/endTime "Minimum value of 0",
	// queryString "Minimum length of 0"), so wire presence — not the
	// value — carries requiredness.
	StartTimeSet   bool
	EndTimeSet     bool
	QueryStringSet bool
	// The three group members mirror the operation's wire shape; Core
	// merges them into the execution name list after the exactly-one-of
	// rule has counted them.
	LogGroupName        string
	LogGroupNames       []string
	LogGroupIdentifiers []string
	Limit               int64
	Region              string
	// UserIdentity is the ARN of the principal that started the query;
	// empty when the caller is anonymous.
	UserIdentity string
}

// queryStringSelectsWithSource reports whether the query's first command
// is a SOURCE command — the selection form a query without group members
// must carry ("the queryString must include a SOURCE command to select
// log groups for the query"). The pipeline compiler keeps SOURCE valid
// only as the first command and the lexer skips leading comments, so the
// first lexed token decides.
func queryStringSelectsWithSource(qs string) bool {
	toks, err := lexQuery(qs)
	if err != nil || len(toks) == 0 {
		return false
	}
	return toks[0].kind == tokIdent && strings.EqualFold(toks[0].text, "source")
}

// startQueryCore validates input, creates the query state, and launches
// the async execution worker under the service's task runner.
func (s *LogsService) startQueryCore(input *StartQueryInput) (string, error) {
	// The members are documented as epoch seconds — "Specified as epoch
	// time, the number of seconds since January 1, 1970" — while the
	// engine, the persisted query state and the read plane all carry
	// milliseconds; the unit changes once, at this seam both planes ride.
	input.StartTime *= 1000
	input.EndTime *= 1000
	// Requiredness rides the presence flags; the zero values are inside
	// the members' domains, and only negatives reject ("Minimum value of
	// 0").
	if !input.QueryStringSet || !input.StartTimeSet || !input.EndTimeSet {
		return "", errRequiredMember("queryString, startTime and endTime")
	}
	if input.StartTime < 0 || input.EndTime < 0 {
		return "", NewLogsError("InvalidParameterException",
			"startTime and endTime must not be negative", 400)
	}
	if err := validateQueryString(input.QueryString, 0); err != nil {
		return "", err
	}
	// The time range is a half-open concept with inclusive ends; a window
	// whose end precedes its start is malformed and rejects with the
	// model-declared InvalidParameterException instead of silently
	// completing with zero rows.
	if input.StartTime >= input.EndTime {
		return "", NewLogsError("InvalidParameterException",
			"startTime must be earlier than endTime", 400)
	}
	// "A StartQuery operation must include one of the following: Either
	// exactly one of the following parameters: logGroupName,
	// logGroupNames, or logGroupIdentifiers Or the queryString must
	// include a SOURCE command to select log groups for the query." Two
	// supplied members reject; none supplied accepts only a SOURCE-first
	// query, whose groups resolve at execution.
	groupMembers := 0
	if input.LogGroupName != "" {
		groupMembers++
	}
	if len(input.LogGroupNames) > 0 {
		groupMembers++
	}
	if len(input.LogGroupIdentifiers) > 0 {
		groupMembers++
	}
	if groupMembers > 1 {
		return "", NewLogsError("InvalidParameterException",
			"Specify exactly one of logGroupName, logGroupNames or logGroupIdentifiers", 400)
	}
	if groupMembers == 0 && !queryStringSelectsWithSource(input.QueryString) {
		return "", errRequiredMember("logGroupName, logGroupNames or logGroupIdentifiers")
	}
	if len(input.LogGroupNames)+len(input.LogGroupIdentifiers) > logsstore.QueryLogGroupsMax {
		return "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("A query can include up to %d log groups", logsstore.QueryLogGroupsMax), 400)
	}
	queryLanguage := input.QueryLanguage
	if queryLanguage == "" {
		queryLanguage = "CWLI"
	}
	if !validateQueryLanguage(queryLanguage) {
		return "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid queryLanguage: %s. Allowed values: CWLI, SQL, PPL", queryLanguage), 400)
	}
	if err := validateQueryPipeline(input.QueryString); err != nil {
		return "", err
	}
	limit32, err := validateListLimit(int32(input.Limit), logsstore.DefaultStartQueryLimit, logsstore.MaxStartQueryLimit)
	if err != nil {
		return "", err
	}

	// "You can have up to 100 concurrent CloudWatch Logs insights
	// queries, including queries that have been added to dashboards" —
	// the registry's non-terminal states (running queries and scheduled
	// executions in progress) occupy the account's concurrency, and the
	// query that would exceed the quota rejects with the operation's
	// declared LimitExceededException instead of starting. The census and
	// the registration run under the admission mutex as one critical
	// section, so two concurrent starts cannot both observe the quota
	// open and both admit.
	s.queryAdmissionMu.Lock()
	defer s.queryAdmissionMu.Unlock()
	running := 0
	s.queries.Range(func(key, value interface{}) bool {
		qs := value.(*queryState)
		qs.mu.RLock()
		status := qs.status
		qs.mu.RUnlock()
		if !isTerminalQueryStatus(status) {
			running++
		}
		return true
	})
	if running >= logsstore.MaxConcurrentQueries {
		return "", NewLogsError("LimitExceededException",
			fmt.Sprintf("You have reached the maximum number of %d concurrent queries", logsstore.MaxConcurrentQueries), 400)
	}

	// The identifiers resolve to names here, so the execution plane and the
	// SOURCE group scan address store name keys whether the caller passed
	// names or ARNs; the raw list stays on the query state untouched.
	var allGroups []string
	if input.LogGroupName != "" {
		allGroups = append(allGroups, input.LogGroupName)
	}
	allGroups = append(allGroups, input.LogGroupNames...)
	allGroups = append(allGroups, resolveLogGroupIdentifiers(input.LogGroupIdentifiers)...)

	// Existence gate: AWS StartQuery fails fast on a log group that does
	// not exist (the operation declares ResourceNotFoundException) rather
	// than returning a queryId that completes with zero rows. Queries that
	// select their groups through a SOURCE command carry no explicit list
	// and stay ungated here — their groups resolve at execution.
	if len(allGroups) > 0 {
		store, err := s.getLogsStoreByRegion(input.Region)
		if err != nil {
			return "", err
		}
		for _, name := range allGroups {
			if _, err := store.GetLogGroup(name); err != nil {
				if errors.Is(err, logsstore.ErrLogGroupNotFound) {
					return "", NewLogsError("ResourceNotFoundException",
						fmt.Sprintf("The specified log group does not exist: %s", name), 400)
				}
				return "", err
			}
		}
	}

	queryId := fmt.Sprintf("query-%d", time.Now().UnixNano())

	qs := &queryState{
		queryId:             queryId,
		logGroupNames:       allGroups,
		logGroupIdentifiers: input.LogGroupIdentifiers,
		scannedGroups:       allGroups,
		startTime:           input.StartTime,
		endTime:             input.EndTime,
		queryString:         input.QueryString,
		queryLanguage:       queryLanguage,
		status:              queryStatusRunning,
		createdAt:           time.Now(),
		region:              input.Region,
		userIdentity:        input.UserIdentity,
	}
	// The record is written before the state goes live: a StartQuery that
	// cannot persist its state has not started, and an unpersisted live
	// state would be an availability the restart path never promised.
	if err := s.persistQueryState(qs); err != nil {
		return "", fmt.Errorf("persist query state: %w", err)
	}
	s.queries.Store(queryId, qs)

	s.spawnTask(func(ctx context.Context) {
		s.executeQuery(ctx, input.Region, queryId, input.QueryString, allGroups, input.StartTime, input.EndTime, int64(limit32))
	})

	return queryId, nil
}

// persistQueryState writes the state's current snapshot through to its
// region's store, mirroring the in-memory registry so DescribeQueries and
// GetQueryResults serve pre-restart queries for the retention window the
// way AWS keeps results queryable. The snapshot is taken under the read
// lock; callers must not hold the write lock. A state with no region has
// no store to write through to — every production state carries the
// StartQuery region or the loaded record's region.
func (s *LogsService) persistQueryState(qs *queryState) error {
	if qs.region == "" {
		return nil
	}
	qs.mu.RLock()
	rec := &logsstore.QueryRecord{
		QueryId:             qs.queryId,
		Region:              qs.region,
		LogGroupNames:       qs.logGroupNames,
		LogGroupIdentifiers: qs.logGroupIdentifiers,
		ScannedGroups:       qs.scannedGroups,
		StartTime:           qs.startTime,
		EndTime:             qs.endTime,
		QueryString:         qs.queryString,
		QueryLanguage:       qs.queryLanguage,
		Status:              qs.status,
		CreatedAtUnixNano:   qs.createdAt.UnixNano(),
		Results:             make([]logsstore.QueryResultRow, len(qs.results)),
		RecordsScanned:      qs.stats.recordsScanned,
		RecordsMatched:      qs.stats.recordsMatched,
		BytesScanned:        qs.stats.bytesScanned,
		UserIdentity:        qs.userIdentity,
		QueryDurationMs:     qs.durationMs,
	}
	for i := range qs.results {
		rec.Results[i] = logsstore.QueryResultRow{
			Columns: qs.results[i].columns,
			Fields:  qs.results[i].fields,
		}
	}
	qs.mu.RUnlock()

	store, err := s.getLogsStoreByRegion(qs.region)
	if err != nil {
		return err
	}
	return store.PutQueryRecord(rec)
}

// loadPersistedQueries rebuilds the in-memory query registry from the
// persisted records, so a restart is not the availability boundary for
// query state. A non-terminal record has no worker behind it after a
// restart — the same reconciliation rule export and import tasks follow:
// Failed with the interruption reason. Records already past the retention
// window are dropped here rather than resurrected until the hourly sweep.
func (s *LogsService) loadPersistedQueries() {
	if s.storageManager == nil {
		return
	}
	cutoff := time.Now().Add(-logsstore.QueryResultRetentionPeriod)
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			continue
		}
		records, err := store.ListQueryRecords()
		if err != nil {
			continue
		}
		for _, rec := range records {
			createdAt := time.Unix(0, rec.CreatedAtUnixNano)
			if createdAt.Before(cutoff) {
				if err := store.DeleteQueryRecord(rec.QueryId); err != nil {
					logs.Error("Failed to drop expired query record",
						logs.String("queryId", rec.QueryId), logs.Err(err))
				}
				continue
			}
			qs := queryStateFromRecord(rec)
			if !isTerminalQueryStatus(qs.status) {
				qs.status = queryStatusFailed
				rec.Status = queryStatusFailed
				if err := store.PutQueryRecord(rec); err != nil {
					logs.Error("Failed to reconcile interrupted query record",
						logs.String("queryId", rec.QueryId), logs.Err(err))
				}
			}
			s.queries.Store(qs.queryId, qs)
		}
	}
}

// queryStateFromRecord rebuilds the live query state from its persisted
// form. Result rows keep their column order through the record's Columns
// half, so a reloaded row renders exactly as the pre-restart original.
func queryStateFromRecord(rec *logsstore.QueryRecord) *queryState {
	qs := &queryState{
		queryId:             rec.QueryId,
		logGroupNames:       rec.LogGroupNames,
		logGroupIdentifiers: rec.LogGroupIdentifiers,
		startTime:           rec.StartTime,
		endTime:             rec.EndTime,
		queryString:         rec.QueryString,
		queryLanguage:       rec.QueryLanguage,
		status:              rec.Status,
		scannedGroups:       rec.ScannedGroups,
		results:             make([]queryResultRow, len(rec.Results)),
		stats: queryStats{
			recordsScanned: rec.RecordsScanned,
			recordsMatched: rec.RecordsMatched,
			bytesScanned:   rec.BytesScanned,
		},
		createdAt:    time.Unix(0, rec.CreatedAtUnixNano).UTC(),
		region:       rec.Region,
		userIdentity: rec.UserIdentity,
		durationMs:   rec.QueryDurationMs,
	}
	for i := range rec.Results {
		fields := rec.Results[i].Fields
		if fields == nil {
			fields = make(map[string]string)
		}
		qs.results[i] = queryResultRow{fields: fields, columns: rec.Results[i].Columns}
	}
	return qs
}

// stopQueryCore validates input and cancels a running query. Stopping a
// query that already ended rejects with InvalidParameterException — the
// StopQuery documentation states an ended query returns an error
// indicating the specified query is not running. Scheduled-query
// executions are stoppable through the same operation per the same
// documentation; their cancel path marks the persisted execution record.
func (s *LogsService) stopQueryCore(queryId string) error {
	if queryId == "" {
		return errRequiredMember("queryId")
	}

	val, ok := s.queries.Load(queryId)
	if !ok {
		found, running := s.cancelScheduledExecution(queryId)
		if !found {
			return NewLogsError("ResourceNotFoundException",
				fmt.Sprintf("Query %s not found", queryId), 400)
		}
		if !running {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Query %s is not running", queryId), 400)
		}
		return nil
	}

	qs := val.(*queryState)
	qs.mu.Lock()
	if isTerminalQueryStatus(qs.status) {
		qs.mu.Unlock()
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Query %s is not running", queryId), 400)
	}
	qs.cancelled = true
	qs.markTerminal(queryStatusCancelled)
	qs.mu.Unlock()
	// The cancellation writes through like every terminal transition: a
	// restarted process must still answer that this query was cancelled.
	if err := s.persistQueryState(qs); err != nil {
		logs.Error("Failed to persist cancelled query state",
			logs.String("queryId", queryId), logs.Err(err))
	}
	return nil
}

// cancelScheduledExecution locates a scheduled-query execution by its
// query ID across the configured regions and, when it is still running,
// marks the persisted record CANCELLED. The running execution observes
// the terminal status at its delivery checkpoint; a finished execution is
// reported back so the caller can reject the stop. found reports whether
// any execution carries the query ID at all.
func (s *LogsService) cancelScheduledExecution(queryId string) (found, running bool) {
	if s.storageManager == nil {
		return false, false
	}
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			continue
		}
		scheduled, err := store.ListScheduledQueries("")
		if err != nil {
			continue
		}
		for _, sq := range scheduled {
			execs, err := store.ListScheduledQueryExecutions(sq.Id, 0, 0)
			if err != nil {
				continue
			}
			for _, exec := range execs {
				if exec.QueryId != queryId {
					continue
				}
				if exec.Status != logsstore.ScheduledExecutionStatusRunning {
					return true, false
				}
				// The RUNNING→CANCELLED flip is one atomic record
				// mutation, so a worker terminal write cannot land
				// between the check and the flip.
				cancelled, err := store.CancelScheduledQueryExecutionIfRunning(exec)
				if err != nil {
					logs.Error("Failed to persist cancelled scheduled query execution",
						logs.String("queryId", queryId), logs.Err(err))
					return true, false
				}
				return true, cancelled
			}
		}
	}
	return false, false
}

// DescribeQueriesInput holds parameters for DescribeQueries.
type DescribeQueriesInput struct {
	StatusFilter  string
	LogGroupName  string
	QueryLanguage string
	NextToken     string
	MaxResults    int32
}

// describeQueriesCore validates input and lists queries with filtering.
// The listing is deterministic — newest first, query id as the tiebreak —
// and pages through the scoped queryPageToken vocabulary, so successive
// pages neither repeat nor miss queries as the sync.Map's random
// iteration order did. Scheduled-query executions appear alongside
// interactive queries (the DescribeQueries documentation states they are
// included), carrying the Scheduled status while in progress and their
// terminal status afterwards.
func (s *LogsService) describeQueriesCore(input *DescribeQueriesInput) ([]*queryState, string, error) {
	maxResults, err := validateListLimit(input.MaxResults, logsstore.DefaultListMaxResults, logsstore.MaxListMaxResults)
	if err != nil {
		return nil, "", err
	}
	// The status and queryLanguage members are enumerations (QueryStatus,
	// QueryLanguage); a value outside the vocabulary rejects with the
	// operation's declared InvalidParameterException. The typed SDK
	// blocks non-enum values client-side, so these rows are raw-HTTP
	// paths, unit-pinned.
	if input.StatusFilter != "" && !validQueryStatusFilter[input.StatusFilter] {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid status: %s. Valid values: Scheduled, Running, Complete, Failed, Cancelled, Timeout, Unknown", input.StatusFilter), 400)
	}
	if !validateQueryLanguage(input.QueryLanguage) {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid queryLanguage: %s. Allowed values: CWLI, SQL, PPL", input.QueryLanguage), 400)
	}

	var allQueries []*queryState
	s.queries.Range(func(key, value interface{}) bool {
		qs := value.(*queryState)
		if !queryMatchesFilters(qs, input) {
			return true
		}
		allQueries = append(allQueries, qs)
		return true
	})
	for _, qs := range s.scheduledExecutionQueries() {
		if queryMatchesFilters(qs, input) {
			allQueries = append(allQueries, qs)
		}
	}

	// Newest first; the query id tiebreak keeps the order total so the
	// value cursor is exact.
	sort.Slice(allQueries, func(i, j int) bool {
		ti, tj := allQueries[i].createdAt.UnixNano(), allQueries[j].createdAt.UnixNano()
		if ti != tj {
			return ti > tj
		}
		return allQueries[i].queryId < allQueries[j].queryId
	})

	var cursor queryPageToken
	if input.NextToken != "" {
		cursor, err = decodeQueryPageToken(input.NextToken)
		if err != nil {
			return nil, "", err
		}
		if cursor.StatusFilter != input.StatusFilter || cursor.LogGroupName != input.LogGroupName ||
			cursor.QueryLanguage != input.QueryLanguage {
			return nil, "", errInvalidPageToken
		}
	}

	// Resume at the first query that sorts strictly after the cursor's
	// (creation time, query id) key in the listing's total order: an
	// older timestamp, or an equal timestamp with a larger query id.
	// The tie-break direction mirrors the sort's ascending-id branch —
	// inverted, an equal-timestamp group spanning a page boundary would
	// re-serve the served smaller ids (a cursor re-mint loop) or, when
	// no smaller id exists, drop the group's remainder.
	startIdx := 0
	if input.NextToken != "" {
		startIdx = len(allQueries)
		for i, qs := range allQueries {
			if qs.createdAt.UnixNano() < cursor.CursorCreateNs ||
				(qs.createdAt.UnixNano() == cursor.CursorCreateNs && qs.queryId > cursor.CursorQueryId) {
				startIdx = i
				break
			}
		}
	}

	endIdx := startIdx + int(maxResults)
	if endIdx > len(allQueries) {
		endIdx = len(allQueries)
	}

	result := &queryPageResult{}
	if startIdx < endIdx {
		result.Items = allQueries[startIdx:endIdx]
		if endIdx < len(allQueries) {
			last := allQueries[endIdx-1]
			result.NextToken, err = encodeQueryPageToken(queryPageToken{
				StatusFilter:   input.StatusFilter,
				LogGroupName:   input.LogGroupName,
				QueryLanguage:  input.QueryLanguage,
				CursorCreateNs: last.createdAt.UnixNano(),
				CursorQueryId:  last.queryId,
			})
			if err != nil {
				return nil, "", err
			}
		}
	}
	return result.Items, result.NextToken, nil
}

// GetQueryResultsInput holds parameters for GetQueryResults.
type GetQueryResultsInput struct {
	QueryId   string
	NextToken string
	MaxItems  int32
}

// GetQueryResultsResult holds the query results page.
type GetQueryResultsResult struct {
	Status        string
	QueryLanguage string
	Results       []queryResultRow
	Stats         queryStats
	// LogGroupsScanned and ResultCount complete the statistics envelope:
	// the scanned-set size and the total rows across every page, both of
	// which live on the query state rather than this page's slice.
	LogGroupsScanned int
	ResultCount      int
	NextToken        string
}

// getQueryResultsCore validates input and retrieves paginated query results.
func (s *LogsService) getQueryResultsCore(input *GetQueryResultsInput) (*GetQueryResultsResult, error) {
	if input.QueryId == "" {
		return nil, errRequiredMember("queryId")
	}

	val, ok := s.queries.Load(input.QueryId)
	if !ok {
		return nil, NewLogsError("ResourceNotFoundException",
			fmt.Sprintf("Query %s not found", input.QueryId), 400)
	}

	qs := val.(*queryState)

	limit32, err := validateListLimit(input.MaxItems, logsstore.DefaultQueryResultsItems, logsstore.MaxQueryResultsItems)
	if err != nil {
		return nil, err
	}
	limit := int(limit32)

	offset := 0
	if input.NextToken != "" {
		cursor, err := decodeQueryResultsPageToken(input.NextToken)
		if err != nil || cursor.QueryId != input.QueryId {
			return nil, NewLogsError("InvalidParameterException",
				"Invalid nextToken", 400)
		}
		offset = cursor.Offset
	}

	// Snapshot under the read lock: the executing worker may be writing
	// the results slice header concurrently.
	qs.mu.RLock()
	status := qs.status
	queryLanguage := qs.queryLanguage
	results := qs.results
	stats := qs.stats
	scannedGroups := len(qs.scannedGroups)
	qs.mu.RUnlock()

	endIdx := offset + limit
	if endIdx > len(results) {
		endIdx = len(results)
	}
	if offset > len(results) {
		offset = len(results)
	}

	result := &GetQueryResultsResult{
		Status:           status,
		QueryLanguage:    queryLanguage,
		Results:          results[offset:endIdx],
		Stats:            stats,
		LogGroupsScanned: scannedGroups,
		ResultCount:      len(results),
	}

	if endIdx < len(results) {
		result.NextToken, err = encodeQueryResultsPageToken(queryResultsPageToken{
			QueryId: input.QueryId, Offset: endIdx,
		})
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// queryMatchesFilters applies the DescribeQueries filter members to one
// query state under its read lock. An absent queryLanguage on the filter
// side matches everything; the CWLI default normalises the state side so
// queries started without an explicit language still match a CWLI filter.
func queryMatchesFilters(qs *queryState, input *DescribeQueriesInput) bool {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	if input.StatusFilter != "" && qs.status != input.StatusFilter {
		return false
	}
	if input.LogGroupName != "" {
		found := false
		for _, n := range qs.logGroupNames {
			if n == input.LogGroupName {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if input.QueryLanguage != "" {
		language := qs.queryLanguage
		if language == "" {
			language = "CWLI"
		}
		if language != input.QueryLanguage {
			return false
		}
	}
	return true
}

// scheduledExecutionQueries materialises the scheduled-query executions
// inside the query-result retention window as query states so
// DescribeQueries lists them alongside interactive queries, per the
// operation documentation. The region set comes from the storage manager,
// not the store cache — listing must not depend on foreground traffic
// having touched the region.
func (s *LogsService) scheduledExecutionQueries() []*queryState {
	if s.storageManager == nil {
		return nil
	}
	windowStart := time.Now().Add(-logsstore.QueryResultRetentionPeriod).UnixMilli()
	var out []*queryState
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			continue
		}
		scheduled, err := store.ListScheduledQueries("")
		if err != nil {
			continue
		}
		for _, sq := range scheduled {
			execs, err := store.ListScheduledQueryExecutions(sq.Id, windowStart, 0)
			if err != nil {
				continue
			}
			for _, exec := range execs {
				out = append(out, scheduledExecutionState(sq, exec))
			}
		}
	}
	return out
}

// scheduledExecutionState maps one persisted scheduled-query execution to
// the query-status vocabulary: an in-progress execution carries the
// Scheduled status the QueryStatus enum defines, a succeeded one Complete,
// a failed one Failed, a cancelled one Cancelled. The input vocabulary is
// the scheduled-execution record constants; mapExecutionStatus (the
// scheduled-query operations file) maps the same inputs onto the
// ExecutionStatus enum GetScheduledQueryHistory reports — two documented
// output enums, one input vocabulary.
func scheduledExecutionState(sq *logsstore.ScheduledQuery, exec *logsstore.ScheduledQueryExecution) *queryState {
	status := queryStatusFailed
	switch exec.Status {
	case logsstore.ScheduledExecutionStatusRunning:
		status = queryStatusScheduled
	case logsstore.ScheduledExecutionStatusSuccess:
		status = queryStatusComplete
	case logsstore.ScheduledExecutionStatusCancelled:
		status = queryStatusCancelled
	}
	queryLanguage := sq.QueryLanguage
	if queryLanguage == "" {
		queryLanguage = "CWLI"
	}
	return &queryState{
		queryId:       exec.QueryId,
		logGroupNames: resolveLogGroupIdentifiers(sq.LogGroupIdentifiers),
		queryString:   sq.QueryString,
		queryLanguage: queryLanguage,
		status:        status,
		createdAt:     time.UnixMilli(exec.TriggerTime).UTC(),
	}
}

type queryPageResult struct {
	Items     []*queryState
	NextToken string
}
