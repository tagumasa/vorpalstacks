package dynamodb

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// This file pins the global-table record's membership lifecycle: the
// account-global record's visibility, the joining conditions, group
// management through both replica surfaces (including member departure and
// its destination cascade), the record's observable settings, and the
// conversion seed race's error mapping. The replication-delivery pins
// (capacity propagation, capture faults, escalation convergence, write
// composition) live in global_table_replication_test.go.

// TestGlobalTableRecordIsVisibleAndReplicatesFromEveryMemberRegion pins the
// account-global treatment of the global-table record: created through one
// member region, the record answers Describe and List from every member
// region, and an item written in the non-creating member region replicates
// to the creating region (and the reverse direction too) — the multi-active
// contract, which a per-region record silently broke for every region but
// the one that ran CreateGlobalTable.
func TestGlobalTableRecordIsVisibleAndReplicatesFromEveryMemberRegion(t *testing.T) {
	svc, eastCtx, westCtx, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)

	for _, reqCtx := range []*request.RequestContext{eastCtx, westCtx} {
		resp, err := svc.DescribeGlobalTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName": "GtPinTable",
		}})
		if err != nil {
			t.Fatalf("describe from %s: %v", reqCtx.GetRegion(), err)
		}
		desc := resp.(map[string]interface{})["GlobalTableDescription"].(map[string]interface{})
		if got := len(desc["ReplicationGroup"].([]map[string]interface{})); got != 2 {
			t.Fatalf("replication group from %s: got %d replicas, want 2", reqCtx.GetRegion(), got)
		}
	}

	listResp, err := svc.ListGlobalTables(context.Background(), westCtx, &request.ParsedRequest{Parameters: map[string]interface{}{}})
	if err != nil {
		t.Fatalf("list from non-creating region: %v", err)
	}
	found := false
	for _, entry := range listResp.(map[string]interface{})["GlobalTables"].([]map[string]interface{}) {
		if entry["GlobalTableName"] == "GtPinTable" {
			found = true
		}
	}
	if !found {
		t.Fatal("ListGlobalTables from the non-creating region does not see the global table")
	}

	put := func(reqCtx *request.RequestContext, id string) {
		t.Helper()
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "GtPinTable",
			"Item":      map[string]interface{}{"id": map[string]interface{}{"S": id}},
		}}); err != nil {
			t.Fatalf("put %s from %s: %v", id, reqCtx.GetRegion(), err)
		}
	}
	put(westCtx, "from-west")
	put(eastCtx, "from-east")
	waitForReplicatedItem(t, svc, eastCtx, "from-west")
	waitForReplicatedItem(t, svc, westCtx, "from-east")
}

// TestListGlobalTablesLimitPresence pins the Limit member's presence
// semantics: an explicit zero is the PositiveIntegerObject range
// violation, while an omitted member keeps the default page.
func TestListGlobalTablesLimitPresence(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)

	if _, err := svc.ListGlobalTables(context.Background(), eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Limit": float64(0),
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("explicit Limit=0: expected ErrInvalidParameter, got %v", err)
	}

	resp, err := svc.ListGlobalTables(context.Background(), eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{}})
	if err != nil {
		t.Fatalf("omitted limit: %v", err)
	}
	tables := resp.(map[string]interface{})["GlobalTables"].([]map[string]interface{})
	if len(tables) == 0 {
		t.Fatal("omitted limit: expected the default page to list the record")
	}
}

