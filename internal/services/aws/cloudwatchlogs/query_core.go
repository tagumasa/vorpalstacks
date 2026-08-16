package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/scheduleexpr"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

func parseTagsFromParams(params map[string]interface{}) map[string]string {
	tags := make(map[string]string)
	if t, ok := params["tags"]; ok {
		if m, ok := t.(map[string]interface{}); ok {
			for k, v := range m {
				if s, ok := v.(string); ok {
					tags[k] = s
				}
			}
		}
	}
	return tags
}

// --- Scheduled Query Core ---

// CreateScheduledQueryInput holds the parsed parameters for
// CreateScheduledQuery.
type CreateScheduledQueryInput struct {
	Name                     string
	Description              string
	QueryString              string
	QueryLanguage            string
	LogGroupIdentifiers      []string
	ScheduleExpression       string
	State                    string
	ExecutionRoleArn         string
	Timezone                 string
	StartTimeOffset          int64
	EndTimeOffset            int64
	ScheduleStartTime        int64
	ScheduleEndTime          int64
	DestinationConfiguration map[string]interface{}
	Tags                     map[string]string
}

// validateScheduledQuerySpec applies the query-spec validation shared by
// CreateScheduledQuery and UpdateScheduledQuery so that an updated query can
// never hold a weaker specification than a newly created one.
func validateScheduledQuerySpec(queryLanguage, queryString, scheduleExpression, executionRoleArn string) error {
	if err := validateQueryString(queryString); err != nil {
		return err
	}
	if err := validateQueryPipeline(queryString); err != nil {
		return err
	}
	if !scheduleexpr.ValidateExpression(scheduleExpression) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid schedule expression: %s. Must be a valid rate(), cron(), or at() expression", scheduleExpression), 400)
	}
	if queryLanguage == "" {
		queryLanguage = "CWLI"
	}
	if !validateQueryLanguage(queryLanguage) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid queryLanguage: %s. Allowed values: CWLI, SQL, PPL", queryLanguage), 400)
	}
	return validateIAMRoleArn(executionRoleArn)
}

// createScheduledQueryCore validates input and persists a new scheduled query.
func (s *LogsService) createScheduledQueryCore(store *logsstore.Store, input *CreateScheduledQueryInput) (*logsstore.ScheduledQuery, error) {
	if err := validateScheduledQueryName(input.Name); err != nil {
		return nil, err
	}
	if input.QueryString == "" || input.ScheduleExpression == "" {
		return nil, ErrMissingParameter
	}
	if err := validateScheduledQuerySpec(input.QueryLanguage, input.QueryString, input.ScheduleExpression, input.ExecutionRoleArn); err != nil {
		return nil, err
	}
	if err := validateLogGroupIdentifierCount(input.LogGroupIdentifiers); err != nil {
		return nil, err
	}
	if input.DestinationConfiguration != nil {
		if err := validateDestinationConfiguration(input.DestinationConfiguration); err != nil {
			return nil, err
		}
	}

	queryLanguage := input.QueryLanguage
	if queryLanguage == "" {
		queryLanguage = "CWLI"
	}

	state := input.State
	if state == "" {
		state = "ENABLED"
	}
	if !validateScheduledQueryState(state) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid state: %s. Allowed values: ENABLED, DISABLED", state), 400)
	}

	id := fmt.Sprintf("sq-%d", time.Now().UnixNano())

	sq := &logsstore.ScheduledQuery{
		Id:                       id,
		Name:                     input.Name,
		Description:              input.Description,
		QueryString:              input.QueryString,
		QueryLanguage:            queryLanguage,
		LogGroupIdentifiers:      input.LogGroupIdentifiers,
		ScheduleExpression:       input.ScheduleExpression,
		ScheduleType:             "CUSTOMER_MANAGED",
		State:                    state,
		ExecutionRoleArn:         input.ExecutionRoleArn,
		Timezone:                 input.Timezone,
		StartTimeOffset:          input.StartTimeOffset,
		EndTimeOffset:            input.EndTimeOffset,
		ScheduleStartTime:        input.ScheduleStartTime,
		ScheduleEndTime:          input.ScheduleEndTime,
		DestinationConfiguration: input.DestinationConfiguration,
		CreationTime:             time.Now().UTC().UnixMilli(),
		Tags:                     input.Tags,
	}

	if err := store.PutScheduledQuery(sq); err != nil {
		return nil, mapStoreError(err)
	}
	return sq, nil
}

func (s *LogsService) deleteScheduledQueryCore(store *logsstore.Store, id string) error {
	if id == "" {
		return ErrMissingParameter
	}
	if err := store.DeleteScheduledQuery(id); err != nil {
		return mapStoreError(err)
	}
	return nil
}

