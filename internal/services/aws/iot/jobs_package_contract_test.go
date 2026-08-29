package iot

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// The AWS SDK validates required members client-side, so several of the
// jobs/package-family required-member negatives never reach the server
// through the typed SDK. The guards below run before any store access, so a
// nil store stands in for the real one.

func TestCreateJobTemplateRequiresDescription(t *testing.T) {
	svc := &IoTService{}
	_, err := svc.createJobTemplateCore(nil, CreateJobTemplateInput{
		JobTemplateID: "tmpl-1",
		Document:      `{"version":1}`,
	})
	if !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for a missing description, got %v", err)
	}
}

func TestCreateOTAUpdateRequiredMembers(t *testing.T) {
	svc := &IoTService{}
	cases := []struct {
		name string
		in   CreateOTAUpdateInput
	}{
		{"missing targets", CreateOTAUpdateInput{
			OtaUpdateID: "ota-1",
			Files:       []interface{}{map[string]interface{}{"fileName": "fw.bin"}},
			RoleArn:     "arn:aws:iam::123456789012:role/ota",
		}},
		{"missing files", CreateOTAUpdateInput{
			OtaUpdateID: "ota-1",
			Targets:     []string{"arn:aws:iot:us-east-1:123456789012:thing/dev"},
			RoleArn:     "arn:aws:iam::123456789012:role/ota",
		}},
		{"missing roleArn", CreateOTAUpdateInput{
			OtaUpdateID: "ota-1",
			Targets:     []string{"arn:aws:iot:us-east-1:123456789012:thing/dev"},
			Files:       []interface{}{map[string]interface{}{"fileName": "fw.bin"}},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := svc.createOTAUpdateCore(nil, tc.in); !errors.Is(err, iotstore.ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestDeleteJobExecutionRequiresNumber(t *testing.T) {
	svc := &IoTService{}
	err := svc.deleteJobExecutionCore(nil, DeleteJobExecutionInput{
		JobID:     "job-1",
		ThingName: "thing-1",
	})
	if !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for a missing executionNumber, got %v", err)
	}
}

func TestListSbomValidationResultsRequiresPathLabels(t *testing.T) {
	svc := &IoTService{}
	if _, err := svc.listSbomValidationResultsCore("pkg", ""); !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for a missing versionName, got %v", err)
	}
	if _, err := svc.listSbomValidationResultsCore("", "1.0"); !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for a missing packageName, got %v", err)
	}
}

func TestUpdatePackageSetAndUnsetRejected(t *testing.T) {
	svc := &IoTService{}
	err := svc.updatePackageCore(nil, UpdatePackageInput{
		PackageName:         "pkg",
		DefaultVersionName:  "1.0",
		UnsetDefaultVersion: true,
	})
	if !errors.Is(err, iotstore.ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for set+unset, got %v", err)
	}
}

func TestParseTimeFilter(t *testing.T) {
	after, before, err := parseTimeFilter(map[string]interface{}{
		"after":  "2026-01-01T00:00:00Z",
		"before": "2030-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("expected a valid filter to parse, got %v", err)
	}
	if after == 0 || before == 0 {
		t.Fatalf("expected non-zero bounds, got after=%d before=%d", after, before)
	}
	if _, _, err := parseTimeFilter(nil); err != nil {
		t.Fatalf("expected an absent filter to pass, got %v", err)
	}
	if _, _, err := parseTimeFilter(map[string]interface{}{
		"after": "not-a-timestamp",
	}); !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for a malformed timestamp, got %v", err)
	}
}

func TestValidateCommandSortOrder(t *testing.T) {
	for _, valid := range []string{"", "ASCENDING", "DESCENDING"} {
		if err := validateCommandSortOrder(valid); err != nil {
			t.Fatalf("expected %q to be accepted, got %v", valid, err)
		}
	}
	if err := validateCommandSortOrder("nope"); !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for an off-enum sort order, got %v", err)
	}
}

func TestCommandDeclaresParameter(t *testing.T) {
	parameters := []interface{}{
		map[string]interface{}{"name": "p1"},
		map[string]interface{}{"name": "p2"},
	}
	if !commandDeclaresParameter(parameters, "p2") {
		t.Fatal("expected p2 to be declared")
	}
	if commandDeclaresParameter(parameters, "p3") {
		t.Fatal("expected p3 to be undeclared")
	}
	if commandDeclaresParameter(nil, "p1") {
		t.Fatal("expected no declaration on an absent parameter list")
	}
}

func TestRecordTimeWithin(t *testing.T) {
	if !recordTimeWithin(int64(100), 0, 0) {
		t.Fatal("absent bounds match everything")
	}
	if !recordTimeWithin(int64(100), int64(50), int64(150)) {
		t.Fatal("100 is within [50,150]")
	}
	if recordTimeWithin(int64(200), int64(50), int64(150)) {
		t.Fatal("200 is outside [50,150]")
	}
}

// TestCreateJobTemplateRequiresDocumentCarrier pins the documented
// conditional requirement: the job document is required unless
// documentSource or jobArn supplies it, so a request carrying none of the
// three carriers is rejected before any store access (nil store stands in).
func TestCreateJobTemplateRequiresDocumentCarrier(t *testing.T) {
	svc := &IoTService{}
	_, err := svc.createJobTemplateCore(nil, CreateJobTemplateInput{
		JobTemplateID: "tmpl-1",
		Description:   "documented description",
	})
	if !errors.Is(err, iotstore.ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for no document carrier, got %v", err)
	}
}

// newJobsExecutionTestStore opens a temporary store for the job-execution
// state-machine pins.
func newJobsExecutionTestStore(t *testing.T) iotstore.IotStoreInterface {
	t.Helper()
	rawStore, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rawStore.Close() })
	return iotstore.NewIotStore(rawStore, "000000000000", "us-east-1", nil)
}

