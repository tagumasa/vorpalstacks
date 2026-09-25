// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestBatchExecuteStatementPerStatementResponses(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "BatchTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	statement := func(stmt string) *request.ParsedRequest {
		return &request.ParsedRequest{Parameters: map[string]interface{}{"Statement": stmt}}
	}
	for _, insert := range []string{
		`INSERT INTO "BatchTable" VALUE {'id': 'a'}`,
		`INSERT INTO "BatchTable" VALUE {'id': 'b'}`,
	} {
		if _, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(insert)); err != nil {
			t.Fatalf("seed %q: %v", insert, err)
		}
	}

	resp, err := svc.BatchExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statements": []interface{}{
			map[string]interface{}{"Statement": `SELECT * FROM "BatchTable" WHERE id = 'a'`},
			map[string]interface{}{"Statement": `SELECT * FROM "MissingTable" WHERE id = 'x'`},
		},
		"ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("batch: %v", err)
	}
	batchResp, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("batch response shape: %#v", resp)
	}
	responses, ok := batchResp["Responses"].([]map[string]interface{})
	if !ok || len(responses) != 2 {
		t.Fatalf("Responses: expected 2 slots, got %#v", batchResp["Responses"])
	}

	if responses[0]["TableName"] != "BatchTable" {
		t.Errorf("slot 0 TableName: %v", responses[0]["TableName"])
	}
	item, ok := responses[0]["Item"].(map[string]interface{})
	if !ok || item["id"].(map[string]interface{})["S"] != "a" {
		t.Errorf("slot 0 Item: %#v", responses[0]["Item"])
	}

	errSlot, ok := responses[1]["Error"].(map[string]interface{})
	if !ok || errSlot["Code"] != "ResourceNotFound" {
		t.Errorf("slot 1 Error: %#v", responses[1]["Error"])
	}
}

func TestPartiQLWriteCommitMaintainsGSIEntriesAndTableMetrics(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:      "PinTable",
		KeySchema: []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "gsik", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexes: []*dbstore.GlobalSecondaryIndex{{
			IndexName: "gsi1",
			KeySchema: []*dbstore.KeySchemaElement{{AttributeName: "gsik", KeyType: dbstore.KeyTypeHash}},
		}},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	statement := func(stmt string) *request.ParsedRequest {
		return &request.ParsedRequest{Parameters: map[string]interface{}{"Statement": stmt}}
	}
	gsiQuery := func(gsiKey string) int {
		t.Helper()
		resp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":                 "PinTable",
			"IndexName":                 "gsi1",
			"KeyConditionExpression":    "gsik = :g",
			"ExpressionAttributeValues": map[string]interface{}{":g": map[string]interface{}{"S": gsiKey}},
		}})
		if err != nil {
			t.Fatalf("gsi query %q: %v", gsiKey, err)
		}
		items, _ := resp.(map[string]interface{})["Items"].([]map[string]interface{})
		return len(items)
	}
	metrics := func() (int64, int64) {
		t.Helper()
		table, err := store.Tables().Get("PinTable")
		if err != nil {
			t.Fatalf("table get: %v", err)
		}
		return table.ItemCount, table.TableSizeBytes
	}

	for _, insert := range []string{
		`INSERT INTO "PinTable" VALUE {'id': 'a', 'gsik': 'g1', 'v': 'x'}`,
		`INSERT INTO "PinTable" VALUE {'id': 'b', 'gsik': 'g2'}`,
	} {
		if _, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(insert)); err != nil {
			t.Fatalf("seed %q: %v", insert, err)
		}
	}
	if got := gsiQuery("g1"); got != 1 {
		t.Fatalf("gsi g1 after inserts: expected 1 item, got %d", got)
	}

	if _, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`UPDATE "PinTable" SET gsik = 'g9' WHERE id = 'a'`)); err != nil {
		t.Fatalf("update: %v", err)
	}
	if got := gsiQuery("g1"); got != 0 {
		t.Fatalf("gsi g1 after update: old index entry not retired, got %d items", got)
	}
	if got := gsiQuery("g9"); got != 1 {
		t.Fatalf("gsi g9 after update: expected the moved item, got %d items", got)
	}

	count, size := metrics()
	if count != 2 {
		t.Fatalf("item count after inserts+update: expected 2, got %d", count)
	}
	itemA := dbstore.CalculateItemSize(map[string]*dbstore.AttributeValue{
		"id": dbstore.StringValue("a"), "gsik": dbstore.StringValue("g9"), "v": dbstore.StringValue("x"),
	})
	itemB := dbstore.CalculateItemSize(map[string]*dbstore.AttributeValue{
		"id": dbstore.StringValue("b"), "gsik": dbstore.StringValue("g2"),
	})
	if size != itemA+itemB {
		t.Fatalf("table size after update: expected %d, got %d", itemA+itemB, size)
	}

	if _, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`DELETE FROM "PinTable" WHERE id = 'a'`)); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := gsiQuery("g9"); got != 0 {
		t.Fatalf("gsi g9 after delete: expected 0 items, got %d", got)
	}
	count, size = metrics()
	if count != 1 || size != itemB {
		t.Fatalf("metrics after delete: expected count 1 size %d, got count %d size %d", itemB, count, size)
	}
}