func (s *LogsService) getScheduledQueryCore(store *logsstore.Store, id string) (*logsstore.ScheduledQuery, error) {
	if id == "" {
		return nil, ErrMissingParameter
	}
	sq, err := store.GetScheduledQuery(id)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return sq, nil
}

// UpdateScheduledQueryInput holds the parsed parameters for
// UpdateScheduledQuery. The Smithy model marks identifier, queryLanguage,
// queryString, scheduleExpression and executionRoleArn as required because
// the update replaces the full specification; the remaining members are
// replaced only when provided.
type UpdateScheduledQueryInput struct {
	Identifier               string
	Description              string
	QueryString              string
	QueryLanguage            string
	LogGroupIdentifiers      []string
	ScheduleExpression       string
	State                    string
	ExecutionRoleArn         string
	Timezone                 string
	StartTimeOffset          *int64
	EndTimeOffset            *int64
	ScheduleStartTime        *int64
	ScheduleEndTime          *int64
	DestinationConfiguration map[string]interface{}
}

// updateScheduledQueryCore validates input and applies the full-specification
// update semantics to the stored scheduled query.
func (s *LogsService) updateScheduledQueryCore(store *logsstore.Store, input *UpdateScheduledQueryInput) (*logsstore.ScheduledQuery, error) {
	if input.Identifier == "" {
		return nil, ErrMissingParameter
	}
	if input.QueryLanguage == "" || input.QueryString == "" || input.ScheduleExpression == "" || input.ExecutionRoleArn == "" {
		return nil, ErrMissingParameter
	}
	if err := validateScheduledQuerySpec(input.QueryLanguage, input.QueryString, input.ScheduleExpression, input.ExecutionRoleArn); err != nil {
		return nil, err
	}
	if input.State != "" && !validateScheduledQueryState(input.State) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid state: %s. Allowed values: ENABLED, DISABLED", input.State), 400)
	}
	if input.LogGroupIdentifiers != nil {
		if err := validateLogGroupIdentifierCount(input.LogGroupIdentifiers); err != nil {
			return nil, err
		}
	}
	if input.DestinationConfiguration != nil {
		if err := validateDestinationConfiguration(input.DestinationConfiguration); err != nil {
			return nil, err
		}
	}

	id := extractIdFromArnOrId(input.Identifier)

	sq, err := store.GetScheduledQuery(id)
	if err != nil {
		return nil, mapStoreError(err)
	}

	sq.QueryString = input.QueryString
	sq.QueryLanguage = input.QueryLanguage
	sq.ScheduleExpression = input.ScheduleExpression
	sq.ExecutionRoleArn = input.ExecutionRoleArn
	if input.Description != "" {
		sq.Description = input.Description
	}
	if input.State != "" {
		sq.State = input.State
	}
	if input.LogGroupIdentifiers != nil {
		sq.LogGroupIdentifiers = input.LogGroupIdentifiers
	}
	if input.Timezone != "" {
		sq.Timezone = input.Timezone
	}
	if input.StartTimeOffset != nil {
		sq.StartTimeOffset = *input.StartTimeOffset
	}
	if input.EndTimeOffset != nil {
		sq.EndTimeOffset = *input.EndTimeOffset
	}
	if input.ScheduleStartTime != nil {
		sq.ScheduleStartTime = *input.ScheduleStartTime
	}
	if input.ScheduleEndTime != nil {
		sq.ScheduleEndTime = *input.ScheduleEndTime
	}
	if input.DestinationConfiguration != nil {
		sq.DestinationConfiguration = input.DestinationConfiguration
	}

	if err := store.PutScheduledQuery(sq); err != nil {
		return nil, mapStoreError(err)
	}
	return sq, nil
}

// ScheduledQueryHistoryResult carries one page of scheduled query executions
// plus the scheduled query itself for response formatting.
type ScheduledQueryHistoryResult struct {
	Executions []*logsstore.ScheduledQueryExecution
	NextMarker string
	Query      *logsstore.ScheduledQuery
}

// GetScheduledQueryHistoryInput holds the parsed parameters for
// GetScheduledQueryHistory.
type GetScheduledQueryHistoryInput struct {
	Identifier        string
	StartTime         int64
	EndTime           int64
	ExecutionStatuses []string
	NextToken         string
	MaxResults        int32
}

