// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// StepFunctionStore provides Step Functions state machine and execution storage.
type StepFunctionStore struct {
	*common.BaseStore
	executionsStore       *common.BaseStore
	executionHistoryStore *common.BaseStore
	activitiesStore       *common.BaseStore
	tasksStore            *common.BaseStore
	versionsStore         *common.BaseStore
	aliasesStore          *common.BaseStore
	mapRunsStore          *common.BaseStore
	*common.TagStore
	arnBuilder         *svcarn.ARNBuilder
	accountID          string
	region             string
	pendingTasks       map[string]chan *ActivityTaskResult
	pendingTasksMu     sync.RWMutex
	activityQueues     map[string]chan *ActivityTask
	activityQueuesMu   sync.RWMutex
	tasksMu            sync.Mutex
	eventIdCounter     int64
	versionCounters    map[string]int64
	versionCountersMu  sync.Mutex
	activeExecutions   map[string]context.CancelFunc
	activeExecutionsMu sync.RWMutex
	createMu           sync.Mutex
	mapRunSeq          int64
}

// NewStepFunctionStore creates a new StepFunctionStore instance.
func NewStepFunctionStore(store storage.BasicStorage, accountID, region string) *StepFunctionStore {
	return &StepFunctionStore{
		BaseStore:             common.NewBaseStore(store.Bucket("stepfunction-statemachines-"+region), "stepfunction-statemachines"),
		executionsStore:       common.NewBaseStore(store.Bucket("stepfunction-executions-"+region), "stepfunction-executions"),
		executionHistoryStore: common.NewBaseStore(store.Bucket("stepfunction-history-"+region), "stepfunction-history"),
		activitiesStore:       common.NewBaseStore(store.Bucket("stepfunction-activities-"+region), "stepfunction-activities"),
		tasksStore:            common.NewBaseStore(store.Bucket("stepfunction-tasks-"+region), "stepfunction-tasks"),
		versionsStore:         common.NewBaseStore(store.Bucket("stepfunction-versions-"+region), "stepfunction-versions"),
		aliasesStore:          common.NewBaseStore(store.Bucket("stepfunction-aliases-"+region), "stepfunction-aliases"),
		mapRunsStore:          common.NewBaseStore(store.Bucket("stepfunction-mapruns-"+region), "stepfunction-mapruns"),
		TagStore:              common.NewTagStoreWithRegion(store, "stepfunction", region),
		arnBuilder:            svcarn.NewARNBuilder(accountID, region),
		accountID:             accountID,
		region:                region,
		pendingTasks:          make(map[string]chan *ActivityTaskResult),
		activityQueues:        make(map[string]chan *ActivityTask),
		activeExecutions:      make(map[string]context.CancelFunc),
		versionCounters:       make(map[string]int64),
	}
}

// GetAccountID returns the AWS account ID.
func (s *StepFunctionStore) GetAccountID() string { return s.accountID }

// GetRegion returns the AWS region.
func (s *StepFunctionStore) GetRegion() string { return s.region }

func (s *StepFunctionStore) buildStateMachineARN(name string) string {
	return s.arnBuilder.StepFunctions().StateMachine(name)
}

func (s *StepFunctionStore) buildActivityARN(name string) string {
	return s.arnBuilder.StepFunctions().Activity(name)
}

// buildExecutionHistoryKey renders the storage key for one history event.
// The event ID is zero-padded to the width of the largest int64 so the
// lexicographic key order equals the numeric event-ID order; an unpadded
// decimal would sort event 10 between 1 and 2.
func (s *StepFunctionStore) buildExecutionHistoryKey(executionArn string, eventId int64) string {
	return fmt.Sprintf("%s:%019d", executionArn, eventId)
}

// CreateStateMachine creates a new state machine in the store.
func (s *StepFunctionStore) CreateStateMachine(ctx context.Context, sm *StateMachine) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if sm.Name == "" {
		return ErrInvalidARN
	}

	arn := s.buildStateMachineARN(sm.Name)
	if s.Exists(arn) {
		return ErrStateMachineAlreadyExists
	}

	now := time.Now().UTC()
	sm.StateMachineArn = arn
	sm.CreationDate = now
	sm.UpdateDate = now
	if sm.Status == "" {
		sm.Status = "ACTIVE"
	}
	if sm.Type == "" {
		sm.Type = "STANDARD"
	}

	return s.Put(arn, sm)
}

