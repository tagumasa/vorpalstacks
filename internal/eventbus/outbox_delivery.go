package eventbus

import (
	"context"
	"fmt"
	"sort"
	"time"

	"vorpalstacks/internal/core/resilience"
)

// Start recovers any pending outbox entries and launches async workers.
// Repeated calls are no-ops; a failed recovery may be retried by calling
// Start again.
func (b *EventBus) Start(ctx context.Context) error {
	b.startMu.Lock()
	defer b.startMu.Unlock()
	if b.startDone {
		return nil
	}

	if b.outbox != nil {
		if err := b.recover(ctx); err != nil {
			return fmt.Errorf("eventbus: recovery failed: %w", err)
		}
	}

	for i := 0; i < AsyncWorkerCount; i++ {
		b.wg.Add(1)
		go b.asyncWorker()
	}

	for i := 0; i < AsyncWorkerCount; i++ {
		b.wg.Add(1)
		go b.directWorker()
	}

	if b.outbox != nil {
		b.wg.Add(1)
		go b.cleanupLoop()
		b.wg.Add(1)
		go b.requeuePendingLoop()
	}

	b.startDone = true
	b.started.Store(true)
	return nil
}

// Shutdown stops the bus, waits for all in-flight handlers to complete,
// and closes the outbox store. Even when ctx expires before the workers
// finish, the outbox is still closed once they do.
func (b *EventBus) Shutdown(ctx context.Context) error {
	b.started.Store(false)
	b.stopOnce.Do(func() { close(b.stopCh) })

	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		if b.outbox != nil {
			if err := b.outbox.Close(); err != nil {
				b.logWarn("outbox close error", "error", err)
			}
		}
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *EventBus) recover(ctx context.Context) error {
	// Reset entries stuck in OutboxProcessing (left behind by a crash
	// or panic) so they are included in the pending list below.
	staleCount, err := b.outbox.ResetStaleProcessing(ctx)
	if err != nil {
		return err
	}
	if staleCount > 0 {
		b.logInfo("reset stale processing outbox entries", "count", staleCount)
	}

	// Recovery walks every page of the pending backlog: a single-page scan
	// would permanently miss entries behind the first page, and the async
	// channel is drained during the walk anyway.
	const pageSize = 1000
	after := ""
	for {
		pending, nextCursor, err := b.outbox.ListPendingFrom(ctx, pageSize, after)
		if err != nil {
			return err
		}

		for _, entry := range pending {
			if entry.Depth >= b.maxEventDepth {
				entry.Status = OutboxFailed
				entry.LastError = "max depth exceeded on recovery"
				if err := b.outbox.UpdateEntry(ctx, entry); err != nil {
					b.logWarn("failed to update outbox entry (max depth)", "event_id", entry.EventID, "error", err)
				}
				continue
			}

			if b.registry != nil {
				event, err := b.registry.Deserialize(entry.EventType, entry.SerializedEvent)
				if err != nil {
					b.logWarn("failed to deserialize outbox entry on recovery", "event_id", entry.EventID, "error", err)
					entry.Status = OutboxFailed
					entry.LastError = err.Error()
					if updateErr := b.outbox.UpdateEntry(ctx, entry); updateErr != nil {
						b.logWarn("failed to update outbox entry (deser fail)", "event_id", entry.EventID, "error", updateErr)
					}
					continue
				}
				_ = event
			}

			select {
			case b.asyncCh <- entry:
			default:
				b.logWarn("async channel full during recovery, will retry later", "event_id", entry.EventID)
			}
		}

		if len(pending) < pageSize {
			break
		}
		after = nextCursor
	}

	return nil
}

func (b *EventBus) asyncWorker() {
	defer b.wg.Done()
	defer func() { resilience.RecoverAndRestart("eventbus asyncWorker", &b.wg, b.asyncWorker) }()
	for {
		select {
		case <-b.stopCh:
			return
		case entry := <-b.asyncCh:
			if !b.started.Load() && entry.Status == OutboxPending {
				continue
			}
			b.processOutboxEntry(entry)
		}
	}
}

func (b *EventBus) directWorker() {
	defer b.wg.Done()
	defer func() { resilience.RecoverAndRestart("eventbus directWorker", &b.wg, b.directWorker) }()
	for {
		select {
		case <-b.stopCh:
			return
		case d := <-b.directCh:
			// Semaphore release and panic containment are handled by
			// dispatchWithSemaphores.
			_, _ = b.dispatchWithSemaphores(d.ctx, d.sub, d.event)
		}
	}
}

