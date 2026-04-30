package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func (r *TestRunner) runSQSMessageTests(ctx context.Context, client *sqs.Client, queueName string) []TestResult {
	var results []TestResult

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

	results = append(results, r.RunTest("sqs", "SendMessage_WithDelaySeconds", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		sendResp, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:     resp.QueueUrl,
			MessageBody:  aws.String("Delayed message"),
			DelaySeconds: 5,
		})
		if err != nil {
			return err
		}
		if sendResp.MessageId == nil || *sendResp.MessageId == "" {
			return fmt.Errorf("SendMessage with DelaySeconds returned nil MessageId")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessage_WithMessageAttributes", func() error {
		attrQueueName := fmt.Sprintf("AttrQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(attrQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(attrQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(attrQueueName),
		})
		if err != nil {
			return err
		}
		sendResp, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    resp.QueueUrl,
			MessageBody: aws.String("Message with attributes"),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"Attr1": {
					DataType:    aws.String("String"),
					StringValue: aws.String("value1"),
				},
				"Attr2": {
					DataType:    aws.String("Number"),
					StringValue: aws.String("42"),
				},
			},
		})
		if err != nil {
			return err
		}
		if sendResp.MessageId == nil || *sendResp.MessageId == "" {
			return fmt.Errorf("SendMessage with MessageAttributes returned nil MessageId")
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:              resp.QueueUrl,
			MessageAttributeNames: []string{"All"},
		})
		if err != nil {
			return fmt.Errorf("receive: %v", err)
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received")
		}
		msg := recvResp.Messages[0]
		if len(msg.MessageAttributes) < 2 {
			return fmt.Errorf("expected at least 2 message attributes, got %d", len(msg.MessageAttributes))
		}
		if msg.MessageAttributes["Attr1"].StringValue == nil || *msg.MessageAttributes["Attr1"].StringValue != "value1" {
			return fmt.Errorf("Attr1 mismatch: got %v", msg.MessageAttributes["Attr1"])
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

	results = append(results, r.RunTest("sqs", "ReceiveMessage_MaxNumberOfMessages", func() error {
		rtQueueName := fmt.Sprintf("RMNQueue-%d", time.Now().UnixNano())
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

		for i := 0; i < 5; i++ {
			client.SendMessage(ctx, &sqs.SendMessageInput{
				QueueUrl:    urlResp.QueueUrl,
				MessageBody: aws.String(fmt.Sprintf("msg-%d", i)),
			})
		}

		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            urlResp.QueueUrl,
			MaxNumberOfMessages: 5,
		})
		if err != nil {
			return fmt.Errorf("receive: %v", err)
		}
		if len(recvResp.Messages) < 5 {
			return fmt.Errorf("expected at least 5 messages, got %d", len(recvResp.Messages))
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage_WaitTimeSeconds", func() error {
		wtQueueName := fmt.Sprintf("WTQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(wtQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(wtQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(wtQueueName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String("wait-time-test-msg"),
		})
		if err != nil {
			return fmt.Errorf("send: %v", err)
		}

		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:        urlResp.QueueUrl,
			WaitTimeSeconds: 1,
		})
		if err != nil {
			return err
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("expected at least 1 message with WaitTimeSeconds=1")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage_VisibilityTimeout", func() error {
		rtQueueName := fmt.Sprintf("RVTQueue-%d", time.Now().UnixNano())
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

		client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String("visibility-test-msg"),
		})

		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:          urlResp.QueueUrl,
			VisibilityTimeout: 120,
		})
		if err != nil {
			return fmt.Errorf("receive: %v", err)
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received")
		}

		recvResp2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: urlResp.QueueUrl,
		})
		if err != nil {
			return fmt.Errorf("second receive: %v", err)
		}
		if len(recvResp2.Messages) > 0 {
			return fmt.Errorf("message should be invisible after 120s visibility timeout, but got %d messages", len(recvResp2.Messages))
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DeleteMessage", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		sendResp, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    resp.QueueUrl,
			MessageBody: aws.String("Message to delete"),
		})
		if err != nil {
			return fmt.Errorf("send: %v", err)
		}
		if sendResp.MessageId == nil {
			return fmt.Errorf("SendMessage returned nil MessageId")
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return fmt.Errorf("receive: %v", err)
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received for DeleteMessage test")
		}
		receiptHandle := recvResp.Messages[0].ReceiptHandle
		_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      resp.QueueUrl,
			ReceiptHandle: receiptHandle,
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		recvResp2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:        resp.QueueUrl,
			WaitTimeSeconds: 1,
		})
		if err != nil {
			return fmt.Errorf("receive after delete: %v", err)
		}
		if len(recvResp2.Messages) > 0 {
			return fmt.Errorf("expected 0 messages after DeleteMessage, got %d", len(recvResp2.Messages))
		}
		return nil
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
			VisibilityTimeout: 0,
		})
		if err != nil {
			return err
		}
		recvResp2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return fmt.Errorf("receive after visibility reset: %v", err)
		}
		if len(recvResp2.Messages) == 0 {
			return fmt.Errorf("expected message to be visible after setting VisibilityTimeout=0")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ChangeMessageVisibilityBatch", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		for i := 0; i < 3; i++ {
			client.SendMessage(ctx, &sqs.SendMessageInput{
				QueueUrl:    resp.QueueUrl,
				MessageBody: aws.String(fmt.Sprintf("CMVB-msg-%d", i)),
			})
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            resp.QueueUrl,
			MaxNumberOfMessages: 10,
		})
		if err != nil {
			return fmt.Errorf("receive: %v", err)
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received for ChangeMessageVisibilityBatch")
		}
		var entries []types.ChangeMessageVisibilityBatchRequestEntry
		for i, msg := range recvResp.Messages {
			entries = append(entries, types.ChangeMessageVisibilityBatchRequestEntry{
				Id:                aws.String(fmt.Sprintf("cmvb%d", i)),
				ReceiptHandle:     msg.ReceiptHandle,
				VisibilityTimeout: 120,
			})
		}
		batchResp, err := client.ChangeMessageVisibilityBatch(ctx, &sqs.ChangeMessageVisibilityBatchInput{
			QueueUrl: resp.QueueUrl,
			Entries:  entries,
		})
		if err != nil {
			return fmt.Errorf("change visibility batch: %v", err)
		}
		if len(batchResp.Successful) != len(entries) {
			return fmt.Errorf("expected %d Successful entries, got %d", len(entries), len(batchResp.Successful))
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ChangeMessageVisibilityBatch_NonExistent", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		batchResp, err := client.ChangeMessageVisibilityBatch(ctx, &sqs.ChangeMessageVisibilityBatchInput{
			QueueUrl: resp.QueueUrl,
			Entries: []types.ChangeMessageVisibilityBatchRequestEntry{
				{
					Id:                aws.String("cmvb-fail"),
					ReceiptHandle:     aws.String("nonexistent-receipt-handle"),
					VisibilityTimeout: 30,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("batch call itself failed: %v", err)
		}
		if len(batchResp.Failed) == 0 {
			return fmt.Errorf("expected Failed entry for non-existent receipt handle")
		}
		if batchResp.Failed[0].Code == nil || !strings.Contains(*batchResp.Failed[0].Code, "ReceiptHandleIsInvalid") {
			if batchResp.Failed[0].Message == nil || !strings.Contains(*batchResp.Failed[0].Message, "ReceiptHandleIsInvalid") {
				return fmt.Errorf("expected ReceiptHandleIsInvalid in failed entry, got Code=%v Message=%v", batchResp.Failed[0].Code, batchResp.Failed[0].Message)
			}
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
		if len(batchResp.Successful) != 2 {
			return fmt.Errorf("expected 2 Successful entries, got %d", len(batchResp.Successful))
		}
		for _, entry := range batchResp.Successful {
			if entry.MessageId == nil || *entry.MessageId == "" {
				return fmt.Errorf("SendMessageBatch entry has nil or empty MessageId")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessageBatch_WithDelaySeconds", func() error {
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
					Id:           aws.String("delayed1"),
					MessageBody:  aws.String("Delayed batch 1"),
					DelaySeconds: 3,
				},
				{
					Id:           aws.String("delayed2"),
					MessageBody:  aws.String("Delayed batch 2"),
					DelaySeconds: 3,
				},
			},
		})
		if err != nil {
			return err
		}
		if len(batchResp.Successful) != 2 {
			return fmt.Errorf("expected 2 successful entries, got %d", len(batchResp.Successful))
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
		if err != nil {
			return err
		}
		recvResp2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            resp.QueueUrl,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     1,
		})
		if err != nil {
			return fmt.Errorf("receive after batch delete: %v", err)
		}
		if len(recvResp2.Messages) > 0 {
			return fmt.Errorf("expected 0 messages after DeleteMessageBatch, got %d", len(recvResp2.Messages))
		}
		client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: resp.QueueUrl})
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

	return results
}