// GetStateMachine retrieves a state machine by its ARN.
func (s *StepFunctionStore) GetStateMachine(ctx context.Context, arn string) (*StateMachine, error) {
	var sm StateMachine
	if err := s.BaseStore.Get(arn, &sm); err != nil {
		return nil, ErrStateMachineNotFound
	}
	return &sm, nil
}

// GetStateMachineByName retrieves a state machine by its name.
func (s *StepFunctionStore) GetStateMachineByName(ctx context.Context, name string) (*StateMachine, error) {
	arn := s.buildStateMachineARN(name)
	return s.GetStateMachine(ctx, arn)
}

// UpdateStateMachine updates an existing state machine.
func (s *StepFunctionStore) UpdateStateMachine(ctx context.Context, sm *StateMachine) error {
	if !s.Exists(sm.StateMachineArn) {
		return ErrStateMachineNotFound
	}
	sm.UpdateDate = time.Now().UTC()
	return s.Put(sm.StateMachineArn, sm)
}

// DeleteStateMachine removes a state machine from the store and cascades
// deletion to its executions, history events, versions, and aliases.
func (s *StepFunctionStore) DeleteStateMachine(ctx context.Context, arn string) error {
	if !s.Exists(arn) {
		return ErrStateMachineNotFound
	}

	smName := extractStateMachineNameFromArn(arn)

	s.activeExecutionsMu.Lock()
	for execArn, cancel := range s.activeExecutions {
		execSmName := extractStateMachineNameFromExecutionArn(execArn)
		if execSmName == smName {
			cancel()
			delete(s.activeExecutions, execArn)
		}
	}
	s.activeExecutionsMu.Unlock()

	s.pendingTasksMu.Lock()
	oldPending := s.pendingTasks
	s.pendingTasks = make(map[string]chan *ActivityTaskResult)
	s.pendingTasksMu.Unlock()
	for token, ch := range oldPending {
		var task ActivityTask
		if err := s.tasksStore.Get(token, &task); err == nil {
			taskSmName := extractStateMachineNameFromExecutionArn(task.ExecutionArn)
			if taskSmName == smName {
				close(ch)
			}
		} else {
			close(ch)
		}
	}

	// Activity queues are NOT closed here. Activities are independent
	// resources in AWS Step Functions; closing their queues when an
	// unrelated state machine is deleted would cause GetActivityTask to
	// receive nil from a closed channel, triggering a nil-pointer panic.
	// Pending tasks belonging to this state machine are already cleaned
	// up above via the pendingTasks swap-and-close logic.

	// Cascade delete executions and their history.
	// Collect keys first to avoid mutating store during ForEach iteration.
	var execKeys []string
	_ = s.executionsStore.ForEach(func(key string, value []byte) error {
		var exec Execution
		if err := json.Unmarshal(value, &exec); err != nil {
			return nil
		}
		if exec.StateMachineArn == arn {
			execKeys = append(execKeys, key)
			s.deleteExecutionHistory(exec.ExecutionArn)
		}
		return nil
	})
	for _, key := range execKeys {
		_ = s.executionsStore.Delete(key)
	}

	// Cascade delete versions. Collect keys first.
	var verKeys []string
	_ = s.versionsStore.ForEach(func(key string, value []byte) error {
		if strings.HasPrefix(key, arn+":") {
			var v StateMachineVersion
			if err := json.Unmarshal(value, &v); err == nil {
				verKeys = append(verKeys, v.StateMachineVersionArn)
			}
		}
		return nil
	})
	for _, key := range verKeys {
		_ = s.versionsStore.Delete(key)
	}

	// Cascade delete aliases. Collect keys first.
	var aliasKeys []string
	_ = s.aliasesStore.ForEach(func(key string, value []byte) error {
		var alias StateMachineAlias
		if err := json.Unmarshal(value, &alias); err != nil {
			return nil
		}
		if alias.StateMachineArn == arn {
			aliasKeys = append(aliasKeys, key)
		}
		return nil
	})
	for _, key := range aliasKeys {
		_ = s.aliasesStore.Delete(key)
	}

	_ = s.TagStore.Delete(arn)
	return s.BaseStore.Delete(arn)
}

