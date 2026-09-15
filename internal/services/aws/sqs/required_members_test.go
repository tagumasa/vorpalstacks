package sqs

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// TestChangeMessageVisibilityRequiredMembers pins the @required contract of
// ChangeMessageVisibilityRequest.VisibilityTimeout ("Values range: 0 to
// 43200", model): an omitted member is rejected with MissingParameter —
// never silently defaulted to 0, which carries the destructive meaning
// "release for immediate redelivery" — while an explicitly provided 0 stays
// a valid set value, and a present non-integer is a wire-type violation.
// The rejection runs before any store access; the release semantics run on
// the real store.
func TestChangeMessageVisibilityRequiredMembers(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "cmv-required"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	if _, err := store.SendMessage(queueURL, sqsstore.NewMessage("cmv-body")); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	received, err := store.ReceiveMessage(queueURL, 1, nil, 0, "")
	if err != nil || len(received) != 1 {
		t.Fatalf("ReceiveMessage: %v (%d messages)", err, len(received))
	}
	receiptHandle := received[0].ReceiptHandle

	baseParams := map[string]interface{}{
		"QueueUrl":      queueURL,
		"ReceiptHandle": receiptHandle,
	}

	// Omitted member — rejected before the store is touched.
	noTimeout := copyParams(baseParams)
	_, err = svc.ChangeMessageVisibility(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ChangeMessageVisibility",
		Parameters: noTimeout,
	})
	if err == nil || !strings.Contains(err.Error(), "MissingParameter") {
		t.Fatalf("omitted VisibilityTimeout: got %v, want MissingParameter", err)
	}

	// Present non-integer — wire-type violation.
	badType := copyParams(baseParams)
	badType["VisibilityTimeout"] = "soon"
	_, err = svc.ChangeMessageVisibility(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ChangeMessageVisibility",
		Parameters: badType,
	})
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer VisibilityTimeout: got %v, want SerializationException", err)
	}

	// Explicit 0 is a valid set value: the message returns to visible.
	zeroTimeout := copyParams(baseParams)
	zeroTimeout["VisibilityTimeout"] = "0"
	if _, err := svc.ChangeMessageVisibility(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ChangeMessageVisibility",
		Parameters: zeroTimeout,
	}); err != nil {
		t.Fatalf("ChangeMessageVisibility with explicit 0: %v", err)
	}
	if visible, _, _ := store.GetMessageCounts(queueURL); visible != 1 {
		t.Fatalf("visible messages after explicit-0 timeout = %d, want 1", visible)
	}
}

// TestChangeMessageVisibilityBatchEntryTimeout pins the entry-level
// VisibilityTimeout contract: the member is Required: No on
// ChangeMessageVisibilityBatchRequestEntry (model and API reference), so an
// omitted value selects the queue's VisibilityTimeout attribute — the same
// default the member documents on ReceiveMessage — while an explicit 0 is
// honoured. The pre-fix behaviour defaulted an omitted entry to 0 and
// released the message at once.
func TestChangeMessageVisibilityBatchEntryTimeout(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":         "cmv-batch-default",
			"Attribute.1.Name":  "VisibilityTimeout",
			"Attribute.1.Value": "5",
		},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	for i := 0; i < 2; i++ {
		if _, err := store.SendMessage(queueURL, sqsstore.NewMessage("cmv-batch-body")); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}
	}
	received, err := store.ReceiveMessage(queueURL, 2, nil, 0, "")
	if err != nil || len(received) != 2 {
		t.Fatalf("ReceiveMessage: %v (%d messages)", err, len(received))
	}

	// Entry "omit" leaves VisibilityTimeout out (queue default 5 seconds);
	// entry "zero" sets it explicitly to 0.
	batchResp, err := svc.ChangeMessageVisibilityBatch(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ChangeMessageVisibilityBatch",
		Parameters: map[string]interface{}{
			"QueueUrl": queueURL,
			"Entries": []interface{}{
				map[string]interface{}{"Id": "omit", "ReceiptHandle": received[0].ReceiptHandle},
				map[string]interface{}{"Id": "zero", "ReceiptHandle": received[1].ReceiptHandle, "VisibilityTimeout": float64(0)},
			},
		},
	})
	if err != nil {
		t.Fatalf("ChangeMessageVisibilityBatch: %v", err)
	}
	batchMap, _ := batchResp.(map[string]interface{})
	if failed, ok := batchMap["Failed"].([]map[string]interface{}); ok && len(failed) != 0 {
		t.Fatalf("batch reported failures: %v", failed)
	}

	visible, notVisible, _ := store.GetMessageCounts(queueURL)
	if visible != 1 || notVisible != 1 {
		t.Fatalf("after batch: visible=%d notVisible=%d, want 1/1 (queue-default entry hidden, explicit-0 entry released)", visible, notVisible)
	}
}

