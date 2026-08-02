package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runSNSEdgeTests(tc *snsTestContext) []TestResult {
	var results []TestResult
	reg := tc.region
	acct := tc.accountID

	results = append(results, r.RunTest("sns", "GetTopicAttributes_NonExistent", func() error {
		_, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic-xyz", reg, acct)),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "DeleteTopic_NonExistent", func() error {
		_, err := tc.client.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-del-topic", reg, acct)),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "Publish_NonExistentTopic", func() error {
		_, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic-xyz", reg, acct)),
			Message:  aws.String("test message"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "Subscribe_NonExistentTopic", func() error {
		_, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-sub-topic", reg, acct)),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:fake-queue", reg, acct)),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "Unsubscribe_NonExistent", func() error {
		_, err := tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{
			SubscriptionArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic:fake-sub-id", reg, acct)),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_NonExistent", func() error {
		_, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic:fake-sub-id", reg, acct)),
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "GetSubscriptionAttributes_NonExistent", func() error {
		_, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-topic:fake-sub-id", reg, acct)),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes_NonExistent", func() error {
		_, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn:       aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-setattr-topic", reg, acct)),
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("test"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	return results
}
