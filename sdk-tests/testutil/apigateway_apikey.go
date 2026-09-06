package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayApiKeyTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var apiKeyValue string
	results = append(results, r.RunTest("apigateway", "CreateApiKey", func() error {
		resp, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
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
		items, err := tc.allApiKeys()
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return fmt.Errorf("expected at least 1 api key")
		}
		found := containsID(items, func(item *types.ApiKey) bool {
			return item.Name != nil && *item.Name == "test-api-key"
		})
		if found == nil {
			return fmt.Errorf("test-api-key not found")
		}
		apiKeyID = *found.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetApiKey", func() error {
		if err := tc.require(apiKeyID); err != nil {
			return err
		}
		resp, err := tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{
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
		if err := tc.require(apiKeyID); err != nil {
			return err
		}
		resp, err := tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
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

		// The /labels row documents add and remove; ApiKey's only
		// string-to-string map is tags, and the value travels as a JSON
		// object.
		_, err = tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/labels"), Value: aws.String(`{"team":"platform"}`)},
			},
		})
		if err != nil {
			return err
		}
		keyResp, err := tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		if err != nil {
			return fmt.Errorf("get api key: %v", err)
		}
		if keyResp.Tags == nil || keyResp.Tags["team"] != "platform" {
			return fmt.Errorf("labels add did not set the tags, got %v", keyResp.Tags)
		}

		_, err = tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/labels")},
			},
		})
		if err != nil {
			return err
		}
		keyResp, err = tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		if err != nil {
			return fmt.Errorf("get api key after remove: %v", err)
		}
		if len(keyResp.Tags) != 0 {
			return fmt.Errorf("labels remove did not clear the tags, got %v", keyResp.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateApiKey_StageAssociation", func() error {
		if err := tc.require(apiKeyID); err != nil {
			return err
		}
		// The /stages row of the official patch table documents add and
		// remove; the value uses the stageKeys member format,
		// restApiId/stageName.
		stage := tc.apiID + "/test"
		_, err := tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/stages"), Value: aws.String(stage)},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		if err != nil {
			return fmt.Errorf("get api key: %v", err)
		}
		found := false
		for _, sk := range resp.StageKeys {
			if sk == stage {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("stage %q missing from stageKeys, got %v", stage, resp.StageKeys)
		}

		// The /stageKeys/ path is not a documented patch form.
		_, err = tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/stageKeys/" + tc.apiID + "~1other")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for the stageKeys path, got: %v", err)
		}

		_, err = tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/stages"), Value: aws.String(stage)},
			},
		})
		if err != nil {
			return err
		}
		resp, err = tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		if err != nil {
			return fmt.Errorf("get api key after remove: %v", err)
		}
		for _, sk := range resp.StageKeys {
			if sk == stage {
				return fmt.Errorf("stage %q still associated after remove, got %v", stage, resp.StageKeys)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteApiKey", func() error {
		if err := tc.require(apiKeyID); err != nil {
			return err
		}
		_, err := tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{
			ApiKey: aws.String(apiKeyID),
		})
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetApiKey should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateApiKey_DefaultEnabled", func() error {
		keyName := tc.uniqueName("default-key")
		resp, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
			Name: aws.String(keyName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{ApiKey: resp.Id})

		if !resp.Enabled {
			return fmt.Errorf("expected enabled=true by default, got false")
		}
		return nil
	}))

	return results
}