// TestSetQueueAttributesAbsentAttributesRejected pins the @required
// contract of SetQueueAttributesRequest.Attributes: both an omitted member
// and an explicitly empty map are rejected with MissingParameter instead of
// succeeding as a validation-free no-op. The rejection precedes the queue
// lookup.
func TestSetQueueAttributesAbsentAttributesRejected(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	cases := map[string]map[string]interface{}{
		"omitted member (query form)": {
			"QueueUrl": "http://localhost:50080/123456789012/no-attrs",
		},
		"empty map (JSON form)": {
			"QueueUrl":   "http://localhost:50080/123456789012/no-attrs",
			"Attributes": map[string]interface{}{},
		},
	}
	for name, params := range cases {
		_, err := svc.SetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "SetQueueAttributes",
			Parameters: params,
		})
		if err == nil || !strings.Contains(err.Error(), "MissingParameter") {
			t.Fatalf("%s: got %v, want MissingParameter", name, err)
		}
	}
}

// TestTagOperationsRequiredMembers pins the @required contracts of
// TagQueueRequest.Tags and UntagQueueRequest.TagKeys: an omitted member, an
// empty JSON map and an empty JSON list are all rejected with
// MissingParameter instead of returning success without tagging anything.
func TestTagOperationsRequiredMembers(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "tag-required"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	tagCases := map[string]map[string]interface{}{
		"omitted Tags":   {"QueueUrl": queueURL},
		"empty Tags map": {"QueueUrl": queueURL, "Tags": map[string]interface{}{}},
	}
	for name, params := range tagCases {
		_, err := svc.TagQueue(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "TagQueue",
			Parameters: params,
		})
		if err == nil || !strings.Contains(err.Error(), "MissingParameter") {
			t.Fatalf("TagQueue %s: got %v, want MissingParameter", name, err)
		}
	}

	untagCases := map[string]map[string]interface{}{
		"omitted TagKeys":    {"QueueUrl": queueURL},
		"empty TagKeys list": {"QueueUrl": queueURL, "TagKeys": []interface{}{}},
	}
	for name, params := range untagCases {
		_, err := svc.UntagQueue(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "UntagQueue",
			Parameters: params,
		})
		if err == nil || !strings.Contains(err.Error(), "MissingParameter") {
			t.Fatalf("UntagQueue %s: got %v, want MissingParameter", name, err)
		}
	}
}

// TestParseIdReceiptBatchEntriesRejectsMissingId pins the @required
// Id member on DeleteMessageBatchRequestEntry and
// ChangeMessageVisibilityBatchRequestEntry: an entry without an Id (or a
// non-object entry) rejects the whole request, aligning the id/receipt arms
// with the send-batch arms — never a silent drop from the batch.
func TestParseIdReceiptBatchEntriesRejectsMissingId(t *testing.T) {
	cases := map[string]map[string]interface{}{
		"entry without Id": {
			"Entries": []interface{}{
				map[string]interface{}{"ReceiptHandle": "handle-1"},
			},
		},
		"non-object entry": {
			"Entries": []interface{}{
				"handle-1",
			},
		},
	}
	for name, params := range cases {
		if _, err := parseDeleteBatchEntries(params); err != ErrInvalidParameterValue {
			t.Fatalf("delete arm, %s: got %v, want ErrInvalidParameterValue", name, err)
		}
		if _, err := parseChangeVisibilityBatchEntries(params); err != ErrInvalidParameterValue {
			t.Fatalf("visibility arm, %s: got %v, want ErrInvalidParameterValue", name, err)
		}
	}
}

// copyParams shallow-copies a parameter map so per-case mutations stay local.
func copyParams(base map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(base))
	for k, v := range base {
		out[k] = v
	}
	return out
}

