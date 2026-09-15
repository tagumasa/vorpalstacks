package sqs

import (
	"errors"
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// TestStoreErrorTableBatchProjectionLaw pins the law that keeps the two
// error surfaces drift-proof: every table row's single projection is an AWS
// error whose JSON error code is the batch Failed-entry Code, except the
// rows carrying an explicit documented override; every mapped store sentinel
// describes a client-induced failure, so SenderFault is true for all of
// them; and an unmapped store failure falls back to InternalError with
// SenderFault false on the batch surface while passing through unchanged on
// the single surface.
func TestStoreErrorTableBatchProjectionLaw(t *testing.T) {
	for _, row := range storeErrorTable {
		awsErr, ok := row.single.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("row for %v: single projection is not an AWS error", row.store)
		}
		wantCode := awsErr.GetCode()
		if row.batchCodeOverride != "" {
			wantCode = row.batchCodeOverride
		}
		code, senderFault := mapStoreErrorToBatchCode(row.store)
		if code != wantCode || !senderFault {
			t.Errorf("row for %v: batch projection = (%q, %v), want (%q, true)",
				row.store, code, senderFault, wantCode)
		}
		if got := convertStoreError(row.store); got != row.single {
			t.Errorf("row for %v: single projection = %v, want %v", row.store, got, row.single)
		}
	}

	// The rows the consolidated table added: the batch surface previously
	// downgraded MessageNotInflight and InvalidMessageContents to
	// (InternalError, false), and the single surface had no arm for the raw
	// entry-ID sentinel at all.
	if code, fault := mapStoreErrorToBatchCode(sqsstore.ErrMessageNotInflight); code != "MessageNotInflight" || !fault {
		t.Errorf("MessageNotInflight batch projection = (%q, %v), want (MessageNotInflight, true)", code, fault)
	}
	if code, fault := mapStoreErrorToBatchCode(sqsstore.ErrInvalidMessageContents); code != "InvalidMessageContents" || !fault {
		t.Errorf("InvalidMessageContents batch projection = (%q, %v), want (InvalidMessageContents, true)", code, fault)
	}
	if got := convertStoreError(sqsstore.ErrInvalidBatchEntryId); got != ErrInvalidBatchEntryId {
		t.Errorf("InvalidBatchEntryId single projection = %v, want ErrInvalidBatchEntryId", got)
	}
	if code, fault := mapStoreErrorToBatchCode(sqsstore.ErrInvalidBatchEntryId); code != "InvalidBatchEntryId" || !fault {
		t.Errorf("InvalidBatchEntryId batch projection = (%q, %v), want (InvalidBatchEntryId, true)", code, fault)
	}

	unknown := errors.New("not a mapped sentinel")
	if code, fault := mapStoreErrorToBatchCode(unknown); code != "InternalError" || fault {
		t.Errorf("unmapped batch projection = (%q, %v), want (InternalError, false)", code, fault)
	}
	if got := convertStoreError(unknown); got != unknown {
		t.Errorf("unmapped single projection = %v, want pass-through of the original error", got)
	}
}

// newErrorParityEnv builds a real store with one standard queue configured
// for the parity cases: MaximumMessageSize 1024 (the minimum) so an
// oversized body needs no large fixture, and VisibilityTimeout 0 so a
// received message is immediately no longer in flight.
func newErrorParityEnv(t *testing.T) (*SQSService, sqsstore.SQSStoreInterface, string) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	store := sqsstore.NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	queue := sqsstore.NewQueue("error-parity-q", "us-east-1", "123456789012")
	queue.MaximumMessageSize = 1024
	queue.VisibilityTimeout = 0
	created, err := store.CreateQueue(queue)
	if err != nil {
		t.Fatalf("create queue: %v", err)
	}
	return &SQSService{}, store, created.URL
}

// sendBatchFailedEntry sends one body through the real SendMessageBatch Core
// and returns its single Failed entry.
func sendBatchFailedEntry(t *testing.T, s *SQSService, store sqsstore.SQSStoreInterface, queueURL, body string) map[string]interface{} {
	t.Helper()
	resp, err := s.sendMessageBatchCore(store, SendMessageBatchInput{
		QueueURL: queueURL,
		Parameters: map[string]interface{}{
			"Entries": []interface{}{map[string]interface{}{"Id": "a", "MessageBody": body}},
		},
	})
	if err != nil {
		t.Fatalf("send batch: %v", err)
	}
	batchResp, _ := resp.(map[string]interface{})
	if batchResp == nil {
		t.Fatalf("send batch: response is %T, want a map", resp)
	}
	failed, _ := batchResp["Failed"].([]map[string]interface{})
	if len(failed) != 1 {
		t.Fatalf("send batch: got %d failed entries, want 1 (response %v)", len(failed), resp)
	}
	return failed[0]
}

