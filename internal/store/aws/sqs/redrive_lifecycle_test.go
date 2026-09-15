package sqs

import (
	"fmt"
	"strings"
	"testing"
)

// Pins for the redrive/dead-letter lifecycle contract: self-targeting and
// allow-policy-violating policies are rejected at write time, the target is
// re-resolved at move time with a deliver-to-consumer fallback, the empty
// RedrivePolicy value clears the association, maxReceiveCount carries the
// documented 1-1,000 range, and both move paths (policy redrive and the
// message-move task) apply one identity policy — the message keeps its ID,
// keeps its original enqueue timestamp on a standard destination and has it
// reset on a FIFO destination.

func newRedriveTestStore(t *testing.T) *SQSStore {
	t.Helper()
	return newSQSTestStore(t)
}

func createRedriveQueue(t *testing.T, store *SQSStore, name string, fifo bool) *Queue {
	t.Helper()
	queue := NewQueue(name, "us-east-1", "123456789012")
	queue.FifoQueue = fifo
	if fifo {
		queue.ContentBasedDeduplication = true
	}
	created, err := store.CreateQueue(queue)
	if err != nil {
		t.Fatalf("create queue %s: %v", name, err)
	}
	return created
}

func redrivePolicyJSON(dlqARN string, maxReceiveCount int) string {
	return fmt.Sprintf(`{"deadLetterTargetArn":%q,"maxReceiveCount":%d}`, dlqARN, maxReceiveCount)
}

// TestSelfRedriveRejectedAtWrite pins the source≠target rule on both write
// paths: a queue whose dead-letter target resolves to itself is rejected, so
// the move-time Put-then-Delete on one key can never destroy a message.
func TestSelfRedriveRejectedAtWrite(t *testing.T) {
	store := newRedriveTestStore(t)
	source := createRedriveQueue(t, store, "self-redrive", false)

	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(source.ARN, 3),
	}); err != ErrInvalidAttributeValue {
		t.Fatalf("SetQueueAttributes self-target: err = %v, want ErrInvalidAttributeValue", err)
	}
	if got := store.Exists(source.URL); !got {
		t.Fatal("source queue disappeared after rejected self-redrive write")
	}

	selfTarget := NewQueue("self-redrive-create", "us-east-1", "123456789012")
	selfTarget.RedrivePolicy = &RedrivePolicy{DeadLetterTargetARN: "arn:aws:sqs:us-east-1:123456789012:self-redrive-create", MaxReceiveCount: 3}
	if _, err := store.CreateQueue(selfTarget); err != ErrInvalidAttributeValue {
		t.Fatalf("CreateQueue self-target: err = %v, want ErrInvalidAttributeValue", err)
	}
}

// TestDeletedDLQMoveFallsBackToConsumer pins the move-time re-resolution: a
// policy whose target has been deleted fails the move and the receive
// delivers the message to the consumer instead of losing it.
func TestDeletedDLQMoveFallsBackToConsumer(t *testing.T) {
	store := newRedriveTestStore(t)
	source := createRedriveQueue(t, store, "fallback-src", false)
	dlq := createRedriveQueue(t, store, "fallback-dlq", false)

	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy":     redrivePolicyJSON(dlq.ARN, 1),
		"VisibilityTimeout": "0",
	}); err != nil {
		t.Fatalf("set redrive policy: %v", err)
	}

	sent, err := store.SendMessage(source.URL, NewMessage("survives a deleted DLQ"))
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if err := store.DeleteQueue(dlq.URL); err != nil {
		t.Fatalf("delete DLQ: %v", err)
	}

	zero := int32(0)
	// First receive hits the over-count trigger, fails to resolve the DLQ,
	// and must still hand the message to the consumer.
	recv, err := store.ReceiveMessage(source.URL, 1, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive: %v", err)
	}
	if len(recv) != 1 || recv[0].ID != sent.ID {
		t.Fatalf("receive with deleted DLQ = %v, want the sent message %s delivered", recv, sent.ID)
	}
	// The second receive re-triggers the move and must deliver again —
	// no message is ever written under the deleted queue's prefix.
	recv2, err := store.ReceiveMessage(source.URL, 1, &zero, 0, "")
	if err != nil {
		t.Fatalf("second receive: %v", err)
	}
	if len(recv2) != 1 || recv2[0].ID != sent.ID {
		t.Fatalf("second receive with deleted DLQ = %v, want redelivery", recv2)
	}
}

