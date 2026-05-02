package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

func (r *TestRunner) runSFNAliasTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	aliasSMName := fmt.Sprintf("AliasSM-%d", time.Now().UnixNano())
	_, aliasRoleARN, aliasRoleCleanup := tc.createRoleForSM("AliasRole")
	defer aliasRoleCleanup()

	aliasDef := `{"Comment":"alias test","StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
	aliasResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(aliasSMName),
		Definition: aws.String(aliasDef),
		RoleArn:    aws.String(aliasRoleARN),
	})
	if err != nil {
		return append(results, TestResult{Service: "stepfunctions", TestName: "AliasSetup", Status: "FAIL", Error: fmt.Sprintf("create SM: %v", err)})
	}
	aliasSMARN := *aliasResp.StateMachineArn
	defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(aliasSMARN)})

	verResp, err := tc.client.PublishStateMachineVersion(tc.ctx, &sfn.PublishStateMachineVersionInput{
		StateMachineArn: aws.String(aliasSMARN),
		Description:     aws.String("alias version"),
	})
	if err != nil {
		return append(results, TestResult{Service: "stepfunctions", TestName: "AliasSetup", Status: "FAIL", Error: fmt.Sprintf("publish version: %v", err)})
	}
	versionARN := *verResp.StateMachineVersionArn
	defer tc.client.DeleteStateMachineVersion(tc.ctx, &sfn.DeleteStateMachineVersionInput{StateMachineVersionArn: aws.String(versionARN)})

	var aliasARN string
	var aliasName string
	results = append(results, r.RunTest("stepfunctions", "CreateStateMachineAlias", func() error {
		aliasName = fmt.Sprintf("PROD-%d", time.Now().UnixNano())
		aliasResp, err := tc.client.CreateStateMachineAlias(tc.ctx, &sfn.CreateStateMachineAliasInput{
			Name:        aws.String(aliasName),
			Description: aws.String("production alias"),
			RoutingConfiguration: []types.RoutingConfigurationListItem{
				{
					StateMachineVersionArn: aws.String(versionARN),
					Weight:                 100,
				},
			},
		})
		if err != nil {
			return err
		}
		if aliasResp.StateMachineAliasArn == nil || *aliasResp.StateMachineAliasArn == "" {
			return fmt.Errorf("alias ARN is nil or empty")
		}
		aliasARN = *aliasResp.StateMachineAliasArn
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachineAlias", func() error {
		if aliasARN == "" {
			return fmt.Errorf("alias ARN not available")
		}
		resp, err := tc.client.DescribeStateMachineAlias(tc.ctx, &sfn.DescribeStateMachineAliasInput{
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
		resp, err := tc.client.ListStateMachineAliases(tc.ctx, &sfn.ListStateMachineAliasesInput{
			StateMachineArn: aws.String(aliasSMARN),
		})
		if err != nil {
			return err
		}
		if len(resp.StateMachineAliases) == 0 {
			return fmt.Errorf("expected at least one alias")
		}
		found := false
		for _, a := range resp.StateMachineAliases {
			if a.StateMachineAliasArn != nil && *a.StateMachineAliasArn == aliasARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created alias not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachineAlias", func() error {
		if aliasARN == "" {
			return fmt.Errorf("alias ARN not available")
		}
		resp, err := tc.client.UpdateStateMachineAlias(tc.ctx, &sfn.UpdateStateMachineAliasInput{
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
		_, err := tc.client.DeleteStateMachineAlias(tc.ctx, &sfn.DeleteStateMachineAliasInput{
			StateMachineAliasArn: aws.String(aliasARN),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachineAlias_Duplicate", func() error {
		dupAliasName := fmt.Sprintf("DUP-%d", time.Now().UnixNano())
		_, err := tc.client.CreateStateMachineAlias(tc.ctx, &sfn.CreateStateMachineAliasInput{
			Name: aws.String(dupAliasName),
			RoutingConfiguration: []types.RoutingConfigurationListItem{
				{
					StateMachineVersionArn: aws.String(versionARN),
					Weight:                 100,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		_, err = tc.client.CreateStateMachineAlias(tc.ctx, &sfn.CreateStateMachineAliasInput{
			Name: aws.String(dupAliasName),
			RoutingConfiguration: []types.RoutingConfigurationListItem{
				{
					StateMachineVersionArn: aws.String(versionARN),
					Weight:                 100,
				},
			},
		})
		if err == nil {
			tc.client.DeleteStateMachineAlias(tc.ctx, &sfn.DeleteStateMachineAliasInput{
				StateMachineAliasArn: aws.String(fmt.Sprintf("arn:aws:states:%s:000000000000:stateMachineAlias:%s", r.region, dupAliasName)),
			})
			return fmt.Errorf("expected error for duplicate alias")
		}
		return nil
	}))

	return results
}
