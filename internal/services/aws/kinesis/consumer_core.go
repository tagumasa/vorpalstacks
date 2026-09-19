package kinesis

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
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
// token; a following page exists exactly when the token is non-nil.
type ListStreamConsumersResult struct {
	Consumers []*kinesisstore.Consumer
	NextToken *string
}

// SubscribeToShardInput is the transport-agnostic input for SubscribeToShard.
type SubscribeToShardInput struct {
	ConsumerARN            string
	ShardId                string
	StartingPositionType   string
	StartingSequenceNumber string
	Timestamp              string
	DryRun                 bool
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
// create-time tags to the consumer ARN. Tags pass the same shared tag-set
// validation every tag write path applies before the consumer is created, so
// an invalid set leaves nothing behind. The per-stream quota is enforced
// inside the store's locked registration path.
func (s *KinesisService) registerStreamConsumerCore(reqCtx *request.RequestContext, input RegisterStreamConsumerInput) (ConsumerResult, error) {
	if input.StreamARN == "" || !validateConsumerName(input.ConsumerName) {
		return ConsumerResult{}, ErrInvalidArgument
	}

	if len(input.Tags) > 0 {
		if err := tags.ValidateTags(input.Tags); err != nil {
			return ConsumerResult{}, ErrInvalidArgument
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return ConsumerResult{}, err
	}

	var tagMap map[string]string
	if len(input.Tags) > 0 {
		tagMap = make(map[string]string, len(input.Tags))
		for _, t := range input.Tags {
			tagMap[t.Key] = t.Value
		}
	}

	consumer, err := store.RegisterStreamConsumer(input.StreamARN, input.ConsumerName, tagMap)
	if err != nil {
		return ConsumerResult{}, s.mapStoreError(err)
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
	// The token "unambiguously identifies the stream": the model's NextToken
	// documentation forbids the identification members from travelling with
	// it, so StreamCreationTimestamp cannot ride a token request.
	if input.NextToken != "" && input.StreamCreationTimestamp != "" {
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

	// The optional StreamCreationTimestamp disambiguates deleted+recreated
	// streams through the shared generation rule.
	if err := s.verifyStreamGeneration(input.StreamCreationTimestamp, stream); err != nil {
		return ListStreamConsumersResult{}, err
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

	// Decode the opaque NextToken — the base64 of the anchor consumer's
	// ARN and the token's own issuance time; one that cannot be decoded is
	// expired — the error ListStreamConsumers declares for it — never a
	// silent restart, and so is a token past the documented validity
	// window or one whose anchor consumer is no longer in the list: the
	// list changed under the client, and restarting from page one would
	// re-deliver entries.
	startOffset := 0
	if input.NextToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(input.NextToken)
		if err != nil {
			return ListStreamConsumersResult{}, ErrExpiredNextToken
		}
		sep := bytes.LastIndexByte(decoded, '#')
		if sep < 0 {
			return ListStreamConsumersResult{}, ErrExpiredNextToken
		}
		issuedAt, err := strconv.ParseInt(string(decoded[sep+1:]), 10, 64)
		if err != nil || time.Since(time.Unix(0, issuedAt)) > kinesisstore.NextTokenValidity {
			return ListStreamConsumersResult{}, ErrExpiredNextToken
		}
		startARN := string(decoded[:sep])
		found := false
		for i, c := range consumers {
			if c.ConsumerARN == startARN {
				startOffset = i + 1
				found = true
				break
			}
		}
		if !found {
			return ListStreamConsumersResult{}, ErrExpiredNextToken
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
		encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s#%d", page[len(page)-1].ConsumerARN, time.Now().UnixNano())))
		nextToken = &encoded
	}

	return ListStreamConsumersResult{
		Consumers: page,
		NextToken: nextToken,
	}, nil
}

// subscribeToShardCore opens a SubscribeToShard event stream: it resolves the
// consumer and its stream, validates the starting position, creates the
// starting iterator and serves records over the returned event-stream
// reader until the shard closes, the subscription's documented lifetime
// expires, the request context is cancelled, or a newer subscription takes
// the consumer-shard pair over.
func (s *KinesisService) subscribeToShardCore(ctx context.Context, reqCtx *request.RequestContext, input SubscribeToShardInput) (*SubscribeToShardEventStream, error) {
	if input.ConsumerARN == "" || input.ShardId == "" || !validateShardId(input.ShardId) {
		return nil, ErrInvalidArgument
	}
	// The StartingPosition member carries the same pairing rules as the
	// shard-iterator types: the type must be one of the model's five
	// values, a sequence-number position must satisfy the SequenceNumber
	// pattern and travel with AT/AFTER_SEQUENCE_NUMBER, and a timestamp
	// position with AT_TIMESTAMP.
	if !validateStartingPosition(input.StartingPositionType, input.StartingSequenceNumber, input.Timestamp != "") {
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
		t, err := parseTimestampMember(input.Timestamp)
		if err != nil {
			return nil, err
		}
		timestamp = &t
	}

	// A DryRun request runs every check — the shard must resolve — then
	// answers the model's dedicated error without opening the subscription.
	if input.DryRun {
		if _, err := store.GetShard(stream.StreamName, input.ShardId); err != nil {
			return nil, s.mapStoreError(err)
		}
		return nil, ErrDryRunOperation
	}

	// The subscription registry enforces the documented pair rules: a
	// repeat subscription for the same consumer and shard within the
	// takeover window answers ResourceInUseException; at or beyond it the
	// new subscription takes over and the previous connection expires.
	key := subscriptionKey(input.ConsumerARN, input.ShardId)
	subCtx, sub, ok := s.subscriptions.begin(key, ctx)
	if !ok {
		return nil, ErrResourceInUse
	}

	iterator, err := store.CreateShardIterator(stream.StreamName, input.ShardId, input.StartingPositionType, input.StartingSequenceNumber, timestamp)
	if err != nil {
		// The call never succeeded, so no tombstone: the takeover window
		// the documentation states runs from a successful call, and a
		// corrected retry is answerable at once.
		s.subscriptions.discard(key, sub)
		return nil, s.mapStoreError(err)
	}

	pr, pw := io.Pipe()
	go s.runSubscribePump(subCtx, pw, key, sub, store, stream, input, iterator)

	return &SubscribeToShardEventStream{reader: pr}, nil
}

// writeCtx runs one event write, abandoning it when ctx ends or the
// subscription's lifetime elapses: io.Pipe writes block until the peer
// reads, so a client that stops reading must not hold the pump — and a
// pump parked on such a write must still observe its documented lifetime,
// which the loop select alone cannot deliver. Closing the pipe unblocks
// the in-flight write and ends the stream for the reader as a clean
// expiry — the documented takeover wording ("the previous connection
// expires"); the pump then exits through its deferred cleanup.
func writeCtx(ctx context.Context, pw *io.PipeWriter, expire <-chan time.Time, write func() error) error {
	done := make(chan error, 1)
	go func() {
		defer func() {
			// recover only works called directly by the deferred function —
			// through a helper it returns nil and the panic continues — so
			// the containment lives here. The buffered channel still
			// receives, so the waiting side ends the stream through its own
			// cleanup instead of hanging on a delivery that will never come.
			if r := recover(); r != nil {
				logs.Error("Event write panicked", logs.Any("panic", r), logs.Any("stack", string(debug.Stack())))
				done <- fmt.Errorf("event write panicked: %v", r)
			}
		}()
		done <- write()
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		pw.Close()
		<-done
		return ctx.Err()
	case <-expire:
		pw.Close()
		<-done
		return errSubscriptionExpired
	}
}

// errSubscriptionExpired reports a write abandoned by the subscription's
// documented lifetime; the pump's deferred cleanup ends the stream.
var errSubscriptionExpired = errors.New("subscription lifetime elapsed")

// runSubscribePump serves one SubscribeToShard event stream. Records carry
// the stream's encryption type (the same value the GetRecords path
// reports — a KMS-encrypted stream must never report NONE here), read
// errors surface under their mapped identity instead of a blanket
// not-found, and a panic ends the stream loudly — logged with its stack
// and an InternalFailure error event when the writer is still usable —
// never the silent swallow.
func (s *KinesisService) runSubscribePump(ctx context.Context, pw *io.PipeWriter, key string, sub *subscription, store *kinesisstore.KinesisStore, stream *kinesisstore.Stream, input SubscribeToShardInput, iterator *kinesisstore.ShardIterator) {
	defer s.subscriptions.end(key, sub)
	defer pw.Close()

	writer := NewSubscribeToShardEventStreamWriter(pw)
	defer func() {
		// recover only works called directly by the deferred function —
		// through a helper it returns nil and the panic continues — so
		// the containment lives here.
		if r := recover(); r != nil {
			logs.Error("SubscribeToShard pump panicked", logs.Any("panic", r), logs.Any("stack", string(debug.Stack())))
			func() {
				defer func() { _ = recover() }()
				_ = writeCtx(ctx, pw, nil, func() error {
					return writer.WriteErrorEvent("InternalFailureException", "internal failure")
				})
			}()
		}
	}()

	// The lifetime timer arms every interactive event write (the loop
	// select alone cannot reach a pump parked on a stalled client) and the
	// loop itself; the closing end event below stays context-bound — the
	// teardown framing belongs to the pump's exit, past the lifetime.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	lifetime := time.NewTimer(s.subscriptions.lifetime())
	defer lifetime.Stop()

	if err := writeCtx(ctx, pw, lifetime.C, writer.WriteInitialResponse); err != nil {
		logs.Error("Failed to write initial response", logs.Err(err))
		return
	}

	encryptionType := resolveEncryptionType(stream)

	includeStart := input.StartingPositionType == "AT_SEQUENCE_NUMBER"
	lastSeqNum := iterator.SequenceNumber
	lastEventTime := time.Now()
	firstIteration := true

	// heartbeatInterval is the maximum gap between events before a
	// heartbeat event is sent to keep the connection alive. AWS does
	// not document the exact interval; 15 seconds is a provisional
	// value that stays well within typical LB/proxy idle timeouts
	// (60 s) while minimising idle traffic.
	const heartbeatInterval = 15 * time.Second

	for {
		page, err := s.readShardPage(store, stream.StreamName, input.ShardId, lastSeqNum, kinesisstore.SubscribePollPageSize, includeStart)
		if err != nil {
			// The error event carries the mapped identity: a genuine
			// not-found reports ResourceNotFoundException; anything
			// else reports its own code, never a blanket not-found.
			mapped := s.mapStoreError(err)
			code := "InternalFailureException"
			var awsErr *awserrors.AWSError
			if errors.As(mapped, &awsErr) {
				code = awsErr.Code
			}
			if writeErr := writeCtx(ctx, pw, lifetime.C, func() error {
				return writer.WriteErrorEvent(code, mapped.Error())
			}); writeErr != nil {
				logs.Error("Failed to write error event", logs.Err(writeErr))
			}
			return
		}
		records := page.Records
		newSeqNum := page.LastSeqNum
		shardClosed := page.ShardClosed
		millisBehindLatest := page.MillisBehindLatest
		// After the first iteration, use strict after-sequence-number semantics
		includeStart = false
		// The advance position moves the cursor on every page — aged-out
		// records included — so the pump cannot loop scanning the same
		// expired window.
		lastSeqNum = newSeqNum

		// Build child shards list when the shard is closed
		var childShards []interface{}
		if shardClosed {
			childShards = formatChildShards(page.ChildShards)
		}

		// Send event when there are records or when the shard closes
		if len(records) > 0 {
			if err := writeCtx(ctx, pw, lifetime.C, func() error {
				return writer.WriteSubscribeToShardEvent(records, newSeqNum, millisBehindLatest, childShards, encryptionType)
			}); err != nil {
				return
			}
			lastEventTime = time.Now()
		} else if shardClosed && len(childShards) > 0 {
			if err := writeCtx(ctx, pw, lifetime.C, func() error {
				return writer.WriteSubscribeToShardEvent(nil, lastSeqNum, 0, childShards, encryptionType)
			}); err != nil {
				return
			}
			lastEventTime = time.Now()
		} else if firstIteration || time.Since(lastEventTime) >= heartbeatInterval {
			// Heartbeat: an event with no records keeps the connection
			// alive during idle periods and hands the client its
			// checkpoint — the ContinuationSequenceNumber documentation
			// exists to "capture your shard progress even when no data is
			// written to the shard". The first iteration fires
			// immediately so an idle shard's first event does not lag the
			// heartbeat interval; AWS Kinesis Data Streams sends periodic
			// heartbeat events for the same purpose; the exact interval
			// is undocumented.
			if err := writeCtx(ctx, pw, lifetime.C, func() error {
				return writer.WriteSubscribeToShardEvent(nil, lastSeqNum, millisBehindLatest, nil, encryptionType)
			}); err != nil {
				return
			}
			lastEventTime = time.Now()
		}
		firstIteration = false

		// Shard closed — send End event and stop
		if shardClosed {
			break
		}

		expired := false
		select {
		case <-ctx.Done():
			return
		case <-lifetime.C:
			// The documented subscription lifetime elapsed; the client
			// renews by calling SubscribeToShard again.
			expired = true
		case <-ticker.C:
		}
		if expired {
			break
		}
	}

	if err := writeCtx(ctx, pw, nil, writer.WriteEndEvent); err != nil {
		logs.Error("Failed to write end event", logs.Err(err))
	}
}
