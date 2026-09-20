package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func (r *TestRunner) runSNSPublishTests(tc *snsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sns", "Publish", func() error {
		topicName := tc.uniqueName("PubTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		resp, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(topicArn),
			Message:  aws.String("Test message"),
		})
		if err != nil {
			return err
		}
		if resp.MessageId == nil || *resp.MessageId == "" {
			return fmt.Errorf("MessageId should be non-empty")
		}
		return nil
	}))

	// Subjects are documented as "UTF-8 text … less than 100 characters
	// long", so a 99-character CJK subject (297 bytes) is legal input.
	results = append(results, r.RunTest("sns", "Publish_SubjectMultibyteAccepted", func() error {
		topicName := tc.uniqueName("PubMbSubjectTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		subject := strings.Repeat("\u65e5", 99)
		resp, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(topicArn),
			Message:  aws.String("body"),
			Subject:  aws.String(subject),
		})
		if err != nil {
			return fmt.Errorf("Publish with 99-character CJK subject: %v", err)
		}
		if resp.MessageId == nil || *resp.MessageId == "" {
			return fmt.Errorf("MessageId should be non-empty")
		}
		return nil
	}))

	// The same documentation bounds subjects to strictly less than 100
	// characters, so a 100-character subject must be rejected even though
	// its byte length is far below any byte-based ceiling.
	results = append(results, r.RunTest("sns", "Publish_SubjectExactLimitRejected", func() error {
		topicName := tc.uniqueName("PubMbSubjLimTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		_, err = tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(topicArn),
			Message:  aws.String("body"),
			Subject:  aws.String(strings.Repeat("\u65e5", 100)),
		})
		if err == nil {
			return fmt.Errorf("expected rejection for 100-character CJK subject")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_FIFO_WithMessageGroupId", func() error {
		fifoName := tc.uniqueName("FifoPubTopic") + ".fifo"
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoName),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		fifoTopicArn := *resp.TopicArn
		defer tc.deleteTopic(fifoTopicArn)

		pubResp, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn:               aws.String(fifoTopicArn),
			Message:                aws.String("FIFO message 1"),
			MessageGroupId:         aws.String("group1"),
			MessageDeduplicationId: aws.String("dedup-fifo-1"),
		})
		if err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		if pubResp.MessageId == nil || *pubResp.MessageId == "" {
			return fmt.Errorf("MessageId is nil or empty")
		}
		if pubResp.SequenceNumber == nil || *pubResp.SequenceNumber == "" {
			return fmt.Errorf("SequenceNumber should be present for FIFO topic")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_FIFO_ContentBasedDedup", func() error {
		fifoDedupName := tc.uniqueName("FifoDedupTopic") + ".fifo"
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoDedupName),
			Attributes: map[string]string{
				"FifoTopic":                 "true",
				"ContentBasedDeduplication": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*resp.TopicArn)

		msg := "dedup test message content"
		pub1, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn:       resp.TopicArn,
			Message:        aws.String(msg),
			MessageGroupId: aws.String("dedup-group"),
		})
		if err != nil {
			return fmt.Errorf("publish 1: %v", err)
		}

		pub2, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn:       resp.TopicArn,
			Message:        aws.String(msg),
			MessageGroupId: aws.String("dedup-group"),
		})
		if err != nil {
			return fmt.Errorf("publish 2: %v", err)
		}

		if *pub1.MessageId != *pub2.MessageId {
			return fmt.Errorf("content-based dedup should return same MessageId: %q vs %q", *pub1.MessageId, *pub2.MessageId)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_FIFO_DeduplicationId", func() error {
		fifoDedupIdName := tc.uniqueName("FifoDedupIdTopic") + ".fifo"
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoDedupIdName),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*resp.TopicArn)

		pub1, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn:               resp.TopicArn,
			Message:                aws.String("msg A"),
			MessageGroupId:         aws.String("group-dedup-id"),
			MessageDeduplicationId: aws.String("dedup-123"),
		})
		if err != nil {
			return fmt.Errorf("publish 1: %v", err)
		}

		pub2, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn:               resp.TopicArn,
			Message:                aws.String("msg B"),
			MessageGroupId:         aws.String("group-dedup-id"),
			MessageDeduplicationId: aws.String("dedup-123"),
		})
		if err != nil {
			return fmt.Errorf("publish 2: %v", err)
		}

		if *pub1.MessageId != *pub2.MessageId {
			return fmt.Errorf("explicit dedup ID should return same MessageId: %q vs %q", *pub1.MessageId, *pub2.MessageId)
		}
		// A dedup hit answers with the ORIGINAL publish's identifiers: the
		// sequence number is assigned to each message and the duplicate
		// publish IS that message — accepted, not delivered, and no new
		// number allocated for it.
		if pub1.SequenceNumber == nil || *pub1.SequenceNumber == "" {
			return fmt.Errorf("FIFO publish response carries no SequenceNumber")
		}
		if pub2.SequenceNumber == nil || *pub2.SequenceNumber != *pub1.SequenceNumber {
			return fmt.Errorf("dedup hit SequenceNumber = %v, want the original %q", pub2.SequenceNumber, *pub1.SequenceNumber)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_WithMessageAttributes", func() error {
		attrTopicArn, err := tc.createTopic(tc.uniqueName("AttrPubTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(attrTopicArn)

		pubResp, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(attrTopicArn),
			Message:  aws.String("message with attrs"),
			Subject:  aws.String("Test Subject"),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"Attr1": {DataType: aws.String("String"), StringValue: aws.String("value1")},
				"Attr2": {DataType: aws.String("Number"), StringValue: aws.String("42")},
			},
		})
		if err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		if pubResp.MessageId == nil || *pubResp.MessageId == "" {
			return fmt.Errorf("MessageId should be non-empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "PublishBatch", func() error {
		topicName := tc.uniqueName("BatchTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		resp, err := tc.client.PublishBatch(tc.ctx, &sns.PublishBatchInput{
			TopicArn: aws.String(topicArn),
			PublishBatchRequestEntries: []types.PublishBatchRequestEntry{
				{Id: aws.String("msg1"), Message: aws.String("Batch message 1")},
				{Id: aws.String("msg2"), Message: aws.String("Batch message 2")},
			},
		})
		if err != nil {
			return err
		}
		if len(resp.Successful) != 2 {
			return fmt.Errorf("expected 2 successful, got %d", len(resp.Successful))
		}
		for _, s := range resp.Successful {
			if s.MessageId == nil || *s.MessageId == "" {
				return fmt.Errorf("successful entry should have MessageId")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "PublishBatch_WithAttributes", func() error {
		batchAttrTopicArn, err := tc.createTopic(tc.uniqueName("BatchAttrTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(batchAttrTopicArn)

		resp, err := tc.client.PublishBatch(tc.ctx, &sns.PublishBatchInput{
			TopicArn: aws.String(batchAttrTopicArn),
			PublishBatchRequestEntries: []types.PublishBatchRequestEntry{
				{
					Id:      aws.String("b1"),
					Message: aws.String("batch attr msg 1"),
					MessageAttributes: map[string]types.MessageAttributeValue{
						"BAttr": {DataType: aws.String("String"), StringValue: aws.String("bval1")},
					},
				},
				{
					Id:      aws.String("b2"),
					Message: aws.String("batch attr msg 2"),
					MessageAttributes: map[string]types.MessageAttributeValue{
						"BAttr": {DataType: aws.String("Number"), StringValue: aws.String("99")},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("publish batch: %v", err)
		}
		if len(resp.Successful) != 2 {
			return fmt.Errorf("expected 2 successful, got %d", len(resp.Successful))
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "PublishBatch_MaxEntries", func() error {
		maxBatchTopicArn, err := tc.createTopic(tc.uniqueName("MaxBatchTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(maxBatchTopicArn)

		entries := make([]types.PublishBatchRequestEntry, 10)
		for i := 0; i < 10; i++ {
			entries[i] = types.PublishBatchRequestEntry{
				Id:      aws.String(fmt.Sprintf("max%d", i)),
				Message: aws.String(fmt.Sprintf("max batch message %d", i)),
			}
		}
		resp, err := tc.client.PublishBatch(tc.ctx, &sns.PublishBatchInput{
			TopicArn:                   aws.String(maxBatchTopicArn),
			PublishBatchRequestEntries: entries,
		})
		if err != nil {
			return fmt.Errorf("publish batch: %v", err)
		}
		if len(resp.Successful) != 10 {
			return fmt.Errorf("expected 10 successful, got %d", len(resp.Successful))
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "PublishBatch_FailedEntry", func() error {
		failBatchFifoName := tc.uniqueName("FailBatchFifo") + ".fifo"
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(failBatchFifoName),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		resp, err := tc.client.PublishBatch(tc.ctx, &sns.PublishBatchInput{
			TopicArn: tResp.TopicArn,
			PublishBatchRequestEntries: []types.PublishBatchRequestEntry{
				{
					Id:                     aws.String("good1"),
					Message:                aws.String("valid message"),
					MessageGroupId:         aws.String("g1"),
					MessageDeduplicationId: aws.String("d1"),
				},
				{
					Id:      aws.String("bad1"),
					Message: aws.String("missing group id"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("publish batch: %v", err)
		}
		if len(resp.Successful) != 1 {
			return fmt.Errorf("expected 1 successful, got %d", len(resp.Successful))
		}
		if len(resp.Failed) != 1 {
			return fmt.Errorf("expected 1 failed, got %d", len(resp.Failed))
		}
		if resp.Failed[0].Id == nil || *resp.Failed[0].Id != "bad1" {
			return fmt.Errorf("failed entry should be bad1")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_FIFO_MissingMessageGroupId", func() error {
		fifoName := tc.uniqueName("FifoNoGroupTopic") + ".fifo"
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoName),
			Attributes: map[string]string{
				"FifoTopic":                 "true",
				"ContentBasedDeduplication": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		_, err = tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: tResp.TopicArn,
			Message:  aws.String("missing group id"),
		})
		if err == nil {
			return fmt.Errorf("expected error for FIFO publish without MessageGroupId")
		}
		return nil
	}))

	// A topic with more subscriptions than one list page (a page holds at
	// most 100) delivers to every subscription: the fan-out walks all
	// pages, so subscription 101 is not silently dropped. Each subscription
	// of the same queue receives its own copy, so counting the delivered
	// bodies counts the reached subscriptions.
	results = append(results, r.RunTest("sns", "Publish_FanOutAllSubscriptionPages", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("FanOutTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		// One distinct queue per subscription: Subscribe is idempotent on
		// the subscription natural key (topic, protocol, endpoint, owner),
		// so a shared queue ARN would collapse the 101 subscriptions into
		// a single record before the fan-out ever ran.
		const want = 101
		type subscriber struct {
			url     *string
			cleanup func()
		}
		subscribers := make([]subscriber, 0, want)
		defer func() {
			for _, s := range subscribers {
				s.cleanup()
			}
		}()
		subscriptionArns := make([]string, 0, want)
		defer func() {
			for _, arn := range subscriptionArns {
				tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: aws.String(arn)})
			}
		}()
		for i := 0; i < want; i++ {
			qURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName(fmt.Sprintf("FanOutQueue%03d", i)), nil)
			if err != nil {
				return fmt.Errorf("create delivery queue %d: %w", i, err)
			}
			subscribers = append(subscribers, subscriber{url: qURL, cleanup: cleanupQueue})
			qARN, err := queueArn(tc.ctx, sqsClient, qURL)
			if err != nil {
				return fmt.Errorf("queue ARN %d: %w", i, err)
			}
			resp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
				TopicArn: aws.String(topicArn),
				Protocol: aws.String("sqs"),
				Endpoint: aws.String(qARN),
			})
			if err != nil {
				return fmt.Errorf("subscribe %d: %w", i, err)
			}
			subscriptionArns = append(subscriptionArns, aws.ToString(resp.SubscriptionArn))
		}

		pubResp, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(topicArn),
			Message:  aws.String("fan-out walks every page"),
		})
		if err != nil {
			return err
		}
		messageID := aws.ToString(pubResp.MessageId)

		// Delivery through the event bus is asynchronous: poll until every
		// subscription's copy arrives in its own queue or the deadline
		// passes.
		deadline := time.Now().Add(60 * time.Second)
		received, matched := 0, 0
		for received < want && time.Now().Before(deadline) {
			for _, s := range subscribers {
				if received == want {
					break
				}
				msgs, err := sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
					QueueUrl:            s.url,
					MaxNumberOfMessages: 10,
					WaitTimeSeconds:     0,
				})
				if err != nil {
					return fmt.Errorf("receive: %w", err)
				}
				for _, m := range msgs.Messages {
					received++
					if strings.Contains(aws.ToString(m.Body), messageID) {
						matched++
					}
				}
			}
		}
		if received != want {
			return fmt.Errorf("fan-out delivered %d of %d subscriptions — a page of the subscription walk went unread", received, want)
		}
		if matched != want {
			return fmt.Errorf("only %d of %d delivered bodies carry the published MessageId", matched, want)
		}
		return nil
	}))

	// FIFO topics deliver same-group messages in sequence order: five
	// publishes to one message group arrive at the subscribed SQS FIFO
	// queue in publish order — the queue receives them in the order SNS
	// dispatched, which the per-group critical section serialises against
	// concurrent publishers.
	results = append(results, r.RunTest("sns", "Publish_FIFO_OrderedDelivery", func() error {
		topicArn, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name:       aws.String(tc.uniqueName("FifoOrderTopic") + ".fifo"),
			Attributes: map[string]string{"FifoTopic": "true"},
		})
		if err != nil {
			return fmt.Errorf("create FIFO topic: %v", err)
		}
		defer tc.deleteTopic(*topicArn.TopicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		queueURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("FifoOrderQueue")+".fifo", map[string]string{
			"FifoQueue":                 "true",
			"ContentBasedDeduplication": "true",
		})
		if err != nil {
			return fmt.Errorf("create delivery queue: %w", err)
		}
		defer cleanupQueue()

		queueARN, err := queueArn(tc.ctx, sqsClient, queueURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}
		subResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: topicArn.TopicArn,
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %w", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: subResp.SubscriptionArn})

		const messages = 5
		for i := 0; i < messages; i++ {
			if _, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
				TopicArn:               topicArn.TopicArn,
				Message:                aws.String(fmt.Sprintf("ordered-%d", i)),
				MessageGroupId:         aws.String("order-group"),
				MessageDeduplicationId: aws.String(fmt.Sprintf("order-dedup-%d", i)),
			}); err != nil {
				return fmt.Errorf("publish %d: %w", i, err)
			}
		}

		// A FIFO receive returns at most one message per group, and the
		// received copy blocks the group's remaining messages until it is
		// deleted — record each body and delete the copy to unblock the
		// next, so the recorded order is the queue's group order.
		var bodies []string
		deadline := time.Now().Add(20 * time.Second)
		for len(bodies) < messages && time.Now().Before(deadline) {
			resp, err := sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            queueURL,
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     1,
			})
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			for _, m := range resp.Messages {
				var envelope struct {
					Message string `json:"Message"`
				}
				if err := json.Unmarshal([]byte(aws.ToString(m.Body)), &envelope); err != nil {
					return fmt.Errorf("delivered body is not the notification envelope: %v", err)
				}
				bodies = append(bodies, envelope.Message)
				if _, err := sqsClient.DeleteMessage(tc.ctx, &sqs.DeleteMessageInput{
					QueueUrl:      queueURL,
					ReceiptHandle: m.ReceiptHandle,
				}); err != nil {
					return fmt.Errorf("delete received copy: %w", err)
				}
			}
		}
		if len(bodies) != messages {
			return fmt.Errorf("FIFO topic delivered %d of %d same-group messages", len(bodies), messages)
		}
		for i, body := range bodies {
			if want := fmt.Sprintf("ordered-%d", i); body != want {
				return fmt.Errorf("arrival %d is %q, want %q — same-group delivery did not preserve publish order (%v)", i, body, want, bodies)
			}
		}
		return nil
	}))

	// MessageGroupId is documented as optional on standard topics: it "is
	// forwarded only to Amazon SQS standard subscriptions to activate fair
	// queues". The publish is accepted and the group identifier reaches the
	// subscribed queue's message.
	results = append(results, r.RunTest("sns", "Publish_StandardTopicMessageGroupIdForwarded", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("FairQueueTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		queueURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("FairQueueDest"), nil)
		if err != nil {
			return fmt.Errorf("create delivery queue: %w", err)
		}
		defer cleanupQueue()

		queueARN, err := queueArn(tc.ctx, sqsClient, queueURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}
		subResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %w", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: subResp.SubscriptionArn})

		pubResp, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn:       aws.String(topicArn),
			Message:        aws.String("fair queues activation message"),
			MessageGroupId: aws.String("tenant-a"),
		})
		if err != nil {
			return fmt.Errorf("MessageGroupId on a standard topic is valid (fair queues): %v", err)
		}
		if pubResp.SequenceNumber != nil {
			return fmt.Errorf("standard publish response carries SequenceNumber %q — the member applies only to FIFO topics", *pubResp.SequenceNumber)
		}

		deadline := time.Now().Add(20 * time.Second)
		for {
			resp, err := sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:                    queueURL,
				MaxNumberOfMessages:         10,
				WaitTimeSeconds:             1,
				MessageSystemAttributeNames: []sqstypes.MessageSystemAttributeName{sqstypes.MessageSystemAttributeNameMessageGroupId},
			})
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			if len(resp.Messages) > 0 {
				if got := resp.Messages[0].Attributes["MessageGroupId"]; got != "tenant-a" {
					return fmt.Errorf("forwarded MessageGroupId = %q, want tenant-a", got)
				}
				return nil
			}
			if !time.Now().Before(deadline) {
				return fmt.Errorf("standard-topic publish with MessageGroupId delivered nothing to the subscribed queue")
			}
		}
	}))

	// The HTTP/S delivery policy's retry ladder is honoured: a subscription
	// policy of two one-second retries drives three delivery attempts at a
	// server that answers the first two notifications with 500 (a
	// documented retryable status) and the third with 200.
	results = append(results, r.RunTest("sns", "DeliveryPolicy_HTTPRetryLadder", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("LadderSdkTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		tokenCh := make(chan string, 1)
		var attempts int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Header.Get("x-amz-sns-message-type") {
			case "SubscriptionConfirmation":
				body, _ := io.ReadAll(r.Body)
				var payload map[string]interface{}
				if json.Unmarshal(body, &payload) == nil {
					if token, ok := payload["Token"].(string); ok {
						select {
						case tokenCh <- token:
						default:
						}
					}
				}
				w.WriteHeader(http.StatusOK)
			default:
				// Notification: fail twice, then accept.
				if atomic.AddInt32(&attempts, 1) <= 2 {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}
		}))
		defer srv.Close()

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(topicArn),
			Protocol:              aws.String("http"),
			Endpoint:              aws.String(srv.URL),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: sResp.SubscriptionArn})

		var token string
		select {
		case token = <-tokenCh:
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timed out waiting for subscription confirmation")
		}
		if _, err := tc.client.ConfirmSubscription(tc.ctx, &sns.ConfirmSubscriptionInput{
			TopicArn: aws.String(topicArn), Token: aws.String(token),
		}); err != nil {
			return fmt.Errorf("confirm: %v", err)
		}

		if _, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
			AttributeName:   aws.String("DeliveryPolicy"),
			AttributeValue:  aws.String(`{"healthyRetryPolicy":{"minDelayTarget":1,"maxDelayTarget":1,"numRetries":2}}`),
		}); err != nil {
			return fmt.Errorf("set delivery policy: %v", err)
		}

		if _, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(topicArn), Message: aws.String("retry ladder body"),
		}); err != nil {
			return err
		}

		deadline := time.Now().Add(15 * time.Second)
		for atomic.LoadInt32(&attempts) < 3 && time.Now().Before(deadline) {
			time.Sleep(50 * time.Millisecond)
		}
		if got := atomic.LoadInt32(&attempts); got != 3 {
			return fmt.Errorf("endpoint saw %d notification attempts, want 3 — the retry ladder did not run", got)
		}
		return nil
	}))

	return results
}
