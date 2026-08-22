package scheduler

import (
	"context"
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
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
	bus            eventbus.Bus
	stores         sync.Map // region → *schedulerstore.SchedulerStore

	// retryStores holds per-region RetryStores for persisted retry records.
	// Records survive server restarts so the at-least-once delivery
	// guarantee is maintained.
	retryStores sync.Map // region → *schedulerstore.RetryStore

	// lastFired tracks the last execution per schedule (key: groupName/name)
	// to prevent duplicate firing when the ticker polls multiple times within
	// the same execution window. The entry carries the expression so a
	// schedule whose expression changed (UpdateSchedule also re-lifecycles
	// a completed one-time schedule) starts a new firing lifecycle instead
	// of inheriting the previous expression's suppression. Slots are
	// released once the delivered boundary is persisted on the schedule
	// record (LastFiredAt), which — unlike this map — survives restarts;
	// the map only guards deliveries in flight in this process and the
	// window where persisting the boundary failed.
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
func (e *Engine) SetEventBus(bus eventbus.Bus) {
	e.bus = bus
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
	// Without this, schedules created with GroupName="" use a phantom
	// "default" group that cannot be deleted via DeleteScheduleGroup.
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
	defer func() { resilience.RecoverAndRestart("scheduler engine", &e.wg, e.run) }()

	ticker := time.NewTicker(schedulerTickerInterval)
	defer ticker.Stop()

	e.checkSchedules()

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
// It iterates active regions via storageManager (not e.stores) because
// e.stores is empty at this point — it is only populated inside
// checkSchedules() which runs in the run() goroutine launched after
// this method returns. The stores created here are cached in e.stores
// via LoadOrStore so checkSchedules() reuses them.
func (e *Engine) ensureDefaultGroups() {
	if e.ctx == nil || e.storageManager == nil {
		return
	}

	regions := e.storageManager.GetActiveRegions()
	for _, region := range regions {
		storage, err := e.storageManager.GetStorage(region)
		if err != nil {
			logs.Debug("Failed to get storage for region",
				logs.String("region", region),
				logs.String("error", err.Error()))
			continue
		}
		store := schedulerstore.NewSchedulerStore(storage, e.accountID, region)
		if actual, loaded := e.stores.LoadOrStore(region, store); loaded {
			store = actual.(*schedulerstore.SchedulerStore)
		}
		if err := store.EnsureDefaultGroup(e.ctx); err != nil {
			logs.Warn("Failed to ensure default schedule group",
				logs.String("region", store.GetRegion()),
				logs.Err(err))
		}
	}
}

// lastFiredEntry records the most recent fire of a schedule: when it
// fired and under which expression. The expression scopes the dedup —
// only boundaries of the same expression are suppressed.
type lastFiredEntry struct {
	firedAt time.Time
	expr    string
}

// lastFiredKey builds the deduplication key used by the lastFired map.
// Region is part of the key so that schedules sharing the same
// group/name across regions do not interfere with each other's
// deduplication state.
func lastFiredKey(region, groupName, name string) string {
	return region + "/" + groupName + "/" + name
}

func (e *Engine) checkSchedules() {
	if e.storageManager == nil {
		return
	}

	regions := e.storageManager.GetActiveRegions()

	// Collect active schedule keys across ALL regions so the lazy cleanup
	// pass at the end does not destroy entries belonging to a region that
	// has not yet been visited, or to a region whose schedules happen to
	// share a legacy (region-less) key with the current region.
	allActiveKeys := make(map[string]bool)

	for _, region := range regions {
		storage, err := e.storageManager.GetStorage(region)
		if err != nil {
			logs.Debug("Failed to get storage for region", logs.String("region", region), logs.String("error", err.Error()))
			continue
		}
		store := schedulerstore.NewSchedulerStore(storage, e.accountID, region)
		if actual, loaded := e.stores.LoadOrStore(region, store); loaded {
			store = actual.(*schedulerstore.SchedulerStore)
		}
		schedules, err := store.GetAllEnabledSchedules(e.ctx)
		if err != nil {
			logs.Debug("Failed to get enabled schedules", logs.String("region", region), logs.String("error", err.Error()))
			continue
		}

		now := time.Now().UTC()

		for _, schedule := range schedules {
			schedule.Region = region
			allActiveKeys[lastFiredKey(region, schedule.GroupName, schedule.Name)] = true
			boundary, due := e.dueBoundary(schedule, now)
			if due {
				// Provisionally reserve the dedup slot so a concurrent tick
				// does not double-fire while the goroutine is still in flight.
				// If executeSchedule fails or panics the goroutine releases
				// the reservation so the next tick can retry the schedule.
				dedupKey := lastFiredKey(region, schedule.GroupName, schedule.Name)
				e.lastFired.Store(dedupKey, lastFiredEntry{firedAt: now, expr: schedule.ScheduleExpression})
				e.wg.Add(1)
				go func(sch *schedulerstore.Schedule, key string, st *schedulerstore.SchedulerStore, fired time.Time) {
					defer e.wg.Done()
					defer func() {
						if r := recover(); r != nil {
							logs.Error("scheduler: panic executing schedule", logs.String("name", sch.Name), logs.Any("panic", r))
							// Release the dedup slot so the next tick can retry the schedule.
							e.lastFired.Delete(key)
						}
					}()
					select {
					case <-e.ctx.Done():
						return
					default:
						if err := e.executeSchedule(e.ctx, sch); err != nil {
							// Release the dedup slot so the next tick can retry the schedule.
							e.lastFired.Delete(key)
							return
						}
						// Persist the delivered boundary on the schedule record
						// so a restart cannot deliver it twice. The in-memory
						// slot then becomes redundant (the persisted marker
						// suppresses the boundary) and is released for the
						// next boundary.
						if err := st.TouchScheduleLastFired(e.ctx, sch.GroupName, sch.Name, fired); err != nil && !errors.Is(err, schedulerstore.ErrScheduleNotFound) {
							// A schedule deleted by its own completion has
							// nothing left to guard; any other failure keeps
							// the in-flight slot, which is then the only
							// duplicate guard until a restart.
							logs.Warn("Failed to persist the delivered boundary of schedule", logs.String("name", sch.Name), logs.Err(err))
							return
						}
						e.lastFired.Delete(key)
					}
				}(schedule, dedupKey, store, boundary)
			}
		}
	}

	// Remove lastFired entries for schedules that no longer exist so the
	// dedup map does not grow unbounded across schedule create/delete
	// cycles. Runs once after the full region sweep so that other regions'
	// entries are not destroyed while their sweep is still iterating.
	e.lastFired.Range(func(key, _ interface{}) bool {
		if k, ok := key.(string); ok {
			if !allActiveKeys[k] {
				e.lastFired.Delete(k)
			}
		}
		return true
	})
}

// resolveScheduleLocation returns the time.Location in which the schedule
// expression should be evaluated. When ScheduleExpressionTimezone is empty or
// invalid, UTC is used (matching the AWS default of UTC).
func resolveScheduleLocation(schedule *schedulerstore.Schedule) *time.Location {
	if schedule.ScheduleExpressionTimezone != "" {
		if loc, err := time.LoadLocation(schedule.ScheduleExpressionTimezone); err == nil {
			return loc
		}
		logs.Debug("Invalid schedule timezone, falling back to UTC",
			logs.String("schedule", schedule.Name),
			logs.String("timezone", schedule.ScheduleExpressionTimezone))
	}
	return time.UTC
}

// dueBoundary returns the schedule boundary this evaluation should deliver,
// if any. It is the decision core of the sweep; shouldExecute wraps it for
// callers that only need the boolean.
func (e *Engine) dueBoundary(schedule *schedulerstore.Schedule, now time.Time) (time.Time, bool) {
	// Convert "now" to the schedule's evaluation timezone so that rate/cron/at
	// expressions are evaluated in the timezone the user configured.
	loc := resolveScheduleLocation(schedule)
	nowLocal := now.In(loc)

	// AWS: "When you configure a one-time schedule, EventBridge Scheduler
	// ignores the StartDate and EndDate you specify for the schedule."
	// Only rate()/cron() honour StartDate/EndDate.
	isAtExpression := strings.HasPrefix(schedule.ScheduleExpression, "at(")
	if !isAtExpression {
		if schedule.StartDate != nil && nowLocal.Before(schedule.StartDate.In(loc)) {
			return time.Time{}, false
		}
		if schedule.EndDate != nil && nowLocal.After(schedule.EndDate.In(loc)) {
			return time.Time{}, false
		}
	}

	boundary, elapsed := scheduleexpr.ElapsedExecutionTime(schedule.ScheduleExpression, nowLocal, schedule.CreationDate, schedule.StartDate)
	if !elapsed {
		return time.Time{}, false
	}

	// The persisted delivered-boundary marker is the source of truth across
	// restarts (the in-memory map dies with the process): an occurrence
	// already recorded as delivered must not fire again.
	if schedule.LastFiredAt != nil && !schedule.LastFiredAt.Before(boundary) {
		return time.Time{}, false
	}

	// Prevent duplicate firing: skip if already executed for this interval.
	// Key includes region so multi-region schedules do not share state.
	// dueBoundary is a pure predicate; the caller reserves the dedup slot
	// after this returns the boundary.
	key := lastFiredKey(schedule.Region, schedule.GroupName, schedule.Name)
	if last, ok := e.lastFired.Load(key); ok {
		// Only a fire under the same expression suppresses: an expression
		// change starts a new firing lifecycle (a completed one-time
		// schedule updated to a new past at() must fire again).
		if entry, ok := last.(lastFiredEntry); ok && entry.expr == schedule.ScheduleExpression && !entry.firedAt.Before(boundary) {
			return time.Time{}, false
		}
	}

	if schedule.FlexibleTimeWindow != nil && schedule.FlexibleTimeWindow.Mode == schedulerstore.FlexibleTimeWindowModeFlexible {
		maxWindow := 1
		if schedule.FlexibleTimeWindow.MaximumWindowInMinutes != nil {
			maxWindow = *schedule.FlexibleTimeWindow.MaximumWindowInMinutes
		}
		// AWS flexible time window starts at the scheduled time and extends
		// forward for MaximumWindowInMinutes. Execution never occurs before
		// the scheduled time (the resolved boundary already lies at or
		// before now); an occurrence whose window has closed is skipped,
		// not resurrected.
		windowEnd := boundary.Add(time.Duration(maxWindow) * time.Minute)
		if !nowLocal.Before(windowEnd) {
			return time.Time{}, false
		}
	}

	// The boundary lies at or before now: fire on this evaluation. A late
	// evaluation (ticker gap longer than the boundary interval) still fires
	// the pending boundary once instead of silently skipping it, and the
	// dedup checks above cap each boundary at one fire.
	return boundary, true
}

// shouldExecute reports whether the schedule should fire at this evaluation.
func (e *Engine) shouldExecute(schedule *schedulerstore.Schedule, now time.Time) bool {
	_, ok := e.dueBoundary(schedule, now)
	return ok
}

func scheduleInput(target *schedulerstore.Target, scheduleName string) string {
	if target.Input != "" {
		return target.Input
	}
	msgPayload := map[string]interface{}{
		"schedule":  scheduleName,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	if msgBytes, err := json.Marshal(msgPayload); err == nil {
		return string(msgBytes)
	}
	return "{}"
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
	if cached, ok := e.stores.Load(region); ok {
		return cached.(*schedulerstore.SchedulerStore)
	}
	storage, err := e.storageManager.GetStorage(region)
	if err != nil {
		return nil
	}
	store := schedulerstore.NewSchedulerStore(storage, e.accountID, region)
	if actual, loaded := e.stores.LoadOrStore(region, store); loaded {
		return actual.(*schedulerstore.SchedulerStore)
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
	if !strings.HasPrefix(schedule.ScheduleExpression, "at(") {
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
	storage, err := e.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("no storage available for region %q: %w", region, err)
	}
	rs := schedulerstore.NewRetryStore(storage, region)
	actual, _ := e.retryStores.LoadOrStore(region, rs)
	return actual.(*schedulerstore.RetryStore), nil
}

// deliverToTarget dispatches a schedule delivery to the appropriate target
// type and returns an error on failure. This is the single entry point for
// all target deliveries, used by both the direct path and the bus path.
func (e *Engine) deliverToTarget(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	_, service, _, _, _ := svcarn.SplitARN(target.Arn)
	switch service {
	case "lambda":
		return e.invokeLambda(ctx, schedule, target)
	case "sqs":
		return e.sendToSQS(ctx, schedule, target)
	case "sns":
		return e.publishToSNS(ctx, schedule, target)
	case "kinesis":
		return e.sendToKinesis(ctx, schedule, target)
	case "states":
		return e.startStepFunctionExecution(ctx, schedule, target)
	case "events":
		return e.sendToEventBridge(ctx, schedule, target)
	case "logs":
		return e.sendToCloudWatchLogs(ctx, schedule, target)
	case "ecs":
		// ECS is an AWS templated target. The ECS service is not yet
		// available on this platform, so delivery fails and the schedule
		// engine's retry/DLQ path takes over. EcsParameters validation
		// is fully implemented, so the schedule itself is valid.
		logs.Error("ECS target delivery failed: ECS service is not available",
			logs.String("targetArn", target.Arn))
		return fmt.Errorf("ecs delivery target %s is not available in this deployment", target.Arn)
	case "firehose":
		// Firehose is an AWS templated target with no sub-parameters.
		// The Firehose service is not yet available on this platform,
		// so delivery fails and the retry/DLQ path takes over.
		logs.Error("Firehose target delivery failed: Firehose service is not available",
			logs.String("targetArn", target.Arn))
		return fmt.Errorf("firehose delivery target %s is not available in this deployment", target.Arn)
	default:
		return fmt.Errorf("unsupported target type: %s", target.Arn)
	}
}

// retryDefaults returns the effective MaximumRetryAttempts and
// MaximumEventAgeInSeconds for a target, applying AWS defaults when the
// RetryPolicy is nil or individual fields are unset.
// Defaults: 185 retries, 86400 seconds (24 hours).
func retryDefaults(target *schedulerstore.Target) (maxRetries int, maxAgeSeconds int) {
	maxRetries = 185
	maxAgeSeconds = 86400
	if target.RetryPolicy != nil {
		if target.RetryPolicy.MaximumRetryAttempts != nil && *target.RetryPolicy.MaximumRetryAttempts >= 0 {
			maxRetries = *target.RetryPolicy.MaximumRetryAttempts
		}
		if target.RetryPolicy.MaximumEventAgeInSeconds != nil && *target.RetryPolicy.MaximumEventAgeInSeconds >= 60 {
			maxAgeSeconds = *target.RetryPolicy.MaximumEventAgeInSeconds
		}
	}
	return
}

// computeRetryBackoff calculates the delay before the next retry attempt
// using exponential backoff with jitter. The base interval doubles with each
// attempt, capped at 1 hour. Jitter is up to 50% of the base interval.
func computeRetryBackoff(attemptCount int) time.Duration {
	if attemptCount < 1 {
		attemptCount = 1
	}
	// Cap the shift to prevent int64 overflow. 1<<63 overflows to a
	// negative value, causing crand.Int to panic. Since the cap below
	// already limits base to 1 hour, attempts beyond ~10 are equivalent.
	shift := attemptCount - 1
	if shift > 12 {
		shift = 12 // 1<<12 seconds = 4096s ≈ 68min, capped to 1h below
	}
	base := time.Duration(1<<uint(shift)) * time.Second
	if base > time.Hour {
		base = time.Hour
	}
	// Add up to 50% jitter using crypto/rand for unpredictability.
	maxJitter := int64(base / 2)
	if maxJitter > 0 {
		n, err := crand.Int(crand.Reader, big.NewInt(maxJitter))
		if err == nil {
			return base + time.Duration(n.Int64())
		}
	}
	return base
}

// deliverWithRetry attempts delivery with an immediate retry, then persists a
// RetryRecord for background retries if the second attempt also fails.
// The RetryRecord survives server restarts so that at-least-once delivery
// is maintained.
//
// Post-execution actions are handled at lifecycle completion via
// maybeActionAfterCompletion:
//   - success on attempt 1 or 2 → applied here
//   - maxRetries=0 DLQ route    → applied here
//   - retry exhausted / event age exceeded → applied in processRetryRecord
//
// For ActionAfterCompletion=DELETE see the AWS blog: "the schedule is
// deleted shortly after its last target invocation" — NOT at fire time;
// for a one-time schedule with no action the completion marker ends its
// firing lifecycle.
func (e *Engine) deliverWithRetry(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	maxRetries, _ := retryDefaults(target)

	// First attempt (immediate).
	err := e.deliverToTarget(ctx, schedule, target)
	if err == nil {
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	// AWS: MaximumRetryAttempts=0 means no retries — only the initial
	// attempt. Route to DLQ immediately on failure.
	if maxRetries == 0 {
		logs.Warn("Scheduler delivery failed, no retries configured",
			logs.String("schedule", schedule.Name),
			logs.String("target", target.Arn),
			logs.Err(err))
		e.routeToDLQ(ctx, schedule, target, "delivery failed (MaximumRetryAttempts=0)")
		// Retry lifecycle is now complete (no retries configured).
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	logs.Warn("Scheduler delivery failed, retrying",
		logs.String("schedule", schedule.Name),
		logs.String("target", target.Arn),
		logs.Err(err))

	// Immediate retry (attempt 2).
	err = e.deliverToTarget(ctx, schedule, target)
	if err == nil {
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	// Both immediate attempts failed — persist for background retry.
	region := schedule.Region
	if region == "" {
		region = defaults.DefaultRegion
	}
	now := time.Now()

	targetJSON, mErr := json.Marshal(target)
	if mErr != nil {
		logs.Error("Failed to serialise target for retry record",
			logs.String("schedule", schedule.Name),
			logs.Err(mErr))
		return
	}

	record := &schedulerstore.RetryRecord{
		ID:                    fmt.Sprintf("retry-%s-%d", uuidString(), now.UnixNano()),
		ScheduleName:          schedule.Name,
		GroupName:             schedule.GroupName,
		Region:                region,
		Target:                string(targetJSON),
		Input:                 scheduleInput(target, schedule.Name),
		AttemptCount:          2, // Two immediate attempts already made.
		CreatedAt:             now,
		NextAttemptAt:         now.Add(computeRetryBackoff(3)),
		ActionAfterCompletion: string(schedule.ActionAfterCompletion),
	}

	rs, rsErr := e.getRetryStore(region)
	if rsErr != nil {
		logs.Error("No retry store available for region, routing to DLQ",
			logs.String("schedule", schedule.Name),
			logs.String("region", region),
			logs.Err(rsErr))
		e.routeToDLQ(ctx, schedule, target, "retry store unavailable: "+rsErr.Error())
		return
	}
	if sErr := rs.SaveRetryRecord(record); sErr != nil {
		logs.Error("Failed to persist retry record",
			logs.String("schedule", schedule.Name),
			logs.Err(sErr))
	}
}

// checkRetries scans all regions for due RetryRecords and attempts redelivery.
// Called on each ticker cycle by the background worker.
func (e *Engine) checkRetries() {
	if e.storageManager == nil {
		return
	}
	regions := e.storageManager.GetActiveRegions()
	now := time.Now()

	for _, region := range regions {
		rs, rsErr := e.getRetryStore(region)
		if rsErr != nil {
			continue
		}
		due, err := rs.GetDueRetryRecords(now)
		if err != nil {
			logs.Debug("Failed to get due retry records",
				logs.String("region", region),
				logs.Err(err))
			continue
		}
		for _, record := range due {
			rec := record // capture for goroutine
			e.wg.Add(1)
			go func() {
				defer e.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						logs.Error("scheduler: panic processing retry record",
							logs.String("recordId", rec.ID),
							logs.Any("panic", r))
					}
				}()
				// Propagate e.ctx so retry processing stops when the engine
				// shuts down. Previously e.ctx was used inside
				// processRetryRecord but the passed context was ignored,
				// causing retry goroutines to outlive the engine.
				e.processRetryRecord(e.ctx, rs, rec, now)
			}()
		}
	}
}

// processRetryRecord attempts redelivery of a single RetryRecord and handles
// the outcome: success -> delete, failure -> update next attempt or route to DLQ.
// The ctx parameter is propagated from checkRetries so that engine shutdown
// cancels in-flight retry processing.
func (e *Engine) processRetryRecord(ctx context.Context, rs *schedulerstore.RetryStore, record *schedulerstore.RetryRecord, now time.Time) {
	var target schedulerstore.Target
	if err := json.Unmarshal([]byte(record.Target), &target); err != nil {
		logs.Error("Failed to deserialise target from retry record",
			logs.String("recordId", record.ID),
			logs.Err(err))
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		return
	}

	schedule := &schedulerstore.Schedule{
		Name:                  record.ScheduleName,
		GroupName:             record.GroupName,
		Region:                record.Region,
		Target:                &target,
		ActionAfterCompletion: schedulerstore.ActionAfterCompletion(record.ActionAfterCompletion),
	}

	// Verify the schedule still exists before retrying. With delayed
	// auto-deletion (AWS-compliant), the schedule remains alive during the
	// entire retry lifecycle, so this check passes for ActionAfterCompletion
	// = DELETE schedules and the retry continues. If the user manually
	// deletes the schedule mid-retry, this check fails and the retry record
	// is discarded — matching the user's intent to cancel.
	schedStore := e.getStoreForSchedule(schedule)
	if schedStore != nil {
		if _, err := schedStore.GetSchedule(ctx, record.GroupName, record.ScheduleName); err != nil {
			logs.Debug("Schedule no longer exists, discarding retry record",
				logs.String("schedule", record.ScheduleName),
				logs.String("recordId", record.ID))
			_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
			return
		}
	}

	maxRetries, maxAgeSeconds := retryDefaults(&target)

	// Check if MaximumEventAgeInSeconds has been exceeded.
	if now.Sub(record.CreatedAt) > time.Duration(maxAgeSeconds)*time.Second {
		logs.Warn("Retry record expired (MaximumEventAgeInSeconds exceeded)",
			logs.String("schedule", record.ScheduleName),
			logs.String("recordId", record.ID),
			logs.Int("attempts", record.AttemptCount))
		e.routeToDLQ(ctx, schedule, &target, "maximum event age exceeded")
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		// Retry lifecycle complete — auto-delete if configured.
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	// Check if MaximumRetryAttempts has been exceeded.
	// AWS: MaximumRetryAttempts=N → N+1 total delivery attempts
	// (1 initial + N retries). Use > not >= so we don't lose one retry.
	if record.AttemptCount > maxRetries {
		logs.Warn("Retry policy exhausted",
			logs.String("schedule", record.ScheduleName),
			logs.String("recordId", record.ID),
			logs.Int("attempts", record.AttemptCount),
			logs.Int("maxRetries", maxRetries))
		e.routeToDLQ(ctx, schedule, &target, "retry policy exhausted")
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		// Retry lifecycle complete — auto-delete if configured.
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	// Attempt delivery.
	err := e.deliverToTarget(ctx, schedule, &target)
	if err == nil {
		logs.Debug("Retry delivery succeeded",
			logs.String("schedule", record.ScheduleName),
			logs.String("recordId", record.ID),
			logs.Int("attempts", record.AttemptCount+1))
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		// Retry lifecycle complete — auto-delete if configured.
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	// Delivery failed — persist the updated record with the new NextAttemptAt
	// BEFORE deleting the old key (keyed by old NextAttemptAt). This order
	// prevents data loss if Save fails after Delete succeeds (e.g. disk full).
	// If Save fails, the old record remains and will be retried at its
	// original NextAttemptAt, which is safe (at-least-once).
	oldNextAttempt := record.NextAttemptAt
	record.AttemptCount++
	record.NextAttemptAt = now.Add(computeRetryBackoff(record.AttemptCount))
	if sErr := rs.SaveRetryRecord(record); sErr != nil {
		logs.Error("Failed to update retry record (old record retained)",
			logs.String("recordId", record.ID),
			logs.Err(sErr))
		return
	}
	_ = rs.DeleteRetryRecord(record.ID, oldNextAttempt)
}

// routeToDLQ sends a failed delivery to the DeadLetterConfig ARN. If no DLQ
// is configured, the message is discarded with an error log (AWS-compliant).
func (e *Engine) routeToDLQ(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target, reason string) {
	if target.DeadLetterConfig == nil || target.DeadLetterConfig.Arn == "" {
		logs.Error("Scheduler delivery permanently failed, no DLQ configured — discarding",
			logs.String("schedule", schedule.Name),
			logs.String("target", target.Arn),
			logs.String("reason", reason))
		return
	}

	dlqArn := target.DeadLetterConfig.Arn
	message := scheduleInput(target, schedule.Name)

	logs.Warn("Routing failed schedule delivery to DLQ",
		logs.String("schedule", schedule.Name),
		logs.String("dlqArn", dlqArn),
		logs.String("reason", reason))

	// DeadLetterConfig must be an SQS queue (AWS specification).
	_, dlqService, _, _, _ := svcarn.SplitARN(dlqArn)
	if dlqService != "sqs" {
		logs.Error("DeadLetterConfig ARN must reference an SQS queue",
			logs.String("dlqArn", dlqArn),
			logs.String("service", dlqService))
		return
	}

	sqsInvoker := e.bus.SQSInvoker()
	if sqsInvoker == nil {
		logs.Error("SQS invoker not available for DLQ delivery",
			logs.String("dlqArn", dlqArn))
		return
	}
	queueName := svcarn.ExtractQueueNameFromARN(dlqArn)
	_, _, dlqRegion, _, _ := svcarn.SplitARN(dlqArn)
	queueURL, qErr := sqsInvoker.GetQueueByName(ctx, dlqRegion, queueName)
	if qErr != nil {
		logs.Error("Failed to resolve DLQ queue URL",
			logs.String("dlqArn", dlqArn),
			logs.Err(qErr))
		return
	}
	// FIFO queues require MessageGroupId. Propagate from the
	// target's SqsParameters if available; otherwise, when the
	// destination queue ends with the FIFO suffix, fall back to
	// the schedule name so the SendMessage call satisfies the
	// SQS FIFO requirement.
	sendOpts := eventbus.SQSSendOptions{}
	if target.SqsParameters != nil && target.SqsParameters.MessageGroupId != "" {
		sendOpts.MessageGroupID = target.SqsParameters.MessageGroupId
	} else if strings.HasSuffix(queueName, ".fifo") {
		sendOpts.MessageGroupID = schedule.Name
	}
	if _, _, err := sqsInvoker.SendMessage(ctx, dlqRegion, queueURL, message, sendOpts); err != nil {
		logs.Error("Failed to send to DLQ",
			logs.String("dlqArn", dlqArn),
			logs.Err(err))
	}
}

// uuidString returns a UUID-like unique string for retry record IDs.
// Uses crypto/rand for uniqueness without external UUID dependency.
func uuidString() string {
	b := make([]byte, 16)
	_, _ = crand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (e *Engine) invokeLambda(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for Lambda invocation", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	input := scheduleInput(target, schedule.Name)

	functionName := svcarn.ExtractFunctionNameFromARN(target.Arn)
	if functionName == "" {
		logs.Debug("Failed to extract function name from ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid Lambda ARN: %s", target.Arn)
	}

	logs.Debug("Invoking Lambda for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName))

	statusCode, _, err := e.bus.LambdaInvoker().InvokeForGateway(ctx, target.Arn, []byte(input))
	if err != nil {
		logs.Debug("Failed to invoke Lambda",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Lambda invocation completed",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName),
		logs.Int("statusCode", int(statusCode)))
	return nil
}

func (e *Engine) sendToSQS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for SQS delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	sqsInvoker := e.bus.SQSInvoker()
	if sqsInvoker == nil {
		return fmt.Errorf("SQS invoker not available")
	}

	queueName := svcarn.ExtractQueueNameFromARN(target.Arn)
	if queueName == "" {
		logs.Debug("Invalid SQS ARN", logs.String("arn", target.Arn))
		return fmt.Errorf("invalid SQS ARN: %s", target.Arn)
	}

	_, _, sqsRegion, _, _ := svcarn.SplitARN(target.Arn)

	queueURL, qErr := sqsInvoker.GetQueueByName(ctx, sqsRegion, queueName)
	if qErr != nil {
		logs.Debug("SQS queue not found", logs.String("queue", queueName), logs.Err(qErr))
		return qErr
	}

	messageBody := scheduleInput(target, schedule.Name)

	// Honour SqsParameters.MessageGroupId for FIFO queues. AWS requires
	// MessageGroupId when the target is a FIFO queue. Fall back to the
	// schedule name if not explicitly set, matching routeToDLQ behaviour.
	sendOpts := eventbus.SQSSendOptions{}
	if target.SqsParameters != nil && target.SqsParameters.MessageGroupId != "" {
		sendOpts.MessageGroupID = target.SqsParameters.MessageGroupId
	} else if strings.HasSuffix(queueName, ".fifo") {
		sendOpts.MessageGroupID = schedule.Name
	}

	logs.Debug("Sending to SQS for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("queue", queueName))

	if _, _, err := sqsInvoker.SendMessage(ctx, sqsRegion, queueURL, messageBody, sendOpts); err != nil {
		logs.Debug("Failed to send to SQS",
			logs.String("schedule", schedule.Name),
			logs.String("queue", queueName),
			logs.Err(err))
		return err
	}
	return nil
}

func (e *Engine) publishToSNS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for SNS delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	snsInvoker := e.bus.SNSInvoker()
	if snsInvoker == nil {
		return fmt.Errorf("SNS invoker not available")
	}

	message := scheduleInput(target, schedule.Name)

	// Delegate to the SNS service's Publish API. This ensures the SNS
	// service handles topic policy evaluation, subscription filtering,
	// message persistence, fan-out to all subscription endpoints (SQS,
	// Lambda, HTTP, etc.), and EventBridge bus event publication.
	messageID, err := snsInvoker.PublishToTopic(ctx, target.Arn, message, "", nil)
	if err != nil {
		logs.Debug("Failed to publish to SNS topic",
			logs.String("schedule", schedule.Name),
			logs.String("topicArn", target.Arn),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("SNS delivery completed",
		logs.String("schedule", schedule.Name),
		logs.String("topic", target.Arn),
		logs.String("messageId", messageID))
	return nil
}

func (e *Engine) sendToCloudWatchLogs(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping CloudWatch Logs delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	logGroup := svcarn.ExtractLogGroupNameFromARN(target.Arn)
	if logGroup == "" {
		logs.Debug("Failed to extract log group from CloudWatch Logs ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid CloudWatch Logs ARN: %s", target.Arn)
	}

	_, _, region, _, resource := svcarn.SplitARN(target.Arn)
	var logStream string
	if idx := strings.LastIndex(resource, ":log-stream:"); idx != -1 {
		logStream = resource[idx+12:]
	} else {
		logStream = fmt.Sprintf("scheduler-%s", schedule.Name)
	}

	message := scheduleInput(target, schedule.Name)

	evt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup:  logGroup,
		LogStream: logStream,
		LogEvents: []eventbus.LogEntry{
			{Timestamp: time.Now().UnixMilli(), Message: message},
		},
	}
	evt.Region = region
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to deliver schedule to CloudWatch Logs",
			logs.String("schedule", schedule.Name),
			logs.String("logGroup", logGroup),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to CloudWatch Logs",
		logs.String("schedule", schedule.Name),
		logs.String("logGroup", logGroup))
	return nil
}

func (e *Engine) sendToKinesis(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for Kinesis delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	kinesisInvoker := e.bus.KinesisInvoker()
	if kinesisInvoker == nil {
		logs.Debug("Kinesis invoker not available", logs.String("schedule", schedule.Name))
		return fmt.Errorf("Kinesis invoker not available")
	}

	// Extract stream name from ARN: arn:aws:kinesis:<region>:<account>:stream/<name>
	_, _, kRegion, _, resource := svcarn.SplitARN(target.Arn)
	streamName := ""
	if idx := strings.LastIndex(resource, "/"); idx != -1 {
		streamName = resource[idx+1:]
	} else {
		streamName = resource
	}
	if streamName == "" {
		logs.Debug("Failed to extract stream name from Kinesis ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid Kinesis ARN: %s", target.Arn)
	}

	// Use PartitionKey from KinesisParameters if provided, otherwise fall
	// back to the schedule name (AWS uses a similar default behaviour).
	partitionKey := schedule.Name
	if target.KinesisParameters != nil && target.KinesisParameters.PartitionKey != "" {
		partitionKey = target.KinesisParameters.PartitionKey
	}

	data := []byte(scheduleInput(target, schedule.Name))

	logs.Debug("Sending to Kinesis for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("stream", streamName))

	if _, err := kinesisInvoker.PutRecord(ctx, streamName, partitionKey, data); err != nil {
		logs.Debug("Failed to send to Kinesis",
			logs.String("schedule", schedule.Name),
			logs.String("stream", streamName),
			logs.Err(err))
		return err
	}

	// Propagate region from the target ARN if the schedule does not carry one.
	if schedule.Region == "" && kRegion != "" {
		schedule.Region = kRegion
	}

	logs.Debug("Kinesis delivery completed",
		logs.String("schedule", schedule.Name),
		logs.String("stream", streamName))
	return nil
}

func (e *Engine) startStepFunctionExecution(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping Step Functions delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	_, _, smRegion, _, _ := svcarn.SplitARN(target.Arn)
	if smRegion == "" {
		smRegion = schedule.Region
	}

	input := scheduleInput(target, schedule.Name)

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: target.Arn,
		Input:           input,
	}
	evt.Region = smRegion
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to start Step Functions execution from schedule",
			logs.String("schedule", schedule.Name),
			logs.String("stateMachineArn", target.Arn),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to Step Functions",
		logs.String("schedule", schedule.Name),
		logs.String("stateMachineArn", target.Arn))
	return nil
}

func (e *Engine) sendToEventBridge(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping EventBridge delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	_, _, ebRegion, _, resource := svcarn.SplitARN(target.Arn)
	if ebRegion == "" {
		ebRegion = schedule.Region
	}

	eventBusName := "default"
	if idx := strings.LastIndex(resource, ":event-bus/"); idx != -1 {
		eventBusName = resource[idx+len(":event-bus/"):]
	}

	input := scheduleInput(target, schedule.Name)

	evt := &eventbus.EventBridgePutEventsEvent{
		EventBusName: eventBusName,
		Input:        input,
	}
	// Populate DetailType and Source from EventBridgeParameters so that
	// EventBridge rules can match on them.
	if target.EventBridgeParameters != nil {
		evt.DetailType = target.EventBridgeParameters.DetailType
		evt.Source = target.EventBridgeParameters.Source
	}
	evt.Region = ebRegion
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to deliver schedule to EventBridge",
			logs.String("schedule", schedule.Name),
			logs.String("eventBus", eventBusName),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to EventBridge",
		logs.String("schedule", schedule.Name),
		logs.String("eventBus", eventBusName))
	return nil
}
