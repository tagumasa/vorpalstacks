package sns

import (
	"testing"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

func strAttr(v string) *snsstore.MessageAttribute {
	return &snsstore.MessageAttribute{Type: "String", StringValue: v}
}

func numAttr(v string) *snsstore.MessageAttribute {
	return &snsstore.MessageAttribute{Type: "Number", StringValue: v}
}

// matchAttr is a shorthand for MessageAttributes-scope matching with no body.
func matchAttr(policy string, attrs map[string]*snsstore.MessageAttribute) bool {
	return matchFilterPolicy(policy, "MessageAttributes", attrs, "")
}

// matchBody is a shorthand for MessageBody-scope matching.
func matchBody(policy string, body string) bool {
	return matchFilterPolicy(policy, "MessageBody", nil, body)
}

func TestMatchFilterPolicy_EmptyPolicyMatchesAll(t *testing.T) {
	attrs := map[string]*snsstore.MessageAttribute{"event": strAttr("click")}
	if !matchAttr("", attrs) {
		t.Error("empty policy should match all messages")
	}
	if !matchAttr("{}", attrs) {
		t.Error("{} policy should match all messages")
	}
}

func TestMatchFilterPolicy_ExactMatch(t *testing.T) {
	policy := `{"event": ["order_created"]}`
	matchAttrs := map[string]*snsstore.MessageAttribute{"event": strAttr("order_created")}
	nomatchAttrs := map[string]*snsstore.MessageAttribute{"event": strAttr("order_deleted")}

	if !matchAttr(policy, matchAttrs) {
		t.Error("should match when attribute equals policy value")
	}
	if matchAttr(policy, nomatchAttrs) {
		t.Error("should not match when attribute differs")
	}
}

func TestMatchFilterPolicy_OrMatch(t *testing.T) {
	policy := `{"event": ["order_created", "order_updated"]}`
	match1 := map[string]*snsstore.MessageAttribute{"event": strAttr("order_created")}
	match2 := map[string]*snsstore.MessageAttribute{"event": strAttr("order_updated")}
	nomatch := map[string]*snsstore.MessageAttribute{"event": strAttr("order_deleted")}

	if !matchAttr(policy, match1) {
		t.Error("should match first value in OR")
	}
	if !matchAttr(policy, match2) {
		t.Error("should match second value in OR")
	}
	if matchAttr(policy, nomatch) {
		t.Error("should not match when value not in list")
	}
}

func TestMatchFilterPolicy_AndAcrossKeys(t *testing.T) {
	policy := `{"event": ["order_created"], "store": ["example_corp"]}`
	match := map[string]*snsstore.MessageAttribute{
		"event": strAttr("order_created"),
		"store": strAttr("example_corp"),
	}
	missingStore := map[string]*snsstore.MessageAttribute{
		"event": strAttr("order_created"),
	}

	if !matchAttr(policy, match) {
		t.Error("should match when all keys match")
	}
	if matchAttr(policy, missingStore) {
		t.Error("should not match when one key is missing")
	}
}

func TestMatchFilterPolicy_Prefix(t *testing.T) {
	policy := `{"event": [{"prefix": "ord-"}]}`
	match := map[string]*snsstore.MessageAttribute{"event": strAttr("ord-12345")}
	nomatch := map[string]*snsstore.MessageAttribute{"event": strAttr("click")}

	if !matchAttr(policy, match) {
		t.Error("prefix should match")
	}
	if matchAttr(policy, nomatch) {
		t.Error("prefix should not match non-matching value")
	}
}

func TestMatchFilterPolicy_AnythingBut(t *testing.T) {
	policy := `{"event": [{"anything-but": ["test"]}]}`
	match := map[string]*snsstore.MessageAttribute{"event": strAttr("production")}
	nomatch := map[string]*snsstore.MessageAttribute{"event": strAttr("test")}

	if !matchAttr(policy, match) {
		t.Error("anything-but should match non-excluded value")
	}
	if matchAttr(policy, nomatch) {
		t.Error("anything-but should not match excluded value")
	}
}

func TestMatchFilterPolicy_NumericRange(t *testing.T) {
	policy := `{"price": [{"numeric": [">=", 0, "<", 100]}]}`
	matchLow := map[string]*snsstore.MessageAttribute{"price": numAttr("50")}
	matchZero := map[string]*snsstore.MessageAttribute{"price": numAttr("0")}
	nomatchHigh := map[string]*snsstore.MessageAttribute{"price": numAttr("100")}
	nomatchNeg := map[string]*snsstore.MessageAttribute{"price": numAttr("-1")}

	if !matchAttr(policy, matchLow) {
		t.Error("50 should be in range [0, 100)")
	}
	if !matchAttr(policy, matchZero) {
		t.Error("0 should be in range [0, 100)")
	}
	if matchAttr(policy, nomatchHigh) {
		t.Error("100 should not be in range [0, 100)")
	}
	if matchAttr(policy, nomatchNeg) {
		t.Error("-1 should not be in range [0, 100)")
	}
}

func TestMatchFilterPolicy_Exists(t *testing.T) {
	policyExists := `{"special": [{"exists": true}]}`
	policyNotExists := `{"special": [{"exists": false}]}`

	withAttr := map[string]*snsstore.MessageAttribute{"special": strAttr("yes")}
	withoutAttr := map[string]*snsstore.MessageAttribute{"other": strAttr("value")}

	if !matchAttr(policyExists, withAttr) {
		t.Error("exists:true should match when attribute present")
	}
	if matchAttr(policyExists, withoutAttr) {
		t.Error("exists:true should not match when attribute absent")
	}
	if matchAttr(policyNotExists, withAttr) {
		t.Error("exists:false should not match when attribute present")
	}
	if !matchAttr(policyNotExists, withoutAttr) {
		t.Error("exists:false should match when attribute absent")
	}
}

func TestMatchFilterPolicy_MissingAttributeFails(t *testing.T) {
	policy := `{"event": ["order_created"]}`
	emptyAttrs := map[string]*snsstore.MessageAttribute{}

	if matchAttr(policy, emptyAttrs) {
		t.Error("missing attribute should fail match (unless exists:false)")
	}
}

func TestMatchFilterPolicy_AnythingBut_AbsentAttribute(t *testing.T) {
	policy := `{"event": [{"anything-but": ["test"]}]}`
	emptyAttrs := map[string]*snsstore.MessageAttribute{}

	if matchAttr(policy, emptyAttrs) {
		t.Error("anything-but should NOT match when attribute is absent")
	}
}

func TestMatchFilterPolicy_InvalidJSON_FailClosed(t *testing.T) {
	attrs := map[string]*snsstore.MessageAttribute{"event": strAttr("click")}
	if matchAttr("not valid json", attrs) {
		t.Error("invalid JSON policy should fail closed (no match)")
	}
}

func TestMatchFilterPolicy_AnythingBut_Array(t *testing.T) {
	policy := `{"event": [{"anything-but": ["test", "debug"]}]}`
	match := map[string]*snsstore.MessageAttribute{"event": strAttr("production")}
	nomatch1 := map[string]*snsstore.MessageAttribute{"event": strAttr("test")}
	nomatch2 := map[string]*snsstore.MessageAttribute{"event": strAttr("debug")}

	if !matchAttr(policy, match) {
		t.Error("anything-but array should match non-excluded value")
	}
	if matchAttr(policy, nomatch1) {
		t.Error("anything-but array should not match excluded 'test'")
	}
	if matchAttr(policy, nomatch2) {
		t.Error("anything-but array should not match excluded 'debug'")
	}
}

func TestMatchFilterPolicy_MessageBodyScope(t *testing.T) {
	policy := `{"event_type": ["order_created"]}`

	if !matchBody(policy, `{"event_type": "order_created", "data": {"id": 42}}`) {
		t.Error("MessageBody scope should match when body property matches")
	}
	if matchBody(policy, `{"event_type": "order_deleted"}`) {
		t.Error("MessageBody scope should not match when body property differs")
	}
	if matchBody(policy, `{"other": "value"}`) {
		t.Error("MessageBody scope should not match when body property is absent")
	}
}

func TestMatchFilterPolicy_MessageBodyScope_Numeric(t *testing.T) {
	policy := `{"quantity": [{"numeric": [">=", 10]}]}`

	if !matchBody(policy, `{"quantity": 50}`) {
		t.Error("MessageBody numeric 50 should be >= 10")
	}
	if matchBody(policy, `{"quantity": 5}`) {
		t.Error("MessageBody numeric 5 should not be >= 10")
	}
}

func TestMatchFilterPolicy_MessageBodyScope_InvalidJSON(t *testing.T) {
	policy := `{"event": ["test"]}`
	if matchBody(policy, "not valid json") {
		t.Error("MessageBody scope with invalid JSON body should fail closed")
	}
}
