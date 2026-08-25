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
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:%s:%s:stateMachine:nonexistent-fake-arn", r.region, r.accountID)),
		})
		if err := AssertErrorContains(err, "StateMachineDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachine_NonExistent", func() error {
		_, err := tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:%s:%s:stateMachine:nonexistent-fake-arn", r.region, r.accountID)),
		})
		if err := AssertErrorContains(err, "StateMachineDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeExecution_NonExistent", func() error {
		_, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(fmt.Sprintf("arn:aws:states:%s:%s:execution:nonexistent:fake-exec", r.region, r.accountID)),
		})
		if err := AssertErrorContains(err, "ExecutionDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeActivity_NonExistent", func() error {
		_, err := tc.client.DescribeActivity(tc.ctx, &sfn.DescribeActivityInput{
			ActivityArn: aws.String(fmt.Sprintf("arn:aws:states:%s:%s:activity:nonexistent-fake-arn", r.region, r.accountID)),
		})
		if err := AssertErrorContains(err, "ActivityDoesNotExist"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineVersion_NonExistent", func() error {
		_, err := tc.client.DeleteStateMachineVersion(tc.ctx, &sfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: aws.String(fmt.Sprintf("arn:aws:states:%s:%s:stateMachine:fake:999", r.region, r.accountID)),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineAlias_NonExistent", func() error {
		_, err := tc.client.DeleteStateMachineAlias(tc.ctx, &sfn.DeleteStateMachineAliasInput{
			StateMachineAliasArn: aws.String(fmt.Sprintf("arn:aws:states:%s:%s:stateMachine:fake:NONEXISTENT", r.region, r.accountID)),
		})
		if err := AssertErrorContains(err, "ResourceNotFound"); err != nil {
			return err
		}
		return nil
	}))

	// Malformed definitions are rejected with InvalidDefinition at
	// creation time; each row exercises a distinct malformation class.
	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine_DefinitionRejections", func() error {
		rejections := []struct {
			label      string
			definition string
		}{
			{"MalformedJSON", "not valid json {{{"},
			{"UnknownStateField", `{"StartAt":"S","States":{"S":{"Type":"Pass","ResultPat":"$.x","End":true}}}`},
			{"MissingTransitionTarget", `{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"Missing"}}}`},
		}
		for _, row := range rejections {
			_, roleARN, roleCleanup := tc.createRoleForSM("InvalidDefRole")
			defer roleCleanup()
			_, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
				Name:       aws.String(fmt.Sprintf("InvalidDef-%s-%d", row.label, time.Now().UnixNano())),
				Definition: aws.String(row.definition),
				RoleArn:    aws.String(roleARN),
			})
			if cerr := expectAWSErrorCode(err, "InvalidDefinition"); cerr != nil {
				return fmt.Errorf("%s: %v", row.label, cerr)
			}
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

	// statusFilter must be a member of the Smithy ExecutionStatus enum;
	// invalid values are rejected instead of silently returning nothing.
	results = append(results, r.RunTest("stepfunctions", "ListExecutions_InvalidStatusFilter_Rejected", func() error {
		smArn, err := tc.createTestSM(fmt.Sprintf("inv-status-%d", time.Now().UnixNano()))
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(smArn)})

		_, err = tc.client.ListExecutions(tc.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(smArn),
			StatusFilter:    "NOT_A_STATUS",
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for invalid statusFilter, got nil")
		}
		return nil
	}))

	// Execution names are limited to 80 characters.
	results = append(results, r.RunTest("stepfunctions", "StartExecution_NameTooLong_Rejected", func() error {
		smArn, err := tc.createTestSM(fmt.Sprintf("long-exec-name-%d", time.Now().UnixNano()))
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(smArn)})

		longName := strings.Repeat("n", 81)
		_, err = tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(smArn),
			Name:            aws.String(longName),
		})
		if err == nil {
			return fmt.Errorf("expected InvalidName for 81-character execution name, got nil")
		}
		return nil
	}))

	// The tagging quota is fifty tags per resource.
	results = append(results, r.RunTest("stepfunctions", "TagResource_TooManyTags_Rejected", func() error {
		smArn, err := tc.createTestSM(fmt.Sprintf("too-many-tags-%d", time.Now().UnixNano()))
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(smArn)})

		tags := make([]types.Tag, 51)
		for i := range tags {
			tags[i] = types.Tag{Key: aws.String(fmt.Sprintf("k%d", i)), Value: aws.String("v")}
		}
		_, err = tc.client.TagResource(tc.ctx, &sfn.TagResourceInput{
			ResourceArn: aws.String(smArn),
			Tags:        tags,
		})
		if err == nil {
			return fmt.Errorf("expected TooManyTags for 51 tags, got nil")
		}
		return nil
	}))

	// Required ARN parameters must fail with InvalidArn rather than a
	// not-found error.
	results = append(results, r.RunTest("stepfunctions", "DescribeActivity_EmptyArn_InvalidArn", func() error {
		_, err := tc.client.DescribeActivity(tc.ctx, &sfn.DescribeActivityInput{
			ActivityArn: aws.String(""),
		})
		if err := AssertErrorContains(err, "InvalidArn"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
