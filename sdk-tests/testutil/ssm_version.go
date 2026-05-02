package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func (r *TestRunner) runSSMVersion(tc *ssmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("ssm", "PutParameter_Overwrite_IncrementsVersion", func() error {
		name := tc.uniqueName("/ver-ow")
		resp1, err := tc.putParam(name, "v1", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		defer tc.deleteParam(name)

		if resp1.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp1.Version)
		}

		resp2, err := tc.putParam(name, "v2", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Overwrite = aws.Bool(true)
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		if resp2.Version != 2 {
			return fmt.Errorf("expected version 2, got %d", resp2.Version)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameterHistory_TwoVersions", func() error {
		name := tc.uniqueName("/ver-hist")
		_, err := tc.putParam(name, "v1", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.putParam(name, "v2", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Overwrite = aws.Bool(true)
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		resp, err := tc.client.GetParameterHistory(tc.ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		if len(resp.Parameters) != 2 {
			return fmt.Errorf("expected 2 history entries, got %d", len(resp.Parameters))
		}
		if resp.Parameters[0].Version != 2 {
			return fmt.Errorf("expected first entry version 2 (newest first), got %d", resp.Parameters[0].Version)
		}
		if resp.Parameters[1].Version != 1 {
			return fmt.Errorf("expected second entry version 1, got %d", resp.Parameters[1].Version)
		}
		if resp.Parameters[0].Value == nil || *resp.Parameters[0].Value != "v2" {
			return fmt.Errorf("expected v2 value, got %q", aws.ToString(resp.Parameters[0].Value))
		}
		if resp.Parameters[1].Value == nil || *resp.Parameters[1].Value != "v1" {
			return fmt.Errorf("expected v1 value, got %q", aws.ToString(resp.Parameters[1].Value))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameterHistory_ContainsLabels", func() error {
		name := tc.uniqueName("/ver-label")
		_, err := tc.putParam(name, "v1", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.client.LabelParameterVersion(tc.ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"golden"},
		})
		if err != nil {
			return fmt.Errorf("label: %v", err)
		}

		resp, err := tc.client.GetParameterHistory(tc.ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 history entry, got %d", len(resp.Parameters))
		}
		if len(resp.Parameters[0].Labels) != 1 || resp.Parameters[0].Labels[0] != "golden" {
			return fmt.Errorf("expected label 'golden', got %v", resp.Parameters[0].Labels)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "LabelParameterVersion_GetByLabel", func() error {
		name := tc.uniqueName("/ver-getlabel")
		_, err := tc.putParam(name, "original", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		labelResp, err := tc.client.LabelParameterVersion(tc.ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"mylabel"},
		})
		if err != nil {
			return fmt.Errorf("label: %v", err)
		}
		if labelResp.ParameterVersion != 1 {
			return fmt.Errorf("expected ParameterVersion 1, got %d", labelResp.ParameterVersion)
		}

		selector := name + ":mylabel"
		resp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{
			Name: aws.String(selector),
		})
		if err != nil {
			return fmt.Errorf("get by label: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != "original" {
			return fmt.Errorf("value mismatch: got %q, want 'original'", aws.ToString(resp.Parameter.Value))
		}
		if resp.Parameter.Selector == nil || *resp.Parameter.Selector != selector {
			return fmt.Errorf("selector mismatch: got %q, want %q", aws.ToString(resp.Parameter.Selector), selector)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "LabelParameterVersion_MovesLabel", func() error {
		name := tc.uniqueName("/ver-move")
		_, err := tc.putParam(name, "v1", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.putParam(name, "v2", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Overwrite = aws.Bool(true)
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		_, err = tc.client.LabelParameterVersion(tc.ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"latest"},
		})
		if err != nil {
			return fmt.Errorf("label v1: %v", err)
		}

		_, err = tc.client.LabelParameterVersion(tc.ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(2),
			Labels:           []string{"latest"},
		})
		if err != nil {
			return fmt.Errorf("label v2: %v", err)
		}

		hist, err := tc.client.GetParameterHistory(tc.ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		for _, p := range hist.Parameters {
			if p.Version == 1 {
				if contains(p.Labels, "latest") {
					return fmt.Errorf("v1 should not have 'latest' label after move")
				}
			}
			if p.Version == 2 {
				if !contains(p.Labels, "latest") {
					return fmt.Errorf("v2 should have 'latest' label")
				}
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "UnlabelParameterVersion", func() error {
		name := tc.uniqueName("/ver-unlabel")
		_, err := tc.putParam(name, "v1", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.client.LabelParameterVersion(tc.ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"alpha", "beta"},
		})
		if err != nil {
			return fmt.Errorf("label: %v", err)
		}

		unlabelResp, err := tc.client.UnlabelParameterVersion(tc.ctx, &ssm.UnlabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"alpha"},
		})
		if err != nil {
			return fmt.Errorf("unlabel: %v", err)
		}
		if len(unlabelResp.RemovedLabels) != 1 || unlabelResp.RemovedLabels[0] != "alpha" {
			return fmt.Errorf("expected RemovedLabels ['alpha'], got %v", unlabelResp.RemovedLabels)
		}

		hist, err := tc.client.GetParameterHistory(tc.ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		if len(hist.Parameters) != 1 {
			return fmt.Errorf("expected 1 history entry, got %d", len(hist.Parameters))
		}
		if contains(hist.Parameters[0].Labels, "alpha") {
			return fmt.Errorf("'alpha' label should be removed")
		}
		if !contains(hist.Parameters[0].Labels, "beta") {
			return fmt.Errorf("'beta' label should still exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter_WithVersionSelector", func() error {
		name := tc.uniqueName("/ver-sel")
		_, err := tc.putParam(name, "version-one", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.putParam(name, "version-two", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Overwrite = aws.Bool(true)
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		resp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{
			Name: aws.String(name + ":1"),
		})
		if err != nil {
			return fmt.Errorf("get v1: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != "version-one" {
			return fmt.Errorf("expected 'version-one', got %q", aws.ToString(resp.Parameter.Value))
		}
		if resp.Parameter.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Parameter.Version)
		}
		if resp.Parameter.ARN == nil || *resp.Parameter.ARN == "" {
			return fmt.Errorf("ARN should not be empty for version selector")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_ReturnsVersionAndTier", func() error {
		name := tc.uniqueName("/resp")
		resp, err := tc.putParam(name, "check", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Tier = types.ParameterTierStandard
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}
		if resp.Tier != types.ParameterTierStandard {
			return fmt.Errorf("expected Standard tier, got %v", resp.Tier)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameterHistory_Pagination", func() error {
		name := tc.uniqueName("/ver-pag")
		for i := 0; i < 5; i++ {
			_, err := tc.putParam(name, fmt.Sprintf("v%d", i), types.ParameterTypeString, func(in *ssm.PutParameterInput) {
				in.Overwrite = aws.Bool(true)
			})
			if err != nil {
				return fmt.Errorf("put v%d: %v", i, err)
			}
		}
		defer tc.deleteParam(name)

		var allVersions []int64
		var nextToken *string
		for {
			resp, err := tc.client.GetParameterHistory(tc.ctx, &ssm.GetParameterHistoryInput{
				Name:       aws.String(name),
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				return fmt.Errorf("history: %v", err)
			}
			for _, p := range resp.Parameters {
				allVersions = append(allVersions, p.Version)
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		if len(allVersions) != 5 {
			return fmt.Errorf("expected 5 history entries, got %d", len(allVersions))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "UnlabelParameterVersion_NonExistentLabel", func() error {
		name := tc.uniqueName("/ver-unlabel-ne")
		_, err := tc.putParam(name, "v1", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.client.LabelParameterVersion(tc.ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"alpha"},
		})
		if err != nil {
			return fmt.Errorf("label: %v", err)
		}

		unlabelResp, err := tc.client.UnlabelParameterVersion(tc.ctx, &ssm.UnlabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"nonexistent"},
		})
		if err != nil {
			return fmt.Errorf("unlabel: %v", err)
		}
		if len(unlabelResp.RemovedLabels) != 0 {
			return fmt.Errorf("expected 0 removed labels, got %v", unlabelResp.RemovedLabels)
		}

		hist, err := tc.client.GetParameterHistory(tc.ctx, &ssm.GetParameterHistoryInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		if !contains(hist.Parameters[0].Labels, "alpha") {
			return fmt.Errorf("'alpha' label should still exist after removing non-existent label")
		}
		return nil
	}))

	return results
}
