package sfn

import (
	"context"
	"errors"
	"fmt"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// newAliasTestStore provisions a state machine with two published versions
// for the alias idempotency tests.
func newAliasTestStore(t *testing.T) (*sfnstore.StepFunctionStore, string, string, string) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")

	ctx := context.Background()
	sm := &sfnstore.StateMachine{
		Name:       "alias-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}
	v1, err := store.PublishStateMachineVersion(ctx, sm.StateMachineArn, "")
	if err != nil {
		t.Fatal(err)
	}

	// A second version for routing-configuration conflicts: publishing
	// follows a revision change, so bump the revision before republishing.
	sm2, err := store.GetStateMachine(ctx, sm.StateMachineArn)
	if err != nil {
		t.Fatal(err)
	}
	sm2.RevisionId = "revision-two"
	if err := store.UpdateStateMachine(ctx, sm2); err != nil {
		t.Fatal(err)
	}
	v2, err := store.PublishStateMachineVersion(ctx, sm.StateMachineArn, "")
	if err != nil {
		t.Fatal(err)
	}
	return store, sm.StateMachineArn, v1.StateMachineVersionArn, v2.StateMachineVersionArn
}

func requireAWSCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %s error, got nil", code)
	}
	var ae *awserrors.AWSError
	if !errors.As(err, &ae) {
		t.Fatalf("expected an AWSError, got %T: %v", err, err)
	}
	if ae.Code != code {
		t.Fatalf("expected %s, got %s (%s)", code, ae.Code, ae.Message)
	}
}

func TestCreateStateMachineAliasIdempotentRetry(t *testing.T) {
	store, _, v1, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	first, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	second, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	})
	if err != nil {
		t.Fatalf("identical retry must return an idempotent success: %v", err)
	}
	if first["stateMachineAliasArn"] != second["stateMachineAliasArn"] {
		t.Fatalf("identical retry returned a different alias ARN: %v vs %v", first["stateMachineAliasArn"], second["stateMachineAliasArn"])
	}
	if first["creationDate"] != second["creationDate"] {
		t.Fatalf("identical retry returned a different creationDate: %v vs %v", first["creationDate"], second["creationDate"])
	}
}

func TestCreateStateMachineAliasSameNameDifferentParametersConflict(t *testing.T) {
	store, _, v1, v2 := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	if _, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	}); err != nil {
		t.Fatalf("first create: %v", err)
	}

	// Same name with a different description conflicts.
	_, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		Description:   "different",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	})
	requireAWSCode(t, err, "ConflictException")

	// Same name routing to a different version conflicts.
	_, err = svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v2, Weight: 100}},
	})
	requireAWSCode(t, err, "ConflictException")

	// Same name with a split routing configuration conflicts.
	_, err = svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name: "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{
			{StateMachineVersionArn: v1, Weight: 50},
			{StateMachineVersionArn: v2, Weight: 50},
		},
	})
	requireAWSCode(t, err, "ConflictException")
}

func TestPublishStateMachineVersionIdempotentPerRevision(t *testing.T) {
	store, smArn, _, v2 := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	sm, err := store.GetStateMachine(ctx, smArn)
	if err != nil {
		t.Fatal(err)
	}

	// The current revision already has v2 (published by the helper);
	// republishing must return it rather than a new version.
	resp, err := svc.publishStateMachineVersionCore(ctx, store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
		RevisionId:      sm.RevisionId,
	})
	if err != nil {
		t.Fatalf("idempotent republish: %v", err)
	}
	if resp["stateMachineVersionArn"] != v2 {
		t.Fatalf("republish returned %v, want the existing %s", resp["stateMachineVersionArn"], v2)
	}
	count, err := store.CountStateMachineVersions(smArn)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected the two seeded versions, got %d", count)
	}
}

func TestPublishStateMachineVersionRevisionMismatchConflicts(t *testing.T) {
	store, smArn, _, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}

	_, err := svc.publishStateMachineVersionCore(context.Background(), store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
		RevisionId:      "stale-revision",
	})
	requireAWSCode(t, err, "ConflictException")
}

func TestPublishStateMachineVersionQuota(t *testing.T) {
	store, smArn, _, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	// Fill the quota with distinct revisions; the helper already published
	// two versions, so only the remainder is seeded here.
	sm, err := store.GetStateMachine(ctx, smArn)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < sfnstore.MaxVersionsPerStateMachine-2; i++ {
		sm.RevisionId = fmt.Sprintf("quota-revision-%d", i)
		if err := store.UpdateStateMachine(ctx, sm); err != nil {
			t.Fatal(err)
		}
		if _, err := store.PublishStateMachineVersion(ctx, smArn, ""); err != nil {
			t.Fatal(err)
		}
	}

	count, err := store.CountStateMachineVersions(smArn)
	if err != nil {
		t.Fatal(err)
	}
	if count != sfnstore.MaxVersionsPerStateMachine {
		t.Fatalf("expected %d seeded versions, got %d", sfnstore.MaxVersionsPerStateMachine, count)
	}

	// A fresh revision past the quota is rejected.
	sm.RevisionId = "quota-exceeding-revision"
	if err := store.UpdateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}
	_, err = svc.publishStateMachineVersionCore(ctx, store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
	})
	requireAWSCode(t, err, "ServiceQuotaExceededException")

	// An idempotent retry of the already-published current revision still
	// succeeds at the quota.
	sm.RevisionId = "quota-revision-0"
	if err := store.UpdateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}
	resp, err := svc.publishStateMachineVersionCore(ctx, store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
	})
	if err != nil {
		t.Fatalf("idempotent retry at the quota: %v", err)
	}
	if resp["stateMachineVersionArn"] == "" {
		t.Fatal("idempotent retry returned no version ARN")
	}
}

func TestUpdateStateMachineAliasRequiresDescriptionOrRouting(t *testing.T) {
	store, _, v1, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	if _, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	}); err != nil {
		t.Fatalf("create: %v", err)
	}
	smArn := smArnOfVersion(t, store, v1)
	alias, err := store.GetStateMachineAliasByName(ctx, smArn, "prod")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: alias.StateMachineAliasArn,
	})
	requireAWSCode(t, err, "ValidationException")

	if _, err := svc.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: alias.StateMachineAliasArn,
		DescriptionProvided:  true,
		Description:          "updated",
	}); err != nil {
		t.Fatalf("description-only update: %v", err)
	}
}

func smArnOfVersion(t *testing.T, store *sfnstore.StepFunctionStore, versionArn string) string {
	t.Helper()
	v, err := store.GetStateMachineVersion(context.Background(), versionArn)
	if err != nil {
		t.Fatal(err)
	}
	return v.StateMachineArn
}

func TestRoutingConfigsEqualMatchesEntriesRegardlessOfOrder(t *testing.T) {
	a := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v1", Weight: 90},
		{StateMachineVersionArn: "v2", Weight: 10},
	}
	sameReordered := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v2", Weight: 10},
		{StateMachineVersionArn: "v1", Weight: 90},
	}
	differentWeights := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v1", Weight: 50},
		{StateMachineVersionArn: "v2", Weight: 50},
	}
	if !routingConfigsEqual(a, sameReordered) {
		t.Fatal("reordered entries must compare equal")
	}
	if routingConfigsEqual(a, differentWeights) {
		t.Fatal("different weights must not compare equal")
	}
	if routingConfigsEqual(a, a[:1]) {
		t.Fatal("different lengths must not compare equal")
	}
}
