package dynamodb

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newIndexStoreFixture creates a store with a table carrying two GSIs (gsi-a
// on attribute "a", gsi-b on attribute "b") for index-entry lifecycle tests.
func newIndexStoreFixture(t *testing.T) *DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	gsiA := &GlobalSecondaryIndex{
		IndexName:  "gsi-a",
		KeySchema:  []*KeySchemaElement{{AttributeName: "a", KeyType: KeyTypeHash}},
		Projection: &Projection{ProjectionType: "ALL"},
	}
	gsiB := &GlobalSecondaryIndex{
		IndexName:  "gsi-b",
		KeySchema:  []*KeySchemaElement{{AttributeName: "b", KeyType: KeyTypeHash}},
		Projection: &Projection{ProjectionType: "ALL"},
	}
	if _, err := store.Tables().Create(CreateTableParams{
		Name:      "IdxTbl",
		KeySchema: []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
		AttributeDefinitions: []*AttributeDefinition{
			{AttributeName: "id", AttributeType: ScalarAttributeTypeS},
			{AttributeName: "a", AttributeType: ScalarAttributeTypeS},
			{AttributeName: "b", AttributeType: ScalarAttributeTypeS},
		},
		BillingMode:            BillingModePayPerRequest,
		GlobalSecondaryIndexes: []*GlobalSecondaryIndex{gsiA, gsiB},
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return store
}

// putIndexFixtureItem writes one item carrying both index key attributes and
// registers its index entries.
func putIndexFixtureItem(t *testing.T, store *DynamoDBStore, id string) {
	t.Helper()
	key := map[string]*AttributeValue{"id": strAttr(id)}
	attrs := map[string]*AttributeValue{"a": strAttr("va"), "b": strAttr("vb")}
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		if err := txn.PutItem("IdxTbl", key, attrs); err != nil {
			return err
		}
		return txn.PutIndexEntries("IdxTbl", &Item{TableName: "IdxTbl", Key: key, Attributes: attrs})
	}); err != nil {
		t.Fatalf("put fixture item %s: %v", id, err)
	}
}

func countIndexEntries(t *testing.T, store *DynamoDBStore, indexName, hashRaw string) int {
	t.Helper()
	count := 0
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		items, err := txn.QueryByGSI("IdxTbl", indexName, EncodeKeyValue(strAttr(hashRaw)), IndexQueryOptions{})
		if err != nil {
			return err
		}
		count = len(items)
		return nil
	}); err != nil {
		t.Fatalf("query %s: %v", indexName, err)
	}
	return count
}

func TestIndexStoreDeleteIndexEntriesForIndex(t *testing.T) {
	store := newIndexStoreFixture(t)
	putIndexFixtureItem(t, store, "k1")
	putIndexFixtureItem(t, store, "k2")

	if got := countIndexEntries(t, store, "gsi-a", "va"); got != 2 {
		t.Fatalf("gsi-a before cleanup = %d entries, want 2", got)
	}
	if got := countIndexEntries(t, store, "gsi-b", "vb"); got != 2 {
		t.Fatalf("gsi-b before cleanup = %d entries, want 2", got)
	}

	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteIndexEntriesForIndex("IdxTbl", "gsi-a")
	}); err != nil {
		t.Fatalf("delete entries for gsi-a: %v", err)
	}

	// Only the deleted index's entries are removed; the sibling index is
	// untouched.
	if got := countIndexEntries(t, store, "gsi-a", "va"); got != 0 {
		t.Fatalf("gsi-a after cleanup = %d entries, want 0", got)
	}
	if got := countIndexEntries(t, store, "gsi-b", "vb"); got != 2 {
		t.Fatalf("gsi-b after cleanup = %d entries, want 2", got)
	}
}

