package authorization

import (
	"context"
	"errors"
	"strings"
	"testing"

	iamstore "vorpalstacks/internal/store/aws/iam"
)

// The Task-role tests exercise the Task Credentials authorisation chain
// over store fakes: the trust leg through the role's assume-role policy
// document, the permission leg through the role's inline policies. The
// attached-policy surface returns none, which is the same fail-closed
// posture an empty attachment list produces in the real store.

const statesTrustPolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"states.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
const lambdaTrustPolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
const invokeAllowPolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"lambda:InvokeFunction","Resource":"arn:aws:lambda:us-east-1:123456789012:function:Echo"}]}`
const sendAllowPolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"sqs:SendMessage","Resource":"arn:aws:sqs:us-east-1:123456789012:queue-a"}]}`

type taskRoleRoles struct {
	iamstore.RoleStoreInterface
	trust map[string]string
}

func (r *taskRoleRoles) GetAssumeRolePolicyDocument(roleName string) (string, error) {
	doc, ok := r.trust[roleName]
	if !ok {
		return "", errors.New("role not found")
	}
	return doc, nil
}

type taskRoleInline struct {
	iamstore.InlinePolicyStoreInterface
	docs map[string]string
}

func (i *taskRoleInline) List(principalType, principalName string) ([]string, error) {
	names := make([]string, 0, len(i.docs))
	for name := range i.docs {
		names = append(names, name)
	}
	return names, nil
}

func (i *taskRoleInline) Get(principalType, principalName, policyName string) (*iamstore.InlinePolicy, error) {
	doc, ok := i.docs[policyName]
	if !ok {
		return nil, errors.New("inline policy not found")
	}
	return &iamstore.InlinePolicy{PolicyDocument: doc}, nil
}

type taskRoleAttached struct {
	iamstore.AttachedPolicyStoreInterface
}

func (taskRoleAttached) ListAttachedPolicies(string, string) ([]string, error) {
	return nil, nil
}

type taskRoleIAMStore struct {
	iamstore.IAMStoreInterface
	roles  *taskRoleRoles
	inline *taskRoleInline
}

func (s *taskRoleIAMStore) Roles() iamstore.RoleStoreInterface { return s.roles }

func (s *taskRoleIAMStore) InlinePolicies() iamstore.InlinePolicyStoreInterface { return s.inline }

func (s *taskRoleIAMStore) AttachedPolicies() iamstore.AttachedPolicyStoreInterface {
	return taskRoleAttached{}
}

func newTaskRoleAuthoriser(t *testing.T, trust map[string]string, inline map[string]string) *TaskRoleAuthoriser {
	t.Helper()
	return NewTaskRoleAuthoriser(&taskRoleIAMStore{
		roles:  &taskRoleRoles{trust: trust},
		inline: &taskRoleInline{docs: inline},
	}, "123456789012")
}

// A role that trusts the States service principal and carries a policy
// allowing the integration action passes both legs of the chain.
func TestTaskRoleAuthoriserAllowsPermittedAction(t *testing.T) {
	a := newTaskRoleAuthoriser(t,
		map[string]string{"TaskRole": statesTrustPolicy},
		map[string]string{"invoke": invokeAllowPolicy})
	err := a.AuthoriseTaskCredentials(context.Background(),
		"arn:aws:iam::123456789012:role/TaskRole",
		"lambda:InvokeFunction",
		"arn:aws:lambda:us-east-1:123456789012:function:Echo")
	if err != nil {
		t.Fatalf("expected the permitted action to pass, got %v", err)
	}
}

// A role with no policies allows nothing: the evaluator's default-deny is
// the fail-closed posture for unauthorised tasks.
func TestTaskRoleAuthoriserDeniesRoleWithoutPolicies(t *testing.T) {
	a := newTaskRoleAuthoriser(t,
		map[string]string{"TaskRole": statesTrustPolicy},
		nil)
	err := a.AuthoriseTaskCredentials(context.Background(),
		"arn:aws:iam::123456789012:role/TaskRole",
		"lambda:InvokeFunction",
		"arn:aws:lambda:us-east-1:123456789012:function:Echo")
	if err == nil {
		t.Fatal("expected a role without policies to deny the action")
	}
}

// A policy scoped to another action does not authorise the requested one.
func TestTaskRoleAuthoriserDeniesUnscopedAction(t *testing.T) {
	a := newTaskRoleAuthoriser(t,
		map[string]string{"TaskRole": statesTrustPolicy},
		map[string]string{"send": sendAllowPolicy})
	err := a.AuthoriseTaskCredentials(context.Background(),
		"arn:aws:iam::123456789012:role/TaskRole",
		"lambda:InvokeFunction",
		"arn:aws:lambda:us-east-1:123456789012:function:Echo")
	if err == nil {
		t.Fatal("expected a policy scoped to sqs:SendMessage to deny lambda:InvokeFunction")
	}
}

// A role whose trust policy names another service principal fails the
// trust leg before the permission leg runs.
func TestTaskRoleAuthoriserRejectsForeignTrust(t *testing.T) {
	a := newTaskRoleAuthoriser(t,
		map[string]string{"TaskRole": lambdaTrustPolicy},
		map[string]string{"invoke": invokeAllowPolicy})
	err := a.AuthoriseTaskCredentials(context.Background(),
		"arn:aws:iam::123456789012:role/TaskRole",
		"lambda:InvokeFunction",
		"arn:aws:lambda:us-east-1:123456789012:function:Echo")
	if err == nil {
		t.Fatal("expected a role not trusting states.amazonaws.com to be rejected")
	}
	if !strings.Contains(err.Error(), "cannot assume role") {
		t.Fatalf("expected the failure to name the assumption, got %v", err)
	}
}

// A role the store cannot resolve fails the trust leg as missing.
func TestTaskRoleAuthoriserRejectsMissingRole(t *testing.T) {
	a := newTaskRoleAuthoriser(t, nil, nil)
	err := a.AuthoriseTaskCredentials(context.Background(),
		"arn:aws:iam::123456789012:role/NoSuchRole",
		"lambda:InvokeFunction",
		"arn:aws:lambda:us-east-1:123456789012:function:Echo")
	if err == nil {
		t.Fatal("expected a missing role to be rejected")
	}
}

// A cross-account role passes the trust leg — the IAM validator's
// established posture for cross-account ARNs — and the permission leg
// still applies, so an unauthorised cross-account task role denies.
func TestTaskRoleAuthoriserCrossAccountRoleStillScopedByPolicies(t *testing.T) {
	a := newTaskRoleAuthoriser(t, nil, nil)
	err := a.AuthoriseTaskCredentials(context.Background(),
		"arn:aws:iam::999999999999:role/ForeignRole",
		"lambda:InvokeFunction",
		"arn:aws:lambda:us-east-1:123456789012:function:Echo")
	if err == nil {
		t.Fatal("expected a cross-account role without policies to deny the action")
	}
	if strings.Contains(err.Error(), "cannot assume role") {
		t.Fatalf("cross-account trust must pass the trust leg, got %v", err)
	}
}
