package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwltypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	kinesistypes "github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulertypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
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

	pollInterval = 300 * time.Millisecond
	defaultPollTimeout = 10 * time.Second
	schedulerPollTimeout = 15 * time.Second
)

func intTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func intRoleARN(roleName string) string {
	return fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
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
	ctx       context.Context
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
		ctx:       context.Background(),
	}
	return ic, nil
}

func (ic *integClients) createLambda(name, roleName string) (string, error) {
	_, err := ic.lambda.CreateFunction(ic.ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String(name),
		Runtime:      lambdatypes.RuntimeNodejs22x,
		Role:         aws.String(intRoleARN(roleName)),
		Handler:      aws.String("index.handler"),
		Code:         &lambdatypes.FunctionCode{ZipFile: []byte(echoHandlerCode)},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", "us-east-1", name), nil
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

func (ic *integClients) createBucket(name string) error {
	_, err := ic.s3.CreateBucket(ic.ctx, &s3.CreateBucketInput{Bucket: aws.String(name)})
	return err
}

func (ic *integClients) deleteBucket(name string) {
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

	results = append(results, r.runEventBridgeToLambda(ic, ts))
	results = append(results, r.runEventBridgeToStepFunctions(ic, ts))
	results = append(results, r.runEventBridgeToSQS(ic, ts))
	results = append(results, r.runEventBridgeToSNS(ic, ts))
	results = append(results, r.runEventBridgeToKinesis(ic, ts))

	results = append(results, r.runESMSQSToLambda(ic, ts))
	results = append(results, r.runESMKinesisToLambda(ic, ts))
	results = append(results, r.SkipTest(integSvc, "ESM_DynamoDBStreams_ToLambda", "DynamoDB Streams ESM not yet implemented on server"))

	results = append(results, r.runAlarmToSNS(ic, ts))
	results = append(results, r.runAlarmToLambda(ic, ts))
	results = append(results, r.runAlarmToStepFunctions(ic, ts))

	results = append(results, r.runSchedulerToLambda(ic, ts))
	results = append(results, r.runSchedulerToSQS(ic, ts))
	results = append(results, r.runSchedulerToSNS(ic, ts))
	results = append(results, r.runSchedulerToStepFunctions(ic, ts))

	results = append(results, r.runSFNTaskLambda(ic, ts))
	results = append(results, r.runSFNTaskSQS(ic, ts))
	results = append(results, r.runSFNTaskSNS(ic, ts))
	results = append(results, r.SkipTest(integSvc, "SFN_Task_DynamoDB", "DynamoDB task integration not yet implemented on server"))

	results = append(results, r.runS3NotificationToLambda(ic, ts))
	results = append(results, r.runS3NotificationToSQS(ic, ts))
	results = append(results, r.runS3NotificationToSNS(ic, ts))
	results = append(results, r.SkipTest(integSvc, "S3_Notification_Kinesis", "S3 notification to Kinesis not yet implemented on server"))

	results = append(results, r.runCWLogsToLambda(ic, ts))
	results = append(results, r.runCWLogsToKinesis(ic, ts))

	results = append(results, r.runSNSToSQS(ic, ts))
	results = append(results, r.runSNSToLambda(ic, ts))

	return results
}

// ---------- EventBridge -> Lambda ----------

func (r *TestRunner) runEventBridgeToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-eb-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-eb-lambda-role-%s", ts)
	busName := fmt.Sprintf("integ-eb-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-rule-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)
	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(fnARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	detail := map[string]string{"test": "eventbridge-lambda"}
	detailJSON, _ := json.Marshal(detail)
	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("IntegrationTest"),
			Detail:       aws.String(string(detailJSON)),
		}},
	})

	return r.pollVerify("EventBridge_Lambda", defaultPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

// ---------- EventBridge -> Step Functions ----------

func (r *TestRunner) runEventBridgeToStepFunctions(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-eb-sfn-role-%s", ts)
	busName := fmt.Sprintf("integ-eb-sfn-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-sfn-rule-%s", ts)
	smName := fmt.Sprintf("integ-eb-sfn-sm-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	_, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(`{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_StepFunctions", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName))})
	}()

	smARN := fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)
	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(smARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("SFNTrigger"),
			Detail:       aws.String(`{"test":"eb-to-sfn"}`),
		}},
	})

	return r.pollVerify("EventBridge_StepFunctions", defaultPollTimeout, func() error {
		resp, err := ic.sfn.ListExecutions(ic.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Executions) == 0 {
			return fmt.Errorf("expected at least 1 execution, got 0")
		}
		return nil
	})
}

