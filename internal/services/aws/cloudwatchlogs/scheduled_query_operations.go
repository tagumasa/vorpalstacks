package cloudwatchlogs

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CreateScheduledQuery creates a scheduled CloudWatch Logs Insights query.
func (s *LogsService) CreateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifiers := request.GetStringList(req.Parameters, "LogGroupIdentifiers")
	// An explicitly empty array is a present member, not an omitted one:
	// the non-nil empty slice carries that distinction to the Core, where
	// the list's @length(min 1) trait rejects it.
	if logGroupIdentifiers == nil && request.HasListParam(req.Parameters, "LogGroupIdentifiers") {
		logGroupIdentifiers = []string{}
	}
	destinationConfiguration := request.GetMapParamLowerFirst(req.Parameters, "DestinationConfiguration")

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

	arn := svcarn.NewARNBuilder(s.accountID, reqCtx.GetRegion()).CloudWatch().ScheduledQuery(sq.Id)
	return map[string]interface{}{
		"scheduledQueryArn": arn,
		"state":             sq.State,
	}, nil
}

// resolveScheduledQueryRef resolves the documented ARN-or-name Identifier
// form of the scheduled-query operations to the record's id key: the ARN
// form through the ARN parser, the bare name form through a record scan —
// the console and CLI send the name, and the records are keyed by the
// opaque id alone. An unmatched name returns the input unchanged so the
// record lookup reports it as not found.
func resolveScheduledQueryRef(store *logsstore.Store, identifier string) (string, error) {
	if strings.HasPrefix(identifier, "arn:") {
		return resolveScheduledQueryIdentifier(identifier), nil
	}
	id, err := store.ScheduledQueryIdByName(identifier)
	if err != nil {
		return "", err
	}
	if id != "" {
		return id, nil
	}
	return identifier, nil
}

// DeleteScheduledQuery deletes a scheduled query.
func (s *LogsService) DeleteScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteScheduledQueryCore(store, request.GetParamLowerFirst(req.Parameters, "Identifier")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateScheduledQuery updates a scheduled query.
func (s *LogsService) UpdateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	logGroupIdentifiers := request.GetStringList(req.Parameters, "LogGroupIdentifiers")
	// An explicitly empty array is a present member, not an omitted one:
	// the non-nil empty slice carries that distinction to the Core, where
	// the list's @length(min 1) trait rejects it.
	if logGroupIdentifiers == nil && request.HasListParam(req.Parameters, "LogGroupIdentifiers") {
		logGroupIdentifiers = []string{}
	}
	destinationConfiguration := request.GetMapParamLowerFirst(req.Parameters, "DestinationConfiguration")
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
		return nil, errValidationMember("identifier")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sq, err := s.getScheduledQueryCore(store, identifier)
	if err != nil {
		return nil, err
	}

	return formatScheduledQuery(sq, reqCtx.GetRegion(), s.accountID), nil
}

// GetScheduledQueryHistory retrieves execution history for a scheduled query.
func (s *LogsService) GetScheduledQueryHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	executionStatuses := request.GetStringList(req.Parameters, "ExecutionStatuses")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Both window members are required; the presence flags distinguish an
	// explicit zero (the members' documented minimum) from an omitted
	// member, which an int64 read alone cannot.
	startTimeSet := false
	for _, key := range []string{"StartTime", "startTime"} {
		if _, ok := req.Parameters[key]; ok {
			startTimeSet = true
			break
		}
	}
	endTimeSet := false
	for _, key := range []string{"EndTime", "endTime"} {
		if _, ok := req.Parameters[key]; ok {
			endTimeSet = true
			break
		}
	}

	history, err := s.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
		Identifier:        request.GetParamLowerFirst(req.Parameters, "Identifier"),
		StartTime:         int64(request.GetIntParam(req.Parameters, "StartTime")),
		EndTime:           int64(request.GetIntParam(req.Parameters, "EndTime")),
		StartTimeSet:      startTimeSet,
		EndTimeSet:        endTimeSet,
		ExecutionStatuses: executionStatuses,
		NextToken:         request.GetParamLowerFirst(req.Parameters, "NextToken"),
		MaxResults:        int32(request.GetIntParam(req.Parameters, "MaxResults")),
	})
	if err != nil {
		return nil, err
	}

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

	arn := svcarn.NewARNBuilder(s.accountID, reqCtx.GetRegion()).CloudWatch().ScheduledQuery(history.Id)

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

// mapExecutionStatus renders the internal execution vocabulary onto the
// ExecutionStatus enum the history records and lastExecutionStatus carry.
// The enum defines no Cancelled member while StopQuery can cancel an
// individual scheduled query execution: the history plane reports such a
// run as Failed — the one in-enum terminal state whose errorMessage member
// can say why — and the query plane's QueryStatus enum (which does define
// Cancelled) keeps reporting it as Cancelled.
func mapExecutionStatus(internalStatus string) string {
	switch internalStatus {
	case logsstore.ScheduledExecutionStatusRunning:
		return "Running"
	case logsstore.ScheduledExecutionStatusSuccess:
		return logsstore.ScheduledQueryStatusComplete
	case logsstore.ScheduledExecutionStatusTimeout:
		return logsstore.ScheduledQueryStatusTimeout
	case logsstore.ScheduledExecutionStatusFailed, logsstore.ScheduledExecutionStatusCancelled:
		return logsstore.ScheduledQueryStatusFailed
	default:
		return internalStatus
	}
}

// ListScheduledQueries lists scheduled queries.
func (s *LogsService) ListScheduledQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.listScheduledQueriesCore(store,
		request.GetParamLowerFirst(req.Parameters, "State"),
		request.GetParamLowerFirst(req.Parameters, "ScheduleType"),
		int32(request.GetIntParam(req.Parameters, "MaxResults")),
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
	arn := svcarn.NewARNBuilder(accountID, region).CloudWatch().ScheduledQuery(sq.Id)
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
	arn := svcarn.NewARNBuilder(accountID, region).CloudWatch().ScheduledQuery(sq.Id)
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
