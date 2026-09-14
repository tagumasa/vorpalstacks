package sfn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	"vorpalstacks/internal/utils/aws/arn"
)

// Executor manages the execution of Step Functions state machines.
type Executor struct {
	store               *sfnstore.StepFunctionStore
	accountID           string
	region              string
	bus                 eventbus.ServiceBus
	taskAuthz           TaskCredentialsAuthoriser
	currentRoleArn      string
	currentExecution    *sfnstore.Execution
	currentStateMachine *sfnstore.StateMachine
	// noPersistHistory keeps logged events out of the store: TestState
	// "executes a state without creating a state machine", so no execution
	// resource exists and no history record may.
	noPersistHistory bool
}

// generateTaskToken returns an unguessable task token. Anyone holding the
// token can report a task result, so it must come from crypto/rand rather
// than a clock, which is both guessable and collision-prone.
func generateTaskToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Unreachable in practice; uuid v4 is also crypto/rand-backed.
		return uuid.New().String()
	}
	return hex.EncodeToString(b)
}

// NewExecutor creates a new Step Functions executor with the given store and event bus.
func NewExecutor(store *sfnstore.StepFunctionStore, bus eventbus.ServiceBus) *Executor {
	return &Executor{
		store: store,
		bus:   bus,
	}
}

// NewExecutorWithStores creates a new Step Functions executor with all dependencies.
// taskAuthz is the Task Credentials authorisation seam; nil leaves task
// invocations on the machine-role behaviour.
func NewExecutorWithStores(store *sfnstore.StepFunctionStore, bus eventbus.ServiceBus, accountID, region string, taskAuthz TaskCredentialsAuthoriser) *Executor {
	return &Executor{
		store:     store,
		bus:       bus,
		accountID: accountID,
		region:    region,
		taskAuthz: taskAuthz,
	}
}

// ExecuteStateMachine executes a state machine with the given execution context.
func (e *Executor) ExecuteStateMachine(ctx context.Context, execution *sfnstore.Execution) error {
	e.currentExecution = execution
	sm, err := e.store.GetStateMachine(ctx, execution.StateMachineArn)
	if err == nil && sm != nil {
		e.currentRoleArn = sm.RoleArn
		e.currentStateMachine = sm
	}

	definition, states, err := e.parseDefinitionForExecution(ctx, execution)
	if err != nil {
		execution.Status = "FAILED"
		execution.Error = "InvalidDefinition"
		logs.Error("Failed to parse state machine definition", logs.String("arn", execution.StateMachineArn), logs.Err(err))
		execution.Cause = "Invalid state machine definition syntax"
		execution.StopDate = time.Now().UTC()
		if updateErr := e.updateExecutionWithRetry(ctx, execution); updateErr != nil {
			logs.Error("Failed to update execution status to FAILED after definition error", logs.Err(updateErr))
		}
		return fmt.Errorf("failed to parse state machine definition: %w", err)
	}

	if definition.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(definition.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	eventId := int64(1)
	roleArn := ""
	if e.currentStateMachine != nil {
		roleArn = e.currentStateMachine.RoleArn
	}
	err = e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: execution.ExecutionArn,
		EventId:      eventId,
		Type:         "ExecutionStarted",
		Timestamp:    time.Now().UTC(),
		ExecutionStartedEventDetails: &sfnstore.ExecutionStartedEventDetails{
			Input:                  execution.Input,
			RoleArn:                roleArn,
			StateMachineAliasArn:   execution.StateMachineAliasArn,
			StateMachineVersionArn: execution.StateMachineVersionArn,
		},
	})
	if err != nil {
		logs.Error("Failed to add ExecutionStarted event", logs.Err(err))
	}

	execCtx := &ExecutionContext{
		Execution:     execution,
		Definition:    definition,
		CurrentState:  definition.StartAt,
		Input:         execution.Input,
		Output:        "",
		EventId:       &eventId,
		States:        states,
		QueryLanguage: definition.QueryLanguage,
		VariableScope: NewVariableScope(nil),
		MapItemIndex:  -1,
	}

	err = e.executeStates(ctx, execCtx)
	e.finalizeExecution(ctx, execution, execCtx, err)
	if err != nil && execution.Status != "ABORTED" {
		return err
	}
	return nil
}

