package dynamodb

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The Cores under test validate their inputs before touching the store, so
// the invalid-input cases run against a nil store: any attempt to use it
// would panic the test, which is itself the assertion that validation
// happens first.

func typedKeyTable(name string) *dbstore.Table {
	return &dbstore.Table{
		Name: name,
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeN},
		},
	}
}

func TestBuildConditionCheckerEvaluatesSyntheticMissingItem(t *testing.T) {
	table := typedKeyTable("t1")

	// An empty expression yields a nil checker — unconditional write.
	checker, err := buildConditionChecker(table.Name, map[string]*dbstore.AttributeValue{"pk": dbstore.NumberValue("1")}, conditionSpec{})
	if err != nil {
		t.Fatalf("empty condition: unexpected error %v", err)
	}
	if checker != nil {
		t.Fatalf("empty condition: expected nil checker, got %v", checker)
	}

	key := map[string]*dbstore.AttributeValue{"pk": dbstore.NumberValue("1")}
	checker, err = buildConditionChecker(table.Name, key, conditionSpec{
		Expr:   "attribute_not_exists(pk)",
		Names:  nil,
		Values: nil,
	})
	if err != nil {
		t.Fatalf("condition parse: unexpected error %v", err)
	}

	// A missing item evaluates against the synthetic empty item, so the
	// condition holds.
	if err := checker(nil, true); err != nil {
		t.Fatalf("attribute_not_exists on missing item: unexpected error %v", err)
	}

	// An existing item holding the attribute fails the same condition.
	existing := &dbstore.Item{
		TableName:  table.Name,
		Key:        key,
		Attributes: key,
	}
	if err := checker(existing, false); !errors.Is(err, ErrConditionalCheckFailed) {
		t.Fatalf("attribute_not_exists on existing item: expected ErrConditionalCheckFailed, got %v", err)
	}
}

func TestValidateReturnValuesNoneAllOld(t *testing.T) {
	for _, rv := range []string{"", "NONE", "ALL_OLD"} {
		if !validateReturnValuesNoneAllOld(rv) {
			t.Errorf("ReturnValues %q must be accepted", rv)
		}
	}
	for _, rv := range []string{"ALL_NEW", "UPDATED_OLD", "UPDATED_NEW", "BOGUS"} {
		if validateReturnValuesNoneAllOld(rv) {
			t.Errorf("ReturnValues %q must be rejected", rv)
		}
	}
}

