package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (tc *cloudwatchTestCtx) alarmTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("cloudwatch", "PutMetricAlarm_VerifyFields", func() error {
		alarmName := tc.uniquePrefix("FieldAlarm")
		testNS := tc.uniquePrefix("FieldNS")
		metricName := tc.uniquePrefix("FieldMetric")
		_, err := tc.client.PutMetricAlarm(tc.ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			AlarmDescription:   aws.String("Test alarm description"),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(2),
			MetricName:         aws.String(metricName),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(80.0),
			Statistic:          types.StatisticAverage,
			ActionsEnabled:     aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("put metric alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		descResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.MetricAlarms) != 1 {
			return fmt.Errorf("expected 1 alarm, got %d", len(descResp.MetricAlarms))
		}
		a := descResp.MetricAlarms[0]
		if a.AlarmName == nil || *a.AlarmName != alarmName {
			return fmt.Errorf("alarm name mismatch: got %v", a.AlarmName)
		}
		if a.AlarmDescription == nil || *a.AlarmDescription != "Test alarm description" {
			return fmt.Errorf("description mismatch: got %v", a.AlarmDescription)
		}
		if a.ComparisonOperator != types.ComparisonOperatorGreaterThanThreshold {
			return fmt.Errorf("comparison operator mismatch: got %s", a.ComparisonOperator)
		}
		if a.EvaluationPeriods == nil || *a.EvaluationPeriods != 2 {
			return fmt.Errorf("evaluation periods mismatch: got %d", func() int32 {
				if a.EvaluationPeriods != nil {
					return *a.EvaluationPeriods
				}
				return -1
			}())
		}
		if a.Threshold == nil || *a.Threshold != 80.0 {
			return fmt.Errorf("threshold mismatch: got %v", a.Threshold)
		}
		if a.Statistic != types.StatisticAverage {
			return fmt.Errorf("statistic mismatch: got %s", a.Statistic)
		}
		if a.Namespace == nil || *a.Namespace != testNS {
			return fmt.Errorf("namespace mismatch: got %v", a.Namespace)
		}
		if a.MetricName == nil || *a.MetricName != metricName {
			return fmt.Errorf("metric name mismatch: got %v", a.MetricName)
		}
		if a.Period == nil || *a.Period != 300 {
			return fmt.Errorf("period mismatch: got %v", a.Period)
		}
		if a.AlarmArn == nil || *a.AlarmArn == "" {
			return fmt.Errorf("alarm ARN is empty")
		}
		if a.StateValue != types.StateValueInsufficientData {
			return fmt.Errorf("initial state should be INSUFFICIENT_DATA, got %s", a.StateValue)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "DescribeAlarms_ForMetricAndPrefix", func() error {
		prefix := tc.uniquePrefix("MultiAlarm")
		testNS := tc.uniquePrefix("MultiNS")

		if err := tc.createAlarm(prefix+"-a", testNS, "CPUUtilization", 50.0); err != nil {
			return fmt.Errorf("create alarm a: %v", err)
		}
		if err := tc.createAlarm(prefix+"-b", testNS, "CPUUtilization", 70.0); err != nil {
			return fmt.Errorf("create alarm b: %v", err)
		}
		if err := tc.createAlarm(prefix+"-c", testNS, "MemoryUtilization", 90.0); err != nil {
			return fmt.Errorf("create alarm c: %v", err)
		}
		defer tc.deleteAlarms(prefix+"-a", prefix+"-b", prefix+"-c")

		prefixResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNamePrefix: aws.String(prefix),
		})
		if err != nil {
			return fmt.Errorf("describe by prefix: %v", err)
		}
		if len(prefixResp.MetricAlarms) != 3 {
			return fmt.Errorf("expected 3 alarms with prefix, got %d", len(prefixResp.MetricAlarms))
		}

		metricResp, err := tc.client.DescribeAlarmsForMetric(tc.ctx, &cloudwatch.DescribeAlarmsForMetricInput{
			Namespace:  aws.String(testNS),
			MetricName: aws.String("CPUUtilization"),
		})
		if err != nil {
			return fmt.Errorf("describe for metric: %v", err)
		}
		if len(metricResp.MetricAlarms) < 2 {
			return fmt.Errorf("expected >= 2 alarms for CPUUtilization, got %d", len(metricResp.MetricAlarms))
		}
		for _, a := range metricResp.MetricAlarms {
			if a.MetricName == nil || *a.MetricName != "CPUUtilization" {
				return fmt.Errorf("unexpected metric name in DescribeAlarmsForMetric: %v", a.MetricName)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "PutMetricAlarm_DeleteAlarm_Verify", func() error {
		alarmName := tc.uniquePrefix("DelVerifyAlarm")
		testNS := tc.uniquePrefix("DelVerifyNS")
		_, err := tc.client.PutMetricAlarm(tc.ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
			AlarmDescription:   aws.String("Test alarm for deletion"),
		})
		if err != nil {
			return fmt.Errorf("put alarm: %v", err)
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
		if descResp.MetricAlarms[0].AlarmDescription == nil || *descResp.MetricAlarms[0].AlarmDescription != "Test alarm for deletion" {
			return fmt.Errorf("alarm description mismatch")
		}

		_, err = tc.client.DeleteAlarms(tc.ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		descResp2, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe after delete: %v", err)
		}
		for _, a := range descResp2.MetricAlarms {
			if a.AlarmName != nil && *a.AlarmName == alarmName {
				return fmt.Errorf("alarm %s should have been deleted", alarmName)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "DescribeAlarms_NonExistent", func() error {
		alarmName := tc.uniquePrefix("NonExistentAlarm")
		resp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("DescribeAlarms should not error for non-existent alarm, got: %v", err)
		}
		if len(resp.MetricAlarms) != 0 {
			return fmt.Errorf("expected no alarms, got %d", len(resp.MetricAlarms))
		}
		return nil
	}))

	return results
}