// TestParseIdReceiptBatchEntriesRejectsMissingReceiptHandle pins the
// @required ReceiptHandle member on the id/receipt batch entry shapes: an
// entry without it rejects the whole request on both wire arms, instead of
// surfacing later as a per-entry Failed ReceiptHandleIsInvalid result.
func TestParseIdReceiptBatchEntriesRejectsMissingReceiptHandle(t *testing.T) {
	jsonCase := map[string]interface{}{
		"Entries": []interface{}{
			map[string]interface{}{"Id": "a"},
		},
	}
	if _, err := parseDeleteBatchEntries(jsonCase); err != ErrMissingParameter {
		t.Fatalf("delete JSON arm: got %v, want ErrMissingParameter", err)
	}
	if _, err := parseChangeVisibilityBatchEntries(jsonCase); err != ErrMissingParameter {
		t.Fatalf("visibility JSON arm: got %v, want ErrMissingParameter", err)
	}

	deleteQuery := map[string]interface{}{
		"DeleteMessageBatchRequestEntry.1.Id": "a",
	}
	if _, err := parseDeleteBatchEntries(deleteQuery); err != ErrMissingParameter {
		t.Fatalf("delete query arm: got %v, want ErrMissingParameter", err)
	}

	visibilityQuery := map[string]interface{}{
		"ChangeMessageVisibilityBatchRequestEntry.1.Id": "a",
	}
	if _, err := parseChangeVisibilityBatchEntries(visibilityQuery); err != ErrMissingParameter {
		t.Fatalf("visibility query arm: got %v, want ErrMissingParameter", err)
	}
}

// TestQueryBatchEntriesWithoutIdRejectWhenMembersRemain pins the query-arm
// truncation rule: the first index without an Id ends the list only when no
// other members of that entry are present — an entry-shaped key set without
// its required Id is a malformed entry that rejects the request, matching
// the JSON arm's rejection.
func TestQueryBatchEntriesWithoutIdRejectWhenMembersRemain(t *testing.T) {
	deleteParams := map[string]interface{}{
		"DeleteMessageBatchRequestEntry.1.Id":            "a",
		"DeleteMessageBatchRequestEntry.1.ReceiptHandle": "h1",
		"DeleteMessageBatchRequestEntry.2.ReceiptHandle": "h2",
	}
	if _, err := parseDeleteBatchEntries(deleteParams); err != ErrInvalidParameterValue {
		t.Fatalf("delete query arm: got %v, want ErrInvalidParameterValue", err)
	}

	sendParams := map[string]interface{}{
		"SendMessageBatchRequestEntry.1.Id":          "a",
		"SendMessageBatchRequestEntry.1.MessageBody": "one",
		"SendMessageBatchRequestEntry.2.MessageBody": "two",
	}
	if _, err := parseBatchSendEntries(sendParams); err != ErrInvalidParameterValue {
		t.Fatalf("send query arm: got %v, want ErrInvalidParameterValue", err)
	}

	contiguous := map[string]interface{}{
		"DeleteMessageBatchRequestEntry.1.Id":            "a",
		"DeleteMessageBatchRequestEntry.1.ReceiptHandle": "h1",
	}
	entries, err := parseDeleteBatchEntries(contiguous)
	if err != nil || len(entries) != 1 {
		t.Fatalf("contiguous query list: got %d entries, err %v, want 1, nil", len(entries), err)
	}
}

// TestAddPermissionRequiredListMembers pins the @required list members of
// AddPermissionRequest: absent AWSAccountIds/Actions reject with the same
// MissingParameter family as the shape's other required members (the Go
// SDK validates the members client-side, so the rejection is pinned here;
// the checks run before any store access, so a nil store is safe).
func TestAddPermissionRequiredListMembers(t *testing.T) {
	s := &SQSService{}
	if err := s.addPermissionCore(nil, AddPermissionInput{
		QueueURL: "http://localhost:50080/000000000000/queue-a",
		Label:    "label-1",
		Actions:  []string{"*"},
	}); err != ErrMissingParameter {
		t.Fatalf("absent AWSAccountIds: got %v, want ErrMissingParameter", err)
	}
	if err := s.addPermissionCore(nil, AddPermissionInput{
		QueueURL:      "http://localhost:50080/000000000000/queue-a",
		Label:         "label-1",
		AWSAccountIDs: []string{"123456789012"},
	}); err != ErrMissingParameter {
		t.Fatalf("absent Actions: got %v, want ErrMissingParameter", err)
	}
}
