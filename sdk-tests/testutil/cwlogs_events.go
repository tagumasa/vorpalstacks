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

	// FilteredLogEvent carries an eventId — a stable identifier, so the
	// same stored event yields the same id on repeated reads.
	results = append(results, tc.runner.RunTest("logs", "FilterLogEvents_EventId", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("EventIdGroup", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", `{"level": "INFO"}`, now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		read := func() (string, error) {
			var nextToken *string
			for i := 0; i < 10; i++ {
				resp, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
					LogGroupName: aws.String(groupName),
					NextToken:    nextToken,
				})
				if err != nil {
					return "", fmt.Errorf("filter log events: %v", err)
				}
				for _, ev := range resp.Events {
					if aws.ToString(ev.Message) == `{"level": "INFO"}` {
						return aws.ToString(ev.EventId), nil
					}
				}
				if resp.NextToken == nil {
					return "", fmt.Errorf("event not found in filter results")
				}
				nextToken = resp.NextToken
			}
			return "", fmt.Errorf("pagination did not converge")
		}
		first, err := read()
		if err != nil {
			return err
		}
		if first == "" {
			return fmt.Errorf("FilteredLogEvent carried an empty eventId")
		}
		second, err := read()
		if err != nil {
			return err
		}
		if first != second {
			return fmt.Errorf("eventId drifted across reads: %q then %q", first, second)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutLogEvents_GetLogEvents_Roundtrip", func() error {
		rtGroupName, cleanupGroup, err := tc.newLogGroupFixture("RTLogGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		rtStreamName := tc.uniquePrefix("RTLogStream")
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
		if resp.Events[0].Timestamp == nil || *resp.Events[0].Timestamp == 0 {
			return fmt.Errorf("timestamp is zero or nil")
		}
		if resp.Events[0].IngestionTime == nil || *resp.Events[0].IngestionTime == 0 {
			return fmt.Errorf("ingestionTime is zero or nil")
		}

		// The identifier member addresses the same group by ARN — the
		// object form DescribeLogGroups reports, trailing ":*" included —
		// and must read the same events.
		groupArn, err := tc.findLogGroupARN(rtGroupName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		arnResp, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupIdentifier: groupArn,
			LogStreamName:      aws.String(rtStreamName),
		})
		if err != nil {
			return fmt.Errorf("get by identifier ARN: %v", err)
		}
		if len(arnResp.Events) == 0 || aws.ToString(arnResp.Events[0].Message) != testMessage {
			return fmt.Errorf("identifier-ARN read returned %d events, first %q", len(arnResp.Events), aws.ToString(arnResp.Events[0].Message))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutLogEvents_MultipleEvents", func() error {
		meName, cleanupGroup, err := tc.newLogGroupFixture("MEGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		meStream := tc.uniquePrefix("MEStream")
		if err := tc.createLogStream(meName, meStream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		_, err = tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
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
		sfhName, cleanupGroup, err := tc.newLogGroupFixture("SFHGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		sfhStream := tc.uniquePrefix("SFHStream")
		if err := tc.createLogStream(sfhName, sfhStream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(sfhName),
			LogStreamName: aws.String(sfhStream),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("first-event"), Timestamp: aws.Int64(ts)},
				{Message: aws.String("second-event"), Timestamp: aws.Int64(ts + 1)},
			},
		}); err != nil {
			return fmt.Errorf("put: %v", err)
		}

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

	results = append(results, tc.runner.RunTest("logs", "GetLogEvents_ForwardPagination", func() error {
		fpName, cleanupGroup, err := tc.newLogGroupFixture("FPLogGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		fpStream := tc.uniquePrefix("FPLogStream")
		if err := tc.createLogStream(fpName, fpStream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		putEvents := []types.InputLogEvent{
			{Message: aws.String("msg-0"), Timestamp: aws.Int64(ts)},
			{Message: aws.String("msg-1"), Timestamp: aws.Int64(ts + 1)},
			{Message: aws.String("msg-2"), Timestamp: aws.Int64(ts + 2)},
			{Message: aws.String("msg-3"), Timestamp: aws.Int64(ts + 3)},
			{Message: aws.String("msg-4"), Timestamp: aws.Int64(ts + 4)},
		}
		if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(fpName),
			LogStreamName: aws.String(fpStream),
			LogEvents:     putEvents,
		}); err != nil {
			return fmt.Errorf("put: %v", err)
		}

		var collected []string
		var nextToken *string
		page := 0
		for {
			page++
			if page > 10 {
				return fmt.Errorf("too many pages")
			}
			resp, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(fpName),
				LogStreamName: aws.String(fpStream),
				StartFromHead: aws.Bool(true),
				Limit:         aws.Int32(2),
				NextToken:     nextToken,
			})
			if err != nil {
				return fmt.Errorf("get page %d: %v", page, err)
			}
			for _, e := range resp.Events {
				collected = append(collected, *e.Message)
			}
			// Termination per the documented contract: the tokens are
			// never null, and at the end of the stream the operation
			// returns the same token that was passed in — an empty page
			// with an unchanged token ends the walk.
			if resp.NextForwardToken == nil || *resp.NextForwardToken == "" {
				return fmt.Errorf("page %d: nextForwardToken must never be null", page)
			}
			if nextToken != nil && *resp.NextForwardToken == *nextToken {
				break
			}
			if len(resp.Events) == 0 {
				return fmt.Errorf("page %d: empty page before token repetition", page)
			}
			nextToken = resp.NextForwardToken
		}

		if len(collected) != 5 {
			return fmt.Errorf("expected 5 events across pages, got %d: %v", len(collected), collected)
		}
		for i, msg := range collected {
			expected := fmt.Sprintf("msg-%d", i)
			if msg != expected {
				return fmt.Errorf("event %d: expected %s, got %s", i, expected, msg)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "FilterLogEvents_WithFilterPattern", func() error {
		fepName, cleanupGroup, err := tc.newLogGroupFixture("FEPGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		fepStream := tc.uniquePrefix("FEPStream")
		if err := tc.createLogStream(fepName, fepStream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		ts := time.Now().UnixMilli()
		if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(fepName),
			LogStreamName: aws.String(fepStream),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("ERROR disk full"), Timestamp: aws.Int64(ts)},
				{Message: aws.String("INFO started"), Timestamp: aws.Int64(ts + 1)},
				{Message: aws.String("ERROR network timeout"), Timestamp: aws.Int64(ts + 2)},
			},
		}); err != nil {
			return fmt.Errorf("put: %v", err)
		}

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

		// The identifier member addresses the same read by ARN.
		fepArn, err := tc.findLogGroupARN(fepName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		arnResp, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
			LogGroupIdentifier: fepArn,
			FilterPattern:      aws.String("ERROR"),
		})
		if err != nil {
			return fmt.Errorf("filter by identifier ARN: %v", err)
		}
		if len(arnResp.Events) != 2 {
			return fmt.Errorf("identifier-ARN filter expected 2 ERROR events, got %d", len(arnResp.Events))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "FilterLogEvents_WithLogStreamNames", func() error {
		flsName, cleanupGroup, err := tc.newLogGroupFixture("FLSGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		flsStream1 := tc.uniquePrefix("FLSStream1")
		flsStream2 := tc.uniquePrefix("FLSStream2")
		if err := tc.createLogStream(flsName, flsStream1); err != nil {
			return fmt.Errorf("create stream 1: %v", err)
		}
		if err := tc.createLogStream(flsName, flsStream2); err != nil {
			return fmt.Errorf("create stream 2: %v", err)
		}

		ts := time.Now().UnixMilli()
		if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream1),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-stream-1"), Timestamp: aws.Int64(ts)},
			},
		}); err != nil {
			return fmt.Errorf("put stream 1: %v", err)
		}
		if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String(flsStream2),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-stream-2"), Timestamp: aws.Int64(ts + 1)},
			},
		}); err != nil {
			return fmt.Errorf("put stream 2: %v", err)
		}

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

		// Prefix mode selects the streams by name prefix — including a
		// prefix that carries '/', which the streams' own names carry.
		for _, svcStream := range []string{"svc/alpha", "svc/beta"} {
			if err := tc.createLogStream(flsName, svcStream); err != nil {
				return fmt.Errorf("create %s: %v", svcStream, err)
			}
		}
		tsSvc := time.Now().UnixMilli()
		if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String("svc/alpha"),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-svc-alpha"), Timestamp: aws.Int64(tsSvc)},
			},
		}); err != nil {
			return fmt.Errorf("put svc/alpha: %v", err)
		}
		if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(flsName),
			LogStreamName: aws.String("svc/beta"),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String("from-svc-beta"), Timestamp: aws.Int64(tsSvc + 1)},
			},
		}); err != nil {
			return fmt.Errorf("put svc/beta: %v", err)
		}
		prefResp, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:        aws.String(flsName),
			LogStreamNamePrefix: aws.String("svc/"),
		})
		if err != nil {
			return fmt.Errorf("filter by stream prefix: %v", err)
		}
		if len(prefResp.Events) != 2 {
			return fmt.Errorf("expected 2 events from 'svc/' streams, got %d", len(prefResp.Events))
		}
		for _, e := range prefResp.Events {
			if e.LogStreamName == nil || !strings.HasPrefix(*e.LogStreamName, "svc/") {
				return fmt.Errorf("event from outside the 'svc/' prefix: %q", aws.ToString(e.LogStreamName))
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "MetricFilterCount_Tracked", func() error {
		mfcName, cleanupGroup, err := tc.newLogGroupFixture("MFCGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

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

		if _, err := tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(mfcName),
			FilterName:    aws.String("CountFilter1"),
			FilterPattern: aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{
				{MetricName: aws.String("E"), MetricNamespace: aws.String("NS"), MetricValue: aws.String("1")},
			},
		}); err != nil {
			return fmt.Errorf("put metric filter: %v", err)
		}

		descResp2, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(mfcName),
		})
		if err != nil {
			return fmt.Errorf("describe after: %v", err)
		}
		// A page that comes back without the group is a failed assertion,
		// not a harness panic: these re-reads carry the same emptiness
		// check the fixture's first read has.
		if len(descResp2.LogGroups) == 0 {
			return fmt.Errorf("group %q missing from the listing after the filter was put", mfcName)
		}
		if descResp2.LogGroups[0].MetricFilterCount == nil || *descResp2.LogGroups[0].MetricFilterCount != 1 {
			return fmt.Errorf("expected 1 filter, got %d", aws.ToInt32(descResp2.LogGroups[0].MetricFilterCount))
		}

		if _, err := tc.client.DeleteMetricFilter(tc.ctx, &cloudwatchlogs.DeleteMetricFilterInput{
			LogGroupName: aws.String(mfcName),
			FilterName:   aws.String("CountFilter1"),
		}); err != nil {
			return fmt.Errorf("delete metric filter: %v", err)
		}

		descResp3, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(mfcName),
		})
		if err != nil {
			return fmt.Errorf("describe after delete: %v", err)
		}
		if len(descResp3.LogGroups) == 0 {
			return fmt.Errorf("group %q missing from the listing after the filter was deleted", mfcName)
		}
		if descResp3.LogGroups[0].MetricFilterCount == nil || *descResp3.LogGroups[0].MetricFilterCount != 0 {
			return fmt.Errorf("expected 0 filters after delete, got %d", aws.ToInt32(descResp3.LogGroups[0].MetricFilterCount))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteLogStream_NonExistent", func() error {
		dlsName, cleanupGroup, err := tc.newLogGroupFixture("DlsGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		_, err = tc.client.DeleteLogStream(tc.ctx, &cloudwatchlogs.DeleteLogStreamInput{
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
