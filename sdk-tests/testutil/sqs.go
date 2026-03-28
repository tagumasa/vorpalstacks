package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSQSTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "sqs",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := sqs.NewFromConfig(cfg)
	ctx := context.Background()

	queueName := fmt.Sprintf("TestQueue-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("sqs", "CreateQueue", func() error {
		resp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		if resp.QueueUrl == nil {
			return fmt.Errorf("CreateQueue returned nil QueueUrl")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueUrl", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		if resp.QueueUrl == nil {
			return fmt.Errorf("GetQueueUrl returned nil QueueUrl")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueAttributes", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if attrResp.Attributes == nil {
			return fmt.Errorf("GetQueueAttributes returned nil Attributes")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessage", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		sendResp, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    resp.QueueUrl,
			MessageBody: aws.String("Test message"),
		})
		if err != nil {
			return err
		}
		if sendResp.MessageId == nil || *sendResp.MessageId == "" {
			return fmt.Errorf("SendMessage returned nil or empty MessageId")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("ReceiveMessage returned empty Messages list")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListQueues", func() error {
		resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{})
		if err != nil {
			return err
		}
		if resp.QueueUrls == nil {
			return fmt.Errorf("ListQueues returned nil QueueUrls")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SetQueueAttributes", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
			QueueUrl: resp.QueueUrl,
			Attributes: map[string]string{
				"VisibilityTimeout": "30",
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sqs", "TagQueue", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.TagQueue(ctx, &sqs.TagQueueInput{
			QueueUrl: resp.QueueUrl,
			Tags: map[string]string{
				"Environment": "Test",
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sqs", "ListQueueTags", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		tagResp, err := client.ListQueueTags(ctx, &sqs.ListQueueTagsInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if tagResp.Tags == nil || len(tagResp.Tags) == 0 {
			return fmt.Errorf("ListQueueTags returned nil or empty Tags (expected Environment=Test)")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "UntagQueue", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.UntagQueue(ctx, &sqs.UntagQueueInput{
			QueueUrl: resp.QueueUrl,
			TagKeys:  []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("sqs", "PurgeQueue", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: resp.QueueUrl,
		})
		return err
	}))

	results = append(results, r.RunTest("sqs", "ChangeMessageVisibility", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    resp.QueueUrl,
			MessageBody: aws.String("Test message for visibility"),
		})
		if err != nil {
			return err
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received from queue")
		}
		receiptHandle := aws.ToString(recvResp.Messages[0].ReceiptHandle)
		if receiptHandle == "" {
			return fmt.Errorf("receipt handle is empty")
		}
		_, err = client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
			QueueUrl:          resp.QueueUrl,
			ReceiptHandle:     aws.String(receiptHandle),
			VisibilityTimeout: 60,
		})
		return err
	}))

	results = append(results, r.RunTest("sqs", "DeleteQueue", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.DeleteQueue(ctx, &sqs.DeleteQueueInput{
			QueueUrl: resp.QueueUrl,
		})
		return err
	}))

	results = append(results, r.RunTest("sqs", "CreateQueue (FIFO)", func() error {
		fifoQueueName := fmt.Sprintf("TestFifoQueue-%d.fifo", time.Now().UnixNano())
		resp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(fifoQueueName),
			Attributes: map[string]string{
				"ContentBasedDeduplication": "true",
				"FifoQueue":                 "true",
			},
		})
		if err != nil {
			return err
		}
		if resp.QueueUrl == nil {
			return fmt.Errorf("CreateQueue (FIFO) returned nil QueueUrl")
		}
		return nil
	}))

	batchQueueName := fmt.Sprintf("TestBatchQueue-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("sqs", "SendMessageBatch", func() error {
		respCreate, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(batchQueueName),
		})
		if err != nil {
			return err
		}
		if respCreate.QueueUrl == nil {
			return fmt.Errorf("CreateQueue for batch returned nil QueueUrl")
		}
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(batchQueueName),
		})
		if err != nil {
			return err
		}
		batchResp, err := client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			QueueUrl: resp.QueueUrl,
			Entries: []types.SendMessageBatchRequestEntry{
				{
					Id:          aws.String("msg1"),
					MessageBody: aws.String("Batch message 1"),
				},
				{
					Id:          aws.String("msg2"),
					MessageBody: aws.String("Batch message 2"),
				},
			},
		})
		if err != nil {
			return err
		}
		if len(batchResp.Successful) == 0 {
			return fmt.Errorf("SendMessageBatch returned empty Successful entries")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DeleteMessageBatch", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(batchQueueName),
		})
		if err != nil {
			return err
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            resp.QueueUrl,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     2,
		})
		if err != nil {
			return err
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received for DeleteMessageBatch test")
		}
		var entries []types.DeleteMessageBatchRequestEntry
		for i, msg := range recvResp.Messages {
			entries = append(entries, types.DeleteMessageBatchRequestEntry{
				Id:            aws.String(fmt.Sprintf("del%d", i)),
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
		_, err = client.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
			QueueUrl: resp.QueueUrl,
			Entries:  entries,
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("sqs", "GetQueueUrl_NonExistent", func() error {
		_, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String("nonexistent-queue-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent queue")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessage_ReceiveRoundtrip", func() error {
		rtQueueName := fmt.Sprintf("RTQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(rtQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(rtQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(rtQueueName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		testBody := "roundtrip-test-message-12345"
		sendResp, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String(testBody),
		})
		if err != nil {
			return fmt.Errorf("send: %v", err)
		}
		if sendResp.MessageId == nil || *sendResp.MessageId == "" {
			return fmt.Errorf("message ID is nil or empty")
		}

		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: urlResp.QueueUrl,
		})
		if err != nil {
			return fmt.Errorf("receive: %v", err)
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received")
		}
		if *recvResp.Messages[0].Body != testBody {
			return fmt.Errorf("message body mismatch: got %q, want %q", *recvResp.Messages[0].Body, testBody)
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
			if len(u) > 0 {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected at least one queue URL")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ChangeMessageVisibility_NonExistent", func() error {
		_, err := client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
			QueueUrl:          aws.String("https://queue.amazonaws.com/000000000000/nonexistent"),
			ReceiptHandle:     aws.String("fake-receipt-handle"),
			VisibilityTimeout: 30,
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent receipt handle")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "CreateQueue_DuplicateName", func() error {
		dupQName := fmt.Sprintf("DupQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(dupQName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(dupQName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(dupQName),
		})
		if err != nil {
			return fmt.Errorf("duplicate queue name should be idempotent, got: %v", err)
		}
		return nil
	}))

	return results
}
