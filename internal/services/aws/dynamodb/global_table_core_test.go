package dynamodb

import (
	"context"
	"errors"
	"testing"
	"time"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// globalTablePlaneFixture builds a two-region plane whose regions both hold
// an identical empty streamed table named like the global table — the
// qualifying shape validateGlobalTableReplica demands — plus one request
// context per member region.
func globalTablePlaneFixture(t *testing.T) (*DynamoDBService, *request.RequestContext, *request.RequestContext, *storage.RegionStorageManager) {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)

	for _, region := range []string{"us-east-1", "eu-west-1"} {
		store, err := svc.GetCachedStoreForRegion(region)
		if err != nil {
			t.Fatalf("region store %s: %v", region, err)
		}
		if _, err := store.Tables().Create(dbstore.CreateTableParams{
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
			t.Fatalf("create replica table in %s: %v", region, err)
		}
	}
	return svc,
		request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1"),
		request.NewRequestContext(context.Background(), sm, "123456789012", "eu-west-1"),
		sm
}

// createPinGlobalTable drives CreateGlobalTable through the wire surface
// from the first member region.
func createPinGlobalTable(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext) {
	t.Helper()
	if _, err := svc.CreateGlobalTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicationGroup": []interface{}{
			map[string]interface{}{"RegionName": "us-east-1"},
			map[string]interface{}{"RegionName": "eu-west-1"},
		},
	}}); err != nil {
		t.Fatalf("create global table: %v", err)
	}
}

// waitForReplicatedItem polls GetItem in the given region until the item
// lands: replication is asynchronous, so the pin waits rather than asserting
// immediacy.
func waitForReplicatedItem(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, id string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "GtPinTable",
			"Key":       map[string]interface{}{"id": map[string]interface{}{"S": id}},
		}})
		if err != nil {
			t.Fatalf("replicated read of %s in %s: %v", id, reqCtx.GetRegion(), err)
		}
		if _, present := resp.(map[string]interface{})["Item"]; present {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("item %s never replicated to %s", id, reqCtx.GetRegion())
}

// TestDescribeGlobalTableDistinguishesAbsenceFromStoreFailure pins the
// error split of both Describe cores: a missing record is the documented
// GlobalTableNotFound, while a store failure on an existing record (a
// corrupt payload the proto reader cannot unmarshal) is reported as the
// storage error it is — never as resource absence.
func TestDescribeGlobalTableDistinguishesAbsenceFromStoreFailure(t *testing.T) {
	svc, eastCtx, _, sm := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	// A record that never existed stays the documented not-found verdict.
	for _, op := range []struct {
		name string
		call func() (interface{}, error)
	}{
		{"DescribeGlobalTable", func() (interface{}, error) {
			return svc.DescribeGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"GlobalTableName": "NoSuchGlobalTable",
			}})
		}},
		{"DescribeGlobalTableSettings", func() (interface{}, error) {
			return svc.DescribeGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"GlobalTableName": "NoSuchGlobalTable",
			}})
		}},
	} {
		if _, err := op.call(); !errors.Is(err, ErrGlobalTableNotFound) {
			t.Fatalf("%s on absent record: expected ErrGlobalTableNotFound, got %v", op.name, err)
		}
	}

	// Overwrite the record's payload with bytes no proto reader accepts: the
	// Describe operations must surface the store failure, not flatten it
	// into GlobalTableNotFound.
	globalStorage, err := sm.GetGlobalStorage()
	if err != nil {
		t.Fatalf("global storage: %v", err)
	}
	if err := globalStorage.Bucket("dynamodb_global_tables").Put([]byte("GtPinTable"), []byte("not-a-protobuf-record")); err != nil {
		t.Fatalf("corrupt record: %v", err)
	}

	for _, op := range []struct {
		name string
		call func() (interface{}, error)
	}{
		{"DescribeGlobalTable", func() (interface{}, error) {
			return svc.DescribeGlobalTable(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"GlobalTableName": "GtPinTable",
			}})
		}},
		{"DescribeGlobalTableSettings", func() (interface{}, error) {
			return svc.DescribeGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"GlobalTableName": "GtPinTable",
			}})
		}},
	} {
		_, err := op.call()
		if err == nil || errors.Is(err, ErrGlobalTableNotFound) {
			t.Fatalf("%s on corrupt record: expected the store failure, got %v", op.name, err)
		}
	}
}

