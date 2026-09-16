package scheduler

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// newLifecycleService builds a real service + engine over a temp region
// storage manager with the store provider wired exactly as BuildEngine
// wires it in production. The cleanup stops the engine and closes the
// cached stores, so each test's ClientTokenStore reaper goroutine is
// released instead of lingering into later tests of the package.
func newLifecycleService(t *testing.T) *SchedulerService {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	t.Cleanup(func() { _ = sm.Close() })
	svc := NewSchedulerService(sm, "000000000000")
	t.Cleanup(func() { _ = svc.StopEngine() })
	svc.BuildEngine()
	return svc
}

// cleanupLoopGoroutines counts the goroutines running the ClientTokenStore
// cleanup loop, read from a full stack dump: the pinned resource is that
// one loop, so unrelated goroutine churn elsewhere in the process cannot
// flip the assertions. The buffer grows until the dump fits, so a truncated
// stack can never under-count.
func cleanupLoopGoroutines() int {
	buf := make([]byte, 1<<20)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return strings.Count(string(buf[:n]), "(*ClientTokenStore).cleanupLoop")
		}
		buf = make([]byte, 2*len(buf))
	}
}

// waitForCleanupLoops polls until the cleanup-loop count reaches want or a
// deadline passes: a stopped loop observes its channel only asynchronously,
// and a fixed sleep would reintroduce the timing dependence the count is
// scoped to remove.
func waitForCleanupLoops(want int) bool {
	deadline := time.Now().Add(2 * time.Second)
	for {
		if cleanupLoopGoroutines() == want {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestEngineStartStopCyclesDoNotLeakStoreGoroutines pins the unified-store
// ownership contract: the engine holds no store cache, so engine start/stop
// cycles create no per-region ClientTokenStore cleanup goroutines. The one
// store created on first use belongs to the service cache and is released
// once, by StopEngine.
func TestEngineStartStopCyclesDoNotLeakStoreGoroutines(t *testing.T) {
	svc := newLifecycleService(t)

	// First cycle: the provider creates one regional store (one cleanup
	// goroutine), owned by the service cache — it survives engine Stop by
	// design and is closed only at StopEngine.
	if err := svc.StartEngine(); err != nil {
		t.Fatalf("start engine: %v", err)
	}
	if err := svc.engine.Stop(); err != nil {
		t.Fatalf("stop engine: %v", err)
	}
	if !waitForCleanupLoops(1) {
		t.Fatalf("store cleanup goroutines after the first cycle = %d, want 1", cleanupLoopGoroutines())
	}

	// Further cycles must add nothing: the provider serves the cached
	// store and the engine owns no stores to leak.
	for i := 0; i < 4; i++ {
		if err := svc.StartEngine(); err != nil {
			t.Fatalf("start engine cycle %d: %v", i+2, err)
		}
		if err := svc.engine.Stop(); err != nil {
			t.Fatalf("stop engine cycle %d: %v", i+2, err)
		}
	}
	if got := cleanupLoopGoroutines(); got != 1 {
		t.Errorf("engine start/stop cycles leaked store cleanup goroutines: %d after five cycles, want 1", got)
	}

	// StopEngine is the single closer: the store's cleanup goroutine exits.
	if err := svc.StopEngine(); err != nil {
		t.Fatalf("stop engine service: %v", err)
	}
	if !waitForCleanupLoops(0) {
		t.Errorf("StopEngine did not release the ClientTokenStore cleanup goroutine: %d remain", cleanupLoopGoroutines())
	}
}

// TestRetryRecordProcessedOnFirstSweep pins the startup-sweep contract: a
// retry record already due at engine start is processed by the pre-loop
// sweep, without waiting a ticker interval. The record is seeded past its
// MaximumEventAgeInSeconds so the first processing expires it and deletes
// it — an observable transition within the poll window.
func TestRetryRecordProcessedOnFirstSweep(t *testing.T) {
	svc := newLifecycleService(t)

	store, err := svc.GetStoreForRegion(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get store: %v", err)
	}
	// A far-future at() schedule: never due, so the schedule sweep stays
	// silent and the only observable engine action is the retry sweep. The
	// schedule must exist so processRetryRecord's existence check passes.
	if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
		Name:                  "retry-first-sweep",
		GroupName:             "default",
		State:                 schedulerstore.ScheduleStateEnabled,
		ScheduleExpression:    "at(2030-01-01T00:00:00)",
		ActionAfterCompletion: "NONE",
		Target: &schedulerstore.Target{
			Arn: "arn:aws:lambda:us-east-1:000000000000:function:retry-first-sweep",
		},
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}

	rs, err := svc.engine.getRetryStore(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get retry store: %v", err)
	}
	now := time.Now()
	record := &schedulerstore.RetryRecord{
		ID:           "retry-first-sweep-record",
		ScheduleName: "retry-first-sweep",
		GroupName:    "default",
		Region:       defaults.DefaultRegion,
		Target:       `{"Arn":"arn:aws:lambda:us-east-1:000000000000:function:retry-first-sweep"}`,
		AttemptCount: 2,
		// Created 25h ago: beyond the 86400s MaximumEventAgeInSeconds
		// default, so the first processing expires and deletes it.
		CreatedAt:     now.Add(-25 * time.Hour),
		NextAttemptAt: now.Add(-time.Minute),
	}
	if err := rs.SaveRetryRecord(record); err != nil {
		t.Fatalf("seed retry record: %v", err)
	}

	if err := svc.StartEngine(); err != nil {
		t.Fatalf("start engine: %v", err)
	}
	defer func() { _ = svc.StopEngine() }()

	// The pre-loop sweep processes within milliseconds; even the TEST_MODE
	// ticker interval (1s) lies outside this window, so passing it pins the
	// startup sweep rather than the first tick.
	deadline := time.Now().Add(600 * time.Millisecond)
	for {
		due, err := rs.GetDueRetryRecords(time.Now())
		if err != nil {
			t.Fatalf("get due retry records: %v", err)
		}
		if len(due) == 0 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("due retry record was not processed by the startup sweep within 600ms of engine start")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// capturingSQSInvoker stands in for the SQS service on the far side of the
// bus for dead-letter pins: every queue resolves and every send succeeds,
// with each sent body kept for assertion.
type capturingSQSInvoker struct {
	sent   int
	bodies []string
}

func (c *capturingSQSInvoker) GetQueueByName(ctx context.Context, region, queueName string) (string, error) {
	return "http://sqs.local/" + queueName, nil
}

func (c *capturingSQSInvoker) GetQueueARN(ctx context.Context, region, queueURL string) (string, error) {
	return "", nil
}

func (c *capturingSQSInvoker) SendMessage(ctx context.Context, region, queueURL, body string, opts invokers.SQSSendOptions) (string, string, error) {
	c.sent++
	c.bodies = append(c.bodies, body)
	return "dlq-message-id", "md5-of-body", nil
}

func (c *capturingSQSInvoker) ReceiveMessage(ctx context.Context, region, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]invokers.ReceivedSQSMessage, error) {
	return nil, nil
}

func (c *capturingSQSInvoker) DeleteMessage(ctx context.Context, region, queueURL, receiptHandle string) error {
	return nil
}

// TestRetryRecordPersistenceFailureRoutesToDLQAndCompletes pins the
// persistence-failure contract: when both immediate attempts fail and the
// RetryRecord cannot be saved, no background retry will ever run for the
// delivery — the engine must close the lifecycle itself, routing the input
// to the dead-letter queue and applying ActionAfterCompletion exactly as it
// does when retries are exhausted. A bare return here would silently drop
// the delivery and leave a one-time schedule unmarked (or a DELETE-action
// schedule alive firing nothing) forever.
func TestRetryRecordPersistenceFailureRoutesToDLQAndCompletes(t *testing.T) {
	svc := newLifecycleService(t)

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	dlq := &capturingSQSInvoker{}
	bus.SetSQSInvoker(dlq)
	// No Lambda invoker registered: both delivery attempts fail, which is
	// the precondition under which the retry persistence path runs.
	svc.engine.SetEventBus(bus)

	// Inject a failing retry store: a RetryStore over a storage handle that
	// has been closed, cached as the engine's store for the region. Every
	// SaveRetryRecord fails, while the schedule store over the healthy
	// manager stays writable so the completion action remains observable.
	failMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("failing storage manager: %v", err)
	}
	failStorage, err := failMgr.GetStorage(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("failing storage: %v", err)
	}
	brokenRetryStore := schedulerstore.NewRetryStore(failStorage, defaults.DefaultRegion)
	if err := failMgr.Close(); err != nil {
		t.Fatalf("close failing storage: %v", err)
	}
	svc.engine.retryStores.Store(defaults.DefaultRegion, brokenRetryStore)

	store, err := svc.GetStoreForRegion(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get store: %v", err)
	}

	targets := func(pin string) *schedulerstore.Target {
		return &schedulerstore.Target{
			Arn:              "arn:aws:lambda:us-east-1:000000000000:function:unwired",
			Input:            `{"pin":"` + pin + `"}`,
			DeadLetterConfig: &schedulerstore.DeadLetterConfig{Arn: "arn:aws:sqs:us-east-1:000000000000:retry-dlq"},
		}
	}
	schedules := []struct {
		name   string
		action string
	}{
		{"retry-persist-fail-none", "NONE"},
		{"retry-persist-fail-delete", "DELETE"},
	}
	for _, s := range schedules {
		if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
			Name:                  s.name,
			GroupName:             "default",
			State:                 schedulerstore.ScheduleStateEnabled,
			ScheduleExpression:    "at(2030-01-01T00:00:00)",
			ActionAfterCompletion: schedulerstore.ActionAfterCompletion(s.action),
			Target:                targets(s.name),
		}); err != nil {
			t.Fatalf("create %s schedule: %v", s.name, err)
		}
		sched, err := store.GetSchedule(t.Context(), "default", s.name)
		if err != nil {
			t.Fatalf("get %s schedule: %v", s.name, err)
		}
		svc.engine.deliverWithRetry(t.Context(), sched, sched.Target)
	}

	if dlq.sent != 2 {
		t.Errorf("DLQ sends = %d, want 2 (one per delivery dropped by the persistence failure)", dlq.sent)
	} else {
		for i, s := range schedules {
			if !strings.Contains(dlq.bodies[i], s.name) {
				t.Errorf("DLQ body %d = %q, want it to carry the input pinning %q", i, dlq.bodies[i], s.name)
			}
		}
	}

	completed, err := store.GetSchedule(t.Context(), "default", "retry-persist-fail-none")
	if err != nil {
		t.Fatalf("get none-action schedule after delivery: %v", err)
	}
	if completed.CompletionDate == nil {
		t.Error("one-time schedule was not marked completed when its retry record could not be persisted")
	}

	_, err = store.GetSchedule(t.Context(), "default", "retry-persist-fail-delete")
	if !errors.Is(err, schedulerstore.ErrScheduleNotFound) {
		t.Errorf("ActionAfterCompletion=DELETE schedule after delivery: error = %v, want ErrScheduleNotFound", err)
	}
}

