package dynamodb

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The Query/Scan data-plane request validation family: parallel-scan
// segment bounds, page Limit handling, key-condition grammar and typing,
// and the exclusive-start-key contract. Each test drives the wire-plane
// method the SDK client itself reaches.

func TestParallelScanRejectsSegmentBeyondTotalSegments(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"Item": map[string]interface{}{
			"id": map[string]interface{}{"S": "A"},
			"sk": map[string]interface{}{"S": "1"},
		},
	}}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// A segment index below the count is a legal segment, whichever
	// segment the seeded item happens to hash into.
	if _, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":     "LegacyTable",
		"Segment":       2,
		"TotalSegments": 3,
	}}); err != nil {
		t.Fatalf("segment 2 of 3: unexpected error: %v", err)
	}

	// Segment == TotalSegments and Segment > TotalSegments both address no
	// segment and must be rejected, not silently return an empty page.
	for _, tc := range []struct{ segment, total int }{
		{3, 3},
		{5, 3},
	} {
		_, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":     "LegacyTable",
			"Segment":       tc.segment,
			"TotalSegments": tc.total,
		}})
		if err == nil || !strings.Contains(err.Error(), "Invalid parameter") {
			t.Fatalf("segment %d of %d: expected Invalid parameter rejection, got %v", tc.segment, tc.total, err)
		}
	}
}

func TestScanHonoursLimitAboveFormerThousandCap(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	// Seed 1200 single-attribute items so a page above the former cap is
	// observable: a Scan with Limit 1100 must return 1100 items and a
	// continuation key, not a silently truncated 1000-item page.
	for i := 0; i < 1200; i++ {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "LegacyTable",
			"Item": map[string]interface{}{
				"id": map[string]interface{}{"S": fmt.Sprintf("item%04d", i)},
				"sk": map[string]interface{}{"S": "1"},
			},
		}}); err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
	}

	scanResp, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"Limit":     1100,
	}})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	resp := scanResp.(map[string]interface{})
	if count := resp["Count"].(int); count != 1100 {
		t.Fatalf("expected Count 1100 (the requested page size), got %v", count)
	}
	if scanned := resp["ScannedCount"].(int); scanned != 1100 {
		t.Fatalf("expected ScannedCount 1100, got %v", scanned)
	}
	if len(resp["Items"].([]map[string]interface{})) != 1100 {
		t.Fatalf("expected 1100 items in the page, got %d", len(resp["Items"].([]map[string]interface{})))
	}
	if _, hasLEK := resp["LastEvaluatedKey"]; !hasLEK {
		t.Fatal("expected LastEvaluatedKey on a truncated page of 1100 of 1200 items")
	}

	// Query shares the preamble: the same Limit above the former cap reads
	// the full partition page.
	queryResp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "LegacyTable",
		"KeyConditionExpression":    "id = :id",
		"ExpressionAttributeValues": map[string]interface{}{":id": map[string]interface{}{"S": "item0001"}},
		"Limit":                     1100,
	}})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	// One matching item exists; the page must not be capped below the
	// requested limit, and the single item must be returned.
	if items := queryResp.(map[string]interface{})["Items"].([]map[string]interface{}); len(items) != 1 {
		t.Fatalf("expected the single matching item, got %d", len(items))
	}
}