// TestDeleteTableRemovesOnlyTheDeletedMemberFromTheGlobalTableRecord pins
// the table cascade's global-table cleanup on the account-global record:
// deleting one member's table removes that member alone — the group and
// its surviving replica answer DescribeGlobalTable from every region's
// view — and only the final member's departure takes the record with it.
func TestDeleteTableRemovesOnlyTheDeletedMemberFromTheGlobalTableRecord(t *testing.T) {
	svc, eastCtx, westCtx, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)

	if _, err := svc.DeleteTable(context.Background(), eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
	}}); err != nil {
		t.Fatalf("delete member table: %v", err)
	}

	for _, reqCtx := range []*request.RequestContext{eastCtx, westCtx} {
		resp, err := svc.DescribeGlobalTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName": "GtPinTable",
		}})
		if err != nil {
			t.Fatalf("describe after one member's delete from %s: %v", reqCtx.GetRegion(), err)
		}
		group := resp.(map[string]interface{})["GlobalTableDescription"].(map[string]interface{})["ReplicationGroup"].([]map[string]interface{})
		if len(group) != 1 || group[0]["RegionName"] != "eu-west-1" {
			t.Fatalf("replication group after one member's delete from %s: %+v", reqCtx.GetRegion(), group)
		}
	}

	// The final member's departure takes the record with it: the emptied
	// group answers GlobalTableNotFoundException from every region's view,
	// the same end state the Delete replica action reaches.
	if _, err := svc.DeleteTable(context.Background(), westCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
	}}); err != nil {
		t.Fatalf("delete final member table: %v", err)
	}
	for _, reqCtx := range []*request.RequestContext{eastCtx, westCtx} {
		_, err := svc.DescribeGlobalTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName": "GtPinTable",
		}})
		if !errors.Is(err, ErrGlobalTableNotFound) {
			t.Fatalf("describe after final member delete from %s: expected GlobalTableNotFound, got %v", reqCtx.GetRegion(), err)
		}
	}
}

// TestReplicaDeleteActionCascadesToDestinationRegion pins the Delete
// action's destination half on the legacy surface: removing a replica takes
// its region's table — with its items — through the deletion cascade while
// the surviving member keeps serving, a deletion-protected replica table
// refuses the action with its membership unchanged, a re-created qualifying
// table re-joins (no orphan blocks the join), and the final member's
// departure takes the global-table record with it.
func TestReplicaDeleteActionCascadesToDestinationRegion(t *testing.T) {
	svc, eastCtx, westCtx, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	// Both members hold data before the removal.
	if _, err := svc.PutItem(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
		"Item":      map[string]interface{}{"id": map[string]interface{}{"S": "before-removal"}},
	}}); err != nil {
		t.Fatalf("put in the surviving member: %v", err)
	}
	waitForReplicatedItem(t, svc, westCtx, "before-removal")

	// The Delete action removes eu-west-1: the record keeps the surviving
	// member, the destination region's table is gone, and the surviving
	// member still serves the replicated item.
	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "eu-west-1"}},
		},
	}}); err != nil {
		t.Fatalf("delete replica: %v", err)
	}
	describeGroup := func() []map[string]interface{} {
		t.Helper()
		resp, err := svc.DescribeGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName": "GtPinTable",
		}})
		if err != nil {
			t.Fatalf("describe global table: %v", err)
		}
		return resp.(map[string]interface{})["GlobalTableDescription"].(map[string]interface{})["ReplicationGroup"].([]map[string]interface{})
	}
	if group := describeGroup(); len(group) != 1 || group[0]["RegionName"] != "us-east-1" {
		t.Fatalf("replication group after the delete action: %+v", group)
	}
	if _, err := svc.DescribeTable(ctx, westCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
	}}); !errors.Is(err, ErrTableNotFound) {
		t.Fatalf("removed replica's table in eu-west-1: expected ErrTableNotFound, got %v", err)
	}
	if _, err := svc.GetItem(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
		"Key":       map[string]interface{}{"id": map[string]interface{}{"S": "before-removal"}},
	}}); err != nil {
		t.Fatalf("surviving member still serves the item: %v", err)
	}

	// A re-created qualifying table re-joins: the cascade left no orphan
	// behind to block the join.
	westStore, err := svc.GetCachedStoreForRegion("eu-west-1")
	if err != nil {
		t.Fatalf("eu-west-1 store: %v", err)
	}
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
		t.Fatalf("re-create eu-west-1 table: %v", err)
	}
	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}},
		},
	}}); err != nil {
		t.Fatalf("re-join after the cascade: %v", err)
	}

	// A deletion-protected replica table refuses the Delete action, with the
	// membership unchanged.
	if _, err := svc.UpdateTable(ctx, westCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "GtPinTable",
		"DeletionProtectionEnabled": true,
	}}); err != nil {
		t.Fatalf("enable deletion protection on the replica table: %v", err)
	}
	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "eu-west-1"}},
		},
	}}); !errors.Is(err, ErrTableDeletionProtected) {
		t.Fatalf("delete of a deletion-protected replica: expected ErrTableDeletionProtected, got %v", err)
	}
	if group := describeGroup(); len(group) != 2 {
		t.Fatalf("membership after the refused delete: got %d replicas, want 2", len(group))
	}

	// Protection lifted, the Delete goes through.
	if _, err := svc.UpdateTable(ctx, westCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "GtPinTable",
		"DeletionProtectionEnabled": false,
	}}); err != nil {
		t.Fatalf("disable deletion protection on the replica table: %v", err)
	}
	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "eu-west-1"}},
		},
	}}); err != nil {
		t.Fatalf("delete replica after the protection lift: %v", err)
	}

	// The final member's departure takes the record with it and deletes the
	// last table.
	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "us-east-1"}},
		},
	}}); err != nil {
		t.Fatalf("delete the final replica: %v", err)
	}
	if _, err := svc.DescribeGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
	}}); !errors.Is(err, ErrGlobalTableNotFound) {
		t.Fatalf("record after the final member's departure: expected ErrGlobalTableNotFound, got %v", err)
	}
	if _, err := svc.DescribeTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
	}}); !errors.Is(err, ErrTableNotFound) {
		t.Fatalf("final member's table: expected ErrTableNotFound, got %v", err)
	}
}

