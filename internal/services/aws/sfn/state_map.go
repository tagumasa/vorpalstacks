package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeMap(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState) (string, string, *ExecutionError) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	processedInput, ipErr := e.applyInputPath(execCtx, execCtx.Input, state.GetInputPath())
	if ipErr != nil {
		return "", "", ipErr
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:             execCtx.Execution.ExecutionArn,
		EventId:                  eventId,
		PreviousEventId:          eventId - 1,
		Type:                     "MapStateEntered",
		Timestamp:                time.Now().UTC(),
		StateEnteredEventDetails: stateEnteredDetails(execCtx, processedInput),
	})

	var itemsArray []interface{}
	var readerArgs itemReaderArgs

	if state.ItemReader != nil {
		readerItems, rargs, rerr := e.readItemReaderItems(ctx, execCtx, state, execCtx.MapItemReaderData)
		if rerr != nil {
			// ItemReader failures surface the documented
			// States.ItemReaderFailed error and pass through the
			// state's Catch handlers like any other runtime error.
			if len(state.Catch) > 0 {
				if catchPolicy := e.findMatchingCatchPolicy(state.Catch, rerr.ErrorCode); catchPolicy != nil {
					if IsJSONataState(state, execCtx.QueryLanguage) {
						return e.executeJSONataCatch(ctx, execCtx, processedInput, rerr.ErrorCode, rerr.Cause, catchPolicy)
					}
					catchOutput, coErr := e.buildCatchOutput(processedInput, rerr.ErrorCode, rerr.Cause, catchPolicy.ResultPath)
					if coErr != nil {
						return "", "", coErr
					}
					return catchOutput, catchPolicy.Next, nil
				}
			}
			return "", "", rerr
		}
		itemsArray = readerItems
		readerArgs = rargs
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
			// "The Map state Items field will accept a JSON array, a
			// JSON object, or a JSONata expression that must evaluate to
			// an array or object": an array iterates its elements, an
			// object its key-value pairs.
			itemised, ierr := jsonDocumentItems(resolved)
			if ierr != nil {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "Items must evaluate to an array or object"}
			}
			itemsArray = itemised
		} else {
			// With ItemReader omitted the dataset is the JSON data from
			// the previous step: "Step Functions will iterate directly
			// over the elements of an array, or the key-value pairs of a
			// JSON object."
			itemised, ierr := jsonDocumentItems(inputData)
			if ierr != nil {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "input must be an array or object"}
			}
			itemsArray = itemised
		}
	} else {
		var rawInput interface{}
		if err := json.Unmarshal([]byte(processedInput), &rawInput); err != nil {
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
		}

		// selectItemsPathValue resolves the ItemsPath selection across the
		// input, the context object and the workflow variables (ItemsPath
		// is one of the documented $$. fields).
		selectItemsPathValue := func(path string) (interface{}, error) {
			if isContextPath(path) {
				return e.getContextValue(execCtx, "", path)
			}
			if name, rest, isVar := parseVariableReference(path); isVar {
				v, found, err := e.resolveVariableRef(execCtx, name, rest)
				if err != nil {
					return nil, err
				}
				if !found {
					return nil, fmt.Errorf("variable path %s selected no value", path)
				}
				return v, nil
			}
			dataMap, ok := rawInput.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("ItemsPath requires object input")
			}
			return getJSONPathValue(dataMap, path)
		}

		var itemsArrayFromPath []interface{}
		switch v := rawInput.(type) {
		case []interface{}:
			if state.ItemsPath == "" || state.ItemsPath == "$" {
				itemsArrayFromPath = v
			} else if isContextOrVariablePath(state.ItemsPath) {
				selected, serr := selectItemsPathValue(state.ItemsPath)
				if serr != nil {
					return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: serr.Error()}
				}
				itemised, ierr := jsonDocumentItems(selected)
				if ierr != nil {
					return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "items is not an array or object"}
				}
				itemsArrayFromPath = itemised
			} else {
				// ASL defines ItemsPath for object input; an array input
				// has no object to resolve the path against.
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: "ItemsPath requires object input"}
			}
		case map[string]interface{}:
			itemsPath := state.ItemsPath
			if itemsPath == "" {
				itemsPath = "$"
			}
			items, err := selectItemsPathValue(itemsPath)
			if err != nil {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItemsPath", Cause: err.Error()}
			}
			// "If the input to the Map state is a JSON object, it runs an
			// iteration for each key-value pair in the object, passing
			// the pair to the iteration as input": a selected array
			// iterates its elements, a selected object its pairs.
			itemised, ierr := jsonDocumentItems(items)
			if ierr != nil {
				return "", "", &ExecutionError{ErrorCode: "States.InvalidItems", Cause: "items is not an array or object"}
			}
			itemsArrayFromPath = itemised
		default:
			return "", "", &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "input must be an array or object"}
		}

		itemsArray = itemsArrayFromPath
	}

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "MapStateStarted",
		Timestamp:       time.Now().UTC(),
		MapStateStartedEventDetails: &sfnstore.MapStateStartedEventDetails{
			Length: int64(len(itemsArray)),
		},
	})

	maxConcurrency := 0
	if f, isNumber := toFloat64(state.MaxConcurrency); isNumber {
		maxConcurrency = int(f)
	} else if state.MaxConcurrency != nil {
		// "In JSONata states, you can specify a JSONata expression that
		// evaluates to an integer" (Distributed-mode MaxConcurrency).
		resolved, _, cerr := e.resolveMapNumericMember(ctx, execCtx, "MaxConcurrency", state.MaxConcurrency, processedInput, 0, 0)
		if cerr != nil {
			return "", "", cerr
		}
		maxConcurrency = int(resolved)
	}
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
		// "The default value is 0, which places no limit on concurrency.
		// Step Functions invokes iterations as concurrently as possible."
		maxConcurrency = len(itemsArray)
	}

	// The tolerated-failure thresholds resolve against the processed input
	// (the same stage MaxConcurrency resolves at), before the run starts —
	// a failing threshold expression surfaces its evaluation error rather
	// than degrading to an unconfigured threshold mid-run.
	thresholds, terr := e.resolveToleratedFailureThresholds(ctx, execCtx, state, processedInput)
	if terr != nil {
		return "", "", terr
	}

	// A configured Label replaces the state name in the Map Run ARN
	// (Distributed Map documentation: "For each Map Run, Step Functions
	// adds the label to the Map Run ARN").
	mapLabel := execCtx.CurrentState
	if state.Label != "" {
		mapLabel = state.Label
	}
	distributed := isDistributedMap(state)
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
	// An inline map keeps its record only as the redrive checkpoint of the
	// completed iterations; it stays invisible on the Map Run API surface.
	mapRunRecord.Inline = !distributed
	// The run belongs to the state's current retry attempt ("When you
	// retry a Map state, it creates a new Map Run"): the attempt stamps
	// the child-execution names so this run's children are fresh
	// executions, not reclamations of a prior attempt's.
	mapRunRecord.Attempt = int64(execCtx.RetryCount)
	reclaimedRun := false
	if execCtx.IsRedrive {
		// A redrive reuses the prior Map Run of this state regardless of how
		// many iterations completed — the run's ARN is its identity across
		// redrives, so a zero-result run must be reclaimed, not replaced by
		// a fresh record (the run list would show two runs for one
		// execution and the lifecycle event would say MapRunStarted). The
		// run identifiers carry a monotonic sequence, so the run with the
		// highest sequence is the most recent cycle of the state — the list
		// order itself cannot be trusted, because a lexicographic key order
		// sorts mapRun-12 before mapRun-9; an unparsable sequence falls
		// back to the latest StartDate.
		existingRuns, lerr := e.store.ListMapRunsByExecution(ctx, execCtx.Execution.ExecutionArn)
		if lerr != nil {
			// A failed run lookup must not degrade into a fresh run: the
			// prior run's succeeded units would then re-dispatch against
			// their already-terminal children. The redrive fails loudly
			// instead and can be retried once the store read works.
			return "", "", &ExecutionError{ErrorCode: "States.Runtime", Cause: fmt.Sprintf("failed to list prior map runs during redrive: %v", lerr)}
		}
		var best *sfnstore.MapRun
		var bestSeq int64
		for _, er := range existingRuns {
			if er.Name != execCtx.CurrentState {
				continue
			}
			seq, _ := sfnstore.MapRunSeqFromARN(er.MapRunArn)
			if best == nil || seq > bestSeq || (seq == bestSeq && er.StartDate > best.StartDate) {
				best, bestSeq = er, seq
			}
		}
		if best != nil {
			mapRunRecord = best
			mapRunRecord.Status = "RUNNING"
			mapRunRecord.StopDate = 0
			// "The redrive count for a redriven Map Run is always greater
			// than 0" — the reclaim itself marks the redrive, and the date
			// records when it last happened.
			mapRunRecord.RedriveCount++
			mapRunRecord.RedriveDate = time.Now().UTC().Unix()
			reclaimedRun = true
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
	// The reclaimed run's ARN is the identity the dispatched units must
	// carry; the freshly generated candidate is discarded when a prior run
	// was reclaimed.
	mapRunArn = mapRunRecord.MapRunArn
	if distributed && mapRunRecord.MapRunArn != "" {
		// The Map Run's own lifecycle events land in the parent history: a
		// fresh run starts, a redriven run marks the redrive. Inline maps
		// have no Map Run on AWS, so they record none of these.
		eventId = execCtx.nextEventId()
		evt := &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    execCtx.Execution.ExecutionArn,
			EventId:         eventId,
			PreviousEventId: eventId - 1,
			Timestamp:       time.Now().UTC(),
		}
		if reclaimedRun {
			evt.Type = "MapRunRedriven"
			evt.MapRunRedrivenEventDetails = &sfnstore.MapRunRedrivenEventDetails{
				MapRunArn:    mapRunRecord.MapRunArn,
				RedriveCount: mapRunRecord.RedriveCount,
			}
		} else {
			evt.Type = "MapRunStarted"
			evt.MapRunStartedEventDetails = &sfnstore.MapRunStartedEventDetails{
				MapRunArn: mapRunRecord.MapRunArn,
			}
		}
		e.logHistoryEvent(ctx, execCtx.Execution, evt)
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
	// The Map context reports each item's provenance as Map.Item.Source,
	// derived from the same resolved reader arguments the dataset was read
	// through — one interpretation of the reader parameters.
	itemSource := itemSourceFromReader(state.ItemReader, readerArgs)
	execCtx.MapItemSource = itemSource
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
				resolved, evalErr := e.applyItemSelectorJSONPath(execCtx, itemSelector, item)
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
	execCtx.MapItemSource = ""

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
	childMetas := make([]mapChildMeta, len(units))

	// Every unit runs the same iterator definition, so the state table is
	// extracted once here rather than re-parsed inside each worker.
	iteratorStates, parseErr := extractStatesFromDefinition(state.GetIterator())

	sem := make(chan struct{}, maxConcurrency)

	for i := range units {
		if results[i] != "" {
			continue
		}
		if parseErr != nil {
			// The parse failed identically for every unit before any
			// observable work; record it per unit so the shared
			// aggregation below owns the terminal outcome, exactly as the
			// per-worker parse did.
			errors[i] = parseErr
			itemsFailed += int64(units[i].ItemCount)
			continue
		}
		wg.Add(1)
		go func(slot int, unit mapWorkUnit) {
			defer wg.Done()
			if distributed && execCtx.IsRedrive {
				// A redriven unit queued at the Map Run's concurrency limit
				// parks its prior child execution as PENDING_REDRIVE for as
				// long as it waits: ListExecutions(mapRunArn,
				// PENDING_REDRIVE) must surface the queued children, and the
				// reclaim below flips them back to RUNNING in the order
				// slots free up.
				select {
				case sem <- struct{}{}:
				default:
					e.parkMapChildPendingRedrive(ctx, execCtx, state, slot, mapRunRecord.Attempt)
					sem <- struct{}{}
				}
			} else {
				sem <- struct{}{}
			}
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
						e.finishMapChildExecution(ctx, childExec, childEventId, "", &ExecutionError{ErrorCode: "States.Runtime", Cause: fmt.Sprintf("internal panic: %v", r)}, false)
					}
				}
			}()

			unitExecution := execCtx.Execution
			unitEventId := execCtx.EventId
			unitStartState := state.GetIterator().StartAt
			unitInput := unit.InputJSON
			if distributed {
				child, base, childResume, derr := e.beginMapChildExecution(ctx, execCtx, state, mapRunArn, unit, slot, state.GetIterator(), mapRunRecord.Attempt)
				if derr != nil {
					// A dispatch failure fails the unit — the unit's
					// topology is distributed by definition, so silently
					// running it inside the parent history would change the
					// observable execution topology on a transient error.
					mu.Lock()
					errors[slot] = derr
					itemsFailed += int64(unit.ItemCount)
					mu.Unlock()
					return
				}
				childExec = child
				childEventId = &base
				unitExecution = childExec
				unitEventId = childEventId
				if childResume != nil && childResume.StateName != "" {
					// A reclaimed child continues from its own failed
					// state with the recorded state's output as input;
					// an empty input means the original unit input.
					unitStartState = childResume.StateName
					if childResume.Input != "" {
						unitInput = childResume.Input
					}
				}
			}
			iteratorCtx := &ExecutionContext{
				Execution:     unitExecution,
				Definition:    state.GetIterator(),
				CurrentState:  unitStartState,
				Input:         unitInput,
				Output:        "",
				EventId:       unitEventId,
				States:        iteratorStates,
				QueryLanguage: execCtx.QueryLanguage,
				VariableScope: execCtx.VariableScope.NewChild(),
				MapItemIndex:  unit.StartIndex,
				MapItemValue:  unit.ContextItem,
				MapItemSource: itemSource,
			}
			// The iteration's variables are scoped to the iteration: the
			// bytes return to the execution budget once it completes.
			defer iteratorCtx.VariableScope.Release()
			if childExec == nil {
				// Inline iterations run inside the parent history, so each
				// one is framed by its MapIteration pair; distributed units
				// carry their own child-execution histories instead. The
				// iteration's name is "the name of the iteration's parent
				// Map state" (MapIterationEventDetails) — never the Label,
				// which only the Map Run ARN carries.
				eventId := execCtx.nextEventId()
				e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
					ExecutionArn:    execCtx.Execution.ExecutionArn,
					EventId:         eventId,
					PreviousEventId: eventId - 1,
					Type:            "MapIterationStarted",
					Timestamp:       time.Now().UTC(),
					MapIterationEventDetails: &sfnstore.MapIterationEventDetails{
						Index: int64(unit.StartIndex),
						Name:  execCtx.CurrentState,
					},
				})
			}
			execErr := e.executeStates(ctx, iteratorCtx)
			if childExec == nil {
				iterationTerminal := "MapIterationFailed"
				if execErr == nil {
					iterationTerminal = "MapIterationSucceeded"
				} else if ctx.Err() != nil {
					iterationTerminal = "MapIterationAborted"
				}
				eventId := execCtx.nextEventId()
				e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
					ExecutionArn:    execCtx.Execution.ExecutionArn,
					EventId:         eventId,
					PreviousEventId: eventId - 1,
					Type:            iterationTerminal,
					Timestamp:       time.Now().UTC(),
					MapIterationEventDetails: &sfnstore.MapIterationEventDetails{
						Index: int64(unit.StartIndex),
						Name:  execCtx.CurrentState,
					},
				})
			}
			if childExec != nil {
				var childFailure *ExecutionError
				childAborted := false
				if execErr != nil {
					if isCanceledError(execErr) {
						// A cancelled child is aborted, not failed: the
						// Map Run item status is Aborted when "the user
						// cancelled the execution".
						childAborted = true
					} else if stateErr, ok := execErr.(*ExecutionError); ok {
						childFailure = stateErr
					} else {
						childFailure = &ExecutionError{ErrorCode: "States.Runtime", Cause: execErr.Error()}
					}
				}
				e.finishMapChildExecution(ctx, childExec, childEventId, iteratorCtx.Output, childFailure, childAborted)
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
	tolerated, exceeded := e.evaluateToleratedFailure(state, thresholds, itemsFailed, mapRunRecord.ItemCounts.Total)
	if firstError != nil && !tolerated {
		// The threshold failure keeps its documented aggregate name; otherwise
		// the Map state reports the iteration's own error identity — an
		// undocumented aggregate would make a Catch on the iteration's error
		// never match.
		failCode := iterationErrorName(firstError)
		if exceeded {
			failCode = "States.ExceedToleratedFailureThreshold"
		}
		// A stopped run is ABORTED (MapRunStatus), matching the
		// MapRunAborted history event; everything else is FAILED.
		if isCanceledError(firstError) {
			mapRunRecord.Status = "ABORTED"
		} else {
			mapRunRecord.Status = "FAILED"
		}
		if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
			logs.Warn("failed to update map run terminal status", logs.Err(err))
		}
		if distributed {
			eventId = execCtx.nextEventId()
			evt := &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Timestamp:       time.Now().UTC(),
			}
			if isCanceledError(firstError) {
				evt.Type = "MapRunAborted"
			} else {
				evt.Type = "MapRunFailed"
				evt.MapRunFailedEventDetails = &sfnstore.MapRunFailedEventDetails{
					Error: failCode,
					Cause: firstError.Error(),
				}
			}
			e.logHistoryEvent(ctx, execCtx.Execution, evt)
		}

		// Retriers run before catchers, as on Task and Parallel; a
		// cancellation is not a failure to retry — the abort path owns the
		// history. Each attempt re-enters the state (a fresh Map Run for
		// distributed maps) and re-emits its Entered/Started pair.
		if !isCanceledError(firstError) && len(state.Retry) > 0 {
			if matchedRetry := e.findMatchingRetryPolicy(state.Retry, failCode); matchedRetry != nil && execCtx.RetryCount < matchedRetry.MaxAttempts {
				if e.sleepForRetry(ctx, matchedRetry, execCtx.RetryCount+1) {
					return "", "", executionInterruptedError(ctx, "Execution interrupted during retry")
				}
				execCtx.RetryCount++
				return e.executeMap(ctx, execCtx, state)
			}
		}

		// A catcher never consumes a cancellation either: the machine is
		// stopping, no recovery transition may run on the dead context.
		if !isCanceledError(firstError) && len(state.Catch) > 0 {
			catchPolicy := e.findMatchingCatchPolicy(state.Catch, failCode)
			if catchPolicy != nil {
				isJSONataCatch := IsJSONataState(state, execCtx.QueryLanguage)
				if isJSONataCatch {
					return e.executeJSONataCatch(ctx, execCtx, processedInput, failCode, firstError.Error(), catchPolicy)
				}
				catchOutput, coErr := e.buildCatchOutput(processedInput, failCode, firstError.Error(), catchPolicy.ResultPath)
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
		return "", "", &ExecutionError{ErrorCode: failCode, Cause: firstError.Error()}
	}

	mapRunRecord.Status = "SUCCEEDED"
	if err := e.store.UpdateMapRun(ctx, mapRunRecord); err != nil {
		logs.Warn("Failed to update map run status", logs.Err(err))
	}
	if distributed {
		eventId = execCtx.nextEventId()
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    execCtx.Execution.ExecutionArn,
			EventId:         eventId,
			PreviousEventId: eventId - 1,
			Type:            "MapRunSucceeded",
			Timestamp:       time.Now().UTC(),
		})
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

	output, outErr := e.applyContainerStateOutput(ctx, execCtx, state, processedInput, output)
	if outErr != nil {
		return "", "", outErr
	}

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:            execCtx.Execution.ExecutionArn,
		EventId:                 eventId,
		PreviousEventId:         eventId - 1,
		Type:                    "MapStateExited",
		Timestamp:               time.Now().UTC(),
		StateExitedEventDetails: stateExitedDetails(execCtx, output),
	})

	return output, state.Next, nil
}