// TestUnrelatedAttributeUpdateNotBrickedByDeletedDLQ pins that the redrive
// target is validated only when the request names a RedrivePolicy: after the
// DLQ is deleted, updates to other attributes succeed and the documented
// recovery path — clearing RedrivePolicy with the empty value — works.
func TestUnrelatedAttributeUpdateNotBrickedByDeletedDLQ(t *testing.T) {
	store := newRedriveTestStore(t)
	source := createRedriveQueue(t, store, "recover-src", false)
	dlq := createRedriveQueue(t, store, "recover-dlq", false)

	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 3),
	}); err != nil {
		t.Fatalf("set redrive policy: %v", err)
	}
	if err := store.DeleteQueue(dlq.URL); err != nil {
		t.Fatalf("delete DLQ: %v", err)
	}

	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"VisibilityTimeout": "35",
	}); err != nil {
		t.Fatalf("unrelated attribute update with deleted DLQ: %v", err)
	}
	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy": "",
	}); err != nil {
		t.Fatalf("clear redrive policy with deleted DLQ: %v", err)
	}

	cleared, err := store.GetQueue(source.URL)
	if err != nil {
		t.Fatalf("get queue: %v", err)
	}
	if cleared.RedrivePolicy != nil {
		t.Errorf("RedrivePolicy after clear = %+v, want nil", cleared.RedrivePolicy)
	}
	if _, stillThere := cleared.Attributes["RedrivePolicy"]; stillThere {
		t.Errorf("raw attribute map still carries RedrivePolicy = %q after clear", cleared.Attributes["RedrivePolicy"])
	}
	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 3),
	}); err == nil {
		t.Error("re-arming a redrive policy whose target stays deleted must fail")
	}
}

// TestParseRedrivePolicyRangeAndClear pins the documented maxReceiveCount
// range on the parser and the write path, and the empty-value clear form.
func TestParseRedrivePolicyRangeAndClear(t *testing.T) {
	rdp, err := ParseRedrivePolicy("")
	if err != nil || rdp != nil {
		t.Fatalf(`ParseRedrivePolicy("") = (%+v, %v), want (nil, nil)`, rdp, err)
	}

	for _, tc := range []struct {
		json string
		want int32
	}{
		{`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":1}`, 1},
		{`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":"1000"}`, 1000},
		{`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq"}`, DefaultMaxReceiveCount},
	} {
		rdp, err := ParseRedrivePolicy(tc.json)
		if err != nil {
			t.Fatalf("ParseRedrivePolicy(%s): %v", tc.json, err)
		}
		if rdp.MaxReceiveCount != tc.want {
			t.Errorf("ParseRedrivePolicy(%s).MaxReceiveCount = %d, want %d", tc.json, rdp.MaxReceiveCount, tc.want)
		}
	}

	for _, json := range []string{
		`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":0}`,
		`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":-3}`,
		`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":1001}`,
		`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":"0"}`,
	} {
		if _, err := ParseRedrivePolicy(json); err == nil {
			t.Errorf("ParseRedrivePolicy(%s) accepted an out-of-range maxReceiveCount", json)
		}
	}

	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "range-dlq", false)
	if err := store.SetQueueAttributes(dlq.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 0),
	}); err == nil {
		t.Error("SetQueueAttributes accepted maxReceiveCount 0 through the write path")
	}
}

