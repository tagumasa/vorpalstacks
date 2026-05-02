package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (tc *cloudwatchTestCtx) edgeTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("cloudwatch", "DeleteAlarms_NonExistent", func() error {
		_, err := tc.client.DeleteAlarms(tc.ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{tc.uniquePrefix("NonExistentAlarm")},
		})
		if err := AssertErrorContains(err, "ResourceNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "PutMetricData_InvalidNamespace", func() error {
		_, err := tc.client.PutMetricData(tc.ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(""),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String("TestMetric"),
					Value:      aws.Float64(1.0),
				},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameter"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "PutCompositeAlarm_DeleteAlarm", func() error {
		alarmName := tc.uniquePrefix("CompDelAlarm")
		_, err := tc.client.PutCompositeAlarm(tc.ctx, &cloudwatch.PutCompositeAlarmInput{
			AlarmName: aws.String(alarmName),
			AlarmRule: aws.String("FALSE"),
		})
		if err != nil {
			return fmt.Errorf("put composite alarm: %v", err)
		}

		descResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmTypes: []types.AlarmType{types.AlarmTypeCompositeAlarm},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.CompositeAlarms) == 0 {
			return fmt.Errorf("expected at least 1 composite alarm, got 0")
		}

		_, err = tc.client.DeleteAlarms(tc.ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		descResp2, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmTypes: []types.AlarmType{types.AlarmTypeCompositeAlarm},
		})
		if err != nil {
			return fmt.Errorf("describe after delete: %v", err)
		}
		for _, a := range descResp2.CompositeAlarms {
			if a.AlarmName != nil && *a.AlarmName == alarmName {
				return fmt.Errorf("alarm %s should have been deleted", alarmName)
			}
		}
		return nil
	}))

	return results
}
