package eventbridge

import (
	"strings"
	"testing"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestPutTargetsRetryPolicyProvidedMemberSemantics pins the age member's
// provided/unprovided distinction: an explicitly supplied
// MaximumEventAgeInSeconds outside 60-86400 — zero included — is a
// per-entry ValidationException, while an omitted member (or a policy
// carrying the attempts member alone) keeps the deployment default and
// passes.
func TestPutTargetsRetryPolicyProvidedMemberSemantics(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := t.Context()

	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	if err := store.CreateRule(ctx, &eventsstore.Rule{Name: "r", EventBusName: "default"}); err != nil {
		t.Fatalf("create rule: %v", err)
	}

	target := func(id string, retry map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"Id":          id,
			"Arn":         "arn:aws:sqs:us-east-1:000000000000:queue",
			"RetryPolicy": retry,
		}
	}

	result, err := svc.putTargetsCore(ctx, store, PutTargetsInput{
		EventBusName:         "default",
		EventBusNameProvided: true,
		Rule:                 "r",
		Region:               "us-east-1",
		Targets: []interface{}{
			target("t-zero", map[string]interface{}{"MaximumEventAgeInSeconds": float64(0)}),
			target("t-attempts-only", map[string]interface{}{"MaximumRetryAttempts": float64(5)}),
			target("t-min", map[string]interface{}{"MaximumEventAgeInSeconds": float64(60)}),
		},
	})
	if err != nil {
		t.Fatalf("putTargetsCore: %v", err)
	}
	if result.FailedEntryCount != 1 {
		t.Fatalf("failedEntryCount = %d, want 1 (entries: %v)", result.FailedEntryCount, result.FailedEntries)
	}
	for _, failed := range result.FailedEntries {
		if failed["ErrorCode"] != "ValidationException" {
			t.Errorf("target %v: error code = %v, want ValidationException", failed["TargetId"], failed["ErrorCode"])
		}
		msg, _ := failed["ErrorMessage"].(string)
		if !strings.Contains(msg, "when provided") {
			t.Errorf("target %v: message %q must state the provided-member range", failed["TargetId"], msg)
		}
	}

	// The bounds hold on both sides of the window; the attempts-only
	// policy stores with the age unset (the deployment default applies at
	// delivery time).
	bounds, err := svc.putTargetsCore(ctx, store, PutTargetsInput{
		EventBusName:         "default",
		EventBusNameProvided: true,
		Rule:                 "r",
		Region:               "us-east-1",
		Targets: []interface{}{
			target("t-max", map[string]interface{}{"MaximumEventAgeInSeconds": float64(86400)}),
			target("t-under", map[string]interface{}{"MaximumEventAgeInSeconds": float64(59)}),
			target("t-over", map[string]interface{}{"MaximumEventAgeInSeconds": float64(86401)}),
		},
	})
	if err != nil {
		t.Fatalf("putTargetsCore bounds: %v", err)
	}
	if bounds.FailedEntryCount != 2 {
		t.Fatalf("bounds failedEntryCount = %d, want 2 (entries: %v)", bounds.FailedEntryCount, bounds.FailedEntries)
	}

	stored, err := store.GetTarget(ctx, "default", "r", "t-attempts-only")
	if err != nil {
		t.Fatalf("get stored target: %v", err)
	}
	if stored.RetryPolicy == nil || stored.RetryPolicy.MaximumRetryAttempts != 5 || stored.RetryPolicy.MaximumEventAgeInSeconds != 0 {
		t.Fatalf("attempts-only policy must store 5 attempts and an unset age, got %+v", stored.RetryPolicy)
	}
}

// TestRemoveTargetsUnknownTargetSucceeds pins the RemoveTargets contract
// for a target ID that does not exist on the rule: the removal succeeds
// (zero failed entries) rather than reporting an InternalFailure entry
// with the raw store error string. The API reference frames a successful
// RemoveTargets as "the target(s) listed in the request are removed" and
// documents no per-entry error code for a missing target.
func TestRemoveTargetsUnknownTargetSucceeds(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := t.Context()

	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	if err := store.CreateRule(ctx, &eventsstore.Rule{Name: "r", EventBusName: "default"}); err != nil {
		t.Fatalf("create rule: %v", err)
	}
	if err := store.PutTarget(ctx, &eventsstore.Target{
		ID:           "real",
		RuleName:     "r",
		EventBusName: "default",
		ARN:          "arn:aws:events:us-east-1:000000000000:event-bus/default",
	}); err != nil {
		t.Fatalf("put target: %v", err)
	}

	// A mixed request: one live target, one that was never there.
	result, err := svc.removeTargetsCore(ctx, store, RemoveTargetsInput{
		EventBusName: "default",
		Rule:         "r",
		Ids:          []string{"real", "ghost"},
	})
	if err != nil {
		t.Fatalf("removeTargetsCore: %v", err)
	}
	if result.FailedEntryCount != 0 {
		t.Errorf("failedEntryCount = %d, want 0 (entries: %v)", result.FailedEntryCount, result.FailedEntries)
	}

	// The live target is gone; repeating the request stays a success.
	if _, err := store.GetTarget(ctx, "default", "r", "real"); err != eventsstore.ErrTargetNotFound {
		t.Errorf("live target not removed (err = %v)", err)
	}
	again, err := svc.removeTargetsCore(ctx, store, RemoveTargetsInput{
		EventBusName: "default",
		Rule:         "r",
		Ids:          []string{"real"},
	})
	if err != nil {
		t.Fatalf("repeat removeTargetsCore: %v", err)
	}
	if again.FailedEntryCount != 0 {
		t.Errorf("repeat failedEntryCount = %d, want 0", again.FailedEntryCount)
	}
}
