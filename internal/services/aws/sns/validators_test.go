package sns

import (
	"fmt"
	"strings"
	"testing"

	snsstore "vorpalstacks/internal/store/aws/sns"
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
		if err := validateFilterPolicy(p, "MessageAttributes"); err != nil {
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
		if err := validateFilterPolicy(p, "MessageAttributes"); err == nil {
			t.Errorf("invalid policy should fail: %q", p)
		}
	}
}

func TestValidateFilterPolicy_TooManyAttributes(t *testing.T) {
	var parts []string
	for i := 0; i < 6; i++ {
		parts = append(parts, fmt.Sprintf(`"attr_%d": ["v"]`, i))
	}
	policy := "{" + strings.Join(parts, ",") + "}"
	if err := validateFilterPolicy(policy, "MessageAttributes"); err == nil {
		t.Error("policy with more than five keys should fail")
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
	if err := validateSubscriptionAttribute("SomeUnknownAttr", "any value", "sqs"); err != nil {
		t.Error("unknown attributes should pass through without validation")
	}
}

func TestValidateTopicAttribute_DataProtectionPolicyReserved(t *testing.T) {
	if err := validateTopicAttribute("DataProtectionPolicy", `{"Version":"2021-06-01"}`, false); err == nil {
		t.Error("DataProtectionPolicy must be rejected via the generic attribute path")
	}
}

func TestValidateSubscriptionAttribute_AuthenticateOnUnsubscribeReserved(t *testing.T) {
	if err := validateSubscriptionAttribute("AuthenticateOnUnsubscribe", "true", "sqs"); err == nil {
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

	if err := validatePublishParams(false, false, snsstore.MaxMessageSize, "body", strings.Repeat(cjk, 99), "", "", "", nil); err != nil {
		t.Errorf("99-character CJK subject rejected: %v", err)
	}
	if err := validatePublishParams(false, false, snsstore.MaxMessageSize, "body", strings.Repeat(cjk, 100), "", "", "", nil); err == nil {
		t.Error("100-character CJK subject accepted")
	}
	if err := validatePublishParams(false, false, snsstore.MaxMessageSize, "body", strings.Repeat(cjk, 101), "", "", "", nil); err == nil {
		t.Error("101-character CJK subject accepted")
	}
}

// TestValidatePublishParamsSubjectContent pins the Subject content rule the
// Publish documentation states beside the length ceiling: "Subjects must be
// UTF-8 text with no line breaks or control characters" — every control
// character (line breaks included) is rejected while ordinary text passes.
func TestValidatePublishParamsSubjectContent(t *testing.T) {
	for _, subject := range []string{"a\nb", "a\rb", "line one\r\nline two", "a\x00b", "a\x1bb", "a\x7fb"} {
		if err := validatePublishParams(false, false, snsstore.MaxMessageSize, "body", subject, "", "", "", nil); err == nil {
			t.Errorf("subject %q with control character accepted", subject)
		}
	}
	for _, subject := range []string{"Plain subject", "Ünïcødé subject ☂", "punctuation !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"} {
		if err := validatePublishParams(false, false, snsstore.MaxMessageSize, "body", subject, "", "", "", nil); err != nil {
			t.Errorf("ordinary subject %q rejected: %v", subject, err)
		}
	}
}

// TestValidatePublishParamsCombinedSize pins that the message ceiling counts
// the combined body and attributes — "Amazon SNS validates the combined size
// of the message body and message attributes against this value" (Publish):
// a body alone under the ceiling plus attributes over it is rejected, and
// the same body without the attributes passes.
func TestValidatePublishParamsCombinedSize(t *testing.T) {
	ceiling := snsstore.MaxMessageSize
	body := strings.Repeat("b", ceiling-100)
	attrs := map[string]*snsstore.MessageAttribute{
		"big": {Type: "String", StringValue: strings.Repeat("v", 200)},
	}
	if err := validatePublishParams(false, false, ceiling, body, "", "", "", "", attrs); err == nil {
		t.Error("body plus attributes over the ceiling accepted")
	}
	if err := validatePublishParams(false, false, ceiling, body, "", "", "", "", nil); err != nil {
		t.Errorf("body alone under the ceiling rejected: %v", err)
	}
}

// TestSubscriptionAttributeValueCap pins the subscription plane's value
// ceiling: AWS documents no value bound for subscription attributes, so the
// platform caps every value — known and unknown names alike — at the topic
// plane's documented bound, keeping the two attribute planes'
// storage-exposure ceilings symmetric.
func TestSubscriptionAttributeValueCap(t *testing.T) {
	atCap := strings.Repeat("v", snsstore.MaxSubscriptionAttributeValueLength)
	overCap := strings.Repeat("v", snsstore.MaxSubscriptionAttributeValueLength+1)

	if err := validateSubscriptionAttribute("FutureAttributeName", atCap, "sqs"); err != nil {
		t.Errorf("unknown attribute at the cap rejected: %v", err)
	}
	for name, value := range map[string]string{
		"FutureAttributeName": overCap,
		"FilterPolicy":        overCap,
	} {
		if err := validateSubscriptionAttribute(name, value, "sqs"); err == nil {
			t.Errorf("%s value over the cap accepted", name)
		}
	}
}

// TestValidatePublishParamsFifoIdentifiers pins the documented
// MessageGroupId and MessageDeduplicationId constraints — up to 128
// alphanumeric characters and punctuation — enforced on every topic type
// (standard topics forward MessageGroupId under the same rules), while
// MessageGroupId itself is optional on standard topics.
func TestValidatePublishParamsFifoIdentifiers(t *testing.T) {
	// A 128-character identifier spanning the whole documented charset:
	// alphanumerics plus every punctuation range !"#…~ .
	charset := "aZ0" + `!"#$%&'()*+,-./` + ":;<=>?@" + `[\]^_` + "`" + "{|}~"
	long := strings.Repeat(charset, 9)[:128]

	if err := validatePublishParams(true, false, snsstore.MaxMessageSize, "body", "", "", long, "d1", nil); err != nil {
		t.Errorf("128-character MessageGroupId over the documented charset rejected: %v", err)
	}
	if err := validatePublishParams(true, false, snsstore.MaxMessageSize, "body", "", "", strings.Repeat("g", 129), "d1", nil); err == nil {
		t.Error("129-character MessageGroupId accepted")
	}
	if err := validatePublishParams(true, false, snsstore.MaxMessageSize, "body", "", "", "group with spaces", "d1", nil); err == nil {
		t.Error("MessageGroupId containing a space accepted — space is outside the documented charset")
	}
	if err := validatePublishParams(false, false, snsstore.MaxMessageSize, "body", "", "", long, "", nil); err != nil {
		t.Errorf("MessageGroupId on a standard topic rejected: %v", err)
	}
	if err := validatePublishParams(false, false, snsstore.MaxMessageSize, "body", "", "", strings.Repeat("g", 129), "", nil); err == nil {
		t.Error("over-length MessageGroupId accepted on a standard topic — the same validation rules apply")
	}
	if err := validatePublishParams(true, false, snsstore.MaxMessageSize, "body", "", "", "g1", strings.Repeat("d", 129), nil); err == nil {
		t.Error("129-character MessageDeduplicationId accepted")
	}
	if err := validatePublishParams(true, false, snsstore.MaxMessageSize, "body", "", "", "g1", "dedup id", nil); err == nil {
		t.Error("MessageDeduplicationId containing a space accepted — space is outside the documented charset")
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
