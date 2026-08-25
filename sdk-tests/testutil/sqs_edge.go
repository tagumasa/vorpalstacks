package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func (r *TestRunner) runSQSEdgeTests(ctx context.Context, client *sqs.Client, queueName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sqs", "GetQueueUrl_NonExistent", func() error {
		_, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String("nonexistent-queue-xyz"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DeleteQueue_NonExistent", func() error {
		_, err := client.DeleteQueue(ctx, &sqs.DeleteQueueInput{
			QueueUrl: aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "PurgeQueue_NonExistent", func() error {
		_, err := client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueAttributes_NonExistent", func() error {
		_, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessage_NonExistent", func() error {
		_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
			MessageBody: aws.String("test"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage_NonExistent", func() error {
		_, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SetQueueAttributes_NonExistent", func() error {
		_, err := client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
			QueueUrl: aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
			Attributes: map[string]string{
				"DelaySeconds": "10",
			},
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "TagQueue_NonExistent", func() error {
		_, err := client.TagQueue(ctx, &sqs.TagQueueInput{
			QueueUrl: aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
			Tags: map[string]string{
				"env": "test",
			},
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "RemovePermission_NonExistent", func() error {
		_, err := client.RemovePermission(ctx, &sqs.RemovePermissionInput{
			QueueUrl: aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
			Label:    aws.String("nonexistent-label"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DeleteMessage_NonExistent", func() error {
		_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
			ReceiptHandle: aws.String("fake-receipt-handle"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ChangeMessageVisibility_NonExistent", func() error {
		_, err := client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
			QueueUrl:          aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/nonexistent", r.accountID)),
			ReceiptHandle:     aws.String("fake-receipt-handle"),
			VisibilityTimeout: 30,
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	// MaxNumberOfMessages > 10 must be rejected (fail-closed).
	results = append(results, r.RunTest("sqs", "ReceiveMessage_MaxNumberOfMessages_TooHigh_Rejected", func() error {
		_, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(fmt.Sprintf("https://queue.amazonaws.com/%s/%s", r.accountID, queueName)),
			MaxNumberOfMessages: 11,
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return err
		}
		return nil
	}))

	// AddPermission with more than 7 actions must be rejected.
	results = append(results, r.RunTest("sqs", "AddPermission_TooManyActions_Rejected", func() error {
		queueURL := fmt.Sprintf("https://queue.amazonaws.com/%s/%s", r.accountID, queueName)
		actions := []string{
			"SendMessage", "ReceiveMessage", "DeleteMessage",
			"ChangeMessageVisibility", "GetQueueAttributes",
			"GetQueueUrl", "PurgeQueue", "SetQueueAttributes",
		}
		_, err := client.AddPermission(ctx, &sqs.AddPermissionInput{
			QueueUrl:      aws.String(queueURL),
			Label:         aws.String("too-many-actions-test"),
			AWSAccountIds: []string{r.accountID},
			Actions:       actions,
		})
		if err := AssertErrorContains(err, "OverLimit"); err != nil {
			return err
		}
		return nil
	}))

	// AddPermission with an invalid label character must be rejected.
	results = append(results, r.RunTest("sqs", "AddPermission_InvalidLabel_Rejected", func() error {
		queueURL := fmt.Sprintf("https://queue.amazonaws.com/%s/%s", r.accountID, queueName)
		_, err := client.AddPermission(ctx, &sqs.AddPermissionInput{
			QueueUrl:      aws.String(queueURL),
			Label:         aws.String("label with spaces!"),
			AWSAccountIds: []string{r.accountID},
			Actions:       []string{"SendMessage"},
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return err
		}
		return nil
	}))

	// CreateQueue with MaximumMessageSize > 1048576 must be rejected.
	results = append(results, r.RunTest("sqs", "CreateQueue_MaximumMessageSize_TooHigh_Rejected", func() error {
		qName := fmt.Sprintf("MsgSizeTest-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(qName),
			Attributes: map[string]string{
				"MaximumMessageSize": "1048577",
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return err
		}
		return nil
	}))

	// CreateQueue with MaximumMessageSize = 1048576 must succeed.
	results = append(results, r.RunTest("sqs", "CreateQueue_MaximumMessageSize_1MiB_Accepted", func() error {
		qName := fmt.Sprintf("MsgSize1MiB-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(qName),
			Attributes: map[string]string{
				"MaximumMessageSize": "1048576",
			},
		})
		if err != nil {
			return fmt.Errorf("create failed: %v", err)
		}
		urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qName)})
		if urlResp.QueueUrl != nil {
			client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
		}
		return nil
	}))

	// StartMessageMoveTask with rate > 500 must be rejected.
	results = append(results, r.RunTest("sqs", "StartMessageMoveTask_RateTooHigh_Rejected", func() error {
		_, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:                    aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:%s:nonexistent-dlq", r.accountID)),
			MaxNumberOfMessagesPerSecond: aws.Int32(501),
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return err
		}
		return nil
	}))

	// SendMessage with 11 message attributes must be rejected (AWS cap: 10).
	results = append(results, r.RunTest("sqs", "SendMessage_TooManyMessageAttributes_Rejected", func() error {
		qName := fmt.Sprintf("TooManyAttrs-%d", time.Now().UnixNano())
		createResp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(qName)})
		if err != nil {
			return fmt.Errorf("create failed: %v", err)
		}
		defer client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: createResp.QueueUrl})

		attrs := make(map[string]types.MessageAttributeValue, 11)
		for i := 0; i < 11; i++ {
			attrs[fmt.Sprintf("Attr%d", i)] = types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("v"),
			}
		}
		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:          createResp.QueueUrl,
			MessageBody:       aws.String("body"),
			MessageAttributes: attrs,
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return err
		}
		return nil
	}))

	// SendMessage attribute names starting with the reserved AWS. prefix
	// must be rejected.
	results = append(results, r.RunTest("sqs", "SendMessage_AttributeNameReservedPrefix_Rejected", func() error {
		qName := fmt.Sprintf("ReservedAttr-%d", time.Now().UnixNano())
		createResp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(qName)})
		if err != nil {
			return fmt.Errorf("create failed: %v", err)
		}
		defer client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: createResp.QueueUrl})

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    createResp.QueueUrl,
			MessageBody: aws.String("body"),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"AWS.Reserved": {
					DataType:    aws.String("String"),
					StringValue: aws.String("v"),
				},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return err
		}
		return nil
	}))

	// ListDeadLetterSourceQueues MaxResults above 1000 must be rejected.
	results = append(results, r.RunTest("sqs", "ListDeadLetterSourceQueues_MaxResultsTooHigh_Rejected", func() error {
		qName := fmt.Sprintf("DlqMax-%d", time.Now().UnixNano())
		createResp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(qName)})
		if err != nil {
			return fmt.Errorf("create failed: %v", err)
		}
		defer client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: createResp.QueueUrl})

		_, err = client.ListDeadLetterSourceQueues(ctx, &sqs.ListDeadLetterSourceQueuesInput{
			QueueUrl:   createResp.QueueUrl,
			MaxResults: aws.Int32(1001),
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
