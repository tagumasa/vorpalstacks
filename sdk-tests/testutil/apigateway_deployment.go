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

func (r *TestRunner) runAPIGatewayDeploymentTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	var deploymentID string
	results = append(results, r.RunTest("apigateway", "CreateDeployment", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId:   aws.String(apiID),
			Description: aws.String("test deployment"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("deployment ID is nil")
		}
		if resp.Description == nil || *resp.Description != "test deployment" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		deploymentID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDeployment", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := client.GetDeployment(ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != deploymentID {
			return fmt.Errorf("deployment ID mismatch, got %v", resp.Id)
		}
		if resp.CreatedDate == nil {
			return fmt.Errorf("createdDate is nil")
		}
		if resp.Description == nil || *resp.Description != "test deployment" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateDeployment", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := client.UpdateDeployment(ctx, &apigateway.UpdateDeploymentInput{
			RestApiId:    aws.String(apiID),
			DeploymentId: aws.String(deploymentID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated deployment"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated deployment" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
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
		if len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 deployment")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateStage", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(apiID),
			StageName:    aws.String("test"),
			DeploymentId: aws.String(deploymentID),
			Description:  aws.String("test stage"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.StageName == nil || *resp.StageName != "test" {
			return fmt.Errorf("stageName mismatch, got %v", resp.StageName)
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != deploymentID {
			return fmt.Errorf("deploymentId mismatch, got %v", resp.DeploymentId)
		}
		if resp.Description == nil || *resp.Description != "test stage" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
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
		if resp.StageName == nil || *resp.StageName != "test" {
			return fmt.Errorf("stage name mismatch, got %v", resp.StageName)
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != deploymentID {
			return fmt.Errorf("deploymentId mismatch, got %v", resp.DeploymentId)
		}
		if resp.Description == nil || *resp.Description != "test stage" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
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
		if len(resp.Item) == 0 {
			return fmt.Errorf("expected at least 1 stage")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.UpdateStage(ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated stage"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated stage" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteStage", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		_, err := client.DeleteStage(ctx, &apigateway.DeleteStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("test"),
		})
		if err == nil {
			return fmt.Errorf("GetStage should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteDeployment", func() error {
		if apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		_, err := client.DeleteDeployment(ctx, &apigateway.DeleteDeploymentInput{
			RestApiId:    aws.String(apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetDeployment(ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err == nil {
			return fmt.Errorf("GetDeployment should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
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

	results = append(results, r.RunTest("apigateway", "Deployment_Stage_FullLifecycle", func() error {
		dsAPI := fmt.Sprintf("DsAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(dsAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		depResp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId:   createResp.Id,
			Description: aws.String("v1 deployment"),
		})
		if err != nil {
			return fmt.Errorf("create deployment: %v", err)
		}

		getDepResp, err := client.GetDeployment(ctx, &apigateway.GetDeploymentInput{
			RestApiId:    createResp.Id,
			DeploymentId: depResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get deployment: %v", err)
		}
		if getDepResp.Description == nil || *getDepResp.Description != "v1 deployment" {
			return fmt.Errorf("deployment description mismatch, got %v", getDepResp.Description)
		}

		_, err = client.UpdateDeployment(ctx, &apigateway.UpdateDeploymentInput{
			RestApiId:    createResp.Id,
			DeploymentId: depResp.Id,
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("v1 updated"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update deployment: %v", err)
		}

		_, err = client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    createResp.Id,
			StageName:    aws.String("production"),
			DeploymentId: depResp.Id,
			Description:  aws.String("production stage"),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		_, err = client.UpdateStage(ctx, &apigateway.UpdateStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("production"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("production updated"),
				},
				{
					Op:    types.OpReplace,
					Path:  aws.String("/variables/env"),
					Value: aws.String("prod"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update stage: %v", err)
		}

		stageResp, err := client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("production"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if stageResp.Variables == nil || stageResp.Variables["env"] != "prod" {
			return fmt.Errorf("stage variables not set, got %v", stageResp.Variables)
		}

		_, err = client.DeleteStage(ctx, &apigateway.DeleteStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("production"),
		})
		if err != nil {
			return fmt.Errorf("delete stage: %v", err)
		}

		_, err = client.DeleteDeployment(ctx, &apigateway.DeleteDeploymentInput{
			RestApiId:    createResp.Id,
			DeploymentId: depResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete deployment: %v", err)
		}
		return nil
	}))

	return results
}
