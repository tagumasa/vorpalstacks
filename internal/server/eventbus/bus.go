package eventbus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/logs"
)

var (
	ErrBusShutdown = fmt.Errorf("eventbus: bus is shut down")
	ErrNilEvent    = fmt.Errorf("eventbus: event must not be nil")
	ErrEmptyType   = fmt.Errorf("eventbus: event type must not be empty")
	ErrMaxDepth    = fmt.Errorf("eventbus: max event depth exceeded")
	ErrNoOutbox    = fmt.Errorf("eventbus: outbox store is nil")
	ErrUnknownSub  = fmt.Errorf("eventbus: subscription not found")
)

const (
	DefaultGlobalConcurrency = 256
	DefaultMaxRetries        = 3
	DefaultMaxEventDepth     = 3
	AsyncWorkerCount         = 8
	CleanupInterval          = 10 * time.Minute
	DeliveredRetention       = 1 * time.Hour
	FailedRetention          = 24 * time.Hour
)

type BusOption func(*EventBus)

func WithOutbox(outbox OutboxStore) BusOption {
	return func(b *EventBus) {
		b.outbox = outbox
	}
}

func WithEventRegistry(registry *EventRegistry) BusOption {
	return func(b *EventBus) {
		b.registry = registry
	}
}

func WithRoleResolver(resolver RoleResolver) BusOption {
	return func(b *EventBus) {
		b.roleResolver = resolver
	}
}

func WithPolicyEvaluator(eval BusPolicyEvaluator) BusOption {
	return func(b *EventBus) {
		b.policyEval = eval
	}
}

func WithGlobalConcurrency(n int) BusOption {
	return func(b *EventBus) {
		b.globalSem = make(chan struct{}, n)
	}
}

func WithBusMaxRetries(n int32) BusOption {
	return func(b *EventBus) {
		b.maxRetries = n
	}
}

func WithMaxEventDepth(n int) BusOption {
	return func(b *EventBus) {
		b.maxEventDepth = n
	}
}

func WithLogger(logger logs.Logger) BusOption {
	return func(b *EventBus) {
		b.logger = logger
	}
}

type Bus interface {
	Publish(ctx context.Context, event Event) error
	PublishSync(ctx context.Context, event Event) (HandlerResult, error)
	Subscribe(handler func(ctx context.Context, event Event) HandlerResult, opts ...SubscribeOption) (string, error)
	Unsubscribe(subscriptionID string) error
	RegisterInvoker(invoker ServiceInvoker)
	GetInvoker(serviceType string) (ServiceInvoker, bool)
	EvaluateTargetPolicy(ctx context.Context, targetARN, serviceType, principal, action, resource string) (bool, error)
	RoleResolver() RoleResolver
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type EventBus struct {
	mu            sync.RWMutex
	subscriptions map[string][]*subscriptionEntry
	outbox        OutboxStore
	registry      *EventRegistry
	roleResolver  RoleResolver
	policyEval    BusPolicyEvaluator
	policyFuncs   map[string]ResourcePolicyFunc
	policyFuncsMu sync.RWMutex
	globalSem     chan struct{}
	maxRetries    int32
	maxEventDepth int
	logger        logs.Logger
	wg            sync.WaitGroup
	started       atomic.Bool
	stopCh        chan struct{}
	invokers      map[string]ServiceInvoker
	invokersMu    sync.RWMutex
	nextSubID     atomic.Int64
	asyncCh       chan *OutboxEntry
}

func NewEventBus(opts ...BusOption) *EventBus {
	b := &EventBus{
		subscriptions: make(map[string][]*subscriptionEntry),
		globalSem:     make(chan struct{}, DefaultGlobalConcurrency),
		maxRetries:    DefaultMaxRetries,
		maxEventDepth: DefaultMaxEventDepth,
		stopCh:        make(chan struct{}),
		invokers:      make(map[string]ServiceInvoker),
		asyncCh:       make(chan *OutboxEntry, 1024),
		policyFuncs:   make(map[string]ResourcePolicyFunc),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

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

func SubscribeTyped[T Event](bus *EventBus, handler func(ctx context.Context, event T) HandlerResult, opts ...SubscribeOption) (string, error) {
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

	wrappedHandler := func(ctx context.Context, event Event) HandlerResult {
		typed, ok := event.(T)
		if !ok {
			return HandlerResult{Error: fmt.Errorf("eventbus: type mismatch for SubscribeTyped")}
		}
		return handler(ctx, typed)
	}

	return bus.subscribeInternal(cfg, wrappedHandler)
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

	entry := &subscriptionEntry{
		id:             subID,
		eventType:      "*",
		filter:         cfg.filter,
		handler:        handler,
		async:          cfg.async,
		priority:       cfg.priority,
		maxConcurrency: cfg.maxConcurrency,
		maxRetries:     cfg.maxRetries,
		sem:            sem,
		authorized:     authorized,
	}

	b.mu.Lock()
	b.subscriptions["*"] = append(b.subscriptions["*"], entry)
	b.mu.Unlock()

	return subID, nil
}

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

func (b *EventBus) PublishSync(ctx context.Context, event Event) (HandlerResult, error) {
	if !b.started.Load() {
		return HandlerResult{}, ErrBusShutdown
	}
	if event == nil {
		return HandlerResult{}, ErrNilEvent
	}
	eventType := event.EventType()
	if eventType == "" {
		return HandlerResult{}, ErrEmptyType
	}

	base := getEventBase(event)
	if base != nil {
		base.Depth++
	} else {
		event.SetEventDepth(event.EventDepth() + 1)
	}
	depth := event.EventDepth()
	if depth >= b.maxEventDepth {
		b.logWarn("dropping event: max depth exceeded", "event_type", eventType, "depth", depth)
		return HandlerResult{}, ErrMaxDepth
	}

	snapshot := b.snapshotSubscriptions(eventType)
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].priority > snapshot[j].priority
	})

	for _, sub := range snapshot {
		if !sub.authorized {
			continue
		}
		if sub.filter != nil && !sub.filter.Match(event) {
			continue
		}

		if !b.acquireSemaphores(sub) {
			return HandlerResult{}, fmt.Errorf("eventbus: semaphore acquisition failed")
		}

		result := b.dispatchHandler(ctx, sub, event)

		b.releaseSemaphores(sub)

		if result.Error != nil {
			return result, nil
		}
	}

	return HandlerResult{}, nil
}

