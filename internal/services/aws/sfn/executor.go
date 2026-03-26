package sfn

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/services/aws/common"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
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
		log.Printf("Failed to parse state machine definition for %s: %v", execution.StateMachineArn, err)
		execution.Cause = "Invalid state machine definition syntax"
		execution.StopDate = time.Now().UTC()
		if err := e.store.UpdateExecution(ctx, execution); err != nil {
			log.Printf("Failed to update execution status to FAILED after definition error: %v", err)
		}
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
		log.Printf("Failed to add ExecutionStarted event: %v", err)
	}

	execCtx := &ExecutionContext{
		Execution:    execution,
		Definition:   definition,
		CurrentState: definition.StartAt,
		Input:        execution.Input,
		Output:       "",
		EventId:      &eventId,
		States:       states,
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
				log.Printf("Failed to add ExecutionTimedOut event: %v", err)
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
				log.Printf("Failed to add ExecutionAborted event: %v", err)
			}
		} else {
			execution.Status = "FAILED"
			execution.Error = "ExecutionFailed"
			log.Printf("State machine execution failed for %s: %v", execution.ExecutionArn, err)
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
				log.Printf("Failed to add ExecutionFailed event: %v", err)
			}
		}

		if err := e.store.UpdateExecution(ctx, execution); err != nil {
			log.Printf("Failed to update execution status to %s: %v", execution.Status, err)
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
		log.Printf("Failed to add ExecutionSucceeded event: %v", err)
	}

	if err := e.store.UpdateExecution(ctx, execution); err != nil {
		log.Printf("Failed to update execution status to SUCCEEDED: %v", err)
	}
	return nil
}

// ExecutionContext holds the context for a state machine execution.
type ExecutionContext struct {
	Execution    *sfnstore.Execution
	Definition   *sfnstore.StateMachineDefinition
	CurrentState string
	Input        string
	Output       string
	EventId      *int64
	States       map[string]sfnstore.State
}

func (ctx *ExecutionContext) nextEventId() int64 {
	return atomic.AddInt64(ctx.EventId, 1)
}

// ExecutionError represents an error that occurred during state machine execution.
type ExecutionError struct {
	Error string
	Cause string
}

// String returns the string representation of the ExecutionError.
func (e *ExecutionError) String() string {
	return e.Error + ": " + e.Cause
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

		switch s := state.(type) {
		case *sfnstore.PassState:
			output, nextState, err = e.executePass(ctx, execCtx, s)
		case *sfnstore.TaskState:
			output, nextState, execErr = e.executeTask(ctx, execCtx, s)
			if execErr != nil {
				return fmt.Errorf("%s: %s", execErr.Error, execErr.Cause)
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
				return fmt.Errorf("%s: %s", execErr.Error, execErr.Cause)
			}
		case *sfnstore.MapState:
			output, nextState, execErr = e.executeMap(ctx, execCtx, s)
			if execErr != nil {
				return fmt.Errorf("%s: %s", execErr.Error, execErr.Cause)
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
	return e.store.AddExecutionHistoryEvent(ctx, event)
}

func (e *Executor) logHistoryEvent(ctx context.Context, execution *sfnstore.Execution, event *sfnstore.ExecutionHistoryEvent) {
	if err := e.addExecutionHistoryEvent(ctx, execution, event); err != nil {
		log.Printf("Failed to add %s history event: %v", event.Type, err)
	}
}