func TestQueryRejectsInvalidKeyConditions(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	for _, sk := range []string{"1", "2", "3"} {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "LegacyTable",
			"Item": map[string]interface{}{
				"id": map[string]interface{}{"S": "A"},
				"sk": map[string]interface{}{"S": sk},
				"v":  map[string]interface{}{"N": sk},
			},
		}}); err != nil {
			t.Fatalf("seed %s: %v", sk, err)
		}
	}

	base := map[string]interface{}{
		"TableName": "LegacyTable",
		"ExpressionAttributeValues": map[string]interface{}{
			":a":  map[string]interface{}{"S": "A"},
			":p":  map[string]interface{}{"S": "2"},
			":lo": map[string]interface{}{"S": "2"},
			":hi": map[string]interface{}{"S": "3"},
		},
	}

	// The documented grammar rejects each of these shapes with the same
	// ValidationException instead of widening or silently emptying the
	// read.
	rejections := []struct{ name, expr string }{
		{"unresolvable sort value drops the condition (was: whole partition)", "id = :a AND sk > :typo"},
		{"unsupported sort operator (was: empty page)", "id = :a AND sk CONTAINS :p"},
		{"OR combination", "id = :a OR sk > :p"},
		{"duplicate partition condition", "id = :a AND id = :a"},
		{"duplicate sort condition", "id = :a AND sk > :p AND sk < :hi"},
		{"begins_with on the partition key", "begins_with(id, :a)"},
		{"uppercase function name (documented case-sensitive)", "id = :a AND BEGINS_WITH(sk, :p)"},
		{"infix begins_with operator (function form alone)", "id = :a AND sk begins_with :lo"},
		{"non-key attribute", "id = :a AND v = :p"},
		{"partition key with non-equality operator", "id > :a AND sk > :p"},
	}
	for _, tc := range rejections {
		params := map[string]interface{}{}
		for k, v := range base {
			params[k] = v
		}
		params["KeyConditionExpression"] = tc.expr
		_, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		if err == nil || !strings.Contains(err.Error(), "Query key condition not supported") {
			t.Fatalf("%s: expected the documented rejection, got %v", tc.name, err)
		}
	}

	// Positive controls: the documented forms still answer.
	positives := []struct {
		name  string
		expr  string
		count int
	}{
		{"partition equality alone", "id = :a", 3},
		{"sort equality", "id = :a AND sk = :p", 1},
		{"lowercase begins_with function form", "id = :a AND begins_with(sk, :lo)", 1},
		{"BETWEEN", "id = :a AND sk BETWEEN :lo AND :hi", 2},
	}
	for _, tc := range positives {
		params := map[string]interface{}{}
		for k, v := range base {
			params[k] = v
		}
		params["KeyConditionExpression"] = tc.expr
		resp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}
		if count := resp.(map[string]interface{})["Count"].(int); count != tc.count {
			t.Fatalf("%s: expected Count %d, got %v", tc.name, tc.count, count)
		}
	}
}

// newTypedKeyTestService creates a table whose sort key and GSI hash are
// Number-typed, so wrong-typed key-condition values are observable on both
// the base table and an index query.
func newTypedKeyTestService(t *testing.T) (*DynamoDBService, *request.RequestContext) {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name: "TypedTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "nsk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "nsk", AttributeType: dbstore.ScalarAttributeTypeN},
			{AttributeName: "gid", AttributeType: dbstore.ScalarAttributeTypeN},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []*dbstore.GlobalSecondaryIndex{
			{
				IndexName: "gsi1",
				KeySchema: []*dbstore.KeySchemaElement{
					{AttributeName: "gid", KeyType: dbstore.KeyTypeHash},
				},
				Projection: &dbstore.Projection{ProjectionType: "ALL"},
			},
		},
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return svc, request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
}

