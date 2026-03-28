package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSSMTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "ssm",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := ssm.NewFromConfig(cfg)
	ctx := context.Background()

	parameterName := fmt.Sprintf("/test/param-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("ssm", "PutParameter", func() error {
		resp, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(parameterName),
			Value: aws.String("test-value"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter", func() error {
		resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String(parameterName),
		})
		if err != nil {
			return err
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameters", func() error {
		resp, err := client.GetParameters(ctx, &ssm.GetParametersInput{
			Names: []string{parameterName},
		})
		if err != nil {
			return err
		}
		if resp.Parameters == nil {
			return fmt.Errorf("parameters list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParametersByPath", func() error {
		resp, err := client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:      aws.String("/test"),
			Recursive: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Parameters == nil {
			return fmt.Errorf("parameters list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters", func() error {
		resp, err := client.DescribeParameters(ctx, &ssm.DescribeParametersInput{
			Filters: []types.ParametersFilter{
				{
					Key:    types.ParametersFilterKeyName,
					Values: []string{parameterName},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Parameters == nil {
			return fmt.Errorf("parameters list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameter", func() error {
		resp, err := client.DeleteParameter(ctx, &ssm.DeleteParameterInput{
			Name: aws.String(parameterName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("ssm", "GetParameter_NonExistent", func() error {
		_, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String("/nonexistent/param-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent parameter")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameter_NonExistent", func() error {
		_, err := client.DeleteParameter(ctx, &ssm.DeleteParameterInput{
			Name: aws.String("/nonexistent/param-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent parameter")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_GetParameter_Roundtrip", func() error {
		rtName := fmt.Sprintf("/rt/param-%d", time.Now().UnixNano())
		rtValue := "roundtrip-value-12345"
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(rtName),
			Value: aws.String(rtValue),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(rtName)})

		resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String(rtName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != rtValue {
			return fmt.Errorf("value mismatch: got %q, want %q", aws.ToString(resp.Parameter.Value), rtValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameters_InvalidNames", func() error {
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String("/valid/param-test"),
			Value: aws.String("valid"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String("/valid/param-test")})

		resp, err := client.GetParameters(ctx, &ssm.GetParametersInput{
			Names: []string{"/valid/param-test", "/nonexistent/param-xyz"},
		})
		if err != nil {
			return fmt.Errorf("get params: %v", err)
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 valid parameter, got %d", len(resp.Parameters))
		}
		if len(resp.InvalidParameters) != 1 {
			return fmt.Errorf("expected 1 invalid parameter, got %d", len(resp.InvalidParameters))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters_ContainsCreated", func() error {
		dpName := fmt.Sprintf("/dp/param-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:        aws.String(dpName),
			Value:       aws.String("desc-test"),
			Type:        types.ParameterTypeString,
			Description: aws.String("Test description for search"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(dpName)})

		resp, err := client.DescribeParameters(ctx, &ssm.DescribeParametersInput{
			Filters: []types.ParametersFilter{
				{
					Key:    types.ParametersFilterKeyName,
					Values: []string{dpName},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 parameter, got %d", len(resp.Parameters))
		}
		if resp.Parameters[0].Description == nil || *resp.Parameters[0].Description != "Test description for search" {
			return fmt.Errorf("description mismatch")
		}
		return nil
	}))

	return results
}
