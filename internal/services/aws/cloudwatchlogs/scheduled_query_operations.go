package cloudwatchlogs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// CreateScheduledQuery creates a scheduled CloudWatch Logs Insights query.
func (s *LogsService) CreateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var logGroupIdentifiers []string
	if idents, ok := req.Parameters["logGroupIdentifiers"]; ok {
		if arr, ok := idents.([]interface{}); ok {
			for _, item := range arr {
				if str, ok := item.(string); ok {
					logGroupIdentifiers = append(logGroupIdentifiers, str)
				}
			}
		}
	}

	var destinationConfiguration map[string]interface{}
	if dc, ok := req.Parameters["destinationConfiguration"]; ok {
		if m, ok := dc.(map[string]interface{}); ok {
			destinationConfiguration = m
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sq, err := s.createScheduledQueryCore(store, &CreateScheduledQueryInput{
		Name:                     request.GetParamLowerFirst(req.Parameters, "Name"),
		Description:              request.GetParamLowerFirst(req.Parameters, "Description"),
		QueryString:              request.GetParamLowerFirst(req.Parameters, "QueryString"),
		QueryLanguage:            request.GetParamLowerFirst(req.Parameters, "QueryLanguage"),
		LogGroupIdentifiers:      logGroupIdentifiers,
		ScheduleExpression:       request.GetParamLowerFirst(req.Parameters, "ScheduleExpression"),
		State:                    request.GetParamLowerFirst(req.Parameters, "State"),
		ExecutionRoleArn:         request.GetParamLowerFirst(req.Parameters, "ExecutionRoleArn"),
		Timezone:                 request.GetParamLowerFirst(req.Parameters, "Timezone"),
		StartTimeOffset:          int64(request.GetIntParam(req.Parameters, "StartTimeOffset")),
		EndTimeOffset:            int64(request.GetIntParam(req.Parameters, "EndTimeOffset")),
		ScheduleStartTime:        int64(request.GetIntParam(req.Parameters, "ScheduleStartTime")),
		ScheduleEndTime:          int64(request.GetIntParam(req.Parameters, "ScheduleEndTime")),
		DestinationConfiguration: destinationConfiguration,
		Tags:                     parseTagsFromParams(req.Parameters),
	})
	if err != nil {
		return nil, err
	}

	arn := fmt.Sprintf("arn:aws:logs:%s:%s:scheduled-query:%s", reqCtx.GetRegion(), s.accountID, sq.Id)
	return map[string]interface{}{
		"scheduledQueryArn": arn,
		"state":             sq.State,
	}, nil
}

// DeleteScheduledQuery deletes a scheduled query.
func (s *LogsService) DeleteScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "Identifier")
	id := extractIdFromArnOrId(identifier)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteScheduledQueryCore(store, id); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// UpdateScheduledQuery updates a scheduled query.
func (s *LogsService) UpdateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var logGroupIdentifiers []string
	if idents, ok := req.Parameters["logGroupIdentifiers"]; ok {
		if arr, ok := idents.([]interface{}); ok {
			for _, item := range arr {
				if str, ok := item.(string); ok {
					logGroupIdentifiers = append(logGroupIdentifiers, str)
				}
			}
		}
	}
	var destinationConfiguration map[string]interface{}
	if dc, ok := req.Parameters["destinationConfiguration"]; ok {
		if m, ok := dc.(map[string]interface{}); ok {
			destinationConfiguration = m
		}
	}
	in := &UpdateScheduledQueryInput{
		Identifier:               request.GetParamLowerFirst(req.Parameters, "Identifier"),
		Description:              request.GetParamLowerFirst(req.Parameters, "Description"),
		QueryString:              request.GetParamLowerFirst(req.Parameters, "QueryString"),
		QueryLanguage:            request.GetParamLowerFirst(req.Parameters, "QueryLanguage"),
		LogGroupIdentifiers:      logGroupIdentifiers,
		ScheduleExpression:       request.GetParamLowerFirst(req.Parameters, "ScheduleExpression"),
		State:                    request.GetParamLowerFirst(req.Parameters, "State"),
		ExecutionRoleArn:         request.GetParamLowerFirst(req.Parameters, "ExecutionRoleArn"),
		Timezone:                 request.GetParamLowerFirst(req.Parameters, "Timezone"),
		DestinationConfiguration: destinationConfiguration,
	}
	if v, present := request.GetIntParamCaseInsensitive(req.Parameters, "StartTimeOffset"); present {
		v64 := int64(v)
		in.StartTimeOffset = &v64
	}
	if v, present := request.GetIntParamCaseInsensitive(req.Parameters, "EndTimeOffset"); present {
		v64 := int64(v)
		in.EndTimeOffset = &v64
	}
	if v, present := request.GetIntParamCaseInsensitive(req.Parameters, "ScheduleStartTime"); present {
		v64 := int64(v)
		in.ScheduleStartTime = &v64
	}
	if v, present := request.GetIntParamCaseInsensitive(req.Parameters, "ScheduleEndTime"); present {
		v64 := int64(v)
		in.ScheduleEndTime = &v64
	}
	if _, ok := req.Parameters["logGroupIdentifiers"]; !ok {
		in.LogGroupIdentifiers = nil
	}

	sq, err := s.updateScheduledQueryCore(store, in)
	if err != nil {
		return nil, err
	}

	return formatScheduledQuery(sq, reqCtx.GetRegion(), s.accountID), nil
}

