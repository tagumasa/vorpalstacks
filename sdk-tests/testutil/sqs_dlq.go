package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func (r *TestRunner) runSQSDLQTests(ctx context.Context, client *sqs.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sqs", "ListDeadLetterSourceQueues_Empty", func() error {
		dlqName := fmt.Sprintf("DLQ-%d", time.Now().UnixNano())
		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(dlqName),
		})
		if err != nil {
			return fmt.Errorf("create dlq: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(dlqName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		urlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(dlqName),
		})
		if err != nil {
			return fmt.Errorf("get url: %v", err)
		}

		dlqResp, err := client.ListDeadLetterSourceQueues(ctx, &sqs.ListDeadLetterSourceQueuesInput{
			QueueUrl: urlResp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if len(dlqResp.QueueUrls) != 0 {
			return fmt.Errorf("expected empty queue URLs for new DLQ, got %d", len(dlqResp.QueueUrls))
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "StartMessageMoveTask", func() error {
		srcDlqName := fmt.Sprintf("SrcDLQ-%d", time.Now().UnixNano())
		destQueueName := fmt.Sprintf("DestQueue-%d", time.Now().UnixNano())

		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("create src: %v", err)
		}
		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("create dest: %v", err)
		}
		defer func() {
			for _, name := range []string{srcDlqName, destQueueName} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(name)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		srcUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("get src url: %v", err)
		}
		destUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("get dest url: %v", err)
		}

		srcAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       srcUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get src attrs: %v", err)
		}
		destAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       destUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get dest attrs: %v", err)
		}

		srcArn := srcAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]
		destArn := destAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]

		taskResp, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:      aws.String(srcArn),
			DestinationArn: aws.String(destArn),
		})
		if err != nil {
			return fmt.Errorf("start task: %v", err)
		}
		if taskResp.TaskHandle == nil || *taskResp.TaskHandle == "" {
			return fmt.Errorf("StartMessageMoveTask returned nil or empty TaskHandle")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "CancelMessageMoveTask", func() error {
		srcDlqName := fmt.Sprintf("CancelDLQ-%d", time.Now().UnixNano())
		destQueueName := fmt.Sprintf("CancelDest-%d", time.Now().UnixNano())

		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("create src: %v", err)
		}
		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("create dest: %v", err)
		}
		defer func() {
			for _, name := range []string{srcDlqName, destQueueName} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(name)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		srcUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("get src url: %v", err)
		}
		destUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("get dest url: %v", err)
		}

		srcAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       srcUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get src attrs: %v", err)
		}
		destAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       destUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get dest attrs: %v", err)
		}

		srcArn := srcAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]
		destArn := destAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]

		taskResp, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:      aws.String(srcArn),
			DestinationArn: aws.String(destArn),
		})
		if err != nil {
			return fmt.Errorf("start task: %v", err)
		}

		cancelResp, err := client.CancelMessageMoveTask(ctx, &sqs.CancelMessageMoveTaskInput{
			TaskHandle: taskResp.TaskHandle,
		})
		if err != nil {
			return fmt.Errorf("cancel task: %v", err)
		}
		if cancelResp.ApproximateNumberOfMessagesMoved < 0 {
			return fmt.Errorf("CancelMessageMoveTask returned negative ApproximateNumberOfMessagesMoved: %d", cancelResp.ApproximateNumberOfMessagesMoved)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListMessageMoveTasks", func() error {
		srcDlqName := fmt.Sprintf("ListDLQ-%d", time.Now().UnixNano())
		destQueueName := fmt.Sprintf("ListDest-%d", time.Now().UnixNano())

		_, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("create src: %v", err)
		}
		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("create dest: %v", err)
		}
		defer func() {
			for _, name := range []string{srcDlqName, destQueueName} {
				urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(name)})
				if urlResp.QueueUrl != nil {
					client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
				}
			}
		}()

		srcUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(srcDlqName),
		})
		if err != nil {
			return fmt.Errorf("get src url: %v", err)
		}
		destUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(destQueueName),
		})
		if err != nil {
			return fmt.Errorf("get dest url: %v", err)
		}

		srcAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       srcUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get src attrs: %v", err)
		}
		destAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       destUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get dest attrs: %v", err)
		}

		srcArn := srcAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]
		destArn := destAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]

		_, err = client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:      aws.String(srcArn),
			DestinationArn: aws.String(destArn),
		})
		if err != nil {
			return fmt.Errorf("start task: %v", err)
		}

		listResp, err := client.ListMessageMoveTasks(ctx, &sqs.ListMessageMoveTasksInput{
			SourceArn: aws.String(srcArn),
		})
		if err != nil {
			return fmt.Errorf("list tasks: %v", err)
		}
		if len(listResp.Results) == 0 {
			return fmt.Errorf("ListMessageMoveTasks returned empty Results")
		}
		task := listResp.Results[0]
		if task.TaskHandle == nil || *task.TaskHandle == "" {
			return fmt.Errorf("ListMessageMoveTasks result has nil or empty TaskHandle")
		}
		if task.SourceArn == nil || *task.SourceArn == "" {
			return fmt.Errorf("ListMessageMoveTasks result has nil or empty SourceArn")
		}
		if task.Status == nil || *task.Status == "" {
			return fmt.Errorf("ListMessageMoveTasks result has nil or empty Status")
		}
		return nil
	}))

	return results
}
