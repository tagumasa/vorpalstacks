package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// fakeInvoke records invocations and replays scripted results in order.
type fakeInvoke struct {
	mu     sync.Mutex
	calls  int32
	script []*lambdastore.InvocationResult
	errs   []error
}

func (f *fakeInvoke) invoke(ctx context.Context, functionRef string, payload []byte) (*lambdastore.InvocationResult, error) {
	i := int(atomic.AddInt32(&f.calls, 1)) - 1
	f.mu.Lock()
	defer f.mu.Unlock()
	var result *lambdastore.InvocationResult
	var err error
	if i < len(f.script) {
		result = f.script[i]
	}
	if i < len(f.errs) {
		err = f.errs[i]
	}
	if result == nil && err == nil {
		result = &lambdastore.InvocationResult{StatusCode: 200, Payload: []byte(`{"ok":true}`)}
	}
	return result, err
}

func (f *fakeInvoke) callCount() int32 {
	return atomic.LoadInt32(&f.calls)
}

// TestInvokeBatchWithRetry_FunctionErrorRetried pins the contract that a
// response carrying a function error classification fails the batch: the
// poller must retry it like any other failure instead of acknowledging the
// event source records.
func TestInvokeBatchWithRetry_FunctionErrorRetried(t *testing.T) {
	fake := &fakeInvoke{
		script: []*lambdastore.InvocationResult{
			{StatusCode: 200, Payload: []byte(`{"errorMessage":"boom"}`), FunctionError: "Handled"},
			{StatusCode: 200, Payload: []byte(`{}`), FunctionError: "Unhandled"},
			{StatusCode: 200, Payload: []byte(`{"ok":true}`)},
		},
	}
	p := &esmPoller{invoke: fake.invoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: 3,
	}
	if err := p.invokeBatchWithRetry(context.Background(), mapping, []byte(`{}`)); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if got := fake.callCount(); got != 3 {
		t.Fatalf("expected 3 invocations (two function errors then success), got %d", got)
	}
}

// TestInvokeBatchWithRetry_FunctionErrorExhausted verifies the exhausted
// retries surface the function error classification rather than a nil
// error, so the caller can distinguish discard semantics.
func TestInvokeBatchWithRetry_FunctionErrorExhausted(t *testing.T) {
	fake := &fakeInvoke{
		script: []*lambdastore.InvocationResult{
			{StatusCode: 200, FunctionError: "Unhandled"},
			{StatusCode: 200, FunctionError: "Unhandled"},
		},
	}
	p := &esmPoller{invoke: fake.invoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: 1,
	}
	err := p.invokeBatchWithRetry(context.Background(), mapping, []byte(`{}`))
	if err == nil {
		t.Fatal("expected an error after exhausted retries, got nil")
	}
	var fnErr *esmFunctionError
	if !errors.As(err, &fnErr) {
		t.Fatalf("expected *esmFunctionError, got %T: %v", err, err)
	}
	if fnErr.classification != "Unhandled" {
		t.Fatalf("expected Unhandled classification, got %q", fnErr.classification)
	}
	if got := fake.callCount(); got != 2 {
		t.Fatalf("expected 2 invocations (MaximumRetryAttempts=1), got %d", got)
	}
}

// TestInvokeBatchWithRetry_ZeroAttemptsBoundsSingleInvoke pins that
// MaximumRetryAttempts=0 allows exactly one invocation before discard.
func TestInvokeBatchWithRetry_ZeroAttemptsBoundsSingleInvoke(t *testing.T) {
	fake := &fakeInvoke{
		errs: []error{fmt.Errorf("connection refused")},
	}
	p := &esmPoller{invoke: fake.invoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: 0,
	}
	if err := p.invokeBatchWithRetry(context.Background(), mapping, []byte(`{}`)); err == nil {
		t.Fatal("expected the infrastructure error to surface")
	}
	if got := fake.callCount(); got != 1 {
		t.Fatalf("expected a single invocation, got %d", got)
	}
}

