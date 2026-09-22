package cloudwatchlogs

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/eventstream"
)

// GetLogGroupFields returns a list of fields found in the log events of
// the specified log group. For JSON-structured logs, field names are
// extracted from the JSON payload. Standard fields (@timestamp, @message,
// @logStream, @log, @ingestionTime) are always included.
func (s *LogsService) GetLogGroupFields(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName, err := logGroupNameOrIdentifier(req.Parameters)
	if err != nil {
		return nil, err
	}
	if err := validateLogGroupName(logGroupName); err != nil {
		return nil, err
	}

	// The time member is epoch seconds per the model documentation; the
	// core centres its ±8-minute window on it and falls back to the most
	// recent 15 minutes when it is omitted.
	timeSeconds := int64(request.GetIntParam(req.Parameters, "Time"))

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	fields, err := s.getLogGroupFieldsCore(store, logGroupName, timeSeconds)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"logGroupFields": fields,
	}, nil
}

// resolveEventPointerCore validates the record a decoded pointer
// addresses against the store and returns its field map: the group, the
// stream and the event itself (an identical message at the pointer's
// timestamp) must all exist. GetLogRecord and GetLogObject share the
// resolution — a well-formed pointer to deleted or fabricated content
// answers the operations' declared ResourceNotFoundException instead of a
// reconstructed 200 record.
func (s *LogsService) resolveEventPointerCore(store *logsstore.Store, accountID, logGroupName, logStreamName, timestamp, message string) (map[string]interface{}, error) {
	if _, err := store.GetLogGroup(logGroupName); err != nil {
		if errors.Is(err, logsstore.ErrLogGroupNotFound) {
			return nil, NewLogsError("ResourceNotFoundException",
				"The specified log group does not exist", 400)
		}
		return nil, mapStoreError(err)
	}
	if _, err := store.GetLogStream(logGroupName, logStreamName); err != nil {
		if errors.Is(err, logsstore.ErrLogStreamNotFound) {
			return nil, NewLogsError("ResourceNotFoundException",
				"The specified log stream does not exist", 400)
		}
		return nil, mapStoreError(err)
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, NewLogsError("InvalidParameterException",
			"Invalid log record pointer: the timestamp component is not an epoch-millisecond value", 400)
	}
	events, err := fetchAllLogEventsEndInclusive(store, logGroupName, logStreamName, ts, ts)
	if err != nil {
		return nil, mapStoreError(err)
	}
	found := false
	for _, e := range events {
		if e.Message == message {
			found = true
			break
		}
	}
	if !found {
		return nil, NewLogsError("ResourceNotFoundException",
			"The specified log record does not exist", 400)
	}

	record := map[string]interface{}{
		"@timestamp": timestamp,
		"@message":   message,
		"@logStream": logStreamName,
		"@log":       accountID + ":" + logGroupName,
	}
	var jsonData map[string]interface{}
	if json.Unmarshal([]byte(message), &jsonData) == nil {
		for k, v := range jsonData {
			record[k] = v
		}
	}
	return record, nil
}

// GetLogRecord retrieves a single log record by its pointer.
// The pointer is a base64-encoded string in the format:
// logGroup|logStream|timestamp|message, with the delimiter escaped inside
// the first two fields (see splitPointer).
func (s *LogsService) GetLogRecord(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logRecordPointer := request.GetParamLowerFirst(req.Parameters, "LogRecordPointer")
	if logRecordPointer == "" {
		return nil, errRequiredMember("logRecordPointer")
	}

	// The Unmask parameter controls whether data-protection-masked fields
	// are returned in their original form. Data protection masking is not
	// implemented on this platform; the parameter is accepted but has no
	// effect on the returned data.
	_ = request.GetBoolParam(req.Parameters, "Unmask")

	parts, err := decodeEventPointer(logRecordPointer)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	record, err := s.resolveEventPointerCore(store, reqCtx.GetAccountID(), parts[0], parts[1], parts[2], parts[3])
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"logRecord": record,
	}, nil
}

// logObjectEventStream is the streamed response of GetLogObject: the
// finite event-stream reader the dispatcher serves with the eventstream
// content type.
type logObjectEventStream struct {
	reader io.Reader
}

// GetStream returns the event-stream reader.
func (e *logObjectEventStream) GetStream() io.Reader {
	return e.reader
}

