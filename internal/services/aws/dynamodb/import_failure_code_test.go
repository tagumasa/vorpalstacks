// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// invokerShapedError reproduces the invoker's wrapping so the failure-code
// mapping sees the error exactly as the S3 plane delivers it.
func invokerShapedError(class error) error {
	return fmt.Errorf("s3 GetObject bucket/key: %w", commonstore.NewStoreError("s3", "get_proto", class))
}

// TestImportFailureCodesFollowTheErrorClass pins the class-driven failure
// codes: an in-job table-creation failure reports TableAlreadyExists only
// for the duplicate-name race, the not-found store class on the S3 plane
// reports S3NoSuchBucket, other S3 errors keep the access code, and the
// export's put failures follow the same mapping.
func TestImportFailureCodesFollowTheErrorClass(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "FcTbl",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	table, err := store.Tables().Get("FcTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}

	runImport := func(target string, invoker *failingS3Invoker) *dbstore.ImportTableDescription {
		t.Helper()
		svc.SetEventBus(ionExportBus{s3: invoker})
		imp, err := store.Imports().Create(table.ARN, table.TableId)
		if err != nil {
			t.Fatalf("create import record: %v", err)
		}
		svc.runImportJob(store, "us-east-1", ImportTableCoreInput{
			TableName:     target,
			KeySchema:     table.KeySchema,
			AttributeDefs: table.AttributeDefinitions,
			BillingMode:   dbstore.BillingModePayPerRequest,
			S3Bucket:      "src-bucket",
			InputFormat:   "DYNAMODB_JSON",
		}, imp.ImportArn)
		final, err := store.Imports().Get(imp.ImportArn)
		if err != nil {
			t.Fatalf("reload import: %v", err)
		}
		return final
	}

	// The class mapping itself: the already-exists sentinel names the
	// duplicate race, anything else is internal.
	if got := importCreateFailureCode(dbstore.ErrTableAlreadyExists); got != "TableAlreadyExists" {
		t.Fatalf("already-exists class: %s, want TableAlreadyExists", got)
	}
	if got := importCreateFailureCode(errors.New("store unavailable")); got != "InternalFailure" {
		t.Fatalf("other class: %s, want InternalFailure", got)
	}

	// The duplicate-name race end to end: the target already exists, so
	// the in-job Create fails with the already-exists class and the code
	// names it.
	dup := runImport("FcTbl", newFailingS3Invoker())
	if dup.ImportStatus != "FAILED" || dup.FailureCode != "TableAlreadyExists" {
		t.Fatalf("duplicate target: %s/%s, want FAILED/TableAlreadyExists", dup.ImportStatus, dup.FailureCode)
	}

	// The not-found class on the S3 plane names the missing resource.
	notFound := invokerShapedError(commonstore.ErrNotFound)
	missing := runImport("FcMissing", &failingS3Invoker{failExists: notFound})
	if missing.FailureCode != "S3NoSuchBucket" {
		t.Fatalf("not-found bucket check: %s, want S3NoSuchBucket", missing.FailureCode)
	}
	missingObj := runImport("FcMissingObj", &failingS3Invoker{failList: notFound})
	if missingObj.FailureCode != "S3NoSuchBucket" {
		t.Fatalf("not-found list: %s, want S3NoSuchBucket", missingObj.FailureCode)
	}

	// Other S3 errors keep the access-class code.
	denied := runImport("FcDenied", &failingS3Invoker{failExists: errors.New("network unreachable")})
	if denied.FailureCode != "S3AccessDenied" {
		t.Fatalf("other bucket check: %s, want S3AccessDenied", denied.FailureCode)
	}

	// The export's put failures follow the same mapping.
	runExport := func(invoker *failingS3Invoker) *dbstore.ExportDescription {
		t.Helper()
		svc.SetEventBus(ionExportBus{s3: invoker})
		export, err := store.Exports().Create(table.ARN, table.TableId, "DYNAMODB_JSON")
		if err != nil {
			t.Fatalf("create export record: %v", err)
		}
		svc.runExportJob(store, ExportTableCoreInput{
			TableArn:     table.ARN,
			TableName:    "FcTbl",
			ExportFormat: "DYNAMODB_JSON",
			S3Bucket:     "dest-bucket",
			Region:       "us-east-1",
			ExportTime:   time.Now(),
		}, export.ExportArn)
		final, err := store.Exports().Get(export.ExportArn)
		if err != nil {
			t.Fatalf("reload export: %v", err)
		}
		return final
	}
	exportMissing := runExport(&failingS3Invoker{failPut: invokerShapedError(commonstore.ErrNotFound)})
	if exportMissing.ExportStatus != "FAILED" || exportMissing.FailureCode != "S3NoSuchBucket" {
		t.Fatalf("export not-found put: %s/%s, want FAILED/S3NoSuchBucket", exportMissing.ExportStatus, exportMissing.FailureCode)
	}
	exportDenied := runExport(&failingS3Invoker{failPut: errors.New("bucket is read-only")})
	if exportDenied.FailureCode != "S3AccessDenied" {
		t.Fatalf("export other put: %s, want S3AccessDenied", exportDenied.FailureCode)
	}
}

func newFailingS3Invoker() *failingS3Invoker {
	return &failingS3Invoker{}
}

// importSourceInvoker serves a fixed set of S3 objects to an import job:
// the bucket exists and the listed keys return their bytes, so the job
// runs its real per-object loop against controlled content.
type importSourceInvoker struct {
	failingS3Invoker
	objects map[string][]byte
	order   []string
}

func (s *importSourceInvoker) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	return s.order, nil
}

