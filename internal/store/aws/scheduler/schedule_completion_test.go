package scheduler

import (
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newCompletionTestStore(t *testing.T) *SchedulerStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewSchedulerStore(st, "000000000000", "us-east-1")
}

// TestCompleteScheduleExcludesFromEnabled pins the completion contract of
// one-time schedules: a completed schedule keeps its wire state (the AWS
// ScheduleState enum has no COMPLETED value) but disappears from the
// engine's enabled sweep, and completing twice is idempotent.
func TestCompleteScheduleExcludesFromEnabled(t *testing.T) {
	store := newCompletionTestStore(t)
	oneTime := &Schedule{
		Name:               "one-time",
		GroupName:          "default",
		State:              ScheduleStateEnabled,
		ScheduleExpression: "at(2027-01-01T00:00:00)",
	}
	if err := store.CreateSchedule(t.Context(), oneTime); err != nil {
		t.Fatalf("create schedule: %v", err)
	}

	// A fresh sweep sees the enabled one-time schedule.
	enabled, err := store.GetAllEnabledSchedules(t.Context())
	if err != nil {
		t.Fatalf("list enabled: %v", err)
	}
	if len(enabled) != 1 || enabled[0].Name != "one-time" {
		t.Fatalf("enabled sweep = %d schedules, want the one-time schedule", len(enabled))
	}

	// Completing twice is idempotent.
	for i := 0; i < 2; i++ {
		if err := store.CompleteSchedule(t.Context(), "default", "one-time"); err != nil {
			t.Fatalf("complete schedule (call %d): %v", i+1, err)
		}
	}

	// The wire state is preserved while the completion marker is set.
	got, err := store.GetSchedule(t.Context(), "default", "one-time")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	if got.State != ScheduleStateEnabled {
		t.Errorf("wire state = %q, want ENABLED preserved", got.State)
	}
	if got.CompletionDate == nil {
		t.Error("completion marker not persisted")
	}

	// The completed schedule no longer fires.
	enabled, err = store.GetAllEnabledSchedules(t.Context())
	if err != nil {
		t.Fatalf("list enabled after completion: %v", err)
	}
	if len(enabled) != 0 {
		t.Errorf("enabled sweep after completion = %d schedules, want 0", len(enabled))
	}

	// Completing an unknown schedule reports the not-found error.
	if err := store.CompleteSchedule(t.Context(), "default", "missing"); err == nil {
		t.Error("completing an unknown schedule did not fail")
	}
}

// TestCompleteScheduleConcurrentWithUpdate pins the record-mutation
// serialisation against the audited lost-update interleaving: a completion
// and a user field update racing on the same record must both survive,
// because each runs its whole read-modify-write cycle inside the record
// lock — never a stale pre-read record.
func TestCompleteScheduleConcurrentWithUpdate(t *testing.T) {
	store := newCompletionTestStore(t)
	oneTime := &Schedule{
		Name:               "one-time",
		GroupName:          "default",
		State:              ScheduleStateEnabled,
		ScheduleExpression: "at(2027-01-01T00:00:00)",
		Description:        "before",
	}
	if err := store.CreateSchedule(t.Context(), oneTime); err != nil {
		t.Fatalf("create schedule: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// The service-side updater pattern: mutate the user's fields through
	// the store-level atomic mutation.
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := store.MutateSchedule(t.Context(), "default", "one-time", func(current *Schedule) error {
				current.Description = "after"
				return nil
			}); err != nil {
				t.Errorf("mutate schedule: %v", err)
				return
			}
		}
	}()
	// The engine-side completion pattern, racing from another goroutine.
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := store.CompleteSchedule(t.Context(), "default", "one-time"); err != nil {
				t.Errorf("complete schedule: %v", err)
				return
			}
		}
	}()
	wg.Wait()

	got, err := store.GetSchedule(t.Context(), "default", "one-time")
	if err != nil {
		t.Fatalf("final get: %v", err)
	}
	if got.Description != "after" {
		t.Errorf("user fields lost to the completion write: description = %q, want %q", got.Description, "after")
	}
	if got.CompletionDate == nil {
		t.Error("completion marker lost to the concurrent update write")
	}
}

