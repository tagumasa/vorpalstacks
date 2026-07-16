package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
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
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_WithMessageAttributes", func() error {
		attrTopicName := tc.uniqueName("AttrPubTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(attrTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		pubResp, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
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
		batchAttrTopicName := tc.uniqueName("BatchAttrTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(batchAttrTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		resp, err := tc.client.PublishBatch(tc.ctx, &sns.PublishBatchInput{
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

	results = append(results, r.RunTest("sns", "PublishBatch_MaxEntries", func() error {
		maxBatchTopicName := tc.uniqueName("MaxBatchTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(maxBatchTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		entries := make([]types.PublishBatchRequestEntry, 10)
		for i := 0; i < 10; i++ {
			entries[i] = types.PublishBatchRequestEntry{
				Id:      aws.String(fmt.Sprintf("max%d", i)),
				Message: aws.String(fmt.Sprintf("max batch message %d", i)),
			}
		}
		resp, err := tc.client.PublishBatch(tc.ctx, &sns.PublishBatchInput{
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

	return results
}
