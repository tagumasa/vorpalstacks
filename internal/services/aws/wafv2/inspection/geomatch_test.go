package inspection

import (
	"testing"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

func TestEvaluateGeoMatch(t *testing.T) {
	acl := func(stmt *wafstore.Statement) *wafstore.WebACL {
		return &wafstore.WebACL{
			Name: "geo-acl", ARN: "arn:aws:wafv2:1:regional/webacl/geo-acl/1",
			DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
			Rules: []*wafstore.Rule{{
				Name: "geo", Priority: 1, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
				Statement: stmt,
			}},
		}
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	// 8.8.8.8 delegates to US; 1.0.0.1 to AU.
	usRule := acl(&wafstore.Statement{GeoMatchStatement: &wafstore.GeoMatchStatement{CountryCodes: []string{"US"}}})
	if r := eval.Evaluate(usRule, &Request{SourceIP: "8.8.8.8"}); r.Action != ActionBlock {
		t.Fatalf("US request against US rule action = %s, want Block", r.Action)
	}
	if r := eval.Evaluate(usRule, &Request{SourceIP: "1.0.0.1"}); r.Action != ActionAllow {
		t.Fatalf("AU request against US rule action = %s, want Allow", r.Action)
	}
	// Country codes compare case-insensitively; the API documents
	// upper-case but tolerates the lower-case form.
	lowerRule := acl(&wafstore.Statement{GeoMatchStatement: &wafstore.GeoMatchStatement{CountryCodes: []string{"us"}}})
	if r := eval.Evaluate(lowerRule, &Request{SourceIP: "8.8.8.8"}); r.Action != ActionBlock {
		t.Fatalf("US request against lowercase rule action = %s, want Block", r.Action)
	}
	// An address with no allocation does not match any code.
	if r := eval.Evaluate(usRule, &Request{SourceIP: "192.0.2.1"}); r.Action != ActionAllow {
		t.Fatalf("documentation-space request action = %s, want Allow", r.Action)
	}
	// ForwardedIPConfig: absent header means the rule is not applied.
	forwarded := acl(&wafstore.Statement{GeoMatchStatement: &wafstore.GeoMatchStatement{
		CountryCodes:      []string{"US"},
		ForwardedIPConfig: &wafstore.ForwardedIPConfig{HeaderName: "X-Forwarded-For", FallbackBehavior: "MATCH"},
	}})
	if r := eval.Evaluate(forwarded, &Request{SourceIP: "8.8.8.8"}); r.Action != ActionAllow {
		t.Fatalf("absent forwarded header action = %s, want Allow (rule not applied)", r.Action)
	}
	if r := eval.Evaluate(forwarded, &Request{SourceIP: "8.8.8.8", Headers: []Header{{Name: "X-Forwarded-For", Value: "1.0.0.1"}}}); r.Action != ActionAllow {
		t.Fatalf("AU forwarded address action = %s, want Allow", r.Action)
	}
	matchAll := acl(&wafstore.Statement{GeoMatchStatement: &wafstore.GeoMatchStatement{
		CountryCodes:      []string{"US"},
		ForwardedIPConfig: &wafstore.ForwardedIPConfig{HeaderName: "X-Forwarded-For", FallbackBehavior: "MATCH"},
	}})
	if r := eval.Evaluate(matchAll, &Request{Headers: []Header{{Name: "X-Forwarded-For", Value: "not-an-ip"}}}); r.Action != ActionBlock {
		t.Fatalf("MATCH fallback action = %s, want Block", r.Action)
	}
}

// A NOT over a geo statement whose forwarded header is absent must not
// invert into a match: the rule is unapplied and the default action
// stands. With the header present and a non-listed country the NOT
// does match.
func TestEvaluateGeoMatchForwardedNotAppliedUnderNot(t *testing.T) {
	notGeo := &wafstore.WebACL{
		Name: "geo-not-acl", ARN: "arn:aws:wafv2:1:regional/webacl/geo-not-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name: "block-non-us", Priority: 1, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{NotStatement: &wafstore.NotStatement{
				Statement: &wafstore.Statement{GeoMatchStatement: &wafstore.GeoMatchStatement{
					CountryCodes:      []string{"US"},
					ForwardedIPConfig: &wafstore.ForwardedIPConfig{HeaderName: "X-Forwarded-For", FallbackBehavior: "NO_MATCH"},
				}},
			}},
		}},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	if r := eval.Evaluate(notGeo, &Request{SourceIP: "8.8.8.8"}); r.Action != ActionAllow {
		t.Fatalf("absent forwarded header under NOT action = %s, want Allow (rule not applied)", r.Action)
	}
	if r := eval.Evaluate(notGeo, &Request{SourceIP: "8.8.8.8", Headers: []Header{{Name: "X-Forwarded-For", Value: "8.8.8.8"}}}); r.Action != ActionAllow {
		t.Fatalf("US forwarded address under NOT action = %s, want Allow", r.Action)
	}
	if r := eval.Evaluate(notGeo, &Request{SourceIP: "8.8.8.8", Headers: []Header{{Name: "X-Forwarded-For", Value: "1.0.0.1"}}}); r.Action != ActionBlock {
		t.Fatalf("AU forwarded address under NOT action = %s, want Block", r.Action)
	}
}

func TestEvaluateAsnMatch(t *testing.T) {
	acl := func(stmt *wafstore.Statement) *wafstore.WebACL {
		return &wafstore.WebACL{
			Name: "asn-acl", ARN: "arn:aws:wafv2:1:regional/webacl/asn-acl/1",
			DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
			Rules: []*wafstore.Rule{{
				Name: "asn", Priority: 1, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
				Statement: stmt,
			}},
		}
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	// 8.8.8.8 originates from AS15169 (Google).
	rule := acl(&wafstore.Statement{AsnMatchStatement: &wafstore.AsnMatchStatement{AsnList: []string{"15169"}}})
	if r := eval.Evaluate(rule, &Request{SourceIP: "8.8.8.8"}); r.Action != ActionBlock {
		t.Fatalf("AS15169 request action = %s, want Block", r.Action)
	}
	other := acl(&wafstore.Statement{AsnMatchStatement: &wafstore.AsnMatchStatement{AsnList: []string{"64512"}}})
	if r := eval.Evaluate(other, &Request{SourceIP: "8.8.8.8"}); r.Action != ActionAllow {
		t.Fatalf("AS15169 request against unrelated rule action = %s, want Allow", r.Action)
	}
}
