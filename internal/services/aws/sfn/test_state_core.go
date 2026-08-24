package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// This file holds the TestState Core: the in-memory single-state test
// run shared by the HTTP API and the admin console handler.

// TestStateMock carries the MockInput members: a mocked result or a mocked
// error output for the state under test, plus the response-validation mode.
type TestStateMock struct {
	Result              string
	ResultProvided      bool
	Error               string
	Cause               string
	ErrorProvided       bool
	FieldValidationMode string
}

// TestStateConfiguration carries the TestStateConfiguration members that
// shape mocked Map and Parallel runs.
type TestStateConfiguration struct {
	ErrorCausedByState       string
	MapItemReaderData        string
	MapIterationFailureCount int32
	RetrierRetryCount        int32
}

// TestStateInput carries every field the in-memory test run reads.
type TestStateInput struct {
	Definition      string
	StateName       string
	Input           string
	InspectionLevel string
	Variables       string
	RoleArn         string
	Context         string
	RevealSecrets   bool
	Mock            *TestStateMock
	StateConfig     *TestStateConfiguration
}

// testStateCore is the single entry point for TestState. The run executes
// the state in-memory; a mock (valid for Task, Map and Parallel states
// only) replaces the state invocation with the mocked result or error
// while the state's output processing and the supplied context object
// still apply.
func (s *StepFunctionService) testStateCore(ctx context.Context, store *sfnstore.StepFunctionStore, in TestStateInput) (map[string]interface{}, error) {
	if in.Definition == "" {
		return nil, NewInvalidDefinitionException("definition is required")
	}
	if in.StateName == "" {
		return nil, NewInvalidDefinitionException("stateName is required")
	}

	input := in.Input
	if input == "" {
		input = "{}"
	}

	inspectionLevel := in.InspectionLevel
	if inspectionLevel == "" {
		inspectionLevel = "INFO"
	}
	if inspectionLevel != "INFO" && inspectionLevel != "DEBUG" && inspectionLevel != "TRACE" {
		return nil, NewValidationException("inspectionLevel must be INFO, DEBUG, or TRACE, got " + inspectionLevel)
	}

	var def sfnstore.StateMachineDefinition
	if err := json.Unmarshal([]byte(in.Definition), &def); err != nil {
		return nil, NewInvalidDefinitionException("definition is not valid JSON: " + err.Error())
	}

	if def.QueryLanguage == "" {
		def.QueryLanguage = "JSONPath"
	}

	// stateName may address a top-level state or one nested inside a
	// Parallel branch or Map processor: testing a state within a Map is
	// the documented iteration-position use case.
	rawState, exists := def.States[in.StateName]
	if !exists {
		rawState, exists = findNestedState(def.States, in.StateName)
	}
	if !exists {
		return nil, NewInvalidDefinitionException(
			fmt.Sprintf("State '%s' not found in definition", in.StateName))
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, store.GetRegion())
	executor.currentStateMachine = &sfnstore.StateMachine{
		StateMachineArn: "arn:aws:states:" + store.GetRegion() + ":" + s.accountID + ":stateMachine:test-sm",
		Name:            "test-sm",
	}

	state, err := executor.parseState(in.StateName, rawState)
	if err != nil {
		return nil, NewInvalidDefinitionException(err.Error())
	}

	// A mock is defined only for the integration-backed state types, and
	// the context object may only accompany a mock. The validation mode
	// follows the MockResponseValidationMode enum (STRICT, PRESENT or
	// NONE). A mock carries either a result or an error output, never
	// both, and revealSecrets cannot combine with a mock.
	if in.Mock != nil {
		switch state.(type) {
		case *sfnstore.TaskState, *sfnstore.MapState, *sfnstore.ParallelState:
		default:
			return nil, NewValidationException("mock can only be specified for Task, Map, or Parallel states")
		}
		switch in.Mock.FieldValidationMode {
		case "", "STRICT", "PRESENT", "NONE":
		default:
			return nil, NewValidationException("fieldValidationMode must be STRICT, PRESENT, or NONE, got " + in.Mock.FieldValidationMode)
		}
		if in.Mock.ResultProvided && in.Mock.ErrorProvided {
			return nil, NewValidationException("mock cannot specify both result and errorOutput")
		}
		if in.RevealSecrets {
			return nil, NewValidationException("revealSecrets cannot be specified when a mock is specified")
		}
	}
	if in.Context != "" && in.Mock == nil {
		return nil, NewValidationException("context can only be specified when a mock is specified")
	}

	// The state configuration names a state in the same definition and
	// carries non-negative counts. The reader data substitutes the raw
	// source bytes ("as found in its original source"), so its format
	// follows the ItemReader InputType, not JSON.
	if in.StateConfig != nil {
		if in.StateConfig.ErrorCausedByState != "" {
			if _, ok := def.States[in.StateConfig.ErrorCausedByState]; !ok {
				return nil, NewValidationException("errorCausedByState must name a state in the definition: " + in.StateConfig.ErrorCausedByState)
			}
		}
		if in.StateConfig.MapIterationFailureCount < 0 || in.StateConfig.RetrierRetryCount < 0 {
			return nil, NewValidationException("mapIterationFailureCount and retrierRetryCount must be non-negative")
		}
	}

	variableScope := NewVariableScope(nil)
	if in.Variables != "" {
		var rawVars map[string]interface{}
		if err := json.Unmarshal([]byte(in.Variables), &rawVars); err != nil {
			return nil, NewInvalidDefinitionException("variables is not valid JSON: " + err.Error())
		}
		vars := make(map[string]interface{}, len(rawVars))
		for k, v := range rawVars {
			vars[strings.TrimPrefix(k, "$")] = v
		}
		if len(vars) > 0 {
			if err := variableScope.SetAll(vars); err != nil {
				return nil, NewInvalidDefinitionException("invalid variables: " + err.Error())
			}
		}
	}

	var suppliedContext map[string]interface{}
	if in.Context != "" {
		if err := json.Unmarshal([]byte(in.Context), &suppliedContext); err != nil {
			return nil, NewValidationException("context is not a valid Context object: " + err.Error())
		}
	}

	testExec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:" + store.GetRegion() + ":" + s.accountID + ":execution:test-state:" + in.StateName,
		StateMachineArn: "arn:aws:states:" + store.GetRegion() + ":" + s.accountID + ":stateMachine:test-sm",
		Name:            "TestState-" + in.StateName,
		Status:          "RUNNING",
		Input:           input,
		StartDate:       time.Now().UTC(),
	}

	eventId := int64(1)
	execCtx := &ExecutionContext{
		Execution:        testExec,
		Definition:       &def,
		CurrentState:     in.StateName,
		Input:            input,
		Output:           "",
		EventId:          &eventId,
		States:           map[string]sfnstore.State{in.StateName: state},
		QueryLanguage:    def.QueryLanguage,
		VariableScope:    variableScope,
		StateEnteredTime: time.Now().UTC(),
		MapItemIndex:     -1,
		SuppliedContext:  suppliedContext,
	}
	// The stateConfiguration reader member substitutes the ItemReader's
	// raw source bytes ("The data read by ItemReader in Distributed Map
	// states as found in its original source").
	if in.StateConfig != nil {
		execCtx.MapItemReaderData = in.StateConfig.MapItemReaderData
	}

	if in.Mock != nil {
		// NONE skips the result validation entirely; STRICT (the default)
		// and PRESENT both require valid JSON.
		if in.Mock.ResultProvided && in.Mock.FieldValidationMode != "NONE" && !json.Valid([]byte(in.Mock.Result)) {
			return nil, NewValidationException("mock result must be a valid JSON value")
		}
		return s.testStateWithMock(executor, execCtx, state, in.Mock, in.StateConfig, inspectionLevel), nil
	}

	var output string
	var execErr *ExecutionError
	var nextState string
	var runErr error

	switch st := state.(type) {
	case *sfnstore.PassState:
		output, nextState, runErr = executor.executePass(ctx, execCtx, st)
	case *sfnstore.ChoiceState:
		nextState, runErr = executor.executeChoice(ctx, execCtx, st)
		if runErr == nil {
			output = execCtx.Input
		}
	case *sfnstore.WaitState:
		output, nextState, runErr = executor.executeWait(ctx, execCtx, st)
	case *sfnstore.FailState:
		output = ""
		runErr = fmt.Errorf("%s: %s", st.Error, st.Cause)
	case *sfnstore.SucceedState:
		output, _, runErr = executor.executeSucceed(ctx, execCtx, st)
	case *sfnstore.TaskState:
		output, nextState, execErr = executor.executeTask(ctx, execCtx, st)
		if execErr != nil {
			return testStateFailureResult(execErr, inspectionLevel, execCtx, output, state, nextState), nil
		}
	case *sfnstore.ParallelState:
		output, nextState, execErr = executor.executeParallel(ctx, execCtx, st)
		if execErr != nil {
			return testStateFailureResult(execErr, inspectionLevel, execCtx, output, state, nextState), nil
		}
	case *sfnstore.MapState:
		output, nextState, execErr = executor.executeMap(ctx, execCtx, st)
		if execErr != nil {
			return testStateFailureResult(execErr, inspectionLevel, execCtx, output, state, nextState), nil
		}
	default:
		return nil, NewInvalidDefinitionException(
			fmt.Sprintf("Unsupported state type: %s", state.GetType()))
	}

	if len(execCtx.PendingAssign) > 0 && execCtx.VariableScope != nil {
		if err := execCtx.VariableScope.SetAll(execCtx.PendingAssign); err != nil {
			logs.Error("sfn: failed to set variables", logs.Err(err))
		}
	}

	result := map[string]interface{}{
		"output":    output,
		"nextState": nextState,
	}

	if runErr != nil {
		result["status"] = "FAILED"
		result["error"] = ""
		result["cause"] = runErr.Error()
	} else {
		result["status"] = "SUCCEEDED"
	}

	if inspectionLevel != "" {
		result["inspectionData"] = buildInspectionData(inspectionLevel, execCtx, output, state)
	}

	return result, nil
}

