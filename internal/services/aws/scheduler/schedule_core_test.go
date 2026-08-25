package scheduler

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// testSchedulerCoreStore opens an isolated store instance for Core-level
// behaviour tests; the default group is materialised exactly as the engine
// does at startup.
func testSchedulerCoreStore(t *testing.T) *schedulerstore.SchedulerStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("storage.Open failed: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	store := schedulerstore.NewSchedulerStore(st, "123456789012", "us-east-1")
	if err := store.EnsureDefaultGroup(context.Background()); err != nil {
		t.Fatalf("EnsureDefaultGroup failed: %v", err)
	}
	return store
}

func int32Ptr(v int32) *int32 { return &v }

// ListSchedules filter members carry documented constraints (State enum,
// NamePrefix pattern/length, NextToken length, GroupName pattern and
// existence) that the Core must enforce.
func TestListSchedulesCoreFilterValidation(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	svc := &SchedulerService{}

	tests := []struct {
		name    string
		in      *ListSchedulesInput
		wantErr bool
	}{
		{"no filters", &ListSchedulesInput{}, false},
		{"state filter enabled", &ListSchedulesInput{State: "ENABLED"}, false},
		{"state filter disabled", &ListSchedulesInput{State: "DISABLED"}, false},
		{"state filter invalid enum", &ListSchedulesInput{State: "BOGUS"}, true},
		{"name prefix valid", &ListSchedulesInput{NamePrefix: "prod-"}, false},
		{"name prefix invalid charset", &ListSchedulesInput{NamePrefix: "bad/name"}, true},
		{"name prefix too long", &ListSchedulesInput{NamePrefix: strings.Repeat("a", 65)}, true},
		{"next token within bound", &ListSchedulesInput{NextToken: strings.Repeat("t", 2048)}, false},
		{"next token beyond bound", &ListSchedulesInput{NextToken: strings.Repeat("t", 2049)}, true},
		{"group name invalid charset", &ListSchedulesInput{GroupName: "bad/name"}, true},
		{"existing group", &ListSchedulesInput{GroupName: "default"}, false},
		{"missing group", &ListSchedulesInput{GroupName: "ghost"}, true},
		{"max results explicit zero", &ListSchedulesInput{MaxResults: int32Ptr(0)}, true},
		{"max results lower bound", &ListSchedulesInput{MaxResults: int32Ptr(1)}, false},
		{"max results upper bound", &ListSchedulesInput{MaxResults: int32Ptr(100)}, false},
		{"max results above range", &ListSchedulesInput{MaxResults: int32Ptr(101)}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.listSchedulesCore(ctx, store, tc.in)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Fatalf("listSchedulesCore error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// A listing scoped to a group that does not exist reports the model's
// ResourceNotFoundException, not an empty page.
func TestListSchedulesCoreGroupMissingIsResourceNotFound(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	svc := &SchedulerService{}

	_, err := svc.listSchedulesCore(ctx, store, &ListSchedulesInput{GroupName: "ghost"})
	if err != ErrScheduleGroupNotFound {
		t.Fatalf("listSchedulesCore error = %v, want ErrScheduleGroupNotFound", err)
	}
}

// An invalid group-name charset is a validation failure, never a
// resource lookup.
func TestListSchedulesCoreGroupPatternIsValidationError(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	svc := &SchedulerService{}

	_, err := svc.listSchedulesCore(ctx, store, &ListSchedulesInput{GroupName: "bad/name"})
	if err != ErrValidation {
		t.Fatalf("listSchedulesCore error = %v, want ErrValidation", err)
	}
}

// seedSchedule stores a minimal valid schedule for identifier-resolution
// tests; the store layer performs no API validation.
func seedSchedule(t *testing.T, store *schedulerstore.SchedulerStore, name, group string) {
	t.Helper()
	err := store.CreateSchedule(context.Background(), &schedulerstore.Schedule{
		Name:               name,
		GroupName:          group,
		ScheduleExpression: "rate(1 hour)",
		State:              schedulerstore.ScheduleStateEnabled,
		Target:             validTarget(),
	})
	if err != nil {
		t.Fatalf("seed CreateSchedule(%s/%s) failed: %v", group, name, err)
	}
}

// The schedule identifier pair (Name pattern, GroupName default resolution
// and pattern) is validated by the Core on every read/update/delete path.
func TestScheduleCoreIdentifierValidation(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	seedSchedule(t, store, "nightly", "default")
	svc := &SchedulerService{}

	tests := []struct {
		name    string
		get     *GetScheduleInput
		del     *DeleteScheduleInput
		wantErr error
	}{
		{name: "get invalid name charset", get: &GetScheduleInput{Name: "bad/name"}, wantErr: ErrValidation},
		{name: "get missing name", get: &GetScheduleInput{}, wantErr: ErrValidation},
		{name: "get invalid group charset", get: &GetScheduleInput{Name: "nightly", GroupName: "bad/name"}, wantErr: ErrValidation},
		{name: "get default group resolution", get: &GetScheduleInput{Name: "nightly"}, wantErr: nil},
		{name: "delete invalid name charset", del: &DeleteScheduleInput{Name: "bad/name"}, wantErr: ErrValidation},
		{name: "delete invalid group charset", del: &DeleteScheduleInput{Name: "nightly", GroupName: "bad/name"}, wantErr: ErrValidation},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			if tc.get != nil {
				_, err = svc.getScheduleCore(ctx, store, tc.get)
			} else {
				err = svc.deleteScheduleCore(ctx, store, tc.del)
			}
			if err != tc.wantErr {
				t.Fatalf("core error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// RoleArn must reference an IAM role (Target shape pattern) and Target.Input
// carries a documented 256 KiB maximum.
func TestValidateTargetRoleArnShapeAndInputSize(t *testing.T) {
	nonIam := validTarget()
	nonIam.RoleArn = "arn:aws:sqs:us-east-1:123456789012:queue"
	if err := validateTarget(nonIam); err == nil {
		t.Fatal("validateTarget accepted a non-IAM RoleArn")
	}

	nonRole := validTarget()
	nonRole.RoleArn = "arn:aws:iam::123456789012:policy/p"
	if err := validateTarget(nonRole); err == nil {
		t.Fatal("validateTarget accepted an IAM non-role RoleArn")
	}

	bigInput := validTarget()
	bigInput.Input = strings.Repeat("a", 262145)
	if err := validateTarget(bigInput); err == nil {
		t.Fatal("validateTarget accepted an Input beyond the 256 KiB maximum")
	}

	maxInput := validTarget()
	maxInput.Input = strings.Repeat("a", 262144)
	if err := validateTarget(maxInput); err != nil {
		t.Fatalf("validateTarget(rejected an Input at the 256 KiB maximum: %v", err)
	}
}

// parseTarget is a pure wire-to-DTO transform: required-member rejection
// (Arn, RoleArn) belongs to the Core validateTarget path alone.
func TestParseTargetIsPureParsing(t *testing.T) {
	target, err := parseTarget(map[string]interface{}{
		"Target": map[string]interface{}{"RoleArn": "arn:aws:iam::123456789012:role/r"},
	})
	if err != nil {
		t.Fatalf("parseTarget(missing Arn) error = %v, want pure parsing with no validation", err)
	}
	if target == nil {
		t.Fatal("parseTarget(missing Arn) returned a nil target")
	}
	// The Core validation path still rejects the required-member absence.
	if err := validateTarget(target); err != ErrInvalidTarget {
		t.Fatalf("validateTarget(missing Arn) error = %v, want ErrInvalidTarget", err)
	}

	absent, err := parseTarget(map[string]interface{}{})
	if err != nil || absent != nil {
		t.Fatalf("parseTarget(no Target member) = (%v, %v), want (nil, nil)", absent, err)
	}
}

// UpdateSchedule and DeleteSchedule honour the ClientToken idempotency
// member: an invalid token is rejected, and a replayed token returns the
// first application's outcome instead of re-executing.
func TestScheduleClientTokenValidationAndReplay(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	seedSchedule(t, store, "nightly", "default")
	svc := &SchedulerService{}

	validSpec := func() *ScheduleSpec {
		return &ScheduleSpec{
			Name:               "nightly",
			ScheduleExpression: "rate(2 hours)",
			State:              "ENABLED",
			Target:             validTarget(),
			FlexibleTimeWindow: validFTW(),
		}
	}

	if _, err := svc.updateScheduleCore(ctx, store, &UpdateScheduleInput{
		Spec:        validSpec(),
		ClientToken: "bad token!",
	}); err == nil {
		t.Fatal("updateScheduleCore accepted an invalid ClientToken pattern")
	}

	first, err := svc.updateScheduleCore(ctx, store, &UpdateScheduleInput{Spec: validSpec(), ClientToken: "update-tok-1"})
	if err != nil {
		t.Fatalf("updateScheduleCore(token) failed: %v", err)
	}
	replay, err := svc.updateScheduleCore(ctx, store, &UpdateScheduleInput{Spec: validSpec(), ClientToken: "update-tok-1"})
	if err != nil {
		t.Fatalf("updateScheduleCore(replayed token) failed: %v", err)
	}
	if replay.ScheduleArn != first.ScheduleArn {
		t.Fatalf("replayed update returned %q, want the first application's ARN %q", replay.ScheduleArn, first.ScheduleArn)
	}

	seedSchedule(t, store, "daily", "default")
	if err := svc.deleteScheduleCore(ctx, store, &DeleteScheduleInput{Name: "daily", ClientToken: "delete-tok-1"}); err != nil {
		t.Fatalf("deleteScheduleCore(token) failed: %v", err)
	}
	if err := svc.deleteScheduleCore(ctx, store, &DeleteScheduleInput{Name: "daily", ClientToken: "delete-tok-1"}); err != nil {
		t.Fatalf("deleteScheduleCore(replayed token) = %v, want the first deletion's success outcome", err)
	}

	// Without a token the deletion is a plain idempotent-trait delete: the
	// second call finds no schedule and reports not-found.
	seedSchedule(t, store, "hourly", "default")
	if err := svc.deleteScheduleCore(ctx, store, &DeleteScheduleInput{Name: "hourly"}); err != nil {
		t.Fatalf("deleteScheduleCore() failed: %v", err)
	}
	if err := svc.deleteScheduleCore(ctx, store, &DeleteScheduleInput{Name: "hourly"}); err != ErrScheduleNotFound {
		t.Fatalf("deleteScheduleCore(second, no token) = %v, want ErrScheduleNotFound", err)
	}
}

// UpdateSchedule validates the identifier pair before the existence probe:
// a malformed identifier is a ValidationException, not a resource lookup.
func TestUpdateScheduleCoreIdentifierValidation(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	seedSchedule(t, store, "nightly", "default")
	svc := &SchedulerService{}

	if _, err := svc.updateScheduleCore(ctx, store, &UpdateScheduleInput{
		Spec: &ScheduleSpec{Name: "bad/name", ScheduleExpression: "rate(1 hour)", Target: validTarget(), FlexibleTimeWindow: validFTW()},
	}); err != ErrValidation {
		t.Fatalf("updateScheduleCore(invalid name) error = %v, want ErrValidation", err)
	}

	result, err := svc.updateScheduleCore(ctx, store, &UpdateScheduleInput{
		Spec: &ScheduleSpec{Name: "nightly", ScheduleExpression: "rate(2 hours)", State: "ENABLED", Target: validTarget(), FlexibleTimeWindow: validFTW()},
	})
	if err != nil {
		t.Fatalf("updateScheduleCore(default group) failed: %v", err)
	}
	if result.ScheduleArn == "" {
		t.Fatal("updateScheduleCore returned an empty ScheduleArn")
	}
}