// ---------- EventBridge -> SQS ----------

func (r *TestRunner) runEventBridgeToSQS(ic *integClients, ts string) TestResult {
	queueName := fmt.Sprintf("integ-eb-sqs-%s", ts)
	busName := fmt.Sprintf("integ-eb-sqs-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-sqs-rule-%s", ts)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(queueARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("SQSTest"),
			Detail:       aws.String(`{"message":"eb-to-sqs"}`),
		}},
	})

	return r.pollVerify("EventBridge_SQS", defaultPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected message in queue, got 0")
		}
		return nil
	})
}

// ---------- EventBridge -> SNS ----------

func (r *TestRunner) runEventBridgeToSNS(ic *integClients, ts string) TestResult {
	topicName := fmt.Sprintf("integ-eb-sns-%s", ts)
	busName := fmt.Sprintf("integ-eb-sns-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-sns-rule-%s", ts)
	queueName := fmt.Sprintf("integ-eb-sns-sqs-%s", ts)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(topicARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("SNSTest"),
			Detail:       aws.String(`{"message":"eb-to-sns"}`),
		}},
	})

	return r.pollVerify("EventBridge_SNS", defaultPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected message in queue (via SNS), got 0")
		}
		return nil
	})
}

// ---------- EventBridge -> Kinesis ----------

func (r *TestRunner) runEventBridgeToKinesis(ic *integClients, ts string) TestResult {
	streamName := fmt.Sprintf("integ-eb-kinesis-%s", ts)
	busName := fmt.Sprintf("integ-eb-kin-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-kin-rule-%s", ts)

	err := ic.createKinesisStream(streamName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_Kinesis", func() error { return fmt.Errorf("create stream: %w", err) })
	}
	defer ic.deleteStream(streamName)

	if err := ic.pollStreamActive(streamName, defaultPollTimeout); err != nil {
		return r.RunTest(integSvc, "EventBridge_Kinesis", func() error { return fmt.Errorf("stream not active: %w", err) })
	}

	streamARN := fmt.Sprintf("arn:aws:kinesis:us-east-1:000000000000:stream/%s", streamName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(streamARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("KinesisTest"),
			Detail:       aws.String(`{"message":"eb-to-kinesis"}`),
		}},
	})

	return r.pollVerify("EventBridge_Kinesis", defaultPollTimeout, func() error {
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
			return fmt.Errorf("expected records in kinesis stream, got 0")
		}
		return nil
	})
}

// ---------- Lambda ESM -> SQS ----------

func (r *TestRunner) runESMSQSToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-esm-sqs-fn-%s", ts)
	roleName := fmt.Sprintf("integ-esm-sqs-role-%s", ts)
	queueName := fmt.Sprintf("integ-esm-sqs-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "ESM_SQS_Lambda", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	_, err = ic.lambda.CreateEventSourceMapping(ic.ctx, &lambda.CreateEventSourceMappingInput{
		FunctionName:   aws.String(fnName),
		EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
		Enabled:        aws.Bool(true),
		BatchSize:      aws.Int32(10),
	})
	if err != nil {
		return r.RunTest(integSvc, "ESM_SQS_Lambda", func() error { return fmt.Errorf("create ESM: %w", err) })
	}

	ic.sqs.SendMessage(ic.ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(`{"test":"esm-sqs"}`),
	})

	return r.pollVerify("ESM_SQS_Lambda", defaultPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 10, 1)
		if err != nil {
			return err
		}
		for _, m := range msgs {
			ic.sqs.DeleteMessage(ic.ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: m.ReceiptHandle,
			})
		}
		if len(msgs) > 0 {
			return fmt.Errorf("expected ESM to consume all messages, got %d remaining", len(msgs))
		}
		return nil
	})
}

// ---------- Lambda ESM -> Kinesis ----------