func TestPartiQLUpdateDeleteZeroMatchContract(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "ZeroTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	statement := func(stmt string, extra map[string]interface{}) *request.ParsedRequest {
		params := map[string]interface{}{"Statement": stmt}
		for k, v := range extra {
			params[k] = v
		}
		return &request.ParsedRequest{Parameters: params}
	}
	allOld := map[string]interface{}{"ReturnValuesOnConditionCheckFailure": "ALL_OLD"}

	if _, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`INSERT INTO "ZeroTable" VALUE {'pk': 'a', 'sk': 's1', 'v': 'x'}`, nil)); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// A WHERE naming only the partition key does not resolve to a single
	// primary-key value.
	_, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`UPDATE "ZeroTable" SET v = 'y' WHERE pk = 'a'`, nil))
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "com.amazon.coral.validate#ValidationException" {
		t.Fatalf("partition-only UPDATE: expected ValidationException, got %v", err)
	}
	if _, err = svc.ExecuteStatement(context.Background(), reqCtx, statement(`DELETE FROM "ZeroTable" WHERE pk = 'a'`, nil)); err == nil {
		t.Fatal("partition-only DELETE: expected rejection")
	}

	// An UPDATE addressing an absent key fails the condition check.
	if _, err = svc.ExecuteStatement(context.Background(), reqCtx, statement(`UPDATE "ZeroTable" SET v = 'y' WHERE pk = 'a' AND sk = 'absent'`, nil)); !errors.Is(err, ErrConditionalCheckFailed) {
		t.Fatalf("absent-key UPDATE: expected ConditionalCheckFailedException, got %v", err)
	}

	// An UPDATE whose further predicates reject the stored item fails the
	// condition check; ALL_OLD carries the failing item's old image.
	_, err = svc.ExecuteStatement(context.Background(), reqCtx, statement(`UPDATE "ZeroTable" SET v = 'y' WHERE pk = 'a' AND sk = 's1' AND v = 'different'`, allOld))
	var enriched *ConditionalCheckFailedError
	if !errors.As(err, &enriched) {
		t.Fatalf("rejected-condition UPDATE: expected the enriched ConditionalCheckFailedException, got %v", err)
	}
	if got := enriched.Item["v"].(map[string]interface{})["S"]; got != "x" {
		t.Fatalf("rejected-condition UPDATE ALL_OLD: expected the old image v=x, got %v", got)
	}

	// A DELETE over an absent key succeeds with zero items deleted: the
	// model response surface is Items alone, carrying no Count member.
	resp, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`DELETE FROM "ZeroTable" WHERE pk = 'a' AND sk = 'absent'`, nil))
	if err != nil {
		t.Fatalf("absent-key DELETE: expected success, got %v", err)
	}
	delResp := resp.(map[string]interface{})
	if items := delResp["Items"].([]map[string]interface{}); len(items) != 0 {
		t.Fatalf("absent-key DELETE: expected zero items, got %v", items)
	}
	if _, ok := delResp["Count"]; ok {
		t.Fatalf("absent-key DELETE: response carries the off-model Count member: %v", delResp["Count"])
	}

	// A DELETE whose further predicates reject the stored item fails the
	// condition check, ALL_OLD carrying the old image.
	_, err = svc.ExecuteStatement(context.Background(), reqCtx, statement(`DELETE FROM "ZeroTable" WHERE pk = 'a' AND sk = 's1' AND v = 'nope'`, allOld))
	if !errors.As(err, &enriched) {
		t.Fatalf("rejected-condition DELETE: expected the enriched ConditionalCheckFailedException, got %v", err)
	}
	if got := enriched.Item["v"].(map[string]interface{})["S"]; got != "x" {
		t.Fatalf("rejected-condition DELETE ALL_OLD: expected the old image v=x, got %v", got)
	}

	// The item survives both rejected statements; a matching UPDATE applies.
	if _, err = svc.ExecuteStatement(context.Background(), reqCtx, statement(`UPDATE "ZeroTable" SET v = 'applied' WHERE pk = 'a' AND sk = 's1'`, nil)); err != nil {
		t.Fatalf("matching UPDATE: %v", err)
	}
	if _, err = svc.ExecuteStatement(context.Background(), reqCtx, statement(`DELETE FROM "ZeroTable" WHERE pk = 'a' AND sk = 's1'`, nil)); err != nil {
		t.Fatalf("matching DELETE: %v", err)
	}
	selResp, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`SELECT * FROM "ZeroTable" WHERE pk = 'a'`, nil))
	if err != nil {
		t.Fatalf("post-delete SELECT: %v", err)
	}
	if items := selResp.(map[string]interface{})["Items"].([]map[string]interface{}); len(items) != 0 {
		t.Fatalf("post-delete SELECT: expected no items, got %v", items)
	}
}

