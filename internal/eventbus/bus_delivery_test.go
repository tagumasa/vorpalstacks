package eventbus

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}

func newOutboxBus(t *testing.T) *EventBus {
	t.Helper()
	store := NewPebbleOutboxStore(newTestDB(t))
	registry := NewEventRegistry()
	registry.Register("test:event", func() Event { return &testEvent{} })
	bus := NewEventBus(WithOutbox(store), WithEventRegistry(registry))
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bus.Shutdown(ctx)
	})
	return bus
}

// A panicking handler must neither crash the PublishSync caller nor leak the
// per-subscription semaphore: the next event for the same subscription must
// still be dispatched.
func TestPublishSyncPanicDoesNotLeakSemaphores(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	calls := atomic.Int32{}
	var boomOnce sync.Once
	if _, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		calls.Add(1)
		boomOnce.Do(func() { panic("boom") })
		return HandlerResult{}
	}, WithMaxConcurrency(1)); err != nil {
		t.Fatal(err)
	}

	res, err := bus.PublishSync(context.Background(), &testEvent{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Error == nil {
		t.Fatal("expected the panic to surface as a handler error")
	}

	res, err = bus.PublishSync(context.Background(), &testEvent{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Error != nil {
		t.Fatalf("second dispatch failed, semaphore leaked: %v", res.Error)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected handler called twice, got %d", calls.Load())
	}
}

// A handler panic on the outbox path must release the semaphore (the
// subscription keeps its WithMaxConcurrency budget) and the entry must be
// retried until MaxRetries is exhausted, ending OutboxFailed — never stuck
// in OutboxProcessing with a permanently blocked subscription.
func TestOutboxHandlerPanicRetriesAndReleasesSemaphore(t *testing.T) {
	bus := newOutboxBus(t)
	store := bus.outbox

	attempts := atomic.Int32{}
	if _, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		attempts.Add(1)
		panic("boom")
	}, WithMaxConcurrency(1)); err != nil {
		t.Fatal(err)
	}

	evt := &testEvent{}
	if err := bus.Publish(context.Background(), evt); err != nil {
		t.Fatal(err)
	}

	waitFor(t, 10*time.Second, func() bool {
		entry, err := store.Read(context.Background(), evt.EventID())
		return err == nil && entry != nil && entry.Status == OutboxFailed
	})

	if got := attempts.Load(); got != int32(bus.maxRetries) {
		t.Fatalf("expected %d attempts before failing, got %d", bus.maxRetries, got)
	}
}

// The outbox round-trip strips EventBase.depth (json:"-"); processOutboxEntry
// must restore it from the entry so cycle prevention survives the hop.
func TestOutboxRestoresEventDepth(t *testing.T) {
	bus := newOutboxBus(t)
	store := bus.outbox

	var gotDepth atomic.Int32
	if _, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		gotDepth.Store(int32(event.EventDepth()))
		return HandlerResult{}
	}); err != nil {
		t.Fatal(err)
	}

	evt := &testEvent{}
	if err := bus.Publish(context.Background(), evt); err != nil {
		t.Fatal(err)
	}

	waitFor(t, 10*time.Second, func() bool {
		entry, err := store.Read(context.Background(), evt.EventID())
		return err == nil && entry != nil && entry.Status == OutboxDelivered
	})

	if gotDepth.Load() != 1 {
		t.Fatalf("expected depth 1 after the publish increment, got %d", gotDepth.Load())
	}
}

// An entry whose subscriber could not be dispatched ("skipped") must not be
// marked delivered; it stays retryable so the event is not silently lost.
// acquireSemaphores only reports failure once stopCh is closed (shutdown),
// so the bus here is left unstarted with stopCh closed to reach the skip
// branch deterministically.
func TestOutboxSkippedEntryNotDelivered(t *testing.T) {
	store := NewPebbleOutboxStore(newTestDB(t))
	registry := NewEventRegistry()
	registry.Register("test:event", func() Event { return &testEvent{} })
	bus := NewEventBus(WithOutbox(store), WithEventRegistry(registry))

	called := atomic.Int32{}
	if _, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		called.Add(1)
		return HandlerResult{}
	}, WithEventType("test:event"), WithMaxConcurrency(1)); err != nil {
		t.Fatal(err)
	}

	// Occupy the subscription's single dispatch slot and close stopCh so
	// the semaphore acquisition takes the shutdown escape.
	bus.mu.RLock()
	subs := bus.subscriptions["test:event"]
	if len(subs) != 1 {
		bus.mu.RUnlock()
		t.Fatalf("expected 1 subscription, got %d", len(subs))
	}
	subs[0].sem <- struct{}{}
	bus.mu.RUnlock()
	defer func() { <-subs[0].sem }()
	bus.stopOnce.Do(func() { close(bus.stopCh) })

	serialized, err := SerializeEvent(&testEvent{})
	if err != nil {
		t.Fatal(err)
	}
	entry := &OutboxEntry{
		EventID:         "evt-skipped-1",
		EventType:       "test:event",
		Depth:           0,
		SerializedEvent: serialized,
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}
	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	bus.processOutboxEntry(entry)

	stored, err := store.Read(context.Background(), "evt-skipped-1")
	if err != nil || stored == nil {
		t.Fatal(err)
	}
	if stored.Status == OutboxDelivered {
		t.Fatal("entry marked delivered despite a skipped subscriber")
	}
	if called.Load() != 0 {
		t.Fatal("handler ran although the semaphore was occupied")
	}
}

