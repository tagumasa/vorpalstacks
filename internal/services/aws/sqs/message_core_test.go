package sqs

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
)

// TestReceiveMessageCoreExplicitMaxRejected pins the MaxNumberOfMessages
// contract "Valid values: 1 to 10. Default: 1.": a member that is present
// on the wire with a value outside the range — including an explicit 0 —
// is rejected, and only an absent member falls back to the default of 1.
// The Go SDK cannot reach the explicit-0 path (its serialiser drops a
// zero-valued MaxNumberOfMessages), so the rejection is pinned here for
// both wire forms: the query-protocol string and the JSON number. The
// validation runs before any store access, so a nil store is safe.
func TestReceiveMessageCoreExplicitMaxRejected(t *testing.T) {
	s := &SQSService{}
	cases := map[string]interface{}{
		"explicit zero (query string form)": "0",
		"explicit zero (JSON number form)":  float64(0),
		"above the maximum":                 float64(11),
	}
	for name, value := range cases {
		_, err := s.receiveMessageCore(nil, ReceiveMessageInput{
			QueueURL: "http://localhost:50080/000000000000/queue-a",
			Parameters: map[string]interface{}{
				"MaxNumberOfMessages": value,
			},
		})
		if err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
			t.Fatalf("%s: expected InvalidParameterValue, got %v", name, err)
		}
	}
}

// TestReceiveMessageCoreNonIntegerMembersRejected pins the wire-type
// contract for the Integer members of ReceiveMessage: a member that is
// present but not an integer is a shape violation rejected with
// SerializationException (the awsJson1_0 protocol SQS serves), never
// silently treated as an omitted member falling back to a default. The
// rejection runs before any store access, so a nil store is safe.
func TestReceiveMessageCoreNonIntegerMembersRejected(t *testing.T) {
	s := &SQSService{}
	queueURL := "http://localhost:50080/000000000000/queue-a"
	for _, member := range []string{"MaxNumberOfMessages", "WaitTimeSeconds", "VisibilityTimeout"} {
		_, err := s.receiveMessageCore(nil, ReceiveMessageInput{
			QueueURL: queueURL,
			Parameters: map[string]interface{}{
				member: "abc",
			},
		})
		if err == nil || !strings.Contains(err.Error(), "SerializationException") {
			t.Fatalf("%s: expected SerializationException, got %v", member, err)
		}
	}

}

// TestSendMessageFifoIdentifiersOnStandardQueue pins the queue-type contract
// of the FIFO identifiers on a standard queue: MessageDeduplicationId
// "applies only to FIFO (first-in-first-out) queues" and is rejected with
// InvalidParameterValue, while MessageGroupId stays accepted — it is
// documented for standard queues as well (fair queues).
func TestSendMessageFifoIdentifiersOnStandardQueue(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "std-queue-type"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	queueURL := resp.(map[string]interface{})["QueueUrl"].(string)

	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":       queueURL,
			"MessageBody":    "group id on a standard queue is valid (fair queues)",
			"MessageGroupId": "tenant-a",
		},
	}); err != nil {
		t.Fatalf("MessageGroupId on standard queue: accepted per the fair-queue contract, got %v", err)
	}

	_, err = svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":               queueURL,
			"MessageBody":            "dedup id on a standard queue",
			"MessageDeduplicationId": "dedup-1",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("MessageDeduplicationId on standard queue: expected InvalidParameterValue, got %v", err)
	}
}

// TestReceiveMessageAttemptIdQueueTypeContract pins the queue-type contract
// of ReceiveRequestAttemptId: it "applies only to FIFO (first-in-first-out)
// queues" (ReceiveMessage API reference) and, like MessageDeduplicationId,
// is rejected with InvalidParameterValue on a standard queue instead of
// silently ignored; on a FIFO queue the attempt id keeps its receive-dedup
// meaning and a plain receive is unaffected.
func TestReceiveMessageAttemptIdQueueTypeContract(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	stdResp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "std-attempt-id"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	stdURL := stdResp.(map[string]interface{})["QueueUrl"].(string)

	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":    stdURL,
			"MessageBody": "seed for the standard-queue receive",
		},
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	_, err = svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ReceiveMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":                stdURL,
			"ReceiveRequestAttemptId": "attempt-1",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("ReceiveRequestAttemptId on standard queue: expected InvalidParameterValue, got %v", err)
	}

	fifoResp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":         "fifo-attempt-id.fifo",
			"Attribute.1.Name":  "FifoQueue",
			"Attribute.1.Value": "true",
		},
	})
	if err != nil {
		t.Fatalf("CreateQueue (fifo): %v", err)
	}
	fifoURL := fifoResp.(map[string]interface{})["QueueUrl"].(string)

	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":               fifoURL,
			"MessageBody":            "fifo seed",
			"MessageGroupId":         "g1",
			"MessageDeduplicationId": "seed-dedup-1",
		},
	}); err != nil {
		t.Fatalf("SendMessage (fifo): %v", err)
	}

	recv, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ReceiveMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":                fifoURL,
			"ReceiveRequestAttemptId": "attempt-fifo-1",
			"VisibilityTimeout":       30,
		},
	})
	if err != nil {
		t.Fatalf("ReceiveRequestAttemptId on FIFO queue: accepted, got %v", err)
	}
	if messages, ok := recv.(map[string]interface{})["Messages"].([]map[string]interface{}); !ok || len(messages) != 1 {
		t.Fatalf("FIFO receive with attempt id: expected the seeded message, got %v", messages)
	}
}
