// Package eventbridge provides AWS EventBridge (CloudWatch Events) service operations for vorpalstacks.
package eventbridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// EventsService provides AWS EventBridge operations.
type EventsService struct {
	storageManager  *storage.RegionStorageManager
	eventsStores    sync.Map // region → *eventsstore.EventsStore
	accountID       string
	bus             eventbus.ServiceBus
	targetSemaphore chan struct{}
	// delivery owns the retry lifecycle for target deliveries whose first
	// attempt failed; it is scheduling only, the attempt semantics live in
	// attemptDelivery/terminalDelivery.
	delivery *deliveryEngine
	// apiDestPacers enforces each API destination's invocation rate limit
	// (name → pacer); apiDestTokens caches OAuth access tokens per
	// connection ARN; apiDestHTTP is the shared API destination client
	// (5-second per-request timeout, TEST_MODE-relaxed TLS), built once at
	// construction so delivery paths only ever read it.
	apiDestPacersMu sync.Mutex
	apiDestPacers   map[string]*apiDestinationPacer
	apiDestTokensMu sync.Mutex
	apiDestTokens   map[string]*oauthToken
	apiDestHTTP     *http.Client
	replayWg        sync.WaitGroup
	replayCancels   sync.Map // replayName → context.CancelFunc
	schedCancel     context.CancelFunc
	schedWg         sync.WaitGroup
	// fireDedup is the scheduler's fired-boundary dedup state; per service
	// instance, never package-level, so two services never share fire
	// reservations.
	fireDedup scheduleFireDedup
}

const targetConcurrencyLimit = 100

// NewEventsService creates a new Events service instance.
// Optional cross-service dependencies should be injected via setter methods
// before registering handlers.
func NewEventsService(storageMgr *storage.RegionStorageManager, accountID string) *EventsService {
	svc := &EventsService{
		storageManager:  storageMgr,
		accountID:       accountID,
		targetSemaphore: make(chan struct{}, targetConcurrencyLimit),
		apiDestPacers:   map[string]*apiDestinationPacer{},
		apiDestTokens:   map[string]*oauthToken{},
		apiDestHTTP:     newAPIDestinationHTTPClient(),
	}
	svc.delivery = newDeliveryEngine(svc.attemptDelivery, svc.terminalDelivery)
	return svc
}

// Close cancels every delivery plane — the retry engine, the scheduler, the
// in-flight replays — and only then waits for the goroutines. Cancelling
// first bounds the wait: sleeping retry timers and replay workers observe
// the cancellation and exit instead of holding shutdown hostage (the former
// wait-then-cancel ordering left the replay cancels as dead letters).
func (s *EventsService) Close() {
	if s.delivery != nil {
		s.delivery.Close()
	}
	if s.schedCancel != nil {
		s.schedCancel()
	}
	s.replayCancels.Range(func(key, value any) bool {
		if cancel, ok := value.(context.CancelFunc); ok {
			cancel()
		}
		s.replayCancels.Delete(key)
		return true
	})
	s.schedWg.Wait()
	s.replayWg.Wait()
}

// SetEventsStore sets a pre-built events store for the given region,
// bypassing per-request store creation.
func (s *EventsService) SetEventsStore(region string, store *eventsstore.EventsStore) {
	if store != nil {
		s.eventsStores.Store(region, store)
	}
}

// SetEventBus injects the event bus and registers the EventBridge delivery handler.
func (s *EventsService) SetEventBus(bus eventbus.ServiceBus) error {
	s.bus = bus
	if _, err := eventbus.SubscribeTyped[*eventbus.EventBridgeDeliveryEvent](bus, s.handleBusDelivery, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("events: subscribe EventBridgeDeliveryEvent: %w", err)
	}
	if _, err := eventbus.SubscribeTyped[*eventbus.EventBridgePutEventsEvent](bus, s.handlePutEventsEvent, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("events: subscribe EventBridgePutEventsEvent: %w", err)
	}
	s.startScheduler()
	return nil
}

