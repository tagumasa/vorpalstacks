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
	"vorpalstacks/internal/core/resilience"
)

var (
	// ErrBusShutdown is returned when publishing to a bus that has been shut down.
	ErrBusShutdown = fmt.Errorf("eventbus: bus is shut down")
	// ErrNilEvent is returned when a nil event is published.
	ErrNilEvent = fmt.Errorf("eventbus: event must not be nil")
	// ErrEmptyType is returned when an event with an empty type string is published.
	ErrEmptyType = fmt.Errorf("eventbus: event type must not be empty")
	// ErrMaxDepth is returned when an event exceeds the maximum allowed dispatch depth.
	ErrMaxDepth = fmt.Errorf("eventbus: max event depth exceeded")
	// ErrNoOutbox is returned when an outbox operation is attempted without a configured store.
	ErrNoOutbox = fmt.Errorf("eventbus: outbox store is nil")
	// ErrUnknownSub is returned when unsubscribing a non-existent subscription ID.
	ErrUnknownSub = fmt.Errorf("eventbus: subscription not found")
)

const (
	// DefaultGlobalConcurrency is the default limit on concurrent handler executions.
	DefaultGlobalConcurrency = 256
	// DefaultMaxRetries is the default number of retry attempts for failed outbox entries.
	DefaultMaxRetries = 3
	// DefaultMaxEventDepth is the default maximum dispatch depth before events are dropped.
	DefaultMaxEventDepth = 3
	// AsyncWorkerCount is the number of goroutines processing the outbox async channel.
	AsyncWorkerCount = 8
	// CleanupInterval is the period between outbox retention cleanups.
	CleanupInterval = 10 * time.Minute
	// DeliveredRetention is how long successfully delivered outbox entries are kept.
	DeliveredRetention = 1 * time.Hour
	// FailedRetention is how long failed outbox entries are kept before purging.
	FailedRetention = 24 * time.Hour
)

type directDispatch struct {
	sub   *subscriptionEntry
	event Event
	ctx   context.Context
}

// BusOption is a functional option used to configure an EventBus.
type BusOption func(*EventBus)

// WithOutbox configures the EventBus to use the given OutboxStore for
// persistent delivery tracking.
func WithOutbox(outbox OutboxStore) BusOption {
	return func(b *EventBus) {
		b.outbox = outbox
	}
}

// WithEventRegistry configures the EventBus to use the given EventRegistry
// for event deserialisation.
func WithEventRegistry(registry *EventRegistry) BusOption {
	return func(b *EventBus) {
		b.registry = registry
	}
}

// WithRoleResolver configures the EventBus to use the given RoleResolver
// for IAM role validation at subscription time.
func WithRoleResolver(resolver RoleResolver) BusOption {
	return func(b *EventBus) {
		b.roleResolver = resolver
	}
}

// WithPolicyEvaluator configures the EventBus to use the given policy
// evaluator for resource-based policy checks at dispatch time.
func WithPolicyEvaluator(eval BusPolicyEvaluator) BusOption {
	return func(b *EventBus) {
		b.policyEval = eval
	}
}

// WithGlobalConcurrency sets the maximum number of concurrent handler
// executions across all subscriptions.
func WithGlobalConcurrency(n int) BusOption {
	return func(b *EventBus) {
		b.globalSem = make(chan struct{}, n)
	}
}

// WithBusMaxRetries sets the default maximum retry attempts for outbox entries.
func WithBusMaxRetries(n int32) BusOption {
	return func(b *EventBus) {
		b.maxRetries = n
	}
}

// WithMaxEventDepth sets the maximum dispatch depth before events are
// silently dropped to prevent infinite cycles.
func WithMaxEventDepth(n int) BusOption {
	return func(b *EventBus) {
		b.maxEventDepth = n
	}
}

// WithLogger configures the EventBus to use the given structured logger.
func WithLogger(logger logs.Logger) BusOption {
	return func(b *EventBus) {
		b.logger = logger
	}
}

