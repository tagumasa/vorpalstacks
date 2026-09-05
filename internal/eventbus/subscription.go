package eventbus

import "context"

type subscriptionEntry struct {
	id             string
	eventType      string
	handler        func(ctx context.Context, event Event) HandlerResult
	async          bool
	priority       int
	maxConcurrency int
	sem            chan struct{}
}

type subscribeConfig struct {
	eventType      string
	async          bool
	priority       int
	maxConcurrency int
}

// SubscribeOption is a functional option used to configure a subscription.
type SubscribeOption func(*subscribeConfig)

// WithAsync configures the subscription for asynchronous delivery.
func WithAsync() SubscribeOption {
	return func(c *subscribeConfig) {
		c.async = true
	}
}

// WithPriority sets the dispatch priority for the subscription; higher
// values are dispatched first.
func WithPriority(p int) SubscribeOption {
	return func(c *subscribeConfig) {
		c.priority = p
	}
}

// WithMaxConcurrency limits the number of concurrent dispatches for this
// subscription using a per-subscription semaphore.
func WithMaxConcurrency(n int) SubscribeOption {
	return func(c *subscribeConfig) {
		c.maxConcurrency = n
	}
}

// WithEventType sets the event type string for subscription routing.
// When set, the subscription is stored under the specific event type key
// for efficient dispatch lookup.
func WithEventType(eventType string) SubscribeOption {
	return func(c *subscribeConfig) {
		c.eventType = eventType
	}
}
