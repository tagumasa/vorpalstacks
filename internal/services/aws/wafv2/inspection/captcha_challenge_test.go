package inspection

import (
	"testing"
	"time"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// secretTokenValidator adapts ParseToken to the evaluator's validator
// contract for tests.
type secretTokenValidator struct {
	secret []byte
}

func (v secretTokenValidator) ValidateToken(value string) (ChallengeToken, bool) {
	return ParseToken(v.secret, value)
}

func uriMatchStatement(fragment string) *wafstore.Statement {
	return &wafstore.Statement{
		ByteMatchStatement: &wafstore.ByteMatchStatement{
			FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
			SearchString:         []byte(fragment),
			PositionalConstraint: "CONTAINS",
			TextTransformations:  []*wafstore.TextTransformation{{Priority: 0, Type: "NONE"}},
		},
	}
}

func captchaRule(name string, action *wafstore.Action, captchaConfig interface{}) *wafstore.Rule {
	return &wafstore.Rule{
		Name:          name,
		Priority:      1,
		Action:        action,
		Statement:     uriMatchStatement("/login"),
		CaptchaConfig: captchaConfig,
		VisibilityConfig: &wafstore.VisibilityConfig{
			SampledRequestsEnabled:   true,
			CloudWatchMetricsEnabled: false,
			MetricName:               name,
		},
	}
}

func allowDefaultAction() interface{} {
	return map[string]interface{}{"Allow": map[string]interface{}{}}
}

func baseRequest(uri string) *Request {
	return &Request{
		Method:  "GET",
		URIPath: uri,
		Headers: []Header{{Name: "Host", Value: "protected.test"}},
		Now:     time.Unix(1700000000, 0),
	}
}

func TestCaptchaActionInterruptsWithoutToken(t *testing.T) {
	secret := []byte("captcha-secret")
	evaluator := NewEvaluator(Resolvers{Token: secretTokenValidator{secret}})
	acl := &wafstore.WebACL{
		Name:          "captcha-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{
			captchaRule("captcha-rule", &wafstore.Action{Captcha: &wafstore.CaptchaAction{}}, nil),
		},
	}
	req := baseRequest("/login")
	req.Headers = append(req.Headers, Header{Name: "Accept", Value: "text/html,application/xhtml+xml"})
	result := evaluator.Evaluate(acl, req)

	if result.Action != ActionCaptcha {
		t.Fatalf("Action = %q, want %q", result.Action, ActionCaptcha)
	}
	if result.CustomResponse == nil || result.CustomResponse.StatusCode != 405 {
		t.Fatalf("interrupt response status = %+v, want 405", result.CustomResponse)
	}
	if len(result.CustomResponse.Headers) != 1 ||
		result.CustomResponse.Headers[0].Name != "x-amzn-waf-action" ||
		result.CustomResponse.Headers[0].Value != "captcha" {
		t.Fatalf("interrupt headers = %+v", result.CustomResponse.Headers)
	}
	if !result.InterstitialRequested {
		t.Fatal("a request accepting text/html must request the interstitial")
	}
	if len(result.MatchedRules) != 1 || result.MatchedRules[0].Action != ActionCaptcha {
		t.Fatalf("matched rules = %+v", result.MatchedRules)
	}
	if inspection := result.MatchedRules[0].Captcha; inspection == nil || inspection.FailureReason != "TOKEN_MISSING" {
		t.Fatalf("captcha inspection = %+v, want TOKEN_MISSING", result.MatchedRules[0].Captcha)
	}
}

func TestCaptchaOmitsInterstitialWithoutHTMLAccept(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	acl := &wafstore.WebACL{
		Name:          "captcha-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{
			captchaRule("captcha-rule", &wafstore.Action{Captcha: &wafstore.CaptchaAction{}}, nil),
		},
	}
	req := baseRequest("/login")
	req.Headers = append(req.Headers, Header{Name: "Accept", Value: "application/json"})
	if result := evaluator.Evaluate(acl, req); result.InterstitialRequested {
		t.Fatal("a request not accepting text/html must not request the interstitial")
	}
}

