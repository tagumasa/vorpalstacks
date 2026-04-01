package testutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
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

	suffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)

	roleName := fmt.Sprintf("TestRole-%s", suffix)
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
		_ = fmt.Errorf("failed to create test role: %v", err)
	}
	defer func() {
		_, _ = iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{
			RoleName: aws.String(roleName),
		})
	}()

	samlRoleName := fmt.Sprintf("SAMLRole-%s", suffix)
	samlTrustPolicy := `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"Federated": "arn:aws:iam::000000000000:saml-provider/TestProvider"},
			"Action": "sts:AssumeRoleWithSAML"
		}]
	}`

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(samlRoleName),
		AssumeRolePolicyDocument: aws.String(samlTrustPolicy),
	})
	if err != nil {
		_ = fmt.Errorf("failed to create SAML test role: %v", err)
	}
	defer func() {
		_, _ = iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{
			RoleName: aws.String(samlRoleName),
		})
	}()

	webIdRoleName := fmt.Sprintf("WebIdRole-%s", suffix)
	webIdTrustPolicy := `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"Federated": "arn:aws:iam::000000000000:oidc-provider/example.com"},
			"Action": "sts:AssumeRoleWithWebIdentity"
		}]
	}`

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(webIdRoleName),
		AssumeRolePolicyDocument: aws.String(webIdTrustPolicy),
	})
	if err != nil {
		_ = fmt.Errorf("failed to create WebIdentity test role: %v", err)
	}
	defer func() {
		_, _ = iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{
			RoleName: aws.String(webIdRoleName),
		})
	}()

	// === BASIC HAPPY PATH TESTS ===

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
		resp, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
			RoleSessionName: aws.String(roleSessionName),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
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

	// === CONTENT VERIFICATION TESTS ===

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

	results = append(results, r.RunTest("sts", "AssumeRole_ContentVerify", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
		resp, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
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
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
		resp, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
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
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`
		resp, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
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

	results = append(results, r.RunTest("sts", "AssumeRole_InvalidDuration", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
		_, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
			RoleSessionName: aws.String("DurationSession"),
			DurationSeconds: aws.Int32(100),
		})
		if err == nil {
			return fmt.Errorf("expected error for duration < 900")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_EmptySessionName", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
		_, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
			RoleSessionName: aws.String(""),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty session name")
		}
		return nil
	}))

	// === AssumeRoleWithSAML TESTS ===

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_Basic", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", samlRoleName)
		resp, err := client.AssumeRoleWithSAML(ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String(roleARN),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String("VGhpcyBpcyBhIGR1bW15IFNBTUwgYXNzZXJ0aW9u"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_ContentVerify", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", samlRoleName)
		resp, err := client.AssumeRoleWithSAML(ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String(roleARN),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String("VGhpcyBpcyBhIGR1bW15IFNBTUwgYXNzZXJ0aW9u"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil || resp.Credentials.AccessKeyId == nil {
			return fmt.Errorf("credentials or access key ID is nil")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		if resp.AssumedRoleUser == nil || resp.AssumedRoleUser.AssumedRoleId == nil {
			return fmt.Errorf("assumed role user is nil")
		}
		if resp.Subject == nil || *resp.Subject == "" {
			return fmt.Errorf("subject is nil or empty")
		}
		if resp.SubjectType == nil || *resp.SubjectType == "" {
			return fmt.Errorf("subject type is nil or empty")
		}
		if resp.Issuer == nil || *resp.Issuer == "" {
			return fmt.Errorf("issuer is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_WithPolicy", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", samlRoleName)
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`
		resp, err := client.AssumeRoleWithSAML(ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String(roleARN),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String("VGhpcyBpcyBhIGR1bW15IFNBTUwgYXNzZXJ0aW9u"),
			Policy:        aws.String(inlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PackedPolicySize == nil || *resp.PackedPolicySize == 0 {
			return fmt.Errorf("PackedPolicySize should be > 0, got: %v", resp.PackedPolicySize)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_InvalidAssertion", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", samlRoleName)
		_, err := client.AssumeRoleWithSAML(ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String(roleARN),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String(""),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty SAML assertion")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_NonExistentRole", func() error {
		_, err := client.AssumeRoleWithSAML(ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String("arn:aws:iam::000000000000:role/NonExistentSAMLRole"),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String("VGhpcyBpcyBhIGR1bW15IFNBTUwgYXNzZXJ0aW9u"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent role")
		}
		return nil
	}))

	// === AssumeRoleWithWebIdentity TESTS ===

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_Basic", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", webIdRoleName)
		resp, err := client.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(roleARN),
			RoleSessionName:  aws.String("WebIdSession"),
			WebIdentityToken: aws.String("dummy-web-identity-token"),
			ProviderId:       aws.String("example.com"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_ContentVerify", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", webIdRoleName)
		resp, err := client.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(roleARN),
			RoleSessionName:  aws.String("WebIdVerifySession"),
			WebIdentityToken: aws.String("dummy-web-identity-token"),
			ProviderId:       aws.String("example.com"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil || resp.Credentials.AccessKeyId == nil {
			return fmt.Errorf("credentials or access key ID is nil")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		if resp.AssumedRoleUser == nil || resp.AssumedRoleUser.AssumedRoleId == nil {
			return fmt.Errorf("assumed role user is nil")
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
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", webIdRoleName)
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"dynamodb:Query","Resource":"*"}]}`
		resp, err := client.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(roleARN),
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

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_EmptyToken", func() error {
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", webIdRoleName)
		_, err := client.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(roleARN),
			RoleSessionName:  aws.String("WebIdSession"),
			WebIdentityToken: aws.String(""),
			ProviderId:       aws.String("example.com"),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty web identity token")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_NonExistentRole", func() error {
		_, err := client.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String("arn:aws:iam::000000000000:role/NonExistentWebIdRole"),
			RoleSessionName:  aws.String("WebIdSession"),
			WebIdentityToken: aws.String("dummy-token"),
			ProviderId:       aws.String("example.com"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent role")
		}
		return nil
	}))

	// === DecodeAuthorizationMessage TESTS ===

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_Basic", func() error {
		originalMsg := `{"ErrorCode":"AccessDenied","Message":"Not authorized"}`
		encoded := base64.StdEncoding.EncodeToString([]byte(originalMsg))
		resp, err := client.DecodeAuthorizationMessage(ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String(encoded),
		})
		if err != nil {
			return err
		}
		if resp.DecodedMessage == nil || *resp.DecodedMessage != originalMsg {
			return fmt.Errorf("decoded message mismatch, got: %v", resp.DecodedMessage)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_PlainText", func() error {
		originalMsg := "Plain text error message"
		encoded := base64.StdEncoding.EncodeToString([]byte(originalMsg))
		resp, err := client.DecodeAuthorizationMessage(ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String(encoded),
		})
		if err != nil {
			return err
		}
		if resp.DecodedMessage == nil || *resp.DecodedMessage != originalMsg {
			return fmt.Errorf("decoded message mismatch, got: %v", resp.DecodedMessage)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_InvalidBase64", func() error {
		_, err := client.DecodeAuthorizationMessage(ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String("not-valid-base64!!!"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid base64")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_Empty", func() error {
		_, err := client.DecodeAuthorizationMessage(ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String(""),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty encoded message")
		}
		return nil
	}))

	// === GetAccessKeyInfo TESTS ===

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_AKIAPrefix", func() error {
		resp, err := client.GetAccessKeyInfo(ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("AKIAIOSFODNN7EXAMPLE"),
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_ASIAPrefix", func() error {
		resp, err := client.GetAccessKeyInfo(ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("ASIAIOSFODNN7EXAMPLE"),
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_UnknownPrefix", func() error {
		resp, err := client.GetAccessKeyInfo(ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("UNKNOWN1234567890"),
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty for unknown prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_Invalid", func() error {
		_, err := client.GetAccessKeyInfo(ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String(""),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty access key")
		}
		return nil
	}))

	// === GetFederationToken TESTS ===

	results = append(results, r.RunTest("sts", "GetFederationToken_Basic", func() error {
		resp, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name: aws.String("TestFederatedUser"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.FederatedUser == nil {
			return fmt.Errorf("federated user is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_ContentVerify", func() error {
		resp, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name: aws.String("FederatedVerify"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil || resp.Credentials.AccessKeyId == nil {
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
		if resp.FederatedUser == nil {
			return fmt.Errorf("federated user is nil")
		}
		if resp.FederatedUser.FederatedUserId == nil || *resp.FederatedUser.FederatedUserId == "" {
			return fmt.Errorf("federated user ID is nil or empty")
		}
		if resp.FederatedUser.Arn == nil || *resp.FederatedUser.Arn == "" {
			return fmt.Errorf("federated user ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_WithPolicy", func() error {
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:*","Resource":"*"}]}`
		resp, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name:   aws.String("FederatedPolicy"),
			Policy: aws.String(inlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PackedPolicySize == nil || *resp.PackedPolicySize == 0 {
			return fmt.Errorf("PackedPolicySize should be > 0, got: %v", resp.PackedPolicySize)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_InvalidName", func() error {
		_, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name: aws.String(""),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty name")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_InvalidPolicy", func() error {
		_, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name:   aws.String("FederatedBadPolicy"),
			Policy: aws.String("not-valid-json"),
		})
		if err == nil {
			return fmt.Errorf("expected error for malformed policy")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_InvalidDuration", func() error {
		_, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name:            aws.String("FederatedBadDuration"),
			DurationSeconds: aws.Int32(100),
		})
		if err == nil {
			return fmt.Errorf("expected error for duration < 900")
		}
		return nil
	}))

	// === GetDelegatedAccessToken TESTS ===

	results = append(results, r.RunTest("sts", "GetDelegatedAccessToken_Basic", func() error {
		resp, err := client.GetDelegatedAccessToken(ctx, &sts.GetDelegatedAccessTokenInput{
			TradeInToken: aws.String("dummy-trade-in-token"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetDelegatedAccessToken_ContentVerify", func() error {
		resp, err := client.GetDelegatedAccessToken(ctx, &sts.GetDelegatedAccessTokenInput{
			TradeInToken: aws.String("dummy-trade-in-token-verify"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil || resp.Credentials.AccessKeyId == nil {
			return fmt.Errorf("credentials or access key ID is nil")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		if resp.AssumedPrincipal == nil || *resp.AssumedPrincipal == "" {
			return fmt.Errorf("assumed principal is nil or empty")
		}
		if resp.PackedPolicySize == nil {
			return fmt.Errorf("packed policy size is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetDelegatedAccessToken_EmptyToken", func() error {
		_, err := client.GetDelegatedAccessToken(ctx, &sts.GetDelegatedAccessTokenInput{
			TradeInToken: aws.String(""),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty trade-in token")
		}
		return nil
	}))

	// === GetSessionToken EXTENDED DURATION TESTS ===

	results = append(results, r.RunTest("sts", "GetSessionToken_ExtendedDuration", func() error {
		resp, err := client.GetSessionToken(ctx, &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int32(86400),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
	}))

	return results
}
