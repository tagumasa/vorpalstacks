package rdsdata

import (
	"context"
	"strings"
	"testing"
)

// TestValidateSQLByteCeiling pins the ExecuteStatement / BatchExecuteStatement
// sql ceiling: the API Reference expresses the bound both as "Maximum length
// of 65536" and as "Maximum length of 64 KB", so the guard measures bytes.
// Both operations route through validateSQL before any engine work.
func TestValidateSQLByteCeiling(t *testing.T) {
	if err := validateSQL(strings.Repeat("a", 65536)); err != nil {
		t.Errorf("65536-byte SQL rejected: %v", err)
	}
	if err := validateSQL(strings.Repeat("a", 65537)); err == nil {
		t.Error("65537-byte SQL accepted")
	}
	if err := validateSQL(""); err != nil {
		t.Errorf("empty SQL rejected: %v", err)
	}
	// A rune-legal multibyte statement whose encoded size exceeds the 64 KB
	// byte ceiling stays rejected.
	if err := validateSQL(strings.Repeat("\u65e5", 30000)); err == nil {
		t.Error("30000-character CJK SQL (90000 bytes) accepted")
	}
}

// TestExecuteStatementInTxCoreRejectsEmptyTransactionID pins the
// in-transaction contract: an empty transactionId is rejected before any
// validation or execution step, so an EventBus caller can never degrade a
// transactional statement into an autocommit execution.
func TestExecuteStatementInTxCoreRejectsEmptyTransactionID(t *testing.T) {
	svc := &RDSDataService{}
	_, err := svc.executeStatementInTxCore(context.Background(), &ExecuteStatementInput{
		ResourceArn:   "arn:aws:rds:us-east-1:123456789012:cluster:example",
		Sql:           "SELECT 1",
		TransactionID: "",
	})
	if err == nil {
		t.Fatal("expected a transactionId validation error, got nil")
	}
	if !strings.Contains(err.Error(), "TransactionNotFoundException") {
		t.Fatalf("error %q is not a TransactionNotFoundException", err.Error())
	}
}
