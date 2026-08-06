package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// RedriveExecution restarts an unsuccessful execution from the failed state.
// The execution keeps its original ARN, name, input, and history; new events
// are appended. Previously-succeeded states are not re-executed.
// https://docs.aws.amazon.com/step-functions/latest/apireference/API_RedriveExecution.html
func (s *StepFunctionService) RedriveExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	executionArn := request.GetParamLowerFirst(req.Parameters, "executionArn")
	if err := validateArnRequired(executionArn, "executionArn"); err != nil {
		return nil, err
	}

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

	if !isRedrivableStatus(exec.Status) {
		return nil, NewInvalidExecutionType(
			fmt.Sprintf("Execution %s is in %s status and cannot be redriven", executionArn, exec.Status))
	}

	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + exec.StateMachineArn)
	}

	definition, err := parseStateMachineDefinition(sm.Definition)
	if err != nil {
		return nil, NewInvalidDefinitionException("Invalid state machine definition: " + err.Error())
	}

	rp, err := determineResumePoint(ctx, store, executionArn, definition)
	if err != nil {
		return nil, fmt.Errorf("failed to determine resume point: %w", err)
	}

	redriveDate := time.Now().UTC()
	exec.Status = "RUNNING"
	exec.Error = ""
	exec.Cause = ""
	exec.StopDate = time.Time{}
	exec.RedriveCount++
	exec.RedriveDate = redriveDate
	exec.Output = ""

	if err := store.UpdateExecution(ctx, exec); err != nil {
		return nil, fmt.Errorf("failed to update execution for redrive: %w", err)
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, reqCtx.GetRegion())
	resumeCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(executionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(executionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in redrive execution", logs.String("arn", executionArn), logs.Any("panic", r))
				exec.Status = "FAILED"
				exec.Error = "States.InternalError"
				exec.Cause = fmt.Sprintf("internal panic: %v", r)
				exec.StopDate = time.Now().UTC()
				_ = store.UpdateExecution(context.Background(), exec)
			}
		}()
		if err := executor.ExecuteStateMachineFromState(resumeCtx, exec, rp.StateName, rp.Input, rp.LastEventId); err != nil {
			logs.Error("sfn: redrive execution failed", logs.String("arn", executionArn), logs.Err(err))
		}
	}()

	return map[string]interface{}{
		"redriveDate": redriveDate.Unix(),
	}, nil
}

// TestState tests a single state within a state machine definition.
// https://docs.aws.amazon.com/step-functions/latest/apireference/API_TestState.html
//
// Optional AWS parameters (roleArn, mock, context, stateConfiguration,
// revealSecrets) are accepted but currently unused. The vorpalstacks
// test runner executes the state in-memory against synthetic inputs, so
// real IAM role assumption and service-call mocking are not exercised.
// The framework ignores unknown request fields, so the parameters are
// simply not read here.
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

	if inspectionLevel == "" {
		inspectionLevel = "INFO"
	}
	if inspectionLevel != "INFO" && inspectionLevel != "DEBUG" && inspectionLevel != "TRACE" {
		return nil, NewInvalidParameterValue("inspectionLevel must be INFO, DEBUG, or TRACE, got " + inspectionLevel)
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

func isErrNotFound(err error) bool {
	return errors.Is(err, sfnstore.ErrExecutionNotFound) || errors.Is(err, sfnstore.ErrStateMachineNotFound)
}