func (r *TestRunner) runESMKinesisToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-esm-kin-fn-%s", ts)
	roleName := fmt.Sprintf("integ-esm-kin-role-%s", ts)
	streamName := fmt.Sprintf("integ-esm-kin-s-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	err := ic.createKinesisStream(streamName)
	if err != nil {
		return r.RunTest(integSvc, "ESM_Kinesis_Lambda", func() error { return fmt.Errorf("create stream: %w", err) })
	}
	defer ic.deleteStream(streamName)

	if err := ic.pollStreamActive(streamName, defaultPollTimeout); err != nil {
		return r.RunTest(integSvc, "ESM_Kinesis_Lambda", func() error { return fmt.Errorf("stream not active: %w", err) })
	}

	esmResp, err := ic.lambda.CreateEventSourceMapping(ic.ctx, &lambda.CreateEventSourceMappingInput{
		FunctionName:     aws.String(fnName),
		EventSourceArn:   aws.String(fmt.Sprintf("arn:aws:kinesis:us-east-1:000000000000:stream/%s", streamName)),
		Enabled:          aws.Bool(true),
		BatchSize:        aws.Int32(100),
		StartingPosition: lambdatypes.EventSourcePositionLatest,
	})
	if err != nil {
		return r.RunTest(integSvc, "ESM_Kinesis_Lambda", func() error { return fmt.Errorf("create ESM: %w", err) })
	}
	esmUUID := esmResp.UUID
	defer func() {
		ic.lambda.DeleteEventSourceMapping(ic.ctx, &lambda.DeleteEventSourceMappingInput{UUID: esmUUID})
	}()

	ic.kinesis.PutRecord(ic.ctx, &kinesis.PutRecordInput{
		StreamName:   aws.String(streamName),
		PartitionKey: aws.String("p1"),
		Data:         []byte(`{"test":"esm-kinesis"}`),
	})

	return r.pollVerify("ESM_Kinesis_Lambda", defaultPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

// ---------- CloudWatch Alarm -> SNS ----------

func (r *TestRunner) runAlarmToSNS(ic *integClients, ts string) TestResult {
	topicName := fmt.Sprintf("integ-alarm-sns-%s", ts)
	queueName := fmt.Sprintf("integ-alarm-sns-q-%s", ts)
	alarmName := fmt.Sprintf("integ-alarm-sns-%s", ts)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	_, err = ic.cw.PutMetricAlarm(ic.ctx, &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		MetricName:         aws.String("CPUUtilization"),
		Namespace:          aws.String("AWS/EC2"),
		Statistic:          cloudwatchtypes.StatisticAverage,
		Period:             aws.Int32(1),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudwatchtypes.ComparisonOperatorGreaterThanThreshold,
		AlarmActions:       []string{topicARN},
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_SNS", func() error { return fmt.Errorf("put alarm: %w", err) })
	}
	defer func() {
		ic.cw.DeleteAlarms(ic.ctx, &cloudwatch.DeleteAlarmsInput{AlarmNames: []string{alarmName}})
	}()

	ic.cw.PutMetricData(ic.ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("AWS/EC2"),
		MetricData: []cloudwatchtypes.MetricDatum{{
			MetricName: aws.String("CPUUtilization"),
			Value:      aws.Float64(100),
		}},
	})

	return r.pollVerify("CWAlarm_SNS", defaultPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected alarm notification in queue, got 0")
		}
		return nil
	})
}

// ---------- CloudWatch Alarm -> Lambda ----------

func (r *TestRunner) runAlarmToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-alarm-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-alarm-lambda-role-%s", ts)
	alarmName := fmt.Sprintf("integ-alarm-lambda-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	_, err := ic.cw.PutMetricAlarm(ic.ctx, &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		MetricName:         aws.String("MemoryUtilization"),
		Namespace:          aws.String("AWS/EC2"),
		Statistic:          cloudwatchtypes.StatisticAverage,
		Period:             aws.Int32(1),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudwatchtypes.ComparisonOperatorGreaterThanThreshold,
		AlarmActions:       []string{fnARN},
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_Lambda", func() error { return fmt.Errorf("put alarm: %w", err) })
	}
	defer func() {
		ic.cw.DeleteAlarms(ic.ctx, &cloudwatch.DeleteAlarmsInput{AlarmNames: []string{alarmName}})
	}()

	ic.cw.PutMetricData(ic.ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("AWS/EC2"),
		MetricData: []cloudwatchtypes.MetricDatum{{
			MetricName: aws.String("MemoryUtilization"),
			Value:      aws.Float64(100),
		}},
	})

	return r.pollVerify("CWAlarm_Lambda", defaultPollTimeout, func() error {
		alarmResp, err := ic.cw.DescribeAlarms(ic.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe alarms: %w", err)
		}
		if len(alarmResp.MetricAlarms) == 0 {
			return fmt.Errorf("alarm %s not found", alarmName)
		}
		if alarmResp.MetricAlarms[0].StateValue != cloudwatchtypes.StateValueAlarm {
			return fmt.Errorf("expected alarm state ALARM, got %s", alarmResp.MetricAlarms[0].StateValue)
		}
		return ic.verifyLambdaInvoked(fnName)
	})
}

