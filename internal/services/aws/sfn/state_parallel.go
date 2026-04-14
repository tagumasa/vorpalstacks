package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeParallel(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ParallelState) (string, string, *ExecutionError) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

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
				Execution:     execCtx.Execution,
				Definition:    b,
				CurrentState:  b.StartAt,
				Input:         processedInput,
				Output:        "",
				EventId:       execCtx.EventId,
				States:        branchStates,
				QueryLanguage: execCtx.QueryLanguage,
				VariableScope: execCtx.VariableScope.NewChild(),
				MapItemIndex:  -1,
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
				if isJSONata {
					return e.executeParallelJSONataCatch(ctx, execCtx, state, processedInput, "States.BranchFailed", firstError.Error(), catchPolicy)
				}
				catchOutput := e.buildCatchOutput(processedInput, "States.BranchFailed", firstError.Error(), catchPolicy.ResultPath)
				return catchOutput, catchPolicy.Next, nil
			}
		}
		return "", "", &ExecutionError{ErrorCode: "States.BranchFailed", Cause: firstError.Error()}
	}

	output := fmt.Sprintf(`[%s]`, strings.Join(results, ","))

	if isJSONata {
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
		var resultData interface{}
		if err := json.Unmarshal([]byte(output), &resultData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidOutput", Cause: "failed to parse output JSON"}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, resultData, nil)

		if len(state.Assign) > 0 {
			evaluated, err := evaluateAssign(ctx, state.Assign, statesVar, execCtx.VariableScope)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
			}
			execCtx.PendingAssign = evaluated
		}

		if state.JSONataOutput == nil && len(state.OutputRaw) > 0 {
			var err error
			state.JSONataOutput, err = resolveJSONataOutput(state)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
			}
		}

		if state.JSONataOutput != nil {
			resolved, err := e.applyJSONataOutput(ctx, state.JSONataOutput, statesVar, execCtx.VariableScope)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
			}
			outputJSON, err := json.Marshal(resolved)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", fmt.Sprintf("failed to marshal: %s", err.Error()))
			}
			output = string(outputJSON)
		}
	} else {
		output = e.applyOutputPath(output, state.GetOutputPath())
	}

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

func (e *Executor) executeParallelJSONataCatch(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ParallelState, processedInput, errorCode, cause string, catchPolicy *sfnstore.CatchPolicy) (string, string, *ExecutionError) {
	errorOutput := map[string]interface{}{
		"Error": errorCode,
		"Cause": cause,
	}

	var inputData interface{}
	if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
		return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
	}
	statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, errorOutput)

	if len(catchPolicy.Assign) > 0 {
		evaluated, err := evaluateAssign(ctx, catchPolicy.Assign, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Catch.Assign", err.Error())
		}
		execCtx.PendingAssign = evaluated
	}

	if catchPolicy.Output != nil {
		resolved, err := e.applyJSONataOutput(ctx, catchPolicy.Output, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Catch.Output", err.Error())
		}
		outputJSON, err := json.Marshal(resolved)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Catch.Output", fmt.Sprintf("failed to marshal: %s", err.Error()))
		}
		return string(outputJSON), catchPolicy.Next, nil
	}

	errorJSON, _ := json.Marshal(errorOutput)
	return string(errorJSON), catchPolicy.Next, nil
}

func getBranchNames(branches []*sfnstore.StateMachineDefinition) []string {
	names := make([]string, len(branches))
	for i, branch := range branches {
		names[i] = branch.StartAt
	}
	return names
}
