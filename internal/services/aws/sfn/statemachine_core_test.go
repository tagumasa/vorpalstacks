package sfn

import (
	"context"
	"fmt"
	"testing"

	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func newCreateTestStore(t *testing.T) *sfnstore.StepFunctionStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")
}

func TestCreateStateMachineIdempotentRetry(t *testing.T) {
	store := newCreateTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	in := CreateStateMachineInput{
		Name:       "idem-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	}
	first, err := svc.createStateMachineCore(ctx, store, in)
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	// A retry differing only in roleArn is still an idempotent request of
	// the previous one.
	retry := in
	retry.RoleArn = "arn:aws:iam::000000000000:role/other"
	second, err := svc.createStateMachineCore(ctx, store, retry)
	if err != nil {
		t.Fatalf("idempotent retry: %v", err)
	}
	if second.StateMachineArn != first.StateMachineArn {
		t.Fatalf("retry returned %s, want the original %s", second.StateMachineArn, first.StateMachineArn)
	}
	if !second.CreationDate.Equal(first.CreationDate) {
		t.Fatalf("retry returned a different creationDate: %v vs %v", second.CreationDate, first.CreationDate)
	}
}

func TestCreateStateMachineSameNameDifferentDefinitionConflicts(t *testing.T) {
	store := newCreateTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	if _, err := svc.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:       "clash-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	}); err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err := svc.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:       "clash-sm",
		Definition: `{"StartAt":"B","States":{"B":{"Type":"Succeed"}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	})
	requireAWSCode(t, err, "StateMachineAlreadyExists")
}

func TestCreateStateMachineIdempotentWithPublish(t *testing.T) {
	store := newCreateTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	in := CreateStateMachineInput{
		Name:               "publish-sm",
		Definition:         `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:            "arn:aws:iam::000000000000:role/sm",
		Publish:            true,
		VersionDescription: "initial",
	}
	first, err := svc.createStateMachineCore(ctx, store, in)
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	if first.StateMachineVersionArn == "" {
		t.Fatal("publishing create returned no version ARN")
	}

	second, err := svc.createStateMachineCore(ctx, store, in)
	if err != nil {
		t.Fatalf("idempotent retry: %v", err)
	}
	if second.StateMachineVersionArn != first.StateMachineVersionArn {
		t.Fatalf("retry version %s, want the original %s", second.StateMachineVersionArn, first.StateMachineVersionArn)
	}

	// publish=false on the retry no longer matches the original request.
	noPublish := in
	noPublish.Publish = false
	_, err = svc.createStateMachineCore(ctx, store, noPublish)
	requireAWSCode(t, err, "StateMachineAlreadyExists")
}

func TestValidateStateMachineDefinitionResultAllowsWarnings(t *testing.T) {
	svc := &StepFunctionService{}

	// A Pass Result that looks like a path yields only the documented
	// PASS_RESULT_IS_STATIC warning; the definition can still be created,
	// so the result stays OK. Severity WARNING keeps the warning in the
	// response.
	warningOnly := `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":"$.payload","End":true}}}`
	resp, err := svc.validateStateMachineDefinitionCore(ValidateStateMachineDefinitionInput{
		Definition: warningOnly,
		Severity:   "WARNING",
	})
	if err != nil {
		t.Fatalf("warning-only definition: %v", err)
	}
	if resp["result"] != "OK" {
		t.Fatalf("warning-only definition must stay OK, got %v", resp["result"])
	}
	diagnostics := resp["diagnostics"].([]map[string]string)
	if len(diagnostics) == 0 {
		t.Fatal("expected the warning diagnostic to be reported")
	}
	foundWarningCode := false
	for _, d := range diagnostics {
		if d["severity"] == "WARNING" && d["code"] == "PASS_RESULT_IS_STATIC" {
			foundWarningCode = true
			if d["location"] != "/States/P/Result" {
				t.Errorf("warning location = %q, want /States/P/Result", d["location"])
			}
		}
	}
	if !foundWarningCode {
		t.Fatalf("PASS_RESULT_IS_STATIC warning missing: %v", diagnostics)
	}

	// The default severity filters warnings out without failing the
	// definition.
	resp, err = svc.validateStateMachineDefinitionCore(ValidateStateMachineDefinitionInput{
		Definition: warningOnly,
	})
	if err != nil {
		t.Fatalf("warning-only definition at default severity: %v", err)
	}
	if resp["result"] != "OK" {
		t.Fatalf("warning-only definition must stay OK at default severity, got %v", resp["result"])
	}
	if diags := resp["diagnostics"].([]map[string]string); len(diags) != 0 {
		t.Fatalf("default severity filters warnings, got %v", diags)
	}

	// Error diagnostics still fail the definition.
	errorDef := `{"StartAt":"Missing","States":{}}`
	resp, err = svc.validateStateMachineDefinitionCore(ValidateStateMachineDefinitionInput{
		Definition: errorDef,
	})
	if err != nil {
		t.Fatalf("error definition: %v", err)
	}
	if resp["result"] != "FAIL" {
		t.Fatalf("error definition must FAIL, got %v", resp["result"])
	}

	// A warning alongside an error also fails.
	mixed := `{"StartAt":"A","States":{"A":{"Type":"Fail"},"B":{"Type":"Pass","Next":"C"}}}`
	resp, err = svc.validateStateMachineDefinitionCore(ValidateStateMachineDefinitionInput{
		Definition: mixed,
	})
	if err != nil {
		t.Fatalf("mixed definition: %v", err)
	}
	if resp["result"] != "FAIL" {
		t.Fatalf("definition with errors must FAIL, got %v", resp["result"])
	}
}

