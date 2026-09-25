// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// createRestorableBackup creates a PAY_PER_REQUEST hash table holding one
// item and an AVAILABLE backup of it, returning the persisted backup.
func createRestorableBackup(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, table, backupName string) *dbstore.Backup {
	t.Helper()
	ctx := context.Background()
	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            table,
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create %s: %v", table, err)
	}
	if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": table,
		"Item":      map[string]interface{}{"pk": map[string]interface{}{"S": "kept"}},
	}}); err != nil {
		t.Fatalf("put item into %s: %v", table, err)
	}
	backup, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName":  table,
		"BackupName": backupName,
	}})
	if err != nil {
		t.Fatalf("create backup %s: %v", backupName, err)
	}
	return backup
}

// TestBackupSnapshotScanIsSerialisable pins the backup collection's
// serialisation: the item snapshot reads through a store View — a pebble
// snapshot — so the whole walk observes one commit point. The pin commits a
// write inside the open View (a commit that lands after the snapshot point
// but before the walk reads, exactly what a concurrent writer produces) and
// asserts the walk is all-or-nothing: everything committed before the point
// is present in full, the later commit is absent — a direct-bucket scan
// would read the later item and fail this.
func TestBackupSnapshotScanIsSerialisable(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "SnapTbl",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create SnapTbl: %v", err)
	}
	if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "SnapTbl",
		"Item":      map[string]interface{}{"pk": map[string]interface{}{"S": "before-point"}},
	}}); err != nil {
		t.Fatalf("put before-point: %v", err)
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}

	seen := map[string]bool{}
	err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		// A write committed after the snapshot point, before the walk.
		if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "SnapTbl",
			"Item":      map[string]interface{}{"pk": map[string]interface{}{"S": "after-point"}},
		}}); err != nil {
			return fmt.Errorf("put after-point inside view: %w", err)
		}
		return txn.Scan("SnapTbl", func(item *dbstore.Item) error {
			seen[*item.Key["pk"].S] = true
			return nil
		})
	})
	if err != nil {
		t.Fatalf("view scan: %v", err)
	}

	if !seen["before-point"] {
		t.Fatalf("snapshot lost the pre-point item, saw %v", seen)
	}
	if seen["after-point"] {
		t.Fatalf("snapshot leaked the post-point commit, saw %v", seen)
	}

	// The later commit is durable — a fresh walk outside any snapshot sees
	// both items, so the absence above is the snapshot, not a lost write.
	if err := store.Items().Scan("SnapTbl", func(item *dbstore.Item) error {
		if !seen[*item.Key["pk"].S] {
			seen[*item.Key["pk"].S] = true
		}
		return nil
	}); err != nil {
		t.Fatalf("direct scan: %v", err)
	}
	if !seen["after-point"] {
		t.Fatalf("post-point commit was not durable, saw %v", seen)
	}
}

