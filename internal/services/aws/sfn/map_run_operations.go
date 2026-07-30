package sfn

import (
	"context"
	"errors"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
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

	if name == "" {
		name = generateExecutionName()
	}

	executionArn := arnutil.NewARNBuilder(s.accountID, reqCtx.GetRegion()).StepFunctions().Execution(arnutil.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, input, "")
	exec.ExecutionArn = executionArn

	if err := store.CreateExecution(ctx, exec); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil, awserrors.NewAWSError("ExecutionAlreadyExists", "An execution with the same name already exists: "+executionArn, 400)
		}
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

	return describeMapRunToResponse(mr), nil
}

// describeMapRunToResponse converts a MapRun record to the DescribeMapRun
// response shape. Matches Smithy DescribeMapRunOutput.
func describeMapRunToResponse(mr *sfnstore.MapRun) map[string]interface{} {
	itemCounts := map[string]interface{}{
		"pending":        mr.ItemCounts.Pending,
		"running":        mr.ItemCounts.Running,
		"succeeded":      mr.ItemCounts.Succeeded,
		"failed":         mr.ItemCounts.Failed,
		"timedOut":       mr.ItemCounts.TimedOut,
		"aborted":        mr.ItemCounts.Aborted,
		"total":          mr.ItemCounts.Total,
		"resultsWritten": mr.ItemCounts.ResultsWritten,
	}
	if mr.ItemCounts.FailuresNotRedrivable != 0 {
		itemCounts["failuresNotRedrivable"] = mr.ItemCounts.FailuresNotRedrivable
	}
	if mr.ItemCounts.PendingRedrive != 0 {
		itemCounts["pendingRedrive"] = mr.ItemCounts.PendingRedrive
	}

	executionCounts := map[string]interface{}{
		"pending":        mr.ExecutionCounts.Pending,
		"running":        mr.ExecutionCounts.Running,
		"succeeded":      mr.ExecutionCounts.Succeeded,
		"failed":         mr.ExecutionCounts.Failed,
		"timedOut":       mr.ExecutionCounts.TimedOut,
		"aborted":        mr.ExecutionCounts.Aborted,
		"total":          mr.ExecutionCounts.Total,
		"resultsWritten": mr.ExecutionCounts.ResultsWritten,
	}
	if mr.ExecutionCounts.FailuresNotRedrivable != 0 {
		executionCounts["failuresNotRedrivable"] = mr.ExecutionCounts.FailuresNotRedrivable
	}
	if mr.ExecutionCounts.PendingRedrive != 0 {
		executionCounts["pendingRedrive"] = mr.ExecutionCounts.PendingRedrive
	}

	resp := map[string]interface{}{
		"mapRunArn":                  mr.MapRunArn,
		"executionArn":               mr.ExecutionArn,
		"status":                     mr.Status,
		"startDate":                  mr.StartDate,
		"maxConcurrency":             mr.MaxConcurrency,
		"itemCounts":                 itemCounts,
		"executionCounts":            executionCounts,
		"toleratedFailurePercentage": mr.ToleratedFailurePercentage,
		"toleratedFailureCount":      mr.ToleratedFailureCount,
		"redriveCount":               mr.RedriveCount,
	}
	if mr.StopDate != 0 {
		resp["stopDate"] = mr.StopDate
	}
	if mr.RedriveDate != 0 {
		resp["redriveDate"] = mr.RedriveDate
	}

	return resp
}

// mapRunListItemToResponse converts a MapRun record to the ListMapRuns
// response shape. Matches Smithy MapRunListItem (5 fields only):
// executionArn, mapRunArn, stateMachineArn, startDate, stopDate.
// Item/execution counts and configuration live only in DescribeMapRunOutput.
func mapRunListItemToResponse(mr *sfnstore.MapRun) map[string]interface{} {
	resp := map[string]interface{}{
		"executionArn":    mr.ExecutionArn,
		"mapRunArn":       mr.MapRunArn,
		"stateMachineArn": mr.StateMachineArn,
		"startDate":       mr.StartDate,
	}
	if mr.StopDate != 0 {
		resp["stopDate"] = mr.StopDate
	}
	return resp
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

	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}

	result, err := store.ListAllMapRuns(ctx, executionArn, limit, nextToken)
	if err != nil {
		return nil, err
	}

	mapRuns := make([]map[string]interface{}, 0, len(result.MapRuns))
	for _, mr := range result.MapRuns {
		mapRuns = append(mapRuns, mapRunListItemToResponse(mr))
	}

	respNextToken := ""
	if result.NextToken != "" && len(result.MapRuns) > 0 {
		respNextToken = result.MapRuns[len(result.MapRuns)-1].MapRunArn
	}

	resp := map[string]interface{}{
		"mapRuns": mapRuns,
	}
	if respNextToken != "" {
		resp["nextToken"] = respNextToken
	}

	return resp, nil
}