func TestPutItemCoreValidatesBeforeStore(t *testing.T) {
	svc := &DynamoDBService{}
	table := typedKeyTable("t1")

	// A string value on a number-typed key attribute is rejected before any
	// store access, with the typed key-mismatch ValidationException.
	_, err := svc.putItemCore(t.Context(), nil, "us-east-1", PutItemCoreInput{
		Table: table,
		Item:  map[string]*dbstore.AttributeValue{"pk": dbstore.StringValue("not-a-number")},
	})
	if err == nil || !strings.Contains(err.Error(), "Type mismatch for key pk") {
		t.Fatalf("key type mismatch: expected typed key-mismatch error, got %v", err)
	}

	// An item missing the key attribute entirely is a missing-key error.
	_, err = svc.putItemCore(t.Context(), nil, "us-east-1", PutItemCoreInput{
		Table: table,
		Item:  map[string]*dbstore.AttributeValue{"other": dbstore.StringValue("v")},
	})
	if !errors.Is(err, ErrMissingKey) {
		t.Fatalf("missing key: expected ErrMissingKey, got %v", err)
	}

	// PutItem recognises only NONE and ALL_OLD.
	_, err = svc.putItemCore(t.Context(), nil, "us-east-1", PutItemCoreInput{
		Table:        table,
		Item:         map[string]*dbstore.AttributeValue{"pk": dbstore.NumberValue("1")},
		ReturnValues: "ALL_NEW",
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("ReturnValues ALL_NEW: expected ErrInvalidParameter, got %v", err)
	}
}

func TestDeleteAndGetItemCoreValidateKeyBeforeStore(t *testing.T) {
	svc := &DynamoDBService{}
	table := typedKeyTable("t1")

	if _, err := svc.deleteItemCore(t.Context(), nil, "us-east-1", DeleteItemCoreInput{
		Table: table,
		Key:   map[string]*dbstore.AttributeValue{"pk": dbstore.StringValue("not-a-number")},
	}); err == nil || !strings.Contains(err.Error(), "Type mismatch for key pk") {
		t.Fatalf("deleteItemCore key type mismatch: expected typed key-mismatch error, got %v", err)
	}

	if _, err := svc.getItemCore(t.Context(), nil, table, map[string]*dbstore.AttributeValue{"pk": dbstore.StringValue("not-a-number")}, false, ""); err == nil || !strings.Contains(err.Error(), "Type mismatch for key pk") {
		t.Fatalf("getItemCore key type mismatch: expected typed key-mismatch error, got %v", err)
	}
}

func TestUpdateItemCoreValidatesInputBeforeStore(t *testing.T) {
	svc := &DynamoDBService{}
	table := typedKeyTable("t1")

	base := func() UpdateItemInput {
		return UpdateItemInput{
			Key:        map[string]*dbstore.AttributeValue{"pk": dbstore.NumberValue("1")},
			UpdateExpr: "SET a = :v",
			ExprAttrValues: map[string]*dbstore.AttributeValue{
				":v": dbstore.StringValue("x"),
			},
		}
	}

	// A key attribute whose type contradicts the schema is rejected.
	bad := base()
	bad.Key = map[string]*dbstore.AttributeValue{"pk": dbstore.StringValue("no")}
	if _, err := svc.updateItemCore(t.Context(), nil, "us-east-1", table, bad); err == nil || !strings.Contains(err.Error(), "Type mismatch for key pk") {
		t.Fatalf("key type mismatch: expected typed key-mismatch error, got %v", err)
	}

	// UpdateExpression and AttributeUpdates are mutually exclusive.
	both := base()
	both.AttrUpdates = map[string]interface{}{"a": map[string]interface{}{"Value": map[string]interface{}{"S": "x"}}}
	if _, err := svc.updateItemCore(t.Context(), nil, "us-east-1", table, both); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("mutually exclusive update members: expected ErrInvalidParameter, got %v", err)
	}

	// UpdateItem recognises five ReturnValues settings; the PutItem-only
	// NONE/ALL_OLD subset is not the whole enum but a bogus value is still
	// rejected.
	bogus := base()
	bogus.ReturnValues = "BOGUS"
	if _, err := svc.updateItemCore(t.Context(), nil, "us-east-1", table, bogus); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("bogus ReturnValues: expected ErrInvalidParameter, got %v", err)
	}
}

func TestValidateBillingModeConsistency(t *testing.T) {
	valid := &dbstore.ProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: 5}
	zeroRead := &dbstore.ProvisionedThroughput{ReadCapacityUnits: 0, WriteCapacityUnits: 5}

	cases := []struct {
		billing dbstore.BillingMode
		pt      *dbstore.ProvisionedThroughput
		want    bool
	}{
		{dbstore.BillingModeProvisioned, valid, true},
		{dbstore.BillingModeProvisioned, nil, false},
		{dbstore.BillingModeProvisioned, zeroRead, false},
		{dbstore.BillingModePayPerRequest, nil, true},
		{dbstore.BillingModePayPerRequest, valid, false},
	}
	for _, tc := range cases {
		if got := validateBillingModeConsistency(tc.billing, tc.pt); got != tc.want {
			t.Errorf("validateBillingModeConsistency(%s, %+v) = %v, want %v", tc.billing, tc.pt, got, tc.want)
		}
	}
}

func TestCreateTableCoreRejectsPayPerRequestWithThroughput(t *testing.T) {
	svc := &DynamoDBService{}
	in := CreateTableInput{
		TableName: "ppr-pt-table",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
		},
		AttributeDefinitions: []*dbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS},
		},
		BillingMode:           dbstore.BillingModePayPerRequest,
		ProvisionedThroughput: &dbstore.ProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: 5},
	}

	if _, err := svc.createTableCore(context.Background(), nil, nil, in); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("PAY_PER_REQUEST with ProvisionedThroughput: expected ErrInvalidParameter, got %v", err)
	}
}