// testStateWithMock serves the mocked branch of TestState with the
// documented semantics. A mocked error drives the state's Retry and Catch
// handlers as they would during a real execution: retrierRetryCount
// pretends that many retry attempts already happened (retryIndex is the
// index of the matching Retry block and the backoff is the next
// attempt's interval), exhaustion falls through to the Catchers, and the
// response status becomes RETRIABLE, CAUGHT_ERROR or FAILED. A mocked
// result replaces the whole state result: for Map and Parallel states
// the input and output processing (ItemsPath extraction, ItemSelector
// transformation, ResultPath, OutputPath) runs without executing the
// iterations or branches. errorCausedByState attributes a Map or
// Parallel error mock to the named inner state; mapIterationFailureCount
// marks that many iterations failed so the tolerated-failure thresholds
// apply.
func (s *StepFunctionService) testStateWithMock(executor *Executor, execCtx *ExecutionContext, state sfnstore.State, mock *TestStateMock, stateConfig *TestStateConfiguration, inspectionLevel string) map[string]interface{} {
	retryCount := int32(0)
	if stateConfig != nil {
		retryCount = stateConfig.RetrierRetryCount
	}

	if mock.ErrorProvided {
		return s.testStateMockedError(executor, execCtx, state, mock, stateConfig, retryCount, inspectionLevel)
	}
	return s.testStateMockedResult(executor, execCtx, state, mock, stateConfig, inspectionLevel)
}

