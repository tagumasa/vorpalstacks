// Package waf holds the cross-service WebACL request-inspection plane:
// the request snapshot types, the inspector and token-exchanger
// contracts, and the helpers that assemble an inspection request from a
// live HTTP request. Protected-resource planes (CloudFront
// distributions, API Gateway stages, AppSync APIs, Cognito user pools)
// consume these through the invoker registry; the WAFv2 service is the
// sole implementer.
package waf

import (
	"context"
	"net/http"
)

// AssociateWebACL references an existing resource hosted by the implementing
// service. WAF consults the registered checkers before creating an
// association and rejects ARNs that a checker owns but cannot resolve, with
// WAFUnavailableEntityException. ARNs whose service namespace no checker
// owns (resource types this platform does not host) keep the
// stub-association semantics.
type WebACLResourceChecker interface {
	// WebACLResourceService returns the ARN service namespace the checker
	// owns (e.g. "apigateway", "appsync", "cognito-idp").
	WebACLResourceService() string
	// WebACLResourceExists reports whether the resource referenced by the
	// ARN exists in the given region.
	WebACLResourceExists(ctx context.Context, region, resourceArn string) bool
}

// WebACLHTTPHeader is one HTTP header exchanged with the WAF inspection
// plane: either a request header to inspect or a header to insert into
// the forwarded request / custom response.
type WebACLHTTPHeader struct {
	Name  string
	Value string
}

// WebACLInspectionRequest carries the HTTP request snapshot that a
// protected-resource plane forwards to the WAF request-inspection
// plane. The body must already be bounded to the body inspection limit
// (16 KB by default) and BodyTruncated must report whether bytes beyond
// that limit existed.
type WebACLInspectionRequest struct {
	Method        string
	URIPath       string
	RawQuery      string
	SourceIP      string
	HTTPVersion   string
	Headers       []WebACLHTTPHeader
	Cookies       []WebACLHTTPHeader
	Body          []byte
	BodyTruncated bool
}

// WebACLInspectionResult is the outcome of inspecting one request. When
// Interrupts reports true (Block, or a Captcha, Challenge or Monetize
// action interrupting the request) the enforcement plane must answer
// with the resolved response (code, body, headers); InsertHeaders must
// be added to the request forwarded to the protected resource for
// non-blocking actions with custom request handling configured.
type WebACLInspectionResult struct {
	Action          string
	ResponseCode    int
	ResponseBody    string
	ResponseHeaders []WebACLHTTPHeader
	InsertHeaders   []WebACLHTTPHeader
}

// Interrupts reports whether the inspection outcome must be answered
// with the resolved response instead of being forwarded to the
// protected resource: Block, plus the Captcha and Challenge actions
// interrupting for a missing, invalid or expired token, plus Monetize's
// payment-required response. A Count outcome — and a Captcha or
// Challenge rule satisfied by a valid token, which behaves like Count —
// continues to the protected resource.
func (r WebACLInspectionResult) Interrupts() bool {
	switch r.Action {
	case "Block", "Captcha", "Challenge", "Monetize":
		return true
	}
	return false
}

// WebACLInspector evaluates the WebACL associated with a protected
// resource against an incoming request. Protected-resource planes
// (CloudFront distributions, API Gateway stages, AppSync APIs, Cognito
// user pools) call this before serving a request; resources without an
// association are allowed through without inspection.
type WebACLInspector interface {
	// InspectWebACLRequest resolves the WebACL associated with the
	// resource ARN and evaluates it against the request. The region
	// selects the regional association and WebACL stores; CLOUDFRONT
	// scope associations live in the global storage.
	InspectWebACLRequest(ctx context.Context, region, resourceArn string, req WebACLInspectionRequest) (WebACLInspectionResult, error)
}

// ChallengeEndpointPath is the reserved request path at which every
// protected listener serves the aws-waf-token exchange. AWS serves the
// challenge scripts and token endpoint from its own edge domains; an
// edge deployment has no such domains, so the exchange lives on the
// protected listener itself and the token cookie scopes to the
// protected host.
const ChallengeEndpointPath = "/awswaf/token"

// WAFTokenExchanger serves the aws-waf-token exchange endpoint that the
// CAPTCHA and Challenge interstitials submit their solutions to.
type WAFTokenExchanger interface {
	// ExchangeWAFToken verifies an interstitial solution and answers
	// with a Set-Cookie of the aws-waf-token. The return value reports
	// whether the request was served; a request the exchanger could not
	// handle must fall through to a 400 response, never to the
	// protected resource.
	ExchangeWAFToken(ctx context.Context, w http.ResponseWriter, r *http.Request) bool
}

// ServeWAFTokenExchange routes a request to the token exchange endpoint
// when it addresses the reserved path. The return value reports whether
// the request was served.
func ServeWAFTokenExchange(ctx context.Context, inspector WebACLInspector, w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != ChallengeEndpointPath || r.Method != http.MethodPost {
		return false
	}
	exchanger, ok := inspector.(WAFTokenExchanger)
	if !ok {
		return false
	}
	return exchanger.ExchangeWAFToken(ctx, w, r)
}