func TestQueryRejectsWrongTypedKeyConditionValues(t *testing.T) {
	svc, reqCtx := newTypedKeyTestService(t)
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TypedTable",
		"Item": map[string]interface{}{
			"id":  map[string]interface{}{"S": "A"},
			"nsk": map[string]interface{}{"N": "1"},
			"gid": map[string]interface{}{"N": "7"},
		},
	}}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// The table types id as S: an N-typed equality value would encode by
	// its own type and silently match nothing.
	_, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "TypedTable",
		"KeyConditionExpression":    "id = :v",
		"ExpressionAttributeValues": map[string]interface{}{":v": map[string]interface{}{"N": "5"}},
	}})
	if err == nil || !strings.Contains(err.Error(), "Condition parameter type does not match schema type") {
		t.Fatalf("N value against S partition key: expected the documented type-mismatch rejection, got %v", err)
	}

	// The sort key is N-typed: an S-typed comparison value filters every
	// item out silently.
	_, err = svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "TypedTable",
		"KeyConditionExpression": "id = :a AND nsk > :lo",
		"ExpressionAttributeValues": map[string]interface{}{
			":a":  map[string]interface{}{"S": "A"},
			":lo": map[string]interface{}{"S": "0"},
		},
	}})
	if err == nil || !strings.Contains(err.Error(), "Condition parameter type does not match schema type") {
		t.Fatalf("S value against N sort key: expected the documented type-mismatch rejection, got %v", err)
	}

	// BETWEEN's upper bound is checked too.
	_, err = svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "TypedTable",
		"KeyConditionExpression": "id = :a AND nsk BETWEEN :lo AND :hi",
		"ExpressionAttributeValues": map[string]interface{}{
			":a":  map[string]interface{}{"S": "A"},
			":lo": map[string]interface{}{"N": "0"},
			":hi": map[string]interface{}{"S": "9"},
		},
	}})
	if err == nil || !strings.Contains(err.Error(), "Condition parameter type does not match schema type") {
		t.Fatalf("BETWEEN with mixed bound types: expected the type-mismatch rejection, got %v", err)
	}

	// begins_with on a Number-typed sort key is documented as unusable.
	_, err = svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "TypedTable",
		"KeyConditionExpression": "id = :a AND begins_with(nsk, :p)",
		"ExpressionAttributeValues": map[string]interface{}{
			":a": map[string]interface{}{"S": "A"},
			":p": map[string]interface{}{"N": "1"},
		},
	}})
	if err == nil || !strings.Contains(err.Error(), "Query key condition not supported") {
		t.Fatalf("begins_with on N sort key: expected the grammar rejection, got %v", err)
	}

	// Index queries check the index key schema: the GSI hashes on the
	// N-typed gid attribute.
	_, err = svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "TypedTable",
		"IndexName":                 "gsi1",
		"KeyConditionExpression":    "gid = :g",
		"ExpressionAttributeValues": map[string]interface{}{":g": map[string]interface{}{"S": "7"}},
	}})
	if err == nil || !strings.Contains(err.Error(), "Condition parameter type does not match schema type") {
		t.Fatalf("S value against N index hash key: expected the type-mismatch rejection, got %v", err)
	}

	// Positive control: correctly typed conditions still answer.
	resp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "TypedTable",
		"KeyConditionExpression": "id = :a AND nsk > :lo",
		"ExpressionAttributeValues": map[string]interface{}{
			":a":  map[string]interface{}{"S": "A"},
			":lo": map[string]interface{}{"N": "0"},
		},
	}})
	if err != nil {
		t.Fatalf("correctly typed query: unexpected error: %v", err)
	}
	if count := resp.(map[string]interface{})["Count"].(int); count != 1 {
		t.Fatalf("correctly typed query: expected Count 1, got %v", count)
	}
}

// newBinarySortKeyTestService creates a table whose sort key is
// Binary-typed — the key schema a begins_with prefix match must serve
// bytewise.
func newBinarySortKeyTestService(t *testing.T) (*DynamoDBService, *request.RequestContext) {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name: "BinarySortTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "bsk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "bsk", AttributeType: dbstore.ScalarAttributeTypeB},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return svc, request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
}

