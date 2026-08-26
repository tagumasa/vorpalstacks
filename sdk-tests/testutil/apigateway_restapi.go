package testutil

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayRestApiTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "GetRestApis", func() error {
		// Walk every page — accumulated APIs from previous runs can push
		// the shared API beyond the first page.
		items, err := tc.allRestApis()
		if err != nil {
			return err
		}
		if items == nil {
			return fmt.Errorf("items list is nil")
		}
		for _, item := range items {
			if item.Name != nil && *item.Name == tc.apiName {
				return nil
			}
		}
		return fmt.Errorf("API not found")
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApi", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String(tc.apiID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != tc.apiName {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Id == nil || *resp.Id != tc.apiID {
			return fmt.Errorf("id mismatch, got %v", resp.Id)
		}
		if resp.CreatedDate == nil {
			return fmt.Errorf("createdDate is nil")
		}
		if resp.MinimumCompressionSize != nil {
			return fmt.Errorf("minimumCompressionSize should be nil when unset, got %d", *resp.MinimumCompressionSize)
		}
		if resp.Version != nil {
			return fmt.Errorf("version should be nil when unset, got %v", *resp.Version)
		}
		if resp.ApiKeySource != "" {
			return fmt.Errorf("apiKeySource should be empty when unset, got %v", resp.ApiKeySource)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.UpdateRestApi(tc.ctx, &apigateway.UpdateRestApiInput{
			RestApiId: aws.String(tc.apiID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("Updated API"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Description == nil || *resp.Description != "Updated API" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}

		// Verify the update persists via a fresh read of the shared API.
		getResp, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String(tc.apiID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Description == nil || *getResp.Description != "Updated API" {
			return fmt.Errorf("description not updated, got %v", getResp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateRestApi_WithPolicy", func() error {
		policyDoc := map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{"Effect": "Allow", "Principal": "*", "Action": "execute-api:Invoke", "Resource": "*"},
			},
		}
		policyBytes, _ := json.Marshal(policyDoc)

		createResp, err := tc.client.CreateRestApi(tc.ctx, &apigateway.CreateRestApiInput{
			Name:   aws.String(tc.uniqueName("PolAPI")),
			Policy: aws.String(string(policyBytes)),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		apiID := aws.ToString(createResp.Id)
		defer tc.deleteAPI(apiID)

		if createResp.Policy == nil || *createResp.Policy == "" {
			return fmt.Errorf("policy not returned in create response")
		}

		resp, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy not returned in get response")
		}
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(*resp.Policy), &parsed); err != nil {
			return fmt.Errorf("policy is not valid JSON: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi_BinaryMediaTypes", func() error {
		apiID, _, err := tc.createAPI(tc.uniqueName("BmAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(apiID)

		_, err = tc.client.UpdateRestApi(tc.ctx, &apigateway.UpdateRestApiInput{
			RestApiId: aws.String(apiID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/binaryMediaTypes"), Value: aws.String("image/png")},
				{Op: types.OpAdd, Path: aws.String("/binaryMediaTypes"), Value: aws.String("application/octet-stream")},
			},
		})
		if err != nil {
			return fmt.Errorf("add binaryMediaTypes: %v", err)
		}

		resp, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.BinaryMediaTypes) != 2 {
			return fmt.Errorf("expected 2 binaryMediaTypes, got %d", len(resp.BinaryMediaTypes))
		}
		found := map[string]bool{}
		for _, mt := range resp.BinaryMediaTypes {
			found[mt] = true
		}
		if !found["image/png"] || !found["application/octet-stream"] {
			return fmt.Errorf("binaryMediaTypes mismatch, got %v", resp.BinaryMediaTypes)
		}

		_, err = tc.client.UpdateRestApi(tc.ctx, &apigateway.UpdateRestApiInput{
			RestApiId: aws.String(apiID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/binaryMediaTypes/image~1png")},
			},
		})
		if err != nil {
			return fmt.Errorf("remove binaryMediaType: %v", err)
		}

		resp2, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return fmt.Errorf("get after remove: %v", err)
		}
		if len(resp2.BinaryMediaTypes) != 1 || resp2.BinaryMediaTypes[0] != "application/octet-stream" {
			return fmt.Errorf("expected 1 binaryMediaType after remove, got %v", resp2.BinaryMediaTypes)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi_MinimumCompressionSize", func() error {
		apiID, _, err := tc.createAPI(tc.uniqueName("McAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(apiID)

		getBefore, _ := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{RestApiId: aws.String(apiID)})
		if getBefore.MinimumCompressionSize != nil {
			return fmt.Errorf("minimumCompressionSize should be nil before setting, got %d", *getBefore.MinimumCompressionSize)
		}

		_, err = tc.client.UpdateRestApi(tc.ctx, &apigateway.UpdateRestApiInput{
			RestApiId: aws.String(apiID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/minimumCompressionSize"), Value: aws.String("2048")},
			},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.MinimumCompressionSize == nil || *resp.MinimumCompressionSize != 2048 {
			return fmt.Errorf("minimumCompressionSize mismatch, got %v", resp.MinimumCompressionSize)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApis_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgAPIs []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagAPI-%s-%d", pgTs, i)
			apiID, _, err := tc.createAPI(name)
			if err != nil {
				return fmt.Errorf("create rest api %s: %v", name, err)
			}
			defer tc.deleteAPI(apiID)
			pgAPIs = append(pgAPIs, apiID)
		}

		allAPIs, err := tc.allRestApis()
		if err != nil {
			return fmt.Errorf("get rest apis page: %v", err)
		}
		count := 0
		for _, item := range allAPIs {
			if item.Name != nil && strings.HasPrefix(*item.Name, "PagAPI-"+pgTs) {
				count++
			}
		}
		if len(pgAPIs) != 5 || count != 5 {
			return fmt.Errorf("expected 5 paginated rest apis, got %d", count)
		}
		return nil
	}))

	return results
}
