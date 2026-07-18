package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/filterpattern"
)

const logEventIDSize = 16

// jan012024Millis is Jan 1, 2024 00:00:00 UTC in epoch milliseconds. AWS
// requires startTime on or after this date when startFromHead=false.
const jan012024Millis = 1704067200000

func logEventToResponse(e *logsstore.OutputLogEvent) map[string]interface{} {
	resp := map[string]interface{}{
		"timestamp":     e.Timestamp,
		"message":       e.Message,
		"ingestionTime": e.IngestionTime,
	}
	if e.LogStreamName != "" {
		resp["logStreamName"] = e.LogStreamName
	}
	return resp
}

// compressJSON serialises v as JSON and gzip-compresses the result.
func compressJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if err := json.NewEncoder(gw).Encode(v); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// --- Log stream CRUD ---

// CreateLogStream creates a new CloudWatch Logs log stream.
func (s *LogsService) CreateLogStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	logStreamName := request.GetParamLowerFirst(req.Parameters, "LogStreamName")

	if logGroupName == "" || logStreamName == "" {
		return nil, ErrMissingParameter
	}

	ls := logsstore.NewLogStream(logStreamName, logGroupName)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.CreateLogStream(ls); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DeleteLogStream deletes a CloudWatch Logs log stream.
func (s *LogsService) DeleteLogStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	logStreamName := request.GetParamLowerFirst(req.Parameters, "LogStreamName")

	if logGroupName == "" || logStreamName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteLogStream(logGroupName, logStreamName); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeLogStreams returns a list of CloudWatch Logs log streams.
func (s *LogsService) DescribeLogStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	prefix := request.GetParamLowerFirst(req.Parameters, "LogStreamNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 50
	}
	orderBy := request.GetParamLowerFirst(req.Parameters, "OrderBy")
	descending := request.GetBoolParam(req.Parameters, "Descending")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if orderBy == "LastEventTime" {
		if prefix != "" {
			return nil, NewLogsError("InvalidParameterException",
				"Cannot specify logStreamNamePrefix when orderBy is LastEventTime", 400)
		}

		allStreams, err := fetchAllLogStreams(store, logGroupName, "")
		if err != nil {
			return nil, mapStoreError(err)
		}

		sort.Slice(allStreams, func(i, j int) bool {
			if descending {
				return allStreams[i].LastEventTs > allStreams[j].LastEventTs
			}
			return allStreams[i].LastEventTs < allStreams[j].LastEventTs
		})

		// Use offset-based pagination for LastEventTime ordering because
		// new events can change the sort order between paginated calls,
		// making key-based markers unreliable (items shift position).
		// DescribeLogStreams uses forward-only pagination (AWS does not
		// support backward pagination for this operation).
		_, offset, err := logsstore.ParsePaginationToken(nextToken)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		if offset > len(allStreams) {
			offset = len(allStreams)
		}

		endIdx := offset + int(limit)
		if endIdx > len(allStreams) {
			endIdx = len(allStreams)
		}

		var pageItems []*logsstore.LogStream
		if offset < len(allStreams) {
			pageItems = allStreams[offset:endIdx]
		}

		logStreams := make([]map[string]interface{}, 0, len(pageItems))
		for _, ls := range pageItems {
			logStreams = append(logStreams, logStreamToMap(ls))
		}

		resp := map[string]interface{}{
			"logStreams": logStreams,
		}
		if endIdx < len(allStreams) {
			resp["nextToken"] = logsstore.EncodePaginationToken(logsstore.PaginationForward, endIdx)
		}
		return resp, nil
	}

	// orderBy != "LastEventTime" (default: LogStreamName)
	if descending {
		allStreams, err := fetchAllLogStreams(store, logGroupName, prefix)
		if err != nil {
			return nil, mapStoreError(err)
		}

		sort.Slice(allStreams, func(i, j int) bool {
			return allStreams[i].Name > allStreams[j].Name
		})

		result := pagination.PaginateSlice(allStreams, nextToken, int(limit), func(ls *logsstore.LogStream) string {
			return ls.Name
		})

		logStreams := make([]map[string]interface{}, 0, len(result.Items))
		for _, ls := range result.Items {
			logStreams = append(logStreams, logStreamToMap(ls))
		}

		resp := map[string]interface{}{
			"logStreams": logStreams,
		}
		if result.NextMarker != "" {
			resp["nextToken"] = result.NextMarker
		}
		return resp, nil
	}

	streams, nextMarker, err := store.ListLogStreams(logGroupName, prefix, nextToken, int(limit))
	if err != nil {
		return nil, mapStoreError(err)
	}

	logStreams := make([]map[string]interface{}, 0, len(streams))
	for _, ls := range streams {
		logStreams = append(logStreams, logStreamToMap(ls))
	}

	resp := map[string]interface{}{
		"logStreams": logStreams,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
}

