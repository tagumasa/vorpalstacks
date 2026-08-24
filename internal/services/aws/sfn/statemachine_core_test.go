package sfn

import (
	"context"
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
