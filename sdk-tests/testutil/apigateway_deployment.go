package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayDeploymentTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var deploymentID string
	results = append(results, r.RunTest("apigateway", "CreateDeployment", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.CreateDeployment(tc.ctx, &apigateway.CreateDeploymentInput{
			RestApiId:   aws.String(tc.apiID),
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
		if tc.apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := tc.client.GetDeployment(tc.ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
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
		if tc.apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := tc.client.UpdateDeployment(tc.ctx, &apigateway.UpdateDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
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

		// Verify the description change persists via a fresh read.
		getResp, err := tc.client.GetDeployment(tc.ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return fmt.Errorf("get deployment: %v", err)
		}
		if getResp.Description == nil || *getResp.Description != "updated deployment" {
			return fmt.Errorf("deployment description mismatch, got %v", getResp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDeployments", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.GetDeployments(tc.ctx, &apigateway.GetDeploymentsInput{
			RestApiId: aws.String(tc.apiID),
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
		if tc.apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		resp, err := tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(tc.apiID),
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
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
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
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.GetStages(tc.ctx, &apigateway.GetStagesInput{
			RestApiId: aws.String(tc.apiID),
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
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated stage"),
				},
				{
					Op:    types.OpReplace,
					Path:  aws.String("/variables/env"),
					Value: aws.String("prod"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated stage" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}

		// Verify the stage variable round-trips via a fresh read.
		stageResp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if stageResp.Variables == nil || stageResp.Variables["env"] != "prod" {
			return fmt.Errorf("stage variables not set, got %v", stageResp.Variables)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteStage", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		_, err := tc.client.DeleteStage(tc.ctx, &apigateway.DeleteStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
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
		if tc.apiID == "" || deploymentID == "" {
			return fmt.Errorf("API ID or deployment ID not available")
		}
		_, err := tc.client.DeleteDeployment(tc.ctx, &apigateway.DeleteDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetDeployment(tc.ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
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
		apiID, _, err := tc.createAPI(tc.uniqueName("CsAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(apiID)

		depID, err := tc.createDeployment(apiID, "")
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		stageDesc := "test stage description"
		_, err = tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(apiID),
			StageName:    aws.String("v1"),
			DeploymentId: aws.String(depID),
			Description:  aws.String(stageDesc),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		resp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if resp.Description == nil || *resp.Description != stageDesc {
			return fmt.Errorf("stage description mismatch, got %v", resp.Description)
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != depID {
			return fmt.Errorf("deployment ID mismatch, got %v", resp.DeploymentId)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_AccessLogDestinationArn_Validated", func() error {
		apiID, _, err := tc.createAPI(tc.uniqueName("AlAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(apiID)

		depID, err := tc.createDeployment(apiID, "")
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		_, err = tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(apiID),
			StageName:    aws.String("access-log"),
			DeploymentId: aws.String(depID),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		// A Lambda function ARN is not a valid access log destination.
		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("access-log"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/accessLogSettings/destinationArn"),
					Value: aws.String("arn:aws:lambda:us-east-1:123456789012:function:logs")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for Lambda destination ARN, got: %v", err)
		}

		// Firehose delivery streams must begin with amazon-apigateway-.
		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("access-log"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/accessLogSettings/destinationArn"),
					Value: aws.String("arn:aws:firehose:us-east-1:123456789012:deliverystream/my-stream")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for non-prefixed Firehose ARN, got: %v", err)
		}

		// A CloudWatch Logs log group ARN is accepted.
		stageResp, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("access-log"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/accessLogSettings/destinationArn"),
					Value: aws.String("arn:aws:logs:us-east-1:123456789012:log-group:my-api-logs")},
			},
		})
		if err != nil {
			return fmt.Errorf("update stage with valid logs ARN: %v", err)
		}
		if stageResp.AccessLogSettings == nil || stageResp.AccessLogSettings.DestinationArn == nil ||
			*stageResp.AccessLogSettings.DestinationArn != "arn:aws:logs:us-east-1:123456789012:log-group:my-api-logs" {
			return fmt.Errorf("accessLogSettings destinationArn not stored, got %v", stageResp.AccessLogSettings)
		}
		return nil
	}))

	return results
}
