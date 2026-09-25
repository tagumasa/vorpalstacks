package dynamodb

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

const testStreamArn = "arn:aws:dynamodb:us-east-1:123456789012:table/Tbl/stream/2026-09-06T00:00:00.000"

// addRecordTxn adds one stream record inside its own transaction, the
// putItemCore shape.
func addRecordTxn(t *testing.T, st storage.TransactionalStorageWith2PC, s *StreamStore, table string, eventID string) *StreamRecord {
	t.Helper()
	var record *StreamRecord
	err := st.Update(context.Background(), func(txn storage.Transaction) error {
		r, err := s.AddRecordTxn(txn, table, testStreamArn, "NEW_AND_OLD_IMAGES",
			StreamEventInsert, map[string]interface{}{"id": eventID}, nil, nil, nil)
		record = r
		return err
	})
	if err != nil {
		t.Fatalf("add record txn: %v", err)
	}
	return record
}

// TestAddRecordTxnSameTxnDistinctSeqs pins the TransactWriteItems shape:
// two records added for one table inside a single transaction must carry
// distinct consecutive sequence numbers — the per-operation counter RMW
// made both collide on one key and the second silently overwrote the
// first.
func TestAddRecordTxnSameTxnDistinctSeqs(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	var first, second *StreamRecord
	err = st.Update(context.Background(), func(txn storage.Transaction) error {
		first, err = store.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
			StreamEventInsert, map[string]interface{}{"id": "a"}, nil, nil, nil)
		if err != nil {
			return err
		}
		second, err = store.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
			StreamEventInsert, map[string]interface{}{"id": "b"}, nil, nil, nil)
		return err
	})
	if err != nil {
		t.Fatalf("add two records in one transaction: %v", err)
	}

	if first.Dynamodb.SequenceNumber == second.Dynamodb.SequenceNumber {
		t.Fatalf("both records carry sequence number %s", first.Dynamodb.SequenceNumber)
	}
	if first.EventID == second.EventID {
		t.Fatalf("both records carry event ID %s", first.EventID)
	}
	records, _, err := store.GetRecords("Tbl", testStreamArn, 0, 10)
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected both records persisted, got %d (the second must not overwrite the first)", len(records))
	}
}

// TestAddRecordTxnConcurrentNoDuplicateKeys pins the concurrent-writer
// shape: transactions running in parallel for one streamed table must
// produce unique record keys — nothing is overwritten, every sequence
// number is distinct.
func TestAddRecordTxnConcurrentNoDuplicateKeys(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	const writers = 16
	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		w := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := st.Update(context.Background(), func(txn storage.Transaction) error {
				_, err := store.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
					StreamEventInsert, map[string]interface{}{"id": fmt.Sprintf("w%d", w)}, nil, nil, nil)
				return err
			})
			if err != nil {
				t.Errorf("writer %d: %v", w, err)
			}
		}()
	}
	wg.Wait()

	records, _, err := store.GetRecords("Tbl", testStreamArn, 0, 1000)
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(records) != writers {
		t.Fatalf("expected %d records, got %d — duplicate keys overwrote records", writers, len(records))
	}
	seen := make(map[string]bool, len(records))
	for _, rec := range records {
		if seen[rec.Dynamodb.SequenceNumber] {
			t.Fatalf("duplicate sequence number %s", rec.Dynamodb.SequenceNumber)
		}
		seen[rec.Dynamodb.SequenceNumber] = true
	}
}

// TestAllocateSeqSeedsFromRecordsAndFloor pins restart seeding now that the
// counter persists no extent: allocation continues from the highest existing
// record key (committed record keys are the stream's extent — the numbering
// never restarts below a number a surviving record used), and when the
// retention sweep has removed the records the persisted floor keeps the
// numbering beyond every number a trim removed.
func TestAllocateSeqSeedsFromRecordsAndFloor(t *testing.T) {
	dir := t.TempDir()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	seed := NewStreamStore(st, "123456789012", "us-east-1")
	// Committed state: a record with sequence 7 exists (and a floor at 3
	// a past trim left behind).
	if err := seed.BaseStore.PutProto(streamSeqKey("Tbl"), streamCounterToProto(streamSeqCounter{TrimmedFloor: 3})); err != nil {
		t.Fatalf("seed counter: %v", err)
	}
	if err := seed.BaseStore.PutProto(streamRecordKey("Tbl", 7), streamRecordToProto(&StreamRecord{})); err != nil {
		t.Fatalf("seed record: %v", err)
	}
	st.Close()

	st2, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("reopen storage: %v", err)
	}
	defer st2.Close()
	restarted := NewStreamStore(st2, "123456789012", "us-east-1")
	var record *StreamRecord
	err = st2.Update(context.Background(), func(txn storage.Transaction) error {
		r, rerr := restarted.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
			StreamEventInsert, map[string]interface{}{"id": "x"}, nil, nil, nil)
		record = r
		return rerr
	})
	if err != nil {
		t.Fatalf("add record after restart: %v", err)
	}
	// The wire form of a sequence number is the model's 21-digit
	// zero-padded rendering; parsed as an integer it is 8.
	if record.Dynamodb.SequenceNumber != FormatStreamSequenceNumber(8) {
		t.Fatalf("first sequence after restart = %s, want %s (highest record key 7, not the floor's 3)", record.Dynamodb.SequenceNumber, FormatStreamSequenceNumber(8))
	}
	latest, err := restarted.GetLatestSequenceForStream("Tbl", "")
	if err != nil {
		t.Fatalf("latest sequence: %v", err)
	}
	if latest != 8 {
		t.Fatalf("persisted latest sequence = %d, want 8", latest)
	}

	// The trim-swept state: no records survive, the floor does — the next
	// number sits above every number the trim removed.
	swept := t.TempDir()
	st3, err := storage.Open(swept)
	if err != nil {
		t.Fatalf("open swept storage: %v", err)
	}
	sweptSeed := NewStreamStore(st3, "123456789012", "us-east-1")
	if err := sweptSeed.BaseStore.PutProto(streamSeqKey("Tbl"), streamCounterToProto(streamSeqCounter{TrimmedFloor: 12})); err != nil {
		t.Fatalf("seed swept counter: %v", err)
	}
	st3.Close()

	st4, err := storage.Open(swept)
	if err != nil {
		t.Fatalf("reopen swept storage: %v", err)
	}
	defer st4.Close()
	afterSweep := NewStreamStore(st4, "123456789012", "us-east-1")
	var sweptRecord *StreamRecord
	err = st4.Update(context.Background(), func(txn storage.Transaction) error {
		r, rerr := afterSweep.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
			StreamEventInsert, map[string]interface{}{"id": "y"}, nil, nil, nil)
		sweptRecord = r
		return rerr
	})
	if err != nil {
		t.Fatalf("add record after sweep: %v", err)
	}
	if sweptRecord.Dynamodb.SequenceNumber != FormatStreamSequenceNumber(13) {
		t.Fatalf("first sequence after a full trim = %s, want %s (above the floor's 12)", sweptRecord.Dynamodb.SequenceNumber, FormatStreamSequenceNumber(13))
	}
}

