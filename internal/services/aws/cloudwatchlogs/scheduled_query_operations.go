package cloudwatchlogs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// CreateScheduledQuery creates a scheduled CloudWatch Logs Insights query.
func (s *LogsService) CreateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	scheduleExpression := request.GetParamLowerFirst(req.Parameters, "ScheduleExpression")

	if name == "" || queryString == "" || scheduleExpression == "" {
		return nil, ErrMissingParameter
	}

	id := fmt.Sprintf("sq-%d", time.Now().UnixNano())
	region := reqCtx.GetRegion()

	var logGroupIdentifiers []string
	if idents, ok := req.Parameters["logGroupIdentifiers"]; ok {
		if arr, ok := idents.([]interface{}); ok {
			for _, item := range arr {
				if s, ok := item.(string); ok {
					logGroupIdentifiers = append(logGroupIdentifiers, s)
				}
			}
		}
	}

	description := request.GetParamLowerFirst(req.Parameters, "Description")
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")
	if queryLanguage == "" {
		queryLanguage = "CWLI"
	}
	executionRoleArn := request.GetParamLowerFirst(req.Parameters, "ExecutionRoleArn")
	timezone := request.GetParamLowerFirst(req.Parameters, "Timezone")
	state := request.GetParamLowerFirst(req.Parameters, "State")
	if state == "" {
		state = "ENABLED"
	}
	startTimeOffset := int64(request.GetIntParam(req.Parameters, "StartTimeOffset"))
	endTimeOffset := int64(request.GetIntParam(req.Parameters, "EndTimeOffset"))
	scheduleStartTime := int64(request.GetIntParam(req.Parameters, "ScheduleStartTime"))
	scheduleEndTime := int64(request.GetIntParam(req.Parameters, "ScheduleEndTime"))

	var destinationConfiguration map[string]interface{}
	if dc, ok := req.Parameters["destinationConfiguration"]; ok {
		if m, ok := dc.(map[string]interface{}); ok {
			destinationConfiguration = m
		}
	}

	var tags map[string]string
	if t, ok := req.Parameters["tags"]; ok {
		if m, ok := t.(map[string]interface{}); ok {
			tags = make(map[string]string)
			for k, v := range m {
				if s, ok := v.(string); ok {
					tags[k] = s
				}
			}
		}
	}

	sq := &logsstore.ScheduledQuery{
		Id:                       id,
		Name:                     name,
		Description:              description,
		QueryString:              queryString,
		QueryLanguage:            queryLanguage,
		LogGroupIdentifiers:      logGroupIdentifiers,
		ScheduleExpression:       scheduleExpression,
		ScheduleType:             "CUSTOMER_MANAGED",
		State:                    state,
		ExecutionRoleArn:         executionRoleArn,
		Timezone:                 timezone,
		StartTimeOffset:          startTimeOffset,
		EndTimeOffset:            endTimeOffset,
		ScheduleStartTime:        scheduleStartTime,
		ScheduleEndTime:          scheduleEndTime,
		DestinationConfiguration: destinationConfiguration,
		CreationTime:             time.Now().UTC().UnixMilli(),
		Tags:                     tags,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.PutScheduledQuery(sq); err != nil {
		return nil, mapStoreError(err)
	}

	arn := fmt.Sprintf("arn:aws:logs:%s:%s:scheduled-query:%s", region, s.accountID, id)
	return map[string]interface{}{
		"scheduledQueryArn": arn,
		"state":             state,
	}, nil
}

// DeleteScheduledQuery deletes a scheduled query.
func (s *LogsService) DeleteScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "Identifier")
	if identifier == "" {
		return nil, ErrMissingParameter
	}

	id := extractIdFromArnOrId(identifier)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteScheduledQuery(id); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// UpdateScheduledQuery updates a scheduled query.
