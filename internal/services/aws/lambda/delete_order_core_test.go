package lambda

import (
	"context"
	"strings"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// TestDeleteAliasReferencedVersionRejectsBeforeSideEffects pins the
// ordering of the qualified delete: DeleteVersion is the authoritative
// check that rejects an alias-referenced version, so it must run before
// the container and sandbox cleanup — a version that survives the
// request keeps its container untouched. The fixture runs against a
// service with no container runtime and a version carrying a container
// ID, so a side-effect-first ordering panics instead of rejecting.
func TestDeleteAliasReferencedVersionRejectsBeforeSideEffects(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)

	fn, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "alias-guard",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
	})
	if err != nil {
		t.Fatalf("create function: %v", err)
	}
	if _, _, err := svc.storeCode("alias-guard", "$LATEST", []byte("zip"), "us-east-1"); err != nil {
		t.Fatalf("seed $LATEST code: %v", err)
	}
	if _, err := svc.publishVersionWithCode(stores, fn, "", "", "us-east-1"); err != nil {
		t.Fatalf("publish version 1: %v", err)
	}
	if _, err := stores.Functions.CreateAliasAtomically("alias-guard", func(f *lambdastore.Function) (*lambdastore.Alias, error) {
		return &lambdastore.Alias{Name: "prod", FunctionVersion: "1"}, nil
	}); err != nil {
		t.Fatalf("create alias: %v", err)
	}
	// The version runs with a live container for this fixture.
	if _, err := stores.Functions.UpdateAtomically("alias-guard", func(f *lambdastore.Function) error {
		for i := range f.Versions {
			if f.Versions[i].Version == "1" {
				f.Versions[i].ContainerID = "ctr-alias-guard"
				f.Versions[i].ContainerImageID = "ctr-alias-guard"
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("attach container to version: %v", err)
	}

	err = svc.deleteFunctionCore(context.Background(), stores, &DeleteFunctionInput{
		FunctionName: "alias-guard",
		Qualifier:    "1",
		Region:       "us-east-1",
	})
	if err == nil || !strings.Contains(err.Error(), "ResourceConflict") {
		t.Fatalf("alias-referenced delete: expected ResourceConflict, got %v", err)
	}
	if _, verr := stores.Functions.GetVersion("alias-guard", "1"); verr != nil {
		t.Fatalf("rejected delete lost the version record: %v", verr)
	}
	stored, err := stores.Functions.Get("alias-guard")
	if err != nil {
		t.Fatalf("reload function: %v", err)
	}
	for _, v := range stored.Versions {
		if v.Version == "1" && v.ContainerID != "ctr-alias-guard" {
			t.Fatalf("rejected delete disturbed the version's container assignment: %q", v.ContainerID)
		}
	}
}
