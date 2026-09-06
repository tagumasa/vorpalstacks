package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayValidationTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "Validation_PutIntegration", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		rows := []struct {
			name     string
			pathPart string
			intType  types.IntegrationType
			wantErr  string // empty means the put must succeed
		}{
			// MOCK integration without URI should succeed
			{name: "mock-without-uri-succeeds", pathPart: "mocktest", intType: types.IntegrationTypeMock},
			// HTTP integration without URI should fail
			{name: "http-without-uri-rejected", pathPart: "httptest", intType: types.IntegrationTypeHttp, wantErr: "expected error for HTTP integration without URI"},
			// Invalid integration type should fail
			{name: "invalid-type-rejected", pathPart: "invalidtest", intType: types.IntegrationType("INVALID_TYPE"), wantErr: "expected error for invalid integration type"},
		}
		for _, row := range rows {
			resID, err := tc.createResourceWithMethod(tc.apiID, aws.ToString(resources.Items[0].Id), row.pathPart, "GET")
			if err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
			_, err = tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
				RestApiId:  aws.String(tc.apiID),
				ResourceId: aws.String(resID),
				HttpMethod: aws.String("GET"),
				Type:       row.intType,
			})
			if row.wantErr == "" {
				if err != nil {
					return fmt.Errorf("%s: %v", row.name, err)
				}
				continue
			}
			if err == nil {
				return fmt.Errorf("%s: %s", row.name, row.wantErr)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Validation_CreateStage_RequiresDeploymentId", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		if err := tc.require(tc.apiID); err != nil {
			return err
		}

		resp, err := tc.client.CreateModel(tc.ctx, &apigateway.CreateModelInput{
			RestApiId:   aws.String(tc.apiID),
			Name:        aws.String(tc.uniqueName("TestModelDefault")),
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