// TestRedriveAllowPolicyEnforcedAtTargetValidation pins the destination
// queue's RedriveAllowPolicy at target-validation time: denyAll refuses every
// source, byQueue admits only the listed source ARNs, and allowAll (or an
// unset policy) is the default. The byQueue-only sourceQueueArns rule is
// pinned at the value-format level.
func TestRedriveAllowPolicyEnforcedAtTargetValidation(t *testing.T) {
	for _, tc := range []struct {
		name          string
		allowPolicy   string
		wantAccepted  bool
		wantFormatErr bool
	}{
		{name: "unset allow policy is allowAll", allowPolicy: "", wantAccepted: true},
		{name: "explicit allowAll", allowPolicy: `{"redrivePermission":"allowAll"}`, wantAccepted: true},
		{name: "denyAll refuses every source", allowPolicy: `{"redrivePermission":"denyAll"}`, wantAccepted: false},
		{name: "byQueue without the source listed", allowPolicy: `{"redrivePermission":"byQueue","sourceQueueArns":["arn:aws:sqs:us-east-1:123456789012:somewhere-else"]}`, wantAccepted: false},
		{name: "byQueue with the source listed", allowPolicy: `{"redrivePermission":"byQueue","sourceQueueArns":["SOURCE_ARN"]}`, wantAccepted: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := newRedriveTestStore(t)
			dlq := createRedriveQueue(t, store, "allow-dlq", false)
			source := createRedriveQueue(t, store, "allow-src", false)

			allowJSON := tc.allowPolicy
			if allowJSON != "" {
				allowJSON = strings.ReplaceAll(allowJSON, "SOURCE_ARN", source.ARN)
				if err := store.SetQueueAttributes(dlq.URL, map[string]string{
					"RedriveAllowPolicy": allowJSON,
				}); err != nil {
					t.Fatalf("set RedriveAllowPolicy: %v", err)
				}
			}

			err := store.SetQueueAttributes(source.URL, map[string]string{
				"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 3),
			})
			if tc.wantAccepted && err != nil {
				t.Fatalf("redrive policy rejected by allow policy %s: %v", allowJSON, err)
			}
			if !tc.wantAccepted && err != ErrInvalidAttributeValue {
				t.Fatalf("redrive policy with allow policy %s: err = %v, want ErrInvalidAttributeValue", allowJSON, err)
			}
		})
	}

	// "You can specify this parameter only when the redrivePermission
	// parameter is set to byQueue."
	if err := validateRedriveAllowPolicyJSON(`{"redrivePermission":"allowAll","sourceQueueArns":["arn:aws:sqs:us-east-1:123456789012:q"]}`); err != ErrInvalidParameterValue {
		t.Errorf("sourceQueueArns without byQueue: err = %v, want ErrInvalidParameterValue", err)
	}
	if err := validateRedriveAllowPolicyJSON(`{"redrivePermission":"denyAll","sourceQueueArns":[]}`); err != nil {
		t.Errorf("empty sourceQueueArns with denyAll: %v", err)
	}
}

