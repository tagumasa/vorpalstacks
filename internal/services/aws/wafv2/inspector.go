package wafv2

import (
	"context"
	"time"

	"strings"

	wafplane "vorpalstacks/internal/common/invokers/waf"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/wafv2/inspection"
	"vorpalstacks/internal/store/aws/waf"
	"vorpalstacks/internal/utils/aws/arn"
)

// inspectionRateTracker serves the request inspection plane: the
// sliding-window counters behind rate-based statements. It is
// package-level and shared across regions, matching the single-process
// edge model. Sampled requests are persisted through the per-region
// sampling stores instead.
var inspectionRateTracker = newRateTracker()

// InspectWebACLRequest resolves the WebACL associated with the resource
// ARN and evaluates it against the request. It implements the
// wafplane.WebACLInspector contract that every protected-resource plane
// calls before serving a request. Resources without an association are
// allowed through without inspection. CloudFront distribution ARNs use
// the global-scope association store (mirroring AssociateWebACL);
// regional resources use the association store of the given region.
func (s *WAFv2Service) InspectWebACLRequest(ctx context.Context, region, resourceArn string, req wafplane.WebACLInspectionRequest) (wafplane.WebACLInspectionResult, error) {
	result := wafplane.WebACLInspectionResult{Action: "Allow"}

	assocStore, err := s.inspectionAssociationStore(resourceArn, region)
	if err != nil {
		return result, err
	}
	assoc, err := assocStore.GetByResourceArn(resourceArn)
	if err != nil {
		if waf.IsNotFound(err) {
			return result, nil
		}
		return result, err
	}
	webACL, err := s.referencedWebACL(assoc.WebACLArn)
	if err != nil {
		return result, err
	}
	if webACL == nil {
		// The associated WebACL no longer exists; there is nothing to
		// enforce, so the request is allowed.
		return result, nil
	}
	return s.evaluateForInspection(webACL, req), nil
}

// inspectionAssociationStore returns the association store that owns
// the resource ARN: the global store for CloudFront distributions, the
// regional store for everything else. The global store is cached in
// the service's store map exactly like the association Core does.
func (s *WAFv2Service) inspectionAssociationStore(resourceArn, region string) (*waf.WebACLAssociationStore, error) {
	if parsed, err := arn.ParseARN(resourceArn); err == nil && parsed.Service == "cloudfront" {
		if cached, ok := s.stores.Load(wafv2GlobalAssocKey); ok {
			if typed, ok := cached.(*waf.WebACLAssociationStore); ok {
				return typed, nil
			}
		}
		if s.storageManager == nil {
			return nil, invalidParamError("wafv2 storage manager not initialised")
		}
		globalStorage, err := s.storageManager.GetGlobalStorage()
		if err != nil {
			return nil, err
		}
		store := waf.NewWebACLAssociationStore(globalStorage)
		if actual, loaded := s.stores.LoadOrStore(wafv2GlobalAssocKey, store); loaded {
			if typed, ok := actual.(*waf.WebACLAssociationStore); ok {
				return typed, nil
			}
		}
		return store, nil
	}
	stores, err := s.GetStoresForRegion(region)
	if err != nil {
		return nil, err
	}
	return stores.associations, nil
}

