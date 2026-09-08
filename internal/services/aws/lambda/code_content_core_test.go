package lambda

import (
	"context"
	"errors"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// TestPublishLayerVersionRejectsEmptyContentMember pins the codeless-version
// hole: a non-nil Content map naming neither ZipFile nor an S3 bucket is
// InvalidParameterValueException, and the failed publish leaves no layer
// shell behind.
func TestPublishLayerVersionRejectsEmptyContentMember(t *testing.T) {
	svc, _ := canonicalRuntimeService(t)
	reqCtx := workingStorageReqCtxForService(t, svc)

	_, _, err := svc.publishLayerVersionCore(context.Background(), reqCtx, &LayerVersionCreateInput{
		LayerName: "empty-content",
		Content:   map[string]interface{}{"Description": "no usable member"},
		Region:    "us-east-1",
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("empty content member: expected InvalidParameterValueException, got %v", err)
	}

	if _, gerr := svc.getOrCreateLambdaStore("us-east-1").Layers.Get("empty-content"); gerr == nil {
		t.Fatal("failed publish must not leave a layer shell behind")
	}
}

// TestPublishLayerVersionRejectsInvalidBase64WithoutShell pins the error
// shape and the side-effect ordering: an undecodable ZipFile is
// InvalidParameterValueException (a client error, not an internal failure),
// and because decoding precedes record creation no layer shell survives the
// rejected request.
func TestPublishLayerVersionRejectsInvalidBase64WithoutShell(t *testing.T) {
	svc, _ := canonicalRuntimeService(t)
	reqCtx := workingStorageReqCtxForService(t, svc)

	_, _, err := svc.publishLayerVersionCore(context.Background(), reqCtx, &LayerVersionCreateInput{
		LayerName: "bad-base64",
		Content:   map[string]interface{}{"ZipFile": "!!!not-base64!!!"},
		Region:    "us-east-1",
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("invalid base64: expected InvalidParameterValueException, got %v", err)
	}

	if _, gerr := svc.getOrCreateLambdaStore("us-east-1").Layers.Get("bad-base64"); gerr == nil {
		t.Fatal("failed publish must not leave a layer shell behind")
	}
}

// TestPrepareCreateFunctionCodeCoreRejectsEmptyCodeMap pins the create-side
// hole: a Code map naming neither ZipFile, S3Bucket nor ImageUri is rejected
// instead of persisting a zero-size archive whose hash the store would have
// to fabricate.
func TestPrepareCreateFunctionCodeCoreRejectsEmptyCodeMap(t *testing.T) {
	svc, _ := canonicalRuntimeService(t)

	_, _, _, err := svc.prepareCreateFunctionCodeCore(context.Background(), "us-east-1", "empty-code",
		map[string]interface{}{}, "Zip")
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("empty code map: expected InvalidParameterValueException, got %v", err)
	}
}

// TestPrepareCreateFunctionCodeCoreDecodesZipFile pins the surviving happy
// path: the decoded archive is persisted under $LATEST with its hash and
// size recorded in the metadata.
func TestPrepareCreateFunctionCodeCoreDecodesZipFile(t *testing.T) {
	svc, _ := canonicalRuntimeService(t)

	meta, imageUri, packageType, err := svc.prepareCreateFunctionCodeCore(context.Background(), "us-east-1", "zip-code",
		map[string]interface{}{"ZipFile": "emhhcg=="}, "Zip")
	if err != nil {
		t.Fatalf("zip file decode: %v", err)
	}
	if imageUri != "" || packageType != "Zip" {
		t.Fatalf("imageUri/packageType %q/%q, want empty/Zip", imageUri, packageType)
	}
	if meta.CodeSize != 4 || meta.CodeSha256 != lambdastore.GenerateCodeHash([]byte("zhar")) {
		t.Fatalf("metadata size/hash %d/%q, want 4 and the archive hash", meta.CodeSize, meta.CodeSha256)
	}
	if meta.CodeLocation == "" {
		t.Fatal("metadata carries no persisted code location")
	}
}

// TestPrepareFunctionCodeUpdateCoreRejectsMissingMember pins the update-side
// required-member contract on the canonical code map: at least one of
// ZipFile, ImageUri or S3Bucket must be present.
func TestPrepareFunctionCodeUpdateCoreRejectsMissingMember(t *testing.T) {
	svc, _ := canonicalRuntimeService(t)

	_, err := svc.prepareFunctionCodeUpdateCore(context.Background(), "us-east-1", "upd-code",
		map[string]interface{}{"S3Key": "k"})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("missing code member: expected InvalidParameterValueException, got %v", err)
	}

	meta, err := svc.prepareFunctionCodeUpdateCore(context.Background(), "us-east-1", "upd-code",
		map[string]interface{}{"ImageUri": "0123.dkr.ecr.image"})
	if err != nil {
		t.Fatalf("image-only update: %v", err)
	}
	if meta.CodeLocation != "" || meta.CodeSize != 0 || meta.CodeSha256 != "" {
		t.Fatalf("image-only update carried zip metadata: %+v", meta)
	}
}