// TestTouchScheduleLastFiredRoundTrip pins the delivered-boundary marker
// contract: the marker only advances and the write never stamps
// LastModificationDate.
func TestTouchScheduleLastFiredRoundTrip(t *testing.T) {
	store := newCompletionTestStore(t)
	if err := store.CreateSchedule(t.Context(), &Schedule{
		Name:               "recurring",
		GroupName:          "default",
		State:              ScheduleStateEnabled,
		ScheduleExpression: "rate(5 minutes)",
		CreationDate:       time.Now().UTC(),
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	created, err := store.GetSchedule(t.Context(), "default", "recurring")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	modBefore := created.LastModificationDate

	newer := time.Now().UTC().Add(-5 * time.Minute)
	older := newer.Add(-5 * time.Minute)
	newest := time.Now().UTC().Add(-1 * time.Minute)

	if err := store.TouchScheduleLastFired(t.Context(), "default", "recurring", older); err != nil {
		t.Fatalf("touch older boundary: %v", err)
	}
	// A stale boundary never regresses the marker.
	if err := store.TouchScheduleLastFired(t.Context(), "default", "recurring", older.Add(-time.Hour)); err != nil {
		t.Fatalf("touch stale boundary: %v", err)
	}
	got, err := store.GetSchedule(t.Context(), "default", "recurring")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	if got.LastFiredAt == nil || !got.LastFiredAt.Equal(older) {
		t.Fatalf("marker = %v, want %v", got.LastFiredAt, older)
	}

	if err := store.TouchScheduleLastFired(t.Context(), "default", "recurring", newer); err != nil {
		t.Fatalf("touch newer boundary: %v", err)
	}
	if err := store.TouchScheduleLastFired(t.Context(), "default", "recurring", newest); err != nil {
		t.Fatalf("touch newest boundary: %v", err)
	}
	got, err = store.GetSchedule(t.Context(), "default", "recurring")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	if got.LastFiredAt == nil || !got.LastFiredAt.Equal(newest) {
		t.Fatalf("marker = %v, want %v", got.LastFiredAt, newest)
	}
	if !got.LastModificationDate.Equal(modBefore) {
		t.Error("touch stamped LastModificationDate")
	}
}

// TestDeleteScheduleNotResurrectedByMutate pins the delete serialisation:
// a record mutation racing the delete can never write the record back
// after the delete removed it.
func TestDeleteScheduleNotResurrectedByMutate(t *testing.T) {
	store := newCompletionTestStore(t)

	for i := 0; i < 50; i++ {
		if err := store.CreateSchedule(t.Context(), &Schedule{
			Name:               "victim",
			GroupName:          "default",
			State:              ScheduleStateEnabled,
			ScheduleExpression: "rate(5 minutes)",
		}); err != nil {
			t.Fatalf("create schedule (iteration %d): %v", i, err)
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			// A not-found error is a legitimate outcome: the delete may
			// have landed first.
			_ = store.MutateSchedule(t.Context(), "default", "victim", func(sch *Schedule) error {
				sch.Description = "after"
				return nil
			})
		}()
		go func() {
			defer wg.Done()
			_ = store.DeleteSchedule(t.Context(), "default", "victim")
		}()
		wg.Wait()

		if _, err := store.GetSchedule(t.Context(), "default", "victim"); err != ErrScheduleNotFound {
			t.Fatalf("iteration %d: schedule resurrected after delete (err = %v)", i, err)
		}
	}
}

// TestDeleteScheduleGroupNotResurrectedByUpdate pins the same contract on
// group records: a group update racing the group delete can never write
// the group back after the delete removed it.
func TestDeleteScheduleGroupNotResurrectedByUpdate(t *testing.T) {
	store := newCompletionTestStore(t)

	for i := 0; i < 50; i++ {
		if err := store.CreateScheduleGroup(t.Context(), &ScheduleGroup{Name: "victim"}); err != nil {
			t.Fatalf("create schedule group (iteration %d): %v", i, err)
		}
		group, err := store.GetScheduleGroup(t.Context(), "victim")
		if err != nil {
			t.Fatalf("get schedule group (iteration %d): %v", i, err)
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			// A not-found error is a legitimate outcome: the delete may
			// have landed first.
			_ = store.UpdateScheduleGroup(t.Context(), group)
		}()
		go func() {
			defer wg.Done()
			_ = store.MarkScheduleGroupDeleting(t.Context(), "victim")
		}()
		wg.Wait()

		// The delete path marks the group DELETING and purges it later;
		// either the group is already gone or it must carry the DELETING
		// state — a concurrent update can never resurrect it as live.
		group, err = store.GetScheduleGroup(t.Context(), "victim")
		if err == nil && group.State != ScheduleGroupStateDeleting {
			t.Fatalf("iteration %d: schedule group resurrected as %s after delete", i, group.State)
		}
		if err != nil && err != ErrScheduleGroupNotFound {
			t.Fatalf("iteration %d: get after delete: %v", i, err)
		}
		// Complete the cascade the engine sweep would perform so the
		// next iteration creates a fresh group.
		_ = store.PurgeDeletedScheduleGroup(t.Context(), "victim")
	}
}
