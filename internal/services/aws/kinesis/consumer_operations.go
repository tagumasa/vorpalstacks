package kinesis

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// SubscribeToShardEventStream represents an event stream for SubscribeToShard operations.
type SubscribeToShardEventStream struct {
	reader io.Reader
}

// GetStream returns the reader for the stream.
func (e *SubscribeToShardEventStream) GetStream() io.Reader {
	return e.reader
}

// GetStreamHeaders returns the HTTP headers for the stream.
func (e *SubscribeToShardEventStream) GetStreamHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/vnd.amazon.eventstream")
	return headers
}

// RegisterStreamConsumer registers a consumer for a Kinesis stream.
func (s *KinesisService) RegisterStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	consumerName := request.GetParamLowerFirst(req.Parameters, "ConsumerName")

	if streamARN == "" || consumerName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	consumer, err := store.RegisterStreamConsumer(streamARN, consumerName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// M14: Apply tags if provided in RegisterStreamConsumer request
	if tagList := tags.ParseTags(req.Parameters, "Tags"); len(tagList) > 0 {
		tagMap := make(map[string]string, len(tagList))
		for _, t := range tagList {
			tagMap[t.Key] = t.Value
		}
		if err := store.Tag(consumer.ConsumerARN, tagMap); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	return map[string]interface{}{
		"Consumer": formatConsumer(consumer),
	}, nil
}

// DeregisterStreamConsumer deregisters a consumer from a Kinesis stream.
func (s *KinesisService) DeregisterStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	consumerARN := request.GetParamLowerFirst(req.Parameters, "ConsumerARN")
	consumerName := request.GetParamLowerFirst(req.Parameters, "ConsumerName")

	if consumerARN == "" && consumerName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if consumerName != "" && streamARN != "" {
		consumer, err := store.GetStreamConsumerByName(streamARN, consumerName)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		consumerARN = consumer.ConsumerARN
	}

	if err := store.DeregisterStreamConsumer(consumerARN); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeStreamConsumer returns details about a Kinesis stream consumer.
func (s *KinesisService) DescribeStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	consumerARN := request.GetParamLowerFirst(req.Parameters, "ConsumerARN")
	consumerName := request.GetParamLowerFirst(req.Parameters, "ConsumerName")

	if consumerARN == "" && consumerName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var consumer *kinesisstore.Consumer
	if consumerARN != "" {
		consumer, err = store.GetStreamConsumer(consumerARN)
	} else if consumerName != "" && streamARN != "" {
		consumer, err = store.GetStreamConsumerByName(streamARN, consumerName)
	} else {
		return nil, ErrInvalidArgument
	}

	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ConsumerDescription": formatConsumer(consumer),
	}, nil
}

// ListStreamConsumers lists consumers of a Kinesis stream.
func (s *KinesisService) ListStreamConsumers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	if streamARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStreamByARN(streamARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// Parse optional StreamCreationTimestamp — used to disambiguate streams
	// with the same name (deleted + recreated). Verify the stream identity.
	if tsStr := request.GetParamLowerFirst(req.Parameters, "StreamCreationTimestamp"); tsStr != "" {
		if unixTs, err := strconv.ParseFloat(tsStr, 64); err == nil {
			// Compare at second precision: stream.CreatedAt has nanosecond
			// resolution but the client timestamp is epoch seconds.
			if stream.CreatedAt.Unix() != int64(unixTs) {
				return nil, s.mapStoreError(kinesisstore.ErrStreamNotFound)
			}
		}
	}

	consumers, err := store.ListStreamConsumers(stream.StreamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// Pagination: parse MaxResults and NextToken
	maxResults := 0
	if _, ok := req.Parameters["MaxResults"]; ok {
		maxResults = request.GetIntParam(req.Parameters, "MaxResults")
	}
	// MaxResults: Smithy range 1-10000. Reject out-of-range values.
	if _, ok := req.Parameters["MaxResults"]; ok {
		if maxResults < 1 || maxResults > 10000 {
			return nil, ErrInvalidArgument
		}
	} else {
		maxResults = 10000
	}

	// Decode opaque NextToken (base64-encoded consumer ARN for resumption)
	startOffset := 0
	if encoded := request.GetStringParam(req.Parameters, "NextToken"); encoded != "" {
		if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
			startARN := string(decoded)
			for i, c := range consumers {
				if c.ConsumerARN == startARN {
					startOffset = i + 1
					break
				}
			}
		}
	}

	// Apply offset and limit
	end := startOffset + maxResults
	hasMore := false
	if end < len(consumers) {
		hasMore = true
	} else {
		end = len(consumers)
	}
	page := consumers[startOffset:end]

	formattedConsumers := make([]map[string]interface{}, len(page))
	for i, c := range page {
		formattedConsumers[i] = formatConsumer(c)
	}

	result := map[string]interface{}{
		"Consumers":        formattedConsumers,
		"HasMoreConsumers": hasMore,
	}
	if hasMore && len(page) > 0 {
		result["NextToken"] = base64.StdEncoding.EncodeToString([]byte(page[len(page)-1].ConsumerARN))
	} else {
		result["NextToken"] = nil
	}

	return result, nil
}

