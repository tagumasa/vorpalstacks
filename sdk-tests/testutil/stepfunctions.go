package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunStepFunctionsTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "stepfunctions",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := sfn.NewFromConfig(cfg)
	ctx := context.Background()

	stateMachineName := fmt.Sprintf("TestStateMachine-%d", time.Now().UnixNano())
	activityName := fmt.Sprintf("TestActivity-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("TestSfnRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)

	trustPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"states.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	form := url.Values{}
	form.Set("Action", "CreateRole")
	form.Set("Version", "2010-05-08")
	form.Set("RoleName", roleName)
	form.Set("AssumeRolePolicyDocument", trustPolicy)
	req, err := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		results = append(results, TestResult{Service: "stepfunctions", TestName: "Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role request: %v", err)})
		return results
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		results = append(results, TestResult{Service: "stepfunctions", TestName: "Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
		return results
	}
	resp.Body.Close()
	defer func() {
		cleanupForm := url.Values{}
		cleanupForm.Set("Action", "DeleteRole")
		cleanupForm.Set("Version", "2010-05-08")
		cleanupForm.Set("RoleName", roleName)
		cleanupReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(cleanupForm.Encode()))
		if cleanupReq != nil {
			cleanupReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if cleanupResp, _ := http.DefaultClient.Do(cleanupReq); cleanupResp != nil {
				cleanupResp.Body.Close()
			}
		}
	}()

	stateMachineDefinition := map[string]interface{}{
		"Comment": "A Hello World example",
		"StartAt": "HelloWorld",
		"States": map[string]interface{}{
			"HelloWorld": map[string]interface{}{
				"Type": "Pass",
				"End":  true,
			},
		},
	}
	definitionJSON, _ := json.Marshal(stateMachineDefinition)

	var stateMachineARN string
	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine", func() error {
		resp, err := client.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(stateMachineName),
			Definition: aws.String(string(definitionJSON)),
			RoleArn:    aws.String(roleARN),
		})
		if err != nil {
			return err
		}
		stateMachineARN = *resp.StateMachineArn
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachine", func() error {
		resp, err := client.DescribeStateMachine(ctx, &sfn.DescribeStateMachineInput{
			StateMachineArn: aws.String(stateMachineARN),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineArn == nil {
			return fmt.Errorf("state machine ARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListStateMachines", func() error {
		resp, err := client.ListStateMachines(ctx, &sfn.ListStateMachinesInput{})
		if err != nil {
			return err
		}
		if resp.StateMachines == nil {
			return fmt.Errorf("state machines list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "StartExecution", func() error {
		input := map[string]string{"message": "test"}
		inputJSON, _ := json.Marshal(input)
		resp, err := client.StartExecution(ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(stateMachineARN),
			Input:           aws.String(string(inputJSON)),
		})
		if err != nil {
			return err
		}
		if resp.ExecutionArn == nil {
			return fmt.Errorf("execution ARN is nil")
		}
		return nil
	}))

	var executionARN string
	results = append(results, r.RunTest("stepfunctions", "ListExecutions", func() error {
		resp, err := client.ListExecutions(ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(stateMachineARN),
		})
		if err != nil {
			return err
		}
		if resp.Executions == nil {
			return fmt.Errorf("executions list is nil")
		}
		if len(resp.Executions) > 0 {
			executionARN = *resp.Executions[0].ExecutionArn
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeExecution", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		resp, err := client.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return err
		}
		if resp.Status == "" {
			return fmt.Errorf("execution status is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "GetExecutionHistory", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		resp, err := client.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return err
		}
		if resp.Events == nil {
			return fmt.Errorf("events list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachine", func() error {
		resp, err := client.UpdateStateMachine(ctx, &sfn.UpdateStateMachineInput{
			StateMachineArn: aws.String(stateMachineARN),
			Definition:      aws.String(string(definitionJSON)),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	var activityARN string
	results = append(results, r.RunTest("stepfunctions", "CreateActivity", func() error {
		resp, err := client.CreateActivity(ctx, &sfn.CreateActivityInput{
			Name: aws.String(activityName),
		})
		if err != nil {
			return err
		}
		activityARN = *resp.ActivityArn
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeActivity", func() error {
		resp, err := client.DescribeActivity(ctx, &sfn.DescribeActivityInput{
			ActivityArn: aws.String(activityARN),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("activity name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListActivities", func() error {
		resp, err := client.ListActivities(ctx, &sfn.ListActivitiesInput{})
		if err != nil {
			return err
		}
		if resp.Activities == nil {
			return fmt.Errorf("activities list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TagResource", func() error {
		resp, err := client.TagResource(ctx, &sfn.TagResourceInput{
			ResourceArn: aws.String(stateMachineARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &sfn.ListTagsForResourceInput{
			ResourceArn: aws.String(stateMachineARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UntagResource", func() error {
		resp, err := client.UntagResource(ctx, &sfn.UntagResourceInput{
			ResourceArn: aws.String(stateMachineARN),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteActivity", func() error {
		resp, err := client.DeleteActivity(ctx, &sfn.DeleteActivityInput{
			ActivityArn: aws.String(activityARN),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachine", func() error {
		resp, err := client.DeleteStateMachine(ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(stateMachineARN),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachine_NonExistent", func() error {
		_, err := client.DescribeStateMachine(ctx, &sfn.DescribeStateMachineInput{
			StateMachineArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:nonexistent-fake-arn"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent state machine")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachine_NonExistent", func() error {
		_, err := client.DeleteStateMachine(ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:nonexistent-fake-arn"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent state machine")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeExecution_NonExistent", func() error {
		_, err := client.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String("arn:aws:states:us-east-1:000000000000:execution:nonexistent:fake-exec"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent execution")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeActivity_NonExistent", func() error {
		_, err := client.DescribeActivity(ctx, &sfn.DescribeActivityInput{
			ActivityArn: aws.String("arn:aws:states:us-east-1:000000000000:activity:nonexistent-fake-arn"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent activity")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine_InvalidDefinition", func() error {
		invalidName := fmt.Sprintf("InvalidSM-%d", time.Now().UnixNano())
		invalidRoleName := fmt.Sprintf("InvalidRole-%d", time.Now().UnixNano())
		invalidForm := url.Values{}
		invalidForm.Set("Action", "CreateRole")
		invalidForm.Set("Version", "2010-05-08")
		invalidForm.Set("RoleName", invalidRoleName)
		invalidForm.Set("AssumeRolePolicyDocument", trustPolicy)
		invalidReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(invalidForm.Encode()))
		if invalidReq != nil {
			invalidReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if invalidResp, err := http.DefaultClient.Do(invalidReq); err == nil {
				invalidResp.Body.Close()
				defer func() {
					delForm := url.Values{}
					delForm.Set("Action", "DeleteRole")
					delForm.Set("Version", "2010-05-08")
					delForm.Set("RoleName", invalidRoleName)
					delReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(delForm.Encode()))
					if delReq != nil {
						delReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
						if delResp, _ := http.DefaultClient.Do(delReq); delResp != nil {
							delResp.Body.Close()
						}
					}
				}()
			}
		}
		_, err := client.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(invalidName),
			Definition: aws.String("not valid json {{{"),
			RoleArn:    aws.String(fmt.Sprintf("arn:aws:iam::000000000000:role/%s", invalidRoleName)),
		})
		if err != nil {
			defer client.DeleteStateMachine(ctx, &sfn.DeleteStateMachineInput{
				StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:%s:000000000000:stateMachine:%s", r.region, invalidName)),
			})
			return fmt.Errorf("server rejected invalid definition: %v", err)
		}
		defer client.DeleteStateMachine(ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:%s:000000000000:stateMachine:%s", r.region, invalidName)),
		})
		return nil
	}))

	var verifySMName string
	var verifySMARN string
	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachine_VerifyDefinition", func() error {
		smName := fmt.Sprintf("VerifySM-%d", time.Now().UnixNano())
		verifySMName = smName
		verifyRoleName := fmt.Sprintf("VerifyRole-%d", time.Now().UnixNano())
		verifyForm := url.Values{}
		verifyForm.Set("Action", "CreateRole")
		verifyForm.Set("Version", "2010-05-08")
		verifyForm.Set("RoleName", verifyRoleName)
		verifyForm.Set("AssumeRolePolicyDocument", trustPolicy)
		verifyReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(verifyForm.Encode()))
		if verifyReq != nil {
			verifyReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if verifyResp, err := http.DefaultClient.Do(verifyReq); err == nil {
				verifyResp.Body.Close()
				defer func() {
					delForm := url.Values{}
					delForm.Set("Action", "DeleteRole")
					delForm.Set("Version", "2010-05-08")
					delForm.Set("RoleName", verifyRoleName)
					delReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(delForm.Encode()))
					if delReq != nil {
						delReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
						if delResp, _ := http.DefaultClient.Do(delReq); delResp != nil {
							delResp.Body.Close()
						}
					}
				}()
			}
		}
		def1 := `{"Comment":"v1","StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
		resp, err := client.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(smName),
			Definition: aws.String(def1),
			RoleArn:    aws.String(fmt.Sprintf("arn:aws:iam::000000000000:role/%s", verifyRoleName)),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		verifySMARN = *resp.StateMachineArn

		def2 := `{"Comment":"v2","StartAt":"B","States":{"B":{"Type":"Pass","Result":"hello","End":true}}}`
		_, err = client.UpdateStateMachine(ctx, &sfn.UpdateStateMachineInput{
			StateMachineArn: aws.String(verifySMARN),
			Definition:      aws.String(def2),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		descResp, err := client.DescribeStateMachine(ctx, &sfn.DescribeStateMachineInput{
			StateMachineArn: aws.String(verifySMARN),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if *descResp.Definition != def2 {
			return fmt.Errorf("definition not updated: got %q, want %q", *descResp.Definition, def2)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Execution_PassStateOutput", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}

		startResp, err := client.StartExecution(ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(verifySMARN),
			Input:           aws.String(`{"value":42}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
			descResp, err := client.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
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

	results = append(results, r.RunTest("stepfunctions", "ListStateMachines_ContainsCreated", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		resp, err := client.ListStateMachines(ctx, &sfn.ListStateMachinesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, sm := range resp.StateMachines {
			if sm.StateMachineArn != nil && *sm.StateMachineArn == verifySMARN {
				found = true
				if sm.Name == nil || *sm.Name != verifySMName {
					return fmt.Errorf("state machine name mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("state machine not found in list")
		}
		return nil
	}))

	// === NEW TESTS: Version/Alias/Additional Operations ===

	results = append(results, r.RunTest("stepfunctions", "ValidateStateMachineDefinition_Valid", func() error {
		validDef := `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
		resp, err := client.ValidateStateMachineDefinition(ctx, &sfn.ValidateStateMachineDefinitionInput{
			Definition: aws.String(validDef),
		})
		if err != nil {
			return err
		}
		if resp.Result != types.ValidateStateMachineDefinitionResultCodeOk {
			return fmt.Errorf("expected OK, got %s", resp.Result)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ValidateStateMachineDefinition_Invalid", func() error {
		invalidDef := `{"StartAt":"Missing","States":{}}`
		resp, err := client.ValidateStateMachineDefinition(ctx, &sfn.ValidateStateMachineDefinitionInput{
			Definition: aws.String(invalidDef),
		})
		if err != nil {
			return err
		}
		if resp.Result == types.ValidateStateMachineDefinitionResultCodeOk {
			return fmt.Errorf("expected FAIL for invalid definition")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "GetStateMachine", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		getForm := url.Values{}
		getForm.Set("stateMachineArn", verifySMARN)
		getReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(getForm.Encode()))
		if getReq != nil {
			getReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			getReq.Header.Set("X-Amz-Target", "AWSStepFunctions.GetStateMachine")
			if getResp, err := http.DefaultClient.Do(getReq); err == nil {
				defer getResp.Body.Close()
				if getResp.StatusCode != 200 {
					return fmt.Errorf("GetStateMachine returned status %d", getResp.StatusCode)
				}
				var result map[string]interface{}
				if err := json.NewDecoder(getResp.Body).Decode(&result); err != nil {
					return fmt.Errorf("failed to decode response: %v", err)
				}
				if result["stateMachineArn"] != verifySMARN {
					return fmt.Errorf("state machine ARN mismatch")
				}
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachineForExecution", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		descExecResp, err := client.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return fmt.Errorf("describe execution: %v", err)
		}
		if descExecResp.StateMachineArn == nil {
			return fmt.Errorf("execution has no state machine ARN")
		}
		if descExecResp.Status != types.ExecutionStatusSucceeded && descExecResp.Status != types.ExecutionStatusRunning {
			return fmt.Errorf("execution not in suitable state: %s", descExecResp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "StartSyncExecution", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		body := map[string]interface{}{
			"stateMachineArn": verifySMARN,
			"input":           `{"sync":true}`,
		}
		bodyBytes, _ := json.Marshal(body)
		syncReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		if syncReq == nil {
			return fmt.Errorf("failed to create request")
		}
		syncReq.Header.Set("Content-Type", "application/x-amz-json-1.0")
		syncReq.Header.Set("X-Amz-Target", "AWSStepFunctions.StartSyncExecution")
		syncResp, err := http.DefaultClient.Do(syncReq)
		if err != nil {
			return err
		}
		defer syncResp.Body.Close()
		if syncResp.StatusCode != 200 {
			var errBody map[string]interface{}
			json.NewDecoder(syncResp.Body).Decode(&errBody)
			return fmt.Errorf("status %d: %v", syncResp.StatusCode, errBody)
		}
		var result map[string]interface{}
		json.NewDecoder(syncResp.Body).Decode(&result)
		if status, _ := result["status"].(string); status != "SUCCEEDED" {
			return fmt.Errorf("expected SUCCEEDED, got %s", status)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_Pass", func() error {
		def := `{"StartAt":"TestPass","States":{"TestPass":{"Type":"Pass","Result":{"hello":"world"},"End":true}}}`
		body := map[string]interface{}{
			"definition": def,
			"stateName":  "TestPass",
			"input":      `{}`,
		}
		bodyBytes, _ := json.Marshal(body)
		testReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		if testReq == nil {
			return fmt.Errorf("failed to create request")
		}
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

	// === Version Operations ===

	var versionARN string
	results = append(results, r.RunTest("stepfunctions", "PublishStateMachineVersion", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		resp, err := client.PublishStateMachineVersion(ctx, &sfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(verifySMARN),
			Description:     aws.String("test version 1"),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineVersionArn == nil {
			return fmt.Errorf("version ARN is nil")
		}
		versionARN = *resp.StateMachineVersionArn
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListStateMachineVersions", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		resp, err := client.ListStateMachineVersions(ctx, &sfn.ListStateMachineVersionsInput{
			StateMachineArn: aws.String(verifySMARN),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineVersions == nil || len(resp.StateMachineVersions) == 0 {
			return fmt.Errorf("expected at least one version")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "PublishStateMachineVersion_Second", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		resp, err := client.PublishStateMachineVersion(ctx, &sfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(verifySMARN),
			Description:     aws.String("test version 2"),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineVersionArn == nil {
			return fmt.Errorf("version ARN is nil")
		}
		if *resp.StateMachineVersionArn == versionARN {
			return fmt.Errorf("second version should have different ARN")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineVersion", func() error {
		if versionARN == "" {
			return fmt.Errorf("version ARN not available")
		}
		_, err := client.DeleteStateMachineVersion(ctx, &sfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: aws.String(versionARN),
		})
		if err != nil {
			return err
		}
		_, err = client.ListStateMachineVersions(ctx, &sfn.ListStateMachineVersionsInput{
			StateMachineArn: aws.String(verifySMARN),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineVersion_NonExistent", func() error {
		_, err := client.DeleteStateMachineVersion(ctx, &sfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:fake:999"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent version")
		}
		return nil
	}))

	// === Alias Operations ===

	var aliasARN string
	var aliasName string
	results = append(results, r.RunTest("stepfunctions", "CreateStateMachineAlias", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		aliasName = fmt.Sprintf("PROD-%d", time.Now().UnixNano())
		listResp, err := client.ListStateMachineVersions(ctx, &sfn.ListStateMachineVersionsInput{
			StateMachineArn: aws.String(verifySMARN),
		})
		if err != nil {
			return fmt.Errorf("list versions: %v", err)
		}
		if len(listResp.StateMachineVersions) == 0 {
			return fmt.Errorf("no versions available for alias")
		}
		latestVersion := *listResp.StateMachineVersions[0].StateMachineVersionArn

		aliasBody := map[string]interface{}{
			"name":        aliasName,
			"description": "production alias",
			"routingConfiguration": []map[string]interface{}{
				{"stateMachineVersionArn": latestVersion, "weight": float64(100)},
			},
		}
		bodyBytes, _ := json.Marshal(aliasBody)
		aliasReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		if aliasReq != nil {
			aliasReq.Header.Set("Content-Type", "application/x-amz-json-1.0")
			aliasReq.Header.Set("X-Amz-Target", "AWSStepFunctions.CreateStateMachineAlias")
			aliasResp, err := http.DefaultClient.Do(aliasReq)
			if err != nil {
				return err
			}
			defer aliasResp.Body.Close()
			if aliasResp.StatusCode != 200 {
				var errBody map[string]interface{}
				json.NewDecoder(aliasResp.Body).Decode(&errBody)
				return fmt.Errorf("status %d: %v", aliasResp.StatusCode, errBody)
			}
			var result map[string]interface{}
			json.NewDecoder(aliasResp.Body).Decode(&result)
			if arn, ok := result["stateMachineAliasArn"].(string); ok {
				aliasARN = arn
			} else {
				return fmt.Errorf("alias ARN not in response")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachineAlias", func() error {
		if aliasARN == "" {
			return fmt.Errorf("alias ARN not available")
		}
		resp, err := client.DescribeStateMachineAlias(ctx, &sfn.DescribeStateMachineAliasInput{
			StateMachineAliasArn: aws.String(aliasARN),
		})
		if err != nil {
			return fmt.Errorf("describe alias error: %v", err)
		}
		if resp.Name == nil || *resp.Name != aliasName {
			return fmt.Errorf("alias name mismatch: got %v, want %q", resp.Name, aliasName)
		}
		if resp.CreationDate.IsZero() {
			return fmt.Errorf("creation date is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListStateMachineAliases", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		listBody := map[string]interface{}{
			"stateMachineArn": verifySMARN,
		}
		listBytes, _ := json.Marshal(listBody)
		listReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewReader(listBytes))
		if listReq != nil {
			listReq.Header.Set("Content-Type", "application/x-amz-json-1.0")
			listReq.Header.Set("X-Amz-Target", "AWSStepFunctions.ListStateMachineAliases")
			listResp, err := http.DefaultClient.Do(listReq)
			if err != nil {
				return err
			}
			defer listResp.Body.Close()
			if listResp.StatusCode != 200 {
				var errBody map[string]interface{}
				json.NewDecoder(listResp.Body).Decode(&errBody)
				return fmt.Errorf("status %d: %v", listResp.StatusCode, errBody)
			}
			var result map[string]interface{}
			json.NewDecoder(listResp.Body).Decode(&result)
			aliases, _ := result["stateMachineAliases"].([]interface{})
			if len(aliases) == 0 {
				return fmt.Errorf("expected at least one alias, got %v (arn=%s)", result, verifySMARN)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachineAlias", func() error {
		if aliasARN == "" {
			return fmt.Errorf("alias ARN not available")
		}
		resp, err := client.UpdateStateMachineAlias(ctx, &sfn.UpdateStateMachineAliasInput{
			StateMachineAliasArn: aws.String(aliasARN),
			Description:          aws.String("updated production alias"),
		})
		if err != nil {
			return err
		}
		if resp.UpdateDate.IsZero() {
			return fmt.Errorf("update date is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineAlias", func() error {
		if aliasARN == "" {
			return fmt.Errorf("alias ARN not available")
		}
		_, err := client.DeleteStateMachineAlias(ctx, &sfn.DeleteStateMachineAliasInput{
			StateMachineAliasArn: aws.String(aliasARN),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineAlias_NonExistent", func() error {
		_, err := client.DeleteStateMachineAlias(ctx, &sfn.DeleteStateMachineAliasInput{
			StateMachineAliasArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:fake:NONEXISTENT"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent alias")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachineAlias_Duplicate", func() error {
		if verifySMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}
		listVerResp, err := client.ListStateMachineVersions(ctx, &sfn.ListStateMachineVersionsInput{
			StateMachineArn: aws.String(verifySMARN),
		})
		if err != nil {
			return fmt.Errorf("list versions: %v", err)
		}
		if len(listVerResp.StateMachineVersions) == 0 {
			return fmt.Errorf("no versions available")
		}
		realVersionArn := *listVerResp.StateMachineVersions[0].StateMachineVersionArn

		dupAliasName := fmt.Sprintf("DUP-%d", time.Now().UnixNano())
		aliasBody := map[string]interface{}{
			"name": dupAliasName,
			"routingConfiguration": []map[string]interface{}{
				{"stateMachineVersionArn": realVersionArn, "weight": float64(100)},
			},
		}
		bodyBytes, _ := json.Marshal(aliasBody)
		aliasReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		if aliasReq != nil {
			aliasReq.Header.Set("Content-Type", "application/x-amz-json-1.0")
			aliasReq.Header.Set("X-Amz-Target", "AWSStepFunctions.CreateStateMachineAlias")
			aliasResp, err := http.DefaultClient.Do(aliasReq)
			if err != nil {
				return fmt.Errorf("first create: %v", err)
			}
			aliasResp.Body.Close()
			if aliasResp.StatusCode != 200 {
				return fmt.Errorf("first create returned %d", aliasResp.StatusCode)
			}
		}
		aliasReq2, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		if aliasReq2 != nil {
			aliasReq2.Header.Set("Content-Type", "application/x-amz-json-1.0")
			aliasReq2.Header.Set("X-Amz-Target", "AWSStepFunctions.CreateStateMachineAlias")
			aliasResp2, err := http.DefaultClient.Do(aliasReq2)
			if err != nil {
				return fmt.Errorf("second create request: %v", err)
			}
			defer aliasResp2.Body.Close()
			if aliasResp2.StatusCode == 200 {
				return fmt.Errorf("expected error for duplicate alias")
			}
		}
		return nil
	}))

	if verifySMARN != "" {
		client.DeleteStateMachine(ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(verifySMARN),
		})
	}

	return results
}
