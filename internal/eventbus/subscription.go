package eventbus

// EventFilter defines a type-safe predicate for filtering events of type T.
type EventFilter[T Event] interface {
	Matches(event T) bool
}

type filterMatcher interface {
	Match(event Event) bool
}

type filterAdapter[T Event] struct {
	inner EventFilter[T]
}

// Match evaluates the type-unsafe event against the underlying typed filter.
func (a *filterAdapter[T]) Match(e Event) bool {
	typed, ok := e.(T)
	if !ok {
		return false
	}
	return a.inner.Matches(typed)
}

type subscriptionEntry struct {
	id             string
	eventType      string
	filter         filterMatcher
	handler        any
	async          bool
	priority       int
	maxConcurrency int
	maxRetries     int32
	sem            chan struct{}
	authorized     bool
}

// AuthzMode controls the level of authorisation enforcement for a subscription.
type AuthzMode int

const (
	// AuthzNone disables authorisation checks for this subscription.
	AuthzNone AuthzMode = iota
	// AuthzResourcePolicy enables resource-based policy evaluation.
	AuthzResourcePolicy
	// AuthzFull enables both resource policy evaluation and IAM role validation.
	AuthzFull
)

type subscribeConfig struct {
	eventType        string
	async            bool
	priority         int
	filter           filterMatcher
	callerPrincipal  string
	targetRoleARN    string
	resourcePolicyFn ResourcePolicyFunc
	authzMode        AuthzMode
	maxConcurrency   int
	maxRetries       int32
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

// WithFilter attaches a type-safe event filter to the subscription.
func WithFilter[T Event](f EventFilter[T]) SubscribeOption {
	return func(c *subscribeConfig) {
		c.filter = &filterAdapter[T]{inner: f}
	}
}

// WithCallerPrincipal sets the IAM principal ARN for resource policy evaluation.
func WithCallerPrincipal(principal string) SubscribeOption {
	return func(c *subscribeConfig) {
		c.callerPrincipal = principal
	}
}

// WithTargetRole sets the IAM role ARN to validate at subscription time.
func WithTargetRole(roleARN string) SubscribeOption {
	return func(c *subscribeConfig) {
		c.targetRoleARN = roleARN
	}
}

// WithResourcePolicy sets the function used to fetch the target's
// resource-based policy for authorisation evaluation.
func WithResourcePolicy(fn ResourcePolicyFunc) SubscribeOption {
	return func(c *subscribeConfig) {
		c.resourcePolicyFn = fn
	}
}

// WithAuthzMode sets the authorisation enforcement level for the subscription.
func WithAuthzMode(mode AuthzMode) SubscribeOption {
	return func(c *subscribeConfig) {
		c.authzMode = mode
	}
}

// WithMaxConcurrency limits the number of concurrent dispatches for this
// subscription using a per-subscription semaphore.
func WithMaxConcurrency(n int) SubscribeOption {
	return func(c *subscribeConfig) {
		c.maxConcurrency = n
	}
}

// WithMaxRetries sets the maximum number of retry attempts for this
// subscription's outbox entries.
func WithMaxRetries(n int32) SubscribeOption {
	return func(c *subscribeConfig) {
		c.maxRetries = n
	}
}
