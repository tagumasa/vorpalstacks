package inspection

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

func tt(priority int, t string) *wafstore.TextTransformation {
	return &wafstore.TextTransformation{Priority: priority, Type: t}
}

func TestTextTransformations(t *testing.T) {
	cases := []struct {
		name string
		tts  []*wafstore.TextTransformation
		in   string
		want string
	}{
		{"NONE keeps value", []*wafstore.TextTransformation{tt(0, "NONE")}, "a b", "a b"},
		{"COMPRESS_WHITE_SPACE collapses", []*wafstore.TextTransformation{tt(0, "COMPRESS_WHITE_SPACE")}, "a \t\n b", "a b"},
		{"COMPRESS_WHITE_SPACE preserves case", []*wafstore.TextTransformation{tt(0, "COMPRESS_WHITE_SPACE")}, "BaD \t BoT", "BaD BoT"},
		{"HTML_ENTITY_DECOD numeric hex", []*wafstore.TextTransformation{tt(0, "HTML_ENTITY_DECODE")}, "&#x3c;script&#x3e;", "<script>"},
		{"HTML_ENTITY_DECOD numeric decimal", []*wafstore.TextTransformation{tt(0, "HTML_ENTITY_DECODE")}, "&#65;", "A"},
		{"HTML_ENTITY_DECODE named case insensitive", []*wafstore.TextTransformation{tt(0, "HTML_ENTITY_DECODE")}, "&QuOt;", "\""},
		{"HTML_ENTITY_DECODE long named entity", []*wafstore.TextTransformation{tt(0, "HTML_ENTITY_DECODE")}, "a&NonBreakingSpace;b", "a\u00a0b"},
		{"HTML_ENTITY_DECODE brace aliases", []*wafstore.TextTransformation{tt(0, "HTML_ENTITY_DECODE")}, "&lcub;a&rcub;", "{a}"},
		{"HTML_ENTITY_DECODE NewLine and Tab", []*wafstore.TextTransformation{tt(0, "HTML_ENTITY_DECODE")}, "a&NewLine;b&Tab;c", "a\nb\tc"},
		{"LOWERCASE", []*wafstore.TextTransformation{tt(0, "LOWERCASE")}, "ABC", "abc"},
		{"UPPERCASE", []*wafstore.TextTransformation{tt(0, "UPPERCASE")}, "abc", "ABC"},
		{"CMD_LINE deletes quotes, carets and backslashes", []*wafstore.TextTransformation{tt(0, "CMD_LINE")}, `c"at^ C:\T`, "cat c:t"},
		{"CMD_LINE replaces comma and semicolon", []*wafstore.TextTransformation{tt(0, "CMD_LINE")}, "a,b;c", "a b c"},
		{"CMD_LINE deletes spaces before slash", []*wafstore.TextTransformation{tt(0, "CMD_LINE")}, "cat /etc/passwd", "cat/etc/passwd"},
		{"CMD_LINE deletes spaces before open paren", []*wafstore.TextTransformation{tt(0, "CMD_LINE")}, "a (b)", "a(b)"},
		{"CMD_LINE lowercases after space deletion", []*wafstore.TextTransformation{tt(0, "CMD_LINE")}, "Cat /Etc", "cat/etc"},
		{"CMD_LINE_UNIX trims and lowercases", []*wafstore.TextTransformation{tt(0, "CMD_LINE_UNIX")}, "  Li\"D \t", "lid"},
		{"CMD_LINE_WIN collapses backslashes", []*wafstore.TextTransformation{tt(0, "CMD_LINE_WIN")}, "c:\\\\windows", "c:\\windows"},
		{"URL_DECODE", []*wafstore.TextTransformation{tt(0, "URL_DECODE")}, "a%20b+c", "a b c"},
		{"URL_DECODE_UNI full width", []*wafstore.TextTransformation{tt(0, "URL_DECODE_UNI")}, "%uff21", "A"},
		{"BASE64_DECODE", []*wafstore.TextTransformation{tt(0, "BASE64_DECODE")}, "QUJD", "ABC"},
		{"BASE64_DECODE_EXT ignores invalid", []*wafstore.TextTransformation{tt(0, "BASE64_DECODE_EXT")}, "Q!U-J!D", "ABC"},
		{"HEX_DECODE", []*wafstore.TextTransformation{tt(0, "HEX_DECODE")}, "414243", "ABC"},
		{"SQL_HEX_DECODE", []*wafstore.TextTransformation{tt(0, "SQL_HEX_DECODE")}, "0x414243", "ABC"},
		{"REPLACE_COMMENTS replaces whole comment", []*wafstore.TextTransformation{tt(0, "REPLACE_COMMENTS")}, "a/*x*/b", "a b"},
		{"REPLACE_COMMENTS unterminated", []*wafstore.TextTransformation{tt(0, "REPLACE_COMMENTS")}, "a/*xyz", "a "},
		{"REMOVE_COMMENTS_CHAR", []*wafstore.TextTransformation{tt(0, "REMOVE_COMMENTS_CHAR")}, "a/*b*/c--d#e", "abcde"},
		{"ESCAPE_SEQ_DECODE", []*wafstore.TextTransformation{tt(0, "ESCAPE_SEQ_DECODE")}, "a\\tb\\x41c", "a\tbAc"},
		{"CSS_DECODE unescapes", []*wafstore.TextTransformation{tt(0, "CSS_DECODE")}, "ja\\vascript", "javascript"},
		{"CSS_DECODE hex escape", []*wafstore.TextTransformation{tt(0, "CSS_DECODE")}, "\\41 bc", "Abc"},
		{"CSS_DECODE caps at four hex digits", []*wafstore.TextTransformation{tt(0, "CSS_DECODE")}, "\\123456", "\u123456"},
		{"JS_DECODE unicode", []*wafstore.TextTransformation{tt(0, "JS_DECODE")}, "\\u0041", "A"},
		{"JS_DECODE keeps unknown escapes removed", []*wafstore.TextTransformation{tt(0, "JS_DECODE")}, "\\d", "d"},
		{"JS_DECODE_EXT keeps unknown escapes", []*wafstore.TextTransformation{tt(0, "JS_DECODE_EXT")}, "\\default", "\\default"},
		{"NORMALIZE_PATH", []*wafstore.TextTransformation{tt(0, "NORMALIZE_PATH")}, "/a//b/./c/../d", "/a/b/d"},
		{"NORMALIZE_PATH_WIN", []*wafstore.TextTransformation{tt(0, "NORMALIZE_PATH_WIN")}, "\\\\a\\\\..\\\\b", "/b"},
		{"REMOVE_NULLS", []*wafstore.TextTransformation{tt(0, "REMOVE_NULLS")}, "a\x00b", "ab"},
		{"REPLACE_NULLS", []*wafstore.TextTransformation{tt(0, "REPLACE_NULLS")}, "a\x00b", "a b"},
		{"REMOVE_WHITESPACE", []*wafstore.TextTransformation{tt(0, "REMOVE_WHITESPACE")}, " a b\tc ", "abc"},
		{"TRIM", []*wafstore.TextTransformation{tt(0, "TRIM")}, "  a b  ", "a b"},
		{"TRIM_LEFT", []*wafstore.TextTransformation{tt(0, "TRIM_LEFT")}, "  a b  ", "a b  "},
		{"TRIM_RIGHT", []*wafstore.TextTransformation{tt(0, "TRIM_RIGHT")}, "  a b  ", "  a b"},
		{"priority ordering applies in sequence", []*wafstore.TextTransformation{tt(1, "LOWERCASE"), tt(0, "UPPERCASE")}, "aBc", "abc"},
		{"MD5 produces raw digest", []*wafstore.TextTransformation{tt(0, "MD5")}, "abc", "\x90\x01\x50\x98\x3c\xd2\x4f\xb0\xd6\x96\x3f\x7d\x28\xe1\x7f\x72"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(applyTextTransformations([]byte(tc.in), tc.tts))
			if got != tc.want {
				t.Fatalf("applyTextTransformations(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestByteMatchConstraints(t *testing.T) {
	cases := []struct {
		constraint string
		candidate  string
		search     string
		want       bool
	}{
		{"EXACTLY", "abc", "abc", true},
		{"EXACTLY", "abcd", "abc", false},
		{"STARTS_WITH", "abcdef", "abc", true},
		{"STARTS_WITH", "xabc", "abc", false},
		{"ENDS_WITH", "abcdef", "def", true},
		{"ENDS_WITH", "defx", "def", false},
		{"CONTAINS", "xxabcxx", "abc", true},
		{"CONTAINS", "abd", "abc", false},
		{"CONTAINS_WORD", "the attack value", "attack", true},
		{"CONTAINS_WORD", "attacked", "attack", false},
		{"CONTAINS_WORD", "re-attack", "attack", true},
		{"CONTAINS_WORD", "attack_x", "attack", false},
		{"CONTAINS_WORD", "attack", "attack", true},
	}
	for _, tc := range cases {
		if got := byteMatchConstraint([]byte(tc.candidate), []byte(tc.search), tc.constraint); got != tc.want {
			t.Errorf("byteMatchConstraint(%q, %q, %s) = %v, want %v", tc.candidate, tc.search, tc.constraint, got, tc.want)
		}
	}
}

func TestExtractFieldToMatch(t *testing.T) {
	req := &Request{
		Method:   "POST",
		URIPath:  "/admin/path",
		RawQuery: "flag=red&name=alice&name=bob",
		Headers: []Header{
			{Name: "Host", Value: "example.test"},
			{Name: "User-Agent", Value: "curl/8"},
		},
		Cookies: []Header{
			{Name: "session", Value: "abc"},
			{Name: "theme", Value: "dark"},
		},
		Body: []byte(`{"user":{"name":"eve"},"roles":["admin","dev"]}`),
	}

	uri := extractFieldToMatch(req, &wafstore.FieldToMatch{UriPath: &wafstore.All{}})
	if len(uri) != 1 || string(uri[0].value) != "/admin/path" {
		t.Fatalf("UriPath candidate = %+v", uri)
	}
	args := extractFieldToMatch(req, &wafstore.FieldToMatch{AllQueryArguments: &wafstore.All{}})
	if len(args) != 3 {
		t.Fatalf("AllQueryArguments candidates = %d, want 3", len(args))
	}
	single := extractFieldToMatch(req, &wafstore.FieldToMatch{SingleQueryArgument: &wafstore.SingleQueryArgument{Name: "name"}})
	if len(single) != 2 || string(single[0].value) != "alice" {
		t.Fatalf("SingleQueryArgument candidates = %+v", single)
	}
	ua := extractFieldToMatch(req, &wafstore.FieldToMatch{SingleHeader: &wafstore.SingleHeader{Name: "user-agent"}})
	if len(ua) != 1 || string(ua[0].value) != "curl/8" {
		t.Fatalf("SingleHeader (case insensitive) candidates = %+v", ua)
	}
	padded := extractFieldToMatch(req, &wafstore.FieldToMatch{SingleHeader: &wafstore.SingleHeader{Name: " user-agent "}})
	if len(padded) != 1 || string(padded[0].value) != "curl/8" {
		t.Fatalf("SingleHeader (trimmed) candidates = %+v", padded)
	}
	headersAll := extractFieldToMatch(req, &wafstore.FieldToMatch{Headers: &wafstore.Headers{
		MatchPattern:     wafstore.HeaderMatchPattern{All: &wafstore.All{}},
		MatchScope:       "ALL",
		OversizeHandling: "CONTINUE",
	}})
	joinedHeaders := joinCandidates(headersAll)
	if !strings.Contains(joinedHeaders, "User-Agent") || !strings.Contains(joinedHeaders, "curl/8") {
		t.Fatalf("Headers ALL candidates = %q, want names and values as separate candidates", joinedHeaders)
	}
	for _, c := range headersAll {
		if c.value[0] == 'U' && string(c.value) != "User-Agent" {
			t.Fatalf("Headers ALL candidate %q is not a bare name or value", c.value)
		}
	}
	includedPadded := extractFieldToMatch(req, &wafstore.FieldToMatch{Headers: &wafstore.Headers{
		MatchPattern:     wafstore.HeaderMatchPattern{IncludedHeaders: []string{" user-agent "}},
		MatchScope:       "VALUE",
		OversizeHandling: "CONTINUE",
	}})
	if len(includedPadded) != 1 || string(includedPadded[0].value) != "curl/8" {
		t.Fatalf("Headers IncludedHeaders (trimmed) candidates = %+v", includedPadded)
	}
	exactCookie := extractFieldToMatch(req, &wafstore.FieldToMatch{Cookies: &wafstore.Cookies{
		MatchPattern:     wafstore.CookieMatchPattern{IncludedCookies: []string{"session"}},
		MatchScope:       "VALUE",
		OversizeHandling: "CONTINUE",
	}})
	if len(exactCookie) != 1 || string(exactCookie[0].value) != "abc" {
		t.Fatalf("Cookies exact-key candidates = %+v", exactCookie)
	}
	caseCookie := extractFieldToMatch(req, &wafstore.FieldToMatch{Cookies: &wafstore.Cookies{
		MatchPattern:     wafstore.CookieMatchPattern{IncludedCookies: []string{"SESSION"}},
		MatchScope:       "VALUE",
		OversizeHandling: "CONTINUE",
	}})
	if len(caseCookie) != 0 {
		t.Fatalf("Cookies case-sensitive key candidates = %+v, want none for a case-mismatched key", caseCookie)
	}
	cookies := extractFieldToMatch(req, &wafstore.FieldToMatch{Cookies: &wafstore.Cookies{
		MatchPattern:     wafstore.CookieMatchPattern{All: &wafstore.All{}},
		MatchScope:       "VALUE",
		OversizeHandling: "CONTINUE",
	}})
	if len(cookies) != 2 || string(cookies[0].value) != "abc" {
		t.Fatalf("Cookies candidates = %+v", cookies)
	}
	jsonValues := extractFieldToMatch(req, &wafstore.FieldToMatch{JsonBody: &wafstore.JsonBody{
		MatchPattern:            wafstore.JsonMatchPattern{IncludedPaths: []string{"/user/name", "/roles/0"}},
		MatchScope:              "VALUE",
		InvalidFallbackBehavior: "NO_MATCH",
	}})
	joined := joinCandidates(jsonValues)
	if !strings.Contains(joined, "eve") || !strings.Contains(joined, "admin") {
		t.Fatalf("JsonBody candidates = %q, want the /user/name value and the /roles/0 member", joined)
	}
	// The Developer Guide's included-path example: path /a/b selects the
	// object at b, so All inspects e, f and g, Keys e and f, Values g;
	// the key b itself is part of the path and is not evaluated.
	widget := &Request{Body: []byte(`{"a":{"c":"d","b":{"e":{"f":"g"}}}}`)}
	for _, tc := range []struct {
		scope string
		want  string
	}{
		{"ALL", "efg"},
		{"KEY", "ef"},
		{"VALUE", "g"},
	} {
		candidates := extractFieldToMatch(widget, &wafstore.FieldToMatch{JsonBody: &wafstore.JsonBody{
			MatchPattern:            wafstore.JsonMatchPattern{IncludedPaths: []string{"/a/b"}},
			MatchScope:              tc.scope,
			InvalidFallbackBehavior: "NO_MATCH",
		}})
		if got := sortedJoinedCandidates(candidates); got != tc.want {
			t.Fatalf("JsonBody /a/b scope %s candidates = %q, want %q", tc.scope, got, tc.want)
		}
	}
	// With the fallback omitted, the valid prefix before the parsing
	// error is evaluated (default None behaviour).
	partial := extractFieldToMatch(&Request{Body: []byte(`{"a":"b","c":`)}, &wafstore.FieldToMatch{JsonBody: &wafstore.JsonBody{
		MatchPattern: wafstore.JsonMatchPattern{All: &wafstore.All{}},
		MatchScope:   "ALL",
	}})
	if got := sortedJoinedCandidates(partial); got != "ab" {
		t.Fatalf("JsonBody partial-parse candidates = %q, want the key a and value b", got)
	}
	order := extractFieldToMatch(req, &wafstore.FieldToMatch{HeaderOrder: &wafstore.HeaderOrderMatch{}})
	if len(order) != 1 || string(order[0].value) != "host:user-agent" {
		t.Fatalf("HeaderOrder candidate = %+v, want colon-separated lower-cased names", order)
	}
	// Invalid JSON follows EVALUATE_AS_STRING.
	bad := &Request{Body: []byte("not json")}
	evalStr := extractFieldToMatch(bad, &wafstore.FieldToMatch{JsonBody: &wafstore.JsonBody{
		MatchPattern:            wafstore.JsonMatchPattern{All: &wafstore.All{}},
		MatchScope:              "ALL",
		InvalidFallbackBehavior: "EVALUATE_AS_STRING",
	}})
	if len(evalStr) != 1 || string(evalStr[0].value) != "not json" {
		t.Fatalf("EVALUATE_AS_STRING candidates = %+v", evalStr)
	}
}

func joinCandidates(candidates []fieldCandidate) string {
	var b strings.Builder
	for _, c := range candidates {
		b.Write(c.value)
	}
	return b.String()
}

// sortedJoinedCandidates renders candidates in sorted order so map
// iteration order cannot destabilise assertions.
func sortedJoinedCandidates(candidates []fieldCandidate) string {
	values := make([]string, 0, len(candidates))
	for _, c := range candidates {
		values = append(values, string(c.value))
	}
	sort.Strings(values)
	return strings.Join(values, "")
}

// fakeRateTracker is a deterministic RateTracker for tests: every key
// counts one hit per Hit call.
type fakeRateTracker struct {
	hits map[string]int64
}

func (f *fakeRateTracker) Hit(key RateKey, window time.Duration, now time.Time) int64 {
	f.hits[key.Value]++
	return f.hits[key.Value]
}

func newTestEvaluator(rate RateTracker, ipSets map[string]*wafstore.IPSet, regexSets map[string]*wafstore.RegexPatternSet, groups map[string]*wafstore.RuleGroup) *Evaluator {
	return NewEvaluator(Resolvers{
		IPSet:     func(arn string) (*wafstore.IPSet, error) { return ipSets[arn], nil },
		RegexSet:  func(arn string) (*wafstore.RegexPatternSet, error) { return regexSets[arn], nil },
		RuleGroup: func(arn string) (*wafstore.RuleGroup, error) { return groups[arn], nil },
		Rate:      rate,
	})
}

func TestEvaluateByteMatchRuleWithDefaultAction(t *testing.T) {
	acl := &wafstore.WebACL{
		Name:          "test-acl",
		ARN:           "arn:aws:wafv2::webacl/test-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name:     "block-attack",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
				SearchString:         []byte("attack"),
				FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
				TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
				PositionalConstraint: "CONTAINS_WORD",
			}},
		}},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	blocked := eval.Evaluate(acl, &Request{URIPath: "/the/attack/here"})
	if blocked.Action != ActionBlock {
		t.Fatalf("action = %s, want Block", blocked.Action)
	}
	if len(blocked.MatchedRules) != 1 || blocked.MatchedRules[0].RuleName != "block-attack" {
		t.Fatalf("matched rules = %+v", blocked.MatchedRules)
	}
	// A rule declared directly in the web ACL carries no within-group
	// name in its sampled requests.
	if blocked.MatchedRules[0].RuleNameWithinRuleGroup != "" {
		t.Fatalf("top-level match carries a within-group name: %+v", blocked.MatchedRules[0])
	}
	if blocked.CustomResponse == nil || blocked.CustomResponse.StatusCode != 403 || blocked.CustomResponse.Body != "" {
		// No custom response is configured, so the default Block
		// response applies: status 403 with no body.
		t.Fatalf("custom response = %+v, want the default 403", blocked.CustomResponse)
	}

	allowed := eval.Evaluate(acl, &Request{URIPath: "/ innocuous"})
	if allowed.Action != ActionAllow {
		t.Fatalf("default action = %s, want Allow", allowed.Action)
	}
}

func TestEvaluateCountContinuesAndAddsLabels(t *testing.T) {
	acl := &wafstore.WebACL{
		Name:          "count-acl",
		ARN:           "arn:aws:wafv2:us-east-1:111122223333:webacl/count-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{
			{
				Name:       "count-first",
				Priority:   1,
				Action:     &wafstore.Action{Count: &wafstore.CountAction{CustomRequestHandling: &wafstore.CustomRequestHandling{InsertHeaders: []wafstore.CustomHTTPHeader{{Name: "x-counted", Value: "1"}}}}},
				RuleLabels: []interface{}{map[string]interface{}{"Name": "suspect"}},
				Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
					SearchString:         []byte("probe"),
					FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
					TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
					PositionalConstraint: "CONTAINS",
				}},
			},
			{
				Name:     "block-on-label",
				Priority: 2,
				Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
				Statement: &wafstore.Statement{LabelMatchStatement: &wafstore.LabelMatchStatement{
					Scope: "LABEL",
					Key:   "suspect",
				}},
			},
		},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)
	result := eval.Evaluate(acl, &Request{URIPath: "/probe"})

	if result.Action != ActionBlock {
		t.Fatalf("action = %s, want Block via label", result.Action)
	}
	qualified := "awswaf:111122223333:webacl:count-acl:suspect"
	if len(result.Labels) != 1 || result.Labels[0] != qualified {
		t.Fatalf("labels = %+v, want %s", result.Labels, qualified)
	}
	if len(result.MatchedRules) != 2 {
		t.Fatalf("matched rules = %+v, want both rules", result.MatchedRules)
	}
	found := false
	for _, h := range result.InsertHeaders {
		if h.Name == "x-amzn-waf-x-counted" {
			found = true
		}
	}
	if !found {
		t.Fatalf("insert headers = %+v, want the count rule's header under the x-amzn-waf- prefix", result.InsertHeaders)
	}
}

