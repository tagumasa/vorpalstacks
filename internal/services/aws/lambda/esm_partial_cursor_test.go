package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// A partial batch response terminates the contiguous consumed prefix: the
// checkpoint rule is "Lambda uses the record with the lowest sequence
// number as the checkpoint. Lambda then retries all records starting from
// that checkpoint." These tests pin that no later chunk (buffered flush),
// later batch (parallel fold) or bisected half may move the checkpoint
// past a reported failure.

func TestPrefixOutcomeClampsAtPartialBatch(t *testing.T) {
	cases := []struct {
		name          string
		outcomes      []batchOutcome
		wantConsumed  string
		wantDelivered bool
		wantDiscarded bool
		wantReported  int
		wantFailure   bool
	}{
		{
			name: "full successes take the last cursor",
			outcomes: []batchOutcome{
				{lastConsumed: "3", delivered: true},
				{lastConsumed: "6", delivered: true},
			},
			wantConsumed:  "6",
			wantDelivered: true,
		},
		{
			name: "a partial batch stops the prefix at its cursor",
			outcomes: []batchOutcome{
				{lastConsumed: "3", delivered: true},
				{lastConsumed: "1", delivered: true, short: true, reported: 1},
				{lastConsumed: "6", delivered: true},
			},
			wantConsumed:  "1",
			wantDelivered: true,
			wantReported:  1,
		},
		{
			name: "a partial batch consuming nothing keeps the earlier cursor",
			outcomes: []batchOutcome{
				{lastConsumed: "3", delivered: true},
				{lastConsumed: "", delivered: true, short: true, reported: 2},
			},
			wantConsumed:  "3",
			wantDelivered: true,
			wantReported:  2,
		},
		{
			name: "a failure after the clamp keeps the clamped cursor",
			outcomes: []batchOutcome{
				{lastConsumed: "1", delivered: true, short: true, reported: 1},
				{lastConsumed: "5", delivered: true},
				{err: context.Canceled},
			},
			wantConsumed:  "1",
			wantDelivered: true,
			wantFailure:   true,
			wantReported:  1,
		},
		{
			name: "flags accumulate across the clamp",
			outcomes: []batchOutcome{
				{lastConsumed: "1", delivered: true, short: true, reported: 1},
				{lastConsumed: "5", discarded: true},
			},
			wantConsumed:  "1",
			wantDelivered: true,
			wantDiscarded: true,
			wantReported:  1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lastConsumed, delivered, discarded, reported, failure := prefixOutcome(tc.outcomes)
			if lastConsumed != tc.wantConsumed {
				t.Fatalf("lastConsumed = %q, want %q", lastConsumed, tc.wantConsumed)
			}
			if delivered != tc.wantDelivered {
				t.Fatalf("delivered = %v, want %v", delivered, tc.wantDelivered)
			}
			if discarded != tc.wantDiscarded {
				t.Fatalf("discarded = %v, want %v", discarded, tc.wantDiscarded)
			}
			if reported != tc.wantReported {
				t.Fatalf("reported = %d, want %d", reported, tc.wantReported)
			}
			if (failure != nil) != tc.wantFailure {
				t.Fatalf("failure = %v, want present %v", failure, tc.wantFailure)
			}
		})
	}
}

// TestDynamoDBBufferedFlushStopsAtPartialChunk drives the buffered
// (batching window) delivery with two chunks: the first chunk ends on a
// partial batch response, so the flush must stop with its cursor — the
// second chunk is not invoked and the next cycle re-reads from the
// partial checkpoint.
func TestDynamoDBBufferedFlushStopsAtPartialChunk(t *testing.T) {
	p, stream, mapping, _ := newDDBAnchorPoller(t, "TRIM_HORIZON")
	mapping.FunctionResponseTypes = []string{"ReportBatchItemFailures"}
	mapping.BatchSize = 2
	mapping.ParallelizationFactor = 2
	mapping.MaximumBatchingWindowInSeconds = 1
	reporter := &scriptedBatchReporter{failSeq: "2", failCalls: 1}
	p.invoke = reporter.invoke

	for i := int64(1); i <= 4; i++ {
		stream.publish(i)
	}

	p.processDynamoDBStreamsMapping(context.Background(), mapping)
	payloads := reporter.snapshot()
	if got := deliveryCount(payloads, "1"); got != 1 {
		t.Fatalf("the record before the failure must be delivered once, got %d deliveries", got)
	}
	if got := deliveryCount(payloads, "2"); got != 1 {
		t.Fatalf("the failed record must be delivered once before its retry, got %d deliveries", got)
	}
	if got := deliveryCount(payloads, "3"); got != 0 {
		t.Fatalf("the chunk after a partial response must not be invoked in the same flush, got %d deliveries", got)
	}
	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "1" {
		t.Fatalf("the flush checkpoint must stop at the partial cursor, got %q", got)
	}

	// The next cycle re-reads from the partial checkpoint: the reported
	// record is retried and, with the report now clean, the buffer drains.
	p.processDynamoDBStreamsMapping(context.Background(), mapping)
	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "4" {
		t.Fatalf("a later clean flush must advance the checkpoint, got %q", got)
	}
	if got := deliveryCount(reporter.snapshot(), "2"); got != 2 {
		t.Fatalf("the failed record must be re-delivered, got %d deliveries", got)
	}
}

