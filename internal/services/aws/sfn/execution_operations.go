package sfn

import (
	"context"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
)

// StartExecution starts an execution of a state machine.
func (s *StepFunctionService) StartExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stateMachineArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	if stateMachineArn == "" {
		stateMachineArn = request.GetParamLowerFirst(req.Parameters, "StateMachineArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.startExecutionCore(ctx, store, StartExecutionInput{
		StateMachineArn: stateMachineArn,
		Name:            request.GetParamLowerFirst(req.Parameters, "name"),
		Input:           request.GetParamLowerFirst(req.Parameters, "input"),
		TraceHeader:     request.GetParamLowerFirst(req.Parameters, "traceHeader"),
	})
	if err != nil {
		return nil, err
	}

	// StartExecutionOutput carries the execution ARN and start date only.
	return map[string]interface{}{
		"executionArn": result.ExecutionArn,
		"startDate":    result.StartDate.Unix(),
	}, nil
}

// StopExecution stops a running execution of a state machine.
func (s *StepFunctionService) StopExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.stopExecutionCore(ctx, store, StopExecutionInput{
		ExecutionArn: request.GetParamLowerFirst(req.Parameters, "executionArn"),
		Error:        request.GetParamLowerFirst(req.Parameters, "error"),
		Cause:        request.GetParamLowerFirst(req.Parameters, "cause"),
	})
}

// DescribeExecution returns the details of an execution.
func (s *StepFunctionService) DescribeExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeExecutionCore(ctx, store, DescribeExecutionInput{
		ExecutionArn: request.GetParamLowerFirst(req.Parameters, "executionArn"),
		IncludedData: request.GetParamLowerFirst(req.Parameters, "includedData"),
	})
}

// ListExecutions returns a list of executions for a state machine.
func (s *StepFunctionService) ListExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	mapRunFilter := request.GetParamLowerFirst(req.Parameters, "mapRunArn")
	result, err := s.listExecutionsCore(ctx, store, ListExecutionsInput{
		StateMachineArn: request.GetParamLowerFirst(req.Parameters, "stateMachineArn"),
		StatusFilter:    request.GetParamLowerFirst(req.Parameters, "statusFilter"),
		MapRunArn:       mapRunFilter,
		RedriveFilter:   request.GetParamLowerFirst(req.Parameters, "redriveFilter"),
		MaxResults:      limit,
		NextToken:       request.GetParamLowerFirst(req.Parameters, "nextToken"),
	})
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
		// The child-execution linkage members are only returned when the
		// list is scoped by mapRunArn (ListExecutions member contract).
		if mapRunFilter != "" {
			if exec.MapRunArn != "" {
				executions[i]["mapRunArn"] = exec.MapRunArn
			}
			if exec.ItemCount != 0 {
				executions[i]["itemCount"] = exec.ItemCount
			}
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
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getExecutionHistoryCore(ctx, store, GetExecutionHistoryInput{
		ExecutionArn:         request.GetParamLowerFirst(req.Parameters, "executionArn"),
		MaxResults:           limit,
		NextToken:            request.GetParamLowerFirst(req.Parameters, "nextToken"),
		IncludeExecutionData: includeExecutionData,
		ReverseOrder:         reverseOrder,
	})
}

// RedriveExecution restarts an unsuccessful execution from the failed
// state; the eligibility contract and resume logic live in the execution
// Core.
func (s *StepFunctionService) RedriveExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.redriveExecutionCore(ctx, store, RedriveExecutionInput{
		ExecutionArn: request.GetParamLowerFirst(req.Parameters, "executionArn"),
		ClientToken:  request.GetParamLowerFirst(req.Parameters, "clientToken"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"redriveDate": result.RedriveDate.Unix(),
	}, nil
}

// generateExecutionName returns the name used when StartExecution omits
// one; AWS generates a UUID in that case.
func generateExecutionName() string {
	return uuid.New().String()
}

// StartSyncExecution starts a synchronous Express state machine
// execution; the validation, EXPRESS-only enforcement and synchronous run
// live in the execution Core.
func (s *StepFunctionService) StartSyncExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stateMachineArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	if stateMachineArn == "" {
		stateMachineArn = request.GetParamLowerFirst(req.Parameters, "StateMachineArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.startSyncExecutionCore(ctx, store, StartSyncExecutionInput{
		StateMachineArn: stateMachineArn,
		Name:            request.GetParamLowerFirst(req.Parameters, "name"),
		Input:           request.GetParamLowerFirst(req.Parameters, "input"),
		TraceHeader:     request.GetParamLowerFirst(req.Parameters, "traceHeader"),
		IncludedData:    request.GetParamLowerFirst(req.Parameters, "includedData"),
	})
}
