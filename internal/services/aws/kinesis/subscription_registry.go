package kinesis

import (
	"context"
	"sync"
	"time"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// subscriptionRegistry tracks the live SubscribeToShard pumps by
// (consumer ARN, shard ID) pair and enforces the subscription rules the
// SubscribeToShard documentation states: a repeat call within the takeover
// window of a successful subscription answers ResourceInUseException, and
// a call at or beyond the window takes the subscription over — the
// previous connection is cancelled and expires.
type subscriptionRegistry struct {
	mu      sync.Mutex
	entries map[string]*subscription

	// ttlOverride exists so tests observe the expiry rule without waiting
	// out the production window; zero selects the documented constant.
	ttlOverride time.Duration

	// windowOverride exists so tests observe the tombstone self-removal
	// without waiting out the production takeover window; zero selects the
	// documented constant.
	windowOverride time.Duration
}

// subscription is one registered event pump. All field writes happen
// under the registry mutex — including the cancel function the owning
// pump installs right after begin returns it.
type subscription struct {
	startedAt time.Time
	cancel    context.CancelFunc

	// ended marks the pump's exit. The entry stays registered as a
	// tombstone for the rest of the takeover window: the documented
	// rule anchors the window to the successful call ("within 5 seconds
	// of a successful call"), not to the connection's lifetime.
	ended bool
}

// subscriptionKey builds the registry key for one consumer-shard pair.
// The consumer ARN is region-qualified, so keys cannot collide across
// regions served by one process.
func subscriptionKey(consumerARN, shardID string) string {
	return consumerARN + "#" + shardID
}

// begin registers a subscription for key, creating its cancellable
// context from parent. It reports ok=false when an entry for the same
// pair began within the takeover window — live or tombstoned, the caller
// answers ResourceInUseException per the documentation. An entry at or
// beyond the window is taken over: a still-live pump's cancel function
// runs under the lock and its stream ends through its own deferred
// cleanup, and the new subscription replaces the entry. Ended entries
// past their window are swept here — the only walk the registry needs,
// since every new subscription passes through this door.
func (r *subscriptionRegistry) begin(key string, parent context.Context) (context.Context, *subscription, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.entries == nil {
		r.entries = make(map[string]*subscription)
	}
	now := time.Now()
	for k, e := range r.entries {
		if e.ended && now.Sub(e.startedAt) >= r.takeoverWindow() {
			delete(r.entries, k)
		}
	}
	if existing, ok := r.entries[key]; ok {
		if now.Sub(existing.startedAt) < r.takeoverWindow() {
			return nil, nil, false
		}
		if existing.cancel != nil {
			existing.cancel()
		}
		delete(r.entries, key)
	}
	sub := &subscription{startedAt: now}
	ctx, cancel := context.WithCancel(parent)
	sub.cancel = cancel
	r.entries[key] = sub
	return ctx, sub, true
}

// end marks the subscription's pump stopped. The entry survives as a
// tombstone for the takeover window; end is idempotent and safe to call
// for a subscription a takeover already replaced — the surviving entry is
// the replacement's own.
func (r *subscriptionRegistry) end(key string, sub *subscription) {
	r.mu.Lock()
	if r.entries[key] == sub {
		sub.ended = true
	}
	r.mu.Unlock()
	// A pair never subscribed to again would keep its tombstone forever —
	// begin's sweep only runs when some other subscription starts. The
	// entry removes itself once the window has definitely elapsed; the
	// identity check keeps a takeover's replacement registered.
	time.AfterFunc(r.takeoverWindow(), func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		if e, ok := r.entries[key]; ok && e == sub && e.ended &&
			time.Since(e.startedAt) >= r.takeoverWindow() {
			delete(r.entries, key)
		}
	})
}

// discard removes a registration whose subscription never became a pump —
// a failed call leaves no tombstone, because the takeover window the
// documentation states anchors to a successful call and this one never
// succeeded. The cancelled context releases anything already holding it.
func (r *subscriptionRegistry) discard(key string, sub *subscription) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.entries[key] == sub {
		if sub.cancel != nil {
			sub.cancel()
		}
		delete(r.entries, key)
	}
}

// takeoverWindow is the documented re-subscription window.
func (r *subscriptionRegistry) takeoverWindow() time.Duration {
	if r.windowOverride > 0 {
		return r.windowOverride
	}
	return kinesisstore.SubscribeTakeoverWindow
}

// lifetime is the documented subscription duration unless a test
// shortened it.
func (r *subscriptionRegistry) lifetime() time.Duration {
	if r.ttlOverride > 0 {
		return r.ttlOverride
	}
	return kinesisstore.SubscribeTTL
}