// testStateMockedError evaluates the Retry/Catch handling for a mocked
// error following the documented TestState examples: errorDetails carry
// retryIndex (the matching Retry block index) plus
// retryBackoffIntervalSeconds while retries remain, or catchIndex when a
// Catcher handles the error.
func (s *StepFunctionService) testStateMockedError(executor *Executor, execCtx *ExecutionContext, state sfnstore.State, mock *TestStateMock, stateConfig *TestStateConfiguration, retryCount int32, inspectionLevel string) map[string]interface{} {
	var retries []*sfnstore.RetryPolicy
	var catches []*sfnstore.CatchPolicy
	var input string
	switch st := state.(type) {
	case *sfnstore.TaskState:
		retries, catches, input = st.Retry, st.Catch, execCtx.Input
	case *sfnstore.MapState:
		retries, catches, input = st.Retry, st.Catch, execCtx.Input
	case *sfnstore.ParallelState:
		retries, catches, input = st.Retry, st.Catch, execCtx.Input
	default:
		input = execCtx.Input
	}

	result := map[string]interface{}{
		"output":    "",
		"status":    "FAILED",
		"error":     mock.Error,
		"cause":     mock.Cause,
		"nextState": "",
	}
	errorDetails := map[string]interface{}{}

	// Retry evaluation: the matching retrier's index and the next
	// attempt's backoff. Attempts already spent are retrierRetryCount
	// plus the initial attempt.
	retryIndex := -1
	for i, policy := range retries {
		if executor.errorMatchesAny(policy.ErrorEquals, mock.Error) {
			retryIndex = i
			break
		}
	}
	if retryIndex >= 0 {
		policy := retries[retryIndex]
		maxAttempts := policy.MaxAttempts
		if maxAttempts == 0 {
			maxAttempts = 3
		}
		errorDetails["retryIndex"] = retryIndex
		errorDetails["retryBackoffIntervalSeconds"] = int32(executor.calculateBackoffInterval(policy, retryCount+1).Seconds())

		attemptsSpent := retryCount + 1
		if stateConfig == nil || attemptsSpent < maxAttempts {
			result["status"] = "RETRIABLE"
			inspection := buildInspectionData(inspectionLevel, execCtx, "", state)
			inspection["errorDetails"] = errorDetails
			result["inspectionData"] = inspection
			return result
		}
	}

	// Retries exhausted (or absent): the Catchers take over.
	for i, policy := range catches {
		if executor.errorMatchesAny(policy.ErrorEquals, mock.Error) {
			result["status"] = "CAUGHT_ERROR"
			result["nextState"] = policy.Next
			result["output"] = executor.buildCatchOutput(input, mock.Error, mock.Cause, policy.ResultPath)
			errorDetails = map[string]interface{}{"catchIndex": i}
			inspection := buildInspectionData(inspectionLevel, execCtx, result["output"].(string), state)
			inspection["errorDetails"] = errorDetails
			result["inspectionData"] = inspection
			return result
		}
	}

	inspection := buildInspectionData(inspectionLevel, execCtx, "", state)
	inspection["errorDetails"] = errorDetails
	result["inspectionData"] = inspection
	return result
}

