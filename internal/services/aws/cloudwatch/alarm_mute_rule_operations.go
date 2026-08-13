package cloudwatch

import (
	"context"
	"errors"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// PutAlarmMuteRule creates or updates an alarm mute rule.
func (s *CloudWatchService) PutAlarmMuteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := getAlarmStringParam(req.Parameters, "Name", "name")
	description := getAlarmStringParam(req.Parameters, "Description", "description")

	scheduleExpr := ""
	if rule, ok := req.Parameters["Rule"]; ok {
		if m, ok := rule.(map[string]interface{}); ok {
			if sched, ok := m["Schedule"]; ok {
				if sm, ok := sched.(map[string]interface{}); ok {
					if expr, ok := sm["Expression"]; ok {
						if e, ok := expr.(string); ok {
							scheduleExpr = e
						}
					}
				}
			}
		}
	}

	var mutedNames []string
	if mt, ok := req.Parameters["MuteTargets"]; ok {
		if m, ok := mt.(map[string]interface{}); ok {
			mutedNames = parseStringArrayParam(m, "AlarmNames", "alarmNames")
		}
	}

	startDate := parseTimestampFromMap(req.Parameters, "StartDate")
	expireDate := parseTimestampFromMap(req.Parameters, "ExpireDate")
	tags, tagErr := parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}

	if err := s.putAlarmMuteRuleCore(store, &PutAlarmMuteRuleInput{
		Name:         name,
		Description:  description,
		ScheduleExpr: scheduleExpr,
		MutedNames:   mutedNames,
		StartDate:    startDate,
		ExpireDate:   expireDate,
		Tags:         tags,
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// DeleteAlarmMuteRule deletes an alarm mute rule.
func (s *CloudWatchService) DeleteAlarmMuteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := getAlarmStringParam(req.Parameters, "AlarmMuteRuleName", "alarmMuteRuleName")

	if err := s.deleteAlarmMuteRuleCore(store, &DeleteAlarmMuteRuleInput{Name: name}); err != nil {
		if errors.Is(err, cwstore.ErrAlarmMuteRuleNotFound) {
			return nil, awserrors.NewResourceNotFoundException("AlarmMuteRule", name)
		}
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// GetAlarmMuteRule returns details for a single alarm mute rule.
func (s *CloudWatchService) GetAlarmMuteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := getAlarmStringParam(req.Parameters, "AlarmMuteRuleName", "alarmMuteRuleName")

	rule, err := s.getAlarmMuteRuleCore(store, &GetAlarmMuteRuleInput{Name: name})
	if err != nil {
		if errors.Is(err, cwstore.ErrAlarmMuteRuleNotFound) {
			return nil, awserrors.NewResourceNotFoundException("AlarmMuteRule", name)
		}
		return nil, fmt.Errorf("failed to get alarm mute rule: %w", err)
	}

	return alarmMuteRuleToResponse(rule), nil
}

// ListAlarmMuteRules lists alarm mute rules.
func (s *CloudWatchService) ListAlarmMuteRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.listAlarmMuteRulesCore(store, &ListAlarmMuteRulesInput{
		AlarmName:  getAlarmStringParam(req.Parameters, "AlarmName", "alarmName"),
		Statuses:   parseStringArrayParam(req.Parameters, "Statuses", "statuses"),
		NextToken:  pagination.GetMarker(req.Parameters, "NextToken"),
		MaxRecords: pagination.GetMaxItems(req.Parameters, 100, "MaxRecords"),
	})
	if err != nil {
		return nil, err
	}

	summaries := make([]map[string]interface{}, 0, len(items))
	for _, r := range items {
		summaries = append(summaries, alarmMuteRuleSummaryToResponse(r))
	}

	resp := map[string]interface{}{
		"AlarmMuteRuleSummaries": summaries,
	}
	if nextMarker != "" {
		resp["NextToken"] = nextMarker
	}
	return resp, nil
}

// alarmMuteRuleToResponse serialises an AlarmMuteRule into the full
// response format for GetAlarmMuteRule.
func alarmMuteRuleToResponse(r *cwstore.AlarmMuteRule) map[string]interface{} {
	resp := map[string]interface{}{
		"Name":   r.Name,
		"Status": r.Status,
	}
	if r.ARN != "" {
		resp["AlarmMuteRuleArn"] = r.ARN
	}
	if r.Description != "" {
		resp["Description"] = r.Description
	}
	if r.ScheduleExpr != "" {
		resp["Rule"] = map[string]interface{}{
			"Schedule": map[string]interface{}{
				"Expression": r.ScheduleExpr,
			},
		}
	}
	if len(r.MutedAlarmNames) > 0 {
		resp["MuteTargets"] = map[string]interface{}{
			"AlarmNames": r.MutedAlarmNames,
		}
	}
	if !r.StartDate.IsZero() {
		resp["StartDate"] = r.StartDate.Format(time.RFC3339)
	}
	if !r.ExpireDate.IsZero() {
		resp["ExpireDate"] = r.ExpireDate.Format(time.RFC3339)
	}
	if r.MuteType != "" {
		resp["MuteType"] = r.MuteType
	}
	if !r.UpdatedAt.IsZero() {
		resp["LastUpdatedTimestamp"] = r.UpdatedAt.Format(time.RFC3339)
	}
	return resp
}

// alarmMuteRuleSummaryToResponse serialises an AlarmMuteRule into the
// summary format for ListAlarmMuteRules.
func alarmMuteRuleSummaryToResponse(r *cwstore.AlarmMuteRule) map[string]interface{} {
	resp := map[string]interface{}{
		"Name":   r.Name,
		"Status": r.Status,
	}
	if r.ARN != "" {
		resp["AlarmMuteRuleArn"] = r.ARN
	}
	if r.MuteType != "" {
		resp["MuteType"] = r.MuteType
	}
	if !r.UpdatedAt.IsZero() {
		resp["LastUpdatedTimestamp"] = r.UpdatedAt.Format(time.RFC3339)
	}
	return resp
}
