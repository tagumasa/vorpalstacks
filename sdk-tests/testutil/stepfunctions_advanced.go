package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

func (r *TestRunner) runSFNAdvancedTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("stepfunctions", "TestState_Pass", func() error {
		def := `{"StartAt":"TestPass","States":{"TestPass":{"Type":"Pass","Result":{"hello":"world"},"End":true}}}`
		body := map[string]interface{}{
			"definition": def,
			"stateName":  "TestPass",
			"input":      `{}`,
		}
		bodyBytes, _ := json.Marshal(body)
		testReq, _ := http.NewRequestWithContext(tc.ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		testReq.Header.Set("Content-Type", "application/x-amz-json-1.0")
		testReq.Header.Set("X-Amz-Target", "AWSStepFunctions.TestState")

		testResp, err := http.DefaultClient.Do(testReq)
		if err != nil {
			return err
		}
		defer testResp.Body.Close()
		if testResp.StatusCode != 200 {
			var errBody map[string]interface{}
			json.NewDecoder(testResp.Body).Decode(&errBody)
			return fmt.Errorf("status %d: %v", testResp.StatusCode, errBody)
		}
		var result map[string]interface{}
		json.NewDecoder(testResp.Body).Decode(&result)
		if status, _ := result["status"].(string); status != "SUCCEEDED" {
			return fmt.Errorf("expected SUCCEEDED, got %s", status)
		}
		if _, ok := result["output"]; !ok {
			return fmt.Errorf("output is missing")
		}
		return nil
	}))

	mapSMName := fmt.Sprintf("MapSM-%d", time.Now().UnixNano())
	_, mapRoleARN, mapRoleCleanup := tc.createRoleForSM("MapRole")
	defer mapRoleCleanup()

	mapDef := `{"Comment":"map test","StartAt":"Map","States":{"Map":{"Type":"Map","Iterator":{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}},"End":true}}}`
	var mapSMARN string
	mapResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(mapSMName),
		Definition: aws.String(mapDef),
		RoleArn:    aws.String(mapRoleARN),
	})
	if err != nil {
		return append(results, TestResult{Service: "stepfunctions", TestName: "MapRunSetup", Status: "FAIL", Error: fmt.Sprintf("create SM: %v", err)})
	}
	mapSMARN = *mapResp.StateMachineArn
	defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(mapSMARN)})

	results = append(results, r.RunTest("stepfunctions", "ListMapRuns", func() error {
		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(mapSMARN),
			Input:           aws.String(`[1,2,3]`),
		})
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}

		var listResp *sfn.ListMapRunsOutput
		for i := 0; i < 20; i++ {
			time.Sleep(100 * time.Millisecond)
			listResp, err = tc.client.ListMapRuns(tc.ctx, &sfn.ListMapRunsInput{
				ExecutionArn: aws.String(*execResp.ExecutionArn),
			})
			if err != nil {
				return err
			}
			if len(listResp.MapRuns) > 0 {
				break
			}
		}
		if listResp == nil || len(listResp.MapRuns) == 0 {
			return fmt.Errorf("no map runs found after polling")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeMapRun", func() error {
		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(mapSMARN),
			Input:           aws.String(`[1,2,3]`),
		})
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}

		var mapRunARN string
		for i := 0; i < 20; i++ {
			time.Sleep(100 * time.Millisecond)
			listResp, err := tc.client.ListMapRuns(tc.ctx, &sfn.ListMapRunsInput{
				ExecutionArn: aws.String(*execResp.ExecutionArn),
			})
			if err != nil {
				return err
			}
			if len(listResp.MapRuns) > 0 && listResp.MapRuns[0].MapRunArn != nil {
				mapRunARN = *listResp.MapRuns[0].MapRunArn
				break
			}
		}
		if mapRunARN == "" {
			return fmt.Errorf("no map runs found after polling")
		}

		descResp, err := tc.client.DescribeMapRun(tc.ctx, &sfn.DescribeMapRunInput{
			MapRunArn: aws.String(mapRunARN),
		})
		if err != nil {
			return err
		}
		if descResp.MapRunArn == nil || *descResp.MapRunArn == "" {
			return fmt.Errorf("map run ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateMapRun", func() error {
		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(mapSMARN),
			Input:           aws.String(`[1,2,3]`),
		})
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}

		var mapRunARN string
		for i := 0; i < 20; i++ {
			time.Sleep(100 * time.Millisecond)
			listResp, err := tc.client.ListMapRuns(tc.ctx, &sfn.ListMapRunsInput{
				ExecutionArn: aws.String(*execResp.ExecutionArn),
			})
			if err != nil {
				return err
			}
			if len(listResp.MapRuns) > 0 && listResp.MapRuns[0].MapRunArn != nil {
				mapRunARN = *listResp.MapRuns[0].MapRunArn
				break
			}
		}
		if mapRunARN == "" {
			return fmt.Errorf("no map runs found after polling")
		}

		_, err = tc.client.UpdateMapRun(tc.ctx, &sfn.UpdateMapRunInput{
			MapRunArn:      aws.String(mapRunARN),
			MaxConcurrency: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "RedriveExecution", func() error {
		failDef := `{"Comment":"fail then redrive","StartAt":"Fail","States":{"Fail":{"Type":"Fail","Error":"TestError","Cause":"Redrive test"}}}`
		failSMName := fmt.Sprintf("RedriveSM-%d", time.Now().UnixNano())
		_, failRoleARN, failRoleCleanup := tc.createRoleForSM("RedriveRole")
		defer failRoleCleanup()

		failResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(failSMName),
			Definition: aws.String(failDef),
			RoleArn:    aws.String(failRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		failSMARN := *failResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(failSMARN)})

		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(failSMARN),
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		redriveResp, err := tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(*execResp.ExecutionArn),
		})
		if err != nil {
			return err
		}
		if redriveResp == nil {
			return fmt.Errorf("redrive response is nil")
		}
		return nil
	}))

	return results
}