func TestEvaluateLabelMatchSemantics(t *testing.T) {
	qualified := "awswaf:111122223333:webacl:match-acl:header:encoding:utf8"
	labelACL := func(scope, key string) *wafstore.WebACL {
		return &wafstore.WebACL{
			Name:          "match-acl",
			ARN:           "arn:aws:wafv2:us-east-1:111122223333:webacl/match-acl/1",
			DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
			Rules: []*wafstore.Rule{
				{
					Name:       "adder",
					Priority:   1,
					Action:     &wafstore.Action{Count: &wafstore.CountAction{}},
					RuleLabels: []interface{}{map[string]interface{}{"Name": "header:encoding:utf8"}},
					Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
						SearchString:         []byte("probe"),
						FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
						TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
						PositionalConstraint: "CONTAINS",
					}},
				},
				{
					Name:     "consumer",
					Priority: 2,
					Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
					Statement: &wafstore.Statement{LabelMatchStatement: &wafstore.LabelMatchStatement{
						Scope: scope,
						Key:   key,
					}},
				},
			},
		}
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	cases := []struct {
		scope string
		key   string
		block bool
	}{
		{"LABEL", "header:encoding:utf8", true},
		{"LABEL", qualified, true},
		{"LABEL", "encoding:utf8", true},
		{"LABEL", "header:encoding2:utf8", false},
		{"LABEL", "HEADER:encoding:utf8", false},
		{"NAMESPACE", "header:", false},
		{"NAMESPACE", "header:encoding", true},
		{"NAMESPACE", "encoding", true},
		{"NAMESPACE", "awswaf:111122223333:webacl:match-acl:header:encoding:", true},
		{"NAMESPACE", "awswaf:111122223333:webacl:match-acl:header:encoding:utf8", false},
		{"NAMESPACE", "header:encoding2", false},
		{"NAMESPACE", "rulegroup:", false},
	}
	for _, tc := range cases {
		result := eval.Evaluate(labelACL(tc.scope, tc.key), &Request{URIPath: "/probe"})
		want := ActionAllow
		if tc.block {
			want = ActionBlock
		}
		if result.Action != want {
			t.Errorf("scope %s key %q action = %s, want %s", tc.scope, tc.key, result.Action, want)
		}
	}
}