func TestCaptchaActionPassesWithValidToken(t *testing.T) {
	secret := []byte("captcha-secret")
	evaluator := NewEvaluator(Resolvers{Token: secretTokenValidator{secret}})
	acl := &wafstore.WebACL{
		Name:          "captcha-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{
			captchaRule("captcha-rule", &wafstore.Action{Captcha: &wafstore.CaptchaAction{}}, nil),
		},
	}
	req := baseRequest("/login")
	solvedAt := req.Now.Unix() - 100
	token := SignToken(secret, ChallengeToken{
		IssuedAt:          solvedAt,
		CaptchaSolvedAt:   solvedAt,
		ChallengeSolvedAt: solvedAt,
		Domains:           []string{"protected.test"},
	})
	req.Cookies = []Header{{Name: TokenCookieName, Value: token}}

	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow — a valid token continues evaluation", result.Action)
	}
	if len(result.MatchedRules) != 1 || result.MatchedRules[0].Action != ActionCaptcha {
		t.Fatalf("matched rules = %+v", result.MatchedRules)
	}
	if inspection := result.MatchedRules[0].Captcha; inspection == nil || inspection.FailureReason != "" || inspection.SolveTimestamp != solvedAt {
		t.Fatalf("captcha inspection = %+v", result.MatchedRules[0].Captcha)
	}
}

func TestCaptchaActionExpiresAfterImmunity(t *testing.T) {
	secret := []byte("captcha-secret")
	evaluator := NewEvaluator(Resolvers{Token: secretTokenValidator{secret}})
	acl := &wafstore.WebACL{
		Name:          "captcha-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{
			captchaRule("captcha-rule", &wafstore.Action{Captcha: &wafstore.CaptchaAction{}}, nil),
		},
	}
	req := baseRequest("/login")
	solvedAt := req.Now.Unix() - int64(wafstore.ImmunityTimeDefault) - 1
	token := SignToken(secret, ChallengeToken{
		CaptchaSolvedAt: solvedAt,
		Domains:         []string{"protected.test"},
	})
	req.Cookies = []Header{{Name: TokenCookieName, Value: token}}

	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionCaptcha {
		t.Fatalf("Action = %q, want Captcha after immunity expiry", result.Action)
	}
	if inspection := result.MatchedRules[0].Captcha; inspection == nil || inspection.FailureReason != "TOKEN_EXPIRED" {
		t.Fatalf("captcha inspection = %+v, want TOKEN_EXPIRED", result.MatchedRules[0].Captcha)
	}
}

func TestCaptchaDomainMismatch(t *testing.T) {
	secret := []byte("captcha-secret")
	evaluator := NewEvaluator(Resolvers{Token: secretTokenValidator{secret}})
	acl := &wafstore.WebACL{
		Name:          "captcha-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{
			captchaRule("captcha-rule", &wafstore.Action{Captcha: &wafstore.CaptchaAction{}}, nil),
		},
	}
	req := baseRequest("/login")
	token := SignToken(secret, ChallengeToken{
		CaptchaSolvedAt: req.Now.Unix(),
		Domains:         []string{"other.test"},
	})
	req.Cookies = []Header{{Name: TokenCookieName, Value: token}}

	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionCaptcha {
		t.Fatalf("Action = %q, want Captcha for a token of another domain", result.Action)
	}
	if inspection := result.MatchedRules[0].Captcha; inspection == nil || inspection.FailureReason != "TOKEN_DOMAIN_MISMATCH" {
		t.Fatalf("captcha inspection = %+v, want TOKEN_DOMAIN_MISMATCH", result.MatchedRules[0].Captcha)
	}
}

func TestChallengeActionInterruptAndPass(t *testing.T) {
	secret := []byte("challenge-secret")
	evaluator := NewEvaluator(Resolvers{Token: secretTokenValidator{secret}})
	acl := &wafstore.WebACL{
		Name:          "challenge-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{
			{
				Name:      "challenge-rule",
				Priority:  1,
				Action:    &wafstore.Action{Challenge: &wafstore.ChallengeAction{}},
				Statement: uriMatchStatement("/login"),
				ChallengeConfig: map[string]interface{}{
					"ImmunityTimeProperty": map[string]interface{}{"ImmunityTime": int64(60)},
				},
				VisibilityConfig: &wafstore.VisibilityConfig{MetricName: "challenge-rule"},
			},
		},
	}

	req := baseRequest("/login")
	result := evaluator.Evaluate(acl, req)
	if result.Action != ActionChallenge || result.CustomResponse.StatusCode != 202 {
		t.Fatalf("without a token the challenge must interrupt with 202, got %+v", result)
	}
	if len(result.CustomResponse.Headers) != 1 || result.CustomResponse.Headers[0].Value != "challenge" {
		t.Fatalf("interrupt headers = %+v", result.CustomResponse.Headers)
	}

	solvedAt := req.Now.Unix() - 30
	token := SignToken(secret, ChallengeToken{
		ChallengeSolvedAt: solvedAt,
		Domains:           []string{"protected.test"},
	})
	req.Cookies = []Header{{Name: TokenCookieName, Value: token}}
	result = evaluator.Evaluate(acl, req)
	if result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow with a freshly solved challenge token", result.Action)
	}
}