func (s *LogsService) UpdateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "Identifier")
	if identifier == "" {
		return nil, ErrMissingParameter
	}

	id := extractIdFromArnOrId(identifier)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sq, err := store.GetScheduledQuery(id)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if v := request.GetParamLowerFirst(req.Parameters, "Description"); v != "" {
		sq.Description = v
	}
	if v := request.GetParamLowerFirst(req.Parameters, "QueryString"); v != "" {
		sq.QueryString = v
	}
	if v := request.GetParamLowerFirst(req.Parameters, "ScheduleExpression"); v != "" {
		sq.ScheduleExpression = v
	}
	if v := request.GetParamLowerFirst(req.Parameters, "State"); v != "" {
		sq.State = v
	}
	if v := request.GetParamLowerFirst(req.Parameters, "ExecutionRoleArn"); v != "" {
		sq.ExecutionRoleArn = v
	}

	if err := store.PutScheduledQuery(sq); err != nil {
		return nil, mapStoreError(err)
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

	sq, err := store.GetScheduledQuery(id)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return formatScheduledQuery(sq, reqCtx.GetRegion(), s.accountID), nil
}

// GetScheduledQueryHistory retrieves execution history for a scheduled query.
func (s *LogsService) GetScheduledQueryHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "Identifier")
	if identifier == "" {
		return nil, ErrMissingParameter
	}

	id := extractIdFromArnOrId(identifier)

	startTime := int64(request.GetIntParam(req.Parameters, "StartTime"))
	endTime := int64(request.GetIntParam(req.Parameters, "EndTime"))

	var executionStatuses []string
	if es, ok := req.Parameters["executionStatuses"]; ok {
		if arr, ok := es.([]interface{}); ok {
			for _, item := range arr {
				if s, ok := item.(string); ok {
					executionStatuses = append(executionStatuses, s)
				}
			}
		}
	}

	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))
	if maxResults <= 0 {
		maxResults = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allExecs, err := store.ListScheduledQueryExecutions(id, startTime, endTime)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if len(executionStatuses) > 0 {
		statusSet := make(map[string]bool)
		for _, s := range executionStatuses {
			statusSet[s] = true
		}
		var filtered []*logsstore.ScheduledQueryExecution
		for _, exec := range allExecs {
			mappedStatus := mapExecutionStatus(exec.Status)
			if statusSet[mappedStatus] {
				filtered = append(filtered, exec)
			}
		}
		allExecs = filtered
	}

	result := pagination.PaginateSlice(allExecs, nextToken, int(maxResults), func(e *logsstore.ScheduledQueryExecution) string {
		return strconv.FormatInt(e.TriggerTime, 10)
	})

	history := make([]map[string]interface{}, len(result.Items))
	for i, exec := range result.Items {
		entry := map[string]interface{}{
			"queryId":            exec.QueryId,
			"triggeredTimestamp": exec.TriggerTime,
			"executionStatus":    mapExecutionStatus(exec.Status),
		}
		if exec.ErrorMessage != "" {
			entry["errorMessage"] = exec.ErrorMessage
		}
		entry["destinations"] = []interface{}{}
		history[i] = entry
	}

	arn := fmt.Sprintf("arn:aws:logs:%s:%s:scheduled-query:%s", reqCtx.GetRegion(), s.accountID, id)

	resp := map[string]interface{}{
		"scheduledQueryArn": arn,
		"triggerHistory":    history,
	}

	sq, sqErr := store.GetScheduledQuery(id)
	if sqErr == nil && sq != nil {
		resp["name"] = sq.Name
	}

	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
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
	state := request.GetParamLowerFirst(req.Parameters, "State")
	scheduleTypeFilter := request.GetParamLowerFirst(req.Parameters, "ScheduleType")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))
	if maxResults <= 0 {
		maxResults = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allQueries, err := store.ListScheduledQueries(state)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if scheduleTypeFilter != "" {
		var filtered []*logsstore.ScheduledQuery
		for _, sq := range allQueries {
			if sq.ScheduleType == scheduleTypeFilter {
				filtered = append(filtered, sq)
			}
		}
		allQueries = filtered
	}

	result := pagination.PaginateSlice(allQueries, nextToken, int(maxResults), func(sq *logsstore.ScheduledQuery) string {
		return sq.Id
	})

	region := reqCtx.GetRegion()
	queries := make([]map[string]interface{}, len(result.Items))
	for i, sq := range result.Items {
		queries[i] = formatScheduledQuerySummary(sq, region, s.accountID)
	}

	resp := map[string]interface{}{
		"scheduledQueries": queries,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
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
	if sq.LastExecutionStatus != nil {
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
	if sq.LastExecutionStatus != nil {
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

		ticker := time.NewTicker(1 * time.Minute)
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

func (s *LogsService) tickScheduledQueries() {
	s.logsStores.Range(func(key, value interface{}) bool {
		store := value.(*logsstore.Store)
		region := key.(string)

		queries, err := store.ListScheduledQueries("ENABLED")
		if err != nil {
			return true
		}

		now := time.Now().UTC().UnixMilli()
		for _, sq := range queries {
			if sq.LastTriggeredTime == 0 {
				s.triggerScheduledQuery(region, store, sq)
				continue
			}

			nextTime := computeNextTriggerTime(sq.ScheduleExpression, sq.LastTriggeredTime)
			if nextTime > 0 && nextTime <= now {
				s.triggerScheduledQuery(region, store, sq)
			}
		}
		return true
	})
}

func (s *LogsService) triggerScheduledQuery(region string, store *logsstore.Store, sq *logsstore.ScheduledQuery) {
	now := time.Now().UTC().UnixMilli()

	exec := &logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          fmt.Sprintf("sq-%s-%d", sq.Id, now),
		TriggerTime:      now,
		Status:           "RUNNING",
	}
	_ = store.PutScheduledQueryExecution(exec)

	defer func() {
		if r := recover(); r != nil {
			exec.Status = "FAILED"
			exec.ErrorMessage = fmt.Sprintf("panic: %v", r)
			_ = store.PutScheduledQueryExecution(exec)
		}
	}()

	endTime := now
	startTime := now - 60*60*1000

	var allEvents []logEventWithContext
	for _, lgName := range sq.LogGroupIdentifiers {
		streams, _, _ := store.ListLogStreams(lgName, "", "", 1000)
		for _, ls := range streams {
			events, _, _, _ := store.GetLogEvents(lgName, ls.Name, startTime, endTime, 10000, true, "")
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

	_, stats := executeQuery(sq.QueryString, allEvents)

	exec.Status = "SUCCESS"
	exec.RecordsScanned = stats.recordsScanned
	exec.RecordsMatched = stats.recordsMatched
	_ = store.PutScheduledQueryExecution(exec)

	sq.LastTriggeredTime = now
	_ = store.PutScheduledQuery(sq)
}

func computeNextTriggerTime(scheduleExpression string, lastTriggered int64) int64 {
	expr := strings.TrimSpace(scheduleExpression)
	exprLower := strings.ToLower(expr)

	if strings.HasPrefix(exprLower, "rate(") {
		closeParen := strings.LastIndex(expr, ")")
		if closeParen < 0 {
			return 0
		}
		inner := strings.TrimSpace(expr[5:closeParen])

		parts := strings.Fields(inner)
		if len(parts) != 2 {
			return 0
		}

		n, err := strconv.Atoi(parts[0])
		if err != nil || n <= 0 {
			return 0
		}

		unit := strings.ToLower(parts[1])
		var durationMs int64
		switch {
		case strings.HasPrefix(unit, "minute"):
			durationMs = int64(n) * 60 * 1000
		case strings.HasPrefix(unit, "hour"):
			durationMs = int64(n) * 60 * 60 * 1000
		case strings.HasPrefix(unit, "day"):
			durationMs = int64(n) * 24 * 60 * 60 * 1000
		default:
			return 0
		}

		return lastTriggered + durationMs
	}

	return 0
}
