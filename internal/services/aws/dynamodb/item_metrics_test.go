package dynamodb

import (
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// metricsTable builds a table record with a string partition key and, when
// lsi is set, one local secondary index, so the item collection gating can
// be exercised without a live store.
func metricsTable(name string, lsi bool) *dbstore.Table {
	tbl := &dbstore.Table{
		Name: name,
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
	}
	if lsi {
		tbl.LocalSecondaryIndexes = []*dbstore.LocalSecondaryIndex{
			{IndexName: "lsi1"},
		}
	}
	return tbl
}

func metricsKey(pk, sk string) map[string]*dbstore.AttributeValue {
	return map[string]*dbstore.AttributeValue{
		"pk": {S: &pk},
		"sk": {S: &sk},
	}
}

// assertMetricsEntry verifies the documented single-entry shape: the
// partition key as the ItemCollectionKey and a two-element size estimate
// range.
func assertMetricsEntry(t *testing.T, entry map[string]interface{}, wantPK string) {
	t.Helper()
	key, ok := entry["ItemCollectionKey"].(map[string]interface{})
	if !ok {
		t.Fatalf("ItemCollectionKey is not a map: %T", entry["ItemCollectionKey"])
	}
	pk, ok := key["pk"].(map[string]interface{})
	if !ok {
		t.Fatalf("ItemCollectionKey[pk] is not a wire attribute: %T", key["pk"])
	}
	if pk["S"] != wantPK {
		t.Fatalf("ItemCollectionKey[pk] = %v, want %q", pk["S"], wantPK)
	}
	rangeGB, ok := entry["SizeEstimateRangeGB"].([]float64)
	if !ok || len(rangeGB) != 2 {
		t.Fatalf("SizeEstimateRangeGB is not a two-element range: %T %v", entry["SizeEstimateRangeGB"], entry["SizeEstimateRangeGB"])
	}
}

func TestBuildItemCollectionMetricsEntry(t *testing.T) {
	table := metricsTable("LSITable", true)
	entry := buildItemCollectionMetricsEntry(table, metricsKey("p1", "a"))
	if entry == nil {
		t.Fatal("expected an entry for an LSI table write")
	}
	assertMetricsEntry(t, entry, "p1")

	if got := buildItemCollectionMetricsEntry(metricsTable("PlainTable", false), metricsKey("p1", "a")); got != nil {
		t.Fatalf("expected no entry for a table without local secondary indexes, got %v", got)
	}
	if got := buildItemCollectionMetricsEntry(nil, metricsKey("p1", "a")); got != nil {
		t.Fatalf("expected no entry for a nil table, got %v", got)
	}
	if got := buildItemCollectionMetricsEntry(table, map[string]*dbstore.AttributeValue{}); got != nil {
		t.Fatalf("expected no entry for a key without the partition attribute, got %v", got)
	}
}

func TestBuildItemCollectionMetricsPerTable(t *testing.T) {
	lsiTable := metricsTable("LSITable", true)
	plainTable := metricsTable("PlainTable", false)

	perTable := buildItemCollectionMetricsPerTable([]itemCollectionWriteRef{
		{tableName: "LSITable", table: lsiTable, key: metricsKey("p1", "a")},
		{tableName: "LSITable", table: lsiTable, key: metricsKey("p1", "b")},
		{tableName: "LSITable", table: lsiTable, key: metricsKey("p2", "a")},
		{tableName: "PlainTable", table: plainTable, key: metricsKey("p1", "a")},
	})

	entries, ok := perTable["LSITable"].([]map[string]interface{})
	if !ok {
		t.Fatalf("LSITable entries missing or wrong type: %T %v", perTable["LSITable"], perTable["LSITable"])
	}
	if len(entries) != 2 {
		t.Fatalf("expected one entry per distinct partition key value, got %d", len(entries))
	}
	assertMetricsEntry(t, entries[0], "p1")
	assertMetricsEntry(t, entries[1], "p2")
	if _, present := perTable["PlainTable"]; present {
		t.Fatalf("expected no entries for a table without local secondary indexes, got %v", perTable["PlainTable"])
	}

	if got := buildItemCollectionMetricsPerTable([]itemCollectionWriteRef{
		{tableName: "PlainTable", table: plainTable, key: metricsKey("p1", "a")},
	}); got != nil {
		t.Fatalf("expected nil when no write touches an LSI table, got %v", got)
	}
}
