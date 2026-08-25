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

// setupMessageMoveTask creates a source dead-letter queue and a destination
// queue named after prefix, resolves both ARNs and starts a message-move task
// from the source to the destination. It returns the two ARNs and the task
// handle together with a cleanup function deleting both queues.
func setupMessageMoveTask(ctx context.Context, client *sqs.Client, prefix string) (string, string, *string, func(), error) {
	ts := time.Now().UnixNano()
	srcURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("%sDLQ-%d", prefix, ts), nil)
	if err != nil {
		return "", "", nil, nil, err
	}
	destURL, _, err := createTestQueue(ctx, client, fmt.Sprintf("%sDest-%d", prefix, ts), nil)
	if err != nil {
		cleanup()
		return "", "", nil, nil, err
	}
	srcArn, err := queueArn(ctx, client, srcURL)
	if err != nil {
		cleanup()
		return "", "", nil, nil, err
	}
	destArn, err := queueArn(ctx, client, destURL)
	if err != nil {
		cleanup()
		return "", "", nil, nil, err
	}
	taskResp, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
		SourceArn:      aws.String(srcArn),
		DestinationArn: aws.String(destArn),
	})
	if err != nil {
		cleanup()
		return "", "", nil, nil, fmt.Errorf("start task: %v", err)
	}
	return srcArn, destArn, taskResp.TaskHandle, cleanup, nil
}

