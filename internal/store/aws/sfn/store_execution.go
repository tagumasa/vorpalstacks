// This file holds the execution CRUD surface, the execution history
// reads and writes, and the active-execution registration that backs
// StopExecution cancellation and redrive draining.

// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

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

	return s.executionsStore.Put(exec.ExecutionArn, exec)
}

// GetExecution retrieves an execution by ARN. A missing record is the
// not-found sentinel; any other storage fault surfaces as itself, so the
// Cores can tell a permanent absence from a transient failure instead of
// answering both with the does-not-exist verdict.
func (s *StepFunctionStore) GetExecution(ctx context.Context, arn string) (*Execution, error) {
	var exec Execution
	if err := s.executionsStore.Get(arn, &exec); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrExecutionNotFound
		}
		return nil, err
	}
	return &exec, nil
}

// UpdateExecution updates an existing execution in the store.
func (s *StepFunctionStore) UpdateExecution(ctx context.Context, exec *Execution) error {
	return s.executionsStore.Put(exec.ExecutionArn, exec)
}

// TransitionExecutionForRedrive applies a redrive transition under the
// creation mutex: the execution is re-read inside the critical section, guard
// decides against the fresh record whether the transition may proceed, and
// apply mutates that fresh record before the single write. Serialising the
// re-read and the write is what makes the transition atomic — two concurrent
// redrives of the same execution cannot both observe the terminal status, so
// exactly one of them transitions it. guard and apply run with createMu held
// and must not call back into mutating store methods.
func (s *StepFunctionStore) TransitionExecutionForRedrive(ctx context.Context, arn string, guard func(fresh *Execution) bool, apply func(fresh *Execution)) (*Execution, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	var fresh Execution
	if err := s.executionsStore.Get(arn, &fresh); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrExecutionNotFound
		}
		return nil, err
	}
	if !guard(&fresh) {
		return nil, ErrExecutionNotRedrivable
	}
	apply(&fresh)
	if err := s.executionsStore.Put(arn, &fresh); err != nil {
		return nil, err
	}
	return &fresh, nil
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
// History keys are "<arn>:<zero-padded event id>", so an own event's key
// is exactly the "<arn>:" prefix followed by digits — a Distributed Map
// child ("parent:M-0:<id>") or a colon-bearing sibling execution name
// contributes a non-digit tail after the prefix and stays excluded, the
// same scope executionHistoryFilter enforces. Counting from the key alone
// needs no record deserialisation, and the call sits on the
// DescribeExecution hot path (redrive-status derivation), so the
// iteration is bounded to the execution's own key prefix rather than
// visiting every key in the history bucket.
func (s *StepFunctionStore) CountExecutionHistory(ctx context.Context, executionArn string) (int, error) {
	prefix := executionArn + ":"
	count := 0
	err := s.executionHistoryStore.ScanPrefix(prefix, func(key string, _ []byte) error {
		rest := strings.TrimPrefix(key, prefix)
		if rest == "" {
			return nil
		}
		for i := 0; i < len(rest); i++ {
			if rest[i] < '0' || rest[i] > '9' {
				return nil
			}
		}
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// LastExecutionHistoryId returns the highest event id recorded for an
// execution, or zero when the history is empty. Resume paths continue the
// id sequence from it. The walk streams the execution's own key prefix —
// bounded iteration, no materialised event slice.
func (s *StepFunctionStore) LastExecutionHistoryId(ctx context.Context, executionArn string) (int64, error) {
	match := executionHistoryFilter(executionArn)
	var max int64
	err := s.executionHistoryStore.ScanPrefix(executionArn+":", func(key string, value []byte) error {
		var event ExecutionHistoryEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		if match(&event) && event.EventId > max {
			max = event.EventId
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return max, nil
}

// executionHistoryFilter matches only the named execution's events:
// Distributed Map child ARNs nest under their parent ("parent:M-0"), so
// the key prefix alone would also return the children's events.
func executionHistoryFilter(executionArn string) func(*ExecutionHistoryEvent) bool {
	return func(e *ExecutionHistoryEvent) bool {
		return e.ExecutionArn == executionArn
	}
}

// AddExecutionHistoryEvent adds a history event to an execution's event log.
// The caller assigns the event ID: IDs are per-execution and sequential, and
// every production caller numbers events through the executor's counter.
func (s *StepFunctionStore) AddExecutionHistoryEvent(ctx context.Context, event *ExecutionHistoryEvent) error {
	key := s.buildExecutionHistoryKey(event.ExecutionArn, event.EventId)
	event.Timestamp = time.Now().UTC()
	return s.executionHistoryStore.Put(key, event)
}

// errHistoryPageFull stops the descending history walk once a page has
// filled; it is a control-flow signal, not a storage failure.
var errHistoryPageFull = errors.New("execution history page full")

// GetExecutionHistory retrieves the history events for an execution in
// ascending event-ID order, or in descending order when reverseOrder is
// set, paginating consistently in the requested direction.
//
// Forward order pages with the raw storage key as the marker (the keys
// zero-pad the event ID, so marker order equals event order). Reverse
// order walks the keys with the descending iterator, dropping Distributed
// Map child events whose ARNs nest under the parent's key prefix. The
// reverse marker anchors the next page at the lowest event ID already
// returned: a fixed anchor keeps pages stable while new events are
// appended, whereas a count from the newest end would shift under growth
// and duplicate the previous page's tail.
func (s *StepFunctionStore) GetExecutionHistory(ctx context.Context, executionArn string, limit int32, nextToken string, reverseOrder bool) ([]*ExecutionHistoryEvent, string, error) {
	prefix := executionArn + ":"

	if !reverseOrder {
		opts := common.ListOptions{
			Prefix:   prefix,
			Marker:   nextToken,
			MaxItems: int(limit),
		}

		result, err := common.List[ExecutionHistoryEvent](s.executionHistoryStore, opts, executionHistoryFilter(executionArn))
		if err != nil {
			return nil, "", err
		}

		return result.Items, result.NextMarker, nil
	}

	before := ""
	if nextToken != "" {
		anchor, parseErr := strconv.ParseInt(nextToken, 10, 64)
		if parseErr != nil || anchor < 0 {
			return nil, "", ErrInvalidToken
		}
		before = s.buildExecutionHistoryKey(executionArn, anchor)
	}

	var page []*ExecutionHistoryEvent
	limitReached := false
	visit := func(key string, value []byte) error {
		if len(page) >= int(limit) {
			limitReached = true
			return errHistoryPageFull
		}
		var event ExecutionHistoryEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return nil
		}
		if event.ExecutionArn != executionArn {
			return nil
		}
		page = append(page, &event)
		return nil
	}
	if err := s.executionHistoryStore.ScanPrefixReverse(prefix, before, visit); err != nil && !errors.Is(err, errHistoryPageFull) {
		return nil, "", err
	}
	if len(page) == 0 {
		return nil, "", nil
	}

	// The marker names the lowest event of this page; the next page walks
	// strictly below it. It is emitted only when the limit stopped the
	// walk, so a page that ends at the oldest event terminates pagination.
	nextMarker := ""
	if limitReached {
		nextMarker = strconv.FormatInt(page[len(page)-1].EventId, 10)
	}
	return page, nextMarker, nil
}

// RegisterExecution registers an active execution with its cancel function
// and returns the registration handle. The handle is channel-identity
// guarded: a superseding registration for the same ARN replaces the entry,
// and only the matching handle's UnregisterExecution removes it — a
// draining goroutine's unregister never deletes its successor's claim.
func (s *StepFunctionStore) RegisterExecution(executionArn string, cancel context.CancelFunc) ExecutionHandle {
	s.activeExecutionsMu.Lock()
	defer s.activeExecutionsMu.Unlock()
	reg := &executionRegistration{cancel: cancel, done: make(chan struct{})}
	s.activeExecutions[executionArn] = reg
	return ExecutionHandle{done: reg.done}
}

// CancelExecution cancels a running execution. The registration stays in
// place: the goroutine's own UnregisterExecution removes it once every
// write has drained, which is what keeps WaitExecutionInactive meaningful
// for a stopped execution.
func (s *StepFunctionStore) CancelExecution(executionArn string) bool {
	s.activeExecutionsMu.Lock()
	reg, exists := s.activeExecutions[executionArn]
	s.activeExecutionsMu.Unlock()
	if exists && reg.cancel != nil {
		reg.cancel()
		return true
	}
	return false
}

// UnregisterExecution releases the execution registration the handle
// identifies. The done channel closes after the entry is observable as
// removed (or after a superseding registration replaced it), so a waiter
// on the channel knows the goroutine's writes have landed.
func (s *StepFunctionStore) UnregisterExecution(executionArn string, handle ExecutionHandle) {
	s.activeExecutionsMu.Lock()
	if reg, exists := s.activeExecutions[executionArn]; exists && reg.done == handle.done {
		delete(s.activeExecutions, executionArn)
	}
	s.activeExecutionsMu.Unlock()
	close(handle.done)
}

// WaitExecutionInactive blocks until no goroutine holds an execution
// registration, or the timeout passes (reporting whether the drain
// completed). A registration spans every history and record write its
// goroutine makes, so once this returns true, numbering new events from
// the store's current maximum is stable — the redrive path relies on
// exactly that ordering.
func (s *StepFunctionStore) WaitExecutionInactive(executionArn string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		s.activeExecutionsMu.RLock()
		reg, exists := s.activeExecutions[executionArn]
		s.activeExecutionsMu.RUnlock()
		if !exists {
			return true
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return false
		}
		select {
		case <-reg.done:
			return true
		case <-time.After(remaining):
			return false
		}
	}
}

// CancelAllExecutions cancels all running executions.
func (s *StepFunctionStore) CancelAllExecutions() {
	s.activeExecutionsMu.Lock()
	for _, reg := range s.activeExecutions {
		reg.cancel()
	}
	s.activeExecutionsMu.Unlock()
}