// A skipped-only outcome must not consume the retry budget: the subscriber
// never ran, so restarts or busy periods cannot drive the entry to
// OutboxFailed. Repeated processing attempts leave it Pending and retryable.
func TestOutboxSkippedEntryDoesNotConsumeRetryBudget(t *testing.T) {
	store := NewPebbleOutboxStore(newTestDB(t))
	registry := NewEventRegistry()
	registry.Register("test:event", func() Event { return &testEvent{} })
	bus := NewEventBus(WithOutbox(store), WithEventRegistry(registry))

	if _, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		return HandlerResult{}
	}, WithEventType("test:event"), WithMaxConcurrency(1)); err != nil {
		t.Fatal(err)
	}

	// Occupy the subscription's single dispatch slot and close stopCh so
	// the semaphore acquisition takes the shutdown escape.
	bus.mu.RLock()
	subs := bus.subscriptions["test:event"]
	if len(subs) != 1 {
		bus.mu.RUnlock()
		t.Fatalf("expected 1 subscription, got %d", len(subs))
	}
	subs[0].sem <- struct{}{}
	bus.mu.RUnlock()
	defer func() { <-subs[0].sem }()
	bus.stopOnce.Do(func() { close(bus.stopCh) })

	serialized, err := SerializeEvent(&testEvent{})
	if err != nil {
		t.Fatal(err)
	}
	entry := &OutboxEntry{
		EventID:         "evt-skipped-budget",
		EventType:       "test:event",
		Depth:           0,
		SerializedEvent: serialized,
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}
	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	// Simulate repeated restart cycles: each attempt skips because the
	// semaphore stays occupied. The budget must remain untouched, so
	// skips alone can never reach OutboxFailed.
	for i := 0; i < 5; i++ {
		bus.processOutboxEntry(entry)
	}

	stored, err := store.Read(context.Background(), "evt-skipped-budget")
	if err != nil || stored == nil {
		t.Fatal(err)
	}
	if stored.Status != OutboxPending {
		t.Fatalf("expected PENDING after repeated skips, got %s", stored.Status)
	}
	if stored.RetryCount != 0 {
		t.Fatalf("skipped attempts must not consume the retry budget, got RetryCount %d", stored.RetryCount)
	}
}

// The requeue scan must paginate: with the async channel (capacity 1024)
// empty and 1500 pending entries, a single-page implementation would
// enqueue at most 1000; a paginated one crosses the page boundary and
// fills the channel to capacity.
func TestRequeuePendingPaginatesBeyondOnePage(t *testing.T) {
	store := NewPebbleOutboxStore(newTestDB(t))
	registry := NewEventRegistry()
	registry.Register("test:event", func() Event { return &testEvent{} })
	bus := NewEventBus(WithOutbox(store), WithEventRegistry(registry))

	serialized, err := SerializeEvent(&testEvent{})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1500; i++ {
		entry := &OutboxEntry{
			EventID:         fmt.Sprintf("evt-%04d", i),
			EventType:       "test:event",
			Depth:           0,
			SerializedEvent: serialized,
			Status:          OutboxPending,
			CreatedAt:       time.Now().UTC(),
			MaxRetries:      3,
			HandlerResults:  map[string]string{},
		}
		if err := store.Write(context.Background(), entry); err != nil {
			t.Fatal(err)
		}
	}

	// Workers are intentionally not started: nothing drains the channel,
	// so requeuePending fills it to capacity and stops.
	bus.started.Store(true)
	bus.requeuePending()

	queued := 0
	seen := make(map[string]bool)
	for {
		select {
		case e := <-bus.asyncCh:
			queued++
			if seen[e.EventID] {
				t.Fatalf("duplicate enqueue: %s", e.EventID)
			}
			seen[e.EventID] = true
		default:
			if queued != cap(bus.asyncCh) {
				t.Fatalf("expected channel filled to capacity %d, got %d enqueued", cap(bus.asyncCh), queued)
			}
			if queued <= 1000 {
				t.Fatal("requeue did not cross the page boundary: entries behind the first page are starved")
			}
			return
		}
	}
}