// TestUpdateGlobalTableSettingsRejectsFractionalCapacity pins the Long
// contract of the two direct capacity members of UpdateGlobalTableSettings:
// fractional input is a ValidationException — the same contract the
// request's GSI capacity members enforce — not a silently truncated value,
// and integral input still applies.
func TestUpdateGlobalTableSettingsRejectsFractionalCapacity(t *testing.T) {
	svc, eastCtx, westCtx, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	settings := func(params map[string]interface{}) (interface{}, error) {
		return svc.UpdateGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: params})
	}

	if _, err := settings(map[string]interface{}{
		"GlobalTableName":                          "GtPinTable",
		"GlobalTableProvisionedWriteCapacityUnits": 1.5,
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("fractional global write capacity: expected ErrInvalidParameter, got %v", err)
	}

	if _, err := settings(map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaSettingsUpdate": []interface{}{
			map[string]interface{}{"RegionName": "eu-west-1", "ReplicaProvisionedReadCapacityUnits": 2.5},
		},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("fractional replica read capacity: expected ErrInvalidParameter, got %v", err)
	}

	// An over-cap pair is refused before the global record commits — the
	// table plane's caps bound the pair the same request would leave each
	// member table in, on both the replica read side and the global write
	// side (an over-cap value that passed here would commit and then fail
	// every member region's propagation).
	maxRead := float64(dbstore.TableMaxReadCapacityUnits) + 1
	if _, err := settings(map[string]interface{}{
		"GlobalTableName":                          "GtPinTable",
		"GlobalTableBillingMode":                   "PROVISIONED",
		"GlobalTableProvisionedWriteCapacityUnits": 5.0,
		"ReplicaSettingsUpdate": []interface{}{
			map[string]interface{}{"RegionName": "us-east-1", "ReplicaProvisionedReadCapacityUnits": maxRead},
		},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("over-cap replica read capacity: expected ErrInvalidParameter, got %v", err)
	}
	maxWrite := float64(dbstore.TableMaxWriteCapacityUnits) + 1
	if _, err := settings(map[string]interface{}{
		"GlobalTableName":                          "GtPinTable",
		"GlobalTableBillingMode":                   "PROVISIONED",
		"GlobalTableProvisionedWriteCapacityUnits": maxWrite,
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("over-cap global write capacity: expected ErrInvalidParameter, got %v", err)
	}
	for _, reqCtx := range []*request.RequestContext{eastCtx, westCtx} {
		table := describeTableOnModePlane(t, svc, reqCtx, "GtPinTable")
		if summary, ok := table["BillingModeSummary"].(map[string]interface{}); !ok || summary["BillingMode"] != "PAY_PER_REQUEST" {
			t.Fatalf("refused over-cap request changed %s's billing mode: %+v", reqCtx.GetRegion(), table["BillingModeSummary"])
		}
	}

	// Integral values apply and reach the tables: the switch to provisioned
	// lands on every member region's table record, so DescribeTable and
	// DescribeGlobalTableSettings report the same state for a replica.
	if _, err := settings(map[string]interface{}{
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
	}); err != nil {
		t.Fatalf("integral settings update: %v", err)
	}
	resp, err := svc.DescribeGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
	}})
	if err != nil {
		t.Fatalf("describe settings: %v", err)
	}
	for _, rs := range resp.(map[string]interface{})["ReplicaSettings"].([]map[string]interface{}) {
		if got := rs["ReplicaProvisionedWriteCapacityUnits"].(int64); got != 5 {
			t.Fatalf("replica %v write capacity: got %d, want 5", rs["RegionName"], got)
		}
		if rs["RegionName"] == "eu-west-1" {
			if got := rs["ReplicaProvisionedReadCapacityUnits"].(int64); got != 3 {
				t.Fatalf("eu-west-1 read capacity: got %d, want 3", got)
			}
		}
	}
	for _, tc := range []struct {
		reqCtx  *request.RequestContext
		read    int64
		gsiRead int64
	}{
		{eastCtx, 7, 2},
		{westCtx, 3, 6},
	} {
		table := describeTableOnModePlane(t, svc, tc.reqCtx, "GtPinTable")
		summary, ok := table["BillingModeSummary"].(map[string]interface{})
		if !ok || summary["BillingMode"] != "PROVISIONED" {
			t.Fatalf("table in %s billing mode summary = %+v, want PROVISIONED", tc.reqCtx.GetRegion(), table["BillingModeSummary"])
		}
		pt := table["ProvisionedThroughput"].(map[string]interface{})
		if got := asInt64(t, pt["ReadCapacityUnits"]); got != tc.read {
			t.Fatalf("table in %s read capacity: got %d, want %d", tc.reqCtx.GetRegion(), got, tc.read)
		}
		if got := asInt64(t, pt["WriteCapacityUnits"]); got != 5 {
			t.Fatalf("table in %s write capacity: got %d, want 5", tc.reqCtx.GetRegion(), got)
		}
		gsi := table["GlobalSecondaryIndexes"].([]map[string]interface{})[0]
		gpt := gsi["ProvisionedThroughput"].(map[string]interface{})
		if got := asInt64(t, gpt["ReadCapacityUnits"]); got != tc.gsiRead {
			t.Fatalf("gsi in %s read capacity: got %d, want %d", tc.reqCtx.GetRegion(), got, tc.gsiRead)
		}
		if got := asInt64(t, gpt["WriteCapacityUnits"]); got != 4 {
			t.Fatalf("gsi in %s write capacity: got %d, want 4", tc.reqCtx.GetRegion(), got)
		}
	}

	// Switching back to on-demand discards the provisioned settings on
	// every replica table, the same semantics the table plane's own update
	// applies.
	if _, err := settings(map[string]interface{}{
		"GlobalTableName":        "GtPinTable",
		"GlobalTableBillingMode": "PAY_PER_REQUEST",
	}); err != nil {
		t.Fatalf("on-demand switch: %v", err)
	}
	for _, reqCtx := range []*request.RequestContext{eastCtx, westCtx} {
		table := describeTableOnModePlane(t, svc, reqCtx, "GtPinTable")
		if pt := table["ProvisionedThroughput"].(map[string]interface{}); asInt64(t, pt["ReadCapacityUnits"]) != 0 || asInt64(t, pt["WriteCapacityUnits"]) != 0 {
			t.Fatalf("table in %s kept provisioned capacity after the on-demand switch", reqCtx.GetRegion())
		}
		if gsi := table["GlobalSecondaryIndexes"].([]map[string]interface{})[0]; gsi["ProvisionedThroughput"] != nil {
			gpt := gsi["ProvisionedThroughput"].(map[string]interface{})
			if asInt64(t, gpt["ReadCapacityUnits"]) != 0 || asInt64(t, gpt["WriteCapacityUnits"]) != 0 {
				t.Fatalf("gsi in %s kept provisioned capacity after the on-demand switch", reqCtx.GetRegion())
			}
		}
	}
}

