package sns

import (
	"strings"
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

// ---------------------------------------------------------------------------
// Set-time grammar and shared-shape pins. The validator accepts exactly the
// shapes the matcher evaluates: the documented grammar with every operator
// and operand form, the numeric matching that normalises ("301.5" equals
// 3.015e2), String.Array elements matching individually, the nested
// payload-scope rules of the constraints page, and the FilterPolicy/
// FilterPolicyScope coupling — while the matcher stays fail-closed on
// operand shapes it cannot evaluate, so an invalid policy never behaves as
// match-all at delivery time.
// ---------------------------------------------------------------------------

// TestFilterPolicyGrammar pins the validator: every documented operator
// and operand form is accepted, and every shape the matcher cannot
// evaluate — or the constraints page forbids — is rejected.
func TestFilterPolicyGrammar(t *testing.T) {
	accepted := map[string]string{
		"exact scalars":         `{"k": ["v", 10, true, null]}`,
		"prefix":                `{"k": [{"prefix": "bas"}]}`,
		"suffix":                `{"k": [{"suffix": "ball"}]}`,
		"wildcard":              `{"k": [{"wildcard": "*ball"}]}`,
		"wildcard triple":       `{"k": [{"wildcard": "*a*b*c"}]}`,
		"equals-ignore-case":    `{"k": [{"equals-ignore-case": "tennis"}]}`,
		"cidr subnet":           `{"k": [{"cidr": "10.0.0.0/24"}]}`,
		"cidr address":          `{"k": [{"cidr": "10.0.0.7"}]}`,
		"numeric range":         `{"n": [{"numeric": [">", 0, "<=", 150]}]}`,
		"numeric exact":         `{"n": [{"numeric": ["=", 301.5]}]}`,
		"numeric negative":      `{"n": [{"numeric": ["<", 0]}]}`,
		"exists true":           `{"k": [{"exists": true}]}`,
		"exists false":          `{"k": [{"exists": false}]}`,
		"anything-but list":     `{"k": [{"anything-but": ["rugby", "tennis", 100]}]}`,
		"anything-but value":    `{"k": [{"anything-but": "rugby"}]}`,
		"anything-but prefix":   `{"k": [{"anything-but": {"prefix": "order-"}}]}`,
		"anything-but suffix":   `{"k": [{"anything-but": {"suffix": "ball"}}]}`,
		"anything-but wildcard": `{"k": [{"anything-but": {"wildcard": "*ball"}}]}`,
		"five keys":             `{"a": [1], "b": [1], "c": [1], "d": [1], "e": [1]}`,
		"combination at 150":    `{"a": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], "b": [1, 2, 3, 4, 5], "c": [1, 2, 3]}`,
		// Ten single-wildcard patterns: field complexity (1+1+...+1) x 10
		// patterns = exactly the 100-point budget.
		"wildcard complexity at 100": `{"a": [{"wildcard": "*1"}, {"wildcard": "*2"}, {"wildcard": "*3"}, {"wildcard": "*4"}, {"wildcard": "*5"}, {"wildcard": "*6"}, {"wildcard": "*7"}, {"wildcard": "*8"}, {"wildcard": "*9"}, {"wildcard": "*0"}]}`,
	}
	for name, policy := range accepted {
		if err := validateFilterPolicy(policy, "MessageAttributes"); err != nil {
			t.Errorf("%s: validateFilterPolicy(%s) = %v, want accepted", name, policy, err)
		}
	}

	rejected := map[string]string{
		"value not array":              `{"k": "v"}`,
		"six keys":                     `{"a": [1], "b": [1], "c": [1], "d": [1], "e": [1], "f": [1]}`,
		"combination past 150":         `{"a": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], "b": [1, 2, 3, 4, 5], "c": [1, 2, 3, 4]}`,
		"wildcard complexity past 100": `{"a": [{"wildcard": "*1"}, {"wildcard": "*2"}, {"wildcard": "*3"}, {"wildcard": "*4"}, {"wildcard": "*5"}, {"wildcard": "*6"}, {"wildcard": "*7"}, {"wildcard": "*8"}, {"wildcard": "*9"}, {"wildcard": "*10"}, {"wildcard": "*11"}]}`,
		"empty attribute name":         `{"": [1]}`,
		"prefix number operand":        `{"k": [{"prefix": 123}]}`,
		"suffix number operand":        `{"k": [{"suffix": 123}]}`,
		"wildcard number operand":      `{"k": [{"wildcard": 123}]}`,
		"equals-ignore-case bool":      `{"k": [{"equals-ignore-case": true}]}`,
		"cidr garbage":                 `{"k": [{"cidr": "not-a-subnet"}]}`,
		"cidr number":                  `{"k": [{"cidr": 10}]}`,
		"exists non-bool":              `{"k": [{"exists": "yes"}]}`,
		"numeric odd array":            `{"k": [{"numeric": [">", 0, "<"]}]}`,
		"numeric unknown operator":     `{"k": [{"numeric": ["!=", 0]}]}`,
		"numeric bound string":         `{"k": [{"numeric": [">", "0"]}]}`,
		"numeric out of range":         `{"k": [{"numeric": ["=", 2000000000]}]}`,
		"numeric accuracy":             `{"k": [{"numeric": ["=", 1.234567]}]}`,
		"exact number out of range":    `{"k": [2000000000]}`,
		"unknown operator":             `{"k": [{"regex": "x"}]}`,
		"two operators":                `{"k": [{"prefix": "a", "suffix": "b"}]}`,
		"wildcard four stars":          `{"k": [{"wildcard": "*a*b*c*d*"}]}`,
		"anything-but empty list":      `{"k": [{"anything-but": []}]}`,
		"anything-but nested numeric":  `{"k": [{"anything-but": {"numeric": [">", 1]}}]}`,
		"anything-but bool":            `{"k": [{"anything-but": true}]}`,
	}
	for name, policy := range rejected {
		if err := validateFilterPolicy(policy, "MessageAttributes"); err == nil {
			t.Errorf("%s: validateFilterPolicy(%s) accepted, want rejection", name, policy)
		}
	}
}

// TestFilterPolicySizeCap pins the documented 256 KB policy-size ceiling.
func TestFilterPolicySizeCap(t *testing.T) {
	big := `{"k": ["` + strings.Repeat("a", snsstore.MaxFilterPolicySizeBytes) + `"]}`
	if err := validateFilterPolicy(big, "MessageAttributes"); err == nil {
		t.Fatal("oversized filter policy accepted, want rejection")
	}
}

func msgAttr(typ, value string) *snsstore.MessageAttribute {
	return &snsstore.MessageAttribute{Type: typ, StringValue: value}
}

// TestFilterPolicyMatcher pins the delivery-time evaluation against the
// documented matching examples: numeric equality normalises across text
// forms (301.5 matches 3.015e2), String.Array elements match individually,
// every string operator behaves, and the null keyword matches a null body
// property.
func TestFilterPolicyMatcher(t *testing.T) {
	cases := []struct {
		name   string
		policy string
		attr   *snsstore.MessageAttribute
		want   bool
	}{
		{"exact string", `{"k": ["rugby"]}`, msgAttr("String", "rugby"), true},
		{"exact string miss", `{"k": ["rugby"]}`, msgAttr("String", "tennis"), false},
		{"numeric equality normalises exponent", `{"k": [{"numeric": ["=", 301.5]}]}`, msgAttr("Number", "3.015e2"), true},
		{"numeric equality normalises decimal", `{"k": [{"numeric": ["=", 10]}]}`, msgAttr("Number", "10.0"), true},
		{"numeric exact scalar normalises", `{"k": [10]}`, msgAttr("Number", "10.0"), true},
		{"numeric exact scalar miss", `{"k": [10]}`, msgAttr("Number", "10.1"), false},
		{"numeric range", `{"k": [{"numeric": [">", 0, "<=", 150]}]}`, msgAttr("Number", "150"), true},
		{"numeric range miss", `{"k": [{"numeric": [">", 0, "<=", 150]}]}`, msgAttr("Number", "150.1"), false},
		{"anything-but scalar", `{"k": [{"anything-but": [100, 500]}]}`, msgAttr("Number", "101"), true},
		{"anything-but scalar excluded", `{"k": [{"anything-but": [100, 500]}]}`, msgAttr("Number", "100"), false},
		{"anything-but array element free", `{"k": [{"anything-but": [100, 500]}]}`, msgAttr("Number.Array", "[100, 50]"), true},
		{"anything-but array all excluded", `{"k": [{"anything-but": [100, 500]}]}`, msgAttr("Number.Array", "[100, 500]"), false},
		{"anything-but nested prefix", `{"k": [{"anything-but": {"prefix": "order-"}}]}`, msgAttr("String", "order-1"), false},
		{"anything-but nested prefix free", `{"k": [{"anything-but": {"prefix": "order-"}}]}`, msgAttr("String", "invoice-1"), true},
		{"anything-but nested wildcard", `{"k": [{"anything-but": {"wildcard": "*ball"}}]}`, msgAttr("String", "baseball"), false},
		{"prefix", `{"k": [{"prefix": "bas"}]}`, msgAttr("String", "basketball"), true},
		{"suffix", `{"k": [{"suffix": "ball"}]}`, msgAttr("String", "football"), true},
		{"wildcard", `{"k": [{"wildcard": "*ball"}]}`, msgAttr("String", "baseball"), true},
		{"wildcard middle", `{"k": [{"wildcard": "log*check"}]}`, msgAttr("String", "log-run-check"), true},
		{"wildcard miss", `{"k": [{"wildcard": "*ball"}]}`, msgAttr("String", "hockey"), false},
		{"equals-ignore-case", `{"k": [{"equals-ignore-case": "tennis"}]}`, msgAttr("String", "TENNIS"), true},
		{"cidr inside", `{"k": [{"cidr": "10.0.0.0/24"}]}`, msgAttr("String", "10.0.0.255"), true},
		{"cidr outside", `{"k": [{"cidr": "10.0.0.0/24"}]}`, msgAttr("String", "10.1.1.0"), false},
		{"cidr single address", `{"k": [{"cidr": "10.0.0.7"}]}`, msgAttr("String", "10.0.0.7"), true},
		{"exists true", `{"k": [{"exists": true}]}`, msgAttr("String", "v"), true},
		{"exists true miss", `{"k": [{"exists": true}]}`, nil, false},
		{"exists false", `{"k": [{"exists": false}]}`, nil, true},
		{"string array element", `{"k": ["tennis"]}`, msgAttr("String.Array", `["rugby", "tennis"]`), true},
		{"string array element miss", `{"k": ["golf"]}`, msgAttr("String.Array", `["rugby", "tennis"]`), false},
		{"bool value", `{"k": [true]}`, msgAttr("String", "true"), true},
		{"missing attribute", `{"k": ["v"]}`, nil, false},
	}
	for _, tc := range cases {
		// A missing attribute is an ABSENT map key — a present key with a
		// nil value would make attrExists true.
		attrs := map[string]*snsstore.MessageAttribute{}
		if tc.attr != nil {
			attrs["k"] = tc.attr
		}
		got := matchFilterPolicy(tc.policy, "MessageAttributes", attrs, "")
		if got != tc.want {
			t.Errorf("%s: match = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestFilterPolicyMatcherBodyScope pins MessageBody-scope matching:
// top-level body properties feed the same grammar, body arrays behave as
// arrays of elements, and a null body property matches the null keyword.
func TestFilterPolicyMatcherBodyScope(t *testing.T) {
	cases := []struct {
		name   string
		policy string
		body   string
		want   bool
	}{
		{"string property", `{"k": ["rugby"]}`, `{"k": "rugby"}`, true},
		{"number property", `{"k": [{"numeric": ["=", 301.5]}]}`, `{"k": 3.015e2}`, true},
		{"bool property", `{"k": [false]}`, `{"k": false}`, true},
		{"null property", `{"k": [null]}`, `{"k": null}`, true},
		{"null property miss", `{"k": [null]}`, `{"k": "v"}`, false},
		{"array property element free", `{"k": [{"anything-but": [100, 500]}]}`, `{"k": [100, 50]}`, true},
		{"array property all excluded", `{"k": [{"anything-but": [100, 500]}]}`, `{"k": [100, 500]}`, false},
		{"missing property", `{"k": ["v"]}`, `{"other": 1}`, false},
	}
	for _, tc := range cases {
		got := matchFilterPolicy(tc.policy, "MessageBody", nil, tc.body)
		if got != tc.want {
			t.Errorf("%s: match = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestFilterPolicyMatcherFailClosed pins the fail-closed matcher: operand
// shapes the matcher cannot evaluate never match — even when a policy
// carrying them somehow reached the store past validation (for example a
// legacy record written before operand validation existed), the
// subscription filters everything rather than matching everything.
func TestFilterPolicyMatcherFailClosed(t *testing.T) {
	broken := map[string]string{
		"prefix non-string":           `{"k": [{"prefix": 123}]}`,
		"anything-but nested unknown": `{"k": [{"anything-but": {"numeric": [">", 1]}}]}`,
		"exists non-bool":             `{"k": [{"exists": "yes"}]}`,
		"numeric odd array":           `{"k": [{"numeric": [">"]}]}`,
		"unknown operator":            `{"k": [{"regex": "x"}]}`,
	}
	for name, policy := range broken {
		attrs := map[string]*snsstore.MessageAttribute{"k": msgAttr("String", "anything")}
		if matchFilterPolicy(policy, "MessageAttributes", attrs, "") {
			t.Errorf("%s: matcher matched everything, want fail-closed no-match", name)
		}
	}
}

// TestFilterPolicyNestedGrammar pins the nested-policy rules of the
// constraints page: nesting is payload-based filtering alone, only leaf
// keys count towards the five-key limit, the combination multiplies in
// each leaf's nested level (the page's example: 4 x 3 x 3 x 2 = 72), and
// a nested object must name at least one property.
func TestFilterPolicyNestedGrammar(t *testing.T) {
	accepted := map[string]string{
		"constraints-page example": `{"key_a": {"key_b": {"key_c": ["value_one", "value_two", "value_three", "value_four"]}}, "key_d": {"key_e": ["value_one", "value_two", "value_three"]}}`,
		"single nested leaf":       `{"key_a": {"key_b": ["value_one"]}}`,
		"five leaves nested":       `{"a": {"b": [1]}, "c": {"d": {"e": [1]}}, "f": [1], "g": [1], "h": [1]}`,
		// (5 values x depth 2) x (5 values x depth 3) = 10 x 15 = 150.
		"combination at 150": `{"a": {"b": [1, 2, 3, 4, 5]}, "c": {"d": {"e": [1, 2, 3, 4, 5]}}}`,
	}
	for name, policy := range accepted {
		if err := validateFilterPolicy(policy, "MessageBody"); err != nil {
			t.Errorf("%s: validateFilterPolicy(%s, MessageBody) = %v, want accepted", name, policy, err)
		}
	}

	rejected := map[string]string{
		"nested under attribute scope": `{"key_a": {"key_b": ["value_one"]}}`,
		"six leaf keys":                `{"a": {"b": [1], "c": [1]}, "d": [1], "e": [1], "f": [1], "g": [1]}`,
		// (6 x 2) x (5 x 3) = 12 x 15 = 180.
		"combination past 150": `{"a": {"b": [1, 2, 3, 4, 5, 6]}, "c": {"d": {"e": [1, 2, 3, 4, 5]}}}`,
		"empty nested object":  `{"key_a": {}}`,
		"empty nested key":     `{"": {"key_b": [1]}}`,
		"nested scalar value":  `{"key_a": {"key_b": "value_one"}}`,
	}
	for name, policy := range rejected {
		scope := "MessageBody"
		if name == "nested under attribute scope" {
			scope = "MessageAttributes"
		}
		if err := validateFilterPolicy(policy, scope); err == nil {
			t.Errorf("%s: validateFilterPolicy(%s, %s) accepted, want rejection", name, policy, scope)
		}
	}
}

// TestFilterPolicyMatcherNestedBodyScope pins nested payload matching: the
// policy's nested paths walk the body's nested objects with AND semantics
// across leaves and OR within a leaf, every operator evaluates at the leaf,
// and a body path that is absent or holds a non-object where the policy
// nests withholds the message.
func TestFilterPolicyMatcherNestedBodyScope(t *testing.T) {
	cases := []struct {
		name   string
		policy string
		body   string
		want   bool
	}{
		{"nested exact", `{"key_a": {"key_b": ["value_one"]}}`, `{"key_a": {"key_b": "value_one"}}`, true},
		{"nested exact miss", `{"key_a": {"key_b": ["value_one"]}}`, `{"key_a": {"key_b": "value_two"}}`, false},
		{"nested path absent", `{"key_a": {"key_b": ["value_one"]}}`, `{"key_a": "scalar"}`, false},
		{"nested parent absent", `{"key_a": {"key_b": ["value_one"]}}`, `{"other": 1}`, false},
		{"two levels with sibling leaves", `{"a": {"b": ["x"]}, "c": ["y"]}`, `{"a": {"b": "x"}, "c": "y"}`, true},
		{"two levels sibling miss", `{"a": {"b": ["x"]}, "c": ["y"]}`, `{"a": {"b": "x"}, "c": "z"}`, false},
		{"three levels", `{"a": {"b": {"c": ["v"]}}}`, `{"a": {"b": {"c": "v"}}}`, true},
		{"nested OR within leaf", `{"a": {"b": ["x", "y"]}}`, `{"a": {"b": "y"}}`, true},
		{"nested numeric", `{"a": {"b": [{"numeric": [">", 10]}]}}`, `{"a": {"b": 50}}`, true},
		{"nested numeric miss", `{"a": {"b": [{"numeric": [">", 10]}]}}`, `{"a": {"b": 5}}`, false},
		{"nested prefix", `{"a": {"b": [{"prefix": "ord-"}]}}`, `{"a": {"b": "ord-123"}}`, true},
		{"nested array property element", `{"a": {"b": ["x"]}}`, `{"a": {"b": ["x", "y"]}}`, true},
		{"nested array property miss", `{"a": {"b": ["z"]}}`, `{"a": {"b": ["x", "y"]}}`, false},
		{"nested null property", `{"a": {"b": [null]}}`, `{"a": {"b": null}}`, true},
		{"nested exists false on absent path", `{"a": {"b": [{"exists": false}]}}`, `{"a": "scalar"}`, true},
		{"nested exists false on present path", `{"a": {"b": [{"exists": false}]}}`, `{"a": {"b": "x"}}`, false},
		{"nested exists true", `{"a": {"b": [{"exists": true}]}}`, `{"a": {"b": "x"}}`, true},
		{"flat leaf beside nested object", `{"a": ["v"], "n": {"b": ["w"]}}`, `{"a": "v", "n": {"b": "w"}}`, true},
	}
	for _, tc := range cases {
		got := matchFilterPolicy(tc.policy, "MessageBody", nil, tc.body)
		if got != tc.want {
			t.Errorf("%s: match = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestFilterPolicyScopePair pins the attribute coupling: an absent scope is
// the documented MessageAttributes default, so a nested policy is refused
// until the same request (or the stored subscription) carries
// FilterPolicyScope MessageBody.
func TestFilterPolicyScopePair(t *testing.T) {
	nested := `{"key_a": {"key_b": ["value_one"]}}`
	flat := `{"key_a": ["value_one"]}`

	if err := validateFilterPolicyScopePair(nested, ""); err == nil {
		t.Error("nested policy with absent scope accepted, want the MessageAttributes default to refuse it")
	}
	if err := validateFilterPolicyScopePair(nested, "MessageAttributes"); err == nil {
		t.Error("nested policy under MessageAttributes scope accepted, want rejection")
	}
	if err := validateFilterPolicyScopePair(nested, "MessageBody"); err != nil {
		t.Errorf("nested policy under MessageBody scope rejected: %v", err)
	}
	if err := validateFilterPolicyScopePair(flat, ""); err != nil {
		t.Errorf("flat policy with absent scope rejected: %v", err)
	}
	if err := validateFilterPolicyScopePair(flat, "MessageBody"); err != nil {
		t.Errorf("flat policy under MessageBody scope rejected: %v", err)
	}
}
