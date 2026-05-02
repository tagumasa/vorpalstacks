package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (r *TestRunner) runSTSWebIdentityTests(tc *stsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_ContentVerify", func() error {
		resp, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(tc.webIdRoleARN()),
			RoleSessionName:  aws.String("WebIdVerifySession"),
			WebIdentityToken: aws.String("dummy-web-identity-token"),
			ProviderId:       aws.String("example.com"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil || resp.Credentials.AccessKeyId == nil || *resp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("credentials or access key ID is nil")
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
		if resp.AssumedRoleUser == nil || resp.AssumedRoleUser.AssumedRoleId == nil {
			return fmt.Errorf("assumed role user is nil")
		}
		if resp.AssumedRoleUser.Arn == nil || *resp.AssumedRoleUser.Arn == "" {
			return fmt.Errorf("assumed role user ARN is nil or empty")
		}
		if resp.SubjectFromWebIdentityToken == nil || *resp.SubjectFromWebIdentityToken == "" {
			return fmt.Errorf("subject from web identity token is nil or empty")
		}
		if resp.Audience == nil || *resp.Audience == "" {
			return fmt.Errorf("audience is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_WithPolicy", func() error {
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"dynamodb:Query","Resource":"*"}]}`
		resp, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(tc.webIdRoleARN()),
			RoleSessionName:  aws.String("WebIdPolicySession"),
			WebIdentityToken: aws.String("dummy-web-identity-token"),
			ProviderId:       aws.String("example.com"),
			Policy:           aws.String(inlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PackedPolicySize == nil || *resp.PackedPolicySize == 0 {
			return fmt.Errorf("PackedPolicySize should be > 0, got: %v", resp.PackedPolicySize)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_WithDuration", func() error {
		resp, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(tc.webIdRoleARN()),
			RoleSessionName:  aws.String("WebIdDurationSession"),
			WebIdentityToken: aws.String("dummy-web-identity-token"),
			ProviderId:       aws.String("example.com"),
			DurationSeconds:  aws.Int32(3600),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil || resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("credentials or expiration invalid")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_EmptyToken", func() error {
		_, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(tc.webIdRoleARN()),
			RoleSessionName:  aws.String("WebIdSession"),
			WebIdentityToken: aws.String(""),
			ProviderId:       aws.String("example.com"),
		})
		if err := AssertErrorContains(err, "InvalidIdentityToken"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_NonExistentRole", func() error {
		_, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String("arn:aws:iam::000000000000:role/NonExistentWebIdRole"),
			RoleSessionName:  aws.String("WebIdSession"),
			WebIdentityToken: aws.String("dummy-token"),
			ProviderId:       aws.String("example.com"),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetWebIdentityToken_Basic", func() error {
		resp, err := tc.client.GetWebIdentityToken(tc.ctx, &sts.GetWebIdentityTokenInput{
			Audience:         []string{"sts.amazonaws.com"},
			SigningAlgorithm: aws.String("RS256"),
		})
		if err != nil {
			return err
		}
		if resp.WebIdentityToken == nil || *resp.WebIdentityToken == "" {
			return fmt.Errorf("web identity token is nil or empty")
		}
		if resp.Expiration == nil || resp.Expiration.IsZero() {
			return fmt.Errorf("expiration is nil or zero")
		}
		return nil
	}))

	return results
}
