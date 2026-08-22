package scheduler

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// TestShouldExecuteNonFlexibleLateBoundary pins that a non-flexible
// schedule fires a boundary missed by a slow ticker on the next evaluation
// instead of silently skipping it, and fires that boundary exactly once.
func TestShouldExecuteNonFlexibleLateBoundary(t *testing.T) {
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:               "late-boundary",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "rate(5 minutes)",
		CreationDate:       time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	// The 12:05 boundary's window passed without an evaluation; the 12:07
	// evaluation must still fire it.
	late := time.Date(2027, 1, 1, 12, 7, 0, 0, time.UTC)
	if !e.shouldExecute(sch, late) {
		t.Error("late evaluation should fire the missed boundary")
	}
	// The caller records the fire; a later evaluation in the same period
	// must not fire again.
	key := lastFiredKey(sch.Region, sch.GroupName, sch.Name)
	e.lastFired.Store(key, lastFiredEntry{firedAt: late, expr: sch.ScheduleExpression})
	if e.shouldExecute(sch, late.Add(30*time.Second)) {
		t.Error("boundary fired twice")
	}
}

// TestShouldExecuteCronFutureMatchStaysSilent pins that a cron schedule
// stays silent until the evaluation minute matches.
func TestShouldExecuteCronFutureMatchStaysSilent(t *testing.T) {
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:               "cron-morning",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "cron(0 6 * * ? *)",
		CreationDate:       time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	if e.shouldExecute(sch, time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)) {
		t.Error("cron schedule fired on a non-matching minute")
	}
	if !e.shouldExecute(sch, time.Date(2027, 1, 2, 6, 0, 0, 0, time.UTC)) {
		t.Error("cron schedule did not fire on its matching minute")
	}
}

// TestShouldExecuteAtPastFiresOnce pins that a one-time at() schedule whose
// minute passed while unobserved still fires on the next evaluation, once.
func TestShouldExecuteAtPastFiresOnce(t *testing.T) {
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:               "one-time",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "at(2027-01-01T11:00:00)",
		CreationDate:       time.Date(2027, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	// The at() time lies two hours in the past at this evaluation.
	past := time.Date(2027, 1, 1, 13, 0, 0, 0, time.UTC)
	if !e.shouldExecute(sch, past) {
		t.Error("past at() schedule never fires")
	}
	key := lastFiredKey(sch.Region, sch.GroupName, sch.Name)
	e.lastFired.Store(key, lastFiredEntry{firedAt: past, expr: sch.ScheduleExpression})
	if e.shouldExecute(sch, past.Add(time.Minute)) {
		t.Error("at() schedule fired twice")
	}
}

// TestShouldExecuteExpressionChangeBreaksDedup pins that the dedup cache
// is scoped to the expression: updating a completed one-time schedule to
// a new past at() starts a new firing lifecycle (UpdateSchedule clears
// the completion marker with it), while a target-only update that keeps
// the expression continues to suppress the fired boundary.
func TestShouldExecuteExpressionChangeBreaksDedup(t *testing.T) {
	e := &Engine{}
	original := &schedulerstore.Schedule{
		Name:               "one-time",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "at(2027-01-01T11:00:00)",
		CreationDate:       time.Date(2027, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	fired := time.Date(2027, 1, 1, 13, 0, 0, 0, time.UTC)
	key := lastFiredKey(original.Region, original.GroupName, original.Name)
	e.lastFired.Store(key, lastFiredEntry{firedAt: fired, expr: original.ScheduleExpression})

	// Target-only update: same expression, the fired boundary stays
	// suppressed.
	if e.shouldExecute(original, fired.Add(time.Minute)) {
		t.Error("target-only update re-fired the suppressed boundary")
	}

	// Expression change to an earlier at(): the new lifecycle's boundary
	// lies before the stored fire time, so only expression-aware dedup
	// lets it through.
	relifecycled := *original
	relifecycled.ScheduleExpression = "at(2027-01-01T10:30:00)"
	if !e.shouldExecute(&relifecycled, fired.Add(time.Minute)) {
		t.Error("expression change did not break the dedup suppression")
	}
}

// TestShouldExecuteFlexibleWindowUnchanged pins the contractual flexible
// time window: execution no earlier than the boundary, within
// MaximumWindowInMinutes of it, and never after the window closes.
func TestShouldExecuteFlexibleWindowUnchanged(t *testing.T) {
	window := 10
	newFlex := func(name, expr string, creation time.Time) *schedulerstore.Schedule {
		return &schedulerstore.Schedule{
			Name:               name,
			GroupName:          "default",
			Region:             "us-east-1",
			ScheduleExpression: expr,
			CreationDate:       creation,
			FlexibleTimeWindow: &schedulerstore.FlexibleTimeWindow{
				Mode:                   schedulerstore.FlexibleTimeWindowModeFlexible,
				MaximumWindowInMinutes: &window,
			},
		}
	}

	// Before the first boundary opens its window: silent.
	e := &Engine{}
	rate := newFlex("flex-rate", "rate(5 minutes)", time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC))
	if e.shouldExecute(rate, time.Date(2027, 1, 1, 11, 59, 0, 0, time.UTC)) {
		t.Error("fired before the flexible window opened")
	}
	// Inside the window of the latest boundary: fires.
	if !e.shouldExecute(rate, time.Date(2027, 1, 1, 12, 7, 0, 0, time.UTC)) {
		t.Error("did not fire inside the flexible window")
	}

	// A one-time at() schedule whose window has closed: the occurrence is
	// skipped, late evaluations do not resurrect it.
	e2 := &Engine{}
	oneTime := newFlex("flex-once", "at(2027-01-01T11:00:00)", time.Date(2027, 1, 1, 10, 0, 0, 0, time.UTC))
	if e2.shouldExecute(oneTime, time.Date(2027, 1, 1, 13, 0, 0, 0, time.UTC)) {
		t.Error("fired after the flexible window closed")
	}
}

// TestShouldExecuteRateWaitsFirstInterval pins the AWS rate() contract:
// the schedule does not fire on the creation boundary — the first
// invocation happens one full interval after creation.
func TestShouldExecuteRateWaitsFirstInterval(t *testing.T) {
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:               "rate-first",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "rate(5 minutes)",
		CreationDate:       time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	if e.shouldExecute(sch, time.Date(2027, 1, 1, 12, 0, 5, 0, time.UTC)) {
		t.Error("rate schedule fired on the creation boundary")
	}
	if !e.shouldExecute(sch, time.Date(2027, 1, 1, 12, 5, 1, 0, time.UTC)) {
		t.Error("rate schedule did not fire one interval after creation")
	}
}

// TestShouldExecuteFlexibleWaitsFirstInterval pins that the flexible
// window of the creation boundary never opens: the first window belongs
// to the first real boundary, one interval after creation.
func TestShouldExecuteFlexibleWaitsFirstInterval(t *testing.T) {
	window := 10
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:               "flex-first",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "rate(5 minutes)",
		CreationDate:       time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC),
		FlexibleTimeWindow: &schedulerstore.FlexibleTimeWindow{
			Mode:                   schedulerstore.FlexibleTimeWindowModeFlexible,
			MaximumWindowInMinutes: &window,
		},
	}
	if e.shouldExecute(sch, time.Date(2027, 1, 1, 12, 3, 0, 0, time.UTC)) {
		t.Error("flexible schedule fired inside the creation boundary's window")
	}
	if !e.shouldExecute(sch, time.Date(2027, 1, 1, 12, 6, 0, 0, time.UTC)) {
		t.Error("flexible schedule did not fire inside the first real boundary's window")
	}
}