// testStateMockedResult runs the state's input and output processing
// against the mocked result without executing the state's work: Map
// iterations and Parallel branches never run, matching the documented
// "input and output processing only" contract.
func (s *StepFunctionService) testStateMockedResult(executor *Executor, execCtx *ExecutionContext, state sfnstore.State, mock *TestStateMock, stateConfig *TestStateConfiguration, inspectionLevel string) map[string]interface{} {
	raw := mock.Result
	nextState, resultPath, resultSelector := mockedStateTransitions(state)

	// Input processing feeds the inspection data: InputPath then
	// Parameters (Task) apply to the state input.
	stateInputPath := ""
	if ip, ok := state.(interface{ GetInputPath() string }); ok {
		stateInputPath = ip.GetInputPath()
	}
	afterInputPath := executor.applyInputPath(execCtx.Input, stateInputPath)
	execCtx.AfterInputPath = &afterInputPath
	if task, ok := state.(*sfnstore.TaskState); ok && task.Parameters != nil {
		if afterParams, perr := executor.applyParameters(execCtx.TaskToken, afterInputPath, task.Parameters); perr == nil {
			execCtx.AfterParameters = &afterParams
		}
	}

	// Map input processing: ItemsPath extraction and ItemSelector run
	// without the iterations; mapIterationFailureCount marks the first N
	// iterations failed and drives the tolerated-failure thresholds.
	var mapItems []interface{}
	if mapState, ok := state.(*sfnstore.MapState); ok {
		items := extractMockMapItems(execCtx, mapState)
		mapItems = items
		itemsJSON, _ := json.Marshal(items)
		afterItemsPath := string(itemsJSON)
		execCtx.AfterItemsPath = &afterItemsPath
		if mapState.ItemBatcher != nil {
			// The batching runs on the extracted items like a real
			// execution, so afterItemBatcher reports the actual unit
			// inputs each child workflow execution would receive.
			if units, berr := executor.buildMapWorkUnits(context.Background(), execCtx, mapState, execCtx.Input, items, items); berr == nil {
				unitInputs := make([]string, len(units))
				for i, u := range units {
					unitInputs[i] = u.InputJSON
				}
				afterItemBatcher := "[" + strings.Join(unitInputs, ",") + "]"
				execCtx.AfterItemBatcher = &afterItemBatcher
			}
		}
		if mapState.MaxConcurrency > 0 {
			execCtx.MaxConcurrencyValue = &mapState.MaxConcurrency
		}
		if mapState.ToleratedFailureCount != nil {
			execCtx.ToleratedFailureCountVal = mapState.ToleratedFailureCount
		}
		if mapState.ToleratedFailurePercentage != nil {
			execCtx.ToleratedFailurePctVal = mapState.ToleratedFailurePercentage
		}

		if stateConfig != nil && stateConfig.MapIterationFailureCount > 0 {
			if int(stateConfig.MapIterationFailureCount) > len(items) {
				// "The value for this field cannot exceed the number of
				// items in the input."
				return map[string]interface{}{
					"output":    "",
					"status":    "FAILED",
					"error":     "ValidationException",
					"cause":     fmt.Sprintf("mapIterationFailureCount %d exceeds the number of items %d", stateConfig.MapIterationFailureCount, len(items)),
					"nextState": "",
				}
			}
			tolerated, exceeded := executor.evaluateToleratedFailure(execCtx, mapState, int64(stateConfig.MapIterationFailureCount), int64(len(items)))
			if !tolerated {
				failCode := "States.IteratorFailed"
				if exceeded {
					failCode = "States.ExceedToleratedFailureThreshold"
				}
				for i, policy := range mapState.Catch {
					if executor.errorMatchesAny(policy.ErrorEquals, failCode) {
						caught := map[string]interface{}{
							"output":    executor.buildCatchOutput(execCtx.Input, failCode, fmt.Sprintf("%d map iterations failed", stateConfig.MapIterationFailureCount), policy.ResultPath),
							"status":    "CAUGHT_ERROR",
							"nextState": policy.Next,
						}
						inspection := buildInspectionData(inspectionLevel, execCtx, caught["output"].(string), state)
						inspection["errorDetails"] = map[string]interface{}{"catchIndex": i}
						caught["inspectionData"] = inspection
						return caught
					}
				}
				failed := map[string]interface{}{
					"output":    "",
					"status":    "FAILED",
					"error":     failCode,
					"cause":     fmt.Sprintf("%d map iterations failed", stateConfig.MapIterationFailureCount),
					"nextState": "",
				}
				failedInspection := buildInspectionData(inspectionLevel, execCtx, "", state)
				failed["inspectionData"] = failedInspection
				return failed
			}
		}
	}

	// Parallel mocked results must be an array with one element per
	// branch, in branch order.
	if parallel, ok := state.(*sfnstore.ParallelState); ok {
		var arr []interface{}
		if err := json.Unmarshal([]byte(raw), &arr); err != nil {
			return map[string]interface{}{
				"output": "", "status": "FAILED",
				"error": "ValidationException",
				"cause": "the Parallel mocked result must be a JSON array with one element per branch",
			}
		}
		if len(arr) != len(parallel.Branches) {
			return map[string]interface{}{
				"output": "", "status": "FAILED",
				"error": "ValidationException",
				"cause": fmt.Sprintf("the Parallel mocked result must have %d elements, one per branch, got %d", len(parallel.Branches), len(arr)),
			}
		}
	}

	if resultSelector != nil {
		if selected, selErr := executor.applyResultSelector(raw, resultSelector, ""); selErr == nil {
			raw = selected
			execCtx.AfterResultSelector = &raw
		}
	}
	output := executor.applyResultPath(execCtx.Input, raw, resultPath)
	execCtx.AfterResultPath = &output
	execCtx.Output = output

	_ = mapItems

	result := map[string]interface{}{
		"output":    output,
		"status":    "SUCCEEDED",
		"nextState": nextState,
	}
	result["inspectionData"] = buildInspectionData(inspectionLevel, execCtx, output, state)
	return result
}