func TestEvaluateIPSetStatement(t *testing.T) {
	ipSets := map[string]*wafstore.IPSet{
		"arn:aws:wafv2::ipset/blocked/1": {
			Addresses:        []string{"10.0.0.0/8", "192.0.2.7"},
			IPAddressVersion: "IPV4",
		},
	}
	acl := &wafstore.WebACL{
		Name:          "ip-acl",
		ARN:           "arn:aws:wafv2::webacl/ip-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name:     "block-net",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{IPSetReferenceStatement: &wafstore.IPSetReferenceStatement{
				ARN: "arn:aws:wafv2::ipset/blocked/1",
			}},
		}},
	}
	eval := newTestEvaluator(nil, ipSets, nil, nil)

	if r := eval.Evaluate(acl, &Request{SourceIP: "10.1.2.3"}); r.Action != ActionBlock {
		t.Fatalf("CIDR member action = %s, want Block", r.Action)
	}
	if r := eval.Evaluate(acl, &Request{SourceIP: "192.0.2.7"}); r.Action != ActionBlock {
		t.Fatalf("bare address member action = %s, want Block", r.Action)
	}
	if r := eval.Evaluate(acl, &Request{SourceIP: "203.0.113.9"}); r.Action != ActionAllow {
		t.Fatalf("non-member action = %s, want Allow", r.Action)
	}
}

