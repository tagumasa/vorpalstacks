package lambda

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
)

// TestLambdaStoreCacheSingleConstructionPath pins the single-construction
// rule for the per-region store cache: the request-plane accessor, the
// admin-plane accessor and the function-store accessor all resolve to the
// one cached lambdaStore instance per region, and a second region gets its
// own instance.
func TestLambdaStoreCacheSingleConstructionPath(t *testing.T) {
	svc := NewLambdaService(nil, "000000000000", "us-east-1", t.TempDir())
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	svc.SetStorageManager(mgr)
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")

	viaRequest, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("request-plane accessor: %v", err)
	}
	viaAdmin := svc.getOrCreateLambdaStore("us-east-1")
	if viaAdmin == nil {
		t.Fatal("admin-plane accessor returned nil")
	}
	viaFunctionStore := svc.getOrCreateFunctionStore("us-east-1")

	if viaRequest != viaAdmin {
		t.Fatal("request-plane and admin-plane accessors returned different instances")
	}
	if viaRequest.Functions != viaFunctionStore {
		t.Fatal("function-store accessor resolved a different FunctionStore")
	}

	cached, ok := svc.storeCache.Load("us-east-1")
	if !ok || cached != viaRequest {
		t.Fatal("cache does not hold the instance the accessors returned")
	}

	viaE, err := svc.getOrCreateLambdaStoreE("us-east-1")
	if err != nil {
		t.Fatalf("error-returning accessor: %v", err)
	}
	if viaE != viaRequest {
		t.Fatal("error-returning accessor built a second instance")
	}

	if other := svc.getOrCreateLambdaStore("eu-west-1"); other == viaRequest {
		t.Fatal("second region reused the first region's instance")
	}
}