// errorMatchesAny reports whether the error code matches any of the
// ErrorEquals patterns.
func (e *Executor) errorMatchesAny(patterns []string, errorCode string) bool {
	for _, pattern := range patterns {
		if e.errorMatchesPattern(errorCode, pattern) {
			return true
		}
	}
	return false
}

// extractMockMapItems extracts the item array a mocked Map test would
// process: the ItemReader dataset when configured ( honouring the raw
// reader override) or the ItemsPath/Items selection from the input.
func extractMockMapItems(execCtx *ExecutionContext, state *sfnstore.MapState) []interface{} {
	if state.ItemReader != nil {
		executor := NewExecutor(nil, nil)
		items, ierr := executor.readItemReaderItems(context.Background(), execCtx, state, execCtx.MapItemReaderData)
		if ierr != nil {
			return nil
		}
		return items
	}

	var rawInput interface{}
	if err := json.Unmarshal([]byte(execCtx.Input), &rawInput); err != nil {
		return nil
	}
	itemsPath := state.ItemsPath
	if itemsPath == "" {
		itemsPath = "$"
	}
	switch v := rawInput.(type) {
	case []interface{}:
		if itemsPath == "$" {
			return v
		}
	case map[string]interface{}:
		if items, err := getJSONPathValue(v, itemsPath); err == nil {
			if arr, ok := items.([]interface{}); ok {
				return arr
			}
			if obj, ok := items.(map[string]interface{}); ok {
				out := make([]interface{}, 0, len(obj))
				for k, val := range obj {
					out = append(out, map[string]interface{}{"key": k, "value": val})
				}
				return out
			}
		}
	}
	return nil
}

