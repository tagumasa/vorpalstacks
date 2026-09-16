package scheduler

import (
	"testing"

	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// The ClientToken protocol pins: every core routes its claim/release through
// claimClientToken, so a released claim must let a retry execute (a stale
// claim would replay the first application's ARN without performing the
// operation), and a replay must return the first outcome without executing.
// The discriminator is the store state after the call: a replay creates and
// deletes nothing.

func testTokenSpec(name, group string) *ScheduleSpec {
	return &ScheduleSpec{
		Name:                  name,
		GroupName:             group,
		ScheduleExpression:    "at(2030-01-01T00:00:00)",
		State:                 "ENABLED",
		ActionAfterCompletion: "NONE",
		Target: &schedulerstore.Target{
			Arn:     "arn:aws:lambda:us-east-1:000000000000:function:token-test",
			RoleArn: "arn:aws:iam::000000000000:role/scheduler-test",
		},
		FlexibleTimeWindow: &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff},
	}
}

// TestCreateScheduleClientTokenReleasedOnErrorPath pins create's release
// site: a create that fails after the claim (missing group) releases the
// token, so the retry with the same token executes instead of replaying.
func TestCreateScheduleClientTokenReleasedOnErrorPath(t *testing.T) {
	store := newFiringStore(t)
	svc := &SchedulerService{}

	if _, err := svc.createScheduleCore(t.Context(), store, &CreateScheduleInput{
		Spec:        testTokenSpec("tok-release", "missing-group"),
		ClientToken: "create-release-1",
		Region:      "us-east-1",
	}); err == nil {
		t.Fatal("create into a missing group was accepted")
	}

	// The group appears; the SAME token must now execute the creation. A
	// stale claim would replay the first attempt's ARN and leave the store
	// without the schedule.
	if err := store.CreateScheduleGroup(t.Context(), &schedulerstore.ScheduleGroup{Name: "missing-group"}); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if _, err := svc.createScheduleCore(t.Context(), store, &CreateScheduleInput{
		Spec:        testTokenSpec("tok-release", "missing-group"),
		ClientToken: "create-release-1",
		Region:      "us-east-1",
	}); err != nil {
		t.Fatalf("retry after release failed: %v", err)
	}
	if _, err := store.GetSchedule(t.Context(), "missing-group", "tok-release"); err != nil {
		t.Errorf("token was not released on the error path — retry replayed instead of creating: %v", err)
	}
}

// TestDeleteScheduleClientTokenReleasedOnErrorPath pins delete's release
// site: a delete that fails after the claim (missing schedule) releases the
// token, so the retry with the same token executes the deletion.
func TestDeleteScheduleClientTokenReleasedOnErrorPath(t *testing.T) {
	store := newFiringStore(t)
	svc := &SchedulerService{}

	if err := svc.deleteScheduleCore(t.Context(), store, &DeleteScheduleInput{
		Name:        "tok-del",
		GroupName:   "default",
		ClientToken: "delete-release-1",
	}); err == nil {
		t.Fatal("delete of a missing schedule was accepted")
	}

	// The schedule appears; the SAME token must now execute the deletion. A
	// stale claim would replay the first attempt's success and leave the
	// schedule in place.
	if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
		Name:               "tok-del",
		GroupName:          "default",
		State:              schedulerstore.ScheduleStateEnabled,
		ScheduleExpression: "at(2030-01-01T00:00:00)",
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	if err := svc.deleteScheduleCore(t.Context(), store, &DeleteScheduleInput{
		Name:        "tok-del",
		GroupName:   "default",
		ClientToken: "delete-release-1",
	}); err != nil {
		t.Fatalf("retry after release failed: %v", err)
	}
	if _, err := store.GetSchedule(t.Context(), "default", "tok-del"); err == nil {
		t.Error("token was not released on the error path — retry replayed instead of deleting")
	}
}

// TestCreateScheduleGroupClientTokenReplayAndRelease pins the group core's
// both sides: a replay returns the first ARN without re-executing (a second
// create with the same token succeeds instead of surfacing AlreadyExists),
// and the AlreadyExists failure path releases the challenger's token so its
// retry executes once the conflict is gone.
func TestCreateScheduleGroupClientTokenReplayAndRelease(t *testing.T) {
	store := newFiringStore(t)
	svc := &SchedulerService{}

	first, err := svc.createScheduleGroupCore(t.Context(), store, &CreateScheduleGroupInput{
		Name:        "tokgrp",
		ClientToken: "group-token-a",
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	// Replay: same token returns the first ARN without executing — without
	// the replay the group's existence would surface AlreadyExists.
	replayed, err := svc.createScheduleGroupCore(t.Context(), store, &CreateScheduleGroupInput{
		Name:        "tokgrp",
		ClientToken: "group-token-a",
	})
	if err != nil {
		t.Fatalf("replay failed: %v", err)
	}
	if replayed.ScheduleGroupArn != first.ScheduleGroupArn {
		t.Errorf("replay returned %q, want the first application's ARN %q", replayed.ScheduleGroupArn, first.ScheduleGroupArn)
	}

	// A different token against the existing group fails AlreadyExists and
	// must be released; once the conflict is purged, the retry executes.
	if _, err := svc.createScheduleGroupCore(t.Context(), store, &CreateScheduleGroupInput{
		Name:        "tokgrp",
		ClientToken: "group-token-b",
	}); err == nil {
		t.Fatal("duplicate group create was accepted")
	}
	if err := store.MarkScheduleGroupDeleting(t.Context(), "tokgrp"); err != nil {
		t.Fatalf("mark deleting: %v", err)
	}
	if err := store.PurgeDeletedScheduleGroup(t.Context(), "tokgrp"); err != nil {
		t.Fatalf("purge: %v", err)
	}
	if _, err := svc.createScheduleGroupCore(t.Context(), store, &CreateScheduleGroupInput{
		Name:        "tokgrp",
		ClientToken: "group-token-b",
	}); err != nil {
		t.Fatalf("retry after release failed: %v", err)
	}
	if _, err := store.GetScheduleGroup(t.Context(), "tokgrp"); err != nil {
		t.Errorf("token was not released on the error path — retry replayed instead of creating: %v", err)
	}
}
