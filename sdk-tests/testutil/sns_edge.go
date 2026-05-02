package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runSNSEdgeTests(tc *snsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sns", "GetTopicAttributes_NonExistent", func() error {
		_, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic-xyz"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "DeleteTopic_NonExistent", func() error {
		_, err := tc.client.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-del-topic"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "Publish_NonExistentTopic", func() error {
		_, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic-xyz"),
			Message:  aws.String("test message"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "Subscribe_NonExistentTopic", func() error {
		_, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-sub-topic"),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:fake-queue"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "Unsubscribe_NonExistent", func() error {
		_, err := tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{
			SubscriptionArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic:fake-sub-id"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_NonExistent", func() error {
		_, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic:fake-sub-id"),
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "GetSubscriptionAttributes_NonExistent", func() error {
		_, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic:fake-sub-id"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes_NonExistent", func() error {
		_, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn:       aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-setattr-topic"),
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("test"),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	return results
}
