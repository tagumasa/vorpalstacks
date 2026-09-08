package lambda

import (
	"context"
	"errors"
	"os"
	"testing"

	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// workingStorageReqCtxForService builds a request context over the
// service's own storage manager, so cores resolving stores from the
// context see the same records as svc.getOrCreateLambdaStore.
func workingStorageReqCtxForService(t *testing.T, svc *LambdaService) *request.RequestContext {
	t.Helper()
	return request.NewRequestContext(context.Background(), svc.storageManager, "000000000000", "us-east-1")
}

// codeDirExists reports whether an on-disk path still exists.
func codeDirExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", path, err)
	}
	return false
}

// TestDeleteFunctionRemovesCodeRoot pins the disk half of function
// deletion: deleting the function removes its whole on-disk code root —
// $LATEST plus every published version archive — instead of leaving the
// archives behind forever.
func TestDeleteFunctionRemovesCodeRoot(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	if _, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "cleanup-full",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
	}); err != nil {
		t.Fatalf("create function: %v", err)
	}
	if _, _, err := svc.storeCode("cleanup-full", "$LATEST", []byte("zip"), "us-east-1"); err != nil {
		t.Fatalf("seed $LATEST: %v", err)
	}
	if _, _, err := svc.storeCode("cleanup-full", "1", []byte("zip"), "us-east-1"); err != nil {
		t.Fatalf("seed version 1: %v", err)
	}

	if err := svc.deleteFunctionCore(context.Background(), stores, &DeleteFunctionInput{
		FunctionName: "cleanup-full",
		Region:       "us-east-1",
	}); err != nil {
		t.Fatalf("delete function: %v", err)
	}

	for _, path := range []string{
		functionCodeDir(svc.dataDir, "us-east-1", "cleanup-full", "$LATEST"),
		functionCodeDir(svc.dataDir, "us-east-1", "cleanup-full", "1"),
	} {
		if codeDirExists(t, path) {
			t.Fatalf("code residue survived function deletion: %s", path)
		}
	}
}

// TestDeleteFunctionVersionRemovesOnlyThatDirectory pins the qualified
// delete: removing version 1 removes exactly that version's archive,
// leaving $LATEST and the other versions executable.
func TestDeleteFunctionVersionRemovesOnlyThatDirectory(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "cleanup-ver",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
	})
	if err != nil {
		t.Fatalf("create function: %v", err)
	}
	if _, _, err := svc.storeCode("cleanup-ver", "$LATEST", []byte("zip"), "us-east-1"); err != nil {
		t.Fatalf("seed $LATEST: %v", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := svc.publishVersionWithCode(stores, fn, "", "", "us-east-1"); err != nil {
			t.Fatalf("publish version %d: %v", i+1, err)
		}
	}

	if err := svc.deleteFunctionCore(context.Background(), stores, &DeleteFunctionInput{
		FunctionName: fn.FunctionName,
		Qualifier:    "1",
		Region:       "us-east-1",
	}); err != nil {
		t.Fatalf("delete version: %v", err)
	}

	if codeDirExists(t, functionCodeDir(svc.dataDir, "us-east-1", "cleanup-ver", "1")) {
		t.Fatalf("deleted version's archive survived")
	}
	for _, v := range []string{"$LATEST", "2"} {
		if !codeDirExists(t, functionCodeDir(svc.dataDir, "us-east-1", "cleanup-ver", v)) {
			t.Fatalf("surviving version %s was removed with the qualified delete", v)
		}
	}
}

// TestDeleteLayerVersionRemovesArchiveAndRoot pins the layer half: a
// version delete removes that version's archive, and the layer root when
// the last version went with it.
func TestDeleteLayerVersionRemovesArchiveAndRoot(t *testing.T) {
	svc, _ := canonicalRuntimeService(t)
	reqCtx := workingStorageReqCtxForService(t, svc)

	publish := func(layer string) {
		t.Helper()
		if _, _, err := svc.publishLayerVersionCore(context.Background(), reqCtx, &LayerVersionCreateInput{
			LayerName: layer,
			Content:   map[string]interface{}{"ZipFile": "emhhcg=="},
			Region:    "us-east-1",
		}); err != nil {
			t.Fatalf("publish layer version: %v", err)
		}
	}
	publish("cleanup-layer")
	publish("cleanup-layer")

	if err := svc.deleteLayerVersionCore(reqCtx, "cleanup-layer", 1); err != nil {
		t.Fatalf("delete layer version: %v", err)
	}
	if codeDirExists(t, layerVersionCodeDir(svc.dataDir, "us-east-1", "cleanup-layer", 1)) {
		t.Fatalf("deleted layer version's archive survived")
	}
	if !codeDirExists(t, layerVersionCodeDir(svc.dataDir, "us-east-1", "cleanup-layer", 2)) {
		t.Fatalf("surviving layer version 2 was removed with the delete")
	}

	if err := svc.deleteLayerVersionCore(reqCtx, "cleanup-layer", 2); err != nil {
		t.Fatalf("delete last layer version: %v", err)
	}
	if codeDirExists(t, layerCodeRoot(svc.dataDir, "us-east-1", "cleanup-layer")) {
		t.Fatalf("empty layer root survived the last version delete")
	}
}

// TestFailedCreateRemovesSeededCodeArchives pins the rollback half: a
// CreateFunction whose publish step fails removes the seeded $LATEST
// archive together with the rolled-back record.
func TestFailedCreateRemovesSeededCodeArchives(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	if _, _, err := svc.storeCode("cleanup-rollback", "$LATEST", []byte("zip"), "us-east-1"); err != nil {
		t.Fatalf("seed $LATEST: %v", err)
	}
	blockCodePath(t, svc.dataDir, "us-east-1", "code", "cleanup-rollback", "1")

	if _, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "cleanup-rollback",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		Publish:      true,
		Region:       "us-east-1",
	}); err == nil {
		t.Fatalf("expected publish failure")
	}

	if _, err := stores.Functions.Get("cleanup-rollback"); !errors.Is(err, lambdastore.ErrFunctionNotFound) {
		t.Fatalf("record survived rollback: %v", err)
	}
	if codeDirExists(t, functionCodeDir(svc.dataDir, "us-east-1", "cleanup-rollback", "$LATEST")) {
		t.Fatalf("seeded $LATEST archive survived the rolled-back create")
	}
}