// seedJobExecution creates a job targeting one thing (materialising the
// QUEUED execution) and then forces the execution record to the given
// status.
func seedJobExecution(t *testing.T, store iotstore.IotStoreInterface, jobID, thingName, status string) {
	t.Helper()
	if _, err := store.CreateJob(&iotstore.Job{
		JobID:   jobID,
		Targets: []string{"arn:aws:iot:us-east-1:000000000000:thing/" + thingName},
	}); err != nil {
		t.Fatalf("create job %s: %v", jobID, err)
	}
	key := jobExecutionRecordKey(jobID, thingName)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		t.Fatalf("load execution: %v", err)
	}
	if !exists {
		t.Fatalf("expected a materialised execution for job %s", jobID)
	}
	rec["status"] = status
	if err := store.PutGeneric(key, rec); err != nil {
		t.Fatalf("seed status %s: %v", status, err)
	}
}

// executionStatus reads the stored status of one job execution record.
func executionStatus(t *testing.T, store iotstore.IotStoreInterface, jobID, thingName string) (string, bool) {
	t.Helper()
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(jobExecutionRecordKey(jobID, thingName), &rec)
	if err != nil {
		t.Fatalf("load execution: %v", err)
	}
	if !exists {
		return "", false
	}
	status, _ := rec["status"].(string)
	return status, true
}

// TestCancelJobExecutionStateGate pins the documented cancel state
// machine: an execution can be canceled while QUEUED, or while IN_PROGRESS
// only with force; every other state is an invalid transition.
func TestCancelJobExecutionStateGate(t *testing.T) {
	svc := &IoTService{}
	cases := []struct {
		name      string
		status    string
		force     bool
		wantError error
	}{
		{"queued without force", "QUEUED", false, nil},
		{"in progress without force", "IN_PROGRESS", false, iotstore.ErrInvalidStateTransition},
		{"in progress with force", "IN_PROGRESS", true, nil},
		{"terminal even with force", "SUCCEEDED", true, iotstore.ErrInvalidStateTransition},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := newJobsExecutionTestStore(t)
			jobID, thing := "job-"+tc.name, "thing-"+tc.name
			seedJobExecution(t, store, jobID, thing, tc.status)
			err := svc.cancelJobExecutionCore(store, CancelJobExecutionInput{
				JobID: jobID, ThingName: thing, Force: tc.force,
			})
			if !errors.Is(err, tc.wantError) {
				t.Fatalf("cancel %s force=%v: got %v, want %v", tc.status, tc.force, err, tc.wantError)
			}
			status, exists := executionStatus(t, store, jobID, thing)
			if !exists {
				t.Fatal("execution record vanished")
			}
			if tc.wantError == nil && status != "CANCELED" {
				t.Fatalf("expected CANCELED, got %s", status)
			}
			if tc.wantError != nil && status != tc.status {
				t.Fatalf("rejected cancel must preserve %s, got %s", tc.status, status)
			}
		})
	}
}

