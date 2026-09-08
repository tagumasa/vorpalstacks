package lambda

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// canonicalRuntimeService builds a service over a fresh empty region store
// so the create/update cores can run their full validation + persistence
// path.
func canonicalRuntimeService(t *testing.T) (*LambdaService, *lambdaStore) {
	t.Helper()
	svc := NewLambdaService(nil, "000000000000", "us-east-1", t.TempDir())
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	svc.SetStorageManager(mgr)
	return svc, svc.getOrCreateLambdaStore("us-east-1")
}

// TestCreateFunctionCoreNormalisesRuntime pins the canonical-form rule on
// the create path: a case-variant Runtime request validates and persists
// as the canonical lowercase value, so the stored configuration can never
// diverge from the runtime image lookup (a mixed-case string used to
// persist as-is and silently execute on provided:al2).
func TestCreateFunctionCoreNormalisesRuntime(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "canonical-create",
		Runtime:      "Python3.12",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
	})
	if err != nil {
		t.Fatalf("create with case-variant runtime: %v", err)
	}
	if string(fn.Runtime) != "python3.12" {
		t.Fatalf("returned runtime %q, want canonical %q", fn.Runtime, "python3.12")
	}
	stored, err := stores.Functions.Get("canonical-create")
	if err != nil {
		t.Fatalf("read back created function: %v", err)
	}
	if string(stored.Runtime) != "python3.12" {
		t.Fatalf("stored runtime %q, want canonical %q", stored.Runtime, "python3.12")
	}
}

// TestUpdateFunctionConfigurationCoreNormalisesRuntime pins the same rule
// on the update path.
func TestUpdateFunctionConfigurationCoreNormalisesRuntime(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	_, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "canonical-update",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	fn, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{FunctionName: "canonical-update", Runtime: "NODEJS22.X"})
	if err != nil {
		t.Fatalf("update with case-variant runtime: %v", err)
	}
	if string(fn.Runtime) != "nodejs22.x" {
		t.Fatalf("returned runtime %q, want canonical %q", fn.Runtime, "nodejs22.x")
	}
	stored, err := stores.Functions.Get("canonical-update")
	if err != nil {
		t.Fatalf("read back updated function: %v", err)
	}
	if string(stored.Runtime) != "nodejs22.x" {
		t.Fatalf("stored runtime %q, want canonical %q", stored.Runtime, "nodejs22.x")
	}
}
