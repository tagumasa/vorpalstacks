package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"vorpalstacks-sdk-tests/config"
)

func (tc *cwlogsTestCtx) metricFilterTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "PutMetricFilter_VerifyFields", func() error {
		mfName, cleanupGroup, err := tc.newLogGroupFixture("MFGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		_, err = tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
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

		resp, err := tc.collectAllMetricFilters(mfName, "")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp))
		}
		mf := resp[0]
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
		dmfName, cleanupGroup, err := tc.newLogGroupFixture("DMFGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		if err := tc.putMetricFilter(dmfName, "TestFilter", "ERROR", "ErrorCount", "vorpalstacks/test"); err != nil {
			return fmt.Errorf("put metric filter: %v", err)
		}

		resp, err := tc.collectAllMetricFilters(dmfName, "")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp))
		}
		if *resp[0].FilterName != "TestFilter" {
			return fmt.Errorf("filter name mismatch: got %q", *resp[0].FilterName)
		}
		if *resp[0].FilterPattern != "ERROR" {
			return fmt.Errorf("filter pattern mismatch: got %q", *resp[0].FilterPattern)
		}
		if len(resp[0].MetricTransformations) != 1 {
			return fmt.Errorf("expected 1 transformation, got %d", len(resp[0].MetricTransformations))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeMetricFilters_FilterNamePrefix", func() error {
		fpfName, cleanupGroup, err := tc.newLogGroupFixture("FPFGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		if err := tc.putMetricFilter(fpfName, "PrefixFilterA", "ERROR", "ErrorA", "test"); err != nil {
			return fmt.Errorf("put filter A: %v", err)
		}
		if err := tc.putMetricFilter(fpfName, "PrefixFilterB", "WARN", "WarnB", "test"); err != nil {
			return fmt.Errorf("put filter B: %v", err)
		}

		resp, err := tc.collectAllMetricFilters(fpfName, "PrefixFilterA")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 1 {
			return fmt.Errorf("expected 1 filter with prefix 'PrefixFilterA', got %d", len(resp))
		}
		if *resp[0].FilterName != "PrefixFilterA" {
			return fmt.Errorf("filter name mismatch: got %q", *resp[0].FilterName)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteMetricFilter_Basic", func() error {
		dmfDelName, cleanupGroup, err := tc.newLogGroupFixture("DMFDelGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		if err := tc.putMetricFilter(dmfDelName, "TempFilter", "ERROR", "Err", "test"); err != nil {
			return fmt.Errorf("put metric filter: %v", err)
		}

		_, err = tc.client.DeleteMetricFilter(tc.ctx, &cloudwatchlogs.DeleteMetricFilterInput{
			LogGroupName: aws.String(dmfDelName),
			FilterName:   aws.String("TempFilter"),
		})
		if err != nil {
			return fmt.Errorf("delete metric filter: %v", err)
		}

		resp, err := tc.collectAllMetricFilters(dmfDelName, "")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 0 {
			return fmt.Errorf("expected 0 filters after delete, got %d", len(resp))
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

		// The pattern's shape allows the empty string (the requiredness
		// PutMetricFilter already implements): a present empty pattern
		// matches every message.
		resp, err = tc.client.TestMetricFilter(tc.ctx, &cloudwatchlogs.TestMetricFilterInput{
			FilterPattern:    aws.String(""),
			LogEventMessages: []string{"anything", "at all"},
		})
		if err != nil {
			return fmt.Errorf("empty pattern: %v", err)
		}
		if len(resp.Matches) != 2 {
			return fmt.Errorf("empty pattern matches all, got %d matches", len(resp.Matches))
		}
		return nil
	}))

	// A metric filter must publish a datapoint that GetMetricStatistics can
	// read back: this pins the whole chain (PutLogEvents → filter
	// evaluation → CloudWatch metric invoker → metric store → API read
	// plane) with immediate visibility.
	results = append(results, tc.runner.RunTest("logs", "MetricFilter_PublishesMetric", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("MFDeliver", "stream")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		metricName := tc.uniquePrefix("ErrorCount")
		if err := tc.putMetricFilter(groupName, "ErrorFilter", "ERROR", metricName, "vorpalstacks/test"); err != nil {
			return fmt.Errorf("put metric filter: %v", err)
		}
		defer func() {
			_, _ = tc.client.DeleteMetricFilter(tc.ctx, &cloudwatchlogs.DeleteMetricFilterInput{
				LogGroupName: aws.String(groupName),
				FilterName:   aws.String("ErrorFilter"),
			})
		}()

		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "stream", "ERROR something failed", now-2000); err != nil {
			return fmt.Errorf("put matching event: %v", err)
		}
		if err := tc.putLogEvent(groupName, "stream", "INFO all good", now-1000); err != nil {
			return fmt.Errorf("put non-matching event: %v", err)
		}

		cwCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: tc.runner.endpoint,
			Region:   tc.region,
		})
		if err != nil {
			return fmt.Errorf("cloudwatch config: %v", err)
		}
		cwClient := cloudwatch.NewFromConfig(cwCfg)
		start := time.Now().Add(-5 * time.Minute)
		end := time.Now().Add(1 * time.Minute)
		for i := 0; i < 40; i++ {
			time.Sleep(250 * time.Millisecond)
			out, err := cwClient.GetMetricStatistics(tc.ctx, &cloudwatch.GetMetricStatisticsInput{
				Namespace:  aws.String("vorpalstacks/test"),
				MetricName: aws.String(metricName),
				StartTime:  aws.Time(start),
				EndTime:    aws.Time(end),
				Period:     aws.Int32(300),
				Statistics: []cwTypes.Statistic{cwTypes.StatisticSum},
			})
			if err != nil {
				continue
			}
			for _, dp := range out.Datapoints {
				if dp.Sum != nil && *dp.Sum >= 1 {
					return nil
				}
			}
		}
		return fmt.Errorf("metric %s not visible via GetMetricStatistics after 10s", metricName)
	}))

	// Log-group names may contain '#': the filters of a group whose name
	// extends another with '#' must stay invisible to the shorter group's
	// listing, in both directions.
	results = append(results, tc.runner.RunTest("logs", "DescribeMetricFilters_HashGroupNameIsolation", func() error {
		base := tc.uniquePrefix("HashGroup")
		nested := base + "#inner"
		if err := tc.createLogGroup(base); err != nil {
			return fmt.Errorf("create base: %v", err)
		}
		defer tc.deleteLogGroup(base)
		if err := tc.createLogGroup(nested); err != nil {
			return fmt.Errorf("create nested: %v", err)
		}
		defer tc.deleteLogGroup(nested)

		if err := tc.putMetricFilter(base, "RootFilter", "ERROR", "HashRoot", "test"); err != nil {
			return fmt.Errorf("put base filter: %v", err)
		}
		defer func() {
			_, _ = tc.client.DeleteMetricFilter(tc.ctx, &cloudwatchlogs.DeleteMetricFilterInput{
				LogGroupName: aws.String(base),
				FilterName:   aws.String("RootFilter"),
			})
		}()
		if err := tc.putMetricFilter(nested, "InnerFilter", "WARN", "HashInner", "test"); err != nil {
			return fmt.Errorf("put nested filter: %v", err)
		}
		defer func() {
			_, _ = tc.client.DeleteMetricFilter(tc.ctx, &cloudwatchlogs.DeleteMetricFilterInput{
				LogGroupName: aws.String(nested),
				FilterName:   aws.String("InnerFilter"),
			})
		}()

		baseResp, err := tc.collectAllMetricFilters(base, "")
		if err != nil {
			return fmt.Errorf("describe base: %v", err)
		}
		if len(baseResp) != 1 || aws.ToString(baseResp[0].FilterName) != "RootFilter" {
			return fmt.Errorf("base listing expected only RootFilter, got %d filters", len(baseResp))
		}

		nestedResp, err := tc.collectAllMetricFilters(nested, "")
		if err != nil {
			return fmt.Errorf("describe nested: %v", err)
		}
		if len(nestedResp) != 1 || aws.ToString(nestedResp[0].FilterName) != "InnerFilter" {
			return fmt.Errorf("nested listing expected only InnerFilter, got %d filters", len(nestedResp))
		}
		return nil
	}))

	return results
}