func TestEvaluateIPSetForwardedIP(t *testing.T) {
	ipSets := map[string]*wafstore.IPSet{
		"arn:aws:wafv2::ipset/proxies/1": {Addresses: []string{"198.51.100.0/24"}},
	}
	acl := &wafstore.WebACL{
		Name:          "fwd-acl",
		ARN:           "arn:aws:wafv2::webacl/fwd-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name:     "block-proxy",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{IPSetReferenceStatement: &wafstore.IPSetReferenceStatement{
				ARN: "arn:aws:wafv2::ipset/proxies/1",
				IPSetForwardedIPConfig: &wafstore.IPSetForwardedIPConfig{
					HeaderName:       "X-Forwarded-For",
					Position:         "FIRST",
					FallbackBehavior: "NO_MATCH",
				},
			}},
		}},
	}
	eval := newTestEvaluator(nil, ipSets, nil, nil)

	req := &Request{SourceIP: "203.0.113.9", Headers: []Header{{Name: "X-Forwarded-For", Value: "198.51.100.4, 203.0.113.1"}}}
	if r := eval.Evaluate(acl, req); r.Action != ActionBlock {
		t.Fatalf("forwarded first-member action = %s, want Block", r.Action)
	}
	noHeader := &Request{SourceIP: "203.0.113.9"}
	if r := eval.Evaluate(acl, noHeader); r.Action != ActionAllow {
		t.Fatalf("missing header with NO_MATCH fallback action = %s, want Allow", r.Action)
	}

	// A MATCH fallback must not fire when the header is absent: the
	// rule is simply not applied to the request.
	matchFallback := &wafstore.WebACL{
		Name:          "fwd-match-acl",
		ARN:           "arn:aws:wafv2::webacl/fwd-match-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name:     "block-proxy",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{IPSetReferenceStatement: &wafstore.IPSetReferenceStatement{
				ARN: "arn:aws:wafv2::ipset/proxies/1",
				IPSetForwardedIPConfig: &wafstore.IPSetForwardedIPConfig{
					HeaderName:       "X-Forwarded-For",
					Position:         "ANY",
					FallbackBehavior: "MATCH",
				},
			}},
		}},
	}
	if r := eval.Evaluate(matchFallback, noHeader); r.Action != ActionAllow {
		t.Fatalf("missing header with MATCH fallback action = %s, want Allow", r.Action)
	}

	// ANY inspects only the last ten addresses when the header carries
	// more, so a member early in the list is not inspected.
	memberEarly := &Request{SourceIP: "203.0.113.9", Headers: []Header{{Name: "X-Forwarded-For",
		Value: "198.51.100.4, 203.0.113.1, 203.0.113.2, 203.0.113.3, 203.0.113.4, 203.0.113.5, 203.0.113.6, 203.0.113.7, 203.0.113.8, 203.0.113.10, 203.0.113.11, 203.0.113.12"}}}
	if r := eval.Evaluate(matchFallback, memberEarly); r.Action != ActionAllow {
		t.Fatalf("ANY position with an early member action = %s, want Allow", r.Action)
	}
	memberLate := &Request{SourceIP: "203.0.113.9", Headers: []Header{{Name: "X-Forwarded-For",
		Value: "203.0.113.1, 203.0.113.2, 203.0.113.3, 203.0.113.4, 203.0.113.5, 203.0.113.6, 203.0.113.7, 203.0.113.8, 203.0.113.9, 203.0.113.10, 203.0.113.11, 198.51.100.4"}}}
	if r := eval.Evaluate(matchFallback, memberLate); r.Action != ActionBlock {
		t.Fatalf("ANY position with a late member action = %s, want Block", r.Action)
	}
}

