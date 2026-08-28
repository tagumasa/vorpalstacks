package iam

import (
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// The typed AWS SDK validates required members client-side, so the
// server-side rejection of an invalid Granularity (optional member) and
// the ServiceNamespace pattern rules are pinned here as unit tests.

func TestGenerateServiceLastAccessedDetailsCoreInvalidGranularity(t *testing.T) {
	s := NewIAMService("000000000000")
	for _, granularity := range []string{"P30D", "P7D", "weekly", "service_level"} {
		_, err := s.generateServiceLastAccessedDetailsCore("us-east-1", nil, &GenerateServiceLastAccessedDetailsInput{
			Arn:         "arn:aws:iam::000000000000:user/example",
			Granularity: granularity,
		})
		if err == nil {
			t.Fatalf("granularity %q must be rejected", granularity)
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("granularity %q: expected AWSError, got %T", granularity, err)
		}
		if awsErr.Code != "InvalidInput" {
			t.Fatalf("granularity %q: got code %q, want InvalidInput", granularity, awsErr.Code)
		}
	}
}

func TestGenerateServiceLastAccessedDetailsCoreEmptyArnRejected(t *testing.T) {
	s := NewIAMService("000000000000")
	_, err := s.generateServiceLastAccessedDetailsCore("us-east-1", nil, &GenerateServiceLastAccessedDetailsInput{})
	if err == nil {
		t.Fatal("an empty Arn must be rejected")
	}
}

// arnType is @length(20, 2048) and jobIDType is a fixed @length(36, 36);
// both bounds are enforced server-side.
func TestServiceLastAccessedLengthValidation(t *testing.T) {
	s := NewIAMService("000000000000")

	for _, arn := range []string{"a", "arn:aws:iam::000000000000:user/" + strings.Repeat("x", 2100)} {
		_, err := s.generateServiceLastAccessedDetailsCore("us-east-1", nil, &GenerateServiceLastAccessedDetailsInput{Arn: arn})
		if err == nil {
			t.Fatalf("Arn of length %d must be rejected", len(arn))
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("Arn length %d: expected AWSError, got %T", len(arn), err)
		}
		if awsErr.Code != "InvalidInput" {
			t.Fatalf("Arn length %d: got code %q, want InvalidInput", len(arn), awsErr.Code)
		}
	}

	for _, jobID := range []string{"short", strings.Repeat("j", 37)} {
		_, err := s.getServiceLastAccessedDetailsCore(nil, jobID)
		if err == nil {
			t.Fatalf("JobId of length %d must be rejected", len(jobID))
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("JobId length %d: expected AWSError, got %T", len(jobID), err)
		}
		if awsErr.Code != "InvalidInput" {
			t.Fatalf("JobId length %d: got code %q, want InvalidInput", len(jobID), awsErr.Code)
		}
	}
}

// A report generated for a group or policy ARN covers the member
// entities' activity, and the per-service entity recorded with the last
// access attempt is the accessing member's ARN — not the report ARN.
func TestReportPrincipalsGroupAndPolicy(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")

	for _, name := range []string{"group-member-a", "group-member-b"} {
		if _, err := store.Users().Create(name, "/", "123456789012", nil); err != nil {
			t.Fatalf("create user %s: %v", name, err)
		}
	}
	if _, err := store.Groups().Create("report-group", "/", "123456789012"); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if err := store.UserGroups().AddUserToGroup("group-member-a", "report-group"); err != nil {
		t.Fatalf("add member: %v", err)
	}

	groupPrincipals := reportPrincipals(store, "arn:aws:iam::123456789012:group/report-group")
	if len(groupPrincipals) != 1 || groupPrincipals[0].arn != "arn:aws:iam::123456789012:user/group-member-a" {
		t.Fatalf("group principals: got %+v, want the member user only", groupPrincipals)
	}

	if _, err := store.Roles().Create("policy-user-role", "/", "123456789012", "{}", "", 3600, nil); err != nil {
		t.Fatalf("create role: %v", err)
	}
	policyArn := "arn:aws:iam::123456789012:policy/report-policy"
	if err := store.AttachedPolicies().Attach(PrincipalTypeRole, "policy-user-role", policyArn); err != nil {
		t.Fatalf("attach policy: %v", err)
	}

	policyPrincipals := reportPrincipals(store, policyArn)
	if len(policyPrincipals) != 1 || policyPrincipals[0].arn != "arn:aws:iam::123456789012:role/policy-user-role" {
		t.Fatalf("policy principals: got %+v, want the attached role only", policyPrincipals)
	}

	userPrincipals := reportPrincipals(store, "arn:aws:iam::123456789012:user/group-member-a")
	if len(userPrincipals) != 1 || userPrincipals[0].arn != "arn:aws:iam::123456789012:user/group-member-a" {
		t.Fatalf("user principals: got %+v, want the entity itself", userPrincipals)
	}
}

func TestGetServiceLastAccessedDetailsWithEntitiesCoreNamespaceValidation(t *testing.T) {
	s := NewIAMService("000000000000")

	if _, err := s.getServiceLastAccessedDetailsWithEntitiesCore(nil, "job", "", "", 100); err == nil {
		t.Fatal("an empty ServiceNamespace must be rejected")
	} else if awsErr, ok := err.(*awserrors.AWSError); !ok {
		t.Fatalf("expected AWSError, got %T", err)
	} else if awsErr.Code != "InvalidInput" || !strings.Contains(awsErr.Message, "ServiceNamespace") {
		t.Fatalf("got code %q message %q, want InvalidInput missing ServiceNamespace", awsErr.Code, awsErr.Message)
	}

	for _, ns := range []string{"bad namespace", "ns!", "no/slash", strings.Repeat("a", 65)} {
		_, err := s.getServiceLastAccessedDetailsWithEntitiesCore(nil, "job", ns, "", 100)
		if err == nil {
			t.Fatalf("ServiceNamespace %q must be rejected", ns)
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("ServiceNamespace %q: expected AWSError, got %T", ns, err)
		}
		if awsErr.Code != "InvalidInput" {
			t.Fatalf("ServiceNamespace %q: got code %q, want InvalidInput", ns, awsErr.Code)
		}
	}
}
