package testutil

import (
	"context"
	"fmt"
	"strings"
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
		if attrResp.Attributes[string(types.QueueAttributeNameVisibilityTimeout)] == "" {
			return fmt.Errorf("VisibilityTimeout is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueAttributes_SpecificAttributes", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: resp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{
				types.QueueAttributeNameQueueArn,
				types.QueueAttributeNameVisibilityTimeout,
			},
		})
		if err != nil {
			return err
		}
		if attrResp.Attributes == nil {
			return fmt.Errorf("GetQueueAttributes returned nil Attributes")
		}
		if _, ok := attrResp.Attributes[string(types.QueueAttributeNameQueueArn)]; !ok {
			return fmt.Errorf("GetQueueAttributes missing QueueArn")
		}
		if _, ok := attrResp.Attributes[string(types.QueueAttributeNameVisibilityTimeout)]; !ok {
			return fmt.Errorf("GetQueueAttributes missing VisibilityTimeout")
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
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:        resp.QueueUrl,
			WaitTimeSeconds: 1,
		})
		if err != nil {
			return err
		}
		_ = recvResp
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
		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    resp.QueueUrl,
			MessageBody: aws.String("Message to delete"),
		})
		if err != nil {
			return fmt.Errorf("send: %v", err)
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

	results = append(results, r.RunTest("sqs", "ListQueues", func() error {
		resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{})
		if err != nil {
			return err
		}
		if resp.QueueUrls == nil {
			return fmt.Errorf("ListQueues returned nil QueueUrls")
		}
		found := false
		for _, u := range resp.QueueUrls {
			if strings.Contains(u, queueName) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created queue %s not found in ListQueues", queueName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListQueues_WithPrefix", func() error {
		prefix := fmt.Sprintf("PrefixTest-%d", time.Now().UnixNano())
		client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(prefix + "-alpha"),
		})
		client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(prefix + "-beta"),
		})
		defer func() {
			for _, suffix := range []string{"-alpha", "-beta"} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(prefix + suffix)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{
			QueueNamePrefix: aws.String(prefix),
		})
		if err != nil {
			return err
		}
		if len(resp.QueueUrls) < 2 {
			return fmt.Errorf("expected at least 2 queues with prefix %q, got %d", prefix, len(resp.QueueUrls))
		}
		for _, u := range resp.QueueUrls {
			if !strings.Contains(u, prefix) {
				return fmt.Errorf("ListQueues returned URL not matching prefix: %s", u)
			}
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

	results = append(results, r.RunTest("sqs", "SetQueueAttributes_MultipleAttrs", func() error {
		rtQueueName := fmt.Sprintf("SMAQueue-%d", time.Now().UnixNano())
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

		_, err = client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
			QueueUrl: urlResp.QueueUrl,
			Attributes: map[string]string{
				"VisibilityTimeout":      "45",
				"MaximumMessageSize":     "1024",
				"MessageRetentionPeriod": "345600",
				"DelaySeconds":           "10",
			},
		})
		if err != nil {
			return fmt.Errorf("set attrs: %v", err)
		}

		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: urlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{
				types.QueueAttributeNameVisibilityTimeout,
				types.QueueAttributeNameMaximumMessageSize,
				types.QueueAttributeNameMessageRetentionPeriod,
				types.QueueAttributeNameDelaySeconds,
			},
		})
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if attrResp.Attributes[string(types.QueueAttributeNameVisibilityTimeout)] != "45" {
			return fmt.Errorf("VisibilityTimeout mismatch: got %s", attrResp.Attributes["VisibilityTimeout"])
		}
		if attrResp.Attributes[string(types.QueueAttributeNameMaximumMessageSize)] != "1024" {
			return fmt.Errorf("MaximumMessageSize mismatch: got %s", attrResp.Attributes["MaximumMessageSize"])
		}
		if attrResp.Attributes[string(types.QueueAttributeNameMessageRetentionPeriod)] != "345600" {
			return fmt.Errorf("MessageRetentionPeriod mismatch: got %s", attrResp.Attributes["MessageRetentionPeriod"])
		}
		if attrResp.Attributes[string(types.QueueAttributeNameDelaySeconds)] != "10" {
			return fmt.Errorf("DelaySeconds mismatch: got %s", attrResp.Attributes["DelaySeconds"])
		}
		return nil
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
		if len(tagResp.Tags) == 0 {
			return fmt.Errorf("ListQueueTags returned nil or empty Tags (expected Environment=Test)")
		}
		if tagResp.Tags["Environment"] != "Test" {
			return fmt.Errorf("expected tag Environment=Test, got %v", tagResp.Tags["Environment"])
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
		if len(batchResp.Successful) == 0 {
			return fmt.Errorf("ChangeMessageVisibilityBatch returned empty Successful entries")
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
		client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: resp.QueueUrl})
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
		client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: resp.QueueUrl})
		return nil
	}))

	// === PERMISSION TESTS ===

	results = append(results, r.RunTest("sqs", "AddPermission", func() error {
		permQueueName := fmt.Sprintf("PermQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(permQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(permQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(permQueueName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		_, err = client.AddPermission(ctx, &sqs.AddPermissionInput{
			QueueUrl:      urlResp.QueueUrl,
			Label:         aws.String("TestPermission"),
			AWSAccountIds: []string{"123456789012"},
			Actions:       []string{"SendMessage", "ReceiveMessage"},
		})
		return err
	}))

	results = append(results, r.RunTest("sqs", "RemovePermission", func() error {
		permQueueName := fmt.Sprintf("RPermQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(permQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(permQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(permQueueName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		_, err = client.AddPermission(ctx, &sqs.AddPermissionInput{
			QueueUrl:      urlResp.QueueUrl,
			Label:         aws.String("RemoveTestPerm"),
			AWSAccountIds: []string{"123456789012"},
			Actions:       []string{"SendMessage"},
		})
		if err != nil {
			return fmt.Errorf("add permission: %v", err)
		}

		_, err = client.RemovePermission(ctx, &sqs.RemovePermissionInput{
			QueueUrl: urlResp.QueueUrl,
			Label:    aws.String("RemoveTestPerm"),
		})
		return err
	}))

	// === DLQ / MESSAGE MOVE TASK TESTS ===

	results = append(results, r.RunTest("sqs", "ListDeadLetterSourceQueues_Empty", func() error {
		dlqName := fmt.Sprintf("DLQ-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(dlqName),
		})
		if err != nil {
			return fmt.Errorf("create dlq: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(dlqName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(dlqName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		dlqResp, err := client.ListDeadLetterSourceQueues(ctx, &sqs.ListDeadLetterSourceQueuesInput{
			QueueUrl: urlResp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if len(dlqResp.QueueUrls) != 0 {
			return fmt.Errorf("expected empty queue URLs for new DLQ, got %d", len(dlqResp.QueueUrls))
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "StartMessageMoveTask", func() error {
		srcDlqName := fmt.Sprintf("SrcDLQ-%d", time.Now().UnixNano())
		destQueueName := fmt.Sprintf("DestQueue-%d", time.Now().UnixNano())

		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("create src: %v", err)
		}
		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("create dest: %v", err)
		}
		defer func() {
			for _, name := range []string{srcDlqName, destQueueName} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(name)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		srcUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("get src url: %v", err)
		}
		destUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("get dest url: %v", err)
		}

		srcAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       srcUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get src attrs: %v", err)
		}
		destAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       destUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get dest attrs: %v", err)
		}

		srcArn := srcAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]
		destArn := destAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]

		taskResp, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:      aws.String(srcArn),
			DestinationArn: aws.String(destArn),
		})
		if err != nil {
			return fmt.Errorf("start task: %v", err)
		}
		if taskResp.TaskHandle == nil || *taskResp.TaskHandle == "" {
			return fmt.Errorf("StartMessageMoveTask returned nil or empty TaskHandle")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "CancelMessageMoveTask", func() error {
		srcDlqName := fmt.Sprintf("CancelDLQ-%d", time.Now().UnixNano())
		destQueueName := fmt.Sprintf("CancelDest-%d", time.Now().UnixNano())

		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("create src: %v", err)
		}
		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("create dest: %v", err)
		}
		defer func() {
			for _, name := range []string{srcDlqName, destQueueName} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(name)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		srcUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("get src url: %v", err)
		}
		destUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("get dest url: %v", err)
		}

		srcAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       srcUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get src attrs: %v", err)
		}
		destAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       destUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get dest attrs: %v", err)
		}

		srcArn := srcAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]
		destArn := destAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]

		taskResp, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:      aws.String(srcArn),
			DestinationArn: aws.String(destArn),
		})
		if err != nil {
			return fmt.Errorf("start task: %v", err)
		}

		cancelResp, err := client.CancelMessageMoveTask(ctx, &sqs.CancelMessageMoveTaskInput{
			TaskHandle: taskResp.TaskHandle,
		})
		if err != nil {
			return fmt.Errorf("cancel task: %v", err)
		}
		_ = cancelResp.ApproximateNumberOfMessagesMoved
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListMessageMoveTasks", func() error {
		srcDlqName := fmt.Sprintf("ListDLQ-%d", time.Now().UnixNano())
		destQueueName := fmt.Sprintf("ListDest-%d", time.Now().UnixNano())

		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("create src: %v", err)
		}
		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("create dest: %v", err)
		}
		defer func() {
			for _, name := range []string{srcDlqName, destQueueName} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(name)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		srcUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("get src url: %v", err)
		}
		destUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("get dest url: %v", err)
		}

		srcAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       srcUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get src attrs: %v", err)
		}
		destAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       destUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get dest attrs: %v", err)
		}

		srcArn := srcAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]
		destArn := destAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]

		_, err = client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:      aws.String(srcArn),
			DestinationArn: aws.String(destArn),
		})
		if err != nil {
			return fmt.Errorf("start task: %v", err)
		}

		listResp, err := client.ListMessageMoveTasks(ctx, &sqs.ListMessageMoveTasksInput{
			SourceArn: aws.String(srcArn),
		})
		if err != nil {
			return fmt.Errorf("list tasks: %v", err)
		}
		if len(listResp.Results) == 0 {
			return fmt.Errorf("ListMessageMoveTasks returned empty Results")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("sqs", "GetQueueUrl_NonExistent", func() error {
		_, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String("nonexistent-queue-xyz"),
		})
		if err := AssertErrorContains(err, "QueueDoesNotExist"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "ReceiptHandleIsInvalid"); err != nil {
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

	results = append(results, r.RunTest("sqs", "ListQueues_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgQueues []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagQ-%s-%d", pgTs, i)
			_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
				QueueName: aws.String(name),
			})
			if err != nil {
				for _, qn := range pgQueues {
					uResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qn)})
					if uResp.QueueUrl != nil {
						client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: uResp.QueueUrl})
					}
				}
				return fmt.Errorf("create queue %s: %v", name, err)
			}
			pgQueues = append(pgQueues, name)
		}

		var allQueues []string
		var nextToken *string
		for {
			resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, qn := range pgQueues {
					uResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qn)})
					if uResp.QueueUrl != nil {
						client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: uResp.QueueUrl})
					}
				}
				return fmt.Errorf("list queues page: %v", err)
			}
			for _, u := range resp.QueueUrls {
				if strings.Contains(u, "PagQ-"+pgTs) {
					allQueues = append(allQueues, u)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, qn := range pgQueues {
			uResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qn)})
			if uResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: uResp.QueueUrl})
			}
		}
		if len(allQueues) != 5 {
			return fmt.Errorf("expected 5 paginated queues, got %d", len(allQueues))
		}
		return nil
	}))

	return results
}