func fetchAllLogStreams(store *logsstore.Store, logGroupName, prefix string) ([]*logsstore.LogStream, error) {
	var all []*logsstore.LogStream
	marker := ""
	for {
		streams, nextMarker, err := store.ListLogStreams(logGroupName, prefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, streams...)
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}
	return all, nil
}

func logStreamToMap(ls *logsstore.LogStream) map[string]interface{} {
	return map[string]interface{}{
		"logStreamName":       ls.Name,
		"arn":                 ls.ARN,
		"creationTime":        ls.CreatedAt.UnixMilli(),
		"firstEventTimestamp": ls.FirstEventTs,
		"lastEventTimestamp":  ls.LastEventTs,
		"lastIngestionTime":   ls.LastIngestionTs,
		"uploadSequenceToken": ls.UploadSequenceToken,
	}
}

// ListLogStreams returns a list of CloudWatch Logs log streams.
func (s *LogsService) ListLogStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeLogStreams(ctx, reqCtx, req)
}

// --- PutLogEvents ---

const (
	// maxEventsTimeSpan is the maximum allowed time span (in milliseconds)
	// for a single PutLogEvents batch. AWS rejects the entire batch if the
	// span between the earliest and latest event exceeds 24 hours.
	maxEventsTimeSpan int64 = 24 * 60 * 60 * 1000

	// tooNewThreshold is the maximum future offset (in milliseconds) for
	// an event timestamp. Events more than 2 hours in the future are
	// rejected individually.
	tooNewThreshold int64 = 2 * 60 * 60 * 1000

	// tooOldThreshold is the maximum age (in milliseconds) for an event
	// timestamp. Events older than 14 days are rejected individually.
	tooOldThreshold int64 = 14 * 24 * 60 * 60 * 1000
)

// validateLogEvents checks that log events satisfy the PutLogEvents
// constraints required by AWS CloudWatch Logs:
//   - Events must be in chronological order (by timestamp).
//   - The time span of the batch must not exceed 24 hours.
//   - Events more than 2 hours in the future or older than 14 days are
//     individually rejected.
//
// Returns the filtered valid events and, if any events were rejected, a map
// suitable for inclusion in the response as rejectedLogEventsInfo.
func validateLogEvents(events []logsstore.LogEntry) ([]logsstore.LogEntry, map[string]interface{}, error) {
	now := time.Now().UnixMilli()

	// Chronological order check must be performed on ALL events in the
	// batch, not just the age-valid subset. AWS rejects the entire batch
	// if any event is out of order, regardless of whether some events are
	// later individually rejected for being too old or too new.
	for i := 1; i < len(events); i++ {
		if events[i].Timestamp < events[i-1].Timestamp {
			return nil, nil, awserrors.NewAWSError("InvalidParameterException",
				"log events in the batch must be in chronological order", 400)
		}
	}

	var valid []logsstore.LogEntry
	var tooOldEndIndex int
	tooNewStartIndex := -1

	for i, e := range events {
		if e.Timestamp > now+tooNewThreshold {
			if tooNewStartIndex == -1 || i < tooNewStartIndex {
				tooNewStartIndex = i
			}
			continue
		}
		if e.Timestamp < now-tooOldThreshold {
			tooOldEndIndex = i + 1
			continue
		}
		valid = append(valid, e)
	}

	if len(valid) == 0 {
		rejected := buildRejectedInfo(tooOldEndIndex, tooNewStartIndex, len(events))
		return nil, rejected, nil
	}

	span := valid[len(valid)-1].Timestamp - valid[0].Timestamp
	if span > maxEventsTimeSpan {
		return nil, nil, awserrors.NewAWSError("InvalidParameterException",
			"Events span must not exceed 24 hours", 400)
	}

	rejected := buildRejectedInfo(tooOldEndIndex, tooNewStartIndex, len(events))
	return valid, rejected, nil
}

