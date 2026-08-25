package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func (r *TestRunner) runSQSPermissionTests(ctx context.Context, client *sqs.Client, queueName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sqs", "TagQueue_ReservedPrefixRejected", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.TagQueue(ctx, &sqs.TagQueueInput{
			QueueUrl: resp.QueueUrl,
			Tags: map[string]string{
				"aws:reserved": "v",
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for aws:-prefixed tag key")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "TagQueue_TooManyTags", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		tags := make(map[string]string, 51)
		for i := 0; i < 51; i++ {
			tags[fmt.Sprintf("Key%d", i)] = "v"
		}
		_, err = client.TagQueue(ctx, &sqs.TagQueueInput{
			QueueUrl: resp.QueueUrl,
			Tags:     tags,
		})
		if err == nil {
			return fmt.Errorf("expected error for 51 tags (quota is 50)")
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
		if err != nil {
			return err
		}
		tagResp, err := client.ListQueueTags(ctx, &sqs.ListQueueTagsInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return fmt.Errorf("list tags after TagQueue: %v", err)
		}
		if tagResp.Tags["Environment"] != "Test" {
			return fmt.Errorf("tag not found after TagQueue: expected Environment=Test, got %v", tagResp.Tags)
		}
		return nil
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
		if err != nil {
			return err
		}
		tagResp, err := client.ListQueueTags(ctx, &sqs.ListQueueTagsInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return fmt.Errorf("list tags after UntagQueue: %v", err)
		}
		if _, ok := tagResp.Tags["Environment"]; ok {
			return fmt.Errorf("tag Environment should be removed after UntagQueue")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "AddPermission_AnyActionName_Accepted", func() error {
		// The Actions parameter documents "the name of any action or `*`";
		// names outside any fixed list must be accepted.
		anyQURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("AnyActionQueue-%d", time.Now().UnixNano()), nil)
		if err != nil {
			return err
		}
		defer cleanup()

		_, err = client.AddPermission(ctx, &sqs.AddPermissionInput{
			QueueUrl:      anyQURL,
			Label:         aws.String("AnyActionPerm"),
			AWSAccountIds: []string{r.accountID},
			Actions:       []string{"CreateQueue", "ListQueues"},
		})
		if err != nil {
			return fmt.Errorf("AddPermission with queue-management action names rejected: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "AddPermission", func() error {
		permQueueURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("PermQueue-%d", time.Now().UnixNano()), nil)
		if err != nil {
			return err
		}
		defer cleanup()

		_, err = client.AddPermission(ctx, &sqs.AddPermissionInput{
			QueueUrl:      permQueueURL,
			Label:         aws.String("TestPermission"),
			AWSAccountIds: []string{r.accountID},
			Actions:       []string{"SendMessage", "ReceiveMessage"},
		})
		if err != nil {
			return err
		}
		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       permQueueURL,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNamePolicy},
		})
		if err != nil {
			return fmt.Errorf("get queue policy after AddPermission: %v", err)
		}
		policy := attrResp.Attributes[string(types.QueueAttributeNamePolicy)]
		if policy == "" {
			return fmt.Errorf("queue Policy attribute is empty after AddPermission")
		}
		if !strings.Contains(policy, "TestPermission") {
			return fmt.Errorf("queue Policy does not contain label 'TestPermission': %s", policy)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "RemovePermission", func() error {
		permQueueURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("RPermQueue-%d", time.Now().UnixNano()), nil)
		if err != nil {
			return err
		}
		defer cleanup()

		_, err = client.AddPermission(ctx, &sqs.AddPermissionInput{
			QueueUrl:      permQueueURL,
			Label:         aws.String("RemoveTestPerm"),
			AWSAccountIds: []string{r.accountID},
			Actions:       []string{"SendMessage"},
		})
		if err != nil {
			return fmt.Errorf("add permission: %v", err)
		}

		_, err = client.RemovePermission(ctx, &sqs.RemovePermissionInput{
			QueueUrl: permQueueURL,
			Label:    aws.String("RemoveTestPerm"),
		})
		if err != nil {
			return err
		}
		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       permQueueURL,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNamePolicy},
		})
		if err != nil {
			return fmt.Errorf("get queue policy after RemovePermission: %v", err)
		}
		policy := attrResp.Attributes[string(types.QueueAttributeNamePolicy)]
		if policy != "" && strings.Contains(policy, "RemoveTestPerm") {
			return fmt.Errorf("queue Policy still contains label 'RemoveTestPerm' after RemovePermission")
		}
		return nil
	}))

	return results
}
