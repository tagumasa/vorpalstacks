package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunAPIGatewayTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "apigateway",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := apigateway.NewFromConfig(cfg)
	ctx := context.Background()

	apiName := fmt.Sprintf("TestAPI-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("apigateway", "CreateRestApi", func() error {
		resp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(apiName),
			Description: aws.String("Test API"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("API ID is nil")
		}
		return nil
	}))

	var apiID string
	results = append(results, r.RunTest("apigateway", "GetRestApis", func() error {
		resp, err := client.GetRestApis(ctx, &apigateway.GetRestApisInput{
			Limit: aws.Int32(500),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil {
			return fmt.Errorf("items list is nil")
		}
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == apiName {
				apiID = *item.Id
				break
			}
		}
		if apiID == "" {
			return fmt.Errorf("API not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApi", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != apiName {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: aws.String(apiID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("Updated API"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	var resourceID string
	results = append(results, r.RunTest("apigateway", "CreateResource", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(apiID),
			ParentId:  aws.String(apiID),
			PathPart:  aws.String("test"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("resource ID is nil")
		}
		resourceID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetResource", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetResource(ctx, &apigateway.GetResourceInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != resourceID {
			return fmt.Errorf("resource ID mismatch, got %v", resp.Id)
		}
		if resp.Path == nil || *resp.Path != "/test" {
			return fmt.Errorf("path mismatch, got %v", resp.Path)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetResources", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetResources(ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil {
			return fmt.Errorf("items list is nil")
		}
		if len(resp.Items) < 2 {
			return fmt.Errorf("expected at least 2 resources, got %d", len(resp.Items))
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateResource", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateResource(ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/pathPart"),
					Value: aws.String("items"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Path == nil || *resp.Path != "/items" {
			return fmt.Errorf("path not updated, got %v", resp.Path)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(apiID),
			ResourceId:        aws.String(resourceID),
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
			ApiKeyRequired:    false,
		})
		if err != nil {
			return err
		}
		if resp.HttpMethod == nil || *resp.HttpMethod != "GET" {
			return fmt.Errorf("httpMethod mismatch, got %v", resp.HttpMethod)
		}
		if resp.AuthorizationType == nil || *resp.AuthorizationType != "NONE" {
			return fmt.Errorf("authorizationType mismatch, got %v", resp.AuthorizationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetMethod(ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return err
		}
		if resp.HttpMethod == nil || *resp.HttpMethod != "GET" {
			return fmt.Errorf("httpMethod mismatch, got %v", resp.HttpMethod)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateMethod(ctx, &apigateway.UpdateMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/authorizationType"),
					Value: aws.String("AWS_IAM"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.AuthorizationType == nil || *resp.AuthorizationType != "AWS_IAM" {
			return fmt.Errorf("authorizationType not updated, got %v", resp.AuthorizationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			Type:       types.IntegrationTypeMock,
			RequestTemplates: map[string]string{
				"application/json": "{\"statusCode\": 200}",
			},
		})
		if err != nil {
			return err
		}
		if resp.Type != types.IntegrationTypeMock {
			return fmt.Errorf("type mismatch, got %v", resp.Type)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return err
		}
		if resp.Type != types.IntegrationTypeMock {
			return fmt.Errorf("type mismatch, got %v", resp.Type)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateIntegration(ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/timeoutInMillis"),
					Value: aws.String("5000"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.TimeoutInMillis != 5000 {
			return fmt.Errorf("timeoutInMillis not updated, got %d", resp.TimeoutInMillis)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutIntegrationResponse(ctx, &apigateway.PutIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			ResponseTemplates: map[string]string{
				"application/json": "{\"message\": \"ok\"}",
			},
			SelectionPattern: aws.String("2\\d{2}"),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetIntegrationResponse(ctx, &apigateway.GetIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateIntegrationResponse(ctx, &apigateway.UpdateIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/selectionPattern"),
					Value: aws.String("ok"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.SelectionPattern == nil || *resp.SelectionPattern != "ok" {
			return fmt.Errorf("selectionPattern not updated, got %v", resp.SelectionPattern)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutMethodResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutMethodResponse(ctx, &apigateway.PutMethodResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			ResponseModels: map[string]string{
				"application/json": "Empty",
			},
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetMethodResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetMethodResponse(ctx, &apigateway.GetMethodResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteMethodResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteMethodResponse(ctx, &apigateway.DeleteMethodResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "DeleteIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteIntegrationResponse(ctx, &apigateway.DeleteIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "DeleteIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteIntegration(ctx, &apigateway.DeleteIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "DeleteMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteMethod(ctx, &apigateway.DeleteMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "DeleteResource", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteResource(ctx, &apigateway.DeleteResourceInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
		})
		return err
	}))

	var deploymentID string
	results = append(results, r.RunTest("apigateway", "CreateDeployment", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId:   aws.String(apiID),
			Description: aws.String("test deployment"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("deployment ID is nil")
		}
		deploymentID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDeployment", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := client.GetDeployment(ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != deploymentID {
			return fmt.Errorf("deployment ID mismatch, got %v", resp.Id)
		}
		if resp.CreatedDate == nil {
			return fmt.Errorf("createdDate is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateDeployment", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := client.UpdateDeployment(ctx, &apigateway.UpdateDeploymentInput{
			RestApiId:    aws.String(apiID),
			DeploymentId: aws.String(deploymentID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated deployment"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated deployment" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDeployments", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetDeployments(ctx, &apigateway.GetDeploymentsInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil || len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 deployment")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateStage", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(apiID),
			StageName:    aws.String("test"),
			DeploymentId: aws.String(deploymentID),
			Description:  aws.String("test stage"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStage", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return err
		}
		if resp.StageName == nil || *resp.StageName != "test" {
			return fmt.Errorf("stage name mismatch, got %v", resp.StageName)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStages", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetStages(ctx, &apigateway.GetStagesInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Item == nil || len(resp.Item) == 0 {
			return fmt.Errorf("expected at least 1 stage")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.UpdateStage(ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated stage"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated stage" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteStage", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		_, err := client.DeleteStage(ctx, &apigateway.DeleteStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("test"),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "DeleteDeployment", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		_, err := client.DeleteDeployment(ctx, &apigateway.DeleteDeploymentInput{
			RestApiId:    aws.String(apiID),
			DeploymentId: aws.String(deploymentID),
		})
		return err
	}))

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
		if resp.Items == nil || len(resp.Items) == 0 {
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
		return err
	}))

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
		if resp.Items == nil || len(resp.Items) == 0 {
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
		return err
	}))

	var authorizerID string
	results = append(results, r.RunTest("apigateway", "CreateAuthorizer", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateAuthorizer(ctx, &apigateway.CreateAuthorizerInput{
			RestApiId:                    aws.String(apiID),
			Name:                         aws.String("test-authorizer"),
			Type:                         types.AuthorizerTypeToken,
			AuthorizerUri:                aws.String("https://example.com/auth"),
			IdentitySource:               aws.String("method.request.header.Authorization"),
			AuthorizerResultTtlInSeconds: aws.Int32(300),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("authorizer ID is nil")
		}
		if resp.Type != types.AuthorizerTypeToken {
			return fmt.Errorf("type mismatch, got %v", resp.Type)
		}
		authorizerID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetAuthorizer", func() error {
		if apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		resp, err := client.GetAuthorizer(ctx, &apigateway.GetAuthorizerInput{
			RestApiId:    aws.String(apiID),
			AuthorizerId: aws.String(authorizerID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "test-authorizer" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateAuthorizer", func() error {
		if apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		resp, err := client.UpdateAuthorizer(ctx, &apigateway.UpdateAuthorizerInput{
			RestApiId:    aws.String(apiID),
			AuthorizerId: aws.String(authorizerID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/name"),
					Value: aws.String("updated-authorizer"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "updated-authorizer" {
			return fmt.Errorf("name not updated, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetAuthorizers", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetAuthorizers(ctx, &apigateway.GetAuthorizersInput{
			RestApiId: aws.String(apiID),
			Limit:     aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil || len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 authorizer")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "TestInvokeAuthorizer", func() error {
		if apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		resp, err := client.TestInvokeAuthorizer(ctx, &apigateway.TestInvokeAuthorizerInput{
			RestApiId:    aws.String(apiID),
			AuthorizerId: aws.String(authorizerID),
			Headers: map[string]string{
				"Authorization": "Bearer test-token",
			},
		})
		if err != nil {
			return err
		}
		if resp.ClientStatus != 200 {
			return fmt.Errorf("expected clientStatus 200, got %d", resp.ClientStatus)
		}
		if resp.Policy == nil {
			return fmt.Errorf("policy is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteAuthorizer", func() error {
		if apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		_, err := client.DeleteAuthorizer(ctx, &apigateway.DeleteAuthorizerInput{
			RestApiId:    aws.String(apiID),
			AuthorizerId: aws.String(authorizerID),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "TestInvokeMethod", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(apiID),
			ParentId:  aws.String(apiID),
			PathPart:  aws.String("mock"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}
		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(apiID),
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("POST"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}
		_, err = client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("POST"),
			Type:       types.IntegrationTypeMock,
			RequestTemplates: map[string]string{
				"application/json": "{\"statusCode\": 200}",
			},
		})
		if err != nil {
			return fmt.Errorf("put integration: %v", err)
		}
		resp, err := client.TestInvokeMethod(ctx, &apigateway.TestInvokeMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("POST"),
			Body:       aws.String(`{"test": "data"}`),
		})
		if err != nil {
			return err
		}
		if resp.Status != 200 {
			return fmt.Errorf("expected status 200, got %d", resp.Status)
		}
		if resp.Log == nil {
			return fmt.Errorf("log is nil")
		}
		return nil
	}))

	var apiKeyValue string
	results = append(results, r.RunTest("apigateway", "CreateApiKey", func() error {
		resp, err := client.CreateApiKey(ctx, &apigateway.CreateApiKeyInput{
			Name:        aws.String("test-api-key"),
			Description: aws.String("Test API key"),
			Enabled:     true,
			Tags: map[string]string{
				"env": "test",
			},
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("api key ID is nil")
		}
		if resp.Value == nil {
			return fmt.Errorf("api key value is nil")
		}
		if !resp.Enabled {
			return fmt.Errorf("expected enabled=true")
		}
		apiKeyValue = *resp.Value
		return nil
	}))

	var apiKeyID string
	results = append(results, r.RunTest("apigateway", "GetApiKeys", func() error {
		resp, err := client.GetApiKeys(ctx, &apigateway.GetApiKeysInput{
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil || len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 api key")
		}
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == "test-api-key" {
				apiKeyID = *item.Id
				break
			}
		}
		if apiKeyID == "" {
			return fmt.Errorf("test-api-key not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetApiKey", func() error {
		if apiKeyID == "" {
			return fmt.Errorf("api key ID not available")
		}
		resp, err := client.GetApiKey(ctx, &apigateway.GetApiKeyInput{
			ApiKey:       aws.String(apiKeyID),
			IncludeValue: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "test-api-key" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Value == nil || *resp.Value != apiKeyValue {
			return fmt.Errorf("value mismatch, got %v", resp.Value)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateApiKey", func() error {
		if apiKeyID == "" {
			return fmt.Errorf("api key ID not available")
		}
		resp, err := client.UpdateApiKey(ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/name"),
					Value: aws.String("updated-api-key"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "updated-api-key" {
			return fmt.Errorf("name not updated, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteApiKey", func() error {
		if apiKeyID == "" {
			return fmt.Errorf("api key ID not available")
		}
		_, err := client.DeleteApiKey(ctx, &apigateway.DeleteApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "CreateUsagePlan", func() error {
		resp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name:        aws.String("test-usage-plan"),
			Description: aws.String("Test usage plan"),
			Throttle: &types.ThrottleSettings{
				BurstLimit: 10,
				RateLimit:  5.0,
			},
			Quota: &types.QuotaSettings{
				Limit:  1000,
				Period: types.QuotaPeriodTypeMonth,
			},
			Tags: map[string]string{
				"team": "backend",
			},
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("usage plan ID is nil")
		}
		if resp.Throttle == nil || resp.Throttle.BurstLimit != 10 {
			return fmt.Errorf("throttle burstLimit mismatch")
		}
		if resp.Quota == nil || resp.Quota.Period != types.QuotaPeriodTypeMonth {
			return fmt.Errorf("quota period mismatch")
		}
		return nil
	}))

	var usagePlanID string
	results = append(results, r.RunTest("apigateway", "GetUsagePlans", func() error {
		resp, err := client.GetUsagePlans(ctx, &apigateway.GetUsagePlansInput{
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil || len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 usage plan")
		}
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == "test-usage-plan" {
				usagePlanID = *item.Id
				break
			}
		}
		if usagePlanID == "" {
			return fmt.Errorf("test-usage-plan not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetUsagePlan", func() error {
		if usagePlanID == "" {
			return fmt.Errorf("usage plan ID not available")
		}
		resp, err := client.GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "test-usage-plan" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateUsagePlan", func() error {
		if usagePlanID == "" {
			return fmt.Errorf("usage plan ID not available")
		}
		resp, err := client.UpdateUsagePlan(ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/name"),
					Value: aws.String("updated-usage-plan"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "updated-usage-plan" {
			return fmt.Errorf("name not updated, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteUsagePlan", func() error {
		if usagePlanID == "" {
			return fmt.Errorf("usage plan ID not available")
		}
		_, err := client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "CreateUsagePlanKey_Lifecycle", func() error {
		keyResp, err := client.CreateApiKey(ctx, &apigateway.CreateApiKeyInput{
			Name:    aws.String("upk-test-key"),
			Enabled: true,
		})
		if err != nil {
			return fmt.Errorf("create api key: %v", err)
		}
		defer client.DeleteApiKey(ctx, &apigateway.DeleteApiKeyInput{ApiKey: keyResp.Id})

		upResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String("upk-test-plan"),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})

		upkResp, err := client.CreateUsagePlanKey(ctx, &apigateway.CreateUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
			KeyType:     aws.String("API_KEY"),
		})
		if err != nil {
			return fmt.Errorf("create usage plan key: %v", err)
		}
		if upkResp.Id == nil {
			return fmt.Errorf("usage plan key ID is nil")
		}

		getResp, err := client.GetUsagePlanKey(ctx, &apigateway.GetUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan key: %v", err)
		}
		if getResp.Type == nil || *getResp.Type != "API_KEY" {
			return fmt.Errorf("type mismatch, got %v", getResp.Type)
		}

		keysResp, err := client.GetUsagePlanKeys(ctx, &apigateway.GetUsagePlanKeysInput{
			UsagePlanId: upResp.Id,
			Limit:       aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("get usage plan keys: %v", err)
		}
		if keysResp.Items == nil || len(keysResp.Items) == 0 {
			return fmt.Errorf("expected at least 1 usage plan key")
		}

		_, err = client.DeleteUsagePlanKey(ctx, &apigateway.DeleteUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete usage plan key: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetUsage", func() error {
		upResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(fmt.Sprintf("usage-plan-%d", time.Now().UnixNano())),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})

		now := time.Now().UTC()
		startDate := now.AddDate(0, -1, 0).Format("2006-01-02")
		endDate := now.Format("2006-01-02")
		resp, err := client.GetUsage(ctx, &apigateway.GetUsageInput{
			UsagePlanId: upResp.Id,
			StartDate:   aws.String(startDate),
			EndDate:     aws.String(endDate),
		})
		if err != nil {
			return err
		}
		if resp.UsagePlanId == nil || *resp.UsagePlanId != *upResp.Id {
			return fmt.Errorf("usagePlanId mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateDomainName", func() error {
		domain := fmt.Sprintf("test-%d.example.com", time.Now().UnixNano())
		resp, err := client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:      aws.String(domain),
			CertificateName: aws.String("test-cert"),
			Tags: map[string]string{
				"domain": "test",
			},
		})
		if err != nil {
			return err
		}
		if resp.DomainName == nil || *resp.DomainName != domain {
			return fmt.Errorf("domain name mismatch, got %v", resp.DomainName)
		}
		if resp.DomainNameId == nil {
			return fmt.Errorf("domain name ID is nil")
		}
		return nil
	}))

	var domainName string
	results = append(results, r.RunTest("apigateway", "GetDomainNames", func() error {
		resp, err := client.GetDomainNames(ctx, &apigateway.GetDomainNamesInput{
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil || len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 domain name")
		}
		domainName = *resp.Items[0].DomainName
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDomainName", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.GetDomainName(ctx, &apigateway.GetDomainNameInput{
			DomainName: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		if resp.DomainName == nil || *resp.DomainName != domainName {
			return fmt.Errorf("domain name mismatch, got %v", resp.DomainName)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateDomainName", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/certificateName"),
					Value: aws.String("updated-cert"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.CertificateName == nil || *resp.CertificateName != "updated-cert" {
			return fmt.Errorf("certificateName not updated, got %v", resp.CertificateName)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateBasePathMapping", func() error {
		if apiID == "" || domainName == "" {
			return fmt.Errorf("API ID or domain name not available")
		}
		resp, err := client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
			DomainName: aws.String(domainName),
			RestApiId:  aws.String(apiID),
			BasePath:   aws.String("v1"),
			Stage:      aws.String("prod"),
		})
		if err != nil {
			return err
		}
		if resp.BasePath == nil || *resp.BasePath != "v1" {
			return fmt.Errorf("basePath mismatch, got %v", resp.BasePath)
		}
		if resp.RestApiId == nil || *resp.RestApiId != apiID {
			return fmt.Errorf("restApiId mismatch, got %v", resp.RestApiId)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetBasePathMappings", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.GetBasePathMappings(ctx, &apigateway.GetBasePathMappingsInput{
			DomainName: aws.String(domainName),
			Limit:      aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil || len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 base path mapping")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetBasePathMapping", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.GetBasePathMapping(ctx, &apigateway.GetBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
		})
		if err != nil {
			return err
		}
		if resp.BasePath == nil || *resp.BasePath != "v1" {
			return fmt.Errorf("basePath mismatch, got %v", resp.BasePath)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateBasePathMapping", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.UpdateBasePathMapping(ctx, &apigateway.UpdateBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/stage"),
					Value: aws.String("staging"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Stage == nil || *resp.Stage != "staging" {
			return fmt.Errorf("stage not updated, got %v", resp.Stage)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteBasePathMapping", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		_, err := client.DeleteBasePathMapping(ctx, &apigateway.DeleteBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "DeleteDomainName", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		_, err := client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{
			DomainName: aws.String(domainName),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "TagResource_UntagResource_ListTags", func() error {
		tagAPI := fmt.Sprintf("TagAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tagAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		arn := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s", r.region, *createResp.Id)

		_, err = client.TagResource(ctx, &apigateway.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}

		tagResp, err := client.GetTags(ctx, &apigateway.GetTagsInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("get tags: %v", err)
		}
		if tagResp.Tags == nil || tagResp.Tags["key1"] != "value1" {
			return fmt.Errorf("tags mismatch, got %v", tagResp.Tags)
		}

		_, err = client.UntagResource(ctx, &apigateway.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"key2"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		tagResp2, err := client.GetTags(ctx, &apigateway.GetTagsInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("get tags after untag: %v", err)
		}
		if _, exists := tagResp2.Tags["key2"]; exists {
			return fmt.Errorf("key2 should have been removed")
		}
		if tagResp2.Tags["key1"] != "value1" {
			return fmt.Errorf("key1 should still exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRestApi", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
			RestApiId: aws.String(apiID),
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

	results = append(results, r.RunTest("apigateway", "GetRestApi_NonExistent", func() error {
		_, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent API")
		}
		var nf *types.NotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRestApi_NonExistent", func() error {
		_, err := client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent API")
		}
		var nf *types.NotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStage_NonExistent", func() error {
		tmpAPI := fmt.Sprintf("TmpAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tmpAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("nonexistent_stage"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent stage")
		}
		var nf *types.NotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi_VerifyUpdate", func() error {
		uaAPI := fmt.Sprintf("UaAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(uaAPI),
			Description: aws.String("original desc"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		newDesc := "updated description v2"
		_, err = client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: createResp.Id,
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String(newDesc),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Description == nil || *resp.Description != newDesc {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateResource_NestedPath", func() error {
		crAPI := fmt.Sprintf("CrAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(crAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		usersResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("users"),
		})
		if err != nil {
			return fmt.Errorf("create users resource: %v", err)
		}

		userIdResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  usersResp.Id,
			PathPart:  aws.String("{userId}"),
		})
		if err != nil {
			return fmt.Errorf("create userId resource: %v", err)
		}
		if userIdResp.Path == nil || *userIdResp.Path != "/users/{userId}" {
			return fmt.Errorf("nested path mismatch, got %v", userIdResp.Path)
		}

		resResp, err := client.GetResources(ctx, &apigateway.GetResourcesInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		if len(resResp.Items) < 3 {
			return fmt.Errorf("expected at least 3 resources (root, users, {userId}), got %d", len(resResp.Items))
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateStage_VerifyConfig", func() error {
		csAPI := fmt.Sprintf("CsAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(csAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		depResp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		stageDesc := "test stage description"
		_, err = client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    createResp.Id,
			StageName:    aws.String("v1"),
			DeploymentId: depResp.Id,
			Description:  aws.String(stageDesc),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		resp, err := client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if resp.Description == nil || *resp.Description != stageDesc {
			return fmt.Errorf("stage description mismatch, got %v", resp.Description)
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != *depResp.Id {
			return fmt.Errorf("deployment ID mismatch, got %v", resp.DeploymentId)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApis_ContainsCreated", func() error {
		gaAPI := fmt.Sprintf("GaAPI-%d", time.Now().UnixNano())
		gaDesc := "searchable description"
		_, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(gaAPI),
			Description: aws.String(gaDesc),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: func() *string {
			resp, _ := client.GetRestApis(ctx, &apigateway.GetRestApisInput{Limit: aws.Int32(500)})
			if resp != nil {
				for _, item := range resp.Items {
					if item.Name != nil && *item.Name == gaAPI {
						return item.Id
					}
				}
			}
			return nil
		}()})

		resp, err := client.GetRestApis(ctx, &apigateway.GetRestApisInput{
			Limit: aws.Int32(500),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == gaAPI {
				found = true
				if item.Description == nil || *item.Description != gaDesc {
					return fmt.Errorf("description mismatch in list, got %v", item.Description)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created API %s not found in GetRestApis", gaAPI)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutMethod_AuthorizationTypes", func() error {
		pmAPI := fmt.Sprintf("PmAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(pmAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("secure"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		for _, authType := range []string{"NONE", "AWS_IAM", "CUSTOM"} {
			_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
				RestApiId:         createResp.Id,
				ResourceId:        resResp.Id,
				HttpMethod:        aws.String("GET"),
				AuthorizationType: aws.String(authType),
			})
			if err != nil {
				return fmt.Errorf("put method with auth %s: %v", authType, err)
			}
			getResp, err := client.GetMethod(ctx, &apigateway.GetMethodInput{
				RestApiId:  createResp.Id,
				ResourceId: resResp.Id,
				HttpMethod: aws.String("GET"),
			})
			if err != nil {
				return fmt.Errorf("get method with auth %s: %v", authType, err)
			}
			if getResp.AuthorizationType == nil || *getResp.AuthorizationType != authType {
				return fmt.Errorf("auth type mismatch for %s, got %v", authType, getResp.AuthorizationType)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegration_Types", func() error {
		itAPI := fmt.Sprintf("ItAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(itAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("inttest"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         createResp.Id,
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("POST"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		for _, intType := range []types.IntegrationType{
			types.IntegrationTypeMock,
			types.IntegrationTypeHttp,
			types.IntegrationTypeHttpProxy,
			types.IntegrationTypeAwsProxy,
		} {
			_, err = client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
				RestApiId:  createResp.Id,
				ResourceId: resResp.Id,
				HttpMethod: aws.String("POST"),
				Type:       intType,
			})
			if err != nil {
				return fmt.Errorf("put integration type %s: %v", intType, err)
			}
			getResp, err := client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
				RestApiId:  createResp.Id,
				ResourceId: resResp.Id,
				HttpMethod: aws.String("POST"),
			})
			if err != nil {
				return fmt.Errorf("get integration type %s: %v", intType, err)
			}
			if getResp.Type != intType {
				return fmt.Errorf("type mismatch, expected %s got %s", intType, getResp.Type)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "MethodWithIntegration_FullLifecycle", func() error {
		lcAPI := fmt.Sprintf("LcAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(lcAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("lifecycle"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         createResp.Id,
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
			OperationName:     aws.String("GetLifecycle"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		_, err = client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			RestApiId:             createResp.Id,
			ResourceId:            resResp.Id,
			HttpMethod:            aws.String("GET"),
			Type:                  types.IntegrationTypeMock,
			IntegrationHttpMethod: aws.String("POST"),
			Uri:                   aws.String("https://httpbin.org/post"),
			RequestParameters:     map[string]string{"integration.request.header.X-Custom": "'static'"},
			RequestTemplates:      map[string]string{"application/json": "{\"statusCode\":200}"},
			PassthroughBehavior:   aws.String("WHEN_NO_MATCH"),
			TimeoutInMillis:       aws.Int32(3000),
			CacheNamespace:        aws.String("lifecycle"),
			CacheKeyParameters:    []string{"header.Authorization"},
		})
		if err != nil {
			return fmt.Errorf("put integration: %v", err)
		}

		getIntResp, err := client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get integration: %v", err)
		}
		if getIntResp.Uri == nil || *getIntResp.Uri != "https://httpbin.org/post" {
			return fmt.Errorf("uri mismatch, got %v", getIntResp.Uri)
		}
		if getIntResp.TimeoutInMillis != 3000 {
			return fmt.Errorf("timeoutInMillis mismatch, got %d", getIntResp.TimeoutInMillis)
		}

		_, err = client.PutIntegrationResponse(ctx, &apigateway.PutIntegrationResponseInput{
			RestApiId:          createResp.Id,
			ResourceId:         resResp.Id,
			HttpMethod:         aws.String("GET"),
			StatusCode:         aws.String("200"),
			ResponseParameters: map[string]string{"method.response.header.Content-Type": "integration.response.header.Content-Type"},
			ResponseTemplates:  map[string]string{"application/json": "$input.json('$')"},
			SelectionPattern:   aws.String("2\\d{2}"),
		})
		if err != nil {
			return fmt.Errorf("put integration response: %v", err)
		}

		_, err = client.PutMethodResponse(ctx, &apigateway.PutMethodResponseInput{
			RestApiId:          createResp.Id,
			ResourceId:         resResp.Id,
			HttpMethod:         aws.String("GET"),
			StatusCode:         aws.String("200"),
			ResponseParameters: map[string]bool{"method.response.header.Content-Type": true},
			ResponseModels:     map[string]string{"application/json": "Empty"},
		})
		if err != nil {
			return fmt.Errorf("put method response: %v", err)
		}

		_, err = client.DeleteMethodResponse(ctx, &apigateway.DeleteMethodResponseInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete method response: %v", err)
		}

		_, err = client.DeleteIntegrationResponse(ctx, &apigateway.DeleteIntegrationResponseInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete integration response: %v", err)
		}

		_, err = client.DeleteIntegration(ctx, &apigateway.DeleteIntegrationInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete integration: %v", err)
		}

		_, err = client.DeleteMethod(ctx, &apigateway.DeleteMethodInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete method: %v", err)
		}

		_, err = client.DeleteResource(ctx, &apigateway.DeleteResourceInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete resource: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Deployment_Stage_FullLifecycle", func() error {
		dsAPI := fmt.Sprintf("DsAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(dsAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		depResp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId:   createResp.Id,
			Description: aws.String("v1 deployment"),
		})
		if err != nil {
			return fmt.Errorf("create deployment: %v", err)
		}

		getDepResp, err := client.GetDeployment(ctx, &apigateway.GetDeploymentInput{
			RestApiId:    createResp.Id,
			DeploymentId: depResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get deployment: %v", err)
		}
		if getDepResp.Description == nil || *getDepResp.Description != "v1 deployment" {
			return fmt.Errorf("deployment description mismatch, got %v", getDepResp.Description)
		}

		_, err = client.UpdateDeployment(ctx, &apigateway.UpdateDeploymentInput{
			RestApiId:    createResp.Id,
			DeploymentId: depResp.Id,
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("v1 updated"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update deployment: %v", err)
		}

		_, err = client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    createResp.Id,
			StageName:    aws.String("production"),
			DeploymentId: depResp.Id,
			Description:  aws.String("production stage"),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		_, err = client.UpdateStage(ctx, &apigateway.UpdateStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("production"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("production updated"),
				},
				{
					Op:    types.OpReplace,
					Path:  aws.String("/variables/env"),
					Value: aws.String("prod"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update stage: %v", err)
		}

		stageResp, err := client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("production"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if stageResp.Variables == nil || stageResp.Variables["env"] != "prod" {
			return fmt.Errorf("stage variables not set, got %v", stageResp.Variables)
		}

		_, err = client.DeleteStage(ctx, &apigateway.DeleteStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("production"),
		})
		if err != nil {
			return fmt.Errorf("delete stage: %v", err)
		}

		_, err = client.DeleteDeployment(ctx, &apigateway.DeleteDeploymentInput{
			RestApiId:    createResp.Id,
			DeploymentId: depResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete deployment: %v", err)
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

	results = append(results, r.RunTest("apigateway", "Authorizer_FullLifecycle", func() error {
		auAPI := fmt.Sprintf("AuAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(auAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		authResp, err := client.CreateAuthorizer(ctx, &apigateway.CreateAuthorizerInput{
			RestApiId:                    createResp.Id,
			Name:                         aws.String("lambda-auth"),
			Type:                         types.AuthorizerTypeToken,
			AuthorizerUri:                aws.String("https://example.com/lambda"),
			IdentitySource:               aws.String("method.request.header.Auth"),
			AuthorizerCredentials:        aws.String("arn:aws:iam::123456789012:role/lambda-auth-role"),
			IdentityValidationExpression: aws.String("Bearer .*"),
			AuthorizerResultTtlInSeconds: aws.Int32(600),
		})
		if err != nil {
			return fmt.Errorf("create authorizer: %v", err)
		}
		if authResp.AuthorizerResultTtlInSeconds == nil || *authResp.AuthorizerResultTtlInSeconds != 600 {
			return fmt.Errorf("ttl mismatch, got %v", authResp.AuthorizerResultTtlInSeconds)
		}

		_, err = client.UpdateAuthorizer(ctx, &apigateway.UpdateAuthorizerInput{
			RestApiId:    createResp.Id,
			AuthorizerId: authResp.Id,
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/authorizerResultTtlInSeconds"),
					Value: aws.String("1200"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update authorizer: %v", err)
		}

		getAuthResp, err := client.GetAuthorizer(ctx, &apigateway.GetAuthorizerInput{
			RestApiId:    createResp.Id,
			AuthorizerId: authResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get authorizer: %v", err)
		}
		if getAuthResp.AuthorizerResultTtlInSeconds == nil || *getAuthResp.AuthorizerResultTtlInSeconds != 1200 {
			return fmt.Errorf("ttl not updated, got %v", getAuthResp.AuthorizerResultTtlInSeconds)
		}

		authListResp, err := client.GetAuthorizers(ctx, &apigateway.GetAuthorizersInput{
			RestApiId: createResp.Id,
			Limit:     aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("list authorizers: %v", err)
		}
		if len(authListResp.Items) == 0 {
			return fmt.Errorf("expected at least 1 authorizer")
		}

		_, err = client.DeleteAuthorizer(ctx, &apigateway.DeleteAuthorizerInput{
			RestApiId:    createResp.Id,
			AuthorizerId: authResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete authorizer: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DomainName_BasePathMapping_FullLifecycle", func() error {
		dbAPI := fmt.Sprintf("DbAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(dbAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		domain := fmt.Sprintf("lc-%d.example.com", time.Now().UnixNano())
		dnResp, err := client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:      aws.String(domain),
			CertificateName: aws.String("lc-cert"),
		})
		if err != nil {
			return fmt.Errorf("create domain: %v", err)
		}

		_, err = client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
			DomainName: aws.String(domain),
			RestApiId:  createResp.Id,
			BasePath:   aws.String("(none)"),
			Stage:      aws.String("prod"),
		})
		if err != nil {
			return fmt.Errorf("create base path mapping: %v", err)
		}

		_, err = client.GetBasePathMappings(ctx, &apigateway.GetBasePathMappingsInput{
			DomainName: aws.String(domain),
		})
		if err != nil {
			return fmt.Errorf("get base path mappings: %v", err)
		}

		_, err = client.DeleteBasePathMapping(ctx, &apigateway.DeleteBasePathMappingInput{
			DomainName: aws.String(domain),
			BasePath:   aws.String("(none)"),
		})
		if err != nil {
			return fmt.Errorf("delete base path mapping: %v", err)
		}

		_ = dnResp
		_, err = client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{
			DomainName: aws.String(domain),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "UsagePlan_WithApiStages", func() error {
		usAPI := fmt.Sprintf("UsAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(usAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		depResp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		_, err = client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    createResp.Id,
			StageName:    aws.String("api-stage"),
			DeploymentId: depResp.Id,
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		upResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(fmt.Sprintf("us-plan-%d", time.Now().UnixNano())),
			ApiStages: []types.ApiStage{
				{
					ApiId: createResp.Id,
					Stage: aws.String("api-stage"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}

		getResp, err := client.GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: upResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}
		if len(getResp.ApiStages) == 0 {
			return fmt.Errorf("expected apiStages to be set")
		}

		_ = upResp
		return nil
	}))

	return results
}
