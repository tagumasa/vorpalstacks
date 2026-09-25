package dynamodb

import (
	"context"
	"errors"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
)

// The boolean request parameters share one strict regime: an absent key
// means the documented default (false for ConsistentRead and
// DeletionProtectionEnabled), and a present non-boolean value is rejected
// with ValidationException rather than silently coerced — the contract the
// AWS JSON protocol enforces by typing booleans at deserialisation. Each
// test drives the wire-plane method the SDK client itself reaches; the
// boolean values themselves are exercised by every other test in the
// package, so the pins here carry the absent-default success and the
// string-form rejection.

func boolPinKey() map[string]interface{} {
	return map[string]interface{}{
		"id": map[string]interface{}{"S": "A"},
		"sk": map[string]interface{}{"S": "1"},
	}
}

func expectInvalidParameter(t *testing.T, what string, err error) {
	t.Helper()
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("%s: expected ErrInvalidParameter, got %v", what, err)
	}
}

func TestGetItemConsistentReadStrictBool(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()

	if _, err := svc.GetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"Key":       boolPinKey(),
	}}); err != nil {
		t.Fatalf("absent ConsistentRead (default false): unexpected error %v", err)
	}

	_, err := svc.GetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":      "LegacyTable",
		"Key":            boolPinKey(),
		"ConsistentRead": "true",
	}})
	expectInvalidParameter(t, "string ConsistentRead", err)
}

func TestScanConsistentReadStrictBool(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()

	if _, err := svc.Scan(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
	}}); err != nil {
		t.Fatalf("absent ConsistentRead (default false): unexpected error %v", err)
	}

	_, err := svc.Scan(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":      "LegacyTable",
		"ConsistentRead": "true",
	}})
	expectInvalidParameter(t, "string ConsistentRead", err)
}

func TestBatchGetItemConsistentReadStrictBool(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()

	if _, err := svc.BatchGetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"LegacyTable": map[string]interface{}{"Keys": []interface{}{boolPinKey()}},
		},
	}}); err != nil {
		t.Fatalf("absent ConsistentRead (default false): unexpected error %v", err)
	}

	_, err := svc.BatchGetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"LegacyTable": map[string]interface{}{"Keys": []interface{}{boolPinKey()}, "ConsistentRead": "true"},
		},
	}})
	expectInvalidParameter(t, "string ConsistentRead", err)
}

func TestExecuteStatementConsistentReadStrictBool(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()
	stmt := `SELECT * FROM "LegacyTable" WHERE id = 'A' AND sk = '1'`

	if _, err := svc.ExecuteStatement(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement": stmt,
	}}); err != nil {
		t.Fatalf("absent ConsistentRead (default false): unexpected error %v", err)
	}

	_, err := svc.ExecuteStatement(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement":      stmt,
		"ConsistentRead": "true",
	}})
	expectInvalidParameter(t, "string ConsistentRead", err)
}

func TestBatchExecuteStatementConsistentReadStrictBool(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()
	stmt := `SELECT * FROM "LegacyTable" WHERE id = 'A' AND sk = '1'`

	resp, err := svc.BatchExecuteStatement(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statements": []interface{}{
			map[string]interface{}{"Statement": stmt, "ConsistentRead": "true"},
		},
	}})
	if err != nil {
		t.Fatalf("per-statement rejection reports in the statement's own response slot, not the request: %v", err)
	}
	responses, ok := resp.(map[string]interface{})["Responses"].([]map[string]interface{})
	if !ok || len(responses) != 1 {
		t.Fatalf("batch response shape: %#v", resp)
	}
	errSlot, ok := responses[0]["Error"].(map[string]interface{})
	if !ok || errSlot["Code"] != "ValidationError" {
		t.Fatalf("string ConsistentRead: expected a per-statement ValidationError, got %#v", responses[0])
	}
	if msg, ok := errSlot["Message"].(string); !ok || !strings.Contains(msg, "ConsistentRead must be a boolean") {
		t.Fatalf("string ConsistentRead: expected the boolean sentence, got %#v", errSlot["Message"])
	}
}

func TestCreateTableDeletionProtectionStrictBool(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()
	create := func(extra map[string]interface{}) error {
		params := map[string]interface{}{
			"TableName":            "BoolPinTable",
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "id", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "id", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		}
		for k, v := range extra {
			params[k] = v
		}
		_, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}

	if err := create(nil); err != nil {
		t.Fatalf("absent DeletionProtectionEnabled (default false): unexpected error %v", err)
	}

	err := create(map[string]interface{}{"DeletionProtectionEnabled": "true"})
	expectInvalidParameter(t, "string DeletionProtectionEnabled", err)
}

func TestUpdateTableDeletionProtectionStrictBool(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "BoolPinUpdateTable",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "id", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "id", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("seed table: %v", err)
	}

	// A present boolean false is an explicit setting, not an absent member.
	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "BoolPinUpdateTable",
		"DeletionProtectionEnabled": false,
	}}); err != nil {
		t.Fatalf("bool false DeletionProtectionEnabled: unexpected error %v", err)
	}

	_, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "BoolPinUpdateTable",
		"DeletionProtectionEnabled": "true",
	}})
	expectInvalidParameter(t, "string DeletionProtectionEnabled", err)
}
