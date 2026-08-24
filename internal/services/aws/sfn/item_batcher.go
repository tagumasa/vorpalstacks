package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// mapWorkUnit is one dispatch unit of a Map state: a consecutive slice of
// the item array and the JSON input its ItemProcessor invocation receives.
// Without an ItemBatcher every item forms its own unit carrying the item
// JSON; with one, units carry {"Items": [...]} combined with the fixed
// BatchInput under the BatchInput key (ItemBatcher documentation).
type mapWorkUnit struct {
	// StartIndex is the index of the unit's first item in the item array;
	// it doubles as the redrive checkpoint key.
	StartIndex int
	// ItemCount is the number of items the unit covers; item counts, the
	// tolerated-failure arithmetic and the Map Run item counters operate
	// on items, never on units.
	ItemCount int
	// InputJSON is the JSON text the ItemProcessor (or child workflow
	// execution) receives for the unit.
	InputJSON string
	// ContextItem carries the original, pre-ItemSelector item for the
	// context object in item-mode units; batch-mode units have none.
	ContextItem interface{}
}

// buildMapWorkUnits groups the processed items into the dispatch units the
// ItemBatcher defines. Consecutive items accumulate while the unit stays
// within the resolved MaxItemsPerBatch and MaxInputBytesPerBatch; the byte
// cap defaults to the 256 KiB child-execution input bound when only an
// item count is configured (ItemBatcher documentation). A single item
// larger than the byte cap still forms its own unit: the batch can not be
// reduced below one item.
func (e *Executor) buildMapWorkUnits(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, processedInput string, processedItems, originalItems []interface{}) ([]mapWorkUnit, *ExecutionError) {
	if state.ItemBatcher == nil {
		units := make([]mapWorkUnit, len(processedItems))
		for i, item := range processedItems {
			itemJSON, err := json.Marshal(item)
			if err != nil {
				return nil, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to marshal map item: " + err.Error()}
			}
			units[i] = mapWorkUnit{StartIndex: i, ItemCount: 1, InputJSON: string(itemJSON)}
			if i < len(originalItems) {
				units[i].ContextItem = originalItems[i]
			}
		}
		return units, nil
	}
	ib := state.ItemBatcher

	maxItems, err := e.resolveMaxItemsPerBatch(ctx, execCtx, state, ib, processedInput)
	if err != nil {
		return nil, err
	}
	maxBytes := int64(sfnstore.MaxExecutionDataBytes)
	if ib.MaxInputBytesPerBatch != nil {
		maxBytes = *ib.MaxInputBytesPerBatch
	} else if ib.MaxInputBytesPerBatchPath != "" {
		resolved, rerr := resolveItemBatcherReference(processedInput, ib.MaxInputBytesPerBatchPath, "MaxInputBytesPerBatchPath")
		if rerr != nil {
			return nil, rerr
		}
		maxBytes, err = positiveInt64(resolved, "MaxInputBytesPerBatchPath")
		if err != nil {
			return nil, err
		}
	}

	var batchInput interface{}
	if ib.BatchInput != nil {
		batchInput = ib.BatchInput
	} else if ib.BatchInputPath != "" {
		resolved, rerr := resolveItemBatcherReference(processedInput, ib.BatchInputPath, "BatchInputPath")
		if rerr != nil {
			return nil, rerr
		}
		if _, isObject := resolved.(map[string]interface{}); !isObject {
			return nil, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "ItemBatcher BatchInputPath must resolve to a JSON object"}
		}
		batchInput = resolved
	}

	units := make([]mapWorkUnit, 0)
	current := make([]interface{}, 0, 16)
	currentBytes := int64(0)
	fixedBytes := int64(0)
	if batchInput != nil {
		if b, merr := json.Marshal(batchInput); merr == nil {
			fixedBytes = int64(len(b)) + int64(len(`"BatchInput":`))
		}
	}
	wrapperBytes := int64(len(`{"Items":[]}`))

	flush := func(startIndex int) {
		if len(current) == 0 {
			return
		}
		payload := map[string]interface{}{"Items": current}
		if batchInput != nil {
			payload["BatchInput"] = batchInput
		}
		b, merr := json.Marshal(payload)
		if merr != nil {
			b = []byte("null")
		}
		units = append(units, mapWorkUnit{
			StartIndex: startIndex,
			ItemCount:  len(current),
			InputJSON:  string(b),
		})
		current = make([]interface{}, 0, 16)
		currentBytes = 0
	}

	for i, item := range processedItems {
		itemJSON, merr := json.Marshal(item)
		if merr != nil {
			return nil, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to marshal map item: " + merr.Error()}
		}
		itemBytes := int64(len(itemJSON))
		if len(current) > 0 {
			countCapped := maxItems > 0 && int64(len(current)) >= maxItems
			bytesCapped := currentBytes+1+itemBytes > maxBytes
			if countCapped || bytesCapped {
				flush(i - len(current))
			}
		}
		if len(current) == 0 {
			// The fixed wrapper and BatchInput overhead count against the
			// byte cap of every unit.
			currentBytes = wrapperBytes + fixedBytes
		}
		if len(current) > 0 {
			currentBytes++ // the comma separating array members
		}
		current = append(current, item)
		currentBytes += itemBytes
	}
	flush(len(processedItems) - len(current))

	return units, nil
}

