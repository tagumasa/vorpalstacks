package lambda

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// blockCodePath creates a regular file where a code or layer archive
// directory would be created, so the archive write fails
// deterministically on a path that cannot be mkdir'd.
func blockCodePath(t *testing.T, dataDir, region, kind, name, version string) {
	t.Helper()
	blockPath := filepath.Join(dataDir, region, kind, name, version)
	if err := os.MkdirAll(filepath.Dir(blockPath), 0755); err != nil {
		t.Fatalf("prepare block parent: %v", err)
	}
	if err := os.WriteFile(blockPath, []byte("blocked"), 0644); err != nil {
		t.Fatalf("block path: %v", err)
	}
}

// TestCreateFunctionCoreRejectsInvalidTagsBeforePersisting pins the
// create-flow ordering: an input-ground tag failure must be rejected
// before the function record is written, so a failed CreateFunction
// never leaves an orphaned record behind.
func TestCreateFunctionCoreRejectsInvalidTagsBeforePersisting(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)

	_, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "tag-reject",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		Tags:         map[string]string{"bad\x00key": "value"},
	})
	if err == nil {
		t.Fatalf("expected tag rejection")
	}
	if _, getErr := stores.Functions.Get("tag-reject"); !errors.Is(getErr, lambdastore.ErrFunctionNotFound) {
		t.Fatalf("failed create left a function record behind: %v", getErr)
	}
}

// TestCreateFunctionRollsBackWhenPublishFails pins the compensation rule
// on the create path: CreateFunction with Publish whose version code
// cannot be persisted must fail without leaving the created function.
func TestCreateFunctionRollsBackWhenPublishFails(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)

	if _, _, err := svc.storeCode("publish-fault", "$LATEST", []byte("zip"), "us-east-1"); err != nil {
		t.Fatalf("seed $LATEST code: %v", err)
	}
	blockCodePath(t, svc.dataDir, "us-east-1", "code", "publish-fault", "1")

	_, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "publish-fault",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		Publish:      true,
	})
	if err == nil {
		t.Fatalf("expected publish failure")
	}
	if _, getErr := stores.Functions.Get("publish-fault"); !errors.Is(getErr, lambdastore.ErrFunctionNotFound) {
		t.Fatalf("failed create-with-publish left a function record behind: %v", getErr)
	}
}

// TestPublishVersionRollsBackWhenCodeStoreFails pins the compensation
// rule on the function publish path: a version whose code archive
// cannot be persisted must not survive the failed request — the orphaned
// record would be permanently uninvocable.
func TestPublishVersionRollsBackWhenCodeStoreFails(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)

	fn, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "ver-fault",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
	})
	if err != nil {
		t.Fatalf("create function: %v", err)
	}
	if _, _, err := svc.storeCode("ver-fault", "$LATEST", []byte("zip"), "us-east-1"); err != nil {
		t.Fatalf("seed $LATEST code: %v", err)
	}
	blockCodePath(t, svc.dataDir, "us-east-1", "code", "ver-fault", "1")

	if _, err := svc.publishVersionWithCode(stores, fn, "", "", "us-east-1"); err == nil {
		t.Fatalf("expected code-store failure")
	}
	if _, verr := stores.Functions.GetVersion("ver-fault", "1"); !errors.Is(verr, lambdastore.ErrVersionNotFound) {
		t.Fatalf("failed publish left an orphaned version: %v", verr)
	}
}

// TestPublishLayerVersionRollsBackWhenCodeStoreFails pins the same rule
// on the layer publish path: when the archive cannot be persisted, the
// version record and a layer shell created by this very request are
// removed, matching single-operation AWS semantics.
func TestPublishLayerVersionRollsBackWhenCodeStoreFails(t *testing.T) {
	dataDir := t.TempDir()
	svc := NewLambdaService(nil, "000000000000", "us-east-1", dataDir)
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	svc.SetStorageManager(mgr)
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")

	blockCodePath(t, dataDir, "us-east-1", "layers", "fault-layer", "1")

	_, _, err = svc.publishLayerVersionCore(context.Background(), reqCtx, &LayerVersionCreateInput{
		LayerName: "fault-layer",
		Content:   map[string]interface{}{"ZipFile": base64.StdEncoding.EncodeToString([]byte("zip"))},
		Region:    "us-east-1",
	})
	if err == nil {
		t.Fatalf("expected layer code-store failure")
	}

	stores := svc.getOrCreateLambdaStore("us-east-1")
	// The whole layer is gone when the shell was created here; a
	// pre-existing layer keeps its record without the failed version.
	if _, verr := stores.Layers.GetVersion("fault-layer", 1); !errors.Is(verr, lambdastore.ErrLayerVersionNotFound) && !errors.Is(verr, lambdastore.ErrLayerNotFound) {
		t.Fatalf("failed publish left an orphaned layer version: %v", verr)
	}
	if _, lerr := stores.Layers.Get("fault-layer"); !errors.Is(lerr, lambdastore.ErrLayerNotFound) {
		t.Fatalf("failed publish left an orphaned layer shell: %v", lerr)
	}
}
