// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestSizeGranularCapacityAccounting(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "CapTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	big := sVal(strings.Repeat("x", 5120))
	key := func(sk string) map[string]interface{} {
		return map[string]interface{}{"id": sVal("a"), "sk": sVal(sk)}
	}
	capacityOf := func(resp interface{}) float64 {
		t.Helper()
		cc, ok := resp.(map[string]interface{})["ConsumedCapacity"]
		if !ok {
			t.Fatal("expected a ConsumedCapacity member")
		}
		switch entry := cc.(type) {
		case map[string]interface{}:
			return entry["CapacityUnits"].(float64)
		case []interface{}:
			return entry[0].(map[string]interface{})["CapacityUnits"].(float64)
		}
		t.Fatalf("unexpected ConsumedCapacity shape: %T", cc)
		return 0
	}

	// GetItem: two units strongly consistent, one eventually consistent,
	// and a projection does not shrink the charge below the read size.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "sk": sVal("big"), "v": big},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapTable", "Key": key("big"), "ConsistentRead": true,
		"ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("strong get: %v", err)
	}
	if got := capacityOf(resp); got != 2.0 {
		t.Fatalf("strong get: expected 2.0 units for a 5128-byte item, got %v", got)
	}
	resp, err = svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapTable", "Key": key("big"),
		"ProjectionExpression":     "v",
		"ReturnConsumedCapacity":   "TOTAL",
		"ExpressionAttributeNames": nil,
	}})
	if err != nil {
		t.Fatalf("eventual get: %v", err)
	}
	if got := capacityOf(resp); got != 1.0 {
		t.Fatalf("eventual get: expected 1.0 unit for a 5128-byte item, got %v", got)
	}

	// BatchGetItem: the same two units, strongly consistent.
	resp, err = svc.BatchGetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"CapTable": map[string]interface{}{"Keys": []interface{}{key("big")}, "ConsistentRead": true},
		},
		"ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("batch get: %v", err)
	}
	if got := capacityOf(resp); got != 2.0 {
		t.Fatalf("batch get: expected 2.0 units, got %v", got)
	}

	// Query and Scan sum the evaluated items' units; both default to
	// eventually consistent reads, so two 5128-byte items charge
	// (2 + 2) / 2.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "sk": sVal("b2"), "v": big},
	}}); err != nil {
		t.Fatalf("seed put b2: %v", err)
	}
	resp, err = svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "CapTable",
		"KeyConditionExpression":    "id = :id",
		"ExpressionAttributeValues": map[string]interface{}{":id": sVal("a")},
		"ReturnConsumedCapacity":    "TOTAL",
	}})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if got := capacityOf(resp); got != 2.0 {
		t.Fatalf("query: expected 2.0 units for two 5128-byte items read eventually consistently, got %v", got)
	}
	resp, err = svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "CapTable",
		"ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if got := capacityOf(resp); got != 2.0 {
		t.Fatalf("scan: expected 2.0 units for two 5128-byte items read eventually consistently, got %v", got)
	}

	// The write planes charge the item's size rounded up to 1 KB
	// multiples: a 5128-byte item costs six units, put or updated, and its
	// deletion six as well.
	resp, err = svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "CapTable",
		"Item":                   map[string]interface{}{"id": sVal("a"), "sk": sVal("w"), "v": big},
		"ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	if got := capacityOf(resp); got != 6.0 {
		t.Fatalf("put: expected 6.0 units for a 5128-byte item, got %v", got)
	}
	resp, err = svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapTable", "Key": key("w"),
		"UpdateExpression":          "SET v = :v",
		"ExpressionAttributeValues": map[string]interface{}{":v": big},
		"ReturnConsumedCapacity":    "TOTAL",
	}})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if got := capacityOf(resp); got != 6.0 {
		t.Fatalf("update: expected 6.0 units for a 5128-byte item, got %v", got)
	}
	resp, err = svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapTable", "Key": key("w"),
		"ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := capacityOf(resp); got != 6.0 {
		t.Fatalf("delete: expected 6.0 units for a 5128-byte item, got %v", got)
	}

	// BatchWriteItem charges the same per-item write units.
	resp, err = svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"CapTable": []interface{}{map[string]interface{}{"PutRequest": map[string]interface{}{
				"Item": map[string]interface{}{"id": sVal("a"), "sk": sVal("bw"), "v": big},
			}}},
		},
		"ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("batch write: %v", err)
	}
	if got := capacityOf(resp); got != 6.0 {
		t.Fatalf("batch write: expected 6.0 units for a 5128-byte item, got %v", got)
	}
}

