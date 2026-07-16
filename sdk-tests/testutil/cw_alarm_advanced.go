package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (tc *cloudwatchTestCtx) alarmAdvancedTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("cloudwatch", "SetAlarmState_Verify", func() error {
		alarmName := tc.uniquePrefix("StateAlarm")
		testNS := tc.uniquePrefix("StateNS")
		if err := tc.createAlarm(alarmName, testNS, "TestMetric", 50.0); err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		_, err := tc.client.SetAlarmState(tc.ctx, &cloudwatch.SetAlarmStateInput{
			AlarmName:   aws.String(alarmName),
			StateValue:  types.StateValueAlarm,
			StateReason: aws.String("Manual test trigger"),
		})
		if err != nil {
			return fmt.Errorf("set alarm state: %v", err)
		}

		descResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.MetricAlarms) != 1 {
			return fmt.Errorf("expected 1 alarm, got %d", len(descResp.MetricAlarms))
		}
		if descResp.MetricAlarms[0].StateValue != types.StateValueAlarm {
			return fmt.Errorf("expected ALARM state, got %s", descResp.MetricAlarms[0].StateValue)
		}
		if descResp.MetricAlarms[0].StateReason == nil || *descResp.MetricAlarms[0].StateReason != "Manual test trigger" {
			return fmt.Errorf("state reason mismatch: got %v", descResp.MetricAlarms[0].StateReason)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "EnableDisableAlarmActions", func() error {
		alarmName := tc.uniquePrefix("ActionsAlarm")
		testNS := tc.uniquePrefix("ActionsNS")
		if err := tc.createAlarm(alarmName, testNS, "TestMetric", 50.0); err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		_, err := tc.client.DisableAlarmActions(tc.ctx, &cloudwatch.DisableAlarmActionsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("disable actions: %v", err)
		}

		descResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.MetricAlarms) != 1 {
			return fmt.Errorf("expected 1 alarm, got %d", len(descResp.MetricAlarms))
		}
		if descResp.MetricAlarms[0].ActionsEnabled != nil && *descResp.MetricAlarms[0].ActionsEnabled {
			return fmt.Errorf("expected ActionsEnabled=false after disable")
		}

		_, err = tc.client.EnableAlarmActions(tc.ctx, &cloudwatch.EnableAlarmActionsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("enable actions: %v", err)
		}

		descResp2, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe after enable: %v", err)
		}
		if descResp2.MetricAlarms[0].ActionsEnabled == nil || !*descResp2.MetricAlarms[0].ActionsEnabled {
			return fmt.Errorf("expected ActionsEnabled=true after enable")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "DescribeAlarmHistory_StateUpdate", func() error {
		alarmName := tc.uniquePrefix("HistAlarm")
		testNS := tc.uniquePrefix("HistNS")
		if err := tc.createAlarm(alarmName, testNS, "TestMetric", 50.0); err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		_, err := tc.client.SetAlarmState(tc.ctx, &cloudwatch.SetAlarmStateInput{
			AlarmName:   aws.String(alarmName),
			StateValue:  types.StateValueAlarm,
			StateReason: aws.String("Manual state change"),
		})
		if err != nil {
			return fmt.Errorf("set state: %v", err)
		}

		histResp, err := tc.client.DescribeAlarmHistory(tc.ctx, &cloudwatch.DescribeAlarmHistoryInput{
			AlarmName:       aws.String(alarmName),
			HistoryItemType: types.HistoryItemTypeStateUpdate,
		})
		if err != nil {
			return fmt.Errorf("describe alarm history: %v", err)
		}
		if len(histResp.AlarmHistoryItems) == 0 {
			return fmt.Errorf("expected alarm history items")
		}
		for _, item := range histResp.AlarmHistoryItems {
			if item.HistoryItemType != types.HistoryItemTypeStateUpdate {
				return fmt.Errorf("expected only StateUpdate items, got %s", item.HistoryItemType)
			}
			if item.AlarmName == nil || *item.AlarmName != alarmName {
				return fmt.Errorf("alarm name mismatch in history: got %v", item.AlarmName)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "PutCompositeAlarm_VerifyFields", func() error {
		alarmName := tc.uniquePrefix("CompAlarm")
		_, err := tc.client.PutCompositeAlarm(tc.ctx, &cloudwatch.PutCompositeAlarmInput{
			AlarmName:        aws.String(alarmName),
			AlarmRule:        aws.String("TRUE"),
			AlarmDescription: aws.String("Composite test alarm"),
			ActionsEnabled:   aws.Bool(true),
			AlarmActions:     []string{"arn:aws:sns:us-east-1:123456789012:my-topic"},
			OKActions:        []string{"arn:aws:sns:us-east-1:123456789012:ok-topic"},
		})
		if err != nil {
			return fmt.Errorf("put composite alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		descResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmTypes: []types.AlarmType{types.AlarmTypeCompositeAlarm},
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.CompositeAlarms) != 1 {
			return fmt.Errorf("expected 1 composite alarm, got %d", len(descResp.CompositeAlarms))
		}
		a := descResp.CompositeAlarms[0]
		if a.AlarmName == nil || *a.AlarmName != alarmName {
			return fmt.Errorf("alarm name mismatch: got %v", a.AlarmName)
		}
		if a.AlarmRule == nil || *a.AlarmRule != "TRUE" {
			return fmt.Errorf("expected AlarmRule=TRUE, got %v", a.AlarmRule)
		}
		if a.AlarmDescription == nil || *a.AlarmDescription != "Composite test alarm" {
			return fmt.Errorf("description mismatch: got %v", a.AlarmDescription)
		}
		if a.AlarmArn == nil || *a.AlarmArn == "" {
			return fmt.Errorf("alarm ARN is empty")
		}
		if a.StateValue != types.StateValueInsufficientData {
			return fmt.Errorf("initial state should be INSUFFICIENT_DATA, got %s", a.StateValue)
		}
		return nil
	}))

	return results
}
