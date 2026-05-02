package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneClusterParamGroupTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "CreateDBClusterParameterGroup", func() error {
		resp, err := tc.client.CreateDBClusterParameterGroup(tc.ctx, &neptune.CreateDBClusterParameterGroupInput{
			DBClusterParameterGroupName: aws.String(tc.paramGroupName),
			DBParameterGroupFamily:      aws.String("neptune1"),
			Description:                 aws.String("Test cluster parameter group"),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterParameterGroup == nil {
			return fmt.Errorf("expected DBClusterParameterGroup in response")
		}
		pg := resp.DBClusterParameterGroup
		if pg.DBClusterParameterGroupName == nil || *pg.DBClusterParameterGroupName != tc.paramGroupName {
			return fmt.Errorf("expected name=%s, got %v", tc.paramGroupName, pg.DBClusterParameterGroupName)
		}
		if pg.DBParameterGroupFamily == nil || *pg.DBParameterGroupFamily != "neptune1" {
			return fmt.Errorf("expected family=neptune1, got %v", pg.DBParameterGroupFamily)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterParameterGroups", func() error {
		resp, err := tc.client.DescribeDBClusterParameterGroups(tc.ctx, &neptune.DescribeDBClusterParameterGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, pg := range resp.DBClusterParameterGroups {
			if pg.DBClusterParameterGroupName != nil && *pg.DBClusterParameterGroupName == tc.paramGroupName {
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

	results = append(results, r.RunTest("neptune", "DescribeDBClusterParameterGroups_FilterByName", func() error {
		resp, err := tc.client.DescribeDBClusterParameterGroups(tc.ctx, &neptune.DescribeDBClusterParameterGroupsInput{
			DBClusterParameterGroupName: aws.String(tc.paramGroupName),
		})
		if err != nil {
			return err
		}
		if len(resp.DBClusterParameterGroups) != 1 {
			return fmt.Errorf("expected 1 parameter group, got %d", len(resp.DBClusterParameterGroups))
		}
		pg := resp.DBClusterParameterGroups[0]
		if pg.Description == nil || *pg.Description != "Test cluster parameter group" {
			return fmt.Errorf("expected Description='Test cluster parameter group', got %v", pg.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterParameters", func() error {
		resp, err := tc.client.DescribeDBClusterParameters(tc.ctx, &neptune.DescribeDBClusterParametersInput{
			DBClusterParameterGroupName: aws.String(tc.paramGroupName),
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

	results = append(results, r.RunTest("neptune", "DescribeEngineDefaultClusterParameters", func() error {
		resp, err := tc.client.DescribeEngineDefaultClusterParameters(tc.ctx, &neptune.DescribeEngineDefaultClusterParametersInput{
			DBParameterGroupFamily: aws.String("neptune1"),
		})
		if err != nil {
			return err
		}
		if resp.EngineDefaults == nil {
			return fmt.Errorf("expected non-nil EngineDefaults")
		}
		if resp.EngineDefaults.Parameters == nil || len(resp.EngineDefaults.Parameters) == 0 {
			return fmt.Errorf("expected non-empty default parameters list")
		}
		return nil
	}))

	return results
}