// TestMoveIdentityPolicyUnifiedAcrossPaths pins one identity policy on both
// move paths: the message keeps its ID everywhere, a standard destination
// keeps the original enqueue timestamp, and a FIFO destination resets it.
func TestMoveIdentityPolicyUnifiedAcrossPaths(t *testing.T) {
	store := newRedriveTestStore(t)
	zero := int32(0)

	// Policy redrive into a standard DLQ: ID and enqueue timestamp survive.
	srcStd := createRedriveQueue(t, store, "ident-src", false)
	dlqStd := createRedriveQueue(t, store, "ident-dlq", false)
	if err := store.SetQueueAttributes(srcStd.URL, map[string]string{
		"RedrivePolicy":     redrivePolicyJSON(dlqStd.ARN, 1),
		"VisibilityTimeout": "0",
	}); err != nil {
		t.Fatalf("set redrive policy: %v", err)
	}
	sentStd, err := store.SendMessage(srcStd.URL, NewMessage("standard identity"))
	if err != nil {
		t.Fatalf("send standard: %v", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := store.ReceiveMessage(srcStd.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("receive standard %d: %v", i, err)
		}
	}
	movedStd, err := store.ReceiveMessage(dlqStd.URL, 1, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive from standard DLQ: %v", err)
	}
	if len(movedStd) != 1 {
		t.Fatalf("standard DLQ receive = %d messages, want 1", len(movedStd))
	}
	if movedStd[0].ID != sentStd.ID {
		t.Errorf("standard redrive: message ID = %s, want the original %s", movedStd[0].ID, sentStd.ID)
	}
	if movedStd[0].SentTimestamp.UnixMilli() != sentStd.SentTimestamp.UnixMilli() {
		t.Errorf("standard redrive: enqueue timestamp = %d, want the original %d",
			movedStd[0].SentTimestamp.UnixMilli(), sentStd.SentTimestamp.UnixMilli())
	}

	// Policy redrive into a FIFO DLQ: ID survives, enqueue timestamp resets.
	srcFifo := createRedriveQueue(t, store, "ident-src.fifo", true)
	dlqFifo := createRedriveQueue(t, store, "ident-dlq.fifo", true)
	if err := store.SetQueueAttributes(srcFifo.URL, map[string]string{
		"RedrivePolicy":     redrivePolicyJSON(dlqFifo.ARN, 1),
		"VisibilityTimeout": "0",
	}); err != nil {
		t.Fatalf("set FIFO redrive policy: %v", err)
	}
	fifoMsg := NewMessage("fifo identity")
	fifoMsg.MessageGroupID = "group-1"
	sentFifo, err := store.SendMessage(srcFifo.URL, fifoMsg)
	if err != nil {
		t.Fatalf("send FIFO: %v", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := store.ReceiveMessage(srcFifo.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("receive FIFO %d: %v", i, err)
		}
	}
	movedFifo, err := store.ReceiveMessage(dlqFifo.URL, 1, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive from FIFO DLQ: %v", err)
	}
	if len(movedFifo) != 1 {
		t.Fatalf("FIFO DLQ receive = %d messages, want 1", len(movedFifo))
	}
	if movedFifo[0].ID != sentFifo.ID {
		t.Errorf("FIFO redrive: message ID = %s, want the original %s", movedFifo[0].ID, sentFifo.ID)
	}
	if !movedFifo[0].SentTimestamp.After(sentFifo.SentTimestamp) {
		t.Errorf("FIFO redrive: enqueue timestamp = %d, want a reset later than %d",
			movedFifo[0].SentTimestamp.UnixMilli(), sentFifo.SentTimestamp.UnixMilli())
	}

	// Message-move task into a custom standard destination: the moved
	// message is the DLQ message itself — same ID, same enqueue timestamp.
	// The DLQ also still holds the earlier policy-redrive copy (visibility 0
	// keeps it receivable), which the task moves under the same policy, so
	// the destination is drained and the seeded copy is located by ID.
	dest := createRedriveQueue(t, store, "ident-dest", false)
	inDLQ, err := store.SendMessage(dlqStd.URL, NewMessage("move-task identity"))
	if err != nil {
		t.Fatalf("seed standard DLQ: %v", err)
	}
	task, err := store.StartMessageMoveTask(dlqStd.ARN, dest.ARN, 0)
	if err != nil {
		t.Fatalf("start move task: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)
	movedAll, err := store.ReceiveMessage(dest.URL, 10, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive from move destination: %v", err)
	}
	var movedTask *Message
	for _, m := range movedAll {
		if m.ID == inDLQ.ID {
			movedTask = m
		}
	}
	if movedTask == nil {
		t.Fatalf("move destination has %d messages, none with the seeded ID %s", len(movedAll), inDLQ.ID)
	}
	if movedTask.SentTimestamp.UnixMilli() != inDLQ.SentTimestamp.UnixMilli() {
		t.Errorf("move task: enqueue timestamp = %d, want the DLQ copy's %d",
			movedTask.SentTimestamp.UnixMilli(), inDLQ.SentTimestamp.UnixMilli())
	}
}

// TestInPlaceMoveDoesNotDestroy pins the same-key guard of the move task: a
// task whose destination is the source queue itself completes without
// deleting the messages (the copy+delete transaction would otherwise Put and
// Delete one key).
func TestInPlaceMoveDoesNotDestroy(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "inplace-dlq", false)
	// The DLQ-only source rule requires a redrive relation before a task
	// can start; the source queue is otherwise unused here.
	inplaceSource := createRedriveQueue(t, store, "inplace-source", false)
	if err := store.SetQueueAttributes(inplaceSource.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 3),
	}); err != nil {
		t.Fatalf("arm in-place redrive policy: %v", err)
	}
	sent, err := store.SendMessage(dlq.URL, NewMessage("stays put"))
	if err != nil {
		t.Fatalf("send: %v", err)
	}

	task, err := store.StartMessageMoveTask(dlq.ARN, dlq.ARN, 0)
	if err != nil {
		t.Fatalf("start in-place move task: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)

	zero := int32(0)
	recv, err := store.ReceiveMessage(dlq.URL, 1, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive after in-place move: %v", err)
	}
	if len(recv) != 1 || recv[0].ID != sent.ID {
		t.Fatalf("receive after in-place move = %v, want the message %s intact", recv, sent.ID)
	}
}

// TestRedrivePreservesApproximateReceiveCountAcrossQueues pins the
// cross-queue accumulation contract: "ApproximateReceiveCount – Returns the
// number of times a message has been received across all queues but not
// deleted" (ReceiveMessage API reference). A message redriven after
// maxReceiveCount=2 source receives carries the accumulated count into the
// DLQ — the DLQ receive reports 4 (three source receives including the
// over-count one, plus the DLQ receive), never a restarted 1.
func TestRedrivePreservesApproximateReceiveCountAcrossQueues(t *testing.T) {
	store := newRedriveTestStore(t)
	source := createRedriveQueue(t, store, "count-src", false)
	dlq := createRedriveQueue(t, store, "count-dlq", false)
	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 2),
	}); err != nil {
		t.Fatalf("SetQueueAttributes: %v", err)
	}

	if _, err := store.SendMessage(source.URL, NewMessage("count probe")); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	zero := int32(0)
	for i := 0; i < 3; i++ {
		if _, err := store.ReceiveMessage(source.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("source receive %d: %v", i, err)
		}
	}

	msgs, err := store.ReceiveMessage(dlq.URL, 1, &zero, 0, "")
	if err != nil {
		t.Fatalf("DLQ receive: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected the redriven message in the DLQ, got %d", len(msgs))
	}
	if got := msgs[0].Attributes["ApproximateReceiveCount"]; got != "4" {
		t.Fatalf("ApproximateReceiveCount after DLQ receive = %q, want \"4\"", got)
	}
	if msgs[0].ApproximateReceiveCount != 4 {
		t.Fatalf("typed ApproximateReceiveCount = %d, want 4", msgs[0].ApproximateReceiveCount)
	}
}