// TestShouldExecuteCronLateBoundary pins that a cron boundary missed by
// a slow ticker fires on the next evaluation instead of being lost, and
// fires exactly once.
func TestShouldExecuteCronLateBoundary(t *testing.T) {
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:               "cron-late",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "cron(0/5 * * * ? *)",
		CreationDate:       time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	// The 12:05 matching minute passed unobserved; the 12:07 evaluation
	// must still fire it.
	late := time.Date(2027, 1, 1, 12, 7, 0, 0, time.UTC)
	if !e.shouldExecute(sch, late) {
		t.Error("late evaluation should fire the missed cron boundary")
	}
	key := lastFiredKey(sch.Region, sch.GroupName, sch.Name)
	e.lastFired.Store(key, lastFiredEntry{firedAt: late, expr: sch.ScheduleExpression})
	if e.shouldExecute(sch, late.Add(30*time.Second)) {
		t.Error("missed cron boundary fired twice")
	}
}

// TestShouldExecuteAtInScheduleTimezone pins that the offset-less at()
// timestamp is interpreted in ScheduleExpressionTimezone: 09:00 Tokyo
// corresponds to 00:00 UTC, so an evaluation at 00:30 UTC has passed it
// while an evaluation the previous UTC evening has not.
func TestShouldExecuteAtInScheduleTimezone(t *testing.T) {
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:                       "one-time-tz",
		GroupName:                  "default",
		Region:                     "us-east-1",
		ScheduleExpression:         "at(2027-01-01T09:00:00)",
		ScheduleExpressionTimezone: "Asia/Tokyo",
		CreationDate:               time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
	}
	if !e.shouldExecute(sch, time.Date(2027, 1, 1, 0, 30, 0, 0, time.UTC)) {
		t.Error("at() schedule with Tokyo timezone did not fire at 09:30 Tokyo time")
	}
	e2 := &Engine{}
	if e2.shouldExecute(sch, time.Date(2026, 12, 31, 15, 30, 0, 0, time.UTC)) {
		t.Error("at() schedule with Tokyo timezone fired before 09:00 Tokyo time")
	}
}