func (r *TestRunner) runSQSDLQTests(ctx context.Context, client *sqs.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sqs", "ListDeadLetterSourceQueues_Empty", func() error {
		dlqURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("DLQ-%d", time.Now().UnixNano()), nil)
		if err != nil {
			return err
		}
		defer cleanup()

		dlqResp, err := client.ListDeadLetterSourceQueues(ctx, &sqs.ListDeadLetterSourceQueuesInput{
			QueueUrl: dlqURL,
		})
		if err != nil {
			return err
		}
		if len(dlqResp.QueueUrls) != 0 {
			return fmt.Errorf("expected empty queue URLs for new DLQ, got %d", len(dlqResp.QueueUrls))
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListDeadLetterSourceQueues_NonEmpty", func() error {
		// The SDK parses the response member as lowercase "queueUrls"; a
		// PascalCase key would silently yield an empty list.
		dlqURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("DLQSrc-%d", time.Now().UnixNano()), nil)
		if err != nil {
			return err
		}
		defer cleanup()
		dlqArn, err := queueArn(ctx, client, dlqURL)
		if err != nil {
			return err
		}

		srcURL, cleanupSrc, err := createTestQueue(ctx, client, fmt.Sprintf("SrcListing-%d", time.Now().UnixNano()), map[string]string{
			"RedrivePolicy": fmt.Sprintf(`{"deadLetterTargetArn":"%s","maxReceiveCount":"3"}`, dlqArn),
		})
		if err != nil {
			cleanup()
			return err
		}
		defer cleanupSrc()

		dlqResp, err := client.ListDeadLetterSourceQueues(ctx, &sqs.ListDeadLetterSourceQueuesInput{
			QueueUrl: dlqURL,
		})
		if err != nil {
			return err
		}
		found := false
		for _, u := range dlqResp.QueueUrls {
			if u == *srcURL {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected source queue URL %s in ListDeadLetterSourceQueues result, got %v", *srcURL, dlqResp.QueueUrls)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DLQ_ReceiveReturnsSourceArn", func() error {
		// A message moved to a dead-letter queue must carry the
		// DeadLetterQueueSourceArn system attribute naming the source queue.
		dlqURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("DLQArn-%d", time.Now().UnixNano()), nil)
		if err != nil {
			return err
		}
		defer cleanup()
		dlqArn, err := queueArn(ctx, client, dlqURL)
		if err != nil {
			return err
		}

		srcURL, cleanupSrc, err := createTestQueue(ctx, client, fmt.Sprintf("SrcArn-%d", time.Now().UnixNano()), map[string]string{
			"RedrivePolicy":     fmt.Sprintf(`{"deadLetterTargetArn":"%s","maxReceiveCount":"1"}`, dlqArn),
			"VisibilityTimeout": "0",
		})
		if err != nil {
			cleanup()
			return err
		}
		defer cleanupSrc()
		srcArn, err := queueArn(ctx, client, srcURL)
		if err != nil {
			return err
		}

		if _, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    srcURL,
			MessageBody: aws.String("message that will be redriven"),
		}); err != nil {
			return fmt.Errorf("send: %v", err)
		}

		// VisibilityTimeout=0 lets each receive redeliver immediately; the
		// second receive therefore exceeds maxReceiveCount=1 and moves the
		// message to the DLQ without any wait between receives.
		if _, err = client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: srcURL}); err != nil {
			return fmt.Errorf("receive 1: %v", err)
		}
		if _, err = client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: srcURL}); err != nil {
			return fmt.Errorf("receive 2: %v", err)
		}

		dlqRecv, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:                    dlqURL,
			MessageSystemAttributeNames: []types.MessageSystemAttributeName{types.MessageSystemAttributeNameAll},
		})
		if err != nil {
			return fmt.Errorf("receive from dlq: %v", err)
		}
		if len(dlqRecv.Messages) != 1 {
			return fmt.Errorf("expected 1 message in DLQ, got %d", len(dlqRecv.Messages))
		}
		got := dlqRecv.Messages[0].Attributes["DeadLetterQueueSourceArn"]
		if got != srcArn {
			return fmt.Errorf("expected DeadLetterQueueSourceArn %q, got %q", srcArn, got)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "StartMessageMoveTask", func() error {
		_, _, taskHandle, cleanup, err := setupMessageMoveTask(ctx, client, "Src")
		if err != nil {
			return err
		}
		defer cleanup()

		if taskHandle == nil || *taskHandle == "" {
			return fmt.Errorf("StartMessageMoveTask returned nil or empty TaskHandle")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "CancelMessageMoveTask", func() error {
		_, _, taskHandle, cleanup, err := setupMessageMoveTask(ctx, client, "Cancel")
		if err != nil {
			return err
		}
		defer cleanup()

		cancelResp, err := client.CancelMessageMoveTask(ctx, &sqs.CancelMessageMoveTaskInput{
			TaskHandle: taskHandle,
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
		srcArn, _, _, cleanup, err := setupMessageMoveTask(ctx, client, "List")
		if err != nil {
			return err
		}
		defer cleanup()

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

	results = append(results, r.RunTest("sqs", "CancelMessageMoveTask_UnknownHandle_ResourceNotFound", func() error {
		_, err := client.CancelMessageMoveTask(ctx, &sqs.CancelMessageMoveTaskInput{
			TaskHandle: aws.String("nonexistent-task-handle"),
		})
		if err == nil {
			return fmt.Errorf("unknown TaskHandle must fail with ResourceNotFoundException")
		}
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return fmt.Errorf("expected ResourceNotFoundException, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "ListMessageMoveTasks_UnsetRate_Omitted", func() error {
		srcArn, _, _, cleanup, err := setupMessageMoveTask(ctx, client, "Unrate")
		if err != nil {
			return err
		}
		defer cleanup()
		// The move task is started without MaxNumberOfMessagesPerSecond:
		// system-optimised rate.

		listResp, err := client.ListMessageMoveTasks(ctx, &sqs.ListMessageMoveTasksInput{
			SourceArn: aws.String(srcArn),
		})
		if err != nil {
			return fmt.Errorf("list tasks: %v", err)
		}
		for _, task := range listResp.Results {
			if task.MaxNumberOfMessagesPerSecond != nil {
				return fmt.Errorf("unset move rate reported as fixed rate %d", *task.MaxNumberOfMessagesPerSecond)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "DLQ_MaxReceiveCount_DeliveredExactlyN", func() error {
		// Verify that a message is delivered exactly maxReceiveCount times
		// before being moved to the DLQ (not maxReceiveCount-1).
		// Regression test for off-by-one bug where >= was used instead of >.

		dlqURL, cleanup, err := createTestQueue(ctx, client, fmt.Sprintf("DLQRecv-%d", time.Now().UnixNano()), nil)
		if err != nil {
			return err
		}
		defer cleanup()
		dlqArn, err := queueArn(ctx, client, dlqURL)
		if err != nil {
			return err
		}

		srcName := fmt.Sprintf("SrcRecv-%d", time.Now().UnixNano())
		redrivePolicy := fmt.Sprintf(`{"deadLetterTargetArn":"%s","maxReceiveCount":"2"}`, dlqArn)
		// Zero visibility lets each receive redeliver immediately, so the
		// receive-count progression needs no sleeps between receives.
		rawSrcURL, cleanupSrc, err := createTestQueue(ctx, client, srcName, map[string]string{
			"RedrivePolicy":     redrivePolicy,
			"VisibilityTimeout": "0",
		})
		if err != nil {
			cleanup()
			return err
		}
		defer cleanupSrc()

		srcURL := *rawSrcURL

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

		// 3rd receive — source should be empty (count=3, 3>2 is true, moved
		// to the DLQ synchronously inside the receive scan; the short wait
		// proves the message does not come back).
		recv3, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(srcURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     1,
		})
		if err != nil {
			return fmt.Errorf("receive 3: %v", err)
		}
		if len(recv3.Messages) != 0 {
			return fmt.Errorf("receive 3: expected 0 messages in source (should be moved to DLQ), got %d", len(recv3.Messages))
		}

		// DLQ should have the message
		dlqRecv, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            dlqURL,
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