// A NOT over a statement whose forwarded header is absent must not
// invert into a match: the rule is unapplied and the default action
// stands. With the header present and a non-member address the NOT
// does match.
func TestEvaluateIPSetForwardedNotAppliedUnderNot(t *testing.T) {
	ipSets := map[string]*wafstore.IPSet{
		"arn:aws:wafv2::ipset/proxies/1": {Addresses: []string{"198.51.100.0/24"}},
	}
	acl := &wafstore.WebACL{
		Name:          "fwd-not-acl",
		ARN:           "arn:aws:wafv2::webacl/fwd-not-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name:     "block-non-proxy",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{NotStatement: &wafstore.NotStatement{
				Statement: &wafstore.Statement{IPSetReferenceStatement: &wafstore.IPSetReferenceStatement{
					ARN: "arn:aws:wafv2::ipset/proxies/1",
					IPSetForwardedIPConfig: &wafstore.IPSetForwardedIPConfig{
						HeaderName:       "X-Forwarded-For",
						Position:         "FIRST",
						FallbackBehavior: "NO_MATCH",
					},
				}},
			}},
		}},
	}
	eval := newTestEvaluator(nil, ipSets, nil, nil)

	noHeader := &Request{SourceIP: "203.0.113.9"}
	if r := eval.Evaluate(acl, noHeader); r.Action != ActionAllow {
		t.Fatalf("absent forwarded header under NOT action = %s, want Allow (rule not applied)", r.Action)
	}
	member := &Request{SourceIP: "203.0.113.9", Headers: []Header{{Name: "X-Forwarded-For", Value: "198.51.100.4"}}}
	if r := eval.Evaluate(acl, member); r.Action != ActionAllow {
		t.Fatalf("forwarded member under NOT action = %s, want Allow", r.Action)
	}
	nonMember := &Request{SourceIP: "203.0.113.9", Headers: []Header{{Name: "X-Forwarded-For", Value: "203.0.113.5"}}}
	if r := eval.Evaluate(acl, nonMember); r.Action != ActionBlock {
		t.Fatalf("forwarded non-member under NOT action = %s, want Block", r.Action)
	}
}

// A rate-based statement with the FORWARDED_IP aggregate key is
// unapplied when the forwarded header is absent, so a NOT of it must
// not invert into a match.
func TestEvaluateRateForwardedNotAppliedUnderNot(t *testing.T) {
	acl := &wafstore.WebACL{
		Name:          "rate-not-acl",
		ARN:           "arn:aws:wafv2::webacl/rate-not-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name:     "block-non-flooder",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{NotStatement: &wafstore.NotStatement{
				Statement: &wafstore.Statement{RateBasedStatement: &wafstore.RateBasedStatement{
					Limit:            100,
					AggregateKeyType: "FORWARDED_IP",
					ForwardedIPConfig: &wafstore.ForwardedIPConfig{
						HeaderName:       "X-Forwarded-For",
						FallbackBehavior: "NO_MATCH",
					},
				}},
			}},
		}},
	}
	tracker := &fakeRateTracker{hits: map[string]int64{}}
	eval := newTestEvaluator(tracker, nil, nil, nil)

	noHeader := &Request{SourceIP: "203.0.113.9", Now: time.Now()}
	if r := eval.Evaluate(acl, noHeader); r.Action != ActionAllow {
		t.Fatalf("absent forwarded header under NOT rate action = %s, want Allow (rule not applied)", r.Action)
	}
	present := &Request{SourceIP: "203.0.113.9", Headers: []Header{{Name: "X-Forwarded-For", Value: "198.51.100.4"}}, Now: time.Now()}
	if r := eval.Evaluate(acl, present); r.Action != ActionBlock {
		t.Fatalf("present forwarded header under NOT rate action = %s, want Block", r.Action)
	}
}

func TestForwardedAddressExtraction(t *testing.T) {
	cases := []struct {
		entry string
		want  string
	}{
		{"2001:db8::1", "2001:db8::1"},
		{"::1", "::1"},
		{" 2001:db8::1 ", "2001:db8::1"},
		{"[2001:db8::1]", "2001:db8::1"},
		{"[2001:db8::1]:443", "2001:db8::1"},
		{"1.2.3.4", "1.2.3.4"},
		{"1.2.3.4:5678", "1.2.3.4"},
		{"not-an-ip", ""},
	}
	for _, tc := range cases {
		if got := firstIPAddress(tc.entry); got != tc.want {
			t.Errorf("firstIPAddress(%q) = %q, want %q", tc.entry, got, tc.want)
		}
	}
	mixed := &Request{Headers: []Header{{Name: "X-Forwarded-For", Value: "2001:db8::1, 1.2.3.4"}}}
	ip, matched := mixed.forwardedIP(&wafstore.ForwardedIPConfig{HeaderName: "X-Forwarded-For"})
	if !matched || ip != "2001:db8::1" {
		t.Fatalf("forwardedIP bare IPv6 = (%q, %v), want (2001:db8::1, true)", ip, matched)
	}
	ported := &Request{Headers: []Header{{Name: "X-Forwarded-For", Value: "1.2.3.4:5678"}}}
	ip, matched = ported.forwardedIP(&wafstore.ForwardedIPConfig{HeaderName: "X-Forwarded-For"})
	if !matched || ip != "1.2.3.4" {
		t.Fatalf("forwardedIP IPv4 with port = (%q, %v), want (1.2.3.4, true)", ip, matched)
	}
}

func TestEvaluateRateBasedStatement(t *testing.T) {
	acl := &wafstore.WebACL{
		Name:          "rate-acl",
		ARN:           "arn:aws:wafv2::webacl/rate-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name:     "rate-limit",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{RateBasedStatement: &wafstore.RateBasedStatement{
				Limit:            2,
				AggregateKeyType: "IP",
			}},
		}},
	}
	tracker := &fakeRateTracker{hits: map[string]int64{}}
	eval := newTestEvaluator(tracker, nil, nil, nil)
	req := &Request{SourceIP: "203.0.113.9", Now: time.Now()}

	if r := eval.Evaluate(acl, req); r.Action != ActionAllow {
		t.Fatalf("first request action = %s, want Allow", r.Action)
	}
	if r := eval.Evaluate(acl, req); r.Action != ActionAllow {
		t.Fatalf("second request action = %s, want Allow", r.Action)
	}
	if r := eval.Evaluate(acl, req); r.Action != ActionBlock {
		t.Fatalf("third request action = %s, want Block once the limit is exceeded", r.Action)
	}
	if tracker.hits["203.0.113.9"] != 3 {
		t.Fatalf("tracker hits = %d, want 3", tracker.hits["203.0.113.9"])
	}
}