// buildRejectedInfo constructs the rejectedLogEventsInfo map from the
// computed too-old and too-new indices. If no events were rejected an
// empty map is returned.
func buildRejectedInfo(tooOldEndIndex, tooNewStartIndex, totalEvents int) map[string]interface{} {
	if tooOldEndIndex == 0 && tooNewStartIndex == -1 {
		return nil
	}
	info := make(map[string]interface{})
	if tooOldEndIndex > 0 {
		info["tooOldLogEventEndIndex"] = tooOldEndIndex
	}
	if tooNewStartIndex >= 0 {
		info["tooNewLogEventStartIndex"] = tooNewStartIndex
	}
	return info
}

// PutLogEvents uploads log events to the specified CloudWatch Logs log stream.
func (s *LogsService) PutLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	logStreamName := request.GetParamLowerFirst(req.Parameters, "LogStreamName")
	sequenceToken := request.GetParamLowerFirst(req.Parameters, "SequenceToken")

	if logGroupName == "" || logStreamName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if sequenceToken != "" {
		ls, err := store.GetLogStream(logGroupName, logStreamName)
		if err != nil {
			return nil, mapStoreError(err)
		}
		if ls.UploadSequenceToken != sequenceToken {
			return nil, awserrors.NewAWSError("InvalidSequenceTokenException",
				fmt.Sprintf("The sequence token is not valid. Expected: %s, Received: %s", ls.UploadSequenceToken, sequenceToken), 400)
		}
	}

	events := parseLogEvents(req)
	if len(events) == 0 {
		return nil, ErrMissingParameter
	}

	if len(events) > logsstore.MaxChunkSize {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Maximum number of log events in a single batch is %d", logsstore.MaxChunkSize), 400)
	}

	// Validate timestamp ordering, time span, and age constraints
	// per AWS CloudWatch Logs PutLogEvents specification.
	validEvents, rejectedInfo, valErr := validateLogEvents(events)
	if valErr != nil {
		return nil, valErr
	}
	if len(validEvents) == 0 {
		// All events were rejected (too old or too new); nothing to write.
		resp := map[string]interface{}{"nextSequenceToken": ""}
		if len(rejectedInfo) > 0 {
			resp["rejectedLogEventsInfo"] = rejectedInfo
		}
		return resp, nil
	}

	nextToken, err := store.PutLogEvents(logGroupName, logStreamName, validEvents)
	if err != nil {
		return nil, mapStoreError(err)
	}

	region := reqCtx.GetRegion()

	// Evaluate metric filters and deliver subscription events asynchronously.
	// These side-effects do not affect the PutLogEvents response and should
	// not block it — AWS processes them in the background after ingestion.
	eventsCopy := make([]logsstore.LogEntry, len(validEvents))
	copy(eventsCopy, validEvents)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("Panic in async metric/subscription processing",
					logs.String("logGroup", logGroupName),
					logs.Any("panic", r))
			}
		}()
		s.evaluateMetricFilters(store, region, logGroupName, eventsCopy)
		s.deliverSubscriptionEvents(store, region, logGroupName, logStreamName, eventsCopy)
	}()

	resp := map[string]interface{}{
		"nextSequenceToken": nextToken,
	}
	if len(rejectedInfo) > 0 {
		resp["rejectedLogEventsInfo"] = rejectedInfo
	}
	return resp, nil
}

