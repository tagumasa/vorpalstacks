// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestTransactWriteRejectsOversizedItems(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnSizeTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	key := map[string]interface{}{"id": sVal("a")}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnSizeTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("small")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	oversized := strings.Repeat("x", dbstore.MaxItemSizeBytes+1)

	// An Update expression grows the small item past the cap.
	_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"TableName":                 "TxnSizeTable",
			"Key":                       key,
			"UpdateExpression":          "SET v = :big",
			"ExpressionAttributeValues": map[string]interface{}{":big": sVal(oversized)},
		}}},
	}})
	var canceled *TransactionCanceledError
	if err == nil || !errors.As(err, &canceled) {
		t.Fatalf("oversizing update: expected TransactionCanceledException, got %v", err)
	}
	if len(canceled.CancellationReasons) != 1 || canceled.CancellationReasons[0].Code != "ValidationError" {
		t.Fatalf("oversizing update: expected a ValidationError cancellation reason, got %+v", canceled.CancellationReasons)
	}
	if !strings.Contains(canceled.CancellationReasons[0].Message, "Item size has exceeded the maximum allowed size") {
		t.Fatalf("oversizing update: expected the item-size message, got %+v", canceled.CancellationReasons[0])
	}

	// A Put of an oversized whole item is the same documented rejection.
	_, err = svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{
			"TableName": "TxnSizeTable",
			"Item":      map[string]interface{}{"id": sVal("b"), "v": sVal(oversized)},
		}}},
	}})
	if err == nil || !errors.As(err, &canceled) || canceled.CancellationReasons[0].Code != "ValidationError" {
		t.Fatalf("oversizing put: expected TransactionCanceledException with a ValidationError reason, got %v", err)
	}

	// The rejected transactions left the stored item untouched and the
	// second item unwritten; an in-cap update still commits.
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnSizeTable", "Key": key,
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	if got := resp.(map[string]interface{})["Item"].(map[string]interface{})["v"].(map[string]interface{})["S"]; got != "small" {
		t.Fatalf("stored item after rejection: expected v=small, got %v", got)
	}
	respB, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnSizeTable", "Key": map[string]interface{}{"id": sVal("b")},
	}})
	if err != nil {
		t.Fatalf("missing-target get: %v", err)
	}
	if _, present := respB.(map[string]interface{})["Item"]; present {
		t.Fatal("oversized put target: expected the item to not exist after the rejected transaction")
	}
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"TableName":                 "TxnSizeTable",
			"Key":                       key,
			"UpdateExpression":          "SET v = :ok",
			"ExpressionAttributeValues": map[string]interface{}{":ok": sVal("committed")},
		}}},
	}}); err != nil {
		t.Fatalf("in-cap update: unexpected error %v", err)
	}
}

func TestTransactWritePutRejectsEmptyKeyValue(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnEmptyKeyTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	txnPut := func(item map[string]interface{}) error {
		_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{
				"TableName": "TxnEmptyKeyTable",
				"Item":      item,
			}}},
		}})
		return err
	}

	err := txnPut(map[string]interface{}{"id": sVal(""), "v": sVal("x")})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty key value: expected the parse-phase ValidationException, got %v", err)
	}
	var canceled *TransactionCanceledError
	if errors.As(err, &canceled) {
		t.Fatal("empty key value: the rejection must not be a TransactionCanceledException")
	}
	// The empty key value is invalid on every plane (an item addressed by
	// it cannot be read back), so the parse-phase rejection above is the
	// observable: the request never reached execution.

	// A well-formed Put in the same position commits.
	if err := txnPut(map[string]interface{}{"id": sVal("ok"), "v": sVal("x")}); err != nil {
		t.Fatalf("valid put: unexpected error %v", err)
	}
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnEmptyKeyTable", "Key": map[string]interface{}{"id": sVal("ok")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	if got := resp.(map[string]interface{})["Item"].(map[string]interface{})["v"].(map[string]interface{})["S"]; got != "x" {
		t.Fatalf("valid put: expected v=x, got %v", got)
	}
}

