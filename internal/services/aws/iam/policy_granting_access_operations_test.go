package iam

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
)

func TestPolicyGrantsServiceNamespace(t *testing.T) {
	parse := func(t *testing.T, doc string) *policy.Document {
		t.Helper()
		parsed, err := policy.ParseDocument(doc)
		if err != nil {
			t.Fatalf("parse policy document: %v", err)
		}
		return parsed
	}

	s3Doc := parse(t, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:*","Resource":"*"}]}`)
	if !policyGrantsServiceNamespace(s3Doc, "s3") {
		t.Error("an Allow on s3:* must grant the s3 namespace")
	}
	if policyGrantsServiceNamespace(s3Doc, "ec2") {
		t.Error("an Allow on s3:* must not grant the ec2 namespace")
	}

	wildcardDoc := parse(t, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`)
	for _, ns := range []string{"s3", "iam", "logs"} {
		if !policyGrantsServiceNamespace(wildcardDoc, ns) {
			t.Errorf("a global wildcard Allow must grant namespace %s", ns)
		}
	}

	denyDoc := parse(t, `{"Version":"2012-10-17","Statement":[{"Effect":"Deny","Action":"s3:*","Resource":"*"}]}`)
	if policyGrantsServiceNamespace(denyDoc, "s3") {
		t.Error("a Deny-only document must not grant any namespace")
	}

	notActionDoc := parse(t, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","NotAction":"iam:*","Resource":"*"}]}`)
	if !policyGrantsServiceNamespace(notActionDoc, "s3") {
		t.Error("an Allow with NotAction grants every namespace at namespace granularity")
	}
}

// The SDK rejects inputs missing required members client-side, so the
// server-side required-parameter checks are pinned here.
func TestListPoliciesGrantingServiceAccessRequiredParams(t *testing.T) {
	s := &IAMService{}

	if _, err := s.ListPoliciesGrantingServiceAccess(context.Background(), nil, &request.ParsedRequest{
		Parameters: map[string]interface{}{},
	}); err == nil {
		t.Error("a missing Arn must be rejected")
	}

	if _, err := s.ListPoliciesGrantingServiceAccess(context.Background(), nil, &request.ParsedRequest{
		Parameters: map[string]interface{}{
			"Arn": "arn:aws:iam::000000000000:user/nobody",
		},
	}); err == nil {
		t.Error("a missing ServiceNamespaces list must be rejected")
	}
}
