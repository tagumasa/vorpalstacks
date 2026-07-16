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

func (r *TestRunner) runSQSQueueTests(ctx context.Context, client *sqs.Client, queueName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sqs", "ListQueues", func() error {
		resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{})
		if err != nil {
			return err
		}
		if resp.QueueUrls == nil {
			return fmt.Errorf("ListQueues returned nil QueueUrls")
		}
		found := false
		for _, u := range resp.QueueUrls {
			if strings.Contains(u, queueName) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created queue %s not found in ListQueues", queueName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListQueues_WithPrefix", func() error {
		prefix := fmt.Sprintf("PrefixTest-%d", time.Now().UnixNano())
		client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(prefix + "-alpha"),
		})
		client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(prefix + "-beta"),
		})
		defer func() {
			for _, suffix := range []string{"-alpha", "-beta"} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(prefix + suffix)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{
			QueueNamePrefix: aws.String(prefix),
		})
		if err != nil {
			return err
		}
		if len(resp.QueueUrls) < 2 {
			return fmt.Errorf("expected at least 2 queues with prefix %q, got %d", prefix, len(resp.QueueUrls))
		}
		for _, u := range resp.QueueUrls {
			if !strings.Contains(u, prefix) {
				return fmt.Errorf("ListQueues returned URL not matching prefix: %s", u)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SetQueueAttributes", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
			QueueUrl: resp.QueueUrl,
			Attributes: map[string]string{
				"VisibilityTimeout": "30",
			},
		})
		if err != nil {
			return err
		}
		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       resp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameVisibilityTimeout},
		})
		if err != nil {
			return fmt.Errorf("get attrs after set: %v", err)
		}
		if attrResp.Attributes[string(types.QueueAttributeNameVisibilityTimeout)] != "30" {
			return fmt.Errorf("VisibilityTimeout not updated: got %s", attrResp.Attributes["VisibilityTimeout"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "SetQueueAttributes_MultipleAttrs", func() error {
		rtQueueName := fmt.Sprintf("SMAQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(rtQueueName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(rtQueueName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(rtQueueName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		_, err = client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
			QueueUrl: urlResp.QueueUrl,
			Attributes: map[string]string{
				"VisibilityTimeout":      "45",
				"MaximumMessageSize":     "1024",
				"MessageRetentionPeriod": "345600",
				"DelaySeconds":           "10",
			},
		})
		if err != nil {
			return fmt.Errorf("set attrs: %v", err)
		}

		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: urlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{
				types.QueueAttributeNameVisibilityTimeout,
				types.QueueAttributeNameMaximumMessageSize,
				types.QueueAttributeNameMessageRetentionPeriod,
				types.QueueAttributeNameDelaySeconds,
			},
		})
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if attrResp.Attributes[string(types.QueueAttributeNameVisibilityTimeout)] != "45" {
			return fmt.Errorf("VisibilityTimeout mismatch: got %s", attrResp.Attributes["VisibilityTimeout"])
		}
		if attrResp.Attributes[string(types.QueueAttributeNameMaximumMessageSize)] != "1024" {
			return fmt.Errorf("MaximumMessageSize mismatch: got %s", attrResp.Attributes["MaximumMessageSize"])
		}
		if attrResp.Attributes[string(types.QueueAttributeNameMessageRetentionPeriod)] != "345600" {
			return fmt.Errorf("MessageRetentionPeriod mismatch: got %s", attrResp.Attributes["MessageRetentionPeriod"])
		}
		if attrResp.Attributes[string(types.QueueAttributeNameDelaySeconds)] != "10" {
			return fmt.Errorf("DelaySeconds mismatch: got %s", attrResp.Attributes["DelaySeconds"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "PurgeQueue", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return err
		}
		recvResp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:        resp.QueueUrl,
			WaitTimeSeconds: 1,
		})
		if err != nil {
			return fmt.Errorf("receive after purge: %v", err)
		}
		if len(recvResp.Messages) > 0 {
			return fmt.Errorf("expected 0 messages after PurgeQueue, got %d", len(recvResp.Messages))
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "CreateQueue (FIFO)", func() error {
		fifoQueueName := fmt.Sprintf("TestFifoQueue-%d.fifo", time.Now().UnixNano())
		resp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(fifoQueueName),
			Attributes: map[string]string{
				"ContentBasedDeduplication": "true",
				"FifoQueue":                 "true",
			},
		})
		if err != nil {
			return err
		}
		if resp.QueueUrl == nil {
			return fmt.Errorf("CreateQueue (FIFO) returned nil QueueUrl")
		}
		client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: resp.QueueUrl})
		return nil
	}))

	results = append(results, r.RunTest("sqs", "CreateQueue_DuplicateName", func() error {
		dupQName := fmt.Sprintf("DupQueue-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(dupQName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(dupQName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(dupQName),
		})
		if err != nil {
			return fmt.Errorf("duplicate queue name should be idempotent, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListQueues_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgQueues []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagQ-%s-%d", pgTs, i)
			_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
				QueueName: aws.String(name),
			})
			if err != nil {
				for _, qn := range pgQueues {
					uResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qn)})
					if uResp.QueueUrl != nil {
						client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: uResp.QueueUrl})
					}
				}
				return fmt.Errorf("create queue %s: %v", name, err)
			}
			pgQueues = append(pgQueues, name)
		}

		var allQueues []string
		var nextToken *string
		for {
			resp, err := client.ListQueues(ctx, &sqs.ListQueuesInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, qn := range pgQueues {
					uResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qn)})
					if uResp.QueueUrl != nil {
						client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: uResp.QueueUrl})
					}
				}
				return fmt.Errorf("list queues page: %v", err)
			}
			for _, u := range resp.QueueUrls {
				if strings.Contains(u, "PagQ-"+pgTs) {
					allQueues = append(allQueues, u)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, qn := range pgQueues {
			uResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(qn)})
			if uResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: uResp.QueueUrl})
			}
		}
		if len(allQueues) != 5 {
			return fmt.Errorf("expected 5 paginated queues, got %d", len(allQueues))
		}
		return nil
	}))

	return results
}
