// This file holds the activity surface: activity CRUD, the activity-task
// and callback-task records, the worker poll queue and the task-result
// waits that back SendTaskSuccess/Failure/Heartbeat.

// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/store/aws/common"
)

// CreateActivity creates a new activity in the store.
func (s *StepFunctionStore) CreateActivity(ctx context.Context, activity *Activity) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if activity.Name == "" {
		return ErrInvalidARN
	}

	arn := s.buildActivityARN(activity.Name)
	if s.activitiesStore.Exists(arn) {
		return ErrActivityAlreadyExists
	}

	activity.ActivityArn = arn
	activity.CreationDate = time.Now().UTC()

	return s.activitiesStore.Put(arn, activity)
}

// GetActivity retrieves an activity by its ARN. A missing record is the
// not-found sentinel; any other storage fault surfaces as itself.
func (s *StepFunctionStore) GetActivity(ctx context.Context, arn string) (*Activity, error) {
	var activity Activity
	if err := s.activitiesStore.Get(arn, &activity); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}
	return &activity, nil
}

// DeleteActivity removes an activity from the store and cascades deletion
// to any pending tasks associated with it.
func (s *StepFunctionStore) DeleteActivity(ctx context.Context, arn string) error {
	if !s.activitiesStore.Exists(arn) {
		return ErrActivityNotFound
	}

	s.activityQueuesMu.Lock()
	if _, exists := s.activityQueues[arn]; exists {
		// The queue is removed but NOT closed: a worker mid-poll still holds
		// the channel reference, and a close would hand its receive nil —
		// the exact hazard DeleteStateMachine documents. The poller's own
		// context bounds its wait; later polls find no activity.
		delete(s.activityQueues, arn)
	}
	s.activityQueuesMu.Unlock()

	// Cascade delete pending tasks. Collect keys first to avoid mutating during ForEach.
	s.deleteActivityTasks(func(task *ActivityTask) bool {
		return task.ActivityArn == arn
	})

	return s.activitiesStore.Delete(arn)
}

// deleteActivityTasks cascades task-record deletion for every task the
// predicate matches: each record's pending-task channel is closed (its
// waiter sees the closure) and the record is removed. Keys are collected
// before any deletion to avoid mutating the store during the ForEach walk.
// Both deletion cascades — DeleteStateMachine for its machine's tasks,
// DeleteActivity for the activity's — run through this one path.
func (s *StepFunctionStore) deleteActivityTasks(match func(*ActivityTask) bool) {
	var taskKeys []string
	_ = s.tasksStore.ForEach(func(key string, value []byte) error {
		var task ActivityTask
		if err := json.Unmarshal(value, &task); err != nil {
			return nil
		}
		if match(&task) {
			taskKeys = append(taskKeys, key)
		}
		return nil
	})
	for _, key := range taskKeys {
		s.pendingTasksMu.Lock()
		if ch, ok := s.pendingTasks[key]; ok {
			close(ch)
			delete(s.pendingTasks, key)
		}
		s.pendingTasksMu.Unlock()
		_ = s.tasksStore.Delete(key)
	}
}

// ListActivities returns a paginated list of activities.
func (s *StepFunctionStore) ListActivities(ctx context.Context, limit int32, nextToken string) (*ActivityListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Activity](s.activitiesStore, opts, nil)
	if err != nil {
		return nil, err
	}

	return &ActivityListResult{
		Activities: result.Items,
		NextToken:  result.NextMarker,
	}, nil
}

