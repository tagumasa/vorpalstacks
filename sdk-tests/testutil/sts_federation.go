package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (r *TestRunner) runSTSFederationTests(tc *stsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sts", "GetFederationToken_ContentVerify", func() error {
		resp, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name: aws.String("FederatedVerify"),
		})
		if err != nil {
			return err
		}
		if err := stsAssertCredentials(resp.Credentials); err != nil {
			return err
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
		resp, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name:   aws.String("FederatedPolicy"),
			Policy: aws.String(inlinePolicy),
		})
		if err != nil {
			return err
		}
		if err := stsAssertPackedPolicySize(resp.PackedPolicySize); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_WithDuration", func() error {
		resp, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name:            aws.String("FederatedDuration"),
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

	results = append(results, r.RunTest("sts", "GetFederationToken_InvalidName", func() error {
		// An empty name is rejected outright; a single-character name
		// fails the userNameType minimum length (2).  Both paths emit
		// InvalidInput.
		_, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name: aws.String(""),
		})
		if err := expectAWSErrorCode(err, "InvalidInput"); err != nil {
			return err
		}
		_, err = tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name: aws.String("x"),
		})
		return expectAWSErrorCode(err, "InvalidInput")
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_InvalidPolicy", func() error {
		_, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name:   aws.String("FederatedBadPolicy"),
			Policy: aws.String("not-valid-json"),
		})
		if err := AssertErrorContains(err, "MalformedPolicyDocument"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_InvalidDuration", func() error {
		_, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name:            aws.String("FederatedBadDuration"),
			DurationSeconds: aws.Int32(100),
		})
		if err := AssertErrorContains(err, "InvalidDuration"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetDelegatedAccessToken_ContentVerify", func() error {
		resp, err := tc.client.GetDelegatedAccessToken(tc.ctx, &sts.GetDelegatedAccessTokenInput{
			TradeInToken: aws.String("dummy-trade-in-token-verify"),
		})
		if err != nil {
			return err
		}
		if err := stsAssertCredentials(resp.Credentials); err != nil {
			return err
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
		_, err := tc.client.GetDelegatedAccessToken(tc.ctx, &sts.GetDelegatedAccessTokenInput{
			TradeInToken: aws.String(""),
		})
		if err := AssertErrorContains(err, "InvalidToken"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
