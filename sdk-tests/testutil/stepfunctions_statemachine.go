package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

func (r *TestRunner) runSFNStateMachineTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	smName := fmt.Sprintf("TestSM-%d", time.Now().UnixNano())
	var smARN string

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine", func() error {
		resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(smName),
			Definition: aws.String(tc.definition),
			RoleArn:    aws.String(tc.roleARN),
			Type:       types.StateMachineTypeStandard,
		})
		if err != nil {
			return err
		}
		if resp.StateMachineArn == nil {
			return fmt.Errorf("state machine ARN is nil")
		}
		smARN = *resp.StateMachineArn
		if smARN == "" {
			return fmt.Errorf("state machine ARN is empty")
		}
		if resp.CreationDate.IsZero() {
			return fmt.Errorf("creation date is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachine", func() error {
		resp, err := tc.client.DescribeStateMachine(tc.ctx, &sfn.DescribeStateMachineInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineArn == nil || *resp.StateMachineArn != smARN {
			return fmt.Errorf("state machine ARN mismatch: got %v", resp.StateMachineArn)
		}
		if resp.Name == nil || *resp.Name != smName {
			return fmt.Errorf("name mismatch: got %v, want %q", resp.Name, smName)
		}
		if resp.Status != types.StateMachineStatusActive {
			return fmt.Errorf("expected ACTIVE status, got %s", resp.Status)
		}
		if resp.RoleArn == nil || *resp.RoleArn != tc.roleARN {
			return fmt.Errorf("role ARN mismatch: got %v", resp.RoleArn)
		}
		if resp.Definition == nil || *resp.Definition == "" {
			return fmt.Errorf("definition is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "GetStateMachine", func() error {
		resp, err := tc.client.DescribeStateMachine(tc.ctx, &sfn.DescribeStateMachineInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineArn == nil || *resp.StateMachineArn != smARN {
			return fmt.Errorf("state machine ARN mismatch")
		}
		if resp.Type != types.StateMachineTypeStandard {
			return fmt.Errorf("type mismatch: got %s", resp.Type)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListStateMachines", func() error {
		resp, err := tc.client.ListStateMachines(tc.ctx, &sfn.ListStateMachinesInput{})
		if err != nil {
			return err
		}
		if resp.StateMachines == nil {
			return fmt.Errorf("state machines list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachine", func() error {
		resp, err := tc.client.UpdateStateMachine(tc.ctx, &sfn.UpdateStateMachineInput{
			StateMachineArn: aws.String(smARN),
			Definition:      aws.String(tc.definition),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.UpdateDate.IsZero() {
			return fmt.Errorf("update date is zero")
		}
		return nil
	}))

	var verifySMARN string
	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachine_VerifyDefinition", func() error {
		vSMName := fmt.Sprintf("VerifySM-%d", time.Now().UnixNano())
		_, vRoleARN, vCleanup := tc.createRoleForSM("VerifyRole")
		defer vCleanup()

		def1 := `{"Comment":"v1","StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
		resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(vSMName),
			Definition: aws.String(def1),
			RoleArn:    aws.String(vRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		verifySMARN = *resp.StateMachineArn

		def2 := `{"Comment":"v2","StartAt":"B","States":{"B":{"Type":"Pass","Result":"hello","End":true}}}`
		updateResp, err := tc.client.UpdateStateMachine(tc.ctx, &sfn.UpdateStateMachineInput{
			StateMachineArn: aws.String(verifySMARN),
			Definition:      aws.String(def2),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if updateResp.UpdateDate.IsZero() {
			return fmt.Errorf("update date is zero")
		}

		descResp, err := tc.client.DescribeStateMachine(tc.ctx, &sfn.DescribeStateMachineInput{
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

	results = append(results, r.RunTest("stepfunctions", "ValidateStateMachineDefinition_Valid", func() error {
		validDef := `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
		resp, err := tc.client.ValidateStateMachineDefinition(tc.ctx, &sfn.ValidateStateMachineDefinitionInput{
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
		resp, err := tc.client.ValidateStateMachineDefinition(tc.ctx, &sfn.ValidateStateMachineDefinitionInput{
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

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachine", func() error {
		resp, err := tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	if verifySMARN != "" {
		tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(verifySMARN),
		})
	}

	return results
}