// TestUpdateGlobalTableSettingsEnforcesTheSwitchRateWindow pins the switch
// quota on the global capacity plane: the fifth provisioned-to-on-demand
// switch inside a rolling 24-hour window is refused before the global
// record commits (the members keep their provisioned state), and the
// propagation that lands a switch on a member record re-enforces the same
// quota under the region's own lock — a member table that crossed the
// window after the validation read is refused, not stamped past the quota.
func TestUpdateGlobalTableSettingsEnforcesTheSwitchRateWindow(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	// A member table at the quota: provisioned, holding a window-full of
	// switch stamps from the last hour.
	stampToQuota := func(table *dbstore.Table) {
		table.BillingMode = dbstore.BillingModeProvisioned
		table.ProvisionedThroughput = &dbstore.ProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: 5}
		table.BillingModeSwitches = make([]time.Time, dbstore.BillingModeSwitchesPerDay)
		for i := range table.BillingModeSwitches {
			table.BillingModeSwitches[i] = time.Now().UTC().Add(-time.Minute)
		}
	}

	// The propagation half, directly: landing the switch on a member record
	// that sits at the quota is refused and leaves the record untouched.
	atQuota := &dbstore.Table{}
	stampToQuota(atQuota)
	if err := applyGlobalCapacityToTable(atQuota, string(dbstore.BillingModePayPerRequest), false, 0, nil, replicaSettingsUpdate{}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("propagation onto a member at the quota: expected ErrInvalidParameter, got %v", err)
	}
	if atQuota.BillingMode != dbstore.BillingModeProvisioned || len(atQuota.BillingModeSwitches) != dbstore.BillingModeSwitchesPerDay {
		t.Fatalf("refused propagation mutated the record: mode=%v stamps=%d", atQuota.BillingMode, len(atQuota.BillingModeSwitches))
	}

	// The request half, end to end: both member tables at the quota, the
	// on-demand settings update is refused before the global record
	// commits and the members keep their provisioned state.
	for _, region := range []string{"us-east-1", "eu-west-1"} {
		store, err := svc.GetCachedStoreForRegion(region)
		if err != nil {
			t.Fatalf("region store %s: %v", region, err)
		}
		if _, err := store.Tables().Update("GtPinTable", func(table *dbstore.Table) error {
			stampToQuota(table)
			return nil
		}); err != nil {
			t.Fatalf("stamp member table in %s: %v", region, err)
		}
	}
	if _, err := svc.UpdateGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName":        "GtPinTable",
		"GlobalTableBillingMode": "PAY_PER_REQUEST",
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("fifth switch inside the window: expected ErrInvalidParameter, got %v", err)
	}
	for _, region := range []string{"us-east-1", "eu-west-1"} {
		store, err := svc.GetCachedStoreForRegion(region)
		if err != nil {
			t.Fatalf("region store %s: %v", region, err)
		}
		table, err := store.Tables().Get("GtPinTable")
		if err != nil {
			t.Fatalf("read member table in %s: %v", region, err)
		}
		if table.BillingMode != dbstore.BillingModeProvisioned || len(table.BillingModeSwitches) != dbstore.BillingModeSwitchesPerDay {
			t.Fatalf("member table in %s after the refusal: mode=%v stamps=%d", region, table.BillingMode, len(table.BillingModeSwitches))
		}
	}
}

