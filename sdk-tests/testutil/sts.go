package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSTSTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "sts",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := sts.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)
	ctx := context.Background()

	// Create test role for AssumeRole test
	roleName := fmt.Sprintf("TestRole-%d", time.Now().UnixNano()%1000000)
	trustPolicy := `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"AWS": "arn:aws:iam::000000000000:root"},
			"Action": "sts:AssumeRole"
		}]
	}`

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
	})
	if err != nil {
		var alreadyExists *types.EntityAlreadyExistsException
		if err.Error() != "EntityAlreadyExists" && err.Error() != "" {
			_ = alreadyExists
		}
	}
	defer func() {
		_, _ = iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{
			RoleName: aws.String(roleName),
		})
	}()

	results = append(results, r.RunTest("sts", "GetCallerIdentity", func() error {
		resp, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return err
		}
		if resp.UserId == nil {
			return fmt.Errorf("user ID is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetSessionToken", func() error {
		resp, err := client.GetSessionToken(ctx, &sts.GetSessionTokenInput{})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
		roleSessionName := "TestSession"
		_, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
			RoleSessionName: aws.String(roleSessionName),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("sts", "AssumeRole_NonExistentRole", func() error {
		_, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/NonExistentRole"),
			RoleSessionName: aws.String("TestSession"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent role")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetCallerIdentity_ContentVerify", func() error {
		resp, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty")
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("ARN is nil or empty")
		}
		if resp.UserId == nil || *resp.UserId == "" {
			return fmt.Errorf("user ID is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetSessionToken_ContentVerify", func() error {
		resp, err := client.GetSessionToken(ctx, &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int32(3600),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.Credentials.AccessKeyId == nil || *resp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("access key ID is nil or empty")
		}
		if resp.Credentials.SecretAccessKey == nil || *resp.Credentials.SecretAccessKey == "" {
			return fmt.Errorf("secret access key is nil or empty")
		}
		if resp.Credentials.SessionToken == nil || *resp.Credentials.SessionToken == "" {
			return fmt.Errorf("session token is nil or empty")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		return nil
	}))

	return results
}
