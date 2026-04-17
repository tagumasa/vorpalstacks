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
		resp, err := client.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(namespace),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(metricName),
					Value:      aws.Float64(42.0),
					Timestamp:  aws.Time(time.Now()),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		resp, err := client.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
			Namespace:  aws.String(namespace),
			MetricName: aws.String(metricName),
			StartTime:  aws.Time(now.Add(-1 * time.Hour)),
			EndTime:    aws.Time(now),
			Period:     aws.Int32(300),
			Statistics: []types.Statistic{types.StatisticAverage},
		})
		if err != nil {
			return err
		}
		if resp.Datapoints == nil {
			return fmt.Errorf("datapoints list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "PutMetricAlarm", func() error {
		alarmName := fmt.Sprintf("TestAlarm-%d", time.Now().UnixNano())
		resp, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String(metricName),
			Namespace:          aws.String(namespace),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "DescribeAlarms", func() error {
		resp, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{})
		if err != nil {
			return err
		}
		for _, alarm := range resp.MetricAlarms {
			if alarm.AlarmName == nil {
				return fmt.Errorf("alarm with nil name in response")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "ListDashboards", func() error {
		resp, err := client.ListDashboards(ctx, &cloudwatch.ListDashboardsInput{})
		if err != nil {
			return err
		}
		if len(resp.DashboardEntries) != 0 {
			return fmt.Errorf("expected no dashboards, got %d", len(resp.DashboardEntries))
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "GetMetricWidgetImage_InvalidWidget", func() error {
		_, err := client.GetMetricWidgetImage(ctx, &cloudwatch.GetMetricWidgetImageInput{
			MetricWidget: aws.String("invalid-json"),
		})
		if err := AssertErrorContains(err, "InvalidParameter"); err != nil {
			return err
		}
		return nil
	}))

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

	results = append(results, r.RunTest("cloudwatch", "GetMetricData_Basic", func() error {
		testNS := fmt.Sprintf("MetricDataNS-%d", time.Now().UnixNano())
		testMetric := fmt.Sprintf("MetricDataMetric-%d", time.Now().UnixNano())
		now := time.Now()

		_, err := client.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(testNS),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(10.0),
					Timestamp:  aws.Time(now.Add(-3 * time.Minute)),
				},
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(20.0),
					Timestamp:  aws.Time(now.Add(-1 * time.Minute)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put metric data: %v", err)
		}

		resp, err := client.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
			StartTime: aws.Time(now.Add(-10 * time.Minute)),
			EndTime:   aws.Time(now.Add(1 * time.Minute)),
			MetricDataQueries: []types.MetricDataQuery{
				{
					Id: aws.String("m1"),
					MetricStat: &types.MetricStat{
						Metric: &types.Metric{
							Namespace:  aws.String(testNS),
							MetricName: aws.String(testMetric),
						},
						Period: aws.Int32(60),
						Stat:   aws.String("Sum"),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("get metric data: %v", err)
		}
		if resp.MetricDataResults == nil {
			return fmt.Errorf("MetricDataResults is nil")
		}
		if len(resp.MetricDataResults) == 0 {
			return fmt.Errorf("expected at least 1 MetricDataResult")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "DescribeAlarmsForMetric_Basic", func() error {
		testNS := fmt.Sprintf("DAMFNS-%d", time.Now().UnixNano())
		alarmName := fmt.Sprintf("DAMFAlarm-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("CPUUtilization"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(80.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}

		resp, err := client.DescribeAlarmsForMetric(ctx, &cloudwatch.DescribeAlarmsForMetricInput{
			Namespace:  aws.String(testNS),
			MetricName: aws.String("CPUUtilization"),
		})
		if err != nil {
			return fmt.Errorf("describe alarms for metric: %v", err)
		}
		if len(resp.MetricAlarms) == 0 {
			return fmt.Errorf("expected at least 1 alarm for metric")
		}
		found := false
		for _, a := range resp.MetricAlarms {
			if a.AlarmName != nil && *a.AlarmName == alarmName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("alarm %s not found in DescribeAlarmsForMetric result", alarmName)
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "SetAlarmState_Basic", func() error {
		alarmName := fmt.Sprintf("SetStateAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("SetStateNS-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}

		_, err = client.SetAlarmState(ctx, &cloudwatch.SetAlarmStateInput{
			AlarmName:   aws.String(alarmName),
			StateValue:  types.StateValueAlarm,
			StateReason: aws.String("Test state change"),
		})
		if err != nil {
			return fmt.Errorf("set alarm state: %v", err)
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
		if descResp.MetricAlarms[0].StateValue != types.StateValueAlarm {
			return fmt.Errorf("expected ALARM state, got %s", descResp.MetricAlarms[0].StateValue)
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "DescribeAlarms_AlarmNamePrefix", func() error {
		prefix := fmt.Sprintf("PrefixAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("PrefixNS-%d", time.Now().UnixNano())

		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(prefix + "-alpha"),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm alpha: %v", err)
		}

		_, err = client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(prefix + "-beta"),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm beta: %v", err)
		}

		_, err = client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String("OtherAlarm-no-match"),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm other: %v", err)
		}

		resp, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNamePrefix: aws.String(prefix),
		})
		if err != nil {
			return fmt.Errorf("describe with prefix: %v", err)
		}
		if len(resp.MetricAlarms) != 2 {
			return fmt.Errorf("expected 2 alarms with prefix %s, got %d", prefix, len(resp.MetricAlarms))
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{prefix + "-alpha", prefix + "-beta", "OtherAlarm-no-match"},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "TagResource_Basic", func() error {
		alarmName := fmt.Sprintf("TagAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("TagNS-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
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
		if len(descResp.MetricAlarms) != 1 || descResp.MetricAlarms[0].AlarmArn == nil {
			return fmt.Errorf("failed to get alarm ARN")
		}
		alarmARN := descResp.MetricAlarms[0].AlarmArn

		_, err = client.TagResource(ctx, &cloudwatch.TagResourceInput{
			ResourceARN: alarmARN,
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Team"), Value: aws.String("platform")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag resource: %v", err)
		}

		tagResp, err := client.ListTagsForResource(ctx, &cloudwatch.ListTagsForResourceInput{
			ResourceARN: alarmARN,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tagResp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tagResp.Tags))
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "UntagResource_Basic", func() error {
		alarmName := fmt.Sprintf("UntagAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("UntagNS-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
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
		if len(descResp.MetricAlarms) != 1 || descResp.MetricAlarms[0].AlarmArn == nil {
			return fmt.Errorf("failed to get alarm ARN")
		}
		alarmARN := descResp.MetricAlarms[0].AlarmArn

		_, err = client.TagResource(ctx, &cloudwatch.TagResourceInput{
			ResourceARN: alarmARN,
			Tags: []types.Tag{
				{Key: aws.String("Keep"), Value: aws.String("yes")},
				{Key: aws.String("Remove"), Value: aws.String("yes")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag resource: %v", err)
		}

		_, err = client.UntagResource(ctx, &cloudwatch.UntagResourceInput{
			ResourceARN: alarmARN,
			TagKeys:     []string{"Remove"},
		})
		if err != nil {
			return fmt.Errorf("untag resource: %v", err)
		}

		tagResp, err := client.ListTagsForResource(ctx, &cloudwatch.ListTagsForResourceInput{
			ResourceARN: alarmARN,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tagResp.Tags) != 1 {
			return fmt.Errorf("expected 1 tag after untag, got %d", len(tagResp.Tags))
		}
		if tagResp.Tags[0].Key == nil || *tagResp.Tags[0].Key != "Keep" {
			return fmt.Errorf("expected Keep tag, got %v", tagResp.Tags[0].Key)
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "GetMetricWidgetImage_Basic", func() error {
		widget := `{"metrics":[["AWS/EC2","CPUUtilization"]]}`
		resp, err := client.GetMetricWidgetImage(ctx, &cloudwatch.GetMetricWidgetImageInput{
			MetricWidget: aws.String(widget),
		})
		if err != nil {
			return fmt.Errorf("get metric widget image: %v", err)
		}
		if len(resp.MetricWidgetImage) == 0 {
			return fmt.Errorf("expected non-empty image")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "EnableAlarmActions_Basic", func() error {
		alarmName := fmt.Sprintf("EnableAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("EnableNS-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}

		_, err = client.DisableAlarmActions(ctx, &cloudwatch.DisableAlarmActionsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("disable actions: %v", err)
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
		if descResp.MetricAlarms[0].ActionsEnabled != nil && *descResp.MetricAlarms[0].ActionsEnabled {
			return fmt.Errorf("expected ActionsEnabled=false after disable")
		}

		_, err = client.EnableAlarmActions(ctx, &cloudwatch.EnableAlarmActionsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("enable actions: %v", err)
		}

		descResp2, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe after enable: %v", err)
		}
		if descResp2.MetricAlarms[0].ActionsEnabled == nil || !*descResp2.MetricAlarms[0].ActionsEnabled {
			return fmt.Errorf("expected ActionsEnabled=true after enable")
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "DescribeAlarmHistory_Basic", func() error {
		alarmName := fmt.Sprintf("HistoryAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("HistoryNS-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}

		_, err = client.SetAlarmState(ctx, &cloudwatch.SetAlarmStateInput{
			AlarmName:   aws.String(alarmName),
			StateValue:  types.StateValueAlarm,
			StateReason: aws.String("Manual alarm state change"),
		})
		if err != nil {
			return fmt.Errorf("set state: %v", err)
		}

		histResp, err := client.DescribeAlarmHistory(ctx, &cloudwatch.DescribeAlarmHistoryInput{
			AlarmName: aws.String(alarmName),
		})
		if err != nil {
			return fmt.Errorf("describe alarm history: %v", err)
		}
		if len(histResp.AlarmHistoryItems) == 0 {
			return fmt.Errorf("expected alarm history items, got 0")
		}
		hasStateUpdate := false
		for _, item := range histResp.AlarmHistoryItems {
			if item.HistoryItemType == types.HistoryItemTypeStateUpdate {
				hasStateUpdate = true
				break
			}
		}
		if !hasStateUpdate {
			return fmt.Errorf("expected StateUpdate history item")
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "DescribeAlarmHistory_FilterByType", func() error {
		alarmName := fmt.Sprintf("HistFilterAlarm-%d", time.Now().UnixNano())
		testNS := fmt.Sprintf("HistFilterNS-%d", time.Now().UnixNano())
		_, err := client.PutMetricAlarm(ctx, &cloudwatch.PutMetricAlarmInput{
			AlarmName:          aws.String(alarmName),
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			EvaluationPeriods:  aws.Int32(1),
			MetricName:         aws.String("TestMetric"),
			Namespace:          aws.String(testNS),
			Period:             aws.Int32(300),
			Threshold:          aws.Float64(50.0),
			Statistic:          types.StatisticAverage,
		})
		if err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}

		_, err = client.SetAlarmState(ctx, &cloudwatch.SetAlarmStateInput{
			AlarmName:   aws.String(alarmName),
			StateValue:  types.StateValueOk,
			StateReason: aws.String("Recovered"),
		})
		if err != nil {
			return fmt.Errorf("set state: %v", err)
		}

		histResp, err := client.DescribeAlarmHistory(ctx, &cloudwatch.DescribeAlarmHistoryInput{
			AlarmName:       aws.String(alarmName),
			HistoryItemType: types.HistoryItemTypeStateUpdate,
		})
		if err != nil {
			return fmt.Errorf("describe alarm history: %v", err)
		}
		for _, item := range histResp.AlarmHistoryItems {
			if item.HistoryItemType != types.HistoryItemTypeStateUpdate {
				return fmt.Errorf("expected only StateUpdate items, got %s", item.HistoryItemType)
			}
		}
		if len(histResp.AlarmHistoryItems) == 0 {
			return fmt.Errorf("expected at least 1 StateUpdate item")
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "PutCompositeAlarm_Basic", func() error {
		alarmName := fmt.Sprintf("CompositeAlarm-%d", time.Now().UnixNano())
		_, err := client.PutCompositeAlarm(ctx, &cloudwatch.PutCompositeAlarmInput{
			AlarmName:        aws.String(alarmName),
			AlarmRule:        aws.String("TRUE"),
			AlarmDescription: aws.String("Test composite alarm"),
			ActionsEnabled:   aws.Bool(true),
			AlarmActions:     []string{"arn:aws:sns:us-east-1:123456789012:my-topic"},
			OKActions:        []string{"arn:aws:sns:us-east-1:123456789012:ok-topic"},
		})
		if err != nil {
			return fmt.Errorf("put composite alarm: %v", err)
		}

		descResp, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmTypes: []types.AlarmType{types.AlarmTypeCompositeAlarm},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.CompositeAlarms) == 0 {
			return fmt.Errorf("expected at least 1 composite alarm, got 0")
		}
		found := false
		for _, a := range descResp.CompositeAlarms {
			if a.AlarmName != nil && *a.AlarmName == alarmName {
				found = true
				if a.AlarmRule == nil || *a.AlarmRule != "TRUE" {
					return fmt.Errorf("expected AlarmRule=TRUE, got %v", a.AlarmRule)
				}
				if a.AlarmDescription == nil || *a.AlarmDescription != "Test composite alarm" {
					return fmt.Errorf("description mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("composite alarm %s not found", alarmName)
		}
		client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "DeleteAlarms_NonExistent", func() error {
		_, err := client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{fmt.Sprintf("NonExistentAlarm-%d", time.Now().UnixNano())},
		})
		if err := AssertErrorContains(err, "ResourceNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudwatch", "PutMetricData_InvalidNamespace", func() error {
		_, err := client.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
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

	results = append(results, r.RunTest("cloudwatch", "PutCompositeAlarm_DeleteAlarm", func() error {
		alarmName := fmt.Sprintf("CompDelAlarm-%d", time.Now().UnixNano())
		_, err := client.PutCompositeAlarm(ctx, &cloudwatch.PutCompositeAlarmInput{
			AlarmName: aws.String(alarmName),
			AlarmRule: aws.String("FALSE"),
		})
		if err != nil {
			return fmt.Errorf("put composite alarm: %v", err)
		}

		descResp, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmTypes: []types.AlarmType{types.AlarmTypeCompositeAlarm},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.CompositeAlarms) == 0 {
			return fmt.Errorf("expected at least 1 composite alarm, got 0")
		}

		_, err = client.DeleteAlarms(ctx, &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		descResp2, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
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
