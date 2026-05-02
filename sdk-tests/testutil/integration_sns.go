package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

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
		return ic.verifyMessageContains(queueURL, "sns-to-sqs")
	})
}

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
		return ic.verifyLambdaLogContains(fnName, "sns-to-lambda")
	})
}
