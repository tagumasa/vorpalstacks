package wafv2

import (
	"encoding/json"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/wafv2/inspection"
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

// WCU cost constants. The per-statement characteristics sections of the
// AWS WAF Developer Guide define them: a string match costs 2 WCUs for
// the anchored positional constraints (EXACTLY, STARTS_WITH, ENDS_WITH)
// and 10 for the substring ones (CONTAINS, CONTAINS_WORD); each text
// transformation adds 10 WCUs; the AllQueryArguments component adds 10
// WCUs; the JSON body component doubles the statement's base cost; a
// rate-based statement costs 2 WCUs plus 30 per custom aggregation key.
// The NONE transformation costs nothing: the console attaches it by
// default to every rule, and the guide's base-cost examples for plain
// rules show the base cost alone, so a no-op transformation cannot be
// charged. These are the single definitions of the values; every other
// site references them.
const (
	wcuByteMatchAnchored  int64 = 2
	wcuByteMatchSubstring int64 = 10
	wcuTextTransformation int64 = 10
	wcuAllQueryArguments  int64 = 10
	wcuRateCustomKey      int64 = 30
	wcuRateBase           int64 = 2
)

// calculateStatementCapacity returns the WCU capacity consumed by a single
// statement. Base WCU values per the per-statement characteristics
// sections of the AWS WAF Developer Guide:
//
//	ByteMatchStatement: 2 (anchored) or 10 (substring) plus modifiers,
//	SqliMatchStatement: 20, XssMatchStatement: 20,
//	SizeConstraintStatement: 1, GeoMatchStatement: 1,
//	RegexMatchStatement: 25, RegexPatternSetRefStatement: 25,
//	IPSetReferenceStatement: 1, LabelMatchStatement: 1,
//	AndStatement/OrStatement: 1 + nested, NotStatement: 1 + nested,
//	RateBasedStatement: 2 + 30 per custom key + scope-down,
//	ManagedRuleGroupStatement: the managed group's catalog WCU
func calculateStatementCapacity(stmt *waf.Statement) int64 {
	if stmt == nil {
		return 0
	}
	if stmt.ByteMatchStatement != nil {
		base := wcuByteMatchAnchored
		switch stmt.ByteMatchStatement.PositionalConstraint {
		case "CONTAINS", "CONTAINS_WORD":
			base = wcuByteMatchSubstring
		}
		return matchStatementCost(base, stmt.ByteMatchStatement.FieldToMatch, stmt.ByteMatchStatement.TextTransformations)
	}
	if stmt.SqliMatchStatement != nil {
		return matchStatementCost(20, stmt.SqliMatchStatement.FieldToMatch, stmt.SqliMatchStatement.TextTransformations)
	}
	if stmt.XssMatchStatement != nil {
		return matchStatementCost(20, stmt.XssMatchStatement.FieldToMatch, stmt.XssMatchStatement.TextTransformations)
	}
	if stmt.SizeConstraintStatement != nil {
		return matchStatementCost(1, stmt.SizeConstraintStatement.FieldToMatch, stmt.SizeConstraintStatement.TextTransformations)
	}
	if stmt.GeoMatchStatement != nil {
		return 1
	}
	if stmt.RegexMatchStatement != nil {
		return matchStatementCost(25, stmt.RegexMatchStatement.FieldToMatch, stmt.RegexMatchStatement.TextTransformations)
	}
	if stmt.RegexPatternSetRefStatement != nil {
		return matchStatementCost(25, stmt.RegexPatternSetRefStatement.FieldToMatch, stmt.RegexPatternSetRefStatement.TextTransformations)
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
		capacity := wcuRateBase + wcuRateCustomKey*int64(len(stmt.RateBasedStatement.CustomKeys))
		if stmt.RateBasedStatement.ScopeDownStatement != nil {
			capacity += calculateStatementCapacity(stmt.RateBasedStatement.ScopeDownStatement)
		}
		return capacity
	}
	if stmt.ManagedRuleGroupStatement != nil {
		// A managed rule group's capacity is its catalog WCU, the same
		// number DescribeManagedRuleGroup reports; it counts against the
		// web ACL's capacity like any other statement.
		if group, ok := inspection.LookupManagedRuleGroup(
			stmt.ManagedRuleGroupStatement.VendorName,
			stmt.ManagedRuleGroupStatement.Name,
		); ok {
			return group.WCU
		}
		return 0
	}
	if stmt.RuleGroupReferenceStatement != nil {
		return 0
	}
	return 1
}

// matchStatementCost applies the request-component and text-
// transformation modifiers the Developer Guide defines for match
// statements: the JSON body component doubles the base cost, the
// AllQueryArguments component adds 10 WCUs, and every text
// transformation other than NONE adds 10 WCUs.
func matchStatementCost(base int64, ftm *waf.FieldToMatch, tts []*waf.TextTransformation) int64 {
	if ftm != nil {
		if ftm.JsonBody != nil {
			base *= 2
		}
		if ftm.AllQueryArguments != nil {
			base += wcuAllQueryArguments
		}
	}
	for _, tt := range tts {
		if tt != nil && tt.Type != "NONE" {
			base += wcuTextTransformation
		}
	}
	return base
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
// required fields for web ACL rules — managed rule group references are
// allowed as a rule's top-level statement. Returns an error if any
// statement contains an invalid enum value or is missing a required
// field.
func parseRules(rulesRaw interface{}) ([]*waf.Rule, error) {
	return parseRulesForEntity(rulesRaw, true)
}

// parseRulesForEntity is parseRules with the entity's placement rules:
// rule groups cannot reference managed rule groups at any depth, per
// the API's statement placement rules.
func parseRulesForEntity(rulesRaw interface{}, allowManagedRuleGroups bool) ([]*waf.Rule, error) {
	rules := convertRules(rulesRaw)
	if err := validateRules(rules, allowManagedRuleGroups); err != nil {
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
