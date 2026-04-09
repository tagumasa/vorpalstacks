package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/server/eventbus"
	"vorpalstacks/internal/services/aws/common"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	"vorpalstacks/internal/utils/aws/arn"
)

// Executor manages the execution of Step Functions state machines.
type Executor struct {
	store               *sfnstore.StepFunctionStore
	lambdaInvoker       common.LambdaInvoker
	sqsStore            sqsstore.SQSStoreInterface
	snsStore            snsstore.SNSStoreInterface
	eventsStore         *eventsstore.EventsStore
	accountID           string
	region              string
	bus                 *eventbus.EventBus
	currentRoleArn      string
	currentExecution    *sfnstore.Execution
	currentStateMachine *sfnstore.StateMachine
}

// NewExecutor creates a new Step Functions executor with the given store and Lambda invoker.
func NewExecutor(store *sfnstore.StepFunctionStore, lambdaInvoker common.LambdaInvoker) *Executor {
	return &Executor{
		store:         store,
		lambdaInvoker: lambdaInvoker,
	}
}

// NewExecutorWithStores creates a new Step Functions executor with all store dependencies.
func NewExecutorWithStores(store *sfnstore.StepFunctionStore, lambdaInvoker common.LambdaInvoker, sqsStore sqsstore.SQSStoreInterface, snsStore snsstore.SNSStoreInterface, eventsStore *eventsstore.EventsStore, accountID, region string) *Executor {
	return &Executor{
		store:         store,
		lambdaInvoker: lambdaInvoker,
		sqsStore:      sqsStore,
		snsStore:      snsStore,
		eventsStore:   eventsStore,
		accountID:     accountID,
		region:        region,
	}
}

// SetEventBus injects the event bus. When set, SNS publish and EventBridge
// PutEvents task integrations route through the bus, fixing fan-out and rule
// matching bugs that exist in the direct store access path.
func (e *Executor) SetEventBus(bus *eventbus.EventBus) {
	e.bus = bus
}

