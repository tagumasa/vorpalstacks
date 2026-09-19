package eventbus

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// panicOnceOutbox delegates to the real store but blows up on the first
// Pending-to-Processing transition, simulating a store-layer defect
// striking inside an async worker's frame.
type panicOnceOutbox struct {
	inner OutboxStore
	fired atomic.Bool
}

func (p *panicOnceOutbox) Write(ctx context.Context, entry *OutboxEntry) error {
	return p.inner.Write(ctx, entry)
}

func (p *panicOnceOutbox) Read(ctx context.Context, eventID string) (*OutboxEntry, error) {
	return p.inner.Read(ctx, eventID)
}

func (p *panicOnceOutbox) UpdateStatus(ctx context.Context, eventID string, from, to OutboxStatus) (bool, error) {
	if from == OutboxPending && to == OutboxProcessing && p.fired.CompareAndSwap(false, true) {
		panic("pin: outbox store blew up mid-transition")
	}
	return p.inner.UpdateStatus(ctx, eventID, from, to)
}

func (p *panicOnceOutbox) UpdateEntry(ctx context.Context, entry *OutboxEntry) error {
	return p.inner.UpdateEntry(ctx, entry)
}

func (p *panicOnceOutbox) ListPendingFrom(ctx context.Context, limit int, afterCursor string) ([]*OutboxEntry, string, error) {
	return p.inner.ListPendingFrom(ctx, limit, afterCursor)
}

func (p *panicOnceOutbox) ResetStaleProcessing(ctx context.Context) (int, error) {
	return p.inner.ResetStaleProcessing(ctx)
}

func (p *panicOnceOutbox) Delete(ctx context.Context, eventID string) error {
	return p.inner.Delete(ctx, eventID)
}

func (p *panicOnceOutbox) Cleanup(ctx context.Context, deliveredBefore time.Time, failedBefore time.Time) (int, error) {
	return p.inner.Cleanup(ctx, deliveredBefore, failedBefore)
}

func (p *panicOnceOutbox) Close() error {
	return p.inner.Close()
}

// A panic inside an async worker's frame (here: the outbox store blowing
// up mid-transition) must not kill the worker goroutine — the containment
// defer restarts it, and the at-least-once machinery recovers the entry:
// it stays Pending and the requeue loop's pending scan re-enqueues it on
// the next periodic pass. The scan interval is shortened for the test
// budget, and the fired flag proves the panic actually struck.
func TestAsyncWorkerPanicsAndStillDelivers(t *testing.T) {
	origInterval := PendingRequeueInterval
	PendingRequeueInterval = 100 * time.Millisecond
	defer func() {
		PendingRequeueInterval = origInterval
	}()

	wrapped := &panicOnceOutbox{inner: NewPebbleOutboxStore(newTestDB(t))}
	registry := NewEventRegistry()
	registry.Register("test:event", func() Event { return &testEvent{} })
	bus := NewEventBus(WithOutbox(wrapped), WithEventRegistry(registry))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bus.Shutdown(ctx)
	}()

	var calls atomic.Int32
	if _, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		calls.Add(1)
		return HandlerResult{}
	}, WithEventType("test:event")); err != nil {
		t.Fatal(err)
	}
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	victim := &testEvent{}
	if err := bus.Publish(context.Background(), victim); err != nil {
		t.Fatal(err)
	}

	waitFor(t, 10*time.Second, func() bool { return calls.Load() >= 1 })
	if !wrapped.fired.Load() {
		t.Fatal("the store panic never fired; the pin ran vacuously")
	}
	entry, err := wrapped.Read(context.Background(), victim.EventID())
	if err != nil || entry == nil || entry.Status != OutboxDelivered {
		t.Fatalf("victim not delivered after the worker restart: entry=%v err=%v", entry, err)
	}
}
