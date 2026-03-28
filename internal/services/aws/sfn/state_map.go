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

func (e *Executor) executeMap(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState) (string, string, *ExecutionError) {
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

	itemsArray, ok := items.([]interface{})
	if !ok {
		return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "items is not an array"}
	}

	maxConcurrency := int(state.MaxConcurrency)
	if maxConcurrency <= 0 {
		maxConcurrency = len(itemsArray)
	}

	mapRunArn := generateMapRunArn(e.region, e.accountID, execCtx.Execution.ExecutionArn, execCtx.CurrentState)
	now := time.Now().UTC()
	mapRunRecord := map[string]interface{}{
		"mapRunArn":           mapRunArn,
		"executionArn":        execCtx.Execution.ExecutionArn,
		"stateMachineArn":     execCtx.Execution.StateMachineArn,
		"name":                execCtx.CurrentState,
		"status":              "RUNNING",
		"startDate":           now.Unix(),
		"itemCount":           int64(len(itemsArray)),
		"maxConcurrency":      int64(maxConcurrency),
		"itemsProcessedCount": int64(0),
		"itemsFailedCount":    int64(0),
		"itemsCancelledCount": int64(0),
	}
	mapRunsMu.Lock()
	mapRuns[execCtx.Execution.ExecutionArn] = append(mapRuns[execCtx.Execution.ExecutionArn], mapRunRecord)
	mapRunsMu.Unlock()

	defer func() {
		mapRunsMu.Lock()
		defer mapRunsMu.Unlock()
		runs, ok := mapRuns[execCtx.Execution.ExecutionArn]
		if !ok {
			return
		}
		for i, r := range runs {
			if arn, ok := r["mapRunArn"].(string); ok && arn == mapRunArn {
				runs[i]["stopDate"] = time.Now().UTC().Unix()
				return
			}
		}
	}()

	var wg sync.WaitGroup
	results := make([]string, len(itemsArray))
	errors := make([]error, len(itemsArray))
	itemsProcessed := int64(0)
	itemsFailed := int64(0)
	var mu sync.Mutex

	sem := make(chan struct{}, maxConcurrency)

	for i, item := range itemsArray {
		wg.Add(1)
		go func(idx int, itemValue interface{}) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			itemJSON, err := json.Marshal(itemValue)
			if err != nil {
				mu.Lock()
				errors[idx] = fmt.Errorf("failed to marshal map item: %w", err)
				itemsFailed++
				mu.Unlock()
				return
			}
			iteratorStates, err := e.extractStatesFromDefinition(state.Iterator)
			if err != nil {
				mu.Lock()
				defer mu.Unlock()
				errors[idx] = err
				itemsFailed++
				return
			}
			iteratorCtx := &ExecutionContext{
				Execution:    execCtx.Execution,
				Definition:   state.Iterator,
				CurrentState: state.Iterator.StartAt,
				Input:        string(itemJSON),
				Output:       "",
				EventId:      execCtx.EventId,
				States:       iteratorStates,
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
		}(i, item)
	}

	wg.Wait()

	mapRunsMu.Lock()
	if runs, ok := mapRuns[execCtx.Execution.ExecutionArn]; ok {
		for i, r := range runs {
			if arn, ok := r["mapRunArn"].(string); ok && arn == mapRunArn {
				runs[i]["itemsProcessedCount"] = itemsProcessed
				runs[i]["itemsFailedCount"] = itemsFailed
			}
		}
	}
	mapRunsMu.Unlock()

	var firstError error
	for _, err := range errors {
		if err != nil {
			firstError = err
			break
		}
	}

	if firstError != nil {
		mapRunsMu.Lock()
		if runs, ok := mapRuns[execCtx.Execution.ExecutionArn]; ok {
			for i, r := range runs {
				if arn, ok := r["mapRunArn"].(string); ok && arn == mapRunArn {
					runs[i]["status"] = "FAILED"
				}
			}
		}
		mapRunsMu.Unlock()

		if len(state.Catch) > 0 {
			catchPolicy := e.findMatchingCatchPolicy(state.Catch, "States.IteratorFailed")
			if catchPolicy != nil {
				catchOutput := e.buildCatchOutput(processedInput, "States.IteratorFailed", firstError.Error(), catchPolicy.ResultPath)
				return catchOutput, catchPolicy.Next, nil
			}
		}
		return "", "", &ExecutionError{ErrorCode: "States.IteratorFailed", Cause: firstError.Error()}
	}

	mapRunsMu.Lock()
	if runs, ok := mapRuns[execCtx.Execution.ExecutionArn]; ok {
		for i, r := range runs {
			if arn, ok := r["mapRunArn"].(string); ok && arn == mapRunArn {
				runs[i]["status"] = "SUCCEEDED"
			}
		}
	}
	mapRunsMu.Unlock()

	output := fmt.Sprintf(`[%s]`, strings.Join(results, ","))
	output = e.applyOutputPath(output, state.GetOutputPath())

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
