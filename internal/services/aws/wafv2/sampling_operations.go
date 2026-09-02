package wafv2

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// GetSampledRequests returns a sample of requests that matched the rule
// identified by its metric name, within the previous three hours.
func (s *WAFv2Service) GetSampledRequests(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	timeWindow := request.GetMapParam(req.Parameters, "TimeWindow")
	_, maxItemsPresent := req.Parameters["MaxItems"]

	result, err := s.getSampledRequestsCore(reqCtx, GetSampledRequestsInput{
		WebACLArn:         request.GetStringParam(req.Parameters, "WebAclArn"),
		RuleMetricName:    request.GetStringParam(req.Parameters, "RuleMetricName"),
		Scope:             request.GetStringParam(req.Parameters, "Scope"),
		StartTime:         parseSampledRequestTime(timeWindow, "StartTime"),
		EndTime:           parseSampledRequestTime(timeWindow, "EndTime"),
		MaxItems:          request.GetInt64Param(req.Parameters, "MaxItems"),
		TimeWindowPresent: timeWindow != nil,
		MaxItemsPresent:   maxItemsPresent,
	})
	if err != nil {
		return nil, err
	}

	sampledRequests := make([]interface{}, 0, len(result.SampledRequests))
	for _, record := range result.SampledRequests {
		headers := make([]interface{}, 0, len(record.Headers))
		for _, h := range record.Headers {
			headers = append(headers, map[string]interface{}{
				"Name":  h.Name,
				"Value": h.Value,
			})
		}
		labels := make([]interface{}, 0, len(record.Labels))
		for _, label := range record.Labels {
			labels = append(labels, map[string]interface{}{"Name": label})
		}
		insertedHeaders := make([]interface{}, 0, len(record.RequestHeadersInserted))
		for _, h := range record.RequestHeadersInserted {
			insertedHeaders = append(insertedHeaders, map[string]interface{}{
				"Name":  h.Name,
				"Value": h.Value,
			})
		}
		entry := map[string]interface{}{
			"Request": map[string]interface{}{
				"ClientIP":    record.ClientIP,
				"URI":         record.URI,
				"Method":      record.Method,
				"HTTPVersion": record.HTTPVersion,
				"Headers":     headers,
			},
			// Every matched request is retained, so each sample
			// represents exactly one request.
			"Weight":    int64(1),
			"Timestamp": record.Timestamp,
			"Action":    record.Action,
		}
		// The within-group name exists only for matches from inside a
		// rule group; a rule declared directly in the web ACL leaves
		// the member absent.
		if record.RuleNameWithinRuleGroup != "" {
			entry["RuleNameWithinRuleGroup"] = record.RuleNameWithinRuleGroup
		}
		if record.ResponseCodeSent != 0 {
			entry["ResponseCodeSent"] = record.ResponseCodeSent
		}
		if len(labels) > 0 {
			entry["Labels"] = labels
		}
		if len(insertedHeaders) > 0 {
			entry["RequestHeadersInserted"] = insertedHeaders
		}
		if record.OverriddenAction != "" {
			entry["OverriddenAction"] = record.OverriddenAction
		}
		if record.Captcha != nil {
			entry["CaptchaResponse"] = sampledTokenResponse(record.Captcha)
		}
		if record.Challenge != nil {
			entry["ChallengeResponse"] = sampledTokenResponse(record.Challenge)
		}
		sampledRequests = append(sampledRequests, entry)
	}

	return map[string]interface{}{
		"SampledRequests": sampledRequests,
		"PopulationSize":  result.PopulationSize,
		"TimeWindow": map[string]interface{}{
			"StartTime": result.StartTime,
			"EndTime":   result.EndTime,
		},
	}, nil
}

// sampledTokenResponse renders one token-inspection outcome as the
// CaptchaResponse or ChallengeResponse member of a sampled request.
func sampledTokenResponse(record *wafstore.TokenInspectionRecord) map[string]interface{} {
	response := map[string]interface{}{}
	if record.SolveTimestamp > 0 {
		response["SolveTimestamp"] = time.Unix(record.SolveTimestamp, 0)
	}
	if record.FailureReason != "" {
		response["FailureReason"] = record.FailureReason
	}
	return response
}

// GetRateBasedStatementManagedKeys returns the IP addresses the
// rate-based statement currently aggregates under.
func (s *WAFv2Service) GetRateBasedStatementManagedKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getRateBasedManagedKeysCore(reqCtx, GetRateBasedManagedKeysInput{
		Scope:             request.GetStringParam(req.Parameters, "Scope"),
		WebACLName:        request.GetStringParam(req.Parameters, "WebACLName"),
		WebACLId:          request.GetStringParam(req.Parameters, "WebACLId"),
		RuleName:          request.GetStringParam(req.Parameters, "RuleName"),
		RuleGroupRuleName: request.GetStringParam(req.Parameters, "RuleGroupRuleName"),
	})
	if err != nil {
		return nil, err
	}

	managedKeys := map[string]interface{}{}
	if len(result.IPv4) > 0 {
		managedKeys["ManagedKeysIPV4"] = map[string]interface{}{
			"IPAddressVersion": "IPV4",
			"Addresses":        result.IPv4,
		}
	}
	if len(result.IPv6) > 0 {
		managedKeys["ManagedKeysIPV6"] = map[string]interface{}{
			"IPAddressVersion": "IPV6",
			"Addresses":        result.IPv6,
		}
	}
	return managedKeys, nil
}

// parseSampledRequestTime reads one TimeWindow member. The JSON
// protocol carries timestamps as epoch seconds, but parsed values may
// already be time.Time; RFC 3339 strings are accepted as well.
func parseSampledRequestTime(timeWindow map[string]interface{}, key string) time.Time {
	if timeWindow == nil {
		return time.Time{}
	}
	switch value := timeWindow[key].(type) {
	case time.Time:
		return value
	case int64:
		return time.Unix(value, 0)
	case int:
		return time.Unix(int64(value), 0)
	case float64:
		return time.Unix(int64(value), 0)
	case string:
		if parsed, err := time.Parse(time.RFC3339, value); err == nil {
			return parsed
		}
	}
	return time.Time{}
}
