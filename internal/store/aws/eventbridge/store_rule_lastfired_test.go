package eventbridge

import (
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newRuleTouchTestStore(t *testing.T) *EventsStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewEventsStore(st, "000000000000", "us-east-1")
}

// TestTouchRuleLastFired pins the persisted fire marker: the boundary is
// recorded on the rule without counting as a modification, the marker
// only ever advances, and touching an unknown rule reports not-found.
func TestTouchRuleLastFired(t *testing.T) {
	store := newRuleTouchTestStore(t)
	rule := &Rule{
		Name:               "scheduled",
		EventBusName:       "default",
		ScheduleExpression: "rate(5 minutes)",
	}
	if err := store.CreateRule(t.Context(), rule); err != nil {
		t.Fatalf("create rule: %v", err)
	}
	created, err := store.GetRule(t.Context(), "default", "scheduled")
	if err != nil {
		t.Fatalf("get rule: %v", err)
	}
	if !created.LastFiredAt.IsZero() {
		t.Fatalf("fresh rule carries a fire marker: %v", created.LastFiredAt)
	}

	fired := time.Date(2026, 8, 22, 12, 5, 0, 0, time.UTC)
	if err := store.TouchRuleLastFired(t.Context(), "default", "scheduled", fired); err != nil {
		t.Fatalf("touch last fired: %v", err)
	}

	got, err := store.GetRule(t.Context(), "default", "scheduled")
	if err != nil {
		t.Fatalf("get rule after touch: %v", err)
	}
	if !got.LastFiredAt.Equal(fired) {
		t.Errorf("lastFiredAt = %v, want %v", got.LastFiredAt, fired)
	}
	if !got.LastModifiedAt.Equal(created.LastModifiedAt) {
		t.Errorf("firing stamped lastModifiedAt: got %v, want %v (unchanged)", got.LastModifiedAt, created.LastModifiedAt)
	}

	// An older in-flight fire cannot regress a newer marker.
	older := fired.Add(-time.Minute)
	if err := store.TouchRuleLastFired(t.Context(), "default", "scheduled", older); err != nil {
		t.Fatalf("touch older boundary: %v", err)
	}
	got, err = store.GetRule(t.Context(), "default", "scheduled")
	if err != nil {
		t.Fatalf("get rule after older touch: %v", err)
	}
	if !got.LastFiredAt.Equal(fired) {
		t.Errorf("older boundary regressed the marker: %v, want %v", got.LastFiredAt, fired)
	}

	// A newer boundary advances the marker.
	newer := fired.Add(5 * time.Minute)
	if err := store.TouchRuleLastFired(t.Context(), "default", "scheduled", newer); err != nil {
		t.Fatalf("touch newer boundary: %v", err)
	}
	got, err = store.GetRule(t.Context(), "default", "scheduled")
	if err != nil {
		t.Fatalf("get rule after newer touch: %v", err)
	}
	if !got.LastFiredAt.Equal(newer) {
		t.Errorf("newer boundary did not advance the marker: %v, want %v", got.LastFiredAt, newer)
	}

	// Touching an unknown rule reports the not-found error.
	if err := store.TouchRuleLastFired(t.Context(), "default", "missing", fired); err != ErrRuleNotFound {
		t.Errorf("touching an unknown rule = %v, want %v", err, ErrRuleNotFound)
	}
}

// TestMutateRuleConcurrentWithTouch pins the record-mutation serialisation:
// a user field update and a delivery-marker write racing on the same rule
// must both survive, because each runs its whole read-modify-write cycle
// inside the record lock — never a stale pre-read record.
func TestMutateRuleConcurrentWithTouch(t *testing.T) {
	store := newRuleTouchTestStore(t)
	if err := store.CreateRule(t.Context(), &Rule{
		Name:               "racy",
		EventBusName:       "default",
		ScheduleExpression: "rate(5 minutes)",
		Description:        "before",
	}); err != nil {
		t.Fatalf("create rule: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// The service-side updater pattern: mutate the user's fields through
	// the store-level atomic mutation.
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := store.MutateRule(t.Context(), "default", "racy", func(rule *Rule) error {
				rule.Description = "after"
				rule.LastModifiedAt = time.Now().UTC()
				return nil
			}); err != nil {
				t.Errorf("mutate rule: %v", err)
				return
			}
		}
	}()
	// The delivery-worker marker pattern, racing from another goroutine.
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := store.TouchRuleLastFired(t.Context(), "default", "racy", time.Now().UTC()); err != nil {
				t.Errorf("touch last fired: %v", err)
				return
			}
		}
	}()
	wg.Wait()

	got, err := store.GetRule(t.Context(), "default", "racy")
	if err != nil {
		t.Fatalf("final get: %v", err)
	}
	if got.Description != "after" {
		t.Errorf("user fields lost to the marker write: description = %q, want %q", got.Description, "after")
	}
	if got.LastFiredAt.IsZero() {
		t.Error("fire marker lost to the concurrent rule update")
	}
}

// TestDeleteRuleNotResurrectedByTouch pins the delete serialisation: a
// delivery-marker write racing the delete can never write the record back
// after the delete removed it.
func TestDeleteRuleNotResurrectedByTouch(t *testing.T) {
	store := newRuleTouchTestStore(t)
	fired := time.Date(2026, 8, 22, 12, 5, 0, 0, time.UTC)

	for i := 0; i < 50; i++ {
		if err := store.CreateRule(t.Context(), &Rule{
			Name:               "victim",
			EventBusName:       "default",
			ScheduleExpression: "rate(5 minutes)",
		}); err != nil {
			t.Fatalf("create rule (iteration %d): %v", i, err)
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			// A not-found error is a legitimate outcome: the delete may
			// have landed first.
			_ = store.TouchRuleLastFired(t.Context(), "default", "victim", fired)
		}()
		go func() {
			defer wg.Done()
			_ = store.DeleteRule(t.Context(), "default", "victim")
		}()
		wg.Wait()

		if _, err := store.GetRule(t.Context(), "default", "victim"); err != ErrRuleNotFound {
			t.Fatalf("iteration %d: rule resurrected after delete (err = %v)", i, err)
		}
	}
}
