// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TestBackupArnCarriesGeneratedId pins the documented shape of the backup
// ARN: the identity segment after /backup/ is a generated id — the creation
// time at millisecond granularity plus an 8-hex-digit suffix, the shape the
// documented example shows (table/Music/backup/01489602797149-73d8d5bc) —
// never the backup's name.
func TestBackupArnCarriesGeneratedId(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	idRe := regexp.MustCompile(`^arn:aws:dynamodb:[^:]*:[^:]*:table/GenIdTbl/backup/[0-9]{14}-[0-9a-f]{8}$`)
	backup := createRestorableBackup(t, svc, reqCtx, "GenIdTbl", "gen-id-bk")
	if !idRe.MatchString(backup.BackupArn) {
		t.Fatalf("backup ARN %q does not carry the generated-id segment", backup.BackupArn)
	}
}

// TestDuplicateBackupNamesCoexist pins the identity model: a backup's
// identity is its ARN's generated id and the documented CreateBackup error
// set carries no duplicate-name rejection, so a second backup under a name
// already in use is created — a distinct ARN, both listed.
func TestDuplicateBackupNamesCoexist(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	first := createRestorableBackup(t, svc, reqCtx, "DupNameTbl", "dup-bk")
	second := createRestorableBackup(t, svc, reqCtx, "DupNameTbl2", "dup-bk")
	if first.BackupArn == second.BackupArn {
		t.Fatalf("same-named backups share an ARN: %s", first.BackupArn)
	}
	if n := countName(listBackupNames(t, svc, reqCtx, ""), "dup-bk"); n != 2 {
		t.Fatalf("duplicate-name backups listed %d times, want both", n)
	}
}

// listBackupNames returns the backup names ListBackups reports for a table.
func listBackupNames(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, table string) []string {
	t.Helper()
	resp, err := svc.ListBackups(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": table,
	}})
	if err != nil {
		t.Fatalf("list backups for %s: %v", table, err)
	}
	summaries := resp.(map[string]interface{})["BackupSummaries"].([]map[string]interface{})
	names := make([]string, 0, len(summaries))
	for _, s := range summaries {
		names = append(names, s["BackupName"].(string))
	}
	return names
}

// TestListBackupsTimeRangeFilterBounds pins the time-range filter: the
// model members TimeRangeLowerBound/TimeRangeUpperBound (epoch bounds)
// filter the listing by backup creation time — presence decides, so an
// omitted bound filters nothing.
func TestListBackupsTimeRangeFilterBounds(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	older := createRestorableBackup(t, svc, reqCtx, "RangeTbl", "range-older")
	newer := createRestorableBackup(t, svc, reqCtx, "RangeTbl2", "range-newer")
	if older.BackupCreationDateTime.After(newer.BackupCreationDateTime) {
		older, newer = newer, older
	}

	list := func(params map[string]interface{}) map[string]bool {
		t.Helper()
		resp, err := svc.ListBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: params})
		if err != nil {
			t.Fatalf("list backups: %v", err)
		}
		names := make(map[string]bool)
		for _, s := range resp.(map[string]interface{})["BackupSummaries"].([]map[string]interface{}) {
			names[s["BackupName"].(string)] = true
		}
		return names
	}
	allNames := func(got map[string]bool) []string {
		var names []string
		for n := range got {
			names = append(names, n)
		}
		return names
	}

	// No bounds: both tables' backups list (the filter is table-scoped, so
	// each query names its own table).
	if got := list(map[string]interface{}{"TableName": "RangeTbl"}); !got["range-older"] {
		t.Fatalf("unbounded list: expected range-older, got %v", allNames(got))
	}

	// A lower bound after the older backup's creation excludes it.
	mid := older.BackupCreationDateTime.Unix() + 1
	if got := list(map[string]interface{}{"TableName": "RangeTbl", "TimeRangeLowerBound": float64(mid)}); got["range-older"] {
		t.Fatalf("lower-bound filter: range-older must be excluded, got %v", allNames(got))
	}

	// An upper bound before the newer backup's creation excludes it.
	pre := newer.BackupCreationDateTime.Unix() - 1
	if got := list(map[string]interface{}{"TableName": "RangeTbl2", "TimeRangeUpperBound": float64(pre)}); got["range-newer"] {
		t.Fatalf("upper-bound filter: range-newer must be excluded, got %v", allNames(got))
	}

	// Bounds spanning both keep each table's backup.
	if got := list(map[string]interface{}{
		"TableName":           "RangeTbl",
		"TimeRangeLowerBound": float64(older.BackupCreationDateTime.Unix() - 1),
		"TimeRangeUpperBound": float64(newer.BackupCreationDateTime.Unix() + 1),
	}); !got["range-older"] {
		t.Fatalf("spanning bounds: expected range-older, got %v", allNames(got))
	}
}

