package cloudwatchlogs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"sort"

	"vorpalstacks/internal/common/request"
)

// GetLogGroupFields returns a list of fields found in the log events of
// the specified log group. For JSON-structured logs, field names are
// extracted from the JSON payload. Standard fields (@timestamp, @message,
// @logStream, @logGroup) are always included.
func (s *LogsService) GetLogGroupFields(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		logGroupName = request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	}
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	centerTime := int64(request.GetIntParam(req.Parameters, "Time"))

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetLogGroup(logGroupName)
	if err != nil {
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

	return map[string]interface{}{
		"logGroupFields": fields,
	}, nil
}

// GetLogRecord retrieves a single log record by its pointer.
// The pointer is a base64-encoded string in the format:
// logGroup|logStream|timestamp|message
func (s *LogsService) GetLogRecord(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logRecordPointer := request.GetParamLowerFirst(req.Parameters, "LogRecordPointer")
	if logRecordPointer == "" {
		return nil, ErrMissingParameter
	}

	unmask := request.GetBoolParam(req.Parameters, "Unmask")
	_ = unmask

	decoded, err := base64.StdEncoding.DecodeString(logRecordPointer)
	if err != nil {
		return nil, NewLogsError("InvalidParameterException",
			"Invalid logRecordPointer", 400)
	}

	parts := splitPointer(string(decoded))
	if len(parts) < 4 {
		return nil, NewLogsError("InvalidParameterException",
			"Invalid logRecordPointer format", 400)
	}

	logGroupName := parts[0]
	logStreamName := parts[1]
	timestamp := parts[2]
	message := parts[3]

	record := map[string]interface{}{
		"@timestamp": timestamp,
		"@message":   message,
		"@logStream": logStreamName,
		"@logGroup":  logGroupName,
	}

	var jsonData map[string]interface{}
	if json.Unmarshal([]byte(message), &jsonData) == nil {
		for k, v := range jsonData {
			record[k] = v
		}
	}

	return map[string]interface{}{
		"logRecord": record,
	}, nil
}

// GetLogObject retrieves a log object by its pointer.
func (s *LogsService) GetLogObject(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logObjectPointer := request.GetParamLowerFirst(req.Parameters, "LogObjectPointer")
	if logObjectPointer == "" {
		return nil, ErrMissingParameter
	}

	unmask := request.GetBoolParam(req.Parameters, "Unmask")
	_ = unmask

	decoded, err := base64.StdEncoding.DecodeString(logObjectPointer)
	if err != nil {
		return nil, NewLogsError("InvalidParameterException",
			"Invalid logObjectPointer", 400)
	}

	parts := splitPointer(string(decoded))
	if len(parts) < 4 {
		return nil, NewLogsError("InvalidParameterException",
			"Invalid logObjectPointer format", 400)
	}

	record := map[string]interface{}{
		"@timestamp": parts[2],
		"@message":   parts[3],
		"@logStream": parts[1],
		"@logGroup":  parts[0],
	}

	var jsonData map[string]interface{}
	if json.Unmarshal([]byte(parts[3]), &jsonData) == nil {
		for k, v := range jsonData {
			record[k] = v
		}
	}

	return record, nil
}

// GetLogFields retrieves the fields available for a specific data source.
func (s *LogsService) GetLogFields(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	dataSourceName := request.GetParamLowerFirst(req.Parameters, "DataSourceName")
	dataSourceType := request.GetParamLowerFirst(req.Parameters, "DataSourceType")

	if dataSourceName == "" || dataSourceType == "" {
		return nil, ErrMissingParameter
	}

	if dataSourceType != "AWS::Logs::LogGroup" {
		return map[string]interface{}{
			"logFields": []interface{}{},
		}, nil
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

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
	sort.Strings(fieldNames)

	fields := make([]map[string]interface{}, len(fieldNames))
	for i, f := range fieldNames {
		fields[i] = map[string]interface{}{"name": f}
	}

	return map[string]interface{}{
		"logFields": fields,
	}, nil
}

func splitPointer(s string) []string {
	var parts []string
	cur := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '|' {
			parts = append(parts, s[cur:i])
			cur = i + 1
		}
	}
	if cur < len(s) {
		parts = append(parts, s[cur:])
	}
	return parts
}
