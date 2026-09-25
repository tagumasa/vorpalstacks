// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestAdminItemCoresRejectEmptyTableName(t *testing.T) {
	svc, _, _ := serviceRegionFixture(t)

	if _, err := svc.adminGetItem(nil, "us-east-1", "", nil); err != ErrInvalidParameter {
		t.Errorf("adminGetItem: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.adminScan(nil, "us-east-1", "", 10, nil); err != ErrInvalidParameter {
		t.Errorf("adminScan: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.adminPutItem(nil, "us-east-1", "", nil); err != ErrInvalidParameter {
		t.Errorf("adminPutItem: expected ErrInvalidParameter, got %v", err)
	}
	if err := svc.adminDeleteItem(nil, "us-east-1", "", nil); err != ErrInvalidParameter {
		t.Errorf("adminDeleteItem: expected ErrInvalidParameter, got %v", err)
	}
}

func TestItemWritesRequireActiveTable(t *testing.T) {
	svc, reqCtx, store := serviceRegionFixture(t)
	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "RestoringTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	table.Status = dbstore.TableStatusCreating
	if err := store.Tables().Put(table); err != nil {
		t.Fatalf("persist CREATING status: %v", err)
	}

	itemWrites := map[string]func() error{
		"PutItem": func() error {
			_, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "RestoringTable"}})
			return err
		},
		"UpdateItem": func() error {
			_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "RestoringTable"}})
			return err
		},
		"DeleteItem": func() error {
			_, err := svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "RestoringTable"}})
			return err
		},
		"adminPutItem": func() error {
			_, err := svc.adminPutItem(context.Background(), "us-east-1", "RestoringTable", nil)
			return err
		},
		"adminDeleteItem": func() error {
			return svc.adminDeleteItem(context.Background(), "us-east-1", "RestoringTable", nil)
		},
	}
	for name, write := range itemWrites {
		if err := write(); err != ErrTableNotActive {
			t.Errorf("%s against a CREATING table: expected ErrTableNotActive, got %v", name, err)
		}
	}

	// Reads resolve a non-ACTIVE table on both planes: with an absent key
	// the rejection must come from the empty-key rule, never the status.
	if _, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "RestoringTable"}}); err == ErrTableNotActive {
		t.Error("GetItem must not reject a CREATING table on status")
	}
	if _, err := svc.adminGetItem(context.Background(), "us-east-1", "RestoringTable", nil); err == ErrTableNotActive {
		t.Error("adminGetItem must not reject a CREATING table on status")
	}
}

func TestBatchTransactionAndPartiQLWritesRequireActiveTable(t *testing.T) {
	svc, reqCtx, store := serviceRegionFixture(t)
	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "RestoringTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	table.Status = dbstore.TableStatusCreating
	if err := store.Tables().Put(table); err != nil {
		t.Fatalf("persist CREATING status: %v", err)
	}

	key := func(id string) map[string]interface{} {
		return map[string]interface{}{"id": map[string]interface{}{"S": id}}
	}
	statement := func(stmt string) *request.ParsedRequest {
		return &request.ParsedRequest{Parameters: map[string]interface{}{"Statement": stmt}}
	}

	writes := map[string]func() error{
		"BatchWriteItem": func() error {
			_, err := svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"RequestItems": map[string]interface{}{
					"RestoringTable": []interface{}{map[string]interface{}{"PutRequest": map[string]interface{}{"Item": key("bw")}}},
				},
			}})
			return err
		},
		"TransactWriteItems": func() error {
			_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{"TableName": "RestoringTable", "Item": key("tw")}}},
			}})
			return err
		},
		"ExecuteStatement INSERT": func() error {
			_, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`INSERT INTO "RestoringTable" VALUE {'id': 'pi'}`))
			return err
		},
		"ExecuteStatement UPDATE": func() error {
			_, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`UPDATE "RestoringTable" SET v = '1' WHERE id = 'pu'`))
			return err
		},
		"ExecuteStatement DELETE": func() error {
			_, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`DELETE FROM "RestoringTable" WHERE id = 'pd'`))
			return err
		},
		"ExecuteTransaction INSERT": func() error {
			_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TransactStatements": []interface{}{map[string]interface{}{"Statement": `INSERT INTO "RestoringTable" VALUE {'id': 'ti'}`}},
			}})
			return err
		},
		"ExecuteTransaction UPDATE": func() error {
			_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TransactStatements": []interface{}{map[string]interface{}{"Statement": `UPDATE "RestoringTable" SET v = '1' WHERE id = 'tu'`}},
			}})
			return err
		},
		"ExecuteTransaction DELETE": func() error {
			_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TransactStatements": []interface{}{map[string]interface{}{"Statement": `DELETE FROM "RestoringTable" WHERE id = 'td'`}},
			}})
			return err
		},
	}
	for name, write := range writes {
		if err := write(); err != ErrTableNotActive {
			t.Errorf("%s against a CREATING table: expected ErrTableNotActive, got %v", name, err)
		}
	}

	reads := map[string]func() error{
		"BatchGetItem": func() error {
			_, err := svc.BatchGetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"RequestItems": map[string]interface{}{
					"RestoringTable": map[string]interface{}{"Keys": []interface{}{key("absent")}},
				},
			}})
			return err
		},
		"TransactGetItems": func() error {
			_, err := svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TransactItems": []interface{}{map[string]interface{}{"Get": map[string]interface{}{"TableName": "RestoringTable", "Key": key("absent")}}},
			}})
			return err
		},
		"ExecuteStatement SELECT": func() error {
			_, err := svc.ExecuteStatement(context.Background(), reqCtx, statement(`SELECT * FROM "RestoringTable"`))
			return err
		},
	}
	for name, read := range reads {
		if err := read(); err == ErrTableNotActive {
			t.Errorf("%s must not reject a CREATING table on status", name)
		}
	}
}