// TestListBackupsAcceptsTableARN pins the documented ARN form of the
// TableName member ("the name or Amazon Resource Name (ARN) of the
// table"): the listing under a table ARN equals the bare-name listing.
func TestListBackupsAcceptsTableARN(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	createRestorableBackup(t, svc, reqCtx, "ArnTbl", "arn-form-bk")

	names := func(tableName string) []string {
		t.Helper()
		resp, err := svc.ListBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": tableName,
		}})
		if err != nil {
			t.Fatalf("list backups for %s: %v", tableName, err)
		}
		var got []string
		for _, s := range resp.(map[string]interface{})["BackupSummaries"].([]map[string]interface{}) {
			got = append(got, s["BackupName"].(string))
		}
		return got
	}

	bare := names("ArnTbl")
	byArn := names("arn:aws:dynamodb:us-east-1:123456789012:table/ArnTbl")
	if len(bare) != 1 || bare[0] != "arn-form-bk" {
		t.Fatalf("bare-name listing: %v", bare)
	}
	if len(byArn) != len(bare) || byArn[0] != bare[0] {
		t.Fatalf("ARN-form listing: %v, want %v", byArn, bare)
	}
}

// TestRestoreDistinguishesAbsenceFromStoreFailure pins the restore
// path's error split: a missing backup answers the documented
// BackupNotFound, while a storage fault on an existing record (a corrupt
// payload the proto reader cannot unmarshal) surfaces as the storage
// error it is — never masquerading as absence.
func TestRestoreDistinguishesAbsenceFromStoreFailure(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	backup := createRestorableBackup(t, svc, reqCtx, "RestoreSplitTbl", "split-bk")

	restore := func(backupArn string) error {
		_, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
			BackupArn:       backupArn,
			TargetTableName: "RestoreSplitTarget",
		})
		return err
	}

	if err := restore("arn:aws:dynamodb:us-east-1:123456789012:table/SplitTbl/backup/01234567890123-deadbef0"); !errors.Is(err, ErrBackupNotFound) {
		t.Fatalf("restore of a missing backup: expected ErrBackupNotFound, got %v", err)
	}

	regionStorage, err := reqCtx.GetStorage()
	if err != nil {
		t.Fatalf("region storage: %v", err)
	}
	backupKey := svcarn.ExtractBackupIdFromARN(backup.BackupArn)
	if err := regionStorage.Bucket("dynamodb_backups-"+reqCtx.GetRegion()).Put([]byte(backupKey), []byte("not-a-protobuf-record")); err != nil {
		t.Fatalf("corrupt backup record: %v", err)
	}

	if err := restore(backup.BackupArn); err == nil || errors.Is(err, ErrBackupNotFound) {
		t.Fatalf("restore over a corrupt record: expected the storage failure, got %v", err)
	}
}

// containsName reports whether the list holds the exact name.
func containsName(names []string, want string) bool {
	for _, n := range names {
		if n == want {
			return true
		}
	}
	return false
}

