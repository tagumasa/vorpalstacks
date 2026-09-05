package eventbus

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/common/invokers"
	waf "vorpalstacks/internal/common/invokers/waf"
	"vorpalstacks/internal/core/logs"
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
	// DeliveredRetention is how long successfully delivered outbox entries are kept.
	DeliveredRetention = 1 * time.Hour
	// FailedRetention is how long failed outbox entries are kept before purging.
	FailedRetention = 24 * time.Hour
)

// CleanupInterval is the period between outbox retention cleanups.
// Overridden to 30s in TEST_MODE for faster stale entry removal.
var CleanupInterval = 10 * time.Minute

// PendingRequeueInterval is how often the requeue loop scans the outbox for
// entries still pending and re-enqueues them, so delivery does not depend
// on a restart.
// Overridden to 2s in TEST_MODE.
var PendingRequeueInterval = 30 * time.Second

func init() {
	if os.Getenv("TEST_MODE") == "true" {
		CleanupInterval = 30 * time.Second
		PendingRequeueInterval = 2 * time.Second
	}
}

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
// executions across all subscriptions. The limit must be positive: an
// unbuffered (zero-capacity) semaphore would block every dispatch until
// shutdown.
func WithGlobalConcurrency(n int) BusOption {
	return func(b *EventBus) {
		if n <= 0 {
			panic("eventbus: WithGlobalConcurrency requires a positive limit")
		}
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

// Bus defines the contract for the internal service event bus: synchronous
// and asynchronous publish, subscription management, and dispatch-time
// cross-service authorisation. The typed invoker registry surface is a
// separate role, declared in invokers.Registry and composed into
// ServiceBus for consumers that need both.
type Bus interface {
	Publish(ctx context.Context, event Event) error
	PublishSync(ctx context.Context, event Event) (HandlerResult, error)
	Subscribe(handler func(ctx context.Context, event Event) HandlerResult, opts ...SubscribeOption) (string, error)
	Unsubscribe(subscriptionID string) error
	EvaluateTargetPolicy(ctx context.Context, targetARN, serviceType, principal, action, resource string) (bool, error)
	RoleResolver() RoleResolver
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// ServiceBus composes the delivery surface with the typed invoker
// registry. Services are handed this combined view so a single injected
// object serves both roles; consumers needing only one side accept the
// narrower interface.
type ServiceBus interface {
	Bus
	invokers.Registry
}

// EventBus is the central implementation of the Bus interface, managing
// subscriptions, outbox persistence, async workers, and invoker dispatch.
type EventBus struct {
	mu                      sync.RWMutex
	subscriptions           map[string][]*subscriptionEntry
	outbox                  OutboxStore
	registry                *EventRegistry
	roleResolver            RoleResolver
	policyEval              BusPolicyEvaluator
	policyFuncs             map[string]ResourcePolicyFunc
	policyFuncsMu           sync.RWMutex
	globalSem               chan struct{}
	maxRetries              int32
	maxEventDepth           int
	logger                  logs.Logger
	wg                      sync.WaitGroup
	started                 atomic.Bool
	startMu                 sync.Mutex
	startDone               bool
	stopOnce                sync.Once
	stopCh                  chan struct{}
	invokersMu              sync.RWMutex
	lambdaInvoker           invokers.LambdaInvoker
	sqsInvoker              invokers.SQSInvoker
	snsInvoker              invokers.SNSInvoker
	kinesisInvoker          invokers.KinesisInvoker
	eventsInvoker           invokers.EventsInvoker
	ec2Invoker              invokers.EC2Invoker
	dynamoDBInvoker         invokers.DynamoDBInvoker
	dynamoDBStreamsInvoker  invokers.DynamoDBStreamsInvoker
	neptuneGraphInvoker     invokers.NeptuneGraphInvoker
	kmsInvoker              invokers.KMSInvoker
	s3Invoker               invokers.S3Invoker
	wafInvoker              invokers.WAFInvoker
	webACLInspector         waf.WebACLInspector
	cloudWatchMetricInvoker invokers.CloudWatchMetricInvoker
	cloudWatchAlarmInvoker  invokers.CloudWatchAlarmInvoker
	timestreamInvoker       invokers.TimestreamInvoker
	cloudTrailInvoker       invokers.CloudTrailInvoker
	logsInvoker             invokers.LogsInvoker
	rdsDataInvoker          invokers.RDSDataInvoker
	cognitoTokenValidator   invokers.CognitoTokenValidator
	subnetUsageCheckers     []invokers.SubnetUsageChecker
	securityGroupCheckers   []invokers.SecurityGroupUsageChecker
	nextSubID               atomic.Int64
	asyncCh                 chan *OutboxEntry
	directCh                chan *directDispatch
	// requeueCursor is where the last requeuePending walk stopped. It is
	// written and read only by the requeuePendingLoop goroutine, so it
	// needs no lock; it lets the next tick resume behind the entries
	// already queued instead of re-filling from the head of the backlog.
	requeueCursor string
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
		asyncCh:       make(chan *OutboxEntry, 1024),
		directCh:      make(chan *directDispatch, 1024),
		policyFuncs:   make(map[string]ResourcePolicyFunc),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
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
	// A random suffix prevents collisions between events of the same type
	// published within the same nanosecond: outbox entries are keyed by
	// event ID, so a collision would silently overwrite one of the events.
	var rnd [8]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return fmt.Sprintf("%s-%d", eventType, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s-%d-%s", eventType, time.Now().UnixNano(), hex.EncodeToString(rnd[:]))
}

// logFields converts alternating key/value pairs into structured log
// fields. A trailing key without a value is dropped.
func logFields(keyvals ...interface{}) []logs.Field {
	fields := make([]logs.Field, 0, len(keyvals)/2)
	for i := 0; i+1 < len(keyvals); i += 2 {
		fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
	}
	return fields
}

func (b *EventBus) logInfo(msg string, keyvals ...interface{}) {
	if b.logger != nil {
		b.logger.Info(msg, logFields(keyvals...)...)
	}
}

func (b *EventBus) logWarn(msg string, keyvals ...interface{}) {
	if b.logger != nil {
		b.logger.Warn(msg, logFields(keyvals...)...)
	}
}

var _ Bus = (*EventBus)(nil)
