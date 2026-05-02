package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

func (r *TestRunner) runSFNEdgeTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachine_NonExistent", func() error {
		_, err := tc.client.DescribeStateMachine(tc.ctx, &sfn.DescribeStateMachineInput{
			StateMachineArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:nonexistent-fake-arn"),
		})
		if err := AssertErrorContains(err, "StateMachineDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachine_NonExistent", func() error {
		_, err := tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:nonexistent-fake-arn"),
		})
		if err := AssertErrorContains(err, "StateMachineDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeExecution_NonExistent", func() error {
		_, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String("arn:aws:states:us-east-1:000000000000:execution:nonexistent:fake-exec"),
		})
		if err := AssertErrorContains(err, "ExecutionDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeActivity_NonExistent", func() error {
		_, err := tc.client.DescribeActivity(tc.ctx, &sfn.DescribeActivityInput{
			ActivityArn: aws.String("arn:aws:states:us-east-1:000000000000:activity:nonexistent-fake-arn"),
		})
		if err := AssertErrorContains(err, "ActivityDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineVersion_NonExistent", func() error {
		_, err := tc.client.DeleteStateMachineVersion(tc.ctx, &sfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:fake:999"),
		})
		if err := AssertErrorContains(err, "StateMachineVersionNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineAlias_NonExistent", func() error {
		_, err := tc.client.DeleteStateMachineAlias(tc.ctx, &sfn.DeleteStateMachineAliasInput{
			StateMachineAliasArn: aws.String("arn:aws:states:us-east-1:000000000000:stateMachine:fake:NONEXISTENT"),
		})
		if err := AssertErrorContains(err, "StateMachineAliasDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine_InvalidDefinition", func() error {
		invalidName := fmt.Sprintf("InvalidSM-%d", time.Now().UnixNano())
		_, invalidRoleARN, invalidRoleCleanup := tc.createRoleForSM("InvalidRole")
		defer invalidRoleCleanup()
		_, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(invalidName),
			Definition: aws.String("not valid json {{{"),
			RoleArn:    aws.String(invalidRoleARN),
		})
		if err == nil {
			return fmt.Errorf("server accepted invalid definition, expected InvalidDefinition error")
		}
		if err := AssertErrorContains(err, "InvalidDefinition"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListStateMachines_ContainsCreated", func() error {
		listSMName := fmt.Sprintf("ListSM-%d", time.Now().UnixNano())
		_, listRoleARN, listRoleCleanup := tc.createRoleForSM("ListRole")
		defer listRoleCleanup()

		createResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(listSMName),
			Definition: aws.String(`{"Comment":"list test","StartAt":"Done","States":{"Done":{"Type":"Pass","End":true}}}`),
			RoleArn:    aws.String(listRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		listSMARN := *createResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(listSMARN)})

		resp, err := tc.client.ListStateMachines(tc.ctx, &sfn.ListStateMachinesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, sm := range resp.StateMachines {
			if sm.StateMachineArn != nil && *sm.StateMachineArn == listSMARN {
				found = true
				if sm.Name == nil || *sm.Name != listSMName {
					return fmt.Errorf("state machine name mismatch")
				}
				if sm.CreationDate.IsZero() {
					return fmt.Errorf("creation date is zero in list")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("state machine not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListStateMachines_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgARNs []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagSM-%s-%d", pgTs, i)
			resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
				Name:       aws.String(name),
				Definition: aws.String(`{"Comment":"pag test","StartAt":"Done","States":{"Done":{"Type":"Pass","End":true}}}`),
				RoleArn:    aws.String(tc.roleARN),
			})
			if err != nil {
				for _, arn := range pgARNs {
					tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(arn)})
				}
				return fmt.Errorf("create state machine %s: %v", name, err)
			}
			pgARNs = append(pgARNs, *resp.StateMachineArn)
		}

		var allNames []string
		var nextToken *string
		for {
			resp, err := tc.client.ListStateMachines(tc.ctx, &sfn.ListStateMachinesInput{
				MaxResults: 2,
				NextToken:  nextToken,
			})
			if err != nil {
				for _, arn := range pgARNs {
					tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(arn)})
				}
				return fmt.Errorf("list state machines page: %v", err)
			}
			for _, sm := range resp.StateMachines {
				if strings.HasPrefix(aws.ToString(sm.Name), "PagSM-"+pgTs) {
					allNames = append(allNames, aws.ToString(sm.Name))
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, arn := range pgARNs {
			tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(arn)})
		}
		if len(allNames) != 5 {
			return fmt.Errorf("expected 5 paginated state machines, got %d", len(allNames))
		}
		return nil
	}))

	_ = types.StateMachineStatusActive
	return results
}