// TestDeleteTableKeepsBackupsAndWritesSystemBackup pins the deletion
// contract's recovery half: deleting a recovery-enabled table writes the
// table-name$DeletedTableBackup system backup (a snapshot of the table
// right before deletion, restorable), deleting a table without recovery
// enabled writes none, on-demand backups of every table survive any
// deletion (a backup is an independently ARN-addressed resource — deleting
// one table may not touch another table's backups), and point-in-time
// restore by name keeps requiring a live source table — the documented
// recovery path for a deleted table is the system backup.
func TestDeleteTableKeepsBackupsAndWritesSystemBackup(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	enablePitr := func(table string) {
		t.Helper()
		if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": table,
			"PointInTimeRecoverySpecification": map[string]interface{}{
				"PointInTimeRecoveryEnabled": true,
			},
		}}); err != nil {
			t.Fatalf("enable pitr on %s: %v", table, err)
		}
	}
	deleteTable := func(table string) {
		t.Helper()
		if _, err := svc.DeleteTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": table,
		}}); err != nil {
			t.Fatalf("delete %s: %v", table, err)
		}
	}

	// The recovery-enabled table, its on-demand backup, and an unrelated
	// table's backup that must survive the deletion.
	pitrBackup := createRestorableBackup(t, svc, reqCtx, "DelPitrTbl", "user-bk")
	createRestorableBackup(t, svc, reqCtx, "OtherTbl", "other-bk")
	enablePitr("DelPitrTbl")

	// A table without recovery enabled, for the negative half.
	createRestorableBackup(t, svc, reqCtx, "DelPlainTbl", "plain-bk")

	deleteTable("DelPitrTbl")

	systemName := "DelPitrTbl" + dbstore.DeletedTableBackupSuffix
	names := listBackupNames(t, svc, reqCtx, "DelPitrTbl")
	if !containsName(names, "user-bk") || !containsName(names, systemName) {
		t.Fatalf("backups after recovery-enabled delete = %v, want user-bk and %s", names, systemName)
	}

	// The system backup restores the deleted table's data.
	systemBackup, err := svc.describeBackupCore(ctx, reqCtx, systemArnFor(t, svc, reqCtx, "DelPitrTbl", systemName))
	if err != nil {
		t.Fatalf("describe system backup: %v", err)
	}
	restored, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn:       systemBackup.BackupArn,
		TargetTableName: "SysRestored",
	})
	if err != nil {
		t.Fatalf("restore from system backup: %v", err)
	}
	if restored.Status != dbstore.TableStatusActive || restored.RestoreSummary == nil ||
		restored.RestoreSummary.SourceBackupArn != systemBackup.BackupArn {
		t.Fatalf("system-restored table = %+v", restored)
	}
	item, err := svc.GetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "SysRestored",
		"Key":       map[string]interface{}{"pk": map[string]interface{}{"S": "kept"}},
	}})
	if err != nil {
		t.Fatalf("get system-restored item: %v", err)
	}
	if resp := item.(map[string]interface{}); resp["Item"] == nil {
		t.Fatalf("system-restored item missing: %+v", resp)
	}

	// Point-in-time restore by name keeps requiring a live source table:
	// the documented recovery path for a deleted table is the system
	// backup alone.
	if _, err := svc.RestoreTableToPointInTime(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"SourceTableName":         "DelPitrTbl",
		"TargetTableName":         "PitrOfDeleted",
		"UseLatestRestorableTime": true,
	}}); !errors.Is(err, ErrTableNotFoundException) {
		t.Fatalf("pitr restore of deleted table: err = %v, want ErrTableNotFoundException", err)
	}

	// The on-demand backup of the deleted table still restores.
	if _, err := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn:       pitrBackup.BackupArn,
		TargetTableName: "UserRestored",
	}); err != nil {
		t.Fatalf("restore from surviving user backup: %v", err)
	}

	// An unrelated table's backup survives the deletion untouched.
	if names := listBackupNames(t, svc, reqCtx, "OtherTbl"); !containsName(names, "other-bk") {
		t.Fatalf("unrelated table's backups after delete = %v, want other-bk", names)
	}

	// A table without recovery enabled leaves no system backup.
	deleteTable("DelPlainTbl")
	if names := listBackupNames(t, svc, reqCtx, "DelPlainTbl"); !containsName(names, "plain-bk") || containsName(names, "DelPlainTbl"+dbstore.DeletedTableBackupSuffix) {
		t.Fatalf("backups after plain delete = %v, want plain-bk alone", names)
	}
}

