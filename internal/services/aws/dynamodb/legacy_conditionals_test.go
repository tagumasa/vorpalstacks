package dynamodb

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The legacy conditional parameter family (Expected, AttributesToGet,
// KeyConditions, QueryFilter, ScanFilter, ConditionalOperator) is
// translated onto the expression machinery that replaced it; these tests
// pin the wire-plane behaviour of each member and the exclusivity rules
// against the expression parameters.

func newLegacyTestService(t *testing.T) (*DynamoDBService, *request.RequestContext) {
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
		Name: "LegacyTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode: dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return svc, request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
}

func legacyPut(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, extra map[string]interface{}) error {
	t.Helper()
	params := map[string]interface{}{
		"TableName": "LegacyTable",
		"Item": map[string]interface{}{
			"id": map[string]interface{}{"S": "A"},
			"sk": map[string]interface{}{"S": "1"},
			"v":  map[string]interface{}{"N": "1"},
		},
	}
	for k, val := range extra {
		params[k] = val
	}
	_, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
	return err
}

func TestExpectedGuardsConditionalWrites(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)

	// Exists:false is the legacy create-guard: the first write succeeds and
	// a second write of the same item fails the guard.
	existsFalse := map[string]interface{}{"Expected": map[string]interface{}{
		"id": map[string]interface{}{"Exists": false},
	}}
	if err := legacyPut(t, svc, reqCtx, existsFalse); err != nil {
		t.Fatalf("Expected Exists:false on a fresh item: unexpected error %v", err)
	}
	if err := legacyPut(t, svc, reqCtx, existsFalse); err == nil || !strings.Contains(err.Error(), "ConditionalCheckFailed") {
		t.Fatalf("Expected Exists:false on an existing item: expected ConditionalCheckFailed, got %v", err)
	}

	// The Value form guards by equality.
	valueGuard := map[string]interface{}{"Expected": map[string]interface{}{
		"v": map[string]interface{}{"Value": map[string]interface{}{"N": "1"}},
	}}
	if err := legacyPut(t, svc, reqCtx, valueGuard); err != nil {
		t.Fatalf("Expected Value matching: unexpected error %v", err)
	}
	mismatched := map[string]interface{}{"Expected": map[string]interface{}{
		"v": map[string]interface{}{"Value": map[string]interface{}{"N": "9"}},
	}}
	if err := legacyPut(t, svc, reqCtx, mismatched); err == nil || !strings.Contains(err.Error(), "ConditionalCheckFailed") {
		t.Fatalf("Expected Value mismatch: expected ConditionalCheckFailed, got %v", err)
	}

	// The ComparisonOperator form compares; ConditionalOperator OR lets one
	// true condition carry a false one.
	orGuard := map[string]interface{}{
		"Expected": map[string]interface{}{
			"v":  map[string]interface{}{"ComparisonOperator": "GT", "AttributeValueList": []interface{}{map[string]interface{}{"N": "9"}}},
			"id": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
		},
		"ConditionalOperator": "OR",
	}
	if err := legacyPut(t, svc, reqCtx, orGuard); err != nil {
		t.Fatalf("Expected OR with one true condition: unexpected error %v", err)
	}
	andGuard := map[string]interface{}{"Expected": orGuard["Expected"]}
	if err := legacyPut(t, svc, reqCtx, andGuard); err == nil || !strings.Contains(err.Error(), "ConditionalCheckFailed") {
		t.Fatalf("Expected AND with one false condition: expected ConditionalCheckFailed, got %v", err)
	}

	// Exclusivity: Expected never rides with ConditionExpression, and the
	// two entry forms never combine in one entry.
	if err := legacyPut(t, svc, reqCtx, map[string]interface{}{
		"Expected":            map[string]interface{}{"v": map[string]interface{}{"Exists": true}},
		"ConditionExpression": "attribute_exists(v)",
	}); err == nil || !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("Expected with ConditionExpression: expected exclusivity rejection, got %v", err)
	}
	if err := legacyPut(t, svc, reqCtx, map[string]interface{}{
		"Expected": map[string]interface{}{
			"v": map[string]interface{}{
				"Value":              map[string]interface{}{"N": "1"},
				"ComparisonOperator": "EQ",
			},
		},
	}); err == nil || !strings.Contains(err.Error(), "cannot carry Value or Exists together") {
		t.Fatalf("Expected mixed entry forms: expected rejection, got %v", err)
	}
	if err := legacyPut(t, svc, reqCtx, map[string]interface{}{
		"Expected": map[string]interface{}{
			"v": map[string]interface{}{"ComparisonOperator": "BETWEEN", "AttributeValueList": []interface{}{
				map[string]interface{}{"N": "1"},
			}},
		},
	}); err == nil || !strings.Contains(err.Error(), "exactly 2") {
		t.Fatalf("Expected BETWEEN with one value: expected arity rejection, got %v", err)
	}

	// Exists is strictly boolean: a present-but-mistyped value is a
	// request error, never a silently ignored guard that falls through to
	// the Value form's behaviour.
	if err := legacyPut(t, svc, reqCtx, map[string]interface{}{
		"Expected": map[string]interface{}{"v": map[string]interface{}{"Exists": "true"}},
	}); err == nil || !strings.Contains(err.Error(), "Exists must be a boolean") {
		t.Fatalf("Expected Exists mistyped: expected rejection, got %v", err)
	}
}

