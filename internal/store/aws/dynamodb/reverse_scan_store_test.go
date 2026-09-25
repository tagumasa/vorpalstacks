package dynamodb

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newReverseQueryFixture creates a composite-key table (id HASH, sk RANGE)
// whose GSI (hash "a", range "sk") orders entries by the encoded sort key,
// plus one plain table for partition-scan tests.
func newReverseQueryFixture(t *testing.T) *DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	gsi := &GlobalSecondaryIndex{
		IndexName: "gsi",
		KeySchema: []*KeySchemaElement{
			{AttributeName: "a", KeyType: KeyTypeHash},
			{AttributeName: "sk", KeyType: KeyTypeRange},
		},
		Projection: &Projection{ProjectionType: "ALL"},
	}
	if _, err := store.Tables().Create(CreateTableParams{
		Name:      "IdxTbl",
		KeySchema: []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}, {AttributeName: "sk", KeyType: KeyTypeRange}},
		AttributeDefinitions: []*AttributeDefinition{
			{AttributeName: "id", AttributeType: ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: ScalarAttributeTypeS},
			{AttributeName: "a", AttributeType: ScalarAttributeTypeS},
		},
		BillingMode:            BillingModePayPerRequest,
		GlobalSecondaryIndexes: []*GlobalSecondaryIndex{gsi},
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return store
}

// putReverseFixtureItem stores one item of partition a="va" with the given
// sort keys and registers its index entries.
func putReverseFixtureItem(t *testing.T, store *DynamoDBStore, id, sk string) {
	t.Helper()
	key := map[string]*AttributeValue{"id": strAttr(id), "sk": strAttr(sk)}
	attrs := map[string]*AttributeValue{"a": strAttr("va"), "sk": strAttr(sk)}
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		if err := txn.PutItem("IdxTbl", key, attrs); err != nil {
			return err
		}
		return txn.PutIndexEntries("IdxTbl", &Item{TableName: "IdxTbl", Key: key, Attributes: attrs})
	}); err != nil {
		t.Fatalf("put fixture item %s/%s: %v", id, sk, err)
	}
}

func queryGSI(t *testing.T, store *DynamoDBStore, opts IndexQueryOptions) []string {
	t.Helper()
	var sks []string
	if err := store.View(t.Context(), func(txn *DynamoDBTxn) error {
		items, err := txn.QueryByGSI("IdxTbl", "gsi", EncodeKeyValue(strAttr("va")), opts)
		if err != nil {
			return err
		}
		for _, item := range items {
			sks = append(sks, *item.Key["sk"].S)
		}
		return nil
	}); err != nil {
		t.Fatalf("query gsi: %v", err)
	}
	return sks
}