// TestFifoRedriveAssignsDLQSequenceNumber pins the FIFO sequence-number
// invariant on the policy redrive path: every enqueued FIFO message carries
// a sequence number from its own queue's counter, so a message redrived
// into a FIFO DLQ must arrive at consumers with a SequenceNumber (a fresh
// one from the DLQ's counter, not the source's).
func TestFifoRedriveAssignsDLQSequenceNumber(t *testing.T) {
	store := newRedriveTestStore(t)
	source := createRedriveQueue(t, store, "seq-src.fifo", true)
	dlq := createRedriveQueue(t, store, "seq-dlq.fifo", true)
	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 1),
	}); err != nil {
		t.Fatalf("SetQueueAttributes: %v", err)
	}

	msg := NewMessage("sequence probe")
	msg.MessageGroupID = "g1"
	sent, err := store.SendMessage(source.URL, msg)
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	sourceSequence := sent.SequenceNumber

	zero := int32(0)
	for i := 0; i < 2; i++ {
		if _, err := store.ReceiveMessage(source.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("source receive %d: %v", i, err)
		}
	}

	dlqMsgs, err := store.ReceiveMessage(dlq.URL, 1, &zero, 0, "")
	if err != nil {
		t.Fatalf("DLQ receive: %v", err)
	}
	if len(dlqMsgs) != 1 {
		t.Fatalf("expected the redriven message in the DLQ, got %d", len(dlqMsgs))
	}
	got := dlqMsgs[0].SequenceNumber
	if got == "" {
		t.Fatal("FIFO DLQ message carries no SequenceNumber")
	}
	if got == sourceSequence {
		t.Fatalf("DLQ SequenceNumber reuses the source's %q", got)
	}
	if got != dlqMsgs[0].Attributes["SequenceNumber"] {
		t.Fatalf("SequenceNumber attribute %q disagrees with the field %q", dlqMsgs[0].Attributes["SequenceNumber"], got)
	}
}