func TestIndexStorePutIndexEntriesForIndexWritesOnlyNamedIndex(t *testing.T) {
	store := newIndexStoreFixture(t)

	// The item is stored but only gsi-a's entry is written; a name missing
	// from the schema writes nothing.
	key := map[string]*AttributeValue{"id": strAttr("k3")}
	attrs := map[string]*AttributeValue{"a": strAttr("va"), "b": strAttr("vb")}
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		if err := txn.PutItem("IdxTbl", key, attrs); err != nil {
			return err
		}
		if err := txn.PutGSIEntriesForIndex("IdxTbl", "gsi-a", &Item{TableName: "IdxTbl", Key: key, Attributes: attrs}); err != nil {
			return err
		}
		return txn.PutGSIEntriesForIndex("IdxTbl", "gsi-absent", &Item{TableName: "IdxTbl", Key: key, Attributes: attrs})
	}); err != nil {
		t.Fatalf("put named index entry: %v", err)
	}

	if got := countIndexEntries(t, store, "gsi-a", "va"); got != 1 {
		t.Fatalf("gsi-a = %d entries, want 1", got)
	}
	if got := countIndexEntries(t, store, "gsi-b", "vb"); got != 0 {
		t.Fatalf("gsi-b must not gain an entry from the named write, got %d", got)
	}
}

// TestIndexStoreLSIEntryLifecycle pins the LSI half of the entry pair:
// PutIndexEntries writes the LSI entry a query resolves, DeleteIndexEntries
// removes it, and an item lacking the LSI sort-key attribute contributes
// no entry.
func TestIndexStoreLSIEntryLifecycle(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(CreateTableParams{
		Name:      "IdxLsiTbl",
		KeySchema: []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
		AttributeDefinitions: []*AttributeDefinition{
			{AttributeName: "id", AttributeType: ScalarAttributeTypeS},
			{AttributeName: "a", AttributeType: ScalarAttributeTypeS},
		},
		BillingMode: BillingModePayPerRequest,
		LocalSecondaryIndexes: []*LocalSecondaryIndex{{
			IndexName:  "lsi-a",
			KeySchema:  []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}, {AttributeName: "a", KeyType: KeyTypeRange}},
			Projection: &Projection{ProjectionType: "ALL"},
		}},
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	put := func(id, sortVal string) {
		t.Helper()
		key := map[string]*AttributeValue{"id": strAttr(id)}
		attrs := map[string]*AttributeValue{}
		if sortVal != "" {
			attrs["a"] = strAttr(sortVal)
		}
		if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
			if err := txn.PutItem("IdxLsiTbl", key, attrs); err != nil {
				return err
			}
			return txn.PutIndexEntries("IdxLsiTbl", &Item{TableName: "IdxLsiTbl", Key: key, Attributes: attrs})
		}); err != nil {
			t.Fatalf("put %s: %v", id, err)
		}
	}
	// The LSI hash key is the table hash key, so the extent is counted per
	// hash value.
	count := func(hashRaw string) int {
		t.Helper()
		n := 0
		if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
			items, err := txn.QueryByLSI("IdxLsiTbl", "lsi-a", EncodeKeyValue(strAttr(hashRaw)), IndexQueryOptions{})
			if err != nil {
				return err
			}
			n = len(items)
			return nil
		}); err != nil {
			t.Fatalf("query lsi: %v", err)
		}
		return n
	}

	put("k1", "va")
	put("k2", "vb")
	if got := count("k1") + count("k2"); got != 2 {
		t.Fatalf("lsi entries after puts = %d, want 2", got)
	}

	// An item without the sort-key attribute holds no LSI entry.
	put("k3", "")
	if got := count("k3"); got != 0 {
		t.Fatalf("sort-key-less item gained an lsi entry = %d, want 0", got)
	}

	// Deleting one item's entries removes exactly its own.
	key1 := map[string]*AttributeValue{"id": strAttr("k1")}
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteIndexEntries("IdxLsiTbl", &Item{TableName: "IdxLsiTbl", Key: key1, Attributes: map[string]*AttributeValue{"a": strAttr("va")}})
	}); err != nil {
		t.Fatalf("delete entries: %v", err)
	}
	if got := count("k1"); got != 0 {
		t.Fatalf("lsi entries of k1 after delete = %d, want 0", got)
	}
	if got := count("k2"); got != 1 {
		t.Fatalf("lsi entries of k2 after delete = %d, want 1", got)
	}
}