func parseLogEvents(req *request.ParsedRequest) []logsstore.LogEntry {
	var events []logsstore.LogEntry
	for i := 1; ; i++ {
		timestamp := int64(request.GetIntParam(req.Parameters, "LogEvents."+strconv.Itoa(i)+".Timestamp"))
		message := request.GetParamLowerFirst(req.Parameters, "LogEvents."+strconv.Itoa(i)+".Message")
		if timestamp == 0 && message == "" {
			break
		}
		events = append(events, logsstore.LogEntry{
			Timestamp: timestamp,
			Message:   message,
		})
	}
	if len(events) > 0 {
		return events
	}
	return parseLogEventsFromMap(req)
}

func parseLogEventsFromMap(req *request.ParsedRequest) []logsstore.LogEntry {
	var logEvents interface{}
	if le, ok := req.Parameters["logEvents"]; ok {
		logEvents = le
	} else if le, ok := req.Parameters["LogEvents"]; ok {
		logEvents = le
	}
	if logEvents == nil {
		return nil
	}

	eventsArray, ok := logEvents.([]interface{})
	if !ok {
		return nil
	}

	var events []logsstore.LogEntry
	for _, e := range eventsArray {
		eventMap, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		ts := int64(0)
		if t, ok := eventMap["timestamp"].(float64); ok {
			ts = int64(t)
		} else if t, ok := eventMap["Timestamp"].(float64); ok {
			ts = int64(t)
		}
		msg := ""
		if m, ok := eventMap["message"].(string); ok {
			msg = m
		} else if m, ok := eventMap["Message"].(string); ok {
			msg = m
		}
		if ts > 0 || msg != "" {
			events = append(events, logsstore.LogEntry{
				Timestamp: ts,
				Message:   msg,
			})
		}
	}

	return events
}

// --- Metric filter evaluation ---

// evaluateMetricFilters evaluates metric filters for the given log group and
// emits CloudWatch metrics for matched log entries.
func (s *LogsService) evaluateMetricFilters(store *logsstore.Store, region, logGroupName string, events []logsstore.LogEntry) {
	if s.cwMetricInvoker == nil {
		return
	}

	filters, _, err := store.ListMetricFilters(logGroupName, "", "", 1000)
	if err != nil || len(filters) == 0 {
		return
	}

	matcher := filterpattern.NewMatcher()
	now := time.Now()

	for _, event := range events {
		for _, filter := range filters {
			matched := matcher.Matches(filter.FilterPattern, event.Message)
			for _, transform := range filter.MetricTransformations {
				var value float64
				var shouldEmit bool

				if matched {
					value = 1.0
					if transform.MetricValue != "" && transform.MetricValue != "1" {
						if v, err := strconv.ParseFloat(transform.MetricValue, 64); err == nil {
							value = v
						}
					}
					shouldEmit = true
				} else if transform.DefaultValueSet {
					value = transform.DefaultValue
					shouldEmit = true
				}

				if shouldEmit {
					ts := time.UnixMilli(event.Timestamp)
					if ts.IsZero() || ts.After(now) {
						ts = now
					}

					if err := s.cwMetricInvoker.PutMetricData(region, transform.MetricNamespace, transform.MetricName, value, ts); err != nil {
						logs.Error("Failed to put metric data", logs.Err(err))
					}
				}
			}
		}
	}
}

// applyMetricFiltersByRegion resolves the store by region and evaluates metric
// filters. Used by bus handlers that lack an HTTP request context.
func (s *LogsService) applyMetricFiltersByRegion(region, logGroupName string, events []logsstore.LogEntry) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return
	}
	s.evaluateMetricFilters(store, region, logGroupName, events)
}

// --- Subscription filter evaluation ---

