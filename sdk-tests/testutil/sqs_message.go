package testutil

import (
	"context"
	"errors"
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

	results = append(results, r.RunTest("sqs", "SendMessage_AttributeListValuesRejected", func() error {
		listAttrQueueName := fmt.Sprintf("ListAttrQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(listAttrQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(listAttrQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(listAttrQueueName),
		})
		if err != nil {
			return err
		}
		// The service model marks StringListValues and BinaryListValues as
		// not implemented; AWS rejects requests carrying them with
		// UnsupportedOperation.
		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    resp.QueueUrl,
			MessageBody: aws.String("Message with list attribute values"),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"Attr1": {
					DataType:         aws.String("String"),
					StringListValues: []string{"value1"},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("SendMessage with StringListValues unexpectedly succeeded")
		}
		var opErr *types.UnsupportedOperation
		if !errors.As(err, &opErr) {
			return fmt.Errorf("expected UnsupportedOperation, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessageBatch_AttributeListValuesRejected", func() error {
		batchAttrQueueName := fmt.Sprintf("BatchAttrQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(batchAttrQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(batchAttrQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(batchAttrQueueName),
		})
		if err != nil {
			return err
		}
		_, err = client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			QueueUrl: resp.QueueUrl,
			Entries: []types.SendMessageBatchRequestEntry{
				{
					Id:          aws.String("entry-1"),
					MessageBody: aws.String("body"),
					MessageAttributes: map[string]types.MessageAttributeValue{
						"Attr1": {
							DataType:         aws.String("String"),
							StringListValues: []string{"value1"},
						},
					},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("SendMessageBatch with StringListValues unexpectedly succeeded")
		}
		var opErr *types.UnsupportedOperation
		if !errors.As(err, &opErr) {
			return fmt.Errorf("expected UnsupportedOperation, got %v", err)
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

	results = append(results, r.RunTest("sqs", "SendMessage_NoAttributes_OmitsAttributesMD5", func() error {
		// The MD5OfMessageAttributes member is only present when the message
		// carries attributes; an always-present md5("") would diverge from
		// the AWS response surface.
		md5QueueName := fmt.Sprintf("Md5Queue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(md5QueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(md5QueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(md5QueueName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		sendResp, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String("no attributes message"),
		})
		if err != nil {
			return err
		}
		if sendResp.MD5OfMessageAttributes != nil {
			return fmt.Errorf("expected MD5OfMessageAttributes to be omitted for a message without attributes, got %q", *sendResp.MD5OfMessageAttributes)
		}

		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: urlResp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if len(recvResp.Messages) != 1 {
			return fmt.Errorf("expected 1 received message, got %d", len(recvResp.Messages))
		}
		if recvResp.Messages[0].MD5OfMessageAttributes != nil {
			return fmt.Errorf("expected received MD5OfMessageAttributes to be omitted, got %q", *recvResp.Messages[0].MD5OfMessageAttributes)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage_FifoSequenceNumber", func() error {
		// FIFO messages expose SequenceNumber as a system attribute on
		// receive when all system attributes are requested.
		fifoName := fmt.Sprintf("FifoSeq-%d.fifo", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(fifoName),
			Attributes: map[string]string{
				"FifoQueue":                 "true",
				"ContentBasedDeduplication": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create fifo: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(fifoName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(fifoName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:               urlResp.QueueUrl,
			MessageBody:            aws.String("fifo sequence number test"),
			MessageGroupId:         aws.String("group-1"),
			MessageDeduplicationId: aws.String("dedup-seq-1"),
		})
		if err != nil {
			return fmt.Errorf("send: %v", err)
		}

		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:                    urlResp.QueueUrl,
			MessageSystemAttributeNames: []types.MessageSystemAttributeName{types.MessageSystemAttributeNameAll},
		})
		if err != nil {
			return err
		}
		if len(recvResp.Messages) != 1 {
			return fmt.Errorf("expected 1 received message, got %d", len(recvResp.Messages))
		}
		seq := recvResp.Messages[0].Attributes["SequenceNumber"]
		if seq == "" {
			return fmt.Errorf("expected non-empty SequenceNumber attribute on FIFO receive, attributes: %v", recvResp.Messages[0].Attributes)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessage_InvalidBodyCharset_Rejected", func() error {
		// A body with characters outside the allowed set must be rejected
		// with InvalidMessageContents.
		csQueueName := fmt.Sprintf("CharsetQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(csQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(csQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(csQueueName)})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String("bad\x00body"),
		})
		if err == nil {
			return fmt.Errorf("body containing NUL must be rejected with InvalidMessageContents")
		}
		if !strings.Contains(err.Error(), "InvalidMessageContents") {
			return fmt.Errorf("expected InvalidMessageContents, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SendMessage_FifoIdOverLength_Rejected", func() error {
		idQName := fmt.Sprintf("IdLenFifo-%d.fifo", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(idQName),
			Attributes: map[string]string{
				"FifoQueue":                 "true",
				"ContentBasedDeduplication": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(idQName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(idQName)})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		longID := strings.Repeat("a", 129)
		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:               urlResp.QueueUrl,
			MessageBody:            aws.String("over-length group id"),
			MessageGroupId:         aws.String(longID),
			MessageDeduplicationId: aws.String("dedup-ok-1"),
		})
		if err == nil {
			return fmt.Errorf("129-character MessageGroupId must be rejected")
		}

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:               urlResp.QueueUrl,
			MessageBody:            aws.String("invalid charset dedup id"),
			MessageGroupId:         aws.String("group-ok"),
			MessageDeduplicationId: aws.String("has space"),
		})
		if err == nil {
			return fmt.Errorf("MessageDeduplicationId with a space must be rejected")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage_LongPoll_EmptyQueue_Waits", func() error {
		lpQueueName := fmt.Sprintf("LongPoll-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(lpQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(lpQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(lpQueueName)})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		start := time.Now()
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:        urlResp.QueueUrl,
			WaitTimeSeconds: 1,
		})
		if err != nil {
			return err
		}
		elapsed := time.Since(start)
		if len(recvResp.Messages) != 0 {
			return fmt.Errorf("expected 0 messages on an empty queue, got %d", len(recvResp.Messages))
		}
		if elapsed < 800*time.Millisecond {
			return fmt.Errorf("expected the call to wait close to 1s, returned after %v", elapsed)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ReceiveMessage_LongPoll_QueueDefault", func() error {
		// An omitted WaitTimeSeconds must fall back to the queue's
		// ReceiveMessageWaitTimeSeconds attribute.
		qdQueueName := fmt.Sprintf("QueueDefaultPoll-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName:  aws.String(qdQueueName),
			Attributes: map[string]string{"ReceiveMessageWaitTimeSeconds": "1"},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qdQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qdQueueName)})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		start := time.Now()
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: urlResp.QueueUrl,
		})
		if err != nil {
			return err
		}
		elapsed := time.Since(start)
		if len(recvResp.Messages) != 0 {
			return fmt.Errorf("expected 0 messages on an empty queue, got %d", len(recvResp.Messages))
		}
		if elapsed < 800*time.Millisecond {
			return fmt.Errorf("expected the queue default wait of 1s, returned after %v", elapsed)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "PurgeQueue_SecondWithin60s_Rejected", func() error {
		purgeQName := fmt.Sprintf("PurgeWin-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(purgeQName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(purgeQName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(purgeQName)})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		if _, err = client.PurgeQueue(ctx, &sqs.PurgeQueueInput{QueueUrl: urlResp.QueueUrl}); err != nil {
			return fmt.Errorf("first purge: %v", err)
		}
		_, err = client.PurgeQueue(ctx, &sqs.PurgeQueueInput{QueueUrl: urlResp.QueueUrl})
		if err == nil {
			return fmt.Errorf("second purge within 60 seconds must fail with PurgeQueueInProgress")
		}
		if !strings.Contains(err.Error(), "PurgeQueueInProgress") {
			return fmt.Errorf("expected PurgeQueueInProgress, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DeleteMessage_OldReceiptHandle_Succeeds", func() error {
		// "If you use an old ReceiptHandle, the request will succeed, but the
		// message might not be deleted."
		oldHQueueName := fmt.Sprintf("OldHandle-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName:  aws.String(oldHQueueName),
			Attributes: map[string]string{"VisibilityTimeout": "0"},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(oldHQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(oldHQueueName)})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		if _, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String("old handle test"),
		}); err != nil {
			return fmt.Errorf("send: %v", err)
		}

		recv1, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: urlResp.QueueUrl})
		if err != nil || len(recv1.Messages) != 1 {
			return fmt.Errorf("receive 1: %v messages=%d", err, len(recv1.Messages))
		}
		firstHandle := *recv1.Messages[0].ReceiptHandle

		time.Sleep(100 * time.Millisecond)
		recv2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: urlResp.QueueUrl})
		if err != nil || len(recv2.Messages) != 1 {
			return fmt.Errorf("receive 2: %v messages=%d", err, len(recv2.Messages))
		}

		if _, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      urlResp.QueueUrl,
			ReceiptHandle: aws.String(firstHandle),
		}); err != nil {
			return fmt.Errorf("delete with the older handle must succeed, got: %v", err)
		}

		recv3, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: urlResp.QueueUrl})
		if err != nil {
			return fmt.Errorf("receive 3: %v", err)
		}
		if len(recv3.Messages) != 0 {
			return fmt.Errorf("message still deliverable after delete with old handle")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DeleteMessage_CrossQueueHandle_Rejected", func() error {
		crossA := fmt.Sprintf("CrossA-%d", time.Now().UnixNano())
		crossB := fmt.Sprintf("CrossB-%d", time.Now().UnixNano())
		for _, name := range []string{crossA, crossB} {
			if _, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(name)}); err != nil {
				return fmt.Errorf("create %s: %v", name, err)
			}
		}
		defer func() {
			for _, name := range []string{crossA, crossB} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(name)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()
		urlA, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(crossA)})
		if err != nil {
			return fmt.Errorf("get url A: %v", err)
		}
		urlB, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(crossB)})
		if err != nil {
			return fmt.Errorf("get url B: %v", err)
		}

		if _, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlA.QueueUrl,
			MessageBody: aws.String("cross-queue handle test"),
		}); err != nil {
			return fmt.Errorf("send: %v", err)
		}
		recvA, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: urlA.QueueUrl})
		if err != nil || len(recvA.Messages) != 1 {
			return fmt.Errorf("receive A: %v messages=%d", err, len(recvA.Messages))
		}
		handleA := *recvA.Messages[0].ReceiptHandle

		_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      urlB.QueueUrl,
			ReceiptHandle: aws.String(handleA),
		})
		if err == nil {
			// The message must still be in queue A.
			return fmt.Errorf("a receipt handle from queue A must not delete via queue B")
		}
		if !strings.Contains(err.Error(), "ReceiptHandleIsInvalid") {
			return fmt.Errorf("expected ReceiptHandleIsInvalid, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ChangeMessageVisibility_NotInflight_Rejected", func() error {
		nifQueueName := fmt.Sprintf("NotInflight-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName:  aws.String(nifQueueName),
			Attributes: map[string]string{"VisibilityTimeout": "0"},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(nifQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(nifQueueName)})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		if _, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String("not in flight test"),
		}); err != nil {
			return fmt.Errorf("send: %v", err)
		}
		recv, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: urlResp.QueueUrl})
		if err != nil || len(recv.Messages) != 1 {
			return fmt.Errorf("receive: %v messages=%d", err, len(recv.Messages))
		}

		// With VisibilityTimeout=0 the message is immediately no longer in
		// flight, so extending its visibility must fail.
		_, err = client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
			QueueUrl:          urlResp.QueueUrl,
			ReceiptHandle:     recv.Messages[0].ReceiptHandle,
			VisibilityTimeout: 30,
		})
		if err == nil {
			return fmt.Errorf("expected MessageNotInflight, got success")
		}
		if !strings.Contains(err.Error(), "MessageNotInflight") {
			return fmt.Errorf("expected MessageNotInflight, got: %v", err)
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
		// A dedicated queue keeps the follow-up receive free of messages
		// other suites leave on the shared runner queue (including delayed
		// sends that become visible during the wait).
		delQueueName := fmt.Sprintf("DelMsg-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(delQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(delQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(delQueueName),
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
		// The deletion is committed server-side before DeleteMessage returns,
		// so an immediate receive proves the message is gone without a wait.
		recvResp2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: resp.QueueUrl,
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

	results = append(results, r.RunTest("sqs", "SendMessageBatch_MissingQueue_RequestLevelError", func() error {
		// QueueDoesNotExist is a request-level error on batch operations:
		// the whole request fails instead of returning per-entry failures.
		_, err := client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			QueueUrl: aws.String("http://localhost:50080/queue/does-not-exist-at-all"),
			Entries: []types.SendMessageBatchRequestEntry{
				{Id: aws.String("a"), MessageBody: aws.String("x")},
			},
		})
		if err == nil {
			return fmt.Errorf("expected QueueDoesNotExist, got success")
		}
		if !strings.Contains(err.Error(), "QueueDoesNotExist") {
			return fmt.Errorf("expected QueueDoesNotExist, got: %v", err)
		}
		return nil
	}))

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
		// The batch deletion is committed server-side before the response,
		// so an immediate receive proves the queue is empty without a wait.
		recvResp2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            resp.QueueUrl,
			MaxNumberOfMessages: 10,
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

	results = append(results, r.RunTest("sqs", "ReceiveMessage_AttributeFiltering", func() error {
		filterQueueName := fmt.Sprintf("FilterQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(filterQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(filterQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(filterQueueName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    urlResp.QueueUrl,
			MessageBody: aws.String("attribute-filter-test"),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"Alpha": {
					DataType:    aws.String("String"),
					StringValue: aws.String("aaa"),
				},
				"Beta": {
					DataType:    aws.String("Number"),
					StringValue: aws.String("99"),
				},
				"Gamma": {
					DataType:    aws.String("String"),
					StringValue: aws.String("ggg"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("send: %v", err)
		}

		// Receive only "Alpha" message attribute and only "SenderId" system attribute
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:              urlResp.QueueUrl,
			MessageAttributeNames: []string{"Alpha"},
			MessageSystemAttributeNames: []types.MessageSystemAttributeName{
				types.MessageSystemAttributeNameSenderId,
			},
		})
		if err != nil {
			return fmt.Errorf("receive: %v", err)
		}
		if len(recvResp.Messages) == 0 {
			return fmt.Errorf("no messages received")
		}
		msg := recvResp.Messages[0]

		// Should have only Alpha in MessageAttributes, not Beta or Gamma
		if len(msg.MessageAttributes) != 1 {
			return fmt.Errorf("expected 1 message attribute, got %d: %+v", len(msg.MessageAttributes), msg.MessageAttributes)
		}
		if _, ok := msg.MessageAttributes["Alpha"]; !ok {
			return fmt.Errorf("expected Alpha in MessageAttributes, got keys: %+v", msg.MessageAttributes)
		}
		if _, ok := msg.MessageAttributes["Beta"]; ok {
			return fmt.Errorf("Beta should NOT be in MessageAttributes (not requested)")
		}
		if _, ok := msg.MessageAttributes["Gamma"]; ok {
			return fmt.Errorf("Gamma should NOT be in MessageAttributes (not requested)")
		}

		// Should have only SenderId in Attributes
		if len(msg.Attributes) != 1 {
			return fmt.Errorf("expected 1 system attribute, got %d: %+v", len(msg.Attributes), msg.Attributes)
		}
		if _, ok := msg.Attributes["SenderId"]; !ok {
			return fmt.Errorf("expected SenderId in Attributes, got keys: %+v", msg.Attributes)
		}
		if _, ok := msg.Attributes["SentTimestamp"]; ok {
			return fmt.Errorf("SentTimestamp should NOT be in Attributes (not requested)")
		}

		return nil
	}))

	return results
}