// TestListStateMachinesValidatesMaxResults pins the Core-level bound so
// both planes enforce the same range: the HTTP handler pre-validates via
// parsePageLimit, and the admin gRPC handler forwards raw — the Core is
// where they meet.
func TestListStateMachinesValidatesMaxResults(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	for _, bad := range []int32{1001, -1} {
		_, err := svc.listStateMachinesCore(ctx, store, ListStateMachinesInput{MaxResults: bad})
		requireAWSCode(t, err, "ValidationException")
	}
	if _, err := svc.listStateMachinesCore(ctx, store, ListStateMachinesInput{MaxResults: 0}); err != nil {
		t.Errorf("default page size rejected: %v", err)
	}
}

// stubRoleProvider backs the Core role-validation pin.
type stubRoleProvider struct {
	exists    map[string]bool
	assumeFor map[string]string
}

func (s *stubRoleProvider) GetAssumeRolePolicyDocument(roleName string) (string, error) {
	if doc, ok := s.assumeFor[roleName]; ok {
		return doc, nil
	}
	return "", fmt.Errorf("role not found: %s", roleName)
}

func (s *stubRoleProvider) Exists(roleName string) bool {
	return s.exists[roleName]
}

// TestCoreValidatesStateMachineRole pins the #29 single-path contract:
// role validation runs inside the Core for every plane, not in the HTTP
// handler — with the provider injected, a nonexistent role fails
// creation identically however the Core is reached, and a nil provider
// (not injected) leaves validation skipped.
func TestCoreValidatesStateMachineRole(t *testing.T) {
	store := newCreateTestStore(t)
	ctx := context.Background()

	withProvider := NewStepFunctionService(nil, "000000000000")
	withProvider.SetRoleProvider(&stubRoleProvider{
		exists:    map[string]bool{"good-role": true},
		assumeFor: map[string]string{"good-role": `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"states.amazonaws.com"},"Action":"sts:AssumeRole"}]}`},
	})
	_, err := withProvider.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:       "role-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/missing-role",
	})
	if err == nil {
		t.Fatal("the Core must reject a role the provider does not carry")
	}
	requireAWSCode(t, err, "ValidationException")

	if _, err := withProvider.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:       "role-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/good-role",
	}); err != nil {
		t.Fatalf("a carried role must create: %v", err)
	}

	// Without the provider injected the Core proceeds (the validator is
	// unavailable, not wrong) — the historical unit-test plane.
	noProvider := &StepFunctionService{}
	if _, err := noProvider.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:       "role-sm-2",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/any",
	}); err != nil {
		t.Fatalf("nil provider must skip validation: %v", err)
	}
}

// TestValidateDefinitionResultSpansTruncatedPage pins the result
// semantics: the OK/FAIL verdict is computed over the definition's full
// diagnostic set, while maxResults bounds only the returned diagnostics
// page — an ERROR beyond the cutoff must still fail a WARNING-severity
// validation.
func TestValidateDefinitionResultSpansTruncatedPage(t *testing.T) {
	svc, _ := newRecoveryService(t)
	def := `{"StartAt":"A1","States":{` +
		`"A1":{"Type":"Pass","Result":"$.looks.like.a.path","Next":"A2"},` +
		`"A2":{"Type":"Pass","Result":"$.looks.like.a.path","Next":"Zbad"},` +
		`"Zbad":{"Type":"Pass","Nope":1,"End":true}}}`
	out, err := svc.validateStateMachineDefinitionCore(ValidateStateMachineDefinitionInput{
		Definition: def, SMType: "STANDARD", Severity: "WARNING", MaxResults: 2})
	if err != nil {
		t.Fatalf("validate core: %v", err)
	}
	if out["result"] != "FAIL" {
		t.Errorf("result = %v with diagnostics %v, want FAIL — the ERROR beyond the cutoff still fails the definition", out["result"], out["diagnostics"])
	}
	if out["truncated"] != true {
		t.Errorf("truncated = %v, want true", out["truncated"])
	}
}