// deliverSubscriptionEvents matches events against subscription filters and
// delivers matched events either via the event bus or directly.
func (s *LogsService) deliverSubscriptionEvents(store *logsstore.Store, region, logGroupName, logStreamName string, events []logsstore.LogEntry) {
	filters, err := store.ListSubscriptionFilters(logGroupName, "")
	if err != nil || len(filters) == 0 {
		return
	}

	matcher := filterpattern.NewMatcher()

	for _, filter := range filters {
		var matched []logsstore.LogEntry
		for _, event := range events {
			if filter.FilterPattern == "" || matcher.Matches(filter.FilterPattern, event.Message) {
				matched = append(matched, event)
			}
		}
		if len(matched) == 0 {
			continue
		}

		payload := s.buildSubscriptionPayload(filter, logGroupName, logStreamName, matched)
		compressed, err := compressJSON(payload)
		if err != nil {
			continue
		}

		if s.bus != nil {
			evt := &eventbus.CloudWatchLogDeliveryEvent{
				LogGroup:       logGroupName,
				LogStream:      logStreamName,
				DestinationArn: filter.DestinationArn,
				Payload:        compressed,
			}
			evt.Region = region
			if err := s.bus.Publish(context.Background(), evt); err != nil {
				logs.Warn("Failed to publish log delivery event", logs.Err(err))
			}
		} else {
			s.deliverDirect(filter.DestinationArn, compressed)
		}
	}
}

// applySubscriptionFiltersByRegion resolves the store by region and evaluates
// subscription filters. Used by bus handlers that lack an HTTP request context.
func (s *LogsService) applySubscriptionFiltersByRegion(region, logGroupName, logStreamName string, events []logsstore.LogEntry) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return
	}
	s.deliverSubscriptionEvents(store, region, logGroupName, logStreamName, events)
}

// --- Subscription payload building ---

func (s *LogsService) buildSubscriptionPayload(
	filter *logsstore.SubscriptionFilter,
	logGroupName, logStreamName string,
	events []logsstore.LogEntry,
) map[string]interface{} {
	logEvents := make([]map[string]interface{}, len(events))
	for i, e := range events {
		logEvents[i] = map[string]interface{}{
			"id":        generateEventID(),
			"timestamp": e.Timestamp,
			"message":   e.Message,
		}
	}

	return map[string]interface{}{
		"owner":               s.accountID,
		"logGroup":            logGroupName,
		"logStream":           logStreamName,
		"subscriptionFilters": []string{filter.FilterName},
		"messageType":         "DATA_MESSAGE",
		"logEvents":           logEvents,
	}
}

// --- Direct delivery (non-bus fallback) ---

// deliverDirect sends compressed subscription payload to Lambda or Kinesis.
func (s *LogsService) deliverDirect(destArn string, compressed []byte) {
	if arn.IsLambdaARN(destArn) {
		s.invokeLambda(destArn, compressed)
	} else if arn.IsKinesisARN(destArn) {
		s.putToKinesis(destArn, compressed)
	}
}

// invokeLambda delivers compressed log data to a Lambda function.
func (s *LogsService) invokeLambda(destArn string, compressed []byte) {
	if s.bus == nil {
		return
	}

	functionName := arn.ExtractFunctionNameFromARN(destArn)
	encodedData := base64.StdEncoding.EncodeToString(compressed)

	payload, err := json.Marshal(map[string]interface{}{
		"awslogs": map[string]interface{}{"data": encodedData},
	})
	if err != nil {
		return
	}

	if _, _, err := s.bus.LambdaInvoker().InvokeForGateway(context.Background(), functionName, payload); err != nil {
		logs.Warn("Failed to invoke Lambda for subscription filter delivery", logs.Err(err))
	}
}

// putToKinesis delivers compressed log data to a Kinesis stream.
func (s *LogsService) putToKinesis(destArn string, compressed []byte) {
	if s.bus == nil {
		return
	}

	streamName := arn.ExtractStreamNameFromARN(destArn)
	encodedData := base64.StdEncoding.EncodeToString(compressed)
	ctx := context.Background()

	shards, err := s.bus.KinesisInvoker().ListShards(ctx, streamName)
	if err != nil || len(shards) == 0 {
		return
	}

	var activeShardID string
	for _, shard := range shards {
		if shard.SequenceNumberRangeEnd == "" {
			activeShardID = shard.ShardID
			break
		}
	}

	if activeShardID == "" {
		return
	}

	envelope, err := json.Marshal(map[string]interface{}{
		"awslogs": map[string]interface{}{"data": encodedData},
	})
	if err != nil {
		return
	}

	b64Envelope := base64.StdEncoding.EncodeToString(envelope)
	if _, err := s.bus.KinesisInvoker().PutRecord(ctx, streamName, activeShardID, []byte(b64Envelope)); err != nil {
		logs.Warn("Failed to deliver subscription filter log events to Kinesis", logs.Err(err))
	}
}

