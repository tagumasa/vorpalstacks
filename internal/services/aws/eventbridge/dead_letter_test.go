package eventbridge

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// The documented DLQ contract at PutTargets: SQS queues only, standard
// queues only, same region as the rule.
func TestValidateDeadLetterQueueARN(t *testing.T) {
	sqsARN := "arn:aws:sqs:us-east-1:000000000000:events-dlq"
	cases := []struct {
		name       string
		arn        string
		ruleRegion string
		ok         bool
	}{
		{"sqs queue", sqsARN, "us-east-1", true},
		{"no dlq configured", "", "us-east-1", true},
		{"region unknown to the caller", sqsARN, "", true},
		{"lambda arn rejected", "arn:aws:lambda:us-east-1:000000000000:function:fn", "us-east-1", false},
		{"sns arn rejected", "arn:aws:sns:us-east-1:000000000000:topic/t", "us-east-1", false},
		{"fifo queue rejected", "arn:aws:sqs:us-east-1:000000000000:events-dlq.fifo", "us-east-1", false},
		{"cross-region queue rejected", "arn:aws:sqs:eu-west-1:000000000000:events-dlq", "us-east-1", false},
	}
	for _, tc := range cases {
		err := validateDeadLetterQueueARN(tc.arn, tc.ruleRegion)
		if tc.ok && err != nil {
			t.Errorf("%s: expected acceptance, got %v", tc.name, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("%s: expected rejection", tc.name)
		}
	}
}

// recordingSQSInvoker captures every SendMessage's queue, body and options.
type recordingSQSInvoker struct {
	mu        sync.Mutex
	failFirst int
	sends     []recordedSQSSend
}

type recordedSQSSend struct {
	region   string
	queueURL string
	body     string
	opts     invokers.SQSSendOptions
}

func (r *recordingSQSInvoker) GetQueueByName(_ context.Context, region, queueName string) (string, error) {
	return "http://vorpalstacks:50080/" + region + "/" + queueName, nil
}

func (r *recordingSQSInvoker) GetQueueARN(_ context.Context, region, queueURL string) (string, error) {
	parts := strings.SplitN(queueURL, "/", 5)
	return "arn:aws:sqs:" + region + ":000000000000:" + parts[len(parts)-1], nil
}

func (r *recordingSQSInvoker) SendMessage(_ context.Context, region, queueURL, body string, opts invokers.SQSSendOptions) (string, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.failFirst > 0 {
		r.failFirst--
		return "", "", errors.New("sqs unavailable")
	}
	r.sends = append(r.sends, recordedSQSSend{region: region, queueURL: queueURL, body: body, opts: opts})
	return "msg-id", "md5", nil
}

func (r *recordingSQSInvoker) ReceiveMessage(_ context.Context, _ string, _ string, _ int32, _ *int32, _ int32) ([]invokers.ReceivedSQSMessage, error) {
	return nil, nil
}

func (r *recordingSQSInvoker) DeleteMessage(_ context.Context, _ string, _ string, _ string) error {
	return nil
}

func (r *recordingSQSInvoker) recorded() []recordedSQSSend {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recordedSQSSend(nil), r.sends...)
}

var _ invokers.SQSInvoker = (*recordingSQSInvoker)(nil)

func newDLQTestService(t *testing.T, sqs *recordingSQSInvoker, kinesisFail int) *EventsService {
	t.Helper()
	t.Setenv("TEST_MODE", "true")
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	t.Cleanup(svc.Close)
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(sqs)
	bus.SetKinesisInvoker(&recordingKinesisInvoker{failNext: kinesisFail})
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}
	return svc
}

// The dead-letter write carries the documented attribute envelope: RULE_ARN
// and TARGET_ARN identify the delivery, ERROR_CODE/ERROR_MESSAGE carry the
// terminal failure, EXHAUSTED_RETRY_CONDITION names the retry condition
// that ended the delivery, RETRY_ATTEMPTS counts the retries, and the trace
// header rides as the AWSTraceHeader message attribute. The message body
// is the transformed target payload.
func TestHandleBusDeliveryDeadLetterEnvelope(t *testing.T) {
	sqs := &recordingSQSInvoker{}
	svc := newDLQTestService(t, sqs, 1000)

	ruleARN := "arn:aws:events:us-east-1:000000000000:rule/orders/push"
	evt := &eventbus.EventBridgeDeliveryEvent{
		RuleARN:              ruleARN,
		TargetARN:            "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		Input:                []byte(`{"id":"evt-dlq"}`),
		EventBridgeEventID:   "evt-dlq",
		TraceHeader:          "Root=1-5e272a90-8a3d1a17d2d0e8a1c1f0a1c2;Parent=94ae789b969f1cc5;Sampled=1",
		RetryPolicySet:       true,
		MaximumRetryAttempts: 0,
		DeadLetterConfigArn:  "arn:aws:sqs:us-east-1:000000000000:events-dlq",
	}
	evt.Region = "us-east-1"

	if res := svc.handleBusDelivery(context.Background(), evt); res.Error != nil {
		t.Fatalf("a successful dead-letter write completes the delivery: %v", res.Error)
	}
	sends := sqs.recorded()
	if len(sends) != 1 {
		t.Fatalf("exactly one dead-letter write must happen, got %d", len(sends))
	}
	send := sends[0]
	if send.body != `{"id":"evt-dlq"}` {
		t.Fatalf("the DLQ body must be the target payload, got %q", send.body)
	}
	attr := func(name string) (string, bool) {
		a, ok := send.opts.TypedMessageAttributes[name]
		return a.StringValue, ok
	}
	for name, want := range map[string]string{
		"RULE_ARN":                  ruleARN,
		"TARGET_ARN":                "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		"ERROR_CODE":                "ERROR_FROM_TARGET",
		"EXHAUSTED_RETRY_CONDITION": "MaximumRetryAttempts",
		"RETRY_ATTEMPTS":            "0",
		"AWSTraceHeader":            evt.TraceHeader,
	} {
		got, ok := attr(name)
		if !ok {
			t.Errorf("attribute %s missing from the DLQ envelope", name)
			continue
		}
		if got != want {
			t.Errorf("attribute %s = %q, want %q", name, got, want)
		}
	}
	if msg, ok := attr("ERROR_MESSAGE"); !ok || msg == "" {
		t.Error("ERROR_MESSAGE must carry the terminal failure text")
	}
}

