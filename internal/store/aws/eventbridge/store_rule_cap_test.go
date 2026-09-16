package eventbridge

import (
	"context"
	"errors"
	"testing"
)

// TestCreateRuleCapped pins the per-bus rule-count quota at the store
// layer: creating a NEW rule on a full bus answers ErrRuleCapReached,
// an existing rule name at the cap still answers ErrRuleAlreadyExists
// (the quota bounds creation alone), the cap is per bus, deleting a
// rule frees capacity, and maxPerBus <= 0 disables the gate.
func TestCreateRuleCapped(t *testing.T) {
	store, _ := newRecordLockTestStores(t)
	ctx := context.Background()

	rule := func(bus, name string) *Rule {
		return &Rule{Name: name, EventBusName: bus, State: RuleStateEnabled}
	}

	if err := store.CreateRuleCapped(ctx, rule("cap-bus", "r1"), 2); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if err := store.CreateRuleCapped(ctx, rule("cap-bus", "r2"), 2); err != nil {
		t.Fatalf("second create: %v", err)
	}
	if err := store.CreateRuleCapped(ctx, rule("cap-bus", "r3"), 2); !errors.Is(err, ErrRuleCapReached) {
		t.Fatalf("third create at cap: got %v, want ErrRuleCapReached", err)
	}
	if err := store.CreateRuleCapped(ctx, rule("cap-bus", "r1"), 2); !errors.Is(err, ErrRuleAlreadyExists) {
		t.Fatalf("existing name at cap: got %v, want ErrRuleAlreadyExists", err)
	}

	// The cap counts one bus: a full cap-bus does not bound another bus.
	if err := store.CreateRuleCapped(ctx, rule("other-bus", "o1"), 2); err != nil {
		t.Fatalf("create on a different bus under a full bus: %v", err)
	}

	// Deleting a rule frees capacity for the next create.
	if err := store.DeleteRule(ctx, "cap-bus", "r2"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := store.CreateRuleCapped(ctx, rule("cap-bus", "r3"), 2); err != nil {
		t.Fatalf("create after freeing capacity: %v", err)
	}

	// maxPerBus <= 0 disables the cap — the uncapped CreateRule contract.
	if err := store.CreateRule(ctx, rule("cap-bus", "r4")); err != nil {
		t.Fatalf("uncapped create at cap: %v", err)
	}
}