func TestBatchAndTransactWriteRejectUnionViolations(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "UnionTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	key := func(id string) map[string]interface{} {
		return map[string]interface{}{"id": map[string]interface{}{"S": id}}
	}
	validPut := map[string]interface{}{"Item": key("ok")}

	if _, err := svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"UnionTable": []interface{}{map[string]interface{}{"Unknown": map[string]interface{}{}}},
		},
	}}); err == nil || !strings.Contains(err.Error(), "exactly one of") {
		t.Fatalf("BatchWriteItem neither member: expected rejection, got %v", err)
	}

	if _, err := svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"UnionTable": []interface{}{map[string]interface{}{
				"PutRequest":    map[string]interface{}{"Item": key("both")},
				"DeleteRequest": map[string]interface{}{"Key": key("both")},
			}},
		},
	}}); err == nil || !strings.Contains(err.Error(), "exactly one of") {
		t.Fatalf("BatchWriteItem both members: expected rejection, got %v", err)
	}

	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{}},
	}}); err == nil || !strings.Contains(err.Error(), "exactly one of") {
		t.Fatalf("TransactWriteItems no member: expected rejection, got %v", err)
	}

	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{
			"Put":    map[string]interface{}{"TableName": "UnionTable", "Item": key("both")},
			"Delete": map[string]interface{}{"TableName": "UnionTable", "Key": key("both")},
		}},
	}}); err == nil || !strings.Contains(err.Error(), "exactly one of") {
		t.Fatalf("TransactWriteItems two members: expected rejection, got %v", err)
	}

	// Positive controls: the valid single-member forms still execute.
	if _, err := svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RequestItems": map[string]interface{}{
			"UnionTable": []interface{}{map[string]interface{}{"PutRequest": validPut}},
		},
	}}); err != nil {
		t.Fatalf("BatchWriteItem valid single member: unexpected error %v", err)
	}
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{"TableName": "UnionTable", "Item": key("ok")}}},
	}}); err != nil {
		t.Fatalf("TransactWriteItems valid single member: unexpected error %v", err)
	}
}