// TestUpdateGlobalTableSettingsIndexAndCapacityValidation pins the
// acceptance rules the settings update inherits from the tables it
// reconfigures: a per-index settings update naming an index no member
// table holds is the model's IndexNotFoundException (never a recorded
// settings entry for a non-existent index), capacity members carried into
// an on-demand result are refused, and a provisioned result whose capacity
// pair or index pairs cannot be completed from the request and the
// existing state is refused before anything commits.
func TestUpdateGlobalTableSettingsIndexAndCapacityValidation(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	settings := func(params map[string]interface{}) (interface{}, error) {
		return svc.UpdateGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: params})
	}

	// A global write-side index update naming an index no member table
	// holds.
	if _, err := settings(map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"GlobalTableGlobalSecondaryIndexSettingsUpdate": []interface{}{
			map[string]interface{}{"IndexName": "noSuchIndex", "ProvisionedWriteCapacityUnits": 4.0},
		},
	}); !errors.Is(err, ErrIndexNotFound) {
		t.Fatalf("global write-side update for a non-existent index: err = %v, want ErrIndexNotFound", err)
	}

	// A replica read-side index update naming an index that replica's
	// table does not hold.
	if _, err := settings(map[string]interface{}{
		"GlobalTableName": "GtPinTable",
		"ReplicaSettingsUpdate": []interface{}{
			map[string]interface{}{
				"RegionName": "us-east-1",
				"ReplicaGlobalSecondaryIndexSettingsUpdate": []interface{}{
					map[string]interface{}{"IndexName": "noSuchIndex", "ProvisionedReadCapacityUnits": 2.0},
				},
			},
		},
	}); !errors.Is(err, ErrIndexNotFound) {
		t.Fatalf("replica read-side update for a non-existent index: err = %v, want ErrIndexNotFound", err)
	}

	// Capacity members carried into an on-demand result are a
	// contradiction.
	if _, err := settings(map[string]interface{}{
		"GlobalTableName":                          "GtPinTable",
		"GlobalTableBillingMode":                   "PAY_PER_REQUEST",
		"GlobalTableProvisionedWriteCapacityUnits": 5.0,
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("write units into an on-demand result: err = %v, want ErrInvalidParameter", err)
	}

	// A provisioned result must complete every capacity pair: write units
	// alone leave the read side of the pair — and of the fixture's on-demand
	// tables' index — unresolved.
	if _, err := settings(map[string]interface{}{
		"GlobalTableName":                          "GtPinTable",
		"GlobalTableBillingMode":                   "PROVISIONED",
		"GlobalTableProvisionedWriteCapacityUnits": 5.0,
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("incomplete provisioned pair: err = %v, want ErrInvalidParameter", err)
	}

	// A refused request mutates nothing: the tables are still on-demand
	// and the settings record carries no capacity.
	eastTable := describeTableOnModePlane(t, svc, eastCtx, "GtPinTable")
	if summary, ok := eastTable["BillingModeSummary"].(map[string]interface{}); ok && summary["BillingMode"] != "PAY_PER_REQUEST" {
		t.Fatalf("refused requests changed the table's billing mode: %+v", summary)
	}
	resp, err := svc.DescribeGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"GlobalTableName": "GtPinTable",
	}})
	if err != nil {
		t.Fatalf("describe settings: %v", err)
	}
	for _, rs := range resp.(map[string]interface{})["ReplicaSettings"].([]map[string]interface{}) {
		if got := rs["ReplicaProvisionedWriteCapacityUnits"].(int64); got != 0 {
			t.Fatalf("refused request left write capacity %d on replica %v", got, rs["RegionName"])
		}
	}
}

