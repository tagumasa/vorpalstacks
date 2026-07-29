package secretsmanager

import (
	"testing"

	"vorpalstacks/internal/common/iam/policy"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

func TestParsePolicyDocument_Valid(t *testing.T) {
	valid := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:root"},"Action":"secretsmanager:GetSecretValue","Resource":"*"}]}`
	doc, err := parsePolicyDocument(valid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if len(doc.Statement) != 1 {
		t.Errorf("expected 1 statement, got %d", len(doc.Statement))
	}
}

func TestParsePolicyDocument_Empty(t *testing.T) {
	_, err := parsePolicyDocument("")
	if err == nil {
		t.Fatal("expected error for empty policy")
	}
}

func TestParsePolicyDocument_MalformedJSON(t *testing.T) {
	_, err := parsePolicyDocument("{not valid json")
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestParsePolicyDocument_NoStatements(t *testing.T) {
	// Empty Statement array is syntactically valid — parsePolicyDocument
	// only checks JSON validity, not content constraints. Content checks
	// are the responsibility of validatePolicyDocument (used by
	// ValidateResourcePolicy) or BlockPublicPolicy (PutResourcePolicy).
	doc, err := parsePolicyDocument(`{"Version":"2012-10-17","Statement":[]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Statement) != 0 {
		t.Errorf("expected 0 statements, got %d", len(doc.Statement))
	}
}

func TestIsPolicyPublic_PrincipalStar(t *testing.T) {
	doc := &policy.Document{
		Version: "2012-10-17",
		Statement: []policy.Statement{
			{
				Effect:    policy.EffectAllow,
				Principal: &policy.Principal{Everyone: true},
				Action:    policy.ActionList{"secretsmanager:GetSecretValue"},
			},
		},
	}
	if !isPolicyPublic(doc) {
		t.Error("expected policy with Principal:* to be public")
	}
}

func TestIsPolicyPublic_PrincipalAWSStar(t *testing.T) {
	doc := &policy.Document{
		Version: "2012-10-17",
		Statement: []policy.Statement{
			{
				Effect:    policy.EffectAllow,
				Principal: &policy.Principal{AWS: policy.StringList{"*"}},
				Action:    policy.ActionList{"secretsmanager:GetSecretValue"},
			},
		},
	}
	if !isPolicyPublic(doc) {
		t.Error("expected policy with Principal AWS:* to be public")
	}
}

func TestIsPolicyPublic_SpecificPrincipal(t *testing.T) {
	doc := &policy.Document{
		Version: "2012-10-17",
		Statement: []policy.Statement{
			{
				Effect:    policy.EffectAllow,
				Principal: &policy.Principal{AWS: policy.StringList{"arn:aws:iam::123456789012:root"}},
				Action:    policy.ActionList{"secretsmanager:GetSecretValue"},
			},
		},
	}
	if isPolicyPublic(doc) {
		t.Error("expected policy with specific principal to NOT be public")
	}
}

func TestIsPolicyPublic_PublicWithCondition(t *testing.T) {
	doc := &policy.Document{
		Version: "2012-10-17",
		Statement: []policy.Statement{
			{
				Effect:    policy.EffectAllow,
				Principal: &policy.Principal{Everyone: true},
				Action:    policy.ActionList{"secretsmanager:GetSecretValue"},
				Condition: policy.ConditionMap{
					"StringEquals": policy.ConditionKeyValue{
						"aws:SourceAccount": policy.ConditionValue{"123456789012"},
					},
				},
			},
		},
	}
	if isPolicyPublic(doc) {
		t.Error("expected policy with Principal:* + Condition to NOT be public")
	}
}

func TestIsPolicyPublic_DenyStatement(t *testing.T) {
	doc := &policy.Document{
		Version: "2012-10-17",
		Statement: []policy.Statement{
			{
				Effect:    policy.EffectDeny,
				Principal: &policy.Principal{Everyone: true},
				Action:    policy.ActionList{"secretsmanager:GetSecretValue"},
			},
		},
	}
	if isPolicyPublic(doc) {
		t.Error("expected Deny statement to NOT be flagged as public")
	}
}

func TestIsPolicyPublic_NilDocument(t *testing.T) {
	if isPolicyPublic(nil) {
		t.Error("expected nil document to NOT be public")
	}
}

func TestValidatePolicyDocument_Valid(t *testing.T) {
	valid := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:root"},"Action":"secretsmanager:GetSecretValue","Resource":"*"}]}`
	checks := validatePolicyDocument(valid)
	if len(checks) != 0 {
		t.Errorf("expected 0 checks for valid policy, got %d: %+v", len(checks), checks)
	}
}

func TestValidatePolicyDocument_Malformed(t *testing.T) {
	checks := validatePolicyDocument("{broken")
	if len(checks) != 1 {
		t.Fatalf("expected 1 check for malformed policy, got %d", len(checks))
	}
	if checks[0].CheckName != checkNameSyntax {
		t.Errorf("expected CheckName %s, got %s", checkNameSyntax, checks[0].CheckName)
	}
}

func TestValidatePolicyDocument_MissingVersion(t *testing.T) {
	noVersion := `{"Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:root"},"Action":"secretsmanager:GetSecretValue","Resource":"*"}]}`
	checks := validatePolicyDocument(noVersion)
	found := false
	for _, c := range checks {
		if c.CheckName == checkNameMissingVersion {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected RESOURCE_POLICY_MISSING_VERSION check")
	}
}

func TestValidatePolicyDocument_EmptyStatements(t *testing.T) {
	// Empty Statement array should be syntactically valid (no syntax
	// error, no missing version when Version is present).
	emptyStmts := `{"Version":"2012-10-17","Statement":[]}`
	checks := validatePolicyDocument(emptyStmts)
	if len(checks) != 0 {
		t.Errorf("expected 0 checks for empty-statement policy with Version, got %d: %+v", len(checks), checks)
	}
}

func TestEnsurePolicyJSONValid_Valid(t *testing.T) {
	valid := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:root"},"Action":"secretsmanager:GetSecretValue","Resource":"*"}]}`
	doc, err := ensurePolicyJSONValid(valid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
}

func TestEnsurePolicyJSONValid_Invalid(t *testing.T) {
	_, err := ensurePolicyJSONValid("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestUpsertReplicationStatus_NewRegion(t *testing.T) {
	var statuses []secretsmanagerstore.ReplicationStatus
	rs := secretsmanagerstore.ReplicationStatus{Region: "us-west-2", Status: "InSync"}
	statuses = upsertReplicationStatus(statuses, rs)
	if len(statuses) != 1 {
		t.Errorf("expected 1 status, got %d", len(statuses))
	}
}

func TestUpsertReplicationStatus_ExistingRegion(t *testing.T) {
	statuses := []secretsmanagerstore.ReplicationStatus{
		{Region: "us-west-2", Status: "InSync"},
	}
	rs := secretsmanagerstore.ReplicationStatus{Region: "us-west-2", Status: "Failed"}
	statuses = upsertReplicationStatus(statuses, rs)
	if len(statuses) != 1 {
		t.Errorf("expected 1 status after upsert, got %d", len(statuses))
	}
	if statuses[0].Status != "Failed" {
		t.Errorf("expected status Failed, got %s", statuses[0].Status)
	}
}

func TestUpsertReplicationStatus_MultipleRegions(t *testing.T) {
	statuses := []secretsmanagerstore.ReplicationStatus{
		{Region: "us-west-2", Status: "InSync"},
		{Region: "eu-west-1", Status: "InSync"},
	}
	rs := secretsmanagerstore.ReplicationStatus{Region: "eu-west-1", Status: "Failed"}
	statuses = upsertReplicationStatus(statuses, rs)
	if len(statuses) != 2 {
		t.Errorf("expected 2 statuses after upsert, got %d", len(statuses))
	}
	for _, s := range statuses {
		if s.Region == "eu-west-1" && s.Status != "Failed" {
			t.Errorf("expected eu-west-1 status Failed, got %s", s.Status)
		}
	}
}
