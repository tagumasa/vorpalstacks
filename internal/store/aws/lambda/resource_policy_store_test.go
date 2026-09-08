package lambda

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// resourcePolicyStore builds a FunctionStore over fresh storage with one
// function created, for the resource-policy operation tests.
func resourcePolicyStore(t *testing.T) *FunctionStore {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	st, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("get storage: %v", err)
	}
	s := NewFunctionStore(st, "000000000000", "us-east-1")
	if _, err := s.Create(&Function{
		FunctionName: "fn",
		Runtime:      RuntimeNodejs22X,
		Role:         "arn:aws:iam::000000000000:role/lambda",
	}); err != nil {
		t.Fatalf("create function: %v", err)
	}
	return s
}

func policyStatements(sids ...string) []FunctionPolicy {
	policies := make([]FunctionPolicy, 0, len(sids))
	for _, sid := range sids {
		policies = append(policies, FunctionPolicy{
			Id:        sid,
			Principal: "*",
			Action:    "lambda:InvokeFunction",
			Resource:  "arn:aws:lambda:us-east-1:000000000000:function:fn",
			Raw:       `{"Sid":"` + sid + `","Effect":"Allow","Principal":"*","Action":"lambda:InvokeFunction"}`,
		})
	}
	return policies
}

// TestSetResourcePolicyReplacesAndVersions pins the whole-document
// replacement semantics: every Put installs exactly the submitted
// statements under a fresh policy revision, and a Put carrying the current
// revision replaces the document again. Statements without a Sid receive a
// generated statement ID.
func TestSetResourcePolicyReplacesAndVersions(t *testing.T) {
	s := resourcePolicyStore(t)

	rev1, err := s.SetResourcePolicy("fn", "", policyStatements("s1", "s2"))
	if err != nil {
		t.Fatalf("set resource policy: %v", err)
	}
	if rev1 == "" {
		t.Fatal("first set returned an empty revision")
	}
	fn, err := s.Get("fn")
	if err != nil {
		t.Fatalf("get function: %v", err)
	}
	if len(fn.Policies) != 2 || fn.Policies[0].Id != "s1" || fn.Policies[1].Id != "s2" {
		t.Fatalf("policy statements %v, want s1 and s2", fn.Policies)
	}

	rev2, err := s.SetResourcePolicy("fn", rev1, policyStatements("s3"))
	if err != nil {
		t.Fatalf("set resource policy with current revision: %v", err)
	}
	if rev2 == rev1 {
		t.Fatal("second set did not rotate the policy revision")
	}
	fn, err = s.Get("fn")
	if err != nil {
		t.Fatalf("get function: %v", err)
	}
	if len(fn.Policies) != 1 || fn.Policies[0].Id != "s3" {
		t.Fatalf("policy statements after replacement %v, want only s3", fn.Policies)
	}

	noSid := []FunctionPolicy{{Raw: `{"Effect":"Allow","Action":"lambda:InvokeFunction"}`}}
	if _, err := s.SetResourcePolicy("fn", rev2, noSid); err != nil {
		t.Fatalf("set resource policy without Sids: %v", err)
	}
	fn, err = s.Get("fn")
	if err != nil {
		t.Fatalf("get function: %v", err)
	}
	if len(fn.Policies) != 1 || fn.Policies[0].Id == "" {
		t.Fatalf("Sid-less statement did not receive a generated ID: %v", fn.Policies)
	}
}

// TestSetResourcePolicyRevisionPrecondition pins the optimistic-locking
// contract: a non-empty expected revision that does not match the current
// policy revision fails without touching the stored policy.
func TestSetResourcePolicyRevisionPrecondition(t *testing.T) {
	s := resourcePolicyStore(t)

	rev, err := s.SetResourcePolicy("fn", "", policyStatements("s1"))
	if err != nil {
		t.Fatalf("set resource policy: %v", err)
	}

	if _, err := s.SetResourcePolicy("fn", "00000000-0000-0000-0000-000000000000", policyStatements("s2")); !errors.Is(err, ErrPolicyRevisionMismatch) {
		t.Fatalf("stale revision error = %v, want ErrPolicyRevisionMismatch", err)
	}
	fn, err := s.Get("fn")
	if err != nil {
		t.Fatalf("get function: %v", err)
	}
	if len(fn.Policies) != 1 || fn.Policies[0].Id != "s1" || fn.PolicyRevisionId != rev {
		t.Fatalf("failed precondition mutated the stored policy: %v rev %q", fn.Policies, fn.PolicyRevisionId)
	}

	// A first Put against a function with no policy carries no current
	// revision, so any non-empty expected revision is a mismatch.
	s2 := resourcePolicyStore(t)
	if _, err := s2.SetResourcePolicy("fn", "00000000-0000-0000-0000-000000000000", policyStatements("s1")); !errors.Is(err, ErrPolicyRevisionMismatch) {
		t.Fatalf("precondition on policy-less function error = %v, want ErrPolicyRevisionMismatch", err)
	}
}

