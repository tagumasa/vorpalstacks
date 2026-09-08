package lambda

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// provisionedConcurrencyStore builds a FunctionStore over a fresh empty
// storage with one function created.
func provisionedConcurrencyStore(t *testing.T) *FunctionStore {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	st, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("get storage: %v", err)
	}
	s := NewFunctionStore(st, "000000000000", "us-east-1")
	if _, err := s.Create(&Function{
		FunctionName: "fn",
		Runtime:      RuntimeNodejs22X,
		Role:         "arn:aws:iam::000000000000:role/lambda",
	}); err != nil {
		t.Fatalf("create function: %v", err)
	}
	return s
}

// TestSetProvisionedConcurrencyReportsInProgress pins the provisioning
// bookkeeping: a Put records the requested capacity as IN_PROGRESS with
// zero allocated and available — the warm capacity is only reported once
// the pre-warmed environment exists.
func TestSetProvisionedConcurrencyReportsInProgress(t *testing.T) {
	s := provisionedConcurrencyStore(t)

	if err := s.SetProvisionedConcurrency("fn", "1", 5); err != nil {
		t.Fatalf("set provisioned concurrency: %v", err)
	}
	pc, err := s.GetProvisionedConcurrency("fn", "1")
	if err != nil {
		t.Fatalf("get provisioned concurrency: %v", err)
	}
	if pc.Status != "IN_PROGRESS" {
		t.Fatalf("Status %q, want IN_PROGRESS", pc.Status)
	}
	if pc.RequestedProvisionedConcurrentExecutions != 5 {
		t.Fatalf("Requested %d, want 5", pc.RequestedProvisionedConcurrentExecutions)
	}
	if pc.AllocatedProvisionedConcurrentExecutions != 0 || pc.AvailableProvisionedConcurrentExecutions != 0 {
		t.Fatalf("Allocated/Available %d/%d, want 0/0",
			pc.AllocatedProvisionedConcurrentExecutions, pc.AvailableProvisionedConcurrentExecutions)
	}
}

// TestMarkProvisionedConcurrencyReadyFlipsStatus pins the readiness
// transition: the pre-warm completion reports the warm capacity that
// actually exists, and an unknown qualifier fails.
func TestMarkProvisionedConcurrencyReadyFlipsStatus(t *testing.T) {
	s := provisionedConcurrencyStore(t)

	if err := s.SetProvisionedConcurrency("fn", "1", 5); err != nil {
		t.Fatalf("set provisioned concurrency: %v", err)
	}
	if err := s.MarkProvisionedConcurrencyReady("fn", "1", 1); err != nil {
		t.Fatalf("mark ready: %v", err)
	}
	pc, err := s.GetProvisionedConcurrency("fn", "1")
	if err != nil {
		t.Fatalf("get provisioned concurrency: %v", err)
	}
	if pc.Status != "READY" {
		t.Fatalf("Status %q, want READY", pc.Status)
	}
	if pc.AllocatedProvisionedConcurrentExecutions != 1 || pc.AvailableProvisionedConcurrentExecutions != 1 {
		t.Fatalf("Allocated/Available %d/%d, want 1/1",
			pc.AllocatedProvisionedConcurrentExecutions, pc.AvailableProvisionedConcurrentExecutions)
	}
	if pc.RequestedProvisionedConcurrentExecutions != 5 {
		t.Fatalf("Requested %d, want 5", pc.RequestedProvisionedConcurrentExecutions)
	}

	if err := s.MarkProvisionedConcurrencyReady("fn", "999", 1); !errors.Is(err, ErrProvisionedConcurrencyNotFound) {
		t.Fatalf("unknown qualifier: expected ErrProvisionedConcurrencyNotFound, got %v", err)
	}
}
