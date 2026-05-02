package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

func (r *TestRunner) runNeptuneInstanceParamGroupTests(tc *neptuneContext) []TestResult {
	var results []TestResult
	instancePGName := fmt.Sprintf("test-pg-%d", tc.ts)

	results = append(results, r.RunTest("neptune", "CreateDBParameterGroup", func() error {
		resp, err := tc.client.CreateDBParameterGroup(tc.ctx, &neptune.CreateDBParameterGroupInput{
			DBParameterGroupName:   aws.String(instancePGName),
			DBParameterGroupFamily: aws.String("neptune1"),
			Description:            aws.String("Test instance parameter group"),
		})
		if err != nil {
			return err
		}
		if resp.DBParameterGroup == nil {
			return fmt.Errorf("expected DBParameterGroup in response")
		}
		pg := resp.DBParameterGroup
		if pg.DBParameterGroupName == nil || *pg.DBParameterGroupName != instancePGName {
			return fmt.Errorf("expected DBParameterGroupName=%s, got %v", instancePGName, pg.DBParameterGroupName)
		}
		if pg.DBParameterGroupFamily == nil || *pg.DBParameterGroupFamily != "neptune1" {
			return fmt.Errorf("expected DBParameterGroupFamily=neptune1, got %v", pg.DBParameterGroupFamily)
		}
		if pg.Description == nil || *pg.Description != "Test instance parameter group" {
			return fmt.Errorf("expected Description='Test instance parameter group', got %v", pg.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBParameterGroups", func() error {
		resp, err := tc.client.DescribeDBParameterGroups(tc.ctx, &neptune.DescribeDBParameterGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, pg := range resp.DBParameterGroups {
			if pg.DBParameterGroupName != nil && *pg.DBParameterGroupName == instancePGName {
				found = true
				if pg.DBParameterGroupFamily == nil || *pg.DBParameterGroupFamily != "neptune1" {
					return fmt.Errorf("expected DBParameterGroupFamily=neptune1, got %v", pg.DBParameterGroupFamily)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created parameter group not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBParameters", func() error {
		resp, err := tc.client.DescribeDBParameters(tc.ctx, &neptune.DescribeDBParametersInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		if err != nil {
			return err
		}
		if resp.Parameters == nil {
			return fmt.Errorf("expected non-nil Parameters list")
		}
		if len(resp.Parameters) == 0 {
			return fmt.Errorf("expected at least one default parameter")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBParameterGroup", func() error {
		_, err := tc.client.ModifyDBParameterGroup(tc.ctx, &neptune.ModifyDBParameterGroupInput{
			DBParameterGroupName: aws.String(instancePGName),
			Parameters: []types.Parameter{
				{ParameterName: aws.String("neptune_query_timeout"), ParameterValue: aws.String("180000")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBParameterGroup_Verify", func() error {
		resp, err := tc.client.DescribeDBParameters(tc.ctx, &neptune.DescribeDBParametersInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		if err != nil {
			return err
		}
		found := false
		for _, p := range resp.Parameters {
			if p.ParameterName != nil && *p.ParameterName == "neptune_query_timeout" {
				found = true
				if p.ParameterValue == nil || *p.ParameterValue != "180000" {
					return fmt.Errorf("expected ParameterValue=180000, got %v", p.ParameterValue)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("expected parameter neptune_query_timeout not found after ModifyDBParameterGroup")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ResetDBParameterGroup", func() error {
		_, err := tc.client.ResetDBParameterGroup(tc.ctx, &neptune.ResetDBParameterGroupInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ResetDBParameterGroup_Verify", func() error {
		resp, err := tc.client.DescribeDBParameters(tc.ctx, &neptune.DescribeDBParametersInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		if err != nil {
			return err
		}
		for _, p := range resp.Parameters {
			if p.ParameterName != nil && *p.ParameterName == "neptune_query_timeout" {
				if p.ParameterValue != nil && *p.ParameterValue == "180000" {
					return fmt.Errorf("expected neptune_query_timeout to be reset but still has value 180000")
				}
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "CopyDBParameterGroup", func() error {
		copyPGName := fmt.Sprintf("test-pg-copy-%d", tc.ts)
		resp, err := tc.client.CopyDBParameterGroup(tc.ctx, &neptune.CopyDBParameterGroupInput{
			SourceDBParameterGroupIdentifier:  aws.String(instancePGName),
			TargetDBParameterGroupIdentifier:  aws.String(copyPGName),
			TargetDBParameterGroupDescription: aws.String("Copied instance parameter group"),
		})
		if err != nil {
			return err
		}
		if resp.DBParameterGroup == nil {
			return fmt.Errorf("expected DBParameterGroup in CopyDBParameterGroup response")
		}
		if resp.DBParameterGroup.DBParameterGroupName == nil || *resp.DBParameterGroup.DBParameterGroupName != copyPGName {
			return fmt.Errorf("expected DBParameterGroupName=%s, got %v", copyPGName, resp.DBParameterGroup.DBParameterGroupName)
		}
		_, _ = tc.client.DeleteDBParameterGroup(tc.ctx, &neptune.DeleteDBParameterGroupInput{
			DBParameterGroupName: aws.String(copyPGName),
		})
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBParameterGroup", func() error {
		_, err := tc.client.DeleteDBParameterGroup(tc.ctx, &neptune.DeleteDBParameterGroupInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBParameterGroup_VerifyDeleted", func() error {
		resp, err := tc.client.DescribeDBParameterGroups(tc.ctx, &neptune.DescribeDBParameterGroupsInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		if err != nil {
			if err := AssertErrorContains(err, "DBParameterGroupNotFoundFault"); err != nil {
				return err
			}
			return nil
		}
		if len(resp.DBParameterGroups) != 0 {
			return fmt.Errorf("expected 0 parameter groups after delete, got %d", len(resp.DBParameterGroups))
		}
		return nil
	}))

	return results
}