// TestDeleteResourcePolicy pins the delete contract: the revision
// precondition applies, a successful delete clears the statements, and
// repeated deletes on a policy-less function succeed without change.
func TestDeleteResourcePolicy(t *testing.T) {
	s := resourcePolicyStore(t)

	rev, err := s.SetResourcePolicy("fn", "", policyStatements("s1"))
	if err != nil {
		t.Fatalf("set resource policy: %v", err)
	}

	if err := s.DeleteResourcePolicy("fn", "00000000-0000-0000-0000-000000000000"); !errors.Is(err, ErrPolicyRevisionMismatch) {
		t.Fatalf("stale revision error = %v, want ErrPolicyRevisionMismatch", err)
	}
	if err := s.DeleteResourcePolicy("fn", rev); err != nil {
		t.Fatalf("delete resource policy: %v", err)
	}
	fn, err := s.Get("fn")
	if err != nil {
		t.Fatalf("get function: %v", err)
	}
	if len(fn.Policies) != 0 {
		t.Fatalf("policies after delete: %v", fn.Policies)
	}

	// The delete is idempotent: a second delete is a no-op that leaves the
	// policy revision (and every other function field) untouched.
	if err := s.DeleteResourcePolicy("fn", ""); err != nil {
		t.Fatalf("second delete: %v", err)
	}
	again, err := s.Get("fn")
	if err != nil {
		t.Fatalf("reload after second delete: %v", err)
	}
	if again.PolicyRevisionId != fn.PolicyRevisionId || again.RevisionId != fn.RevisionId {
		t.Fatalf("second delete mutated revisions: policy %q→%q function %q→%q",
			fn.PolicyRevisionId, again.PolicyRevisionId, fn.RevisionId, again.RevisionId)
	}

	s2 := resourcePolicyStore(t)
	if err := s2.DeleteResourcePolicy("fn", ""); err != nil {
		t.Fatalf("delete on never-policied function: %v", err)
	}
}

// TestPermissionMutationsBumpPolicyRevision pins the single-revision
// invariant: the statement-level AddPermission and RemovePermission paths
// rotate the policy revision as well, so a revision held from
// GetResourcePolicy cannot silently overwrite a later AddPermission. The
// function's own RevisionId and LastModified stay untouched — the policy
// is versioned independently of the function configuration.
func TestPermissionMutationsBumpPolicyRevision(t *testing.T) {
	s := resourcePolicyStore(t)

	before, err := s.Get("fn")
	if err != nil {
		t.Fatalf("get function before mutations: %v", err)
	}

	if err := s.AddPolicyAtomically("fn", &FunctionPolicy{
		Id:        "s1",
		Principal: "s3.amazonaws.com",
		Action:    "lambda:InvokeFunction",
		Resource:  "arn:aws:lambda:us-east-1:000000000000:function:fn",
	}); err != nil {
		t.Fatalf("add permission: %v", err)
	}
	fn, err := s.Get("fn")
	if err != nil {
		t.Fatalf("get function: %v", err)
	}
	if fn.PolicyRevisionId == "" {
		t.Fatal("AddPermission did not set a policy revision")
	}
	if fn.RevisionId != before.RevisionId || !fn.LastModified.Equal(before.LastModified) {
		t.Fatal("AddPermission rotated the function revision")
	}

	if err := s.RemovePolicy("fn", "s1"); err != nil {
		t.Fatalf("remove permission: %v", err)
	}
	fn, err = s.Get("fn")
	if err != nil {
		t.Fatalf("get function: %v", err)
	}
	if len(fn.Policies) != 0 {
		t.Fatalf("statement not removed: %v", fn.Policies)
	}
	if fn.PolicyRevisionId == "" {
		t.Fatal("RemovePolicy lost the policy revision")
	}
	if fn.RevisionId != before.RevisionId || !fn.LastModified.Equal(before.LastModified) {
		t.Fatal("RemovePolicy rotated the function revision")
	}
}