// TestReplicaGSIOverrideWithUnresolvableReferenceRejected pins the
// reference guard of the GSI-override path: a global-table record whose
// replica regions no longer resolve a member table (out-of-band
// deletions) makes the reference lookup answer nil, and the replica
// Update action's GlobalSecondaryIndexes member must then answer the
// validation error — never dereference the nil reference.
func TestReplicaGSIOverrideWithUnresolvableReferenceRejected(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)

	// Out-of-band: the record's replication group now names a region
	// whose store holds no member table, so no replica resolves a
	// reference.
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	if _, err := store.GlobalTables().Update("GtPinTable", func(gt *dbstore.GlobalTable) error {
		gt.ReplicationGroup = []*dbstore.Replica{{RegionName: "ap-southeast-2", ReplicaStatus: "ACTIVE"}}
		return nil
	}); err != nil {
		t.Fatalf("seed unresolvable membership: %v", err)
	}

	_, err = svc.UpdateGlobalTable(context.Background(), eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Update": map[string]interface{}{
				"RegionName": "ap-southeast-2",
				"GlobalSecondaryIndexes": []interface{}{
					map[string]interface{}{"IndexName": "gid-index"},
				},
			}},
		},
	}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("GSI override with unresolvable reference: expected ErrInvalidParameter, got %v", err)
	}
}

