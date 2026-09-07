package dynamodb

import (
	"context"
	"fmt"
	"sync"
	"testing"

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
	records, _, err := store.GetRecords("Tbl", 0, 10)
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

	records, _, err := store.GetRecords("Tbl", 0, 1000)
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

// TestAllocateSeqSeedsFromRecords pins restart seeding: allocation must
// continue from the highest existing record key, not from the persisted
// counter alone, because transactions commit out of allocation order and
// can leave the counter behind the records.
func TestAllocateSeqSeedsFromRecords(t *testing.T) {
	dir := t.TempDir()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	seed := NewStreamStore(st, "123456789012", "us-east-1")
	// Committed state after a commit-order inversion: the counter says 3
	// while a record with sequence 7 exists.
	if err := seed.BaseStore.PutProto(streamSeqKey("Tbl"), streamCounterToProto(streamSeqCounter{LastSeq: 3})); err != nil {
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
	record, err := restarted.AddRecord("Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
		StreamEventInsert, map[string]interface{}{"id": "x"}, nil, nil, nil)
	if err != nil {
		t.Fatalf("add record after restart: %v", err)
	}
	if record.Dynamodb.SequenceNumber != "8" {
		t.Fatalf("first sequence after restart = %s, want 8 (highest record key 7, not the counter's 3)", record.Dynamodb.SequenceNumber)
	}
	// The persisted counter catches up to the allocation.
	latest, err := restarted.GetLatestSequence("Tbl")
	if err != nil {
		t.Fatalf("latest sequence: %v", err)
	}
	if latest != 8 {
		t.Fatalf("persisted latest sequence = %d, want 8", latest)
	}
}
