// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"testing"

	"context"
	"errors"
	"sync"
	"time"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TestGlobalSettingsPropagationSurvivesAMemberRegionFailure pins the
// post-commit semantics of the capacity propagation: the global record has
// already committed when the member-region loop runs, so one region's
// failure must not fail the request (whose committed state stands) nor stop
// the remaining regions — every outcome is either the pre-commit
// TableNotFoundException of a race that deleted the member table before
// validation read it, or a successful response whose surviving members carry
// the propagated capacity.
func TestGlobalSettingsPropagationSurvivesAMemberRegionFailure(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	westStore, err := svc.GetCachedStoreForRegion("eu-west-1")
	if err != nil {
		t.Fatalf("west store: %v", err)
	}
	eastStore, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("east store: %v", err)
	}
	deleteWestTable := func() {
		t.Helper()
		// The cascade runs under the table's record lock: the racing
		// propagation writes the same record through the lock, and the
		// delete must serialise with it exactly as any other writer.
		_ = westStore.Tables().WithTableLock("GtPinTable", func() error {
			return westStore.Update(context.Background(), func(txn *dbstore.DynamoDBTxn) error {
				return txn.DeleteTableCascade("GtPinTable")
			})
		})
	}
	recreateWestTable := func() {
		t.Helper()
		deleteWestTable()
		if _, err := westStore.Tables().Create(dbstore.CreateTableParams{
			Name:      "GtPinTable",
			KeySchema: []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
			AttributeDefinitions: []*dbstore.AttributeDefinition{
				{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
				{AttributeName: "gid", AttributeType: dbstore.ScalarAttributeTypeS},
			},
			BillingMode: dbstore.BillingModePayPerRequest,
			GlobalSecondaryIndexes: []*dbstore.GlobalSecondaryIndex{{
				IndexName:  "gsiPin",
				KeySchema:  []*dbstore.KeySchemaElement{{AttributeName: "gid", KeyType: dbstore.KeyTypeHash}},
				Projection: &dbstore.Projection{ProjectionType: "ALL"},
			}},
			StreamSpecification: &dbstore.StreamSpecification{StreamEnabled: true, StreamViewType: dbstore.StreamViewTypeNewAndOldImages},
		}); err != nil {
			t.Fatalf("recreate west table: %v", err)
		}
	}

	// The deleter goroutines race the requests by design, but their
	// lifetime stays inside the test: the fixture's cleanup closes the
	// region storage, and a racer still sleeping past that point would
	// run its cascade against closed storage and panic. The join is a
	// cleanup of its own rather than a wait at the end of the loop: a
	// mid-test failure exits through the cleanups too, skipping any
	// in-body wait, and the cleanup ordering (last registered, first
	// run) places the join strictly before the fixture's close on every
	// exit path. The wait coordinates nothing but the join — the racers
	// still race the requests uncoordinated.
	var deleters sync.WaitGroup
	t.Cleanup(deleters.Wait)
	for round := 0; round < 12; round++ {
		// The per-round reset below switches the east table back to
		// on-demand, the one direction the platform rate limits (four
		// switches per rolling 24 hours) — scaffolding state this test
		// neither owns nor asserts, so each round starts by clearing the
		// stamps: twelve resets against a quota of four would otherwise
		// abort the run the moment four request-won rounds had stamped
		// the table, with the quota's own enforcement left to its
		// dedicated pins.
		if _, err := eastStore.Tables().Update("GtPinTable", func(table *dbstore.Table) error {
			table.BillingModeSwitches = nil
			return nil
		}); err != nil {
			t.Fatalf("round %d clear switch stamps: %v", round, err)
		}
		recreateWestTable()
		// The deleter races the request: wherever it lands, the contract
		// above must hold.
		deleters.Add(1)
		go func() {
			defer deleters.Done()
			time.Sleep(300 * time.Microsecond)
			deleteWestTable()
		}()

		_, err := svc.UpdateGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName":                          "GtPinTable",
			"GlobalTableBillingMode":                   "PROVISIONED",
			"GlobalTableProvisionedWriteCapacityUnits": 5.0,
			"GlobalTableGlobalSecondaryIndexSettingsUpdate": []interface{}{
				map[string]interface{}{"IndexName": "gsiPin", "ProvisionedWriteCapacityUnits": 4.0},
			},
			"ReplicaSettingsUpdate": []interface{}{
				map[string]interface{}{
					"RegionName":                          "us-east-1",
					"ReplicaProvisionedReadCapacityUnits": 7.0,
					"ReplicaGlobalSecondaryIndexSettingsUpdate": []interface{}{
						map[string]interface{}{"IndexName": "gsiPin", "ProvisionedReadCapacityUnits": 2.0},
					},
				},
				map[string]interface{}{
					"RegionName":                          "eu-west-1",
					"ReplicaProvisionedReadCapacityUnits": 3.0,
					"ReplicaGlobalSecondaryIndexSettingsUpdate": []interface{}{
						map[string]interface{}{"IndexName": "gsiPin", "ProvisionedReadCapacityUnits": 6.0},
					},
				},
			},
		}})
		if err != nil && !errors.Is(err, ErrTableNotFoundException) {
			t.Fatalf("round %d: a member-region propagation failure escaped the best-effort contract: %v", round, err)
		}
		if err == nil {
			// The committed round reached the surviving member: the
			// requesting region's table carries the propagated capacity.
			table := describeTableOnModePlane(t, svc, eastCtx, "GtPinTable")
			summary, ok := table["BillingModeSummary"].(map[string]interface{})
			if !ok || summary["BillingMode"] != "PROVISIONED" {
				t.Fatalf("round %d: east billing mode summary = %+v, want PROVISIONED", round, table["BillingModeSummary"])
			}
			pt := table["ProvisionedThroughput"].(map[string]interface{})
			if got := asInt64(t, pt["WriteCapacityUnits"]); got != 5 {
				t.Fatalf("round %d: east write capacity: got %d, want 5", round, got)
			}
		}

		// The next round starts from an on-demand east table again, so the
		// provisioned switch stays a switch rather than a no-op.
		if _, err := svc.UpdateGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName":        "GtPinTable",
			"GlobalTableBillingMode": "PAY_PER_REQUEST",
		}}); err != nil && !errors.Is(err, ErrTableNotFoundException) {
			t.Fatalf("round %d on-demand reset: %v", round, err)
		}
	}
}

