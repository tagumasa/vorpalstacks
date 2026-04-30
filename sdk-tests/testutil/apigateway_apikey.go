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

func (r *TestRunner) runAPIGatewayApiKeyTests(ctx context.Context, client *apigateway.Client) []TestResult {
	var results []TestResult

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
		if resp.Name == nil || *resp.Name != "test-api-key" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Value == nil {
			return fmt.Errorf("api key value is nil")
		}
		if resp.Description == nil || *resp.Description != "Test API key" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		if !resp.Enabled {
			return fmt.Errorf("expected enabled=true")
		}
		if resp.Tags == nil || resp.Tags["env"] != "test" {
			return fmt.Errorf("tags mismatch, got %v", resp.Tags)
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
		if len(resp.Items) == 0 {
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
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetApiKey(ctx, &apigateway.GetApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		if err == nil {
			return fmt.Errorf("GetApiKey should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateApiKey_DefaultEnabled", func() error {
		keyName := fmt.Sprintf("default-key-%d", time.Now().UnixNano())
		resp, err := client.CreateApiKey(ctx, &apigateway.CreateApiKeyInput{
			Name: aws.String(keyName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteApiKey(ctx, &apigateway.DeleteApiKeyInput{ApiKey: resp.Id})

		if !resp.Enabled {
			return fmt.Errorf("expected enabled=true by default, got false")
		}
		return nil
	}))

	return results
}