// TestPurgeStaleKinesisCheckpoints_CoversDynamoDBKeys pins the cleanup of
// both checkpoint key forms: mappings that are deleted or disabled must
// not leak their in-memory checkpoints, whether keyed for Kinesis or for
// DynamoDB Streams.
func TestPurgeStaleKinesisCheckpoints_CoversDynamoDBKeys(t *testing.T) {
	p := &esmPoller{kinesisCP: map[string]string{
		"11111111-1111-1111-1111-111111111111:stream:s0": "9",
		"ddb:22222222-2222-2222-2222-222222222222":       "15",
		"ddb:33333333-3333-3333-3333-333333333333":       "3",
		"44444444-4444-4444-4444-444444444444:stream:s0": "1",
	}}
	p.purgeStaleKinesisCheckpoints(map[string]struct{}{
		"11111111-1111-1111-1111-111111111111": {},
		"22222222-2222-2222-2222-222222222222": {},
	})
	for _, keep := range []string{
		"11111111-1111-1111-1111-111111111111:stream:s0",
		"ddb:22222222-2222-2222-2222-222222222222",
	} {
		if _, ok := p.kinesisCP[keep]; !ok {
			t.Fatalf("active mapping checkpoint %q must survive the purge", keep)
		}
	}
	for _, drop := range []string{
		"ddb:33333333-3333-3333-3333-333333333333",
		"44444444-4444-4444-4444-444444444444:stream:s0",
	} {
		if _, ok := p.kinesisCP[drop]; ok {
			t.Fatalf("inactive mapping checkpoint %q must be purged", drop)
		}
	}
}

// TestCycleReport_Precedence pins the reporting precedence of one poll
// cycle: a failure outranks a discard, which outranks the deferred success
// report, and an idle cycle never reports success at all.
func TestCycleReport_Precedence(t *testing.T) {
	var r cycleReport
	if r.shouldReportSuccess() {
		t.Fatal("an idle cycle must not report success")
	}
	r.recordProcessed()
	if !r.shouldReportSuccess() {
		t.Fatal("a delivering cycle without failures must report success")
	}
	r.recordDiscard()
	if r.shouldReportSuccess() {
		t.Fatal("a discard in the same cycle must suppress the success report")
	}
	r.recordFailure()
	if r.shouldReportSuccess() {
		t.Fatal("a failure in the same cycle must suppress the success report")
	}

	var lateFailure cycleReport
	lateFailure.recordProcessed()
	lateFailure.recordFailure()
	if lateFailure.shouldReportSuccess() {
		t.Fatal("a failure outranks the success of the same cycle regardless of order")
	}
}

// TestInvokeBatchWithRetry_InfiniteRetriesSingleAttemptPerPoll pins the
// documented -1 default: one attempt per poll cycle; persistence of the
// checkpoint across cycles provides the retry-until-expiry behaviour.
func TestInvokeBatchWithRetry_InfiniteRetriesSingleAttemptPerPoll(t *testing.T) {
	fake := &fakeInvoke{
		errs: []error{fmt.Errorf("throttled")},
	}
	p := &esmPoller{invoke: fake.invoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: -1,
	}
	if err := p.invokeBatchWithRetry(context.Background(), mapping, []byte(`{}`)); err == nil {
		t.Fatal("expected the infrastructure error to surface")
	}
	if got := fake.callCount(); got != 1 {
		t.Fatalf("expected a single in-cycle attempt for -1, got %d", got)
	}
}

// TestInvokeBatchWithRetry_SuccessFirstAttempt verifies a clean response
// with no function error succeeds immediately without retry pauses.
func TestInvokeBatchWithRetry_SuccessFirstAttempt(t *testing.T) {
	fake := &fakeInvoke{}
	p := &esmPoller{invoke: fake.invoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: 5,
	}
	if err := p.invokeBatchWithRetry(context.Background(), mapping, []byte(`{}`)); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if got := fake.callCount(); got != 1 {
		t.Fatalf("expected a single invocation, got %d", got)
	}
}