// TestRestoreRejectsIncompleteBackups pins the restore-completeness guards:
// a backup still in CREATING (the crash window between the record's initial
// write and the snapshot save) is refused before the target table is
// created, an AVAILABLE record with no snapshot is treated as corrupted
// rather than empty so the restore fails instead of completing an empty
// ACTIVE table, and an intact backup still restores with its item and
// summary intact.
func TestRestoreRejectsIncompleteBackups(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	backup := createRestorableBackup(t, svc, reqCtx, "RestSrc", "incomplete-bk")

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}

	backup.BackupStatus = dbstore.BackupStatusCreating
	if err := store.Backups().Put(backup); err != nil {
		t.Fatalf("mark backup CREATING: %v", err)
	}
	if _, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn:       backup.BackupArn,
		TargetTableName: "RestoredCreating",
	}); !errors.Is(err, ErrBackupNotAvailable) {
		t.Fatalf("restoring a CREATING backup: err = %v, want ErrBackupNotAvailable", err)
	}
	if store.Tables().Exists("RestoredCreating") {
		t.Fatal("refused restore created its target table")
	}

	backup.BackupStatus = dbstore.BackupStatusAvailable
	if err := store.Backups().Put(backup); err != nil {
		t.Fatalf("mark backup AVAILABLE: %v", err)
	}
	if err := store.Backups().DeleteSnapshot(backup.BackupArn); err != nil {
		t.Fatalf("drop snapshot: %v", err)
	}
	if _, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn:       backup.BackupArn,
		TargetTableName: "RestoredSnapshotless",
	}); err == nil {
		t.Fatal("restoring a snapshotless backup reported success")
	}
	if half, getErr := store.Tables().Get("RestoredSnapshotless"); getErr == nil && half.Status == dbstore.TableStatusActive {
		t.Fatalf("failed restore left an ACTIVE table: %+v", half)
	}

	// Positive control: an intact backup restores to an ACTIVE table with
	// the snapshot's item present and the source recorded in the summary.
	intact := createRestorableBackup(t, svc, reqCtx, "RestSrc2", "intact-bk")
	restored, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn:       intact.BackupArn,
		TargetTableName: "RestoredIntact",
	})
	if err != nil {
		t.Fatalf("restore intact backup: %v", err)
	}
	if restored.Status != dbstore.TableStatusActive || restored.RestoreSummary == nil ||
		restored.RestoreSummary.SourceBackupArn != intact.BackupArn {
		t.Fatalf("restored table = %+v", restored)
	}
	item, err := svc.GetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "RestoredIntact",
		"Key":       map[string]interface{}{"pk": map[string]interface{}{"S": "kept"}},
	}})
	if err != nil {
		t.Fatalf("get restored item: %v", err)
	}
	if resp := item.(map[string]interface{}); resp["Item"] == nil {
		t.Fatalf("restored item missing: %+v", resp)
	}
}