func TestExecuteTransactionZeroMatchStatement(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "TxnZeroTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	if _, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement": `INSERT INTO "TxnZeroTable" VALUE {'pk': 'a', 'sk': 's1', 'v': 'x'}`,
	}}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactStatements": []interface{}{map[string]interface{}{
			"Statement": `UPDATE "TxnZeroTable" SET v = 'y' WHERE pk = 'a' AND sk = 's1' AND v = 'nope'`,
		}},
	}})
	var canceled *TransactionCanceledError
	if !errors.As(err, &canceled) {
		t.Fatalf("zero-match UPDATE transaction: expected TransactionCanceledException, got %v", err)
	}
	if len(canceled.CancellationReasons) != 1 || canceled.CancellationReasons[0].Code != "ConditionalCheckFailed" {
		t.Fatalf("zero-match UPDATE transaction: expected a ConditionalCheckFailed reason, got %+v", canceled.CancellationReasons)
	}

	resp, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactStatements": []interface{}{map[string]interface{}{
			"Statement": `DELETE FROM "TxnZeroTable" WHERE pk = 'a' AND sk = 'absent'`,
		}},
	}})
	if err != nil {
		t.Fatalf("zero-key DELETE transaction: expected success, got %v", err)
	}
	responses := resp.(map[string]interface{})["Responses"].([]map[string]interface{})
	if len(responses) != 1 {
		t.Fatalf("zero-key DELETE transaction: expected one response slot, got %v", responses)
	}
	if item := responses[0]["Item"]; item != nil {
		t.Fatalf("zero-key DELETE transaction: expected an empty Item slot, got %v", item)
	}

	selResp, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement": `SELECT * FROM "TxnZeroTable" WHERE pk = 'a'`,
	}})
	if err != nil {
		t.Fatalf("post-transaction SELECT: %v", err)
	}
	items := selResp.(map[string]interface{})["Items"].([]map[string]interface{})
	if len(items) != 1 || items[0]["v"].(map[string]interface{})["S"] != "x" {
		t.Fatalf("post-transaction SELECT: expected the untouched item, got %v", items)
	}
}

