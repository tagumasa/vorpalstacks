package sfn

import (
	"context"
	"encoding/json"
	"testing"
)

// TestRetryPolicyUnmarshalDefaults pins the retrier decode: MaxAttempts
// absent means the documented default of three, while an explicit zero
// keeps its never-retry meaning ("A value of 0 specifies that the error is
// never retried"); JitterStrategy decodes onto the policy.
func TestRetryPolicyUnmarshalDefaults(t *testing.T) {
	var absent RetryPolicy
	if err := json.Unmarshal([]byte(`{"ErrorEquals":["States.ALL"]}`), &absent); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if absent.MaxAttempts != 3 {
		t.Errorf("absent MaxAttempts = %d, want the documented default 3", absent.MaxAttempts)
	}

	var explicit RetryPolicy
	if err := json.Unmarshal([]byte(`{"ErrorEquals":["States.ALL"],"MaxAttempts":0}`), &explicit); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if explicit.MaxAttempts != 0 {
		t.Errorf("explicit MaxAttempts 0 = %d, want 0 (never retried)", explicit.MaxAttempts)
	}

	var full RetryPolicy
	if err := json.Unmarshal([]byte(`{"ErrorEquals":["States.ALL"],"MaxAttempts":2,"JitterStrategy":"FULL"}`), &full); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if full.MaxAttempts != 2 || full.JitterStrategy != "FULL" {
		t.Errorf("decoded = %+v, want MaxAttempts 2 and JitterStrategy FULL", full)
	}
}

// TestVersionCounterResetsAfterMachineDeletion pins that deleting a state
// machine drops its in-memory version counter: a same-name machine's first
// publish must number from 1 again (recovering from the now-empty persisted
// records) instead of resuming at the deleted machine's high-water mark.
func TestVersionCounterResetsAfterMachineDeletion(t *testing.T) {
	store := newHistoryTestStore(t)
	ctx := context.Background()

	sm := &StateMachine{Name: "version-recycle"}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}
	first, err := store.PublishStateMachineVersion(ctx, sm.StateMachineArn, "v1")
	if err != nil {
		t.Fatal(err)
	}
	if first.Version != 1 {
		t.Fatalf("first publish version = %d, want 1", first.Version)
	}
	if err := store.DeleteStateMachine(ctx, sm.StateMachineArn); err != nil {
		t.Fatal(err)
	}

	again := &StateMachine{Name: "version-recycle"}
	if err := store.CreateStateMachine(ctx, again); err != nil {
		t.Fatal(err)
	}
	recreated, err := store.PublishStateMachineVersion(ctx, again.StateMachineArn, "fresh")
	if err != nil {
		t.Fatal(err)
	}
	if recreated.Version != 1 {
		t.Fatalf("recreated machine's first publish version = %d, want 1 (stale high-water mark survived the deletion)", recreated.Version)
	}
}