// resolveEventBusName canonicalises a name-or-ARN EventBusName value to the
// bus name (model member documentation: "The name or ARN of the event bus to
// receive the event"). A value that is not an event-bus ARN is returned
// unchanged and governs by its own name validity.
func resolveEventBusName(nameOrARN string) string {
	if name := svcarn.ExtractEventBusNameFromARN(nameOrARN); name != "" {
		return name
	}
	return nameOrARN
}

// GetStoreForRegion returns a cached or newly-created store for the given region.
// This is used by the admin handler and event bus subscribers that do not have
// a RequestContext.
func (s *EventsService) GetStoreForRegion(region string) (*eventsstore.EventsStore, error) {
	if cached, ok := s.eventsStores.Load(region); ok {
		if typed, ok := cached.(*eventsstore.EventsStore); ok {
			return typed, nil
		}
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}
	es := eventsstore.NewEventsStore(st, s.accountID, region)
	ensureDefaultEventBus(es)
	if actual, loaded := s.eventsStores.LoadOrStore(region, es); loaded {
		if typed, ok := actual.(*eventsstore.EventsStore); ok {
			return typed, nil
		}
	}
	// The first store of a region in this process is also the recovery
	// point for replays a previous process left in non-terminal states.
	s.recoverAbandonedReplays(es)
	return es, nil
}

// recoverAbandonedReplays converges replay records a previous process left
// in non-terminal states: no goroutine survives the restart to advance
// them. STARTING and RUNNING converge to FAILED, CANCELLING to the
// CANCELLED terminal its worker never wrote — through the same state-gated
// mutation the live worker's abort path uses, so a replay this process
// created after the store landed cannot be touched.
func (s *EventsService) recoverAbandonedReplays(es *eventsstore.EventsStore) {
	ctx := context.Background()
	nextToken := ""
	for {
		result, err := es.ListReplays(ctx, "", "", "", eventsstore.ListLimitMaximum, nextToken)
		if err != nil {
			logs.Warn("eventbridge: replay recovery failed to list replays", logs.Err(err))
			return
		}
		for _, replay := range result.Replays {
			switch replay.State {
			case eventsstore.ReplayStateStarting, eventsstore.ReplayStateRunning, eventsstore.ReplayStateCancelling:
				s.recordAbortedReplay(es, replay.Name)
			}
		}
		if result.NextToken == "" {
			return
		}
		nextToken = result.NextToken
	}
}

// ensureDefaultEventBus creates the default event bus when the store has
// none yet. Every AWS account carries a default bus in every region, so both
// store-acquisition planes (the request-scoped store and the region cache)
// seed it lazily; the idempotent create tolerates the already-present case.
func ensureDefaultEventBus(es *eventsstore.EventsStore) {
	defaultBus := &eventsstore.EventBus{Name: "default"}
	if err := es.CreateEventBus(context.Background(), defaultBus); err != nil {
		if err != eventsstore.ErrEventBusAlreadyExists {
			logs.Warn("eventbridge: failed to auto-create default event bus", logs.Err(err))
		}
	}
}

