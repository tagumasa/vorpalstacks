package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	return results
}
