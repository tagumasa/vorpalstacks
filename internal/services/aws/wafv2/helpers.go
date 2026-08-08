package wafv2

import (
	"encoding/json"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/aws/waf"
)

// marshalUnmarshal converts a value to JSON and back into a different type.
// This is used for conversion between raw maps from the AWS request layer
// and the typed WAF store structs. Because both sides use PascalCase JSON
// keys (AWS convention), the roundtrip preserves all fields faithfully.
func marshalUnmarshal[In any, Out any](in In) (Out, bool) {
	var zero Out
	data, err := json.Marshal(in)
	if err != nil {
		return zero, false
	}
	var out Out
	if err := json.Unmarshal(data, &out); err != nil {
		return zero, false
	}
	return out, true
}

// calculateStatementCapacity returns the WCU capacity consumed by a single
// statement. Base WCU values per AWS WAF documentation:
//
//	ByteMatchStatement: 1, SqliMatchStatement: 20, XssMatchStatement: 20,
//	SizeConstraintStatement: 1, GeoMatchStatement: 1,
//	RegexMatchStatement: 25, RegexPatternSetRefStatement: 25,
//	IPSetReferenceStatement: 1, LabelMatchStatement: 1,
//	AndStatement/OrStatement: 1 + sum of nested,
//	NotStatement: 1 + nested, RateBasedStatement: 2 + nested,
//	ManagedRuleGroupStatement: 0 (capacity is provided by DescribeManagedRuleGroup)
func calculateStatementCapacity(stmt *waf.Statement) int64 {
	if stmt == nil {
		return 0
	}
	if stmt.ByteMatchStatement != nil {
		return 1
	}
	if stmt.SqliMatchStatement != nil {
		return 20
	}
	if stmt.XssMatchStatement != nil {
		return 20
	}
	if stmt.SizeConstraintStatement != nil {
		return 1
	}
	if stmt.GeoMatchStatement != nil {
		return 1
	}
	if stmt.RegexMatchStatement != nil {
		return 25
	}
	if stmt.RegexPatternSetRefStatement != nil {
		return 25
	}
	if stmt.IPSetReferenceStatement != nil {
		return 1
	}
	if stmt.LabelMatchStatement != nil {
		return 1
	}
	if stmt.AsnMatchStatement != nil {
		return 1
	}
	if stmt.AndStatement != nil {
		var total int64 = 1
		for _, s := range stmt.AndStatement.Statements {
			total += calculateStatementCapacity(s)
		}
		return total
	}
	if stmt.OrStatement != nil {
		var maxNested int64
		for _, s := range stmt.OrStatement.Statements {
			if c := calculateStatementCapacity(s); c > maxNested {
				maxNested = c
			}
		}
		return 1 + maxNested
	}
	if stmt.NotStatement != nil {
		return 1 + calculateStatementCapacity(stmt.NotStatement.Statement)
	}
	if stmt.RateBasedStatement != nil {
		var nested int64
		if stmt.RateBasedStatement.ScopeDownStatement != nil {
			nested = calculateStatementCapacity(stmt.RateBasedStatement.ScopeDownStatement)
		}
		return 2 + nested
	}
	if stmt.ManagedRuleGroupStatement != nil {
		return 0
	}
	if stmt.RuleGroupReferenceStatement != nil {
		return 0
	}
	return 1
}

