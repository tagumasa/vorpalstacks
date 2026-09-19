package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// schedulerTickerInterval is the interval between schedule checks.
// Reduced to 1 second in test mode for faster test execution.
var schedulerTickerInterval = 1 * time.Minute

func init() {
	if os.Getenv("TEST_MODE") == "true" {
		schedulerTickerInterval = 1 * time.Second
	}
}

// Engine manages scheduled task execution for EventBridge Scheduler.
type Engine struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	bus            eventbus.ServiceBus

	// storeProvider returns the shared per-region SchedulerStore from the
	// owning service's cache. The engine keeps no store cache of its own:
	// one construction path and one closer (the service's StopEngine) serve
	// both planes, so no store goroutine is created or leaked per engine
	// start/stop cycle.
	storeProvider func(region string) (*schedulerstore.SchedulerStore, error)

	// retryStores holds per-region RetryStores for persisted retry records.
	// Records survive server restarts so the at-least-once delivery
	// guarantee is maintained.
	retryStores sync.Map // region → *schedulerstore.RetryStore

	// lastFired tracks the last execution per schedule (key:
	// region/groupName/name, built by lastFiredKey — region-scoped so
	// same-named schedules across regions keep independent deduplication
	// state) to prevent duplicate firing when the ticker polls multiple
	// times within the same execution window. The entry carries the
	// expression so a schedule whose expression changed (UpdateSchedule
	// also re-lifecycles a completed one-time schedule) starts a new
	// firing lifecycle instead of inheriting the previous expression's
	// suppression. Slots are released once the delivered boundary is
	// persisted on the schedule record (LastFiredAt), which — unlike this
	// map — survives restarts; the map only guards deliveries in flight
	// in this process and the window where persisting the boundary
	// failed.
	lastFired sync.Map // string → lastFiredEntry

	running   bool
	runningMu sync.RWMutex
	stopChan  chan struct{}
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewEngine creates a new scheduler engine with the given store dependencies.
func NewEngine(
	storageManager *storage.RegionStorageManager,
	accountID string,
) *Engine {
	return &Engine{
		storageManager: storageManager,
		accountID:      accountID,
		stopChan:       make(chan struct{}),
	}
}

// SetEventBus injects the event bus for publishing scheduler lifecycle events.
func (e *Engine) SetEventBus(bus eventbus.ServiceBus) {
	e.bus = bus
}

// SetStoreProvider injects the service's store-for-region function, making
// the service's cache the engine's single store source.
func (e *Engine) SetStoreProvider(provider func(region string) (*schedulerstore.SchedulerStore, error)) {
	e.storeProvider = provider
}

// Start starts the scheduler engine.
func (e *Engine) Start() error {
	e.runningMu.Lock()
	defer e.runningMu.Unlock()
	if e.running {
		return nil
	}
	e.running = true
	e.stopChan = make(chan struct{})
	e.ctx, e.cancel = context.WithCancel(context.Background())

	// Ensure the default schedule group exists in every regional store.
	// Without this, schedules whose omitted GroupName resolves to the
	// default group reference a phantom record that ListScheduleGroups
	// never shows and DeleteScheduleGroup cannot address.
	e.ensureDefaultGroups()

	e.wg.Add(1)
	go e.run()

	logs.Debug("Scheduler engine started")
	return nil
}

// Stop stops the scheduler engine.
func (e *Engine) Stop() error {
	e.runningMu.Lock()
	defer e.runningMu.Unlock()
	if !e.running {
		return nil
	}
	e.running = false
	if e.cancel != nil {
		e.cancel()
	}
	close(e.stopChan)
	e.wg.Wait()

	logs.Debug("Scheduler engine stopped")
	return nil
}

func (e *Engine) run() {
	defer e.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			resilience.RestartAfterPanic("scheduler engine", r, &e.wg, e.run)
		}
	}()

	ticker := time.NewTicker(schedulerTickerInterval)
	defer ticker.Stop()

	// Startup sweep: schedules and already-due retry records are processed
	// before the first tick, so records that became due while the engine was
	// down do not wait a full ticker interval.
	e.checkSchedules()
	e.checkRetries()

	for {
		select {
		case <-e.stopChan:
			return
		case <-ticker.C:
			e.checkSchedules()
			// Process pending retries from previous failed deliveries.
			e.checkRetries()
		}
	}
}

