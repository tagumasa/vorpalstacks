package sqs

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
)

// Pins for the message-move-task lifecycle: recovery of tasks orphaned by an
// unclean shutdown, the single-active rule under concurrency, the DLQ-only
// source contract, deduplication gated on the destination's opt-in exactly
// like SendMessage, in-flight messages left to their consumer, the rate
// throttle across the documented 1-500 range, and termination when a batch
// can move nothing.

// armRedrive points source's RedrivePolicy at dlq with maxReceiveCount 1, so
// a second receive of the same message redrives it into the DLQ.
func armRedrive(t *testing.T, store *SQSStore, source, dlq *Queue) {
	t.Helper()
	if err := store.SetQueueAttributes(source.URL, map[string]string{
		"RedrivePolicy":     redrivePolicyJSON(dlq.ARN, 1),
		"VisibilityTimeout": "0",
	}); err != nil {
		t.Fatalf("arm redrive policy on %s: %v", source.URL, err)
	}
}

// redrivePastCount sends one dedup-ID-free message to the redrive-armed
// source and receives it past maxReceiveCount so the policy redrive lands it
// in the DLQ carrying the source's ARN as its DeadLetterQueueSourceArn. FIFO
// sources carry a message group (FIFO sends require one); the redriven copy
// carries no MessageDeduplicationId, which is what the move path's
// body-hash deduplication branch keys on.
func redrivePastCount(t *testing.T, store *SQSStore, source *Queue, body string) {
	t.Helper()
	msg := NewMessage(body)
	if source.FifoQueue {
		msg.MessageGroupID = "move-group"
	}
	if _, err := store.SendMessage(source.URL, msg); err != nil {
		t.Fatalf("send to %s: %v", source.URL, err)
	}
	zero := int32(0)
	for i := 0; i < 2; i++ {
		if _, err := store.ReceiveMessage(source.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("receive %d on %s: %v", i, source.URL, err)
		}
	}
}

// loadMoveTaskRecord reads one message-move-task record by its store key —
// the white-box reader the lifecycle pins use in place of a single-task
// getter on the store surface.
func loadMoveTaskRecord(store *SQSStore, taskID string) (*MessageMoveTask, error) {
	var taskPb pb.MessageMoveTask
	if err := store.tasksStore.GetProto(taskID, &taskPb); err != nil {
		return nil, err
	}
	return ProtoToMessageMoveTask(&taskPb), nil
}

// waitMoveTaskTerminalStatus polls until the task reaches the wanted
// terminal status and returns the final record.
func waitMoveTaskTerminalStatus(t *testing.T, store *SQSStore, taskID, want string) *MessageMoveTask {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		task, err := loadMoveTaskRecord(store, taskID)
		if err == nil {
			switch task.Status {
			case MoveTaskStatusCompleted, MoveTaskStatusFailed, MoveTaskStatusCancelled:
				if task.Status != want {
					t.Fatalf("move task %s ended %s (%s), want %s", taskID, task.Status, task.FailureReason, want)
				}
				return task
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("move task %s did not reach %s within the deadline", taskID, want)
	return nil
}

// TestMoveTaskRecoveryAfterUncleanShutdown pins the construction-time
// recovery sweep: RUNNING and CANCELLING records persisted by a process that
// died without finalising them are terminalised at the next start, so the
// single-active slot frees up and a stale CANCELLING task stops being
// un-cancellable.
func TestMoveTaskRecoveryAfterUncleanShutdown(t *testing.T) {
	dir := t.TempDir()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	store := NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	dlq := createRedriveQueue(t, store, "recover-dlq", false)
	source := createRedriveQueue(t, store, "recover-source", false)
	armRedrive(t, store, source, dlq)

	// Two records an unclean shutdown could leave behind: a task still
	// RUNNING and one whose cancel was recorded but never carried out.
	now := time.Now().UTC()
	stale := []*MessageMoveTask{
		{TaskId: "stale-running", SourceQueueARN: dlq.ARN, DestinationQueueARN: source.ARN, Status: MoveTaskStatusRunning, StartTime: now},
		{TaskId: "stale-cancelling", SourceQueueARN: dlq.ARN, DestinationQueueARN: source.ARN, Status: MoveTaskStatusCancelling, StartTime: now},
	}
	for _, task := range stale {
		if err := store.tasksStore.PutProto(task.TaskId, MessageMoveTaskToProto(task)); err != nil {
			t.Fatalf("seed stale task %s: %v", task.TaskId, err)
		}
	}
	// The crash: no worker exists for these records and the store closes
	// without finalising them.
	store.Close()
	st.Close()

	st2, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("reopen storage: %v", err)
	}
	store2 := NewSQSStore(st2, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		store2.Close()
		st2.Close()
	})

	recovered, err := loadMoveTaskRecord(store2, "stale-running")
	if err != nil {
		t.Fatalf("load recovered RUNNING task: %v", err)
	}
	if recovered.Status != MoveTaskStatusFailed || recovered.FailureReason == "" {
		t.Errorf("recovered RUNNING task: status = %s, failureReason = %q, want FAILED with a reason",
			recovered.Status, recovered.FailureReason)
	}
	recoveredCancel, err := loadMoveTaskRecord(store2, "stale-cancelling")
	if err != nil {
		t.Fatalf("load recovered CANCELLING task: %v", err)
	}
	if recoveredCancel.Status != MoveTaskStatusCancelled {
		t.Errorf("recovered CANCELLING task: status = %s, want CANCELLED", recoveredCancel.Status)
	}

	// Cancel of a recovered (terminal) task is a terminal-state rejection,
	// not an eternal CANCELLING. This must be asserted before the next
	// start: starting a new task purges the source's terminal records.
	if _, err := store2.CancelMessageMoveTask("stale-running"); err != ErrTaskAlreadyTerminal {
		t.Errorf("cancel recovered task: err = %v, want ErrTaskAlreadyTerminal", err)
	}

	// The single-active slot is free again: a new task starts and runs to
	// completion on the empty DLQ.
	started, err := store2.StartMessageMoveTask(dlq.ARN, "", 0)
	if err != nil {
		t.Fatalf("start after recovery: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store2, started.TaskId, MoveTaskStatusCompleted)
}

// TestStartMessageMoveTaskSingleActive pins the check-then-put race fix:
// concurrent starts for one source admit exactly one task; every loser gets
// ErrOverLimit.
func TestStartMessageMoveTaskSingleActive(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "single-dlq", false)
	source := createRedriveQueue(t, store, "single-source", false)
	armRedrive(t, store, source, dlq)
	if _, err := store.SendMessage(dlq.URL, NewMessage("single active")); err != nil {
		t.Fatalf("seed DLQ: %v", err)
	}

	const attempts = 8
	var mu sync.Mutex
	var started, rejected int
	var wg sync.WaitGroup
	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.StartMessageMoveTask(dlq.ARN, "", 0)
			mu.Lock()
			defer mu.Unlock()
			switch err {
			case nil:
				started++
			case ErrOverLimit:
				rejected++
			default:
				t.Errorf("concurrent start: unexpected error %v", err)
			}
		}()
	}
	wg.Wait()
	if started != 1 || rejected != attempts-1 {
		t.Errorf("concurrent starts: %d started, %d rejected; want exactly 1 started and %d rejected",
			started, rejected, attempts-1)
	}
}