// TestRestoreFromBackupAppliesOverrideMembers pins the override family the
// from-backup restore owes: the billing mode and SSE override apply, a
// present-empty GSI override list excludes every global secondary index
// while an absent member keeps them, the LSI override selects from the
// backup's existing indexes, the on-demand pair applies on an on-demand
// resulting mode and is rejected on a provisioned one, a mode switch back
// to provisioned requires the initial throughput values, and overrides
// naming unknown indexes are validation errors.
func TestRestoreFromBackupAppliesOverrideMembers(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	// A provisioned source with one GSI and one LSI.
	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "OvProvSrc",
		"KeySchema": []interface{}{
			map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"},
		},
		"AttributeDefinitions": []interface{}{
			map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"},
			map[string]interface{}{"AttributeName": "sk", "AttributeType": "S"},
			map[string]interface{}{"AttributeName": "gsik", "AttributeType": "S"},
		},
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 3.0, "WriteCapacityUnits": 4.0},
		"GlobalSecondaryIndexes": []interface{}{map[string]interface{}{
			"IndexName":             "gsi-1",
			"KeySchema":             []interface{}{map[string]interface{}{"AttributeName": "gsik", "KeyType": "HASH"}},
			"Projection":            map[string]interface{}{"ProjectionType": "ALL"},
			"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 2.0, "WriteCapacityUnits": 3.0},
		}},
		"LocalSecondaryIndexes": []interface{}{map[string]interface{}{
			"IndexName": "lsi-1",
			"KeySchema": []interface{}{
				map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"},
				map[string]interface{}{"AttributeName": "sk", "KeyType": "RANGE"},
			},
			"Projection": map[string]interface{}{"ProjectionType": "ALL"},
		}},
	}}); err != nil {
		t.Fatalf("create provisioned source: %v", err)
	}
	provBackup, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName": "OvProvSrc", "BackupName": "ov-prov-bk",
	}})
	if err != nil {
		t.Fatalf("backup provisioned source: %v", err)
	}

	// An on-demand source for the mode-switch and on-demand pairs.
	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "OvOndSrc",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create on-demand source: %v", err)
	}
	ondBackup, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName": "OvOndSrc", "BackupName": "ov-ond-bk",
	}})
	if err != nil {
		t.Fatalf("backup on-demand source: %v", err)
	}

	restore := func(backupArn, target string, extra map[string]interface{}) (map[string]interface{}, error) {
		params := map[string]interface{}{"BackupArn": backupArn, "TargetTableName": target}
		for k, v := range extra {
			params[k] = v
		}
		if _, err := svc.RestoreTableFromBackup(ctx, reqCtx, &request.ParsedRequest{Parameters: params}); err != nil {
			return nil, err
		}
		return describeTableOnModePlane(t, svc, reqCtx, target), nil
	}

	// Mode switch to on-demand, GSI exclusion via a present-empty list,
	// LSI kept, and SSE override applied.
	restored, err := restore(provBackup.BackupArn, "OvRestoredOnDemand", map[string]interface{}{
		"BillingModeOverride":          "PAY_PER_REQUEST",
		"GlobalSecondaryIndexOverride": []interface{}{},
		"SSESpecificationOverride":     map[string]interface{}{"Enabled": true},
	})
	if err != nil {
		t.Fatalf("restore with overrides: %v", err)
	}
	if restored["BillingModeSummary"].(map[string]interface{})["BillingMode"] != "PAY_PER_REQUEST" {
		t.Fatalf("restored billing mode: %+v", restored["BillingModeSummary"])
	}
	if restored["GlobalSecondaryIndexes"] != nil {
		t.Fatalf("present-empty GSI override left indexes: %+v", restored["GlobalSecondaryIndexes"])
	}
	if lsis := restored["LocalSecondaryIndexes"].([]map[string]interface{}); len(lsis) != 1 || lsis[0]["IndexName"] != "lsi-1" {
		t.Fatalf("LSI override dropped the backup's index: %+v", restored["LocalSecondaryIndexes"])
	}
	sse, ok := restored["SSEDescription"].(map[string]interface{})
	if !ok || sse["SSEType"] != "KMS" {
		t.Fatalf("SSE override not applied: %+v", restored["SSEDescription"])
	}

	// An absent GSI override keeps the backup's indexes.
	restored, err = restore(provBackup.BackupArn, "OvRestoredKept", nil)
	if err != nil {
		t.Fatalf("restore without overrides: %v", err)
	}
	if gsis := restored["GlobalSecondaryIndexes"].([]map[string]interface{}); len(gsis) != 1 || gsis[0]["IndexName"] != "gsi-1" {
		t.Fatalf("absent GSI override dropped indexes: %+v", restored["GlobalSecondaryIndexes"])
	}

	// A mode switch to on-demand without a GSI override keeps the backup's
	// index but clears its provisioned settings — secondary indexes inherit
	// the table's capacity mode.
	restored, err = restore(provBackup.BackupArn, "OvRestoredModeOnly", map[string]interface{}{
		"BillingModeOverride": "PAY_PER_REQUEST",
	})
	if err != nil {
		t.Fatalf("restore with mode override only: %v", err)
	}
	gsis := restored["GlobalSecondaryIndexes"].([]map[string]interface{})
	if len(gsis) != 1 || gsis[0]["IndexName"] != "gsi-1" {
		t.Fatalf("mode-only switch dropped indexes: %+v", restored["GlobalSecondaryIndexes"])
	}
	gsiPt := gsis[0]["ProvisionedThroughput"].(map[string]interface{})
	if asInt64(t, gsiPt["ReadCapacityUnits"]) != 0 || asInt64(t, gsiPt["WriteCapacityUnits"]) != 0 {
		t.Fatalf("mode switch left the index provisioned: %+v", gsiPt)
	}

	// The on-demand pair applies on an on-demand resulting mode.
	restored, err = restore(ondBackup.BackupArn, "OvRestoredCaps", map[string]interface{}{
		"OnDemandThroughputOverride": map[string]interface{}{"MaxReadRequestUnits": 9.0, "MaxWriteRequestUnits": 8.0},
	})
	if err != nil {
		t.Fatalf("restore with on-demand override: %v", err)
	}
	odt := restored["OnDemandThroughput"].(map[string]interface{})
	if asInt64(t, odt["MaxReadRequestUnits"]) != 9 || asInt64(t, odt["MaxWriteRequestUnits"]) != 8 {
		t.Fatalf("on-demand override not applied: %+v", odt)
	}

	// The on-demand pair is rejected on a provisioned resulting mode.
	if _, err := restore(provBackup.BackupArn, "OvRestoredBadCaps", map[string]interface{}{
		"OnDemandThroughputOverride": map[string]interface{}{"MaxReadRequestUnits": 9.0},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("on-demand override on provisioned base: err = %v, want ErrInvalidParameter", err)
	}

	// Switching the restored table back to provisioned requires the
	// initial throughput values.
	if _, err := restore(ondBackup.BackupArn, "OvRestoredSwitchBack", map[string]interface{}{
		"BillingModeOverride": "PROVISIONED",
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("mode switch back without throughput: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := restore(ondBackup.BackupArn, "OvRestoredSwitchBackOk", map[string]interface{}{
		"BillingModeOverride":           "PROVISIONED",
		"ProvisionedThroughputOverride": map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 6.0},
	}); err != nil {
		t.Fatalf("mode switch back with throughput: %v", err)
	}

	// Overrides may only name indexes the backup holds.
	if _, err := restore(provBackup.BackupArn, "OvRestoredUnknownGsi", map[string]interface{}{
		"GlobalSecondaryIndexOverride": []interface{}{map[string]interface{}{"IndexName": "nope"}},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("unknown GSI override: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := restore(provBackup.BackupArn, "OvRestoredUnknownLsi", map[string]interface{}{
		"LocalSecondaryIndexOverride": []interface{}{map[string]interface{}{"IndexName": "nope"}},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("unknown LSI override: err = %v, want ErrInvalidParameter", err)
	}

	// Override entries carry member-level settings: a known GSI keeps its
	// place with a replacement projection (and its identical key schema
	// accepted), the LSI override applies its projection the same way, and
	// a changed key schema is rejected.
	restored, err = restore(provBackup.BackupArn, "OvRestoredProj", map[string]interface{}{
		"GlobalSecondaryIndexOverride": []interface{}{map[string]interface{}{
			"IndexName":  "gsi-1",
			"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "gsik", "KeyType": "HASH"}},
			"Projection": map[string]interface{}{"ProjectionType": "INCLUDE", "NonKeyAttributes": []interface{}{"a", "b"}},
		}},
		"LocalSecondaryIndexOverride": []interface{}{map[string]interface{}{
			"IndexName":  "lsi-1",
			"Projection": map[string]interface{}{"ProjectionType": "KEYS_ONLY"},
		}},
	})
	if err != nil {
		t.Fatalf("member-level overrides: %v", err)
	}
	gsis = restored["GlobalSecondaryIndexes"].([]map[string]interface{})
	if len(gsis) != 1 {
		t.Fatalf("member-level GSI override count: %+v", gsis)
	}
	gsiProj := gsis[0]["Projection"].(map[string]interface{})
	nka, nkaOK := gsiProj["NonKeyAttributes"].([]string)
	if gsiProj["ProjectionType"] != "INCLUDE" || !nkaOK || len(nka) != 2 {
		t.Fatalf("GSI projection override not applied: %+v", gsiProj)
	}
	lsis := restored["LocalSecondaryIndexes"].([]map[string]interface{})
	if len(lsis) != 1 || lsis[0]["Projection"].(map[string]interface{})["ProjectionType"] != "KEYS_ONLY" {
		t.Fatalf("LSI projection override not applied: %+v", lsis)
	}

	// An on-demand source with a GSI: the index-level on-demand and warm
	// pairs of an override entry must reach the restored index — the parser
	// accepts them and the selection copy is their only path onto the
	// restored record.
	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "OvOndGsiSrc",
		"KeySchema": []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{
			map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"},
			map[string]interface{}{"AttributeName": "gsik", "AttributeType": "S"},
		},
		"BillingMode": "PAY_PER_REQUEST",
		"GlobalSecondaryIndexes": []interface{}{map[string]interface{}{
			"IndexName":  "gsi-ond",
			"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "gsik", "KeyType": "HASH"}},
			"Projection": map[string]interface{}{"ProjectionType": "ALL"},
		}},
	}}); err != nil {
		t.Fatalf("create on-demand GSI source: %v", err)
	}
	ondGsiBackup, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName": "OvOndGsiSrc", "BackupName": "ov-ond-gsi-bk",
	}})
	if err != nil {
		t.Fatalf("backup on-demand GSI source: %v", err)
	}
	restored, err = restore(ondGsiBackup.BackupArn, "OvRestoredGsiCaps", map[string]interface{}{
		"GlobalSecondaryIndexOverride": []interface{}{map[string]interface{}{
			"IndexName":          "gsi-ond",
			"OnDemandThroughput": map[string]interface{}{"MaxReadRequestUnits": 40.0, "MaxWriteRequestUnits": 30.0},
			"WarmThroughput":     map[string]interface{}{"ReadUnitsPerSecond": 41.0, "WriteUnitsPerSecond": 31.0},
		}},
	})
	if err != nil {
		t.Fatalf("restore with index on-demand/warm override: %v", err)
	}
	gsis = restored["GlobalSecondaryIndexes"].([]map[string]interface{})
	if len(gsis) != 1 || gsis[0]["IndexName"] != "gsi-ond" {
		t.Fatalf("index on-demand override selection: %+v", restored["GlobalSecondaryIndexes"])
	}
	gsiOdt := gsis[0]["OnDemandThroughput"].(map[string]interface{})
	if asInt64(t, gsiOdt["MaxReadRequestUnits"]) != 40 || asInt64(t, gsiOdt["MaxWriteRequestUnits"]) != 30 {
		t.Fatalf("GSI on-demand override not applied: %+v", gsiOdt)
	}
	gsiWarm := gsis[0]["WarmThroughput"].(map[string]interface{})
	if asInt64(t, gsiWarm["ReadUnitsPerSecond"]) != 41 || asInt64(t, gsiWarm["WriteUnitsPerSecond"]) != 31 {
		t.Fatalf("GSI warm override not applied: %+v", gsiWarm)
	}

	if _, err := restore(provBackup.BackupArn, "OvRestoredBadSchema", map[string]interface{}{
		"GlobalSecondaryIndexOverride": []interface{}{map[string]interface{}{
			"IndexName": "gsi-1",
			"KeySchema": []interface{}{map[string]interface{}{"AttributeName": "other", "KeyType": "HASH"}},
		}},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("changed key schema override: err = %v, want ErrInvalidParameter", err)
	}

	// The mode override is a closed enum on the from-backup surface: a value
	// outside PROVISIONED | PAY_PER_REQUEST is rejected before any member is
	// applied — no target table is left behind carrying an invalid mode.
	for _, bad := range []string{"garbage", "provisioned"} {
		if _, err := restore(ondBackup.BackupArn, "OvRestoredBadMode", map[string]interface{}{
			"BillingModeOverride": bad,
		}); !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("billing mode override %q: err = %v, want ErrInvalidParameter", bad, err)
		}
		if store, err := svc.store(reqCtx); err == nil && store.Tables().Exists("OvRestoredBadMode") {
			t.Fatalf("billing mode override %q left a target table behind", bad)
		}
	}
}

