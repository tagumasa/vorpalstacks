// Package inspection implements the AWS WAF request-inspection engine:
// it evaluates a WebACL's rules against an incoming web request and
// produces the action WAF would apply. The engine is pure — store access
// is injected through resolver functions — so it can be unit tested
// without storage and shared by every protected-resource plane.
package inspection

import (
	"strings"
	"time"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// Header is a single raw request header. The raw (non-canonicalised)
// name is preserved because WAF header matching is case insensitive over
// the wire form, and repeated headers must be kept as separate entries.
type Header struct {
	Name  string
	Value string
}

// Request is the snapshot of a web request that the evaluator inspects.
// The enforcement plane fills it from the incoming HTTP request; the
// Body must already be bounded to the body inspection limit and
// BodyTruncated must report whether bytes beyond that limit existed.
type Request struct {
	Method        string
	URIPath       string
	RawQuery      string
	Headers       []Header
	Cookies       []Header
	SourceIP      string
	HTTPVersion   string
	Body          []byte
	BodyTruncated bool
	Now           time.Time
}

// headerValue returns the values of every header with the given name,
// compared case insensitively after trimming leading and trailing
// spaces from both the request header name and the given name, per the
// Developer Guide's single-header match rule.
func (r *Request) headerValues(name string) []string {
	trimmed := trimASCIISpace(name)
	var values []string
	for _, h := range r.Headers {
		if stringsEqualFold(trimASCIISpace(h.Name), trimmed) {
			values = append(values, h.Value)
		}
	}
	return values
}

// hasHeader reports whether a header with the given name exists, using
// the same trimmed case-insensitive comparison as headerValues.
func (r *Request) hasHeader(name string) bool {
	return len(r.headerValues(name)) > 0
}

// cookieValue returns the value of the first cookie carrying the exact
// given name. Cookie names compare case sensitively, as sent on the
// wire.
func (r *Request) cookieValue(name string) string {
	for _, c := range r.Cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

// requestAcceptsHTML reports whether the request's Accept header
// includes text/html — the condition the CaptchaAction and
// ChallengeAction documentation attaches the JavaScript interstitial to.
func requestAcceptsHTML(req *Request) bool {
	for _, value := range req.headerValues("Accept") {
		if strings.Contains(value, "text/html") {
			return true
		}
	}
	return false
}

// RuleLabel is the persisted form of a rule's label declaration. The
// store keeps RuleLabels as an untyped interface because the value
// round-trips through JSON storage; normaliseRuleLabels converts it.
type RuleLabel struct {
	Name string `json:"Name"`
}

// CustomResponseBody is one entry of a WebACL's CustomResponseBodies map.
// The store keeps the whole map as an untyped interface for the same
// JSON round-trip reason.
type CustomResponseBody struct {
	ContentType string `json:"ContentType"`
	Content     string `json:"Content"`
}

// stringsEqualFold is a small ASCII-fold comparator used for header and
// label comparisons; http.CanonicalHeaderKey is avoided because header
// names may use non-canonical casing on the wire.
func stringsEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// forwardedIP resolves the client IP according to a statement's
// ForwardedIPConfig. AWS WAF inspects the first IP address in the
// configured header (AWS WAF Developer Guide, forwarded IP address
// handling). A request whose configured header is absent is not
// applied the rule at all, regardless of the fallback behaviour; the
// fallback governs only a present header that holds no parseable
// address: MATCH then treats the statement as matched.
func (r *Request) forwardedIP(cfg *wafstore.ForwardedIPConfig) (ip string, matched bool) {
	if cfg == nil || cfg.HeaderName == "" {
		return "", false
	}
	if !r.hasHeader(cfg.HeaderName) {
		return "", false
	}
	for _, value := range r.headerValues(cfg.HeaderName) {
		first := firstIPAddress(value)
		if first != "" {
			return first, true
		}
	}
	return "", cfg.FallbackBehavior == "MATCH"
}

// firstIPAddress extracts the first parseable IP address from a
// comma-separated forwarded header value, through the shared
// parseForwardedAddress entry parser so bare IPv6, bracketed and
// port-suffixed forms behave identically to the IP-set path.
func firstIPAddress(value string) string {
	start := 0
	for i := 0; i <= len(value); i++ {
		if i < len(value) && value[i] != ',' {
			continue
		}
		if addr := parseForwardedAddress(value[start:i]); addr.IsValid() {
			return addr.String()
		}
		start = i + 1
	}
	return ""
}
