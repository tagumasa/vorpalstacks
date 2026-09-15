package sqs

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// TestSendMessageDelaySecondsStrict pins the Integer-member contract of
// SendMessage.DelaySeconds: a present value that is not an integer, or a
// magnitude beyond the int32 domain (the model's Integer members are all
// int32), is a wire-type violation (SerializationException) rejected before
// anything is sent, while an omitted member applies the queue's DelaySeconds
// attribute and a legitimate integer value is honoured.
func TestSendMessageDelaySecondsStrict(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":         "delay-strict",
			"Attribute.1.Name":  "DelaySeconds",
			"Attribute.1.Value": "30",
		},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	// Present non-integer — wire-type violation, nothing is sent.
	_, err = svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":     queueURL,
			"MessageBody":  "garbage-delay",
			"DelaySeconds": "soon",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer DelaySeconds: got %v, want SerializationException", err)
	}
	if visible, _, delayed := store.GetMessageCounts(queueURL); visible != 0 || delayed != 0 {
		t.Fatalf("counts after rejected send = visible %d, delayed %d, want 0/0", visible, delayed)
	}

	// A string magnitude beyond the int32 domain (2^33) is the same
	// wire-type violation — it must reject, never wrap to 0 at the int32
	// conversion and fall back to the queue's delay.
	_, err = svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":     queueURL,
			"MessageBody":  "overflow-delay",
			"DelaySeconds": "8589934592",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("beyond-int32 DelaySeconds: got %v, want SerializationException", err)
	}
	if visible, _, delayed := store.GetMessageCounts(queueURL); visible != 0 || delayed != 0 {
		t.Fatalf("counts after beyond-int32 send = visible %d, delayed %d, want 0/0", visible, delayed)
	}

	// Omitted member — the queue's DelaySeconds attribute (30) applies.
	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":    queueURL,
			"MessageBody": "default-delay",
		},
	}); err != nil {
		t.Fatalf("SendMessage without DelaySeconds: %v", err)
	}
	if visible, _, delayed := store.GetMessageCounts(queueURL); visible != 0 || delayed != 1 {
		t.Fatalf("counts after omitted DelaySeconds = visible %d, delayed %d, want 0/1", visible, delayed)
	}

	// Legitimate integer value — honoured as a per-message delay.
	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":     queueURL,
			"MessageBody":  "explicit-delay",
			"DelaySeconds": "60",
		},
	}); err != nil {
		t.Fatalf("SendMessage with DelaySeconds 60: %v", err)
	}
	if visible, _, delayed := store.GetMessageCounts(queueURL); visible != 0 || delayed != 2 {
		t.Fatalf("counts after explicit DelaySeconds = visible %d, delayed %d, want 0/2", visible, delayed)
	}
}

// TestSendMessageBatchEntryDelaySecondsStrict pins the same Integer-member
// contract on the JSON arm of SendMessageBatch entries: a present non-integer
// is a request-level wire-type violation (nothing is sent), an omitted member
// applies the queue's DelaySeconds attribute, and a JSON number is honoured.
func TestSendMessageBatchEntryDelaySecondsStrict(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":         "delay-batch-strict",
			"Attribute.1.Name":  "DelaySeconds",
			"Attribute.1.Value": "30",
		},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	// Present non-integer in one entry — the whole request is rejected.
	_, err = svc.SendMessageBatch(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessageBatch",
		Parameters: map[string]interface{}{
			"QueueUrl": queueURL,
			"Entries": []interface{}{
				map[string]interface{}{"Id": "ok", "MessageBody": "fine"},
				map[string]interface{}{"Id": "bad", "MessageBody": "garbage", "DelaySeconds": "soon"},
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer entry DelaySeconds: got %v, want SerializationException", err)
	}
	if visible, _, delayed := store.GetMessageCounts(queueURL); visible != 0 || delayed != 0 {
		t.Fatalf("counts after rejected batch = visible %d, delayed %d, want 0/0", visible, delayed)
	}

	// Omitted member — the queue's DelaySeconds attribute (30) applies.
	if _, err := svc.SendMessageBatch(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessageBatch",
		Parameters: map[string]interface{}{
			"QueueUrl": queueURL,
			"Entries": []interface{}{
				map[string]interface{}{"Id": "omit", "MessageBody": "default-delay"},
			},
		},
	}); err != nil {
		t.Fatalf("SendMessageBatch without entry DelaySeconds: %v", err)
	}
	if visible, _, delayed := store.GetMessageCounts(queueURL); visible != 0 || delayed != 1 {
		t.Fatalf("counts after omitted entry DelaySeconds = visible %d, delayed %d, want 0/1", visible, delayed)
	}

	// JSON number — honoured as a per-message delay.
	if _, err := svc.SendMessageBatch(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessageBatch",
		Parameters: map[string]interface{}{
			"QueueUrl": queueURL,
			"Entries": []interface{}{
				map[string]interface{}{"Id": "num", "MessageBody": "explicit-delay", "DelaySeconds": float64(60)},
			},
		},
	}); err != nil {
		t.Fatalf("SendMessageBatch with entry DelaySeconds 60: %v", err)
	}
	if visible, _, delayed := store.GetMessageCounts(queueURL); visible != 0 || delayed != 2 {
		t.Fatalf("counts after explicit entry DelaySeconds = visible %d, delayed %d, want 0/2", visible, delayed)
	}
}

