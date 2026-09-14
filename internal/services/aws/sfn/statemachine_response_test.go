package sfn

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestUpdateStateMachineResponseMembers pins the UpdateStateMachineOutput
// member set exactly: updateDate and revisionId (stateMachineVersionArn
// when a version was published) — the shape carries no stateMachineArn.
func TestUpdateStateMachineResponseMembers(t *testing.T) {
	ctx := context.Background()

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewStepFunctionService(mgr, "000000000000")
	reqCtx := request.NewRequestContext(ctx, mgr, "000000000000", "us-east-1")
	handlerStore, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	created, err := svc.createStateMachineCore(ctx, handlerStore, CreateStateMachineInput{
		Name:       "upd-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	resp, err := svc.UpdateStateMachine(ctx, reqCtx, &request.ParsedRequest{
		Parameters: map[string]interface{}{
			"stateMachineArn": created.StateMachineArn,
			"definition":      `{"StartAt":"B","States":{"B":{"Type":"Pass","End":true}}}`,
		},
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	result, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("response is %T, want the serialised member map", resp)
	}
	if _, present := result["stateMachineArn"]; present {
		t.Errorf("response carries stateMachineArn, which UpdateStateMachineOutput does not model: %v", result)
	}
	if _, present := result["updateDate"]; !present {
		t.Errorf("response lacks updateDate: %v", result)
	}
	if _, present := result["revisionId"]; !present {
		t.Errorf("response lacks revisionId: %v", result)
	}
}

// TestDescribeForExecutionMapChildLabel pins the label member: it is
// "returned only if the executionArn is a child workflow execution that
// was started by a Distributed Map state", and its value is the Map
// state's label the Map Run ARN carries.
func TestDescribeForExecutionMapChildLabel(t *testing.T) {
	store := newCreateTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	created, err := svc.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:       "map-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	plain := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:map-sm:plain",
		StateMachineArn: created.StateMachineArn,
		Name:            "plain", Status: "SUCCEEDED", Input: "{}",
	}
	if err := store.CreateExecution(ctx, plain); err != nil {
		t.Fatalf("create plain execution: %v", err)
	}
	resp, err := svc.describeStateMachineForExecutionCore(ctx, store, DescribeStateMachineForExecutionInput{
		ExecutionArn: plain.ExecutionArn,
	})
	if err != nil {
		t.Fatalf("describe plain: %v", err)
	}
	if _, present := resp["label"]; present {
		t.Errorf("plain execution response carries label: %v", resp)
	}

	child := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:map-sm:parent:MyLabel-0",
		StateMachineArn: created.StateMachineArn,
		Name:            "parent:MyLabel-0", Status: "SUCCEEDED", Input: "{}",
		MapRunArn: "arn:aws:states:us-east-1:000000000000:mapRun:map-sm/MyLabel/mapRun-1-20260913000000",
	}
	if err := store.CreateExecution(ctx, child); err != nil {
		t.Fatalf("create child execution: %v", err)
	}
	resp, err = svc.describeStateMachineForExecutionCore(ctx, store, DescribeStateMachineForExecutionInput{
		ExecutionArn: child.ExecutionArn,
	})
	if err != nil {
		t.Fatalf("describe child: %v", err)
	}
	if resp["label"] != "MyLabel" {
		t.Errorf("child label = %v, want MyLabel from the Map Run ARN", resp["label"])
	}
}

// TestDescribeStateMachineQualifiedLabel pins label's documented presence
// rule — "This parameter is present only if the stateMachineArn specified
// in input is a qualified state machine ARN": the alias-qualified describe
// carries the user-defined alias name, the version-qualified describe the
// auto-generated version number, and an unqualified ARN omits the member.
func TestDescribeStateMachineQualifiedLabel(t *testing.T) {
	store, smArn, v1, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	created, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	})
	if err != nil {
		t.Fatalf("create alias: %v", err)
	}
	aliasArn, _ := created["stateMachineAliasArn"].(string)
	if aliasArn == "" {
		t.Fatal("create alias returned no stateMachineAliasArn")
	}

	plain, err := svc.describeStateMachineCore(ctx, store, DescribeStateMachineInput{StateMachineArn: smArn})
	if err != nil {
		t.Fatalf("describe unqualified: %v", err)
	}
	if _, present := plain["label"]; present {
		t.Fatalf("unqualified ARN must omit label, got %v", plain["label"])
	}

	byVersion, err := svc.describeStateMachineCore(ctx, store, DescribeStateMachineInput{StateMachineArn: v1})
	if err != nil {
		t.Fatalf("describe version-qualified: %v", err)
	}
	if label, _ := byVersion["label"].(string); label != "1" {
		t.Fatalf("version-qualified label = %v, want the version number 1", byVersion["label"])
	}

	byAlias, err := svc.describeStateMachineCore(ctx, store, DescribeStateMachineInput{StateMachineArn: aliasArn})
	if err != nil {
		t.Fatalf("describe alias-qualified: %v", err)
	}
	if label, _ := byAlias["label"].(string); label != "prod" {
		t.Fatalf("alias-qualified label = %v, want the alias name prod", byAlias["label"])
	}
}