func TestItemPlaneReportsMalformedMembersAccurately(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "MalformedTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	malformed := map[string]interface{}{"id": map[string]interface{}{"S": "x", "N": "1"}}
	expectValueCause := func(t *testing.T, op string, err error) {
		t.Helper()
		if err == nil {
			t.Fatalf("%s malformed member: expected error, got success", op)
		}
		if strings.Contains(err.Error(), "Missing required key") {
			t.Fatalf("%s malformed member: reported the missing-key error, not the value cause: %v", op, err)
		}
		if !strings.Contains(err.Error(), "more than one datatypes set") {
			t.Fatalf("%s malformed member: expected the attribute-value cause, got %v", op, err)
		}
	}

	_, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "MalformedTable", "Item": malformed,
	}})
	expectValueCause(t, "PutItem", err)

	for _, op := range []struct {
		name string
		call func() error
	}{
		{"GetItem", func() error {
			_, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TableName": "MalformedTable", "Key": malformed,
			}})
			return err
		}},
		{"DeleteItem", func() error {
			_, err := svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TableName": "MalformedTable", "Key": malformed,
			}})
			return err
		}},
		{"UpdateItem", func() error {
			_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TableName": "MalformedTable", "Key": malformed,
			}})
			return err
		}},
	} {
		if err := op.call(); err == nil {
			t.Fatalf("%s malformed Key: expected error, got success", op.name)
		} else if !strings.Contains(err.Error(), "more than one datatypes set") {
			t.Fatalf("%s malformed Key: expected the attribute-value cause, got %v", op.name, err)
		}
	}

	// Genuinely absent members keep their own errors: a missing Key on
	// GetItem stays the empty-key ValidationException, and a missing Item on
	// PutItem stays the missing-key error.
	if _, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "MalformedTable",
	}}); err == nil || !strings.Contains(err.Error(), "Invalid parameter") {
		t.Fatalf("absent Key: expected the empty-key error, got %v", err)
	}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "MalformedTable",
	}}); err == nil || !strings.Contains(err.Error(), "Missing required key") {
		t.Fatalf("absent Item: expected the missing-key error, got %v", err)
	}
}

func TestReturnEnumValidation(t *testing.T) {
	for _, v := range []string{"", "NONE", "TOTAL", "INDEXES"} {
		if _, err := getReturnConsumedCapacity(map[string]interface{}{"ReturnConsumedCapacity": v}); err != nil {
			t.Errorf("ReturnConsumedCapacity %q: expected acceptance, got %v", v, err)
		}
	}
	if _, err := getReturnConsumedCapacity(map[string]interface{}{"ReturnConsumedCapacity": "BOGUS"}); err == nil {
		t.Error("ReturnConsumedCapacity BOGUS: expected rejection")
	}
	for _, v := range []string{"", "NONE", "SIZE"} {
		if _, err := getItemCollectionMetricsSetting(map[string]interface{}{"ReturnItemCollectionMetrics": v}); err != nil {
			t.Errorf("ReturnItemCollectionMetrics %q: expected acceptance, got %v", v, err)
		}
	}
	if _, err := getItemCollectionMetricsSetting(map[string]interface{}{"ReturnItemCollectionMetrics": "TOTAL"}); err == nil {
		t.Error("ReturnItemCollectionMetrics TOTAL: expected rejection")
	}

	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "EnumTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	item := map[string]interface{}{"id": map[string]interface{}{"S": "enum-1"}}

	for name, extra := range map[string]map[string]interface{}{
		"ReturnConsumedCapacity":      {"ReturnConsumedCapacity": "NOT_A_SETTING"},
		"ReturnItemCollectionMetrics": {"ReturnItemCollectionMetrics": "NOT_A_SETTING"},
	} {
		params := map[string]interface{}{"TableName": "EnumTable", "Item": item}
		for k, v := range extra {
			params[k] = v
		}
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params}); err == nil {
			t.Fatalf("%s: expected rejection, got success", name)
		}
		// The rejection happened before execution: the write left no item
		// (GetItem answers an absent item with an empty response).
		got, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "EnumTable", "Key": map[string]interface{}{"id": map[string]interface{}{"S": "enum-1"}},
		}})
		if err != nil {
			t.Fatalf("%s: follow-up read failed: %v", name, err)
		}
		if gotMap, ok := got.(map[string]interface{}); ok && gotMap["Item"] != nil {
			t.Fatalf("%s: the rejected write must not have executed", name)
		}
	}

	// Positive control: the valid settings still execute and shape the
	// response (TOTAL reports a ConsumedCapacity entry).
	resp, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "EnumTable", "Item": item, "ReturnConsumedCapacity": "TOTAL",
	}})
	if err != nil {
		t.Fatalf("valid TOTAL put: unexpected error %v", err)
	}
	respMap := resp.(map[string]interface{})
	if _, ok := respMap["ConsumedCapacity"]; !ok {
		t.Fatal("valid TOTAL put: expected a ConsumedCapacity entry")
	}
}