// referencedWebACL loads a WebACL by its ARN, resolving the store set
// from the region embedded in the ARN itself so inspections are
// correct even when the protecting resource lives in another region.
// A missing WebACL returns (nil, nil).
func (s *WAFv2Service) referencedWebACL(webACLArn string) (*waf.WebACL, error) {
	region := s.arnRegion(webACLArn)
	stores, err := s.GetStoresForRegion(region)
	if err != nil {
		return nil, err
	}
	webACL, err := stores.webACLs.GetByARN(webACLArn)
	if err != nil {
		if waf.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return webACL, nil
}

// arnRegion extracts the region slot of an ARN, falling back to the
// service's configured region when the ARN carries none.
func (s *WAFv2Service) arnRegion(arnString string) string {
	if parsed, err := arn.ParseARN(arnString); err == nil && parsed.Region != "" {
		return parsed.Region
	}
	return s.region
}

// evaluateForInspection runs the evaluation engine and converts the
// outcome to the cross-service result shape. Matched rules with
// sampling enabled are recorded for GetSampledRequests.
func (s *WAFv2Service) evaluateForInspection(webACL *waf.WebACL, req wafplane.WebACLInspectionRequest) wafplane.WebACLInspectionResult {
	inspReq := &inspection.Request{
		Method:        req.Method,
		URIPath:       req.URIPath,
		RawQuery:      req.RawQuery,
		SourceIP:      req.SourceIP,
		HTTPVersion:   req.HTTPVersion,
		Body:          req.Body,
		BodyTruncated: req.BodyTruncated,
		Now:           time.Now(),
	}
	for _, h := range req.Headers {
		inspReq.Headers = append(inspReq.Headers, inspection.Header{Name: h.Name, Value: h.Value})
	}
	for _, c := range req.Cookies {
		inspReq.Cookies = append(inspReq.Cookies, inspection.Header{Name: c.Name, Value: c.Value})
	}

	evaluator := inspection.NewEvaluator(inspection.Resolvers{
		IPSet:     func(arnString string) (*waf.IPSet, error) { return s.referencedIPSet(arnString) },
		RegexSet:  func(arnString string) (*waf.RegexPatternSet, error) { return s.referencedRegexSet(arnString) },
		RuleGroup: func(arnString string) (*waf.RuleGroup, error) { return s.referencedRuleGroup(arnString) },
		Rate:      inspectionRateTracker,
		Token:     s.inspectionTokenValidator(),
	})
	outcome := evaluator.Evaluate(webACL, inspReq)
	if (outcome.Action == inspection.ActionCaptcha || outcome.Action == inspection.ActionChallenge) && outcome.InterstitialRequested {
		// The interrupting response carries the interstitial page only
		// for clients that accept text/html; the page embeds a fresh
		// single-use challenge the client exchanges for the token
		// cookie at the reserved endpoint.
		if body := s.challengeInterstitial(outcome.Action); body != "" {
			outcome.CustomResponse.Body = body
		}
	}
	if len(outcome.Unsupported) > 0 {
		// Managed rules whose inspection inputs exist only inside AWS
		// (the data-dependent rule groups) and managed group names the
		// catalog does not carry are treated as non-matching; surface
		// the affected rule names so the operator sees the inspection
		// is not enforcing them.
		names := outcome.Unsupported
		if len(names) > 5 {
			names = names[:5]
		}
		logs.Warn("wafv2 rules with platform-unsupported statements were treated as non-matching",
			logs.String("webacl", webACL.ARN), logs.String("rules", strings.Join(names, ",")))
	}
	s.recordSamples(webACL, req, outcome)

	result := wafplane.WebACLInspectionResult{Action: outcome.Action}
	if outcome.CustomResponse != nil {
		result.ResponseCode = outcome.CustomResponse.StatusCode
		result.ResponseBody = outcome.CustomResponse.Body
		for _, h := range outcome.CustomResponse.Headers {
			result.ResponseHeaders = append(result.ResponseHeaders, wafplane.WebACLHTTPHeader{Name: h.Name, Value: h.Value})
		}
	}
	for _, h := range outcome.InsertHeaders {
		result.InsertHeaders = append(result.InsertHeaders, wafplane.WebACLHTTPHeader{Name: h.Name, Value: h.Value})
	}
	return result
}

// inspectionTokenValidator builds the evaluator's token validator from
// the persistent signing secret. A nil return leaves every token check
// failing — Captcha and Challenge rules then always interrupt — which
// is the fail-closed outcome when the secret is unavailable.
func (s *WAFv2Service) inspectionTokenValidator() inspection.TokenValidator {
	store, err := s.tokenStoreLoad()
	if err != nil {
		logs.Warn("wafv2 token store unavailable, captcha and challenge tokens cannot be verified",
			logs.Err(err))
		return nil
	}
	secret, err := store.SigningKey()
	if err != nil {
		logs.Warn("wafv2 token signing key unavailable, captcha and challenge tokens cannot be verified",
			logs.Err(err))
		return nil
	}
	return serviceTokenValidator{secret: secret}
}

// recordSamples persists one sampled-request record per matched rule
// whose visibility configuration enables sampling. Records live in the
// sampling store of the WebACL's own region (the same resolution the
// WebACL load uses), so they survive a server restart within the
// retention window. The terminating action's record is what
// GetSampledRequests surfaces as the request's disposition.
func (s *WAFv2Service) recordSamples(webACL *waf.WebACL, req wafplane.WebACLInspectionRequest, outcome *inspection.Result) {
	if !samplingEnabled(webACL.VisibilityConfig) {
		return
	}
	stores, err := s.GetStoresForRegion(s.arnRegion(webACL.ARN))
	if err != nil {
		logs.Warn("wafv2 sampling store unavailable, sample dropped", logs.String("webacl", webACL.ARN), logs.Err(err))
		return
	}
	s.startSampleRetentionSweep()
	// The response code, applied labels and inserted headers describe
	// the request as a whole, so every per-rule record of this request
	// carries the same values. Every terminating action resolves a
	// response; an allowed request passes through to the origin, which
	// answers 200 unless the origin itself fails.
	responseCode := 200
	if outcome.CustomResponse != nil {
		responseCode = outcome.CustomResponse.StatusCode
	}
	inserted := make([]waf.SampledHTTPHeader, 0, len(outcome.InsertHeaders))
	for _, h := range outcome.InsertHeaders {
		inserted = append(inserted, waf.SampledHTTPHeader{Name: h.Name, Value: h.Value})
	}
	for _, matched := range outcome.MatchedRules {
		// A rule's own visibility configuration overrides the web
		// ACL's sampling setting for that rule.
		if matched.SampledRequestsEnabled != nil && !*matched.SampledRequestsEnabled {
			continue
		}
		metricName := matched.MetricName
		if metricName == "" {
			metricName = matched.RuleName
		}
		record := waf.SampledRequest{
			RuleNameWithinRuleGroup: matched.RuleNameWithinRuleGroup,
			MetricName:              metricName,
			Action:                  strings.ToUpper(matched.Action),
			Timestamp:               time.Now(),
			ClientIP:                req.SourceIP,
			URI:                     req.URIPath,
			Method:                  req.Method,
			HTTPVersion:             req.HTTPVersion,
			Headers:                 sampledHeaders(req.Headers),
			ResponseCodeSent:        responseCode,
			Labels:                  outcome.Labels,
			RequestHeadersInserted:  inserted,
			OverriddenAction:        strings.ToUpper(matched.OverriddenAction),
			Captcha:                 sampledTokenInspection(matched.Captcha),
			Challenge:               sampledTokenInspection(matched.Challenge),
		}
		if err := stores.samples.Record(webACL.ARN, metricName, record); err != nil {
			logs.Warn("wafv2 sampled-request persistence failed, sample dropped", logs.String("webacl", webACL.ARN), logs.Err(err))
		}
	}
}

func samplingEnabled(vc *waf.VisibilityConfig) bool {
	return vc != nil && vc.SampledRequestsEnabled
}

func sampledHeaders(headers []wafplane.WebACLHTTPHeader) []waf.SampledHTTPHeader {
	out := make([]waf.SampledHTTPHeader, 0, len(headers))
	for _, h := range headers {
		out = append(out, waf.SampledHTTPHeader{Name: h.Name, Value: h.Value})
	}
	return out
}

// sampledTokenInspection converts one token-inspection outcome to its
// persisted form; a nil outcome stays nil so the record omits it.
func sampledTokenInspection(inspection *inspection.TokenInspection) *waf.TokenInspectionRecord {
	if inspection == nil {
		return nil
	}
	return &waf.TokenInspectionRecord{
		SolveTimestamp: inspection.SolveTimestamp,
		FailureReason:  inspection.FailureReason,
	}
}

// referencedIPSet loads an IP set referenced by a statement, resolving
// the store set from the ARN's own region.
func (s *WAFv2Service) referencedIPSet(arnString string) (*waf.IPSet, error) {
	stores, err := s.GetStoresForRegion(s.arnRegion(arnString))
	if err != nil {
		return nil, err
	}
	return stores.ipSets.GetByARN(arnString)
}

// referencedRegexSet loads a regex pattern set referenced by a
// statement.
func (s *WAFv2Service) referencedRegexSet(arnString string) (*waf.RegexPatternSet, error) {
	stores, err := s.GetStoresForRegion(s.arnRegion(arnString))
	if err != nil {
		return nil, err
	}
	return stores.regexPatternSets.GetByARN(arnString)
}

// referencedRuleGroup loads a rule group referenced by a statement.
func (s *WAFv2Service) referencedRuleGroup(arnString string) (*waf.RuleGroup, error) {
	stores, err := s.GetStoresForRegion(s.arnRegion(arnString))
	if err != nil {
		return nil, err
	}
	return stores.ruleGroups.GetByARN(arnString)
}
