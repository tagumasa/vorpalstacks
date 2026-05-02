package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func (tc *cwlogsTestCtx) eventTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "PutLogEvents_GetLogEvents_Roundtrip", func() error {
		rtGroupName := tc.uniquePrefix("RTLogGroup")
		rtStreamName := tc.uniquePrefix("RTLogStream")
		if err := tc.createLogGroup(rtGroupName); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(rtGroupName)

		if err := tc.createLogStream(rtGroupName, rtStreamName); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		testMessage := "roundtrip-log-message-verify-12345"
		ts := time.Now().UnixMilli()
		if err := tc.putLogEvent(rtGroupName, rtStreamName, testMessage, ts); err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
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

	results = append(results, tc.runner.RunTest("logs", "PutLogEvents_MultipleEvents", func() error {
		meName := tc.uniquePrefix("MEGroup")
		meStream := tc.uniquePrefix("MEStream")
		if err := tc.createLogGroup(meName); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(meName)

		if err := tc.createLogStream(meName, meStream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		_, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
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

		resp, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
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

	results = append(results, tc.runner.RunTest("logs", "GetLogEvents_StartFromHead", func() error {
		sfhName := tc.uniquePrefix("SFHGroup")
		sfhStream := tc.uniquePrefix("SFHStream")
		if err := tc.createLogGroup(sfhName); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(sfhName)

		if err := tc.createLogStream(sfhName, sfhStream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(sfhName),
			LogStreamName: aws.String(sfhStream),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("first-event"), Timestamp: aws.Int64(ts)},
				{Message: aws.String("second-event"), Timestamp: aws.Int64(ts + 1)},
			},
		})

		resp, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
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

	results = append(results, tc.runner.RunTest("logs", "FilterLogEvents_WithFilterPattern", func() error {
		fepName := tc.uniquePrefix("FEPGroup")
		fepStream := tc.uniquePrefix("FEPStream")
		if err := tc.createLogGroup(fepName); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(fepName)

		if err := tc.createLogStream(fepName, fepStream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(fepName),
			LogStreamName: aws.String(fepStream),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("ERROR disk full"), Timestamp: aws.Int64(ts)},
				{Message: aws.String("INFO started"), Timestamp: aws.Int64(ts + 1)},
				{Message: aws.String("ERROR network timeout"), Timestamp: aws.Int64(ts + 2)},
			},
		})

		resp, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
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

	results = append(results, tc.runner.RunTest("logs", "FilterLogEvents_WithLogStreamNames", func() error {
		flsName := tc.uniquePrefix("FLSGroup")
		flsStream1 := tc.uniquePrefix("FLSStream1")
		flsStream2 := tc.uniquePrefix("FLSStream2")
		if err := tc.createLogGroup(flsName); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(flsName)

		tc.createLogStream(flsName, flsStream1)
		tc.createLogStream(flsName, flsStream2)

		ts := time.Now().UnixMilli()
		tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream1),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-stream-1"), Timestamp: aws.Int64(ts)},
			},
		})
		tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream2),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-stream-2"), Timestamp: aws.Int64(ts + 1)},
			},
		})

		resp, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
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

	results = append(results, tc.runner.RunTest("logs", "MetricFilterCount_Tracked", func() error {
		mfcName := tc.uniquePrefix("MFCGroup")
		if err := tc.createLogGroup(mfcName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(mfcName)

		descResp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
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

		tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(mfcName),
			FilterName:    aws.String("CountFilter1"),
			FilterPattern: aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{
				{MetricName: aws.String("E"), MetricNamespace: aws.String("NS"), MetricValue: aws.String("1")},
			},
		})

		descResp2, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(mfcName),
		})
		if err != nil {
			return fmt.Errorf("describe after: %v", err)
		}
		if descResp2.LogGroups[0].MetricFilterCount == nil || *descResp2.LogGroups[0].MetricFilterCount != 1 {
			return fmt.Errorf("expected 1 filter, got %d", aws.ToInt32(descResp2.LogGroups[0].MetricFilterCount))
		}

		tc.client.DeleteMetricFilter(tc.ctx, &cloudwatchlogs.DeleteMetricFilterInput{
			LogGroupName: aws.String(mfcName),
			FilterName:   aws.String("CountFilter1"),
		})

		descResp3, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
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

	results = append(results, tc.runner.RunTest("logs", "DeleteLogStream_NonExistent", func() error {
		dlsName := tc.uniquePrefix("DlsGroup")
		if err := tc.createLogGroup(dlsName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(dlsName)

		_, err := tc.client.DeleteLogStream(tc.ctx, &cloudwatchlogs.DeleteLogStreamInput{
			LogGroupName:  aws.String(dlsName),
			LogStreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
