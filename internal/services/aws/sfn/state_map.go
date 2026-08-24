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

	if state.ItemReader != nil {
		readerItems, rerr := e.readItemReaderItems(ctx, execCtx, state, execCtx.MapItemReaderData)
		if rerr != nil {
			// ItemReader failures surface the documented
			// States.ItemReaderFailed error and pass through the
			// state's Catch handlers like any other runtime error.
			if len(state.Catch) > 0 {
				if catchPolicy := e.findMatchingCatchPolicy(state.Catch, rerr.ErrorCode); catchPolicy != nil {
					if IsJSONataState(state, execCtx.QueryLanguage) {
						return e.executeMapJSONataCatch(ctx, execCtx, state, processedInput, rerr.ErrorCode, rerr.Cause, catchPolicy)
					}
					catchOutput := e.buildCatchOutput(processedInput, rerr.ErrorCode, rerr.Cause, catchPolicy.ResultPath)
					return catchOutput, catchPolicy.Next, nil
				}
			}
			return "", "", rerr
		}
		itemsArray = readerItems
	} else if isJSONata {
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
		var rawInput interface{}
		if err := json.Unmarshal([]byte(processedInput), &rawInput); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}

		var itemsArrayFromPath []interface{}
		switch v := rawInput.(type) {
		case []interface{}:
			if state.ItemsPath == "" || state.ItemsPath == "$" {
				itemsArrayFromPath = v
			} else {
				inputMap, ok := v[0].(map[string]interface{})
				if !ok {
					return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: "ItemsPath requires object input"}
				}
				items, err := getJSONPathValue(inputMap, state.ItemsPath)
				if err != nil {
					return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: err.Error()}
				}
				var ok2 bool
				itemsArrayFromPath, ok2 = items.([]interface{})
				if !ok2 {
					return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "items is not an array"}
				}
			}
		case map[string]interface{}:
			itemsPath := state.ItemsPath
			if itemsPath == "" {
				itemsPath = "$"
			}
			items, err := getJSONPathValue(v, itemsPath)
			if err != nil {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: err.Error()}
			}
			var ok2 bool
			itemsArrayFromPath, ok2 = items.([]interface{})
			if !ok2 {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "items is not an array"}
			}
		default:
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "input must be an array or object"}
		}

		itemsArray = itemsArrayFromPath
	}

	maxConcurrency := int(state.MaxConcurrency)
	if maxConcurrency <= 0 && state.MaxConcurrencyPath != "" && !isJSONata {
		// MaxConcurrencyPath is a reference path selecting a
		// non-negative integer from the state input.
		resolved, cerr := resolveMapConcurrencyPath(processedInput, state.MaxConcurrencyPath)
		if cerr != nil {
			return "", "", cerr
		}
		maxConcurrency = resolved
	}
	if maxConcurrency <= 0 {
		maxConcurrency = len(itemsArray)
	}

	// A configured Label replaces the state name in the Map Run ARN
	// (Distributed Map documentation: "For each Map Run, Step Functions
	// adds the label to the Map Run ARN").
	mapLabel := execCtx.CurrentState
	if state.Label != "" {
		mapLabel = state.Label
	}
	mapRunArn := generateMapRunArn(e.store, e.region, e.accountID, execCtx.Execution.ExecutionArn, mapLabel)
	now := time.Now().UTC()
	total := int64(len(itemsArray))
	mapRunRecord := &sfnstore.MapRun{
		MapRunArn:       mapRunArn,
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		StateMachineArn: execCtx.Execution.StateMachineArn,
		Name:            execCtx.CurrentState,
		Status:          "RUNNING",
		StartDate:       now.Unix(),
		ItemCounts: sfnstore.MapRunItemCounts{
			Pending: total,
			Total:   total,
		},
		ExecutionCounts: sfnstore.MapRunExecutionCounts{
			Pending: total,
			Total:   total,
		},
		MaxConcurrency: int64(maxConcurrency),
	}
	if execCtx.IsRedrive {
		existingRuns, _ := e.store.ListMapRunsByExecution(ctx, execCtx.Execution.ExecutionArn)
		for _, er := range existingRuns {
			if er.Name == execCtx.CurrentState && len(er.CompletedResults) > 0 {
				mapRunRecord = er
				mapRunRecord.Status = "RUNNING"
				mapRunRecord.StopDate = 0
				break
			}
		}
	}
	if mapRunRecord.MapRunArn != "" && !execCtx.IsRedrive {
		if err := e.store.CreateMapRun(ctx, mapRunRecord); err != nil {
			logs.Warn("failed to create map run record", logs.Err(err))
		}
	} else if mapRunRecord.MapRunArn != "" {
		if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
			logs.Warn("failed to update map run record for redrive", logs.Err(err))
		}
	}
	if mapRunRecord.CompletedResults == nil {
		mapRunRecord.CompletedResults = make(map[int]string)
	}

	defer func() {
		mapRunRecord.StopDate = time.Now().UTC().Unix()
		if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
			logs.Error("sfn: failed to update map run status", logs.Err(err))
		}
	}()

	// Parameters is the pre-ItemSelector name of the per-item payload
	// template ("ItemSelector replaces Parameters in a Map state"), so a
	// definition written against the legacy name feeds the same path.
	itemSelector := state.ItemSelector
	if itemSelector == nil && state.Parameters != nil {
		itemSelector = state.Parameters.Values
	}
	processedItems := make([]interface{}, len(itemsArray))
	if itemSelector != nil {
		for i, item := range itemsArray {
			execCtx.MapItemIndex = i
			execCtx.MapItemValue = item
			var selected interface{}
			if isJSONata {
				// JSONata template failures are query evaluation errors.
				resolved, err := e.applyItemSelector(ctx, execCtx, itemSelector, item)
				if err != nil {
					return "", "", e.newQueryEvalError(ctx, execCtx, "ItemSelector", err.Error())
				}
				selected = resolved
			} else {
				// JSONPath context failures are runtime errors; the
				// classifier lives in applyItemSelectorJSONPath.
				resolved, evalErr := e.applyItemSelectorJSONPath(itemSelector, item)
				if evalErr != nil {
					return "", "", evalErr
				}
				selected = resolved
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

	// The ItemBatcher groups consecutive items into the work units each
	// child workflow execution receives; without one every item is its
	// own unit carrying the item JSON.
	units, batchErr := e.buildMapWorkUnits(ctx, execCtx, state, processedInput, processedItems, itemsArray)
	if batchErr != nil {
		return "", "", batchErr
	}
	if state.ItemBatcher != nil && execCtx.AfterItemBatcher == nil {
		unitInputs := make([]string, len(units))
		for i, u := range units {
			unitInputs[i] = u.InputJSON
		}
		joined := "[" + strings.Join(unitInputs, ",") + "]"
		execCtx.AfterItemBatcher = &joined
	}

	var wg sync.WaitGroup
	results := make([]string, len(units))
	errors := make([]error, len(units))
	itemsProcessed := int64(0)
	itemsFailed := int64(0)
	var mu sync.Mutex

	for i, unit := range units {
		if cached, ok := mapRunRecord.CompletedResults[unit.StartIndex]; ok && cached != "" {
			results[i] = cached
			itemsProcessed += int64(unit.ItemCount)
		}
	}

	// Distributed mode dispatches every unit as its own child workflow
	// execution; the collected child identities feed the ResultWriter
	// export records.
	distributed := isDistributedMap(state)
	childMetas := make([]mapChildMeta, len(units))

	sem := make(chan struct{}, maxConcurrency)

	for i := range units {
		if results[i] != "" {
			continue
		}
		wg.Add(1)
		go func(slot int, unit mapWorkUnit) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			var childExec *sfnstore.Execution
			var childEventId *int64

			defer func() {
				if r := recover(); r != nil {
					logs.Error("sfn: panic in map worker", logs.Int("index", unit.StartIndex), logs.Any("panic", r), logs.String("stack", string(debug.Stack())))
					mu.Lock()
					errors[slot] = fmt.Errorf("internal panic: %v", r)
					itemsFailed += int64(unit.ItemCount)
					mu.Unlock()
					if childExec != nil {
						e.finishMapChildExecution(ctx, childExec, childEventId, "", &ExecutionError{ErrorCode: "States.Runtime", Cause: fmt.Sprintf("internal panic: %v", r)})
					}
				}
			}()

			iteratorStates, err := e.extractStatesFromDefinition(state.GetIterator())
			if err != nil {
				mu.Lock()
				defer mu.Unlock()
				errors[slot] = err
				itemsFailed += int64(unit.ItemCount)
				return
			}
			unitExecution := execCtx.Execution
			unitEventId := execCtx.EventId
			if distributed {
				childExec = e.beginMapChildExecution(ctx, execCtx, state, mapRunArn, unit, slot)
				if childExec != nil {
					one := int64(1)
					childEventId = &one
					unitExecution = childExec
					unitEventId = childEventId
				}
			}
			iteratorCtx := &ExecutionContext{
				Execution:     unitExecution,
				Definition:    state.GetIterator(),
				CurrentState:  state.GetIterator().StartAt,
				Input:         unit.InputJSON,
				Output:        "",
				EventId:       unitEventId,
				States:        iteratorStates,
				QueryLanguage: execCtx.QueryLanguage,
				VariableScope: execCtx.VariableScope.NewChild(),
				MapItemIndex:  unit.StartIndex,
				MapItemValue:  unit.ContextItem,
			}
			execErr := e.executeStates(ctx, iteratorCtx)
			if childExec != nil {
				var childFailure *ExecutionError
				if execErr != nil {
					if stateErr, ok := execErr.(*ExecutionError); ok {
						childFailure = stateErr
					} else {
						childFailure = &ExecutionError{ErrorCode: "States.Runtime", Cause: execErr.Error()}
					}
				}
				e.finishMapChildExecution(ctx, childExec, childEventId, iteratorCtx.Output, childFailure)
				mu.Lock()
				childMetas[slot] = mapChildMeta{Arn: childExec.ExecutionArn, Name: childExec.Name, RedriveCount: childExec.RedriveCount}
				mu.Unlock()
			}
			mu.Lock()
			defer mu.Unlock()
			errors[slot] = execErr
			if execErr == nil {
				results[slot] = iteratorCtx.Output
				itemsProcessed += int64(unit.ItemCount)
				mapRunRecord.CompletedResults[unit.StartIndex] = iteratorCtx.Output
			} else {
				itemsFailed += int64(unit.ItemCount)
			}
		}(i, units[i])
	}

	wg.Wait()

	totalItems := mapRunRecord.ItemCounts.Total
	// The item counters roll up items while the execution counters roll up
	// child workflow executions: with an ItemBatcher one unit covers
	// several items, without one the two are the same population.
	executionsSucceeded := int64(0)
	for _, err := range errors {
		if err == nil {
			executionsSucceeded++
		}
	}
	mapRunRecord.ItemCounts.Succeeded = itemsProcessed
	mapRunRecord.ItemCounts.Failed = itemsFailed
	mapRunRecord.ItemCounts.Running = 0
	mapRunRecord.ItemCounts.Pending = 0
	mapRunRecord.ItemCounts.Total = totalItems
	mapRunRecord.ExecutionCounts.Succeeded = executionsSucceeded
	mapRunRecord.ExecutionCounts.Failed = int64(len(errors)) - executionsSucceeded
	mapRunRecord.ExecutionCounts.Running = 0
	mapRunRecord.ExecutionCounts.Pending = 0
	mapRunRecord.ExecutionCounts.Total = int64(len(units))
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

	// Distributed mode applies the tolerated-failure thresholds: a Map Run
	// fails with States.ExceedToleratedFailureThreshold only when the
	// failed-item count or percentage exceeds the configured threshold,
	// and succeeds otherwise with the failed iterations tolerated.
	tolerated, exceeded := e.evaluateToleratedFailure(execCtx, state, itemsFailed, totalItems)
	if firstError != nil && !tolerated {
		failCode := "States.IteratorFailed"
		if exceeded {
			failCode = "States.ExceedToleratedFailureThreshold"
		}
		mapRunRecord.Status = "FAILED"
		if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
			logs.Warn("failed to update map run to FAILED", logs.Err(err))
		}

		if len(state.Catch) > 0 {
			catchPolicy := e.findMatchingCatchPolicy(state.Catch, failCode)
			if catchPolicy != nil {
				isJSONataCatch := IsJSONataState(state, execCtx.QueryLanguage)
				if isJSONataCatch {
					return e.executeMapJSONataCatch(ctx, execCtx, state, processedInput, failCode, firstError.Error(), catchPolicy)
				}
				catchOutput := e.buildCatchOutput(processedInput, failCode, firstError.Error(), catchPolicy.ResultPath)
				return catchOutput, catchPolicy.Next, nil
			}
		}
		return "", "", &ExecutionError{ErrorCode: failCode, Cause: firstError.Error()}
	}

	mapRunRecord.Status = "SUCCEEDED"
	if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
		logs.Warn("Failed to update map run status", logs.Err(err))
	}

	// Failed iterations tolerated by the configured thresholds leave null
	// placeholders in the assembled result array; a successful iteration
	// always produces JSON text, so an empty slot marks a failure.
	rendered := make([]string, len(results))
	for i, r := range results {
		if r == "" {
			rendered[i] = "null"
		} else {
			rendered[i] = r
		}
	}
	output := fmt.Sprintf(`[%s]`, strings.Join(rendered, ","))

	// A configured ResultWriter exports the per-unit execution records to
	// S3 and replaces the raw result with the Map Run ARN plus the export
	// location, before ResultPath and OutputPath processing. A unit input
	// is the item JSON or the batched {"Items": [...]} payload.
	if state.ResultWriter != nil {
		unitInputs := make([]string, len(units))
		for i, unit := range units {
			unitInputs[i] = unit.InputJSON
		}
		exported, werr := e.writeMapResultWriter(ctx, execCtx, state, mapRunRecord, rendered, errors, unitInputs, childMetas)
		if werr != nil {
			return "", "", werr
		}
		output = exported
	}

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
		if state.ResultSelector != nil {
			// The selector receives the Map result array as its data,
			// before ResultPath folds it into the state input.
			selected, selErr := e.applyResultSelector(output, state.ResultSelector, "")
			if selErr != nil {
				return "", "", selErr
			}
			output = selected
		}
		if state.ResultPath != "" {
			output = e.applyResultPath(processedInput, output, state.ResultPath)
		}
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
