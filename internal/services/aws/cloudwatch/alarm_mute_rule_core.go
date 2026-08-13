package cloudwatch

import (
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
)

// PutAlarmMuteRuleInput holds parameters for PutAlarmMuteRule.
type PutAlarmMuteRuleInput struct {
	Name         string
	Description  string
	ScheduleExpr string
	MutedNames   []string
	StartDate    time.Time
	ExpireDate   time.Time
	Tags         map[string]string
}

// DeleteAlarmMuteRuleInput holds parameters for DeleteAlarmMuteRule.
type DeleteAlarmMuteRuleInput struct {
	Name string
}

// GetAlarmMuteRuleInput holds parameters for GetAlarmMuteRule.
type GetAlarmMuteRuleInput struct {
	Name string
}

// ListAlarmMuteRulesInput holds parameters for ListAlarmMuteRules.
type ListAlarmMuteRulesInput struct {
	AlarmName  string
	Statuses   []string
	NextToken  string
	MaxRecords int
}

// putAlarmMuteRuleCore validates input and creates an alarm mute rule.
func (s *CloudWatchService) putAlarmMuteRuleCore(stores *cloudwatchStores, input *PutAlarmMuteRuleInput) error {
	if input.Name == "" {
		return awserrors.NewMissingParameter("Name is required")
	}
	if input.ScheduleExpr == "" {
		return awserrors.NewInvalidParameterValueException(
			"Rule.Schedule.Expression is required")
	}

	rule := &cwstore.AlarmMuteRule{
		Name:            input.Name,
		Description:     input.Description,
		ScheduleExpr:    input.ScheduleExpr,
		MutedAlarmNames: input.MutedNames,
		StartDate:       input.StartDate,
		ExpireDate:      input.ExpireDate,
		Tags:            input.Tags,
	}

	if _, err := stores.alarmMuteRules.PutAlarmMuteRule(rule); err != nil {
		return fmt.Errorf("failed to put alarm mute rule: %w", err)
	}
	return nil
}

// deleteAlarmMuteRuleCore validates input and deletes an alarm mute rule.
func (s *CloudWatchService) deleteAlarmMuteRuleCore(stores *cloudwatchStores, input *DeleteAlarmMuteRuleInput) error {
	if input.Name == "" {
		return awserrors.NewMissingParameter("AlarmMuteRuleName is required")
	}
	return stores.alarmMuteRules.DeleteAlarmMuteRule(input.Name)
}

// getAlarmMuteRuleCore validates input and retrieves an alarm mute rule.
func (s *CloudWatchService) getAlarmMuteRuleCore(stores *cloudwatchStores, input *GetAlarmMuteRuleInput) (*cwstore.AlarmMuteRule, error) {
	if input.Name == "" {
		return nil, awserrors.NewMissingParameter("AlarmMuteRuleName is required")
	}
	return stores.alarmMuteRules.GetAlarmMuteRule(input.Name)
}

// listAlarmMuteRulesCore validates input and lists alarm mute rules.
func (s *CloudWatchService) listAlarmMuteRulesCore(stores *cloudwatchStores, input *ListAlarmMuteRulesInput) ([]*cwstore.AlarmMuteRule, string, error) {
	maxRecords := input.MaxRecords
	if maxRecords == 0 {
		maxRecords = 100
	}
	opts := common.ListOptions{Marker: input.NextToken, MaxItems: maxRecords}
	result, err := stores.alarmMuteRules.ListAlarmMuteRulesPaginated(input.AlarmName, input.Statuses, opts)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list alarm mute rules: %w", err)
	}
	return result.Items, result.NextMarker, nil
}