func TestEvaluateRateBasedCustomKeys(t *testing.T) {
	customACL := func(customKeys []*wafstore.RateBasedStatementCustomKey) *wafstore.WebACL {
		return &wafstore.WebACL{
			Name:          "custom-rate-acl",
			ARN:           "arn:aws:wafv2::webacl/custom-rate-acl/1",
			DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
			Rules: []*wafstore.Rule{{
				Name:     "rate-limit",
				Priority: 1,
				Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
				Statement: &wafstore.Statement{RateBasedStatement: &wafstore.RateBasedStatement{
					Limit:            2,
					AggregateKeyType: "CUSTOM_KEYS",
					CustomKeys:       customKeys,
				}},
			}},
		}
	}
	now := time.Now()

	// Aggregation combines every key into one instance: an IP that
	// alternates methods never exceeds the per-combination limit until
	// one (IP, method) pair alone passes it.
	composite := customACL([]*wafstore.RateBasedStatementCustomKey{
		{IP: &wafstore.RateLimitEmptyKey{}},
		{HTTPMethod: &wafstore.RateLimitEmptyKey{}},
	})
	tracker := &fakeRateTracker{hits: map[string]int64{}}
	eval := newTestEvaluator(tracker, nil, nil, nil)
	get := &Request{SourceIP: "203.0.113.9", Method: "GET", Now: now}
	post := &Request{SourceIP: "203.0.113.9", Method: "POST", Now: now}
	for i := 0; i < 2; i++ {
		if r := eval.Evaluate(composite, get); r.Action != ActionAllow {
			t.Fatalf("composite GET %d action = %s, want Allow", i+1, r.Action)
		}
		if r := eval.Evaluate(composite, post); r.Action != ActionAllow {
			t.Fatalf("composite POST %d action = %s, want Allow", i+1, r.Action)
		}
	}
	// Two GETs and two POSTs so far; the third GET exceeds the limit
	// for the (IP, GET) combination only.
	if r := eval.Evaluate(composite, get); r.Action != ActionBlock {
		t.Fatalf("composite third GET action = %s, want Block", r.Action)
	}

	// A request missing any key component is omitted from the
	// aggregation entirely.
	cookieACL := customACL([]*wafstore.RateBasedStatementCustomKey{
		{Cookie: &wafstore.RateLimitCookieKey{Name: "session"}},
		{HTTPMethod: &wafstore.RateLimitEmptyKey{}},
	})
	cookieTracker := &fakeRateTracker{hits: map[string]int64{}}
	cookieEval := newTestEvaluator(cookieTracker, nil, nil, nil)
	for i := 0; i < 5; i++ {
		if r := cookieEval.Evaluate(cookieACL, &Request{Method: "GET", Now: now}); r.Action != ActionAllow {
			t.Fatalf("missing-component request %d action = %s, want Allow", i, r.Action)
		}
	}
	if len(cookieTracker.hits) != 0 {
		t.Fatalf("missing-component requests were aggregated: %+v", cookieTracker.hits)
	}
	withCookie := &Request{Method: "GET", Cookies: []Header{{Name: "session", Value: "s1"}}, Now: now}
	for _, want := range []string{ActionAllow, ActionAllow, ActionBlock} {
		if r := cookieEval.Evaluate(cookieACL, withCookie); r.Action != want {
			t.Fatalf("cookie-aggregated request action = %s, want %s", r.Action, want)
		}
	}

	// The API wire form of CustomKeys round-trips into the typed union
	// the evaluator consumes (an AWS-format rule saved through the
	// service layer lands in exactly this shape).
	wire := []byte(`{"Name":"rate-limit","Priority":1,
		"Action":{"Block":{}},
		"Statement":{"RateBasedStatement":{"Limit":2,"AggregateKeyType":"CUSTOM_KEYS",
		"CustomKeys":[{"Header":{"Name":"X-Tenant","TextTransformations":[{"Priority":0,"Type":"NONE"}]}}]}}}`)
	var wireRule wafstore.Rule
	if err := json.Unmarshal(wire, &wireRule); err != nil {
		t.Fatalf("unmarshal wire rule: %v", err)
	}
	wireACL := &wafstore.WebACL{
		Name:          "wire-rate-acl",
		ARN:           "arn:aws:wafv2::webacl/wire-rate-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules:         []*wafstore.Rule{&wireRule},
	}
	wireEval := newTestEvaluator(&fakeRateTracker{hits: map[string]int64{}}, nil, nil, nil)
	tenant := &Request{Headers: []Header{{Name: "X-Tenant", Value: "tenant-a"}}, Now: now}
	for _, want := range []string{ActionAllow, ActionAllow, ActionBlock} {
		if r := wireEval.Evaluate(wireACL, tenant); r.Action != want {
			t.Fatalf("wire-form header-key request action = %s, want %s", r.Action, want)
		}
	}
}

func TestEvaluateRuleGroupReferenceWithOverride(t *testing.T) {
	group := &wafstore.RuleGroup{
		Name: "shared-group",
		ARN:  "arn:aws:wafv2::rulegroup/shared-group/1",
		Rules: []*wafstore.Rule{{
			Name:     "inner-block",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
				SearchString:         []byte("bad"),
				FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
				TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
				PositionalConstraint: "CONTAINS",
			}},
		}},
	}
	groups := map[string]*wafstore.RuleGroup{group.ARN: group}

	baseACL := func(override *wafstore.Action) *wafstore.WebACL {
		return &wafstore.WebACL{
			Name:          "rg-acl",
			ARN:           "arn:aws:wafv2::webacl/rg-acl/1",
			DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
			Rules: []*wafstore.Rule{{
				Name:           "use-group",
				Priority:       1,
				OverrideAction: override,
				Statement:      &wafstore.Statement{RuleGroupReferenceStatement: &wafstore.RuleGroupReferenceStatement{ARN: group.ARN}},
			}},
		}
	}
	eval := newTestEvaluator(nil, nil, nil, groups)
	req := &Request{URIPath: "/bad-path"}

	none := eval.Evaluate(baseACL(&wafstore.Action{Block: nil, Count: nil, Allow: nil, Captcha: nil, Challenge: nil, Monetize: nil}), req)
	if none.Action != ActionBlock {
		t.Fatalf("OverrideAction None action = %s, want Block from inner rule", none.Action)
	}
	counted := eval.Evaluate(baseACL(&wafstore.Action{Count: &wafstore.CountAction{}}), req)
	if counted.Action != ActionAllow {
		t.Fatalf("OverrideAction Count action = %s, want Allow (counted, not terminated)", counted.Action)
	}
	if len(counted.MatchedRules) == 0 || counted.MatchedRules[len(counted.MatchedRules)-1].Action != ActionCount {
		t.Fatalf("matched rules = %+v, want a count record", counted.MatchedRules)
	}
}

