package cloudwatchlogs

import (
	"sync"
	"testing"
	"time"
)

func newRaceScheduledQuery(id string) *ScheduledQuery {
	now := time.Now().UTC().UnixMilli()
	return &ScheduledQuery{
		Id:                 id,
		Name:               "race-" + id,
		QueryString:        "fields @message",
		ScheduleExpression: "rate(1 minute)",
		State:              "ENABLED",
		CreationTime:       now,
		LastUpdatedTime:    now,
	}
}

// TestMutateScheduledQueryConcurrentWithTouch pins that a user-driven
// mutation and a delivery touch racing on the same record both survive:
// the record lock serialises the two read-modify-write cycles.
func TestMutateScheduledQueryConcurrentWithTouch(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.PutScheduledQuery(newRaceScheduledQuery("race-mutate")); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = s.MutateScheduledQuery("race-mutate", func(sq *ScheduledQuery) error {
				sq.Description = "after"
				return nil
			})
		}()
		go func() {
			defer wg.Done()
			stamp := time.Now().UTC().UnixMilli()
			_ = s.TouchScheduledQueryDelivery("race-mutate", stamp, stamp, ScheduledQueryStatusComplete)
		}()
	}
	wg.Wait()

	sq, err := s.GetScheduledQuery("race-mutate")
	if err != nil {
		t.Fatal(err)
	}
	if sq.Description != "after" {
		t.Fatalf("mutation lost: description = %q", sq.Description)
	}
	if sq.LastTriggeredTime == 0 || sq.LastExecutionStatus == "" {
		t.Fatalf("delivery touch lost: lastTriggeredTime = %d, lastExecutionStatus = %q",
			sq.LastTriggeredTime, sq.LastExecutionStatus)
	}
}

// TestDeleteScheduledQueryNotResurrectedByTouch pins that a delivery
// touch racing a delete cannot write the record back: once the delete
// wins, the record stays gone.
func TestDeleteScheduledQueryNotResurrectedByTouch(t *testing.T) {
	s := newLogsTestStore(t)
	for i := 0; i < 50; i++ {
		const id = "race-delete"
		if err := s.PutScheduledQuery(newRaceScheduledQuery(id)); err != nil {
			t.Fatal(err)
		}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			stamp := time.Now().UTC().UnixMilli()
			_ = s.TouchScheduledQueryDelivery(id, stamp, stamp, ScheduledQueryStatusComplete)
		}()
		go func() {
			defer wg.Done()
			_ = s.DeleteScheduledQuery(id)
		}()
		wg.Wait()

		if _, err := s.GetScheduledQuery(id); err != ErrResourceNotFound {
			t.Fatalf("iteration %d: record resurrected after delete (err = %v)", i, err)
		}
	}
}

// TestTouchScheduledQueryDeliveryMarkerOnlyAdvances pins that the
// consumed-boundary marker only moves forward and that a delivery touch
// neither advances lastUpdatedTime nor drops the execution fields.
func TestTouchScheduledQueryDeliveryMarkerOnlyAdvances(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.PutScheduledQuery(newRaceScheduledQuery("touch-roundtrip")); err != nil {
		t.Fatal(err)
	}
	created, err := s.GetScheduledQuery("touch-roundtrip")
	if err != nil {
		t.Fatal(err)
	}

	newer := time.Now().UTC().UnixMilli()
	older := time.Now().Add(-time.Hour).UTC().UnixMilli()
	if err := s.TouchScheduledQueryDelivery("touch-roundtrip", newer, newer, ScheduledQueryStatusComplete); err != nil {
		t.Fatal(err)
	}
	if err := s.TouchScheduledQueryDelivery("touch-roundtrip", older, older, ScheduledQueryStatusFailed); err != nil {
		t.Fatal(err)
	}

	got, err := s.GetScheduledQuery("touch-roundtrip")
	if err != nil {
		t.Fatal(err)
	}
	if got.LastExecutedBoundary != newer {
		t.Fatalf("marker regressed: %d, want %d", got.LastExecutedBoundary, newer)
	}
	if got.LastTriggeredTime != older || got.LastExecutionStatus != ScheduledQueryStatusFailed {
		t.Fatalf("execution fields = (%d, %q), want the latest touch (%d, %q)",
			got.LastTriggeredTime, got.LastExecutionStatus, older, ScheduledQueryStatusFailed)
	}
	if got.LastUpdatedTime != created.LastUpdatedTime {
		t.Fatalf("lastUpdatedTime advanced on a delivery touch: %d -> %d", created.LastUpdatedTime, got.LastUpdatedTime)
	}
}
