package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/filterpattern"
)

func logEventToResponse(e *logsstore.OutputLogEvent) map[string]interface{} {
	return map[string]interface{}{
		"timestamp":     e.Timestamp,
		"message":       e.Message,
		"ingestionTime": e.IngestionTime,
	}
}

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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streams, nextMarker, err := store.ListLogStreams(logGroupName, prefix, nextToken, int(limit))
	if err != nil {
		return nil, mapStoreError(err)
	}

	logStreams := make([]map[string]interface{}, 0)
	for _, ls := range streams {
		logStreams = append(logStreams, map[string]interface{}{
			"logStreamName":       ls.Name,
			"arn":                 ls.ARN,
			"creationTime":        ls.CreatedAt.UnixMilli(),
			"firstEventTimestamp": ls.FirstEventTs,
			"lastEventTimestamp":  ls.LastEventTs,
			"lastIngestionTime":   ls.LastIngestionTs,
			"uploadSequenceToken": ls.UploadSequenceToken,
		})
	}

	response := map[string]interface{}{
		"logStreams": logStreams,
	}
	if nextMarker != "" {
		response["nextToken"] = nextMarker
	}

	return response, nil
}

// ListLogStreams returns a list of CloudWatch Logs log streams.
func (s *LogsService) ListLogStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeLogStreams(ctx, reqCtx, req)
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
			return nil, NewLogsError("InvalidSequenceTokenException",
				fmt.Sprintf("The sequence token is not valid. Expected: %s, Received: %s", ls.UploadSequenceToken, sequenceToken), 400)
		}
	}

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

	if len(events) == 0 {
		events = parseLogEventsFromMap(req)
	}

	nextToken, err := store.PutLogEvents(logGroupName, logStreamName, events)
	if err != nil {
		return nil, mapStoreError(err)
	}

	s.applyMetricFilters(reqCtx, logGroupName, events)

	s.applySubscriptionFilters(reqCtx, logGroupName, logStreamName, events)

	return map[string]interface{}{
		"nextSequenceToken": nextToken,
	}, nil
}

func (s *LogsService) applyMetricFilters(reqCtx *request.RequestContext, logGroupName string, events []logsstore.LogEntry) {
	store, err := s.store(reqCtx)
	if err != nil {
		return
	}

	filters, _, err := store.ListMetricFilters(logGroupName, "", "", 1000)
	if err != nil || len(filters) == 0 {
		return
	}

	cwMetricStore, err := s.getMetricStore(reqCtx)
	if err != nil {
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

					datum := cwstore.MetricDatum{
						MetricName: transform.MetricName,
						Value:      value,
						Timestamp:  ts,
					}
					if err := cwMetricStore.PutMetricData(transform.MetricNamespace, []cwstore.MetricDatum{datum}); err != nil {
						log.Printf("Failed to put metric data: %v", err)
					}
				}
			}
		}
	}
}

