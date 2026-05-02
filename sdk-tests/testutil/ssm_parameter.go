package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func (r *TestRunner) runSSMParameterCRUD(tc *ssmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("ssm", "PutParameter", func() error {
		name := tc.uniqueName("/test")
		resp, err := tc.putParam(name, "test-value", types.ParameterTypeString)
		if err != nil {
			return err
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

	results = append(results, r.RunTest("ssm", "PutParameter_WithDescription", func() error {
		name := tc.uniqueName("/desc")
		_, err := tc.putParam(name, "val", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Description = aws.String("my description")
		})
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)

		resp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != "val" {
			return fmt.Errorf("value mismatch: got %q", aws.ToString(resp.Parameter.Value))
		}
		if resp.Parameter.Name == nil || *resp.Parameter.Name != name {
			return fmt.Errorf("name mismatch: got %q", aws.ToString(resp.Parameter.Name))
		}
		if resp.Parameter.Type != types.ParameterTypeString {
			return fmt.Errorf("type mismatch: got %v", resp.Parameter.Type)
		}
		if resp.Parameter.Version != 1 {
			return fmt.Errorf("version mismatch: got %d", resp.Parameter.Version)
		}
		if resp.Parameter.ARN == nil || *resp.Parameter.ARN == "" {
			return fmt.Errorf("ARN should not be empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_StringList", func() error {
		name := tc.uniqueName("/sl")
		resp, err := tc.putParam(name, "a,b,c", types.ParameterTypeStringList)
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)
		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}

		gp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if gp.Parameter.Type != types.ParameterTypeStringList {
			return fmt.Errorf("expected StringList type, got %v", gp.Parameter.Type)
		}
		if gp.Parameter.Value == nil || *gp.Parameter.Value != "a,b,c" {
			return fmt.Errorf("value mismatch: got %q", aws.ToString(gp.Parameter.Value))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_DataType", func() error {
		name := tc.uniqueName("/dt")
		_, err := tc.putParam(name, "ami-12345", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.DataType = aws.String("aws:ec2:image")
		})
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)

		resp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter_Roundtrip", func() error {
		name := tc.uniqueName("/rt")
		value := "roundtrip-value-12345"
		_, err := tc.putParam(name, value, types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		resp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != value {
			return fmt.Errorf("value mismatch: got %q, want %q", aws.ToString(resp.Parameter.Value), value)
		}
		if resp.Parameter.Name == nil || *resp.Parameter.Name != name {
			return fmt.Errorf("name mismatch: got %q", aws.ToString(resp.Parameter.Name))
		}
		if resp.Parameter.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Parameter.Version)
		}
		if resp.Parameter.Type != types.ParameterTypeString {
			return fmt.Errorf("type mismatch: got %v", resp.Parameter.Type)
		}
		if resp.Parameter.ARN == nil || *resp.Parameter.ARN == "" {
			return fmt.Errorf("ARN should not be empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter_NonExistent", func() error {
		_, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{
			Name: aws.String("/nonexistent/param-xyz"),
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameter", func() error {
		name := tc.uniqueName("/del")
		_, err := tc.putParam(name, "delete-me", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		_, err = tc.client.DeleteParameter(tc.ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		_, err = tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return fmt.Errorf("parameter should be gone after delete: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameter_NonExistent", func() error {
		_, err := tc.client.DeleteParameter(tc.ctx, &ssm.DeleteParameterInput{
			Name: aws.String("/nonexistent/param-xyz"),
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_Duplicate_NoOverwrite", func() error {
		name := tc.uniqueName("/dup")
		_, err := tc.putParam(name, "first", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put first: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.putParam(name, "second", types.ParameterTypeString)
		if err == nil {
			return fmt.Errorf("expected error when creating duplicate without Overwrite")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_SecureString", func() error {
		name := tc.uniqueName("/sec")
		resp, err := tc.putParam(name, "secret-value", types.ParameterTypeSecureString)
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)
		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}

		gp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if gp.Parameter.Type != types.ParameterTypeSecureString {
			return fmt.Errorf("expected SecureString type, got %v", gp.Parameter.Type)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter_WithDecryption", func() error {
		name := tc.uniqueName("/dec")
		_, err := tc.putParam(name, "my-secret", types.ParameterTypeSecureString)
		if err != nil {
			return err
		}
		defer tc.deleteParam(name)

		gp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("get with decryption: %v", err)
		}
		if gp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if gp.Parameter.Type != types.ParameterTypeSecureString {
			return fmt.Errorf("type mismatch: got %v", gp.Parameter.Type)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_Overwrite_UpdatesValue", func() error {
		name := tc.uniqueName("/ow-val")
		_, err := tc.putParam(name, "original", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Description = aws.String("v1 desc")
		})
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.putParam(name, "updated", types.ParameterTypeString, func(in *ssm.PutParameterInput) {
			in.Overwrite = aws.Bool(true)
			in.Description = aws.String("v2 desc")
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		gp, err := tc.client.GetParameter(tc.ctx, &ssm.GetParameterInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if gp.Parameter.Value == nil || *gp.Parameter.Value != "updated" {
			return fmt.Errorf("value should be 'updated', got %q", aws.ToString(gp.Parameter.Value))
		}
		if gp.Parameter.Version != 2 {
			return fmt.Errorf("expected version 2, got %d", gp.Parameter.Version)
		}
		return nil
	}))

	return results
}
