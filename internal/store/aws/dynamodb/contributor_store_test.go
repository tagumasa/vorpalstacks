package dynamodb

import (
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newContributorTestStore(t *testing.T) (*ContributorStore, storage.Storage) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewContributorStore(st, "us-east-1"), st
}

func TestContributorStoreTopKeysWindowAndOrder(t *testing.T) {
	store, st := newContributorTestStore(t)
	base := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)

	record := func(table, keyStr string, at time.Time, units float64) {
		t.Helper()
		if err := st.Update(t.Context(), func(txn storage.Transaction) error {
			return RecordAccessTxn(txn, "us-east-1", table, ContributorLayoutPartitionKey, keyStr, at, 1, units)
		}); err != nil {
			t.Fatalf("record: %v", err)
		}
	}

	hotKey := `["s:hot"]`
	coldKey := `["s:cold"]`
	for i := 0; i < 5; i++ {
		record("Tbl", hotKey, base.Add(time.Duration(i)*time.Minute), 3)
	}
	record("Tbl", coldKey, base, 1)
	// Outside the report window.
	record("Tbl", hotKey, base.Add(-2*time.Hour), 3)

	stats, err := store.TopKeys("Tbl", ContributorLayoutPartitionKey, base.Add(-time.Minute), base.Add(6*time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys: %v", err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 keys in window, got %d", len(stats))
	}
	if stats[0].Key != hotKey {
		t.Fatalf("expected the hot key first, got %q", stats[0].Key)
	}
	if stats[0].Count != 5 || stats[0].Units != 15 {
		t.Fatalf("expected 5 accesses / 15 units, got %d/%.0f", stats[0].Count, stats[0].Units)
	}

	limited, err := store.TopKeys("Tbl", ContributorLayoutPartitionKey, base.Add(-time.Minute), base.Add(6*time.Minute), 1)
	if err != nil {
		t.Fatalf("top keys limited: %v", err)
	}
	if len(limited) != 1 || limited[0].Key != hotKey {
		t.Fatalf("expected only the hot key, got %v", limited)
	}
}

func TestContributorStoreSweepTableOlderThan(t *testing.T) {
	store, st := newContributorTestStore(t)
	base := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)

	record := func(table, keyStr string, at time.Time, units float64) {
		t.Helper()
		if err := st.Update(t.Context(), func(txn storage.Transaction) error {
			return RecordAccessTxn(txn, "us-east-1", table, ContributorLayoutPartitionKey, keyStr, at, 1, units)
		}); err != nil {
			t.Fatalf("record: %v", err)
		}
	}

	record("Tbl", `["s:a"]`, base, 1)
	fresh := base.Add(2 * time.Hour)
	record("Tbl", `["s:b"]`, fresh, 1)
	record("Other", `["s:a"]`, base, 1)

	if err := store.SweepTableOlderThan("Tbl", base.Add(time.Hour)); err != nil {
		t.Fatalf("sweep: %v", err)
	}
	stats, err := store.TopKeys("Tbl", ContributorLayoutPartitionKey, base.Add(-time.Hour), base.Add(3*time.Hour), 10)
	if err != nil {
		t.Fatalf("top keys after sweep: %v", err)
	}
	if len(stats) != 1 || stats[0].Key != `["s:b"]` {
		t.Fatalf("expected only the fresh counter on Tbl, got %v", stats)
	}
	// The sweep is per table: the other table keeps its counter.
	other, err := store.TopKeys("Other", ContributorLayoutPartitionKey, base.Add(-time.Hour), base.Add(3*time.Hour), 10)
	if err != nil {
		t.Fatalf("top keys other: %v", err)
	}
	if len(other) != 1 {
		t.Fatalf("expected the other table's counter to survive, got %v", other)
	}
}

func TestContributorLayoutsAndKeyString(t *testing.T) {
	pk := "pk"
	sk := "sk"

	partitionOnly := &Table{Name: "A", KeySchema: []*KeySchemaElement{{AttributeName: pk, KeyType: KeyTypeHash}}}
	if layouts := ContributorLayouts(partitionOnly); len(layouts) != 1 || layouts[0] != ContributorLayoutPartitionKey {
		t.Fatalf("expected [PKC] for a hash-only table, got %v", layouts)
	}

	composite := &Table{Name: "B", KeySchema: []*KeySchemaElement{{AttributeName: pk, KeyType: KeyTypeHash}, {AttributeName: sk, KeyType: KeyTypeRange}}}
	if layouts := ContributorLayouts(composite); len(layouts) != 2 || layouts[1] != ContributorLayoutFullKey {
		t.Fatalf("expected [PKC SKC] for a composite table, got %v", layouts)
	}

	s := "v"
	n := "7"
	key := map[string]*AttributeValue{"pk": {S: &s}, "sk": {N: &n}}
	if got := ContributorKeyString(composite, key, ContributorLayoutPartitionKey); got != `["s:v"]` {
		t.Fatalf("expected partition-only key [\"s:v\"], got %s", got)
	}
	if got := ContributorKeyString(composite, key, ContributorLayoutFullKey); got != `["s:v","n:7"]` {
		t.Fatalf("expected full key [\"s:v\",\"n:7\"], got %s", got)
	}
}

func TestRecordAccessTxnCommitsWithTransaction(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	at := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	if err := st.Update(t.Context(), func(txn storage.Transaction) error {
		return RecordAccessTxn(txn, "us-east-1", "Tbl", ContributorLayoutPartitionKey, `["s:x"]`, at, 1, 3.0)
	}); err != nil {
		t.Fatalf("record in txn: %v", err)
	}

	store := NewContributorStore(st, "us-east-1")
	stats, err := store.TopKeys("Tbl", ContributorLayoutPartitionKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys: %v", err)
	}
	if len(stats) != 1 || stats[0].Units != 3 {
		t.Fatalf("expected the transactional counter (3 units), got %v", stats)
	}
}

func TestRecordContributorReadsAtomicity(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(
		"Tbl",
		[]*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
		[]*AttributeDefinition{{AttributeName: "id", AttributeType: ScalarAttributeTypeS}},
		BillingModePayPerRequest, nil, nil, nil, nil, nil, false,
	); err != nil {
		t.Fatalf("create table: %v", err)
	}
	tbl, err := store.Tables().Get("Tbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	tbl.ContributorInsightsEnabled = true
	if err := store.Tables().Put(tbl); err != nil {
		t.Fatalf("enable insights: %v", err)
	}

	const goroutines = 40
	const perGoroutine = 25
	key := map[string]*AttributeValue{"id": strAttr("k")}
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				if err := store.RecordContributorReads(t.Context(), "Tbl", []map[string]*AttributeValue{key}); err != nil {
					t.Errorf("record reads: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	at := time.Now()
	stats, err := store.Contributors().TopKeys("Tbl", ContributorLayoutPartitionKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys: %v", err)
	}
	if len(stats) != 1 || stats[0].Count != goroutines*perGoroutine || stats[0].Units != goroutines*perGoroutine {
		t.Fatalf("expected %d atomically counted read events, got %+v", goroutines*perGoroutine, stats)
	}
}