// kinesisTestItem builds a Kinesis wire-form batch item carrying the given
// sequence number and data.
func kinesisTestItem(seq, data string) streamBatchItem {
	return streamBatchItem{
		record: map[string]interface{}{
			"kinesis": map[string]interface{}{
				"sequenceNumber": seq,
				"data":           data,
			},
		},
		seq: seq,
	}
}

// poisonAwareInvoke fails (Unhandled) for any batch whose payload carries
// the poison marker, succeeding for clean batches.
func poisonAwareInvoke(ctx context.Context, functionRef string, payload []byte) (*lambdastore.InvocationResult, error) {
	if bytes.Contains(payload, []byte("poison")) {
		return &lambdastore.InvocationResult{StatusCode: 200, FunctionError: "Unhandled"}, nil
	}
	return &lambdastore.InvocationResult{StatusCode: 200, Payload: []byte(`{}`)}, nil
}

// TestProcessStreamBatch_BisectIsolatesPoisonRecord pins the documented
// bisection walk: a failing batch splits in two, halves retry recursively,
// and a single failing record is discarded so later records still process.
func TestProcessStreamBatch_BisectIsolatesPoisonRecord(t *testing.T) {
	items := []streamBatchItem{
		kinesisTestItem("1", "ok"), kinesisTestItem("2", "ok"),
		kinesisTestItem("3", "ok"), kinesisTestItem("4", "ok"),
		kinesisTestItem("5", "ok"), kinesisTestItem("6", "poison"),
		kinesisTestItem("7", "ok"), kinesisTestItem("8", "ok"),
	}
	p := &esmPoller{invoke: poisonAwareInvoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:                "arn:aws:lambda:us-east-1:123456789012:function:fn",
		BisectBatchOnFunctionError: true,
		MaximumRetryAttempts:       0,
	}
	outcome := p.processStreamBatch(context.Background(), mapping, testStreamSource("kinesis", "s1"), items)
	if outcome.err != nil {
		t.Fatalf("bisection should consume the batch despite the poison record, got %v", outcome.err)
	}
	if !outcome.discarded {
		t.Fatal("the isolated poison record should be reported as discarded")
	}
	if outcome.lastConsumed != "8" {
		t.Fatalf("checkpoint should advance to the last record, got %q", outcome.lastConsumed)
	}
}

// TestProcessStreamBatch_NoBisectDiscardsWholeBatch pins that without
// bisection a finite retry budget discards the entire failed batch.
func TestProcessStreamBatch_NoBisectDiscardsWholeBatch(t *testing.T) {
	items := []streamBatchItem{
		kinesisTestItem("1", "ok"), kinesisTestItem("2", "poison"),
	}
	p := &esmPoller{invoke: poisonAwareInvoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: 0,
	}
	outcome := p.processStreamBatch(context.Background(), mapping, testStreamSource("kinesis", "s1"), items)
	if outcome.err != nil {
		t.Fatalf("finite budget should discard, got %v", outcome.err)
	}
	if !outcome.discarded {
		t.Fatal("exhausted finite budget should report discarded")
	}
	if outcome.lastConsumed != "2" {
		t.Fatalf("discarded batch still advances the checkpoint, got %q", outcome.lastConsumed)
	}
}

// TestProcessStreamBatch_InfiniteBudgetBlocksAndConsumesNothing pins the
// default -1 budget: the failure surfaces and nothing is consumed, so the
// next poll re-reads the batch.
func TestProcessStreamBatch_InfiniteBudgetBlocksAndConsumesNothing(t *testing.T) {
	items := []streamBatchItem{kinesisTestItem("1", "poison")}
	p := &esmPoller{invoke: poisonAwareInvoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: -1,
	}
	outcome := p.processStreamBatch(context.Background(), mapping, testStreamSource("kinesis", "s1"), items)
	if outcome.err == nil {
		t.Fatal("infinite budget should keep reporting the failure")
	}
	if outcome.discarded {
		t.Fatal("infinite budget must not discard")
	}
	if outcome.lastConsumed != "" {
		t.Fatalf("nothing may be consumed, got checkpoint %q", outcome.lastConsumed)
	}
}