func TestTransactWriteRejectsDuplicateItemActions(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnDupTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	key := func(id string) map[string]interface{} { return map[string]interface{}{"id": sVal(id)} }

	// duplicateExpect pins the documented cancellation envelope: the
	// reasons are ordered per requested item, the offending action's slot
	// carries the ValidationError pair, and every earlier slot is None.
	duplicateExpect := func(t *testing.T, err error, failLabel string) {
		t.Helper()
		var canceled *TransactionCanceledError
		if !errors.As(err, &canceled) {
			t.Fatalf("%s: expected TransactionCanceledException, got %v", failLabel, err)
		}
		if len(canceled.CancellationReasons) != 2 {
			t.Fatalf("%s: expected 2 cancellation reasons, got %d", failLabel, len(canceled.CancellationReasons))
		}
		if got := canceled.CancellationReasons[0].Code; got != "None" {
			t.Fatalf("%s: first reason code = %q, want None", failLabel, got)
		}
		want := CancellationReason{Code: "ValidationError", Message: "One or more parameter values were invalid."}
		if got := canceled.CancellationReasons[1]; got.Code != want.Code || got.Message != want.Message {
			t.Fatalf("%s: second reason = %+v, want %+v", failLabel, got, want)
		}
	}

	// Two Puts on one item: the cancellation is answered before any
	// execution, so nothing is written.
	_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{
			map[string]interface{}{"Put": map[string]interface{}{"TableName": "TxnDupTable", "Item": map[string]interface{}{"id": sVal("a"), "v": sVal("1")}}},
			map[string]interface{}{"Put": map[string]interface{}{"TableName": "TxnDupTable", "Item": map[string]interface{}{"id": sVal("a"), "v": sVal("2")}}},
		},
	}})
	if err == nil {
		t.Fatal("duplicate put: expected the same-item cancellation")
	}
	duplicateExpect(t, err, "duplicate put")

	// A ConditionCheck on an item a write action already targets is the
	// same duplicate.
	_, err = svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{
			map[string]interface{}{"Update": map[string]interface{}{
				"TableName": "TxnDupTable", "Key": key("a"),
				"UpdateExpression":          "SET v = :v",
				"ExpressionAttributeValues": map[string]interface{}{":v": sVal("1")},
			}},
			map[string]interface{}{"ConditionCheck": map[string]interface{}{
				"TableName": "TxnDupTable", "Key": key("a"),
				"ConditionExpression": "attribute_exists(v)",
			}},
		},
	}})
	if err == nil {
		t.Fatal("update plus condition check: expected the same-item cancellation")
	}
	duplicateExpect(t, err, "update plus condition check")

	// Two ConditionChecks on one item are the same duplicate — the rule
	// is uniform across all four actions.
	conditionCheck := func() map[string]interface{} {
		return map[string]interface{}{"ConditionCheck": map[string]interface{}{
			"TableName":           "TxnDupTable",
			"Key":                 key("a"),
			"ConditionExpression": "attribute_exists(v)",
		}}
	}
	_, err = svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{conditionCheck(), conditionCheck()},
	}})
	if err == nil {
		t.Fatal("condition check pair: expected the same-item cancellation")
	}
	duplicateExpect(t, err, "condition check pair")

	// Nothing was written by the rejected requests; distinct items in one
	// transaction keep committing.
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnDupTable", "Key": key("a"),
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	if _, present := resp.(map[string]interface{})["Item"]; present {
		t.Fatal("duplicate put: expected no item after the rejected transactions")
	}
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{
			map[string]interface{}{"Put": map[string]interface{}{"TableName": "TxnDupTable", "Item": map[string]interface{}{"id": sVal("a"), "v": sVal("ok")}}},
			map[string]interface{}{"Put": map[string]interface{}{"TableName": "TxnDupTable", "Item": map[string]interface{}{"id": sVal("b"), "v": sVal("ok")}}},
		},
	}}); err != nil {
		t.Fatalf("distinct items: unexpected error %v", err)
	}
}