func TestConsumedCapacityShapesPerPlane(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "CapShapeTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	key := map[string]interface{}{"id": sVal("a"), "sk": sVal("b")}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapShapeTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "sk": sVal("b"), "v": sVal("x")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}
	assertMemberSet := func(label string, entry map[string]interface{}, want ...string) {
		t.Helper()
		wantSet := map[string]bool{}
		for _, w := range want {
			wantSet[w] = true
		}
		if len(entry) != len(wantSet) {
			t.Fatalf("%s: expected members %v, got %v", label, want, entry)
		}
		for _, w := range want {
			if _, ok := entry[w]; !ok {
				t.Fatalf("%s: expected member %q, got %v", label, w, entry)
			}
		}
	}

	// GetItem: TOTAL aggregate only; INDEXES adds the table detail.
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapShapeTable", "Key": key, "ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("get TOTAL: %v", err)
	}
	assertMemberSet("get TOTAL", resp.(map[string]interface{})["ConsumedCapacity"].(map[string]interface{}), "TableName", "CapacityUnits")
	resp, err = svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CapShapeTable", "Key": key, "ReturnConsumedCapacity": "INDEXES",
	}})
	if err != nil {
		t.Fatalf("get INDEXES: %v", err)
	}
	getIndexes := resp.(map[string]interface{})["ConsumedCapacity"].(map[string]interface{})
	assertMemberSet("get INDEXES", getIndexes, "TableName", "CapacityUnits", "Table")

	// BatchGetItem answers one entry per table with the same shape rules.
	resp, err = svc.BatchGetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"CapShapeTable": map[string]interface{}{"Keys": []interface{}{key}},
		},
		"ReturnConsumedCapacity": "INDEXES",
	}})
	if err != nil {
		t.Fatalf("batch get INDEXES: %v", err)
	}
	batchEntry := resp.(map[string]interface{})["ConsumedCapacity"].([]interface{})[0].(map[string]interface{})
	assertMemberSet("batch get INDEXES", batchEntry, "TableName", "CapacityUnits", "Table")

	// Query and Scan: TOTAL aggregate only; INDEXES adds the table detail.
	query := func(setting string) map[string]interface{} {
		t.Helper()
		resp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":                 "CapShapeTable",
			"KeyConditionExpression":    "id = :id",
			"ExpressionAttributeValues": map[string]interface{}{":id": sVal("a")},
			"ReturnConsumedCapacity":    setting,
		}})
		if err != nil {
			t.Fatalf("query %s: %v", setting, err)
		}
		return resp.(map[string]interface{})["ConsumedCapacity"].(map[string]interface{})
	}
	assertMemberSet("query TOTAL", query("TOTAL"), "TableName", "CapacityUnits")
	assertMemberSet("query INDEXES", query("INDEXES"), "TableName", "CapacityUnits", "Table")
	scan := func(setting string) map[string]interface{} {
		t.Helper()
		resp, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "CapShapeTable", "ReturnConsumedCapacity": setting,
		}})
		if err != nil {
			t.Fatalf("scan %s: %v", setting, err)
		}
		return resp.(map[string]interface{})["ConsumedCapacity"].(map[string]interface{})
	}
	assertMemberSet("scan TOTAL", scan("TOTAL"), "TableName", "CapacityUnits")
	assertMemberSet("scan INDEXES", scan("INDEXES"), "TableName", "CapacityUnits", "Table")

	// The transactional read keeps its read-specific member on both
	// settings, adding the table detail under INDEXES.
	tgi := func(setting string) map[string]interface{} {
		t.Helper()
		resp, err := svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Get": map[string]interface{}{
				"TableName": "CapShapeTable", "Key": key,
			}}},
			"ReturnConsumedCapacity": setting,
		}})
		if err != nil {
			t.Fatalf("transact get %s: %v", setting, err)
		}
		return resp.(map[string]interface{})["ConsumedCapacity"].([]map[string]interface{})[0]
	}
	assertMemberSet("transact get TOTAL", tgi("TOTAL"), "TableName", "CapacityUnits", "ReadCapacityUnits")
	assertMemberSet("transact get INDEXES", tgi("INDEXES"), "TableName", "CapacityUnits", "ReadCapacityUnits", "Table")
}