func TestPartiQLInsertDuplicateKeyErrorCodes(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "DupTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	insert := func() error {
		_, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"Statement": `INSERT INTO "DupTable" VALUE {'id': 'a'}`,
		}})
		return err
	}
	if err := insert(); err != nil {
		t.Fatalf("first insert: %v", err)
	}

	err := insert()
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "com.amazonaws.dynamodb.v20120810#DuplicateItemException" {
		t.Fatalf("duplicate insert: expected DuplicateItemException, got %v", err)
	}

	_, err = svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactStatements": []interface{}{map[string]interface{}{
			"Statement": `INSERT INTO "DupTable" VALUE {'id': 'a'}`,
		}},
	}})
	var canceled *TransactionCanceledError
	if !errors.As(err, &canceled) {
		t.Fatalf("duplicate insert transaction: expected TransactionCanceledException, got %v", err)
	}
	if len(canceled.CancellationReasons) != 1 || canceled.CancellationReasons[0].Code != "ConditionalCheckFailed" {
		t.Fatalf("duplicate insert transaction: expected a ConditionalCheckFailed reason, got %+v", canceled.CancellationReasons)
	}
}

func TestBatchExecuteStatementEnforcesBatchRules(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "RuleTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	if _, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement": `INSERT INTO "RuleTable" VALUE {'pk': 'a', 'sk': 's1', 'v': 'x'}`,
	}}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// A mixed batch is rejected before anything executes.
	_, err := svc.BatchExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statements": []interface{}{
			map[string]interface{}{"Statement": `SELECT * FROM "RuleTable" WHERE pk = 'a' AND sk = 's1'`},
			map[string]interface{}{"Statement": `INSERT INTO "RuleTable" VALUE {'pk': 'a', 'sk': 's2'}`},
		},
	}})
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "com.amazon.coral.validate#ValidationException" {
		t.Fatalf("mixed batch: expected the request-level ValidationException, got %v", err)
	}
	if !strings.Contains(apiErr.Message, "cannot mix both in one batch") {
		t.Fatalf("mixed batch: expected the model's rule sentence, got %q", apiErr.Message)
	}

	// The rejected batch executed nothing: the INSERT's item is absent.
	selResp, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement": `SELECT * FROM "RuleTable" WHERE pk = 'a' AND sk = 's2'`,
	}})
	if err != nil {
		t.Fatalf("post-mix SELECT: %v", err)
	}
	if items := selResp.(map[string]interface{})["Items"].([]map[string]interface{}); len(items) != 0 {
		t.Fatalf("rejected mixed batch executed its write statement: %v", items)
	}

	// An all-read batch carrying a SELECT without all-key equality answers
	// that statement with a per-statement ValidationError; the conforming
	// statement of the same batch still returns its item.
	resp, err := svc.BatchExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statements": []interface{}{
			map[string]interface{}{"Statement": `SELECT * FROM "RuleTable" WHERE pk = 'a' AND sk = 's1'`},
			map[string]interface{}{"Statement": `SELECT * FROM "RuleTable" WHERE pk = 'a'`},
		},
	}})
	if err != nil {
		t.Fatalf("rule-violating read batch: unexpected request error: %v", err)
	}
	batchResp, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("read batch response shape: %#v", resp)
	}
	responses, ok := batchResp["Responses"].([]map[string]interface{})
	if !ok || len(responses) != 2 {
		t.Fatalf("read batch: expected 2 response slots, got %#v", batchResp["Responses"])
	}
	if item := responses[0]["Item"]; item == nil {
		t.Fatalf("conforming SELECT: expected its item, got %#v", responses[0])
	}
	errSlot, ok := responses[1]["Error"].(map[string]interface{})
	if !ok || errSlot["Code"] != "ValidationError" {
		t.Fatalf("violating SELECT: expected a per-statement ValidationError, got %#v", responses[1])
	}
	if !strings.Contains(errSlot["Message"].(string), "equality condition on all key attributes") {
		t.Fatalf("violating SELECT: expected the key-rule sentence, got %q", errSlot["Message"])
	}
}

