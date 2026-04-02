package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCloudWatchLogsTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "logs",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cloudwatchlogs.NewFromConfig(cfg)
	ctx := context.Background()

	logGroupName := fmt.Sprintf("TestLogGroup-%d", time.Now().UnixNano())
	logStreamName := fmt.Sprintf("TestLogStream-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("logs", "CreateLogGroup", func() error {
		resp, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeLogGroups", func() error {
		resp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{})
		if err != nil {
			return err
		}
		if resp.LogGroups == nil {
			return fmt.Errorf("log groups list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeLogStreams", func() error {
		resp, err := client.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if resp.LogStreams == nil {
			return fmt.Errorf("log streams list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "CreateLogStream", func() error {
		resp, err := client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "PutLogEvents", func() error {
		resp, err := client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
			LogEvents: []types.InputLogEvent{
				{
					Message:   aws.String("Test log message"),
					Timestamp: aws.Int64(time.Now().UnixMilli()),
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

	results = append(results, r.RunTest("logs", "GetLogEvents", func() error {
		resp, err := client.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			return err
		}
		if resp.Events == nil {
			return fmt.Errorf("events list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "FilterLogEvents", func() error {
		resp, err := client.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if resp.Events == nil {
			return fmt.Errorf("events list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "PutRetentionPolicy", func() error {
		resp, err := client.PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(logGroupName),
			RetentionInDays: aws.Int32(7),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DeleteLogStream", func() error {
		resp, err := client.DeleteLogStream(ctx, &cloudwatchlogs.DeleteLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DeleteLogGroup", func() error {
		resp, err := client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("logs", "CreateLogGroup_Duplicate", func() error {
		dupGroupName := fmt.Sprintf("DupLogGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dupGroupName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dupGroupName),
		})

		_, err = client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dupGroupName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate log group")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DeleteLogGroup_NonExistent", func() error {
		_, err := client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String("nonexistent-log-group-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent log group")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "PutLogEvents_GetLogEvents_Roundtrip", func() error {
		rtGroupName := fmt.Sprintf("RTLogGroup-%d", time.Now().UnixNano())
		rtStreamName := fmt.Sprintf("RTLogStream-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(rtGroupName),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(rtGroupName),
		})

		_, err = client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(rtGroupName),
			LogStreamName: aws.String(rtStreamName),
		})
		if err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		testMessage := "roundtrip-log-message-verify-12345"
		ts := time.Now().UnixMilli()
		_, err = client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(rtGroupName),
			LogStreamName: aws.String(rtStreamName),
			LogEvents: []types.InputLogEvent{
				{
					Message:   aws.String(testMessage),
					Timestamp: aws.Int64(ts),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(rtGroupName),
			LogStreamName: aws.String(rtStreamName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.Events) == 0 {
			return fmt.Errorf("no events returned")
		}
		if resp.Events[0].Message == nil || *resp.Events[0].Message != testMessage {
			return fmt.Errorf("message mismatch: got %q, want %q", aws.ToString(resp.Events[0].Message), testMessage)
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeLogGroups_ContainsCreated", func() error {
		dlgName := fmt.Sprintf("DLGGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dlgName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dlgName),
		})

		resp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(dlgName),
		})
		if err != nil {
			return err
		}
		if len(resp.LogGroups) != 1 {
			return fmt.Errorf("expected 1 log group, got %d", len(resp.LogGroups))
		}
		if *resp.LogGroups[0].LogGroupName != dlgName {
			return fmt.Errorf("log group name mismatch")
		}
		return nil
	}))

	// === TAGGING TESTS ===

	results = append(results, r.RunTest("logs", "TagResource_Basic", func() error {
		tgName := fmt.Sprintf("TagGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(tgName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(tgName),
		})

		descResp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(tgName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) == 0 || descResp.LogGroups[0].LogGroupArn == nil {
			return fmt.Errorf("log group arn not found")
		}

		_, err = client.TagResource(ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: descResp.LogGroups[0].LogGroupArn,
			Tags: map[string]string{
				"Environment": "test",
				"Team":        "vorpalstacks",
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "ListTagsForResource_Basic", func() error {
		ltName := fmt.Sprintf("ListTagGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(ltName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(ltName),
		})

		descResp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(ltName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) == 0 || descResp.LogGroups[0].LogGroupArn == nil {
			return fmt.Errorf("log group arn not found")
		}

		client.TagResource(ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: descResp.LogGroups[0].LogGroupArn,
			Tags: map[string]string{
				"Key1": "Value1",
				"Key2": "Value2",
			},
		})

		tagResp, err := client.ListTagsForResource(ctx, &cloudwatchlogs.ListTagsForResourceInput{
			ResourceArn: descResp.LogGroups[0].LogGroupArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if tagResp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if tagResp.Tags["Key1"] != "Value1" {
			return fmt.Errorf("Key1 mismatch: got %q", tagResp.Tags["Key1"])
		}
		if tagResp.Tags["Key2"] != "Value2" {
			return fmt.Errorf("Key2 mismatch: got %q", tagResp.Tags["Key2"])
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "UntagResource_Basic", func() error {
		utName := fmt.Sprintf("UntagGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(utName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(utName),
		})

		descResp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(utName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) == 0 || descResp.LogGroups[0].LogGroupArn == nil {
			return fmt.Errorf("log group arn not found")
		}

		client.TagResource(ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: descResp.LogGroups[0].LogGroupArn,
			Tags: map[string]string{
				"RemoveMe":  "yes",
				"KeepMe":    "no",
				"KeepMeToo": "also-no",
			},
		})

		_, err = client.UntagResource(ctx, &cloudwatchlogs.UntagResourceInput{
			ResourceArn: descResp.LogGroups[0].LogGroupArn,
			TagKeys:     []string{"RemoveMe"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		tagResp, err := client.ListTagsForResource(ctx, &cloudwatchlogs.ListTagsForResourceInput{
			ResourceArn: descResp.LogGroups[0].LogGroupArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if _, ok := tagResp.Tags["RemoveMe"]; ok {
			return fmt.Errorf("RemoveMe tag should have been removed")
		}
		if tagResp.Tags["KeepMe"] != "no" {
			return fmt.Errorf("KeepMe tag should still exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "TagLogGroup_Basic", func() error {
		tlgName := fmt.Sprintf("TagLGGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(tlgName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(tlgName),
		})

		_, err = client.TagLogGroup(ctx, &cloudwatchlogs.TagLogGroupInput{
			LogGroupName: aws.String(tlgName),
			Tags: map[string]string{
				"DeprecatedTag": "yes",
			},
		})
		if err != nil {
			return fmt.Errorf("tag log group: %v", err)
		}

		tagResp, err := client.ListTagsLogGroup(ctx, &cloudwatchlogs.ListTagsLogGroupInput{
			LogGroupName: aws.String(tlgName),
		})
		if err != nil {
			return fmt.Errorf("list tags log group: %v", err)
		}
		if tagResp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if tagResp.Tags["DeprecatedTag"] != "yes" {
			return fmt.Errorf("DeprecatedTag mismatch: got %q", tagResp.Tags["DeprecatedTag"])
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "CreateLogGroup_WithTags", func() error {
		cwtName := fmt.Sprintf("CWTagGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(cwtName),
			Tags: map[string]string{
				"CreatedBy": "sdk-test",
				"CreatedAt": "2026-04-02",
				"Automated": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create with tags: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(cwtName),
		})

		tagResp, err := client.ListTagsLogGroup(ctx, &cloudwatchlogs.ListTagsLogGroupInput{
			LogGroupName: aws.String(cwtName),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tagResp.Tags) != 3 {
			return fmt.Errorf("expected 3 tags, got %d", len(tagResp.Tags))
		}
		if tagResp.Tags["CreatedBy"] != "sdk-test" {
			return fmt.Errorf("CreatedBy mismatch")
		}
		return nil
	}))

	// === RETENTION TESTS ===

	results = append(results, r.RunTest("logs", "DeleteRetentionPolicy_Basic", func() error {
		drName := fmt.Sprintf("DelRetGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(drName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(drName),
		})

		_, err = client.PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(drName),
			RetentionInDays: aws.Int32(14),
		})
		if err != nil {
			return fmt.Errorf("put retention: %v", err)
		}

		_, err = client.DeleteRetentionPolicy(ctx, &cloudwatchlogs.DeleteRetentionPolicyInput{
			LogGroupName: aws.String(drName),
		})
		if err != nil {
			return fmt.Errorf("delete retention: %v", err)
		}

		descResp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(drName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) == 0 {
			return fmt.Errorf("log group not found")
		}
		if descResp.LogGroups[0].RetentionInDays != nil && *descResp.LogGroups[0].RetentionInDays != 0 {
			return fmt.Errorf("expected retention 0 after delete, got %d", *descResp.LogGroups[0].RetentionInDays)
		}
		return nil
	}))

	// === METRIC FILTER TESTS ===

	results = append(results, r.RunTest("logs", "PutMetricFilter_Basic", func() error {
		mfName := fmt.Sprintf("MFGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(mfName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(mfName),
		})

		_, err = client.PutMetricFilter(ctx, &cloudwatchlogs.PutMetricFilterInput{
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
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeMetricFilters_Basic", func() error {
		dmfName := fmt.Sprintf("DMFGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dmfName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dmfName),
		})

		_, err = client.PutMetricFilter(ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(dmfName),
			FilterName:    aws.String("TestFilter"),
			FilterPattern: aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{
				{
					MetricName:      aws.String("ErrorCount"),
					MetricNamespace: aws.String("vorpalstacks/test"),
					MetricValue:     aws.String("1"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put metric filter: %v", err)
		}

		resp, err := client.DescribeMetricFilters(ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
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

	results = append(results, r.RunTest("logs", "DescribeMetricFilters_FilterNamePrefix", func() error {
		fpfName := fmt.Sprintf("FPFGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(fpfName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(fpfName),
		})

		client.PutMetricFilter(ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(fpfName),
			FilterName:    aws.String("PrefixFilterA"),
			FilterPattern: aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{
				{
					MetricName:      aws.String("ErrorA"),
					MetricNamespace: aws.String("test"),
					MetricValue:     aws.String("1"),
				},
			},
		})
		client.PutMetricFilter(ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(fpfName),
			FilterName:    aws.String("PrefixFilterB"),
			FilterPattern: aws.String("WARN"),
			MetricTransformations: []types.MetricTransformation{
				{
					MetricName:      aws.String("WarnB"),
					MetricNamespace: aws.String("test"),
					MetricValue:     aws.String("1"),
				},
			},
		})

		resp, err := client.DescribeMetricFilters(ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
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

	results = append(results, r.RunTest("logs", "DeleteMetricFilter_Basic", func() error {
		dmfDelName := fmt.Sprintf("DMFDelGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dmfDelName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dmfDelName),
		})

		client.PutMetricFilter(ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(dmfDelName),
			FilterName:    aws.String("TempFilter"),
			FilterPattern: aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{
				{
					MetricName:      aws.String("Err"),
					MetricNamespace: aws.String("test"),
					MetricValue:     aws.String("1"),
				},
			},
		})

		_, err = client.DeleteMetricFilter(ctx, &cloudwatchlogs.DeleteMetricFilterInput{
			LogGroupName: aws.String(dmfDelName),
			FilterName:   aws.String("TempFilter"),
		})
		if err != nil {
			return fmt.Errorf("delete metric filter: %v", err)
		}

		resp, err := client.DescribeMetricFilters(ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
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

	results = append(results, r.RunTest("logs", "TestMetricFilter_Basic", func() error {
		resp, err := client.TestMetricFilter(ctx, &cloudwatchlogs.TestMetricFilterInput{
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

	// === SUBSCRIPTION FILTER TESTS ===

	results = append(results, r.RunTest("logs", "PutSubscriptionFilter_Basic", func() error {
		sfName := fmt.Sprintf("SFGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(sfName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(sfName),
		})

		_, err = client.PutSubscriptionFilter(ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(sfName),
			FilterName:     aws.String("TestSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:test-func"),
			RoleArn:        aws.String("arn:aws:iam::000000000000:role/test-role"),
			Distribution:   types.DistributionByLogStream,
		})
		if err != nil {
			return fmt.Errorf("put subscription filter: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeSubscriptionFilters_Basic", func() error {
		dsfName := fmt.Sprintf("DSFGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dsfName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dsfName),
		})

		client.PutSubscriptionFilter(ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(dsfName),
			FilterName:     aws.String("DescSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:test"),
			RoleArn:        aws.String("arn:aws:iam::000000000000:role/test"),
		})

		resp, err := client.DescribeSubscriptionFilters(ctx, &cloudwatchlogs.DescribeSubscriptionFiltersInput{
			LogGroupName: aws.String(dsfName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.SubscriptionFilters) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp.SubscriptionFilters))
		}
		if *resp.SubscriptionFilters[0].FilterName != "DescSub" {
			return fmt.Errorf("filter name mismatch: got %q", *resp.SubscriptionFilters[0].FilterName)
		}
		if *resp.SubscriptionFilters[0].DestinationArn != "arn:aws:lambda:us-east-1:000000000000:function:test" {
			return fmt.Errorf("destination arn mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DeleteSubscriptionFilter_Basic", func() error {
		delSFName := fmt.Sprintf("DelSFGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(delSFName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(delSFName),
		})

		client.PutSubscriptionFilter(ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(delSFName),
			FilterName:     aws.String("DelSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:test"),
			RoleArn:        aws.String("arn:aws:iam::000000000000:role/test"),
		})

		_, err = client.DeleteSubscriptionFilter(ctx, &cloudwatchlogs.DeleteSubscriptionFilterInput{
			LogGroupName: aws.String(delSFName),
			FilterName:   aws.String("DelSub"),
		})
		if err != nil {
			return fmt.Errorf("delete subscription filter: %v", err)
		}

		resp, err := client.DescribeSubscriptionFilters(ctx, &cloudwatchlogs.DescribeSubscriptionFiltersInput{
			LogGroupName: aws.String(delSFName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.SubscriptionFilters) != 0 {
			return fmt.Errorf("expected 0 filters after delete, got %d", len(resp.SubscriptionFilters))
		}
		return nil
	}))

	// === DESTINATION TESTS ===

	results = append(results, r.RunTest("logs", "PutDestination_Basic", func() error {
		_, err := client.PutDestination(ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(fmt.Sprintf("TestDest-%d", time.Now().UnixNano())),
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/dest-role"),
			TargetArn:       aws.String("arn:aws:kinesis:us-east-1:000000000000:stream/test-stream"),
		})
		if err != nil {
			return fmt.Errorf("put destination: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeDestinations_Basic", func() error {
		ddName := fmt.Sprintf("DescDest-%d", time.Now().UnixNano())
		_, err := client.PutDestination(ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(ddName),
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/dd-role"),
			TargetArn:       aws.String("arn:aws:kinesis:us-east-1:000000000000:stream/dd-stream"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteDestination(ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(ddName),
		})

		resp, err := client.DescribeDestinations(ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(ddName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) != 1 {
			return fmt.Errorf("expected 1 destination, got %d", len(resp.Destinations))
		}
		if *resp.Destinations[0].DestinationName != ddName {
			return fmt.Errorf("name mismatch: got %q", *resp.Destinations[0].DestinationName)
		}
		if resp.Destinations[0].Arn == nil || *resp.Destinations[0].Arn == "" {
			return fmt.Errorf("ARN is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "PutDestinationPolicy_Basic", func() error {
		pdpName := fmt.Sprintf("PdpDest-%d", time.Now().UnixNano())
		_, err := client.PutDestination(ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(pdpName),
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/pdp-role"),
			TargetArn:       aws.String("arn:aws:kinesis:us-east-1:000000000000:stream/pdp-stream"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteDestination(ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(pdpName),
		})

		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"logs:PutSubscriptionFilter"}]}`
		_, err = client.PutDestinationPolicy(ctx, &cloudwatchlogs.PutDestinationPolicyInput{
			DestinationName: aws.String(pdpName),
			AccessPolicy:    aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}

		resp, err := client.DescribeDestinations(ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(pdpName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) == 0 {
			return fmt.Errorf("destination not found")
		}
		if resp.Destinations[0].AccessPolicy == nil || *resp.Destinations[0].AccessPolicy != policy {
			return fmt.Errorf("access policy mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "PutDestination_UpdateInPlace", func() error {
		udpName := fmt.Sprintf("UdpDest-%d", time.Now().UnixNano())
		_, err := client.PutDestination(ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(udpName),
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/original-role"),
			TargetArn:       aws.String("arn:aws:kinesis:us-east-1:000000000000:stream/original"),
		})
		if err != nil {
			return fmt.Errorf("initial put: %v", err)
		}
		defer client.DeleteDestination(ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(udpName),
		})

		_, err = client.PutDestination(ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(udpName),
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/updated-role"),
			TargetArn:       aws.String("arn:aws:kinesis:us-east-1:000000000000:stream/updated"),
		})
		if err != nil {
			return fmt.Errorf("update put: %v", err)
		}

		resp, err := client.DescribeDestinations(ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(udpName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) == 0 {
			return fmt.Errorf("destination not found")
		}
		if *resp.Destinations[0].RoleArn != "arn:aws:iam::000000000000:role/updated-role" {
			return fmt.Errorf("roleArn not updated: got %q", *resp.Destinations[0].RoleArn)
		}
		if *resp.Destinations[0].TargetArn != "arn:aws:kinesis:us-east-1:000000000000:stream/updated" {
			return fmt.Errorf("targetArn not updated: got %q", *resp.Destinations[0].TargetArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DeleteDestination_Basic", func() error {
		ddelName := fmt.Sprintf("DelDest-%d", time.Now().UnixNano())
		_, err := client.PutDestination(ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(ddelName),
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/ddel-role"),
			TargetArn:       aws.String("arn:aws:kinesis:us-east-1:000000000000:stream/ddel-stream"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = client.DeleteDestination(ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(ddelName),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		resp, err := client.DescribeDestinations(ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(ddelName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) != 0 {
			return fmt.Errorf("expected 0 destinations after delete, got %d", len(resp.Destinations))
		}
		return nil
	}))

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("logs", "DescribeLogGroups_Pagination", func() error {
		pgPrefix := fmt.Sprintf("PagGroup-%d-", time.Now().UnixNano())
		groupNames := []string{
			pgPrefix + "alpha",
			pgPrefix + "beta",
			pgPrefix + "gamma",
		}
		for _, name := range groupNames {
			client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
				LogGroupName: aws.String(name),
			})
			defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
				LogGroupName: aws.String(name),
			})
		}

		var allGroups []types.LogGroup
		var nextToken *string
		for {
			resp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
				LogGroupNamePrefix: aws.String(pgPrefix),
				Limit:              aws.Int32(2),
				NextToken:          nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe page: %v", err)
			}
			allGroups = append(allGroups, resp.LogGroups...)
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}

		if len(allGroups) != 3 {
			return fmt.Errorf("expected 3 groups across pages, got %d", len(allGroups))
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeLogStreams_Pagination", func() error {
		psName := fmt.Sprintf("PagStreamGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(psName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(psName),
		})

		streamNames := []string{psName + "-s1", psName + "-s2", psName + "-s3"}
		for _, sn := range streamNames {
			client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
				LogGroupName:  aws.String(psName),
				LogStreamName: aws.String(sn),
			})
		}

		var allStreams []types.LogStream
		var nextToken *string
		for {
			resp, err := client.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
				LogGroupName: aws.String(psName),
				Limit:        aws.Int32(2),
				NextToken:    nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe page: %v", err)
			}
			allStreams = append(allStreams, resp.LogStreams...)
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}

		if len(allStreams) != 3 {
			return fmt.Errorf("expected 3 streams across pages, got %d", len(allStreams))
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DescribeLogStreams_NamePrefix", func() error {
		npName := fmt.Sprintf("NpStreamGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(npName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(npName),
		})

		client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(npName),
			LogStreamName: aws.String("app-server-1"),
		})
		client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(npName),
			LogStreamName: aws.String("app-server-2"),
		})
		client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(npName),
			LogStreamName: aws.String("db-server-1"),
		})

		resp, err := client.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(npName),
			LogStreamNamePrefix: aws.String("app-"),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.LogStreams) != 2 {
			return fmt.Errorf("expected 2 streams with prefix 'app-', got %d", len(resp.LogStreams))
		}
		for _, ls := range resp.LogStreams {
			if !strings.HasPrefix(*ls.LogStreamName, "app-") {
				return fmt.Errorf("unexpected stream: %q", *ls.LogStreamName)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "DeleteLogStream_NonExistent", func() error {
		dlsName := fmt.Sprintf("DlsGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dlsName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dlsName),
		})

		_, err = client.DeleteLogStream(ctx, &cloudwatchlogs.DeleteLogStreamInput{
			LogGroupName:  aws.String(dlsName),
			LogStreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent log stream")
		}
		return nil
	}))

	// === LOG EVENTS EXTENDED TESTS ===

	results = append(results, r.RunTest("logs", "PutLogEvents_MultipleEvents", func() error {
		meName := fmt.Sprintf("MEGroup-%d", time.Now().UnixNano())
		meStream := fmt.Sprintf("MEStream-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(meName),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(meName),
		})

		_, err = client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(meName),
			LogStreamName: aws.String(meStream),
		})
		if err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		_, err = client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(meName),
			LogStreamName: aws.String(meStream),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("event-1"), Timestamp: aws.Int64(ts)},
				{Message: aws.String("event-2"), Timestamp: aws.Int64(ts + 1)},
				{Message: aws.String("event-3"), Timestamp: aws.Int64(ts + 2)},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(meName),
			LogStreamName: aws.String(meStream),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.Events) != 3 {
			return fmt.Errorf("expected 3 events, got %d", len(resp.Events))
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "GetLogEvents_StartFromHead", func() error {
		sfhName := fmt.Sprintf("SFHGroup-%d", time.Now().UnixNano())
		sfhStream := fmt.Sprintf("SFHStream-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(sfhName),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(sfhName),
		})

		_, err = client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(sfhName),
			LogStreamName: aws.String(sfhStream),
		})
		if err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(sfhName),
			LogStreamName: aws.String(sfhStream),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("first-event"), Timestamp: aws.Int64(ts)},
				{Message: aws.String("second-event"), Timestamp: aws.Int64(ts + 1)},
			},
		})

		resp, err := client.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(sfhName),
			LogStreamName: aws.String(sfhStream),
			StartFromHead: aws.Bool(true),
			Limit:         aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.Events) == 0 {
			return fmt.Errorf("no events returned")
		}
		if *resp.Events[0].Message != "first-event" {
			return fmt.Errorf("expected first-event when StartFromHead=true, got %q", *resp.Events[0].Message)
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "FilterLogEvents_WithFilterPattern", func() error {
		fepName := fmt.Sprintf("FEPGroup-%d", time.Now().UnixNano())
		fepStream := fmt.Sprintf("FEPStream-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(fepName),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(fepName),
		})

		_, err = client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(fepName),
			LogStreamName: aws.String(fepStream),
		})
		if err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(fepName),
			LogStreamName: aws.String(fepStream),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("ERROR disk full"), Timestamp: aws.Int64(ts)},
				{Message: aws.String("INFO started"), Timestamp: aws.Int64(ts + 1)},
				{Message: aws.String("ERROR network timeout"), Timestamp: aws.Int64(ts + 2)},
			},
		})

		resp, err := client.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:  aws.String(fepName),
			FilterPattern: aws.String("ERROR"),
		})
		if err != nil {
			return fmt.Errorf("filter: %v", err)
		}
		if len(resp.Events) != 2 {
			return fmt.Errorf("expected 2 ERROR events, got %d", len(resp.Events))
		}
		for _, e := range resp.Events {
			if !strings.Contains(*e.Message, "ERROR") {
				return fmt.Errorf("non-ERROR event in results: %q", *e.Message)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "FilterLogEvents_WithLogStreamNames", func() error {
		flsName := fmt.Sprintf("FLSGroup-%d", time.Now().UnixNano())
		flsStream1 := fmt.Sprintf("FLSStream1-%d", time.Now().UnixNano())
		flsStream2 := fmt.Sprintf("FLSStream2-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(flsName),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(flsName),
		})

		client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream1),
		})
		client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream2),
		})

		ts := time.Now().UnixMilli()
		client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream1),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-stream-1"), Timestamp: aws.Int64(ts)},
			},
		})
		client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream2),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-stream-2"), Timestamp: aws.Int64(ts + 1)},
			},
		})

		resp, err := client.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:   aws.String(flsName),
			LogStreamNames: []string{flsStream1},
		})
		if err != nil {
			return fmt.Errorf("filter: %v", err)
		}
		if len(resp.Events) != 1 {
			return fmt.Errorf("expected 1 event from stream1, got %d", len(resp.Events))
		}
		if resp.Events[0].LogStreamName == nil || *resp.Events[0].LogStreamName != flsStream1 {
			return fmt.Errorf("logStreamName mismatch: got %q", aws.ToString(resp.Events[0].LogStreamName))
		}
		return nil
	}))

	results = append(results, r.RunTest("logs", "MetricFilterCount_Tracked", func() error {
		mfcName := fmt.Sprintf("MFCGroup-%d", time.Now().UnixNano())
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(mfcName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(mfcName),
		})

		descResp, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(mfcName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) == 0 {
			return fmt.Errorf("log group not found")
		}
		if descResp.LogGroups[0].MetricFilterCount != nil && *descResp.LogGroups[0].MetricFilterCount != 0 {
			return fmt.Errorf("expected 0 filters, got %d", *descResp.LogGroups[0].MetricFilterCount)
		}

		client.PutMetricFilter(ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(mfcName),
			FilterName:    aws.String("CountFilter1"),
			FilterPattern: aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{
				{MetricName: aws.String("E"), MetricNamespace: aws.String("NS"), MetricValue: aws.String("1")},
			},
		})

		descResp2, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(mfcName),
		})
		if err != nil {
			return fmt.Errorf("describe after: %v", err)
		}
		if descResp2.LogGroups[0].MetricFilterCount == nil || *descResp2.LogGroups[0].MetricFilterCount != 1 {
			return fmt.Errorf("expected 1 filter, got %d", aws.ToInt32(descResp2.LogGroups[0].MetricFilterCount))
		}

		client.DeleteMetricFilter(ctx, &cloudwatchlogs.DeleteMetricFilterInput{
			LogGroupName: aws.String(mfcName),
			FilterName:   aws.String("CountFilter1"),
		})

		descResp3, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(mfcName),
		})
		if err != nil {
			return fmt.Errorf("describe after delete: %v", err)
		}
		if descResp3.LogGroups[0].MetricFilterCount == nil || *descResp3.LogGroups[0].MetricFilterCount != 0 {
			return fmt.Errorf("expected 0 filters after delete, got %d", aws.ToInt32(descResp3.LogGroups[0].MetricFilterCount))
		}
		return nil
	}))

	return results
}
