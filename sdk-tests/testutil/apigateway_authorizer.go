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

func (r *TestRunner) runAPIGatewayAuthorizerTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

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
		if len(resp.Items) == 0 {
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
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetAuthorizer(ctx, &apigateway.GetAuthorizerInput{
			RestApiId:    aws.String(apiID),
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

	return results
}
