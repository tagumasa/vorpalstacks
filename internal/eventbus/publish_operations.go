package eventbus

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// preparePublish enforces the preconditions shared by Publish and
// PublishSync — the bus must be started, the event non-nil with a
// non-empty type — and increments the event's dispatch depth, rejecting
// events at or beyond the cycle-prevention bound. Handlers publishing
// derived events inherit the incremented depth.
func (b *EventBus) preparePublish(event Event) error {
	if !b.started.Load() {
		return ErrBusShutdown
	}
	if event == nil {
		return ErrNilEvent
	}
	eventType := event.EventType()
	if eventType == "" {
		return ErrEmptyType
	}

	base := getEventBase(event)
	if base != nil {
		base.depth.Add(1)
	} else {
		event.SetEventDepth(event.EventDepth() + 1)
	}
	depth := event.EventDepth()
	if depth >= b.maxEventDepth {
		b.logWarn("dropping event: max depth exceeded", "event_type", eventType, "depth", depth)
		return ErrMaxDepth
	}
	return nil
}

// PublishSync publishes an event and dispatches it synchronously to all
// matching subscribers, returning the first handler result or error.
func (b *EventBus) PublishSync(ctx context.Context, event Event) (HandlerResult, error) {
	if err := b.preparePublish(event); err != nil {
		return HandlerResult{}, err
	}

	snapshot := b.snapshotSubscriptions(event.EventType())
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].priority > snapshot[j].priority
	})

	var lastResult HandlerResult

	for _, sub := range snapshot {
		result, ok := b.dispatchWithSemaphores(ctx, sub, event)
		if !ok {
			return HandlerResult{}, fmt.Errorf("eventbus: semaphore acquisition failed")
		}

		if result.Error != nil {
			return result, nil
		}
		lastResult = result
	}

	return lastResult, nil
}

// Publish enqueues an event for asynchronous delivery, persisting it to the
// outbox store if one is configured.
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	if err := b.preparePublish(event); err != nil {
		return err
	}
	eventType := event.EventType()

	if b.outbox == nil {
		return b.dispatchAsyncDirect(ctx, event)
	}

	// The async worker deserialises outbox entries via the registry; a
	// nil registry would silently fail every event. Fail fast here so
	// the misconfiguration is visible to the caller.
	if b.registry == nil {
		return fmt.Errorf("eventbus: outbox configured without EventRegistry")
	}

	serialized, err := SerializeEvent(event)
	if err != nil {
		return fmt.Errorf("eventbus: failed to serialize event: %w", err)
	}

	eventID := event.EventID()
	if eventID == "" {
		eventID = generateEventID(eventType)
		if base := getEventBase(event); base != nil {
			base.ID = eventID
		}
	}

	entry := &OutboxEntry{
		EventID:         eventID,
		EventType:       eventType,
		Depth:           event.EventDepth(),
		SerializedEvent: serialized,
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      b.maxRetries,
		HandlerResults:  make(map[string]string),
	}

	// The outbox write deliberately uses a background context: the
	// at-least-once contract requires the event to be persisted even when
	// the publishing caller's context is cancelled right after Publish
	// returns.
	if err := b.outbox.Write(context.Background(), entry); err != nil {
		return fmt.Errorf("eventbus: failed to write to outbox: %w", err)
	}

	select {
	case b.asyncCh <- entry:
	default:
		// The entry stays OutboxPending in the store; the requeue loop
		// re-enqueues it on its next scan.
		b.logWarn("async channel full, entry remains pending for requeue scan", "event_id", entry.EventID)
	}

	return nil
}

func (b *EventBus) dispatchAsyncDirect(ctx context.Context, event Event) error {
	snapshot := b.snapshotSubscriptions(event.EventType())
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].priority > snapshot[j].priority
	})

	for _, sub := range snapshot {
		if !sub.async {
			continue
		}
		select {
		case b.directCh <- &directDispatch{sub: sub, event: event, ctx: ctx}:
		case <-b.stopCh:
			return nil
		}
	}

	return nil
}

func (b *EventBus) dispatchHandler(ctx context.Context, sub *subscriptionEntry, event Event) HandlerResult {
	return sub.handler(ctx, event)
}

// dispatchWithSemaphores acquires the subscription and global concurrency
// slots, runs the handler, and releases the slots even when the handler
// panics; a panic is converted into a HandlerResult error so one faulty
// handler can neither leak the semaphore tokens permanently nor take down
// the worker goroutine. At the point of recovery the goroutine holds only
// the two semaphore slots for this single dispatch, so converting the panic
// is safe. acquired is assigned before the handler runs because a panicking
// handler aborts the evaluation of a return statement, which would leave
// the named results at their zero values.
func (b *EventBus) dispatchWithSemaphores(ctx context.Context, sub *subscriptionEntry, event Event) (result HandlerResult, acquired bool) {
	if !b.acquireSemaphores(sub) {
		return HandlerResult{}, false
	}
	defer b.releaseSemaphores(sub)
	defer func() {
		if r := recover(); r != nil {
			result = HandlerResult{Error: fmt.Errorf("eventbus: handler panicked: %v", r)}
			b.logWarn("event handler panicked", "subscription", sub.id, "event_type", event.EventType(), "panic", r)
		}
	}()
	acquired = true
	result = b.dispatchHandler(ctx, sub, event)
	return result, acquired
}

func (b *EventBus) acquireSemaphores(sub *subscriptionEntry) bool {
	select {
	case b.globalSem <- struct{}{}:
	case <-b.stopCh:
		return false
	}

	if sub.sem != nil {
		select {
		case sub.sem <- struct{}{}:
		case <-b.stopCh:
			<-b.globalSem
			return false
		}
	}

	return true
}

func (b *EventBus) releaseSemaphores(sub *subscriptionEntry) {
	if sub.sem != nil {
		<-sub.sem
	}
	<-b.globalSem
}