// CreateActivityTask creates a new task for an activity. A caller that
// already minted the task token (so it could embed it in the task input
// via $$.Task.Token) keeps its token; otherwise a fresh one is generated.
func (s *StepFunctionStore) CreateActivityTask(task *ActivityTask) error {
	if task.TaskToken == "" {
		task.TaskToken = uuid.New().String()
	}
	task.Status = "PENDING"
	task.CreatedAt = time.Now().UTC()

	if err := s.tasksStore.Put(task.TaskToken, task); err != nil {
		return err
	}

	s.activityQueuesMu.Lock()
	queue := s.activityQueueLocked(task.ActivityArn)

	select {
	case queue <- task:
		s.activityQueuesMu.Unlock()
		return nil
	default:
		_ = s.tasksStore.Delete(task.TaskToken)
		s.activityQueuesMu.Unlock()
		return ErrActivityQueueFull
	}
}

// CreateCallbackTask registers a task-token record for a .waitForTaskToken
// integration attempt. Unlike CreateActivityTask it never enqueues onto an
// activity poll queue: the token reaches its consumer through the
// integration payload, and SendTaskSuccess/SendTaskFailure/SendTaskHeartbeat
// resolve the record by token alone.
func (s *StepFunctionStore) CreateCallbackTask(task *ActivityTask) error {
	if task.TaskToken == "" {
		return ErrTaskNotFound
	}
	// The consumer holds the token as soon as the integration payload
	// leaves, so the record is RUNNING from creation: SendTaskHeartbeat
	// requires a RUNNING task, and the heartbeat window anchors at
	// CreatedAt. Activity tasks reach RUNNING only when a worker claims
	// them through the poll.
	task.Status = "RUNNING"
	task.CreatedAt = time.Now().UTC()
	return s.tasksStore.Put(task.TaskToken, task)
}

// DeleteCallbackTask removes a callback task record whose submission
// failed. The token was never delivered to any consumer, so no later
// report must resolve it.
func (s *StepFunctionStore) DeleteCallbackTask(taskToken string) {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()
	_ = s.tasksStore.Delete(taskToken)
}

// ActivityTaskPollTimeout bounds how long GetActivityTask holds the
// request open waiting for a task. The Step Functions API reference fixes
// the maximum hold at 60 seconds: "The maximum time the service holds on
// to the request before responding is 60 seconds. If no task is available
// within 60 seconds, the poll returns a taskToken with an empty string."
// It is a variable so tests can shorten the wait.
var ActivityTaskPollTimeout = 60 * time.Second

// activityQueueCapacity bounds an activity's poll queue: a producer that
// finds the queue full fails with ErrActivityQueueFull rather than
// unbalancing the poll contract.
const activityQueueCapacity = 100

// activityQueueLocked returns the activity's poll queue, creating it on
// first use. The caller holds activityQueuesMu.
func (s *StepFunctionStore) activityQueueLocked(activityArn string) chan *ActivityTask {
	queue, exists := s.activityQueues[activityArn]
	if !exists {
		queue = make(chan *ActivityTask, activityQueueCapacity)
		s.activityQueues[activityArn] = queue
	}
	return queue
}

// GetActivityTask retrieves a task from an activity queue, blocking until
// a task becomes available or the caller's context ends (the service
// layer bounds the wait with ActivityTaskPollTimeout).
func (s *StepFunctionStore) GetActivityTask(ctx context.Context, activityArn string, workerName string) (*ActivityTask, error) {
	s.activityQueuesMu.Lock()
	queue := s.activityQueueLocked(activityArn)
	s.activityQueuesMu.Unlock()

	select {
	case task := <-queue:
		if task == nil {
			// A closed queue must never crash the poller; an empty result
			// is the documented no-task answer.
			return nil, nil
		}
		task.Status = "RUNNING"
		task.WorkerName = workerName
		if err := s.tasksStore.Put(task.TaskToken, task); err != nil {
			return nil, err
		}
		return task, nil
	case <-ctx.Done():
		return nil, nil
	}
}

