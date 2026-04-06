package eventbus

type EventFilter[T Event] interface {
	Matches(event T) bool
}

type filterMatcher interface {
	Match(event Event) bool
}

type filterAdapter[T Event] struct {
	inner EventFilter[T]
}

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

type AuthzMode int

const (
	AuthzNone AuthzMode = iota
	AuthzResourcePolicy
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

type SubscribeOption func(*subscribeConfig)

func WithAsync() SubscribeOption {
	return func(c *subscribeConfig) {
		c.async = true
	}
}

func WithPriority(p int) SubscribeOption {
	return func(c *subscribeConfig) {
		c.priority = p
	}
}

func WithFilter[T Event](f EventFilter[T]) SubscribeOption {
	return func(c *subscribeConfig) {
		c.filter = &filterAdapter[T]{inner: f}
	}
}

func WithCallerPrincipal(principal string) SubscribeOption {
	return func(c *subscribeConfig) {
		c.callerPrincipal = principal
	}
}

func WithTargetRole(roleARN string) SubscribeOption {
	return func(c *subscribeConfig) {
		c.targetRoleARN = roleARN
	}
}

func WithResourcePolicy(fn ResourcePolicyFunc) SubscribeOption {
	return func(c *subscribeConfig) {
		c.resourcePolicyFn = fn
	}
}

func WithAuthzMode(mode AuthzMode) SubscribeOption {
	return func(c *subscribeConfig) {
		c.authzMode = mode
	}
}

func WithMaxConcurrency(n int) SubscribeOption {
	return func(c *subscribeConfig) {
		c.maxConcurrency = n
	}
}

func WithMaxRetries(n int32) SubscribeOption {
	return func(c *subscribeConfig) {
		c.maxRetries = n
	}
}