// ---------- CloudWatch Alarm -> Step Functions ----------

func (r *TestRunner) runAlarmToStepFunctions(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-alarm-sfn-role-%s", ts)
	alarmName := fmt.Sprintf("integ-alarm-sfn-%s", ts)
	smName := fmt.Sprintf("integ-alarm-sfn-sm-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	_, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(`{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_StepFunctions", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)),
		})
	}()

	smARN := fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)

	_, err = ic.cw.PutMetricAlarm(ic.ctx, &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		MetricName:         aws.String("MemoryUtilization"),
		Namespace:          aws.String("AWS/EC2"),
		Statistic:          cloudwatchtypes.StatisticAverage,
		Period:             aws.Int32(1),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudwatchtypes.ComparisonOperatorGreaterThanThreshold,
		AlarmActions:       []string{smARN},
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_StepFunctions", func() error { return fmt.Errorf("put alarm: %w", err) })
	}
	defer func() {
		ic.cw.DeleteAlarms(ic.ctx, &cloudwatch.DeleteAlarmsInput{AlarmNames: []string{alarmName}})
	}()

	ic.cw.PutMetricData(ic.ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("AWS/EC2"),
		MetricData: []cloudwatchtypes.MetricDatum{{
			MetricName: aws.String("DiskSpaceUtilization"),
			Value:      aws.Float64(100),
		}},
	})

	return r.pollVerify("CWAlarm_StepFunctions", defaultPollTimeout, func() error {
		resp, err := ic.sfn.ListExecutions(ic.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Executions) == 0 {
			return fmt.Errorf("expected at least 1 execution from alarm, got 0")
		}
		return nil
	})
}

// ---------- Scheduler -> Lambda ----------

func (r *TestRunner) runSchedulerToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-sched-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-sched-role-%s", ts)
	lambdaRoleName := fmt.Sprintf("integ-sched-lambda-fn-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-lambda-%s", ts)
	groupName := fmt.Sprintf("integ-sched-group-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	IAMCreateRole(ic.iam, lambdaRoleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, lambdaRoleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	ic.createLambda(fnName, lambdaRoleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	_, err := ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(fnARN),
			RoleArn: aws.String(intRoleARN(roleName)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Lambda", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Lambda", schedulerPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

// ---------- Scheduler -> SQS ----------

func (r *TestRunner) runSchedulerToSQS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-sqs-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-sqs-%s", ts)
	groupName := fmt.Sprintf("integ-sched-sqs-group-%s", ts)
	queueName := fmt.Sprintf("integ-sched-sqs-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(queueARN),
			RoleArn: aws.String(intRoleARN(roleName)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SQS", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_SQS", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected message from scheduler in queue, got 0")
		}
		return nil
	})
}

// ---------- Scheduler -> SNS ----------

func (r *TestRunner) runSchedulerToSNS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-sns-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-sns-%s", ts)
	groupName := fmt.Sprintf("integ-sched-sns-group-%s", ts)
	topicName := fmt.Sprintf("integ-sched-sns-t-%s", ts)
	queueName := fmt.Sprintf("integ-sched-sns-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(topicARN),
			RoleArn: aws.String(intRoleARN(roleName)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SNS", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_SNS", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected message from scheduler (via SNS) in queue, got 0")
		}
		return nil
	})
}

// ---------- Scheduler -> Step Functions ----------