func TestAttributesToGetProjectsReads(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	if err := legacyPut(t, svc, reqCtx, nil); err != nil {
		t.Fatalf("seed: %v", err)
	}

	got, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":       "LegacyTable",
		"Key":             map[string]interface{}{"id": map[string]interface{}{"S": "A"}, "sk": map[string]interface{}{"S": "1"}},
		"AttributesToGet": []interface{}{"v"},
	}})
	if err != nil {
		t.Fatalf("AttributesToGet read: %v", err)
	}
	item := got.(map[string]interface{})["Item"].(map[string]interface{})
	if _, present := item["v"]; !present {
		t.Fatalf("AttributesToGet read: requested attribute v missing: %v", item)
	}
	if _, present := item["id"]; present {
		t.Fatalf("AttributesToGet read: unrequested attribute id projected in: %v", item)
	}

	if _, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "LegacyTable",
		"Key":                  map[string]interface{}{"id": map[string]interface{}{"S": "A"}, "sk": map[string]interface{}{"S": "1"}},
		"AttributesToGet":      []interface{}{"v"},
		"ProjectionExpression": "v",
	}}); err == nil || !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("AttributesToGet with ProjectionExpression: expected exclusivity rejection, got %v", err)
	}
}

// AttributesToGet entries are plain attribute names: a name that itself
// contains '.' addresses the top-level attribute so named, never a nested
// document path — the member predates expression attribute names and
// cannot descend into a List or Map.
func TestAttributesToGetAddressesTopLevelNameOnly(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"Item": map[string]interface{}{
			"id":          map[string]interface{}{"S": "A"},
			"sk":          map[string]interface{}{"S": "1"},
			"dotted.name": map[string]interface{}{"S": "top-level"},
			"dotted":      map[string]interface{}{"M": map[string]interface{}{"name": map[string]interface{}{"S": "nested"}}},
		},
	}}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	got, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":       "LegacyTable",
		"Key":             map[string]interface{}{"id": map[string]interface{}{"S": "A"}, "sk": map[string]interface{}{"S": "1"}},
		"AttributesToGet": []interface{}{"dotted.name"},
	}})
	if err != nil {
		t.Fatalf("AttributesToGet read: %v", err)
	}
	item := got.(map[string]interface{})["Item"].(map[string]interface{})
	if av, ok := item["dotted.name"].(map[string]interface{}); !ok || av["S"] != "top-level" {
		t.Fatalf("projected item = %v, want the top-level attribute dotted.name", item)
	}
	if _, hasNested := item["dotted"]; hasNested {
		t.Fatalf("projected item = %v, the name must not re-split into the nested path", item)
	}
}

