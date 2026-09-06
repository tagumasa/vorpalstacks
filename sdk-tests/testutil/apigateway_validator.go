package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayValidatorTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var validatorID string
	results = append(results, r.RunTest("apigateway", "CreateRequestValidator", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, validatorID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, validatorID); err != nil {
			return err
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

		// The /name row documents replace only: add rejects.
		_, err = tc.client.UpdateRequestValidator(tc.ctx, &apigateway.UpdateRequestValidatorInput{
			RestApiId:          aws.String(tc.apiID),
			RequestValidatorId: aws.String(validatorID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/name"), Value: aws.String("nope")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for add on /name, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRequestValidators", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, validatorID); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetRequestValidator should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	return results
}
