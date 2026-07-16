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

type cloudwatchTestCtx struct {
	client *cloudwatch.Client
	ctx    context.Context
	runner *TestRunner
}

func (r *TestRunner) RunCloudWatchTests() []TestResult {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{
			Service:  "cloudwatch",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		}}
	}

	tc := &cloudwatchTestCtx{
		client: cloudwatch.NewFromConfig(cfg),
		ctx:    context.Background(),
		runner: r,
	}

	return tc.runAll()
}

func (tc *cloudwatchTestCtx) runAll() []TestResult {
	var results []TestResult

	results = append(results, tc.metricTests()...)
	results = append(results, tc.alarmTests()...)
	results = append(results, tc.alarmAdvancedTests()...)
	results = append(results, tc.dashboardTests()...)
	results = append(results, tc.widgetTests()...)
	results = append(results, tc.tagTests()...)
	results = append(results, tc.edgeTests()...)

	return results
}

func (tc *cloudwatchTestCtx) uniquePrefix(base string) string {
	return fmt.Sprintf("%s-%d", base, time.Now().UnixNano())
}

func (tc *cloudwatchTestCtx) createAlarm(alarmName, testNS, metricName string, threshold float64) error {
	_, err := tc.client.PutMetricAlarm(tc.ctx, &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
		EvaluationPeriods:  aws.Int32(1),
		MetricName:         aws.String(metricName),
		Namespace:          aws.String(testNS),
		Period:             aws.Int32(300),
		Threshold:          aws.Float64(threshold),
		Statistic:          types.StatisticAverage,
	})
	return err
}

func (tc *cloudwatchTestCtx) deleteAlarms(names ...string) {
	tc.client.DeleteAlarms(tc.ctx, &cloudwatch.DeleteAlarmsInput{
		AlarmNames: names,
	})
}
