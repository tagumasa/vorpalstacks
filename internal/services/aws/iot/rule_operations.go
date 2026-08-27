package iot

import (
	"context"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/iot/rules"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	props := unwrapProps(req.Parameters, "topicRulePayload")
	in := TopicRuleInput{
		RuleName:         request.GetParamCaseInsensitive(req.Parameters, "ruleName"),
		SQL:              request.GetParamCaseInsensitive(props, "sql"),
		Description:      request.GetParamCaseInsensitive(props, "description"),
		AwsIotSqlVersion: request.GetParamCaseInsensitive(props, "awsIotSqlVersion"),
		RuleDisabled:     request.GetBoolParam(props, "ruleDisabled"),
		Actions:          extractActionsFromProps(props),
		ErrorAction:      request.GetMapParamCaseInsensitive(props, "errorAction"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createTopicRuleCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ruleArn":  result.RuleArn,
		"ruleName": result.RuleName,
	}, nil
}

// GetTopicRule retrieves a rule definition. The model output shape carries
// both the wrapped rule member and the ruleArn.
func (s *IoTService) GetTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rule, err := s.describeTopicRuleCore(store, request.GetParamCaseInsensitive(req.Parameters, "ruleName"))
	if err != nil {
		return nil, err
	}

	// AWS wraps the rule in a "rule" member (GetTopicRuleOutput.rule);
	// returning it flat leaves the SDK Rule field nil. ruleArn is a sibling
	// output member (GetTopicRuleOutput.ruleArn).
	return map[string]interface{}{"rule": ruleToResponse(rule), "ruleArn": rule.ARN}, nil
}

func (s *IoTService) ReplaceTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	props := unwrapProps(req.Parameters, "topicRulePayload")
	in := TopicRuleInput{
		RuleName:         request.GetParamCaseInsensitive(req.Parameters, "ruleName"),
		SQL:              request.GetParamCaseInsensitive(props, "sql"),
		Description:      request.GetParamCaseInsensitive(props, "description"),
		AwsIotSqlVersion: request.GetParamCaseInsensitive(props, "awsIotSqlVersion"),
		RuleDisabled:     request.GetBoolParam(props, "ruleDisabled"),
		Actions:          extractActionsFromProps(props),
		ErrorAction:      request.GetMapParamCaseInsensitive(props, "errorAction"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.replaceTopicRuleCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteTopicRuleCore(store, request.GetParamCaseInsensitive(req.Parameters, "ruleName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListTopicRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	nextToken := request.GetParamCaseInsensitive(req.Parameters, "nextToken")
	if nextToken == "" {
		nextToken = request.GetParamCaseInsensitive(req.Parameters, "marker")
	}
	maxResults := request.GetIntParam(req.Parameters, "maxResults")
	if maxResults <= 0 {
		maxResults = request.GetIntParam(req.Parameters, "pageSize")
	}

	result, err := s.listTopicRulesCore(store, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Rules))
	for _, r := range result.Rules {
		items = append(items, ruleToResponse(r))
	}

	return listResponse("rules", items, result.NextToken), nil
}

func (s *IoTService) EnableTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.enableTopicRuleCore(store, request.GetParamCaseInsensitive(req.Parameters, "ruleName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DisableTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.disableTopicRuleCore(store, request.GetParamCaseInsensitive(req.Parameters, "ruleName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
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
