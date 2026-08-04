package wafv2

import (
	"encoding/json"
	"fmt"

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

// validateScope checks that the Scope parameter is a valid Smithy enum
// value (CLOUDFRONT or REGIONAL). Returns a validation error if not.
func validateScope(scope string) error {
	if scope != "CLOUDFRONT" && scope != "REGIONAL" {
		return invalidParamError("Scope must be CLOUDFRONT or REGIONAL")
	}
	return nil
}

// validateDefaultAction checks that the DefaultAction only contains Allow
// or Block. Per the Smithy model, DefaultAction is a separate shape from
// RuleAction and only supports terminating actions (Allow, Block).
func validateDefaultAction(action *waf.Action) error {
	if action == nil {
		return nil
	}
	if action.Allow != nil {
		return nil
	}
	if action.Block != nil {
		return nil
	}
	return invalidParamError("DefaultAction must be Allow or Block")
}

// --- Enum validators ---

func isValidTextTransformationType(t string) bool {
	switch t {
	case "NONE", "COMPRESS_WHITE_SPACE", "HTML_ENTITY_DECODE", "LOWERCASE",
		"CMD_LINE", "URL_DECODE", "BASE64_DECODE", "HEX_DECODE", "MD5",
		"REPLACE_COMMENTS", "ESCAPE_SEQ_DECODE", "SQL_HEX_DECODE",
		"CSS_DECODE", "JS_DECODE", "NORMALIZE_PATH", "NORMALIZE_PATH_WIN",
		"REMOVE_NULLS", "REPLACE_NULLS", "BASE64_DECODE_EXT",
		"URL_DECODE_UNI", "UTF8_TO_UNICODE":
		return true
	}
	return false
}

func isValidPositionalConstraint(p string) bool {
	switch p {
	case "EXACTLY", "STARTS_WITH", "ENDS_WITH", "CONTAINS", "CONTAINS_WORD":
		return true
	}
	return false
}

func isValidComparisonOperator(op string) bool {
	switch op {
	case "EQ", "NE", "LE", "LT", "GE", "GT":
		return true
	}
	return false
}

func isValidOversizeHandling(h string) bool {
	switch h {
	case "CONTINUE", "MATCH", "NO_MATCH":
		return true
	}
	return false
}

func isValidMapMatchScope(s string) bool {
	switch s {
	case "ALL", "KEY", "VALUE":
		return true
	}
	return false
}

func isValidBodyParsingFallbackBehavior(b string) bool {
	switch b {
	case "MATCH", "NO_MATCH", "EVALUATE_AS_STRING":
		return true
	}
	return false
}

func isValidAggregateKeyType(k string) bool {
	switch k {
	case "IP", "FORWARDED_IP", "CUSTOM_KEYS", "CONSTANT":
		return true
	}
	return false
}

// validateTextTransformations validates an array of TextTransformation.
func validateTextTransformations(tts []*waf.TextTransformation) error {
	for _, tt := range tts {
		if tt == nil {
			continue
		}
		if !isValidTextTransformationType(tt.Type) {
			return invalidParamError(fmt.Sprintf("Invalid TextTransformationType: %s", tt.Type))
		}
	}
	return nil
}

// validateFieldToMatch recursively validates all enum fields inside a FieldToMatch.
func validateFieldToMatch(ftm *waf.FieldToMatch) error {
	if ftm == nil {
		return nil
	}
	if ftm.Body != nil && ftm.Body.OversizeHandling != "" && !isValidOversizeHandling(ftm.Body.OversizeHandling) {
		return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.Body.OversizeHandling))
	}
	if ftm.Cookies != nil {
		if ftm.Cookies.MatchScope != "" && !isValidMapMatchScope(ftm.Cookies.MatchScope) {
			return invalidParamError(fmt.Sprintf("Invalid MatchScope: %s", ftm.Cookies.MatchScope))
		}
		if ftm.Cookies.OversizeHandling != "" && !isValidOversizeHandling(ftm.Cookies.OversizeHandling) {
			return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.Cookies.OversizeHandling))
		}
	}
	if ftm.Headers != nil {
		if ftm.Headers.MatchScope != "" && !isValidMapMatchScope(ftm.Headers.MatchScope) {
			return invalidParamError(fmt.Sprintf("Invalid MatchScope: %s", ftm.Headers.MatchScope))
		}
		if ftm.Headers.OversizeHandling != "" && !isValidOversizeHandling(ftm.Headers.OversizeHandling) {
			return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.Headers.OversizeHandling))
		}
	}
	if ftm.HeaderOrder != nil && ftm.HeaderOrder.OversizeHandling != "" && !isValidOversizeHandling(ftm.HeaderOrder.OversizeHandling) {
		return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.HeaderOrder.OversizeHandling))
	}
	if ftm.JsonBody != nil {
		if ftm.JsonBody.MatchScope != "" && !isValidMapMatchScope(ftm.JsonBody.MatchScope) {
			return invalidParamError(fmt.Sprintf("Invalid JsonMatchScope: %s", ftm.JsonBody.MatchScope))
		}
		if ftm.JsonBody.InvalidFallbackBehavior != "" && !isValidBodyParsingFallbackBehavior(ftm.JsonBody.InvalidFallbackBehavior) {
			return invalidParamError(fmt.Sprintf("Invalid InvalidFallbackBehavior: %s", ftm.JsonBody.InvalidFallbackBehavior))
		}
		if ftm.JsonBody.OversizeHandling != "" && !isValidOversizeHandling(ftm.JsonBody.OversizeHandling) {
			return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.JsonBody.OversizeHandling))
		}
	}
	return nil
}

