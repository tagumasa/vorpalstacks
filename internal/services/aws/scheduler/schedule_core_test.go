package scheduler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
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
	store := schedulerstore.NewSchedulerStore(st, "123456789012", "us-east-1")
	// The store's own Close stops the idempotency-token reaper goroutine;
	// the storage handle is released only after it.
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	if err := store.EnsureDefaultGroup(context.Background()); err != nil {
		t.Fatalf("EnsureDefaultGroup failed: %v", err)
	}
	return store
}

// assertResourceNotFound pins the modelled ResourceNotFoundException: the
// error identity is the AWS error code (both transport planes map by code),
// and the message names the resource the request addressed.
func assertResourceNotFound(t *testing.T, err error, resource, identifier string) {
	t.Helper()
	if err == nil {
		t.Fatalf("want %s %q not-found, got nil", strings.ToLower(resource), identifier)
	}
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "ResourceNotFoundException" {
		t.Fatalf("error = %v, want ResourceNotFoundException", err)
	}
	if !strings.Contains(awsErr.Message, identifier) {
		t.Fatalf("not-found message %q does not name the identifier %q", awsErr.Message, identifier)
	}
}

// A replayed DeleteScheduleGroup ClientToken reports the first deletion's
// outcome even after the cascade has purged the group record — the plain
// DELETING-state idempotency cannot cover that window.
func TestDeleteScheduleGroupClientTokenReplayAfterPurge(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	svc := &SchedulerService{}

	if err := store.CreateScheduleGroup(ctx, &schedulerstore.ScheduleGroup{Name: "replay-group"}); err != nil {
		t.Fatalf("CreateScheduleGroup failed: %v", err)
	}
	if err := svc.deleteScheduleGroupCore(ctx, store, &DeleteScheduleGroupInput{Name: "replay-group", ClientToken: "group-delete-tok-1"}); err != nil {
		t.Fatalf("deleteScheduleGroupCore(token) failed: %v", err)
	}
	// Complete the cascade synchronously, as the engine sweep would.
	if err := store.PurgeDeletedScheduleGroup(ctx, "replay-group"); err != nil {
		t.Fatalf("PurgeDeletedScheduleGroup failed: %v", err)
	}
	// The replayed token reports the first deletion's success; without it
	// the purged group is not-found.
	if err := svc.deleteScheduleGroupCore(ctx, store, &DeleteScheduleGroupInput{Name: "replay-group", ClientToken: "group-delete-tok-1"}); err != nil {
		t.Fatalf("deleteScheduleGroupCore(replayed token) = %v, want the first deletion's success outcome", err)
	}
	err := svc.deleteScheduleGroupCore(ctx, store, &DeleteScheduleGroupInput{Name: "replay-group"})
	assertResourceNotFound(t, err, "Schedule group", "replay-group")
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
// ResourceNotFoundException naming the group, not an empty page.
func TestListSchedulesCoreGroupMissingIsResourceNotFound(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	svc := &SchedulerService{}

	_, err := svc.listSchedulesCore(ctx, store, &ListSchedulesInput{GroupName: "ghost"})
	assertResourceNotFound(t, err, "Schedule group", "ghost")
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
		t.Fatalf("validateTarget rejected an Input at the 256 KiB maximum: %v", err)
	}
}

