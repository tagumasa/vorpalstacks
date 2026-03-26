package sfn

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeParallel(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ParallelState) (string, string, *ExecutionError) {
	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ParallelStateEntered",
		Timestamp:       time.Now().UTC(),
		ParallelStateEnteredEventDetails: &sfnstore.ParallelStateEnteredEventDetails{
			Input:    processedInput,
			Name:     execCtx.CurrentState,
			Branches: getBranchNames(state.Branches),
		},
	})

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]string, len(state.Branches))
	errors := make([]error, len(state.Branches))

	for i, branch := range state.Branches {
		wg.Add(1)
		go func(idx int, b *sfnstore.StateMachineDefinition) {
			defer wg.Done()
			branchStates, err := e.extractStatesFromDefinition(b)
			if err != nil {
				mu.Lock()
				errors[idx] = err
				mu.Unlock()
				return
			}
			branchCtx := &ExecutionContext{
				Execution:    execCtx.Execution,
				Definition:   b,
				CurrentState: b.StartAt,
				Input:        processedInput,
				Output:       "",
				EventId:      execCtx.EventId,
				States:       branchStates,
			}
			execErr := e.executeStates(ctx, branchCtx)
			mu.Lock()
			defer mu.Unlock()
			errors[idx] = execErr
			if execErr == nil {
				results[idx] = branchCtx.Output
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
		if len(state.Catch) > 0 {
			catchPolicy := e.findMatchingCatchPolicy(state.Catch, "States.BranchFailed")
			if catchPolicy != nil {
				catchOutput := e.buildCatchOutput(processedInput, "States.BranchFailed", firstError.Error(), catchPolicy.ResultPath)
				return catchOutput, catchPolicy.Next, nil
			}
		}
		return "", "", &ExecutionError{Error: "States.BranchFailed", Cause: firstError.Error()}
	}

	output := fmt.Sprintf(`[%s]`, strings.Join(results, ","))
	output = e.applyOutputPath(output, state.GetOutputPath())

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ParallelStateExited",
		Timestamp:       time.Now().UTC(),
		ParallelStateExitedEventDetails: &sfnstore.ParallelStateExitedEventDetails{
			Output: output,
			Name:   execCtx.CurrentState,
		},
	})

	return output, state.Next, nil
}

func getBranchNames(branches []*sfnstore.StateMachineDefinition) []string {
	names := make([]string, len(branches))
	for i, branch := range branches {
		names[i] = branch.StartAt
	}
	return names
}
