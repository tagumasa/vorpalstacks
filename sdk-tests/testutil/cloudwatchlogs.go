package testutil

import (
	"context"
	"fmt"
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
		_, err := client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		return err
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
		_, err := client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		return err
	}))

	results = append(results, r.RunTest("logs", "PutLogEvents", func() error {
		_, err := client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
			LogEvents: []types.InputLogEvent{
				{
					Message:   aws.String("Test log message"),
					Timestamp: aws.Int64(time.Now().UnixMilli()),
				},
			},
		})
		return err
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
		_, err := client.PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(logGroupName),
			RetentionInDays: aws.Int32(7),
		})
		return err
	}))

	results = append(results, r.RunTest("logs", "DeleteLogStream", func() error {
		_, err := client.DeleteLogStream(ctx, &cloudwatchlogs.DeleteLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		return err
	}))

	results = append(results, r.RunTest("logs", "DeleteLogGroup", func() error {
		_, err := client.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		return err
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

	return results
}
