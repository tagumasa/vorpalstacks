package scheduler

import (
	"context"
	"errors"
	"testing"
)

// TestMutateScheduleInLiveGroupRefusesDeletingGroup pins the update path's
// write-side guard: once a group is DELETING, MutateScheduleInLiveGroup
// refuses the mutation inside the same critical section that would write
// the record (the probe-to-write race the create path already forecloses),
// while the plain MutateSchedule keeps serving the engine's own markers on
// the doomed schedules until the cascade removes them.
func TestMutateScheduleInLiveGroupRefusesDeletingGroup(t *testing.T) {
	store := newCompletionTestStore(t)
	ctx := context.Background()

	if err := store.CreateScheduleGroup(ctx, &ScheduleGroup{Name: "doomed"}); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if err := store.CreateSchedule(ctx, &Schedule{
		Name:               "live-group-pin",
		GroupName:          "doomed",
		State:              ScheduleStateEnabled,
		ScheduleExpression: "at(2030-01-01T00:00:00)",
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}

	if err := store.MarkScheduleGroupDeleting(ctx, "doomed"); err != nil {
		t.Fatalf("mark deleting: %v", err)
	}

	err := store.MutateScheduleInLiveGroup(ctx, "doomed", "live-group-pin", func(*Schedule) error { return nil })
	if !errors.Is(err, ErrScheduleGroupDeleting) {
		t.Fatalf("MutateScheduleInLiveGroup in a DELETING group: error = %v, want ErrScheduleGroupDeleting", err)
	}

	// The engine's completion marker is not a user update: it must keep
	// working until the cascade deletes the schedule.
	if err := store.CompleteSchedule(ctx, "doomed", "live-group-pin"); err != nil {
		t.Fatalf("plain MutateSchedule marker in a DELETING group failed: %v", err)
	}
}

// TestMutateScheduleInLiveGroupReportsPurgedGroup pins the sentinel a fully
// completed cascade leaves behind: with the group record purged, the locked
// group check answers before the schedule read, so the live-group mutate
// surfaces the group's not-found sentinel even though the schedule is gone
// too — the core layer must map that sentinel to the group's not-found
// shape, never to an internal-server fault.
func TestMutateScheduleInLiveGroupReportsPurgedGroup(t *testing.T) {
	store := newCompletionTestStore(t)
	ctx := context.Background()

	if err := store.CreateScheduleGroup(ctx, &ScheduleGroup{Name: "purged"}); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if err := store.CreateSchedule(ctx, &Schedule{
		Name:               "purged-pin",
		GroupName:          "purged",
		State:              ScheduleStateEnabled,
		ScheduleExpression: "at(2030-01-01T00:00:00)",
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	for _, step := range []func() error{
		func() error { return store.MarkScheduleGroupDeleting(ctx, "purged") },
		func() error { return store.DeleteSchedulesInGroup(ctx, "purged") },
		func() error { return store.PurgeDeletedScheduleGroup(ctx, "purged") },
	} {
		if err := step(); err != nil {
			t.Fatalf("cascade step: %v", err)
		}
	}

	err := store.MutateScheduleInLiveGroup(ctx, "purged", "purged-pin", func(*Schedule) error { return nil })
	if !errors.Is(err, ErrScheduleGroupNotFound) {
		t.Fatalf("MutateScheduleInLiveGroup after the group's purge: error = %v, want ErrScheduleGroupNotFound", err)
	}
	// The plain mutate reads the schedule first, so the same absent pair
	// answers with the schedule's own sentinel there.
	if err := store.MutateSchedule(ctx, "purged", "purged-pin", func(*Schedule) error { return nil }); !errors.Is(err, ErrScheduleNotFound) {
		t.Fatalf("plain MutateSchedule after the group's purge: error = %v, want ErrScheduleNotFound", err)
	}
}