// TestDescribeGlobalTableEchoesStoredReplicaSettings pins the
// ReplicaDescription shape: the per-replica settings an Update action
// persists — the KMS key identifier, both read-throughput overrides, the
// per-index capacity overrides, and the table class override — are
// observable on the DescribeGlobalTable face, with the untouched replica
// carrying region and status alone.
func TestDescribeGlobalTableEchoesStoredReplicaSettings(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{
			map[string]interface{}{"Update": map[string]interface{}{
				"RegionName":         "us-east-1",
				"KMSMasterKeyId":     "alias/geopolitical",
				"TableClassOverride": "STANDARD_INFREQUENT_ACCESS",
				"ProvisionedThroughputOverride": map[string]interface{}{
					"ReadCapacityUnits": float64(25),
				},
				"GlobalSecondaryIndexes": []interface{}{
					map[string]interface{}{
						"IndexName": "gsiPin",
						"ProvisionedThroughputOverride": map[string]interface{}{
							"ReadCapacityUnits": float64(15),
						},
					},
				},
			}},
		},
	}}); err != nil {
		t.Fatalf("update replica settings: %v", err)
	}

	resp, err := svc.DescribeGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
	}})
	if err != nil {
		t.Fatalf("describe global table: %v", err)
	}
	group := resp.(map[string]interface{})["GlobalTableDescription"].(map[string]interface{})["ReplicationGroup"].([]map[string]interface{})
	if len(group) != 2 {
		t.Fatalf("replication group: expected 2 replicas, got %#v", group)
	}

	var updated, untouched map[string]interface{}
	for _, replica := range group {
		if replica["RegionName"] == "us-east-1" {
			updated = replica
		} else {
			untouched = replica
		}
	}
	if updated == nil || untouched == nil {
		t.Fatalf("replication group regions: %#v", group)
	}

	if updated["KMSMasterKeyId"] != "alias/geopolitical" {
		t.Fatalf("KMSMasterKeyId: %v", updated["KMSMasterKeyId"])
	}
	ptOverride := updated["ProvisionedThroughputOverride"].(map[string]interface{})
	if ptOverride["ReadCapacityUnits"] != int64(25) {
		t.Fatalf("ProvisionedThroughputOverride: %#v", ptOverride)
	}
	classSummary := updated["ReplicaTableClassSummary"].(map[string]interface{})
	if classSummary["TableClass"] != "STANDARD_INFREQUENT_ACCESS" {
		t.Fatalf("ReplicaTableClassSummary: %#v", classSummary)
	}
	if _, has := classSummary["LastUpdateDateTime"]; !has {
		t.Fatalf("ReplicaTableClassSummary: expected LastUpdateDateTime, got %#v", classSummary)
	}
	gsiList := updated["GlobalSecondaryIndexes"].([]map[string]interface{})
	if len(gsiList) != 1 || gsiList[0]["IndexName"] != "gsiPin" {
		t.Fatalf("GlobalSecondaryIndexes: %#v", gsiList)
	}
	gsiOverride := gsiList[0]["ProvisionedThroughputOverride"].(map[string]interface{})
	if gsiOverride["ReadCapacityUnits"] != int64(15) {
		t.Fatalf("GSI ProvisionedThroughputOverride: %#v", gsiOverride)
	}

	// The untouched replica carries the members every replica has —
	// region and status — and none of the overrides it never stored.
	if len(untouched) != 2 {
		t.Fatalf("untouched replica: expected region and status alone, got %#v", untouched)
	}
}

// TestJoiningReplicaMustMatchWriteCapacity pins the capacity half of the
// documented joining conditions: a replica table whose write capacity
// settings differ from the reference replica's — different provisioned
// write units, or a different billing mode — is rejected by both
// CreateGlobalTable and UpdateGlobalTable's replica-create, while a
// matching table joins.
func TestJoiningReplicaMustMatchWriteCapacity(t *testing.T) {
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	ctx := context.Background()
	eastCtx := request.NewRequestContext(ctx, sm, "123456789012", "us-east-1")

	// Regions: us-east-1 and eu-west-1 form the initial pair (WCU 5);
	// us-west-2 carries WCU 10 (unit mismatch), af-south-1 WCU 5 (match),
	// ap-northeast-1 is on-demand (billing-mode mismatch).
	create := func(region string, units int64, onDemand bool) {
		t.Helper()
		store, err := svc.GetCachedStoreForRegion(region)
		if err != nil {
			t.Fatalf("region store %s: %v", region, err)
		}
		params := dbstore.CreateTableParams{
			Name:                 "GtCapTable",
			KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
			AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
			StreamSpecification:  &dbstore.StreamSpecification{StreamEnabled: true, StreamViewType: dbstore.StreamViewTypeNewAndOldImages},
		}
		if onDemand {
			params.BillingMode = dbstore.BillingModePayPerRequest
		} else {
			params.BillingMode = dbstore.BillingModeProvisioned
			params.ProvisionedThroughput = &dbstore.ProvisionedThroughput{ReadCapacityUnits: units, WriteCapacityUnits: units}
		}
		if _, err := store.Tables().Create(params); err != nil {
			t.Fatalf("create table in %s: %v", region, err)
		}
	}
	create("us-east-1", 5, false)
	create("eu-west-1", 5, false)
	create("us-west-2", 10, false)
	create("af-south-1", 5, false)
	create("ap-northeast-1", 0, true)

	replicationGroup := func(regions ...string) []interface{} {
		var group []interface{}
		for _, region := range regions {
			group = append(group, map[string]interface{}{"RegionName": region})
		}
		return group
	}

	// A single-region group is not a replication relationship: the
	// documented two-or-more rule rejects it outright, before any
	// per-table validation runs.
	if _, err := svc.CreateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName":  "GtCapTable",
		"ReplicationGroup": replicationGroup("us-east-1"),
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("single-region create: expected ErrInvalidParameter, got %v", err)
	}

	// CreateGlobalTable rejects the pair whose write units differ.
	if _, err := svc.CreateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName":  "GtCapTable",
		"ReplicationGroup": replicationGroup("us-east-1", "us-west-2"),
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("create with mismatched write units: expected ErrInvalidParameter, got %v", err)
	}

	// The matching pair forms the global table.
	if _, err := svc.CreateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName":  "GtCapTable",
		"ReplicationGroup": replicationGroup("us-east-1", "eu-west-1"),
	}}); err != nil {
		t.Fatalf("create with matching write units: %v", err)
	}

	addReplica := func(region string) error {
		_, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName": "GtCapTable",
			"ReplicaUpdates": []interface{}{
				map[string]interface{}{"Create": map[string]interface{}{"RegionName": region}},
			},
		}})
		return err
	}

	if err := addReplica("us-west-2"); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("join with different write units: expected ErrInvalidParameter, got %v", err)
	}
	if err := addReplica("ap-northeast-1"); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("join with different billing mode: expected ErrInvalidParameter, got %v", err)
	}
	if err := addReplica("af-south-1"); err != nil {
		t.Fatalf("join with matching write units: %v", err)
	}

	resp, err := svc.DescribeGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtCapTable",
	}})
	if err != nil {
		t.Fatalf("describe after join: %v", err)
	}
	desc := resp.(map[string]interface{})["GlobalTableDescription"].(map[string]interface{})
	if got := len(desc["ReplicationGroup"].([]map[string]interface{})); got != 3 {
		t.Fatalf("replication group after the accepted join: got %d replicas, want 3", got)
	}
}