// TestKinesisParallelFoldClampsAtPartialBatch drives the unbuffered
// Kinesis path with ParallelizationFactor 2: the first batch ends on a
// partial batch response and the second batch succeeds concurrently — the
// fold must clamp the checkpoint to the partial cursor and the cycle must
// report the partial response instead of a clean success.
func TestKinesisParallelFoldClampsAtPartialBatch(t *testing.T) {
	const (
		streamARN = "arn:aws:kinesis:us-east-1:123456789012:stream/pfold-stream"
		funcARN   = "arn:aws:lambda:us-east-1:123456789012:function:pfold-fn"
	)
	esmStore := lambdastore.NewEventSourceStore(&memStorage{bucket: newMemBucket()}, "123456789012", "us-east-1")
	created, err := esmStore.Create(&lambdastore.EventSourceMapping{
		EventSourceArn:        streamARN,
		FunctionArn:           funcARN,
		StartingPosition:      "TRIM_HORIZON",
		BatchSize:             2,
		ParallelizationFactor: 2,
		State:                 "Enabled",
		FunctionResponseTypes: []string{"ReportBatchItemFailures"},
	})
	if err != nil {
		t.Fatalf("create mapping: %v", err)
	}

	stream := &scriptedKinesisStream{
		shards: []eventbus.ShardInfo{{ShardID: "shard-0"}},
		floor:  "0",
	}
	for _, seq := range []string{"1", "2", "3", "4"} {
		stream.publish(seq)
	}
	reporter := &scriptedBatchReporter{failSeq: "2", failCalls: 1, needle: `"sequenceNumber"`}
	p := &esmPoller{
		bus:       &fakeKinesisBus{invoker: stream},
		invoke:    reporter.invoke,
		kinesisCP: make(map[string]string),
		buffers:   make(map[string]*streamBuffer),
		esmStore:  esmStore,
		// The rendered records carry the invoke identity ARN, which is
		// built from the service account ID.
		lambdaSvc: &LambdaService{accountID: "123456789012"},
	}

	p.processKinesisMapping(context.Background(), created)

	cpKey := fmt.Sprintf("%s:%s:%s", created.UUID, "pfold-stream", "shard-0")
	if got := p.kinesisCP[cpKey]; got != "1" {
		t.Fatalf("the fold must clamp the checkpoint at the partial cursor, got %q", got)
	}
	stored, err := esmStore.Get(created.UUID)
	if err != nil {
		t.Fatalf("get mapping: %v", err)
	}
	if !strings.Contains(stored.LastProcessingResult, "record(s) reported in batchItemFailures") {
		t.Fatalf("a partial batch response must be reported, got %q", stored.LastProcessingResult)
	}
}

// countingInvoke records every payload and answers from a script keyed by
// payload content, so bisect tests can fail the whole batch once and then
// answer the halves.
type countingInvoke struct {
	mu       sync.Mutex
	payloads []string
	// answer returns the invocation result for one payload.
	answer func(payload string) (*lambdastore.InvocationResult, error)
}

func (f *countingInvoke) invoke(_ context.Context, _ string, payload []byte) (*lambdastore.InvocationResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.payloads = append(f.payloads, string(payload))
	return f.answer(string(payload))
}

func (f *countingInvoke) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.payloads)
}

func ddbWireItem(seq string) streamBatchItem {
	return streamBatchItem{
		record: map[string]interface{}{
			"dynamodb": map[string]interface{}{"SequenceNumber": seq},
		},
		seq: seq,
	}
}

