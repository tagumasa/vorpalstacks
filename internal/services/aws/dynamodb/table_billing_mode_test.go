// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// billingModePlaneFixture builds the single-region table plane: a service
// over a temp-dir storage manager and a request context for its region.
func billingModePlaneFixture(t *testing.T) (*DynamoDBService, *request.RequestContext) {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	return svc, request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
}

func describeTableOnModePlane(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, name string) map[string]interface{} {
	t.Helper()
	resp, err := svc.DescribeTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": name}})
	if err != nil {
		t.Fatalf("describe %s: %v", name, err)
	}
	return resp.(map[string]interface{})["Table"].(map[string]interface{})
}

// gsiCapacity reads a described table's lone GSI capacity pair. The stored
// values render as int64 while an absent pair renders as zero literals, so
// the read normalises both.
func gsiCapacity(t *testing.T, table map[string]interface{}, indexName string) (int64, int64) {
	t.Helper()
	gsis := table["GlobalSecondaryIndexes"].([]map[string]interface{})
	for _, gsi := range gsis {
		if gsi["IndexName"] != indexName {
			continue
		}
		pt := gsi["ProvisionedThroughput"].(map[string]interface{})
		return asInt64(t, pt["ReadCapacityUnits"]), asInt64(t, pt["WriteCapacityUnits"])
	}
	t.Fatalf("index %s not found among %v", indexName, gsis)
	return 0, 0
}

func asInt64(t *testing.T, v interface{}) int64 {
	t.Helper()
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	}
	t.Fatalf("capacity value type %T=%v", v, v)
	return 0
}

