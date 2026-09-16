package eventbridge

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// countedAttempt is a test attempt callback that fails the first failN
// calls (or every call when permanent), counting invocations.
type countedAttempt struct {
	calls     atomic.Int32
	failN     int32
	permanent bool
}

func (c *countedAttempt) attempt(_ context.Context, _ *deliveryJob) error {
	n := c.calls.Add(1)
	if n <= c.failN {
		if c.permanent {
			return permanentDeliveryError("counted permanent failure")
		}
		return errors.New("counted transient failure")
	}
	return nil
}

func quickJob() *deliveryJob {
	return &deliveryJob{
		region:     "us-east-1",
		event:      &eventsstore.Event{ID: "evt-engine"},
		target:     eventsstore.Target{ARN: "arn:aws:kinesis:us-east-1:000000000000:stream/demo"},
		payload:    []byte(`{}`),
		maxRetries: 3,
		deadline:   time.Now().Add(5 * time.Second),
		backoff:    time.Millisecond,
		maxBackoff: 2 * time.Millisecond,
	}
}

// A transient first failure is retried by the engine until it succeeds;
// the worker never sleeps, the schedule drives the retry.
func TestDeliveryEngineRetriesUntilSuccess(t *testing.T) {
	counter := &countedAttempt{failN: 2}
	var terminalCalls atomic.Int32
	engine := newDeliveryEngine(
		func(ctx context.Context, job *deliveryJob) error { return counter.attempt(ctx, job) },
		func(ctx context.Context, job *deliveryJob, err error) error {
			terminalCalls.Add(1)
			return nil
		},
	)
	go engine.runJob(quickJob())
	if !waitFor(time.Second, func() bool { return counter.calls.Load() == 3 }) {
		t.Fatalf("expected 3 attempts (2 failures + 1 success), got %d", counter.calls.Load())
	}
	if terminalCalls.Load() != 0 {
		t.Fatalf("a succeeded delivery must not reach terminal handling, got %d calls", terminalCalls.Load())
	}
	engine.Close()
}

// A permanent failure terminates on the first attempt: no retry budget is
// consumed on a target type that has no delivery path.
func TestDeliveryEnginePermanentFailureSkipsRetries(t *testing.T) {
	counter := &countedAttempt{failN: 1000, permanent: true}
	var terminalErrs []error
	var mu sync.Mutex
	engine := newDeliveryEngine(
		func(ctx context.Context, job *deliveryJob) error { return counter.attempt(ctx, job) },
		func(ctx context.Context, job *deliveryJob, err error) error {
			mu.Lock()
			defer mu.Unlock()
			terminalErrs = append(terminalErrs, err)
			return nil
		},
	)
	job := quickJob()
	engine.runJob(job)
	if got := counter.calls.Load(); got != 1 {
		t.Fatalf("a permanent failure must terminate after one attempt, got %d", got)
	}
	if len(terminalErrs) != 1 || !errors.Is(terminalErrs[0], errPermanentDelivery) {
		t.Fatalf("the terminal handling must see the permanent error once, got %v", terminalErrs)
	}
	if engine.pending() != 0 {
		t.Fatalf("no retry may be scheduled for a permanent failure, %d pending", engine.pending())
	}
	engine.Close()
}

// Retry exhaustion terminates the delivery exactly once after the full
// attempt budget.
func TestDeliveryEngineExhaustionTerminates(t *testing.T) {
	counter := &countedAttempt{failN: 1000}
	var terminalCalls atomic.Int32
	engine := newDeliveryEngine(
		func(ctx context.Context, job *deliveryJob) error { return counter.attempt(ctx, job) },
		func(ctx context.Context, job *deliveryJob, err error) error {
			terminalCalls.Add(1)
			return nil
		},
	)
	go engine.runJob(quickJob()) // maxRetries 3 → 4 total attempts
	if !waitFor(time.Second, func() bool { return counter.calls.Load() == 4 && terminalCalls.Load() == 1 }) {
		t.Fatalf("expected the full budget of 4 attempts and one terminal call, got %d attempts / %d terminal calls",
			counter.calls.Load(), terminalCalls.Load())
	}
	engine.Close()
}

// waitFor polls cond until it holds or the timeout expires.
func waitFor(timeout time.Duration, cond func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return cond()
}

// Close cancels pending retries: a job scheduled for a distant retry dies
// with the engine instead of holding shutdown hostage.
func TestDeliveryEngineCloseCancelsPendingRetries(t *testing.T) {
	engine := newDeliveryEngine(
		func(context.Context, *deliveryJob) error { return nil },
		func(context.Context, *deliveryJob, error) error { return nil },
	)
	job := quickJob()
	job.nextDue = time.Now().Add(time.Minute)
	if err := engine.schedule(job); err != nil {
		t.Fatalf("schedule: %v", err)
	}
	closed := make(chan struct{})
	go func() {
		engine.Close()
		close(closed)
	}()
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Fatal("Close must cancel the pending retry instead of waiting it out")
	}
}