// TestStartMessageMoveTaskRequiresDLQSource pins the SourceArn contract:
// "Currently, only ARNs of dead-letter queues (DLQs) whose sources are other
// Amazon SQS queues are accepted" — a queue that nothing redrives into is
// rejected whether or not a destination is named.
func TestStartMessageMoveTaskRequiresDLQSource(t *testing.T) {
	store := newRedriveTestStore(t)
	plain := createRedriveQueue(t, store, "not-a-dlq", false)
	other := createRedriveQueue(t, store, "also-plain", false)

	if _, err := store.StartMessageMoveTask(plain.ARN, other.ARN, 0); err != ErrUnsupportedOperation {
		t.Errorf("start with named destination: err = %v, want ErrUnsupportedOperation", err)
	}
	if _, err := store.StartMessageMoveTask(plain.ARN, "", 0); err != ErrUnsupportedOperation {
		t.Errorf("start with omitted destination: err = %v, want ErrUnsupportedOperation", err)
	}
}

// TestMoveTaskDedupGatedOnDestinationCBD pins the SendMessage gating on the
// move path: a FIFO destination without ContentBasedDeduplication receives
// every same-body copy (unconditional body-hash deduplication used to delete
// the duplicates silently); a CBD destination receives them too, because a
// move registers the window without consulting it — only a later same-key
// send into the destination is suppressed.
func TestMoveTaskDedupGatedOnDestinationCBD(t *testing.T) {
	store := newRedriveTestStore(t)
	zero := int32(0)
	dlq := createRedriveQueue(t, store, "gate-dlq.fifo", true)

	// Non-CBD destination: two redriven same-body copies both deliver.
	sa := createRedriveQueue(t, store, "gate-a.fifo", true)
	sb := createRedriveQueue(t, store, "gate-b.fifo", true)
	armRedrive(t, store, sa, dlq)
	armRedrive(t, store, sb, dlq)
	plainDest := NewQueue("gate-plain-dest.fifo", "us-east-1", "123456789012")
	plainDest.FifoQueue = true
	plainDest, err := store.CreateQueue(plainDest)
	if err != nil {
		t.Fatalf("create plain FIFO destination: %v", err)
	}
	redrivePastCount(t, store, sa, "same body")
	redrivePastCount(t, store, sb, "same body")

	task1, err := store.StartMessageMoveTask(dlq.ARN, plainDest.ARN, 0)
	if err != nil {
		t.Fatalf("start move to non-CBD destination: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store, task1.TaskId, MoveTaskStatusCompleted)
	// The copies share a message group, and a FIFO receive delivers at most
	// one message per group, so the delivered stock is counted rather than
	// received. Nothing has been received from this destination, so both
	// copies are visible.
	delivered, _, _ := store.GetMessageCounts(plainDest.URL)
	if delivered != 2 {
		t.Errorf("non-CBD destination holds %d visible of 2 same-body copies; deduplication was applied without the destination's opt-in", delivered)
	}

	// CBD destination: both same-body copies relocate — a move registers the
	// window without consulting it, so only a later same-key SEND into the
	// destination is suppressed. The delivered stock is counted (the copies
	// share a message group, and a FIFO receive delivers at most one message
	// per group), which is also what distinguishes this from a consult
	// regression: suppression would leave one copy behind.
	sc := createRedriveQueue(t, store, "gate-c.fifo", true)
	sd := createRedriveQueue(t, store, "gate-d.fifo", true)
	armRedrive(t, store, sc, dlq)
	armRedrive(t, store, sd, dlq)
	cbdDest := createRedriveQueue(t, store, "gate-cbd-dest.fifo", true)
	redrivePastCount(t, store, sc, "same body too")
	redrivePastCount(t, store, sd, "same body too")

	task2, err := store.StartMessageMoveTask(dlq.ARN, cbdDest.ARN, 0)
	if err != nil {
		t.Fatalf("start move to CBD destination: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, store, task2.TaskId, MoveTaskStatusCompleted)
	deliveredCBD, err := store.ReceiveMessage(cbdDest.URL, 10, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive from CBD destination: %v", err)
	}
	if len(deliveredCBD) != 1 {
		t.Errorf("CBD destination delivered %d copies, want 1 per receive (one message per group)", len(deliveredCBD))
	}
	cbdStock, cbdInFlight, _ := store.GetMessageCounts(cbdDest.URL)
	if cbdStock+cbdInFlight != 2 {
		t.Errorf("CBD destination holds %d of 2 same-body copies; the move consulted the window it should only register", cbdStock+cbdInFlight)
	}
	if final.MovedMessages != 2 {
		t.Errorf("CBD move: moved = %d, want 2 (both copies relocated)", final.MovedMessages)
	}
}

// TestMoveTaskSkipsInFlightMessages pins the in-flight policy: a message
// under a live receive is not stolen — the task waits out the visibility
// window and moves the message once it is visible again, and the consumer's
// held receipt handle no longer resolves afterwards.
func TestMoveTaskSkipsInFlightMessages(t *testing.T) {
	store := newRedriveTestStore(t)
	zero := int32(0)
	dlq := createRedriveQueue(t, store, "inflight-dlq", false)
	source := createRedriveQueue(t, store, "inflight-source", false)
	armRedrive(t, store, source, dlq)
	dest := createRedriveQueue(t, store, "inflight-dest", false)
	sent, err := store.SendMessage(dlq.URL, NewMessage("mid-flight"))
	if err != nil {
		t.Fatalf("seed DLQ: %v", err)
	}

	visibility := int32(1)
	held, err := store.ReceiveMessage(dlq.URL, 1, &visibility, 0, "")
	if err != nil || len(held) != 1 {
		t.Fatalf("receive in-flight: %v (%d messages)", err, len(held))
	}
	handle := held[0].ReceiptHandle

	start := time.Now()
	task, err := store.StartMessageMoveTask(dlq.ARN, dest.ARN, 0)
	if err != nil {
		t.Fatalf("start move: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)
	// The visibility window (1 s) must have elapsed before the move could
	// take the message; without the skip the task drains it in the first
	// batch within ~100 ms. The bound sits just under the window a correct
	// run cannot beat, so a regression needs machine load an order of
	// magnitude beyond the drain time to slip through.
	if elapsed := time.Since(start); elapsed < 900*time.Millisecond {
		t.Errorf("move of an in-flight message completed in %v; the visibility window was not waited out", elapsed)
	}
	if final.MovedMessages != 1 || final.Status != MoveTaskStatusCompleted {
		t.Errorf("move task: status %s moved %d, want COMPLETED with 1 moved", final.Status, final.MovedMessages)
	}
	arrived, err := store.ReceiveMessage(dest.URL, 1, &zero, 0, "")
	if err != nil || len(arrived) != 1 || arrived[0].ID != sent.ID {
		t.Fatalf("destination receive: %v (%d messages), want the moved %s", err, len(arrived), sent.ID)
	}
	if err := store.DeleteMessage(dlq.URL, handle); err == nil {
		t.Error("delete with the pre-move receipt handle succeeded; the handle outlived the message it was issued for")
	}
}

// TestMoveTaskRateThrottlesAboveBatchSize pins the rate formula across the
// documented range: a rate at or above the batch size still throttles — 101
// messages at 500/s need at least one 200 ms inter-batch window.
func TestMoveTaskRateThrottlesAboveBatchSize(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "rate-dlq", false)
	source := createRedriveQueue(t, store, "rate-source", false)
	armRedrive(t, store, source, dlq)
	dest := createRedriveQueue(t, store, "rate-dest", false)
	for i := 0; i <= moveTaskBatchSize; i++ {
		if _, err := store.SendMessage(dlq.URL, NewMessage(fmt.Sprintf("payload %d", i))); err != nil {
			t.Fatalf("seed message %d: %v", i, err)
		}
	}

	start := time.Now()
	task, err := store.StartMessageMoveTask(dlq.ARN, dest.ARN, 500)
	if err != nil {
		t.Fatalf("start move at 500/s: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)
	if final.MovedMessages != int32(moveTaskBatchSize+1) {
		t.Errorf("moved = %d, want %d", final.MovedMessages, moveTaskBatchSize+1)
	}
	if elapsed := time.Since(start); elapsed < 200*time.Millisecond {
		t.Errorf("101-message move at 500/s finished in %v; the inter-batch throttle did not engage", elapsed)
	}
}

// TestMoveTaskFailureCapTerminates pins the termination rule: a batch that
// moved nothing but failed messages finalises the task FAILED instead of
// re-listing the same failures forever, and the failed messages stay in the
// source queue. The transaction failures are injected with the
// failingUpdateStorage wrapper from lifecycle_test.go.
func TestMoveTaskFailureCapTerminates(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	// The world is seeded on a healthy store, which is then closed so the
	// failing store is the only actor on the directory. The DLQ message is
	// produced through the redrive path (not a direct send) so it carries the
	// DeadLetterQueueSourceArn the omitted-destination routing keys on — a
	// directly-seeded message would be unroutable and never reach the failing
	// transaction the cap is pinned against.
	healthy := NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	dlq := createRedriveQueue(t, healthy, "cap-dlq", false)
	source := createRedriveQueue(t, healthy, "cap-source", false)
	armRedrive(t, healthy, source, dlq)
	if _, err := healthy.SendMessage(source.URL, NewMessage("cannot move")); err != nil {
		t.Fatalf("seed source: %v", err)
	}
	zero := int32(0)
	for i := 0; i < 2; i++ {
		if _, err := healthy.ReceiveMessage(source.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("redrive receive %d: %v", i, err)
		}
	}
	healthy.Close()

	failing := NewSQSStore(&failingUpdateStorage{TransactionalStorage: st, fail: true},
		"123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		failing.Close()
		st.Close()
	})

	task, err := failing.StartMessageMoveTask(dlq.ARN, "", 0)
	if err != nil {
		t.Fatalf("start move on failing transactions: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, failing, task.TaskId, MoveTaskStatusFailed)
	if final.FailureMessages != 1 {
		t.Errorf("failed count = %d, want 1", final.FailureMessages)
	}
	visible, _, _ := failing.GetMessageCounts(dlq.URL)
	if visible != 1 {
		t.Errorf("source queue holds %d visible messages after the failed task, want the 1 unmoved", visible)
	}
}

// TestDeleteQueueCleansMoveTaskRecords pins the residue rule: deleting a
// queue removes the move-task records that name it as source or destination.
func TestDeleteQueueCleansMoveTaskRecords(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "residue-dlq", false)
	source := createRedriveQueue(t, store, "residue-source", false)
	armRedrive(t, store, source, dlq)

	now := time.Now().UTC()
	records := []*MessageMoveTask{
		{TaskId: "residue-running", SourceQueueARN: dlq.ARN, DestinationQueueARN: source.ARN, Status: MoveTaskStatusRunning, StartTime: now},
		{TaskId: "residue-terminal", SourceQueueARN: source.ARN, DestinationQueueARN: dlq.ARN, Status: MoveTaskStatusCompleted, StartTime: now},
	}
	for _, task := range records {
		if err := store.tasksStore.PutProto(task.TaskId, MessageMoveTaskToProto(task)); err != nil {
			t.Fatalf("seed task %s: %v", task.TaskId, err)
		}
	}

	if err := store.DeleteQueue(dlq.URL); err != nil {
		t.Fatalf("delete DLQ: %v", err)
	}
	if tasks, err := store.ListMessageMoveTasks(dlq.ARN, 10); err != nil || len(tasks) != 0 {
		t.Errorf("tasks naming the deleted queue as source: %v (%d records), want none", err, len(tasks))
	}
	if tasks, err := store.ListMessageMoveTasks(source.ARN, 10); err != nil || len(tasks) != 0 {
		t.Errorf("tasks naming the deleted queue as destination: %v (%d records), want none", err, len(tasks))
	}
}

// TestFindSourceQueueForDLQPaginates pins the full-bucket scan: with more
// queues than one list page, a redrive source that sorts into the second
// page is still found (the DLQ-only source rule depends on it).
func TestFindSourceQueueForDLQPaginates(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "page-dlq", false)
	for i := 0; i < 105; i++ {
		createRedriveQueue(t, store, fmt.Sprintf("page-filler-%03d", i), false)
	}
	last := createRedriveQueue(t, store, "zz-page-source", false)
	if err := store.SetQueueAttributes(last.URL, map[string]string{
		"RedrivePolicy": redrivePolicyJSON(dlq.ARN, 3),
	}); err != nil {
		t.Fatalf("arm page-crossing redrive policy: %v", err)
	}

	if got := store.findSourceQueueForDLQ(dlq.URL); got != last.URL {
		t.Errorf("findSourceQueueForDLQ = %q, want the second-page resident %q", got, last.URL)
	}
}

// TestMoveTaskCompletionCountsAndTerminalCancel pins the completed record's
// counts and the terminal-state rejection of CancelMessageMoveTask.
func TestMoveTaskCompletionCountsAndTerminalCancel(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "counts-dlq", false)
	source := createRedriveQueue(t, store, "counts-source", false)
	armRedrive(t, store, source, dlq)
	dest := createRedriveQueue(t, store, "counts-dest", false)
	for i := 0; i < 3; i++ {
		if _, err := store.SendMessage(dlq.URL, NewMessage(fmt.Sprintf("counted %d", i))); err != nil {
			t.Fatalf("seed message %d: %v", i, err)
		}
	}

	task, err := store.StartMessageMoveTask(dlq.ARN, dest.ARN, 0)
	if err != nil {
		t.Fatalf("start move: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)
	if final.MovedMessages != 3 || final.FailureMessages != 0 || final.ApproximateNumberOfMessagesToMove != 3 {
		t.Errorf("completed task counts: moved %d, failed %d, toMove %d; want 3/0/3",
			final.MovedMessages, final.FailureMessages, final.ApproximateNumberOfMessagesToMove)
	}
	if _, err := store.CancelMessageMoveTask(task.TaskId); err != ErrTaskAlreadyTerminal {
		t.Errorf("cancel completed task: err = %v, want ErrTaskAlreadyTerminal", err)
	}
}

// TestMoveWorkerStopsWhenRecordDeleted pins the zombie-worker rule:
// DeleteQueue removes move-task records without signalling the worker, so
// the worker's own record check is the stop — a missing record (a worker
// that is unobservable and uncancellable by definition) stops the loop. The
// source URL may already name a recreated different incarnation, whose
// messages the zombie must never relocate.
func TestMoveWorkerStopsWhenRecordDeleted(t *testing.T) {
	store := newSQSTestStore(t)
	dlq := createRedriveQueue(t, store, "zombie-dlq", false)
	source := createRedriveQueue(t, store, "zombie-src", false)
	armRedrive(t, store, source, dlq)

	// One in-flight message keeps the worker polling (it skips messages
	// under a live receive instead of completing).
	inflight := NewMessage("in flight")
	if _, err := store.SendMessage(dlq.URL, inflight); err != nil {
		t.Fatalf("send in-flight: %v", err)
	}
	long := int32(120)
	if _, err := store.ReceiveMessage(dlq.URL, 1, &long, 0, ""); err != nil {
		t.Fatalf("receive in-flight: %v", err)
	}

	task, err := store.StartMessageMoveTask(dlq.ARN, source.ARN, 0)
	if err != nil {
		t.Fatalf("StartMessageMoveTask: %v", err)
	}
	// Give the worker one poll cycle, then remove its record (what
	// DeleteQueue does through deleteMoveTasksReferencingQueue).
	time.Sleep(400 * time.Millisecond)
	if err := store.tasksStore.Delete(task.TaskId); err != nil {
		t.Fatalf("delete task record: %v", err)
	}

	// A message that becomes movable after the record vanished stays put:
	// the zombie must not relocate it.
	fresh := NewMessage("arrived after the record died")
	if _, err := store.SendMessage(dlq.URL, fresh); err != nil {
		t.Fatalf("send fresh: %v", err)
	}
	time.Sleep(800 * time.Millisecond)
	visible, notVisible, _ := store.GetMessageCounts(dlq.URL)
	if visible != 1 || notVisible != 1 {
		t.Fatalf("dlq counts (visible %d, notVisible %d): the zombie relocated something (want 1/1)", visible, notVisible)
	}
	if moved, _, _ := store.GetMessageCounts(source.URL); moved != 0 {
		t.Fatalf("destination received %d messages from a zombie worker", moved)
	}
}

// flakyTaskReadStorage fails the move-task bucket's reads while armed,
// standing in for a transient storage read failure inside the worker's
// record check. Message-plane reads pass through untouched so the task can
// still move once the record becomes readable again.
type flakyTaskReadStorage struct {
	storage.TransactionalStorage
	tasksBucket string
	armed       atomic.Bool
}

func (f *flakyTaskReadStorage) Bucket(name string) storage.Bucket {
	inner := f.TransactionalStorage.Bucket(name)
	if name != f.tasksBucket {
		return inner
	}
	return &flakyTaskReadBucket{Bucket: inner, armed: &f.armed}
}

type flakyTaskReadBucket struct {
	storage.Bucket
	armed *atomic.Bool
}

func (b *flakyTaskReadBucket) Get(k []byte) ([]byte, error) {
	if b.armed.Load() {
		return nil, errors.New("injected transient task-record read failure")
	}
	return b.Bucket.Get(k)
}

// TestMoveWorkerRetriesTransientRecordReadError pins the read-error leg of
// the worker's record check: while the task record reads fail, the worker
// keeps retrying (the record stays RUNNING, nothing moves) instead of
// stopping and wedging the source's single-active slot until a restart;
// once reads recover the task completes normally.
func TestMoveWorkerRetriesTransientRecordReadError(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	wrapped := &flakyTaskReadStorage{TransactionalStorage: st, tasksBucket: auxiliaryBucketName("sqs-move-tasks", "us-east-1")}
	store := NewSQSStore(wrapped, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	dlq := createRedriveQueue(t, store, "flaky-dlq", false)
	source := createRedriveQueue(t, store, "flaky-src", false)
	armRedrive(t, store, source, dlq)
	dest := createRedriveQueue(t, store, "flaky-dest", false)
	if _, err := store.SendMessage(dlq.URL, NewMessage("movable under flaky reads")); err != nil {
		t.Fatalf("seed DLQ: %v", err)
	}

	wrapped.armed.Store(true)
	task, err := store.StartMessageMoveTask(dlq.ARN, dest.ARN, 0)
	if err != nil {
		t.Fatalf("StartMessageMoveTask: %v", err)
	}

	// Two full retry cycles with the reads failing: the record must stay
	// RUNNING (the worker retried, not stopped) and the message unmoved.
	time.Sleep(600 * time.Millisecond)
	tasks, err := store.ListMessageMoveTasks(dlq.ARN, 10)
	if err != nil || len(tasks) != 1 {
		t.Fatalf("ListMessageMoveTasks under failing reads: %v (%d tasks)", err, len(tasks))
	}
	if tasks[0].Status != MoveTaskStatusRunning {
		t.Fatalf("task ended %s while its record reads were failing; the worker stopped instead of retrying", tasks[0].Status)
	}
	if visible, _, _ := store.GetMessageCounts(dlq.URL); visible != 1 {
		t.Fatalf("DLQ holds %d visible messages under failing reads, want 1 (nothing may move)", visible)
	}

	wrapped.armed.Store(false)
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)
	if final.MovedMessages != 1 {
		t.Fatalf("after reads recovered: moved = %d, want 1", final.MovedMessages)
	}
	if visible, _, _ := store.GetMessageCounts(dest.URL); visible != 1 {
		t.Fatalf("destination holds %d messages after recovery, want 1", visible)
	}
}

// TestMoveTaskCancelObservedInsideThrottleWindow pins the throttle window's
// cancel responsiveness: with the rate at its minimum (one message per
// second, a 100 s window per batch) a cancel marked right after the batch
// must finalise CANCELLED within seconds — observed by re-reading the task
// record between window slices — not after the whole window has slept out.
func TestMoveTaskCancelObservedInsideThrottleWindow(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "cancel-throttle-dlq", false)
	source := createRedriveQueue(t, store, "cancel-throttle-src", false)
	armRedrive(t, store, source, dlq)
	dest := createRedriveQueue(t, store, "cancel-throttle-dest", false)

	// One movable message and one held in flight: the first batch moves the
	// former and skips the latter (done stays false), so the worker enters
	// the rate window with work outstanding.
	if _, err := store.SendMessage(dlq.URL, NewMessage("moved before cancel")); err != nil {
		t.Fatalf("seed movable: %v", err)
	}
	if _, err := store.SendMessage(dlq.URL, NewMessage("held in flight")); err != nil {
		t.Fatalf("seed in-flight: %v", err)
	}
	long := int32(60)
	if _, err := store.ReceiveMessage(dlq.URL, 1, &long, 0, ""); err != nil {
		t.Fatalf("hold in-flight message: %v", err)
	}

	task, err := store.StartMessageMoveTask(dlq.ARN, dest.ARN, 1)
	if err != nil {
		t.Fatalf("StartMessageMoveTask at rate 1: %v", err)
	}
	// Let the first batch run (the worker yields 100 ms before it), then
	// cancel inside the 100 s window.
	time.Sleep(300 * time.Millisecond)
	cancelled, err := store.CancelMessageMoveTask(task.TaskId)
	if err != nil {
		t.Fatalf("CancelMessageMoveTask inside the window: %v", err)
	}
	if cancelled.Status != MoveTaskStatusCancelling {
		t.Fatalf("cancel returned status %s, want CANCELLING", cancelled.Status)
	}

	// The 5 s poller deadline is the discriminator: an unobserved cancel
	// sits out the full 100 s window before finalising.
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCancelled)
	if final.MovedMessages != 1 {
		t.Fatalf("cancelled task reports moved = %d, want 1 (the pre-cancel batch)", final.MovedMessages)
	}
	if visible, _, _ := store.GetMessageCounts(dest.URL); visible != 1 {
		t.Fatalf("destination holds %d messages after cancel, want 1", visible)
	}
}

// failingFirstUpdateStorage fails the first storage.Update call made while
// armed, standing in for a per-message transaction failure that clears
// right after: exactly one message's relocation fails while its siblings
// move in the same batch.
type failingFirstUpdateStorage struct {
	storage.TransactionalStorage
	armed    atomic.Bool
	consumed atomic.Bool
}

func (f *failingFirstUpdateStorage) Update(ctx context.Context, fn func(txn storage.Transaction) error) error {
	if f.armed.Load() && !f.consumed.Swap(true) {
		return errors.New("injected one-shot move transaction failure")
	}
	return f.TransactionalStorage.Update(ctx, fn)
}

// TestMoveTaskDoneWithFailuresReportsFailed pins the completion-side
// accounting: a fully walked source that moved some messages but failed one
// finalises FAILED with both counts and keeps the failed message in the
// source — COMPLETED would end the retries and hide the residual failure,
// the same class of hidden work the all-fail rule exists to surface.
func TestMoveTaskDoneWithFailuresReportsFailed(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	wrapped := &failingFirstUpdateStorage{TransactionalStorage: st}
	store := NewSQSStore(wrapped, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	dlq := createRedriveQueue(t, store, "residual-dlq", false)
	source := createRedriveQueue(t, store, "residual-src", false)
	armRedrive(t, store, source, dlq)
	dest := createRedriveQueue(t, store, "residual-dest", false)
	for i := 0; i < 3; i++ {
		if _, err := store.SendMessage(dlq.URL, NewMessage(fmt.Sprintf("residual %d", i))); err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
	}

	wrapped.armed.Store(true)
	task, err := store.StartMessageMoveTask(dlq.ARN, dest.ARN, 0)
	if err != nil {
		t.Fatalf("StartMessageMoveTask: %v", err)
	}

	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusFailed)
	if final.MovedMessages != 2 {
		t.Fatalf("moved = %d, want 2 (the failure's siblings in the same batch)", final.MovedMessages)
	}
	if final.FailureMessages != 1 {
		t.Fatalf("failed = %d, want 1 (the distinct failed message, not one per attempt)", final.FailureMessages)
	}
	if final.FailureReason == "" {
		t.Fatal("FAILED without a reason for the residual failures")
	}
	if visible, _, _ := store.GetMessageCounts(dlq.URL); visible != 1 {
		t.Fatalf("DLQ holds %d visible messages after the task, want 1 (the failed copy stays)", visible)
	}
	if visible, _, _ := store.GetMessageCounts(dest.URL); visible != 2 {
		t.Fatalf("destination holds %d messages, want 2", visible)
	}
}

// TestMoveTaskDropsPriorHopSystemAttributes pins the relocation's
// provenance rewrite: the moved copy carries neither the previous hop's
// DeadLetterQueueSourceArn (it names the redrive relation of the queue the
// copy just left — on a move home it would name the destination itself) nor
// a SqsManagedSseEnabled entry (a queue-level property the receive path
// re-derives from the destination's own setting), both of which the clone
// of the source record's attribute map would otherwise carry along.
func TestMoveTaskDropsPriorHopSystemAttributes(t *testing.T) {
	store := newRedriveTestStore(t)
	zero := int32(0)
	dlq := createRedriveQueue(t, store, "hop-dlq", false)
	source := createRedriveQueue(t, store, "hop-src", false)
	armRedrive(t, store, source, dlq)
	if _, err := store.SendMessage(source.URL, NewMessage("carries hop attributes")); err != nil {
		t.Fatalf("seed source: %v", err)
	}
	// Two receives past maxReceiveCount 1: the policy redrive lands the copy
	// in the DLQ with DeadLetterQueueSourceArn stamped by moveToDLQ.
	for i := 0; i < 2; i++ {
		if _, err := store.ReceiveMessage(source.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("receive %d: %v", i, err)
		}
	}

	// Seed the stale SSE entry on the DLQ record the way a receive on an
	// SSE-enabled queue stamps one into the stored map.
	var record pb.Message
	dlqMsgs, err := store.ReceiveMessage(dlq.URL, 1, &zero, 0, "")
	if err != nil || len(dlqMsgs) != 1 {
		t.Fatalf("receive the redriven copy: %v (%d)", err, len(dlqMsgs))
	}
	if err := store.messagesStore.GetProto(messageKey(dlq.URL, dlqMsgs[0].ID), &record); err != nil {
		t.Fatalf("load DLQ record: %v", err)
	}
	record.Attributes["SqsManagedSseEnabled"] = "true"
	if err := store.messagesStore.PutProto(messageKey(dlq.URL, record.Id), &record); err != nil {
		t.Fatalf("seed stale SSE entry: %v", err)
	}

	task, err := store.StartMessageMoveTask(dlq.ARN, source.ARN, 0)
	if err != nil {
		t.Fatalf("StartMessageMoveTask home: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)

	// The moved-home copy returns with its accumulated receive count past
	// the armed maxReceiveCount, so a later receive on the source would
	// redrive it straight back — the documented clear disarms that.
	if err := store.SetQueueAttributes(source.URL, map[string]string{"RedrivePolicy": ""}); err != nil {
		t.Fatalf("clear redrive policy: %v", err)
	}

	home, err := store.ReceiveMessage(source.URL, 1, &zero, 0, "")
	if err != nil || len(home) != 1 {
		t.Fatalf("receive moved-home copy: %v (%d)", err, len(home))
	}
	if _, ok := home[0].Attributes[deadLetterQueueSourceArnAttr]; ok {
		t.Error("moved-home copy still carries the previous hop's DeadLetterQueueSourceArn")
	}
	if _, ok := home[0].Attributes["SqsManagedSseEnabled"]; ok {
		t.Error("moved-home copy still carries the source queue's SqsManagedSseEnabled value")
	}
}

// TestStartMessageMoveTaskRefusedWhileClosing pins the shutdown race: with
// the closing flag set (Close's first critical section), a start is refused
// and leaves no RUNNING record — the worker can never be spawned against a
// Wait that already observes zero.
func TestStartMessageMoveTaskRefusedWhileClosing(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "closing-dlq", false)
	source := createRedriveQueue(t, store, "closing-src", false)
	armRedrive(t, store, source, dlq)

	store.taskMu.Lock()
	store.closing = true
	store.taskMu.Unlock()

	if _, err := store.StartMessageMoveTask(dlq.ARN, source.ARN, 0); err == nil {
		t.Fatal("StartMessageMoveTask during shutdown: expected refusal")
	}
	tasks, err := store.ListMessageMoveTasks(dlq.ARN, 10)
	if err != nil {
		t.Fatalf("ListMessageMoveTasks: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("refused start left %d task records", len(tasks))
	}
}

// TestMoveWorkerPanicDoesNotPoisonMessagePlane pins the deferred unlock:
// a panic inside the move batch is recovered by the worker's own recover
// and the message mutex is released — the plane stays usable instead of
// deadlocking every later receive, delete and visibility change.
func TestMoveWorkerPanicDoesNotPoisonMessagePlane(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	wrapped := &panickingUpdateStorage{TransactionalStorage: st}
	store := NewSQSStore(wrapped, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	dlq := createRedriveQueue(t, store, "panic-dlq", false)
	source := createRedriveQueue(t, store, "panic-src", false)
	armRedrive(t, store, source, dlq)

	movable := NewMessage("movable")
	if _, err := store.SendMessage(dlq.URL, movable); err != nil {
		t.Fatalf("send: %v", err)
	}

	wrapped.armed = true
	task, err := store.StartMessageMoveTask(dlq.ARN, source.ARN, 0)
	if err != nil {
		t.Fatalf("StartMessageMoveTask: %v", err)
	}

	// The worker finalises FAILED after its recovered panic.
	waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusFailed)

	// The message plane must still answer: a receive on an unrelated queue
	// completes (it would block forever on a poisoned mutex).
	probe, err := store.CreateQueue(NewQueue("panic-probe", "us-east-1", "123456789012"))
	if err != nil {
		t.Fatalf("create probe queue: %v", err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := store.ReceiveMessage(probe.URL, 1, nil, 0, "")
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("post-panic receive: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("message plane deadlocked after a recovered move panic (mutex poisoned)")
	}
}

// panickingUpdateStorage injects a panic inside storage.Update once armed,
// standing in for a storage-layer panic under the move batch.
type panickingUpdateStorage struct {
	storage.TransactionalStorage
	armed bool
}

func (p *panickingUpdateStorage) Update(ctx context.Context, fn func(txn storage.Transaction) error) error {
	if p.armed {
		panic("injected move-batch panic")
	}
	return p.TransactionalStorage.Update(ctx, fn)
}

// TestMoveTaskOmittedDestinationRoutesPerMessage pins the omitted-DestinationArn
// contract: "the messages will be redriven back to their respective original
// source queues" (model) — a DLQ serving two sources returns each message to
// ITS source, not to whichever queue the source scan happens to find first,
// and the task record reports no destination ("this field value will be NULL").
func TestMoveTaskOmittedDestinationRoutesPerMessage(t *testing.T) {
	store := newRedriveTestStore(t)
	zero := int32(0)
	dlq := createRedriveQueue(t, store, "permsg-dlq", false)
	sourceA := createRedriveQueue(t, store, "permsg-a", false)
	sourceB := createRedriveQueue(t, store, "permsg-b", false)
	armRedrive(t, store, sourceA, dlq)
	armRedrive(t, store, sourceB, dlq)
	redrivePastCount(t, store, sourceA, "grew in source a")
	redrivePastCount(t, store, sourceB, "grew in source b")

	task, err := store.StartMessageMoveTask(dlq.ARN, "", 0)
	if err != nil {
		t.Fatalf("start omitted-destination move: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)
	if final.MovedMessages != 2 {
		t.Errorf("moved = %d, want 2", final.MovedMessages)
	}
	if final.DestinationQueueARN != "" {
		t.Errorf("task record DestinationQueueARN = %q, want empty (the request named no destination)", final.DestinationQueueARN)
	}

	// The moved copies keep their accumulated ApproximateReceiveCount, so a
	// receive on an still-armed source would immediately re-redrive them.
	// Disarm both sources before the assertion receives.
	for _, src := range []*Queue{sourceA, sourceB} {
		if err := store.SetQueueAttributes(src.URL, map[string]string{"RedrivePolicy": ""}); err != nil {
			t.Fatalf("disarm %s: %v", src.URL, err)
		}
	}

	fromA, err := store.ReceiveMessage(sourceA.URL, 10, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive from source a: %v", err)
	}
	fromB, err := store.ReceiveMessage(sourceB.URL, 10, &zero, 0, "")
	if err != nil {
		t.Fatalf("receive from source b: %v", err)
	}
	if len(fromA) != 1 || fromA[0].Body != "grew in source a" {
		t.Errorf("source a holds %d messages (first body %q); want exactly its own redriven message", len(fromA), firstBody(fromA))
	}
	if len(fromB) != 1 || fromB[0].Body != "grew in source b" {
		t.Errorf("source b holds %d messages (first body %q); want exactly its own redriven message", len(fromB), firstBody(fromB))
	}
	if visible, _, _ := store.GetMessageCounts(dlq.URL); visible != 0 {
		t.Errorf("DLQ still holds %d messages after the per-message move", visible)
	}
}

// TestMoveTaskOmittedDestinationLeavesOriginlessMessages pins the no-home
// rule: a message sent directly to the DLQ (no DeadLetterQueueSourceArn) has
// no original source to return to and stays there — the task completes with
// nothing moved instead of either delivering it to an arbitrary queue or
// re-listing it forever.
func TestMoveTaskOmittedDestinationLeavesOriginlessMessages(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "nohome-dlq", false)
	source := createRedriveQueue(t, store, "nohome-source", false)
	armRedrive(t, store, source, dlq)
	if _, err := store.SendMessage(dlq.URL, NewMessage("sent straight to the DLQ")); err != nil {
		t.Fatalf("seed DLQ: %v", err)
	}

	task, err := store.StartMessageMoveTask(dlq.ARN, "", 0)
	if err != nil {
		t.Fatalf("start omitted-destination move: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)
	if final.MovedMessages != 0 {
		t.Errorf("moved = %d, want 0 (the message has no original source)", final.MovedMessages)
	}
	visible, _, _ := store.GetMessageCounts(dlq.URL)
	if visible != 1 {
		t.Errorf("DLQ holds %d visible messages, want the 1 origin-less message untouched", visible)
	}
	if moved, _, _ := store.GetMessageCounts(source.URL); moved != 0 {
		t.Errorf("source queue received %d messages an origin-less message had no right to visit", moved)
	}
}

// TestMoveTaskWalksPastUnroutableMessages pins the listing walk: origin-less
// messages (no DeadLetterQueueSourceArn) stay in the DLQ forever, so a front
// page full of them never drains — the move task must walk past them to the
// movable messages beyond instead of reporting COMPLETED with the remaining
// work hidden. Seeded keys use "!"-prefixed ids, which sort before any uuid,
// so the whole first page is unroutable and the movable message sits beyond
// it in key order.
func TestMoveTaskWalksPastUnroutableMessages(t *testing.T) {
	store := newRedriveTestStore(t)
	source := createRedriveQueue(t, store, "walk-past-src", false)
	dlq := createRedriveQueue(t, store, "walk-past-dlq", false)
	armRedrive(t, store, source, dlq)

	redrivePastCount(t, store, source, "the movable message beyond the unroutable page")
	if visible, notVisible, delayed := store.GetMessageCounts(dlq.URL); visible+notVisible+delayed != 1 {
		t.Fatalf("DLQ holds %d/%d/%d after the redrive, want 1", visible, notVisible, delayed)
	}

	for i := 0; i < moveTaskBatchSize+5; i++ {
		originless := NewMessage(fmt.Sprintf("origin-less %d", i))
		originless.ID = fmt.Sprintf("!originless-%03d", i)
		originless.QueueURL = dlq.URL
		originless.QueueARN = dlq.ARN
		if err := store.messagesStore.PutProto(messageKey(dlq.URL, originless.ID), MessageToProto(originless)); err != nil {
			t.Fatalf("seed origin-less %d: %v", i, err)
		}
	}

	task, err := store.StartMessageMoveTask(dlq.ARN, "", 0)
	if err != nil {
		t.Fatalf("start move task: %v", err)
	}
	final := waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)

	if visible, notVisible, delayed := store.GetMessageCounts(source.URL); visible+notVisible+delayed != 1 {
		t.Fatalf("source holds %d/%d/%d after the task; the movable message beyond the unroutable front page must move home", visible, notVisible, delayed)
	}
	if visible, notVisible, delayed := store.GetMessageCounts(dlq.URL); visible+notVisible+delayed != moveTaskBatchSize+5 {
		t.Fatalf("DLQ holds %d/%d/%d after the task; the origin-less messages stay (%d expected)", visible, notVisible, delayed, moveTaskBatchSize+5)
	}
	if final.MovedMessages != 1 {
		t.Fatalf("task reports %d moved; COMPLETED must reflect exactly the one movable message", final.MovedMessages)
	}
}

// TestStartMessageMoveTaskMalformedARNs pins the malformed-identifier rule:
// an ARN argument that does not parse as a queue ARN is rejected as
// InvalidAddress (the operation's documented malformed-identifier error),
// never silently treated as an omitted member or a missing queue.
func TestStartMessageMoveTaskMalformedARNs(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "malformed-dlq", false)
	source := createRedriveQueue(t, store, "malformed-source", false)
	armRedrive(t, store, source, dlq)

	if _, err := store.StartMessageMoveTask("not-an-arn", "", 0); err != ErrInvalidAddress {
		t.Errorf("malformed SourceArn: err = %v, want ErrInvalidAddress", err)
	}
	if _, err := store.StartMessageMoveTask(dlq.ARN, "also not an arn", 0); err != ErrInvalidAddress {
		t.Errorf("malformed DestinationArn: err = %v, want ErrInvalidAddress", err)
	}
}

// firstBody returns the first message's body for failure messages.
func firstBody(messages []*Message) string {
	if len(messages) == 0 {
		return ""
	}
	return messages[0].Body
}

// TestMoveTaskClearsSequenceNumberOnStandardDestination pins the FIFO-only
// contract on moves: a FIFO message (sequence number assigned by its own
// queue's counter, as every FIFO enqueue carries; explicit deduplication ID
// from its send) moved to a standard destination arrives without either —
// "This parameter applies only to FIFO (first-in-first-out) queues" (model,
// on SequenceNumber and MessageDeduplicationId both) — mirroring moveToDLQ's
// dlqSequence staying empty for non-FIFO targets. MessageGroupId is NOT
// FIFO-only (the model documents it for standard queues' fair queues) and
// stays.
func TestMoveTaskClearsSequenceNumberOnStandardDestination(t *testing.T) {
	store := newRedriveTestStore(t)
	zero := int32(0)
	dlq := createRedriveQueue(t, store, "seq-dlq.fifo", true)
	source := createRedriveQueue(t, store, "seq-src.fifo", true)
	armRedrive(t, store, source, dlq)
	carrier := NewMessage("sequence and dedup-id carrier")
	carrier.MessageGroupID = "move-group"
	carrier.MessageDeduplicationID = "explicit-dedup-id"
	if _, err := store.SendMessage(source.URL, carrier); err != nil {
		t.Fatalf("send to source: %v", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := store.ReceiveMessage(source.URL, 1, &zero, 0, ""); err != nil {
			t.Fatalf("receive %d on source: %v", i, err)
		}
	}
	dest := createRedriveQueue(t, store, "seq-dest", false)

	// The redriven copy in the FIFO DLQ carries both FIFO-only members.
	inDLQ, err := store.ReceiveMessage(dlq.URL, 1, &zero, 0, "")
	if err != nil || len(inDLQ) != 1 {
		t.Fatalf("receive from DLQ: %v (%d messages)", err, len(inDLQ))
	}
	if inDLQ[0].SequenceNumber == "" {
		t.Fatal("redriven FIFO message carries no SequenceNumber to clear")
	}
	if inDLQ[0].MessageDeduplicationID == "" {
		t.Fatal("redriven FIFO message carries no MessageDeduplicationId to clear")
	}
	// The receive above left the message in flight under visibility 0: it is
	// visible again immediately, so the move task can take it.
	task, err := store.StartMessageMoveTask(dlq.ARN, dest.ARN, 0)
	if err != nil {
		t.Fatalf("start move: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store, task.TaskId, MoveTaskStatusCompleted)

	arrived, err := store.ReceiveMessage(dest.URL, 10, &zero, 0, "")
	if err != nil || len(arrived) != 1 {
		t.Fatalf("receive from standard destination: %v (%d messages)", err, len(arrived))
	}
	if arrived[0].SequenceNumber != "" {
		t.Errorf("standard destination message carries SequenceNumber %q; the field is FIFO-only", arrived[0].SequenceNumber)
	}
	if _, ok := arrived[0].Attributes["SequenceNumber"]; ok {
		t.Errorf("standard destination message emits a SequenceNumber attribute; the attribute is FIFO-only")
	}
	if arrived[0].MessageDeduplicationID != "" {
		t.Errorf("standard destination message carries MessageDeduplicationId %q; the field is FIFO-only", arrived[0].MessageDeduplicationID)
	}
	if _, ok := arrived[0].Attributes["MessageDeduplicationId"]; ok {
		t.Errorf("standard destination message emits a MessageDeduplicationId attribute; the attribute is FIFO-only")
	}
}