// ExecuteStateMachineFromState resumes a previously-interrupted execution
// starting from the given state. It does not emit an ExecutionStarted
// event (that already exists from the original attempt); the caller
// records the ExecutionRedriven event when the resume is a user-initiated
// redrive, and the event history continues from lastEventId. The
// IsRedrive flag on the ExecutionContext signals to Map and Parallel
// states that they should consult stored checkpoints for
// already-completed iterations or branches. armMachineTimeout re-arms the
// state-machine-level timeout for restart recovery only: the documented
// redrive contract resets it ("When you redrive an execution, the state
// machine level timeout, if defined, is reset to 0"), so a redrive passes
// false and runs without it.
func (e *Executor) ExecuteStateMachineFromState(ctx context.Context, execution *sfnstore.Execution, startState string, startInput string, lastEventId int64, armMachineTimeout bool) error {
	e.currentExecution = execution
	sm, err := e.store.GetStateMachine(ctx, execution.StateMachineArn)
	if err == nil && sm != nil {
		e.currentRoleArn = sm.RoleArn
		e.currentStateMachine = sm
	}

	definition, states, err := e.parseDefinitionForExecution(ctx, execution)
	if err != nil {
		execution.Status = "FAILED"
		execution.Error = "InvalidDefinition"
		execution.Cause = "Invalid state machine definition syntax"
		execution.StopDate = time.Now().UTC()
		if updateErr := e.updateExecutionWithRetry(ctx, execution); updateErr != nil {
			logs.Error("Failed to update execution status to FAILED after definition error", logs.Err(updateErr))
		}
		return fmt.Errorf("failed to parse state machine definition: %w", err)
	}

	if armMachineTimeout && definition.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(definition.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	// The counter is add-and-fetch, so it starts at the last id the history
	// already holds: the first state event after the resume takes
	// lastEventId+1 and the id sequence stays gapless.
	eventId := lastEventId

	resumeInput := startInput
	if resumeInput == "" {
		resumeInput = execution.Input
	}

	execCtx := &ExecutionContext{
		Execution:     execution,
		Definition:    definition,
		CurrentState:  startState,
		Input:         resumeInput,
		Output:        "",
		EventId:       &eventId,
		States:        states,
		QueryLanguage: definition.QueryLanguage,
		VariableScope: NewVariableScope(nil),
		MapItemIndex:  -1,
		IsRedrive:     true,
	}

	err = e.executeStates(ctx, execCtx)
	e.finalizeExecution(ctx, execution, execCtx, err)
	if err != nil && execution.Status != "ABORTED" {
		return err
	}
	return nil
}

// finalizeExecution updates the execution record and appends the appropriate
// terminal history event based on the execution outcome.
func (e *Executor) finalizeExecution(ctx context.Context, execution *sfnstore.Execution, execCtx *ExecutionContext, execErr error) {
	if execErr == nil {
		execution.Status = "SUCCEEDED"
		execution.StopDate = time.Now().UTC()
		execution.Output = execCtx.Output

		if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: execution.ExecutionArn,
			EventId:      execCtx.nextEventId(),
			Type:         "ExecutionSucceeded",
			Timestamp:    time.Now().UTC(),
			ExecutionSucceededEventDetails: &sfnstore.ExecutionSucceededEventDetails{
				Output: execCtx.Output,
			},
		}); err != nil {
			logs.Error("Failed to add ExecutionSucceeded event", logs.Err(err))
		}

		if err := e.store.UpdateExecution(ctx, execution); err != nil {
			logs.Error("Failed to update execution status to SUCCEEDED", logs.Err(err))
		}
		return
	}

	if execution.Status == "ABORTED" {
		if err := e.store.UpdateExecution(ctx, execution); err != nil {
			logs.Error("Failed to update execution status", logs.String("status", execution.Status), logs.Err(err))
		}
		return
	}

	if ctx.Err() == context.DeadlineExceeded {
		execution.Status = "TIMED_OUT"
		execution.Error = "ExecutionTimedOut"
		execution.Cause = "Execution timed out"
		execution.StopDate = time.Now().UTC()

		if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: execution.ExecutionArn,
			EventId:      execCtx.nextEventId(),
			Type:         "ExecutionTimedOut",
			Timestamp:    time.Now().UTC(),
			ExecutionTimedOutEventDetails: &sfnstore.ExecutionTimedOutEventDetails{
				Error: execution.Error,
				Cause: execution.Cause,
			},
		}); err != nil {
			logs.Error("Failed to add ExecutionTimedOut event", logs.Err(err))
		}
	} else if ctx.Err() == context.Canceled {
		execution.Status = "ABORTED"
		execution.StopDate = time.Now().UTC()
		// A StopExecution that landed first persisted the caller's error
		// and cause on the stored record; the terminal event reports that
		// pair and the record keeps it. A cancellation without a stop (the
		// shutdown sweep) still reports the generic pair.
		recordAlreadyTerminal := false
		if stored, gerr := e.store.GetExecution(context.Background(), execution.ExecutionArn); gerr == nil && isTerminalStatus(stored.Status) {
			recordAlreadyTerminal = true
			execution.Error = stored.Error
			execution.Cause = stored.Cause
		}
		if execution.Error == "" {
			execution.Error = "ExecutionAborted"
		}
		if execution.Cause == "" {
			execution.Cause = "Execution was aborted by StopExecution"
		}

		// The state that was running when the cancellation landed records
		// its modelled Aborted variant before the execution's terminal
		// event closes the history.
		e.logStateAbortedEvent(ctx, execCtx)

		if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: execution.ExecutionArn,
			EventId:      execCtx.nextEventId(),
			Type:         "ExecutionAborted",
			Timestamp:    time.Now().UTC(),
			ExecutionAbortedEventDetails: &sfnstore.ExecutionAbortedEventDetails{
				Error: execution.Error,
				Cause: execution.Cause,
			},
		}); err != nil {
			logs.Error("Failed to add ExecutionAborted event", logs.Err(err))
		}
		if recordAlreadyTerminal {
			return
		}
	} else {
		execution.Status = "FAILED"
		// Surface the state's error name — States.Runtime, task error
		// codes, and so on — rather than a generic marker: the
		// DescribeExecution error field carries the error string of the
		// failure.
		var stateErr *ExecutionError
		if errors.As(execErr, &stateErr) {
			execution.Error = stateErr.ErrorCode
			execution.Cause = stateErr.Cause
		} else {
			execution.Error = "ExecutionFailed"
			execution.Cause = "An internal error occurred during execution"
		}
		logs.Error("State machine execution failed", logs.String("arn", execution.ExecutionArn), logs.Err(execErr))
		execution.StopDate = time.Now().UTC()

		if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: execution.ExecutionArn,
			EventId:      execCtx.nextEventId(),
			Type:         "ExecutionFailed",
			Timestamp:    time.Now().UTC(),
			ExecutionFailedEventDetails: &sfnstore.ExecutionFailedEventDetails{
				Error: execution.Error,
				Cause: execution.Cause,
			},
		}); err != nil {
			logs.Error("Failed to add ExecutionFailed event", logs.Err(err))
		}
	}

	if err := e.store.UpdateExecution(ctx, execution); err != nil {
		logs.Error("Failed to update execution status", logs.String("status", execution.Status), logs.Err(err))
	}
}

