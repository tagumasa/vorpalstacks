package sfn

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func newResolverTestStore(t *testing.T) *sfnstore.StepFunctionStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")
	sm := &sfnstore.StateMachine{
		Name:       "resolver-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	}
	if err := store.CreateStateMachine(context.Background(), sm); err != nil {
		t.Fatal(err)
	}
	return store
}

func TestResolveStateMachineReferenceBase(t *testing.T) {
	store := newResolverTestStore(t)
	arn := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm"

	ref, err := resolveStateMachineReference(context.Background(), store, arn)
	if err != nil {
		t.Fatalf("base resolution failed: %v", err)
	}
	if ref.Version != nil || ref.Alias != nil {
		t.Fatalf("unqualified ARN must not resolve to a version or alias")
	}
	if ref.definition() != ref.StateMachine.Definition {
		t.Fatalf("base reference must use the live definition")
	}
}

func TestResolveStateMachineReferenceVersion(t *testing.T) {
	store := newResolverTestStore(t)
	smArn := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm"

	version, err := store.PublishStateMachineVersion(context.Background(), smArn, "")
	if err != nil {
		t.Fatal(err)
	}

	ref, err := resolveStateMachineReference(context.Background(), store, version.StateMachineVersionArn)
	if err != nil {
		t.Fatalf("version resolution failed: %v", err)
	}
	if ref.Version == nil || ref.Version.StateMachineVersionArn != version.StateMachineVersionArn {
		t.Fatalf("version qualifier must resolve to the version record")
	}
	if ref.Alias != nil {
		t.Fatalf("version ARN must not resolve to an alias")
	}
}

func TestResolveStateMachineReferenceAlias(t *testing.T) {
	store := newResolverTestStore(t)
	smArn := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm"

	version, err := store.PublishStateMachineVersion(context.Background(), smArn, "")
	if err != nil {
		t.Fatal(err)
	}
	alias := &sfnstore.StateMachineAlias{
		StateMachineArn: smArn,
		Name:            "PROD",
		RoutingConfiguration: []sfnstore.RoutingConfiguration{
			{StateMachineVersionArn: version.StateMachineVersionArn, Weight: 100},
		},
	}
	if err := store.CreateStateMachineAlias(context.Background(), alias); err != nil {
		t.Fatal(err)
	}

	want := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm:PROD"
	if alias.StateMachineAliasArn != want {
		t.Fatalf("alias ARN = %q, want %q", alias.StateMachineAliasArn, want)
	}

	ref, err := resolveStateMachineReference(context.Background(), store, alias.StateMachineAliasArn)
	if err != nil {
		t.Fatalf("alias resolution failed: %v", err)
	}
	if ref.Alias == nil || ref.Alias.Name != "PROD" {
		t.Fatalf("alias qualifier must resolve to the alias record")
	}
	if ref.Version == nil || ref.Version.StateMachineVersionArn != version.StateMachineVersionArn {
		t.Fatalf("alias must route to the routed version")
	}
	if ref.definition() != version.Definition {
		t.Fatalf("alias reference must run the version snapshot")
	}
}

func TestResolveStateMachineReferenceUnknownQualifier(t *testing.T) {
	store := newResolverTestStore(t)

	unknownVersion := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm:99"
	if _, err := resolveStateMachineReference(context.Background(), store, unknownVersion); err == nil {
		t.Fatal("unknown version qualifier must fail")
	}

	unknownAlias := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm:NOSUCH"
	if _, err := resolveStateMachineReference(context.Background(), store, unknownAlias); err == nil {
		t.Fatal("unknown alias qualifier must fail")
	}
}

func TestResolveStateMachineReferenceMalformedAndLabel(t *testing.T) {
	store := newResolverTestStore(t)

	if _, err := resolveStateMachineReference(context.Background(), store, "arn:aws:iam::000000000000:role/x"); err == nil {
		t.Fatal("non-States ARN must be rejected")
	}
	label := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm/mapLabel"
	if _, err := resolveStateMachineReference(context.Background(), store, label); err == nil {
		t.Fatal("Distributed Map label ARN must be rejected")
	}
}

func TestSelectVersionByWeightSingleAndSplit(t *testing.T) {
	single := []sfnstore.RoutingConfiguration{{StateMachineVersionArn: "v1", Weight: 100}}
	for i := 0; i < 8; i++ {
		got, err := selectVersionByWeight(single)
		if err != nil || got != "v1" {
			t.Fatalf("single entry must be deterministic, got %q err %v", got, err)
		}
	}

	split := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v1", Weight: 50},
		{StateMachineVersionArn: "v2", Weight: 50},
	}
	seen := map[string]bool{}
	for i := 0; i < 200 && len(seen) < 2; i++ {
		got, err := selectVersionByWeight(split)
		if err != nil {
			t.Fatalf("split selection failed: %v", err)
		}
		if got != "v1" && got != "v2" {
			t.Fatalf("unexpected version %q", got)
		}
		seen[got] = true
	}
	if len(seen) != 2 {
		t.Fatal("a 50/50 split must reach both versions over repeated picks")
	}
}
