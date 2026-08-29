package eventbridge

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/scheduleexpr"
	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

func ruleToMap(r *eventsstore.Rule, includeTimestamps bool) map[string]interface{} {
	result := map[string]interface{}{
		"Arn":          r.ARN,
		"Name":         r.Name,
		"EventBusName": r.EventBusName,
		"State":        string(r.State),
	}
	if includeTimestamps {
		result["CreationTime"] = r.CreatedAt.Unix()
		result["LastModifiedTime"] = r.LastModifiedAt.Unix()
	}
	if r.Description != "" {
		result["Description"] = r.Description
	}
	if r.EventPattern != "" {
		result["EventPattern"] = r.EventPattern
	}
	if r.ScheduleExpression != "" {
		result["ScheduleExpression"] = r.ScheduleExpression
	}
	if r.RoleARN != "" {
		result["RoleArn"] = r.RoleARN
	}
	if r.ManagedBy != "" {
		result["ManagedBy"] = r.ManagedBy
	}
	if r.CreatedBy != "" {
		result["CreatedBy"] = r.CreatedBy
	}
	return result
}

func isValidEventPattern(pattern string) bool {
	if pattern == "" {
		return true
	}
	var js map[string]interface{}
	return json.Unmarshal([]byte(pattern), &js) == nil
}

func isValidScheduleExpression(expr string) bool {
	return scheduleexpr.ValidateRuleExpression(expr)
}

// parsePutRuleInput reads the PutRule wire request into the
// transport-agnostic Core input.
func parsePutRuleInput(req *request.ParsedRequest, eventBusName string, eventBusNameProvided bool, createdBy string) PutRuleInput {
	input := PutRuleInput{
		Name:                 request.GetParamLowerFirst(req.Parameters, "Name"),
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		CreatedBy:            createdBy,
	}
	if desc, ok := req.Parameters["Description"].(string); ok {
		input.DescriptionSet = true
		input.Description = desc
	}
	if pattern, ok := req.Parameters["EventPattern"].(string); ok {
		input.EventPatternSet = true
		input.EventPattern = pattern
	}
	if schedule, ok := req.Parameters["ScheduleExpression"].(string); ok {
		input.ScheduleExpressionSet = true
		input.ScheduleExpression = schedule
	}
	if state, ok := req.Parameters["State"].(string); ok {
		input.StateSet = true
		input.State = state
	}
	if roleArn, ok := req.Parameters["RoleArn"].(string); ok {
		input.RoleArnSet = true
		input.RoleArn = roleArn
	}
	input.Tags = tagutil.ParseTags(req.Parameters, "Tags")
	return input
}

// PutRule creates or updates a rule on the specified event bus.
// Supports event patterns and schedule expressions (cron/rate).
func (s *EventsService) PutRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	input := parsePutRuleInput(req, eventBusName, eventBusNameProvided, "arn:aws:iam::"+reqCtx.GetAccountID()+":root")
	input.IAMValidator = reqCtx.GetIAMValidator()

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putRuleCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RuleArn": result.RuleArn,
	}, nil
}

// DeleteRule removes a rule from the event bus.
// Rules with targets cannot be deleted until targets are removed.
func (s *EventsService) DeleteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	force, _ := req.Parameters["Force"].(bool)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteRuleCore(ctx, store, DeleteRuleInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		Name:                 name,
		Force:                force,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeRule returns details about a rule including its state,
// schedule expression, event pattern, and tags.
func (s *EventsService) DescribeRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeRuleCore(ctx, store, DescribeRuleInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		Name:                 name,
	})
	if err != nil {
		return nil, err
	}

	response := ruleToMap(result.Rule, false)
	if len(result.Tags) > 0 {
		response["Tags"] = tagListToMaps(result.Tags)
	}

	return response, nil
}

// ListRules returns a list of rules for the specified event bus,
// optionally filtered by name prefix with pagination support.
func (s *EventsService) ListRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	eventBusName, eventBusNameProvided := eventBusNameParam(req)
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listRulesCore(ctx, store, ListRulesInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		NamePrefix:           namePrefix,
		Limit:                limit,
		NextToken:            nextToken,
	})
	if err != nil {
		return nil, err
	}

	rules := make([]map[string]interface{}, len(result.Rules))
	for i, r := range result.Rules {
		rules[i] = ruleToMap(r, true)
	}

	response := map[string]interface{}{
		"Rules": rules,
	}

	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}

// EnableRule enables a rule so it can match and deliver events to its targets.
func (s *EventsService) EnableRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setRuleStateCore(ctx, store, SetRuleStateInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		Name:                 name,
		State:                eventsstore.RuleStateEnabled,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisableRule disables a rule so it no longer delivers events to its targets.
func (s *EventsService) DisableRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setRuleStateCore(ctx, store, SetRuleStateInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		Name:                 name,
		State:                eventsstore.RuleStateDisabled,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListRuleNamesByTarget returns the list of rules that have the specified target.
// It scans all rules on the event bus, checks each rule's targets for a match,
// then applies client-side pagination to the matched rule names.
func (s *EventsService) ListRuleNamesByTarget(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	input := ListRuleNamesByTargetInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		TargetArn:            request.GetStringParam(req.Parameters, "TargetArn"),
		Limit:                int32(request.GetIntParam(req.Parameters, "Limit")),
		NextToken:            pagination.GetMarker(req.Parameters, "NextToken"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listRuleNamesByTargetCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"RuleNames": result.RuleNames,
	}

	if result.NextToken != "" {
		pagination.SetNextToken(resp, "NextToken", result.NextToken)
	}

	return resp, nil
}
