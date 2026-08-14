package wafv2

import (
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