func (s *EventsService) handleBusDelivery(ctx context.Context, evt *eventbus.EventBridgeDeliveryEvent) eventbus.HandlerResult {
	if _, _, _, _, resource := svcarn.SplitARN(evt.TargetARN); strings.HasPrefix(resource, "event-bus/") {
		// The hop counter is carried on the delivery event because every
		// hop re-publishes a fresh chain root: without it, a bus pair that
		// targets each other amplifies one inbound event without bound.
		if evt.HopDepth >= maxCrossBusDepth {
			err := fmt.Errorf("cross-bus delivery depth %d reached the maximum of %d for event %s", evt.HopDepth, maxCrossBusDepth, evt.EventBridgeEventID)
			logs.Error("eventbridge: refusing cross-bus delivery",
				logs.String("targetArn", evt.TargetARN),
				logs.Err(err))
			return eventbus.HandlerResult{Error: err}
		}
		return s.handleEventBusDelivery(ctx, evt)
	}

	target := eventsstore.Target{
		ARN: evt.TargetARN,
		ID:  evt.TargetID,
	}
	if evt.DeadLetterConfigArn != "" {
		target.DeadLetterConfig = &eventsstore.DeadLetterConfig{Arn: evt.DeadLetterConfigArn}
	}
	// RetryPolicySet carries the policy's presence across the bus: an
	// explicit MaximumRetryAttempts of 0 is the documented "single
	// attempt" setting, so presence cannot be inferred from non-zero
	// values.
	if evt.RetryPolicySet {
		target.RetryPolicy = &eventsstore.RetryPolicy{
			MaximumRetryAttempts:     evt.MaximumRetryAttempts,
			MaximumEventAgeInSeconds: evt.MaximumEventAgeInSeconds,
		}
	}
	if evt.SqsMessageGroupId != "" {
		target.SqsParameters = &eventsstore.SqsParameters{MessageGroupId: evt.SqsMessageGroupId}
	}
	if evt.AppSyncGraphQLOperation != "" {
		target.AppSyncParameters = &eventsstore.AppSyncParameters{GraphQLOperation: evt.AppSyncGraphQLOperation}
	}
	if len(evt.TargetHttpParameters) > 0 {
		var httpParams eventsstore.HttpParameters
		if err := json.Unmarshal(evt.TargetHttpParameters, &httpParams); err == nil {
			target.HttpParameters = &httpParams
		}
	}

	// The event's identity comes from the delivery event, not the payload:
	// a target's Input/InputTransformer may have replaced the payload with
	// a constant, while the Kinesis partition-key default and the log-stream
	// name must still derive from the event's own ID. Publishers that
	// predate the identity field fall back to the bus identity.
	eventID := evt.EventBridgeEventID
	if eventID == "" {
		eventID = evt.EventID()
	}
	event := &eventsstore.Event{
		ID: eventID,
	}

	// The first attempt runs inline so the healthy path needs no engine
	// hop and the bus worker is held only for one bounded invocation; the
	// retry schedule belongs to the delivery engine, never to the worker.
	job := s.newDeliveryJob(evt.Region, event, target, evt.Input)
	job.attempts = 1
	job.ruleARN = evt.RuleARN
	job.traceHeader = evt.TraceHeader
	job.partitionKey = evt.KinesisPartitionKey
	if err := s.attemptDelivery(ctx, job); err == nil {
		return eventbus.HandlerResult{}
	} else if errors.Is(err, errPermanentDelivery) || job.exhausted() {
		if dlqErr := s.terminalDelivery(ctx, job, err); dlqErr != nil {
			// The retry budget is spent and the dead-letter copy failed:
			// the event is lost unless the transport-level retry re-drives
			// the delivery, so the failure must reach the outbox verdict.
			return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: dead-letter routing failed for target %s: %w", target.ARN, dlqErr)}
		}
		// The AWS contract terminates a delivery here — exhausted events
		// are dropped (dead-lettered when configured); reporting the
		// exhaustion as a handler failure would re-drive the delivery
		// through the transport retry budget and duplicate the terminal
		// handling on every pass.
		return eventbus.HandlerResult{}
	} else if err := s.delivery.schedule(job); err != nil {
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: delivery engine rejected the retry for target %s: %w", target.ARN, err)}
	}
	return eventbus.HandlerResult{}
}

