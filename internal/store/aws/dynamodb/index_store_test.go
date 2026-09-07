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

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
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