// TestRetentionFloorSurvivesRecordCommits pins the trim floor's single
// writer: the retention sweep alone writes the counter, so a record
// transaction that allocated its sequence before a trim and committed after
// it cannot rewrite the floor the trim raised — the record path once wrote
// the whole counter inside its carrying transaction, last-writer-wins over
// the sweep.
func TestRetentionFloorSurvivesRecordCommits(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	// Three committed records the trim will remove.
	for i := 0; i < 3; i++ {
		addRecordTxn(t, st, store, "Tbl", fmt.Sprintf("c%d", i))
	}

	allocated := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- st.Update(context.Background(), func(txn storage.Transaction) error {
			_, err := store.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
				StreamEventInsert, map[string]interface{}{"id": "racer"}, nil, nil, nil)
			if err != nil {
				return err
			}
			allocated <- struct{}{}
			<-release
			return nil
		})
	}()
	<-allocated

	// The trim removes the three committed records and raises the floor
	// while the racing record's transaction is still open.
	if err := store.TrimOlderThan("Tbl", time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("trim: %v", err)
	}
	floor, err := store.OldestSequence("Tbl")
	if err != nil {
		t.Fatalf("floor after trim: %v", err)
	}
	if floor != 3 {
		t.Fatalf("floor after trim = %d, want 3", floor)
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("racing transaction: %v", err)
	}
	// The racing record's commit must not rewrite the floor its allocation
	// predated: the sweep is the counter's only writer.
	floor, err = store.OldestSequence("Tbl")
	if err != nil {
		t.Fatalf("floor after commit: %v", err)
	}
	if floor != 3 {
		t.Fatalf("floor after the racing record committed = %d, want 3 (the record path must not rewrite the counter)", floor)
	}
	records, _, err := store.GetRecords("Tbl", testStreamArn, 0, 10)
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(records) != 1 || records[0].Dynamodb.SequenceNumber != FormatStreamSequenceNumber(4) {
		t.Fatalf("the racing record must survive the trim it outlived, got %v", records)
	}
}

// TestGetLatestSequenceSurvivesInvertedCommitOrder pins the stream extent
// under the inverted commit order last-writer-wins batches allow: the
// transaction that allocated the lower sequence commits after the one that
// allocated the higher. LATEST iterators and DescribeStream's ending
// number must still report the higher committed record — the extent is
// derived from the record keys, not from the counter the slower
// transaction overwrites.
func TestGetLatestSequenceSurvivesInvertedCommitOrder(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewStreamStore(st, "123456789012", "us-east-1")
	aAllocated := make(chan struct{})
	releaseA := make(chan struct{})
	done := make(chan error, 2)

	// A allocates sequence 1 and holds its transaction open.
	go func() {
		done <- st.Update(context.Background(), func(txn storage.Transaction) error {
			_, err := store.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
				StreamEventInsert, map[string]interface{}{"id": "a"}, nil, nil, nil)
			if err != nil {
				return err
			}
			aAllocated <- struct{}{}
			<-releaseA
			return nil
		})
	}()
	<-aAllocated

	// B allocates sequence 2 and commits while A is still open.
	go func() {
		done <- st.Update(context.Background(), func(txn storage.Transaction) error {
			_, err := store.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
				StreamEventInsert, map[string]interface{}{"id": "b"}, nil, nil, nil)
			return err
		})
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		visible, _, err := store.GetRecords("Tbl", testStreamArn, 1, 1)
		if err != nil {
			t.Fatalf("poll records: %v", err)
		}
		if len(visible) == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the second transaction never committed — the storage engine serialised the two Updates")
		}
		time.Sleep(2 * time.Millisecond)
	}

	close(releaseA)
	if err := <-done; err != nil {
		t.Fatalf("first transaction: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatalf("second transaction: %v", err)
	}

	latest, err := store.GetLatestSequenceForStream("Tbl", "")
	if err != nil {
		t.Fatalf("latest sequence: %v", err)
	}
	if latest != 2 {
		t.Fatalf("latest sequence after inverted commit order = %d, want 2 — the committed higher record must win over the slower transaction's counter", latest)
	}
	records, _, err := store.GetRecords("Tbl", testStreamArn, 0, 10)
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected both committed records, got %d", len(records))
	}
}
