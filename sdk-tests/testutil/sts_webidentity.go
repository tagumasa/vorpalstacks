package testutil

import (
	"encoding/base64"
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
		if err := stsAssertCredentials(resp.Credentials); err != nil {
			return err
		}
		if err := stsAssertAssumedRoleUser(resp.AssumedRoleUser); err != nil {
			return err
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
		if err := stsAssertPackedPolicySize(resp.PackedPolicySize); err != nil {
			return err
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
			RoleArn:          aws.String(fmt.Sprintf("arn:aws:iam::%s:role/NonExistentWebIdRole", tc.accountID)),
			RoleSessionName:  aws.String("WebIdSession"),
			WebIdentityToken: aws.String("dummy-token"),
			ProviderId:       aws.String("example.com"),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_ShortProviderId", func() error {
		// ProviderId shorter than urlType min (4 chars) should return
		// ValidationError, not InvalidIdentityToken.
		_, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(tc.webIdRoleARN()),
			RoleSessionName:  aws.String("WebIdShortProvider"),
			WebIdentityToken: aws.String("dummy-web-identity-token"),
			ProviderId:       aws.String("ab"),
		})
		if err := AssertErrorContains(err, "ValidationError"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_ShortSubjectFallback", func() error {
		// roleSessionName "ab" (2 chars) is valid for roleSessionNameType
		// (min 2) but below webIdentitySubjectType min (6).  The server
		// must pad the fallback Subject to meet the constraint.
		resp, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(tc.webIdRoleARN()),
			RoleSessionName:  aws.String("ab"),
			WebIdentityToken: aws.String("dummy-web-identity-token"),
			ProviderId:       aws.String("example.com"),
		})
		if err != nil {
			return err
		}
		if resp.SubjectFromWebIdentityToken == nil || *resp.SubjectFromWebIdentityToken == "" {
			return fmt.Errorf("subject is nil or empty")
		}
		if len(*resp.SubjectFromWebIdentityToken) < 6 {
			return fmt.Errorf("subject too short: %d chars (min 6), got: %s",
				len(*resp.SubjectFromWebIdentityToken), *resp.SubjectFromWebIdentityToken)
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

	// When the WebIdentityToken is a parseable JWT with an iss claim,
	// the Provider field in the response must contain the iss value,
	// not the caller-supplied ProviderId.
	results = append(results, r.RunTest("sts", "AssumeRoleWithWebIdentity_ProviderFromISS", func() error {
		// Build a minimal JWT with iss/sub/aud claims.
		payload := `{"iss":"https://test.oidc.example.com","sub":"test-subject-123456","aud":"test-audience","exp":9999999999}`
		encodedPayload := base64.RawURLEncoding.EncodeToString([]byte(payload))
		jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9." + encodedPayload + ".sig"

		resp, err := tc.client.AssumeRoleWithWebIdentity(tc.ctx, &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(tc.webIdRoleARN()),
			RoleSessionName:  aws.String("IssProviderSession"),
			WebIdentityToken: aws.String(jwtToken),
			ProviderId:       aws.String("example.com"),
		})
		if err != nil {
			return err
		}
		if resp.Provider == nil {
			return fmt.Errorf("Provider field is nil")
		}
		if *resp.Provider != "https://test.oidc.example.com" {
			return fmt.Errorf("expected Provider to be JWT iss value, got: %s", *resp.Provider)
		}
		return nil
	}))

	return results
}
