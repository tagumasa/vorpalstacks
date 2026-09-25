package dynamodb

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// panicScanItems makes the backfill's item scan panic, exercising the
// synchronous-path panic conversion: the helpers run on the request's
// goroutine, so a panic must surface as the returned error, never as a
// swallowed log entry that leaves UpdateTable succeeding over a
// half-built index.
type panicScanItems struct {
	dbstore.ItemStoreInterface
}

func (panicScanItems) Scan(string, func(*dbstore.Item) error) error {
	panic("scan exploded")
}

type panicBackfillStore struct {
	dbstore.DynamoDBStoreInterface
}

func (panicBackfillStore) Items() dbstore.ItemStoreInterface {
	return panicScanItems{}
}

func TestBackfillPanicsConvertToErrors(t *testing.T) {
	var svc DynamoDBService
	store := panicBackfillStore{}

	gsiErr := svc.backfillGSI(context.Background(), store, "t", "gsi-index")
	if gsiErr == nil {
		t.Fatalf("GSI backfill: panic swallowed — want an error")
	}
	if !strings.Contains(gsiErr.Error(), "panicked") {
		t.Fatalf("GSI backfill: error %q does not name the panic", gsiErr)
	}

	vectorErr := svc.backfillVectorIndex(context.Background(), store, "t", "vec-index")
	if vectorErr == nil {
		t.Fatalf("vector backfill: panic swallowed — want an error")
	}
	if !strings.Contains(vectorErr.Error(), "panicked") {
		t.Fatalf("vector backfill: error %q does not name the panic", vectorErr)
	}
}

// itemsPanicStore embeds the real store so every non-scan operation runs
// against real storage; only the backfill's item scan panics.
type itemsPanicStore struct {
	dbstore.DynamoDBStoreInterface
}

func (itemsPanicStore) Items() dbstore.ItemStoreInterface {
	return panicScanItems{}
}

// TestUpdateTableSweepsDeletedVectorEntriesDespiteBackfillFailure pins the
// ordering of the post-commit stages: both deletion sweeps of
// already-committed metadata run ahead of the backfills, so a request that
// adds a GSI (whose backfill panics) and deletes a vector index still
// removes the deleted index's entries — a re-sent request carries no
// deletion to re-run, and a re-created same-name index must not inherit
// stale entries.
func TestUpdateTableSweepsDeletedVectorEntriesDespiteBackfillFailure(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	real := dbstore.NewDynamoDBStore(st, st, "123456789012", "us-east-1")

	if _, err := real.Tables().Create(dbstore.CreateTableParams{
		Name:                 "SweepTbl",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	attachVec := func() {
		t.Helper()
		if _, err := real.Tables().Update("SweepTbl", func(table *dbstore.Table) error {
			table.VectorIndexes = []*dbstore.VectorIndex{{
				IndexName:           "vec",
				VectorAttributeName: "embedding",
				Dimensions:          2,
				DistanceFunction:    "COSINE",
				Projection:          &dbstore.Projection{ProjectionType: "ALL"},
				IndexStatus:         dbstore.IndexStatusActive,
			}}
			return nil
		}); err != nil {
			t.Fatalf("attach vector index: %v", err)
		}
	}
	attachVec()

	// One item with a vector, indexed through the transactional write path.
	one := "1"
	key := map[string]*dbstore.AttributeValue{"id": {S: &one}}
	attrs := map[string]*dbstore.AttributeValue{
		"embedding": {L: []*dbstore.AttributeValue{{N: &one}, {N: &one}}},
	}
	if err := real.Update(context.Background(), func(txn *dbstore.DynamoDBTxn) error {
		if err := txn.PutItem("SweepTbl", key, attrs); err != nil {
			return err
		}
		return txn.PutIndexEntries("SweepTbl", &dbstore.Item{TableName: "SweepTbl", Key: key, Attributes: attrs})
	}); err != nil {
		t.Fatalf("put indexed item: %v", err)
	}

	var svc DynamoDBService
	_, err = svc.updateTableCore(context.Background(), nil, itemsPanicStore{real}, UpdateTableInput{
		TableName: "SweepTbl",
		GSIUpdates: []interface{}{
			map[string]interface{}{
				"Create": map[string]interface{}{
					"IndexName":  "gsi-late",
					"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "id", "KeyType": "HASH"}},
					"Projection": map[string]interface{}{"ProjectionType": "ALL"},
				},
			},
		},
		VectorIndexUpdates: []interface{}{
			map[string]interface{}{"Delete": map[string]interface{}{"IndexName": "vec"}},
		},
	})
	if err == nil {
		t.Fatal("update with a panicking GSI backfill: succeeded, want the converted panic error")
	}
	if !strings.Contains(err.Error(), "panicked") {
		t.Fatalf("update error %q does not name the panic", err)
	}

	// The deleted index's entries are gone despite the failed request:
	// re-attach the same-name index and search — a stranded entry would
	// resurface as the re-created index's inherited hit.
	attachVec()
	hits := 0
	if err := real.View(context.Background(), func(txn *dbstore.DynamoDBTxn) error {
		found, err := txn.VectorTopK("SweepTbl", "vec", []float64{1, 1}, 10, nil)
		if err != nil {
			return err
		}
		hits = len(found)
		return nil
	}); err != nil {
		t.Fatalf("vector search after re-attach: %v", err)
	}
	if hits != 0 {
		t.Fatalf("re-created index inherited %d stale entries, want none", hits)
	}
}

// TestUpdateTableRejectsDuplicateGSICreate pins the added-index name
// uniqueness the GSI guide states ("The name must be unique among all the
// indexes on the table"): recreating an existing index answers the
// recreate-an-existing-resource conflict, the same refusal the vector
// update path applies.
func TestUpdateTableRejectsDuplicateGSICreate(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	real := dbstore.NewDynamoDBStore(st, st, "123456789012", "us-east-1")

	if _, err := real.Tables().Create(dbstore.CreateTableParams{
		Name:                 "DupTbl",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	gsiCreate := []interface{}{
		map[string]interface{}{
			"Create": map[string]interface{}{
				"IndexName":  "gsi-dup",
				"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "id", "KeyType": "HASH"}},
				"Projection": map[string]interface{}{"ProjectionType": "ALL"},
			},
		},
	}
	var svc DynamoDBService
	if _, err := svc.updateTableCore(context.Background(), nil, real, UpdateTableInput{TableName: "DupTbl", GSIUpdates: gsiCreate}); err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err = svc.updateTableCore(context.Background(), nil, real, UpdateTableInput{TableName: "DupTbl", GSIUpdates: gsiCreate})
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "com.amazonaws.dynamodb.v20120810#ResourceInUseException" {
		t.Fatalf("duplicate create: got %v, want ResourceInUseException", err)
	}
}
