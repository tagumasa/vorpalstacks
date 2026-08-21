package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// checkActionsSuppressor evaluates whether a composite alarm's actions
// should be suppressed based on its ActionsSuppressor configuration. The
// suppressor alarm must be in ALARM state for at least WaitPeriod seconds
// before suppression begins. After the suppressor transitions out of ALARM,
// suppression continues for ExtensionPeriod seconds.
//
// Returns (true, reason) if actions should be suppressed.
func checkActionsSuppressor(alarm *cwstore.Alarm, alarmStore *cwstore.AlarmStore) (bool, string) {
	suppressorName := extractAlarmNameFromARN(alarm.ActionsSuppressor)
	if suppressorName == "" {
		suppressorName = alarm.ActionsSuppressor
	}

	suppressor, err := alarmStore.GetAlarm(suppressorName)
	if err != nil || suppressor == nil {
		return false, ""
	}

	now := time.Now().UTC()
	if suppressor.State == "ALARM" {
		waitDuration := time.Duration(alarm.ActionsSuppressorWaitPeriod) * time.Second
		if now.Sub(suppressor.StateUpdatedTimestamp) >= waitDuration {
			return true, "WaitPeriod"
		}
	} else {
		extDuration := time.Duration(alarm.ActionsSuppressorExtPeriod) * time.Second
		if extDuration > 0 && now.Sub(suppressor.StateUpdatedTimestamp) < extDuration {
			return true, "ExtensionPeriod"
		}
	}
	return false, ""
}

// extractAlarmNameFromARN extracts the alarm name from a CloudWatch alarm
// ARN, whose resource field is alarm:<alarm-name>
// (arn:<partition>:cloudwatch:<region>:<account>:alarm:<name>). Returns
// empty string if the ARN format is not recognised.
func extractAlarmNameFromARN(arn string) string {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	if strings.HasPrefix(resource, "alarm:") {
		return strings.TrimPrefix(resource, "alarm:")
	}
	return ""
}

// handleAlarmStateTransition is called when the evaluator detects an alarm
// state change. It updates the alarm state in the store, records alarm
// history, publishes a CloudWatchAlarmStateEvent via the event bus, and
// dispatches alarm actions to SNS topics and Lambda functions.
func (s *CloudWatchService) handleAlarmStateTransition(ctx context.Context, result *alarmEvalResult, alarmStore *cwstore.AlarmStore, muteRuleStore *cwstore.AlarmMuteRuleStore) {
	alarm := result.alarm

	if err := alarmStore.SetAlarmState(alarm.Name, result.newState, result.reason, ""); err != nil {
		s.log("failed to set alarm state", "alarm", alarm.Name, "new_state", result.newState, "error", err)
		return
	}

	historyType := cwstore.AlarmTypeMetricAlarm
	if alarm.AlarmType == cwstore.AlarmTypeCompositeAlarm {
		historyType = cwstore.AlarmTypeCompositeAlarm
	}

	if err := alarmStore.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarm.Name,
		AlarmType:       historyType,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeStateUpdate,
		HistorySummary:  result.reason,
	}); err != nil {
		s.log("failed to add alarm history", "alarm", alarm.Name, "error", err)
	}

	if !alarm.ActionsEnabled {
		return
	}

	// Alarm Mute Rules: if any ACTIVE mute rule targets this alarm,
	// suppress all alarm actions.
	if isAlarmMuted(result.alarm.Name, muteRuleStore) {
		return
	}

	// ActionsSuppressor: when the suppressor alarm is in ALARM state and
	// has been so for at least ActionsSuppressorWaitPeriod seconds, the
	// composite alarm's actions are suppressed. When the suppressor
	// transitions out of ALARM, suppression continues for
	// ActionsSuppressorExtensionPeriod seconds.
	if alarm.ActionsSuppressor != "" {
		if suppressed, reason := checkActionsSuppressor(alarm, alarmStore); suppressed {
			if err := alarmStore.SetAlarmActionsSuppressed(alarm.Name, reason); err != nil {
				s.log("failed to set actions suppressed", "alarm", alarm.Name, "error", err)
			}
			result.actionsToFire = nil
		}
	}

	s.publishAlarmStateEvent(ctx, result)
	s.dispatchAlarmActions(ctx, result)
}

// publishAlarmStateEvent publishes a CloudWatchAlarmStateEvent to the
// event bus. The event carries the alarm ARN, previous state, new state,
// and the reason for the transition.
func (s *CloudWatchService) publishAlarmStateEvent(ctx context.Context, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	_, _, alarmRegion, _, _ := svcarn.SplitARN(result.alarm.ARN)
	if alarmRegion == "" {
		alarmRegion = s.region
	}

	evt := &eventbus.CloudWatchAlarmStateEvent{
		EventBase: eventbus.EventBase{
			Timestamp: time.Now().UTC(),
			Source:    "aws:cloudwatch",
			Region:    alarmRegion,
			AccountID: s.accountID,
			Caller: eventbus.CallerContext{
				ServicePrincipal: "cloudwatch.amazonaws.com",
				AccountID:        s.accountID,
			},
		},
		AlarmName:     result.alarm.Name,
		AlarmARN:      result.alarm.ARN,
		PreviousState: result.oldState,
		NewState:      result.newState,
		Reason:        result.reason,
	}

	if err := s.bus.Publish(ctx, evt); err != nil {
		logs.Warn("failed to publish alarm state change event", logs.String("alarmName", result.alarm.Name), logs.Err(err))
	}
}

