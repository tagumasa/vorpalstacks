package testutil

import (
	"encoding/base64"
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

		// The remaining rows run on a private authorizer: the type row, the
		// ARN-enforced authorizerUri row, identitySource and
		// authorizerCredentials. The shared fixture keeps its shape for the
		// invoke test that follows.
		ownResp, err := tc.client.CreateAuthorizer(tc.ctx, &apigateway.CreateAuthorizerInput{
			RestApiId:      aws.String(tc.apiID),
			Name:           aws.String(tc.uniqueName("row-authorizer")),
			Type:           types.AuthorizerTypeToken,
			AuthorizerUri:  aws.String("https://example.com/auth"),
			IdentitySource: aws.String("method.request.header.Authorization"),
		})
		if err != nil {
			return fmt.Errorf("create own authorizer: %v", err)
		}
		defer tc.client.DeleteAuthorizer(tc.ctx, &apigateway.DeleteAuthorizerInput{
			RestApiId: aws.String(tc.apiID), AuthorizerId: ownResp.Id,
		})
		ownUpd, err := tc.client.UpdateAuthorizer(tc.ctx, &apigateway.UpdateAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: ownResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/type"), Value: aws.String("REQUEST")},
				{Op: types.OpReplace, Path: aws.String("/authorizerUri"), Value: aws.String("arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/row-fn/invocations")},
				{Op: types.OpAdd, Path: aws.String("/identitySource"), Value: aws.String("method.request.header.X-Row")},
				{Op: types.OpReplace, Path: aws.String("/authorizerCredentials"), Value: aws.String("arn:aws:iam::123456789012:role/apigateway-rows")},
			},
		})
		if err != nil {
			return fmt.Errorf("row patch: %v", err)
		}
		if ownUpd.Type != types.AuthorizerTypeRequest {
			return fmt.Errorf("type row not applied, got %v", ownUpd.Type)
		}
		if aws.ToString(ownUpd.AuthorizerUri) != "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/row-fn/invocations" {
			return fmt.Errorf("authorizerUri row not applied, got %v", ownUpd.AuthorizerUri)
		}
		if aws.ToString(ownUpd.IdentitySource) != "method.request.header.X-Row" {
			return fmt.Errorf("identitySource row not applied, got %v", ownUpd.IdentitySource)
		}
		if aws.ToString(ownUpd.AuthorizerCredentials) != "arn:aws:iam::123456789012:role/apigateway-rows" {
			return fmt.Errorf("authorizerCredentials row not applied, got %v", ownUpd.AuthorizerCredentials)
		}
		// The authorizerUri row rejects a non-ARN value.
		_, err = tc.client.UpdateAuthorizer(tc.ctx, &apigateway.UpdateAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: ownResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/authorizerUri"), Value: aws.String("https://example.com/other")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for a non-ARN authorizerUri, got: %v", err)
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
		found := false
		for _, a := range items {
			if aws.ToString(a.Id) == authorizerID {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("created authorizer %q not found in list", authorizerID)
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

		// Claims is the Cognito-path member: the caller's own token claims
		// are reported for a COGNITO_USER_POOLS authorizer.
		cognitoResp, err := tc.client.CreateAuthorizer(tc.ctx, &apigateway.CreateAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			Name:         aws.String(tc.uniqueName("cognito-authorizer")),
			Type:         types.AuthorizerTypeCognitoUserPools,
			ProviderARNs: []string{"arn:aws:cognito-idp:us-east-1:000000000000:userpool/us-east-1_ABC123"},
		})
		if err != nil {
			return fmt.Errorf("create cognito authorizer: %v", err)
		}
		defer tc.client.DeleteAuthorizer(tc.ctx, &apigateway.DeleteAuthorizerInput{
			RestApiId: aws.String(tc.apiID), AuthorizerId: cognitoResp.Id,
		})

		b64 := func(v string) string { return base64.RawURLEncoding.EncodeToString([]byte(v)) }
		jwt := strings.Join([]string{
			b64(`{"alg":"none","typ":"JWT"}`),
			b64(`{"sub":"pin-user","iat":1700000000}`),
			b64("signature"),
		}, ".")
		cognito, err := tc.client.TestInvokeAuthorizer(tc.ctx, &apigateway.TestInvokeAuthorizerInput{
			RestApiId:    aws.String(tc.apiID),
			AuthorizerId: cognitoResp.Id,
			Headers:      map[string]string{"Authorization": jwt},
		})
		if err != nil {
			return err
		}
		if cognito.ClientStatus != 200 {
			return fmt.Errorf("expected clientStatus 200 for cognito authorizer, got %d", cognito.ClientStatus)
		}
		if cognito.Claims["sub"] != "pin-user" {
			return fmt.Errorf("token claims not reported, got %v", cognito.Claims)
		}
		if cognito.Claims["iat"] != "1700000000" {
			return fmt.Errorf("numeric claim not stringified, got %v", cognito.Claims)
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