// ExecutionContext holds the context for a state machine execution.
type ExecutionContext struct {
	Execution        *sfnstore.Execution
	Definition       *sfnstore.StateMachineDefinition
	CurrentState     string
	Input            string
	Output           string
	EventId          *int64
	States           map[string]sfnstore.State
	QueryLanguage    string
	VariableScope    *VariableScope
	PendingAssign    map[string]interface{}
	StateEnteredTime time.Time
	RetryCount       int32
	MapItemIndex     int
	MapItemValue     interface{}
	// MapItemKey carries the current iteration's key when a Map iterates an
	// object dataset; MapItemSource carries the Item provenance the context
	// object reports as Map.Item.Source ("STATE_DATA" for state input, the
	// Amazon S3 URI for ItemReader sources).
	MapItemKey        string
	MapItemSource     string
	AfterArguments    *string
	AfterItemSelector *string
	IsRedrive         bool
	// TaskToken carries the token minted for the activity-task attempt
	// being dispatched. Like RetryCount it is attempt-scoped state: it is
	// set before each attempt's input evaluation and cleared when the task
	// state finishes, so the context object only exposes a Task section
	// while a token actually backs it. Map and Parallel branches each own
	// their ExecutionContext, so concurrent branches never share a token.
	TaskToken string
	// MapItemReaderData substitutes the ItemReader's raw source bytes for
	// this execution. TestState populates it from the stateConfiguration
	// mapItemReaderData member; normal executions leave it empty so the
	// reader fetches from S3.
	MapItemReaderData string
	// The After* members feed the TestState inspectionData shape: they
	// record the intermediate data-processing results of the state under
	// test (afterInputPath, afterParameters, afterResultSelector,
	// afterResultPath, afterItemsPath, afterItemBatcher) alongside the Map
	// runtime settings TestState reports (maxConcurrency, tolerated
	// failure).
	AfterInputPath           *string
	AfterParameters          *string
	AfterResultSelector      *string
	AfterResultPath          *string
	AfterItemsPath           *string
	AfterItemBatcher         *string
	MaxConcurrencyValue      *int32
	ToleratedFailureCountVal *int64
	ToleratedFailurePctVal   *float64
	// SuppliedContext, when set, replaces the derived context object for
	// the run. TestState injects the caller's context parameter through
	// it so states under test resolve $$.Context against the supplied
	// object instead of the synthetic one.
	SuppliedContext map[string]interface{}
}