// TestErrorParitySingleVsBatchRealPath pins that the same store failure
// surfaces with the same AWS error code on the single operation and on the
// batch Failed entry, through the real store paths: a disallowed-charset
// body (InvalidMessageContents), an oversized body (InvalidParameterValue —
// the oversized family never shared a code with the charset failure), and a
// visibility change on a message that is no longer in flight
// (MessageNotInflight with SenderFault true on the batch entry).
func TestErrorParitySingleVsBatchRealPath(t *testing.T) {
	s, store, queueURL := newErrorParityEnv(t)

	// Charset failure: rejected as the request-level error on SendMessage
	// and reported as a Failed entry with the same code on SendMessageBatch.
	badBody := "bad body \u0000"
	if _, err := s.sendMessageCore(store, SendMessageInput{
		QueueURL:    queueURL,
		MessageBody: badBody,
		Parameters:  map[string]interface{}{},
	}); err != ErrInvalidMessageContents {
		t.Fatalf("charset single send: got %v, want ErrInvalidMessageContents", err)
	}
	entry := sendBatchFailedEntry(t, s, store, queueURL, badBody)
	if entry["Code"] != "InvalidMessageContents" || entry["SenderFault"] != true {
		t.Errorf("charset batch entry: got (%v, %v), want (InvalidMessageContents, true)",
			entry["Code"], entry["SenderFault"])
	}

	// Oversized body: the single send returns the oversized-message error
	// whose code is InvalidParameterValue, and the batch entry reports the
	// same code instead of a free-form one.
	bigBody := strings.Repeat("x", 2048)
	if _, err := s.sendMessageCore(store, SendMessageInput{
		QueueURL:    queueURL,
		MessageBody: bigBody,
		Parameters:  map[string]interface{}{},
	}); err != ErrMessageTooLarge {
		t.Fatalf("oversize single send: got %v, want ErrMessageTooLarge", err)
	}
	if code := ErrMessageTooLarge.GetCode(); code != "InvalidParameterValue" {
		t.Errorf("oversize error code = %q, want InvalidParameterValue", code)
	}
	entry = sendBatchFailedEntry(t, s, store, queueURL, bigBody)
	if entry["Code"] != "InvalidParameterValue" || entry["SenderFault"] != true {
		t.Errorf("oversize batch entry: got (%v, %v), want (InvalidParameterValue, true)",
			entry["Code"], entry["SenderFault"])
	}

	// Not-in-flight: after a receive on a queue with VisibilityTimeout 0 the
	// message is immediately no longer in flight, so extending its visibility
	// fails on the single operation with MessageNotInflight and on the batch
	// surface with a Failed entry carrying SenderFault true.
	if _, err := s.sendMessageCore(store, SendMessageInput{
		QueueURL:    queueURL,
		MessageBody: "inflight probe",
		Parameters:  map[string]interface{}{},
	}); err != nil {
		t.Fatalf("seed send: %v", err)
	}
	received, err := store.ReceiveMessage(queueURL, 1, nil, 0, "")
	if err != nil || len(received) != 1 {
		t.Fatalf("receive: err=%v messages=%d", err, len(received))
	}
	handle := received[0].ReceiptHandle
	if err := s.changeMessageVisibilityCore(store, ChangeMessageVisibilityInput{
		QueueURL:             queueURL,
		ReceiptHandle:        handle,
		VisibilityTimeout:    30,
		VisibilityTimeoutSet: true,
	}); err != ErrMessageNotInflight {
		t.Fatalf("not-in-flight single CMV: got %v, want ErrMessageNotInflight", err)
	}
	resp, err := s.changeMessageVisibilityBatchCore(store, ChangeMessageVisibilityBatchInput{
		QueueURL: queueURL,
		Parameters: map[string]interface{}{
			"Entries": []interface{}{map[string]interface{}{
				"Id":                "a",
				"ReceiptHandle":     handle,
				"VisibilityTimeout": 30,
			}},
		},
	})
	if err != nil {
		t.Fatalf("not-in-flight CMV batch: %v", err)
	}
	batchResp, _ := resp.(map[string]interface{})
	if batchResp == nil {
		t.Fatalf("not-in-flight CMV batch: response is %T, want a map", resp)
	}
	failed, _ := batchResp["Failed"].([]map[string]interface{})
	if len(failed) != 1 || failed[0]["Code"] != "MessageNotInflight" || failed[0]["SenderFault"] != true {
		t.Errorf("not-in-flight batch entry: got %v, want one entry with Code MessageNotInflight and SenderFault true", failed)
	}
}