// TestListQueuesMaxResultsStrict pins the Integer-member contract of
// ListQueues.MaxResults: a present non-integer is rejected as a wire-type
// violation, an omitted member keeps the default page (and no NextToken,
// which only a set MaxResults elicits).
func TestListQueuesMaxResultsStrict(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	if _, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "lq-strict"},
	}); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}

	_, err := svc.ListQueues(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ListQueues",
		Parameters: map[string]interface{}{"MaxResults": "many"},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer MaxResults: got %v, want SerializationException", err)
	}

	resp, err := svc.ListQueues(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ListQueues",
		Parameters: map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("ListQueues without MaxResults: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	urls, _ := respMap["QueueUrls"].([]string)
	if len(urls) != 1 {
		t.Fatalf("absent MaxResults: QueueUrls = %v, want the created queue", urls)
	}
	if _, hasToken := respMap["NextToken"]; hasToken {
		t.Fatalf("absent MaxResults must not elicit NextToken: %v", respMap)
	}
}

// TestListDeadLetterSourceQueuesMaxResultsStrict pins the Integer-member
// contract of ListDeadLetterSourceQueues.MaxResults: a present non-integer
// or a magnitude beyond the int32 domain is rejected, an omitted member
// keeps the documented default page.
func TestListDeadLetterSourceQueuesMaxResultsStrict(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "dlq-src-strict"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	_, err = svc.ListDeadLetterSourceQueues(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ListDeadLetterSourceQueues",
		Parameters: map[string]interface{}{
			"QueueUrl":   queueURL,
			"MaxResults": "many",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer MaxResults: got %v, want SerializationException", err)
	}

	// 2^32+1 wraps to 1 at an int32 conversion — it must reject as the same
	// wire-type violation instead of being accepted as a one-item page.
	_, err = svc.ListDeadLetterSourceQueues(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ListDeadLetterSourceQueues",
		Parameters: map[string]interface{}{
			"QueueUrl":   queueURL,
			"MaxResults": "4294967297",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("beyond-int32 MaxResults: got %v, want SerializationException", err)
	}

	if _, err := svc.ListDeadLetterSourceQueues(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ListDeadLetterSourceQueues",
		Parameters: map[string]interface{}{"QueueUrl": queueURL},
	}); err != nil {
		t.Fatalf("ListDeadLetterSourceQueues without MaxResults: %v", err)
	}
}

// TestListMessageMoveTasksMaxResultsStrict pins the Integer-member contract
// of ListMessageMoveTasks.MaxResults: a present non-integer is rejected, an
// omitted member keeps the documented default (the most recent task).
func TestListMessageMoveTasksMaxResultsStrict(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "lmt-strict"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)
	queue, err := store.GetQueue(queueURL)
	if err != nil {
		t.Fatalf("GetQueue: %v", err)
	}

	_, err = svc.ListMessageMoveTasks(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ListMessageMoveTasks",
		Parameters: map[string]interface{}{
			"SourceArn":  queue.ARN,
			"MaxResults": "many",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer MaxResults: got %v, want SerializationException", err)
	}

	if _, err := svc.ListMessageMoveTasks(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ListMessageMoveTasks",
		Parameters: map[string]interface{}{"SourceArn": queue.ARN},
	}); err != nil {
		t.Fatalf("ListMessageMoveTasks without MaxResults: %v", err)
	}
}

// TestStartMessageMoveTaskRateStrict pins the Integer-member contract of
// StartMessageMoveTask.MaxNumberOfMessagesPerSecond on a real redrive pair:
// a present non-integer is rejected, an omitted member starts the task with
// the system-optimised variable rate (the unset 0 the Core forwards).
func TestStartMessageMoveTaskRateStrict(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	dlqResp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "rate-strict-dlq"},
	})
	if err != nil {
		t.Fatalf("CreateQueue DLQ: %v", err)
	}
	dlqMap, _ := dlqResp.(map[string]interface{})
	dlqURL, _ := dlqMap["QueueUrl"].(string)
	dlq, err := store.GetQueue(dlqURL)
	if err != nil {
		t.Fatalf("GetQueue DLQ: %v", err)
	}

	rdp := `{"deadLetterTargetArn":"` + dlq.ARN + `","maxReceiveCount":3}`
	if _, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":         "rate-strict-source",
			"Attribute.1.Name":  "RedrivePolicy",
			"Attribute.1.Value": rdp,
		},
	}); err != nil {
		t.Fatalf("CreateQueue source with RedrivePolicy: %v", err)
	}

	_, err = svc.StartMessageMoveTask(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "StartMessageMoveTask",
		Parameters: map[string]interface{}{
			"SourceArn":                    dlq.ARN,
			"MaxNumberOfMessagesPerSecond": "fast",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer rate: got %v, want SerializationException", err)
	}

	resp, err := svc.StartMessageMoveTask(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "StartMessageMoveTask",
		Parameters: map[string]interface{}{
			"SourceArn": dlq.ARN,
		},
	})
	if err != nil {
		t.Fatalf("StartMessageMoveTask without rate: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	handle, _ := respMap["TaskHandle"].(string)
	if handle == "" {
		t.Fatalf("absent rate must still start the task: %v", respMap)
	}
	// The omitted rate forwards the unset 0 to the store: the started record
	// carries no rate. The wire response cannot show this (it emits
	// MaxNumberOfMessagesPerSecond only when positive), so the pin reads the
	// typed record.
	tasks, lerr := store.ListMessageMoveTasks(dlq.ARN, 10)
	if lerr != nil {
		t.Fatalf("ListMessageMoveTasks: %v", lerr)
	}
	if len(tasks) == 0 {
		t.Fatal("ListMessageMoveTasks returned no task for the just-started source")
	}
	if tasks[0].MaxNumberOfMessages != 0 {
		t.Fatalf("started task rate = %d; the omitted member must forward the unset 0", tasks[0].MaxNumberOfMessages)
	}
	// Best-effort cancel so the move-task goroutine does not outlive the
	// test store (the cleanup stack runs this before store.Close()).
	t.Cleanup(func() {
		_, _ = svc.CancelMessageMoveTask(context.Background(), reqCtx, &request.ParsedRequest{
			Operation: "CancelMessageMoveTask",
			Parameters: map[string]interface{}{
				"TaskHandle": handle,
			},
		})
	})
}

// TestReceiveMessageVisibilityTimeoutRange pins the services-side range
// check of ReceiveMessage.VisibilityTimeout (the queue attribute's
// documented 0–43200 domain): out-of-range values are rejected here like
// its sibling members, while an omitted member applies the queue's
// VisibilityTimeout attribute and an explicit 0 returns the message to
// immediate visibility.
func TestReceiveMessageVisibilityTimeoutRange(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "vt-range"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	for _, bad := range []string{"-1", "43201"} {
		_, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
			Operation: "ReceiveMessage",
			Parameters: map[string]interface{}{
				"QueueUrl":          queueURL,
				"VisibilityTimeout": bad,
			},
		})
		if err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
			t.Fatalf("VisibilityTimeout %s: got %v, want InvalidParameterValue", bad, err)
		}
	}

	// Omitted member — the queue's default VisibilityTimeout (30) applies:
	// the received message is hidden from subsequent receives.
	if _, err := store.SendMessage(queueURL, sqsstore.NewMessage("vt-body")); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	rcvResp, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ReceiveMessage",
		Parameters: map[string]interface{}{
			"QueueUrl": queueURL,
		},
	})
	if err != nil {
		t.Fatalf("ReceiveMessage without VisibilityTimeout: %v", err)
	}
	rcvMap, _ := rcvResp.(map[string]interface{})
	if msgs, _ := rcvMap["Messages"].([]map[string]interface{}); len(msgs) != 1 {
		t.Fatalf("absent VisibilityTimeout: messages = %v, want 1", rcvMap)
	}
	if visible, notVisible, _ := store.GetMessageCounts(queueURL); visible != 0 || notVisible != 1 {
		t.Fatalf("counts after absent-VisibilityTimeout receive = visible %d, notVisible %d, want 0/1", visible, notVisible)
	}

	// Explicit 0 — the message is received and stays immediately visible.
	if _, err := store.SendMessage(queueURL, sqsstore.NewMessage("vt-zero")); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	zeroResp, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ReceiveMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":          queueURL,
			"VisibilityTimeout": "0",
		},
	})
	if err != nil {
		t.Fatalf("ReceiveMessage with VisibilityTimeout 0: %v", err)
	}
	zeroMap, _ := zeroResp.(map[string]interface{})
	if msgs, _ := zeroMap["Messages"].([]map[string]interface{}); len(msgs) != 1 {
		t.Fatalf("explicit-0 VisibilityTimeout: messages = %v, want 1", zeroMap)
	}
	// The first message stays hidden by the queue default; the explicit-0
	// message is visible again at once.
	if visible, notVisible, _ := store.GetMessageCounts(queueURL); visible != 1 || notVisible != 1 {
		t.Fatalf("counts after explicit-0 receive = visible %d, notVisible %d, want 1/1", visible, notVisible)
	}
}

