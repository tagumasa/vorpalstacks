package dynamodb

import (
	"regexp"
	"strings"
	"testing"
)

// exportImportIdPattern matches the documented export-ARN ID shape:
// a millisecond-granularity timestamp followed by a random 8-hex-digit
// suffix (the documentation's examples: 01234567890123-a1b2c3d4,
// 01695097218000-d6299cbd).
var exportImportIdPattern = regexp.MustCompile(`^\d{14}-[0-9a-f]{8}$`)

// Two exports of the same table created back-to-back inside one
// wall-clock second mint distinct ARNs, keep both records, and embed
// the table's resource path the documented format prescribes — the ARN
// is the record's bucket key, so a shared one would overwrite the
// first record.
func TestExportStoreCreateMintsUniqueTableEmbeddingArn(t *testing.T) {
	store := newBackupGlobalStore(t)
	tableArn := "arn:aws:dynamodb:us-east-1:123456789012:table/Tbl"

	first, err := store.Exports().Create(tableArn, "tbl-id", "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("first export create: %v", err)
	}
	second, err := store.Exports().Create(tableArn, "tbl-id", "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("second export create: %v", err)
	}
	if first.ExportArn == second.ExportArn {
		t.Fatalf("same-second exports minted the same ARN %q — the second record overwrote the first", first.ExportArn)
	}
	for _, export := range []*ExportDescription{first, second} {
		if !strings.HasPrefix(export.ExportArn, tableArn+"/export/") {
			t.Fatalf("export ARN %q lacks the documented table-embedding form %s/export/<id>", export.ExportArn, tableArn)
		}
		id := strings.TrimPrefix(export.ExportArn, tableArn+"/export/")
		if !exportImportIdPattern.MatchString(id) {
			t.Fatalf("export ARN id segment %q does not match the documented timestamp-suffix shape", id)
		}
		if _, err := store.Exports().Get(export.ExportArn); err != nil {
			t.Fatalf("get export %q: %v", export.ExportArn, err)
		}
	}
	exports, _, err := store.Exports().List(tableArn, "", 100)
	if err != nil {
		t.Fatalf("list exports: %v", err)
	}
	if len(exports) != 2 {
		t.Fatalf("same-second exports: list saw %d records, want 2", len(exports))
	}
}

// The import family's ARN minting mirrors the export pin: two
// back-to-back imports of the same table survive as distinct records
// under the table-embedding form.
func TestImportStoreCreateMintsUniqueTableEmbeddingArn(t *testing.T) {
	store := newBackupGlobalStore(t)
	tableArn := "arn:aws:dynamodb:us-east-1:123456789012:table/Tbl"

	first, err := store.Imports().Create(tableArn, "tbl-id")
	if err != nil {
		t.Fatalf("first import create: %v", err)
	}
	second, err := store.Imports().Create(tableArn, "tbl-id")
	if err != nil {
		t.Fatalf("second import create: %v", err)
	}
	if first.ImportArn == second.ImportArn {
		t.Fatalf("same-second imports minted the same ARN %q — the second record overwrote the first", first.ImportArn)
	}
	for _, imp := range []*ImportTableDescription{first, second} {
		if !strings.HasPrefix(imp.ImportArn, tableArn+"/import/") {
			t.Fatalf("import ARN %q lacks the documented table-embedding form %s/import/<id>", imp.ImportArn, tableArn)
		}
		id := strings.TrimPrefix(imp.ImportArn, tableArn+"/import/")
		if !exportImportIdPattern.MatchString(id) {
			t.Fatalf("import ARN id segment %q does not match the documented timestamp-suffix shape", id)
		}
		if _, err := store.Imports().Get(imp.ImportArn); err != nil {
			t.Fatalf("get import %q: %v", imp.ImportArn, err)
		}
	}
	imports, _, err := store.Imports().List(tableArn, "", 100)
	if err != nil {
		t.Fatalf("list imports: %v", err)
	}
	if len(imports) != 2 {
		t.Fatalf("same-second imports: list saw %d records, want 2", len(imports))
	}
}

// TestDeleteTableCascadeRemovesJobRecords pins the cascade's job-record
// rule through the shared deletion core: a table's deletion removes the
// export records whose TableArn addresses the table — matched through the
// ARN parser, since the ARN grammar places a colon before the resource —
// and leaves the records of every other table in place.
func TestDeleteTableCascadeRemovesJobRecords(t *testing.T) {
	store := newBackupGlobalStore(t)

	ownExport, err := store.Exports().Create("arn:aws:dynamodb:us-east-1:123456789012:table/CascadeTbl", "id-1", ExportFormatDynamoDBJSON)
	if err != nil {
		t.Fatalf("create own export: %v", err)
	}
	foreignExport, err := store.Exports().Create("arn:aws:dynamodb:us-east-1:123456789012:table/OtherTbl", "id-2", ExportFormatDynamoDBJSON)
	if err != nil {
		t.Fatalf("create foreign export: %v", err)
	}

	err = store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteTableCascade("CascadeTbl")
	})
	if err != nil {
		t.Fatalf("cascade delete: %v", err)
	}

	if _, err := store.Exports().Get(ownExport.ExportArn); err == nil {
		t.Fatal("the deleted table's export record must be gone")
	}
	if _, err := store.Exports().Get(foreignExport.ExportArn); err != nil {
		t.Fatalf("the other table's export record must survive: %v", err)
	}
}