func partialReport(ids ...string) []byte {
	failures := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		failures = append(failures, map[string]interface{}{"itemIdentifier": id})
	}
	raw, _ := json.Marshal(map[string]interface{}{"batchItemFailures": failures})
	return raw
}

// TestStreamBatchBisectStopsAfterPartialHalf pins the bisect composition:
// when the first half of a bisected batch ends on a partial batch
// response, the second half is not invoked — its records re-read from the
// partial checkpoint.
func TestStreamBatchBisectStopsAfterPartialHalf(t *testing.T) {
	p, _, mapping, _ := newDDBAnchorPoller(t, "TRIM_HORIZON")
	mapping.FunctionResponseTypes = []string{"ReportBatchItemFailures"}
	mapping.BisectBatchOnFunctionError = true
	invoke := &countingInvoke{answer: func(payload string) (*lambdastore.InvocationResult, error) {
		isFullBatch := strings.Contains(payload, `"SequenceNumber":"1"`) && strings.Contains(payload, `"SequenceNumber":"4"`)
		switch {
		case isFullBatch:
			return &lambdastore.InvocationResult{StatusCode: 200, FunctionError: "Handled"}, nil
		case strings.Contains(payload, `"SequenceNumber":"2"`):
			return &lambdastore.InvocationResult{StatusCode: 200, Payload: partialReport("2")}, nil
		default:
			return &lambdastore.InvocationResult{StatusCode: 200, Payload: partialReport()}, nil
		}
	}}
	p.invoke = invoke.invoke

	items := []streamBatchItem{ddbWireItem("1"), ddbWireItem("2"), ddbWireItem("3"), ddbWireItem("4")}
	outcome := p.processStreamBatch(context.Background(), mapping, streamSource{kind: streamSourceDynamoDB}, items)

	if outcome.err != nil {
		t.Fatalf("the batch must not fail: %v", outcome.err)
	}
	if !outcome.short {
		t.Fatalf("a partial first half must terminate the batch's consumption")
	}
	if outcome.lastConsumed != "1" {
		t.Fatalf("the batch cursor must stop at the partial cursor, got %q", outcome.lastConsumed)
	}
	if got := invoke.count(); got != 2 {
		t.Fatalf("the whole batch and the first half must be the only invocations, got %d", got)
	}
	if outcome.reported != 1 {
		t.Fatalf("one reported record must be counted, got %d", outcome.reported)
	}
}

// TestStreamBatchBisectPropagatesSecondHalfPartial pins the other bisect
// composition: a fully consumed first half followed by a second half
// whose response reports its first record keeps the batch cursor at the
// first half's end.
func TestStreamBatchBisectPropagatesSecondHalfPartial(t *testing.T) {
	p, _, mapping, _ := newDDBAnchorPoller(t, "TRIM_HORIZON")
	mapping.FunctionResponseTypes = []string{"ReportBatchItemFailures"}
	mapping.BisectBatchOnFunctionError = true
	invoke := &countingInvoke{answer: func(payload string) (*lambdastore.InvocationResult, error) {
		isFullBatch := strings.Contains(payload, `"SequenceNumber":"1"`) && strings.Contains(payload, `"SequenceNumber":"4"`)
		switch {
		case isFullBatch:
			return &lambdastore.InvocationResult{StatusCode: 200, FunctionError: "Handled"}, nil
		case strings.Contains(payload, `"SequenceNumber":"2"`):
			return &lambdastore.InvocationResult{StatusCode: 200, Payload: partialReport()}, nil
		default:
			return &lambdastore.InvocationResult{StatusCode: 200, Payload: partialReport("3")}, nil
		}
	}}
	p.invoke = invoke.invoke

	items := []streamBatchItem{ddbWireItem("1"), ddbWireItem("2"), ddbWireItem("3"), ddbWireItem("4")}
	outcome := p.processStreamBatch(context.Background(), mapping, streamSource{kind: streamSourceDynamoDB}, items)

	if outcome.err != nil {
		t.Fatalf("the batch must not fail: %v", outcome.err)
	}
	if !outcome.short {
		t.Fatalf("a partial second half must terminate the batch's consumption")
	}
	if outcome.lastConsumed != "2" {
		t.Fatalf("the batch cursor must stay at the first half's end, got %q", outcome.lastConsumed)
	}
	if outcome.reported != 1 {
		t.Fatalf("one reported record must be counted, got %d", outcome.reported)
	}
}