func TestUpdateTimeToLiveCoreValidatesAttributeName(t *testing.T) {
	svc := &DynamoDBService{}

	if _, err := svc.updateTimeToLiveCore(t.Context(), nil, UpdateTimeToLiveInput{
		TableName:     "t1",
		SpecPresent:   true,
		Enabled:       true,
		AttributeName: "",
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty attribute name: expected ErrInvalidParameter, got %v", err)
	}

	// The model bounds the TTL attribute name at 255 characters.
	if _, err := svc.updateTimeToLiveCore(t.Context(), nil, UpdateTimeToLiveInput{
		TableName:     "t1",
		SpecPresent:   true,
		Enabled:       true,
		AttributeName: repeatChar('a', 256),
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("256-char attribute name: expected ErrInvalidParameter, got %v", err)
	}
}

func TestUpdateTimeToLiveCoreRejectsAbsentAndMalformedSpec(t *testing.T) {
	svc := &DynamoDBService{}

	// The wire carried no TimeToLiveSpecification map at all.
	if _, err := svc.updateTimeToLiveCore(t.Context(), nil, UpdateTimeToLiveInput{
		TableName:     "t1",
		Enabled:       true,
		AttributeName: "ttl",
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("absent specification: expected ErrInvalidParameter, got %v", err)
	}

	// Enabled arrived with the wrong JSON type.
	if _, err := svc.updateTimeToLiveCore(t.Context(), nil, UpdateTimeToLiveInput{
		TableName:     "t1",
		SpecPresent:   true,
		Malformed:     true,
		Enabled:       true,
		AttributeName: "ttl",
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("malformed Enabled: expected ErrInvalidParameter, got %v", err)
	}
}

func TestParseTimeToLiveSpec(t *testing.T) {
	in := parseTimeToLiveSpec(nil)
	if in.SpecPresent || in.Malformed || in.Enabled || in.AttributeName != "" {
		t.Fatalf("absent specification: got %+v", in)
	}

	// An empty specification map leaves Enabled at its default (false) and
	// the attribute name empty — both rejections belong to the Core.
	in = parseTimeToLiveSpec(map[string]interface{}{})
	if !in.SpecPresent || in.Malformed || in.Enabled || in.AttributeName != "" {
		t.Fatalf("empty specification: got %+v", in)
	}

	in = parseTimeToLiveSpec(map[string]interface{}{"Enabled": true, "AttributeName": "ttl"})
	if !in.SpecPresent || !in.Enabled || in.AttributeName != "ttl" || in.Malformed {
		t.Fatalf("well-formed specification: got %+v", in)
	}

	in = parseTimeToLiveSpec(map[string]interface{}{"Enabled": "yes", "AttributeName": "ttl"})
	if !in.SpecPresent || !in.Malformed {
		t.Fatalf("non-bool Enabled must set Malformed: got %+v", in)
	}

	in = parseTimeToLiveSpec(map[string]interface{}{"Enabled": true, "AttributeName": 7})
	if !in.SpecPresent || in.AttributeName != "" {
		t.Fatalf("non-string AttributeName must extract empty: got %+v", in)
	}
}

func TestListTablesCoreValidatesLimitAndMarker(t *testing.T) {
	svc := &DynamoDBService{}

	// The model's ListTablesInputLimit bounds Limit to 1-100; both planes
	// hand an absent limit through as the documented maximum.
	for _, limit := range []int{0, -1, 101} {
		if _, _, err := svc.listTablesCore(nil, "", limit); !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("limit %d: expected ErrInvalidParameter, got %v", limit, err)
		}
	}

	// The marker is an ExclusiveStartTableName and must be a valid table
	// name.
	if _, _, err := svc.listTablesCore(nil, "bad name!", 10); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("malformed marker: expected ErrInvalidParameter, got %v", err)
	}
}

func repeatChar(c byte, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = c
	}
	return string(b)
}

// A projection document path with a malformed bracket segment is rejected
// at parse time — the same validateBracketIndex contract the
// update-expression paths enforce — instead of silently parsing to index 0.
func TestParseProjectionExpressionRejectsMalformedBracketPath(t *testing.T) {
	bad := map[string]interface{}{
		"ProjectionExpression": "doc.list[x]",
	}
	if _, err := parseProjectionExpression(bad); err == nil {
		t.Fatal("malformed bracket path accepted")
	}
	negative := map[string]interface{}{
		"ProjectionExpression": "doc.list[-1]",
	}
	if _, err := parseProjectionExpression(negative); err == nil {
		t.Fatal("negative bracket index accepted")
	}
	unclosed := map[string]interface{}{
		"ProjectionExpression": "doc.list[2",
	}
	if _, err := parseProjectionExpression(unclosed); err == nil {
		t.Fatal("unclosed bracket accepted")
	}
	good := map[string]interface{}{
		"ProjectionExpression": "doc.list[2].name",
	}
	projection, err := parseProjectionExpression(good)
	if err != nil {
		t.Fatalf("valid bracket path rejected: %v", err)
	}
	if len(projection) != 1 || len(projection[0]) != 4 {
		t.Fatalf("projection=%v, want one four-segment path", projection)
	}

	// An undefined expression attribute name is a validation error on the
	// projection plane, exactly as on the condition and update planes: the
	// alias may not silently stand as a literal name the projection then
	// fails to find.
	undefinedAlias := map[string]interface{}{
		"ProjectionExpression": "#p",
	}
	if _, err := parseProjectionExpression(undefinedAlias); err == nil {
		t.Fatal("undefined alias accepted: the projection must reject a #name the names map does not define")
	}
}

// A bracket index too large to represent never wraps into a small valid
// index: the accumulation bails at the first digit that would overflow,
// so an absurd document-path index is a validation error instead of
// silently addressing an unrelated list position.
func TestBracketIndexOverflowIsRejected(t *testing.T) {
	if _, err := validateBracketIndex("18446744073709551617"); err == nil {
		t.Fatal("2^64+1 wraps to a small index when accumulated unchecked; must be rejected")
	}
	if _, err := validateBracketIndex(strings.Repeat("9", 25)); err == nil {
		t.Fatal("a 25-digit index must be rejected, not wrapped")
	}
	maxText := strconv.FormatInt(int64(math.MaxInt), 10)
	idx, err := validateBracketIndex(maxText)
	if err != nil {
		t.Fatalf("the largest representable index rejected: %v", err)
	}
	if idx != math.MaxInt {
		t.Fatalf("idx = %d, want the largest representable index", idx)
	}
}

// A whole-token projection alias stands for exactly one attribute name: when
// the documented name itself contains '.', the resolved name must stay a
// single top-level path segment instead of being re-split into a document
// path. The naming-rules page grounds the pin — attribute names containing
// dots are addressable through expression attribute names. The compound
// form (#p.#n) keeps descending: each of its segments resolves separately.
func TestGetItemProjectionAliasAddressesDottedNameAsOneSegment(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()

	seed := func(t *testing.T) {
		t.Helper()
		if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "LegacyTable",
			"Item": map[string]interface{}{
				"id":          map[string]interface{}{"S": "A"},
				"sk":          map[string]interface{}{"S": "1"},
				"dotted.name": map[string]interface{}{"S": "top-level"},
				"dotted":      map[string]interface{}{"M": map[string]interface{}{"name": map[string]interface{}{"S": "nested"}}},
			},
		}}); err != nil {
			t.Fatalf("seed item: %v", err)
		}
	}
	key := map[string]interface{}{"id": map[string]interface{}{"S": "A"}, "sk": map[string]interface{}{"S": "1"}}

	t.Run("whole-token alias", func(t *testing.T) {
		seed(t)
		resp, err := svc.GetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":                "LegacyTable",
			"Key":                      key,
			"ProjectionExpression":     "#d",
			"ExpressionAttributeNames": map[string]interface{}{"#d": "dotted.name"},
		}})
		if err != nil {
			t.Fatalf("GetItem: %v", err)
		}
		item := resp.(map[string]interface{})["Item"].(map[string]interface{})
		if av, ok := item["dotted.name"].(map[string]interface{}); !ok || av["S"] != "top-level" {
			t.Fatalf("projected item = %v, want the top-level attribute dotted.name", item)
		}
		if _, hasNested := item["dotted"]; hasNested {
			t.Fatalf("projected item = %v, the alias must not re-split into the nested path", item)
		}
	})

	t.Run("compound alias still descends", func(t *testing.T) {
		seed(t)
		resp, err := svc.GetItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":            "LegacyTable",
			"Key":                  key,
			"ProjectionExpression": "#p.#n",
			"ExpressionAttributeNames": map[string]interface{}{
				"#p": "dotted",
				"#n": "name",
			},
		}})
		if err != nil {
			t.Fatalf("GetItem: %v", err)
		}
		item := resp.(map[string]interface{})["Item"].(map[string]interface{})
		outer, ok := item["dotted"].(map[string]interface{})
		if !ok {
			t.Fatalf("projected item = %v, want the nested map under dotted", item)
		}
		if av, ok := outer["M"].(map[string]interface{}); !ok {
			t.Fatalf("dotted = %v, want the M-shaped nested value", outer)
		} else if inner, ok := av["name"].(map[string]interface{}); !ok || inner["S"] != "nested" {
			t.Fatalf("dotted.name = %v, want the nested value", av)
		}
	})
}

