package testutil

import (
	"bytes"
	"context"
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

	authApiId, authApiKeyId, err := setupGraphQLAuthAPI(client, ctx, res.uid)
	if err != nil {
		return append(results, TestResult{
			Service:  "appsync",
			TestName: "GraphQLAuth_Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}

	// Best-effort cleanup of the dedicated API key and API; ignore errors.
	defer func() {
		if authApiKeyId != "" {
			client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
				ApiId: aws.String(authApiId),
				Id:    aws.String(authApiKeyId),
			})
		}
		client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{
			ApiId: aws.String(authApiId),
		})
	}()

	// Build the GraphQL endpoint URL.
	graphqlURL := fmt.Sprintf("%s/v1/apis/%s/graphql", r.endpoint, authApiId)

	// Helper to send a raw GraphQL POST with an optional API key.
	sendGraphQL := func(apiKey string) (*http.Response, []byte, error) {
		body := []byte(`{"query":"{ __typename }"}`)
		req, err := http.NewRequestWithContext(ctx, "POST", graphqlURL, bytes.NewReader(body))
		if err != nil {
			return nil, nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if apiKey != "" {
			req.Header.Set("x-api-key", apiKey)
		}
		resp, err := r.client.Do(req)
		if err != nil {
			return nil, nil, fmt.Errorf("send request: %w", err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, nil, fmt.Errorf("read response body: %w", err)
		}
		return resp, respBody, nil
	}

	// Test: Missing API key -> 401 UnauthorizedException
	results = append(results, r.RunTest("appsync", "GraphQL_MissingApiKey", func() error {
		resp, respBody, err := sendGraphQL("")
		if err != nil {
			return fmt.Errorf("send GraphQL: %w", err)
		}
		return assertUnauthorized(resp, respBody)
	}))

	// Test: Invalid API key -> 401 UnauthorizedException
	results = append(results, r.RunTest("appsync", "GraphQL_InvalidApiKey", func() error {
		resp, respBody, err := sendGraphQL("invalid-key-value")
		if err != nil {
			return fmt.Errorf("send GraphQL: %w", err)
		}
		return assertUnauthorized(resp, respBody)
	}))

	// Test: Valid API key -> should NOT be 401
	results = append(results, r.RunTest("appsync", "GraphQL_ValidApiKey", func() error {
		resp, respBody, err := sendGraphQL(authApiKeyId)
		if err != nil {
			return fmt.Errorf("send GraphQL: %w", err)
		}
		if resp == nil {
			return fmt.Errorf("no response")
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("valid API key was rejected with 401 (body: %s)", string(respBody))
		}
		return nil
	}))

	return results
}

// setupGraphQLAuthAPI provisions the dedicated GraphQL API and its API key
// backing the auth tests. It runs at section start and its failure surfaces
// as the GraphQLAuth_Setup FAIL row.
func setupGraphQLAuthAPI(client *appsync.Client, ctx context.Context, uid int64) (string, string, error) {
	resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
		Name:               aws.String(fmt.Sprintf("test-gqlauth-%d", uid)),
		AuthenticationType: types.AuthenticationTypeApiKey,
	})
	if err != nil {
		return "", "", err
	}
	if resp.GraphqlApi == nil || resp.GraphqlApi.ApiId == nil || *resp.GraphqlApi.ApiId == "" {
		return "", "", fmt.Errorf("expected GraphqlApi with ApiId in CreateGraphqlApi response")
	}
	apiId := *resp.GraphqlApi.ApiId

	keyResp, err := client.CreateApiKey(ctx, &appsync.CreateApiKeyInput{
		ApiId: aws.String(apiId),
	})
	if err != nil {
		_, _ = client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(apiId)})
		return "", "", fmt.Errorf("failed to create API key: %v", err)
	}
	if keyResp.ApiKey == nil || keyResp.ApiKey.Id == nil || *keyResp.ApiKey.Id == "" {
		_, _ = client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(apiId)})
		return "", "", fmt.Errorf("expected ApiKey with Id in CreateApiKey response")
	}
	return apiId, *keyResp.ApiKey.Id, nil
}

// assertUnauthorized verifies that a GraphQL data-plane response rejected
// the request with 401 and an UnauthorizedException error entry.
func assertUnauthorized(resp *http.Response, respBody []byte) error {
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
}