// GetScheduledQuery retrieves a scheduled query.
func (s *LogsService) GetScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "Identifier")
	if identifier == "" {
		return nil, ErrMissingParameter
	}

	id := extractIdFromArnOrId(identifier)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sq, err := s.getScheduledQueryCore(store, id)
	if err != nil {
		return nil, err
	}

	return formatScheduledQuery(sq, reqCtx.GetRegion(), s.accountID), nil
}

// GetScheduledQueryHistory retrieves execution history for a scheduled query.
func (s *LogsService) GetScheduledQueryHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var executionStatuses []string
	if es, ok := req.Parameters["executionStatuses"]; ok {
		if arr, ok := es.([]interface{}); ok {
			for _, item := range arr {
				if str, ok := item.(string); ok {
					executionStatuses = append(executionStatuses, str)
				}
			}
		}
	}

	maxResults, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "MaxResults")), 50, 1000)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	history, err := s.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
		Identifier:        request.GetParamLowerFirst(req.Parameters, "Identifier"),
		StartTime:         int64(request.GetIntParam(req.Parameters, "StartTime")),
		EndTime:           int64(request.GetIntParam(req.Parameters, "EndTime")),
		ExecutionStatuses: executionStatuses,
		NextToken:         request.GetParamLowerFirst(req.Parameters, "NextToken"),
		MaxResults:        maxResults,
	})
	if err != nil {
		return nil, err
	}

	id := extractIdFromArnOrId(request.GetParamLowerFirst(req.Parameters, "Identifier"))

	entries := make([]map[string]interface{}, len(history.Executions))
	for i, exec := range history.Executions {
		entry := map[string]interface{}{
			"queryId":            exec.QueryId,
			"triggeredTimestamp": exec.TriggerTime,
			"executionStatus":    mapExecutionStatus(exec.Status),
		}
		if exec.ErrorMessage != "" {
			entry["errorMessage"] = exec.ErrorMessage
		}
		destinations := make([]interface{}, 0, len(exec.Destinations))
		for _, dest := range exec.Destinations {
			d := map[string]interface{}{
				"destinationType":       dest.DestinationType,
				"destinationIdentifier": dest.DestinationIdentifier,
				"status":                dest.Status,
			}
			if dest.ProcessedIdentifier != "" {
				d["processedIdentifier"] = dest.ProcessedIdentifier
			}
			if dest.ErrorMessage != "" {
				d["errorMessage"] = dest.ErrorMessage
			}
			destinations = append(destinations, d)
		}
		entry["destinations"] = destinations
		entries[i] = entry
	}

	arn := fmt.Sprintf("arn:aws:logs:%s:%s:scheduled-query:%s", reqCtx.GetRegion(), s.accountID, id)

	resp := map[string]interface{}{
		"scheduledQueryArn": arn,
		"triggerHistory":    entries,
	}
	if history.Query != nil {
		resp["name"] = history.Query.Name
	}

	if history.NextMarker != "" {
		resp["nextToken"] = history.NextMarker
	}

	return resp, nil
}

