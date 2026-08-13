package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
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

func (s *LogsService) createScheduledQueryCore(store *logsstore.Store, region, accountID string, req *request.ParsedRequest) (*logsstore.ScheduledQuery, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if err := validateScheduledQueryName(name); err != nil {
		return nil, err
	}

	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	if err := validateQueryString(queryString); err != nil {
		return nil, err
	}

	scheduleExpression := request.GetParamLowerFirst(req.Parameters, "ScheduleExpression")
	if scheduleExpression == "" {
		return nil, ErrMissingParameter
	}
	if !scheduleexpr.ValidateExpression(scheduleExpression) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid schedule expression: %s", scheduleExpression), 400)
	}

	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")
	if queryLanguage == "" {
		queryLanguage = "CWLI"
	}
	if !validateQueryLanguage(queryLanguage) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid queryLanguage: %s. Allowed values: CWLI, SQL, PPL", queryLanguage), 400)
	}

	state := request.GetParamLowerFirst(req.Parameters, "State")
	if state == "" {
		state = "ENABLED"
	}
	if !validateScheduledQueryState(state) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid state: %s. Allowed values: ENABLED, DISABLED", state), 400)
	}

	tags := parseTagsFromParams(req.Parameters)

	sq := &logsstore.ScheduledQuery{
		Name:               name,
		Description:        request.GetParamLowerFirst(req.Parameters, "Description"),
		QueryString:        queryString,
		QueryLanguage:      queryLanguage,
		ScheduleExpression: scheduleExpression,
		State:              state,
		CreationTime:       time.Now().UTC().UnixMilli(),
		Tags:               tags,
	}

	if v := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifiers"); v != "" {
		sq.LogGroupIdentifiers = []string{v}
	}
	if strs := request.GetStringList(req.Parameters, "LogGroupIdentifiers"); len(strs) > 0 {
		sq.LogGroupIdentifiers = strs
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

func (s *LogsService) listScheduledQueriesCore(store *logsstore.Store, stateFilter string, maxResults int32, nextToken string) ([]*logsstore.ScheduledQuery, string, error) {
	if maxResults <= 0 {
		maxResults = 50
	}
	all, err := store.ListScheduledQueries(stateFilter)
	if err != nil {
		return nil, "", mapStoreError(err)
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

	standardFields := []string{"@timestamp", "@message", "@logStream", "@logGroup", "@ingestionTime"}
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
