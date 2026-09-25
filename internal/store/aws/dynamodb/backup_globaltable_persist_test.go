package dynamodb

import (
	"fmt"
	"sync"
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
	return NewDynamoDBStore(st, st, "123456789012", "us-east-1")
}

// Backups persist as protobuf records and their item snapshots round-trip
// with key and attribute maps intact. A user backup's record is born
// CREATING — its snapshot is taken after the record exists, so the
// persisted state never presents a snapshot-less backup as restorable; the
// creating operation writes AVAILABLE once the snapshot lands.
func TestBackupAndSnapshotProtoRoundTrip(t *testing.T) {
	store := newBackupGlobalStore(t)
	backups := store.Backups()

	backup, err := backups.Create("snap-backup", "Tbl", "arn:aws:dynamodb:us-east-1:123456789012:table/Tbl", 128)
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}
	if backup.BackupStatus != BackupStatusCreating || backup.BackupType != BackupTypeUser {
		t.Fatalf("created backup = %+v", backup)
	}

	loaded, err := backups.Get(backup.BackupArn)
	if err != nil {
		t.Fatalf("get backup: %v", err)
	}
	if loaded.BackupName != "snap-backup" || loaded.BackupArn != backup.BackupArn || loaded.BackupSizeBytes != 128 ||
		loaded.BackupStatus != BackupStatusCreating ||
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

// Concurrent Create calls with the same BackupName all succeed: a backup's
// identity is the id its creation mints — the ARN's segment and the record
// key — and the documented CreateBackup error set carries no duplicate-name
// rejection, so same-named creations coexist. The pin keeps the minted ids
// honest under concurrency: every racer gets its own record, addressed by
// its own ARN, whatever the interleaving.
func TestConcurrentSameNameBackupCreatesAllPersist(t *testing.T) {
	store := newBackupGlobalStore(t)
	backups := store.Backups()

	for round := 0; round < 25; round++ {
		name := fmt.Sprintf("dup-backup-%d", round)
		const racers = 4
		start := make(chan struct{})
		type outcome struct {
			backup *Backup
			err    error
		}
		results := make(chan outcome, racers)
		for i := 0; i < racers; i++ {
			go func() {
				<-start
				backup, err := backups.Create(name, "Tbl", "arn:aws:dynamodb:us-east-1:123456789012:table/Tbl", 128)
				results <- outcome{backup: backup, err: err}
			}()
		}
		close(start)
		arns := map[string]bool{}
		for i := 0; i < racers; i++ {
			got := <-results
			if got.err != nil {
				t.Fatalf("round %d: create error: %v", round, got.err)
			}
			if arns[got.backup.BackupArn] {
				t.Fatalf("round %d: two creations share an ARN: %s", round, got.backup.BackupArn)
			}
			arns[got.backup.BackupArn] = true
			byArn, err := backups.Get(got.backup.BackupArn)
			if err != nil || byArn.BackupName != name {
				t.Fatalf("round %d: read-back by ARN: %v, %+v", round, err, byArn)
			}
		}
	}
}

// TestResourceListPagination pins the continuation-marker contract every
// resource family's List shares: a limit below the record count yields a
// page plus the next marker, and following the marker serves the rest and
// ends with an empty marker.
func TestResourceListPagination(t *testing.T) {
	store := newBackupGlobalStore(t)

	mustCreateTable := func(name string) {
		t.Helper()
		if _, err := store.Tables().Create(CreateTableParams{
			Name:                 name,
			KeySchema:            []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
			AttributeDefinitions: []*AttributeDefinition{{AttributeName: "id", AttributeType: ScalarAttributeTypeS}},
			BillingMode:          BillingModePayPerRequest,
		}); err != nil {
			t.Fatalf("create table %s: %v", name, err)
		}
	}
	for _, name := range []string{"pag-a", "pag-b", "pag-c"} {
		mustCreateTable(name)
	}
	tables, next, err := store.Tables().List("", 2)
	if err != nil || len(tables) != 2 || next == "" {
		t.Fatalf("table list page one: %d items, next %q, err %v", len(tables), next, err)
	}
	tables, next, err = store.Tables().List(next, 2)
	if err != nil || len(tables) != 1 || next != "" {
		t.Fatalf("table list page two: %d items, next %q, err %v", len(tables), next, err)
	}

	for _, name := range []string{"pag-a", "pag-b", "pag-c"} {
		if _, err := store.GlobalTables().Create(name, nil); err != nil {
			t.Fatalf("create global table %s: %v", name, err)
		}
	}
	globals, next, err := store.GlobalTables().List("", 2)
	if err != nil || len(globals) != 2 || next == "" {
		t.Fatalf("global table list page one: %d items, next %q, err %v", len(globals), next, err)
	}
	globals, next, err = store.GlobalTables().List(next, 2)
	if err != nil || len(globals) != 1 || next != "" {
		t.Fatalf("global table list page two: %d items, next %q, err %v", len(globals), next, err)
	}

	tableArn := "arn:aws:dynamodb:us-east-1:123456789012:table/PagTbl"
	for _, name := range []string{"pag-a", "pag-b", "pag-c"} {
		if _, err := store.Backups().Create(name, "PagTbl", tableArn, 8); err != nil {
			t.Fatalf("create backup %s: %v", name, err)
		}
	}
	// Exports and imports are seeded through Put with distinct ARNs:
	// Create keys the record by an ARN carrying a second-granularity
	// timestamp, so three creations inside one second would collapse into
	// one key before the pagination under test is reached.
	for _, id := range []string{"pag-1", "pag-2", "pag-3"} {
		if err := store.Exports().Put(&ExportDescription{
			ExportArn: "arn:aws:dynamodb:us-east-1:123456789012:export/" + id,
			TableArn:  tableArn,
		}); err != nil {
			t.Fatalf("put export %s: %v", id, err)
		}
		if err := store.Imports().Put(&ImportTableDescription{
			ImportArn: "arn:aws:dynamodb:us-east-1:123456789012:import/" + id,
			TableArn:  tableArn,
		}); err != nil {
			t.Fatalf("put import %s: %v", id, err)
		}
	}

	backups, next, err := store.Backups().List("", 2, "")
	if err != nil || len(backups) != 2 || next == "" {
		t.Fatalf("backup list page one: %d items, next %q, err %v", len(backups), next, err)
	}
	backups, next, err = store.Backups().List(next, 2, "")
	if err != nil || len(backups) != 1 || next != "" {
		t.Fatalf("backup list page two: %d items, next %q, err %v", len(backups), next, err)
	}

	exports, exportNext, err := store.Exports().List(tableArn, "", 2)
	if err != nil || len(exports) != 2 || exportNext == "" {
		t.Fatalf("export list page one: %d items, next %q, err %v", len(exports), exportNext, err)
	}
	exports, exportNext, err = store.Exports().List(tableArn, exportNext, 2)
	if err != nil || len(exports) != 1 || exportNext != "" {
		t.Fatalf("export list page two: %d items, next %q, err %v", len(exports), exportNext, err)
	}

	imports, importNext, err := store.Imports().List(tableArn, "", 2)
	if err != nil || len(imports) != 2 || importNext == "" {
		t.Fatalf("import list page one: %d items, next %q, err %v", len(imports), importNext, err)
	}
	imports, importNext, err = store.Imports().List(tableArn, importNext, 2)
	if err != nil || len(imports) != 1 || importNext != "" {
		t.Fatalf("import list page two: %d items, next %q, err %v", len(imports), importNext, err)
	}
}

