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

	results = append(results, r.RunTest("sqs", "DLQ_MaxReceiveCount_DeliveredExactlyN", func() error {
		// Verify that a message is delivered exactly maxReceiveCount times
		// before being moved to the DLQ (not maxReceiveCount-1).
		// Regression test for off-by-one bug where >= was used instead of >.

		dlqName := fmt.Sprintf("DLQRecv-%d", time.Now().UnixNano())
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

		dlqUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(dlqName),
		})
		if err != nil {
			return fmt.Errorf("get dlq url: %v", err)
		}

		dlqAttrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       dlqUrlResp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
		})
		if err != nil {
			return fmt.Errorf("get dlq arn: %v", err)
		}
		dlqArn := dlqAttrs.Attributes[string(types.QueueAttributeNameQueueArn)]

		srcName := fmt.Sprintf("SrcRecv-%d", time.Now().UnixNano())
		redrivePolicy := fmt.Sprintf(`{"deadLetterTargetArn":"%s","maxReceiveCount":"2"}`, dlqArn)
		_, err = client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(srcName),
			Attributes: map[string]string{
				"RedrivePolicy":     redrivePolicy,
				"VisibilityTimeout": "1",
			},
		})
		if err != nil {
			return fmt.Errorf("create src: %v", err)
		}
		defer func() {
			urlResp, _ := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(srcName)})
			if urlResp.QueueUrl != nil {
				client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: urlResp.QueueUrl})
			}
		}()

		srcUrlResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(srcName),
		})
		if err != nil {
			return fmt.Errorf("get src url: %v", err)
		}
		srcURL := *srcUrlResp.QueueUrl

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(srcURL),
			MessageBody: aws.String("DLQ off-by-one regression test"),
		})
		if err != nil {
			return fmt.Errorf("send message: %v", err)
		}

		// 1st receive — should succeed (count=1, 1>2 is false)
		recv1, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(srcURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     2,
		})
		if err != nil {
			return fmt.Errorf("receive 1: %v", err)
		}
		if len(recv1.Messages) != 1 {
			return fmt.Errorf("receive 1: expected 1 message, got %d", len(recv1.Messages))
		}

		time.Sleep(2 * time.Second)

		// 2nd receive — should STILL succeed (count=2, 2>2 is false)
		recv2, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(srcURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     2,
		})
		if err != nil {
			return fmt.Errorf("receive 2: %v", err)
		}
		if len(recv2.Messages) != 1 {
			return fmt.Errorf("receive 2: expected 1 message (maxReceiveCount=2), got %d — off-by-one bug", len(recv2.Messages))
		}

		time.Sleep(2 * time.Second)

		// 3rd receive — source should be empty (count=3, 3>2 is true, moved to DLQ)
		recv3, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(srcURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     2,
		})
		if err != nil {
			return fmt.Errorf("receive 3: %v", err)
		}
		if len(recv3.Messages) != 0 {
			return fmt.Errorf("receive 3: expected 0 messages in source (should be moved to DLQ), got %d", len(recv3.Messages))
		}

		// DLQ should have the message
		dlqRecv, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            dlqUrlResp.QueueUrl,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     2,
		})
		if err != nil {
			return fmt.Errorf("receive from DLQ: %v", err)
		}
		if len(dlqRecv.Messages) != 1 {
			return fmt.Errorf("DLQ: expected 1 message, got %d", len(dlqRecv.Messages))
		}

		return nil
	}))

	return results
}
