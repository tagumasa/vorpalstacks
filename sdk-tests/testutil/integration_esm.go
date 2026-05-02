package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

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
