package dynamodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestSingletonCancellationReasonClassification(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		code    string
		message string
		cancels bool
	}{
		{
			name:    "duplicate-key insert cancels",
			err:     ErrConditionalCheckFailed,
			code:    "ConditionalCheckFailed",
			message: "The conditional request failed",
		},
		{
			name:    "clause type mismatch cancels",
			err:     ErrTypeMismatch,
			code:    "ValidationError",
			message: "Type mismatch for attribute to update.",
		},
		{
			name:    "wrapped clause mismatch still detected",
			err:     fmt.Errorf("apply clause: %w", ErrTypeMismatch),
			code:    "ValidationError",
			message: "Type mismatch for attribute to update.",
		},
		{
			name: "multi-item match cancels",
			err: NewAPIError("com.amazon.coral.validate#ValidationException",
				"UPDATE statement must match exactly one item", http.StatusBadRequest),
			code:    "ValidationError",
			message: "UPDATE statement must match exactly one item",
		},
		{
			name:    "coral validation cancels as ValidationError",
			err:     ErrInvalidParameter,
			code:    "ValidationError",
			message: "Invalid parameter",
		},
		{
			name:    "unknown table cancels as ResourceNotFound",
			err:     ErrTableNotFound,
			code:    "ResourceNotFound",
			message: "Requested resource not found: Table not found",
		},
		{
			name:    "storage fault cancels as InternalServerError",
			err:     errors.New("pebble: closed"),
			code:    "InternalServerError",
			message: "pebble: closed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reason := singletonCancellationReason(tc.err)
			if reason.Code != tc.code {
				t.Fatalf("code = %q, want %q", reason.Code, tc.code)
			}
			if reason.Message != tc.message {
				t.Fatalf("message = %q, want %q", reason.Message, tc.message)
			}
		})
	}
}

// The statement-error classifier serves both consuming shapes; the pin
// fixes the per-plane renderings the consolidation preserved: the batch
// plane's enum-bound default (InternalServerError) against the
// cancellation plane's free-string default (the exception's own short
// code), and the known-class renderings carrying the exception's own
// message on both planes alike.
func TestStatementErrorPairPerPlaneDefaults(t *testing.T) {
	resourceInUse := NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException",
		"Table is not in ACTIVE state", http.StatusBadRequest)
	if code, _ := partiqlStatementErrorPair(resourceInUse, partiqlErrorBatch); code != "InternalServerError" {
		t.Fatalf("batch default: code = %q, want InternalServerError", code)
	}
	if code, _ := partiqlStatementErrorPair(resourceInUse, partiqlErrorCancellation); code != "ResourceInUse" {
		t.Fatalf("cancellation default: code = %q, want ResourceInUse", code)
	}

	duplicate := NewAPIError("com.amazonaws.dynamodb.v20120810#DuplicateItemException",
		"There was an attempt to insert an item with the same primary key as an item that already exists", http.StatusBadRequest)
	for _, plane := range []partiqlErrorPlane{partiqlErrorBatch, partiqlErrorCancellation} {
		code, message := partiqlStatementErrorPair(duplicate, plane)
		if code != "DuplicateItem" || message != duplicate.Message {
			t.Fatalf("plane %d duplicate: %q/%q", plane, code, message)
		}
	}

	conflict := NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionConflictException",
		"Operation was rejected because there is an ongoing transaction for the item", http.StatusBadRequest)
	if code, message := partiqlStatementErrorPair(conflict, partiqlErrorBatch); code != "TransactionConflict" || message != conflict.Message {
		t.Fatalf("batch conflict: %q/%q", code, message)
	}
}