// TestRetryRecordDiscardedOnlyOnConfirmedNotFound pins the existence-probe
// contract of processRetryRecord: a storage failure on the probe retains the
// pending retry record for the next sweep (deleting it would drop a delivery
// without DLQ routing), while a confirmed not-found — the user deleted the
// schedule mid-retry — discards it.
func TestRetryRecordDiscardedOnlyOnConfirmedNotFound(t *testing.T) {
	// Retention: the probe's schedule store is faulted, so GetSchedule
	// reports a storage failure rather than the not-found sentinel.
	svc := newLifecycleService(t)

	rs, err := svc.engine.getRetryStore(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get retry store: %v", err)
	}
	now := time.Now()
	record := &schedulerstore.RetryRecord{
		ID:            "retry-probe-storage-fail",
		ScheduleName:  "probe-storage-fail",
		GroupName:     "default",
		Region:        defaults.DefaultRegion,
		Target:        `{"Arn":"arn:aws:lambda:us-east-1:000000000000:function:probe"}`,
		AttemptCount:  2,
		CreatedAt:     now,
		NextAttemptAt: now.Add(-time.Minute),
	}
	if err := rs.SaveRetryRecord(record); err != nil {
		t.Fatalf("seed retry record: %v", err)
	}

	failMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("failing storage manager: %v", err)
	}
	failStorage, err := failMgr.GetStorage(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("failing regional storage: %v", err)
	}
	faultedStore := schedulerstore.NewSchedulerStore(failStorage, "000000000000", defaults.DefaultRegion)
	// The construction spawned the idempotency-token reaper over storage
	// that is about to be faulted; stop it so it cannot linger past the
	// test.
	t.Cleanup(func() { faultedStore.Close() })
	if err := failMgr.Close(); err != nil {
		t.Fatalf("close failing storage manager: %v", err)
	}
	svc.engine.SetStoreProvider(func(region string) (*schedulerstore.SchedulerStore, error) {
		return faultedStore, nil
	})

	svc.engine.processRetryRecord(t.Context(), rs, record, now)

	due, err := rs.GetDueRetryRecords(now.Add(time.Hour))
	if err != nil {
		t.Fatalf("list retry records after faulted probe: %v", err)
	}
	retained := false
	for _, rec := range due {
		if rec.ID == record.ID {
			retained = true
		}
	}
	if !retained {
		t.Error("retry record was discarded on a storage failure of the existence probe, want retained")
	}

	// Discard: the store is healthy and the schedule genuinely does not
	// exist, so the probe's not-found is confirmed.
	svc2 := newLifecycleService(t)
	rs2, err := svc2.engine.getRetryStore(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get retry store (healthy): %v", err)
	}
	gone := &schedulerstore.RetryRecord{
		ID:            "retry-probe-not-found",
		ScheduleName:  "deleted-mid-retry",
		GroupName:     "default",
		Region:        defaults.DefaultRegion,
		Target:        `{"Arn":"arn:aws:lambda:us-east-1:000000000000:function:probe"}`,
		AttemptCount:  2,
		CreatedAt:     now,
		NextAttemptAt: now.Add(-time.Minute),
	}
	if err := rs2.SaveRetryRecord(gone); err != nil {
		t.Fatalf("seed retry record (healthy): %v", err)
	}
	svc2.engine.processRetryRecord(t.Context(), rs2, gone, now)
	due2, err := rs2.GetDueRetryRecords(now.Add(time.Hour))
	if err != nil {
		t.Fatalf("list retry records after not-found probe: %v", err)
	}
	for _, rec := range due2 {
		if rec.ID == gone.ID {
			t.Error("retry record for a deleted schedule was retained, want discarded")
		}
	}
}

