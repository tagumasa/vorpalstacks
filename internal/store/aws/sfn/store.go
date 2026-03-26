// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"context"
	"fmt"
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
	activeExecutions   map[string]context.CancelFunc
	activeExecutionsMu sync.RWMutex
	createMu           sync.Mutex
}

// NewStepFunctionStore creates a new StepFunctionStore instance.
func NewStepFunctionStore(store storage.BasicStorage, accountID, region string) *StepFunctionStore {
	return &StepFunctionStore{
		BaseStore:             common.NewBaseStore(store.Bucket("stepfunction-statemachines-"+region), "stepfunction-statemachines"),
		executionsStore:       common.NewBaseStore(store.Bucket("stepfunction-executions-"+region), "stepfunction-executions"),
		executionHistoryStore: common.NewBaseStore(store.Bucket("stepfunction-history-"+region), "stepfunction-history"),
		activitiesStore:       common.NewBaseStore(store.Bucket("stepfunction-activities-"+region), "stepfunction-activities"),
		tasksStore:            common.NewBaseStore(store.Bucket("stepfunction-tasks-"+region), "stepfunction-tasks"),
		TagStore:              common.NewTagStoreWithRegion(store, "stepfunction", region),
		arnBuilder:            svcarn.NewARNBuilder(accountID, region),
		accountID:             accountID,
		region:                region,
		pendingTasks:          make(map[string]chan *ActivityTaskResult),
		activityQueues:        make(map[string]chan *ActivityTask),
		activeExecutions:      make(map[string]context.CancelFunc),
	}
}

// GetAccountID returns the AWS account ID.
func (s *StepFunctionStore) GetAccountID() string { return s.accountID }

// GetRegion returns the AWS region.
func (s *StepFunctionStore) GetRegion() string { return s.region }

func (s *StepFunctionStore) buildStateMachineARN(name string) string {
	return s.arnBuilder.StepFunctions().StateMachine(name)
}

func (s *StepFunctionStore) buildExecutionARN(stateMachineName, executionName string) string {
	return s.arnBuilder.StepFunctions().Execution(stateMachineName, executionName)
}

func (s *StepFunctionStore) buildActivityARN(name string) string {
	return s.arnBuilder.StepFunctions().Activity(name)
}

func (s *StepFunctionStore) buildExecutionHistoryKey(executionArn string, eventId int64) string {
	return fmt.Sprintf("%s:%d", executionArn, eventId)
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

// DeleteStateMachine removes a state machine from the store.
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
	for _, ch := range oldPending {
		close(ch)
	}

	s.activityQueuesMu.Lock()
	for arn, ch := range s.activityQueues {
	drainLoop:
		for {
			select {
			case <-ch:
			default:
				delete(s.activityQueues, arn)
				break drainLoop
			}
		}
	}
	s.activityQueuesMu.Unlock()

	return s.BaseStore.Delete(arn)
}

func extractStateMachineNameFromArn(arn string) string {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	if strings.HasPrefix(resource, "stateMachine:") {
		return strings.TrimPrefix(resource, "stateMachine:")
	}
	return ""
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

// ListExecutions returns a paginated list of executions for a state machine.
func (s *StepFunctionStore) ListExecutions(ctx context.Context, stateMachineArn string, statusFilter string, limit int32, nextToken string) (*ExecutionListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Execution](s.executionsStore, opts, func(e *Execution) bool {
		if stateMachineArn != "" && e.StateMachineArn != stateMachineArn {
			return false
		}
		if statusFilter != "" && e.Status != statusFilter {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &ExecutionListResult{
		Executions: result.Items,
		NextToken:  result.NextMarker,
	}, nil
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

// GetExecutionHistory retrieves the history events for an execution.
func (s *StepFunctionStore) GetExecutionHistory(ctx context.Context, executionArn string, limit int32, nextToken string) ([]*ExecutionHistoryEvent, string, error) {
	prefix := executionArn + ":"
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

// DeleteActivity removes an activity from the store.
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

// CreateActivityTask creates a new task for an activity.
func (s *StepFunctionStore) CreateActivityTask(task *ActivityTask) error {
	task.TaskToken = uuid.New().String()
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
	s.activityQueuesMu.Unlock()

	select {
	case queue <- task:
		return nil
	default:
		_ = s.tasksStore.Delete(task.TaskToken)
		return ErrActivityQueueFull
	}
}

// GetActivityTask retrieves a task from an activity queue.
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

// CompleteActivityTask marks an activity task as completed with output.
func (s *StepFunctionStore) CompleteActivityTask(taskToken string, output string) error {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	var task ActivityTask
	if err := s.tasksStore.Get(taskToken, &task); err != nil {
		return ErrTaskNotFound
	}

	task.Status = "SUCCEEDED"
	task.Output = output
	task.CompletedAt = time.Now().UTC()

	if err := s.tasksStore.Put(taskToken, &task); err != nil {
		return err
	}

	s.pendingTasksMu.Lock()
	if ch, ok := s.pendingTasks[taskToken]; ok {
		ch <- &ActivityTaskResult{TaskToken: taskToken, Output: output}
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

	task.Status = "FAILED"
	task.Error = errorMsg
	task.Cause = cause
	task.CompletedAt = time.Now().UTC()

	if err := s.tasksStore.Put(taskToken, &task); err != nil {
		return err
	}

	s.pendingTasksMu.Lock()
	if ch, ok := s.pendingTasks[taskToken]; ok {
		ch <- &ActivityTaskResult{TaskToken: taskToken, Error: fmt.Errorf("%s: %s", errorMsg, cause)}
	}
	s.pendingTasksMu.Unlock()

	return nil
}

// WaitForTaskResult waits for the result of an activity task.
func (s *StepFunctionStore) WaitForTaskResult(ctx context.Context, taskToken string, timeout time.Duration) (*ActivityTaskResult, error) {
	s.pendingTasksMu.Lock()
	ch := make(chan *ActivityTaskResult, 1)
	s.pendingTasks[taskToken] = ch
	s.pendingTasksMu.Unlock()

	defer func() {
		s.pendingTasksMu.Lock()
		delete(s.pendingTasks, taskToken)
		s.pendingTasksMu.Unlock()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case result := <-ch:
		return result, nil
	case <-timer.C:
		return nil, ErrTaskTimeout
	case <-ctx.Done():
		return nil, ctx.Err()
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
