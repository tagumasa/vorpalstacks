package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/filterpattern"
)

const logEventIDSize = 16

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
	input := CreateLogStreamInput{
		LogGroupName:  request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		LogStreamName: request.GetParamLowerFirst(req.Parameters, "LogStreamName"),
		Region:        reqCtx.GetRegion(),
	}

	if err := s.createLogStreamCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteLogStream deletes a CloudWatch Logs log stream.
func (s *LogsService) DeleteLogStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := DeleteLogStreamInput{
		LogGroupName:  request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		LogStreamName: request.GetParamLowerFirst(req.Parameters, "LogStreamName"),
		Region:        reqCtx.GetRegion(),
	}

	if err := s.deleteLogStreamCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeLogStreams returns a list of CloudWatch Logs log streams.
func (s *LogsService) DescribeLogStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		logGroupName = request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
		if logGroupName == "" {
			return nil, ErrMissingParameter
		}
	}

	input := DescribeLogStreamsInput{
		LogGroupName:        logGroupName,
		LogStreamNamePrefix: request.GetParamLowerFirst(req.Parameters, "LogStreamNamePrefix"),
		OrderBy:             request.GetParamLowerFirst(req.Parameters, "OrderBy"),
		Descending:          request.GetBoolParam(req.Parameters, "Descending"),
		NextToken:           request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Limit:               int32(request.GetIntParam(req.Parameters, "Limit")),
		Region:              reqCtx.GetRegion(),
	}

	result, err := s.describeLogStreamsCore(input)
	if err != nil {
		return nil, err
	}

	logStreams := make([]map[string]interface{}, 0, len(result.LogStreams))
	for _, ls := range result.LogStreams {
		logStreams = append(logStreams, logStreamToMap(ls))
	}

	resp := map[string]interface{}{
		"logStreams": logStreams,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
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

// PutLogEvents uploads log events to the specified CloudWatch Logs log stream.
func (s *LogsService) PutLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	events := parseLogEvents(req)

	input := PutLogEventsInput{
		LogGroupName:  request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		LogStreamName: request.GetParamLowerFirst(req.Parameters, "LogStreamName"),
		Events:        events,
		Region:        reqCtx.GetRegion(),
	}

	result, err := s.putLogEventsCore(input)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"nextSequenceToken": result.NextSequenceToken,
	}
	if len(result.RejectedLogEvents) > 0 {
		resp["rejectedLogEventsInfo"] = result.RejectedLogEvents
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
		logs.Warn("cwMetricInvoker is nil, skipping metric filter evaluation",
			logs.String("logGroup", logGroupName))
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
	if logGroupName == "" {
		logGroupName = request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	}

	input := GetLogEventsInput{
		LogGroupName:  logGroupName,
		LogStreamName: request.GetParamLowerFirst(req.Parameters, "LogStreamName"),
		StartTime:     int64(request.GetIntParam(req.Parameters, "StartTime")),
		EndTime:       int64(request.GetIntParam(req.Parameters, "EndTime")),
		Limit:         int32(request.GetIntParam(req.Parameters, "Limit")),
		StartFromHead: request.GetBoolParam(req.Parameters, "StartFromHead"),
		NextToken:     request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Region:        reqCtx.GetRegion(),
	}

	result, err := s.getLogEventsCore(input)
	if err != nil {
		return nil, err
	}

	outputEvents := make([]map[string]interface{}, 0, len(result.Events))
	for _, e := range result.Events {
		outputEvents = append(outputEvents, logEventToResponse(e))
	}

	return map[string]interface{}{
		"events":            outputEvents,
		"nextForwardToken":  result.NextForwardToken,
		"nextBackwardToken": result.NextBackwardToken,
	}, nil
}

// FilterLogEvents filters log events from the specified CloudWatch Logs log group.
func (s *LogsService) FilterLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		logGroupName = request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	}

	input := FilterLogEventsInput{
		LogGroupName:      logGroupName,
		LogStreamNames:    request.GetStringList(req.Parameters, "LogStreamNames"),
		LogStreamNamePref: request.GetParamLowerFirst(req.Parameters, "LogStreamNamePrefix"),
		StartTime:         int64(request.GetIntParam(req.Parameters, "StartTime")),
		EndTime:           int64(request.GetIntParam(req.Parameters, "EndTime")),
		FilterPattern:     request.GetParamLowerFirst(req.Parameters, "FilterPattern"),
		Limit:             int32(request.GetIntParam(req.Parameters, "Limit")),
		StartFromHead:     request.GetBoolParamDefault(req.Parameters, "StartFromHead", true),
		NextToken:         request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Region:            reqCtx.GetRegion(),
	}

	result, err := s.filterLogEventsCore(input)
	if err != nil {
		return nil, err
	}

	outputEvents := make([]map[string]interface{}, 0, len(result.Events))
	for _, e := range result.Events {
		outputEvents = append(outputEvents, logEventToResponse(e))
	}

	searchedStreamNames := make([]map[string]interface{}, 0, len(result.SearchedStreams))
	for name := range result.SearchedStreams {
		searchedStreamNames = append(searchedStreamNames, map[string]interface{}{
			"logStreamName":      name,
			"searchedCompletely": true,
		})
	}

	resp := map[string]interface{}{
		"events":             outputEvents,
		"searchedLogStreams": searchedStreamNames,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
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
