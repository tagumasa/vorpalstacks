package dynamodb

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func strAttr(v string) *AttributeValue {
	return &AttributeValue{S: &v}
}

func TestJournalStoreReverseReplayOrder(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	region := "us-east-1"
	store := NewJournalStore(st, region)
	base := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	key := map[string]*AttributeValue{"id": strAttr("a")}

	appendRecord := func(at time.Time, operation string, before map[string]*AttributeValue) {
		t.Helper()
		if err := st.Update(t.Context(), func(txn storage.Transaction) error {
			return appendJournalTxnAt(txn, region, "JournalTbl", operation, key, before, at)
		}); err != nil {
			t.Fatalf("append journal record at %v: %v", at, err)
		}
	}

	appendRecord(base, JournalOperationPut, nil)
	appendRecord(base.Add(time.Second), JournalOperationPut, map[string]*AttributeValue{"id": strAttr("a"), "v": strAttr("1")})
	appendRecord(base.Add(2*time.Second), JournalOperationPut, map[string]*AttributeValue{"id": strAttr("a"), "v": strAttr("2")})
	appendRecord(base.Add(3*time.Second), JournalOperationDelete, map[string]*AttributeValue{"id": strAttr("a"), "v": strAttr("3")})

	var replayed []*JournalChange
	if err := store.ReverseReplay("JournalTbl", base.Add(time.Second), func(change *JournalChange) error {
		replayed = append(replayed, change)
		return nil
	}); err != nil {
		t.Fatalf("reverse replay: %v", err)
	}

	if len(replayed) != 2 {
		t.Fatalf("expected 2 changes newer than the cut-off, got %d", len(replayed))
	}
	// Newest change first so the caller can undo mutations in reverse order.
	if replayed[0].Operation != JournalOperationDelete || replayed[0].BeforeImage["v"].S == nil || *replayed[0].BeforeImage["v"].S != "3" {
		t.Fatalf("newest change must be the delete at t=3s restoring v=3, got op=%q before=%v", replayed[0].Operation, replayed[0].BeforeImage)
	}
	if replayed[1].Operation != JournalOperationPut || replayed[1].BeforeImage["v"].S == nil || *replayed[1].BeforeImage["v"].S != "2" {
		t.Fatalf("second change must be the put at t=2s restoring v=2, got op=%q before=%v", replayed[1].Operation, replayed[1].BeforeImage)
	}

	// A cut-off after every record replays nothing.
	count := 0
	if err := store.ReverseReplay("JournalTbl", base.Add(time.Hour), func(change *JournalChange) error {
		count++
		return nil
	}); err != nil {
		t.Fatalf("reverse replay past all records: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no changes past every record, got %d", count)
	}
}

func TestJournalStoreDeleteOlderThan(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	region := "us-east-1"
	store := NewJournalStore(st, region)
	base := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	key := map[string]*AttributeValue{"id": strAttr("b")}

	for i := range 3 {
		at := base.Add(time.Duration(i) * time.Minute)
		if err := st.Update(t.Context(), func(txn storage.Transaction) error {
			return appendJournalTxnAt(txn, region, "PruneTbl", JournalOperationPut, key, nil, at)
		}); err != nil {
			t.Fatalf("append journal record %d: %v", i, err)
		}
	}

	removed, err := store.DeleteOlderThan("PruneTbl", base.Add(time.Minute))
	if err != nil {
		t.Fatalf("prune: %v", err)
	}
	if removed != 2 {
		t.Fatalf("expected 2 pruned records (t=0 and t=1 are at or before the cut-off), got %d", removed)
	}

	remaining := 0
	if err := store.ReverseReplay("PruneTbl", base, func(change *JournalChange) error {
		remaining++
		return nil
	}); err != nil {
		t.Fatalf("replay after prune: %v", err)
	}
	if remaining != 1 {
		t.Fatalf("expected 1 surviving record, got %d", remaining)
	}

	// Pruning is idempotent.
	removed, err = store.DeleteOlderThan("PruneTbl", base.Add(time.Minute))
	if err != nil || removed != 0 {
		t.Fatalf("second prune must remove nothing, got removed=%d err=%v", removed, err)
	}
}