// TestGlobalTableDeleteIfEmpty pins the atomic conditional delete: a
// populated record is kept, an empty one is taken, an absent one is a
// no-op — and a member re-add that commits against a racing
// DeleteIfEmpty is never erased, because the emptiness read and the
// delete hold one lock across both steps.
func TestGlobalTableDeleteIfEmpty(t *testing.T) {
	store := newBackupGlobalStore(t)
	gts := store.GlobalTables()

	if _, err := gts.Create("GtIfEmpty", []*Replica{{RegionName: "us-east-1"}}); err != nil {
		t.Fatalf("create populated record: %v", err)
	}
	deleted, err := gts.DeleteIfEmpty("GtIfEmpty")
	if err != nil || deleted {
		t.Fatalf("populated record: DeleteIfEmpty = (%v, %v), want (false, nil)", deleted, err)
	}
	if !gts.Exists("GtIfEmpty") {
		t.Fatal("populated record was deleted")
	}

	if _, err := gts.Update("GtIfEmpty", func(gt *GlobalTable) error {
		gt.ReplicationGroup = nil
		return nil
	}); err != nil {
		t.Fatalf("empty the record: %v", err)
	}
	deleted, err = gts.DeleteIfEmpty("GtIfEmpty")
	if err != nil || !deleted {
		t.Fatalf("empty record: DeleteIfEmpty = (%v, %v), want (true, nil)", deleted, err)
	}
	if gts.Exists("GtIfEmpty") {
		t.Fatal("empty record survived DeleteIfEmpty")
	}

	deleted, err = gts.DeleteIfEmpty("GtIfEmpty")
	if err != nil || deleted {
		t.Fatalf("absent record: DeleteIfEmpty = (%v, %v), want (false, nil)", deleted, err)
	}
}

// TestGlobalTableDeleteIfEmptySurvivesConcurrentReAdd races the
// conditional delete against a member re-add: whenever the re-add
// commits, the record and its member must survive the race — the
// check-then-delete window the single lock hold closes.
func TestGlobalTableDeleteIfEmptySurvivesConcurrentReAdd(t *testing.T) {
	store := newBackupGlobalStore(t)
	gts := store.GlobalTables()

	for round := 0; round < 50; round++ {
		if _, err := gts.Create("GtRace", nil); err != nil {
			t.Fatalf("round %d: seed empty record: %v", round, err)
		}
		start := make(chan struct{})
		var wg sync.WaitGroup
		added := false
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			if _, err := gts.Update("GtRace", func(gt *GlobalTable) error {
				gt.ReplicationGroup = []*Replica{{RegionName: "us-east-1"}}
				return nil
			}); err == nil {
				added = true
			}
		}()
		go func() {
			defer wg.Done()
			<-start
			_, _ = gts.DeleteIfEmpty("GtRace")
		}()
		close(start)
		wg.Wait()

		if added {
			gt, err := gts.Get("GtRace")
			if err != nil {
				t.Fatalf("round %d: committed re-add erased the record: %v", round, err)
			}
			if len(gt.ReplicationGroup) != 1 {
				t.Fatalf("round %d: committed re-add lost its member: %+v", round, gt.ReplicationGroup)
			}
		}
		if gts.Exists("GtRace") {
			if err := gts.Delete("GtRace"); err != nil {
				t.Fatalf("round %d: cleanup: %v", round, err)
			}
		}
	}
}
