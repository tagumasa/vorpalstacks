package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwltypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func (r *TestRunner) runCWLogsToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-cwl-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-cwl-role-%s", ts)
	logGroupName := fmt.Sprintf("/integ/cwl-lambda/%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)})
	defer func() {
		ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroupName)})
	}()

	_, err := ic.cwl.PutSubscriptionFilter(ic.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   aws.String(logGroupName),
		FilterName:     aws.String("integ-lambda-sub"),
		FilterPattern:  aws.String("[...]"),
		DestinationArn: aws.String(fnARN),
	})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_Lambda", func() error { return fmt.Errorf("put subscription filter: %w", err) })
	}

	ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
	})

	ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
		LogEvents: []cwltypes.InputLogEvent{
			{Message: aws.String("integration test log message"), Timestamp: aws.Int64(time.Now().UnixMilli())},
		},
	})

	return r.pollVerify("CWLogs_Lambda", defaultPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

func (r *TestRunner) runCWLogsToKinesis(ic *integClients, ts string) TestResult {
	streamName := fmt.Sprintf("integ-cwl-kin-%s", ts)
	logGroupName := fmt.Sprintf("/integ/cwl-kinesis/%s", ts)

	err := ic.createKinesisStream(streamName)
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("create stream: %w", err) })
	}
	defer ic.deleteStream(streamName)

	if err := ic.pollStreamActive(streamName, defaultPollTimeout); err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("stream not active: %w", err) })
	}

	streamARN := fmt.Sprintf("arn:aws:kinesis:us-east-1:000000000000:stream/%s", streamName)

	ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)})
	defer func() {
		ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroupName)})
	}()

	ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
	})

	_, err = ic.cwl.PutSubscriptionFilter(ic.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   aws.String(logGroupName),
		FilterName:     aws.String("integ-kinesis-sub"),
		FilterPattern:  aws.String("[...]"),
		DestinationArn: aws.String(streamARN),
	})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("put subscription filter: %w", err) })
	}

	ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
		LogEvents: []cwltypes.InputLogEvent{
			{Message: aws.String("kinesis subscription test"), Timestamp: aws.Int64(time.Now().UnixMilli())},
		},
	})

	return r.pollVerify("CWLogs_Kinesis", defaultPollTimeout, func() error {
		streamDesc, err := ic.describeStream(streamName)
		if err != nil {
			return fmt.Errorf("describe stream: %w", err)
		}
		if len(streamDesc.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards in stream")
		}
		shardID := streamDesc.StreamDescription.Shards[0].ShardId
		iter, err := ic.createIteratorFromHorizon(streamName, *shardID)
		if err != nil {
			return fmt.Errorf("create iterator: %w", err)
		}
		records, err := ic.getRecords(iter)
		if err != nil {
			return fmt.Errorf("get records: %w", err)
		}
		if len(records.Records) == 0 {
			return fmt.Errorf("expected records from CW Logs subscription, got 0")
		}
		if !strings.Contains(string(records.Records[0].Data), "awslogs") {
			return fmt.Errorf("expected awslogs envelope, got: %s", string(records.Records[0].Data))
		}
		return nil
	})
}
