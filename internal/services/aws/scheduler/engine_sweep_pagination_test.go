package scheduler

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// TestFiringSweepCoversSchedulesBeyondTheFirstPage pins the sweep's
// unbounded walk: with more enabled schedules in a region than the shared
// list layer's default page size, every one of them — including the
// key-order last — is evaluated and fired by the sweep. The paged form
// clamps a zero MaxItems to the default page size and discards the
// marker, which silently starved the schedules beyond the first page.
func TestFiringSweepCoversSchedulesBeyondTheFirstPage(t *testing.T) {
	svc := newLifecycleService(t)
	if err := svc.StartEngine(); err != nil {
		t.Fatalf("start engine: %v", err)
	}
	defer func() { _ = svc.StopEngine() }()

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })

	var mu sync.Mutex
	fired := make(map[string]bool)
	// The engine publishes ScheduleFiredEvent through the bus's async
	// path, which dispatches only to subscribers registered WithAsync —
	// the same option the service's production subscription carries.
	if _, err := eventbus.SubscribeTyped[*eventbus.ScheduleFiredEvent](bus,
		func(ctx context.Context, evt *eventbus.ScheduleFiredEvent) eventbus.HandlerResult {
			mu.Lock()
			fired[evt.ScheduleName] = true
			mu.Unlock()
			return eventbus.HandlerResult{}
		}, eventbus.WithAsync()); err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	svc.engine.SetEventBus(bus)

	store, err := svc.GetStoreForRegion(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get store: %v", err)
	}

	// One schedule more than the default page size, all immediately due
	// (past at() expressions fire on the next evaluation), zero-padded so
	// key order matches name order and the key-order last is exactly the
	// 101st.
	const count = pagination.DefaultMaxItems + 1
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("sweep-%03d", i)
		if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
			Name:                  name,
			GroupName:             "default",
			State:                 schedulerstore.ScheduleStateEnabled,
			ScheduleExpression:    "at(2020-01-01T00:00:00)",
			ActionAfterCompletion: "NONE",
			Target: &schedulerstore.Target{
				Arn:   "arn:aws:events:us-east-1:000000000000:event-bus/vocab-sweep",
				Input: fmt.Sprintf(`{"pin":%q}`, name),
			},
		}); err != nil {
			t.Fatalf("create schedule %s: %v", name, err)
		}
	}

	svc.engine.checkSchedules()

	deadline := time.Now().Add(3 * time.Second)
	for {
		mu.Lock()
		got := len(fired)
		missing := make([]string, 0)
		for i := 0; i < count; i++ {
			if !fired[fmt.Sprintf("sweep-%03d", i)] {
				missing = append(missing, fmt.Sprintf("sweep-%03d", i))
			}
		}
		mu.Unlock()
		if len(missing) == 0 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("sweep fired %d of %d schedules; missing (first 5): %v", got, count, missing[:min(5, len(missing))])
		}
		time.Sleep(20 * time.Millisecond)
	}
}
