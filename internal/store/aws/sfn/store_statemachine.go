// This file holds the state-machine CRUD surface and the deletion cascade:
// removing a state machine drops its executions, histories, map runs,
// versions, aliases and tags.

// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

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

// GetStateMachine retrieves a state machine by its ARN. A missing record is
// the not-found sentinel; any other storage fault surfaces as itself so the
// Cores can tell a permanent absence from a transient failure.
func (s *StepFunctionStore) GetStateMachine(ctx context.Context, arn string) (*StateMachine, error) {
	var sm StateMachine
	if err := s.BaseStore.Get(arn, &sm); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrStateMachineNotFound
		}
		return nil, err
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
	for execArn, reg := range s.activeExecutions {
		execSmName := extractStateMachineNameFromExecutionArn(execArn)
		if execSmName == smName {
			reg.cancel()
			delete(s.activeExecutions, execArn)
		}
	}
	s.activeExecutionsMu.Unlock()

	// Foreign waiters survive this machine's deletion: dropping them here
	// would strand the waiter for the rest of its task timeout even after
	// its worker reports. Deciding foreign-versus-own and closing within a
	// single lock hold leaves no window in which a report for a foreign
	// token finds no registered waiter and drops its result.
	s.pendingTasksMu.Lock()
	kept := make(map[string]chan *ActivityTaskResult)
	for token, ch := range s.pendingTasks {
		var task ActivityTask
		foreign := false
		if err := s.tasksStore.Get(token, &task); err == nil {
			foreign = extractStateMachineNameFromExecutionArn(task.ExecutionArn) != smName
		}
		if foreign {
			kept[token] = ch
			continue
		}
		close(ch)
	}
	s.pendingTasks = kept
	s.pendingTasksMu.Unlock()

	// Activity queues are NOT closed here. Activities are independent
	// resources in AWS Step Functions; closing their queues when an
	// unrelated state machine is deleted would cause GetActivityTask to
	// receive nil from a closed channel, triggering a nil-pointer panic.
	// Pending tasks belonging to this state machine are already cleaned
	// up above via the pendingTasks close logic.

	// The execution cascade runs under the creators' mutex: CreateExecution
	// checks existence and writes under createMu, so a StartExecution
	// racing this deletion could otherwise land its execution after the
	// cascade scan — a permanent orphan under a machine that no longer
	// exists, which the restart recovery would then fail on every boot.
	s.createMu.Lock()
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
	s.createMu.Unlock()

	// Cascade delete the machine's activity-task records: the executions
	// are gone, so GetActivityTaskByToken and SendTaskHeartbeat must not
	// resolve tasks whose execution no longer exists — the same cascade
	// DeleteActivity runs for its own tasks.
	s.deleteActivityTasks(func(task *ActivityTask) bool {
		return extractStateMachineNameFromExecutionArn(task.ExecutionArn) == smName
	})

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
	// The version counter entry goes with the versions it counted: leaving
	// the stale high-water mark would make a same-name machine's first
	// publish resume at staleMax+1 instead of recovering from the persisted
	// records (which are now empty).
	s.versionCountersMu.Lock()
	delete(s.versionCounters, arn)
	s.versionCountersMu.Unlock()

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

	// Cascade delete map runs of the machine's executions: runs are keyed
	// by their own ARN, so they otherwise survive as orphans that
	// DescribeMapRun still resolves and unfiltered listings still return.
	var mapRunKeys []string
	_ = s.mapRunsStore.ForEach(func(key string, value []byte) error {
		var mr MapRun
		if err := json.Unmarshal(value, &mr); err != nil {
			return nil
		}
		if mr.StateMachineArn == arn {
			mapRunKeys = append(mapRunKeys, key)
		}
		return nil
	})
	for _, key := range mapRunKeys {
		_ = s.mapRunsStore.Delete(key)
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

// deleteExecutionHistory drops an execution's history events, keyed by the
// execution ARN prefix in the dedicated history bucket.
func (s *StepFunctionStore) deleteExecutionHistory(executionArn string) {
	_ = s.executionHistoryStore.DeleteByPrefix(executionArn + ":")
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
