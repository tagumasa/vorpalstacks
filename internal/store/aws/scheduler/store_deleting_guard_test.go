package scheduler

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// newDeletingGuardStore opens a throwaway store for the DELETING-guard pins.
func newDeletingGuardStore(t *testing.T) *SchedulerStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	store := NewSchedulerStore(st, "000000000000", "us-east-1")
	// The store's own Close stops the idempotency-token reaper goroutine;
	// the storage handle is released only after it.
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	return store
}

// TestCreateScheduleRefusedWhileGroupDeleting pins the relocated
// DELETING-refusal guard on the live schedule-creation path: a group that
// exists but is being deleted refuses new members, an absent explicit group
// refuses them, and the default group tolerates an unseeded record (it can
// never be deleted, so absence is not a refusal for it). The group name
// arrives resolved from the service layer — an empty name is rejected here
// rather than defaulted.
func TestCreateScheduleRefusedWhileGroupDeleting(t *testing.T) {
	store := newDeletingGuardStore(t)
	ctx := t.Context()

	if err := store.CreateScheduleGroup(ctx, &ScheduleGroup{Name: "dying"}); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if err := store.MarkScheduleGroupDeleting(ctx, "dying"); err != nil {
		t.Fatalf("mark deleting: %v", err)
	}

	err := store.CreateSchedule(ctx, &Schedule{
		Name:               "doomed",
		GroupName:          "dying",
		ScheduleExpression: "rate(1 hour)",
	})
	if !errors.Is(err, ErrScheduleGroupDeleting) {
		t.Errorf("CreateSchedule into DELETING group: error = %v, want ErrScheduleGroupDeleting", err)
	}
	if err := store.CreateSchedule(ctx, &Schedule{
		Name:               "orphan",
		GroupName:          "ghost",
		ScheduleExpression: "rate(1 hour)",
	}); !errors.Is(err, ErrScheduleGroupNotFound) {
		t.Errorf("CreateSchedule into absent group: error = %v, want ErrScheduleGroupNotFound", err)
	}
	// The default group's record this bare test store never seeded is not a
	// refusal: the group can never be deleted, so its absence cannot mean a
	// lifecycle conflict.
	if err := store.CreateSchedule(ctx, &Schedule{
		Name:               "unseeded-default",
		GroupName:          DefaultGroupName,
		ScheduleExpression: "rate(1 hour)",
	}); err != nil {
		t.Errorf("CreateSchedule into the unseeded default group: %v", err)
	}
	// An unresolved group name is a malformed identity, not a defaulted one.
	if err := store.CreateSchedule(ctx, &Schedule{
		Name:               "unresolved-group",
		ScheduleExpression: "rate(1 hour)",
	}); !errors.Is(err, ErrInvalidName) {
		t.Errorf("CreateSchedule with empty group: error = %v, want ErrInvalidName", err)
	}
}

// TestConcurrentCreateAndGroupDeletionUnderOneLock interleaves schedule
// creation with the full DELETING cascade (mark, member deletion, purge,
// group recreation) under the race detector. The single record lock must
// keep them serialised: every refusal names a live lifecycle condition, and
// once the storm settles the cascade converges — a final synchronous pass
// leaves the group and every member gone.
func TestConcurrentCreateAndGroupDeletionUnderOneLock(t *testing.T) {
	store := newDeletingGuardStore(t)
	ctx := t.Context()
	if err := store.CreateScheduleGroup(ctx, &ScheduleGroup{Name: "storm"}); err != nil {
		t.Fatalf("create group: %v", err)
	}

	const rounds = 50
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < rounds; i++ {
			_ = store.MarkScheduleGroupDeleting(ctx, "storm")
			_ = store.DeleteSchedulesInGroup(ctx, "storm")
			_ = store.PurgeDeletedScheduleGroup(ctx, "storm")
			_ = store.CreateScheduleGroup(ctx, &ScheduleGroup{Name: "storm"})
		}
	}()
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < rounds; i++ {
				err := store.CreateSchedule(ctx, &Schedule{
					Name:               fmt.Sprintf("storm-%d-%d", worker, i),
					GroupName:          "storm",
					State:              ScheduleStateEnabled,
					ScheduleExpression: "rate(1 hour)",
				})
				if err == nil {
					continue
				}
				if !errors.Is(err, ErrScheduleGroupDeleting) &&
					!errors.Is(err, ErrScheduleGroupNotFound) &&
					!errors.Is(err, ErrScheduleGroupAlreadyExists) &&
					!errors.Is(err, ErrScheduleAlreadyExists) {
					t.Errorf("CreateSchedule in the storm returned %v, want nil or a lifecycle sentinel", err)
				}
			}
		}(w)
	}
	wg.Wait()

	// Convergence: with the writers drained, one cascade pass must empty the
	// group completely — no acknowledged member can survive its group's
	// deletion. The survivors' enumeration mirrors the cascade's own walk —
	// the group key prefix with no state filter — because the engine's
	// enabled view drops non-enabled records, and a state-filtered listing
	// here could hide exactly the survivors this assertion exists to catch.
	if err := store.MarkScheduleGroupDeleting(ctx, "storm"); err != nil && !errors.Is(err, ErrScheduleGroupNotFound) {
		t.Fatalf("final mark: %v", err)
	}
	if err := store.DeleteSchedulesInGroup(ctx, "storm"); err != nil {
		t.Fatalf("final member deletion: %v", err)
	}
	if err := store.PurgeDeletedScheduleGroup(ctx, "storm"); err != nil {
		t.Fatalf("final purge: %v", err)
	}
	remaining, err := storecommon.ListMatching[Schedule](store.schedulesStore, "storm:", nil)
	if err != nil {
		t.Fatalf("list remaining: %v", err)
	}
	for _, sch := range remaining {
		t.Errorf("schedule %s survived its group's deletion", sch.Name)
	}
}