func TestEvaluateRuleGroupReferenceExcludedAndActionOverrides(t *testing.T) {
	group := &wafstore.RuleGroup{
		Name: "shared-group",
		ARN:  "arn:aws:wafv2::rulegroup/shared-group/1",
		Rules: []*wafstore.Rule{{
			Name:     "inner-block",
			Priority: 1,
			Action:   &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
				SearchString:         []byte("bad"),
				FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
				TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
				PositionalConstraint: "CONTAINS",
			}},
		}},
	}
	groups := map[string]*wafstore.RuleGroup{group.ARN: group}
	refACL := func(ref *wafstore.RuleGroupReferenceStatement) *wafstore.WebACL {
		return &wafstore.WebACL{
			Name:          "rg-acl",
			ARN:           "arn:aws:wafv2::webacl/rg-acl/1",
			DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
			Rules: []*wafstore.Rule{{
				Name:      "use-group",
				Priority:  1,
				Statement: &wafstore.Statement{RuleGroupReferenceStatement: ref},
			}},
		}
	}
	eval := newTestEvaluator(nil, nil, nil, groups)
	req := &Request{URIPath: "/bad-path"}

	// ExcludedRules converts the inner terminating action to Count; the
	// match still records the within-group rule name, and an exclusion is
	// not an action override.
	excluded := eval.Evaluate(refACL(&wafstore.RuleGroupReferenceStatement{
		ARN:           group.ARN,
		ExcludedRules: []wafstore.ExcludedRule{{Name: "inner-block"}},
	}), req)
	if excluded.Action != ActionAllow || len(excluded.MatchedRules) != 1 || excluded.MatchedRules[0].Action != ActionCount {
		t.Fatalf("excluded evaluation = (%s, %+v), want Allow with one Count record", excluded.Action, excluded.MatchedRules)
	}
	if entry := excluded.MatchedRules[0]; entry.RuleNameWithinRuleGroup != "shared-group#inner-block" ||
		entry.RuleName != "inner-block" || entry.OverriddenAction != "" {
		t.Fatalf("excluded match entry = %+v", entry)
	}

	// A rule action override replaces the inner action and reports the
	// configured action it replaced.
	overridden := eval.Evaluate(refACL(&wafstore.RuleGroupReferenceStatement{
		ARN: group.ARN,
		RuleActionOverrides: []wafstore.RuleActionOverride{{
			Name:        "inner-block",
			ActionToUse: &wafstore.Action{Count: &wafstore.CountAction{}},
		}},
	}), req)
	if overridden.Action != ActionAllow || len(overridden.MatchedRules) != 1 || overridden.MatchedRules[0].Action != ActionCount {
		t.Fatalf("overridden evaluation = (%s, %+v), want Allow with the override applied", overridden.Action, overridden.MatchedRules)
	}
	if entry := overridden.MatchedRules[0]; entry.OverriddenAction != ActionBlock {
		t.Fatalf("overridden match entry = %+v, want the configured Block as the overridden action", entry)
	}

	// An override naming no rule of the group leaves evaluation to the
	// configured actions; the write path rejects such names for
	// customer-owned groups before the reference is ever stored.
	unknown := eval.Evaluate(refACL(&wafstore.RuleGroupReferenceStatement{
		ARN: group.ARN,
		RuleActionOverrides: []wafstore.RuleActionOverride{{
			Name:        "no-such-rule",
			ActionToUse: &wafstore.Action{Count: &wafstore.CountAction{}},
		}},
	}), req)
	if unknown.Action != ActionBlock {
		t.Fatalf("unknown-name override action = %s, want Block from the inner rule", unknown.Action)
	}
}

func TestEvaluatePriorityOrderAndCustomResponse(t *testing.T) {
	acl := &wafstore.WebACL{
		Name: "prio-acl",
		ARN:  "arn:aws:wafv2::webacl/prio-acl/1",
		CustomResponseBodies: map[string]interface{}{
			"denied": map[string]interface{}{"ContentType": "text/plain", "Content": "denied by policy"},
		},
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{
			{
				Name:     "allow-known",
				Priority: 1,
				Action:   &wafstore.Action{Allow: &wafstore.AllowAction{}},
				Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
					SearchString:         []byte("/safe"),
					FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
					TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
					PositionalConstraint: "STARTS_WITH",
				}},
			},
			{
				Name:     "block-all-admin",
				Priority: 2,
				Action: &wafstore.Action{Block: &wafstore.BlockAction{CustomResponse: &wafstore.CustomResponse{
					ResponseCode:          418,
					CustomResponseBodyKey: "denied",
				}}},
				Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
					SearchString:         []byte("/admin"),
					FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
					TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
					PositionalConstraint: "STARTS_WITH",
				}},
			},
		},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	safe := eval.Evaluate(acl, &Request{URIPath: "/safe/admin"})
	if safe.Action != ActionAllow {
		t.Fatalf("higher-priority allow action = %s, want Allow", safe.Action)
	}
	admin := eval.Evaluate(acl, &Request{URIPath: "/admin/x"})
	if admin.Action != ActionBlock {
		t.Fatalf("admin action = %s, want Block", admin.Action)
	}
	if admin.CustomResponse == nil || admin.CustomResponse.StatusCode != 418 || admin.CustomResponse.Body != "denied by policy" {
		t.Fatalf("custom response = %+v", admin.CustomResponse)
	}
}

func TestEvaluateLogicalStatements(t *testing.T) {
	andStmt := &wafstore.Statement{AndStatement: &wafstore.AndStatement{Statements: []*wafstore.Statement{
		{ByteMatchStatement: &wafstore.ByteMatchStatement{
			SearchString: []byte("a"), FieldToMatch: &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
			TextTransformations: []*wafstore.TextTransformation{tt(0, "NONE")}, PositionalConstraint: "CONTAINS",
		}},
		{ByteMatchStatement: &wafstore.ByteMatchStatement{
			SearchString: []byte("b"), FieldToMatch: &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
			TextTransformations: []*wafstore.TextTransformation{tt(0, "NONE")}, PositionalConstraint: "CONTAINS",
		}},
	}}}
	notStmt := &wafstore.Statement{NotStatement: &wafstore.NotStatement{Statement: andStmt}}
	acl := &wafstore.WebACL{
		Name: "logic-acl", ARN: "arn:aws:wafv2::webacl/logic-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name: "not-both", Priority: 1,
			Action:    &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: notStmt,
		}},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	if r := eval.Evaluate(acl, &Request{URIPath: "/a"}); r.Action != ActionBlock {
		t.Fatalf("single-letter path action = %s, want Block (AND fails, NOT succeeds)", r.Action)
	}
	if r := eval.Evaluate(acl, &Request{URIPath: "/ab"}); r.Action != ActionAllow {
		t.Fatalf("both-letters path action = %s, want Allow (AND holds, NOT fails)", r.Action)
	}
}

