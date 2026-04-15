package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSNSTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "sns",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := sns.NewFromConfig(cfg)
	ctx := context.Background()

	topicName := fmt.Sprintf("TestTopic-%d", time.Now().UnixNano())
	var topicArn string

	results = append(results, r.RunTest("sns", "CreateTopic", func() error {
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(topicName),
		})
		if err != nil {
			return err
		}
		topicArn = *resp.TopicArn
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListTopics", func() error {
		resp, err := client.ListTopics(ctx, &sns.ListTopicsInput{})
		if err != nil {
			return err
		}
		if resp.Topics == nil {
			return fmt.Errorf("Topics is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetTopicAttributes", func() error {
		resp, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String(topicArn),
		})
		if err != nil {
			return err
		}
		if resp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes", func() error {
		_, err := client.SetTopicAttributes(ctx, &sns.SetTopicAttributesInput{
			TopicArn:       aws.String(topicArn),
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("Test Topic"),
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "Publish", func() error {
		resp, err := client.Publish(ctx, &sns.PublishInput{
			TopicArn: aws.String(topicArn),
			Message:  aws.String("Test message"),
		})
		if err != nil {
			return err
		}
		if resp.MessageId == nil {
			return fmt.Errorf("MessageId is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "AddPermission", func() error {
		_, err := client.AddPermission(ctx, &sns.AddPermissionInput{
			TopicArn:     aws.String(topicArn),
			Label:        aws.String("TestPermission"),
			AWSAccountId: []string{"000000000000"},
			ActionName:   []string{"Publish"},
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "RemovePermission", func() error {
		_, err := client.RemovePermission(ctx, &sns.RemovePermissionInput{
			TopicArn: aws.String(topicArn),
			Label:    aws.String("TestPermission"),
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "Subscribe", func() error {
		resp, err := client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("email"),
			Endpoint: aws.String("test@example.com"),
		})
		if err != nil {
			return err
		}
		if resp.SubscriptionArn == nil {
			return fmt.Errorf("SubscriptionArn is nil")
		}
		if *resp.SubscriptionArn != "pending confirmation" {
			client.Unsubscribe(ctx, &sns.UnsubscribeInput{SubscriptionArn: resp.SubscriptionArn})
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptions", func() error {
		resp, err := client.ListSubscriptions(ctx, &sns.ListSubscriptionsInput{})
		if err != nil {
			return err
		}
		if resp.Subscriptions == nil {
			return fmt.Errorf("Subscriptions is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "TagResource", func() error {
		_, err := client.TagResource(ctx, &sns.TagResourceInput{
			ResourceArn: aws.String(topicArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("Test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &sns.ListTagsForResourceInput{
			ResourceArn: aws.String(topicArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("Tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &sns.UntagResourceInput{
			ResourceArn: aws.String(topicArn),
			TagKeys:     []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "PublishBatch", func() error {
		resp, err := client.PublishBatch(ctx, &sns.PublishBatchInput{
			TopicArn: aws.String(topicArn),
			PublishBatchRequestEntries: []types.PublishBatchRequestEntry{
				{
					Id:      aws.String("msg1"),
					Message: aws.String("Batch message 1"),
				},
				{
					Id:      aws.String("msg2"),
					Message: aws.String("Batch message 2"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Successful == nil {
			return fmt.Errorf("Successful is nil")
		}
		if len(resp.Successful) == 0 {
			return fmt.Errorf("expected at least one successful publish, got 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "CreateTopic (FIFO)", func() error {
		fifoTopicName := fmt.Sprintf("TestFifoTopic-%d.fifo", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoTopicName),
			Attributes: map[string]string{
				"ContentBasedDeduplication": "true",
				"FifoTopic":                 "true",
			},
		})
		if err != nil {
			return err
		}
		if resp.TopicArn == nil {
			return fmt.Errorf("TopicArn is nil")
		}
		client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp.TopicArn})
		return nil
	}))

	// === GROUP A: CORE LIFECYCLE ===

	var delTopicArn string
	results = append(results, r.RunTest("sns", "DeleteTopic", func() error {
		delTopicName := fmt.Sprintf("DelTopic-%d", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(delTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		delTopicArn = *resp.TopicArn
		_, err = client.DeleteTopic(ctx, &sns.DeleteTopicInput{
			TopicArn: resp.TopicArn,
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "DeleteTopic_VerifyGone", func() error {
		_, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String(delTopicArn),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "DeleteTopic_NonExistent", func() error {
		_, err := client.DeleteTopic(ctx, &sns.DeleteTopicInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-del-topic"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	var unsubSubArn string
	results = append(results, r.RunTest("sns", "Unsubscribe", func() error {
		unsubTopicName := fmt.Sprintf("UnsubTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(unsubTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		sResp, err := client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:fake-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		unsubSubArn = *sResp.SubscriptionArn

		_, err = client.Unsubscribe(ctx, &sns.UnsubscribeInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "Unsubscribe_VerifyGone", func() error {
		_, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(unsubSubArn),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	var sqsSubArn string
	var sqsSubTopicArn string
	results = append(results, r.RunTest("sns", "Subscribe_SQS_AutoConfirmed", func() error {
		sqsTopicName := fmt.Sprintf("SqsTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(sqsTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		sqsSubTopicArn = *tResp.TopicArn

		sResp, err := client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:auto-confirm-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		sqsSubArn = *sResp.SubscriptionArn

		getResp, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "false" {
			return fmt.Errorf("SQS subscription should be auto-confirmed, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetSubscriptionAttributes", func() error {
		resp, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(sqsSubArn),
		})
		if err != nil {
			return err
		}
		if resp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if resp.Attributes["SubscriptionArn"] != sqsSubArn {
			return fmt.Errorf("SubscriptionArn mismatch: got %q", resp.Attributes["SubscriptionArn"])
		}
		if resp.Attributes["Protocol"] != "sqs" {
			return fmt.Errorf("Protocol mismatch: got %q", resp.Attributes["Protocol"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes", func() error {
		_, err := client.SetSubscriptionAttributes(ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(sqsSubArn),
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(sqsSubArn),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["RawMessageDelivery"] != "true" {
			return fmt.Errorf("RawMessageDelivery not set: got %q", getResp.Attributes["RawMessageDelivery"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptionsByTopic", func() error {
		lstTopicName := fmt.Sprintf("LstSubTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(lstTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		_, err = client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:list-sub-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		listResp, err := client.ListSubscriptionsByTopic(ctx, &sns.ListSubscriptionsByTopicInput{
			TopicArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("list by topic: %v", err)
		}
		if len(listResp.Subscriptions) == 0 {
			return fmt.Errorf("expected at least one subscription")
		}
		found := false
		for _, sub := range listResp.Subscriptions {
			if sub.TopicArn != nil && *sub.TopicArn == *tResp.TopicArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("subscription for topic not found in ListSubscriptionsByTopic")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptions_ContainsCreated", func() error {
		lsTopicName := fmt.Sprintf("LSTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(lsTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		sResp, err := client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:list-all-sub-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		listResp, err := client.ListSubscriptions(ctx, &sns.ListSubscriptionsInput{})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, sub := range listResp.Subscriptions {
			if sub.SubscriptionArn != nil && *sub.SubscriptionArn == *sResp.SubscriptionArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created subscription not found in ListSubscriptions")
		}
		return nil
	}))

	// === GROUP B: FIFO PUBLISH ===

	var fifoTopicArn string
	results = append(results, r.RunTest("sns", "Publish_FIFO_WithMessageGroupId", func() error {
		fifoName := fmt.Sprintf("FifoPubTopic-%d.fifo", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoName),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		fifoTopicArn = *resp.TopicArn
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: aws.String(fifoTopicArn)})

		pubResp, err := client.Publish(ctx, &sns.PublishInput{
			TopicArn:               aws.String(fifoTopicArn),
			Message:                aws.String("FIFO message 1"),
			MessageGroupId:         aws.String("group1"),
			MessageDeduplicationId: aws.String("dedup-fifo-1"),
		})
		if err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		if pubResp.MessageId == nil {
			return fmt.Errorf("MessageId is nil")
		}
		if pubResp.SequenceNumber == nil {
			return fmt.Errorf("SequenceNumber is nil for FIFO topic")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_FIFO_ContentBasedDedup", func() error {
		fifoDedupName := fmt.Sprintf("FifoDedupTopic-%d.fifo", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoDedupName),
			Attributes: map[string]string{
				"FifoTopic":                 "true",
				"ContentBasedDeduplication": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp.TopicArn})

		msg := "dedup test message content"
		pub1, err := client.Publish(ctx, &sns.PublishInput{
			TopicArn:       resp.TopicArn,
			Message:        aws.String(msg),
			MessageGroupId: aws.String("dedup-group"),
		})
		if err != nil {
			return fmt.Errorf("publish 1: %v", err)
		}

		pub2, err := client.Publish(ctx, &sns.PublishInput{
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
		fifoDedupIdName := fmt.Sprintf("FifoDedupIdTopic-%d.fifo", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoDedupIdName),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp.TopicArn})

		pub1, err := client.Publish(ctx, &sns.PublishInput{
			TopicArn:               resp.TopicArn,
			Message:                aws.String("msg A"),
			MessageGroupId:         aws.String("group-dedup-id"),
			MessageDeduplicationId: aws.String("dedup-123"),
		})
		if err != nil {
			return fmt.Errorf("publish 1: %v", err)
		}

		pub2, err := client.Publish(ctx, &sns.PublishInput{
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
		return nil
	}))

	// === GROUP C: PUBLISH WITH MESSAGE ATTRIBUTES ===

	results = append(results, r.RunTest("sns", "Publish_WithMessageAttributes", func() error {
		attrTopicName := fmt.Sprintf("AttrPubTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(attrTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		pubResp, err := client.Publish(ctx, &sns.PublishInput{
			TopicArn: tResp.TopicArn,
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
		if pubResp.MessageId == nil {
			return fmt.Errorf("MessageId is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "PublishBatch_WithAttributes", func() error {
		batchAttrTopicName := fmt.Sprintf("BatchAttrTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(batchAttrTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		resp, err := client.PublishBatch(ctx, &sns.PublishBatchInput{
			TopicArn: tResp.TopicArn,
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

	// === GROUP D: PUBLISHBATCH EDGE CASES ===

	results = append(results, r.RunTest("sns", "PublishBatch_MaxEntries", func() error {
		maxBatchTopicName := fmt.Sprintf("MaxBatchTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(maxBatchTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		entries := make([]types.PublishBatchRequestEntry, 10)
		for i := 0; i < 10; i++ {
			entries[i] = types.PublishBatchRequestEntry{
				Id:      aws.String(fmt.Sprintf("max%d", i)),
				Message: aws.String(fmt.Sprintf("max batch message %d", i)),
			}
		}
		resp, err := client.PublishBatch(ctx, &sns.PublishBatchInput{
			TopicArn:                   tResp.TopicArn,
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
		failBatchFifoName := fmt.Sprintf("FailBatchFifo-%d.fifo", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(failBatchFifoName),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		resp, err := client.PublishBatch(ctx, &sns.PublishBatchInput{
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

	// === GROUP E: DATA PROTECTION POLICY ===

	results = append(results, r.RunTest("sns", "PutDataProtectionPolicy", func() error {
		dppTopicName := fmt.Sprintf("DppTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(dppTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		policy := `{"Name":"data-protection-policy","Description":"Test policy","Version":"2021-06-01"}`
		_, err = client.PutDataProtectionPolicy(ctx, &sns.PutDataProtectionPolicyInput{
			ResourceArn:          tResp.TopicArn,
			DataProtectionPolicy: aws.String(policy),
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "GetDataProtectionPolicy", func() error {
		dppGetTopicName := fmt.Sprintf("DppGetTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(dppGetTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		policy := `{"Name":"round-trip-policy","Version":"2021-06-01"}`
		_, err = client.PutDataProtectionPolicy(ctx, &sns.PutDataProtectionPolicyInput{
			ResourceArn:          tResp.TopicArn,
			DataProtectionPolicy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		getResp, err := client.GetDataProtectionPolicy(ctx, &sns.GetDataProtectionPolicyInput{
			ResourceArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.DataProtectionPolicy == nil {
			return fmt.Errorf("DataProtectionPolicy is nil")
		}
		if *getResp.DataProtectionPolicy != policy {
			return fmt.Errorf("policy mismatch: got %q", *getResp.DataProtectionPolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetDataProtectionPolicy_NonExistent", func() error {
		_, err := client.GetDataProtectionPolicy(ctx, &sns.GetDataProtectionPolicyInput{
			ResourceArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-dpp-topic"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	// === GROUP F: PLATFORM APPLICATIONS ===

	var platformAppArn string
	results = append(results, r.RunTest("sns", "CreatePlatformApplication", func() error {
		appName := fmt.Sprintf("TestApp-%d", time.Now().UnixNano())
		resp, err := client.CreatePlatformApplication(ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(appName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "fake-credential",
			},
		})
		if err != nil {
			return err
		}
		if resp.PlatformApplicationArn == nil {
			return fmt.Errorf("PlatformApplicationArn is nil")
		}
		platformAppArn = *resp.PlatformApplicationArn
		return nil
	}))

	results = append(results, r.RunTest("sns", "CreatePlatformApplication_Duplicate", func() error {
		if platformAppArn == "" {
			return fmt.Errorf("platformAppArn not set from previous test")
		}
		appName := fmt.Sprintf("TestApp-%d", time.Now().UnixNano())
		dupResp, err := client.CreatePlatformApplication(ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(appName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "fake-credential",
			},
		})
		if err != nil {
			return fmt.Errorf("create with unique name should succeed: %v", err)
		}
		client.DeletePlatformApplication(ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: dupResp.PlatformApplicationArn,
		})
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListPlatformApplications", func() error {
		resp, err := client.ListPlatformApplications(ctx, &sns.ListPlatformApplicationsInput{})
		if err != nil {
			return err
		}
		if resp.PlatformApplications == nil {
			return fmt.Errorf("PlatformApplications is nil")
		}
		found := false
		for _, app := range resp.PlatformApplications {
			if app.PlatformApplicationArn != nil && *app.PlatformApplicationArn == platformAppArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created platform application not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetPlatformApplicationAttributes", func() error {
		resp, err := client.GetPlatformApplicationAttributes(ctx, &sns.GetPlatformApplicationAttributesInput{
			PlatformApplicationArn: aws.String(platformAppArn),
		})
		if err != nil {
			return err
		}
		if resp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if resp.Attributes["PlatformCredential"] != "fake-credential" {
			return fmt.Errorf("PlatformCredential mismatch: got %q", resp.Attributes["PlatformCredential"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetPlatformApplicationAttributes", func() error {
		_, err := client.SetPlatformApplicationAttributes(ctx, &sns.SetPlatformApplicationAttributesInput{
			PlatformApplicationArn: aws.String(platformAppArn),
			Attributes: map[string]string{
				"EventEndpointCreated": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := client.GetPlatformApplicationAttributes(ctx, &sns.GetPlatformApplicationAttributesInput{
			PlatformApplicationArn: aws.String(platformAppArn),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["EventEndpointCreated"] != "true" {
			return fmt.Errorf("EventEndpointCreated not set: got %q", getResp.Attributes["EventEndpointCreated"])
		}
		return nil
	}))

	// === GROUP G: PLATFORM ENDPOINTS ===

	var endpointArn string
	results = append(results, r.RunTest("sns", "CreatePlatformEndpoint", func() error {
		resp, err := client.CreatePlatformEndpoint(ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(platformAppArn),
			Token:                  aws.String("fake-device-token-12345"),
			CustomUserData:         aws.String("user-data-123"),
		})
		if err != nil {
			return err
		}
		if resp.EndpointArn == nil {
			return fmt.Errorf("EndpointArn is nil")
		}
		endpointArn = *resp.EndpointArn
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListEndpointsByPlatformApplication", func() error {
		resp, err := client.ListEndpointsByPlatformApplication(ctx, &sns.ListEndpointsByPlatformApplicationInput{
			PlatformApplicationArn: aws.String(platformAppArn),
		})
		if err != nil {
			return err
		}
		if resp.Endpoints == nil {
			return fmt.Errorf("Endpoints is nil")
		}
		if len(resp.Endpoints) == 0 {
			return fmt.Errorf("expected at least one endpoint")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetEndpointAttributes", func() error {
		resp, err := client.GetEndpointAttributes(ctx, &sns.GetEndpointAttributesInput{
			EndpointArn: aws.String(endpointArn),
		})
		if err != nil {
			return err
		}
		if resp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if resp.Attributes["Token"] != "fake-device-token-12345" {
			return fmt.Errorf("Token mismatch: got %q", resp.Attributes["Token"])
		}
		if resp.Attributes["Enabled"] != "true" {
			return fmt.Errorf("Enabled should be true by default, got %q", resp.Attributes["Enabled"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetEndpointAttributes", func() error {
		_, err := client.SetEndpointAttributes(ctx, &sns.SetEndpointAttributesInput{
			EndpointArn: aws.String(endpointArn),
			Attributes: map[string]string{
				"CustomUserData": "updated-user-data",
				"Enabled":        "false",
			},
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := client.GetEndpointAttributes(ctx, &sns.GetEndpointAttributesInput{
			EndpointArn: aws.String(endpointArn),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["CustomUserData"] != "updated-user-data" {
			return fmt.Errorf("CustomUserData mismatch: got %q", getResp.Attributes["CustomUserData"])
		}
		if getResp.Attributes["Enabled"] != "false" {
			return fmt.Errorf("Enabled should be false, got %q", getResp.Attributes["Enabled"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "DeletePlatformApplication_Cascade", func() error {
		cascadeAppName := fmt.Sprintf("CascadeApp-%d", time.Now().UnixNano())
		appResp, err := client.CreatePlatformApplication(ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(cascadeAppName),
			Platform: aws.String("APNS"),
			Attributes: map[string]string{
				"PlatformCredential": "cascade-cred",
			},
		})
		if err != nil {
			return fmt.Errorf("create app: %v", err)
		}
		cascadeAppArn := *appResp.PlatformApplicationArn

		epResp, err := client.CreatePlatformEndpoint(ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(cascadeAppArn),
			Token:                  aws.String("cascade-token-123"),
		})
		if err != nil {
			return fmt.Errorf("create endpoint: %v", err)
		}
		cascadeEndpointArn := *epResp.EndpointArn

		_, err = client.DeletePlatformApplication(ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: aws.String(cascadeAppArn),
		})
		if err != nil {
			return fmt.Errorf("delete app: %v", err)
		}

		_, err = client.GetEndpointAttributes(ctx, &sns.GetEndpointAttributesInput{
			EndpointArn: aws.String(cascadeEndpointArn),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "DeleteEndpoint", func() error {
		delEpAppName := fmt.Sprintf("DelEpApp-%d", time.Now().UnixNano())
		appResp, err := client.CreatePlatformApplication(ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(delEpAppName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "del-ep-cred",
			},
		})
		if err != nil {
			return fmt.Errorf("create app: %v", err)
		}
		defer client.DeletePlatformApplication(ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
		})

		epResp, err := client.CreatePlatformEndpoint(ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
			Token:                  aws.String("del-ep-token"),
		})
		if err != nil {
			return fmt.Errorf("create endpoint: %v", err)
		}

		_, err = client.DeleteEndpoint(ctx, &sns.DeleteEndpointInput{
			EndpointArn: epResp.EndpointArn,
		})
		if err != nil {
			return fmt.Errorf("delete endpoint: %v", err)
		}

		_, err = client.GetEndpointAttributes(ctx, &sns.GetEndpointAttributesInput{
			EndpointArn: epResp.EndpointArn,
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	// === GROUP H: ERROR CASES ===

	results = append(results, r.RunTest("sns", "GetTopicAttributes_NonExistent", func() error {
		_, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic-xyz"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_NonExistentTopic", func() error {
		_, err := client.Publish(ctx, &sns.PublishInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic-xyz"),
			Message:  aws.String("test message"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "CreateTopic_DuplicateIdempotent", func() error {
		dupTopicName := fmt.Sprintf("DupTopic-%d", time.Now().UnixNano())
		resp1, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(dupTopicName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp1.TopicArn})

		resp2, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(dupTopicName),
		})
		if err != nil {
			return fmt.Errorf("duplicate create should be idempotent, got: %v", err)
		}
		if *resp1.TopicArn != *resp2.TopicArn {
			return fmt.Errorf("ARN mismatch: %q vs %q", *resp1.TopicArn, *resp2.TopicArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListTopics_ContainsCreated", func() error {
		ltTopicName := fmt.Sprintf("LTTopic-%d", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(ltTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp.TopicArn})

		listResp, err := client.ListTopics(ctx, &sns.ListTopicsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, t := range listResp.Topics {
			if t.TopicArn != nil && *t.TopicArn == *resp.TopicArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created topic not found in ListTopics")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes_GetVerify", func() error {
		attrTopicName := fmt.Sprintf("AttrTopic-%d", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(attrTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp.TopicArn})

		_, err = client.SetTopicAttributes(ctx, &sns.SetTopicAttributesInput{
			TopicArn:       resp.TopicArn,
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("MyDisplayName"),
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
			TopicArn: resp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes == nil {
			return fmt.Errorf("attributes is nil")
		}
		if getResp.Attributes["DisplayName"] != "MyDisplayName" {
			return fmt.Errorf("DisplayName mismatch: got %q", getResp.Attributes["DisplayName"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Subscribe_NonExistentTopic", func() error {
		_, err := client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-sub-topic"),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:fake-queue"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Unsubscribe_NonExistent", func() error {
		_, err := client.Unsubscribe(ctx, &sns.UnsubscribeInput{
			SubscriptionArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic:fake-sub-id"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_NonExistent", func() error {
		_, err := client.SetSubscriptionAttributes(ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic:fake-sub-id"),
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetSubscriptionAttributes_NonExistent", func() error {
		_, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic:fake-sub-id"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes_NonExistent", func() error {
		_, err := client.SetTopicAttributes(ctx, &sns.SetTopicAttributesInput{
			TopicArn:       aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-setattr-topic"),
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("test"),
		})
		if err := AssertErrorContains(err, "NotFound"); err != nil {
			return err
		}
		return nil
	}))

	// === GROUP I: SUBSCRIPTION EDGE CASES ===

	results = append(results, r.RunTest("sns", "Subscribe_EmailPendingConfirmation", func() error {
		emailTopicName := fmt.Sprintf("EmailTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(emailTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		sResp, err := client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("email"),
			Endpoint: aws.String("pending@example.com"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		getResp, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "true" {
			return fmt.Errorf("email subscription should be pending, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Subscribe_ApplicationPendingConfirmation", func() error {
		appTopicName := fmt.Sprintf("AppTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(appTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		sResp, err := client.Subscribe(ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("application"),
			Endpoint: aws.String("arn:aws:sns:us-east-1:000000000000:app/FAKE/fake-endpoint"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		getResp, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "true" {
			return fmt.Errorf("application subscription should be pending, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetTopicAttributes_FifoAttributes", func() error {
		fifoAttrName := fmt.Sprintf("FifoAttrTopic-%d.fifo", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoAttrName),
			Attributes: map[string]string{
				"FifoTopic":                 "true",
				"ContentBasedDeduplication": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		getResp, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
			TopicArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["FifoTopic"] != "true" {
			return fmt.Errorf("FifoTopic should be true, got %q", getResp.Attributes["FifoTopic"])
		}
		if getResp.Attributes["ContentBasedDeduplication"] != "true" {
			return fmt.Errorf("ContentBasedDeduplication should be true, got %q", getResp.Attributes["ContentBasedDeduplication"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetTopicAttributes_PolicyDefault", func() error {
		policyTopicName := fmt.Sprintf("PolicyTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(policyTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		getResp, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
			TopicArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		policy, ok := getResp.Attributes["Policy"]
		if !ok || policy == "" {
			return fmt.Errorf("default Policy should be set")
		}
		if !strings.Contains(policy, "Version") {
			return fmt.Errorf("Policy should contain Version, got: %s", policy)
		}
		if !strings.Contains(*tResp.TopicArn, policyTopicName) {
			return fmt.Errorf("Policy should reference topic ARN")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "CreateTopic_WithTags", func() error {
		tagTopicName := fmt.Sprintf("TagTopic-%d", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(tagTopicName),
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("Prod")},
				{Key: aws.String("Team"), Value: aws.String("Backend")},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp.TopicArn})

		listResp, err := client.ListTagsForResource(ctx, &sns.ListTagsForResourceInput{
			ResourceArn: resp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.Tags) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(listResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "TagResource_MultipleTags", func() error {
		multiTagTopicName := fmt.Sprintf("MultiTagTopic-%d", time.Now().UnixNano())
		resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(multiTagTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: resp.TopicArn})

		_, err = client.TagResource(ctx, &sns.TagResourceInput{
			ResourceArn: resp.TopicArn,
			Tags: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("Val1")},
				{Key: aws.String("Key2"), Value: aws.String("Val2")},
				{Key: aws.String("Key3"), Value: aws.String("Val3")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}

		listResp, err := client.ListTagsForResource(ctx, &sns.ListTagsForResourceInput{
			ResourceArn: resp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.Tags) < 3 {
			return fmt.Errorf("expected at least 3 tags, got %d", len(listResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptionsByTopic_Empty", func() error {
		emptySubTopicName := fmt.Sprintf("EmptySubTopic-%d", time.Now().UnixNano())
		tResp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(emptySubTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: tResp.TopicArn})

		listResp, err := client.ListSubscriptionsByTopic(ctx, &sns.ListSubscriptionsByTopicInput{
			TopicArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(listResp.Subscriptions) != 0 {
			return fmt.Errorf("expected 0 subscriptions for new topic, got %d", len(listResp.Subscriptions))
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListTopics_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgTopicARNs []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagTopic-%s-%d", pgTs, i)
			resp, err := client.CreateTopic(ctx, &sns.CreateTopicInput{Name: aws.String(name)})
			if err != nil {
				for _, arn := range pgTopicARNs {
					client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: aws.String(arn)})
				}
				return fmt.Errorf("create topic %s: %v", name, err)
			}
			pgTopicARNs = append(pgTopicARNs, *resp.TopicArn)
		}

		var allTopics []string
		var nextToken *string
		for {
			resp, err := client.ListTopics(ctx, &sns.ListTopicsInput{NextToken: nextToken})
			if err != nil {
				for _, arn := range pgTopicARNs {
					client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: aws.String(arn)})
				}
				return fmt.Errorf("list topics page: %v", err)
			}
			for _, t := range resp.Topics {
				if strings.Contains(aws.ToString(t.TopicArn), "PagTopic-"+pgTs) {
					allTopics = append(allTopics, aws.ToString(t.TopicArn))
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, arn := range pgTopicARNs {
			client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: aws.String(arn)})
		}
		if len(allTopics) != 5 {
			return fmt.Errorf("expected 5 paginated topics, got %d", len(allTopics))
		}
		return nil
	}))

	if sqsSubTopicArn != "" {
		client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: aws.String(sqsSubTopicArn)})
	}

	if sqsSubArn != "" {
		client.Unsubscribe(ctx, &sns.UnsubscribeInput{SubscriptionArn: aws.String(sqsSubArn)})
	}

	if endpointArn != "" {
		client.DeleteEndpoint(ctx, &sns.DeleteEndpointInput{EndpointArn: aws.String(endpointArn)})
	}

	if platformAppArn != "" {
		client.DeletePlatformApplication(ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: aws.String(platformAppArn),
		})
	}

	if topicArn != "" {
		client.DeleteTopic(ctx, &sns.DeleteTopicInput{TopicArn: aws.String(topicArn)})
	}

	return results
}
