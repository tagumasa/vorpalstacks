package sfn

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// GetActivityTask is a bounded long poll: with no task available the
// handler must return an empty taskToken at the poll deadline (60 seconds
// per the Step Functions API reference) rather than blocking until the
// client disconnects. The deadline is shortened here so the test is fast.
func TestGetActivityTaskReturnsEmptyAtPollDeadline(t *testing.T) {
	origTimeout := sfnstore.ActivityTaskPollTimeout
	sfnstore.ActivityTaskPollTimeout = 50 * time.Millisecond
	t.Cleanup(func() { sfnstore.ActivityTaskPollTimeout = origTimeout })

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewStepFunctionService(mgr, "000000000000")
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")

	activityArn := "arn:aws:states:us-east-1:000000000000:activity:poll"
	if _, err := svc.CreateActivity(context.Background(), reqCtx, &request.ParsedRequest{
		Parameters: map[string]interface{}{"name": "poll"},
	}); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	result, err := svc.GetActivityTask(context.Background(), reqCtx, &request.ParsedRequest{
		Parameters: map[string]interface{}{"activityArn": activityArn},
	})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}

	resp, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected response type %T", result)
	}
	if resp["taskToken"] != "" {
		t.Fatalf("expected empty taskToken at poll deadline, got %v", resp["taskToken"])
	}
	if elapsed > 5*time.Second {
		t.Fatalf("poll held the request for %v; the deadline was not applied", elapsed)
	}
}
