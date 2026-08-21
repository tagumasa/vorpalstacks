package eventbridge

import (
	"context"
	"encoding/json"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
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

// PutRule creates or updates a rule on the specified event bus.
// Supports event patterns and schedule expressions (cron/rate).
func (s *EventsService) PutRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}
	if !validateResourceName(name, "rule") {
		return nil, awserrors.NewValidationException("Rule name must match the pattern and be 1-64 characters")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Check if event bus exists, auto-create default event bus if needed
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		if err == eventsstore.ErrEventBusNotFound {
			if eventBusName == "default" {
				if createErr := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: eventBusName}); createErr != nil {
					return nil, createErr
				}
			} else {
				return nil, NewResourceNotFoundException("Event bus '" + eventBusName + "' does not exist")
			}
		} else {
			return nil, err
		}
	}

	rule := &eventsstore.Rule{
		Name:         name,
		EventBusName: eventBusName,
		CreatedBy:    "arn:aws:iam::" + reqCtx.GetAccountID() + ":root",
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		if !validateDescription(desc) {
			return nil, errDescriptionTooLong()
		}
		rule.Description = desc
	}

	if pattern, ok := req.Parameters["EventPattern"].(string); ok {
		if !validateEventPatternLength(pattern) {
			return nil, awserrors.NewValidationException("EventPattern must be at most 4096 characters")
		}
		if !isValidEventPattern(pattern) {
			return nil, awserrors.NewInvalidEventPatternException("EventPattern must be valid JSON")
		}
		rule.EventPattern = pattern
	}

	if schedule, ok := req.Parameters["ScheduleExpression"].(string); ok {
		if !isValidScheduleExpression(schedule) {
			return nil, awserrors.NewValidationException("ScheduleExpression must be a valid rate() or cron() expression")
		}
		rule.ScheduleExpression = schedule
	}

	// AWS requires a rule to contain at least an EventPattern or
	// ScheduleExpression.  A rule with neither is rejected with
	// ValidationException.
	if rule.EventPattern == "" && rule.ScheduleExpression == "" {
		return nil, awserrors.NewValidationException("A rule must contain at least an EventPattern or ScheduleExpression")
	}

	if state, ok := req.Parameters["State"].(string); ok {
		if !validateRuleState(state) {
			return nil, awserrors.NewValidationException("State must be 'ENABLED', 'DISABLED', or 'ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS'")
		}
		rule.State = eventsstore.RuleState(state)
	} else {
		rule.State = eventsstore.RuleStateEnabled
	}

	if roleArn, ok := req.Parameters["RoleArn"].(string); ok {
		if roleArn != "" {
			if s.bus != nil {
				if rr := s.bus.RoleResolver(); rr != nil {
					if err := rr.ValidateRole(ctx, roleArn); err != nil {
						return nil, err
					}
				}
			}
			if validator := reqCtx.GetIAMValidator(); validator != nil {
				if err := validator.ValidateRoleForService(ctx, roleArn, iam.ServicePrincipalEvents); err != nil {
					return nil, err
				}
			}
		}
		rule.RoleARN = roleArn
	}

	if err := store.CreateRule(ctx, rule); err != nil {
		if err == eventsstore.ErrRuleAlreadyExists {
			existingRule, err := store.GetRule(ctx, eventBusName, name)
			if err != nil {
				return nil, err
			}
			if existingRule != nil {
				if desc, ok := req.Parameters["Description"].(string); ok {
					if !validateDescription(desc) {
						return nil, errDescriptionTooLong()
					}
					existingRule.Description = desc
				}
				if pattern, ok := req.Parameters["EventPattern"].(string); ok {
					if !validateEventPatternLength(pattern) {
						return nil, awserrors.NewValidationException("EventPattern must be at most 4096 characters")
					}
					if !isValidEventPattern(pattern) {
						return nil, awserrors.NewInvalidEventPatternException("EventPattern must be valid JSON")
					}
					existingRule.EventPattern = pattern
				}
				if schedule, ok := req.Parameters["ScheduleExpression"].(string); ok {
					existingRule.ScheduleExpression = schedule
				}
				if roleArn, ok := req.Parameters["RoleArn"].(string); ok {
					existingRule.RoleARN = roleArn
				}
				if state, ok := req.Parameters["State"].(string); ok {
					existingRule.State = eventsstore.RuleState(state)
				}
				if err := store.UpdateRule(ctx, existingRule); err != nil {
					return nil, err
				}
				if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 {
					if err := store.TagStore.TagFromSlice(existingRule.ARN, tags); err != nil {
						return nil, err
					}
				}
				return map[string]interface{}{
					"RuleArn": existingRule.ARN,
				}, nil
			}
			return nil, awserrors.NewResourceAlreadyExistsException("Rule '" + name + "'")
		}
		return nil, err
	}

	if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 {
		if err := store.TagStore.TagFromSlice(rule.ARN, tags); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"RuleArn": rule.ARN,
	}, nil
}

