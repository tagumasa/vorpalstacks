package sfn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

func generateMapRunArn(store *sfnstore.StepFunctionStore, region, accountID, executionArn, stateName string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("states", fmt.Sprintf("mapRun:%s/%s/%s",
		arnutil.ExtractStateMachineNameFromARN(executionArn), stateName, generateMapRunID(store)))
}

func generateMapRunID(store *sfnstore.StepFunctionStore) string {
	id := store.NextMapRunSeq()
	return fmt.Sprintf("mapRun-%d-%s", id, time.Now().UTC().Format("20060102150405"))
}

// StartSyncExecution starts a synchronous Step Functions execution.
func (s *StepFunctionService) StartSyncExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stateMachineArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	name := request.GetParamLowerFirst(req.Parameters, "name")
	input := request.GetParamLowerFirst(req.Parameters, "input")

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

	if sm.Type == "EXPRESS" {
		return nil, NewInvalidExecutionType("EXPRESS state machines do not support synchronous executions")
	}

	if name == "" {
		name = generateExecutionName()
	}

	executionArn := arnutil.NewARNBuilder(s.accountID, reqCtx.GetRegion()).StepFunctions().Execution(arnutil.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, input, "")
	exec.ExecutionArn = executionArn

	if err := store.CreateExecution(ctx, exec); err != nil {
		return nil, err
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, reqCtx.GetRegion())
	execErr := executor.ExecuteStateMachine(ctx, exec)

	updated, err := store.GetExecution(ctx, executionArn)
	if err != nil {
		updated = exec
	}

	result := map[string]interface{}{
		"executionArn":     updated.ExecutionArn,
		"stateMachineArn":  updated.StateMachineArn,
		"name":             updated.Name,
		"startDate":        updated.StartDate.Unix(),
		"stopDate":         updated.StopDate.Unix(),
		"status":           updated.Status,
		"stateMachineType": sm.Type,
	}

	if updated.Output != "" {
		result["output"] = updated.Output
	}
	if execErr != nil {
		result["error"] = updated.Error
		result["cause"] = updated.Cause
	}

	return result, nil
}

// DescribeMapRun retrieves the details of a Step Functions map run.
func (s *StepFunctionService) DescribeMapRun(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	mapRunArn := request.GetParamLowerFirst(req.Parameters, "mapRunArn")
	if mapRunArn == "" {
		return nil, NewInvalidName("mapRunArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	mr, err := store.GetMapRun(ctx, mapRunArn)
	if err != nil {
		return nil, NewMapRunDoesNotExist("Map Run does not exist: " + mapRunArn)
	}

	return map[string]interface{}{
		"mapRunArn":           mr.MapRunArn,
		"executionArn":        mr.ExecutionArn,
		"stateMachineArn":     mr.StateMachineArn,
		"name":                mr.Name,
		"status":              mr.Status,
		"startDate":           mr.StartDate,
		"stopDate":            mr.StopDate,
		"itemCount":           mr.ItemCount,
		"maxConcurrency":      mr.MaxConcurrency,
		"itemsProcessedCount": mr.ItemsProcessedCount,
		"itemsFailedCount":    mr.ItemsFailedCount,
		"itemsCancelledCount": mr.ItemsCancelledCount,
	}, nil
}

// ListMapRuns lists map runs, optionally filtered by execution ARN.
// Supports pagination via nextToken/maxResults, returning nextToken
// when more results are available than the requested page size.
func (s *StepFunctionService) ListMapRuns(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	executionArn := request.GetParamLowerFirst(req.Parameters, "executionArn")
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit <= 0 {
		limit = 100
	}

	result, err := store.ListAllMapRuns(ctx, executionArn, limit, nextToken)
	if err != nil {
		return nil, err
	}

	mapRuns := make([]map[string]interface{}, 0, len(result.MapRuns))
	for _, mr := range result.MapRuns {
		mapRuns = append(mapRuns, map[string]interface{}{
			"mapRunArn":           mr.MapRunArn,
			"executionArn":        mr.ExecutionArn,
			"stateMachineArn":     mr.StateMachineArn,
			"name":                mr.Name,
			"status":              mr.Status,
			"startDate":           mr.StartDate,
			"stopDate":            mr.StopDate,
			"itemCount":           mr.ItemCount,
			"maxConcurrency":      mr.MaxConcurrency,
			"itemsProcessedCount": mr.ItemsProcessedCount,
			"itemsFailedCount":    mr.ItemsFailedCount,
			"itemsCancelledCount": mr.ItemsCancelledCount,
		})
	}

	respNextToken := ""
	if result.NextToken != "" && len(result.MapRuns) > 0 {
		respNextToken = result.MapRuns[len(result.MapRuns)-1].MapRunArn
	}

	return map[string]interface{}{
		"mapRuns":   mapRuns,
		"nextToken": respNextToken,
	}, nil
}