// A shutdown race strands the retry on a dead queue: schedule's post-push
// cancellation check returns the error (which runJob reports), and the
// scheduler's exit drain reports jobs queued before it left. The stranded
// job itself is unreachable by design — it dies with the process — so the
// pin is the error return that keeps the loss on a reporting path.
func TestDeliveryEngineScheduleAfterCloseReports(t *testing.T) {
	engine := newDeliveryEngine(
		func(context.Context, *deliveryJob) error { return nil },
		func(context.Context, *deliveryJob, error) error { return nil },
	)
	engine.Close()
	if err := engine.schedule(quickJob()); err == nil {
		t.Fatal("scheduling on a shut-down engine must return the shutdown error for the caller to report")
	}
	// runJob on a shut-down engine terminates without hanging and without
	// queueing anything the dead scheduler would have to drain.
	engine.runJob(&deliveryJob{
		region:     "us-east-1",
		event:      &eventsstore.Event{ID: "evt-lost"},
		target:     eventsstore.Target{ARN: "arn:aws:kinesis:us-east-1:000000000000:stream/demo"},
		payload:    []byte(`{}`),
		maxRetries: 3,
		deadline:   time.Now().Add(time.Second),
		backoff:    time.Millisecond,
		maxBackoff: 2 * time.Millisecond,
	})
	if n := engine.pending(); n != 0 {
		t.Fatalf("a retry scheduled after shutdown must not sit silently on the dead queue, pending = %d", n)
	}
}

// A permanently-unavailable target type (firehose before the service
// exists) must fail fast: the synchronous driver returns the terminal
// error without sleeping out the production backoff.
func TestSyncDispatchPermanentTargetFailsFast(t *testing.T) {
	// No TEST_MODE: the production defaults (185 retries / 24 h) apply, so
	// a non-fail-fast implementation would block for the test timeout.
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	start := time.Now()
	err = svc.dispatchToTarget(ctx, "us-east-1", "",
		&eventsstore.Event{ID: "evt-failfast"},
		eventsstore.Target{ARN: "arn:aws:firehose:us-east-1:000000000000:deliverystream/none"},
		[]byte(`{}`))
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("a firehose target must fail in this deployment")
	}
	if !errors.Is(err, errPermanentDelivery) {
		t.Fatalf("the firehose stub must be classified permanent, got: %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("a permanent failure must fail fast, took %s", elapsed)
	}
}

// The bus-less sync arm surfaces its delivery outcome: a rule whose target
// fails permanently (firehose before the service exists) makes deliverEvent
// return the terminal error, so the PutEvents entry and the scheduled fire
// report the loss instead of recording a silent success.
func TestSyncArmSurfacesTargetFailure(t *testing.T) {
	// No TEST_MODE: the production defaults apply; the permanent failure
	// still terminates immediately.
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	store, err := svc.GetStoreForRegion("us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	rule := &eventsstore.Rule{
		Name: "sync-fail",
		// A pattern rule: pattern-less rules are scheduled rules and never
		// match bus events, so the surface-failure pin needs a pattern
		// matching the delivered event's source.
		EventPattern: `{"source":["s"]}`, EventBusName: "default", State: eventsstore.RuleStateEnabled,
	}
	if err := store.CreateRule(ctx, rule); err != nil {
		t.Fatalf("create rule: %v", err)
	}
	if err := store.PutTarget(ctx, &eventsstore.Target{
		ID:           "t1",
		RuleName:     "sync-fail",
		EventBusName: "default",
		ARN:          "arn:aws:firehose:us-east-1:000000000000:deliverystream/none",
	}); err != nil {
		t.Fatalf("put target: %v", err)
	}

	err = svc.deliverEvent(ctx, store, &eventsstore.Event{
		ID: "evt-syncsurf", DetailType: "t", Source: "s", Region: "us-east-1",
		EventBusName: "default", Detail: map[string]interface{}{},
	}, "default", "us-east-1", 0)
	if err == nil {
		t.Fatal("the synchronous arm must surface the target's terminal error, got nil")
	}
	if !errors.Is(err, errPermanentDelivery) {
		t.Fatalf("the surfaced error is the delivery failure, got: %v", err)
	}
}

