package inspection

import (
	"time"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// Terminating action values applied to a request.
const (
	ActionAllow = "Allow"
	ActionBlock = "Block"
	ActionCount = "Count"
	// ActionCaptcha and ActionChallenge interrupt the request for a
	// token check when the request carries no valid, unexpired token;
	// with one they behave like Count. ActionMonetize interrupts with
	// HTTP 402.
	ActionCaptcha   = "Captcha"
	ActionChallenge = "Challenge"
	ActionMonetize  = "Monetize"
)

// ResolvedResponse is the fully resolved custom response for a Block
// action: the status code, the body looked up from the WebACL's
// CustomResponseBodies, and the response headers.
type ResolvedResponse struct {
	StatusCode int
	Body       string
	Headers    []wafstore.CustomHTTPHeader
}

// MatchedRule records one rule that matched during evaluation.
type MatchedRule struct {
	RuleName   string
	RuleGroup  string // non-empty when the match came from a referenced rule group
	MetricName string
	Action     string
	// RuleNameWithinRuleGroup is the sampled-request rule name of the
	// match: <rule group name>#<rule name> for customer-owned groups
	// and <vendor name>#<group name>#<rule name> for managed ones. It
	// is empty when the matching rule sits directly in the web ACL,
	// where the sampled request carries no such name.
	RuleNameWithinRuleGroup string
	// OverriddenAction is the action a rule group rule was configured
	// with when a rule action override replaced it; empty when no
	// override applied to the match.
	OverriddenAction string
	// SampledRequestsEnabled carries the rule's own visibility
	// configuration when it declares one; nil defers the sampling
	// decision to the web ACL's configuration.
	SampledRequestsEnabled *bool
	// Captcha and Challenge carry the token-inspection outcome of a
	// matching Captcha or Challenge rule, for the sampled-request
	// records. A nil value means the rule's action was neither.
	Captcha   *TokenInspection
	Challenge *TokenInspection
}

// TokenInspection reports the outcome of the aws-waf-token check a
// Captcha or Challenge rule ran: the kind's solve timestamp carried by
// the supplied token, and the reason the token failed when it did.
type TokenInspection struct {
	SolveTimestamp int64  // unix seconds; 0 when the token carried no solve
	FailureReason  string // empty when the token satisfied the check
}

// Result is the outcome of evaluating a WebACL against a request.
// Action is the terminating action AWS WAF applies to the request
// (Allow, Block, or one of the interrupting Captcha, Challenge and
// Monetize actions, each with a resolved response); MatchedRules lists
// every rule that matched, in evaluation order. Unsupported names rules
// whose statements the platform cannot evaluate; they are treated as
// non-matching. InterstitialRequested reports whether the request's
// Accept header permits the JavaScript interstitial of an interrupting
// Captcha or Challenge response — the action documentation includes the
// interstitial only for requests carrying Accept: text/html.
type Result struct {
	Action                string
	MatchedRules          []MatchedRule
	Labels                []string
	InsertHeaders         []wafstore.CustomHTTPHeader
	CustomResponse        *ResolvedResponse
	InterstitialRequested bool
	Unsupported           []string
}

// RateKey identifies one rate-based aggregation bucket.
type RateKey struct {
	WebACLARN string
	RuleName  string
	Value     string
}

// RateTracker counts requests per rate key inside a sliding evaluation
// window. The implementation is provided by the caller; Hit increments
// the counter for the key and returns the number of requests counted
// for the key within the window ending at now.
type RateTracker interface {
	Hit(key RateKey, window time.Duration, now time.Time) int64
}

// Resolvers provide the referenced-entity lookups the evaluator needs.
// They are function values so the engine itself stays free of store
// construction and can be driven by unit tests with in-memory data.
type Resolvers struct {
	IPSet     func(arn string) (*wafstore.IPSet, error)
	RegexSet  func(arn string) (*wafstore.RegexPatternSet, error)
	RuleGroup func(arn string) (*wafstore.RuleGroup, error)
	Rate      RateTracker
	// Token verifies aws-waf-token cookie values for the Captcha and
	// Challenge actions. A nil Token leaves every token check failing,
	// so those actions always interrupt.
	Token TokenValidator
}