// The three string enums the table cores persist must be validated before
// they reach the store: a raw BillingMode/TableClass/StreamViewType cast
// would otherwise persist an out-of-enum value forever.
func TestTablePlaneEnumValidation(t *testing.T) {
	if validateBillingModeValue("FOO") {
		t.Error("invalid billing mode accepted")
	}
	if !validateBillingModeValue(dbstore.BillingModePayPerRequest) {
		t.Error("PAY_PER_REQUEST rejected")
	}
	if validateTableClassValue("STANDARD_X") {
		t.Error("invalid table class accepted")
	}
	if !validateTableClassValue("STANDARD_INFREQUENT_ACCESS") {
		t.Error("STANDARD_INFREQUENT_ACCESS rejected")
	}
	badView := &dbstore.StreamSpecification{StreamEnabled: true, StreamViewType: dbstore.StreamViewType("ALL")}
	if err := validateStreamSpecification(badView); err == nil {
		t.Error("invalid stream view type accepted")
	}
	if err := validateStreamSpecification(&dbstore.StreamSpecification{StreamEnabled: true, StreamViewType: dbstore.StreamViewTypeKeysOnly}); err != nil {
		t.Errorf("KEYS_ONLY rejected: %v", err)
	}
	if err := validateStreamSpecification(&dbstore.StreamSpecification{StreamEnabled: false}); err != nil {
		t.Errorf("disabled specification rejected: %v", err)
	}
}

