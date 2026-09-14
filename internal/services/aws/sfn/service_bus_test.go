package sfn

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
)

// TestStartExecutionEventAnswersStoreFaultAsFailure pins the bus contract on
// the cross-service start-execution handler: a storage fault while acquiring
// the region store is a failed delivery reported in HandlerResult.Error.
// Synchronous publishers (the Scheduler engine publishes via PublishSync and
// folds result.Error into their retry and dead-letter policy) would
// otherwise record the failed start as a successful delivery and lose it
// silently — the sibling core-error path already carries the propagation
// contract this test extends to the store-acquisition fault.
func TestStartExecutionEventAnswersStoreFaultAsFailure(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("occupies the region-directory name"), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: blocker})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewStepFunctionService(mgr, "000000000000")

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Input:           "{}",
	}
	evt.Region = "us-east-1"
	result := svc.handleStartExecutionEvent(context.Background(), evt)
	if result.Error == nil {
		t.Fatal("a store-acquisition fault must be reported in HandlerResult.Error, not recorded as success")
	}
}