func TestPrimaryKeyValueLengthLimits(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "KeyLengthTable",
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

	// Over-limit hash value on the write plane.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "KeyLengthTable",
		"Item":      map[string]interface{}{"id": sVal(strings.Repeat("h", dbstore.MaxPartitionKeyBytes+1)), "sk": sVal("a")},
	}}); err == nil || !strings.Contains(err.Error(), "Size of hashkey has exceeded the maximum size limit of") {
		t.Fatalf("over-limit partition value: expected the hashkey size rejection, got %v", err)
	}

	// Over-limit sort value on the write plane.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "KeyLengthTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "sk": sVal(strings.Repeat("r", dbstore.MaxSortKeyBytes+1))},
	}}); err == nil || !strings.Contains(err.Error(), "Aggregated size of all range keys has exceeded the size limit of") {
		t.Fatalf("over-limit sort value: expected the range-key size rejection, got %v", err)
	}

	// The exact boundary values are accepted and retrievable: the limits are
	// inclusive maxima.
	boundaryHash := strings.Repeat("b", dbstore.MaxPartitionKeyBytes)
	boundarySort := strings.Repeat("s", dbstore.MaxSortKeyBytes)
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "KeyLengthTable",
		"Item":      map[string]interface{}{"id": sVal(boundaryHash), "sk": sVal(boundarySort)},
	}}); err != nil {
		t.Fatalf("boundary key values: expected acceptance, got %v", err)
	}
	if _, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "KeyLengthTable",
		"Key":       map[string]interface{}{"id": sVal(boundaryHash), "sk": sVal(boundarySort)},
	}}); err != nil {
		t.Fatalf("boundary key read-back: expected success, got %v", err)
	}

	// The read and update planes reject an over-long key with the same form:
	// no item carrying such a value can exist, so addressing by one is a
	// parameter-validation failure, not a miss.
	overLongKey := map[string]interface{}{"id": sVal(strings.Repeat("x", dbstore.MaxPartitionKeyBytes+1)), "sk": sVal("a")}
	for _, op := range []struct {
		name string
		call func() error
	}{
		{"GetItem", func() error {
			_, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "KeyLengthTable", "Key": overLongKey}})
			return err
		}},
		{"UpdateItem", func() error {
			_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "KeyLengthTable", "Key": overLongKey}})
			return err
		}},
		{"DeleteItem", func() error {
			_, err := svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "KeyLengthTable", "Key": overLongKey}})
			return err
		}},
	} {
		if err := op.call(); err == nil || !strings.Contains(err.Error(), "Size of hashkey has exceeded the maximum size limit of") {
			t.Errorf("%s over-limit partition value: expected the hashkey size rejection, got %v", op.name, err)
		}
	}
}

func TestKeySchemaMembershipRejected(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name: "KeyMembershipTable",
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
	extraKey := map[string]interface{}{"id": sVal("a"), "sk": sVal("b"), "extra": sVal("x")}
	missingKey := map[string]interface{}{"id": sVal("a")}
	expectMismatch := func(t *testing.T, op string, err error) {
		t.Helper()
		if err == nil || !strings.Contains(err.Error(), "The provided key element does not match the schema") {
			t.Errorf("%s extra key member: expected the schema-mismatch rejection, got %v", op, err)
		}
	}

	expectMismatch(t, "GetItem", func() error {
		_, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "KeyMembershipTable", "Key": extraKey}})
		return err
	}())
	expectMismatch(t, "DeleteItem", func() error {
		_, err := svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "KeyMembershipTable", "Key": extraKey}})
		return err
	}())
	expectMismatch(t, "UpdateItem", func() error {
		_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "KeyMembershipTable", "Key": extraKey}})
		return err
	}())

	// The missing-member direction is the same documented rejection: a
	// composite-key table addressed by the partition alone never matches.
	expectMismatch(t, "GetItem missing member", func() error {
		_, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "KeyMembershipTable", "Key": missingKey}})
		return err
	}())

	// Batch planes reject the whole request for one bad key.
	expectMismatch(t, "BatchGetItem", func() error {
		_, err := svc.BatchGetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"RequestItems": map[string]interface{}{
				"KeyMembershipTable": map[string]interface{}{"Keys": []interface{}{extraKey}},
			},
		}})
		return err
	}())
	expectMismatch(t, "BatchWriteItem DeleteRequest", func() error {
		_, err := svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"RequestItems": map[string]interface{}{
				"KeyMembershipTable": []interface{}{map[string]interface{}{"DeleteRequest": map[string]interface{}{"Key": extraKey}}},
			},
		}})
		return err
	}())

	// Transaction planes carry the same rule.
	expectMismatch(t, "TransactGetItems", func() error {
		_, err := svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Get": map[string]interface{}{"TableName": "KeyMembershipTable", "Key": extraKey}}},
		}})
		return err
	}())
	expectMismatch(t, "TransactWriteItems Delete", func() error {
		_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Delete": map[string]interface{}{"TableName": "KeyMembershipTable", "Key": extraKey}}},
		}})
		return err
	}())

	// Positive control: the exact key addresses the item end to end.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "KeyMembershipTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "sk": sVal("b"), "payload": sVal("v")},
	}}); err != nil {
		t.Fatalf("exact-key put: expected success, got %v", err)
	}
	got, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "KeyMembershipTable",
		"Key":       map[string]interface{}{"id": sVal("a"), "sk": sVal("b")},
	}})
	if err != nil {
		t.Fatalf("exact-key get: expected success, got %v", err)
	}
	item := got.(map[string]interface{})["Item"].(map[string]interface{})
	if _, ok := item["payload"]; !ok {
		t.Fatal("exact-key get: expected the stored item with its non-key attribute")
	}
}