// The wire contract for the envelope: unaffected statements carry the
// literal code "None" and omit the message entirely, while the failing
// statement's reason keeps both members.
func TestTransactionCanceledEnvelopeRendersNoneWithoutMessage(t *testing.T) {
	canceled := NewTransactionCanceledError("Transaction canceled", []CancellationReason{
		{Code: "None"},
		{Code: "ConditionalCheckFailed", Message: "The conditional request failed"},
	})

	var payload struct {
		Type                string `json:"__type"`
		CancellationReasons []struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"CancellationReasons"`
	}
	if err := json.Unmarshal([]byte(canceled.ToJSON()), &payload); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if !strings.HasSuffix(payload.Type, "TransactionCanceledException") {
		t.Fatalf("type = %q, want TransactionCanceledException", payload.Type)
	}
	if len(payload.CancellationReasons) != 2 {
		t.Fatalf("reasons = %d, want 2", len(payload.CancellationReasons))
	}
	if payload.CancellationReasons[0].Code != "None" {
		t.Fatalf("unaffected reason code = %q, want None", payload.CancellationReasons[0].Code)
	}
	if payload.CancellationReasons[0].Message != "" {
		t.Fatalf("unaffected reason must omit the message, got %q", payload.CancellationReasons[0].Message)
	}
	if payload.CancellationReasons[1].Code != "ConditionalCheckFailed" ||
		payload.CancellationReasons[1].Message != "The conditional request failed" {
		t.Fatalf("failing reason = %+v, want ConditionalCheckFailed with message", payload.CancellationReasons[1])
	}
	if strings.Contains(canceled.ToJSON(), `"Message":""`) {
		t.Fatalf("None reason must be serialised without a Message member")
	}
}

// TestEmptyResponseSlotsOmitTheItemMember pins the presence convention on
// the multi-statement planes: a response slot with nothing to carry — a
// write statement's slot, a read that matched no item, a transaction get
// for a missing key — omits the Item member entirely; the AWS JSON
// protocol never serialises an unset member as an explicit null.
func TestEmptyResponseSlotsOmitTheItemMember(t *testing.T) {
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	defer sm.Close()
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)

	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "EmptySlotTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "pk", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	reqCtx := request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")

	// A write-only transaction: every statement's slot carries no Item
	// member at all.
	resp, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactStatements": []interface{}{
			map[string]interface{}{"Statement": `INSERT INTO "EmptySlotTable" VALUE {'pk': 'a'}`},
		},
	}})
	if err != nil {
		t.Fatalf("write transaction: %v", err)
	}
	slots := resp.(map[string]interface{})["Responses"].([]map[string]interface{})
	if len(slots) != 1 {
		t.Fatalf("write transaction: expected one slot, got %#v", slots)
	}
	if _, present := slots[0]["Item"]; present {
		t.Fatalf("write statement slot: Item member must stay absent, got %#v", slots[0])
	}

	// A transaction read for a missing key: the empty slot omits the
	// member rather than carrying a null.
	resp, err = svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{
			map[string]interface{}{"Get": map[string]interface{}{
				"Key":       map[string]interface{}{"pk": map[string]interface{}{"S": "absent"}},
				"TableName": "EmptySlotTable",
			}},
		},
	}})
	if err != nil {
		t.Fatalf("transact get: %v", err)
	}
	responses := resp.(map[string]interface{})["Responses"].([]map[string]interface{})
	if len(responses) != 1 {
		t.Fatalf("transact get: expected one response, got %#v", responses)
	}
	if _, present := responses[0]["Item"]; present {
		t.Fatalf("missing-key transact get: Item member must stay absent, got %#v", responses[0])
	}

	// A batch read that matches no item: same omission.
	resp, err = svc.BatchExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statements": []interface{}{
			map[string]interface{}{"Statement": `SELECT * FROM "EmptySlotTable" WHERE pk = 'absent'`},
		},
	}})
	if err != nil {
		t.Fatalf("batch read: %v", err)
	}
	batchSlots := resp.(map[string]interface{})["Responses"].([]map[string]interface{})
	if len(batchSlots) != 1 {
		t.Fatalf("batch read: expected one slot, got %#v", batchSlots)
	}
	if _, present := batchSlots[0]["Item"]; present {
		t.Fatalf("no-match batch read: Item member must stay absent, got %#v", batchSlots[0])
	}
}