// A caller context cancelled mid-retry still dead-letters: the synchronous
// dispatch loop gives the event the same terminal handling the engine path
// provides, and the write rides a detached bounded context so the
// cancellation that ended the retry cannot also kill it.
func TestSyncDispatchDeadLettersOnCancellation(t *testing.T) {
	sqs := &recordingSQSInvoker{}
	svc := newDLQTestService(t, sqs, 1000)

	event := &eventsstore.Event{ID: "evt-cancel"}
	target := eventsstore.Target{
		ARN:              "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		DeadLetterConfig: &eventsstore.DeadLetterConfig{Arn: "arn:aws:sqs:us-east-1:000000000000:events-dlq"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := svc.dispatchToTarget(ctx, "us-east-1", "arn:aws:events:us-east-1:000000000000:rule/orders/push", event, target, []byte(`{"id":"evt-cancel"}`))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("the cancelled dispatch must return the context error, got %v", err)
	}
	sends := sqs.recorded()
	if len(sends) != 1 {
		t.Fatalf("the cancelled delivery must dead-letter exactly once, got %d sends", len(sends))
	}
	if sends[0].body != `{"id":"evt-cancel"}` {
		t.Fatalf("the DLQ body must be the target payload, got %q", sends[0].body)
	}
}

// A permanent failure dead-letters with NO_RESOURCE and without a retry
// condition: no retry budget was consumed.
func TestHandleBusDeliveryDeadLetterPermanentNoRetryCondition(t *testing.T) {
	sqs := &recordingSQSInvoker{}
	svc := newDLQTestService(t, sqs, 0) // kinesis invoker never used

	evt := &eventbus.EventBridgeDeliveryEvent{
		TargetARN:           "arn:aws:firehose:us-east-1:000000000000:deliverystream/none",
		Input:               []byte(`{}`),
		EventBridgeEventID:  "evt-dlq-perm",
		DeadLetterConfigArn: "arn:aws:sqs:us-east-1:000000000000:events-dlq",
	}
	evt.Region = "us-east-1"

	if res := svc.handleBusDelivery(context.Background(), evt); res.Error != nil {
		t.Fatalf("a successful dead-letter write completes the delivery: %v", res.Error)
	}
	sends := sqs.recorded()
	if len(sends) != 1 {
		t.Fatalf("exactly one dead-letter write must happen, got %d", len(sends))
	}
	attrs := sends[0].opts.TypedMessageAttributes
	if got := attrs["ERROR_CODE"].StringValue; got != "NO_RESOURCE" {
		t.Fatalf("a permanent failure must code NO_RESOURCE, got %q", got)
	}
	if _, ok := attrs["EXHAUSTED_RETRY_CONDITION"]; ok {
		t.Fatal("a permanent failure exhausted no retry condition; the attribute must be absent")
	}
	if got := attrs["RETRY_ATTEMPTS"].StringValue; got != "0" {
		t.Fatalf("a permanent failure performed no retries, got RETRY_ATTEMPTS=%s", got)
	}
}

// The dead-letter write is bounded-retried: a transient queue fault during
// the routing does not lose the last durable copy of the event.
func TestRouteToDeadLetterRetriesTransientWriteFailure(t *testing.T) {
	sqs := &recordingSQSInvoker{failFirst: 2}
	svc := newDLQTestService(t, sqs, 0)

	job := svc.newDeliveryJob("us-east-1", &eventsstore.Event{ID: "evt-dlq-retry"},
		eventsstore.Target{
			ARN:              "arn:aws:firehose:us-east-1:000000000000:deliverystream/none",
			DeadLetterConfig: &eventsstore.DeadLetterConfig{Arn: "arn:aws:sqs:us-east-1:000000000000:events-dlq"},
		}, []byte(`{}`))

	start := time.Now()
	if err := svc.routeToDeadLetter(context.Background(), job, errors.New("primary failed")); err != nil {
		t.Fatalf("the retried write must eventually succeed: %v", err)
	}
	if sends := sqs.recorded(); len(sends) != 1 {
		t.Fatalf("the queue must receive exactly one copy after the retries, got %d", len(sends))
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("the bounded retry must stay short, took %s", elapsed)
	}
}