func (ctx *ExecutionContext) nextEventId() int64 {
	return atomic.AddInt64(ctx.EventId, 1)
}

// GetEffectiveQueryLanguage returns the query language for a state, checking
// the state-level override first, then the machine-level default, and
// falling back to "JSONPath".
func GetEffectiveQueryLanguage(state sfnstore.State, defaultLang string) string {
	if ql := state.GetQueryLanguage(); ql != "" {
		return ql
	}
	if defaultLang != "" {
		return defaultLang
	}
	return "JSONPath"
}

// IsJSONataState reports whether a state uses JSONata as its query language,
// checking the state-level override and machine-level default.
func IsJSONataState(state sfnstore.State, defaultLang string) bool {
	return GetEffectiveQueryLanguage(state, defaultLang) == "JSONata"
}

// ExecutionError represents an error that occurred during state machine execution.
type ExecutionError struct {
	ErrorCode string
	Cause     string
	// underlying carries a Go-level cause for the errors.Is/As chain: the
	// mid-execution cancellation vehicles wrap context.Canceled so every
	// container layer recognises an interrupted execution through the
	// platform's isCanceledError test instead of re-deriving it per layer.
	underlying error
}

// Error returns the error representation of the ExecutionError.
func (e *ExecutionError) Error() string {
	return e.ErrorCode + ": " + e.Cause
}

// Unwrap exposes the underlying Go-level cause, if any.
func (e *ExecutionError) Unwrap() error {
	return e.underlying
}

// executionInterruptedError builds the vehicle for an execution whose
// context ended mid-state. A cancellation (StopExecution or the shutdown
// sweep) wraps context.Canceled, so container aggregates classify it as an
// abort rather than a failure and Catch gates refuse to consume it; a
// machine-level deadline keeps the plain States.Timeout identity.
func executionInterruptedError(ctx context.Context, cause string) *ExecutionError {
	if ctx.Err() == context.Canceled {
		return &ExecutionError{ErrorCode: "States.Timeout", Cause: cause, underlying: context.Canceled}
	}
	return &ExecutionError{ErrorCode: "States.Timeout", Cause: cause}
}

// logStateAbortedEvent records the modelled Aborted variant of the state
// that was running when the execution was cancelled. Only the long-running
// state types carry an Aborted variant in the HistoryEventType enum; the
// instant states (Pass, Choice, Fail, Succeed) terminate before a
// cancellation can observe them.
func (e *Executor) logStateAbortedEvent(ctx context.Context, execCtx *ExecutionContext) {
	if execCtx == nil {
		return
	}
	state, exists := execCtx.States[execCtx.CurrentState]
	if !exists {
		return
	}
	var eventType string
	switch state.(type) {
	case *sfnstore.TaskState:
		eventType = "TaskStateAborted"
	case *sfnstore.WaitState:
		eventType = "WaitStateAborted"
	case *sfnstore.MapState:
		eventType = "MapStateAborted"
	case *sfnstore.ParallelState:
		eventType = "ParallelStateAborted"
	default:
		return
	}
	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            eventType,
		Timestamp:       time.Now().UTC(),
	})
}

