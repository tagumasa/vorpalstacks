package wafv2

import (
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// GetSampledRequestsInput is the transport-agnostic input for the
// GetSampledRequests operation. Presence flags separate an omitted
// member from its zero value because the required-member validation
// distinguishes the two.
type GetSampledRequestsInput struct {
	WebACLArn         string
	RuleMetricName    string
	Scope             string
	StartTime         time.Time
	EndTime           time.Time
	MaxItems          int64
	TimeWindowPresent bool
	MaxItemsPresent   bool
	// Now allows tests to pin the clamping instant; zero means
	// time.Now().
	Now time.Time
}

// GetSampledRequestsResult carries the retained samples of the rule
// inside the requested window plus the effective (clamped) window.
type GetSampledRequestsResult struct {
	SampledRequests []wafstore.SampledRequest
	PopulationSize  int64
	StartTime       time.Time
	EndTime         time.Time
}

// getSampledRequestsCore validates the request and reads the sampled
// requests of one rule metric from the sampling store of the WebACL's
// own region. The retention store keeps every matched request up to
// SamplingPopulationDepth per rule and counts every match in minute
// buckets, so PopulationSize is the number of matched requests in the
// window capped at the depth — the population the MaxItems
// documentation draws its sample from (the first 5,000 requests of the
// time range).
func (s *WAFv2Service) getSampledRequestsCore(reqCtx *request.RequestContext, input GetSampledRequestsInput) (*GetSampledRequestsResult, error) {
	if input.WebACLArn == "" {
		return nil, invalidParamError("WebAclArn is required")
	}
	if input.RuleMetricName == "" {
		return nil, invalidParamError("RuleMetricName is required")
	}
	if err := validateScope(input.Scope); err != nil {
		return nil, err
	}
	if !input.TimeWindowPresent {
		return nil, invalidParamError("TimeWindow is required")
	}
	if input.StartTime.IsZero() {
		return nil, invalidParamError("TimeWindow StartTime is required")
	}
	if input.EndTime.IsZero() {
		return nil, invalidParamError("TimeWindow EndTime is required")
	}
	if !input.EndTime.After(input.StartTime) {
		return nil, invalidParamError("TimeWindow EndTime must be later than StartTime")
	}
	if !input.MaxItemsPresent {
		return nil, invalidParamError("MaxItems is required")
	}
	if input.MaxItems < 1 || input.MaxItems > wafstore.MaxSampledRequests {
		return nil, invalidParamError(fmt.Sprintf("MaxItems must be between 1 and %d", wafstore.MaxSampledRequests))
	}

	stores, err := s.GetStoresForRegion(s.arnRegion(input.WebACLArn))
	if err != nil {
		return nil, err
	}
	webACL, err := stores.webACLs.GetByARN(input.WebACLArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	// Samples are retrievable only for the previous three hours; a
	// start earlier than that is clamped, per the AWS WAF API
	// Reference.
	now := input.Now
	if now.IsZero() {
		now = time.Now()
	}
	start := input.StartTime
	if cutoff := now.Add(-wafstore.SampleRetention); start.Before(cutoff) {
		start = cutoff
	}

	records, err := stores.samples.Query(webACL.ARN, input.RuleMetricName, start, input.EndTime, int(input.MaxItems))
	if err != nil {
		return nil, err
	}
	population := stores.samples.CountPopulation(webACL.ARN, input.RuleMetricName, start, input.EndTime)
	if population > wafstore.SamplingPopulationDepth {
		population = wafstore.SamplingPopulationDepth
	}
	return &GetSampledRequestsResult{
		SampledRequests: records,
		PopulationSize:  population,
		StartTime:       start,
		EndTime:         input.EndTime,
	}, nil
}

// GetRateBasedManagedKeysInput is the transport-agnostic input for the
// GetRateBasedStatementManagedKeys operation.
type GetRateBasedManagedKeysInput struct {
	Scope             string
	WebACLName        string
	WebACLId          string
	RuleName          string
	RuleGroupRuleName string
}

// ManagedKeysResult carries the currently tracked rate-aggregation
// addresses split by IP version.
type ManagedKeysResult struct {
	IPv4 []string
	IPv6 []string
}

// getRateBasedManagedKeysCore validates the request, locates the named
// rate-based rule (top level, or inside the rule group referenced by
// the web ACL rule named RuleGroupRuleName), and returns the
// aggregation keys the inspection tracker currently holds for it.
// Rate aggregation keys address the rule by its own name, so the
// nested form resolves to the same keys as the bare rule name.
func (s *WAFv2Service) getRateBasedManagedKeysCore(reqCtx *request.RequestContext, input GetRateBasedManagedKeysInput) (*ManagedKeysResult, error) {
	if err := validateScope(input.Scope); err != nil {
		return nil, err
	}
	if input.WebACLName == "" {
		return nil, invalidParamError("WebACLName is required")
	}
	if input.WebACLId == "" {
		return nil, invalidParamError("WebACLId is required")
	}
	if input.RuleName == "" {
		return nil, invalidParamError("RuleName is required")
	}

	stores, err := s.storeForScope(reqCtx, input.Scope)
	if err != nil {
		return nil, err
	}
	webACL, err := stores.webACLs.Get(input.WebACLId)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}
	if webACL.Name != input.WebACLName {
		return nil, notFoundError("WebACL")
	}

	rateStmt, err := s.locateRateStatement(stores, webACL, input.RuleName, input.RuleGroupRuleName)
	if err != nil {
		return nil, err
	}
	if rateStmt == nil {
		return nil, notFoundError("Rule")
	}
	// The managed-keys call exists only for rules that aggregate on
	// the client or forwarded client address; an omitted aggregate key
	// defaults to IP.
	aggregate := rateStmt.AggregateKeyType
	if aggregate == "" {
		aggregate = "IP"
	}
	if aggregate != "IP" && aggregate != "FORWARDED_IP" {
		return nil, newAPIError("WAFUnsupportedAggregateKeyTypeException",
			"The rule that you've named doesn't aggregate solely on the IP address or solely on the forwarded IP address. This call is only available for rate-based rules with an AggregateKeyType setting of IP or FORWARDED_IP.",
			400)
	}

	ipv4, ipv6 := inspectionRateTracker.ActiveIPKeys(webACL.ARN, input.RuleName)
	return &ManagedKeysResult{IPv4: ipv4, IPv6: ipv6}, nil
}

// locateRateStatement finds the rate-based statement of the named
// rule: a top-level web ACL rule when RuleGroupRuleName is empty,
// otherwise the rule of that name inside the rule group referenced by
// the web ACL rule named RuleGroupRuleName. A nil statement with a nil
// error means the rule exists but is not rate based.
func (s *WAFv2Service) locateRateStatement(stores *wafv2Stores, webACL *wafstore.WebACL, ruleName, ruleGroupRuleName string) (*wafstore.RateBasedStatement, error) {
	if ruleGroupRuleName == "" {
		for _, rule := range webACL.Rules {
			if rule != nil && rule.Name == ruleName && rule.Statement != nil {
				return rule.Statement.RateBasedStatement, nil
			}
		}
		return nil, nil
	}

	for _, rule := range webACL.Rules {
		if rule == nil || rule.Statement == nil || rule.Statement.RuleGroupReferenceStatement == nil {
			continue
		}
		if rule.Name != ruleGroupRuleName {
			continue
		}
		group, err := stores.ruleGroups.GetByARN(rule.Statement.RuleGroupReferenceStatement.ARN)
		if err != nil {
			return nil, err
		}
		for _, inner := range group.Rules {
			if inner != nil && inner.Name == ruleName && inner.Statement != nil {
				return inner.Statement.RateBasedStatement, nil
			}
		}
		return nil, nil
	}
	return nil, nil
}
