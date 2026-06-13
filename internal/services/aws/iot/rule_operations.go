package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamCaseInsensitive(req.Parameters, "ruleName")
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetRule(ruleName); err == nil {
		return nil, iotstore.ErrResourceAlreadyExists
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

	created, err := store.CreateRule(rule)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
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

	existing, err := store.GetRule(ruleName)
	if err != nil {
		return nil, iotstore.ErrRuleNotFound
	}

	props := unwrapProps(req.Parameters, "topicRulePayload")

	existing.SQL = request.GetParamCaseInsensitive(props, "sql")
	existing.Description = request.GetParamCaseInsensitive(props, "description")
	if ver := request.GetParamCaseInsensitive(props, "awsIotSqlVersion"); ver != "" {
		existing.AwsIotSqlVersion = ver
	}
	existing.RuleDisabled = request.GetBoolParam(props, "ruleDisabled")

	if err := store.UpdateRule(existing); err != nil {
		return nil, iotstore.ErrInvalidRequest
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
		return nil, iotstore.ErrRuleNotFound
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

	resp := map[string]interface{}{
		"rules": items,
	}
	if rules.NextMarker != "" {
		resp["nextToken"] = rules.NextMarker
	}
	return resp, nil
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

	rule, err := store.GetRule(ruleName)
	if err != nil {
		return nil, iotstore.ErrRuleNotFound
	}

	rule.RuleDisabled = false
	if err := store.UpdateRule(rule); err != nil {
		return nil, iotstore.ErrInvalidRequest
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

	rule, err := store.GetRule(ruleName)
	if err != nil {
		return nil, iotstore.ErrRuleNotFound
	}

	rule.RuleDisabled = true
	if err := store.UpdateRule(rule); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) GetTopicRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeTopicRule(ctx, reqCtx, req)
}

func ruleToResponse(r *iotstore.TopicRule) map[string]interface{} {
	return map[string]interface{}{
		"ruleArn":          r.ARN,
		"ruleName":         r.RuleName,
		"sql":              r.SQL,
		"topicPattern":     r.TopicPattern,
		"description":      r.Description,
		"ruleDisabled":     r.RuleDisabled,
		"createdAt":        r.CreatedAt,
		"awsIotSqlVersion": r.AwsIotSqlVersion,
	}
}
