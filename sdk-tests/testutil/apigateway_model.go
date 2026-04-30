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

func (r *TestRunner) runAPIGatewayModelTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "CreateModel", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateModel(ctx, &apigateway.CreateModelInput{
			RestApiId:   aws.String(apiID),
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
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetModel(ctx, &apigateway.GetModelInput{
			RestApiId: aws.String(apiID),
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
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.UpdateModel(ctx, &apigateway.UpdateModelInput{
			RestApiId: aws.String(apiID),
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
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetModels", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetModels(ctx, &apigateway.GetModelsInput{
			RestApiId: aws.String(apiID),
			Limit:     aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 model")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteModel", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		_, err := client.DeleteModel(ctx, &apigateway.DeleteModelInput{
			RestApiId: aws.String(apiID),
			ModelName: aws.String("UserModel"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetModel(ctx, &apigateway.GetModelInput{
			RestApiId: aws.String(apiID),
			ModelName: aws.String("UserModel"),
		})
		if err == nil {
			return fmt.Errorf("GetModel should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Model_FullLifecycle", func() error {
		mlAPI := fmt.Sprintf("MlAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(mlAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.CreateModel(ctx, &apigateway.CreateModelInput{
			RestApiId:   createResp.Id,
			Name:        aws.String("ErrorModel"),
			ContentType: aws.String("application/json"),
			Description: aws.String("Error response"),
			Schema:      aws.String(`{"type":"object","properties":{"message":{"type":"string"}}}`),
		})
		if err != nil {
			return fmt.Errorf("create model: %v", err)
		}

		getResp, err := client.GetModel(ctx, &apigateway.GetModelInput{
			RestApiId: createResp.Id,
			ModelName: aws.String("ErrorModel"),
		})
		if err != nil {
			return fmt.Errorf("get model: %v", err)
		}
		if getResp.ContentType == nil || *getResp.ContentType != "application/json" {
			return fmt.Errorf("contentType mismatch, got %v", getResp.ContentType)
		}

		_, err = client.UpdateModel(ctx, &apigateway.UpdateModelInput{
			RestApiId: createResp.Id,
			ModelName: aws.String("ErrorModel"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/schema"),
					Value: aws.String(`{"type":"object"}`),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update model: %v", err)
		}

		modelsResp, err := client.GetModels(ctx, &apigateway.GetModelsInput{
			RestApiId: createResp.Id,
			Limit:     aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("list models: %v", err)
		}
		if len(modelsResp.Items) == 0 {
			return fmt.Errorf("expected at least 1 model")
		}

		_, err = client.DeleteModel(ctx, &apigateway.DeleteModelInput{
			RestApiId: createResp.Id,
			ModelName: aws.String("ErrorModel"),
		})
		if err != nil {
			return fmt.Errorf("delete model: %v", err)
		}
		return nil
	}))

	return results
}