// TestReplicaDeletionSurvivesAPostCommitCascadeFailure pins the
// post-commit contract of the Delete action: the membership record has
// already committed when the destination region's table cascade runs, so
// a cascade failure must not error the response — the committed
// membership is the truth, and erroring would hand the caller a retry
// that answers ReplicaNotFound. The failure is injected by corrupting the
// destination table's stored record so the cascade cannot read it.
func TestReplicaDeletionSurvivesAPostCommitCascadeFailure(t *testing.T) {
	svc, eastCtx, _, sm := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	westStorage, err := sm.GetStorage("eu-west-1")
	if err != nil {
		t.Fatalf("west region storage: %v", err)
	}
	if err := westStorage.Bucket("dynamodb_tables-eu-west-1").Put([]byte("GtPinTable"), []byte("not-a-protobuf-record")); err != nil {
		t.Fatalf("corrupt west table record: %v", err)
	}

	resp, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "eu-west-1"}},
		},
	}})
	if err != nil {
		t.Fatalf("delete replica over a failing cascade: %v", err)
	}
	group := resp.(map[string]interface{})["GlobalTableDescription"].(map[string]interface{})["ReplicationGroup"].([]map[string]interface{})
	if len(group) != 1 || group[0]["RegionName"] != "us-east-1" {
		t.Fatalf("committed membership after the failed cascade: %#v", group)
	}

	// The retry the contract protects: the same Delete action answers the
	// documented replica-not-found verdict, not the cascade's storage
	// error — the member is committed out of the record.
	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "eu-west-1"}},
		},
	}}); !errors.Is(err, ErrReplicaNotFound) {
		t.Fatalf("retry after the failed cascade: expected ErrReplicaNotFound, got %v", err)
	}
}

func TestReplicaOpsUseCanonicalWriteComposition(t *testing.T) {
	svc, _, store := serviceRegionFixture(t)
	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "RepTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	key := map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("rep1")}
	metrics := func() (int64, int64) {
		t.Helper()
		current, err := store.Tables().Get("RepTable")
		if err != nil {
			t.Fatalf("table get: %v", err)
		}
		return current.ItemCount, current.TableSizeBytes
	}

	attrs := map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("rep1"), "v": dbstore.StringValue("one")}
	if err := svc.replicaPutOp(table, key, attrs)(context.Background(), store); err != nil {
		t.Fatalf("replica put: %v", err)
	}
	count, size := metrics()
	wantSize := dbstore.CalculateItemSize(attrs)
	if count != 1 || size != wantSize {
		t.Fatalf("after first put: expected count 1 size %d, got count %d size %d", wantSize, count, size)
	}

	bigger := map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("rep1"), "v": dbstore.StringValue("much longer replacement value")}
	if err := svc.replicaPutOp(table, key, bigger)(context.Background(), store); err != nil {
		t.Fatalf("replica re-put: %v", err)
	}
	count, size = metrics()
	wantBigger := dbstore.CalculateItemSize(bigger)
	if count != 1 || size != wantBigger {
		t.Fatalf("after re-put: expected count 1 size %d, got count %d size %d", wantBigger, count, size)
	}

	if err := svc.replicaDeleteOp(table, key)(context.Background(), store); err != nil {
		t.Fatalf("replica delete: %v", err)
	}
	if _, err := store.Items().Get("RepTable", key); !dbstore.IsItemNotFound(err) {
		t.Fatalf("after delete: expected item not found, got %v", err)
	}
	if count, size = metrics(); count != 0 || size != 0 {
		t.Fatalf("after delete: expected count 0 size 0, got count %d size %d", count, size)
	}
}

