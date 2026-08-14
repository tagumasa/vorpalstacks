package testutil

import (
	"fmt"
	"strings"

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

	// The AWS-documented DataType value set includes aws:ssm:integration.
	results = append(results, r.RunTest("ssm", "PutParameter_DataType_Integration", func() error {
		name := tc.uniqueName("/dt-int")
		resp, err := tc.putParam(name, "integration-val", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.DataType = aws.String("aws:ssm:integration")
		})
		if err != nil {
			return fmt.Errorf("put with aws:ssm:integration: %v", err)
		}
		defer tc.deleteParam(name)
		if resp.Tier == "" {
			return fmt.Errorf("empty tier in response")
		}

		gp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if gp.Parameter.DataType == nil || *gp.Parameter.DataType != "aws:ssm:integration" {
			return fmt.Errorf("DataType mismatch: %v", gp.Parameter.DataType)
		}
		return nil
	}))

	// Smithy marks PutParameter Value as required; empty values are rejected.
	results = append(results, r.RunTest("ssm", "PutParameter_EmptyValue_Rejected", func() error {
		name := tc.uniqueName("/empty-val")
		defer tc.deleteParam(name)
		_, err := tc.putParam(name, "", types.ParameterTypeString)
		if err == nil {
			return fmt.Errorf("expected error for empty Value")
		}
		return nil
	}))

	// The documented 2048-character maximum includes 1037 reserved
	// characters; caller-specified names are capped at 1011.
	results = append(results, r.RunTest("ssm", "PutParameter_NameTooLong_Rejected", func() error {
		defer tc.deleteParam("long-name")
		_, err := tc.putParam(strings.Repeat("a", 1012), "val", types.ParameterTypeString)
		if err == nil {
			return fmt.Errorf("expected error for 1012-character name")
		}
		return nil
	}))

	// Parameter hierarchies are limited to fifteen levels.
	results = append(results, r.RunTest("ssm", "PutParameter_HierarchyTooDeep_Rejected", func() error {
		segments := make([]string, 16)
		for i := range segments {
			segments[i] = "lvl"
		}
		name := "/" + strings.Join(segments, "/")
		defer tc.deleteParam(name)
		_, err := tc.putParam(name, "val", types.ParameterTypeString)
		if err == nil {
			return fmt.Errorf("expected error for 16-level hierarchy")
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