// TestPITRRestoreAppliesOnDemandAndVectorOverrides pins the two override
// members the point-in-time restore owed: the on-demand maximum pair, and
// the vector index override's exclusion contract (absent keeps all, a
// present-empty list excludes all, a provided definition must match the
// source's).
func TestPITRRestoreAppliesOnDemandAndVectorOverrides(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "PitrOvSrc",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
		"VectorIndexes":        []interface{}{vectorIndexParams("vec", "embedding", float64(2), "COSINE")},
	}}); err != nil {
		t.Fatalf("create vector source: %v", err)
	}
	if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "PitrOvSrc",
		"Item": map[string]interface{}{
			"pk":        map[string]interface{}{"S": "one"},
			"embedding": map[string]interface{}{"L": []interface{}{map[string]interface{}{"N": "1"}, map[string]interface{}{"N": "0"}}},
		},
	}}); err != nil {
		t.Fatalf("put vector item: %v", err)
	}
	if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "PitrOvSrc",
		"PointInTimeRecoverySpecification": map[string]interface{}{
			"PointInTimeRecoveryEnabled": true,
		},
	}}); err != nil {
		t.Fatalf("enable pitr: %v", err)
	}

	pitrRestore := func(target string, extra map[string]interface{}) (map[string]interface{}, error) {
		params := map[string]interface{}{"SourceTableName": "PitrOvSrc", "TargetTableName": target, "UseLatestRestorableTime": true}
		for k, v := range extra {
			params[k] = v
		}
		if _, err := svc.RestoreTableToPointInTime(ctx, reqCtx, &request.ParsedRequest{Parameters: params}); err != nil {
			return nil, err
		}
		return describeTableOnModePlane(t, svc, reqCtx, target), nil
	}

	// The on-demand pair applies and the vector index carries over when
	// the override is absent.
	restored, err := pitrRestore("PitrOvKept", map[string]interface{}{
		"OnDemandThroughputOverride": map[string]interface{}{"MaxReadRequestUnits": 7.0},
	})
	if err != nil {
		t.Fatalf("pitr restore with on-demand override: %v", err)
	}
	odt := restored["OnDemandThroughput"].(map[string]interface{})
	if asInt64(t, odt["MaxReadRequestUnits"]) != 7 {
		t.Fatalf("on-demand override not applied: %+v", odt)
	}
	vis := restored["VectorIndexes"].([]map[string]interface{})
	if len(vis) != 1 || vis[0]["IndexName"] != "vec" {
		t.Fatalf("vector index not carried over: %+v", restored["VectorIndexes"])
	}

	// The mode override is the same closed enum on this surface: an invalid
	// value is rejected before the target is created.
	if _, err := pitrRestore("PitrOvBadMode", map[string]interface{}{
		"BillingModeOverride": "garbage",
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("PITR billing mode override: err = %v, want ErrInvalidParameter", err)
	}

	// A present-empty override excludes every vector index.
	restored, err = pitrRestore("PitrOvExcluded", map[string]interface{}{
		"VectorIndexOverride": []interface{}{},
	})
	if err != nil {
		t.Fatalf("pitr restore with empty vector override: %v", err)
	}
	if restored["VectorIndexes"] != nil {
		t.Fatalf("present-empty vector override left indexes: %+v", restored["VectorIndexes"])
	}

	// A provided definition must match the source's.
	if _, err := pitrRestore("PitrOvMismatch", map[string]interface{}{
		"VectorIndexOverride": []interface{}{vectorIndexParams("vec", "embedding", float64(5), "COSINE")},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("mismatched vector override: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := pitrRestore("PitrOvUnknown", map[string]interface{}{
		"VectorIndexOverride": []interface{}{vectorIndexParams("other", "embedding", float64(2), "COSINE")},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("unknown vector override: err = %v, want ErrInvalidParameter", err)
	}
}

// TestPITRRestoreRejectsMistypedUseLatestRestorableTime pins the boolean
// member's strict contract: a present-but-mistyped value is a request
// error, never a silently ignored one that falls through to the
// RestoreDateTime branch.
func TestPITRRestoreRejectsMistypedUseLatestRestorableTime(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "PitrBoolSrc",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create source: %v", err)
	}
	if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                        "PitrBoolSrc",
		"PointInTimeRecoverySpecification": map[string]interface{}{"PointInTimeRecoveryEnabled": true},
	}}); err != nil {
		t.Fatalf("enable pitr: %v", err)
	}

	_, err := svc.RestoreTableToPointInTime(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"SourceTableName":         "PitrBoolSrc",
		"TargetTableName":         "PitrBoolTgt",
		"UseLatestRestorableTime": "latest",
	}})
	if err == nil || !strings.Contains(err.Error(), "UseLatestRestorableTime must be a boolean") {
		t.Fatalf("mistyped UseLatestRestorableTime: expected rejection, got %v", err)
	}
}