func TestTransactGetRejectsDuplicateTargets(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnGetDupTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	key := func(id string) map[string]interface{} { return map[string]interface{}{"id": sVal(id)} }
	get := func(id string) map[string]interface{} {
		return map[string]interface{}{"Get": map[string]interface{}{"TableName": "TxnGetDupTable", "Key": key(id)}}
	}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnGetDupTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("x")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	_, err := svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{get("a"), get("a")},
	}})
	if err == nil || !strings.Contains(err.Error(), "Transaction request cannot include multiple operations on one item") {
		t.Fatalf("duplicate get: expected the same-item ValidationException, got %v", err)
	}

	// Distinct items — including a missing one, which answers an empty
	// slot — keep serving two responses.
	resp, err := svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{get("a"), get("b")},
	}})
	if err != nil {
		t.Fatalf("distinct gets: unexpected error %v", err)
	}
	responses := resp.(map[string]interface{})["Responses"].([]map[string]interface{})
	if len(responses) != 2 {
		t.Fatalf("distinct gets: expected two responses, got %d", len(responses))
	}
	if got := responses[0]["Item"].(map[string]interface{})["v"].(map[string]interface{})["S"]; got != "x" {
		t.Fatalf("distinct gets: expected the stored v=x, got %v", got)
	}
	if item := responses[1]["Item"]; item != nil {
		t.Fatalf("distinct gets: expected the missing item to answer an empty slot, got %v", item)
	}
}

// TestTransactGetConsumedCapacityFollowsRequestOrder pins the
// ConsumedCapacity entry order on the transactional read plane:
// TransactItems is an ordered list, so the per-table entries answer in
// first-appearance request order, never the accumulation map's random
// order. The request names the lexicographically later table first, so a
// sorted or map-ordered build answers differently; the loop makes the
// assertion decisive — a single draw of a random order could accidentally
// match.
func TestTransactGetConsumedCapacityFollowsRequestOrder(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnCapOrderB",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "TxnCapOrderA",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create second table: %v", err)
	}
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	get := func(table string) map[string]interface{} {
		return map[string]interface{}{"Get": map[string]interface{}{
			"TableName": table,
			"Key":       map[string]interface{}{"id": sVal("a")},
		}}
	}
	for i := 0; i < 20; i++ {
		resp, err := svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems":          []interface{}{get("TxnCapOrderB"), get("TxnCapOrderA")},
			"ReturnConsumedCapacity": "TOTAL",
		}})
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		caps, ok := resp.(map[string]interface{})["ConsumedCapacity"].([]map[string]interface{})
		if !ok || len(caps) != 2 {
			t.Fatalf("iteration %d: expected two ConsumedCapacity entries, got %T %v", i, resp.(map[string]interface{})["ConsumedCapacity"], resp.(map[string]interface{})["ConsumedCapacity"])
		}
		if got := caps[0]["TableName"]; got != "TxnCapOrderB" || caps[1]["TableName"] != "TxnCapOrderA" {
			t.Fatalf("iteration %d: entries must follow first-appearance request order [TxnCapOrderB, TxnCapOrderA], got [%v, %v]", i, got, caps[1]["TableName"])
		}
	}
}

// TestBatchWriteConsumedCapacityAnswersInSortedOrder pins the entry order
// on the batch write plane: RequestItems is a map carrying no request
// order, so the per-table entries answer in sorted table-name order — the
// batch read plane's choice, deterministic where the map's random
// iteration is not. The loop makes the assertion decisive against a
// map-ordered build.
func TestBatchWriteConsumedCapacityAnswersInSortedOrder(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "BatchWriteCapOrderB",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "BatchWriteCapOrderA",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create second table: %v", err)
	}
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	put := func() interface{} {
		return map[string]interface{}{"PutRequest": map[string]interface{}{"Item": map[string]interface{}{"id": sVal("a")}}}
	}
	requestItems := map[string]interface{}{
		"BatchWriteCapOrderB": []interface{}{put()},
		"BatchWriteCapOrderA": []interface{}{put()},
	}
	for i := 0; i < 20; i++ {
		resp, err := svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"RequestItems":           requestItems,
			"ReturnConsumedCapacity": "TOTAL",
		}})
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		capsRaw, ok := resp.(map[string]interface{})["ConsumedCapacity"].([]interface{})
		if !ok || len(capsRaw) != 2 {
			t.Fatalf("iteration %d: expected two ConsumedCapacity entries, got %T %v", i, resp.(map[string]interface{})["ConsumedCapacity"], resp.(map[string]interface{})["ConsumedCapacity"])
		}
		first := capsRaw[0].(map[string]interface{})["TableName"]
		second := capsRaw[1].(map[string]interface{})["TableName"]
		if first != "BatchWriteCapOrderA" || second != "BatchWriteCapOrderB" {
			t.Fatalf("iteration %d: entries must answer in sorted order [BatchWriteCapOrderA, BatchWriteCapOrderB], got [%v, %v]", i, first, second)
		}
	}
}

