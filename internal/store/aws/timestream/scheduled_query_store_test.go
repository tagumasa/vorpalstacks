package timestream

import (
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newScheduledQueryTestStore(t *testing.T) *ScheduledQueryStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewScheduledQueryStore(st, "000000000000", "us-east-1")
}

// TestScheduledQueryRecordWritesSerialised pins that concurrent record
// writers (the engine's run bookkeeping and a state update) each run their
// whole read-modify-write cycle inside the record lock, so no writer loses
// another's fields.
func TestScheduledQueryRecordWritesSerialised(t *testing.T) {
	store := newScheduledQueryTestStore(t)
	if _, err := store.CreateScheduledQuery(
		"concurrent",
		"SELECT 1",
		&ScheduleConfiguration{ScheduleExpression: "rate(5 minutes)"},
		nil, "", "", nil, nil, "",
	); err != nil {
		t.Fatalf("create scheduled query: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := store.UpdateLastRun("concurrent", "AUTO_SUCCEEDED", time.Now().UTC()); err != nil {
				t.Errorf("update last run: %v", err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := store.UpdateNextRunTime("concurrent", time.Now().UTC().Add(time.Minute)); err != nil {
				t.Errorf("update next run time: %v", err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if _, err := store.UpdateScheduledQuery("concurrent", ScheduledQueryStatusDisabled); err != nil {
				t.Errorf("update scheduled query: %v", err)
				return
			}
		}
	}()
	wg.Wait()

	got, err := store.GetScheduledQuery("concurrent")
	if err != nil {
		t.Fatalf("get scheduled query: %v", err)
	}
	if got.PreviousRunTime.IsZero() {
		t.Error("PreviousRunTime lost to a concurrent record write")
	}
	if got.LastRunStatus == "" {
		t.Error("LastRunStatus lost to a concurrent record write")
	}
	if got.NextRunTime.IsZero() {
		t.Error("NextRunTime lost to a concurrent record write")
	}
	if got.ScheduledQueryStatus != ScheduledQueryStatusDisabled {
		t.Error("ScheduledQueryStatus lost to a concurrent record write")
	}
}

// TestDeleteScheduledQueryNotResurrectedByRunWrite pins the delete
// serialisation: a run-bookkeeping write racing the delete can never write
// the record back after the delete removed it.
func TestDeleteScheduledQueryNotResurrectedByRunWrite(t *testing.T) {
	store := newScheduledQueryTestStore(t)

	for i := 0; i < 50; i++ {
		if _, err := store.CreateScheduledQuery(
			"victim",
			"SELECT 1",
			&ScheduleConfiguration{ScheduleExpression: "rate(5 minutes)"},
			nil, "", "", nil, nil, "",
		); err != nil {
			t.Fatalf("create scheduled query (iteration %d): %v", i, err)
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			// A not-found error is a legitimate outcome: the delete may
			// have landed first.
			_ = store.UpdateLastRun("victim", "AUTO_SUCCEEDED", time.Now().UTC())
		}()
		go func() {
			defer wg.Done()
			_ = store.DeleteScheduledQuery("victim")
		}()
		wg.Wait()

		if _, err := store.GetScheduledQuery("victim"); err != ErrScheduledQueryNotFound {
			t.Fatalf("iteration %d: scheduled query resurrected after delete (err = %v)", i, err)
		}
	}
}