func extractStateMachineNameFromArn(arn string) string {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	if strings.HasPrefix(resource, "stateMachine:") {
		return strings.TrimPrefix(resource, "stateMachine:")
	}
	return ""
}
func (s *StepFunctionStore) deleteExecutionHistory(executionArn string) {
	prefix := executionArn + ":"
	var histKeys []string
	_ = s.executionHistoryStore.ForEach(func(key string, value []byte) error {
		if strings.HasPrefix(key, prefix) {
			histKeys = append(histKeys, key)
		}
		return nil
	})
	for _, key := range histKeys {
		_ = s.executionHistoryStore.Delete(key)
	}
}
func extractStateMachineNameFromExecutionArn(arn string) string {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	if strings.HasPrefix(resource, "execution:") {
		rest := strings.TrimPrefix(resource, "execution:")
		if idx := strings.Index(rest, ":"); idx > 0 {
			return rest[:idx]
		}
		return rest
	}
	return ""
}

// ListStateMachines returns a paginated list of state machines.
func (s *StepFunctionStore) ListStateMachines(ctx context.Context, limit int32, nextToken string) (*StateMachineListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[StateMachine](s.BaseStore, opts, nil)
	if err != nil {
		return nil, err
	}

	return &StateMachineListResult{
		StateMachines: result.Items,
		NextToken:     result.NextMarker,
	}, nil
}

// CreateExecution creates a new execution for a state machine.
func (s *StepFunctionStore) CreateExecution(ctx context.Context, exec *Execution) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if exec.ExecutionArn == "" {
		return ErrInvalidARN
	}

	if s.executionsStore.Exists(exec.ExecutionArn) {
		return ErrExecutionAlreadyExists
	}

	now := time.Now().UTC()
	exec.StartDate = now
	exec.Status = "RUNNING"
	if exec.InputDetails == nil {
		exec.InputDetails = &ExecutionInputDetails{
			Included: exec.Input != "",
			Type:     "JSON",
		}
	}

	return s.executionsStore.Put(exec.ExecutionArn, exec)
}

// GetExecution retrieves an execution by its ARN.
func (s *StepFunctionStore) GetExecution(ctx context.Context, arn string) (*Execution, error) {
	var exec Execution
	if err := s.executionsStore.Get(arn, &exec); err != nil {
		return nil, ErrExecutionNotFound
	}
	return &exec, nil
}

// UpdateExecution updates an existing execution in the store.
func (s *StepFunctionStore) UpdateExecution(ctx context.Context, exec *Execution) error {
	return s.executionsStore.Put(exec.ExecutionArn, exec)
}

// executionMatches reports whether an execution passes the ListExecutions
// filters: state machine, status, map run and redrive state.
func executionMatches(e *Execution, stateMachineArn, statusFilter, mapRunArn, redriveFilter string) bool {
	if stateMachineArn != "" && e.StateMachineArn != stateMachineArn {
		return false
	}
	if statusFilter != "" && e.Status != statusFilter {
		return false
	}
	if mapRunArn != "" && e.MapRunArn != mapRunArn {
		return false
	}
	if redriveFilter != "" {
		switch redriveFilter {
		case "REDRIVEN":
			if e.RedriveCount == 0 {
				return false
			}
		case "NOT_REDRIVEN":
			if e.RedriveCount != 0 {
				return false
			}
		}
	}
	return true
}

// ListExecutions returns a paginated list of executions for a state machine.
func (s *StepFunctionStore) ListExecutions(ctx context.Context, stateMachineArn string, statusFilter string, mapRunArn string, redriveFilter string, limit int32, nextToken string) (*ExecutionListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Execution](s.executionsStore, opts, func(e *Execution) bool {
		return executionMatches(e, stateMachineArn, statusFilter, mapRunArn, redriveFilter)
	})
	if err != nil {
		return nil, err
	}

	return &ExecutionListResult{
		Executions: result.Items,
		NextToken:  result.NextMarker,
	}, nil
}

// ListAllExecutions returns every execution matching the filters without
// pagination. The ListExecutions contract orders results by time rather
// than storage key ("Results are sorted by time, with the most recent
// execution first"), so the service layer fetches the full match set
// before sorting and paging it.
func (s *StepFunctionStore) ListAllExecutions(ctx context.Context, stateMachineArn, statusFilter, mapRunArn, redriveFilter string) ([]*Execution, error) {
	return common.ListMatching[Execution](s.executionsStore, "", func(e *Execution) bool {
		return executionMatches(e, stateMachineArn, statusFilter, mapRunArn, redriveFilter)
	})
}