func TestKeyConditionsDriveQuery(t *testing.T) {
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

	// Partition equality alone reads the whole partition.
	resp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"KeyConditions": map[string]interface{}{
			"id": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
		},
	}})
	if err != nil {
		t.Fatalf("KeyConditions partition query: %v", err)
	}
	if count := resp.(map[string]interface{})["Count"].(int); count != 3 {
		t.Fatalf("KeyConditions partition query: expected 3 items, got %v", count)
	}

	// A sort-key BETWEEN narrows the read.
	resp, err = svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"KeyConditions": map[string]interface{}{
			"id": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
			"sk": map[string]interface{}{"ComparisonOperator": "BETWEEN", "AttributeValueList": []interface{}{map[string]interface{}{"S": "1"}, map[string]interface{}{"S": "2"}}},
		},
	}})
	if err != nil {
		t.Fatalf("KeyConditions sort-range query: %v", err)
	}
	if count := resp.(map[string]interface{})["Count"].(int); count != 2 {
		t.Fatalf("KeyConditions sort-range query: expected 2 items, got %v", count)
	}

	// Conditions outside the key schema, and mixing with the expression
	// parameter, are rejections.
	badCases := []struct {
		name   string
		extra  map[string]interface{}
		expect string
	}{
		{"non-key attribute", map[string]interface{}{
			"KeyConditions": map[string]interface{}{
				"v":  map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"N": "1"}}},
				"id": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
			},
		}, "key schema"},
		{"partition not EQ", map[string]interface{}{
			"KeyConditions": map[string]interface{}{
				"id": map[string]interface{}{"ComparisonOperator": "GT", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
			},
		}, "must use EQ"},
		{"missing partition", map[string]interface{}{
			"KeyConditions": map[string]interface{}{
				"sk": map[string]interface{}{"ComparisonOperator": "GT", "AttributeValueList": []interface{}{map[string]interface{}{"S": "0"}}},
			},
		}, "missed key schema element"},
		{"both families", map[string]interface{}{
			"KeyConditionExpression": "id = :a",
			"ExpressionAttributeValues": map[string]interface{}{
				":a": map[string]interface{}{"S": "A"},
			},
			"KeyConditions": map[string]interface{}{
				"id": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
			},
		}, "cannot be used together"},
	}
	for _, tc := range badCases {
		params := map[string]interface{}{"TableName": "LegacyTable"}
		for k, val := range tc.extra {
			params[k] = val
		}
		_, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		if err == nil || !strings.Contains(err.Error(), tc.expect) {
			t.Fatalf("%s: expected %q rejection, got %v", tc.name, tc.expect, err)
		}
	}
}

