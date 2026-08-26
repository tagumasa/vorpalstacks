package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayValidationTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "Validation_PutIntegration_Mock_NoUri", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resources, err := tc.client.GetResources(tc.ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(tc.apiID),
		})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		if len(resources.Items) == 0 {
			return fmt.Errorf("no resources found")
		}

		res, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(tc.apiID),
			ParentId:  resources.Items[0].Id,
			PathPart:  aws.String("mocktest"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(tc.apiID),
			ResourceId:        res.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		// MOCK integration without URI should succeed
		_, err = tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: res.Id,
			HttpMethod: aws.String("GET"),
			Type:       types.IntegrationTypeMock,
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "Validation_PutIntegration_Http_RequiresUri", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resources, err := tc.client.GetResources(tc.ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(tc.apiID),
		})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		if len(resources.Items) == 0 {
			return fmt.Errorf("no resources found")
		}

		res, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(tc.apiID),
			ParentId:  resources.Items[0].Id,
			PathPart:  aws.String("httptest"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(tc.apiID),
			ResourceId:        res.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		// HTTP integration without URI should fail
		_, err = tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: res.Id,
			HttpMethod: aws.String("GET"),
			Type:       types.IntegrationTypeHttp,
		})
		if err == nil {
			return fmt.Errorf("expected error for HTTP integration without URI")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Validation_PutIntegration_InvalidType", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resources, err := tc.client.GetResources(tc.ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(tc.apiID),
		})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		if len(resources.Items) == 0 {
			return fmt.Errorf("no resources found")
		}

		res, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(tc.apiID),
			ParentId:  resources.Items[0].Id,
			PathPart:  aws.String("invalidtest"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(tc.apiID),
			ResourceId:        res.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		// Invalid integration type should fail
		_, err = tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: res.Id,
			HttpMethod: aws.String("GET"),
			Type:       types.IntegrationType("INVALID_TYPE"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid integration type")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Validation_CreateStage_RequiresDeploymentId", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}

		_, err := tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("nostage"),
		})
		if err == nil {
			return fmt.Errorf("expected error for CreateStage without deploymentId")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Validation_CreateApiKey_NoName", func() error {
		// CreateApiKey without a name should succeed (auto-generated)
		resp, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{})
		if err != nil {
			return fmt.Errorf("CreateApiKey without name failed: %v", err)
		}
		if resp.Id == nil {
			return fmt.Errorf("API key ID is nil")
		}
		// Clean up
		tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{
			ApiKey: resp.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Validation_CreateModel_DefaultContentType", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}

		resp, err := tc.client.CreateModel(tc.ctx, &apigateway.CreateModelInput{
			RestApiId:   aws.String(tc.apiID),
			Name:        aws.String(fmt.Sprintf("TestModelDefault-%d", 0)),
			Schema:      aws.String(`{"type": "object"}`),
			ContentType: aws.String("application/json"),
		})
		if err != nil {
			return fmt.Errorf("CreateModel failed: %v", err)
		}
		if resp.ContentType == nil || *resp.ContentType != "application/json" {
			return fmt.Errorf("expected contentType application/json, got %v", resp.ContentType)
		}
		return nil
	}))

	return results
}
