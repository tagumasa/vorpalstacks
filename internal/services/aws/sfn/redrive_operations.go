package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// RedriveExecution restarts an unsuccessful execution from the failed state.
// https://docs.aws.amazon.com/step-functions/latest/apireference/API_RedriveExecution.html
func (s *StepFunctionService) RedriveExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	executionArn := request.GetParamLowerFirst(req.Parameters, "executionArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	exec, err := store.GetExecution(ctx, executionArn)
	if err != nil {
		if isErrNotFound(err) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + executionArn)
		}
		return nil, err
	}

	redrivableStatuses := map[string]bool{
		"FAILED":    true,
		"TIMED_OUT": true,
		"ABORTED":   true,
	}
	if !redrivableStatuses[exec.Status] {
		return nil, NewInvalidExecutionType(
			fmt.Sprintf("Execution %s is in %s status and cannot be redriven", executionArn, exec.Status))
	}

	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + exec.StateMachineArn)
	}

	history, _, err := store.GetExecutionHistory(ctx, executionArn, 100, "")
	if err != nil {
		return nil, err
	}
	_ = history

	newExecutionArn := arnutil.NewARNBuilder(s.accountID, reqCtx.GetRegion()).StepFunctions().Execution(
		extractStateMachineName(sm.StateMachineArn), exec.Name+"-redrive-"+strconv.FormatInt(time.Now().UnixNano(), 10))

	newExec := sfnstore.NewExecution(exec.StateMachineArn, exec.Name+"-redrive", exec.Input, "")
	newExec.ExecutionArn = newExecutionArn
	newExec.TraceHeader = exec.TraceHeader

	if err := store.CreateExecution(ctx, newExec); err != nil {
		return nil, err
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, reqCtx.GetRegion())
	execCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(newExecutionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(newExecutionArn)
		if err := executor.ExecuteStateMachine(execCtx, newExec); err != nil {
			logs.Error("sfn: redrive execution failed", logs.String("arn", newExecutionArn), logs.Err(err))
		}
	}()

	redriveDate := time.Now().UTC()
	return map[string]interface{}{
		"executionArn":    newExec.ExecutionArn,
		"stateMachineArn": newExec.StateMachineArn,
		"name":            newExec.Name,
		"startDate":       newExec.StartDate.Unix(),
		"redriveDate":     redriveDate.Unix(),
	}, nil
}

// TestState tests a single state within a state machine definition.
// https://docs.aws.amazon.com/step-functions/latest/apireference/API_TestState.html
func (s *StepFunctionService) TestState(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	stateName := request.GetParamLowerFirst(req.Parameters, "stateName")
	input := request.GetParamLowerFirst(req.Parameters, "input")
	inspectionLevel := request.GetParamLowerFirst(req.Parameters, "inspectionLevel")
	variablesParam := request.GetParamLowerFirst(req.Parameters, "variables")

	if definition == "" {
		return nil, NewInvalidDefinitionException("definition is required")
	}
	if stateName == "" {
		return nil, NewInvalidDefinitionException("stateName is required")
	}

	if input == "" {
		input = "{}"
	}

	var def sfnstore.StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return nil, NewInvalidDefinitionException("definition is not valid JSON: " + err.Error())
	}

	if def.QueryLanguage == "" {
		def.QueryLanguage = "JSONPath"
	}

	rawState, exists := def.States[stateName]
	if !exists {
		return nil, NewInvalidDefinitionException(
			fmt.Sprintf("State '%s' not found in definition", stateName))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, reqCtx.GetRegion())
	executor.currentStateMachine = &sfnstore.StateMachine{
		StateMachineArn: "arn:aws:states:" + reqCtx.GetRegion() + ":" + s.accountID + ":stateMachine:test-sm",
		Name:            "test-sm",
	}

	state, err := executor.parseState(stateName, rawState)
	if err != nil {
		return nil, NewInvalidDefinitionException(err.Error())
	}

	variableScope := NewVariableScope(nil)
	if variablesParam != "" {
		var rawVars map[string]interface{}
		if err := json.Unmarshal([]byte(variablesParam), &rawVars); err != nil {
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

	testExec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:" + reqCtx.GetRegion() + ":" + s.accountID + ":execution:test-state:" + stateName,
		StateMachineArn: "arn:aws:states:" + reqCtx.GetRegion() + ":" + s.accountID + ":stateMachine:test-sm",
		Name:            "TestState-" + stateName,
		Status:          "RUNNING",
		Input:           input,
		StartDate:       time.Now().UTC(),
	}

	eventId := int64(1)
	execCtx := &ExecutionContext{
		Execution:        testExec,
		Definition:       &def,
		CurrentState:     stateName,
		Input:            input,
		Output:           "",
		EventId:          &eventId,
		States:           map[string]sfnstore.State{stateName: state},
		QueryLanguage:    def.QueryLanguage,
		VariableScope:    variableScope,
		StateEnteredTime: time.Now().UTC(),
		MapItemIndex:     -1,
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
			return errResult, nil
		}
	case *sfnstore.ParallelState:
		output, nextState, execErr = executor.executeParallel(ctx, execCtx, st)
		if execErr != nil {
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
			return errResult, nil
		}
	case *sfnstore.MapState:
		output, nextState, execErr = executor.executeMap(ctx, execCtx, st)
		if execErr != nil {
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
			return errResult, nil
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

func buildInspectionData(inspectionLevel string, execCtx *ExecutionContext, output string, state sfnstore.State) map[string]interface{} {
	data := map[string]interface{}{
		"input":  execCtx.Input,
		"output": output,
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
	}

	return data
}

func extractStateMachineName(arn string) string {
	for i := len(arn) - 1; i >= 0; i-- {
		if arn[i] == ':' {
			return arn[i+1:]
		}
	}
	return arn
}

func isErrNotFound(err error) bool {
	return err != nil && (err.Error() == "execution not found" ||
		err.Error() == "state machine not found")
}