func (b *EventBus) processOutboxEntry(entry *OutboxEntry) {
	ctx := context.Background()

	// Handler results describe the current attempt only. Subscriber IDs
	// (sub-N) are reassigned on every process start, so keys persisted by
	// a previous attempt or a previous process generation must not leak
	// into this attempt's verdict — a stale "failed" key would fail an
	// entry whose handlers all succeeded, a stale "skipped" key would pin
	// it Pending forever.
	entry.HandlerResults = make(map[string]string)

	updated, err := b.outbox.UpdateStatus(ctx, entry.EventID, OutboxPending, OutboxProcessing)
	if err != nil || !updated {
		return
	}

	snapshot := b.snapshotSubscriptions(entry.EventType)
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].priority > snapshot[j].priority
	})

	event, err := b.deserializeEntry(entry)
	if err != nil {
		b.failEntry(ctx, entry, err.Error())
		return
	}
	// Serialisation strips the in-memory depth counter (json:"-"); restore
	// it from the entry so handlers publishing derived events increment
	// from the true depth and the maxEventDepth cycle prevention applies
	// across outbox hops.
	event.SetEventDepth(entry.Depth)

	allSkipped := true
	for _, sub := range snapshot {
		allSkipped = false

		result, acquired := b.dispatchWithSemaphores(ctx, sub, event)
		if !acquired {
			entry.HandlerResults[sub.id] = "skipped"
			continue
		}

		if result.Error != nil {
			entry.HandlerResults[sub.id] = "failed"
		} else {
			entry.HandlerResults[sub.id] = "delivered"
		}
	}

	if allSkipped {
		entry.Status = OutboxDelivered
		now := time.Now().UTC()
		entry.DeliveredAt = &now
		b.persistEntry(ctx, entry)
		return
	}

	hasFailure := false
	hasSkipped := false
	for _, v := range entry.HandlerResults {
		switch v {
		case "failed":
			hasFailure = true
		case "skipped":
			hasSkipped = true
		}
	}

	switch {
	case hasFailure:
		// A handler ran and returned an error: a genuine delivery failure
		// that consumes the retry budget.
		entry.RetryCount++
		if entry.RetryCount >= entry.MaxRetries {
			entry.Status = OutboxFailed
			b.persistEntry(ctx, entry)
		} else {
			entry.Status = OutboxPending
			b.persistEntry(ctx, entry)
			select {
			case b.asyncCh <- entry:
			default:
			}
		}
	case hasSkipped:
		// "skipped" means the subscriber never ran (semaphore saturated or
		// shutdown in progress): the event is not lost, it simply has not
		// been delivered yet. The entry returns to Pending without
		// consuming the retry budget — otherwise every restart or busy
		// period would burn budget and a handful of deployments could
		// drive a never-attempted entry to OutboxFailed. It is not marked
		// Delivered either, because the skipped subscriber has not
		// received the event; the requeue loop re-enqueues it once
		// capacity returns.
		entry.Status = OutboxPending
		b.persistEntry(ctx, entry)
	default:
		now := time.Now().UTC()
		entry.Status = OutboxDelivered
		entry.DeliveredAt = &now
		b.persistEntry(ctx, entry)
	}
}

func (b *EventBus) failEntry(ctx context.Context, entry *OutboxEntry, reason string) {
	entry.Status = OutboxFailed
	entry.LastError = reason
	if err := b.outbox.UpdateEntry(ctx, entry); err != nil {
		b.logWarn("failed to update outbox entry (fail)", "event_id", entry.EventID, "error", err)
	}
}

func (b *EventBus) persistEntry(ctx context.Context, entry *OutboxEntry) {
	if err := b.outbox.UpdateEntry(ctx, entry); err != nil {
		b.logWarn("failed to persist outbox entry", "event_id", entry.EventID, "status", entry.Status, "error", err)
	}
}

func (b *EventBus) deserializeEntry(entry *OutboxEntry) (Event, error) {
	if b.registry == nil {
		return nil, fmt.Errorf("eventbus: no event registry configured")
	}
	return b.registry.Deserialize(entry.EventType, entry.SerializedEvent)
}

func (b *EventBus) cleanupLoop() {
	defer b.wg.Done()
	defer func() { resilience.RecoverAndRestart("eventbus cleanupLoop", &b.wg, b.cleanupLoop) }()
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-b.stopCh:
			return
		case <-ticker.C:
			now := time.Now().UTC()
			deliveredBefore := now.Add(-DeliveredRetention)
			failedBefore := now.Add(-FailedRetention)
			n, err := b.outbox.Cleanup(context.Background(), deliveredBefore, failedBefore)
			if err != nil {
				b.logWarn("outbox cleanup error", "error", err)
			} else if n > 0 {
				b.logInfo("outbox cleanup completed", "purged", n)
			}
		}
	}
}

// requeuePendingLoop periodically re-enqueues outbox entries that are still
// OutboxPending — entries dropped because the async channel was full, or
// retried while the channel was saturated. Without this loop such entries
// would sit unprocessed until the next restart, violating the at-least-once
// delivery contract.
func (b *EventBus) requeuePendingLoop() {
	defer b.wg.Done()
	defer func() { resilience.RecoverAndRestart("eventbus requeuePendingLoop", &b.wg, b.requeuePendingLoop) }()
	ticker := time.NewTicker(PendingRequeueInterval)
	defer ticker.Stop()

	for {
		select {
		case <-b.stopCh:
			return
		case <-ticker.C:
			b.requeuePending()
		}
	}
}

// requeuePending re-enqueues outbox entries that are still OutboxPending —
// entries dropped because the async channel was full, or retried while the
// channel was saturated. Without this loop such entries would sit
// unprocessed until the next restart, violating the at-least-once delivery
// contract.
func (b *EventBus) requeuePending() {
	if !b.started.Load() {
		return
	}
	// Walk every page of the backlog, resuming from where the previous
	// tick stopped. When the async channel saturates mid-walk the scan
	// stops and the rotation cursor (b.requeueCursor, owned by this loop's
	// single goroutine) remembers the position: the next tick continues
	// behind the entries already queued instead of always re-filling from
	// the head, which would starve the tail of a large backlog. Once a
	// walk completes the cursor wraps back to the head. Re-enqueueing an
	// already-queued entry is safe because processOutboxEntry transitions
	// it via a compare-and-set.
	const pageSize = 1000
	after := b.requeueCursor
	for {
		pending, nextCursor, err := b.outbox.ListPendingFrom(context.Background(), pageSize, after)
		if err != nil {
			b.logWarn("failed to list pending outbox entries for requeue", "error", err)
			return
		}
		for _, entry := range pending {
			select {
			case b.asyncCh <- entry:
			default:
				// Channel full; the next tick resumes from this page so
				// the entries behind it are not starved.
				b.requeueCursor = after
				return
			}
		}
		if len(pending) < pageSize {
			b.requeueCursor = ""
			return
		}
		after = nextCursor
	}
}