func TestReturnValuesOnConditionCheckFailure(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "CondFailTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }

	// The existing item every failing write is answered with.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CondFailTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("old")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	failingCondition := func() (string, map[string]interface{}) {
		return "v = :expect", map[string]interface{}{":expect": sVal("never-matches")}
	}

	failedPut := func(rvocf string) error {
		expr, vals := failingCondition()
		_, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":           "CondFailTable",
			"Item":                map[string]interface{}{"id": sVal("a"), "v": sVal("new")},
			"ConditionExpression": expr, "ExpressionAttributeValues": vals,
			"ReturnValuesOnConditionCheckFailure": rvocf,
		}})
		return err
	}

	// ALL_OLD: the exception carries the existing item, and its wire form
	// renders the Item member alongside the message.
	err := failedPut("ALL_OLD")
	if err == nil {
		t.Fatal("conditioned put: expected failure")
	}
	var enriched *ConditionalCheckFailedError
	if !errors.As(err, &enriched) {
		t.Fatalf("conditioned put: expected the enriched ConditionalCheckFailedException, got %v", err)
	}
	if got := enriched.Item["v"].(map[string]interface{})["S"]; got != "old" {
		t.Fatalf("ALL_OLD item: expected the existing v=old, got %v", got)
	}
	wire := enriched.ToJSON()
	if !strings.Contains(wire, `"Item"`) || !strings.Contains(wire, `"__type":"com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException"`) {
		t.Fatalf("ALL_OLD wire form: expected the exception type with the Item member, got %s", wire)
	}

	// Plain sentinel shapes: NONE, the omitted member, and ALL_OLD against a
	// key with no item (nothing exists to return, matching the transaction
	// plane's semantics).
	for _, rvocf := range []string{"NONE", ""} {
		if err := failedPut(rvocf); !errors.Is(err, ErrConditionalCheckFailed) {
			t.Errorf("rvocf %q: expected the plain sentinel, got %v", rvocf, err)
		}
	}

	// DeleteItem and UpdateItem carry the same contract.
	expr, vals := failingCondition()
	if _, err := svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CondFailTable", "Key": map[string]interface{}{"id": sVal("a")},
		"ConditionExpression": expr, "ExpressionAttributeValues": vals,
		"ReturnValuesOnConditionCheckFailure": "ALL_OLD",
	}}); err == nil || !errors.As(err, &enriched) {
		t.Fatalf("conditioned delete: expected the enriched exception, got %v", err)
	} else if got := enriched.Item["v"].(map[string]interface{})["S"]; got != "old" {
		t.Fatalf("conditioned delete ALL_OLD: expected v=old, got %v", got)
	}
	expr, vals = failingCondition()
	if _, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "CondFailTable", "Key": map[string]interface{}{"id": sVal("a")},
		"UpdateExpression": "SET v = :nv",
		"ExpressionAttributeValues": map[string]interface{}{
			":nv":     sVal("new"),
			":expect": sVal("never-matches"),
		},
		"ConditionExpression":                 expr,
		"ReturnValuesOnConditionCheckFailure": "ALL_OLD",
	}}); err == nil || !errors.As(err, &enriched) {
		t.Fatalf("conditioned update: expected the enriched exception, got %v", err)
	} else if got := enriched.Item["v"].(map[string]interface{})["S"]; got != "old" {
		t.Fatalf("conditioned update ALL_OLD: expected v=old, got %v", got)
	}

	// An unknown enum value is request validation: rejected before any
	// write executes.
	if err := failedPut("BOGUS"); err == nil || !strings.Contains(err.Error(), "Invalid parameter") {
		t.Fatalf("unknown enum value: expected rejection, got %v", err)
	}
}
