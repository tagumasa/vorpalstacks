package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	kinesistypes "github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"vorpalstacks-sdk-tests/config"
)

const (
	integSvc = "integration"

	lambdaTrustPolicy    = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	sfnTrustPolicy       = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"states.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	schedulerTrustPolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"scheduler.amazonaws.com"},"Action":"sts:AssumeRole"}]}`

	echoHandlerCode = `exports.handler = async (event) => { return JSON.stringify(event); };`

	pollInterval         = 300 * time.Millisecond
	integTestTimeout     = 3 * time.Minute
	defaultPollTimeout   = 10 * time.Second
	schedulerPollTimeout = 15 * time.Second
)

func intTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func intRoleARN(roleName, accountID string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
}

type integClients struct {
	lambda    *lambda.Client
	eb        *eventbridge.Client
	cw        *cloudwatch.Client
	cwl       *cloudwatchlogs.Client
	sfn       *sfn.Client
	scheduler *scheduler.Client
	sns       *sns.Client
	sqs       *sqs.Client
	kinesis   *kinesis.Client
	s3        *s3.Client
	iam       *iam.Client
	dynamodb  *dynamodb.Client
	kms       *kms.Client
	ctx       context.Context
	region    string
	accountID string
}

func (r *TestRunner) newIntegClients() (*integClients, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, err
	}
	ic := &integClients{
		lambda:    lambda.NewFromConfig(cfg),
		eb:        eventbridge.NewFromConfig(cfg),
		cw:        cloudwatch.NewFromConfig(cfg),
		cwl:       cloudwatchlogs.NewFromConfig(cfg),
		sfn:       sfn.NewFromConfig(cfg),
		scheduler: scheduler.NewFromConfig(cfg),
		sns:       sns.NewFromConfig(cfg),
		sqs:       sqs.NewFromConfig(cfg),
		kinesis:   kinesis.NewFromConfig(cfg),
		s3:        s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true }),
		iam:       iam.NewFromConfig(cfg),
		dynamodb:  dynamodb.NewFromConfig(cfg),
		kms:       kms.NewFromConfig(cfg),
		ctx:       context.Background(),
		region:    r.region,
		accountID: r.accountID,
	}
	return ic, nil
}

func (ic *integClients) createLambda(name, roleName string) (string, error) {
	zipCode, err := zipLambdaCode(echoHandlerCode)
	if err != nil {
		return "", fmt.Errorf("zip lambda code: %w", err)
	}
	_, err = ic.lambda.CreateFunction(ic.ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String(name),
		Runtime:      lambdatypes.RuntimeNodejs22x,
		Role:         aws.String(intRoleARN(roleName, ic.accountID)),
		Handler:      aws.String("index.handler"),
		Code:         &lambdatypes.FunctionCode{ZipFile: zipCode},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", ic.region, ic.accountID, name), nil
}

func (ic *integClients) deleteLambda(name string) {
	ic.lambda.DeleteFunction(ic.ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(name)})
}

func (ic *integClients) createQueue(name string) (string, error) {
	resp, err := ic.sqs.CreateQueue(ic.ctx, &sqs.CreateQueueInput{QueueName: aws.String(name)})
	if err != nil {
		return "", err
	}
	return *resp.QueueUrl, nil
}

func (ic *integClients) deleteQueue(queueURL string) {
	ic.sqs.DeleteQueue(ic.ctx, &sqs.DeleteQueueInput{QueueUrl: aws.String(queueURL)})
}

func (ic *integClients) receiveMessages(queueURL string, maxMessages int32, waitSeconds int32) ([]sqstypes.Message, error) {
	resp, err := ic.sqs.ReceiveMessage(ic.ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitSeconds,
	})
	if err != nil {
		return nil, err
	}
	return resp.Messages, nil
}

func (ic *integClients) createTopic(name string) (string, error) {
	resp, err := ic.sns.CreateTopic(ic.ctx, &sns.CreateTopicInput{Name: aws.String(name)})
	if err != nil {
		return "", err
	}
	return *resp.TopicArn, nil
}

func (ic *integClients) deleteTopic(topicARN string) {
	ic.sns.DeleteTopic(ic.ctx, &sns.DeleteTopicInput{TopicArn: aws.String(topicARN)})
}

func (ic *integClients) createKinesisStream(name string) error {
	_, err := ic.kinesis.CreateStream(ic.ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(name),
		ShardCount: aws.Int32(1),
	})
	return err
}

func (ic *integClients) deleteStream(name string) {
	ic.kinesis.DeleteStream(ic.ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(name)})
}

func (ic *integClients) describeStream(name string) (*kinesis.DescribeStreamOutput, error) {
	return ic.kinesis.DescribeStream(ic.ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(name)})
}

func (ic *integClients) createIteratorFromHorizon(streamName, shardID string) (string, error) {
	return ic.createIteratorWithType(streamName, shardID, kinesistypes.ShardIteratorTypeTrimHorizon)
}