func (b *EventBus) Publish(ctx context.Context, event Event) error {
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
		base.Depth++
	} else {
		event.SetEventDepth(event.EventDepth() + 1)
	}
	depth := event.EventDepth()
	if depth >= b.maxEventDepth {
		b.logWarn("dropping event: max depth exceeded", "event_type", eventType, "depth", depth)
		return ErrMaxDepth
	}

	if b.outbox == nil {
		return b.dispatchAsyncDirect(ctx, event)
	}

	serialized, err := SerializeEvent(event)
	if err != nil {
		return fmt.Errorf("eventbus: failed to serialize event: %w", err)
	}

	eventID := event.EventID()
	if eventID == "" {
		eventID = generateEventID(eventType)
		if base != nil {
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

	if err := b.outbox.Write(context.Background(), entry); err != nil {
		return fmt.Errorf("eventbus: failed to write to outbox: %w", err)
	}

	select {
	case b.asyncCh <- entry:
	default:
		b.logWarn("async channel full, event will be recovered from outbox", "event_id", entry.EventID)
	}

	return nil
}

func (b *EventBus) dispatchAsyncDirect(ctx context.Context, event Event) error {
	snapshot := b.snapshotSubscriptions(event.EventType())
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].priority > snapshot[j].priority
	})

	for _, sub := range snapshot {
		if !sub.authorized || !sub.async {
			continue
		}
		b.wg.Add(1)
		go func(s *subscriptionEntry) {
			defer b.wg.Done()
			if !b.acquireSemaphores(s) {
				return
			}
			b.dispatchHandler(ctx, s, event)
			b.releaseSemaphores(s)
		}(sub)
	}

	return nil
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

func (b *EventBus) RegisterInvoker(invoker ServiceInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.invokers[invoker.ServiceType()] = invoker
}

func (b *EventBus) GetInvoker(serviceType string) (ServiceInvoker, bool) {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	inv, ok := b.invokers[serviceType]
	return inv, ok
}

func (b *EventBus) RoleResolver() RoleResolver {
	return b.roleResolver
}

func (b *EventBus) Start(ctx context.Context) error {
	if b.outbox != nil {
		if err := b.recover(ctx); err != nil {
			return fmt.Errorf("eventbus: recovery failed: %w", err)
		}
	}

	for i := 0; i < AsyncWorkerCount; i++ {
		b.wg.Add(1)
		go b.asyncWorker()
	}

	if b.outbox != nil {
		b.wg.Add(1)
		go b.cleanupLoop()
	}

	b.started.Store(true)
	return nil
}

