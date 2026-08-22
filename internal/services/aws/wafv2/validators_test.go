package wafv2

import (
	"strings"
	"testing"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

func TestValidateVisibilityConfig(t *testing.T) {
	valid := &wafstore.VisibilityConfig{MetricName: "test-metric"}
	if err := validateVisibilityConfig(valid); err != nil {
		t.Fatalf("valid config err = %v", err)
	}
	if err := validateVisibilityConfig(nil); err == nil {
		t.Fatal("nil config: expected error, got nil")
	}

	cases := []struct {
		name string
		vc   *wafstore.VisibilityConfig
	}{
		{"empty MetricName", &wafstore.VisibilityConfig{MetricName: ""}},
		{"too long MetricName", &wafstore.VisibilityConfig{MetricName: string(make([]byte, 256))}},
		{"invalid characters", &wafstore.VisibilityConfig{MetricName: "bad metric!"}},
		{"space rejected", &wafstore.VisibilityConfig{MetricName: "has space"}},
	}
	for _, tc := range cases {
		if err := validateVisibilityConfig(tc.vc); err == nil {
			t.Fatalf("%s: expected error, got nil", tc.name)
		}
	}

	// Pattern-allowed punctuation must pass.
	allowed := &wafstore.VisibilityConfig{MetricName: "My#Metric:1.0-2/x"}
	if err := validateVisibilityConfig(allowed); err != nil {
		t.Fatalf("allowed punctuation err = %v", err)
	}
}

func TestValidResourceTypes(t *testing.T) {
	for _, rt := range []string{
		"APPLICATION_LOAD_BALANCER", "API_GATEWAY", "APPSYNC",
		"COGNITO_USER_POOL", "APP_RUNNER_SERVICE", "VERIFIED_ACCESS_INSTANCE",
		"AMPLIFY", "AGENTCORE_GATEWAY",
	} {
		if !validResourceTypes[rt] {
			t.Errorf("expected %s to be valid", rt)
		}
	}
	for _, rt := range []string{"", "EC2_INSTANCE", "application_load_balancer", "S3"} {
		if validResourceTypes[rt] {
			t.Errorf("expected %s to be invalid", rt)
		}
	}
}

func TestCalculateRulesCapacity(t *testing.T) {
	// ByteMatchStatement costs 1 WCU; AndStatement costs 1 plus the
	// sum of nested statements.
	rules := []*wafstore.Rule{
		{
			Name: "byte-match",
			Statement: &wafstore.Statement{
				ByteMatchStatement: &wafstore.ByteMatchStatement{SearchString: []byte("x")},
			},
		},
		{
			Name: "and",
			Statement: &wafstore.Statement{
				AndStatement: &wafstore.AndStatement{
					Statements: []*wafstore.Statement{
						{ByteMatchStatement: &wafstore.ByteMatchStatement{SearchString: []byte("y")}},
						{ByteMatchStatement: &wafstore.ByteMatchStatement{SearchString: []byte("z")}},
					},
				},
			},
		},
	}
	// 1 + (1 + 1 + 1) = 4
	if got := calculateRulesCapacity(rules); got != 4 {
		t.Fatalf("capacity = %d, want 4", got)
	}
	if got := calculateRulesCapacity(nil); got != 0 {
		t.Fatalf("nil rules capacity = %d, want 0", got)
	}
}

func TestEnsureRuleGroupNotReferenced(t *testing.T) {
	// ARN-based scan helper: nil/empty ARN is a no-op.
	if err := ensureRuleGroupNotReferenced(nil, ""); err != nil {
		t.Fatalf("empty ARN err = %v, want nil", err)
	}
}

func TestEnsureNotAssociated(t *testing.T) {
	// No association stores and an empty ARN are both no-ops.
	if err := ensureNotAssociated(nil, ""); err != nil {
		t.Fatalf("empty ARN err = %v, want nil", err)
	}
	if err := ensureNotAssociated(nil, "arn:aws:wafv2::regional/webacl/x/y"); err != nil {
		t.Fatalf("nil stores err = %v, want nil", err)
	}
}

func TestValidateDefaultAction(t *testing.T) {
	if err := validateDefaultAction(nil); err == nil {
		t.Fatal("nil action: expected error, got nil")
	}
	if err := validateDefaultAction(&wafstore.Action{Allow: &wafstore.AllowAction{}}); err != nil {
		t.Fatalf("allow action err = %v", err)
	}
	if err := validateDefaultAction(&wafstore.Action{Block: &wafstore.BlockAction{}}); err != nil {
		t.Fatalf("block action err = %v", err)
	}
	if err := validateDefaultAction(&wafstore.Action{Count: &wafstore.CountAction{}}); err == nil {
		t.Fatal("count action: expected error, got nil")
	}
}

// TestValidateEntityNamePattern pins the Smithy EntityName @pattern
// ^[\w\-]+$ alongside the length gate: spaces, punctuation outside the
// class, and multibyte characters are rejected.
func TestValidateEntityNamePattern(t *testing.T) {
	for _, name := range []string{"abc-DEF_1", "ipset", "123"} {
		if err := validateEntityName(name); err != nil {
			t.Errorf("valid name %q rejected: %v", name, err)
		}
	}
	for _, name := range []string{"bad name", "name!", "bad/name", "\u65e5\u672c\u8a9e", "a.b"} {
		if err := validateEntityName(name); err == nil {
			t.Errorf("invalid name %q accepted", name)
		}
	}
}

// TestValidateEntityDescriptionPattern pins the Smithy EntityDescription
// @pattern: the value must start and end with a class character and carry
// at least one middle character (the pattern's minimum is 3 characters).
// An empty description stays accepted: the protocol layer cannot
// distinguish an omitted optional member from an explicitly empty one.
func TestValidateEntityDescriptionPattern(t *testing.T) {
	for _, desc := range []string{"valid desc", "a:b", "x-y", "WebACL for tests 2026"} {
		if err := validateEntityDescription(desc); err != nil {
			t.Errorf("valid description %q rejected: %v", desc, err)
		}
	}
	for _, desc := range []string{"a", "ab", " leading", "trailing ", "desc!", "\u65e5\u672c"} {
		if err := validateEntityDescription(desc); err == nil {
			t.Errorf("invalid description %q accepted", desc)
		}
	}
	if err := validateEntityDescription(""); err != nil {
		t.Errorf("empty description rejected: %v", err)
	}
	// A multibyte description over 256 bytes but under 256 characters
	// passes the length gate and is rejected by the pattern instead, so
	// the reported character count must reflect runes, not bytes.
	cjkLong := strings.Repeat("\u65e5", 100)
	err := validateEntityDescription(cjkLong)
	if err == nil {
		t.Error("100-character CJK description accepted")
	} else if !strings.Contains(err.Error(), "must start and end with an allowed character") {
		t.Errorf("100-character CJK description rejected on length, not pattern: %v", err)
	}
}

// TestValidateTokenDomainPattern pins the Smithy TokenDomain @pattern
// ^[\w./-]+$ on WebACL API-key token domains.
func TestValidateTokenDomainPattern(t *testing.T) {
	valid := []interface{}{"abc.com", "store.abc.com", "a-b.example"}
	if err := validateTokenDomains(valid); err != nil {
		t.Errorf("valid token domains rejected: %v", err)
	}
	for _, dom := range []string{"a b", "abc!", "abc\\com"} {
		invalid := []interface{}{dom}
		if err := validateTokenDomains(invalid); err == nil {
			t.Errorf("invalid token domain %q accepted", dom)
		}
	}
}

// TestValidateCustomResponseBodiesKeyPattern pins that CustomResponseBodies
// map keys follow the EntityName shape (the map key targets the same
// ^[\w\-]+$ pattern).
func TestValidateCustomResponseBodiesKeyPattern(t *testing.T) {
	if err := validateCustomResponseBodies(map[string]interface{}{"ok-key_1": nil}); err != nil {
		t.Errorf("valid key rejected: %v", err)
	}
	for _, key := range []string{"bad key", "key!", "key/name"} {
		if err := validateCustomResponseBodies(map[string]interface{}{key: nil}); err == nil {
			t.Errorf("invalid key %q accepted", key)
		}
	}
}

// TestValidateStatementRegexStringUnicodeLengths pins that the
// RegexMatchStatement RegexString bound follows the Smithy
// RegexPatternString @length(1, 512) trait counted in Unicode characters:
// the shape's pattern is ".*", so multibyte regex patterns are valid input
// and must not be rejected on byte length.
func TestValidateStatementRegexStringUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	inRange := &wafstore.Statement{RegexMatchStatement: &wafstore.RegexMatchStatement{
		RegexString: strings.Repeat(cjk, 200),
	}}
	if err := validateStatement(inRange); err != nil {
		t.Errorf("200-character CJK RegexString rejected: %v", err)
	}

	overRange := &wafstore.Statement{RegexMatchStatement: &wafstore.RegexMatchStatement{
		RegexString: strings.Repeat(cjk, 513),
	}}
	if err := validateStatement(overRange); err == nil {
		t.Error("513-character CJK RegexString accepted")
	}
}