func TestDynamoDBTxnJournalsPITRWrites(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(
		"PitrTbl",
		[]*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
		[]*AttributeDefinition{{AttributeName: "id", AttributeType: ScalarAttributeTypeS}},
		BillingModePayPerRequest, nil, nil, nil, nil, nil, false,
	); err != nil {
		t.Fatalf("create table: %v", err)
	}
	// Without recovery enabled no journal records are appended.
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.PutItem("PitrTbl", map[string]*AttributeValue{"id": strAttr("k")}, map[string]*AttributeValue{"v": strAttr("0")})
	}); err != nil {
		t.Fatalf("put without recovery: %v", err)
	}
	journaled := 0
	if err := store.Journal().ReverseReplay("PitrTbl", time.Time{}, func(change *JournalChange) error {
		journaled++
		return nil
	}); err != nil {
		t.Fatalf("replay without recovery: %v", err)
	}
	if journaled != 0 {
		t.Fatalf("writes without recovery must not be journaled, got %d records", journaled)
	}

	if err := store.Tables().SetPointInTimeRecovery("PitrTbl", &PointInTimeRecoveryDescription{
		Status:                     PITRStatusEnabled,
		EarliestRestorableDateTime: time.Now().Add(-time.Minute),
	}); err != nil {
		t.Fatalf("enable recovery: %v", err)
	}

	// Insert, overwrite, then delete: three journal records carrying the
	// before-image needed to undo each mutation.
	steps := []struct {
		put   bool
		attrs map[string]*AttributeValue
	}{
		{put: true, attrs: map[string]*AttributeValue{"v": strAttr("1")}},
		{put: true, attrs: map[string]*AttributeValue{"v": strAttr("2")}},
		{put: false},
	}
	for i, step := range steps {
		if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
			key := map[string]*AttributeValue{"id": strAttr("k")}
			if step.put {
				return txn.PutItem("PitrTbl", key, step.attrs)
			}
			return txn.DeleteItem("PitrTbl", key)
		}); err != nil {
			t.Fatalf("mutation step %d: %v", i, err)
		}
	}

	var replayed []*JournalChange
	if err := store.Journal().ReverseReplay("PitrTbl", time.Now().Add(-time.Minute), func(change *JournalChange) error {
		replayed = append(replayed, change)
		return nil
	}); err != nil {
		t.Fatalf("replay with recovery: %v", err)
	}
	if len(replayed) != 3 {
		t.Fatalf("expected 3 journaled mutations, got %d", len(replayed))
	}
	if replayed[0].Operation != JournalOperationDelete || replayed[0].BeforeImage["v"].S == nil || *replayed[0].BeforeImage["v"].S != "2" {
		t.Fatalf("newest record must be the delete with before-image v=2, got op=%q before=%v", replayed[0].Operation, replayed[0].BeforeImage)
	}
	if replayed[1].Operation != JournalOperationPut || replayed[1].BeforeImage["v"].S == nil || *replayed[1].BeforeImage["v"].S != "1" {
		t.Fatalf("middle record must be the put with before-image v=1, got op=%q before=%v", replayed[1].Operation, replayed[1].BeforeImage)
	}
	if replayed[2].Operation != JournalOperationPut || replayed[2].BeforeImage["v"].S == nil || *replayed[2].BeforeImage["v"].S != "0" {
		t.Fatalf("oldest record must be the first put with the pre-recovery item as before-image, got op=%q before=%v", replayed[2].Operation, replayed[2].BeforeImage)
	}

	// Undoing the journal in replay order must restore the state before
	// the first journaled write (the item as it was with v=0).
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		for _, change := range replayed {
			if change.BeforeImage == nil {
				if err := txn.DeleteItem("PitrTbl", change.Key); err != nil {
					return err
				}
				continue
			}
			if err := txn.PutItem("PitrTbl", change.Key, change.BeforeImage); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("undo replay: %v", err)
	}
	item, err := store.Items().Get("PitrTbl", map[string]*AttributeValue{"id": strAttr("k")})
	if err != nil {
		t.Fatalf("item must exist after undoing the whole journal: %v", err)
	}
	if item.Attributes["v"].S == nil || *item.Attributes["v"].S != "0" {
		t.Fatalf("item must hold the pre-recovery value v=0 after undo, got %v", item.Attributes["v"])
	}
}
