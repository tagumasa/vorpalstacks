package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iotTypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
	iotdataplane "github.com/aws/aws-sdk-go-v2/service/iotdataplane"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"vorpalstacks-sdk-tests/config"
)

// runIoTIntegrationRuleActionTests creates a TopicRule with an SQS action,
// publishes a message via the iot-data-plane Publish endpoint, and verifies
// the message arrives in the SQS queue. This exercises the full chain:
// HTTP Publish → rule executor → SQL evaluation → action dispatcher → SQS.
func (r *TestRunner) runIoTIntegrationRuleActionTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "Integration_TopicRule_SQSAction", func() error {
		ctx, cancel := context.WithTimeout(tc.ctx, 30*time.Second)
		defer cancel()

		// 1. Create SQS queue.
		sqsCfg, _ := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		sqsClient := sqs.NewFromConfig(sqsCfg)
		queueName := uniqueName("iot-rule-sqs")
		createOut, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return fmt.Errorf("CreateQueue: %w", err)
		}
		queueURL := *createOut.QueueUrl
		defer func() {
			_, _ = sqsClient.DeleteQueue(tc.ctx, &sqs.DeleteQueueInput{QueueUrl: aws.String(queueURL)})
		}()

		// Create TopicRule with SQS action.
		ruleName := uniqueName("iot-rule-sqs")
		topicFilter := "test/integration/" + queueName
		_, err = tc.client.CreateTopicRule(ctx, &iot.CreateTopicRuleInput{
			RuleName: aws.String(ruleName),
			TopicRulePayload: &iotTypes.TopicRulePayload{
				Sql:          aws.String(fmt.Sprintf("SELECT * FROM '%s'", topicFilter)),
				RuleDisabled: aws.Bool(false),
				Actions: []iotTypes.Action{{
					Sqs: &iotTypes.SqsAction{
						RoleArn:   aws.String("arn:aws:iam::000000000000:role/test"),
						QueueUrl:  aws.String(queueURL),
						UseBase64: aws.Bool(false),
					},
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateTopicRule: %w", err)
		}
		defer func() {
			_, _ = tc.client.DeleteTopicRule(tc.ctx, &iot.DeleteTopicRuleInput{RuleName: aws.String(ruleName)})
		}()

		// Wait for the rule to be registered in the executor.
		time.Sleep(2 * time.Second)

		// 3. Publish a message via the iot-data-plane Publish endpoint.
		dataClient := iotdataplane.NewFromConfig(sqsCfg)
		payload := `{"device":"sensor-1","value":42}`
		_, err = dataClient.Publish(ctx, &iotdataplane.PublishInput{
			Topic:   aws.String(topicFilter),
			Payload: []byte(payload),
		})
		if err != nil {
			return fmt.Errorf("iotdataplane.Publish: %w", err)
		}

		// 4. Poll SQS for the message.
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			recvOut, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     1,
			})
			if err != nil {
				continue
			}
			for _, msg := range recvOut.Messages {
				if msg.Body != nil && *msg.Body != "" {
					return nil // Message received — integration verified.
				}
			}
		}
		return fmt.Errorf("message not received in SQS queue after 10s")
	}))

	results = append(results, r.RunTest("iot", "Integration_TopicRule_CloudWatchLogsAction", func() error {
		ctx, cancel := context.WithTimeout(tc.ctx, 30*time.Second)
		defer cancel()

		cwlCfg, _ := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		cwlClient := cloudwatchlogs.NewFromConfig(cwlCfg)

		// The CloudWatch Logs action requires a pre-existing log group (the
		// AWS role permissions cover CreateLogStream, not CreateLogGroup).
		logGroupName := "/iot-test/" + uniqueName("rule-logs")
		if _, err := cwlClient.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		}); err != nil {
			return fmt.Errorf("CreateLogGroup: %w", err)
		}
		defer func() {
			_, _ = cwlClient.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
				LogGroupName: aws.String(logGroupName),
			})
		}()

		// Create a TopicRule with a CloudWatch Logs action targeting the
		// configured group.
		ruleName := uniqueName("iot-rule-cwl")
		topicFilter := "test/integration/cwl/" + ruleName
		_, err := tc.client.CreateTopicRule(ctx, &iot.CreateTopicRuleInput{
			RuleName: aws.String(ruleName),
			TopicRulePayload: &iotTypes.TopicRulePayload{
				Sql:          aws.String(fmt.Sprintf("SELECT * FROM '%s'", topicFilter)),
				RuleDisabled: aws.Bool(false),
				Actions: []iotTypes.Action{{
					CloudwatchLogs: &iotTypes.CloudwatchLogsAction{
						LogGroupName: aws.String(logGroupName),
						RoleArn:      aws.String("arn:aws:iam::000000000000:role/test"),
					},
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateTopicRule: %w", err)
		}
		defer func() {
			_, _ = tc.client.DeleteTopicRule(tc.ctx, &iot.DeleteTopicRuleInput{RuleName: aws.String(ruleName)})
		}()

		// Wait for the rule to be registered in the executor.
		time.Sleep(2 * time.Second)

		dataClient := iotdataplane.NewFromConfig(cwlCfg)
		payload := `{"device":"sensor-cwl","value":7}`
		if _, err := dataClient.Publish(ctx, &iotdataplane.PublishInput{
			Topic:   aws.String(topicFilter),
			Payload: []byte(payload),
		}); err != nil {
			return fmt.Errorf("iotdataplane.Publish: %w", err)
		}

		// Poll the CloudWatch Logs API: the message must be readable from the
		// configured group, immediately, through the API read plane.
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			streams, err := cwlClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
				LogGroupName: aws.String(logGroupName),
			})
			if err != nil {
				continue
			}
			for _, stream := range streams.LogStreams {
				if stream.LogStreamName == nil {
					continue
				}
				out, err := cwlClient.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
					LogGroupName:  aws.String(logGroupName),
					LogStreamName: stream.LogStreamName,
					StartFromHead: aws.Bool(true),
				})
				if err != nil {
					continue
				}
				for _, ev := range out.Events {
					if ev.Message != nil && strings.Contains(*ev.Message, "sensor-cwl") {
						return nil
					}
				}
			}
		}
		return fmt.Errorf("event not visible in configured log group %s after 10s", logGroupName)
	}))

	return results
}