// findNestedState searches Parallel branches and Map processors for the
// named state so TestState can address states inside flow states.
func findNestedState(states map[string]interface{}, name string) (interface{}, bool) {
	for _, raw := range states {
		stateMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		for _, container := range []string{"Branches", "Iterator", "ItemProcessor"} {
			switch members := stateMap[container].(type) {
			case []interface{}:
				for _, member := range members {
					if branch, ok := member.(map[string]interface{}); ok {
						if nested, ok := branch["States"].(map[string]interface{}); ok {
							if found, ok := nested[name]; ok {
								return found, true
							}
							if found, ok := findNestedState(nested, name); ok {
								return found, true
							}
						}
					}
				}
			case map[string]interface{}:
				if nested, ok := members["States"].(map[string]interface{}); ok {
					if found, ok := nested[name]; ok {
						return found, true
					}
					if found, ok := findNestedState(nested, name); ok {
						return found, true
					}
				}
			}
		}
	}
	return nil, false
}

// mockedStateTransitions reads the transition and output-processing
// fields shared by the mockable state types.
func mockedStateTransitions(state sfnstore.State) (nextState, resultPath string, resultSelector *sfnstore.ResultSelector) {
	switch st := state.(type) {
	case *sfnstore.TaskState:
		return st.Next, st.ResultPath, st.ResultSelector
	case *sfnstore.MapState:
		return st.Next, st.ResultPath, nil
	case *sfnstore.ParallelState:
		return st.Next, st.ResultPath, nil
	}
	return "", "", nil
}

// testStateFailureResult renders the FAILED branch of the TestState
// response with inspection data.
func testStateFailureResult(execErr *ExecutionError, inspectionLevel string, execCtx *ExecutionContext, output string, state sfnstore.State, nextState string) map[string]interface{} {
	errResult := map[string]interface{}{
		"output":    "",
		"status":    "FAILED",
		"error":     execErr.ErrorCode,
		"cause":     execErr.Cause,
		"nextState": nextState,
	}
	if inspectionLevel != "" {
		errResult["inspectionData"] = buildInspectionData(inspectionLevel, execCtx, output, state)
	}
	return errResult
}

func buildInspectionData(inspectionLevel string, execCtx *ExecutionContext, output string, state sfnstore.State) map[string]interface{} {
	// The InspectionData shape exposes the state's raw result as "result";
	// there is no "output" member.
	data := map[string]interface{}{
		"input":  execCtx.Input,
		"result": output,
	}

	if inspectionLevel == "DEBUG" || inspectionLevel == "TRACE" {
		if execCtx.VariableScope != nil {
			allVars := execCtx.VariableScope.GetAll()
			if len(allVars) > 0 {
				varsJSON, _ := json.Marshal(allVars)
				data["variables"] = string(varsJSON)
			}
		}

		if execCtx.AfterArguments != nil {
			data["afterArguments"] = *execCtx.AfterArguments
		}

		if execCtx.AfterItemSelector != nil {
			data["afterItemSelector"] = *execCtx.AfterItemSelector
		}

		if execCtx.AfterInputPath != nil {
			data["afterInputPath"] = *execCtx.AfterInputPath
		}
		if execCtx.AfterParameters != nil {
			data["afterParameters"] = *execCtx.AfterParameters
		}
		if execCtx.AfterResultSelector != nil {
			data["afterResultSelector"] = *execCtx.AfterResultSelector
		}
		if execCtx.AfterResultPath != nil {
			data["afterResultPath"] = *execCtx.AfterResultPath
		}

		// Map-specific members: the effective items extraction, the
		// batching, the concurrency setting and the tolerated-failure
		// thresholds.
		if _, isMap := state.(*sfnstore.MapState); isMap {
			if execCtx.AfterItemsPath != nil {
				data["afterItemsPath"] = *execCtx.AfterItemsPath
			}
			if execCtx.AfterItemBatcher != nil {
				data["afterItemBatcher"] = *execCtx.AfterItemBatcher
			}
			if execCtx.MaxConcurrencyValue != nil {
				data["maxConcurrency"] = *execCtx.MaxConcurrencyValue
			}
			if execCtx.ToleratedFailureCountVal != nil {
				data["toleratedFailureCount"] = *execCtx.ToleratedFailureCountVal
			}
			if execCtx.ToleratedFailurePctVal != nil {
				data["toleratedFailurePercentage"] = *execCtx.ToleratedFailurePctVal
			}
		}
	}

	return data
}
