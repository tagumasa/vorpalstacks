package testutil

import (
	"fmt"

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

		listResp, err := tc.client.ListTagsForResource(tc.ctx, &sns.ListTagsForResourceInput{
			ResourceArn: resp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.Tags) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(listResp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range listResp.Tags {
			tagMap[*t.Key] = *t.Value
		}
		if tagMap["Env"] != "Prod" {
			return fmt.Errorf("Env tag mismatch: got %q", tagMap["Env"])
		}
		if tagMap["Team"] != "Backend" {
			return fmt.Errorf("Team tag mismatch: got %q", tagMap["Team"])
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

		resp, err := tc.client.ListTagsForResource(tc.ctx, &sns.ListTagsForResourceInput{
			ResourceArn: aws.String(topicArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("Tags is nil")
		}
		found := false
		for _, t := range resp.Tags {
			if *t.Key == "ListKey" && *t.Value == "ListVal" {
				found = true
				break
			}
		}
		if !found {
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

		listResp, err := tc.client.ListTagsForResource(tc.ctx, &sns.ListTagsForResourceInput{
			ResourceArn: aws.String(topicArn),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		for _, t := range listResp.Tags {
			if *t.Key == "RemoveMe" {
				return fmt.Errorf("tag RemoveMe should have been removed")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "TagResource_MultipleTags", func() error {
		multiTagTopicName := tc.uniqueName("MultiTagTopic")
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(multiTagTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*resp.TopicArn)

		_, err = tc.client.TagResource(tc.ctx, &sns.TagResourceInput{
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

		listResp, err := tc.client.ListTagsForResource(tc.ctx, &sns.ListTagsForResourceInput{
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

	return results
}