func mapExecutionStatus(internalStatus string) string {
	switch internalStatus {
	case "RUNNING":
		return "Running"
	case "SUCCESS":
		return "Complete"
	case "FAILED":
		return "Failed"
	case "INVALID_QUERY":
		return "InvalidQuery"
	case "TIMEOUT":
		return "Timeout"
	default:
		return internalStatus
	}
}

// ListScheduledQueries lists scheduled queries.
func (s *LogsService) ListScheduledQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "MaxResults")), 50, 1000)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.listScheduledQueriesCore(store,
		request.GetParamLowerFirst(req.Parameters, "State"),
		request.GetParamLowerFirst(req.Parameters, "ScheduleType"),
		maxResults,
		request.GetParamLowerFirst(req.Parameters, "NextToken"))
	if err != nil {
		return nil, err
	}

	region := reqCtx.GetRegion()
	queries := make([]map[string]interface{}, len(items))
	for i, sq := range items {
		queries[i] = formatScheduledQuerySummary(sq, region, s.accountID)
	}

	resp := map[string]interface{}{
		"scheduledQueries": queries,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
}

func formatScheduledQuery(sq *logsstore.ScheduledQuery, region, accountID string) map[string]interface{} {
	arn := fmt.Sprintf("arn:aws:logs:%s:%s:scheduled-query:%s", region, accountID, sq.Id)
	result := map[string]interface{}{
		"scheduledQueryArn":  arn,
		"name":               sq.Name,
		"queryString":        sq.QueryString,
		"scheduleExpression": sq.ScheduleExpression,
		"state":              sq.State,
		"creationTime":       sq.CreationTime,
		"lastUpdatedTime":    sq.LastUpdatedTime,
		"scheduleType":       sq.ScheduleType,
	}
	if sq.Description != "" {
		result["description"] = sq.Description
	}
	if sq.QueryLanguage != "" {
		result["queryLanguage"] = sq.QueryLanguage
	}
	if len(sq.LogGroupIdentifiers) > 0 {
		result["logGroupIdentifiers"] = sq.LogGroupIdentifiers
	}
	if sq.ExecutionRoleArn != "" {
		result["executionRoleArn"] = sq.ExecutionRoleArn
	}
	if sq.Timezone != "" {
		result["timezone"] = sq.Timezone
	}
	if sq.StartTimeOffset != 0 {
		result["startTimeOffset"] = sq.StartTimeOffset
	}
	if sq.EndTimeOffset != 0 {
		result["endTimeOffset"] = sq.EndTimeOffset
	}
	if sq.ScheduleStartTime != 0 {
		result["scheduleStartTime"] = sq.ScheduleStartTime
	}
	if sq.ScheduleEndTime != 0 {
		result["scheduleEndTime"] = sq.ScheduleEndTime
	}
	if sq.DestinationConfiguration != nil {
		result["destinationConfiguration"] = sq.DestinationConfiguration
	}
	if sq.LastTriggeredTime > 0 {
		result["lastTriggeredTime"] = sq.LastTriggeredTime
	}
	if sq.LastExecutionStatus != "" {
		result["lastExecutionStatus"] = sq.LastExecutionStatus
	}
	return result
}

func formatScheduledQuerySummary(sq *logsstore.ScheduledQuery, region, accountID string) map[string]interface{} {
	arn := fmt.Sprintf("arn:aws:logs:%s:%s:scheduled-query:%s", region, accountID, sq.Id)
	result := map[string]interface{}{
		"scheduledQueryArn":        arn,
		"name":                     sq.Name,
		"state":                    sq.State,
		"scheduleType":             sq.ScheduleType,
		"scheduleExpression":       sq.ScheduleExpression,
		"timezone":                 sq.Timezone,
		"destinationConfiguration": sq.DestinationConfiguration,
		"creationTime":             sq.CreationTime,
		"lastUpdatedTime":          sq.LastUpdatedTime,
	}
	if sq.LastTriggeredTime > 0 {
		result["lastTriggeredTime"] = sq.LastTriggeredTime
	}
	if sq.LastExecutionStatus != "" {
		result["lastExecutionStatus"] = sq.LastExecutionStatus
	}
	return result
}