// TestGSIThroughputFollowsTheTableBillingMode pins the inheritance rule the
// capacity-mode switch owes: switching the table to on-demand clears the
// indexes' provisioned settings along with the table's own (DescribeTable
// renders zeros, not the pre-switch units), an index update may not attach
// provisioned settings to an on-demand table's index, switching back
// requires the initial provisioned values for the table and its indexes,
// and creation holds the same settings-per-mode contract in both directions.
func TestGSIThroughputFollowsTheTableBillingMode(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	name := "GsiModeTable"

	_, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            name,
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}, map[string]interface{}{"AttributeName": "gpk", "AttributeType": "S"}},
		"BillingMode":          "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{
			"ReadCapacityUnits":  5.0,
			"WriteCapacityUnits": 7.0,
		},
		"GlobalSecondaryIndexes": []interface{}{map[string]interface{}{
			"IndexName":             "gsi",
			"KeySchema":             []interface{}{map[string]interface{}{"AttributeName": "gpk", "KeyType": "HASH"}},
			"Projection":            map[string]interface{}{"ProjectionType": "ALL"},
			"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 3.0, "WriteCapacityUnits": 4.0},
		}},
	}})
	if err != nil {
		t.Fatalf("create provisioned table with GSI: %v", err)
	}

	if rcu, wcu := gsiCapacity(t, describeTableOnModePlane(t, svc, reqCtx, name), "gsi"); rcu != 3 || wcu != 4 {
		t.Fatalf("created GSI capacity: got %d/%d, want 3/4", rcu, wcu)
	}

	update := func(params map[string]interface{}) (map[string]interface{}, error) {
		resp, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: params})
		if err != nil {
			return nil, err
		}
		return resp.(map[string]interface{})["TableDescription"].(map[string]interface{}), nil
	}

	// The switch to on-demand takes the index's provisioned settings with
	// the table's own; the response and a fresh describe both render zeros.
	if _, err := update(map[string]interface{}{"TableName": name, "BillingMode": "PAY_PER_REQUEST"}); err != nil {
		t.Fatalf("switch to on-demand: %v", err)
	}
	desc := describeTableOnModePlane(t, svc, reqCtx, name)
	if mode := desc["BillingModeSummary"].(map[string]interface{})["BillingMode"]; mode != "PAY_PER_REQUEST" {
		t.Fatalf("billing mode after switch: %v", mode)
	}
	if rcu, wcu := gsiCapacity(t, desc, "gsi"); rcu != 0 || wcu != 0 {
		t.Fatalf("GSI capacity on the on-demand table: got %d/%d, want 0/0", rcu, wcu)
	}
	if pt := desc["ProvisionedThroughput"].(map[string]interface{}); asInt64(t, pt["ReadCapacityUnits"]) != 0 || asInt64(t, pt["WriteCapacityUnits"]) != 0 {
		t.Fatalf("table capacity on the on-demand table: %v", pt)
	}

	// An index update may not attach provisioned settings to the on-demand
	// table's index — the settings-per-mode contract the switch just used.
	if _, err := update(map[string]interface{}{
		"TableName": name,
		"GlobalSecondaryIndexUpdates": []interface{}{map[string]interface{}{
			"Update": map[string]interface{}{
				"IndexName":             "gsi",
				"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 1.0, "WriteCapacityUnits": 1.0},
			},
		}},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("attach provisioned settings to an on-demand index: expected ErrInvalidParameter, got %v", err)
	}

	// Switching back requires the initial provisioned values for the table
	// and its indexes: the table-level pair alone does not restore the
	// index's settings, so the switch is refused.
	if _, err := update(map[string]interface{}{
		"TableName":             name,
		"BillingMode":           "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 5.0},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("switch back without index capacity: expected ErrInvalidParameter, got %v", err)
	}

	// With the index's settings riding the same request, the switch lands
	// and both levels report their new units.
	if _, err := update(map[string]interface{}{
		"TableName":             name,
		"BillingMode":           "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 5.0},
		"GlobalSecondaryIndexUpdates": []interface{}{map[string]interface{}{
			"Update": map[string]interface{}{
				"IndexName":             "gsi",
				"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 8.0, "WriteCapacityUnits": 9.0},
			},
		}},
	}); err != nil {
		t.Fatalf("switch back with index capacity: %v", err)
	}
	desc = describeTableOnModePlane(t, svc, reqCtx, name)
	if mode := desc["BillingModeSummary"].(map[string]interface{})["BillingMode"]; mode != "PROVISIONED" {
		t.Fatalf("billing mode after switch back: %v", mode)
	}
	if rcu, wcu := gsiCapacity(t, desc, "gsi"); rcu != 8 || wcu != 9 {
		t.Fatalf("GSI capacity after switch back: got %d/%d, want 8/9", rcu, wcu)
	}
	if pt := desc["ProvisionedThroughput"].(map[string]interface{}); asInt64(t, pt["ReadCapacityUnits"]) != 5 || asInt64(t, pt["WriteCapacityUnits"]) != 5 {
		t.Fatalf("table capacity after switch back: %v", pt)
	}

	// Creation holds the same contract in both directions: an on-demand
	// table's index carries no provisioned settings, a provisioned table's
	// index must carry them.
	_, err = svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "GsiOndemandCreate",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}, map[string]interface{}{"AttributeName": "gpk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
		"GlobalSecondaryIndexes": []interface{}{map[string]interface{}{
			"IndexName":             "gsi",
			"KeySchema":             []interface{}{map[string]interface{}{"AttributeName": "gpk", "KeyType": "HASH"}},
			"Projection":            map[string]interface{}{"ProjectionType": "ALL"},
			"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 3.0, "WriteCapacityUnits": 4.0},
		}},
	}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("on-demand creation with provisioned index settings: expected ErrInvalidParameter, got %v", err)
	}

	_, err = svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":             "GsiProvCreate",
		"KeySchema":             []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions":  []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}, map[string]interface{}{"AttributeName": "gpk", "AttributeType": "S"}},
		"BillingMode":           "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 5.0},
		"GlobalSecondaryIndexes": []interface{}{map[string]interface{}{
			"IndexName":  "gsi",
			"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "gpk", "KeyType": "HASH"}},
			"Projection": map[string]interface{}{"ProjectionType": "ALL"},
		}},
	}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("provisioned creation without index capacity: expected ErrInvalidParameter, got %v", err)
	}
}

// TestBillingModeSwitchBackRequiresTableCapacity pins the table-level half
// of the switch-back rule the model documents ("When switching from
// pay-per-request to provisioned capacity, initial provisioned capacity
// values must be set"): an on-demand table cannot switch back on the billing
// mode member alone.
func TestBillingModeSwitchBackRequiresTableCapacity(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "PlainOndemandTable",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create on-demand table: %v", err)
	}

	_, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":   "PlainOndemandTable",
		"BillingMode": "PROVISIONED",
	}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("switch back on the mode member alone: expected ErrInvalidParameter, got %v", err)
	}
}

