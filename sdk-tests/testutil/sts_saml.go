package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (r *TestRunner) runSTSSAMLTests(tc *stsTestContext) []TestResult {
	var results []TestResult

	samlAssertion := "VGhpcyBpcyBhIGR1bW15IFNBTUwgYXNzZXJ0aW9u"

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_ContentVerify", func() error {
		resp, err := tc.client.AssumeRoleWithSAML(tc.ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String(tc.samlRoleARN()),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String(samlAssertion),
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
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`
		resp, err := tc.client.AssumeRoleWithSAML(tc.ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String(tc.samlRoleARN()),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String(samlAssertion),
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

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_WithDuration", func() error {
		resp, err := tc.client.AssumeRoleWithSAML(tc.ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:         aws.String(tc.samlRoleARN()),
			PrincipalArn:    aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion:   aws.String(samlAssertion),
			DurationSeconds: aws.Int32(3600),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil || resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("credentials or expiration invalid")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_InvalidAssertion", func() error {
		_, err := tc.client.AssumeRoleWithSAML(tc.ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String(tc.samlRoleARN()),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String(""),
		})
		if err := AssertErrorContains(err, "InvalidSAMLAssertion"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithSAML_NonExistentRole", func() error {
		_, err := tc.client.AssumeRoleWithSAML(tc.ctx, &sts.AssumeRoleWithSAMLInput{
			RoleArn:       aws.String("arn:aws:iam::000000000000:role/NonExistentSAMLRole"),
			PrincipalArn:  aws.String("arn:aws:iam::000000000000:saml-provider/TestProvider"),
			SAMLAssertion: aws.String(samlAssertion),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
