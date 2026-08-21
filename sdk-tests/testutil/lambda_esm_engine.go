package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	kinesistypes "github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"vorpalstacks-sdk-tests/config"
)

// runLambdaESMEngineTests exercises event source mapping delivery
// semantics end to end against real queues and streams: function-error
// acknowledgement, bisection, tumbling windows and parallelization.
func runLambdaESMEngineTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
	cwlClient *cloudwatchlogs.Client,
	createIAMRole func(string) error,
	deleteIAMRole func(string),
) []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "ESM_Engine_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to load config: %v", err)}}
	}
	sqsClient := sqs.NewFromConfig(cfg)
	kinesisClient := kinesis.NewFromConfig(cfg)
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	// waitForProcessingResult polls the mapping until its
	// LastProcessingResult contains want, failing fast when the poller
	// reports success instead.
	waitForProcessingResult := func(uuid, want, failOn string, timeout time.Duration) (string, error) {
		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			getResp, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{
				UUID: aws.String(uuid),
			})
			if err != nil {
				return "", fmt.Errorf("get mapping: %v", err)
			}
			last := aws.ToString(getResp.LastProcessingResult)
			if strings.Contains(last, want) {
				return last, nil
			}
			if failOn != "" && strings.Contains(last, failOn) {
				return last, fmt.Errorf("unexpected LastProcessingResult=%q", last)
			}
			time.Sleep(1 * time.Second)
		}
		return "", fmt.Errorf("timed out waiting for LastProcessingResult to contain %q", want)
	}

	results = append(results, r.RunTest("lambda", "ESM_FunctionError_NotAcknowledged", func() error {
		suffix := time.Now().UnixNano()
		fnName := fmt.Sprintf("EsmFailFn-%d", suffix)
		roleName := fmt.Sprintf("EsmFailRole-%d", suffix)
		queueName := fmt.Sprintf("esm-fail-queue-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		// A short visibility timeout returns the unacknowledged message to
		// the queue quickly so the poller keeps retrying it.
		qResp, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName:  aws.String(queueName),
			Attributes: map[string]string{"VisibilityTimeout": "2"},
		})
		if err != nil {
			return fmt.Errorf("create queue: %v", err)
		}
		defer sqsClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: qResp.QueueUrl})

		code, err := zipLambdaCode("exports.handler = async () => { throw new Error('poison'); };")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		if _, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    qResp.QueueUrl,
			MessageBody: aws.String(`{"id":"poison-1"}`),
		}); err != nil {
			return fmt.Errorf("send message: %v", err)
		}

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(fnName),
			EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:%s", r.region, r.accountID, queueName)),
			BatchSize:      aws.Int32(10),
			Enabled:        aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		// The mapping must report the function error as its last processing
		// result. A poller that acknowledges a failed batch deletes the
		// message and reports no errors instead.
		_, err = waitForProcessingResult(*mapping.UUID, "function error", "No errors", 30*time.Second)
		return err
	}))

	results = append(results, r.RunTest("lambda", "ESM_UpdateReplacesFunctionResponseTypes", func() error {
		suffix := time.Now().UnixNano()
		queueName := fmt.Sprintf("esm-frt-queue-%d", suffix)
		fnName := fmt.Sprintf("EsmFrtFn-%d", suffix)
		roleName := fmt.Sprintf("EsmFrtRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		qResp, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return fmt.Errorf("create queue: %v", err)
		}
		defer sqsClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: qResp.QueueUrl})

		code, err := zipLambdaCode("exports.handler = async () => ({});")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		// The mapping stays disabled: the response-type list is pure
		// configuration and needs no event flow.
		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:          aws.String(fnName),
			EventSourceArn:        aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:%s", r.region, r.accountID, queueName)),
			FunctionResponseTypes: []types.FunctionResponseType{types.FunctionResponseTypeReportBatchItemFailures},
			Enabled:               aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		for i := 0; i < 2; i++ {
			if _, err := client.UpdateEventSourceMapping(ctx, &lambda.UpdateEventSourceMappingInput{
				UUID:                  mapping.UUID,
				FunctionResponseTypes: []types.FunctionResponseType{types.FunctionResponseTypeReportBatchItemFailures},
			}); err != nil {
				return fmt.Errorf("update mapping: %v", err)
			}
		}

		final, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{
			UUID: mapping.UUID,
		})
		if err != nil {
			return err
		}
		if len(final.FunctionResponseTypes) != 1 {
			return fmt.Errorf("updates must replace the response-type list, got %d entries: %v",
				len(final.FunctionResponseTypes), final.FunctionResponseTypes)
		}
		if final.FunctionResponseTypes[0] != types.FunctionResponseTypeReportBatchItemFailures {
			return fmt.Errorf("unexpected response type %q", final.FunctionResponseTypes[0])
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_Kinesis_BisectDiscardsPoisonRecord", func() error {
		suffix := time.Now().UnixNano()
		streamName := fmt.Sprintf("esm-bisect-stream-%d", suffix)
		fnName := fmt.Sprintf("EsmBisectFn-%d", suffix)
		roleName := fmt.Sprintf("EsmBisectRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if _, err := kinesisClient.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		defer kinesisClient.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

		streamARN := ""
		deadline := time.Now().Add(10 * time.Second)
		for {
			out, err := kinesisClient.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				streamARN = aws.ToString(out.StreamDescription.StreamARN)
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("stream did not become ACTIVE: %v", err)
			}
			time.Sleep(150 * time.Millisecond)
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		// Kinesis record data reaches the function base64-encoded, matching
		// the documented Kinesis event format; the handler decodes before
		// comparing.
		code, err := zipLambdaCode(`exports.handler = async (event) => {
			for (const r of event.Records) {
				if (Buffer.from(r.kinesis.data, 'base64').toString() === 'poison') {
					throw new Error('poison record');
				}
			}
			return 'ok';
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		// Seven healthy records with a poison record in the middle of the
		// batch: bisection must isolate and discard only the poison record.
		for i := 1; i <= 7; i++ {
			if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
				StreamName:   aws.String(streamName),
				PartitionKey: aws.String(fmt.Sprintf("pk-%d", i)),
				Data:         []byte(fmt.Sprintf("ok-%d", i)),
			}); err != nil {
				return fmt.Errorf("put record %d: %v", i, err)
			}
		}
		if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(streamName),
			PartitionKey: aws.String("pk-poison"),
			Data:         []byte("poison"),
		}); err != nil {
			return fmt.Errorf("put poison record: %v", err)
		}

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:               aws.String(fnName),
			EventSourceArn:             aws.String(streamARN),
			StartingPosition:           types.EventSourcePositionTrimHorizon,
			BatchSize:                  aws.Int32(10),
			BisectBatchOnFunctionError: aws.Bool(true),
			MaximumRetryAttempts:       aws.Int32(0),
			Enabled:                    aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		// Bisection narrows the failure to the poison record, discards it
		// after the exhausted budget and advances past it.
		if _, err := waitForProcessingResult(*mapping.UUID, "discarded", "", 30*time.Second); err != nil {
			return err
		}

		// A record appended after the discard must process normally: the
		// shard is not blocked by the discarded poison record.
		if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(streamName),
			PartitionKey: aws.String("pk-after"),
			Data:         []byte("ok-after"),
		}); err != nil {
			return fmt.Errorf("put post-discard record: %v", err)
		}
		_, err = waitForProcessingResult(*mapping.UUID, "No errors", "", 30*time.Second)
		return err
	}))

	// waitForWindowEvent polls the function's CloudWatch logs for a logged
	// WINDOW_EVENT line whose parsed payload satisfies check.
	waitForWindowEvent := func(fnName string, check func(map[string]interface{}) bool, timeout time.Duration) (map[string]interface{}, error) {
		logGroupName := "/aws/lambda/" + fnName
		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			var nextToken *string
			for {
				out, err := cwlClient.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
					LogGroupName: aws.String(logGroupName),
					Limit:        aws.Int32(100),
					NextToken:    nextToken,
				})
				if err != nil {
					break
				}
				for _, ev := range out.Events {
					msg := aws.ToString(ev.Message)
					if !strings.HasPrefix(msg, "WINDOW_EVENT ") {
						continue
					}
					var parsed map[string]interface{}
					if json.Unmarshal([]byte(strings.TrimPrefix(msg, "WINDOW_EVENT ")), &parsed) != nil {
						continue
					}
					if check(parsed) {
						return parsed, nil
					}
				}
				if out.NextToken == nil {
					break
				}
				nextToken = out.NextToken
			}
			time.Sleep(1 * time.Second)
		}
		return nil, fmt.Errorf("timed out waiting for a matching WINDOW_EVENT log in %s", logGroupName)
	}

	results = append(results, r.RunTest("lambda", "ESM_Kinesis_OnFailureDestination_SQS", func() error {
		suffix := time.Now().UnixNano()
		streamName := fmt.Sprintf("esm-dest-stream-%d", suffix)
		fnName := fmt.Sprintf("EsmDestFn-%d", suffix)
		roleName := fmt.Sprintf("EsmDestRole-%d", suffix)
		queueName := fmt.Sprintf("esm-dest-queue-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if _, err := kinesisClient.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		defer kinesisClient.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

		streamARN := ""
		deadline := time.Now().Add(10 * time.Second)
		for {
			out, err := kinesisClient.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				streamARN = aws.ToString(out.StreamDescription.StreamARN)
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("stream did not become ACTIVE: %v", err)
			}
			time.Sleep(150 * time.Millisecond)
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		qResp, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(queueName)})
		if err != nil {
			return fmt.Errorf("create destination queue: %v", err)
		}
		defer sqsClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: qResp.QueueUrl})
		queueARN := fmt.Sprintf("arn:aws:sqs:%s:%s:%s", r.region, r.accountID, queueName)

		// The function always fails, so the whole batch is discarded once
		// the single configured retry attempt is spent.
		code, err := zipLambdaCode("exports.handler = async () => { throw new Error('always fails'); };")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		for i := 1; i <= 3; i++ {
			if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
				StreamName:   aws.String(streamName),
				PartitionKey: aws.String(fmt.Sprintf("pk-%d", i)),
				Data:         []byte(fmt.Sprintf("fail-%d", i)),
			}); err != nil {
				return fmt.Errorf("put record %d: %v", i, err)
			}
		}

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:         aws.String(fnName),
			EventSourceArn:       aws.String(streamARN),
			StartingPosition:     types.EventSourcePositionTrimHorizon,
			BatchSize:            aws.Int32(10),
			MaximumRetryAttempts: aws.Int32(0),
			DestinationConfig: &types.DestinationConfig{
				OnFailure: &types.OnFailure{Destination: aws.String(queueARN)},
			},
			Enabled: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		// The discarded batch must arrive on the destination queue as the
		// documented invocation record.
		var record map[string]interface{}
		waitDeadline := time.Now().Add(30 * time.Second)
		for record == nil && time.Now().Before(waitDeadline) {
			out, rerr := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:        qResp.QueueUrl,
				WaitTimeSeconds: 2,
			})
			if rerr != nil {
				return fmt.Errorf("receive from destination queue: %v", rerr)
			}
			for _, msg := range out.Messages {
				var parsed map[string]interface{}
				if json.Unmarshal([]byte(aws.ToString(msg.Body)), &parsed) != nil {
					continue
				}
				if _, ok := parsed["KinesisBatchInfo"]; ok {
					record = parsed
					break
				}
			}
		}
		if record == nil {
			return fmt.Errorf("no invocation record arrived on destination queue %s within 30s", queueARN)
		}

		reqCtx, ok := record["requestContext"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("record lacks requestContext: %v", record)
		}
		if reqCtx["condition"] != "RetryAttemptsExhausted" {
			return fmt.Errorf("condition = %v, want RetryAttemptsExhausted", reqCtx["condition"])
		}
		if reqCtx["functionArn"] == "" {
			return fmt.Errorf("functionArn empty in %v", reqCtx)
		}
		if reqCtx["approximateInvokeCount"] != float64(1) {
			return fmt.Errorf("approximateInvokeCount = %v, want 1 (MaximumRetryAttempts=0)", reqCtx["approximateInvokeCount"])
		}
		if record["version"] != "1.0" {
			return fmt.Errorf("version = %v, want 1.0", record["version"])
		}
		if _, has := record["payload"]; has {
			return fmt.Errorf("SQS destination record must not carry a payload member: %v", record)
		}
		batch, ok := record["KinesisBatchInfo"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("record lacks KinesisBatchInfo: %v", record)
		}
		if batch["batchSize"] != float64(3) {
			return fmt.Errorf("KinesisBatchInfo.batchSize = %v, want 3", batch["batchSize"])
		}
		if batch["streamArn"] != streamARN {
			return fmt.Errorf("KinesisBatchInfo.streamArn = %v, want %s", batch["streamArn"], streamARN)
		}
		if batch["shardId"] == "" || batch["startSequenceNumber"] == "" || batch["endSequenceNumber"] == "" {
			return fmt.Errorf("KinesisBatchInfo lacks identification: %v", batch)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_Tags_CreateTagUntag", func() error {
		suffix := time.Now().UnixNano()
		queueName := fmt.Sprintf("esm-tag-queue-%d", suffix)
		fnName := fmt.Sprintf("EsmTagFn-%d", suffix)
		roleName := fmt.Sprintf("EsmTagRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		qResp, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(queueName)})
		if err != nil {
			return fmt.Errorf("create queue: %v", err)
		}
		defer sqsClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: qResp.QueueUrl})

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode("exports.handler = async () => 'ok';")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)
		functionARN := fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", r.region, r.accountID, fnName)

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(fnName),
			EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:%s", r.region, r.accountID, queueName)),
			BatchSize:      aws.Int32(10),
			Tags: map[string]string{
				"Team": "core",
				"Env":  "esm",
			},
			Enabled: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})
		esmARN := fmt.Sprintf("arn:aws:lambda:%s:%s:event-source-mapping:%s", r.region, r.accountID, aws.ToString(mapping.UUID))

		// Tags supplied at creation are readable through the mapping ARN.
		listResp, err := client.ListTags(ctx, &lambda.ListTagsInput{Resource: aws.String(esmARN)})
		if err != nil {
			return fmt.Errorf("ListTags on event source mapping: %v", err)
		}
		if listResp.Tags["Team"] != "core" || listResp.Tags["Env"] != "esm" {
			return fmt.Errorf("creation tags missing, got %v", listResp.Tags)
		}

		// TagResource and UntagResource operate on the mapping ARN too.
		if _, err := client.TagResource(ctx, &lambda.TagResourceInput{
			Resource: aws.String(esmARN),
			Tags:     map[string]string{"Extra": "1"},
		}); err != nil {
			return fmt.Errorf("TagResource on event source mapping: %v", err)
		}
		if _, err := client.UntagResource(ctx, &lambda.UntagResourceInput{
			Resource: aws.String(esmARN),
			TagKeys:  []string{"Env"},
		}); err != nil {
			return fmt.Errorf("UntagResource on event source mapping: %v", err)
		}
		listResp, err = client.ListTags(ctx, &lambda.ListTagsInput{Resource: aws.String(esmARN)})
		if err != nil {
			return fmt.Errorf("ListTags after tag/untag: %v", err)
		}
		if listResp.Tags["Team"] != "core" || listResp.Tags["Extra"] != "1" {
			return fmt.Errorf("tag/untag result wrong, got %v", listResp.Tags)
		}
		if _, ok := listResp.Tags["Env"]; ok {
			return fmt.Errorf("untagged key Env still present: %v", listResp.Tags)
		}

		// "Lambda does not support adding tags to function aliases or
		// versions."
		if _, err := client.ListTags(ctx, &lambda.ListTagsInput{
			Resource: aws.String(functionARN + ":$LATEST"),
		}); err == nil {
			return fmt.Errorf("ListTags on a qualified function ARN must be rejected")
		} else if !strings.Contains(err.Error(), "InvalidParameterValueException") {
			return fmt.Errorf("qualified function ARN rejection = %v, want InvalidParameterValueException", err)
		}

		// Deleting the mapping drops its tags.
		if _, err := client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID}); err != nil {
			return fmt.Errorf("delete mapping: %v", err)
		}
		if _, err := client.ListTags(ctx, &lambda.ListTagsInput{Resource: aws.String(esmARN)}); err == nil {
			return fmt.Errorf("ListTags on a deleted mapping must fail")
		} else if !strings.Contains(err.Error(), "ResourceNotFoundException") {
			return fmt.Errorf("deleted mapping ListTags = %v, want ResourceNotFoundException", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_Kinesis_TumblingWindow_FinalInvoke", func() error {
		suffix := time.Now().UnixNano()
		streamName := fmt.Sprintf("esm-window-stream-%d", suffix)
		fnName := fmt.Sprintf("EsmWindowFn-%d", suffix)
		roleName := fmt.Sprintf("EsmWindowRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if _, err := kinesisClient.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		defer kinesisClient.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

		streamARN := ""
		deadline := time.Now().Add(10 * time.Second)
		for {
			out, err := kinesisClient.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				streamARN = aws.ToString(out.StreamDescription.StreamARN)
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("stream did not become ACTIVE: %v", err)
			}
			time.Sleep(150 * time.Millisecond)
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		// The handler echoes the whole time-window event into the logs and
		// answers with the state the response contract requires.
		code, err := zipLambdaCode(`exports.handler = async (event) => {
			console.log('WINDOW_EVENT ' + JSON.stringify(event));
			return { state: { n: (event.Records || []).length } };
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		// A one-second window makes every record appended a few seconds
		// apart land in a distinct window.
		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:            aws.String(fnName),
			EventSourceArn:          aws.String(streamARN),
			StartingPosition:        types.EventSourcePositionTrimHorizon,
			BatchSize:               aws.Int32(10),
			TumblingWindowInSeconds: aws.Int32(1),
			Enabled:                 aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(streamName),
			PartitionKey: aws.String("pk-w1"),
			Data:         []byte("w1"),
		}); err != nil {
			return fmt.Errorf("put first record: %v", err)
		}

		// Wait until the first record was aggregated mid-window before
		// appending the next one, so the two records provably sit in
		// different windows and the rollover delivers the final invoke.
		if _, err := waitForProcessingResult(*mapping.UUID, "No errors", "", 30*time.Second); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
		if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(streamName),
			PartitionKey: aws.String("pk-w2"),
			Data:         []byte("w2"),
		}); err != nil {
			return fmt.Errorf("put second record: %v", err)
		}

		// The final invocation for the first window carries the aggregated
		// state of the first record (n=1), no records, and the final flag.
		event, err := waitForWindowEvent(fnName, func(ev map[string]interface{}) bool {
			final, _ := ev["isFinalInvokeForWindow"].(bool)
			if !final {
				return false
			}
			state, _ := ev["state"].(map[string]interface{})
			n, _ := state["n"].(float64)
			return n == 1
		}, 30*time.Second)
		if err != nil {
			return err
		}
		window, _ := event["window"].(map[string]interface{})
		start, _ := window["start"].(string)
		end, _ := window["end"].(string)
		if start == "" || end == "" {
			return fmt.Errorf("final invoke carries no window boundaries: %+v", event)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_Kinesis_ParallelizationFactor_ExactlyOnce", func() error {
		suffix := time.Now().UnixNano()
		streamName := fmt.Sprintf("esm-pf-stream-%d", suffix)
		fnName := fmt.Sprintf("EsmPfFn-%d", suffix)
		roleName := fmt.Sprintf("EsmPfRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if _, err := kinesisClient.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		defer kinesisClient.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

		streamARN := ""
		deadline := time.Now().Add(10 * time.Second)
		for {
			out, err := kinesisClient.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				streamARN = aws.ToString(out.StreamDescription.StreamARN)
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("stream did not become ACTIVE: %v", err)
			}
			time.Sleep(150 * time.Millisecond)
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode(`exports.handler = async (event) => {
			const seen = (event.Records || []).map(r => Buffer.from(r.kinesis.data, 'base64').toString());
			console.log('PF_EVENT ' + JSON.stringify(seen));
			return {};
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:          aws.String(fnName),
			EventSourceArn:        aws.String(streamARN),
			StartingPosition:      types.EventSourcePositionTrimHorizon,
			BatchSize:             aws.Int32(5),
			ParallelizationFactor: aws.Int32(5),
			Enabled:               aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		expected := make([]string, 0, 12)
		for i := 1; i <= 12; i++ {
			data := fmt.Sprintf("pf-%d", i)
			expected = append(expected, data)
			if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
				StreamName:   aws.String(streamName),
				PartitionKey: aws.String(fmt.Sprintf("pk-%d", i)),
				Data:         []byte(data),
			}); err != nil {
				return fmt.Errorf("put record %d: %v", i, err)
			}
		}

		// The poller reports success only after every fetched batch was
		// delivered, so this pins completion of the concurrent delivery.
		if _, err := waitForProcessingResult(*mapping.UUID, "No errors", "", 30*time.Second); err != nil {
			return err
		}

		// A checkpoint resume that re-includes its boundary record would
		// deliver one record twice; the scan fails fast on any duplicate.
		// Full log coverage is best-effort: CloudWatch Logs ingestion can
		// lag a concurrent batch by tens of seconds (subscriber semaphore
		// skip plus outbox requeue), which is tracked separately, while a
		// genuine duplicate always shows up within seconds of the resume.
		logGroupName := "/aws/lambda/" + fnName
		scanDeadline := time.Now().Add(20 * time.Second)
		for time.Now().Before(scanDeadline) {
			counts := make(map[string]int)
			var nextToken *string
			for {
				out, err := cwlClient.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
					LogGroupName: aws.String(logGroupName),
					Limit:        aws.Int32(100),
					NextToken:    nextToken,
				})
				if err != nil {
					break
				}
				for _, ev := range out.Events {
					msg := aws.ToString(ev.Message)
					if !strings.HasPrefix(msg, "PF_EVENT ") {
						continue
					}
					var datas []string
					if json.Unmarshal([]byte(strings.TrimPrefix(msg, "PF_EVENT ")), &datas) != nil {
						continue
					}
					for _, d := range datas {
						counts[d]++
					}
				}
				if out.NextToken == nil {
					break
				}
				nextToken = out.NextToken
			}
			complete := true
			for _, want := range expected {
				switch counts[want] {
				case 0:
					complete = false
				case 1:
				default:
					return fmt.Errorf("record %q was delivered %d times", want, counts[want])
				}
			}
			if complete {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_BatchSizeOverTen_RequiresWindow", func() error {
		suffix := time.Now().UnixNano()
		queueName := fmt.Sprintf("esm-window-pair-queue-%d", suffix)
		fnName := fmt.Sprintf("EsmPairFn-%d", suffix)
		roleName := fmt.Sprintf("EsmPairRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode(`exports.handler = async () => ({});`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		qResp, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(queueName)})
		if err != nil {
			return fmt.Errorf("create queue: %v", err)
		}
		defer sqsClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: qResp.QueueUrl})

		queueARN := fmt.Sprintf("arn:aws:sqs:%s:%s:%s", r.region, r.accountID, queueName)

		// "when you set BatchSize to a value greater than 10, you must set
		// MaximumBatchingWindowInSeconds to at least 1" — the pairing must
		// be enforced on create.
		_, err = client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(fnName),
			EventSourceArn: aws.String(queueARN),
			BatchSize:      aws.Int32(100),
			Enabled:        aws.Bool(true),
		})
		if err == nil {
			return fmt.Errorf("BatchSize 100 without a batching window must be rejected")
		}
		if cerr := AssertErrorContains(err, "InvalidParameterValueException"); cerr != nil {
			return cerr
		}

		// The same pair with a window of one second is accepted.
		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:                   aws.String(fnName),
			EventSourceArn:                 aws.String(queueARN),
			BatchSize:                      aws.Int32(100),
			MaximumBatchingWindowInSeconds: aws.Int32(1),
			Enabled:                        aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("BatchSize 100 with a window must be accepted: %v", err)
		}
		client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		// A stored batch size above 10 also constrains a window-lowering
		// update.
		mapping2, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:                   aws.String(fnName),
			EventSourceArn:                 aws.String(queueARN),
			BatchSize:                      aws.Int32(100),
			MaximumBatchingWindowInSeconds: aws.Int32(1),
			Enabled:                        aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("re-create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping2.UUID})
		_, err = client.UpdateEventSourceMapping(ctx, &lambda.UpdateEventSourceMappingInput{
			UUID:                           mapping2.UUID,
			MaximumBatchingWindowInSeconds: aws.Int32(0),
		})
		if err == nil {
			return fmt.Errorf("lowering the window under a BatchSize above 10 must be rejected")
		}
		if cerr := AssertErrorContains(err, "InvalidParameterValueException"); cerr != nil {
			return cerr
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_Kinesis_BatchingWindow_HoldsPartialBatch", func() error {
		suffix := time.Now().UnixNano()
		streamName := fmt.Sprintf("esm-gather-stream-%d", suffix)
		fnName := fmt.Sprintf("EsmGatherFn-%d", suffix)
		roleName := fmt.Sprintf("EsmGatherRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if _, err := kinesisClient.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		defer kinesisClient.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

		streamARN := ""
		deadline := time.Now().Add(10 * time.Second)
		for {
			out, err := kinesisClient.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				streamARN = aws.ToString(out.StreamDescription.StreamARN)
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("stream did not become ACTIVE: %v", err)
			}
			time.Sleep(150 * time.Millisecond)
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode(`exports.handler = async () => ({});`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		// A five-second window with a batch size of ten: three records are
		// a partial batch and must be held until the window elapses.
		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:                   aws.String(fnName),
			EventSourceArn:                 aws.String(streamARN),
			StartingPosition:               types.EventSourcePositionTrimHorizon,
			BatchSize:                      aws.Int32(10),
			MaximumBatchingWindowInSeconds: aws.Int32(5),
			Enabled:                        aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		for i := 1; i <= 3; i++ {
			if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
				StreamName:   aws.String(streamName),
				PartitionKey: aws.String(fmt.Sprintf("gk-%d", i)),
				Data:         []byte(fmt.Sprintf("gather-%d", i)),
			}); err != nil {
				return fmt.Errorf("put record %d: %v", i, err)
			}
		}

		// Inside the window the records are gathered, not delivered: an
		// immediate delivery would report success within a couple of
		// seconds.
		time.Sleep(3 * time.Second)
		getResp, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{UUID: mapping.UUID})
		if err != nil {
			return fmt.Errorf("get mapping: %v", err)
		}
		if strings.Contains(aws.ToString(getResp.LastProcessingResult), "No errors") {
			return fmt.Errorf("partial batch was delivered inside the batching window (LastProcessingResult=%q)",
				aws.ToString(getResp.LastProcessingResult))
		}

		// Once the window elapses the gathered records flush.
		_, err = waitForProcessingResult(*mapping.UUID, "No errors", "", 15*time.Second)
		return err
	}))

	// deliveredIds scans the function's CloudWatch logs for every DELIVERED
	// line the handler emitted, keyed by the delivered record id.
	deliveredIds := func(fnName string) (map[string]bool, error) {
		ids := map[string]bool{}
		var nextToken *string
		for {
			out, err := cwlClient.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName: aws.String("/aws/lambda/" + fnName),
				Limit:        aws.Int32(100),
				NextToken:    nextToken,
			})
			if err != nil {
				return nil, err
			}
			for _, ev := range out.Events {
				msg := aws.ToString(ev.Message)
				if strings.HasPrefix(msg, "DELIVERED ") {
					ids[strings.TrimPrefix(msg, "DELIVERED ")] = true
				}
			}
			if out.NextToken == nil {
				break
			}
			nextToken = out.NextToken
		}
		return ids, nil
	}

	// deliveredCounts counts every DELIVERED line per id, so a retry loop
	// (a record the function keeps reporting as failed) is distinguishable
	// from a one-shot delivery.
	deliveredCounts := func(fnName string) (map[string]int, error) {
		counts := map[string]int{}
		var nextToken *string
		for {
			out, err := cwlClient.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName: aws.String("/aws/lambda/" + fnName),
				Limit:        aws.Int32(100),
				NextToken:    nextToken,
			})
			if err != nil {
				return nil, err
			}
			for _, ev := range out.Events {
				msg := aws.ToString(ev.Message)
				if strings.HasPrefix(msg, "DELIVERED ") {
					counts[strings.TrimPrefix(msg, "DELIVERED ")]++
				}
			}
			if out.NextToken == nil {
				break
			}
			nextToken = out.NextToken
		}
		return counts, nil
	}

	waitForDelivery := func(fnName, id string, timeout time.Duration) error {
		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			ids, err := deliveredIds(fnName)
			if err == nil && ids[id] {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
		return fmt.Errorf("timed out waiting for delivery of %q", id)
	}

	waitForDeliveryCount := func(fnName, id string, min int, timeout time.Duration) error {
		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			counts, err := deliveredCounts(fnName)
			if err == nil && counts[id] >= min {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
		return fmt.Errorf("timed out waiting for %d deliveries of %q", min, id)
	}

	// createStreamTable provisions a DynamoDB table with a NEW_AND_OLD_IMAGES
	// stream and waits for it to become active.
	createStreamTable := func(tableName string) (string, error) {
		createResp, err := dynamodbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			},
			KeySchema: []dynamodbtypes.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: dynamodbtypes.KeyTypeHash},
			},
			BillingMode: dynamodbtypes.BillingModePayPerRequest,
			StreamSpecification: &dynamodbtypes.StreamSpecification{
				StreamEnabled:  aws.Bool(true),
				StreamViewType: dynamodbtypes.StreamViewTypeNewAndOldImages,
			},
		})
		if err != nil {
			return "", fmt.Errorf("create table: %v", err)
		}
		streamARN := aws.ToString(createResp.TableDescription.LatestStreamArn)
		deadline := time.Now().Add(15 * time.Second)
		for {
			out, derr := dynamodbClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
			if derr == nil && out.Table.TableStatus == dynamodbtypes.TableStatusActive {
				return streamARN, nil
			}
			if time.Now().After(deadline) {
				return "", fmt.Errorf("table did not become ACTIVE: %v", derr)
			}
			time.Sleep(200 * time.Millisecond)
		}
	}

	results = append(results, r.RunTest("lambda", "ESM_DynamoDB_LatestSkipsHistory", func() error {
		suffix := time.Now().UnixNano()
		tableName := fmt.Sprintf("esm-ddb-latest-%d", suffix)
		fnName := fmt.Sprintf("EsmDdbLatestFn-%d", suffix)
		roleName := fmt.Sprintf("EsmDdbLatestRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		streamARN, err := createStreamTable(tableName)
		if err != nil {
			return err
		}
		defer dynamodbClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

		put := func(id string) error {
			_, err := dynamodbClient.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item:      map[string]dynamodbtypes.AttributeValue{"id": &dynamodbtypes.AttributeValueMemberS{Value: id}},
			})
			if err != nil {
				return fmt.Errorf("put %s: %v", id, err)
			}
			return nil
		}

		// History that predates the mapping. Under LATEST the newest record
		// is the anchor and is delivered; the older records are skipped.
		for _, id := range []string{"ddb-old-1", "ddb-old-2", "ddb-anchor"} {
			if err := put(id); err != nil {
				return err
			}
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode(`exports.handler = async (event) => {
			for (const r of event.Records) {
				console.log('DELIVERED ' + r.dynamodb.NewImage.id.S);
			}
			return 'ok';
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:     aws.String(fnName),
			EventSourceArn:   aws.String(streamARN),
			StartingPosition: types.EventSourcePositionLatest,
			BatchSize:        aws.Int32(10),
			Enabled:          aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		if err := waitForDelivery(fnName, "ddb-anchor", 60*time.Second); err != nil {
			return err
		}
		if err := put("ddb-after"); err != nil {
			return err
		}
		if err := waitForDelivery(fnName, "ddb-after", 60*time.Second); err != nil {
			return err
		}

		ids, err := deliveredIds(fnName)
		if err != nil {
			return fmt.Errorf("scan deliveries: %v", err)
		}
		for _, old := range []string{"ddb-old-1", "ddb-old-2"} {
			if ids[old] {
				return fmt.Errorf("record %q predating the mapping must not be delivered under LATEST", old)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_DynamoDB_RejectsAtTimestamp", func() error {
		suffix := time.Now().UnixNano()
		tableName := fmt.Sprintf("esm-ddb-ats-%d", suffix)
		fnName := fmt.Sprintf("EsmDdbAtsFn-%d", suffix)
		roleName := fmt.Sprintf("EsmDdbAtsRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		streamARN, err := createStreamTable(tableName)
		if err != nil {
			return err
		}
		defer dynamodbClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode(`exports.handler = async () => 'ok';`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		// AT_TIMESTAMP is a stream start position for Kinesis-family
		// sources only; a DynamoDB stream source must reject it.
		_, err = client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:              aws.String(fnName),
			EventSourceArn:            aws.String(streamARN),
			StartingPosition:          types.EventSourcePositionAtTimestamp,
			StartingPositionTimestamp: aws.Time(time.Now().Add(-time.Minute)),
			Enabled:                   aws.Bool(true),
		})
		if err == nil {
			return fmt.Errorf("AT_TIMESTAMP must be rejected for DynamoDB stream sources")
		}
		return AssertErrorContains(err, "InvalidParameterValueException")
	}))

	results = append(results, r.RunTest("lambda", "ESM_Kinesis_PartialBatchResponse", func() error {
		suffix := time.Now().UnixNano()
		streamName := fmt.Sprintf("esm-pbr-stream-%d", suffix)
		fnName := fmt.Sprintf("EsmPbrFn-%d", suffix)
		roleName := fmt.Sprintf("EsmPbrRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if _, err := kinesisClient.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		defer kinesisClient.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

		streamARN := ""
		deadline := time.Now().Add(10 * time.Second)
		for {
			out, err := kinesisClient.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				streamARN = aws.ToString(out.StreamDescription.StreamARN)
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("stream did not become ACTIVE: %v", err)
			}
			time.Sleep(150 * time.Millisecond)
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		// The handler reports the poison record's sequence number in its
		// partial batch response; everything else is a success.
		code, err := zipLambdaCode(`exports.handler = async (event) => {
			const failures = [];
			for (const r of event.Records) {
				const data = Buffer.from(r.kinesis.data, 'base64').toString();
				console.log('DELIVERED ' + data);
				if (data === 'pbr-poison') {
					failures.push({ itemIdentifier: r.kinesis.sequenceNumber });
				}
			}
			return { batchItemFailures: failures };
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:     aws.String(fnName),
			EventSourceArn:   aws.String(streamARN),
			StartingPosition: types.EventSourcePositionTrimHorizon,
			BatchSize:        aws.Int32(10),
			Enabled:          aws.Bool(true),
			FunctionResponseTypes: []types.FunctionResponseType{
				types.FunctionResponseTypeReportBatchItemFailures,
			},
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		for _, rec := range []string{"pbr-ok-1", "pbr-ok-2", "pbr-poison"} {
			if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
				StreamName:   aws.String(streamName),
				PartitionKey: aws.String("pk-" + rec),
				Data:         []byte(rec),
			}); err != nil {
				return fmt.Errorf("put record %s: %v", rec, err)
			}
		}

		// The poison record is retried: the checkpoint cannot advance past a
		// reported failure, so the record is re-delivered.
		if err := waitForDeliveryCount(fnName, "pbr-poison", 2, 60*time.Second); err != nil {
			return err
		}

		// The successfully processed records must not be re-delivered once
		// the reported failure holds the checkpoint behind them.
		counts, err := deliveredCounts(fnName)
		if err != nil {
			return fmt.Errorf("scan deliveries: %v", err)
		}
		for _, ok := range []string{"pbr-ok-1", "pbr-ok-2"} {
			if counts[ok] != 1 {
				return fmt.Errorf("record %q must be delivered exactly once, got %d deliveries", ok, counts[ok])
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_SQS_PartialBatchResponse", func() error {
		suffix := time.Now().UnixNano()
		queueName := fmt.Sprintf("esm-pbr-queue-%d", suffix)
		fnName := fmt.Sprintf("EsmPbrSqsFn-%d", suffix)
		roleName := fmt.Sprintf("EsmPbrSqsRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		// A short visibility timeout returns the reported message quickly so
		// the retry loop is observable within the test window.
		qResp, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName:  aws.String(queueName),
			Attributes: map[string]string{"VisibilityTimeout": "2"},
		})
		if err != nil {
			return fmt.Errorf("create queue: %v", err)
		}
		defer sqsClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: qResp.QueueUrl})

		// The handler reports the poison message id; the messages it does
		// not report are deleted from the queue.
		code, err := zipLambdaCode(`exports.handler = async (event) => {
			const failures = [];
			for (const m of event.Records) {
				console.log('DELIVERED ' + m.body);
				if (m.body === 'pbr-poison') {
					failures.push({ itemIdentifier: m.messageId });
				}
			}
			return { batchItemFailures: failures };
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(fnName),
			EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:%s", r.region, r.accountID, queueName)),
			BatchSize:      aws.Int32(10),
			Enabled:        aws.Bool(true),
			FunctionResponseTypes: []types.FunctionResponseType{
				types.FunctionResponseTypeReportBatchItemFailures,
			},
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		for _, body := range []string{"pbr-ok-1", "pbr-poison", "pbr-ok-2"} {
			if _, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
				QueueUrl:    qResp.QueueUrl,
				MessageBody: aws.String(body),
			}); err != nil {
				return fmt.Errorf("send %s: %v", body, err)
			}
		}

		// The reported message becomes visible again: the poller keeps
		// retrying only it.
		if err := waitForDeliveryCount(fnName, "pbr-poison", 2, 60*time.Second); err != nil {
			return err
		}

		// The unreported messages are deleted: they are delivered exactly
		// once while the poison message loops.
		counts, err := deliveredCounts(fnName)
		if err != nil {
			return fmt.Errorf("scan deliveries: %v", err)
		}
		for _, ok := range []string{"pbr-ok-1", "pbr-ok-2"} {
			if counts[ok] != 1 {
				return fmt.Errorf("message %q must be delivered exactly once, got %d deliveries", ok, counts[ok])
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_Kinesis_ParallelizationFactor_PartialBatchResponse", func() error {
		suffix := time.Now().UnixNano()
		streamName := fmt.Sprintf("esm-pbrpf-stream-%d", suffix)
		fnName := fmt.Sprintf("EsmPbrPfFn-%d", suffix)
		roleName := fmt.Sprintf("EsmPbrPfRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		if _, err := kinesisClient.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		defer kinesisClient.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

		streamARN := ""
		deadline := time.Now().Add(10 * time.Second)
		for {
			out, err := kinesisClient.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				streamARN = aws.ToString(out.StreamDescription.StreamARN)
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("stream did not become ACTIVE: %v", err)
			}
			time.Sleep(150 * time.Millisecond)
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode(`exports.handler = async (event) => {
			const failures = [];
			for (const r of event.Records) {
				const data = Buffer.from(r.kinesis.data, 'base64').toString();
				console.log('DELIVERED ' + data);
				if (data === 'pbrpf-poison') {
					failures.push({ itemIdentifier: r.kinesis.sequenceNumber });
				}
			}
			return { batchItemFailures: failures };
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		// Two concurrent batches per shard: the first carries the poison
		// record, the second succeeds — the checkpoint must stay at the
		// partial cursor of the first batch.
		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:          aws.String(fnName),
			EventSourceArn:        aws.String(streamARN),
			StartingPosition:      types.EventSourcePositionTrimHorizon,
			BatchSize:             aws.Int32(2),
			ParallelizationFactor: aws.Int32(2),
			Enabled:               aws.Bool(true),
			FunctionResponseTypes: []types.FunctionResponseType{
				types.FunctionResponseTypeReportBatchItemFailures,
			},
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		for _, rec := range []string{"pbrpf-ok-1", "pbrpf-poison", "pbrpf-ok-2", "pbrpf-ok-3"} {
			if _, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
				StreamName:   aws.String(streamName),
				PartitionKey: aws.String("pk-" + rec),
				Data:         []byte(rec),
			}); err != nil {
				return fmt.Errorf("put record %s: %v", rec, err)
			}
		}

		// The reported record is retried: a concurrently succeeding batch
		// must not move the checkpoint past the reported failure.
		if err := waitForDeliveryCount(fnName, "pbrpf-poison", 2, 60*time.Second); err != nil {
			return err
		}

		// The partial response must surface in the processing result
		// instead of a clean success.
		if _, err := waitForProcessingResult(aws.ToString(mapping.UUID), "batchItemFailures", "", 30*time.Second); err != nil {
			return err
		}

		// The record before every clamp is delivered exactly once; records
		// after the checkpoint are re-read with the poison (at-least-once).
		counts, err := deliveredCounts(fnName)
		if err != nil {
			return fmt.Errorf("scan deliveries: %v", err)
		}
		if counts["pbrpf-ok-1"] != 1 {
			return fmt.Errorf("record %q must be delivered exactly once, got %d deliveries", "pbrpf-ok-1", counts["pbrpf-ok-1"])
		}
		for _, ok := range []string{"pbrpf-ok-2", "pbrpf-ok-3"} {
			if counts[ok] < 1 {
				return fmt.Errorf("record %q must be delivered at least once, got %d deliveries", ok, counts[ok])
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ESM_DynamoDB_BatchingWindow_PartialBatchResponse", func() error {
		suffix := time.Now().UnixNano()
		tableName := fmt.Sprintf("esm-pbrbw-table-%d", suffix)
		fnName := fmt.Sprintf("EsmPbrBwFn-%d", suffix)
		roleName := fmt.Sprintf("EsmPbrBwRole-%d", suffix)
		roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

		streamARN, err := createStreamTable(tableName)
		if err != nil {
			return err
		}
		defer dynamodbClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

		put := func(id string) error {
			_, err := dynamodbClient.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item:      map[string]dynamodbtypes.AttributeValue{"id": &dynamodbtypes.AttributeValueMemberS{Value: id}},
			})
			return err
		}

		if err := createIAMRole(roleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(roleName)

		code, err := zipLambdaCode(`exports.handler = async (event) => {
			const failures = [];
			for (const r of event.Records) {
				console.log('DELIVERED ' + r.dynamodb.NewImage.id.S);
				if (r.dynamodb.NewImage.id.S === 'pbrbw-poison') {
					failures.push({ itemIdentifier: r.dynamodb.SequenceNumber });
				}
			}
			return { batchItemFailures: failures };
		};`)
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: code},
		}); err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(cwlClient, ctx, fnName)

		// The batching window gathers four records into two chunks; the
		// first chunk carries the poison record — its partial cursor must
		// stop the flush, so the second chunk waits for the retry instead
		// of jumping the checkpoint past the failure.
		mapping, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:                   aws.String(fnName),
			EventSourceArn:                 aws.String(streamARN),
			StartingPosition:               types.EventSourcePositionTrimHorizon,
			BatchSize:                      aws.Int32(2),
			ParallelizationFactor:          aws.Int32(2),
			MaximumBatchingWindowInSeconds: aws.Int32(1),
			Enabled:                        aws.Bool(true),
			FunctionResponseTypes: []types.FunctionResponseType{
				types.FunctionResponseTypeReportBatchItemFailures,
			},
		})
		if err != nil {
			return fmt.Errorf("create mapping: %v", err)
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: mapping.UUID})

		for _, id := range []string{"pbrbw-ok-1", "pbrbw-poison", "pbrbw-ok-2", "pbrbw-ok-3"} {
			if err := put(id); err != nil {
				return fmt.Errorf("put %s: %v", id, err)
			}
		}

		// The reported record is retried: the flush checkpoint must stop
		// at the partial cursor of the first chunk.
		if err := waitForDeliveryCount(fnName, "pbrbw-poison", 2, 60*time.Second); err != nil {
			return err
		}

		// The partial response must surface in the processing result.
		if _, err := waitForProcessingResult(aws.ToString(mapping.UUID), "batchItemFailures", "", 30*time.Second); err != nil {
			return err
		}

		counts, err := deliveredCounts(fnName)
		if err != nil {
			return fmt.Errorf("scan deliveries: %v", err)
		}
		if counts["pbrbw-ok-1"] != 1 {
			return fmt.Errorf("record %q must be delivered exactly once, got %d deliveries", "pbrbw-ok-1", counts["pbrbw-ok-1"])
		}
		return nil
	}))

	return results
}
