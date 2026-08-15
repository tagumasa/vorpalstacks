package sfn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

func (s *StepFunctionService) StartExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stateMachineArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	name := request.GetParamLowerFirst(req.Parameters, "name")
	input := request.GetParamLowerFirst(req.Parameters, "input")
	traceHeader := request.GetParamLowerFirst(req.Parameters, "traceHeader")

	if stateMachineArn == "" {
		stateMachineArn = request.GetParamLowerFirst(req.Parameters, "StateMachineArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sm, err := store.GetStateMachine(ctx, stateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + stateMachineArn)
		}
		return nil, err
	}

	if name == "" {
		name = generateExecutionName()
	} else if len(name) > sfnstore.MaxResourceNameLength {
		return nil, NewInvalidName(fmt.Sprintf("name must be 1-80 characters, got %d", len(name)))
	}

	executionArn := arnutil.NewARNBuilder(s.accountID, reqCtx.GetRegion()).StepFunctions().Execution(arnutil.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, input, traceHeader)
	exec.ExecutionArn = executionArn

	if err := store.CreateExecution(ctx, exec); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil, NewExecutionAlreadyExists("An execution with the same name already exists: " + executionArn)
		}
		return nil, err
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, reqCtx.GetRegion())
	execCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(executionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(executionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in execution", logs.String("arn", executionArn), logs.Any("panic", r))
				exec.Status = "FAILED"
				exec.Error = "States.InternalError"
				exec.Cause = fmt.Sprintf("internal panic: %v", r)
				exec.StopDate = time.Now().UTC()
				_ = store.UpdateExecution(context.Background(), exec)
			}
		}()
		if err := executor.ExecuteStateMachine(execCtx, exec); err != nil {
			logs.Error("sfn: execution error", logs.String("arn", executionArn), logs.Err(err))
		}
	}()

	return map[string]interface{}{
		"executionArn":    exec.ExecutionArn,
		"startDate":       exec.StartDate.Unix(),
		"stateMachineArn": exec.StateMachineArn,
	}, nil
}

// StopExecution stops a running execution of a state machine.
func (s *StepFunctionService) StopExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "executionArn")
	if err := validateArnRequired(arn, "executionArn"); err != nil {
		return nil, err
	}
	errorMsg := request.GetParamLowerFirst(req.Parameters, "error")
	cause := request.GetParamLowerFirst(req.Parameters, "cause")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exec, err := store.GetExecution(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + arn)
		}
		return nil, err
	}

	if isTerminalStatus(exec.Status) {
		return map[string]interface{}{
			"stopDate": exec.StopDate.Unix(),
		}, nil
	}

	store.CancelExecution(arn)

	exec.Status = "ABORTED"
	exec.StopDate = time.Now().UTC()
	exec.Error = errorMsg
	exec.Cause = cause

	if err := store.UpdateExecution(ctx, exec); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stopDate": exec.StopDate.Unix(),
	}, nil
}

// DescribeExecution returns the details of an execution.
func (s *StepFunctionService) DescribeExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "executionArn")
	if err := validateArnRequired(arn, "executionArn"); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exec, err := store.GetExecution(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + arn)
		}
		return nil, err
	}

	return executionToResponse(exec), nil
}

