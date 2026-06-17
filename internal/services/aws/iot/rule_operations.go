package iot

import (
	"context"
	"log/slog"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// actionsToList converts a flat action map (keyed by type) into a list of
// single-key maps, matching the AWS IoT rule action representation.
func actionsToList(m map[string]interface{}) []map[string]interface{} {
	if len(m) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(m))
	for k, v := range m {
		result = append(result, map[string]interface{}{k: v})
	}
	return result
}

func (s *IoTService) CreateTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamCaseInsensitive(req.Parameters, "ruleName")
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "topicRulePayload")

	rule := &iotstore.TopicRule{
		RuleName:         ruleName,
		SQL:              request.GetParamCaseInsensitive(props, "sql"),
		Description:      request.GetParamCaseInsensitive(props, "description"),
		AwsIotSqlVersion: request.GetParamCaseInsensitive(props, "awsIotSqlVersion"),
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
	}
	rule.RuleDisabled = request.GetBoolParam(props, "ruleDisabled")

	// Extract action configurations from the topicRulePayload.
	rule.Actions = extractActionsFromProps(props)
	if ea := request.GetMapParamCaseInsensitive(props, "errorAction"); ea != nil {
		rule.ErrorAction = ea
	}

	created, err := store.CreateRule(rule)
	if err != nil {
		return nil, err
	}

	if !created.RuleDisabled && s.executor != nil && len(created.Actions) > 0 {
		if err := s.executor.AddRule(created.RuleName, created.TopicPattern, created.SQL, actionsToList(created.Actions)); err != nil {
			slog.Warn("rule created but executor registration failed", "rule", created.RuleName, "error", err)
		}
	}

	return map[string]interface{}{
		"ruleArn":  created.ARN,
		"ruleName": created.RuleName,
	}, nil
}

func (s *IoTService) DescribeTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamCaseInsensitive(req.Parameters, "ruleName")
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rule, err := store.GetRule(ruleName)
	if err != nil {
		return nil, iotstore.ErrRuleNotFound
	}

	return ruleToResponse(rule), nil
}

func (s *IoTService) ReplaceTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamCaseInsensitive(req.Parameters, "ruleName")
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "topicRulePayload")
	disabled := request.GetBoolParam(props, "ruleDisabled")
	actionsMap := extractActionsFromProps(props)
	errorAction := request.GetMapParamCaseInsensitive(props, "errorAction")

	opts := iotstore.RuleUpdateOpts{
		SQL:              request.GetParamCaseInsensitive(props, "sql"),
		Description:      request.GetParamCaseInsensitive(props, "description"),
		AwsIotSqlVersion: request.GetParamCaseInsensitive(props, "awsIotSqlVersion"),
		RuleDisabled:     &disabled,
		Actions:          actionsMap,
		ErrorAction:      errorAction,
	}

	updated, err := store.UpdateRule(ruleName, opts)
	if err != nil {
		return nil, err
	}

	if s.executor != nil {
		s.executor.RemoveRule(ruleName)
		if !updated.RuleDisabled && updated.SQL != "" && len(updated.Actions) > 0 {
			if err := s.executor.AddRule(updated.RuleName, updated.TopicPattern, updated.SQL, actionsToList(updated.Actions)); err != nil {
				slog.Warn("rule replaced but executor registration failed", "rule", ruleName, "error", err)
			}
		}
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamCaseInsensitive(req.Parameters, "ruleName")
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteRule(ruleName); err != nil {
		return nil, err
	}

	if s.executor != nil {
		s.executor.RemoveRule(ruleName)
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListTopicRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rules, err := store.ListRules(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(rules.Items))
	for _, r := range rules.Items {
		items = append(items, ruleToResponse(r))
	}

	return listResponse("rules", items, rules.NextMarker), nil
}

func (s *IoTService) EnableTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamCaseInsensitive(req.Parameters, "ruleName")
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	disabled := false
	opts := iotstore.RuleUpdateOpts{RuleDisabled: &disabled}
	updated, err := store.UpdateRule(ruleName, opts)
	if err != nil {
		return nil, err
	}

	if s.executor != nil && updated.SQL != "" && len(updated.Actions) > 0 {
		if err := s.executor.AddRule(updated.RuleName, updated.TopicPattern, updated.SQL, actionsToList(updated.Actions)); err != nil {
			slog.Warn("rule enabled but executor registration failed", "rule", ruleName, "error", err)
		}
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DisableTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamCaseInsensitive(req.Parameters, "ruleName")
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	disabled := true
	opts := iotstore.RuleUpdateOpts{RuleDisabled: &disabled}
	_, err = store.UpdateRule(ruleName, opts)
	if err != nil {
		return nil, err
	}

	if s.executor != nil {
		s.executor.RemoveRule(ruleName)
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) GetTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeTopicRule(ctx, reqCtx, req)
}

func ruleToResponse(r *iotstore.TopicRule) map[string]interface{} {
	resp := map[string]interface{}{
		"ruleArn":          r.ARN,
		"ruleName":         r.RuleName,
		"sql":              r.SQL,
		"topicPattern":     r.TopicPattern,
		"description":      r.Description,
		"ruleDisabled":     r.RuleDisabled,
		"createdAt":        r.CreatedAt,
		"awsIotSqlVersion": r.AwsIotSqlVersion,
	}
	if len(r.Actions) > 0 {
		resp["actions"] = r.Actions
	}
	if len(r.ErrorAction) > 0 {
		resp["errorAction"] = r.ErrorAction
	}
	return resp
}

// extractActionsFromProps extracts action configurations from a topicRulePayload.
// AWS IoT sends actions as a list of single-key maps; this converts to a flat map.
func extractActionsFromProps(props map[string]interface{}) map[string]interface{} {
	if actionsList := request.GetListParamLowerFirst(props, "actions"); len(actionsList) > 0 {
		actionsMap := make(map[string]interface{})
		for _, item := range actionsList {
			for k, v := range item {
				actionsMap[k] = v
			}
		}
		return actionsMap
	}
	return request.GetMapParamCaseInsensitive(props, "actions")
}