// TestBusDeliveryCompletesOnetimeSchedule pins the production wiring of the
// one-time completion marker end to end: the firing publishes through the
// bus, the service's own subscription reconstructs the schedule WITH its
// expression, and a successful delivery ends the lifecycle — the completed
// schedule leaves the engine's enabled sweep. A reconstruction that dropped
// the expression left the marker unwritable on the only production
// delivery path.
func TestBusDeliveryCompletesOnetimeSchedule(t *testing.T) {
	svc := newLifecycleService(t)

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	bus.SetLambdaInvoker(&recordingUniversalLambdaInvoker{})
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("set event bus: %v", err)
	}

	store, err := svc.GetStoreForRegion(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get store: %v", err)
	}
	if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
		Name:               "bus-complete-pin",
		GroupName:          "default",
		State:              schedulerstore.ScheduleStateEnabled,
		ScheduleExpression: "at(2020-01-01T00:00:00)",
		Target: &schedulerstore.Target{
			Arn: "arn:aws:lambda:us-east-1:000000000000:function:bus-complete-pin",
		},
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	loaded, err := store.GetSchedule(t.Context(), "default", "bus-complete-pin")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	loaded.Region = defaults.DefaultRegion

	if err := svc.engine.executeSchedule(t.Context(), loaded); err != nil {
		t.Fatalf("executeSchedule: %v", err)
	}

	deadline := time.Now().Add(3 * time.Second)
	for {
		got, err := store.GetSchedule(t.Context(), "default", "bus-complete-pin")
		if err != nil {
			t.Fatalf("get schedule after delivery: %v", err)
		}
		if got.CompletionDate != nil {
			enabled, err := store.GetAllEnabledSchedules(t.Context())
			if err != nil {
				t.Fatalf("list enabled: %v", err)
			}
			for _, sch := range enabled {
				if sch.Name == "bus-complete-pin" {
					t.Error("completed one-time schedule stayed in the enabled sweep")
				}
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("one-time schedule was not marked completed through the bus delivery path")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestRetryDeliveryCompletesOnetimeSchedule pins the retry-path half of the
// completion wiring: a one-time schedule whose first attempts failed ends
// its lifecycle when the retried delivery succeeds, the fire-time
// expression riding the retry record beside ActionAfterCompletion.
func TestRetryDeliveryCompletesOnetimeSchedule(t *testing.T) {
	svc := newLifecycleService(t)

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	bus.SetLambdaInvoker(&recordingUniversalLambdaInvoker{})
	svc.engine.SetEventBus(bus)

	store, err := svc.GetStoreForRegion(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get store: %v", err)
	}
	if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
		Name:               "retry-complete-pin",
		GroupName:          "default",
		State:              schedulerstore.ScheduleStateEnabled,
		ScheduleExpression: "at(2020-01-01T00:00:00)",
		Target: &schedulerstore.Target{
			Arn: "arn:aws:lambda:us-east-1:000000000000:function:retry-complete-pin",
		},
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}

	rs, err := svc.engine.getRetryStore(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get retry store: %v", err)
	}
	now := time.Now()
	record := &schedulerstore.RetryRecord{
		ID:                 "retry-complete-pin-record",
		ScheduleName:       "retry-complete-pin",
		GroupName:          "default",
		Region:             defaults.DefaultRegion,
		Target:             `{"Arn":"arn:aws:lambda:us-east-1:000000000000:function:retry-complete-pin"}`,
		AttemptCount:       2,
		CreatedAt:          now,
		NextAttemptAt:      now.Add(-time.Minute),
		ScheduleExpression: "at(2020-01-01T00:00:00)",
	}
	if err := rs.SaveRetryRecord(record); err != nil {
		t.Fatalf("seed retry record: %v", err)
	}

	svc.engine.processRetryRecord(t.Context(), rs, record, now)

	got, err := store.GetSchedule(t.Context(), "default", "retry-complete-pin")
	if err != nil {
		t.Fatalf("get schedule after retry: %v", err)
	}
	if got.CompletionDate == nil {
		t.Error("one-time schedule was not marked completed by the successful retry delivery")
	}
	due, err := rs.GetDueRetryRecords(now.Add(time.Hour))
	if err != nil {
		t.Fatalf("list retry records: %v", err)
	}
	for _, rec := range due {
		if rec.ID == record.ID {
			t.Error("retry record survived its successful delivery")
		}
	}
}