// dispatchAlarmActions iterates over the action ARNs for the new state and
// dispatches notifications to SNS topics (via the event bus) and Lambda
// functions (via the direct invoker). Region and account ID are extracted
// from each action ARN to support cross-region alarm actions.
func (s *CloudWatchService) dispatchAlarmActions(ctx context.Context, result *alarmEvalResult) {
	for _, actionArn := range result.actionsToFire {
		switch svcarn.GetServiceFromARN(actionArn) {
		case "sns":
			s.dispatchAlarmToSNS(ctx, actionArn, result)
		case "lambda":
			s.dispatchAlarmToLambda(ctx, actionArn, result)
		case "states":
			s.dispatchAlarmToStepFunctions(ctx, actionArn, result)
		case "sqs":
			s.dispatchAlarmToSQS(ctx, actionArn, result)
		}
	}
}

// dispatchAlarmToSNS publishes the alarm state change notification to an
// SNS topic via the event bus. Region and account ID are extracted from the
// topic ARN.
func (s *CloudWatchService) dispatchAlarmToSNS(ctx context.Context, topicArn string, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, topicArn, "sns", "cloudwatch.amazonaws.com", "sns:Publish", topicArn)
	if evalErr != nil {
		s.log("resource policy evaluation failed for alarm SNS delivery, dropping notification", "topicArn", topicArn, "error", evalErr)
		return
	}
	if !allowed {
		return
	}

	_, _, region, accountID, _ := svcarn.SplitARN(topicArn)
	// SNS message IDs are UUIDs; a clock-derived ID collides when several
	// alarms notify in the same nanosecond.
	messageID := uuid.New().String()

	payload := buildAlarmNotificationPayload(result)
	messageBytes, _ := json.Marshal(payload)

	snsEvt := &eventbus.SNSDeliveryEvent{
		EventBase: eventbus.EventBase{
			Timestamp: time.Now().UTC(),
			Source:    "aws:cloudwatch",
			Region:    region,
			AccountID: accountID,
			Caller: eventbus.CallerContext{
				ServicePrincipal: "cloudwatch.amazonaws.com",
				AccountID:        accountID,
			},
		},
		TopicARN:  topicArn,
		MessageID: messageID,
		Message:   string(messageBytes),
		Subject:   fmt.Sprintf("ALARM: \"%s\" in %s", result.alarm.Name, result.newState),
	}
	snsEvt.Region = region

	if err := s.bus.Publish(ctx, snsEvt); err != nil {
		logs.Warn("failed to publish alarm SNS notification", logs.String("alarmName", result.alarm.Name), logs.Err(err))
	}
}

// dispatchAlarmToLambda invokes a Lambda function with the alarm state
// change notification payload. The function name is extracted from the
// function ARN.
func (s *CloudWatchService) dispatchAlarmToLambda(ctx context.Context, functionArn string, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, functionArn, "lambda", "cloudwatch.amazonaws.com", "lambda:InvokeFunction", functionArn)
	if evalErr != nil {
		s.log("resource policy evaluation failed for alarm Lambda delivery, dropping notification", "functionArn", functionArn, "error", evalErr)
		return
	}
	if !allowed {
		return
	}

	fnName := svcarn.ExtractFunctionNameFromARN(functionArn)
	payload := buildAlarmNotificationPayload(result)
	payloadBytes, _ := json.Marshal(payload)

	_, _, err := s.bus.LambdaInvoker().InvokeForGateway(ctx, functionArn, payloadBytes)
	if err != nil {
		s.log("lambda dispatch failed for alarm action", "alarm", result.alarm.Name, "function", fnName, "error", err)
	}
}

func (s *CloudWatchService) dispatchAlarmToStepFunctions(ctx context.Context, stateMachineArn string, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	_, _, smRegion, _, _ := svcarn.SplitARN(stateMachineArn)
	if smRegion == "" {
		smRegion = s.region
	}

	payload := buildAlarmNotificationPayload(result)
	payloadBytes, _ := json.Marshal(payload)

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: stateMachineArn,
		Input:           string(payloadBytes),
	}
	evt.Region = smRegion
	evt.AccountID = s.accountID

	if err := s.bus.Publish(ctx, evt); err != nil {
		logs.Warn("failed to publish alarm Step Function event", logs.String("alarmName", result.alarm.Name), logs.Err(err))
	}
}

