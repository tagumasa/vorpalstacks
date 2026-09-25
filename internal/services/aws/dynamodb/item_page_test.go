package dynamodb

import (
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func compositeTable() *dbstore.Table {
	return &dbstore.Table{
		Name: "T",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "gsi_pk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexes: []*dbstore.GlobalSecondaryIndex{
			{
				IndexName: "gsi",
				KeySchema: []*dbstore.KeySchemaElement{
					{AttributeName: "gsi_pk", KeyType: dbstore.KeyTypeHash},
					{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
				},
			},
		},
	}
}

func TestPrimaryKeyFromStartKey(t *testing.T) {
	table := compositeTable()

	full := map[string]*dbstore.AttributeValue{
		"pk":     dbstore.StringValue("p1"),
		"sk":     dbstore.StringValue("s1"),
		"gsi_pk": dbstore.StringValue("g1"),
	}
	pk := primaryKeyFromStartKey(table, full)
	if pk == nil || len(pk) != 2 || pk["pk"] == nil || pk["sk"] == nil {
		t.Fatalf("complete start key: expected both primary key attributes, got %v", pk)
	}

	// Dropping any primary key member makes the start key unable to anchor
	// pagination at an item's storage position.
	for _, missing := range []string{"pk", "sk"} {
		partial := map[string]*dbstore.AttributeValue{}
		for k, v := range full {
			partial[k] = v
		}
		delete(partial, missing)
		if got := primaryKeyFromStartKey(table, partial); got != nil {
			t.Fatalf("start key without %s: expected nil, got %v", missing, got)
		}
	}
}

func TestIndexMarkerFromStartKeyMirrorsIndexKeyComposition(t *testing.T) {
	table := compositeTable()

	esk := map[string]*dbstore.AttributeValue{
		"pk":     dbstore.StringValue("p1"),
		"sk":     dbstore.StringValue("s1"),
		"gsi_pk": dbstore.StringValue("g1"),
	}
	encodedHash := dbstore.EncodeKeyValue(dbstore.StringValue("g1"))
	// The expected marker is the key the index store itself composes for
	// the same entry — the resumer must anchor at the writer's key, not a
	// second spelling of the layout.
	entry := &dbstore.Item{
		TableName: table.Name,
		Key: map[string]*dbstore.AttributeValue{
			"pk": dbstore.StringValue("p1"),
			"sk": dbstore.StringValue("s1"),
		},
		Attributes: map[string]*dbstore.AttributeValue{
			"gsi_pk": dbstore.StringValue("g1"),
		},
	}
	want := dbstore.NewIndexStore("us-east-1").BuildGSIKey(table, table.GlobalSecondaryIndexes[0], entry)

	if got := indexMarkerFromStartKey(table, "gsi", encodedHash, esk); got != want {
		t.Fatalf("index marker:\n got %q\nwant %q", got, want)
	}

	// A start key without the index sort key cannot name one index entry.
	noSort := map[string]*dbstore.AttributeValue{"pk": dbstore.StringValue("p1"), "gsi_pk": dbstore.StringValue("g1")}
	if got := indexMarkerFromStartKey(table, "gsi", encodedHash, noSort); got != "" {
		t.Fatalf("start key without sort key: expected empty marker, got %q", got)
	}
	if got := indexMarkerFromStartKey(table, "gsi", encodedHash, nil); got != "" {
		t.Fatalf("nil start key: expected empty marker, got %q", got)
	}
}

func TestSortKeyConditionMatches(t *testing.T) {
	item := &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{
		"sk": dbstore.NumberValue("10"),
	}}

	if !sortKeyConditionMatches(item, nil) {
		t.Fatalf("nil condition must match every item")
	}

	eq := &sortKeyCondition{attrName: "sk", op: "=", value: dbstore.NumberValue("10")}
	lt := &sortKeyCondition{attrName: "sk", op: "<", value: dbstore.NumberValue("20")}
	gt := &sortKeyCondition{attrName: "sk", op: ">", value: dbstore.NumberValue("20")}
	if !sortKeyConditionMatches(item, eq) || !sortKeyConditionMatches(item, lt) {
		t.Fatalf("sk=10 must match =10 and <20")
	}
	if sortKeyConditionMatches(item, gt) {
		t.Fatalf("sk=10 must not match >20")
	}

	// An item without the sort-key attribute never matches.
	other := &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{"x": dbstore.StringValue("v")}}
	if sortKeyConditionMatches(other, eq) {
		t.Fatalf("item without sort key must not match")
	}
}