// TestMapScheduleRecordWriteError pins the shared record-write mapping:
// every sentinel the locked writes can raise maps to its modelled
// exception, so a group purged mid-flight answers 404 on the update path
// exactly as it already does on the create path, and only genuine storage
// faults fall to the internal-server catch-all.
func TestMapScheduleRecordWriteError(t *testing.T) {
	cases := []struct {
		name       string
		in         error
		code       string
		httpStatus int
	}{
		{"duplicate name", schedulerstore.ErrScheduleAlreadyExists, "ConflictException", http.StatusConflict},
		{"group deleting", schedulerstore.ErrScheduleGroupDeleting, "ConflictException", http.StatusConflict},
		{"group purged mid-flight", schedulerstore.ErrScheduleGroupNotFound, "ResourceNotFoundException", http.StatusNotFound},
		{"schedule deleted mid-flight", schedulerstore.ErrScheduleNotFound, "ResourceNotFoundException", http.StatusNotFound},
		{"storage fault", errors.New("storage handle closed"), "InternalServerException", http.StatusInternalServerError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mapped := mapScheduleRecordWriteError(tc.in, "grp", "sch")
			awsErr, ok := mapped.(*awserrors.AWSError)
			if !ok {
				t.Fatalf("mapped error type = %T, want *awserrors.AWSError", mapped)
			}
			if awsErr.Code != tc.code {
				t.Errorf("code = %q, want %q", awsErr.Code, tc.code)
			}
			if awsErr.HTTPStatus != tc.httpStatus {
				t.Errorf("HTTP status = %d, want %d", awsErr.HTTPStatus, tc.httpStatus)
			}
		})
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
	// second call finds no schedule and reports not-found naming it.
	seedSchedule(t, store, "hourly", "default")
	if err := svc.deleteScheduleCore(ctx, store, &DeleteScheduleInput{Name: "hourly"}); err != nil {
		t.Fatalf("deleteScheduleCore() failed: %v", err)
	}
	err = svc.deleteScheduleCore(ctx, store, &DeleteScheduleInput{Name: "hourly"})
	assertResourceNotFound(t, err, "Schedule", "hourly")
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

// A group in DELETING refuses both schedule creation and schedule update:
// the engine cascade deletes the group's members within one sweep, so a
// request acknowledged now would be silently destroyed otherwise.
func TestScheduleCoresRejectDeletingGroup(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	svc := &SchedulerService{}

	if err := store.CreateScheduleGroup(ctx, &schedulerstore.ScheduleGroup{Name: "dying"}); err != nil {
		t.Fatalf("create group: %v", err)
	}
	seedSchedule(t, store, "existing", "dying")
	if err := store.MarkScheduleGroupDeleting(ctx, "dying"); err != nil {
		t.Fatalf("mark deleting: %v", err)
	}

	if _, err := svc.createScheduleCore(ctx, store, &CreateScheduleInput{
		Spec:   testTokenSpec("late-create", "dying"),
		Region: "us-east-1",
	}); err != ErrScheduleGroupDeleting {
		t.Fatalf("createScheduleCore into DELETING group: error = %v, want ErrScheduleGroupDeleting", err)
	}

	if _, err := svc.updateScheduleCore(ctx, store, &UpdateScheduleInput{
		Spec:   testTokenSpec("existing", "dying"),
		Region: "us-east-1",
	}); err != ErrScheduleGroupDeleting {
		t.Fatalf("updateScheduleCore in DELETING group: error = %v, want ErrScheduleGroupDeleting", err)
	}
}

// fakeRoleProvider fakes the two-method RolePolicyProvider the IAM validator
// consults: every role resolves to the provider's trust document.
type fakeRoleProvider struct{ trustDoc string }

func (f fakeRoleProvider) GetAssumeRolePolicyDocument(roleName string) (string, error) {
	return f.trustDoc, nil
}

func (f fakeRoleProvider) Exists(roleName string) bool { return true }

// Role validation failures must surface the scheduler model's error
// vocabulary — ValidationException naming the role — never the shared IAM
// validator's Lambda-flavoured default identity, whose
// InvalidParameterValueException/InvalidArn codes the scheduler model does
// not define. The pin drives a role that exists but trusts only Lambda, so
// the failure is the trust-policy refusal.
func TestScheduleCoresSurfaceValidationExceptionOnUntrustedRole(t *testing.T) {
	ctx := context.Background()
	store := testSchedulerCoreStore(t)
	svc := &SchedulerService{}
	lambdaTrustDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	validator := iam.NewIAMValidator(fakeRoleProvider{trustDoc: lambdaTrustDoc}, "123456789012")
	roleArn := "arn:aws:iam::123456789012:role/untrusted"

	assertRoleValidationError := func(t *testing.T, err error) {
		t.Helper()
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) {
			t.Fatalf("role failure is not an AWSError: %v", err)
		}
		if awsErr.Code != "ValidationException" {
			t.Fatalf("role failure code = %q, want ValidationException", awsErr.Code)
		}
		if !strings.Contains(awsErr.Message, roleArn) {
			t.Fatalf("role failure message %q does not name the role ARN", awsErr.Message)
		}
	}

	createSpec := validTarget()
	createSpec.RoleArn = roleArn
	_, err := svc.createScheduleCore(ctx, store, &CreateScheduleInput{
		Spec: &ScheduleSpec{
			Name:               "untrusted-role-create",
			ScheduleExpression: "rate(1 hour)",
			State:              "ENABLED",
			Target:             createSpec,
			FlexibleTimeWindow: validFTW(),
		},
		IAMValidator: validator,
	})
	assertRoleValidationError(t, err)

	seedSchedule(t, store, "existing", "default")
	updateSpec := validTarget()
	updateSpec.RoleArn = roleArn
	_, err = svc.updateScheduleCore(ctx, store, &UpdateScheduleInput{
		Spec: &ScheduleSpec{
			Name:               "existing",
			ScheduleExpression: "rate(1 hour)",
			State:              "ENABLED",
			Target:             updateSpec,
			FlexibleTimeWindow: validFTW(),
		},
		IAMValidator: validator,
	})
	assertRoleValidationError(t, err)
}