func TestImmunityRuleLevelOverridesACLLevel(t *testing.T) {
	secret := []byte("challenge-secret")
	solvedAt := time.Unix(1700000000, 0).Unix() - 400
	tokenValue := SignToken(secret, ChallengeToken{
		ChallengeSolvedAt: solvedAt,
		Domains:           []string{"protected.test"},
	})
	challengeRule := func(ruleImmunity int64) *wafstore.Rule {
		return &wafstore.Rule{
			Name:      "challenge-rule",
			Priority:  1,
			Action:    &wafstore.Action{Challenge: &wafstore.ChallengeAction{}},
			Statement: uriMatchStatement("/login"),
			ChallengeConfig: map[string]interface{}{
				"ImmunityTimeProperty": map[string]interface{}{"ImmunityTime": ruleImmunity},
			},
			VisibilityConfig: &wafstore.VisibilityConfig{MetricName: "challenge-rule"},
		}
	}

	// A rule immunity of 3600 keeps a 400-second-old solve valid even
	// though the web ACL caps challenges at 300 seconds.
	evaluator := NewEvaluator(Resolvers{Token: secretTokenValidator{secret}})
	acl := &wafstore.WebACL{
		Name:            "challenge-acl",
		DefaultAction:   allowDefaultAction(),
		ChallengeConfig: map[string]interface{}{"ImmunityTimeProperty": map[string]interface{}{"ImmunityTime": int64(300)}},
		Rules:           []*wafstore.Rule{challengeRule(3600)},
	}
	req := baseRequest("/login")
	req.Cookies = []Header{{Name: TokenCookieName, Value: tokenValue}}
	if result := evaluator.Evaluate(acl, req); result.Action != ActionAllow {
		t.Fatalf("Action = %q, want Allow — the rule's 3600-second immunity must override the ACL's 300", result.Action)
	}

	// The same solve has expired under a 300-second rule immunity even
	// though the web ACL would allow an hour.
	acl.Rules = []*wafstore.Rule{challengeRule(300)}
	acl.ChallengeConfig = map[string]interface{}{
		"ImmunityTimeProperty": map[string]interface{}{"ImmunityTime": int64(3600)},
	}
	if result := evaluator.Evaluate(acl, req); result.Action != ActionChallenge {
		t.Fatalf("Action = %q, want Challenge — the rule's 300-second immunity must override the ACL's 3600", result.Action)
	}
}

func TestCaptchaRuleInsideRuleGroupInterrupts(t *testing.T) {
	secret := []byte("captcha-secret")
	group := &wafstore.RuleGroup{
		Name: "captcha-group",
		ARN:  "arn:aws:wafv2:us-east-1:111122223333:regional/rulegroup/captcha-group/abc",
		Rules: []*wafstore.Rule{
			captchaRule("inner-captcha", &wafstore.Action{Captcha: &wafstore.CaptchaAction{}}, nil),
		},
	}
	evaluator := NewEvaluator(Resolvers{
		Token:     secretTokenValidator{secret},
		RuleGroup: func(arn string) (*wafstore.RuleGroup, error) { return group, nil },
	})
	acl := &wafstore.WebACL{
		Name:          "group-acl",
		DefaultAction: allowDefaultAction(),
		Rules: []*wafstore.Rule{
			{
				Name:     "group-reference",
				Priority: 1,
				Statement: &wafstore.Statement{
					RuleGroupReferenceStatement: &wafstore.RuleGroupReferenceStatement{
						ARN: group.ARN,
					},
				},
			},
		},
	}
	result := evaluator.Evaluate(acl, baseRequest("/login"))
	if result.Action != ActionCaptcha {
		t.Fatalf("Action = %q, want Captcha from the inner rule", result.Action)
	}
	if len(result.MatchedRules) != 1 || result.MatchedRules[0].RuleName != "inner-captcha" ||
		result.MatchedRules[0].RuleGroup != group.Name {
		t.Fatalf("matched rules = %+v", result.MatchedRules)
	}
}