// Handler results persisted by an earlier attempt or an earlier process
// generation must not poison the current attempt's verdict: a stale
// "failed" key from a subscriber ID that no longer exists must not fail an
// entry whose current handlers all succeeded.
func TestProcessOutboxEntryIgnoresStaleHandlerResults(t *testing.T) {
	store := NewPebbleOutboxStore(newTestDB(t))
	registry := NewEventRegistry()
	registry.Register("test:event", func() Event { return &testEvent{} })
	bus := NewEventBus(WithOutbox(store), WithEventRegistry(registry))
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	if _, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		return HandlerResult{}
	}, WithEventType("test:event")); err != nil {
		t.Fatal(err)
	}

	serialized, err := SerializeEvent(&testEvent{})
	if err != nil {
		t.Fatal(err)
	}
	entry := &OutboxEntry{
		EventID:         "evt-stale-results",
		EventType:       "test:event",
		Depth:           0,
		SerializedEvent: serialized,
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		// Results left behind by a previous process generation whose
		// "sub-0" denoted a different, failing subscriber.
		HandlerResults: map[string]string{"sub-0": "failed", "sub-1": "skipped"},
	}
	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	bus.processOutboxEntry(entry)

	stored, err := store.Read(context.Background(), "evt-stale-results")
	if err != nil || stored == nil {
		t.Fatal(err)
	}
	if stored.Status != OutboxDelivered {
		t.Fatalf("stale handler results poisoned the verdict: status=%s RetryCount=%d", stored.Status, stored.RetryCount)
	}
	if stored.RetryCount != 0 {
		t.Fatalf("a fully successful attempt must not consume the retry budget, got RetryCount %d", stored.RetryCount)
	}
	for id, result := range stored.HandlerResults {
		if result != "delivered" {
			t.Fatalf("stale result for %s survived the attempt: %s", id, result)
		}
	}
}

// When the async channel saturates mid-walk, the next requeue tick must
// resume behind the entries already queued instead of re-filling from the
// head: otherwise the tail of a large backlog starves forever.
func TestRequeuePendingResumesBehindQueuedEntries(t *testing.T) {
	store := NewPebbleOutboxStore(newTestDB(t))
	registry := NewEventRegistry()
	registry.Register("test:event", func() Event { return &testEvent{} })
	bus := NewEventBus(WithOutbox(store), WithEventRegistry(registry))

	serialized, err := SerializeEvent(&testEvent{})
	if err != nil {
		t.Fatal(err)
	}
	const total = 1500
	for i := 0; i < total; i++ {
		entry := &OutboxEntry{
			EventID:         fmt.Sprintf("evt-%04d", i),
			EventType:       "test:event",
			Depth:           0,
			SerializedEvent: serialized,
			Status:          OutboxPending,
			CreatedAt:       time.Now().UTC(),
			MaxRetries:      3,
			HandlerResults:  map[string]string{},
		}
		if err := store.Write(context.Background(), entry); err != nil {
			t.Fatal(err)
		}
	}

	// First walk fills the channel to capacity and stops mid-page.
	bus.started.Store(true)
	bus.requeuePending()
	for {
		select {
		case <-bus.asyncCh:
		default:
			goto drained
		}
	}
drained:

	// The second walk must reach the tail instead of re-filling the head.
	bus.requeuePending()
	tailSeen := false
	seen := make(map[string]bool)
	for {
		select {
		case e := <-bus.asyncCh:
			if seen[e.EventID] {
				t.Fatalf("duplicate enqueue: %s", e.EventID)
			}
			seen[e.EventID] = true
			if e.EventID == fmt.Sprintf("evt-%04d", total-1) {
				tailSeen = true
			}
		default:
			if !tailSeen {
				t.Fatalf("second requeue walk starved the tail: %d entries, tail missing", len(seen))
			}
			if len(seen) > cap(bus.asyncCh) {
				t.Fatalf("unexpected enqueue count %d", len(seen))
			}
			return
		}
	}
}

// Repeated Start calls must not launch a second set of workers.
func TestStartIsIdempotent(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	// wg counts: 2*AsyncWorkerCount workers. A second Start would add 16
	// more; Shutdown would then block on them.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := bus.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown after double start failed or hung: %v", err)
	}
}

// Event IDs must be unique even when generated concurrently.
func TestGenerateEventIDConcurrentUnique(t *testing.T) {
	const goroutines = 16
	const perGoroutine = 200

	ids := make(chan string, goroutines*perGoroutine)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				ids <- generateEventID("test:event")
			}
		}()
	}
	wg.Wait()
	close(ids)

	seen := make(map[string]bool, goroutines*perGoroutine)
	for id := range ids {
		if seen[id] {
			t.Fatalf("duplicate event ID generated: %s", id)
		}
		seen[id] = true
	}
}