// resolveMaxItemsPerBatch resolves the per-unit item ceiling: a literal
// count, a reference path into the state input, or a JSONata expression
// evaluated against it. Zero means unlimited.
func (e *Executor) resolveMaxItemsPerBatch(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, ib *sfnstore.ItemBatcherConfig, processedInput string) (int64, *ExecutionError) {
	switch v := ib.MaxItemsPerBatch.(type) {
	case nil:
		if ib.MaxItemsPerBatchPath == "" {
			return 0, nil
		}
		resolved, err := resolveItemBatcherReference(processedInput, ib.MaxItemsPerBatchPath, "MaxItemsPerBatchPath")
		if err != nil {
			return 0, err
		}
		value, verr := positiveInt64(resolved, "MaxItemsPerBatchPath")
		if verr != nil {
			return 0, verr
		}
		return value, nil
	case float64:
		if v != math.Trunc(v) || v < 1 {
			return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "ItemBatcher MaxItemsPerBatch must be a positive integer"}
		}
		return int64(v), nil
	case string:
		// JSONata states may carry an expression that evaluates to a
		// positive integer (ItemBatcher documentation).
		var inputData interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
			return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		vars := buildVarsMap(statesVar, execCtx.VariableScope)
		resolved, err := ResolveTemplate(ctx, v, nil, vars)
		if err != nil {
			return 0, e.newQueryEvalError(ctx, execCtx, "MaxItemsPerBatch", err.Error())
		}
		value, verr := positiveInt64(resolved, "MaxItemsPerBatch")
		if verr != nil {
			return 0, verr
		}
		return value, nil
	default:
		return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "ItemBatcher MaxItemsPerBatch must be a positive integer or a JSONata expression"}
	}
}

// resolveItemBatcherReference resolves a reference path against the
// state's processed input.
func resolveItemBatcherReference(processedInput, path, field string) (interface{}, *ExecutionError) {
	var input interface{}
	if err := json.Unmarshal([]byte(processedInput), &input); err != nil {
		return nil, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("failed to parse input JSON for ItemBatcher %s", field)}
	}
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("ItemBatcher %s requires object input", field)}
	}
	value, err := getJSONPathValue(inputMap, path)
	if err != nil {
		return nil, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("ItemBatcher %s failed to resolve: %s", field, err.Error())}
	}
	return value, nil
}

// positiveInt64 asserts a resolved value is a positive integer.
func positiveInt64(value interface{}, field string) (int64, *ExecutionError) {
	number, ok := value.(float64)
	if !ok || number != math.Trunc(number) || number < 1 {
		return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("ItemBatcher %s must resolve to a positive integer", field)}
	}
	return int64(number), nil
}

// resolveMapConcurrencyPath resolves the MaxConcurrencyPath reference path
// against the state input; the selected field must hold a non-negative
// integer, where zero means no concurrency limit.
func resolveMapConcurrencyPath(processedInput, path string) (int, *ExecutionError) {
	value, err := resolveItemBatcherReference(processedInput, path, "MaxConcurrencyPath")
	if err != nil {
		return 0, err
	}
	number, ok := value.(float64)
	if !ok || number != math.Trunc(number) || number < 0 || number > math.MaxInt32 {
		return 0, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "MaxConcurrencyPath must resolve to a non-negative integer"}
	}
	return int(number), nil
}
