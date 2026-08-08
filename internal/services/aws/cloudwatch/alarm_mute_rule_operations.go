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
	"vorpalstacks/internal/store/aws/common"
)

// PutAlarmMuteRule creates or updates an alarm mute rule that suppresses
// alarm actions for specified alarms during a time window.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_PutAlarmMuteRule.html
func (s *CloudWatchService) PutAlarmMuteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := getAlarmStringParam(req.Parameters, "Name", "name")
	if name == "" {
		return nil, awserrors.NewMissingParameter("Name is required")
	}

	description := getAlarmStringParam(req.Parameters, "Description", "description")

	// Parse Rule.Schedule.Expression.
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

	// Parse MuteTargets.AlarmNames.
	var mutedNames []string
	if mt, ok := req.Parameters["MuteTargets"]; ok {
		if m, ok := mt.(map[string]interface{}); ok {
			mutedNames = parseStringArrayParam(m, "AlarmNames", "alarmNames")
		}
	}

	if scheduleExpr == "" {
		return nil, awserrors.NewInvalidParameterValueException(
			"Rule.Schedule.Expression is required")
	}

	startDate := parseTimestampFromMap(req.Parameters, "StartDate")
	expireDate := parseTimestampFromMap(req.Parameters, "ExpireDate")
	tags, tagErr := parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}

	rule := &cwstore.AlarmMuteRule{
		Name:            name,
		Description:     description,
		ScheduleExpr:    scheduleExpr,
		MutedAlarmNames: mutedNames,
		StartDate:       startDate,
		ExpireDate:      expireDate,
		Tags:            tags,
	}

	result, err := store.alarmMuteRules.PutAlarmMuteRule(rule)
	if err != nil {
		return nil, fmt.Errorf("failed to put alarm mute rule: %w", err)
	}

	_ = result
	return map[string]interface{}{}, nil
}

// DeleteAlarmMuteRule deletes an alarm mute rule.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_DeleteAlarmMuteRule.html
func (s *CloudWatchService) DeleteAlarmMuteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := getAlarmStringParam(req.Parameters, "AlarmMuteRuleName", "alarmMuteRuleName")
	if name == "" {
		return nil, awserrors.NewMissingParameter("AlarmMuteRuleName is required")
	}

	if err := store.alarmMuteRules.DeleteAlarmMuteRule(name); err != nil {
		if errors.Is(err, cwstore.ErrAlarmMuteRuleNotFound) {
			return nil, awserrors.NewResourceNotFoundException("AlarmMuteRule", name)
		}
		return nil, fmt.Errorf("failed to delete alarm mute rule: %w", err)
	}

	return map[string]interface{}{}, nil
}

// GetAlarmMuteRule returns details for a single alarm mute rule.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_GetAlarmMuteRule.html
func (s *CloudWatchService) GetAlarmMuteRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := getAlarmStringParam(req.Parameters, "AlarmMuteRuleName", "alarmMuteRuleName")
	if name == "" {
		return nil, awserrors.NewMissingParameter("AlarmMuteRuleName is required")
	}

	rule, err := store.alarmMuteRules.GetAlarmMuteRule(name)
	if err != nil {
		if errors.Is(err, cwstore.ErrAlarmMuteRuleNotFound) {
			return nil, awserrors.NewResourceNotFoundException("AlarmMuteRule", name)
		}
		return nil, fmt.Errorf("failed to get alarm mute rule: %w", err)
	}

	return alarmMuteRuleToResponse(rule), nil
}

// ListAlarmMuteRules lists alarm mute rules, optionally filtered by
// alarm name and statuses.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_ListAlarmMuteRules.html
func (s *CloudWatchService) ListAlarmMuteRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	statuses := parseStringArrayParam(req.Parameters, "Statuses", "statuses")

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxRecords")

	opts := common.ListOptions{Marker: marker, MaxItems: maxResults}
	result, err := store.alarmMuteRules.ListAlarmMuteRulesPaginated(alarmName, statuses, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list alarm mute rules: %w", err)
	}

	summaries := make([]map[string]interface{}, 0, len(result.Items))
	for _, r := range result.Items {
		summaries = append(summaries, alarmMuteRuleSummaryToResponse(r))
	}

	resp := map[string]interface{}{
		"AlarmMuteRuleSummaries": summaries,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
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
// summary format for ListAlarmMuteRules.  The summary contains fewer
// fields than the full GetAlarmMuteRule response — only Name, Status,
// ARN, MuteType, and LastUpdatedTimestamp.
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
