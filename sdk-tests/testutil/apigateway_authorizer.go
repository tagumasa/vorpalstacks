package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayAuthorizerTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var authorizerID string
	results = append(results, r.RunTest("apigateway", "CreateAuthorizer", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, authorizerID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, authorizerID); err != nil {
			return err
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

		// The /name row documents replace only: add rejects.
		_, err = tc.client.UpdateAuthorizer(tc.ctx, &apigateway.UpdateAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: aws.String(authorizerID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/name"), Value: aws.String("nope")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for add on /name, got: %v", err)
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
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, authorizerID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, authorizerID); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetAuthorizer should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	return results
}
