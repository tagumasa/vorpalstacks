package authorization

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/aws/iam"
)

// The Root-session tests only exercise the session path; permanent-key
// lookups must fail so Authorize falls through to authorizeSession.
type rootTestAccessKeys struct {
	iam.AccessKeyStoreInterface
}

func (rootTestAccessKeys) Get(string) (*iam.AccessKey, error) {
	return nil, errors.New("not a permanent IAM key")
}

type rootTestIAMStore struct {
	iam.IAMStoreInterface
}

func (rootTestIAMStore) AccessKeys() iam.AccessKeyStoreInterface {
	return rootTestAccessKeys{}
}

type rootTestResolver struct {
	creds *auth.SessionCredentials
}

func (r *rootTestResolver) ResolveSession(string) (*auth.SessionCredentials, error) {
	return r.creds, nil
}

const rootTaskPolicy = `{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": ["iam:GetLoginProfile", "iam:CreateLoginProfile"],
    "Resource": "*"
  }]
}`

func newRootSessionAuthorizer(t *testing.T, creds *auth.SessionCredentials) *Authorizer {
	t.Helper()
	a := NewAuthorizer(rootTestIAMStore{}, &rootTestResolver{creds: creds})
	t.Cleanup(a.Stop)
	return a
}

func rootSessionRequest(t *testing.T, operation string) (*request.RequestContext, *request.ParsedRequest, *http.Request) {
	t.Helper()
	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	parsedReq := &request.ParsedRequest{
		Operation:   operation,
		Parameters:  map[string]interface{}{},
		AccessKeyID: "ASIAROOTSESSIONTEST",
	}
	r, _ := http.NewRequest("POST", "http://localhost/"+operation, nil)
	return reqCtx, parsedReq, r
}

// An AssumeRoot session must be allowed only for the actions its task
// policy grants, not bypass evaluation entirely.
func TestRootSessionScopedByTaskPolicy(t *testing.T) {
	a := newRootSessionAuthorizer(t, &auth.SessionCredentials{
		AccessKeyID:   "ASIAROOTSESSIONTEST",
		PrincipalType: "Root",
		PrincipalArn:  "arn:aws:iam::123456789012:root",
		Policy:        rootTaskPolicy,
	})

	reqCtx, parsedReq, r := rootSessionRequest(t, "GetLoginProfile")
	allowed, err := a.Authorize(context.Background(), reqCtx, parsedReq, "iam", r)
	if err != nil {
		t.Fatal(err)
	}
	if !allowed {
		t.Fatal("task-policy action GetLoginProfile denied for Root session")
	}

	reqCtx, parsedReq, r = rootSessionRequest(t, "DeleteUser")
	allowed, err = a.Authorize(context.Background(), reqCtx, parsedReq, "iam", r)
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Fatal("action outside the task policy allowed for Root session")
	}
}

// A Root session without any task policy must be denied entirely — it must
// not fall back to the unrestricted root bypass.
func TestRootSessionWithoutTaskPolicyDenied(t *testing.T) {
	a := newRootSessionAuthorizer(t, &auth.SessionCredentials{
		AccessKeyID:   "ASIAROOTSESSIONTEST",
		PrincipalType: "Root",
		PrincipalArn:  "arn:aws:iam::123456789012:root",
	})

	reqCtx, parsedReq, r := rootSessionRequest(t, "GetLoginProfile")
	allowed, err := a.Authorize(context.Background(), reqCtx, parsedReq, "iam", r)
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Fatal("Root session without a task policy was allowed")
	}
}
