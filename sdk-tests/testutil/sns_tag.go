package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func (r *TestRunner) runSNSTagTests(tc *snsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sns", "CreateTopic_WithTags", func() error {
		tagTopicName := tc.uniqueName("TagTopic")
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(tagTopicName),
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("Prod")},
				{Key: aws.String("Team"), Value: aws.String("Backend")},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*resp.TopicArn)

		tags, err := tc.listTags(*resp.TopicArn)
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tags) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(tags))
		}
		if v, ok := snsTagValue(tags, "Env"); !ok || v != "Prod" {
			return fmt.Errorf("Env tag mismatch: got %q", v)
		}
		if v, ok := snsTagValue(tags, "Team"); !ok || v != "Backend" {
			return fmt.Errorf("Team tag mismatch: got %q", v)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "TagResource", func() error {
		topicName := tc.uniqueName("TagResTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		_, err = tc.client.TagResource(tc.ctx, &sns.TagResourceInput{
			ResourceArn: aws.String(topicArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("Test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "ListTagsForResource", func() error {
		topicName := tc.uniqueName("ListTagTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		_, err = tc.client.TagResource(tc.ctx, &sns.TagResourceInput{
			ResourceArn: aws.String(topicArn),
			Tags: []types.Tag{
				{Key: aws.String("ListKey"), Value: aws.String("ListVal")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}

		tags, err := tc.listTags(topicArn)
		if err != nil {
			return err
		}
		if v, ok := snsTagValue(tags, "ListKey"); !ok || v != "ListVal" {
			return fmt.Errorf("tag ListKey=ListVal not found in ListTagsForResource response")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "UntagResource", func() error {
		topicName := tc.uniqueName("UntagTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		_, err = tc.client.TagResource(tc.ctx, &sns.TagResourceInput{
			ResourceArn: aws.String(topicArn),
			Tags: []types.Tag{
				{Key: aws.String("RemoveMe"), Value: aws.String("gone")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}

		_, err = tc.client.UntagResource(tc.ctx, &sns.UntagResourceInput{
			ResourceArn: aws.String(topicArn),
			TagKeys:     []string{"RemoveMe"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		tags, err := tc.listTags(topicArn)
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if _, ok := snsTagValue(tags, "RemoveMe"); ok {
			return fmt.Errorf("tag RemoveMe should have been removed")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "TagResource_MultipleTags", func() error {
		multiTagTopicArn, err := tc.createTopic(tc.uniqueName("MultiTagTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(multiTagTopicArn)

		_, err = tc.client.TagResource(tc.ctx, &sns.TagResourceInput{
			ResourceArn: aws.String(multiTagTopicArn),
			Tags: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("Val1")},
				{Key: aws.String("Key2"), Value: aws.String("Val2")},
				{Key: aws.String("Key3"), Value: aws.String("Val3")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}

		tags, err := tc.listTags(multiTagTopicArn)
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tags) < 3 {
			return fmt.Errorf("expected at least 3 tags, got %d", len(tags))
		}
		return nil
	}))

	// Tag operations against a topic that does not exist fail with the
	// model's ResourceNotFound wire code, as the service model specifies.
	results = append(results, r.RunTest("sns", "TagResource_NonExistentTopic", func() error {
		arn := fmt.Sprintf("arn:aws:sns:%s:%s:no-such-topic-%d", tc.region, tc.accountID, time.Now().UnixNano())
		_, err := tc.client.TagResource(tc.ctx, &sns.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags:        []types.Tag{{Key: aws.String("Environment"), Value: aws.String("test")}},
		})
		if err := tc.expectResourceNotFound("TagResource", err); err != nil {
			return err
		}
		_, err = tc.client.UntagResource(tc.ctx, &sns.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"Environment"},
		})
		if err := tc.expectResourceNotFound("UntagResource", err); err != nil {
			return err
		}
		_, err = tc.listTags(arn)
		return tc.expectResourceNotFound("ListTagsForResource", err)
	}))

	return results
}