// TestListDeadLetterSourceQueuesExplicitZeroRejected pins the documented
// "Value range is 1 to 1000" for ListDeadLetterSourceQueues MaxResults: an
// explicitly supplied 0 rejects exactly like ListQueues, while an omitted
// member selects the default page. The checks run before any store access,
// so a nil store is safe.
func TestListDeadLetterSourceQueuesExplicitZeroRejected(t *testing.T) {
	s := &SQSService{}
	queueURL := "http://localhost:50080/000000000000/queue-a"
	if _, err := s.listDeadLetterSourceQueuesCore(nil, ListDeadLetterSourceQueuesInput{
		QueueURL:      queueURL,
		MaxResults:    0,
		MaxResultsSet: true,
	}); err != ErrInvalidParameterValue {
		t.Fatalf("explicit MaxResults=0: got %v, want ErrInvalidParameterValue", err)
	}
	if _, err := s.listDeadLetterSourceQueuesCore(nil, ListDeadLetterSourceQueuesInput{
		QueueURL:   queueURL,
		MaxResults: 1001,
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("MaxResults above the range on a nil store: expected InvalidParameterValue before store access, got %v", err)
	}
}

// TestReceiveMessageFractionalAndOverflowRejected pins the strict-integer
// contract end to end for the JSON number forms: a fractional
// MaxNumberOfMessages and one beyond the int32 domain are wire-type
// violations rejected with SerializationException, never truncated (2.5 →
// 2) or wrapped past the range check (4294967297 → 1). The validation runs
// before any store access, so a nil store is safe.
func TestReceiveMessageFractionalAndOverflowRejected(t *testing.T) {
	s := &SQSService{}
	for name, value := range map[string]float64{
		"fractional":     2.5,
		"int32 overflow": 4294967297,
	} {
		_, err := s.receiveMessageCore(nil, ReceiveMessageInput{
			QueueURL:   "http://localhost:50080/000000000000/queue-a",
			Parameters: map[string]interface{}{"MaxNumberOfMessages": value},
		})
		if err == nil || !strings.Contains(err.Error(), "SerializationException") {
			t.Fatalf("%s: expected SerializationException, got %v", name, err)
		}
	}
}
