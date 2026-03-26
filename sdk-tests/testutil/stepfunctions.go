package testutil

import (
	"context"
	"encoding/json"
	"fmt"
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
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/TestRole")

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
		_, err := client.StartExecution(ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(stateMachineARN),
			Input:           aws.String(string(inputJSON)),
		})
		return err
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
		_, err := client.UpdateStateMachine(ctx, &sfn.UpdateStateMachineInput{
			StateMachineArn: aws.String(stateMachineARN),
			Definition:      aws.String(string(definitionJSON)),
		})
		return err
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
		_, err := client.TagResource(ctx, &sfn.TagResourceInput{
			ResourceArn: aws.String(stateMachineARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
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
		_, err := client.UntagResource(ctx, &sfn.UntagResourceInput{
			ResourceArn: aws.String(stateMachineARN),
			TagKeys:     []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteActivity", func() error {
		_, err := client.DeleteActivity(ctx, &sfn.DeleteActivityInput{
			ActivityArn: aws.String(activityARN),
		})
		return err
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachine", func() error {
		_, err := client.DeleteStateMachine(ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(stateMachineARN),
		})
		return err
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
		_, err := client.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(invalidName),
			Definition: aws.String("not valid json {{{"),
			RoleArn:    aws.String("arn:aws:iam::000000000000:role/TestRole"),
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
		def1 := `{"Comment":"v1","StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
		resp, err := client.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(smName),
			Definition: aws.String(def1),
			RoleArn:    aws.String("arn:aws:iam::000000000000:role/TestRole"),
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

	if verifySMARN != "" {
		client.DeleteStateMachine(ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(verifySMARN),
		})
	}

	return results
}
