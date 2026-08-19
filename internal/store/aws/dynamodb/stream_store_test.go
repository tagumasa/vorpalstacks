package dynamodb

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func TestStreamStoreTrimOlderThan(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	tbl := "TrimTbl"
	keys := map[string]interface{}{"pk": "a"}

	addRecord := func() {
		t.Helper()
		if _, err := store.AddRecord(tbl, "arn:aws:dynamodb:us-east-1:123456789012:table/"+tbl+"/stream/1",
			"NEW_AND_OLD_IMAGES", StreamEventInsert, keys, nil, nil, nil); err != nil {
			t.Fatalf("add record: %v", err)
		}
	}
	for i := 0; i < 3; i++ {
		addRecord()
	}

	// A cut-off in the past trims nothing: the records were just created.
	if err := store.TrimOlderThan(tbl, time.Now().Add(-time.Hour)); err != nil {
		t.Fatalf("trim with past cut-off: %v", err)
	}
	floor, err := store.OldestSequence(tbl)
	if err != nil {
		t.Fatalf("read floor: %v", err)
	}
	if floor != 0 {
		t.Fatalf("expected floor 0 before trimming, got %d", floor)
	}
	records, _, err := store.GetRecords(tbl, 0, 10)
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 records before trimming, got %d", len(records))
	}

	// A future cut-off trims every record and advances the floor to the
	// highest removed sequence number.
	if err := store.TrimOlderThan(tbl, time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("trim with future cut-off: %v", err)
	}
	floor, err = store.OldestSequence(tbl)
	if err != nil {
		t.Fatalf("read floor after trim: %v", err)
	}
	if floor != 3 {
		t.Fatalf("expected floor 3 after trimming, got %d", floor)
	}
	records, _, err = store.GetRecords(tbl, 0, 10)
	if err != nil {
		t.Fatalf("get records after trim: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records after trimming, got %d", len(records))
	}

	// New records continue after the floor and remain readable from it.
	addRecord()
	records, _, err = store.GetRecords(tbl, floor, 10)
	if err != nil {
		t.Fatalf("get records from floor: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected the post-trim record, got %d", len(records))
	}

	// Re-running with a past cut-off must not roll the floor back.
	if err := store.TrimOlderThan(tbl, time.Now().Add(-time.Hour)); err != nil {
		t.Fatalf("re-trim with past cut-off: %v", err)
	}
	floor, err = store.OldestSequence(tbl)
	if err != nil {
		t.Fatalf("read floor after re-trim: %v", err)
	}
	if floor != 3 {
		t.Fatalf("floor must never decrease, got %d", floor)
	}
}

func TestStreamStoreOldestSequenceWithoutRecords(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	floor, err := store.OldestSequence("NeverSeenTable")
	if err != nil {
		t.Fatalf("oldest sequence on unknown table: %v", err)
	}
	if floor != 0 {
		t.Fatalf("expected floor 0 for a table without records, got %d", floor)
	}
}