// --- persisted delivered-boundary marker ---

func newFiringStore(t *testing.T) *schedulerstore.SchedulerStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return schedulerstore.NewSchedulerStore(st, "000000000000", "us-east-1")
}

// TestDeliveredBoundaryMarkerSuppressesRefire pins that the record's
// delivered-boundary marker suppresses the boundary it records without any
// in-memory state, and never suppresses a newer boundary.
func TestDeliveredBoundaryMarkerSuppressesRefire(t *testing.T) {
	e := &Engine{}
	sch := &schedulerstore.Schedule{
		Name:               "marker",
		GroupName:          "default",
		Region:             "us-east-1",
		ScheduleExpression: "rate(5 minutes)",
		CreationDate:       time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	boundary := time.Date(2027, 1, 1, 12, 5, 0, 0, time.UTC)
	if _, ok := e.dueBoundary(sch, boundary.Add(2*time.Minute)); !ok {
		t.Fatal("schedule not due before the boundary was delivered")
	}
	delivered := boundary
	sch.LastFiredAt = &delivered
	if _, ok := e.dueBoundary(sch, boundary.Add(2*time.Minute)); ok {
		t.Error("marker did not suppress the delivered boundary")
	}
	older := boundary.Add(-1 * time.Minute)
	sch.LastFiredAt = &older
	if _, ok := e.dueBoundary(sch, boundary.Add(2*time.Minute)); !ok {
		t.Error("marker suppressed a boundary newer than itself")
	}
}

// TestDeliveredBoundaryNotRefiredAfterRestart pins the restart contract:
// once the delivery path persists the boundary on the record, a fresh
// engine (empty in-memory state) reading the record back must not deliver
// the boundary again — for rate() and cron() schedules alike — while a
// genuinely new boundary still fires.
func TestDeliveredBoundaryNotRefiredAfterRestart(t *testing.T) {
	for _, tc := range []struct {
		name string
		expr string
	}{
		{"rate", "rate(5 minutes)"},
		{"cron", "cron(0/5 * * * ? *)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := newFiringStore(t)
			now := time.Now().UTC()
			if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
				Name:               "recurring",
				GroupName:          "default",
				State:              schedulerstore.ScheduleStateEnabled,
				ScheduleExpression: tc.expr,
			}); err != nil {
				t.Fatalf("create schedule: %v", err)
			}
			// CreateSchedule stamps CreationDate with now; backdate it so a
			// boundary has already elapsed.
			if err := store.MutateSchedule(t.Context(), "default", "recurring", func(s *schedulerstore.Schedule) error {
				s.CreationDate = now.Add(-11 * time.Minute)
				return nil
			}); err != nil {
				t.Fatalf("backdate creation: %v", err)
			}

			loaded, err := store.GetSchedule(t.Context(), "default", "recurring")
			if err != nil {
				t.Fatalf("get schedule: %v", err)
			}
			e1 := &Engine{}
			boundary, ok := e1.dueBoundary(loaded, now)
			if !ok {
				t.Fatal("boundary not due before delivery")
			}
			// What the delivery goroutine does after a successful delivery.
			if err := store.TouchScheduleLastFired(t.Context(), "default", "recurring", boundary); err != nil {
				t.Fatalf("touch last fired: %v", err)
			}

			// Restart: the engine's in-memory state is gone; the record is
			// read back from the store.
			restarted, err := store.GetSchedule(t.Context(), "default", "recurring")
			if err != nil {
				t.Fatalf("get schedule after restart: %v", err)
			}
			e2 := &Engine{}
			if _, ok := e2.dueBoundary(restarted, now); ok {
				t.Error("delivered boundary re-fired after restart")
			}
			// A genuinely new boundary still fires after the restart.
			if _, ok := e2.dueBoundary(restarted, now.Add(6*time.Minute)); !ok {
				t.Error("new boundary did not fire after restart")
			}
		})
	}
}