// TestBatchGetConsumedCapacityAnswersInSortedOrder pins the entry order on
// the batch read plane: RequestItems is a map carrying no request order,
// so the per-table entries answer in sorted table-name order —
// deterministic where the accumulation map's random iteration is not. The
// loop makes the assertion decisive against a map-ordered build.
func TestBatchGetConsumedCapacityAnswersInSortedOrder(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "BatchCapOrderB",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "BatchCapOrderA",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create second table: %v", err)
	}
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	requestItems := map[string]interface{}{
		"BatchCapOrderB": map[string]interface{}{"Keys": []interface{}{map[string]interface{}{"id": sVal("a")}}},
		"BatchCapOrderA": map[string]interface{}{"Keys": []interface{}{map[string]interface{}{"id": sVal("a")}}},
	}
	for i := 0; i < 20; i++ {
		resp, err := svc.BatchGetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"RequestItems":           requestItems,
			"ReturnConsumedCapacity": "TOTAL",
		}})
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		capsRaw, ok := resp.(map[string]interface{})["ConsumedCapacity"].([]interface{})
		if !ok || len(capsRaw) != 2 {
			t.Fatalf("iteration %d: expected two ConsumedCapacity entries, got %T %v", i, resp.(map[string]interface{})["ConsumedCapacity"], resp.(map[string]interface{})["ConsumedCapacity"])
		}
		first := capsRaw[0].(map[string]interface{})["TableName"]
		second := capsRaw[1].(map[string]interface{})["TableName"]
		if first != "BatchCapOrderA" || second != "BatchCapOrderB" {
			t.Fatalf("iteration %d: entries must answer in sorted order [BatchCapOrderA, BatchCapOrderB], got [%v, %v]", i, first, second)
		}
	}
}

func TestTransactAggregateSizeLimit(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnAggTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	// Each item is just under the 400 KB per-item cap; eleven of them
	// carry the request past the 4 MB aggregate.
	payload := strings.Repeat("x", int(dbstore.MaxItemSizeBytes)-64)
	putOp := func(id string) map[string]interface{} {
		return map[string]interface{}{"Put": map[string]interface{}{
			"TableName": "TxnAggTable",
			"Item":      map[string]interface{}{"id": sVal(id), "v": sVal(payload)},
		}}
	}

	var oversize []interface{}
	for i := 0; i < 11; i++ {
		oversize = append(oversize, putOp(fmt.Sprintf("w%d", i)))
	}
	_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": oversize,
	}})
	if err == nil || !strings.Contains(err.Error(), "The aggregate size of the items in the transaction exceeds 4 MB") {
		t.Fatalf("oversize write aggregate: expected the aggregate ValidationException, got %v", err)
	}
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnAggTable", "Key": map[string]interface{}{"id": sVal("w0")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	if _, present := resp.(map[string]interface{})["Item"]; present {
		t.Fatal("oversize write aggregate: expected no item after the rejected transaction")
	}

	// The read plane: seed eleven sub-cap items and request them all in
	// one TransactGetItems.
	for i := 0; i < 11; i++ {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "TxnAggTable",
			"Item":      map[string]interface{}{"id": sVal(fmt.Sprintf("g%d", i)), "v": sVal(payload)},
		}}); err != nil {
			t.Fatalf("seed put g%d: %v", i, err)
		}
	}
	var gets []interface{}
	for i := 0; i < 11; i++ {
		gets = append(gets, map[string]interface{}{"Get": map[string]interface{}{
			"TableName": "TxnAggTable",
			"Key":       map[string]interface{}{"id": sVal(fmt.Sprintf("g%d", i))},
		}})
	}
	_, err = svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": gets,
	}})
	if err == nil || !strings.Contains(err.Error(), "The aggregate size of the items in the transaction exceeded 4 MB") {
		t.Fatalf("oversize read aggregate: expected the aggregate ValidationException, got %v", err)
	}

	// A sub-limit aggregate on the write plane keeps committing.
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{putOp("ok")},
	}}); err != nil {
		t.Fatalf("single-item write: unexpected error %v", err)
	}
}