// --- GetLogEvents / FilterLogEvents ---

// GetLogEvents retrieves log events from the specified CloudWatch Logs log stream.
func (s *LogsService) GetLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	logStreamName := request.GetParamLowerFirst(req.Parameters, "LogStreamName")

	if logGroupName == "" || logStreamName == "" {
		return nil, ErrMissingParameter
	}

	startTime := int64(request.GetIntParam(req.Parameters, "StartTime"))
	endTime := int64(request.GetIntParam(req.Parameters, "EndTime"))
	limit := int(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 10000
	}
	startFromHead := request.GetBoolParam(req.Parameters, "StartFromHead")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	events, nextForwardToken, nextBackwardToken, err := store.GetLogEvents(logGroupName, logStreamName, startTime, endTime, limit, startFromHead, nextToken)
	if err != nil {
		return nil, mapStoreError(err)
	}

	outputEvents := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		outputEvents = append(outputEvents, logEventToResponse(e))
	}

	return map[string]interface{}{
		"events":            outputEvents,
		"nextForwardToken":  nextForwardToken,
		"nextBackwardToken": nextBackwardToken,
	}, nil
}

// FilterLogEvents filters log events from the specified CloudWatch Logs log group.
func (s *LogsService) FilterLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	logStreamNames := request.GetStringList(req.Parameters, "LogStreamNames")
	logStreamNamePrefix := request.GetParamLowerFirst(req.Parameters, "LogStreamNamePrefix")

	if len(logStreamNames) > 0 && logStreamNamePrefix != "" {
		return nil, NewLogsError("InvalidParameterException",
			"Cannot specify both logStreamNames and logStreamNamePrefix", 400)
	}

	startTime := int64(request.GetIntParam(req.Parameters, "StartTime"))
	endTime := int64(request.GetIntParam(req.Parameters, "EndTime"))
	filterPattern := request.GetParamLowerFirst(req.Parameters, "FilterPattern")
	limit := int(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 10000
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	// AWS spec: startFromHead defaults to true (oldest first / ascending).
	startFromHead := request.GetBoolParamDefault(req.Parameters, "StartFromHead", true)

	// AWS spec: startFromHead=false is supported only when startTime is on or
	// after Jan 1, 2024 00:00:00 UTC.
	if !startFromHead && startTime > 0 && startTime < jan012024Millis {
		return nil, NewLogsError("InvalidParameterException",
			"Setting startFromHead to false is supported only when startTime is on or after Jan 1, 2024 00:00:00 UTC", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if logStreamNamePrefix != "" {
		prefixStreams, err := fetchAllLogStreams(store, logGroupName, logStreamNamePrefix)
		if err != nil {
			return nil, mapStoreError(err)
		}
		for _, ls := range prefixStreams {
			logStreamNames = append(logStreamNames, ls.Name)
		}
	}

	events, searchedStreams, nextMarker, err := store.FilterLogEvents(logGroupName, logStreamNames, startTime, endTime, filterPattern, limit, startFromHead, nextToken)
	if err != nil {
		return nil, mapStoreError(err)
	}

	outputEvents := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		outputEvents = append(outputEvents, logEventToResponse(e))
	}

	searchedStreamNames := make([]map[string]interface{}, 0, len(searchedStreams))
	for name := range searchedStreams {
		searchedStreamNames = append(searchedStreamNames, map[string]interface{}{
			"logStreamName":      name,
			"searchedCompletely": true,
		})
	}

	resp := map[string]interface{}{
		"events":             outputEvents,
		"searchedLogStreams": searchedStreamNames,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
}

// --- Helpers ---

func generateEventID() string {
	now := time.Now().UnixNano()
	b := make([]byte, logEventIDSize)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%x", now)
	}
	return fmt.Sprintf("%x%x", now, b)
}