// TestExecuteTransactionPerStatementCapacity pins the batch
// transaction's ConsumedCapacity: each write statement charges the
// transactional two-unit write minimum and each read statement its
// transactional read charge, aggregated per table in statement
// first-appearance order — the statement count shows in the units, the
// entry order in the statement order.
func TestExecuteTransactionPerStatementCapacity(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "ExecTxnCapWriteTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "ExecTxnCapReadTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create read table: %v", err)
	}
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "ExecTxnCapReadTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("x")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	exec := func(statements ...string) []map[string]interface{} {
		t.Helper()
		stmts := make([]interface{}, len(statements))
		for i, s := range statements {
			stmts[i] = map[string]interface{}{"Statement": s}
		}
		resp, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactStatements":     stmts,
			"ReturnConsumedCapacity": "TOTAL",
		}})
		if err != nil {
			t.Fatalf("ExecuteTransaction: %v", err)
		}
		cc, ok := resp.(map[string]interface{})["ConsumedCapacity"].([]map[string]interface{})
		if !ok {
			t.Fatalf("expected ConsumedCapacity entries, got %#v", resp)
		}
		return cc
	}

	one := exec(`INSERT INTO "ExecTxnCapWriteTable" VALUE {'id': 'a'}`)
	if len(one) != 1 || one[0]["TableName"] != "ExecTxnCapWriteTable" || one[0]["CapacityUnits"].(float64) != 2 {
		t.Fatalf("single write: expected 2 units on the write table, got %#v", one)
	}

	// A transaction is all-read or all-write; each plane charges its own
	// statement kind. Two writes on one table aggregate to four units,
	// and the entry order follows statement first-appearance order.
	writes := exec(
		`INSERT INTO "ExecTxnCapWriteTable" VALUE {'id': 'b'}`,
		`INSERT INTO "ExecTxnCapReadTable" VALUE {'id': 'z'}`,
		`INSERT INTO "ExecTxnCapWriteTable" VALUE {'id': 'c'}`,
	)
	if len(writes) != 2 {
		t.Fatalf("write transaction: expected 2 capacity entries, got %#v", writes)
	}
	if writes[0]["TableName"] != "ExecTxnCapWriteTable" || writes[0]["CapacityUnits"].(float64) != 4 || writes[0]["WriteCapacityUnits"].(float64) != 4 {
		t.Fatalf("write table entry: expected 4 write units, got %#v", writes[0])
	}
	if _, has := writes[0]["ReadCapacityUnits"]; has {
		t.Fatalf("write table entry: unexpected ReadCapacityUnits, got %#v", writes[0])
	}
	if writes[1]["TableName"] != "ExecTxnCapReadTable" || writes[1]["CapacityUnits"].(float64) != 2 || writes[1]["WriteCapacityUnits"].(float64) != 2 {
		t.Fatalf("second write table entry: expected 2 write units, got %#v", writes[1])
	}

	reads := exec(
		`SELECT * FROM "ExecTxnCapWriteTable" WHERE id = 'a'`,
		`SELECT * FROM "ExecTxnCapReadTable" WHERE id = 'a'`,
	)
	if len(reads) != 2 {
		t.Fatalf("read transaction: expected 2 capacity entries, got %#v", reads)
	}
	if reads[0]["TableName"] != "ExecTxnCapWriteTable" || reads[0]["CapacityUnits"].(float64) != 2 || reads[0]["ReadCapacityUnits"].(float64) != 2 {
		t.Fatalf("read table entry: expected 2 read units, got %#v", reads[0])
	}
	if _, has := reads[0]["WriteCapacityUnits"]; has {
		t.Fatalf("read table entry: unexpected WriteCapacityUnits, got %#v", reads[0])
	}
}