// DeleteRule removes a rule from the event bus.
// Rules with targets cannot be deleted until targets are removed.
func (s *EventsService) DeleteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rule, err := store.GetRule(ctx, eventBusName, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	force, _ := req.Parameters["Force"].(bool)

	// Check if rule has targets. When Force is false (the default), AWS
	// rejects the delete with a DependencyViolation-style error so that
	// callers must explicitly remove targets first. When Force is true the
	// targets are cascade-deleted before the rule itself is removed.
	targetsResult, err := store.ListTargetsByRule(ctx, eventBusName, name, 1, "")
	if err != nil {
		return nil, err
	}
	if len(targetsResult.Targets) > 0 && !force {
		return nil, awserrors.NewValidationException("Rule '" + name + "' has targets. Remove targets before deleting the rule, or retry with Force=true.")
	}
	if len(targetsResult.Targets) > 0 && force {
		// Cascade-delete targets before removing the rule.
		allToken := ""
		for {
			page, err := store.ListTargetsByRule(ctx, eventBusName, name, 1000, allToken)
			if err != nil {
				return nil, err
			}
			for _, t := range page.Targets {
				if err := store.DeleteTarget(ctx, eventBusName, name, t.ID); err != nil {
					return nil, err
				}
			}
			if page.NextToken == "" {
				break
			}
			allToken = page.NextToken
		}
	}

	if err := store.DeleteRule(ctx, eventBusName, name); err != nil {
		return nil, err
	}

	// Clean up scheduler state for the deleted rule.
	lastFireTimes.Delete(rule.ARN)

	_ = store.TagStore.Delete(rule.ARN)

	return response.EmptyResponse(), nil
}

// DescribeRule returns details about a rule including its state,
// schedule expression, event pattern, and tags.
func (s *EventsService) DescribeRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rule, err := store.GetRule(ctx, eventBusName, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	result := ruleToMap(rule, false)

	if tagSlice, err := store.TagStore.ListAsSlice(rule.ARN); err == nil && len(tagSlice) > 0 {
		tagMaps := make([]map[string]string, 0, len(tagSlice))
		for _, t := range tagSlice {
			tagMaps = append(tagMaps, map[string]string{"Key": t.Key, "Value": t.Value})
		}
		result["Tags"] = tagMaps
	}

	return result, nil
}

// ListRules returns a list of rules for the specified event bus,
// optionally filtered by name prefix with pagination support.
func (s *EventsService) ListRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listRulesCore(ctx, store, ListRulesInput{
		EventBusName: eventBusName,
		NamePrefix:   namePrefix,
		Limit:        limit,
		NextToken:    nextToken,
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
	if name == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rule, err := store.GetRule(ctx, eventBusName, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	rule.State = eventsstore.RuleStateEnabled
	if err := store.UpdateRule(ctx, rule); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisableRule disables a rule so it no longer delivers events to its targets.
func (s *EventsService) DisableRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rule, err := store.GetRule(ctx, eventBusName, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	rule.State = eventsstore.RuleStateDisabled
	if err := store.UpdateRule(ctx, rule); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListRuleNamesByTarget returns the list of rules that have the specified target.
// It scans all rules on the event bus, checks each rule's targets for a match,
// then applies client-side pagination to the matched rule names.
func (s *EventsService) ListRuleNamesByTarget(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	targetArn := request.GetStringParam(req.Parameters, "TargetArn")
	if targetArn == "" {
		return nil, awserrors.NewValidationException("TargetArn is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var allRuleNames []string
	scanToken := ""
	for {
		rulesResult, err := store.ListRules(ctx, eventBusName, "", 100, scanToken)
		if err != nil {
			return nil, err
		}
		for _, rule := range rulesResult.Rules {
			targets, err := store.ListTargetsByRule(ctx, eventBusName, rule.Name, 100, "")
			if err != nil {
				// Swallowing the error here would silently omit rules
				// from the result, under-reporting which rules reference
				// the target. Fail the whole listing instead.
				return nil, err
			}
			for _, t := range targets.Targets {
				if t.ARN == targetArn {
					allRuleNames = append(allRuleNames, rule.Name)
					break
				}
			}
		}
		if rulesResult.NextToken == "" {
			break
		}
		scanToken = rulesResult.NextToken
	}

	start := 0
	if nextToken != "" {
		// AWS returns InvalidParameterException when the supplied
		// NextToken is not a recognised opaque cursor. Our pagination
		// scheme uses a zero-based integer offset, so reject any token
		// that fails to parse as a non-negative integer.
		idx, err := strconv.Atoi(nextToken)
		if err != nil || idx < 0 {
			return nil, awserrors.NewInvalidParameterException("Invalid NextToken: " + nextToken)
		}
		start = idx
	}
	end := start + int(limit)
	// Clamp both endpoints to the result-set length to prevent slice
	// bounds panics when a stale NextToken offset exceeds the current
	// length (e.g. rules were deleted between paginated calls).
	if start > len(allRuleNames) {
		start = len(allRuleNames)
	}
	if end > len(allRuleNames) {
		end = len(allRuleNames)
	}

	ruleNames := allRuleNames[start:end]

	resp := map[string]interface{}{
		"RuleNames": ruleNames,
	}

	if end < len(allRuleNames) {
		pagination.SetNextToken(resp, "NextToken", strconv.Itoa(end))
	}

	return resp, nil
}
