package testutil

import (
	"context"
	"fmt"
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
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("sns", "GetTopicAttributes_NonExistent", func() error {
		_, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent topic")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Publish_NonExistentTopic", func() error {
		_, err := client.Publish(ctx, &sns.PublishInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:nonexistent-topic-xyz"),
			Message:  aws.String("test message"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent topic")
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

	return results
}
