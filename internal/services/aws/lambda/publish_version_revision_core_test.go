package lambda

import (
	"context"
	"errors"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// publishVersionFixture creates a function with stored $LATEST code so
// publishVersionWithCode can persist the version snapshot, then returns the
// created function carrying its current RevisionId.
func publishVersionFixture(t *testing.T, svc *LambdaService, stores *lambdaStore, name string) *lambdastore.Function {
	t.Helper()
	code := []byte("exports.handler = async () => {};")
	if _, _, err := svc.storeCode(name, "$LATEST", code, "us-east-1"); err != nil {
		t.Fatalf("store code: %v", err)
	}
	fn, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: name,
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		CodeLocation: name + "-latest.zip",
		CodeSize:     int64(len(code)),
	})
	if err != nil {
		t.Fatalf("create fixture function: %v", err)
	}
	return fn
}

// TestPublishVersionRejectsStaleRevisionId pins the modelled RevisionId
// precondition on PublishVersion: publishing with a revision that does not
// match the function's current revision fails with
// PreconditionFailedException (HTTP 412), exactly like the update paths.
func TestPublishVersionRejectsStaleRevisionId(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn := publishVersionFixture(t, svc, stores, "publish-stale")

	_, err := svc.publishVersionWithCode(stores, fn, "", "stale-revision", "us-east-1")
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "PreconditionFailedException" {
		t.Fatalf("stale RevisionId: expected PreconditionFailedException, got %v", err)
	}
}

// TestPublishVersionAcceptsMatchingRevisionId pins the positive path: a
// matching (or omitted) RevisionId publishes the version.
func TestPublishVersionAcceptsMatchingRevisionId(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn := publishVersionFixture(t, svc, stores, "publish-match")

	version, err := svc.publishVersionWithCode(stores, fn, "", fn.RevisionId, "us-east-1")
	if err != nil {
		t.Fatalf("publish with matching RevisionId: %v", err)
	}
	if version == nil || version.Version == "" {
		t.Fatalf("expected published version, got %+v", version)
	}

	// An omitted RevisionId performs no precondition check.
	if _, err := svc.publishVersionWithCode(stores, fn, "", "", "us-east-1"); err != nil {
		t.Fatalf("publish with omitted RevisionId: %v", err)
	}
}