// begins_with on a Binary sort key matches bytewise prefixes on both
// request planes: the key-face exclusion is Number alone ("You cannot use
// this function with a sort key that is of type Number" — the Query
// key-condition grammar), and the legacy KeyConditions face types both the
// operand and the target "String or Binary (not a Number or a set type)".
// The two planes share the key-condition engine, and a String operand
// against the Binary key takes the schema type-mismatch rejection.
func TestQueryBeginsWithBinarySortKey(t *testing.T) {
	svc, reqCtx := newBinarySortKeyTestService(t)
	for _, sk := range []string{"AQE=", "AQI=", "AgE="} { // 0x01 0x01, 0x01 0x02, 0x02 0x01
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "BinarySortTable",
			"Item": map[string]interface{}{
				"id":  map[string]interface{}{"S": "A"},
				"bsk": map[string]interface{}{"B": sk},
			},
		}}); err != nil {
			t.Fatalf("seed %s: %v", sk, err)
		}
	}

	resp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "BinarySortTable",
		"KeyConditionExpression": "id = :a AND begins_with(bsk, :p)",
		"ExpressionAttributeValues": map[string]interface{}{
			":a": map[string]interface{}{"S": "A"},
			":p": map[string]interface{}{"B": "AQ=="}, // the single byte 0x01
		},
	}})
	if err != nil {
		t.Fatalf("begins_with on a Binary sort key: %v", err)
	}
	if count := resp.(map[string]interface{})["Count"].(int); count != 2 {
		t.Fatalf("begins_with binary prefix: expected Count 2, got %v", count)
	}

	legacy, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "BinarySortTable",
		"KeyConditions": map[string]interface{}{
			"id":  map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
			"bsk": map[string]interface{}{"ComparisonOperator": "BEGINS_WITH", "AttributeValueList": []interface{}{map[string]interface{}{"B": "AQ=="}}},
		},
	}})
	if err != nil {
		t.Fatalf("legacy BEGINS_WITH on a Binary sort key: %v", err)
	}
	if count := legacy.(map[string]interface{})["Count"].(int); count != 2 {
		t.Fatalf("legacy BEGINS_WITH binary prefix: expected Count 2, got %v", count)
	}

	_, err = svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":              "BinarySortTable",
		"KeyConditionExpression": "id = :a AND begins_with(bsk, :p)",
		"ExpressionAttributeValues": map[string]interface{}{
			":a": map[string]interface{}{"S": "A"},
			":p": map[string]interface{}{"S": "\x01"},
		},
	}})
	if err == nil || !strings.Contains(err.Error(), "Condition parameter type does not match schema type") {
		t.Fatalf("S operand against B sort key: expected the type-mismatch rejection, got %v", err)
	}
}

// newIndexEskTestService creates a table whose GSI carries a Number-typed
// sort attribute, so a wrong-typed echoed index key in an ExclusiveStartKey
// is observable on an index query.
func newIndexEskTestService(t *testing.T) (*DynamoDBService, *request.RequestContext) {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name: "IndexEskTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "gid", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "gsk", AttributeType: dbstore.ScalarAttributeTypeN},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []*dbstore.GlobalSecondaryIndex{
			{
				IndexName: "gesi",
				KeySchema: []*dbstore.KeySchemaElement{
					{AttributeName: "gid", KeyType: dbstore.KeyTypeHash},
					{AttributeName: "gsk", KeyType: dbstore.KeyTypeRange},
				},
				Projection: &dbstore.Projection{ProjectionType: "ALL"},
			},
		},
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return svc, request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
}