func extractIdFromArnOrId(identifier string) string {
	parts := strings.Split(identifier, ":")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return identifier
}

// scheduledQueryTickerInterval is the interval between scheduled-query
// evaluations. TEST_MODE shortens it so integration tests observe the
// AWS first-interval contract (the first run happens one full interval
// after creation) without also waiting out the evaluation phase — the
// same convention as the scheduler and timestream-query engines.
var scheduledQueryTickerInterval = 1 * time.Minute

// recordDelivery stamps the execution outcome on the stored record: the
// consumed boundary and the execution clock. A query deleted while its
// execution was in flight has no record left to stamp, which is benign.
func recordDelivery(store *logsstore.Store, id string, boundary, executedAt int64, status string) {
	if err := store.TouchScheduledQueryDelivery(id, boundary, executedAt, status); err != nil && !errors.Is(err, logsstore.ErrResourceNotFound) {
		logs.Error("Failed to record the scheduled query execution outcome",
			logs.String("scheduledQueryId", id),
			logs.Err(err))
	}
}

func init() {
	if os.Getenv("TEST_MODE") == "true" {
		scheduledQueryTickerInterval = 1 * time.Second
	}
}

// startScheduledQueryWorker runs a background goroutine that evaluates
// schedule expressions and triggers enabled scheduled queries.
func (s *LogsService) startScheduledQueryWorker() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("PANIC in scheduled query worker, restarting",
					logs.Any("panic", r))
				go s.startScheduledQueryWorker()
			}
		}()

		ticker := time.NewTicker(scheduledQueryTickerInterval)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.tickScheduledQueries()
			}
		}
	}()
}

// scheduledQueryDue reports whether an ENABLED scheduled query should run
// at the evaluation time now. The schedule expression is evaluated in the
// query's configured timezone (UTC when unset or invalid), inside the
// optional scheduleStartTime/scheduleEndTime execution window. Each
// boundary runs exactly once: the boundary is the latest elapsed
// execution instant of the expression — rate() never runs on the creation
// boundary (the first run is one full interval after creation), cron()
// recovers a matching minute missed between evaluations, at() runs once
// its timestamp is reached — and the query runs when that boundary is
// later than the last consumed boundary (LastExecutedBoundary, zero
// meaning never run). The marker holds the boundary value, never the
// execution clock: an execution that runs late must not suppress the
// next unexecuted boundary, and it survives restarts. The evaluated
// boundary is returned so the caller stamps exactly what it consumed.
func scheduledQueryDue(sq *logsstore.ScheduledQuery, now time.Time) (time.Time, bool) {
	nowMillis := now.UnixMilli()
	if sq.ScheduleStartTime != 0 && nowMillis < sq.ScheduleStartTime {
		return time.Time{}, false
	}
	if sq.ScheduleEndTime != 0 && nowMillis > sq.ScheduleEndTime {
		return time.Time{}, false
	}
	creationTime := time.UnixMilli(sq.CreationTime).UTC()
	boundary, elapsed := scheduleexpr.ElapsedExecutionTime(sq.ScheduleExpression, now.In(scheduledQueryLocation(sq)), creationTime, nil, scheduleexpr.RateFiresAfterFirstInterval)
	if !elapsed {
		return time.Time{}, false
	}
	if !boundary.After(time.UnixMilli(sq.LastExecutedBoundary).UTC()) {
		return time.Time{}, false
	}
	return boundary, true
}

// scheduledQueryLocation resolves the timezone the schedule expression is
// evaluated in: the query's configured timezone, or UTC when unset or
// invalid.
func scheduledQueryLocation(sq *logsstore.ScheduledQuery) *time.Location {
	if sq.Timezone != "" {
		if loc, err := time.LoadLocation(sq.Timezone); err == nil {
			return loc
		}
		logs.Debug("Invalid scheduled query timezone, falling back to UTC",
			logs.String("scheduledQuery", sq.Name),
			logs.String("timezone", sq.Timezone))
	}
	return time.UTC
}