// CountExecutionHistory returns the number of history events recorded for
// an execution. Redrive eligibility requires fewer than 24,999 events so
// the ExecutionRedriven event and at least one more event still fit.
func (s *StepFunctionStore) CountExecutionHistory(ctx context.Context, executionArn string) (int, error) {
	events, err := common.ListMatching[ExecutionHistoryEvent](s.executionHistoryStore, executionArn+":", nil)
	if err != nil {
		return 0, err
	}
	return len(events), nil
}

// AddExecutionHistoryEvent adds a history event to an execution's event log.
func (s *StepFunctionStore) AddExecutionHistoryEvent(ctx context.Context, event *ExecutionHistoryEvent) error {
	if event.EventId == 0 {
		event.EventId = atomic.AddInt64(&s.eventIdCounter, 1)
	}
	key := s.buildExecutionHistoryKey(event.ExecutionArn, event.EventId)
	event.Timestamp = time.Now().UTC()
	return s.executionHistoryStore.Put(key, event)
}

// GetExecutionHistory retrieves the history events for an execution in
// ascending event-ID order, or in descending order when reverseOrder is
// set, paginating consistently in the requested direction.
//
// Forward order pages with the raw storage key as the marker (the keys
// zero-pad the event ID, so marker order equals event order). Reverse
// order collects the execution's events, sorts them numerically and pages
// from the newest end; execution histories are bounded in size and the
// shared BaseStore offers no descending iterator, so collecting is the
// simplest direction-correct implementation. The reverse marker anchors
// the next page to the lowest event ID already returned: a fixed anchor
// keeps pages stable while new events are appended, whereas a count from
// the newest end would shift under growth and duplicate the previous
// page's tail.
func (s *StepFunctionStore) GetExecutionHistory(ctx context.Context, executionArn string, limit int32, nextToken string, reverseOrder bool) ([]*ExecutionHistoryEvent, string, error) {
	prefix := executionArn + ":"

	if !reverseOrder {
		opts := common.ListOptions{
			Prefix:   prefix,
			Marker:   nextToken,
			MaxItems: int(limit),
		}

		result, err := common.List[ExecutionHistoryEvent](s.executionHistoryStore, opts, nil)
		if err != nil {
			return nil, "", err
		}

		return result.Items, result.NextMarker, nil
	}

	events, err := common.ListMatching[ExecutionHistoryEvent](s.executionHistoryStore, prefix, nil)
	if err != nil {
		return nil, "", err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].EventId < events[j].EventId })

	// above bounds the page at the first event at or past the anchor; the
	// page then takes the events immediately below it.
	above := len(events)
	if nextToken != "" {
		anchor, parseErr := strconv.ParseInt(nextToken, 10, 64)
		if parseErr != nil || anchor < 0 {
			return nil, "", ErrInvalidToken
		}
		above = sort.Search(len(events), func(i int) bool { return events[i].EventId >= anchor })
	}
	if above == 0 {
		return nil, "", nil
	}

	start := above - int(limit)
	if start < 0 {
		start = 0
	}
	page := events[start:above]

	reversed := make([]*ExecutionHistoryEvent, len(page))
	for i, e := range page {
		reversed[len(page)-1-i] = e
	}

	nextMarker := ""
	if start > 0 {
		nextMarker = strconv.FormatInt(events[start].EventId, 10)
	}
	return reversed, nextMarker, nil
}

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

