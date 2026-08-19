package dynamodb

import (
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func TestRecordContributorQuerySingleEvent(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(
		"Tbl",
		[]*KeySchemaElement{{AttributeName: "pk", KeyType: KeyTypeHash}, {AttributeName: "sk", KeyType: KeyTypeRange}},
		[]*AttributeDefinition{{AttributeName: "pk", AttributeType: ScalarAttributeTypeS}, {AttributeName: "sk", AttributeType: ScalarAttributeTypeS}},
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

	key := map[string]*AttributeValue{"pk": strAttr("A")}
	if err := store.RecordContributorQuery(t.Context(), "Tbl", key); err != nil {
		t.Fatalf("record query: %v", err)
	}

	at := time.Now()
	stats, err := store.Contributors().TopKeys("Tbl", ContributorLayoutPartitionKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys: %v", err)
	}
	if len(stats) != 1 || stats[0].Count != 1 || stats[0].Units != ContributorReadUnits {
		t.Fatalf("expected one query event on the partition series, got %+v", stats)
	}
	full, err := store.Contributors().TopKeys("Tbl", ContributorLayoutFullKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys full: %v", err)
	}
	if len(full) != 0 {
		t.Fatalf("a query must not touch the full-key series, got %+v", full)
	}
}

// TestRecordContributorReadsAggregatesSameKey hammers one partition with
// two same-partition items in a single read batch: both reads must be
// counted although they target the same counter key inside one
// transaction.
func TestRecordContributorReadsAggregatesSameKey(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(
		"Tbl",
		[]*KeySchemaElement{{AttributeName: "pk", KeyType: KeyTypeHash}, {AttributeName: "sk", KeyType: KeyTypeRange}},
		[]*AttributeDefinition{{AttributeName: "pk", AttributeType: ScalarAttributeTypeS}, {AttributeName: "sk", AttributeType: ScalarAttributeTypeS}},
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

	keys := []map[string]*AttributeValue{
		{"pk": strAttr("A"), "sk": strAttr("a1")},
		{"pk": strAttr("A"), "sk": strAttr("a2")},
	}
	if err := store.RecordContributorReads(t.Context(), "Tbl", keys); err != nil {
		t.Fatalf("record reads: %v", err)
	}

	at := time.Now()
	stats, err := store.Contributors().TopKeys("Tbl", ContributorLayoutPartitionKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys: %v", err)
	}
	if len(stats) != 1 || stats[0].Key != `["s:A"]` || stats[0].Count != 2 || stats[0].Units != 2*ContributorReadUnits {
		t.Fatalf("expected both same-partition reads counted on the partition series, got %+v", stats)
	}
}

// TestDeleteTableCascadeDropsContributorCounters verifies the counter
// bucket is emptied with the table, so a same-name recreation starts from
// a clean aggregation instead of inheriting the dropped table's counts.
func TestDeleteTableCascadeDropsContributorCounters(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(
		"DropTbl",
		[]*KeySchemaElement{{AttributeName: "pk", KeyType: KeyTypeHash}},
		[]*AttributeDefinition{{AttributeName: "pk", AttributeType: ScalarAttributeTypeS}},
		BillingModePayPerRequest, nil, nil, nil, nil, nil, false,
	); err != nil {
		t.Fatalf("create table: %v", err)
	}
	tbl, err := store.Tables().Get("DropTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	tbl.ContributorInsightsEnabled = true
	if err := store.Tables().Put(tbl); err != nil {
		t.Fatalf("enable insights: %v", err)
	}

	key := map[string]*AttributeValue{"pk": strAttr("k")}
	if err := store.RecordContributorReads(t.Context(), "DropTbl", []map[string]*AttributeValue{key}); err != nil {
		t.Fatalf("record reads: %v", err)
	}
	at := time.Now()
	stats, err := store.Contributors().TopKeys("DropTbl", ContributorLayoutPartitionKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected the recorded counter before the drop, got %+v", stats)
	}

	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteTableCascade("DropTbl")
	}); err != nil {
		t.Fatalf("cascade delete: %v", err)
	}
	stats, err = store.Contributors().TopKeys("DropTbl", ContributorLayoutPartitionKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys after drop: %v", err)
	}
	if len(stats) != 0 {
		t.Fatalf("expected no surviving counters after the drop, got %+v", stats)
	}
}

// TestContributorWriteFlushAtomicity hammers item writes on an
// insights-enabled table from concurrent transactions: every write must be
// counted exactly once in the access counters.
func TestContributorWriteFlushAtomicity(t *testing.T) {
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

	const goroutines = 20
	const perGoroutine = 25
	key := map[string]*AttributeValue{"id": strAttr("k")}
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
					return txn.PutItem("Tbl", key, map[string]*AttributeValue{"v": strAttr("1")})
				}); err != nil {
					t.Errorf("put item: %v", err)
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
	total := goroutines * perGoroutine
	if len(stats) != 1 || stats[0].Count != int64(total) || stats[0].Units != float64(total)*ContributorWriteUnits {
		t.Fatalf("expected %d counted writes (%.0f units), got %+v", total, float64(total)*ContributorWriteUnits, stats)
	}
}
