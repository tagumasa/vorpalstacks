package testutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

// runSFNMapFeatureTests pins the Distributed Map ItemBatcher data plane,
// the child-execution dispatch surfaced through ListExecutions with a
// mapRunArn filter, the closed-world state-field rejection and the
// direct-member InputPath/OutputPath/Fail-path contracts.
func (r *TestRunner) runSFNMapFeatureTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	// latestMapRunFor finds the Map Run the given parent execution
	// spawned, traversing every page of the listing.
	latestMapRunFor := func(executionArn string) (string, error) {
		var nextToken *string
		var best string
		var bestStart time.Time
		for {
			page, err := tc.client.ListMapRuns(tc.ctx, &awssfn.ListMapRunsInput{
				ExecutionArn: aws.String(executionArn),
				NextToken:    nextToken,
			})
			if err != nil {
				return "", err
			}
			for _, mr := range page.MapRuns {
				if aws.ToString(mr.ExecutionArn) == executionArn && mr.StartDate.After(bestStart) {
					best = aws.ToString(mr.MapRunArn)
					bestStart = *mr.StartDate
				}
			}
			if page.NextToken == nil {
				break
			}
			nextToken = page.NextToken
		}
		if best == "" {
			return "", fmt.Errorf("no map run found for %s", executionArn)
		}
		return best, nil
	}

	// The batched Distributed Map used by the batching and child-execution
	// assertions: five items in batches of two, with a fixed BatchInput.
	batchState := map[string]interface{}{
		"Type":      "Map",
		"Label":     "dist",
		"ItemsPath": "$.v",
		"ItemBatcher": map[string]interface{}{
			"MaxItemsPerBatch": 2,
			"BatchInput":       map[string]interface{}{"factCheck": "December 2022"},
		},
		"ItemProcessor": map[string]interface{}{
			"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
			"StartAt":         "W",
			"States": map[string]interface{}{
				"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
			},
		},
		"End": true,
	}
	batchSM, err := tc.createSingleStateSM("MapBatch-"+ts, batchState)
	if err != nil {
		return []TestResult{{Service: "stepfunctions", TestName: "MapFeatureSetup", Status: "FAIL", Error: fmt.Sprintf("create SM: %v", err)}}
	}
	defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(batchSM)})

	parentName := "feat-parent-" + ts
	parentExecArn := ""
	results = append(results, r.RunTest("stepfunctions", "ItemBatcher_BatchesAndCounts", func() error {
		execArn, output, oerr := tc.runWithInput(batchSM, parentName, `{"v":[1,2,3,4,5]}`)
		if oerr != nil {
			return oerr
		}
		parentExecArn = execArn
		var units []map[string]interface{}
		if err := json.Unmarshal([]byte(output), &units); err != nil {
			return fmt.Errorf("output not a batch array: %v (%s)", err, output)
		}
		if len(units) != 3 {
			return fmt.Errorf("expected 3 batch outputs, got %d (%s)", len(units), output)
		}
		for i, want := range []int{2, 2, 1} {
			items, ok := units[i]["Items"].([]interface{})
			if !ok || len(items) != want {
				return fmt.Errorf("unit %d Items = %v, want %d items", i, units[i]["Items"], want)
			}
			bi, ok := units[i]["BatchInput"].(map[string]interface{})
			if !ok || bi["factCheck"] != "December 2022" {
				return fmt.Errorf("unit %d BatchInput = %v", i, units[i]["BatchInput"])
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ItemBatcher_MapRunItemAndExecutionCounts", func() error {
		if parentExecArn == "" {
			return fmt.Errorf("parent execution ARN unavailable")
		}
		mapRunArn, lerr := latestMapRunFor(parentExecArn)
		if lerr != nil {
			return lerr
		}
		desc, err := tc.client.DescribeMapRun(tc.ctx, &awssfn.DescribeMapRunInput{MapRunArn: aws.String(mapRunArn)})
		if err != nil {
			return err
		}
		if desc.ItemCounts.Total != 5 || desc.ItemCounts.Succeeded != 5 {
			return fmt.Errorf("item counts = %+v, want 5/5", desc.ItemCounts)
		}
		if desc.ExecutionCounts.Total != 3 || desc.ExecutionCounts.Succeeded != 3 {
			return fmt.Errorf("execution counts = %+v, want 3/3", desc.ExecutionCounts)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DistributedMap_ChildExecutionsListedByMapRun", func() error {
		if parentExecArn == "" {
			return fmt.Errorf("parent execution ARN unavailable")
		}
		mapRunArn, merr := latestMapRunFor(parentExecArn)
		if merr != nil {
			return merr
		}
		var children []types.ExecutionListItem
		var nextToken *string
		for {
			page, err := tc.client.ListExecutions(tc.ctx, &awssfn.ListExecutionsInput{
				MapRunArn: aws.String(mapRunArn),
				NextToken: nextToken,
			})
			if err != nil {
				return err
			}
			children = append(children, page.Executions...)
			if page.NextToken == nil {
				break
			}
			nextToken = page.NextToken
		}
		if len(children) != 3 {
			return fmt.Errorf("expected 3 child executions, got %d", len(children))
		}
		// The listing is sorted by start time, so the ordinal order of
		// near-simultaneous children is not guaranteed; assert the set.
		byName := map[string]types.ExecutionListItem{}
		for _, child := range children {
			byName[aws.ToString(child.Name)] = child
			if aws.ToString(child.MapRunArn) != mapRunArn {
				return fmt.Errorf("child %q mapRunArn = %q", aws.ToString(child.Name), aws.ToString(child.MapRunArn))
			}
			if child.Status != types.ExecutionStatusSucceeded {
				return fmt.Errorf("child %q status = %s", aws.ToString(child.Name), child.Status)
			}
		}
		for ordinal, wantItems := range map[int]int32{0: 2, 1: 2, 2: 1} {
			name := fmt.Sprintf("%s:dist-%d", parentName, ordinal)
			child, ok := byName[name]
			if !ok {
				return fmt.Errorf("child %q missing from the map run listing (have %v)", name, byNameKeys(byName))
			}
			if aws.ToInt32(child.ItemCount) != wantItems {
				return fmt.Errorf("child %q itemCount = %d, want %d", name, aws.ToInt32(child.ItemCount), wantItems)
			}
		}
		// A child answers DescribeExecution and GetExecutionHistory with
		// its own record.
		desc, err := tc.client.DescribeExecution(tc.ctx, &awssfn.DescribeExecutionInput{
			ExecutionArn: children[0].ExecutionArn,
		})
		if err != nil {
			return fmt.Errorf("describe child: %v", err)
		}
		if desc.Status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("child describe status = %s", desc.Status)
		}
		history, err := tc.client.GetExecutionHistory(tc.ctx, &awssfn.GetExecutionHistoryInput{
			ExecutionArn: children[0].ExecutionArn,
		})
		if err != nil {
			return fmt.Errorf("child history: %v", err)
		}
		if len(history.Events) == 0 {
			return fmt.Errorf("child history is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "PassState_DirectInputOutputPath", func() error {
		smARN, err := tc.createSingleStateSM("MapIOPath-"+ts, map[string]interface{}{
			"Type":       "Pass",
			"InputPath":  "$.keep",
			"OutputPath": "$.n",
			"End":        true,
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, output, rerr := tc.runWithInput(smARN, "io-"+ts, `{"keep":{"n":42},"drop":1}`)
		if rerr != nil {
			return rerr
		}
		if output != "42" {
			return fmt.Errorf("output = %s, want 42", output)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "FailState_ErrorCausePathResolved", func() error {
		smARN, err := tc.createSingleStateSM("MapFailPath-"+ts, map[string]interface{}{
			"Type":      "Fail",
			"ErrorPath": "$.code",
			"CausePath": "$.reason",
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		execArn, serr := tc.startExecution(smARN, "failpath-"+ts, `{"code":"CustomError","reason":"resolved cause"}`)
		if serr != nil {
			return serr
		}
		desc, werr := tc.awaitTerminal(execArn, 200*time.Millisecond, 30)
		if werr != nil {
			return werr
		}
		if desc.Status != types.ExecutionStatusFailed {
			return fmt.Errorf("status = %s, want FAILED", desc.Status)
		}
		if aws.ToString(desc.Error) != "CustomError" {
			return fmt.Errorf("error = %q, want CustomError", aws.ToString(desc.Error))
		}
		if aws.ToString(desc.Cause) != "resolved cause" {
			return fmt.Errorf("cause = %q, want 'resolved cause'", aws.ToString(desc.Cause))
		}
		return nil
	}))

	return results
}

// byNameKeys renders the child names of a name-keyed map for failure
// messages.
func byNameKeys(m map[string]types.ExecutionListItem) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
