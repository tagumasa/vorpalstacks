package waf

import (
	"context"
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateRule creates a new WAF rule with the specified name and metric name.
func (s *WAFService) CreateRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidParameterValue", "Name is required", 400)
	}
	metricName := request.GetStringParam(req.Parameters, "MetricName")
	if metricName == "" {
		return nil, NewAPIError("InvalidParameterValue", "MetricName is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	rule := &wafstore.Rule{
		ID:         id,
		Name:       name,
		MetricName: metricName,
	}

	if err := stores.ruleGroups.CreateRule(id, rule); err != nil {
		return nil, NewAPIError("WAFStaleDataException", fmt.Sprintf("Failed to create rule: %v", err), 400)
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
		"Rule": map[string]interface{}{
			"RuleId":     rule.ID,
			"Name":       rule.Name,
			"MetricName": rule.MetricName,
		},
	}, nil
}

// GetRule returns the details of the specified WAF rule.
func (s *WAFService) GetRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "RuleId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "RuleId is required", 400)
	}

	rule, err := stores.ruleGroups.GetRule(id)
	if err != nil {
		return nil, NewAPIError("WAFNonexistentItemException", "Rule not found", 400)
	}

	return map[string]interface{}{
		"Rule": map[string]interface{}{
			"RuleId":     rule.ID,
			"Name":       rule.Name,
			"MetricName": rule.MetricName,
			"Predicates": rule.Predicates,
		},
	}, nil
}

// UpdateRule inserts or deletes predicates within the specified WAF rule.
func (s *WAFService) UpdateRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "RuleId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "RuleId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rule, err := stores.ruleGroups.GetRule(id)
	if err != nil {
		return nil, NewAPIError("WAFNonexistentItemException", "Rule not found", 400)
	}

	if updatesRaw := req.Parameters["Updates"]; updatesRaw != nil {
		if updates, ok := updatesRaw.([]interface{}); ok {
			for _, u := range updates {
				if update, ok := u.(map[string]interface{}); ok {
					action := request.GetStringParam(update, "Action")
					if pred, ok := update["Predicate"].(map[string]interface{}); ok {
						if action == "INSERT" {
							rule.Predicates = append(rule.Predicates, pred)
						} else if action == "DELETE" {
							rule.Predicates = removePredicate(rule.Predicates, pred)
						}
					}
				}
			}
		}
	}

	if err := stores.ruleGroups.CreateRule(id, rule); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}

func removePredicate(predicates []interface{}, target map[string]interface{}) []interface{} {
	tDataId, _ := target["DataId"].(string)
	result := make([]interface{}, 0, len(predicates))
	for _, p := range predicates {
		if pm, ok := p.(map[string]interface{}); ok {
			if pm["DataId"] == tDataId {
				continue
			}
		}
		result = append(result, p)
	}
	return result
}

// DeleteRule permanently deletes the specified WAF rule.
func (s *WAFService) DeleteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "RuleId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "RuleId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	if err := stores.ruleGroups.DeleteRule(id); err != nil {
		return nil, NewAPIError("WAFNonexistentItemException", "Rule not found", 400)
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}
