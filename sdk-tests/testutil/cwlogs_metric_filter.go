package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func (tc *cwlogsTestCtx) metricFilterTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "PutMetricFilter_VerifyFields", func() error {
		mfName := tc.uniquePrefix("MFGroup")
		if err := tc.createLogGroup(mfName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(mfName)

		_, err := tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(mfName),
			FilterName:    aws.String("ErrorFilter"),
			FilterPattern: aws.String("[ip, user, timestamp, request, status_code=*, bytes=*]"),
			MetricTransformations: []types.MetricTransformation{
				{
					MetricName:      aws.String("ErrorCount"),
					MetricNamespace: aws.String("vorpalstacks/test"),
					MetricValue:     aws.String("1"),
					DefaultValue:    aws.Float64(0.0),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put metric filter: %v", err)
		}

		resp, err := tc.client.DescribeMetricFilters(tc.ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName: aws.String(mfName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.MetricFilters) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp.MetricFilters))
		}
		mf := resp.MetricFilters[0]
		if mf.FilterName == nil || *mf.FilterName != "ErrorFilter" {
			return fmt.Errorf("filterName mismatch: got %q", aws.ToString(mf.FilterName))
		}
		if mf.FilterPattern == nil || *mf.FilterPattern != "[ip, user, timestamp, request, status_code=*, bytes=*]" {
			return fmt.Errorf("filterPattern mismatch: got %q", aws.ToString(mf.FilterPattern))
		}
		if len(mf.MetricTransformations) != 1 {
			return fmt.Errorf("expected 1 transformation, got %d", len(mf.MetricTransformations))
		}
		mt := mf.MetricTransformations[0]
		if mt.MetricName == nil || *mt.MetricName != "ErrorCount" {
			return fmt.Errorf("metricName mismatch: got %q", aws.ToString(mt.MetricName))
		}
		if mt.MetricNamespace == nil || *mt.MetricNamespace != "vorpalstacks/test" {
			return fmt.Errorf("metricNamespace mismatch: got %q", aws.ToString(mt.MetricNamespace))
		}
		if mt.MetricValue == nil || *mt.MetricValue != "1" {
			return fmt.Errorf("metricValue mismatch: got %q", aws.ToString(mt.MetricValue))
		}
		if mt.DefaultValue == nil || *mt.DefaultValue != 0.0 {
			return fmt.Errorf("defaultValue mismatch: got %v", aws.ToFloat64(mt.DefaultValue))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeMetricFilters_Basic", func() error {
		dmfName := tc.uniquePrefix("DMFGroup")
		if err := tc.createLogGroup(dmfName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(dmfName)

		if err := tc.putMetricFilter(dmfName, "TestFilter", "ERROR", "ErrorCount", "vorpalstacks/test"); err != nil {
			return fmt.Errorf("put metric filter: %v", err)
		}

		resp, err := tc.client.DescribeMetricFilters(tc.ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName: aws.String(dmfName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.MetricFilters) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp.MetricFilters))
		}
		if *resp.MetricFilters[0].FilterName != "TestFilter" {
			return fmt.Errorf("filter name mismatch: got %q", *resp.MetricFilters[0].FilterName)
		}
		if *resp.MetricFilters[0].FilterPattern != "ERROR" {
			return fmt.Errorf("filter pattern mismatch: got %q", *resp.MetricFilters[0].FilterPattern)
		}
		if len(resp.MetricFilters[0].MetricTransformations) != 1 {
			return fmt.Errorf("expected 1 transformation, got %d", len(resp.MetricFilters[0].MetricTransformations))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeMetricFilters_FilterNamePrefix", func() error {
		fpfName := tc.uniquePrefix("FPFGroup")
		if err := tc.createLogGroup(fpfName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(fpfName)

		tc.putMetricFilter(fpfName, "PrefixFilterA", "ERROR", "ErrorA", "test")
		tc.putMetricFilter(fpfName, "PrefixFilterB", "WARN", "WarnB", "test")

		resp, err := tc.client.DescribeMetricFilters(tc.ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName:     aws.String(fpfName),
			FilterNamePrefix: aws.String("PrefixFilterA"),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.MetricFilters) != 1 {
			return fmt.Errorf("expected 1 filter with prefix 'PrefixFilterA', got %d", len(resp.MetricFilters))
		}
		if *resp.MetricFilters[0].FilterName != "PrefixFilterA" {
			return fmt.Errorf("filter name mismatch: got %q", *resp.MetricFilters[0].FilterName)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteMetricFilter_Basic", func() error {
		dmfDelName := tc.uniquePrefix("DMFDelGroup")
		if err := tc.createLogGroup(dmfDelName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(dmfDelName)

		tc.putMetricFilter(dmfDelName, "TempFilter", "ERROR", "Err", "test")

		_, err := tc.client.DeleteMetricFilter(tc.ctx, &cloudwatchlogs.DeleteMetricFilterInput{
			LogGroupName: aws.String(dmfDelName),
			FilterName:   aws.String("TempFilter"),
		})
		if err != nil {
			return fmt.Errorf("delete metric filter: %v", err)
		}

		resp, err := tc.client.DescribeMetricFilters(tc.ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName: aws.String(dmfDelName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.MetricFilters) != 0 {
			return fmt.Errorf("expected 0 filters after delete, got %d", len(resp.MetricFilters))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "TestMetricFilter_Basic", func() error {
		resp, err := tc.client.TestMetricFilter(tc.ctx, &cloudwatchlogs.TestMetricFilterInput{
			FilterPattern: aws.String("ERROR"),
			LogEventMessages: []string{
				"ERROR something went wrong",
				"INFO all good",
				"[ERROR] critical failure",
			},
		})
		if err != nil {
			return fmt.Errorf("test metric filter: %v", err)
		}
		if resp.Matches == nil {
			return fmt.Errorf("matches is nil")
		}
		if len(resp.Matches) != 2 {
			return fmt.Errorf("expected 2 matches, got %d", len(resp.Matches))
		}
		for _, m := range resp.Matches {
			if m.EventMessage == nil || !strings.Contains(*m.EventMessage, "ERROR") {
				return fmt.Errorf("unexpected match: %q", aws.ToString(m.EventMessage))
			}
		}
		return nil
	}))

	return results
}