func (s *LogsService) tickScheduledQueries() {
	now := time.Now().UTC()

	s.logsStores.Range(func(key, value interface{}) bool {
		store := value.(*logsstore.Store)
		region := key.(string)

		queries, err := store.ListScheduledQueries("ENABLED")
		if err != nil {
			return true
		}

		for _, sq := range queries {
			// The boundary the evaluation consumes is the exact value
			// the execution stamps; recomputing it later against the
			// execution clock could advance past an unexecuted
			// boundary and suppress it.
			if boundary, due := scheduledQueryDue(sq, now); due {
				s.triggerScheduledQuery(region, store, sq, boundary)
			}
		}
		return true
	})
}

func (s *LogsService) triggerScheduledQuery(region string, store *logsstore.Store, sq *logsstore.ScheduledQuery, boundary time.Time) {
	now := time.Now().UTC().UnixMilli()

	exec := &logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          fmt.Sprintf("sq-%s-%d", sq.Id, now),
		TriggerTime:      now,
		Status:           "RUNNING",
	}
	if err := store.PutScheduledQueryExecution(exec); err != nil {
		logs.Error("Failed to persist scheduled query execution (RUNNING)",
			logs.String("scheduledQueryId", sq.Id),
			logs.Err(err))
		return
	}

	defer func() {
		if r := recover(); r != nil {
			exec.Status = "FAILED"
			exec.ErrorMessage = fmt.Sprintf("panic: %v", r)
			if err := store.PutScheduledQueryExecution(exec); err != nil {
				logs.Error("Failed to persist scheduled query execution (FAILED after panic)",
					logs.String("scheduledQueryId", sq.Id),
					logs.Err(err))
			}
		}
	}()

	endTime := now
	startTime := now - 60*60*1000
	if sq.StartTimeOffset > 0 {
		startTime = now - sq.StartTimeOffset
	}
	if sq.EndTimeOffset > 0 {
		endTime = now - sq.EndTimeOffset
	}

	fail := func(message string) {
		exec.Status = "FAILED"
		exec.ErrorMessage = message
		if err := store.PutScheduledQueryExecution(exec); err != nil {
			logs.Error("Failed to persist scheduled query execution (FAILED)",
				logs.String("scheduledQueryId", sq.Id),
				logs.Err(err))
		}
	}

	events := fetchLogEvents(store, sq.LogGroupIdentifiers, startTime, endTime)
	ctx := &execContext{
		startTime:     startTime,
		endTime:       endTime,
		accountID:     s.accountID,
		defaultGroups: sq.LogGroupIdentifiers,
		events:        events,
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
	rows, err := executeQueryContext(ctx, sq.QueryString)
	if err != nil {
		fail(fmt.Sprintf("query failed: %v", err))
		return
	}
	exec.Destinations = s.deliverScheduledQueryResults(region, store, sq, exec.QueryId, rows)
	for _, dest := range exec.Destinations {
		if dest.Status != destinationStatusComplete {
			exec.Status = "FAILED"
			exec.ErrorMessage = fmt.Sprintf("destination delivery failed: %s", dest.ErrorMessage)
			break
		}
	}
	if exec.Status == "FAILED" {
		if err := store.PutScheduledQueryExecution(exec); err != nil {
			logs.Error("Failed to persist scheduled query execution (delivery FAILED)",
				logs.String("scheduledQueryId", sq.Id),
				logs.Err(err))
		}
		// Record the trigger so a persistently failing destination retries
		// on schedule instead of on every worker tick.
		recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusFailed)
		return
	}

	stats := queryStats{
		recordsScanned: int64(len(ctx.events)),
	}
	for _, e := range ctx.events {
		stats.bytesScanned += int64(len(e.message))
	}
	stats.recordsMatched = stats.recordsScanned

	exec.Status = "SUCCESS"
	exec.RecordsScanned = stats.recordsScanned
	exec.RecordsMatched = stats.recordsMatched
	if err := store.PutScheduledQueryExecution(exec); err != nil {
		logs.Error("Failed to persist scheduled query execution (SUCCESS)",
			logs.String("scheduledQueryId", sq.Id),
			logs.Err(err))
	}

	recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusComplete)
}
