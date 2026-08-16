package cloudwatchlogs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"sort"
	"strings"

	"vorpalstacks/internal/common/request"
)

// GetLogGroupFields returns a list of fields found in the log events of
// the specified log group. For JSON-structured logs, field names are
// extracted from the JSON payload. Standard fields (@timestamp, @message,
// @logStream, @log, @ingestionTime) are always included.
func (s *LogsService) GetLogGroupFields(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		logGroupName = request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	}
	if err := validateLogGroupName(logGroupName); err != nil {
		return nil, err
	}

	centerTime := int64(request.GetIntParam(req.Parameters, "Time"))

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	fields, err := s.getLogGroupFieldsCore(store, logGroupName, centerTime)
	if err != nil {
		return nil, err
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

	// The Unmask parameter controls whether data-protection-masked fields
	// are returned in their original form. Data protection masking is not
	// implemented on this platform; the parameter is accepted but has no
	// effect on the returned data.
	_ = request.GetBoolParam(req.Parameters, "Unmask")

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
		"@log":       reqCtx.GetAccountID() + ":" + logGroupName,
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

	// Data protection masking is not implemented on this platform;
	// the Unmask parameter is accepted but has no effect.
	_ = request.GetBoolParam(req.Parameters, "Unmask")

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
		"@log":       reqCtx.GetAccountID() + ":" + parts[0],
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

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	fieldNames, err := s.getLogFieldsCore(store, dataSourceName, dataSourceType)
	if err != nil {
		return nil, err
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

// splitPointer splits a decoded log record pointer into its
// logGroup|logStream|timestamp|message parts. Log group and stream names
// cannot contain the delimiter, so the first three splits are unambiguous;
// the message is everything after the third delimiter, including any '|'
// characters it contains.
func splitPointer(s string) []string {
	return strings.SplitN(s, "|", 4)
}