// TestCreateBackupDuplicateNameCoexists pins the identity model on the
// create path: the name is a label and the documented CreateBackup error
// set carries no duplicate-name rejection, so a second backup under a name
// already in use — on the same table or another — is created, carrying its
// own generated id and ARN.
func TestCreateBackupDuplicateNameCoexists(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	first := createRestorableBackup(t, svc, reqCtx, "DupSrc", "dup-name-bk")

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "DupOther",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create DupOther: %v", err)
	}

	second, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName":  "DupOther",
		"BackupName": "dup-name-bk",
	}})
	if err != nil {
		t.Fatalf("duplicate backup on DupOther: %v", err)
	}
	if second.BackupArn == first.BackupArn {
		t.Fatalf("same-named backups share an ARN: %s", first.BackupArn)
	}
	sameTable, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName":  "DupSrc",
		"BackupName": "dup-name-bk",
	}})
	if err != nil {
		t.Fatalf("duplicate backup on DupSrc: %v", err)
	}
	if sameTable.BackupArn == first.BackupArn || sameTable.BackupArn == second.BackupArn {
		t.Fatalf("same-table same-name backup shares an ARN: %s", sameTable.BackupArn)
	}
}

// TestRestoreStatusTransitionsPreserveConcurrentCounterWrites pins the
// locked-transition contract of the shared restore pipeline: a counter write
// a concurrent flow lands on the freshly created target record must survive
// the CREATING downgrade and reach the final ACTIVE record. The downgrade
// runs through the locked read-modify-write Update — a wholesale write of
// the Create-time record would drop whatever the racing writer committed in
// the window, so every interleaving must end with both effects present.
func TestRestoreStatusTransitionsPreserveConcurrentCounterWrites(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	backup := createRestorableBackup(t, svc, reqCtx, "RaceSrc", "race-bk")
	if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "RaceSrc",
		"PointInTimeRecoverySpecification": map[string]interface{}{
			"PointInTimeRecoveryEnabled": true,
		},
	}}); err != nil {
		t.Fatalf("enable pitr: %v", err)
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}

	// The single-item restore baseline the racing write adds to.
	control, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn: backup.BackupArn, TargetTableName: "RaceControl",
	})
	if err != nil {
		t.Fatalf("control restore: %v", err)
	}
	if control.ItemCount != 1 {
		t.Fatalf("control restore item count = %d, want 1", control.ItemCount)
	}

	// racer waits for the target table to appear and then lands a racing
	// counter Update on it. The wait is bounded — a restore that never
	// creates its target fails the round through the returned error
	// instead of spinning forever — and the Update's error is carried to
	// the test goroutine rather than swallowed.
	racer := func(target string) chan error {
		done := make(chan error, 1)
		go func() {
			deadline := time.Now().Add(5 * time.Second)
			for !store.Tables().Exists(target) {
				if time.Now().After(deadline) {
					done <- fmt.Errorf("racer: target %s never appeared", target)
					return
				}
			}
			if _, err := store.Tables().Update(target, func(table *dbstore.Table) error {
				table.ItemCount += 1000
				return nil
			}); err != nil {
				done <- fmt.Errorf("racer update %s: %w", target, err)
				return
			}
			done <- nil
		}()
		return done
	}

	for round := 0; round < 8; round++ {
		target := fmt.Sprintf("RaceFrom%d", round)
		done := racer(target)
		restored, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
			BackupArn: backup.BackupArn, TargetTableName: target,
		})
		if rerr := <-done; rerr != nil {
			t.Fatalf("round %d racer: %v", round, rerr)
		}
		if err != nil {
			t.Fatalf("round %d restore: %v", round, err)
		}
		// The committed record is the artefact both effects must be present
		// in: the racing Update serialises against the pipeline's locked
		// stages, and when it lands after the final ACTIVE stamp the
		// pipeline's return snapshot predates it while the store carries
		// both — so the assertion reads the record fresh, as the PITR half
		// below does.
		final, err := store.Tables().Get(target)
		if err != nil {
			t.Fatalf("round %d get target: %v", round, err)
		}
		if final.ItemCount != 1+1000 {
			t.Fatalf("round %d from-backup item count = %d, want %d (racing counter write lost)",
				round, final.ItemCount, 1+1000)
		}
		if restored.Status != dbstore.TableStatusActive {
			t.Fatalf("round %d restored status = %v, want ACTIVE", round, restored.Status)
		}

		target = fmt.Sprintf("RacePitr%d", round)
		done = racer(target)
		_, err = svc.restoreTableToPointInTimeCore(ctx, reqCtx, restoreTableToPointInTimeInput{Parameters: map[string]interface{}{
			"SourceTableName": "RaceSrc", "TargetTableName": target, "UseLatestRestorableTime": true,
		}})
		if rerr := <-done; rerr != nil {
			t.Fatalf("round %d pitr racer: %v", round, rerr)
		}
		if err != nil {
			t.Fatalf("round %d pitr restore: %v", round, err)
		}
		final, err = store.Tables().Get(target)
		if err != nil {
			t.Fatalf("round %d pitr get target: %v", round, err)
		}
		if final.ItemCount != 1+1000 {
			t.Fatalf("round %d pitr item count = %d, want %d (racing counter write lost)",
				round, final.ItemCount, 1+1000)
		}
	}
}

