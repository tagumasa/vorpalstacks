package wafv2

import (
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/store/aws/waf"
)

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

func convertVisibilityConfig(m map[string]interface{}) *waf.VisibilityConfig {
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

func convertVisibilityConfigToResponse(vc *waf.VisibilityConfig) map[string]interface{} {
	if vc == nil {
		return nil
	}
	return map[string]interface{}{
		"SampledRequestsEnabled":   vc.SampledRequestsEnabled,
		"CloudWatchMetricsEnabled": vc.CloudWatchMetricsEnabled,
		"MetricName":               vc.MetricName,
	}
}

func convertAction(m map[string]interface{}) *waf.Action {
	if m == nil {
		return &waf.Action{}
	}
	action := &waf.Action{}
	if _, ok := m["Allow"]; ok {
		action.Allow = &waf.AllowAction{}
	}
	if _, ok := m["Block"]; ok {
		action.Block = &waf.BlockAction{}
	}
	if _, ok := m["Count"]; ok {
		action.Count = &waf.CountAction{}
	}
	return action
}

func convertActionToResponse(a interface{}) map[string]interface{} {
	if a == nil {
		return map[string]interface{}{"Allow": map[string]interface{}{}}
	}
	if action, ok := a.(*waf.Action); ok {
		result := map[string]interface{}{}
		if action.Allow != nil {
			result["Allow"] = map[string]interface{}{}
		} else if action.Block != nil {
			result["Block"] = map[string]interface{}{}
		} else if action.Count != nil {
			result["Count"] = map[string]interface{}{}
		}
		if len(result) == 0 {
			result["Allow"] = map[string]interface{}{}
		}
		return result
	}
	if m, ok := a.(map[string]interface{}); ok {
		result := map[string]interface{}{}
		for k, v := range m {
			switch k {
			case "allow", "Allow":
				result["Allow"] = v
			case "block", "Block":
				result["Block"] = v
			case "count", "Count":
				result["Count"] = v
			case "captcha", "Captcha":
				result["Captcha"] = v
			case "challenge", "Challenge":
				result["Challenge"] = v
			default:
				result[k] = v
			}
		}
		if len(result) == 0 {
			result["Allow"] = map[string]interface{}{}
		}
		return result
	}
	return map[string]interface{}{"Allow": map[string]interface{}{}}
}

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
		rules = append(rules, rule)
	}
	return rules
}

func convertStatement(m map[string]interface{}) *waf.Statement {
	if m == nil {
		return nil
	}
	stmt := &waf.Statement{}
	if _, ok := m["IPSetReferenceStatement"]; ok {
		if ipSetMap, ok := m["IPSetReferenceStatement"].(map[string]interface{}); ok {
			arn := request.GetStringParam(ipSetMap, "ARN")
			stmt.IPSetReferenceStatement = &waf.IPSetReferenceStatement{
				IPSetArn: arn,
			}
			if ipConfig, ok := ipSetMap["IPSetForwardedIPConfig"].(map[string]interface{}); ok {
				stmt.IPSetReferenceStatement.IPSetForwardedIPConfig = &waf.IPSetForwardedIPConfig{
					HeaderName:       request.GetStringParam(ipConfig, "HeaderName"),
					Position:         request.GetStringParam(ipConfig, "Position"),
					FallbackBehavior: request.GetStringParam(ipConfig, "FallbackBehavior"),
				}
			}
		}
	}
	if _, ok := m["RegexPatternSetReferenceStatement"]; ok {
		if rpsMap, ok := m["RegexPatternSetReferenceStatement"].(map[string]interface{}); ok {
			stmt.RegexMatchStatement = &waf.RegexMatchStatement{
				RegexPatternSetID: request.GetStringParam(rpsMap, "ARN"),
			}
		}
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
	if _, ok := m["ManagedRuleGroupStatement"]; ok {
		if mrgMap, ok := m["ManagedRuleGroupStatement"].(map[string]interface{}); ok {
			stmt.ManagedRuleGroupStatement = &waf.ManagedRuleGroupStatement{
				Name:       request.GetStringParam(mrgMap, "Name"),
				VendorName: request.GetStringParam(mrgMap, "VendorName"),
			}
		}
	}
	return stmt
}

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
			"Name":      r.Name,
			"Priority":  r.Priority,
			"Statement": convertStatementToResponse(r.Statement),
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
		result = append(result, m)
	}
	return result
}

func convertStatementToResponse(s *waf.Statement) map[string]interface{} {
	if s == nil {
		return nil
	}
	result := map[string]interface{}{}
	if s.IPSetReferenceStatement != nil {
		ipSetRef := map[string]interface{}{
			"ARN": s.IPSetReferenceStatement.IPSetArn,
		}
		if s.IPSetReferenceStatement.IPSetForwardedIPConfig != nil {
			ipSetRef["IPSetForwardedIPConfig"] = map[string]interface{}{
				"HeaderName":       s.IPSetReferenceStatement.IPSetForwardedIPConfig.HeaderName,
				"Position":         s.IPSetReferenceStatement.IPSetForwardedIPConfig.Position,
				"FallbackBehavior": s.IPSetReferenceStatement.IPSetForwardedIPConfig.FallbackBehavior,
			}
		}
		result["IPSetReferenceStatement"] = ipSetRef
	}
	if s.RegexMatchStatement != nil {
		result["RegexPatternSetReferenceStatement"] = map[string]interface{}{
			"ARN": s.RegexMatchStatement.RegexPatternSetID,
		}
	}
	if s.ManagedRuleGroupStatement != nil {
		result["ManagedRuleGroupStatement"] = map[string]interface{}{
			"Name":       s.ManagedRuleGroupStatement.Name,
			"VendorName": s.ManagedRuleGroupStatement.VendorName,
		}
	}
	return result
}

func mapToAction(m map[string]interface{}) *waf.Action {
	action := &waf.Action{}
	for k := range m {
		switch k {
		case "allow", "Allow":
			action.Allow = &waf.AllowAction{}
		case "block", "Block":
			action.Block = &waf.BlockAction{}
		case "count", "Count":
			action.Count = &waf.CountAction{}
		case "captcha", "Captcha":
			action.Captcha = &waf.CaptchaAction{}
		case "challenge", "Challenge":
			action.Challenge = &waf.ChallengeAction{}
		}
	}
	return action
}
