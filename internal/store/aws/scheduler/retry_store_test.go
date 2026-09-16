package scheduler

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// TestGetDueRetryRecordsReapsCorruptRecords pins the due scan's corruption
// reap: a record whose stored value fails deserialisation is deleted rather
// than skipped. The corrupt key sorts at or before every future cutoff, so a
// skip-only scan would revisit it on every engine tick — permanently wasted
// work and an unbounded accumulation of unreadable records.
func TestGetDueRetryRecordsReapsCorruptRecords(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	rs := NewRetryStore(st, "us-east-1")

	now := time.Now()
	healthy := &RetryRecord{
		ID:            "healthy-1",
		ScheduleName:  "reap-pin",
		GroupName:     "default",
		NextAttemptAt: now.Add(-time.Minute),
	}
	if err := rs.SaveRetryRecord(healthy); err != nil {
		t.Fatalf("save healthy record: %v", err)
	}
	corruptKey := retryStorageKey(now.Add(-time.Second), "corrupt-1")
	if err := rs.store.PutRaw(corruptKey, []byte("not json")); err != nil {
		t.Fatalf("seed corrupt record: %v", err)
	}

	due, err := rs.GetDueRetryRecords(now)
	if err != nil {
		t.Fatalf("get due records: %v", err)
	}
	if len(due) != 1 || due[0].ID != "healthy-1" {
		t.Fatalf("due records = %+v, want only the healthy record", due)
	}
	if rs.store.Exists(corruptKey) {
		t.Error("corrupt record still present after the due scan, want reaped")
	}
}