func TestIndexQueryRejectsWrongTypedStartKeyIndexAttributes(t *testing.T) {
	svc, reqCtx := newIndexEskTestService(t)
	for _, gsk := range []string{"1", "2"} {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "IndexEskTable",
			"Item": map[string]interface{}{
				"id":  map[string]interface{}{"S": "A"},
				"sk":  map[string]interface{}{"S": gsk},
				"gid": map[string]interface{}{"S": "G"},
				"gsk": map[string]interface{}{"N": gsk},
			},
		}}); err != nil {
			t.Fatalf("seed %s: %v", gsk, err)
		}
	}

	base := map[string]interface{}{
		"TableName":              "IndexEskTable",
		"IndexName":              "gesi",
		"KeyConditionExpression": "gid = :g",
		"ExpressionAttributeValues": map[string]interface{}{
			":g": map[string]interface{}{"S": "G"},
		},
	}

	// A start key echoing the read target's own LastEvaluatedKey shape,
	// with the index sort attribute flipped to the wrong type: the walk
	// must not resume at a position the value's own encoding dictates.
	bad := map[string]interface{}{}
	for k, v := range base {
		bad[k] = v
	}
	bad["ExclusiveStartKey"] = map[string]interface{}{
		"id":  map[string]interface{}{"S": "A"},
		"sk":  map[string]interface{}{"S": "1"},
		"gid": map[string]interface{}{"S": "G"},
		"gsk": map[string]interface{}{"S": "1"},
	}
	_, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: bad})
	if err == nil || !strings.Contains(err.Error(), "Type mismatch for key gsk expected: N actual: S") {
		t.Fatalf("wrong-typed index sort in start key: expected the key type-mismatch rejection, got %v", err)
	}

	// The correctly typed start key resumes the walk after the first item.
	good := map[string]interface{}{}
	for k, v := range base {
		good[k] = v
	}
	good["ExclusiveStartKey"] = map[string]interface{}{
		"id":  map[string]interface{}{"S": "A"},
		"sk":  map[string]interface{}{"S": "1"},
		"gid": map[string]interface{}{"S": "G"},
		"gsk": map[string]interface{}{"N": "1"},
	}
	resp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: good})
	if err != nil {
		t.Fatalf("correctly typed start key: unexpected error: %v", err)
	}
	if count := resp.(map[string]interface{})["Count"].(int); count != 1 {
		t.Fatalf("correctly typed start key: expected Count 1 (second of 2 items), got %v", count)
	}

	// An index Scan shares the preamble, so the same wrong-typed echo is
	// rejected on the scan plane too.
	scanBad := map[string]interface{}{
		"TableName": "IndexEskTable",
		"IndexName": "gesi",
	}
	scanBad["ExclusiveStartKey"] = bad["ExclusiveStartKey"]
	_, err = svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: scanBad})
	if err == nil || !strings.Contains(err.Error(), "Type mismatch for key gsk expected: N actual: S") {
		t.Fatalf("index scan with wrong-typed index sort in start key: expected the key type-mismatch rejection, got %v", err)
	}
}

func TestIndexReadLastEvaluatedKeySurvivesProjection(t *testing.T) {
	svc, reqCtx := newIndexEskTestService(t)
	for _, gsk := range []string{"1", "2", "3"} {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "IndexEskTable",
			"Item": map[string]interface{}{
				"id":  map[string]interface{}{"S": "A"},
				"sk":  map[string]interface{}{"S": gsk},
				"gid": map[string]interface{}{"S": "G"},
				"gsk": map[string]interface{}{"N": gsk},
				"v":   map[string]interface{}{"S": "val" + gsk},
			},
		}}); err != nil {
			t.Fatalf("seed %s: %v", gsk, err)
		}
	}

	// The projection omits every key attribute of the read target; the
	// LastEvaluatedKey must still name the read target's full key schema,
	// or the echoed start key cannot resume the walk on page 2.
	page1, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "IndexEskTable",
		"IndexName":                 "gesi",
		"KeyConditionExpression":    "gid = :g",
		"ExpressionAttributeValues": map[string]interface{}{":g": map[string]interface{}{"S": "G"}},
		"ExpressionAttributeNames":  map[string]interface{}{"#v": "v"},
		"ProjectionExpression":      "#v",
		"Limit":                     2,
	}})
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	resp := page1.(map[string]interface{})
	lek, hasLEK := resp["LastEvaluatedKey"]
	if !hasLEK {
		t.Fatal("page 1 of 2 pages: expected LastEvaluatedKey")
	}
	lekMap := lek.(map[string]interface{})
	for _, key := range []string{"id", "sk", "gid", "gsk"} {
		if _, ok := lekMap[key]; !ok {
			t.Fatalf("LastEvaluatedKey lacks the read target key attribute %q (have %v)", key, lekMap)
		}
	}

	page2, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "IndexEskTable",
		"IndexName":                 "gesi",
		"KeyConditionExpression":    "gid = :g",
		"ExpressionAttributeValues": map[string]interface{}{":g": map[string]interface{}{"S": "G"}},
		"ExpressionAttributeNames":  map[string]interface{}{"#v": "v"},
		"ProjectionExpression":      "#v",
		"ExclusiveStartKey":         lekMap,
	}})
	if err != nil {
		t.Fatalf("page 2 with the echoed key: %v", err)
	}
	if count := page2.(map[string]interface{})["Count"].(int); count != 1 {
		t.Fatalf("page 2: expected Count 1 (third of 3 items), got %v", count)
	}

	// The Scan plane composes its LastEvaluatedKey through the same
	// merge before its own projection application.
	scanPage1, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                "IndexEskTable",
		"IndexName":                "gesi",
		"ExpressionAttributeNames": map[string]interface{}{"#v": "v"},
		"ProjectionExpression":     "#v",
		"Limit":                    2,
	}})
	if err != nil {
		t.Fatalf("scan page 1: %v", err)
	}
	scanLEK, hasScanLEK := scanPage1.(map[string]interface{})["LastEvaluatedKey"]
	if !hasScanLEK {
		t.Fatal("scan page 1 of 2 pages: expected LastEvaluatedKey")
	}
	scanLEKMap := scanLEK.(map[string]interface{})
	if _, ok := scanLEKMap["gsk"]; !ok {
		t.Fatalf("scan LastEvaluatedKey lacks the index sort attribute (have %v)", scanLEKMap)
	}
}

