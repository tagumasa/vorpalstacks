package dynamodb

import (
	"errors"
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func testARNBuilder() *svcarn.DynamoDBBuilder {
	return svcarn.NewARNBuilder("123456789012", "us-east-1").DynamoDB()
}

func gsiFixture(name string) *dbstore.GlobalSecondaryIndex {
	return &dbstore.GlobalSecondaryIndex{
		IndexName:  name,
		KeySchema:  []*dbstore.KeySchemaElement{{AttributeName: "a", KeyType: dbstore.KeyTypeHash}},
		Projection: &dbstore.Projection{ProjectionType: "ALL"},
	}
}

func gsiCreateUpdate(name string) interface{} {
	return map[string]interface{}{
		"Create": map[string]interface{}{
			"IndexName":  name,
			"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "a", "KeyType": "HASH"}},
			"Projection": map[string]interface{}{"ProjectionType": "ALL"},
		},
	}
}

func gsiDeleteUpdate(name string) interface{} {
	return map[string]interface{}{"Delete": map[string]interface{}{"IndexName": name}}
}

func indexNames(gsis []*dbstore.GlobalSecondaryIndex) []string {
	names := make([]string, 0, len(gsis))
	for _, g := range gsis {
		names = append(names, g.IndexName)
	}
	return names
}

func TestApplyGSIUpdatesReturnsDeletedIndexNames(t *testing.T) {
	existing := []*dbstore.GlobalSecondaryIndex{gsiFixture("gsi-1"), gsiFixture("gsi-2")}

	updated, deleted, err := applyGSIUpdates(testARNBuilder(), "T", existing, []interface{}{gsiDeleteUpdate("gsi-1")})
	if err != nil {
		t.Fatalf("delete update: %v", err)
	}
	if len(deleted) != 1 || deleted[0] != "gsi-1" {
		t.Fatalf("deleted names = %v, want [gsi-1]", deleted)
	}
	if names := indexNames(updated); len(names) != 1 || names[0] != "gsi-2" {
		t.Fatalf("updated indexes = %v, want [gsi-2]", names)
	}

	// A delete followed by a same-name create in one request leaves the
	// index in the final schema, so its entries must survive for the
	// backfill to rebuild on top of.
	updated, deleted, err = applyGSIUpdates(testARNBuilder(), "T", existing, []interface{}{
		gsiDeleteUpdate("gsi-1"),
		gsiCreateUpdate("gsi-1"),
	})
	if err != nil {
		t.Fatalf("delete+create update: %v", err)
	}
	if len(deleted) != 0 {
		t.Fatalf("re-created index must not be reported deleted, got %v", deleted)
	}
	if names := indexNames(updated); len(names) != 2 {
		t.Fatalf("updated indexes = %v, want gsi-1 and gsi-2", names)
	}

	// Deleting an unknown index stays rejected.
	_, _, err = applyGSIUpdates(testARNBuilder(), "T", existing, []interface{}{gsiDeleteUpdate("nope")})
	if !errors.Is(err, ErrIndexNotFound) {
		t.Fatalf("delete of unknown index = %v, want ErrIndexNotFound", err)
	}
}
