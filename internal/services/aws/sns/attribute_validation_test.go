package sns

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateFilterPolicy_Valid(t *testing.T) {
	validPolicies := []string{
		``,
		`{}`,
		`{"event": ["order_created"]}`,
		`{"event": ["order_created", "order_updated"], "store": ["example_corp"]}`,
		`{"event": [{"prefix": "ord-"}]}`,
		`{"event": [{"anything-but": ["test"]}]}`,
		`{"price": [{"numeric": [">=", 0, "<", 100]}]}`,
		`{"special": [{"exists": true}]}`,
	}

	for _, p := range validPolicies {
		if err := validateFilterPolicy(p); err != nil {
			t.Errorf("valid policy should pass: %q, got error: %v", p, err)
		}
	}
}

func TestValidateFilterPolicy_Invalid(t *testing.T) {
	invalidPolicies := []string{
		"not json",
		`{"event": "not an array"}`,
		`{"event": [{"unknown-op": true}]}`,
		`{"event": [{"exists": false, "prefix": "test-"}]}`,
		`{"": ["empty-key"]}`,
	}

	for _, p := range invalidPolicies {
		if err := validateFilterPolicy(p); err == nil {
			t.Errorf("invalid policy should fail: %q", p)
		}
	}
}

func TestValidateFilterPolicy_TooManyAttributes(t *testing.T) {
	var parts []string
	for i := 0; i < 101; i++ {
		parts = append(parts, fmt.Sprintf(`"attr_%d": ["v"]`, i))
	}
	policy := "{" + strings.Join(parts, ",") + "}"
	if err := validateFilterPolicy(policy); err == nil {
		t.Error("policy with >100 attributes should fail")
	}
}

func TestValidateFilterPolicyScope(t *testing.T) {
	valid := []string{"MessageAttributes", "MessageBodyAttributes"}
	invalid := []string{"", "Invalid", "messageAttributes", "messageattributes"}

	for _, v := range valid {
		if err := validateFilterPolicyScope(v); err != nil {
			t.Errorf("valid scope %q should pass: %v", v, err)
		}
	}
	for _, v := range invalid {
		if err := validateFilterPolicyScope(v); err == nil {
			t.Errorf("invalid scope %q should fail", v)
		}
	}
}

func TestValidateRedrivePolicy_Valid(t *testing.T) {
	valid := []string{
		``,
		`{"deadLetterTargetArn": "arn:aws:sqs:us-east-1:123456789012:MyDLQ"}`,
	}
	for _, v := range valid {
		if err := validateRedrivePolicy(v); err != nil {
			t.Errorf("valid redrive policy should pass: %q, got: %v", v, err)
		}
	}
}

func TestValidateRedrivePolicy_Invalid(t *testing.T) {
	invalid := []string{
		"not json",
		`{}`,
		`{"deadLetterTargetArn": ""}`,
		`{"maxReceiveCount": 3}`,
	}
	for _, v := range invalid {
		if err := validateRedrivePolicy(v); err == nil {
			t.Errorf("invalid redrive policy should fail: %q", v)
		}
	}
}

func TestValidateSubscriptionAttribute_UnknownPassthrough(t *testing.T) {
	if err := validateSubscriptionAttribute("SomeUnknownAttr", "any value"); err != nil {
		t.Error("unknown attributes should pass through without validation")
	}
}