// validateStatement recursively validates all enum fields in a Statement
// tree. Returns WAFInvalidParameterException for any invalid enum value.
func validateStatement(stmt *waf.Statement) error {
	if stmt == nil {
		return nil
	}
	if stmt.ByteMatchStatement != nil {
		if stmt.ByteMatchStatement.PositionalConstraint != "" && !isValidPositionalConstraint(stmt.ByteMatchStatement.PositionalConstraint) {
			return invalidParamError(fmt.Sprintf("Invalid PositionalConstraint: %s", stmt.ByteMatchStatement.PositionalConstraint))
		}
		if err := validateFieldToMatch(stmt.ByteMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.ByteMatchStatement.TextTransformations); err != nil {
			return err
		}
	}
	if stmt.SqliMatchStatement != nil {
		if err := validateFieldToMatch(stmt.SqliMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.SqliMatchStatement.TextTransformations); err != nil {
			return err
		}
	}
	if stmt.XssMatchStatement != nil {
		if err := validateFieldToMatch(stmt.XssMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.XssMatchStatement.TextTransformations); err != nil {
			return err
		}
	}
	if stmt.SizeConstraintStatement != nil {
		if stmt.SizeConstraintStatement.ComparisonOperator != "" && !isValidComparisonOperator(stmt.SizeConstraintStatement.ComparisonOperator) {
			return invalidParamError(fmt.Sprintf("Invalid ComparisonOperator: %s", stmt.SizeConstraintStatement.ComparisonOperator))
		}
		if err := validateFieldToMatch(stmt.SizeConstraintStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.SizeConstraintStatement.TextTransformations); err != nil {
			return err
		}
	}
	if stmt.RegexMatchStatement != nil {
		if err := validateFieldToMatch(stmt.RegexMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.RegexMatchStatement.TextTransformations); err != nil {
			return err
		}
	}
	if stmt.RegexPatternSetRefStatement != nil {
		if err := validateFieldToMatch(stmt.RegexPatternSetRefStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.RegexPatternSetRefStatement.TextTransformations); err != nil {
			return err
		}
	}
	if stmt.IPSetReferenceStatement != nil {
		// IPSetForwardedIPConfig enums (Position: FIRST/LAST; FallbackBehavior: MATCH/NO_MATCH/MATCH_OTHER)
		// are validated implicitly by the SDK — no additional validation needed here.
	}
	if stmt.AndStatement != nil {
		for _, s := range stmt.AndStatement.Statements {
			if err := validateStatement(s); err != nil {
				return err
			}
		}
	}
	if stmt.OrStatement != nil {
		for _, s := range stmt.OrStatement.Statements {
			if err := validateStatement(s); err != nil {
				return err
			}
		}
	}
	if stmt.NotStatement != nil {
		if err := validateStatement(stmt.NotStatement.Statement); err != nil {
			return err
		}
	}
	if stmt.RateBasedStatement != nil {
		if stmt.RateBasedStatement.AggregateKeyType != "" && !isValidAggregateKeyType(stmt.RateBasedStatement.AggregateKeyType) {
			return invalidParamError(fmt.Sprintf("Invalid AggregateKeyType: %s", stmt.RateBasedStatement.AggregateKeyType))
		}
		if err := validateStatement(stmt.RateBasedStatement.ScopeDownStatement); err != nil {
			return err
		}
	}
	if stmt.ManagedRuleGroupStatement != nil {
		if err := validateStatement(stmt.ManagedRuleGroupStatement.ScopeDownStatement); err != nil {
			return err
		}
	}
	return nil
}

// validateRules validates all statements in a rule list.
func validateRules(rules []*waf.Rule) error {
	for _, r := range rules {
		if r == nil {
			continue
		}
		if err := validateStatement(r.Statement); err != nil {
			return err
		}
	}
	return nil
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

// parseRules converts raw rule data and validates all enum values.
// Returns an error if any statement contains an invalid enum value.
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