func TestLegacyFiltersDriveQueryAndScan(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	tags := map[string][]interface{}{"1": {"alpha"}, "2": {"beta", "gamma"}, "3": {"delta"}}
	for _, sk := range []string{"1", "2", "3"} {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "LegacyTable",
			"Item": map[string]interface{}{
				"id":   map[string]interface{}{"S": "A"},
				"sk":   map[string]interface{}{"S": sk},
				"v":    map[string]interface{}{"N": sk},
				"tags": map[string]interface{}{"SS": tags[sk]},
			},
		}}); err != nil {
			t.Fatalf("seed %s: %v", sk, err)
		}
	}

	queryResp, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"KeyConditions": map[string]interface{}{
			"id": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
		},
		"QueryFilter": map[string]interface{}{
			"v": map[string]interface{}{"ComparisonOperator": "GT", "AttributeValueList": []interface{}{map[string]interface{}{"N": "1"}}},
		},
	}})
	if err != nil {
		t.Fatalf("QueryFilter query: %v", err)
	}
	if count := queryResp.(map[string]interface{})["Count"].(int); count != 2 {
		t.Fatalf("QueryFilter query: expected 2 filtered items, got %v", count)
	}

	scanResp, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"ScanFilter": map[string]interface{}{
			"v": map[string]interface{}{"ComparisonOperator": "LE", "AttributeValueList": []interface{}{map[string]interface{}{"N": "1"}}},
		},
	}})
	if err != nil {
		t.Fatalf("ScanFilter scan: %v", err)
	}
	if count := scanResp.(map[string]interface{})["Count"].(int); count != 1 {
		t.Fatalf("ScanFilter scan: expected 1 filtered item, got %v", count)
	}

	// Legacy IN is per-element membership: a scalar candidate matches when
	// it is a member of a set attribute; a non-set attribute keeps the
	// equality-any reading.
	inResp, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"ScanFilter": map[string]interface{}{
			"tags": map[string]interface{}{"ComparisonOperator": "IN", "AttributeValueList": []interface{}{map[string]interface{}{"S": "beta"}}},
		},
	}})
	if err != nil {
		t.Fatalf("ScanFilter IN membership scan: %v", err)
	}
	if count := inResp.(map[string]interface{})["Count"].(int); count != 1 {
		t.Fatalf("ScanFilter IN membership scan: expected 1 item (tags contains beta), got %v", count)
	}

	inMultiResp, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"ScanFilter": map[string]interface{}{
			"tags": map[string]interface{}{"ComparisonOperator": "IN", "AttributeValueList": []interface{}{
				map[string]interface{}{"S": "alpha"}, map[string]interface{}{"S": "delta"},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("ScanFilter IN multi-candidate scan: %v", err)
	}
	if count := inMultiResp.(map[string]interface{})["Count"].(int); count != 2 {
		t.Fatalf("ScanFilter IN multi-candidate scan: expected 2 items (alpha or delta member), got %v", count)
	}

	scalarInResp, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"ScanFilter": map[string]interface{}{
			"sk": map[string]interface{}{"ComparisonOperator": "IN", "AttributeValueList": []interface{}{
				map[string]interface{}{"S": "2"}, map[string]interface{}{"S": "3"},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("ScanFilter IN scalar scan: %v", err)
	}
	if count := scalarInResp.(map[string]interface{})["Count"].(int); count != 2 {
		t.Fatalf("ScanFilter IN scalar scan: expected 2 items (sk equals 2 or 3), got %v", count)
	}

	// Exclusivity with FilterExpression, and the OR combinator.
	if _, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "LegacyTable",
		"FilterExpression":          "v = :one",
		"ExpressionAttributeValues": map[string]interface{}{":one": map[string]interface{}{"N": "1"}},
		"ScanFilter": map[string]interface{}{
			"v": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"N": "1"}}},
		},
	}}); err == nil || !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("ScanFilter with FilterExpression: expected exclusivity rejection, got %v", err)
	}

	orResp, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"ScanFilter": map[string]interface{}{
			"v":  map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"N": "1"}}},
			"sk": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "3"}}},
		},
		"ConditionalOperator": "OR",
	}})
	if err != nil {
		t.Fatalf("ScanFilter OR scan: %v", err)
	}
	if count := orResp.(map[string]interface{})["Count"].(int); count != 2 {
		t.Fatalf("ScanFilter OR scan: expected 2 items (one per true disjunct), got %v", count)
	}

	// The expression attribute maps are rejected alongside either legacy
	// filter with the same exclusivity the Expected path applies — the
	// legacy filter carries its values inline, so a supplied expression
	// map is a mixed request, never a silently ignored member.
	mixedQuery := map[string]interface{}{
		"TableName": "LegacyTable",
		"KeyConditions": map[string]interface{}{
			"id": map[string]interface{}{"ComparisonOperator": "EQ", "AttributeValueList": []interface{}{map[string]interface{}{"S": "A"}}},
		},
		"QueryFilter": map[string]interface{}{
			"v": map[string]interface{}{"ComparisonOperator": "GT", "AttributeValueList": []interface{}{map[string]interface{}{"N": "1"}}},
		},
		"ExpressionAttributeValues": map[string]interface{}{":one": map[string]interface{}{"N": "1"}},
	}
	if _, err := svc.Query(context.Background(), reqCtx, &request.ParsedRequest{Parameters: mixedQuery}); err == nil || !strings.Contains(err.Error(), "cannot be used with ExpressionAttributeNames or ExpressionAttributeValues") {
		t.Fatalf("QueryFilter with expression values: expected exclusivity rejection, got %v", err)
	}
	mixedScan := map[string]interface{}{
		"TableName": "LegacyTable",
		"ScanFilter": map[string]interface{}{
			"v": map[string]interface{}{"ComparisonOperator": "LE", "AttributeValueList": []interface{}{map[string]interface{}{"N": "1"}}},
		},
		"ExpressionAttributeNames": map[string]interface{}{"#v": "v"},
	}
	if _, err := svc.Scan(context.Background(), reqCtx, &request.ParsedRequest{Parameters: mixedScan}); err == nil || !strings.Contains(err.Error(), "cannot be used with ExpressionAttributeNames or ExpressionAttributeValues") {
		t.Fatalf("ScanFilter with expression names: expected exclusivity rejection, got %v", err)
	}
}