// An explicit MaximumRetryAttempts of 0 crosses the bus as a present
// policy (RetryPolicySet): exactly one attempt runs, and the failure
// terminates inline instead of falling back to the 185-retry defaults.
func TestHandleBusDeliveryRetryZeroIsSingleAttempt(t *testing.T) {
	t.Setenv("TEST_MODE", "true")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()
	invoker := &recordingKinesisInvoker{failNext: 1000}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}

	evt := &eventbus.EventBridgeDeliveryEvent{
		TargetARN:                "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		Input:                    []byte(`{}`),
		EventBridgeEventID:       "evt-retry-zero",
		RetryPolicySet:           true,
		MaximumRetryAttempts:     0,
		MaximumEventAgeInSeconds: 60,
	}
	evt.Region = "us-east-1"

	res := svc.handleBusDelivery(context.Background(), evt)
	if res.Error != nil {
		t.Fatalf("a spent zero-retry budget without DLQ is the documented drop, not a handler failure: %v", res.Error)
	}
	if got := len(invoker.keys()); got != 1 {
		t.Fatalf("retry-0 must run exactly one attempt, got %d", got)
	}
	if svc.delivery.pending() != 0 {
		t.Fatalf("retry-0 must not schedule engine retries, %d pending", svc.delivery.pending())
	}
}

// A transient first failure hands the retry to the engine; the handler
// answers success (the AWS retry lifecycle now owns the delivery) and the
// engine completes it.
func TestHandleBusDeliveryTransientFailureDefersToEngine(t *testing.T) {
	t.Setenv("TEST_MODE", "true")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()
	invoker := &recordingKinesisInvoker{failNext: 1}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}

	evt := &eventbus.EventBridgeDeliveryEvent{
		TargetARN:          "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		Input:              []byte(`{}`),
		EventBridgeEventID: "evt-transient",
	}
	evt.Region = "us-east-1"

	res := svc.handleBusDelivery(context.Background(), evt)
	if res.Error != nil {
		t.Fatalf("a deferred retry is not a handler failure: %v", res.Error)
	}
	deadline := time.Now().Add(5 * time.Second)
	for len(invoker.keys()) < 2 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := len(invoker.keys()); got != 2 {
		t.Fatalf("the engine must complete the retry, got %d attempts", got)
	}
}

// When the engine is shut down, a transient failure cannot be deferred:
// the handler reports the rejection so the transport-level outbox retry
// re-drives the delivery instead of recording it as completed.
func TestHandleBusDeliveryReportsEngineRejection(t *testing.T) {
	t.Setenv("TEST_MODE", "true")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	invoker := &recordingKinesisInvoker{failNext: 1000}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}
	svc.Close() // engine gone; the inline attempt still runs

	evt := &eventbus.EventBridgeDeliveryEvent{
		TargetARN:          "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		Input:              []byte(`{}`),
		EventBridgeEventID: "evt-rejected",
	}
	evt.Region = "us-east-1"

	res := svc.handleBusDelivery(context.Background(), evt)
	if res.Error == nil {
		t.Fatal("an undetectable retry handoff must be reported as a handler failure")
	}
	if got := len(invoker.keys()); got != 1 {
		t.Fatalf("only the inline attempt may run, got %d", got)
	}
}

// A spent terminal delivery whose dead-letter write also fails is a loss:
// the handler must report it so the outbox retry re-drives the delivery.
func TestHandleBusDeliveryReportsDeadLetterFailure(t *testing.T) {
	t.Setenv("TEST_MODE", "true")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()
	invoker := &recordingKinesisInvoker{failNext: 1000}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}

	evt := &eventbus.EventBridgeDeliveryEvent{
		TargetARN:            "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		Input:                []byte(`{}`),
		EventBridgeEventID:   "evt-dlq-fail",
		RetryPolicySet:       true,
		MaximumRetryAttempts: 0,
		DeadLetterConfigArn:  "arn:aws:sqs:us-east-1:000000000000:no-such-queue",
	}
	evt.Region = "us-east-1"

	res := svc.handleBusDelivery(context.Background(), evt)
	if res.Error == nil {
		t.Fatal("a failed dead-letter write after a spent budget must be reported as a handler failure")
	}
}

// The cross-bus hop counter rides the delivery event: at or beyond the
// bound the handler refuses the delivery instead of re-publishing the
// cycle, and one hop below the bound still delivers.
func TestHandleBusDeliveryHopDepthGuard(t *testing.T) {
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()

	busARN := "arn:aws:events:us-east-1:000000000000:event-bus/default"
	atLimit := &eventbus.EventBridgeDeliveryEvent{TargetARN: busARN, Input: []byte(`{}`), HopDepth: maxCrossBusDepth}
	atLimit.Region = "us-east-1"
	if res := svc.handleBusDelivery(context.Background(), atLimit); res.Error == nil {
		t.Fatal("a delivery at the cross-bus depth bound must be refused")
	}

	belowLimit := &eventbus.EventBridgeDeliveryEvent{TargetARN: busARN, Input: []byte(`{}`), HopDepth: maxCrossBusDepth - 1}
	belowLimit.Region = "us-east-1"
	if res := svc.handleBusDelivery(context.Background(), belowLimit); res.Error != nil {
		t.Fatalf("a delivery below the bound must proceed: %v", res.Error)
	}
}