func TestEvaluateSQLiAndXSS(t *testing.T) {
	acl := &wafstore.WebACL{
		Name: "sqli-acl", ARN: "arn:aws:wafv2::webacl/sqli-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{
			{
				Name: "sqli", Priority: 1, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
				Statement: &wafstore.Statement{SqliMatchStatement: &wafstore.SqliMatchStatement{
					FieldToMatch:        &wafstore.FieldToMatch{QueryString: &wafstore.All{}},
					TextTransformations: []*wafstore.TextTransformation{tt(0, "URL_DECODE"), tt(1, "LOWERCASE")},
				}},
			},
			{
				Name: "xss", Priority: 2, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
				Statement: &wafstore.Statement{XssMatchStatement: &wafstore.XssMatchStatement{
					FieldToMatch:        &wafstore.FieldToMatch{Body: &wafstore.Body{OversizeHandling: "CONTINUE"}},
					TextTransformations: []*wafstore.TextTransformation{tt(0, "HTML_ENTITY_DECODE")},
				}},
			},
		},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	if r := eval.Evaluate(acl, &Request{RawQuery: "q=1%20UNION%20SELECT"}); r.Action != ActionBlock {
		t.Fatalf("SQLi query action = %s, want Block", r.Action)
	}
	if r := eval.Evaluate(acl, &Request{Body: []byte("&#x3c;script&#x3e;alert(1)")}); r.Action != ActionBlock {
		t.Fatalf("XSS body action = %s, want Block", r.Action)
	}
	if r := eval.Evaluate(acl, &Request{RawQuery: "q=harmless", Body: []byte("plain")}); r.Action != ActionAllow {
		t.Fatalf("benign action = %s, want Allow", r.Action)
	}
}

func TestSQLiSensitivityDefaultIsLow(t *testing.T) {
	// The extended pattern set runs only at HIGH sensitivity; the
	// API model documents LOW as the default, so an unset level must
	// not detect an extended-only pattern.
	extendedOnly := []byte("q=information_schema.tables")
	if sqliDetected(extendedOnly, "") {
		t.Fatal("unset sensitivity must behave as LOW and skip extended patterns")
	}
	if sqliDetected(extendedOnly, "LOW") {
		t.Fatal("LOW sensitivity must skip extended patterns")
	}
	if !sqliDetected(extendedOnly, "HIGH") {
		t.Fatal("HIGH sensitivity must run extended patterns")
	}
	if !sqliDetected([]byte("a=1 UNION SELECT * FROM t"), "") {
		t.Fatal("core patterns must run at every level")
	}
}

func TestMatchedRulesCarryRuleSamplingSetting(t *testing.T) {
	acl := &wafstore.WebACL{
		Name: "sample-acl", ARN: "arn:aws:wafv2::webacl/sample-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{
			{
				Name: "count-with-own-config", Priority: 1,
				Action: &wafstore.Action{Count: &wafstore.CountAction{}},
				Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
					SearchString:         []byte("/a"),
					FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
					TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
					PositionalConstraint: "STARTS_WITH",
				}},
				VisibilityConfig: &wafstore.VisibilityConfig{MetricName: "own", SampledRequestsEnabled: false},
			},
			{
				Name: "count-no-config", Priority: 2,
				Action: &wafstore.Action{Count: &wafstore.CountAction{}},
				Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
					SearchString:         []byte("/a"),
					FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
					TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
					PositionalConstraint: "STARTS_WITH",
				}},
			},
		},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)
	out := eval.Evaluate(acl, &Request{URIPath: "/a"})
	if len(out.MatchedRules) != 2 {
		t.Fatalf("matched rules = %d, want 2", len(out.MatchedRules))
	}
	first := out.MatchedRules[0]
	if first.SampledRequestsEnabled == nil || *first.SampledRequestsEnabled {
		t.Fatalf("rule with own visibility config must report its own sampling setting, got %+v", first.SampledRequestsEnabled)
	}
	second := out.MatchedRules[1]
	if second.SampledRequestsEnabled != nil {
		t.Fatalf("rule without visibility config must defer to the web ACL setting, got %+v", second.SampledRequestsEnabled)
	}
}

func TestEvaluateSizeConstraint(t *testing.T) {
	acl := &wafstore.WebACL{
		Name: "size-acl", ARN: "arn:aws:wafv2::webacl/size-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name: "too-big", Priority: 1, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{SizeConstraintStatement: &wafstore.SizeConstraintStatement{
				FieldToMatch:        &wafstore.FieldToMatch{QueryString: &wafstore.All{}},
				TextTransformations: []*wafstore.TextTransformation{tt(0, "NONE")},
				ComparisonOperator:  "GT",
				Size:                10,
			}},
		}},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)

	if r := eval.Evaluate(acl, &Request{RawQuery: "a=1"}); r.Action != ActionAllow {
		t.Fatalf("short query action = %s, want Allow", r.Action)
	}
	if r := eval.Evaluate(acl, &Request{RawQuery: "aaaaaaaaaaa=bbbbbbbbbb"}); r.Action != ActionBlock {
		t.Fatalf("long query action = %s, want Block", r.Action)
	}
}

func TestEvaluateUnsupportedStatementsAreRecorded(t *testing.T) {
	acl := &wafstore.WebACL{
		Name: "managed-acl", ARN: "arn:aws:wafv2::webacl/managed-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name: "managed", Priority: 1, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{ManagedRuleGroupStatement: &wafstore.ManagedRuleGroupStatement{
				VendorName: "AWS", Name: "ManagedRulesCommonRuleSet",
			}},
		}},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)
	result := eval.Evaluate(acl, &Request{SourceIP: "203.0.113.9"})

	if result.Action != ActionAllow {
		t.Fatalf("action = %s, want Allow with the managed rule group treated as non-matching", result.Action)
	}
	if len(result.Unsupported) != 1 || result.Unsupported[0] != "managed" {
		t.Fatalf("unsupported = %+v, want the managed rule name", result.Unsupported)
	}
}

func TestEvaluateHeaderOversizeHandling(t *testing.T) {
	headers := make([]Header, 250)
	for i := range headers {
		headers[i] = Header{Name: "X-Pad", Value: "v"}
	}
	req := &Request{Headers: headers}
	acl := &wafstore.WebACL{
		Name: "oversize-acl", ARN: "arn:aws:wafv2::webacl/oversize-acl/1",
		DefaultAction: &wafstore.Action{Allow: &wafstore.AllowAction{}},
		Rules: []*wafstore.Rule{{
			Name: "oversize-match", Priority: 1, Action: &wafstore.Action{Block: &wafstore.BlockAction{}},
			Statement: &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
				SearchString: []byte("never-present-value"),
				FieldToMatch: &wafstore.FieldToMatch{Headers: &wafstore.Headers{
					MatchPattern:     wafstore.HeaderMatchPattern{All: &wafstore.All{}},
					MatchScope:       "VALUE",
					OversizeHandling: "MATCH",
				}},
				TextTransformations:  []*wafstore.TextTransformation{tt(0, "NONE")},
				PositionalConstraint: "CONTAINS",
			}},
		}},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)
	if r := eval.Evaluate(acl, req); r.Action != ActionBlock {
		t.Fatalf("oversize MATCH action = %s, want Block", r.Action)
	}
}

func TestEvaluateDefaultBlockAction(t *testing.T) {
	acl := &wafstore.WebACL{
		Name: "deny-acl", ARN: "arn:aws:wafv2::webacl/deny-acl/1",
		DefaultAction: &wafstore.Action{Block: &wafstore.BlockAction{}},
	}
	eval := newTestEvaluator(nil, nil, nil, nil)
	if r := eval.Evaluate(acl, &Request{}); r.Action != ActionBlock {
		t.Fatalf("default block action = %s, want Block", r.Action)
	}
}