// ensureDefaultGroups creates the default schedule group in every
// regional store if it does not already exist. This is called once
// during Engine.Start, before run() is launched.
//
// It iterates active regions via storageManager; the per-region stores
// come from the service's shared cache (storeForRegion), so this pass
// also warms the cache that the sweeps later reuse.
func (e *Engine) ensureDefaultGroups() {
	if e.ctx == nil || e.storageManager == nil {
		return
	}

	regions := e.storageManager.GetActiveRegions()
	for _, region := range regions {
		store := e.storeForRegion(region)
		if store == nil {
			continue
		}
		if err := store.EnsureDefaultGroup(e.ctx); err != nil {
			logs.Warn("Failed to ensure default schedule group",
				logs.String("region", store.GetRegion()),
				logs.Err(err))
		}
	}
}

func (e *Engine) executeSchedule(ctx context.Context, schedule *schedulerstore.Schedule) error {
	logs.Debug("Executing schedule",
		logs.String("name", schedule.Name),
		logs.String("group", schedule.GroupName))

	if schedule.Target == nil {
		logs.Debug("Schedule has no target", logs.String("schedule", schedule.Name))
		return nil
	}

	target := schedule.Target
	targetArn := target.Arn
	region := schedule.Region
	if region == "" {
		region = defaults.DefaultRegion
	}

	if e.bus != nil {
		input := scheduleInput(target, schedule.Name)
		schedEvt := &eventbus.ScheduleFiredEvent{
			ScheduleName:          schedule.Name,
			ScheduleArn:           schedule.ARN,
			GroupName:             schedule.GroupName,
			TargetArn:             targetArn,
			Input:                 input,
			ActionAfterCompletion: string(schedule.ActionAfterCompletion),
			ScheduleExpression:    schedule.ScheduleExpression,
		}
		// Serialise the full target so the bus handler has access to all
		// sub-parameters (SqsParameters, KinesisParameters, etc.) without
		// needing to re-fetch the schedule from the store.
		if payloadBytes, err := json.Marshal(target); err == nil {
			schedEvt.TargetPayload = string(payloadBytes)
		}
		schedEvt.Region = region
		if err := e.bus.Publish(context.Background(), schedEvt); err != nil {
			logs.Warn("Failed to publish schedule fired event to event bus",
				logs.String("schedule", schedule.Name),
				logs.String("target", targetArn),
				logs.Err(err))
			return err
		}
	} else {
		// Direct delivery path: attempt delivery with retry. The retry
		// chain covers immediate retries, persisted background retries,
		// and DLQ routing on permanent failure.
		e.deliverWithRetry(ctx, schedule, target)
	}

	// NOTE: post-execution actions are handled at delivery lifecycle
	// completion (success or retry exhaustion) by maybeActionAfterCompletion,
	// not here: ActionAfterCompletion=DELETE deletes the schedule (AWS
	// deletes it "shortly after its last target invocation", i.e. after the
	// retry policy terminates — deleting here would orphan retry records and
	// break at-least-once delivery semantics), and a one-time schedule with
	// no action is marked completed so it never fires again.

	return nil
}

func (e *Engine) getStoreForSchedule(schedule *schedulerstore.Schedule) *schedulerstore.SchedulerStore {
	region := schedule.Region
	if region == "" {
		region = defaults.DefaultRegion
	}
	return e.storeForRegion(region)
}

// storeForRegion returns the shared per-region SchedulerStore from the
// service's cache. A nil provider (an engine constructed without one, e.g.
// in unit tests) yields nil and callers skip the region.
func (e *Engine) storeForRegion(region string) *schedulerstore.SchedulerStore {
	if e.storeProvider == nil {
		return nil
	}
	store, err := e.storeProvider(region)
	if err != nil {
		logs.Debug("Failed to get store for region",
			logs.String("region", region),
			logs.String("error", err.Error()))
		return nil
	}
	return store
}

