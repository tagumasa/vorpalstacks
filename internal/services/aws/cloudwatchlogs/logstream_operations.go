package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

// formatLogEvent shapes one OutputLogEvent, whose modelled members are
// exactly {ingestionTime, message, timestamp}; logStreamName belongs to
// FilteredLogEvent alone, so the FilterLogEvents handler adds it.
func formatLogEvent(e *logsstore.OutputLogEvent) map[string]interface{} {
	return map[string]interface{}{
		"timestamp":     e.Timestamp,
		"message":       e.Message,
		"ingestionTime": e.IngestionTime,
	}
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
	logGroupName, err := logGroupNameOrIdentifier(req.Parameters)
	if err != nil {
		return nil, err
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
		logStreams = append(logStreams, formatLogStream(ls))
	}

	resp := map[string]interface{}{
		"logStreams": logStreams,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

func formatLogStream(ls *logsstore.LogStream) map[string]interface{} {
	return map[string]interface{}{
		"logStreamName": ls.Name,
		"arn":           ls.ARN,
		"creationTime":  ls.CreatedAt.UnixMilli(),
		// storedBytes is no longer supported for log streams and is
		// always reported as zero (the LogStream shape's deprecated
		// member documentation); the log-group member keeps its value.
		"storedBytes":         int64(0),
		"firstEventTimestamp": ls.FirstEventTs,
		"lastEventTimestamp":  ls.LastEventTs,
		"lastIngestionTime":   ls.LastIngestionTs,
		"uploadSequenceToken": ls.UploadSequenceToken,
	}
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

	// nextSequenceToken carries a minimum length of 1 on the wire shape;
	// a batch whose every event was rejected produces no token to
	// continue from, so the member is omitted rather than sent empty.
	resp := map[string]interface{}{}
	if result.NextSequenceToken != "" {
		resp["nextSequenceToken"] = result.NextSequenceToken
	}
	if len(result.RejectedLogEvents) > 0 {
		resp["rejectedLogEventsInfo"] = result.RejectedLogEvents
	}
	return resp, nil
}

func parseLogEvents(req *request.ParsedRequest) []PutLogEvent {
	var events []PutLogEvent
	for i := 1; ; i++ {
		tsKey := "LogEvents." + strconv.Itoa(i) + ".Timestamp"
		// The Timestamp shape's range starts at zero, so an explicit
		// timestamp of 0 is a legal member value: presence, not the
		// value, decides whether the member was sent. An unparseable
		// timestamp keeps the absent treatment — the seam then rejects
		// the request for the missing required member.
		ts, tsPresent, tsErr := request.GetIntParamStrictCaseInsensitive(req.Parameters, tsKey)
		message := request.GetParamLowerFirst(req.Parameters, "LogEvents."+strconv.Itoa(i)+".Message")
		if !tsPresent && message == "" {
			break
		}
		events = append(events, PutLogEvent{
			LogEntry: logsstore.LogEntry{
				Timestamp: int64(ts),
				Message:   message,
			},
			TimestampSet: tsPresent && tsErr == nil,
		})
	}
	if len(events) > 0 {
		return events
	}
	return parseLogEventsFromMap(req)
}

func parseLogEventsFromMap(req *request.ParsedRequest) []PutLogEvent {
	eventsArray := request.GetArrayParamLowerFirst(req.Parameters, "LogEvents")

	var events []PutLogEvent
	for _, e := range eventsArray {
		eventMap, ok := e.(map[string]interface{})
		if !ok {
			// A non-object element cannot carry the required members; an
			// empty event reaches the seam's member validation and
			// rejects the request instead of silently dropping it.
			events = append(events, PutLogEvent{})
			continue
		}
		ts := int64(0)
		tsSet := false
		if t, ok := eventMap["timestamp"].(float64); ok {
			ts = int64(t)
			tsSet = true
		} else if _, ok := eventMap["timestamp"]; ok {
			tsSet = true
		} else if t, ok := eventMap["Timestamp"].(float64); ok {
			ts = int64(t)
			tsSet = true
		} else if _, ok := eventMap["Timestamp"]; ok {
			tsSet = true
		}
		msg := ""
		if m, ok := eventMap["message"].(string); ok {
			msg = m
		} else if m, ok := eventMap["Message"].(string); ok {
			msg = m
		}
		// Both members are required on the wire shape, so an element
		// missing one or both still yields its event: the seam's member
		// validation rejects the request rather than the element being
		// silently dropped from the batch.
		events = append(events, PutLogEvent{
			LogEntry: logsstore.LogEntry{
				Timestamp: ts,
				Message:   msg,
			},
			TimestampSet: tsSet,
		})
	}

	return events
}

// --- Subscription filter evaluation ---

// deliverSubscriptionEvents matches events against subscription filters and
// delivers matched events either via the event bus or, bus-less, directly
// through the shared dispatch.
func (s *LogsService) deliverSubscriptionEvents(store *logsstore.Store, region, logGroupName, logStreamName string, events []logsstore.LogEntry, transformed map[string]string) {
	filters, err := store.ListSubscriptionFilters(logGroupName, "")
	if err != nil || len(filters) == 0 {
		return
	}

	matcher := filterpattern.NewMatcher()

	for _, filter := range filters {
		// The criteria member gates which events the filter processes
		// ("specifies which log events should be processed by this
		// subscription filter based on system fields such as source
		// account and source region") — a filter whose criteria the
		// batch's system fields fail never evaluates.
		if !evalFieldSelectionCriteria(filter.FieldSelectionCriteria, region, s.accountID) {
			continue
		}
		var matched []logsstore.LogEntry
		for _, event := range events {
			// The delivered message follows the filter's transform
			// selection: a transformed-selecting subscription delivers
			// the transformed form (falling back to the original for
			// events the recipe left untransformed).
			message := transformedForFilters(event.Message,
				logsstore.TransformedMessageDigest(event.Timestamp, logStreamName, event.Message), transformed, filter.ApplyOnTransformedLogs)
			if filter.FilterPattern == "" || matcher.Matches(filter.FilterPattern, message) {
				matched = append(matched, logsstore.LogEntry{
					Timestamp:     event.Timestamp,
					Message:       message,
					IngestionTime: event.IngestionTime,
				})
			}
		}
		if len(matched) == 0 {
			continue
		}

		payload := s.buildSubscriptionPayload(filter, region, logGroupName, logStreamName, matched)
		compressed, err := compressJSON(payload)
		if err != nil {
			continue
		}

		if s.eventBus() != nil {
			evt := &eventbus.CloudWatchLogDeliveryEvent{
				LogGroup:       logGroupName,
				LogStream:      logStreamName,
				DestinationArn: filter.DestinationArn,
				Distribution:   filter.Distribution,
				Payload:        compressed,
			}
			evt.Region = region
			if err := s.eventBus().Publish(context.Background(), evt); err != nil {
				logs.Warn("Failed to publish log delivery event", logs.Err(err))
				// The event never reached the delivery handler: the batch
				// rides the retry window like any other failed delivery.
				s.retryFailedDelivery(store, region, filter.DestinationArn, logGroupName, logStreamName, filter.Distribution, compressed, err)
			}
		} else if err := s.dispatchSubscriptionDelivery(region, filter.DestinationArn, logGroupName, logStreamName, filter.Distribution, compressed); err != nil {
			s.retryFailedDelivery(store, region, filter.DestinationArn, logGroupName, logStreamName, filter.Distribution, compressed, err)
		}
	}
}

// --- Subscription payload building ---

func (s *LogsService) buildSubscriptionPayload(
	filter *logsstore.SubscriptionFilter,
	region, logGroupName, logStreamName string,
	events []logsstore.LogEntry,
) map[string]interface{} {
	logEvents := make([]map[string]interface{}, len(events))
	for i, e := range events {
		entry := map[string]interface{}{
			"id":        generateEventID(),
			"timestamp": e.Timestamp,
			"message":   e.Message,
		}
		// "A list of system fields to include in the log events sent to
		// the subscription destination. Valid values are @aws.account,
		// @aws.region, and @source.log."
		for _, field := range filter.EmitSystemFields {
			switch field {
			case "@aws.account":
				entry["@aws.account"] = s.accountID
			case "@aws.region":
				entry["@aws.region"] = region
			case "@source.log":
				entry["@source.log"] = logGroupName
			}
		}
		logEvents[i] = entry
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

// invokeLambda delivers compressed log data to a Lambda function. The
// documented Lambda invocation payload is the JSON envelope
// {"awslogs":{"data": BASE64ENCODED_GZIP_COMPRESSED_DATA}} — the base64
// wrap exists because a Lambda payload must be JSON, so the gzip bytes
// cannot ride in it raw (CloudWatch Logs User Guide, "Log group-level
// subscription filters", Lambda example). A non-nil return is a delivery
// failure the caller retries.
func (s *LogsService) invokeLambda(destArn string, compressed []byte) error {
	if s.eventBus() == nil {
		logs.Warn("Cannot deliver subscription payload to Lambda: event bus not configured",
			logs.String("destinationArn", destArn))
		return fmt.Errorf("event bus not configured")
	}

	functionName := arn.ExtractFunctionNameFromARN(destArn)
	invoker := s.eventBus().LambdaInvoker()
	if invoker == nil {
		logs.Warn("Cannot deliver subscription payload to Lambda: no Lambda invoker configured",
			logs.String("destinationArn", destArn))
		return fmt.Errorf("no lambda invoker configured")
	}
	encodedData := base64.StdEncoding.EncodeToString(compressed)

	payload, err := json.Marshal(map[string]interface{}{
		"awslogs": map[string]interface{}{"data": encodedData},
	})
	if err != nil {
		return fmt.Errorf("marshal lambda payload: %w", err)
	}

	if _, _, err := invoker.InvokeForGateway(context.Background(), functionName, payload); err != nil {
		logs.Warn("Failed to invoke Lambda for subscription filter delivery", logs.Err(err))
		return err
	}
	return nil
}

// putToKinesis delivers compressed log data to a Kinesis stream. The
// delivered blob is the gzipped subscription payload as-is — the
// documented consumer pipeline is `base64 -d | zcat` over the Kinesis
// API's base64 presentation of the record (CloudWatch Logs User Guide,
// "Log group-level subscription filters", Kinesis example) — one gzip
// layer, no JSON envelope or extra base64 wrap; the wrapper is the Lambda
// invocation envelope alone. A non-nil return is a delivery failure the
// caller retries: an unprobeable, shardless or closed stream is exactly
// the throttled-shape destination the retry window exists for.
func (s *LogsService) putToKinesis(destArn, logGroup, logStream, distribution string, compressed []byte) error {
	if s.eventBus() == nil {
		logs.Warn("Cannot deliver subscription payload to Kinesis: event bus not configured",
			logs.String("destinationArn", destArn))
		return fmt.Errorf("event bus not configured")
	}

	streamName := arn.ExtractStreamNameFromARN(destArn)
	_, _, destRegion, _, _ := arn.SplitARN(destArn)
	ctx := context.Background()
	invoker := s.eventBus().KinesisInvoker()
	if invoker == nil {
		logs.Warn("Cannot deliver subscription payload to Kinesis: no Kinesis invoker configured",
			logs.String("destinationArn", destArn))
		return fmt.Errorf("no kinesis invoker configured")
	}

	// The destination ARN addresses the stream in its own region: the shard
	// probe and the record write must resolve the same regional store, or a
	// cross-region destination fails the probe against the wrong region and
	// the delivery vanishes.
	shards, err := invoker.ListShards(ctx, destRegion, streamName)
	if err != nil {
		logs.Warn("Failed to list shards for subscription filter delivery to Kinesis",
			logs.String("stream", streamName),
			logs.String("region", destRegion),
			logs.Err(err))
		return fmt.Errorf("list shards: %w", err)
	}

	if len(shards) == 0 {
		logs.Warn("Subscription filter Kinesis destination has no shards",
			logs.String("stream", streamName),
			logs.String("region", destRegion))
		return fmt.Errorf("kinesis stream %s has no shards", streamName)
	}

	// The probe guards writability: the destination must expose an open
	// shard before the write is attempted. Hash placement, not the probe,
	// chooses the shard the record lands on.
	openShard := false
	for _, shard := range shards {
		if shard.SequenceNumberRangeEnd == "" {
			openShard = true
			break
		}
	}

	if !openShard {
		logs.Warn("Subscription filter Kinesis destination has no open shard",
			logs.String("stream", streamName),
			logs.String("region", destRegion))
		return fmt.Errorf("kinesis stream %s has no open shard", streamName)
	}

	// The platform's Kinesis storage carries a record's data as the JSON
	// wire form — the base64 text an SDK PutRecord arrives as and
	// GetRecords returns verbatim — so the delivery writes the same form:
	// the blob a client's base64 decoding yields is the documented gzipped
	// payload, one gzip layer with no JSON envelope (that envelope is the
	// Lambda invocation format alone).
	wireData := base64.StdEncoding.EncodeToString(compressed)
	if _, _, err := invoker.PutRecord(ctx, destRegion, streamName, subscriptionPartitionKey(logGroup, logStream, distribution), []byte(wireData)); err != nil {
		logs.Warn("Failed to deliver subscription filter log events to Kinesis", logs.Err(err))
		return err
	}
	return nil
}

// subscriptionPartitionKey derives the delivery's partition key from its
// log-group/log-stream identity and the filter's distribution. The default
// distribution is ByLogStream — the CloudWatch Logs User Guide's Kinesis
// example notes "By default, the stream filter distribution is by log
// stream", which causes throttling on hot streams and is why the guide
// recommends Random for an even spread — so ByLogStream keeps a stable,
// data-derived key (the SHA-256 hex digest of the group and stream names,
// bounded well inside the partition-key length limit; a log stream name
// alone may run to 512 characters): a retry of the same delivery keeps its
// shard, and distinct log streams group onto shards as the documented
// distribution describes. Random takes a fresh random key per delivery and
// lets hash placement spread the deliveries.
func subscriptionPartitionKey(logGroup, logStream, distribution string) string {
	if distribution == "Random" {
		return randomHexWithFallback(16, func() string {
			digest := sha256.Sum256([]byte(logGroup + "\x00" + logStream + "\x00" + "random"))
			return hex.EncodeToString(digest[:])
		})
	}
	digest := sha256.Sum256([]byte(logGroup + "\x00" + logStream))
	return hex.EncodeToString(digest[:])
}

// --- GetLogEvents / FilterLogEvents ---

// GetLogEvents retrieves log events from the specified CloudWatch Logs log stream.
func (s *LogsService) GetLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName, err := logGroupNameOrIdentifier(req.Parameters)
	if err != nil {
		return nil, err
	}
	// The Unmask parameter controls whether data-protection-masked fields
	// are returned in their original form. Data protection masking is not
	// implemented on this platform; the parameter is accepted but has no
	// effect, the same form the GetLogRecord/GetLogObject paths carry.
	_ = request.GetBoolParam(req.Parameters, "Unmask")

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
		outputEvents = append(outputEvents, formatLogEvent(e))
	}

	return map[string]interface{}{
		"events":            outputEvents,
		"nextForwardToken":  result.NextForwardToken,
		"nextBackwardToken": result.NextBackwardToken,
	}, nil
}

// FilterLogEvents filters log events from the specified CloudWatch Logs log group.
func (s *LogsService) FilterLogEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName, err := logGroupNameOrIdentifier(req.Parameters)
	if err != nil {
		return nil, err
	}
	// The Unmask parameter controls whether data-protection-masked fields
	// are returned in their original form. Data protection masking is not
	// implemented on this platform; the parameter is accepted but has no
	// effect, the same form the GetLogRecord/GetLogObject paths carry.
	_ = request.GetBoolParam(req.Parameters, "Unmask")

	// An explicitly empty array is a present member, not an omitted one:
	// the non-nil empty slice carries that distinction to the Core, where
	// the stream list's @length(min 1) trait rejects it.
	streamNames := request.GetStringList(req.Parameters, "LogStreamNames")
	if streamNames == nil && request.HasListParam(req.Parameters, "LogStreamNames") {
		streamNames = []string{}
	}

	input := FilterLogEventsInput{
		LogGroupName:      logGroupName,
		LogStreamNames:    streamNames,
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
		ev := formatLogEvent(e)
		// FilteredLogEvent carries an eventId (a digest over the event's
		// coordinates, stable across reads) and the logStreamName member
		// the OutputLogEvent shape does not have; both are minted here
		// alone.
		ev["eventId"] = filterEventID(logGroupName, e.LogStreamName, e.Timestamp, e.Message, e.Ordinal)
		ev["logStreamName"] = e.LogStreamName
		outputEvents = append(outputEvents, ev)
	}

	// searchedLogStreams is documented as unsupported since May 15, 2020:
	// the member is always present and always an empty list.
	resp := map[string]interface{}{
		"events":             outputEvents,
		"searchedLogStreams": []map[string]interface{}{},
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// --- Helpers ---

// filterEventID mints the deterministic event identifier the
// FilteredLogEvent shape carries: a hex digest over the event's
// coordinates plus the read engine's per-event ordinal, so the same
// stored event yields the same id on every read while byte-identical
// duplicate events — distinct events whose content coordinates alone
// would collide — keep distinct ids.
func filterEventID(group, stream string, ts int64, msg, ordinal string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d|%s|%s", group, stream, ts, msg, ordinal)))
	return hex.EncodeToString(sum[:])
}

func generateEventID() string {
	// The clock prefix keeps ids unique even when the random suffix
	// degrades to nothing (the shared minter's fallback shape).
	suffix := randomHexWithFallback(logEventIDSize, func() string { return "" })
	return fmt.Sprintf("%x%s", time.Now().UnixNano(), suffix)
}
