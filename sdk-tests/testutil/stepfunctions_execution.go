package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

func (r *TestRunner) runSFNExecutionTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	execSMName := fmt.Sprintf("ExecSM-%d", time.Now().UnixNano())
	_, execRoleARN, execRoleCleanup := tc.createRoleForSM("ExecRole")
	defer execRoleCleanup()

	passDef := `{"Comment":"pass","StartAt":"B","States":{"B":{"Type":"Pass","Result":"hello","End":true}}}`
	var execSMARN string
	results = append(results, r.RunTest("stepfunctions", "StartExecution", func() error {
		resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(execSMName),
			Definition: aws.String(passDef),
			RoleArn:    aws.String(execRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create SM: %v", err)
		}
		execSMARN = *resp.StateMachineArn

		input := map[string]string{"message": "test"}
		inputJSON, _ := json.Marshal(input)
		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(execSMARN),
			Input:           aws.String(string(inputJSON)),
		})
		if err != nil {
			return err
		}
		if startResp.ExecutionArn == nil {
			return fmt.Errorf("execution ARN is nil")
		}
		if startResp.StartDate.IsZero() {
			return fmt.Errorf("start date is zero")
		}
		return nil
	}))

	var executionARN string
	results = append(results, r.RunTest("stepfunctions", "ListExecutions", func() error {
		resp, err := tc.client.ListExecutions(tc.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(execSMARN),
		})
		if err != nil {
			return err
		}
		if resp.Executions == nil {
			return fmt.Errorf("executions list is nil")
		}
		for _, ex := range resp.Executions {
			if ex.ExecutionArn != nil {
				executionARN = *ex.ExecutionArn
				break
			}
		}
		if executionARN == "" {
			return fmt.Errorf("no execution found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeExecution", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		resp, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return err
		}
		if resp.Status == "" {
			return fmt.Errorf("execution status is empty")
		}
		if resp.StateMachineArn == nil || *resp.StateMachineArn != execSMARN {
			return fmt.Errorf("state machine ARN mismatch in execution")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "GetExecutionHistory", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		resp, err := tc.client.GetExecutionHistory(tc.ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return err
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("events list is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachineForExecution", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		descResp, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return fmt.Errorf("describe execution: %v", err)
		}
		if descResp.StateMachineArn == nil {
			return fmt.Errorf("execution has no state machine ARN")
		}
		if descResp.Status != types.ExecutionStatusSucceeded && descResp.Status != types.ExecutionStatusRunning {
			return fmt.Errorf("execution not in suitable state: %s", descResp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Execution_PassStateOutput", func() error {
		if execSMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(execSMARN),
			Input:           aws.String(`{"value":42}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
			descResp, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: startResp.ExecutionArn,
			})
			if err != nil {
				return fmt.Errorf("describe: %v", err)
			}
			if descResp.Status == types.ExecutionStatusSucceeded {
				if descResp.Output == nil {
					return fmt.Errorf("execution output is nil")
				}
				if *descResp.Output != `"hello"` {
					return fmt.Errorf("expected output %q, got %q", `"hello"`, *descResp.Output)
				}
				return nil
			}
			if descResp.Status == types.ExecutionStatusFailed || descResp.Status == types.ExecutionStatusAborted {
				return fmt.Errorf("execution failed with status %s", descResp.Status)
			}
		}
		return fmt.Errorf("execution did not complete in time")
	}))

	results = append(results, r.RunTest("stepfunctions", "StopExecution", func() error {
		longDef := `{"Comment":"long","StartAt":"Wait","States":{"Wait":{"Type":"Wait","Seconds":300,"End":true}}}`
		longSMName := fmt.Sprintf("LongSM-%d", time.Now().UnixNano())
		_, longRoleARN, longRoleCleanup := tc.createRoleForSM("LongRole")
		defer longRoleCleanup()

		longResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(longSMName),
			Definition: aws.String(longDef),
			RoleArn:    aws.String(longRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		longSMARN := *longResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(longSMARN)})

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(longSMARN),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		stopResp, err := tc.client.StopExecution(tc.ctx, &sfn.StopExecutionInput{
			ExecutionArn: aws.String(*startResp.ExecutionArn),
			Error:        aws.String("TestError"),
			Cause:        aws.String("Integration test stop"),
		})
		if err != nil {
			return err
		}
		if stopResp.StopDate.IsZero() {
			return fmt.Errorf("stop date is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "StartSyncExecution", func() error {
		syncSMName := fmt.Sprintf("SyncSM-%d", time.Now().UnixNano())
		_, syncRoleARN, syncRoleCleanup := tc.createRoleForSM("SyncRole")
		defer syncRoleCleanup()

		syncResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(syncSMName),
			Definition: aws.String(passDef),
			RoleArn:    aws.String(syncRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		syncSMARN := *syncResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(syncSMARN)})

		body := map[string]interface{}{
			"stateMachineArn": syncSMARN,
			"input":           `{"sync":true}`,
		}
		bodyBytes, _ := json.Marshal(body)
		syncReq, _ := http.NewRequestWithContext(tc.ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		syncReq.Header.Set("Content-Type", "application/x-amz-json-1.0")
		syncReq.Header.Set("X-Amz-Target", "AWSStepFunctions.StartSyncExecution")

		httpResp, err := http.DefaultClient.Do(syncReq)
		if err != nil {
			return err
		}
		defer httpResp.Body.Close()
		if httpResp.StatusCode != 200 {
			var errBody map[string]interface{}
			json.NewDecoder(httpResp.Body).Decode(&errBody)
			return fmt.Errorf("status %d: %v", httpResp.StatusCode, errBody)
		}
		var result map[string]interface{}
		json.NewDecoder(httpResp.Body).Decode(&result)
		if status, _ := result["status"].(string); status != "SUCCEEDED" {
			return fmt.Errorf("expected SUCCEEDED, got %s", status)
		}
		if _, ok := result["executionArn"]; !ok {
			return fmt.Errorf("executionArn missing from response")
		}
		return nil
	}))

	tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{
		StateMachineArn: aws.String(execSMARN),
	})

	return results
}