// TestDeleteJobExecutionForceGate pins the documented delete rule: without
// force only terminal executions may be deleted; force lifts the
// restriction for non-terminal states.
func TestDeleteJobExecutionForceGate(t *testing.T) {
	svc := &IoTService{}
	t.Run("queued without force is rejected", func(t *testing.T) {
		store := newJobsExecutionTestStore(t)
		seedJobExecution(t, store, "job-queued", "thing-queued", "QUEUED")
		err := svc.deleteJobExecutionCore(store, DeleteJobExecutionInput{
			JobID: "job-queued", ThingName: "thing-queued", ExecutionNumber: 1,
		})
		if !errors.Is(err, iotstore.ErrInvalidStateTransition) {
			t.Fatalf("expected ErrInvalidStateTransition, got %v", err)
		}
		if _, exists := executionStatus(t, store, "job-queued", "thing-queued"); !exists {
			t.Fatal("rejected delete must preserve the execution record")
		}
	})
	t.Run("queued with force is deleted", func(t *testing.T) {
		store := newJobsExecutionTestStore(t)
		seedJobExecution(t, store, "job-forced", "thing-forced", "QUEUED")
		err := svc.deleteJobExecutionCore(store, DeleteJobExecutionInput{
			JobID: "job-forced", ThingName: "thing-forced", ExecutionNumber: 1, Force: true,
		})
		if err != nil {
			t.Fatalf("forced delete: %v", err)
		}
		if _, exists := executionStatus(t, store, "job-forced", "thing-forced"); exists {
			t.Fatal("expected the execution record to be deleted")
		}
	})
	t.Run("canceled terminal deletes without force", func(t *testing.T) {
		store := newJobsExecutionTestStore(t)
		seedJobExecution(t, store, "job-canceled", "thing-canceled", "CANCELED")
		err := svc.deleteJobExecutionCore(store, DeleteJobExecutionInput{
			JobID: "job-canceled", ThingName: "thing-canceled", ExecutionNumber: 1,
		})
		if err != nil {
			t.Fatalf("terminal delete: %v", err)
		}
	})
}

// TestCancelJobExecutionExplicitZeroVersion pins that an explicitly
// supplied expectedVersion of zero is a version conflict, not an alias for
// an omitted member: stored versions start at one.
func TestCancelJobExecutionExplicitZeroVersion(t *testing.T) {
	svc := &IoTService{}
	store := newJobsExecutionTestStore(t)
	seedJobExecution(t, store, "job-zero", "thing-zero", "QUEUED")
	err := svc.cancelJobExecutionCore(store, CancelJobExecutionInput{
		JobID: "job-zero", ThingName: "thing-zero",
		ExpectedVersion:         0,
		ExpectedVersionProvided: true,
	})
	if !errors.Is(err, iotstore.ErrVersionConflict) {
		t.Fatalf("expected ErrVersionConflict for an explicit zero, got %v", err)
	}
	if status, _ := executionStatus(t, store, "job-zero", "thing-zero"); status != "QUEUED" {
		t.Fatalf("rejected cancel must preserve QUEUED, got %s", status)
	}
}

// TestOTAUpdateMaterialisesJob pins the awsIotJobId contract: the create
// materialises the IoT job it names, and the update cannot be deleted
// without forceDeleteAWSJob while that job is not terminal.
func TestOTAUpdateMaterialisesJob(t *testing.T) {
	svc := &IoTService{}
	store := newJobsExecutionTestStore(t)
	thing := "arn:aws:iot:us-east-1:000000000000:thing/ota-thing"
	result, err := svc.createOTAUpdateCore(store, CreateOTAUpdateInput{
		OtaUpdateID: "ota-1",
		RoleArn:     "arn:aws:iam::000000000000:role/ota",
		Targets:     []string{thing},
		Files:       []interface{}{map[string]interface{}{"fileName": "fw.bin"}},
	})
	if err != nil {
		t.Fatalf("create OTA update: %v", err)
	}
	job, err := store.GetJob(result.AwsIotJobID)
	if err != nil {
		t.Fatalf("the awsIotJobId must identify a real job: %v", err)
	}
	if len(job.Targets) != 1 || job.Targets[0] != thing {
		t.Fatalf("job targets = %#v", job.Targets)
	}
	// A non-terminal job blocks the delete without the force flag.
	err = svc.deleteOTAUpdateCore(store, DeleteOTAUpdateInput{OtaUpdateID: "ota-1"})
	if !errors.Is(err, iotstore.ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for a non-terminal job, got %v", err)
	}
	// The force flag deletes both the update and the job.
	if err := svc.deleteOTAUpdateCore(store, DeleteOTAUpdateInput{OtaUpdateID: "ota-1", ForceDeleteAWSJob: true}); err != nil {
		t.Fatalf("forced delete: %v", err)
	}
	if _, err := store.GetJob(result.AwsIotJobID); err == nil {
		t.Fatal("expected the job to be deleted with the update")
	}
}