func TestBatchExecuteStatementConsumedCapacityAggregation(t *testing.T) {
	svc, reqCtx, store := serviceRegionFixture(t)
	for _, name := range []string{"CapTableA", "CapTableB"} {
		if _, err := store.Tables().Create(dbstore.CreateTableParams{
			Name:                 name,
			KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
			AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
			BillingMode:          dbstore.BillingModePayPerRequest,
		}); err != nil {
			t.Fatalf("create table %s: %v", name, err)
		}
	}

	batch := func(rcc string, statements ...string) map[string]interface{} {
		stmts := make([]interface{}, 0, len(statements))
		for _, s := range statements {
			stmts = append(stmts, map[string]interface{}{"Statement": s})
		}
		params := map[string]interface{}{"Statements": stmts}
		if rcc != "" {
			params["ReturnConsumedCapacity"] = rcc
		}
		resp, err := svc.BatchExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		if err != nil {
			t.Fatalf("batch (%s): %v", rcc, err)
		}
		return resp.(map[string]interface{})
	}

	// Seed the three items through a write batch (its own capacity entry is
	// not under assertion here).
	batch("TOTAL",
		`INSERT INTO "CapTableA" VALUE {'id': 'a1'}`,
		`INSERT INTO "CapTableA" VALUE {'id': 'a2'}`,
		`INSERT INTO "CapTableB" VALUE {'id': 'b1'}`,
	)

	// A read batch over two tables: the eventual read unit is 0.5 per
	// statement, so table A sums two statements to 1.0 and table B carries
	// its single statement's 0.5 — two entries, one per table.
	readResp := batch("TOTAL",
		`SELECT * FROM "CapTableA" WHERE id = 'a1'`,
		`SELECT * FROM "CapTableA" WHERE id = 'a2'`,
		`SELECT * FROM "CapTableB" WHERE id = 'b1'`,
	)
	ccList, ok := readResp["ConsumedCapacity"].([]map[string]interface{})
	if !ok || len(ccList) != 2 {
		t.Fatalf("read batch: expected 2 aggregated capacity entries, got %#v", readResp["ConsumedCapacity"])
	}
	if ccList[0]["TableName"] != "CapTableA" || ccList[0]["CapacityUnits"] != 1.0 {
		t.Fatalf("read batch table A: expected summed 1.0, got %#v", ccList[0])
	}
	if ccList[1]["TableName"] != "CapTableB" || ccList[1]["CapacityUnits"] != 0.5 {
		t.Fatalf("read batch table B: expected 0.5, got %#v", ccList[1])
	}

	// A write batch on one table: each write consumes 1.0, summed into a
	// single entry.
	writeResp := batch("TOTAL",
		`INSERT INTO "CapTableA" VALUE {'id': 'w1'}`,
		`INSERT INTO "CapTableA" VALUE {'id': 'w2'}`,
	)
	ccList, ok = writeResp["ConsumedCapacity"].([]map[string]interface{})
	if !ok || len(ccList) != 1 {
		t.Fatalf("write batch: expected 1 aggregated capacity entry, got %#v", writeResp["ConsumedCapacity"])
	}
	if ccList[0]["TableName"] != "CapTableA" || ccList[0]["CapacityUnits"] != 2.0 {
		t.Fatalf("write batch: expected summed 2.0 on CapTableA, got %#v", ccList[0])
	}

	// Without the setting the batch reports no capacity member.
	noneResp := batch("", `SELECT * FROM "CapTableA" WHERE id = 'a1'`)
	if _, has := noneResp["ConsumedCapacity"]; has {
		t.Fatalf("unset ReturnConsumedCapacity: expected no capacity member, got %#v", noneResp["ConsumedCapacity"])
	}
}