func TestExecuteTransactionRejectsReadWriteMixAsValidation(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "MixTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactStatements": []interface{}{
			map[string]interface{}{"Statement": `SELECT * FROM "MixTable" WHERE id = 'a'`},
			map[string]interface{}{"Statement": `INSERT INTO "MixTable" VALUE {'id': 'a'}`},
		},
	}})
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "com.amazon.coral.validate#ValidationException" {
		t.Fatalf("mixed transaction: expected the request-level ValidationException, got %v", err)
	}
	if !strings.Contains(apiErr.Message, "cannot mix both in one transaction") {
		t.Fatalf("mixed transaction: expected the model's rule sentence, got %q", apiErr.Message)
	}

	// The rejected transaction executed nothing: the INSERT's item is absent.
	selResp, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement": `SELECT * FROM "MixTable" WHERE id = 'a'`,
	}})
	if err != nil {
		t.Fatalf("post-mix SELECT: %v", err)
	}
	if items := selResp.(map[string]interface{})["Items"].([]map[string]interface{}); len(items) != 0 {
		t.Fatalf("rejected mixed transaction executed its write statement: %v", items)
	}
}

func TestExecuteStatementResponseSurfaceContract(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "SurfaceTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	exec := func(params map[string]interface{}) (interface{}, error) {
		return svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
	}

	assertNoAccountingMembers := func(label string, resp interface{}) map[string]interface{} {
		t.Helper()
		m, ok := resp.(map[string]interface{})
		if !ok {
			t.Fatalf("%s: response shape: %#v", label, resp)
		}
		if _, ok := m["Count"]; ok {
			t.Fatalf("%s: response carries the off-model Count member: %v", label, m["Count"])
		}
		if _, ok := m["ScannedCount"]; ok {
			t.Fatalf("%s: response carries the off-model ScannedCount member: %v", label, m["ScannedCount"])
		}
		return m
	}

	if _, err := exec(map[string]interface{}{
		"Statement": `INSERT INTO "SurfaceTable" VALUE {'pk': 'a', 'v': 'x'}`,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	selResp, err := exec(map[string]interface{}{
		"Statement": `SELECT * FROM "SurfaceTable" WHERE pk = 'a'`,
	})
	if err != nil {
		t.Fatalf("SELECT: %v", err)
	}
	sel := assertNoAccountingMembers("SELECT", selResp)
	if items := sel["Items"].([]map[string]interface{}); len(items) != 1 {
		t.Fatalf("SELECT: expected the seeded item, got %v", sel["Items"])
	}

	insResp, err := exec(map[string]interface{}{
		"Statement": `INSERT INTO "SurfaceTable" VALUE {'pk': 'b'}`,
	})
	if err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	assertNoAccountingMembers("INSERT", insResp)

	updResp, err := exec(map[string]interface{}{
		"Statement": `UPDATE "SurfaceTable" SET v = 'y' WHERE pk = 'a'`,
	})
	if err != nil {
		t.Fatalf("UPDATE: %v", err)
	}
	assertNoAccountingMembers("UPDATE", updResp)

	// The DELETE zero-match tier answers SUCCESS with zero items and the
	// same clean surface.
	delResp, err := exec(map[string]interface{}{
		"Statement": `DELETE FROM "SurfaceTable" WHERE pk = 'missing'`,
	})
	if err != nil {
		t.Fatalf("zero-match DELETE: %v", err)
	}
	del := assertNoAccountingMembers("zero-match DELETE", delResp)
	if items := del["Items"].([]map[string]interface{}); len(items) != 0 {
		t.Fatalf("zero-match DELETE: expected zero items, got %v", del["Items"])
	}

	// An explicitly-present Limit below the PositiveIntegerObject range is
	// rejected as a parameter-validation failure, not read as "no limit".
	for _, bad := range []interface{}{0, -3} {
		_, err := exec(map[string]interface{}{
			"Statement": `SELECT * FROM "SurfaceTable"`,
			"Limit":     bad,
		})
		apiErr, ok := err.(*APIError)
		if !ok || apiErr.Code != "com.amazon.coral.validate#ValidationException" {
			t.Fatalf("Limit=%v: expected ValidationException, got %v", bad, err)
		}
	}

	// An in-range explicit Limit and an absent Limit both execute.
	for _, params := range []map[string]interface{}{
		{"Statement": `SELECT * FROM "SurfaceTable"`, "Limit": 1},
		{"Statement": `SELECT * FROM "SurfaceTable"`},
	} {
		if _, err := exec(params); err != nil {
			t.Fatalf("valid statement %v: %v", params["Limit"], err)
		}
	}

	// ResourceInUseException reports through the enum's InternalServerError
	// code with the exception's own message.
	code, msg := partiqlStatementErrorPair(NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table is being restored", http.StatusBadRequest), partiqlErrorBatch)
	if code != "InternalServerError" || msg != "Table is being restored" {
		t.Fatalf("ResourceInUseException mapping: expected InternalServerError with the message, got %q/%q", code, msg)
	}
}

func TestPartiQLWritePlaneValidationsMatchSingleItemPlane(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "SizeTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "pk", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	exec := func(stmt string, params []interface{}) error {
		p := map[string]interface{}{"Statement": stmt}
		if params != nil {
			p["Parameters"] = params
		}
		_, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: p})
		return err
	}
	execTxn := func(stmt string, params []interface{}) error {
		m := map[string]interface{}{"Statement": stmt}
		if params != nil {
			m["Parameters"] = params
		}
		_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactStatements": []interface{}{m},
		}})
		return err
	}

	// A wrong-typed primary key is rejected on both engines: the item's key
	// attribute types must match the schema definitions. The single-item
	// engine answers the mismatch directly; the transaction engine cancels
	// with the mismatch as the statement's ValidationError reason.
	if err := exec(`INSERT INTO "SizeTable" VALUE {'pk': 5, 'v': 'x'}`, nil); err == nil || !strings.Contains(err.Error(), "Type mismatch for key pk") {
		t.Fatalf("ExecuteStatement wrong-typed key INSERT: expected the key type mismatch, got %v", err)
	}
	err := execTxn(`INSERT INTO "SizeTable" VALUE {'pk': 5, 'v': 'x'}`, nil)
	var canceled *TransactionCanceledError
	if !errors.As(err, &canceled) || len(canceled.CancellationReasons) != 1 ||
		canceled.CancellationReasons[0].Code != "ValidationError" ||
		!strings.Contains(canceled.CancellationReasons[0].Message, "Type mismatch for key pk") {
		t.Fatalf("ExecuteTransaction wrong-typed key INSERT: expected a ValidationError cancellation carrying the mismatch, got %v", err)
	}

	// An INSERT whose composed item exceeds the documented item size is
	// rejected on both engines.
	big := strings.Repeat("x", 420*1024)
	if err := exec(`INSERT INTO "SizeTable" VALUE {'pk': 'big', 'v': ?}`, []interface{}{map[string]interface{}{"S": big}}); err == nil || !strings.Contains(err.Error(), "Invalid parameter") {
		t.Fatalf("ExecuteStatement oversized INSERT: expected rejection, got %v", err)
	}
	err = execTxn(`INSERT INTO "SizeTable" VALUE {'pk': 'big', 'v': ?}`, []interface{}{map[string]interface{}{"S": big}})
	if !errors.As(err, &canceled) || len(canceled.CancellationReasons) != 1 ||
		canceled.CancellationReasons[0].Code != "ValidationError" {
		t.Fatalf("ExecuteTransaction oversized INSERT: expected a ValidationError cancellation, got %v", err)
	}

	// A post-clause UPDATE image grown past the limit is rejected and the
	// stored item is left untouched.
	medium := strings.Repeat("y", 300*1024)
	if err := exec(`INSERT INTO "SizeTable" VALUE {'pk': 'grow', 'v': 'seed'}`, nil); err != nil {
		t.Fatalf("seed insert: %v", err)
	}
	if err := exec(`UPDATE "SizeTable" SET v = ? WHERE pk = 'grow'`, []interface{}{map[string]interface{}{"S": medium}}); err != nil {
		t.Fatalf("in-limit UPDATE: %v", err)
	}
	if err := exec(`UPDATE "SizeTable" SET w = ? WHERE pk = 'grow'`, []interface{}{map[string]interface{}{"S": medium}}); err == nil || !strings.Contains(err.Error(), "Invalid parameter") {
		t.Fatalf("oversized UPDATE: expected rejection, got %v", err)
	}
	item, err := store.Items().Get("SizeTable", map[string]*dbstore.AttributeValue{"pk": dbstore.StringValue("grow")})
	if err != nil {
		t.Fatalf("post-rejection read: %v", err)
	}
	if _, has := item.Attributes["w"]; has {
		t.Fatalf("oversized UPDATE: the rejected statement left its attribute behind")
	}
}