// resolveMapNumericMember resolves a Map numeric member that is a literal
// number or, per the Distributed-mode field contracts, "In JSONata
// states, you can specify a JSONata expression that evaluates to an
// integer". max <= 0 means no upper bound. The second return reports
// whether a member was present.
func (e *Executor) resolveMapNumericMember(ctx context.Context, execCtx *ExecutionContext, field string, raw interface{}, processedInput string, min, max float64) (float64, bool, *ExecutionError) {
	switch v := raw.(type) {
	case nil:
		return 0, false, nil
	case string:
		if !IsExpression(v) {
			return 0, true, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("Map state %s must be a number or a JSONata expression", field)}
		}
		var inputData interface{}
		if processedInput != "" {
			if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
				return 0, true, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: "failed to parse input JSON"}
			}
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		vars := buildVarsMap(statesVar, execCtx.VariableScope)
		resolved, err := EvaluateJSONata(ctx, UnwrapExpression(v), nil, vars)
		if err != nil {
			return 0, true, e.newQueryEvalError(ctx, execCtx, field, err.Error())
		}
		bound := fmt.Sprintf("at least %v", min)
		if max > 0 {
			bound = fmt.Sprintf("from %v to %v", min, max)
		}
		n, ok := toFloat64(resolved)
		if !ok || n != math.Trunc(n) || n < min || (max > 0 && n > max) {
			return 0, true, e.newQueryEvalError(ctx, execCtx, field, "the expression must evaluate to an integer "+bound)
		}
		return n, true, nil
	}
	// Any other value must be a literal number (the JSON wire form is
	// float64; programmatic construction may carry any Go numeric).
	v, isNumber := toFloat64(raw)
	if !isNumber {
		return 0, true, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("Map state %s must be a number or a JSONata expression", field)}
	}
	if v < min || (max > 0 && v > max) {
		return 0, true, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("Map state %s is outside the acceptable range", field)}
	}
	return v, true, nil
}