// isCanceledError reports whether an error chain carries a context
// cancellation. Callers inside Map and Parallel hold a local variable named
// errors (the per-unit error slice), so this helper keeps the package
// qualifier reachable.
func isCanceledError(err error) bool {
	return errors.Is(err, context.Canceled)
}

func (e *Executor) executeStates(ctx context.Context, execCtx *ExecutionContext) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		state, exists := execCtx.States[execCtx.CurrentState]
		if !exists {
			return fmt.Errorf("state not found: %s", execCtx.CurrentState)
		}

		var nextState string
		var output string
		var err error
		var execErr *ExecutionError

		execCtx.PendingAssign = nil
		execCtx.StateEnteredTime = time.Now().UTC()
		execCtx.RetryCount = 0

		switch s := state.(type) {
		case *sfnstore.PassState:
			output, nextState, err = e.executePass(ctx, execCtx, s)
		case *sfnstore.TaskState:
			output, nextState, execErr = e.executeTask(ctx, execCtx, s)
			if execErr != nil {
				return execErr
			}
		case *sfnstore.ChoiceState:
			output, nextState, err = e.executeChoice(ctx, execCtx, s)
		case *sfnstore.WaitState:
			output, nextState, err = e.executeWait(ctx, execCtx, s)
		case *sfnstore.ParallelState:
			output, nextState, execErr = e.executeParallel(ctx, execCtx, s)
			if execErr != nil {
				return execErr
			}
		case *sfnstore.MapState:
			output, nextState, execErr = e.executeMap(ctx, execCtx, s)
			if execErr != nil {
				return execErr
			}
		case *sfnstore.FailState:
			return e.executeFail(ctx, execCtx, s)
		case *sfnstore.SucceedState:
			output, _, err = e.executeSucceed(ctx, execCtx, s)
			if err != nil {
				return err
			}
			execCtx.Output = output
			return nil
		default:
			return fmt.Errorf("unknown state type: %s", state.GetType())
		}

		if err != nil {
			return fmt.Errorf("state execution failed: %w", err)
		}

		if len(execCtx.PendingAssign) > 0 && execCtx.VariableScope != nil {
			if assignErr := execCtx.VariableScope.SetAll(execCtx.PendingAssign); assignErr != nil {
				return fmt.Errorf("failed to apply Assign in state %s: %w", execCtx.CurrentState, assignErr)
			}
		}

		execCtx.Input = output
		execCtx.Output = output

		if state.GetEnd() {
			break
		}

		if nextState == "" {
			nextState = state.GetNext()
		}
		if nextState == "" {
			break
		}
		execCtx.CurrentState = nextState
	}

	return nil
}

func (e *Executor) addExecutionHistoryEvent(ctx context.Context, execution *sfnstore.Execution, event *sfnstore.ExecutionHistoryEvent) error {
	event.ExecutionArn = execution.ExecutionArn
	err := e.store.AddExecutionHistoryEvent(ctx, event)
	if err == nil {
		e.publishHistoryToCloudWatchLogs(execution, event)
	}
	return err
}

func (e *Executor) logHistoryEvent(ctx context.Context, execution *sfnstore.Execution, event *sfnstore.ExecutionHistoryEvent) {
	if e.noPersistHistory {
		return
	}
	if err := e.addExecutionHistoryEvent(ctx, execution, event); err != nil {
		logs.Error("Failed to add history event", logs.String("type", event.Type), logs.Err(err))
	}
}

