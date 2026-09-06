package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayModelTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "CreateModel", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.CreateModel(tc.ctx, &apigateway.CreateModelInput{
			RestApiId:   aws.String(tc.apiID),
			Name:        aws.String("UserModel"),
			ContentType: aws.String("application/json"),
			Description: aws.String("User model"),
			Schema:      aws.String(`{"type":"object"}`),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("model ID is nil")
		}
		if resp.Name == nil || *resp.Name != "UserModel" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.ContentType == nil || *resp.ContentType != "application/json" {
			return fmt.Errorf("contentType mismatch, got %v", resp.ContentType)
		}
		if resp.Description == nil || *resp.Description != "User model" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		if resp.Schema == nil || *resp.Schema != `{"type":"object"}` {
			return fmt.Errorf("schema mismatch, got %v", resp.Schema)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetModel", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.GetModel(tc.ctx, &apigateway.GetModelInput{
			RestApiId: aws.String(tc.apiID),
			ModelName: aws.String("UserModel"),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "UserModel" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.ContentType == nil || *resp.ContentType != "application/json" {
			return fmt.Errorf("contentType mismatch, got %v", resp.ContentType)
		}
		if resp.Description == nil || *resp.Description != "User model" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		if resp.Schema == nil || *resp.Schema != `{"type":"object"}` {
			return fmt.Errorf("schema mismatch, got %v", resp.Schema)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateModel", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.UpdateModel(tc.ctx, &apigateway.UpdateModelInput{
			RestApiId: aws.String(tc.apiID),
			ModelName: aws.String("UserModel"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated model"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated model" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}

		// The /description row documents replace only: remove rejects.
		_, err = tc.client.UpdateModel(tc.ctx, &apigateway.UpdateModelInput{
			RestApiId: aws.String(tc.apiID),
			ModelName: aws.String("UserModel"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/description")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for remove on /description, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetModels", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		items, err := tc.allModels(tc.apiID)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return fmt.Errorf("expected at least 1 model")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteModel", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		_, err := tc.client.DeleteModel(tc.ctx, &apigateway.DeleteModelInput{
			RestApiId: aws.String(tc.apiID),
			ModelName: aws.String("UserModel"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetModel(tc.ctx, &apigateway.GetModelInput{
			RestApiId: aws.String(tc.apiID),
			ModelName: aws.String("UserModel"),
		})
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetModel should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	return results
}
