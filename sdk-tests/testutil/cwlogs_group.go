package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func (tc *cwlogsTestCtx) basicTests() []TestResult {
	var results []TestResult

	logGroupName := tc.uniquePrefix("TestLogGroup")
	logStreamName := tc.uniquePrefix("TestLogStream")

	results = append(results, tc.runner.RunTest("logs", "CreateLogGroup_VerifyFields", func() error {
		resp, err := tc.client.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}

		descResp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(logGroupName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) != 1 {
			return fmt.Errorf("expected 1 log group, got %d", len(descResp.LogGroups))
		}
		lg := descResp.LogGroups[0]
		if lg.LogGroupName == nil || *lg.LogGroupName != logGroupName {
			return fmt.Errorf("logGroupName mismatch: got %q", aws.ToString(lg.LogGroupName))
		}
		if lg.Arn == nil || !strings.Contains(*lg.Arn, logGroupName) {
			return fmt.Errorf("arn missing or does not contain group name: %q", aws.ToString(lg.Arn))
		}
		if lg.CreationTime == nil || *lg.CreationTime == 0 {
			return fmt.Errorf("creationTime is zero or nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeLogGroups_ContainsCreated", func() error {
		resp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if len(resp.LogGroups) != 1 {
			return fmt.Errorf("expected 1 log group, got %d", len(resp.LogGroups))
		}
		lg := resp.LogGroups[0]
		if lg.LogGroupName == nil || *lg.LogGroupName != logGroupName {
			return fmt.Errorf("logGroupName mismatch: got %q, want %q", aws.ToString(lg.LogGroupName), logGroupName)
		}
		if lg.Arn == nil || *lg.Arn == "" {
			return fmt.Errorf("arn is nil or empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "CreateLogStream_VerifyFields", func() error {
		_, err := tc.client.CreateLogStream(tc.ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.client.DescribeLogStreams(tc.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(logGroupName),
			LogStreamNamePrefix: aws.String(logStreamName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogStreams) != 1 {
			return fmt.Errorf("expected 1 log stream, got %d", len(descResp.LogStreams))
		}
		ls := descResp.LogStreams[0]
		if ls.LogStreamName == nil || *ls.LogStreamName != logStreamName {
			return fmt.Errorf("logStreamName mismatch: got %q", aws.ToString(ls.LogStreamName))
		}
		if ls.Arn == nil || !strings.Contains(*ls.Arn, logStreamName) {
			return fmt.Errorf("arn missing or does not contain stream name: %q", aws.ToString(ls.Arn))
		}
		if ls.CreationTime == nil || *ls.CreationTime == 0 {
			return fmt.Errorf("creationTime is zero or nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutLogEvents_VerifyNextToken", func() error {
		resp, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
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

	results = append(results, tc.runner.RunTest("logs", "GetLogEvents_VerifyMessage", func() error {
		resp, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			return err
		}
		if len(resp.Events) == 0 {
			return fmt.Errorf("no events returned")
		}
		if resp.Events[0].Message == nil || *resp.Events[0].Message != "Test log message" {
			return fmt.Errorf("message mismatch: got %q, want %q", aws.ToString(resp.Events[0].Message), "Test log message")
		}
		if resp.Events[0].Timestamp == nil || *resp.Events[0].Timestamp == 0 {
			return fmt.Errorf("timestamp is zero or nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "FilterLogEvents_VerifyMessage", func() error {
		resp, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if len(resp.Events) == 0 {
			return fmt.Errorf("no events returned")
		}
		if resp.Events[0].Message == nil || *resp.Events[0].Message != "Test log message" {
			return fmt.Errorf("message mismatch: got %q, want %q", aws.ToString(resp.Events[0].Message), "Test log message")
		}
		if resp.Events[0].LogStreamName == nil || *resp.Events[0].LogStreamName != logStreamName {
			return fmt.Errorf("logStreamName mismatch: got %q", aws.ToString(resp.Events[0].LogStreamName))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutRetentionPolicy_VerifyRetention", func() error {
		_, err := tc.client.PutRetentionPolicy(tc.ctx, &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(logGroupName),
			RetentionInDays: aws.Int32(7),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(logGroupName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) == 0 {
			return fmt.Errorf("log group not found")
		}
		if descResp.LogGroups[0].RetentionInDays == nil || *descResp.LogGroups[0].RetentionInDays != 7 {
			return fmt.Errorf("retentionInDays mismatch: got %d, want 7", aws.ToInt32(descResp.LogGroups[0].RetentionInDays))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteLogStream_VerifyGone", func() error {
		_, err := tc.client.DeleteLogStream(tc.ctx, &cloudwatchlogs.DeleteLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.client.DescribeLogStreams(tc.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(logGroupName),
			LogStreamNamePrefix: aws.String(logStreamName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		for _, ls := range descResp.LogStreams {
			if ls.LogStreamName != nil && *ls.LogStreamName == logStreamName {
				return fmt.Errorf("log stream still exists after delete")
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteLogGroup_VerifyGone", func() error {
		_, err := tc.client.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(logGroupName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		for _, lg := range descResp.LogGroups {
			if lg.LogGroupName != nil && *lg.LogGroupName == logGroupName {
				return fmt.Errorf("log group still exists after delete")
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "CreateLogGroup_Duplicate", func() error {
		dupGroupName := tc.uniquePrefix("DupLogGroup")
		if err := tc.createLogGroup(dupGroupName); err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.deleteLogGroup(dupGroupName)

		_, err := tc.client.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dupGroupName),
		})
		if err := AssertErrorContains(err, "ResourceAlreadyExistsException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteLogGroup_NonExistent", func() error {
		_, err := tc.client.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String("nonexistent-log-group-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