// GetActivityTaskByToken retrieves an activity task by its token. A missing
// record is the not-found sentinel; any other storage fault surfaces as
// itself.
func (s *StepFunctionStore) GetActivityTaskByToken(taskToken string) (*ActivityTask, error) {
	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

// HeartbeatActivityTask records the current time as the last heartbeat for
// the task. Workers call this via SendTaskHeartbeat to keep the task alive
// beyond the heartbeat interval. The actual timeout enforcement happens in
// WaitForTaskResult, which polls LastHeartbeatAt. A missing record is the
// not-found sentinel; any other storage fault surfaces as itself.
func (s *StepFunctionStore) HeartbeatActivityTask(taskToken string) error {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		if common.IsNotFound(err) {
			return ErrTaskNotFound
		}
		return err
	}

	task.LastHeartbeatAt = time.Now().UTC()

	if err := s.tasksStore.Put(taskToken, &task); err != nil {
		return err
	}

	return nil
}

// CompleteActivityTask marks an activity task as completed with output. A
// missing record is the not-found sentinel; any other storage fault
// surfaces as itself.
func (s *StepFunctionStore) CompleteActivityTask(taskToken string, output string) error {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		if common.IsNotFound(err) {
			return ErrTaskNotFound
		}
		return err
	}

	// A terminal task (already reported, or abandoned when its attempt
	// timed out) must not be overwritten: a worker holding a stale token
	// from an earlier attempt must not be able to complete the retry.
	if isTerminalTaskStatus(task.Status) {
		return ErrTaskNotRunning
	}

	task.Status = "SUCCEEDED"
	task.Output = output
	task.CompletedAt = time.Now().UTC()

	if err := s.tasksStore.Put(taskToken, &task); err != nil {
		return err
	}

	s.pendingTasksMu.Lock()
	if ch, ok := s.pendingTasks[taskToken]; ok {
		select {
		case ch <- &ActivityTaskResult{TaskToken: taskToken, Output: output, WorkerName: task.WorkerName}:
		default:
		}
	}
	s.pendingTasksMu.Unlock()

	return nil
}

// FailActivityTask marks an activity task as failed with an error. A
// missing record is the not-found sentinel; any other storage fault
// surfaces as itself.
func (s *StepFunctionStore) FailActivityTask(taskToken string, errorMsg string, cause string) error {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		if common.IsNotFound(err) {
			return ErrTaskNotFound
		}
		return err
	}

	if isTerminalTaskStatus(task.Status) {
		return ErrTaskNotRunning
	}

	task.Status = "FAILED"
	task.Error = errorMsg
	task.Cause = cause
	task.CompletedAt = time.Now().UTC()

	if err := s.tasksStore.Put(taskToken, &task); err != nil {
		return err
	}

	s.pendingTasksMu.Lock()
	if ch, ok := s.pendingTasks[taskToken]; ok {
		select {
		case ch <- &ActivityTaskResult{TaskToken: taskToken, Error: errors.New(errorMsg), Cause: cause, WorkerName: task.WorkerName}:
		default:
		}
	}
	s.pendingTasksMu.Unlock()

	return nil
}

// isTerminalTaskStatus reports whether an activity task record can no
// longer accept a worker report. Once a task has succeeded, failed, or
// been abandoned by a timeout, the token is spent.
func isTerminalTaskStatus(status string) bool {
	switch status {
	case "SUCCEEDED", "FAILED", "TIMED_OUT":
		return true
	}
	return false
}

// markTaskTimedOut flips a task record to TIMED_OUT when the executor
// stops waiting for it, so a worker that later presents the token is
// rejected instead of silently overwriting a dead record.
func (s *StepFunctionStore) markTaskTimedOut(taskToken string) {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		return
	}
	if isTerminalTaskStatus(task.Status) {
		return
	}
	task.Status = "TIMED_OUT"
	_ = s.tasksStore.Put(taskToken, &task)
}

