package dynamodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestIsValidDynamoDBNumber(t *testing.T) {
	valid := []string{
		"0",
		"-1.5",
		"+3",
		"1e10",
		"9.9999999999999999999999999999999999999E+125",
		"1E-130",
		"0.000000000000000000000000000000000000000001",
		"10000000000000000000000000000000000000000",
	}
	for _, value := range valid {
		assert.True(t, isValidDynamoDBNumber(value), "expected %q to be valid", value)
	}

	invalid := []string{
		"",
		"abc",
		"1/3",     // big.Rat fraction form, not part of the grammar
		"1e400",   // magnitude above 9.99...E+125
		"-1E-131", // magnitude below 1E-130
		"1.23456789012345678901234567890123456789", // 39 significant digits
		"1.2.3",
		"--1",
	}
	for _, value := range invalid {
		assert.False(t, isValidDynamoDBNumber(value), "expected %q to be invalid", value)
	}
}

func TestCountSignificantDigits(t *testing.T) {
	assert.Equal(t, 0, dbstore.CountSignificantDigits("0"))
	assert.Equal(t, 0, dbstore.CountSignificantDigits("0.000"))
	assert.Equal(t, 1, dbstore.CountSignificantDigits("1000"))
	assert.Equal(t, 3, dbstore.CountSignificantDigits("1.23"))
	assert.Equal(t, 1, dbstore.CountSignificantDigits("0.0001"))
	assert.Equal(t, 2, dbstore.CountSignificantDigits("-2.5e10"))
}

func TestSetDuplicateDetection(t *testing.T) {
	assert.True(t, hasDuplicateString([]string{"a", "b", "a"}))
	assert.False(t, hasDuplicateString([]string{"a", "b", ""}))
	assert.True(t, hasDuplicateNumber([]string{"1", "2", "1.0"}))
	assert.False(t, hasDuplicateNumber([]string{"1", "2", "3"}))
	assert.True(t, hasDuplicateBinary([][]byte{[]byte("x"), []byte("x")}))
	assert.False(t, hasDuplicateBinary([][]byte{[]byte("x"), []byte("")}))
}

func TestResolveNameStrict(t *testing.T) {
	resolved, err := resolveNameStrict("#n", map[string]string{"#n": "name"})
	require.NoError(t, err)
	assert.Equal(t, "name", resolved)

	resolved, err = resolveNameStrict("plain", nil)
	require.NoError(t, err)
	assert.Equal(t, "plain", resolved)

	_, err = resolveNameStrict("#missing", map[string]string{"#n": "name"})
	assert.Error(t, err)

	_, err = resolveNameStrict("#missing", nil)
	assert.Error(t, err)
}

func TestUpdateExpressionStrictness(t *testing.T) {
	values := map[string]*dbstore.AttributeValue{
		":v": dbstore.NumberValue("2"),
		":w": dbstore.NumberValue("3"),
	}

	attrs := map[string]*dbstore.AttributeValue{
		"src": dbstore.StringValue("copied"),
	}

	// The multiply operator is not part of SET arithmetic.
	_, err := applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :v * :w", nil, values)
	assert.Error(t, err)

	// An undefined value placeholder must not silently skip the action.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :missing", nil, values)
	assert.Error(t, err)

	// An undefined name placeholder must not pass through as a literal.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET #missing = :v", nil, values)
	assert.Error(t, err)

	// SET a = src copies the existing attribute value.
	updated, err := applyUpdateExpressionWithTracking(attrs, "SET a = src", nil, nil)
	require.NoError(t, err)
	assert.Contains(t, updated, "a")
	assert.Equal(t, "copied", *attrs["a"].S)

	// DELETE with a mismatched operand type is rejected.
	ssAttrs := map[string]*dbstore.AttributeValue{
		"tags": dbstore.StringSet([]string{"a", "b"}),
	}
	_, err = applyUpdateExpressionWithTracking(ssAttrs, "DELETE tags :nums", nil, map[string]*dbstore.AttributeValue{
		":nums": dbstore.NumberSet([]string{"1"}),
	})
	assert.Error(t, err)

	// DELETE of matching elements still works.
	updated, err = applyUpdateExpressionWithTracking(ssAttrs, "DELETE tags :del", nil, map[string]*dbstore.AttributeValue{
		":del": dbstore.StringSet([]string{"a"}),
	})
	require.NoError(t, err)
	assert.Contains(t, updated, "tags")
	assert.Equal(t, []string{"b"}, ssAttrs["tags"].SS)

	// An undefined name placeholder in REMOVE must not pass through as a
	// literal attribute name.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "REMOVE #missing", nil, nil)
	assert.Error(t, err)

	// REMOVE with a defined name placeholder removes the attribute.
	rmAttrs := map[string]*dbstore.AttributeValue{
		"gsik": dbstore.StringValue("g"),
	}
	updated, err = applyUpdateExpressionWithTracking(rmAttrs, "REMOVE #name", map[string]string{"#name": "gsik"}, nil)
	require.NoError(t, err)
	assert.Contains(t, updated, "gsik")
	assert.NotContains(t, rmAttrs, "gsik")
}

// TestAdminItemCoresRejectEmptyTableName pins the empty-table rejection of
// the admin item operations: an omitted TableName is a client error
// (ValidationException) owned by describeTableCore's name validation, which
// both planes share.
func TestAdminItemCoresRejectEmptyTableName(t *testing.T) {
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	defer sm.Close()
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)

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

// TestItemWritesRequireActiveTable pins the write-requires-ACTIVE rule on
// both planes: a table persisted in CREATING state — what a restore holds
// while copying items — rejects every item write, HTTP and admin console
// alike, while reads resolve the same table on both planes.
func TestItemWritesRequireActiveTable(t *testing.T) {
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

	reqCtx := request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
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

// TestBatchTransactionAndPartiQLWritesRequireActiveTable extends the
// write-requires-ACTIVE rule to every batched write surface: BatchWriteItem,
// TransactWriteItems and each PartiQL write engine (single-statement and
// transactional) reject a table persisted in CREATING state with the same
// error as the single-item writes, while the batched reads stay
// status-independent like GetItem and Scan.
func TestBatchTransactionAndPartiQLWritesRequireActiveTable(t *testing.T) {
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

	reqCtx := request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
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
