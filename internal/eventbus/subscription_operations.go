package eventbus

import (
	"context"
	"fmt"
)

// Subscribe registers a handler for events, returning a subscription ID.
func (b *EventBus) Subscribe(handler func(ctx context.Context, event Event) HandlerResult, opts ...SubscribeOption) (string, error) {
	if handler == nil {
		return "", fmt.Errorf("eventbus: handler must not be nil")
	}

	cfg := &subscribeConfig{
		authzMode: AuthzNone,
		async:     false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return b.subscribeInternal(cfg, handler)
}

// SubscribeTyped registers a type-safe handler that receives only events of
// type T, automatically filtering mismatched events.
func SubscribeTyped[T Event](bus Bus, handler func(ctx context.Context, event T) HandlerResult, opts ...SubscribeOption) (string, error) {
	if handler == nil {
		return "", fmt.Errorf("eventbus: handler must not be nil")
	}

	var zero T
	allOpts := make([]SubscribeOption, 0, len(opts)+2)
	allOpts = append(allOpts, WithEventType(zero.EventType()))
	allOpts = append(allOpts, opts...)

	wrappedHandler := func(ctx context.Context, event Event) HandlerResult {
		typed, ok := event.(T)
		if !ok {
			return HandlerResult{Error: fmt.Errorf("eventbus: type mismatch for SubscribeTyped")}
		}
		return handler(ctx, typed)
	}

	return bus.Subscribe(wrappedHandler, allOpts...)
}

func (b *EventBus) subscribeInternal(cfg *subscribeConfig, handler func(ctx context.Context, event Event) HandlerResult) (string, error) {

	if cfg.authzMode == AuthzFull && cfg.targetRoleARN != "" && b.roleResolver != nil {
		if err := b.roleResolver.ValidateRole(context.Background(), cfg.targetRoleARN); err != nil {
			return "", fmt.Errorf("eventbus: role validation failed for %q: %w", cfg.targetRoleARN, err)
		}
	}

	if cfg.authzMode >= AuthzResourcePolicy && cfg.resourcePolicyFn != nil && b.policyEval != nil {
		if cfg.callerPrincipal == "" {
			return "", fmt.Errorf("eventbus: caller principal required for resource policy evaluation")
		}
	}

	// Resource policy evaluation at subscribe time is limited: the bus
	// subscribes to event types (e.g. "sns:deliver"), not to specific
	// targets. The actual target ARN is only known at dispatch time when
	// the handler looks up notification configurations. Therefore the
	// resource policy check here validates the caller principal against
	// the source resource's own policy (if any), while target-specific
	// policy evaluation happens at dispatch time via
	// EvaluateTargetPolicy. This is an intentional architectural
	// exception to decision #3 in the plan.
	authorized := true
	if cfg.authzMode >= AuthzResourcePolicy && cfg.resourcePolicyFn != nil && b.policyEval != nil {
		policyDoc, err := cfg.resourcePolicyFn(context.Background(), "")
		if err != nil {
			return "", fmt.Errorf("eventbus: failed to fetch resource policy: %w", err)
		}
		if policyDoc != nil && len(policyDoc.Statement) > 0 {
			allowed, err := b.policyEval.Evaluate(context.Background(), policyDoc, cfg.callerPrincipal, "eventbus:Subscribe", "*")
			if err != nil {
				return "", fmt.Errorf("eventbus: policy evaluation failed: %w", err)
			}
			if !allowed {
				authorized = false
			}
		}
	}

	subID := fmt.Sprintf("sub-%d", b.nextSubID.Add(1))

	var sem chan struct{}
	if cfg.maxConcurrency > 0 {
		sem = make(chan struct{}, cfg.maxConcurrency)
	}

	et := cfg.eventType
	if et == "" {
		et = "*"
	}

	entry := &subscriptionEntry{
		id:             subID,
		eventType:      et,
		filter:         cfg.filter,
		handler:        handler,
		async:          cfg.async,
		priority:       cfg.priority,
		maxConcurrency: cfg.maxConcurrency,
		sem:            sem,
		authorized:     authorized,
	}

	b.mu.Lock()
	b.subscriptions[et] = append(b.subscriptions[et], entry)
	b.mu.Unlock()

	return subID, nil
}

// Unsubscribe removes a previously registered subscription by its ID.
func (b *EventBus) Unsubscribe(subscriptionID string) error {
	if subscriptionID == "" {
		return ErrUnknownSub
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for eventType, entries := range b.subscriptions {
		for i, entry := range entries {
			if entry.id == subscriptionID {
				b.subscriptions[eventType] = append(entries[:i], entries[i+1:]...)
				if len(b.subscriptions[eventType]) == 0 {
					delete(b.subscriptions, eventType)
				}
				return nil
			}
		}
	}

	return ErrUnknownSub
}

func (b *EventBus) snapshotSubscriptions(eventType string) []*subscriptionEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	typeEntries := b.subscriptions[eventType]
	wildcardEntries := b.subscriptions["*"]

	total := len(typeEntries) + len(wildcardEntries)
	if total == 0 {
		return nil
	}

	snapshot := make([]*subscriptionEntry, 0, total)
	snapshot = append(snapshot, typeEntries...)
	snapshot = append(snapshot, wildcardEntries...)
	return snapshot
}

// SetResourcePolicyFunc registers a ResourcePolicyFunc for the given
// service type (e.g. "lambda", "sqs", "sns", "events"). This function
// is called by EvaluateTargetPolicy to fetch the target's resource-based
// policy before dispatch. Must be called before any events are published.
func (b *EventBus) SetResourcePolicyFunc(serviceType string, fn ResourcePolicyFunc) {
	b.policyFuncsMu.Lock()
	defer b.policyFuncsMu.Unlock()
	b.policyFuncs[serviceType] = fn
}

// EvaluateTargetPolicy checks whether the specified principal is authorised
// to perform the given action on the target resource according to the
// target's resource-based policy.
//
// This is the primary authorisation point for cross-service delivery. Role
// validation happens at configuration time (PutTargets, PutRule, etc.) via
// RoleResolver; resource policy evaluation happens here at dispatch time
// because the actual target ARN is only known when the handler processes
// the event.
//
// Returns (true, nil) when:
//   - no ResourcePolicyFunc is registered for the service type (opt-in);
//   - the target has no policy (open resource); or
//   - the policy explicitly allows the action.
//
// Returns (false, nil) when the policy explicitly denies the action or
// no matching Allow statement is found (default-deny). Returns (false, err)
// when the policy cannot be fetched or parsed (fail-closed).
func (b *EventBus) EvaluateTargetPolicy(ctx context.Context, targetARN, serviceType, principal, action, resource string) (bool, error) {
	b.policyFuncsMu.RLock()
	fn, ok := b.policyFuncs[serviceType]
	b.policyFuncsMu.RUnlock()

	if !ok || fn == nil {
		return true, nil
	}

	policyDoc, err := fn(ctx, targetARN)
	if err != nil {
		b.logWarn("failed to fetch resource policy", "target_arn", targetARN, "error", err)
		return false, nil
	}

	if policyDoc == nil || len(policyDoc.Statement) == 0 {
		return true, nil
	}

	if b.policyEval == nil {
		return true, nil
	}

	allowed, err := b.policyEval.Evaluate(ctx, policyDoc, principal, action, resource)
	if err != nil {
		b.logWarn("policy evaluation failed", "target_arn", targetARN, "error", err)
		return false, nil
	}

	return allowed, nil
}