// getScheduledQueryHistoryCore validates input and returns a page of the
// scheduled query's execution history.
func (s *LogsService) getScheduledQueryHistoryCore(store *logsstore.Store, input *GetScheduledQueryHistoryInput) (*ScheduledQueryHistoryResult, error) {
	if input.Identifier == "" {
		return nil, ErrMissingParameter
	}
	if input.MaxResults <= 0 {
		input.MaxResults = 50
	}

	id := extractIdFromArnOrId(input.Identifier)

	allExecs, err := store.ListScheduledQueryExecutions(id, input.StartTime, input.EndTime)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if len(input.ExecutionStatuses) > 0 {
		statusSet := make(map[string]bool)
		for _, st := range input.ExecutionStatuses {
			statusSet[st] = true
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

	result := pagination.PaginateSlice(allExecs, input.NextToken, int(input.MaxResults), func(e *logsstore.ScheduledQueryExecution) string {
		return strconv.FormatInt(e.TriggerTime, 10)
	})

	history := &ScheduledQueryHistoryResult{
		Executions: result.Items,
		NextMarker: result.NextMarker,
	}
	// The scheduled query lookup is best-effort for the response name; a
	// missing query does not invalidate already-recorded history entries.
	if sq, sqErr := store.GetScheduledQuery(id); sqErr == nil {
		history.Query = sq
	}
	return history, nil
}

// listScheduledQueriesCore validates input and returns a page of scheduled
// queries filtered by state and schedule type.
func (s *LogsService) listScheduledQueriesCore(store *logsstore.Store, stateFilter, scheduleTypeFilter string, maxResults int32, nextToken string) ([]*logsstore.ScheduledQuery, string, error) {
	if maxResults <= 0 {
		maxResults = 50
	}
	all, err := store.ListScheduledQueries(stateFilter)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	if scheduleTypeFilter != "" {
		var filtered []*logsstore.ScheduledQuery
		for _, sq := range all {
			if sq.ScheduleType == scheduleTypeFilter {
				filtered = append(filtered, sq)
			}
		}
		all = filtered
	}

	result := pagination.PaginateSlice(all, nextToken, int(maxResults), func(sq *logsstore.ScheduledQuery) string {
		return sq.Id
	})
	return result.Items, result.NextMarker, nil
}

// getLogGroupFieldsCore validates input and retrieves field names from
// recent log events in the specified log group.
func (s *LogsService) getLogGroupFieldsCore(store *logsstore.Store, logGroupName string, centerTime int64) ([]map[string]interface{}, error) {
	if _, err := store.GetLogGroup(logGroupName); err != nil {
		return nil, mapStoreError(err)
	}

	standardFields := []string{"@timestamp", "@message", "@logStream", "@log", "@ingestionTime"}
	jsonFields := make(map[string]bool)

	var startTime, endTime int64
	if centerTime > 0 {
		startTime = centerTime - 8*60*1000
		endTime = centerTime + 8*60*1000
	}
	streams, _, _ := store.ListLogStreams(logGroupName, "", "", 50)
	for _, ls := range streams {
		events, _, _, _ := store.GetLogEvents(logGroupName, ls.Name, startTime, endTime, 50, true, "")
		for _, evt := range events {
			var data map[string]interface{}
			if json.Unmarshal([]byte(evt.Message), &data) == nil {
				for k := range data {
					jsonFields[k] = true
				}
			}
		}
	}

	fields := make([]map[string]interface{}, 0, len(standardFields)+len(jsonFields))
	for _, f := range standardFields {
		fields = append(fields, map[string]interface{}{
			"name":    f,
			"percent": 100,
		})
	}

	jsonFieldNames := make([]string, 0, len(jsonFields))
	for f := range jsonFields {
		jsonFieldNames = append(jsonFieldNames, f)
	}
	sort.Strings(jsonFieldNames)
	for _, f := range jsonFieldNames {
		fields = append(fields, map[string]interface{}{
			"name":    f,
			"percent": 50,
		})
	}

	return fields, nil
}

// getLogFieldsCore validates input and retrieves field names from log events.
func (s *LogsService) getLogFieldsCore(store *logsstore.Store, dataSourceName, dataSourceType string) ([]string, error) {
	jsonFields := make(map[string]bool)
	streams, _, _ := store.ListLogStreams(dataSourceName, "", "", 20)
	for _, ls := range streams {
		events, _, _, _ := store.GetLogEvents(dataSourceName, ls.Name, 0, 0, 20, true, "")
		for _, evt := range events {
			var data map[string]interface{}
			if json.Unmarshal([]byte(evt.Message), &data) == nil {
				for k := range data {
					jsonFields[k] = true
				}
			}
		}
	}

	fieldNames := []string{"@timestamp", "@message", "@logStream"}
	for f := range jsonFields {
		fieldNames = append(fieldNames, f)
	}
	return fieldNames, nil
}