// GetActivity retrieves an activity by its ARN.
func (s *StepFunctionStore) GetActivity(ctx context.Context, arn string) (*Activity, error) {
	var activity Activity
	if err := s.activitiesStore.Get(arn, &activity); err != nil {
		return nil, ErrActivityNotFound
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
	if ch, exists := s.activityQueues[arn]; exists {
		delete(s.activityQueues, arn)
		close(ch)
	}
	s.activityQueuesMu.Unlock()

	// Cascade delete pending tasks. Collect keys first to avoid mutating during ForEach.
	var taskKeys []string
	_ = s.tasksStore.ForEach(func(key string, value []byte) error {
		var task ActivityTask
		if err := json.Unmarshal(value, &task); err != nil {
			return nil
		}
		if task.ActivityArn == arn {
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

	return s.activitiesStore.Delete(arn)
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
	queue, exists := s.activityQueues[task.ActivityArn]
	if !exists {
		queue = make(chan *ActivityTask, 100)
		s.activityQueues[task.ActivityArn] = queue
	}

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

// ActivityTaskPollTimeout bounds how long GetActivityTask holds the
// request open waiting for a task. The Step Functions API reference fixes
// the maximum hold at 60 seconds: "The maximum time the service holds on
// to the request before responding is 60 seconds. If no task is available
// within 60 seconds, the poll returns a taskToken with an empty string."
// It is a variable so tests can shorten the wait.
var ActivityTaskPollTimeout = 60 * time.Second

// GetActivityTask retrieves a task from an activity queue, blocking until
// a task becomes available or the caller's context ends (the service
// layer bounds the wait with ActivityTaskPollTimeout).
func (s *StepFunctionStore) GetActivityTask(ctx context.Context, activityArn string, workerName string) (*ActivityTask, error) {
	s.activityQueuesMu.Lock()
	queue, exists := s.activityQueues[activityArn]
	if !exists {
		queue = make(chan *ActivityTask, 100)
		s.activityQueues[activityArn] = queue
	}
	s.activityQueuesMu.Unlock()

	select {
	case task := <-queue:
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

// GetActivityTaskByToken retrieves an activity task by its token.
func (s *StepFunctionStore) GetActivityTaskByToken(taskToken string) (*ActivityTask, error) {
	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

// HeartbeatActivityTask records the current time as the last heartbeat for
// the task. Workers call this via SendTaskHeartbeat to keep the task alive
// beyond the heartbeat interval. The actual timeout enforcement happens in
// WaitForTaskResult, which polls LastHeartbeatAt.
func (s *StepFunctionStore) HeartbeatActivityTask(taskToken string) error {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		return ErrTaskNotFound
	}

	task.LastHeartbeatAt = time.Now().UTC()

	if err := s.tasksStore.Put(taskToken, &task); err != nil {
		return err
	}

	return nil
}

// CompleteActivityTask marks an activity task as completed with output.
func (s *StepFunctionStore) CompleteActivityTask(taskToken string, output string) error {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		return ErrTaskNotFound
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
		case ch <- &ActivityTaskResult{TaskToken: taskToken, Output: output}:
		default:
		}
	}
	s.pendingTasksMu.Unlock()

	return nil
}

// FailActivityTask marks an activity task as failed with an error.
func (s *StepFunctionStore) FailActivityTask(taskToken string, errorMsg string, cause string) error {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		return ErrTaskNotFound
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
		case ch <- &ActivityTaskResult{TaskToken: taskToken, Error: fmt.Errorf("%s: %s", errorMsg, cause)}:
		default:
		}
	}
	s.pendingTasksMu.Unlock()

	return nil
}

// WaitForTaskResult waits for the result of an activity task.
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

// getActivityTaskRecord reads an activity task record by its token.
func (s *StepFunctionStore) getActivityTaskRecord(taskToken string) (*ActivityTask, error) {
	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

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
	if task, err := s.getActivityTaskRecord(taskToken); err == nil {
		switch task.Status {
		case "SUCCEEDED":
			return &ActivityTaskResult{TaskToken: taskToken, Output: task.Output}, nil
		case "FAILED":
			return &ActivityTaskResult{TaskToken: taskToken, Error: fmt.Errorf("%s: %s", task.Error, task.Cause)}, nil
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
			return result, nil
		case <-timer.C:
			s.markTaskTimedOut(taskToken)
			return nil, ErrTaskTimeout
		case <-hbTimer.C:
			var task ActivityTask
			if err := s.tasksStore.Get(taskToken, &task); err != nil {
				return nil, ErrTaskNotFound
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
		InputDetails: &ExecutionInputDetails{
			Included: input != "",
			Type:     "JSON",
		},
		StartDate: time.Now().UTC(),
	}
}

// NewExecutionHistoryEvent creates a new execution history event.
func NewExecutionHistoryEvent(executionArn string, eventType string, previousEventId int64) *ExecutionHistoryEvent {
	return &ExecutionHistoryEvent{
		ExecutionArn:    executionArn,
		Type:            eventType,
		PreviousEventId: previousEventId,
		Timestamp:       time.Now().UTC(),
	}
}

// RegisterExecution registers an active execution with its cancel function.
func (s *StepFunctionStore) RegisterExecution(executionArn string, cancel context.CancelFunc) {
	s.activeExecutionsMu.Lock()
	s.activeExecutions[executionArn] = cancel
	s.activeExecutionsMu.Unlock()
}

// CancelExecution cancels a running execution.
func (s *StepFunctionStore) CancelExecution(executionArn string) bool {
	s.activeExecutionsMu.Lock()
	cancel, exists := s.activeExecutions[executionArn]
	if exists {
		delete(s.activeExecutions, executionArn)
	}
	s.activeExecutionsMu.Unlock()

	if exists && cancel != nil {
		cancel()
		return true
	}
	return false
}

// UnregisterExecution removes an execution from the active executions list.
func (s *StepFunctionStore) UnregisterExecution(executionArn string) {
	s.activeExecutionsMu.Lock()
	delete(s.activeExecutions, executionArn)
	s.activeExecutionsMu.Unlock()
}

// CancelAllExecutions cancels all running executions.
func (s *StepFunctionStore) CancelAllExecutions() {
	s.activeExecutionsMu.Lock()
	for _, cancel := range s.activeExecutions {
		cancel()
	}
	s.activeExecutionsMu.Unlock()
}

func (s *StepFunctionStore) buildVersionARN(smArn string, version int64) string {
	return smArn + fmt.Sprintf(":%d", version)
}

func (s *StepFunctionStore) buildAliasARN(smArn, aliasName string) string {
	// The alias ARN extends the state machine ARN with the alias name, so
	// alias names are namespaced per state machine rather than per account.
	return smArn + ":" + aliasName
}

func (s *StepFunctionStore) nextVersionNumber(smArn string) int64 {
	s.versionCountersMu.Lock()
	defer s.versionCountersMu.Unlock()
	s.versionCounters[smArn]++
	return s.versionCounters[smArn]
}

func (s *StepFunctionStore) recoverVersionCounter(smArn string) {
	versions, err := common.ListMatching[StateMachineVersion](s.versionsStore, smArn+":", nil)
	if err != nil {
		return
	}
	var maxVersion int64
	for _, v := range versions {
		if v.Version > maxVersion {
			maxVersion = v.Version
		}
	}
	s.versionCountersMu.Lock()
	if maxVersion > s.versionCounters[smArn] {
		s.versionCounters[smArn] = maxVersion
	}
	s.versionCountersMu.Unlock()
}

// PublishStateMachineVersion publishes the state machine's current
// revision as a version. Publishing is idempotent per revision: when a
// version of the current revision already exists it is returned instead of
// publishing a duplicate.
func (s *StepFunctionStore) PublishStateMachineVersion(ctx context.Context, smArn string, description string) (*StateMachineVersion, error) {
	sm, err := s.GetStateMachine(ctx, smArn)
	if err != nil {
		return nil, ErrStateMachineNotFound
	}

	if existing, err := s.FindVersionByRevision(smArn, sm.RevisionId); err == nil {
		return existing, nil
	} else if !errors.Is(err, ErrStateMachineVersionNotFound) {
		return nil, err
	}

	s.versionCountersMu.Lock()
	if _, exists := s.versionCounters[smArn]; !exists {
		s.versionCountersMu.Unlock()
		s.recoverVersionCounter(smArn)
	} else {
		s.versionCountersMu.Unlock()
	}

	version := s.nextVersionNumber(smArn)
	versionArn := s.buildVersionARN(smArn, version)

	v := &StateMachineVersion{
		StateMachineVersionArn: versionArn,
		StateMachineArn:        smArn,
		Version:                version,
		Description:            description,
		CreationDate:           time.Now().UTC(),
		Definition:             sm.Definition,
		RevisionId:             sm.RevisionId,
	}

	if err := s.versionsStore.Put(versionArn, v); err != nil {
		return nil, err
	}

	return v, nil
}

// FindVersionByRevision returns the version published for the given state
// machine revision, or ErrStateMachineVersionNotFound when the revision has
// never been published. When several records claim the same revision the
// highest version number wins, keeping the result deterministic.
func (s *StepFunctionStore) FindVersionByRevision(smArn, revisionId string) (*StateMachineVersion, error) {
	versions, err := common.ListMatching[StateMachineVersion](s.versionsStore, smArn+":", func(v *StateMachineVersion) bool {
		return v.RevisionId == revisionId
	})
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, ErrStateMachineVersionNotFound
	}
	newest := versions[0]
	for _, v := range versions[1:] {
		if v.Version > newest.Version {
			newest = v
		}
	}
	return newest, nil
}

// CountStateMachineVersions returns the number of published versions of a
// state machine; the quota is one thousand versions per state machine.
func (s *StepFunctionStore) CountStateMachineVersions(smArn string) (int, error) {
	versions, err := common.ListMatching[StateMachineVersion](s.versionsStore, smArn+":", nil)
	if err != nil {
		return 0, err
	}
	return len(versions), nil
}

// GetStateMachineVersion retrieves a state machine version by its ARN.
func (s *StepFunctionStore) GetStateMachineVersion(ctx context.Context, arn string) (*StateMachineVersion, error) {
	var v StateMachineVersion
	if err := s.versionsStore.Get(arn, &v); err != nil {
		return nil, ErrStateMachineVersionNotFound
	}
	return &v, nil
}

// GetStateMachineVersionByNumber retrieves a version of a state machine by
// its sequential version number.
func (s *StepFunctionStore) GetStateMachineVersionByNumber(ctx context.Context, smArn string, version int64) (*StateMachineVersion, error) {
	return s.GetStateMachineVersion(ctx, s.buildVersionARN(smArn, version))
}

// DeleteStateMachineVersion removes a state machine version from the store.
// A version an alias routing configuration still points at cannot be
// deleted; the caller maps the sentinel to ConflictException.
func (s *StepFunctionStore) DeleteStateMachineVersion(ctx context.Context, arn string) error {
	if !s.versionsStore.Exists(arn) {
		return ErrStateMachineVersionNotFound
	}
	aliases, err := common.ListMatching[StateMachineAlias](s.aliasesStore, "", func(a *StateMachineAlias) bool {
		for _, rc := range a.RoutingConfiguration {
			if rc.StateMachineVersionArn == arn {
				return true
			}
		}
		return false
	})
	if err != nil {
		return err
	}
	if len(aliases) > 0 {
		return ErrStateMachineVersionInUse
	}
	return s.versionsStore.Delete(arn)
}

// ListStateMachineVersions returns a paginated list of versions for a state machine.
func (s *StepFunctionStore) ListStateMachineVersions(ctx context.Context, smArn string, limit int32, nextToken string) (*StateMachineVersionListResult, error) {
	opts := common.ListOptions{
		Prefix:   smArn + ":",
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[StateMachineVersion](s.versionsStore, opts, nil)
	if err != nil {
		return nil, err
	}

	return &StateMachineVersionListResult{
		Versions:  result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// CreateStateMachineAlias creates a new alias for a state machine.
func (s *StepFunctionStore) CreateStateMachineAlias(ctx context.Context, alias *StateMachineAlias) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	if alias.Name == "" {
		return ErrInvalidARN
	}

	aliasArn := s.buildAliasARN(alias.StateMachineArn, alias.Name)
	storageKey := aliasArn
	if s.aliasesStore.Exists(storageKey) {
		return ErrStateMachineAliasAlreadyExists
	}

	now := time.Now().UTC()
	alias.StateMachineAliasArn = aliasArn
	alias.CreationDate = now
	alias.UpdateDate = now

	return s.aliasesStore.Put(storageKey, alias)
}

// GetStateMachineAlias retrieves a state machine alias by its ARN.
func (s *StepFunctionStore) GetStateMachineAlias(ctx context.Context, arn string) (*StateMachineAlias, error) {
	var alias StateMachineAlias
	if err := s.aliasesStore.Get(arn, &alias); err != nil {
		return nil, ErrStateMachineAliasNotFound
	}
	return &alias, nil
}

// GetStateMachineAliasByName retrieves a state machine alias by its state
// machine ARN and alias name.
func (s *StepFunctionStore) GetStateMachineAliasByName(ctx context.Context, smArn, name string) (*StateMachineAlias, error) {
	return s.GetStateMachineAlias(ctx, s.buildAliasARN(smArn, name))
}

// UpdateStateMachineAlias updates an existing state machine alias.
func (s *StepFunctionStore) UpdateStateMachineAlias(ctx context.Context, alias *StateMachineAlias) error {
	if !s.aliasesStore.Exists(alias.StateMachineAliasArn) {
		return ErrStateMachineAliasNotFound
	}
	alias.UpdateDate = time.Now().UTC()
	return s.aliasesStore.Put(alias.StateMachineAliasArn, alias)
}

// DeleteStateMachineAlias removes a state machine alias from the store.
func (s *StepFunctionStore) DeleteStateMachineAlias(ctx context.Context, arn string) error {
	if !s.aliasesStore.Exists(arn) {
		return ErrStateMachineAliasNotFound
	}
	return s.aliasesStore.Delete(arn)
}

// ListStateMachineAliases returns a paginated list of aliases for a state machine.
func (s *StepFunctionStore) ListStateMachineAliases(ctx context.Context, smArn string, limit int32, nextToken string) (*StateMachineAliasListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[StateMachineAlias](s.aliasesStore, opts, func(a *StateMachineAlias) bool {
		return a.StateMachineArn == smArn
	})
	if err != nil {
		return nil, err
	}

	return &StateMachineAliasListResult{
		Aliases:   result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// CountStateMachineAliases returns the number of aliases defined for a
// state machine; the quota is one hundred aliases per state machine.
func (s *StepFunctionStore) CountStateMachineAliases(smArn string) (int, error) {
	aliases, err := common.ListMatching[StateMachineAlias](s.aliasesStore, "", func(a *StateMachineAlias) bool {
		return a.StateMachineArn == smArn
	})
	if err != nil {
		return 0, err
	}
	return len(aliases), nil
}

func (s *StepFunctionStore) nextMapRunSeq() int64 {
	if atomic.LoadInt64(&s.mapRunSeq) == 0 {
		s.recoverMapRunSeq()
	}
	return atomic.AddInt64(&s.mapRunSeq, 1)
}

// NextMapRunSeq returns the next sequential identifier for a map run.
func (s *StepFunctionStore) NextMapRunSeq() int64 {
	return s.nextMapRunSeq()
}

func (s *StepFunctionStore) recoverMapRunSeq() {
	mapRuns, err := common.ListMatching[MapRun](s.mapRunsStore, "", nil)
	if err != nil {
		return
	}
	var maxSeq int64
	for _, mr := range mapRuns {
		arn := mr.MapRunArn
		if idx := strings.LastIndex(arn, "/mapRun-"); idx >= 0 {
			rest := arn[idx+7:]
			if spaceIdx := strings.IndexByte(rest, '-'); spaceIdx > 0 {
				if n, err := strconv.ParseInt(rest[:spaceIdx], 10, 64); err == nil && n > maxSeq {
					maxSeq = n
				}
			}
		}
	}
	if maxSeq > s.mapRunSeq {
		atomic.StoreInt64(&s.mapRunSeq, maxSeq)
	}
}

// CreateMapRun persists a new map run in Pebble-backed storage.
func (s *StepFunctionStore) CreateMapRun(ctx context.Context, mr *MapRun) error {
	return s.mapRunsStore.Put(mr.MapRunArn, mr)
}

// UpdateMapRun persists an updated map run record.
func (s *StepFunctionStore) UpdateMapRun(ctx context.Context, mr *MapRun) error {
	return s.mapRunsStore.Put(mr.MapRunArn, mr)
}

// GetMapRun retrieves a map run by its ARN.
func (s *StepFunctionStore) GetMapRun(ctx context.Context, mapRunArn string) (*MapRun, error) {
	var mr MapRun
	if err := s.mapRunsStore.Get(mapRunArn, &mr); err != nil {
		return nil, ErrMapRunNotFound
	}
	return &mr, nil
}

// ListMapRunsByExecution returns all map runs for a given execution ARN.
func (s *StepFunctionStore) ListMapRunsByExecution(ctx context.Context, executionArn string) ([]*MapRun, error) {
	return common.ListMatching[MapRun](s.mapRunsStore, "", func(mr *MapRun) bool {
		return mr.ExecutionArn == executionArn
	})
}

// ListAllMapRuns returns all map runs, optionally filtered by execution ARN.
func (s *StepFunctionStore) ListAllMapRuns(ctx context.Context, executionArn string, limit int32, nextToken string) (*MapRunListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[MapRun](s.mapRunsStore, opts, func(mr *MapRun) bool {
		if executionArn != "" && mr.ExecutionArn != executionArn {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &MapRunListResult{
		MapRuns:   result.Items,
		NextToken: result.NextMarker,
	}, nil
}