// Bus defines the contract for the internal service event bus, supporting
// both synchronous and asynchronous publish, subscription management, and
// cross-service authorisation.
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

	LambdaInvoker() LambdaInvoker
	SQSInvoker() SQSInvoker
	SNSInvoker() SNSInvoker
	KinesisInvoker() KinesisInvoker
	EventsInvoker() EventsInvoker
	EC2Invoker() EC2Invoker
	DynamoDBInvoker() DynamoDBInvoker
	NeptuneGraphInvoker() NeptuneGraphInvoker
	KMSInvoker() KMSInvoker

	SetLambdaInvoker(invoker LambdaInvoker)
	SetSQSInvoker(invoker SQSInvoker)
	SetSNSInvoker(invoker SNSInvoker)
	SetKinesisInvoker(invoker KinesisInvoker)
	SetEC2Invoker(invoker EC2Invoker)
	SetEventsInvoker(invoker EventsInvoker)
	SetDynamoDBInvoker(invoker DynamoDBInvoker)
	SetNeptuneGraphInvoker(invoker NeptuneGraphInvoker)
	SetKMSInvoker(invoker KMSInvoker)
}

// EventBus is the central implementation of the Bus interface, managing
// subscriptions, outbox persistence, async workers, and invoker dispatch.
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
	stopOnce      sync.Once
	stopCh        chan struct{}
	invokers      map[string]ServiceInvoker
	invokersMu    sync.RWMutex

	lambdaInvoker       LambdaInvoker
	sqsInvoker          SQSInvoker
	snsInvoker          SNSInvoker
	kinesisInvoker      KinesisInvoker
	eventsInvoker       EventsInvoker
	ec2Invoker          EC2Invoker
	dynamoDBInvoker     DynamoDBInvoker
	neptuneGraphInvoker NeptuneGraphInvoker
	kmsInvoker          KMSInvoker
	nextSubID           atomic.Int64
	asyncCh             chan *OutboxEntry
	directCh            chan *directDispatch
}

// NewEventBus creates a new EventBus with sensible defaults, applying all
// provided functional options.
func NewEventBus(opts ...BusOption) *EventBus {
	b := &EventBus{
		subscriptions: make(map[string][]*subscriptionEntry),
		globalSem:     make(chan struct{}, DefaultGlobalConcurrency),
		maxRetries:    DefaultMaxRetries,
		maxEventDepth: DefaultMaxEventDepth,
		stopCh:        make(chan struct{}),
		invokers:      make(map[string]ServiceInvoker),
		asyncCh:       make(chan *OutboxEntry, 1024),
		directCh:      make(chan *directDispatch, 1024),
		policyFuncs:   make(map[string]ResourcePolicyFunc),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

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
		maxRetries:     cfg.maxRetries,
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

// PublishSync publishes an event and dispatches it synchronously to all
// matching subscribers, returning the first handler result or error.
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
		base.depth.Add(1)
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

	var lastResult HandlerResult

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
		lastResult = result
	}

	return lastResult, nil
}

// Publish enqueues an event for asynchronous delivery, persisting it to the
// outbox store if one is configured.
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
		base.depth.Add(1)
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
		select {
		case b.directCh <- &directDispatch{sub: sub, event: event, ctx: ctx}:
		case <-b.stopCh:
			b.wg.Done()
		}
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

// RegisterInvoker registers a ServiceInvoker for the given service type.
func (b *EventBus) RegisterInvoker(invoker ServiceInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.invokers[invoker.ServiceType()] = invoker
}

// GetInvoker retrieves the ServiceInvoker registered for the given service type.
func (b *EventBus) GetInvoker(serviceType string) (ServiceInvoker, bool) {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	inv, ok := b.invokers[serviceType]
	return inv, ok
}

// SetLambdaInjector sets the Lambda invoker used for dispatching Lambda function
// invocations from bus events.
func (b *EventBus) SetLambdaInvoker(invoker LambdaInvoker) { b.lambdaInvoker = invoker }

// LambdaInvoker returns the configured Lambda invoker.
func (b *EventBus) LambdaInvoker() LambdaInvoker { return b.lambdaInvoker }

// SetSQSInvoker sets the SQS invoker used for dispatching SQS SendMessage calls
// from bus events.
func (b *EventBus) SetSQSInvoker(invoker SQSInvoker) { b.sqsInvoker = invoker }

// SQSInvoker returns the configured SQS invoker.
func (b *EventBus) SQSInvoker() SQSInvoker { return b.sqsInvoker }