// dispatchAlarmToSQS delivers the alarm notification to an SQS queue
// via the SQS invoker. The queue name is extracted from the queue ARN.
func (s *CloudWatchService) dispatchAlarmToSQS(ctx context.Context, queueArn string, result *alarmEvalResult) {
	if s.bus == nil || s.bus.SQSInvoker() == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, queueArn, "sqs", "cloudwatch.amazonaws.com", "sqs:SendMessage", queueArn)
	if evalErr != nil {
		s.log("resource policy evaluation failed for alarm SQS delivery, dropping notification", "queueArn", queueArn, "error", evalErr)
		return
	}
	if !allowed {
		return
	}

	_, _, sqsRegion, _, queueName := svcarn.SplitARN(queueArn)
	if queueName == "" {
		return
	}

	queueURL, err := s.bus.SQSInvoker().GetQueueByName(ctx, sqsRegion, queueName)
	if err != nil {
		s.log("failed to resolve SQS queue for alarm action", "queueArn", queueArn, "error", err)
		return
	}

	payload := buildAlarmNotificationPayload(result)
	messageBytes, _ := json.Marshal(payload)

	_, _, err = s.bus.SQSInvoker().SendMessage(ctx, sqsRegion, queueURL, string(messageBytes), eventbus.SQSSendOptions{})
	if err != nil {
		s.log("SQS dispatch failed for alarm action", "alarm", result.alarm.Name, "queue", queueName, "error", err)
	}
}

// operatorPhrase returns a human-readable phrase describing the comparison
// direction, suitable for inclusion in alarm state change reason strings.
func operatorPhrase(operator string) string {
	switch operator {
	case "GreaterThanOrEqualToThreshold":
		return "were at or above"
	case "GreaterThanThreshold":
		return "were above"
	case "LessThanOrEqualToThreshold":
		return "were at or below"
	case "LessThanThreshold":
		return "were below"
	case "LessThanLowerOrGreaterThanUpperThreshold":
		return "were outside the anomaly band"
	case "LessThanLowerThreshold":
		return "were below the anomaly band lower bound"
	case "GreaterThanUpperThreshold":
		return "were above the anomaly band upper bound"
	default:
		return "crossed"
	}
}

// buildAlarmNotificationPayload constructs the CloudWatch alarm
// notification payload matching the format AWS sends to SNS topics and
// Lambda functions. This includes the alarm description, metric details,
// and state transition information.
func buildAlarmNotificationPayload(result *alarmEvalResult) map[string]interface{} {
	alarm := result.alarm
	now := time.Now().UTC()

	_, _, alarmRegion, _, _ := svcarn.SplitARN(alarm.ARN)

	return map[string]interface{}{
		"AlarmName":          alarm.Name,
		"AlarmArn":           alarm.ARN,
		"AlarmDescription":   alarm.AlarmDescription,
		"AlarmConfiguration": buildAlarmConfiguration(alarm),
		"PreviousState": map[string]interface{}{
			"StateValue":      result.oldState,
			"StateReason":     "",
			"StateReasonData": "",
		},
		"NewState": map[string]interface{}{
			"StateValue":      result.newState,
			"StateReason":     result.reason,
			"StateReasonData": "",
			"TriggeredTime":   now.Format(time.RFC3339),
		},
		"NewStateReason":     result.reason,
		"StateChangeTime":    now.Format(time.RFC3339),
		"Region":             alarmRegion,
		"MetricName":         alarm.MetricName,
		"Namespace":          alarm.Namespace,
		"Statistic":          alarm.Statistic,
		"Period":             alarm.Period,
		"EvaluationPeriods":  alarm.EvaluationPeriods,
		"Threshold":          alarm.Threshold,
		"ComparisonOperator": alarm.ComparisonOperator,
		"TreatMissingData":   alarm.TreatMissingData,
	}
}

// buildAlarmConfiguration serialises the alarm's key configuration fields
// into a nested map for inclusion in the notification payload.
func buildAlarmConfiguration(alarm *cwstore.Alarm) map[string]interface{} {
	config := map[string]interface{}{
		"AlarmName":          alarm.Name,
		"AlarmArn":           alarm.ARN,
		"AlarmType":          alarm.AlarmType,
		"MetricName":         alarm.MetricName,
		"Namespace":          alarm.Namespace,
		"Statistic":          alarm.Statistic,
		"Period":             alarm.Period,
		"EvaluationPeriods":  alarm.EvaluationPeriods,
		"Threshold":          alarm.Threshold,
		"ComparisonOperator": alarm.ComparisonOperator,
		"TreatMissingData":   alarm.TreatMissingData,
		"ActionsEnabled":     alarm.ActionsEnabled,
	}

	if alarm.DatapointsToAlarm > 0 {
		config["DatapointsToAlarm"] = alarm.DatapointsToAlarm
	}
	if len(alarm.Dimensions) > 0 {
		dims := make([]map[string]string, len(alarm.Dimensions))
		for i, d := range alarm.Dimensions {
			dims[i] = map[string]string{"Name": d.Name, "Value": d.Value}
		}
		config["Dimensions"] = dims
	}
	if len(alarm.AlarmActions) > 0 {
		config["AlarmActions"] = alarm.AlarmActions
	}
	if len(alarm.OKActions) > 0 {
		config["OKActions"] = alarm.OKActions
	}
	if len(alarm.InsufficientDataActions) > 0 {
		config["InsufficientDataActions"] = alarm.InsufficientDataActions
	}

	return config
}