// TestAttributeUpdatesMemberValidation pins the legacy AttributeUpdates
// member contract: an omitted Action defaults to PUT, an unknown Action and
// a Value missing under PUT/ADD are ValidationExceptions (Value is optional
// only for DELETE), an invalid Value rejects instead of storing a nil
// attribute, DELETE without Value removes the attribute, and DELETE with a
// set Value subtracts the set's elements from the stored set instead of
// destroying the attribute.
func TestAttributeUpdatesMemberValidation(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)

	seed := map[string]interface{}{
		"id": map[string]interface{}{"S": "A"},
		"sk": map[string]interface{}{"S": "1"},
		"v":  map[string]interface{}{"N": "1"},
		"ss": map[string]interface{}{"SS": []interface{}{"a", "b", "c"}},
	}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable", "Item": seed,
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	update := func(updates map[string]interface{}) error {
		_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":        "LegacyTable",
			"Key":              map[string]interface{}{"id": map[string]interface{}{"S": "A"}, "sk": map[string]interface{}{"S": "1"}},
			"AttributeUpdates": updates,
		}})
		return err
	}
	getAttr := func(name string) interface{} {
		resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "LegacyTable",
			"Key":       map[string]interface{}{"id": map[string]interface{}{"S": "A"}, "sk": map[string]interface{}{"S": "1"}},
		}})
		if err != nil {
			t.Fatalf("read-back get: %v", err)
		}
		item := resp.(map[string]interface{})["Item"].(map[string]interface{})
		return item[name]
	}

	// PUT without Value and ADD without Value are rejected; Value is
	// optional only for DELETE.
	if err := update(map[string]interface{}{"v": map[string]interface{}{"Action": "PUT"}}); err == nil || !strings.Contains(err.Error(), "Value must be specified for AttributeUpdates action PUT") {
		t.Fatalf("PUT without Value: expected rejection, got %v", err)
	}
	if err := update(map[string]interface{}{"v": map[string]interface{}{"Action": "ADD"}}); err == nil || !strings.Contains(err.Error(), "Value must be specified for AttributeUpdates action ADD") {
		t.Fatalf("ADD without Value: expected rejection, got %v", err)
	}

	// An unknown Action is the enum rejection; a member that is not an
	// AttributeValueUpdate structure is rejected rather than skipped.
	if err := update(map[string]interface{}{"v": map[string]interface{}{"Action": "BOGUS", "Value": map[string]interface{}{"N": "2"}}}); err == nil || !strings.Contains(err.Error(), "must be one of ADD, DELETE, PUT") {
		t.Fatalf("unknown Action: expected the enum rejection, got %v", err)
	}
	if err := update(map[string]interface{}{"v": "not-a-structure"}); err == nil {
		t.Fatal("non-structure member: expected rejection, got success")
	}

	// An invalid Value rejects the update; nothing is stored, not even a
	// nil attribute.
	if err := update(map[string]interface{}{"v": map[string]interface{}{"Action": "PUT", "Value": map[string]interface{}{"S": "x", "N": "1"}}}); err == nil {
		t.Fatal("double-typed Value: expected rejection, got success")
	}

	// An omitted Action defaults to PUT: the value lands.
	if err := update(map[string]interface{}{"extra": map[string]interface{}{"Value": map[string]interface{}{"S": "by-default"}}}); err != nil {
		t.Fatalf("default PUT: unexpected error %v", err)
	}
	if got := getAttr("extra").(map[string]interface{})["S"]; got != "by-default" {
		t.Fatalf("default PUT: expected extra=by-default, got %v", got)
	}

	// DELETE with a set Value subtracts the elements: [a,b,c] - [a,c] = [b].
	if err := update(map[string]interface{}{"ss": map[string]interface{}{"Action": "DELETE", "Value": map[string]interface{}{"SS": []interface{}{"a", "c"}}}}); err != nil {
		t.Fatalf("DELETE set subtraction: unexpected error %v", err)
	}
	remaining, ok := getAttr("ss").(map[string]interface{})["SS"].([]string)
	if !ok || len(remaining) != 1 || remaining[0] != "b" {
		t.Fatalf("DELETE set subtraction: expected [b], got %v", remaining)
	}

	// DELETE without Value removes the attribute outright.
	if err := update(map[string]interface{}{"ss": map[string]interface{}{"Action": "DELETE"}}); err != nil {
		t.Fatalf("DELETE without Value: unexpected error %v", err)
	}
	if got := getAttr("ss"); got != nil {
		t.Fatalf("DELETE without Value: expected the attribute gone, got %v", got)
	}
}
