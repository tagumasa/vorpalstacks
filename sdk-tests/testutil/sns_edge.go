package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func (r *TestRunner) runSNSEdgeTests(tc *snsTestContext) []TestResult {
	var results []TestResult
	reg := tc.region
	acct := tc.accountID

	results = append(results, r.RunTest("sns", "GetTopicAttributes_NonExistent", func() error {
		_, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic-xyz", reg, acct)),
		})
		return tc.expectNotFound("GetTopicAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "DeleteTopic_NonExistent", func() error {
		_, err := tc.client.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-del-topic", reg, acct)),
		})
		return tc.expectNotFound("DeleteTopic", err)
	}))

	results = append(results, r.RunTest("sns", "Publish_NonExistentTopic", func() error {
		_, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic-xyz", reg, acct)),
			Message:  aws.String("test message"),
		})
		return tc.expectNotFound("Publish", err)
	}))

	results = append(results, r.RunTest("sns", "Subscribe_NonExistentTopic", func() error {
		_, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-sub-topic", reg, acct)),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:fake-queue", reg, acct)),
		})
		return tc.expectNotFound("Subscribe", err)
	}))

	results = append(results, r.RunTest("sns", "Unsubscribe_NonExistent", func() error {
		_, err := tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{
			SubscriptionArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic:fake-sub-id", reg, acct)),
		})
		return tc.expectNotFound("Unsubscribe", err)
	}))

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_NonExistent", func() error {
		_, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic:fake-sub-id", reg, acct)),
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		})
		return tc.expectNotFound("SetSubscriptionAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "GetSubscriptionAttributes_NonExistent", func() error {
		_, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic:fake-sub-id", reg, acct)),
		})
		return tc.expectNotFound("GetSubscriptionAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes_NonExistent", func() error {
		_, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn:       aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-setattr-topic", reg, acct)),
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("test"),
		})
		return tc.expectNotFound("SetTopicAttributes", err)
	}))

	// --- New validation edge tests ---

	// Invalid protocol must be rejected at Subscribe time.
	results = append(results, r.RunTest("sns", "Subscribe_InvalidProtocol", func() error {
		_, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:any-topic", reg, acct)),
			Protocol: aws.String("foobar"),
			Endpoint: aws.String("https://example.com/hook"),
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	// Endpoint format must be validated per protocol.
	results = append(results, r.RunTest("sns", "Subscribe_InvalidEndpoint", func() error {
		_, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:any-topic", reg, acct)),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("not-a-valid-queue-url"),
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	// Reserved AWS prefix must be rejected.
	results = append(results, r.RunTest("sns", "CreateTopic_ReservedPrefix", func() error {
		_, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String("aws-reserved-test-topic"),
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	// FifoTopic=true with non-.fifo name must be rejected.
	results = append(results, r.RunTest("sns", "CreateTopic_FifoConsistencyMismatch", func() error {
		_, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(tc.uniqueName("no-suffix-fifo-test")),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	// More than 10 message attributes must be rejected.
	results = append(results, r.RunTest("sns", "Publish_TooManyMessageAttributes", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("too-many-attrs-test"))
		if err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}
		defer tc.deleteTopic(topicArn)

		msgAttrs := make(map[string]types.MessageAttributeValue, 11)
		for i := 0; i < 11; i++ {
			key := fmt.Sprintf("attr%d", i)
			msgAttrs[key] = types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("v"),
			}
		}

		_, err = tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn:          aws.String(topicArn),
			Message:           aws.String("test"),
			MessageAttributes: msgAttrs,
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	// BatchRequestTooLong: two entries whose combined size exceeds the 256 KB
	// batch limit. Each individual message passes the single-message size
	// check, but the total triggers BatchRequestTooLong as a top-level error.
	// The two-pass design ensures NO entry is delivered before the error.
	results = append(results, r.RunTest("sns", "PublishBatch_BatchRequestTooLong", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("BatchTooLong"))
		if err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}
		defer tc.deleteTopic(topicArn)

		bigMsg := strings.Repeat("x", 140000)
		_, err = tc.client.PublishBatch(tc.ctx, &sns.PublishBatchInput{
			TopicArn: aws.String(topicArn),
			PublishBatchRequestEntries: []types.PublishBatchRequestEntry{
				{Id: aws.String("e1"), Message: aws.String(bigMsg)},
				{Id: aws.String("e2"), Message: aws.String(bigMsg)},
			},
		})
		return expectAWSErrorCode(err, "BatchRequestTooLong")
	}))

	// CreatePlatformApplication must enforce the same attribute value length
	// cap as SetPlatformApplicationAttributes (8192 chars).
	results = append(results, r.RunTest("sns", "CreatePlatformApplication_AttributeValueTooLong", func() error {
		_, err := tc.client.CreatePlatformApplication(tc.ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(tc.uniqueName("TooLongAttrApp")),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": strings.Repeat("x", 8193),
			},
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	return results
}