// TestProcessStreamBatch_CleanBatchSucceeds pins the happy path: the whole
// batch succeeds without discard and checkpoints its last record.
func TestProcessStreamBatch_CleanBatchSucceeds(t *testing.T) {
	items := []streamBatchItem{kinesisTestItem("1", "ok"), kinesisTestItem("2", "ok")}
	p := &esmPoller{invoke: poisonAwareInvoke}
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:fn",
		MaximumRetryAttempts: 2,
	}
	outcome := p.processStreamBatch(context.Background(), mapping, testStreamSource("kinesis", "s1"), items)
	if outcome.err != nil {
		t.Fatalf("expected success, got %v", outcome.err)
	}
	if outcome.discarded {
		t.Fatal("clean batch must not report discard")
	}
	if outcome.lastConsumed != "2" {
		t.Fatalf("checkpoint should be the last record, got %q", outcome.lastConsumed)
	}
}

// TestMarshalStreamBatch_RendersRecordsEnvelope verifies the batch payload
// envelope and record order.
func TestMarshalStreamBatch_RendersRecordsEnvelope(t *testing.T) {
	payload := marshalStreamBatch([]streamBatchItem{
		kinesisTestItem("1", "ok"), kinesisTestItem("2", "ok"),
	})
	var decoded struct {
		Records []struct {
			Kinesis struct {
				SequenceNumber string `json:"sequenceNumber"`
			} `json:"kinesis"`
		} `json:"Records"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if len(decoded.Records) != 2 || decoded.Records[0].Kinesis.SequenceNumber != "1" {
		t.Fatalf("unexpected payload: %s", payload)
	}
}

// capturedInvoke records every payload and answers from a script keyed by
// call index; nil script entries answer with a one-item state response.
type capturedInvoke struct {
	mu       sync.Mutex
	payloads [][]byte
	script   []func(call int, payload []byte) (*lambdastore.InvocationResult, error)
}

func (c *capturedInvoke) invoke(ctx context.Context, functionRef string, payload []byte) (*lambdastore.InvocationResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payloads = append(c.payloads, payload)
	i := len(c.payloads) - 1
	if i < len(c.script) && c.script[i] != nil {
		return c.script[i](i, payload)
	}
	return &lambdastore.InvocationResult{StatusCode: 200, Payload: []byte(`{"state":{"n":1}}`)}, nil
}

// decodedWindowEnvelope is the parsed form of one time-window payload.
type decodedWindowEnvelope struct {
	Records                 []json.RawMessage `json:"Records"`
	Window                  streamWindowSpan  `json:"window"`
	State                   json.RawMessage   `json:"state"`
	ShardID                 string            `json:"shardId"`
	EventSourceARN          string            `json:"eventSourceARN"`
	IsFinalInvokeForWindow  bool              `json:"isFinalInvokeForWindow"`
	IsWindowTerminatedEarly bool              `json:"isWindowTerminatedEarly"`
}

func decodeWindowPayload(t *testing.T, payload []byte) decodedWindowEnvelope {
	t.Helper()
	var env decodedWindowEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		t.Fatalf("unmarshal window payload %s: %v", payload, err)
	}
	return env
}

func newWindowTestPoller(inv *capturedInvoke) *esmPoller {
	return &esmPoller{
		invoke:    inv.invoke,
		kinesisCP: make(map[string]string),
		windows:   make(map[string]*shardWindow),
	}
}

// testStreamSource builds the stream identity the batch pipeline tests
// deliver under.
func testStreamSource(kind, shardID string) streamSource {
	if kind == "dynamodb" {
		return streamSource{
			kind:      streamSourceDynamoDB,
			streamArn: "arn:aws:dynamodb:us-east-1:123456789012:table/t/stream/2020-01-01T00:00:00.000",
			shardID:   shardID,
		}
	}
	return streamSource{
		kind:      streamSourceKinesis,
		streamArn: "arn:aws:kinesis:us-east-1:123456789012:stream/s",
		shardID:   shardID,
	}
}

// TestWindowStartOf_AlignsToEpochBoundaries pins the boundary arithmetic:
// windows are epoch-aligned floors of the insertion time.
func TestWindowStartOf_AlignsToEpochBoundaries(t *testing.T) {
	cases := []struct {
		t, windowSeconds, want int64
	}{
		{179, 60, 120},
		{180, 60, 180},
		{0, 60, 0},
		{7, 100, 0},
		{307, 100, 300},
	}
	for _, tc := range cases {
		if got := windowStartOf(tc.t, tc.windowSeconds); got != tc.want {
			t.Fatalf("windowStartOf(%d, %d) = %d, want %d", tc.t, tc.windowSeconds, got, tc.want)
		}
	}
	if got := windowStartOf(500, 0); got != 0 {
		t.Fatalf("windowStartOf with no window must be 0, got %d", got)
	}
}

// TestSplitByWindow_GroupsConsecutiveWindows partitions items by their
// window start while preserving order.
func TestSplitByWindow_GroupsConsecutiveWindows(t *testing.T) {
	items := []windowedStreamItem{
		{item: streamBatchItem{seq: "1"}, windowStart: 100},
		{item: streamBatchItem{seq: "2"}, windowStart: 100},
		{item: streamBatchItem{seq: "3"}, windowStart: 300},
	}
	groups := splitByWindow(items)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if len(groups[0]) != 2 || groups[0][0].item.seq != "1" || groups[0][1].item.seq != "2" {
		t.Fatalf("first group wrong: %+v", groups[0])
	}
	if len(groups[1]) != 1 || groups[1][0].item.seq != "3" {
		t.Fatalf("second group wrong: %+v", groups[1])
	}
}

// TestMarshalWindowedBatch_EnvelopeShape pins the documented time-window
// event shape: window boundaries as ISO-8601 UTC, state passthrough, shard
// and source identification, and the final/early flags.
func TestMarshalWindowedBatch_EnvelopeShape(t *testing.T) {
	payload := marshalWindowedBatch(windowDelivery{
		span:      spanOf(1607497440, 60),
		shardID:   "shardId-000000000001",
		sourceARN: "arn:aws:kinesis:us-east-1:123456789012:stream/test",
		final:     true,
		early:     true,
	}, json.RawMessage(`{"n":5}`), []streamBatchItem{kinesisTestItem("42", "x")})

	env := decodeWindowPayload(t, payload)
	if env.Window.Start != "2020-12-09T07:04:00Z" || env.Window.End != "2020-12-09T07:05:00Z" {
		t.Fatalf("window span wrong: %+v", env.Window)
	}
	if string(env.State) != `{"n":5}` {
		t.Fatalf("state must pass through verbatim, got %s", env.State)
	}
	if env.ShardID != "shardId-000000000001" || env.EventSourceARN != "arn:aws:kinesis:us-east-1:123456789012:stream/test" {
		t.Fatalf("identification wrong: %+v", env)
	}
	if !env.IsFinalInvokeForWindow || !env.IsWindowTerminatedEarly {
		t.Fatalf("final/early flags wrong: %+v", env)
	}
	if len(env.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(env.Records))
	}

	// An absent state renders as an empty object so event.state always
	// exists for the function.
	empty := marshalWindowedBatch(windowDelivery{span: spanOf(0, 60)}, nil, nil)
	env = decodeWindowPayload(t, empty)
	if string(env.State) != "{}" {
		t.Fatalf("initial state must be {}, got %s", env.State)
	}
}

// TestExtractWindowState pins the response contract: the state member is
// required and returned verbatim.
func TestExtractWindowState(t *testing.T) {
	state, err := extractWindowState([]byte(`{"state":{"a":1},"batchItemFailures":[]}`))
	if err != nil {
		t.Fatalf("valid response rejected: %v", err)
	}
	if string(state) != `{"a":1}` {
		t.Fatalf("state not verbatim: %s", state)
	}
	if _, err := extractWindowState([]byte(`{"batchItemFailures":[]}`)); err == nil {
		t.Fatal("response without state must fail")
	}
	if _, err := extractWindowState([]byte(`ok`)); err == nil {
		t.Fatal("non-object response must fail")
	}
	// The runtime appends the return value after any console output, so
	// the state object may sit behind log lines.
	state, err = extractWindowState([]byte("WINDOW_EVENT {\"a\":1}\n{\"state\":{\"n\":1}}"))
	if err != nil {
		t.Fatalf("logged response rejected: %v", err)
	}
	if string(state) != `{"n":1}` {
		t.Fatalf("state from logged response not verbatim: %s", state)
	}
}

// TestProcessStreamWindow_TwoCycleRolloverDeliversFinalInvoke pins the
// core state machine: records are aggregated mid-window across cycles, the
// window close delivers a final invocation carrying the accumulated state
// with no records, the durable checkpoint only advances at that point, and
// the next window starts with a fresh state.
func TestProcessStreamWindow_TwoCycleRolloverDeliversFinalInvoke(t *testing.T) {
	inv := &capturedInvoke{}
	p := newWindowTestPoller(inv)
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:             "arn:aws:lambda:us-east-1:123456789012:function:win",
		TumblingWindowInSeconds: 100,
	}
	const key = "u1:stream:s1"

	// Cycle 1: two records in the 100..200 window.
	items := []windowedStreamItem{
		{item: kinesisTestItem("1", "a"), windowStart: 100},
		{item: kinesisTestItem("2", "a"), windowStart: 100},
	}
	res, err := p.processStreamWindow(context.Background(), mapping, key, testStreamSource("kinesis", "s1"), items, "2")
	if err != nil {
		t.Fatalf("cycle 1: %v", err)
	}
	if res.discarded || !res.processedAny {
		t.Fatalf("cycle 1 result wrong: %+v", res)
	}
	if len(inv.payloads) != 1 {
		t.Fatalf("cycle 1 must invoke once, got %d", len(inv.payloads))
	}
	env := decodeWindowPayload(t, inv.payloads[0])
	if env.IsFinalInvokeForWindow || env.Window.Start != "1970-01-01T00:01:40Z" || len(env.Records) != 2 || string(env.State) != "{}" {
		t.Fatalf("cycle 1 payload wrong: %s", inv.payloads[0])
	}
	p.kinesisCPMu.RLock()
	cp := p.kinesisCP[key]
	p.kinesisCPMu.RUnlock()
	if cp != "" {
		t.Fatalf("checkpoint must not advance mid-window, got %q", cp)
	}

	// Cycle 2: a record from the 300..400 window closes the first window.
	items = []windowedStreamItem{{item: kinesisTestItem("3", "b"), windowStart: 300}}
	if _, err := p.processStreamWindow(context.Background(), mapping, key, testStreamSource("kinesis", "s1"), items, "3"); err != nil {
		t.Fatalf("cycle 2: %v", err)
	}
	if len(inv.payloads) != 3 {
		t.Fatalf("cycle 2 must deliver close + new window, got %d payloads", len(inv.payloads))
	}
	closing := decodeWindowPayload(t, inv.payloads[1])
	if !closing.IsFinalInvokeForWindow || closing.IsWindowTerminatedEarly || len(closing.Records) != 0 {
		t.Fatalf("closing payload wrong: %s", inv.payloads[1])
	}
	if string(closing.State) != `{"n":1}` {
		t.Fatalf("closing must carry the accumulated state, got %s", closing.State)
	}
	if closing.Window.Start != "1970-01-01T00:01:40Z" || closing.ShardID != "s1" {
		t.Fatalf("closing identification wrong: %+v", closing)
	}
	fresh := decodeWindowPayload(t, inv.payloads[2])
	if fresh.IsFinalInvokeForWindow || string(fresh.State) != "{}" || fresh.Window.Start != "1970-01-01T00:05:00Z" {
		t.Fatalf("new window must start fresh: %s", inv.payloads[2])
	}
	p.kinesisCPMu.RLock()
	cp = p.kinesisCP[key]
	p.kinesisCPMu.RUnlock()
	if cp != "2" {
		t.Fatalf("checkpoint must advance to the closed window's last record, got %q", cp)
	}
	p.windowsMu.Lock()
	open := p.windows[key]
	p.windowsMu.Unlock()
	if open == nil || open.windowStart != 300 || open.readSeq != "3" {
		t.Fatalf("new window state wrong: %+v", open)
	}
}

// TestProcessStreamWindow_SingleBatchRolloverMarksFinalChunk pins the
// same-batch case: when a later window's records are already visible, the
// current window's chunk is delivered with the final flag directly.
func TestProcessStreamWindow_SingleBatchRolloverMarksFinalChunk(t *testing.T) {
	inv := &capturedInvoke{}
	p := newWindowTestPoller(inv)
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:             "arn:aws:lambda:us-east-1:123456789012:function:win",
		TumblingWindowInSeconds: 100,
	}
	items := []windowedStreamItem{
		{item: kinesisTestItem("1", "a"), windowStart: 100},
		{item: kinesisTestItem("2", "a"), windowStart: 100},
		{item: kinesisTestItem("3", "b"), windowStart: 300},
	}
	if _, err := p.processStreamWindow(context.Background(), mapping, "k", testStreamSource("kinesis", "s1"), items, "3"); err != nil {
		t.Fatalf("process: %v", err)
	}
	if len(inv.payloads) != 2 {
		t.Fatalf("expected 2 invocations, got %d", len(inv.payloads))
	}
	first := decodeWindowPayload(t, inv.payloads[0])
	if !first.IsFinalInvokeForWindow || len(first.Records) != 2 || string(first.State) != "{}" {
		t.Fatalf("first chunk must be final: %s", inv.payloads[0])
	}
	second := decodeWindowPayload(t, inv.payloads[1])
	if second.IsFinalInvokeForWindow || second.Window.Start != "1970-01-01T00:05:00Z" {
		t.Fatalf("second chunk opens the next window: %s", inv.payloads[1])
	}
	p.kinesisCPMu.RLock()
	cp := p.kinesisCP["k"]
	p.kinesisCPMu.RUnlock()
	if cp != "2" {
		t.Fatalf("checkpoint must advance to the final chunk, got %q", cp)
	}
}

// TestProcessStreamWindow_InactivityCloseAfterGrace pins the closure rule:
// the final invocation fires only once the window's end time has passed
// AND the inactivity grace period elapsed.
func TestProcessStreamWindow_InactivityCloseAfterGrace(t *testing.T) {
	inv := &capturedInvoke{}
	p := newWindowTestPoller(inv)
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:             "arn:aws:lambda:us-east-1:123456789012:function:win",
		TumblingWindowInSeconds: 1,
	}
	items := []windowedStreamItem{{item: kinesisTestItem("1", "a"), windowStart: 100}}
	if _, err := p.processStreamWindow(context.Background(), mapping, "k", testStreamSource("kinesis", "s1"), items, "1"); err != nil {
		t.Fatalf("cycle 1: %v", err)
	}

	p.windowsMu.Lock()
	win := p.windows["k"]
	p.windowsMu.Unlock()
	if win == nil {
		t.Fatal("window must stay open after a mid-window delivery")
	}
	// Before the grace period elapses the window stays open.
	if _, err := p.processStreamWindow(context.Background(), mapping, "k", testStreamSource("kinesis", "s1"), nil, ""); err != nil {
		t.Fatalf("grace cycle: %v", err)
	}
	if len(inv.payloads) != 1 {
		t.Fatalf("no invocation expected during the grace period, got %d", len(inv.payloads))
	}

	p.windowsMu.Lock()
	win.lastActivity = time.Now().Add(-3 * time.Minute)
	p.windowsMu.Unlock()
	if _, err := p.processStreamWindow(context.Background(), mapping, "k", testStreamSource("kinesis", "s1"), nil, ""); err != nil {
		t.Fatalf("close cycle: %v", err)
	}
	if len(inv.payloads) != 2 {
		t.Fatalf("close must invoke once, got %d", len(inv.payloads))
	}
	closing := decodeWindowPayload(t, inv.payloads[1])
	if !closing.IsFinalInvokeForWindow || closing.IsWindowTerminatedEarly || len(closing.Records) != 0 || string(closing.State) != `{"n":1}` {
		t.Fatalf("closing payload wrong: %s", inv.payloads[1])
	}
	p.windowsMu.Lock()
	stillOpen := p.windows["k"]
	p.windowsMu.Unlock()
	if stillOpen != nil {
		t.Fatal("window state must be dropped after the close")
	}
	p.kinesisCPMu.RLock()
	cp := p.kinesisCP["k"]
	p.kinesisCPMu.RUnlock()
	if cp != "1" {
		t.Fatalf("checkpoint must advance at the close, got %q", cp)
	}
}

// TestProcessStreamWindow_OversizeStateTerminatesEarly pins the 1 MB
// ceiling: an oversized aggregated state closes the window with the early
// termination flag.
func TestProcessStreamWindow_OversizeStateTerminatesEarly(t *testing.T) {
	oversize := []byte(`{"state":"` + strings.Repeat("x", maxWindowStateBytes) + `"}`)
	inv := &capturedInvoke{script: []func(call int, payload []byte) (*lambdastore.InvocationResult, error){
		func(call int, payload []byte) (*lambdastore.InvocationResult, error) {
			return &lambdastore.InvocationResult{StatusCode: 200, Payload: oversize}, nil
		},
		nil,
	}}
	p := newWindowTestPoller(inv)
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:             "arn:aws:lambda:us-east-1:123456789012:function:win",
		TumblingWindowInSeconds: 100,
	}
	items := []windowedStreamItem{{item: kinesisTestItem("1", "a"), windowStart: 100}}
	if _, err := p.processStreamWindow(context.Background(), mapping, "k", testStreamSource("kinesis", "s1"), items, "1"); err != nil {
		t.Fatalf("process: %v", err)
	}
	if len(inv.payloads) != 2 {
		t.Fatalf("oversize state must terminate the window early, got %d payloads", len(inv.payloads))
	}
	terminated := decodeWindowPayload(t, inv.payloads[1])
	if !terminated.IsFinalInvokeForWindow || !terminated.IsWindowTerminatedEarly {
		t.Fatalf("termination flags wrong: %s", inv.payloads[1])
	}
	p.windowsMu.Lock()
	stillOpen := p.windows["k"]
	p.windowsMu.Unlock()
	if stillOpen != nil {
		t.Fatal("terminated window must be dropped")
	}
}

// TestProcessStreamWindow_MissingStateDiscardsChunk pins the response
// contract from the function side: a response without the state member is
// a failed invocation; after the retry budget is exhausted the chunk is
// discarded but the window stays open for later chunks.
func TestProcessStreamWindow_MissingStateDiscardsChunk(t *testing.T) {
	inv := &capturedInvoke{script: []func(call int, payload []byte) (*lambdastore.InvocationResult, error){
		func(call int, payload []byte) (*lambdastore.InvocationResult, error) {
			return &lambdastore.InvocationResult{StatusCode: 200, Payload: []byte(`{"ok":true}`)}, nil
		},
	}}
	p := newWindowTestPoller(inv)
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:             "arn:aws:lambda:us-east-1:123456789012:function:win",
		TumblingWindowInSeconds: 100,
		MaximumRetryAttempts:    0,
	}
	items := []windowedStreamItem{{item: kinesisTestItem("1", "a"), windowStart: 100}}
	res, err := p.processStreamWindow(context.Background(), mapping, "k", testStreamSource("kinesis", "s1"), items, "1")
	if err != nil {
		t.Fatalf("exhausted budget must discard, not block: %v", err)
	}
	if !res.discarded {
		t.Fatal("chunk must be reported as discarded")
	}
	if len(inv.payloads) != 1 {
		t.Fatalf("zero budget allows exactly one attempt, got %d", len(inv.payloads))
	}
	p.windowsMu.Lock()
	win := p.windows["k"]
	p.windowsMu.Unlock()
	if win == nil {
		t.Fatal("window must stay open after a discarded chunk")
	}
	if win.readSeq != "1" {
		t.Fatalf("read position must advance past the discarded chunk, got %q", win.readSeq)
	}
}
