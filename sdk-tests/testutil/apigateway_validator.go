package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayValidatorTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	var validatorID string
	results = append(results, r.RunTest("apigateway", "CreateRequestValidator", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateRequestValidator(ctx, &apigateway.CreateRequestValidatorInput{
			RestApiId:                 aws.String(apiID),
			Name:                      aws.String("test-validator"),
			ValidateRequestBody:       true,
			ValidateRequestParameters: true,
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("validator ID is nil")
		}
		if !resp.ValidateRequestBody {
			return fmt.Errorf("validateRequestBody should be true")
		}
		if !resp.ValidateRequestParameters {
			return fmt.Errorf("validateRequestParameters should be true")
		}
		validatorID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRequestValidator", func() error {
		if apiID == "" || validatorID == "" {
			return fmt.Errorf("API ID or validator ID not available")
		}
		resp, err := client.GetRequestValidator(ctx, &apigateway.GetRequestValidatorInput{
			RestApiId:          aws.String(apiID),
			RequestValidatorId: aws.String(validatorID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "test-validator" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRequestValidator", func() error {
		if apiID == "" || validatorID == "" {
			return fmt.Errorf("API ID or validator ID not available")
		}
		resp, err := client.UpdateRequestValidator(ctx, &apigateway.UpdateRequestValidatorInput{
			RestApiId:          aws.String(apiID),
			RequestValidatorId: aws.String(validatorID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/name"),
					Value: aws.String("updated-validator"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "updated-validator" {
			return fmt.Errorf("name not updated, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRequestValidators", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetRequestValidators(ctx, &apigateway.GetRequestValidatorsInput{
			RestApiId: aws.String(apiID),
			Limit:     aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 validator")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRequestValidator", func() error {
		if apiID == "" || validatorID == "" {
			return fmt.Errorf("API ID or validator ID not available")
		}
		_, err := client.DeleteRequestValidator(ctx, &apigateway.DeleteRequestValidatorInput{
			RestApiId:          aws.String(apiID),
			RequestValidatorId: aws.String(validatorID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetRequestValidator(ctx, &apigateway.GetRequestValidatorInput{
			RestApiId:          aws.String(apiID),
			RequestValidatorId: aws.String(validatorID),
		})
		if err == nil {
			return fmt.Errorf("GetRequestValidator should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "RequestValidator_FullLifecycle", func() error {
		rvAPI := fmt.Sprintf("RvAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(rvAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.CreateRequestValidator(ctx, &apigateway.CreateRequestValidatorInput{
			RestApiId:                 createResp.Id,
			Name:                      aws.String("body-only"),
			ValidateRequestBody:       true,
			ValidateRequestParameters: false,
		})
		if err != nil {
			return fmt.Errorf("create body validator: %v", err)
		}

		_, err = client.CreateRequestValidator(ctx, &apigateway.CreateRequestValidatorInput{
			RestApiId:                 createResp.Id,
			Name:                      aws.String("params-only"),
			ValidateRequestBody:       false,
			ValidateRequestParameters: true,
		})
		if err != nil {
			return fmt.Errorf("create params validator: %v", err)
		}

		rvListResp, err := client.GetRequestValidators(ctx, &apigateway.GetRequestValidatorsInput{
			RestApiId: createResp.Id,
			Limit:     aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("list validators: %v", err)
		}
		if len(rvListResp.Items) < 2 {
			return fmt.Errorf("expected at least 2 validators, got %d", len(rvListResp.Items))
		}
		return nil
	}))

	return results
}