// WaitForTaskResult waits for the result of an activity task.
func (s *StepFunctionStore) WaitForTaskResult(ctx context.Context, taskToken string, timeout time.Duration, heartbeatTimeout time.Duration) (*ActivityTaskResult, error) {
	s.pendingTasksMu.Lock()
	ch := make(chan *ActivityTaskResult, 1)
	s.pendingTasks[taskToken] = ch
	s.pendingTasksMu.Unlock()

	defer func() {
		s.pendingTasksMu.Lock()
		delete(s.pendingTasks, taskToken)
		s.pendingTasksMu.Unlock()
	}()

	// A worker may complete the task before this channel is registered —
	// the completion notify only reaches a registered waiter. Re-checking
	// the persisted status after registration closes that window: any
	// completion that landed earlier is visible here, and any later
	// completion finds the registered channel. Without this check a
	// fast completion sleeps the executor until the task timeout.
	if task, err := s.GetActivityTaskByToken(taskToken); err == nil {
		switch task.Status {
		case "SUCCEEDED":
			return &ActivityTaskResult{TaskToken: taskToken, Output: task.Output, WorkerName: task.WorkerName}, nil
		case "FAILED":
			return &ActivityTaskResult{TaskToken: taskToken, Error: errors.New(task.Error), Cause: task.Cause, WorkerName: task.WorkerName}, nil
		case "TIMED_OUT":
			return nil, ErrTaskTimeout
		}
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// Heartbeat enforcement. When heartbeatTimeout > 0, the caller must
	// send a heartbeat at least once per heartbeatTimeout. SendTaskHeartbeat
	// updates LastHeartbeatAt; this loop polls for updates.
	if heartbeatTimeout <= 0 {
		select {
		case result := <-ch:
			if result == nil {
				// The waiter was orphaned (its task record was
				// cascade-deleted): a nil result must surface as a task
				// error, not a nil dereference at the consumer.
				return nil, ErrTaskNotFound
			}
			return result, nil
		case <-timer.C:
			s.markTaskTimedOut(taskToken)
			return nil, ErrTaskTimeout
		case <-ctx.Done():
			s.markTaskTimedOut(taskToken)
			return nil, ctx.Err()
		}
	}

	hbTimer := time.NewTimer(heartbeatTimeout)
	defer hbTimer.Stop()

	for {
		select {
		case result := <-ch:
			if result == nil {
				// The waiter was orphaned (its task record was
				// cascade-deleted): a nil result must surface as a task
				// error, not a nil dereference at the consumer.
				return nil, ErrTaskNotFound
			}
			return result, nil
		case <-timer.C:
			s.markTaskTimedOut(taskToken)
			return nil, ErrTaskTimeout
		case <-hbTimer.C:
			// The re-read goes through the discriminating getter: a
			// missing record keeps the not-found sentinel while a storage
			// fault surfaces as itself.
			task, err := s.GetActivityTaskByToken(taskToken)
			if err != nil {
				return nil, err
			}
			lastHB := task.LastHeartbeatAt
			if lastHB.IsZero() {
				lastHB = task.CreatedAt
			}
			if time.Since(lastHB) > heartbeatTimeout {
				s.markTaskTimedOut(taskToken)
				return nil, ErrHeartbeatTimeout
			}
			hbTimer.Reset(heartbeatTimeout)
		case <-ctx.Done():
			s.markTaskTimedOut(taskToken)
			return nil, ctx.Err()
		}
	}
}

// NewExecution creates a new execution for a state machine.
func NewExecution(stateMachineArn, name, input, traceHeader string) *Execution {
	return &Execution{
		StateMachineArn: stateMachineArn,
		Name:            name,
		Input:           input,
		TraceHeader:     traceHeader,
		Status:          "RUNNING",
		StartDate:       time.Now().UTC(),
	}
}

// executionRegistration is one goroutine's claim on running an execution:
// the cancel function StopExecution needs, and a done channel closed at
// unregister time, after the goroutine's last history and record write.
type executionRegistration struct {
	cancel context.CancelFunc
	done   chan struct{}
}

// ExecutionHandle identifies one execution registration; only the store
// can close the channel it wraps.
type ExecutionHandle struct {
	done chan struct{}
}
