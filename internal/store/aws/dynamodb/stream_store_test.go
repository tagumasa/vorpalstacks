package dynamodb

import (
	"context"
	"strconv"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// addRecordTo adds one stream record inside its own transaction — the
// putItemCore shape AddRecordTxn exists for.
func addRecordTo(t *testing.T, st storage.TransactionalStorageWith2PC, s *StreamStore, table, streamArn string) *StreamRecord {
	t.Helper()
	var record *StreamRecord
	if err := st.Update(context.Background(), func(txn storage.Transaction) error {
		r, err := s.AddRecordTxn(txn, table, streamArn, "NEW_AND_OLD_IMAGES",
			StreamEventInsert, map[string]interface{}{"pk": "a"}, nil, nil, nil)
		record = r
		return err
	}); err != nil {
		t.Fatalf("add record txn: %v", err)
	}
	return record
}

func TestStreamStoreTrimOlderThan(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	tbl := "TrimTbl"
	streamArn := "arn:aws:dynamodb:us-east-1:123456789012:table/" + tbl + "/stream/1"

	addRecord := func() {
		t.Helper()
		addRecordTo(t, st, store, tbl, streamArn)
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
	records, _, err := store.GetRecords(tbl, streamArn, 0, 10)
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
	records, _, err = store.GetRecords(tbl, streamArn, 0, 10)
	if err != nil {
		t.Fatalf("get records after trim: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records after trimming, got %d", len(records))
	}

	// New records continue after the floor and remain readable from it.
	addRecord()
	records, _, err = store.GetRecords(tbl, streamArn, floor, 10)
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

// TestHighestSequenceSkipsForeignShortKeys pins the reverse extent scan's
// length guard: a record key always ends in the 20-digit zero-padded
// sequence, so a foreign or corrupt key shorter than that under the table
// prefix cannot be a record — the extent read skips it rather than
// panicking on the suffix slice, and the highest committed record still
// wins. The foreign key sorts above every record key, so the reverse scan
// meets it first.
func TestHighestSequenceSkipsForeignShortKeys(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	tbl := "T"
	streamArn := "arn:aws:dynamodb:us-east-1:123456789012:table/" + tbl + "/stream/1"
	addRecordTo(t, st, store, tbl, streamArn)

	bucket := st.Bucket(streamBucketName("us-east-1"))
	if err := bucket.Put([]byte(tbl+KeySep+"zz"), []byte("foreign")); err != nil {
		t.Fatalf("seed foreign key: %v", err)
	}

	seq, err := store.GetLatestSequenceForStream(tbl, streamArn)
	if err != nil {
		t.Fatalf("latest sequence with foreign short key: %v", err)
	}
	if seq != 1 {
		t.Fatalf("latest sequence = %d, want the record's 1", seq)
	}
}

// TestExtentCountsOnlyTheNamedGeneration pins the generation-scoped extent:
// residue of a superseded generation the best-effort sweep missed lives
// under the same table prefix with a higher sequence than the current
// generation's own records, and the extent the current generation reports
// — its EndingSequenceNumber and LATEST position — must not be advanced by
// it, while the allocator's cross-generation seeding keeps counting every
// record so numbering stays unique and increasing.
func TestExtentCountsOnlyTheNamedGeneration(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	tbl := "GenTbl"
	currentArn := "arn:aws:dynamodb:us-east-1:123456789012:table/" + tbl + "/stream/2026-09-24T00:00:00.000"
	supersededArn := "arn:aws:dynamodb:us-east-1:123456789012:table/" + tbl + "/stream/2026-09-01T00:00:00.000"
	add := func(streamArn string) {
		t.Helper()
		addRecordTo(t, st, store, tbl, streamArn)
	}

	add(currentArn)
	add(supersededArn) // residue: higher sequence, wrong generation
	add(currentArn)

	if seq, err := store.GetLatestSequenceForStream(tbl, currentArn); err != nil || seq != 3 {
		t.Fatalf("current generation extent = %d/%v, want 3 (the residue's 2 must not advance it)", seq, err)
	}
	if seq, err := store.GetLatestSequenceForStream(tbl, supersededArn); err != nil || seq != 2 {
		t.Fatalf("superseded generation extent = %d/%v, want its own 2", seq, err)
	}

	// The allocator keeps numbering above every record, residue included:
	// the next record of the current generation continues past the residue.
	record := addRecordTo(t, st, store, tbl, currentArn)
	if got := record.Dynamodb.SequenceNumber; got != FormatStreamSequenceNumber(4) {
		t.Fatalf("post-residue sequence = %s, want 4 (numbering crosses generations)", got)
	}
}

// TestStreamStoreGenerationSpace pins the generation semantics of the
// record space: reads serve only the generation whose ARN they name, and
// the generation reset empties the space so the successor generation
// starts fresh while its numbering stays unique and increasing.
func TestStreamStoreGenerationSpace(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	tbl := "GenTbl"
	firstArn := "arn:aws:dynamodb:us-east-1:123456789012:table/" + tbl + "/stream/2026-01-01T00:00:00.000"
	secondArn := "arn:aws:dynamodb:us-east-1:123456789012:table/" + tbl + "/stream/2026-01-02T00:00:00.000"

	for range 2 {
		addRecordTo(t, st, store, tbl, firstArn)
	}

	// The owning generation reads its records; another generation of the
	// same table reads none of them.
	first, _, err := store.GetRecords(tbl, firstArn, 0, 10)
	if err != nil {
		t.Fatalf("read first generation: %v", err)
	}
	if len(first) != 2 {
		t.Fatalf("first generation must read its 2 records, got %d", len(first))
	}
	second, nextSeq, err := store.GetRecords(tbl, secondArn, 0, 10)
	if err != nil {
		t.Fatalf("read second generation: %v", err)
	}
	if len(second) != 0 {
		t.Fatalf("second generation must read none of the first generation's records, got %d", len(second))
	}
	if nextSeq != 0 {
		t.Fatalf("a read that serves nothing must keep its position, got next seq %d", nextSeq)
	}

	// The generation reset empties the space entirely.
	if err := store.DeleteTableRecords(tbl); err != nil {
		t.Fatalf("reset generation space: %v", err)
	}
	if latest, err := store.GetLatestSequenceForStream(tbl, ""); err != nil || latest != 0 {
		t.Fatalf("reset must leave an empty record space, latest seq = %d (err %v)", latest, err)
	}
	empty, _, err := store.GetRecords(tbl, firstArn, 0, 10)
	if err != nil {
		t.Fatalf("read after reset: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("the superseded generation's records must be gone after reset, got %d", len(empty))
	}

	// The successor generation's records stay unique and increasing — the
	// allocator keeps its high-water mark across the reset — and read back
	// under their own ARN.
	addRecordTo(t, st, store, tbl, secondArn)
	afterReset, _, err := store.GetRecords(tbl, secondArn, 0, 10)
	if err != nil {
		t.Fatalf("read second generation after reset: %v", err)
	}
	if len(afterReset) != 1 {
		t.Fatalf("second generation must read its own record after reset, got %d", len(afterReset))
	}
	if seq, err := strconv.ParseInt(afterReset[0].Dynamodb.SequenceNumber, 10, 64); err != nil || seq <= 2 {
		t.Fatalf("the successor generation's numbering must stay above the reset space's extent, got seq %d (err %v)", seq, err)
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

// The shard-iterator signing key is generated once, persisted in the stream
// bucket, and read back unchanged by a later store instance over the same
// storage — restarts must not invalidate iterators issued before them.
func TestIteratorSigningKeyGeneratedOnceAndPersisted(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	first := NewStreamStore(st, "123456789012", "us-east-1")
	key1, err := first.IteratorSigningKey()
	if err != nil {
		t.Fatalf("first key: %v", err)
	}
	if len(key1) != 32 {
		t.Fatalf("key length = %d, want 32", len(key1))
	}
	again, err := first.IteratorSigningKey()
	if err != nil {
		t.Fatalf("cached key: %v", err)
	}
	if string(key1) != string(again) {
		t.Fatalf("cached key must be stable within one store")
	}

	reopened := NewStreamStore(st, "123456789012", "us-east-1")
	key2, err := reopened.IteratorSigningKey()
	if err != nil {
		t.Fatalf("reopened key: %v", err)
	}
	if string(key1) != string(key2) {
		t.Fatalf("key must survive a store reopen")
	}
}