func (s *EventsService) handleEventBusDelivery(ctx context.Context, evt *eventbus.EventBridgeDeliveryEvent) eventbus.HandlerResult {
	var envelope map[string]interface{}
	if err := json.Unmarshal(evt.Input, &envelope); err != nil {
		logs.Warn("eventbridge: failed to unmarshal event for rule matching", logs.String("error", err.Error()))
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: unmarshal bus delivery input: %w", err)}
	}

	eventBusName := "default"
	if name := svcarn.ExtractEventBusNameFromARN(evt.TargetARN); name != "" {
		eventBusName = name
	}

	event := eventFromEnvelope(envelope, eventBusName)

	// The platform receipt stamp rides the delivery event across the
	// bus: a rule-dispatched bus target keeps its original
	// ingestion-time, while a publisher that never went through ingress
	// (an SDK integration publishing events:deliver directly) is
	// stamped on arrival — the moment the platform received it.
	if evt.EventIngestionTime.IsZero() {
		event.IngestionTime = time.Now().UTC()
	} else {
		event.IngestionTime = evt.EventIngestionTime
	}

	es, err := s.GetStoreForRegion(evt.Region)
	if err != nil {
		logs.Warn("eventbridge: failed to get store for bus delivery",
			logs.String("region", evt.Region),
			logs.Err(err))
		return eventbus.HandlerResult{Error: err}
	}

	// The event-bus target is validated at PutTargets time, but the bus can
	// be deleted afterwards; delivering onto a vanished bus must reach the
	// publisher's retry and dead-letter policy, not vanish with a success
	// verdict.
	if _, err := es.GetEventBus(ctx, eventBusName); err != nil {
		logs.Warn("eventbridge: bus delivery targets a non-existent bus",
			logs.String("region", evt.Region),
			logs.String("eventBus", eventBusName))
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: bus delivery targets a non-existent bus %q", eventBusName)}
	}

	// The unified ingress archives the bus-delivered event like every other
	// ingress path ("EventBridge sends events that match the event pattern
	// to the archive" — the archives user guide states a pattern filter,
	// never an ingress-path one); the hop depth stays on the delivery so the
	// cross-bus cycle guard keeps counting.
	if err := s.deliverEvent(ctx, es, event, eventBusName, evt.Region, evt.HopDepth+1); err != nil {
		logs.Warn("eventbridge: failed to deliver event via rule matching", logs.String("error", err.Error()))
		// Synchronous publishers (the Scheduler engine publishes via
		// PublishSync) rely on the handler error to drive their retry and
		// dead-letter policy, so the failure must be propagated.
		return eventbus.HandlerResult{Error: err}
	}

	return eventbus.HandlerResult{}
}