// TestAbsentCrossServiceInvokerFailsClosed pins the unified posture of the
// optional cross-service validators: a built engine whose bus carries no
// invoker for the service fails the creation request closed for both EC2
// (VPC configuration) and KMS (key existence), while a service without an
// engine — the construction window before BuildEngine, or store-only unit
// harnesses — skips the cross-service existence check (the modelled trait
// validation in validators.go still applies there; documented at both
// sites).
func TestAbsentCrossServiceInvokerFailsClosed(t *testing.T) {
	svc := newLifecycleService(t)
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	svc.engine.SetEventBus(bus)

	vpcTarget := &schedulerstore.Target{
		Arn: "arn:aws:ecs:us-east-1:000000000000:cluster/pin",
		EcsParameters: &schedulerstore.EcsParameters{
			TaskDefinitionArn: "arn:aws:ecs:us-east-1:000000000000:task-definition/family:1",
			NetworkConfiguration: &schedulerstore.NetworkConfiguration{
				AwsVpcConfiguration: &schedulerstore.AwsVpcConfiguration{
					Subnets: []string{"subnet-00000000000000001"},
				},
			},
		},
	}
	for name, err := range map[string]error{
		"absent EC2 invoker": svc.validateVpcConfig(context.Background(), "us-east-1", vpcTarget),
		"absent KMS invoker": svc.validateKmsKey(context.Background(), "us-east-1",
			"arn:aws:kms:us-east-1:000000000000:key/abcd-1234"),
	} {
		if err == nil {
			t.Errorf("%s: validation must fail closed", name)
			continue
		}
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.Code != "ValidationException" {
			t.Errorf("%s: error is %v, want ValidationException", name, err)
		}
	}

	bare := &SchedulerService{}
	if err := bare.validateVpcConfig(context.Background(), "us-east-1", vpcTarget); err != nil {
		t.Errorf("engine-less service must skip the VPC existence check, got %v", err)
	}
	if err := bare.validateKmsKey(context.Background(), "us-east-1",
		"arn:aws:kms:us-east-1:000000000000:key/abcd-1234"); err != nil {
		t.Errorf("engine-less service must skip the KMS key check, got %v", err)
	}
}