// TestMalformedFilterExpressionFailsTheRequest pins the read-plane
// contract for FilterExpression: the expression is request validation —
// an operator outside the grammar, an undefined substitution or a string
// literal answers ValidationException, never a successful response whose
// filter silently evaluated false for every item.
func TestMalformedFilterExpressionFailsTheRequest(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	for _, sk := range []string{"1", "2", "3"} {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "LegacyTable",
			"Item": map[string]interface{}{
				"id": map[string]interface{}{"S": "A"},
				"sk": map[string]interface{}{"S": sk},
				"v":  map[string]interface{}{"N": sk},
			},
		}}); err != nil {
			t.Fatalf("seed %s: %v", sk, err)
		}
	}

	values := map[string]interface{}{
		":a": map[string]interface{}{"S": "A"},
		":p": map[string]interface{}{"N": "2"},
	}
	rejections := []struct{ name, expr string }{
		{"operator outside the grammar", "v NE :p"},
		{"undefined placeholder", "v = :missing"},
		{"undefined alias", "#missing = :p"},
		{"string literal operand", "v = '2'"},
	}
	for _, tc := range rejections {
		queryParams := map[string]interface{}{
			"TableName":                 "LegacyTable",
			"KeyConditionExpression":    "id = :a",
			"FilterExpression":          tc.expr,
			"ExpressionAttributeValues": values,
		}
		if _, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: queryParams}); err == nil || !strings.Contains(err.Error(), "ValidationException") {
			t.Fatalf("Query %s: expected a ValidationException, got %v", tc.name, err)
		}
		scanParams := map[string]interface{}{
			"TableName":                 "LegacyTable",
			"FilterExpression":          tc.expr,
			"ExpressionAttributeValues": values,
		}
		if _, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: scanParams}); err == nil || !strings.Contains(err.Error(), "ValidationException") {
			t.Fatalf("Scan %s: expected a ValidationException, got %v", tc.name, err)
		}
	}

	// Positive control: a well-formed filter still narrows the page.
	queryResp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "LegacyTable",
		"KeyConditionExpression":    "id = :a",
		"FilterExpression":          "v = :p",
		"ExpressionAttributeValues": values,
	}})
	if err != nil {
		t.Fatalf("well-formed filter query: %v", err)
	}
	if count := queryResp.(map[string]interface{})["Count"].(int); count != 1 {
		t.Fatalf("well-formed filter query: expected 1 filtered item, got %v", count)
	}
	if scanned := queryResp.(map[string]interface{})["ScannedCount"].(int); scanned != 3 {
		t.Fatalf("well-formed filter query: expected 3 scanned items, got %v", scanned)
	}
}
