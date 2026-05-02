package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

func (r *TestRunner) runSTSAssumeTests(tc *stsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sts", "AssumeRole_ContentVerify", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("VerifySession"),
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
		if resp.AssumedRoleUser == nil {
			return fmt.Errorf("assumed role user is nil")
		}
		if resp.AssumedRoleUser.AssumedRoleId == nil || *resp.AssumedRoleUser.AssumedRoleId == "" {
			return fmt.Errorf("assumed role ID is nil or empty")
		}
		if resp.AssumedRoleUser.Arn == nil || *resp.AssumedRoleUser.Arn == "" {
			return fmt.Errorf("assumed role user ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithSourceIdentity", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("SourceIdSession"),
			SourceIdentity:  aws.String("AdminUser"),
		})
		if err != nil {
			return err
		}
		if resp.SourceIdentity == nil || *resp.SourceIdentity != "AdminUser" {
			return fmt.Errorf("SourceIdentity not returned correctly, got: %v", resp.SourceIdentity)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithPolicy", func() error {
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("PolicySession"),
			Policy:          aws.String(inlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PackedPolicySize == nil || *resp.PackedPolicySize == 0 {
			return fmt.Errorf("PackedPolicySize should be > 0, got: %v", resp.PackedPolicySize)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithDurationSeconds", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("DurationSession"),
			DurationSeconds: aws.Int32(3600),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_NonExistentRole", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/NonExistentRole"),
			RoleSessionName: aws.String("TestSession"),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_InvalidDuration", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("DurationSession"),
			DurationSeconds: aws.Int32(100),
		})
		if err := AssertErrorContains(err, "InvalidDuration"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_EmptySessionName", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String(""),
		})
		if err := AssertErrorContains(err, "ValidationError"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_InvalidArn", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String("not-a-valid-arn"),
			RoleSessionName: aws.String("TestSession"),
		})
		if err := AssertErrorContains(err, "InvalidRoleArn"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoot_Basic", func() error {
		resp, err := tc.client.AssumeRoot(tc.ctx, &sts.AssumeRootInput{
			TargetPrincipal: aws.String("000000000000"),
			TaskPolicyArn: &types.PolicyDescriptorType{
				Arn: aws.String("arn:aws:iam::aws:policy/IAMAuditRootUserCredentials"),
			},
			DurationSeconds: aws.Int32(900),
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
