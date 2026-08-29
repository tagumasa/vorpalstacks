package kinesis

import (
	"context"
	"encoding/base64"
	"io"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// RegisterStreamConsumerInput is the transport-agnostic input for
// RegisterStreamConsumer.
type RegisterStreamConsumerInput struct {
	StreamARN    string
	ConsumerName string
	Tags         []tags.Tag
}

// ConsumerResult carries a resolved stream consumer.
type ConsumerResult struct {
	Consumer *kinesisstore.Consumer
}

// DeregisterStreamConsumerInput is the transport-agnostic input for
// DeregisterStreamConsumer.
type DeregisterStreamConsumerInput struct {
	StreamARN    string
	ConsumerARN  string
	ConsumerName string
}

// DescribeStreamConsumerInput is the transport-agnostic input for
// DescribeStreamConsumer.
type DescribeStreamConsumerInput struct {
	StreamARN    string
	ConsumerARN  string
	ConsumerName string
}

// ListStreamConsumersInput is the transport-agnostic input for
// ListStreamConsumers.
type ListStreamConsumersInput struct {
	StreamARN               string
	StreamCreationTimestamp string
	MaxResults              int
	HasMaxResults           bool
	NextToken               string
}

// ListStreamConsumersResult carries a page of consumers plus the resumption
// token.
type ListStreamConsumersResult struct {
	Consumers []*kinesisstore.Consumer
	HasMore   bool
	NextToken *string
}

// SubscribeToShardInput is the transport-agnostic input for SubscribeToShard.
type SubscribeToShardInput struct {
	ConsumerARN            string
	ShardId                string
	StartingPositionType   string
	StartingSequenceNumber string
	Timestamp              string
}

// The MaxResults input window and the effective page follow two different
// documented rules: the ListStreamConsumersInputLimit range trait accepts
// 1-10000 on the wire, while the MaxResults member documentation fixes the
// default page at 100 and returns at most 100 results per call.
const (
	DefaultListStreamConsumersResults = 100
	MaxListStreamConsumersResults     = 100
)

// registerStreamConsumerCore registers a consumer on a stream and applies the
// create-time tags to the consumer ARN.
func (s *KinesisService) registerStreamConsumerCore(reqCtx *request.RequestContext, input RegisterStreamConsumerInput) (ConsumerResult, error) {
	if input.StreamARN == "" || !validateConsumerName(input.ConsumerName) {
		return ConsumerResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return ConsumerResult{}, err
	}

	consumer, err := store.RegisterStreamConsumer(input.StreamARN, input.ConsumerName)
	if err != nil {
		return ConsumerResult{}, s.mapStoreError(err)
	}

	if len(input.Tags) > 0 {
		tagMap := make(map[string]string, len(input.Tags))
		for _, t := range input.Tags {
			tagMap[t.Key] = t.Value
		}
		if err := store.Tag(consumer.ConsumerARN, tagMap); err != nil {
			return ConsumerResult{}, s.mapStoreError(err)
		}
	}

	return ConsumerResult{Consumer: consumer}, nil
}

// deregisterStreamConsumerCore deregisters a consumer identified by ARN or by
// the stream-ARN + consumer-name pair.
func (s *KinesisService) deregisterStreamConsumerCore(reqCtx *request.RequestContext, input DeregisterStreamConsumerInput) error {
	if input.ConsumerARN == "" && input.ConsumerName == "" {
		return ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	consumerARN := input.ConsumerARN
	if input.ConsumerName != "" && input.StreamARN != "" {
		consumer, err := store.GetStreamConsumerByName(input.StreamARN, input.ConsumerName)
		if err != nil {
			return s.mapStoreError(err)
		}
		consumerARN = consumer.ConsumerARN
	}

	if err := store.DeregisterStreamConsumer(consumerARN); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// describeStreamConsumerCore resolves a consumer by ARN or by the stream-ARN +
// consumer-name pair.
func (s *KinesisService) describeStreamConsumerCore(reqCtx *request.RequestContext, input DescribeStreamConsumerInput) (ConsumerResult, error) {
	if input.ConsumerARN == "" && input.ConsumerName == "" {
		return ConsumerResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return ConsumerResult{}, err
	}

	var consumer *kinesisstore.Consumer
	if input.ConsumerARN != "" {
		consumer, err = store.GetStreamConsumer(input.ConsumerARN)
	} else if input.ConsumerName != "" && input.StreamARN != "" {
		consumer, err = store.GetStreamConsumerByName(input.StreamARN, input.ConsumerName)
	} else {
		return ConsumerResult{}, ErrInvalidArgument
	}

	if err != nil {
		return ConsumerResult{}, s.mapStoreError(err)
	}

	return ConsumerResult{Consumer: consumer}, nil
}

// listStreamConsumersCore lists the consumers of a stream with the documented
// MaxResults window and opaque ARN-based resumption tokens.
func (s *KinesisService) listStreamConsumersCore(reqCtx *request.RequestContext, input ListStreamConsumersInput) (ListStreamConsumersResult, error) {
	if input.StreamARN == "" {
		return ListStreamConsumersResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return ListStreamConsumersResult{}, err
	}

	stream, err := store.GetStreamByARN(input.StreamARN)
	if err != nil {
		return ListStreamConsumersResult{}, s.mapStoreError(err)
	}

	// Parse optional StreamCreationTimestamp — used to disambiguate streams
	// with the same name (deleted + recreated). Verify the stream identity.
	if input.StreamCreationTimestamp != "" {
		if unixTs, err := strconv.ParseFloat(input.StreamCreationTimestamp, 64); err == nil {
			// Compare at second precision: stream.CreatedAt has nanosecond
			// resolution but the client timestamp is epoch seconds.
			if stream.CreatedAt.Unix() != int64(unixTs) {
				return ListStreamConsumersResult{}, s.mapStoreError(kinesisstore.ErrStreamNotFound)
			}
		}
	}

	consumers, err := store.ListStreamConsumers(stream.StreamName)
	if err != nil {
		return ListStreamConsumersResult{}, s.mapStoreError(err)
	}

	// The range trait window is enforced on the provided value, then the
	// effective page applies the documented cap.
	maxResults := input.MaxResults
	if input.HasMaxResults {
		if !validateListStreamConsumersLimit(maxResults) {
			return ListStreamConsumersResult{}, ErrInvalidArgument
		}
		if maxResults > MaxListStreamConsumersResults {
			maxResults = MaxListStreamConsumersResults
		}
	} else {
		maxResults = DefaultListStreamConsumersResults
	}

	// Decode opaque NextToken (base64-encoded consumer ARN for resumption)
	startOffset := 0
	if input.NextToken != "" {
		if decoded, err := base64.StdEncoding.DecodeString(input.NextToken); err == nil {
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

	var nextToken *string
	if hasMore && len(page) > 0 {
		encoded := base64.StdEncoding.EncodeToString([]byte(page[len(page)-1].ConsumerARN))
		nextToken = &encoded
	}

	return ListStreamConsumersResult{
		Consumers: page,
		HasMore:   hasMore,
		NextToken: nextToken,
	}, nil
}

// subscribeToShardCore opens a SubscribeToShard event stream: it resolves the
// consumer and its stream, creates the starting iterator and serves records
// over the returned event-stream reader until the shard closes or the request
// context is cancelled.
func (s *KinesisService) subscribeToShardCore(ctx context.Context, reqCtx *request.RequestContext, input SubscribeToShardInput) (*SubscribeToShardEventStream, error) {
	if input.ConsumerARN == "" || input.ShardId == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	consumer, err := store.GetStreamConsumer(input.ConsumerARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream, err := store.GetStreamByARN(consumer.StreamARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var timestamp *time.Time
	if input.Timestamp != "" {
		if unixTs, err := strconv.ParseInt(input.Timestamp, 10, 64); err == nil {
			t := time.Unix(unixTs, 0).UTC()
			timestamp = &t
		}
	}

	iterator, err := store.CreateShardIterator(stream.StreamName, input.ShardId, input.StartingPositionType, input.StartingSequenceNumber, timestamp)
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

		includeStart := input.StartingPositionType == "AT_SEQUENCE_NUMBER"
		lastSeqNum := iterator.SequenceNumber
		lastEventTime := time.Now()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// heartbeatInterval is the maximum gap between events before a
		// heartbeat event is sent to keep the connection alive. AWS does
		// not document the exact interval; 15 seconds is a provisional
		// value that stays well within typical LB/proxy idle timeouts
		// (60 s) while minimising idle traffic.
		const heartbeatInterval = 15 * time.Second

		for {
			records, newSeqNum, err := store.GetRecords(stream.StreamName, input.ShardId, lastSeqNum, 1000, includeStart)
			if err != nil {
				if writeErr := writer.WriteResourceNotFoundException(err.Error()); writeErr != nil {
					logs.Error("Failed to write error event", logs.Err(writeErr))
				}
				return
			}
			// After the first iteration, use strict after-sequence-number semantics
			includeStart = false

			// Check whether the shard has been closed (split or merged)
			shard, shardErr := store.GetShard(stream.StreamName, input.ShardId)
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
				children, _ := store.GetChildShards(stream.StreamName, input.ShardId)
				childShards = formatChildShards(children)
			}

			// Send event when there are records or when the shard closes
			if len(records) > 0 {
				if err := writer.WriteSubscribeToShardEvent(records, newSeqNum, millisBehindLatest, childShards); err != nil {
					return
				}
				lastSeqNum = newSeqNum
				lastEventTime = time.Now()
			} else if shardClosed && len(childShards) > 0 {
				if err := writer.WriteSubscribeToShardEvent(nil, lastSeqNum, 0, childShards); err != nil {
					return
				}
				lastEventTime = time.Now()
			} else if time.Since(lastEventTime) >= heartbeatInterval {
				// Heartbeat: send an event with no records to keep the
				// connection alive during idle periods. AWS Kinesis Data
				// Streams sends periodic heartbeat events for the same
				// purpose; the exact interval is undocumented.
				if err := writer.WriteSubscribeToShardEvent(nil, lastSeqNum, millisBehindLatest, nil); err != nil {
					return
				}
				lastEventTime = time.Now()
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