func (r *TestRunner) runSchedulerToStepFunctions(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-sfn-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-sfn-%s", ts)
	groupName := fmt.Sprintf("integ-sched-sfn-group-%s", ts)
	smName := fmt.Sprintf("integ-sched-sfn-sm-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	IAMCreateRole(ic.iam, fmt.Sprintf("%s-sfn", roleName), sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, fmt.Sprintf("%s-sfn", roleName))

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	_, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(fmt.Sprintf("%s-sfn", roleName))),
		Definition: aws.String(`{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_StepFunctions", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)),
		})
	}()

	smARN := fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(smARN),
			RoleArn: aws.String(intRoleARN(roleName)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_StepFunctions", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_StepFunctions", schedulerPollTimeout, func() error {
		resp, err := ic.sfn.ListExecutions(ic.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Executions) == 0 {
			return fmt.Errorf("expected at least 1 execution from scheduler, got 0")
		}
		return nil
	})
}

// ---------- Step Functions Task -> Lambda ----------

func (r *TestRunner) runSFNTaskLambda(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-lambda-role-%s", ts)
	lambdaRoleName := fmt.Sprintf("integ-sfn-lambda-fn-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-lambda-%s", ts)
	fnName := fmt.Sprintf("integ-sfn-lambda-fn-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	IAMCreateRole(ic.iam, lambdaRoleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, lambdaRoleName)

	ic.createLambda(fnName, lambdaRoleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	definition := fmt.Sprintf(`{
		"StartAt":"InvokeLambda",
		"States":{
			"InvokeLambda":{
				"Type":"Task",
				"Resource":"%s",
				"End":true
			}
		}
	}`, fnARN)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(definition),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_Lambda", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	smARN := createResp.StateMachineArn
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smARN})
	}()

	execResp, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: smARN,
		Input:           aws.String(`{"test":"sfn-task-lambda"}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_Lambda", func() error { return fmt.Errorf("start execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_Lambda", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: execResp.ExecutionArn,
		})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}
		return nil
	})
}

// ---------- Step Functions Task -> SQS ----------

func (r *TestRunner) runSFNTaskSQS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-sqs-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-sqs-%s", ts)
	queueName := fmt.Sprintf("integ-sfn-sqs-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	definition := `{
		"StartAt":"SendMsg",
		"States":{
			"SendMsg":{
				"Type":"Task",
				"Resource":"arn:aws:states:::sqs:sendMessage",
				"Parameters":{
					"QueueUrl":"QUEUE_URL_PLACEHOLDER",
					"MessageBody":{"test":"sfn-to-sqs"}
				},
				"End":true
			}
		}
	}`
	definition = strings.Replace(definition, "QUEUE_URL_PLACEHOLDER", queueURL, 1)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(definition),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SQS", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	smARN := createResp.StateMachineArn
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smARN})
	}()

	startResp, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: smARN,
		Input:           aws.String(`{}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SQS", func() error { return fmt.Errorf("start execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_SQS", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: startResp.ExecutionArn,
		})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return fmt.Errorf("receive messages: %w", err)
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected SQS message from SFN Task, got 0")
		}
		return nil
	})
}

// ---------- Step Functions Task -> SNS ----------

func (r *TestRunner) runSFNTaskSNS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-sns-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-sns-%s", ts)
	topicName := fmt.Sprintf("integ-sfn-sns-t-%s", ts)
	queueName := fmt.Sprintf("integ-sfn-sns-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	definition := fmt.Sprintf(`{
		"StartAt":"Publish",
		"States":{
			"Publish":{
				"Type":"Task",
				"Resource":"arn:aws:states:::sns:publish",
				"Parameters":{
					"TopicArn":"%s",
					"Message":{"test":"sfn-to-sns"}
				},
				"End":true
			}
		}
	}`, topicARN)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(definition),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	smARN := createResp.StateMachineArn
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smARN})
	}()

	startResp, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: smARN,
		Input:           aws.String(`{}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("start execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_SNS", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: startResp.ExecutionArn,
		})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return fmt.Errorf("receive messages: %w", err)
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected SNS message from SFN Task in queue, got 0")
		}
		return nil
	})
}

// ---------- S3 Notification -> Lambda ----------