// systemArnFor returns the backup ARN ListBackups reports under the given
// name for a table.
func systemArnFor(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, table, backupName string) string {
	t.Helper()
	resp, err := svc.ListBackups(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": table,
	}})
	if err != nil {
		t.Fatalf("list backups for arn lookup: %v", err)
	}
	for _, s := range resp.(map[string]interface{})["BackupSummaries"].([]map[string]interface{}) {
		if s["BackupName"] == backupName {
			return s["BackupArn"].(string)
		}
	}
	t.Fatalf("backup %s not listed for %s", backupName, table)
	return ""
}

// TestSystemBackupSweeperDeletesExpired pins the retention sweep: a system
// backup older than the documented 35-day window is deleted with its
// snapshot, a fresh one is kept, and an aged on-demand backup is outside
// the sweep entirely.
func TestSystemBackupSweeperDeletesExpired(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	createRestorableBackup(t, svc, reqCtx, "AgedTbl", "aged-user-bk")
	if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "AgedTbl",
		"PointInTimeRecoverySpecification": map[string]interface{}{
			"PointInTimeRecoveryEnabled": true,
		},
	}}); err != nil {
		t.Fatalf("enable pitr on AgedTbl: %v", err)
	}
	if _, err := svc.DeleteTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "AgedTbl",
	}}); err != nil {
		t.Fatalf("delete AgedTbl: %v", err)
	}
	createRestorableBackup(t, svc, reqCtx, "FreshTbl", "fresh-user-bk")
	if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "FreshTbl",
		"PointInTimeRecoverySpecification": map[string]interface{}{
			"PointInTimeRecoveryEnabled": true,
		},
	}}); err != nil {
		t.Fatalf("enable pitr on FreshTbl: %v", err)
	}
	if _, err := svc.DeleteTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "FreshTbl",
	}}); err != nil {
		t.Fatalf("delete FreshTbl: %v", err)
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	ageBackup := func(name string) *dbstore.Backup {
		t.Helper()
		b, err := findBackupByName(t, store, name)
		if err != nil {
			t.Fatalf("get backup %s: %v", name, err)
		}
		b.BackupCreationDateTime = time.Now().UTC().Add(-deletedTableBackupRetentionDays*24*time.Hour - time.Hour)
		if err := store.Backups().Put(b); err != nil {
			t.Fatalf("age backup %s: %v", name, err)
		}
		return b
	}
	agedSystem := ageBackup("AgedTbl" + dbstore.DeletedTableBackupSuffix)
	agedUser := ageBackup("aged-user-bk")

	svc.sweepStoreSystemBackups(store)

	if _, err := store.Backups().Get(agedSystem.BackupArn); err == nil {
		t.Fatal("expired system backup survived the sweep")
	}
	if items, err := store.Backups().GetSnapshot(agedSystem.BackupArn); err == nil && len(items) > 0 {
		t.Fatalf("expired system backup snapshot survived: %d items", len(items))
	}
	if _, err := findBackupByName(t, store, "FreshTbl"+dbstore.DeletedTableBackupSuffix); err != nil {
		t.Fatalf("fresh system backup swept: %v", err)
	}
	if _, err := store.Backups().Get(agedUser.BackupArn); err != nil {
		t.Fatalf("aged on-demand backup swept: %v", err)
	}
}

// findBackupsByName collects every backup record carrying the exact name —
// the name is a label carried by any number of records, so lookups walk
// the list.
func findBackupsByName(t *testing.T, store dbstore.DynamoDBStoreInterface, name string) []*dbstore.Backup {
	t.Helper()
	var matches []*dbstore.Backup
	marker := ""
	for {
		backups, next, err := store.Backups().List(marker, listBackupsMinFetchSize, "")
		if err != nil {
			t.Fatalf("list backups looking for %s: %v", name, err)
		}
		for _, b := range backups {
			if b.BackupName == name {
				matches = append(matches, b)
			}
		}
		if next == "" {
			return matches
		}
		marker = next
	}
}