func (ic *integClients) createIteratorWithType(streamName, shardID string, iterType kinesistypes.ShardIteratorType) (string, error) {
	resp, err := ic.kinesis.GetShardIterator(ic.ctx, &kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamName),
		ShardId:           aws.String(shardID),
		ShardIteratorType: iterType,
	})
	if err != nil {
		return "", err
	}
	return *resp.ShardIterator, nil
}

func (ic *integClients) getRecords(iterator string) (*kinesis.GetRecordsOutput, error) {
	return ic.kinesis.GetRecords(ic.ctx, &kinesis.GetRecordsInput{ShardIterator: aws.String(iterator)})
}

func (ic *integClients) verifyLambdaInvoked(fnName string) error {
	logGroupName := fmt.Sprintf("/aws/lambda/%s", fnName)
	resp, err := ic.cwl.DescribeLogStreams(ic.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		return fmt.Errorf("lambda %s was not invoked (no CW Logs at %s): %w", fnName, logGroupName, err)
	}
	if len(resp.LogStreams) == 0 {
		return fmt.Errorf("lambda %s was not invoked (no log streams in %s)", fnName, logGroupName)
	}
	return nil
}

func (ic *integClients) verifyLambdaLogContains(fnName, substr string) error {
	logGroupName := fmt.Sprintf("/aws/lambda/%s", fnName)
	streams, err := ic.cwl.DescribeLogStreams(ic.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		return fmt.Errorf("lambda %s: describe log streams: %w", fnName, err)
	}
	if len(streams.LogStreams) == 0 {
		return fmt.Errorf("lambda %s: no log streams", fnName)
	}
	logs, err := ic.cwl.GetLogEvents(ic.ctx, &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: streams.LogStreams[0].LogStreamName,
		Limit:         aws.Int32(50),
	})
	if err != nil {
		return fmt.Errorf("lambda %s: get log events: %w", fnName, err)
	}
	for _, ev := range logs.Events {
		if strings.Contains(aws.ToString(ev.Message), substr) {
			return nil
		}
	}
	return fmt.Errorf("lambda %s: log does not contain %q", fnName, substr)
}

func (ic *integClients) verifyMessageContains(queueURL, substr string) error {
	msgs, err := ic.receiveMessages(queueURL, 5, 3)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return fmt.Errorf("expected message in queue, got 0")
	}
	for _, m := range msgs {
		if strings.Contains(aws.ToString(m.Body), substr) {
			return nil
		}
	}
	return fmt.Errorf("no message contains %q", substr)
}

func (ic *integClients) createBucket(name string) error {
	_, err := ic.s3.CreateBucket(ic.ctx, &s3.CreateBucketInput{Bucket: aws.String(name)})
	return err
}

func (ic *integClients) deleteBucket(name string) {
	var contToken *string
	for {
		listResp, err := ic.s3.ListObjectsV2(ic.ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(name),
			ContinuationToken: contToken,
		})
		if err != nil {
			break
		}
		if len(listResp.Contents) > 0 {
			var objs []s3types.ObjectIdentifier
			for _, o := range listResp.Contents {
				objs = append(objs, s3types.ObjectIdentifier{Key: o.Key})
			}
			ic.s3.DeleteObjects(ic.ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(name),
				Delete: &s3types.Delete{Objects: objs},
			})
		}
		if !aws.ToBool(listResp.IsTruncated) {
			break
		}
		contToken = listResp.NextContinuationToken
	}
	ic.s3.DeleteBucket(ic.ctx, &s3.DeleteBucketInput{Bucket: aws.String(name)})
}

func (ic *integClients) putObject(bucket, key string, data []byte) error {
	_, err := ic.s3.PutObject(ic.ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(string(data)),
	})
	return err
}

// runIntegWithTimeout wraps an integration test function with a deadline.
// Integration test functions make many HTTP calls and, unlike SDK tests, are
// not individually wrapped in RunTest. Without this guard, a single hanging
// HTTP call (e.g. server unresponsive) blocks the entire integration phase
// indefinitely because context.Background has no deadline and the leaked
// goroutine keeps running.
func (r *TestRunner) runIntegWithTimeout(name string, fn func() TestResult) TestResult {
	if r.verbose {
		fmt.Printf("  Running: %s...\n", name)
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), integTestTimeout)
	defer cancel()

	done := make(chan TestResult, 1)
	go func() {
		result := fn()
		result.Duration = time.Since(start)
		if result.TestName == "" {
			result.TestName = name
		}
		done <- result
	}()

	select {
	case res := <-done:
		if r.verbose {
			if res.Status == "PASS" {
				fmt.Printf("  ✓ %s (%.2fs)\n", name, res.Duration.Seconds())
			} else {
				fmt.Printf("  ✗ %s: %s\n", name, res.Error)
			}
		}
		return res
	case <-ctx.Done():
		if r.verbose {
			fmt.Printf("  ✗ %s: timed out after %v\n", name, integTestTimeout)
		}
		return TestResult{
			Service:  integSvc,
			TestName: name,
			Status:   "FAIL",
			Error:    fmt.Sprintf("integration test timed out after %v", integTestTimeout),
			Duration: integTestTimeout,
		}
	}
}