// SetSNSInvoker sets the SNS invoker used for dispatching SNS Publish calls
// from bus events.
func (b *EventBus) SetSNSInvoker(invoker SNSInvoker) { b.snsInvoker = invoker }

// SNSInvoker returns the configured SNS invoker.
func (b *EventBus) SNSInvoker() SNSInvoker { return b.snsInvoker }

// SetKinesisInvoker sets the Kinesis invoker used for dispatching Kinesis
// PutRecord calls from bus events.
func (b *EventBus) SetKinesisInvoker(invoker KinesisInvoker) { b.kinesisInvoker = invoker }

// KinesisInvoker returns the configured Kinesis invoker.
func (b *EventBus) KinesisInvoker() KinesisInvoker { return b.kinesisInvoker }

// SetEventsInvoker sets the EventBridge invoker used for dispatching
// EventBridge PutEvents calls from bus events.
func (b *EventBus) SetEventsInvoker(invoker EventsInvoker) { b.eventsInvoker = invoker }

// EventsInvoker returns the configured EventBridge invoker.
func (b *EventBus) EventsInvoker() EventsInvoker { return b.eventsInvoker }

// SetEC2Invoker sets the EC2 invoker used for dispatching EC2 API calls
// from bus events.
func (b *EventBus) SetEC2Invoker(invoker EC2Invoker) { b.ec2Invoker = invoker }

// EC2Invoker returns the configured EC2 invoker.
func (b *EventBus) EC2Invoker() EC2Invoker { return b.ec2Invoker }

// SetDynamoDBInvoker sets the DynamoDB invoker used for dispatching DynamoDB
// item operations from bus events (e.g. AppSync GraphQL resolvers).
func (b *EventBus) SetDynamoDBInvoker(invoker DynamoDBInvoker) { b.dynamoDBInvoker = invoker }

// DynamoDBInvoker returns the configured DynamoDB invoker.
func (b *EventBus) DynamoDBInvoker() DynamoDBInvoker { return b.dynamoDBInvoker }

// SetNeptuneGraphInvoker sets the NeptuneGraph invoker used for dispatching
// graph queries from bus consumers (e.g. AppSync GraphQL resolvers).
func (b *EventBus) SetNeptuneGraphInvoker(invoker NeptuneGraphInvoker) {
	b.neptuneGraphInvoker = invoker
}

// NeptuneGraphInvoker returns the configured NeptuneGraph invoker.
func (b *EventBus) NeptuneGraphInvoker() NeptuneGraphInvoker { return b.neptuneGraphInvoker }

// SetKMSInvoker sets the KMS invoker used for encryption key operations.
func (b *EventBus) SetKMSInvoker(invoker KMSInvoker) { b.kmsInvoker = invoker }

// KMSInvoker returns the configured KMS invoker.
func (b *EventBus) KMSInvoker() KMSInvoker { return b.kmsInvoker }

// RoleResolver returns the configured RoleResolver, or nil if none was set.
func (b *EventBus) RoleResolver() RoleResolver {
	return b.roleResolver
}

// Start recovers any pending outbox entries and launches async workers.
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

	for i := 0; i < AsyncWorkerCount; i++ {
		b.wg.Add(1)
		go b.directWorker()
	}

	if b.outbox != nil {
		b.wg.Add(1)
		go b.cleanupLoop()
	}

	b.started.Store(true)
	return nil
}

// Shutdown stops the bus, waits for all in-flight handlers to complete,
// and closes the outbox store.
func (b *EventBus) Shutdown(ctx context.Context) error {
	b.started.Store(false)
	b.stopOnce.Do(func() { close(b.stopCh) })

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
			if !b.acquireSemaphores(d.sub) {
				b.wg.Done()
				continue
			}
			b.dispatchHandler(d.ctx, d.sub, d.event)
			b.releaseSemaphores(d.sub)
			b.wg.Done()
		}
	}
}

func (b *EventBus) processOutboxEntry(entry *OutboxEntry) {
	ctx := context.Background()

	if entry.HandlerResults == nil {
		entry.HandlerResults = make(map[string]string)
	}

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
		b.persistEntry(ctx, entry)
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
			b.persistEntry(ctx, entry)
		} else {
			entry.Status = OutboxPending
			b.persistEntry(ctx, entry)
			select {
			case b.asyncCh <- entry:
			default:
			}
		}
	} else {
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
