package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayAuthorizerTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var authorizerID string
	results = append(results, r.RunTest("apigateway", "CreateAuthorizer", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.CreateAuthorizer(tc.ctx, &apigateway.CreateAuthorizerInput{
			RestApiId:                    aws.String(tc.apiID),
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
		if resp.Name == nil || *resp.Name != "test-authorizer" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Type != types.AuthorizerTypeToken {
			return fmt.Errorf("type mismatch, got %v", resp.Type)
		}
		if resp.AuthorizerUri == nil || *resp.AuthorizerUri != "https://example.com/auth" {
			return fmt.Errorf("authorizerUri mismatch, got %v", resp.AuthorizerUri)
		}
		if resp.IdentitySource == nil || *resp.IdentitySource != "method.request.header.Authorization" {
			return fmt.Errorf("identitySource mismatch, got %v", resp.IdentitySource)
		}
		if resp.AuthorizerResultTtlInSeconds == nil || *resp.AuthorizerResultTtlInSeconds != 300 {
			return fmt.Errorf("authorizerResultTtlInSeconds mismatch, got %v", resp.AuthorizerResultTtlInSeconds)
		}
		authorizerID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetAuthorizer", func() error {
		if tc.apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		resp, err := tc.client.GetAuthorizer(tc.ctx, &apigateway.GetAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: aws.String(authorizerID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "test-authorizer" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Type != types.AuthorizerTypeToken {
			return fmt.Errorf("type mismatch, got %v", resp.Type)
		}
		if resp.AuthorizerUri == nil || *resp.AuthorizerUri != "https://example.com/auth" {
			return fmt.Errorf("authorizerUri mismatch, got %v", resp.AuthorizerUri)
		}
		if resp.IdentitySource == nil || *resp.IdentitySource != "method.request.header.Authorization" {
			return fmt.Errorf("identitySource mismatch, got %v", resp.IdentitySource)
		}
		if resp.AuthorizerResultTtlInSeconds == nil || *resp.AuthorizerResultTtlInSeconds != 300 {
			return fmt.Errorf("authorizerResultTtlInSeconds mismatch, got %v", resp.AuthorizerResultTtlInSeconds)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateAuthorizer", func() error {
		if tc.apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		resp, err := tc.client.UpdateAuthorizer(tc.ctx, &apigateway.UpdateAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
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

		// Verify a TTL change persists via a fresh read of the authorizer.
		_, err = tc.client.UpdateAuthorizer(tc.ctx, &apigateway.UpdateAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: aws.String(authorizerID),
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

		getResp, err := tc.client.GetAuthorizer(tc.ctx, &apigateway.GetAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: aws.String(authorizerID),
		})
		if err != nil {
			return fmt.Errorf("get authorizer: %v", err)
		}
		if getResp.AuthorizerResultTtlInSeconds == nil || *getResp.AuthorizerResultTtlInSeconds != 1200 {
			return fmt.Errorf("ttl not updated, got %v", getResp.AuthorizerResultTtlInSeconds)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetAuthorizers", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		items, err := tc.allAuthorizers(tc.apiID)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return fmt.Errorf("expected at least 1 authorizer")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "TestInvokeAuthorizer", func() error {
		if tc.apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		resp, err := tc.client.TestInvokeAuthorizer(tc.ctx, &apigateway.TestInvokeAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
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
		if tc.apiID == "" || authorizerID == "" {
			return fmt.Errorf("API ID or authorizer ID not available")
		}
		_, err := tc.client.DeleteAuthorizer(tc.ctx, &apigateway.DeleteAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: aws.String(authorizerID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetAuthorizer(tc.ctx, &apigateway.GetAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: aws.String(authorizerID),
		})
		if err == nil {
			return fmt.Errorf("GetAuthorizer should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	return results
}