// TestUpdateGlobalTableSettingsReplicaListBounds pins the model's 1-50
// entry bound on the ReplicaSettingsUpdate list through the wire surface
// the SDK client itself reaches.
func TestUpdateGlobalTableSettingsReplicaListBounds(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	ctx := context.Background()

	settings := func(params map[string]interface{}) (interface{}, error) {
		return svc.UpdateGlobalTableSettings(ctx, eastCtx, &request.ParsedRequest{Parameters: params})
	}

	// An explicitly empty list is a request, not an absent member.
	if _, err := settings(map[string]interface{}{
		"GlobalTableName":       "GtPinTable",
		"ReplicaSettingsUpdate": []interface{}{},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty replica settings update list: expected ErrInvalidParameter, got %v", err)
	}

	updates := make([]interface{}, 51)
	for i := range updates {
		updates[i] = map[string]interface{}{"RegionName": "us-east-1"}
	}
	if _, err := settings(map[string]interface{}{
		"GlobalTableName":       "GtPinTable",
		"ReplicaSettingsUpdate": updates,
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("51-entry replica settings update list: expected ErrInvalidParameter, got %v", err)
	}
}

// TestUpdateTableRefusesUnsupportedReplicaFamilyMembers pins the policy on
// the remaining replica-family members: EVENTUAL consistency alongside
// Create actions is accepted as the platform's own replication mode, STRONG
// consistency and the consistency member without Create actions are
// refused, and witness updates and the settings-replication mode — features
// with no substrate on this platform — are refused rather than accepted
// without effect.
func TestUpdateTableRefusesUnsupportedReplicaFamilyMembers(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	ctx := context.Background()

	update := func(params map[string]interface{}) error {
		_, err := svc.UpdateTable(ctx, eastCtx, &request.ParsedRequest{Parameters: params})
		return err
	}
	createAction := []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}}}

	if err := update(map[string]interface{}{"TableName": "GtPinTable", "MultiRegionConsistency": "STRONG", "ReplicaUpdates": createAction}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("STRONG consistency: expected ErrInvalidParameter, got %v", err)
	}
	if err := update(map[string]interface{}{"TableName": "GtPinTable", "MultiRegionConsistency": "EVENTUAL"}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("EVENTUAL without Create actions: expected ErrInvalidParameter, got %v", err)
	}
	if err := update(map[string]interface{}{"TableName": "GtPinTable", "GlobalTableWitnessUpdates": []interface{}{map[string]interface{}{"Create": map[string]interface{}{"RegionName": "eu-west-1"}}}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("witness update: expected ErrInvalidParameter, got %v", err)
	}
	if err := update(map[string]interface{}{"TableName": "GtPinTable", "GlobalTableSettingsReplicationMode": "ENABLED"}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("settings replication mode: expected ErrInvalidParameter, got %v", err)
	}

	// EVENTUAL alongside Create actions is the platform's own mode: the
	// request stands.
	if err := update(map[string]interface{}{"TableName": "GtPinTable", "MultiRegionConsistency": "EVENTUAL", "ReplicaUpdates": createAction}); err != nil {
		t.Fatalf("EVENTUAL with Create actions: %v", err)
	}
}