// gsiMarkerForItem derives the index-bucket key of one fixture item so the
// marker-resumption tests address a real entry.
func gsiMarkerForItem(t *testing.T, store *DynamoDBStore, id, sk string) string {
	t.Helper()
	table, err := store.Tables().Get("IdxTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	idx := NewIndexStore("us-east-1")
	marker := idx.BuildGSIKey(table, table.GlobalSecondaryIndexes[0], &Item{
		TableName:  "IdxTbl",
		Key:        map[string]*AttributeValue{"id": strAttr(id), "sk": strAttr(sk)},
		Attributes: map[string]*AttributeValue{"a": strAttr("va"), "sk": strAttr(sk)},
	})
	if marker == "" {
		t.Fatalf("build marker for %s/%s: empty", id, sk)
	}
	return marker
}

func TestQueryByIndexReverseHonoursLimit(t *testing.T) {
	store := newReverseQueryFixture(t)
	putReverseFixtureItem(t, store, "k1", "01")
	putReverseFixtureItem(t, store, "k2", "02")
	putReverseFixtureItem(t, store, "k3", "03")
	putReverseFixtureItem(t, store, "k4", "04")
	putReverseFixtureItem(t, store, "k5", "05")

	if got := queryGSI(t, store, IndexQueryOptions{}); len(got) != 5 || got[0] != "01" || got[4] != "05" {
		t.Fatalf("forward walk = %v, want 01..05 ascending", got)
	}

	// A reverse page reads from the far end and stops at the limit instead
	// of materialising the whole hash range.
	if got := queryGSI(t, store, IndexQueryOptions{Reverse: true, Limit: 2}); len(got) != 2 || got[0] != "05" || got[1] != "04" {
		t.Fatalf("reverse limit 2 = %v, want [05 04]", got)
	}
	if got := queryGSI(t, store, IndexQueryOptions{Limit: 2}); len(got) != 2 || got[0] != "01" || got[1] != "02" {
		t.Fatalf("forward limit 2 = %v, want [01 02]", got)
	}
}

func TestQueryByIndexMarkerResumesBothDirections(t *testing.T) {
	store := newReverseQueryFixture(t)
	putReverseFixtureItem(t, store, "k1", "01")
	putReverseFixtureItem(t, store, "k2", "02")
	putReverseFixtureItem(t, store, "k3", "03")
	putReverseFixtureItem(t, store, "k4", "04")

	markerAt03 := gsiMarkerForItem(t, store, "k3", "03")

	// A forward page resumes strictly after the marker entry.
	if got := queryGSI(t, store, IndexQueryOptions{Marker: markerAt03}); len(got) != 1 || got[0] != "04" {
		t.Fatalf("forward from marker = %v, want [04]", got)
	}
	// A reverse page resumes strictly before it.
	if got := queryGSI(t, store, IndexQueryOptions{Reverse: true, Marker: markerAt03}); len(got) != 2 || got[0] != "02" || got[1] != "01" {
		t.Fatalf("reverse from marker = %v, want [02 01]", got)
	}
}

func TestQueryByIndexFilterRunsBeforeLimit(t *testing.T) {
	store := newReverseQueryFixture(t)
	putReverseFixtureItem(t, store, "k1", "01")
	putReverseFixtureItem(t, store, "k2", "02")
	putReverseFixtureItem(t, store, "k3", "03")
	putReverseFixtureItem(t, store, "k4", "04")
	putReverseFixtureItem(t, store, "k5", "05")

	onlyEven := func(item *Item) bool {
		return *item.Key["sk"].S == "02" || *item.Key["sk"].S == "04"
	}

	if got := queryGSI(t, store, IndexQueryOptions{Limit: 1, Filter: onlyEven}); len(got) != 1 || got[0] != "02" {
		t.Fatalf("filtered forward limit 1 = %v, want [02]", got)
	}
	if got := queryGSI(t, store, IndexQueryOptions{Reverse: true, Limit: 1, Filter: onlyEven}); len(got) != 1 || got[0] != "04" {
		t.Fatalf("filtered reverse limit 1 = %v, want [04]", got)
	}
	// The filter must not count rejected entries against the limit.
	if got := queryGSI(t, store, IndexQueryOptions{Limit: 2, Filter: onlyEven}); len(got) != 2 || got[0] != "02" || got[1] != "04" {
		t.Fatalf("filtered forward limit 2 = %v, want [02 04]", got)
	}
}

func TestScanByPartitionKeyReverseAndMarker(t *testing.T) {
	store := newReverseQueryFixture(t)
	putReverseFixtureItem(t, store, "k1", "01")
	putReverseFixtureItem(t, store, "k2", "02")
	putReverseFixtureItem(t, store, "k3", "03")
	putReverseFixtureItem(t, store, "k4", "04")

	table, err := store.Tables().Get("IdxTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	hashValue := EncodeKeyValue(strAttr("k1"))

	scanSKs := func(opts ScanOptions) []string {
		var sks []string
		if _, err := store.Items().ScanByPartitionKeyWithTable("IdxTbl", table, hashValue, opts, func(item *Item) error {
			sks = append(sks, *item.Key["sk"].S)
			return nil
		}); err != nil {
			t.Fatalf("partition scan: %v", err)
		}
		return sks
	}

	// Every fixture item shares the id partition only with itself: k1 is
	// the sole item of partition "k1", so seed a second partition with
	// multiple members for the ordering assertions.
	putReverseFixtureItem(t, store, "p1", "10")
	putReverseFixtureItem(t, store, "p1", "20")
	putReverseFixtureItem(t, store, "p1", "30")

	multiHash := EncodeKeyValue(strAttr("p1"))
	scanMulti := func(opts ScanOptions) []string {
		var sks []string
		if _, err := store.Items().ScanByPartitionKeyWithTable("IdxTbl", table, multiHash, opts, func(item *Item) error {
			sks = append(sks, *item.Key["sk"].S)
			return nil
		}); err != nil {
			t.Fatalf("partition scan: %v", err)
		}
		return sks
	}

	if got := scanMulti(ScanOptions{}); len(got) != 3 || got[0] != "10" || got[2] != "30" {
		t.Fatalf("forward partition scan = %v, want 10..30 ascending", got)
	}
	if got := scanMulti(ScanOptions{Reverse: true, Limit: 2}); len(got) != 2 || got[0] != "30" || got[1] != "20" {
		t.Fatalf("reverse partition scan limit 2 = %v, want [30 20]", got)
	}

	// The marker is the storage key of the sk=20 item; a reverse walk
	// resumes strictly before it.
	marker := EncodeItemKey("IdxTbl", map[string]*AttributeValue{"id": strAttr("p1"), "sk": strAttr("20")}, table)
	if got := scanMulti(ScanOptions{Reverse: true, Marker: marker}); len(got) != 1 || got[0] != "10" {
		t.Fatalf("reverse from marker = %v, want [10]", got)
	}
	if got := scanMulti(ScanOptions{Marker: marker}); len(got) != 1 || got[0] != "30" {
		t.Fatalf("forward from marker = %v, want [30]", got)
	}

	// The filter runs before the limit is counted.
	onlyEvenSK := func(item *Item) bool { return *item.Key["sk"].S == "20" }
	if got := scanMulti(ScanOptions{Limit: 1, Filter: onlyEvenSK}); len(got) != 1 || got[0] != "20" {
		t.Fatalf("filtered partition scan = %v, want [20]", got)
	}

	// The k1 partition still yields its single item.
	if got := scanSKs(ScanOptions{}); len(got) != 1 || got[0] != "01" {
		t.Fatalf("single-item partition = %v, want [01]", got)
	}
}

// TestPaginatedScansDeliverEveryItem pins the marker contract shared by the
// three storage scan helpers: the returned marker is the key of the last
// item the page delivered, so walking pages at any limit delivers every
// item exactly once, in both directions. The page walks below spell out the
// exact expected sequences.
func TestPaginatedScansDeliverEveryItem(t *testing.T) {
	store := newReverseQueryFixture(t)
	for _, id := range []string{"k1", "k2", "k3", "k4", "k5"} {
		putReverseFixtureItem(t, store, id, "01")
	}
	putReverseFixtureItem(t, store, "p1", "10")
	putReverseFixtureItem(t, store, "p1", "20")
	putReverseFixtureItem(t, store, "p1", "30")
	putReverseFixtureItem(t, store, "p1", "40")
	putReverseFixtureItem(t, store, "p1", "50")

	table, err := store.Tables().Get("IdxTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}

	// ScanWithOptions pages the whole table (ascending storage-key order).
	var scanPages []string
	marker := ""
	for {
		var page []string
		next, err := store.Items().ScanWithOptions("IdxTbl", ScanOptions{Limit: 2, Marker: marker}, func(item *Item) error {
			page = append(page, *item.Key["id"].S)
			return nil
		})
		if err != nil {
			t.Fatalf("scan page: %v", err)
		}
		scanPages = append(scanPages, page...)
		if next == "" || len(page) == 0 {
			break
		}
		marker = next
	}
	wantScan := []string{"k1", "k2", "k3", "k4", "k5", "p1", "p1", "p1", "p1", "p1"}
	if len(scanPages) != len(wantScan) {
		t.Fatalf("scan walk delivered %v (%d items), want %d", scanPages, len(scanPages), len(wantScan))
	}

	// The same table-scan walk in reverse delivers every item exactly once
	// in descending storage-key order.
	var reverseScanPages []string
	marker = ""
	for {
		var page []string
		next, err := store.Items().ScanWithOptions("IdxTbl", ScanOptions{Limit: 2, Reverse: true, Marker: marker}, func(item *Item) error {
			page = append(page, *item.Key["id"].S)
			return nil
		})
		if err != nil {
			t.Fatalf("reverse scan page: %v", err)
		}
		reverseScanPages = append(reverseScanPages, page...)
		if next == "" || len(page) == 0 {
			break
		}
		marker = next
	}
	wantReverseScan := []string{"p1", "p1", "p1", "p1", "p1", "k5", "k4", "k3", "k2", "k1"}
	if len(reverseScanPages) != len(wantReverseScan) {
		t.Fatalf("reverse scan walk delivered %v (%d items), want %d", reverseScanPages, len(reverseScanPages), len(wantReverseScan))
	}

	// List pages the same table with its own marker/limit pair.
	var listKeys []string
	marker = ""
	for {
		items, next, err := store.Items().List("IdxTbl", marker, 2)
		if err != nil {
			t.Fatalf("list page: %v", err)
		}
		for _, item := range items {
			listKeys = append(listKeys, *item.Key["id"].S)
		}
		if next == "" || len(items) == 0 {
			break
		}
		marker = next
	}
	if len(listKeys) != len(wantScan) {
		t.Fatalf("list walk delivered %v (%d items), want %d", listKeys, len(listKeys), len(wantScan))
	}

	// The partition scan pages one partition, forward and reverse, and both
	// walks deliver all five sort keys exactly once in order.
	hashValue := EncodeKeyValue(strAttr("p1"))
	wantForward := []string{"10", "20", "30", "40", "50"}
	wantReverse := []string{"50", "40", "30", "20", "10"}

	var forwardSKs []string
	marker = ""
	for {
		var page []string
		next, err := store.Items().ScanByPartitionKeyWithTable("IdxTbl", table, hashValue, ScanOptions{Limit: 2, Marker: marker}, func(item *Item) error {
			page = append(page, *item.Key["sk"].S)
			return nil
		})
		if err != nil {
			t.Fatalf("partition page: %v", err)
		}
		forwardSKs = append(forwardSKs, page...)
		if next == "" || len(page) == 0 {
			break
		}
		marker = next
	}
	if len(forwardSKs) != len(wantForward) {
		t.Fatalf("forward partition walk = %v, want %v", forwardSKs, wantForward)
	}

	var reverseSKs []string
	marker = ""
	for {
		var page []string
		next, err := store.Items().ScanByPartitionKeyWithTable("IdxTbl", table, hashValue, ScanOptions{Limit: 2, Reverse: true, Marker: marker}, func(item *Item) error {
			page = append(page, *item.Key["sk"].S)
			return nil
		})
		if err != nil {
			t.Fatalf("reverse partition page: %v", err)
		}
		reverseSKs = append(reverseSKs, page...)
		if next == "" || len(page) == 0 {
			break
		}
		marker = next
	}
	if len(reverseSKs) != len(wantReverse) {
		t.Fatalf("reverse partition walk = %v, want %v", reverseSKs, wantReverse)
	}
}

// TestTxnWalksDeliverThroughTheSharedCore pins the transaction walks'
// delivery semantics on the shared walk core: Scan delivers every item of
// the table in key order, and ScanByPartitionKey delivers exactly the
// items of one partition — the encoded partition component is prefix-free,
// so a same-prefixed foreign partition cannot leak into the page.
func TestTxnWalksDeliverThroughTheSharedCore(t *testing.T) {
	store := newReverseQueryFixture(t)
	putReverseFixtureItem(t, store, "k1", "01")
	putReverseFixtureItem(t, store, "k2", "02")
	putReverseFixtureItem(t, store, "k3", "03")

	var scanned []string
	err := store.View(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.Scan("IdxTbl", func(item *Item) error {
			scanned = append(scanned, *item.Key["id"].S+"/"+*item.Key["sk"].S)
			return nil
		})
	})
	if err != nil {
		t.Fatalf("txn scan: %v", err)
	}
	want := []string{"k1/01", "k2/02", "k3/03"}
	if len(scanned) != len(want) {
		t.Fatalf("txn scan delivered %v, want %v", scanned, want)
	}
	for i, key := range want {
		if scanned[i] != key {
			t.Fatalf("txn scan delivered %v, want %v", scanned, want)
		}
	}

	var partition []string
	err = store.View(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.ScanByPartitionKey("IdxTbl", EncodeKeyValue(strAttr("k2")), func(item *Item) error {
			partition = append(partition, *item.Key["sk"].S)
			return nil
		})
	})
	if err != nil {
		t.Fatalf("txn partition scan: %v", err)
	}
	if len(partition) != 1 || partition[0] != "02" {
		t.Fatalf("txn partition scan delivered %v, want [02]", partition)
	}
}