func (s *EventsService) handlePutEventsEvent(ctx context.Context, evt *eventbus.EventBridgePutEventsEvent) eventbus.HandlerResult {
	region := evt.Region
	if region == "" {
		// An event without a region has no bus to land on; report the
		// drop so the publisher's retry and dead-letter policy sees it.
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents bus event carries no region")}
	}

	es, err := s.GetStoreForRegion(region)
	if err != nil {
		logs.Warn("eventbridge: failed to get storage for putEvents bus event",
			logs.String("region", region),
			logs.Err(err))
		return eventbus.HandlerResult{Error: err}
	}

	var inputMap map[string]interface{}
	if err := json.Unmarshal([]byte(evt.Input), &inputMap); err != nil {
		logs.Warn("eventbridge: failed to unmarshal putEvents input",
			logs.String("region", region),
			logs.Err(err))
		// PutEvents rejects a malformed payload at invocation time; the
		// bus path reports the same rejection instead of recording the
		// drop as a successful delivery.
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: unmarshal putEvents input: %w", err)}
	}

	// Source and DetailType can come from two places:
	// 1. The event struct itself (e.g. Scheduler's EventBridgeParameters)
	// 2. The Input JSON map (e.g. S3 notifications, direct PutEvents)
	// Struct fields take precedence over Input map values.
	source := evt.Source
	detailType := evt.DetailType
	if source == "" {
		source, _ = inputMap["Source"].(string)
	}
	if detailType == "" {
		detailType, _ = inputMap["DetailType"].(string)
	}
	if source == "" || detailType == "" {
		logs.Warn("eventbridge: putEvents input missing Source or DetailType",
			logs.String("region", region))
		// PutEvents requires Source and DetailType; a bus event carrying
		// neither is a failed delivery, not a successful no-op.
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents input missing Source or DetailType")}
	}

	var detail map[string]interface{}
	if detailRaw, ok := inputMap["Detail"]; ok {
		switch d := detailRaw.(type) {
		case map[string]interface{}:
			detail = d
		case string:
			// Detail is "a valid JSON object" (model member
			// documentation); the string form resolves through the same
			// parse the PutEvents API plane uses, and a malformed or
			// non-object value fails the delivery rather than silently
			// minting an empty detail.
			parsed, err := parseDetailObject(d)
			if err != nil {
				return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents Detail must be a valid JSON object: %w", err)}
			}
			detail = parsed
		default:
			return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents Detail must be a valid JSON object")}
		}
	} else {
		// Publishers carry the detail payload as the Input itself (the
		// scheduler's target input, the S3 notification detail); the
		// "Detail" wrapper key appears only when a publisher wraps it
		// explicitly. Falling back to the whole map keeps the delivered
		// event's detail populated for both forms.
		detail = inputMap
	}

	eventBusName := "default"
	if evt.EventBusName != "" {
		// EventBusName accepts the name or ARN form on both ingress
		// planes; rules and archives are keyed by the canonical name.
		eventBusName = resolveEventBusName(evt.EventBusName)
	}

	// The internal ingress plane reports a missing bus as a delivery
	// failure so the publisher's retry and dead-letter policy observes
	// the drop. The PutEvents API plane instead documents a silent drop
	// for the same condition (see putEventsCore); the asymmetry is the
	// internal-plane contract, not an oversight.
	if _, err := es.GetEventBus(ctx, eventBusName); err != nil {
		logs.Warn("eventbridge: putEvents bus event targets a non-existent bus",
			logs.String("region", region),
			logs.String("eventBus", eventBusName))
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents bus event targets a non-existent bus %q", eventBusName)}
	}

	// The same member bounds the PutEvents API plane enforces per entry:
	// DetailType 1-128 characters, Source 1-256 (API reference member
	// documentation).
	if !validateDetailType(detailType) {
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents DetailType must be between 1 and 128 characters")}
	}
	if !validateSource(source) {
		return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents Source must be between 1 and 256 characters")}
	}

	// A publisher-supplied Time member carries the PutEvents entry
	// semantics (parseEntryTimestamp); without one the arrival instant
	// applies.
	eventTime := time.Now().UTC()
	if timeVal, ok := inputMap["Time"]; ok {
		parsed, err := parseEntryTimestamp(timeVal)
		if err != nil {
			return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents %s", err.Error())}
		}
		eventTime = parsed
	}

	event := &eventsstore.Event{
		ID:         generateEventID(),
		Version:    "0",
		DetailType: detailType,
		Source:     source,
		Account:    s.accountID,
		// A publisher-supplied Time member carries the same semantics as
		// the PutEvents API entry loop (parseEntryTimestamp); the receipt
		// stamp below is always the arrival, never publisher-suppliable.
		Time:          eventTime,
		IngestionTime: time.Now().UTC(),
		Region:        region,
		Detail:        detail,
		EventBusName:  eventBusName,
	}

	if resources, ok := inputMap["Resources"].([]interface{}); ok {
		for _, r := range resources {
			rStr, ok := r.(string)
			if !ok {
				return eventbus.HandlerResult{Error: fmt.Errorf("eventbridge: putEvents Resources must be an array of strings")}
			}
			event.Resources = append(event.Resources, rStr)
		}
	}

	// The unified ingress: the event archives onto the bus's enabled
	// archives before rule matching, exactly like the PutEvents API plane.
	if err := s.deliverEvent(ctx, es, event, eventBusName, region, 0); err != nil {
		logs.Warn("eventbridge: failed to deliver putEvents bus event",
			logs.String("region", region),
			logs.String("error", err.Error()))
		// Synchronous publishers (the Scheduler engine publishes via
		// PublishSync) rely on the handler error to drive their retry and
		// dead-letter policy, so the failure must be propagated.
		return eventbus.HandlerResult{Error: err}
	}

	// The assigned event id rides the payload so synchronous publishers
	// can report it (the Step Functions putEvents integration answers the
	// entry's EventId from it); asynchronous publishers ignore it.
	return eventbus.HandlerResult{Payload: []byte(event.ID)}
}

func (s *EventsService) store(ctx *request.RequestContext) (*eventsstore.EventsStore, error) {
	return storecommon.GetOrCreateStoreE(&s.eventsStores, ctx.GetRegion(), func() (*eventsstore.EventsStore, error) {
		storage, err := ctx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		es := eventsstore.NewEventsStore(storage, s.accountID, ctx.GetRegion())
		ensureDefaultEventBus(es)
		// The request-scoped plane's first store of a region carries the
		// same recovery duty as GetStoreForRegion's.
		s.recoverAbandonedReplays(es)
		return es, nil
	})
}

