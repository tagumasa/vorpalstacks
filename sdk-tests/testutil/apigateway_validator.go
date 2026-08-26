package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayValidatorTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var validatorID string
	results = append(results, r.RunTest("apigateway", "CreateRequestValidator", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.CreateRequestValidator(tc.ctx, &apigateway.CreateRequestValidatorInput{
			RestApiId:                 aws.String(tc.apiID),
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
		if tc.apiID == "" || validatorID == "" {
			return fmt.Errorf("API ID or validator ID not available")
		}
		resp, err := tc.client.GetRequestValidator(tc.ctx, &apigateway.GetRequestValidatorInput{
			RestApiId:          aws.String(tc.apiID),
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
		if tc.apiID == "" || validatorID == "" {
			return fmt.Errorf("API ID or validator ID not available")
		}
		resp, err := tc.client.UpdateRequestValidator(tc.ctx, &apigateway.UpdateRequestValidatorInput{
			RestApiId:          aws.String(tc.apiID),
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
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}

		// Exercise both single-flag combinations alongside the both-flags
		// validator created above, then require the list to see them all.
		_, err := tc.client.CreateRequestValidator(tc.ctx, &apigateway.CreateRequestValidatorInput{
			RestApiId:                 aws.String(tc.apiID),
			Name:                      aws.String("body-only"),
			ValidateRequestBody:       true,
			ValidateRequestParameters: false,
		})
		if err != nil {
			return fmt.Errorf("create body validator: %v", err)
		}

		_, err = tc.client.CreateRequestValidator(tc.ctx, &apigateway.CreateRequestValidatorInput{
			RestApiId:                 aws.String(tc.apiID),
			Name:                      aws.String("params-only"),
			ValidateRequestBody:       false,
			ValidateRequestParameters: true,
		})
		if err != nil {
			return fmt.Errorf("create params validator: %v", err)
		}

		items, err := tc.allRequestValidators(tc.apiID)
		if err != nil {
			return err
		}
		if len(items) < 2 {
			return fmt.Errorf("expected at least 2 validators, got %d", len(items))
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRequestValidator", func() error {
		if tc.apiID == "" || validatorID == "" {
			return fmt.Errorf("API ID or validator ID not available")
		}
		_, err := tc.client.DeleteRequestValidator(tc.ctx, &apigateway.DeleteRequestValidatorInput{
			RestApiId:          aws.String(tc.apiID),
			RequestValidatorId: aws.String(validatorID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetRequestValidator(tc.ctx, &apigateway.GetRequestValidatorInput{
			RestApiId:          aws.String(tc.apiID),
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

	return results
}
