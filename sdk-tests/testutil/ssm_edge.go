package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func (r *TestRunner) runSSMEdge(tc *ssmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("ssm", "PutParameter_ReservedName", func() error {
		_, err := tc.putParam("aws:test/reserved", "val", types.ParameterTypeString)
		if err == nil {
			return fmt.Errorf("expected error for reserved 'aws:' parameter name prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_AdvancedTier", func() error {
		name := tc.uniqueName("/adv")
		resp, err := tc.putParam(name, "advanced-val", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Tier = types.ParameterTierAdvanced
		})
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)
		if resp.Tier != types.ParameterTierAdvanced {
			return fmt.Errorf("expected Advanced tier, got %v", resp.Tier)
		}

		gp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if gp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_AllowedPattern", func() error {
		name := tc.uniqueName("/pat")
		_, err := tc.putParam(name, "123", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.AllowedPattern = aws.String("^[0-9]+$")
		})
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)

		desc, err := tc.client.DescribeParameters(tc.ctx, &ssm.DescribeParametersInput{
			Filters: []types.ParametersFilter{
				{Key: types.ParametersFilterKeyName, Values: []string{name}},
			},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		found := false
		for _, p := range desc.Parameters {
			if p.Name != nil && *p.Name == name {
				found = true
				if p.AllowedPattern == nil || *p.AllowedPattern != "^[0-9]+$" {
					return fmt.Errorf("allowedPattern mismatch: got %q", aws.ToString(p.AllowedPattern))
				}
			}
		}
		if !found {
			return fmt.Errorf("parameter not found in describe")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameterHistory_NonExistent", func() error {
		_, err := tc.client.GetParameterHistory(tc.ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String("/nonexistent/history-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent parameter history")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "LabelParameterVersion_NonExistentParam", func() error {
		_, err := tc.client.LabelParameterVersion(tc.ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String("/nonexistent/label-xyz"),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"test"},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent parameter")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParametersByPath_Pagination", func() error {
		ppTs := tc.ts + "-pp"
		var ppParams []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("/pathpag/%s/p-%d", ppTs, i)
			_, err := tc.putParam(name, "pv", types.ParameterTypeString)
			if err != nil {
				for _, pn := range ppParams {
					tc.deleteParam(pn)
				}
				return fmt.Errorf("put %s: %v", name, err)
			}
			ppParams = append(ppParams, name)
		}

		var found int
		var nextToken *string
		for {
			resp, err := tc.client.GetParametersByPath(tc.ctx, &ssm.GetParametersByPathInput{
				Path:       aws.String("/pathpag/" + ppTs),
				Recursive:  aws.Bool(true),
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, pn := range ppParams {
					tc.deleteParam(pn)
				}
				return fmt.Errorf("get by path: %v", err)
			}
			for _, p := range resp.Parameters {
				if p.Name != nil {
					for _, pp := range ppParams {
						if *p.Name == pp {
							found++
						}
					}
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}

		for _, pn := range ppParams {
			tc.deleteParam(pn)
		}
		if found != 5 {
			return fmt.Errorf("expected 5 parameters, found %d", found)
		}
		return nil
	}))

	return results
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
