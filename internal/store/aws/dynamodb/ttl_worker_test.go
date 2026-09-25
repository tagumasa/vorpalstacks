package dynamodb

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newTTLTestStore opens a store with one string-keyed table whose TTL is
// enabled on the "ttl" attribute and whose OLD_IMAGE stream is on — the
// shape the TTL worker walks in production.
func newTTLTestStore(t *testing.T) *DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	if _, err := store.tables.Create(CreateTableParams{
		Name:                 "TtlTbl",
		KeySchema:            []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
		AttributeDefinitions: []*AttributeDefinition{{AttributeName: "id", AttributeType: ScalarAttributeTypeS}},
		BillingMode:          BillingModePayPerRequest,
		StreamSpecification: &StreamSpecification{
			StreamEnabled:  true,
			StreamViewType: StreamViewTypeOldImage,
		},
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if err := store.tables.SetTimeToLive("TtlTbl", &TimeToLiveSpecification{
		Enabled:       true,
		AttributeName: "ttl",
	}); err != nil {
		t.Fatalf("enable ttl: %v", err)
	}
	return store
}

func putTTLItem(t *testing.T, store *DynamoDBStore, id string, attributes map[string]*AttributeValue) {
	t.Helper()
	key := map[string]*AttributeValue{"id": {S: &id}}
	err := store.Update(context.Background(), func(txn *DynamoDBTxn) error {
		return txn.PutItem("TtlTbl", key, attributes)
	})
	if err != nil {
		t.Fatalf("put item %s: %v", id, err)
	}
}

// readTTLItem returns the live item or nil when the key holds nothing.
func readTTLItem(t *testing.T, store *DynamoDBStore, id string) *Item {
	t.Helper()
	key := map[string]*AttributeValue{"id": {S: &id}}
	var item *Item
	err := store.Update(context.Background(), func(txn *DynamoDBTxn) error {
		it, err := txn.GetItem("TtlTbl", key)
		if err != nil {
			if IsItemNotFound(err) {
				return nil
			}
			return err
		}
		item = it
		return nil
	})
	if err != nil {
		t.Fatalf("read item %s: %v", id, err)
	}
	return item
}

// A full cleanup pass deletes exactly the items expired at the pass's
// decision clock: future, absent and non-numeric TTL attributes all survive.
func TestTTLCleanupDeletesOnlyExpiredItems(t *testing.T) {
	store := newTTLTestStore(t)
	const now = int64(1_000_000)

	putTTLItem(t, store, "expired", map[string]*AttributeValue{
		"ttl": numAttr("999900"),
		"v":   strAttr("expired"),
	})
	putTTLItem(t, store, "future", map[string]*AttributeValue{
		"ttl": numAttr("1003600"),
		"v":   strAttr("future"),
	})
	putTTLItem(t, store, "farfuture", map[string]*AttributeValue{
		"ttl": numAttr("9.9E+125"),
		"v":   strAttr("farfuture"),
	})
	putTTLItem(t, store, "absent", map[string]*AttributeValue{
		"v": strAttr("absent"),
	})
	putTTLItem(t, store, "nonnumeric", map[string]*AttributeValue{
		"ttl": numAttr("soon"),
		"v":   strAttr("nonnumeric"),
	})

	table, err := store.tables.Get("TtlTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	w := &ttlWorker{store: store, ctx: context.Background()}
	w.cleanupTableTTL(table, "ttl", now)

	if it := readTTLItem(t, store, "expired"); it != nil {
		t.Fatalf("expired item survived the cleanup pass: %+v", it)
	}
	for _, id := range []string{"future", "farfuture", "absent", "nonnumeric"} {
		if it := readTTLItem(t, store, id); it == nil {
			t.Fatalf("item %s deleted despite a non-expiring TTL attribute", id)
		}
	}
}

// The expiry decision stays in the float domain: a TTL number beyond the
// int64 range is a valid DynamoDB Number far in the future, and converting
// such a float to int64 is implementation-dependent — amd64 answers MinInt64,
// which reads as long expired — so the predicate must compare the parsed
// float against the clock without the conversion.
func TestTTLExpiryComparesBeyondInt64InTheFloatDomain(t *testing.T) {
	const now = int64(1_000_000)
	farFuture := map[string]*AttributeValue{"ttl": numAttr("9.9E+125")}
	if ttlExpired(farFuture, "ttl", now) {
		t.Fatal("a TTL beyond int64 read as expired")
	}
	justBeyondInt64 := map[string]*AttributeValue{"ttl": numAttr("9.3E+18")}
	if ttlExpired(justBeyondInt64, "ttl", now) {
		t.Fatal("a TTL just beyond the int64 ceiling read as expired")
	}
	deepPast := map[string]*AttributeValue{"ttl": numAttr("-9.9E+125")}
	if !ttlExpired(deepPast, "ttl", now) {
		t.Fatal("a TTL far in the past read as unexpired")
	}
}

// A write that lands between the scan's snapshot and the deletion
// transaction must survive the pass: extending the TTL attribute into the
// future (or removing it) makes the item no longer expired, so the TTL
// process does not delete it — the documented contract the in-transaction
// re-check restores. Each subtest drives the transaction phase directly
// with the stale scan key, which is the deterministic form of the race.
func TestTTLCleanupSparesItemsRewrittenAfterTheScanDecision(t *testing.T) {
	store := newTTLTestStore(t)
	const now = int64(1_000_000)

	table, err := store.tables.Get("TtlTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	w := &ttlWorker{store: store, ctx: context.Background()}

	t.Run("ttl extended into the future", func(t *testing.T) {
		putTTLItem(t, store, "extended", map[string]*AttributeValue{
			"ttl": numAttr("999900"),
			"v":   strAttr("old"),
		})
		key := map[string]*AttributeValue{"id": {S: ptrTTLString("extended")}}
		// The client's extend write commits before the transaction opens.
		putTTLItem(t, store, "extended", map[string]*AttributeValue{
			"ttl": numAttr("2000000"),
			"v":   strAttr("new"),
		})

		deleted, err := w.deleteIfStillExpired(table, "ttl", now, key)
		if err != nil {
			t.Fatalf("deleteIfStillExpired: %v", err)
		}
		if deleted {
			t.Fatal("item reported deleted after its TTL was extended")
		}
		it := readTTLItem(t, store, "extended")
		if it == nil {
			t.Fatal("item deleted despite an extended TTL attribute — the extend write was lost")
		}
		if got := it.Attributes["ttl"]; got == nil || got.N == nil || *got.N != "2000000" {
			t.Fatalf("surviving item lost the extended TTL attribute: %+v", got)
		}
		if got := it.Attributes["v"]; got == nil || got.S == nil || *got.S != "new" {
			t.Fatalf("surviving item lost the extend write's value: %+v", got)
		}
		// A spared item emits no stream record: nothing was removed.
		records, _, err := store.streams.GetRecords("TtlTbl", table.StreamArn, 0, 100)
		if err != nil {
			t.Fatalf("get records: %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("spared item emitted %d stream records", len(records))
		}
	})

	t.Run("ttl attribute removed", func(t *testing.T) {
		putTTLItem(t, store, "removed", map[string]*AttributeValue{
			"ttl": numAttr("999900"),
			"v":   strAttr("old"),
		})
		key := map[string]*AttributeValue{"id": {S: ptrTTLString("removed")}}
		putTTLItem(t, store, "removed", map[string]*AttributeValue{
			"v": strAttr("new"),
		})

		deleted, err := w.deleteIfStillExpired(table, "ttl", now, key)
		if err != nil {
			t.Fatalf("deleteIfStillExpired: %v", err)
		}
		if deleted {
			t.Fatal("item reported deleted after its TTL attribute was removed")
		}
		if it := readTTLItem(t, store, "removed"); it == nil {
			t.Fatal("item deleted despite its TTL attribute being removed — the rewrite was lost")
		}
	})

	t.Run("still expired is deleted", func(t *testing.T) {
		putTTLItem(t, store, "doomed", map[string]*AttributeValue{
			"ttl": numAttr("999900"),
		})
		key := map[string]*AttributeValue{"id": {S: ptrTTLString("doomed")}}

		deleted, err := w.deleteIfStillExpired(table, "ttl", now, key)
		if err != nil {
			t.Fatalf("deleteIfStillExpired: %v", err)
		}
		if !deleted {
			t.Fatal("genuinely expired item spared by the re-check")
		}
		if it := readTTLItem(t, store, "doomed"); it != nil {
			t.Fatalf("genuinely expired item survived: %+v", it)
		}
	})

	t.Run("deleted by a client before the transaction", func(t *testing.T) {
		putTTLItem(t, store, "gone", map[string]*AttributeValue{
			"ttl": numAttr("999900"),
		})
		key := map[string]*AttributeValue{"id": {S: ptrTTLString("gone")}}
		err := store.Update(context.Background(), func(txn *DynamoDBTxn) error {
			return txn.DeleteItem("TtlTbl", key)
		})
		if err != nil {
			t.Fatalf("client delete: %v", err)
		}

		deleted, err := w.deleteIfStillExpired(table, "ttl", now, key)
		if err != nil {
			t.Fatalf("deleteIfStillExpired after client delete: %v", err)
		}
		if deleted {
			t.Fatal("client-deleted item counted as a TTL deletion")
		}
	})
}

// The REMOVE record describes the item the transaction actually deleted —
// keys and old image from the re-read record, the TTL service identity
// attached so consumers can tell expiry from client deletes.
func TestTTLCleanupStreamRecordDescribesTheDeletedItem(t *testing.T) {
	store := newTTLTestStore(t)
	const now = int64(1_000_000)

	putTTLItem(t, store, "victim", map[string]*AttributeValue{
		"ttl": numAttr("999900"),
		"v":   strAttr("old"),
	})

	table, err := store.tables.Get("TtlTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	w := &ttlWorker{store: store, ctx: context.Background()}
	w.cleanupTableTTL(table, "ttl", now)

	records, _, err := store.streams.GetRecords("TtlTbl", table.StreamArn, 0, 100)
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected one REMOVE record, got %d", len(records))
	}
	r := records[0]
	if r.EventName != StreamEventRemove {
		t.Fatalf("event name = %q, want REMOVE", r.EventName)
	}
	if id, ok := r.Dynamodb.Keys["id"].(map[string]interface{}); !ok ||
		id["S"] != "victim" {
		t.Fatalf("record keys = %+v, want the deleted item's key", r.Dynamodb.Keys)
	}
	if v, ok := r.Dynamodb.OldImage["v"].(map[string]interface{}); !ok || v["S"] != "old" {
		t.Fatalf("record old image = %+v, want the deleted item's attributes", r.Dynamodb.OldImage)
	}
	if r.UserIdentity == nil || r.UserIdentity.Type != "Service" ||
		r.UserIdentity.PrincipalID != "dynamodb.amazonaws.com" {
		t.Fatalf("record user identity = %+v, want the TTL service identity", r.UserIdentity)
	}
}

// ptrTTLString gives each key literal its own pointer, matching how live
// keys are built at their use sites.
func ptrTTLString(s string) *string { return &s }

// A stream-capture failure aborts the whole TTL removal: the delete and
// its REMOVE record are one transaction ("the whole removal describes one
// consistent item"), so the capture fault leaves the item in place for
// the next pass instead of committing a delete no stream consumer ever
// sees. The counter record's corrupted payload is the injected capture
// fault — the TTL delete is the table's first stream write (the store's
// own item paths do not stream), so the capture seeds its sequence from
// that record and fails on it; healing the record lets the next pass run
// the whole removal.
func TestTTLCleanupStreamFailureAbortsTheDelete(t *testing.T) {
	store := newTTLTestStore(t)
	const now = int64(1_000_000)

	putTTLItem(t, store, "victim", map[string]*AttributeValue{
		"ttl": numAttr("999900"),
		"v":   strAttr("old"),
	})

	table, err := store.tables.Get("TtlTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	streamBucket := store.Storage().Bucket(streamBucketName("us-east-1"))
	counterKey := []byte(streamSeqKey("TtlTbl"))
	if err := streamBucket.Put(counterKey, []byte("not-a-proto-record")); err != nil {
		t.Fatalf("corrupt stream sequence counter: %v", err)
	}

	w := &ttlWorker{store: store, ctx: context.Background()}
	w.cleanupTableTTL(table, "ttl", now)

	if item := readTTLItem(t, store, "victim"); item == nil {
		t.Fatal("item deleted despite the stream-capture failure")
	}

	// Heal the counter: the next pass runs the whole removal — item gone,
	// exactly one REMOVE record.
	if err := streamBucket.Delete(counterKey); err != nil {
		t.Fatalf("heal stream sequence counter: %v", err)
	}
	w.cleanupTableTTL(table, "ttl", now)
	if item := readTTLItem(t, store, "victim"); item != nil {
		t.Fatal("item survived the post-heal pass")
	}
	records, _, err := store.streams.GetRecords("TtlTbl", table.StreamArn, 0, 100)
	if err != nil {
		t.Fatalf("get records after the post-heal pass: %v", err)
	}
	if len(records) != 1 || records[0].EventName != StreamEventRemove {
		t.Fatalf("records after the post-heal pass = %d, want one REMOVE", len(records))
	}
}

// The tick containment keeps the worker alive across a panicking pass: a
// store-level fault inside one cleanup tick is logged and the ticker keeps
// ticking instead of taking the process down. The nil store makes the very
// first store call of the pass panic, which stands deterministically for
// any panic inside the pass.
func TestTTLWorkerTickContainsPanics(t *testing.T) {
	w := &ttlWorker{ctx: context.Background()}
	w.tick()
}