// DescribeLimits reports the quota constants; the throughput validator
// enforces the same constants on every create and update path.
func TestProvisionedThroughputQuotaBounds(t *testing.T) {
	over := &dbstore.ProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: dbstore.TableMaxWriteCapacityUnits + 1}
	if validateProvisionedThroughputValues(over) {
		t.Error("write units above the table quota accepted")
	}
	atBound := &dbstore.ProvisionedThroughput{ReadCapacityUnits: dbstore.TableMaxReadCapacityUnits, WriteCapacityUnits: 1}
	if !validateProvisionedThroughputValues(atBound) {
		t.Error("read units at the table quota rejected")
	}
}

// DescribeTable's Replicas member reports the global table's replication
// group; a standalone table reports no replicas.
func TestBuildTableDescriptionRendersReplicas(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)
	svc := &DynamoDBService{}
	table, err := store.Tables().Get("ImpTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}

	desc := svc.buildTableDescription(table, svc.replicasForTable(store, table.Name))
	if _, ok := desc["Replicas"]; ok {
		t.Fatalf("standalone table replicas = %v, want the member absent", desc["Replicas"])
	}

	if _, err := store.GlobalTables().Create("ImpTbl", []*dbstore.Replica{
		{RegionName: "us-east-1", ReplicaStatus: "ACTIVE"},
		{RegionName: "eu-west-1", ReplicaStatus: "ACTIVE"},
	}); err != nil {
		t.Fatalf("create global table: %v", err)
	}
	desc = svc.buildTableDescription(table, svc.replicasForTable(store, table.Name))
	replicas, ok := desc["Replicas"].([]interface{})
	if !ok || len(replicas) != 2 {
		t.Fatalf("global table replicas = %v, want two entries", desc["Replicas"])
	}
	first := replicas[0].(map[string]interface{})
	if first["RegionName"] != "us-east-1" || first["ReplicaStatus"] != "ACTIVE" {
		t.Fatalf("replica entry = %v", first)
	}
}

