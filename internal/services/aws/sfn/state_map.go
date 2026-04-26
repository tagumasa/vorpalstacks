package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeMap(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState) (string, string, *ExecutionError) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "MapStateEntered",
		Timestamp:       time.Now().UTC(),
		MapStateEnteredEventDetails: &sfnstore.MapStateEnteredEventDetails{
			Input: processedInput,
			Name:  execCtx.CurrentState,
		},
	})

	var itemsArray []interface{}

	if isJSONata {
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)

		if state.Items != nil {
			vars := buildVarsMap(statesVar, execCtx.VariableScope)
			resolved, err := ResolveTemplate(ctx, state.Items, nil, vars)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Items", err.Error())
			}
			switch v := resolved.(type) {
			case []interface{}:
				itemsArray = v
			default:
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "Items must evaluate to an array"}
			}
		} else {
			var inputDataMap map[string]interface{}
			if err := json.Unmarshal([]byte(processedInput), &inputDataMap); err != nil {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
			}
			itemsPath := state.ItemsPath
			if itemsPath == "" {
				itemsPath = "$"
			}
			items, err := getJSONPathValue(inputDataMap, itemsPath)
			if err != nil {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: err.Error()}
			}
			var ok bool
			itemsArray, ok = items.([]interface{})
			if !ok {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "items is not an array"}
			}
		}
	} else {
		var inputData map[string]interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}

		itemsPath := state.ItemsPath
		if itemsPath == "" {
			itemsPath = "$"
		}

		items, err := getJSONPathValue(inputData, itemsPath)
		if err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: err.Error()}
		}

		var ok bool
		itemsArray, ok = items.([]interface{})
		if !ok {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "items is not an array"}
		}
	}

	maxConcurrency := int(state.MaxConcurrency)
	if maxConcurrency <= 0 {
		maxConcurrency = len(itemsArray)
	}

	mapRunArn := generateMapRunArn(e.store, e.region, e.accountID, execCtx.Execution.ExecutionArn, execCtx.CurrentState)
	now := time.Now().UTC()
	mapRunRecord := &sfnstore.MapRun{
		MapRunArn:           mapRunArn,
		ExecutionArn:        execCtx.Execution.ExecutionArn,
		StateMachineArn:     execCtx.Execution.StateMachineArn,
		Name:                execCtx.CurrentState,
		Status:              "RUNNING",
		StartDate:           now.Unix(),
		ItemCount:           int64(len(itemsArray)),
		MaxConcurrency:      int64(maxConcurrency),
		ItemsProcessedCount: 0,
		ItemsFailedCount:    0,
		ItemsCancelledCount: 0,
	}
	if err := e.store.CreateMapRun(ctx, mapRunRecord); err != nil {
		logs.Warn("failed to create map run record", logs.Err(err))
	}

	defer func() {
		mapRunRecord.StopDate = time.Now().UTC().Unix()
		if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
			logs.Error("sfn: failed to update map run status", logs.Err(err))
		}
	}()

	var wg sync.WaitGroup
	results := make([]string, len(itemsArray))
	errors := make([]error, len(itemsArray))
	itemsProcessed := int64(0)
	itemsFailed := int64(0)
	var mu sync.Mutex

	processedItems := make([]interface{}, len(itemsArray))
	if state.ItemSelector != nil {
		for i, item := range itemsArray {
			execCtx.MapItemIndex = i
			execCtx.MapItemValue = item
			var selected interface{}
			var err error
			if isJSONata {
				selected, err = e.applyItemSelector(ctx, execCtx, state.ItemSelector, item, true)
			} else {
				selected, err = e.applyItemSelectorJSONPath(state.ItemSelector, item)
			}
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "ItemSelector", err.Error())
			}
			processedItems[i] = selected
		}
		selJSON, _ := json.Marshal(processedItems)
		s := string(selJSON)
		execCtx.AfterItemSelector = &s
		execCtx.MapItemIndex = -1
		execCtx.MapItemValue = nil
	} else {
		copy(processedItems, itemsArray)
	}

	sem := make(chan struct{}, maxConcurrency)

	for i := range processedItems {
		wg.Add(1)
		go func(idx int, itemValue interface{}, originalItem interface{}) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			defer func() {
				if r := recover(); r != nil {
					logs.Error("sfn: panic in map worker", logs.Int("index", idx), logs.Any("panic", r), logs.String("stack", string(debug.Stack())))
					mu.Lock()
					errors[idx] = fmt.Errorf("internal panic: %v", r)
					itemsFailed++
					mu.Unlock()
				}
			}()

			itemJSON, err := json.Marshal(itemValue)
			if err != nil {
				mu.Lock()
				errors[idx] = fmt.Errorf("failed to marshal map item: %w", err)
				itemsFailed++
				mu.Unlock()
				return
			}
			iteratorStates, err := e.extractStatesFromDefinition(state.GetIterator())
			if err != nil {
				mu.Lock()
				defer mu.Unlock()
				errors[idx] = err
				itemsFailed++
				return
			}
			iteratorCtx := &ExecutionContext{
				Execution:     execCtx.Execution,
				Definition:    state.GetIterator(),
				CurrentState:  state.GetIterator().StartAt,
				Input:         string(itemJSON),
				Output:        "",
				EventId:       execCtx.EventId,
				States:        iteratorStates,
				QueryLanguage: execCtx.QueryLanguage,
				VariableScope: execCtx.VariableScope.NewChild(),
				MapItemIndex:  idx,
				MapItemValue:  originalItem,
			}
			execErr := e.executeStates(ctx, iteratorCtx)
			mu.Lock()
			defer mu.Unlock()
			errors[idx] = execErr
			if execErr == nil {
				results[idx] = iteratorCtx.Output
				itemsProcessed++
			} else {
				itemsFailed++
			}
		}(i, processedItems[i], itemsArray[i])
	}

	wg.Wait()

	mapRunRecord.ItemsProcessedCount = itemsProcessed
	mapRunRecord.ItemsFailedCount = itemsFailed
	if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
		logs.Warn("failed to update map run after completion", logs.Err(err))
	}

	var firstError error
	for _, err := range errors {
		if err != nil {
			firstError = err
			break
		}
	}

	if firstError != nil {
		mapRunRecord.Status = "FAILED"
		if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
			logs.Warn("failed to update map run to FAILED", logs.Err(err))
		}

		if len(state.Catch) > 0 {
			catchPolicy := e.findMatchingCatchPolicy(state.Catch, "States.IteratorFailed")
			if catchPolicy != nil {
				isJSONataCatch := IsJSONataState(state, execCtx.QueryLanguage)
				if isJSONataCatch {
					return e.executeMapJSONataCatch(ctx, execCtx, state, processedInput, "States.IteratorFailed", firstError.Error(), catchPolicy)
				}
				catchOutput := e.buildCatchOutput(processedInput, "States.IteratorFailed", firstError.Error(), catchPolicy.ResultPath)
				return catchOutput, catchPolicy.Next, nil
			}
		}
		return "", "", &ExecutionError{ErrorCode: "States.IteratorFailed", Cause: firstError.Error()}
	}

	mapRunRecord.Status = "SUCCEEDED"
	if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
		logs.Warn("Failed to update map run status", logs.Err(err))
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
		Type:            "MapStateExited",
		Timestamp:       time.Now().UTC(),
		MapStateExitedEventDetails: &sfnstore.MapStateExitedEventDetails{
			Output:         output,
			Name:           execCtx.CurrentState,
			ItemsProcessed: itemsProcessed,
			ItemsFailed:    itemsFailed,
		},
	})

	return output, state.Next, nil
}

func (e *Executor) executeMapJSONataCatch(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, processedInput, errorCode, cause string, catchPolicy *sfnstore.CatchPolicy) (string, string, *ExecutionError) {
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

	errorJSON, err := json.Marshal(errorOutput)
	if err != nil {
		errorJSON = []byte(`{"error":"failed to marshal error output"}`)
	}
	return string(errorJSON), catchPolicy.Next, nil
}