func (s *importSourceInvoker) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	data, ok := s.objects[key]
	if !ok {
		return nil, fmt.Errorf("no such object %s", key)
	}
	return data, nil
}

// TestImportFailurePersistsCountedFigures pins the failure path's counter
// contract: a job that already processed and imported objects before
// failing mid-run persists the counted figures — the counters are
// observable state a client polling DescribeImport reads off either
// terminal state, never discarded by the failure.
func TestImportFailurePersistsCountedFigures(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "CountSrcTbl",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create source table: %v", err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	table, err := store.Tables().Get("CountSrcTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}

	var good bytes.Buffer
	gz := gzip.NewWriter(&good)
	if _, err := gz.Write([]byte(`{"Item":{"pk":{"S":"ok1"},"v":{"S":"one"}}}` + "\n")); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}

	svc.SetEventBus(ionExportBus{s3: &importSourceInvoker{
		objects: map[string][]byte{"a.json.gz": good.Bytes(), "b.json.gz": []byte("not a gzip stream")},
		order:   []string{"a.json.gz", "b.json.gz"},
	}})
	imp, err := store.Imports().Create(table.ARN, table.TableId)
	if err != nil {
		t.Fatalf("create import record: %v", err)
	}
	svc.runImportJob(store, "us-east-1", ImportTableCoreInput{
		TableName:            "ImpCountTarget",
		KeySchema:            table.KeySchema,
		AttributeDefs:        table.AttributeDefinitions,
		BillingMode:          dbstore.BillingModePayPerRequest,
		S3Bucket:             "src-bucket",
		InputFormat:          "DYNAMODB_JSON",
		InputCompressionType: "GZIP",
	}, imp.ImportArn)

	final, err := store.Imports().Get(imp.ImportArn)
	if err != nil {
		t.Fatalf("reload import: %v", err)
	}
	if final.ImportStatus != "FAILED" || final.FailureCode != "GZIPError" {
		t.Fatalf("mid-run failure = %s/%s, want FAILED/GZIPError", final.ImportStatus, final.FailureCode)
	}
	if final.ProcessedItemCount != 1 || final.ImportedItemCount != 1 {
		t.Fatalf("counted figures on the failed record = processed %d imported %d, want 1/1",
			final.ProcessedItemCount, final.ImportedItemCount)
	}
	if final.ProcessedSizeBytes == 0 {
		t.Fatal("the failed record discarded the processed size of the first object")
	}
}