// maybeAutoDelete deletes the schedule when ActionAfterCompletion is DELETE.
// It is invoked at delivery lifecycle completion points (success or retry
// exhaustion) so that the schedule remains alive during the entire retry
// lifecycle, matching AWS EventBridge Scheduler semantics.
//
// AWS documentation: "the schedule is deleted shortly after its last target
// invocation" — meaning the final successful attempt or the last unsuccessful
// retry once the retry policy is exhausted.
//
// Idempotent: calling twice (e.g. success path racing with retry exhaustion)
// is safe — the second DeleteSchedule call returns NotFound and the error is
// logged at Debug level.
func (e *Engine) maybeAutoDelete(ctx context.Context, schedule *schedulerstore.Schedule) {
	if schedule == nil {
		return
	}
	if schedule.ActionAfterCompletion != schedulerstore.ActionAfterCompletionDelete {
		return
	}
	store := e.getStoreForSchedule(schedule)
	if store == nil {
		logs.Warn("No store available for auto-delete after completion",
			logs.String("schedule", schedule.Name),
			logs.String("group", schedule.GroupName),
			logs.String("region", schedule.Region))
		return
	}
	if err := store.DeleteSchedule(ctx, schedule.GroupName, schedule.Name); err != nil {
		logs.Debug("Failed to auto-delete schedule after completion",
			logs.String("schedule", schedule.Name),
			logs.String("group", schedule.GroupName),
			logs.String("error", err.Error()))
	}
}

// maybeActionAfterCompletion applies the post-execution action at the
// delivery lifecycle completion points (success or retry exhaustion):
// ActionAfterCompletion=DELETE removes the schedule via maybeAutoDelete,
// and a one-time schedule with no action is marked completed so it never
// fires again.
func (e *Engine) maybeActionAfterCompletion(ctx context.Context, schedule *schedulerstore.Schedule) {
	if schedule == nil {
		return
	}
	if schedule.ActionAfterCompletion == schedulerstore.ActionAfterCompletionDelete {
		e.maybeAutoDelete(ctx, schedule)
		return
	}
	e.maybeCompleteOnetime(ctx, schedule)
}

// maybeCompleteOnetime ends the firing lifecycle of a one-time at()
// schedule. The AWS ScheduleState enum has no COMPLETED value, so the
// wire state is preserved and completion is recorded as an internal
// persisted marker — without it the in-memory dedup map would forget the
// fire on the next server restart and run the one-time schedule again.
func (e *Engine) maybeCompleteOnetime(ctx context.Context, schedule *schedulerstore.Schedule) {
	if !scheduleexpr.IsAtExpression(schedule.ScheduleExpression) {
		return
	}
	store := e.getStoreForSchedule(schedule)
	if store == nil {
		logs.Warn("No store available to complete one-time schedule",
			logs.String("schedule", schedule.Name),
			logs.String("group", schedule.GroupName),
			logs.String("region", schedule.Region))
		return
	}
	if err := store.CompleteSchedule(ctx, schedule.GroupName, schedule.Name); err != nil {
		logs.Debug("Failed to mark one-time schedule completed",
			logs.String("schedule", schedule.Name),
			logs.String("group", schedule.GroupName),
			logs.String("error", err.Error()))
	}
}

// getRetryStore returns the RetryStore for the given region, creating it
// lazily on first access.
func (e *Engine) getRetryStore(region string) (*schedulerstore.RetryStore, error) {
	if region == "" {
		region = defaults.DefaultRegion
	}
	if cached, ok := e.retryStores.Load(region); ok {
		return cached.(*schedulerstore.RetryStore), nil
	}
	// The sweep entry guards a nil storage manager, but this accessor is
	// its own contract: a construction without one reports it instead of
	// dereferencing nil.
	if e.storageManager == nil {
		return nil, fmt.Errorf("scheduler engine has no storage manager")
	}
	storage, err := e.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("no storage available for region %q: %w", region, err)
	}
	rs := schedulerstore.NewRetryStore(storage, region)
	actual, _ := e.retryStores.LoadOrStore(region, rs)
	return actual.(*schedulerstore.RetryStore), nil
}
