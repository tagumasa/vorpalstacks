package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (r *TestRunner) runSTSIdentityTests(tc *stsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sts", "GetCallerIdentity_ContentVerify", func() error {
		resp, err := tc.client.GetCallerIdentity(tc.ctx, &sts.GetCallerIdentityInput{})
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
		resp, err := tc.client.GetSessionToken(tc.ctx, &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int32(3600),
		})
		if err != nil {
			return err
		}
		return stsAssertCredentials(resp.Credentials)
	}))

	results = append(results, r.RunTest("sts", "GetSessionToken_RootDurationCap", func() error {
		// Root user requesting >3600 should be rejected.
		_, err := tc.client.GetSessionToken(tc.ctx, &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int32(7200),
		})
		return AssertErrorContains(err, "ValidationError")
	}))

	results = append(results, r.RunTest("sts", "GetSessionToken_InvalidDuration", func() error {
		_, err := tc.client.GetSessionToken(tc.ctx, &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int32(100),
		})
		if err := AssertErrorContains(err, "InvalidDuration"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
