package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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
			QueueUrl: aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "PurgeQueue_NonExistent", func() error {
		_, err := client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueAttributes_NonExistent", func() error {
		_, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessage_NonExistent", func() error {
		_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
			MessageBody: aws.String("test"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage_NonExistent", func() error {
		_, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SetQueueAttributes_NonExistent", func() error {
		_, err := client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
			QueueUrl: aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
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
			QueueUrl: aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
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
			QueueUrl: aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
			Label:    aws.String("nonexistent-label"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DeleteMessage_NonExistent", func() error {
		_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
			ReceiptHandle: aws.String("fake-receipt-handle"),
		})
		if err := AssertErrorContains(err, "ReceiptHandleIsInvalid"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ChangeMessageVisibility_NonExistent", func() error {
		_, err := client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
			QueueUrl:          aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
			ReceiptHandle:     aws.String("fake-receipt-handle"),
			VisibilityTimeout: 30,
		})
		if err := AssertErrorContains(err, "ReceiptHandleIsInvalid"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListQueues_ContainsCreated", func() error {
		lqName := fmt.Sprintf("LQTest-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(lqName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(lqName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{})
		if err != nil {
			return err
		}
		if resp.QueueUrls == nil {
			return fmt.Errorf("queue URLs is nil")
		}
		found := false
		for _, u := range resp.QueueUrls {
			if strings.Contains(u, lqName) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created queue %s not found in ListQueues", lqName)
		}
		return nil
	}))

	return results
}