// DescribeTimeToLive surfaces the stored TTL state machine value, so a
// freshly-enabled table reports ENABLING until the transition commits.
func TestDescribeTimeToLiveSurfacesStoredStatus(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)
	if err := store.Tables().SetTimeToLive("ImpTbl", &dbstore.TimeToLiveSpecification{
		Enabled:       true,
		AttributeName: "exp",
		Status:        dbstore.TTLStatusEnabling,
	}); err != nil {
		t.Fatalf("set ttl: %v", err)
	}
	ttl, err := store.Tables().GetTimeToLive("ImpTbl")
	if err != nil || ttl == nil {
		t.Fatalf("get ttl: %v %v", ttl, err)
	}
	status := "DISABLED"
	if ttl.Status != "" {
		status = string(ttl.Status)
	} else if ttl.Enabled {
		status = "ENABLED"
	}
	if status != "ENABLING" {
		t.Fatalf("derived status = %s, want ENABLING from the stored state machine", status)
	}
}

// TestTransactItemReadUnits pins the per-item transactional read charge:
// twice the strongly consistent rate (one read to prepare the transaction,
// one to commit it) over the item's size rounded up to 4 KB multiples, with
// a one-unit minimum that also covers absent items.
func TestTransactItemReadUnits(t *testing.T) {
	cases := map[int64]float64{
		0:     2,
		1:     2,
		4096:  2,
		4097:  4,
		8192:  4,
		12288: 6,
	}
	for size, want := range cases {
		if got := transactItemReadUnits(size); got != want {
			t.Errorf("size %d: units = %v, want %v", size, got, want)
		}
	}
}

// TestTransactItemWriteUnits pins the per-item transactional write charge:
// twice the standard rate (one write to prepare the transaction, one to
// commit it) over the item's size rounded up to 1 KB multiples, with a
// one-unit minimum that also covers deleting an absent item.
func TestTransactItemWriteUnits(t *testing.T) {
	cases := map[int64]float64{
		0:    2,
		1:    2,
		1024: 2,
		1025: 4,
		2048: 4,
		5120: 10,
	}
	for size, want := range cases {
		if got := transactItemWriteUnits(size); got != want {
			t.Errorf("size %d: units = %v, want %v", size, got, want)
		}
	}
}

