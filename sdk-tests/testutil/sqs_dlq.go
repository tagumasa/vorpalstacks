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

// setupMoveTaskQueues creates the dead-letter queue, its redrive-source
// feeder and the destination queue a move task needs, resolving the two ARNs
// and returning a cleanup deleting all three queues. The feeder exists
// because the SourceArn contract accepts "only ARNs of dead-letter queues
// (DLQs) whose sources are other Amazon SQS queues".
func setupMoveTaskQueues(ctx context.Context, client *sqs.Client, prefix string) (srcURL, destURL *string, srcArn, destArn string, cleanup func(), err error) {
	ts := time.Now().UnixNano()
	srcURL, srcCleanup, err := createTestQueue(ctx, client, fmt.Sprintf("%sDLQ-%d", prefix, ts), nil)
	if err != nil {
		return nil, nil, "", "", nil, err
	}
	destURL, destCleanup, err := createTestQueue(ctx, client, fmt.Sprintf("%sDest-%d", prefix, ts), nil)
	if err != nil {
		srcCleanup()
		return nil, nil, "", "", nil, err
	}
	feederURL, feederCleanup, err := createTestQueue(ctx, client, fmt.Sprintf("%sFeeder-%d", prefix, ts), nil)
	if err != nil {
		destCleanup()
		srcCleanup()
		return nil, nil, "", "", nil, err
	}
	cleanup = func() {
		feederCleanup()
		destCleanup()
		srcCleanup()
	}
	srcArn, err = queueArn(ctx, client, srcURL)
	if err != nil {
		cleanup()
		return nil, nil, "", "", nil, err
	}
	destArn, err = queueArn(ctx, client, destURL)
	if err != nil {
		cleanup()
		return nil, nil, "", "", nil, err
	}
	if _, err := client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
		QueueUrl: feederURL,
		Attributes: map[string]string{
			"RedrivePolicy": fmt.Sprintf(`{"deadLetterTargetArn":%q,"maxReceiveCount":3}`, srcArn),
		},
	}); err != nil {
		cleanup()
		return nil, nil, "", "", nil, fmt.Errorf("arm feeder redrive policy: %v", err)
	}
	return srcURL, destURL, srcArn, destArn, cleanup, nil
}