// ListExecutions returns a list of executions for a state machine.
func (s *StepFunctionService) ListExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stateMachineArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	statusFilter := request.GetParamLowerFirst(req.Parameters, "statusFilter")
	mapRunArn := request.GetParamLowerFirst(req.Parameters, "mapRunArn")
	redriveFilter := request.GetParamLowerFirst(req.Parameters, "redriveFilter")
	if redriveFilter != "" && redriveFilter != "REDRIVEN" && redriveFilter != "NOT_REDRIVEN" {
		return nil, NewValidationException("redriveFilter must be REDRIVEN or NOT_REDRIVEN, got " + redriveFilter)
	}
	if !isValidExecutionStatus(statusFilter) {
		return nil, NewValidationException("statusFilter must be one of RUNNING, SUCCEEDED, FAILED, TIMED_OUT, ABORTED, PENDING_REDRIVE, got " + statusFilter)
	}
	// PENDING_REDRIVE lists child workflow executions awaiting redrive;
	// those only exist in the scope of a Map Run, so the documented
	// contract requires mapRunArn and rejects a stateMachineArn pairing
	// with a validation exception — unconditionally, even when a
	// mapRunArn is also present.
	if statusFilter == "PENDING_REDRIVE" && stateMachineArn != "" {
		return nil, NewValidationException("statusFilter PENDING_REDRIVE requires mapRunArn; providing stateMachineArn with PENDING_REDRIVE is not supported")
	}
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.GetStateMachine(ctx, stateMachineArn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + stateMachineArn)
	}
	result, err := store.ListExecutions(ctx, stateMachineArn, statusFilter, mapRunArn, redriveFilter, limit, nextToken)
	if err != nil {
		return nil, err
	}

	executions := make([]map[string]interface{}, len(result.Executions))
	for i, exec := range result.Executions {
		executions[i] = map[string]interface{}{
			"executionArn":    exec.ExecutionArn,
			"stateMachineArn": exec.StateMachineArn,
			"name":            exec.Name,
			"status":          exec.Status,
			"startDate":       exec.StartDate.Unix(),
		}
		if !exec.StopDate.IsZero() {
			executions[i]["stopDate"] = exec.StopDate.Unix()
		}
		if exec.MapRunArn != "" {
			executions[i]["mapRunArn"] = exec.MapRunArn
		}
		if exec.ItemCount != 0 {
			executions[i]["itemCount"] = exec.ItemCount
		}
		if exec.RedriveCount != 0 {
			executions[i]["redriveCount"] = exec.RedriveCount
		}
		if !exec.RedriveDate.IsZero() {
			executions[i]["redriveDate"] = exec.RedriveDate.Unix()
		}
		if exec.StateMachineAliasArn != "" {
			executions[i]["stateMachineAliasArn"] = exec.StateMachineAliasArn
		}
		if exec.StateMachineVersionArn != "" {
			executions[i]["stateMachineVersionArn"] = exec.StateMachineVersionArn
		}
	}

	response := map[string]interface{}{
		"executions": executions,
	}

	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}

	return response, nil
}

// GetExecutionHistory returns the history of an execution.
func (s *StepFunctionService) GetExecutionHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "executionArn")
	if err := validateArnRequired(arn, "executionArn"); err != nil {
		return nil, err
	}
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.GetExecution(ctx, arn)
	if err != nil {
		return nil, NewExecutionDoesNotExist("Execution Does not exist: " + arn)
	}
	includeExecutionData := true
	if v, ok := req.Parameters["includeExecutionData"]; ok {
		if vBool, ok := v.(bool); ok {
			includeExecutionData = vBool
		}
	}

	reverseOrder := false
	if v, ok := req.Parameters["reverseOrder"]; ok {
		if vBool, ok := v.(bool); ok {
			reverseOrder = vBool
		}
	}

	// Reverse order must paginate in reverse as a whole: the store serves
	// newest-first pages with a direction-consistent marker, so reversing
	// an ascending page in place would scramble the global order.
	events, nextTokenResult, err := store.GetExecutionHistory(ctx, arn, limit, nextToken, reverseOrder)
	if err != nil {
		if errors.Is(err, sfnstore.ErrInvalidToken) {
			return nil, NewInvalidToken("Invalid nextToken: " + nextToken)
		}
		return nil, err
	}

	history := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		history = append(history, historyEventToResponse(event, includeExecutionData))
	}

	response := map[string]interface{}{
		"events": history,
	}

	if nextTokenResult != "" {
		response["nextToken"] = nextTokenResult
	}

	return response, nil
}

// generateExecutionName returns the name used when StartExecution omits
// one; AWS generates a UUID in that case.
func generateExecutionName() string {
	return uuid.New().String()
}
