package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func (r *TestRunner) runSSMParameterList(tc *ssmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("ssm", "GetParameters", func() error {
		name := tc.uniqueName("/gp")
		_, err := tc.putParam(name, "multi-val", types.ParameterTypeString)
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)

		resp, err := tc.client.GetParameters(tc.ctx, &ssm.GetParametersInput{
			Names: []string{name},
		})
		if err != nil {
			return err
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 parameter, got %d", len(resp.Parameters))
		}
		if resp.Parameters[0].Name == nil || *resp.Parameters[0].Name != name {
			return fmt.Errorf("name mismatch: got %q", aws.ToString(resp.Parameters[0].Name))
		}
		if resp.Parameters[0].Value == nil || *resp.Parameters[0].Value != "multi-val" {
			return fmt.Errorf("value mismatch: got %q", aws.ToString(resp.Parameters[0].Value))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameters_InvalidNames", func() error {
		name := tc.uniqueName("/valid")
		_, err := tc.putParam(name, "valid", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		resp, err := tc.client.GetParameters(tc.ctx, &ssm.GetParametersInput{
			Names: []string{name, "/nonexistent/param-xyz"},
		})
		if err != nil {
			return fmt.Errorf("get params: %v", err)
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 valid parameter, got %d", len(resp.Parameters))
		}
		if resp.Parameters[0].Name == nil || *resp.Parameters[0].Name != name {
			return fmt.Errorf("returned wrong parameter: got %q", aws.ToString(resp.Parameters[0].Name))
		}
		if len(resp.InvalidParameters) != 1 {
			return fmt.Errorf("expected 1 invalid parameter, got %d", len(resp.InvalidParameters))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParametersByPath", func() error {
		pathTs := tc.ts + "-pp"
		pathName := fmt.Sprintf("/path-test/%s/param", pathTs)
		_, err := tc.putParam(pathName, "path-val", types.ParameterTypeString)
		if err != nil {
			return err
		}
		defer tc.deleteParam(pathName)

		var found bool
		var nextToken *string
		for {
			resp, err := tc.client.GetParametersByPath(tc.ctx, &ssm.GetParametersByPathInput{
				Path:      aws.String("/path-test/" + pathTs),
				Recursive: aws.Bool(true),
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("get by path: %v", err)
			}
			for _, p := range resp.Parameters {
				if p.Name != nil && *p.Name == pathName {
					found = true
					if p.Value == nil || *p.Value != "path-val" {
						return fmt.Errorf("value mismatch: got %q", aws.ToString(p.Value))
					}
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		if !found {
			return fmt.Errorf("parameter %s not found by path", pathName)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParametersByPath_NonRecursive", func() error {
		nrTs := tc.ts + "-nr"
		base := fmt.Sprintf("/nr/%s/param", nrTs)
		_, err := tc.putParam(base, "direct", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put direct: %v", err)
		}
		defer tc.deleteParam(base)

		_, err = tc.putParam(base+"/child", "child", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put child: %v", err)
		}
		defer tc.deleteParam(base + "/child")

		var nextToken *string
		for {
			resp, err := tc.client.GetParametersByPath(tc.ctx, &ssm.GetParametersByPathInput{
				Path:      aws.String("/nr/" + nrTs),
				Recursive: aws.Bool(false),
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("get by path: %v", err)
			}
			for _, p := range resp.Parameters {
				if strings.Contains(aws.ToString(p.Name), "/child") {
					return fmt.Errorf("non-recursive should not return children, got %q", aws.ToString(p.Name))
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters", func() error {
		dpName := tc.uniqueName("/dp")
		_, err := tc.putParam(dpName, "desc-test", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Description = aws.String("Test description for search")
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(dpName)

		var found bool
		var nextToken *string
		for {
			resp, err := tc.client.DescribeParameters(tc.ctx, &ssm.DescribeParametersInput{
				Filters: []types.ParametersFilter{
					{Key: types.ParametersFilterKeyName, Values: []string{dpName}},
				},
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe: %v", err)
			}
			for _, p := range resp.Parameters {
				if p.Name != nil && *p.Name == dpName {
					found = true
					if p.Description == nil || *p.Description != "Test description for search" {
						return fmt.Errorf("description mismatch: got %q", aws.ToString(p.Description))
					}
					if p.Type != types.ParameterTypeString {
						return fmt.Errorf("type mismatch: got %v", p.Type)
					}
					if p.Version != 1 {
						return fmt.Errorf("version mismatch: got %d", p.Version)
					}
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		if !found {
			return fmt.Errorf("parameter %s not found in DescribeParameters", dpName)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters_TypeFilter", func() error {
		strName := tc.uniqueName("/tf-str")
		slName := tc.uniqueName("/tf-sl")
		_, err := tc.putParam(strName, "x", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put string: %v", err)
		}
		defer tc.deleteParam(strName)
		_, err = tc.putParam(slName, "a,b,c", types.ParameterTypeStringList)
		if err != nil {
			return fmt.Errorf("put stringlist: %v", err)
		}
		defer tc.deleteParam(slName)

		var nextToken *string
		for {
			resp, err := tc.client.DescribeParameters(tc.ctx, &ssm.DescribeParametersInput{
				Filters: []types.ParametersFilter{
					{Key: types.ParametersFilterKeyType, Values: []string{"String"}},
				},
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe: %v", err)
			}
			for _, p := range resp.Parameters {
				if p.Type != types.ParameterTypeString {
					return fmt.Errorf("type filter returned non-String type: %v", p.Type)
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters_Pagination", func() error {
		pgTs := tc.ts + "-pg"
		var pgParams []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("/pag/%s/param-%d", pgTs, i)
			_, err := tc.putParam(name, "pagval", types.ParameterTypeString)
			if err != nil {
				for _, pn := range pgParams {
					tc.deleteParam(pn)
				}
				return fmt.Errorf("put parameter %s: %v", name, err)
			}
			pgParams = append(pgParams, name)
		}

		var allParams []string
		var nextToken *string
		for {
			resp, err := tc.client.DescribeParameters(tc.ctx, &ssm.DescribeParametersInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, pn := range pgParams {
					tc.deleteParam(pn)
				}
				return fmt.Errorf("describe parameters page: %v", err)
			}
			for _, p := range resp.Parameters {
				if p.Name != nil && strings.Contains(*p.Name, "/pag/"+pgTs) {
					allParams = append(allParams, *p.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, pn := range pgParams {
			tc.deleteParam(pn)
		}
		if len(allParams) != 5 {
			return fmt.Errorf("expected 5 paginated parameters, got %d", len(allParams))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameters_WithDecryption", func() error {
		name := tc.uniqueName("/gp-dec")
		_, err := tc.putParam(name, "plaintext", types.ParameterTypeString)
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)

		resp, err := tc.client.GetParameters(tc.ctx, &ssm.GetParametersInput{
			Names:          []string{name},
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 parameter, got %d", len(resp.Parameters))
		}
		if resp.Parameters[0].Value == nil || *resp.Parameters[0].Value != "plaintext" {
			return fmt.Errorf("value mismatch: got %q", aws.ToString(resp.Parameters[0].Value))
		}
		return nil
	}))

	return results
}