// TestBillingModeSwitchRateWindow pins the documented switch quota: four
// provisioned-to-on-demand switches inside a rolling 24-hour window pass,
// the fifth is refused, and the on-demand-to-provisioned direction stays
// unlimited ("You can switch tables from on-demand mode to provisioned
// capacity mode at any time") — the interleaved sequence exercises exactly
// that asymmetry.
func TestBillingModeSwitchRateWindow(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	name := "SwitchRateTable"

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":             name,
		"KeySchema":             []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions":  []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":           "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 5.0},
	}}); err != nil {
		t.Fatalf("create provisioned table: %v", err)
	}

	switchMode := func(mode string, pt map[string]interface{}) error {
		params := map[string]interface{}{"TableName": name, "BillingMode": mode}
		if pt != nil {
			params["ProvisionedThroughput"] = pt
		}
		_, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}
	provisioned := map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 5.0}

	for i := 0; i < dbstore.BillingModeSwitchesPerDay; i++ {
		if err := switchMode("PAY_PER_REQUEST", nil); err != nil {
			t.Fatalf("switch %d to on-demand: %v", i+1, err)
		}
		if err := switchMode("PROVISIONED", provisioned); err != nil {
			t.Fatalf("switch %d back to provisioned: %v", i+1, err)
		}
	}

	if err := switchMode("PAY_PER_REQUEST", nil); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("fifth switch inside the window: expected ErrInvalidParameter, got %v", err)
	}
}

