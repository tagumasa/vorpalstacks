// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"testing"

	"vorpalstacks/internal/common/request"
)

// describeTableForCapacity describes a table the way a client does.
func describeTableForCapacity(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext, name string) map[string]interface{} {
	t.Helper()
	resp, err := svc.DescribeTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": name}})
	if err != nil {
		t.Fatalf("describe %s: %v", name, err)
	}
	return resp.(map[string]interface{})["Table"].(map[string]interface{})
}

// TestUpdateTableAppliesCapacityMembers pins the capacity members the model
// defines on UpdateTable — OnDemandThroughput and WarmThroughput — which the
// creation plane already parsed: an update carrying them replaces the stored
// pairs, and DescribeTable reports both, the on-demand pair now rendering on
// the table description the way the model's output member carries it.
func TestUpdateTableAppliesCapacityMembers(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	name := "CapacityUpdateTable"

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            name,
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":          name,
		"OnDemandThroughput": map[string]interface{}{"MaxReadRequestUnits": 100.0, "MaxWriteRequestUnits": 200.0},
		"WarmThroughput":     map[string]interface{}{"ReadUnitsPerSecond": 300.0, "WriteUnitsPerSecond": 400.0},
	}}); err != nil {
		t.Fatalf("update capacity: %v", err)
	}

	desc := describeTableForCapacity(t, svc, reqCtx, name)
	odt, ok := desc["OnDemandThroughput"].(map[string]interface{})
	if !ok || odt["MaxReadRequestUnits"].(int64) != 100 || odt["MaxWriteRequestUnits"].(int64) != 200 {
		t.Fatalf("on-demand pair after update: %v", desc["OnDemandThroughput"])
	}
	wt, ok := desc["WarmThroughput"].(map[string]interface{})
	if !ok || wt["ReadUnitsPerSecond"].(int64) != 300 || wt["WriteUnitsPerSecond"].(int64) != 400 {
		t.Fatalf("warm pair after update: %v", desc["WarmThroughput"])
	}

	// The creation plane renders its stored on-demand pair through the same
	// member — the positive control for the render path.
	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "CapacityCreateTable",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
		"OnDemandThroughput":   map[string]interface{}{"MaxReadRequestUnits": 11.0, "MaxWriteRequestUnits": 12.0},
	}}); err != nil {
		t.Fatalf("create with on-demand pair: %v", err)
	}
	odt, ok = describeTableForCapacity(t, svc, reqCtx, "CapacityCreateTable")["OnDemandThroughput"].(map[string]interface{})
	if !ok || odt["MaxReadRequestUnits"].(int64) != 11 || odt["MaxWriteRequestUnits"].(int64) != 12 {
		t.Fatalf("on-demand pair after create: %v", odt)
	}

	// The model's member bounds: a present member is a limit of at least
	// 1 or the removal sentinel -1; zero and other negatives are
	// malformed, a present-empty pair names no limit at all, and an
	// absent member stays absent in the response.
	for _, tc := range []struct {
		label string
		pair  map[string]interface{}
	}{
		{"zero limit", map[string]interface{}{"MaxReadRequestUnits": 0.0}},
		{"negative limit", map[string]interface{}{"MaxReadRequestUnits": -5.0}},
		{"fractional limit", map[string]interface{}{"MaxWriteRequestUnits": 2.5}},
		{"empty pair", map[string]interface{}{}},
	} {
		if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":          name,
			"OnDemandThroughput": tc.pair,
		}}); !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("%s: expected ErrInvalidParameter, got %v", tc.label, err)
		}
	}

	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":          name,
		"OnDemandThroughput": map[string]interface{}{"MaxReadRequestUnits": 7.0},
	}}); err != nil {
		t.Fatalf("single-member pair: %v", err)
	}
	odt, ok = describeTableForCapacity(t, svc, reqCtx, name)["OnDemandThroughput"].(map[string]interface{})
	if !ok || odt["MaxReadRequestUnits"].(int64) != 7 {
		t.Fatalf("single-member pair: %v", odt)
	}
	if _, rendered := odt["MaxWriteRequestUnits"]; rendered {
		t.Fatalf("absent member rendered: %v", odt)
	}

	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":          name,
		"OnDemandThroughput": map[string]interface{}{"MaxReadRequestUnits": -1.0},
	}}); err != nil {
		t.Fatalf("removal sentinel: %v", err)
	}
	odt, ok = describeTableForCapacity(t, svc, reqCtx, name)["OnDemandThroughput"].(map[string]interface{})
	if !ok || odt["MaxReadRequestUnits"].(int64) != -1 {
		t.Fatalf("removal sentinel: %v", odt)
	}
}

