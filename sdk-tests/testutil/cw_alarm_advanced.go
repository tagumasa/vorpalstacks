package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (tc *cloudwatchTestCtx) alarmAdvancedTests() []TestResult {
	var results []TestResult
	reg := tc.runner.Region()
	acct := tc.runner.AccountID()

	results = append(results, tc.runner.RunTest("cloudwatch", "SetAlarmState_Verify", func() error {
		alarmName := tc.uniquePrefix("StateAlarm")
		testNS := tc.uniquePrefix("StateNS")
		// A manually set alarm state holds only until the next background
		// evaluation, and under TEST_MODE the evaluator re-evaluates every
		// alarm once per second. On a dataless metric that evaluation is
		// INSUFFICIENT_DATA, so a tick landing between the manual set and
		// the read-back below overwrote the manual ALARM. Seed breaching
		// data (100 over the threshold of 50) so every evaluation agrees
		// with the manual state: createAlarm pins Period=300 and the
		// evaluator window is the last completed 300 s bucket, so one
		// point inside that bucket and one in the current open bucket keep
		// every evaluation for the rest of this block and all of the next
		// one breaching.
		now := time.Now().UTC()
		lastCompletedBucket := now.Truncate(300 * time.Second).Add(-1 * time.Second)
		if _, err := tc.client.PutMetricData(tc.ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(testNS),
			MetricData: []types.MetricDatum{
				{MetricName: aws.String("TestMetric"), Value: aws.Float64(100), Timestamp: aws.Time(lastCompletedBucket)},
				{MetricName: aws.String("TestMetric"), Value: aws.Float64(100), Timestamp: aws.Time(now)},
			},
		}); err != nil {
			return fmt.Errorf("seed metric data: %v", err)
		}
		if err := tc.createAlarm(alarmName, testNS, "TestMetric", 50.0); err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		// Let the evaluator settle the alarm in ALARM before the manual
		// override: an evaluation transition that lands after the manual
		// set would replace the manual state reason even though the state
		// value agrees. Once the alarm is in ALARM every later evaluation
		// is ALARM to ALARM, which the evaluator does not write.
		if err := tc.waitForAlarmState(alarmName, types.StateValueAlarm, 10*time.Second); err != nil {
			return err
		}

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
			AlarmActions:     []string{fmt.Sprintf("arn:aws:sns:%s:%s:my-topic", reg, acct)},
			OKActions:        []string{fmt.Sprintf("arn:aws:sns:%s:%s:ok-topic", reg, acct)},
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

// waitForAlarmState polls DescribeAlarms until the named alarm reports
// the wanted state value. Alarm state is written asynchronously by the
// background evaluator, so a test that depends on an evaluated state
// must wait for the evaluator's write instead of racing it.
func (tc *cloudwatchTestCtx) waitForAlarmState(alarmName string, want types.StateValue, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		resp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe while waiting for %s: %v", want, err)
		}
		if len(resp.MetricAlarms) == 1 && resp.MetricAlarms[0].StateValue == want {
			return nil
		}
		if time.Now().After(deadline) {
			got := "no such alarm"
			if len(resp.MetricAlarms) == 1 {
				got = string(resp.MetricAlarms[0].StateValue)
			}
			return fmt.Errorf("alarm %s never reached %s within %s (last seen: %s)", alarmName, want, timeout, got)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