// convertVisibilityConfig converts a raw map to a typed VisibilityConfig.
// Returns nil if the input is empty — callers should validate that
// VisibilityConfig is provided where required by the Smithy model.
func convertVisibilityConfig(m map[string]interface{}) *waf.VisibilityConfig {
	if len(m) == 0 {
		return nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	var vc waf.VisibilityConfig
	if err := json.Unmarshal(data, &vc); err != nil {
		return nil
	}
	return &vc
}

// convertVisibilityConfigToResponse serialises a VisibilityConfig back to
// a raw map for the API response.
func convertVisibilityConfigToResponse(vc *waf.VisibilityConfig) map[string]interface{} {
	if vc == nil {
		return nil
	}
	data, err := json.Marshal(vc)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

// convertAction converts a raw map to a typed Action. Handles all six
// AWS action types: Allow, Block, Count, Captcha, Challenge, Monetize.
func convertAction(m map[string]interface{}) *waf.Action {
	if len(m) == 0 {
		return nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	var a waf.Action
	if err := json.Unmarshal(data, &a); err != nil {
		return nil
	}
	return &a
}

// convertActionToResponse serialises an Action back to a raw map. The
// input may be either *waf.Action (pre-storage) or map[string]interface{}
// (post-deserialisation from Pebble storage).
func convertActionToResponse(a interface{}) map[string]interface{} {
	if a == nil {
		return nil
	}
	data, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	if len(result) == 0 {
		return nil
	}
	return result
}

// convertRules converts a raw Rules array from the AWS request into typed
// Rule structs. Each rule's Action, OverrideAction, Statement, and
// VisibilityConfig are fully populated via JSON roundtrip.
func convertRules(rulesRaw interface{}) []*waf.Rule {
	if rulesRaw == nil {
		return nil
	}
	rulesArr, ok := rulesRaw.([]interface{})
	if !ok {
		return nil
	}
	rules := make([]*waf.Rule, 0, len(rulesArr))
	for _, r := range rulesArr {
		rMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		rule := &waf.Rule{
			Name:     request.GetStringParam(rMap, "Name"),
			Priority: int32(request.GetIntParam(rMap, "Priority")),
		}
		if actionRaw, ok := rMap["Action"]; ok {
			if actionMap, ok := actionRaw.(map[string]interface{}); ok {
				rule.Action = convertAction(actionMap)
			}
		}
		if overrideActionRaw, ok := rMap["OverrideAction"]; ok {
			if overrideActionMap, ok := overrideActionRaw.(map[string]interface{}); ok {
				rule.OverrideAction = convertAction(overrideActionMap)
			}
		}
		if vcRaw, ok := rMap["VisibilityConfig"]; ok {
			if vcMap, ok := vcRaw.(map[string]interface{}); ok {
				rule.VisibilityConfig = convertVisibilityConfig(vcMap)
			}
		}
		if statementRaw, ok := rMap["Statement"]; ok {
			if statementMap, ok := statementRaw.(map[string]interface{}); ok {
				rule.Statement = convertStatement(statementMap)
			}
		}
		if rlRaw, ok := rMap["RuleLabels"]; ok {
			rule.RuleLabels = rlRaw
		}
		if ccRaw, ok := rMap["CaptchaConfig"]; ok {
			rule.CaptchaConfig = ccRaw
		}
		if chRaw, ok := rMap["ChallengeConfig"]; ok {
			rule.ChallengeConfig = chRaw
		}
		rules = append(rules, rule)
	}
	return rules
}

// parseRules converts raw rule data and validates all enum values and
// required fields. Returns an error if any statement contains an invalid
// enum value or is missing a required field.
func parseRules(rulesRaw interface{}) ([]*waf.Rule, error) {
	rules := convertRules(rulesRaw)
	if err := validateRules(rules); err != nil {
		return nil, err
	}
	return rules, nil
}

// convertStatement converts a raw map to a typed Statement via JSON
// roundtrip. This faithfully preserves all sixteen Smithy statement types
// (ByteMatchStatement, GeoMatchStatement, RateBasedStatement,
// RegexPatternSetReferenceStatement, AndStatement, OrStatement, etc.)
// with their full data including nested FieldToMatch and
// TextTransformations.
func convertStatement(m map[string]interface{}) *waf.Statement {
	if len(m) == 0 {
		return nil
	}
	stmt, ok := marshalUnmarshal[map[string]interface{}, waf.Statement](m)
	if !ok {
		return nil
	}
	return &stmt
}

// convertRulesToResponse serialises typed Rules back to a slice of raw
// maps for the API response.
func convertRulesToResponse(rules []*waf.Rule) []interface{} {
	if rules == nil {
		return nil
	}
	result := make([]interface{}, 0, len(rules))
	for _, r := range rules {
		if r == nil {
			continue
		}
		m := map[string]interface{}{
			"Name":     r.Name,
			"Priority": r.Priority,
		}
		if r.Statement != nil {
			m["Statement"] = convertStatementToResponse(r.Statement)
		}
		if r.Action != nil {
			m["Action"] = convertActionToResponse(r.Action)
		}
		if r.OverrideAction != nil {
			m["OverrideAction"] = convertActionToResponse(r.OverrideAction)
		}
		if r.VisibilityConfig != nil {
			m["VisibilityConfig"] = convertVisibilityConfigToResponse(r.VisibilityConfig)
		}
		if r.RuleLabels != nil {
			m["RuleLabels"] = r.RuleLabels
		}
		if r.CaptchaConfig != nil {
			m["CaptchaConfig"] = r.CaptchaConfig
		}
		if r.ChallengeConfig != nil {
			m["ChallengeConfig"] = r.ChallengeConfig
		}
		result = append(result, m)
	}
	return result
}

// convertStatementToResponse serialises a typed Statement back to a raw
// map for the API response. Only the non-nil statement type is included
// (AWS Statement is a union — exactly one type should be present).
func convertStatementToResponse(s *waf.Statement) map[string]interface{} {
	if s == nil {
		return nil
	}
	data, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	if len(result) == 0 {
		return nil
	}
	return result
}

func buildWebACLSummary(wa *waf.WebACL) map[string]interface{} {
	return map[string]interface{}{
		"Id":          wa.ID,
		"Name":        wa.Name,
		"ARN":         wa.ARN,
		"Description": wa.Description,
		"LockToken":   wa.LockToken,
	}
}

func buildWebACLSummaryList(webACLs []*waf.WebACL) []interface{} {
	result := make([]interface{}, 0, len(webACLs))
	for _, wa := range webACLs {
		result = append(result, buildWebACLSummary(wa))
	}
	return result
}

func buildRuleGroupSummary(rg *waf.RuleGroup) map[string]interface{} {
	return map[string]interface{}{
		"Id":          rg.ID,
		"Name":        rg.Name,
		"ARN":         rg.ARN,
		"Description": rg.Description,
		"LockToken":   rg.LockToken,
	}
}

func buildIPSetSummary(ips *waf.IPSet) map[string]interface{} {
	return map[string]interface{}{
		"Id":          ips.ID,
		"Name":        ips.Name,
		"ARN":         ips.ARN,
		"Description": ips.Description,
		"LockToken":   ips.LockToken,
	}
}

func buildRegexPatternSetSummary(rps *waf.RegexPatternSet) map[string]interface{} {
	return map[string]interface{}{
		"Id":          rps.ID,
		"Name":        rps.Name,
		"ARN":         rps.ARN,
		"Description": rps.Description,
		"LockToken":   rps.LockToken,
	}
}