func findBackupByName(t *testing.T, store dbstore.DynamoDBStoreInterface, name string) (*dbstore.Backup, error) {
	matches := findBackupsByName(t, store, name)
	if len(matches) == 0 {
		return nil, dbstore.ErrBackupNotFound
	}
	return matches[0], nil
}

// countName counts the exact-name occurrences in a backup-name list.
func countName(names []string, want string) int {
	n := 0
	for _, name := range names {
		if name == want {
			n++
		}
	}
	return n
}

// TestSystemBackupRedeleteRetainsBothGenerations pins the identity model of
// a repeated delete cycle: the system backup's NAME follows the documented
// table-name$DeletedTableBackup convention, but its identity is a
// generated id in its ARN, so recreating the same table name and deleting
// it again inside the 35-day window leaves BOTH generations — each listed
// under the same name, each addressable by its own ARN, each restoring its
// own generation's items. The documented retention ("retains it for 35
// days") is exactly what an overwrite would break: the first delete's
// recovery point must survive the second delete.
func TestSystemBackupRedeleteRetainsBothGenerations(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	createTable := func(table string) {
		t.Helper()
		if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":            table,
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		}}); err != nil {
			t.Fatalf("create %s: %v", table, err)
		}
	}
	enablePitr := func(table string) {
		t.Helper()
		if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": table,
			"PointInTimeRecoverySpecification": map[string]interface{}{
				"PointInTimeRecoveryEnabled": true,
			},
		}}); err != nil {
			t.Fatalf("enable pitr on %s: %v", table, err)
		}
	}
	deleteTable := func(table string) {
		t.Helper()
		if _, err := svc.DeleteTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": table,
		}}); err != nil {
			t.Fatalf("delete %s: %v", table, err)
		}
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	systemName := "ReTbl" + dbstore.DeletedTableBackupSuffix

	createTable("ReTbl")
	enablePitr("ReTbl")
	deleteTable("ReTbl")

	// The recreated table is a new instance carrying a distinguishing item.
	createTable("ReTbl")
	if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "ReTbl",
		"Item":      map[string]interface{}{"pk": map[string]interface{}{"S": "second-generation"}},
	}}); err != nil {
		t.Fatalf("put second-generation item: %v", err)
	}
	enablePitr("ReTbl")
	deleteTable("ReTbl")

	generations := findBackupsByName(t, store, systemName)
	if len(generations) != 2 {
		t.Fatalf("system backups under %s: %d records, want both generations", systemName, len(generations))
	}
	if generations[0].BackupArn == generations[1].BackupArn {
		t.Fatalf("both generations share an ARN: %s", generations[0].BackupArn)
	}
	if generations[0].SourceTableId == generations[1].SourceTableId {
		t.Fatalf("both generations carry table id %s — the records are not distinct", generations[0].SourceTableId)
	}
	if n := countName(listBackupNames(t, svc, reqCtx, "ReTbl"), systemName); n != 2 {
		t.Fatalf("system backup listed %d times, want both generations", n)
	}

	// Each generation restores its own items: exactly one of the two ARNs —
	// the second generation's — yields the recreated table's distinguishing
	// item, the other restores the empty first incarnation.
	withItem := 0
	for i, gen := range generations {
		target := fmt.Sprintf("ReTblRestored%d", i)
		restored, rErr := svc.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
			BackupArn:       gen.BackupArn,
			TargetTableName: target,
		})
		if rErr != nil {
			t.Fatalf("restore from %s: %v", gen.BackupArn, rErr)
		}
		if restored.Status != dbstore.TableStatusActive {
			t.Fatalf("restored table status = %s", restored.Status)
		}
		item, gErr := svc.GetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": target,
			"Key":       map[string]interface{}{"pk": map[string]interface{}{"S": "second-generation"}},
		}})
		if gErr != nil {
			t.Fatalf("get restored item from %s: %v", gen.BackupArn, gErr)
		}
		if item.(map[string]interface{})["Item"] != nil {
			withItem++
		}
	}
	if withItem != 1 {
		t.Fatalf("exactly the second generation's restore must carry the distinguishing item, got %d of 2", withItem)
	}
}