func TestExecuteTransactionReadShapeAndTableResolution(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "TxnReadTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	for _, pk := range []string{"a", "b"} {
		if _, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"Statement": `INSERT INTO "TxnReadTable" VALUE {'pk': '` + pk + `', 'v': 'shared'}`,
		}}); err != nil {
			t.Fatalf("seed %s: %v", pk, err)
		}
	}

	// A single-match read still answers its item.
	resp, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactStatements": []interface{}{map[string]interface{}{
			"Statement": `SELECT * FROM "TxnReadTable" WHERE pk = 'a' AND v = 'shared'`,
		}},
	}})
	if err != nil {
		t.Fatalf("single-match read: %v", err)
	}
	responses := resp.(map[string]interface{})["Responses"].([]map[string]interface{})
	if len(responses) != 1 || responses[0]["Item"] == nil {
		t.Fatalf("single-match read response = %v, want one non-nil Item", responses)
	}

	// A multi-match read is rejected: the shape cannot carry the matches.
	_, err = svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactStatements": []interface{}{map[string]interface{}{
			"Statement": `SELECT * FROM "TxnReadTable" WHERE v = 'shared'`,
		}},
	}})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !strings.Contains(apiErr.Code, "ValidationException") {
		t.Fatalf("multi-match read: expected ValidationException, got %v", err)
	}

	// A read or write statement naming a missing table answers
	// ResourceNotFoundException directly.
	for _, stmt := range []string{
		`SELECT * FROM "NoTxnTable" WHERE pk = 'a'`,
		`INSERT INTO "NoTxnTable" VALUE {'pk': 'a'}`,
	} {
		if _, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactStatements": []interface{}{map[string]interface{}{"Statement": stmt}},
		}}); !errors.Is(err, ErrResourceNotFound) {
			t.Fatalf("statement %q on a missing table: err = %v, want ErrResourceNotFound", stmt, err)
		}
	}
}
