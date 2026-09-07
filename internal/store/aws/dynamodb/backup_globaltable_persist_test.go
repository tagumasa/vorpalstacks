package dynamodb

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// newBackupGlobalStore opens a store with one hash-only table for
// backup and global-table persistence tests.
func newBackupGlobalStore(t *testing.T) *DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewDynamoDBStore(st, "123456789012", "us-east-1")
}

// Backups persist as protobuf records and their item snapshots round-trip
// with key and attribute maps intact.
func TestBackupAndSnapshotProtoRoundTrip(t *testing.T) {
	store := newBackupGlobalStore(t)
	backups := store.Backups()

	backup, err := backups.Create("snap-backup", "Tbl", "arn:aws:dynamodb:us-east-1:123456789012:table/Tbl", 128)
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}
	if backup.BackupStatus != BackupStatusAvailable || backup.BackupType != BackupTypeUser {
		t.Fatalf("created backup = %+v", backup)
	}

	loaded, err := backups.GetByName("snap-backup")
	if err != nil {
		t.Fatalf("get backup: %v", err)
	}
	if loaded.BackupArn != backup.BackupArn || loaded.BackupSizeBytes != 128 ||
		!loaded.BackupCreationDateTime.Equal(backup.BackupCreationDateTime) {
		t.Fatalf("round-trip mismatch: %+v", loaded)
	}

	byArn, err := backups.Get(backup.BackupArn)
	if err != nil || byArn.BackupName != "snap-backup" {
		t.Fatalf("get by arn: %v, %+v", err, byArn)
	}

	strVal := "v1"
	items := []*Item{{
		TableName: "Tbl",
		Key:       map[string]*AttributeValue{"id": {S: &strVal}},
		Attributes: map[string]*AttributeValue{
			"id":  {S: &strVal},
			"num": {N: &strVal},
		},
	}}
	if err := backups.SaveSnapshot("snap-backup", items); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	loadedItems, err := backups.GetSnapshot("snap-backup")
	if err != nil {
		t.Fatalf("get snapshot: %v", err)
	}
	if len(loadedItems) != 1 || loadedItems[0].TableName != "Tbl" ||
		loadedItems[0].Key["id"] == nil || loadedItems[0].Key["id"].S == nil ||
		*loadedItems[0].Key["id"].S != "v1" || loadedItems[0].Attributes["num"] == nil {
		t.Fatalf("snapshot round-trip mismatch: %+v", loadedItems)
	}

	// The snapshot key must stay invisible to ListBackups.
	listed, _, err := backups.List("", 100, "")
	if err != nil {
		t.Fatalf("list backups: %v", err)
	}
	if len(listed) != 1 || listed[0].BackupName != "snap-backup" {
		t.Fatalf("snapshot key leaked into list: %+v", listed)
	}
}

// Global tables persist as protobuf records with the typed auto-scaling
// settings, per-index settings and table-class members intact.
func TestGlobalTableProtoRoundTrip(t *testing.T) {
	store := newBackupGlobalStore(t)
	globalTables := store.GlobalTables()

	units := int64(5)
	updatedAt := time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC)
	created, err := globalTables.Create("gt", []*Replica{{
		RegionName:                    "us-east-1",
		ReplicaStatus:                 "ACTIVE",
		ProvisionedReadCapacityUnits:  10,
		ProvisionedWriteCapacityUnits: 20,
		ReadAutoScalingSettings:       &AutoScalingSettingsDescription{MinimumUnits: &units},
		GlobalSecondaryIndexReadSettings: []IndexAutoScalingSettings{{
			IndexName:                    "gsi-1",
			ProvisionedReadCapacityUnits: &units,
			Read:                         &AutoScalingSettingsDescription{MinimumUnits: &units},
		}},
		TableClass:            "STANDARD",
		TableClassLastUpdated: &updatedAt,
	}})
	if err != nil {
		t.Fatalf("create global table: %v", err)
	}

	if _, err := globalTables.Update("gt", func(gt *GlobalTable) error {
		gt.WriteAutoScalingSettings = &AutoScalingSettingsDescription{MinimumUnits: &units}
		gt.GlobalSecondaryIndexWriteSettings = []IndexAutoScalingSettings{{
			IndexName:                     "gsi-1",
			ProvisionedWriteCapacityUnits: &units,
			Write:                         &AutoScalingSettingsDescription{MinimumUnits: &units},
		}}
		return nil
	}); err != nil {
		t.Fatalf("update global table: %v", err)
	}

	loaded, err := globalTables.Get("gt")
	if err != nil {
		t.Fatalf("get global table: %v", err)
	}
	if loaded.GlobalTableArn != created.GlobalTableArn || !loaded.CreationDateTime.Equal(created.CreationDateTime) {
		t.Fatalf("round-trip mismatch: %+v", loaded)
	}
	if loaded.WriteAutoScalingSettings == nil || loaded.WriteAutoScalingSettings.MinimumUnits == nil ||
		*loaded.WriteAutoScalingSettings.MinimumUnits != 5 {
		t.Fatalf("global write settings lost: %+v", loaded.WriteAutoScalingSettings)
	}
	if len(loaded.GlobalSecondaryIndexWriteSettings) != 1 ||
		loaded.GlobalSecondaryIndexWriteSettings[0].IndexName != "gsi-1" ||
		loaded.GlobalSecondaryIndexWriteSettings[0].Write == nil {
		t.Fatalf("global GSI write settings lost: %+v", loaded.GlobalSecondaryIndexWriteSettings)
	}
	replica := loaded.ReplicationGroup[0]
	if replica.ReadAutoScalingSettings == nil || replica.ReadAutoScalingSettings.MinimumUnits == nil ||
		*replica.ReadAutoScalingSettings.MinimumUnits != 5 {
		t.Fatalf("replica read settings lost: %+v", replica.ReadAutoScalingSettings)
	}
	if len(replica.GlobalSecondaryIndexReadSettings) != 1 ||
		replica.GlobalSecondaryIndexReadSettings[0].Read == nil {
		t.Fatalf("replica GSI read settings lost: %+v", replica.GlobalSecondaryIndexReadSettings)
	}
	if replica.TableClass != "STANDARD" || replica.TableClassLastUpdated == nil ||
		!replica.TableClassLastUpdated.Equal(updatedAt) {
		t.Fatalf("replica table class lost: %q, %v", replica.TableClass, replica.TableClassLastUpdated)
	}

	listed, _, err := globalTables.List("", 100)
	if err != nil || len(listed) != 1 || listed[0].GlobalTableName != "gt" {
		t.Fatalf("list global tables: %v, %+v", err, listed)
	}
}