func (r *TestRunner) pollVerify(testName string, timeout time.Duration, verify func() error) TestResult {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := verify(); err == nil {
			return TestResult{
				Service:  integSvc,
				TestName: testName,
				Status:   "PASS",
				Duration: timeout - time.Until(deadline),
			}
		} else {
			lastErr = err
		}
		time.Sleep(pollInterval)
	}
	errMsg := "timeout waiting for condition"
	if lastErr != nil {
		errMsg = lastErr.Error()
	}
	return TestResult{
		Service:  integSvc,
		TestName: testName,
		Status:   "FAIL",
		Error:    errMsg,
		Duration: timeout,
	}
}

func (ic *integClients) pollStreamActive(streamName string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		desc, err := ic.describeStream(streamName)
		if err != nil {
			return err
		}
		if desc.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
			return nil
		}
		time.Sleep(pollInterval)
	}
	return fmt.Errorf("stream %s not active after %v", streamName, timeout)
}

func (r *TestRunner) RunIntegrationTests() []TestResult {
	var results []TestResult

	ic, err := r.newIntegClients()
	if err != nil {
		return append(results, SetupFailResult(integSvc, "Failed to create clients: %v", err))
	}

	ts := intTimestamp()

	results = append(results, r.runIntegWithTimeout("EventBridge_Lambda", func() TestResult {
		return r.runEventBridgeToLambda(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("EventBridge_StepFunctions", func() TestResult {
		return r.runEventBridgeToStepFunctions(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("EventBridge_SQS", func() TestResult {
		return r.runEventBridgeToSQS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("EventBridge_SNS", func() TestResult {
		return r.runEventBridgeToSNS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("EventBridge_Kinesis", func() TestResult {
		return r.runEventBridgeToKinesis(ic, ts)
	}))

	results = append(results, r.runIntegWithTimeout("ESM_SQS_Lambda", func() TestResult {
		return r.runESMSQSToLambda(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("ESM_Kinesis_Lambda", func() TestResult {
		return r.runESMKinesisToLambda(ic, ts)
	}))

	results = append(results, r.runIntegWithTimeout("Alarm_SNS", func() TestResult {
		return r.runAlarmToSNS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("Alarm_Lambda", func() TestResult {
		return r.runAlarmToLambda(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("Alarm_StepFunctions", func() TestResult {
		return r.runAlarmToStepFunctions(ic, ts)
	}))

	results = append(results, r.runIntegWithTimeout("Scheduler_Lambda", func() TestResult {
		return r.runSchedulerToLambda(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("Scheduler_SQS", func() TestResult {
		return r.runSchedulerToSQS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("Scheduler_SNS", func() TestResult {
		return r.runSchedulerToSNS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("Scheduler_StepFunctions", func() TestResult {
		return r.runSchedulerToStepFunctions(ic, ts)
	}))

	results = append(results, r.runIntegWithTimeout("SFNTask_Lambda", func() TestResult {
		return r.runSFNTaskLambda(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("SFNTask_SQS", func() TestResult {
		return r.runSFNTaskSQS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("SFNTask_SNS", func() TestResult {
		return r.runSFNTaskSNS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("SFNTask_EventBridge", func() TestResult {
		return r.runSFNTaskEventBridge(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("SFNTask_DynamoDB", func() TestResult {
		return r.runSFNTaskDynamoDB(ic, ts)
	}))

	results = append(results, r.runIntegWithTimeout("S3_Notification_Lambda", func() TestResult {
		return r.runS3NotificationToLambda(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("S3_Notification_SQS", func() TestResult {
		return r.runS3NotificationToSQS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("S3_ObjectTagging_Notification_SQS", func() TestResult {
		return r.runS3TaggingNotificationToSQS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("S3_Notification_SNS", func() TestResult {
		return r.runS3NotificationToSNS(ic, ts)
	}))

	results = append(results, r.runIntegWithTimeout("CWLogs_Lambda", func() TestResult {
		return r.runCWLogsToLambda(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("CWLogs_Kinesis", func() TestResult {
		return r.runCWLogsToKinesis(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("CWLogs_LookupTable_KMS", func() TestResult {
		return r.runCWLogsLookupTableKMS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("CWLogs_ScheduledQuery_S3", func() TestResult {
		return r.runCWLogsScheduledQueryS3(ic, ts)
	}))

	results = append(results, r.runIntegWithTimeout("SNS_SQS", func() TestResult {
		return r.runSNSToSQS(ic, ts)
	}))
	results = append(results, r.runIntegWithTimeout("SNS_Lambda", func() TestResult {
		return r.runSNSToLambda(ic, ts)
	}))

	// Neptune direct protocol tests already wrap each sub-test in RunTest
	// (60s per-test timeout). The cluster creation is the only unguarded
	// part; under normal conditions it completes in <1s.
	results = append(results, r.RunNeptuneDirectProtocolTests()...)

	return results
}