func TestFlushRestoreChunkUsesCanonicalWriteComposition(t *testing.T) {
	svc, _, store := serviceRegionFixture(t)
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "RestoreChunk",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	items := []*dbstore.Item{
		{TableName: "RestoreChunk", Key: map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("r1")},
			Attributes: map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("r1"), "v": dbstore.StringValue("one")}},
		{TableName: "RestoreChunk", Key: map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("r2")},
			Attributes: map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("r2"), "v": dbstore.StringValue("two")}},
	}
	if err := svc.flushRestoreChunk(context.Background(), store, "RestoreChunk", items); err != nil {
		t.Fatalf("flush chunk: %v", err)
	}

	for _, item := range items {
		got, err := store.Items().Get("RestoreChunk", item.Key)
		if err != nil {
			t.Fatalf("get %s: %v", *item.Key["id"].S, err)
		}
		if *got.Attributes["v"].S != *item.Attributes["v"].S {
			t.Fatalf("item %s: expected v=%s, got %v", *item.Key["id"].S, *item.Attributes["v"].S, got.Attributes["v"])
		}
	}
	table, err := store.Tables().Get("RestoreChunk")
	if err != nil {
		t.Fatalf("table get: %v", err)
	}
	var wantSize int64
	for _, item := range items {
		wantSize += dbstore.CalculateItemSize(item.Attributes)
	}
	if table.ItemCount != 2 || table.TableSizeBytes != wantSize {
		t.Fatalf("metrics: expected count 2 size %d, got count %d size %d", wantSize, table.ItemCount, table.TableSizeBytes)
	}
}