// ExecuteStateMachine executes a state machine with the given execution context.
func (e *Executor) ExecuteStateMachine(ctx context.Context, execution *sfnstore.Execution) error {
	e.currentExecution = execution
	sm, err := e.store.GetStateMachine(ctx, execution.StateMachineArn)
	if err == nil && sm != nil {
		e.currentRoleArn = sm.RoleArn
		e.currentStateMachine = sm
	}

	definition, states, err := e.parseDefinition(ctx, execution.StateMachineArn)
	if err != nil {
		execution.Status = "FAILED"
		execution.Error = "InvalidDefinition"
		logs.Error("Failed to parse state machine definition", logs.String("arn", execution.StateMachineArn), logs.Err(err))
		execution.Cause = "Invalid state machine definition syntax"
		execution.StopDate = time.Now().UTC()
		if updateErr := e.store.UpdateExecution(ctx, execution); updateErr != nil {
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
			Input:           execution.Input,
			RoleArn:         roleArn,
			StateMachineArn: execution.StateMachineArn,
			Name:            execution.Name,
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
	if err != nil {
		if execution.Status == "ABORTED" {
			return nil
		}

		if ctx.Err() == context.DeadlineExceeded {
			execution.Status = "TIMED_OUT"
			execution.Error = "ExecutionTimedOut"
			execution.Cause = "Execution timed out"
			execution.StopDate = time.Now().UTC()

			if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn: execution.ExecutionArn,
				EventId:      *execCtx.EventId,
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
			execution.Error = "ExecutionAborted"
			execution.Cause = "Execution was aborted by StopExecution"
			execution.StopDate = time.Now().UTC()

			if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn: execution.ExecutionArn,
				EventId:      *execCtx.EventId,
				Type:         "ExecutionAborted",
				Timestamp:    time.Now().UTC(),
				ExecutionAbortedEventDetails: &sfnstore.ExecutionAbortedEventDetails{
					Error: execution.Error,
					Cause: execution.Cause,
				},
			}); err != nil {
				logs.Error("Failed to add ExecutionAborted event", logs.Err(err))
			}
		} else {
			execution.Status = "FAILED"
			execution.Error = "ExecutionFailed"
			logs.Error("State machine execution failed", logs.String("arn", execution.ExecutionArn), logs.Err(err))
			execution.Cause = "An internal error occurred during execution"
			execution.StopDate = time.Now().UTC()

			if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn: execution.ExecutionArn,
				EventId:      *execCtx.EventId,
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
		return err
	}

	execution.Status = "SUCCEEDED"
	execution.StopDate = time.Now().UTC()
	execution.Output = execCtx.Output

	if err := e.addExecutionHistoryEvent(ctx, execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: execution.ExecutionArn,
		EventId:      *execCtx.EventId,
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
	return nil
}

// ExecutionContext holds the context for a state machine execution.
type ExecutionContext struct {
	Execution         *sfnstore.Execution
	Definition        *sfnstore.StateMachineDefinition
	CurrentState      string
	Input             string
	Output            string
	EventId           *int64
	States            map[string]sfnstore.State
	QueryLanguage     string
	VariableScope     *VariableScope
	PendingAssign     map[string]interface{}
	StateEnteredTime  time.Time
	RetryCount        int32
	MapItemIndex      int
	MapItemValue      interface{}
	AfterArguments    *string
	AfterItemSelector *string
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
}

// Error returns the error representation of the ExecutionError.
func (e *ExecutionError) Error() string {
	return e.ErrorCode + ": " + e.Cause
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
				return fmt.Errorf("%s: %s", execErr.ErrorCode, execErr.Cause)
			}
		case *sfnstore.ChoiceState:
			nextState, err = e.executeChoice(ctx, execCtx, s)
			if err == nil {
				output = execCtx.Input
			}
		case *sfnstore.WaitState:
			output, nextState, err = e.executeWait(ctx, execCtx, s)
		case *sfnstore.ParallelState:
			output, nextState, execErr = e.executeParallel(ctx, execCtx, s)
			if execErr != nil {
				return fmt.Errorf("%s: %s", execErr.ErrorCode, execErr.Cause)
			}
		case *sfnstore.MapState:
			output, nextState, execErr = e.executeMap(ctx, execCtx, s)
			if execErr != nil {
				return fmt.Errorf("%s: %s", execErr.ErrorCode, execErr.Cause)
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
	if err := e.addExecutionHistoryEvent(ctx, execution, event); err != nil {
		logs.Error("Failed to add history event", logs.String("type", event.Type), logs.Err(err))
	}
}

func (e *Executor) buildContextObject(execCtx *ExecutionContext) map[string]interface{} {
	ctx := map[string]interface{}{}

	if execCtx.Execution != nil {
		var execInput interface{}
		if execCtx.Execution.Input != "" {
			json.Unmarshal([]byte(execCtx.Execution.Input), &execInput)
		}
		ctx["Execution"] = map[string]interface{}{
			"Id":        execCtx.Execution.ExecutionArn,
			"Name":      execCtx.Execution.Name,
			"RoleArn":   e.extractExecutionRoleArn(),
			"StartTime": execCtx.Execution.StartDate.Format(time.RFC3339),
			"Input":     execInput,
		}
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

	if execCtx.MapItemIndex >= 0 {
		ctx["Map"] = map[string]interface{}{
			"Item": map[string]interface{}{
				"Index": execCtx.MapItemIndex,
				"Value": execCtx.MapItemValue,
			},
		}
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

	_, _, region, _, _ := arn.SplitARN(logGroupArn)
	logGroup := arn.ExtractLogGroupNameFromARN(logGroupArn)
	if logGroup == "" {
		return
	}

	execParts := strings.Split(execution.ExecutionArn, ":")
	execName := execution.Name
	if len(execParts) > 0 {
		if last := execParts[len(execParts)-1]; last != "" {
			execName = last
		}
	}
	logStream := fmt.Sprintf("%s-%s", e.currentStateMachine.Name, execName)

	eventBytes, err := json.Marshal(event)
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
