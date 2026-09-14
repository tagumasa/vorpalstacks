package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"

	"vorpalstacks/internal/core/logs"
)

func (e *Executor) executeParallel(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ParallelState) (string, string, *ExecutionError) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	processedInput, ipErr := e.applyInputPath(execCtx, execCtx.Input, state.GetInputPath())
	if ipErr != nil {
		return "", "", ipErr
	}

	if !isJSONata && state.Parameters != nil {
		// A Parallel Parameters payload template transforms the state
		// input once; every branch receives the transformed value.
		applied, evalErr := e.applyParameters(execCtx, "", processedInput, state.Parameters)
		if evalErr != nil {
			return "", "", evalErr
		}
		processedInput = applied
	}

	if isJSONata && state.Arguments != nil {
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		argsInput, err := e.applyJSONataArguments(ctx, state.Arguments, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Arguments", err.Error())
		}
		processedInput = argsInput
		execCtx.AfterArguments = &processedInput
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:             execCtx.Execution.ExecutionArn,
		EventId:                  eventId,
		PreviousEventId:          eventId - 1,
		Type:                     "ParallelStateEntered",
		Timestamp:                time.Now().UTC(),
		StateEnteredEventDetails: stateEnteredDetails(execCtx, processedInput),
	})

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ParallelStateStarted",
		Timestamp:       time.Now().UTC(),
	})

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]string, len(state.Branches))
	errors := make([]error, len(state.Branches))

	// Branch definitions are fixed at creation time, so the state tables
	// are extracted once here rather than re-parsed inside each branch
	// worker; a branch whose parse fails records the error and lets the
	// shared aggregation own the outcome, exactly as the per-worker parse
	// did.
	branchStates := make([]map[string]sfnstore.State, len(state.Branches))
	branchParseErrs := make([]error, len(state.Branches))
	for i, branch := range state.Branches {
		states, err := extractStatesFromDefinition(branch)
		branchStates[i] = states
		branchParseErrs[i] = err
	}

	if execCtx.IsRedrive && execCtx.Execution.ParallelCheckpoints != nil {
		if cp, ok := execCtx.Execution.ParallelCheckpoints[execCtx.CurrentState]; ok {
			for idx, cached := range cp.BranchResults {
				if idx >= 0 && idx < len(results) {
					results[idx] = cached
				}
			}
		}
	}

	for i, branch := range state.Branches {
		if results[i] != "" {
			continue
		}
		if branchParseErrs[i] != nil {
			errors[i] = branchParseErrs[i]
			continue
		}
		wg.Add(1)
		go func(idx int, b *sfnstore.StateMachineDefinition) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logs.Error("sfn: panic in parallel branch worker", logs.Int("index", idx), logs.Any("panic", r), logs.String("stack", string(debug.Stack())))
					mu.Lock()
					errors[idx] = fmt.Errorf("internal panic: %v", r)
					mu.Unlock()
				}
			}()
			branchCtx := &ExecutionContext{
				Execution:     execCtx.Execution,
				Definition:    b,
				CurrentState:  b.StartAt,
				Input:         processedInput,
				Output:        "",
				EventId:       execCtx.EventId,
				States:        branchStates[idx],
				QueryLanguage: execCtx.QueryLanguage,
				VariableScope: execCtx.VariableScope.NewChild(),
				MapItemIndex:  -1,
			}
			// The branch's variables are scoped to the branch: the bytes
			// return to the execution budget once it completes.
			defer branchCtx.VariableScope.Release()
			execErr := e.executeStates(ctx, branchCtx)
			mu.Lock()
			defer mu.Unlock()
			errors[idx] = execErr
			if execErr == nil {
				results[idx] = branchCtx.Output
				// Branch results are checkpointed on every run, not only on
				// redrives: the original failed run must leave its successful
				// branches behind so the FIRST redrive re-runs only the
				// failed ones ("Reschedules and redrives only those branches
				// that failed or aborted"). The checkpoint reaches the store
				// with the terminal record write.
				if execCtx.Execution.ParallelCheckpoints == nil {
					execCtx.Execution.ParallelCheckpoints = make(map[string]*sfnstore.ParallelCheckpoint)
				}
				cp, ok := execCtx.Execution.ParallelCheckpoints[execCtx.CurrentState]
				if !ok {
					cp = &sfnstore.ParallelCheckpoint{BranchResults: make(map[int]string)}
					execCtx.Execution.ParallelCheckpoints[execCtx.CurrentState] = cp
				}
				cp.BranchResults[idx] = branchCtx.Output
			}
		}(i, branch)
	}

	wg.Wait()

	var firstError error
	for _, err := range errors {
		if err != nil {
			firstError = err
			break
		}
	}

	if firstError != nil {
		// The state-level terminal event carries no detail members in the
		// model; an aborted parallel records the Aborted variant instead of
		// Failed.
		eventId = execCtx.nextEventId()
		terminalType := "ParallelStateFailed"
		if isCanceledError(firstError) {
			terminalType = "ParallelStateAborted"
		}
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    execCtx.Execution.ExecutionArn,
			EventId:         eventId,
			PreviousEventId: eventId - 1,
			Type:            terminalType,
			Timestamp:       time.Now().UTC(),
		})

		// Retriers run before catchers ("Step Functions uses any appropriate
		// retriers first. If the retry policy fails to resolve the error,
		// Step Functions applies the matching catcher transition"); a
		// cancellation is not a failure to retry — the abort path owns the
		// history. Each attempt re-enters the state, re-emitting its
		// Entered/Started pair.
		if !isCanceledError(firstError) && len(state.Retry) > 0 {
			if matchedRetry := e.findMatchingRetryPolicy(state.Retry, "States.BranchFailed"); matchedRetry != nil && execCtx.RetryCount < matchedRetry.MaxAttempts {
				if e.sleepForRetry(ctx, matchedRetry, execCtx.RetryCount+1) {
					return "", "", executionInterruptedError(ctx, "Execution interrupted during retry")
				}
				execCtx.RetryCount++
				return e.executeParallel(ctx, execCtx, state)
			}
		}

		// A catcher never consumes a cancellation either: the machine is
		// stopping, no recovery transition may run on the dead context.
		if !isCanceledError(firstError) && len(state.Catch) > 0 {
			catchPolicy := e.findMatchingCatchPolicy(state.Catch, "States.BranchFailed")
			if catchPolicy != nil {
				if isJSONata {
					return e.executeJSONataCatch(ctx, execCtx, processedInput, "States.BranchFailed", firstError.Error(), catchPolicy)
				}
				catchOutput, coErr := e.buildCatchOutput(processedInput, "States.BranchFailed", firstError.Error(), catchPolicy.ResultPath)
				if coErr != nil {
					return "", "", coErr
				}
				if len(catchPolicy.Assign) > 0 {
					if err := e.applyJSONPathCatchAssign(execCtx, catchPolicy.Assign, catchOutput); err != nil {
						return "", "", err
					}
				}
				return catchOutput, catchPolicy.Next, nil
			}
		}
		if isCanceledError(firstError) {
			// The interruption itself propagates as a cancellation: the
			// abort path owns the history and the execution's terminal
			// classification.
			return "", "", executionInterruptedError(ctx, firstError.Error())
		}
		return "", "", &ExecutionError{ErrorCode: "States.BranchFailed", Cause: firstError.Error()}
	}

	output := fmt.Sprintf(`[%s]`, strings.Join(results, ","))

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ParallelStateSucceeded",
		Timestamp:       time.Now().UTC(),
	})

	output, outErr := e.applyContainerStateOutput(ctx, execCtx, state, processedInput, output)
	if outErr != nil {
		return "", "", outErr
	}

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:            execCtx.Execution.ExecutionArn,
		EventId:                 eventId,
		PreviousEventId:         eventId - 1,
		Type:                    "ParallelStateExited",
		Timestamp:               time.Now().UTC(),
		StateExitedEventDetails: stateExitedDetails(execCtx, output),
	})

	return output, state.Next, nil
}