// setupMessageMoveTask creates the move-task queues (see setupMoveTaskQueues)
// and starts a task from the DLQ to the destination. It returns the DLQ and
// destination ARNs and the task handle together with the cleanup function.
func setupMessageMoveTask(ctx context.Context, client *sqs.Client, prefix string) (string, string, *string, func(), error) {
	_, _, srcArn, destArn, cleanup, err := setupMoveTaskQueues(ctx, client, prefix)
	if err != nil {
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
		// A message held in flight keeps the task RUNNING (the worker waits
		// out its visibility window), so the cancel deterministically lands
		// on a live task instead of racing an empty source's instant
		// completion.
		srcURL, _, srcArn, _, cleanup, err := setupMoveTaskQueues(ctx, client, "Cancel")
		if err != nil {
			return err
		}
		defer cleanup()

		if _, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    srcURL,
			MessageBody: aws.String("held in flight across the cancel"),
		}); err != nil {
			return fmt.Errorf("seed DLQ message: %v", err)
		}
		recv, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            srcURL,
			MaxNumberOfMessages: 1,
			VisibilityTimeout:   60,
		})
		if err != nil || len(recv.Messages) != 1 {
			return fmt.Errorf("hold DLQ message in flight: %v (%d messages)", err, len(recv.Messages))
		}

		taskResp, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn: aws.String(srcArn),
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
		if cancelResp.ApproximateNumberOfMessagesMoved != 0 {
			return fmt.Errorf("cancel reported %d messages moved; the only DLQ message is in flight and must not have moved", cancelResp.ApproximateNumberOfMessagesMoved)
		}

		// The cancel must be carried out, not just recorded: the task
		// reaches CANCELLED well inside the held visibility window.
		deadline := time.Now().Add(10 * time.Second)
		for {
			listResp, err := client.ListMessageMoveTasks(ctx, &sqs.ListMessageMoveTasksInput{
				SourceArn: aws.String(srcArn),
			})
			if err != nil {
				return fmt.Errorf("list tasks after cancel: %v", err)
			}
			if len(listResp.Results) == 0 {
				return fmt.Errorf("ListMessageMoveTasks returned empty Results after cancel")
			}
			status := aws.ToString(listResp.Results[0].Status)
			if status == "CANCELLED" {
				return nil
			}
			if status != "CANCELLING" && status != "RUNNING" {
				return fmt.Errorf("cancelled task ended %s, want CANCELLED", status)
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("task still %s 10s after the cancel was accepted", status)
			}
			time.Sleep(100 * time.Millisecond)
		}
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
		if len(listResp.Results) == 0 {
			return fmt.Errorf("ListMessageMoveTasks returned empty Results; the unset-rate task record is missing")
		}
		for _, task := range listResp.Results {
			if task.MaxNumberOfMessagesPerSecond != nil {
				return fmt.Errorf("unset move rate reported as fixed rate %d", *task.MaxNumberOfMessagesPerSecond)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "StartMessageMoveTask_OmittedDestination_RoutesPerMessage", func() error {
		// An omitted DestinationArn redrives every message to its own
		// original source queue ("the messages will be redriven back to their
		// respective original source queues") and the task entries report no
		// destination ("this field value will be NULL").
		ts := time.Now().UnixNano()
		dlqURL, cleanupDLQ, err := createTestQueue(ctx, client, fmt.Sprintf("PerMsgDLQ-%d", ts), nil)
		if err != nil {
			return err
		}
		defer cleanupDLQ()
		dlqArn, err := queueArn(ctx, client, dlqURL)
		if err != nil {
			return err
		}

		makeSource := func(name string, body string) (string, func(), error) {
			srcURL, srcCleanup, err := createTestQueue(ctx, client, name, map[string]string{
				"RedrivePolicy":     fmt.Sprintf(`{"deadLetterTargetArn":"%s","maxReceiveCount":"1"}`, dlqArn),
				"VisibilityTimeout": "0",
			})
			if err != nil {
				return "", nil, err
			}
			if _, err := client.SendMessage(ctx, &sqs.SendMessageInput{
				QueueUrl:    srcURL,
				MessageBody: aws.String(body),
			}); err != nil {
				srcCleanup()
				return "", nil, fmt.Errorf("send to %s: %v", name, err)
			}
			// Two receives with zero visibility exceed maxReceiveCount=1 and
			// redrive the message into the DLQ synchronously.
			for i := 0; i < 2; i++ {
				if _, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: srcURL}); err != nil {
					srcCleanup()
					return "", nil, fmt.Errorf("receive %d on %s: %v", i, name, err)
				}
			}
			return *srcURL, srcCleanup, nil
		}
		srcAURL, cleanupA, err := makeSource(fmt.Sprintf("PerMsgA-%d", ts), "grew in source a")
		if err != nil {
			return err
		}
		defer cleanupA()
		srcBURL, cleanupB, err := makeSource(fmt.Sprintf("PerMsgB-%d", ts), "grew in source b")
		if err != nil {
			return err
		}
		defer cleanupB()

		// The task must start while the sources still carry their redrive
		// policies — they are what proves the source is a dead-letter queue.
		if _, err := client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn: aws.String(dlqArn),
		}); err != nil {
			return fmt.Errorf("start omitted-destination task: %v", err)
		}

		// Wait for the task to complete, then verify no listed entry carries a
		// destination.
		deadline := time.Now().Add(10 * time.Second)
		completed := false
		for time.Now().Before(deadline) {
			listResp, err := client.ListMessageMoveTasks(ctx, &sqs.ListMessageMoveTasksInput{
				SourceArn: aws.String(dlqArn),
			})
			if err != nil {
				return fmt.Errorf("list tasks: %v", err)
			}
			for _, task := range listResp.Results {
				if task.Status != nil && *task.Status == "COMPLETED" {
					completed = true
				}
				if task.DestinationArn != nil {
					return fmt.Errorf("omitted destination reported as %q", *task.DestinationArn)
				}
			}
			if completed {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if !completed {
			return fmt.Errorf("omitted-destination move task did not complete within the deadline")
		}

		// The moved copies keep their accumulated receive counts, so the
		// assertion receives would re-redrive instantly on an armed source;
		// disarm now that the task is done.
		for _, srcURL := range []string{srcAURL, srcBURL} {
			if _, err := client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
				QueueUrl:   aws.String(srcURL),
				Attributes: map[string]string{"RedrivePolicy": ""},
			}); err != nil {
				return fmt.Errorf("disarm source: %v", err)
			}
		}

		recvBody := func(queueURL string) (string, error) {
			recv, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:        aws.String(queueURL),
				WaitTimeSeconds: 2,
			})
			if err != nil {
				return "", err
			}
			if len(recv.Messages) != 1 {
				return "", fmt.Errorf("expected 1 message, got %d", len(recv.Messages))
			}
			return *recv.Messages[0].Body, nil
		}
		bodyA, err := recvBody(srcAURL)
		if err != nil {
			return fmt.Errorf("receive from source a: %v", err)
		}
		if bodyA != "grew in source a" {
			return fmt.Errorf("source a received %q; its own redriven message went elsewhere", bodyA)
		}
		bodyB, err := recvBody(srcBURL)
		if err != nil {
			return fmt.Errorf("receive from source b: %v", err)
		}
		if bodyB != "grew in source b" {
			return fmt.Errorf("source b received %q; its own redriven message went elsewhere", bodyB)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "StartMessageMoveTask_MalformedDestination_Rejected", func() error {
		// A DestinationArn that does not parse is a malformed identifier, not
		// an omitted member: it must reject instead of silently falling back
		// to the redrive homes.
		ts := time.Now().UnixNano()
		dlqURL, cleanupDLQ, err := createTestQueue(ctx, client, fmt.Sprintf("MalDLQ-%d", ts), nil)
		if err != nil {
			return err
		}
		defer cleanupDLQ()
		dlqArn, err := queueArn(ctx, client, dlqURL)
		if err != nil {
			return err
		}
		feederURL, cleanupFeeder, err := createTestQueue(ctx, client, fmt.Sprintf("MalFeeder-%d", ts), map[string]string{
			"RedrivePolicy": fmt.Sprintf(`{"deadLetterTargetArn":"%s","maxReceiveCount":"3"}`, dlqArn),
		})
		if err != nil {
			return err
		}
		defer cleanupFeeder()
		_ = feederURL

		_, err = client.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
			SourceArn:      aws.String(dlqArn),
			DestinationArn: aws.String("not-an-arn"),
		})
		if err == nil {
			return fmt.Errorf("malformed DestinationArn accepted")
		}
		if !strings.Contains(err.Error(), "InvalidAddress") {
			return fmt.Errorf("expected InvalidAddress for a malformed DestinationArn, got: %v", err)
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
