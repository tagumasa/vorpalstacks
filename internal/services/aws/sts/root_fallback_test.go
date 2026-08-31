package sts

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
)

func newStsTestEnv(t *testing.T) (*STSService, *request.RequestContext) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	return NewSTSService(), reqCtx
}

// Session-issuing operations must reject callers whose identity cannot be
// resolved instead of defaulting them to the account root ARN: a session
// minted with a root principal ARN bypasses policy evaluation in the
// authorizer.
func TestSessionOperationsRejectUnresolvableCaller(t *testing.T) {
	t.Setenv("TEST_MODE", "false")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	svc, reqCtx := newStsTestEnv(t)
	noAuth := &request.ParsedRequest{Parameters: map[string]interface{}{}}

	if _, _, err := svc.resolveCallerArnOrReject(reqCtx, callerAccessKeyID(noAuth)); !errors.Is(err, ErrInvalidClientTokenId) {
		t.Fatalf("resolveCallerArnOrReject: got %v, want ErrInvalidClientTokenId", err)
	}
	if _, err := svc.GetSessionToken(context.Background(), reqCtx, noAuth); !errors.Is(err, ErrInvalidClientTokenId) {
		t.Fatalf("GetSessionToken: got %v, want ErrInvalidClientTokenId", err)
	}
	if _, err := svc.GetFederationToken(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Name": "feduser",
	}}); !errors.Is(err, ErrInvalidClientTokenId) {
		t.Fatalf("GetFederationToken: got %v, want ErrInvalidClientTokenId", err)
	}
	if _, err := svc.GetWebIdentityToken(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Audience.member.1": "audience",
		"SigningAlgorithm":  "RS256",
	}}); !errors.Is(err, ErrInvalidClientTokenId) {
		t.Fatalf("GetWebIdentityToken: got %v, want ErrInvalidClientTokenId", err)
	}
	if _, err := svc.AssumeRole(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"RoleArn":         "arn:aws:iam::000000000000:role/testrole",
		"RoleSessionName": "session",
	}}); !errors.Is(err, ErrInvalidClientTokenId) {
		t.Fatalf("AssumeRole: got %v, want ErrInvalidClientTokenId", err)
	}
}

// The configured root credential resolves to the root principal by
// construction: the signature middleware's static verifier accepts exactly
// that key, and its registration in the IAM access-key store is not
// guaranteed.
func TestResolveCallerRecognisesConfiguredRootKey(t *testing.T) {
	t.Setenv("TEST_MODE", "false")
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIACONFIGUREDROOT")
	svc, reqCtx := newStsTestEnv(t)

	req := &request.ParsedRequest{
		Headers:    http.Header{"X-Amz-Access-Key": []string{"AKIACONFIGUREDROOT"}},
		Parameters: map[string]interface{}{},
	}
	arn, name, err := svc.resolveCallerArnOrReject(reqCtx, callerAccessKeyID(req))
	if err != nil {
		t.Fatal(err)
	}
	if arn != "arn:aws:iam::000000000000:root" {
		t.Fatalf("got ARN %s, want the root ARN", arn)
	}
	if name != iam.RootUserName {
		t.Fatalf("got name %s, want %s", name, iam.RootUserName)
	}

	// Any other unregistered key is still rejected in production mode.
	other := &request.ParsedRequest{
		Headers:    http.Header{"X-Amz-Access-Key": []string{"AKIAUNKNOWN"}},
		Parameters: map[string]interface{}{},
	}
	if _, _, err := svc.resolveCallerArnOrReject(reqCtx, callerAccessKeyID(other)); !errors.Is(err, ErrInvalidClientTokenId) {
		t.Fatalf("unknown key: got %v, want ErrInvalidClientTokenId", err)
	}
}

// TEST_MODE runs the documented SDK-test procedure with signature
// verification disabled and dummy credentials no store knows; unresolvable
// callers fall back to root there, mirroring GetCallerIdentity.
func TestResolveCallerTestModeFallsBackToRoot(t *testing.T) {
	t.Setenv("TEST_MODE", "true")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	svc, reqCtx := newStsTestEnv(t)

	req := &request.ParsedRequest{
		Headers:    http.Header{"X-Amz-Access-Key": []string{"test"}},
		Parameters: map[string]interface{}{},
	}
	arn, name, err := svc.resolveCallerArnOrReject(reqCtx, callerAccessKeyID(req))
	if err != nil {
		t.Fatal(err)
	}
	if arn != "arn:aws:iam::000000000000:root" || name != iam.RootUserName {
		t.Fatalf("got (%s, %s), want the root principal", arn, name)
	}
}
