package sfn

import (
	"context"
	"errors"
	"fmt"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// This file holds the map-run-operation Core path shared by the HTTP API
// and the admin console handler, plus the map run ARN generators the
// execution engine uses.

// generateMapRunArn renders the ARN of a map run spawned by the given
// execution's state machine.
func generateMapRunArn(store *sfnstore.StepFunctionStore, region, accountID, executionArn, stateName string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("states", fmt.Sprintf("mapRun:%s/%s/%s",
		svcarn.ExtractStateMachineNameFromARN(executionArn), stateName, generateMapRunID(store)))
}

// generateMapRunID renders the unique identifier portion of a map run ARN.
func generateMapRunID(store *sfnstore.StepFunctionStore) string {
	id := store.NextMapRunSeq()
	return fmt.Sprintf("mapRun-%d-%s", id, time.Now().UTC().Format("20060102150405"))
}

// UpdateMapRunInput carries every field that UpdateMapRun needs; the
// pointer members distinguish "not provided" from an explicit zero.
type UpdateMapRunInput struct {
	MapRunArn                  string
	MaxConcurrency             *int64
	ToleratedFailureCount      *int64
	ToleratedFailurePercentage *float32
}

// describeMapRunCore is the single entry point for DescribeMapRun.
func (s *StepFunctionService) describeMapRunCore(ctx context.Context, store *sfnstore.StepFunctionStore, mapRunArn string) (map[string]interface{}, error) {
	if err := validateArnRequired(mapRunArn, "mapRunArn"); err != nil {
		return nil, err
	}

	mr, err := store.GetMapRun(ctx, mapRunArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrMapRunNotFound) {
			return nil, NewResourceNotFound("Map Run does not exist: " + mapRunArn)
		}
		return nil, err
	}
	return describeMapRunToResponse(mr), nil
}

// listMapRunsCore is the single entry point for ListMapRuns; listing the
// map runs of an unknown execution fails with ExecutionDoesNotExist
// rather than an empty list.
func (s *StepFunctionService) listMapRunsCore(ctx context.Context, store *sfnstore.StepFunctionStore, executionArn string, maxResults int32, nextToken string) (map[string]interface{}, error) {
	if err := validateArnRequired(executionArn, "executionArn"); err != nil {
		return nil, err
	}
	if err := validateMaxResults(maxResults, 0, sfnstore.MaxPageSize, "maxResults"); err != nil {
		return nil, err
	}

	if _, err := store.GetExecution(ctx, executionArn); err != nil {
		return nil, NewExecutionDoesNotExist("Execution Does not exist: " + executionArn)
	}

	result, err := store.ListAllMapRuns(ctx, executionArn, maxResults, nextToken)
	if err != nil {
		return nil, err
	}

	mapRuns := make([]map[string]interface{}, 0, len(result.MapRuns))
	for _, mr := range result.MapRuns {
		mapRuns = append(mapRuns, mapRunListItemToResponse(mr))
	}

	resp := map[string]interface{}{"mapRuns": mapRuns}
	if result.NextToken != "" && len(result.MapRuns) > 0 {
		resp["nextToken"] = result.NextToken
	}
	return resp, nil
}

// updateMapRunCore is the single entry point for UpdateMapRun; the
// numeric overrides obey their Smithy ranges.
func (s *StepFunctionService) updateMapRunCore(ctx context.Context, store *sfnstore.StepFunctionStore, in UpdateMapRunInput) error {
	if err := validateArnRequired(in.MapRunArn, "mapRunArn"); err != nil {
		return err
	}

	mr, err := store.GetMapRun(ctx, in.MapRunArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrMapRunNotFound) {
			return NewResourceNotFound("Map Run does not exist: " + in.MapRunArn)
		}
		return err
	}

	maxConcurrency := mr.MaxConcurrency
	toleratedFailureCount := mr.ToleratedFailureCount
	toleratedFailurePercentage := mr.ToleratedFailurePercentage
	if in.MaxConcurrency != nil {
		maxConcurrency = *in.MaxConcurrency
	}
	if in.ToleratedFailureCount != nil {
		toleratedFailureCount = *in.ToleratedFailureCount
	}
	if in.ToleratedFailurePercentage != nil {
		toleratedFailurePercentage = *in.ToleratedFailurePercentage
	}

	if err := validateMapRunUpdateParams(maxConcurrency, toleratedFailureCount, toleratedFailurePercentage); err != nil {
		return err
	}

	mr.MaxConcurrency = maxConcurrency
	mr.ToleratedFailureCount = toleratedFailureCount
	mr.ToleratedFailurePercentage = toleratedFailurePercentage

	if err := store.UpdateMapRun(ctx, mr); err != nil {
		return err
	}
	return nil
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