// TestBuildTransactWriteResponsePerItemCapacity pins the per-table capacity
// aggregation of a committed transaction: each written item contributes its
// own doubled 1 KB-rounded charge, a condition check contributes the doubled
// 4 KB-rounded read of the checked item, and a table touched by both kinds
// reports both member fields under their sum. A replay reports the recorded
// per-table read units of the transaction's items — the same items, sized
// as reads — in TransactItems first-appearance order.
func TestBuildTransactWriteResponsePerItemCapacity(t *testing.T) {
	bigValue := strings.Repeat("x", 2500)
	ops := []writeOperation{
		{opType: "Put", tableName: "T1", itemData: map[string]*dbstore.AttributeValue{
			"pk": dbstore.StringValue("a"), "body": dbstore.StringValue(bigValue),
		}},
		{opType: "Delete", tableName: "T1", itemSize: 8 * 1024},
		{opType: "ConditionCheck", tableName: "T2", itemSize: 5000},
	}
	params := map[string]interface{}{"ReturnConsumedCapacity": "TOTAL"}

	resp := buildTransactWriteResponse(params, ops, false, nil)
	caps, ok := resp["ConsumedCapacity"].([]map[string]interface{})
	if !ok || len(caps) != 2 {
		t.Fatalf("consumed capacities = %v", resp["ConsumedCapacity"])
	}
	if caps[0]["TableName"] != "T1" || caps[1]["TableName"] != "T2" {
		t.Fatalf("capacities not in TransactItems first-appearance order: %v", caps)
	}
	byTable := map[string]map[string]interface{}{}
	for _, c := range caps {
		byTable[c["TableName"].(string)] = c
	}

	// T1: the Put item sizes 2507 bytes (3 KB rounded, 6 units) and the
	// deleted item 8 KB (8 KB at write granularity, 16 units) — 22 write
	// units, no read units.
	t1 := byTable["T1"]
	if t1["WriteCapacityUnits"].(float64) != 22 {
		t.Fatalf("T1 write units = %v, want 22", t1["WriteCapacityUnits"])
	}
	if _, present := t1["ReadCapacityUnits"]; present {
		t.Fatalf("T1 must not report read units: %v", t1)
	}
	if t1["CapacityUnits"].(float64) != 22 {
		t.Fatalf("T1 capacity units = %v, want 22", t1["CapacityUnits"])
	}

	// T2: the condition-checked item sizes 5000 bytes (2 units) at read
	// granularity — 4 read units, no write units.
	t2 := byTable["T2"]
	if t2["ReadCapacityUnits"].(float64) != 4 {
		t.Fatalf("T2 read units = %v, want 4", t2["ReadCapacityUnits"])
	}
	if _, present := t2["WriteCapacityUnits"]; present {
		t.Fatalf("T2 must not report write units: %v", t2)
	}

	replay := buildTransactWriteResponse(params, ops, true, transactReplayReadUnits(ops))
	replayCaps, ok := replay["ConsumedCapacity"].([]map[string]interface{})
	if !ok || len(replayCaps) != 2 {
		t.Fatalf("replay capacities = %v", replay["ConsumedCapacity"])
	}
	if replayCaps[0]["TableName"] != "T1" || replayCaps[1]["TableName"] != "T2" {
		t.Fatalf("replay capacities not in TransactItems first-appearance order: %v", replayCaps)
	}
	// T1 replay: the Put item (2507 bytes) reads as one doubled unit and
	// the deleted item (8 KB) as two doubled units — 6 read units.
	if got := replayCaps[0]["ReadCapacityUnits"].(float64); got != 6 {
		t.Fatalf("T1 replay read units = %v, want 6", got)
	}
	if got := replayCaps[0]["CapacityUnits"].(float64); got != 6 {
		t.Fatalf("T1 replay capacity units = %v, want 6", got)
	}
	// T2 replay: the checked item (5000 bytes) reads as two doubled units.
	if got := replayCaps[1]["ReadCapacityUnits"].(float64); got != 4 {
		t.Fatalf("T2 replay read units = %v, want 4", got)
	}
	if got := replayCaps[1]["CapacityUnits"].(float64); got != 4 {
		t.Fatalf("T2 replay capacity units = %v, want 4", got)
	}
	for _, c := range replayCaps {
		if _, present := c["WriteCapacityUnits"]; present {
			t.Fatalf("replay must not report write units: %v", c)
		}
	}
}

// TestTransactReplayReadUnits pins the replay's read sizing: every item of
// the transaction is read at the doubled 4 KB granularity — an 8 KB
// single-item transaction replays with 4 read units, matching the
// documented "read capacity units consumed in reading the item" — and an
// absent delete target still reads the minimum doubled unit, because the
// two underlying reads happen regardless.
func TestTransactReplayReadUnits(t *testing.T) {
	// body of 8185 bytes plus attribute overhead sizes the item at 8192.
	ops := []writeOperation{
		{opType: "Put", tableName: "T", itemData: map[string]*dbstore.AttributeValue{
			"pk": dbstore.StringValue("a"), "body": dbstore.StringValue(strings.Repeat("x", 8185)),
		}},
	}
	units := transactReplayReadUnits(ops)
	if len(units) != 1 || units["T"] != 4 {
		t.Fatalf("replay read units = %v, want T:4", units)
	}

	units = transactReplayReadUnits([]writeOperation{{opType: "Delete", tableName: "T"}})
	if units["T"] != 2 {
		t.Fatalf("absent delete replay read units = %v, want T:2", units)
	}
}