func (e *Executor) buildContextObject(execCtx *ExecutionContext) map[string]interface{} {
	if execCtx == nil {
		return map[string]interface{}{}
	}
	// A supplied context object replaces the derived one entirely: the
	// TestState context parameter represents the exact Context object for
	// the state under test.
	if execCtx.SuppliedContext != nil {
		return execCtx.SuppliedContext
	}

	ctx := map[string]interface{}{}

	if execCtx.Execution != nil {
		var execInput interface{}
		if execCtx.Execution.Input != "" {
			if err := json.Unmarshal([]byte(execCtx.Execution.Input), &execInput); err != nil {
				execInput = nil
			}
		}
		// RedriveTime exists in the Context object only for a redriven
		// execution ("RedriveTime Context object is only available if
		// you've redriven an execution").
		execution := map[string]interface{}{
			"Id":           execCtx.Execution.ExecutionArn,
			"Name":         execCtx.Execution.Name,
			"RoleArn":      e.extractExecutionRoleArn(),
			"StartTime":    execCtx.Execution.StartDate.Format(time.RFC3339),
			"Input":        execInput,
			"RedriveCount": execCtx.Execution.RedriveCount,
		}
		if !execCtx.Execution.RedriveDate.IsZero() {
			execution["RedriveTime"] = execCtx.Execution.RedriveDate.Format(time.RFC3339)
		}
		ctx["Execution"] = execution
	}

	if e.currentStateMachine != nil {
		ctx["StateMachine"] = map[string]interface{}{
			"Id":   e.currentStateMachine.StateMachineArn,
			"Name": e.currentStateMachine.Name,
		}
	}

	ctx["State"] = map[string]interface{}{
		"Name":        execCtx.CurrentState,
		"EnteredTime": execCtx.StateEnteredTime.Format(time.RFC3339),
		"RetryCount":  execCtx.RetryCount,
	}

	// The Task section exists in the context object only while a task
	// attempt with a token is being processed; the AWS context object
	// leaves it unpopulated elsewhere. JSONata Arguments reference the
	// token as $states.context.Task.Token, JSONPath Parameters as
	// $$.Task.Token.
	if execCtx.TaskToken != "" {
		ctx["Task"] = map[string]interface{}{
			"Token": execCtx.TaskToken,
		}
	}

	if execCtx.MapItemIndex >= 0 {
		// Source reports the item's provenance: "STATE_DATA" for state
		// input, the Amazon S3 URI for ItemReader sources; Key exists only
		// for object datasets.
		source := execCtx.MapItemSource
		if source == "" {
			source = "STATE_DATA"
		}
		item := map[string]interface{}{
			"Index":  execCtx.MapItemIndex,
			"Value":  execCtx.MapItemValue,
			"Source": source,
		}
		if execCtx.MapItemKey != "" {
			item["Key"] = execCtx.MapItemKey
		}
		ctx["Map"] = map[string]interface{}{"Item": item}
	}

	return ctx
}

func (e *Executor) buildStatesVarWithContext(execCtx *ExecutionContext, input, result, errorOutput interface{}) map[string]interface{} {
	ctxObj := e.buildContextObject(execCtx)
	return BuildStatesVar(input, result, errorOutput, ctxObj)
}

func (e *Executor) newQueryEvalError(ctx context.Context, execCtx *ExecutionContext, location, cause string) *ExecutionError {
	if execCtx.Execution != nil {
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: execCtx.Execution.ExecutionArn,
			EventId:      execCtx.nextEventId(),
			Type:         "EvaluationFailed",
			Timestamp:    time.Now().UTC(),
			EvaluationFailedEventDetails: &sfnstore.EvaluationFailedEventDetails{
				State:    execCtx.CurrentState,
				Cause:    cause,
				Error:    "States.QueryEvaluationError",
				Location: location,
			},
		})
	}
	return &ExecutionError{ErrorCode: "States.QueryEvaluationError", Cause: cause}
}

// newJSONPathEvalError classifies a JSONPath input/output processing
// failure. Step Functions reports JSONPath evaluation failures as
// States.Runtime — an unprocessable runtime exception that is not retriable
// and cannot be caught by a States.ALL catcher. States.QueryEvaluationError
// is reserved for JSONata expression failures, so a JSONPath state must
// never surface it. Every JSONPath evaluator funnels its failures through
// this classifier so the error code cannot diverge per state type.
func newJSONPathEvalError(field string, err error) *ExecutionError {
	return &ExecutionError{ErrorCode: "States.Runtime", Cause: fmt.Sprintf("%s: %s", field, err.Error())}
}