// TestCompletedScheduleUpdateRefiresConsistently pins the unified update
// contract: UpdateSchedule overrides the whole schedule and resets
// execution state, so updating a completed one-time schedule — even with
// the same expression and only a Description change — clears the completion
// marker and the delivered-boundary marker, and the schedule fires again
// the same way a restarted process would see it.
func TestCompletedScheduleUpdateRefiresConsistently(t *testing.T) {
	store := newFiringStore(t)
	now := time.Now().UTC()
	atExpr := "at(" + now.Add(-time.Hour).UTC().Format("2006-01-02T15:04:05") + ")"
	if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
		Name:               "one-time",
		GroupName:          "default",
		State:              schedulerstore.ScheduleStateEnabled,
		ScheduleExpression: atExpr,
		CreationDate:       now.Add(-2 * time.Hour),
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}

	// Lifecycle one: the boundary is delivered and the lifecycle completes.
	loaded, err := store.GetSchedule(t.Context(), "default", "one-time")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	e := &Engine{}
	boundary, ok := e.dueBoundary(loaded, now)
	if !ok {
		t.Fatal("past at() schedule not due")
	}
	if err := store.TouchScheduleLastFired(t.Context(), "default", "one-time", boundary); err != nil {
		t.Fatalf("touch last fired: %v", err)
	}
	if err := store.CompleteSchedule(t.Context(), "default", "one-time"); err != nil {
		t.Fatalf("complete schedule: %v", err)
	}

	// Update with the SAME expression, changing only the Description.
	svc := &SchedulerService{}
	if _, err := svc.updateScheduleCore(t.Context(), store, &UpdateScheduleInput{Spec: &ScheduleSpec{
		Name:                  "one-time",
		GroupName:             "default",
		ScheduleExpression:    atExpr,
		Description:           "after",
		State:                 "ENABLED",
		ActionAfterCompletion: "NONE",
		Target: &schedulerstore.Target{
			Arn:     "arn:aws:lambda:us-east-1:000000000000:function:refire-test",
			RoleArn: "arn:aws:iam::000000000000:role/scheduler-test",
		},
	}}); err != nil {
		t.Fatalf("updateScheduleCore: %v", err)
	}

	updated, err := store.GetSchedule(t.Context(), "default", "one-time")
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if updated.Description != "after" {
		t.Errorf("description = %q, want %q", updated.Description, "after")
	}
	if updated.CompletionDate != nil {
		t.Error("update did not re-lifecycle the completed schedule")
	}
	if updated.LastFiredAt != nil {
		t.Error("update did not reset the delivered-boundary marker")
	}
	// The re-lifecycled schedule fires again, identically for a fresh
	// engine (the restart view) and the running process, whose slot was
	// released after the first delivery.
	if _, ok := (&Engine{}).dueBoundary(updated, now); !ok {
		t.Error("re-lifecycled schedule did not fire again")
	}
}

// TestActiveScheduleUpdateKeepsDeliveredMarker pins the other side of the
// update contract: updating an active schedule that keeps its expression
// must not reset the delivered-boundary marker, so the already delivered
// boundary is not re-fired.
func TestActiveScheduleUpdateKeepsDeliveredMarker(t *testing.T) {
	store := newFiringStore(t)
	now := time.Now().UTC()
	if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
		Name:               "recurring",
		GroupName:          "default",
		State:              schedulerstore.ScheduleStateEnabled,
		ScheduleExpression: "rate(5 minutes)",
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	// CreateSchedule stamps CreationDate with now; backdate it so a
	// boundary has already elapsed.
	if err := store.MutateSchedule(t.Context(), "default", "recurring", func(s *schedulerstore.Schedule) error {
		s.CreationDate = now.Add(-11 * time.Minute)
		return nil
	}); err != nil {
		t.Fatalf("backdate creation: %v", err)
	}
	loaded, err := store.GetSchedule(t.Context(), "default", "recurring")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	boundary, ok := (&Engine{}).dueBoundary(loaded, now)
	if !ok {
		t.Fatal("boundary not due before delivery")
	}
	if err := store.TouchScheduleLastFired(t.Context(), "default", "recurring", boundary); err != nil {
		t.Fatalf("touch last fired: %v", err)
	}

	svc := &SchedulerService{}
	if _, err := svc.updateScheduleCore(t.Context(), store, &UpdateScheduleInput{Spec: &ScheduleSpec{
		Name:                  "recurring",
		GroupName:             "default",
		ScheduleExpression:    "rate(5 minutes)",
		Description:           "after",
		State:                 "ENABLED",
		ActionAfterCompletion: "NONE",
		Target: &schedulerstore.Target{
			Arn:     "arn:aws:lambda:us-east-1:000000000000:function:refire-test",
			RoleArn: "arn:aws:iam::000000000000:role/scheduler-test",
		},
	}}); err != nil {
		t.Fatalf("updateScheduleCore: %v", err)
	}

	updated, err := store.GetSchedule(t.Context(), "default", "recurring")
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if updated.LastFiredAt == nil {
		t.Fatal("delivered-boundary marker lost to an expression-keeping update")
	}
	if _, ok := (&Engine{}).dueBoundary(updated, now); ok {
		t.Error("expression-keeping update re-fired the delivered boundary")
	}
}
