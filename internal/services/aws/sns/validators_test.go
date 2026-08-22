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
	valid := []string{"MessageAttributes", "MessageBody"}
	invalid := []string{"", "Invalid", "messageAttributes", "messageattributes", "MessageBodyAttributes"}

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

func TestValidateTopicAttribute_DataProtectionPolicyReserved(t *testing.T) {
	if err := validateTopicAttribute("DataProtectionPolicy", `{"Version":"2021-06-01"}`); err == nil {
		t.Error("DataProtectionPolicy must be rejected via the generic attribute path")
	}
}

func TestValidateSubscriptionAttribute_AuthenticateOnUnsubscribeReserved(t *testing.T) {
	if err := validateSubscriptionAttribute("AuthenticateOnUnsubscribe", "true"); err == nil {
		t.Error("AuthenticateOnUnsubscribe must be rejected via SetSubscriptionAttributes")
	}
}

// TestValidatePublishParamsUnicodeSubject pins that the Publish Subject
// ceiling counts Unicode characters and tops out at 99: the AWS
// documentation states subjects are "UTF-8 text … less than 100 characters
// long", so 100 characters are rejected and rune-legal multibyte subjects
// must not be rejected on byte length.
func TestValidatePublishParamsUnicodeSubject(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validatePublishParams(false, false, "body", strings.Repeat(cjk, 99), "", "", ""); err != nil {
		t.Errorf("99-character CJK subject rejected: %v", err)
	}
	if err := validatePublishParams(false, false, "body", strings.Repeat(cjk, 100), "", "", ""); err == nil {
		t.Error("100-character CJK subject accepted")
	}
	if err := validatePublishParams(false, false, "body", strings.Repeat(cjk, 101), "", "", ""); err == nil {
		t.Error("101-character CJK subject accepted")
	}
}

// TestValidatePlatformApplicationNameLength pins the documented
// CreatePlatformApplication name ceiling of 256 characters (the model's
// member documentation states "between 1 and 256 characters long").
func TestValidatePlatformApplicationNameLength(t *testing.T) {
	if err := validatePlatformApplicationName(strings.Repeat("a", 200)); err != nil {
		t.Errorf("200-character platform application name rejected: %v", err)
	}
	if err := validatePlatformApplicationName(strings.Repeat("a", 256)); err != nil {
		t.Errorf("256-character platform application name rejected: %v", err)
	}
	if err := validatePlatformApplicationName(strings.Repeat("a", 257)); err == nil {
		t.Error("257-character platform application name accepted")
	}
	if err := validatePlatformApplicationName(""); err == nil {
		t.Error("empty platform application name accepted")
	}
	for _, name := range []string{"MyApp.prod_1", "a", "0-9_underscore.dot"} {
		if err := validatePlatformApplicationName(name); err != nil {
			t.Errorf("valid platform application name %q rejected: %v", name, err)
		}
	}
	for _, name := range []string{"bad name", "name!", "slash/name", "\u65e5\u672c", "sp@ce"} {
		if err := validatePlatformApplicationName(name); err == nil {
			t.Errorf("invalid platform application name %q accepted", name)
		}
	}
}
