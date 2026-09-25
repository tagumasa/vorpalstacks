// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"testing"

	"vorpalstacks/internal/common/request"
)

// TestDescriptionsEmitStoredMembers pins the field-completeness gaps the
// descriptions owed: DescribeBackup's SourceTableDetails carries the
// backup record's billing mode, ListBackups' BackupSummary carries the
// same mode as a BillingModeSummary, and DescribeExport emits the record's
// ClientToken — the member DescribeImport already emitted on its side.
func TestDescriptionsEmitStoredMembers(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "DescTbl",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	backup, err := svc.createBackupCore(ctx, reqCtx, createBackupInput{Parameters: map[string]interface{}{
		"TableName": "DescTbl", "BackupName": "desc-bk",
	}})
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}

	describe, err := svc.DescribeBackup(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"BackupArn": backup.BackupArn,
	}})
	if err != nil {
		t.Fatalf("describe backup: %v", err)
	}
	source := describe.(map[string]interface{})["BackupDescription"].(map[string]interface{})["SourceTableDetails"].(map[string]interface{})
	if source["BillingMode"] != "PAY_PER_REQUEST" {
		t.Fatalf("SourceTableDetails.BillingMode = %v, want PAY_PER_REQUEST", source["BillingMode"])
	}

	listed, err := svc.ListBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "DescTbl",
	}})
	if err != nil {
		t.Fatalf("list backups: %v", err)
	}
	summary := listed.(map[string]interface{})["BackupSummaries"].([]map[string]interface{})[0]
	bms, ok := summary["BillingModeSummary"].(map[string]interface{})
	if !ok || bms["BillingMode"] != "PAY_PER_REQUEST" {
		t.Fatalf("BackupSummary.BillingModeSummary = %v, want PAY_PER_REQUEST", summary["BillingModeSummary"])
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	table, err := store.Tables().Get("DescTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	export, err := store.Exports().Create(table.ARN, table.TableId, "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("create export record: %v", err)
	}
	export.ClientToken = "desc-token"
	if err := store.Exports().Put(export); err != nil {
		t.Fatalf("put export record: %v", err)
	}

	exportDesc, err := svc.DescribeExport(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ExportArn": export.ExportArn,
	}})
	if err != nil {
		t.Fatalf("describe export: %v", err)
	}
	ed := exportDesc.(map[string]interface{})["ExportDescription"].(map[string]interface{})
	if ed["ClientToken"] != "desc-token" {
		t.Fatalf("ExportDescription.ClientToken = %v, want desc-token", ed["ClientToken"])
	}
}

// TestDescribeByArnErrorPaths pins the shared describe shape's error
// halves on every ARN-addressed family: a too-short ARN fails the length
// validation, and a length-valid but unknown ARN returns the family's
// not-found sentinel.
func TestDescribeByArnErrorPaths(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	unknown := "arn:aws:dynamodb:us-east-1:123456789012:table/UnknownTbl/export/01234567-89ab-cdef-0123-456789abcdef"

	if _, err := svc.describeBackupCore(ctx, reqCtx, "too-short"); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("backup describe with a short ARN: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := svc.describeExportCore(ctx, reqCtx, "too-short"); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("export describe with a short ARN: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := svc.describeImportCore(ctx, reqCtx, "too-short"); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("import describe with a short ARN: err = %v, want ErrInvalidParameter", err)
	}

	if _, err := svc.describeBackupCore(ctx, reqCtx, unknown); !errors.Is(err, ErrBackupNotFound) {
		t.Errorf("backup describe of an unknown ARN: err = %v, want ErrBackupNotFound", err)
	}
	if _, err := svc.describeExportCore(ctx, reqCtx, unknown); !errors.Is(err, ErrExportNotFound) {
		t.Errorf("export describe of an unknown ARN: err = %v, want ErrExportNotFound", err)
	}
	if _, err := svc.describeImportCore(ctx, reqCtx, unknown); !errors.Is(err, ErrImportNotFound) {
		t.Errorf("import describe of an unknown ARN: err = %v, want ErrImportNotFound", err)
	}
}
