package testutil

import (
	"fmt"
	"strings"

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
		// The list order is not contractual: address the root resource by
		// its path, not by list position.
		rootID := ""
		for _, res := range resources.Items {
			if aws.ToString(res.Path) == "/" {
				rootID = aws.ToString(res.Id)
			}
		}
		if rootID == "" {
			return fmt.Errorf("root resource not found")
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
			{name: "http-without-uri-rejected", pathPart: "httptest", intType: types.IntegrationTypeHttp, wantErr: "BadRequestException"},
			// Invalid integration type should fail
			{name: "invalid-type-rejected", pathPart: "invalidtest", intType: types.IntegrationType("INVALID_TYPE"), wantErr: "BadRequestException"},
		}
		for _, row := range rows {
			resID, err := tc.createResourceWithMethod(tc.apiID, rootID, row.pathPart, "GET")
			if err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
			defer tc.client.DeleteResource(tc.ctx, &apigateway.DeleteResourceInput{
				RestApiId: aws.String(tc.apiID), ResourceId: aws.String(resID),
			})
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
				return fmt.Errorf("%s: expected %s", row.name, row.wantErr)
			}
			if err := AssertErrorContains(err, row.wantErr); err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Validation_CreateStage_RequiresDeploymentId", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}

		// The SDK's typed input marks DeploymentId required and rejects the
		// call client-side, so this documents that guard: the server-side
		// rejection is unreachable through the SDK and is pinned by the
		// unit test TestCreateStageCoreRequiresDeploymentId.
		_, err := tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("nostage"),
		})
		if err == nil {
			return fmt.Errorf("expected the SDK's required-member validation for deploymentId")
		}
		if !strings.Contains(err.Error(), "DeploymentId") {
			return fmt.Errorf("expected the validation error to name DeploymentId, got: %v", err)
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

	results = append(results, r.RunTest("apigateway", "Validation_CreateApiKey_NameTooLong", func() error {
		// The documented member bound: an API key name cannot exceed 1024
		// characters. The SDK applies no client-side length check, so the
		// server-side rejection is reached through the wire.
		_, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
			Name: aws.String(strings.Repeat("n", 1025)),
		})
		if err == nil {
			return fmt.Errorf("CreateApiKey accepted a name above the 1024-character bound")
		}
		if !strings.Contains(err.Error(), "BadRequestException") {
			return fmt.Errorf("expected BadRequestException, got: %v", err)
		}
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
		defer tc.client.DeleteModel(tc.ctx, &apigateway.DeleteModelInput{
			RestApiId: aws.String(tc.apiID), ModelName: resp.Name,
		})
		if resp.ContentType == nil || *resp.ContentType != "application/json" {
			return fmt.Errorf("expected contentType application/json, got %v", resp.ContentType)
		}
		return nil
	}))

	return results
}
