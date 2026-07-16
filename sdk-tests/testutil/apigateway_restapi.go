package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayRestApiTests(ctx context.Context, client *apigateway.Client, apiID, apiName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "GetRestApis", func() error {
		resp, err := client.GetRestApis(ctx, &apigateway.GetRestApisInput{
			Limit: aws.Int32(500),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil {
			return fmt.Errorf("items list is nil")
		}
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == apiName {
				return nil
			}
		}
		return fmt.Errorf("API not found")
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApi", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != apiName {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Id == nil || *resp.Id != apiID {
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
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: aws.String(apiID),
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
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi_VerifyUpdate", func() error {
		uaAPI := fmt.Sprintf("UaAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(uaAPI),
			Description: aws.String("original desc"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		newDesc := "updated description v2"
		_, err = client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: createResp.Id,
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String(newDesc),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Description == nil || *resp.Description != newDesc {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApis_ContainsCreated", func() error {
		gaAPI := fmt.Sprintf("GaAPI-%d", time.Now().UnixNano())
		gaDesc := "searchable description"
		_, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(gaAPI),
			Description: aws.String(gaDesc),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: func() *string {
			resp, _ := client.GetRestApis(ctx, &apigateway.GetRestApisInput{Limit: aws.Int32(500)})
			if resp != nil {
				for _, item := range resp.Items {
					if item.Name != nil && *item.Name == gaAPI {
						return item.Id
					}
				}
			}
			return nil
		}()})

		resp, err := client.GetRestApis(ctx, &apigateway.GetRestApisInput{
			Limit: aws.Int32(500),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == gaAPI {
				found = true
				if item.Description == nil || *item.Description != gaDesc {
					return fmt.Errorf("description mismatch in list, got %v", item.Description)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created API %s not found in GetRestApis", gaAPI)
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
		polAPI := fmt.Sprintf("PolAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:   aws.String(polAPI),
			Policy: aws.String(string(policyBytes)),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		if createResp.Policy == nil || *createResp.Policy == "" {
			return fmt.Errorf("policy not returned in create response")
		}

		getResp, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Policy == nil || *getResp.Policy == "" {
			return fmt.Errorf("policy not returned in get response")
		}
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(*getResp.Policy), &parsed); err != nil {
			return fmt.Errorf("policy is not valid JSON: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi_BinaryMediaTypes", func() error {
		bmAPI := fmt.Sprintf("BmAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(bmAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: createResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/binaryMediaTypes"), Value: aws.String("image/png")},
				{Op: types.OpAdd, Path: aws.String("/binaryMediaTypes"), Value: aws.String("application/octet-stream")},
			},
		})
		if err != nil {
			return fmt.Errorf("add binaryMediaTypes: %v", err)
		}

		resp, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
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

		_, err = client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: createResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/binaryMediaTypes/image~1png")},
			},
		})
		if err != nil {
			return fmt.Errorf("remove binaryMediaType: %v", err)
		}

		resp2, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
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
		mcAPI := fmt.Sprintf("McAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(mcAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		getBefore, _ := client.GetRestApi(ctx, &apigateway.GetRestApiInput{RestApiId: createResp.Id})
		if getBefore.MinimumCompressionSize != nil {
			return fmt.Errorf("minimumCompressionSize should be nil before setting, got %d", *getBefore.MinimumCompressionSize)
		}

		_, err = client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: createResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/minimumCompressionSize"), Value: aws.String("2048")},
			},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
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
			resp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
				Name:        aws.String(name),
				Description: aws.String("pagination test"),
			})
			if err != nil {
				for _, id := range pgAPIs {
					client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(id)})
				}
				return fmt.Errorf("create rest api %s: %v", name, err)
			}
			pgAPIs = append(pgAPIs, *resp.Id)
		}

		var allAPIs []string
		var position *string
		for {
			resp, err := client.GetRestApis(ctx, &apigateway.GetRestApisInput{
				Limit:    aws.Int32(2),
				Position: position,
			})
			if err != nil {
				for _, id := range pgAPIs {
					client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(id)})
				}
				return fmt.Errorf("get rest apis page: %v", err)
			}
			for _, item := range resp.Items {
				if item.Name != nil && strings.HasPrefix(*item.Name, "PagAPI-"+pgTs) {
					allAPIs = append(allAPIs, *item.Name)
				}
			}
			if resp.Position != nil && *resp.Position != "" {
				position = resp.Position
			} else {
				break
			}
		}

		for _, id := range pgAPIs {
			client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(id)})
		}
		if len(allAPIs) != 5 {
			return fmt.Errorf("expected 5 paginated rest apis, got %d", len(allAPIs))
		}
		return nil
	}))

	return results
}
