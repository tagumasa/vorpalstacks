package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

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
		return ic.verifyMessageContains(queueURL, bucketName)
	})
}

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
		return ic.verifyMessageContains(queueURL, bucketName)
	})
}
