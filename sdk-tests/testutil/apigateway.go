package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunAPIGatewayTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "apigateway",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := apigateway.NewFromConfig(cfg)
	ctx := context.Background()

	apiName := fmt.Sprintf("TestAPI-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("apigateway", "CreateRestApi", func() error {
		_, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(apiName),
			Description: aws.String("Test API"),
		})
		return err
	}))

	var apiID string
	results = append(results, r.RunTest("apigateway", "GetRestApis", func() error {
		resp, err := client.GetRestApis(ctx, &apigateway.GetRestApisInput{
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil {
			return fmt.Errorf("items list is nil")
		}
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == apiName {
				apiID = *item.Id
				break
			}
		}
		if apiID == "" {
			return fmt.Errorf("API not found")
		}
		return nil
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
		if resp.Name == nil {
			return fmt.Errorf("API name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateRestApi", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		_, err := client.UpdateRestApi(ctx, &apigateway.UpdateRestApiInput{
			RestApiId: aws.String(apiID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("Updated API"),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "CreateResource", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(apiID),
			ParentId:  aws.String(apiID),
			PathPart:  aws.String("test"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("resource ID is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetResources", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetResources(ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil {
			return fmt.Errorf("items list is nil")
		}
		return nil
	}))

	var deploymentID string
	results = append(results, r.RunTest("apigateway", "CreateDeployment", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Id != nil {
			deploymentID = *resp.Id
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDeployments", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetDeployments(ctx, &apigateway.GetDeploymentsInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil {
			return fmt.Errorf("items list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateStage", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		if deploymentID == "" {
			return fmt.Errorf("Deployment ID not available")
		}
		_, err := client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(apiID),
			StageName:    aws.String("test"),
			DeploymentId: aws.String(deploymentID),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "GetStage", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return err
		}
		if resp.StageName == nil {
			return fmt.Errorf("stage name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStages", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetStages(ctx, &apigateway.GetStagesInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Item == nil {
			return fmt.Errorf("stage item is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRestApi", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		_, err := client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
			RestApiId: aws.String(apiID),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("apigateway", "GetRestApi_NonExistent", func() error {
		_, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent API")
		}
		var nf *types.NotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRestApi_NonExistent", func() error {
		_, err := client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent API")
		}
		var nf *types.NotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStage_NonExistent", func() error {
		tmpAPI := fmt.Sprintf("TmpAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tmpAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("nonexistent_stage"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent stage")
		}
		var nf *types.NotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
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

	results = append(results, r.RunTest("apigateway", "CreateResource_NestedPath", func() error {
		crAPI := fmt.Sprintf("CrAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(crAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		usersResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("users"),
		})
		if err != nil {
			return fmt.Errorf("create users resource: %v", err)
		}

		userIdResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  usersResp.Id,
			PathPart:  aws.String("{userId}"),
		})
		if err != nil {
			return fmt.Errorf("create userId resource: %v", err)
		}
		if userIdResp.Path == nil || *userIdResp.Path != "/users/{userId}" {
			return fmt.Errorf("nested path mismatch, got %v", userIdResp.Path)
		}

		resResp, err := client.GetResources(ctx, &apigateway.GetResourcesInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		if len(resResp.Items) < 3 {
			return fmt.Errorf("expected at least 3 resources (root, users, {userId}), got %d", len(resResp.Items))
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateStage_VerifyConfig", func() error {
		csAPI := fmt.Sprintf("CsAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(csAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		depResp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		stageDesc := "test stage description"
		_, err = client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    createResp.Id,
			StageName:    aws.String("v1"),
			DeploymentId: depResp.Id,
			Description:  aws.String(stageDesc),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		resp, err := client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if resp.Description == nil || *resp.Description != stageDesc {
			return fmt.Errorf("stage description mismatch, got %v", resp.Description)
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != *depResp.Id {
			return fmt.Errorf("deployment ID mismatch, got %v", resp.DeploymentId)
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

	return results
}