// SubscribeToShard subscribes a consumer to receive records from a Kinesis shard.
func (s *KinesisService) SubscribeToShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	consumerARN := request.GetParamLowerFirst(req.Parameters, "ConsumerARN")
	shardID := request.GetParamLowerFirst(req.Parameters, "ShardId")

	if consumerARN == "" || shardID == "" {
		return nil, ErrInvalidArgument
	}

	startingPosition := ""
	seqNum := ""
	tsStr := ""
	if sp := request.GetMapParam(req.Parameters, "StartingPosition"); sp != nil {
		startingPosition = request.GetStringParam(sp, "Type")
		seqNum = request.GetStringParam(sp, "SequenceNumber")
		tsStr = request.GetStringParam(sp, "Timestamp")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	consumer, err := store.GetStreamConsumer(consumerARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream, err := store.GetStreamByARN(consumer.StreamARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var timestamp *time.Time
	if tsStr != "" {
		if unixTs, err := strconv.ParseInt(tsStr, 10, 64); err == nil {
			t := time.Unix(unixTs, 0).UTC()
			timestamp = &t
		}
	}

	iterator, err := store.CreateShardIterator(stream.StreamName, shardID, startingPosition, seqNum, timestamp)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		defer func() { recover() }()

		writer := NewSubscribeToShardEventStreamWriter(pw)

		if err := writer.WriteInitialResponse(); err != nil {
			logs.Error("Failed to write initial response", logs.Err(err))
			return
		}

		includeStart := startingPosition == "AT_SEQUENCE_NUMBER"
		lastSeqNum := iterator.SequenceNumber
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			records, newSeqNum, err := store.GetRecords(stream.StreamName, shardID, lastSeqNum, 1000, includeStart)
			if err != nil {
				if writeErr := writer.WriteResourceNotFoundException(err.Error()); writeErr != nil {
					logs.Error("Failed to write error event", logs.Err(writeErr))
				}
				return
			}
			// After the first iteration, use strict after-sequence-number semantics
			includeStart = false

			// Check whether the shard has been closed (split or merged)
			shard, shardErr := store.GetShard(stream.StreamName, shardID)
			if shardErr != nil {
				logs.Error("SubscribeToShard: failed to get shard for closure check", logs.Err(shardErr))
				break
			}
			shardClosed := shard != nil && shard.SequenceNumberRange != nil && shard.SequenceNumberRange.EndingSequenceNumber != ""

			// Calculate MillisBehindLatest from the last record's arrival time
			var millisBehindLatest int64
			if len(records) > 0 {
				last := records[len(records)-1]
				millisBehindLatest = time.Since(last.ApproximateArrivalTimestamp).Milliseconds()
				if millisBehindLatest < 0 {
					millisBehindLatest = 0
				}
			}

			// Build child shards list when the shard is closed
			var childShards []interface{}
			if shardClosed {
				children, _ := store.GetChildShards(stream.StreamName, shardID)
				childShards = formatChildShards(children)
			}

			// Send event when there are records or when the shard closes
			if len(records) > 0 {
				if err := writer.WriteSubscribeToShardEvent(records, newSeqNum, millisBehindLatest, childShards); err != nil {
					return
				}
				lastSeqNum = newSeqNum
			} else if shardClosed && len(childShards) > 0 {
				if err := writer.WriteSubscribeToShardEvent(nil, lastSeqNum, 0, childShards); err != nil {
					return
				}
			}

			// Shard closed — send End event and stop
			if shardClosed {
				break
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}

		if err := writer.WriteEndEvent(); err != nil {
			logs.Error("Failed to write end event", logs.Err(err))
		}
	}()

	return &SubscribeToShardEventStream{reader: pr}, nil
}