func (r *TestRunner) runS3NotificationToLambda(ic *integClients, ts string) TestResult {
	bucketName := fmt.Sprintf("integ-s3-lambda-%s", strings.ToLower(ts))
	fnName := fmt.Sprintf("integ-s3-lambda-fn-%s", ts)
	roleName := fmt.Sprintf("integ-s3-lambda-role-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	err := ic.createBucket(bucketName)
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_Lambda", func() error { return fmt.Errorf("create bucket: %w", err) })
	}
	defer ic.deleteBucket(bucketName)

	_, err = ic.s3.PutBucketNotificationConfiguration(ic.ctx, &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(bucketName),
		NotificationConfiguration: &s3types.NotificationConfiguration{
			LambdaFunctionConfigurations: []s3types.LambdaFunctionConfiguration{
				{
					LambdaFunctionArn: aws.String(fnARN),
					Events:            []s3types.Event{s3types.EventS3ObjectCreatedPut},
				},
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_Lambda", func() error { return fmt.Errorf("put notification config: %w", err) })
	}

	ic.putObject(bucketName, "test-key.txt", []byte("test-data"))

	return r.pollVerify("S3_Notification_Lambda", defaultPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

// ---------- S3 Notification -> SQS ----------

func (r *TestRunner) runS3NotificationToSQS(ic *integClients, ts string) TestResult {
	bucketName := fmt.Sprintf("integ-s3-sqs-%s", strings.ToLower(ts))
	queueName := fmt.Sprintf("integ-s3-sqs-q-%s", ts)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)

	err = ic.createBucket(bucketName)
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_SQS", func() error { return fmt.Errorf("create bucket: %w", err) })
	}
	defer ic.deleteBucket(bucketName)

	_, err = ic.s3.PutBucketNotificationConfiguration(ic.ctx, &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(bucketName),
		NotificationConfiguration: &s3types.NotificationConfiguration{
			QueueConfigurations: []s3types.QueueConfiguration{
				{
					QueueArn: aws.String(queueARN),
					Events:   []s3types.Event{s3types.EventS3ObjectCreatedPut},
				},
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_SQS", func() error { return fmt.Errorf("put notification config: %w", err) })
	}

	ic.putObject(bucketName, "test-key.txt", []byte("test-data"))

	return r.pollVerify("S3_Notification_SQS", defaultPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected S3 notification in queue, got 0")
		}
		return nil
	})
}

// ---------- S3 Notification -> SNS ----------

func (r *TestRunner) runS3NotificationToSNS(ic *integClients, ts string) TestResult {
	bucketName := fmt.Sprintf("integ-s3-sns-%s", strings.ToLower(ts))
	topicName := fmt.Sprintf("integ-s3-sns-t-%s", ts)
	queueName := fmt.Sprintf("integ-s3-sns-q-%s", ts)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	err = ic.createBucket(bucketName)
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_SNS", func() error { return fmt.Errorf("create bucket: %w", err) })
	}
	defer ic.deleteBucket(bucketName)

	_, err = ic.s3.PutBucketNotificationConfiguration(ic.ctx, &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(bucketName),
		NotificationConfiguration: &s3types.NotificationConfiguration{
			TopicConfigurations: []s3types.TopicConfiguration{
				{
					TopicArn: aws.String(topicARN),
					Events:   []s3types.Event{s3types.EventS3ObjectCreatedPut},
				},
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "S3_Notification_SNS", func() error { return fmt.Errorf("put notification config: %w", err) })
	}

	ic.putObject(bucketName, "test-key.txt", []byte("test-data"))

	return r.pollVerify("S3_Notification_SNS", defaultPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected S3 notification (via SNS) in queue, got 0")
		}
		return nil
	})
}

// ---------- CloudWatch Logs -> Lambda ----------

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

// ---------- CloudWatch Logs -> Kinesis ----------

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

// ---------- SNS -> SQS ----------

func (r *TestRunner) runSNSToSQS(ic *integClients, ts string) TestResult {
	topicName := fmt.Sprintf("integ-sns-sqs-t-%s", ts)
	queueName := fmt.Sprintf("integ-sns-sqs-q-%s", ts)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "SNS_SQS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "SNS_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	_, err = ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})
	if err != nil {
		return r.RunTest(integSvc, "SNS_SQS", func() error { return fmt.Errorf("subscribe: %w", err) })
	}

	_, err = ic.sns.Publish(ic.ctx, &sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(`{"test":"sns-to-sqs"}`),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
			"TestType": {DataType: aws.String("String"), StringValue: aws.String("integration")},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "SNS_SQS", func() error { return fmt.Errorf("publish: %w", err) })
	}

	return r.pollVerify("SNS_SQS", defaultPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected message from SNS in queue, got 0")
		}
		return nil
	})
}

// ---------- SNS -> Lambda ----------

func (r *TestRunner) runSNSToLambda(ic *integClients, ts string) TestResult {
	topicName := fmt.Sprintf("integ-sns-lambda-t-%s", ts)
	fnName := fmt.Sprintf("integ-sns-lambda-fn-%s", ts)
	roleName := fmt.Sprintf("integ-sns-lambda-role-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "SNS_Lambda", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	_, err = ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("lambda"),
		Endpoint: aws.String(fnARN),
	})
	if err != nil {
		return r.RunTest(integSvc, "SNS_Lambda", func() error { return fmt.Errorf("subscribe: %w", err) })
	}

	_, err = ic.sns.Publish(ic.ctx, &sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(`{"test":"sns-to-lambda"}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SNS_Lambda", func() error { return fmt.Errorf("publish: %w", err) })
	}

	return r.pollVerify("SNS_Lambda", defaultPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}
