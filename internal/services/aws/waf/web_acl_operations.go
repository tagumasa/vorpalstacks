package waf

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateWebACL creates a new web ACL with the specified default action and rules.
func (s *WAFService) CreateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidParameterValue", "Name is required", 400)
	}
	metricName := request.GetStringParam(req.Parameters, "MetricName")
	if metricName == "" {
		return nil, NewAPIError("InvalidParameterValue", "MetricName is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	var defaultAction interface{}
	if daRaw := req.Parameters["DefaultAction"]; daRaw != nil {
		if da, ok := daRaw.(map[string]interface{}); ok {
			defaultAction = da
		}
	}

	var rules []interface{}
	if rulesRaw := req.Parameters["Rules"]; rulesRaw != nil {
		if arr, ok := rulesRaw.([]interface{}); ok {
			rules = arr
		}
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	webACL := &wafstore.WebACL{
		ID:            id,
		Name:          name,
		MetricName:    metricName,
		DefaultAction: defaultAction,
		ARN:           stores.arnBuilder.BuildWebACLARN(id, ""),
	}

	_, err = stores.webACLs.Get(id)
	if err == nil {
		return nil, NewAPIError("AlreadyExistsException", fmt.Sprintf("WebACL already exists: %s", name), 400)
	}

	if tags := parseTags(req.Parameters["Tags"]); len(tags) > 0 {
		webACL.Tags = tags
	}

	if len(rules) > 0 {
		webACL.Rules = convertToStoreRules(rules)
	}

	if err := stores.webACLs.Put(webACL.ID, webACL); err != nil {
		logs.Warn("failed to persist WebACL", logs.String("id", webACL.ID), logs.Err(err))
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
		"WebACL": map[string]interface{}{
			"WebACLId":      webACL.ID,
			"Name":          webACL.Name,
			"MetricName":    webACL.MetricName,
			"DefaultAction": webACL.DefaultAction,
			"Rules":         convertRulesToResponse(webACL.Rules),
			"WebACLArn":     webACL.ARN,
		},
	}, nil
}

// GetWebACL returns the details of the specified web ACL.
func (s *WAFService) GetWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "WebACLId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "WebACLId is required", 400)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	webACL, err := stores.webACLs.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "WebACL not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"WebACL": map[string]interface{}{
			"WebACLId":      webACL.ID,
			"Name":          webACL.Name,
			"MetricName":    webACL.MetricName,
			"DefaultAction": webACL.DefaultAction,
			"Rules":         convertRulesToResponse(webACL.Rules),
			"WebACLArn":     webACL.ARN,
		},
	}, nil
}

// ListWebACLs returns a paginated list of all web ACLs.
func (s *WAFService) ListWebACLs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := request.GetIntParam(req.Parameters, "Limit")
	if maxItems == 0 {
		maxItems = 100
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.webACLs.List("", maxItems)
	if err != nil {
		return nil, err
	}

	webACLs := make([]interface{}, 0, len(result.WebACLs))
	for _, wa := range result.WebACLs {
		webACLs = append(webACLs, map[string]interface{}{
			"WebACLId": wa.ID,
			"Name":     wa.Name,
		})
	}

	return map[string]interface{}{
		"WebACLs":    webACLs,
		"NextMarker": result.NextMarker,
	}, nil
}

// UpdateWebACL inserts or deletes activated rules and updates the default action of the specified web ACL.
func (s *WAFService) UpdateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "WebACLId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "WebACLId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	webACL, err := stores.webACLs.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "WebACL not found", 400)
		}
		return nil, err
	}

	if updatesRaw := req.Parameters["Updates"]; updatesRaw != nil {
		if updates, ok := updatesRaw.([]interface{}); ok {
			for _, u := range updates {
				if update, ok := u.(map[string]interface{}); ok {
					action := request.GetStringParam(update, "Action")
					if action == "INSERT" {
						if activatedRule, ok := update["ActivatedRule"].(map[string]interface{}); ok {
							webACL.Rules = append(webACL.Rules, convertMapToRule(activatedRule))
						}
					} else if action == "DELETE" {
						if activatedRule, ok := update["ActivatedRule"].(map[string]interface{}); ok {
							webACL.Rules = removeStoreRule(webACL.Rules, activatedRule)
						}
					}
				}
			}
		}
	}

	if daRaw := req.Parameters["DefaultAction"]; daRaw != nil {
		if da, ok := daRaw.(map[string]interface{}); ok {
			webACL.DefaultAction = da
		}
	}

	if err := stores.webACLs.Put(id, webACL); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}

func removeActivatedRule(rules []interface{}, target map[string]interface{}) []interface{} {
	tRuleId, _ := target["RuleId"].(string)
	result := make([]interface{}, 0, len(rules))
	for _, r := range rules {
		if rm, ok := r.(map[string]interface{}); ok {
			if rm["RuleId"] == tRuleId {
				continue
			}
		}
		result = append(result, r)
	}
	return result
}

func removeStoreRule(rules []*wafstore.Rule, target map[string]interface{}) []*wafstore.Rule {
	tRuleId, _ := target["RuleId"].(string)
	result := make([]*wafstore.Rule, 0, len(rules))
	for _, r := range rules {
		if r.ID == tRuleId {
			continue
		}
		result = append(result, r)
	}
	return result
}

func convertToStoreRules(rules []interface{}) []*wafstore.Rule {
	result := make([]*wafstore.Rule, 0, len(rules))
	for _, r := range rules {
		if rm, ok := r.(map[string]interface{}); ok {
			result = append(result, convertMapToRule(rm))
		}
	}
	return result
}

func convertMapToRule(m map[string]interface{}) *wafstore.Rule {
	rule := &wafstore.Rule{}
	if v, ok := m["RuleId"].(string); ok {
		rule.ID = v
	}
	if v, ok := m["Name"].(string); ok {
		rule.Name = v
	}
	if v, ok := m["MetricName"].(string); ok {
		rule.MetricName = v
	}
	if v, ok := m["Priority"].(float64); ok {
		rule.Priority = int32(v)
	}
	if v, ok := m["Priority"].(int32); ok {
		rule.Priority = v
	}
	if _, ok := m["Type"].(string); ok {
		if action, ok := m["Action"].(map[string]interface{}); ok {
			rule.Action = action
		}
	}
	if preds, ok := m["Predicates"].([]interface{}); ok {
		rule.Predicates = preds
	}
	return rule
}

func convertRulesToResponse(rules []*wafstore.Rule) []interface{} {
	if len(rules) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(rules))
	for _, r := range rules {
		m := map[string]interface{}{
			"RuleId":     r.ID,
			"Name":       r.Name,
			"MetricName": r.MetricName,
			"Priority":   r.Priority,
			"Type":       "REGULAR",
		}
		if r.Action != nil {
			m["Action"] = r.Action
		}
		if len(r.Predicates) > 0 {
			m["Predicates"] = r.Predicates
		}
		result = append(result, m)
	}
	return result
}

// DeleteWebACL permanently deletes the specified web ACL.
func (s *WAFService) DeleteWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "WebACLId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "WebACLId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.webACLs.Delete(id, ""); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "WebACL not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}

// CheckCapacity returns the capacity consumed by the specified web ACL rules.
func (s *WAFService) CheckCapacity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"Capacity": 10,
	}, nil
}