// TestUpdateTableReplicaUpdatesManageTheReplicationGroup pins the modern
// replica-management surface the model defines on UpdateTable: a Create
// action on a regional table converts it into a global table with the
// requesting region as the implicit first member, further Create and Delete
// actions manage the group, the Update action applies the per-replica
// overrides the record carries and refuses the ones it cannot, and the
// legacy UpdateGlobalTable operation keeps requiring an existing record —
// the seeding belongs to the modern surface alone.
func TestUpdateTableReplicaUpdatesManageTheReplicationGroup(t *testing.T) {
	svc, eastCtx, westCtx, sm := globalTablePlaneFixture(t)
	ctx := context.Background()

	replicaRegions := func() map[string]bool {
		t.Helper()
		resp, err := svc.DescribeGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"GlobalTableName": "GtPinTable"}})
		if err != nil {
			t.Fatalf("describe global table: %v", err)
		}
		group := resp.(map[string]interface{})["GlobalTableDescription"].(map[string]interface{})["ReplicationGroup"].([]map[string]interface{})
		regions := make(map[string]bool, len(group))
		for _, replica := range group {
			regions[replica["RegionName"].(string)] = true
		}
		return regions
	}
	update := func(params map[string]interface{}) error {
		_, err := svc.UpdateTable(ctx, eastCtx, &request.ParsedRequest{Parameters: params})
		return err
	}

	// The replica update list carries a documented minimum of one entry:
	// a present-but-empty list is rejected rather than read as an omission.
	if err := update(map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty replica updates: expected ErrInvalidParameter, got %v", err)
	}

	// The ACTIVE gate covers the membership pass: a table outside ACTIVE
	// refuses replica updates with the same ResourceInUseException every
	// other UpdateTable member meets, and the refused action leaves the
	// table regional (a later Create action on the restored-ACTIVE table
	// still performs the conversion).
	eastStore, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("east store: %v", err)
	}
	if _, err := eastStore.Tables().Update("GtPinTable", func(table *dbstore.Table) error {
		table.Status = dbstore.TableStatusCreating
		return nil
	}); err != nil {
		t.Fatalf("flip table to CREATING: %v", err)
	}
	if err := update(map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}}},
	}); !errors.Is(err, ErrTableNotActive) {
		t.Fatalf("replica update on a CREATING table: expected ErrTableNotActive, got %v", err)
	}
	if _, err := eastStore.Tables().Update("GtPinTable", func(table *dbstore.Table) error {
		table.Status = dbstore.TableStatusActive
		return nil
	}); err != nil {
		t.Fatalf("restore table to ACTIVE: %v", err)
	}

	// The first Create action converts the regional table: the requesting
	// region is the implicit first member.
	resp, err := svc.UpdateTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}}},
	}})
	if err != nil {
		t.Fatalf("modern create replica: %v", err)
	}
	replicas := resp.(map[string]interface{})["TableDescription"].(map[string]interface{})["Replicas"].([]interface{})
	if len(replicas) != 2 {
		t.Fatalf("UpdateTable response replicas: %+v", replicas)
	}
	if regions := replicaRegions(); !regions["us-east-1"] || !regions["eu-west-1"] {
		t.Fatalf("replication group after the create action: %+v", regions)
	}

	// The seeded record drives real replication.
	if _, err := svc.PutItem(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
		"Item":      map[string]interface{}{"id": map[string]interface{}{"S": "modern-flow"}},
	}}); err != nil {
		t.Fatalf("put in the seed region: %v", err)
	}
	waitForReplicatedItem(t, svc, westCtx, "modern-flow")

	// A third region joins through the same surface.
	thirdStore, err := svc.GetCachedStoreForRegion("us-west-2")
	if err != nil {
		t.Fatalf("third region store: %v", err)
	}
	if _, err := thirdStore.Tables().Create(dbstore.CreateTableParams{
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
		t.Fatalf("create third-region table: %v", err)
	}
	if err := update(map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "us-west-2"}}},
	}); err != nil {
		t.Fatalf("create third replica: %v", err)
	}
	if regions := replicaRegions(); len(regions) != 3 || !regions["us-west-2"] {
		t.Fatalf("replication group after the third member: %+v", regions)
	}

	// The Update action applies the overrides the record carries. The
	// throughput override is the read side alone in the model's shape —
	// write capacity is global across a replication group — so the
	// documented single-member request applies and a write member the
	// shape does not define is never recorded.
	if err := update(map[string]interface{}{
		"TableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"RegionName":                    "us-west-2",
			"ProvisionedThroughputOverride": map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 7.0},
			"TableClassOverride":            "STANDARD_INFREQUENT_ACCESS",
		}}},
	}); err != nil {
		t.Fatalf("update replica overrides: %v", err)
	}
	if err := update(map[string]interface{}{
		"TableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"RegionName":                    "us-west-2",
			"ProvisionedThroughputOverride": map[string]interface{}{"ReadCapacityUnits": 6.0},
		}}},
	}); err != nil {
		t.Fatalf("documented single-member override: %v", err)
	}
	if err := update(map[string]interface{}{
		"TableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"RegionName":                    "us-west-2",
			"ProvisionedThroughputOverride": map[string]interface{}{"ReadCapacityUnits": 0.0},
		}}},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("zero read override: expected ErrInvalidParameter, got %v", err)
	}
	settingsResp, err := svc.DescribeGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"GlobalTableName": "GtPinTable"}})
	if err != nil {
		t.Fatalf("describe global table settings: %v", err)
	}
	settings := settingsResp.(map[string]interface{})["ReplicaSettings"].([]map[string]interface{})
	var west *map[string]interface{}
	for i := range settings {
		if settings[i]["RegionName"] == "us-west-2" {
			west = &settings[i]
		}
	}
	if west == nil || (*west)["ReplicaProvisionedReadCapacityUnits"].(int64) != 6 {
		t.Fatalf("replica capacity override: %+v", west)
	}
	class, ok := (*west)["ReplicaTableClassSummary"].(map[string]interface{})
	if !ok || class["TableClass"] != "STANDARD_INFREQUENT_ACCESS" {
		t.Fatalf("replica table class override: %+v", (*west)["ReplicaTableClassSummary"])
	}

	// The Update action applies the documented members the model defines:
	// the KMS key identifier and the on-demand read maximum record on the
	// replica, and a per-index override names an index the reference table
	// actually carries — an unknown index name is the index not-found
	// error, not a recorded entry.
	kmsUpdate := map[string]interface{}{
		"TableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"RegionName":                 "us-west-2",
			"KMSMasterKeyId":             "alias/replica-key",
			"OnDemandThroughputOverride": map[string]interface{}{"MaxReadRequestUnits": 9.0},
			"GlobalSecondaryIndexes": []interface{}{map[string]interface{}{
				"IndexName":                     "gsiPin",
				"ProvisionedThroughputOverride": map[string]interface{}{"ReadCapacityUnits": 3.0},
			}},
		}}},
	}
	kmsResp, err := svc.UpdateTable(ctx, eastCtx, &request.ParsedRequest{Parameters: kmsUpdate})
	if err != nil {
		t.Fatalf("update replica member overrides: %v", err)
	}
	replicasEcho := kmsResp.(map[string]interface{})["TableDescription"].(map[string]interface{})["Replicas"].([]interface{})
	var westEcho map[string]interface{}
	for _, raw := range replicasEcho {
		entry := raw.(map[string]interface{})
		if entry["RegionName"] == "us-west-2" {
			westEcho = entry
		}
	}
	if westEcho == nil || westEcho["KMSMasterKeyId"] != "alias/replica-key" {
		t.Fatalf("replica KMS override echo: %+v", westEcho)
	}
	if odt, ok := westEcho["OnDemandThroughputOverride"].(map[string]interface{}); !ok || odt["MaxReadRequestUnits"].(int64) != 9 {
		t.Fatalf("replica on-demand override echo: %+v", westEcho["OnDemandThroughputOverride"])
	}
	gsiEcho, ok := westEcho["GlobalSecondaryIndexes"].([]map[string]interface{})
	if !ok || len(gsiEcho) != 1 || gsiEcho[0]["IndexName"] != "gsiPin" {
		t.Fatalf("replica GSI override echo: %+v", westEcho["GlobalSecondaryIndexes"])
	}
	if pt, ok := gsiEcho[0]["ProvisionedThroughputOverride"].(map[string]interface{}); !ok || pt["ReadCapacityUnits"].(int64) != 3 {
		t.Fatalf("replica GSI throughput echo: %+v", gsiEcho[0]["ProvisionedThroughputOverride"])
	}
	if err := update(map[string]interface{}{
		"TableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"RegionName": "us-west-2",
			"GlobalSecondaryIndexes": []interface{}{map[string]interface{}{
				"IndexName": "no-such-index",
			}},
		}}},
	}); !errors.Is(err, ErrIndexNotFound) {
		t.Fatalf("unknown GSI override index: expected ErrIndexNotFound, got %v", err)
	}

	// The KMS key identifier decodes through the typed-member path and
	// validates to the SSE plane's identifier standard: a non-string
	// value, a malformed identifier, and an ARN that is no key
	// identifier (another service's ARN, or a KMS ARN whose resource
	// names no key or alias) are all validation errors, never silently
	// skipped updates.
	for _, bad := range []interface{}{float64(42), "alias/", "not/alias/or/arn",
		"arn:aws:sqs:us-east-1:123456789012:queue/leak",
		"arn:aws:kms:us-east-1:123456789012:policy/p1"} {
		if err := update(map[string]interface{}{
			"TableName": "GtPinTable",
			"ReplicaUpdates": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
				"RegionName":     "us-west-2",
				"KMSMasterKeyId": bad,
			}}},
		}); !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("KMSMasterKeyId %v: expected ErrInvalidParameter, got %v", bad, err)
		}
	}

	// A full KMS key ARN is the identifier standard's ARN form — accepted
	// and stored as given.
	if err := update(map[string]interface{}{
		"TableName": "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"RegionName":     "us-west-2",
			"KMSMasterKeyId": "arn:aws:kms:us-east-1:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
		}}},
	}); err != nil {
		t.Fatalf("KMS key ARN update: %v", err)
	}

	// A Delete action removes the third member.
	if err := update(map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "us-west-2"}}},
	}); err != nil {
		t.Fatalf("delete replica: %v", err)
	}
	if regions := replicaRegions(); len(regions) != 2 || regions["us-west-2"] {
		t.Fatalf("replication group after the delete action: %+v", regions)
	}

	// The removed member's table is deleted in its own region — the Delete
	// action's destination half — while the surviving members keep theirs.
	west2Ctx := request.NewRequestContext(ctx, sm, "123456789012", "us-west-2")
	if _, err := svc.DescribeTable(ctx, west2Ctx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
	}}); !errors.Is(err, ErrTableNotFound) {
		t.Fatalf("removed replica's table in us-west-2: expected ErrTableNotFound, got %v", err)
	}
	if _, err := svc.DescribeTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "GtPinTable",
	}}); err != nil {
		t.Fatalf("surviving member's table in us-east-1: %v", err)
	}

	// A Delete action naming the REQUESTING region completes the request:
	// the membership change commits and this region's own table cascades
	// away behind it, so the response renders the departed table's
	// description in the DELETING state rather than answering the
	// completed deletion with a not-found error.
	ownResp, err := svc.UpdateTable(ctx, westCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "eu-west-1"}}},
	}})
	if err != nil {
		t.Fatalf("own-region delete action: %v", err)
	}
	ownDesc := ownResp.(map[string]interface{})["TableDescription"].(map[string]interface{})
	if ownDesc["TableStatus"] != "DELETING" {
		t.Fatalf("own-region delete action status: %v", ownDesc["TableStatus"])
	}
	if regions := replicaRegions(); len(regions) != 1 || !regions["us-east-1"] {
		t.Fatalf("replication group after the own-region delete: %+v", regions)
	}

	// Duplicate create and unknown-region delete both refuse.
	if err := update(map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}}},
	}); err == nil {
		t.Fatal("duplicate create action: expected refusal")
	}
	if err := update(map[string]interface{}{
		"TableName":      "GtPinTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{"Delete": map[string]interface{}{"RegionName": "us-west-2"}}},
	}); err == nil {
		t.Fatal("delete of a non-member: expected refusal")
	}

	// The legacy operation keeps requiring an existing record: the seeding
	// belongs to the modern surface alone.
	if _, err := svc.UpdateGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "NoRecordTable",
		"ReplicaUpdates":  []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}}},
	}}); !errors.Is(err, ErrGlobalTableNotFound) {
		t.Fatalf("legacy update on a recordless table: expected ErrGlobalTableNotFound, got %v", err)
	}
}

