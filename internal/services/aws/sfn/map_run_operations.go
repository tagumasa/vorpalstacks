package sfn

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arncommon "vorpalstacks/internal/utils/aws/arn"
)

var (
	mapRunsMu sync.Mutex
	mapRuns   = map[string][]map[string]interface{}{}
	mapRunSeq int64
)

func generateMapRunArn(region, accountID, executionArn, stateName string) string {
	return fmt.Sprintf("arn:aws:states:%s:%s:mapRun:%s/%s/%s",
		region, accountID, arncommon.ExtractStateMachineNameFromARN(executionArn), stateName, generateMapRunID())
}

func generateMapRunID() string {
	mapRunsMu.Lock()
	mapRunSeq++
	id := mapRunSeq
	mapRunsMu.Unlock()
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

	executionArn := fmt.Sprintf("arn:aws:states:%s:%s:execution:%s:%s",
		reqCtx.GetRegion(), s.accountID,
		arncommon.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, input, "")
	exec.ExecutionArn = executionArn

	if err := store.CreateExecution(ctx, exec); err != nil {
		return nil, err
	}

	sqsStore := s.sqsStore
	if sqsStore == nil {
		sqsStore = reqCtx.GetSQSStore()
	}
	snsStore := s.snsStore
	if snsStore == nil {
		snsStore = reqCtx.GetSNSStore()
	}
	eventsStore := s.eventsStore
	if eventsStore == nil {
		if ebStore := reqCtx.GetEventBridgeStore(); ebStore != nil {
			if concrete, ok := ebStore.(*eventsstore.EventsStore); ok {
				eventsStore = concrete
			}
		}
	}

	executor := NewExecutorWithStores(store, s.lambdaInvoker, sqsStore, snsStore, eventsStore, s.accountID, reqCtx.GetRegion())
	executor.SetEventBus(s.bus)
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

	mapRunsMu.Lock()
	defer mapRunsMu.Unlock()

	for _, runs := range mapRuns {
		for _, run := range runs {
			if arn, ok := run["mapRunArn"].(string); ok && arn == mapRunArn {
				return run, nil
			}
		}
	}

	return nil, NewMapRunDoesNotExist("Map Run does not exist: " + mapRunArn)
}

// ListMapRuns lists map runs, optionally filtered by execution ARN.
func (s *StepFunctionService) ListMapRuns(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	executionArn := request.GetParamLowerFirst(req.Parameters, "executionArn")

	mapRunsMu.Lock()
	defer mapRunsMu.Unlock()

	var matchingRuns []map[string]interface{}
	for execArn, runs := range mapRuns {
		if executionArn != "" && execArn != executionArn {
			continue
		}
		matchingRuns = append(matchingRuns, runs...)
	}

	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit <= 0 {
		limit = 100
	}

	if int32(len(matchingRuns)) > limit {
		matchingRuns = matchingRuns[:limit]
	}

	return map[string]interface{}{
		"mapRuns": matchingRuns,
	}, nil
}