// TestTableCreationBillingModeContractAtEveryEntryPoint pins the billing-mode
// contract on the two creation entry points that used to bypass it. The
// import plane's omitted mode resolves to the documented provisioned default
// and a provisioned request without a throughput pair is a validation error
// (no invented inline pair). The admin plane's mode member carries explicit
// presence (a non-required enum is emitted optional): an omitted member
// routes through the same empty value the HTTP plane's omission produces —
// the shared contract's provisioned default, identical outcomes on both
// planes — while the shared contract governs the pair against the named
// mode in both directions: a provisioned request without one is rejected,
// and an on-demand request carrying one is rejected instead of silently
// discarding it.
func TestTableCreationBillingModeContractAtEveryEntryPoint(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	adminReq := func(name string, mode *pb.BillingMode, pt *pb.ProvisionedThroughput) *pb.CreateTableInput {
		return &pb.CreateTableInput{
			Tablename:             name,
			Keyschema:             []*pb.KeySchemaElement{{Attributename: "pk", Keytype: pb.KeyType_KEY_TYPE_HASH}},
			Attributedefinitions:  []*pb.AttributeDefinition{{Attributename: "pk", Attributetype: pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_S}},
			Billingmode:           mode,
			Provisionedthroughput: pt,
		}
	}
	mode := func(m pb.BillingMode) *pb.BillingMode { return &m }

	// The omitted admin member yields the provisioned default, identical to
	// the HTTP plane: without a pair the provisioned contract rejects, with
	// a pair the created table is provisioned.
	if _, err := svc.adminCreateTable(ctx, reqCtx.GetRegion(), adminReq("AdmBillOmitA", nil, nil)); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("admin omitted mode without a pair: err = %v, want ErrInvalidParameter", err)
	}
	omitDesc, err := svc.adminCreateTable(ctx, reqCtx.GetRegion(), adminReq("AdmBillOmitB", nil, &pb.ProvisionedThroughput{Readcapacityunits: 5, Writecapacityunits: 5}))
	if err != nil {
		t.Fatalf("admin omitted mode with a pair: %v", err)
	}
	if got := omitDesc.GetBillingmodesummary().GetBillingmode(); got != pb.BillingMode_BILLING_MODE_PROVISIONED {
		t.Fatalf("admin omitted mode summary = %v, want PROVISIONED", got)
	}

	// Admin surface: the shared contract governs the pair against the named
	// mode in both directions.
	if _, err := svc.adminCreateTable(context.Background(), reqCtx.GetRegion(), adminReq("AdmBillA", mode(pb.BillingMode_BILLING_MODE_PROVISIONED), nil)); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("admin provisioned without a pair: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := svc.adminCreateTable(context.Background(), reqCtx.GetRegion(), adminReq("AdmBillB", mode(pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST), &pb.ProvisionedThroughput{Readcapacityunits: 5, Writecapacityunits: 5})); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("admin on-demand with a pair: err = %v, want ErrInvalidParameter", err)
	}
	desc, err := svc.adminCreateTable(context.Background(), reqCtx.GetRegion(), adminReq("AdmBillC", mode(pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST), nil))
	if err != nil {
		t.Fatalf("admin on-demand create: %v", err)
	}
	if got := desc.GetBillingmodesummary().GetBillingmode(); got != pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST {
		t.Fatalf("admin on-demand summary = %v, want PAY_PER_REQUEST", got)
	}

	// Import surface: the omitted mode resolves to the same provisioned
	// default, and a provisioned request without a pair is rejected — the
	// previous behaviour silently created an on-demand table (omitted mode)
	// or invented an inline 5/5 pair (provisioned without a pair).
	importReq := func(table string, creationExtra map[string]interface{}) importTableInput {
		tcp := map[string]interface{}{
			"TableName":            table,
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		}
		for k, v := range creationExtra {
			tcp[k] = v
		}
		return importTableInput{Parameters: map[string]interface{}{
			"S3BucketSource":          map[string]interface{}{"S3Bucket": "import-bucket"},
			"InputFormat":             "DYNAMODB_JSON",
			"TableCreationParameters": tcp,
		}}
	}
	if _, err := svc.importTableCore(ctx, reqCtx, importReq("ImpBillA", nil)); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("import omitted mode without a pair: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := svc.importTableCore(ctx, reqCtx, importReq("ImpBillB", map[string]interface{}{"BillingMode": "PROVISIONED"})); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("import provisioned without a pair: err = %v, want ErrInvalidParameter", err)
	}

	// A provisioned import that carries its pair is accepted, and the job's
	// created table carries the request's pair — not an invented one. The
	// job reaches the S3 plane and terminates FAILED there (a bare service
	// has no S3 invoker), after the table-creation step.
	resp, err := svc.importTableCore(ctx, reqCtx, importReq("ImpBillC", map[string]interface{}{
		"BillingMode":           "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 7.0, "WriteCapacityUnits": 3.0},
	}))
	if err != nil {
		t.Fatalf("import provisioned with a pair: %v", err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		imp, getErr := store.Imports().Get(resp.Import.ImportArn)
		if getErr != nil {
			t.Fatalf("reload import: %v", getErr)
		}
		if imp.ImportStatus != "IN_PROGRESS" {
			if imp.ImportStatus != "FAILED" {
				t.Fatalf("import job = %s/%s, want FAILED", imp.ImportStatus, imp.FailureCode)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("import job never left IN_PROGRESS")
		}
		time.Sleep(20 * time.Millisecond)
	}
	created, err := store.Tables().Get("ImpBillC")
	if err != nil {
		t.Fatalf("read import-created table: %v", err)
	}
	if created.BillingMode != dbstore.BillingModeProvisioned ||
		created.ProvisionedThroughput == nil ||
		created.ProvisionedThroughput.ReadCapacityUnits != 7 ||
		created.ProvisionedThroughput.WriteCapacityUnits != 3 {
		t.Fatalf("import-created table mode/pair = %v/%+v, want PROVISIONED 7/3", created.BillingMode, created.ProvisionedThroughput)
	}
}

func TestBillingModeRollingSwitchWindow(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "RollTbl",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "pk", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModeProvisioned,
		ProvisionedThroughput: &dbstore.ProvisionedThroughput{
			ReadCapacityUnits:  5,
			WriteCapacityUnits: 5,
		},
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	switchMode := func() error {
		_, err := svc.UpdateTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":   "RollTbl",
			"BillingMode": "PAY_PER_REQUEST",
		}})
		return err
	}
	// Fresh stamps at the cap refuse the switch.
	if _, err := store.Tables().Update("RollTbl", func(table *dbstore.Table) error {
		table.BillingModeSwitches = make([]time.Time, dbstore.BillingModeSwitchesPerDay)
		for i := range table.BillingModeSwitches {
			table.BillingModeSwitches[i] = time.Now().UTC().Add(-time.Minute)
		}
		return nil
	}); err != nil {
		t.Fatalf("seed fresh stamps: %v", err)
	}
	if err := switchMode(); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("switch at the fresh-stamp cap: err = %v, want ErrInvalidParameter", err)
	}

	// Aged-out stamps free the window.
	if _, err := store.Tables().Update("RollTbl", func(table *dbstore.Table) error {
		for i := range table.BillingModeSwitches {
			table.BillingModeSwitches[i] = time.Now().UTC().Add(-25 * time.Hour)
		}
		return nil
	}); err != nil {
		t.Fatalf("age the stamps: %v", err)
	}
	if err := switchMode(); err != nil {
		t.Fatalf("switch after the stamps aged out: %v", err)
	}
	table, err := store.Tables().Get("RollTbl")
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if table.BillingMode != dbstore.BillingModePayPerRequest {
		t.Fatalf("billing mode = %s, want PAY_PER_REQUEST", table.BillingMode)
	}
	if len(table.BillingModeSwitches) != 1 {
		t.Fatalf("stamps = %d, want 1 (aged stamps dropped, the new switch recorded alone)", len(table.BillingModeSwitches))
	}
}