// TestReplicaConversionSeedRaceAnswersMappedError pins the seed race's
// error mapping: two concurrent replica-conversion requests on the same
// fresh regional table race the record seeding, and every losing answer
// is a mapped family error (GlobalTableAlreadyExists from the seeded
// Create race, ReplicaAlreadyExists from the membership Update race) or
// success — never the raw storage error the raced Create returns.
func TestReplicaConversionSeedRaceAnswersMappedError(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	ctx := context.Background()

	convert := func(name string) error {
		_, err := svc.UpdateTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":      name,
			"ReplicaUpdates": []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}}},
		}})
		return err
	}

	for round := 0; round < 25; round++ {
		name := "GtRaceTable-" + time.Now().Format("150405.000000000")
		for _, region := range []string{"us-east-1", "eu-west-1"} {
			store, err := svc.GetCachedStoreForRegion(region)
			if err != nil {
				t.Fatalf("region store %s: %v", region, err)
			}
			if _, err := store.Tables().Create(dbstore.CreateTableParams{
				Name:      name,
				KeySchema: []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
				AttributeDefinitions: []*dbstore.AttributeDefinition{
					{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
				},
				BillingMode:         dbstore.BillingModePayPerRequest,
				StreamSpecification: &dbstore.StreamSpecification{StreamEnabled: true, StreamViewType: dbstore.StreamViewTypeNewAndOldImages},
			}); err != nil {
				t.Fatalf("round %d: create table in %s: %v", round, region, err)
			}
		}

		start := make(chan struct{})
		errs := make([]error, 2)
		var wg sync.WaitGroup
		for i := range errs {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				<-start
				errs[i] = convert(name)
			}(i)
		}
		close(start)
		wg.Wait()

		for i, err := range errs {
			if err == nil || errors.Is(err, ErrGlobalTableAlreadyExists) || errors.Is(err, ErrReplicaAlreadyExists) {
				continue
			}
			t.Fatalf("round %d caller %d: unmapped error from the seed race: %v", round, i, err)
		}
	}
}