// GetStreamHeaders returns the HTTP headers of the event stream.
func (e *logObjectEventStream) GetStreamHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/vnd.amazon.eventstream")
	return headers
}

// GetLogObject retrieves a log object by its pointer. The response's one
// modelled member is fieldStream — the @streaming union
// GetLogObjectResponseStream — so the wire form is one awsJson1.1 event
// stream: an initial-response message, then one fields event whose
// payload is the FieldsData document (the data blob in its base64 JSON
// form carries "the fields and values of the log event in a structured
// format that can be parsed and processed by the client"), then the
// stream ends.
func (s *LogsService) GetLogObject(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logObjectPointer := request.GetParamLowerFirst(req.Parameters, "LogObjectPointer")
	if logObjectPointer == "" {
		return nil, errRequiredMember("logObjectPointer")
	}

	// Data protection masking is not implemented on this platform;
	// the Unmask parameter is accepted but has no effect.
	_ = request.GetBoolParam(req.Parameters, "Unmask")

	parts, err := decodeEventPointer(logObjectPointer)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	record, err := s.resolveEventPointerCore(store, reqCtx.GetAccountID(), parts[0], parts[1], parts[2], parts[3])
	if err != nil {
		return nil, err
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return nil, NewLogsError("InternalServerException", "Failed to serialise the log object", 500)
	}
	// The FieldsData document: the data blob is the record's JSON in the
	// base64 form a blob carries on the awsJson1.1 wire.
	payload, err := json.Marshal(map[string]interface{}{
		"data": base64.StdEncoding.EncodeToString(recordJSON),
	})
	if err != nil {
		return nil, NewLogsError("InternalServerException", "Failed to serialise the log object", 500)
	}

	var buf bytes.Buffer
	encoder := eventstream.NewEncoder(&buf)
	if err := encoder.WriteInitialResponse([]byte("{}")); err != nil {
		return nil, NewLogsError("InternalServerException", "Failed to encode the log object stream", 500)
	}
	if err := encoder.WriteEvent("fields", "application/json", payload); err != nil {
		return nil, NewLogsError("InternalServerException", "Failed to encode the log object stream", 500)
	}
	return &logObjectEventStream{reader: &buf}, nil
}

// GetLogFields retrieves the fields available for a specific data source.
// Each item carries the modelled logFieldName member.
func (s *LogsService) GetLogFields(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	dataSourceName := request.GetParamLowerFirst(req.Parameters, "DataSourceName")
	dataSourceType := request.GetParamLowerFirst(req.Parameters, "DataSourceType")

	if dataSourceName == "" || dataSourceType == "" {
		return nil, errRequiredMember("dataSourceName and dataSourceType")
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
		fields[i] = map[string]interface{}{"logFieldName": f}
	}

	return map[string]interface{}{
		"logFields": fields,
	}, nil
}

// decodeEventPointer base64-decodes a log record pointer and splits it
// into its logGroup|logStream|timestamp|message parts, un-escaping the
// delimiter encoding the first two fields carry (see escapePointerField).
// The message is the last field and stays raw — a '|' inside it belongs
// to the message itself.
func decodeEventPointer(encoded string) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, NewLogsError("InvalidParameterException",
			"Invalid log record pointer", 400)
	}
	parts := splitPointer(string(decoded))
	if len(parts) < 4 {
		return nil, NewLogsError("InvalidParameterException",
			"Invalid log record pointer format", 400)
	}
	parts[0] = unescapePointerField(parts[0])
	parts[1] = unescapePointerField(parts[1])
	return parts, nil
}

// splitPointer splits a decoded log record pointer on its first three
// unescaped delimiters. The group and stream fields escape the delimiter
// and the escape character itself, so a log stream name containing '|'
// cannot shift the subsequent fields; the message is everything after the
// third delimiter, including any '|' characters it contains.
func splitPointer(s string) []string {
	var parts []string
	start, field := 0, 0
	for i := 0; i < len(s) && field < 3; i++ {
		if s[i] == '\\' && i+1 < len(s) {
			i++
			continue
		}
		if s[i] == '|' {
			parts = append(parts, s[start:i])
			start = i + 1
			field++
		}
	}
	if field < 3 {
		return strings.Split(s, "|")
	}
	return append(parts, s[start:])
}

// unescapePointerField reverses escapePointerField.
func unescapePointerField(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && (s[i+1] == '|' || s[i+1] == '\\') {
			b.WriteByte(s[i+1])
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