// TestGSIActionCapacityMembersAreApplied pins the OnDemandThroughput and
// WarmThroughput members the model defines on the GSI create and update
// actions: a created index carries the pairs the request gives it, and an
// update action replaces them.
func TestGSIActionCapacityMembersAreApplied(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	name := "GsiCapacityTable"

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            name,
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}, map[string]interface{}{"AttributeName": "gpk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	gsiDesc := func() map[string]interface{} {
		t.Helper()
		gsis := describeTableForCapacity(t, svc, reqCtx, name)["GlobalSecondaryIndexes"].([]map[string]interface{})
		if len(gsis) != 1 {
			t.Fatalf("expected one GSI, got %d", len(gsis))
		}
		return gsis[0]
	}

	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": name,
		"GlobalSecondaryIndexUpdates": []interface{}{map[string]interface{}{
			"Create": map[string]interface{}{
				"IndexName":          "gsi",
				"KeySchema":          []interface{}{map[string]interface{}{"AttributeName": "gpk", "KeyType": "HASH"}},
				"Projection":         map[string]interface{}{"ProjectionType": "ALL"},
				"OnDemandThroughput": map[string]interface{}{"MaxReadRequestUnits": 30.0, "MaxWriteRequestUnits": 40.0},
				"WarmThroughput":     map[string]interface{}{"ReadUnitsPerSecond": 50.0, "WriteUnitsPerSecond": 60.0},
			},
		}},
	}}); err != nil {
		t.Fatalf("create GSI with capacity members: %v", err)
	}
	gsi := gsiDesc()
	odt, ok := gsi["OnDemandThroughput"].(map[string]interface{})
	if !ok || odt["MaxReadRequestUnits"].(int64) != 30 || odt["MaxWriteRequestUnits"].(int64) != 40 {
		t.Fatalf("GSI on-demand pair after create action: %v", gsi["OnDemandThroughput"])
	}
	wt, ok := gsi["WarmThroughput"].(map[string]interface{})
	if !ok || wt["ReadUnitsPerSecond"].(int64) != 50 || wt["WriteUnitsPerSecond"].(int64) != 60 {
		t.Fatalf("GSI warm pair after create action: %v", gsi["WarmThroughput"])
	}

	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": name,
		"GlobalSecondaryIndexUpdates": []interface{}{map[string]interface{}{
			"Update": map[string]interface{}{
				"IndexName":          "gsi",
				"OnDemandThroughput": map[string]interface{}{"MaxReadRequestUnits": 70.0, "MaxWriteRequestUnits": 80.0},
				"WarmThroughput":     map[string]interface{}{"ReadUnitsPerSecond": 90.0, "WriteUnitsPerSecond": 100.0},
			},
		}},
	}}); err != nil {
		t.Fatalf("update GSI capacity members: %v", err)
	}
	gsi = gsiDesc()
	odt, ok = gsi["OnDemandThroughput"].(map[string]interface{})
	if !ok || odt["MaxReadRequestUnits"].(int64) != 70 || odt["MaxWriteRequestUnits"].(int64) != 80 {
		t.Fatalf("GSI on-demand pair after update action: %v", gsi["OnDemandThroughput"])
	}
	wt, ok = gsi["WarmThroughput"].(map[string]interface{})
	if !ok || wt["ReadUnitsPerSecond"].(int64) != 90 || wt["WriteUnitsPerSecond"].(int64) != 100 {
		t.Fatalf("GSI warm pair after update action: %v", gsi["WarmThroughput"])
	}
}

// TestCreateTableAppliesInlineResourcePolicy pins the ResourcePolicy member
// the model defines on CreateTableInput: the inline policy attaches through
// the same core the PutResourcePolicy operation uses, and an explicitly
// empty policy document is refused before the table exists.
func TestCreateTableAppliesInlineResourcePolicy(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"dynamodb:GetItem","Resource":"*"}]}`
	resp, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "PolicyTable",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
		"ResourcePolicy":       policy,
	}})
	if err != nil {
		t.Fatalf("create with inline policy: %v", err)
	}
	tableArn := resp.(map[string]interface{})["TableDescription"].(map[string]interface{})["TableArn"].(string)

	got, err := svc.GetResourcePolicy(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"ResourceArn": tableArn}})
	if err != nil {
		t.Fatalf("get resource policy: %v", err)
	}
	if got.(map[string]interface{})["Policy"].(string) != policy {
		t.Fatalf("stored policy: %v", got)
	}
	if rev, ok := got.(map[string]interface{})["RevisionId"].(string); !ok || rev == "" {
		t.Fatalf("policy revision: %v", got)
	}

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "PolicyEmptyTable",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
		"ResourcePolicy":       "",
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty inline policy: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.DescribeTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "PolicyEmptyTable"}}); err == nil {
		t.Fatal("refused create left a table behind")
	}
}
