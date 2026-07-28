package iot

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/iot/rules"
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
		CreatedAt:        fmt.Sprintf("%d", time.Now().UTC().Unix()),
	}
	rule.RuleDisabled = request.GetBoolParam(props, "ruleDisabled")

	// Validate the IoT SQL statement before persisting.
	if rule.SQL != "" {
		if _, err := rules.NewParser(rule.SQL).Parse(); err != nil {
			return nil, iotstore.ErrSqlParse
		}
	}

	// Extract action configurations from the topicRulePayload.
	rule.Actions = extractActionsFromProps(props)
	if ea := request.GetMapParamCaseInsensitive(props, "errorAction"); ea != nil {
		rule.ErrorAction = ea
	}

	created, err := store.CreateRule(rule)
	if err != nil {
		return nil, err
	}

	if executor := s.executorForReq(reqCtx); !created.RuleDisabled && executor != nil && len(created.Actions) > 0 {
		if err := executor.AddRule(created.RuleName, created.TopicPattern, created.SQL, actionsToList(created.Actions), created.ErrorAction); err != nil {
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
		return nil, err
	}

	// AWS wraps the rule in a "rule" member (DescribeTopicRuleOutput.rule /
	// GetTopicRuleOutput.rule); returning it flat leaves the SDK Rule field nil.
	return map[string]interface{}{"rule": ruleToResponse(rule)}, nil
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

	// Validate the IoT SQL statement before persisting.
	if opts.SQL != "" {
		if _, err := rules.NewParser(opts.SQL).Parse(); err != nil {
			return nil, iotstore.ErrSqlParse
		}
	}

	updated, err := store.UpdateRule(ruleName, opts)
	if err != nil {
		return nil, err
	}

	if executor := s.executorForReq(reqCtx); executor != nil {
		executor.RemoveRule(ruleName)
		if !updated.RuleDisabled && updated.SQL != "" && len(updated.Actions) > 0 {
			if err := executor.AddRule(updated.RuleName, updated.TopicPattern, updated.SQL, actionsToList(updated.Actions), updated.ErrorAction); err != nil {
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

	arn := iotstore.BuildRuleARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), ruleName)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteRule(ruleName); err != nil {
		return nil, err
	}

	if executor := s.executorForReq(reqCtx); executor != nil {
		executor.RemoveRule(ruleName)
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

	if executor := s.executorForReq(reqCtx); executor != nil && updated.SQL != "" && len(updated.Actions) > 0 {
		if err := executor.AddRule(updated.RuleName, updated.TopicPattern, updated.SQL, actionsToList(updated.Actions), updated.ErrorAction); err != nil {
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

	if executor := s.executorForReq(reqCtx); executor != nil {
		executor.RemoveRule(ruleName)
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) GetTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeTopicRule(ctx, reqCtx, req)
}

func ruleToResponse(r *iotstore.TopicRule) map[string]interface{} {
	// Smithy: createdAt is CreatedAtDate (timestamp). SDK expects JSON Number.
	// Handle both Unix int string (new) and RFC3339 string (legacy).
	createdAt, _ := strconv.ParseInt(r.CreatedAt, 10, 64)
	if createdAt == 0 && r.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, r.CreatedAt); err == nil {
			createdAt = t.Unix()
		}
	}
	// TopicRule (output shape) includes topicPattern; TopicRulePayload (input)
	// does not. Derive from the SQL FROM clause when not already stored.
	topicPattern := r.TopicPattern
	if topicPattern == "" && r.SQL != "" {
		if parsed, err := rules.NewParser(r.SQL).Parse(); err == nil {
			topicPattern = parsed.FromTopic
		}
	}
	resp := map[string]interface{}{
		"ruleArn":          r.ARN,
		"ruleName":         r.RuleName,
		"sql":              r.SQL,
		"topicPattern":     topicPattern,
		"description":      r.Description,
		"ruleDisabled":     r.RuleDisabled,
		"createdAt":        createdAt,
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