// streamPlaneFixture stands up one streamed table with two committed item
// writes, the state every Streams read pin below starts from.
func streamPlaneFixture(t *testing.T) (*DynamoDBService, dbstore.DynamoDBStoreInterface, *request.RequestContext, *dbstore.Table) {
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
	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "StreamsPinTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
		StreamSpecification:  &dbstore.StreamSpecification{StreamEnabled: true, StreamViewType: dbstore.StreamViewTypeNewAndOldImages},
	})
	if err != nil {
		t.Fatalf("create streamed table: %v", err)
	}
	reqCtx := request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
	for _, id := range []string{"a", "b"} {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "StreamsPinTable",
			"Item":      map[string]interface{}{"id": map[string]interface{}{"S": id}},
		}}); err != nil {
			t.Fatalf("seed put %s: %v", id, err)
		}
	}
	return svc, store, reqCtx, table
}

// streamShardID reads the stream's single shard id the way a client does:
// from DescribeStream's response.
func streamShardID(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, streamArn string) string {
	t.Helper()
	resp, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"StreamArn": streamArn}})
	if err != nil {
		t.Fatalf("describe stream: %v", err)
	}
	shards := resp.(map[string]interface{})["StreamDescription"].(map[string]interface{})["Shards"].([]interface{})
	if len(shards) != 1 {
		t.Fatalf("expected the single shard, got %d", len(shards))
	}
	return shards[0].(map[string]interface{})["ShardId"].(string)
}

// TestDescribeStreamValidatesAndHonoursRequestMembers pins the
// DescribeStream request members against the model: the Limit range rejects
// an explicit below-one value, the ExclusiveStartShardId length bounds hold
// and the cursor positions the page after the matching shard, and the
// ShardFilter accepts the enum's lone CHILD_SHARDS value — whose answer over
// a never-splitting single shard is an empty list — while rejecting any
// other type. The shard's sequence numbers ride the wire in the model's
// 21-40 character form.

// TestGetShardIteratorValidatesShardIdAndSequenceNumber pins the
// GetShardIterator members against the model: a ShardId outside the 28-65
// length bounds is a validation error and a well-formed id that is not the
// stream's shard answers ResourceNotFound; the SequenceNumber length is the
// model's 21-40 and the padded wire form round-trips through
// AT/AFTER_SEQUENCE_NUMBER; and a position whose first record the retention
// window already removed is the documented TrimmedDataAccessException at
// iterator-issuing time.

// streamsFaultStore wraps a real store so that every Update hands its
// callback a transaction whose buckets all delegate to the carrying
// transaction except the streams bucket, whose writes fail — the
// storage-fault class the stream record write can hit inside the item
// write's own transaction.
type streamsFaultStore struct {
	dbstore.DynamoDBStoreInterface
}

func (w streamsFaultStore) Update(ctx context.Context, fn func(txn *dbstore.DynamoDBTxn) error) error {
	return w.DynamoDBStoreInterface.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		return fn(w.NewTxn(faultingStreamsTxn{txn.RawTxn()}))
	})
}

type faultingStreamsTxn struct {
	storage.Transaction
}

func (t faultingStreamsTxn) Bucket(name string) storage.Bucket {
	if name == "dynamodb_streams-us-east-1" {
		return faultingStreamsBucket{}
	}
	return t.Transaction.Bucket(name)
}

type faultingStreamsBucket struct {
	storage.Bucket
}

func (faultingStreamsBucket) Put(key, value []byte) error {
	return errors.New("injected stream-record write failure")
}

// TestStreamRecordWriteFailureFailsTheItemWrite pins the capture
// atomicity: when the stream record write fails inside the carrying
// transaction, the item write itself fails and nothing commits — a streamed
// table's item change never lands without the stream record its
// StreamSpecification promises.

// TestShardIteratorIsBoundToItsStreamGeneration pins the iterator's stream
// identity and the generation boundary: UpdateTable regenerates the stream
// ARN on every re-enable, an iterator issued for the superseded generation —
// whose signature, table, and position are all still valid — must answer
// ResourceNotFound against the successor stream, exactly as an iterator of a
// stream that no longer exists does, and the re-enabled stream is a new
// generation that starts empty — its record space carries none of the
// superseded generation's records, and its own writes read back under its
// own ARN.

// TestListStreamsRejectsBelowOneLimit pins the ListStreams Limit range: an
// explicit value below the model's minimum of 1 is a ValidationException,
// while an omitted Limit takes the default page.
