// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// tableIdPattern is the UUID shape the Smithy TableId pattern documents
// (^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$).
var tableIdPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// TestTableIdIdentity pins the minted table identity: DescribeTable emits
// a pattern-valid UUID, the backup surfaces (DescribeBackup's
// SourceTableDetails and ListBackups' BackupSummary) carry the same
// identity as the source table, the export record and its
// manifest-summary tableId field carry it, the import record carries the
// created target table's identity, and a same-name re-created table mints
// a distinct identity.
func TestTableIdIdentity(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	create := func(name string) string {
		t.Helper()
		if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":            name,
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		}}); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
		desc := describeTableOnModePlane(t, svc, reqCtx, name)
		id, ok := desc["TableId"].(string)
		if !ok || !tableIdPattern.MatchString(id) {
			t.Fatalf("TableId = %v, want a pattern-valid UUID", desc["TableId"])
		}
		return id
	}

	tableId := create("IdTbl")
	if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "IdTbl",
		"Item":      map[string]interface{}{"pk": map[string]interface{}{"S": "one"}},
	}}); err != nil {
		t.Fatalf("put item: %v", err)
	}

	// The backup surfaces carry the source table's identity.
	backup, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName": "IdTbl", "BackupName": "id-bk",
	}})
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}
	if backup.SourceTableId != tableId {
		t.Fatalf("backup SourceTableId = %q, want the table's %q", backup.SourceTableId, tableId)
	}
	describe, err := svc.DescribeBackup(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"BackupArn": backup.BackupArn,
	}})
	if err != nil {
		t.Fatalf("describe backup: %v", err)
	}
	source := describe.(map[string]interface{})["BackupDescription"].(map[string]interface{})["SourceTableDetails"].(map[string]interface{})
	if source["TableId"] != tableId {
		t.Fatalf("SourceTableDetails.TableId = %v, want %q", source["TableId"], tableId)
	}
	listed, err := svc.ListBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "IdTbl",
	}})
	if err != nil {
		t.Fatalf("list backups: %v", err)
	}
	summary := listed.(map[string]interface{})["BackupSummaries"].([]map[string]interface{})[0]
	if summary["TableId"] != tableId {
		t.Fatalf("BackupSummary.TableId = %v, want %q", summary["TableId"], tableId)
	}

	// The export record and its manifest carry the identity.
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	table, err := store.Tables().Get("IdTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	recorder := newRecordingS3Invoker()
	svc.SetEventBus(ionExportBus{s3: recorder})
	export, err := store.Exports().Create(table.ARN, table.TableId, "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("create export record: %v", err)
	}
	svc.runExportJob(store, ExportTableCoreInput{
		TableArn:     table.ARN,
		TableName:    "IdTbl",
		ExportFormat: "DYNAMODB_JSON",
		S3Bucket:     "dest-bucket",
		Region:       "us-east-1",
		ExportTime:   time.Now(),
	}, export.ExportArn)
	finalExport, err := store.Exports().Get(export.ExportArn)
	if err != nil {
		t.Fatalf("reload export: %v", err)
	}
	if finalExport.TableId != tableId {
		t.Fatalf("export TableId = %q, want %q", finalExport.TableId, tableId)
	}
	for key, data := range recorder.puts {
		if len(key) > len("/manifest-summary.json") && key[len(key)-len("/manifest-summary.json"):] == "/manifest-summary.json" {
			if !strings.Contains(string(data), `"tableId":"`+tableId+`"`) {
				t.Fatalf("manifest-summary tableId does not carry the identity: %s", data)
			}
		}
	}

	// The import record carries the created target table's identity.
	imp, err := store.Imports().Create(table.ARN, "")
	if err != nil {
		t.Fatalf("create import record: %v", err)
	}
	svc.runImportJob(store, "us-east-1", ImportTableCoreInput{
		TableName:     "IdImported",
		KeySchema:     table.KeySchema,
		AttributeDefs: table.AttributeDefinitions,
		BillingMode:   dbstore.BillingModePayPerRequest,
		S3Bucket:      "src-bucket",
		InputFormat:   "DYNAMODB_JSON",
	}, imp.ImportArn)
	finalImport, err := store.Imports().Get(imp.ImportArn)
	if err != nil {
		t.Fatalf("reload import: %v", err)
	}
	imported, err := store.Tables().Get("IdImported")
	if err != nil {
		t.Fatalf("get imported table: %v", err)
	}
	if imported.TableId == "" || !tableIdPattern.MatchString(imported.TableId) {
		t.Fatalf("imported table TableId = %q, want a minted UUID", imported.TableId)
	}
	if finalImport.TableId != imported.TableId {
		t.Fatalf("import TableId = %q, want the imported table's %q", finalImport.TableId, imported.TableId)
	}

	// A same-name re-creation mints a fresh identity.
	if _, err := svc.DeleteTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "IdTbl",
	}}); err != nil {
		t.Fatalf("delete IdTbl: %v", err)
	}
	if regenerated := create("IdTbl"); regenerated == tableId {
		t.Fatal("re-created table reused the dead table's identity")
	}
}
