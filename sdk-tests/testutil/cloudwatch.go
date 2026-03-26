package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCloudWatchTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudwatch",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cloudwatch.NewFromConfig(cfg)
	ctx := context.Background()

	namespace := fmt.Sprintf("TestNamespace-%d", time.Now().UnixNano())
	metricName := fmt.Sprintf("TestMetric-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("cloudwatch", "PutMetricData", func() error {
		_, err := client.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(namespace),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(metricName),
					Value:      aws.Float64(42.0),
					Timestamp:  aws.Time(time.Now()),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("cloudwatch", "ListMetrics", func() error {
		resp, err := client.ListMetrics(ctx, &cloudwatch.ListMetricsInput{
			Namespace: aws.String(namespace),
		})
		if err != nil {
			return err
		}
		if resp.Metrics == nil {
			return fmt.Errorf("metrics list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "GetMetricStatistics", func() error {
		now := time.Now()
		_, err := client.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
			Namespace:  aws.String(namespace),
			MetricName: aws.String(metricName),
			StartTime:  aws.Time(now.Add(-1 * time.Hour)),
			EndTime:    aws.Time(now),
			Period:     aws.Int32(300),
			Statistics: []types.Statistic{types.StatisticAverage},
		})
		return err
	}))

	results = append(results, r.RunTest("cloudwatch", "PutMetricAlarm", func() error {
		alarmName := fmt.Sprintf("TestAlarm-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String(metricName),
			Namespace:          aws.String(namespace),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		return err
	}))

	results = append(results, r.RunTest("cloudwatch", "DescribeAlarms", func() error {
		resp, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{})
		if err != nil {
			return err
		}
		if resp.MetricAlarms == nil {
			return fmt.Errorf("metric alarms list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "ListDashboards", func() error {
		resp, err := client.ListDashboards(ctx, &cloudwatch.ListDashboardsInput{})
		if err != nil {
			return err
		}
		_ = resp
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("cloudwatch", "PutMetricData_GetMetricStatistics_Roundtrip", func() error {
		testNS := fmt.Sprintf("RoundtripNS-%d", time.Now().UnixNano())
		testMetric := fmt.Sprintf("RoundtripMetric-%d", time.Now().UnixNano())
		now := time.Now()

		_, err := client.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(testNS),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(42.0),
					Unit:       types.StandardUnitNone,
					Timestamp:  aws.Time(now.Add(-5 * time.Minute)),
				},
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(58.0),
					Unit:       types.StandardUnitNone,
					Timestamp:  aws.Time(now.Add(-2 * time.Minute)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		listResp, err := client.ListMetrics(ctx, &cloudwatch.ListMetricsInput{
			Namespace:  aws.String(testNS),
			MetricName: aws.String(testMetric),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(listResp.Metrics) == 0 {
			return fmt.Errorf("metric not found in ListMetrics")
		}

		_, err = client.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
			Namespace:  aws.String(testNS),
			MetricName: aws.String(testMetric),
			StartTime:  aws.Time(now.Add(-10 * time.Minute)),
			EndTime:    aws.Time(now.Add(1 * time.Minute)),
			Period:     aws.Int32(60),
			Statistics: []types.Statistic{types.StatisticSum},
		})
		if err != nil {
			return fmt.Errorf("get stats: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "DescribeAlarms_NonExistent", func() error {
		alarmName := fmt.Sprintf("NonExistentAlarm-%d", time.Now().UnixNano())
		resp, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
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

	results = append(results, r.RunTest("cloudwatch", "PutMetricAlarm_DeleteAlarm", func() error {
		alarmName := fmt.Sprintf("DeleteAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("AlarmNS-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
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

		descResp, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
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

		_, err = client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		return nil
	}))

	return results
}