func parseLogEventsFromMap(req *request.ParsedRequest) []logsstore.LogEntry {
	var events []logsstore.LogEntry

	var logEvents interface{}
	if le, ok := req.Parameters["logEvents"]; ok {
		logEvents = le
	} else if le, ok := req.Parameters["LogEvents"]; ok {
		logEvents = le
	}

	if logEvents != nil {
		if eventsArray, ok := logEvents.([]interface{}); ok {
			for _, e := range eventsArray {
				if eventMap, ok := e.(map[string]interface{}); ok {
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
			}
		}
	}

	return events
}

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
	startFromHead := request.GetParamLowerFirst(req.Parameters, "StartFromHead") == "true"
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	events, nextForwardToken, startIndex, err := store.GetLogEvents(logGroupName, logStreamName, startTime, endTime, limit, startFromHead, nextToken)
	if err != nil {
		return nil, mapStoreError(err)
	}

	outputEvents := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		outputEvents = append(outputEvents, logEventToResponse(e))
	}

	nextBackwardToken := ""
	if len(events) > 0 && startIndex > 0 {
		nextBackwardToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", startIndex)))
	}

	return map[string]interface{}{
		"events":            outputEvents,
		"nextForwardToken":  nextForwardToken,
		"nextBackwardToken": nextBackwardToken,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FilterLogEvents filters log events from the specified CloudWatch Logs log group.
func (s *LogsService) FilterLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	var logStreamNames []string
	for i := 1; ; i++ {
		name := request.GetParamLowerFirst(req.Parameters, "LogStreamNames."+strconv.Itoa(i))
		if name == "" {
			break
		}
		logStreamNames = append(logStreamNames, name)
	}

	startTime := int64(request.GetIntParam(req.Parameters, "StartTime"))
	endTime := int64(request.GetIntParam(req.Parameters, "EndTime"))
	filterPattern := request.GetParamLowerFirst(req.Parameters, "FilterPattern")
	limit := int(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 10000
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	events, searchedStreams, err := store.FilterLogEvents(logGroupName, logStreamNames, startTime, endTime, filterPattern, limit)
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

	return map[string]interface{}{
		"events":             outputEvents,
		"searchedLogStreams": searchedStreamNames,
		"nextToken":          "",
	}, nil
}

func (s *LogsService) applySubscriptionFilters(reqCtx *request.RequestContext, logGroupName, logStreamName string, events []logsstore.LogEntry) {
	store, err := s.store(reqCtx)
	if err != nil {
		return
	}

	filters, err := store.ListSubscriptionFilters(logGroupName, "")
	if err != nil || len(filters) == 0 {
		return
	}

	matcher := filterpattern.NewMatcher()

	for _, filter := range filters {
		var matchedEvents []logsstore.LogEntry
		for _, event := range events {
			if filter.FilterPattern == "" || matcher.Matches(filter.FilterPattern, event.Message) {
				matchedEvents = append(matchedEvents, event)
			}
		}

		if len(matchedEvents) == 0 {
			continue
		}

		s.deliverToDestination(reqCtx, filter, logGroupName, logStreamName, matchedEvents)
	}
}

func (s *LogsService) deliverToDestination(
	reqCtx *request.RequestContext,
	filter *logsstore.SubscriptionFilter,
	logGroupName, logStreamName string,
	events []logsstore.LogEntry,
) {
	payload := s.buildSubscriptionPayload(filter, logGroupName, logStreamName, events)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if err := json.NewEncoder(gw).Encode(payload); err != nil {
		return
	}
	if err := gw.Close(); err != nil {
		return
	}
	compressed := buf.Bytes()

	if arn.IsLambdaARN(filter.DestinationArn) {
		s.deliverToLambda(reqCtx, filter.DestinationArn, compressed)
	} else if arn.IsKinesisARN(filter.DestinationArn) {
		s.deliverToKinesis(reqCtx, filter.DestinationArn, compressed)
	}
}

func (s *LogsService) buildSubscriptionPayload(
	filter *logsstore.SubscriptionFilter,
	logGroupName, logStreamName string,
	events []logsstore.LogEntry,
) map[string]interface{} {
	var logEvents []map[string]interface{}
	for _, e := range events {
		logEvents = append(logEvents, map[string]interface{}{
			"id":        generateEventID(),
			"timestamp": e.Timestamp,
			"message":   e.Message,
		})
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

func (s *LogsService) deliverToLambda(reqCtx *request.RequestContext, destArn string, compressedData []byte) {
	_, _, region, _, _ := arn.SplitARN(destArn)
	invoker := s.getLambdaInvoker(region)
	if invoker == nil {
		invoker = s.getLambdaInvoker(reqCtx.GetRegion())
	}
	if invoker == nil {
		return
	}

	functionName := arn.ExtractFunctionNameFromARN(destArn)
	encodedData := base64.StdEncoding.EncodeToString(compressedData)

	payload := map[string]interface{}{
		"awslogs": map[string]interface{}{
			"data": encodedData,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	_, _, err = invoker.InvokeForGateway(context.Background(), functionName, payloadBytes)
	if err != nil {
		return
	}
}

func (s *LogsService) deliverToKinesis(reqCtx *request.RequestContext, destArn string, compressedData []byte) {
	_, _, region, _, _ := arn.SplitARN(destArn)
	kinesisStore, ok := s.getKinesisStore(region)
	if !ok {
		kinesisStore, ok = s.getKinesisStore(reqCtx.GetRegion())
		if !ok {
			return
		}
	}

	streamName := arn.ExtractStreamNameFromARN(destArn)
	partitionKey := generatePartitionKey()
	encodedData := base64.StdEncoding.EncodeToString(compressedData)

	shards, err := kinesisStore.ListShards(streamName, nil, "", 0)
	if err != nil || len(shards) == 0 {
		return
	}

	var activeShardID string
	for _, shard := range shards {
		if shard.SequenceNumberRange.EndingSequenceNumber == "" {
			activeShardID = shard.ShardID
			break
		}
	}

	if activeShardID == "" {
		return
	}

	_, err = kinesisStore.PutRecord(streamName, activeShardID, partitionKey, encodedData)
	if err != nil {
		return
	}
}

func generateEventID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		b = []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	}
	return fmt.Sprintf("%x%x", time.Now().UnixNano(), b)
}

func generatePartitionKey() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		b = []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	}
	return fmt.Sprintf("%x", b)
}
