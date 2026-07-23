package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

// runAppSyncGraphQLAuthTests verifies API Key authentication on the
// GraphQL data-plane endpoint (POST /v1/apis/{apiId}/graphql).
// Uses a dedicated API to avoid interference from other tests that
// may change the authenticationType.
func (r *TestRunner) runAppSyncGraphQLAuthTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client
	uid := res.uid

	// Create a dedicated GraphQL API with API_KEY auth.
	var authApiId string
	var authApiKeyId string

	results = append(results, r.RunTest("appsync", "GraphQLAuth_Setup", func() error {
		resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
			Name:               aws.String(fmt.Sprintf("test-gqlauth-%d", uid)),
			AuthenticationType: types.AuthenticationTypeApiKey,
		})
		if err != nil {
			return err
		}
		authApiId = *resp.GraphqlApi.ApiId

		// Create an API key for this API.
		keyResp, err := client.CreateApiKey(ctx, &appsync.CreateApiKeyInput{
			ApiId: aws.String(authApiId),
		})
		if err != nil {
			return err
		}
		authApiKeyId = *keyResp.ApiKey.Id
		return nil
	}))

	if authApiId == "" {
		return results
	}

	// Build the GraphQL endpoint URL.
	graphqlURL := fmt.Sprintf("%s/v1/apis/%s/graphql", r.endpoint, authApiId)

	// Helper to send a raw GraphQL POST with an optional API key.
	sendGraphQL := func(apiKey string) (*http.Response, []byte) {
		body := []byte(`{"query":"{ __typename }"}`)
		req, err := http.NewRequestWithContext(ctx, "POST", graphqlURL, bytes.NewReader(body))
		if err != nil {
			return nil, nil
		}
		req.Header.Set("Content-Type", "application/json")
		if apiKey != "" {
			req.Header.Set("x-api-key", apiKey)
		}
		resp, err := r.client.Do(req)
		if err != nil {
			return nil, nil
		}
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		return resp, respBody
	}

	// Test: Missing API key -> 401 UnauthorizedException
	results = append(results, r.RunTest("appsync", "GraphQL_MissingApiKey", func() error {
		resp, respBody := sendGraphQL("")
		if resp == nil {
			return fmt.Errorf("no response")
		}
		if resp.StatusCode != http.StatusUnauthorized {
			return fmt.Errorf("expected 401, got %d (body: %s)", resp.StatusCode, string(respBody))
		}
		var result struct {
			Errors []struct {
				ErrorType string `json:"errorType"`
				Message   string `json:"message"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return fmt.Errorf("failed to parse error body: %w", err)
		}
		if len(result.Errors) == 0 || result.Errors[0].ErrorType != "UnauthorizedException" {
			return fmt.Errorf("expected UnauthorizedException error, got: %s", string(respBody))
		}
		return nil
	}))

	// Test: Invalid API key -> 401 UnauthorizedException
	results = append(results, r.RunTest("appsync", "GraphQL_InvalidApiKey", func() error {
		resp, respBody := sendGraphQL("invalid-key-value")
		if resp == nil {
			return fmt.Errorf("no response")
		}
		if resp.StatusCode != http.StatusUnauthorized {
			return fmt.Errorf("expected 401, got %d (body: %s)", resp.StatusCode, string(respBody))
		}
		var result struct {
			Errors []struct {
				ErrorType string `json:"errorType"`
				Message   string `json:"message"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return fmt.Errorf("failed to parse error body: %w", err)
		}
		if len(result.Errors) == 0 || result.Errors[0].ErrorType != "UnauthorizedException" {
			return fmt.Errorf("expected UnauthorizedException error, got: %s", string(respBody))
		}
		return nil
	}))

	// Test: Valid API key -> should NOT be 401
	results = append(results, r.RunTest("appsync", "GraphQL_ValidApiKey", func() error {
		resp, respBody := sendGraphQL(authApiKeyId)
		if resp == nil {
			return fmt.Errorf("no response")
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("valid API key was rejected with 401 (body: %s)", string(respBody))
		}
		return nil
	}))

	// Cleanup: delete the API key and API.
	results = append(results, r.RunTest("appsync", "GraphQLAuth_Cleanup", func() error {
		// Best-effort cleanup; ignore errors.
		if authApiKeyId != "" {
			client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
				ApiId: aws.String(authApiId),
				Id:    aws.String(authApiKeyId),
			})
		}
		client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{
			ApiId: aws.String(authApiId),
		})
		return nil
	}))

	return results
}
