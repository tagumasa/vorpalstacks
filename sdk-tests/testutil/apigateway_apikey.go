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

		// nameQuery is a substring filter on key names; includeValues
		// controls whether the response carries the key values.
		byName, err := tc.client.GetApiKeys(tc.ctx, &apigateway.GetApiKeysInput{
			NameQuery: aws.String("est-api-ke"),
		})
		if err != nil {
			return err
		}
		if containsID(byName.Items, func(item *types.ApiKey) bool { return item.Id != nil && *item.Id == apiKeyID }) == nil {
			return fmt.Errorf("nameQuery substring did not match the fixture key")
		}
		byMissing, err := tc.client.GetApiKeys(tc.ctx, &apigateway.GetApiKeysInput{
			NameQuery: aws.String("no-such-key-name"),
		})
		if err != nil {
			return err
		}
		if len(byMissing.Items) != 0 {
			return fmt.Errorf("nameQuery with no match returned %d items", len(byMissing.Items))
		}
		withValues, err := tc.client.GetApiKeys(tc.ctx, &apigateway.GetApiKeysInput{
			NameQuery:     aws.String("test-api-key"),
			IncludeValues: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if len(withValues.Items) == 0 || withValues.Items[0].Value == nil || *withValues.Items[0].Value != apiKeyValue {
			return fmt.Errorf("includeValues did not return the key value, got %+v", withValues.Items)
		}
		withoutValues, err := tc.client.GetApiKeys(tc.ctx, &apigateway.GetApiKeysInput{
			NameQuery: aws.String("test-api-key"),
		})
		if err != nil {
			return err
		}
		if len(withoutValues.Items) == 0 || withoutValues.Items[0].Value != nil {
			return fmt.Errorf("values must stay excluded without includeValues, got %+v", withoutValues.Items)
		}

		// customerId matches the key's customer identifier exactly.
		customerResp, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
			Name:       aws.String(tc.uniqueName("customer-key")),
			CustomerId: aws.String("portal-123"),
		})
		if err != nil {
			return fmt.Errorf("create customer key: %v", err)
		}
		defer tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{ApiKey: customerResp.Id})
		byCustomer, err := tc.client.GetApiKeys(tc.ctx, &apigateway.GetApiKeysInput{
			CustomerId: aws.String("portal-123"),
		})
		if err != nil {
			return err
		}
		if len(byCustomer.Items) != 1 || byCustomer.Items[0].Id == nil || *byCustomer.Items[0].Id != *customerResp.Id {
			return fmt.Errorf("customerId filter did not match exactly the customer key, got %+v", byCustomer.Items)
		}
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

		// The remaining scalar rows: description, enabled, and customerId.
		scalarResp, err := tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/description"), Value: aws.String("row coverage")},
				{Op: types.OpReplace, Path: aws.String("/enabled"), Value: aws.String("false")},
				{Op: types.OpReplace, Path: aws.String("/customerId"), Value: aws.String("cust-123")},
			},
		})
		if err != nil {
			return fmt.Errorf("scalar row patch: %v", err)
		}
		if aws.ToString(scalarResp.Description) != "row coverage" {
			return fmt.Errorf("description row not applied, got %v", scalarResp.Description)
		}
		if scalarResp.Enabled {
			return fmt.Errorf("enabled row not applied, got %v", scalarResp.Enabled)
		}
		if aws.ToString(scalarResp.CustomerId) != "cust-123" {
			return fmt.Errorf("customerId row not applied, got %v", scalarResp.CustomerId)
		}
		// Re-enable the key so the later stage-association test exercises
		// an active key.
		if _, err := tc.client.UpdateApiKey(tc.ctx, &apigateway.UpdateApiKeyInput{
			ApiKey: aws.String(apiKeyID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/enabled"), Value: aws.String("true")},
			},
		}); err != nil {
			return fmt.Errorf("re-enable: %v", err)
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

	results = append(results, r.RunTest("apigateway", "ImportApiKeys_CSV", func() error {
		planID, err := tc.createOwnUsagePlan("importplan")
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.deleteUsagePlan(planID)

		keyValue := "importedkey1234abcdefghij0123456789"
		firstImport := "name,key,description,usageplanIds\n" +
			"ImportedFirst," + keyValue + ",first description," + planID + "\n"
		imported, err := tc.client.ImportApiKeys(tc.ctx, &apigateway.ImportApiKeysInput{
			Format: types.ApiKeysFormatCsv,
			Body:   []byte(firstImport),
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{ApiKey: aws.String(keyValue)})
		if len(imported.Ids) != 1 || imported.Ids[0] != keyValue {
			return fmt.Errorf("imported ids = %v, want [%s]", imported.Ids, keyValue)
		}

		stored, err := tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{ApiKey: aws.String(keyValue)})
		if err != nil {
			return fmt.Errorf("get imported key: %v", err)
		}
		if aws.ToString(stored.Name) != "ImportedFirst" || aws.ToString(stored.Description) != "first description" {
			return fmt.Errorf("imported key payload mismatch: %+v", stored)
		}
		keysResp, err := tc.client.GetUsagePlanKeys(tc.ctx, &apigateway.GetUsagePlanKeysInput{
			UsagePlanId: aws.String(planID),
		})
		if err != nil {
			return fmt.Errorf("get usage plan keys: %v", err)
		}
		associated := false
		for _, k := range keysResp.Items {
			if aws.ToString(k.Id) == keyValue {
				associated = true
			}
		}
		if !associated {
			return fmt.Errorf("imported key not associated with the usage plan")
		}

		// Re-importing the same key value (re-referencing the plan)
		// overwrites the stored key and reports the existing plan
		// association as a warning.
		secondImport := "name,key,description,usageplanIds\nImportedRenamed," + keyValue + ",second description," + planID + "\n"
		reimported, err := tc.client.ImportApiKeys(tc.ctx, &apigateway.ImportApiKeysInput{
			Format: types.ApiKeysFormatCsv,
			Body:   []byte(secondImport),
		})
		if err != nil {
			return err
		}
		if len(reimported.Ids) != 1 || len(reimported.Warnings) != 1 {
			return fmt.Errorf("re-import ids/warnings = %v / %v", reimported.Ids, reimported.Warnings)
		}
		renamed, err := tc.client.GetApiKey(tc.ctx, &apigateway.GetApiKeyInput{ApiKey: aws.String(keyValue)})
		if err != nil {
			return err
		}
		if aws.ToString(renamed.Name) != "ImportedRenamed" || aws.ToString(renamed.Description) != "second description" {
			return fmt.Errorf("re-import did not overwrite: %+v", renamed)
		}

		// With failonwarnings, an invalid row rejects the whole import.
		_, err = tc.client.ImportApiKeys(tc.ctx, &apigateway.ImportApiKeysInput{
			Format:         types.ApiKeysFormatCsv,
			FailOnWarnings: true,
			Body:           []byte("name,key\nGood," + keyValue + "\nBad,short\n"),
		})
		if aerr := AssertErrorContains(err, "BadRequestException"); aerr != nil {
			return fmt.Errorf("expected BadRequestException with failonwarnings, got: %v", aerr)
		}
		return nil
	}))

	return results
}