func TestBatchGetRejectsUnparseableKeys(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "BatchBadKeyTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "BatchBadKeyTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("x")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	// A key whose value is not an attribute-value shape rejects the whole
	// request, even when a well-formed key sits beside it.
	_, err := svc.BatchGetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"BatchBadKeyTable": map[string]interface{}{"Keys": []interface{}{
				map[string]interface{}{"id": sVal("a")},
				map[string]interface{}{"id": "bare-string"},
			}},
		},
	}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("malformed key: expected the whole-request ValidationException, got %v", err)
	}

	// The well-formed shape keeps serving, with no UnprocessedKeys echo.
	resp, err := svc.BatchGetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"BatchBadKeyTable": map[string]interface{}{"Keys": []interface{}{
				map[string]interface{}{"id": sVal("a")},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("valid key: unexpected error %v", err)
	}
	tableResp := resp.(map[string]interface{})["Responses"].(map[string]interface{})["BatchBadKeyTable"]
	if got := tableResp.([]map[string]interface{})[0]["v"].(map[string]interface{})["S"]; got != "x" {
		t.Fatalf("valid key: expected v=x, got %v", got)
	}
}

func TestTransactWriteExecutorValidationIsValidationException(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnExprTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	key := map[string]interface{}{"id": sVal("a")}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnExprTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("old")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}
	txnUpdate := func(expr string, values map[string]interface{}) error {
		params := map[string]interface{}{
			"TableName": "TxnExprTable", "Key": key,
			"UpdateExpression": expr,
		}
		if values != nil {
			params["ExpressionAttributeValues"] = values
		}
		_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Update": params}},
		}})
		return err
	}

	err := txnUpdate("SET id = :v", map[string]interface{}{":v": sVal("b")})
	if err == nil || !strings.Contains(err.Error(), "Cannot update attribute id") {
		t.Fatalf("key-attribute target: expected the key-schema ValidationException, got %v", err)
	}
	var canceled *TransactionCanceledError
	if errors.As(err, &canceled) {
		t.Fatal("key-attribute target: the rejection must not be a TransactionCanceledException")
	}

	err = txnUpdate("SET v = :missing", nil)
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("undefined value placeholder: expected the ValidationException, got %v", err)
	}
	if errors.As(err, &canceled) {
		t.Fatal("undefined value placeholder: the rejection must not be a TransactionCanceledException")
	}

	// The stored item is untouched by both rejected requests; a valid
	// update in the same position commits.
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnExprTable", "Key": key,
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	if got := resp.(map[string]interface{})["Item"].(map[string]interface{})["v"].(map[string]interface{})["S"]; got != "old" {
		t.Fatalf("stored item after rejections: expected v=old, got %v", got)
	}
	if err := txnUpdate("SET v = :ok", map[string]interface{}{":ok": sVal("new")}); err != nil {
		t.Fatalf("valid update: unexpected error %v", err)
	}
}

func TestTransactWriteItemsListBound(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)

	items := make([]interface{}, transactMaxItems+1)
	for i := range items {
		items[i] = map[string]interface{}{}
	}
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": items,
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("101-entry TransactItems: err = %v, want ErrInvalidParameter", err)
	}
}