// RegisterHandlers registers the Events service handlers with the dispatcher.
func (s *EventsService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("eventbridge", "CreateEventBus", s.CreateEventBus)
	d.RegisterHandlerForService("eventbridge", "DeleteEventBus", s.DeleteEventBus)
	d.RegisterHandlerForService("eventbridge", "DescribeEventBus", s.DescribeEventBus)
	d.RegisterHandlerForService("eventbridge", "ListEventBuses", s.ListEventBuses)
	d.RegisterHandlerForService("eventbridge", "UpdateEventBus", s.UpdateEventBus)

	d.RegisterHandlerForService("eventbridge", "PutRule", s.PutRule)
	d.RegisterHandlerForService("eventbridge", "DeleteRule", s.DeleteRule)
	d.RegisterHandlerForService("eventbridge", "DescribeRule", s.DescribeRule)
	d.RegisterHandlerForService("eventbridge", "ListRules", s.ListRules)
	d.RegisterHandlerForService("eventbridge", "EnableRule", s.EnableRule)
	d.RegisterHandlerForService("eventbridge", "DisableRule", s.DisableRule)

	d.RegisterHandlerForService("eventbridge", "PutTargets", s.PutTargets)
	d.RegisterHandlerForService("eventbridge", "RemoveTargets", s.RemoveTargets)
	d.RegisterHandlerForService("eventbridge", "ListTargetsByRule", s.ListTargetsByRule)
	d.RegisterHandlerForService("eventbridge", "ListRuleNamesByTarget", s.ListRuleNamesByTarget)

	d.RegisterHandlerForService("eventbridge", "PutEvents", s.PutEvents)

	d.RegisterHandlerForService("eventbridge", "TagResource", s.TagResource)
	d.RegisterHandlerForService("eventbridge", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("eventbridge", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("eventbridge", "CreateArchive", s.CreateArchive)
	d.RegisterHandlerForService("eventbridge", "DeleteArchive", s.DeleteArchive)
	d.RegisterHandlerForService("eventbridge", "DescribeArchive", s.DescribeArchive)
	d.RegisterHandlerForService("eventbridge", "UpdateArchive", s.UpdateArchive)
	d.RegisterHandlerForService("eventbridge", "ListArchives", s.ListArchives)

	d.RegisterHandlerForService("eventbridge", "StartReplay", s.StartReplay)
	d.RegisterHandlerForService("eventbridge", "DescribeReplay", s.DescribeReplay)
	d.RegisterHandlerForService("eventbridge", "ListReplays", s.ListReplays)
	d.RegisterHandlerForService("eventbridge", "CancelReplay", s.CancelReplay)

	d.RegisterHandlerForService("eventbridge", "CreateConnection", s.CreateConnection)
	d.RegisterHandlerForService("eventbridge", "DeleteConnection", s.DeleteConnection)
	d.RegisterHandlerForService("eventbridge", "DescribeConnection", s.DescribeConnection)
	d.RegisterHandlerForService("eventbridge", "UpdateConnection", s.UpdateConnection)
	d.RegisterHandlerForService("eventbridge", "DeauthorizeConnection", s.DeauthorizeConnection)
	d.RegisterHandlerForService("eventbridge", "ListConnections", s.ListConnections)

	d.RegisterHandlerForService("eventbridge", "CreateApiDestination", s.CreateApiDestination)
	d.RegisterHandlerForService("eventbridge", "DeleteApiDestination", s.DeleteApiDestination)
	d.RegisterHandlerForService("eventbridge", "DescribeApiDestination", s.DescribeApiDestination)
	d.RegisterHandlerForService("eventbridge", "UpdateApiDestination", s.UpdateApiDestination)
	d.RegisterHandlerForService("eventbridge", "ListApiDestinations", s.ListApiDestinations)

	d.RegisterHandlerForService("eventbridge", "TestEventPattern", s.TestEventPattern)

	d.RegisterHandlerForService("eventbridge", "PutPermission", s.PutPermission)
	d.RegisterHandlerForService("eventbridge", "RemovePermission", s.RemovePermission)
}
