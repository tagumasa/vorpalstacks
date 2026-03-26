// Package waf provides WAF (Web Application Firewall) service operations for vorpalstacks.
package waf

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/store/aws/waf"
)

func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("crypto/rand failed: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func parseTags(tagsRaw interface{}) []waf.Tag {
	var tags []waf.Tag
	switch v := tagsRaw.(type) {
	case []interface{}:
		for _, t := range v {
			if m, ok := t.(map[string]interface{}); ok {
				tags = append(tags, waf.Tag{
					Key:   request.GetStringParam(m, "Key"),
					Value: request.GetStringParam(m, "Value"),
				})
			}
		}
	}
	return tags
}

func convertPBToStoreVisibilityConfig(m map[string]interface{}) *waf.VisibilityConfig {
	if m == nil {
		return &waf.VisibilityConfig{
			SampledRequestsEnabled:   true,
			CloudWatchMetricsEnabled: true,
			MetricName:               "unnamed",
		}
	}
	return &waf.VisibilityConfig{
		SampledRequestsEnabled:   request.GetBoolParam(m, "SampledRequestsEnabled"),
		CloudWatchMetricsEnabled: request.GetBoolParam(m, "CloudWatchMetricsEnabled"),
		MetricName:               request.GetStringParam(m, "MetricName"),
	}
}

func convertStoreVisibilityConfigToPB(vc *waf.VisibilityConfig) map[string]interface{} {
	if vc == nil {
		return nil
	}
	return map[string]interface{}{
		"SampledRequestsEnabled":   vc.SampledRequestsEnabled,
		"CloudWatchMetricsEnabled": vc.CloudWatchMetricsEnabled,
		"MetricName":               vc.MetricName,
	}
}

func convertPBToStoreAction(m map[string]interface{}) *waf.Action {
	if m == nil {
		return &waf.Action{}
	}
	action := &waf.Action{}
	actionType := request.GetStringParam(m, "ActionType")
	if actionType == "" {
		actionType = request.GetStringParam(m, "Type")
	}
	switch actionType {
	case "ALLOW":
		action.Allow = &waf.AllowAction{}
	case "BLOCK":
		action.Block = &waf.BlockAction{}
	case "COUNT":
		action.Count = &waf.CountAction{}
	case "CAPTCHA":
		action.Captcha = &waf.CaptchaAction{}
	case "CHALLENGE":
		action.Challenge = &waf.ChallengeAction{}
	}
	if _, ok := m["Allow"]; ok {
		action.Allow = &waf.AllowAction{}
	}
	if _, ok := m["Block"]; ok {
		action.Block = &waf.BlockAction{}
	}
	if _, ok := m["Count"]; ok {
		action.Count = &waf.CountAction{}
	}
	if _, ok := m["Captcha"]; ok {
		action.Captcha = &waf.CaptchaAction{}
	}
	if _, ok := m["Challenge"]; ok {
		action.Challenge = &waf.ChallengeAction{}
	}
	return action
}

func convertStoreActionToPB(a *waf.Action) map[string]interface{} {
	if a == nil {
		return nil
	}
	result := map[string]interface{}{}
	if a.Allow != nil {
		result["Allow"] = map[string]interface{}{}
	} else if a.Block != nil {
		result["Block"] = map[string]interface{}{}
	} else if a.Count != nil {
		result["Count"] = map[string]interface{}{}
	} else if a.Captcha != nil {
		result["Captcha"] = map[string]interface{}{}
	} else if a.Challenge != nil {
		result["Challenge"] = map[string]interface{}{}
	}
	return result
}

func convertStoreRulesToPB(rules []*waf.Rule) []interface{} {
	if rules == nil {
		return nil
	}
	result := make([]interface{}, len(rules))
	for i, r := range rules {
		result[i] = map[string]interface{}{
			"Name":      r.Name,
			"Priority":  r.Priority,
			"Statement": convertStoreStatementToPB(r.Statement),
			"Action": func() interface{} {
				if a, ok := r.Action.(*waf.Action); ok {
					return convertStoreActionToPB(a)
				}
				return r.Action
			}(),
			"OverrideAction":   convertStoreActionToPB(r.OverrideAction),
			"VisibilityConfig": convertStoreVisibilityConfigToPB(r.VisibilityConfig),
		}
	}
	return result
}

func convertPBToStoreRules(rulesRaw interface{}) []*waf.Rule {
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
				rule.Action = convertPBToStoreAction(actionMap)
			}
		}
		if overrideActionRaw, ok := rMap["OverrideAction"]; ok {
			if overrideActionMap, ok := overrideActionRaw.(map[string]interface{}); ok {
				rule.OverrideAction = convertPBToStoreAction(overrideActionMap)
			}
		}
		if vcRaw, ok := rMap["VisibilityConfig"]; ok {
			if vcMap, ok := vcRaw.(map[string]interface{}); ok {
				rule.VisibilityConfig = convertPBToStoreVisibilityConfig(vcMap)
			}
		}
		if statementRaw, ok := rMap["Statement"]; ok {
			if statementMap, ok := statementRaw.(map[string]interface{}); ok {
				rule.Statement = convertPBToStoreStatement(statementMap)
			}
		}
		rules = append(rules, rule)
	}
	return rules
}

func convertPBToStoreStatement(m map[string]interface{}) *waf.Statement {
	if m == nil {
		return nil
	}
	stmt := &waf.Statement{}
	if _, ok := m["IPSetReferenceStatement"]; ok {
		if ipSetMap, ok := m["IPSetReferenceStatement"].(map[string]interface{}); ok {
			stmt.IPSetReferenceStatement = &waf.IPSetReferenceStatement{
				IPSetArn: request.GetStringParam(ipSetMap, "ARN"),
			}
		}
	}
	if _, ok := m["RegexPatternSetReferenceStatement"]; ok {
		stmt.RegexMatchStatement = &waf.RegexMatchStatement{}
	}
	if _, ok := m["ByteMatchStatement"]; ok {
		stmt.ByteMatchStatement = &waf.ByteMatchStatement{}
	}
	if _, ok := m["GeoMatchStatement"]; ok {
		stmt.GeoMatchStatement = &waf.GeoMatchStatement{}
	}
	if _, ok := m["RateBasedStatement"]; ok {
		stmt.RateBasedStatement = &waf.RateBasedStatement{}
	}
	if _, ok := m["SizeConstraintStatement"]; ok {
		stmt.SizeConstraintStatement = &waf.SizeConstraintStatement{}
	}
	if _, ok := m["ManagedRuleGroupStatement"]; ok {
		stmt.ManagedRuleGroupStatement = &waf.ManagedRuleGroupStatement{}
	}
	return stmt
}

func convertStoreStatementToPB(s *waf.Statement) map[string]interface{} {
	if s == nil {
		return nil
	}
	result := map[string]interface{}{}
	if s.IPSetReferenceStatement != nil {
		result["IPSetReferenceStatement"] = map[string]interface{}{
			"ARN": s.IPSetReferenceStatement.IPSetArn,
		}
	}
	return result
}