func (b *EventBus) Shutdown(ctx context.Context) error {
	b.started.Store(false)
	close(b.stopCh)

	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if b.outbox != nil {
			if err := b.outbox.Close(); err != nil {
				b.logWarn("outbox close error", "error", err)
			}
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *EventBus) recover(ctx context.Context) error {
	pending, err := b.outbox.ListPending(ctx, 1000)
	if err != nil {
		return err
	}

	for _, entry := range pending {
		if entry.Depth >= b.maxEventDepth {
			entry.Status = OutboxFailed
			entry.LastError = "max depth exceeded on recovery"
			_ = b.outbox.UpdateEntry(ctx, entry)
			continue
		}

		if b.registry != nil {
			event, err := b.registry.Deserialize(entry.EventType, entry.SerializedEvent)
			if err != nil {
				b.logWarn("failed to deserialize outbox entry on recovery", "event_id", entry.EventID, "error", err)
				entry.Status = OutboxFailed
				entry.LastError = err.Error()
				_ = b.outbox.UpdateEntry(ctx, entry)
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

	return nil
}

func (b *EventBus) asyncWorker() {
	defer b.wg.Done()
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

func (b *EventBus) processOutboxEntry(entry *OutboxEntry) {
	ctx := context.Background()

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

	allSkipped := true
	for _, sub := range snapshot {
		if !sub.authorized {
			continue
		}
		if sub.filter != nil && !sub.filter.Match(event) {
			continue
		}
		allSkipped = false

		if !b.acquireSemaphores(sub) {
			entry.HandlerResults[sub.id] = "skipped"
			continue
		}

		result := b.dispatchHandler(ctx, sub, event)
		b.releaseSemaphores(sub)

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
		_ = b.outbox.UpdateEntry(ctx, entry)
		return
	}

	hasFailure := false
	for _, v := range entry.HandlerResults {
		if v == "failed" {
			hasFailure = true
			break
		}
	}

	if hasFailure {
		entry.RetryCount++
		if entry.RetryCount >= entry.MaxRetries {
			entry.Status = OutboxFailed
			_ = b.outbox.UpdateEntry(ctx, entry)
		} else {
			entry.Status = OutboxPending
			_ = b.outbox.UpdateEntry(ctx, entry)
			select {
			case b.asyncCh <- entry:
			default:
			}
		}
	} else {
		now := time.Now().UTC()
		entry.Status = OutboxDelivered
		entry.DeliveredAt = &now
		_ = b.outbox.UpdateEntry(ctx, entry)
	}
}

func (b *EventBus) failEntry(ctx context.Context, entry *OutboxEntry, reason string) {
	entry.Status = OutboxFailed
	entry.LastError = reason
	_ = b.outbox.UpdateEntry(ctx, entry)
}

func (b *EventBus) deserializeEntry(entry *OutboxEntry) (Event, error) {
	if b.registry == nil {
		return nil, fmt.Errorf("eventbus: no event registry configured")
	}
	return b.registry.Deserialize(entry.EventType, entry.SerializedEvent)
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

func (b *EventBus) dispatchHandler(ctx context.Context, sub *subscriptionEntry, event Event) HandlerResult {
	switch h := sub.handler.(type) {
	case func(context.Context, Event) HandlerResult:
		return h(ctx, event)
	case func(context.Context, *ServiceInvokeRequest) HandlerResult:
		if typed, ok := event.(*ServiceInvokeRequest); ok {
			return h(ctx, typed)
		}
		return HandlerResult{Error: fmt.Errorf("eventbus: type assertion failed for *ServiceInvokeRequest")}
	default:
		return HandlerResult{Error: fmt.Errorf("eventbus: unsupported handler type %T", sub.handler)}
	}
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

func (b *EventBus) cleanupLoop() {
	defer b.wg.Done()
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

func getEventBase(event Event) *EventBase {
	if base, ok := event.(*EventBase); ok {
		return base
	}
	if base, ok := event.(interface{ getEventBase() *EventBase }); ok {
		return base.getEventBase()
	}
	return nil
}

func generateEventID(eventType string) string {
	return fmt.Sprintf("%s-%d", eventType, time.Now().UnixNano())
}

func (b *EventBus) logInfo(msg string, keyvals ...interface{}) {
	if b.logger != nil {
		fields := make([]logs.Field, 0, len(keyvals)/2)
		for i := 0; i+1 < len(keyvals); i += 2 {
			fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
		}
		b.logger.Info(msg, fields...)
	}
}

func (b *EventBus) logWarn(msg string, keyvals ...interface{}) {
	if b.logger != nil {
		fields := make([]logs.Field, 0, len(keyvals)/2)
		for i := 0; i+1 < len(keyvals); i += 2 {
			fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
		}
		b.logger.Warn(msg, fields...)
	}
}

var _ Bus = (*EventBus)(nil)
var _ error = ErrBusShutdown
var _ error = ErrNilEvent
var _ error = ErrEmptyType
var _ error = ErrMaxDepth
var _ error = ErrNoOutbox
var _ error = ErrUnknownSub
var _ = errors.Is