// newRuntimeValueError reports a run-time value that violates a state's
// value contract — for example a JSONata wait expression evaluating to a
// non-timestamp. Like JSONPath evaluation failures this is an
// unprocessable runtime exception surfaced as States.Runtime; it is not
// a JSONata evaluation failure, so States.QueryEvaluationError does not
// apply.
func newRuntimeValueError(field, cause string) *ExecutionError {
	return &ExecutionError{ErrorCode: "States.Runtime", Cause: fmt.Sprintf("%s: %s", field, cause)}
}

// historyEventLogsAtLevel reports whether a history event is published to
// the logging destination at the configured level, per the CloudWatch Logs
// level table in the developer guide: the execution-terminal failures
// (ExecutionFailed, ExecutionAborted, ExecutionTimedOut) form the FATAL
// class; state- and task-level failures — event types whose names end in
// Failed, Aborted or TimedOut, plus FailStateEntered — form the ERROR
// class; everything else is logged at ALL only, and OFF logs nothing.
func historyEventLogsAtLevel(level, eventType string) bool {
	executionTerminal := eventType == "ExecutionFailed" || eventType == "ExecutionAborted" || eventType == "ExecutionTimedOut"
	switch level {
	case "ALL":
		return true
	case "OFF":
		return false
	case "FATAL":
		return executionTerminal
	case "ERROR":
		return executionTerminal ||
			eventType == "FailStateEntered" ||
			strings.HasSuffix(eventType, "Failed") ||
			strings.HasSuffix(eventType, "Aborted") ||
			strings.HasSuffix(eventType, "TimedOut")
	}
	return false
}

func (e *Executor) publishHistoryToCloudWatchLogs(execution *sfnstore.Execution, event *sfnstore.ExecutionHistoryEvent) {
	if e.bus == nil || e.currentStateMachine == nil {
		return
	}
	lc := e.currentStateMachine.LoggingConfiguration
	if lc == nil || len(lc.Destinations) == 0 {
		return
	}

	logGroupArn := ""
	for _, dest := range lc.Destinations {
		if dest.CloudWatchLogsLogGroup != nil && dest.CloudWatchLogsLogGroup.LogGroupArn != "" {
			logGroupArn = dest.CloudWatchLogsLogGroup.LogGroupArn
			break
		}
	}
	if logGroupArn == "" {
		return
	}

	if !historyEventLogsAtLevel(lc.Level, event.Type) {
		return
	}

	_, _, region, _, _ := arn.SplitARN(logGroupArn)
	logGroup := arn.ExtractLogGroupNameFromARN(logGroupArn)
	if logGroup == "" {
		return
	}

	_, _, _, _, execResource := arn.SplitARN(execution.ExecutionArn)
	execName := execution.Name
	if segs := strings.Split(execResource, ":"); len(segs) > 0 && segs[len(segs)-1] != "" {
		execName = segs[len(segs)-1]
	}
	logStream := fmt.Sprintf("%s-%s", e.currentStateMachine.Name, execName)

	// The published record is the HistoryEvent API shape, so the
	// includeExecutionData configuration governs its payload members.
	record := historyEventToResponse(event, lc.IncludeExecutionData)
	eventBytes, err := json.Marshal(record)
	if err != nil {
		return
	}

	evt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup:  logGroup,
		LogStream: logStream,
		LogEvents: []eventbus.LogEntry{
			{Timestamp: event.Timestamp.UnixMilli(), Message: string(eventBytes)},
		},
	}
	evt.Region = region
	evt.AccountID = e.accountID

	if err := e.bus.Publish(context.Background(), evt); err != nil {
		logs.Debug("Failed to publish execution history to CloudWatch Logs",
			logs.String("executionArn", execution.ExecutionArn),
			logs.Err(err))
	}
}

// updateExecutionWithRetry persists the execution record, retrying on
// transient store errors. The first attempt is synchronous; the retry
// uses a fresh context.Background() so the call still completes even when
// the caller's context is already cancelled (e.g. parsing failed late
// during shutdown).
func (e *Executor) updateExecutionWithRetry(ctx context.Context, execution *sfnstore.Execution) error {
	err := e.store.UpdateExecution(ctx, execution)
	if err == nil {
		return nil
	}
	retryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if retryErr := e.store.UpdateExecution(retryCtx, execution); retryErr != nil {
		return fmt.Errorf("update failed (initial: %v, retry: %v)", err, retryErr)
	}
	return nil
}