// TestReplicaStreamCaptureFaultRollsBackTheItemWrite pins the replica
// capture's fault rule: the item write and the replica's stream record are
// one update, so an in-transaction GetTable read fault (not the table's
// absence) aborts the whole transaction — the destination keeps neither
// half — and re-running the operation once the fault heals performs both.
func TestReplicaStreamCaptureFaultRollsBackTheItemWrite(t *testing.T) {
	svc, _, store := serviceRegionFixture(t)
	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "RepFaultTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
		StreamSpecification: &dbstore.StreamSpecification{
			StreamEnabled:  true,
			StreamViewType: dbstore.StreamViewTypeNewImage,
		},
	})
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	// Corrupt the destination's table record: the in-transaction GetTable
	// of the capture half fails on the undecodable payload — the injected
	// read fault.
	tableBucket := store.Storage().Bucket("dynamodb_tables-us-east-1")
	tableKey := []byte("RepFaultTable")
	original, err := tableBucket.Get(tableKey)
	if err != nil || original == nil {
		t.Fatalf("read table record: %v", err)
	}
	if err := tableBucket.Put(tableKey, []byte("undecodable-table-record")); err != nil {
		t.Fatalf("corrupt table record: %v", err)
	}

	key := map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("fault1")}
	attrs := map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("fault1"), "v": dbstore.StringValue("one")}
	if err := svc.replicaPutOp(table, key, attrs)(context.Background(), store); err == nil {
		t.Fatal("replica put under a corrupted table record succeeded")
	}
	// The corrupt table record fails every schema-aware read, so the
	// absence proof walks the item bucket's prefix directly: no item key
	// may exist under the table while the fault stood.
	itemBucket := store.Storage().Bucket("dynamodb_items-us-east-1")
	if n := bucketPrefixCount(itemBucket, "RepFaultTable"+dbstore.KeySep); n != 0 {
		t.Fatalf("item committed despite the capture read fault: %d item key(s) under the table", n)
	}

	// Heal the record: the retried operation performs both halves — the
	// item lands and the replica's own stream carries the record.
	if err := tableBucket.Put(tableKey, original); err != nil {
		t.Fatalf("heal table record: %v", err)
	}
	if err := svc.replicaPutOp(table, key, attrs)(context.Background(), store); err != nil {
		t.Fatalf("replica put after heal: %v", err)
	}
	if _, err := store.Items().Get("RepFaultTable", key); err != nil {
		t.Fatalf("item missing after the healed retry: %v", err)
	}
	healed, err := store.Tables().Get("RepFaultTable")
	if err != nil {
		t.Fatalf("table get after heal: %v", err)
	}
	records, _, err := store.Streams().GetRecords("RepFaultTable", healed.StreamArn, 0, 10)
	if err != nil {
		t.Fatalf("stream records after the healed retry: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("stream records after the healed retry = %d, want the one replication record", len(records))
	}
}

// bucketPrefixCount counts the keys a bucket holds under one prefix.
func bucketPrefixCount(bucket storage.Bucket, prefix string) int {
	n := 0
	iter := bucket.ScanPrefix([]byte(prefix))
	for iter.Next() {
		n++
	}
	iter.Close()
	return n
}

// TestExhaustedReplicationEscalatesAndConverges pins the convergence
// worker: a delivery whose fault outlives the in-flight bounded retries is
// escalated rather than logged away, and once the fault heals the escalated
// delivery lands — the replicas converge. The fault is the destination's
// corrupted table record; healing it while the escalation is queued is the
// heal the worker must pick up.
func TestExhaustedReplicationEscalatesAndConverges(t *testing.T) {
	svc, eastCtx, westCtx, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)

	westStore, err := svc.GetCachedStoreForRegion("eu-west-1")
	if err != nil {
		t.Fatalf("west store: %v", err)
	}
	tableBucket := westStore.Storage().Bucket("dynamodb_tables-eu-west-1")
	tableKey := []byte("GtPinTable")
	original, err := tableBucket.Get(tableKey)
	if err != nil || original == nil {
		t.Fatalf("read west table record: %v", err)
	}
	if err := tableBucket.Put(tableKey, []byte("undecodable-table-record")); err != nil {
		t.Fatalf("corrupt west table record: %v", err)
	}

	// The source-side put succeeds; the delivery to the corrupted
	// destination exhausts its in-flight retries and escalates.
	if _, err := svc.PutItem(context.Background(), eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
		"Item":      map[string]interface{}{"id": map[string]interface{}{"S": "converge-me"}},
	}}); err != nil {
		t.Fatalf("source put under fault: %v", err)
	}

	// Wait until the escalation is observable (the first enqueue), then
	// heal the destination's table record: the escalated delivery's
	// retries keep failing until they see the healed record.
	deadline := time.Now().Add(5 * time.Second)
	for svc.pendingEscalationCount() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("exhausted delivery never escalated")
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := tableBucket.Put(tableKey, original); err != nil {
		t.Fatalf("heal west table record: %v", err)
	}

	// Convergence: the replicated item becomes readable from the healed
	// destination region.
	waitForReplicatedItem(t, svc, westCtx, "converge-me")
}
